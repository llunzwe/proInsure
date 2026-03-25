package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceCustomerLifetimeValue represents comprehensive CLV calculations
type DeviceCustomerLifetimeValue struct {
	database.BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// CLV Metrics
	CLV          float64 `gorm:"not null" json:"clv"` // total lifetime value
	CurrentCLV   float64 `json:"current_clv"`         // current calculated CLV
	PotentialCLV float64 `json:"potential_clv"`       // potential future CLV
	RealizedCLV  float64 `json:"realized_clv"`        // already realized CLV

	LastCalculated time.Time `gorm:"not null" json:"last_calculated"`

	// CLV Components
	PremiumRevenue   float64 `json:"premium_revenue"`   // insurance premiums
	ClaimCosts       float64 `json:"claim_costs"`       // total claims paid
	AncillaryRevenue float64 `json:"ancillary_revenue"` // accessories, services, etc.
	ReferralRevenue  float64 `json:"referral_revenue"`  // revenue from referrals
	UpgradeRevenue   float64 `json:"upgrade_revenue"`   // device upgrades/replacements

	// Customer Metrics
	CustomerTenure    int     `json:"customer_tenure"`    // days as customer
	PurchaseFrequency float64 `json:"purchase_frequency"` // purchases per month
	AverageOrderValue float64 `json:"average_order_value"`
	RetentionScore    float64 `json:"retention_score"` // 0-100

	// Risk & Churn Analysis
	ChurnProbability float64 `json:"churn_probability"`                   // 0-100
	ChurnRiskFactors string  `gorm:"type:json" json:"churn_risk_factors"` // JSON risk factors
	TimeToChurn      int     `json:"time_to_churn,omitempty"`             // days until likely churn

	// Segmentation
	CustomerSegment string `json:"customer_segment"` // high_value, medium_value, low_value, at_risk
	LoyaltyTier     string `json:"loyalty_tier"`     // platinum, gold, silver, bronze

	// Predictive Metrics
	PredictedRevenue   float64 `json:"predicted_revenue"` // next 12 months
	UpsellPotential    float64 `json:"upsell_potential"`  // additional revenue potential
	CrossSellPotential float64 `json:"cross_sell_potential"`

	// Model Information
	CLVModelVersion   string  `json:"clv_model_version"`
	CalculationMethod string  `json:"calculation_method"` // historical, predictive, hybrid
	ConfidenceScore   float64 `json:"confidence_score"`   // 0-100

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceChurnPrediction represents churn risk analysis
type DeviceChurnPrediction struct {
	database.BaseModel

	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	PredictionID string    `gorm:"uniqueIndex;not null" json:"prediction_id"`

	// Churn Assessment
	ChurnProbability   float64    `gorm:"not null" json:"churn_probability"`        // 0-100
	ChurnRiskLevel     string     `gorm:"type:varchar(20)" json:"churn_risk_level"` // low, medium, high, critical
	PredictedChurnDate *time.Time `json:"predicted_churn_date,omitempty"`

	PredictionDate time.Time `gorm:"not null" json:"prediction_date"`
	ValidUntil     time.Time `json:"valid_until"`

	// Risk Factors
	RiskFactors string `gorm:"type:json" json:"risk_factors"` // JSON risk indicators
	RiskWeights string `gorm:"type:json" json:"risk_weights"` // JSON factor weightings

	// Behavioral Indicators
	UsageDecline     bool `json:"usage_decline"`
	PaymentDelays    int  `json:"payment_delays"`
	SupportTickets   int  `json:"support_tickets"`
	NegativeFeedback bool `json:"negative_feedback"`

	// Retention Actions
	RecommendedActions     string  `gorm:"type:json" json:"recommended_actions"` // JSON retention strategies
	ActionPriority         string  `gorm:"type:varchar(20)" json:"action_priority"`
	EstimatedEffectiveness float64 `json:"estimated_effectiveness"` // 0-100

	// Intervention Tracking
	InterventionTaken bool       `json:"intervention_taken"`
	InterventionType  string     `json:"intervention_type,omitempty"`
	InterventionDate  *time.Time `json:"intervention_date,omitempty"`
	InterventionCost  float64    `json:"intervention_cost,omitempty"`

	// Outcome Tracking
	ChurnOccurred    bool       `json:"churn_occurred,omitempty"`
	ActualChurnDate  *time.Time `json:"actual_churn_date,omitempty"`
	RetentionSuccess bool       `json:"retention_success,omitempty"`

	// Model Information
	AIModelVersion     string   `json:"ai_model_version"`
	PredictionAccuracy *float64 `json:"prediction_accuracy,omitempty"` // retrospective accuracy

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DevicePredictiveAnalytics represents comprehensive predictive analytics
type DevicePredictiveAnalytics struct {
	database.BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Claims Prediction
	ClaimsPrediction   string  `gorm:"type:json" json:"claims_prediction"` // JSON claims forecast
	ClaimProbability   float64 `json:"claim_probability"`                  // next 12 months
	ExpectedClaimValue float64 `json:"expected_claim_value"`

	// Revenue Forecasting
	RevenueForecast string `gorm:"type:json" json:"revenue_forecast"` // JSON revenue projections
	PremiumForecast string `gorm:"type:json" json:"premium_forecast"` // JSON premium forecasts
	ServiceRevenue  string `gorm:"type:json" json:"service_revenue"`  // JSON service revenue

	// Risk Trend Analysis
	RiskTrends       string  `gorm:"type:json" json:"risk_trends"` // JSON risk evolution
	RiskVelocity     float64 `json:"risk_velocity"`                // rate of risk change
	RiskAcceleration float64 `json:"risk_acceleration"`            // acceleration of risk change

	// Behavioral Forecasting
	BehaviorForecast string `gorm:"type:json" json:"behavior_forecast"` // JSON behavior predictions
	UsagePatterns    string `gorm:"type:json" json:"usage_patterns"`    // JSON usage trend predictions
	LoyaltyForecast  string `gorm:"type:json" json:"loyalty_forecast"`  // JSON loyalty predictions

	// Market Intelligence
	MarketPosition   string  `gorm:"type:json" json:"market_position"` // JSON competitive positioning
	PriceSensitivity float64 `json:"price_sensitivity"`
	BrandAffinity    float64 `json:"brand_affinity"` // 0-100

	// Operational Forecasting
	SupportTickets   string `gorm:"type:json" json:"support_tickets"`   // JSON support demand forecast
	MaintenanceNeeds string `gorm:"type:json" json:"maintenance_needs"` // JSON maintenance predictions

	// Model Performance
	OverallAccuracy     float64 `json:"overall_accuracy"`                      // 0-100
	ModelVersions       string  `gorm:"type:json" json:"model_versions"`       // JSON model version tracking
	ConfidenceIntervals string  `gorm:"type:json" json:"confidence_intervals"` // JSON prediction confidence

	LastUpdated     time.Time `gorm:"not null" json:"last_updated"`
	ForecastHorizon int       `json:"forecast_horizon"` // months

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceProfitabilityMetrics represents detailed profitability analysis
type DeviceProfitabilityMetrics struct {
	database.BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	MetricID string    `gorm:"uniqueIndex;not null" json:"metric_id"`

	// Revenue Streams
	PremiumRevenue   float64 `json:"premium_revenue"`
	ServiceRevenue   float64 `json:"service_revenue"`
	AncillaryRevenue float64 `json:"ancillary_revenue"`
	TotalRevenue     float64 `json:"total_revenue"`

	// Cost Components
	AcquisitionCost    float64 `json:"acquisition_cost"`
	UnderwritingCost   float64 `json:"underwriting_cost"`
	ClaimsCost         float64 `json:"claims_cost"`
	AdministrativeCost float64 `json:"administrative_cost"`
	TotalCost          float64 `json:"total_cost"`

	// Profitability Metrics
	GrossMargin   float64 `json:"gross_margin"`   // percentage
	NetMargin     float64 `json:"net_margin"`     // percentage
	ROI           float64 `json:"roi"`            // percentage
	PaybackPeriod int     `json:"payback_period"` // months

	// Risk-Adjusted Metrics
	RiskAdjustedMargin float64 `json:"risk_adjusted_margin"`
	SharpeRatio        float64 `json:"sharpe_ratio"`  // risk-adjusted return
	ValueAtRisk        float64 `json:"value_at_risk"` // 95% VaR

	// Lifetime Metrics
	LifetimeValue  float64 `json:"lifetime_value"`
	LifetimeCost   float64 `json:"lifetime_cost"`
	LifetimeProfit float64 `json:"lifetime_profit"`

	// Segmentation
	ProfitabilitySegment string `gorm:"type:varchar(20)" json:"profitability_segment"` // high, medium, low, loss
	PerformanceQuartile  int    `json:"performance_quartile"`                          // 1-4, 1 being best

	// Trend Analysis
	RevenueTrend string `gorm:"type:json" json:"revenue_trend"` // JSON trend data
	CostTrend    string `gorm:"type:json" json:"cost_trend"`    // JSON trend data
	MarginTrend  string `gorm:"type:json" json:"margin_trend"`  // JSON trend data

	// Benchmarking
	IndustryAverageMargin float64 `json:"industry_average_margin"`
	PeerGroupRanking      int     `json:"peer_group_ranking"`
	PerformanceIndex      float64 `json:"performance_index"` // relative to peers

	CalculationDate time.Time `gorm:"not null" json:"calculation_date"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceMarketIntelligence represents market and competitive analysis
type DeviceMarketIntelligence struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	IntelligenceID string    `gorm:"uniqueIndex;not null" json:"intelligence_id"`

	// Market Position
	MarketShare         float64 `json:"market_share"`                                 // percentage
	BrandStrength       float64 `json:"brand_strength"`                               // 0-100
	CompetitivePosition string  `gorm:"type:varchar(20)" json:"competitive_position"` // leader, challenger, follower, niche

	// Pricing Intelligence
	PricePosition    string  `gorm:"type:varchar(20)" json:"price_position"` // premium, mid-range, budget
	PriceElasticity  float64 `json:"price_elasticity"`
	ValueProposition string  `json:"value_proposition"`

	// Competitive Analysis
	CompetitorAnalysis string  `gorm:"type:json" json:"competitor_analysis"` // JSON competitor data
	ThreatLevel        string  `gorm:"type:varchar(20)" json:"threat_level"` // low, medium, high
	OpportunityScore   float64 `json:"opportunity_score"`                    // 0-100

	// Customer Insights
	CustomerSatisfaction float64 `json:"customer_satisfaction"` // 0-100
	NetPromoterScore     float64 `json:"net_promoter_score"`    // -100 to 100
	CustomerRetention    float64 `json:"customer_retention"`    // percentage

	// Market Trends
	MarketTrends     string `gorm:"type:json" json:"market_trends"`     // JSON trend analysis
	TechnologyTrends string `gorm:"type:json" json:"technology_trends"` // JSON tech developments
	RegulatoryTrends string `gorm:"type:json" json:"regulatory_trends"` // JSON regulatory changes

	// Economic Indicators
	EconomicImpact       string  `gorm:"type:json" json:"economic_impact"` // JSON economic factors
	InflationSensitivity float64 `json:"inflation_sensitivity"`
	CurrencyRisk         float64 `json:"currency_risk"`

	// SWOT Analysis
	Strengths     string `gorm:"type:json" json:"strengths"`     // JSON strengths
	Weaknesses    string `gorm:"type:json" json:"weaknesses"`    // JSON weaknesses
	Opportunities string `gorm:"type:json" json:"opportunities"` // JSON opportunities
	Threats       string `gorm:"type:json" json:"threats"`       // JSON threats

	// Strategic Recommendations
	StrategicActions string `gorm:"type:json" json:"strategic_actions"`     // JSON recommended actions
	PriorityLevel    string `gorm:"type:varchar(20)" json:"priority_level"` // high, medium, low

	IntelligenceDate time.Time `gorm:"not null" json:"intelligence_date"`
	DataFreshness    int       `json:"data_freshness"` // days since last update

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceRiskTrendAnalysis represents comprehensive risk evolution tracking
type DeviceRiskTrendAnalysis struct {
	database.BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Risk Evolution
	RiskTrajectory   string  `gorm:"type:json" json:"risk_trajectory"` // JSON risk over time
	RiskVelocity     float64 `json:"risk_velocity"`                    // rate of risk change
	RiskAcceleration float64 `json:"risk_acceleration"`                // acceleration of change

	// Trend Analysis
	TrendDirection string  `gorm:"type:varchar(20)" json:"trend_direction"` // improving, stable, deteriorating
	TrendStrength  float64 `json:"trend_strength"`                          // 0-100, strength of trend
	TrendDuration  int     `json:"trend_duration"`                          // days trend has been active

	// Risk Components
	TheftRiskTrend       string `gorm:"type:json" json:"theft_risk_trend"`
	DamageRiskTrend      string `gorm:"type:json" json:"damage_risk_trend"`
	OperationalRiskTrend string `gorm:"type:json" json:"operational_risk_trend"`
	ComplianceRiskTrend  string `gorm:"type:json" json:"compliance_risk_trend"`

	// Predictive Indicators
	LeadingIndicators string `gorm:"type:json" json:"leading_indicators"` // JSON predictive signals
	LaggingIndicators string `gorm:"type:json" json:"lagging_indicators"` // JSON outcome signals

	// Risk Drivers
	PrimaryRiskDrivers string `gorm:"type:json" json:"primary_risk_drivers"` // JSON main risk factors
	EmergingRisks      string `gorm:"type:json" json:"emerging_risks"`       // JSON new risk factors
	MitigatedRisks     string `gorm:"type:json" json:"mitigated_risks"`      // JSON controlled risks

	// Intervention Impact
	RiskInterventions         string `gorm:"type:json" json:"risk_interventions"`         // JSON risk mitigation actions
	InterventionEffectiveness string `gorm:"type:json" json:"intervention_effectiveness"` // JSON effectiveness metrics

	// Forecasting
	RiskForecast     string   `gorm:"type:json" json:"risk_forecast"` // JSON future risk projections
	ForecastHorizon  int      `json:"forecast_horizon"`               // months
	ForecastAccuracy *float64 `json:"forecast_accuracy,omitempty"`    // retrospective accuracy

	LastAnalyzed   time.Time `gorm:"not null" json:"last_analyzed"`
	AnalysisPeriod int       `json:"analysis_period"` // days of data analyzed

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// TableName returns the table name for DeviceCustomerLifetimeValue
func (DeviceCustomerLifetimeValue) TableName() string {
	return "device_customer_lifetime_value"
}

// TableName returns the table name for DeviceChurnPrediction
func (DeviceChurnPrediction) TableName() string {
	return "device_churn_predictions"
}

// TableName returns the table name for DevicePredictiveAnalytics
func (DevicePredictiveAnalytics) TableName() string {
	return "device_predictive_analytics"
}

// TableName returns the table name for DeviceProfitabilityMetrics
func (DeviceProfitabilityMetrics) TableName() string {
	return "device_profitability_metrics"
}

// TableName returns the table name for DeviceMarketIntelligence
func (DeviceMarketIntelligence) TableName() string {
	return "device_market_intelligence"
}

// TableName returns the table name for DeviceRiskTrendAnalysis
func (DeviceRiskTrendAnalysis) TableName() string {
	return "device_risk_trend_analysis"
}
