package tablet

import (
	"fmt"
	"strings"

	"smartsure/internal/domain/models/device/categories/base"
)

// TabletSpec defines specifications specific to tablet devices
type TabletSpec struct {
	base.BaseCategorySpec

	// Display specifications
	ScreenSize       float64  `json:"screen_size_inches"`   // in inches (typically 7-13 inches)
	ScreenType       string   `json:"screen_type"`          // LCD, IPS, OLED, AMOLED
	Resolution       string   `json:"resolution"`           // e.g., "2048x1536", "2732x2048"
	PPI              int      `json:"ppi"`                  // pixels per inch
	RefreshRate      int      `json:"refresh_rate_hz"`      // Hz (60, 90, 120)
	TouchSampling    int      `json:"touch_sampling_rate"`  // Hz
	ScreenProtection string   `json:"screen_protection"`    // Gorilla Glass version or equivalent
	HDRSupport       []string `json:"hdr_support"`          // HDR10, HDR10+, Dolby Vision
	Brightness       int      `json:"peak_brightness_nits"` // nits
	IsMultiTouch    bool     `json:"is_multi_touch"`       // Always true for tablets
	StylusSupport    bool     `json:"stylus_support"`       // Apple Pencil, S Pen, etc.
	StylusType       string   `json:"stylus_type"`          // Active, Passive, Magnetic

	// Performance specifications
	Processor       string `json:"processor"`
	RAM             int    `json:"ram_gb"`              // in GB (2-16GB typical)
	StorageCapacity int    `json:"storage_capacity_gb"` // in GB (32GB-2TB typical)
	StorageType     string `json:"storage_type"`        // SSD, eMMC, NVMe
	ExpandableStorage bool `json:"expandable_storage"`   // microSD support

	// Camera specifications (typically less emphasis than phones)
	RearCameras     []CameraSpec `json:"rear_cameras"`
	FrontCameras    []CameraSpec `json:"front_cameras"`
	VideoCapability string       `json:"video_capability"` // e.g., "4K@30fps"

	// Battery specifications
	BatteryCapacity       int    `json:"battery_capacity_mah"`     // 3000-10000mAh typical
	ChargingWattage       int    `json:"charging_wattage"`         // max charging speed in watts
	WirelessCharging      bool   `json:"wireless_charging"`
	WirelessChargingWatts int    `json:"wireless_charging_watts"`
	BatteryType           string `json:"battery_type"`             // Li-Ion, Li-Po
	FastChargingProtocol  string `json:"fast_charging_protocol"`   // USB-C PD, proprietary
	BatteryLifeHours      int    `json:"battery_life_hours"`       // estimated hours of use

	// Connectivity
	WiFiStandards     []string `json:"wifi_standards"`     // WiFi 5, WiFi 6, WiFi 6E
	BluetoothVersion  string   `json:"bluetooth_version"`  // 4.2, 5.0, 5.1, 5.2
	CellularSupport   bool     `json:"cellular_support"`
	NetworkSupport    []string `json:"network_support"`   // LTE, 5G
	Is5GCapable       bool     `json:"is_5g_capable"`
	USBType           string   `json:"usb_type"`           // USB-C, Lightning, micro-USB
	HasHeadphoneJack  bool     `json:"has_headphone_jack"`
	NFCSupport        bool     `json:"nfc_support"`
	GPS               bool     `json:"gps"`
	BiometricAuth     []string `json:"biometric_auth"`     // Face ID, Touch ID, fingerprint

	// Audio specifications
	Speakers     []SpeakerSpec `json:"speakers"`
	Microphones  []MicSpec     `json:"microphones"`
	AudioJack    string        `json:"audio_jack"`      // 3.5mm, USB-C, Lightning
	DolbyAtmos   bool          `json:"dolby_atmos"`
	SpatialAudio bool          `json:"spatial_audio"`

	// Ports and Expansion
	USBPorts        []USBPortSpec `json:"usb_ports"`
	SDCardSlots     int           `json:"sd_card_slots"`
	SmartCardReader bool          `json:"smart_card_reader"`

	// Tablet-specific features
	Is2In1Convertible bool     `json:"is_2in1_convertible"` // Can detach keyboard
	KeyboardIncluded  bool     `json:"keyboard_included"`
	HasKickstand      bool     `json:"has_kickstand"`
	HasPogoPins       bool     `json:"has_pogo_pins"` // For docking stations
	OperatingModes    []string `json:"operating_modes"` // Tablet, Laptop, Tent, Stand modes

	// Build and Design
	Weight         float64  `json:"weight_grams"`         // in grams
	Thickness      float64  `json:"thickness_mm"`         // in mm
	BuildMaterials []string `json:"build_materials"`      // Aluminum, Plastic, Glass
	ColorOptions   []string `json:"color_options"`
	IPRating       string   `json:"ip_rating"`            // IP68, IP67, etc.
	DustResistance bool     `json:"dust_resistance"`
	WaterResistance bool    `json:"water_resistance"`

	// Software and Ecosystem
	OperatingSystem string   `json:"operating_system"`    // iPadOS, Android, Windows
	OSVersion       string   `json:"os_version"`
	AppStore        string   `json:"app_store"`           // App Store, Google Play, Microsoft Store
	EcosystemApps   []string `json:"ecosystem_apps"`      // Pre-installed apps
	ProductivityApps []string `json:"productivity_apps"`   // Office apps, note-taking, etc.

	// Accessories compatibility
	CompatibleCases     []string `json:"compatible_cases"`
	CompatibleKeyboards []string `json:"compatible_keyboards"`
	CompatibleStyluses  []string `json:"compatible_styluses"`
	CompatibleDocks     []string `json:"compatible_docks"`
}

