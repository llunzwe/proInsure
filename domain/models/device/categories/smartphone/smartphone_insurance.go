package smartphone

import (
	"fmt"
	"math"
	"time"
	
	"smartsure/internal/domain/models/device/categories/base"
)

// SmartphoneInsuranceCalculator implements insurance calculations for smartphones
type SmartphoneInsuranceCalculator struct {
	base.BaseInsuranceCalculator
	screenRepairWeight  float64
	batteryRepairWeight float64
	waterDamageWeight   float64
}

// NewSmartphoneInsuranceCalculator creates a new smartphone insurance calculator
func NewSmartphoneInsuranceCalculator() *SmartphoneInsuranceCalculator {
	return &SmartphoneInsuranceCalculator{
		BaseInsuranceCalculator: base.BaseInsuranceCalculator{
			BasePremiumRate:    0.08, // 8% of device value annually
			CategoryMultiplier: 1.2,  // Smartphones have 20% higher base rate
			RiskThresholds: map[string]float64{
				"low":      30.0,
				"medium":   50.0,
				"high":     70.0,
				"critical": 90.0,
			},
		},
		screenRepairWeight:  0.35, // Screen damage is 35% of claims
		batteryRepairWeight: 0.15, // Battery issues are 15% of claims
		waterDamageWeight:   0.20, // Water damage is 20% of claims
	}
}

// CalculatePremium calculates insurance premium for a smartphone
func (c *SmartphoneInsuranceCalculator) CalculatePremium(spec base.CategorySpec, coverage base.CoverageLevel) (float64, error) {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0, fmt.Errorf("invalid specification type for smartphone")
	}

	// Get base premium
	basePremium := c.GetBasePremium(phoneSpec.MarketValue)

	// Apply age depreciation
	ageMultiplier := c.CalculateAgeDepreciation(phoneSpec.ReleaseDate, phoneSpec.Depreciation)
	basePremium *= ageMultiplier

	// Apply specification-based multipliers
	specMultiplier := c.calculateSpecMultiplier(phoneSpec)
	basePremium *= specMultiplier

	// Apply coverage level multiplier
	coverageMultiplier := c.getCoverageMultiplier(coverage)
	basePremium *= coverageMultiplier

	// Apply risk factors
	riskFactors := phoneSpec.GetInsuranceFactors()
	riskMultiplier := c.calculateRiskMultiplier(riskFactors)
	basePremium *= riskMultiplier

	// Monthly premium
	monthlyPremium := basePremium / 12

	// Apply minimum and maximum limits
	if monthlyPremium < 5.0 {
		monthlyPremium = 5.0
	}
	if monthlyPremium > 100.0 {
		monthlyPremium = 100.0
	}

	return math.Round(monthlyPremium*100) / 100, nil
}

// CalculateDeductible calculates deductible for a claim type
func (c *SmartphoneInsuranceCalculator) CalculateDeductible(spec base.CategorySpec, claimType base.ClaimType) (float64, error) {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0, fmt.Errorf("invalid specification type for smartphone")
	}

	baseDeductible := phoneSpec.MarketValue * 0.1 // 10% of device value

	// Adjust based on claim type
	switch claimType {
	case base.ClaimTypeScreen:
		baseDeductible *= 0.5 // Screen repairs have lower deductible
	case base.ClaimTypeTheft, base.ClaimTypeLoss:
		baseDeductible *= 1.5 // Theft/loss have higher deductible
	case base.ClaimTypeWater:
		if phoneSpec.WaterResistance == "IP68" || phoneSpec.WaterResistance == "IP67" {
			baseDeductible *= 1.2 // Higher deductible for water damage on resistant phones
		}
	case base.ClaimTypeBattery:
		baseDeductible *= 0.7 // Battery issues have moderate deductible
	}

	// Apply limits
	if baseDeductible < 50.0 {
		baseDeductible = 50.0
	}
	if baseDeductible > 500.0 {
		baseDeductible = 500.0
	}

	return math.Round(baseDeductible), nil
}

