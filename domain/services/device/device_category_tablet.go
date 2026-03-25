package device

import (
	"time"

	"smartsure/internal/domain/models"
)

// TabletDevice represents a tablet with specific features
type TabletDevice struct {
	*models.Device

	// Display specifications
	DisplaySpecs TabletDisplaySpecs `json:"display_specs"`

	// Productivity features
	ProductivityFeatures TabletProductivity `json:"productivity_features"`

	// Content consumption
	MediaFeatures MediaCapabilities `json:"media_features"`

	// Connectivity options
	ConnectivityOptions TabletConnectivity `json:"connectivity_options"`
}

// TabletDisplaySpecs represents tablet display specifications
type TabletDisplaySpecs struct {
	ScreenSize        float64 `json:"screen_size"` // in inches
	Resolution        string  `json:"resolution"`
	DisplayType       string  `json:"display_type"` // LCD, OLED, MiniLED
	RefreshRate       int     `json:"refresh_rate"`
	HasProMotion      bool    `json:"has_promotion"`
	PeakBrightness    int     `json:"peak_brightness"`
	HDRSupport        bool    `json:"hdr_support"`
	ScreenToBodyRatio float64 `json:"screen_to_body_ratio"`
	HasTrueTone       bool    `json:"has_true_tone"`
}

// TabletProductivity represents productivity features
type TabletProductivity struct {
	HasKeyboardSupport     bool     `json:"has_keyboard_support"`
	HasStylusSupport       bool     `json:"has_stylus_support"`
	StylusModel            string   `json:"stylus_model"`
	HasDesktopMode         bool     `json:"has_desktop_mode"`
	MultitaskingSupport    bool     `json:"multitasking_support"`
	ExternalDisplaySupport bool     `json:"external_display_support"`
	ProductivityApps       []string `json:"productivity_apps"`
}

// MediaCapabilities represents media consumption features
type MediaCapabilities struct {
	SpeakerCount     int      `json:"speaker_count"`
	HasDolbyAtmos    bool     `json:"has_dolby_atmos"`
	HasHeadphoneJack bool     `json:"has_headphone_jack"`
	VideoCodecs      []string `json:"video_codecs"`
	CameraCount      int      `json:"camera_count"`
	HasLiDAR         bool     `json:"has_lidar"`
	ContentServices  []string `json:"content_services"`
}

// TabletConnectivity represents connectivity options
type TabletConnectivity struct {
	HasCellular      bool     `json:"has_cellular"`
	Has5G            bool     `json:"has_5g"`
	WiFiVersion      string   `json:"wifi_version"` // WiFi 6, WiFi 6E, WiFi 7
	BluetoothVersion string   `json:"bluetooth_version"`
	USBType          string   `json:"usb_type"` // USB-C, Lightning
	HasThunderbolt   bool     `json:"has_thunderbolt"`
	SIMType          string   `json:"sim_type"` // eSIM, Nano-SIM
	ConnectorTypes   []string `json:"connector_types"`
}

// NewTabletDevice creates a new tablet device
func NewTabletDevice(base *models.Device) *TabletDevice {
	return &TabletDevice{
		Device: base,
	}
}

// GetCategory returns the device category
func (t *TabletDevice) GetCategory() string {
	return "tablet"
}

// GetSubCategory returns the device sub-category
func (t *TabletDevice) GetSubCategory() string {
	// Determine based on features and price
	if t.ProductivityFeatures.HasKeyboardSupport && t.ProductivityFeatures.HasStylusSupport {
		return "pro"
	} else if t.DisplaySpecs.ScreenSize > 11 {
		return "large"
	} else if t.Device.DeviceFinancial.PurchasePrice.Amount > 800 {
		return "premium"
	} else if t.DisplaySpecs.ScreenSize < 9 {
		return "mini"
	}
	return "standard"
}

// GetSpecificFeatures returns tablet-specific features
func (t *TabletDevice) GetSpecificFeatures() map[string]interface{} {
	return map[string]interface{}{
		"display":       t.DisplaySpecs,
		"productivity":  t.ProductivityFeatures,
		"media":         t.MediaFeatures,
		"connectivity":  t.ConnectivityOptions,
		"screen_size":   t.DisplaySpecs.ScreenSize,
		"has_stylus":    t.ProductivityFeatures.HasStylusSupport,
		"has_keyboard":  t.ProductivityFeatures.HasKeyboardSupport,
		"has_cellular":  t.ConnectivityOptions.HasCellular,
		"speaker_count": t.MediaFeatures.SpeakerCount,
	}
}

