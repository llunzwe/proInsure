package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserFraudDetection contains comprehensive fraud detection and prevention data
type UserFraudDetection struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Fraud Scoring
	OverallFraudScore     float64    `gorm:"default:0" json:"overall_fraud_score"`
	IdentityFraudScore    float64    `gorm:"default:0" json:"identity_fraud_score"`
	TransactionFraudScore float64    `gorm:"default:0" json:"transaction_fraud_score"`
	ClaimFraudScore       float64    `gorm:"default:0" json:"claim_fraud_score"`
	AccountFraudScore     float64    `gorm:"default:0" json:"account_fraud_score"`
	DocumentFraudScore    float64    `gorm:"default:0" json:"document_fraud_score"`
	BehavioralFraudScore  float64    `gorm:"default:0" json:"behavioral_fraud_score"`
	NetworkFraudScore     float64    `gorm:"default:0" json:"network_fraud_score"`
	FraudRiskCategory     string     `gorm:"type:varchar(20)" json:"fraud_risk_category"` // low/medium/high/critical
	LastFraudAssessment   *time.Time `json:"last_fraud_assessment"`
	NextFraudAssessment   *time.Time `json:"next_fraud_assessment"`

	// Fraud Patterns
	KnownFraudPatterns []string                 `gorm:"type:json" json:"known_fraud_patterns"`
	SuspiciousPatterns []string                 `gorm:"type:json" json:"suspicious_patterns"`
	PatternMatchScore  float64                  `gorm:"default:0" json:"pattern_match_score"`
	AnomalyPatterns    []map[string]interface{} `gorm:"type:json" json:"anomaly_patterns"`
	FraudIndicators    map[string]bool          `gorm:"type:json" json:"fraud_indicators"`
	RuleViolations     []string                 `gorm:"type:json" json:"rule_violations"`
	MLFraudPrediction  float64                  `gorm:"default:0" json:"ml_fraud_prediction"`
	FraudModelVersion  string                   `gorm:"type:varchar(20)" json:"fraud_model_version"`

	// Identity Verification
	IdentityVerified     bool                     `gorm:"default:false" json:"identity_verified"`
	VerificationMethod   string                   `gorm:"type:varchar(50)" json:"verification_method"`
	VerificationDate     *time.Time               `json:"verification_date"`
	VerificationAttempts int                      `gorm:"default:0" json:"verification_attempts"`
	FailedVerifications  int                      `gorm:"default:0" json:"failed_verifications"`
	DocumentsVerified    []string                 `gorm:"type:json" json:"documents_verified"`
	BiometricVerified    bool                     `gorm:"default:false" json:"biometric_verified"`
	LivenessCheckPassed  bool                     `gorm:"default:false" json:"liveness_check_passed"`
	FaceMatchScore       float64                  `gorm:"default:0" json:"face_match_score"`
	DocumentAuthenticity float64                  `gorm:"default:0" json:"document_authenticity"`
	AddressVerified      bool                     `gorm:"default:false" json:"address_verified"`
	PhoneVerified        bool                     `gorm:"default:false" json:"phone_verified"`
	EmailVerified        bool                     `gorm:"default:false" json:"email_verified"`
	SocialMediaVerified  bool                     `gorm:"default:false" json:"social_media_verified"`
	ReferenceChecks      []map[string]interface{} `gorm:"type:json" json:"reference_checks"`

	// Blacklist & Watchlist
	BlacklistStatus     bool       `gorm:"default:false;index" json:"blacklist_status"`
	BlacklistReason     string     `gorm:"type:text" json:"blacklist_reason"`
	BlacklistDate       *time.Time `json:"blacklist_date"`
	BlacklistExpiry     *time.Time `json:"blacklist_expiry"`
	WatchlistStatus     bool       `gorm:"default:false;index" json:"watchlist_status"`
	WatchlistCategories []string   `gorm:"type:json" json:"watchlist_categories"`
	SanctionsListCheck  bool       `gorm:"default:false" json:"sanctions_list_check"`
	PEPCheck            bool       `gorm:"default:false" json:"pep_check"`
	AdverseMediaCheck   bool       `gorm:"default:false" json:"adverse_media_check"`
	GlobalDatabaseCheck bool       `gorm:"default:false" json:"global_database_check"`

	// Investigation History
	UnderInvestigation     bool                   `gorm:"default:false;index" json:"under_investigation"`
	InvestigationID        *uuid.UUID             `gorm:"type:uuid" json:"investigation_id"`
	InvestigationStartDate *time.Time             `json:"investigation_start_date"`
	InvestigationEndDate   *time.Time             `json:"investigation_end_date"`
	InvestigationStatus    string                 `gorm:"type:varchar(50)" json:"investigation_status"`
	InvestigationFindings  map[string]interface{} `gorm:"type:json" json:"investigation_findings"`
	InvestigatorNotes      string                 `gorm:"type:text" json:"investigator_notes"`
	EvidenceCollected      []string               `gorm:"type:json" json:"evidence_collected"`
	InvestigationOutcome   string                 `gorm:"type:varchar(50)" json:"investigation_outcome"`
	ActionsTaken           []string               `gorm:"type:json" json:"actions_taken"`

	// Network Analysis
	LinkedAccounts        []uuid.UUID              `gorm:"type:json" json:"linked_accounts"`
	SharedDevices         []string                 `gorm:"type:json" json:"shared_devices"`
	SharedPaymentMethods  []string                 `gorm:"type:json" json:"shared_payment_methods"`
	NetworkRiskScore      float64                  `gorm:"default:0" json:"network_risk_score"`
	FraudRingMembership   bool                     `gorm:"default:false" json:"fraud_ring_membership"`
	FraudRingID           *uuid.UUID               `gorm:"type:uuid" json:"fraud_ring_id"`
	NetworkCentrality     float64                  `gorm:"default:0" json:"network_centrality"`
	SuspiciousConnections int                      `gorm:"default:0" json:"suspicious_connections"`
	CollisionPatterns     []map[string]interface{} `gorm:"type:json" json:"collision_patterns"`
}

