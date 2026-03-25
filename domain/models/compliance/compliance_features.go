package compliance

import (
	"time"
	
	"github.com/google/uuid"
	"github.com/lib/pq"
	
	"smartsure/pkg/database"
)

// ComplianceIncident represents a security, privacy, or compliance incident
type ComplianceIncident struct {
	database.BaseModel
	ComplianceID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	IncidentNumber       string         `gorm:"uniqueIndex;not null" json:"incident_number"`
	IncidentType         string         `gorm:"index;not null" json:"incident_type"`
	Severity             string         `gorm:"index;not null" json:"severity"`
	Status               string         `gorm:"index;not null;default:'open'" json:"status"`
	Title                string         `gorm:"not null" json:"title"`
	Description          string         `gorm:"type:text" json:"description"`
	DetectedDate         time.Time      `gorm:"not null" json:"detected_date"`
	ReportedDate         time.Time      `json:"reported_date"`
	ContainedDate        time.Time      `json:"contained_date"`
	EradicatedDate       time.Time      `json:"eradicated_date"`
	RecoveredDate        time.Time      `json:"recovered_date"`
	ClosedDate           time.Time      `json:"closed_date"`
	ImpactAssessment     string         `gorm:"type:text" json:"impact_assessment"`
	RootCause            string         `gorm:"type:text" json:"root_cause"`
	AffectedSystems      pq.StringArray `gorm:"type:text[]" json:"affected_systems"`
	AffectedUsers        int            `json:"affected_users"`
	DataCompromised      bool           `json:"data_compromised"`
	DataTypes            pq.StringArray `gorm:"type:text[]" json:"data_types_compromised"`
	ResponseTeam         pq.StringArray `gorm:"type:text[]" json:"response_team"`
	ResponseActions      pq.StringArray `gorm:"type:text[]" json:"response_actions"`
	LessonsLearned       string         `gorm:"type:text" json:"lessons_learned"`
	PreventiveMeasures   pq.StringArray `gorm:"type:text[]" json:"preventive_measures"`
	NotificationRequired bool           `json:"notification_required"`
	NotifiedParties      pq.StringArray `gorm:"type:text[]" json:"notified_parties"`
	NotificationDate     time.Time      `json:"notification_date"`
	RegulatoryReported   bool           `json:"regulatory_reported"`
	ReportedTo           pq.StringArray `gorm:"type:text[]" json:"reported_to"`
	EstimatedCost        float64        `json:"estimated_cost"`
	ActualCost           float64        `json:"actual_cost"`
	InsuranceClaimed     bool           `json:"insurance_claimed"`
	InsuranceAmount      float64        `json:"insurance_amount"`
	EvidenceCollected    bool           `json:"evidence_collected"`
	ForensicsCompleted   bool           `json:"forensics_completed"`
	LegalInvolvement     bool           `json:"legal_involvement"`
	MediaExposure        bool           `json:"media_exposure"`
	ReputationalImpact   string         `json:"reputational_impact"`
	RecurrenceLikelihood string         `json:"recurrence_likelihood"`
	IncidentManager      uuid.UUID      `json:"incident_manager"`
	PostIncidentReview   time.Time      `json:"post_incident_review_date"`
	ReviewCompleted      bool           `json:"review_completed"`
	RelatedIncidents     pq.StringArray `gorm:"type:text[]" json:"related_incidents"`
	TTD                  int            `json:"time_to_detect_hours"`  // Time to Detect
	TTC                  int            `json:"time_to_contain_hours"` // Time to Contain
	TTR                  int            `json:"time_to_recover_hours"` // Time to Recover
	MTTD                 float64        `json:"mean_time_to_detect"`   // Mean Time to Detect
	MTTR                 float64        `json:"mean_time_to_recover"`  // Mean Time to Recover
}

