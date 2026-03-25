package compliance

import (
	"time"
	
	"github.com/google/uuid"
	"github.com/lib/pq"
	
	"smartsure/pkg/database"
)

// ComplianceBusinessContinuityPlan represents a business continuity plan
type ComplianceBusinessContinuityPlan struct {
	database.BaseModel
	ComplianceID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	PlanNumber           string         `gorm:"uniqueIndex;not null" json:"plan_number"`
	PlanName             string         `gorm:"not null" json:"plan_name"`
	PlanType             string         `json:"plan_type"` // disaster_recovery, emergency_response, crisis_management
	Scope                pq.StringArray `gorm:"type:text[]" json:"scope"`
	Status               string         `gorm:"index;not null;default:'draft'" json:"status"`
	Version              string         `json:"version"`
	CriticalProcesses    pq.StringArray `gorm:"type:text[]" json:"critical_processes"`
	CriticalSystems      pq.StringArray `gorm:"type:text[]" json:"critical_systems"`
	CriticalPersonnel    pq.StringArray `gorm:"type:text[]" json:"critical_personnel"`
	RTO                  int            `json:"rto_hours"`  // Recovery Time Objective
	RPO                  int            `json:"rpo_hours"`  // Recovery Point Objective
	MAO                  int            `json:"mao_hours"`  // Maximum Acceptable Outage
	MTPD                 int            `json:"mtpd_hours"` // Maximum Tolerable Period of Disruption
	RecoveryStrategies   pq.StringArray `gorm:"type:text[]" json:"recovery_strategies"`
	BackupSites          pq.StringArray `gorm:"type:text[]" json:"backup_sites"`
	AlternateProcesses   pq.StringArray `gorm:"type:text[]" json:"alternate_processes"`
	CommunicationPlan    string         `gorm:"type:text" json:"communication_plan"`
	ContactList          pq.StringArray `gorm:"type:text[]" json:"contact_list"`
	VendorDependencies   pq.StringArray `gorm:"type:text[]" json:"vendor_dependencies"`
	ResourceRequirements map[string]int `gorm:"type:jsonb" json:"resource_requirements"`
	TestingSchedule      string         `json:"testing_schedule"`
	LastTestDate         time.Time      `json:"last_test_date"`
	NextTestDate         time.Time      `gorm:"index" json:"next_test_date"`
	TestScenarios        pq.StringArray `gorm:"type:text[]" json:"test_scenarios"`
	TestResults          string         `json:"test_results"`
	LessonsLearned       string         `gorm:"type:text" json:"lessons_learned"`
	ImprovementActions   pq.StringArray `gorm:"type:text[]" json:"improvement_actions"`
	ActivationCriteria   pq.StringArray `gorm:"type:text[]" json:"activation_criteria"`
	ActivationAuthority  uuid.UUID      `json:"activation_authority"`
	LastActivation       time.Time      `json:"last_activation_date"`
	ActivationCount      int            `json:"activation_count"`
	ReviewDate           time.Time      `json:"review_date"`
	ApprovalDate         time.Time      `json:"approval_date"`
	ApprovedBy           uuid.UUID      `json:"approved_by"`
	TrainingRequired     bool           `json:"training_required"`
	TrainingCompleted    bool           `json:"training_completed"`
}

// ComplianceControlAssessment represents a control effectiveness assessment
type ComplianceControlAssessment struct {
	database.BaseModel
	ComplianceID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	AssessmentNumber       string         `gorm:"uniqueIndex;not null" json:"assessment_number"`
	ControlID              string         `gorm:"index;not null" json:"control_id"`
	ControlName            string         `json:"control_name"`
	AssessmentType         string         `json:"assessment_type"`   // self, internal, external
	AssessmentMethod       string         `json:"assessment_method"` // test, review, observation, inquiry
	Status                 string         `gorm:"index;not null;default:'planned'" json:"status"`
	DesignEffectiveness    string         `json:"design_effectiveness"`
	OperatingEffectiveness string         `json:"operating_effectiveness"`
	OverallEffectiveness   string         `json:"overall_effectiveness"`
	TestProcedures         pq.StringArray `gorm:"type:text[]" json:"test_procedures"`
	SampleSize             int            `json:"sample_size"`
	ExceptionsFound        int            `json:"exceptions_found"`
	ExceptionRate          float64        `json:"exception_rate"`
	RootCauses             pq.StringArray `gorm:"type:text[]" json:"root_causes"`
	Recommendations        pq.StringArray `gorm:"type:text[]" json:"recommendations"`
	ManagementResponse     string         `gorm:"type:text" json:"management_response"`
	RemediationPlan        string         `gorm:"type:text" json:"remediation_plan"`
	RemediationDeadline    time.Time      `json:"remediation_deadline"`
	RemediationStatus      string         `json:"remediation_status"`
	AssessmentDate         time.Time      `json:"assessment_date"`
	AssessedBy             uuid.UUID      `json:"assessed_by"`
	ReviewedBy             uuid.UUID      `json:"reviewed_by"`
	Evidence               pq.StringArray `gorm:"type:text[]" json:"evidence"`
	WorkPapers             pq.StringArray `gorm:"type:text[]" json:"work_papers"`
	FollowUpRequired       bool           `json:"follow_up_required"`
	FollowUpDate           time.Time      `json:"follow_up_date"`
	PriorAssessment        string         `json:"prior_assessment_reference"`
	TrendAnalysis          string         `json:"trend_analysis"`
}

