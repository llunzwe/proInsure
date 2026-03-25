package smartwatch

import (
	"math"
	"strings"
	"time"

	"smartsure/internal/domain/models/device/categories/base"
)

// SmartwatchRiskCalculator provides smartwatch-specific risk and valuation calculations
// Integrates with existing device valuation and risk models
type SmartwatchRiskCalculator struct{}

// NewSmartwatchRiskCalculator creates a new smartwatch risk calculator
func NewSmartwatchRiskCalculator() *SmartwatchRiskCalculator {
	return &SmartwatchRiskCalculator{}
}

// CalculateRiskScore calculates risk score for smartwatch
func (c *SmartwatchRiskCalculator) CalculateRiskScore(spec base.CategorySpec, factors base.RiskFactors) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 75.0 // High default risk for invalid spec
	}

	baseRisk := 30.0 // Base risk for smartwatches

	// Usage pattern risk based on usage intensity
	usageRisk := 0.0
	switch factors.UsageIntensity {
	case "high":
		usageRisk = 20.0
	case "medium":
		usageRisk = 12.0
	case "low":
		usageRisk = 5.0
	}

	// Environment risk from custom factors
	envRisk := 0.0
	if outdoor, exists := factors.CustomFactors["outdoor_usage"]; exists && outdoor > 0 {
		envRisk += 10
	}
	if sports, exists := factors.CustomFactors["sports_usage"]; exists && sports > 0 {
		envRisk += 15 // High risk for sports activities
	}
	if water, exists := factors.CustomFactors["water_exposure"]; exists && water > 0 {
		// Adjust based on water resistance
		waterPenalty := 20.0
		switch swSpec.WaterResistance {
		case "10ATM", "100m":
			waterPenalty = 5.0
		case "5ATM", "50m", "IP68":
			waterPenalty = 10.0
		case "3ATM", "30m", "IP67":
			waterPenalty = 15.0
		}
		envRisk += waterPenalty
	}

	// Previous damage increases risk
	damageRisk := float64(factors.PreviousClaims) * 10.0

	// Age-related risk
	ageMonths := time.Since(swSpec.PurchaseDate).Hours() / 24 / 30
	ageRisk := ageMonths * 0.5
	if ageRisk > 20 {
		ageRisk = 20
	}

	// Material quality reduces risk
	materialMitigation := 0.0
	switch swSpec.CaseMaterial {
	case "titanium":
		materialMitigation = -15
	case "steel", "stainless_steel":
		materialMitigation = -10
	case "aluminum":
		materialMitigation = -5
	}

	// Health monitoring features (increased value = slightly higher risk)
	healthRisk := 0.0
	if swSpec.ECGSensor || swSpec.BloodOxygenSensor {
		healthRisk = 5.0 // Premium features attract more attention/theft
	}

	totalRisk := baseRisk + usageRisk + envRisk + damageRisk + ageRisk + materialMitigation + healthRisk

	// Normalize to 0-100 scale
	if totalRisk < 0 {
		totalRisk = 0
	}
	if totalRisk > 100 {
		totalRisk = 100
	}

	return totalRisk
}

// CalculateDepreciation calculates depreciated value for smartwatch
func (c *SmartwatchRiskCalculator) CalculateDepreciation(originalValue float64, purchaseDate time.Time) float64 {
	ageYears := time.Since(purchaseDate).Hours() / 24 / 365

	// Smartwatches depreciate faster than smartphones due to:
	// 1. Battery degradation (non-replaceable)
	// 2. Fashion/style changes
	// 3. Rapid feature evolution

	depreciation := 0.0

	if ageYears < 1 {
		// First year: 35% depreciation
		depreciation = originalValue * (0.35 * ageYears)
	} else if ageYears < 2 {
		// Second year: additional 25%
		firstYear := originalValue * 0.35
		secondYear := originalValue * 0.25 * (ageYears - 1)
		depreciation = firstYear + secondYear
	} else if ageYears < 3 {
		// Third year: additional 15%
		firstYear := originalValue * 0.35
		secondYear := originalValue * 0.25
		thirdYear := originalValue * 0.15 * (ageYears - 2)
		depreciation = firstYear + secondYear + thirdYear
	} else {
		// After 3 years: 5% per year, minimum 10% value
		firstThreeYears := originalValue * 0.75
		additionalYears := originalValue * 0.05 * (ageYears - 3)
		depreciation = firstThreeYears + additionalYears
		if depreciation > originalValue*0.9 {
			depreciation = originalValue * 0.9
		}
	}

	currentValue := originalValue - depreciation
	if currentValue < originalValue*0.1 {
		currentValue = originalValue * 0.1 // Minimum 10% residual value
	}

	return currentValue
}

