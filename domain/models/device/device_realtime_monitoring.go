package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceMonitoringSession represents active real-time monitoring sessions
type DeviceMonitoringSession struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	SessionID      string    `gorm:"uniqueIndex;not null" json:"session_id"`

	// Session Details
	SessionType    string    `gorm:"type:varchar(30);not null" json:"session_type"` // diagnostic, continuous, event_driven, predictive
	Status         string    `gorm:"type:varchar(20);not null;default:'starting'" json:"status"` // starting, active, paused, stopped, error

	StartedAt      time.Time `gorm:"not null" json:"started_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
	Duration       int        `json:"duration,omitempty"` // seconds

	// Monitoring Configuration
	MonitoringRules string `gorm:"type:json" json:"monitoring_rules"` // JSON monitoring rules
	SamplingRate    int    `gorm:"default:60" json:"sampling_rate"`   // seconds between samples
	DataRetention   int    `gorm:"default:86400" json:"data_retention"` // seconds to keep data

	// Active Rules & Thresholds
	ActiveRules     string `gorm:"type:json" json:"active_rules,omitempty"` // JSON currently active monitoring rules
	Thresholds      string `gorm:"type:json" json:"thresholds,omitempty"` // JSON alert thresholds

	// Performance Metrics
	DataPointsCollected int64 `gorm:"default:0" json:"data_points_collected"`
	AlertsTriggered     int   `gorm:"default:0" json:"alerts_triggered"`
	ProcessingLatency   int   `json:"processing_latency,omitempty"` // milliseconds
	Throughput          float64 `json:"throughput,omitempty"` // data points per second

	// Error Tracking
	ErrorCount      int     `gorm:"default:0" json:"error_count"`
	LastError       string  `json:"last_error,omitempty"`
	ErrorRate       float64 `json:"error_rate,omitempty"` // errors per minute

	// Resource Usage
	MemoryUsage     float64 `json:"memory_usage,omitempty"` // MB
	CPUUsage        float64 `json:"cpu_usage,omitempty"`    // percentage
	NetworkUsage    float64 `json:"network_usage,omitempty"` // KB/s

	// Session Metadata
	InitiatedBy     uuid.UUID `gorm:"type:uuid;not null" json:"initiated_by"`
	InitiatedByName string    `json:"initiated_by_name,omitempty"`
	Priority        string    `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, critical

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	ActiveAlerts   []DeviceLiveAlert `gorm:"foreignKey:SessionID;references:SessionID" json:"active_alerts,omitempty"`
}

// DeviceRealTimeMetric represents live metric data points
type DeviceRealTimeMetric struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	SessionID      string    `gorm:"index" json:"session_id,omitempty"`

	// Metric Details
	MetricType     string    `gorm:"type:varchar(50);not null" json:"metric_type"` // cpu_usage, memory_usage, battery_level, etc.
	MetricName     string    `gorm:"not null" json:"metric_name"`
	MetricValue    float64   `gorm:"not null" json:"metric_value"`

	Timestamp      time.Time `gorm:"not null;index" json:"timestamp"`

	// Metric Metadata
	Unit           string    `json:"unit,omitempty"` // %, MB, V, etc.
	Quality        string    `gorm:"type:varchar(20);default:'good'" json:"quality"` // good, fair, poor
	Source         string    `json:"source,omitempty"` // sensor, calculation, external

	// Contextual Data
	ContextData    string    `gorm:"type:json" json:"context_data,omitempty"` // JSON additional context
	Tags           string    `gorm:"type:json" json:"tags,omitempty"` // JSON tags for categorization

	// Aggregation Flags
	IsAggregated   bool      `gorm:"default:false" json:"is_aggregated"`
	AggregationWindow int    `json:"aggregation_window,omitempty"` // seconds
	AggregationMethod string `json:"aggregation_method,omitempty"` // avg, min, max, sum
}