// ComplianceRisk represents an identified compliance risk
type ComplianceRisk struct {
	database.BaseModel
	ComplianceID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	RiskNumber           string         `gorm:"uniqueIndex;not null" json:"risk_number"`
	RiskName             string         `gorm:"not null" json:"risk_name"`
	RiskDescription      string         `gorm:"type:text" json:"risk_description"`
	RiskCategory         string         `gorm:"index;not null" json:"risk_category"`
	RiskType             string         `gorm:"index" json:"risk_type"`
	Status               string         `gorm:"index;not null;default:'identified'" json:"status"`
	Likelihood           int            `json:"likelihood"` // 1-5
	Impact               int            `json:"impact"`     // 1-5
	InherentRiskScore    float64        `json:"inherent_risk_score"`
	CurrentControls      pq.StringArray `gorm:"type:text[]" json:"current_controls"`
	ControlEffectiveness string         `json:"control_effectiveness"`
	ResidualRiskScore    float64        `json:"residual_risk_score"`
	RiskAppetite         float64        `json:"risk_appetite"`
	RiskTolerance        float64        `json:"risk_tolerance"`
	Treatment            string         `gorm:"index" json:"treatment"`
	TreatmentPlan        string         `gorm:"type:text" json:"treatment_plan"`
	MitigationActions    pq.StringArray `gorm:"type:text[]" json:"mitigation_actions"`
	TargetRiskScore      float64        `json:"target_risk_score"`
	RiskOwner            uuid.UUID      `gorm:"index" json:"risk_owner"`
	ActionOwner          uuid.UUID      `json:"action_owner"`
	IdentifiedDate       time.Time      `json:"identified_date"`
	AssessmentDate       time.Time      `json:"assessment_date"`
	ReviewDate           time.Time      `json:"review_date"`
	NextReviewDate       time.Time      `gorm:"index" json:"next_review_date"`
	TargetCompletionDate time.Time      `json:"target_completion_date"`
	ActualCompletionDate time.Time      `json:"actual_completion_date"`
	TriggerEvents        pq.StringArray `gorm:"type:text[]" json:"trigger_events"`
	RiskIndicators       pq.StringArray `gorm:"type:text[]" json:"risk_indicators"`
	EarlyWarningSignals  pq.StringArray `gorm:"type:text[]" json:"early_warning_signals"`
	RelatedRisks         pq.StringArray `gorm:"type:text[]" json:"related_risks"`
	RegulatoryImpact     bool           `json:"regulatory_impact"`
	FinancialImpact      float64        `json:"financial_impact"`
	OperationalImpact    string         `json:"operational_impact"`
	ReputationalImpact   string         `json:"reputational_impact"`
	MonitoringRequired   bool           `json:"monitoring_required"`
	MonitoringFrequency  string         `json:"monitoring_frequency"`
	EscalationRequired   bool           `json:"escalation_required"`
	EscalatedTo          uuid.UUID      `json:"escalated_to"`
	AcceptanceApproval   bool           `json:"acceptance_approval"`
	ApprovedBy           uuid.UUID      `json:"approved_by"`
	ApprovalDate         time.Time      `json:"approval_date"`
}

