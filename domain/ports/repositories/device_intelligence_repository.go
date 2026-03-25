package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// DeviceIntelligenceRepository defines the interface for device intelligence operations
type DeviceIntelligenceRepository interface {
	// IoT Connection operations
	CreateIoTConnection(ctx context.Context, connection *models.IoTConnection) error
	GetIoTConnection(ctx context.Context, connectionID string) (*models.IoTConnection, error)
	GetDeviceIoTConnection(ctx context.Context, deviceID uuid.UUID) (*models.IoTConnection, error)
	UpdateIoTConnection(ctx context.Context, connection *models.IoTConnection) error
	DeleteIoTConnection(ctx context.Context, connectionID string) error
	GetActiveIoTConnections(ctx context.Context) ([]*models.IoTConnection, error)

	// IoT Sensor Data operations
	CreateIoTSensorData(ctx context.Context, data *models.IoTSensorData) error
	GetIoTSensorData(ctx context.Context, deviceID uuid.UUID, sensorID string, limit int) ([]*models.IoTSensorData, error)
	GetIoTSensorDataByTimeRange(ctx context.Context, deviceID uuid.UUID, sensorType string, start, end time.Time) ([]*models.IoTSensorData, error)
	GetLatestIoTSensorData(ctx context.Context, deviceID uuid.UUID) ([]*models.IoTSensorData, error)

	// Predictive Maintenance operations
	CreatePredictiveAlert(ctx context.Context, alert *models.PredictiveMaintenanceAlert) error
	GetPredictiveAlerts(ctx context.Context, deviceID uuid.UUID, status string, limit int) ([]*models.PredictiveMaintenanceAlert, error)
	GetActivePredictiveAlerts(ctx context.Context, deviceID uuid.UUID) ([]*models.PredictiveMaintenanceAlert, error)
	UpdatePredictiveAlert(ctx context.Context, alert *models.PredictiveMaintenanceAlert) error
	ResolvePredictiveAlert(ctx context.Context, alertID string, resolvedBy uuid.UUID, notes string) error

	// Device Health operations
	CreateOrUpdateDeviceHealthScore(ctx context.Context, score *models.DeviceHealthScore) error
	GetDeviceHealthScore(ctx context.Context, deviceID uuid.UUID) (*models.DeviceHealthScore, error)
	GetDevicesByHealthRisk(ctx context.Context, riskLevel string, limit int) ([]*models.DeviceHealthScore, error)

	// Usage Pattern operations
	CreateUsagePattern(ctx context.Context, pattern *models.DeviceUsagePattern) error
	GetUsagePatterns(ctx context.Context, deviceID uuid.UUID, patternType string, limit int) ([]*models.DeviceUsagePattern, error)
	GetLatestUsagePattern(ctx context.Context, deviceID uuid.UUID, patternType string) (*models.DeviceUsagePattern, error)
	UpdateUsagePattern(ctx context.Context, pattern *models.DeviceUsagePattern) error

	// IoT Command operations
	CreateIoTCommand(ctx context.Context, command *models.IoTCommand) error
	GetIoTCommand(ctx context.Context, commandID string) (*models.IoTCommand, error)
	GetDeviceIoTCommands(ctx context.Context, deviceID uuid.UUID, status string, limit int) ([]*models.IoTCommand, error)
	UpdateIoTCommand(ctx context.Context, command *models.IoTCommand) error
}

// deviceIntelligenceRepository implements DeviceIntelligenceRepository
type deviceIntelligenceRepository struct {
	db *gorm.DB
}

// NewDeviceIntelligenceRepository creates a new device intelligence repository
func NewDeviceIntelligenceRepository(db *gorm.DB) DeviceIntelligenceRepository {
	return &deviceIntelligenceRepository{db: db}
}

// CreateIoTConnection creates a new IoT connection
func (r *deviceIntelligenceRepository) CreateIoTConnection(ctx context.Context, connection *models.IoTConnection) error {
	return r.db.WithContext(ctx).Create(connection).Error
}

// GetIoTConnection gets an IoT connection by ID
func (r *deviceIntelligenceRepository) GetIoTConnection(ctx context.Context, connectionID string) (*models.IoTConnection, error) {
	var connection models.IoTConnection
	err := r.db.WithContext(ctx).Where("connection_id = ?", connectionID).First(&connection).Error
	return &connection, err
}