// DeviceLiveAlert represents active real-time alerts
type DeviceLiveAlert struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	SessionID      string    `gorm:"index" json:"session_id,omitempty"`
	AlertID        string    `gorm:"uniqueIndex;not null" json:"alert_id"`

	// Alert Classification
	AlertType      string    `gorm:"type:varchar(50);not null" json:"alert_type"` // performance, security, hardware, software
	AlertCategory  string    `gorm:"type:varchar(30);not null" json:"alert_category"` // warning, error, critical, info

	// Alert Details
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `gorm:"type:text;not null" json:"description"`
	Message        string    `gorm:"type:text" json:"message"`

	// Severity & Priority
	Severity       string    `gorm:"type:varchar(20);not null" json:"severity"` // low, medium, high, critical
	Priority       string    `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, urgent

	// Status Tracking
	Status         string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, acknowledged, resolved, suppressed
	Acknowledged   bool      `gorm:"default:false" json:"acknowledged"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	AcknowledgedBy *uuid.UUID `gorm:"type:uuid" json:"acknowledged_by,omitempty"`

	Resolved       bool      `gorm:"default:false" json:"resolved"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy     *uuid.UUID `gorm:"type:uuid" json:"resolved_by,omitempty"`

	// Trigger Information
	TriggeredAt    time.Time `gorm:"not null" json:"triggered_at"`
	TriggerSource  string    `json:"trigger_source,omitempty"` // monitoring_rule, threshold, anomaly_detection
	TriggerRule    string    `json:"trigger_rule,omitempty"` // rule that triggered alert

	// Trigger Values
	TriggerValue   float64   `json:"trigger_value,omitempty"`
	ThresholdValue float64   `json:"threshold_value,omitempty"`
	BaselineValue  float64   `json:"baseline_value,omitempty"`

	// Evidence & Context
	Evidence       string    `gorm:"type:json" json:"evidence,omitempty"` // JSON supporting data
	ContextData    string    `gorm:"type:json" json:"context_data,omitempty"` // JSON contextual information

	// Automated Actions
	AutoActions    string    `gorm:"type:json" json:"auto_actions,omitempty"` // JSON automated responses taken
	RecommendedActions string `gorm:"type:json" json:"recommended_actions,omitempty"` // JSON recommended actions

	// Escalation
	EscalationLevel int       `gorm:"default:0" json:"escalation_level"`
	EscalationRules string    `gorm:"type:json" json:"escalation_rules,omitempty"` // JSON escalation criteria
	NextEscalation  *time.Time `json:"next_escalation,omitempty"`

	// Notification Tracking
	NotificationsSent string `gorm:"type:json" json:"notifications_sent,omitempty"` // JSON notification history
	LastNotification  *time.Time `json:"last_notification,omitempty"`

	// Resolution Details
	ResolutionNotes string    `gorm:"type:text" json:"resolution_notes,omitempty"`
	ResolutionTime  int       `json:"resolution_time,omitempty"` // minutes to resolve
	ResolutionCost  float64   `json:"resolution_cost,omitempty"`

	// Impact Assessment
	BusinessImpact  string    `gorm:"type:varchar(20)" json:"business_impact"` // none, low, medium, high, critical
	UserImpact      string    `gorm:"type:varchar(20)" json:"user_impact"`     // none, low, medium, high, critical
	ServiceImpact   string    `gorm:"type:varchar(20)" json:"service_impact"`  // none, low, medium, high, critical

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	Session        *DeviceMonitoringSession `gorm:"foreignKey:SessionID;references:SessionID" json:"session,omitempty"`
}

// DeviceAlertConfiguration represents alert configuration rules
type DeviceAlertConfiguration struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ConfigID       string    `gorm:"uniqueIndex;not null" json:"config_id"`

	// Configuration Details
	ConfigName     string    `gorm:"not null" json:"config_name"`
	ConfigType     string    `gorm:"type:varchar(30);not null" json:"config_type"` // threshold, anomaly, pattern, composite
	Description    string    `gorm:"type:text" json:"description"`

	// Rule Definition
	RuleDefinition string    `gorm:"type:json;not null" json:"rule_definition"` // JSON rule logic
	RuleParameters string    `gorm:"type:json" json:"rule_parameters,omitempty"` // JSON rule parameters

	// Thresholds & Conditions
	Thresholds     string    `gorm:"type:json" json:"thresholds,omitempty"` // JSON threshold values
	Conditions     string    `gorm:"type:json" json:"conditions,omitempty"` // JSON condition logic

	// Alert Properties
	AlertTitle     string    `json:"alert_title,omitempty"` // template for alert title
	AlertMessage   string    `json:"alert_message,omitempty"` // template for alert message
	DefaultSeverity string   `gorm:"type:varchar(20);default:'medium'" json:"default_severity"`
	DefaultPriority string   `gorm:"type:varchar(20);default:'normal'" json:"default_priority"`

	// Behavior Settings
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	CooldownPeriod int       `gorm:"default:300" json:"cooldown_period"` // seconds between alerts
	MaxFrequency   int       `gorm:"default:10" json:"max_frequency"` // max alerts per hour
	SuppressionRules string  `gorm:"type:json" json:"suppression_rules,omitempty"` // JSON suppression logic

	// Notification Settings
	NotificationChannels string `gorm:"type:json" json:"notification_channels,omitempty"` // JSON notification methods
	EscalationRules string `gorm:"type:json" json:"escalation_rules,omitempty"` // JSON escalation logic
	AutoActions     string    `gorm:"type:json" json:"auto_actions,omitempty"` // JSON automated responses

	// Performance & Reliability
	FalsePositiveRate float64 `json:"false_positive_rate,omitempty"` // historical false positive rate
	AlertAccuracy     float64 `json:"alert_accuracy,omitempty"`      // historical accuracy
	ResponseTime      int     `json:"response_time,omitempty"`       // average response time in seconds

	// Maintenance
	LastTested      *time.Time `json:"last_tested,omitempty"`
	TestResults     string     `gorm:"type:json" json:"test_results,omitempty"` // JSON test results
	NeedsReview     bool       `gorm:"default:false" json:"needs_review"`

	// Audit
	CreatedBy       uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	LastModifiedBy  *uuid.UUID `gorm:"type:uuid" json:"last_modified_by,omitempty"`
	LastModifiedAt  time.Time `gorm:"autoUpdateTime" json:"last_modified_at"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceAlertHistory represents historical alert data for analytics
