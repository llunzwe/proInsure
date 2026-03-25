package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceClaimHistory tracks complete claim history for a device
type DeviceClaimHistory struct {
	database.BaseModel
	DeviceID             uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id"`
	ClaimID              uuid.UUID  `gorm:"type:uuid;not null;index" json:"claim_id"`
	ClaimType            string     `gorm:"type:varchar(50);not null" json:"claim_type"`   // theft, damage, loss, malfunction
	ClaimStatus          string     `gorm:"type:varchar(50);not null" json:"claim_status"` // pending, approved, rejected, paid
	ClaimAmount          float64    `json:"claim_amount"`
	ApprovedAmount       float64    `json:"approved_amount"`
	DeductibleApplied    float64    `json:"deductible_applied"`
	ClaimDate            time.Time  `json:"claim_date"`
	ResolutionDate       *time.Time `json:"resolution_date"`
	ProcessingDays       int        `json:"processing_days"`
	DamageType           string     `json:"damage_type"`                            // screen, water, physical, battery, other
	RepairOrReplace      string     `json:"repair_or_replace"`                      // repair, replace, cash_settlement
	DocumentsProvided    string     `gorm:"type:json" json:"documents_provided"`    // JSON array
	PreviousClaims       int        `json:"previous_claims"`                        // Count at time of claim
	SuspiciousIndicators string     `gorm:"type:json" json:"suspicious_indicators"` // JSON array
	InvestigationNotes   string     `gorm:"type:text" json:"investigation_notes"`
}

// DevicePremiumCalculation manages dynamic premium calculations
type DevicePremiumCalculation struct {
	database.BaseModel
	DeviceID            uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id"`
	PolicyID            uuid.UUID  `gorm:"type:uuid;index" json:"policy_id"`
	BasePremium         float64    `json:"base_premium"`
	RiskAdjustedPremium float64    `json:"risk_adjusted_premium"`
	FinalPremium        float64    `json:"final_premium"`
	CalculationDate     time.Time  `json:"calculation_date"`
	EffectiveDate       time.Time  `json:"effective_date"`
	ExpiryDate          *time.Time `json:"expiry_date"`

	// Risk Factors
	AgeMultiplier          float64 `json:"age_multiplier"`
	ConditionMultiplier    float64 `json:"condition_multiplier"`
	LocationMultiplier     float64 `json:"location_multiplier"`
	UsageMultiplier        float64 `json:"usage_multiplier"`
	ClaimHistoryMultiplier float64 `json:"claim_history_multiplier"`

	// Discounts
	LoyaltyDiscount     float64 `json:"loyalty_discount"`
	NoClaimBonus        float64 `json:"no_claim_bonus"`
	MultiDeviceDiscount float64 `json:"multi_device_discount"`
	ReferralDiscount    float64 `json:"referral_discount"`
	PrepaymentDiscount  float64 `json:"prepayment_discount"`
	AutoRenewalDiscount float64 `json:"auto_renewal_discount"`
	TotalDiscounts      float64 `json:"total_discounts"`

	// Payment History
	PaymentStatus    string     `gorm:"type:varchar(50)" json:"payment_status"` // current, late, defaulted
	LatePaymentCount int        `json:"late_payment_count"`
	LastPaymentDate  *time.Time `json:"last_payment_date"`
	NextPaymentDue   *time.Time `json:"next_payment_due"`

	// Adjustment Triggers
	LastAdjustmentReason string `json:"last_adjustment_reason"`
	AdjustmentHistory    string `gorm:"type:json" json:"adjustment_history"` // JSON array
}

// DeviceRiskProfile provides comprehensive risk assessment
type DeviceRiskProfile struct {
	database.BaseModel
	DeviceID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	OverallRiskScore   float64   `json:"overall_risk_score"`                    // 0-100
	RiskCategory       string    `gorm:"type:varchar(20)" json:"risk_category"` // low, medium, high, very_high
	LastAssessmentDate time.Time `json:"last_assessment_date"`
	NextAssessmentDue  time.Time `json:"next_assessment_due"`

	// Risk Components
	TheftRiskScore      float64 `json:"theft_risk_score"`
	DamageRiskScore     float64 `json:"damage_risk_score"`
	FraudRiskScore      float64 `json:"fraud_risk_score"`
	UsageRiskScore      float64 `json:"usage_risk_score"`
	LocationRiskScore   float64 `json:"location_risk_score"`
	BehavioralRiskScore float64 `json:"behavioral_risk_score"`

	// Risk Indicators
	HighValueDevice  bool `json:"high_value_device"`
	FrequentTraveler bool `json:"frequent_traveler"`
	HighRiskLocation bool `json:"high_risk_location"`
	PreviousClaimer  bool `json:"previous_claimer"`
	MultipleDevices  bool `json:"multiple_devices"`
	BusinessUse      bool `json:"business_use"`

	// Predictions
	ClaimProbability   float64 `json:"claim_probability"` // Next 12 months
	TheftProbability   float64 `json:"theft_probability"`
	DamageProbability  float64 `json:"damage_probability"`
	ExpectedLossAmount float64 `json:"expected_loss_amount"`

	// Historical Data
	RiskScoreHistory  string `gorm:"type:json" json:"risk_score_history"` // JSON array
	RiskEvents        string `gorm:"type:json" json:"risk_events"`        // JSON array
	MitigationActions string `gorm:"type:json" json:"mitigation_actions"` // JSON array
}

