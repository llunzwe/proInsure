package models

import (
	"errors"
	"time"

	"smartsure/internal/domain/models/policy"
	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Policy represents a comprehensive insurance policy in the database with enterprise-grade features
type Policy struct {
	database.BaseModel

	// Embedded structs for organization
	PolicyIdentification   policy.PolicyIdentification   `gorm:"embedded" json:"identification"`
	PolicyClassification   policy.PolicyClassification   `gorm:"embedded" json:"classification"`
	PolicyCoverageDetails  policy.PolicyCoverageDetails  `gorm:"embedded" json:"coverage_details"`
	PolicyPricing          policy.PolicyPricing          `gorm:"embedded" json:"pricing"`
	PolicyDiscounts        policy.PolicyDiscounts        `gorm:"embedded" json:"discount_info"`
	PolicyLoadings         policy.PolicyLoadings         `gorm:"embedded" json:"loadings"`
	PolicyPaymentInfo      policy.PolicyPaymentInfo      `gorm:"embedded" json:"payment_info"`
	PolicyLifecycle        policy.PolicyLifecycle        `gorm:"embedded" json:"lifecycle"`
	PolicyRiskAssessment   policy.PolicyRiskAssessment   `gorm:"embedded" json:"risk_assessment"`
	PolicyUnderwritingInfo policy.PolicyUnderwritingInfo `gorm:"embedded" json:"underwriting_info"`
	PolicyCompliance       policy.PolicyCompliance       `gorm:"embedded" json:"compliance"`
	PolicyDocumentation    policy.PolicyDocumentation    `gorm:"embedded" json:"documentation"`
	PolicyCommunication    policy.PolicyCommunication    `gorm:"embedded" json:"communication"`
	PolicyAnalytics        policy.PolicyAnalytics        `gorm:"embedded" json:"analytics"`
	PolicyTermsConditions  policy.PolicyTermsConditions  `gorm:"embedded" json:"terms_conditions"`
	PolicyMetadata         policy.PolicyMetadata         `gorm:"embedded" json:"metadata"`
	PolicyAudit            policy.PolicyAudit            `gorm:"embedded" json:"audit"`
	PolicyRelationships    policy.PolicyRelationships    `gorm:"embedded" json:"relationships"`

	// Relationships
	Customer            User                        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Device              *Device                     `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	Product             *policy.Product             `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quote               *policy.Quote               `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	CorporateAccount    *CorporateAccount           `gorm:"foreignKey:CorporateAccountID" json:"corporate_account,omitempty"`
	Agent               *User                       `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Broker              *Partner                    `gorm:"foreignKey:BrokerID" json:"broker,omitempty"`
	Underwriter         *User                       `gorm:"foreignKey:UnderwriterID" json:"underwriter,omitempty"`
	Claims              []Claim                     `gorm:"foreignKey:PolicyID" json:"claims,omitempty"`
	Payments            []Payment                   `gorm:"foreignKey:PolicyID" json:"payments,omitempty"`
	PolicyDocuments     []Document                  `gorm:"foreignKey:PolicyID" json:"documents,omitempty"`
	PolicyModifications []policy.PolicyModification `gorm:"foreignKey:PolicyID" json:"modifications,omitempty"`
	// Supporting Model Relationships
	Bundle             *policy.PolicyBundle                  `gorm:"foreignKey:BundleID" json:"bundle,omitempty"`
	PolicyEndorsements []policy.PolicyEndorsement            `gorm:"foreignKey:PolicyID" json:"policy_endorsements,omitempty"`
	PolicyRiders       []policy.PolicyRider                  `gorm:"foreignKey:PolicyID" json:"policy_riders,omitempty"`
	RenewalInfo        *policy.PolicyRenewal                 `gorm:"foreignKey:OriginalPolicyID" json:"renewal_info,omitempty"`
	PaymentSchedules   []policy.PolicyPaymentSchedule        `gorm:"foreignKey:PolicyID" json:"payment_schedules,omitempty"`
	PolicyLimits       []policy.PolicyLimit                  `gorm:"foreignKey:PolicyID" json:"policy_limits,omitempty"`
	PolicyExclusions   []policy.PolicyExclusion              `gorm:"foreignKey:PolicyID" json:"policy_exclusions,omitempty"`
	Underwriting       *policy.PolicyUnderwriting            `gorm:"foreignKey:PolicyID" json:"underwriting,omitempty"`
	Benefits           []policy.PolicyBenefit                `gorm:"foreignKey:PolicyID" json:"benefits,omitempty"`
	Discounts          []policy.PolicyDiscount               `gorm:"foreignKey:PolicyID" json:"discounts,omitempty"`
	CommunicationPref  *policy.PolicyCommunicationPreference `gorm:"foreignKey:PolicyID" json:"communication_pref,omitempty"`

	// Smartphone-Specific Feature Relationships (Refactored as separate models)
	// One-to-one relationships with CASCADE DELETE - features are deleted when policy is deleted
	PolicyCoverage              *policy.PolicyCoverage              `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"policy_coverage,omitempty"`
	PolicyServiceOptions        *policy.PolicyServiceOptions        `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"service_options,omitempty"`
	PolicyInternationalCoverage *policy.PolicyInternationalCoverage `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"international_coverage,omitempty"`
	PolicyClaimLimits           *policy.PolicyClaimLimits           `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"claim_limits,omitempty"`
	PolicyLoyaltyProgram        *policy.PolicyLoyaltyProgram        `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"loyalty_program,omitempty"`
	PolicySmartFeatures         *policy.PolicySmartFeatures         `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"smart_features,omitempty"`
	PolicyFamilyGroup           *policy.PolicyFamilyGroup           `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"family_group,omitempty"`
	PolicyEnvironmental         *policy.PolicyEnvironmental         `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"environmental,omitempty"`
	PolicyCorporate             *policy.PolicyCorporate             `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"corporate,omitempty"`

	// Enterprise Features - Financial & Risk Management
	PolicyReinsurance *policy.PolicyReinsurance `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"reinsurance,omitempty"`
	PolicyCoinsurance *policy.PolicyCoinsurance `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"coinsurance,omitempty"`
	PolicyReserves    *policy.PolicyReserves    `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"reserves,omitempty"`
	PolicyInvestment  *policy.PolicyInvestment  `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"investment,omitempty"`
	PolicyCommission  *policy.PolicyCommission  `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"commission,omitempty"`

	// Legal & Compliance
	PolicyLegal            *policy.PolicyLegal            `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"legal,omitempty"`
	PolicyRegulatoryFiling *policy.PolicyRegulatoryFiling `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"regulatory_filing,omitempty"`

	// Analytics & Customer Experience
	PolicyPredictiveAnalytics *policy.PolicyPredictiveAnalytics `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"predictive_analytics,omitempty"`
	PolicyCustomerJourney     *policy.PolicyCustomerJourney     `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"customer_journey,omitempty"`

	// Advanced Features
	PolicyIntegrations  *policy.PolicyIntegrations  `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"integrations,omitempty"`
	PolicyTelematics    *policy.PolicyTelematics    `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"telematics,omitempty"`
	PolicyMultiCurrency *policy.PolicyMultiCurrency `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"multi_currency,omitempty"`
	PolicyAutomation    *policy.PolicyAutomation    `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"automation,omitempty"`

	// Risk & Experience Management
	PolicyRiskAccumulation *policy.PolicyRiskAccumulation `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"risk_accumulation,omitempty"`
	PolicyHierarchy        *policy.PolicyHierarchy        `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"hierarchy,omitempty"`
	PolicyExperienceRating *policy.PolicyExperienceRating `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"experience_rating,omitempty"`

	// Business Development
	PolicyOpportunities *policy.PolicyOpportunities `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"opportunities,omitempty"`
	PolicyPrivacy       *policy.PolicyPrivacy       `gorm:"foreignKey:PolicyID;constraint:OnDelete:CASCADE" json:"privacy,omitempty"`
}

// TableName returns the table name for Policy model
func (Policy) TableName() string {
	return "policies"
}

// BeforeCreate handles pre-creation logic with business identifier generation
func (p *Policy) BeforeCreate(tx *gorm.DB) error {
	// Call parent BeforeCreate to handle UUID generation
	if err := p.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	// Generate policy number using BusinessIdentifierGenerator
	if p.PolicyIdentification.PolicyNumber == "" {
		generator := database.NewBusinessIdentifierGenerator()
		p.PolicyIdentification.PolicyNumber = generator.GeneratePolicyNumber()
	}
	// Set default values
	if p.PolicyIdentification.PolicyVersion == 0 {
		p.PolicyIdentification.PolicyVersion = 1
	}
	if p.PolicyPricing.Currency == "" {
		p.PolicyPricing.Currency = policy.CurrencyUSD
	}
	if p.PolicyPaymentInfo.PaymentFrequency == "" {
		p.PolicyPaymentInfo.PaymentFrequency = policy.PaymentFrequencyMonthly
	}
	if p.PolicyLifecycle.Status == "" {
		p.PolicyLifecycle.Status = policy.PolicyStatusPending
	}
	if p.PolicyLifecycle.UnderwritingStatus == "" {
		p.PolicyLifecycle.UnderwritingStatus = policy.UnderwritingStatusPending
	}
	if p.PolicyRiskAssessment.RiskScore == 0 {
		p.PolicyRiskAssessment.RiskScore = 50.0
	}
	return nil
}

// ============================================
// LIFECYCLE MANAGEMENT METHODS
// ============================================

// IsActive checks if the policy is currently active
func (p *Policy) IsActive() bool {
	now := time.Now()
	return p.PolicyLifecycle.Status == policy.PolicyStatusActive &&
		now.After(p.PolicyLifecycle.EffectiveDate) &&
		now.Before(p.PolicyLifecycle.ExpirationDate) &&
		p.PolicyPaymentInfo.PaymentStatus == policy.PaymentStatusPaid
}

// IsExpired checks if the policy has expired
func (p *Policy) IsExpired() bool {
	return time.Now().After(p.PolicyLifecycle.ExpirationDate)
}

// IsSuspended checks if the policy is suspended
func (p *Policy) IsSuspended() bool {
	return p.PolicyLifecycle.Status == policy.PolicyStatusSuspended && p.PolicyLifecycle.SuspensionDate != nil
}

// IsCancelled checks if the policy has been cancelled
func (p *Policy) IsCancelled() bool {
	return p.PolicyLifecycle.Status == policy.PolicyStatusCancelled && p.PolicyLifecycle.CancellationDate != nil
}

// IsPendingRenewal checks if policy needs renewal
func (p *Policy) IsPendingRenewal() bool {
	daysUntilExpiry := p.DaysUntilExpiry()
	return daysUntilExpiry > 0 && daysUntilExpiry <= 30 && p.PolicyPaymentInfo.AutoRenewal
}

// DaysUntilExpiry returns the number of days until policy expires
func (p *Policy) DaysUntilExpiry() int {
	if p.IsExpired() {
		return 0
	}
	return int(time.Until(p.PolicyLifecycle.ExpirationDate).Hours() / 24)
}

// GetPolicyAge returns the age of the policy in days
func (p *Policy) GetPolicyAge() int {
	if p.PolicyLifecycle.EffectiveDate.IsZero() {
		return 0
	}
	return int(time.Since(p.PolicyLifecycle.EffectiveDate).Hours() / 24)
}

// CanFileClaim checks if a claim can be filed against this policy
func (p *Policy) CanFileClaim() (bool, string) {
	if !p.IsActive() {
		return false, "Policy is not active"
	}
	if p.PolicyPaymentInfo.PaymentStatus != policy.PaymentStatusPaid {
		return false, "Policy payment is overdue"
	}
	if p.PolicyCoverageDetails.RemainingLimit.Amount <= 0 {
		return false, "Policy coverage limit exhausted"
	}
	if p.PolicyLifecycle.UnderwritingStatus == policy.UnderwritingStatusRejected {
		return false, "Policy underwriting rejected"
	}
	return true, ""
}

// CanRenew checks if policy can be renewed
func (p *Policy) CanRenew() (bool, string) {
	if p.PolicyLifecycle.Status == policy.PolicyStatusCancelled {
		return false, "Cancelled policies cannot be renewed"
	}
	if p.PolicyRiskAssessment.FraudRiskScore > 80 {
		return false, "High fraud risk detected"
	}
	if p.PolicyRiskAssessment.LossRatio > 150 {
		return false, "Loss ratio too high for renewal"
	}
	if p.PolicyPaymentInfo.OutstandingAmount.Amount > 0 {
		return false, "Outstanding payment required"
	}
	return true, ""
}

// CanModify checks if policy can be modified
func (p *Policy) CanModify() bool {
	return p.PolicyLifecycle.Status == policy.PolicyStatusActive &&
		p.PolicyLifecycle.UnderwritingStatus == policy.UnderwritingStatusApproved &&
		len(p.Claims) == 0 // No active claims
}

// ============================================
// PRICING & CALCULATION METHODS
// ============================================

// CalculateTotalPremium calculates the total premium including all fees
func (p *Policy) CalculateTotalPremium() float64 {
	total := p.PolicyPricing.FinalPremium.Amount + p.PolicyPricing.PremiumTax.Amount +
		p.PolicyPricing.AdminFee.Amount + p.PolicyPricing.ProcessingFee.Amount

	// Apply discounts
	total -= p.PolicyDiscounts.DiscountAmount.Amount
	total -= p.PolicyDiscounts.LoyaltyDiscount.Amount
	total -= p.PolicyDiscounts.BundleDiscount.Amount
	total -= p.PolicyDiscounts.NoClaimsBonus.Amount
	total -= p.PolicyDiscounts.CorporateDiscount.Amount

	// Apply loadings
	total += p.PolicyLoadings.LoadingAmount.Amount

	if total < 0 {
		return 0
	}
	return total
}

// CalculateProRatedPremium calculates pro-rated premium for partial period
func (p *Policy) CalculateProRatedPremium(startDate, endDate time.Time) float64 {
	totalDays := p.PolicyLifecycle.ExpirationDate.Sub(p.PolicyLifecycle.EffectiveDate).Hours() / 24
	periodDays := endDate.Sub(startDate).Hours() / 24

	if totalDays <= 0 {
		return 0
	}

	return p.PolicyPricing.FinalPremium.Amount * (periodDays / totalDays)
}

// CalculateRefund calculates refund amount for cancellation
func (p *Policy) CalculateRefund(cancellationDate time.Time) float64 {
	if cancellationDate.Before(p.PolicyLifecycle.EffectiveDate) {
		return p.PolicyPricing.TotalAmount.Amount // Full refund if cancelled before effective date
	}

	if cancellationDate.After(p.PolicyLifecycle.ExpirationDate) {
		return 0
	}

	usedDays := cancellationDate.Sub(p.PolicyLifecycle.EffectiveDate).Hours() / 24
	totalDays := p.PolicyLifecycle.ExpirationDate.Sub(p.PolicyLifecycle.EffectiveDate).Hours() / 24

	if totalDays <= 0 {
		return 0
	}

	unusedPremium := p.PolicyPricing.FinalPremium.Amount * ((totalDays - usedDays) / totalDays)

	// Apply cancellation penalty (10%)
	penalty := unusedPremium * 0.10
	refund := unusedPremium - penalty - p.PolicyPricing.AdminFee.Amount

	if refund < 0 {
		return 0
	}
	return refund
}

// GetEffectiveDeductible returns the effective deductible amount
func (p *Policy) GetEffectiveDeductible(claimAmount float64) float64 {
	if p.PolicyCoverageDetails.DeductibleType == policy.DeductibleTypePercentage {
		return claimAmount * (p.PolicyCoverageDetails.Deductible.Amount / 100)
	}
	return p.PolicyCoverageDetails.Deductible.Amount
}

// CalculateClaimPayout calculates the payout for a claim
func (p *Policy) CalculateClaimPayout(claimAmount float64) float64 {
	if !p.IsActive() {
		return 0
	}

	// Apply deductible
	effectiveDeductible := p.GetEffectiveDeductible(claimAmount)
	payableAmount := claimAmount - effectiveDeductible

	if payableAmount <= 0 {
		return 0
	}

	// Apply coinsurance
	if p.PolicyCoverageDetails.CoinsurancePercent > 0 {
		payableAmount = payableAmount * (p.PolicyCoverageDetails.CoinsurancePercent / 100)
	}

	// Check coverage limit
	if payableAmount > p.PolicyCoverageDetails.CoverageAmount.Amount {
		payableAmount = p.PolicyCoverageDetails.CoverageAmount.Amount
	}

	// Check remaining limit
	if payableAmount > p.PolicyCoverageDetails.RemainingLimit.Amount {
		payableAmount = p.PolicyCoverageDetails.RemainingLimit.Amount
	}

	// Check out of pocket maximum
	if p.PolicyCoverageDetails.OutOfPocketMax.Amount > 0 && payableAmount > p.PolicyCoverageDetails.OutOfPocketMax.Amount {
		payableAmount = p.PolicyCoverageDetails.OutOfPocketMax.Amount
	}

	return payableAmount
}

// ============================================
// RISK ASSESSMENT METHODS
// ============================================

// CalculateRiskScore calculates the overall risk score
func (p *Policy) CalculateRiskScore() float64 {
	weights := map[string]float64{
		"device":   0.30,
		"user":     0.25,
		"location": 0.15,
		"behavior": 0.15,
		"fraud":    0.15,
	}

	totalScore := p.PolicyRiskAssessment.DeviceRiskScore*weights["device"] +
		p.PolicyRiskAssessment.UserRiskScore*weights["user"] +
		p.PolicyRiskAssessment.LocationRiskScore*weights["location"] +
		p.PolicyRiskAssessment.BehaviorRiskScore*weights["behavior"] +
		p.PolicyRiskAssessment.FraudRiskScore*weights["fraud"]

	// Adjust for claims history
	if p.PolicyRiskAssessment.ClaimFrequency > 2 {
		totalScore += float64(p.PolicyRiskAssessment.ClaimFrequency) * 5
	}

	if totalScore > 100 {
		totalScore = 100
	}

	return totalScore
}

// GetRiskCategory returns the risk category based on score
func (p *Policy) GetRiskCategory() string {
	score := p.CalculateRiskScore()

	switch {
	case score < 25:
		return "low"
	case score < 50:
		return "medium"
	case score < 75:
		return "high"
	default:
		return "very_high"
	}
}

// IsHighRisk checks if the policy is high risk
func (p *Policy) IsHighRisk() bool {
	return p.CalculateRiskScore() > 70 ||
		p.PolicyRiskAssessment.FraudRiskScore > 60 ||
		p.PolicyRiskAssessment.LossRatio > 100 ||
		p.PolicyRiskAssessment.ClaimFrequency > 3
}

// RequiresManualUnderwriting checks if manual underwriting is needed
func (p *Policy) RequiresManualUnderwriting() bool {
	return p.PolicyCoverageDetails.CoverageAmount.Amount > 5000 ||
		p.IsHighRisk() ||
		p.PolicyClassification.PolicyCategory == "corporate" ||
		p.PolicyUnderwritingInfo.RequiresInspection
}

// ============================================
// PAYMENT & BILLING METHODS
// ============================================

// IsPaymentOverdue checks if payment is overdue
func (p *Policy) IsPaymentOverdue() bool {
	if p.PolicyPaymentInfo.PaymentStatus == policy.PaymentStatusPaid {
		return false
	}

	if p.PolicyPaymentInfo.NextBillingDate != nil && time.Now().After(*p.PolicyPaymentInfo.NextBillingDate) {
		gracePeriodEnd := p.PolicyPaymentInfo.NextBillingDate.AddDate(0, 0, p.PolicyPaymentInfo.PaymentGracePeriod)
		return time.Now().After(gracePeriodEnd)
	}

	return p.PolicyPaymentInfo.OutstandingAmount.Amount > 0
}

// GetMonthlyPremium returns the monthly premium amount
func (p *Policy) GetMonthlyPremium() float64 {
	switch p.PolicyPaymentInfo.PaymentFrequency {
	case policy.PaymentFrequencyMonthly:
		return p.PolicyPricing.FinalPremium.Amount
	case policy.PaymentFrequencyQuarterly:
		return p.PolicyPricing.FinalPremium.Amount / 3
	case policy.PaymentFrequencySemiAnnual:
		return p.PolicyPricing.FinalPremium.Amount / 6
	case policy.PaymentFrequencyAnnual:
		return p.PolicyPricing.FinalPremium.Amount / 12
	default:
		return p.PolicyPricing.FinalPremium.Amount
	}
}

// GetNextPaymentAmount returns the next payment amount due
func (p *Policy) GetNextPaymentAmount() float64 {
	baseAmount := p.GetMonthlyPremium()

	// Add outstanding amount if any
	if p.PolicyPaymentInfo.OutstandingAmount.Amount > 0 {
		baseAmount += p.PolicyPaymentInfo.OutstandingAmount.Amount
	}

	return baseAmount
}

// ============================================
// COMPLIANCE & VALIDATION METHODS
// ============================================

// IsCompliant checks if policy meets regulatory requirements
func (p *Policy) IsCompliant() bool {
	return p.PolicyCompliance.ComplianceStatus == "compliant" &&
		p.PolicyCompliance.KYCStatus == "verified" &&
		p.PolicyCompliance.AMLStatus == "cleared" &&
		!p.PolicyCompliance.SanctionsCheck
}

// RequiresKYC checks if KYC verification is required
func (p *Policy) RequiresKYC() bool {
	return p.PolicyCoverageDetails.CoverageAmount.Amount > 10000 ||
		p.PolicyClassification.PolicyCategory == "corporate" ||
		p.PolicyClassification.BusinessLine == "commercial"
}

// NeedsComplianceReview checks if compliance review is needed
func (p *Policy) NeedsComplianceReview() bool {
	return p.PolicyCompliance.ComplianceStatus == "pending" ||
		p.PolicyCompliance.KYCStatus == "pending" ||
		p.PolicyCompliance.AMLStatus == "pending" ||
		p.PolicyRiskAssessment.FraudRiskScore > 70
}

// ============================================
// ANALYTICS & METRICS METHODS
// ============================================

// CalculateLossRatio calculates the loss ratio
func (p *Policy) CalculateLossRatio() float64 {
	if p.PolicyPricing.FinalPremium.Amount <= 0 {
		return 0
	}

	totalClaimsAmount := 0.0
	for _, claim := range p.Claims {
		if claim.ClaimLifecycle.Status == "approved" || claim.ClaimLifecycle.Status == "paid" {
			totalClaimsAmount += claim.ClaimFinancial.ClaimedAmount
		}
	}

	return (totalClaimsAmount / p.PolicyPricing.FinalPremium.Amount) * 100
}

// GetProfitability calculates policy profitability
func (p *Policy) GetProfitability() float64 {
	totalRevenue := p.CalculateTotalPremium()
	totalClaims := 0.0

	for _, claim := range p.Claims {
		if claim.ClaimLifecycle.Status == "paid" {
			totalClaims += claim.ClaimFinancial.ApprovedAmount
		}
	}

	operationalCost := totalRevenue * 0.20 // 20% operational cost assumption
	profit := totalRevenue - totalClaims - operationalCost

	if totalRevenue > 0 {
		return (profit / totalRevenue) * 100
	}
	return 0
}

// GetRetentionProbability calculates the probability of policy renewal
func (p *Policy) GetRetentionProbability() float64 {
	score := 100.0

	// Negative factors
	if p.PolicyRiskAssessment.ClaimFrequency > 2 {
		score -= float64(p.PolicyRiskAssessment.ClaimFrequency) * 10
	}
	if p.PolicyPaymentInfo.PaymentStatus != policy.PaymentStatusPaid {
		score -= 20
	}
	if p.PolicyAnalytics.CustomerSatisfaction < 3 {
		score -= 30
	}

	// Positive factors
	if p.PolicyPaymentInfo.AutoRenewal {
		score += 20
	}
	if p.PolicyDiscounts.NoClaimsBonus.Amount > 0 {
		score += 15
	}
	if p.GetPolicyAge() > 365 {
		score += 10
	}

	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}

	return score
}

