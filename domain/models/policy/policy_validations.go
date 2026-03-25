package policy

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================
// POLICY VALIDATION SERVICE
// ============================================

// PolicyValidator provides comprehensive validation for policies
type PolicyValidator struct {
	db     *gorm.DB
	policy interface{} // *models.Policy - using interface{} to avoid import cycle
}

// NewPolicyValidator creates a new validator instance
func NewPolicyValidator(db *gorm.DB, policy interface{}) *PolicyValidator {
	return &PolicyValidator{
		db:     db,
		policy: policy,
	}
}

// ============================================
// CORE VALIDATIONS
// ============================================

// ValidateComprehensive performs all validation checks
func (pv *PolicyValidator) ValidateComprehensive() error {
	// Required field validations
	if err := pv.ValidateRequiredFields(); err != nil {
		return err
	}

	// Business rule validations
	if err := pv.ValidateBusinessRules(); err != nil {
		return err
	}

	// Financial validations
	if err := pv.ValidateFinancials(); err != nil {
		return err
	}

	// Date validations
	if err := pv.ValidateDates(); err != nil {
		return err
	}

	// Coverage validations
	if err := pv.ValidateCoverage(); err != nil {
		return err
	}

	// Risk validations
	if err := pv.ValidateRisk(); err != nil {
		return err
	}

	// Compliance validations
	if err := pv.ValidateCompliance(); err != nil {
		return err
	}

	return nil
}

// ValidateRequiredFields checks all required fields are present
func (pv *PolicyValidator) ValidateRequiredFields() error {
	// Policy identification
	if pv.getNestedStringField("PolicyIdentification", "PolicyNumber") == "" {
		return errors.New("policy number is required")
	}

	// Customer and device
	if pv.getNestedUUIDField("PolicyRelationships", "CustomerID") == uuid.Nil {
		return errors.New("customer ID is required")
	}
	if pv.getNestedUUIDField("PolicyRelationships", "DeviceID") == uuid.Nil {
		return errors.New("device ID is required")
	}

	// Policy type and classification
	if pv.getNestedStringField("PolicyClassification", "PolicyType") == "" {
		return errors.New("policy type is required")
	}
	if pv.getNestedStringField("PolicyClassification", "CoverageType") == "" {
		return errors.New("coverage type is required")
	}

	// Dates - using reflection to check if zero
	effectiveDateVal := pv.getPolicyField("PolicyLifecycle", "EffectiveDate")
	if !effectiveDateVal.IsValid() || effectiveDateVal.IsZero() {
		return errors.New("effective date is required")
	}
	expirationDateVal := pv.getPolicyField("PolicyLifecycle", "ExpirationDate")
	if !expirationDateVal.IsValid() || expirationDateVal.IsZero() {
		return errors.New("expiration date is required")
	}

	// Financial
	if pv.getNestedFloat64Field("PolicyPricing", "FinalPremium", "Amount") <= 0 {
		return errors.New("premium must be positive")
	}
	if pv.getNestedFloat64Field("PolicyCoverageDetails", "CoverageAmount", "Amount") <= 0 {
		return errors.New("coverage amount must be positive")
	}

	return nil
}

// ============================================
// BUSINESS RULE VALIDATIONS
// ============================================

// ValidateBusinessRules checks business logic constraints
func (pv *PolicyValidator) ValidateBusinessRules() error {
	// Check device eligibility
	if err := pv.ValidateDeviceEligibility(); err != nil {
		return err
	}

	// Check customer eligibility
	if err := pv.ValidateCustomerEligibility(); err != nil {
		return err
	}

	// Validate policy limits
	if err := pv.ValidatePolicyLimits(); err != nil {
		return err
	}

	// Validate underwriting rules
	if err := pv.ValidateUnderwritingRules(); err != nil {
		return err
	}

	return nil
}

