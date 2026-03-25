package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// CustomerPortal represents customer self-service portal sessions
type CustomerPortal struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	SessionToken      string    `gorm:"uniqueIndex;not null" json:"session_token"`
	IPAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	LastActivity      time.Time `gorm:"not null" json:"last_activity"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	LoginMethod       string    `json:"login_method"` // password, biometric, sso
	DeviceFingerprint string    `json:"device_fingerprint"`
	GeolocationData   string    `json:"geolocation_data"`  // JSON object
	SessionDuration   int       `json:"session_duration"`  // seconds
	ActionsPerformed  string    `json:"actions_performed"` // JSON array
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ExpiresAt         time.Time `gorm:"not null" json:"expires_at"`

	// Relationships
	User       models.User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Activities []PortalActivity `gorm:"foreignKey:SessionID" json:"activities,omitempty"`
}

// PortalActivity represents user activities in the customer portal
type PortalActivity struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SessionID     uuid.UUID  `gorm:"type:uuid;not null" json:"session_id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	ActivityType  string     `gorm:"not null" json:"activity_type"` // view_policy, file_claim, update_profile, etc.
	ResourceID    *uuid.UUID `gorm:"type:uuid" json:"resource_id"`  // ID of the resource accessed
	ResourceType  string     `json:"resource_type"`                 // policy, claim, device, etc.
	ActionDetails string     `json:"action_details"`                // JSON object with specific action data
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	Duration      int        `json:"duration"` // seconds spent on activity
	IsSuccessful  bool       `gorm:"default:true" json:"is_successful"`
	ErrorMessage  string     `json:"error_message"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Session CustomerPortal `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	User    models.User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ChatbotConversation represents AI chatbot conversations
type ChatbotConversation struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID              *uuid.UUID `gorm:"type:uuid" json:"user_id"` // null for anonymous users
	SessionID           string     `gorm:"not null" json:"session_id"`
	ConversationID      string     `gorm:"uniqueIndex;not null" json:"conversation_id"`
	Channel             string     `gorm:"not null" json:"channel"` // web, mobile, whatsapp, telegram
	Language            string     `gorm:"default:'en'" json:"language"`
	StartedAt           time.Time  `gorm:"not null" json:"started_at"`
	EndedAt             *time.Time `json:"ended_at"`
	IsActive            bool       `gorm:"default:true" json:"is_active"`
	SatisfactionScore   int        `json:"satisfaction_score"` // 1-5
	WasResolved         bool       `gorm:"default:false" json:"was_resolved"`
	EscalatedToHuman    bool       `gorm:"default:false" json:"escalated_to_human"`
	EscalationReason    string     `json:"escalation_reason"`
	TotalMessages       int        `gorm:"default:0" json:"total_messages"`
	UserMessages        int        `gorm:"default:0" json:"user_messages"`
	BotMessages         int        `gorm:"default:0" json:"bot_messages"`
	ConversationSummary string     `json:"conversation_summary"`
	Intent              string     `json:"intent"`   // policy_inquiry, claim_status, general_support
	Entities            string     `json:"entities"` // JSON object with extracted entities
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User *models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// Messages relationship is defined in ChatbotMessage to avoid circular dependency
}

// ChatbotMessage represents individual messages in chatbot conversations
type ChatbotMessage struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ConversationID  uuid.UUID `gorm:"type:uuid;not null" json:"conversation_id"`
	MessageID       string    `gorm:"uniqueIndex;not null" json:"message_id"`
	MessageType     string    `gorm:"not null" json:"message_type"` // user, bot, system
	Content         string    `gorm:"not null" json:"content"`
	Intent          string    `json:"intent"`
	Confidence      float64   `json:"confidence"`
	Entities        string    `json:"entities"`      // JSON object
	ResponseTime    int       `json:"response_time"` // milliseconds
	IsProcessed     bool      `gorm:"default:true" json:"is_processed"`
	ProcessingError string    `json:"processing_error"`
	Attachments     string    `json:"attachments"`   // JSON array of attachment URLs
	QuickReplies    string    `json:"quick_replies"` // JSON array of suggested responses
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Conversation ChatbotConversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

