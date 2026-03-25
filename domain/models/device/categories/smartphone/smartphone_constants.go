package smartphone

// ============================================
// SMARTPHONE-SPECIFIC CONSTANTS
// These constants are unique to smartphones
// Common constants are in device_constants.go
// ============================================
const (
	// Screen sizes
	MinScreenSize = 3.0
	MaxScreenSize = 8.0

	// Storage limits
	MinStorage = 8    // GB
	MaxStorage = 2048 // GB

	// RAM limits
	MinRAM = 1  // GB
	MaxRAM = 32 // GB

	// Battery capacity
	MinBatteryCapacity = 1000  // mAh
	MaxBatteryCapacity = 10000 // mAh

	// Camera limits
	MinMegapixels = 0.3
	MaxMegapixels = 200.0

	// Charging speeds
	StandardCharging  = 5.0  // Watts
	FastCharging      = 18.0 // Watts
	SuperFastCharging = 30.0 // Watts

	// Risk multipliers
	ScreenRiskMultiplier  = 0.35
	BatteryRiskMultiplier = 0.15
	WaterDamageMultiplier = 0.20
	TheftRiskMultiplier   = 0.15
	DropRiskMultiplier    = 0.25

	// Insurance rates
	SmartphoneBasePremiumRate    = 0.08 // 8% annually
	SmartphoneCategoryMultiplier = 1.2  // 20% higher than base

	// Depreciation
	SmartphoneDepreciationRate = 0.30 // 30% per year
	MinResaleValuePercent      = 0.10 // 10% of original value

	// Coverage limits
	MaxClaimsPerYear     = 3
	MinDeductiblePercent = 0.05 // 5% of device value
	MaxDeductiblePercent = 0.15 // 15% of device value

	// Repair cost limits
	MaxScreenRepairCost     = 500.0
	MaxBatteryRepairCost    = 150.0
	MaxCameraRepairCost     = 300.0
	MaxChargingPortCost     = 100.0
	MaxMotherboardCostRatio = 0.6 // 60% of device value

	// Age limits
	MaxInsurableAgeMonths  = 60 // 5 years
	NewDeviceAgeMonths     = 1  // 1 month
	WarrantyPeriodMonths   = 12 // Standard warranty
	ExtendedWarrantyMonths = 24 // Extended warranty
)

// Display types (Smartphone-specific)
const (
	DisplayTypeLCD         = "LCD"
	DisplayTypeOLED        = "OLED"
	DisplayTypeAMOLED      = "AMOLED"
	DisplayTypeRetina      = "Retina"
	DisplayTypeSuperAMOLED = "Super AMOLED"
	DisplayTypeLTPO        = "LTPO"
	DisplayTypeIPS         = "IPS"
)

// Operating systems (Smartphone-specific)
const (
	OSiOS           = "iOS"
	OSAndroid       = "Android"
	OSWindowsPhone  = "Windows Phone"
	OSBlackBerry    = "BlackBerry OS"
	OSKaiOS         = "KaiOS"
	OSHarmonyOS     = "HarmonyOS"
	OSCustomAndroid = "Custom Android" // LineageOS, etc.
)

// NOTE: Common biometric types (fingerprint, face_id, iris, pin, pattern) are defined in device_constants.go
// Smartphone-specific biometric features:
const (
	BiometricInDisplayFingerprint  = "in_display_fingerprint"
	BiometricUltrasonicFingerprint = "ultrasonic_fingerprint"
	BiometricVoice                 = "voice" // Voice unlock (smartphone-specific)
)

// Special features (Smartphone-specific)
// NOTE: Common connectivity features (5G, NFC) are referenced via Network5G and ConnectivityNFC in device_constants.go
const (
	FeatureFoldable         = "foldable"
	FeatureFlip             = "flip"
	FeatureStylus           = "stylus"
	FeatureDualScreen       = "dual_screen"
	FeatureWirelessCharging = "wireless_charging"
	FeatureReverseWireless  = "reverse_wireless_charging"
	FeatureIRBlaster        = "ir_blaster"
	FeatureHeadphoneJack    = "headphone_jack"
)

// Camera types (Smartphone-specific)
const (
	CameraWide       = "wide"
	CameraUltraWide  = "ultrawide"
	CameraTelephoto  = "telephoto"
	CameraMacro      = "macro"
	CameraDepth      = "depth"
	CameraMonochrome = "monochrome"
	CameraToF        = "tof" // Time of Flight
	CameraPeriscope  = "periscope"
)

// Protection levels
const (
	ProtectionNone   = "none"
	ProtectionLow    = "low"
	ProtectionMedium = "medium"
	ProtectionHigh   = "high"
)

// Environmental risk levels
const (
	EnvironmentalRiskLow    = "low"
	EnvironmentalRiskMedium = "medium"
	EnvironmentalRiskHigh   = "high"
)

// Manufacturers
var SupportedManufacturers = []string{
	"Apple",
	"Samsung",
	"Google",
	"OnePlus",
	"Xiaomi",
	"Huawei",
	"Sony",
	"LG",
	"Motorola",
	"Nokia",
	"Asus",
	"Oppo",
	"Vivo",
	"Realme",
	"Honor",
	"ZTE",
	"BlackBerry",
	"HTC",
	"Lenovo",
	"TCL",
}

// Premium manufacturers (retain value better)
var PremiumManufacturers = []string{
	"Apple",
	"Samsung",
	"Google",
}

// NOTE: Water resistance ratings are defined in device_constants.go as common constants

// Screen protection types
var ScreenProtectionTypes = []string{
	"None",
	"Gorilla Glass",
	"Gorilla Glass 2",
	"Gorilla Glass 3",
	"Gorilla Glass 4",
	"Gorilla Glass 5",
	"Gorilla Glass 6",
	"Gorilla Glass Victus",
	"Gorilla Glass Victus+",
	"Gorilla Glass Victus 2",
	"Dragontrail",
	"Sapphire",
	"Ceramic Shield",
}

// Battery types
var BatteryTypes = []string{
	"Li-Ion",
	"Li-Po",
	"Li-Polymer",
	"Graphene",
}

// Charging standards
var ChargingStandards = []string{
	"USB-C",
	"Lightning",
	"Micro-USB",
	"Wireless Qi",
	"MagSafe",
	"Proprietary",
}
