package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// IoTConnection represents an IoT device connection for real-time communication
type IoTConnection struct {
	database.BaseModel

	DeviceID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"device_id"`
	ConnectionID   string             `gorm:"uniqueIndex;not null" json:"connection_id"`

	// Connection Details
	Protocol       string             `gorm:"type:varchar(20);not null" json:"protocol"` // mqtt, websocket, http, coap
	Status         string             `gorm:"type:varchar(20);not null;default:'disconnected'" json:"status"` // connected, disconnected, error, reconnecting

	// Network Information
	IPAddress      string             `gorm:"not null" json:"ip_address"`
	MACAddress     string             `json:"mac_address,omitempty"`
	NetworkType    string             `json:"network_type,omitempty"` // wifi, cellular, ethernet, bluetooth
	CellTowerID    string             `json:"cell_tower_id,omitempty"`
	WiFiSSID       string             `json:"wifi_ssid,omitempty"`

	// Connection Metadata
	UserAgent      string             `json:"user_agent,omitempty"`
	FirmwareVersion string            `json:"firmware_version,omitempty"`
	SDKVersion     string             `json:"sdk_version,omitempty"`
	DeviceModel    string             `json:"device_model,omitempty"`

	// Connection Timing
	ConnectedAt    time.Time          `gorm:"not null" json:"connected_at"`
	LastHeartbeat  time.Time          `gorm:"not null" json:"last_heartbeat"`
	LastSeen       time.Time          `gorm:"not null" json:"last_seen"`
	DisconnectedAt *time.Time         `json:"disconnected_at,omitempty"`

	// Connection Quality
	SignalStrength int                `json:"signal_strength,omitempty"` // -100 to 0 for WiFi, 0-4 for cellular
	Latency        int                `json:"latency,omitempty"`        // milliseconds
	PacketLoss     float64            `json:"packet_loss,omitempty"`    // percentage 0-100
	Bandwidth      float64            `json:"bandwidth,omitempty"`      // Mbps

	// Security
	CertificateFingerprint string     `json:"certificate_fingerprint,omitempty"`
	EncryptionEnabled      bool       `gorm:"default:true" json:"encryption_enabled"`
	AuthenticationMethod   string     `json:"authentication_method,omitempty"` // jwt, certificate, api_key

	// Connection History
	ConnectionCount        int        `gorm:"default:1" json:"connection_count"`
	TotalConnectedTime     int64      `json:"total_connected_time"` // seconds
	AverageSessionDuration int64      `json:"average_session_duration"` // seconds
	LongestSessionDuration int64      `json:"longest_session_duration"` // seconds
	ReconnectionCount      int        `gorm:"default:0" json:"reconnection_count"`

	// Data Transfer
	TotalDataSent     int64 `json:"total_data_sent"`     // bytes
	TotalDataReceived int64 `json:"total_data_received"` // bytes
	DataTransferRate  float64 `json:"data_transfer_rate"` // bytes per second

	// Error Tracking
	LastError          string     `json:"last_error,omitempty"`
	ErrorCount         int        `gorm:"default:0" json:"error_count"`
	LastErrorTime      *time.Time `json:"last_error_time,omitempty"`

	// Configuration
	HeartbeatInterval int        `gorm:"default:30" json:"heartbeat_interval"` // seconds
	ReconnectDelay    int        `gorm:"default:5" json:"reconnect_delay"`    // seconds
	MaxRetries        int        `gorm:"default:3" json:"max_retries"`

	// Relationships
	SensorDataStreams  []IoTSensorDataStream `gorm:"foreignKey:ConnectionID;references:ConnectionID" json:"sensor_data_streams,omitempty"`
	ActiveCommands     []IoTCommand          `gorm:"foreignKey:ConnectionID;references:ConnectionID;where:status IN ('pending','sent')" json:"active_commands,omitempty"`
}

// IoTSensorDataStream represents a sensor data stream configuration
type IoTSensorDataStream struct {
	database.BaseModel

	ConnectionID   string    `gorm:"not null;index" json:"connection_id"`
	StreamID       string    `gorm:"uniqueIndex;not null" json:"stream_id"`
	SensorType     string    `gorm:"type:varchar(50);not null" json:"sensor_type"` // accelerometer, gyroscope, gps, battery, etc.

	// Stream Configuration
	SampleRate     int       `gorm:"default:1" json:"sample_rate"`      // samples per second
	DataFormat     string    `gorm:"default:'json'" json:"data_format"` // json, binary, protobuf
	Compression    string    `gorm:"default:'none'" json:"compression"` // none, gzip, lz4
	Encryption     string    `gorm:"default:'tls'" json:"encryption"`   // none, tls, aes256

	// Data Quality
	RequiredAccuracy   float64 `json:"required_accuracy,omitempty"`   // minimum accuracy required
	RequiredPrecision  float64 `json:"required_precision,omitempty"`  // minimum precision required
	QualityThreshold   float64 `json:"quality_threshold,omitempty"`   // 0-100, reject data below this

	// Filtering & Processing
	EnableFiltering    bool      `gorm:"default:false" json:"enable_filtering"`
	FilterConfig       string    `gorm:"type:json" json:"filter_config,omitempty"` // JSON configuration
	EnableAggregation  bool      `gorm:"default:false" json:"enable_aggregation"`
	AggregationWindow  int       `gorm:"default:60" json:"aggregation_window"` // seconds
	EnableAnomalyDetection bool  `gorm:"default:false" json:"enable_anomaly_detection"`

	// Status
	Status             string    `gorm:"type:varchar(20);default:'inactive'" json:"status"` // active, inactive, error, paused
	LastDataReceived   *time.Time `json:"last_data_received,omitempty"`
	DataPointsReceived int64     `gorm:"default:0" json:"data_points_received"`
	DataPointsProcessed int64    `gorm:"default:0" json:"data_points_processed"`
	DataPointsRejected int64     `gorm:"default:0" json:"data_points_rejected"`

	// Performance Metrics
	AverageLatency     int       `json:"average_latency,omitempty"` // milliseconds
	Throughput         float64   `json:"throughput,omitempty"`     // data points per second
	ErrorRate          float64   `json:"error_rate,omitempty"`      // percentage

	// Relationships
	Connection         IoTConnection `gorm:"foreignKey:ConnectionID;references:ConnectionID" json:"connection,omitempty"`
	SensorData         []IoTSensorData `gorm:"foreignKey:StreamID;references:StreamID" json:"sensor_data,omitempty"`
}

// IoTSensorData represents individual sensor data points
type IoTSensorData struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	StreamID       string    `gorm:"index" json:"stream_id"`
	SensorID       string    `gorm:"not null;index" json:"sensor_id"`
	SensorType     string    `gorm:"type:varchar(50);not null" json:"sensor_type"`

	// Sensor Reading
	Value          float64   `gorm:"not null" json:"value"`
	Unit           string    `gorm:"type:varchar(20)" json:"unit,omitempty"`
	Timestamp      time.Time `gorm:"not null;index" json:"timestamp"`

	// Location Data (if applicable)
	Latitude       float64   `json:"latitude,omitempty"`
	Longitude      float64   `json:"longitude,omitempty"`
	Altitude       float64   `json:"altitude,omitempty"`
	LocationAccuracy float64 `json:"location_accuracy,omitempty"` // meters

	// Sensor Metadata
	Accuracy       float64   `json:"sensor_accuracy,omitempty"` // sensor accuracy 0-100
	Precision      float64   `json:"precision,omitempty"`      // measurement precision
	RawData        string    `gorm:"type:json" json:"raw_data,omitempty"` // original raw data

	// Quality Assessment
	IsValid        bool      `gorm:"default:true" json:"is_valid"`
	QualityScore   float64   `gorm:"default:1.0" json:"quality_score"` // 0-1, overall quality
	ValidationErrors string  `gorm:"type:json" json:"validation_errors,omitempty"` // JSON array of errors

	// Processing Status
	IsProcessed    bool      `gorm:"default:false;index" json:"is_processed"`
	ProcessedAt    *time.Time `json:"processed_at,omitempty"`
	ProcessingErrors string  `gorm:"type:json" json:"processing_errors,omitempty"`

	// Anomaly Detection
	IsAnomaly      bool      `gorm:"default:false" json:"is_anomaly"`
	AnomalyScore   float64   `json:"anomaly_score,omitempty"`
	AnomalyType    string    `json:"anomaly_type,omitempty"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	DataStream     IoTSensorDataStream `gorm:"foreignKey:StreamID;references:StreamID" json:"data_stream,omitempty"`
}

