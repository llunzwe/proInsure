package smartwatch

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartwatchPerformance tracks smartwatch-specific performance metrics
type SmartwatchPerformance struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordDate time.Time `json:"record_date"`

	// Processor Performance
	ProcessorModel    string  `json:"processor_model"` // S8, S9, W5+ Gen 1, Exynos W920
	ProcessorCores    int     `json:"processor_cores"`
	ProcessorSpeed    float64 `json:"processor_speed"` // GHz
	CoprocessorActive bool    `gorm:"default:false" json:"coprocessor_active"`

	// Memory & Storage
	RAMTotal     int `json:"ram_total"`     // MB
	RAMAvailable int `json:"ram_available"` // MB
	StorageTotal int `json:"storage_total"` // GB
	StorageUsed  int `json:"storage_used"`  // GB
	AppCount     int `json:"app_count"`

	// Battery Performance
	BatteryCapacity int     `json:"battery_capacity"` // mAh
	BatteryHealth   int     `json:"battery_health"`   // percentage
	BatteryCycles   int     `json:"battery_cycles"`
	BatteryLife     float64 `json:"battery_life"`  // hours
	ChargingTime    float64 `json:"charging_time"` // hours to full
	PowerSaveMode   bool    `gorm:"default:false" json:"power_save_mode"`

	// Display Performance
	DisplayType     string `json:"display_type"` // OLED, LTPO, Retina
	AlwaysOnDisplay bool   `gorm:"default:false" json:"always_on_display"`
	Brightness      int    `json:"brightness"`   // percentage
	RefreshRate     int    `json:"refresh_rate"` // Hz (1Hz for AOD, 60Hz normal)
	RaiseToWake     bool   `gorm:"default:false" json:"raise_to_wake"`
	ScreenTimeout   int    `json:"screen_timeout"` // seconds

	// Health Sensor Performance
	HeartRateAccuracy    float64 `json:"heart_rate_accuracy"` // percentage
	HeartRateSampling    int     `json:"heart_rate_sampling"` // samples per minute
	SpO2Accuracy         float64 `json:"spo2_accuracy"`       // percentage
	SpO2Available        bool    `gorm:"default:false" json:"spo2_available"`
	ECGAvailable         bool    `gorm:"default:false" json:"ecg_available"`
	ECGReadings          int     `json:"ecg_readings"` // total readings taken
	TemperatureSensor    bool    `gorm:"default:false" json:"temperature_sensor"`
	BloodPressureCapable bool    `gorm:"default:false" json:"blood_pressure_capable"`

	// Fitness Tracking Performance
	StepAccuracy         float64 `json:"step_accuracy"` // percentage
	StepsToday           int     `json:"steps_today"`
	CaloriesAccuracy     float64 `json:"calories_accuracy"` // percentage
	CaloriesToday        int     `json:"calories_today"`
	ActiveMinutes        int     `json:"active_minutes"`
	StandHours           int     `json:"stand_hours"`
	WorkoutsRecorded     int     `json:"workouts_recorded"`
	GPSAccuracy          float64 `json:"gps_accuracy"` // meters
	AutoWorkoutDetection bool    `gorm:"default:false" json:"auto_workout_detection"`

	// Sleep Tracking
	SleepTrackingEnabled bool    `gorm:"default:false" json:"sleep_tracking_enabled"`
	SleepAccuracy        float64 `json:"sleep_accuracy"`                    // percentage
	LastNightSleep       float64 `json:"last_night_sleep"`                  // hours
	SleepStages          bool    `gorm:"default:false" json:"sleep_stages"` // REM, deep, light
	SleepBreathing       bool    `gorm:"default:false" json:"sleep_breathing"`

	// Environmental Sensors
	AltimeterAccuracy  float64 `json:"altimeter_accuracy"` // meters
	BarometerAvailable bool    `gorm:"default:false" json:"barometer_available"`
	CompassAccuracy    float64 `json:"compass_accuracy"` // degrees
	NoiseMonitoring    bool    `gorm:"default:false" json:"noise_monitoring"`
	NoiseLevel         int     `json:"noise_level"` // dB

	// App Performance
	AppLaunchTime       float64 `json:"app_launch_time"`   // seconds average
	AppResponseTime     float64 `json:"app_response_time"` // milliseconds
	AppCrashCount       int     `json:"app_crash_count"`
	ComplicationLoad    int     `json:"complication_load"`    // number of complications
	WatchFaceComplexity string  `json:"watchface_complexity"` // simple, moderate, complex

	// Connectivity Performance
	BluetoothStability  float64 `json:"bluetooth_stability"`   // percentage uptime
	WiFiConnectionTime  float64 `json:"wifi_connection_time"`  // seconds
	CellularSignal      int     `json:"cellular_signal"`       // dBm (if LTE model)
	PhoneConnectionLost int     `json:"phone_connection_lost"` // times disconnected
	SyncLatency         float64 `json:"sync_latency"`          // seconds

	// Notification Performance
	NotificationsReceived int     `json:"notifications_received"`
	NotificationDelay     float64 `json:"notification_delay"` // seconds
	HapticIntensity       string  `json:"haptic_intensity"`   // light, medium, strong
	SoundEnabled          bool    `gorm:"default:false" json:"sound_enabled"`

	// Water Resistance
	WaterLockEnabled bool    `gorm:"default:false" json:"water_lock_enabled"`
	WaterEjections   int     `json:"water_ejections"` // times used
	SwimTracking     bool    `gorm:"default:false" json:"swim_tracking"`
	DepthReached     float64 `json:"depth_reached"` // meters

	// Fall Detection & Safety
	FallDetectionEnabled bool `gorm:"default:false" json:"fall_detection_enabled"`
	FallsDetected        int  `json:"falls_detected"`
	SOSActivations       int  `json:"sos_activations"`
	EmergencyContacts    int  `json:"emergency_contacts"`
	CrashDetection       bool `gorm:"default:false" json:"crash_detection"`

	// Voice Assistant
	VoiceAssistantEnabled bool    `gorm:"default:false" json:"voice_assistant_enabled"`
	VoiceCommands         int     `json:"voice_commands"`
	VoiceResponseTime     float64 `json:"voice_response_time"` // seconds
	DictationAccuracy     float64 `json:"dictation_accuracy"`  // percentage

	// System Performance
	SystemUptime     int        `json:"system_uptime"` // hours
	RestartCount     int        `json:"restart_count"`
	FreezeCount      int        `json:"freeze_count"`
	UpdatesPending   int        `json:"updates_pending"`
	LastUpdate       *time.Time `json:"last_update"`
	PerformanceScore float64    `json:"performance_score"` // 0-100

	// Thermal Management
	Temperature       float64 `json:"temperature"` // Celsius
	ThermalThrottling bool    `gorm:"default:false" json:"thermal_throttling"`

	// Usage Patterns
	DailyActiveHours float64 `json:"daily_active_hours"`
	ScreenOnTime     float64 `json:"screen_on_time"`    // hours
	InteractionCount int     `json:"interaction_count"` // taps, swipes, crown turns
	CrownRotations   int     `json:"crown_rotations"`
	ButtonPresses    int     `json:"button_presses"`
}

