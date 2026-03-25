package smartphone

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartphoneSparePart represents smartphone-specific spare parts
type SmartphoneSparePart struct {
	database.BaseModel

	DeviceID           uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	PartType           string    `gorm:"type:varchar(50);not null" json:"part_type"`
	PartNumber         string    `gorm:"uniqueIndex" json:"part_number"`
	PartName           string    `json:"part_name"`
	Manufacturer       string    `json:"manufacturer"`
	SupplierName       string    `json:"supplier_name"`
	SupplierPartNumber string    `json:"supplier_part_number"`

	// Compatibility
	CompatibleModels   string `gorm:"type:json" json:"compatible_models"` // JSON array of phone models
	CompatibilityNotes string `json:"compatibility_notes"`
	IsUniversal        bool   `gorm:"default:false" json:"is_universal"`
	ModelGeneration    string `json:"model_generation"` // iPhone 12 series, Galaxy S21 series

	// Quality & authenticity
	Quality           string `gorm:"type:varchar(20);default:'oem'" json:"quality"` // oem, genuine, aftermarket, refurbished
	IsOriginalPart    bool   `gorm:"default:false" json:"is_original_part"`
	IsCertified       bool   `gorm:"default:false" json:"is_certified"`
	CertificationBody string `json:"certification_body"`
	QualityGrade      string `gorm:"type:varchar(5)" json:"quality_grade"` // A+, A, B, C

	// Inventory & pricing
	QuantityInStock int     `json:"quantity_in_stock"`
	MinimumStock    int     `json:"minimum_stock"`
	ReorderPoint    int     `json:"reorder_point"`
	UnitCost        float64 `json:"unit_cost"`
	RetailPrice     float64 `json:"retail_price"`
	WholesalePrice  float64 `json:"wholesale_price"`
	Currency        string  `gorm:"default:'USD'" json:"currency"`

	// Smartphone Display Components
	DisplayType        string  `json:"display_type"` // LCD, OLED, AMOLED, Super AMOLED, LTPO
	DisplaySize        float64 `json:"display_size"` // inches
	Resolution         string  `json:"resolution"`   // 1080x2340, 1170x2532, etc
	RefreshRate        int     `json:"refresh_rate"` // 60Hz, 90Hz, 120Hz, 144Hz
	TouchLayerIncluded bool    `gorm:"default:false" json:"touch_layer_included"`
	WithFrame          bool    `gorm:"default:false" json:"with_frame"`
	ColorAccuracy      string  `json:"color_accuracy"`  // sRGB percentage
	PeakBrightness     int     `json:"peak_brightness"` // nits

	// Camera Components
	CameraType           string  `json:"camera_type"`       // front, rear_main, rear_wide, rear_telephoto, rear_macro
	CameraMegapixels     float64 `json:"camera_megapixels"` // 12MP, 48MP, 108MP, etc
	CameraAperture       string  `json:"camera_aperture"`   // f/1.8, f/2.0, etc
	OpticalStabilization bool    `gorm:"default:false" json:"optical_stabilization"`
	AutofocusType        string  `json:"autofocus_type"` // PDAF, laser, contrast
	WithFlexCable        bool    `gorm:"default:false" json:"with_flex_cable"`

	// Battery Components
	BatteryCapacity      int     `json:"battery_capacity"`      // mAh
	BatteryVoltage       float64 `json:"battery_voltage"`       // V
	BatteryType          string  `json:"battery_type"`          // Li-Ion, Li-Po
	BatteryCycles        int     `json:"battery_cycles"`        // for refurbished
	FastChargingSupport  string  `json:"fast_charging_support"` // wattage supported
	WirelessChargingCoil bool    `gorm:"default:false" json:"wireless_charging_coil"`

	// Charging Port Components
	ChargingPortType    string `json:"charging_port_type"`    // USB-C, Lightning, Micro-USB
	ChargingPortVersion string `json:"charging_port_version"` // USB 2.0, USB 3.0, Thunderbolt
	WithFlexRibbon      bool   `gorm:"default:false" json:"with_flex_ribbon"`
	WaterResistantSeal  bool   `gorm:"default:false" json:"water_resistant_seal"`

	// Audio Components
	SpeakerType       string `json:"speaker_type"`    // earpiece, loud_speaker, bottom_speaker
	MicrophoneType    string `json:"microphone_type"` // primary, secondary, noise_cancelling
	AudioJackIncluded bool   `gorm:"default:false" json:"audio_jack_included"`

	// Button Components
	ButtonType            string `json:"button_type"`      // power, volume_up, volume_down, home, mute_switch
	ButtonMechanism       string `json:"button_mechanism"` // physical, capacitive, haptic
	WithFlexCableAttached bool   `gorm:"default:false" json:"with_flex_cable_attached"`

	// Sensor Components
	SensorType       string `json:"sensor_type"`       // proximity, ambient_light, fingerprint, face_id
	BiometricType    string `json:"biometric_type"`    // optical, ultrasonic, capacitive
	SensorGeneration string `json:"sensor_generation"` // for Face ID/Touch ID generations

	// Board & Chip Components
	BoardType       string `json:"board_type"`       // motherboard, logic_board, sub_board
	ChipsetModel    string `json:"chipset_model"`    // A15, Snapdragon 888, Exynos 2100
	StorageCapacity int    `json:"storage_capacity"` // GB for storage chips
	RAMCapacity     int    `json:"ram_capacity"`     // GB for RAM chips

	// Housing & Frame Components
	HousingMaterial string `json:"housing_material"` // aluminum, glass, plastic, ceramic, titanium
	HousingColor    string `json:"housing_color"`
	FrameIncluded   bool   `gorm:"default:false" json:"frame_included"`
	BackGlassType   string `json:"back_glass_type"` // standard, frosted, ceramic_shield

	// Small Parts & Hardware
	SIMTrayType        string `json:"sim_tray_type"` // single, dual, eSIM
	ScrewKit           bool   `gorm:"default:false" json:"screw_kit"`
	AdhesiveType       string `json:"adhesive_type"` // pre-cut, liquid, tape
	WaterDamageSticker bool   `gorm:"default:false" json:"water_damage_sticker"`

	// Installation Requirements
	RequiresProfessional bool   `gorm:"default:false" json:"requires_professional"`
	RequiresCalibration  bool   `gorm:"default:false" json:"requires_calibration"`
	RequiresProgramming  bool   `gorm:"default:false" json:"requires_programming"`
	EstimatedRepairTime  int    `json:"estimated_repair_time"`                   // minutes
	DifficultyLevel      string `json:"difficulty_level"`                        // easy, medium, hard, expert
	SpecialToolsRequired string `gorm:"type:json" json:"special_tools_required"` // JSON array

	// Warranty & lifecycle
	WarrantyDays     int        `json:"warranty_days"`
	ExpectedLifespan int        `json:"expected_lifespan_days"`
	ManufactureDate  *time.Time `json:"manufacture_date"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	BatchNumber      string     `json:"batch_number"`

	// Documentation
	DatasheetURL         string `json:"datasheet_url"`
	InstallationGuideURL string `json:"installation_guide_url"`
	PhotoURLs            string `gorm:"type:json" json:"photo_urls"`
	VideoGuideURL        string `json:"video_guide_url"`
	CalibrationGuideURL  string `json:"calibration_guide_url"`

	// Status & availability
	Status       string     `gorm:"type:varchar(20);default:'available'" json:"status"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	IsRecalled   bool       `gorm:"default:false" json:"is_recalled"`
	RecallDate   *time.Time `json:"recall_date"`
	RecallReason string     `json:"recall_reason"`

	// Performance metrics
	DefectRate           float64 `json:"defect_rate"`           // Percentage
	ReturnRate           float64 `json:"return_rate"`           // Percentage
	CustomerSatisfaction float64 `json:"customer_satisfaction"` // 0-5 rating
	InstallSuccessRate   float64 `json:"install_success_rate"`  // Percentage
}

