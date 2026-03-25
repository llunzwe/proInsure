package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserUnderwriting manages underwriting data and risk assessment
type UserUnderwriting struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Underwriting Status
	UnderwritingStatus string     `gorm:"type:varchar(20)" json:"underwriting_status"`
	UnderwritingScore  float64    `gorm:"default:0" json:"underwriting_score"`
	RiskCategory       string     `gorm:"type:varchar(20)" json:"risk_category"`
	UnderwritingDate   *time.Time `json:"underwriting_date"`
	ExpiryDate         *time.Time `json:"expiry_date"`
	ReviewRequired     bool       `gorm:"default:false" json:"review_required"`

	// Risk Assessment
	RiskScore       float64                  `gorm:"default:0" json:"risk_score"`
	RiskFactors     []map[string]interface{} `gorm:"type:json" json:"risk_factors"`
	RiskMitigations []string                 `gorm:"type:json" json:"risk_mitigations"`
	AcceptableRisk  bool                     `gorm:"default:true" json:"acceptable_risk"`
	RiskTolerance   string                   `gorm:"type:varchar(20)" json:"risk_tolerance"`

	// Medical Underwriting
	MedicalStatus         string                   `gorm:"type:varchar(20)" json:"medical_status"`
	HealthScore           float64                  `gorm:"default:0" json:"health_score"`
	PreExistingConditions []string                 `gorm:"type:json" json:"pre_existing_conditions"`
	MedicalExamRequired   bool                     `gorm:"default:false" json:"medical_exam_required"`
	MedicalExamDate       *time.Time               `json:"medical_exam_date"`
	MedicalReports        []map[string]interface{} `gorm:"type:json" json:"medical_reports"`

	// Financial Underwriting
	FinancialRiskScore float64 `gorm:"default:0" json:"financial_risk_score"`
	IncomeVerified     bool    `gorm:"default:false" json:"income_verified"`
	AssetVerified      bool    `gorm:"default:false" json:"asset_verified"`
	DebtServiceRatio   float64 `gorm:"default:0" json:"debt_service_ratio"`
	FinancialCapacity  string  `gorm:"type:varchar(20)" json:"financial_capacity"`

	// Behavioral Underwriting
	BehavioralScore float64                  `gorm:"default:0" json:"behavioral_score"`
	ClaimHistory    []map[string]interface{} `gorm:"type:json" json:"claim_history"`
	PaymentHistory  string                   `gorm:"type:varchar(20)" json:"payment_history"`
	LoyaltyScore    float64                  `gorm:"default:0" json:"loyalty_score"`

	// Decision Details
	UnderwritingDecision string             `gorm:"type:varchar(20)" json:"underwriting_decision"`
	DecisionDate         *time.Time         `json:"decision_date"`
	DecisionReason       string             `gorm:"type:text" json:"decision_reason"`
	Conditions           []string           `gorm:"type:json" json:"conditions"`
	Exclusions           []string           `gorm:"type:json" json:"exclusions"`
	LoadingFactors       map[string]float64 `gorm:"type:json" json:"loading_factors"`

	// Manual Review
	ManualReviewRequired bool       `gorm:"default:false" json:"manual_review_required"`
	ReviewerID           *uuid.UUID `gorm:"type:uuid" json:"reviewer_id"`
	ReviewDate           *time.Time `json:"review_date"`
	ReviewNotes          string     `gorm:"type:text" json:"review_notes"`
	OverrideApplied      bool       `gorm:"default:false" json:"override_applied"`
	OverrideReason       string     `gorm:"type:text" json:"override_reason"`
}

// UserUnderwritingHistory tracks underwriting history and changes
type UserUnderwritingHistory struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	UnderwriterID *uuid.UUID `gorm:"type:uuid" json:"underwriter_id"`

	// History Entry
	EventType        string                 `gorm:"type:varchar(50)" json:"event_type"`
	EventDate        time.Time              `json:"event_date"`
	PreviousScore    float64                `gorm:"default:0" json:"previous_score"`
	NewScore         float64                `gorm:"default:0" json:"new_score"`
	PreviousCategory string                 `gorm:"type:varchar(20)" json:"previous_category"`
	NewCategory      string                 `gorm:"type:varchar(20)" json:"new_category"`
	ChangeReason     string                 `gorm:"type:text" json:"change_reason"`
	DataUsed         map[string]interface{} `gorm:"type:json" json:"data_used"`
	ModelVersion     string                 `gorm:"type:varchar(20)" json:"model_version"`
	Confidence       float64                `gorm:"default:0" json:"confidence"`
}

