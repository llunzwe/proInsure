package shared

import (
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RiskProfile represents comprehensive risk assessment for entities
type RiskProfile struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProfileCode string    `gorm:"uniqueIndex;not null" json:"profile_code"`

	// Entity Reference (polymorphic)
	EntityType string    `gorm:"not null" json:"entity_type"` // user, device, policy, claim, partner
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`

	// Risk Scores (0-100)
	OverallRiskScore    float64 `gorm:"default:0" json:"overall_risk_score"`
	FraudRiskScore      float64 `gorm:"default:0" json:"fraud_risk_score"`
	CreditRiskScore     float64 `gorm:"default:0" json:"credit_risk_score"`
	ClaimRiskScore      float64 `gorm:"default:0" json:"claim_risk_score"`
	BehavioralRiskScore float64 `gorm:"default:0" json:"behavioral_risk_score"`
	DeviceRiskScore     float64 `gorm:"default:0" json:"device_risk_score"`
	GeographicRiskScore float64 `gorm:"default:0" json:"geographic_risk_score"`

	// Risk Level
	RiskLevel    string `gorm:"default:'low'" json:"risk_level"` // low, medium, high, critical
	RiskCategory string `json:"risk_category"`                   // standard, elevated, restricted, blacklisted

	// Risk Factors (JSON arrays of contributing factors)
	PositiveFactors string `gorm:"type:json" json:"positive_factors"`
	NegativeFactors string `gorm:"type:json" json:"negative_factors"`
	RiskIndicators  string `gorm:"type:json" json:"risk_indicators"`

	// Historical Data
	PreviousRiskScore float64 `json:"previous_risk_score"`
	ScoreChange       float64 `json:"score_change"`
	TrendDirection    string  `json:"trend_direction"` // improving, stable, deteriorating

	// Assessment Details
	LastAssessmentDate  time.Time `json:"last_assessment_date"`
	NextAssessmentDate  time.Time `json:"next_assessment_date"`
	AssessmentFrequency string    `gorm:"default:'monthly'" json:"assessment_frequency"` // daily, weekly, monthly, quarterly
	AssessmentMethod    string    `json:"assessment_method"`                             // automated, manual, hybrid

	// Actions & Recommendations
	RequiresReview  bool       `gorm:"default:false" json:"requires_review"`
	ReviewedBy      *uuid.UUID `gorm:"type:uuid" json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	ReviewNotes     string     `json:"review_notes"`
	Recommendations string     `gorm:"type:json" json:"recommendations"`

	// Impact on Business Rules
	PremiumAdjustment  float64 `json:"premium_adjustment"` // percentage adjustment
	CoverageLimit      float64 `json:"coverage_limit"`
	DeductibleIncrease float64 `json:"deductible_increase"`
	RequiresDeposit    bool    `gorm:"default:false" json:"requires_deposit"`
	DepositAmount      float64 `json:"deposit_amount"`

	// Monitoring
	MonitoringLevel string  `gorm:"default:'standard'" json:"monitoring_level"` // standard, enhanced, intensive
	AlertsEnabled   bool    `gorm:"default:true" json:"alerts_enabled"`
	AlertThreshold  float64 `gorm:"default:70" json:"alert_threshold"`

	// Metadata
	DataSources     string  `gorm:"type:json" json:"data_sources"` // JSON array of data sources used
	ConfidenceLevel float64 `json:"confidence_level"`              // 0-100% confidence in assessment

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	RiskFactors []RiskFactor `gorm:"foreignKey:RiskProfileID" json:"risk_factors,omitempty"`
	RiskEvents  []RiskEvent  `gorm:"foreignKey:RiskProfileID" json:"risk_events,omitempty"`
}

