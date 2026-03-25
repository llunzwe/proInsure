package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceComplianceStatus tracks regulatory compliance for devices
type DeviceComplianceStatus struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ComplianceDate time.Time `json:"compliance_date"`

	// Overall Compliance
	OverallStatus      string    `json:"overall_status"`   // compliant, partial, non-compliant, pending
	ComplianceScore    float64   `json:"compliance_score"` // 0-100
	LastAssessmentDate time.Time `json:"last_assessment_date"`
	NextAssessmentDate time.Time `json:"next_assessment_date"`

	// Regional Compliance
	RegionalCompliance string `gorm:"type:json" json:"regional_compliance"` // JSON object by region
	CountryCompliance  string `gorm:"type:json" json:"country_compliance"`  // JSON object by country
	StateCompliance    string `gorm:"type:json" json:"state_compliance"`    // JSON object by state/province

	// Industry Standards
	ISOCompliance     bool   `json:"iso_compliance"`
	ISOCertifications string `gorm:"type:json" json:"iso_certifications"` // JSON array
	NISTCompliance    bool   `json:"nist_compliance"`
	SOCCompliance     string `json:"soc_compliance"` // SOC1, SOC2, SOC3
	PCIDSSCompliance  bool   `json:"pci_dss_compliance"`
	HIPAACompliance   bool   `json:"hipaa_compliance"`

	// Certification Status
	Certifications        string `gorm:"type:json" json:"certifications"`         // JSON array
	CertificationExpiry   string `gorm:"type:json" json:"certification_expiry"`   // JSON object
	PendingCertifications string `gorm:"type:json" json:"pending_certifications"` // JSON array
	ExpiredCertifications string `gorm:"type:json" json:"expired_certifications"` // JSON array

	// Audit History
	LastAuditDate   *time.Time `json:"last_audit_date"`
	LastAuditResult string     `json:"last_audit_result"`
	AuditHistory    string     `gorm:"type:json" json:"audit_history"` // JSON array
	NextAuditDate   *time.Time `json:"next_audit_date"`
	AuditorName     string     `json:"auditor_name"`

	// Non-Compliance Tracking
	NonComplianceCount int    `json:"non_compliance_count"`
	ActiveViolations   string `gorm:"type:json" json:"active_violations"`   // JSON array
	ResolvedViolations string `gorm:"type:json" json:"resolved_violations"` // JSON array
	ViolationSeverity  string `json:"violation_severity"`                   // low, medium, high, critical

	// Remediation
	RemediationRequired bool       `json:"remediation_required"`
	RemediationPlan     string     `gorm:"type:json" json:"remediation_plan"` // JSON object
	RemediationDeadline *time.Time `json:"remediation_deadline"`
	RemediationStatus   string     `json:"remediation_status"` // pending, in_progress, completed
	RemediationCost     float64    `json:"remediation_cost"`

	// Regulatory Changes
	ImpactedByChanges  bool   `json:"impacted_by_changes"`
	PendingChanges     string `gorm:"type:json" json:"pending_changes"`     // JSON array
	ChangeDeadlines    string `gorm:"type:json" json:"change_deadlines"`    // JSON object
	AdaptationRequired string `gorm:"type:json" json:"adaptation_required"` // JSON array

	// Cost Tracking
	ComplianceCost      float64 `json:"compliance_cost"`
	PenaltyCost         float64 `json:"penalty_cost"`
	AuditCost           float64 `json:"audit_cost"`
	CertificationCost   float64 `json:"certification_cost"`
	TotalComplianceCost float64 `json:"total_compliance_cost"`
}

