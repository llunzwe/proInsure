package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserComplianceProfile manages comprehensive compliance and regulatory data
type UserComplianceProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// KYC Status
	KYCStatus             string     `gorm:"type:varchar(20)" json:"kyc_status"`
	KYCLevel              string     `gorm:"type:varchar(20)" json:"kyc_level"`
	KYCCompletionDate     *time.Time `json:"kyc_completion_date"`
	KYCExpiryDate         *time.Time `json:"kyc_expiry_date"`
	KYCProvider           string     `gorm:"type:varchar(50)" json:"kyc_provider"`
	KYCReferenceNumber    string     `gorm:"type:varchar(100)" json:"kyc_reference_number"`
	KYCDocuments          []string   `gorm:"type:json" json:"kyc_documents"`
	KYCVerificationMethod string     `gorm:"type:varchar(50)" json:"kyc_verification_method"`
	KYCScore              float64    `gorm:"default:0" json:"kyc_score"`

	// AML Status
	AMLStatus                 string                   `gorm:"type:varchar(20)" json:"aml_status"`
	AMLCheckDate              *time.Time               `json:"aml_check_date"`
	AMLRiskLevel              string                   `gorm:"type:varchar(20)" json:"aml_risk_level"`
	AMLProvider               string                   `gorm:"type:varchar(50)" json:"aml_provider"`
	AMLReportID               string                   `gorm:"type:varchar(100)" json:"aml_report_id"`
	AMLFlags                  []string                 `gorm:"type:json" json:"aml_flags"`
	SuspiciousActivityReports []map[string]interface{} `gorm:"type:json" json:"suspicious_activity_reports"`

	// Sanctions & PEP
	SanctionsScreened    bool                     `gorm:"default:false" json:"sanctions_screened"`
	SanctionsCheckDate   *time.Time               `json:"sanctions_check_date"`
	SanctionListMatches  []string                 `gorm:"type:json" json:"sanction_list_matches"`
	PEPStatus            bool                     `gorm:"default:false" json:"pep_status"`
	PEPCheckDate         *time.Time               `json:"pep_check_date"`
	PEPRelationships     []map[string]interface{} `gorm:"type:json" json:"pep_relationships"`
	AdverseMediaCheck    bool                     `gorm:"default:false" json:"adverse_media_check"`
	AdverseMediaFindings []string                 `gorm:"type:json" json:"adverse_media_findings"`

	// Regulatory Compliance
	RegulatoryJurisdictions []string          `gorm:"type:json" json:"regulatory_jurisdictions"`
	ComplianceStatus        map[string]string `gorm:"type:json" json:"compliance_status"`
	LicenseRequirements     map[string]bool   `gorm:"type:json" json:"license_requirements"`
	RegulatoryRestrictions  []string          `gorm:"type:json" json:"regulatory_restrictions"`
	CrossBorderRestrictions []string          `gorm:"type:json" json:"cross_border_restrictions"`

	// Data Protection
	GDPRConsent              bool       `gorm:"default:false" json:"gdpr_consent"`
	GDPRConsentDate          *time.Time `json:"gdpr_consent_date"`
	CCPAConsent              bool       `gorm:"default:false" json:"ccpa_consent"`
	CCPAConsentDate          *time.Time `json:"ccpa_consent_date"`
	DataProtectionOfficer    string     `gorm:"type:varchar(100)" json:"data_protection_officer"`
	DataProcessingAgreements []string   `gorm:"type:json" json:"data_processing_agreements"`

	// Audit Trail
	ComplianceAudits   []map[string]interface{} `gorm:"type:json" json:"compliance_audits"`
	LastAuditDate      *time.Time               `json:"last_audit_date"`
	NextAuditDate      *time.Time               `json:"next_audit_date"`
	AuditFindings      []string                 `gorm:"type:json" json:"audit_findings"`
	RemediationActions []map[string]interface{} `gorm:"type:json" json:"remediation_actions"`

	// Risk Assessment
	ComplianceRiskScore  float64    `gorm:"default:0" json:"compliance_risk_score"`
	RiskFactors          []string   `gorm:"type:json" json:"risk_factors"`
	EnhancedDueDiligence bool       `gorm:"default:false" json:"enhanced_due_diligence"`
	EDDReason            string     `gorm:"type:text" json:"edd_reason"`
	EDDCompletionDate    *time.Time `json:"edd_completion_date"`
}

