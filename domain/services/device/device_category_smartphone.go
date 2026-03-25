package device

import (
	"time"

	"smartsure/internal/domain/models"
)

// SmartphoneDevice represents a smartphone with specific features
type SmartphoneDevice struct {
	*models.Device

	// Camera specifications
	CameraSpecs CameraSpecification `json:"camera_specs"`

	// Network capabilities
	NetworkCapabilities NetworkCapabilities `json:"network_capabilities"`

	// Gaming features
	GamingFeatures GamingFeatures `json:"gaming_features"`

	// Display features
	DisplayFeatures DisplayFeatures `json:"display_features"`
}

// CameraSpecification represents camera specs
type CameraSpecification struct {
	MainCameraMegapixels  float64  `json:"main_camera_megapixels"`
	FrontCameraMegapixels float64  `json:"front_camera_megapixels"`
	NumberOfCameras       int      `json:"number_of_cameras"`
	HasOpticalZoom        bool     `json:"has_optical_zoom"`
	OpticalZoomLevel      float64  `json:"optical_zoom_level"`
	HasNightMode          bool     `json:"has_night_mode"`
	VideoResolution       string   `json:"video_resolution"` // 4K, 8K, etc.
	HasStabilization      bool     `json:"has_stabilization"`
	CameraFeatures        []string `json:"camera_features"`
}

// NetworkCapabilities represents network features
type NetworkCapabilities struct {
	Has5G            bool     `json:"has_5g"`
	Supported5GBands []string `json:"supported_5g_bands"`
	Has4G            bool     `json:"has_4g"`
	DualSIM          bool     `json:"dual_sim"`
	TripleSIM        bool     `json:"triple_sim"`
	eSIMSupport      bool     `json:"esim_support"`
	WiFi6Support     bool     `json:"wifi6_support"`
	WiFi6ESupport    bool     `json:"wifi6e_support"`
}

// GamingFeatures represents gaming capabilities
type GamingFeatures struct {
	HasGamingMode     bool   `json:"has_gaming_mode"`
	RefreshRate       int    `json:"refresh_rate"` // 60Hz, 90Hz, 120Hz, 144Hz
	TouchSamplingRate int    `json:"touch_sampling_rate"`
	HasVaporCooling   bool   `json:"has_vapor_cooling"`
	GameOptimization  bool   `json:"game_optimization"`
	GPUModel          string `json:"gpu_model"`
	BenchmarkScore    int    `json:"benchmark_score"`
}

// DisplayFeatures represents display specifications
type DisplayFeatures struct {
	ScreenSize        float64 `json:"screen_size"`  // in inches
	Resolution        string  `json:"resolution"`   // e.g., "2400x1080"
	DisplayType       string  `json:"display_type"` // AMOLED, LCD, etc.
	RefreshRate       int     `json:"refresh_rate"`
	PeakBrightness    int     `json:"peak_brightness"` // in nits
	HDRSupport        bool    `json:"hdr_support"`
	AlwaysOnDisplay   bool    `json:"always_on_display"`
	CurvedDisplay     bool    `json:"curved_display"`
	ScreenToBodyRatio float64 `json:"screen_to_body_ratio"`
}

// NewSmartphoneDevice creates a new smartphone device
func NewSmartphoneDevice(base *models.Device) *SmartphoneDevice {
	return &SmartphoneDevice{
		Device: base,
	}
}

// GetCategory returns the device category
func (s *SmartphoneDevice) GetCategory() string {
	return "smartphone"
}

// GetSubCategory returns the device sub-category
func (s *SmartphoneDevice) GetSubCategory() string {
	// Determine based on price/features
	if s.Device.DeviceFinancial.PurchasePrice.Amount > 1000 {
		return "flagship"
	} else if s.Device.DeviceFinancial.PurchasePrice.Amount > 500 {
		return "premium"
	} else if s.Device.DeviceFinancial.PurchasePrice.Amount > 300 {
		return "midrange"
	}
	return "budget"
}