// AssessRisk performs risk assessment for a smartphone
func (c *SmartphoneInsuranceCalculator) AssessRisk(spec base.CategorySpec, history base.ClaimHistory) (base.RiskAssessment, error) {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return base.RiskAssessment{}, fmt.Errorf("invalid specification type for smartphone")
	}

	assessment := base.RiskAssessment{
		Factors:         make(map[string]float64),
		Recommendations: []string{},
		ValidUntil:      time.Now().AddDate(0, 6, 0), // Valid for 6 months
	}

	// Calculate base risk score
	score := 0.0

	// Device age risk
	ageInMonths := time.Since(phoneSpec.ReleaseDate).Hours() / 24 / 30
	ageRisk := math.Min(ageInMonths*1.5, 30.0)
	assessment.Factors["age"] = ageRisk
	score += ageRisk

	// Screen size risk (larger screens = higher risk)
	screenRisk := (phoneSpec.ScreenSize - 5.0) * 10.0
	if screenRisk < 0 {
		screenRisk = 0
	}
	assessment.Factors["screen_size"] = screenRisk
	score += screenRisk

	// Build quality risk
	buildRisk := 20.0
	if phoneSpec.WaterResistance != "" {
		buildRisk -= 5.0
	}
	if phoneSpec.ScreenProtection != "" {
		buildRisk -= 5.0
	}
	assessment.Factors["build_quality"] = buildRisk
	score += buildRisk

	// Claim history risk
	if history.TotalClaims > 0 {
		claimRisk := float64(history.TotalClaims) * 15.0
		if history.FraudSuspicion {
			claimRisk *= 2.0
		}
		assessment.Factors["claim_history"] = claimRisk
		score += claimRisk
	}

	// Value risk (expensive phones = higher risk)
	valueRisk := (phoneSpec.MarketValue / 100.0)
	if valueRisk > 20.0 {
		valueRisk = 20.0
	}
	assessment.Factors["device_value"] = valueRisk
	score += valueRisk

	// Foldable risk
	if phoneSpec.IsFoldable() {
		foldableRisk := 25.0
		assessment.Factors["foldable"] = foldableRisk
		score += foldableRisk
		assessment.Recommendations = append(assessment.Recommendations,
			"Foldable device detected - recommend premium protection plan")
	}

	// Set risk level
	assessment.Score = math.Min(score, 100.0)
	if assessment.Score < 30 {
		assessment.Level = "low"
	} else if assessment.Score < 50 {
		assessment.Level = "medium"
		assessment.Recommendations = append(assessment.Recommendations,
			"Consider adding screen protection")
	} else if assessment.Score < 70 {
		assessment.Level = "high"
		assessment.Recommendations = append(assessment.Recommendations,
			"Recommend comprehensive coverage",
			"Suggest protective case and screen protector")
	} else {
		assessment.Level = "critical"
		assessment.RequiresReview = true
		assessment.Recommendations = append(assessment.Recommendations,
			"Manual review required",
			"Consider premium tier coverage only")
	}

	return assessment, nil
}

