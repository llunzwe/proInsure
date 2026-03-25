package categories

import (
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models/device/categories/base"
)

// DeviceModelBridge bridges the hybrid category system with existing device models
// This allows the hybrid system to leverage all existing device sub-entity models
// like DeviceNetworkProfile, DeviceDataUsage, DeviceWiFiProfile, etc.
type DeviceModelBridge struct {
	// Core device identification (from existing device model)
	DeviceID     uuid.UUID `json:"device_id"`
	IMEI         string    `json:"imei"`
	SerialNumber string    `json:"serial_number"`

	// Category-specific spec (hybrid system)
	CategorySpec base.CategorySpec `json:"-"`

	// References to existing device sub-entities
	// These would be loaded from the service layer to avoid circular imports
	NetworkProfileID   *uuid.UUID `json:"network_profile_id,omitempty"`
	DataUsageID        *uuid.UUID `json:"data_usage_id,omitempty"`
	WiFiProfileID      *uuid.UUID `json:"wifi_profile_id,omitempty"`
	BluetoothProfileID *uuid.UUID `json:"bluetooth_profile_id,omitempty"`
	SecurityProfileID  *uuid.UUID `json:"security_profile_id,omitempty"`
	PerformanceDataID  *uuid.UUID `json:"performance_data_id,omitempty"`
	QualityMetricsID   *uuid.UUID `json:"quality_metrics_id,omitempty"`
	SustainabilityID   *uuid.UUID `json:"sustainability_id,omitempty"`
	EmergencyDataID    *uuid.UUID `json:"emergency_data_id,omitempty"`
	ComplianceStatusID *uuid.UUID `json:"compliance_status_id,omitempty"`
}

// DeviceModelAdapter adapts existing device models for use with the category system
type DeviceModelAdapter struct {
	bridge *DeviceModelBridge
}

// NewDeviceModelAdapter creates a new adapter
func NewDeviceModelAdapter(deviceID uuid.UUID, category base.CategoryType) *DeviceModelAdapter {
	return &DeviceModelAdapter{
		bridge: &DeviceModelBridge{
			DeviceID: deviceID,
		},
	}
}

// AttachCategorySpec attaches a category specification to the device
func (a *DeviceModelAdapter) AttachCategorySpec(spec base.CategorySpec) {
	a.bridge.CategorySpec = spec
}

// GetDeviceID returns the device ID
func (a *DeviceModelAdapter) GetDeviceID() uuid.UUID {
	return a.bridge.DeviceID
}

// GetCategorySpec returns the attached category specification
func (a *DeviceModelAdapter) GetCategorySpec() base.CategorySpec {
	return a.bridge.CategorySpec
}

// SetNetworkProfile sets reference to existing DeviceNetworkProfile
func (a *DeviceModelAdapter) SetNetworkProfile(id uuid.UUID) {
	a.bridge.NetworkProfileID = &id
}

// SetDataUsage sets reference to existing DeviceDataUsage
func (a *DeviceModelAdapter) SetDataUsage(id uuid.UUID) {
	a.bridge.DataUsageID = &id
}

// SetWiFiProfile sets reference to existing DeviceWiFiProfile
func (a *DeviceModelAdapter) SetWiFiProfile(id uuid.UUID) {
	a.bridge.WiFiProfileID = &id
}

// SetBluetoothProfile sets reference to existing DeviceBluetoothProfile
func (a *DeviceModelAdapter) SetBluetoothProfile(id uuid.UUID) {
	a.bridge.BluetoothProfileID = &id
}

// SetSecurityProfile sets reference to existing DeviceSecurity
func (a *DeviceModelAdapter) SetSecurityProfile(id uuid.UUID) {
	a.bridge.SecurityProfileID = &id
}

// SetPerformanceData sets reference to existing DevicePerformance
func (a *DeviceModelAdapter) SetPerformanceData(id uuid.UUID) {
	a.bridge.PerformanceDataID = &id
}

// SetQualityMetrics sets reference to existing DeviceQuality
func (a *DeviceModelAdapter) SetQualityMetrics(id uuid.UUID) {
	a.bridge.QualityMetricsID = &id
}

// SetSustainability sets reference to existing DeviceSustainability
func (a *DeviceModelAdapter) SetSustainability(id uuid.UUID) {
	a.bridge.SustainabilityID = &id
}

// SetEmergencyData sets reference to existing DeviceEmergency
func (a *DeviceModelAdapter) SetEmergencyData(id uuid.UUID) {
	a.bridge.EmergencyDataID = &id
}