// ValidateForCategory validates tablet-specific requirements
func (t *TabletDevice) ValidateForCategory() error {
	if t.Device.DeviceClassification.Category != "tablet" {
		return ErrInvalidCategory
	}

	return nil
}

// GetCategorySpecificRisks returns tablet-specific risks
func (t *TabletDevice) GetCategorySpecificRisks() []CategoryRisk {
	risks := []CategoryRisk{}

	// Screen damage risk (larger screens more vulnerable)
	screenRisk := 1.10
	if t.DisplaySpecs.ScreenSize > 11 {
		screenRisk = 1.20
	}
	risks = append(risks, CategoryRisk{
		Type:        "screen_damage",
		Description: "Large screen damage risk",
		Severity:    "high",
		Impact:      screenRisk,
	})

	// Bend damage risk (thin tablets)
	risks = append(risks, CategoryRisk{
		Type:        "bend_damage",
		Description: "Risk of bending from pressure",
		Severity:    "medium",
		Impact:      1.08,
	})

	// Port damage risk (frequent plugging/unplugging)
	if t.ProductivityFeatures.HasKeyboardSupport {
		risks = append(risks, CategoryRisk{
			Type:        "port_damage",
			Description: "Connector wear from accessories",
			Severity:    "low",
			Impact:      1.05,
		})
	}

	// Battery degradation risk
	if t.Device.DeviceFinancial.PurchaseDate != nil {
		deviceAge := time.Since(*t.Device.DeviceFinancial.PurchaseDate)
		if deviceAge > 2*365*24*time.Hour {
			risks = append(risks, CategoryRisk{
				Type:        "battery_degradation",
				Description: "Battery capacity reduction",
				Severity:    "medium",
				Impact:      1.06,
			})
		}
	}

	// Accessory damage risk
	if t.ProductivityFeatures.HasStylusSupport {
		risks = append(risks, CategoryRisk{
			Type:        "accessory_damage",
			Description: "Stylus/keyboard damage or loss",
			Severity:    "low",
			Impact:      1.04,
		})
	}

	return risks
}

// CalculateCategoryPremiumAdjustment calculates tablet-specific premium adjustment
func (t *TabletDevice) CalculateCategoryPremiumAdjustment() float64 {
	adjustment := 1.0

	// Pro tablets are handled more carefully
	if t.GetSubCategory() == "pro" {
		adjustment *= 0.94
	}

	// Large tablets have higher repair costs
	if t.DisplaySpecs.ScreenSize > 11 {
		adjustment *= 1.12
	}

	// Mini tablets are more portable (higher risk)
	if t.GetSubCategory() == "mini" {
		adjustment *= 1.06
	}

	// Cellular models are more expensive to repair
	if t.ConnectivityOptions.HasCellular {
		adjustment *= 1.08
	}

	// Stylus support increases complexity
	if t.ProductivityFeatures.HasStylusSupport {
		adjustment *= 1.05
	}

	return adjustment
}

// GetCategoryDepreciationRate returns tablet-specific depreciation rate
func (t *TabletDevice) GetCategoryDepreciationRate() float64 {
	// Tablets depreciate moderately
	baseRate := 0.28 // 28% per year

	// Pro tablets hold value better
	if t.GetSubCategory() == "pro" {
		baseRate = 0.22
	} else if t.GetSubCategory() == "mini" {
		baseRate = 0.35
	}

	// iPads hold value better than Android tablets
	if t.Device.DeviceClassification.Brand == "Apple" {
		baseRate *= 0.85
	}

	return baseRate
}