// UserAnomalyDetection tracks behavioral anomalies and deviations
type UserAnomalyDetection struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Anomaly Detection
	AnomalyType         string                 `gorm:"type:varchar(50)" json:"anomaly_type"`
	AnomalyScore        float64                `gorm:"default:0" json:"anomaly_score"`
	AnomalySeverity     string                 `gorm:"type:varchar(20)" json:"anomaly_severity"`
	DetectionTimestamp  time.Time              `json:"detection_timestamp"`
	AnomalyDescription  string                 `gorm:"type:text" json:"anomaly_description"`
	BaselineBehavior    map[string]interface{} `gorm:"type:json" json:"baseline_behavior"`
	DeviationMetrics    map[string]float64     `gorm:"type:json" json:"deviation_metrics"`
	ContextualFactors   map[string]interface{} `gorm:"type:json" json:"contextual_factors"`
	FalsePositive       bool                   `gorm:"default:false" json:"false_positive"`
	ConfirmedAnomaly    bool                   `gorm:"default:false" json:"confirmed_anomaly"`
	ResolutionStatus    string                 `gorm:"type:varchar(50)" json:"resolution_status"`
	ResolutionTimestamp *time.Time             `json:"resolution_timestamp"`
	PreventiveMeasures  []string               `gorm:"type:json" json:"preventive_measures"`
	RecurrenceCount     int                    `gorm:"default:0" json:"recurrence_count"`
}

// UserIdentityVerification contains identity verification history and results
type UserIdentityVerification struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Verification Details
	VerificationID        string                 `gorm:"type:varchar(100);uniqueIndex" json:"verification_id"`
	VerificationType      string                 `gorm:"type:varchar(50)" json:"verification_type"`
	VerificationProvider  string                 `gorm:"type:varchar(50)" json:"verification_provider"`
	VerificationStatus    string                 `gorm:"type:varchar(20)" json:"verification_status"`
	VerificationScore     float64                `gorm:"default:0" json:"verification_score"`
	VerificationTimestamp time.Time              `json:"verification_timestamp"`
	DocumentType          string                 `gorm:"type:varchar(50)" json:"document_type"`
	DocumentNumber        string                 `gorm:"type:varchar(100)" json:"document_number"`
	DocumentCountry       string                 `gorm:"type:varchar(2)" json:"document_country"`
	DocumentExpiry        *time.Time             `json:"document_expiry"`
	BiometricData         map[string]interface{} `gorm:"type:json" json:"biometric_data"`
	VerificationResults   map[string]interface{} `gorm:"type:json" json:"verification_results"`
	ManualReviewRequired  bool                   `gorm:"default:false" json:"manual_review_required"`
	ManualReviewCompleted bool                   `gorm:"default:false" json:"manual_review_completed"`
	ReviewerNotes         string                 `gorm:"type:text" json:"reviewer_notes"`
}

// UserFraudInvestigation contains detailed fraud investigation records
type UserFraudInvestigation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	InvestigatorID *uuid.UUID `gorm:"type:uuid" json:"investigator_id"`

	// Investigation Details
	CaseNumber             string                   `gorm:"type:varchar(50);uniqueIndex" json:"case_number"`
	Priority               string                   `gorm:"type:varchar(20)" json:"priority"`
	Status                 string                   `gorm:"type:varchar(50)" json:"status"`
	Type                   string                   `gorm:"type:varchar(50)" json:"type"`
	Reason                 string                   `gorm:"type:text" json:"reason"`
	StartDate              time.Time                `json:"start_date"`
	EndDate                *time.Time               `json:"end_date"`
	EstimatedLoss          decimal.Decimal          `gorm:"type:decimal(15,2)" json:"estimated_loss"`
	RecoveredAmount        decimal.Decimal          `gorm:"type:decimal(15,2)" json:"recovered_amount"`
	Evidence               []map[string]interface{} `gorm:"type:json" json:"evidence"`
	Interviews             []map[string]interface{} `gorm:"type:json" json:"interviews"`
	DataAnalysis           map[string]interface{}   `gorm:"type:json" json:"data_analysis"`
	Findings               string                   `gorm:"type:text" json:"findings"`
	Recommendations        []string                 `gorm:"type:json" json:"recommendations"`
	LegalActionTaken       bool                     `gorm:"default:false" json:"legal_action_taken"`
	LawEnforcementNotified bool                     `gorm:"default:false" json:"law_enforcement_notified"`
	CaseClosureReason      string                   `gorm:"type:text" json:"case_closure_reason"`
}