// ============================================
// CORPORATE & BUNDLE METHODS
// ============================================

// IsCorporatePolicy checks if this is a corporate policy
func (p *Policy) IsCorporatePolicy() bool {
	if p.PolicyCorporate != nil {
		return p.PolicyCorporate.IsCorporate()
	}
	return p.PolicyIdentification.CorporateAccountID != nil || p.PolicyClassification.PolicyCategory == "corporate"
}

// IsBundledPolicy checks if this is part of a bundle
func (p *Policy) IsBundledPolicy() bool {
	return p.PolicyIdentification.BundleID != nil
}

// IsRenewalPolicy checks if this is a renewal of another policy
func (p *Policy) IsRenewalPolicy() bool {
	return p.PolicyIdentification.ParentPolicyID != nil
}

// IsBYODPolicy checks if this is a BYOD (Bring Your Own Device) policy
func (p *Policy) IsBYODPolicy() bool {
	if p.PolicyCorporate == nil {
		return false
	}
	return p.PolicyCorporate.IsBYOD()
}

// GetCorporateDiscount returns the total corporate discount
func (p *Policy) GetCorporateDiscount() float64 {
	if p.PolicyCorporate == nil {
		return 0
	}
	return p.PolicyCorporate.GetTotalDiscount()
}

// RequiresMDMEnrollment checks if MDM enrollment is required
func (p *Policy) RequiresMDMEnrollment() bool {
	if p.PolicyCorporate == nil {
		return false
	}
	return p.PolicyCorporate.RequiresMDM()
}

