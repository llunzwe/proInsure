package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"smartsure/pkg/errors"
	"smartsure/pkg/logger"
)

// WebhookHandler handles webhook subscriptions and deliveries
type WebhookHandler struct {
	db            *gorm.DB
	logger        *logger.Logger
	deliveryQueue chan *WebhookDelivery
	workers       int
	wg            sync.WaitGroup
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(db *gorm.DB, logger *logger.Logger) *WebhookHandler {
	h := &WebhookHandler{
		db:            db,
		logger:        logger,
		deliveryQueue: make(chan *WebhookDelivery, 10000),
		workers:       10,
	}

	// Start webhook delivery workers
	h.startDeliveryWorkers()

	return h
}

// WebhookSubscription represents a webhook subscription
type WebhookSubscription struct {
	ID           uuid.UUID              `json:"id" gorm:"type:uuid;primary_key"`
	Name         string                 `json:"name" binding:"required"`
	URL          string                 `json:"url" binding:"required,url"`
	Secret       string                 `json:"secret,omitempty"`                // For HMAC signature
	Events       []string               `json:"events" binding:"required,min=1"` // Event types to subscribe to
	Active       bool                   `json:"active"`
	Headers      map[string]string      `json:"headers,omitempty" gorm:"type:jsonb"`
	MaxRetries   int                    `json:"max_retries"`
	RetryDelay   int                    `json:"retry_delay"` // Seconds
	Metadata     map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	LastDelivery *time.Time             `json:"last_delivery,omitempty"`
	FailureCount int                    `json:"failure_count"`
	SuccessCount int                    `json:"success_count"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID             uuid.UUID              `json:"id" gorm:"type:uuid;primary_key"`
	SubscriptionID uuid.UUID              `json:"subscription_id" gorm:"type:uuid"`
	EventType      string                 `json:"event_type"`
	Payload        map[string]interface{} `json:"payload" gorm:"type:jsonb"`
	Status         string                 `json:"status"` // pending, success, failed
	Attempts       int                    `json:"attempts"`
	LastAttempt    *time.Time             `json:"last_attempt,omitempty"`
	NextRetry      *time.Time             `json:"next_retry,omitempty"`
	ResponseCode   int                    `json:"response_code,omitempty"`
	ResponseBody   string                 `json:"response_body,omitempty"`
	Error          string                 `json:"error,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	DeliveredAt    *time.Time             `json:"delivered_at,omitempty"`
}

// WebhookEvent represents an event that can trigger webhooks
type WebhookEvent struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Source    string                 `json:"source"`
	Version   string                 `json:"version"`
}

// CreateSubscription creates a new webhook subscription
// @Summary Create webhook subscription
// @Description Subscribe to receive webhook notifications for specified events
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param subscription body WebhookSubscription true "Webhook subscription"
// @Success 201 {object} WebhookSubscription
// @Router /api/v2/webhooks/subscriptions [post]
func (h *WebhookHandler) CreateSubscription(c *gin.Context) {
	var subscription WebhookSubscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		errors.ValidationFailed(c, "Invalid webhook subscription", map[string][]string{
			"subscription": {err.Error()},
		})
		return
	}

	// Validate events
	validEvents := h.getValidEvents()
	for _, event := range subscription.Events {
		if !contains(validEvents, event) {
			errors.ValidationFailed(c, "Invalid event type", map[string][]string{
				"events": {fmt.Sprintf("Invalid event type: %s", event)},
			})
			return
		}
	}

	// Generate secret if not provided
	if subscription.Secret == "" {
		subscription.Secret = h.generateSecret()
	}

	// Set defaults
	subscription.ID = uuid.New()
	subscription.Active = true
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()
	if subscription.MaxRetries == 0 {
		subscription.MaxRetries = 3
	}
	if subscription.RetryDelay == 0 {
		subscription.RetryDelay = 60
	}

	// Save to database
	if err := h.db.Create(&subscription).Error; err != nil {
		errors.InternalServerError(c, err)
		return
	}

	// Test webhook with ping event
	go h.sendTestWebhook(&subscription)

	c.JSON(http.StatusCreated, subscription)
}

