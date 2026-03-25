package smartphone

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"smartsure/internal/domain/models/device/categories/base"
)

// SmartphoneSpec represents smartphone-specific specifications
type SmartphoneSpec struct {
	base.BaseCategorySpec

	// Display specifications
	ScreenSize       float64  `json:"screen_size_inches"`   // in inches
	ScreenType       string   `json:"screen_type"`          // OLED, LCD, AMOLED, LTPO, IPS
	Resolution       string   `json:"resolution"`           // e.g., "1920x1080"
	PPI              int      `json:"ppi"`                  // pixels per inch
	RefreshRate      int      `json:"refresh_rate_hz"`      // Hz (60, 90, 120, 144)
	TouchSampling    int      `json:"touch_sampling_rate"`  // Hz
	ScreenProtection string   `json:"screen_protection"`    // Gorilla Glass version
	HDRSupport       []string `json:"hdr_support"`          // HDR10, HDR10+, Dolby Vision
	Brightness       int      `json:"peak_brightness_nits"` // nits

	// Performance specifications
	Processor         string `json:"processor"`
	RAM               int    `json:"ram_gb"`              // in GB
	StorageCapacity   int    `json:"storage_capacity_gb"` // in GB
	ExpandableStorage bool   `json:"expandable_storage"`
	StorageType       string `json:"storage_type"` // UFS 3.1, UFS 4.0, eMMC

	// Camera specifications
	RearCameras     []CameraSpec `json:"rear_cameras"`
	FrontCameras    []CameraSpec `json:"front_cameras"`
	VideoCapability string       `json:"video_capability"` // e.g., "4K@60fps"

	// Battery specifications
	BatteryCapacity       int    `json:"battery_capacity_mah"`
	ChargingWattage       int    `json:"charging_wattage"` // max charging speed in watts
	WirelessCharging      bool   `json:"wireless_charging"`
	WirelessChargingWatts int    `json:"wireless_charging_watts"`
	ReverseWireless       bool   `json:"reverse_wireless_charging"`
	BatteryType           string `json:"battery_type"`           // Li-Ion, Li-Po
	FastChargingProtocol  string `json:"fast_charging_protocol"` // PD, QC, SuperVOOC, etc.

	// Connectivity
	NetworkSupport   []string `json:"network_support"` // 3G, 4G, 5G
	Is5GCapable      bool     `json:"is_5g_capable"`
	Network5GBands   []string `json:"5g_bands"`          // n1, n78, n79, etc.
	WiFiStandards    []string `json:"wifi_standards"`    // WiFi 6, WiFi 6E, WiFi 7
	BluetoothVersion string   `json:"bluetooth_version"` // 5.0, 5.3, etc.
	NFC              bool     `json:"nfc"`
	DualSIM          bool     `json:"dual_sim"`
	DualSIMType      string   `json:"dual_sim_type"` // physical, eSIM, hybrid

	// Security features
	BiometricTypes []string `json:"biometric_types"` // fingerprint, face_id
	SecurityChip   string   `json:"security_chip"`

	// Build quality
	WaterResistance string   `json:"water_resistance"` // IP rating
	Color           string   `json:"color"`            // device color
	Material        string   `json:"material"`         // glass, metal, plastic
	Materials       []string `json:"materials"`        // glass, aluminum, plastic
	Colors          []string `json:"available_colors"`

	// Operating System
	OSVersion       string `json:"os_version"`
	OSUpdateSupport int    `json:"os_update_years"` // Years of OS updates

	// Special features
	SpecialFeatures []string `json:"special_features"` // stylus, foldable, etc.
}

// CameraSpec represents camera specifications
type CameraSpec struct {
	Megapixels float64  `json:"megapixels"`
	Aperture   string   `json:"aperture"`
	SensorSize string   `json:"sensor_size"`
	Features   []string `json:"features"` // OIS, autofocus, etc.
	Type       string   `json:"type"`     // wide, ultrawide, telephoto, macro
}

// NewSmartphoneSpec creates a new smartphone specification
func NewSmartphoneSpec() *SmartphoneSpec {
	spec := &SmartphoneSpec{}
	spec.Category = base.CategorySmartphone
	spec.ID = uuid.New()
	return spec
}

// Validate validates smartphone specifications
func (s *SmartphoneSpec) Validate() error {
	// Validate base specifications
	if s.Manufacturer == "" {
		return fmt.Errorf("manufacturer is required")
	}
	if s.Model == "" {
		return fmt.Errorf("model is required")
	}

	// Validate screen size
	if s.ScreenSize < 3.0 || s.ScreenSize > 8.0 {
		return fmt.Errorf("screen size must be between 3.0 and 8.0 inches")
	}

	// Validate RAM
	if s.RAM < 1 || s.RAM > 32 {
		return fmt.Errorf("RAM must be between 1 and 32 GB")
	}

	// Validate storage
	if s.StorageCapacity < 8 || s.StorageCapacity > 2048 {
		return fmt.Errorf("storage must be between 8 and 2048 GB")
	}

	// Validate battery capacity
	if s.BatteryCapacity < 1000 || s.BatteryCapacity > 10000 {
		return fmt.Errorf("battery capacity must be between 1000 and 10000 mAh")
	}

	return nil
}

