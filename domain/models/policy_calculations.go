package models

import (
	"math"
	"time"

	"smartsure/internal/domain/models/policy"

	"gorm.io/gorm"
)

// ============================================
// POLICY CALCULATION SERVICE
// ============================================

// PolicyCalculator provides comprehensive calculation services for policies
type PolicyCalculator struct {
	policy interface{} // *models.Policy - using interface{} to avoid import cycle
	db     *gorm.DB
}

// NewPolicyCalculator creates a new calculator instance
func NewPolicyCalculator(policy interface{}, db *gorm.DB) *PolicyCalculator {
	return &PolicyCalculator{
		policy: policy,
		db:     db,
	}
}

// getPolicy safely type asserts and returns the policy
func (pc *PolicyCalculator) getPolicy() *Policy {
	if pol, ok := pc.policy.(*Policy); ok {
		return pol
	}
	return nil
}

// safeGetCurrency returns currency from policy or default
func (pc *PolicyCalculator) safeGetCurrency() string {
	if policy := pc.getPolicy(); policy != nil {
		return string(policy.PolicyPricing.Currency)
	}
	return "USD"
}

// safeGetCoverageAmount returns coverage amount from policy or default
func (pc *PolicyCalculator) safeGetCoverageAmount() float64 {
	if policy := pc.getPolicy(); policy != nil {
		return policy.PolicyCoverageDetails.CoverageAmount.Amount
	}
	return 1000.0
}

// ============================================
// PREMIUM CALCULATIONS
// ============================================

// CalculateBasePremium calculates the base premium based on coverage
func (pc *PolicyCalculator) CalculateBasePremium() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(100.0, "USD") // Fallback for invalid type
	}

	coverage := pol.PolicyCoverageDetails.CoverageAmount.Amount

	// Base rate calculation (2-5% of coverage based on type)
	baseRate := pc.getBaseRate()
	basePremium := coverage * baseRate

	// Adjust for policy term
	termMultiplier := pc.getTermMultiplier()
	basePremium *= termMultiplier

	// Minimum premium
	if basePremium < 50 {
		basePremium = 50
	}

	return policy.NewMoney(basePremium, policy.CurrencyCode(string(pol.PolicyPricing.Currency)))
}

// CalculateRiskAdjustedPremium applies risk factors to base premium
func (pc *PolicyCalculator) CalculateRiskAdjustedPremium() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(100.0, "USD")
	}

	basePremium := pol.PolicyPricing.BasePremium.Amount

	// Get risk multiplier based on risk score
	riskMultiplier := pc.getRiskMultiplier()

	// Apply risk adjustment
	adjustedPremium := basePremium * riskMultiplier

	// Apply device-specific factors
	deviceMultiplier := pc.getDeviceMultiplier()
	adjustedPremium *= deviceMultiplier

	// Apply location factors
	locationMultiplier := pc.getLocationMultiplier()
	adjustedPremium *= locationMultiplier

	return policy.NewMoney(adjustedPremium, policy.CurrencyCode(pc.safeGetCurrency()))
}

// CalculateFinalPremium calculates final premium with all discounts and loadings
func (pc *PolicyCalculator) CalculateFinalPremium() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(100.0, "USD")
	}

	// Start with risk-adjusted premium
	premium := pol.PolicyPricing.RiskAdjustedPremium.Amount

	// Apply loadings
	loadings := pc.CalculateTotalLoadings()
	premium += loadings.Amount

	// Apply discounts
	discounts := pc.CalculateTotalDiscounts()
	premium -= discounts.Amount

	// Ensure minimum premium
	minPremium := pol.PolicyCoverageDetails.CoverageAmount.Amount * 0.01 // 1% minimum
	if premium < minPremium {
		premium = minPremium
	}

	return policy.NewMoney(premium, policy.CurrencyCode(string(pol.PolicyPricing.Currency)))
}