// HasCorporatePremiumSupport checks if corporate premium support is included
func (p *Policy) HasCorporatePremiumSupport() bool {
	if p.PolicyCorporate == nil {
		return false
	}
	return p.PolicyCorporate.HasPremiumSupport()
}

// GetTotalDiscounts returns the total discount amount
func (p *Policy) GetTotalDiscounts() float64 {
	total := p.PolicyDiscounts.DiscountAmount.Amount + p.PolicyDiscounts.LoyaltyDiscount.Amount +
		p.PolicyDiscounts.BundleDiscount.Amount + p.PolicyDiscounts.NoClaimsBonus.Amount + p.PolicyDiscounts.CorporateDiscount.Amount

	// Add discounts from PolicyDiscount records
	for _, discount := range p.Discounts {
		if discount.IsValid() {
			total += discount.DiscountAmount
		}
	}

	return total
}

// HasActiveRiders checks if policy has active riders
func (p *Policy) HasActiveRiders() bool {
	for _, rider := range p.PolicyRiders {
		if rider.Status == "active" {
			return true
		}
	}
	return false
}

// GetTotalRiderPremium calculates total premium from all riders
func (p *Policy) GetTotalRiderPremium() float64 {
	total := 0.0
	for _, rider := range p.PolicyRiders {
		if rider.Status == "active" {
			total += rider.Premium
		}
	}
	return total
}