// ValidateDeviceEligibility checks if device can be insured
func (pv *PolicyValidator) ValidateDeviceEligibility() error {
	if pv.db == nil {
		return nil // Skip DB validations in unit tests
	}

	// Use reflection to access device data without direct import
	deviceID := pv.getNestedUUIDField("PolicyRelationships", "DeviceID")
	if deviceID == uuid.Nil {
		return errors.New("device ID not found in policy")
	}

	// Query device using raw SQL to avoid import cycle
	var deviceStatus string
	var currentValue float64
	var purchaseDate *time.Time
	var statusInfo map[string]interface{}

	if err := pv.db.Raw(`
		SELECT 
			status,
			current_value_amount as current_value,
			purchase_date,
			status_info
		FROM devices 
		WHERE id = ?
	`, deviceID).Scan(&struct {
		Status       string
		CurrentValue float64
		PurchaseDate *time.Time
		StatusInfo   map[string]interface{}
	}{
		Status:       deviceStatus,
		CurrentValue: currentValue,
		PurchaseDate: purchaseDate,
		StatusInfo:   statusInfo,
	}).Error; err != nil {
		return fmt.Errorf("device not found: %v", err)
	}

	// Check device age
	if purchaseDate != nil {
		deviceAge := int(time.Since(*purchaseDate).Hours() / 24)
		if deviceAge > 365*5 { // 5 years
			return errors.New("device is too old for insurance (max 5 years)")
		}
	}

	// Check device value
	if currentValue < 100 {
		return errors.New("device value too low for insurance (min $100)")
	}
	if currentValue > 50000 {
		return errors.New("device value too high for standard insurance (max $50,000)")
	}

	// Check device status
	if deviceStatus == "stolen" || deviceStatus == "blacklisted" {
		return fmt.Errorf("cannot insure device with status: %s", deviceStatus)
	}

	return nil
}

// ValidateCustomerEligibility checks if customer ID is valid
// Full customer eligibility validation should be done at the service layer
func (pv *PolicyValidator) ValidateCustomerEligibility() error {
	// Basic validation - ensure customer ID is provided
	if pv.getNestedUUIDField("PolicyRelationships", "CustomerID") == uuid.Nil {
		return errors.New("customer ID is required")
	}

	// Additional customer validation (status, age, fraud score, etc.)
	// should be performed at the service layer to avoid circular imports

	return nil
}

// ValidatePolicyLimits checks policy-specific limits
func (pv *PolicyValidator) ValidatePolicyLimits() error {
	// Coverage amount limits by policy type
	maxCoverage := map[PolicyType]float64{
		PolicyTypeBasic:         5000,
		PolicyTypeStandard:      15000,
		PolicyTypeComprehensive: 30000,
		PolicyTypePremium:       50000,
		PolicyTypePlatinum:      100000,
		PolicyTypeEnterprise:    500000,
	}

	policyType := PolicyType(pv.getNestedStringField("PolicyClassification", "PolicyType"))
	if limit, exists := maxCoverage[policyType]; exists {
		coverageAmount := pv.getNestedFloat64Field("PolicyCoverageDetails", "CoverageAmount", "Amount")
		if coverageAmount > limit {
			return fmt.Errorf("coverage amount exceeds limit for %s policy (max $%.2f)",
				policyType, limit)
		}
	}

	// Deductible limits
	coverageAmount := pv.getNestedFloat64Field("PolicyCoverageDetails", "CoverageAmount", "Amount")
	minDeductible := coverageAmount * 0.01 // 1% minimum
	maxDeductible := coverageAmount * 0.25 // 25% maximum

	deductibleAmount := pv.getNestedFloat64Field("PolicyCoverageDetails", "Deductible", "Amount")
	if deductibleAmount < minDeductible {
		return fmt.Errorf("deductible too low (min $%.2f)", minDeductible)
	}
	if deductibleAmount > maxDeductible {
		return fmt.Errorf("deductible too high (max $%.2f)", maxDeductible)
	}

	return nil
}

// ValidateUnderwritingRules checks underwriting requirements
func (pv *PolicyValidator) ValidateUnderwritingRules() error {
	// High-value policies require manual underwriting
	coverageAmount := pv.getNestedFloat64Field("PolicyCoverageDetails", "CoverageAmount", "Amount")
	if coverageAmount > 25000 {
		underwritingStatus := pv.getNestedStringField("PolicyLifecycle", "UnderwritingStatus")
		if underwritingStatus != string(UnderwritingStatusApproved) {
			return errors.New("high-value policy requires underwriting approval")
		}
	}

	// High-risk policies require inspection
	riskScore := pv.getNestedFloat64Field("PolicyRiskAssessment", "", "RiskScore")
	if riskScore > 70 {
		requiresInspection := pv.getNestedBoolField("PolicyUnderwritingInfo", "RequiresInspection")
		if !requiresInspection {
			return errors.New("high-risk policy requires inspection")
		}
	}

	return nil
}