// RiskFactor represents individual risk factors contributing to the profile
type RiskFactor struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RiskProfileID uuid.UUID `gorm:"type:uuid;not null" json:"risk_profile_id"`

	// Factor Details
	FactorCode string `gorm:"not null" json:"factor_code"`
	FactorName string `gorm:"not null" json:"factor_name"`
	Category   string `gorm:"not null" json:"category"` // demographic, behavioral, historical, external
	Type       string `gorm:"not null" json:"type"`     // positive, negative, neutral

	// Impact
	Weight float64 `gorm:"default:1.0" json:"weight"`      // Factor weight in calculation
	Score  float64 `json:"score"`                          // Individual factor score
	Impact string  `gorm:"default:'medium'" json:"impact"` // low, medium, high

	// Value & Threshold
	CurrentValue     string `json:"current_value"`
	ThresholdValue   string `json:"threshold_value"`
	IsAboveThreshold bool   `gorm:"default:false" json:"is_above_threshold"`

	// Temporal
	ValidFrom  time.Time  `json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	RiskProfile RiskProfile `gorm:"foreignKey:RiskProfileID" json:"risk_profile,omitempty"`
}

// RiskEvent represents events that affect risk assessment
type RiskEvent struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RiskProfileID uuid.UUID `gorm:"type:uuid;not null" json:"risk_profile_id"`

	// Event Details
	EventType        string `gorm:"not null" json:"event_type"` // claim, payment_default, policy_lapse, fraud_alert
	EventDescription string `json:"event_description"`
	Severity         string `gorm:"default:'medium'" json:"severity"` // low, medium, high, critical

	// Impact
	ScoreImpact  float64 `json:"score_impact"`  // Change in risk score
	DurationDays int     `json:"duration_days"` // How long the impact lasts

	// Timing
	OccurredAt time.Time  `json:"occurred_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`

	// Source
	SourceType string     `json:"source_type"` // system, manual, external
	SourceID   *uuid.UUID `gorm:"type:uuid" json:"source_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	RiskProfile RiskProfile `gorm:"foreignKey:RiskProfileID" json:"risk_profile,omitempty"`
}

// RiskMatrix defines risk calculation rules and thresholds
type RiskMatrix struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MatrixCode  string    `gorm:"uniqueIndex;not null" json:"matrix_code"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`

	// Applicable To
	EntityType  string `gorm:"not null" json:"entity_type"` // user, device, policy, claim
	ProductType string `json:"product_type"`
	Region      string `json:"region"`

	// Risk Categories & Weights
	FraudWeight      float64 `gorm:"default:0.25" json:"fraud_weight"`
	CreditWeight     float64 `gorm:"default:0.20" json:"credit_weight"`
	ClaimWeight      float64 `gorm:"default:0.25" json:"claim_weight"`
	BehavioralWeight float64 `gorm:"default:0.15" json:"behavioral_weight"`
	DeviceWeight     float64 `gorm:"default:0.10" json:"device_weight"`
	GeographicWeight float64 `gorm:"default:0.05" json:"geographic_weight"`

	// Risk Level Thresholds
	LowThreshold      float64 `gorm:"default:30" json:"low_threshold"`
	MediumThreshold   float64 `gorm:"default:50" json:"medium_threshold"`
	HighThreshold     float64 `gorm:"default:70" json:"high_threshold"`
	CriticalThreshold float64 `gorm:"default:85" json:"critical_threshold"`

	// Business Rules
	Rules string `gorm:"type:json" json:"rules"` // JSON array of rule definitions

	// Premium Adjustments by Risk Level
	LowPremiumAdjustment      float64 `gorm:"default:0" json:"low_premium_adjustment"`
	MediumPremiumAdjustment   float64 `gorm:"default:10" json:"medium_premium_adjustment"`
	HighPremiumAdjustment     float64 `gorm:"default:25" json:"high_premium_adjustment"`
	CriticalPremiumAdjustment float64 `gorm:"default:50" json:"critical_premium_adjustment"`

	// Coverage Limits by Risk Level
	LowCoverageLimit      float64 `json:"low_coverage_limit"`
	MediumCoverageLimit   float64 `json:"medium_coverage_limit"`
	HighCoverageLimit     float64 `json:"high_coverage_limit"`
	CriticalCoverageLimit float64 `json:"critical_coverage_limit"`

	// Deductible Adjustments
	LowDeductible      float64 `json:"low_deductible"`
	MediumDeductible   float64 `json:"medium_deductible"`
	HighDeductible     float64 `json:"high_deductible"`
	CriticalDeductible float64 `json:"critical_deductible"`

	// Status
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveUntil *time.Time `json:"effective_until"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	RiskRules []RiskRule `gorm:"foreignKey:MatrixID" json:"risk_rules,omitempty"`
}