// CameraSpec defines camera specifications for tablets
type CameraSpec struct {
	Megapixels    float64 `json:"megapixels"`
	Aperture      string  `json:"aperture"`       // f/1.8, f/2.2, etc.
	FocalLength   int     `json:"focal_length"`  // in mm
	SensorSize    string  `json:"sensor_size"`   // 1/2.55", 1/3.1", etc.
	HasOIS        bool    `json:"has_ois"`       // Optical Image Stabilization
	HasAutofocus  bool    `json:"has_autofocus"`
	Flash         string  `json:"flash"`         // LED, None
	VideoFeatures []string `json:"video_features"` // 4K, slow-motion, etc.
}

// SpeakerSpec defines speaker specifications
type SpeakerSpec struct {
	Type        string  `json:"type"`         // Stereo, Mono, Quad
	Position    string  `json:"position"`     // Bottom, Top, Side
	Wattage     float64 `json:"wattage"`      // in watts
	FrequencyHz string  `json:"frequency_hz"` // e.g., "20-20000"
}

// MicSpec defines microphone specifications
type MicSpec struct {
	Type     string `json:"type"`      // Array, Single
	Position string `json:"position"`  // Bottom, Top
	Quality  string `json:"quality"`   // High, Medium, Low
}

// USBPortSpec defines USB port specifications
type USBPortSpec struct {
	Type    string `json:"type"`    // USB-A, USB-C
	Version string `json:"version"` // 2.0, 3.0, 3.1, 3.2
	Features []string `json:"features"` // Power Delivery, DisplayPort, etc.
}

