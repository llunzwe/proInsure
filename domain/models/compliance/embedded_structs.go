package compliance

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ComplianceIdentification contains core identification fields
type ComplianceIdentification struct {
	ComplianceNumber  string         `gorm:"uniqueIndex;not null" json:"compliance_number"`
	ComplianceType    string         `gorm:"index;not null" json:"compliance_type"` // ISO27001, GDPR, etc.
	Framework         string         `gorm:"index" json:"framework"`                // NIST, COBIT, etc.
	Version           string         `json:"version"`
	Title             string         `gorm:"not null" json:"title"`
	Description       string         `gorm:"type:text" json:"description"`
	Category          string         `gorm:"index" json:"category"` // security, privacy, quality, etc.
	SubCategory       string         `json:"sub_category"`
	ReferenceNumber   string         `gorm:"index" json:"reference_number"` // External reference
	InternalReference string         `json:"internal_reference"`
	Tags              pq.StringArray `gorm:"type:text[]" json:"tags"`
}

// ComplianceStatus contains status and lifecycle information
type ComplianceStatus struct {
	Status               string    `gorm:"index;not null;default:'pending_assessment'" json:"status"`
	ComplianceScore      float64   `json:"compliance_score"` // 0-100
	MaturityLevel        int       `json:"maturity_level"`   // 1-5
	EffectivenessRating  string    `json:"effectiveness_rating"`
	LastAssessmentDate   time.Time `json:"last_assessment_date"`
	NextAssessmentDate   time.Time `gorm:"index" json:"next_assessment_date"`
	CertificationDate    time.Time `json:"certification_date"`
	CertificationExpiry  time.Time `gorm:"index" json:"certification_expiry"`
	RemediationDeadline  time.Time `gorm:"index" json:"remediation_deadline"`
	IsActive             bool      `gorm:"default:true;index" json:"is_active"`
	IsCritical           bool      `gorm:"index" json:"is_critical"`
	RequiresAttestation  bool      `json:"requires_attestation"`
	AttestationCompleted bool      `json:"attestation_completed"`
	AttestationDate      time.Time `json:"attestation_date"`
	AttestedBy           uuid.UUID `json:"attested_by"`
}

// ComplianceRiskProfile contains risk assessment data
type ComplianceRiskProfile struct {
	RiskLevel            string         `gorm:"index" json:"risk_level"`
	RiskScore            float64        `json:"risk_score"` // 0-100
	InherentRisk         float64        `json:"inherent_risk"`
	ResidualRisk         float64        `json:"residual_risk"`
	RiskAppetite         float64        `json:"risk_appetite"`
	RiskTolerance        float64        `json:"risk_tolerance"`
	RiskCategory         string         `json:"risk_category"`
	RiskTreatment        string         `json:"risk_treatment"`
	RiskOwner            uuid.UUID      `json:"risk_owner"`
	ControlEffectiveness string         `json:"control_effectiveness"`
	ThreatLevel          string         `json:"threat_level"`
	VulnerabilityLevel   string         `json:"vulnerability_level"`
	ImpactLevel          string         `json:"impact_level"`
	LikelihoodLevel      string         `json:"likelihood_level"`
	RiskIndicators       pq.StringArray `gorm:"type:text[]" json:"risk_indicators"`
}

// ComplianceRegulatory contains regulatory and legal requirements
type ComplianceRegulatory struct {
	RegulatoryBodies       pq.StringArray       `gorm:"type:text[]" json:"regulatory_bodies"`
	Jurisdictions          pq.StringArray       `gorm:"type:text[]" json:"jurisdictions"`
	LegalRequirements      pq.StringArray       `gorm:"type:text[]" json:"legal_requirements"`
	IndustryStandards      pq.StringArray       `gorm:"type:text[]" json:"industry_standards"`
	ContractualObligations pq.StringArray       `gorm:"type:text[]" json:"contractual_obligations"`
	RegulatoryDeadlines    map[string]time.Time `gorm:"type:jsonb" json:"regulatory_deadlines"`
	PenaltyAmount          float64              `json:"penalty_amount"`
	PenaltyCurrency        string               `json:"penalty_currency"`
	EnforcementActions     int                  `json:"enforcement_actions"`
	RegulatoryNotices      int                  `json:"regulatory_notices"`
	ComplianceOfficer      uuid.UUID            `json:"compliance_officer"`
	LegalCounsel           uuid.UUID            `json:"legal_counsel"`
	LastRegulatoryReview   time.Time            `json:"last_regulatory_review"`
	NextRegulatoryReview   time.Time            `gorm:"index" json:"next_regulatory_review"`
}