// DeviceLegalHolds manages legal restrictions and holds on devices
type DeviceLegalHolds struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	HoldID   string    `gorm:"uniqueIndex" json:"hold_id"`

	// Hold Information
	HoldType         string `json:"hold_type"`   // litigation, investigation, regulatory, preservation
	HoldStatus       string `json:"hold_status"` // active, pending, released, expired
	HoldReason       string `gorm:"type:text" json:"hold_reason"`
	IssuingAuthority string `json:"issuing_authority"`

	// Court Order Details
	CourtOrderNumber  string     `json:"court_order_number"`
	CourtJurisdiction string     `json:"court_jurisdiction"`
	JudgeName         string     `json:"judge_name"`
	OrderDate         *time.Time `json:"order_date"`

	// Legal Case Association
	CaseNumber    string `json:"case_number"`
	CaseType      string `json:"case_type"` // civil, criminal, regulatory
	PlaintiffName string `json:"plaintiff_name"`
	DefendantName string `json:"defendant_name"`
	LegalCounsel  string `json:"legal_counsel"`

	// Hold Duration
	HoldStartDate    time.Time  `json:"hold_start_date"`
	HoldEndDate      *time.Time `json:"hold_end_date"`
	ExpectedDuration int        `json:"expected_duration"` // days
	ExtensionCount   int        `json:"extension_count"`
	AutoRenew        bool       `json:"auto_renew"`

	// Evidence Preservation
	EvidencePreserved  bool       `json:"evidence_preserved"`
	PreservationMethod string     `json:"preservation_method"`
	DataSnapshot       string     `json:"data_snapshot"` // reference to backup
	PreservationDate   *time.Time `json:"preservation_date"`

	// Release Authorization
	ReleaseAuthorized   bool       `json:"release_authorized"`
	ReleaseAuthorizedBy string     `json:"release_authorized_by"`
	ReleaseDate         *time.Time `json:"release_date"`
	ReleaseReason       string     `json:"release_reason"`
	ReleaseDocuments    string     `gorm:"type:json" json:"release_documents"` // JSON array

	// Chain of Custody
	CustodyLog       string     `gorm:"type:json" json:"custody_log"` // JSON array
	CurrentCustodian string     `json:"current_custodian"`
	CustodyTransfers int        `json:"custody_transfers"`
	LastTransferDate *time.Time `json:"last_transfer_date"`

	// Documentation
	LegalDocuments    string `gorm:"type:json" json:"legal_documents"`    // JSON array of URLs
	Correspondence    string `gorm:"type:json" json:"correspondence"`     // JSON array
	NotificationsSent string `gorm:"type:json" json:"notifications_sent"` // JSON array

	// Risk & Cost
	LitigationRisk     string  `json:"litigation_risk"` // low, medium, high, critical
	EstimatedLiability float64 `json:"estimated_liability"`
	LegalCosts         float64 `json:"legal_costs"`
	ComplianceCosts    float64 `json:"compliance_costs"`
}