// CalculateTotalDiscounts calculates all applicable discounts
func (pc *PolicyCalculator) CalculateTotalDiscounts() policy.Money {
	total := 0.0
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	basePremium := pol.PolicyPricing.BasePremium.Amount

	// Loyalty discount
	if pol.PolicyDiscounts.LoyaltyDiscount.Amount == 0 {
		loyaltyDiscount := pc.calculateLoyaltyDiscount(basePremium)
		total += loyaltyDiscount
	} else {
		total += pol.PolicyDiscounts.LoyaltyDiscount.Amount
	}

	// No claims bonus
	if pol.PolicyDiscounts.NoClaimsBonus.Amount == 0 {
		ncbDiscount := pc.calculateNoClaimsBonus(basePremium)
		total += ncbDiscount
	} else {
		total += pol.PolicyDiscounts.NoClaimsBonus.Amount
	}

	// Bundle discount
	if pol.IsBundledPolicy() {
		bundleDiscount := basePremium * 0.10 // 10% bundle discount
		total += bundleDiscount
	}

	// Corporate discount
	if pol.IsCorporatePolicy() {
		corporateDiscount := basePremium * 0.15 // 15% corporate discount
		total += corporateDiscount
	}

	// Early payment discount
	if pol.PolicyPaymentInfo.PaymentFrequency == policy.PaymentFrequencyAnnual {
		earlyPayDiscount := basePremium * 0.05 // 5% for annual payment
		total += earlyPayDiscount
	}

	// Add any configured discounts
	total += pol.PolicyDiscounts.DiscountAmount.Amount

	// Apply discount percentage if set
	if pol.PolicyDiscounts.DiscountPercentage > 0 {
		percentageDiscount := basePremium * (pol.PolicyDiscounts.DiscountPercentage / 100)
		total += percentageDiscount
	}

	// Cap total discount at 70% of base premium
	maxDiscount := basePremium * 0.70
	if total > maxDiscount {
		total = maxDiscount
	}

	return policy.NewMoney(total, pol.PolicyPricing.Currency)
}

// CalculateTotalLoadings calculates all applicable loadings
func (pc *PolicyCalculator) CalculateTotalLoadings() policy.Money {
	total := 0.0
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	basePremium := pol.PolicyPricing.BasePremium.Amount

	// High risk loading
	if pol.PolicyRiskAssessment.RiskScore > 70 {
		riskLoading := basePremium * 0.15 // 15% high risk loading
		total += riskLoading
	} else if pol.PolicyRiskAssessment.RiskScore > 85 {
		riskLoading := basePremium * 0.25 // 25% very high risk loading
		total += riskLoading
	}

	// Frequent claims loading
	if pol.PolicyRiskAssessment.ClaimFrequency > 3 {
		claimsLoading := basePremium * 0.10 // 10% frequent claims loading
		total += claimsLoading
	}

	// Add any configured loadings
	total += pol.PolicyLoadings.LoadingAmount.Amount

	// Apply loading percentage if set
	if pol.PolicyLoadings.LoadingPercentage > 0 {
		percentageLoading := basePremium * (pol.PolicyLoadings.LoadingPercentage / 100)
		total += percentageLoading
	}

	return policy.NewMoney(total, pol.PolicyPricing.Currency)
}

// ============================================
// TAX AND FEE CALCULATIONS
// ============================================

// CalculatePremiumTax calculates tax on premium
func (pc *PolicyCalculator) CalculatePremiumTax() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	premium := pol.PolicyPricing.FinalPremium.Amount
	taxRate := pc.getTaxRate()
	tax := premium * taxRate
	return policy.NewMoney(tax, pol.PolicyPricing.Currency)
}

// CalculateAdminFee calculates administrative fee
func (pc *PolicyCalculator) CalculateAdminFee() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(25.0, "USD")
	}
	// Base admin fee
	baseFee := 25.0

	// Additional fee for complex policies
	if pol.PolicyClassification.PolicyType == policy.PolicyTypeEnterprise {
		baseFee += 50.0
	} else if pol.PolicyClassification.PolicyType == policy.PolicyTypePlatinum {
		baseFee += 25.0
	}

	// Monthly payment processing fee
	if pol.PolicyPaymentInfo.PaymentFrequency == policy.PaymentFrequencyMonthly {
		baseFee += 5.0
	}

	return policy.NewMoney(baseFee, pol.PolicyPricing.Currency)
}

