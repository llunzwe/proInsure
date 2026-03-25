package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserPaymentMethod represents a user's saved payment method
type UserPaymentMethod struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Type           string         `gorm:"type:varchar(20);not null" json:"type"` // card/bank/wallet/crypto
	Provider       string         `json:"provider"`                              // visa/mastercard/paypal/stripe
	Last4          string         `json:"last4"`
	ExpiryMonth    int            `json:"expiry_month,omitempty"`
	ExpiryYear     int            `json:"expiry_year,omitempty"`
	HolderName     string         `json:"holder_name"`
	BillingAddress string         `gorm:"type:json" json:"billing_address"`
	IsDefault      bool           `gorm:"default:false" json:"is_default"`
	IsVerified     bool           `gorm:"default:false" json:"is_verified"`
	TokenID        string         `json:"-"`           // Payment gateway token
	Fingerprint    string         `json:"fingerprint"` // For duplicate detection
	Status         string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	FailureCount   int            `gorm:"default:0" json:"failure_count"`
	LastUsedAt     *time.Time     `json:"last_used_at"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserInsuranceHistory represents a user's insurance history with other providers
type UserInsuranceHistory struct {
	ID                 uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	UserID             uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	InsurerName        string          `gorm:"not null" json:"insurer_name"`
	PolicyType         string          `json:"policy_type"`
	PolicyNumber       string          `json:"policy_number"`
	StartDate          time.Time       `json:"start_date"`
	EndDate            *time.Time      `json:"end_date"`
	PremiumAmount      decimal.Decimal `gorm:"type:decimal(10,2)" json:"premium_amount"`
	ClaimsCount        int             `gorm:"default:0" json:"claims_count"`
	ClaimsAmount       decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"claims_amount"`
	CancellationReason string          `json:"cancellation_reason"`
	Rating             int             `json:"rating"` // 1-5 stars
	WouldRecommend     bool            `json:"would_recommend"`
	Notes              string          `gorm:"type:text" json:"notes"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserSupportTicket represents a customer support ticket
type UserSupportTicket struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	TicketNumber       string         `gorm:"uniqueIndex;not null" json:"ticket_number"`
	Category           string         `gorm:"not null" json:"category"`         // billing/claim/technical/general
	Priority           string         `gorm:"default:'normal'" json:"priority"` // low/normal/high/urgent
	Subject            string         `gorm:"not null" json:"subject"`
	Description        string         `gorm:"type:text;not null" json:"description"`
	Status             string         `gorm:"default:'open'" json:"status"` // open/pending/resolved/closed
	AssignedTo         *uuid.UUID     `gorm:"type:uuid" json:"assigned_to"`
	Resolution         string         `gorm:"type:text" json:"resolution"`
	SatisfactionRating *int           `json:"satisfaction_rating"` // 1-5
	ResponseTime       int            `json:"response_time_hours"`
	ResolutionTime     int            `json:"resolution_time_hours"`
	EscalatedAt        *time.Time     `json:"escalated_at"`
	ResolvedAt         *time.Time     `json:"resolved_at"`
	ClosedAt           *time.Time     `json:"closed_at"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserNotificationLog tracks all notifications sent to user
type UserNotificationLog struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Type          string     `gorm:"not null" json:"type"` // email/sms/push/in-app/whatsapp
	Category      string     `json:"category"`             // marketing/transactional/alert/reminder
	Subject       string     `json:"subject"`
	Content       string     `gorm:"type:text" json:"content"`
	Channel       string     `json:"channel"`
	Status        string     `gorm:"default:'pending'" json:"status"` // pending/sent/delivered/failed/bounced
	DeliveredAt   *time.Time `json:"delivered_at"`
	ReadAt        *time.Time `json:"read_at"`
	ClickedAt     *time.Time `json:"clicked_at"`
	FailureReason string     `json:"failure_reason"`
	Metadata      string     `gorm:"type:json" json:"metadata"`
	CampaignID    string     `json:"campaign_id"`
	TemplateID    string     `json:"template_id"`
	Priority      int        `gorm:"default:0" json:"priority"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserAuditLog tracks all user actions for compliance
type UserAuditLog struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Action       string     `gorm:"not null" json:"action"`
	EntityType   string     `json:"entity_type"`
	EntityID     *uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	OldValue     string     `gorm:"type:json" json:"old_value"`
	NewValue     string     `gorm:"type:json" json:"new_value"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	SessionID    string     `json:"session_id"`
	Result       string     `json:"result"` // success/failure
	ErrorMessage string     `json:"error_message"`
	Duration     int        `json:"duration_ms"`
	Metadata     string     `gorm:"type:json" json:"metadata"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// UserSecurityEvent tracks security-related events
type UserSecurityEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	EventType   string    `gorm:"not null" json:"event_type"`     // login_attempt/password_change/2fa_enabled/suspicious_activity
	Severity    string    `gorm:"default:'info'" json:"severity"` // info/warning/critical
	Success     bool      `gorm:"default:false" json:"success"`
	IPAddress   string    `json:"ip_address"`
	Location    string    `json:"location"`
	DeviceInfo  string    `gorm:"type:json" json:"device_info"`
	RiskScore   float64   `json:"risk_score"`
	ActionTaken string    `json:"action_taken"` // blocked/allowed/challenged
	Details     string    `gorm:"type:text" json:"details"`
	Metadata    string    `gorm:"type:json" json:"metadata"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// UserPreference stores user preferences and settings