// DeviceExportControls manages international trade compliance
type DeviceExportControls struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Trade Compliance Status
	ExportStatus     string `json:"export_status"`     // approved, restricted, prohibited, pending
	ImportStatus     string `json:"import_status"`     // approved, restricted, prohibited, pending
	ComplianceStatus string `json:"compliance_status"` // compliant, non-compliant, under_review

	// License Requirements
	ExportLicenseRequired bool       `json:"export_license_required"`
	ExportLicenseNumber   string     `json:"export_license_number"`
	LicenseExpiryDate     *time.Time `json:"license_expiry_date"`
	LicenseType           string     `json:"license_type"`
	IssuingAgency         string     `json:"issuing_agency"`

	// Embargo & Sanctions
	EmbargoStatus       bool   `json:"embargo_status"`
	SanctionedCountries string `gorm:"type:json" json:"sanctioned_countries"` // JSON array
	EntityListStatus    bool   `json:"entity_list_status"`
	DeniedParties       string `gorm:"type:json" json:"denied_parties"` // JSON array

	// Classification
	ECCNNumber            string `json:"eccn_number"`    // Export Control Classification Number
	HSTariffCode          string `json:"hs_tariff_code"` // Harmonized System code
	DualUseClassification string `json:"dual_use_classification"`
	MilitaryEndUse        bool   `json:"military_end_use"`
	EncryptionStatus      string `json:"encryption_status"`

	// Country Restrictions
	RestrictedCountries string `gorm:"type:json" json:"restricted_countries"` // JSON array
	ApprovedCountries   string `gorm:"type:json" json:"approved_countries"`   // JSON array
	CountryLicenses     string `gorm:"type:json" json:"country_licenses"`     // JSON object

	// Technology Transfer
	TechnologyTransfer   bool   `json:"technology_transfer"`
	TransferRestrictions string `gorm:"type:json" json:"transfer_restrictions"` // JSON array
	ITARControlled       bool   `json:"itar_controlled"`                        // International Traffic in Arms
	EARControlled        bool   `json:"ear_controlled"`                         // Export Administration Regulations

	// Customs Documentation
	CustomsDeclarations  string `gorm:"type:json" json:"customs_declarations"`   // JSON array
	ShippingDocuments    string `gorm:"type:json" json:"shipping_documents"`     // JSON array
	CertificatesOfOrigin string `gorm:"type:json" json:"certificates_of_origin"` // JSON array

	// Trade History
	ExportHistory string `gorm:"type:json" json:"export_history"` // JSON array
	ImportHistory string `gorm:"type:json" json:"import_history"` // JSON array
	TotalExports  int    `json:"total_exports"`
	TotalImports  int    `json:"total_imports"`

	// Compliance Violations
	ViolationCount    int     `json:"violation_count"`
	ActiveViolations  string  `gorm:"type:json" json:"active_violations"` // JSON array
	PenaltiesImposed  float64 `json:"penalties_imposed"`
	CorrectiveActions string  `gorm:"type:json" json:"corrective_actions"` // JSON array
}

// DeviceDataPrivacy manages privacy compliance and data protection
type DeviceDataPrivacy struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID   uuid.UUID `gorm:"type:uuid;index" json:"user_id"`

	// GDPR Compliance
	GDPRCompliant   bool       `json:"gdpr_compliant"`
	GDPRConsentDate *time.Time `json:"gdpr_consent_date"`
	GDPRBasis       string     `json:"gdpr_basis"` // consent, contract, legal, vital, public, legitimate
	DataController  string     `json:"data_controller"`
	DataProcessor   string     `json:"data_processor"`

	// CCPA Compliance
	CCPACompliant      bool       `json:"ccpa_compliant"`
	CCPAOptOut         bool       `json:"ccpa_opt_out"`
	SaleOfData         bool       `json:"sale_of_data"`
	CCPADisclosureDate *time.Time `json:"ccpa_disclosure_date"`

	// Privacy Settings
	DataCollectionEnabled bool `json:"data_collection_enabled"`
	AnalyticsEnabled      bool `json:"analytics_enabled"`
	MarketingEnabled      bool `json:"marketing_enabled"`
	ThirdPartySharing     bool `json:"third_party_sharing"`
	LocationTracking      bool `json:"location_tracking"`

	// Data Retention
	RetentionPeriod   int        `json:"retention_period"` // days
	RetentionPolicy   string     `json:"retention_policy"`
	DataMinimization  bool       `json:"data_minimization"`
	AutoDeleteEnabled bool       `json:"auto_delete_enabled"`
	LastPurgeDate     *time.Time `json:"last_purge_date"`

	// Consent Management
	ConsentVersion   string    `json:"consent_version"`
	ConsentGiven     bool      `json:"consent_given"`
	ConsentTimestamp time.Time `json:"consent_timestamp"`
	ConsentWithdrawn bool      `json:"consent_withdrawn"`
	ConsentHistory   string    `gorm:"type:json" json:"consent_history"` // JSON array

	// Data Breaches
	BreachCount         int        `json:"breach_count"`
	LastBreachDate      *time.Time `json:"last_breach_date"`
	BreachNotifications string     `gorm:"type:json" json:"breach_notifications"` // JSON array
	DataCompromised     string     `gorm:"type:json" json:"data_compromised"`     // JSON array

	// Privacy Impact Assessment
	PIARequired        bool       `json:"pia_required"`
	PIACompleted       bool       `json:"pia_completed"`
	PIADate            *time.Time `json:"pia_date"`
	PIARiskLevel       string     `json:"pia_risk_level"`                       // low, medium, high
	PIARecommendations string     `gorm:"type:json" json:"pia_recommendations"` // JSON array

	// Third-Party Sharing
	DataRecipients       string `gorm:"type:json" json:"data_recipients"` // JSON array
	DataTransfers        string `gorm:"type:json" json:"data_transfers"`  // JSON array
	CrossBorderTransfers bool   `json:"cross_border_transfers"`
	BCRsInPlace          bool   `json:"bcrs_in_place"` // Binding Corporate Rules
	SCCsUsed             bool   `json:"sccs_used"`     // Standard Contractual Clauses

	// User Rights Requests
	AccessRequests      int `json:"access_requests"`
	DeletionRequests    int `json:"deletion_requests"`
	PortabilityRequests int `json:"portability_requests"`
	CorrectionRequests  int `json:"correction_requests"`
	ObjectionRequests   int `json:"objection_requests"`

	// Violations & Penalties
	PrivacyViolations   int     `json:"privacy_violations"`
	RegulatoryFines     float64 `json:"regulatory_fines"`
	ComplianceWarnings  string  `gorm:"type:json" json:"compliance_warnings"` // JSON array
	RemediationRequired bool    `json:"remediation_required"`
}

