package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/pkg/notification"
)

// DeviceIntelligenceService handles device intelligence operations
type DeviceIntelligenceService struct {
	db              *gorm.DB
	deviceIntelRepo repositories.DeviceIntelligenceRepository
	notificationSvc *notification.Service
	logger          *logrus.Logger
}

// NewDeviceIntelligenceService creates a new device intelligence service
func NewDeviceIntelligenceService(db *gorm.DB, notificationSvc *notification.Service, logger *logrus.Logger) *DeviceIntelligenceService {
	return &DeviceIntelligenceService{
		db:              db,
		deviceIntelRepo: repositories.NewDeviceIntelligenceRepository(db),
		notificationSvc: notificationSvc,
		logger:          logger,
	}
}

// IoTConnectRequest represents IoT connection request
type IoTConnectRequest struct {
	DeviceID               uuid.UUID `json:"device_id" binding:"required"`
	Protocol               string    `json:"protocol" binding:"required"` // mqtt, websocket, http
	IPAddress              string    `json:"ip_address" binding:"required"`
	MACAddress             string    `json:"mac_address,omitempty"`
	FirmwareVersion        string    `json:"firmware_version,omitempty"`
	UserAgent              string    `json:"user_agent,omitempty"`
	Location               string    `json:"location,omitempty"`
	NetworkType            string    `json:"network_type,omitempty"`
	CertificateFingerprint string    `json:"certificate_fingerprint,omitempty"`
}

// IoTConnectResponse represents IoT connection response
type IoTConnectResponse struct {
	ConnectionID string `json:"connection_id"`
	Status       string `json:"status"`
}

// ConnectIoTDevice connects an IoT device
func (s *DeviceIntelligenceService) ConnectIoTDevice(ctx context.Context, req *IoTConnectRequest) (*IoTConnectResponse, error) {
	// Check if device already has an active connection
	existing, err := s.deviceIntelRepo.GetDeviceIoTConnection(ctx, req.DeviceID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing connection: %w", err)
	}

	if existing != nil && existing.Status == "connected" {
		// Update existing connection
		existing.LastHeartbeat = time.Now()
		existing.LastSeen = time.Now()
		if err := s.deviceIntelRepo.UpdateIoTConnection(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update existing connection: %w", err)
		}
		return &IoTConnectResponse{
			ConnectionID: existing.ConnectionID,
			Status:       "reconnected",
		}, nil
	}

	// Create new connection
	connectionID := s.generateConnectionID()
	connection := &models.IoTConnection{
		DeviceID:               req.DeviceID,
		ConnectionID:           connectionID,
		Protocol:               req.Protocol,
		Status:                 "connected",
		IPAddress:              req.IPAddress,
		MACAddress:             req.MACAddress,
		FirmwareVersion:        req.FirmwareVersion,
		LastHeartbeat:          time.Now(),
		LastSeen:               time.Now(),
		ConnectedAt:            time.Now(),
		UserAgent:              req.UserAgent,
		Location:               req.Location,
		NetworkType:            req.NetworkType,
		CertificateFingerprint: req.CertificateFingerprint,
		EncryptionEnabled:      true,
	}

	if err := s.deviceIntelRepo.CreateIoTConnection(ctx, connection); err != nil {
		return nil, fmt.Errorf("failed to create IoT connection: %w", err)
	}

	return &IoTConnectResponse{
		ConnectionID: connectionID,
		Status:       "connected",
	}, nil
}

// DisconnectIoTDevice disconnects an IoT device
func (s *DeviceIntelligenceService) DisconnectIoTDevice(ctx context.Context, connectionID string) error {
	connection, err := s.deviceIntelRepo.GetIoTConnection(ctx, connectionID)
	if err != nil {
		return fmt.Errorf("failed to get IoT connection: %w", err)
	}

	now := time.Now()
	connection.Status = "disconnected"
	connection.DisconnectedAt = &now

	return s.deviceIntelRepo.UpdateIoTConnection(ctx, connection)
}

// SendIoTCommand sends a command to an IoT device
func (s *DeviceIntelligenceService) SendIoTCommand(ctx context.Context, deviceID uuid.UUID, commandType string, commandData map[string]interface{}, priority string, initiatedBy uuid.UUID) (*models.IoTCommand, error) {
	commandID := s.generateCommandID()
	command := &models.IoTCommand{
		DeviceID:    deviceID,
		CommandID:   commandID,
		CommandType: commandType,
		Status:      "pending",
		Priority:    priority,
		CommandData: fmt.Sprintf("%v", commandData), // Convert to JSON string
		InitiatedBy: initiatedBy,
		MaxRetries:  3,
	}

	if err := s.deviceIntelRepo.CreateIoTCommand(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to create IoT command: %w", err)
	}

	// TODO: Send command to IoT device via MQTT/WebSocket/etc.

	return command, nil
}