// ============================================
// FINANCIAL VALIDATIONS
// ============================================

// ValidateFinancials checks financial data consistency
func (pv *PolicyValidator) ValidateFinancials() error {
	// Premium validations
	basePremium := pv.getNestedFloat64Field("PolicyPricing", "BasePremium", "Amount")
	if basePremium <= 0 {
		return errors.New("base premium must be positive")
	}

	// Final premium must be at least base premium minus discounts
	minPremium := basePremium * 0.3 // Can't discount more than 70%
	finalPremium := pv.getNestedFloat64Field("PolicyPricing", "FinalPremium", "Amount")
	if finalPremium < minPremium {
		return fmt.Errorf("final premium too low (min $%.2f)", minPremium)
	}

	// Premium tax validation
	premiumTax := pv.getNestedFloat64Field("PolicyPricing", "PremiumTax", "Amount")
	if premiumTax < 0 {
		return errors.New("premium tax cannot be negative")
	}

	// Outstanding amount validation
	outstandingAmount := pv.getNestedFloat64Field("PolicyPaymentInfo", "OutstandingAmount", "Amount")
	if outstandingAmount < 0 {
		return errors.New("outstanding amount cannot be negative")
	}

	return nil
}

// ============================================
// DATE VALIDATIONS
// ============================================

// ValidateDates checks date consistency and business rules
func (pv *PolicyValidator) ValidateDates() error {
	now := time.Now()

	// Get effective date using reflection
	effectiveDateVal := pv.getNestedTimeField("PolicyLifecycle", "EffectiveDate")
	if !effectiveDateVal.IsValid() {
		return errors.New("effective date not found")
	}
	effectiveDate := effectiveDateVal.Interface().(time.Time)

	// Get expiration date using reflection
	expirationDateVal := pv.getNestedTimeField("PolicyLifecycle", "ExpirationDate")
	if !expirationDateVal.IsValid() {
		return errors.New("expiration date not found")
	}
	expirationDate := expirationDateVal.Interface().(time.Time)

	// Effective date validations
	if effectiveDate.Before(now.AddDate(0, 0, -30)) {
		return errors.New("effective date cannot be more than 30 days in the past")
	}
	if effectiveDate.After(now.AddDate(0, 0, 90)) {
		return errors.New("effective date cannot be more than 90 days in the future")
	}

	// Expiration date validations
	if !expirationDate.After(effectiveDate) {
		return errors.New("expiration date must be after effective date")
	}

	// Policy term validation
	minTerm := 30      // 30 days minimum
	maxTerm := 365 * 3 // 3 years maximum
	term := int(expirationDate.Sub(effectiveDate).Hours() / 24)

	if term < minTerm {
		return fmt.Errorf("policy term too short (min %d days)", minTerm)
	}
	if term > maxTerm {
		return fmt.Errorf("policy term too long (max %d days)", maxTerm)
	}

	return nil
}

// ============================================
// COVERAGE VALIDATIONS
// ============================================

// ValidateCoverage checks coverage configuration
func (pv *PolicyValidator) ValidateCoverage() error {
	// Coverage amount validations
	coverageAmount := pv.getNestedFloat64Field("PolicyCoverageDetails", "CoverageAmount", "Amount")
	if coverageAmount <= 0 {
		return errors.New("coverage amount must be positive")
	}

	// Deductible validations
	deductibleAmount := pv.getNestedFloat64Field("PolicyCoverageDetails", "Deductible", "Amount")
	if deductibleAmount < 0 {
		return errors.New("deductible cannot be negative")
	}
	if deductibleAmount > coverageAmount {
		return errors.New("deductible cannot exceed coverage amount")
	}

	// Coinsurance validation
	coinsurancePercent := pv.getNestedFloat64Field("PolicyCoverageDetails", "", "CoinsurancePercent")
	if coinsurancePercent < 0 || coinsurancePercent > 100 {
		return errors.New("coinsurance percentage must be between 0 and 100")
	}

	// Out of pocket max validation
	outOfPocketMax := pv.getNestedFloat64Field("PolicyCoverageDetails", "OutOfPocketMax", "Amount")
	if outOfPocketMax < 0 {
		return errors.New("out of pocket maximum cannot be negative")
	}

	// Remaining limit validation
	remainingLimit := pv.getNestedFloat64Field("PolicyCoverageDetails", "RemainingLimit", "Amount")
	if remainingLimit > coverageAmount {
		return errors.New("remaining limit cannot exceed coverage amount")
	}

	return nil
}

