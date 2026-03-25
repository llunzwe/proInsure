package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserCommunicationHistory tracks all communications with users
type UserCommunicationHistory struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Communication Metrics
	TotalCommunications   int `gorm:"default:0" json:"total_communications"`
	EmailsSent            int `gorm:"default:0" json:"emails_sent"`
	SMSSent               int `gorm:"default:0" json:"sms_sent"`
	PushNotificationsSent int `gorm:"default:0" json:"push_notifications_sent"`
	WhatsAppSent          int `gorm:"default:0" json:"whatsapp_sent"`
	InAppMessages         int `gorm:"default:0" json:"in_app_messages"`
	PhoneCalls            int `gorm:"default:0" json:"phone_calls"`
	VideoCallsSent        int `gorm:"default:0" json:"video_calls"`
	PostalMailSent        int `gorm:"default:0" json:"postal_mail_sent"`

	// Engagement Metrics
	EmailOpenRate       float64 `gorm:"default:0" json:"email_open_rate"`
	EmailClickRate      float64 `gorm:"default:0" json:"email_click_rate"`
	SMSReadRate         float64 `gorm:"default:0" json:"sms_read_rate"`
	PushEngagementRate  float64 `gorm:"default:0" json:"push_engagement_rate"`
	ResponseRate        float64 `gorm:"default:0" json:"response_rate"`
	AverageResponseTime int     `json:"average_response_time_minutes"`
	UnsubscribeRate     float64 `gorm:"default:0" json:"unsubscribe_rate"`
	BounceRate          float64 `gorm:"default:0" json:"bounce_rate"`
	ConversionRate      float64 `gorm:"default:0" json:"conversion_rate"`

	// Channel Preferences
	PreferredChannels      []string           `gorm:"type:json" json:"preferred_channels"`
	ChannelEffectiveness   map[string]float64 `gorm:"type:json" json:"channel_effectiveness"`
	OptOutChannels         []string           `gorm:"type:json" json:"opt_out_channels"`
	ChannelFrequencyLimits map[string]int     `gorm:"type:json" json:"channel_frequency_limits"`
	BestTimeToContact      map[string]string  `gorm:"type:json" json:"best_time_to_contact"`
	LastChannelUsed        string             `gorm:"type:varchar(50)" json:"last_channel_used"`
	LastContactDate        *time.Time         `json:"last_contact_date"`
	NextScheduledContact   *time.Time         `json:"next_scheduled_contact"`
}

// UserCommunicationLog logs individual communication events
type UserCommunicationLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	CampaignID *uuid.UUID `gorm:"type:uuid" json:"campaign_id"`
	TemplateID *uuid.UUID `gorm:"type:uuid" json:"template_id"`
	SenderID   *uuid.UUID `gorm:"type:uuid" json:"sender_id"`

	// Communication Details
	MessageID        string                 `gorm:"type:varchar(100);uniqueIndex" json:"message_id"`
	Channel          string                 `gorm:"type:varchar(50)" json:"channel"`
	Type             string                 `gorm:"type:varchar(50)" json:"type"` // transactional/promotional/notification/alert
	Subject          string                 `json:"subject"`
	Content          string                 `gorm:"type:text" json:"content"`
	Status           string                 `gorm:"type:varchar(20)" json:"status"`    // sent/delivered/read/failed
	Direction        string                 `gorm:"type:varchar(10)" json:"direction"` // inbound/outbound
	Priority         string                 `gorm:"type:varchar(20)" json:"priority"`
	SentAt           time.Time              `json:"sent_at"`
	DeliveredAt      *time.Time             `json:"delivered_at"`
	ReadAt           *time.Time             `json:"read_at"`
	RespondedAt      *time.Time             `json:"responded_at"`
	BouncedAt        *time.Time             `json:"bounced_at"`
	UnsubscribedAt   *time.Time             `json:"unsubscribed_at"`
	FailureReason    string                 `gorm:"type:text" json:"failure_reason"`
	RetryCount       int                    `gorm:"default:0" json:"retry_count"`
	Metadata         map[string]interface{} `gorm:"type:json" json:"metadata"`
	TrackingData     map[string]interface{} `gorm:"type:json" json:"tracking_data"`
	ResponseData     map[string]interface{} `gorm:"type:json" json:"response_data"`
	Cost             decimal.Decimal        `gorm:"type:decimal(10,4)" json:"cost"`
	BillingReference string                 `gorm:"type:varchar(100)" json:"billing_reference"`
}