// CalculateProcessingFee calculates processing fee
func (pc *PolicyCalculator) CalculateProcessingFee() policy.Money {
	// Base processing fee
	fee := 10.0

	// Online discount
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(10.0, "USD")
	}
	if pol.PolicyClassification.Channel == policy.ChannelOnline ||
		pol.PolicyClassification.Channel == policy.ChannelMobile {
		fee *= 0.5 // 50% discount for self-service
	}

	// Rush processing
	daysUntilEffective := int(pol.PolicyLifecycle.EffectiveDate.Sub(time.Now()).Hours() / 24)
	if daysUntilEffective < 2 {
		fee += 25.0 // Rush fee
	}

	return policy.NewMoney(fee, pol.PolicyPricing.Currency)
}

// CalculateTotalAmount calculates the total amount including all fees and taxes
func (pc *PolicyCalculator) CalculateTotalAmount() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	total := pol.PolicyPricing.FinalPremium.Amount
	total += pol.PolicyPricing.PremiumTax.Amount
	total += pol.PolicyPricing.AdminFee.Amount
	total += pol.PolicyPricing.ProcessingFee.Amount

	return policy.NewMoney(total, pol.PolicyPricing.Currency)
}

// ============================================
// REFUND CALCULATIONS
// ============================================

// CalculateCancellationRefund calculates refund for policy cancellation
func (pc *PolicyCalculator) CalculateCancellationRefund(cancellationDate time.Time) policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	// No refund if cancelled after expiration
	if cancellationDate.After(pol.PolicyLifecycle.ExpirationDate) {
		return policy.NewMoney(0, pol.PolicyPricing.Currency)
	}

	// Full refund if cancelled before effective date
	if cancellationDate.Before(pol.PolicyLifecycle.EffectiveDate) {
		return pol.PolicyPricing.TotalAmount
	}

	// Calculate pro-rated refund
	totalDays := pol.PolicyLifecycle.ExpirationDate.Sub(pol.PolicyLifecycle.EffectiveDate).Hours() / 24
	usedDays := cancellationDate.Sub(pol.PolicyLifecycle.EffectiveDate).Hours() / 24
	unusedDays := totalDays - usedDays

	if unusedDays <= 0 {
		return policy.NewMoney(0, pol.PolicyPricing.Currency)
	}

	// Calculate unused premium
	unusedPremium := pol.PolicyPricing.FinalPremium.Amount * (unusedDays / totalDays)

	// Apply cancellation penalty
	penalty := unusedPremium * 0.10 // 10% cancellation penalty

	// Apply short-rate cancellation table
	shortRatePenalty := pc.getShortRatePenalty(usedDays, totalDays)
	penalty = math.Max(penalty, shortRatePenalty)

	// Calculate final refund
	refund := unusedPremium - penalty - pol.PolicyPricing.AdminFee.Amount

	if refund < 0 {
		return policy.NewMoney(0, pol.PolicyPricing.Currency)
	}

	return policy.NewMoney(refund, pol.PolicyPricing.Currency)
}

// CalculateModificationRefund calculates refund for policy modification
func (pc *PolicyCalculator) CalculateModificationRefund(oldCoverage, newCoverage policy.Money, effectiveDate time.Time) policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	if newCoverage.Amount >= oldCoverage.Amount {
		return policy.NewMoney(0, pol.PolicyPricing.Currency) // No refund for coverage increase
	}

	// Calculate premium difference
	coverageReduction := oldCoverage.Amount - newCoverage.Amount
	premiumReduction := (coverageReduction / oldCoverage.Amount) * pol.PolicyPricing.FinalPremium.Amount

	// Pro-rate for remaining term
	remainingDays := pol.PolicyLifecycle.ExpirationDate.Sub(effectiveDate).Hours() / 24
	totalDays := pol.PolicyLifecycle.ExpirationDate.Sub(pol.PolicyLifecycle.EffectiveDate).Hours() / 24

	if remainingDays <= 0 || totalDays <= 0 {
		return policy.NewMoney(0, pol.PolicyPricing.Currency)
	}

	refund := premiumReduction * (remainingDays / totalDays)

	// Apply modification fee
	modificationFee := 25.0
	refund -= modificationFee

	if refund < 0 {
		return policy.NewMoney(0, pol.PolicyPricing.Currency)
	}

	return policy.NewMoney(refund, pol.PolicyPricing.Currency)
}

