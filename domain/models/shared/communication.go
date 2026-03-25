package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Communication represents a unified communication record for all entity types
type Communication struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CommunicationCode string    `gorm:"uniqueIndex;not null" json:"communication_code"`

	// Recipient Information (polymorphic - can be any entity)
	RecipientType  string    `gorm:"not null" json:"recipient_type"` // user, partner, customer, agent, vendor
	RecipientID    uuid.UUID `gorm:"type:uuid;not null" json:"recipient_id"`
	RecipientEmail string    `json:"recipient_email"`
	RecipientPhone string    `json:"recipient_phone"`
	RecipientName  string    `json:"recipient_name"`

	// Communication Type
	Type     string `gorm:"not null" json:"type"`             // email, sms, push, in_app, whatsapp, call
	Category string `gorm:"not null" json:"category"`         // transactional, promotional, alert, reminder, notification
	Priority string `gorm:"default:'normal'" json:"priority"` // low, normal, high, urgent

	// Template
	TemplateID   *uuid.UUID `gorm:"type:uuid" json:"template_id"`
	TemplateName string     `json:"template_name"`
	Language     string     `gorm:"default:'en'" json:"language"`

	// Content
	Subject     string `json:"subject"`
	Body        string `gorm:"type:text" json:"body"`
	HTMLBody    string `gorm:"type:text" json:"html_body"`
	Variables   string `gorm:"type:json" json:"variables"`   // JSON object of template variables
	Attachments string `gorm:"type:json" json:"attachments"` // JSON array of attachment URLs

	// Context (what triggered this communication)
	TriggerType       string     `json:"trigger_type"`                // claim_update, payment_reminder, policy_renewal, etc.
	TriggerID         *uuid.UUID `gorm:"type:uuid" json:"trigger_id"` // ID of the triggering entity
	RelatedEntityType string     `json:"related_entity_type"`         // policy, claim, payment, etc.
	RelatedEntityID   *uuid.UUID `gorm:"type:uuid" json:"related_entity_id"`

	// Scheduling
	ScheduledFor *time.Time `json:"scheduled_for"`
	SendAfter    *time.Time `json:"send_after"`
	ExpiresAt    *time.Time `json:"expires_at"`

	// Status & Tracking
	Status      string     `gorm:"default:'pending'" json:"status"` // pending, queued, sending, sent, delivered, failed, cancelled
	SentAt      *time.Time `json:"sent_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	OpenedAt    *time.Time `json:"opened_at"`
	ClickedAt   *time.Time `json:"clicked_at"`

	// Delivery Details
	Provider          string     `json:"provider"` // sendgrid, twilio, firebase, internal
	ProviderMessageID string     `json:"provider_message_id"`
	DeliveryAttempts  int        `gorm:"default:0" json:"delivery_attempts"`
	LastAttemptAt     *time.Time `json:"last_attempt_at"`
	FailureReason     string     `json:"failure_reason"`

	// Response Tracking
	ResponseRequired   bool       `gorm:"default:false" json:"response_required"`
	ResponseDeadline   *time.Time `json:"response_deadline"`
	ResponseReceived   bool       `gorm:"default:false" json:"response_received"`
	ResponseReceivedAt *time.Time `json:"response_received_at"`
	ResponseContent    string     `json:"response_content"`

	// Preferences & Compliance
	UnsubscribeToken string `json:"unsubscribe_token"`
	ConsentVerified  bool   `gorm:"default:false" json:"consent_verified"`
	OptOutRequested  bool   `gorm:"default:false" json:"opt_out_requested"`

	// Metadata
	Tags     string `gorm:"type:json" json:"tags"` // JSON array
	Metadata string `gorm:"type:json" json:"metadata"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Template *CommunicationTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Events   []CommunicationEvent   `gorm:"foreignKey:CommunicationID" json:"events,omitempty"`
}