// GetSpecificFeatures returns smartphone-specific features
func (s *SmartphoneDevice) GetSpecificFeatures() map[string]interface{} {
	return map[string]interface{}{
		"camera":       s.CameraSpecs,
		"network":      s.NetworkCapabilities,
		"gaming":       s.GamingFeatures,
		"display":      s.DisplayFeatures,
		"has_5g":       s.NetworkCapabilities.Has5G,
		"dual_sim":     s.NetworkCapabilities.DualSIM,
		"refresh_rate": s.DisplayFeatures.RefreshRate,
		"camera_count": s.CameraSpecs.NumberOfCameras,
	}
}

// ValidateForCategory validates smartphone-specific requirements
func (s *SmartphoneDevice) ValidateForCategory() error {
	// Smartphone-specific validation
	if s.Device.DeviceIdentification.IMEI == "" {
		return ErrInvalidIMEI
	}

	if s.Device.DeviceClassification.Category != "smartphone" {
		return ErrInvalidCategory
	}

	return nil
}

// GetCategorySpecificRisks returns smartphone-specific risks
func (s *SmartphoneDevice) GetCategorySpecificRisks() []CategoryRisk {
	risks := []CategoryRisk{}

	// Screen damage risk (most common for smartphones)
	risks = append(risks, CategoryRisk{
		Type:        "screen_damage",
		Description: "Risk of screen damage from drops",
		Severity:    "high",
		Impact:      1.15, // 15% premium increase
	})

	// Water damage risk
	waterResistance := s.Device.DeviceSpecifications.WaterResistance
	if waterResistance == "" || waterResistance == "none" {
		risks = append(risks, CategoryRisk{
			Type:        "water_damage",
			Description: "No water resistance rating",
			Severity:    "medium",
			Impact:      1.10,
		})
	}

	// Theft risk based on value
	if s.Device.DeviceFinancial.CurrentValue.Amount > 1000 {
		risks = append(risks, CategoryRisk{
			Type:        "theft",
			Description: "High-value device theft risk",
			Severity:    "high",
			Impact:      1.20,
		})
	}

	// Battery degradation risk
	if s.Device.DeviceFinancial.PurchaseDate != nil {
		deviceAge := time.Since(*s.Device.DeviceFinancial.PurchaseDate)
		if deviceAge > 2*365*24*time.Hour {
			risks = append(risks, CategoryRisk{
				Type:        "battery_degradation",
				Description: "Battery degradation due to age",
				Severity:    "medium",
				Impact:      1.05,
			})
		}
	}

	// 5G-specific risks
	if s.NetworkCapabilities.Has5G {
		risks = append(risks, CategoryRisk{
			Type:        "overheating",
			Description: "5G usage can cause overheating",
			Severity:    "low",
			Impact:      1.02,
		})
	}

	return risks
}

// CalculateCategoryPremiumAdjustment calculates smartphone-specific premium adjustment
func (s *SmartphoneDevice) CalculateCategoryPremiumAdjustment() float64 {
	adjustment := 1.0

	// Flagship devices get better care typically
	if s.GetSubCategory() == "flagship" {
		adjustment *= 0.95
	}

	// Multiple cameras increase repair cost
	if s.CameraSpecs.NumberOfCameras > 3 {
		adjustment *= 1.10
	}

	// High refresh rate screens are more expensive
	if s.DisplayFeatures.RefreshRate > 90 {
		adjustment *= 1.05
	}

	// Gaming phones have higher usage
	if s.GamingFeatures.HasGamingMode {
		adjustment *= 1.08
	}

	// Dual SIM increases complexity
	if s.NetworkCapabilities.DualSIM {
		adjustment *= 1.03
	}

	return adjustment
}

// GetCategoryDepreciationRate returns smartphone-specific depreciation rate
func (s *SmartphoneDevice) GetCategoryDepreciationRate() float64 {
	// Smartphones depreciate faster than other devices
	baseRate := 0.30 // 30% per year

	// Flagship devices hold value better
	if s.GetSubCategory() == "flagship" {
		baseRate = 0.25
	} else if s.GetSubCategory() == "budget" {
		baseRate = 0.40
	}

	// Android devices depreciate faster than iOS
	if s.Device.DeviceSpecifications.OperatingSystem == "Android" {
		baseRate *= 1.1
	}

	return baseRate
}