// GetCategorySpecificCoverage returns tablet-specific coverage options
func (t *TabletDevice) GetCategorySpecificCoverage() []CoverageType {
	coverage := []CoverageType{
		{
			Type:        "screen_damage",
			Name:        "Display Protection",
			Description: "Covers screen cracks and display damage",
			MaxAmount:   600,
			Deductible:  75,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "bend_damage",
			Name:        "Structural Damage Protection",
			Description: "Covers bending and structural damage",
			MaxAmount:   500,
			Deductible:  100,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "liquid_damage",
			Name:        "Liquid Damage Protection",
			Description: "Covers spills and liquid damage",
			MaxAmount:   400,
			Deductible:  75,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "battery_replacement",
			Name:        "Battery Protection",
			Description: "Covers battery replacement",
			MaxAmount:   150,
			Deductible:  0,
			IsRequired:  false,
			IsAvailable: true,
		},
	}

	// Add stylus coverage if applicable
	if t.ProductivityFeatures.HasStylusSupport {
		coverage = append(coverage, CoverageType{
			Type:        "stylus_replacement",
			Name:        "Stylus Protection",
			Description: "Covers stylus damage or loss",
			MaxAmount:   150,
			Deductible:  25,
			IsRequired:  false,
			IsAvailable: true,
		})
	}

	// Add keyboard coverage if applicable
	if t.ProductivityFeatures.HasKeyboardSupport {
		coverage = append(coverage, CoverageType{
			Type:        "keyboard_replacement",
			Name:        "Keyboard Protection",
			Description: "Covers keyboard damage",
			MaxAmount:   200,
			Deductible:  30,
			IsRequired:  false,
			IsAvailable: true,
		})
	}

	return coverage
}

// IsEligibleForCategoryPrograms checks eligibility for tablet-specific programs
func (t *TabletDevice) IsEligibleForCategoryPrograms() map[string]bool {
	programs := make(map[string]bool)

	// Default to now if purchase date not set
	purchaseDate := time.Now()
	if t.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *t.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	// Trade-in program
	programs["trade_in"] = deviceAge < 3*365*24*time.Hour &&
		t.Device.DevicePhysicalCondition.Grade != "F"

	// Upgrade program
	programs["upgrade"] = deviceAge > 1*365*24*time.Hour &&
		t.Device.DeviceFinancial.CurrentValue.Amount > 300

	// Battery replacement program
	programs["battery_replacement"] = deviceAge > 2*365*24*time.Hour

	// Screen repair program
	programs["screen_repair"] = true

	// Productivity bundle (keyboard + stylus)
	programs["productivity_bundle"] = t.GetSubCategory() == "pro" ||
		t.GetSubCategory() == "premium"

	// Student discount program
	programs["student_program"] = true

	// Creative professional program
	programs["creative_pro"] = t.ProductivityFeatures.HasStylusSupport &&
		t.DisplaySpecs.ScreenSize > 10

	return programs
}

// GetCategoryMaintenanceSchedule returns tablet-specific maintenance schedule
func (t *TabletDevice) GetCategoryMaintenanceSchedule() MaintenanceSchedule {
	now := time.Now()

	// Default to now if purchase date not set
	purchaseDate := now
	if t.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *t.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	intervals := []MaintenanceInterval{
		{
			Type:        "software_update",
			Description: "OS and app updates",
			Frequency:   30 * 24 * time.Hour, // Monthly
			IsCritical:  true,
			Cost:        0,
		},
		{
			Type:        "screen_cleaning",
			Description: "Professional screen cleaning",
			Frequency:   90 * 24 * time.Hour, // Quarterly
			IsCritical:  false,
			Cost:        0,
		},
		{
			Type:        "port_cleaning",
			Description: "Clean charging and accessory ports",
			Frequency:   180 * 24 * time.Hour, // Semi-annually
			IsCritical:  false,
			Cost:        0,
		},
		{
			Type:        "battery_calibration",
			Description: "Battery health check",
			Frequency:   365 * 24 * time.Hour, // Annually
			IsCritical:  false,
			Cost:        0,
		},
	}

	// Add accessory maintenance if applicable
	if t.ProductivityFeatures.HasKeyboardSupport {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "keyboard_check",
			Description: "Keyboard connection and key check",
			Frequency:   180 * 24 * time.Hour, // Semi-annually
			IsCritical:  false,
			Cost:        0,
		})
	}

	if t.ProductivityFeatures.HasStylusSupport {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "stylus_tip_replacement",
			Description: "Replace stylus tip if worn",
			Frequency:   365 * 24 * time.Hour, // Annually
			IsCritical:  false,
			Cost:        15,
		})
	}

	// Add battery replacement after 3 years
	if deviceAge > 3*365*24*time.Hour {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "battery_replacement",
			Description: "Consider battery replacement",
			Frequency:   0, // One-time
			IsCritical:  true,
			Cost:        120,
		})
	}

	return MaintenanceSchedule{
		Intervals: intervals,
		NextDue:   now.Add(30 * 24 * time.Hour),
		LastDone:  now,
	}
}