// Validate performs comprehensive validation of tablet specifications
func (t *TabletSpec) Validate() error {
	// Screen size validation (tablets typically 7-13 inches)
	if t.ScreenSize < 7.0 || t.ScreenSize > 13.0 {
		return fmt.Errorf("tablet screen size must be between 7-13 inches, got %.1f", t.ScreenSize)
	}

	// PPI validation
	if t.PPI < 200 || t.PPI > 500 {
		return fmt.Errorf("tablet PPI must be between 200-500, got %d", t.PPI)
	}

	// RAM validation (tablets typically 2-16GB)
	if t.RAM < 2 || t.RAM > 16 {
		return fmt.Errorf("tablet RAM must be between 2-16GB, got %dGB", t.RAM)
	}

	// Storage validation (32GB-2TB)
	if t.StorageCapacity < 32 || t.StorageCapacity > 2048 {
		return fmt.Errorf("tablet storage must be between 32GB-2TB, got %dGB", t.StorageCapacity)
	}

	// Battery capacity validation (3000-10000mAh typical)
	if t.BatteryCapacity < 3000 || t.BatteryCapacity > 10000 {
		return fmt.Errorf("tablet battery capacity must be between 3000-10000mAh, got %dmAh", t.BatteryCapacity)
	}

	// Weight validation (tablets typically 300-800g)
	if t.Weight < 300 || t.Weight > 800 {
		return fmt.Errorf("tablet weight must be between 300-800g, got %.1fg", t.Weight)
	}

	// Thickness validation (tablets typically 5-10mm)
	if t.Thickness < 5.0 || t.Thickness > 10.0 {
		return fmt.Errorf("tablet thickness must be between 5-10mm, got %.1fmm", t.Thickness)
	}

	// Required fields validation
	if t.Processor == "" {
		return fmt.Errorf("processor specification is required")
	}

	if t.OperatingSystem == "" {
		return fmt.Errorf("operating system specification is required")
	}

	return nil
}

// GetCategory returns the device category
func (t *TabletSpec) GetCategory() base.CategoryType {
	return base.CategoryTablet
}

// GetDisplayInfo returns formatted display information
func (t *TabletSpec) GetDisplayInfo() string {
	return fmt.Sprintf("%.1f\" %s %s %dp", t.ScreenSize, t.ScreenType, t.Resolution, t.PPI)
}

// GetPerformanceInfo returns formatted performance information
func (t *TabletSpec) GetPerformanceInfo() string {
	return fmt.Sprintf("%s, %dGB RAM, %dGB Storage", t.Processor, t.RAM, t.StorageCapacity)
}

// GetBatteryInfo returns formatted battery information
func (t *TabletSpec) GetBatteryInfo() string {
	return fmt.Sprintf("%dmAh, %dW charging, %d hours battery life", t.BatteryCapacity, t.ChargingWattage, t.BatteryLifeHours)
}

// GetConnectivityInfo returns formatted connectivity information
func (t *TabletSpec) GetConnectivityInfo() string {
	features := []string{}

	if t.WiFiStandards != nil && len(t.WiFiStandards) > 0 {
		features = append(features, fmt.Sprintf("WiFi %s", strings.Join(t.WiFiStandards, "/")))
	}

	if t.CellularSupport {
		networkStr := "Cellular"
		if t.Is5GCapable {
			networkStr = "5G"
		}
		features = append(features, networkStr)
	}

	if t.BluetoothVersion != "" {
		features = append(features, fmt.Sprintf("Bluetooth %s", t.BluetoothVersion))
	}

	if t.NFCSupport {
		features = append(features, "NFC")
	}

	return strings.Join(features, ", ")
}

// GetBuildInfo returns formatted build information
func (t *TabletSpec) GetBuildInfo() string {
	return fmt.Sprintf("%.1fg, %.1fmm thick, %s", t.Weight, t.Thickness, strings.Join(t.BuildMaterials, "/"))
}

// SupportsFeature checks if the tablet supports a specific feature
func (t *TabletSpec) SupportsFeature(feature string) bool {
	switch strings.ToLower(feature) {
	case "cellular":
		return t.CellularSupport
	case "5g":
		return t.Is5GCapable
	case "wireless_charging":
		return t.WirelessCharging
	case "stylus":
		return t.StylusSupport
	case "nfc":
		return t.NFCSupport
	case "gps":
		return t.GPS
	case "expandable_storage":
		return t.ExpandableStorage
	case "keyboard":
		return t.KeyboardIncluded
	case "2in1":
		return t.Is2In1Convertible
	default:
		return false
	}
}