// ComplianceConsentRecord represents a data subject consent record
type ComplianceConsentRecord struct {
	database.BaseModel
	ComplianceID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"compliance_id"`
	UserID               uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	ConsentID            string          `gorm:"uniqueIndex;not null" json:"consent_id"`
	ConsentType          string          `gorm:"index;not null" json:"consent_type"`
	ConsentStatus        string          `gorm:"index;not null" json:"consent_status"`
	ConsentVersion       string          `gorm:"not null" json:"consent_version"`
	ConsentDate          time.Time       `gorm:"not null" json:"consent_date"`
	WithdrawalDate       time.Time       `json:"withdrawal_date"`
	ExpiryDate           time.Time       `gorm:"index" json:"expiry_date"`
	Purposes             pq.StringArray  `gorm:"type:text[]" json:"purposes"`
	DataCategories       pq.StringArray  `gorm:"type:text[]" json:"data_categories"`
	ProcessingActivities pq.StringArray  `gorm:"type:text[]" json:"processing_activities"`
	ThirdParties         pq.StringArray  `gorm:"type:text[]" json:"third_parties"`
	ConsentMethod        string          `json:"consent_method"`    // explicit, implicit, opt-in, opt-out
	CollectionMethod     string          `json:"collection_method"` // web, mobile, paper, verbal
	CollectionPoint      string          `json:"collection_point"`
	IPAddress            string          `json:"ip_address"`
	UserAgent            string          `json:"user_agent"`
	Language             string          `json:"language"`
	ConsentText          string          `gorm:"type:text" json:"consent_text"`
	PrivacyNoticeVersion string          `json:"privacy_notice_version"`
	PrivacyNoticeLink    string          `json:"privacy_notice_link"`
	ParentalConsent      bool            `json:"parental_consent"`
	GuardianID           uuid.UUID       `json:"guardian_id"`
	IsMinor              bool            `json:"is_minor"`
	AgeVerified          bool            `json:"age_verified"`
	DoubleOptIn          bool            `json:"double_opt_in"`
	ConfirmationSent     bool            `json:"confirmation_sent"`
	ConfirmationDate     time.Time       `json:"confirmation_date"`
	Preferences          map[string]bool `gorm:"type:jsonb" json:"preferences"`
	GranularConsents     map[string]bool `gorm:"type:jsonb" json:"granular_consents"`
	ConsentProof         string          `json:"consent_proof"` // hash or reference
	LegalBasis           string          `json:"legal_basis"`
	LegitimateInterests  pq.StringArray  `gorm:"type:text[]" json:"legitimate_interests"`
	RetentionPeriod      int             `json:"retention_period_days"`
	AutoRenewal          bool            `json:"auto_renewal"`
	RenewalNotification  int             `json:"renewal_notification_days"`
	LastUpdated          time.Time       `json:"last_updated"`
}

// CompliancePrivacyImpactAssessment represents a Privacy Impact Assessment
type CompliancePrivacyImpactAssessment struct {
	database.BaseModel
	ComplianceID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	PIANumber             string         `gorm:"uniqueIndex;not null" json:"pia_number"`
	ProjectName           string         `gorm:"not null" json:"project_name"`
	ProjectDescription    string         `gorm:"type:text" json:"project_description"`
	AssessmentType        string         `json:"assessment_type"` // initial, review, update
	Status                string         `gorm:"index;not null;default:'draft'" json:"status"`
	DataFlowDescription   string         `gorm:"type:text" json:"data_flow_description"`
	DataCategories        pq.StringArray `gorm:"type:text[]" json:"data_categories"`
	DataSources           pq.StringArray `gorm:"type:text[]" json:"data_sources"`
	ProcessingPurposes    pq.StringArray `gorm:"type:text[]" json:"processing_purposes"`
	DataRecipients        pq.StringArray `gorm:"type:text[]" json:"data_recipients"`
	RetentionPeriods      map[string]int `gorm:"type:jsonb" json:"retention_periods"`
	SecurityMeasures      pq.StringArray `gorm:"type:text[]" json:"security_measures"`
	RisksIdentified       pq.StringArray `gorm:"type:text[]" json:"risks_identified"`
	MitigationMeasures    pq.StringArray `gorm:"type:text[]" json:"mitigation_measures"`
	ResidualRisks         pq.StringArray `gorm:"type:text[]" json:"residual_risks"`
	PrivacyByDesign       bool           `json:"privacy_by_design"`
	DataMinimization      bool           `json:"data_minimization"`
	ConsentRequired       bool           `json:"consent_required"`
	LegalBasis            string         `json:"legal_basis"`
	CrossBorderTransfer   bool           `json:"cross_border_transfer"`
	TransferMechanisms    pq.StringArray `gorm:"type:text[]" json:"transfer_mechanisms"`
	ThirdPartyInvolvement bool           `json:"third_party_involvement"`
	ThirdParties          pq.StringArray `gorm:"type:text[]" json:"third_parties"`
	DataSubjectRights     pq.StringArray `gorm:"type:text[]" json:"data_subject_rights"`
	AssessmentDate        time.Time      `json:"assessment_date"`
	ReviewDate            time.Time      `json:"review_date"`
	ApprovalDate          time.Time      `json:"approval_date"`
	AssessedBy            uuid.UUID      `json:"assessed_by"`
	ReviewedBy            uuid.UUID      `json:"reviewed_by"`
	ApprovedBy            uuid.UUID      `json:"approved_by"`
	DPOReview             bool           `json:"dpo_review"`
	DPOReviewDate         time.Time      `json:"dpo_review_date"`
	DPORecommendations    string         `gorm:"type:text" json:"dpo_recommendations"`
	HighRisk              bool           `json:"high_risk"`
	ConsultationRequired  bool           `json:"consultation_required"`
	ConsultationDate      time.Time      `json:"consultation_date"`
	ConsultationOutcome   string         `gorm:"type:text" json:"consultation_outcome"`
}

// ComplianceQualityMetric represents quality management metrics
type ComplianceQualityMetric struct {
	database.BaseModel
	ComplianceID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"compliance_id"`
	MetricName         string         `gorm:"not null;index" json:"metric_name"`
	MetricType         string         `gorm:"index" json:"metric_type"`
	MetricCategory     string         `json:"metric_category"`
	MeasurementPeriod  string         `json:"measurement_period"`
	Target             float64        `json:"target"`
	Actual             float64        `json:"actual"`
	Variance           float64        `json:"variance"`
	VariancePercentage float64        `json:"variance_percentage"`
	Trend              string         `json:"trend"` // improving, stable, declining
	Status             string         `gorm:"index" json:"status"`
	DataSource         string         `json:"data_source"`
	CollectionMethod   string         `json:"collection_method"`
	MeasurementDate    time.Time      `gorm:"index" json:"measurement_date"`
	PreviousValue      float64        `json:"previous_value"`
	BaselineValue      float64        `json:"baseline_value"`
	BestValue          float64        `json:"best_value"`
	WorstValue         float64        `json:"worst_value"`
	IndustryBenchmark  float64        `json:"industry_benchmark"`
	InternalBenchmark  float64        `json:"internal_benchmark"`
	ImprovementActions pq.StringArray `gorm:"type:text[]" json:"improvement_actions"`
	ResponsiblePerson  uuid.UUID      `json:"responsible_person"`
	ReviewFrequency    string         `json:"review_frequency"`
	LastReviewDate     time.Time      `json:"last_review_date"`
	NextReviewDate     time.Time      `gorm:"index" json:"next_review_date"`
	CriticalThreshold  float64        `json:"critical_threshold"`
	WarningThreshold   float64        `json:"warning_threshold"`
	AlertsEnabled      bool           `json:"alerts_enabled"`
	AlertsSent         int            `json:"alerts_sent"`
	CorrectiveActions  int            `json:"corrective_actions"`
	PreventiveActions  int            `json:"preventive_actions"`
}
