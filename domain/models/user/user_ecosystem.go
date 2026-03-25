package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserDeviceEcosystem tracks all devices owned/used by the user
type UserDeviceEcosystem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Device Portfolio
	TotalDevicesOwned  int             `gorm:"default:0" json:"total_devices_owned"`
	ActiveDevices      int             `gorm:"default:0" json:"active_devices"`
	RetiredDevices     int             `gorm:"default:0" json:"retired_devices"`
	DeviceCategories   []string        `gorm:"type:json" json:"device_categories"`
	PreferredBrands    []string        `gorm:"type:json" json:"preferred_brands"`
	AverageDeviceAge   float64         `gorm:"default:0" json:"average_device_age_years"`
	AverageDeviceValue decimal.Decimal `gorm:"type:decimal(10,2)" json:"average_device_value"`
	TotalDeviceValue   decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_device_value"`

	// Usage Patterns
	PrimaryDeviceID        *uuid.UUID         `gorm:"type:uuid" json:"primary_device_id"`
	DeviceUsageHours       map[string]float64 `gorm:"type:json" json:"device_usage_hours"`
	DeviceSwitchingPattern string             `gorm:"type:varchar(50)" json:"device_switching_pattern"`
	MultiDeviceUser        bool               `gorm:"default:false" json:"multi_device_user"`
	DeviceSharingEnabled   bool               `gorm:"default:false" json:"device_sharing_enabled"`
	SharedWithUsers        []uuid.UUID        `gorm:"type:json" json:"shared_with_users"`

	// Upgrade Behavior
	UpgradeFrequency       float64                  `gorm:"default:0" json:"upgrade_frequency_years"`
	LastUpgradeDate        *time.Time               `json:"last_upgrade_date"`
	NextUpgradeEligibility *time.Time               `json:"next_upgrade_eligibility"`
	UpgradePreferences     map[string]interface{}   `gorm:"type:json" json:"upgrade_preferences"`
	TradeInHistory         []map[string]interface{} `gorm:"type:json" json:"trade_in_history"`
	EarlyAdopter           bool                     `gorm:"default:false" json:"early_adopter"`

	// Ecosystem Integration
	ConnectedDevices     []map[string]interface{} `gorm:"type:json" json:"connected_devices"`
	SmartHomeIntegration bool                     `gorm:"default:false" json:"smart_home_integration"`
	WearableDevices      []string                 `gorm:"type:json" json:"wearable_devices"`
	IoTDevices           []string                 `gorm:"type:json" json:"iot_devices"`
	CrossDeviceSync      bool                     `gorm:"default:false" json:"cross_device_sync"`
	CloudBackupEnabled   bool                     `gorm:"default:false" json:"cloud_backup_enabled"`

	// Device Care
	DeviceCareScore              float64  `gorm:"default:0" json:"device_care_score"`
	MaintenanceScheduleAdherence float64  `gorm:"default:0" json:"maintenance_schedule_adherence"`
	ProtectionProductsUsed       []string `gorm:"type:json" json:"protection_products_used"`
	RepairHistory                int      `gorm:"default:0" json:"repair_history_count"`
	AverageDeviceLifespan        float64  `gorm:"default:0" json:"average_device_lifespan_years"`
}