// ============================================
// RISK VALIDATIONS
// ============================================

// ValidateRisk checks risk assessment data
func (pv *PolicyValidator) ValidateRisk() error {
	// Risk score validation
	riskScore := pv.getNestedFloat64Field("PolicyRiskAssessment", "", "RiskScore")
	if riskScore < 0 || riskScore > 100 {
		return fmt.Errorf("risk score must be between 0 and 100")
	}

	// Component risk scores validation
	riskScores := []struct {
		name  string
		score float64
	}{
		{"device risk", pv.getNestedFloat64Field("PolicyRiskAssessment", "", "DeviceRiskScore")},
		{"user risk", pv.getNestedFloat64Field("PolicyRiskAssessment", "", "UserRiskScore")},
		{"location risk", pv.getNestedFloat64Field("PolicyRiskAssessment", "", "LocationRiskScore")},
		{"behavior risk", pv.getNestedFloat64Field("PolicyRiskAssessment", "", "BehaviorRiskScore")},
		{"fraud risk", pv.getNestedFloat64Field("PolicyRiskAssessment", "", "FraudRiskScore")},
	}

	for _, rs := range riskScores {
		if rs.score < 0 || rs.score > 100 {
			return fmt.Errorf("%s score must be between 0 and 100", rs.name)
		}
	}

	// Loss ratio validation
	lossRatio := pv.getNestedFloat64Field("PolicyRiskAssessment", "", "LossRatio")
	if lossRatio < 0 {
		return errors.New("loss ratio cannot be negative")
	}

	// Claim frequency validation
	claimFrequency := pv.getNestedFloat64Field("PolicyRiskAssessment", "", "ClaimFrequency")
	if claimFrequency < 0 {
		return errors.New("claim frequency cannot be negative")
	}

	return nil
}

// ============================================
// COMPLIANCE VALIDATIONS
// ============================================

// ValidateCompliance checks regulatory compliance
func (pv *PolicyValidator) ValidateCompliance() error {
	// KYC validation
	totalAmount := pv.getNestedFloat64Field("PolicyPricing", "TotalAmount", "Amount")
	if totalAmount > 10000 {
		kycStatus := pv.getNestedStringField("PolicyCompliance", "KYCStatus")
		if kycStatus != "verified" {
			return errors.New("KYC verification required for high-value policies")
		}
	}

	// AML validation
	if totalAmount > 50000 {
		amlStatus := pv.getNestedStringField("PolicyCompliance", "AMLStatus")
		if amlStatus != "clear" {
			return errors.New("AML clearance required for policies over $50,000")
		}
	}

	// Sanctions check
	sanctionsCheck := pv.getNestedBoolField("PolicyCompliance", "SanctionsCheck")
	if sanctionsCheck {
		return errors.New("policy cannot be issued to sanctioned entity")
	}

	// Data privacy consent
	dataPrivacyConsent := pv.getNestedBoolField("PolicyCompliance", "DataPrivacyConsent")
	if !dataPrivacyConsent {
		return errors.New("data privacy consent is required")
	}

	// Regional compliance
	region := pv.getNestedStringField("PolicyCompliance", "RegulatoryRegion")
	switch region {
	case "US-CA": // California
		gracePeriod := pv.getNestedIntField("PolicyPaymentInfo", "PaymentGracePeriod")
		if gracePeriod < 10 {
			return errors.New("California requires minimum 10-day grace period")
		}
	case "EU": // European Union
		gdprCompliant := pv.getNestedBoolField("PolicyCompliance", "GDPRCompliant")
		if !gdprCompliant {
			return errors.New("GDPR compliance required for EU policies")
		}
	}

	return nil
}