// DeviceCoverageOptions defines available coverage tiers and options
type DeviceCoverageOptions struct {
	database.BaseModel
	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	PolicyID     uuid.UUID `gorm:"type:uuid;index" json:"policy_id"`
	CoverageTier string    `gorm:"type:varchar(50);not null" json:"coverage_tier"` // basic, standard, premium, comprehensive

	// Coverage Limits
	MaxClaimAmount     float64 `json:"max_claim_amount"`
	AnnualClaimLimit   float64 `json:"annual_claim_limit"`
	LifetimeClaimLimit float64 `json:"lifetime_claim_limit"`
	PerIncidentLimit   float64 `json:"per_incident_limit"`

	// Coverage Types
	TheftCoverage       bool `gorm:"default:true" json:"theft_coverage"`
	DamageCoverage      bool `gorm:"default:true" json:"damage_coverage"`
	LossCoverage        bool `gorm:"default:false" json:"loss_coverage"`
	MalfunctionCoverage bool `gorm:"default:false" json:"malfunction_coverage"`

	// Add-on Coverage
	ScreenProtection      bool `json:"screen_protection"`
	WaterDamageProtection bool `json:"water_damage_protection"`
	InternationalCoverage bool `json:"international_coverage"`
	AccessoriesCoverage   bool `json:"accessories_coverage"`
	DataRecoveryCoverage  bool `json:"data_recovery_coverage"`
	CyberProtection       bool `json:"cyber_protection"`
	ExtendedWarranty      bool `json:"extended_warranty"`

	// Special Provisions
	NewForOld           bool `json:"new_for_old"` // Replacement with new device
	WorldwideCoverage   bool `json:"worldwide_coverage"`
	BusinessUseCoverage bool `json:"business_use_coverage"`
	LoanerDevice        bool `json:"loaner_device"`
	ExpressReplacement  bool `json:"express_replacement"`

	// Coverage History
	UpgradeHistory      string `gorm:"type:json" json:"upgrade_history"`      // JSON array
	DowngradeHistory    string `gorm:"type:json" json:"downgrade_history"`    // JSON array
	TemporaryExtensions string `gorm:"type:json" json:"temporary_extensions"` // JSON array

	// Coverage Gaps
	IdentifiedGaps    string `gorm:"type:json" json:"identified_gaps"`    // JSON array
	RecommendedAddons string `gorm:"type:json" json:"recommended_addons"` // JSON array
}

// DeviceDeductibles manages deductible information
type DeviceDeductibles struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	PolicyID uuid.UUID `gorm:"type:uuid;index" json:"policy_id"`

	// Current Deductibles
	StandardDeductible float64 `json:"standard_deductible"`
	TheftDeductible    float64 `json:"theft_deductible"`
	DamageDeductible   float64 `json:"damage_deductible"`
	ScreenDeductible   float64 `json:"screen_deductible"`

	// Variable Options
	PercentageDeductible bool    `json:"percentage_deductible"` // As % of claim
	DeductiblePercentage float64 `json:"deductible_percentage"`
	MinimumDeductible    float64 `json:"minimum_deductible"`
	MaximumDeductible    float64 `json:"maximum_deductible"`

	// Deductible Credits
	ClaimFreeMonths     int     `json:"claim_free_months"`
	DeductibleReduction float64 `json:"deductible_reduction"`
	LoyaltyCredit       float64 `json:"loyalty_credit"`
	PrepaymentCredit    float64 `json:"prepayment_credit"`
	TotalCredits        float64 `json:"total_credits"`

	// Payment Tracking
	DeductiblesPaid       float64    `json:"deductibles_paid"`
	LastDeductiblePaid    *time.Time `json:"last_deductible_paid"`
	OutstandingDeductible float64    `json:"outstanding_deductible"`

	// Annual Management
	AnnualDeductibleUsed float64   `json:"annual_deductible_used"`
	AnnualMaxDeductible  float64   `json:"annual_max_deductible"`
	ResetDate            time.Time `json:"reset_date"`

	// Family Aggregation
	FamilyPlanID         *uuid.UUID `gorm:"type:uuid" json:"family_plan_id"`
	FamilyDeductiblePool float64    `json:"family_deductible_pool"`
	FamilyContribution   float64    `json:"family_contribution"`

	// History
	DeductibleHistory string `gorm:"type:json" json:"deductible_history"` // JSON array
	WaiverHistory     string `gorm:"type:json" json:"waiver_history"`     // JSON array
}

