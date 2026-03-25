package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PolicyLegal represents legal and litigation matters for a policy
type PolicyLegal struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Litigation
	ActiveLitigations   int            `gorm:"type:int" json:"active_litigations_count"`
	LitigationDetails   datatypes.JSON `gorm:"type:json" json:"litigation_details"` // []Litigation
	TotalLitigationCost Money          `gorm:"embedded;embeddedPrefix:litigation_cost_" json:"total_litigation_cost"`
	LitigationReserve   Money          `gorm:"embedded;embeddedPrefix:litigation_reserve_" json:"litigation_reserve"`

	// Disputes
	DisputeCount          int            `gorm:"type:int" json:"dispute_count"`
	DisputeDetails        datatypes.JSON `gorm:"type:json" json:"dispute_details"` // []Dispute
	PendingDisputes       int            `gorm:"type:int" json:"pending_disputes"`
	ResolvedDisputes      int            `gorm:"type:int" json:"resolved_disputes"`
	DisputeResolutionRate float64        `gorm:"type:decimal(5,2)" json:"dispute_resolution_rate"`

	// Arbitration
	ArbitrationCases    int            `gorm:"type:int" json:"arbitration_cases"`
	ArbitrationDetails  datatypes.JSON `gorm:"type:json" json:"arbitration_details"` // []Arbitration
	ArbitrationCosts    Money          `gorm:"embedded;embeddedPrefix:arbitration_costs_" json:"arbitration_costs"`
	ArbitrationOutcomes datatypes.JSON `gorm:"type:json" json:"arbitration_outcomes"`

	// Settlements
	SettlementCount       int            `gorm:"type:int" json:"settlement_count"`
	TotalSettlementAmount Money          `gorm:"embedded;embeddedPrefix:total_settlement_" json:"total_settlement_amount"`
	PendingSettlements    Money          `gorm:"embedded;embeddedPrefix:pending_settlements_" json:"pending_settlements"`
	SettlementAgreements  datatypes.JSON `gorm:"type:json" json:"settlement_agreements"` // []Settlement

	// Legal Holds
	LegalHoldStatus     bool           `gorm:"type:boolean;default:false" json:"legal_hold_status"`
	LegalHoldReason     string         `gorm:"type:text" json:"legal_hold_reason"`
	LegalHoldDate       *time.Time     `gorm:"type:timestamp" json:"legal_hold_date,omitempty"`
	LegalHoldExpiry     *time.Time     `gorm:"type:timestamp" json:"legal_hold_expiry,omitempty"`
	PreservationNotices datatypes.JSON `gorm:"type:json" json:"preservation_notices"`

	// Subrogation
	SubrogationCases     int            `gorm:"type:int" json:"subrogation_cases"`
	SubrogationDetails   datatypes.JSON `gorm:"type:json" json:"subrogation_details"` // []Subrogation
	SubrogationRecovered Money          `gorm:"embedded;embeddedPrefix:subrogation_recovered_" json:"subrogation_recovered"`
	SubrogationPending   Money          `gorm:"embedded;embeddedPrefix:subrogation_pending_" json:"subrogation_pending"`
	SubrogationSuccess   float64        `gorm:"type:decimal(5,2)" json:"subrogation_success_rate"`

	// Recovery Actions
	RecoveryActions int            `gorm:"type:int" json:"recovery_actions"`
	RecoveryDetails datatypes.JSON `gorm:"type:json" json:"recovery_details"` // []Recovery
	TotalRecovered  Money          `gorm:"embedded;embeddedPrefix:total_recovered_" json:"total_recovered"`
	RecoveryPending Money          `gorm:"embedded;embeddedPrefix:recovery_pending_" json:"recovery_pending"`

	// Legal Representation
	LegalCounsel   string `gorm:"type:varchar(255)" json:"legal_counsel"`
	LawFirm        string `gorm:"type:varchar(255)" json:"law_firm"`
	LegalFees      Money  `gorm:"embedded;embeddedPrefix:legal_fees_" json:"legal_fees"`
	RetainerAmount Money  `gorm:"embedded;embeddedPrefix:retainer_" json:"retainer_amount"`

	// Court Information
	CourtCases   datatypes.JSON `gorm:"type:json" json:"court_cases"`
	Jurisdiction string         `gorm:"type:varchar(255)" json:"jurisdiction"`
	VenueDetails string         `gorm:"type:text" json:"venue_details"`

	// Compliance & Documentation
	LegalDocuments       datatypes.JSON `gorm:"type:json" json:"legal_documents"` // []LegalDocument
	ComplianceIssues     datatypes.JSON `gorm:"type:json" json:"compliance_issues"`
	RegulatoryViolations int            `gorm:"type:int" json:"regulatory_violations"`
	PenaltiesAmount      Money          `gorm:"embedded;embeddedPrefix:penalties_" json:"penalties_amount"`

	// Status
	LegalStatus    string    `gorm:"type:varchar(50)" json:"legal_status"`
	RiskAssessment string    `gorm:"type:varchar(50)" json:"risk_assessment"` // low, medium, high, critical
	LastReviewDate time.Time `gorm:"type:timestamp" json:"last_review_date"`
	NextReviewDate time.Time `gorm:"type:timestamp" json:"next_review_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// PolicyRegulatoryFiling represents regulatory filing requirements
type PolicyRegulatoryFiling struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Filing Requirements
	FilingRequirements datatypes.JSON `gorm:"type:json" json:"filing_requirements"` // []FilingRequirement
	StatutoryReports   datatypes.JSON `gorm:"type:json" json:"statutory_reports"`   // []Report
	RequiredFilings    int            `gorm:"type:int" json:"required_filings_count"`
	CompletedFilings   int            `gorm:"type:int" json:"completed_filings_count"`
	PendingFilings     int            `gorm:"type:int" json:"pending_filings_count"`

	// Capital Requirements
	RegulatoryCapital    Money   `gorm:"embedded;embeddedPrefix:regulatory_capital_" json:"regulatory_capital"`
	MinimumCapitalReq    Money   `gorm:"embedded;embeddedPrefix:mcr_" json:"minimum_capital_requirement"`
	SolvencyCapitalReq   Money   `gorm:"embedded;embeddedPrefix:scr_" json:"solvency_capital_requirement"`
	ActualCapital        Money   `gorm:"embedded;embeddedPrefix:actual_capital_" json:"actual_capital"`
	CapitalAdequacyRatio float64 `gorm:"type:decimal(10,4)" json:"capital_adequacy_ratio"`

	// Solvency Metrics
	SolvencyRatio    float64   `gorm:"type:decimal(10,4)" json:"solvency_ratio"`
	SolvencyStatus   string    `gorm:"type:varchar(50)" json:"solvency_status"`
	SolvencyTestDate time.Time `gorm:"type:timestamp" json:"solvency_test_date"`
	OwnFundsAmount   Money     `gorm:"embedded;embeddedPrefix:own_funds_" json:"own_funds_amount"`
	TierOneCapital   Money     `gorm:"embedded;embeddedPrefix:tier_one_capital_" json:"tier_one_capital"`
	TierTwoCapital   Money     `gorm:"embedded;embeddedPrefix:tier_two_capital_" json:"tier_two_capital"`

	// Stress Testing
	StressTestResults   datatypes.JSON `gorm:"type:json" json:"stress_test_results"` // []StressTestResult
	LastStressTestDate  time.Time      `gorm:"type:timestamp" json:"last_stress_test_date"`
	StressTestScenarios datatypes.JSON `gorm:"type:json" json:"stress_test_scenarios"`
	WorstCaseCapital    Money          `gorm:"embedded;embeddedPrefix:worst_case_capital_" json:"worst_case_capital"`

	// QIS (Quantitative Impact Studies)
	QISParticipation  bool           `gorm:"type:boolean;default:false" json:"qis_participation"`
	QISResults        datatypes.JSON `gorm:"type:json" json:"qis_results"` // []QISResult
	QISSubmissionDate *time.Time     `gorm:"type:timestamp" json:"qis_submission_date,omitempty"`

	// ORSA (Own Risk and Solvency Assessment)
	ORSAReport       datatypes.JSON `gorm:"type:json" json:"orsa_report"`
	ORSAApprovalDate *time.Time     `gorm:"type:timestamp" json:"orsa_approval_date,omitempty"`
	ORSAFrequency    string         `gorm:"type:varchar(50)" json:"orsa_frequency"`
	NextORSADate     time.Time      `gorm:"type:timestamp" json:"next_orsa_date"`

	// Regulatory Reporting
	RegulatorName      string    `gorm:"type:varchar(255)" json:"regulator_name"`
	RegulatoryRegion   string    `gorm:"type:varchar(100)" json:"regulatory_region"`
	ReportingFrequency string    `gorm:"type:varchar(50)" json:"reporting_frequency"`
	LastReportingDate  time.Time `gorm:"type:timestamp" json:"last_reporting_date"`
	NextReportingDate  time.Time `gorm:"type:timestamp" json:"next_reporting_date"`

	// Compliance Metrics
	ComplianceScore  float64        `gorm:"type:decimal(5,2)" json:"compliance_score"`
	ComplianceStatus string         `gorm:"type:varchar(50)" json:"compliance_status"`
	ComplianceIssues datatypes.JSON `gorm:"type:json" json:"compliance_issues"`
	RemediationPlans datatypes.JSON `gorm:"type:json" json:"remediation_plans"`

	// Regulatory Actions
	RegulatoryActions datatypes.JSON `gorm:"type:json" json:"regulatory_actions"`
	Sanctions         datatypes.JSON `gorm:"type:json" json:"sanctions"`
	Fines             Money          `gorm:"embedded;embeddedPrefix:fines_" json:"fines"`
	Warnings          int            `gorm:"type:int" json:"warnings_count"`

	// Pillar Reports (Solvency II)
	Pillar1Requirements datatypes.JSON `gorm:"type:json" json:"pillar1_requirements"`
	Pillar2Requirements datatypes.JSON `gorm:"type:json" json:"pillar2_requirements"`
	Pillar3Disclosures  datatypes.JSON `gorm:"type:json" json:"pillar3_disclosures"`

	// Documentation
	RegulatoryDocuments datatypes.JSON `gorm:"type:json" json:"regulatory_documents"`
	FilingHistory       datatypes.JSON `gorm:"type:json" json:"filing_history"`
	CorrespondenceLog   datatypes.JSON `gorm:"type:json" json:"correspondence_log"`

	// Status
	FilingStatus   string     `gorm:"type:varchar(50)" json:"filing_status"`
	ReviewStatus   string     `gorm:"type:varchar(50)" json:"review_status"`
	ApprovalStatus string     `gorm:"type:varchar(50)" json:"approval_status"`
	LastAuditDate  *time.Time `gorm:"type:timestamp" json:"last_audit_date,omitempty"`
	NextAuditDate  *time.Time `gorm:"type:timestamp" json:"next_audit_date,omitempty"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// =====================================