// UserPolicyPortfolio manages all insurance policies for the user
type UserPolicyPortfolio struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Portfolio Overview
	TotalPolicies   int             `gorm:"default:0" json:"total_policies"`
	ActivePolicies  int             `gorm:"default:0" json:"active_policies"`
	ExpiredPolicies int             `gorm:"default:0" json:"expired_policies"`
	PolicyTypes     []string        `gorm:"type:json" json:"policy_types"`
	TotalCoverage   decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_coverage"`
	TotalPremium    decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_premium"`

	// Bundle Management
	BundledPolicies      bool                       `gorm:"default:false" json:"bundled_policies"`
	BundleID             *uuid.UUID                 `gorm:"type:uuid" json:"bundle_id"`
	BundleDiscount       decimal.Decimal            `gorm:"type:decimal(5,2)" json:"bundle_discount"`
	BundleComponents     []uuid.UUID                `gorm:"type:json" json:"bundle_components"`
	CrossPolicyDiscounts map[string]decimal.Decimal `gorm:"type:json" json:"cross_policy_discounts"`

	// Coverage Analysis
	CoverageGaps         []string                 `gorm:"type:json" json:"coverage_gaps"`
	OverlapingCoverages  []string                 `gorm:"type:json" json:"overlaping_coverages"`
	CoverageEfficiency   float64                  `gorm:"default:0" json:"coverage_efficiency"`
	UnderinsuredAreas    []string                 `gorm:"type:json" json:"underinsured_areas"`
	OverinsuredAreas     []string                 `gorm:"type:json" json:"overinsured_areas"`
	RecommendedCoverages []map[string]interface{} `gorm:"type:json" json:"recommended_coverages"`

	// Portfolio Performance
	ClaimFrequency       float64 `gorm:"default:0" json:"claim_frequency"`
	ClaimRatio           float64 `gorm:"default:0" json:"claim_ratio"`
	ProfitabilityScore   float64 `gorm:"default:0" json:"profitability_score"`
	RetentionScore       float64 `gorm:"default:0" json:"retention_score"`
	PortfolioHealthScore float64 `gorm:"default:0" json:"portfolio_health_score"`

	// Renewal Management
	UpcomingRenewals   []map[string]interface{}   `gorm:"type:json" json:"upcoming_renewals"`
	AutoRenewalEnabled bool                       `gorm:"default:false" json:"auto_renewal_enabled"`
	RenewalReminders   map[string]time.Time       `gorm:"type:json" json:"renewal_reminders"`
	RenewalDiscounts   map[string]decimal.Decimal `gorm:"type:json" json:"renewal_discounts"`
	LapsedPolicies     []uuid.UUID                `gorm:"type:json" json:"lapsed_policies"`
}

// UserClaimPatterns analyzes claim behavior and patterns
type UserClaimPatterns struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Claim Statistics
	TotalClaims         int             `gorm:"default:0" json:"total_claims"`
	ApprovedClaims      int             `gorm:"default:0" json:"approved_claims"`
	DeniedClaims        int             `gorm:"default:0" json:"denied_claims"`
	PendingClaims       int             `gorm:"default:0" json:"pending_claims"`
	WithdrawnClaims     int             `gorm:"default:0" json:"withdrawn_claims"`
	TotalClaimAmount    decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_claim_amount"`
	TotalApprovedAmount decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_approved_amount"`

	// Claim Patterns
	ClaimFrequency     float64            `gorm:"default:0" json:"claim_frequency"`
	AverageClaimAmount decimal.Decimal    `gorm:"type:decimal(10,2)" json:"average_claim_amount"`
	ClaimTypes         map[string]int     `gorm:"type:json" json:"claim_types"`
	SeasonalPatterns   map[string]float64 `gorm:"type:json" json:"seasonal_patterns"`
	TimeOfDayPatterns  map[string]int     `gorm:"type:json" json:"time_of_day_patterns"`
	LocationPatterns   map[string]int     `gorm:"type:json" json:"location_patterns"`

	// Fraud Indicators
	SuspiciousPatterns []string `gorm:"type:json" json:"suspicious_patterns"`
	FraudScore         float64  `gorm:"default:0" json:"fraud_score"`
	InvestigationCount int      `gorm:"default:0" json:"investigation_count"`
	FraudulentClaims   int      `gorm:"default:0" json:"fraudulent_claims"`
	RedFlags           []string `gorm:"type:json" json:"red_flags"`

	// Behavior Analysis
	ClaimVelocity        float64  `gorm:"default:0" json:"claim_velocity"`
	RepeatClaimAreas     []string `gorm:"type:json" json:"repeat_claim_areas"`
	ClaimComplexity      float64  `gorm:"default:0" json:"claim_complexity"`
	DocumentationQuality float64  `gorm:"default:0" json:"documentation_quality"`
	ResponseTimeliness   float64  `gorm:"default:0" json:"response_timeliness"`

	// Predictive Insights
	NextClaimProbability   float64    `gorm:"default:0" json:"next_claim_probability"`
	EstimatedNextClaimDate *time.Time `json:"estimated_next_claim_date"`
	RiskCategory           string     `gorm:"type:varchar(20)" json:"risk_category"`
	PreventiveMeasures     []string   `gorm:"type:json" json:"preventive_measures"`
}

