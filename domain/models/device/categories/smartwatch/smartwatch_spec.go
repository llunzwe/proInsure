package smartwatch

import (
	"encoding/json"
	"fmt"
	"time"

	"smartsure/internal/domain/models/device/categories/base"
)

// SmartwatchSpec represents smartwatch-specific specifications
// It acts as a bridge between the hybrid system and existing device models
type SmartwatchSpec struct {
	base.BaseCategorySpec

	// Display specifications
	ScreenSize       float64 `json:"screen_size_mm"` // in mm (diagonal)
	ScreenType       string  `json:"screen_type"`    // OLED, LCD, AMOLED, LTPO, E-ink, Retina
	DisplayType      string  `json:"display_type"`   // OLED, LCD, AMOLED, LTPO, E-ink, Retina (alias for compatibility)
	Resolution       string  `json:"resolution"`     // e.g., "368x448"
	PPI              int     `json:"ppi"`            // pixels per inch
	AlwaysOnDisplay  bool    `json:"always_on_display"`
	TouchScreen      bool    `json:"touch_screen"`
	Brightness       int     `json:"peak_brightness_nits"` // typically 1000-2000 nits
	ScreenProtection string  `json:"screen_protection"`    // Sapphire crystal, Ion-X glass, etc.

	// Health sensors
	HeartRateSensor   bool     `json:"heart_rate_sensor"`
	ECGSensor         bool     `json:"ecg_sensor"`
	BloodOxygenSensor bool     `json:"blood_oxygen_sensor"`
	TemperatureSensor bool     `json:"temperature_sensor"`
	HealthFeatures    []string `json:"health_features"` // sleep tracking, stress, etc.
	BiometricTypes    []string `json:"biometric_types"` // wrist detection, PIN, pattern

	// Fitness tracking
	GPSBuiltIn       bool     `json:"gps_built_in"`
	ActivityTracking []string `json:"activity_tracking"` // running, swimming, cycling, etc.
	WaterResistance  string   `json:"water_resistance"`  // e.g., "5ATM", "IP68"

	// Battery
	BatteryCapacity  int     `json:"battery_capacity_mah"` // in mAh (typically 200-500)
	BatteryLife      int     `json:"battery_life_days"`
	ChargingMethod   string  `json:"charging_method"` // magnetic, wireless, pins, USB-C
	ChargingTime     int     `json:"charging_time_min"`
	ChargingWattage  float64 `json:"charging_wattage"`  // typically 2.5W-5W
	WirelessCharging bool    `json:"wireless_charging"` // Qi compatible
	PowerSavingMode  bool    `json:"power_saving_mode"`
	BatteryType      string  `json:"battery_type"` // Li-Ion, Li-Po

	// Connectivity
	Cellular         bool     `json:"cellular"` // LTE/5G support
	eSIM             bool     `json:"esim"`     // eSIM support
	WiFi             bool     `json:"wifi"`
	WiFiStandards    []string `json:"wifi_standards"`    // 802.11b/g/n, WiFi 6
	Bluetooth        bool     `json:"bluetooth"`         // Bluetooth support
	BluetoothVersion string   `json:"bluetooth_version"` // 5.0, 5.3, etc.
	NFC              bool     `json:"nfc"`               // for payments

	// Smart features
	VoiceAssistant  string `json:"voice_assistant"`     // Siri, Google, Alexa, Bixby
	MobilePayments  bool   `json:"mobile_payments"`     // Apple Pay, Google Pay, etc.
	StorageCapacity int    `json:"storage_capacity_gb"` // internal storage in GB
	RAM             int    `json:"ram_mb"`              // RAM in MB (typically 512-2048)
	MusicStorage    int    `json:"music_storage_gb"`    // dedicated music storage
	AppStore        bool   `json:"app_store"`
	Processor       string `json:"processor"` // S8, Snapdragon W5+, Exynos, etc.

	// Build
	CaseMaterial string   `json:"case_material"` // aluminum, steel, titanium, plastic, ceramic
	CaseSize     []int    `json:"case_sizes_mm"` // e.g., [38, 40, 41, 42, 44, 45, 49]
	BandType     string   `json:"band_type"`     // sport, leather, metal, milanese, solo loop
	BandSize     []string `json:"band_sizes"`    // S/M, M/L, etc.
	Weight       int      `json:"weight_grams"`
	Color        string   `json:"color"`       // case color
	BandColors   []string `json:"band_colors"` // available band colors

	// Operating System
	OS            string   `json:"os"` // watchOS, Wear OS, Tizen, etc.
	OSVersion     string   `json:"os_version"`
	Compatibility []string `json:"compatibility"` // iOS, Android, etc.

	// Additional fields for insurance
	PurchaseDate time.Time `json:"purchase_date"`
	Condition    string    `json:"condition"`         // excellent, good, fair, poor
	DisplaySize  float64   `json:"display_size_inch"` // Display size in inches
}