// HasPendingEndorsements checks for pending endorsements
func (p *Policy) HasPendingEndorsements() bool {
	for _, endorsement := range p.PolicyEndorsements {
		if endorsement.Status == "pending" {
			return true
		}
	}
	return false
}

// GetNextPaymentDue returns the next payment due from schedule
func (p *Policy) GetNextPaymentDue() *policy.PolicyPaymentSchedule {
	var nextPayment *policy.PolicyPaymentSchedule

	for i := range p.PaymentSchedules {
		schedule := &p.PaymentSchedules[i]
		if schedule.Status == "pending" {
			if nextPayment == nil || schedule.DueDate.Before(nextPayment.DueDate) {
				nextPayment = schedule
			}
		}
	}

	return nextPayment
}

// GetActiveBenefits returns all currently active benefits
func (p *Policy) GetActiveBenefits() []policy.PolicyBenefit {
	activeBenefits := []policy.PolicyBenefit{}

	for _, benefit := range p.Benefits {
		if benefit.IsAvailable() {
			activeBenefits = append(activeBenefits, benefit)
		}
	}

	return activeBenefits
}

// CheckLimitForPeril checks if coverage limit is available for a peril
func (p *Policy) CheckLimitForPeril(perilType string, amount float64) (bool, float64) {
	for _, limit := range p.PolicyLimits {
		if limit.PerilType == perilType {
			if limit.RemainingAmount >= amount {
				return true, limit.RemainingAmount
			}
			return false, limit.RemainingAmount
		}
	}
	// No specific limit found, use general coverage amount
	return p.PolicyCoverageDetails.RemainingLimit.Amount >= amount, p.PolicyCoverageDetails.RemainingLimit.Amount
}

