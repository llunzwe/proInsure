package device

import (
	"math"
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceDepreciationCurve tracks depreciation patterns and market value over time
type DeviceDepreciationCurve struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ModelID  string    `gorm:"type:varchar(100);index" json:"model_id"`
	Brand    string    `gorm:"type:varchar(50)" json:"brand"`

	// Depreciation Patterns
	DepreciationType       string  `gorm:"type:varchar(50)" json:"depreciation_type"` // linear, declining_balance, sum_of_years
	AnnualDepreciationRate float64 `json:"annual_depreciation_rate"`
	FirstYearDepreciation  float64 `json:"first_year_depreciation"`
	DepreciationCurve      string  `gorm:"type:json" json:"depreciation_curve"` // JSON array of monthly values

	// Market Value Tracking
	OriginalValue          float64   `json:"original_value"`
	CurrentMarketValue     float64   `json:"current_market_value"`
	PredictedValue6Months  float64   `json:"predicted_value_6_months"`
	PredictedValue12Months float64   `json:"predicted_value_12_months"`
	LastValuationDate      time.Time `json:"last_valuation_date"`

	// Value Factors
	SeasonalAdjustment float64 `json:"seasonal_adjustment"`  // Multiplier for seasonal fluctuations
	BrandRetentionRate float64 `json:"brand_retention_rate"` // 0-1 scale
	ObsolescenceFactor float64 `json:"obsolescence_factor"`  // Impact of new models
	ConditionImpact    float64 `json:"condition_impact"`     // Condition-based adjustment
	RegionalVariation  float64 `json:"regional_variation"`   // Regional price differences
	SupplyDemandImpact float64 `json:"supply_demand_impact"` // Market supply/demand factor

	// Historical Data
	ValueHistory     string `gorm:"type:json" json:"value_history"`      // JSON array of historical values
	MarketEvents     string `gorm:"type:json" json:"market_events"`      // JSON array of events affecting value
	NewModelReleases string `gorm:"type:json" json:"new_model_releases"` // JSON array of release impacts

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceInsurableValue manages insurance valuation and coverage limits
type DeviceInsurableValue struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	PolicyID uuid.UUID `gorm:"type:uuid;index" json:"policy_id"`

	// Valuation Methods
	ValuationMethod     string  `gorm:"type:varchar(50)" json:"valuation_method"` // agreed_value, actual_cash, replacement_cost
	AgreedValue         float64 `json:"agreed_value"`
	ActualCashValue     float64 `json:"actual_cash_value"`
	ReplacementCost     float64 `json:"replacement_cost"`
	CurrentInsuredValue float64 `json:"current_insured_value"`

	// Depreciation Schedule
	DepreciationSchedule string    `gorm:"type:json" json:"depreciation_schedule"` // JSON array
	DepreciatedValue     float64   `json:"depreciated_value"`
	DepreciationMethod   string    `json:"depreciation_method"`
	EffectiveDate        time.Time `json:"effective_date"`

	// Coverage Limits
	MaxCoverageLimit    float64 `json:"max_coverage_limit"`
	MinCoverageLimit    float64 `json:"min_coverage_limit"`
	RecommendedCoverage float64 `json:"recommended_coverage"`
	CurrentCoverage     float64 `json:"current_coverage"`

	// Insurance Analysis
	UnderInsured          bool    `json:"under_insured"`
	UnderInsuranceAmount  float64 `json:"under_insurance_amount"`
	OverInsured           bool    `json:"over_insured"`
	OverInsuranceAmount   float64 `json:"over_insurance_amount"`
	InsuranceToValueRatio float64 `json:"insurance_to_value_ratio"`

	// Appraisal & Verification
	AppraisalRequired  bool       `json:"appraisal_required"`
	LastAppraisalDate  *time.Time `json:"last_appraisal_date"`
	AppraisalValue     float64    `json:"appraisal_value"`
	AppraisalSource    string     `json:"appraisal_source"`
	VerificationStatus string     `json:"verification_status"`
	VerificationDate   *time.Time `json:"verification_date"`

	// Dispute Resolution
	ValueDisputed     bool    `json:"value_disputed"`
	DisputeReason     string  `json:"dispute_reason"`
	DisputeResolution string  `json:"dispute_resolution"`
	FinalAgreedValue  float64 `json:"final_agreed_value"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// Policy should be loaded via service layer using PolicyID to avoid circular import
}

// DeviceTotalCostOwnership tracks all costs associated with device ownership
type DeviceTotalCostOwnership struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	CalculationDate time.Time `json:"calculation_date"`
	OwnershipPeriod int       `json:"ownership_period"` // in months

	// Direct Costs
	PurchasePrice     float64 `json:"purchase_price"`
	InsurancePremiums float64 `json:"insurance_premiums"`
	RepairCosts       float64 `json:"repair_costs"`
	MaintenanceCosts  float64 `json:"maintenance_costs"`
	AccessoryCosts    float64 `json:"accessory_costs"`
	ServicePlanCosts  float64 `json:"service_plan_costs"`
	WarrantyCosts     float64 `json:"warranty_costs"`

	// Indirect Costs
	FinancingCosts  float64 `json:"financing_costs"`  // Interest if financed
	OpportunityCost float64 `json:"opportunity_cost"` // Lost investment returns
	DowntimeCosts   float64 `json:"downtime_costs"`   // Productivity loss
	DataPlanCosts   float64 `json:"data_plan_costs"`

	// Depreciation
	TotalDepreciation  float64 `json:"total_depreciation"`
	AnnualDepreciation float64 `json:"annual_depreciation"`
	ResidualValue      float64 `json:"residual_value"`

	// Disposal
	DisposalCost   float64 `json:"disposal_cost"`
	RecyclingFee   float64 `json:"recycling_fee"`
	DataWipingCost float64 `json:"data_wiping_cost"`

	// Total Cost Analysis
	TotalDirectCosts     float64 `json:"total_direct_costs"`
	TotalIndirectCosts   float64 `json:"total_indirect_costs"`
	TotalCostOfOwnership float64 `json:"total_cost_of_ownership"`
	MonthlyCost          float64 `json:"monthly_cost"`
	DailyCost            float64 `json:"daily_cost"`

	// Optimization
	CostSavingOpportunities string `gorm:"type:json" json:"cost_saving_opportunities"` // JSON array
	OptimalReplacementTime  int    `json:"optimal_replacement_time"`                   // in months
	BreakEvenPoint          int    `json:"break_even_point"`                           // months until TCO exceeds value

	// Comparisons
	IndustryAverageTCO float64 `json:"industry_average_tco"`
	TCOPercentile      float64 `json:"tco_percentile"` // Where this TCO falls in market

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceResaleValue tracks current and projected resale values
type DeviceResaleValue struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ValuationDate time.Time `json:"valuation_date"`

	// Current Values
	CurrentMarketValue float64 `json:"current_market_value"`
	TradeInValue       float64 `json:"trade_in_value"`
	PrivateSaleValue   float64 `json:"private_sale_value"`
	WholesaleValue     float64 `json:"wholesale_value"`
	SalvageValue       float64 `json:"salvage_value"`
	InstantCashValue   float64 `json:"instant_cash_value"` // Quick sale value

	// Regional Variations
	LocalMarketValue     float64 `json:"local_market_value"`
	NationalAverageValue float64 `json:"national_average_value"`
	RegionalPremium      float64 `json:"regional_premium"`    // +/- adjustment
	MarketDemandScore    float64 `json:"market_demand_score"` // 0-100

	// Condition Impact
	PerfectConditionValue float64 `json:"perfect_condition_value"`
	CurrentConditionValue float64 `json:"current_condition_value"`
	ConditionDiscount     float64 `json:"condition_discount"`
	RefurbishmentCost     float64 `json:"refurbishment_cost"`
	PostRefurbValue       float64 `json:"post_refurb_value"`

	// Timing Optimization
	BestSellTime         *time.Time `json:"best_sell_time"`
	SeasonalPeakValue    float64    `json:"seasonal_peak_value"`
	ProjectedValue30Days float64    `json:"projected_value_30_days"`
	ProjectedValue90Days float64    `json:"projected_value_90_days"`
	ValueDeclineRate     float64    `json:"value_decline_rate"` // Monthly %

	// Value Enhancement
	ValueEnhancements string  `gorm:"type:json" json:"value_enhancements"` // JSON array of tips
	MaxPotentialValue float64 `json:"max_potential_value"`
	EnhancementCost   float64 `json:"enhancement_cost"`
	EnhancementROI    float64 `json:"enhancement_roi"`

	// Market Indicators
	ListingsCount     int     `json:"listings_count"`
	AverageDaysToSell int     `json:"average_days_to_sell"`
	PriceVolatility   float64 `json:"price_volatility"`
	MarketTrend       string  `json:"market_trend"` // rising, stable, falling

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceFinancialRisk assesses financial risk and loss potential
type DeviceFinancialRisk struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	AssessmentDate time.Time `json:"assessment_date"`

	// Loss Potential
	MaximumLossAmount  float64 `json:"maximum_loss_amount"`
	ExpectedLossAmount float64 `json:"expected_loss_amount"`
	ProbableLossAmount float64 `json:"probable_loss_amount"`
	ValueAtRisk        float64 `json:"value_at_risk"`   // VaR calculation
	ConditionalVaR     float64 `json:"conditional_var"` // CVaR/Expected Shortfall

	// Risk Metrics
	LossFrequency        float64 `json:"loss_frequency"` // Expected events per year
	LossSeverity         float64 `json:"loss_severity"`  // Average loss amount
	CatastrophicLossProb float64 `json:"catastrophic_loss_prob"`
	TotalLossProbability float64 `json:"total_loss_probability"`

	// Premium Calculations
	BasePremium          float64 `json:"base_premium"`
	RiskAdjustedPremium  float64 `json:"risk_adjusted_premium"`
	RiskLoadingFactor    float64 `json:"risk_loading_factor"`
	ExpenseLoadingFactor float64 `json:"expense_loading_factor"`
	ProfitMargin         float64 `json:"profit_margin"`
	FinalPremium         float64 `json:"final_premium"`

	// Portfolio Impact
	PortfolioContribution  float64 `json:"portfolio_contribution"` // Risk contribution to portfolio
	DiversificationBenefit float64 `json:"diversification_benefit"`
	ConcentrationRisk      float64 `json:"concentration_risk"`
	CorrelationFactor      float64 `json:"correlation_factor"`

	// Risk Mitigation
	MitigationMeasures  string  `gorm:"type:json" json:"mitigation_measures"` // JSON array
	MitigationValue     float64 `json:"mitigation_value"`                     // Risk reduction value
	ResidualRisk        float64 `json:"residual_risk"`
	RiskRetentionAmount float64 `json:"risk_retention_amount"`

	// Insurance Analysis
	InsuranceGap           float64 `json:"insurance_gap"`
	OptimalCoverage        float64 `json:"optimal_coverage"`
	RiskTransferEfficiency float64 `json:"risk_transfer_efficiency"`
	SelfInsuranceViability bool    `json:"self_insurance_viability"`

	// Financial Scoring
	FinancialImpactScore float64 `json:"financial_impact_score"` // 0-100
	RiskRating           string  `json:"risk_rating"`            // AAA to D
	Creditworthiness     float64 `json:"creditworthiness"`

	// Stress Testing
	StressTestScenarios string  `gorm:"type:json" json:"stress_test_scenarios"` // JSON array
	WorstCaseScenario   float64 `json:"worst_case_scenario"`
	BestCaseScenario    float64 `json:"best_case_scenario"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceDepreciationCurve
func (ddc *DeviceDepreciationCurve) CalculateCurrentValue(deviceAge int) float64 {
	if ddc.DepreciationType == "linear" {
		yearlyDepreciation := ddc.OriginalValue * ddc.AnnualDepreciationRate
		currentValue := ddc.OriginalValue - (yearlyDepreciation * float64(deviceAge) / 12)
		if currentValue < 0 {
			currentValue = ddc.OriginalValue * 0.1 // Minimum 10% residual
		}
		return currentValue
	}

	if ddc.DepreciationType == "declining_balance" {
		years := float64(deviceAge) / 12
		currentValue := ddc.OriginalValue * math.Pow(1-ddc.AnnualDepreciationRate, years)
		return currentValue
	}

	return ddc.CurrentMarketValue
}

func (ddc *DeviceDepreciationCurve) ApplySeasonalAdjustment(baseValue float64) float64 {
	return baseValue * ddc.SeasonalAdjustment
}

func (ddc *DeviceDepreciationCurve) PredictFutureValue(monthsAhead int) float64 {
	monthlyRate := ddc.AnnualDepreciationRate / 12
	futureValue := ddc.CurrentMarketValue * math.Pow(1-monthlyRate, float64(monthsAhead))

	// Apply market factors
	futureValue *= ddc.BrandRetentionRate
	futureValue *= (1 - ddc.ObsolescenceFactor)

	return futureValue
}

// Methods for DeviceInsurableValue
func (div *DeviceInsurableValue) CalculateInsuranceGap() float64 {
	gap := div.ActualCashValue - div.CurrentCoverage
	if gap > 0 {
		div.UnderInsured = true
		div.UnderInsuranceAmount = gap
	} else if gap < 0 {
		div.OverInsured = true
		div.OverInsuranceAmount = -gap
	}
	return gap
}

func (div *DeviceInsurableValue) GetRecommendedCoverage() float64 {
	// 80% rule for insurance
	recommended := div.ReplacementCost * 0.8

	if recommended > div.MaxCoverageLimit {
		recommended = div.MaxCoverageLimit
	}
	if recommended < div.MinCoverageLimit {
		recommended = div.MinCoverageLimit
	}

	div.RecommendedCoverage = recommended
	return recommended
}

func (div *DeviceInsurableValue) RequiresAppraisal() bool {
	// Require appraisal if value disputed or high value
	return div.ValueDisputed ||
		div.CurrentInsuredValue > 5000 ||
		(div.LastAppraisalDate != nil && time.Since(*div.LastAppraisalDate) > 365*24*time.Hour)
}

// Methods for DeviceTotalCostOwnership
func (tco *DeviceTotalCostOwnership) CalculateTotalCost() float64 {
	// Direct costs
	tco.TotalDirectCosts = tco.PurchasePrice + tco.InsurancePremiums +
		tco.RepairCosts + tco.MaintenanceCosts + tco.AccessoryCosts +
		tco.ServicePlanCosts + tco.WarrantyCosts

	// Indirect costs
	tco.TotalIndirectCosts = tco.FinancingCosts + tco.OpportunityCost +
		tco.DowntimeCosts + tco.DataPlanCosts + tco.TotalDepreciation

	// Disposal costs
	disposalCosts := tco.DisposalCost + tco.RecyclingFee + tco.DataWipingCost

	tco.TotalCostOfOwnership = tco.TotalDirectCosts + tco.TotalIndirectCosts + disposalCosts

	if tco.OwnershipPeriod > 0 {
		tco.MonthlyCost = tco.TotalCostOfOwnership / float64(tco.OwnershipPeriod)
		tco.DailyCost = tco.MonthlyCost / 30
	}

	return tco.TotalCostOfOwnership
}

func (tco *DeviceTotalCostOwnership) CalculateOptimalReplacement() int {
	// Calculate when TCO exceeds device value
	if tco.ResidualValue > 0 && tco.MonthlyCost > 0 {
		monthsUntilBreakeven := int(tco.ResidualValue / tco.MonthlyCost)
		tco.BreakEvenPoint = monthsUntilBreakeven

		// Optimal replacement is typically 20% before breakeven
		tco.OptimalReplacementTime = int(float64(monthsUntilBreakeven) * 0.8)
	}

	return tco.OptimalReplacementTime
}

func (tco *DeviceTotalCostOwnership) IsAboveAverage() bool {
	return tco.TotalCostOfOwnership > tco.IndustryAverageTCO
}

// Methods for DeviceResaleValue
func (drv *DeviceResaleValue) CalculateBestSellTime() *time.Time {
	// Consider seasonal peaks and value decline rate
	now := time.Now()

	// If value declining rapidly, sell sooner
	if drv.ValueDeclineRate > 5 { // >5% monthly decline
		bestTime := now.Add(30 * 24 * time.Hour)
		drv.BestSellTime = &bestTime
	} else {
		// Find next seasonal peak (typically November-December for electronics)
		month := now.Month()
		year := now.Year()
		if month > 11 {
			year++
		}
		bestTime := time.Date(year, 11, 15, 0, 0, 0, 0, time.UTC)
		drv.BestSellTime = &bestTime
	}

	return drv.BestSellTime
}

func (drv *DeviceResaleValue) CalculateROI() float64 {
	if drv.EnhancementCost > 0 {
		valueIncrease := drv.PostRefurbValue - drv.CurrentConditionValue
		drv.EnhancementROI = (valueIncrease - drv.EnhancementCost) / drv.EnhancementCost * 100
	}
	return drv.EnhancementROI
}

func (drv *DeviceResaleValue) GetQuickSaleValue() float64 {
	// Instant cash value is typically 70-80% of market value
	drv.InstantCashValue = drv.CurrentMarketValue * 0.75
	return drv.InstantCashValue
}

// Methods for DeviceFinancialRisk
func (dfr *DeviceFinancialRisk) CalculateValueAtRisk(confidenceLevel float64) float64 {
	// Simple VaR calculation
	zScore := 1.645 // 95% confidence level
	if confidenceLevel == 0.99 {
		zScore = 2.326
	}

	dfr.ValueAtRisk = dfr.ExpectedLossAmount + (zScore * dfr.LossSeverity)

	// Calculate Conditional VaR (expected loss beyond VaR)
	dfr.ConditionalVaR = dfr.ValueAtRisk * 1.2

	return dfr.ValueAtRisk
}

func (dfr *DeviceFinancialRisk) CalculateRiskAdjustedPremium() float64 {
	// Base premium calculation
	purePremium := dfr.LossFrequency * dfr.LossSeverity

	// Add risk loading
	riskLoading := purePremium * dfr.RiskLoadingFactor

	// Add expense loading
	expenseLoading := purePremium * dfr.ExpenseLoadingFactor

	// Add profit margin
	profit := purePremium * dfr.ProfitMargin

	dfr.RiskAdjustedPremium = purePremium + riskLoading + expenseLoading + profit
	dfr.FinalPremium = dfr.RiskAdjustedPremium

	return dfr.FinalPremium
}

func (dfr *DeviceFinancialRisk) AssignRiskRating() string {
	ratings := []string{"AAA", "AA", "A", "BBB", "BB", "B", "CCC", "CC", "C", "D"}

	// Map financial impact score to rating
	index := int((100 - dfr.FinancialImpactScore) / 10)
	if index < 0 {
		index = 0
	}
	if index >= len(ratings) {
		index = len(ratings) - 1
	}

	dfr.RiskRating = ratings[index]
	return dfr.RiskRating
}