// GamificationProfile represents user gamification and loyalty data
type GamificationProfile struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID            uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	LoyaltyTier       string     `gorm:"default:'bronze'" json:"loyalty_tier"` // bronze, silver, gold, platinum
	TotalPoints       int        `gorm:"default:0" json:"total_points"`
	AvailablePoints   int        `gorm:"default:0" json:"available_points"`
	LifetimePoints    int        `gorm:"default:0" json:"lifetime_points"`
	NoClaimBonus      int        `gorm:"default:0" json:"no_claim_bonus"` // months without claims
	ReferralCount     int        `gorm:"default:0" json:"referral_count"`
	SafeDriverScore   float64    `gorm:"default:0" json:"safe_driver_score"`
	EngagementScore   float64    `gorm:"default:0" json:"engagement_score"`
	LastActivityDate  *time.Time `json:"last_activity_date"`
	NextTierThreshold int        `json:"next_tier_threshold"`
	TierUpgradeDate   *time.Time `json:"tier_upgrade_date"`
	AnnualSpending    float64    `gorm:"default:0" json:"annual_spending"`
	DiscountEarned    float64    `gorm:"default:0" json:"discount_earned"`
	BadgesEarned      string     `json:"badges_earned"` // JSON array of badge IDs
	Achievements      string     `json:"achievements"`  // JSON array of achievement objects
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User              models.User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PointTransactions []PointTransaction `gorm:"foreignKey:UserID" json:"point_transactions,omitempty"`
	Rewards           []RewardRedemption `gorm:"foreignKey:UserID" json:"rewards,omitempty"`
}

// PointTransaction represents loyalty point transactions
type PointTransaction struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	TransactionType string     `gorm:"not null" json:"transaction_type"` // earned, redeemed, expired, adjusted
	Points          int        `gorm:"not null" json:"points"`           // positive for earned, negative for redeemed
	Reason          string     `gorm:"not null" json:"reason"`
	ReferenceID     *uuid.UUID `gorm:"type:uuid" json:"reference_id"` // policy, claim, referral ID
	ReferenceType   string     `json:"reference_type"`                // policy, claim, referral, manual
	Description     string     `json:"description"`
	ExpiresAt       *time.Time `json:"expires_at"`
	IsProcessed     bool       `gorm:"default:true" json:"is_processed"`
	ProcessedBy     *uuid.UUID `gorm:"type:uuid" json:"processed_by"` // admin user ID for manual adjustments
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User            models.User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProcessedByUser *models.User `gorm:"foreignKey:ProcessedBy" json:"processed_by_user,omitempty"`
}