// IsExcluded checks if a specific condition/peril is excluded
func (p *Policy) IsExcluded(exclusionType string) bool {
	for _, exclusion := range p.PolicyExclusions {
		if exclusion.ExclusionType == exclusionType &&
			exclusion.OverriddenBy == nil {
			return true
		}
	}
	return false
}

// ============================================
// SMARTPHONE-SPECIFIC COVERAGE METHODS
// ============================================

// HasScreenProtection checks if screen damage is covered
func (p *Policy) HasScreenProtection() bool {
	return p.IsActive() && p.PolicyCoverage != nil && p.PolicyCoverage.ScreenProtection
}

// HasWaterDamageProtection checks if water damage is covered
func (p *Policy) HasWaterDamageProtection() bool {
	return p.IsActive() && p.PolicyCoverage != nil && p.PolicyCoverage.WaterDamageProtection
}

// HasTheftProtection checks if theft is covered
func (p *Policy) HasTheftProtection() bool {
	return p.IsActive() && p.PolicyCoverage != nil && p.PolicyCoverage.TheftProtection
}

// HasLossProtection checks if loss is covered
func (p *Policy) HasLossProtection() bool {
	return p.IsActive() && p.PolicyCoverage != nil && p.PolicyCoverage.LossProtection
}

// GetCoverageForDamageType returns if specific damage type is covered
func (p *Policy) GetCoverageForDamageType(damageType string) bool {
	if !p.IsActive() || p.PolicyCoverage == nil {
		return false
	}

	return p.PolicyCoverage.HasProtectionFor(damageType)
}

