package smartwatch

// ============================================
// SMARTWATCH SPECIFICATION CONSTANTS
// ============================================

// Screen Size Categories (in mm)
const (
	ScreenSizeSmall  = 38.0 // 38mm and below
	ScreenSizeMedium = 41.0 // 40-41mm
	ScreenSizeLarge  = 44.0 // 42-45mm
	ScreenSizeXLarge = 49.0 // 46mm and above
)

// Screen Types
const (
	ScreenTypeLCD    = "LCD"
	ScreenTypeOLED   = "OLED"
	ScreenTypeAMOLED = "AMOLED"
	ScreenTypeLTPO   = "LTPO"
	ScreenTypeRetina = "Retina"
	ScreenTypeEInk   = "E-Ink"
)

// Screen Protection Types
const (
	ScreenProtectionSapphireCrystal = "Sapphire Crystal"
	ScreenProtectionIonXGlass       = "Ion-X Glass"
	ScreenProtectionGorillaGlass    = "Gorilla Glass"
	ScreenProtectionMineralGlass    = "Mineral Glass"
	ScreenProtectionPlastic         = "Plastic"
)

// ============================================
// HEALTH SENSOR CONSTANTS
// ============================================

// Health Features
const (
	HealthFeatureSleepTracking    = "sleep_tracking"
	HealthFeatureStressMonitoring = "stress_monitoring"
	HealthFeatureCycleTracking    = "cycle_tracking"
	HealthFeatureFallDetection    = "fall_detection"
	HealthFeatureAFibDetection    = "afib_detection"
	HealthFeatureBloodPressure    = "blood_pressure"
	HealthFeatureRespiratoryRate  = "respiratory_rate"
	HealthFeatureHRV              = "heart_rate_variability"
	HealthFeatureVO2Max           = "vo2_max"
	HealthFeatureBodyTemperature  = "body_temperature"
	HealthFeatureBloodGlucose     = "blood_glucose"
	HealthFeatureHydration        = "hydration_tracking"
)

// NOTE: Common biometric types (pin, pattern) are defined in device_constants.go
// Smartwatch-specific biometric features:
const (
	BiometricWristDetection = "wrist_detection"    // Auto lock when removed from wrist
	BiometricECGAuth        = "ecg_authentication" // ECG-based authentication
)

// ============================================
// FITNESS TRACKING CONSTANTS
// ============================================

// Activity Types
const (
	ActivityRunning          = "running"
	ActivityCycling          = "cycling"
	ActivitySwimming         = "swimming"
	ActivityYoga             = "yoga"
	ActivityHiking           = "hiking"
	ActivityStrengthTraining = "strength_training"
	ActivityDancing          = "dancing"
	ActivityRowing           = "rowing"
	ActivityGolf             = "golf"
	ActivityTennis           = "tennis"
	ActivitySkiing           = "skiing"
	ActivitySurfing          = "surfing"
	ActivityClimbing         = "climbing"
	ActivityMartialArts      = "martial_arts"
)

// Water Resistance Ratings (Smartwatch-specific)
// NOTE: Common water resistance ratings (IP67, IP68, 3ATM, 5ATM, 10ATM) are defined in device_constants.go
const (
	WaterResistance20ATM     = "20ATM"      // High-speed water sports (watch-specific)
	WaterResistanceSwimProof = "Swim Proof" // Watch-specific terminology
	WaterResistanceDiveProof = "Dive Proof" // Watch-specific for dive computers
)

// ============================================
// BATTERY CONSTANTS
// ============================================

// Battery Capacity Ranges (mAh)
const (
	BatteryCapacityVeryLow  = 200 // Below 200mAh
	BatteryCapacityLow      = 300 // 200-300mAh
	BatteryCapacityMedium   = 400 // 300-400mAh
	BatteryCapacityHigh     = 500 // 400-500mAh
	BatteryCapacityVeryHigh = 600 // Above 500mAh
)