// RewardRedemption represents reward redemptions
type RewardRedemption struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID            uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	RewardType        string     `gorm:"not null" json:"reward_type"` // discount, cashback, gift_card, premium_credit
	RewardValue       float64    `gorm:"not null" json:"reward_value"`
	PointsUsed        int        `gorm:"not null" json:"points_used"`
	Status            string     `gorm:"not null;default:'pending'" json:"status"` // pending, approved, redeemed, expired, cancelled
	RedemptionCode    string     `gorm:"uniqueIndex" json:"redemption_code"`
	ExpiresAt         time.Time  `gorm:"not null" json:"expires_at"`
	RedeemedAt        *time.Time `json:"redeemed_at"`
	AppliedToPolicyID *uuid.UUID `gorm:"type:uuid" json:"applied_to_policy_id"`
	Description       string     `json:"description"`
	Terms             string     `json:"terms"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User            models.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AppliedToPolicy *models.Policy `gorm:"foreignKey:AppliedToPolicyID" json:"applied_to_policy,omitempty"`
}

// NotificationPreference has been moved to communication.go for unified communication system
// See: internal/domain/models/communication.go

// NotificationLog represents sent notifications
type NotificationLog struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	NotificationType string     `gorm:"not null" json:"notification_type"` // email, sms, push, whatsapp
	Channel          string     `gorm:"not null" json:"channel"`
	Subject          string     `json:"subject"`
	Content          string     `gorm:"not null" json:"content"`
	Recipient        string     `gorm:"not null" json:"recipient"`                // email address, phone number, device token
	Status           string     `gorm:"not null;default:'pending'" json:"status"` // pending, sent, delivered, failed, bounced
	SentAt           *time.Time `json:"sent_at"`
	DeliveredAt      *time.Time `json:"delivered_at"`
	OpenedAt         *time.Time `json:"opened_at"`
	ClickedAt        *time.Time `json:"clicked_at"`
	ErrorMessage     string     `json:"error_message"`
	RetryCount       int        `gorm:"default:0" json:"retry_count"`
	MaxRetries       int        `gorm:"default:3" json:"max_retries"`
	Priority         string     `gorm:"default:'normal'" json:"priority"` // low, normal, high, urgent
	ReferenceID      *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	ReferenceType    string     `json:"reference_type"` // policy, claim, payment, etc.
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName methods
func (CustomerPortal) TableName() string {
	return "customer_portals"
}

func (PortalActivity) TableName() string {
	return "portal_activities"
}

func (ChatbotConversation) TableName() string {
	return "chatbot_conversations"
}

func (ChatbotMessage) TableName() string {
	return "chatbot_messages"
}

func (GamificationProfile) TableName() string {
	return "gamification_profiles"
}

func (PointTransaction) TableName() string {
	return "point_transactions"
}

func (RewardRedemption) TableName() string {
	return "reward_redemptions"
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}

// BeforeCreate hooks
func (cp *CustomerPortal) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == uuid.Nil {
		cp.ID = uuid.New()
	}
	return nil
}

func (pa *PortalActivity) BeforeCreate(tx *gorm.DB) error {
	if pa.ID == uuid.Nil {
		pa.ID = uuid.New()
	}
	return nil
}

func (cc *ChatbotConversation) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	if cc.ConversationID == "" {
		cc.ConversationID = "CHAT-" + uuid.New().String()[:8]
	}
	return nil
}

func (cm *ChatbotMessage) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == uuid.Nil {
		cm.ID = uuid.New()
	}
	if cm.MessageID == "" {
		cm.MessageID = "MSG-" + uuid.New().String()[:8]
	}
	return nil
}

func (gp *GamificationProfile) BeforeCreate(tx *gorm.DB) error {
	if gp.ID == uuid.Nil {
		gp.ID = uuid.New()
	}
	return nil
}

func (pt *PointTransaction) BeforeCreate(tx *gorm.DB) error {
	if pt.ID == uuid.Nil {
		pt.ID = uuid.New()
	}
	return nil
}

func (rr *RewardRedemption) BeforeCreate(tx *gorm.DB) error {
	if rr.ID == uuid.Nil {
		rr.ID = uuid.New()
	}
	if rr.RedemptionCode == "" {
		rr.RedemptionCode = "REWARD-" + uuid.New().String()[:8]
	}
	return nil
}

func (nl *NotificationLog) BeforeCreate(tx *gorm.DB) error {
	if nl.ID == uuid.Nil {
		nl.ID = uuid.New()
	}
	return nil
}

// Business logic methods for CustomerPortal
func (cp *CustomerPortal) IsExpired() bool {
	return time.Now().After(cp.ExpiresAt)
}

func (cp *CustomerPortal) ExtendSession(duration time.Duration) {
	cp.ExpiresAt = time.Now().Add(duration)
	cp.LastActivity = time.Now()
}

func (cp *CustomerPortal) Terminate() {
	cp.IsActive = false
}

// Business logic methods for GamificationProfile
func (gp *GamificationProfile) AddPoints(points int, reason string) {
	gp.TotalPoints += points
	gp.AvailablePoints += points
	gp.LifetimePoints += points
	gp.UpdateTier()
}

func (gp *GamificationProfile) RedeemPoints(points int) bool {
	if gp.AvailablePoints >= points {
		gp.AvailablePoints -= points
		return true
	}
	return false
}

func (gp *GamificationProfile) UpdateTier() {
	if gp.LifetimePoints >= 10000 {
		gp.LoyaltyTier = "platinum"
		gp.NextTierThreshold = 0
	} else if gp.LifetimePoints >= 5000 {
		gp.LoyaltyTier = "gold"
		gp.NextTierThreshold = 10000 - gp.LifetimePoints
	} else if gp.LifetimePoints >= 1000 {
		gp.LoyaltyTier = "silver"
		gp.NextTierThreshold = 5000 - gp.LifetimePoints
	} else {
		gp.LoyaltyTier = "bronze"
		gp.NextTierThreshold = 1000 - gp.LifetimePoints
	}
}

func (gp *GamificationProfile) CalculateDiscountPercentage() float64 {
	switch gp.LoyaltyTier {
	case "platinum":
		return 0.15 // 15% discount
	case "gold":
		return 0.10 // 10% discount
	case "silver":
		return 0.05 // 5% discount
	default:
		return 0.0 // No discount for bronze
	}
}

// Business logic methods for ChatbotConversation
func (cc *ChatbotConversation) EndConversation(resolved bool, satisfactionScore int) {
	now := time.Now()
	cc.EndedAt = &now
	cc.IsActive = false
	cc.WasResolved = resolved
	cc.SatisfactionScore = satisfactionScore
}

func (cc *ChatbotConversation) EscalateToHuman(reason string) {
	cc.EscalatedToHuman = true
	cc.EscalationReason = reason
}

func (cc *ChatbotConversation) AddMessage(messageType, content string) {
	cc.TotalMessages++
	if messageType == "user" {
		cc.UserMessages++
	} else if messageType == "bot" {
		cc.BotMessages++
	}
}

// Business logic methods for NotificationLog
func (nl *NotificationLog) MarkAsSent() {
	nl.Status = "sent"
	now := time.Now()
	nl.SentAt = &now
}

func (nl *NotificationLog) MarkAsDelivered() {
	nl.Status = "delivered"
	now := time.Now()
	nl.DeliveredAt = &now
}

func (nl *NotificationLog) MarkAsFailed(errorMessage string) {
	nl.Status = "failed"
	nl.ErrorMessage = errorMessage
	nl.RetryCount++
}

func (nl *NotificationLog) CanRetry() bool {
	return nl.RetryCount < nl.MaxRetries && nl.Status == "failed"
}
