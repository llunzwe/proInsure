package device

import (
	"time"

	"smartsure/internal/domain/models"
)

// LaptopDevice represents a laptop with specific features
type LaptopDevice struct {
	*models.Device

	// Processor specifications
	ProcessorSpecs ProcessorSpecification `json:"processor_specs"`

	// Memory and storage
	SystemSpecs SystemSpecification `json:"system_specs"`

	// Display configuration
	LaptopDisplay LaptopDisplaySpecs `json:"display_specs"`

	// Port configuration
	PortConfiguration PortSpecs `json:"port_configuration"`

	// Business features
	BusinessFeatures LaptopBusinessFeatures `json:"business_features"`
}

// ProcessorSpecification represents CPU specifications
type ProcessorSpecification struct {
	Brand           string  `json:"brand"` // Intel, AMD, Apple
	Model           string  `json:"model"`
	Generation      int     `json:"generation"`
	CoreCount       int     `json:"core_count"`
	ThreadCount     int     `json:"thread_count"`
	BaseClockSpeed  float64 `json:"base_clock_ghz"`
	TurboClockSpeed float64 `json:"turbo_clock_ghz"`
	HasGPU          bool    `json:"has_integrated_gpu"`
	GPUModel        string  `json:"gpu_model"`
	TDP             int     `json:"tdp_watts"`
}

// SystemSpecification represents memory and storage
type SystemSpecification struct {
	RAMSize            int    `json:"ram_size_gb"`
	RAMType            string `json:"ram_type"` // DDR4, DDR5, LPDDR4X
	RAMSpeed           int    `json:"ram_speed_mhz"`
	MaxRAM             int    `json:"max_ram_gb"`
	StorageType        string `json:"storage_type"` // SSD, HDD, Hybrid
	StorageSize        int    `json:"storage_size_gb"`
	StorageSpeed       string `json:"storage_speed"` // NVMe, SATA
	HasSecondaryDrive  bool   `json:"has_secondary_drive"`
	SecondaryDriveSize int    `json:"secondary_drive_size_gb"`
	GraphicsCard       string `json:"graphics_card"`
	DedicatedVRAM      int    `json:"dedicated_vram_gb"`
}

// LaptopDisplaySpecs represents display specifications
type LaptopDisplaySpecs struct {
	ScreenSize       float64 `json:"screen_size"`  // in inches
	Resolution       string  `json:"resolution"`   // e.g., "1920x1080", "3840x2160"
	DisplayType      string  `json:"display_type"` // IPS, OLED, TN
	RefreshRate      int     `json:"refresh_rate"`
	Touchscreen      bool    `json:"touchscreen"`
	ColorGamut       string  `json:"color_gamut"` // sRGB%, AdobeRGB%
	PeakBrightness   int     `json:"peak_brightness"`
	HasPrivacyScreen bool    `json:"has_privacy_screen"`
	AspectRatio      string  `json:"aspect_ratio"` // 16:9, 16:10, 3:2
}

// PortSpecs represents port configuration
type PortSpecs struct {
	USBTypeAPorts   int      `json:"usb_type_a_ports"`
	USBTypeCPorts   int      `json:"usb_type_c_ports"`
	Thunderbolt     int      `json:"thunderbolt_ports"`
	HDMI            bool     `json:"has_hdmi"`
	DisplayPort     bool     `json:"has_displayport"`
	AudioJack       bool     `json:"has_audio_jack"`
	SDCardReader    bool     `json:"has_sd_card_reader"`
	Ethernet        bool     `json:"has_ethernet"`
	ProprietaryPort bool     `json:"has_proprietary_port"`
	PortTypes       []string `json:"port_types"`
}

// LaptopBusinessFeatures represents business-oriented features
type LaptopBusinessFeatures struct {
	Category            string   `json:"category"` // business, gaming, consumer, workstation
	HasFingerprint      bool     `json:"has_fingerprint"`
	HasSmartCard        bool     `json:"has_smart_card_reader"`
	HasTPM              bool     `json:"has_tpm"`
	HasIRCamera         bool     `json:"has_ir_camera"`
	HasBacklitKeyboard  bool     `json:"has_backlit_keyboard"`
	KeyboardType        string   `json:"keyboard_type"` // standard, mechanical, butterfly
	HasNumPad           bool     `json:"has_numpad"`
	MilSpecRated        bool     `json:"mil_spec_rated"`
	SecurityFeatures    []string `json:"security_features"`
	BatteryLife         int      `json:"battery_life_hours"`
	HasRemovableBattery bool     `json:"has_removable_battery"`
	Weight              float64  `json:"weight_kg"`
	Thickness           float64  `json:"thickness_mm"`
}