// ListSubscriptions lists all webhook subscriptions
// @Summary List webhook subscriptions
// @Description Get a list of all webhook subscriptions
// @Tags Webhooks
// @Produce json
// @Param active query bool false "Filter by active status"
// @Param event query string false "Filter by event type"
// @Success 200 {array} WebhookSubscription
// @Router /api/v2/webhooks/subscriptions [get]
func (h *WebhookHandler) ListSubscriptions(c *gin.Context) {
	var subscriptions []WebhookSubscription
	query := h.db

	// Apply filters
	if active := c.Query("active"); active != "" {
		query = query.Where("active = ?", active == "true")
	}

	if event := c.Query("event"); event != "" {
		query = query.Where("? = ANY(events)", event)
	}

	if err := query.Find(&subscriptions).Error; err != nil {
		errors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// GetSubscription gets a specific webhook subscription
// @Summary Get webhook subscription
// @Description Get details of a specific webhook subscription
// @Tags Webhooks
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} WebhookSubscription
// @Router /api/v2/webhooks/subscriptions/{id} [get]
func (h *WebhookHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")

	var subscription WebhookSubscription
	if err := h.db.First(&subscription, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.ResourceNotFound(c, "WebhookSubscription", id)
		} else {
			errors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// UpdateSubscription updates a webhook subscription
// @Summary Update webhook subscription
// @Description Update an existing webhook subscription
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body WebhookSubscription true "Updated subscription"
// @Success 200 {object} WebhookSubscription
// @Router /api/v2/webhooks/subscriptions/{id} [put]
func (h *WebhookHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")

	var subscription WebhookSubscription
	if err := h.db.First(&subscription, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.ResourceNotFound(c, "WebhookSubscription", id)
		} else {
			errors.InternalServerError(c, err)
		}
		return
	}

	var updates WebhookSubscription
	if err := c.ShouldBindJSON(&updates); err != nil {
		errors.ValidationFailed(c, "Invalid webhook subscription", map[string][]string{
			"subscription": {err.Error()},
		})
		return
	}

	// Update fields
	subscription.Name = updates.Name
	subscription.URL = updates.URL
	subscription.Events = updates.Events
	subscription.Active = updates.Active
	subscription.Headers = updates.Headers
	subscription.MaxRetries = updates.MaxRetries
	subscription.RetryDelay = updates.RetryDelay
	subscription.Metadata = updates.Metadata
	subscription.UpdatedAt = time.Now()

	if err := h.db.Save(&subscription).Error; err != nil {
		errors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription deletes a webhook subscription
// @Summary Delete webhook subscription
// @Description Delete a webhook subscription
// @Tags Webhooks
// @Param id path string true "Subscription ID"
// @Success 204
// @Router /api/v2/webhooks/subscriptions/{id} [delete]
func (h *WebhookHandler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	result := h.db.Delete(&WebhookSubscription{}, "id = ?", id)
	if result.Error != nil {
		errors.InternalServerError(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		errors.ResourceNotFound(c, "WebhookSubscription", id)
		return
	}

	c.Status(http.StatusNoContent)
}

// TestWebhook tests a webhook subscription
// @Summary Test webhook
// @Description Send a test webhook to verify the endpoint is working
// @Tags Webhooks
// @Param id path string true "Subscription ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/webhooks/subscriptions/{id}/test [post]
func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	id := c.Param("id")

	var subscription WebhookSubscription
	if err := h.db.First(&subscription, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.ResourceNotFound(c, "WebhookSubscription", id)
		} else {
			errors.InternalServerError(c, err)
		}
		return
	}

	// Send test webhook
	testEvent := &WebhookEvent{
		Type:      "webhook.test",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message":         "This is a test webhook",
			"subscription_id": subscription.ID,
		},
		Source:  "api",
		Version: "v2",
	}

	delivery := h.createDelivery(&subscription, testEvent)
	result := h.deliverWebhook(&subscription, delivery)

	c.JSON(http.StatusOK, gin.H{
		"success":       result.Status == "success",
		"status_code":   result.ResponseCode,
		"response_time": result.LastAttempt.Sub(result.CreatedAt).Milliseconds(),
		"message": func() string {
			if result.Status == "success" {
				return "Webhook delivered successfully"
			}
			return result.Error
		}(),
	})
}

// GetDeliveryHistory gets webhook delivery history
// @Summary Get webhook delivery history
// @Description Get delivery history for a webhook subscription
// @Tags Webhooks
// @Produce json
// @Param id path string true "Subscription ID"
// @Param status query string false "Filter by status (pending, success, failed)"
// @Success 200 {array} WebhookDelivery
// @Router /api/v2/webhooks/subscriptions/{id}/deliveries [get]
func (h *WebhookHandler) GetDeliveryHistory(c *gin.Context) {
	id := c.Param("id")

	var deliveries []WebhookDelivery
	query := h.db.Where("subscription_id = ?", id)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC").Limit(100)

	if err := query.Find(&deliveries).Error; err != nil {
		errors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

// RetryDelivery retries a failed webhook delivery
// @Summary Retry webhook delivery
// @Description Retry a failed webhook delivery
// @Tags Webhooks
// @Param delivery_id path string true "Delivery ID"
// @Success 200 {object} WebhookDelivery
// @Router /api/v2/webhooks/deliveries/{delivery_id}/retry [post]
func (h *WebhookHandler) RetryDelivery(c *gin.Context) {
	deliveryID := c.Param("delivery_id")

	var delivery WebhookDelivery
	if err := h.db.First(&delivery, "id = ?", deliveryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.ResourceNotFound(c, "WebhookDelivery", deliveryID)
		} else {
			errors.InternalServerError(c, err)
		}
		return
	}

	if delivery.Status == "success" {
		errors.NewProblemDetails(
			"https://api.smartsure.com/errors/already-delivered",
			"Already Delivered",
			http.StatusBadRequest,
			"This webhook has already been delivered successfully",
			c.Request.URL.Path,
		).Send(c)
		return
	}

	// Get subscription
	var subscription WebhookSubscription
	if err := h.db.First(&subscription, "id = ?", delivery.SubscriptionID).Error; err != nil {
		errors.InternalServerError(c, err)
		return
	}

	// Retry delivery
	result := h.deliverWebhook(&subscription, &delivery)

	c.JSON(http.StatusOK, result)
}

// GetWebhookEvents gets available webhook events
// @Summary Get available webhook events
// @Description Get a list of all available webhook event types
// @Tags Webhooks
// @Produce json
// @Success 200 {object} map[string][]string
// @Router /api/v2/webhooks/events [get]
func (h *WebhookHandler) GetWebhookEvents(c *gin.Context) {
	events := map[string][]string{
		"device": {
			"device.created",
			"device.updated",
			"device.deleted",
			"device.verified",
			"device.stolen",
		},
		"policy": {
			"policy.created",
			"policy.updated",
			"policy.cancelled",
			"policy.expired",
			"policy.renewed",
		},
		"claim": {
			"claim.created",
			"claim.updated",
			"claim.approved",
			"claim.rejected",
			"claim.paid",
		},
		"payment": {
			"payment.created",
			"payment.succeeded",
			"payment.failed",
			"payment.refunded",
		},
		"subscription": {
			"subscription.created",
			"subscription.updated",
			"subscription.cancelled",
			"subscription.renewed",
		},
		"user": {
			"user.created",
			"user.updated",
			"user.deleted",
			"user.suspended",
			"user.activated",
		},
		"fraud": {
			"fraud.detected",
			"fraud.cleared",
			"fraud.investigated",
		},
		"webhook": {
			"webhook.test",
			"webhook.ping",
		},
	}

	c.JSON(http.StatusOK, events)
}

// TriggerEvent triggers a webhook event (internal use)
func (h *WebhookHandler) TriggerEvent(eventType string, data map[string]interface{}) {
	event := &WebhookEvent{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
		Source:    "api",
		Version:   "v2",
	}

	// Find all active subscriptions for this event
	var subscriptions []WebhookSubscription
	h.db.Where("active = ? AND ? = ANY(events)", true, eventType).Find(&subscriptions)

	// Queue deliveries for each subscription
	for _, subscription := range subscriptions {
		delivery := h.createDelivery(&subscription, event)
		h.deliveryQueue <- delivery
	}
}

// Private methods

func (h *WebhookHandler) startDeliveryWorkers() {
	for i := 0; i < h.workers; i++ {
		h.wg.Add(1)
		go h.deliveryWorker()
	}
}

func (h *WebhookHandler) deliveryWorker() {
	defer h.wg.Done()

	for delivery := range h.deliveryQueue {
		var subscription WebhookSubscription
		if err := h.db.First(&subscription, "id = ?", delivery.SubscriptionID).Error; err != nil {
			h.logger.Error("Failed to find subscription", "error", err, "subscription_id", delivery.SubscriptionID)
			continue
		}

		h.deliverWebhook(&subscription, delivery)
	}
}

func (h *WebhookHandler) createDelivery(subscription *WebhookSubscription, event *WebhookEvent) *WebhookDelivery {
	delivery := &WebhookDelivery{
		ID:             uuid.New(),
		SubscriptionID: subscription.ID,
		EventType:      event.Type,
		Payload: map[string]interface{}{
			"id":        uuid.New().String(),
			"type":      event.Type,
			"timestamp": event.Timestamp,
			"data":      event.Data,
			"source":    event.Source,
			"version":   event.Version,
		},
		Status:    "pending",
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	h.db.Create(delivery)
	return delivery
}

func (h *WebhookHandler) deliverWebhook(subscription *WebhookSubscription, delivery *WebhookDelivery) *WebhookDelivery {
	delivery.Attempts++
	now := time.Now()
	delivery.LastAttempt = &now

	// Prepare request
	payload, _ := json.Marshal(delivery.Payload)
	req, err := http.NewRequest("POST", subscription.URL, bytes.NewBuffer(payload))
	if err != nil {
		delivery.Status = "failed"
		delivery.Error = err.Error()
		h.db.Save(delivery)
		return delivery
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-ID", delivery.ID.String())
	req.Header.Set("X-Webhook-Event", delivery.EventType)
	req.Header.Set("X-Webhook-Timestamp", delivery.CreatedAt.Format(time.RFC3339))

	// Add custom headers
	for key, value := range subscription.Headers {
		req.Header.Set(key, value)
	}

	// Add signature if secret is configured
	if subscription.Secret != "" {
		signature := h.generateSignature(subscription.Secret, payload)
		req.Header.Set("X-Webhook-Signature", signature)
	}

	// Send request with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		delivery.Status = "failed"
		delivery.Error = err.Error()

		// Schedule retry if attempts < max retries
		if delivery.Attempts < subscription.MaxRetries {
			retryTime := time.Now().Add(time.Duration(subscription.RetryDelay) * time.Second)
			delivery.NextRetry = &retryTime

			// Re-queue for retry
			go func() {
				time.Sleep(time.Duration(subscription.RetryDelay) * time.Second)
				h.deliveryQueue <- delivery
			}()
		}
	} else {
		defer resp.Body.Close()
		delivery.ResponseCode = resp.StatusCode

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			delivery.Status = "success"
			deliveredAt := time.Now()
			delivery.DeliveredAt = &deliveredAt

			// Update subscription stats
			subscription.SuccessCount++
			subscription.LastDelivery = &deliveredAt
			h.db.Save(subscription)
		} else {
			delivery.Status = "failed"
			delivery.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)

			// Schedule retry if needed
			if delivery.Attempts < subscription.MaxRetries {
				retryTime := time.Now().Add(time.Duration(subscription.RetryDelay) * time.Second)
				delivery.NextRetry = &retryTime

				go func() {
					time.Sleep(time.Duration(subscription.RetryDelay) * time.Second)
					h.deliveryQueue <- delivery
				}()
			} else {
				// Update failure count
				subscription.FailureCount++
				h.db.Save(subscription)
			}
		}
	}

	h.db.Save(delivery)
	return delivery
}

func (h *WebhookHandler) generateSignature(secret string, payload []byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(payload)
	return "sha256=" + hex.EncodeToString(hash.Sum(nil))
}

func (h *WebhookHandler) generateSecret() string {
	return uuid.New().String()
}

func (h *WebhookHandler) getValidEvents() []string {
	return []string{
		"device.created", "device.updated", "device.deleted", "device.verified", "device.stolen",
		"policy.created", "policy.updated", "policy.cancelled", "policy.expired", "policy.renewed",
		"claim.created", "claim.updated", "claim.approved", "claim.rejected", "claim.paid",
		"payment.created", "payment.succeeded", "payment.failed", "payment.refunded",
		"subscription.created", "subscription.updated", "subscription.cancelled", "subscription.renewed",
		"user.created", "user.updated", "user.deleted", "user.suspended", "user.activated",
		"fraud.detected", "fraud.cleared", "fraud.investigated",
		"webhook.test", "webhook.ping",
	}
}

func (h *WebhookHandler) sendTestWebhook(subscription *WebhookSubscription) {
	event := &WebhookEvent{
		Type:      "webhook.ping",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message":         "Webhook subscription created successfully",
			"subscription_id": subscription.ID,
		},
		Source:  "api",
		Version: "v2",
	}

	delivery := h.createDelivery(subscription, event)
	h.deliverWebhook(subscription, delivery)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