// ComplianceVendorAssessment represents third-party vendor compliance assessment
type ComplianceVendorAssessment struct {
	database.BaseModel
	ComplianceID             uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	VendorID                 uuid.UUID      `gorm:"type:uuid;not null;index" json:"vendor_id"`
	AssessmentNumber         string         `gorm:"uniqueIndex;not null" json:"assessment_number"`
	VendorName               string         `gorm:"not null" json:"vendor_name"`
	VendorType               string         `json:"vendor_type"`
	CriticalityLevel         string         `json:"criticality_level"`
	ServiceDescription       string         `gorm:"type:text" json:"service_description"`
	DataAccess               bool           `json:"data_access"`
	DataTypes                pq.StringArray `gorm:"type:text[]" json:"data_types_accessed"`
	ComplianceRequirements   pq.StringArray `gorm:"type:text[]" json:"compliance_requirements"`
	CertificationsRequired   pq.StringArray `gorm:"type:text[]" json:"certifications_required"`
	CertificationsHeld       pq.StringArray `gorm:"type:text[]" json:"certifications_held"`
	AssessmentScore          float64        `json:"assessment_score"`
	RiskRating               string         `json:"risk_rating"`
	IssuesIdentified         int            `json:"issues_identified"`
	CriticalIssues           int            `json:"critical_issues"`
	RemediationRequired      bool           `json:"remediation_required"`
	RemediationPlan          string         `gorm:"type:text" json:"remediation_plan"`
	ContractClauses          pq.StringArray `gorm:"type:text[]" json:"contract_clauses"`
	SLACompliance            float64        `json:"sla_compliance"`
	IncidentsReported        int            `json:"incidents_reported"`
	LastIncidentDate         time.Time      `json:"last_incident_date"`
	AuditRights              bool           `json:"audit_rights"`
	LastAuditDate            time.Time      `json:"last_audit_date"`
	NextAssessmentDate       time.Time      `gorm:"index" json:"next_assessment_date"`
	AssessmentFrequency      string         `json:"assessment_frequency"`
	DueDiligenceCompleted    bool           `json:"due_diligence_completed"`
	BackgroundCheckCompleted bool           `json:"background_check_completed"`
	InsuranceVerified        bool           `json:"insurance_verified"`
	InsuranceAmount          float64        `json:"insurance_amount"`
	SubcontractorsAllowed    bool           `json:"subcontractors_allowed"`
	SubcontractorsList       pq.StringArray `gorm:"type:text[]" json:"subcontractors_list"`
	GeographicRestrictions   pq.StringArray `gorm:"type:text[]" json:"geographic_restrictions"`
	DataLocations            pq.StringArray `gorm:"type:text[]" json:"data_locations"`
}