// Business logic methods

// IsHighPerformanceWatch checks if high-end model
func (sp *SmartwatchPerformance) IsHighPerformanceWatch() bool {
	// Check processor (latest Apple S9, Snapdragon W5+)
	highEndProcessors := []string{"S9", "S8", "W5+", "Exynos W930"}
	for _, proc := range highEndProcessors {
		if sp.ProcessorModel == proc {
			return true
		}
	}
	// Check RAM (2GB+ is high-end)
	if sp.RAMTotal >= 2000 {
		return true
	}
	// Check storage (32GB+ is high-end)
	if sp.StorageTotal >= 32 {
		return true
	}
	return false
}

// HasHealthIssues checks for health tracking problems
func (sp *SmartwatchPerformance) HasHealthIssues() bool {
	// Low accuracy readings
	if sp.HeartRateAccuracy < 90 {
		return true
	}
	if sp.SpO2Available && sp.SpO2Accuracy < 85 {
		return true
	}
	// Fitness tracking issues
	if sp.StepAccuracy < 85 {
		return true
	}
	if sp.CaloriesAccuracy < 80 {
		return true
	}
	return false
}

// GetBatteryEfficiency calculates battery efficiency
func (sp *SmartwatchPerformance) GetBatteryEfficiency() float64 {
	if sp.BatteryCapacity == 0 {
		return 0
	}

	// Calculate efficiency based on battery life per mAh
	efficiency := sp.BatteryLife / (float64(sp.BatteryCapacity) / 100)

	// Normalize to 0-100 scale
	// Assuming 6 hours per 100mAh is excellent
	normalizedScore := efficiency * 16.67

	// Adjust for always-on display
	if sp.AlwaysOnDisplay {
		normalizedScore *= 1.2 // Bonus for efficiency with AOD
	}

	if normalizedScore > 100 {
		normalizedScore = 100
	}

	return normalizedScore
}