// RecordSensorData records sensor data from IoT device
func (s *DeviceIntelligenceService) RecordSensorData(ctx context.Context, deviceID uuid.UUID, sensorID, sensorType string, value float64, unit string, accuracy, precision float64) error {
	sensorData := &models.IoTSensorData{
		DeviceID:     deviceID,
		SensorID:     sensorID,
		SensorType:   sensorType,
		Value:        value,
		Unit:         unit,
		Timestamp:    time.Now(),
		Accuracy:     accuracy,
		Precision:    precision,
		IsValid:      s.validateSensorData(sensorType, value),
		QualityScore: s.calculateDataQuality(accuracy, precision),
	}

	if err := s.deviceIntelRepo.CreateIoTSensorData(ctx, sensorData); err != nil {
		return fmt.Errorf("failed to record sensor data: %w", err)
	}

	// Check for anomalies and predictive maintenance triggers
	go s.analyzeSensorData(sensorData)

	return nil
}

// GetDeviceHealthScore gets the health score for a device
func (s *DeviceIntelligenceService) GetDeviceHealthScore(ctx context.Context, deviceID uuid.UUID) (*models.DeviceHealthScore, error) {
	score, err := s.deviceIntelRepo.GetDeviceHealthScore(ctx, deviceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Calculate initial health score
			return s.calculateDeviceHealthScore(ctx, deviceID)
		}
		return nil, fmt.Errorf("failed to get device health score: %w", err)
	}

	// Check if score needs refresh (older than 24 hours)
	if time.Since(score.LastCalculated) > 24*time.Hour {
		return s.calculateDeviceHealthScore(ctx, deviceID)
	}

	return score, nil
}

// CreatePredictiveAlert creates a predictive maintenance alert
func (s *DeviceIntelligenceService) CreatePredictiveAlert(ctx context.Context, deviceID uuid.UUID, alertType, severity, title, description string, confidenceScore float64, triggeringData map[string]interface{}) (*models.PredictiveMaintenanceAlert, error) {
	alertID := s.generateAlertID()
	alert := &models.PredictiveMaintenanceAlert{
		DeviceID:         deviceID,
		AlertID:          alertID,
		AlertType:        alertType,
		Severity:         severity,
		Status:           "active",
		Title:            title,
		Description:      description,
		ConfidenceScore:  confidenceScore,
		PredictionModel:  "ai-v1",
		AlgorithmVersion: "1.0.0",
	}

	if triggeringData != nil {
		alert.TriggeringSensors = fmt.Sprintf("%v", triggeringData)
		alert.TriggeringValues = fmt.Sprintf("%v", triggeringData)
	}

	if err := s.deviceIntelRepo.CreatePredictiveAlert(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to create predictive alert: %w", err)
	}

	// Send notification
	s.sendPredictiveAlertNotification(alert)

	return alert, nil
}

// ResolvePredictiveAlert resolves a predictive alert
func (s *DeviceIntelligenceService) ResolvePredictiveAlert(ctx context.Context, alertID string, resolvedBy uuid.UUID, resolutionNotes string) error {
	return s.deviceIntelRepo.ResolvePredictiveAlert(ctx, alertID, resolvedBy, resolutionNotes)
}

// GetActivePredictiveAlerts gets active alerts for a device
func (s *DeviceIntelligenceService) GetActivePredictiveAlerts(ctx context.Context, deviceID uuid.UUID) ([]*models.PredictiveMaintenanceAlert, error) {
	return s.deviceIntelRepo.GetActivePredictiveAlerts(ctx, deviceID)
}

// GetPredictiveAlerts gets alerts for a device with filters
func (s *DeviceIntelligenceService) GetPredictiveAlerts(ctx context.Context, deviceID uuid.UUID, status string, limit int) ([]*models.PredictiveMaintenanceAlert, error) {
	return s.deviceIntelRepo.GetPredictiveAlerts(ctx, deviceID, status, limit)
}

// AnalyzeUsagePatterns analyzes usage patterns for predictive maintenance
func (s *DeviceIntelligenceService) AnalyzeUsagePatterns(ctx context.Context, deviceID uuid.UUID) error {
	// Get recent sensor data for analysis
	sensorData, err := s.deviceIntelRepo.GetLatestIoTSensorData(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get sensor data for analysis: %w", err)
	}

	// Analyze patterns and create usage pattern record
	pattern := s.calculateUsagePattern(sensorData)
	if err := s.deviceIntelRepo.CreateUsagePattern(ctx, pattern); err != nil {
		return fmt.Errorf("failed to create usage pattern: %w", err)
	}

	return nil
}