// DeviceRegulatoryReporting manages mandatory reporting requirements
type DeviceRegulatoryReporting struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Reporting Requirements
	ReportingRequired  bool   `json:"reporting_required"`
	ReportingFrequency string `json:"reporting_frequency"`                 // daily, weekly, monthly, quarterly, annual
	ReportingAgencies  string `gorm:"type:json" json:"reporting_agencies"` // JSON array
	ReportTypes        string `gorm:"type:json" json:"report_types"`       // JSON array

	// Submission Tracking
	LastSubmissionDate *time.Time `json:"last_submission_date"`
	NextSubmissionDate time.Time  `json:"next_submission_date"`
	SubmissionsCount   int        `json:"submissions_count"`
	SubmissionHistory  string     `gorm:"type:json" json:"submission_history"` // JSON array

	// Incident Reporting
	IncidentReports    string `gorm:"type:json" json:"incident_reports"`    // JSON array
	MandatoryIncidents string `gorm:"type:json" json:"mandatory_incidents"` // JSON array
	VoluntaryReports   int    `json:"voluntary_reports"`
	IncidentDeadlines  string `gorm:"type:json" json:"incident_deadlines"` // JSON object

	// Compliance Deadlines
	UpcomingDeadlines  string `gorm:"type:json" json:"upcoming_deadlines"` // JSON array
	MissedDeadlines    int    `json:"missed_deadlines"`
	DeadlineExtensions int    `json:"deadline_extensions"`
	CriticalDeadlines  string `gorm:"type:json" json:"critical_deadlines"` // JSON array

	// Reporting Metrics
	ReportingAccuracy float64 `json:"reporting_accuracy"` // percentage
	TimelinessScore   float64 `json:"timeliness_score"`   // 0-100
	CompletenessScore float64 `json:"completeness_score"` // 0-100
	QualityScore      float64 `json:"quality_score"`      // 0-100

	// Penalties & Fines
	PenaltiesIncurred float64 `json:"penalties_incurred"`
	PenaltyReasons    string  `gorm:"type:json" json:"penalty_reasons"` // JSON array
	AppealsFiled      int     `json:"appeals_filed"`
	AppealsWon        int     `json:"appeals_won"`

	// Corrective Actions
	CorrectiveActions    string `gorm:"type:json" json:"corrective_actions"` // JSON array
	ActionDeadlines      string `gorm:"type:json" json:"action_deadlines"`   // JSON object
	ActionStatus         string `gorm:"type:json" json:"action_status"`      // JSON object
	VerificationRequired bool   `json:"verification_required"`

	// Regulatory Correspondence
	InboundCorrespondence  string `gorm:"type:json" json:"inbound_correspondence"`  // JSON array
	OutboundCorrespondence string `gorm:"type:json" json:"outbound_correspondence"` // JSON array
	OpenInquiries          int    `json:"open_inquiries"`
	ResponseDeadlines      string `gorm:"type:json" json:"response_deadlines"` // JSON object

	// Audit Findings
	AuditFindings     string `gorm:"type:json" json:"audit_findings"`     // JSON array
	FindingSeverity   string `gorm:"type:json" json:"finding_severity"`   // JSON object
	RemediationStatus string `gorm:"type:json" json:"remediation_status"` // JSON object
	FollowUpRequired  bool   `json:"follow_up_required"`

	// Certificates & Attestations
	ComplianceCertificates string `gorm:"type:json" json:"compliance_certificates"` // JSON array
	CertificateExpiry      string `gorm:"type:json" json:"certificate_expiry"`      // JSON object
	AttestationsDue        string `gorm:"type:json" json:"attestations_due"`        // JSON array
}

