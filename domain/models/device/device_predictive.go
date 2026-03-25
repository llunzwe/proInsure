package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PredictiveMaintenanceAlert represents AI-powered maintenance predictions
type PredictiveMaintenanceAlert struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AlertID        string    `gorm:"uniqueIndex;not null" json:"alert_id"`

	// Alert Classification
	AlertType      string    `gorm:"type:varchar(50);not null" json:"alert_type"` // battery, screen, hardware, software, connectivity, performance
	SubType        string    `json:"sub_type,omitempty"` // specific component or issue
	Category       string    `gorm:"type:varchar(30);not null" json:"category"`  // preventive, predictive, corrective

	// Severity & Priority
	Severity       string    `gorm:"type:varchar(20);not null" json:"severity"` // low, medium, high, critical
	Priority       string    `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, urgent
	UrgencyScore   float64   `gorm:"not null" json:"urgency_score"` // 0-100, calculated score

	// Alert Status
	Status         string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, acknowledged, resolved, dismissed, false_positive
	AcknlowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	AcknlowledgedBy *uuid.UUID `gorm:"type:uuid" json:"acknowledged_by,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy     *uuid.UUID `gorm:"type:uuid" json:"resolved_by,omitempty"`

	// Prediction Details
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `gorm:"type:text;not null" json:"description"`
	PredictedFailureDate *time.Time `json:"predicted_failure_date,omitempty"`
	TimeToFailure  int       `json:"time_to_failure,omitempty"` // days until predicted failure
	ConfidenceScore float64  `gorm:"not null" json:"confidence_score"` // 0-100, AI confidence

	// AI Model Information
	PredictionModel string    `json:"prediction_model,omitempty"` // model name/version
	AlgorithmVersion string   `json:"algorithm_version,omitempty"`
	FeatureImportance string  `gorm:"type:json" json:"feature_importance,omitempty"` // JSON feature weights

	// Triggering Data
	TriggeringSensors []string `gorm:"type:json" json:"triggering_sensors,omitempty"` // sensors that triggered alert
	TriggeringValues  string   `gorm:"type:json" json:"triggering_values,omitempty"`  // JSON sensor values that triggered
	TriggeringPatterns string `gorm:"type:json" json:"triggering_patterns,omitempty"` // JSON patterns detected
	DataTimeframe     string   `json:"data_timeframe,omitempty"` // time period analyzed

	// Evidence & Supporting Data
	SupportingEvidence string `gorm:"type:json" json:"supporting_evidence,omitempty"` // JSON evidence data
	HistoricalPatterns string `gorm:"type:json" json:"historical_patterns,omitempty"` // JSON similar past cases
	ComparativeAnalysis string `gorm:"type:json" json:"comparative_analysis,omitempty"` // JSON comparison with similar devices

	// Recommended Actions
	RecommendedActions string  `gorm:"type:json" json:"recommended_actions,omitempty"` // JSON action recommendations
	ActionPriority     string  `gorm:"type:varchar(20)" json:"action_priority,omitempty"`
	EstimatedCost      float64 `json:"estimated_cost,omitempty"`
	EstimatedDowntime  int     `json:"estimated_downtime,omitempty"` // minutes

	// Resolution Details
	ResolutionNotes    string  `gorm:"type:text" json:"resolution_notes,omitempty"`
	ActualResolution   string  `json:"actual_resolution,omitempty"`
	ResolutionCost     float64 `json:"resolution_cost,omitempty"`
	ResolutionTime     int     `json:"resolution_time,omitempty"` // minutes taken to resolve

	// Validation & Accuracy Tracking
	WasAccurate       bool     `json:"was_accurate,omitempty"`       // true if prediction was correct
	FalsePositiveReason string `gorm:"type:text" json:"false_positive_reason,omitempty"`
	AccuracyFeedback   string `json:"accuracy_feedback,omitempty"`

	// Escalation & Notification
	EscalationLevel   int      `gorm:"default:0" json:"escalation_level"`
	LastEscalation    *time.Time `json:"last_escalation,omitempty"`
	NotificationSent  bool     `gorm:"default:false" json:"notification_sent"`
	NotificationChannels []string `gorm:"type:json" json:"notification_channels,omitempty"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceHealthScore represents comprehensive device health assessment
type DeviceHealthScore struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Overall Health Assessment
	OverallScore   float64   `gorm:"not null" json:"overall_score"` // 0-100, composite score
	HealthGrade    string    `gorm:"type:varchar(5)" json:"health_grade"` // A+, A, A-, B+, B, B-, C+, C, C-, D, F
	HealthStatus   string    `gorm:"type:varchar(20)" json:"health_status"` // excellent, good, fair, poor, critical

	LastCalculated time.Time `gorm:"not null" json:"last_calculated"`
	NextAssessment time.Time `json:"next_assessment"`

	// Component Health Scores (0-100)
	BatteryHealth       float64 `json:"battery_health"`
	ScreenHealth        float64 `json:"screen_health"`
	HardwareHealth      float64 `json:"hardware_health"`
	SoftwareHealth      float64 `json:"software_health"`
	StorageHealth       float64 `json:"storage_health"`
	PerformanceHealth   float64 `json:"performance_health"`
	ConnectivityHealth  float64 `json:"connectivity_health"`
	SecurityHealth      float64 `json:"security_health"`
	ThermalHealth       float64 `json:"thermal_health"`
	UsageHealth         float64 `json:"usage_health"`

	// Health Score Breakdown
	ScoreBreakdown string `gorm:"type:json" json:"score_breakdown"` // JSON detailed breakdown
	ComponentWeights string `gorm:"type:json" json:"component_weights"` // JSON weightings used

	// Risk Assessment
	RiskLevel       string  `gorm:"type:varchar(20)" json:"risk_level"` // low, medium, high, critical
	RiskScore       float64 `json:"risk_score"` // 0-100
	RiskFactors     string  `gorm:"type:json" json:"risk_factors"` // JSON risk factors identified
	RiskTrends      string  `gorm:"type:json" json:"risk_trends"` // JSON risk trend analysis

	// Failure Predictions
	PredictedFailures string `gorm:"type:json" json:"predicted_failures,omitempty"` // JSON predicted failure modes
	TimeToFailure     string `gorm:"type:json" json:"time_to_failure,omitempty"` // JSON time estimates for failures

	// Recommendations
	Recommendations     string `gorm:"type:json" json:"recommendations,omitempty"` // JSON health improvement recommendations
	PriorityActions     string `gorm:"type:json" json:"priority_actions,omitempty"` // JSON urgent actions needed
	MaintenanceSchedule string `gorm:"type:json" json:"maintenance_schedule,omitempty"` // JSON recommended maintenance

	// Assessment Metadata
	AssessmentMethod string `json:"assessment_method"` // automated, manual, hybrid
	DataSources      string `gorm:"type:json" json:"data_sources"` // JSON sources used for assessment
	AssessmentVersion string `json:"assessment_version"`
	ConfidenceLevel   float64 `json:"confidence_level"` // 0-100, confidence in assessment

	// Historical Tracking
	ScoreHistory     string `gorm:"type:json" json:"score_history,omitempty"` // JSON historical scores
	TrendAnalysis    string `json:"trend_analysis,omitempty"` // improving, stable, declining

	// Benchmarking
	IndustryAverage  float64 `json:"industry_average,omitempty"`
	PeerComparison   string  `gorm:"type:json" json:"peer_comparison,omitempty"` // JSON comparison with similar devices

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceFailurePrediction represents specific failure mode predictions
type DeviceFailurePrediction struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	PredictionID   string    `gorm:"uniqueIndex;not null" json:"prediction_id"`

	// Failure Details
	FailureMode    string    `gorm:"type:varchar(50);not null" json:"failure_mode"` // battery_failure, screen_crack, hardware_fault, etc.
	FailureComponent string  `json:"failure_component,omitempty"`
	FailureSeverity string   `gorm:"type:varchar(20)" json:"failure_severity"` // minor, moderate, major, catastrophic

	// Prediction Details
	PredictedDate  *time.Time `json:"predicted_date,omitempty"`
	TimeToFailure  int        `json:"time_to_failure,omitempty"` // days
	ConfidenceScore float64   `gorm:"not null" json:"confidence_score"` // 0-100
	PredictionModel string    `json:"prediction_model,omitempty"`

	// Evidence & Reasoning
	TriggeringFactors string `gorm:"type:json" json:"triggering_factors,omitempty"` // JSON factors causing prediction
	SupportingData    string `gorm:"type:json" json:"supporting_data,omitempty"` // JSON data supporting prediction
	HistoricalPatterns string `gorm:"type:json" json:"historical_patterns,omitempty"` // JSON similar historical cases

	// Risk Assessment
	RiskImpact      string `gorm:"type:varchar(20)" json:"risk_impact"` // low, medium, high, critical
	RiskProbability float64 `json:"risk_probability"` // 0-100
	RiskCost        float64 `json:"risk_cost,omitempty"` // estimated financial impact

	// Prevention & Mitigation
	PreventionActions string `gorm:"type:json" json:"prevention_actions,omitempty"` // JSON preventive measures
	MitigationPlan    string `gorm:"type:json" json:"mitigation_plan,omitempty"` // JSON mitigation strategies
	BackupPlans       string `gorm:"type:json" json:"backup_plans,omitempty"` // JSON contingency plans

	// Status Tracking
	Status           string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active, prevented, occurred, false_positive
	Outcome          string    `json:"outcome,omitempty"`
	ActualFailureDate *time.Time `json:"actual_failure_date,omitempty"`
	PreventionCost   float64   `json:"prevention_cost,omitempty"`
	ActualDamageCost float64   `json:"actual_damage_cost,omitempty"`

	// Validation
	WasAccurate      bool     `json:"was_accurate,omitempty"`
	AccuracyRating   float64  `json:"accuracy_rating,omitempty"` // 0-100
	FeedbackNotes    string   `gorm:"type:text" json:"feedback_notes,omitempty"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	RelatedAlert   *PredictiveMaintenanceAlert `gorm:"foreignKey:AlertID;references:AlertID" json:"related_alert,omitempty"`
}

