package sparepart

// Part Types
const (
	// Display & Touch
	PartTypeScreen      = "screen"
	PartTypeTouchScreen = "touch_screen"
	PartTypeDigitizer   = "digitizer"
	PartTypeLCD         = "lcd_panel"
	PartTypeOLED        = "oled_panel"
	PartTypeBacklight   = "backlight"

	// Power & Charging
	PartTypeBattery              = "battery"
	PartTypeChargePort           = "charge_port"
	PartTypeWirelessChargingCoil = "wireless_charging_coil"
	PartTypeChargingIC           = "charging_ic"
	PartTypePowerManagementIC    = "power_management_ic"

	// Physical Buttons & Controls
	PartTypeButton       = "button"
	PartTypePowerButton  = "power_button"
	PartTypeVolumeButton = "volume_button"
	PartTypeHomeButton   = "home_button"
	PartTypeMuteSwitch   = "mute_switch"
	PartTypeAlertSlider  = "alert_slider"
	PartTypeSideButton   = "side_button"

	// Camera System
	PartTypeCamera          = "camera"
	PartTypeFrontCamera     = "front_camera"
	PartTypeRearCamera      = "rear_camera"
	PartTypeCameraLens      = "camera_lens"
	PartTypeCameraLensGlass = "camera_lens_glass"
	PartTypeCameraModule    = "camera_module"

	// Audio Components
	PartTypeSpeaker         = "speaker"
	PartTypeEarpieceSpeaker = "earpiece_speaker"
	PartTypeLoudSpeaker     = "loud_speaker"
	PartTypeMicrophone      = "microphone"
	PartTypeNoiseCancelMic  = "noise_cancel_microphone"
	PartTypeAudioJack       = "audio_jack"
	PartTypeAudioIC         = "audio_ic"

	// Sensors
	PartTypeSensor             = "sensor"
	PartTypeProximitySensor    = "proximity_sensor"
	PartTypeAmbientLightSensor = "ambient_light_sensor"
	PartTypeFingerprintSensor  = "fingerprint_sensor"
	PartTypeFaceIDModule       = "face_id_module"
	PartTypeAccelerometer      = "accelerometer"
	PartTypeGyroscope          = "gyroscope"
	PartTypeMagnetometer       = "magnetometer"
	PartTypeBarometer          = "barometer"

	// Antennas & Connectivity
	PartTypeAntenna          = "antenna"
	PartTypeWiFiAntenna      = "wifi_antenna"
	PartTypeBluetoothAntenna = "bluetooth_antenna"
	PartTypeCellularAntenna  = "cellular_antenna"
	PartType5GAntenna        = "5g_antenna"
	PartTypeNFCAntenna       = "nfc_antenna"
	PartTypeGPSAntenna       = "gps_antenna"

	// Flex Cables
	PartTypeFlex             = "flex_cable"
	PartTypePowerButtonFlex  = "power_button_flex"
	PartTypeVolumeButtonFlex = "volume_button_flex"
	PartTypeDisplayFlex      = "display_flex"
	PartTypeChargingFlex     = "charging_flex"

	// Internal Components
	PartTypeMotherboard = "motherboard"
	PartTypeLogicBoard  = "logic_board"
	PartTypeMemory      = "memory"
	PartTypeStorage     = "storage"
	PartTypeRAM         = "ram"
	PartTypeNAND        = "nand_flash"

	// ICs and Chips
	PartTypeTouchControllerIC = "touch_controller_ic"
	PartTypeDisplayDriverIC   = "display_driver_ic"
	PartTypeBasebandProcessor = "baseband_processor"
	PartTypeWiFiBluetoothChip = "wifi_bluetooth_chip"

	// Physical Components
	PartTypeHousing   = "housing"
	PartTypeBackGlass = "back_glass"
	PartTypeRearCover = "rear_cover"
	PartTypeFrame     = "frame"
	PartTypeSIM       = "sim_tray"

	// Small Parts & Hardware
	PartTypeScrew                = "screw"
	PartTypeAdhesive             = "adhesive"
	PartTypeThermalPaste         = "thermal_paste"
	PartTypeThermalPad           = "thermal_pad"
	PartTypeWaterDamageIndicator = "water_damage_indicator"
	PartTypeMesh                 = "speaker_mesh"
	PartTypeGasket               = "gasket"
	PartTypeSeal                 = "seal"
	PartTypeCableClip            = "cable_clip"
	PartTypeBracket              = "bracket"
	PartTypeEMIShield            = "emi_shield"
	PartTypeFoamPad              = "foam_pad"

	// Other
	PartTypeVibrator  = "vibrator"
	PartTypeCooling   = "cooling"
	PartTypeConnector = "connector"
	PartTypeOther     = "other"
)