// CompliancePrivacy contains privacy and data protection information
type CompliancePrivacy struct {
	DataCategories       pq.StringArray `gorm:"type:text[]" json:"data_categories"`
	ProcessingPurposes   pq.StringArray `gorm:"type:text[]" json:"processing_purposes"`
	LegalBasis           string         `json:"legal_basis"`
	DataRetentionPeriod  int            `json:"data_retention_period_days"`
	DataSubjectRights    pq.StringArray `gorm:"type:text[]" json:"data_subject_rights"`
	ConsentRequired      bool           `json:"consent_required"`
	ConsentObtained      bool           `json:"consent_obtained"`
	PIIPresent           bool           `gorm:"index" json:"pii_present"`
	SensitiveDataPresent bool           `gorm:"index" json:"sensitive_data_present"`
	CrossBorderTransfer  bool           `json:"cross_border_transfer"`
	TransferMechanisms   pq.StringArray `gorm:"type:text[]" json:"transfer_mechanisms"`
	PrivacyNoticeVersion string         `json:"privacy_notice_version"`
	LastPIADate          time.Time      `json:"last_pia_date"`
	NextPIADate          time.Time      `gorm:"index" json:"next_pia_date"`
	DPOApproval          bool           `json:"dpo_approval"`
	DPOApprovedBy        uuid.UUID      `json:"dpo_approved_by"`
}

// ComplianceAudit contains audit and assessment information
type ComplianceAudit struct {
	AuditType            string         `gorm:"index" json:"audit_type"`
	AuditScope           pq.StringArray `gorm:"type:text[]" json:"audit_scope"`
	AuditCriteria        pq.StringArray `gorm:"type:text[]" json:"audit_criteria"`
	AuditFindings        int            `json:"audit_findings"`
	CriticalFindings     int            `json:"critical_findings"`
	MajorFindings        int            `json:"major_findings"`
	MinorFindings        int            `json:"minor_findings"`
	Observations         int            `json:"observations"`
	NonConformities      int            `json:"non_conformities"`
	CorrectiveActions    int            `json:"corrective_actions"`
	PreventiveActions    int            `json:"preventive_actions"`
	AuditStatus          string         `gorm:"index" json:"audit_status"`
	AuditorName          string         `json:"auditor_name"`
	AuditorOrganization  string         `json:"auditor_organization"`
	AuditStartDate       time.Time      `json:"audit_start_date"`
	AuditEndDate         time.Time      `json:"audit_end_date"`
	AuditReportDate      time.Time      `json:"audit_report_date"`
	ManagementReviewDate time.Time      `json:"management_review_date"`
	NextAuditDate        time.Time      `gorm:"index" json:"next_audit_date"`
}

// ComplianceControl contains control measures and safeguards
type ComplianceControl struct {
	ControlType          string         `json:"control_type"`
	ControlCategory      string         `json:"control_category"`
	ControlObjectives    pq.StringArray `gorm:"type:text[]" json:"control_objectives"`
	ControlMeasures      pq.StringArray `gorm:"type:text[]" json:"control_measures"`
	ImplementationStatus string         `json:"implementation_status"`
	TestingFrequency     string         `json:"testing_frequency"`
	LastTestDate         time.Time      `json:"last_test_date"`
	NextTestDate         time.Time      `gorm:"index" json:"next_test_date"`
	TestResults          string         `json:"test_results"`
	Effectiveness        string         `json:"effectiveness"`
	AutomationLevel      string         `json:"automation_level"` // manual, semi-automated, automated
	ControlOwner         uuid.UUID      `json:"control_owner"`
	ControlReviewer      uuid.UUID      `json:"control_reviewer"`
	Compensating         bool           `json:"compensating"`
	CompensatingControls pq.StringArray `gorm:"type:text[]" json:"compensating_controls"`
}

// ComplianceEvidence contains documentation and proof of compliance
type ComplianceEvidence struct {
	EvidenceTypes       pq.StringArray `gorm:"type:text[]" json:"evidence_types"`
	DocumentIDs         pq.StringArray `gorm:"type:text[]" json:"document_ids"`
	EvidenceLocation    string         `json:"evidence_location"`
	RetentionPeriod     int            `json:"retention_period_years"`
	CollectionMethod    string         `json:"collection_method"`
	CollectionFrequency string         `json:"collection_frequency"`
	LastCollectionDate  time.Time      `json:"last_collection_date"`
	NextCollectionDate  time.Time      `gorm:"index" json:"next_collection_date"`
	EvidenceQuality     string         `json:"evidence_quality"` // high, medium, low
	VerificationStatus  string         `json:"verification_status"`
	VerifiedBy          uuid.UUID      `json:"verified_by"`
	VerificationDate    time.Time      `json:"verification_date"`
	DigitalSignature    string         `json:"digital_signature"`
	ChainOfCustody      bool           `json:"chain_of_custody"`
}

