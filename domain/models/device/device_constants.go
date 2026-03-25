package device

// ============================================
// COMMON DEVICE CONSTANTS
// These apply to all device categories
// ============================================

// Device Status Constants
const (
	DeviceStatusActive   DeviceStatus = "active"
	DeviceStatusInactive DeviceStatus = "inactive"
	DeviceStatusStolen   DeviceStatus = "stolen"
	DeviceStatusLost     DeviceStatus = "lost"
	DeviceStatusDamaged  DeviceStatus = "damaged"
	DeviceStatusRetired  DeviceStatus = "retired"
)

// Device Condition Constants
const (
	DeviceConditionExcellent DeviceCondition = "excellent"
	DeviceConditionLikeNew   DeviceCondition = "like_new"
	DeviceConditionGood      DeviceCondition = "good"
	DeviceConditionFair      DeviceCondition = "fair"
	DeviceConditionPoor      DeviceCondition = "poor"
	DeviceConditionBroken    DeviceCondition = "broken"
)

// Device Grade Constants
const (
	DeviceGradeA      DeviceGrade = "A"
	DeviceGradeAMinus DeviceGrade = "A-"
	DeviceGradeBPlus  DeviceGrade = "B+"
	DeviceGradeB      DeviceGrade = "B"
	DeviceGradeBMinus DeviceGrade = "B-"
	DeviceGradeCPlus  DeviceGrade = "C+"
	DeviceGradeC      DeviceGrade = "C"
	DeviceGradeD      DeviceGrade = "D"
	DeviceGradeF      DeviceGrade = "F"
)

// Device Segment Constants
const (
	DeviceSegmentFlagship = "flagship"
	DeviceSegmentPremium  = "premium"
	DeviceSegmentMidRange = "mid_range"
	DeviceSegmentBudget   = "budget"
)

// Screen Condition Constants
const (
	ScreenConditionPerfect        ScreenCondition = "perfect"
	ScreenConditionMinorScratches ScreenCondition = "minor_scratches"
	ScreenConditionCracked        ScreenCondition = "cracked"
	ScreenConditionBroken         ScreenCondition = "broken"
)

// Body Condition Constants
const (
	BodyConditionPerfect   BodyCondition = "perfect"
	BodyConditionMinorWear BodyCondition = "minor_wear"
	BodyConditionDamaged   BodyCondition = "damaged"
)

// Theft Risk Level Constants
const (
	TheftRiskLow      = "low"
	TheftRiskMedium   = "medium"
	TheftRiskHigh     = "high"
	TheftRiskVeryHigh = "very_high"
)

// Blacklist Status Constants
const (
	BlacklistStatusClean    BlacklistStatus = "clean"
	BlacklistStatusBlocked  BlacklistStatus = "blocked"
	BlacklistStatusChecking BlacklistStatus = "checking"
)

// Authenticity Status Constants
const (
	AuthenticityStatusVerified   AuthenticityStatus = "verified"
	AuthenticityStatusUnverified AuthenticityStatus = "unverified"
	AuthenticityStatusChecking   AuthenticityStatus = "checking"
	AuthenticityStatusFailed     AuthenticityStatus = "failed"
	AuthenticityStatusFake       AuthenticityStatus = "fake"
)

// Network Status Constants
const (
	NetworkStatusLocked   NetworkStatus = "locked"
	NetworkStatusUnlocked NetworkStatus = "unlocked"
)

// Water Damage Indicator Constants
const (
	WaterDamageWhite = "white"
	WaterDamageRed   = "red"
	WaterDamagePink  = "pink"
)

// Component Condition Constants
const (
	ComponentWorking      = "working"
	ComponentIntermittent = "intermittent"
	ComponentDamaged      = "damaged"
	ComponentNotWorking   = "not_working"
)

// Camera Condition Constants
const (
	CameraAllWorking = "all_working"
	CameraFrontIssue = "front_issue"
	CameraRearIssue  = "rear_issue"
	CameraBothIssue  = "both_issue"
)

// Ownership Type Constants
const (
	OwnershipTypePersonal  = "personal"
	OwnershipTypeCorporate = "corporate"
	OwnershipTypeBYOD      = "byod"
)

// Assignment Type Constants
const (
	AssignmentTypePermanent = "permanent"
	AssignmentTypeTemporary = "temporary"
	AssignmentTypeLoaner    = "loaner"
)

// Parts Availability Constants
const (
	PartsReadilyAvailable = "readily_available"
	PartsLimited          = "limited"
	PartsScarce           = "scarce"
	PartsDiscontinued     = "discontinued"
)