// Business logic methods

// IsDisplayPart checks if this is a display-related part
func (sp *SmartphoneSparePart) IsDisplayPart() bool {
	displayParts := []string{"screen", "display", "lcd", "oled", "digitizer", "touch_screen", "glass"}
	for _, part := range displayParts {
		if sp.PartType == part {
			return true
		}
	}
	return sp.DisplayType != ""
}

// IsCameraPart checks if this is a camera-related part
func (sp *SmartphoneSparePart) IsCameraPart() bool {
	cameraParts := []string{"camera", "camera_lens", "camera_module", "camera_glass"}
	for _, part := range cameraParts {
		if sp.PartType == part {
			return true
		}
	}
	return sp.CameraType != ""
}

// IsBatteryPart checks if this is a battery-related part
func (sp *SmartphoneSparePart) IsBatteryPart() bool {
	return sp.PartType == "battery" || sp.BatteryCapacity > 0
}

// IsHighValuePart checks if this is a high-value part
func (sp *SmartphoneSparePart) IsHighValuePart() bool {
	// High value based on price
	if sp.RetailPrice > 150 {
		return true
	}
	// High value based on type
	highValueParts := []string{"motherboard", "logic_board", "display", "oled_display"}
	for _, part := range highValueParts {
		if sp.PartType == part {
			return true
		}
	}
	return false
}

