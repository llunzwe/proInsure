package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// ComplianceRule represents regulatory compliance rules
type ComplianceRule struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	RuleCode           string         `gorm:"uniqueIndex;not null" json:"rule_code"`
	RuleName           string         `gorm:"not null" json:"rule_name"`
	RuleCategory       string         `gorm:"not null" json:"rule_category"`   // data_protection, financial, insurance, security
	Jurisdiction       string         `gorm:"not null" json:"jurisdiction"`    // US, EU, UK, ZA, etc.
	RegulatoryBody     string         `gorm:"not null" json:"regulatory_body"` // GDPR, CCPA, FSCA, FCA, etc.
	Description        string         `gorm:"not null" json:"description"`
	Requirements       string         `json:"requirements"`                     // JSON array of specific requirements
	ComplianceLevel    string         `gorm:"not null" json:"compliance_level"` // mandatory, recommended, optional
	EffectiveDate      time.Time      `gorm:"not null" json:"effective_date"`
	ExpiryDate         *time.Time     `json:"expiry_date"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	AutomationLevel    string         `gorm:"default:'manual'" json:"automation_level"` // manual, semi_automated, fully_automated
	CheckFrequency     string         `json:"check_frequency"`                          // daily, weekly, monthly, quarterly, annual
	PenaltyDescription string         `json:"penalty_description"`
	MaxPenaltyAmount   float64        `json:"max_penalty_amount"`
	LastReviewDate     *time.Time     `json:"last_review_date"`
	NextReviewDate     *time.Time     `json:"next_review_date"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ComplianceChecks []ComplianceCheck     `gorm:"foreignKey:RuleID" json:"compliance_checks,omitempty"`
	Violations       []ComplianceViolation `gorm:"foreignKey:RuleID" json:"violations,omitempty"`
}

// ComplianceCheck represents automated compliance checks
type ComplianceCheck struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RuleID             uuid.UUID  `gorm:"type:uuid;not null" json:"rule_id"`
	CheckID            string     `gorm:"uniqueIndex;not null" json:"check_id"`
	CheckType          string     `gorm:"not null" json:"check_type"` // policy_creation, claim_processing, data_access, payment
	EntityType         string     `json:"entity_type"`                // user, policy, claim, transaction
	EntityID           *uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	CheckDate          time.Time  `gorm:"not null" json:"check_date"`
	CheckStatus        string     `gorm:"not null;default:'pending'" json:"check_status"` // pending, passed, failed, warning
	ComplianceScore    float64    `json:"compliance_score"`                               // 0-100
	CheckResults       string     `json:"check_results"`                                  // JSON object with detailed results
	FailureReasons     string     `json:"failure_reasons"`                                // JSON array of failure reasons
	Recommendations    string     `json:"recommendations"`                                // JSON array of recommendations
	AutoRemediation    bool       `gorm:"default:false" json:"auto_remediation"`
	RemediationActions string     `json:"remediation_actions"` // JSON array of actions taken
	IsManualReview     bool       `gorm:"default:false" json:"is_manual_review"`
	ReviewedBy         *uuid.UUID `gorm:"type:uuid" json:"reviewed_by"`
	ReviewDate         *time.Time `json:"review_date"`
	ReviewNotes        string     `json:"review_notes"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Rule           ComplianceRule  `gorm:"foreignKey:RuleID" json:"rule,omitempty"`
	ReviewedByUser *models.User   `gorm:"foreignKey:ReviewedBy" json:"reviewed_by_user,omitempty"`
}

// ComplianceViolation represents compliance violations
type ComplianceViolation struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RuleID              uuid.UUID  `gorm:"type:uuid;not null" json:"rule_id"`
	ViolationID         string     `gorm:"uniqueIndex;not null" json:"violation_id"`
	ViolationType       string     `gorm:"not null" json:"violation_type"` // data_breach, unauthorized_access, policy_violation
	Severity            string     `gorm:"not null" json:"severity"`       // low, medium, high, critical
	DetectionDate       time.Time  `gorm:"not null" json:"detection_date"`
	DetectionMethod     string     `json:"detection_method"` // automated, manual, reported
	EntityType          string     `json:"entity_type"`
	EntityID            *uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	Description         string     `gorm:"not null" json:"description"`
	ImpactAssessment    string     `json:"impact_assessment"` // JSON object
	RootCause           string     `json:"root_cause"`
	Status              string     `gorm:"not null;default:'open'" json:"status"` // open, investigating, resolved, closed
	AssignedTo          *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	ResolutionDate      *time.Time `json:"resolution_date"`
	ResolutionActions   string     `json:"resolution_actions"` // JSON array
	PreventiveActions   string     `json:"preventive_actions"` // JSON array
	ReportedToRegulator bool       `gorm:"default:false" json:"reported_to_regulator"`
	ReportingDate       *time.Time `json:"reporting_date"`
	RegulatoryResponse  string     `json:"regulatory_response"`
	FinancialImpact     float64    `gorm:"default:0" json:"financial_impact"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Rule           ComplianceRule  `gorm:"foreignKey:RuleID" json:"rule,omitempty"`
	AssignedToUser *models.User    `gorm:"foreignKey:AssignedTo" json:"assigned_to_user,omitempty"`
}