// ============================================
// CLAIM CALCULATIONS
// ============================================

// CalculateDeductible calculates the applicable deductible for a claim
func (pc *PolicyCalculator) CalculateDeductible(claimAmount float64) policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	if pol.PolicyCoverageDetails.DeductibleType == policy.DeductibleTypePercentage {
		deductible := claimAmount * (pol.PolicyCoverageDetails.Deductible.Amount / 100)
		return policy.NewMoney(deductible, pol.PolicyPricing.Currency)
	}
	return pol.PolicyCoverageDetails.Deductible
}

// CalculateClaimPayout calculates the payout amount for a claim
func (pc *PolicyCalculator) CalculateClaimPayout(claimAmount float64) policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	// Apply deductible
	deductible := pc.CalculateDeductible(claimAmount)
	payableAmount := claimAmount - deductible.Amount

	if payableAmount <= 0 {
		return policy.NewMoney(0, pol.PolicyPricing.Currency)
	}

	// Apply coinsurance
	if pol.PolicyCoverageDetails.CoinsurancePercent > 0 {
		coinsuranceAmount := payableAmount * (1 - pol.PolicyCoverageDetails.CoinsurancePercent/100)
		payableAmount = payableAmount - coinsuranceAmount
	}

	// Check coverage limit
	if payableAmount > pol.PolicyCoverageDetails.CoverageAmount.Amount {
		payableAmount = pol.PolicyCoverageDetails.CoverageAmount.Amount
	}

	// Check remaining limit
	if payableAmount > pol.PolicyCoverageDetails.RemainingLimit.Amount {
		payableAmount = pol.PolicyCoverageDetails.RemainingLimit.Amount
	}

	// Check out of pocket maximum
	if pol.PolicyCoverageDetails.OutOfPocketMax.Amount > 0 {
		// Calculate customer's out of pocket for this claim
		customerPortion := claimAmount - payableAmount
		if customerPortion > pol.PolicyCoverageDetails.OutOfPocketMax.Amount {
			// Increase payout so customer doesn't exceed OOP max
			payableAmount = claimAmount - pol.PolicyCoverageDetails.OutOfPocketMax.Amount
		}
	}

	return policy.NewMoney(payableAmount, pol.PolicyPricing.Currency)
}

// CalculateMaxAnnualPayout calculates maximum payout for the policy year
func (pc *PolicyCalculator) CalculateMaxAnnualPayout() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	// Start with coverage amount
	maxPayout := pol.PolicyCoverageDetails.CoverageAmount.Amount

	// Apply aggregate limit if set
	if pol.PolicyCoverageDetails.AggregateLimit.Amount > 0 {
		maxPayout = math.Min(maxPayout, pol.PolicyCoverageDetails.AggregateLimit.Amount)
	}

	// Check coverage limits
	if pol.PolicyCoverageDetails.CoverageLimits.Annual.Amount > 0 {
		maxPayout = math.Min(maxPayout, pol.PolicyCoverageDetails.CoverageLimits.Annual.Amount)
	}

	return policy.NewMoney(maxPayout, pol.PolicyPricing.Currency)
}

// ============================================
// RENEWAL CALCULATIONS
// ============================================

// CalculateRenewalPremium calculates premium for policy renewal
func (pc *PolicyCalculator) CalculateRenewalPremium() policy.Money {
	pol := pc.getPolicy()
	if pol == nil {
		return policy.NewMoney(0.0, "USD")
	}
	currentPremium := pol.PolicyPricing.FinalPremium.Amount

	// Get adjustment factor based on claims history and risk
	adjustmentFactor := 1.0 // Default adjustment factor
	if pol.PolicyRiskAssessment.ClaimFrequency > 0 {
		adjustmentFactor = 1.1 // 10% increase for claims
	}

	// Apply adjustment
	renewalPremium := currentPremium * adjustmentFactor

	// Apply inflation adjustment (3% annual)
	policyAge := int(pol.PolicyLifecycle.ExpirationDate.Sub(pol.PolicyLifecycle.EffectiveDate).Hours() / 24)
	if policyAge > 365 {
		years := float64(policyAge) / 365
		inflationFactor := math.Pow(1.03, years)
		renewalPremium *= inflationFactor
	}

	// Apply market adjustment
	marketFactor := pc.getMarketAdjustmentFactor()
	renewalPremium *= marketFactor

	// Round to nearest dollar
	renewalPremium = math.Round(renewalPremium)

	return policy.NewMoney(renewalPremium, pol.PolicyPricing.Currency)
}