type DeviceAlertHistory struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AlertID        string    `gorm:"not null;index" json:"alert_id"`
	ConfigID       string    `gorm:"index" json:"config_id,omitempty"`

	// Alert Details (snapshot at time of alert)
	AlertType      string    `gorm:"type:varchar(50);not null" json:"alert_type"`
	Severity       string    `gorm:"type:varchar(20);not null" json:"severity"`
	Priority       string    `gorm:"type:varchar(20)" json:"priority"`
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`

	// Timing
	TriggeredAt    time.Time `gorm:"not null;index" json:"triggered_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`

	// Resolution Metrics
	TimeToAcknowledge int `json:"time_to_acknowledge,omitempty"` // seconds
	TimeToResolve     int `json:"time_to_resolve,omitempty"`     // seconds
	WasEscalated      bool `gorm:"default:false" json:"was_escalated"`

	// Impact & Cost
	BusinessImpact  string  `json:"business_impact,omitempty"`
	UserImpact      string  `json:"user_impact,omitempty"`
	EstimatedCost   float64 `json:"estimated_cost,omitempty"`
	ActualCost      float64 `json:"actual_cost,omitempty"`

	// Quality Assessment
	WasValidAlert   bool    `json:"was_valid_alert"` // true if alert was legitimate
	FalsePositive   bool    `json:"false_positive"`
	AlertAccuracy   float64 `json:"alert_accuracy,omitempty"` // 0-100

	// Response Details
	ResponseActions string  `gorm:"type:json" json:"response_actions,omitempty"` // JSON actions taken
	ResponseQuality float64 `json:"response_quality,omitempty"` // 0-100, quality of response

	// Learning Data
	SimilarAlerts   string  `gorm:"type:json" json:"similar_alerts,omitempty"` // JSON similar historical alerts
	Patterns        string  `gorm:"type:json" json:"patterns,omitempty"` // JSON patterns identified

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	Config         *DeviceAlertConfiguration `gorm:"foreignKey:ConfigID;references:ConfigID" json:"config,omitempty"`
}

// DeviceEscalationRule represents alert escalation rules
type DeviceEscalationRule struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RuleID         string    `gorm:"uniqueIndex;not null" json:"rule_id"`

	// Rule Definition
	RuleName       string    `gorm:"not null" json:"rule_name"`
	Description    string    `gorm:"type:text" json:"description"`
	Enabled        bool      `gorm:"default:true" json:"enabled"`

	// Trigger Conditions
	TriggerConditions string `gorm:"type:json;not null" json:"trigger_conditions"` // JSON trigger logic
	SeverityThreshold string `gorm:"type:varchar(20)" json:"severity_threshold"` // minimum severity to trigger

	// Escalation Steps
	EscalationSteps string `gorm:"type:json;not null" json:"escalation_steps"` // JSON escalation sequence
	TimeDelays      string `gorm:"type:json" json:"time_delays,omitempty"` // JSON time delays between steps

	// Notification Targets
	PrimaryContacts   string `gorm:"type:json" json:"primary_contacts,omitempty"` // JSON primary notification targets
	SecondaryContacts string `gorm:"type:json" json:"secondary_contacts,omitempty"` // JSON secondary targets
	EscalationChannels string `gorm:"type:json" json:"escalation_channels,omitempty"` // JSON notification channels

	// Business Hours
	BusinessHoursOnly bool   `gorm:"default:false" json:"business_hours_only"`
	TimeZone          string `json:"time_zone,omitempty"`

	// Performance Tracking
	TimesTriggered     int       `gorm:"default:0" json:"times_triggered"`
	LastTriggered      *time.Time `json:"last_triggered,omitempty"`
	AverageResponseTime int      `json:"average_response_time,omitempty"` // seconds

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// TableName returns the table name for DeviceMonitoringSession
func (DeviceMonitoringSession) TableName() string {
	return "device_monitoring_sessions"
}

// TableName returns the table name for DeviceRealTimeMetric
func (DeviceRealTimeMetric) TableName() string {
	return "device_realtime_metrics"
}

// TableName returns the table name for DeviceLiveAlert
func (DeviceLiveAlert) TableName() string {
	return "device_live_alerts"
}

// TableName returns the table name for DeviceAlertConfiguration
func (DeviceAlertConfiguration) TableName() string {
	return "device_alert_configurations"
}

// TableName returns the table name for DeviceAlertHistory
func (DeviceAlertHistory) TableName() string {
	return "device_alert_history"
}

// TableName returns the table name for DeviceEscalationRule
func (DeviceEscalationRule) TableName() string {
	return "device_escalation_rules"
}
