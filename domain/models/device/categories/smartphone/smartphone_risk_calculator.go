package smartphone

import (
	"math"
	"time"

	"smartsure/internal/domain/models/device/categories/base"
)

// SmartphoneRiskCalculator calculates risk scores for smartphones
type SmartphoneRiskCalculator struct {
	base.BaseCalculator
}

// NewSmartphoneRiskCalculator creates a new risk calculator
func NewSmartphoneRiskCalculator() *SmartphoneRiskCalculator {
	calc := &SmartphoneRiskCalculator{
		BaseCalculator: *base.NewBaseCalculator(),
	}

	// Customize risk weights for smartphones
	calc.RiskWeights["screen_damage"] = 0.30
	calc.RiskWeights["water_damage"] = 0.20
	calc.RiskWeights["theft_risk"] = 0.15
	calc.RiskWeights["drop_risk"] = 0.25
	calc.RiskWeights["battery_risk"] = 0.10

	return calc
}

// CalculateRiskScore calculates comprehensive risk score for a smartphone
func (c *SmartphoneRiskCalculator) CalculateRiskScore(spec base.CategorySpec, factors base.RiskFactors) float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 50.0 // Default medium risk
	}

	score := 0.0

	// Screen damage risk
	screenRisk := c.calculateScreenRisk(phoneSpec, factors)
	score += screenRisk * c.RiskWeights["screen_damage"]

	// Water damage risk
	waterRisk := c.calculateWaterRisk(phoneSpec, factors)
	score += waterRisk * c.RiskWeights["water_damage"]

	// Theft risk
	theftRisk := c.calculateTheftRisk(phoneSpec, factors)
	score += theftRisk * c.RiskWeights["theft_risk"]

	// Drop/physical damage risk
	dropRisk := c.calculateDropRisk(phoneSpec, factors)
	score += dropRisk * c.RiskWeights["drop_risk"]

	// Battery degradation risk
	batteryRisk := c.calculateBatteryRisk(phoneSpec, factors)
	score += batteryRisk * c.RiskWeights["battery_risk"]

	// Apply age factor
	ageMultiplier := c.calculateAgeFactor(phoneSpec.ReleaseDate)
	score *= ageMultiplier

	// Apply user profile factor
	userMultiplier := c.getUserProfileMultiplier(factors.UserProfile)
	score *= userMultiplier

	// Normalize score to 0-100
	score = math.Min(math.Max(score, 0), 100)

	return score
}

// calculateScreenRisk assesses screen damage risk
func (c *SmartphoneRiskCalculator) calculateScreenRisk(spec *SmartphoneSpec, factors base.RiskFactors) float64 {
	risk := 0.0

	// Larger screens have higher risk
	sizeRisk := (spec.ScreenSize - 5.0) * 20.0
	if sizeRisk < 0 {
		sizeRisk = 0
	}
	risk += sizeRisk

	// Foldable screens are much higher risk
	if spec.IsFoldable() {
		risk += 40.0
	}

	// Protection reduces risk
	if spec.ScreenProtection != "" {
		risk -= 10.0
		if spec.ScreenProtection == "Gorilla Glass Victus" ||
			spec.ScreenProtection == "Gorilla Glass Victus+" {
			risk -= 5.0
		}
	}

	// User protection level
	switch factors.ProtectionLevel {
	case "high":
		risk *= 0.5
	case "medium":
		risk *= 0.75
	case "low", "none":
		risk *= 1.2
	}

	// Usage intensity
	switch factors.UsageIntensity {
	case "high":
		risk *= 1.3
	case "medium":
		risk *= 1.0
	case "low":
		risk *= 0.7
	}

	return math.Min(risk, 100)
}

// calculateWaterRisk assesses water damage risk
func (c *SmartphoneRiskCalculator) calculateWaterRisk(spec *SmartphoneSpec, factors base.RiskFactors) float64 {
	risk := 50.0 // Base water risk

	// Water resistance rating reduces risk
	switch spec.WaterResistance {
	case "IP68":
		risk *= 0.3
	case "IP67":
		risk *= 0.4
	case "IP54", "IP55":
		risk *= 0.7
	case "none", "":
		risk *= 1.2
	}

	// Environmental factors
	switch factors.EnvironmentalRisk {
	case "high": // Humid, coastal, rainy areas
		risk *= 1.5
	case "medium":
		risk *= 1.0
	case "low": // Dry areas
		risk *= 0.7
	}

	// User behavior
	if factors.HandlingRisk == "high" {
		risk *= 1.3
	}

	return math.Min(risk, 100)
}

// calculateTheftRisk assesses theft risk
func (c *SmartphoneRiskCalculator) calculateTheftRisk(spec *SmartphoneSpec, factors base.RiskFactors) float64 {
	risk := 0.0

	// High-value devices have higher theft risk
	if spec.MarketValue > 1000 {
		risk += 40.0
	} else if spec.MarketValue > 500 {
		risk += 25.0
	} else {
		risk += 10.0
	}

	// Popular models have higher theft risk
	if spec.IsHighEnd() {
		risk += 20.0
	}

	// Location risk
	risk += factors.LocationRisk

	// User profile
	switch factors.UserProfile {
	case "student":
		risk *= 1.2
	case "professional":
		risk *= 0.9
	case "senior":
		risk *= 0.7
	}

	return math.Min(risk, 100)
}