// RequiresSpecialHandling checks if part needs special handling
func (sp *SmartphoneSparePart) RequiresSpecialHandling() bool {
	// Battery requires special handling
	if sp.IsBatteryPart() {
		return true
	}
	// Display requires careful handling
	if sp.IsDisplayPart() {
		return true
	}
	// Biometric sensors require special handling
	if sp.BiometricType != "" {
		return true
	}
	return sp.RequiresProfessional || sp.RequiresCalibration || sp.RequiresProgramming
}

// IsOEMQuality checks if part is OEM quality
func (sp *SmartphoneSparePart) IsOEMQuality() bool {
	return sp.Quality == "oem" || sp.IsOriginalPart
}

// GetInstallationComplexity returns installation complexity score
func (sp *SmartphoneSparePart) GetInstallationComplexity() int {
	complexity := 1 // Base complexity

	if sp.RequiresProfessional {
		complexity += 3
	}
	if sp.RequiresCalibration {
		complexity += 2
	}
	if sp.RequiresProgramming {
		complexity += 2
	}

	switch sp.DifficultyLevel {
	case "easy":
		complexity += 1
	case "medium":
		complexity += 2
	case "hard":
		complexity += 3
	case "expert":
		complexity += 4
	}

	return complexity
}

// IsHighDemandPart checks if part is in high demand
func (sp *SmartphoneSparePart) IsHighDemandPart() bool {
	// Screen and battery are always high demand
	if sp.IsDisplayPart() || sp.IsBatteryPart() {
		return true
	}
	// Charging ports are high demand
	if sp.PartType == "charging_port" || sp.ChargingPortType != "" {
		return true
	}
	// Check stock levels
	if sp.QuantityInStock < sp.MinimumStock {
		return true
	}
	return false
}

// NeedsReorder checks if part needs reordering
func (sp *SmartphoneSparePart) NeedsReorder() bool {
	return sp.QuantityInStock <= sp.ReorderPoint
}

// GetProfitMargin calculates profit margin
func (sp *SmartphoneSparePart) GetProfitMargin() float64 {
	if sp.UnitCost == 0 {
		return 0
	}
	return ((sp.RetailPrice - sp.UnitCost) / sp.UnitCost) * 100
}

// IsCompatibleWithModel checks if part is compatible with a specific model
func (sp *SmartphoneSparePart) IsCompatibleWithModel(model string) bool {
	// This would parse the CompatibleModels JSON and check
	// For now, return true as placeholder
	return true
}

// GetWarrantyStatus returns warranty status for the part
func (sp *SmartphoneSparePart) GetWarrantyStatus() string {
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

// EstimateInstallationCost estimates the installation cost
func (sp *SmartphoneSparePart) EstimateInstallationCost() float64 {
	baseCost := 20.0 // Base installation cost

	// Adjust based on complexity
	complexity := sp.GetInstallationComplexity()
	baseCost += float64(complexity) * 10

	// Adjust based on time
	if sp.EstimatedRepairTime > 0 {
		baseCost += float64(sp.EstimatedRepairTime) * 0.5 // $0.50 per minute
	}

	// Premium for certain parts
	if sp.IsDisplayPart() {
		baseCost *= 1.5
	}
	if sp.RequiresProgramming || sp.RequiresCalibration {
		baseCost *= 1.3
	}

	return baseCost
}

// IsCriticalComponent checks if this is a critical component
func (sp *SmartphoneSparePart) IsCriticalComponent() bool {
	criticalParts := []string{
		"motherboard", "logic_board", "display", "battery",
		"charging_port", "power_management_ic",
	}
	for _, part := range criticalParts {
		if sp.PartType == part {
			return true
		}
	}
	return false
}

// GetQualityScore returns a quality score for the part
func (sp *SmartphoneSparePart) GetQualityScore() float64 {
	score := 0.0

	// Quality grade scoring
	switch sp.QualityGrade {
	case "A+":
		score = 100
	case "A":
		score = 90
	case "B":
		score = 75
	case "C":
		score = 60
	default:
		score = 50
	}

	// Adjust for OEM parts
	if sp.IsOEMQuality() {
		score *= 1.2
	}

	// Adjust for certification
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
