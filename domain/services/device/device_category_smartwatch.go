package device

import (
	"time"

	"smartsure/internal/domain/models"
)

// SmartwatchDevice represents a smartwatch with specific features
type SmartwatchDevice struct {
	*models.Device

	// Health sensors
	HealthSensors HealthSensorSpecs `json:"health_sensors"`

	// Fitness features
	FitnessFeatures FitnessCapabilities `json:"fitness_features"`

	// Battery features
	BatteryFeatures BatterySpecs `json:"battery_features"`

	// Connectivity
	ConnectivityFeatures WatchConnectivity `json:"connectivity_features"`
}

// HealthSensorSpecs represents health sensor specifications
type HealthSensorSpecs struct {
	HasHeartRateMonitor  bool     `json:"has_heart_rate_monitor"`
	HasECG               bool     `json:"has_ecg"`
	HasBloodOxygen       bool     `json:"has_blood_oxygen"`
	HasBloodPressure     bool     `json:"has_blood_pressure"`
	HasTemperatureSensor bool     `json:"has_temperature_sensor"`
	HasSleepTracking     bool     `json:"has_sleep_tracking"`
	HasStressMonitoring  bool     `json:"has_stress_monitoring"`
	HealthFeatures       []string `json:"health_features"`
}

// FitnessCapabilities represents fitness tracking features
type FitnessCapabilities struct {
	HasGPS                  bool     `json:"has_gps"`
	HasAltimeter            bool     `json:"has_altimeter"`
	HasBarometer            bool     `json:"has_barometer"`
	HasCompass              bool     `json:"has_compass"`
	SportModes              int      `json:"sport_modes_count"`
	WaterResistanceATM      int      `json:"water_resistance_atm"`
	HasSwimmingTracking     bool     `json:"has_swimming_tracking"`
	HasAutoWorkoutDetection bool     `json:"has_auto_workout_detection"`
	FitnessFeatures         []string `json:"fitness_features"`
}

// BatterySpecs represents battery specifications
type BatterySpecs struct {
	BatteryLife         int  `json:"battery_life_days"`
	ChargingTime        int  `json:"charging_time_minutes"`
	HasWirelessCharging bool `json:"has_wireless_charging"`
	HasFastCharging     bool `json:"has_fast_charging"`
	PowerSavingMode     bool `json:"has_power_saving_mode"`
	BatteryCapacity     int  `json:"battery_capacity_mah"`
}

// WatchConnectivity represents connectivity features
type WatchConnectivity struct {
	HasCellular      bool   `json:"has_cellular"`
	HasWiFi          bool   `json:"has_wifi"`
	BluetoothVersion string `json:"bluetooth_version"`
	HasNFC           bool   `json:"has_nfc"`
	HaseSIM          bool   `json:"has_esim"`
	CompatibleOS     string `json:"compatible_os"` // iOS, Android, Both
}

// NewSmartwatchDevice creates a new smartwatch device
func NewSmartwatchDevice(base *models.Device) *SmartwatchDevice {
	return &SmartwatchDevice{
		Device: base,
	}
}

// GetCategory returns the device category
func (w *SmartwatchDevice) GetCategory() string {
	return "smartwatch"
}

// GetSubCategory returns the device sub-category
func (w *SmartwatchDevice) GetSubCategory() string {
	// Determine based on features
	if w.HealthSensors.HasECG && w.HealthSensors.HasBloodOxygen {
		return "health_focused"
	} else if w.FitnessFeatures.SportModes > 50 {
		return "sport_focused"
	} else if w.Device.DeviceFinancial.PurchasePrice.Amount > 500 {
		return "premium"
	}
	return "basic"
}

// GetSpecificFeatures returns smartwatch-specific features
func (w *SmartwatchDevice) GetSpecificFeatures() map[string]interface{} {
	return map[string]interface{}{
		"health_sensors":   w.HealthSensors,
		"fitness":          w.FitnessFeatures,
		"battery":          w.BatteryFeatures,
		"connectivity":     w.ConnectivityFeatures,
		"has_ecg":          w.HealthSensors.HasECG,
		"has_gps":          w.FitnessFeatures.HasGPS,
		"has_cellular":     w.ConnectivityFeatures.HasCellular,
		"water_resistance": w.FitnessFeatures.WaterResistanceATM,
		"battery_days":     w.BatteryFeatures.BatteryLife,
	}
}