// GetDeviceAnalytics gets comprehensive analytics for a device
func (s *DeviceIntelligenceService) GetDeviceAnalytics(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Get health score
	healthScore, err := s.GetDeviceHealthScore(ctx, deviceID)
	if err != nil {
		s.logger.Warnf("Failed to get health score for analytics: %v", err)
	} else {
		analytics["health_score"] = healthScore
	}

	// Get active alerts
	alerts, err := s.GetActivePredictiveAlerts(ctx, deviceID)
	if err != nil {
		s.logger.Warnf("Failed to get active alerts for analytics: %v", err)
	} else {
		analytics["active_alerts"] = alerts
	}

	// Get recent sensor data summary
	sensorData, err := s.deviceIntelRepo.GetLatestIoTSensorData(ctx, deviceID)
	if err != nil {
		s.logger.Warnf("Failed to get sensor data for analytics: %v", err)
	} else {
		analytics["sensor_summary"] = s.summarizeSensorData(sensorData)
	}

	// Get IoT connection status
	connection, err := s.deviceIntelRepo.GetDeviceIoTConnection(ctx, deviceID)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Warnf("Failed to get IoT connection for analytics: %v", err)
	} else {
		analytics["iot_connection"] = connection
	}

	return analytics, nil
}

// Private helper methods

func (s *DeviceIntelligenceService) generateConnectionID() string {
	return fmt.Sprintf("iot-%s", uuid.New().String()[:12])
}

func (s *DeviceIntelligenceService) generateCommandID() string {
	return fmt.Sprintf("cmd-%s", uuid.New().String()[:12])
}

func (s *DeviceIntelligenceService) generateAlertID() string {
	return fmt.Sprintf("alert-%s", uuid.New().String()[:12])
}

func (s *DeviceIntelligenceService) validateSensorData(sensorType string, value float64) bool {
	// Basic validation based on sensor type
	switch sensorType {
	case "battery_level":
		return value >= 0 && value <= 100
	case "temperature":
		return value >= -50 && value <= 100
	case "accelerometer_x", "accelerometer_y", "accelerometer_z":
		return value >= -20 && value <= 20
	case "gyroscope_x", "gyroscope_y", "gyroscope_z":
		return value >= -10 && value <= 10
	default:
		return true // Allow unknown sensor types
	}
}

func (s *DeviceIntelligenceService) calculateDataQuality(accuracy, precision float64) float64 {
	if accuracy == 0 || precision == 0 {
		return 0.5 // Default quality
	}
	return math.Min(accuracy/100.0, precision/100.0)
}

func (s *DeviceIntelligenceService) analyzeSensorData(data *models.IoTSensorData) {
	ctx := context.Background()

	// Simple anomaly detection (placeholder for ML model)
	if s.detectAnomaly(data) {
		alert, err := s.CreatePredictiveAlert(ctx, data.DeviceID, "anomaly", "medium",
			"Sensor Anomaly Detected", fmt.Sprintf("Unusual %s reading detected", data.SensorType),
			0.75, map[string]interface{}{
				"sensor_type": data.SensorType,
				"value":       data.Value,
				"timestamp":   data.Timestamp,
			})
		if err != nil {
			s.logger.Errorf("Failed to create anomaly alert: %v", err)
		} else {
			s.logger.Infof("Created anomaly alert: %s", alert.AlertID)
		}
	}
}

func (s *DeviceIntelligenceService) detectAnomaly(data *models.IoTSensorData) bool {
	// Simple threshold-based anomaly detection (placeholder)
	switch data.SensorType {
	case "battery_level":
		return data.Value < 10 // Low battery alert
	case "temperature":
		return data.Value > 80 || data.Value < -10 // Extreme temperature
	default:
		return false
	}
}

func (s *DeviceIntelligenceService) calculateDeviceHealthScore(ctx context.Context, deviceID uuid.UUID) (*models.DeviceHealthScore, error) {
	// Get recent sensor data for health calculation
	sensorData, err := s.deviceIntelRepo.GetLatestIoTSensorData(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor data for health calculation: %w", err)
	}

	// Calculate component scores (simplified)
	batteryHealth := s.calculateBatteryHealth(sensorData)
	screenHealth := s.calculateScreenHealth(sensorData)
	hardwareHealth := s.calculateHardwareHealth(sensorData)
	softwareHealth := s.calculateSoftwareHealth(sensorData)
	usageHealth := s.calculateUsageHealth(sensorData)

	// Calculate overall score (weighted average)
	overallScore := (batteryHealth*0.3 + screenHealth*0.2 + hardwareHealth*0.25 + softwareHealth*0.15 + usageHealth*0.1)

	// Determine risk level
	riskLevel := s.determineRiskLevel(overallScore)

	score := &models.DeviceHealthScore{
		DeviceID:       deviceID,
		OverallScore:   overallScore,
		LastCalculated: time.Now(),
		BatteryHealth:  batteryHealth,
		ScreenHealth:   screenHealth,
		HardwareHealth: hardwareHealth,
		SoftwareHealth: softwareHealth,
		UsageHealth:    usageHealth,
		RiskLevel:      riskLevel,
	}

	if err := s.deviceIntelRepo.CreateOrUpdateDeviceHealthScore(ctx, score); err != nil {
		return nil, fmt.Errorf("failed to save health score: %w", err)
	}

	return score, nil
}

