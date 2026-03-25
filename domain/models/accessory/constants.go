package accessory

// Accessory Types
const (
	// Cases & Protection
	AccessoryTypeCase                = "case"
	AccessoryTypeScreenProtector     = "screen_protector"
	AccessoryTypeCover               = "cover"
	AccessoryTypeWalletCase          = "wallet_case"
	AccessoryTypeFolioCase           = "folio_case"
	AccessoryTypeBumperCase          = "bumper_case"
	AccessoryTypeArmbandCase         = "armband_case"
	AccessoryTypeBatteryCase         = "battery_case"
	AccessoryTypeWaterproofPouch     = "waterproof_pouch"
	AccessoryTypePrivacyScreen       = "privacy_screen"
	AccessoryTypeCameraLensProtector = "camera_lens_protector"
	AccessoryTypePortDustPlug        = "port_dust_plug"
	AccessoryTypeSkinWrap            = "skin_wrap"

	// Power & Charging
	AccessoryTypeCharger                  = "charger"
	AccessoryTypeWirelessCharger          = "wireless_charger"
	AccessoryTypeCarCharger               = "car_charger"
	AccessoryTypeMagneticCharger          = "magnetic_charger"
	AccessoryTypeSolarCharger             = "solar_charger"
	AccessoryTypeChargingStation          = "charging_station"
	AccessoryTypeWirelessChargingReceiver = "wireless_charging_receiver"
	AccessoryTypePowerBank                = "power_bank"
	AccessoryTypeCable                    = "cable"

	// Audio
	AccessoryTypeHeadphones          = "headphones"
	AccessoryTypeTrueWirelessEarbuds = "true_wireless_earbuds"
	AccessoryTypeBluetoothHeadset    = "bluetooth_headset"
	AccessoryTypeExternalMicrophone  = "external_microphone"
	AccessoryTypeAudioAdapter        = "audio_adapter"
	AccessoryTypeFMTransmitter       = "fm_transmitter"
	AccessoryTypeAudioSplitter       = "audio_splitter"
	AccessoryTypeSpeaker             = "speaker"

	// Photography & Video
	AccessoryTypeLens             = "lens"
	AccessoryTypeSelfiStick       = "selfie_stick"
	AccessoryTypeTripod           = "tripod"
	AccessoryTypeRingLight        = "ring_light"
	AccessoryTypeGimbalStabilizer = "gimbal_stabilizer"
	AccessoryTypeExternalFlash    = "external_flash"
	AccessoryTypeLightReflector   = "light_reflector"

	// Connectivity & Adapters
	AccessoryTypeAdapter         = "adapter"
	AccessoryTypeOTGAdapter      = "otg_adapter"
	AccessoryTypeCardReader      = "card_reader"
	AccessoryTypeHDMIAdapter     = "hdmi_adapter"
	AccessoryTypeEthernetAdapter = "ethernet_adapter"
	AccessoryTypeUSBHub          = "usb_hub"
	AccessoryTypeDisplayDongle   = "display_dongle"
	AccessoryTypeDock            = "dock"

	// Gaming & Entertainment
	AccessoryTypeGameController = "game_controller"
	AccessoryTypeCoolingFan     = "cooling_fan"
	AccessoryTypeTriggerButton  = "trigger_button"
	AccessoryTypeVRHeadset      = "vr_headset"

	// Mounts & Holders
	AccessoryTypeMount          = "mount"
	AccessoryTypeCarMount       = "car_mount"
	AccessoryTypeBikeMount      = "bike_mount"
	AccessoryTypeDashboardMount = "dashboard_mount"
	AccessoryTypeStand          = "stand"
	AccessoryTypeGrip           = "grip"
	AccessoryTypePopSocket      = "pop_socket"
	AccessoryTypeRingHolder     = "ring_holder"
	AccessoryTypeMagneticHolder = "magnetic_holder"

	// Input Devices
	AccessoryTypeKeyboard = "keyboard"
	AccessoryTypeMouse    = "mouse"
	AccessoryTypeStylus   = "stylus"

	// Storage
	AccessoryTypeMemoryCard = "memory_card"
	AccessoryTypeUSBDrive   = "usb_drive"

	// Lifestyle & Tools
	AccessoryTypeSmartwatchBand = "smartwatch_band"
	AccessoryTypeSIMCardTool    = "sim_card_tool"
	AccessoryTypeCleaningKit    = "screen_cleaning_kit"

	AccessoryTypeOther = "other"
)

