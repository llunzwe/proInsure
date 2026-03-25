package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceFraudPatterns tracks historical fraud patterns and signatures
type DeviceFraudPatterns struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	PatternType   string    `gorm:"type:varchar(50)" json:"pattern_type"` // claim_pattern, behavior_pattern, network_pattern
	PatternName   string    `gorm:"type:varchar(100)" json:"pattern_name"`
	DetectionDate time.Time `json:"detection_date"`

	// Pattern Identification
	FraudSignature string  `gorm:"type:json" json:"fraud_signature"` // JSON object of pattern characteristics
	MatchScore     float64 `json:"match_score"`                      // 0-100 similarity to known fraud
	KnownSchemeID  string  `json:"known_scheme_id"`                  // Reference to fraud scheme database
	SchemeCategory string  `json:"scheme_category"`                  // insurance_fraud, identity_theft, device_theft

	// Fraud Ring Detection
	FraudRingID      *uuid.UUID `gorm:"type:uuid" json:"fraud_ring_id"`
	RingRole         string     `json:"ring_role"`                        // orchestrator, participant, victim
	RelatedDevices   string     `gorm:"type:json" json:"related_devices"` // JSON array of device IDs
	CrossDeviceScore float64    `json:"cross_device_score"`               // Correlation strength

	// Temporal Analysis
	TimePattern          string     `json:"time_pattern"`      // day_of_week, time_of_day, seasonal
	FrequencyPattern     string     `json:"frequency_pattern"` // regular, irregular, burst
	TemporalAnomalyScore float64    `json:"temporal_anomaly_score"`
	LastOccurrence       *time.Time `json:"last_occurrence"`
	OccurrenceCount      int        `json:"occurrence_count"`

	// Geographic Analysis
	GeographicHotspot   bool    `json:"geographic_hotspot"`
	LocationPattern     string  `gorm:"type:json" json:"location_pattern"` // JSON array of locations
	GeoAnomalyScore     float64 `json:"geo_anomaly_score"`
	CrossBorderActivity bool    `json:"cross_border_activity"`

	// Evolution Tracking
	MethodEvolution string  `gorm:"type:json" json:"method_evolution"` // JSON array of method changes
	AdaptationRate  float64 `json:"adaptation_rate"`                   // How quickly pattern changes
	Countermeasures string  `gorm:"type:json" json:"countermeasures"`  // JSON array

	// Financial Impact
	TotalLossAmount float64 `json:"total_loss_amount"`
	RecoveredAmount float64 `json:"recovered_amount"`
	RecoveryRate    float64 `json:"recovery_rate"`
	PotentialLoss   float64 `json:"potential_loss"`

	// Status
	Status          string     `gorm:"type:varchar(50)" json:"status"` // active, inactive, monitoring
	InvestigationID *uuid.UUID `gorm:"type:uuid" json:"investigation_id"`
	RiskLevel       string     `json:"risk_level"` // critical, high, medium, low

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceAnomalyDetection handles real-time anomaly detection
type DeviceAnomalyDetection struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	DetectionTime time.Time `json:"detection_time"`

	// Baseline Behavior
	BaselineEstablished bool      `json:"baseline_established"`
	BaselinePeriodDays  int       `json:"baseline_period_days"`
	BaselineData        string    `gorm:"type:json" json:"baseline_data"` // JSON object
	LastBaselineUpdate  time.Time `json:"last_baseline_update"`

	// Anomaly Detection
	AnomalyType     string  `gorm:"type:varchar(50)" json:"anomaly_type"` // usage, location, claim, transaction
	AnomalyScore    float64 `json:"anomaly_score"`                        // 0-100
	DeviationAmount float64 `json:"deviation_amount"`                     // Standard deviations from baseline
	ConfidenceLevel float64 `json:"confidence_level"`

	// Thresholds
	AlertThreshold       float64 `json:"alert_threshold"`
	BlockThreshold       float64 `json:"block_threshold"`
	InvestigateThreshold float64 `json:"investigate_threshold"`
	ThresholdBreached    string  `json:"threshold_breached"` // none, alert, investigate, block

	// Classification
	AnomalyCategory   string  `json:"anomaly_category"` // benign, suspicious, malicious
	AnomalySeverity   string  `json:"anomaly_severity"` // critical, high, medium, low
	FalsePositiveProb float64 `json:"false_positive_prob"`

	// Alert Management
	AlertGenerated        bool       `json:"alert_generated"`
	AlertID               *uuid.UUID `gorm:"type:uuid" json:"alert_id"`
	AlertPriority         int        `json:"alert_priority"` // 1-10
	InvestigationPriority int        `json:"investigation_priority"`

	// Machine Learning
	MLModelVersion    string  `json:"ml_model_version"`
	ModelConfidence   float64 `json:"model_confidence"`
	FeatureImportance string  `gorm:"type:json" json:"feature_importance"` // JSON object
	AdaptiveThreshold bool    `json:"adaptive_threshold"`
	LearningRate      float64 `json:"learning_rate"`

	// Response
	AutoResponse      bool       `json:"auto_response"`
	ResponseAction    string     `json:"response_action"` // monitor, flag, block, investigate
	ResponseTime      *time.Time `json:"response_time"`
	MitigationApplied string     `gorm:"type:json" json:"mitigation_applied"` // JSON array

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceIdentityVerification manages identity verification history
type DeviceIdentityVerification struct {
	database.BaseModel
	DeviceID         uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID           uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	VerificationDate time.Time `json:"verification_date"`

	// Multi-factor Authentication
	MFAEnabled     bool       `json:"mfa_enabled"`
	MFAMethods     string     `gorm:"type:json" json:"mfa_methods"` // JSON array
	MFASuccessRate float64    `json:"mfa_success_rate"`
	LastMFASuccess *time.Time `json:"last_mfa_success"`
	LastMFAFailure *time.Time `json:"last_mfa_failure"`

	// Biometric Verification
	BiometricEnabled    bool    `json:"biometric_enabled"`
	FingerprintVerified bool    `json:"fingerprint_verified"`
	FaceIDVerified      bool    `json:"face_id_verified"`
	VoiceVerified       bool    `json:"voice_verified"`
	BiometricScore      float64 `json:"biometric_score"`

	// Document Verification
	DocumentType     string     `json:"document_type"` // passport, driver_license, national_id
	DocumentVerified bool       `json:"document_verified"`
	DocumentNumber   string     `json:"document_number"` // Encrypted
	DocumentExpiry   *time.Time `json:"document_expiry"`
	DocumentScore    float64    `json:"document_score"`

	// Video Verification
	VideoVerified  bool    `json:"video_verified"`
	VideoSessionID string  `json:"video_session_id"`
	LivenessCheck  bool    `json:"liveness_check"`
	VideoScore     float64 `json:"video_score"`

	// Behavioral Biometrics
	TypingPattern     string  `gorm:"type:json" json:"typing_pattern"` // JSON object
	SwipePattern      string  `gorm:"type:json" json:"swipe_pattern"`  // JSON object
	BehavioralScore   float64 `json:"behavioral_score"`
	BehaviorDeviation float64 `json:"behavior_deviation"`

	// Device Fingerprinting
	DeviceFingerprint     string    `json:"device_fingerprint"`
	FingerprintStable     bool      `json:"fingerprint_stable"`
	FingerprintChanges    int       `json:"fingerprint_changes"`
	LastFingerprintUpdate time.Time `json:"last_fingerprint_update"`

	// SIM & Network
	SIMSwapDetected     bool       `json:"sim_swap_detected"`
	SIMSwapDate         *time.Time `json:"sim_swap_date"`
	CarrierVerified     bool       `json:"carrier_verified"`
	PhoneNumberVerified bool       `json:"phone_number_verified"`

	// Verification Summary
	OverallConfidence  float64 `json:"overall_confidence"`               // 0-100
	VerificationStatus string  `json:"verification_status"`              // verified, partial, failed, pending
	FailureReasons     string  `gorm:"type:json" json:"failure_reasons"` // JSON array
	RetryCount         int     `json:"retry_count"`
	MaxRetriesReached  bool    `json:"max_retries_reached"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// DeviceFraudInvestigation manages fraud investigation cases
type DeviceFraudInvestigation struct {
	database.BaseModel
	DeviceID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id"`
	CaseNumber string     `gorm:"uniqueIndex" json:"case_number"`
	OpenedDate time.Time  `json:"opened_date"`
	ClosedDate *time.Time `json:"closed_date"`

	// Case Management
	CaseStatus     string     `gorm:"type:varchar(50)" json:"case_status"` // open, investigating, pending, closed
	CasePriority   string     `json:"case_priority"`                       // critical, high, medium, low
	CaseType       string     `json:"case_type"`                           // claim_fraud, identity_fraud, device_fraud
	InvestigatorID *uuid.UUID `gorm:"type:uuid" json:"investigator_id"`
	TeamID         *uuid.UUID `gorm:"type:uuid" json:"team_id"`

	// Evidence Collection
	EvidenceItems      string `gorm:"type:json" json:"evidence_items"`      // JSON array
	DocumentsCollected string `gorm:"type:json" json:"documents_collected"` // JSON array
	WitnessStatements  string `gorm:"type:json" json:"witness_statements"`  // JSON array
	DigitalEvidence    string `gorm:"type:json" json:"digital_evidence"`    // JSON array

	// Investigation Timeline
	TimelineEvents     string `gorm:"type:json" json:"timeline_events"` // JSON array
	KeyFindings        string `gorm:"type:text" json:"key_findings"`
	InvestigationNotes string `gorm:"type:text" json:"investigation_notes"`
	NextSteps          string `gorm:"type:json" json:"next_steps"` // JSON array

	// Outcome
	InvestigationResult string  `json:"investigation_result"` // confirmed_fraud, false_positive, inconclusive
	FraudConfirmed      bool    `json:"fraud_confirmed"`
	FraudAmount         float64 `json:"fraud_amount"`
	RecoveredAmount     float64 `json:"recovered_amount"`

	// Legal Action
	LegalActionTaken       bool   `json:"legal_action_taken"`
	LegalCaseNumber        string `json:"legal_case_number"`
	LawEnforcementNotified bool   `json:"law_enforcement_notified"`
	PoliceReportNumber     string `json:"police_report_number"`
	ProsecutionStatus      string `json:"prosecution_status"`

	// Recovery Efforts
	RecoveryInitiated bool   `json:"recovery_initiated"`
	RecoveryMethod    string `json:"recovery_method"`
	RecoveryStatus    string `json:"recovery_status"`
	RecoveryTimeline  string `gorm:"type:json" json:"recovery_timeline"` // JSON array

	// Cost Tracking
	InvestigationCost float64 `json:"investigation_cost"`
	LegalCost         float64 `json:"legal_cost"`
	RecoveryCost      float64 `json:"recovery_cost"`
	TotalCost         float64 `json:"total_cost"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// Investigator should be loaded via service layer using InvestigatorID to avoid circular import
}

// DeviceBlacklistManagement handles device blacklist status across systems
type DeviceBlacklistManagement struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	IMEI     string    `gorm:"index" json:"imei"`

	// Blacklist Status
	BlacklistStatus string     `gorm:"type:varchar(50)" json:"blacklist_status"` // clean, blacklisted, graylisted, pending
	BlacklistDate   *time.Time `json:"blacklist_date"`
	BlacklistReason string     `json:"blacklist_reason"`
	BlacklistSource string     `json:"blacklist_source"` // internal, carrier, interpol, gsma

	// Database Checks
	GlobalBlacklist   bool   `json:"global_blacklist"`
	RegionalBlacklist string `gorm:"type:json" json:"regional_blacklist"` // JSON array of regions
	CarrierBlacklist  string `gorm:"type:json" json:"carrier_blacklist"`  // JSON array of carriers
	INTERPOLStatus    string `json:"interpol_status"`
	GSMAStatus        string `json:"gsma_status"`
	PrivateDBStatus   string `gorm:"type:json" json:"private_db_status"` // JSON object

	// Reason Codes
	ReasonCode     string    `json:"reason_code"` // THEFT, FRAUD, LOST, NON_PAYMENT
	DetailedReason string    `gorm:"type:text" json:"detailed_reason"`
	ReportingParty string    `json:"reporting_party"`
	ReportDate     time.Time `json:"report_date"`

	// Removal Management
	RemovalRequested   bool       `json:"removal_requested"`
	RemovalRequestDate *time.Time `json:"removal_request_date"`
	RemovalReason      string     `json:"removal_reason"`
	RemovalApproved    bool       `json:"removal_approved"`
	RemovalDate        *time.Time `json:"removal_date"`

	// False Blacklist Detection
	FalseBlacklistCheck    bool       `json:"false_blacklist_check"`
	FalsePositiveSuspected bool       `json:"false_positive_suspected"`
	VerificationStatus     string     `json:"verification_status"`
	VerificationDate       *time.Time `json:"verification_date"`

	// Impact Assessment
	ImpactLevel        string  `json:"impact_level"`                       // critical, high, medium, low
	ServicesAffected   string  `gorm:"type:json" json:"services_affected"` // JSON array
	FinancialImpact    float64 `json:"financial_impact"`
	ReputationalImpact string  `json:"reputational_impact"`

	// Cross-border Sync
	SyncStatus      string     `json:"sync_status"` // synced, pending, failed
	LastSyncDate    *time.Time `json:"last_sync_date"`
	SyncedDatabases string     `gorm:"type:json" json:"synced_databases"` // JSON array
	SyncErrors      string     `gorm:"type:json" json:"sync_errors"`      // JSON array

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceFraudPrevention tracks prevention measures and effectiveness
type DeviceFraudPrevention struct {
	database.BaseModel
	DeviceID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	LastAssessmentDate time.Time `json:"last_assessment_date"`

	// Prevention Measures
	ActiveMeasures       string  `gorm:"type:json" json:"active_measures"`       // JSON array
	MeasureEffectiveness string  `gorm:"type:json" json:"measure_effectiveness"` // JSON object with scores
	SecurityScore        float64 `json:"security_score"`                         // 0-100

	// Security Features
	SecurityFeatures     string `gorm:"type:json" json:"security_features"`     // JSON object
	FeatureUtilization   string `gorm:"type:json" json:"feature_utilization"`   // JSON object with usage %
	FeatureEffectiveness string `gorm:"type:json" json:"feature_effectiveness"` // JSON object

	// User Education
	EducationCompleted bool       `json:"education_completed"`
	EducationModules   string     `gorm:"type:json" json:"education_modules"` // JSON array
	EducationScore     float64    `json:"education_score"`
	LastEducationDate  *time.Time `json:"last_education_date"`

	// Awareness Scoring
	FraudAwarenessScore float64 `json:"fraud_awareness_score"`
	SecurityHygiene     float64 `json:"security_hygiene"`
	RiskBehaviorScore   float64 `json:"risk_behavior_score"`
	ImprovementAreas    string  `gorm:"type:json" json:"improvement_areas"` // JSON array

	// Effectiveness Metrics
	PreventedIncidents int     `json:"prevented_incidents"`
	DetectedThreats    int     `json:"detected_threats"`
	FalsePositives     int     `json:"false_positives"`
	MissedThreats      int     `json:"missed_threats"`
	PreventionRate     float64 `json:"prevention_rate"`

	// Cost-Benefit Analysis
	PreventionCost  float64 `json:"prevention_cost"`
	LossPrevented   float64 `json:"loss_prevented"`
	ROI             float64 `json:"roi"`
	CostPerIncident float64 `json:"cost_per_incident"`

	// Strategy Optimization
	CurrentStrategy       string    `gorm:"type:json" json:"current_strategy"`    // JSON object
	RecommendedChanges    string    `gorm:"type:json" json:"recommended_changes"` // JSON array
	OptimizationPotential float64   `json:"optimization_potential"`
	NextReviewDate        time.Time `json:"next_review_date"`

	// Risk Mitigation
	MitigationActions   string  `gorm:"type:json" json:"mitigation_actions"` // JSON array
	RiskReduction       float64 `json:"risk_reduction"`                      // Percentage reduction
	ResidualRisk        float64 `json:"residual_risk"`
	AcceptableRiskLevel float64 `json:"acceptable_risk_level"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceFraudPatterns
func (dfp *DeviceFraudPatterns) IsHighRisk() bool {
	return dfp.RiskLevel == "critical" || dfp.RiskLevel == "high"
}

func (dfp *DeviceFraudPatterns) CalculateRecoveryRate() float64 {
	if dfp.TotalLossAmount > 0 {
		dfp.RecoveryRate = (dfp.RecoveredAmount / dfp.TotalLossAmount) * 100
	}
	return dfp.RecoveryRate
}

func (dfp *DeviceFraudPatterns) IsPartOfRing() bool {
	return dfp.FraudRingID != nil
}

func (dfp *DeviceFraudPatterns) UpdateRiskLevel() {
	if dfp.MatchScore > 80 || dfp.CrossDeviceScore > 75 {
		dfp.RiskLevel = "critical"
	} else if dfp.MatchScore > 60 || dfp.CrossDeviceScore > 50 {
		dfp.RiskLevel = "high"
	} else if dfp.MatchScore > 40 {
		dfp.RiskLevel = "medium"
	} else {
		dfp.RiskLevel = "low"
	}
}

// Methods for DeviceAnomalyDetection
func (dad *DeviceAnomalyDetection) ShouldTriggerAlert() bool {
	return dad.AnomalyScore > dad.AlertThreshold && dad.FalsePositiveProb < 0.3
}

func (dad *DeviceAnomalyDetection) ShouldBlock() bool {
	return dad.AnomalyScore > dad.BlockThreshold && dad.AnomalySeverity == "critical"
}

func (dad *DeviceAnomalyDetection) RequiresInvestigation() bool {
	return dad.AnomalyScore > dad.InvestigateThreshold || dad.AnomalyCategory == "malicious"
}

func (dad *DeviceAnomalyDetection) UpdateThresholds(learningRate float64) {
	// Adaptive threshold adjustment based on false positive rate
	if dad.FalsePositiveProb > 0.5 && dad.AdaptiveThreshold {
		dad.AlertThreshold *= (1 + learningRate)
		dad.InvestigateThreshold *= (1 + learningRate*0.5)
	}
}

// Methods for DeviceIdentityVerification
func (div *DeviceIdentityVerification) IsFullyVerified() bool {
	return div.VerificationStatus == "verified" && div.OverallConfidence > 90
}

func (div *DeviceIdentityVerification) CalculateOverallConfidence() float64 {
	confidence := 0.0
	factors := 0

	if div.MFAEnabled {
		confidence += div.MFASuccessRate
		factors++
	}
	if div.BiometricEnabled {
		confidence += div.BiometricScore
		factors++
	}
	if div.DocumentVerified {
		confidence += div.DocumentScore
		factors++
	}
	if div.VideoVerified {
		confidence += div.VideoScore
		factors++
	}

	if factors > 0 {
		div.OverallConfidence = confidence / float64(factors)
	}

	// Apply penalties
	if div.SIMSwapDetected {
		div.OverallConfidence *= 0.5
	}
	if !div.FingerprintStable {
		div.OverallConfidence *= 0.8
	}

	return div.OverallConfidence
}

func (div *DeviceIdentityVerification) HasRecentSIMSwap() bool {
	if div.SIMSwapDate == nil {
		return false
	}
	return time.Since(*div.SIMSwapDate) < 7*24*time.Hour // Within last 7 days
}

// Methods for DeviceFraudInvestigation
func (dfi *DeviceFraudInvestigation) IsOpen() bool {
	return dfi.CaseStatus == "open" || dfi.CaseStatus == "investigating"
}

func (dfi *DeviceFraudInvestigation) CalculateTotalCost() float64 {
	dfi.TotalCost = dfi.InvestigationCost + dfi.LegalCost + dfi.RecoveryCost
	return dfi.TotalCost
}

func (dfi *DeviceFraudInvestigation) GetRecoveryRate() float64 {
	if dfi.FraudAmount > 0 {
		return (dfi.RecoveredAmount / dfi.FraudAmount) * 100
	}
	return 0
}

func (dfi *DeviceFraudInvestigation) DaysSinceOpened() int {
	return int(time.Since(dfi.OpenedDate).Hours() / 24)
}

// Methods for DeviceBlacklistManagement
func (dbm *DeviceBlacklistManagement) IsBlacklisted() bool {
	return dbm.BlacklistStatus == "blacklisted" || dbm.GlobalBlacklist
}

func (dbm *DeviceBlacklistManagement) CanBeRemoved() bool {
	return dbm.RemovalRequested && !dbm.RemovalApproved &&
		dbm.ReasonCode != "THEFT" && dbm.ReasonCode != "FRAUD"
}

func (dbm *DeviceBlacklistManagement) GetImpactSeverity() int {
	switch dbm.ImpactLevel {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

func (dbm *DeviceBlacklistManagement) NeedsSyncUpdate() bool {
	if dbm.LastSyncDate == nil {
		return true
	}
	return time.Since(*dbm.LastSyncDate) > 24*time.Hour
}

// Methods for DeviceFraudPrevention
func (dfp *DeviceFraudPrevention) CalculateROI() float64 {
	if dfp.PreventionCost > 0 {
		dfp.ROI = ((dfp.LossPrevented - dfp.PreventionCost) / dfp.PreventionCost) * 100
	}
	return dfp.ROI
}

func (dfp *DeviceFraudPrevention) CalculatePreventionRate() float64 {
	totalAttempts := dfp.PreventedIncidents + dfp.MissedThreats
	if totalAttempts > 0 {
		dfp.PreventionRate = (float64(dfp.PreventedIncidents) / float64(totalAttempts)) * 100
	}
	return dfp.PreventionRate
}

func (dfp *DeviceFraudPrevention) IsEffective() bool {
	return dfp.PreventionRate > 85 && dfp.ROI > 100
}

func (dfp *DeviceFraudPrevention) NeedsReview() bool {
	return time.Since(dfp.LastAssessmentDate) > 30*24*time.Hour ||
		dfp.SecurityScore < 60 ||
		dfp.PreventionRate < 70
}

func (dfp *DeviceFraudPrevention) GetRiskLevel() string {
	if dfp.ResidualRisk > dfp.AcceptableRiskLevel*1.5 {
		return "high"
	} else if dfp.ResidualRisk > dfp.AcceptableRiskLevel {
		return "medium"
	}
	return "low"
}