// Part Categories
const (
	CategoryDisplay      = "display"
	CategoryPowerSystem  = "power_system"
	CategoryAudio        = "audio"
	CategoryMainBoard    = "main_board"
	CategoryCameraSensor = "camera_system"
	CategoryConnectivity = "connectivity"
	CategoryMechanical   = "mechanical"
	CategoryInternal     = "internal"
	CategoryExternal     = "external"
	CategoryConsumable   = "consumable"
)

// Quality Grades
const (
	QualityOEM         = "oem"
	QualityOriginal    = "original"
	QualityAftermarket = "aftermarket"
	QualityRefurbished = "refurbished"
	QualityUsed        = "used"
	QualityCopy        = "copy"
	QualityAAA         = "aaa_grade"
	QualityAA          = "aa_grade"
	QualityA           = "a_grade"
)

// Condition Status
const (
	ConditionNew         = "new"
	ConditionLikeNew     = "like_new"
	ConditionRefurbished = "refurbished"
	ConditionUsed        = "used"
	ConditionWorking     = "working_pull"
	ConditionAsIs        = "as_is"
	ConditionDefective   = "defective"
	ConditionForParts    = "for_parts"
)

// Stock Status
const (
	StockAvailable  = "available"
	StockLow        = "low"
	StockCritical   = "critical"
	StockOutOfStock = "out_of_stock"
	StockOnOrder    = "on_order"
	StockReserved   = "reserved"
	StockAllocated  = "allocated"
	StockQuarantine = "quarantine"
	StockObsolete   = "obsolete"
)

// Usage Priority
const (
	PriorityEmergency = "emergency"
	PriorityCritical  = "critical"
	PriorityHigh      = "high"
	PriorityNormal    = "normal"
	PriorityLow       = "low"
)

// Warranty Types
const (
	WarrantyNone     = "none"
	Warranty30Days   = "30_days"
	Warranty60Days   = "60_days"
	Warranty90Days   = "90_days"
	Warranty180Days  = "180_days"
	Warranty1Year    = "1_year"
	WarrantyLifetime = "lifetime"
)

// Supplier Types
const (
	SupplierOEM           = "oem"
	SupplierAuthorized    = "authorized"
	SupplierDistributor   = "distributor"
	SupplierWholesaler    = "wholesaler"
	SupplierFactory       = "factory"
	SupplierLocal         = "local"
	SupplierInternational = "international"
)

// Certification Types
const (
	CertificationISO       = "iso"
	CertificationCE        = "ce"
	CertificationFCC       = "fcc"
	CertificationRoHS      = "rohs"
	CertificationUL        = "ul"
	CertificationOEM       = "oem_certified"
	CertificationQualified = "qualified"
)

// Compatibility Levels
const (
	CompatibilityExact       = "exact_match"
	CompatibilityCompatible  = "compatible"
	CompatibilityAlternative = "alternative"
	CompatibilityUniversal   = "universal"
	CompatibilityLimited     = "limited"
	CompatibilityNone        = "incompatible"
)

// Movement Types
const (
	MovementReceive    = "receive"
	MovementIssue      = "issue"
	MovementReturn     = "return"
	MovementTransfer   = "transfer"
	MovementAdjustment = "adjustment"
	MovementScrap      = "scrap"
	MovementConsume    = "consume"
	MovementReserve    = "reserve"
)

