package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MarketDataRepository defines the interface for market data operations
type MarketDataRepository interface {
	// === Market Price Data ===

	// GetCurrentMarketPrice gets current market price for device model
	GetCurrentMarketPrice(ctx context.Context, brand, model string) (decimal.Decimal, error)

	// GetMarketPriceHistory gets price history for a device model
	GetMarketPriceHistory(ctx context.Context, brand, model string, startDate, endDate time.Time) ([]*MarketPrice, error)

	// GetMarketTrend gets market trend for device category
	GetMarketTrend(ctx context.Context, category string, period int) (*MarketTrend, error)

	// UpdateMarketPrice updates market price for a device model
	UpdateMarketPrice(ctx context.Context, price *MarketPrice) error

	// BulkUpdateMarketPrices updates multiple market prices
	BulkUpdateMarketPrices(ctx context.Context, prices []*MarketPrice) error

	// === Depreciation Data ===

	// GetDepreciationCurve gets depreciation curve for device model
	GetDepreciationCurve(ctx context.Context, brand, model string) (*DepreciationCurve, error)

	// GetCategoryDepreciationRate gets average depreciation rate for category
	GetCategoryDepreciationRate(ctx context.Context, category string) (float64, error)

	// CalculateDepreciatedValue calculates depreciated value
	CalculateDepreciatedValue(ctx context.Context, originalValue decimal.Decimal, brand, model string, age time.Duration) (decimal.Decimal, error)

	// === Trade-In Values ===

	// GetTradeInMatrix gets trade-in value matrix by condition
	GetTradeInMatrix(ctx context.Context, brand, model string) (*TradeInMatrix, error)

	// GetCompetitorTradeInValues gets competitor trade-in values for comparison
	GetCompetitorTradeInValues(ctx context.Context, brand, model string) ([]*CompetitorTradeIn, error)

	// === Market Analytics ===

	// GetMarketShare gets market share by brand/category
	GetMarketShare(ctx context.Context, timeframe string) ([]*MarketShare, error)

	// GetPopularModels gets most popular models by category
	GetPopularModels(ctx context.Context, category string, limit int) ([]*PopularModel, error)

	// GetPriceDistribution gets price distribution for category
	GetPriceDistribution(ctx context.Context, category string) (*PriceDistribution, error)

	// === Seasonal Data ===

	// GetSeasonalFactors gets seasonal adjustment factors
	GetSeasonalFactors(ctx context.Context, month time.Month) (*SeasonalFactors, error)

	// GetHolidayPricing gets special holiday pricing adjustments
	GetHolidayPricing(ctx context.Context, holiday string) (*HolidayPricing, error)

	// === Regional Data ===

	// GetRegionalPricing gets regional price variations
	GetRegionalPricing(ctx context.Context, region string, brand, model string) (*RegionalPrice, error)

	// GetRegionalDemand gets regional demand index
	GetRegionalDemand(ctx context.Context, region string, category string) (float64, error)

	// === Competitive Intelligence ===

	// GetCompetitorPricing gets competitor insurance pricing
	GetCompetitorPricing(ctx context.Context, deviceValue decimal.Decimal, category string) ([]*CompetitorPrice, error)

	// GetMarketPremiumRates gets average market premium rates
	GetMarketPremiumRates(ctx context.Context, category string) (*PremiumRates, error)

	// === Supply Chain Data ===

	// GetRepairCosts gets average repair costs by component
	GetRepairCosts(ctx context.Context, brand, model string) (*RepairCostMatrix, error)

	// GetPartAvailability gets spare part availability index
	GetPartAvailability(ctx context.Context, brand, model string) (float64, error)

	// GetSupplierPricing gets supplier pricing for parts
	GetSupplierPricing(ctx context.Context, partType string) ([]*SupplierPrice, error)

	// === Risk Data ===

	// GetCategoryRiskFactors gets risk factors for device category
	GetCategoryRiskFactors(ctx context.Context, category string) (*RiskFactors, error)

	// GetModelFailureRates gets failure rates for device model
	GetModelFailureRates(ctx context.Context, brand, model string) (*FailureRates, error)

	// GetTheftStatistics gets theft statistics by model/region
	GetTheftStatistics(ctx context.Context, brand, model, region string) (*TheftStats, error)

	// === Forecasting ===

	// ForecastPrice forecasts future price for device
	ForecastPrice(ctx context.Context, brand, model string, futureDate time.Time) (decimal.Decimal, error)

	// ForecastDemand forecasts demand for category
	ForecastDemand(ctx context.Context, category string, period int) (*DemandForecast, error)

	// === Data Management ===

	// RefreshMarketData refreshes market data from external sources
	RefreshMarketData(ctx context.Context, source string) error

	// GetDataFreshness gets freshness of market data
	GetDataFreshness(ctx context.Context, dataType string) (time.Time, error)

	// ValidateMarketData validates market data integrity
	ValidateMarketData(ctx context.Context) ([]*DataValidationError, error)
}