// UserRiskEvolution tracks risk score changes over time
type UserRiskEvolution struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Risk Timeline
	RiskHistory     []map[string]interface{} `gorm:"type:json" json:"risk_history"`
	ScoreTrend      string                   `gorm:"type:varchar(20)" json:"score_trend"` // improving/stable/deteriorating
	VolatilityScore float64                  `gorm:"default:0" json:"volatility_score"`
	StabilityPeriod int                      `json:"stability_period_days"`

	// Risk Components
	ComponentScores  map[string]float64 `gorm:"type:json" json:"component_scores"`
	ComponentWeights map[string]float64 `gorm:"type:json" json:"component_weights"`
	ComponentTrends  map[string]string  `gorm:"type:json" json:"component_trends"`

	// Predictions
	PredictedRiskScore   float64  `gorm:"default:0" json:"predicted_risk_score"`
	PredictionConfidence float64  `gorm:"default:0" json:"prediction_confidence"`
	PredictionHorizon    int      `json:"prediction_horizon_days"`
	RiskDrivers          []string `gorm:"type:json" json:"risk_drivers"`

	// Benchmarking
	PeerComparison    float64 `gorm:"default:0" json:"peer_comparison"`
	IndustryBenchmark float64 `gorm:"default:0" json:"industry_benchmark"`
	PercentileRank    float64 `gorm:"default:0" json:"percentile_rank"`

	// Alerts
	RiskAlerts        []map[string]interface{} `gorm:"type:json" json:"risk_alerts"`
	ThresholdBreaches []map[string]interface{} `gorm:"type:json" json:"threshold_breaches"`
	EarlyWarnings     []string                 `gorm:"type:json" json:"early_warnings"`
}

// UserReserves manages reserve requirements and allocations
type UserReserves struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Reserve Amounts
	TotalReserves       decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_reserves"`
	ClaimReserves       decimal.Decimal `gorm:"type:decimal(15,2)" json:"claim_reserves"`
	PremiumReserves     decimal.Decimal `gorm:"type:decimal(15,2)" json:"premium_reserves"`
	TechnicalReserves   decimal.Decimal `gorm:"type:decimal(15,2)" json:"technical_reserves"`
	ContingencyReserves decimal.Decimal `gorm:"type:decimal(15,2)" json:"contingency_reserves"`

	// Reserve Calculations
	RequiredReserves decimal.Decimal `gorm:"type:decimal(15,2)" json:"required_reserves"`
	ActualReserves   decimal.Decimal `gorm:"type:decimal(15,2)" json:"actual_reserves"`
	ReserveAdequacy  float64         `gorm:"default:0" json:"reserve_adequacy"`
	ReserveRatio     float64         `gorm:"default:0" json:"reserve_ratio"`

	// Collateral
	CollateralType    string          `gorm:"type:varchar(50)" json:"collateral_type"`
	CollateralValue   decimal.Decimal `gorm:"type:decimal(15,2)" json:"collateral_value"`
	CollateralStatus  string          `gorm:"type:varchar(20)" json:"collateral_status"`
	LastValuationDate *time.Time      `json:"last_valuation_date"`

	// Guarantees
	GuaranteeAmount  decimal.Decimal        `gorm:"type:decimal(15,2)" json:"guarantee_amount"`
	GuaranteeType    string                 `gorm:"type:varchar(50)" json:"guarantee_type"`
	GuarantorDetails map[string]interface{} `gorm:"type:json" json:"guarantor_details"`
	GuaranteeExpiry  *time.Time             `json:"guarantee_expiry"`

	// Adjustments
	ReserveAdjustments []map[string]interface{} `gorm:"type:json" json:"reserve_adjustments"`
	LastAdjustmentDate *time.Time               `json:"last_adjustment_date"`
	AdjustmentReason   string                   `gorm:"type:text" json:"adjustment_reason"`
}