// UserBlacklistManagement manages user blacklist status and history
type UserBlacklistManagement struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID        uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	BlacklistedBy *uuid.UUID `gorm:"type:uuid" json:"blacklisted_by"`

	// Blacklist Details
	Status                 string            `gorm:"type:varchar(20)" json:"status"` // active/expired/removed
	Severity               string            `gorm:"type:varchar(20)" json:"severity"`
	Categories             []string          `gorm:"type:json" json:"categories"`
	Reason                 string            `gorm:"type:text" json:"reason"`
	Evidence               []string          `gorm:"type:json" json:"evidence"`
	StartDate              time.Time         `json:"start_date"`
	EndDate                *time.Time        `json:"end_date"`
	Duration               int               `json:"duration_days"`
	ReviewDate             *time.Time        `json:"review_date"`
	AppealEligible         bool              `gorm:"default:false" json:"appeal_eligible"`
	AppealSubmitted        bool              `gorm:"default:false" json:"appeal_submitted"`
	AppealDate             *time.Time        `json:"appeal_date"`
	AppealOutcome          string            `gorm:"type:varchar(50)" json:"appeal_outcome"`
	RemovalReason          string            `gorm:"type:text" json:"removal_reason"`
	ImpactedServices       []string          `gorm:"type:json" json:"impacted_services"`
	RelatedIncidents       []uuid.UUID       `gorm:"type:json" json:"related_incidents"`
	CrossPlatformBlacklist bool              `gorm:"default:false" json:"cross_platform_blacklist"`
	ExternalReferences     map[string]string `gorm:"type:json" json:"external_references"`
}

// UserFraudPrevention contains proactive fraud prevention measures
type UserFraudPrevention struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Prevention Measures
	PreventionLevel      string                     `gorm:"type:varchar(20)" json:"prevention_level"`
	ActiveMeasures       []string                   `gorm:"type:json" json:"active_measures"`
	VelocityLimits       map[string]int             `gorm:"type:json" json:"velocity_limits"`
	TransactionLimits    map[string]decimal.Decimal `gorm:"type:json" json:"transaction_limits"`
	GeofencingEnabled    bool                       `gorm:"default:false" json:"geofencing_enabled"`
	AllowedGeolocations  []map[string]float64       `gorm:"type:json" json:"allowed_geolocations"`
	DeviceRestrictions   []string                   `gorm:"type:json" json:"device_restrictions"`
	TimeRestrictions     map[string]string          `gorm:"type:json" json:"time_restrictions"`
	RequiresMFA          bool                       `gorm:"default:false" json:"requires_mfa"`
	RequiresBiometric    bool                       `gorm:"default:false" json:"requires_biometric"`
	RequiresManualReview bool                       `gorm:"default:false" json:"requires_manual_review"`
	StepUpAuthentication bool                       `gorm:"default:false" json:"step_up_authentication"`
	CoolingPeriod        int                        `json:"cooling_period_hours"`
	MaxFailedAttempts    int                        `gorm:"default:3" json:"max_failed_attempts"`
	SessionTimeout       int                        `json:"session_timeout_minutes"`
	IPWhitelisting       []string                   `gorm:"type:json" json:"ip_whitelisting"`
	BehavioralBaseline   map[string]interface{}     `gorm:"type:json" json:"behavioral_baseline"`
	RiskThresholds       map[string]float64         `gorm:"type:json" json:"risk_thresholds"`
	AutoBlockEnabled     bool                       `gorm:"default:true" json:"auto_block_enabled"`
	LastSecurityReview   *time.Time                 `json:"last_security_review"`
}

// TableName returns the table name
func (UserFraudDetection) TableName() string {
	return "user_fraud_detection"
}

// TableName returns the table name
func (UserAnomalyDetection) TableName() string {
	return "user_anomaly_detection"
}

// TableName returns the table name
func (UserIdentityVerification) TableName() string {
	return "user_identity_verification"
}

// TableName returns the table name
func (UserFraudInvestigation) TableName() string {
	return "user_fraud_investigation"
}

// TableName returns the table name
func (UserBlacklistManagement) TableName() string {
	return "user_blacklist_management"
}

// TableName returns the table name
func (UserFraudPrevention) TableName() string {
	return "user_fraud_prevention"
}
