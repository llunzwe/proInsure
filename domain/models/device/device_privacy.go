package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// DevicePrivacy represents privacy and consent management
type DevicePrivacy struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// GDPR Compliance
	GDPRConsent        bool       `gorm:"type:boolean;default:false" json:"gdpr_consent"`
	GDPRConsentDate    *time.Time `gorm:"type:timestamp" json:"gdpr_consent_date,omitempty"`
	GDPRConsentVersion string     `gorm:"type:varchar(50)" json:"gdpr_consent_version"`
	DataController     string     `gorm:"type:varchar(255)" json:"data_controller"`
	DataProcessor      string     `gorm:"type:varchar(255)" json:"data_processor"`
	LawfulBasis        string     `gorm:"type:varchar(100)" json:"lawful_basis"`

	// CCPA Compliance
	CCPAOptOut         bool       `gorm:"type:boolean;default:false" json:"ccpa_opt_out"`
	CCPAOptOutDate     *time.Time `gorm:"type:timestamp" json:"ccpa_opt_out_date,omitempty"`
	DoNotSellData      bool       `gorm:"type:boolean;default:false" json:"do_not_sell_data"`
	CaliforniaResident bool       `gorm:"type:boolean;default:false" json:"california_resident"`

	// Data Retention
	RetentionPeriodDays int            `gorm:"type:int" json:"retention_period_days"`
	DataExpirationDate  *time.Time     `gorm:"type:timestamp" json:"data_expiration_date,omitempty"`
	RetentionPolicy     datatypes.JSON `gorm:"type:json" json:"retention_policy"` // RetentionPolicy
	AutoDeleteEnabled   bool           `gorm:"type:boolean;default:false" json:"auto_delete_enabled"`
	DataCategories      datatypes.JSON `gorm:"type:json" json:"data_categories"` // []DataCategory

	// Consent Records
	ConsentRecords     datatypes.JSON `gorm:"type:json" json:"consent_records"`     // []ConsentRecord
	ConsentHistory     datatypes.JSON `gorm:"type:json" json:"consent_history"`     // []ConsentChange
	ActiveConsents     datatypes.JSON `gorm:"type:json" json:"active_consents"`     // map[string]bool
	ConsentGranularity datatypes.JSON `gorm:"type:json" json:"consent_granularity"` // map[string]ConsentDetail
	WithdrawnConsents  datatypes.JSON `gorm:"type:json" json:"withdrawn_consents"`  // []WithdrawnConsent

	// Right to be Forgotten
	DeletionRequests   datatypes.JSON `gorm:"type:json" json:"deletion_requests"` // []DeletionRequest
	DeletionStatus     string         `gorm:"type:varchar(50)" json:"deletion_status"`
	DeletionScheduled  *time.Time     `gorm:"type:timestamp" json:"deletion_scheduled,omitempty"`
	PartialDeletion    bool           `gorm:"type:boolean;default:false" json:"partial_deletion"`
	DeletionExceptions datatypes.JSON `gorm:"type:json" json:"deletion_exceptions"` // []Exception

	// Data Portability
	PortabilityRequests datatypes.JSON `gorm:"type:json" json:"portability_requests"` // []PortabilityRequest
	DataExports         datatypes.JSON `gorm:"type:json" json:"data_exports"`         // []Export
	ExportFormat        string         `gorm:"type:varchar(50)" json:"export_format"` // json, csv, xml
	LastExportDate      *time.Time     `gorm:"type:timestamp" json:"last_export_date,omitempty"`

	// Access Rights
	AccessRequests     datatypes.JSON `gorm:"type:json" json:"access_requests"`     // []AccessRequest
	DataAccessLog      datatypes.JSON `gorm:"type:json" json:"data_access_log"`     // []AccessLog
	AccessRestrictions datatypes.JSON `gorm:"type:json" json:"access_restrictions"` // []Restriction

	// Privacy Preferences
	PrivacySettings   datatypes.JSON `gorm:"type:json" json:"privacy_settings"`         // map[string]interface{}
	DataSharingPrefs  datatypes.JSON `gorm:"type:json" json:"data_sharing_preferences"` // map[string]bool
	MarketingConsent  bool           `gorm:"type:boolean;default:false" json:"marketing_consent"`
	AnalyticsConsent  bool           `gorm:"type:boolean;default:false" json:"analytics_consent"`
	ThirdPartyConsent bool           `gorm:"type:boolean;default:false" json:"third_party_consent"`

	// Data Minimization
	MinimalDataMode     bool           `gorm:"type:boolean;default:false" json:"minimal_data_mode"`
	AnonymizationLevel  string         `gorm:"type:varchar(50)" json:"anonymization_level"` // none, partial, full
	PseudonymizationKey string         `gorm:"type:varchar(255)" json:"pseudonymization_key"`
	DataMaskingRules    datatypes.JSON `gorm:"type:json" json:"data_masking_rules"` // []MaskingRule

	// Audit & Compliance
	PrivacyAudits       datatypes.JSON `gorm:"type:json" json:"privacy_audits"` // []Audit
	ComplianceStatus    string         `gorm:"type:varchar(50)" json:"compliance_status"`
	LastComplianceCheck time.Time      `gorm:"type:timestamp" json:"last_compliance_check"`
	ViolationReports    datatypes.JSON `gorm:"type:json" json:"violation_reports"` // []Violation

	// Status
	PrivacyStatus string    `gorm:"type:varchar(50)" json:"privacy_status"`
	LastUpdated   time.Time `gorm:"type:timestamp" json:"last_updated"`
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceAdvancedPricing represents advanced pricing models
type DeviceAdvancedPricing struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Dynamic Pricing History
	PricingHistory    datatypes.JSON `gorm:"type:json" json:"pricing_history"` // []PricePoint
	CurrentPrice      float64        `gorm:"type:decimal(15,2)" json:"current_price"`
	BasePrice         float64        `gorm:"type:decimal(15,2)" json:"base_price"`
	PriceChangeDate   time.Time      `gorm:"type:timestamp" json:"price_change_date"`
	PriceChangeReason string         `gorm:"type:text" json:"price_change_reason"`

	// Surge Pricing
	SurgeEnabled       bool           `gorm:"type:boolean;default:false" json:"surge_enabled"`
	SurgeMultiplier    float64        `gorm:"type:decimal(5,2)" json:"surge_multiplier"`
	SurgeConditions    datatypes.JSON `gorm:"type:json" json:"surge_conditions"` // []Condition
	SurgeHistory       datatypes.JSON `gorm:"type:json" json:"surge_history"`    // []SurgeEvent
	CurrentSurgeLevel  float64        `gorm:"type:decimal(5,2)" json:"current_surge_level"`
	MaxSurgeMultiplier float64        `gorm:"type:decimal(5,2)" json:"max_surge_multiplier"`

	// Bundle Pricing
	BundleDiscounts    datatypes.JSON `gorm:"type:json" json:"bundle_discounts"` // []BundleDiscount
	ActiveBundles      datatypes.JSON `gorm:"type:json" json:"active_bundles"`   // []Bundle
	BundleSavings      float64        `gorm:"type:decimal(15,2)" json:"bundle_savings"`
	BundleOptimization datatypes.JSON `gorm:"type:json" json:"bundle_optimization"` // Optimization
	CrossSellBundles   datatypes.JSON `gorm:"type:json" json:"cross_sell_bundles"`  // []Bundle

	// Loyalty Pricing Tiers
	LoyaltyTier        string         `gorm:"type:varchar(50)" json:"loyalty_tier"` // bronze, silver, gold, platinum
	TierDiscount       float64        `gorm:"type:decimal(5,2)" json:"tier_discount"`
	TierBenefits       datatypes.JSON `gorm:"type:json" json:"tier_benefits"` // []Benefit
	PointsToNextTier   int            `gorm:"type:int" json:"points_to_next_tier"`
	TierExpirationDate *time.Time     `gorm:"type:timestamp" json:"tier_expiration_date,omitempty"`

	// Geographic Pricing
	GeographicZone     string         `gorm:"type:varchar(100)" json:"geographic_zone"`
	ZonePricing        datatypes.JSON `gorm:"type:json" json:"zone_pricing"` // map[string]float64
	RegionalAdjustment float64        `gorm:"type:decimal(5,2)" json:"regional_adjustment"`
	LocalTaxes         float64        `gorm:"type:decimal(15,2)" json:"local_taxes"`
	CurrencyAdjustment float64        `gorm:"type:decimal(5,2)" json:"currency_adjustment"`

	// Time-Based Pricing
	TimeBasedPricing bool           `gorm:"type:boolean;default:false" json:"time_based_pricing"`
	PeakHoursPricing datatypes.JSON `gorm:"type:json" json:"peak_hours_pricing"` // []PeakPricing
	SeasonalPricing  datatypes.JSON `gorm:"type:json" json:"seasonal_pricing"`   // []SeasonalPrice
	WeekendPricing   float64        `gorm:"type:decimal(15,2)" json:"weekend_pricing"`
	HolidayPricing   datatypes.JSON `gorm:"type:json" json:"holiday_pricing"` // map[string]float64

	// Promotional Pricing
	ActivePromotions    datatypes.JSON `gorm:"type:json" json:"active_promotions"`  // []Promotion
	PromoCodeHistory    datatypes.JSON `gorm:"type:json" json:"promo_code_history"` // []PromoCode
	DiscountAmount      float64        `gorm:"type:decimal(15,2)" json:"discount_amount"`
	DiscountPercentage  float64        `gorm:"type:decimal(5,2)" json:"discount_percentage"`
	PromoExpirationDate *time.Time     `gorm:"type:timestamp" json:"promo_expiration_date,omitempty"`

	// Competition-Based Pricing
	CompetitorPrices datatypes.JSON `gorm:"type:json" json:"competitor_prices"` // map[string]float64
	PriceMatching    bool           `gorm:"type:boolean;default:false" json:"price_matching"`
	PricePosition    string         `gorm:"type:varchar(50)" json:"price_position"` // lowest, competitive, premium
	PriceElasticity  float64        `gorm:"type:decimal(10,4)" json:"price_elasticity"`
	OptimalPrice     float64        `gorm:"type:decimal(15,2)" json:"optimal_price"`

	// Personalized Pricing
	PersonalizedPrice float64        `gorm:"type:decimal(15,2)" json:"personalized_price"`
	CustomerSegment   string         `gorm:"type:varchar(100)" json:"customer_segment"`
	WillingnessToPay  float64        `gorm:"type:decimal(15,2)" json:"willingness_to_pay"`
	PricePreferences  datatypes.JSON `gorm:"type:json" json:"price_preferences"` // Preferences
	PriceSensitivity  float64        `gorm:"type:decimal(5,2)" json:"price_sensitivity"`

	// A/B Testing
	PriceExperiments datatypes.JSON `gorm:"type:json" json:"price_experiments"` // []Experiment
	TestGroup        string         `gorm:"type:varchar(50)" json:"test_group"`
	ExperimentPrice  float64        `gorm:"type:decimal(15,2)" json:"experiment_price"`
	ConversionRate   float64        `gorm:"type:decimal(5,2)" json:"conversion_rate"`

	// Status
	PricingStatus   string    `gorm:"type:varchar(50)" json:"pricing_status"`
	LastCalculation time.Time `gorm:"type:timestamp" json:"last_calculation"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// HasValidConsent checks if device has valid privacy consent
func (dp *DevicePrivacy) HasValidConsent() bool {
	return dp.GDPRConsent || (!dp.CCPAOptOut && !dp.DoNotSellData)
}

// IsCompliant checks privacy compliance
func (dp *DevicePrivacy) IsCompliant() bool {
	return dp.ComplianceStatus == "compliant" &&
		dp.ViolationReports == nil
}

// RequiresDeletion checks if data requires deletion
func (dp *DevicePrivacy) RequiresDeletion() bool {
	if dp.DeletionScheduled != nil {
		return time.Now().After(*dp.DeletionScheduled)
	}
	if dp.DataExpirationDate != nil {
		return time.Now().After(*dp.DataExpirationDate)
	}
	return dp.DeletionStatus == "pending"
}

// IsMinimalDataMode checks if minimal data collection is enabled
func (dp *DevicePrivacy) IsMinimalDataMode() bool {
	return dp.MinimalDataMode || dp.AnonymizationLevel == "full"
}

// HasMarketingConsent checks marketing consent
func (dp *DevicePrivacy) HasMarketingConsent() bool {
	return dp.MarketingConsent && !dp.CCPAOptOut && dp.GDPRConsent
}

// HasSurgePricing checks if surge pricing is active
func (dap *DeviceAdvancedPricing) HasSurgePricing() bool {
	return dap.SurgeEnabled && dap.CurrentSurgeLevel > 1.0
}

// GetEffectivePrice calculates the effective price
func (dap *DeviceAdvancedPricing) GetEffectivePrice() float64 {
	price := dap.CurrentPrice

	// Apply surge multiplier
	if dap.HasSurgePricing() {
		price = price * dap.CurrentSurgeLevel
	}

	// Apply tier discount
	if dap.TierDiscount > 0 {
		price = price * (1 - dap.TierDiscount/100)
	}

	// Apply regional adjustment
	price = price * (1 + dap.RegionalAdjustment/100)

	// Apply discounts
	if dap.DiscountPercentage > 0 {
		price = price * (1 - dap.DiscountPercentage/100)
	} else if dap.DiscountAmount > 0 {
		price = price - dap.DiscountAmount
	}

	return price
}

// IsEligibleForLoyaltyDiscount checks loyalty discount eligibility
func (dap *DeviceAdvancedPricing) IsEligibleForLoyaltyDiscount() bool {
	return dap.LoyaltyTier != "" && dap.LoyaltyTier != "bronze"
}

// HasActivePromotion checks for active promotions
func (dap *DeviceAdvancedPricing) HasActivePromotion() bool {
	if dap.PromoExpirationDate != nil {
		return time.Now().Before(*dap.PromoExpirationDate)
	}
	return dap.ActivePromotions != nil
}

// IsPriceCompetitive checks if price is competitive
func (dap *DeviceAdvancedPricing) IsPriceCompetitive() bool {
	return dap.PricePosition == "lowest" || dap.PricePosition == "competitive"
}