// SanctionsScreening represents sanctions and watchlist screening
type SanctionsScreening struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ScreeningID       string     `gorm:"uniqueIndex;not null" json:"screening_id"`
	EntityType        string     `gorm:"not null" json:"entity_type"` // user, device, transaction
	EntityID          uuid.UUID  `gorm:"type:uuid;not null" json:"entity_id"`
	ScreeningType     string     `gorm:"not null" json:"screening_type"` // ofac, eu_sanctions, un_sanctions, pep, adverse_media
	ScreeningDate     time.Time  `gorm:"not null" json:"screening_date"`
	ScreeningProvider string     `json:"screening_provider"` // internal, worldcheck, dowjones, etc.
	SearchCriteria    string     `json:"search_criteria"`    // JSON object with search parameters
	MatchesFound      int        `gorm:"default:0" json:"matches_found"`
	HighestRiskScore  float64    `gorm:"default:0" json:"highest_risk_score"`         // 0-100
	ScreeningResults  string     `json:"screening_results"`                           // JSON array of matches
	Status            string     `gorm:"not null;default:'clear'" json:"status"`      // clear, potential_match, confirmed_match, false_positive
	ReviewStatus      string     `gorm:"default:'not_required'" json:"review_status"` // not_required, pending, reviewed, escalated
	ReviewedBy        *uuid.UUID `gorm:"type:uuid" json:"reviewed_by"`
	ReviewDate        *time.Time `json:"review_date"`
	ReviewDecision    string     `json:"review_decision"` // approve, reject, escalate
	ReviewNotes       string     `json:"review_notes"`
	IsBlocked         bool       `gorm:"default:false" json:"is_blocked"`
	BlockReason       string     `json:"block_reason"`
	UnblockDate       *time.Time `json:"unblock_date"`
	NextScreeningDate *time.Time `json:"next_screening_date"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ReviewedByUser *models.User `gorm:"foreignKey:ReviewedBy" json:"reviewed_by_user,omitempty"`
}

// DataProtectionRecord represents GDPR/data protection compliance
type DataProtectionRecord struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RecordID            string     `gorm:"uniqueIndex;not null" json:"record_id"`
	DataSubjectID       uuid.UUID  `gorm:"type:uuid;not null" json:"data_subject_id"` // User ID
	DataCategory        string     `gorm:"not null" json:"data_category"`             // personal, sensitive, financial, health
	DataType            string     `gorm:"not null" json:"data_type"`                 // name, email, phone, location, biometric
	ProcessingPurpose   string     `gorm:"not null" json:"processing_purpose"`        // insurance_contract, claims_processing, marketing
	LegalBasis          string     `gorm:"not null" json:"legal_basis"`               // consent, contract, legal_obligation, legitimate_interest
	ConsentDate         *time.Time `json:"consent_date"`
	ConsentWithdrawn    bool       `gorm:"default:false" json:"consent_withdrawn"`
	ConsentWithdrawDate *time.Time `json:"consent_withdraw_date"`
	RetentionPeriod     int        `json:"retention_period"` // days
	RetentionReason     string     `json:"retention_reason"`
	DeletionDate        *time.Time `json:"deletion_date"`
	IsDeleted           bool       `gorm:"default:false" json:"is_deleted"`
	DeletionMethod      string     `json:"deletion_method"` // soft_delete, hard_delete, anonymization
	DataLocation        string     `json:"data_location"`   // country/region where data is stored
	ThirdPartySharing   bool       `gorm:"default:false" json:"third_party_sharing"`
	ThirdParties        string     `json:"third_parties"` // JSON array of third parties
	IsEncrypted         bool       `gorm:"default:true" json:"is_encrypted"`
	EncryptionMethod    string     `json:"encryption_method"`
	AccessLog           string     `json:"access_log"` // JSON array of access records
	LastAccessed        *time.Time `json:"last_accessed"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	DataSubject    models.User         `gorm:"foreignKey:DataSubjectID" json:"data_subject,omitempty"`
	AccessRequests []DataAccessRequest `gorm:"foreignKey:RecordID" json:"access_requests,omitempty"`
}