// Battery Life Categories (days)
const (
	BatteryLifeLessThanDay = 0.5 // Less than 1 day
	BatteryLifeOneDay      = 1   // 1 day
	BatteryLifeTwoDays     = 2   // 2 days
	BatteryLifeWeek        = 7   // 1 week
	BatteryLifeTwoWeeks    = 14  // 2 weeks
	BatteryLifeMonth       = 30  // 1 month
)

// Charging Methods
const (
	ChargingMethodMagnetic    = "magnetic"
	ChargingMethodWireless    = "wireless"
	ChargingMethodPins        = "pins"
	ChargingMethodUSBC        = "usb-c"
	ChargingMethodProprietary = "proprietary"
	ChargingMethodSolar       = "solar"
	ChargingMethodKinetic     = "kinetic"
)

// Charging Speed Categories (minutes to full)
const (
	ChargingSpeedVeryFast = 30  // Under 30 minutes
	ChargingSpeedFast     = 60  // 30-60 minutes
	ChargingSpeedNormal   = 90  // 60-90 minutes
	ChargingSpeedSlow     = 120 // 90-120 minutes
	ChargingSpeedVerySlow = 180 // Over 120 minutes
)

// ============================================
// CONNECTIVITY CONSTANTS
// ============================================

// Connectivity Types (Smartwatch-specific)
// NOTE: Common connectivity types (bluetooth, wifi, nfc, gps) are defined in device_constants.go
const (
	ConnectivityANT = "ant+" // ANT+ is specific to fitness devices/watches
)

// ============================================
// BUILD QUALITY CONSTANTS
// ============================================

// Case Materials
const (
	CaseMaterialAluminum = "aluminum"
	CaseMaterialSteel    = "stainless_steel"
	CaseMaterialTitanium = "titanium"
	CaseMaterialCeramic  = "ceramic"
	CaseMaterialPlastic  = "plastic"
	CaseMaterialResin    = "resin"
	CaseMaterialCarbon   = "carbon_fiber"
	CaseMaterialGold     = "gold"
	CaseMaterialPlatinum = "platinum"
)

// Band Types
const (
	BandTypeSport    = "sport"
	BandTypeLeather  = "leather"
	BandTypeMetal    = "metal"
	BandTypeMilanese = "milanese"
	BandTypeSoloLoop = "solo_loop"
	BandTypeFabric   = "fabric"
	BandTypeSilicone = "silicone"
	BandTypeNylon    = "nylon"
	BandTypeCanvas   = "canvas"
	BandTypeRubber   = "rubber"
)

// Case Size Categories (mm)
const (
	CaseSizeXSmall  = 38 // 38mm and below
	CaseSizeSmall   = 40 // 40-41mm
	CaseSizeMedium  = 42 // 42-43mm
	CaseSizeLarge   = 44 // 44-45mm
	CaseSizeXLarge  = 46 // 46-47mm
	CaseSizeXXLarge = 49 // 49mm and above
)

// Weight Categories (grams)
const (
	WeightUltraLight = 25 // Under 25g
	WeightLight      = 35 // 25-35g
	WeightMedium     = 45 // 35-45g
	WeightHeavy      = 60 // 45-60g
	WeightVeryHeavy  = 80 // Over 60g
)

// ============================================
// OPERATING SYSTEM CONSTANTS
// ============================================

// Operating Systems
const (
	OSWatchOS     = "watchOS"
	OSWearOS      = "Wear OS"
	OSTizen       = "Tizen"
	OSFitbitOS    = "Fitbit OS"
	OSGarminOS    = "Garmin OS"
	OSAmazfit     = "Amazfit OS"
	OSHarmonyOS   = "HarmonyOS"
	OSProprietary = "Proprietary"
)

// Voice Assistants
const (
	VoiceAssistantSiri    = "Siri"
	VoiceAssistantGoogle  = "Google Assistant"
	VoiceAssistantAlexa   = "Alexa"
	VoiceAssistantBixby   = "Bixby"
	VoiceAssistantCortana = "Cortana"
	VoiceAssistantNone    = "None"
)

// ============================================
// PROCESSOR CONSTANTS
// ============================================