// CommunicationTemplate represents reusable message templates
type CommunicationTemplate struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TemplateCode string    `gorm:"uniqueIndex;not null" json:"template_code"`
	Name         string    `gorm:"not null" json:"name"`

	// Template Type
	Type     string `gorm:"not null" json:"type"`     // email, sms, push, in_app
	Category string `gorm:"not null" json:"category"` // claim, payment, policy, general
	Purpose  string `json:"purpose"`

	// Content
	Subject  string `json:"subject"`
	Body     string `gorm:"type:text;not null" json:"body"`
	HTMLBody string `gorm:"type:text" json:"html_body"`

	// Variables that can be used in template
	Variables         string `gorm:"type:json" json:"variables"`          // JSON schema of available variables
	RequiredVariables string `gorm:"type:json" json:"required_variables"` // JSON array of required variable names

	// Multi-language Support
	Language         string     `gorm:"default:'en'" json:"language"`
	IsDefault        bool       `gorm:"default:false" json:"is_default"`
	ParentTemplateID *uuid.UUID `gorm:"type:uuid" json:"parent_template_id"` // For language variants

	// Settings
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	RequiresApproval bool       `gorm:"default:false" json:"requires_approval"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at"`

	// Tracking
	UseCount    int        `gorm:"default:0" json:"use_count"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	SuccessRate float64    `gorm:"default:0" json:"success_rate"`

	// Version Control
	Version      int    `gorm:"default:1" json:"version"`
	VersionNotes string `json:"version_notes"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Translations []CommunicationTemplate `gorm:"foreignKey:ParentTemplateID" json:"translations,omitempty"`
}

// CommunicationEvent tracks events related to communications
type CommunicationEvent struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CommunicationID uuid.UUID `gorm:"type:uuid;not null" json:"communication_id"`
	EventType       string    `gorm:"not null" json:"event_type"` // queued, sent, delivered, opened, clicked, bounced, failed
	EventData       string    `gorm:"type:json" json:"event_data"`
	IPAddress       string    `json:"ip_address"`
	UserAgent       string    `json:"user_agent"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Communication Communication `gorm:"foreignKey:CommunicationID" json:"communication,omitempty"`
}

// NotificationPreference manages user/partner communication preferences
type NotificationPreference struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// Entity (can be user, partner, etc.)
	EntityType string    `gorm:"not null" json:"entity_type"` // user, partner, customer
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`

	// Channel Preferences
	EmailEnabled    bool `gorm:"default:true" json:"email_enabled"`
	SMSEnabled      bool `gorm:"default:true" json:"sms_enabled"`
	PushEnabled     bool `gorm:"default:true" json:"push_enabled"`
	InAppEnabled    bool `gorm:"default:true" json:"in_app_enabled"`
	WhatsAppEnabled bool `gorm:"default:false" json:"whatsapp_enabled"`

	// Category Preferences
	TransactionalAlerts bool `gorm:"default:true" json:"transactional_alerts"`
	PromotionalMessages bool `gorm:"default:false" json:"promotional_messages"`
	PolicyUpdates       bool `gorm:"default:true" json:"policy_updates"`
	ClaimUpdates        bool `gorm:"default:true" json:"claim_updates"`
	PaymentReminders    bool `gorm:"default:true" json:"payment_reminders"`
	SecurityAlerts      bool `gorm:"default:true" json:"security_alerts"`

	// Timing Preferences
	PreferredTimeZone string `gorm:"default:'UTC'" json:"preferred_timezone"`
	QuietHoursStart   string `json:"quiet_hours_start"`               // HH:MM format
	QuietHoursEnd     string `json:"quiet_hours_end"`                 // HH:MM format
	PreferredDays     string `gorm:"type:json" json:"preferred_days"` // JSON array of days

	// Frequency Control
	MaxEmailsPerDay int    `gorm:"default:10" json:"max_emails_per_day"`
	MaxSMSPerDay    int    `gorm:"default:5" json:"max_sms_per_day"`
	MaxPushPerDay   int    `gorm:"default:20" json:"max_push_per_day"`
	DigestFrequency string `gorm:"default:'immediate'" json:"digest_frequency"` // immediate, daily, weekly

	// Language & Format
	PreferredLanguage string `gorm:"default:'en'" json:"preferred_language"`
	PreferredFormat   string `gorm:"default:'html'" json:"preferred_format"` // text, html

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CommunicationCampaign represents bulk communication campaigns
type CommunicationCampaign struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignCode string    `gorm:"uniqueIndex;not null" json:"campaign_code"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`

	// Campaign Type
	Type     string `gorm:"not null" json:"type"` // email, sms, push, multi_channel
	Category string `json:"category"`             // promotional, informational, alert

	// Target Audience
	TargetType          string `gorm:"not null" json:"target_type"`      // all, segment, specific
	TargetCriteria      string `gorm:"type:json" json:"target_criteria"` // JSON criteria for segment
	TargetList          string `gorm:"type:json" json:"target_list"`     // JSON array of specific IDs
	EstimatedRecipients int    `json:"estimated_recipients"`

	// Content
	TemplateID uuid.UUID `gorm:"type:uuid" json:"template_id"`
	Subject    string    `json:"subject"`
	Content    string    `gorm:"type:text" json:"content"`

	// Scheduling
	ScheduledFor  time.Time `json:"scheduled_for"`
	TimeZoneAware bool      `gorm:"default:false" json:"timezone_aware"`
	SendInBatches bool      `gorm:"default:false" json:"send_in_batches"`
	BatchSize     int       `json:"batch_size"`
	BatchInterval int       `json:"batch_interval"` // minutes between batches

	// Status
	Status      string     `gorm:"default:'draft'" json:"status"` // draft, scheduled, running, completed, cancelled
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`

	// Results
	TotalSent      int `gorm:"default:0" json:"total_sent"`
	TotalDelivered int `gorm:"default:0" json:"total_delivered"`
	TotalOpened    int `gorm:"default:0" json:"total_opened"`
	TotalClicked   int `gorm:"default:0" json:"total_clicked"`
	TotalFailed    int `gorm:"default:0" json:"total_failed"`

	// Performance
	DeliveryRate float64 `gorm:"default:0" json:"delivery_rate"`
	OpenRate     float64 `gorm:"default:0" json:"open_rate"`
	ClickRate    float64 `gorm:"default:0" json:"click_rate"`

	CreatedBy  uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	ApprovedBy *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Template       CommunicationTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Communications []Communication       `gorm:"foreignKey:TriggerID" json:"communications,omitempty"`
}

// CommunicationQueue manages the queue for pending communications
type CommunicationQueue struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CommunicationID uuid.UUID  `gorm:"type:uuid;not null" json:"communication_id"`
	Priority        int        `gorm:"default:5" json:"priority"` // 1-10, 1 being highest
	ScheduledFor    time.Time  `json:"scheduled_for"`
	RetryCount      int        `gorm:"default:0" json:"retry_count"`
	MaxRetries      int        `gorm:"default:3" json:"max_retries"`
	Status          string     `gorm:"default:'pending'" json:"status"` // pending, processing, completed, failed
	ProcessedAt     *time.Time `json:"processed_at"`
	NextRetryAt     *time.Time `json:"next_retry_at"`
	FailureReason   string     `json:"failure_reason"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Communication Communication `gorm:"foreignKey:CommunicationID" json:"communication,omitempty"`
}