// ComplianceMonitoring contains continuous monitoring information
type ComplianceMonitoring struct {
	MonitoringEnabled   bool               `json:"monitoring_enabled"`
	MonitoringType      string             `json:"monitoring_type"` // continuous, periodic, event-driven
	MonitoringFrequency string             `json:"monitoring_frequency"`
	KPIs                pq.StringArray     `gorm:"type:text[]" json:"kpis"`
	KRIs                pq.StringArray     `gorm:"type:text[]" json:"kris"` // Key Risk Indicators
	Thresholds          map[string]float64 `gorm:"type:jsonb" json:"thresholds"`
	AlertsEnabled       bool               `json:"alerts_enabled"`
	AlertRecipients     pq.StringArray     `gorm:"type:text[]" json:"alert_recipients"`
	EscalationPath      pq.StringArray     `gorm:"type:text[]" json:"escalation_path"`
	LastMonitoringDate  time.Time          `json:"last_monitoring_date"`
	NextMonitoringDate  time.Time          `gorm:"index" json:"next_monitoring_date"`
	AutomatedChecks     int                `json:"automated_checks"`
	ManualChecks        int                `json:"manual_checks"`
	MonitoringDashboard string             `json:"monitoring_dashboard_url"`
	RealtimeMonitoring  bool               `json:"realtime_monitoring"`
}

// ComplianceReporting contains reporting and communication details
type ComplianceReporting struct {
	ReportingRequired   bool           `json:"reporting_required"`
	ReportingFrequency  string         `json:"reporting_frequency"`
	ReportTypes         pq.StringArray `gorm:"type:text[]" json:"report_types"`
	ReportRecipients    pq.StringArray `gorm:"type:text[]" json:"report_recipients"`
	InternalReporting   bool           `json:"internal_reporting"`
	ExternalReporting   bool           `json:"external_reporting"`
	RegulatoryReporting bool           `json:"regulatory_reporting"`
	LastReportDate      time.Time      `json:"last_report_date"`
	NextReportDate      time.Time      `gorm:"index" json:"next_report_date"`
	ReportFormat        string         `json:"report_format"` // PDF, XML, JSON, etc.
	ReportTemplate      string         `json:"report_template_id"`
	AutomatedReporting  bool           `json:"automated_reporting"`
	BoardReporting      bool           `json:"board_reporting"`
	PublicDisclosure    bool           `json:"public_disclosure"`
	MetricsIncluded     pq.StringArray `gorm:"type:text[]" json:"metrics_included"`
}

// ComplianceTraining contains training and awareness information
type ComplianceTraining struct {
	TrainingRequired      bool           `json:"training_required"`
	TrainingTypes         pq.StringArray `gorm:"type:text[]" json:"training_types"`
	TargetAudience        pq.StringArray `gorm:"type:text[]" json:"target_audience"`
	TrainingFrequency     string         `json:"training_frequency"`
	CompletionRate        float64        `json:"completion_rate"` // 0-100
	PassRate              float64        `json:"pass_rate"`       // 0-100
	AverageScore          float64        `json:"average_score"`
	LastTrainingDate      time.Time      `json:"last_training_date"`
	NextTrainingDate      time.Time      `gorm:"index" json:"next_training_date"`
	TrainingProvider      string         `json:"training_provider"`
	TrainingMethod        string         `json:"training_method"` // online, classroom, hybrid
	CertificationRequired bool           `json:"certification_required"`
	RefresherRequired     bool           `json:"refresher_required"`
	TrainingMaterials     pq.StringArray `gorm:"type:text[]" json:"training_materials"`
	TrainingRecords       int            `json:"training_records"`
}

// ComplianceRelationships contains relationships to other entities
type ComplianceRelationships struct {
	OrganizationID       uuid.UUID      `gorm:"index" json:"organization_id"`
	DepartmentID         uuid.UUID      `json:"department_id"`
	ProcessID            uuid.UUID      `json:"process_id"`
	SystemID             uuid.UUID      `json:"system_id"`
	VendorID             uuid.UUID      `json:"vendor_id"`
	ContractID           uuid.UUID      `json:"contract_id"`
	ProjectID            uuid.UUID      `json:"project_id"`
	PolicyIDs            pq.StringArray `gorm:"type:text[]" json:"policy_ids"`
	ProcedureIDs         pq.StringArray `gorm:"type:text[]" json:"procedure_ids"`
	RelatedCompliances   pq.StringArray `gorm:"type:text[]" json:"related_compliances"`
	DependentCompliances pq.StringArray `gorm:"type:text[]" json:"dependent_compliances"`
	AffectedUsers        int            `json:"affected_users"`
	AffectedSystems      int            `json:"affected_systems"`
	AffectedProcesses    int            `json:"affected_processes"`
}

// ComplianceMetadata contains additional metadata and custom fields
type ComplianceMetadata struct {
	Source         string                 `json:"source"`
	ImportedFrom   string                 `json:"imported_from"`
	ExternalID     string                 `gorm:"index" json:"external_id"`
	CustomFields   map[string]interface{} `gorm:"type:jsonb" json:"custom_fields"`
	Notes          string                 `gorm:"type:text" json:"notes"`
	InternalNotes  string                 `gorm:"type:text" json:"internal_notes"`
	Cost           float64                `json:"cost"`
	Budget         float64                `json:"budget"`
	ROI            float64                `json:"roi"`
	Priority       int                    `gorm:"index" json:"priority"`
	Visibility     string                 `json:"visibility"` // public, internal, confidential
	Classification string                 `json:"classification"`
	Language       string                 `json:"language"`
	TimeZone       string                 `json:"timezone"`
}