// METHODS
// =====================================

// HasActiveLitigation checks if there are active litigation cases
func (pl *PolicyLegal) HasActiveLitigation() bool {
	return pl.ActiveLitigations > 0
}

// IsUnderLegalHold checks if policy is under legal hold
func (pl *PolicyLegal) IsUnderLegalHold() bool {
	if !pl.LegalHoldStatus {
		return false
	}
	if pl.LegalHoldExpiry != nil {
		return time.Now().Before(*pl.LegalHoldExpiry)
	}
	return true
}

// GetTotalLegalExposure calculates total legal exposure
func (pl *PolicyLegal) GetTotalLegalExposure() float64 {
	return pl.TotalLitigationCost.Amount + pl.ArbitrationCosts.Amount +
		pl.PendingSettlements.Amount + pl.SubrogationPending.Amount +
		pl.LegalFees.Amount + pl.PenaltiesAmount.Amount
}

// GetSubrogationRecoveryRate calculates subrogation recovery rate
func (pl *PolicyLegal) GetSubrogationRecoveryRate() float64 {
	totalSubrogation := pl.SubrogationRecovered.Amount + pl.SubrogationPending.Amount
	if totalSubrogation > 0 {
		return (pl.SubrogationRecovered.Amount / totalSubrogation) * 100
	}
	return 0
}

