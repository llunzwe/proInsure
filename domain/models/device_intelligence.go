package models

import (
	"time"

	"github.com/google/uuid"
	"smartsure/pkg/database"
)

// IoTConnection represents an IoT device connection
type IoTConnection struct {
	database.BaseModel

	DeviceID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"device_id"`
	ConnectionID   string             `gorm:"uniqueIndex;not null" json:"connection_id"`
	Protocol       string             `gorm:"type:varchar(20);not null" json:"protocol"` // mqtt, websocket, http
	Status         string             `gorm:"type:varchar(20);not null" json:"status"`   // connected, disconnected, error
	IPAddress      string             `gorm:"not null" json:"ip_address"`
	MACAddress     string             `json:"mac_address,omitempty"`
	FirmwareVersion string            `json:"firmware_version,omitempty"`
	LastHeartbeat  time.Time          `gorm:"not null" json:"last_heartbeat"`
	LastSeen       time.Time          `gorm:"not null" json:"last_seen"`
	ConnectedAt    time.Time          `gorm:"not null" json:"connected_at"`
	DisconnectedAt *time.Time         `json:"disconnected_at,omitempty"`

	// Connection metadata
	UserAgent      string             `json:"user_agent,omitempty"`
	Location       string             `json:"location,omitempty"`
	NetworkType    string             `json:"network_type,omitempty"` // wifi, cellular, ethernet

	// Security
	CertificateFingerprint string     `json:"certificate_fingerprint,omitempty"`
	EncryptionEnabled      bool       `gorm:"default:true" json:"encryption_enabled"`

	// Relationships
	Device         Device             `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// IoTSensorData represents sensor data from IoT devices
type IoTSensorData struct {
	database.BaseModel

	DeviceID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"device_id"`
	SensorID       string             `gorm:"not null;index" json:"sensor_id"`
	SensorType     string             `gorm:"type:varchar(50);not null" json:"sensor_type"` // accelerometer, gyroscope, gps, etc.

	Value          float64            `gorm:"not null" json:"value"`
	Unit           string             `gorm:"type:varchar(20)" json:"unit,omitempty"`
	Timestamp      time.Time          `gorm:"not null;index" json:"timestamp"`

	// Sensor metadata
	Accuracy       float64            `json:"accuracy,omitempty"`
	Precision      float64            `json:"precision,omitempty"`
	RawData        string             `gorm:"type:jsonb" json:"raw_data,omitempty"`

	// Quality indicators
	IsValid        bool               `gorm:"default:true" json:"is_valid"`
	QualityScore   float64            `gorm:"default:1.0" json:"quality_score"`

	// Relationships
	Device         Device             `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// PredictiveMaintenanceAlert represents a predictive maintenance alert
type PredictiveMaintenanceAlert struct {
	database.BaseModel

	DeviceID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"device_id"`
	AlertID        string             `gorm:"uniqueIndex;not null" json:"alert_id"`
	AlertType      string             `gorm:"type:varchar(50);not null" json:"alert_type"` // battery, screen, hardware, etc.

	Severity       string             `gorm:"type:varchar(20);not null" json:"severity"` // low, medium, high, critical
	Status         string             `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, resolved, dismissed

	Title          string             `gorm:"not null" json:"title"`
	Description    string             `gorm:"type:text" json:"description"`
	PredictedFailureDate *time.Time   `json:"predicted_failure_date,omitempty"`

	// Prediction data
	ConfidenceScore float64           `gorm:"not null" json:"confidence_score"`
	PredictionModel string            `json:"prediction_model,omitempty"`
	AlgorithmVersion string           `json:"algorithm_version,omitempty"`

	// Sensor data that triggered the alert
	TriggeringSensors string          `gorm:"type:jsonb" json:"triggering_sensors,omitempty"`
	TriggeringValues  string          `gorm:"type:jsonb" json:"triggering_values,omitempty"`

	// Resolution
	ResolvedAt     *time.Time         `json:"resolved_at,omitempty"`
	ResolutionNotes string            `gorm:"type:text" json:"resolution_notes,omitempty"`
	ResolvedBy     *uuid.UUID         `gorm:"type:uuid" json:"resolved_by,omitempty"`

	// Actions taken
	ActionsTaken   string             `gorm:"type:jsonb" json:"actions_taken,omitempty"`

	// Relationships
	Device         Device             `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// DeviceHealthScore represents the overall health score of a device
type DeviceHealthScore struct {
	database.BaseModel

	DeviceID       uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	OverallScore   float64            `gorm:"not null" json:"overall_score"` // 0-100
	ScoreBreakdown string             `gorm:"type:jsonb" json:"score_breakdown"` // JSON object with component scores

	LastCalculated time.Time          `gorm:"not null" json:"last_calculated"`

	// Component scores
	BatteryHealth  float64            `json:"battery_health"`
	ScreenHealth   float64            `json:"screen_health"`
	HardwareHealth float64            `json:"hardware_health"`
	SoftwareHealth float64            `json:"software_health"`
	UsageHealth    float64            `json:"usage_health"`

	// Risk assessment
	RiskLevel      string             `gorm:"type:varchar(20)" json:"risk_level"` // low, medium, high, critical
	RiskFactors    string             `gorm:"type:jsonb" json:"risk_factors"`

	// Recommendations
	Recommendations string            `gorm:"type:jsonb" json:"recommendations"`

	// Relationships
	Device         Device             `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// DeviceUsagePattern represents usage patterns of a device
type DeviceUsagePattern struct {
	database.BaseModel

	DeviceID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"device_id"`
	PatternID      string             `gorm:"uniqueIndex;not null" json:"pattern_id"`

	PatternType    string             `gorm:"type:varchar(50);not null" json:"pattern_type"` // daily, weekly, anomaly
	TimeRange      string             `gorm:"type:varchar(20);not null" json:"time_range"`   // hour, day, week, month

	StartDate      time.Time          `gorm:"not null" json:"start_date"`
	EndDate        time.Time          `gorm:"not null" json:"end_date"`

	// Usage metrics
	AverageUsageHours float64         `json:"average_usage_hours"`
	PeakUsageHours    float64         `json:"peak_usage_hours"`
	UsageFrequency    float64         `json:"usage_frequency"`

	// Pattern data
	PatternData    string             `gorm:"type:jsonb" json:"pattern_data"`
	AnomaliesDetected bool           `gorm:"default:false" json:"anomalies_detected"`

	// Statistical data
	Mean           float64            `json:"mean"`
	StdDev         float64            `json:"std_dev"`
	MinValue       float64            `json:"min_value"`
	MaxValue       float64            `json:"max_value"`

	// Relationships
	Device         Device             `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// IoTCommand represents a command sent to an IoT device
type IoTCommand struct {
	database.BaseModel

	DeviceID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"device_id"`
	CommandID      string             `gorm:"uniqueIndex;not null" json:"command_id"`
	CommandType    string             `gorm:"type:varchar(50);not null" json:"command_type"` // ping, update, reboot, etc.

	Status         string             `gorm:"type:varchar(20);not null" json:"status"` // pending, sent, delivered, executed, failed
	Priority       string             `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, urgent

	CommandData    string             `gorm:"type:jsonb" json:"command_data"`
	ResponseData   string             `gorm:"type:jsonb" json:"response_data,omitempty"`

	SentAt         *time.Time         `json:"sent_at,omitempty"`
	DeliveredAt    *time.Time         `json:"delivered_at,omitempty"`
	ExecutedAt     *time.Time         `json:"executed_at,omitempty"`
	FailedAt       *time.Time         `json:"failed_at,omitempty"`

	ErrorMessage   string             `gorm:"type:text" json:"error_message,omitempty"`
	RetryCount     int                `gorm:"default:0" json:"retry_count"`
	MaxRetries     int                `gorm:"default:3" json:"max_retries"`

	// Audit
	InitiatedBy    uuid.UUID          `gorm:"type:uuid;not null" json:"initiated_by"`

	// Relationships
	Device         Device             `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// TableName returns the table name for IoTConnection
func (IoTConnection) TableName() string {
	return "iot_connections"
}

// TableName returns the table name for IoTSensorData
func (IoTSensorData) TableName() string {
	return "iot_sensor_data"
}

// TableName returns the table name for PredictiveMaintenanceAlert
func (PredictiveMaintenanceAlert) TableName() string {
	return "predictive_maintenance_alerts"
}

// TableName returns the table name for DeviceHealthScore
func (DeviceHealthScore) TableName() string {
	return "device_health_scores"
}

// TableName returns the table name for DeviceUsagePattern
func (DeviceUsagePattern) TableName() string {
	return "device_usage_patterns"
}

// TableName returns the table name for IoTCommand
func (IoTCommand) TableName() string {
	return "iot_commands"
}
