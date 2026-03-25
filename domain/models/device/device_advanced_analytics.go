package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceMLModel represents machine learning models for device analytics
type DeviceMLModel struct {
	database.BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ModelID  string    `gorm:"uniqueIndex;not null" json:"model_id"`

	// Model Information
	ModelName   string `gorm:"not null" json:"model_name"`
	ModelType   string `gorm:"type:varchar(30);not null" json:"model_type"` // classification, regression, clustering, anomaly_detection
	Description string `gorm:"type:text" json:"description"`

	// Model Metadata
	Algorithm string `json:"algorithm,omitempty"` // specific algorithm used
	Framework string `json:"framework,omitempty"` // tensorflow, pytorch, scikit-learn
	Version   string `gorm:"default:'1.0.0'" json:"version"`

	// Model Performance
	Accuracy  float64 `json:"accuracy,omitempty"`  // 0-100
	Precision float64 `json:"precision,omitempty"` // 0-100
	Recall    float64 `json:"recall,omitempty"`    // 0-100
	F1Score   float64 `json:"f1_score,omitempty"`  // 0-100
	AUC       float64 `json:"auc,omitempty"`       // area under curve
	MAE       float64 `json:"mae,omitempty"`       // mean absolute error
	MSE       float64 `json:"mse,omitempty"`       // mean squared error

	// Training Information
	TrainingDataSize int       `json:"training_data_size,omitempty"`
	TrainingDuration int       `json:"training_duration,omitempty"` // seconds
	TrainingCost     float64   `json:"training_cost,omitempty"`     // computational cost
	TrainedAt        time.Time `gorm:"not null" json:"trained_at"`

	// Feature Information
	Features           string `gorm:"type:json" json:"features,omitempty"`            // JSON list of features used
	FeatureImportance  string `gorm:"type:json" json:"feature_importance,omitempty"`  // JSON feature importance scores
	FeatureEngineering string `gorm:"type:json" json:"feature_engineering,omitempty"` // JSON feature engineering steps

	// Model Artifacts
	ModelArtifactPath   string `json:"model_artifact_path,omitempty"`  // path to stored model
	ModelSize           int64  `json:"model_size,omitempty"`           // size in bytes
	SerializationFormat string `json:"serialization_format,omitempty"` // pickle, onnx, h5

	// Validation & Testing
	CrossValidationScores string `gorm:"type:json" json:"cross_validation_scores,omitempty"` // JSON CV results
	TestResults           string `gorm:"type:json" json:"test_results,omitempty"`            // JSON test results
	ConfusionMatrix       string `gorm:"type:json" json:"confusion_matrix,omitempty"`        // JSON confusion matrix

	// Deployment Information
	IsDeployed            bool       `gorm:"default:false" json:"is_deployed"`
	DeployedAt            *time.Time `json:"deployed_at,omitempty"`
	DeploymentEnvironment string     `json:"deployment_environment,omitempty"` // production, staging, development

	// Monitoring & Maintenance
	LastUsed       *time.Time `json:"last_used,omitempty"`
	UsageCount     int        `gorm:"default:0" json:"usage_count"`
	DriftDetected  bool       `gorm:"default:false" json:"drift_detected"`
	LastDriftCheck *time.Time `json:"last_drift_check,omitempty"`

	// Lifecycle Management
	Status           string     `gorm:"type:varchar(20);default:'active'" json:"status"` // active, retired, archived
	RetiredAt        *time.Time `json:"retired_at,omitempty"`
	ReplacementModel string     `json:"replacement_model,omitempty"` // ID of replacement model

	// Audit & Governance
	CreatedBy        uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	GovernanceStatus string     `gorm:"type:varchar(20)" json:"governance_status"` // pending, approved, rejected

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceABTest represents A/B testing frameworks for device features
type DeviceABTest struct {
	database.BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	TestID   string    `gorm:"uniqueIndex;not null" json:"test_id"`

	// Test Definition
	TestName    string `gorm:"not null" json:"test_name"`
	Description string `gorm:"type:text" json:"description"`
	TestType    string `gorm:"type:varchar(30);not null" json:"test_type"` // feature_flag, parameter, algorithm, ui

	// Test Configuration
	TestHypothesis   string `gorm:"type:text" json:"test_hypothesis"`
	PrimaryMetric    string `json:"primary_metric,omitempty"`                     // primary success metric
	SecondaryMetrics string `gorm:"type:json" json:"secondary_metrics,omitempty"` // JSON secondary metrics

	// Variants Definition
	Variants            string `gorm:"type:json;not null" json:"variants"`              // JSON test variants
	ControlVariant      string `json:"control_variant,omitempty"`                       // control variant ID
	VariantDistribution string `gorm:"type:json" json:"variant_distribution,omitempty"` // JSON traffic distribution

	// Targeting & Segmentation
	TargetCriteria    string `gorm:"type:json" json:"target_criteria,omitempty"`    // JSON targeting rules
	SegmentDefinition string `gorm:"type:json" json:"segment_definition,omitempty"` // JSON user segments
	ExclusionCriteria string `gorm:"type:json" json:"exclusion_criteria,omitempty"` // JSON exclusion rules

	// Test Execution
	Status          string     `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // draft, running, paused, completed, cancelled
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	PlannedDuration int        `json:"planned_duration,omitempty"` // days

	// Sample Size & Statistical Power
	RequiredSampleSize int     `json:"required_sample_size,omitempty"`
	CurrentSampleSize  int     `gorm:"default:0" json:"current_sample_size"`
	StatisticalPower   float64 `json:"statistical_power,omitempty"` // 0-100
	ConfidenceLevel    float64 `json:"confidence_level,omitempty"`  // 0-100

	// Results & Analysis
	TestResults         string `gorm:"type:json" json:"test_results,omitempty"`                // JSON test results
	StatisticalAnalysis string `gorm:"type:json" json:"statistical_analysis,omitempty"`        // JSON statistical analysis
	WinnerDetermination string `gorm:"type:varchar(20)" json:"winner_determination,omitempty"` // automatic, manual, inconclusive

	// Winning Variant
	WinnerVariant      string     `json:"winner_variant,omitempty"`
	ConfidenceInWinner float64    `json:"confidence_in_winner,omitempty"` // 0-100
	WinnerDeclaredAt   *time.Time `json:"winner_declared_at,omitempty"`

	// Impact Assessment
	BusinessImpact  string `gorm:"type:json" json:"business_impact,omitempty"`  // JSON business impact metrics
	UserImpact      string `gorm:"type:json" json:"user_impact,omitempty"`      // JSON user experience impact
	TechnicalImpact string `gorm:"type:json" json:"technical_impact,omitempty"` // JSON technical impact

	// Rollout Planning
	RolloutPlan       string  `gorm:"type:json" json:"rollout_plan,omitempty"`     // JSON rollout strategy
	RolloutPercentage float64 `json:"rollout_percentage,omitempty"`                // percentage of users
	RolloutSchedule   string  `gorm:"type:json" json:"rollout_schedule,omitempty"` // JSON rollout timeline

	// Learning & Insights
	KeyInsights     string `gorm:"type:json" json:"key_insights,omitempty"`    // JSON key learnings
	Recommendations string `gorm:"type:json" json:"recommendations,omitempty"` // JSON future recommendations
	FollowUpTests   string `gorm:"type:json" json:"follow_up_tests,omitempty"` // JSON suggested follow-up tests

	// Audit & Compliance
	CreatedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	RequiresApproval  bool       `gorm:"default:true" json:"requires_approval"`
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	EthicalReview     bool       `gorm:"default:false" json:"ethical_review"`
	EthicalApprovedAt *time.Time `json:"ethical_approved_at,omitempty"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceExperimentTracking represents experiment tracking and analytics
type DeviceExperimentTracking struct {
	database.BaseModel

	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ExperimentID string    `gorm:"not null;index" json:"experiment_id"`
	EventID      string    `gorm:"uniqueIndex;not null" json:"event_id"`

	// Experiment Context
	TestID      string `gorm:"index" json:"test_id,omitempty"`
	VariantID   string `gorm:"not null" json:"variant_id"`
	UserSegment string `json:"user_segment,omitempty"`

	// Event Information
	EventType     string `gorm:"type:varchar(50);not null" json:"event_type"` // page_view, click, conversion, error
	EventName     string `json:"event_name,omitempty"`
	EventCategory string `json:"event_category,omitempty"`

	// Event Data
	EventProperties string    `gorm:"type:json" json:"event_properties,omitempty"` // JSON event properties
	EventValue      float64   `json:"event_value,omitempty"`
	EventTimestamp  time.Time `gorm:"not null" json:"event_timestamp"`

	// User Context
	UserID         uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserProperties string    `gorm:"type:json" json:"user_properties,omitempty"` // JSON user properties

	// Technical Context
	DeviceInfo  string `gorm:"type:json" json:"device_info,omitempty"` // JSON device information
	AppVersion  string `json:"app_version,omitempty"`
	OSVersion   string `json:"os_version,omitempty"`
	NetworkType string `json:"network_type,omitempty"`

	// Attribution
	AttributionSource string `json:"attribution_source,omitempty"`
	CampaignID        string `json:"campaign_id,omitempty"`
	Channel           string `json:"channel,omitempty"`

	// Quality Assurance
	EventQuality     string  `gorm:"type:varchar(20);default:'valid'" json:"event_quality"` // valid, suspicious, invalid
	ValidationErrors string  `gorm:"type:json" json:"validation_errors,omitempty"`          // JSON validation issues
	DataCompleteness float64 `json:"data_completeness,omitempty"`                           // 0-100

	// Processing Status
	IsProcessed      bool       `gorm:"default:false;index" json:"is_processed"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
	ProcessingErrors string     `gorm:"type:json" json:"processing_errors,omitempty"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceFeatureFlag represents feature flags for gradual rollouts
type DeviceFeatureFlag struct {
	database.BaseModel

	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	FlagID   string    `gorm:"uniqueIndex;not null" json:"flag_id"`

	// Flag Definition
	FlagName    string `gorm:"not null" json:"flag_name"`
	Description string `gorm:"type:text" json:"description"`
	FlagType    string `gorm:"type:varchar(20);not null" json:"flag_type"` // release_toggle, experiment, permission, operational

	// Flag Configuration
	DefaultValue      bool    `gorm:"default:false" json:"default_value"`
	RolloutPercentage float64 `gorm:"default:0" json:"rollout_percentage"`                           // 0-100
	RolloutStrategy   string  `gorm:"type:varchar(30);default:'percentage'" json:"rollout_strategy"` // percentage, user_id, gradual

	// Targeting Rules
	TargetRules    string `gorm:"type:json" json:"target_rules,omitempty"`    // JSON targeting criteria
	SegmentRules   string `gorm:"type:json" json:"segment_rules,omitempty"`   // JSON user segment rules
	ExclusionRules string `gorm:"type:json" json:"exclusion_rules,omitempty"` // JSON exclusion criteria

	// Status & Lifecycle
	Status     string     `gorm:"type:varchar(20);default:'disabled'" json:"status"` // disabled, enabled, archived
	EnabledAt  *time.Time `json:"enabled_at,omitempty"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`

	// Rollout Tracking
	CurrentRollout  float64 `json:"current_rollout,omitempty"`                   // actual rollout percentage
	RolloutSchedule string  `gorm:"type:json" json:"rollout_schedule,omitempty"` // JSON rollout timeline
	AutoRollout     bool    `gorm:"default:false" json:"auto_rollout"`

	// Monitoring & Analytics
	UsageMetrics  string  `gorm:"type:json" json:"usage_metrics,omitempty"`  // JSON usage statistics
	ImpactMetrics string  `gorm:"type:json" json:"impact_metrics,omitempty"` // JSON impact analysis
	ErrorRate     float64 `json:"error_rate,omitempty"`                      // error rate with flag enabled

	// Safety & Rollback
	SafeToEnable     bool   `gorm:"default:true" json:"safe_to_enable"`
	RollbackPlan     string `gorm:"type:text" json:"rollback_plan,omitempty"`
	EmergencyDisable bool   `gorm:"default:false" json:"emergency_disable"`

	// Dependencies
	PrerequisiteFlags string `gorm:"type:json" json:"prerequisite_flags,omitempty"` // JSON required flags
	ConflictingFlags  string `gorm:"type:json" json:"conflicting_flags,omitempty"`  // JSON conflicting flags

	// Audit
	CreatedBy  uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	ApprovedBy *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceModelValidation represents model validation and performance tracking
type DeviceModelValidation struct {
	database.BaseModel

	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ValidationID string    `gorm:"uniqueIndex;not null" json:"validation_id"`
	ModelID      string    `gorm:"index" json:"model_id,omitempty"` // associated ML model

	// Validation Context
	ValidationType string `gorm:"type:varchar(30);not null" json:"validation_type"` // backtesting, cross_validation, holdout, online
	Description    string `gorm:"type:text" json:"description"`

	// Validation Period
	StartDate  time.Time `gorm:"not null" json:"start_date"`
	EndDate    time.Time `gorm:"not null" json:"end_date"`
	DataPoints int       `json:"data_points,omitempty"` // number of data points validated

	// Performance Metrics
	Accuracy  float64 `json:"accuracy,omitempty"`  // 0-100
	Precision float64 `json:"precision,omitempty"` // 0-100
	Recall    float64 `json:"recall,omitempty"`    // 0-100
	F1Score   float64 `json:"f1_score,omitempty"`  // 0-100
	AUC       float64 `json:"auc,omitempty"`
	MAE       float64 `json:"mae,omitempty"`
	RMSE      float64 `json:"rmse,omitempty"`

	// Validation Results
	TestResults      string `gorm:"type:json" json:"test_results,omitempty"`      // JSON detailed results
	ConfusionMatrix  string `gorm:"type:json" json:"confusion_matrix,omitempty"`  // JSON confusion matrix
	ResidualAnalysis string `gorm:"type:json" json:"residual_analysis,omitempty"` // JSON residual analysis

	// Model Stability
	StabilityScore float64 `json:"stability_score,omitempty"`                  // 0-100
	DriftDetection string  `gorm:"type:json" json:"drift_detection,omitempty"` // JSON drift analysis
	BiasAnalysis   string  `gorm:"type:json" json:"bias_analysis,omitempty"`   // JSON bias analysis

	// Comparative Analysis
	BaselineComparison string `gorm:"type:json" json:"baseline_comparison,omitempty"` // JSON comparison with baseline
	ChampionChallenger string `gorm:"type:json" json:"champion_challenger,omitempty"` // JSON champion-challenger analysis

	// Recommendations
	ModelRecommendations string `gorm:"type:json" json:"model_recommendations,omitempty"` // JSON model improvement suggestions
	DataQualityIssues    string `gorm:"type:json" json:"data_quality_issues,omitempty"`   // JSON data quality problems identified

	// Status & Conclusions
	Status      string    `gorm:"type:varchar(20);default:'completed'" json:"status"` // pending, running, completed, failed
	Conclusion  string    `gorm:"type:text" json:"conclusion,omitempty"`
	ValidatedBy uuid.UUID `gorm:"type:uuid;not null" json:"validated_by"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// TableName returns the table name for DeviceMLModel
func (DeviceMLModel) TableName() string {
	return "device_ml_models"
}

// TableName returns the table name for DeviceABTest
func (DeviceABTest) TableName() string {
	return "device_ab_tests"
}

// TableName returns the table name for DeviceExperimentTracking
func (DeviceExperimentTracking) TableName() string {
	return "device_experiment_tracking"
}

// TableName returns the table name for DeviceFeatureFlag
func (DeviceFeatureFlag) TableName() string {
	return "device_feature_flags"
}

// TableName returns the table name for DeviceModelValidation
func (DeviceModelValidation) TableName() string {
	return "device_model_validations"
}