// CalculateResaleValue calculates resale value based on condition
func (c *SmartwatchRiskCalculator) CalculateResaleValue(spec base.CategorySpec, condition string) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 0
	}

	// Get depreciated value first
	depreciatedValue := c.CalculateDepreciation(swSpec.MarketValue, swSpec.PurchaseDate)

	// Apply condition multiplier
	conditionMultiplier := 1.0
	switch condition {
	case "excellent", "mint", "like_new":
		conditionMultiplier = 0.95
	case "good":
		conditionMultiplier = 0.80
	case "fair":
		conditionMultiplier = 0.60
	case "poor":
		conditionMultiplier = 0.40
	case "broken", "parts":
		conditionMultiplier = 0.20
	}

	// Brand value retention
	brandMultiplier := 0.7 // Default
	switch swSpec.Manufacturer {
	case "Apple":
		brandMultiplier = 0.9 // Apple Watches retain value well
	case "Samsung", "Garmin":
		brandMultiplier = 0.8
	case "Fitbit", "Amazfit":
		brandMultiplier = 0.6
	}

	// Special features add value
	featureBonus := 1.0
	if swSpec.Cellular {
		featureBonus += 0.1
	}
	if swSpec.ECGSensor {
		featureBonus += 0.05
	}
	if swSpec.CaseMaterial == "titanium" || swSpec.CaseMaterial == "steel" {
		featureBonus += 0.1
	}

	resaleValue := depreciatedValue * conditionMultiplier * brandMultiplier * featureBonus

	// Round to nearest $5
	return math.Round(resaleValue/5) * 5
}

// CalculateRepairCost estimates repair cost for issues
func (c *SmartwatchRiskCalculator) CalculateRepairCost(spec base.CategorySpec, issueType string) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 100 // Default repair cost
	}

	// Base repair costs as percentage of device value
	repairCosts := map[string]float64{
		"screen_crack":       swSpec.MarketValue * 0.35,
		"screen_replacement": swSpec.MarketValue * 0.40,
		"battery_issue":      swSpec.MarketValue * 0.30,
		"water_damage":       swSpec.MarketValue * 0.60,
		"band_broken":        30.0, // Fixed cost for bands
		"button_stuck":       swSpec.MarketValue * 0.20,
		"sensor_failure":     swSpec.MarketValue * 0.35,
		"charging_issue":     swSpec.MarketValue * 0.25,
		"software_issue":     50.0, // Fixed cost for software repairs
		"speaker_issue":      swSpec.MarketValue * 0.25,
		"vibration_motor":    swSpec.MarketValue * 0.20,
	}

	cost, exists := repairCosts[issueType]
	if !exists {
		// Default to 30% of value for unknown issues
		cost = swSpec.MarketValue * 0.30
	}

	// Premium materials increase repair cost
	if swSpec.CaseMaterial == "titanium" || swSpec.CaseMaterial == "ceramic" {
		cost *= 1.3
	}

	// Round to nearest $10
	return math.Round(cost/10) * 10
}

// GetDepreciationCurve returns depreciation curve data for smartwatch
func (c *SmartwatchRiskCalculator) GetDepreciationCurve(spec base.CategorySpec) base.DepreciationCurve {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return base.DepreciationCurve{}
	}

	return base.DepreciationCurve{
		Model:        swSpec.Model,
		InitialValue: swSpec.MarketValue,
		CurveType:    "exponential",
		AnnualRate:   0.35, // 35% annual depreciation for smartwatches
		MinValue:     swSpec.MarketValue * 0.10,
		DataPoints: []base.DataPoint{
			{Age: 0, Value: 1.00},
			{Age: 6, Value: 0.75},
			{Age: 12, Value: 0.65},
			{Age: 18, Value: 0.52},
			{Age: 24, Value: 0.40},
			{Age: 30, Value: 0.32},
			{Age: 36, Value: 0.25},
			{Age: 48, Value: 0.18},
			{Age: 60, Value: 0.10},
		},
		LastUpdated: time.Now(),
	}
}