// Request Status
const (
	RequestPending   = "pending"
	RequestApproved  = "approved"
	RequestAllocated = "allocated"
	RequestIssued    = "issued"
	RequestCompleted = "completed"
	RequestCancelled = "cancelled"
	RequestRejected  = "rejected"
)

// Testing Status
const (
	TestingNotRequired = "not_required"
	TestingPending     = "pending"
	TestingInProgress  = "in_progress"
	TestingPassed      = "passed"
	TestingFailed      = "failed"
	TestingPartial     = "partial_pass"
)

// Defect Types
const (
	DefectNone          = "none"
	DefectCosmetic      = "cosmetic"
	DefectFunctional    = "functional"
	DefectCompatibility = "compatibility"
	DefectPerformance   = "performance"
	DefectManufacturing = "manufacturing"
	DefectPackaging     = "packaging"
)

// Return Reasons
const (
	ReturnDefective      = "defective"
	ReturnWrongPart      = "wrong_part"
	ReturnNotCompatible  = "not_compatible"
	ReturnNotNeeded      = "not_needed"
	ReturnCustomerReturn = "customer_return"
	ReturnWarranty       = "warranty"
	ReturnDamaged        = "damaged"
)

// Alert Types
const (
	AlertLowStock      = "low_stock"
	AlertStockout      = "stockout"
	AlertExpiring      = "expiring"
	AlertQualityIssue  = "quality_issue"
	AlertPriceIncrease = "price_increase"
	AlertSupplierIssue = "supplier_issue"
	AlertHighDemand    = "high_demand"
	AlertObsolete      = "obsolete_warning"
)

// Cost Categories
const (
	CostPurchase   = "purchase"
	CostShipping   = "shipping"
	CostHandling   = "handling"
	CostStorage    = "storage"
	CostDuty       = "import_duty"
	CostTax        = "tax"
	CostInspection = "inspection"
)

// Default Values
const (
	DefaultMinStock          = 5
	DefaultMaxStock          = 500
	DefaultReorderPoint      = 10
	DefaultReorderQuantity   = 50
	DefaultLeadTimeDays      = 7
	DefaultWarrantyDays      = 90
	DefaultShelfLifeMonths   = 24
	DefaultCriticalThreshold = 3
	DefaultLowThreshold      = 10
	DefaultMarkupPercentage  = 0.50 // 50% markup
	DefaultScrapPercentage   = 0.02 // 2% scrap rate
	DefaultDefectRate        = 0.01 // 1% defect rate
)

// Thresholds
const (
	EmergencyStockLevel   = 2
	CriticalStockLevel    = 5
	LowStockLevel         = 15
	OptimalStockLevel     = 50
	MaxStockLevel         = 200
	HighUsageThreshold    = 100 // parts per month
	SlowMovingDays        = 60
	ObsoleteDays          = 180
	HighCostThreshold     = 100.00
	QualityScoreThreshold = 85.0
	DefectRateWarning     = 0.03 // 3%
	SupplierRatingMinimum = 4.0
)

// Status Values
const (
	StatusActive       = "active"
	StatusInactive     = "inactive"
	StatusPending      = "pending"
	StatusApproved     = "approved"
	StatusSuspended    = "suspended"
	StatusDiscontinued = "discontinued"
)

// Validation Messages
const (
	ErrInvalidPartNumber   = "invalid part number format"
	ErrDuplicatePartNumber = "part number already exists"
	ErrInvalidQuantity     = "quantity must be positive"
	ErrInsufficientStock   = "insufficient stock available"
	ErrIncompatiblePart    = "part not compatible with device"
	ErrExpiredPart         = "part has exceeded shelf life"
	ErrQualityCheckFailed  = "part failed quality check"
	ErrSupplierNotApproved = "supplier not approved for this part"
	ErrExceedsMaxStock     = "quantity exceeds maximum stock level"
	ErrBelowMinStock       = "stock below minimum level"
)