// ComplianceTrainingRecord represents compliance training completion records
type ComplianceTrainingRecord struct {
	database.BaseModel
	ComplianceID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	UserID              uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TrainingID          string         `gorm:"index;not null" json:"training_id"`
	TrainingName        string         `gorm:"not null" json:"training_name"`
	TrainingType        string         `json:"training_type"`
	TrainingCategory    string         `json:"training_category"`
	DeliveryMethod      string         `json:"delivery_method"`
	Duration            int            `json:"duration_minutes"`
	EnrollmentDate      time.Time      `json:"enrollment_date"`
	StartDate           time.Time      `json:"start_date"`
	CompletionDate      time.Time      `json:"completion_date"`
	ExpiryDate          time.Time      `gorm:"index" json:"expiry_date"`
	Status              string         `gorm:"index" json:"status"`
	Progress            float64        `json:"progress_percentage"`
	Score               float64        `json:"score"`
	PassingScore        float64        `json:"passing_score"`
	Passed              bool           `json:"passed"`
	Attempts            int            `json:"attempts"`
	CertificateIssued   bool           `json:"certificate_issued"`
	CertificateNumber   string         `json:"certificate_number"`
	CertificateExpiry   time.Time      `json:"certificate_expiry"`
	RefresherRequired   bool           `json:"refresher_required"`
	RefresherDate       time.Time      `json:"refresher_date"`
	Mandatory           bool           `json:"mandatory"`
	Department          string         `json:"department"`
	JobRole             string         `json:"job_role"`
	Supervisor          uuid.UUID      `json:"supervisor"`
	TrainingProvider    string         `json:"training_provider"`
	Cost                float64        `json:"cost"`
	Feedback            string         `gorm:"type:text" json:"feedback"`
	Rating              int            `json:"rating"`
	ComplianceTopics    pq.StringArray `gorm:"type:text[]" json:"compliance_topics"`
	LearningObjectives  pq.StringArray `gorm:"type:text[]" json:"learning_objectives"`
	AssessmentQuestions int            `json:"assessment_questions"`
	TimeSpent           int            `json:"time_spent_minutes"`
}

// ComplianceObligationRegister represents regulatory and contractual obligations
type ComplianceObligationRegister struct {
	database.BaseModel
	ComplianceID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	ObligationNumber        string         `gorm:"uniqueIndex;not null" json:"obligation_number"`
	ObligationTitle         string         `gorm:"not null" json:"obligation_title"`
	ObligationType          string         `gorm:"index" json:"obligation_type"` // regulatory, contractual, voluntary
	Source                  string         `json:"source"`
	SourceDocument          string         `json:"source_document"`
	RegulatoryBody          string         `json:"regulatory_body"`
	Jurisdiction            string         `json:"jurisdiction"`
	EffectiveDate           time.Time      `json:"effective_date"`
	EndDate                 time.Time      `json:"end_date"`
	Description             string         `gorm:"type:text" json:"description"`
	Requirements            pq.StringArray `gorm:"type:text[]" json:"requirements"`
	Applicability           string         `json:"applicability"`
	BusinessUnits           pq.StringArray `gorm:"type:text[]" json:"business_units"`
	Processes               pq.StringArray `gorm:"type:text[]" json:"processes"`
	Systems                 pq.StringArray `gorm:"type:text[]" json:"systems"`
	ComplianceStatus        string         `gorm:"index" json:"compliance_status"`
	ImplementationStatus    string         `json:"implementation_status"`
	Controls                pq.StringArray `gorm:"type:text[]" json:"controls"`
	Evidence                pq.StringArray `gorm:"type:text[]" json:"evidence"`
	ResponsiblePerson       uuid.UUID      `json:"responsible_person"`
	ComplianceOwner         uuid.UUID      `json:"compliance_owner"`
	FrequencyOfReview       string         `json:"frequency_of_review"`
	LastReviewDate          time.Time      `json:"last_review_date"`
	NextReviewDate          time.Time      `gorm:"index" json:"next_review_date"`
	PenaltyForNonCompliance string         `gorm:"type:text" json:"penalty_for_non_compliance"`
	MaximumPenalty          float64        `json:"maximum_penalty"`
	ReportingRequired       bool           `json:"reporting_required"`
	ReportingFrequency      string         `json:"reporting_frequency"`
	LastReportDate          time.Time      `json:"last_report_date"`
	Changes                 pq.StringArray `gorm:"type:text[]" json:"changes"`
	ChangeNotification      time.Time      `json:"change_notification_date"`
	TransitionPeriod        int            `json:"transition_period_days"`
	Exceptions              pq.StringArray `gorm:"type:text[]" json:"exceptions"`
	Waivers                 pq.StringArray `gorm:"type:text[]" json:"waivers"`
	RelatedObligations      pq.StringArray `gorm:"type:text[]" json:"related_obligations"`
	ConflictingObligations  pq.StringArray `gorm:"type:text[]" json:"conflicting_obligations"`
}