// GetCoverageOptions returns available coverage options for smartphones
func (c *SmartphoneInsuranceCalculator) GetCoverageOptions(spec base.CategorySpec) []base.CoverageOption {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return []base.CoverageOption{}
	}

	options := []base.CoverageOption{}
	baseValue := phoneSpec.MarketValue

	// Basic coverage
	basic := base.CoverageOption{
		Name:           "Basic Protection",
		Level:          base.CoverageBasic,
		MonthlyPremium: baseValue * 0.005, // 0.5% per month
		AnnualPremium:  baseValue * 0.05,  // 5% per year
		Deductible:     baseValue * 0.15,  // 15% deductible
		MaxCoverage:    baseValue * 0.8,   // 80% max coverage
		CoveredPerils: []base.ClaimType{
			base.ClaimTypeScreen,
			base.ClaimTypeDamage,
		},
		Exclusions: []string{
			"Theft", "Loss", "Water damage", "Cosmetic damage",
		},
		WaitingPeriod: 30,
		ClaimLimit:    1,
		Features: []string{
			"Screen repair", "Accidental damage",
		},
	}
	options = append(options, basic)

	// Standard coverage
	standard := base.CoverageOption{
		Name:           "Standard Protection",
		Level:          base.CoverageStandard,
		MonthlyPremium: baseValue * 0.008,
		AnnualPremium:  baseValue * 0.08,
		Deductible:     baseValue * 0.10,
		MaxCoverage:    baseValue * 0.9,
		CoveredPerils: []base.ClaimType{
			base.ClaimTypeScreen,
			base.ClaimTypeDamage,
			base.ClaimTypeMalfunction,
			base.ClaimTypeBattery,
		},
		Exclusions: []string{
			"Theft", "Loss", "Intentional damage",
		},
		WaitingPeriod: 15,
		ClaimLimit:    2,
		Features: []string{
			"Screen repair", "Accidental damage",
			"Mechanical breakdown", "Battery replacement",
		},
	}
	options = append(options, standard)

	// Premium coverage
	premium := base.CoverageOption{
		Name:           "Premium Protection",
		Level:          base.CoveragePremium,
		MonthlyPremium: baseValue * 0.012,
		AnnualPremium:  baseValue * 0.12,
		Deductible:     baseValue * 0.08,
		MaxCoverage:    baseValue * 0.95,
		CoveredPerils: []base.ClaimType{
			base.ClaimTypeScreen,
			base.ClaimTypeDamage,
			base.ClaimTypeMalfunction,
			base.ClaimTypeBattery,
			base.ClaimTypeWater,
			base.ClaimTypeTheft,
		},
		Exclusions: []string{
			"Loss", "Intentional damage",
		},
		WaitingPeriod: 7,
		ClaimLimit:    3,
		Features: []string{
			"All damage types", "Theft protection",
			"Water damage", "Express replacement",
			"Worldwide coverage",
		},
	}
	options = append(options, premium)

	// Comprehensive coverage (for high-end devices)
	if phoneSpec.IsHighEnd() {
		comprehensive := base.CoverageOption{
			Name:           "Comprehensive Protection",
			Level:          base.CoverageComprehensive,
			MonthlyPremium: baseValue * 0.015,
			AnnualPremium:  baseValue * 0.15,
			Deductible:     baseValue * 0.05,
			MaxCoverage:    baseValue,
			CoveredPerils: []base.ClaimType{
				base.ClaimTypeScreen,
				base.ClaimTypeDamage,
				base.ClaimTypeMalfunction,
				base.ClaimTypeBattery,
				base.ClaimTypeWater,
				base.ClaimTypeTheft,
				base.ClaimTypeLoss,
			},
			Exclusions:    []string{},
			WaitingPeriod: 0,
			ClaimLimit:    -1, // Unlimited
			Features: []string{
				"Full coverage", "Zero waiting period",
				"Unlimited claims", "Loss protection",
				"Premium support", "Same-day replacement",
				"Data recovery service",
			},
		}
		options = append(options, comprehensive)
	}

	return options
}

// Helper methods

func (c *SmartphoneInsuranceCalculator) calculateSpecMultiplier(spec *SmartphoneSpec) float64 {
	multiplier := 1.0

	// Screen size factor
	if spec.ScreenSize > 6.5 {
		multiplier *= 1.1
	}

	// Foldable factor
	if spec.IsFoldable() {
		multiplier *= 1.5
	}

	// 5G factor
	if spec.Has5G() {
		multiplier *= 1.05
	}

	// High-end factor
	if spec.IsHighEnd() {
		multiplier *= 1.15
	}

	// Water resistance reduction
	if spec.WaterResistance == "IP68" {
		multiplier *= 0.95
	}

	return multiplier
}

