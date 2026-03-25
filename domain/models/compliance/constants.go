package compliance

import "time"

// ISO Standards
const (
	ISO27001 = "ISO/IEC 27001:2022"   // Information Security Management
	ISO27701 = "ISO/IEC 27701:2019"   // Privacy Information Management
	ISO31000 = "ISO 31000:2018"       // Risk Management
	ISO9001  = "ISO 9001:2015"        // Quality Management System
	ISO22301 = "ISO 22301:2019"       // Business Continuity Management
	ISO19600 = "ISO 19600:2014"       // Compliance Management
	ISO20000 = "ISO/IEC 20000-1:2018" // IT Service Management
	ISO27017 = "ISO/IEC 27017:2015"   // Cloud Security
	ISO27018 = "ISO/IEC 27018:2019"   // PII Protection in Clouds
	ISO37001 = "ISO 37001:2016"       // Anti-bribery Management
	ISO14001 = "ISO 14001:2015"       // Environmental Management
	ISO26000 = "ISO 26000:2010"       // Social Responsibility
	ISO38500 = "ISO/IEC 38500:2015"   // IT Governance
)

// Compliance Statuses
const (
	ComplianceStatusCompliant    = "compliant"
	ComplianceStatusNonCompliant = "non_compliant"
	ComplianceStatusPartial      = "partial"
	ComplianceStatusPending      = "pending_assessment"
	ComplianceStatusExpired      = "expired"
	ComplianceStatusRemediation  = "under_remediation"
	ComplianceStatusException    = "exception_granted"
)

// Risk Levels
const (
	RiskLevelCritical = "critical"
	RiskLevelHigh     = "high"
	RiskLevelMedium   = "medium"
	RiskLevelLow      = "low"
	RiskLevelMinimal  = "minimal"
)

// Risk Categories
const (
	RiskCategoryOperational   = "operational"
	RiskCategoryFinancial     = "financial"
	RiskCategoryCompliance    = "compliance"
	RiskCategoryStrategic     = "strategic"
	RiskCategoryReputational  = "reputational"
	RiskCategoryCyber         = "cyber"
	RiskCategoryLegal         = "legal"
	RiskCategoryEnvironmental = "environmental"
)

// Risk Treatment Options
const (
	RiskTreatmentAccept   = "accept"
	RiskTreatmentMitigate = "mitigate"
	RiskTreatmentTransfer = "transfer"
	RiskTreatmentAvoid    = "avoid"
	RiskTreatmentMonitor  = "monitor"
)

// Incident Types
const (
	IncidentTypeSecurity    = "security"
	IncidentTypePrivacy     = "privacy"
	IncidentTypeOperational = "operational"
	IncidentTypeCompliance  = "compliance"
	IncidentTypeData        = "data_breach"
	IncidentTypeSystem      = "system_failure"
	IncidentTypeFraud       = "fraud"
	IncidentTypeThirdParty  = "third_party"
)

// Incident Severities
const (
	IncidentSeverityCritical = "critical"
	IncidentSeverityHigh     = "high"
	IncidentSeverityMedium   = "medium"
	IncidentSeverityLow      = "low"
	IncidentSeverityInfo     = "informational"
)

// Incident Statuses
const (
	IncidentStatusOpen          = "open"
	IncidentStatusInvestigating = "investigating"
	IncidentStatusContained     = "contained"
	IncidentStatusEradicated    = "eradicated"
	IncidentStatusRecovering    = "recovering"
	IncidentStatusClosed        = "closed"
	IncidentStatusPostIncident  = "post_incident_review"
)

// Privacy Rights (GDPR/CCPA)
const (
	PrivacyRightAccess            = "right_to_access"
	PrivacyRightRectification     = "right_to_rectification"
	PrivacyRightErasure           = "right_to_erasure"
	PrivacyRightPortability       = "right_to_portability"
	PrivacyRightRestriction       = "right_to_restriction"
	PrivacyRightObject            = "right_to_object"
	PrivacyRightOptOut            = "right_to_opt_out"
	PrivacyRightNonDiscrimination = "right_to_non_discrimination"
)

// Consent Types
const (
	ConsentTypeProcessing      = "data_processing"
	ConsentTypeMarketing       = "marketing"
	ConsentTypeAnalytics       = "analytics"
	ConsentTypeThirdParty      = "third_party_sharing"
	ConsentTypeCookies         = "cookies"
	ConsentTypeLocation        = "location_tracking"
	ConsentTypeBiometric       = "biometric_data"
	ConsentTypeSpecialCategory = "special_category_data"
)

// Consent Statuses
const (
	ConsentStatusGranted   = "granted"
	ConsentStatusWithdrawn = "withdrawn"
	ConsentStatusExpired   = "expired"
	ConsentStatusPending   = "pending"
	ConsentStatusPartial   = "partial"
)

// Audit Types
const (
	AuditTypeInternal      = "internal"
	AuditTypeExternal      = "external"
	AuditTypeCertification = "certification"
	AuditTypeSurveillance  = "surveillance"
	AuditTypeCompliance    = "compliance"
	AuditTypeSecurity      = "security"
	AuditTypePrivacy       = "privacy"
	AuditTypeQuality       = "quality"
)

// Audit Statuses
const (
	AuditStatusScheduled  = "scheduled"
	AuditStatusInProgress = "in_progress"
	AuditStatusCompleted  = "completed"
	AuditStatusReporting  = "reporting"
	AuditStatusClosed     = "closed"
	AuditStatusCancelled  = "cancelled"
)

// Control Types
const (
	ControlTypePreventive   = "preventive"
	ControlTypeDetective    = "detective"
	ControlTypeCorrective   = "corrective"
	ControlTypeCompensating = "compensating"
	ControlTypeDirective    = "directive"
)