// EstimateMarketDemand estimates market demand score
func (c *SmartwatchRiskCalculator) EstimateMarketDemand(spec base.CategorySpec) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 50.0
	}

	demandScore := 50.0 // Base demand

	// Brand popularity
	switch swSpec.Manufacturer {
	case "Apple":
		demandScore += 30
	case "Samsung", "Garmin":
		demandScore += 20
	case "Fitbit", "Huawei":
		demandScore += 10
	}

	// Feature demand
	if swSpec.Cellular {
		demandScore += 10
	}
	if swSpec.ECGSensor || swSpec.BloodOxygenSensor {
		demandScore += 5
	}
	if swSpec.GPSBuiltIn {
		demandScore += 5
	}

	// Age impact on demand
	ageYears := time.Since(swSpec.PurchaseDate).Hours() / 24 / 365
	if ageYears < 1 {
		demandScore += 10
	} else if ageYears > 2 {
		demandScore -= 15
	}

	// Normalize to 0-100
	if demandScore > 100 {
		demandScore = 100
	}
	if demandScore < 0 {
		demandScore = 0
	}

	return demandScore
}

// CalculateMaintenanceCost estimates annual maintenance cost for smartwatch
func (c *SmartwatchRiskCalculator) CalculateMaintenanceCost(spec base.CategorySpec) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 50.0 // Default maintenance cost
	}

	baseCost := 30.0 // Base annual maintenance cost

	// Screen protection cost (smartwatches need regular screen care)
	if swSpec.ScreenType == "AMOLED" || swSpec.ScreenType == "OLED" {
		baseCost += 15.0 // Premium displays need more care
	}

	// Battery maintenance (smartwatches have sealed batteries)
	if swSpec.BatteryCapacity < 300 {
		baseCost += 10.0 // Smaller batteries may need more frequent care
	}

	// Health sensors maintenance
	if len(swSpec.HealthFeatures) > 5 {
		baseCost += 12.0 // More sensors = more maintenance
	}

	// Waterproof models need more maintenance
	if swSpec.WaterResistance == "IP68" || swSpec.WaterResistance == "5ATM" {
		baseCost += 8.0
	}

	// Connectivity complexity
	if swSpec.GPSBuiltIn && swSpec.Cellular {
		baseCost += 10.0 // Dual connectivity increases maintenance
	}

	// Premium brands have higher maintenance costs
	switch swSpec.Manufacturer {
	case "Apple", "Garmin", "Fitbit":
		baseCost *= 1.3 // 30% premium for brand reputation
	}

	return baseCost
}

// CalculateRepairCostEstimate estimates repair costs for smartwatch components
func (c *SmartwatchRiskCalculator) CalculateRepairCostEstimate(component string, spec base.CategorySpec) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 100.0 // Default estimate
	}

	component = strings.ToLower(component)

	// Base cost estimates for smartwatch components
	switch component {
	case "screen", "display":
		// Screen repair - most expensive component
		baseCost := 150.0
		if swSpec.ScreenType == "AMOLED" || swSpec.ScreenType == "OLED" {
			baseCost *= 1.3 // Premium display technology
		}
		if swSpec.PPI > 300 {
			baseCost *= 1.2 // High PPI screens cost more
		}
		return baseCost

	case "battery":
		baseCost := 80.0
		if swSpec.BatteryCapacity < 200 {
			baseCost *= 0.8 // Smaller batteries cheaper
		} else if swSpec.BatteryCapacity > 400 {
			baseCost *= 1.2 // Larger batteries more expensive
		}
		return baseCost

	case "charging_port", "charging_cable":
		return 60.0

	case "case", "housing":
		baseCost := 100.0
		if swSpec.WaterResistance == "IP68" || swSpec.WaterResistance == "5ATM" {
			baseCost *= 1.4 // Waterproof cases more expensive
		}
		return baseCost

	case "strap", "band", "wristband":
		return 40.0

	case "buttons", "crown", "digital_crown":
		return 35.0

	case "sensors", "health_sensors":
		// Depends on which sensors
		for _, feature := range swSpec.HealthFeatures {
			if strings.Contains(strings.ToLower(feature), "ecg") {
				return 90.0 // ECG sensor expensive
			}
		}
		return 50.0

	case "gps", "cellular":
		return 75.0 // Connectivity module repair

	case "mainboard", "logic_board":
		return 200.0 // Major repair

	default:
		return 50.0 // Default for unknown components
	}
}