// UserPaymentEcosystem manages payment methods and schedules
type UserPaymentEcosystem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Payment Methods
	TotalPaymentMethods  int         `gorm:"default:0" json:"total_payment_methods"`
	ActivePaymentMethods int         `gorm:"default:0" json:"active_payment_methods"`
	PreferredMethod      string      `gorm:"type:varchar(50)" json:"preferred_method"`
	PaymentMethodTypes   []string    `gorm:"type:json" json:"payment_method_types"`
	DefaultPaymentMethod *uuid.UUID  `gorm:"type:uuid" json:"default_payment_method"`
	BackupPaymentMethods []uuid.UUID `gorm:"type:json" json:"backup_payment_methods"`

	// Payment Schedules
	ScheduledPayments  []map[string]interface{} `gorm:"type:json" json:"scheduled_payments"`
	RecurringPayments  []map[string]interface{} `gorm:"type:json" json:"recurring_payments"`
	AutoPayEnrollments []uuid.UUID              `gorm:"type:json" json:"auto_pay_enrollments"`
	PaymentCalendar    map[string]interface{}   `gorm:"type:json" json:"payment_calendar"`

	// Payment History
	TotalPaymentsMade int             `gorm:"default:0" json:"total_payments_made"`
	TotalAmountPaid   decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_amount_paid"`
	OnTimePaymentRate float64         `gorm:"default:0" json:"on_time_payment_rate"`
	FailedPayments    int             `gorm:"default:0" json:"failed_payments"`
	RefundedAmount    decimal.Decimal `gorm:"type:decimal(15,2)" json:"refunded_amount"`

	// Digital Wallets
	DigitalWallets     []map[string]interface{} `gorm:"type:json" json:"digital_wallets"`
	WalletIntegrations []string                 `gorm:"type:json" json:"wallet_integrations"`
	TokenizedCards     int                      `gorm:"default:0" json:"tokenized_cards"`
	BiometricPayments  bool                     `gorm:"default:false" json:"biometric_payments"`

	// Payment Preferences
	PreferredCurrency     string         `gorm:"type:varchar(3)" json:"preferred_currency"`
	MultiCurrencyEnabled  bool           `gorm:"default:false" json:"multi_currency_enabled"`
	InstallmentPreference bool           `gorm:"default:false" json:"installment_preference"`
	PaymentNotifications  bool           `gorm:"default:true" json:"payment_notifications"`
	PaymentReminders      map[string]int `gorm:"type:json" json:"payment_reminders"`
}

// UserReferralNetwork tracks referral relationships and rewards
type UserReferralNetwork struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	ReferredBy *uuid.UUID `gorm:"type:uuid" json:"referred_by"`

	// Referral Metrics
	ReferralCode           string  `gorm:"type:varchar(20);uniqueIndex" json:"referral_code"`
	TotalReferrals         int     `gorm:"default:0" json:"total_referrals"`
	SuccessfulReferrals    int     `gorm:"default:0" json:"successful_referrals"`
	PendingReferrals       int     `gorm:"default:0" json:"pending_referrals"`
	ReferralConversionRate float64 `gorm:"default:0" json:"referral_conversion_rate"`

	// Network Structure
	DirectReferrals   []uuid.UUID            `gorm:"type:json" json:"direct_referrals"`
	IndirectReferrals []uuid.UUID            `gorm:"type:json" json:"indirect_referrals"`
	NetworkDepth      int                    `gorm:"default:0" json:"network_depth"`
	NetworkSize       int                    `gorm:"default:0" json:"network_size"`
	ReferralTree      map[string]interface{} `gorm:"type:json" json:"referral_tree"`

	// Rewards & Earnings
	TotalEarnings    decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_earnings"`
	PendingEarnings  decimal.Decimal `gorm:"type:decimal(10,2)" json:"pending_earnings"`
	PaidEarnings     decimal.Decimal `gorm:"type:decimal(15,2)" json:"paid_earnings"`
	RewardPoints     int             `gorm:"default:0" json:"reward_points"`
	RewardTier       string          `gorm:"type:varchar(20)" json:"reward_tier"`
	BonusEligibility bool            `gorm:"default:false" json:"bonus_eligibility"`

	// Performance
	ReferralQuality       float64         `gorm:"default:0" json:"referral_quality"`
	AverageReferralValue  decimal.Decimal `gorm:"type:decimal(10,2)" json:"average_referral_value"`
	ReferralLifetimeValue decimal.Decimal `gorm:"type:decimal(15,2)" json:"referral_lifetime_value"`
	TopPerformer          bool            `gorm:"default:false" json:"top_performer"`
	InfluencerStatus      bool            `gorm:"default:false" json:"influencer_status"`

	// Campaigns
	ActiveCampaigns     []string                 `gorm:"type:json" json:"active_campaigns"`
	CampaignPerformance map[string]interface{}   `gorm:"type:json" json:"campaign_performance"`
	SpecialOffers       []map[string]interface{} `gorm:"type:json" json:"special_offers"`
}