// ValidateForCategory validates smartwatch-specific requirements
func (w *SmartwatchDevice) ValidateForCategory() error {
	if w.Device.DeviceClassification.Category != "smartwatch" &&
		w.Device.DeviceClassification.Category != "wearable" {
		return ErrInvalidCategory
	}

	return nil
}

// GetCategorySpecificRisks returns smartwatch-specific risks
func (w *SmartwatchDevice) GetCategorySpecificRisks() []CategoryRisk {
	risks := []CategoryRisk{}

	// Water damage risk (despite water resistance)
	if w.FitnessFeatures.WaterResistanceATM < 5 {
		risks = append(risks, CategoryRisk{
			Type:        "water_damage",
			Description: "Limited water resistance for swimming",
			Severity:    "medium",
			Impact:      1.12,
		})
	}

	// Strap/band damage risk
	risks = append(risks, CategoryRisk{
		Type:        "strap_damage",
		Description: "Wear and tear on straps/bands",
		Severity:    "low",
		Impact:      1.05,
	})

	// Sensor malfunction risk
	if w.HealthSensors.HasECG || w.HealthSensors.HasBloodOxygen {
		risks = append(risks, CategoryRisk{
			Type:        "sensor_malfunction",
			Description: "Advanced health sensor failure risk",
			Severity:    "medium",
			Impact:      1.08,
		})
	}

	// Battery degradation (faster for watches)
	if w.Device.DeviceFinancial.PurchaseDate != nil {
		deviceAge := time.Since(*w.Device.DeviceFinancial.PurchaseDate)
		if deviceAge > 1*365*24*time.Hour {
			risks = append(risks, CategoryRisk{
				Type:        "battery_degradation",
				Description: "Smartwatch battery degradation",
				Severity:    "high",
				Impact:      1.15,
			})
		}
	}

	// Screen damage (smaller but more exposed)
	risks = append(risks, CategoryRisk{
		Type:        "screen_damage",
		Description: "Screen damage from daily wear",
		Severity:    "medium",
		Impact:      1.10,
	})

	return risks
}

// CalculateCategoryPremiumAdjustment calculates smartwatch-specific premium adjustment
func (w *SmartwatchDevice) CalculateCategoryPremiumAdjustment() float64 {
	adjustment := 1.0

	// Health-focused watches need more careful handling
	if w.GetSubCategory() == "health_focused" {
		adjustment *= 0.92
	}

	// Sport watches face more wear
	if w.GetSubCategory() == "sport_focused" {
		adjustment *= 1.15
	}

	// Cellular models are more complex
	if w.ConnectivityFeatures.HasCellular {
		adjustment *= 1.10
	}

	// Advanced health sensors increase repair cost
	if w.HealthSensors.HasECG {
		adjustment *= 1.08
	}

	// Higher water resistance means more active use
	if w.FitnessFeatures.WaterResistanceATM >= 10 {
		adjustment *= 1.12
	}

	return adjustment
}

// GetCategoryDepreciationRate returns smartwatch-specific depreciation rate
func (w *SmartwatchDevice) GetCategoryDepreciationRate() float64 {
	// Smartwatches depreciate faster due to battery and technology
	baseRate := 0.35 // 35% per year

	// Premium watches hold value better
	if w.GetSubCategory() == "premium" {
		baseRate = 0.30
	} else if w.GetSubCategory() == "basic" {
		baseRate = 0.45
	}

	// Apple watches hold value better
	if w.Device.DeviceClassification.Brand == "Apple" {
		baseRate *= 0.85
	}

	return baseRate
}