// Accessory Categories
const (
	CategoryProtection   = "protection"
	CategoryPower        = "power"
	CategoryAudio        = "audio"
	CategoryStorage      = "storage"
	CategoryConnectivity = "connectivity"
	CategoryInput        = "input"
	CategoryMounting     = "mounting"
	CategoryPhotography  = "photography"
	CategoryGaming       = "gaming"
	CategoryFitness      = "fitness"
)

// Quality Grades
const (
	QualityOriginal    = "original"
	QualityOEM         = "oem"
	QualityPremium     = "premium"
	QualityStandard    = "standard"
	QualityBudget      = "budget"
	QualityRefurbished = "refurbished"
)

// Condition Status
const (
	ConditionNew       = "new"
	ConditionLikeNew   = "like_new"
	ConditionExcellent = "excellent"
	ConditionGood      = "good"
	ConditionFair      = "fair"
	ConditionPoor      = "poor"
	ConditionDamaged   = "damaged"
	ConditionForParts  = "for_parts"
)

// Stock Status
const (
	StockInStock      = "in_stock"
	StockLowStock     = "low_stock"
	StockOutOfStock   = "out_of_stock"
	StockBackorder    = "backorder"
	StockDiscontinued = "discontinued"
	StockReserved     = "reserved"
	StockInTransit    = "in_transit"
)

// Warranty Status
const (
	WarrantyActive        = "active"
	WarrantyExpired       = "expired"
	WarrantyClaimed       = "claimed"
	WarrantyVoid          = "void"
	WarrantyNotApplicable = "not_applicable"
)

// Certification Types
const (
	CertificationCE       = "ce"
	CertificationFCC      = "fcc"
	CertificationRoHS     = "rohs"
	CertificationUL       = "ul"
	CertificationISO      = "iso"
	CertificationMFI      = "mfi"       // Made for iPhone/iPad
	CertificationQI       = "qi"        // Wireless charging
	CertificationIPRating = "ip_rating" // Water/dust resistance
	CertificationMILSTD   = "mil_std"   // Military standard
)

// Pricing Types
const (
	PricingFixed       = "fixed"
	PricingTiered      = "tiered"
	PricingDynamic     = "dynamic"
	PricingBundled     = "bundled"
	PricingPromotional = "promotional"
	PricingClearance   = "clearance"
)

// Supplier Types
const (
	SupplierManufacturer = "manufacturer"
	SupplierDistributor  = "distributor"
	SupplierWholesaler   = "wholesaler"
	SupplierRetailer     = "retailer"
	SupplierDropshipper  = "dropshipper"
	SupplierMarketplace  = "marketplace"
)

// Return Status
const (
	ReturnNone       = "none"
	ReturnRequested  = "requested"
	ReturnApproved   = "approved"
	ReturnShipped    = "shipped"
	ReturnReceived   = "received"
	ReturnProcessing = "processing"
	ReturnCompleted  = "completed"
	ReturnRejected   = "rejected"
)

// Promotion Types
const (
	PromotionDiscount  = "discount"
	PromotionBOGO      = "bogo" // Buy one get one
	PromotionBundle    = "bundle"
	PromotionFreeShip  = "free_shipping"
	PromotionCashback  = "cashback"
	PromotionLoyalty   = "loyalty"
	PromotionSeasonal  = "seasonal"
	PromotionClearance = "clearance_sale"
)