// calculateDropRisk assesses physical damage risk
func (c *SmartphoneRiskCalculator) calculateDropRisk(spec *SmartphoneSpec, factors base.RiskFactors) float64 {
	risk := 40.0 // Base drop risk

	// Build materials affect risk
	glassCount := 0
	for _, material := range spec.Materials {
		if material == "glass" {
			glassCount++
		}
	}
	risk += float64(glassCount) * 10.0

	// Size and weight factor
	if spec.ScreenSize > 6.5 {
		risk += 10.0 // Harder to handle
	}

	// Protection level
	switch factors.ProtectionLevel {
	case "high": // Case + screen protector
		risk *= 0.4
	case "medium": // Basic case
		risk *= 0.7
	case "low", "none":
		risk *= 1.3
	}

	// User handling risk
	switch factors.HandlingRisk {
	case "high":
		risk *= 1.5
	case "medium":
		risk *= 1.0
	case "low":
		risk *= 0.6
	}

	// Usage intensity
	switch factors.UsageIntensity {
	case "high":
		risk *= 1.2
	case "medium":
		risk *= 1.0
	case "low":
		risk *= 0.8
	}

	return math.Min(risk, 100)
}

// calculateBatteryRisk assesses battery degradation risk
func (c *SmartphoneRiskCalculator) calculateBatteryRisk(spec *SmartphoneSpec, factors base.RiskFactors) float64 {
	risk := 0.0

	// Age is primary factor for battery risk
	ageInMonths := time.Since(spec.ReleaseDate).Hours() / 24 / 30
	risk += ageInMonths * 2.0 // 2 points per month

	// Fast charging increases battery wear
	if spec.ChargingWattage > 30 {
		risk += 10.0
	}

	// Usage intensity affects battery life
	switch factors.UsageIntensity {
	case "high":
		risk += 20.0
	case "medium":
		risk += 10.0
	case "low":
		risk += 5.0
	}

	// Battery capacity (larger batteries may degrade faster)
	if spec.BatteryCapacity > 5000 {
		risk += 5.0
	}

	return math.Min(risk, 100)
}

// calculateAgeFactor returns age-based risk multiplier
func (c *SmartphoneRiskCalculator) calculateAgeFactor(releaseDate time.Time) float64 {
	ageInMonths := time.Since(releaseDate).Hours() / 24 / 30

	if ageInMonths < 6 {
		return 0.8 // New devices have lower risk
	} else if ageInMonths < 12 {
		return 0.9
	} else if ageInMonths < 24 {
		return 1.0
	} else if ageInMonths < 36 {
		return 1.2
	} else if ageInMonths < 48 {
		return 1.4
	}
	return 1.6 // Devices older than 4 years
}

// getUserProfileMultiplier returns user profile risk multiplier
func (c *SmartphoneRiskCalculator) getUserProfileMultiplier(profile string) float64 {
	multipliers := map[string]float64{
		"careful":      0.7,
		"average":      1.0,
		"careless":     1.5,
		"professional": 0.8,
		"student":      1.2,
		"senior":       0.9,
		"child":        1.8,
	}

	if mult, exists := multipliers[profile]; exists {
		return mult
	}
	return 1.0
}

// CalculateDepreciation calculates current value after depreciation
func (c *SmartphoneRiskCalculator) CalculateDepreciation(originalValue float64, purchaseDate time.Time) float64 {
	return c.BaseCalculator.CalculateDepreciation(originalValue, purchaseDate, base.CategorySmartphone)
}

// CalculateResaleValue estimates resale value
func (c *SmartphoneRiskCalculator) CalculateResaleValue(spec base.CategorySpec, condition string) float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0
	}

	// Start with depreciated value
	currentValue := c.CalculateDepreciation(phoneSpec.MarketValue, phoneSpec.ReleaseDate)

	// Apply condition multiplier
	conditionMultipliers := map[string]float64{
		"mint":      0.95,
		"excellent": 0.85,
		"good":      0.70,
		"fair":      0.50,
		"poor":      0.30,
		"broken":    0.10,
	}

	if mult, exists := conditionMultipliers[condition]; exists {
		currentValue *= mult
	} else {
		currentValue *= 0.5
	}

	// High-end phones retain value better
	if phoneSpec.IsHighEnd() {
		currentValue *= 1.1
	}

	// Popular brands retain value better
	brandMultipliers := map[string]float64{
		"Apple":   1.15,
		"Samsung": 1.05,
		"Google":  1.05,
		"OnePlus": 0.95,
	}

	if mult, exists := brandMultipliers[phoneSpec.Manufacturer]; exists {
		currentValue *= mult
	}

	return currentValue
}