// Processor Types
const (
	ProcessorAppleS8        = "Apple S8"
	ProcessorAppleS9        = "Apple S9"
	ProcessorSnapdragonW5   = "Snapdragon W5+"
	ProcessorSnapdragon4100 = "Snapdragon 4100+"
	ProcessorExynos9110     = "Exynos 9110"
	ProcessorExynos920      = "Exynos W920"
	ProcessorMediatekMT2601 = "Mediatek MT2601"
	ProcessorKirin          = "Kirin A1"
)

// ============================================
// PRICE TIER CONSTANTS
// ============================================

// Price Categories (USD)
const (
	PriceTierBudget      = 100.0   // Under $100
	PriceTierLowEnd      = 200.0   // $100-200
	PriceTierMidRange    = 350.0   // $200-350
	PriceTierPremium     = 500.0   // $350-500
	PriceTierHighEnd     = 800.0   // $500-800
	PriceTierLuxury      = 1500.0  // $800-1500
	PriceTierUltraLuxury = 10000.0 // Above $1500
)

// ============================================
// RISK ASSESSMENT CONSTANTS
// ============================================

// Risk Factors
const (
	RiskFactorSportUsage     = 1.3 // Higher risk for sports watches
	RiskFactorDailyWear      = 1.0 // Normal wear risk
	RiskFactorOccasionalWear = 0.7 // Lower risk for occasional use
	RiskFactorLuxuryMaterial = 0.8 // Premium materials are more durable
	RiskFactorPlasticBuild   = 1.2 // Plastic builds have higher damage risk
	RiskFactorWaterSports    = 1.5 // Water activities increase risk
	RiskFactorChildUser      = 1.4 // Kids' watches have higher damage risk
)

// Depreciation Rates (annual percentage)
const (
	DepreciationRateLuxury   = 0.15 // 15% per year for luxury watches
	DepreciationRatePremium  = 0.20 // 20% per year for premium watches
	DepreciationRateMidRange = 0.25 // 25% per year for mid-range
	DepreciationRateBudget   = 0.35 // 35% per year for budget watches
	DepreciationRateSport    = 0.30 // 30% per year for sport watches
)

// ============================================
// REPAIR COST CONSTANTS
// ============================================

// Repair Cost Percentages (of device value)
const (
	RepairCostScreenReplacement    = 0.40 // 40% of device value
	RepairCostBatteryReplacement   = 0.25 // 25% of device value
	RepairCostCrownReplacement     = 0.15 // 15% of device value
	RepairCostSensorRepair         = 0.30 // 30% of device value
	RepairCostWaterDamage          = 0.60 // 60% of device value
	RepairCostMotherboardRepair    = 0.70 // 70% of device value
	RepairCostBandReplacement      = 0.10 // 10% of device value
	RepairCostBackCoverReplacement = 0.20 // 20% of device value
)

// ============================================
// VALIDATION CONSTANTS
// ============================================

// Validation Limits
const (
	MinBatteryCapacity = 100  // Minimum 100mAh
	MaxBatteryCapacity = 1000 // Maximum 1000mAh
	MinScreenSize      = 20.0 // Minimum 20mm
	MaxScreenSize      = 60.0 // Maximum 60mm
	MinWeight          = 10   // Minimum 10 grams
	MaxWeight          = 200  // Maximum 200 grams
	MinStorageCapacity = 0    // Some watches have no storage
	MaxStorageCapacity = 64   // Maximum 64GB
	MinRAM             = 256  // Minimum 256MB
	MaxRAM             = 4096 // Maximum 4GB
)

// ============================================
// FEATURE FLAGS
// ============================================

// Premium Features
const (
	FeatureECG              = "ecg_sensor"
	FeatureBloodOxygen      = "blood_oxygen"
	FeatureAlwaysOnDisplay  = "always_on_display"
	FeatureCellular         = "cellular"
	FeatureeSIM             = "esim"
	FeatureSapphireGlass    = "sapphire_glass"
	FeatureTitaniumBuild    = "titanium_build"
	FeatureWirelessCharging = "wireless_charging"
	FeatureSolarCharging    = "solar_charging"
	FeatureGolfMode         = "golf_mode"
	FeatureDiveMode         = "dive_mode"
	FeatureMusicStorage     = "music_storage"
)
