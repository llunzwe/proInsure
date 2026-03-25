package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/shared"
	"smartsure/pkg/logger"
)

// CommunicationService handles communication and notification operations
type CommunicationService struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewCommunicationService creates a new communication service
func NewCommunicationService(db *gorm.DB, logger *logger.Logger) *CommunicationService {
	return &CommunicationService{
		db:     db,
		logger: logger,
	}
}

// SendCommunication sends a communication to a recipient
func (s *CommunicationService) SendCommunication(ctx context.Context, comm *shared.Communication) error {
	s.logger.Info("Sending communication", "type", comm.Type, "recipient", comm.RecipientEmail)

	// Check recipient preferences
	canSend, err := s.checkPreferences(ctx, comm)
	if err != nil {
		s.logger.Warn("Failed to check preferences", "error", err)
		// Continue anyway for critical messages
		if comm.Priority != "urgent" && comm.Category != "transactional" {
			return fmt.Errorf("failed to verify preferences: %w", err)
		}
	}

	if !canSend {
		s.logger.Info("Communication blocked by user preferences", "recipient", comm.RecipientID)
		comm.Status = "blocked"
		return s.db.WithContext(ctx).Create(comm).Error
	}

	// Apply template if specified
	if comm.TemplateID != nil {
		if err := s.applyTemplate(ctx, comm); err != nil {
			s.logger.Error("Failed to apply template", "error", err)
		}
	}

	// Queue the communication
	comm.Status = "queued"
	if err := s.db.WithContext(ctx).Create(comm).Error; err != nil {
		return fmt.Errorf("failed to create communication: %w", err)
	}

	// Add to queue for processing
	queue := &shared.CommunicationQueue{
		CommunicationID: comm.ID,
		Priority:        s.getPriorityValue(comm.Priority),
		ScheduledFor:    time.Now(),
		Status:          "pending",
	}

	if comm.ScheduledFor != nil {
		queue.ScheduledFor = *comm.ScheduledFor
	}

	if err := s.db.WithContext(ctx).Create(queue).Error; err != nil {
		s.logger.Error("Failed to queue communication", "error", err)
	}

	// Process immediately if urgent
	if comm.Priority == "urgent" {
		go s.processCommunication(context.Background(), comm.ID)
	}

	return nil
}

// GetCommunication retrieves a communication by ID
func (s *CommunicationService) GetCommunication(ctx context.Context, commID uuid.UUID) (*shared.Communication, error) {
	var comm shared.Communication

	err := s.db.WithContext(ctx).
		Preload("Template").
		Preload("Events").
		First(&comm, "id = ?", commID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("communication not found")
		}
		return nil, fmt.Errorf("failed to get communication: %w", err)
	}

	return &comm, nil
}

// ListCommunications lists communications with filters
func (s *CommunicationService) ListCommunications(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*shared.Communication, int64, error) {
	query := s.db.WithContext(ctx).Model(&shared.Communication{})

	// Apply filters
	if recipientType, ok := filters["recipient_type"].(string); ok && recipientType != "" {
		query = query.Where("recipient_type = ?", recipientType)
	}
	if recipientID, ok := filters["recipient_id"].(uuid.UUID); ok {
		query = query.Where("recipient_id = ?", recipientID)
	}
	if commType, ok := filters["type"].(string); ok && commType != "" {
		query = query.Where("type = ?", commType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count communications: %w", err)
	}

	var communications []*shared.Communication
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&communications).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list communications: %w", err)
	}

	return communications, total, nil
}

// CreateTemplate creates a communication template
func (s *CommunicationService) CreateTemplate(ctx context.Context, template *shared.CommunicationTemplate) error {
	s.logger.Info("Creating communication template", "name", template.Name, "type", template.Type)

	// Check if template with same code exists
	var existing shared.CommunicationTemplate
	err := s.db.WithContext(ctx).Where("template_code = ?", template.TemplateCode).First(&existing).Error
	if err == nil {
		return fmt.Errorf("template with code %s already exists", template.TemplateCode)
	}

	// Create template
	if err := s.db.WithContext(ctx).Create(template).Error; err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	s.logger.Info("Template created successfully", "template_id", template.ID)
	return nil
}

