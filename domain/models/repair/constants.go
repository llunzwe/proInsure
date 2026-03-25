package repair

// RepairShop Certification Levels
const (
	CertificationBronze   = "bronze"
	CertificationSilver   = "silver"
	CertificationGold     = "gold"
	CertificationPlatinum = "platinum"
	CertificationDiamond  = "diamond"
)

// RepairShop Types
const (
	ShopTypeAuthorized   = "authorized"
	ShopTypePartner      = "partner"
	ShopTypeIndependent  = "independent"
	ShopTypeFranchise    = "franchise"
	ShopTypeManufacturer = "manufacturer"
	ShopTypeMobile       = "mobile"
	ShopTypePopUp        = "pop_up"
	ShopTypeKiosk        = "kiosk"
)

// Repair Types
const (
	RepairTypeScreen         = "screen"
	RepairTypeBattery        = "battery"
	RepairTypeCamera         = "camera"
	RepairTypeWaterDamage    = "water_damage"
	RepairTypePort           = "port"
	RepairTypeSpeaker        = "speaker"
	RepairTypeMicrophone     = "microphone"
	RepairTypeButton         = "button"
	RepairTypeMotherboard    = "motherboard"
	RepairTypeSoftware       = "software"
	RepairTypeFrame          = "frame"
	RepairTypeBackGlass      = "back_glass"
	RepairTypeAntenna        = "antenna"
	RepairTypeSensor         = "sensor"
	RepairTypeVibrationMotor = "vibration_motor"
	RepairTypeDiagnostic     = "diagnostic"
	RepairTypePreventive     = "preventive"
	RepairTypeCosmetic       = "cosmetic"
)

// Repair Status
const (
	RepairStatusScheduled      = "scheduled"
	RepairStatusCheckedIn      = "checked_in"
	RepairStatusDiagnosing     = "diagnosing"
	RepairStatusWaitingParts   = "waiting_parts"
	RepairStatusInProgress     = "in_progress"
	RepairStatusQualityCheck   = "quality_check"
	RepairStatusCompleted      = "completed"
	RepairStatusReadyForPickup = "ready_for_pickup"
	RepairStatusDelivered      = "delivered"
	RepairStatusCancelled      = "cancelled"
	RepairStatusDelayed        = "delayed"
	RepairStatusOnHold         = "on_hold"
	RepairStatusEscalated      = "escalated"
	RepairStatusFailed         = "failed"
)

// Repair Priority
const (
	PriorityLow       = "low"
	PriorityNormal    = "normal"
	PriorityHigh      = "high"
	PriorityUrgent    = "urgent"
	PriorityCritical  = "critical"
	PriorityEmergency = "emergency"
	PriorityVIP       = "vip"
)

// Service Types
const (
	ServiceInStore      = "in_store"
	ServiceMobile       = "mobile"
	ServicePickup       = "pickup"
	ServiceDropOff      = "drop_off"
	ServiceMailIn       = "mail_in"
	ServiceRemote       = "remote"
	ServiceOnSite       = "on_site"
	ServiceExpress      = "express"
	ServiceWhileYouWait = "while_you_wait"
)

// Part Quality
const (
	PartOriginal      = "original"
	PartOEM           = "oem"
	PartAftermarket   = "aftermarket"
	PartRefurbished   = "refurbished"
	PartThirdParty    = "third_party"
	PartCustomerParts = "customer_parts"
)

// Technician Capacity
const (
	DefaultTechnicianDailyCapacity = 8 // Number of repairs a technician can handle per day
)

// Warranty Types
const (
	WarrantyNone         = "none"
	WarrantyLabor        = "labor"
	WarrantyParts        = "parts"
	WarrantyFull         = "full"
	WarrantyLimited      = "limited"
	WarrantyExtended     = "extended"
	WarrantyManufacturer = "manufacturer"
	WarrantyLifetime     = "lifetime"
)

// Device Condition (for assessment)
const (
	ConditionNew       = "new"
	ConditionLikeNew   = "like_new"
	ConditionExcellent = "excellent"
	ConditionGood      = "good"
	ConditionFair      = "fair"
	ConditionPoor      = "poor"
	ConditionBroken    = "broken"
	ConditionSalvage   = "salvage"
)

// Replacement Type
const (
	ReplacementLikeForLike = "like_for_like"
	ReplacementUpgrade     = "upgrade"
	ReplacementDowngrade   = "downgrade"
	ReplacementSimilar     = "similar"
	ReplacementTemporary   = "temporary"
	ReplacementRefurbished = "refurbished"
	ReplacementNew         = "new"
)

