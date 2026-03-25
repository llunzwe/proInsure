package smartwatch

import (
	"errors"
	"math"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models/device/categories/base"
)

// SmartwatchInsuranceCalculator provides smartwatch-specific insurance calculations
// It leverages existing device models for comprehensive risk assessment
type SmartwatchInsuranceCalculator struct{}

// NewSmartwatchInsuranceCalculator creates a new smartwatch insurance calculator
func NewSmartwatchInsuranceCalculator() *SmartwatchInsuranceCalculator {
	return &SmartwatchInsuranceCalculator{}
}

// CalculatePremium calculates monthly premium for smartwatch
func (c *SmartwatchInsuranceCalculator) CalculatePremium(spec base.CategorySpec, coverage base.CoverageLevel) (float64, error) {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 0, errors.New("invalid smartwatch specification")
	}

	// Base premium calculation
	basePremium := swSpec.MarketValue * 0.03 // 3% of value per month

	// Adjust for coverage level
	coverageMultiplier := 1.0
	switch coverage {
	case base.CoverageBasic:
		coverageMultiplier = 0.7
	case base.CoverageStandard:
		coverageMultiplier = 1.0
	case base.CoveragePremium:
		coverageMultiplier = 1.5
	case base.CoverageComprehensive:
		coverageMultiplier = 2.0
	}

	// Get insurance factors
	factors := swSpec.GetInsuranceFactors()

	// Calculate risk adjustments
	riskMultiplier := 1.0
	riskMultiplier *= (1.0 + factors["screen_risk"]*0.2)
	riskMultiplier *= (1.0 + factors["water_damage_risk"]*0.3)
	riskMultiplier *= (1.0 + factors["battery_risk"]*0.1)
	riskMultiplier *= (1.0 - factors["material_quality"]*0.1) // Better materials reduce premium

	// Health sensor premium (medical-grade sensors increase value)
	if healthScore := swSpec.GetHealthFeatureScore(); healthScore > 3 {
		riskMultiplier *= 1.15 // Premium devices with health features
	}

	// Activity tracking increases risk (more outdoor use)
	if len(swSpec.ActivityTracking) > 5 {
		riskMultiplier *= 1.1
	}

	// Cellular models have higher premium (more complex, higher value)
	if swSpec.Cellular {
		riskMultiplier *= 1.2
	}

	premium := basePremium * coverageMultiplier * riskMultiplier

	// Minimum and maximum premiums
	if premium < 5.0 {
		premium = 5.0
	}
	if premium > 50.0 {
		premium = 50.0
	}

	return premium, nil
}

// CalculateDeductible calculates deductible amount for smartwatch
func (c *SmartwatchInsuranceCalculator) CalculateDeductible(spec base.CategorySpec, claimType base.ClaimType) (float64, error) {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 0, errors.New("invalid smartwatch specification")
	}

	baseDeductible := swSpec.MarketValue * 0.1 // 10% of value

	// Adjust by claim type
	switch claimType {
	case base.ClaimTypeDamage:
		return baseDeductible, nil
	case base.ClaimTypeTheft:
		return baseDeductible * 1.5, nil
	case base.ClaimTypeLoss:
		return baseDeductible * 2.0, nil
	case base.ClaimTypeMalfunction:
		return baseDeductible * 0.5, nil
	case base.ClaimTypeWater:
		// Lower deductible if water-resistant
		if swSpec.WaterResistance == "10ATM" || swSpec.WaterResistance == "5ATM" {
			return baseDeductible * 0.7, nil
		}
		return baseDeductible * 1.2, nil
	default:
		return baseDeductible, nil
	}
}

