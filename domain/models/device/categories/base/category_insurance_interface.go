package base

import (
	"time"
	
	"github.com/google/uuid"
)

// InsuranceCalculator defines the interface for category-specific insurance calculations
type InsuranceCalculator interface {
	// CalculatePremium calculates the insurance premium based on device specs
	CalculatePremium(spec CategorySpec, coverage CoverageLevel) (float64, error)

	// CalculateDeductible calculates the deductible amount
	CalculateDeductible(spec CategorySpec, claimType ClaimType) (float64, error)

	// AssessRisk performs risk assessment and returns risk score
	AssessRisk(spec CategorySpec, history ClaimHistory) (RiskAssessment, error)

	// GetCoverageOptions returns available coverage options for the category
	GetCoverageOptions(spec CategorySpec) []CoverageOption

	// ValidateEligibility checks if device is eligible for insurance
	ValidateEligibility(spec CategorySpec) (bool, []string)

	// CalculateMaxCoverage determines maximum insurable value
	CalculateMaxCoverage(spec CategorySpec) float64

	// GetExclusions returns policy exclusions for this category
	GetExclusions() []string

	// CalculateClaimPayout calculates the claim payout amount
	CalculateClaimPayout(spec CategorySpec, claim ClaimDetails) (float64, error)

	// GetPremiumFactors returns factors affecting premium calculation
	GetPremiumFactors(spec CategorySpec) map[string]float64
}

// CoverageLevel represents insurance coverage levels
type CoverageLevel string

const (
	CoverageBasic         CoverageLevel = "basic"
	CoverageStandard      CoverageLevel = "standard"
	CoveragePremium       CoverageLevel = "premium"
	CoverageComprehensive CoverageLevel = "comprehensive"
)

// ClaimType represents types of insurance claims
type ClaimType string

const (
	ClaimTypeDamage      ClaimType = "damage"
	ClaimTypeTheft       ClaimType = "theft"
	ClaimTypeLoss        ClaimType = "loss"
	ClaimTypeMalfunction ClaimType = "malfunction"
	ClaimTypeScreen      ClaimType = "screen"
	ClaimTypeBattery     ClaimType = "battery"
	ClaimTypeWater       ClaimType = "water"
)

// RiskAssessment contains risk evaluation results
type RiskAssessment struct {
	Score           float64            `json:"score"`
	Level           string             `json:"level"`
	Factors         map[string]float64 `json:"factors"`
	Recommendations []string           `json:"recommendations"`
	ValidUntil      time.Time          `json:"valid_until"`
	RequiresReview  bool               `json:"requires_review"`
}

// CoverageOption represents an insurance coverage option
type CoverageOption struct {
	ID             uuid.UUID     `json:"id"`
	Name           string        `json:"name"`
	Level          CoverageLevel `json:"level"`
	MonthlyPremium float64       `json:"monthly_premium"`
	AnnualPremium  float64       `json:"annual_premium"`
	Deductible     float64       `json:"deductible"`
	MaxCoverage    float64       `json:"max_coverage"`
	CoveredPerils  []ClaimType   `json:"covered_perils"`
	Exclusions     []string      `json:"exclusions"`
	WaitingPeriod  int           `json:"waiting_period_days"`
	ClaimLimit     int           `json:"claim_limit_per_year"`
	Features       []string      `json:"features"`
}

// ClaimHistory represents historical claim data
type ClaimHistory struct {
	TotalClaims     int               `json:"total_claims"`
	ApprovedClaims  int               `json:"approved_claims"`
	RejectedClaims  int               `json:"rejected_claims"`
	TotalPayout     float64           `json:"total_payout"`
	LastClaimDate   *time.Time        `json:"last_claim_date"`
	ClaimsByType    map[ClaimType]int `json:"claims_by_type"`
	AverageInterval int               `json:"average_interval_days"`
	FraudSuspicion  bool              `json:"fraud_suspicion"`
}

// ClaimDetails contains claim information
type ClaimDetails struct {
	ID             uuid.UUID `json:"id"`
	Type           ClaimType `json:"type"`
	Date           time.Time `json:"date"`
	Description    string    `json:"description"`
	EstimatedCost  float64   `json:"estimated_cost"`
	SupportingDocs []string  `json:"supporting_docs"`
	Location       string    `json:"location"`
	Circumstances  string    `json:"circumstances"`
	IsPreExisting  bool      `json:"is_pre_existing"`
}

// BaseInsuranceCalculator provides common implementation for insurance calculations
type BaseInsuranceCalculator struct {
	BasePremiumRate    float64
	CategoryMultiplier float64
	RiskThresholds     map[string]float64
}

// GetBasePremium calculates base premium
func (b *BaseInsuranceCalculator) GetBasePremium(marketValue float64) float64 {
	return marketValue * b.BasePremiumRate * b.CategoryMultiplier
}

// ApplyRiskMultiplier applies risk factor to premium
func (b *BaseInsuranceCalculator) ApplyRiskMultiplier(basePremium float64, riskScore float64) float64 {
	riskMultiplier := 1.0
	if riskScore > 80 {
		riskMultiplier = 2.0
	} else if riskScore > 60 {
		riskMultiplier = 1.5
	} else if riskScore > 40 {
		riskMultiplier = 1.2
	}
	return basePremium * riskMultiplier
}

// CalculateAgeDepreciation calculates depreciation based on age
func (b *BaseInsuranceCalculator) CalculateAgeDepreciation(releaseDate time.Time, depreciationRate float64) float64 {
	ageInYears := time.Since(releaseDate).Hours() / 24 / 365
	depreciation := 1.0 - (depreciationRate * ageInYears)
	if depreciation < 0.2 {
		depreciation = 0.2 // Minimum 20% of original value
	}
	return depreciation
}

// InsuranceProfile represents a device's insurance profile
type InsuranceProfile struct {
	DeviceID            uuid.UUID          `json:"device_id"`
	CategoryType        CategoryType       `json:"category_type"`
	EligibilityStatus   string             `json:"eligibility_status"`
	RiskAssessment      RiskAssessment     `json:"risk_assessment"`
	RecommendedCoverage CoverageLevel      `json:"recommended_coverage"`
	AvailableOptions    []CoverageOption   `json:"available_options"`
	PremiumFactors      map[string]float64 `json:"premium_factors"`
	Restrictions        []string           `json:"restrictions"`
	LastUpdated         time.Time          `json:"last_updated"`
}

// PremiumFactors contains factors affecting premium calculation
type PremiumFactors struct {
	BaseRate            float64 `json:"base_rate"`
	AgeMultiplier       float64 `json:"age_multiplier"`
	ConditionMultiplier float64 `json:"condition_multiplier"`
	LocationMultiplier  float64 `json:"location_multiplier"`
	HistoryMultiplier   float64 `json:"history_multiplier"`
	CategoryMultiplier  float64 `json:"category_multiplier"`
	SeasonalAdjustment  float64 `json:"seasonal_adjustment"`
	LoyaltyDiscount     float64 `json:"loyalty_discount"`
}