// UserLoyaltyJourney tracks loyalty program participation and progress
type UserLoyaltyJourney struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Loyalty Status
	LoyaltyTier         string     `gorm:"type:varchar(20)" json:"loyalty_tier"`
	LoyaltyPoints       int        `gorm:"default:0" json:"loyalty_points"`
	LifetimePoints      int        `gorm:"default:0" json:"lifetime_points"`
	PointsExpiring      int        `gorm:"default:0" json:"points_expiring"`
	PointsExpiryDate    *time.Time `json:"points_expiry_date"`
	TierStatus          string     `gorm:"type:varchar(20)" json:"tier_status"`
	NextTierPoints      int        `gorm:"default:0" json:"next_tier_points"`
	TierProgressPercent float64    `gorm:"default:0" json:"tier_progress_percent"`

	// Milestones & Achievements
	MilestonesCompleted []string                 `gorm:"type:json" json:"milestones_completed"`
	CurrentMilestones   []map[string]interface{} `gorm:"type:json" json:"current_milestones"`
	Achievements        []map[string]interface{} `gorm:"type:json" json:"achievements"`
	Badges              []string                 `gorm:"type:json" json:"badges"`
	SpecialRecognitions []string                 `gorm:"type:json" json:"special_recognitions"`

	// Benefits & Rewards
	ActiveBenefits      []map[string]interface{} `gorm:"type:json" json:"active_benefits"`
	RedeemedRewards     []map[string]interface{} `gorm:"type:json" json:"redeemed_rewards"`
	AvailableRewards    []map[string]interface{} `gorm:"type:json" json:"available_rewards"`
	ExclusiveOffers     []map[string]interface{} `gorm:"type:json" json:"exclusive_offers"`
	PersonalizedRewards []map[string]interface{} `gorm:"type:json" json:"personalized_rewards"`

	// Engagement
	ProgramJoinDate   time.Time  `json:"program_join_date"`
	LastActivityDate  *time.Time `json:"last_activity_date"`
	ActivityStreak    int        `gorm:"default:0" json:"activity_streak"`
	EngagementLevel   string     `gorm:"type:varchar(20)" json:"engagement_level"`
	ParticipationRate float64    `gorm:"default:0" json:"participation_rate"`

	// Value Metrics
	CustomerValue  decimal.Decimal `gorm:"type:decimal(15,2)" json:"customer_value"`
	RewardValue    decimal.Decimal `gorm:"type:decimal(10,2)" json:"reward_value"`
	DiscountsSaved decimal.Decimal `gorm:"type:decimal(10,2)" json:"discounts_saved"`
	ROIScore       float64         `gorm:"default:0" json:"roi_score"`
}

// TableName returns the table name
func (UserDeviceEcosystem) TableName() string {
	return "user_device_ecosystem"
}

// TableName returns the table name
func (UserPolicyPortfolio) TableName() string {
	return "user_policy_portfolio"
}

// TableName returns the table name
func (UserClaimPatterns) TableName() string {
	return "user_claim_patterns"
}

// TableName returns the table name
func (UserPaymentEcosystem) TableName() string {
	return "user_payment_ecosystem"
}

// TableName returns the table name
func (UserReferralNetwork) TableName() string {
	return "user_referral_network"
}

// TableName returns the table name
func (UserLoyaltyJourney) TableName() string {
	return "user_loyalty_journey"
}