// RiskRule defines specific risk calculation rules
type RiskRule struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MatrixID uuid.UUID `gorm:"type:uuid;not null" json:"matrix_id"`

	// Rule Definition
	RuleCode    string `gorm:"uniqueIndex;not null" json:"rule_code"`
	RuleName    string `gorm:"not null" json:"rule_name"`
	Description string `json:"description"`
	Category    string `json:"category"`

	// Condition
	ConditionType     string `gorm:"not null" json:"condition_type"` // threshold, range, pattern, formula
	ConditionField    string `json:"condition_field"`
	ConditionOperator string `json:"condition_operator"` // >, <, >=, <=, ==, !=, contains, matches
	ConditionValue    string `json:"condition_value"`
	ConditionFormula  string `json:"condition_formula"` // Complex formula for calculation

	// Action
	ActionType      string  `gorm:"not null" json:"action_type"` // adjust_score, set_flag, require_review, reject
	ActionValue     string  `json:"action_value"`
	ScoreAdjustment float64 `json:"score_adjustment"`

	// Priority & Weight
	Priority int     `gorm:"default:5" json:"priority"`
	Weight   float64 `gorm:"default:1.0" json:"weight"`

	// Status
	IsActive    bool `gorm:"default:true" json:"is_active"`
	IsMandatory bool `gorm:"default:false" json:"is_mandatory"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Matrix RiskMatrix `gorm:"foreignKey:MatrixID" json:"matrix,omitempty"`
}

// RiskAssessmentLog tracks all risk assessments performed
type RiskAssessmentLog struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// Entity
	EntityType string    `gorm:"not null" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`

	// Assessment
	AssessmentType string `json:"assessment_type"` // scheduled, triggered, manual
	TriggerEvent   string `json:"trigger_event"`

	// Scores
	PreviousScore float64 `json:"previous_score"`
	NewScore      float64 `json:"new_score"`
	ScoreChange   float64 `json:"score_change"`

	// Risk Level
	PreviousLevel string `json:"previous_level"`
	NewLevel      string `json:"new_level"`

	// Factors
	FactorsAnalyzed int    `json:"factors_analyzed"`
	PositiveFactors string `gorm:"type:json" json:"positive_factors"`
	NegativeFactors string `gorm:"type:json" json:"negative_factors"`

	// Rules Applied
	RulesApplied   string `gorm:"type:json" json:"rules_applied"`
	RulesTriggered string `gorm:"type:json" json:"rules_triggered"`

	// Duration
	ProcessingTime int `json:"processing_time"` // milliseconds

	// Result
	ActionsTaken  string `gorm:"type:json" json:"actions_taken"`
	Notifications string `gorm:"type:json" json:"notifications"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Table names
func (RiskProfile) TableName() string       { return "risk_profiles" }
func (RiskFactor) TableName() string        { return "risk_factors" }
func (RiskEvent) TableName() string         { return "risk_events" }
func (RiskMatrix) TableName() string        { return "risk_matrices" }
func (RiskRule) TableName() string          { return "risk_rules" }
func (RiskAssessmentLog) TableName() string { return "risk_assessment_logs" }

// BeforeCreate hooks
func (rp *RiskProfile) BeforeCreate(tx *gorm.DB) error {
	if rp.ID == uuid.Nil {
		rp.ID = uuid.New()
	}
	if rp.ProfileCode == "" {
		rp.ProfileCode = "RISK-" + rp.EntityType + "-" + uuid.New().String()[:8]
	}
	rp.LastAssessmentDate = time.Now()
	rp.NextAssessmentDate = rp.CalculateNextAssessment()
	return nil
}

func (rf *RiskFactor) BeforeCreate(tx *gorm.DB) error {
	if rf.ID == uuid.Nil {
		rf.ID = uuid.New()
	}
	rf.ValidFrom = time.Now()
	return nil
}

func (re *RiskEvent) BeforeCreate(tx *gorm.DB) error {
	if re.ID == uuid.Nil {
		re.ID = uuid.New()
	}
	return nil
}

func (rm *RiskMatrix) BeforeCreate(tx *gorm.DB) error {
	if rm.ID == uuid.Nil {
		rm.ID = uuid.New()
	}
	if rm.MatrixCode == "" {
		rm.MatrixCode = "MATRIX-" + uuid.New().String()[:8]
	}
	return nil
}