// UserNotificationPreferences manages user notification preferences
type UserNotificationPreferences struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Global Settings
	NotificationsEnabled bool              `gorm:"default:true" json:"notifications_enabled"`
	DoNotDisturb         bool              `gorm:"default:false" json:"do_not_disturb"`
	DoNotDisturbStart    string            `gorm:"type:varchar(10)" json:"do_not_disturb_start"` // HH:MM format
	DoNotDisturbEnd      string            `gorm:"type:varchar(10)" json:"do_not_disturb_end"`
	QuietHours           map[string]string `gorm:"type:json" json:"quiet_hours"`
	Timezone             string            `gorm:"type:varchar(50)" json:"timezone"`
	Language             string            `gorm:"type:varchar(10)" json:"language"`

	// Channel Preferences
	EmailEnabled      bool   `gorm:"default:true" json:"email_enabled"`
	EmailFrequency    string `gorm:"type:varchar(20)" json:"email_frequency"` // immediate/daily/weekly/monthly
	EmailDigest       bool   `gorm:"default:false" json:"email_digest"`
	SMSEnabled        bool   `gorm:"default:false" json:"sms_enabled"`
	SMSFrequency      string `gorm:"type:varchar(20)" json:"sms_frequency"`
	PushEnabled       bool   `gorm:"default:true" json:"push_enabled"`
	PushSound         bool   `gorm:"default:true" json:"push_sound"`
	PushVibration     bool   `gorm:"default:true" json:"push_vibration"`
	WhatsAppEnabled   bool   `gorm:"default:false" json:"whatsapp_enabled"`
	InAppEnabled      bool   `gorm:"default:true" json:"in_app_enabled"`
	PhoneCallsEnabled bool   `gorm:"default:false" json:"phone_calls_enabled"`
	PostalMailEnabled bool   `gorm:"default:false" json:"postal_mail_enabled"`

	// Category Preferences
	MarketingOptIn      bool `gorm:"default:false" json:"marketing_opt_in"`
	TransactionalOptIn  bool `gorm:"default:true" json:"transactional_opt_in"`
	AlertsOptIn         bool `gorm:"default:true" json:"alerts_opt_in"`
	NewsletterOptIn     bool `gorm:"default:false" json:"newsletter_opt_in"`
	ProductUpdatesOptIn bool `gorm:"default:true" json:"product_updates_opt_in"`
	SecurityAlertsOptIn bool `gorm:"default:true" json:"security_alerts_opt_in"`
	BillingAlertsOptIn  bool `gorm:"default:true" json:"billing_alerts_opt_in"`
	PolicyAlertsOptIn   bool `gorm:"default:true" json:"policy_alerts_opt_in"`
	ClaimAlertsOptIn    bool `gorm:"default:true" json:"claim_alerts_opt_in"`

	// Notification Rules
	NotificationRules map[string]interface{} `gorm:"type:json" json:"notification_rules"`
	CustomPreferences map[string]interface{} `gorm:"type:json" json:"custom_preferences"`
	FrequencyCaps     map[string]int         `gorm:"type:json" json:"frequency_caps"`
	LastUpdated       time.Time              `json:"last_updated"`
	ConsentDate       *time.Time             `json:"consent_date"`
	ConsentVersion    string                 `gorm:"type:varchar(20)" json:"consent_version"`
}

// UserCampaignEngagement tracks user engagement with marketing campaigns
type UserCampaignEngagement struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CampaignID uuid.UUID `gorm:"type:uuid;not null;index" json:"campaign_id"`

	// Engagement Details
	CampaignName               string                   `gorm:"type:varchar(100)" json:"campaign_name"`
	CampaignType               string                   `gorm:"type:varchar(50)" json:"campaign_type"`
	FirstExposure              time.Time                `json:"first_exposure"`
	LastExposure               *time.Time               `json:"last_exposure"`
	ExposureCount              int                      `gorm:"default:0" json:"exposure_count"`
	Impressions                int                      `gorm:"default:0" json:"impressions"`
	Clicks                     int                      `gorm:"default:0" json:"clicks"`
	Conversions                int                      `gorm:"default:0" json:"conversions"`
	Revenue                    decimal.Decimal          `gorm:"type:decimal(15,2)" json:"revenue"`
	EngagementScore            float64                  `gorm:"default:0" json:"engagement_score"`
	ResponseTime               int                      `json:"response_time_hours"`
	DevicesUsed                []string                 `gorm:"type:json" json:"devices_used"`
	ChannelsEngaged            []string                 `gorm:"type:json" json:"channels_engaged"`
	Actions                    []map[string]interface{} `gorm:"type:json" json:"actions"`
	AttributedPurchases        []uuid.UUID              `gorm:"type:json" json:"attributed_purchases"`
	ViralReach                 int                      `gorm:"default:0" json:"viral_reach"`
	SharedCount                int                      `gorm:"default:0" json:"shared_count"`
	FeedbackProvided           bool                     `gorm:"default:false" json:"feedback_provided"`
	FeedbackScore              *int                     `json:"feedback_score"`
	UnsubscribedDuringCampaign bool                     `gorm:"default:false" json:"unsubscribed_during_campaign"`
}