func (c *SmartphoneInsuranceCalculator) getCoverageMultiplier(coverage base.CoverageLevel) float64 {
	multipliers := map[base.CoverageLevel]float64{
		base.CoverageBasic:         0.8,
		base.CoverageStandard:      1.0,
		base.CoveragePremium:       1.3,
		base.CoverageComprehensive: 1.6,
	}

	if mult, exists := multipliers[coverage]; exists {
		return mult
	}
	return 1.0
}

func (c *SmartphoneInsuranceCalculator) calculateRiskMultiplier(factors map[string]float64) float64 {
	multiplier := 1.0

	for key, value := range factors {
		switch key {
		case "screen_risk":
			multiplier *= (1.0 + value*0.1)
		case "fragility":
			multiplier *= value
		case "repairability":
			multiplier *= value
		case "value_factor":
			multiplier *= value
		}
	}

	return multiplier
}

// ValidateEligibility checks if smartphone is eligible for insurance
func (c *SmartphoneInsuranceCalculator) ValidateEligibility(spec base.CategorySpec) (bool, []string) {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return false, []string{"Invalid specification type"}
	}

	issues := []string{}
	eligible := true

	// Check age
	ageInDays := time.Since(phoneSpec.ReleaseDate).Hours() / 24
	if ageInDays > 1825 { // 5 years
		issues = append(issues, "Device is too old (>5 years)")
		eligible = false
	}

	// Check value
	if phoneSpec.MarketValue < 100 {
		issues = append(issues, "Device value too low (<$100)")
		eligible = false
	}
	if phoneSpec.MarketValue > 5000 {
		issues = append(issues, "Device value exceeds maximum limit (>$5000)")
		eligible = false
	}

	// Check if model is supported
	if phoneSpec.Manufacturer == "" || phoneSpec.Model == "" {
		issues = append(issues, "Device manufacturer or model information missing")
		eligible = false
	}

	return eligible, issues
}

// CalculateMaxCoverage determines maximum insurable value
func (c *SmartphoneInsuranceCalculator) CalculateMaxCoverage(spec base.CategorySpec) float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0
	}

	// Depreciated value
	ageInYears := time.Since(phoneSpec.ReleaseDate).Hours() / 24 / 365
	depreciation := 1.0 - (phoneSpec.Depreciation * ageInYears)
	if depreciation < 0.2 {
		depreciation = 0.2
	}

	return phoneSpec.MarketValue * depreciation
}

// GetExclusions returns policy exclusions
func (c *SmartphoneInsuranceCalculator) GetExclusions() []string {
	return []string{
		"Intentional damage",
		"Cosmetic damage that doesn't affect functionality",
		"Software issues",
		"Damage from unauthorized repairs",
		"Damage from jailbreaking or rooting",
		"Pre-existing conditions",
		"Normal wear and tear",
	}
}

// CalculateClaimPayout calculates the claim payout amount
func (c *SmartphoneInsuranceCalculator) CalculateClaimPayout(spec base.CategorySpec, claim base.ClaimDetails) (float64, error) {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return 0, fmt.Errorf("invalid specification type")
	}

	repairCosts := phoneSpec.GetRepairCosts()
	payout := 0.0

	switch claim.Type {
	case base.ClaimTypeScreen:
		payout = repairCosts["screen"]
	case base.ClaimTypeBattery:
		payout = repairCosts["battery"]
	case base.ClaimTypeTheft, base.ClaimTypeLoss:
		payout = c.CalculateMaxCoverage(spec)
	default:
		payout = claim.EstimatedCost
	}

	// Apply max coverage limit
	maxCoverage := c.CalculateMaxCoverage(spec)
	if payout > maxCoverage {
		payout = maxCoverage
	}

	return payout, nil
}

// GetPremiumFactors returns factors affecting premium
func (c *SmartphoneInsuranceCalculator) GetPremiumFactors(spec base.CategorySpec) map[string]float64 {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return map[string]float64{}
	}

	return phoneSpec.GetInsuranceFactors()
}