func (rr *RiskRule) BeforeCreate(tx *gorm.DB) error {
	if rr.ID == uuid.Nil {
		rr.ID = uuid.New()
	}
	if rr.RuleCode == "" {
		rr.RuleCode = "RULE-" + uuid.New().String()[:8]
	}
	return nil
}

func (ral *RiskAssessmentLog) BeforeCreate(tx *gorm.DB) error {
	if ral.ID == uuid.Nil {
		ral.ID = uuid.New()
	}
	ral.ScoreChange = ral.NewScore - ral.PreviousScore
	return nil
}

// Business Logic Methods

// CalculateOverallRisk calculates the weighted overall risk score
func (rp *RiskProfile) CalculateOverallRisk(matrix *RiskMatrix) {
	overall := 0.0

	if matrix != nil {
		overall += rp.FraudRiskScore * matrix.FraudWeight
		overall += rp.CreditRiskScore * matrix.CreditWeight
		overall += rp.ClaimRiskScore * matrix.ClaimWeight
		overall += rp.BehavioralRiskScore * matrix.BehavioralWeight
		overall += rp.DeviceRiskScore * matrix.DeviceWeight
		overall += rp.GeographicRiskScore * matrix.GeographicWeight
	} else {
		// Default weights if no matrix provided
		overall = (rp.FraudRiskScore + rp.CreditRiskScore + rp.ClaimRiskScore +
			rp.BehavioralRiskScore + rp.DeviceRiskScore + rp.GeographicRiskScore) / 6
	}

	rp.OverallRiskScore = math.Min(100, math.Max(0, overall))
	rp.UpdateRiskLevel(matrix)
}

// UpdateRiskLevel updates the risk level based on score and matrix
func (rp *RiskProfile) UpdateRiskLevel(matrix *RiskMatrix) {
	if matrix == nil {
		// Default thresholds
		if rp.OverallRiskScore < 30 {
			rp.RiskLevel = "low"
		} else if rp.OverallRiskScore < 50 {
			rp.RiskLevel = "medium"
		} else if rp.OverallRiskScore < 70 {
			rp.RiskLevel = "high"
		} else {
			rp.RiskLevel = "critical"
		}
	} else {
		if rp.OverallRiskScore < matrix.LowThreshold {
			rp.RiskLevel = "low"
		} else if rp.OverallRiskScore < matrix.MediumThreshold {
			rp.RiskLevel = "medium"
		} else if rp.OverallRiskScore < matrix.HighThreshold {
			rp.RiskLevel = "high"
		} else {
			rp.RiskLevel = "critical"
		}
	}

	// Update risk category
	switch rp.RiskLevel {
	case "low":
		rp.RiskCategory = "standard"
	case "medium":
		rp.RiskCategory = "elevated"
	case "high":
		rp.RiskCategory = "restricted"
	case "critical":
		rp.RiskCategory = "blacklisted"
	}
}

// CalculateNextAssessment calculates the next assessment date
func (rp *RiskProfile) CalculateNextAssessment() time.Time {
	switch rp.AssessmentFrequency {
	case "daily":
		return time.Now().AddDate(0, 0, 1)
	case "weekly":
		return time.Now().AddDate(0, 0, 7)
	case "monthly":
		return time.Now().AddDate(0, 1, 0)
	case "quarterly":
		return time.Now().AddDate(0, 3, 0)
	default:
		return time.Now().AddDate(0, 1, 0) // Default monthly
	}
}

// IsHighRisk checks if profile is high risk
func (rp *RiskProfile) IsHighRisk() bool {
	return rp.RiskLevel == "high" || rp.RiskLevel == "critical"
}

// RequiresManualReview checks if manual review is required
func (rp *RiskProfile) RequiresManualReview() bool {
	return rp.RequiresReview || rp.ConfidenceLevel < 70 || rp.IsHighRisk()
}

// UpdateTrend updates the trend direction
func (rp *RiskProfile) UpdateTrend() {
	change := rp.OverallRiskScore - rp.PreviousRiskScore

	if math.Abs(change) < 5 {
		rp.TrendDirection = "stable"
	} else if change > 0 {
		rp.TrendDirection = "deteriorating"
	} else {
		rp.TrendDirection = "improving"
	}

	rp.ScoreChange = change
}