// UserRegulatoryFiling manages regulatory filings and reporting
type UserRegulatoryFiling struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Filing Details
	FilingID              string                 `gorm:"type:varchar(100);uniqueIndex" json:"filing_id"`
	FilingType            string                 `gorm:"type:varchar(50)" json:"filing_type"`
	RegulatoryBody        string                 `gorm:"type:varchar(100)" json:"regulatory_body"`
	Jurisdiction          string                 `gorm:"type:varchar(50)" json:"jurisdiction"`
	FilingDate            time.Time              `json:"filing_date"`
	DueDate               time.Time              `json:"due_date"`
	SubmissionDate        *time.Time             `json:"submission_date"`
	Status                string                 `gorm:"type:varchar(20)" json:"status"`
	ReferenceNumber       string                 `gorm:"type:varchar(100)" json:"reference_number"`
	DocumentLinks         []string               `gorm:"type:json" json:"document_links"`
	FilingData            map[string]interface{} `gorm:"type:json" json:"filing_data"`
	Acknowledgment        string                 `gorm:"type:text" json:"acknowledgment"`
	ComplianceCertificate string                 `gorm:"type:varchar(255)" json:"compliance_certificate"`
}

// UserAuditLog maintains comprehensive audit trail
type UserAuditLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	PerformedBy *uuid.UUID `gorm:"type:uuid" json:"performed_by"`

	// Audit Entry
	EventType        string                 `gorm:"type:varchar(50)" json:"event_type"`
	EventCategory    string                 `gorm:"type:varchar(50)" json:"event_category"`
	EventDescription string                 `gorm:"type:text" json:"event_description"`
	EventTimestamp   time.Time              `json:"event_timestamp"`
	EntityType       string                 `gorm:"type:varchar(50)" json:"entity_type"`
	EntityID         string                 `gorm:"type:varchar(100)" json:"entity_id"`
	OldValues        map[string]interface{} `gorm:"type:json" json:"old_values"`
	NewValues        map[string]interface{} `gorm:"type:json" json:"new_values"`
	IPAddress        string                 `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent        string                 `gorm:"type:varchar(255)" json:"user_agent"`
	SessionID        string                 `gorm:"type:varchar(100)" json:"session_id"`
	DeviceID         string                 `gorm:"type:varchar(100)" json:"device_id"`
	Location         map[string]interface{} `gorm:"type:json" json:"location"`
	RiskScore        float64                `gorm:"default:0" json:"risk_score"`
	Flagged          bool                   `gorm:"default:false" json:"flagged"`
	ReviewRequired   bool                   `gorm:"default:false" json:"review_required"`
	ReviewedBy       *uuid.UUID             `gorm:"type:uuid" json:"reviewed_by"`
	ReviewDate       *time.Time             `json:"review_date"`
	ComplianceImpact string                 `gorm:"type:text" json:"compliance_impact"`
}

// UserLegalHolds manages legal preservation requirements
type UserLegalHolds struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Legal Hold Details
	HoldID               string                   `gorm:"type:varchar(100);uniqueIndex" json:"hold_id"`
	HoldType             string                   `gorm:"type:varchar(50)" json:"hold_type"`
	CaseNumber           string                   `gorm:"type:varchar(100)" json:"case_number"`
	Court                string                   `gorm:"type:varchar(100)" json:"court"`
	Jurisdiction         string                   `gorm:"type:varchar(50)" json:"jurisdiction"`
	StartDate            time.Time                `json:"start_date"`
	EndDate              *time.Time               `json:"end_date"`
	Status               string                   `gorm:"type:varchar(20)" json:"status"`
	Custodian            string                   `gorm:"type:varchar(100)" json:"custodian"`
	LegalCounsel         string                   `gorm:"type:varchar(100)" json:"legal_counsel"`
	DataScope            []string                 `gorm:"type:json" json:"data_scope"`
	PreservationMethod   string                   `gorm:"type:varchar(50)" json:"preservation_method"`
	DataLocation         []string                 `gorm:"type:json" json:"data_location"`
	AccessRestrictions   []string                 `gorm:"type:json" json:"access_restrictions"`
	ChainOfCustody       []map[string]interface{} `gorm:"type:json" json:"chain_of_custody"`
	ReleaseAuthorization string                   `gorm:"type:text" json:"release_authorization"`
	ComplianceNotes      string                   `gorm:"type:text" json:"compliance_notes"`
}

// UserDataRetention manages data retention policies
type UserDataRetention struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Retention Policy
	RetentionPolicy string                 `gorm:"type:varchar(100)" json:"retention_policy"`
	RetentionPeriod int                    `json:"retention_period_days"`
	DataCategories  map[string]int         `gorm:"type:json" json:"data_categories"`
	RetentionRules  map[string]interface{} `gorm:"type:json" json:"retention_rules"`

	// Deletion Schedule
	DeletionScheduled  bool       `gorm:"default:false" json:"deletion_scheduled"`
	DeletionDate       *time.Time `json:"deletion_date"`
	DeletionMethod     string     `gorm:"type:varchar(50)" json:"deletion_method"`
	DeletionScope      []string   `gorm:"type:json" json:"deletion_scope"`
	DeletionExclusions []string   `gorm:"type:json" json:"deletion_exclusions"`

	// Archive Management
	ArchiveEnabled    bool       `gorm:"default:false" json:"archive_enabled"`
	ArchiveLocation   string     `gorm:"type:varchar(255)" json:"archive_location"`
	ArchiveDate       *time.Time `json:"archive_date"`
	ArchiveRetention  int        `json:"archive_retention_days"`
	RestoreCapability bool       `gorm:"default:true" json:"restore_capability"`

	// Compliance
	ComplianceReasons      []string   `gorm:"type:json" json:"compliance_reasons"`
	RegulatoryRequirements []string   `gorm:"type:json" json:"regulatory_requirements"`
	ExemptionStatus        bool       `gorm:"default:false" json:"exemption_status"`
	ExemptionReason        string     `gorm:"type:text" json:"exemption_reason"`
	LastReviewDate         *time.Time `json:"last_review_date"`
}

// UserCrossBorderCompliance manages international compliance
type UserCrossBorderCompliance struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// International Compliance
	OperatingCountries   []string                 `gorm:"type:json" json:"operating_countries"`
	DataResidency        map[string]string        `gorm:"type:json" json:"data_residency"`
	CrossBorderTransfers []map[string]interface{} `gorm:"type:json" json:"cross_border_transfers"`
	TransferMechanisms   []string                 `gorm:"type:json" json:"transfer_mechanisms"`

	// Export Controls
	ExportControlStatus string                   `gorm:"type:varchar(20)" json:"export_control_status"`
	ExportLicenses      []map[string]interface{} `gorm:"type:json" json:"export_licenses"`
	RestrictedCountries []string                 `gorm:"type:json" json:"restricted_countries"`
	DualUseGoods        bool                     `gorm:"default:false" json:"dual_use_goods"`

	// Tax Compliance
	TaxResidencies   []string           `gorm:"type:json" json:"tax_residencies"`
	TaxIdentifiers   map[string]string  `gorm:"type:json" json:"tax_identifiers"`
	WithholdingRates map[string]float64 `gorm:"type:json" json:"withholding_rates"`
	TaxTreaties      []string           `gorm:"type:json" json:"tax_treaties"`
	FATCACompliance  bool               `gorm:"default:false" json:"fatca_compliance"`
	CRSReporting     bool               `gorm:"default:false" json:"crs_reporting"`

	// Currency Controls
	CurrencyRestrictions map[string]interface{} `gorm:"type:json" json:"currency_restrictions"`
	ExchangeControls     []string               `gorm:"type:json" json:"exchange_controls"`
	RepatriationRules    map[string]interface{} `gorm:"type:json" json:"repatriation_rules"`

	// Local Requirements
	LocalCompliance     map[string]interface{} `gorm:"type:json" json:"local_compliance"`
	LocalRepresentation map[string]interface{} `gorm:"type:json" json:"local_representation"`
	LocalFilings        map[string]interface{} `gorm:"type:json" json:"local_filings"`
}

// UserComplianceReporting manages compliance reporting
type UserComplianceReporting struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Reporting Configuration
	ReportingEnabled   bool     `gorm:"default:true" json:"reporting_enabled"`
	ReportingFrequency string   `gorm:"type:varchar(20)" json:"reporting_frequency"`
	ReportTypes        []string `gorm:"type:json" json:"report_types"`
	Recipients         []string `gorm:"type:json" json:"recipients"`

	// Report History
	GeneratedReports []map[string]interface{} `gorm:"type:json" json:"generated_reports"`
	LastReportDate   *time.Time               `json:"last_report_date"`
	NextReportDate   *time.Time               `json:"next_report_date"`

	// Metrics
	ComplianceScore   float64                `gorm:"default:0" json:"compliance_score"`
	ComplianceMetrics map[string]float64     `gorm:"type:json" json:"compliance_metrics"`
	TrendAnalysis     map[string]interface{} `gorm:"type:json" json:"trend_analysis"`
	Benchmarks        map[string]float64     `gorm:"type:json" json:"benchmarks"`

	// Certifications
	Certifications      []map[string]interface{} `gorm:"type:json" json:"certifications"`
	CertificationExpiry map[string]time.Time     `gorm:"type:json" json:"certification_expiry"`
	RenewalReminders    map[string]time.Time     `gorm:"type:json" json:"renewal_reminders"`
}

// TableName returns the table name
func (UserComplianceProfile) TableName() string {
	return "user_compliance_profiles"
}

// TableName returns the table name
func (UserRegulatoryFiling) TableName() string {
	return "user_regulatory_filings"
}

// TableName returns the table name
func (UserAuditLog) TableName() string {
	return "user_audit_logs"
}

// TableName returns the table name
func (UserLegalHolds) TableName() string {
	return "user_legal_holds"
}

// TableName returns the table name
func (UserDataRetention) TableName() string {
	return "user_data_retention"
}

// TableName returns the table name
func (UserCrossBorderCompliance) TableName() string {
	return "user_cross_border_compliance"
}

// TableName returns the table name
func (UserComplianceReporting) TableName() string {
	return "user_compliance_reporting"
}
