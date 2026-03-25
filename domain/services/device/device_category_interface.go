package device

import (
	"errors"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

var (
	ErrInvalidCategory = errors.New("invalid device category")
	ErrInvalidIMEI     = errors.New("invalid IMEI")
)

// DeviceCategory represents the base interface for all device categories
type DeviceCategory interface {
	// Category identification
	GetCategory() string
	GetSubCategory() string

	// Category-specific features
	GetSpecificFeatures() map[string]interface{}
	ValidateForCategory() error
	GetCategorySpecificRisks() []CategoryRisk

	// Category-specific calculations
	CalculateCategoryPremiumAdjustment() float64
	GetCategoryDepreciationRate() float64
	GetCategorySpecificCoverage() []CoverageType

	// Category-specific eligibility
	IsEligibleForCategoryPrograms() map[string]bool
	GetCategoryMaintenanceSchedule() MaintenanceSchedule
}

// CategoryRisk represents a category-specific risk
type CategoryRisk struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"` // low, medium, high, critical
	Impact      float64 `json:"impact"`   // Impact on premium (multiplier)
}

// CoverageType represents a type of coverage
type CoverageType struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MaxAmount   float64 `json:"max_amount"`
	Deductible  float64 `json:"deductible"`
	IsRequired  bool    `json:"is_required"`
	IsAvailable bool    `json:"is_available"`
}

// MaintenanceSchedule represents device maintenance schedule
type MaintenanceSchedule struct {
	Intervals []MaintenanceInterval `json:"intervals"`
	NextDue   time.Time             `json:"next_due"`
	LastDone  time.Time             `json:"last_done"`
}

// MaintenanceInterval represents a maintenance interval
type MaintenanceInterval struct {
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Frequency   time.Duration `json:"frequency"`
	IsCritical  bool          `json:"is_critical"`
	Cost        float64       `json:"estimated_cost"`
}

// === Category Factory ===

// DeviceCategoryFactory creates category-specific device instances
type DeviceCategoryFactory struct{}

// NewDeviceCategoryFactory creates a new category factory
func NewDeviceCategoryFactory() *DeviceCategoryFactory {
	return &DeviceCategoryFactory{}
}

// CreateCategoryDevice creates a category-specific device based on category
func (f *DeviceCategoryFactory) CreateCategoryDevice(device *models.Device) (DeviceCategory, error) {
	category := device.DeviceClassification.Category

	switch category {
	case "smartphone":
		return NewSmartphoneDevice(device), nil
	case "smartwatch", "wearable":
		return NewSmartwatchDevice(device), nil
	case "tablet":
		return NewTabletDevice(device), nil
	case "laptop":
		return NewLaptopDevice(device), nil
	default:
		return NewGenericDevice(device), nil
	}
}

// GetDeviceCategory determines device category from device model
func (f *DeviceCategoryFactory) GetDeviceCategory(device *models.Device) string {
	// Use classification category if available
	if device.DeviceClassification.Category != "" {
		return string(device.DeviceClassification.Category)
	}

	// Fallback to model analysis
	model := device.DeviceClassification.Model
	_ = device.DeviceClassification.Brand // reserved for future use

	// Smartphone detection
	if contains(model, []string{"iPhone", "Galaxy", "Pixel", "OnePlus", "Xiaomi"}) {
		return "smartphone"
	}

	// Smartwatch detection
	if contains(model, []string{"Watch", "Band", "Fit", "Gear"}) {
		return "smartwatch"
	}

	// Tablet detection
	if contains(model, []string{"iPad", "Tab", "Tablet"}) {
		return "tablet"
	}

	// Laptop detection
	if contains(model, []string{"MacBook", "ThinkPad", "XPS", "Surface"}) {
		return "laptop"
	}

	return "unknown"
}

// === Generic Device Implementation (Base) ===

// GenericDevice represents a generic device implementation
type GenericDevice struct {
	*models.Device
}

// NewGenericDevice creates a new generic device
func NewGenericDevice(base *models.Device) *GenericDevice {
	return &GenericDevice{Device: base}
}

// GetCategory returns the device category
func (g *GenericDevice) GetCategory() string {
	if g.Device.DeviceClassification.Category != "" {
		return string(g.Device.DeviceClassification.Category)
	}
	return "generic"
}

// GetSubCategory returns the device sub-category
func (g *GenericDevice) GetSubCategory() string {
	return "standard"
}

// GetSpecificFeatures returns generic features
func (g *GenericDevice) GetSpecificFeatures() map[string]interface{} {
	return map[string]interface{}{
		"brand":  g.Device.DeviceClassification.Brand,
		"model":  g.Device.DeviceClassification.Model,
		"os":     g.Device.DeviceSpecifications.OperatingSystem,
		"status": string(g.Device.DeviceStatusInfo.Status),
	}
}

// ValidateForCategory validates generic requirements
func (g *GenericDevice) ValidateForCategory() error {
	if g.Device.ID == uuid.Nil {
		return errors.New("device ID required")
	}
	return nil
}

// GetCategorySpecificRisks returns generic risks
func (g *GenericDevice) GetCategorySpecificRisks() []CategoryRisk {
	return []CategoryRisk{
		{
			Type:        "general_damage",
			Description: "General wear and tear",
			Severity:    "medium",
			Impact:      1.0,
		},
	}
}

// CalculateCategoryPremiumAdjustment returns generic adjustment
func (g *GenericDevice) CalculateCategoryPremiumAdjustment() float64 {
	return 1.0 // No adjustment for generic
}

// GetCategoryDepreciationRate returns generic depreciation rate
func (g *GenericDevice) GetCategoryDepreciationRate() float64 {
	return 0.25 // 25% per year standard
}

// GetCategorySpecificCoverage returns generic coverage
func (g *GenericDevice) GetCategorySpecificCoverage() []CoverageType {
	return []CoverageType{
		{
			Type:        "damage",
			Name:        "Damage Protection",
			Description: "Covers accidental damage",
			MaxAmount:   g.Device.DeviceFinancial.CurrentValue.Amount,
			Deductible:  50,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "theft",
			Name:        "Theft Protection",
			Description: "Covers device theft",
			MaxAmount:   g.Device.DeviceFinancial.CurrentValue.Amount,
			Deductible:  100,
			IsRequired:  false,
			IsAvailable: true,
		},
	}
}

// IsEligibleForCategoryPrograms returns generic eligibility
func (g *GenericDevice) IsEligibleForCategoryPrograms() map[string]bool {
	return map[string]bool{
		"trade_in":    true,
		"repair":      true,
		"replacement": true,
	}
}

// GetCategoryMaintenanceSchedule returns generic maintenance
func (g *GenericDevice) GetCategoryMaintenanceSchedule() MaintenanceSchedule {
	return MaintenanceSchedule{
		Intervals: []MaintenanceInterval{
			{
				Type:        "inspection",
				Description: "General device inspection",
				Frequency:   180 * 24 * time.Hour,
				IsCritical:  false,
				Cost:        0,
			},
		},
		NextDue:  time.Now().Add(180 * 24 * time.Hour),
		LastDone: time.Now(),
	}
}

// Helper function
func contains(str string, substrs []string) bool {
	for _, substr := range substrs {
		if len(str) >= len(substr) && str[:len(substr)] == substr {
			return true
		}
	}
	return false
}