// NeedsBatteryReplacement checks battery health
func (sp *SmartwatchPerformance) NeedsBatteryReplacement() bool {
	// Battery health below 80%
	if sp.BatteryHealth < 80 {
		return true
	}
	// High cycle count (watches have smaller batteries)
	if sp.BatteryCycles > 300 {
		return true
	}
	// Poor battery life
	if sp.BatteryLife < 12 && !sp.AlwaysOnDisplay {
		return true
	}
	return false
}

// GetFitnessAccuracyScore rates fitness tracking
func (sp *SmartwatchPerformance) GetFitnessAccuracyScore() float64 {
	score := 0.0
	components := 0

	// Step accuracy
	score += sp.StepAccuracy
	components++

	// Calorie accuracy
	score += sp.CaloriesAccuracy
	components++

	// Heart rate accuracy
	score += sp.HeartRateAccuracy
	components++

	// GPS accuracy (convert to percentage)
	if sp.GPSAccuracy > 0 {
		gpsScore := 100 - (sp.GPSAccuracy * 10) // Assuming <10m is good
		if gpsScore < 0 {
			gpsScore = 0
		}
		score += gpsScore
		components++
	}

	// SpO2 accuracy if available
	if sp.SpO2Available {
		score += sp.SpO2Accuracy
		components++
	}

	if components > 0 {
		return score / float64(components)
	}
	return 0
}