// GetDeviceIoTConnection gets the IoT connection for a device
func (r *deviceIntelligenceRepository) GetDeviceIoTConnection(ctx context.Context, deviceID uuid.UUID) (*models.IoTConnection, error) {
	var connection models.IoTConnection
	err := r.db.WithContext(ctx).Where("device_id = ? AND status = 'connected'", deviceID).First(&connection).Error
	return &connection, err
}

// UpdateIoTConnection updates an IoT connection
func (r *deviceIntelligenceRepository) UpdateIoTConnection(ctx context.Context, connection *models.IoTConnection) error {
	return r.db.WithContext(ctx).Save(connection).Error
}

// DeleteIoTConnection deletes an IoT connection
func (r *deviceIntelligenceRepository) DeleteIoTConnection(ctx context.Context, connectionID string) error {
	return r.db.WithContext(ctx).Delete(&models.IoTConnection{}, "connection_id = ?", connectionID).Error
}

// GetActiveIoTConnections gets all active IoT connections
func (r *deviceIntelligenceRepository) GetActiveIoTConnections(ctx context.Context) ([]*models.IoTConnection, error) {
	var connections []*models.IoTConnection
	err := r.db.WithContext(ctx).Where("status = 'connected'").Find(&connections).Error
	return connections, err
}

// CreateIoTSensorData creates IoT sensor data
func (r *deviceIntelligenceRepository) CreateIoTSensorData(ctx context.Context, data *models.IoTSensorData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// GetIoTSensorData gets sensor data for a device and sensor
func (r *deviceIntelligenceRepository) GetIoTSensorData(ctx context.Context, deviceID uuid.UUID, sensorID string, limit int) ([]*models.IoTSensorData, error) {
	var data []*models.IoTSensorData
	query := r.db.WithContext(ctx).Where("device_id = ? AND sensor_id = ?", deviceID, sensorID).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&data).Error
	return data, err
}

// GetIoTSensorDataByTimeRange gets sensor data within a time range
func (r *deviceIntelligenceRepository) GetIoTSensorDataByTimeRange(ctx context.Context, deviceID uuid.UUID, sensorType string, start, end time.Time) ([]*models.IoTSensorData, error) {
	var data []*models.IoTSensorData
	err := r.db.WithContext(ctx).Where("device_id = ? AND sensor_type = ? AND timestamp BETWEEN ? AND ?",
		deviceID, sensorType, start, end).Order("timestamp ASC").Find(&data).Error
	return data, err
}

// GetLatestIoTSensorData gets the latest sensor data for all sensors on a device
func (r *deviceIntelligenceRepository) GetLatestIoTSensorData(ctx context.Context, deviceID uuid.UUID) ([]*models.IoTSensorData, error) {
	var data []*models.IoTSensorData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Where("timestamp = (SELECT MAX(timestamp) FROM iot_sensor_data WHERE device_id = ? AND sensor_type = iot_sensor_data.sensor_type)", deviceID).
		Find(&data).Error
	return data, err
}

// CreatePredictiveAlert creates a predictive maintenance alert
func (r *deviceIntelligenceRepository) CreatePredictiveAlert(ctx context.Context, alert *models.PredictiveMaintenanceAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// GetPredictiveAlerts gets predictive alerts for a device
func (r *deviceIntelligenceRepository) GetPredictiveAlerts(ctx context.Context, deviceID uuid.UUID, status string, limit int) ([]*models.PredictiveMaintenanceAlert, error) {
	var alerts []*models.PredictiveMaintenanceAlert
	query := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("created_at DESC")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&alerts).Error
	return alerts, err
}

// GetActivePredictiveAlerts gets active predictive alerts for a device
func (r *deviceIntelligenceRepository) GetActivePredictiveAlerts(ctx context.Context, deviceID uuid.UUID) ([]*models.PredictiveMaintenanceAlert, error) {
	var alerts []*models.PredictiveMaintenanceAlert
	err := r.db.WithContext(ctx).Where("device_id = ? AND status = 'active'", deviceID).
		Order("severity DESC, created_at DESC").Find(&alerts).Error
	return alerts, err
}