// === Market Data Models ===

// MarketPrice represents market price data
type MarketPrice struct {
	ID         uuid.UUID       `json:"id"`
	Brand      string          `json:"brand"`
	Model      string          `json:"model"`
	Category   string          `json:"category"`
	Price      decimal.Decimal `json:"price"`
	Currency   string          `json:"currency"`
	Source     string          `json:"source"`
	RecordedAt time.Time       `json:"recorded_at"`
	ValidUntil time.Time       `json:"valid_until"`
}

// MarketTrend represents market trend data
type MarketTrend struct {
	Category      string   `json:"category"`
	TrendType     string   `json:"trend_type"` // rising, falling, stable
	ChangePercent float64  `json:"change_percent"`
	PeriodDays    int      `json:"period_days"`
	Confidence    float64  `json:"confidence"`
	Factors       []string `json:"factors"`
}

// DepreciationCurve represents depreciation data
type DepreciationCurve struct {
	Brand            string          `json:"brand"`
	Model            string          `json:"model"`
	InitialValue     decimal.Decimal `json:"initial_value"`
	CurveType        string          `json:"curve_type"` // linear, exponential, stepped
	AnnualRate       float64         `json:"annual_rate"`
	MonthlyRates     map[int]float64 `json:"monthly_rates"`
	ResidualValue    decimal.Decimal `json:"residual_value"`
	UsefulLifeMonths int             `json:"useful_life_months"`
}

// TradeInMatrix represents trade-in values by condition
type TradeInMatrix struct {
	Brand           string                     `json:"brand"`
	Model           string                     `json:"model"`
	BaseValue       decimal.Decimal            `json:"base_value"`
	ConditionValues map[string]decimal.Decimal `json:"condition_values"` // A, B, C, D, F grades
	AgeAdjustments  map[int]float64            `json:"age_adjustments"`  // by months
	UpdatedAt       time.Time                  `json:"updated_at"`
}

// CompetitorTradeIn represents competitor trade-in offer
type CompetitorTradeIn struct {
	Competitor  string          `json:"competitor"`
	Brand       string          `json:"brand"`
	Model       string          `json:"model"`
	OfferAmount decimal.Decimal `json:"offer_amount"`
	Conditions  []string        `json:"conditions"`
	ValidUntil  time.Time       `json:"valid_until"`
}

// MarketShare represents market share data
type MarketShare struct {
	Brand        string  `json:"brand"`
	Category     string  `json:"category"`
	SharePercent float64 `json:"share_percent"`
	Trend        string  `json:"trend"` // increasing, decreasing, stable
	Quarter      string  `json:"quarter"`
	Year         int     `json:"year"`
}

// PopularModel represents a popular device model
type PopularModel struct {
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Category    string  `json:"category"`
	Popularity  float64 `json:"popularity_score"`
	SalesRank   int     `json:"sales_rank"`
	ReviewScore float64 `json:"review_score"`
	TrendingUp  bool    `json:"trending_up"`
}

// PriceDistribution represents price distribution statistics
type PriceDistribution struct {
	Category     string          `json:"category"`
	Min          decimal.Decimal `json:"min"`
	Max          decimal.Decimal `json:"max"`
	Mean         decimal.Decimal `json:"mean"`
	Median       decimal.Decimal `json:"median"`
	Percentile25 decimal.Decimal `json:"percentile_25"`
	Percentile75 decimal.Decimal `json:"percentile_75"`
	StdDev       decimal.Decimal `json:"std_dev"`
}

// SeasonalFactors represents seasonal adjustment factors
type SeasonalFactors struct {
	Month            time.Month `json:"month"`
	DemandMultiplier float64    `json:"demand_multiplier"`
	PriceMultiplier  float64    `json:"price_multiplier"`
	ClaimMultiplier  float64    `json:"claim_multiplier"`
	Notes            string     `json:"notes"`
}