// SetComplianceStatus sets reference to existing DeviceCompliance
func (a *DeviceModelAdapter) SetComplianceStatus(id uuid.UUID) {
	a.bridge.ComplianceStatusID = &id
}

// DeviceModelIntegrationService integrates the hybrid system with existing models
// This service would be used in the application layer to coordinate between systems
type DeviceModelIntegrationService struct {
	categoryService *CategoryService
}

// NewDeviceModelIntegrationService creates a new integration service
func NewDeviceModelIntegrationService(categoryService *CategoryService) *DeviceModelIntegrationService {
	return &DeviceModelIntegrationService{
		categoryService: categoryService,
	}
}

// CreateDeviceWithCategory creates a device using both the existing model and category system
func (s *DeviceModelIntegrationService) CreateDeviceWithCategory(
	deviceData map[string]interface{},
	categorySpec base.CategorySpec,
) (*DeviceModelAdapter, error) {
	// Extract device ID (would be created by the existing device service)
	deviceID, ok := deviceData["id"].(uuid.UUID)
	if !ok {
		deviceID = uuid.New()
	}

	// Create adapter
	adapter := NewDeviceModelAdapter(deviceID, categorySpec.GetCategory())
	adapter.AttachCategorySpec(categorySpec)

	// The actual device creation would happen in the existing device service
	// This adapter just bridges the two systems

	return adapter, nil
}

// EnrichDeviceWithCategoryFeatures enriches existing device with category-specific features
func (s *DeviceModelIntegrationService) EnrichDeviceWithCategoryFeatures(
	deviceID uuid.UUID,
	category base.CategoryType,
	specifications map[string]interface{},
) error {
	// Create category spec from specifications
	spec, err := s.categoryService.CreateSpec(category)
	if err != nil {
		return err
	}

	// Validate the specifications
	if err := s.categoryService.ValidateDevice(category, spec); err != nil {
		return err
	}

	// The enrichment would update the existing device model
	// with category-specific data stored in JSONB

	return nil
}

// CalculateInsuranceForDevice calculates insurance using both systems
func (s *DeviceModelIntegrationService) CalculateInsuranceForDevice(
	deviceID uuid.UUID,
	spec base.CategorySpec,
	coverage base.CoverageLevel,
) (InsuranceCalculation, error) {
	// Use category system for premium calculation
	premium, err := s.categoryService.CalculatePremium(spec, coverage)
	if err != nil {
		return InsuranceCalculation{}, err
	}

	// Get risk assessment
	claimHistory := base.ClaimHistory{} // Would be loaded from existing claim models
	riskAssessment, err := s.categoryService.AssessRisk(spec, claimHistory)
	if err != nil {
		return InsuranceCalculation{}, err
	}

	return InsuranceCalculation{
		DeviceID:       deviceID,
		Category:       spec.GetCategory(),
		Premium:        premium,
		RiskAssessment: riskAssessment,
		CalculatedAt:   time.Now(),
	}, nil
}

// InsuranceCalculation represents the result of insurance calculation
type InsuranceCalculation struct {
	DeviceID       uuid.UUID           `json:"device_id"`
	Category       base.CategoryType   `json:"category"`
	Premium        float64             `json:"premium"`
	RiskAssessment base.RiskAssessment `json:"risk_assessment"`
	CalculatedAt   time.Time           `json:"calculated_at"`
}

// GetDeviceWithAllFeatures retrieves device with all sub-entity features
// This would be implemented in the application service layer
func (s *DeviceModelIntegrationService) GetDeviceWithAllFeatures(deviceID uuid.UUID) (map[string]interface{}, error) {
	// This method would:
	// 1. Load the core device from existing device model
	// 2. Load all related sub-entities (network, performance, quality, etc.)
	// 3. Load the category-specific spec from hybrid system
	// 4. Combine everything into a comprehensive response

	result := make(map[string]interface{})
	result["device_id"] = deviceID

	// The actual implementation would load from repositories
	// Example structure:
	result["core"] = map[string]interface{}{
		"imei":          "123456789012345",
		"serial_number": "SN123456",
		"status":        "active",
	}

	result["network"] = map[string]interface{}{
		// From DeviceNetworkProfile
		"carrier":        "Verizon",
		"signal_quality": 85,
		"5g_capable":     true,
	}

	result["performance"] = map[string]interface{}{
		// From DevicePerformance
		"battery_health": 92,
		"storage_used":   45.5,
		"cpu_usage":      23.4,
	}

	result["category_features"] = map[string]interface{}{
		// From category spec
		"category":     "smartphone",
		"manufacturer": "Apple",
		"model":        "iPhone 14 Pro",
	}

	return result, nil
}