// DeviceInsuranceAudit tracks insurance compliance and audits
type DeviceInsuranceAudit struct {
	database.BaseModel
	DeviceID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id"`
	AuditType string     `gorm:"type:varchar(50);not null" json:"audit_type"` // eligibility, compliance, renewal, claim, fraud
	AuditDate time.Time  `json:"audit_date"`
	AuditorID *uuid.UUID `gorm:"type:uuid" json:"auditor_id"`

	// Audit Results
	AuditStatus     string  `gorm:"type:varchar(50)" json:"audit_status"` // passed, failed, conditional, pending
	ComplianceScore float64 `json:"compliance_score"`                     // 0-100
	RiskFlags       string  `gorm:"type:json" json:"risk_flags"`          // JSON array

	// Eligibility Checks
	EligibilityStatus  bool   `json:"eligibility_status"`
	EligibilityReasons string `gorm:"type:json" json:"eligibility_reasons"` // JSON array
	RequiredDocuments  string `gorm:"type:json" json:"required_documents"`  // JSON array
	DocumentsVerified  string `gorm:"type:json" json:"documents_verified"`  // JSON array

	// Compliance Verification
	PolicyCompliance   bool `json:"policy_compliance"`
	PremiumCompliance  bool `json:"premium_compliance"`
	ClaimCompliance    bool `json:"claim_compliance"`
	DocumentCompliance bool `json:"document_compliance"`

	// Violations & Exceptions
	ViolationsFound   string `gorm:"type:json" json:"violations_found"`   // JSON array
	ExceptionsGranted string `gorm:"type:json" json:"exceptions_granted"` // JSON array
	CorrectiveActions string `gorm:"type:json" json:"corrective_actions"` // JSON array

	// Third-party Verification
	ThirdPartyRequired bool       `json:"third_party_required"`
	ThirdPartyProvider string     `json:"third_party_provider"`
	VerificationStatus string     `json:"verification_status"`
	VerificationDate   *time.Time `json:"verification_date"`

	// Follow-up
	NextAuditDate    *time.Time `json:"next_audit_date"`
	FollowUpRequired bool       `json:"follow_up_required"`
	Notes            string     `gorm:"type:text" json:"notes"`
}

