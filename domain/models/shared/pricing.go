package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// UsageBasedInsurance represents usage-based insurance policies with dynamic pricing
type UsageBasedInsurance struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PolicyID               uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`
	UserID                 uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	DeviceID               uuid.UUID      `gorm:"type:uuid;not null" json:"device_id"`
	PricingModel           string         `gorm:"not null" json:"pricing_model"` // pay_per_use, risk_based, seasonal, activity_based
	BasePremium            float64        `gorm:"not null" json:"base_premium"`
	CurrentPremium         float64        `gorm:"not null" json:"current_premium"`
	UsageMetrics           string         `json:"usage_metrics"`                    // JSON object with usage data
	RiskFactors            string         `json:"risk_factors"`                     // JSON object with risk factors
	ActivityPatterns       string         `json:"activity_patterns"`                // JSON object
	LocationRisk           float64        `gorm:"default:1.0" json:"location_risk"` // multiplier
	TimeRisk               float64        `gorm:"default:1.0" json:"time_risk"`     // multiplier
	BehaviorRisk           float64        `gorm:"default:1.0" json:"behavior_risk"` // multiplier
	LastRecalculation      time.Time      `json:"last_recalculation"`
	NextRecalculation      time.Time      `json:"next_recalculation"`
	RecalculationFrequency string         `gorm:"default:'monthly'" json:"recalculation_frequency"` // daily, weekly, monthly, quarterly
	IsActive               bool           `gorm:"default:true" json:"is_active"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Policy       models.Policy `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	User         models.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device       models.Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	UsageRecords []UsageRecord `gorm:"foreignKey:UBIId" json:"usage_records,omitempty"`
}

// UsageRecord represents individual usage records for UBI calculation
type UsageRecord struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UBIId      uuid.UUID `gorm:"type:uuid;not null" json:"ubi_id"`
	RecordType string    `gorm:"not null" json:"record_type"` // location, activity, time, behavior
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`
	Location   string    `json:"location"`
	Latitude   *float64  `json:"latitude"`
	Longitude  *float64  `json:"longitude"`
	Activity   string    `json:"activity"` // travel, sports, work, home
	RiskScore  float64   `json:"risk_score"`
	Duration   int       `json:"duration"` // minutes
	Metadata   string    `json:"metadata"` // JSON object with additional data
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	UBI UsageBasedInsurance `gorm:"foreignKey:UBIId" json:"ubi,omitempty"`
}

// MicroInsurance represents short-term and event-specific insurance products
type MicroInsurance struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	DeviceID       uuid.UUID      `gorm:"type:uuid;not null" json:"device_id"`
	ProductType    string         `gorm:"not null" json:"product_type"` // travel, event, temporary, activity
	EventName      string         `json:"event_name"`
	EventLocation  string         `json:"event_location"`
	CoverageType   string         `gorm:"not null" json:"coverage_type"` // theft, damage, loss, all
	Premium        float64        `gorm:"not null" json:"premium"`
	CoverageAmount float64        `gorm:"not null" json:"coverage_amount"`
	Deductible     float64        `gorm:"default:0" json:"deductible"`
	StartDate      time.Time      `gorm:"not null" json:"start_date"`
	EndDate        time.Time      `gorm:"not null" json:"end_date"`
	Duration       int            `json:"duration"`                                // hours
	Status         string         `gorm:"not null;default:'active'" json:"status"` // active, expired, cancelled, claimed
	PurchaseMethod string         `json:"purchase_method"`                         // app, web, api, partner
	PaymentStatus  string         `gorm:"default:'pending'" json:"payment_status"` // pending, paid, failed, refunded
	Terms          string         `json:"terms"`                                   // JSON object with terms and conditions
	Restrictions   string         `json:"restrictions"`                            // JSON object with restrictions
	IsAutoRenew    bool           `gorm:"default:false" json:"is_auto_renew"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User   models.User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device models.Device  `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	Claims []models.Claim `gorm:"foreignKey:MicroInsuranceID" json:"claims,omitempty"`
}

// PricingRule represents dynamic pricing rules for different scenarios
type PricingRule struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	RuleName           string         `gorm:"not null;uniqueIndex" json:"rule_name"`
	RuleType           string         `gorm:"not null" json:"rule_type"` // location, time, activity, device, user_profile
	Description        string         `json:"description"`
	Conditions         string         `gorm:"not null" json:"conditions"` // JSON object with conditions
	PriceMultiplier    float64        `gorm:"not null" json:"price_multiplier"`
	Priority           int            `gorm:"default:0" json:"priority"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	ValidFrom          time.Time      `json:"valid_from"`
	ValidTo            *time.Time     `json:"valid_to"`
	ApplicableProducts string         `json:"applicable_products"` // JSON array of product types
	CreatedBy          uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// SeasonalPricing represents seasonal pricing adjustments
type SeasonalPricing struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Season          string         `gorm:"not null" json:"season"` // spring, summer, autumn, winter, holiday
	StartDate       time.Time      `gorm:"not null" json:"start_date"`
	EndDate         time.Time      `gorm:"not null" json:"end_date"`
	PriceAdjustment float64        `gorm:"not null" json:"price_adjustment"` // multiplier
	CoverageTypes   string         `json:"coverage_types"`                   // JSON array of applicable coverage types
	Regions         string         `json:"regions"`                          // JSON array of applicable regions
	Description     string         `json:"description"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// ActivityRisk represents risk levels for different activities