// CanClaimAccessories checks if accessories can be claimed
func (p *Policy) CanClaimAccessories(claimAmount float64) bool {
	return p.PolicyCoverage != nil && p.PolicyCoverage.CanClaimAccessories(claimAmount)
}

// HasExpressService checks if express services are available
func (p *Policy) HasExpressService() bool {
	return p.PolicyServiceOptions != nil && p.PolicyServiceOptions.HasExpressService()
}

// ============================================
// INTERNATIONAL COVERAGE METHODS
// ============================================

// IsValidInCountry checks if policy is valid in a specific country
func (p *Policy) IsValidInCountry(countryCode string) bool {
	if p.PolicyInternationalCoverage == nil {
		return countryCode == "US" // Default to US only
	}

	return p.PolicyInternationalCoverage.IsValidInCountry(countryCode)
}

// GetRemainingTravelDays returns remaining travel days
func (p *Policy) GetRemainingTravelDays() int {
	if p.PolicyInternationalCoverage == nil {
		return 0
	}
	return p.PolicyInternationalCoverage.GetRemainingTravelDays()
}

// ============================================
// CLAIM FREQUENCY METHODS
// ============================================

// CanFileNewClaim checks if a new claim can be filed
func (p *Policy) CanFileNewClaim() bool {
	if !p.IsActive() {
		return false
	}

	if p.PolicyClaimLimits == nil {
		return true // No limits defined
	}

	canFile, _ := p.PolicyClaimLimits.CanFileNewClaim("", 0)
	return canFile
}

// GetRemainingClaims returns number of claims remaining this year
func (p *Policy) GetRemainingClaims() int {
	if p.PolicyClaimLimits == nil {
		return 0
	}
	return p.PolicyClaimLimits.GetRemainingClaims()
}

// IsInFirstClaimWaitingPeriod checks if in waiting period
func (p *Policy) IsInFirstClaimWaitingPeriod() bool {
	if p.PolicyClaimLimits == nil {
		return false
	}

	if p.PolicyClaimLimits.ClaimsThisYear > 0 {
		return false
	}

	daysSinceStart := int(time.Since(p.PolicyLifecycle.EffectiveDate).Hours() / 24)
	return daysSinceStart < p.PolicyClaimLimits.FirstClaimWaiting
}

// GetClaimFrequency returns claims per year rate
func (p *Policy) GetClaimFrequency() float64 {
	if p.PolicyClaimLimits == nil {
		return 0
	}
	return p.PolicyClaimLimits.GetClaimFrequency()
}

// ============================================
// LOYALTY & REWARDS METHODS
// ============================================

// GetLoyaltyTierBenefits returns benefits based on tier
func (p *Policy) GetLoyaltyTierBenefits() float64 {
	if p.PolicyLoyaltyProgram == nil {
		return 0
	}
	return p.PolicyLoyaltyProgram.GetTierBenefits()
}

// CalculateLoyaltyPoints calculates points earned
func (p *Policy) CalculateLoyaltyPoints() int {
	if p.PolicyLoyaltyProgram == nil {
		return 0
	}

	return p.PolicyLoyaltyProgram.CalculatePointsForActivity("policy_renewal", p.PolicyPricing.FinalPremium.Amount)
}

// HasEarnedReferralBonus checks if referral bonus is earned
func (p *Policy) HasEarnedReferralBonus() bool {
	if p.PolicyLoyaltyProgram == nil {
		return false
	}
	return p.PolicyLoyaltyProgram.HasEarnedReferralBonus()
}

// ============================================
// FAMILY & GROUP METHODS
// ============================================

// IsFamilyPolicy checks if this is a family policy
func (p *Policy) IsFamilyPolicy() bool {
	if p.PolicyFamilyGroup == nil {
		return false
	}
	return p.PolicyFamilyGroup.IsFamilyPolicy()
}

// GetGroupDiscountRate returns group discount rate
func (p *Policy) GetGroupDiscountRate() float64 {
	if p.PolicyFamilyGroup == nil {
		return 0
	}
	return p.PolicyFamilyGroup.GetGroupDiscountRate()
}

// HasSpecialDiscount checks for special discounts
func (p *Policy) HasSpecialDiscount() bool {
	if p.PolicyFamilyGroup == nil {
		return false
	}
	return p.PolicyFamilyGroup.HasSpecialDiscount()
}

// GetSpecialDiscountAmount calculates total special discounts
func (p *Policy) GetSpecialDiscountAmount() float64 {
	if p.PolicyFamilyGroup == nil {
		return 0
	}
	discountRate := p.PolicyFamilyGroup.GetTotalDiscountRate()
	return p.PolicyPricing.FinalPremium.Amount * discountRate
}

// ============================================
// SMART FEATURES METHODS
// ============================================

// HasSmartFeatures checks if smart features are enabled
func (p *Policy) HasSmartFeatures() bool {
	if p.PolicySmartFeatures == nil {
		return false
	}
	return p.PolicySmartFeatures.HasSmartFeatures()
}

// IsEligibleForPreventiveMaintenance checks eligibility
func (p *Policy) IsEligibleForPreventiveMaintenance() bool {
	if p.PolicySmartFeatures == nil {
		return false
	}
	return p.PolicySmartFeatures.IsEligibleForPreventiveMaintenance()
}