// DataAccessRequest represents data subject access requests (GDPR Article 15)
type DataAccessRequest struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RequestID          string     `gorm:"uniqueIndex;not null" json:"request_id"`
	RecordID           *uuid.UUID `gorm:"type:uuid" json:"record_id"`
	DataSubjectID      uuid.UUID  `gorm:"type:uuid;not null" json:"data_subject_id"`
	RequestType        string     `gorm:"not null" json:"request_type"` // access, rectification, erasure, portability, restriction
	RequestDate        time.Time  `gorm:"not null" json:"request_date"`
	RequestChannel     string     `json:"request_channel"` // email, portal, phone, letter
	RequestDetails     string     `json:"request_details"`
	IdentityVerified   bool       `gorm:"default:false" json:"identity_verified"`
	VerificationMethod string     `json:"verification_method"`
	VerificationDate   *time.Time `json:"verification_date"`
	Status             string     `gorm:"not null;default:'received'" json:"status"` // received, verifying, processing, completed, rejected
	AssignedTo         *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	DueDate            time.Time  `gorm:"not null" json:"due_date"` // 30 days from request
	CompletionDate     *time.Time `json:"completion_date"`
	ResponseMethod     string     `json:"response_method"` // email, portal, secure_download
	ResponseData       string     `json:"response_data"`   // JSON object or file path
	RejectionReason    string     `json:"rejection_reason"`
	ProcessingNotes    string     `json:"processing_notes"`
	IsCompliant        bool       `gorm:"default:true" json:"is_compliant"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	DataSubject    models.User  `gorm:"foreignKey:DataSubjectID" json:"data_subject,omitempty"`
	AssignedToUser *models.User `gorm:"foreignKey:AssignedTo" json:"assigned_to_user,omitempty"`
}

// SecurityIncident represents security incidents and breaches
type SecurityIncident struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	IncidentID           string     `gorm:"uniqueIndex;not null" json:"incident_id"`
	IncidentType         string     `gorm:"not null" json:"incident_type"`         // data_breach, unauthorized_access, malware, phishing, ddos
	Severity             string     `gorm:"not null" json:"severity"`              // low, medium, high, critical
	Status               string     `gorm:"not null;default:'open'" json:"status"` // open, investigating, contained, resolved, closed
	DetectionDate        time.Time  `gorm:"not null" json:"detection_date"`
	DetectionMethod      string     `json:"detection_method"` // automated, manual, third_party, user_report
	ReportedBy           *uuid.UUID `gorm:"type:uuid" json:"reported_by"`
	AssignedTo           *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	Title                string     `gorm:"not null" json:"title"`
	Description          string     `gorm:"not null" json:"description"`
	AffectedSystems      string     `json:"affected_systems"` // JSON array
	AffectedUsers        string     `json:"affected_users"`   // JSON array of user IDs
	DataCompromised      bool       `gorm:"default:false" json:"data_compromised"`
	DataTypes            string     `json:"data_types"`       // JSON array of compromised data types
	EstimatedImpact      string     `json:"estimated_impact"` // JSON object
	RootCause            string     `json:"root_cause"`
	ContainmentActions   string     `json:"containment_actions"` // JSON array
	ContainmentDate      *time.Time `json:"containment_date"`
	EradicationActions   string     `json:"eradication_actions"` // JSON array
	EradicationDate      *time.Time `json:"eradication_date"`
	RecoveryActions      string     `json:"recovery_actions"` // JSON array
	RecoveryDate         *time.Time `json:"recovery_date"`
	LessonsLearned       string     `json:"lessons_learned"`
	PreventiveActions    string     `json:"preventive_actions"` // JSON array
	NotificationRequired bool       `gorm:"default:false" json:"notification_required"`
	NotificationDate     *time.Time `json:"notification_date"`
	NotificationMethod   string     `json:"notification_method"`
	RegulatoryReported   bool       `gorm:"default:false" json:"regulatory_reported"`
	RegulatoryReportDate *time.Time `json:"regulatory_report_date"`
	FinancialImpact      float64    `gorm:"default:0" json:"financial_impact"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ReportedByUser *models.User `gorm:"foreignKey:ReportedBy" json:"reported_by_user,omitempty"`
	AssignedToUser *models.User `gorm:"foreignKey:AssignedTo" json:"assigned_to_user,omitempty"`
	// Activities relationship is defined in IncidentActivity to avoid circular dependency
}