// DeviceMaintenanceSchedule represents AI-generated maintenance recommendations
type DeviceMaintenanceSchedule struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Schedule Overview
	ScheduleVersion string    `json:"schedule_version"`
	LastUpdated     time.Time `gorm:"not null" json:"last_updated"`
	NextReview      time.Time `json:"next_review"`

	// Maintenance Categories
	PreventiveMaintenance string `gorm:"type:json" json:"preventive_maintenance"` // JSON preventive tasks
	PredictiveMaintenance string `gorm:"type:json" json:"predictive_maintenance"` // JSON AI-recommended tasks
	CorrectiveMaintenance string `gorm:"type:json" json:"corrective_maintenance"` // JSON repair tasks

	// Schedule Details
	MaintenanceTasks string `gorm:"type:json" json:"maintenance_tasks"` // JSON detailed task list
	TaskPriorities   string `gorm:"type:json" json:"task_priorities"`   // JSON priority rankings
	EstimatedCosts   string `gorm:"type:json" json:"estimated_costs"`   // JSON cost estimates
	TimeEstimates    string `gorm:"type:json" json:"time_estimates"`    // JSON time requirements

	// Scheduling
	RecommendedSchedule string `gorm:"type:json" json:"recommended_schedule"` // JSON optimal timing
	FlexibilityWindows  string `gorm:"type:json" json:"flexibility_windows"`  // JSON scheduling flexibility
	DependencyMapping   string `gorm:"type:json" json:"dependency_mapping"`   // JSON task dependencies

	// Resource Requirements
	RequiredSkills     string `gorm:"type:json" json:"required_skills"`     // JSON skill requirements
	RequiredParts      string `gorm:"type:json" json:"required_parts"`      // JSON parts needed
	RequiredTools      string `gorm:"type:json" json:"required_tools"`      // JSON tools needed

	// Cost Optimization
	CostBenefitAnalysis string `gorm:"type:json" json:"cost_benefit_analysis"` // JSON cost-benefit data
	BudgetConstraints   string `gorm:"type:json" json:"budget_constraints"`   // JSON budget considerations
	PrioritizationLogic string `json:"prioritization_logic"` // how tasks are prioritized

	// Execution Tracking
	CompletedTasks     string `gorm:"type:json" json:"completed_tasks,omitempty"` // JSON completed tasks
	PendingTasks       string `gorm:"type:json" json:"pending_tasks,omitempty"`   // JSON pending tasks
	OverdueTasks       string `gorm:"type:json" json:"overdue_tasks,omitempty"`   // JSON overdue tasks

	// Performance Metrics
	ScheduleCompliance float64 `json:"schedule_compliance,omitempty"` // percentage of tasks completed on time
	CostVariance       float64 `json:"cost_variance,omitempty"`       // percentage over/under budget
	EffectivenessScore float64 `json:"effectiveness_score,omitempty"` // 0-100, how effective the schedule was

	// AI Model Information
	AIModelVersion     string  `json:"ai_model_version,omitempty"`
	RecommendationConfidence float64 `json:"recommendation_confidence,omitempty"` // 0-100

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceUpgradeRecommendation represents AI-powered upgrade suggestions
type DeviceUpgradeRecommendation struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecommendationID string  `gorm:"uniqueIndex;not null" json:"recommendation_id"`

	// Recommendation Details
	UpgradeType    string    `gorm:"type:varchar(50);not null" json:"upgrade_type"` // hardware, software, firmware, accessory
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`

	// Benefits & Impact
	PerformanceImprovement float64 `json:"performance_improvement,omitempty"` // percentage improvement
	BatteryLifeImprovement  float64 `json:"battery_life_improvement,omitempty"` // percentage improvement
	ReliabilityImprovement  float64 `json:"reliability_improvement,omitempty"` // percentage improvement
	SecurityEnhancement     bool    `json:"security_enhancement,omitempty"`
	NewFeatures            []string `gorm:"type:json" json:"new_features,omitempty"`

	// Cost Analysis
	UpgradeCost     float64 `json:"upgrade_cost,omitempty"`
	CostBenefitRatio float64 `json:"cost_benefit_ratio,omitempty"`
	PaybackPeriod   int     `json:"payback_period,omitempty"` // months
	ROI             float64 `json:"roi,omitempty"` // percentage

	// Timing & Urgency
	RecommendedDate *time.Time `json:"recommended_date,omitempty"`
	UrgencyLevel    string    `gorm:"type:varchar(20)" json:"urgency_level"` // low, medium, high, critical
	TimeSensitivity string    `json:"time_sensitivity,omitempty"` // time window for optimal upgrade

	// Technical Details
	CompatibilityCheck string `gorm:"type:json" json:"compatibility_check,omitempty"` // JSON compatibility analysis
	Prerequisites      string `gorm:"type:json" json:"prerequisites,omitempty"` // JSON requirements
	RiskAssessment     string `gorm:"type:json" json:"risk_assessment,omitempty"` // JSON upgrade risks

	// Implementation
	ImplementationSteps string `gorm:"type:json" json:"implementation_steps,omitempty"` // JSON step-by-step guide
	EstimatedDuration  int     `json:"estimated_duration,omitempty"` // minutes
	RequiredDowntime   int     `json:"required_downtime,omitempty"` // minutes

	// Status Tracking
	Status         string    `gorm:"type:varchar(20);default:'recommended'" json:"status"` // recommended, planned, in_progress, completed, cancelled
	ImplementedAt  *time.Time `json:"implemented_at,omitempty"`
	ImplementedBy  *uuid.UUID `gorm:"type:uuid" json:"implemented_by,omitempty"`
	Outcome        string    `gorm:"type:text" json:"outcome,omitempty"`

	// AI Model Information
	ConfidenceScore float64 `json:"confidence_score,omitempty"` // 0-100
	AIModelVersion  string  `json:"ai_model_version,omitempty"`
	RecommendationBasis string `gorm:"type:json" json:"recommendation_basis,omitempty"` // JSON reasoning

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceValueForecast represents AI-powered value forecasting
type DeviceValueForecast struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ForecastID     string    `gorm:"uniqueIndex;not null" json:"forecast_id"`

	// Forecast Details
	ForecastType   string    `gorm:"type:varchar(30);not null" json:"forecast_type"` // resale_value, insurance_value, trade_in_value
	ForecastDate   time.Time `gorm:"not null" json:"forecast_date"`
	ForecastHorizon int      `json:"forecast_horizon"` // months into future

	// Value Projections
	CurrentValue   float64   `json:"current_value"`
	ProjectedValue float64   `json:"projected_value"`
	ValueChange    float64   `json:"value_change"` // absolute change
	ValueChangePct float64   `json:"value_change_pct"` // percentage change

	// Forecast Confidence
	ConfidenceScore float64 `gorm:"not null" json:"confidence_score"` // 0-100
	ConfidenceInterval string `gorm:"type:json" json:"confidence_interval,omitempty"` // JSON confidence bounds
	MarginOfError   float64 `json:"margin_of_error,omitempty"`

	// Factors Influencing Forecast
	MarketConditions string `gorm:"type:json" json:"market_conditions,omitempty"` // JSON market factors
	DeviceCondition  string `gorm:"type:json" json:"device_condition,omitempty"` // JSON condition factors
	EconomicFactors  string `gorm:"type:json" json:"economic_factors,omitempty"` // JSON economic indicators

	// Historical Data
	HistoricalTrend string `gorm:"type:json" json:"historical_trend,omitempty"` // JSON historical value changes
	ComparableSales string `gorm:"type:json" json:"comparable_sales,omitempty"` // JSON similar device sales

	// Model Information
	AIModelVersion  string  `json:"ai_model_version,omitempty"`
	ForecastModel   string  `json:"forecast_model,omitempty"` // linear_regression, time_series, neural_network
	FeatureWeights  string  `gorm:"type:json" json:"feature_weights,omitempty"` // JSON feature importance

	// Validation
	BacktestedAccuracy float64 `json:"backtested_accuracy,omitempty"` // 0-100
	ActualVsPredicted  string  `gorm:"type:json" json:"actual_vs_predicted,omitempty"` // JSON validation data

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// TableName returns the table name for PredictiveMaintenanceAlert
func (PredictiveMaintenanceAlert) TableName() string {
	return "predictive_maintenance_alerts"
}

// TableName returns the table name for DeviceHealthScore
func (DeviceHealthScore) TableName() string {
	return "device_health_scores"
}

// TableName returns the table name for DeviceFailurePrediction
func (DeviceFailurePrediction) TableName() string {
	return "device_failure_predictions"
}

// TableName returns the table name for DeviceMaintenanceSchedule
func (DeviceMaintenanceSchedule) TableName() string {
	return "device_maintenance_schedules"
}

// TableName returns the table name for DeviceUpgradeRecommendation
func (DeviceUpgradeRecommendation) TableName() string {
	return "device_upgrade_recommendations"
}

// TableName returns the table name for DeviceValueForecast
func (DeviceValueForecast) TableName() string {
	return "device_value_forecasts"
}