// ComplianceDataSubjectRequest represents GDPR/CCPA data subject requests
type ComplianceDataSubjectRequest struct {
	database.BaseModel
	ComplianceID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	RequestNumber      string         `gorm:"uniqueIndex;not null" json:"request_number"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	RequestType        string         `gorm:"index;not null" json:"request_type"` // access, rectification, erasure, portability, etc.
	Status             string         `gorm:"index;not null;default:'received'" json:"status"`
	SubmittedDate      time.Time      `gorm:"not null" json:"submitted_date"`
	VerificationDate   time.Time      `json:"verification_date"`
	ProcessingDate     time.Time      `json:"processing_date"`
	CompletionDate     time.Time      `json:"completion_date"`
	DueDate            time.Time      `gorm:"index" json:"due_date"`
	IdentityVerified   bool           `json:"identity_verified"`
	VerificationMethod string         `json:"verification_method"`
	RequestDetails     string         `gorm:"type:text" json:"request_details"`
	DataCategories     pq.StringArray `gorm:"type:text[]" json:"data_categories"`
	Systems            pq.StringArray `gorm:"type:text[]" json:"systems_searched"`
	DataFound          bool           `json:"data_found"`
	DataProvided       bool           `json:"data_provided"`
	FormatProvided     string         `json:"format_provided"`
	DeliveryMethod     string         `json:"delivery_method"`
	RejectionReason    string         `gorm:"type:text" json:"rejection_reason"`
	ExtensionRequired  bool           `json:"extension_required"`
	ExtensionReason    string         `json:"extension_reason"`
	ExtensionApproved  bool           `json:"extension_approved"`
	NewDueDate         time.Time      `json:"new_due_date"`
	ProcessedBy        uuid.UUID      `json:"processed_by"`
	ApprovedBy         uuid.UUID      `json:"approved_by"`
	CommunicationLog   pq.StringArray `gorm:"type:text[]" json:"communication_log"`
	Cost               float64        `json:"cost"`
	ChargedToUser      bool           `json:"charged_to_user"`
	ComplexRequest     bool           `json:"complex_request"`
	LegalConsultation  bool           `json:"legal_consultation"`
	ThirdPartyInvolved bool           `json:"third_party_involved"`
	ThirdParties       pq.StringArray `gorm:"type:text[]" json:"third_parties"`
	ResponseLetterID   string         `json:"response_letter_id"`
	AuditTrail         pq.StringArray `gorm:"type:text[]" json:"audit_trail"`
}

// ComplianceSecurityControl represents security controls implementation
type ComplianceSecurityControl struct {
	database.BaseModel
	ComplianceID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	ControlNumber        string         `gorm:"uniqueIndex;not null" json:"control_number"`
	ControlFamily        string         `gorm:"index" json:"control_family"`
	ControlTitle         string         `gorm:"not null" json:"control_title"`
	ControlDescription   string         `gorm:"type:text" json:"control_description"`
	ControlType          string         `json:"control_type"`
	Framework            string         `gorm:"index" json:"framework"` // NIST, ISO27001, CIS, etc.
	FrameworkReference   string         `json:"framework_reference"`
	Priority             string         `json:"priority"`
	ImplementationStatus string         `gorm:"index" json:"implementation_status"`
	ImplementationDate   time.Time      `json:"implementation_date"`
	Automation           string         `json:"automation"` // manual, semi-automated, automated
	Technology           pq.StringArray `gorm:"type:text[]" json:"technology"`
	ResponsibleTeam      string         `json:"responsible_team"`
	ControlOwner         uuid.UUID      `json:"control_owner"`
	TestingRequired      bool           `json:"testing_required"`
	TestingFrequency     string         `json:"testing_frequency"`
	LastTestDate         time.Time      `json:"last_test_date"`
	NextTestDate         time.Time      `gorm:"index" json:"next_test_date"`
	TestResults          string         `json:"test_results"`
	Effectiveness        string         `json:"effectiveness"`
	Maturity             int            `json:"maturity_level"` // 1-5
	Gaps                 pq.StringArray `gorm:"type:text[]" json:"gaps"`
	RemediationPlan      string         `gorm:"type:text" json:"remediation_plan"`
	RemediationDeadline  time.Time      `json:"remediation_deadline"`
	CompensatingControls pq.StringArray `gorm:"type:text[]" json:"compensating_controls"`
	Dependencies         pq.StringArray `gorm:"type:text[]" json:"dependencies"`
	RelatedControls      pq.StringArray `gorm:"type:text[]" json:"related_controls"`
	Threats              pq.StringArray `gorm:"type:text[]" json:"threats_addressed"`
	Risks                pq.StringArray `gorm:"type:text[]" json:"risks_mitigated"`
	CostOfImplementation float64        `json:"cost_of_implementation"`
	CostOfMaintenance    float64        `json:"cost_of_maintenance"`
	BusinessImpact       string         `json:"business_impact"`
	Documentation        pq.StringArray `gorm:"type:text[]" json:"documentation"`
	Evidence             pq.StringArray `gorm:"type:text[]" json:"evidence"`
	ReviewFrequency      string         `json:"review_frequency"`
	LastReviewDate       time.Time      `json:"last_review_date"`
	NextReviewDate       time.Time      `gorm:"index" json:"next_review_date"`
}