// DeviceClaimValidation handles automated claim validation
type DeviceClaimValidation struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ClaimID        uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`
	ValidationDate time.Time `json:"validation_date"`

	// Validation Status
	OverallStatus   string  `gorm:"type:varchar(50)" json:"overall_status"` // valid, invalid, suspicious, pending_review
	ValidationScore float64 `json:"validation_score"`                       // 0-100
	AutoApproved    bool    `json:"auto_approved"`
	RequiresReview  bool    `json:"requires_review"`

	// Document Validation
	DocumentsRequired string `gorm:"type:json" json:"documents_required"` // JSON array
	DocumentsProvided string `gorm:"type:json" json:"documents_provided"` // JSON array
	DocumentsVerified bool   `json:"documents_verified"`

	// Evidence Validation
	PhotosProvided int  `json:"photos_provided"`
	PhotosVerified bool `json:"photos_verified"`
	VideoProvided  bool `json:"video_provided"`
	VideoVerified  bool `json:"video_verified"`

	// Report Verification
	PoliceReportRequired bool   `json:"police_report_required"`
	PoliceReportProvided bool   `json:"police_report_provided"`
	PoliceReportVerified bool   `json:"police_report_verified"`
	ReportNumber         string `json:"report_number"`

	// Damage Assessment
	DamageConsistent    bool    `json:"damage_consistent"`
	RepairEstimateValid bool    `json:"repair_estimate_valid"`
	EstimateAmount      float64 `json:"estimate_amount"`
	ApprovedAmount      float64 `json:"approved_amount"`

	// Cross-checking
	PreviousClaimsCheck bool `json:"previous_claims_check"`
	DuplicateClaimFound bool `json:"duplicate_claim_found"`
	TimelineValid       bool `json:"timeline_valid"`
	LocationValid       bool `json:"location_valid"`

	// Ownership Verification
	OwnershipProofValid  bool `json:"ownership_proof_valid"`
	DeviceConditionMatch bool `json:"device_condition_match"`
	IMEIVerified         bool `json:"imei_verified"`

	// Validation Rules Applied
	RulesApplied   string `gorm:"type:json" json:"rules_applied"`   // JSON array
	RuleViolations string `gorm:"type:json" json:"rule_violations"` // JSON array

	// AI/ML Scoring
	FraudProbability float64 `json:"fraud_probability"`
	AnomalyScore     float64 `json:"anomaly_score"`
	ConfidenceLevel  float64 `json:"confidence_level"`

	// Notes
	SystemNotes   string `gorm:"type:text" json:"system_notes"`
	ReviewerNotes string `gorm:"type:text" json:"reviewer_notes"`
}

// Methods for DeviceClaimHistory
func (dch *DeviceClaimHistory) CalculateClaimFrequency(months int) float64 {
	// Implementation for calculating claim frequency
	return 0.0
}

func (dch *DeviceClaimHistory) IsSuspicious() bool {
	// Check for suspicious patterns
	return dch.SuspiciousIndicators != "" && dch.SuspiciousIndicators != "[]"
}

// Methods for DevicePremiumCalculation
func (dpc *DevicePremiumCalculation) CalculateFinalPremium() float64 {
	finalPremium := dpc.BasePremium

	// Apply risk multipliers
	finalPremium *= dpc.AgeMultiplier
	finalPremium *= dpc.ConditionMultiplier
	finalPremium *= dpc.LocationMultiplier
	finalPremium *= dpc.UsageMultiplier
	finalPremium *= dpc.ClaimHistoryMultiplier

	// Apply discounts
	finalPremium -= dpc.TotalDiscounts

	if finalPremium < 0 {
		finalPremium = 0
	}

	dpc.FinalPremium = finalPremium
	return finalPremium
}

func (dpc *DevicePremiumCalculation) IsPaymentCurrent() bool {
	return dpc.PaymentStatus == "current"
}

// Methods for DeviceRiskProfile
func (drp *DeviceRiskProfile) CalculateOverallRisk() float64 {
	// Weight different risk components
	weights := map[string]float64{
		"theft":      0.25,
		"damage":     0.25,
		"fraud":      0.20,
		"usage":      0.10,
		"location":   0.10,
		"behavioral": 0.10,
	}

	overall := drp.TheftRiskScore*weights["theft"] +
		drp.DamageRiskScore*weights["damage"] +
		drp.FraudRiskScore*weights["fraud"] +
		drp.UsageRiskScore*weights["usage"] +
		drp.LocationRiskScore*weights["location"] +
		drp.BehavioralRiskScore*weights["behavioral"]

	drp.OverallRiskScore = overall

	// Set category
	switch {
	case overall < 25:
		drp.RiskCategory = "low"
	case overall < 50:
		drp.RiskCategory = "medium"
	case overall < 75:
		drp.RiskCategory = "high"
	default:
		drp.RiskCategory = "very_high"
	}

	return overall
}

func (drp *DeviceRiskProfile) RequiresReassessment() bool {
	return time.Since(drp.LastAssessmentDate) > 30*24*time.Hour
}

// Methods for DeviceCoverageOptions
func (dco *DeviceCoverageOptions) HasComprehensiveCoverage() bool {
	return dco.CoverageTier == "comprehensive" || dco.CoverageTier == "premium"
}

func (dco *DeviceCoverageOptions) GetTotalCoverageValue() float64 {
	return dco.MaxClaimAmount
}

// Methods for DeviceDeductibles
func (dd *DeviceDeductibles) CalculateEffectiveDeductible(claimType string) float64 {
	var baseDeductible float64

	switch claimType {
	case "theft":
		baseDeductible = dd.TheftDeductible
	case "damage":
		baseDeductible = dd.DamageDeductible
	case "screen":
		baseDeductible = dd.ScreenDeductible
	default:
		baseDeductible = dd.StandardDeductible
	}

	// Apply credits
	effectiveDeductible := baseDeductible - dd.TotalCredits

	if effectiveDeductible < dd.MinimumDeductible {
		effectiveDeductible = dd.MinimumDeductible
	}

	return effectiveDeductible
}

func (dd *DeviceDeductibles) HasOutstandingDeductible() bool {
	return dd.OutstandingDeductible > 0
}

// Methods for DeviceInsuranceAudit
func (dia *DeviceInsuranceAudit) IsPassed() bool {
	return dia.AuditStatus == "passed"
}

func (dia *DeviceInsuranceAudit) RequiresAction() bool {
	return dia.AuditStatus == "failed" || dia.FollowUpRequired
}

// Methods for DeviceClaimValidation
func (dcv *DeviceClaimValidation) IsValid() bool {
	return dcv.OverallStatus == "valid" && !dcv.RequiresReview
}

func (dcv *DeviceClaimValidation) CanAutoApprove() bool {
	return dcv.AutoApproved &&
		dcv.DocumentsVerified &&
		!dcv.DuplicateClaimFound &&
		dcv.FraudProbability < 0.3
}
