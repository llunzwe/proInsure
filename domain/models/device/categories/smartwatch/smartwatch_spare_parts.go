package smartwatch

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartwatchSparePart represents smartwatch-specific spare parts
type SmartwatchSparePart struct {
	database.BaseModel

	DeviceID           uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	PartType           string    `gorm:"type:varchar(50);not null" json:"part_type"`
	PartNumber         string    `gorm:"uniqueIndex" json:"part_number"`
	PartName           string    `json:"part_name"`
	Manufacturer       string    `json:"manufacturer"`
	SupplierName       string    `json:"supplier_name"`
	SupplierPartNumber string    `json:"supplier_part_number"`

	// Compatibility
	CompatibleModels   string `gorm:"type:json" json:"compatible_models"` // JSON array of watch models
	CompatibilityNotes string `json:"compatibility_notes"`
	WatchGeneration    string `json:"watch_generation"` // Series 8, Series 9, Gen 6, etc
	WatchSize          string `json:"watch_size"`       // 38mm, 40mm, 41mm, 42mm, 44mm, 45mm, 49mm

	// Quality & authenticity
	Quality        string `gorm:"type:varchar(20);default:'oem'" json:"quality"`
	IsOriginalPart bool   `gorm:"default:false" json:"is_original_part"`
	IsCertified    bool   `gorm:"default:false" json:"is_certified"`
	QualityGrade   string `gorm:"type:varchar(5)" json:"quality_grade"`

	// Inventory & pricing
	QuantityInStock int     `json:"quantity_in_stock"`
	MinimumStock    int     `json:"minimum_stock"`
	ReorderPoint    int     `json:"reorder_point"`
	UnitCost        float64 `json:"unit_cost"`
	RetailPrice     float64 `json:"retail_price"`
	Currency        string  `gorm:"default:'USD'" json:"currency"`

	// Display Components
	DisplayType        string  `json:"display_type"` // OLED, LTPO OLED, Retina
	DisplaySize        float64 `json:"display_size"` // inches
	Resolution         string  `json:"resolution"`   // 368x448, 396x484, etc
	AlwaysOnCapable    bool    `gorm:"default:false" json:"always_on_capable"`
	TouchLayerIncluded bool    `gorm:"default:false" json:"touch_layer_included"`
	ForceTouch         bool    `gorm:"default:false" json:"force_touch"`
	DisplayBrightness  int     `json:"display_brightness"` // nits
	SapphireCrystal    bool    `gorm:"default:false" json:"sapphire_crystal"`
	IonXGlass          bool    `gorm:"default:false" json:"ion_x_glass"`

	// Battery Components
	BatteryCapacity   int     `json:"battery_capacity"` // mAh
	BatteryVoltage    float64 `json:"battery_voltage"`  // V
	BatteryType       string  `json:"battery_type"`     // Li-Ion, Li-Po
	BatteryShape      string  `json:"battery_shape"`    // rectangular, circular
	FastChargeSupport bool    `gorm:"default:false" json:"fast_charge_support"`

	// Health Sensors
	SensorType        string `json:"sensor_type"`       // heart_rate, SpO2, ECG, temperature
	SensorGeneration  string `json:"sensor_generation"` // 1st gen, 2nd gen, etc
	OpticalSensors    int    `json:"optical_sensors"`   // number of optical sensors
	ElectricalSensors bool   `gorm:"default:false" json:"electrical_sensors"`
	SensorAccuracy    string `json:"sensor_accuracy"` // medical_grade, fitness_grade

	// Digital Crown & Buttons
	DigitalCrownIncluded bool   `gorm:"default:false" json:"digital_crown_included"`
	HapticFeedback       bool   `gorm:"default:false" json:"haptic_feedback"`
	RotaryEncoder        bool   `gorm:"default:false" json:"rotary_encoder"`
	SideButton           bool   `gorm:"default:false" json:"side_button"`
	FunctionButton       bool   `gorm:"default:false" json:"function_button"`
	ButtonMechanism      string `json:"button_mechanism"` // physical, capacitive, pressure

	// Haptic Components
	HapticMotor        bool   `gorm:"default:false" json:"haptic_motor"`
	TapticEngine       bool   `gorm:"default:false" json:"taptic_engine"`
	LinearActuator     bool   `gorm:"default:false" json:"linear_actuator"`
	VibrationIntensity string `json:"vibration_intensity"` // light, medium, strong

	// Audio Components
	Speaker           bool `gorm:"default:false" json:"speaker"`
	Microphone        bool `gorm:"default:false" json:"microphone"`
	WaterEjection     bool `gorm:"default:false" json:"water_ejection"`
	SpeakerWaterproof bool `gorm:"default:false" json:"speaker_waterproof"`

	// Charging Components
	ChargingCoil      bool   `gorm:"default:false" json:"charging_coil"`
	ChargingPins      int    `json:"charging_pins"` // number of pins
	MagneticAlignment bool   `gorm:"default:false" json:"magnetic_alignment"`
	ChargingConnector string `json:"charging_connector"` // proprietary, pogo_pins

	// Case & Housing
	CaseMaterial          string `json:"case_material"` // aluminum, steel, titanium, ceramic
	CaseColor             string `json:"case_color"`
	BackCrystal           bool   `gorm:"default:false" json:"back_crystal"`
	BackCaseMaterial      string `json:"back_case_material"`      // ceramic, composite
	WaterResistanceRating string `json:"water_resistance_rating"` // 50m, 100m, IP68

	// Band Connector
	BandConnectorType    string `json:"band_connector_type"`    // slide, click, pin
	BandReleaseMechanism string `json:"band_release_mechanism"` // button_release, slide_release
	ConnectorMaterial    string `json:"connector_material"`

	// Antenna Components
	AntennaType       string `json:"antenna_type"` // Bluetooth, WiFi, Cellular, GPS
	AntennaIntegrated bool   `gorm:"default:false" json:"antenna_integrated"`

	// Board & Processing
	ProcessorType       string `json:"processor_type"` // S8, S9, Snapdragon W5
	CoprocessorIncluded bool   `gorm:"default:false" json:"coprocessor_included"`
	MemoryIncluded      bool   `gorm:"default:false" json:"memory_included"`
	StorageIncluded     bool   `gorm:"default:false" json:"storage_included"`

	// Small Parts
	SealGasket     bool `gorm:"default:false" json:"seal_gasket"`
	Screws         bool `gorm:"default:false" json:"screws"`
	AdhesiveStrips bool `gorm:"default:false" json:"adhesive_strips"`
	FlexCable      bool `gorm:"default:false" json:"flex_cable"`

	// Installation Requirements
	RequiresProfessional bool   `gorm:"default:true" json:"requires_professional"`
	RequiresCalibration  bool   `gorm:"default:false" json:"requires_calibration"`
	RequiresPairing      bool   `gorm:"default:false" json:"requires_pairing"`
	EstimatedRepairTime  int    `json:"estimated_repair_time"` // minutes
	DifficultyLevel      string `json:"difficulty_level"`      // easy, medium, hard, expert
	SpecialToolsRequired string `gorm:"type:json" json:"special_tools_required"`

	// Warranty & lifecycle
	WarrantyDays     int        `json:"warranty_days"`
	ExpectedLifespan int        `json:"expected_lifespan_days"`
	ManufactureDate  *time.Time `json:"manufacture_date"`
	BatchNumber      string     `json:"batch_number"`

	// Documentation
	DatasheetURL         string `json:"datasheet_url"`
	InstallationGuideURL string `json:"installation_guide_url"`
	CalibrationGuideURL  string `json:"calibration_guide_url"`
	PhotoURLs            string `gorm:"type:json" json:"photo_urls"`

	// Status & availability
	Status        string     `gorm:"type:varchar(20);default:'available'" json:"status"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	LastOrderDate *time.Time `json:"last_order_date"`

	// Performance metrics
	DefectRate           float64 `json:"defect_rate"`           // Percentage
	ReturnRate           float64 `json:"return_rate"`           // Percentage
	CustomerSatisfaction float64 `json:"customer_satisfaction"` // 0-5 rating
	InstallSuccessRate   float64 `json:"install_success_rate"`  // Percentage
}

// Business logic methods

// IsDisplayPart checks if this is a display component
func (sp *SmartwatchSparePart) IsDisplayPart() bool {
	displayParts := []string{"display", "screen", "glass", "crystal", "digitizer", "touch_layer"}
	for _, part := range displayParts {
		if sp.PartType == part {
			return true
		}
	}
	return sp.DisplayType != "" || sp.SapphireCrystal || sp.IonXGlass
}

// IsBatteryPart checks if this is a battery
func (sp *SmartwatchSparePart) IsBatteryPart() bool {
	return sp.PartType == "battery" || sp.BatteryCapacity > 0
}

// IsSensorPart checks if this is a health sensor
func (sp *SmartwatchSparePart) IsSensorPart() bool {
	sensorParts := []string{"sensor", "heart_rate_sensor", "spo2_sensor", "ecg_sensor", "temperature_sensor"}
	for _, part := range sensorParts {
		if sp.PartType == part {
			return true
		}
	}
	return sp.SensorType != "" || sp.OpticalSensors > 0
}

// IsHighValuePart checks if this is expensive
func (sp *SmartwatchSparePart) IsHighValuePart() bool {
	// Price threshold
	if sp.RetailPrice > 100 {
		return true
	}
	// Specific high-value parts
	highValueParts := []string{"display", "motherboard", "processor", "complete_case"}
	for _, part := range highValueParts {
		if sp.PartType == part {
			return true
		}
	}
	// Premium materials
	if sp.CaseMaterial == "titanium" || sp.CaseMaterial == "ceramic" || sp.SapphireCrystal {
		return true
	}
	return false
}

// IsCriticalComponent checks if critical for operation
func (sp *SmartwatchSparePart) IsCriticalComponent() bool {
	criticalParts := []string{
		"battery", "display", "motherboard", "processor",
		"charging_coil", "digital_crown",
	}
	for _, part := range criticalParts {
		if sp.PartType == part {
			return true
		}
	}
	return false
}

// RequiresSpecialHandling checks special handling needs
func (sp *SmartwatchSparePart) RequiresSpecialHandling() bool {
	// Battery always requires special handling
	if sp.IsBatteryPart() {
		return true
	}
	// Display requires careful handling
	if sp.IsDisplayPart() {
		return true
	}
	// Health sensors are sensitive
	if sp.IsSensorPart() {
		return true
	}
	return sp.RequiresProfessional || sp.RequiresCalibration
}

// GetWaterResistanceImpact checks if part affects water resistance
func (sp *SmartwatchSparePart) GetWaterResistanceImpact() bool {
	// These parts affect water resistance when replaced
	waterParts := []string{
		"display", "back_crystal", "seal_gasket", "case",
		"digital_crown", "side_button", "speaker",
	}
	for _, part := range waterParts {
		if sp.PartType == part {
			return true
		}
	}
	return sp.SealGasket || sp.BackCrystal
}

// NeedsReorder checks if reorder needed
func (sp *SmartwatchSparePart) NeedsReorder() bool {
	return sp.QuantityInStock <= sp.ReorderPoint
}

// IsCompatibleWithModel checks model compatibility
func (sp *SmartwatchSparePart) IsCompatibleWithModel(model string, size string) bool {
	// Check size first
	if sp.WatchSize != "" && sp.WatchSize != size {
		return false
	}
	// Would parse CompatibleModels JSON
	// Simplified for now
	return true
}

// GetInstallationComplexity returns complexity score
func (sp *SmartwatchSparePart) GetInstallationComplexity() int {
	complexity := 1 // Base

	// Professional requirement adds complexity
	if sp.RequiresProfessional {
		complexity += 3
	}

	// Calibration adds complexity
	if sp.RequiresCalibration {
		complexity += 2
	}

	// Pairing requirement
	if sp.RequiresPairing {
		complexity += 1
	}

	// Difficulty level
	switch sp.DifficultyLevel {
	case "easy":
		complexity += 1
	case "medium":
		complexity += 2
	case "hard":
		complexity += 3
	case "expert":
		complexity += 5
	}

	// Water resistance impact
	if sp.GetWaterResistanceImpact() {
		complexity += 2
	}

	return complexity
}

// EstimateInstallationCost estimates repair cost
func (sp *SmartwatchSparePart) EstimateInstallationCost() float64 {
	baseCost := 30.0 // Base for smartwatch repair

	// Complexity multiplier
	complexity := sp.GetInstallationComplexity()
	baseCost += float64(complexity) * 10

	// Time-based cost
	if sp.EstimatedRepairTime > 0 {
		baseCost += float64(sp.EstimatedRepairTime) * 0.75 // $0.75 per minute
	}

	// Premium for certain parts
	if sp.IsDisplayPart() {
		baseCost *= 1.5
	}
	if sp.IsSensorPart() {
		baseCost *= 1.3
	}
	if sp.RequiresCalibration {
		baseCost += 25 // Calibration fee
	}

	return baseCost
}

// GetQualityScore returns quality rating
func (sp *SmartwatchSparePart) GetQualityScore() float64 {
	score := 50.0 // Base

	// Grade scoring
	switch sp.QualityGrade {
	case "A+":
		score = 100
	case "A":
		score = 90
	case "B":
		score = 75
	case "C":
		score = 60
	}

	// OEM bonus
	if sp.IsOriginalPart {
		score *= 1.2
	}

	// Certification bonus
	if sp.IsCertified {
		score *= 1.1
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	// Factor in defect rate
	if sp.DefectRate > 0 {
		score *= (1 - sp.DefectRate/100)
	}

	return score
}

// IsHealthSensorPart checks if health-related
func (sp *SmartwatchSparePart) IsHealthSensorPart() bool {
	healthSensors := []string{"heart_rate", "SpO2", "ECG", "temperature", "blood_pressure"}
	for _, sensor := range healthSensors {
		if sp.SensorType == sensor {
			return true
		}
	}
	return sp.ElectricalSensors || sp.SensorAccuracy == "medical_grade"
}

// GetWarrantyStatus returns warranty status
func (sp *SmartwatchSparePart) GetWarrantyStatus() string {
	if sp.WarrantyDays == 0 {
		return "no_warranty"
	}
	if sp.ManufactureDate == nil {
		return "unknown"
	}

	warrantyExpiry := sp.ManufactureDate.AddDate(0, 0, sp.WarrantyDays)
	if time.Now().After(warrantyExpiry) {
		return "expired"
	}

	daysRemaining := int(time.Until(warrantyExpiry).Hours() / 24)
	if daysRemaining <= 30 {
		return "expiring_soon"
	}

	return "active"
}

// IsInDemand checks if part is high demand
func (sp *SmartwatchSparePart) IsInDemand() bool {
	// Battery and display always high demand
	if sp.IsBatteryPart() || sp.IsDisplayPart() {
		return true
	}
	// Health sensors increasingly popular
	if sp.IsHealthSensorPart() {
		return true
	}
	// Low stock indicates demand
	if sp.QuantityInStock < sp.MinimumStock {
		return true
	}
	return false
}
