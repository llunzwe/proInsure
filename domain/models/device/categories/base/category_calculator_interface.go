package base

import (
	"time"
)

// CategoryCalculator defines calculation interface for category-specific operations
type CategoryCalculator interface {
	// CalculateDepreciation calculates current depreciated value
	CalculateDepreciation(originalValue float64, purchaseDate time.Time) float64

	// CalculateResaleValue estimates current resale value
	CalculateResaleValue(spec CategorySpec, condition string) float64

	// CalculateRepairCostEstimate estimates repair costs
	CalculateRepairCostEstimate(component string, spec CategorySpec) float64

	// CalculateTotalCostOfOwnership calculates TCO
	CalculateTotalCostOfOwnership(spec CategorySpec, years int) float64

	// CalculateUpgradeValue calculates trade-in value for upgrades
	CalculateUpgradeValue(spec CategorySpec, targetModel string) float64

	// CalculateRiskScore calculates device-specific risk score
	CalculateRiskScore(spec CategorySpec, factors RiskFactors) float64

	// CalculateMaintenanceCost estimates annual maintenance cost
	CalculateMaintenanceCost(spec CategorySpec) float64

	// CalculateWarrantyCost calculates extended warranty cost
	CalculateWarrantyCost(spec CategorySpec, months int) float64

	// GetDepreciationCurve returns depreciation curve data
	GetDepreciationCurve(spec CategorySpec) DepreciationCurve

	// GetPriceHistory returns historical price data
	GetPriceHistory(model string) []PricePoint
}

// RiskFactors contains factors for risk calculation
type RiskFactors struct {
	DeviceAge         float64            `json:"device_age_months"`
	UsageIntensity    string             `json:"usage_intensity"` // low, medium, high
	PreviousClaims    int                `json:"previous_claims"`
	EnvironmentalRisk string             `json:"environmental_risk"`
	HandlingRisk      string             `json:"handling_risk"`
	LocationRisk      float64            `json:"location_risk_score"`
	UserProfile       string             `json:"user_profile"`
	ProtectionLevel   string             `json:"protection_level"`
	CustomFactors     map[string]float64 `json:"custom_factors"`
}

// DepreciationCurve represents device depreciation over time
type DepreciationCurve struct {
	Model        string      `json:"model"`
	InitialValue float64     `json:"initial_value"`
	CurveType    string      `json:"curve_type"` // linear, exponential, stepped
	AnnualRate   float64     `json:"annual_rate"`
	MinValue     float64     `json:"min_value"`
	DataPoints   []DataPoint `json:"data_points"`
	LastUpdated  time.Time   `json:"last_updated"`
}

// DataPoint represents a point on depreciation curve
type DataPoint struct {
	Age   int     `json:"age_months"`
	Value float64 `json:"value"`
	Note  string  `json:"note,omitempty"`
}

// PricePoint represents historical price data
type PricePoint struct {
	Date       time.Time `json:"date"`
	Price      float64   `json:"price"`
	Source     string    `json:"source"`
	Condition  string    `json:"condition"`
	MarketType string    `json:"market_type"` // retail, wholesale, refurbished
}

// BaseCalculator provides common calculation functionality
type BaseCalculator struct {
	DepreciationRates map[CategoryType]float64
	RiskWeights       map[string]float64
	MaintenanceRates  map[CategoryType]float64
	WarrantyRates     map[CategoryType]float64
}

// NewBaseCalculator creates a new base calculator
func NewBaseCalculator() *BaseCalculator {
	return &BaseCalculator{
		DepreciationRates: map[CategoryType]float64{
			CategorySmartphone: 0.30, // 30% per year
			CategorySmartwatch: 0.35, // 35% per year
			CategoryTablet:     0.25, // 25% per year
			CategoryLaptop:     0.20, // 20% per year
			CategoryWearable:   0.40, // 40% per year
			CategoryIoT:        0.35, // 35% per year
			CategoryAccessory:  0.50, // 50% per year
		},
		RiskWeights: map[string]float64{
			"device_age":      0.20,
			"usage_intensity": 0.15,
			"previous_claims": 0.25,
			"environmental":   0.10,
			"handling":        0.10,
			"location":        0.10,
			"protection":      0.10,
		},
		MaintenanceRates: map[CategoryType]float64{
			CategorySmartphone: 0.05, // 5% of value per year
			CategorySmartwatch: 0.04, // 4% of value per year
			CategoryTablet:     0.03, // 3% of value per year
			CategoryLaptop:     0.06, // 6% of value per year
		},
		WarrantyRates: map[CategoryType]float64{
			CategorySmartphone: 0.10, // 10% of value for extended warranty
			CategorySmartwatch: 0.12, // 12% of value for extended warranty
			CategoryTablet:     0.08, // 8% of value for extended warranty
			CategoryLaptop:     0.15, // 15% of value for extended warranty
		},
	}
}

// CalculateDepreciation calculates depreciated value
func (b *BaseCalculator) CalculateDepreciation(originalValue float64, purchaseDate time.Time, category CategoryType) float64 {
	ageInYears := time.Since(purchaseDate).Hours() / 24 / 365
	rate, exists := b.DepreciationRates[category]
	if !exists {
		rate = 0.30 // Default 30%
	}

	depreciation := originalValue * (1 - (rate * ageInYears))
	minValue := originalValue * 0.10 // Minimum 10% of original value

	if depreciation < minValue {
		return minValue
	}
	return depreciation
}

// CalculateRiskScore calculates weighted risk score
func (b *BaseCalculator) CalculateRiskScore(factors RiskFactors) float64 {
	score := 0.0

	// Age factor
	ageScore := factors.DeviceAge * 2.0 // 2 points per month
	if ageScore > 100 {
		ageScore = 100
	}
	score += ageScore * b.RiskWeights["device_age"]

	// Usage intensity factor
	usageScore := 0.0
	switch factors.UsageIntensity {
	case "high":
		usageScore = 80.0
	case "medium":
		usageScore = 50.0
	case "low":
		usageScore = 20.0
	}
	score += usageScore * b.RiskWeights["usage_intensity"]

	// Previous claims factor
	claimsScore := float64(factors.PreviousClaims) * 30.0
	if claimsScore > 100 {
		claimsScore = 100
	}
	score += claimsScore * b.RiskWeights["previous_claims"]

	// Add other factors
	score += factors.LocationRisk * b.RiskWeights["location"]

	// Apply custom factors if any
	for key, value := range factors.CustomFactors {
		if weight, exists := b.RiskWeights[key]; exists {
			score += value * weight
		}
	}

	// Normalize to 0-100 scale
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// CalculationResult represents the result of a calculation
type CalculationResult struct {
	Value      float64            `json:"value"`
	Currency   string             `json:"currency"`
	Breakdown  map[string]float64 `json:"breakdown"`
	Confidence float64            `json:"confidence"`
	ValidUntil time.Time          `json:"valid_until"`
	Notes      []string           `json:"notes"`
}

// MarketAnalysis contains market analysis data
type MarketAnalysis struct {
	AveragePrice float64    `json:"average_price"`
	MedianPrice  float64    `json:"median_price"`
	PriceRange   [2]float64 `json:"price_range"`
	Volatility   float64    `json:"volatility"`
	Trend        string     `json:"trend"`        // rising, falling, stable
	DemandLevel  string     `json:"demand_level"` // high, medium, low
	SupplyLevel  string     `json:"supply_level"`
	LastUpdated  time.Time  `json:"last_updated"`
}