// UserConversation tracks conversational interactions
type UserConversation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	AgentID  *uuid.UUID `gorm:"type:uuid" json:"agent_id"`
	BotID    *uuid.UUID `gorm:"type:uuid" json:"bot_id"`
	TicketID *uuid.UUID `gorm:"type:uuid" json:"ticket_id"`

	// Conversation Details
	ConversationID      string                   `gorm:"type:varchar(100);uniqueIndex" json:"conversation_id"`
	Channel             string                   `gorm:"type:varchar(50)" json:"channel"`
	Type                string                   `gorm:"type:varchar(50)" json:"type"` // support/sales/inquiry/complaint
	Status              string                   `gorm:"type:varchar(20)" json:"status"`
	StartTime           time.Time                `json:"start_time"`
	EndTime             *time.Time               `json:"end_time"`
	Duration            int                      `json:"duration_seconds"`
	MessageCount        int                      `gorm:"default:0" json:"message_count"`
	Messages            []map[string]interface{} `gorm:"type:json" json:"messages"`
	Transcript          string                   `gorm:"type:text" json:"transcript"`
	Language            string                   `gorm:"type:varchar(10)" json:"language"`
	Sentiment           string                   `gorm:"type:varchar(20)" json:"sentiment"`
	SentimentScore      float64                  `gorm:"default:0" json:"sentiment_score"`
	Topics              []string                 `gorm:"type:json" json:"topics"`
	Intent              string                   `gorm:"type:varchar(100)" json:"intent"`
	Resolution          string                   `gorm:"type:text" json:"resolution"`
	SatisfactionScore   *int                     `json:"satisfaction_score"`
	TransferCount       int                      `gorm:"default:0" json:"transfer_count"`
	EscalationRequired  bool                     `gorm:"default:false" json:"escalation_required"`
	AIAssisted          bool                     `gorm:"default:false" json:"ai_assisted"`
	HandoffToHuman      bool                     `gorm:"default:false" json:"handoff_to_human"`
	ConversationQuality float64                  `gorm:"default:0" json:"conversation_quality"`
}

// UserMessageTemplate manages personalized message templates
type UserMessageTemplate struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Template Details
	TemplateName       string                 `gorm:"type:varchar(100)" json:"template_name"`
	TemplateType       string                 `gorm:"type:varchar(50)" json:"template_type"`
	Channel            string                 `gorm:"type:varchar(50)" json:"channel"`
	Language           string                 `gorm:"type:varchar(10)" json:"language"`
	Subject            string                 `json:"subject"`
	Content            string                 `gorm:"type:text" json:"content"`
	Variables          []string               `gorm:"type:json" json:"variables"`
	Personalization    map[string]interface{} `gorm:"type:json" json:"personalization"`
	Active             bool                   `gorm:"default:true" json:"active"`
	UsageCount         int                    `gorm:"default:0" json:"usage_count"`
	LastUsed           *time.Time             `json:"last_used"`
	PerformanceMetrics map[string]float64     `gorm:"type:json" json:"performance_metrics"`
	ABTestVariant      string                 `gorm:"type:varchar(10)" json:"ab_test_variant"`
	VersionNumber      int                    `gorm:"default:1" json:"version_number"`
	ApprovedBy         *uuid.UUID             `gorm:"type:uuid" json:"approved_by"`
	ApprovalDate       *time.Time             `json:"approval_date"`
}

// UserSentiment tracks user sentiment over time
type UserSentiment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Sentiment Analysis
	CurrentSentiment     string                   `gorm:"type:varchar(20)" json:"current_sentiment"` // positive/neutral/negative
	SentimentScore       float64                  `gorm:"default:0" json:"sentiment_score"`          // -1 to 1
	SentimentTrend       string                   `gorm:"type:varchar(20)" json:"sentiment_trend"`   // improving/stable/declining
	LastAnalysisDate     time.Time                `json:"last_analysis_date"`
	DataSources          []string                 `gorm:"type:json" json:"data_sources"`
	PositiveIndicators   []string                 `gorm:"type:json" json:"positive_indicators"`
	NegativeIndicators   []string                 `gorm:"type:json" json:"negative_indicators"`
	EmotionDistribution  map[string]float64       `gorm:"type:json" json:"emotion_distribution"`
	TopicSentiment       map[string]float64       `gorm:"type:json" json:"topic_sentiment"`
	ChannelSentiment     map[string]float64       `gorm:"type:json" json:"channel_sentiment"`
	SentimentHistory     []map[string]interface{} `gorm:"type:json" json:"sentiment_history"`
	TriggerEvents        []string                 `gorm:"type:json" json:"trigger_events"`
	RecommendedActions   []string                 `gorm:"type:json" json:"recommended_actions"`
	InterventionRequired bool                     `gorm:"default:false" json:"intervention_required"`
	LastIntervention     *time.Time               `json:"last_intervention"`
	InterventionOutcome  string                   `gorm:"type:text" json:"intervention_outcome"`
}

// TableName returns the table name
func (UserCommunicationHistory) TableName() string {
	return "user_communication_history"
}

// TableName returns the table name
func (UserCommunicationLog) TableName() string {
	return "user_communication_logs"
}

// TableName returns the table name
func (UserNotificationPreferences) TableName() string {
	return "user_notification_preferences"
}

// TableName returns the table name
func (UserCampaignEngagement) TableName() string {
	return "user_campaign_engagement"
}

// TableName returns the table name
func (UserConversation) TableName() string {
	return "user_conversations"
}

// TableName returns the table name
func (UserMessageTemplate) TableName() string {
	return "user_message_templates"
}

// TableName returns the table name
func (UserSentiment) TableName() string {
	return "user_sentiment"
}
