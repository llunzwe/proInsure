package smartphone

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartphonePerformance tracks smartphone-specific performance metrics
type SmartphonePerformance struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordDate time.Time `json:"record_date"`

	// CPU Performance
	CPUModel        string  `json:"cpu_model"`         // Snapdragon 8 Gen 2, A16 Bionic, etc
	CPUCores        int     `json:"cpu_cores"`         // 8, 10, etc
	CPUMaxFrequency float64 `json:"cpu_max_frequency"` // GHz
	CPUUsageAverage float64 `json:"cpu_usage_average"` // percentage
	CPUThrottling   bool    `gorm:"default:false" json:"cpu_throttling"`
	CPUTemperature  float64 `json:"cpu_temperature"` // Celsius

	// GPU Performance
	GPUModel        string  `json:"gpu_model"`         // Adreno 740, Apple GPU, Mali-G710
	GPUUsageAverage float64 `json:"gpu_usage_average"` // percentage
	GPUTemperature  float64 `json:"gpu_temperature"`   // Celsius

	// Gaming Performance
	GamingFPS       float64 `json:"gaming_fps"`       // frames per second
	GamingStability float64 `json:"gaming_stability"` // FPS consistency percentage
	ThermalThrottle int     `json:"thermal_throttle"` // seconds until throttling
	GameLoadTime    float64 `json:"game_load_time"`   // seconds
	TouchLatency    float64 `json:"touch_latency"`    // milliseconds

	// App Performance
	AppLaunchTimeCold float64 `json:"app_launch_time_cold"` // seconds
	AppLaunchTimeWarm float64 `json:"app_launch_time_warm"` // seconds
	AppCrashCount     int     `json:"app_crash_count"`
	AppFreezeCount    int     `json:"app_freeze_count"`
	BackgroundApps    int     `json:"background_apps"`

	// Memory Performance
	RAMTotal       int     `json:"ram_total"`       // GB
	RAMAvailable   int     `json:"ram_available"`   // GB
	RAMUsage       float64 `json:"ram_usage"`       // percentage
	MemoryPressure string  `json:"memory_pressure"` // low, medium, high, critical
	SwapUsage      float64 `json:"swap_usage"`      // percentage
	MemoryLeaks    int     `json:"memory_leaks"`    // detected leaks

	// Storage Performance
	StorageTotal int     `json:"storage_total"` // GB
	StorageUsed  int     `json:"storage_used"`  // GB
	StorageFree  int     `json:"storage_free"`  // GB
	StorageSpeed string  `json:"storage_speed"` // UFS 3.1, UFS 4.0, NVMe
	ReadSpeed    float64 `json:"read_speed"`    // MB/s
	WriteSpeed   float64 `json:"write_speed"`   // MB/s
	IOPSRead     int     `json:"iops_read"`     // Input/output operations per second
	IOPSWrite    int     `json:"iops_write"`

	// Battery Performance
	BatteryCapacity int     `json:"battery_capacity"` // mAh
	BatteryHealth   int     `json:"battery_health"`   // percentage
	BatteryCycles   int     `json:"battery_cycles"`
	ChargingSpeed   int     `json:"charging_speed"` // watts
	ScreenOnTime    float64 `json:"screen_on_time"` // hours
	StandbyTime     float64 `json:"standby_time"`   // hours
	BatteryDrain    float64 `json:"battery_drain"`  // percentage per hour

	// Display Performance
	RefreshRate     int    `json:"refresh_rate"`   // Hz (60, 90, 120, 144)
	TouchSampling   int    `json:"touch_sampling"` // Hz (120, 240, 360, 480)
	Brightness      int    `json:"brightness"`     // nits
	ColorGamut      string `json:"color_gamut"`    // sRGB, DCI-P3, etc
	HDRSupport      string `json:"hdr_support"`    // HDR10, HDR10+, Dolby Vision
	AlwaysOnDisplay bool   `gorm:"default:false" json:"always_on_display"`

	// Camera Performance
	CameraLaunchTime float64 `json:"camera_launch_time"` // seconds
	PhotoProcessTime float64 `json:"photo_process_time"` // seconds
	VideoStability   string  `json:"video_stability"`    // none, EIS, OIS, both
	LowLightScore    float64 `json:"low_light_score"`    // 0-100
	AutoFocusSpeed   float64 `json:"autofocus_speed"`    // milliseconds

	// AI/ML Performance
	AIChipset       string  `json:"ai_chipset"`         // Neural Engine, Tensor, etc
	AIOperations    float64 `json:"ai_operations"`      // TOPS (trillion operations per second)
	MLModelLoadTime float64 `json:"ml_model_load_time"` // seconds
	MLInferenceTime float64 `json:"ml_inference_time"`  // milliseconds

	// Benchmark Scores
	AntutuScore      int     `json:"antutu_score"`
	Geekbench5Single int     `json:"geekbench5_single"`
	Geekbench5Multi  int     `json:"geekbench5_multi"`
	GFXBenchScore    float64 `json:"gfxbench_score"`
	PCMarkScore      int     `json:"pcmark_score"`

	// System Health
	SystemUptime     int     `json:"system_uptime"` // hours
	RebootCount      int     `json:"reboot_count"`
	SystemErrors     int     `json:"system_errors"`
	KernelPanics     int     `json:"kernel_panics"`
	PerformanceScore float64 `json:"performance_score"` // 0-100 overall score

	// Thermal Management
	ThermalState   string  `json:"thermal_state"`   // normal, fair, serious, critical
	CoolingSystem  string  `json:"cooling_system"`  // passive, vapor_chamber, graphite
	MaxTemperature float64 `json:"max_temperature"` // Celsius
	AvgTemperature float64 `json:"avg_temperature"` // Celsius

	// Multitasking
	AppsInMemory    int     `json:"apps_in_memory"`
	AppSwitchTime   float64 `json:"app_switch_time"` // milliseconds
	SplitScreenApps int     `json:"split_screen_apps"`
	PiPSupported    bool    `gorm:"default:false" json:"pip_supported"`
}