// GetCategorySpecificCoverage returns smartphone-specific coverage options
func (s *SmartphoneDevice) GetCategorySpecificCoverage() []CoverageType {
	return []CoverageType{
		{
			Type:        "screen_damage",
			Name:        "Screen Protection",
			Description: "Covers screen cracks and touch issues",
			MaxAmount:   500,
			Deductible:  50,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "water_damage",
			Name:        "Liquid Damage Protection",
			Description: "Covers water and liquid damage",
			MaxAmount:   800,
			Deductible:  75,
			IsRequired:  false,
			IsAvailable: s.Device.DeviceSpecifications.WaterResistance == "",
		},
		{
			Type:        "theft",
			Name:        "Theft Protection",
			Description: "Covers device theft with police report",
			MaxAmount:   s.Device.DeviceFinancial.CurrentValue.Amount,
			Deductible:  100,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "battery_replacement",
			Name:        "Battery Protection",
			Description: "Covers battery replacement when health drops below 80%",
			MaxAmount:   150,
			Deductible:  0,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "camera_damage",
			Name:        "Camera Protection",
			Description: "Covers camera lens and sensor damage",
			MaxAmount:   300,
			Deductible:  50,
			IsRequired:  false,
			IsAvailable: s.CameraSpecs.NumberOfCameras > 1,
		},
	}
}

// IsEligibleForCategoryPrograms checks eligibility for smartphone-specific programs
func (s *SmartphoneDevice) IsEligibleForCategoryPrograms() map[string]bool {
	programs := make(map[string]bool)

	// Default to now if purchase date not set
	purchaseDate := time.Now()
	if s.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *s.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	// Trade-in program
	programs["trade_in"] = deviceAge < 3*365*24*time.Hour &&
		s.Device.DevicePhysicalCondition.Grade != "F"

	// Upgrade program
	programs["upgrade"] = deviceAge > 1*365*24*time.Hour &&
		s.Device.DeviceFinancial.CurrentValue.Amount > 200

	// Battery replacement program
	programs["battery_replacement"] = deviceAge > 18*30*24*time.Hour // 18 months

	// Screen repair program
	programs["screen_repair"] = true // Always available for smartphones

	// 5G upgrade program
	programs["5g_upgrade"] = !s.NetworkCapabilities.Has5G &&
		s.GetSubCategory() != "budget"

	// Premium care program
	programs["premium_care"] = s.GetSubCategory() == "flagship" ||
		s.GetSubCategory() == "premium"

	return programs
}

// GetCategoryMaintenanceSchedule returns smartphone-specific maintenance schedule
func (s *SmartphoneDevice) GetCategoryMaintenanceSchedule() MaintenanceSchedule {
	now := time.Now()

	// Default to now if purchase date not set
	purchaseDate := now
	if s.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *s.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	intervals := []MaintenanceInterval{
		{
			Type:        "software_update",
			Description: "Security and OS updates",
			Frequency:   30 * 24 * time.Hour, // Monthly
			IsCritical:  true,
			Cost:        0,
		},
		{
			Type:        "cache_cleanup",
			Description: "Clear cache and optimize storage",
			Frequency:   90 * 24 * time.Hour, // Quarterly
			IsCritical:  false,
			Cost:        0,
		},
		{
			Type:        "battery_calibration",
			Description: "Battery health check and calibration",
			Frequency:   180 * 24 * time.Hour, // Semi-annually
			IsCritical:  false,
			Cost:        0,
		},
		{
			Type:        "screen_protector_replacement",
			Description: "Replace screen protector if damaged",
			Frequency:   365 * 24 * time.Hour, // Annually
			IsCritical:  false,
			Cost:        20,
		},
		{
			Type:        "port_cleaning",
			Description: "Clean charging port and speakers",
			Frequency:   180 * 24 * time.Hour, // Semi-annually
			IsCritical:  false,
			Cost:        0,
		},
	}

	// Add battery replacement after 2 years
	if deviceAge > 2*365*24*time.Hour {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "battery_replacement",
			Description: "Consider battery replacement",
			Frequency:   0, // One-time
			IsCritical:  true,
			Cost:        100,
		})
	}

	return MaintenanceSchedule{
		Intervals: intervals,
		NextDue:   now.Add(30 * 24 * time.Hour),
		LastDone:  now,
	}
}