// DeviceSecurityCompliance tracks security standard adherence
type DeviceSecurityCompliance struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Security Standards
	ISO27001Compliant bool    `json:"iso27001_compliant"`
	NISTCompliant     bool    `json:"nist_compliant"`
	CISCompliant      bool    `json:"cis_compliant"`
	OWASPCompliant    bool    `json:"owasp_compliant"`
	ComplianceScore   float64 `json:"compliance_score"` // 0-100

	// Vulnerability Assessment
	LastAssessmentDate   *time.Time `json:"last_assessment_date"`
	VulnerabilitiesFound int        `json:"vulnerabilities_found"`
	CriticalVulns        int        `json:"critical_vulns"`
	HighVulns            int        `json:"high_vulns"`
	MediumVulns          int        `json:"medium_vulns"`
	LowVulns             int        `json:"low_vulns"`

	// Patch Management
	PatchCompliance float64    `json:"patch_compliance"` // percentage
	PendingPatches  int        `json:"pending_patches"`
	CriticalPatches int        `json:"critical_patches"`
	LastPatchDate   *time.Time `json:"last_patch_date"`
	PatchWindow     string     `json:"patch_window"`

	// Access Control
	AccessControlCompliant bool       `json:"access_control_compliant"`
	MFAEnabled             bool       `json:"mfa_enabled"`
	RBACImplemented        bool       `json:"rbac_implemented"` // Role-Based Access Control
	PrivilegedAccounts     int        `json:"privileged_accounts"`
	LastAccessReview       *time.Time `json:"last_access_review"`

	// Encryption Compliance
	EncryptionAtRest    bool       `json:"encryption_at_rest"`
	EncryptionInTransit bool       `json:"encryption_in_transit"`
	EncryptionAlgorithm string     `json:"encryption_algorithm"`
	KeyManagement       string     `json:"key_management"`
	CertificateExpiry   *time.Time `json:"certificate_expiry"`

	// Authentication
	PasswordPolicy     string `json:"password_policy"`
	PasswordComplexity bool   `json:"password_complexity"`
	SessionTimeout     int    `json:"session_timeout"` // minutes
	AccountLockout     bool   `json:"account_lockout"`
	BiometricAuth      bool   `json:"biometric_auth"`

	// Security Audits
	LastAuditDate       *time.Time `json:"last_audit_date"`
	AuditFindings       string     `gorm:"type:json" json:"audit_findings"` // JSON array
	RemediationRequired bool       `json:"remediation_required"`
	RemediationDeadline *time.Time `json:"remediation_deadline"`
	AuditScore          float64    `json:"audit_score"`

	// Penetration Testing
	LastPenTestDate   *time.Time `json:"last_pentest_date"`
	PenTestFindings   string     `gorm:"type:json" json:"pentest_findings"` // JSON array
	ExploitableVulns  int        `json:"exploitable_vulns"`
	RemediationStatus string     `json:"remediation_status"`

	// Security Incidents
	IncidentCount    int        `json:"incident_count"`
	LastIncidentDate *time.Time `json:"last_incident_date"`
	IncidentHistory  string     `gorm:"type:json" json:"incident_history"` // JSON array
	MTTR             float64    `json:"mttr"`                              // Mean Time To Respond (hours)
	MTTD             float64    `json:"mttd"`                              // Mean Time To Detect (hours)

	// Compliance Tracking
	NonCompliantAreas string  `gorm:"type:json" json:"non_compliant_areas"` // JSON array
	ComplianceGaps    string  `gorm:"type:json" json:"compliance_gaps"`     // JSON array
	RemediationPlan   string  `gorm:"type:json" json:"remediation_plan"`    // JSON object
	ComplianceCost    float64 `json:"compliance_cost"`
}

