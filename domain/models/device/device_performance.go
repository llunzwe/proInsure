package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DevicePerformanceMetrics tracks overall system performance
type DevicePerformanceMetrics struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	MeasurementDate time.Time `json:"measurement_date"`

	// CPU Usage Patterns
	CPUUsageAverage float64 `json:"cpu_usage_average"`                  // percentage
	CPUUsagePeak    float64 `json:"cpu_usage_peak"`                     // percentage
	CPUIdleTime     float64 `json:"cpu_idle_time"`                      // percentage
	CPUUsageHistory string  `gorm:"type:json" json:"cpu_usage_history"` // JSON array

	// RAM Utilization
	RAMUsageAverage float64 `json:"ram_usage_average"` // percentage
	RAMUsagePeak    float64 `json:"ram_usage_peak"`    // percentage
	RAMAvailable    int64   `json:"ram_available"`     // MB
	RAMTotal        int64   `json:"ram_total"`         // MB
	RAMTrend        string  `json:"ram_trend"`         // increasing, stable, decreasing

	// Storage Usage
	StorageUsed       int64   `json:"storage_used"`        // GB
	StorageTotal      int64   `json:"storage_total"`       // GB
	StorageGrowthRate float64 `json:"storage_growth_rate"` // GB per month
	StorageType       string  `json:"storage_type"`        // SSD, eMMC, UFS

	// App Performance
	AppLaunchTimeAvg  float64 `json:"app_launch_time_avg"`  // seconds
	AppLaunchTimeCold float64 `json:"app_launch_time_cold"` // seconds
	AppLaunchTimeWarm float64 `json:"app_launch_time_warm"` // seconds
	AppCrashCount     int     `json:"app_crash_count"`
	AppFreezeCount    int     `json:"app_freeze_count"`

	// System Responsiveness
	ResponsivenessScore float64 `json:"responsiveness_score"` // 0-100
	TouchLatency        float64 `json:"touch_latency"`        // ms
	ScrollingSmooth     float64 `json:"scrolling_smooth"`     // FPS
	AnimationFPS        float64 `json:"animation_fps"`
	UILagIncidents      int     `json:"ui_lag_incidents"`

	// Multitasking Performance
	MultitaskingScore float64 `json:"multitasking_score"` // 0-100
	AppSwitchTime     float64 `json:"app_switch_time"`    // seconds
	BackgroundApps    int     `json:"background_apps"`
	MemoryPressure    string  `json:"memory_pressure"` // low, medium, high

	// Gaming Performance
	GamingFPS         float64 `json:"gaming_fps"`
	GamingStability   float64 `json:"gaming_stability"` // percentage
	GraphicsScore     float64 `json:"graphics_score"`   // benchmark
	GameLoadTime      float64 `json:"game_load_time"`   // seconds
	ThermalThrottling bool    `json:"thermal_throttling"`

	// Media Performance
	VideoPlaybackQuality string  `json:"video_playback_quality"` // 4K, 1080p, 720p
	VideoDroppedFrames   int     `json:"video_dropped_frames"`
	AudioLatency         float64 `json:"audio_latency"`         // ms
	CameraStartupTime    float64 `json:"camera_startup_time"`   // seconds
	PhotoProcessingTime  float64 `json:"photo_processing_time"` // seconds

	// Overall Performance
	OverallScore         float64    `json:"overall_score"`     // 0-100
	PerformanceTrend     string     `json:"performance_trend"` // improving, stable, degrading
	BenchmarkScore       int        `json:"benchmark_score"`
	ComparisonPercentile float64    `json:"comparison_percentile"` // percentile vs similar devices
	LastOptimization     *time.Time `json:"last_optimization"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceBatteryAnalytics provides detailed battery health tracking
type DeviceBatteryAnalytics struct {
	database.BaseModel
	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AnalysisDate time.Time `json:"analysis_date"`

	// Battery Health
	BatteryHealth    float64 `json:"battery_health"`                     // percentage
	DesignCapacity   int     `json:"design_capacity"`                    // mAh
	CurrentCapacity  int     `json:"current_capacity"`                   // mAh
	DegradationRate  float64 `json:"degradation_rate"`                   // percentage per month
	DegradationCurve string  `gorm:"type:json" json:"degradation_curve"` // JSON array

	// Charge Cycles
	ChargeCycles       int `json:"charge_cycles"`
	FullChargeCount    int `json:"full_charge_count"`
	PartialChargeCount int `json:"partial_charge_count"`
	DeepDischargeCount int `json:"deep_discharge_count"`
	OptimalChargeCount int `json:"optimal_charge_count"` // 20-80%

	// Usage Patterns
	ScreenOnTime      float64 `json:"screen_on_time"`                       // hours per day
	ScreenOffTime     float64 `json:"screen_off_time"`                      // hours per day
	ActiveUsageTime   float64 `json:"active_usage_time"`                    // hours
	StandbyTime       float64 `json:"standby_time"`                         // hours
	DailyUsagePattern string  `gorm:"type:json" json:"daily_usage_pattern"` // JSON object

	// App Consumption
	TopBatteryApps      string  `gorm:"type:json" json:"top_battery_apps"`      // JSON array
	AppBatteryBreakdown string  `gorm:"type:json" json:"app_battery_breakdown"` // JSON object
	BackgroundDrain     float64 `json:"background_drain"`                       // percentage
	SystemDrain         float64 `json:"system_drain"`                           // percentage

	// Standby Analysis
	StandbyDrain   float64 `json:"standby_drain"`   // percentage per hour
	DeepSleepTime  float64 `json:"deep_sleep_time"` // hours
	WakelocksCount int     `json:"wakelocks_count"`
	NetworkDrain   float64 `json:"network_drain"` // percentage

	// Charging Patterns
	FastChargingUsage     int     `json:"fast_charging_usage"`  // count
	FastChargingImpact    float64 `json:"fast_charging_impact"` // degradation factor
	WirelessChargingUsage int     `json:"wireless_charging_usage"`
	WirelessEfficiency    float64 `json:"wireless_efficiency"`              // percentage
	ChargingHabits        string  `gorm:"type:json" json:"charging_habits"` // JSON object

	// Temperature Impact
	ChargingTemperature float64 `json:"charging_temperature"` // avg celsius
	UsageTemperature    float64 `json:"usage_temperature"`    // avg celsius
	TemperatureEvents   int     `json:"temperature_events"`
	ThermalImpact       float64 `json:"thermal_impact"` // degradation factor

	// Optimization
	OptimalChargingTime  string  `json:"optimal_charging_time"`
	RecommendedSettings  string  `gorm:"type:json" json:"recommended_settings"` // JSON array
	PowerSavingPotential float64 `json:"power_saving_potential"`                // percentage
	EstimatedLifespan    int     `json:"estimated_lifespan"`                    // months

	// Replacement Prediction
	ReplacementNeeded    bool       `json:"replacement_needed"`
	ReplacementUrgency   string     `json:"replacement_urgency"` // low, medium, high, critical
	PredictedFailureDate *time.Time `json:"predicted_failure_date"`
	ReplacementCost      float64    `json:"replacement_cost"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceTemperatureHistory tracks thermal events and patterns