// GetCategorySpecificCoverage returns smartwatch-specific coverage options
func (w *SmartwatchDevice) GetCategorySpecificCoverage() []CoverageType {
	return []CoverageType{
		{
			Type:        "screen_damage",
			Name:        "Display Protection",
			Description: "Covers screen cracks and display issues",
			MaxAmount:   300,
			Deductible:  30,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "water_damage",
			Name:        "Water Damage Protection",
			Description: "Covers water damage beyond rated resistance",
			MaxAmount:   400,
			Deductible:  50,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "strap_replacement",
			Name:        "Strap/Band Protection",
			Description: "Covers strap and band replacement",
			MaxAmount:   100,
			Deductible:  0,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "battery_replacement",
			Name:        "Battery Protection",
			Description: "Covers battery replacement when health drops",
			MaxAmount:   100,
			Deductible:  0,
			IsRequired:  false,
			IsAvailable: true,
		},
		{
			Type:        "sensor_malfunction",
			Name:        "Sensor Protection",
			Description: "Covers health sensor malfunctions",
			MaxAmount:   200,
			Deductible:  25,
			IsRequired:  false,
			IsAvailable: w.HealthSensors.HasECG || w.HealthSensors.HasBloodOxygen,
		},
	}
}

// IsEligibleForCategoryPrograms checks eligibility for smartwatch-specific programs
func (w *SmartwatchDevice) IsEligibleForCategoryPrograms() map[string]bool {
	programs := make(map[string]bool)

	// Default to now if purchase date not set
	purchaseDate := time.Now()
	if w.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *w.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	// Trade-in program (shorter eligibility for watches)
	programs["trade_in"] = deviceAge < 2*365*24*time.Hour &&
		w.Device.DevicePhysicalCondition.Grade != "F"

	// Battery replacement (more critical for watches)
	programs["battery_replacement"] = deviceAge > 1*365*24*time.Hour

	// Strap replacement program
	programs["strap_replacement"] = true

	// Health sensor calibration
	programs["sensor_calibration"] = w.HealthSensors.HasECG ||
		w.HealthSensors.HasBloodOxygen

	// Fitness tracker upgrade
	programs["fitness_upgrade"] = !w.FitnessFeatures.HasGPS &&
		w.GetSubCategory() != "premium"

	// Sport band program
	programs["sport_band"] = w.GetSubCategory() == "sport_focused"

	return programs
}

// GetCategoryMaintenanceSchedule returns smartwatch-specific maintenance schedule
func (w *SmartwatchDevice) GetCategoryMaintenanceSchedule() MaintenanceSchedule {
	now := time.Now()

	// Default to now if purchase date not set
	purchaseDate := now
	if w.Device.DeviceFinancial.PurchaseDate != nil {
		purchaseDate = *w.Device.DeviceFinancial.PurchaseDate
	}

	deviceAge := time.Since(purchaseDate)

	intervals := []MaintenanceInterval{
		{
			Type:        "software_update",
			Description: "WatchOS/WearOS updates",
			Frequency:   30 * 24 * time.Hour, // Monthly
			IsCritical:  true,
			Cost:        0,
		},
		{
			Type:        "strap_inspection",
			Description: "Check strap/band condition",
			Frequency:   90 * 24 * time.Hour, // Quarterly
			IsCritical:  false,
			Cost:        0,
		},
		{
			Type:        "sensor_cleaning",
			Description: "Clean health sensors",
			Frequency:   30 * 24 * time.Hour, // Monthly
			IsCritical:  false,
			Cost:        0,
		},
		{
			Type:        "water_seal_check",
			Description: "Verify water resistance seals",
			Frequency:   180 * 24 * time.Hour, // Semi-annually
			IsCritical:  true,
			Cost:        30,
		},
		{
			Type:        "battery_optimization",
			Description: "Battery health check and optimization",
			Frequency:   90 * 24 * time.Hour, // Quarterly
			IsCritical:  false,
			Cost:        0,
		},
	}

	// Add battery replacement after 1.5 years for watches
	if deviceAge > 18*30*24*time.Hour {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "battery_replacement",
			Description: "Consider battery replacement",
			Frequency:   0, // One-time
			IsCritical:  true,
			Cost:        80,
		})
	}

	// Add sensor calibration for health watches
	if w.HealthSensors.HasECG {
		intervals = append(intervals, MaintenanceInterval{
			Type:        "sensor_calibration",
			Description: "Calibrate health sensors",
			Frequency:   365 * 24 * time.Hour, // Annually
			IsCritical:  false,
			Cost:        50,
		})
	}

	return MaintenanceSchedule{
		Intervals: intervals,
		NextDue:   now.Add(30 * 24 * time.Hour),
		LastDone:  now,
	}
}