// NewLaptopDevice creates a new laptop device
func NewLaptopDevice(base *models.Device) *LaptopDevice {
	return &LaptopDevice{
		Device: base,
	}
}

// GetCategory returns the device category
func (l *LaptopDevice) GetCategory() string {
	return "laptop"
}

// GetSubCategory returns the device sub-category
func (l *LaptopDevice) GetSubCategory() string {
	// Determine based on features and specs
	if l.SystemSpecs.DedicatedVRAM >= 4 && l.ProcessorSpecs.CoreCount >= 6 {
		return "gaming"
	} else if l.BusinessFeatures.HasTPM && l.BusinessFeatures.HasFingerprint {
		return "business"
	} else if l.ProcessorSpecs.CoreCount >= 8 && l.SystemSpecs.RAMSize >= 32 {
		return "workstation"
	} else if l.LaptopDisplay.Touchscreen && l.BusinessFeatures.Weight < 1.5 {
		return "ultrabook"
	} else if l.Device.DeviceFinancial.PurchasePrice.Amount > 1500 {
		return "premium"
	}
	return "consumer"
}

// GetSpecificFeatures returns laptop-specific features
func (l *LaptopDevice) GetSpecificFeatures() map[string]interface{} {
	return map[string]interface{}{
		"processor":         l.ProcessorSpecs,
		"system":            l.SystemSpecs,
		"display":           l.LaptopDisplay,
		"ports":             l.PortConfiguration,
		"business":          l.BusinessFeatures,
		"cpu_model":         l.ProcessorSpecs.Model,
		"ram_size":          l.SystemSpecs.RAMSize,
		"storage_type":      l.SystemSpecs.StorageType,
		"screen_size":       l.LaptopDisplay.ScreenSize,
		"has_touchscreen":   l.LaptopDisplay.Touchscreen,
		"has_dedicated_gpu": l.SystemSpecs.DedicatedVRAM > 0,
		"category":          l.BusinessFeatures.Category,
		"weight_kg":         l.BusinessFeatures.Weight,
	}
}

// ValidateForCategory validates laptop-specific requirements
func (l *LaptopDevice) ValidateForCategory() error {
	if l.Device.DeviceClassification.Category != "laptop" &&
		l.Device.DeviceClassification.Category != "notebook" {
		return ErrInvalidCategory
	}

	return nil
}

// GetCategorySpecificRisks returns laptop-specific risks
func (l *LaptopDevice) GetCategorySpecificRisks() []CategoryRisk {
	risks := []CategoryRisk{}

	// Screen damage risk
	risks = append(risks, CategoryRisk{
		Type:        "screen_damage",
		Description: "Screen damage from closure pressure or impact",
		Severity:    "high",
		Impact:      1.12,
	})

	// Hinge damage risk
	risks = append(risks, CategoryRisk{
		Type:        "hinge_damage",
		Description: "Display hinge wear and failure",
		Severity:    "medium",
		Impact:      1.10,
	})

	// Liquid damage risk (keyboard spills)
	risks = append(risks, CategoryRisk{
		Type:        "liquid_damage",
		Description: "Keyboard liquid spills",
		Severity:    "high",
		Impact:      1.15,
	})

	// Overheating risk (gaming/workstation laptops)
	if l.GetSubCategory() == "gaming" || l.GetSubCategory() == "workstation" {
		risks = append(risks, CategoryRisk{
			Type:        "overheating",
			Description: "Thermal stress from high performance usage",
			Severity:    "medium",
			Impact:      1.08,
		})
	}

	// Port wear risk
	if l.PortConfiguration.USBTypeAPorts > 2 || l.PortConfiguration.USBTypeCPorts > 2 {
		risks = append(risks, CategoryRisk{
			Type:        "port_wear",
			Description: "Port damage from frequent plugging/unplugging",
			Severity:    "low",
			Impact:      1.05,
		})
	}

	// Battery degradation risk
	if l.Device.DeviceFinancial.PurchaseDate != nil {
		deviceAge := time.Since(*l.Device.DeviceFinancial.PurchaseDate)
		if deviceAge > 2*365*24*time.Hour {
			risks = append(risks, CategoryRisk{
				Type:        "battery_degradation",
				Description: "Battery capacity reduction",
				Severity:    "medium",
				Impact:      1.07,
			})
		}
	}

	// Keyboard wear risk
	risks = append(risks, CategoryRisk{
		Type:        "keyboard_wear",
		Description: "Key failure from extensive use",
		Severity:    "low",
		Impact:      1.04,
	})

	// Drop damage risk (ultrabooks/thin laptops)
	if l.BusinessFeatures.Weight < 1.5 {
		risks = append(risks, CategoryRisk{
			Type:        "drop_damage",
			Description: "Higher drop damage risk for ultralight devices",
			Severity:    "medium",
			Impact:      1.09,
		})
	}

	return risks
}