// IncidentActivity represents activities during incident response
type IncidentActivity struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	IncidentID   uuid.UUID `gorm:"type:uuid;not null" json:"incident_id"`
	ActivityType string    `gorm:"not null" json:"activity_type"` // investigation, containment, communication, recovery
	Description  string    `gorm:"not null" json:"description"`
	PerformedBy  uuid.UUID `gorm:"type:uuid;not null" json:"performed_by"`
	ActivityDate time.Time `gorm:"not null" json:"activity_date"`
	Duration     int       `json:"duration"`               // minutes
	Status       string    `gorm:"not null" json:"status"` // planned, in_progress, completed, cancelled
	Outcome      string    `json:"outcome"`
	Evidence     string    `json:"evidence"` // JSON array of evidence/documents
	NextActions  string    `json:"next_actions"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Incident        SecurityIncident `gorm:"foreignKey:IncidentID" json:"incident,omitempty"`
	PerformedByUser models.User      `gorm:"foreignKey:PerformedBy" json:"performed_by_user,omitempty"`
}

// SecurityAuditLog represents comprehensive security audit logging
type SecurityAuditLog struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	LogID           string     `gorm:"uniqueIndex;not null" json:"log_id"`
	UserID          *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	SessionID       string     `json:"session_id"`
	Action          string     `gorm:"not null" json:"action"`
	Resource        string     `json:"resource"`
	ResourceID      *uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	Method          string     `json:"method"` // GET, POST, PUT, DELETE
	Endpoint        string     `json:"endpoint"`
	IPAddress       string     `json:"ip_address"`
	UserAgent       string     `json:"user_agent"`
	RequestData     string     `json:"request_data"` // JSON object
	ResponseStatus  int        `json:"response_status"`
	ResponseData    string     `json:"response_data"` // JSON object
	Duration        int        `json:"duration"`      // milliseconds
	IsSuccessful    bool       `json:"is_successful"`
	ErrorMessage    string     `json:"error_message"`
	RiskScore       float64    `json:"risk_score"` // 0-100
	IsHighRisk      bool       `gorm:"default:false" json:"is_high_risk"`
	GeolocationData string     `json:"geolocation_data"` // JSON object
	DeviceInfo      string     `json:"device_info"`      // JSON object
	Timestamp       time.Time  `gorm:"not null" json:"timestamp"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User *models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName methods
func (ComplianceRule) TableName() string {
	return "compliance_rules"
}

func (ComplianceCheck) TableName() string {
	return "compliance_checks"
}

func (ComplianceViolation) TableName() string {
	return "compliance_violations"
}

func (SanctionsScreening) TableName() string {
	return "sanctions_screenings"
}

func (DataProtectionRecord) TableName() string {
	return "data_protection_records"
}

func (DataAccessRequest) TableName() string {
	return "data_access_requests"
}

func (SecurityIncident) TableName() string {
	return "security_incidents"
}

func (IncidentActivity) TableName() string {
	return "incident_activities"
}

func (SecurityAuditLog) TableName() string {
	return "security_audit_logs"
}

// BeforeCreate hooks
func (cr *ComplianceRule) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == uuid.Nil {
		cr.ID = uuid.New()
	}
	return nil
}

func (cc *ComplianceCheck) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	if cc.CheckID == "" {
		cc.CheckID = "CHK-" + uuid.New().String()[:8]
	}
	return nil
}

func (cv *ComplianceViolation) BeforeCreate(tx *gorm.DB) error {
	if cv.ID == uuid.Nil {
		cv.ID = uuid.New()
	}
	if cv.ViolationID == "" {
		cv.ViolationID = "VIO-" + uuid.New().String()[:8]
	}
	return nil
}

func (ss *SanctionsScreening) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	if ss.ScreeningID == "" {
		ss.ScreeningID = "SCR-" + uuid.New().String()[:8]
	}
	return nil
}

func (dpr *DataProtectionRecord) BeforeCreate(tx *gorm.DB) error {
	if dpr.ID == uuid.Nil {
		dpr.ID = uuid.New()
	}
	if dpr.RecordID == "" {
		dpr.RecordID = "DPR-" + uuid.New().String()[:8]
	}
	return nil
}

func (dar *DataAccessRequest) BeforeCreate(tx *gorm.DB) error {
	if dar.ID == uuid.Nil {
		dar.ID = uuid.New()
	}
	if dar.RequestID == "" {
		dar.RequestID = "DAR-" + uuid.New().String()[:8]
	}
	// Set due date to 30 days from request date (GDPR requirement)
	dar.DueDate = dar.RequestDate.AddDate(0, 0, 30)
	return nil
}