// Methods for DeviceComplianceStatus
func (dcs *DeviceComplianceStatus) IsCompliant() bool {
	return dcs.OverallStatus == "compliant" && dcs.ComplianceScore >= 80
}

func (dcs *DeviceComplianceStatus) HasViolations() bool {
	return dcs.NonComplianceCount > 0 || dcs.ActiveViolations != "[]"
}

func (dcs *DeviceComplianceStatus) NeedsRemediation() bool {
	return dcs.RemediationRequired && dcs.RemediationStatus != "completed"
}

func (dcs *DeviceComplianceStatus) CalculateTotalCost() float64 {
	dcs.TotalComplianceCost = dcs.ComplianceCost + dcs.PenaltyCost + dcs.AuditCost + dcs.CertificationCost + dcs.RemediationCost
	return dcs.TotalComplianceCost
}

func (dcs *DeviceComplianceStatus) IsAuditDue() bool {
	if dcs.NextAuditDate != nil {
		return time.Since(*dcs.NextAuditDate) > 0
	}
	return true
}

// Methods for DeviceLegalHolds
func (dlh *DeviceLegalHolds) IsActive() bool {
	return dlh.HoldStatus == "active"
}

func (dlh *DeviceLegalHolds) CanRelease() bool {
	return dlh.ReleaseAuthorized && !dlh.IsActive()
}

func (dlh *DeviceLegalHolds) GetHoldDuration() int {
	if dlh.HoldEndDate != nil {
		return int(dlh.HoldEndDate.Sub(dlh.HoldStartDate).Hours() / 24)
	}
	return int(time.Since(dlh.HoldStartDate).Hours() / 24)
}

func (dlh *DeviceLegalHolds) IsHighRisk() bool {
	return dlh.LitigationRisk == "high" || dlh.LitigationRisk == "critical"
}

func (dlh *DeviceLegalHolds) GetTotalCost() float64 {
	return dlh.LegalCosts + dlh.ComplianceCosts
}

// Methods for DeviceExportControls
func (dec *DeviceExportControls) IsExportApproved() bool {
	return dec.ExportStatus == "approved" && (!dec.ExportLicenseRequired || dec.ExportLicenseNumber != "")
}

func (dec *DeviceExportControls) IsRestricted() bool {
	return dec.ExportStatus == "restricted" || dec.ImportStatus == "restricted" ||
		dec.EmbargoStatus || dec.EntityListStatus
}

func (dec *DeviceExportControls) HasViolations() bool {
	return dec.ViolationCount > 0
}