type DeviceTemperatureHistory struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordedAt time.Time `json:"recorded_at"`

	// Current Temperature
	CPUTemperature     float64 `json:"cpu_temperature"`     // celsius
	BatteryTemperature float64 `json:"battery_temperature"` // celsius
	GPUTemperature     float64 `json:"gpu_temperature"`     // celsius
	SkinTemperature    float64 `json:"skin_temperature"`    // device surface
	AmbientTemperature float64 `json:"ambient_temperature"`

	// Temperature Patterns
	AverageTemperature float64 `json:"average_temperature"`
	PeakTemperature    float64 `json:"peak_temperature"`
	MinTemperature     float64 `json:"min_temperature"`
	TemperatureHistory string  `gorm:"type:json" json:"temperature_history"` // JSON array

	// Overheating Events
	OverheatingIncidents int        `json:"overheating_incidents"`
	LastOverheating      *time.Time `json:"last_overheating"`
	OverheatingDuration  int        `json:"overheating_duration"`                // minutes total
	OverheatingCauses    string     `gorm:"type:json" json:"overheating_causes"` // JSON array

	// Thermal Throttling
	ThrottlingEvents    int     `json:"throttling_events"`
	ThrottlingDuration  int     `json:"throttling_duration"`  // minutes
	PerformanceImpact   float64 `json:"performance_impact"`   // percentage
	ThrottleTemperature float64 `json:"throttle_temperature"` // trigger temp

	// Usage Temperature
	ChargingTemperature float64 `json:"charging_temp_avg"`
	GamingTemperature   float64 `json:"gaming_temp_avg"`
	VideoTemperature    float64 `json:"video_temp_avg"`
	CameraTemperature   float64 `json:"camera_temp_avg"`
	IdleTemperature     float64 `json:"idle_temp_avg"`

	// Temperature Peaks
	ChargingPeak     float64 `json:"charging_peak"`
	GamingPeak       float64 `json:"gaming_peak"`
	VideoPeak        float64 `json:"video_peak"`
	NavigationPeak   float64 `json:"navigation_peak"`
	MultitaskingPeak float64 `json:"multitasking_peak"`

	// Cooling System
	CoolingType       string  `json:"cooling_type"`       // passive, active, liquid
	CoolingEfficiency float64 `json:"cooling_efficiency"` // percentage
	FanSpeed          int     `json:"fan_speed"`          // RPM if applicable
	ThermalPasteAge   int     `json:"thermal_paste_age"`  // months

	// Temperature Impact
	BatteryTempImpact    float64 `json:"battery_temp_impact"` // degradation factor
	PerformanceTempRatio float64 `json:"performance_temp_ratio"`
	UserComfortScore     float64 `json:"user_comfort_score"` // 0-100

	// Shutdowns
	ThermalShutdowns    int        `json:"thermal_shutdowns"`
	LastShutdown        *time.Time `json:"last_shutdown"`
	ShutdownTemperature float64    `json:"shutdown_temperature"`
	EmergencyCooldowns  int        `json:"emergency_cooldowns"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceStorageHealth monitors storage degradation and performance
type DeviceStorageHealth struct {
	database.BaseModel
	DeviceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	CheckDate time.Time `json:"check_date"`

	// Storage Health
	HealthStatus    float64 `json:"health_status"`    // percentage
	WearLevel       float64 `json:"wear_level"`       // percentage
	LifetimeWrites  int64   `json:"lifetime_writes"`  // TB
	LifetimeReads   int64   `json:"lifetime_reads"`   // TB
	DegradationRate float64 `json:"degradation_rate"` // percentage per year

	// Performance Metrics
	ReadSpeed        float64 `json:"read_speed"`         // MB/s
	WriteSpeed       float64 `json:"write_speed"`        // MB/s
	RandomReadSpeed  float64 `json:"random_read_speed"`  // IOPS
	RandomWriteSpeed float64 `json:"random_write_speed"` // IOPS
	SpeedDegradation float64 `json:"speed_degradation"`  // percentage

	// Bad Sectors
	BadSectorCount      int `json:"bad_sector_count"`
	RemappedSectors     int `json:"remapped_sectors"`
	PendingSectors      int `json:"pending_sectors"`
	UncorrectableErrors int `json:"uncorrectable_errors"`

	// Fragmentation
	FragmentationLevel  float64    `json:"fragmentation_level"` // percentage
	FileFragments       int        `json:"file_fragments"`
	LastDefragmentation *time.Time `json:"last_defragmentation"`
	DefragNeeded        bool       `json:"defrag_needed"`

	// Cache Performance
	CacheHitRate      float64 `json:"cache_hit_rate"` // percentage
	CacheSize         int64   `json:"cache_size"`     // MB
	CacheEfficiency   float64 `json:"cache_efficiency"`
	WriteCacheEnabled bool    `json:"write_cache_enabled"`

	// Storage Usage Patterns
	SystemDataSize int64   `json:"system_data_size"` // GB
	AppDataSize    int64   `json:"app_data_size"`    // GB
	MediaDataSize  int64   `json:"media_data_size"`  // GB
	CacheDataSize  int64   `json:"cache_data_size"`  // GB
	DataGrowthRate float64 `json:"data_growth_rate"` // GB per month

	// App Data Patterns
	AppDataBreakdown string `gorm:"type:json" json:"app_data_breakdown"` // JSON object
	LargestApps      string `gorm:"type:json" json:"largest_apps"`       // JSON array
	FastGrowingApps  string `gorm:"type:json" json:"fast_growing_apps"`  // JSON array

	// Cloud Sync
	CloudSyncEnabled bool       `json:"cloud_sync_enabled"`
	CloudDataSize    int64      `json:"cloud_data_size"` // GB
	SyncFrequency    string     `json:"sync_frequency"`  // realtime, hourly, daily
	LastSyncTime     *time.Time `json:"last_sync_time"`

	// Optimization
	CleanupPotential   int64  `json:"cleanup_potential"` // GB
	DuplicateFiles     int    `json:"duplicate_files"`
	OrphanedFiles      int    `json:"orphaned_files"`
	TempFilesSize      int64  `json:"temp_files_size"`                      // MB
	RecommendedActions string `gorm:"type:json" json:"recommended_actions"` // JSON array

	// Failure Prediction
	FailureRisk       string `json:"failure_risk"`                        // low, medium, high, critical
	PredictedLifespan int    `json:"predicted_lifespan"`                  // months
	WarningIndicators string `gorm:"type:json" json:"warning_indicators"` // JSON array
	BackupRecommended bool   `json:"backup_recommended"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceMemoryManagement tracks RAM usage and optimization