// ============================================
// HELPER METHODS
// ============================================

// getBaseRate returns the base rate based on policy type
func (pc *PolicyCalculator) getBaseRate() float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 0.025 // Default 2.5%
	}
	rates := map[policy.PolicyType]float64{
		policy.PolicyTypeBasic:         0.02,  // 2%
		policy.PolicyTypeStandard:      0.025, // 2.5%
		policy.PolicyTypeComprehensive: 0.03,  // 3%
		policy.PolicyTypePremium:       0.035, // 3.5%
		policy.PolicyTypePlatinum:      0.04,  // 4%
		policy.PolicyTypeEnterprise:    0.045, // 4.5%
	}

	if rate, exists := rates[pol.PolicyClassification.PolicyType]; exists {
		return rate
	}
	return 0.025 // Default 2.5%
}

// getTermMultiplier returns multiplier based on policy term
func (pc *PolicyCalculator) getTermMultiplier() float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 1.0
	}
	days := pol.PolicyLifecycle.ExpirationDate.Sub(pol.PolicyLifecycle.EffectiveDate).Hours() / 24
	years := days / 365

	if years <= 1 {
		return 1.0
	} else if years <= 2 {
		return 1.9 // Small discount for 2-year term
	} else if years <= 3 {
		return 2.7 // Larger discount for 3-year term
	}
	return 3.5 // Maximum for longer terms
}

// getRiskMultiplier returns multiplier based on risk score
func (pc *PolicyCalculator) getRiskMultiplier() float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 1.0
	}
	score := pol.PolicyRiskAssessment.RiskScore

	if score < 20 {
		return 0.8 // 20% discount
	} else if score < 40 {
		return 0.9 // 10% discount
	} else if score < 60 {
		return 1.0 // Standard rate
	} else if score < 80 {
		return 1.3 // 30% increase
	} else if score < 90 {
		return 1.5 // 50% increase
	}
	return 2.0 // 100% increase for very high risk
}

// getDeviceMultiplier returns multiplier based on device factors
func (pc *PolicyCalculator) getDeviceMultiplier() float64 {
	if pc.db == nil {
		return 1.0
	}

	pol := pc.getPolicy()
	if pol == nil {
		return 1.0
	}
	var device Device
	if err := pc.db.First(&device, pol.PolicyRelationships.DeviceID).Error; err != nil {
		return 1.0 // Default if device not found
	}

	multiplier := 1.0

	// Age factor - calculate from purchase date if available
	if device.DeviceFinancial.PurchaseDate != nil {
		age := int(time.Since(*device.DeviceFinancial.PurchaseDate).Hours() / 24)
		if age < 180 { // Less than 6 months
			multiplier *= 0.9
		} else if age > 365*3 { // More than 3 years
			multiplier *= 1.2
		}
	}

	// Value factor
	currentValue := device.GetCurrentValue()
	if currentValue > 2000 {
		multiplier *= 1.1
	} else if currentValue < 500 {
		multiplier *= 0.9
	}

	// Brand factor
	premiumBrands := []string{"Apple", "Samsung", "Google"}
	brand := device.GetBrand()
	for _, premiumBrand := range premiumBrands {
		if brand == premiumBrand {
			multiplier *= 1.05
			break
		}
	}

	return multiplier
}

// getLocationMultiplier returns multiplier based on location
func (pc *PolicyCalculator) getLocationMultiplier() float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 1.0
	}
	// This would typically use geocoding and crime statistics
	// For now, return based on risk score
	locationRisk := pol.PolicyRiskAssessment.LocationRiskScore

	if locationRisk < 30 {
		return 0.95
	} else if locationRisk < 70 {
		return 1.0
	}
	return 1.15
}