// Business logic methods

// IsHighPerformanceDevice checks if device is high-performance
func (sp *SmartphonePerformance) IsHighPerformanceDevice() bool {
	// Check benchmark scores
	if sp.AntutuScore > 1000000 || sp.Geekbench5Single > 1500 {
		return true
	}
	// Check specs
	if sp.RAMTotal >= 12 && sp.CPUMaxFrequency >= 3.0 {
		return true
	}
	// Check display
	if sp.RefreshRate >= 120 {
		return true
	}
	return false
}

// IsGamingCapable checks if device is suitable for gaming
func (sp *SmartphonePerformance) IsGamingCapable() bool {
	if sp.GamingFPS >= 60 && sp.GamingStability >= 90 {
		return true
	}
	if sp.RAMTotal >= 8 && sp.RefreshRate >= 90 {
		return true
	}
	return sp.TouchSampling >= 240 && sp.TouchLatency < 20
}

// HasPerformanceIssues checks for performance problems
func (sp *SmartphonePerformance) HasPerformanceIssues() bool {
	// Check for throttling
	if sp.CPUThrottling || sp.ThermalState == "critical" {
		return true
	}
	// Check for memory issues
	if sp.MemoryPressure == "critical" || sp.MemoryLeaks > 0 {
		return true
	}
	// Check for crashes
	if sp.AppCrashCount > 5 || sp.KernelPanics > 0 {
		return true
	}
	// Check storage
	if sp.StorageFree < 1 { // Less than 1GB free
		return true
	}
	return false
}

// GetBatteryEfficiency calculates battery efficiency score
func (sp *SmartphonePerformance) GetBatteryEfficiency() float64 {
	if sp.BatteryCapacity == 0 {
		return 0
	}

	// Calculate efficiency based on screen-on time per mAh
	efficiency := (sp.ScreenOnTime * 1000) / float64(sp.BatteryCapacity)

	// Normalize to 0-100 scale
	// Assuming 2 hours per 1000mAh is excellent
	normalizedScore := efficiency * 50

	if normalizedScore > 100 {
		normalizedScore = 100
	}

	return normalizedScore
}

// NeedsBatteryReplacement checks if battery needs replacement
func (sp *SmartphonePerformance) NeedsBatteryReplacement() bool {
	// Battery health below 80% typically needs replacement
	if sp.BatteryHealth < 80 {
		return true
	}
	// High cycle count
	if sp.BatteryCycles > 500 {
		return true
	}
	// Poor screen-on time
	if sp.ScreenOnTime < 3 && sp.BatteryCapacity > 3000 {
		return true
	}
	return false
}

// GetStorageHealthScore calculates storage health
func (sp *SmartphonePerformance) GetStorageHealthScore() float64 {
	score := 100.0

	// Deduct for low free space
	freePercentage := float64(sp.StorageFree) / float64(sp.StorageTotal) * 100
	if freePercentage < 10 {
		score -= 30
	} else if freePercentage < 20 {
		score -= 15
	}

	// Deduct for slow speeds (assuming UFS 3.0 baseline)
	if sp.ReadSpeed < 1000 { // MB/s
		score -= 20
	}
	if sp.WriteSpeed < 500 {
		score -= 20
	}

	// Deduct for low IOPS
	if sp.IOPSRead < 50000 {
		score -= 10
	}

	if score < 0 {
		score = 0
	}

	return score
}