// GetInsuranceFactors returns factors specific to smartphone insurance
func (s *SmartphoneSpec) GetInsuranceFactors() map[string]float64 {
	factors := make(map[string]float64)

	// Screen size factor (larger screens = higher risk)
	factors["screen_risk"] = s.ScreenSize / 6.0

	// Value factor based on specs
	valueScore := 0.0
	if s.RAM >= 12 {
		valueScore += 0.3
	}
	if s.StorageCapacity >= 256 {
		valueScore += 0.2
	}
	if len(s.RearCameras) >= 3 {
		valueScore += 0.2
	}
	if s.Has5G() {
		valueScore += 0.2
	}
	if s.ScreenType == "OLED" || s.ScreenType == "AMOLED" {
		valueScore += 0.1
	}
	factors["value_factor"] = 1.0 + valueScore

	// Fragility factor
	fragilityScore := 1.0
	if s.WaterResistance == "IP68" || s.WaterResistance == "IP67" {
		fragilityScore *= 0.8
	}
	if s.ScreenProtection != "" {
		fragilityScore *= 0.9
	}
	factors["fragility"] = fragilityScore

	// Repairability factor
	repairabilityScore := 1.0
	if s.IsFoldable() {
		repairabilityScore *= 1.5 // Foldables are harder to repair
	}
	factors["repairability"] = repairabilityScore

	return factors
}

// Has5G checks if the device supports 5G
func (s *SmartphoneSpec) Has5G() bool {
	for _, network := range s.NetworkSupport {
		if network == "5G" {
			return true
		}
	}
	return false
}

// IsFoldable checks if the device is foldable
func (s *SmartphoneSpec) IsFoldable() bool {
	for _, feature := range s.SpecialFeatures {
		if feature == "foldable" || feature == "flip" {
			return true
		}
	}
	return false
}

// IsHighEnd determines if the smartphone is high-end based on specifications
func (s *SmartphoneSpec) IsHighEnd() bool {
	highEndScore := 0

	if s.MarketValue > 800 {
		highEndScore++
	}
	if s.RAM >= 8 {
		highEndScore++
	}
	if s.StorageCapacity >= 256 {
		highEndScore++
	}
	if s.Has5G() {
		highEndScore++
	}
	if s.ScreenType == "OLED" || s.ScreenType == "AMOLED" {
		highEndScore++
	}
	if len(s.RearCameras) >= 3 {
		highEndScore++
	}

	return highEndScore >= 4
}

// GetRepairCosts returns estimated repair costs for smartphone components
func (s *SmartphoneSpec) GetRepairCosts() map[string]float64 {
	costs := make(map[string]float64)

	// Base costs adjusted by device value
	valueFactor := s.MarketValue / 1000.0
	if valueFactor < 0.5 {
		valueFactor = 0.5
	}
	if valueFactor > 2.0 {
		valueFactor = 2.0
	}

	// Screen replacement
	screenCost := 150.0
	if s.ScreenType == "OLED" || s.ScreenType == "AMOLED" {
		screenCost *= 1.5
	}
	if s.IsFoldable() {
		screenCost *= 3.0
	}
	costs["screen"] = screenCost * valueFactor

	// Battery replacement
	costs["battery"] = 80.0 * valueFactor

	// Camera module
	cameraCost := 100.0
	if len(s.RearCameras) >= 3 {
		cameraCost *= 1.5
	}
	costs["camera"] = cameraCost * valueFactor

	// Charging port
	costs["charging_port"] = 60.0 * valueFactor

	// Back glass
	costs["back_glass"] = 80.0 * valueFactor

	// Motherboard
	costs["motherboard"] = s.MarketValue * 0.6

	// Other components
	costs["speaker"] = 40.0 * valueFactor
	costs["microphone"] = 30.0 * valueFactor
	costs["buttons"] = 25.0 * valueFactor

	return costs
}

// ToJSON converts smartphone spec to JSON
func (s *SmartphoneSpec) ToJSON() (json.RawMessage, error) {
	return json.Marshal(s)
}

// FromJSON populates smartphone spec from JSON
func (s *SmartphoneSpec) FromJSON(data json.RawMessage) error {
	return json.Unmarshal(data, s)
}

// GetWarrantyPeriod returns warranty period in months
func (s *SmartphoneSpec) GetWarrantyPeriod() int {
	if s.Warranty > 0 {
		return s.Warranty
	}
	// Default warranty periods based on value
	if s.IsHighEnd() {
		return 24
	}
	return 12
}

// GetCompatibleAccessories returns compatible accessory types
func (s *SmartphoneSpec) GetCompatibleAccessories() []string {
	accessories := []string{
		"case",
		"screen_protector",
		"charger",
		"cable",
		"headphones",
		"car_mount",
		"power_bank",
	}

	if s.WirelessCharging {
		accessories = append(accessories, "wireless_charger")
	}

	if s.NFC {
		accessories = append(accessories, "nfc_tags")
	}

	for _, feature := range s.SpecialFeatures {
		if feature == "stylus" {
			accessories = append(accessories, "stylus", "stylus_tips")
			break
		}
	}

	return accessories
}

// GetReleaseDate returns the release date
func (s *SmartphoneSpec) GetReleaseDate() time.Time {
	return s.ReleaseDate
}