// GetSafetyScoreDiscount returns discount based on safety score
func (p *Policy) GetSafetyScoreDiscount() float64 {
	if p.PolicySmartFeatures == nil {
		return 0
	}
	return p.PolicySmartFeatures.GetSafetyScoreDiscount()
}

// ============================================
// ENVIRONMENTAL METHODS
// ============================================

// IsGreenPolicy checks if this is an eco-friendly policy
func (p *Policy) IsGreenPolicy() bool {
	if p.PolicyEnvironmental == nil {
		return false
	}
	return p.PolicyEnvironmental.IsGreenPolicy()
}

// GetGreenDiscount returns environmental discount
func (p *Policy) GetGreenDiscount() float64 {
	if p.PolicyEnvironmental == nil {
		return 0
	}
	return p.PolicyEnvironmental.GetGreenDiscountRate() * p.PolicyPricing.FinalPremium.Amount
}

// CalculateSustainabilityScore calculates sustainability score
func (p *Policy) CalculateSustainabilityScore() float64 {
	if p.PolicyEnvironmental == nil {
		return 0
	}
	return p.PolicyEnvironmental.CalculateSustainabilityScore()
}

// ============================================
// REPLACEMENT SERVICE METHODS
// ============================================

// GetReplacementType returns the type of replacement available
func (p *Policy) GetReplacementType() string {
	if !p.IsActive() || p.PolicyServiceOptions == nil {
		return "none"
	}
	return p.PolicyServiceOptions.ReplacementType
}

// HasLoanerDevice checks if loaner device is available
func (p *Policy) HasLoanerDevice() bool {
	return p.IsActive() && p.PolicyServiceOptions != nil && p.PolicyServiceOptions.HasLoanerDevice()
}

// HasHomeService checks if home service is available
func (p *Policy) HasHomeService() bool {
	return p.IsActive() && p.PolicyServiceOptions != nil && p.PolicyServiceOptions.HasHomeService()
}

// GetServiceLevel returns the service level
func (p *Policy) GetServiceLevel() string {
	if p.PolicyServiceOptions == nil {
		return "standard"
	}
	return p.PolicyServiceOptions.GetServiceLevel()
}

// ============================================
// HELPER METHODS
// ============================================

// generatePolicyNumber is deprecated - use BusinessIdentifierGenerator instead
// Kept for backward compatibility
func generatePolicyNumber() string {
	generator := database.NewBusinessIdentifierGenerator()
	return generator.GeneratePolicyNumber()
}

// ============================================
// POLICY RELATIONSHIP METHODS (Prevent Circular References)
// ============================================

// GetRenewals retrieves all renewal policies for this policy (prevents circular reference)
func (p *Policy) GetRenewals(db *gorm.DB) ([]Policy, error) {
	var renewals []Policy
	err := db.Where("parent_policy_id = ?", p.ID).Find(&renewals).Error
	return renewals, err
}

// GetLatestRenewal retrieves the most recent renewal for this policy
func (p *Policy) GetLatestRenewal(db *gorm.DB) (*Policy, error) {
	var renewal Policy
	err := db.Where("parent_policy_id = ?", p.ID).
		Order("created_at DESC").
		First(&renewal).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &renewal, err
}