func (si *SecurityIncident) BeforeCreate(tx *gorm.DB) error {
	if si.ID == uuid.Nil {
		si.ID = uuid.New()
	}
	if si.IncidentID == "" {
		si.IncidentID = "INC-" + uuid.New().String()[:8]
	}
	return nil
}

func (ia *IncidentActivity) BeforeCreate(tx *gorm.DB) error {
	if ia.ID == uuid.Nil {
		ia.ID = uuid.New()
	}
	return nil
}

func (al *SecurityAuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	if al.LogID == "" {
		al.LogID = "LOG-" + uuid.New().String()[:8]
	}
	return nil
}

// Business logic methods for ComplianceRule
func (cr *ComplianceRule) CheckActive() bool {
	now := time.Now()
	if cr.ExpiryDate != nil {
		return cr.IsActive && now.After(cr.EffectiveDate) && now.Before(*cr.ExpiryDate)
	}
	return cr.IsActive && now.After(cr.EffectiveDate)
}

func (cr *ComplianceRule) RequiresReview() bool {
	if cr.NextReviewDate == nil {
		return false
	}
	return time.Now().After(*cr.NextReviewDate)
}

// Business logic methods for ComplianceCheck
func (cc *ComplianceCheck) IsPassed() bool {
	return cc.CheckStatus == "passed"
}

func (cc *ComplianceCheck) RequiresManualReview() bool {
	return cc.CheckStatus == "warning" || cc.IsManualReview
}

func (cc *ComplianceCheck) Pass() {
	cc.CheckStatus = "passed"
	cc.ComplianceScore = 100
}

func (cc *ComplianceCheck) Fail(reasons []string) {
	cc.CheckStatus = "failed"
	cc.ComplianceScore = 0
	// Convert reasons to JSON string
	cc.FailureReasons = "[]" // Placeholder - would serialize reasons array
}

// Business logic methods for SanctionsScreening
func (ss *SanctionsScreening) HasMatches() bool {
	return ss.MatchesFound > 0
}

func (ss *SanctionsScreening) IsHighRisk() bool {
	return ss.HighestRiskScore >= 80 || ss.Status == "confirmed_match"
}

func (ss *SanctionsScreening) Block(reason string) {
	ss.IsBlocked = true
	ss.BlockReason = reason
	ss.Status = "confirmed_match"
}

func (ss *SanctionsScreening) Unblock() {
	ss.IsBlocked = false
	now := time.Now()
	ss.UnblockDate = &now
	ss.Status = "false_positive"
}

// Business logic methods for DataProtectionRecord
func (dpr *DataProtectionRecord) IsRetentionExpired() bool {
	if dpr.RetentionPeriod <= 0 {
		return false
	}
	expiryDate := dpr.CreatedAt.AddDate(0, 0, dpr.RetentionPeriod)
	return time.Now().After(expiryDate)
}

func (dpr *DataProtectionRecord) SoftDelete() {
	dpr.IsDeleted = true
	dpr.DeletionMethod = "soft_delete"
	now := time.Now()
	dpr.DeletionDate = &now
}

func (dpr *DataProtectionRecord) WithdrawConsent() {
	dpr.ConsentWithdrawn = true
	now := time.Now()
	dpr.ConsentWithdrawDate = &now
}

// Business logic methods for SecurityIncident
func (si *SecurityIncident) IsCritical() bool {
	return si.Severity == "critical"
}

func (si *SecurityIncident) IsDataBreach() bool {
	return si.DataCompromised
}

func (si *SecurityIncident) RequiresNotification() bool {
	return si.NotificationRequired || si.IsCritical() || si.IsDataBreach()
}

func (si *SecurityIncident) Contain() {
	si.Status = "contained"
	now := time.Now()
	si.ContainmentDate = &now
}

func (si *SecurityIncident) Resolve() {
	si.Status = "resolved"
	now := time.Now()
	si.RecoveryDate = &now
}

// Business logic methods for SecurityAuditLog
func (al *SecurityAuditLog) CheckHighRisk() bool {
	return al.IsHighRisk || al.RiskScore >= 70
}

func (al *SecurityAuditLog) CalculateRiskScore() {
	score := 0.0

	// Failed requests increase risk
	if !al.IsSuccessful {
		score += 20
	}

	// Certain actions are inherently risky
	switch al.Action {
	case "delete", "modify_sensitive_data", "admin_access":
		score += 30
	case "bulk_export", "data_access":
		score += 20
	}

	// High response times might indicate issues
	if al.Duration > 5000 { // 5 seconds
		score += 10
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	al.RiskScore = score
	al.IsHighRisk = score >= 70
}