func (dec *DeviceExportControls) RequiresSpecialLicense() bool {
	return dec.DualUseClassification != "" || dec.MilitaryEndUse ||
		dec.ITARControlled || dec.EARControlled
}

func (dec *DeviceExportControls) IsLicenseExpired() bool {
	if dec.LicenseExpiryDate != nil {
		return time.Since(*dec.LicenseExpiryDate) > 0
	}
	return false
}

// Methods for DeviceDataPrivacy
func (ddp *DeviceDataPrivacy) IsPrivacyCompliant() bool {
	return ddp.GDPRCompliant && ddp.CCPACompliant && ddp.ConsentGiven
}

func (ddp *DeviceDataPrivacy) HasDataBreach() bool {
	return ddp.BreachCount > 0
}

func (ddp *DeviceDataPrivacy) RequiresPIA() bool {
	return ddp.PIARequired && !ddp.PIACompleted
}

func (ddp *DeviceDataPrivacy) GetTotalRequests() int {
	return ddp.AccessRequests + ddp.DeletionRequests + ddp.PortabilityRequests +
		ddp.CorrectionRequests + ddp.ObjectionRequests
}

func (ddp *DeviceDataPrivacy) IsHighRisk() bool {
	return ddp.PIARiskLevel == "high" || ddp.BreachCount > 0 || ddp.PrivacyViolations > 0
}

// Methods for DeviceRegulatoryReporting
func (drr *DeviceRegulatoryReporting) HasMissedDeadlines() bool {
	return drr.MissedDeadlines > 0
}

func (drr *DeviceRegulatoryReporting) IsReportingDue() bool {
	return time.Since(drr.NextSubmissionDate) > 0
}

func (drr *DeviceRegulatoryReporting) GetComplianceScore() float64 {
	score := (drr.ReportingAccuracy + drr.TimelinessScore + drr.CompletenessScore + drr.QualityScore) / 4
	if drr.MissedDeadlines > 0 {
		score *= 0.8 // Penalty for missed deadlines
	}
	return score
}

func (drr *DeviceRegulatoryReporting) HasOpenInquiries() bool {
	return drr.OpenInquiries > 0
}

func (drr *DeviceRegulatoryReporting) GetAppealsSuccessRate() float64 {
	if drr.AppealsFiled > 0 {
		return float64(drr.AppealsWon) / float64(drr.AppealsFiled) * 100
	}
	return 0
}

// Methods for DeviceSecurityCompliance
func (dsc *DeviceSecurityCompliance) IsSecurityCompliant() bool {
	return dsc.ComplianceScore >= 80 && dsc.CriticalVulns == 0
}

func (dsc *DeviceSecurityCompliance) HasCriticalVulnerabilities() bool {
	return dsc.CriticalVulns > 0
}

func (dsc *DeviceSecurityCompliance) NeedsPatchUpdate() bool {
	return dsc.PendingPatches > 0 || dsc.CriticalPatches > 0
}

func (dsc *DeviceSecurityCompliance) IsEncryptionCompliant() bool {
	return dsc.EncryptionAtRest && dsc.EncryptionInTransit
}

func (dsc *DeviceSecurityCompliance) GetVulnerabilityScore() float64 {
	totalVulns := float64(dsc.CriticalVulns*4 + dsc.HighVulns*3 + dsc.MediumVulns*2 + dsc.LowVulns)
	if totalVulns > 0 {
		return 100 - (totalVulns * 2) // Deduct points based on vulnerabilities
	}
	return 100
}

func (dsc *DeviceSecurityCompliance) IsCertificateExpiring() bool {
	if dsc.CertificateExpiry != nil {
		return time.Until(*dsc.CertificateExpiry) < 30*24*time.Hour // Within 30 days
	}
	return false
}

func (dsc *DeviceSecurityCompliance) HasSecurityIncidents() bool {
	return dsc.IncidentCount > 0
}