// NewSmartwatchSpec creates a new smartwatch specification
func NewSmartwatchSpec() *SmartwatchSpec {
	spec := &SmartwatchSpec{}
	spec.Category = base.CategorySmartwatch
	spec.ID = spec.GetSpecID()
	return spec
}

// Validate validates smartwatch specifications
func (s *SmartwatchSpec) Validate() error {
	// Validate base specifications
	if s.Manufacturer == "" {
		return fmt.Errorf("manufacturer is required")
	}
	if s.Model == "" {
		return fmt.Errorf("model is required")
	}

	// Validate screen size (typical smartwatch range)
	if s.ScreenSize < 20 || s.ScreenSize > 60 {
		return fmt.Errorf("screen size must be between 20mm and 60mm")
	}

	// Validate battery life
	if s.BatteryLife < 0 || s.BatteryLife > 30 {
		return fmt.Errorf("battery life must be between 0 and 30 days")
	}

	// Validate weight
	if s.Weight < 20 || s.Weight > 200 {
		return fmt.Errorf("weight must be between 20g and 200g")
	}

	return nil
}

// GetInsuranceFactors returns factors specific to smartwatch insurance
func (s *SmartwatchSpec) GetInsuranceFactors() map[string]float64 {
	factors := make(map[string]float64)

	// Screen size factor (larger screens = higher risk)
	factors["screen_risk"] = s.ScreenSize / 45.0

	// Value factor based on features
	valueScore := 0.0
	if s.Cellular {
		valueScore += 0.25
	}
	if s.ECGSensor {
		valueScore += 0.15
	}
	if s.BloodOxygenSensor {
		valueScore += 0.1
	}
	if s.GPSBuiltIn {
		valueScore += 0.1
	}
	if s.MobilePayments {
		valueScore += 0.1
	}

	// Material quality factor
	materialScore := 0.5
	switch s.CaseMaterial {
	case "titanium":
		materialScore = 1.0
	case "steel", "stainless_steel":
		materialScore = 0.8
	case "aluminum":
		materialScore = 0.6
	case "plastic", "resin":
		materialScore = 0.3
	}

	factors["value_score"] = valueScore
	factors["material_quality"] = materialScore

	// Water resistance factor (better resistance = lower risk)
	waterScore := 0.5
	switch s.WaterResistance {
	case "10ATM", "100m":
		waterScore = 0.2
	case "5ATM", "50m", "IP68":
		waterScore = 0.3
	case "3ATM", "30m", "IP67":
		waterScore = 0.5
	case "IPX7", "IPX6":
		waterScore = 0.7
	default:
		waterScore = 1.0
	}
	factors["water_damage_risk"] = waterScore

	// Battery replacement difficulty
	factors["battery_risk"] = 0.8 // Most smartwatches have non-replaceable batteries

	return factors
}

// GetRepairCostEstimate estimates repair cost for common issues
func (s *SmartwatchSpec) GetRepairCostEstimate(issueType string) float64 {
	baseCost := s.MarketValue * 0.3

	switch issueType {
	case "screen_crack":
		return baseCost * 0.6
	case "battery_replacement":
		return baseCost * 0.4
	case "water_damage":
		return baseCost * 0.8
	case "band_replacement":
		return 50.0 // Fixed cost for band
	case "sensor_failure":
		return baseCost * 0.5
	case "button_malfunction":
		return baseCost * 0.3
	default:
		return baseCost
	}
}

// ToJSON converts specification to JSON
func (s *SmartwatchSpec) ToJSON() (json.RawMessage, error) {
	return json.Marshal(s)
}

// FromJSON loads specification from JSON
func (s *SmartwatchSpec) FromJSON(data json.RawMessage) error {
	return json.Unmarshal(data, s)
}

// GetHealthFeatureScore calculates a health feature score for risk assessment
func (s *SmartwatchSpec) GetHealthFeatureScore() float64 {
	score := 0.0

	if s.HeartRateSensor {
		score += 1.0
	}
	if s.ECGSensor {
		score += 2.0
	}
	if s.BloodOxygenSensor {
		score += 1.5
	}
	if s.TemperatureSensor {
		score += 1.0
	}

	// Add points for each health feature
	score += float64(len(s.HealthFeatures)) * 0.5

	return score
}
