package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceLifecycle represents the complete lifecycle tracking of a device
type DeviceLifecycle struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	CurrentStage   string    `gorm:"type:varchar(50);default:'new'" json:"current_stage"` // new, active, aging, end_of_life, recycled, disposed
	CurrentOwnerID uuid.UUID `gorm:"type:uuid" json:"current_owner_id"`

	// Ownership History
	TotalOwners              int        `json:"total_owners"`
	OwnershipTransfers       int        `json:"ownership_transfers"`
	FirstActivation          time.Time  `json:"first_activation"`
	LastTransferDate         *time.Time `json:"last_transfer_date"`
	AverageOwnershipDuration int        `json:"average_ownership_duration"` // in days

	// Claims & Repairs
	TotalClaimsCount  int     `json:"total_claims_count"`
	ApprovedClaims    int     `json:"approved_claims"`
	RejectedClaims    int     `json:"rejected_claims"`
	TotalClaimsValue  float64 `json:"total_claims_value"`
	TotalRepairsCount int     `json:"total_repairs_count"`
	MajorRepairs      int     `json:"major_repairs"`
	MinorRepairs      int     `json:"minor_repairs"`
	TotalRepairCost   float64 `json:"total_repair_cost"`

	// Financial Metrics
	TotalRevenueGenerated float64 `json:"total_revenue_generated"`
	PremiumsCollected     float64 `json:"premiums_collected"`
	ClaimsPaidOut         float64 `json:"claims_paid_out"`
	RepairRevenue         float64 `json:"repair_revenue"`
	RentalRevenue         float64 `json:"rental_revenue"`
	FinancingRevenue      float64 `json:"financing_revenue"`
	TradeInRevenue        float64 `json:"trade_in_revenue"`
	NetProfitability      float64 `json:"net_profitability"`

	// Environmental & Sustainability
	CarbonFootprint      float64 `json:"carbon_footprint"` // in kg CO2
	RecyclingEligible    bool    `gorm:"default:false" json:"recycling_eligible"`
	RecyclingValue       float64 `json:"recycling_value"`
	DisposalMethod       string  `json:"disposal_method"`       // recycle, refurbish, donate, destroy
	EnvironmentalScore   float64 `json:"environmental_score"`   // 0-100
	SustainabilityRating string  `json:"sustainability_rating"` // A+, A, B, C, D

	// Usage Metrics
	TotalActiveMonths  int        `json:"total_active_months"`
	InactiveMonths     int        `json:"inactive_months"`
	UsageIntensity     string     `json:"usage_intensity"` // light, moderate, heavy
	PeakValueDate      *time.Time `json:"peak_value_date"`
	PeakValue          float64    `json:"peak_value"`
	CurrentMarketValue float64    `json:"current_market_value"`
	DepreciationRate   float64    `json:"depreciation_rate"` // % per year

	// Risk & Compliance
	RiskScore         float64 `json:"risk_score"`        // 0-100
	ComplianceStatus  string  `json:"compliance_status"` // compliant, non_compliant, review_needed
	RegulatoryIssues  int     `json:"regulatory_issues"`
	SecurityIncidents int     `json:"security_incidents"`
	FraudAttempts     int     `json:"fraud_attempts"`

	// Quality & Performance
	QualityScore            float64 `json:"quality_score"`              // 0-100
	ReliabilityScore        float64 `json:"reliability_score"`          // 0-100
	FailureRate             float64 `json:"failure_rate"`               // % of time with issues
	MeanTimeBetweenFailures int     `json:"mean_time_between_failures"` // in days
	CustomerSatisfaction    float64 `json:"customer_satisfaction"`      // 0-5

	// Milestones & Events
	WarrantyExpired    bool       `gorm:"default:false" json:"warranty_expired"`
	WarrantyExpiryDate *time.Time `json:"warranty_expiry_date"`
	EOLDate            *time.Time `json:"eol_date"` // End of Life date
	LastMajorUpdate    *time.Time `json:"last_major_update"`
	NextMilestone      string     `json:"next_milestone"`
	NextMilestoneDate  *time.Time `json:"next_milestone_date"`

	// Predictions
	PredictedEOL        *time.Time `json:"predicted_eol"`
	PredictedNextClaim  *time.Time `json:"predicted_next_claim"`
	PredictedNextRepair *time.Time `json:"predicted_next_repair"`
	ChurnRisk           float64    `json:"churn_risk"` // 0-100
	UpgradeRecommended  bool       `gorm:"default:false" json:"upgrade_recommended"`
}