type ActivityRisk struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ActivityName       string         `gorm:"not null;uniqueIndex" json:"activity_name"`
	Category           string         `gorm:"not null" json:"category"`   // sports, travel, work, entertainment
	RiskLevel          string         `gorm:"not null" json:"risk_level"` // low, medium, high, extreme
	RiskScore          float64        `gorm:"not null" json:"risk_score"` // 0-100
	PriceMultiplier    float64        `gorm:"not null" json:"price_multiplier"`
	Description        string         `json:"description"`
	SafetyRequirements string         `json:"safety_requirements"` // JSON array
	ExcludedCoverage   string         `json:"excluded_coverage"`   // JSON array
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName methods
func (UsageBasedInsurance) TableName() string {
	return "usage_based_insurances"
}

func (UsageRecord) TableName() string {
	return "usage_records"
}

func (MicroInsurance) TableName() string {
	return "micro_insurances"
}

func (PricingRule) TableName() string {
	return "pricing_rules"
}

func (SeasonalPricing) TableName() string {
	return "seasonal_pricings"
}

func (ActivityRisk) TableName() string {
	return "activity_risks"
}

// BeforeCreate hooks
func (ubi *UsageBasedInsurance) BeforeCreate(tx *gorm.DB) error {
	if ubi.ID == uuid.Nil {
		ubi.ID = uuid.New()
	}
	return nil
}

func (ur *UsageRecord) BeforeCreate(tx *gorm.DB) error {
	if ur.ID == uuid.Nil {
		ur.ID = uuid.New()
	}
	return nil
}

func (mi *MicroInsurance) BeforeCreate(tx *gorm.DB) error {
	if mi.ID == uuid.Nil {
		mi.ID = uuid.New()
	}
	return nil
}

func (pr *PricingRule) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}

func (sp *SeasonalPricing) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

func (ar *ActivityRisk) BeforeCreate(tx *gorm.DB) error {
	if ar.ID == uuid.Nil {
		ar.ID = uuid.New()
	}
	return nil
}

// Business logic methods for UsageBasedInsurance
func (ubi *UsageBasedInsurance) CalculatePremium() float64 {
	return ubi.BasePremium * ubi.LocationRisk * ubi.TimeRisk * ubi.BehaviorRisk
}

func (ubi *UsageBasedInsurance) UpdatePremium() {
	ubi.CurrentPremium = ubi.CalculatePremium()
	ubi.LastRecalculation = time.Now()

	// Set next recalculation based on frequency
	switch ubi.RecalculationFrequency {
	case "daily":
		ubi.NextRecalculation = time.Now().AddDate(0, 0, 1)
	case "weekly":
		ubi.NextRecalculation = time.Now().AddDate(0, 0, 7)
	case "monthly":
		ubi.NextRecalculation = time.Now().AddDate(0, 1, 0)
	case "quarterly":
		ubi.NextRecalculation = time.Now().AddDate(0, 3, 0)
	default:
		ubi.NextRecalculation = time.Now().AddDate(0, 1, 0)
	}
}

func (ubi *UsageBasedInsurance) IsRecalculationDue() bool {
	return time.Now().After(ubi.NextRecalculation)
}

// Business logic methods for MicroInsurance
func (mi *MicroInsurance) IsActive() bool {
	now := time.Now()
	return mi.Status == "active" && now.After(mi.StartDate) && now.Before(mi.EndDate)
}

func (mi *MicroInsurance) IsExpired() bool {
	return time.Now().After(mi.EndDate)
}

func (mi *MicroInsurance) Cancel() {
	mi.Status = "cancelled"
}

func (mi *MicroInsurance) Expire() {
	mi.Status = "expired"
}

func (mi *MicroInsurance) GetRemainingDuration() time.Duration {
	if mi.IsExpired() {
		return 0
	}
	return time.Until(mi.EndDate)
}

// Business logic methods for PricingRule
func (pr *PricingRule) IsValid() bool {
	now := time.Now()
	if pr.ValidTo == nil {
		return pr.IsActive && now.After(pr.ValidFrom)
	}
	return pr.IsActive && now.After(pr.ValidFrom) && now.Before(*pr.ValidTo)
}

func (pr *PricingRule) Deactivate() {
	pr.IsActive = false
}

// Business logic methods for SeasonalPricing
func (sp *SeasonalPricing) IsCurrentlyActive() bool {
	now := time.Now()
	return sp.IsActive && now.After(sp.StartDate) && now.Before(sp.EndDate)
}

// Business logic methods for ActivityRisk
func (ar *ActivityRisk) IsHighRisk() bool {
	return ar.RiskLevel == "high" || ar.RiskLevel == "extreme"
}

func (ar *ActivityRisk) RequiresSpecialTerms() bool {
	return ar.IsHighRisk()
}