// Table names
func (Communication) TableName() string          { return "communications" }
func (CommunicationTemplate) TableName() string  { return "communication_templates" }
func (CommunicationEvent) TableName() string     { return "communication_events" }
func (NotificationPreference) TableName() string { return "notification_preferences" }
func (CommunicationCampaign) TableName() string  { return "communication_campaigns" }
func (CommunicationQueue) TableName() string     { return "communication_queue" }

// BeforeCreate hooks
func (c *Communication) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.CommunicationCode == "" {
		c.CommunicationCode = "COMM-" + time.Now().Format("20060102150405") + "-" + uuid.New().String()[:6]
	}
	if c.UnsubscribeToken == "" {
		c.UnsubscribeToken = uuid.New().String()
	}
	return nil
}

func (ct *CommunicationTemplate) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == uuid.Nil {
		ct.ID = uuid.New()
	}
	if ct.TemplateCode == "" {
		ct.TemplateCode = "TMPL-" + uuid.New().String()[:8]
	}
	return nil
}

func (ce *CommunicationEvent) BeforeCreate(tx *gorm.DB) error {
	if ce.ID == uuid.Nil {
		ce.ID = uuid.New()
	}
	return nil
}

func (np *NotificationPreference) BeforeCreate(tx *gorm.DB) error {
	if np.ID == uuid.Nil {
		np.ID = uuid.New()
	}
	return nil
}

func (cc *CommunicationCampaign) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	if cc.CampaignCode == "" {
		cc.CampaignCode = "CAMP-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	return nil
}

func (cq *CommunicationQueue) BeforeCreate(tx *gorm.DB) error {
	if cq.ID == uuid.Nil {
		cq.ID = uuid.New()
	}
	return nil
}

// Business Logic Methods

// Send marks the communication as sent
func (c *Communication) Send() {
	c.Status = "sent"
	now := time.Now()
	c.SentAt = &now
}

// MarkDelivered marks the communication as delivered
func (c *Communication) MarkDelivered() {
	c.Status = "delivered"
	now := time.Now()
	c.DeliveredAt = &now
}

// MarkOpened marks the communication as opened
func (c *Communication) MarkOpened() {
	now := time.Now()
	c.OpenedAt = &now
}