// Replacement Reason
const (
	ReasonTheft        = "theft"
	ReasonDamage       = "damage"
	ReasonLoss         = "loss"
	ReasonDefect       = "defect"
	ReasonBeyondRepair = "beyond_repair"
	ReasonObsolete     = "obsolete"
	ReasonTrade        = "trade"
)

// Delivery Method
const (
	DeliveryExpress   = "express"
	DeliveryStandard  = "standard"
	DeliveryOvernight = "overnight"
	DeliverySameDay   = "same_day"
	DeliveryTwoDay    = "two_day"
	DeliveryPickup    = "pickup"
	DeliveryScheduled = "scheduled"
	DeliveryPriority  = "priority"
)

// Loan Status
const (
	LoanStatusActive      = "active"
	LoanStatusReturned    = "returned"
	LoanStatusOverdue     = "overdue"
	LoanStatusLost        = "lost"
	LoanStatusDamaged     = "damaged"
	LoanStatusExtended    = "extended"
	LoanStatusTransferred = "transferred"
	LoanStatusWrittenOff  = "written_off"
)

// Quality Ratings
const (
	QualityExcellent = 5
	QualityGood      = 4
	QualityAverage   = 3
	QualityPoor      = 2
	QualityTerrible  = 1
)

// Review Status
const (
	ReviewPending  = "pending"
	ReviewVerified = "verified"
	ReviewFlagged  = "flagged"
	ReviewRemoved  = "removed"
	ReviewFeatured = "featured"
)

// Shop Status
const (
	ShopStatusActive      = "active"
	ShopStatusInactive    = "inactive"
	ShopStatusSuspended   = "suspended"
	ShopStatusProbation   = "probation"
	ShopStatusBlacklisted = "blacklisted"
	ShopStatusOnboarding  = "onboarding"
	ShopStatusUnderReview = "under_review"
	ShopStatusClosed      = "closed"
	ShopStatusHoliday     = "holiday"
	ShopStatusMaintenance = "maintenance"
)

// SLA Levels
const (
	SLAStandard = "standard"
	SLAPremium  = "premium"
	SLAPriority = "priority"
	SLAVIP      = "vip"
	SLACustom   = "custom"
)

// Payment Status
const (
	PaymentPending   = "pending"
	PaymentCompleted = "completed"
	PaymentFailed    = "failed"
	PaymentRefunded  = "refunded"
	PaymentDisputed  = "disputed"
	PaymentPartial   = "partial"
)

// Insurance Coverage
const (
	CoverageFull       = "full"
	CoveragePartial    = "partial"
	CoverageDeductible = "deductible"
	CoverageNone       = "none"
	CoverageReimbursed = "reimbursed"
	CoveragePending    = "pending"
)

// Default Values
const (
	DefaultCapacityPerDay        = 10
	DefaultWarrantyPeriodDays    = 90
	DefaultServiceRadiusKm       = 20
	DefaultMaxLoanDurationDays   = 14
	DefaultSecurityDeposit       = 100.0
	DefaultMinimumRating         = 4.0
	DefaultQualityThreshold      = 0.8
	DefaultOnTimeThreshold       = 0.9
	DefaultResponseTimeMinutes   = 30
	DefaultDiagnosticTimeMinutes = 45
	DefaultExpressHours          = 2
	DefaultStandardDays          = 3
	DefaultPickupWindowHours     = 48
)

// Price Multipliers
const (
	ExpressFeeMultiplier        = 1.5
	UrgentFeeMultiplier         = 2.0
	WeekendFeeMultiplier        = 1.25
	HolidayFeeMultiplier        = 1.5
	MobileServiceMultiplier     = 1.3
	AfterHoursMultiplier        = 1.4
	VIPServiceMultiplier        = 2.5
	BulkDiscountMultiplier      = 0.85
	LoyaltyDiscountMultiplier   = 0.9
	InsuranceDiscountMultiplier = 0.8
)

// Thresholds
const (
	HighValueThreshold         = 1000.0
	MaxRepairAttempts          = 3
	MaxLoanExtensions          = 2
	MinReviewLength            = 10
	MaxReviewLength            = 1000
	MinRatingForPromotion      = 4.5
	MaxDelayDays               = 7
	MaxOverdueDays             = 30
	MinInventoryThreshold      = 5
	CriticalInventoryThreshold = 2
)

// Time Limits (in hours)
const (
	DiagnosticTimeLimit       = 2
	RepairTimeLimit           = 24
	ExpressRepairTimeLimit    = 4
	QualityCheckTimeLimit     = 1
	CustomerResponseTimeLimit = 72
	PartsOrderTimeLimit       = 48
	EscalationTimeLimit       = 12
	RefundProcessTimeLimit    = 168 // 7 days
)