// IoTCommand represents a command sent to an IoT device
type IoTCommand struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ConnectionID   string    `gorm:"index" json:"connection_id"`
	CommandID      string    `gorm:"uniqueIndex;not null" json:"command_id"`

	// Command Details
	CommandType    string    `gorm:"type:varchar(50);not null" json:"command_type"` // ping, update, reboot, diagnostics, config_update
	CommandData    string    `gorm:"type:json" json:"command_data"` // JSON command payload

	// Execution Details
	Status         string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, sent, delivered, executed, failed, cancelled
	Priority       string    `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, urgent

	// Timing
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	ExecutedAt     *time.Time `json:"executed_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	ExpiredAt      *time.Time `json:"expired_at,omitempty"`

	// Response & Results
	ResponseData   string    `gorm:"type:json" json:"response_data,omitempty"`
	ExecutionTime  int       `json:"execution_time,omitempty"` // milliseconds
	ExitCode       int       `json:"exit_code,omitempty"`

	// Error Handling
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`
	ErrorCode      string    `json:"error_code,omitempty"`
	RetryCount     int       `gorm:"default:0" json:"retry_count"`
	MaxRetries     int       `gorm:"default:3" json:"max_retries"`

	// Audit
	InitiatedBy    uuid.UUID `gorm:"type:uuid;not null" json:"initiated_by"`
	InitiatedByName string   `json:"initiated_by_name,omitempty"`

	// Timeout & Expiration
	TimeoutSeconds int       `gorm:"default:30" json:"timeout_seconds"`
	ExpiresAt      time.Time `gorm:"not null" json:"expires_at"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	Connection     IoTConnection `gorm:"foreignKey:ConnectionID;references:ConnectionID" json:"connection,omitempty"`
}