// ============================================
// DEFAULT VALUES (COMMON)
// ============================================

const (
	DefaultCurrency               = "USD"
	DefaultRiskScoreThreshold     = 50.0
	DefaultHighValueThreshold     = 1000.0
	DefaultPremiumValueThreshold  = 1500.0
	DefaultMinInsurableValue      = 100.0
	DefaultFlagshipMinValue       = 200.0
	DefaultInspectionDays         = 180
	DefaultMaxDeviceAgeDays       = 1460 // 4 years
	DefaultMaxNonFlagshipAge      = 1095 // 3 years
	DefaultBaseInsuranceRate      = 0.05 // 5% per year
	DefaultMinMonthlyPremium      = 5.0
	DefaultMaxMonthlyPremium      = 100.0
	DefaultBatteryHealthThreshold = 80
)

// Depreciation Rates by Segment
const (
	DepreciationRateFlagship = 0.18 // 18% per year
	DepreciationRatePremium  = 0.20 // 20% per year
	DepreciationRateMidRange = 0.25 // 25% per year
	DepreciationRateBudget   = 0.30 // 30% per year
	DepreciationMinFlagship  = 0.15 // 15% minimum value
	DepreciationMinDefault   = 0.10 // 10% minimum value
)

// Risk Score Weights
const (
	RiskScoreMaxAge        = 25.0
	RiskScoreMaxCondition  = 25.0
	RiskScoreMaxTheft      = 20.0
	RiskScoreMaxValue      = 15.0
	RiskScoreMaxRepairs    = 10.0
	RiskScoreSecurityBonus = 5.0
	RiskScoreMaxTotal      = 100.0
)

// Condition Score Multipliers
const (
	ConditionMultiplierExcellent = 1.0
	ConditionMultiplierGood      = 0.85
	ConditionMultiplierFair      = 0.7
	ConditionMultiplierPoor      = 0.5
)

// Trade-In Value Multipliers
const (
	TradeInMultiplierExcellent = 0.9
	TradeInMultiplierGood      = 0.75
	TradeInMultiplierFair      = 0.6
	TradeInMultiplierPoor      = 0.4
	TradeInScreenCracked       = 0.7
	TradeInScreenBroken        = 0.5
	TradeInNoBoxPenalty        = 0.95
	TradeInNoReceiptPenalty    = 0.95
	TradeInMinValue            = 50.0
)

// ============================================
// COMMON USER PROFILES
// ============================================

const (
	UserProfileCareful      = "careful"
	UserProfileAverage      = "average"
	UserProfileCareless     = "careless"
	UserProfileProfessional = "professional"
	UserProfileStudent      = "student"
	UserProfileSenior       = "senior"
	UserProfileChild        = "child"
	UserProfileAthlete      = "athlete"
	UserProfileFieldWorker  = "field_worker"
)

// ============================================
// COMMON WATER RESISTANCE RATINGS
// ============================================

// Water Resistance Constants
const (
	WaterResistanceNone  WaterResistance = "none"
	WaterResistanceIPX4  WaterResistance = "IPX4"  // Splash resistant
	WaterResistanceIPX7  WaterResistance = "IPX7"  // Immersion up to 1m
	WaterResistanceIP67  WaterResistance = "IP67"  // Dust tight, 1m water
	WaterResistanceIP68  WaterResistance = "IP68"  // Dust tight, 1.5m+ water
	WaterResistanceIP69K WaterResistance = "IP69K" // High pressure/temperature
	WaterResistance3ATM  WaterResistance = "3ATM"  // 30m pressure (watches)
	WaterResistance5ATM  WaterResistance = "5ATM"  // 50m pressure (watches)
	WaterResistance10ATM WaterResistance = "10ATM" // 100m pressure (watches)
)

// ============================================
// COMMON CONNECTIVITY
// ============================================

const (
	// Connectivity types
	ConnectivityBluetooth = "bluetooth"
	ConnectivityWiFi      = "wifi"
	ConnectivityNFC       = "nfc"
	ConnectivityGPS       = "gps"
	ConnectivityCellular  = "cellular"

	// Bluetooth versions (common)
	BluetoothVersion40 = "4.0"
	BluetoothVersion42 = "4.2"
	BluetoothVersion50 = "5.0"
	BluetoothVersion51 = "5.1"
	BluetoothVersion52 = "5.2"
	BluetoothVersion53 = "5.3"

	// Network generations
	Network2G  = "2G"
	Network3G  = "3G"
	Network4G  = "4G"
	NetworkLTE = "LTE"
	Network5G  = "5G"
)