// GetHealthMonitoringScore rates health features
func (sp *SmartwatchPerformance) GetHealthMonitoringScore() float64 {
	score := 50.0 // Base score

	// Heart rate monitoring
	if sp.HeartRateAccuracy > 95 {
		score += 10
	} else if sp.HeartRateAccuracy > 90 {
		score += 5
	}

	// SpO2 monitoring
	if sp.SpO2Available {
		score += 10
		if sp.SpO2Accuracy > 90 {
			score += 5
		}
	}

	// ECG capability
	if sp.ECGAvailable {
		score += 15
	}

	// Temperature sensor
	if sp.TemperatureSensor {
		score += 5
	}

	// Blood pressure
	if sp.BloodPressureCapable {
		score += 10
	}

	// Sleep tracking
	if sp.SleepTrackingEnabled {
		score += 5
		if sp.SleepStages {
			score += 5
		}
	}

	// Fall detection
	if sp.FallDetectionEnabled {
		score += 5
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// IsConnectivityStable checks connection stability
func (sp *SmartwatchPerformance) IsConnectivityStable() bool {
	// Bluetooth stability check
	if sp.BluetoothStability < 95 {
		return false
	}
	// Frequent disconnections
	if sp.PhoneConnectionLost > 5 {
		return false
	}
	// High sync latency
	if sp.SyncLatency > 5 {
		return false
	}
	return true
}

// HasPerformanceIssues checks for problems
func (sp *SmartwatchPerformance) HasPerformanceIssues() bool {
	// App crashes
	if sp.AppCrashCount > 3 {
		return true
	}
	// System freezes
	if sp.FreezeCount > 0 {
		return true
	}
	// Frequent restarts
	if sp.RestartCount > 2 {
		return true
	}
	// Slow app launch
	if sp.AppLaunchTime > 3 {
		return true
	}
	// Thermal throttling
	if sp.ThermalThrottling {
		return true
	}
	// Low RAM
	ramUsagePercent := float64(sp.RAMTotal-sp.RAMAvailable) / float64(sp.RAMTotal) * 100
	if ramUsagePercent > 90 {
		return true
	}
	return false
}

// GetDisplayEfficiency rates display efficiency
func (sp *SmartwatchPerformance) GetDisplayEfficiency() float64 {
	score := 50.0

	// LTPO displays are more efficient
	if sp.DisplayType == "LTPO" || sp.DisplayType == "LTPO OLED" {
		score += 20
	}

	// Always-on display management
	if sp.AlwaysOnDisplay {
		// Low refresh rate for AOD is good
		if sp.RefreshRate <= 1 {
			score += 15
		}
	}

	// Brightness optimization
	if sp.Brightness < 50 {
		score += 10
	} else if sp.Brightness > 80 {
		score -= 10
	}

	// Auto features
	if sp.RaiseToWake {
		score += 5
	}

	// Screen timeout
	if sp.ScreenTimeout <= 15 {
		score += 5
	}

	return score
}

// GetOverallPerformanceScore calculates overall score
func (sp *SmartwatchPerformance) GetOverallPerformanceScore() float64 {
	if sp.PerformanceScore > 0 {
		return sp.PerformanceScore
	}

	// Component scores with weights
	batteryScore := sp.GetBatteryEfficiency() * 0.25
	fitnessScore := sp.GetFitnessAccuracyScore() * 0.20
	healthScore := sp.GetHealthMonitoringScore() * 0.20
	displayScore := sp.GetDisplayEfficiency() * 0.10

	// Connectivity score
	connectivityScore := 100.0
	if !sp.IsConnectivityStable() {
		connectivityScore = 60.0
	}
	connectivityScore *= 0.15

	// Performance score
	performanceScore := 100.0
	if sp.HasPerformanceIssues() {
		performanceScore = 50.0
	}
	performanceScore *= 0.10

	overall := batteryScore + fitnessScore + healthScore +
		displayScore + connectivityScore + performanceScore

	if overall > 100 {
		overall = 100
	}

	return overall
}

// ShouldOptimizeSettings checks if optimization needed
func (sp *SmartwatchPerformance) ShouldOptimizeSettings() bool {
	// Battery life too short
	if sp.BatteryLife < 18 {
		return true
	}
	// Performance issues
	if sp.HasPerformanceIssues() {
		return true
	}
	// Low overall score
	if sp.GetOverallPerformanceScore() < 60 {
		return true
	}
	// Too many complications
	if sp.ComplicationLoad > 8 {
		return true
	}
	return false
}

// IsSwimProof checks swim tracking capability
func (sp *SmartwatchPerformance) IsSwimProof() bool {
	return sp.SwimTracking && sp.WaterEjections > 0
}

// GetSafetyScore rates safety features
func (sp *SmartwatchPerformance) GetSafetyScore() float64 {
	score := 0.0

	if sp.FallDetectionEnabled {
		score += 30
	}
	if sp.SOSActivations >= 0 {
		score += 20
	}
	if sp.EmergencyContacts > 0 {
		score += 20
	}
	if sp.CrashDetection {
		score += 30
	}

	return score
}