// IsHighRisk determines if legal risk is high
func (pl *PolicyLegal) IsHighRisk() bool {
	return pl.RiskAssessment == "high" || pl.RiskAssessment == "critical" ||
		pl.ActiveLitigations > 2 || pl.RegulatoryViolations > 0
}

// IsSolvencyCompliant checks if solvency requirements are met
func (prf *PolicyRegulatoryFiling) IsSolvencyCompliant() bool {
	return prf.SolvencyRatio >= 1.0 &&
		prf.ActualCapital.Amount >= prf.MinimumCapitalReq.Amount &&
		prf.ActualCapital.Amount >= prf.SolvencyCapitalReq.Amount
}

// GetCapitalBuffer calculates capital buffer percentage
func (prf *PolicyRegulatoryFiling) GetCapitalBuffer() float64 {
	if prf.MinimumCapitalReq.Amount > 0 {
		return ((prf.ActualCapital.Amount - prf.MinimumCapitalReq.Amount) /
			prf.MinimumCapitalReq.Amount) * 100
	}
	return 0
}

// IsFilingOverdue checks if any filings are overdue
func (prf *PolicyRegulatoryFiling) IsFilingOverdue() bool {
	return time.Now().After(prf.NextReportingDate) || prf.PendingFilings > 0
}

// GetCompliancePercentage calculates compliance percentage
func (prf *PolicyRegulatoryFiling) GetCompliancePercentage() float64 {
	if prf.RequiredFilings > 0 {
		return (float64(prf.CompletedFilings) / float64(prf.RequiredFilings)) * 100
	}
	return 100
}

// RequiresORSA checks if ORSA is required
func (prf *PolicyRegulatoryFiling) RequiresORSA() bool {
	return time.Now().After(prf.NextORSADate) || prf.ORSAReport == nil
}

// HasRegulatoryActions checks for regulatory actions
func (prf *PolicyRegulatoryFiling) HasRegulatoryActions() bool {
	return prf.RegulatoryActions != nil || prf.Sanctions != nil ||
		prf.Fines.Amount > 0 || prf.Warnings > 0
}