// IoTCommandTemplate represents reusable command templates
type IoTCommandTemplate struct {
	database.BaseModel

	TemplateID     string    `gorm:"uniqueIndex;not null" json:"template_id"`
	Name           string    `gorm:"not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	CommandType    string    `gorm:"type:varchar(50);not null" json:"command_type"`

	// Template Configuration
	TemplateData   string    `gorm:"type:json" json:"template_data"` // JSON template with placeholders
	Parameters     string    `gorm:"type:json" json:"parameters"`    // JSON parameter definitions
	DefaultValues  string    `gorm:"type:json" json:"default_values"` // JSON default parameter values

	// Validation
	ParameterValidation string `gorm:"type:json" json:"parameter_validation"` // JSON validation rules
	RequiresApproval    bool   `gorm:"default:false" json:"requires_approval"`

	// Usage Tracking
	UsageCount      int       `gorm:"default:0" json:"usage_count"`
	LastUsed        *time.Time `json:"last_used,omitempty"`
	SuccessRate     float64   `json:"success_rate,omitempty"`

	// Categories & Tags
	Category        string    `json:"category,omitempty"`
	Tags            []string  `gorm:"type:json" json:"tags,omitempty"`

	// Audit
	CreatedBy       uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	IsPublic        bool      `gorm:"default:false" json:"is_public"`
}

// IoTDeviceConfiguration represents device-specific configuration
type IoTDeviceConfiguration struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Communication Settings
	HeartbeatInterval    int    `gorm:"default:30" json:"heartbeat_interval"`     // seconds
	DataReportingInterval int   `gorm:"default:60" json:"data_reporting_interval"` // seconds
	PreferredProtocol    string `gorm:"default:'mqtt'" json:"preferred_protocol"`
	FallbackProtocol     string `json:"fallback_protocol,omitempty"`

	// Data Collection Settings
	EnabledSensors       []string `gorm:"type:json" json:"enabled_sensors"`       // enabled sensor types
	SensorConfigurations string   `gorm:"type:json" json:"sensor_configurations"` // JSON sensor configs
	DataFilters          string   `gorm:"type:json" json:"data_filters"`          // JSON filtering rules

	// Power Management
	LowPowerMode         bool     `gorm:"default:false" json:"low_power_mode"`
	SleepInterval        int      `json:"sleep_interval,omitempty"` // seconds
	WakeOnMotion         bool     `gorm:"default:false" json:"wake_on_motion"`

	// Security Settings
	RequireEncryption    bool     `gorm:"default:true" json:"require_encryption"`
	AllowedIPs           []string `gorm:"type:json" json:"allowed_ips,omitempty"`
	CertificateRotation  int      `json:"certificate_rotation,omitempty"` // days

	// Alerting & Monitoring
	EnableRemoteAlerts   bool     `gorm:"default:true" json:"enable_remote_alerts"`
	AlertThresholds      string   `gorm:"type:json" json:"alert_thresholds"` // JSON threshold configs
	MonitoringRules      string   `gorm:"type:json" json:"monitoring_rules"` // JSON monitoring rules

	// Firmware & Updates
	AutoUpdateEnabled    bool     `gorm:"default:true" json:"auto_update_enabled"`
	UpdateWindow         string   `json:"update_window,omitempty"` // cron expression
	UpdatePriority       string   `gorm:"default:'normal'" json:"update_priority"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// TableName returns the table name for IoTConnection
func (IoTConnection) TableName() string {
	return "iot_connections"
}

// TableName returns the table name for IoTSensorDataStream
func (IoTSensorDataStream) TableName() string {
	return "iot_sensor_data_streams"
}

// TableName returns the table name for IoTSensorData
func (IoTSensorData) TableName() string {
	return "iot_sensor_data"
}

// TableName returns the table name for IoTCommand
func (IoTCommand) TableName() string {
	return "iot_commands"
}

// TableName returns the table name for IoTCommandTemplate
func (IoTCommandTemplate) TableName() string {
	return "iot_command_templates"
}

// TableName returns the table name for IoTDeviceConfiguration
func (IoTDeviceConfiguration) TableName() string {
	return "iot_device_configurations"
}