// HolidayPricing represents holiday pricing adjustments
type HolidayPricing struct {
	Holiday          string    `json:"holiday"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	PriceAdjustment  float64   `json:"price_adjustment"`
	DemandMultiplier float64   `json:"demand_multiplier"`
}

// RegionalPrice represents regional price variation
type RegionalPrice struct {
	Region        string          `json:"region"`
	Brand         string          `json:"brand"`
	Model         string          `json:"model"`
	BasePrice     decimal.Decimal `json:"base_price"`
	RegionalPrice decimal.Decimal `json:"regional_price"`
	Adjustment    float64         `json:"adjustment_percent"`
	Factors       []string        `json:"factors"`
}

// CompetitorPrice represents competitor pricing
type CompetitorPrice struct {
	Competitor       string          `json:"competitor"`
	DeviceValue      decimal.Decimal `json:"device_value"`
	MonthlyPremium   decimal.Decimal `json:"monthly_premium"`
	AnnualPremium    decimal.Decimal `json:"annual_premium"`
	Deductible       decimal.Decimal `json:"deductible"`
	CoverageFeatures []string        `json:"coverage_features"`
}

// PremiumRates represents market premium rates
type PremiumRates struct {
	Category         string  `json:"category"`
	AverageRate      float64 `json:"average_rate"` // as percentage of value
	MinRate          float64 `json:"min_rate"`
	MaxRate          float64 `json:"max_rate"`
	MedianRate       float64 `json:"median_rate"`
	MarketLeaderRate float64 `json:"market_leader_rate"`
}

// RepairCostMatrix represents repair costs by component
type RepairCostMatrix struct {
	Brand             string                     `json:"brand"`
	Model             string                     `json:"model"`
	ComponentCosts    map[string]decimal.Decimal `json:"component_costs"`
	LaborCosts        map[string]decimal.Decimal `json:"labor_costs"`
	AverageRepairCost decimal.Decimal            `json:"average_repair_cost"`
	MaxRepairCost     decimal.Decimal            `json:"max_repair_cost"`
}

// SupplierPrice represents supplier pricing
type SupplierPrice struct {
	Supplier     string          `json:"supplier"`
	PartType     string          `json:"part_type"`
	Brand        string          `json:"brand"`
	UnitPrice    decimal.Decimal `json:"unit_price"`
	BulkPrice    decimal.Decimal `json:"bulk_price"`
	MinQuantity  int             `json:"min_quantity"`
	LeadTimeDays int             `json:"lead_time_days"`
}

// RiskFactors represents risk factors for category
type RiskFactors struct {
	Category         string             `json:"category"`
	TheftRisk        float64            `json:"theft_risk"`
	DamageRisk       float64            `json:"damage_risk"`
	ObsolescenceRisk float64            `json:"obsolescence_risk"`
	RegulatoryRisk   float64            `json:"regulatory_risk"`
	ComponentRisks   map[string]float64 `json:"component_risks"`
}

// FailureRates represents failure rate data
type FailureRates struct {
	Brand              string             `json:"brand"`
	Model              string             `json:"model"`
	OverallFailureRate float64            `json:"overall_failure_rate"`
	ComponentFailures  map[string]float64 `json:"component_failures"`
	AgeFailureCurve    map[int]float64    `json:"age_failure_curve"` // by months
	MTBF               int                `json:"mtbf_hours"`        // Mean Time Between Failures
}

// TheftStats represents theft statistics
type TheftStats struct {
	Brand           string   `json:"brand"`
	Model           string   `json:"model"`
	Region          string   `json:"region"`
	TheftRate       float64  `json:"theft_rate"` // per 1000 devices
	RecoveryRate    float64  `json:"recovery_rate"`
	AverageTimeDays float64  `json:"average_time_days"`
	HotspotAreas    []string `json:"hotspot_areas"`
}

// DemandForecast represents demand forecast
type DemandForecast struct {
	Category       string    `json:"category"`
	ForecastPeriod int       `json:"forecast_period_days"`
	ExpectedDemand float64   `json:"expected_demand"`
	LowerBound     float64   `json:"lower_bound"`
	UpperBound     float64   `json:"upper_bound"`
	Confidence     float64   `json:"confidence"`
	GeneratedAt    time.Time `json:"generated_at"`
}

// DataValidationError represents data validation error
type DataValidationError struct {
	DataType    string    `json:"data_type"`
	ErrorType   string    `json:"error_type"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	RecordCount int       `json:"record_count"`
	DetectedAt  time.Time `json:"detected_at"`
}