// GetParentPolicy retrieves the parent policy if this is a renewal
func (p *Policy) GetParentPolicy(db *gorm.DB) (*Policy, error) {
	if p.PolicyIdentification.ParentPolicyID == nil {
		return nil, nil
	}

	var parent Policy
	err := db.First(&parent, p.PolicyIdentification.ParentPolicyID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &parent, err
}

// GetCustomer returns the customer for this policy, using cached value if available
func (p *Policy) GetCustomer(db *gorm.DB) (*User, error) {
	// If already loaded, return cached value
	if p.Customer.ID != uuid.Nil {
		return &p.Customer, nil
	}

	// Load from database
	var customer User
	err := db.First(&customer, p.PolicyRelationships.CustomerID).Error
	if err != nil {
		return nil, err
	}

	p.Customer = customer
	return &customer, nil
}

// SetStatus updates the policy status with validation
func (p *Policy) SetStatus(status policy.PolicyStatus) error {
	if !status.IsValid() {
		return errors.New("invalid policy status")
	}

	p.PolicyLifecycle.Status = status
	now := time.Now()
	p.PolicyLifecycle.LastModifiedDate = &now

	switch status {
	case policy.PolicyStatusActive:
		p.PolicyLifecycle.ActivationDate = &now
	case policy.PolicyStatusCancelled:
		p.PolicyLifecycle.CancellationDate = &now
	case policy.PolicyStatusSuspended:
		p.PolicyLifecycle.SuspensionDate = &now
	}

	return nil
}

// UpdateRemainingLimit updates the remaining coverage limit after a claim
func (p *Policy) UpdateRemainingLimit(claimAmount float64) {
	p.PolicyCoverageDetails.RemainingLimit.Amount -= claimAmount
	if p.PolicyCoverageDetails.RemainingLimit.Amount < 0 {
		p.PolicyCoverageDetails.RemainingLimit.Amount = 0
	}
}

// ============================================
// VALIDATION METHODS
// ============================================

// Validate performs comprehensive validation on the policy
func (p *Policy) Validate() error {
	// Basic required fields
	if p.PolicyIdentification.PolicyNumber == "" {
		return errors.New("policy number is required")
	}
	if p.PolicyRelationships.CustomerID == uuid.Nil {
		return errors.New("customer ID is required")
	}
	if p.PolicyRelationships.DeviceID == uuid.Nil {
		return errors.New("device ID is required")
	}

	// Coverage validation
	if p.PolicyCoverageDetails.CoverageAmount.Amount <= 0 {
		return errors.New("coverage amount must be positive")
	}
	if p.PolicyPricing.FinalPremium.Amount < 0 {
		return errors.New("premium cannot be negative")
	}
	if p.PolicyCoverageDetails.Deductible.Amount < 0 {
		return errors.New("deductible cannot be negative")
	}
	if p.PolicyCoverageDetails.Deductible.Amount > p.PolicyCoverageDetails.CoverageAmount.Amount {
		return errors.New("deductible cannot exceed coverage amount")
	}

	// Date validation
	if p.PolicyLifecycle.EffectiveDate.IsZero() {
		return errors.New("effective date is required")
	}
	if p.PolicyLifecycle.ExpirationDate.IsZero() {
		return errors.New("expiration date is required")
	}
	if !p.PolicyLifecycle.ExpirationDate.After(p.PolicyLifecycle.EffectiveDate) {
		return errors.New("expiration date must be after effective date")
	}

	// Risk score validation
	if p.PolicyRiskAssessment.RiskScore < 0 || p.PolicyRiskAssessment.RiskScore > 100 {
		return errors.New("risk score must be between 0 and 100")
	}

	// Payment frequency validation
	validFrequencies := []policy.PaymentFrequency{
		policy.PaymentFrequencyMonthly,
		policy.PaymentFrequencyQuarterly,
		policy.PaymentFrequencySemiAnnual,
		policy.PaymentFrequencyAnnual,
	}
	validFreq := false
	for _, freq := range validFrequencies {
		if p.PolicyPaymentInfo.PaymentFrequency == freq {
			validFreq = true
			break
		}
	}
	if !validFreq {
		return errors.New("invalid payment frequency")
	}

	return nil
}

// ShouldSendRenewalNotice checks if renewal notice should be sent
func (p *Policy) ShouldSendRenewalNotice() bool {
	if !p.PolicyPaymentInfo.AutoRenewal || p.PolicyLifecycle.Status != policy.PolicyStatusActive {
		return false
	}

	daysUntilExpiry := p.DaysUntilExpiry()
	// Send notices at 60, 30, 15, and 7 days before expiration
	return daysUntilExpiry == 60 || daysUntilExpiry == 30 ||
		daysUntilExpiry == 15 || daysUntilExpiry == 7
}

// GetCoverageUtilization returns the percentage of coverage used
func (p *Policy) GetCoverageUtilization() float64 {
	if p.PolicyCoverageDetails.CoverageAmount.Amount <= 0 {
		return 0
	}
	used := p.PolicyCoverageDetails.CoverageAmount.Amount - p.PolicyCoverageDetails.RemainingLimit.Amount
	return (used / p.PolicyCoverageDetails.CoverageAmount.Amount) * 100
}

// RequiresRenewalReview checks if policy needs review before renewal
func (p *Policy) RequiresRenewalReview() bool {
	return p.PolicyRiskAssessment.LossRatio > 80 ||
		p.PolicyRiskAssessment.ClaimFrequency > 2 ||
		p.PolicyRiskAssessment.FraudRiskScore > 50 ||
		p.PolicyAnalytics.ChurnRisk > 70
}

// GetPremiumAdjustmentFactor calculates premium adjustment for renewal
func (p *Policy) GetPremiumAdjustmentFactor() float64 {
	factor := 1.0

	// Claim history adjustment
	if p.PolicyRiskAssessment.ClaimFrequency == 0 {
		factor -= 0.10 // 10% discount for no claims
	} else if p.PolicyRiskAssessment.ClaimFrequency > 2 {
		factor += 0.20 // 20% increase for frequent claims
	}

	// Loss ratio adjustment
	if p.PolicyRiskAssessment.LossRatio > 100 {
		factor += 0.25
	} else if p.PolicyRiskAssessment.LossRatio > 75 {
		factor += 0.15
	}

	// Risk score adjustment
	riskCategory := p.GetRiskCategory()
	switch riskCategory {
	case "very_high":
		factor += 0.30
	case "high":
		factor += 0.15
	case "low":
		factor -= 0.05
	}

	// Loyalty adjustment
	if p.GetPolicyAge() > 365*2 { // 2+ years
		factor -= 0.05
	}

	// Ensure factor is within reasonable bounds
	if factor < 0.5 {
		factor = 0.5 // Maximum 50% discount
	}
	if factor > 2.0 {
		factor = 2.0 // Maximum 100% increase
	}

	return factor
}

// ============================================
// SIMPLE GETTER METHODS (Domain Model Pattern)
// ============================================
// Note: Complex business logic should be in /internal/domain/services/policy/

// GetID returns the policy ID
func (p *Policy) GetID() uuid.UUID {
	return p.ID
}

// GetPolicyNumber returns the policy number
func (p *Policy) GetPolicyNumber() string {
	return p.PolicyIdentification.PolicyNumber
}

// GetStatus returns the policy status
func (p *Policy) GetStatus() string {
	return string(p.PolicyLifecycle.Status)
}

// GetEffectiveDate returns the policy effective date
func (p *Policy) GetEffectiveDate() time.Time {
	return p.PolicyLifecycle.EffectiveDate
}

// GetExpirationDate returns the policy expiration date
func (p *Policy) GetExpirationDate() time.Time {
	return p.PolicyLifecycle.ExpirationDate
}

// GetCustomerID returns the customer ID
func (p *Policy) GetCustomerID() uuid.UUID {
	return p.PolicyRelationships.CustomerID
}

// GetDeviceID returns the device ID
func (p *Policy) GetDeviceID() uuid.UUID {
	return p.PolicyRelationships.DeviceID
}

// GetPremiumAmount returns the final premium amount
func (p *Policy) GetPremiumAmount() float64 {
	return p.PolicyPricing.FinalPremium.Amount
}

// GetCoverageAmount returns the coverage amount
func (p *Policy) GetCoverageAmount() float64 {
	return p.PolicyCoverageDetails.CoverageAmount.Amount
}

// GetDeductible returns the deductible amount
func (p *Policy) GetDeductible() float64 {
	return p.PolicyCoverageDetails.Deductible.Amount
}

// GetPaymentStatus returns the payment status
func (p *Policy) GetPaymentStatus() string {
	return string(p.PolicyPaymentInfo.PaymentStatus)
}