// AssessRisk performs risk assessment for smartwatch
func (c *SmartwatchInsuranceCalculator) AssessRisk(spec base.CategorySpec, history base.ClaimHistory) (base.RiskAssessment, error) {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return base.RiskAssessment{}, errors.New("invalid smartwatch specification")
	}

	assessment := base.RiskAssessment{
		Score:      50.0, // Base score
		Level:      "medium",
		Factors:    make(map[string]float64),
		ValidUntil: time.Now().AddDate(0, 6, 0), // Valid for 6 months
	}

	// Device age factor
	deviceAge := time.Since(swSpec.PurchaseDate).Hours() / 24 / 365
	if deviceAge > 2 {
		assessment.Score += 15
		assessment.Factors["age"] = deviceAge
	}

	// Water resistance factor
	waterRisk := 20.0
	switch swSpec.WaterResistance {
	case "10ATM", "100m":
		waterRisk = 5.0
	case "5ATM", "50m", "IP68":
		waterRisk = 10.0
	case "3ATM", "30m", "IP67":
		waterRisk = 15.0
	}
	assessment.Score += waterRisk
	assessment.Factors["water_risk"] = waterRisk

	// Activity usage increases risk
	if len(swSpec.ActivityTracking) > 3 {
		assessment.Score += 10
		assessment.Factors["activity_usage"] = float64(len(swSpec.ActivityTracking))
	}

	// Claims history impact
	claimCount := history.TotalClaims
	if claimCount > 0 {
		assessment.Score += float64(claimCount) * 10
		assessment.Factors["claim_history"] = float64(claimCount)
	}

	// Material quality reduces risk
	if swSpec.CaseMaterial == "titanium" || swSpec.CaseMaterial == "steel" {
		assessment.Score -= 10
		assessment.Factors["material_quality"] = -10
	}

	// Set risk level
	switch {
	case assessment.Score < 30:
		assessment.Level = "low"
	case assessment.Score < 60:
		assessment.Level = "medium"
	case assessment.Score < 80:
		assessment.Level = "high"
	default:
		assessment.Level = "very_high"
		assessment.RequiresReview = true
	}

	// Add recommendations
	if assessment.Score > 60 {
		assessment.Recommendations = append(assessment.Recommendations,
			"Consider protective case for active use",
			"Regular software updates recommended",
			"Avoid extreme sports without proper protection")
	}

	if waterRisk > 15 {
		assessment.Recommendations = append(assessment.Recommendations,
			"Limited water resistance - avoid swimming/showering")
	}

	return assessment, nil
}

// GetCoverageOptions returns available coverage options for smartwatch
func (c *SmartwatchInsuranceCalculator) GetCoverageOptions(spec base.CategorySpec) []base.CoverageOption {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return []base.CoverageOption{}
	}

	options := []base.CoverageOption{
		{
			ID:             uuid.New(),
			Level:          base.CoverageBasic,
			Name:           "Basic Protection",
			MonthlyPremium: 0, // Will be calculated
			AnnualPremium:  0,
			Deductible:     swSpec.MarketValue * 0.15,
			MaxCoverage:    swSpec.MarketValue * 0.8,
			CoveredPerils: []base.ClaimType{
				base.ClaimTypeDamage,
				base.ClaimTypeMalfunction,
				base.ClaimTypeBattery,
			},
			Exclusions: []string{
				"Loss",
				"Theft",
				"Cosmetic damage",
			},
			WaitingPeriod: 30,
			ClaimLimit:    2,
			Features: []string{
				"Accidental damage coverage",
				"Mechanical breakdown protection",
				"Battery failure coverage",
			},
		},
		{
			ID:             uuid.New(),
			Level:          base.CoverageStandard,
			Name:           "Standard Protection",
			MonthlyPremium: 0,
			AnnualPremium:  0,
			Deductible:     swSpec.MarketValue * 0.1,
			MaxCoverage:    swSpec.MarketValue,
			CoveredPerils: []base.ClaimType{
				base.ClaimTypeDamage,
				base.ClaimTypeMalfunction,
				base.ClaimTypeBattery,
				base.ClaimTypeTheft,
				base.ClaimTypeWater,
			},
			Exclusions: []string{
				"Loss",
				"Intentional damage",
			},
			WaitingPeriod: 14,
			ClaimLimit:    3,
			Features: []string{
				"Comprehensive damage coverage",
				"Theft protection",
				"Water damage coverage (if rated)",
			},
		},
		{
			ID:             uuid.New(),
			Level:          base.CoveragePremium,
			Name:           "Premium Protection",
			MonthlyPremium: 0,
			AnnualPremium:  0,
			Deductible:     swSpec.MarketValue * 0.05,
			MaxCoverage:    swSpec.MarketValue * 1.1,
			CoveredPerils: []base.ClaimType{
				base.ClaimTypeDamage,
				base.ClaimTypeMalfunction,
				base.ClaimTypeBattery,
				base.ClaimTypeTheft,
				base.ClaimTypeWater,
				base.ClaimTypeLoss,
				base.ClaimTypeScreen,
			},
			Exclusions: []string{
				"Intentional damage",
			},
			WaitingPeriod: 7,
			ClaimLimit:    4,
			Features: []string{
				"Full coverage including loss",
				"Express replacement service",
				"Band damage coverage",
				"Worldwide protection",
			},
		},
	}

	// Calculate premiums for each option
	for i := range options {
		if premium, err := c.CalculatePremium(spec, options[i].Level); err == nil {
			options[i].MonthlyPremium = premium
			options[i].AnnualPremium = premium * 11 // 1 month discount for annual payment
		}
	}

	return options
}