// GetTemplate retrieves a template by ID
func (s *CommunicationService) GetTemplate(ctx context.Context, templateID uuid.UUID) (*shared.CommunicationTemplate, error) {
	var template shared.CommunicationTemplate

	err := s.db.WithContext(ctx).
		Preload("Translations").
		First(&template, "id = ?", templateID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &template, nil
}

// UpdatePreferences updates notification preferences
func (s *CommunicationService) UpdatePreferences(ctx context.Context, entityType string, entityID uuid.UUID, updates map[string]interface{}) error {
	s.logger.Info("Updating notification preferences", "entity_type", entityType, "entity_id", entityID)

	var prefs shared.NotificationPreference
	err := s.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&prefs).Error

	if err == gorm.ErrRecordNotFound {
		// Create new preferences
		prefs = shared.NotificationPreference{
			EntityType: entityType,
			EntityID:   entityID,
		}
		if err := s.db.WithContext(ctx).Create(&prefs).Error; err != nil {
			return fmt.Errorf("failed to create preferences: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to get preferences: %w", err)
	}

	// Update preferences
	if err := s.db.WithContext(ctx).Model(&prefs).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

// GetPreferences gets notification preferences
func (s *CommunicationService) GetPreferences(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.NotificationPreference, error) {
	var prefs shared.NotificationPreference

	err := s.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		First(&prefs).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default preferences
			return &shared.NotificationPreference{
				EntityType:          entityType,
				EntityID:            entityID,
				EmailEnabled:        true,
				SMSEnabled:          true,
				PushEnabled:         true,
				TransactionalAlerts: true,
			}, nil
		}
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	return &prefs, nil
}

// CreateCampaign creates a communication campaign
func (s *CommunicationService) CreateCampaign(ctx context.Context, campaign *shared.CommunicationCampaign) error {
	s.logger.Info("Creating communication campaign", "name", campaign.Name, "type", campaign.Type)

	// Validate campaign
	if campaign.TemplateID == uuid.Nil {
		return fmt.Errorf("template is required for campaign")
	}

	// Estimate recipients based on criteria
	recipients, err := s.estimateRecipients(ctx, campaign.TargetType, campaign.TargetCriteria)
	if err != nil {
		s.logger.Warn("Failed to estimate recipients", "error", err)
	}
	campaign.EstimatedRecipients = recipients

	// Create campaign
	if err := s.db.WithContext(ctx).Create(campaign).Error; err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	s.logger.Info("Campaign created successfully", "campaign_id", campaign.ID)
	return nil
}

// StartCampaign starts a scheduled campaign
func (s *CommunicationService) StartCampaign(ctx context.Context, campaignID uuid.UUID) error {
	s.logger.Info("Starting campaign", "campaign_id", campaignID)

	return s.db.Transaction(func(tx *gorm.DB) error {
		var campaign shared.CommunicationCampaign
		if err := tx.WithContext(ctx).First(&campaign, "id = ?", campaignID).Error; err != nil {
			return fmt.Errorf("campaign not found: %w", err)
		}

		if campaign.Status != "scheduled" && campaign.Status != "draft" {
			return fmt.Errorf("campaign cannot be started in current status: %s", campaign.Status)
		}

		// Start the campaign
		campaign.StartCampaign()
		if err := tx.WithContext(ctx).Save(&campaign).Error; err != nil {
			return fmt.Errorf("failed to update campaign: %w", err)
		}

		// Queue campaign communications
		go s.processCampaign(context.Background(), campaignID)

		return nil
	})
}

// ProcessQueue processes pending communications in the queue
func (s *CommunicationService) ProcessQueue(ctx context.Context) error {
	var queueItems []shared.CommunicationQueue

	// Get pending items that should be processed
	err := s.db.WithContext(ctx).
		Where("status = ? AND scheduled_for <= ?", "pending", time.Now()).
		Order("priority ASC, scheduled_for ASC").
		Limit(100).
		Find(&queueItems).Error

	if err != nil {
		return fmt.Errorf("failed to get queue items: %w", err)
	}

	for _, item := range queueItems {
		if item.ShouldProcess() {
			go s.processCommunication(context.Background(), item.CommunicationID)
		}
	}

	return nil
}

// TrackEvent tracks a communication event
func (s *CommunicationService) TrackEvent(ctx context.Context, commID uuid.UUID, eventType string, eventData map[string]interface{}) error {
	eventDataJSON, _ := json.Marshal(eventData)

	event := &shared.CommunicationEvent{
		CommunicationID: commID,
		EventType:       eventType,
		EventData:       string(eventDataJSON),
		IPAddress:       s.extractIPAddress(eventData),
		UserAgent:       s.extractUserAgent(eventData),
	}

	if err := s.db.WithContext(ctx).Create(event).Error; err != nil {
		return fmt.Errorf("failed to track event: %w", err)
	}

	// Update communication status based on event
	go s.updateCommunicationStatus(context.Background(), commID, eventType)

	return nil
}

// Helper functions

func (s *CommunicationService) checkPreferences(ctx context.Context, comm *shared.Communication) (bool, error) {
	prefs, err := s.GetPreferences(ctx, comm.RecipientType, comm.RecipientID)
	if err != nil {
		return true, err // Allow if can't check preferences
	}

	return prefs.CanSendNow(comm.Type, comm.Category), nil
}

func (s *CommunicationService) applyTemplate(ctx context.Context, comm *shared.Communication) error {
	template, err := s.GetTemplate(ctx, *comm.TemplateID)
	if err != nil {
		return err
	}

	// Apply template content
	comm.Subject = template.Subject
	comm.Body = s.replaceVariables(template.Body, comm.Variables)
	comm.HTMLBody = s.replaceVariables(template.HTMLBody, comm.Variables)

	// Update template usage
	template.UseTemplate()
	s.db.WithContext(ctx).Save(template)

	return nil
}

func (s *CommunicationService) replaceVariables(content string, variables string) string {
	if variables == "" {
		return content
	}

	var vars map[string]interface{}
	if err := json.Unmarshal([]byte(variables), &vars); err != nil {
		return content
	}

	// Simple variable replacement - in production, use a template engine
	result := content
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = fmt.Sprintf(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

func (s *CommunicationService) getPriorityValue(priority string) int {
	switch priority {
	case "urgent":
		return 1
	case "high":
		return 3
	case "normal":
		return 5
	case "low":
		return 8
	default:
		return 5
	}
}

func (s *CommunicationService) processCommunication(ctx context.Context, commID uuid.UUID) {
	comm, err := s.GetCommunication(ctx, commID)
	if err != nil {
		s.logger.Error("Failed to get communication for processing", "error", err, "comm_id", commID)
		return
	}

	// Send via appropriate channel
	switch comm.Type {
	case "email":
		err = s.sendEmail(ctx, comm)
	case "sms":
		err = s.sendSMS(ctx, comm)
	case "push":
		err = s.sendPush(ctx, comm)
	case "in_app":
		err = s.sendInApp(ctx, comm)
	default:
		err = fmt.Errorf("unsupported communication type: %s", comm.Type)
	}

	if err != nil {
		comm.MarkFailed(err.Error())
		s.db.WithContext(ctx).Save(comm)

		// Retry if needed
		if comm.ShouldRetry() {
			s.scheduleRetry(ctx, commID)
		}
	} else {
		comm.Send()
		s.db.WithContext(ctx).Save(comm)
	}
}

func (s *CommunicationService) sendEmail(ctx context.Context, comm *shared.Communication) error {
	// Integration with email service provider
	s.logger.Info("Sending email", "to", comm.RecipientEmail, "subject", comm.Subject)
	// Actual implementation would use SendGrid, AWS SES, etc.
	return nil
}

func (s *CommunicationService) sendSMS(ctx context.Context, comm *shared.Communication) error {
	// Integration with SMS provider
	s.logger.Info("Sending SMS", "to", comm.RecipientPhone)
	// Actual implementation would use Twilio, MessageBird, etc.
	return nil
}

func (s *CommunicationService) sendPush(ctx context.Context, comm *shared.Communication) error {
	// Integration with push notification service
	s.logger.Info("Sending push notification", "to", comm.RecipientID)
	// Actual implementation would use Firebase, OneSignal, etc.
	return nil
}

func (s *CommunicationService) sendInApp(ctx context.Context, comm *shared.Communication) error {
	// Store in-app notification
	s.logger.Info("Creating in-app notification", "for", comm.RecipientID)
	comm.Status = "delivered"
	return nil
}

func (s *CommunicationService) scheduleRetry(ctx context.Context, commID uuid.UUID) {
	var queue shared.CommunicationQueue
	if err := s.db.WithContext(ctx).Where("communication_id = ?", commID).First(&queue).Error; err != nil {
		s.logger.Error("Failed to get queue item for retry", "error", err)
		return
	}

	queue.IncrementRetry()
	s.db.WithContext(ctx).Save(&queue)
}

func (s *CommunicationService) processCampaign(ctx context.Context, campaignID uuid.UUID) {
	s.logger.Info("Processing campaign", "campaign_id", campaignID)

	var campaign shared.CommunicationCampaign
	if err := s.db.WithContext(ctx).First(&campaign, "id = ?", campaignID).Error; err != nil {
		s.logger.Error("Failed to get campaign", "error", err)
		return
	}

	// Get target recipients
	recipients, err := s.getCampaignRecipients(ctx, &campaign)
	if err != nil {
		s.logger.Error("Failed to get campaign recipients", "error", err)
		return
	}

	// Create communications for each recipient
	for _, recipient := range recipients {
		comm := &shared.Communication{
			RecipientType:     recipient.Type,
			RecipientID:       recipient.ID,
			RecipientEmail:    recipient.Email,
			RecipientPhone:    recipient.Phone,
			RecipientName:     recipient.Name,
			Type:              campaign.Type,
			Category:          campaign.Category,
			Priority:          "normal",
			TemplateID:        &campaign.TemplateID,
			Subject:           campaign.Subject,
			Body:              campaign.Content,
			TriggerType:       "campaign",
			TriggerID:         &campaign.ID,
			RelatedEntityType: "campaign",
			RelatedEntityID:   &campaign.ID,
		}

		if err := s.SendCommunication(ctx, comm); err != nil {
			s.logger.Error("Failed to send campaign communication", "error", err)
			campaign.TotalFailed++
		} else {
			campaign.TotalSent++
		}
	}

	// Complete campaign
	campaign.CompleteCampaign()
	s.db.WithContext(ctx).Save(&campaign)
}

func (s *CommunicationService) estimateRecipients(ctx context.Context, targetType string, criteria string) (int, error) {
	// Parse criteria and estimate recipient count
	// Simplified implementation
	switch targetType {
	case "all":
		var count int64
		s.db.WithContext(ctx).Model(&models.User{}).Count(&count)
		return int(count), nil
	case "segment":
		// Parse criteria JSON and apply filters
		return 0, nil
	default:
		return 0, nil
	}
}

type CampaignRecipient struct {
	ID    uuid.UUID
	Type  string
	Email string
	Phone string
	Name  string
}

func (s *CommunicationService) getCampaignRecipients(ctx context.Context, campaign *shared.CommunicationCampaign) ([]CampaignRecipient, error) {
	var recipients []CampaignRecipient

	switch campaign.TargetType {
	case "all":
		var users []models.User
		s.db.WithContext(ctx).Find(&users)
		for _, user := range users {
			recipients = append(recipients, CampaignRecipient{
				ID:    user.ID,
				Type:  "user",
				Email: user.Email,
				Phone: user.PhoneNumber,
				Name:  user.FirstName + " " + user.LastName,
			})
		}
	case "specific":
		// Parse target list
		var targetIDs []string
		json.Unmarshal([]byte(campaign.TargetList), &targetIDs)
		// Fetch specific users
	}

	return recipients, nil
}

func (s *CommunicationService) updateCommunicationStatus(ctx context.Context, commID uuid.UUID, eventType string) {
	var comm shared.Communication
	if err := s.db.WithContext(ctx).First(&comm, "id = ?", commID).Error; err != nil {
		return
	}

	switch eventType {
	case "delivered":
		comm.MarkDelivered()
	case "opened":
		comm.MarkOpened()
	case "clicked":
		comm.MarkClicked()
	case "bounced", "failed":
		comm.MarkFailed("Event: " + eventType)
	}

	s.db.WithContext(ctx).Save(&comm)
}

func (s *CommunicationService) extractIPAddress(data map[string]interface{}) string {
	if ip, ok := data["ip_address"].(string); ok {
		return ip
	}
	return ""
}

func (s *CommunicationService) extractUserAgent(data map[string]interface{}) string {
	if ua, ok := data["user_agent"].(string); ok {
		return ua
	}
	return ""
}