// CalculateRepairCostEstimate estimates repair cost for a component
func (c *SmartphoneRiskCalculator) CalculateRepairCostEstimate(component string, spec base.CategorySpec) float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0
	}

	costs := phoneSpec.GetRepairCosts()
	if cost, exists := costs[component]; exists {
		return cost
	}

	// Default estimate based on device value
	return phoneSpec.MarketValue * 0.2
}

// CalculateTotalCostOfOwnership calculates TCO over specified years
func (c *SmartphoneRiskCalculator) CalculateTotalCostOfOwnership(spec base.CategorySpec, years int) float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0
	}

	tco := phoneSpec.MarketValue // Initial cost

	// Annual insurance (estimated at 8% of value)
	tco += phoneSpec.MarketValue * 0.08 * float64(years)

	// Expected repairs (based on risk)
	riskScore := c.CalculateRiskScore(spec, base.RiskFactors{
		UsageIntensity: "medium",
		UserProfile:    "average",
	})
	expectedRepairs := (riskScore / 100) * phoneSpec.MarketValue * 0.3 * float64(years)
	tco += expectedRepairs

	// Accessories (cases, screen protectors, etc.)
	tco += 50.0 * float64(years)

	// Minus resale value
	resaleValue := c.CalculateResaleValue(spec, "good")
	tco -= resaleValue

	return tco
}

// CalculateUpgradeValue calculates trade-in value for upgrade
func (c *SmartphoneRiskCalculator) CalculateUpgradeValue(spec base.CategorySpec, targetModel string) float64 {
	// Base it on resale value with bonus for brand loyalty
	resaleValue := c.CalculateResaleValue(spec, "good")

	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return resaleValue
	}

	// Add loyalty bonus if staying with same brand
	// (In real implementation, would check targetModel brand)
	loyaltyBonus := resaleValue * 0.1

	// Additional bonus for newer devices
	ageInMonths := time.Since(phoneSpec.ReleaseDate).Hours() / 24 / 30
	if ageInMonths < 12 {
		loyaltyBonus += resaleValue * 0.05
	}

	return resaleValue + loyaltyBonus
}

// CalculateMaintenanceCost estimates annual maintenance cost
func (c *SmartphoneRiskCalculator) CalculateMaintenanceCost(spec base.CategorySpec) float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0
	}

	// Base maintenance as percentage of value
	rate := c.MaintenanceRates[base.CategorySmartphone]
	maintenanceCost := phoneSpec.MarketValue * rate

	// Adjust for device age
	ageInYears := time.Since(phoneSpec.ReleaseDate).Hours() / 24 / 365
	if ageInYears > 2 {
		maintenanceCost *= 1.5
	} else if ageInYears > 3 {
		maintenanceCost *= 2.0
	}

	return maintenanceCost
}

// CalculateWarrantyCost calculates extended warranty cost
func (c *SmartphoneRiskCalculator) CalculateWarrantyCost(spec base.CategorySpec, months int) float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0
	}

	// Base warranty rate
	rate := c.WarrantyRates[base.CategorySmartphone]
	yearlyCost := phoneSpec.MarketValue * rate
	monthlyCost := yearlyCost / 12

	// Adjust for high-end devices
	if phoneSpec.IsHighEnd() {
		monthlyCost *= 1.2
	}

	// Adjust for foldables
	if phoneSpec.IsFoldable() {
		monthlyCost *= 1.5
	}

	return monthlyCost * float64(months)
}

// GetDepreciationCurve returns depreciation curve for the device
func (c *SmartphoneRiskCalculator) GetDepreciationCurve(spec base.CategorySpec) base.DepreciationCurve {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return base.DepreciationCurve{}
	}

	curve := base.DepreciationCurve{
		Model:        phoneSpec.Model,
		InitialValue: phoneSpec.MarketValue,
		CurveType:    "exponential",
		AnnualRate:   c.DepreciationRates[base.CategorySmartphone],
		MinValue:     phoneSpec.MarketValue * 0.1,
		LastUpdated:  time.Now(),
		DataPoints:   []base.DataPoint{},
	}

	// Generate data points for 5 years
	for month := 0; month <= 60; month += 6 {
		ageInYears := float64(month) / 12.0
		value := phoneSpec.MarketValue * math.Pow(1-curve.AnnualRate, ageInYears)
		if value < curve.MinValue {
			value = curve.MinValue
		}

		curve.DataPoints = append(curve.DataPoints, base.DataPoint{
			Age:   month,
			Value: value,
		})
	}

	return curve
}

// GetPriceHistory returns historical price data (mock implementation)
func (c *SmartphoneRiskCalculator) GetPriceHistory(model string) []base.PricePoint {
	// In real implementation, this would fetch from database
	return []base.PricePoint{
		{
			Date:       time.Now().AddDate(0, -6, 0),
			Price:      1000,
			Source:     "market_average",
			Condition:  "new",
			MarketType: "retail",
		},
		{
			Date:       time.Now().AddDate(0, -3, 0),
			Price:      900,
			Source:     "market_average",
			Condition:  "new",
			MarketType: "retail",
		},
		{
			Date:       time.Now(),
			Price:      800,
			Source:     "market_average",
			Condition:  "new",
			MarketType: "retail",
		},
	}
}