// ValidateEligibility checks if smartwatch is eligible for insurance
func (c *SmartwatchInsuranceCalculator) ValidateEligibility(spec base.CategorySpec) (bool, []string) {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return false, []string{"Invalid smartwatch specification"}
	}

	issues := []string{}
	eligible := true

	// Check age
	deviceAge := time.Since(swSpec.PurchaseDate).Hours() / 24 / 365
	if deviceAge > 3 {
		issues = append(issues, "Device is over 3 years old")
		eligible = false
	}

	// Check value
	if swSpec.MarketValue < 50 {
		issues = append(issues, "Device value too low for insurance")
		eligible = false
	}

	if swSpec.MarketValue > 2000 {
		issues = append(issues, "Device requires special underwriting")
	}

	// Check condition
	if swSpec.Condition != "" && swSpec.Condition != "excellent" && swSpec.Condition != "good" {
		issues = append(issues, "Device must be in good or excellent condition")
		eligible = false
	}

	// OS support check
	if swSpec.OS == "proprietary" || swSpec.OS == "discontinued" {
		issues = append(issues, "Operating system not supported")
		eligible = false
	}

	return eligible, issues
}

// CalculateMaxCoverage determines maximum insurable value for smartwatch
func (c *SmartwatchInsuranceCalculator) CalculateMaxCoverage(spec base.CategorySpec) float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 0
	}

	// Base max coverage is market value
	maxCoverage := swSpec.MarketValue

	// Premium models can have higher coverage
	if swSpec.MarketValue > 500 {
		maxCoverage = swSpec.MarketValue * 1.1 // 110% for premium models
	}

	// Medical-grade devices get extra coverage
	if healthScore := swSpec.GetHealthFeatureScore(); healthScore > 4 {
		maxCoverage = swSpec.MarketValue * 1.15
	}

	// Cap at reasonable maximum
	if maxCoverage > 2000 {
		maxCoverage = 2000
	}

	return maxCoverage
}

// GetExclusions returns policy exclusions for smartwatch category
func (c *SmartwatchInsuranceCalculator) GetExclusions() []string {
	return []string{
		"Intentional damage or destruction",
		"Normal wear and tear",
		"Cosmetic damage not affecting functionality",
		"Software issues not caused by hardware failure",
		"Band or strap damage (unless specifically covered)",
		"Loss of data or applications",
		"Damage from unauthorized repairs",
		"Damage from jailbreaking or rooting",
		"Pre-existing conditions",
		"Damage during professional sports",
		"Nuclear hazards, war, or terrorism",
		"Damage from extreme sports (unless covered)",
		"Mysterious disappearance",
	}
}