type UserPreference struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Category    string    `gorm:"not null" json:"category"` // notification/privacy/display/payment
	Key         string    `gorm:"not null" json:"key"`
	Value       string    `gorm:"type:text" json:"value"`
	Type        string    `json:"type"` // boolean/string/number/json
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Composite unique index on UserID + Category + Key
}

// UserDocument stores user-uploaded documents
type UserDocument struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Type            string         `gorm:"not null" json:"type"` // id_proof/address_proof/income_proof/medical
	Name            string         `gorm:"not null" json:"name"`
	FileName        string         `json:"file_name"`
	FileSize        int64          `json:"file_size"`
	MimeType        string         `json:"mime_type"`
	StoragePath     string         `json:"-"`
	StorageURL      string         `json:"storage_url"`
	Hash            string         `json:"hash"`                            // For duplicate detection
	Status          string         `gorm:"default:'pending'" json:"status"` // pending/verified/rejected
	VerifiedBy      *uuid.UUID     `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt      *time.Time     `json:"verified_at"`
	RejectionReason string         `json:"rejection_reason"`
	ExpiryDate      *time.Time     `json:"expiry_date"`
	Metadata        string         `gorm:"type:json" json:"metadata"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserActivity tracks user engagement and behavior
type UserActivity struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Type        string    `gorm:"not null" json:"type"` // page_view/feature_use/search/purchase
	Category    string    `json:"category"`
	Action      string    `json:"action"`
	Label       string    `json:"label"`
	Value       float64   `json:"value"`
	Duration    int       `json:"duration_seconds"`
	PageURL     string    `json:"page_url"`
	ReferrerURL string    `json:"referrer_url"`
	SearchQuery string    `json:"search_query"`
	ResultCount int       `json:"result_count"`
	Platform    string    `json:"platform"` // web/ios/android
	AppVersion  string    `json:"app_version"`
	SessionID   string    `json:"session_id"`
	Metadata    string    `gorm:"type:json" json:"metadata"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// UserReward represents loyalty rewards and points
type UserReward struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	Type        string          `gorm:"not null" json:"type"` // points/cashback/discount/voucher
	Category    string          `json:"category"`             // referral/purchase/loyalty/achievement
	Points      int             `gorm:"default:0" json:"points"`
	CashValue   decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"cash_value"`
	Description string          `json:"description"`
	Status      string          `gorm:"default:'active'" json:"status"` // active/redeemed/expired
	EarnedFrom  string          `json:"earned_from"`                    // Source of reward
	RedeemedAt  *time.Time      `json:"redeemed_at"`
	RedeemedFor string          `json:"redeemed_for"`
	ExpiresAt   *time.Time      `json:"expires_at"`
	Metadata    string          `gorm:"type:json" json:"metadata"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// === Table Names ===

func (UserPaymentMethod) TableName() string    { return "user_payment_methods" }
func (UserInsuranceHistory) TableName() string { return "user_insurance_histories" }
func (UserSupportTicket) TableName() string    { return "user_support_tickets" }
func (UserNotificationLog) TableName() string  { return "user_notification_logs" }
func (UserAuditLog) TableName() string         { return "user_audit_logs" }
func (UserSecurityEvent) TableName() string    { return "user_security_events" }
func (UserPreference) TableName() string       { return "user_preferences" }
func (UserDocument) TableName() string         { return "user_documents" }
func (UserActivity) TableName() string         { return "user_activities" }
func (UserReward) TableName() string           { return "user_rewards" }

// === Model Methods ===

// BeforeCreate hooks for UUID generation
func (m *UserPaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserInsuranceHistory) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserSupportTicket) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.TicketNumber == "" {
		m.TicketNumber = generateTicketNumber()
	}
	return nil
}

func (m *UserNotificationLog) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserAuditLog) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserSecurityEvent) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserPreference) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserDocument) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserActivity) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *UserReward) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// Helper function to generate ticket number
func generateTicketNumber() string {
	return "TKT-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:8]
}

// IsExpired checks if payment method is expired
func (pm *UserPaymentMethod) IsExpired() bool {
	if pm.Type != "card" {
		return false
	}
	now := time.Now()
	return pm.ExpiryYear < now.Year() || (pm.ExpiryYear == now.Year() && pm.ExpiryMonth < int(now.Month()))
}

// IsExpired checks if reward is expired
func (r *UserReward) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// CanRedeem checks if reward can be redeemed
func (r *UserReward) CanRedeem() bool {
	return r.Status == "active" && !r.IsExpired()
}