func (s *DeviceIntelligenceService) calculateBatteryHealth(sensorData []*models.IoTSensorData) float64 {
	for _, data := range sensorData {
		if data.SensorType == "battery_level" {
			// Battery health decreases as level drops
			return math.Max(0, data.Value)
		}
	}
	return 50 // Default if no battery data
}

func (s *DeviceIntelligenceService) calculateScreenHealth(sensorData []*models.IoTSensorData) float64 {
	// Placeholder - would analyze touch/accelerometer data for screen issues
	return 85
}

func (s *DeviceIntelligenceService) calculateHardwareHealth(sensorData []*models.IoTSensorData) float64 {
	// Placeholder - would analyze accelerometer/gyroscope data for hardware wear
	return 78
}

func (s *DeviceIntelligenceService) calculateSoftwareHealth(sensorData []*models.IoTSensorData) float64 {
	// Placeholder - would analyze crash reports, app usage patterns
	return 92
}

func (s *DeviceIntelligenceService) calculateUsageHealth(sensorData []*models.IoTSensorData) float64 {
	// Placeholder - would analyze usage patterns for optimal health
	return 88
}

func (s *DeviceIntelligenceService) determineRiskLevel(score float64) string {
	switch {
	case score >= 80:
		return "low"
	case score >= 60:
		return "medium"
	case score >= 40:
		return "high"
	default:
		return "critical"
	}
}

func (s *DeviceIntelligenceService) calculateUsagePattern(sensorData []*models.IoTSensorData) *models.DeviceUsagePattern {
	// Simplified usage pattern calculation
	values := make([]float64, len(sensorData))
	for i, data := range sensorData {
		values[i] = data.Value
	}

	mean, stdDev := s.calculateStats(values)

	return &models.DeviceUsagePattern{
		DeviceID:          sensorData[0].DeviceID,
		PatternID:         s.generatePatternID(),
		PatternType:       "daily",
		TimeRange:         "day",
		StartDate:         time.Now().Add(-24 * time.Hour),
		EndDate:           time.Now(),
		AverageUsageHours: mean / 24, // Placeholder calculation
		AnomaliesDetected: stdDev > mean*0.5,
		Mean:              mean,
		StdDev:            stdDev,
		MinValue:          s.minValue(values),
		MaxValue:          s.maxValue(values),
	}
}

func (s *DeviceIntelligenceService) calculateStats(values []float64) (mean, stdDev float64) {
	if len(values) == 0 {
		return 0, 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))

	sumSq := 0.0
	for _, v := range values {
		sumSq += (v - mean) * (v - mean)
	}
	stdDev = math.Sqrt(sumSq / float64(len(values)))

	return mean, stdDev
}

func (s *DeviceIntelligenceService) minValue(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func (s *DeviceIntelligenceService) maxValue(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func (s *DeviceIntelligenceService) generatePatternID() string {
	return fmt.Sprintf("pattern-%s", uuid.New().String()[:12])
}

func (s *DeviceIntelligenceService) summarizeSensorData(sensorData []*models.IoTSensorData) map[string]interface{} {
	summary := make(map[string]interface{})
	sensorTypes := make(map[string][]float64)

	for _, data := range sensorData {
		sensorTypes[data.SensorType] = append(sensorTypes[data.SensorType], data.Value)
	}

	for sensorType, values := range sensorTypes {
		mean, stdDev := s.calculateStats(values)
		summary[sensorType] = map[string]interface{}{
			"latest_value": values[len(values)-1],
			"mean":         mean,
			"std_dev":      stdDev,
			"count":        len(values),
			"min":          s.minValue(values),
			"max":          s.maxValue(values),
		}
	}

	return summary
}

func (s *DeviceIntelligenceService) sendPredictiveAlertNotification(alert *models.PredictiveMaintenanceAlert) {
	// Send notification about predictive alert
	message := fmt.Sprintf("Predictive Alert: %s - %s", alert.Title, alert.Description)
	// TODO: Send to device owner via push notification/email
	s.logger.Infof("Predictive alert notification: %s", message)
}