// CalculateTotalCostOfOwnership calculates total cost of ownership over years
func (c *SmartwatchRiskCalculator) CalculateTotalCostOfOwnership(spec base.CategorySpec, years int) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 500.0 // Default TCO
	}

	// Base purchase price (estimated)
	purchasePrice := 300.0 // Average smartwatch price
	if swSpec.Manufacturer == "Apple" {
		purchasePrice = 400.0
	} else if swSpec.Manufacturer == "Garmin" {
		purchasePrice = 350.0
	}

	// Annual maintenance cost
	annualMaintenance := c.CalculateMaintenanceCost(spec)

	// Expected repairs over lifetime (based on usage patterns)
	expectedRepairs := 0.5 * float64(years) // 0.5 repairs per year on average

	// Average repair cost
	avgRepairCost := 75.0

	// Insurance premiums (estimated monthly)
	monthlyPremium := purchasePrice * 0.015 // 1.5% of value per month

	// Total cost calculation
	totalCost := purchasePrice
	totalCost += annualMaintenance * float64(years)
	totalCost += expectedRepairs * avgRepairCost
	totalCost += monthlyPremium * 12 * float64(years)

	return totalCost
}

// CalculateUpgradeValue calculates trade-in value when upgrading to a new model
func (c *SmartwatchRiskCalculator) CalculateUpgradeValue(spec base.CategorySpec, targetModel string) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 50.0 // Default trade-in value
	}

	// Base current value (simplified calculation)
	baseValue := 200.0 // Base smartwatch value
	if swSpec.Manufacturer == "Apple" {
		baseValue = 250.0
	} else if swSpec.Manufacturer == "Garmin" {
		baseValue = 220.0
	}

	// Age depreciation (assume 2 years average ownership)
	ageDepreciation := 0.3 // 30% depreciation per year
	currentValue := baseValue * (1.0 - ageDepreciation*2)

	// Condition adjustment
	conditionMultiplier := 1.0
	if strings.Contains(swSpec.Condition, "good") {
		conditionMultiplier = 0.9
	} else if strings.Contains(swSpec.Condition, "fair") {
		conditionMultiplier = 0.7
	} else if strings.Contains(swSpec.Condition, "poor") {
		conditionMultiplier = 0.5
	}

	currentValue *= conditionMultiplier

	// Target model adjustment (upgrading to premium models increases trade-in value)
	targetMultiplier := 1.0
	if strings.Contains(targetModel, "Ultra") || strings.Contains(targetModel, "Pro") {
		targetMultiplier = 1.2
	} else if strings.Contains(targetModel, "SE") || strings.Contains(targetModel, "Mini") {
		targetMultiplier = 0.9
	}

	return currentValue * targetMultiplier
}

// CalculateWarrantyCost calculates extended warranty cost for specified months
func (c *SmartwatchRiskCalculator) CalculateWarrantyCost(spec base.CategorySpec, months int) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 50.0 // Default warranty cost
	}

	// Base warranty cost per month
	baseCostPerMonth := 5.0
	if swSpec.Manufacturer == "Apple" {
		baseCostPerMonth = 8.0
	} else if swSpec.Manufacturer == "Garmin" {
		baseCostPerMonth = 6.0
	}

	// Adjust for device features (proxy for value)
	featureMultiplier := 1.0
	if swSpec.ECGSensor {
		featureMultiplier += 0.3 // Premium feature
	}
	if swSpec.BloodOxygenSensor {
		featureMultiplier += 0.2
	}
	if swSpec.ScreenSize > 40 {
		featureMultiplier += 0.2 // Larger screen = premium
	}

	baseCostPerMonth *= featureMultiplier

	// Calculate total cost
	totalCost := baseCostPerMonth * float64(months/12)

	// Apply bulk discount for longer terms
	if months >= 24 {
		totalCost *= 0.9 // 10% discount for 2+ years
	}

	return totalCost
}

// GetPriceHistory returns historical price data for smartwatch models
func (c *SmartwatchRiskCalculator) GetPriceHistory(model string) []base.PricePoint {
	// Mock historical price data for smartwatches
	// In a real implementation, this would query a price history database
	now := time.Now()

	// Base price based on model
	basePrice := 300.0
	if strings.Contains(model, "Apple") {
		basePrice = 400.0
	} else if strings.Contains(model, "Garmin") {
		basePrice = 350.0
	} else if strings.Contains(model, "Samsung") {
		basePrice = 320.0
	}

	// Generate historical prices over 2 years
	prices := []base.PricePoint{}
	for i := 24; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)

		// Seasonal adjustment
		seasonalMultiplier := 1.0
		month := date.Month()
		if month == 11 || month == 12 { // Holiday season
			seasonalMultiplier = 1.1
		} else if month == 1 || month == 2 { // Post-holiday dip
			seasonalMultiplier = 0.95
		}

		// Depreciation over time (slight decrease)
		timeMultiplier := 1.0 - float64(i)*0.02 // 2% depreciation per month

		price := basePrice * seasonalMultiplier * timeMultiplier

		prices = append(prices, base.PricePoint{
			Date:  date,
			Price: price,
		})
	}

	return prices
}