// Compatibility Types
const (
	CompatibilityUniversal = "universal"
	CompatibilityBrand     = "brand_specific"
	CompatibilityModel     = "model_specific"
	CompatibilityOS        = "os_specific"
	CompatibilityGeneric   = "generic"
)

// Package Types
const (
	PackageRetail = "retail"
	PackageBulk   = "bulk"
	PackageOEM    = "oem"
	PackageBundle = "bundle"
	PackageGift   = "gift"
	PackageEco    = "eco_friendly"
)

// Review Status
const (
	ReviewPending  = "pending"
	ReviewApproved = "approved"
	ReviewRejected = "rejected"
	ReviewFlagged  = "flagged"
	ReviewVerified = "verified_purchase"
)

// Alert Types
const (
	AlertLowStock    = "low_stock"
	AlertExpiring    = "expiring_soon"
	AlertPriceChange = "price_change"
	AlertNewArrival  = "new_arrival"
	AlertBackInStock = "back_in_stock"
	AlertRecall      = "recall"
	AlertDefect      = "defect_report"
)

// Movement Types
const (
	MovementInbound    = "inbound"
	MovementOutbound   = "outbound"
	MovementTransfer   = "transfer"
	MovementAdjustment = "adjustment"
	MovementReturn     = "return"
	MovementDamage     = "damage"
	MovementLoss       = "loss"
	MovementSale       = "sale"
)

// Claim Types
const (
	ClaimWarranty       = "warranty"
	ClaimDamage         = "damage"
	ClaimDefective      = "defective"
	ClaimMissing        = "missing"
	ClaimWrongItem      = "wrong_item"
	ClaimNotAsDescribed = "not_as_described"
)

// Insurance Coverage
const (
	InsuranceNone       = "none"
	InsuranceBasic      = "basic"
	InsuranceExtended   = "extended"
	InsurancePremium    = "premium"
	InsuranceAccidental = "accidental_damage"
	InsuranceTheft      = "theft_protection"
)

// Default Values
const (
	DefaultWarrantyMonths   = 12
	DefaultReturnDays       = 30
	DefaultMinStock         = 10
	DefaultMaxStock         = 1000
	DefaultReorderPoint     = 20
	DefaultLeadTimeDays     = 7
	DefaultShelfLifeMonths  = 36
	DefaultDepreciationRate = 0.15 // 15% per year
	DefaultMarkupPercentage = 0.30 // 30% markup
	DefaultDiscountMax      = 0.50 // 50% max discount
	DefaultBulkMinQuantity  = 10
	DefaultBulkDiscountRate = 0.10 // 10% bulk discount
)

// Thresholds
const (
	LowStockThreshold      = 20
	CriticalStockThreshold = 5
	SlowMovingDays         = 90
	DeadStockDays          = 180
	ExpiryWarningDays      = 30
	HighValueThreshold     = 500.00
	PopularityThreshold    = 4.0 // Rating
	MinReviewCount         = 10
	ReturnRateWarning      = 0.10 // 10%
	DefectRateWarning      = 0.05 // 5%
)

// Status Values
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusPending   = "pending"
	StatusSuspended = "suspended"
	StatusArchived  = "archived"
	StatusDeleted   = "deleted"
)

// Validation Messages
const (
	ErrInvalidSKU            = "invalid SKU format"
	ErrDuplicateSKU          = "SKU already exists"
	ErrInvalidPrice          = "price must be positive"
	ErrInvalidStock          = "stock quantity cannot be negative"
	ErrIncompatibleAccessory = "accessory not compatible with device"
	ErrExceededMaxDiscount   = "discount exceeds maximum allowed"
	ErrInvalidBarcode        = "invalid barcode format"
	ErrExpiredProduct        = "product has expired"
	ErrInsufficientStock     = "insufficient stock available"
	ErrInvalidSupplier       = "invalid supplier information"
)