// CalculateCategoryPremiumAdjustment calculates laptop-specific premium adjustment
func (l *LaptopDevice) CalculateCategoryPremiumAdjustment() float64 {
	adjustment := 1.0

	// Business laptops are handled more carefully
	if l.GetSubCategory() == "business" {
		adjustment *= 0.93
	}

	// Gaming laptops face more stress
	if l.GetSubCategory() == "gaming" {
		adjustment *= 1.18
	}

	// Workstations are expensive to repair
	if l.GetSubCategory() == "workstation" {
		adjustment *= 1.15
	}

	// Ultrabooks are more fragile
	if l.GetSubCategory() == "ultrabook" {
		adjustment *= 1.10
	}

	// Touchscreen adds complexity
	if l.LaptopDisplay.Touchscreen {
		adjustment *= 1.07
	}

	// Dedicated GPU increases repair cost
	if l.SystemSpecs.DedicatedVRAM > 0 {
		adjustment *= 1.08
	}

	// MIL-spec rated devices are more durable
	if l.BusinessFeatures.MilSpecRated {
		adjustment *= 0.90
	}

	// Removable battery reduces some risks
	if l.BusinessFeatures.HasRemovableBattery {
		adjustment *= 0.97
	}

	return adjustment
}

// GetCategoryDepreciationRate returns laptop-specific depreciation rate
func (l *LaptopDevice) GetCategoryDepreciationRate() float64 {
	// Laptops depreciate moderately fast
	baseRate := 0.32 // 32% per year

	// Business laptops hold value better
	if l.GetSubCategory() == "business" {
		baseRate = 0.25
	} else if l.GetSubCategory() == "gaming" {
		baseRate = 0.40 // Gaming laptops depreciate faster
	} else if l.GetSubCategory() == "workstation" {
		baseRate = 0.28
	} else if l.GetSubCategory() == "consumer" {
		baseRate = 0.38
	}

	// Apple laptops hold value better
	if l.Device.DeviceClassification.Brand == "Apple" {
		baseRate *= 0.80
	}

	// High-end processors depreciate slower
	if l.ProcessorSpecs.CoreCount >= 8 {
		baseRate *= 0.95
	}

	return baseRate
}

// GetCategorySpecificCoverage returns laptop-specific coverage options
func (l *LaptopDevice) GetCategorySpecificCoverage() []CoverageType {
	coverage := []CoverageType{
		{
			Type:        "screen_damage",
			Name:        "Display Protection",
			Description: "Covers screen cracks and display damage",
			MaxAmount:   700,
			Deductible:  100,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "liquid_damage",
			Name:        "Liquid Damage Protection",
			Description: "Covers spills and liquid damage to keyboard/internals",
			MaxAmount:   600,
			Deductible:  75,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "accidental_damage",
			Name:        "Accidental Damage Protection",
			Description: "Covers drops, impacts, and physical damage",
			MaxAmount:   800,
			Deductible:  100,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "hinge_repair",
			Name:        "Hinge Protection",
			Description: "Covers display hinge repair or replacement",
			MaxAmount:   250,
			Deductible:  50,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "keyboard_replacement",
			Name:        "Keyboard Protection",
			Description: "Covers keyboard failure or damage",
			MaxAmount:   200,
			Deductible:  40,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "battery_replacement",
			Name:        "Battery Protection",
			Description: "Covers battery replacement when capacity drops below 80%",
			MaxAmount:   180,
			Deductible:  0,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "port_repair",
			Name:        "Port Protection",
			Description: "Covers USB, charging port, and connector repairs",
			MaxAmount:   150,
			Deductible:  30,
			IsRequired:  false,
			IsAvailable: l.PortConfiguration.USBTypeCPorts > 0 || l.PortConfiguration.Thunderbolt > 0,
		},
	}

	// Add overheating protection for gaming/workstation laptops
	if l.GetSubCategory() == "gaming" || l.GetSubCategory() == "workstation" {
		coverage = append(coverage, CoverageType{
			Type:        "thermal_damage",
			Name:        "Thermal Protection",
			Description: "Covers overheating and thermal damage",
			MaxAmount:   400,
			Deductible:  75,
			IsRequired:  false,
			IsAvailable: true,
		})
	}

	// Add data recovery for business laptops
	if l.GetSubCategory() == "business" || l.GetSubCategory() == "workstation" {
		coverage = append(coverage, CoverageType{
			Type:        "data_recovery",
			Name:        "Data Recovery Protection",
			Description: "Professional data recovery services",
			MaxAmount:   500,
			Deductible:  100,
			IsRequired:  false,
			IsAvailable: true,
		})
	}

	return coverage
}