// IsOverheating checks if device is overheating
func (sp *SmartphonePerformance) IsOverheating() bool {
	// CPU temperature above 45°C is concerning for sustained use
	if sp.CPUTemperature > 45 {
		return true
	}
	// GPU temperature above 50°C
	if sp.GPUTemperature > 50 {
		return true
	}
	// Thermal state check
	if sp.ThermalState == "serious" || sp.ThermalState == "critical" {
		return true
	}
	// Average temperature above 40°C
	if sp.AvgTemperature > 40 {
		return true
	}
	return false
}

// GetDisplayQualityScore rates display quality
func (sp *SmartphonePerformance) GetDisplayQualityScore() float64 {
	score := 50.0 // Base score

	// Refresh rate scoring
	switch {
	case sp.RefreshRate >= 144:
		score += 20
	case sp.RefreshRate >= 120:
		score += 15
	case sp.RefreshRate >= 90:
		score += 10
	case sp.RefreshRate >= 60:
		score += 5
	}

	// Touch sampling scoring
	if sp.TouchSampling >= 480 {
		score += 10
	} else if sp.TouchSampling >= 240 {
		score += 7
	} else if sp.TouchSampling >= 120 {
		score += 4
	}

	// HDR support
	switch sp.HDRSupport {
	case "Dolby Vision":
		score += 10
	case "HDR10+":
		score += 8
	case "HDR10":
		score += 5
	}

	// Color gamut
	if sp.ColorGamut == "DCI-P3" {
		score += 5
	}

	// Brightness (assuming 1000+ nits is excellent)
	if sp.Brightness >= 1000 {
		score += 5
	}

	return score
}

// GetCameraPerformanceScore rates camera performance
func (sp *SmartphonePerformance) GetCameraPerformanceScore() float64 {
	score := 50.0

	// Launch time (faster is better)
	if sp.CameraLaunchTime < 1.0 {
		score += 10
	} else if sp.CameraLaunchTime < 2.0 {
		score += 5
	}

	// Processing time
	if sp.PhotoProcessTime < 0.5 {
		score += 10
	} else if sp.PhotoProcessTime < 1.0 {
		score += 5
	}

	// Video stability
	switch sp.VideoStability {
	case "both":
		score += 15
	case "OIS":
		score += 10
	case "EIS":
		score += 7
	}

	// Low light performance
	score += sp.LowLightScore * 0.15

	// Autofocus speed (faster is better)
	if sp.AutoFocusSpeed < 100 {
		score += 10
	} else if sp.AutoFocusSpeed < 200 {
		score += 5
	}

	return score
}

// GetOverallPerformanceScore calculates overall performance
func (sp *SmartphonePerformance) GetOverallPerformanceScore() float64 {
	if sp.PerformanceScore > 0 {
		return sp.PerformanceScore
	}

	// Weight different aspects
	cpuScore := 0.0
	if sp.Geekbench5Multi > 0 {
		cpuScore = float64(sp.Geekbench5Multi) / 50 // Normalize to ~100
	}

	memoryScore := 100 - sp.RAMUsage
	if sp.MemoryPressure == "high" {
		memoryScore *= 0.7
	} else if sp.MemoryPressure == "critical" {
		memoryScore *= 0.5
	}

	storageScore := sp.GetStorageHealthScore()
	batteryScore := sp.GetBatteryEfficiency()
	displayScore := sp.GetDisplayQualityScore()

	// Weighted average
	overall := (cpuScore*0.25 + memoryScore*0.20 +
		storageScore*0.15 + batteryScore*0.20 + displayScore*0.20)

	// Apply penalties
	if sp.HasPerformanceIssues() {
		overall *= 0.8
	}
	if sp.IsOverheating() {
		overall *= 0.9
	}

	if overall > 100 {
		overall = 100
	}

	return overall
}

// ShouldOptimize checks if device needs optimization
func (sp *SmartphonePerformance) ShouldOptimize() bool {
	// Low overall score
	if sp.GetOverallPerformanceScore() < 60 {
		return true
	}
	// High RAM usage
	if sp.RAMUsage > 90 {
		return true
	}
	// Low storage
	if sp.StorageFree < 2 {
		return true
	}
	// Many background apps
	if sp.BackgroundApps > 20 {
		return true
	}
	return false
}