type DeviceMemoryManagement struct {
	database.BaseModel
	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AnalysisTime time.Time `json:"analysis_time"`

	// RAM Configuration
	TotalRAM     int64  `json:"total_ram"`     // MB
	AvailableRAM int64  `json:"available_ram"` // MB
	UsedRAM      int64  `json:"used_ram"`      // MB
	RAMType      string `json:"ram_type"`      // DDR3, DDR4, LPDDR4X, LPDDR5
	RAMSpeed     int    `json:"ram_speed"`     // MHz

	// Usage Patterns
	AverageUsage float64 `json:"average_usage"`                  // percentage
	PeakUsage    float64 `json:"peak_usage"`                     // percentage
	UsageHistory string  `gorm:"type:json" json:"usage_history"` // JSON array
	UsagePattern string  `json:"usage_pattern"`                  // light, moderate, heavy

	// Memory Leaks
	LeaksDetected int    `json:"leaks_detected"`
	LeakingApps   string `gorm:"type:json" json:"leaking_apps"` // JSON array
	MemoryLost    int64  `json:"memory_lost"`                   // MB
	LeakSeverity  string `json:"leak_severity"`                 // low, medium, high

	// App Memory Usage
	AppMemoryBreakdown  string `gorm:"type:json" json:"app_memory_breakdown"` // JSON object
	TopMemoryApps       string `gorm:"type:json" json:"top_memory_apps"`      // JSON array
	BackgroundAppMemory int64  `json:"background_app_memory"`                 // MB
	ForegroundAppMemory int64  `json:"foreground_app_memory"`                 // MB

	// Background Processes
	BackgroundProcesses   int   `json:"background_processes"`
	BackgroundMemoryUsage int64 `json:"background_memory_usage"` // MB
	KilledProcesses       int   `json:"killed_processes"`
	ProcessRestarts       int   `json:"process_restarts"`

	// Memory Optimization
	OptimizationEvents int        `json:"optimization_events"`
	LastOptimization   *time.Time `json:"last_optimization"`
	MemoryReclaimed    int64      `json:"memory_reclaimed"`    // MB
	OptimizationImpact float64    `json:"optimization_impact"` // percentage improvement

	// Low Memory Events
	LowMemoryIncidents   int        `json:"low_memory_incidents"`
	CriticalMemoryEvents int        `json:"critical_memory_events"`
	LastLowMemory        *time.Time `json:"last_low_memory"`
	AppCrashesFromMemory int        `json:"app_crashes_from_memory"`

	// Memory Pressure
	CurrentPressure  string `json:"current_pressure"`                  // none, low, medium, high, critical
	PressureHistory  string `gorm:"type:json" json:"pressure_history"` // JSON array
	PressureDuration int    `json:"pressure_duration"`                 // minutes
	PressureResponse string `json:"pressure_response"`                 // actions taken

	// Swap Usage
	SwapEnabled    bool    `json:"swap_enabled"`
	SwapSize       int64   `json:"swap_size"`       // MB
	SwapUsed       int64   `json:"swap_used"`       // MB
	SwapEfficiency float64 `json:"swap_efficiency"` // percentage

	// Performance Tips
	OptimizationTips string  `gorm:"type:json" json:"optimization_tips"` // JSON array
	RecommendedRAM   int64   `json:"recommended_ram"`                    // MB
	UpgradeNeeded    bool    `json:"upgrade_needed"`
	MemoryScore      float64 `json:"memory_score"` // 0-100

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceProcessorHealth monitors CPU performance and health
type DeviceProcessorHealth struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	MeasurementTime time.Time `json:"measurement_time"`

	// Processor Info
	ProcessorModel string `json:"processor_model"`
	CoreCount      int    `json:"core_count"`
	ThreadCount    int    `json:"thread_count"`
	Architecture   string `json:"architecture"` // ARM64, x86_64
	ProcessNode    int    `json:"process_node"` // nm

	// Throttling Events
	ThrottleCount       int        `json:"throttle_count"`
	ThermalThrottleTime int        `json:"thermal_throttle_time"` // minutes
	PowerThrottleTime   int        `json:"power_throttle_time"`   // minutes
	LastThrottle        *time.Time `json:"last_throttle"`
	ThrottleReason      string     `json:"throttle_reason"`

	// Core Utilization
	CoreUsagePattern string  `gorm:"type:json" json:"core_usage_pattern"` // JSON object
	BigCoreUsage     float64 `json:"big_core_usage"`                      // percentage
	LittleCoreUsage  float64 `json:"little_core_usage"`                   // percentage
	CoreBalancing    float64 `json:"core_balancing"`                      // efficiency score

	// Clock Speed
	CurrentClockSpeed float64 `json:"current_clock_speed"`              // GHz
	MaxClockSpeed     float64 `json:"max_clock_speed"`                  // GHz
	MinClockSpeed     float64 `json:"min_clock_speed"`                  // GHz
	AverageClockSpeed float64 `json:"average_clock_speed"`              // GHz
	ClockVariation    string  `gorm:"type:json" json:"clock_variation"` // JSON array

	// Thermal Management
	ThermalZone     float64 `json:"thermal_zone"`   // celsius
	ThermalPolicy   string  `json:"thermal_policy"` // performance, balanced, quiet
	CoolingState    int     `json:"cooling_state"`
	ThermalHeadroom float64 `json:"thermal_headroom"` // celsius

	// Power Efficiency
	PowerConsumption float64 `json:"power_consumption"` // watts
	PowerEfficiency  float64 `json:"power_efficiency"`  // performance per watt
	IdlePower        float64 `json:"idle_power"`        // watts
	PeakPower        float64 `json:"peak_power"`        // watts

	// Benchmark Scores
	SingleCoreScore  int    `json:"single_core_score"`
	MultiCoreScore   int    `json:"multi_core_score"`
	ScoreTrend       string `json:"score_trend"`                        // improving, stable, declining
	BenchmarkHistory string `gorm:"type:json" json:"benchmark_history"` // JSON array

	// AI/ML Processing
	AIAccelerator     bool    `json:"ai_accelerator"`
	NPUUsage          float64 `json:"npu_usage"`        // percentage
	AIWorkloadTime    float64 `json:"ai_workload_time"` // hours
	MLOperationsCount int64   `json:"ml_operations_count"`

	// Graphics Processing
	GPUIntegrated    bool    `json:"gpu_integrated"`
	GPUUsage         float64 `json:"gpu_usage"`          // percentage
	GraphicsWorkload float64 `json:"graphics_workload"`  // hours
	VideoDecodeUsage float64 `json:"video_decode_usage"` // percentage

	// Aging Effects
	ProcessorAge           int     `json:"processor_age"`                         // months
	PerformanceDegradation float64 `json:"performance_degradation"`               // percentage
	EstimatedLifespan      int     `json:"estimated_lifespan"`                    // months
	DegradationTimeline    string  `gorm:"type:json" json:"degradation_timeline"` // JSON array

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DevicePerformanceMetrics
func (dpm *DevicePerformanceMetrics) IsPerformanceGood() bool {
	return dpm.OverallScore >= 75 && dpm.ResponsivenessScore >= 70 && dpm.AppCrashCount < 5
}

func (dpm *DevicePerformanceMetrics) NeedsOptimization() bool {
	return dpm.PerformanceTrend == "degrading" || dpm.OverallScore < 60 ||
		dpm.MemoryPressure == "high" || dpm.UILagIncidents > 10
}

func (dpm *DevicePerformanceMetrics) HasStorageIssues() bool {
	remainingStorage := dpm.StorageTotal - dpm.StorageUsed
	percentUsed := float64(dpm.StorageUsed) / float64(dpm.StorageTotal) * 100
	return percentUsed > 90 || remainingStorage < 2 // Less than 2GB free
}

func (dpm *DevicePerformanceMetrics) GetPerformanceCategory() string {
	if dpm.OverallScore >= 85 {
		return "excellent"
	} else if dpm.OverallScore >= 70 {
		return "good"
	} else if dpm.OverallScore >= 50 {
		return "fair"
	}
	return "poor"
}

func (dpm *DevicePerformanceMetrics) IsGamingCapable() bool {
	return dpm.GamingFPS >= 30 && dpm.GraphicsScore >= 60 && !dpm.ThermalThrottling
}

// Methods for DeviceBatteryAnalytics
func (dba *DeviceBatteryAnalytics) IsBatteryHealthy() bool {
	return dba.BatteryHealth >= 80 && dba.DegradationRate < 2.0 // Less than 2% per month
}

func (dba *DeviceBatteryAnalytics) NeedsReplacement() bool {
	return dba.ReplacementNeeded || dba.BatteryHealth < 70 ||
		dba.ReplacementUrgency == "high" || dba.ReplacementUrgency == "critical"
}

func (dba *DeviceBatteryAnalytics) GetBatteryLifeCategory() string {
	screenTime := dba.ScreenOnTime
	if screenTime >= 8 {
		return "excellent"
	} else if screenTime >= 6 {
		return "good"
	} else if screenTime >= 4 {
		return "average"
	}
	return "poor"
}

func (dba *DeviceBatteryAnalytics) HasChargingIssues() bool {
	return dba.FastChargingImpact > 5.0 || dba.WirelessEfficiency < 70 ||
		dba.ChargingTemperature > 40
}

func (dba *DeviceBatteryAnalytics) GetRemainingLifespan() int {
	if dba.EstimatedLifespan > 0 {
		return dba.EstimatedLifespan
	}
	// Estimate based on health and degradation rate
	if dba.DegradationRate > 0 {
		remainingHealth := dba.BatteryHealth - 70 // 70% is replacement threshold
		monthsRemaining := int(remainingHealth / dba.DegradationRate)
		if monthsRemaining < 0 {
			return 0
		}
		return monthsRemaining
	}
	return 24 // Default 2 years
}

// Methods for DeviceTemperatureHistory
func (dth *DeviceTemperatureHistory) HasOverheatingProblem() bool {
	return dth.OverheatingIncidents > 5 || dth.ThermalShutdowns > 0 ||
		dth.PeakTemperature > 50
}

func (dth *DeviceTemperatureHistory) IsThrottling() bool {
	return dth.ThrottlingEvents > 0 && dth.PerformanceImpact > 10
}

func (dth *DeviceTemperatureHistory) GetThermalStatus() string {
	if dth.AverageTemperature < 35 {
		return "cool"
	} else if dth.AverageTemperature < 40 {
		return "normal"
	} else if dth.AverageTemperature < 45 {
		return "warm"
	}
	return "hot"
}

func (dth *DeviceTemperatureHistory) NeedsCoolingImprovement() bool {
	return dth.CoolingEfficiency < 70 || dth.ThermalShutdowns > 0 ||
		dth.UserComfortScore < 60
}

func (dth *DeviceTemperatureHistory) GetHighestUsageTemp() float64 {
	temps := []float64{
		dth.ChargingPeak,
		dth.GamingPeak,
		dth.VideoPeak,
		dth.NavigationPeak,
		dth.MultitaskingPeak,
	}

	max := temps[0]
	for _, temp := range temps {
		if temp > max {
			max = temp
		}
	}
	return max
}

// Methods for DeviceStorageHealth
func (dsh *DeviceStorageHealth) IsHealthy() bool {
	return dsh.HealthStatus >= 90 && dsh.WearLevel < 20 && dsh.BadSectorCount == 0
}

func (dsh *DeviceStorageHealth) NeedsAttention() bool {
	return dsh.FailureRisk == "high" || dsh.FailureRisk == "critical" ||
		dsh.BackupRecommended || dsh.BadSectorCount > 5
}

func (dsh *DeviceStorageHealth) GetSpaceEfficiency() float64 {
	totalData := dsh.SystemDataSize + dsh.AppDataSize + dsh.MediaDataSize
	if totalData > 0 {
		usefulData := totalData - dsh.CacheDataSize
		return float64(usefulData) / float64(totalData) * 100
	}
	return 100
}

func (dsh *DeviceStorageHealth) HasPerformanceDegradation() bool {
	return dsh.SpeedDegradation > 20 || dsh.FragmentationLevel > 30 ||
		dsh.CacheHitRate < 70
}

func (dsh *DeviceStorageHealth) GetCleanupRecommendation() string {
	if dsh.CleanupPotential > 10 {
		return "high_priority"
	} else if dsh.CleanupPotential > 5 {
		return "recommended"
	} else if dsh.CleanupPotential > 2 {
		return "optional"
	}
	return "not_needed"
}

// Methods for DeviceMemoryManagement
func (dmm *DeviceMemoryManagement) HasMemoryPressure() bool {
	return dmm.CurrentPressure == "high" || dmm.CurrentPressure == "critical" ||
		dmm.LowMemoryIncidents > 10
}

func (dmm *DeviceMemoryManagement) GetMemoryUsagePercent() float64 {
	if dmm.TotalRAM > 0 {
		return float64(dmm.UsedRAM) / float64(dmm.TotalRAM) * 100
	}
	return 0
}

func (dmm *DeviceMemoryManagement) HasMemoryLeaks() bool {
	return dmm.LeaksDetected > 0 && (dmm.LeakSeverity == "medium" || dmm.LeakSeverity == "high")
}

func (dmm *DeviceMemoryManagement) NeedsMoreRAM() bool {
	return dmm.UpgradeNeeded || dmm.CriticalMemoryEvents > 5 ||
		dmm.AppCrashesFromMemory > 3 || dmm.GetMemoryUsagePercent() > 90
}

func (dmm *DeviceMemoryManagement) GetMemoryEfficiency() float64 {
	if dmm.MemoryScore > 0 {
		return dmm.MemoryScore
	}
	// Calculate based on usage and optimization
	efficiency := 100.0
	efficiency -= float64(dmm.LeaksDetected * 5)
	efficiency -= float64(dmm.LowMemoryIncidents * 2)
	efficiency -= (dmm.GetMemoryUsagePercent() - 70) // Penalty for usage over 70%

	if efficiency < 0 {
		efficiency = 0
	}
	return efficiency
}

// Methods for DeviceProcessorHealth
func (dph *DeviceProcessorHealth) IsThrottling() bool {
	return dph.ThrottleCount > 0 || dph.ThermalThrottleTime > 0 || dph.PowerThrottleTime > 0
}

func (dph *DeviceProcessorHealth) GetPerformanceRatio() float64 {
	if dph.MaxClockSpeed > 0 {
		return (dph.AverageClockSpeed / dph.MaxClockSpeed) * 100
	}
	return 100
}

func (dph *DeviceProcessorHealth) HasAging() bool {
	return dph.PerformanceDegradation > 10 || dph.ProcessorAge > 36 // More than 3 years
}

func (dph *DeviceProcessorHealth) IsHighPerformance() bool {
	return dph.SingleCoreScore > 1000 && dph.MultiCoreScore > 4000 &&
		dph.PowerEfficiency > 80
}

func (dph *DeviceProcessorHealth) GetThermalHeadroom() float64 {
	if dph.ThermalHeadroom > 0 {
		return dph.ThermalHeadroom
	}
	// Estimate based on current temp and throttle temp (assuming throttle at 85°C)
	throttleTemp := 85.0
	return throttleTemp - dph.ThermalZone
}