// ApplyRiskAdjustments applies risk-based adjustments
func (rp *RiskProfile) ApplyRiskAdjustments(matrix *RiskMatrix) {
	if matrix == nil {
		return
	}

	switch rp.RiskLevel {
	case "low":
		rp.PremiumAdjustment = matrix.LowPremiumAdjustment
		rp.CoverageLimit = matrix.LowCoverageLimit
		rp.DeductibleIncrease = matrix.LowDeductible
	case "medium":
		rp.PremiumAdjustment = matrix.MediumPremiumAdjustment
		rp.CoverageLimit = matrix.MediumCoverageLimit
		rp.DeductibleIncrease = matrix.MediumDeductible
	case "high":
		rp.PremiumAdjustment = matrix.HighPremiumAdjustment
		rp.CoverageLimit = matrix.HighCoverageLimit
		rp.DeductibleIncrease = matrix.HighDeductible
		rp.RequiresDeposit = true
		rp.DepositAmount = matrix.HighCoverageLimit * 0.1 // 10% deposit
	case "critical":
		rp.PremiumAdjustment = matrix.CriticalPremiumAdjustment
		rp.CoverageLimit = matrix.CriticalCoverageLimit
		rp.DeductibleIncrease = matrix.CriticalDeductible
		rp.RequiresDeposit = true
		rp.DepositAmount = matrix.CriticalCoverageLimit * 0.2 // 20% deposit
	}
}

// ApplyFactor applies a risk factor to the profile
func (rf *RiskFactor) ApplyFactor(profile *RiskProfile) {
	impact := rf.Score * rf.Weight

	switch rf.Category {
	case "fraud":
		profile.FraudRiskScore += impact
	case "credit":
		profile.CreditRiskScore += impact
	case "claim":
		profile.ClaimRiskScore += impact
	case "behavioral":
		profile.BehavioralRiskScore += impact
	case "device":
		profile.DeviceRiskScore += impact
	case "geographic":
		profile.GeographicRiskScore += impact
	}
}

// IsExpired checks if factor has expired
func (rf *RiskFactor) IsExpired() bool {
	if rf.ValidUntil == nil {
		return false
	}
	return time.Now().After(*rf.ValidUntil)
}

// ApplyEvent applies a risk event impact
func (re *RiskEvent) ApplyEvent(profile *RiskProfile) {
	if !re.IsActive {
		return
	}

	if re.ExpiresAt != nil && time.Now().After(*re.ExpiresAt) {
		re.IsActive = false
		return
	}

	// Apply score impact based on severity
	multiplier := 1.0
	switch re.Severity {
	case "critical":
		multiplier = 2.0
	case "high":
		multiplier = 1.5
	case "medium":
		multiplier = 1.0
	case "low":
		multiplier = 0.5
	}

	profile.OverallRiskScore += re.ScoreImpact * multiplier
	profile.OverallRiskScore = math.Min(100, math.Max(0, profile.OverallRiskScore))
}

// EvaluateRule evaluates if a risk rule applies
func (rr *RiskRule) EvaluateRule(value interface{}) bool {
	// This would implement complex rule evaluation logic
	// Simplified for demonstration
	return true
}

// GetPremiumAdjustment gets the premium adjustment for a risk level
func (rm *RiskMatrix) GetPremiumAdjustment(riskLevel string) float64 {
	switch riskLevel {
	case "low":
		return rm.LowPremiumAdjustment
	case "medium":
		return rm.MediumPremiumAdjustment
	case "high":
		return rm.HighPremiumAdjustment
	case "critical":
		return rm.CriticalPremiumAdjustment
	default:
		return 0
	}
}

// GetCoverageLimit gets the coverage limit for a risk level
func (rm *RiskMatrix) GetCoverageLimit(riskLevel string) float64 {
	switch riskLevel {
	case "low":
		return rm.LowCoverageLimit
	case "medium":
		return rm.MediumCoverageLimit
	case "high":
		return rm.HighCoverageLimit
	case "critical":
		return rm.CriticalCoverageLimit
	default:
		return rm.MediumCoverageLimit
	}
}