// Control Effectiveness
const (
	ControlEffectivenessEffective     = "effective"
	ControlEffectivenessPartial       = "partially_effective"
	ControlEffectivenessIneffective   = "ineffective"
	ControlEffectivenessNotTested     = "not_tested"
	ControlEffectivenessNotApplicable = "not_applicable"
)

// Business Continuity Statuses
const (
	BCStatusActive    = "active"
	BCStatusActivated = "activated"
	BCStatusTesting   = "testing"
	BCStatusUpdating  = "updating"
	BCStatusApproved  = "approved"
	BCStatusDraft     = "draft"
)

// Quality Metrics
const (
	QualityMetricCustomerSatisfaction = "customer_satisfaction"
	QualityMetricProcessEfficiency    = "process_efficiency"
	QualityMetricDefectRate           = "defect_rate"
	QualityMetricFirstCallResolution  = "first_call_resolution"
	QualityMetricResponseTime         = "response_time"
	QualityMetricComplianceRate       = "compliance_rate"
)

// Regulatory Bodies
const (
	RegulatoryGDPR   = "GDPR"    // EU General Data Protection Regulation
	RegulatoryCCPA   = "CCPA"    // California Consumer Privacy Act
	RegulatoryHIPAA  = "HIPAA"   // Health Insurance Portability and Accountability Act
	RegulatoryPCIDSS = "PCI-DSS" // Payment Card Industry Data Security Standard
	RegulatorySOX    = "SOX"     // Sarbanes-Oxley Act
	RegulatoryFCPA   = "FCPA"    // Foreign Corrupt Practices Act
	RegulatoryAML    = "AML"     // Anti-Money Laundering
	RegulatoryKYC    = "KYC"     // Know Your Customer
	RegulatoryLGPD   = "LGPD"    // Lei Geral de Proteção de Dados (Brazil)
	RegulatoryPIPEDA = "PIPEDA"  // Personal Information Protection and Electronic Documents Act (Canada)
	RegulatoryAPP    = "APP"     // Australian Privacy Principles
)

// Training Types
const (
	TrainingTypeSecurity     = "security_awareness"
	TrainingTypePrivacy      = "privacy_awareness"
	TrainingTypeCompliance   = "compliance"
	TrainingTypeIncident     = "incident_response"
	TrainingTypeBCP          = "business_continuity"
	TrainingTypeAntiTrust    = "anti_trust"
	TrainingTypePhishing     = "phishing_awareness"
	TrainingTypeDataHandling = "data_handling"
)

// Notification Types
const (
	NotificationTypeBreach        = "breach_notification"
	NotificationTypeIncident      = "incident_notification"
	NotificationTypeCompliance    = "compliance_update"
	NotificationTypeAudit         = "audit_notification"
	NotificationTypeCertification = "certification_update"
	NotificationTypeRegulatory    = "regulatory_change"
	NotificationTypeTraining      = "training_reminder"
)

// Time Limits (as per regulations)
const (
	GDPRBreachNotificationHours = 72 * time.Hour           // GDPR breach notification
	GDPRResponseDays            = 30 * 24 * time.Hour      // GDPR data subject request response
	CCPAResponseDays            = 45 * 24 * time.Hour      // CCPA consumer request response
	HIPAABreachNotificationDays = 60 * 24 * time.Hour      // HIPAA breach notification
	PCIDSSRetentionDays         = 365 * 24 * time.Hour     // PCI-DSS data retention
	DefaultRetentionPeriod      = 7 * 365 * 24 * time.Hour // 7 years default
)

// Compliance Thresholds
const (
	MinPasswordLength      = 12
	MaxLoginAttempts       = 5
	SessionTimeoutMinutes  = 30
	PasswordExpiryDays     = 90
	MFARequiredRiskScore   = 70
	MaxDataExportSizeMB    = 100
	EncryptionKeyLength    = 256
	BackupRetentionDays    = 90
	AuditLogRetentionYears = 7
)

// Evidence Types
const (
	EvidenceTypeDocument    = "document"
	EvidenceTypeScreenshot  = "screenshot"
	EvidenceTypeLog         = "log_file"
	EvidenceTypeReport      = "report"
	EvidenceTypeCertificate = "certificate"
	EvidenceTypeAttestation = "attestation"
	EvidenceTypeTestResult  = "test_result"
	EvidenceTypeApproval    = "approval"
)

// Compliance Frameworks
const (
	FrameworkNIST  = "NIST"         // National Institute of Standards and Technology
	FrameworkCOBIT = "COBIT"        // Control Objectives for Information Technologies
	FrameworkITIL  = "ITIL"         // Information Technology Infrastructure Library
	FrameworkCIS   = "CIS"          // Center for Internet Security
	FrameworkOWASP = "OWASP"        // Open Web Application Security Project
	FrameworkMITRE = "MITRE_ATT&CK" // MITRE Adversarial Tactics, Techniques & Common Knowledge
)

// Response Actions
const (
	ResponseActionInvestigate = "investigate"
	ResponseActionContain     = "contain"
	ResponseActionEradicate   = "eradicate"
	ResponseActionRecover     = "recover"
	ResponseActionReport      = "report"
	ResponseActionEscalate    = "escalate"
	ResponseActionIsolate     = "isolate"
	ResponseActionNotify      = "notify"
)

// Assessment Types
const (
	AssessmentTypeRisk          = "risk_assessment"
	AssessmentTypeVulnerability = "vulnerability_assessment"
	AssessmentTypePrivacy       = "privacy_impact_assessment"
	AssessmentTypeThirdParty    = "third_party_assessment"
	AssessmentTypeSecurity      = "security_assessment"
	AssessmentTypeCompliance    = "compliance_assessment"
	AssessmentTypePenetration   = "penetration_test"
)