// UpdatePredictiveAlert updates a predictive alert
func (r *deviceIntelligenceRepository) UpdatePredictiveAlert(ctx context.Context, alert *models.PredictiveMaintenanceAlert) error {
	return r.db.WithContext(ctx).Save(alert).Error
}

// ResolvePredictiveAlert resolves a predictive alert
func (r *deviceIntelligenceRepository) ResolvePredictiveAlert(ctx context.Context, alertID string, resolvedBy uuid.UUID, notes string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.PredictiveMaintenanceAlert{}).
		Where("alert_id = ?", alertID).
		Updates(map[string]interface{}{
			"status":           "resolved",
			"resolved_at":      now,
			"resolved_by":      resolvedBy,
			"resolution_notes": notes,
		}).Error
}

// CreateOrUpdateDeviceHealthScore creates or updates device health score
func (r *deviceIntelligenceRepository) CreateOrUpdateDeviceHealthScore(ctx context.Context, score *models.DeviceHealthScore) error {
	return r.db.WithContext(ctx).Save(score).Error
}

// GetDeviceHealthScore gets device health score
func (r *deviceIntelligenceRepository) GetDeviceHealthScore(ctx context.Context, deviceID uuid.UUID) (*models.DeviceHealthScore, error) {
	var score models.DeviceHealthScore
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&score).Error
	return &score, err
}

// GetDevicesByHealthRisk gets devices by health risk level
func (r *deviceIntelligenceRepository) GetDevicesByHealthRisk(ctx context.Context, riskLevel string, limit int) ([]*models.DeviceHealthScore, error) {
	var scores []*models.DeviceHealthScore
	query := r.db.WithContext(ctx).Where("risk_level = ?", riskLevel).Order("overall_score ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&scores).Error
	return scores, err
}

// CreateUsagePattern creates a usage pattern
func (r *deviceIntelligenceRepository) CreateUsagePattern(ctx context.Context, pattern *models.DeviceUsagePattern) error {
	return r.db.WithContext(ctx).Create(pattern).Error
}

// GetUsagePatterns gets usage patterns for a device
func (r *deviceIntelligenceRepository) GetUsagePatterns(ctx context.Context, deviceID uuid.UUID, patternType string, limit int) ([]*models.DeviceUsagePattern, error) {
	var patterns []*models.DeviceUsagePattern
	query := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("created_at DESC")
	if patternType != "" {
		query = query.Where("pattern_type = ?", patternType)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&patterns).Error
	return patterns, err
}

// GetLatestUsagePattern gets the latest usage pattern for a device and type
func (r *deviceIntelligenceRepository) GetLatestUsagePattern(ctx context.Context, deviceID uuid.UUID, patternType string) (*models.DeviceUsagePattern, error) {
	var pattern models.DeviceUsagePattern
	err := r.db.WithContext(ctx).Where("device_id = ? AND pattern_type = ?", deviceID, patternType).
		Order("created_at DESC").First(&pattern).Error
	return &pattern, err
}

// UpdateUsagePattern updates a usage pattern
func (r *deviceIntelligenceRepository) UpdateUsagePattern(ctx context.Context, pattern *models.DeviceUsagePattern) error {
	return r.db.WithContext(ctx).Save(pattern).Error
}

// CreateIoTCommand creates an IoT command
func (r *deviceIntelligenceRepository) CreateIoTCommand(ctx context.Context, command *models.IoTCommand) error {
	return r.db.WithContext(ctx).Create(command).Error
}

// GetIoTCommand gets an IoT command by ID
func (r *deviceIntelligenceRepository) GetIoTCommand(ctx context.Context, commandID string) (*models.IoTCommand, error) {
	var command models.IoTCommand
	err := r.db.WithContext(ctx).Where("command_id = ?", commandID).First(&command).Error
	return &command, err
}

// GetDeviceIoTCommands gets IoT commands for a device
func (r *deviceIntelligenceRepository) GetDeviceIoTCommands(ctx context.Context, deviceID uuid.UUID, status string, limit int) ([]*models.IoTCommand, error) {
	var commands []*models.IoTCommand
	query := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("created_at DESC")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&commands).Error
	return commands, err
}

// UpdateIoTCommand updates an IoT command
func (r *deviceIntelligenceRepository) UpdateIoTCommand(ctx context.Context, command *models.IoTCommand) error {
	return r.db.WithContext(ctx).Save(command).Error
}