// IsEligibleForCategoryPrograms checks eligibility for laptop-specific programs
func (l *LaptopDevice) IsEligibleForCategoryPrograms() map[string]bool {
	programs := make(map[string]bool)

	// Default to now if purchase date not set
	purchaseDate := time.Now()
	if l.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *l.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	// Trade-in program
	programs["trade_in"] = deviceAge < 4*365*24*time.Hour &&
		l.Device.DevicePhysicalCondition.Grade != "F"

	// Upgrade program
	programs["upgrade"] = deviceAge > 2*365*24*time.Hour &&
		l.Device.DeviceFinancial.CurrentValue.Amount > 400

	// Battery replacement program
	programs["battery_replacement"] = deviceAge > 2*365*24*time.Hour

	// Screen repair program
	programs["screen_repair"] = true

	// Keyboard replacement program
	programs["keyboard_replacement"] = true

	// Business refresh program
	programs["business_refresh"] = l.GetSubCategory() == "business" &&
		deviceAge > 3*365*24*time.Hour

	// Gaming upgrade program
	programs["gaming_upgrade"] = l.GetSubCategory() == "gaming" &&
		l.SystemSpecs.DedicatedVRAM < 8

	// Student discount program
	programs["student_program"] = l.GetSubCategory() == "consumer" ||
		l.GetSubCategory() == "ultrabook"

	// Professional warranty extension
	programs["pro_warranty"] = l.GetSubCategory() == "workstation" ||
		l.GetSubCategory() == "business"

	// SSD upgrade program (for HDD laptops)
	programs["ssd_upgrade"] = l.SystemSpecs.StorageType == "HDD"

	// RAM upgrade program
	programs["ram_upgrade"] = l.SystemSpecs.RAMSize < l.SystemSpecs.MaxRAM

	return programs
}

// GetCategoryMaintenanceSchedule returns laptop-specific maintenance schedule
func (l *LaptopDevice) GetCategoryMaintenanceSchedule() MaintenanceSchedule {
	now := time.Now()

	// Default to now if purchase date not set
	purchaseDate := now
	if l.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *l.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	intervals := []MaintenanceInterval{
		{
			Type:        "software_update",
			Description: "OS and driver updates",
			Frequency:   30 * 24 * time.Hour, // Monthly
			IsCritical:  true,
			Cost:        0,
		},
		{
			Type:        "thermal_cleaning",
			Description: "Clean fans and thermal vents",
			Frequency:   180 * 24 * time.Hour, // Semi-annually
			IsCritical:  true,
			Cost:        30,
		},
		{
			Type:        "keyboard_cleaning",
			Description: "Keyboard and trackpad cleaning",
			Frequency:   90 * 24 * time.Hour, // Quarterly
			IsCritical:  false,
			Cost:        0,
		},
		{
			Type:        "hinge_inspection",
			Description: "Check display hinges for wear",
			Frequency:   365 * 24 * time.Hour, // Annually
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
			Type:        "port_inspection",
			Description: "Check all ports for damage or debris",
			Frequency:   180 * 24 * time.Hour, // Semi-annually
			IsCritical:  false,
			Cost:        0,
		},
	}

	// Add thermal paste replacement for gaming/workstation laptops
	if l.GetSubCategory() == "gaming" || l.GetSubCategory() == "workstation" {
		if deviceAge > 2*365*24*time.Hour {
			intervals = append(intervals, MaintenanceInterval{
				Type:        "thermal_paste",
				Description: "Replace thermal paste for CPU/GPU",
				Frequency:   0, // One-time
				IsCritical:  true,
				Cost:        80,
			})
		}
	}

	// Add battery replacement after 3 years
	if deviceAge > 3*365*24*time.Hour {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "battery_replacement",
			Description: "Consider battery replacement",
			Frequency:   0, // One-time
			IsCritical:  true,
			Cost:        150,
		})
	}

	// Add SSD health check for laptops with SSDs
	if l.SystemSpecs.StorageType == "SSD" {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "ssd_health_check",
			Description: "Check SSD health and wear level",
			Frequency:   365 * 24 * time.Hour, // Annually
			IsCritical:  false,
			Cost:        0,
		})
	}

	return MaintenanceSchedule{
		Intervals: intervals,
		NextDue:   now.Add(30 * 24 * time.Hour),
		LastDone:  now,
	}
}