// UserUnderwritingModels tracks ML models used for underwriting
type UserUnderwritingModels struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Model Information
	ActiveModels    []map[string]interface{} `gorm:"type:json" json:"active_models"`
	ModelVersions   map[string]string        `gorm:"type:json" json:"model_versions"`
	ModelScores     map[string]float64       `gorm:"type:json" json:"model_scores"`
	ModelConfidence map[string]float64       `gorm:"type:json" json:"model_confidence"`

	// Feature Importance
	FeatureImportance    map[string]float64     `gorm:"type:json" json:"feature_importance"`
	FeatureContributions map[string]interface{} `gorm:"type:json" json:"feature_contributions"`
	MissingFeatures      []string               `gorm:"type:json" json:"missing_features"`
	FeatureQuality       map[string]float64     `gorm:"type:json" json:"feature_quality"`

	// Model Performance
	ModelAccuracy  map[string]float64 `gorm:"type:json" json:"model_accuracy"`
	ModelPrecision map[string]float64 `gorm:"type:json" json:"model_precision"`
	ModelRecall    map[string]float64 `gorm:"type:json" json:"model_recall"`
	ModelF1Score   map[string]float64 `gorm:"type:json" json:"model_f1_score"`

	// Explanations
	DecisionExplanation map[string]interface{}   `gorm:"type:json" json:"decision_explanation"`
	RuleBasedDecisions  []map[string]interface{} `gorm:"type:json" json:"rule_based_decisions"`
	ModelInterpretation map[string]interface{}   `gorm:"type:json" json:"model_interpretation"`

	// A/B Testing
	TestGroup        string                 `gorm:"type:varchar(50)" json:"test_group"`
	ExperimentID     string                 `gorm:"type:varchar(100)" json:"experiment_id"`
	ModelComparisons map[string]interface{} `gorm:"type:json" json:"model_comparisons"`
}

// UserPricingModel manages individualized pricing
type UserPricingModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Base Pricing
	BasePremium         decimal.Decimal `gorm:"type:decimal(10,2)" json:"base_premium"`
	RiskAdjustedPremium decimal.Decimal `gorm:"type:decimal(10,2)" json:"risk_adjusted_premium"`
	FinalPremium        decimal.Decimal `gorm:"type:decimal(10,2)" json:"final_premium"`
	Currency            string          `gorm:"type:varchar(3)" json:"currency"`

	// Adjustments
	RiskLoading     decimal.Decimal `gorm:"type:decimal(5,2)" json:"risk_loading"`
	ExpenseLoading  decimal.Decimal `gorm:"type:decimal(5,2)" json:"expense_loading"`
	ProfitMargin    decimal.Decimal `gorm:"type:decimal(5,2)" json:"profit_margin"`
	DiscountApplied decimal.Decimal `gorm:"type:decimal(5,2)" json:"discount_applied"`
	TaxesAndFees    decimal.Decimal `gorm:"type:decimal(10,2)" json:"taxes_and_fees"`

	// Factors
	PricingFactors     map[string]float64 `gorm:"type:json" json:"pricing_factors"`
	CompetitiveFactors map[string]float64 `gorm:"type:json" json:"competitive_factors"`
	MarketFactors      map[string]float64 `gorm:"type:json" json:"market_factors"`
	PersonalFactors    map[string]float64 `gorm:"type:json" json:"personal_factors"`

	// Dynamic Pricing
	DynamicPricingEnabled bool            `gorm:"default:false" json:"dynamic_pricing_enabled"`
	PriceElasticity       float64         `gorm:"default:0" json:"price_elasticity"`
	WillingnessToPay      decimal.Decimal `gorm:"type:decimal(10,2)" json:"willingness_to_pay"`
	OptimalPrice          decimal.Decimal `gorm:"type:decimal(10,2)" json:"optimal_price"`

	// Pricing History
	PriceHistory      []map[string]interface{} `gorm:"type:json" json:"price_history"`
	LastPriceChange   *time.Time               `json:"last_price_change"`
	PriceChangeReason string                   `gorm:"type:text" json:"price_change_reason"`
	PriceValidUntil   *time.Time               `json:"price_valid_until"`

	// Competitiveness
	MarketComparison float64                    `gorm:"default:0" json:"market_comparison"`
	CompetitorPrices map[string]decimal.Decimal `gorm:"type:json" json:"competitor_prices"`
	PricePosition    string                     `gorm:"type:varchar(20)" json:"price_position"`
	WinRate          float64                    `gorm:"default:0" json:"win_rate"`
}

// TableName returns the table name
func (UserUnderwriting) TableName() string {
	return "user_underwriting"
}

// TableName returns the table name
func (UserUnderwritingHistory) TableName() string {
	return "user_underwriting_history"
}

// TableName returns the table name
func (UserRiskEvolution) TableName() string {
	return "user_risk_evolution"
}

// TableName returns the table name
func (UserReserves) TableName() string {
	return "user_reserves"
}

// TableName returns the table name
func (UserUnderwritingModels) TableName() string {
	return "user_underwriting_models"
}

// TableName returns the table name
func (UserPricingModel) TableName() string {
	return "user_pricing_models"
}