// LifecycleEvent represents significant events in device lifecycle
type LifecycleEvent struct {
	database.BaseModel
	DeviceID         uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	EventType        string     `gorm:"not null" json:"event_type"` // activation, transfer, repair, claim, upgrade, recycle
	EventDate        time.Time  `json:"event_date"`
	EventDescription string     `json:"event_description"`
	EventValue       float64    `json:"event_value"` // Financial impact
	UserID           *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	RelatedEntityID  *uuid.UUID `gorm:"type:uuid" json:"related_entity_id"` // claim_id, repair_id, etc.
	ImpactScore      float64    `json:"impact_score"`                       // 0-100 impact on lifecycle

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// OwnershipHistory tracks device ownership changes
type OwnershipHistory struct {
	database.BaseModel
	DeviceID          uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	PreviousOwnerID   *uuid.UUID `gorm:"type:uuid" json:"previous_owner_id"`
	NewOwnerID        uuid.UUID  `gorm:"type:uuid;not null" json:"new_owner_id"`
	TransferDate      time.Time  `json:"transfer_date"`
	TransferType      string     `json:"transfer_type"` // sale, gift, inheritance, trade_in
	TransferValue     float64    `json:"transfer_value"`
	OwnershipDuration int        `json:"ownership_duration"` // days with previous owner
	DeviceCondition   string     `json:"device_condition"`   // at transfer

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// PreviousOwner and NewOwner should be loaded via service layer using PreviousOwnerID and NewOwnerID to avoid circular import
}

// TableName returns the table name
func (t *DeviceLifecycle) TableName() string {
	return "device_lifecycles"
}

func (t *LifecycleEvent) TableName() string {
	return "lifecycle_events"
}

func (t *OwnershipHistory) TableName() string {
	return "ownership_histories"
}

// BeforeCreate handles pre-creation logic
func (dl *DeviceLifecycle) BeforeCreate(tx *gorm.DB) error {
	if err := dl.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// UpdateStage updates the lifecycle stage
func (dl *DeviceLifecycle) UpdateStage(newStage string) {
	dl.CurrentStage = newStage

	if newStage == "end_of_life" && dl.EOLDate == nil {
		now := time.Now()
		dl.EOLDate = &now
	}
}

// CalculateNetProfitability calculates net profitability
func (dl *DeviceLifecycle) CalculateNetProfitability() {
	totalRevenue := dl.PremiumsCollected + dl.RepairRevenue + dl.RentalRevenue +
		dl.FinancingRevenue + dl.TradeInRevenue
	totalCosts := dl.ClaimsPaidOut + dl.TotalRepairCost
	dl.NetProfitability = totalRevenue - totalCosts
	dl.TotalRevenueGenerated = totalRevenue
}

// CalculateEnvironmentalScore calculates environmental impact score
func (dl *DeviceLifecycle) CalculateEnvironmentalScore() {
	score := 100.0

	// Deduct for carbon footprint
	score -= dl.CarbonFootprint * 0.1

	// Deduct for repairs (resource consumption)
	score -= float64(dl.TotalRepairsCount) * 2

	// Add points for longevity
	score += float64(dl.TotalActiveMonths) * 0.5

	// Ensure score is between 0 and 100
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	dl.EnvironmentalScore = score

	// Set rating based on score
	switch {
	case score >= 90:
		dl.SustainabilityRating = "A+"
	case score >= 80:
		dl.SustainabilityRating = "A"
	case score >= 70:
		dl.SustainabilityRating = "B"
	case score >= 60:
		dl.SustainabilityRating = "C"
	default:
		dl.SustainabilityRating = "D"
	}
}

// CalculateMTBF calculates Mean Time Between Failures
func (dl *DeviceLifecycle) CalculateMTBF() {
	if dl.TotalRepairsCount == 0 {
		dl.MeanTimeBetweenFailures = dl.TotalActiveMonths * 30 // No failures
		return
	}

	dl.MeanTimeBetweenFailures = (dl.TotalActiveMonths * 30) / dl.TotalRepairsCount
}

// AddOwnershipTransfer records a new ownership transfer
func (dl *DeviceLifecycle) AddOwnershipTransfer() {
	dl.TotalOwners++
	dl.OwnershipTransfers++
	now := time.Now()
	dl.LastTransferDate = &now

	// Recalculate average ownership duration
	if dl.TotalOwners > 0 {
		totalDays := int(time.Since(dl.FirstActivation).Hours() / 24)
		dl.AverageOwnershipDuration = totalDays / dl.TotalOwners
	}
}

// ShouldUpgrade determines if device should be upgraded
func (dl *DeviceLifecycle) ShouldUpgrade() bool {
	// Check multiple factors
	if dl.CurrentStage == "end_of_life" {
		return true
	}

	if dl.TotalRepairsCount > 5 || dl.MajorRepairs > 2 {
		return true
	}

	if dl.FailureRate > 0.3 { // More than 30% failure rate
		return true
	}

	if dl.CurrentMarketValue < dl.PeakValue*0.2 { // Less than 20% of peak value
		return true
	}

	dl.UpgradeRecommended = true
	return true
}

// IsEndOfLife checks if device has reached end of life
func (dl *DeviceLifecycle) IsEndOfLife() bool {
	if dl.EOLDate != nil && time.Now().After(*dl.EOLDate) {
		return true
	}

	// Check other EOL criteria
	if dl.TotalActiveMonths > 60 { // 5 years
		return true
	}

	if dl.CurrentMarketValue <= 0 {
		return true
	}

	return dl.CurrentStage == "end_of_life"
}