// MarkClicked marks the communication as clicked
func (c *Communication) MarkClicked() {
	now := time.Now()
	c.ClickedAt = &now
}

// MarkFailed marks the communication as failed
func (c *Communication) MarkFailed(reason string) {
	c.Status = "failed"
	c.FailureReason = reason
	c.DeliveryAttempts++
	now := time.Now()
	c.LastAttemptAt = &now
}

// ShouldRetry checks if communication should be retried
func (c *Communication) ShouldRetry() bool {
	return c.Status == "failed" && c.DeliveryAttempts < 3
}

// IsExpired checks if communication has expired
func (c *Communication) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// RequiresResponse checks if response is required and pending
func (c *Communication) RequiresResponse() bool {
	if !c.ResponseRequired || c.ResponseReceived {
		return false
	}
	if c.ResponseDeadline != nil && time.Now().After(*c.ResponseDeadline) {
		return false
	}
	return true
}

// UseTemplate increments the template use count
func (ct *CommunicationTemplate) UseTemplate() {
	ct.UseCount++
	now := time.Now()
	ct.LastUsedAt = &now
}

// UpdateSuccessRate updates the template success rate
func (ct *CommunicationTemplate) UpdateSuccessRate(wasSuccessful bool) {
	if ct.UseCount == 0 {
		if wasSuccessful {
			ct.SuccessRate = 100.0
		} else {
			ct.SuccessRate = 0.0
		}
	} else {
		// Calculate weighted average
		totalSuccess := ct.SuccessRate * float64(ct.UseCount-1) / 100.0
		if wasSuccessful {
			totalSuccess += 1.0
		}
		ct.SuccessRate = (totalSuccess / float64(ct.UseCount)) * 100.0
	}
}

// CanSendNow checks if communication can be sent based on preferences
func (np *NotificationPreference) CanSendNow(channelType string, category string) bool {
	// Check channel enabled
	switch channelType {
	case "email":
		if !np.EmailEnabled {
			return false
		}
	case "sms":
		if !np.SMSEnabled {
			return false
		}
	case "push":
		if !np.PushEnabled {
			return false
		}
	case "in_app":
		if !np.InAppEnabled {
			return false
		}
	case "whatsapp":
		if !np.WhatsAppEnabled {
			return false
		}
	}

	// Check category preferences
	switch category {
	case "transactional":
		return np.TransactionalAlerts
	case "promotional":
		return np.PromotionalMessages
	case "policy":
		return np.PolicyUpdates
	case "claim":
		return np.ClaimUpdates
	case "payment":
		return np.PaymentReminders
	case "security":
		return np.SecurityAlerts
	}

	return true
}

// IsInQuietHours checks if current time is in quiet hours
func (np *NotificationPreference) IsInQuietHours() bool {
	if np.QuietHoursStart == "" || np.QuietHoursEnd == "" {
		return false
	}

	// This would need proper time parsing and timezone handling
	// Simplified for demonstration
	return false
}

// StartCampaign starts the campaign
func (cc *CommunicationCampaign) StartCampaign() {
	cc.Status = "running"
	now := time.Now()
	cc.StartedAt = &now
}

// CompleteCampaign marks the campaign as completed
func (cc *CommunicationCampaign) CompleteCampaign() {
	cc.Status = "completed"
	now := time.Now()
	cc.CompletedAt = &now

	// Calculate rates
	if cc.TotalSent > 0 {
		cc.DeliveryRate = float64(cc.TotalDelivered) / float64(cc.TotalSent) * 100
		cc.OpenRate = float64(cc.TotalOpened) / float64(cc.TotalDelivered) * 100
		cc.ClickRate = float64(cc.TotalClicked) / float64(cc.TotalOpened) * 100
	}
}

// ShouldProcess checks if queue item should be processed
func (cq *CommunicationQueue) ShouldProcess() bool {
	if cq.Status != "pending" {
		return false
	}
	return time.Now().After(cq.ScheduledFor) || time.Now().Equal(cq.ScheduledFor)
}

// IncrementRetry increments retry count and sets next retry time
func (cq *CommunicationQueue) IncrementRetry() {
	cq.RetryCount++
	// Exponential backoff: 5min, 15min, 45min
	backoffMinutes := 5 * (1 << uint(cq.RetryCount))
	nextRetry := time.Now().Add(time.Duration(backoffMinutes) * time.Minute)
	cq.NextRetryAt = &nextRetry

	if cq.RetryCount >= cq.MaxRetries {
		cq.Status = "failed"
	}
}