// CalculateClaimPayout calculates the claim payout amount for smartwatch
func (c *SmartwatchInsuranceCalculator) CalculateClaimPayout(spec base.CategorySpec, claim base.ClaimDetails) (float64, error) {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return 0, errors.New("invalid smartwatch specification")
	}

	// Check if pre-existing condition
	if claim.IsPreExisting {
		return 0, errors.New("pre-existing conditions not covered")
	}

	// Base payout is estimated cost
	payout := claim.EstimatedCost

	// Apply depreciation based on device age
	deviceAge := time.Since(swSpec.PurchaseDate).Hours() / 24 / 365
	depreciationRate := 0.2 * deviceAge                       // 20% per year
	depreciationFactor := math.Max(0.3, 1.0-depreciationRate) // Minimum 30% value

	// Adjust payout based on claim type
	switch claim.Type {
	case base.ClaimTypeTheft, base.ClaimTypeLoss:
		// Total loss - pay depreciated value
		payout = swSpec.MarketValue * depreciationFactor

	case base.ClaimTypeDamage:
		// Repair cost, capped at depreciated value
		payout = math.Min(payout, swSpec.MarketValue*depreciationFactor)

	case base.ClaimTypeScreen:
		// Screen repair is common, standard rate
		if payout > swSpec.MarketValue*0.3 {
			payout = swSpec.MarketValue * 0.3
		}

	case base.ClaimTypeBattery:
		// Battery replacement, capped
		if payout > swSpec.MarketValue*0.2 {
			payout = swSpec.MarketValue * 0.2
		}

	case base.ClaimTypeWater:
		// Water damage depends on water resistance
		if swSpec.WaterResistance == "10ATM" || swSpec.WaterResistance == "5ATM" {
			// Should be water resistant, reduce payout
			payout *= 0.7
		}

	case base.ClaimTypeMalfunction:
		// Mechanical failure, check warranty
		if deviceAge < 1 {
			// Still under warranty, reduce payout
			payout *= 0.5
		}
	}

	// Apply deductible (this should be done at claims processing, but calculate here for estimate)
	deductible, _ := c.CalculateDeductible(spec, claim.Type)
	payout = math.Max(0, payout-deductible)

	// Cap at max coverage
	maxCoverage := c.CalculateMaxCoverage(spec)
	if payout > maxCoverage {
		payout = maxCoverage
	}

	return payout, nil
}

// GetPremiumFactors returns factors affecting premium calculation for smartwatch
func (c *SmartwatchInsuranceCalculator) GetPremiumFactors(spec base.CategorySpec) map[string]float64 {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return map[string]float64{}
	}

	factors := make(map[string]float64)

	// Base factors
	factors["base_rate"] = 0.03 // 3% of value
	factors["market_value"] = swSpec.MarketValue

	// Age factor
	deviceAge := time.Since(swSpec.PurchaseDate).Hours() / 24 / 365
	factors["age_years"] = deviceAge
	factors["age_multiplier"] = 1.0 + (deviceAge * 0.2) // 20% increase per year

	// Water resistance factor
	waterMultiplier := 1.0
	switch swSpec.WaterResistance {
	case "10ATM", "100m":
		waterMultiplier = 0.9 // Lower risk
	case "5ATM", "50m", "IP68":
		waterMultiplier = 0.95
	case "3ATM", "30m", "IP67":
		waterMultiplier = 1.0
	default:
		waterMultiplier = 1.2 // Higher risk
	}
	factors["water_resistance_multiplier"] = waterMultiplier

	// Material quality factor
	materialMultiplier := 1.0
	switch swSpec.CaseMaterial {
	case "titanium":
		materialMultiplier = 0.85
	case "steel", "stainless steel":
		materialMultiplier = 0.9
	case "aluminum":
		materialMultiplier = 1.0
	case "plastic", "resin":
		materialMultiplier = 1.1
	}
	factors["material_multiplier"] = materialMultiplier

	// Health features factor (more sensors = higher value = higher premium)
	healthScore := swSpec.GetHealthFeatureScore()
	factors["health_feature_score"] = float64(healthScore)
	if healthScore > 3 {
		factors["health_premium_multiplier"] = 1.15
	} else {
		factors["health_premium_multiplier"] = 1.0
	}

	// Activity usage factor
	activityCount := float64(len(swSpec.ActivityTracking))
	factors["activity_count"] = activityCount
	if activityCount > 5 {
		factors["activity_risk_multiplier"] = 1.1
	} else {
		factors["activity_risk_multiplier"] = 1.0
	}

	// Connectivity factor
	if swSpec.Cellular {
		factors["cellular_multiplier"] = 1.2
	} else {
		factors["cellular_multiplier"] = 1.0
	}

	// Display size factor (larger = more fragile)
	if swSpec.DisplaySize > 1.5 {
		factors["display_size_multiplier"] = 1.1
	} else {
		factors["display_size_multiplier"] = 1.0
	}

	// Battery capacity factor
	if swSpec.BatteryCapacity > 500 {
		factors["battery_risk_multiplier"] = 1.05
	} else {
		factors["battery_risk_multiplier"] = 1.0
	}

	return factors
}