// getTaxRate returns the applicable tax rate
func (pc *PolicyCalculator) getTaxRate() float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 0.06 // Default 6%
	}
	// Tax rate by jurisdiction
	jurisdiction := pol.PolicyCompliance.TaxJurisdiction

	taxRates := map[string]float64{
		"US-CA": 0.0825, // California
		"US-NY": 0.08,   // New York
		"US-TX": 0.0625, // Texas
		"US-FL": 0.06,   // Florida
		"EU":    0.20,   // EU average
		"UK":    0.12,   // UK insurance premium tax
	}

	if rate, exists := taxRates[jurisdiction]; exists {
		return rate
	}
	return 0.06 // Default 6%
}

// getShortRatePenalty calculates short-rate cancellation penalty
func (pc *PolicyCalculator) getShortRatePenalty(usedDays, totalDays float64) float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 0.0
	}
	percentUsed := (usedDays / totalDays) * 100

	// Short-rate penalty table
	if percentUsed < 10 {
		return pol.PolicyPricing.FinalPremium.Amount * 0.10
	} else if percentUsed < 25 {
		return pol.PolicyPricing.FinalPremium.Amount * 0.20
	} else if percentUsed < 50 {
		return pol.PolicyPricing.FinalPremium.Amount * 0.35
	} else if percentUsed < 75 {
		return pol.PolicyPricing.FinalPremium.Amount * 0.50
	}
	return pol.PolicyPricing.FinalPremium.Amount * 0.75
}

// getMarketAdjustmentFactor returns market-based adjustment factor
func (pc *PolicyCalculator) getMarketAdjustmentFactor() float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 1.0
	}
	// This would typically use market data and competitor analysis
	// For now, return a simple factor based on policy type
	switch pol.PolicyClassification.PolicyType {
	case policy.PolicyTypeBasic:
		return 0.98 // Competitive pricing for basic
	case policy.PolicyTypeEnterprise:
		return 1.05 // Premium pricing for enterprise
	default:
		return 1.0
	}
}

// calculateLoyaltyDiscount calculates loyalty discount based on customer history
func (pc *PolicyCalculator) calculateLoyaltyDiscount(basePremium float64) float64 {
	if pc.db == nil {
		return 0
	}

	// Count customer's active policies
	pol := pc.getPolicy()
	if pol == nil {
		return 0
	}
	var policyCount int64
	pc.db.Model(&Policy{}).
		Where("customer_id = ? AND status = ?",
			pol.PolicyRelationships.CustomerID,
			"active").
		Count(&policyCount)

	discount := 0.0
	if policyCount >= 3 {
		discount = basePremium * 0.10 // 10% for 3+ policies
	} else if policyCount >= 2 {
		discount = basePremium * 0.05 // 5% for 2 policies
	}

	// Additional discount for long-term customers
	// Simplified implementation - in production this would query customer join date
	// For now, apply basic loyalty discount based on policy count
	if policyCount >= 5 {
		discount += basePremium * 0.15 // 15%
	} else if policyCount >= 3 {
		discount += basePremium * 0.10 // 10%
	}

	return discount
}

// calculateNoClaimsBonus calculates bonus for claim-free history
func (pc *PolicyCalculator) calculateNoClaimsBonus(basePremium float64) float64 {
	pol := pc.getPolicy()
	if pol == nil {
		return 0
	}
	if pol.PolicyRiskAssessment.ClaimFrequency > 0 {
		return 0 // No bonus if claims exist
	}

	// Calculate based on policy age
	policyAge := int(pol.PolicyLifecycle.ExpirationDate.Sub(pol.PolicyLifecycle.EffectiveDate).Hours() / 24)
	if policyAge < 365 {
		return 0 // No bonus in first year
	}

	years := float64(policyAge) / 365.0
	bonusRate := years * 0.05 // 5% per year
	if bonusRate > 0.30 {
		bonusRate = 0.30 // Cap at 30% maximum
	}

	return basePremium * bonusRate
}
