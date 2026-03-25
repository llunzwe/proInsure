package accessory

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AccessoryReview represents customer reviews
type AccessoryReview struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID    uuid.UUID  `gorm:"type:uuid;not null" json:"accessory_id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	OrderID        *uuid.UUID `gorm:"type:uuid" json:"order_id,omitempty"`
	Rating         int        `gorm:"not null" json:"rating"`
	Title          string     `json:"title"`
	Comment        string     `json:"comment"`
	Pros           string     `json:"pros,omitempty"`
	Cons           string     `json:"cons,omitempty"`
	Images         string     `json:"images,omitempty"` // JSON array
	VideoURL       string     `json:"video_url,omitempty"`
	IsVerified     bool       `gorm:"default:false" json:"is_verified"`
	HelpfulCount   int        `gorm:"default:0" json:"helpful_count"`
	UnhelpfulCount int        `gorm:"default:0" json:"unhelpful_count"`
	Status         string     `gorm:"default:'pending'" json:"status"`
	Response       string     `json:"response,omitempty"`
	ResponseDate   *time.Time `json:"response_date,omitempty"`
	FlagCount      int        `gorm:"default:0" json:"flag_count"`
	FlagReason     string     `json:"flag_reason,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessoryPromotion represents promotional campaigns
type AccessoryPromotion struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID       uuid.UUID       `gorm:"type:uuid;not null" json:"accessory_id"`
	PromotionCode     string          `gorm:"uniqueIndex;not null" json:"promotion_code"`
	PromotionName     string          `gorm:"not null" json:"promotion_name"`
	PromotionType     string          `gorm:"not null" json:"promotion_type"`
	DiscountPercent   float64         `json:"discount_percent,omitempty"`
	DiscountAmount    decimal.Decimal `sql:"type:decimal(10,2)" json:"discount_amount,omitempty"`
	StartDate         time.Time       `gorm:"not null" json:"start_date"`
	EndDate           time.Time       `gorm:"not null" json:"end_date"`
	MinQuantity       int             `json:"min_quantity,omitempty"`
	MaxQuantity       int             `json:"max_quantity,omitempty"`
	UsageLimit        int             `json:"usage_limit,omitempty"`
	UsedCount         int             `gorm:"default:0" json:"used_count"`
	TargetCustomers   string          `json:"target_customers,omitempty"`   // JSON array
	ExcludedCustomers string          `json:"excluded_customers,omitempty"` // JSON array
	RequiresCoupon    bool            `gorm:"default:false" json:"requires_coupon"`
	CouponCode        string          `json:"coupon_code,omitempty"`
	StackingAllowed   bool            `gorm:"default:false" json:"stacking_allowed"`
	IsActive          bool            `gorm:"default:true" json:"is_active"`
	CreatedBy         uuid.UUID       `gorm:"type:uuid" json:"created_by"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessoryBundle represents product bundles
type AccessoryBundle struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	BundleSKU       string          `gorm:"uniqueIndex;not null" json:"bundle_sku"`
	BundleName      string          `gorm:"not null" json:"bundle_name"`
	Description     string          `json:"description"`
	MainAccessoryID uuid.UUID       `gorm:"type:uuid;not null" json:"main_accessory_id"`
	BundledItems    string          `gorm:"not null" json:"bundled_items"` // JSON array of IDs and quantities
	TotalValue      decimal.Decimal `sql:"type:decimal(10,2)" json:"total_value"`
	BundlePrice     decimal.Decimal `sql:"type:decimal(10,2)" json:"bundle_price"`
	Savings         decimal.Decimal `sql:"type:decimal(10,2)" json:"savings"`
	SavingsPercent  float64         `json:"savings_percent"`
	ImageURL        string          `json:"image_url,omitempty"`
	Stock           int             `json:"stock"`
	MinStock        int             `json:"min_stock"`
	IsActive        bool            `gorm:"default:true" json:"is_active"`
	IsFeatured      bool            `gorm:"default:false" json:"is_featured"`
	ValidFrom       *time.Time      `json:"valid_from,omitempty"`
	ValidUntil      *time.Time      `json:"valid_until,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
}

// AccessoryMovement represents stock movements
type AccessoryMovement struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID   uuid.UUID       `gorm:"type:uuid;not null" json:"accessory_id"`
	MovementType  string          `gorm:"not null" json:"movement_type"`
	ReferenceType string          `json:"reference_type,omitempty"`
	ReferenceID   string          `json:"reference_id,omitempty"`
	Quantity      int             `gorm:"not null" json:"quantity"`
	FromLocation  string          `json:"from_location,omitempty"`
	ToLocation    string          `json:"to_location,omitempty"`
	UnitCost      decimal.Decimal `sql:"type:decimal(10,2)" json:"unit_cost"`
	TotalCost     decimal.Decimal `sql:"type:decimal(10,2)" json:"total_cost"`
	StockBefore   int             `json:"stock_before"`
	StockAfter    int             `json:"stock_after"`
	Reason        string          `json:"reason,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	PerformedBy   uuid.UUID       `gorm:"type:uuid" json:"performed_by"`
	ApprovedBy    *uuid.UUID      `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt    *time.Time      `json:"approved_at,omitempty"`
	BatchNumber   string          `json:"batch_number,omitempty"`
	ExpiryDate    *time.Time      `json:"expiry_date,omitempty"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// AccessoryClaim represents warranty or damage claims
type AccessoryClaim struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID     uuid.UUID       `gorm:"type:uuid;not null" json:"accessory_id"`
	UserID          uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	OrderID         uuid.UUID       `gorm:"type:uuid" json:"order_id"`
	ClaimNumber     string          `gorm:"uniqueIndex;not null" json:"claim_number"`
	ClaimType       string          `gorm:"not null" json:"claim_type"`
	Status          string          `gorm:"default:'pending'" json:"status"`
	Description     string          `gorm:"not null" json:"description"`
	Images          string          `json:"images,omitempty"` // JSON array
	PurchaseDate    time.Time       `json:"purchase_date"`
	ClaimDate       time.Time       `json:"claim_date"`
	Resolution      string          `json:"resolution,omitempty"`
	ResolutionDate  *time.Time      `json:"resolution_date,omitempty"`
	ReplacementSent bool            `gorm:"default:false" json:"replacement_sent"`
	RefundAmount    decimal.Decimal `sql:"type:decimal(10,2)" json:"refund_amount"`
	RefundStatus    string          `json:"refund_status,omitempty"`
	TrackingNumber  string          `json:"tracking_number,omitempty"`
	Notes           string          `json:"notes,omitempty"`
	ProcessedBy     *uuid.UUID      `gorm:"type:uuid" json:"processed_by,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessoryPriceHistory tracks price changes
type AccessoryPriceHistory struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID   uuid.UUID       `gorm:"type:uuid;not null" json:"accessory_id"`
	OldPrice      decimal.Decimal `sql:"type:decimal(10,2)" json:"old_price"`
	NewPrice      decimal.Decimal `sql:"type:decimal(10,2)" json:"new_price"`
	PriceType     string          `json:"price_type"` // retail, wholesale, sale, etc.
	ChangeReason  string          `json:"change_reason,omitempty"`
	ChangePercent float64         `json:"change_percent"`
	EffectiveDate time.Time       `json:"effective_date"`
	EndDate       *time.Time      `json:"end_date,omitempty"`
	ApprovedBy    uuid.UUID       `gorm:"type:uuid" json:"approved_by"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// AccessoryCompatibilityMatrix detailed compatibility mapping
type AccessoryCompatibilityMatrix struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID        uuid.UUID  `gorm:"type:uuid;not null" json:"accessory_id"`
	DeviceBrand        string     `json:"device_brand"`
	DeviceModel        string     `json:"device_model"`
	DeviceCategory     string     `json:"device_category,omitempty"`
	CompatibilityLevel string     `json:"compatibility_level"` // full, partial, none
	TestedDate         *time.Time `json:"tested_date,omitempty"`
	TestedBy           string     `json:"tested_by,omitempty"`
	Notes              string     `json:"notes,omitempty"`
	Limitations        string     `json:"limitations,omitempty"`
	RequiredAdapter    string     `json:"required_adapter,omitempty"`
	IsVerified         bool       `gorm:"default:false" json:"is_verified"`
	UserReported       bool       `gorm:"default:false" json:"user_reported"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessoryAlert represents inventory or price alerts
type AccessoryAlert struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID      uuid.UUID  `gorm:"type:uuid;not null" json:"accessory_id"`
	AlertType        string     `gorm:"not null" json:"alert_type"`
	AlertLevel       string     `gorm:"default:'info'" json:"alert_level"` // info, warning, critical
	Title            string     `gorm:"not null" json:"title"`
	Message          string     `json:"message"`
	TriggerValue     string     `json:"trigger_value,omitempty"`
	CurrentValue     string     `json:"current_value,omitempty"`
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	IsAcknowledged   bool       `gorm:"default:false" json:"is_acknowledged"`
	AcknowledgedBy   *uuid.UUID `gorm:"type:uuid" json:"acknowledged_by,omitempty"`
	AcknowledgedAt   *time.Time `json:"acknowledged_at,omitempty"`
	ActionRequired   string     `json:"action_required,omitempty"`
	ActionTaken      string     `json:"action_taken,omitempty"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	AutoResolve      bool       `gorm:"default:false" json:"auto_resolve"`
	NotificationSent bool       `gorm:"default:false" json:"notification_sent"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessoryRecommendation product recommendations
type AccessoryRecommendation struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID        uuid.UUID `gorm:"type:uuid;not null" json:"accessory_id"`
	RecommendedID      uuid.UUID `gorm:"type:uuid;not null" json:"recommended_id"`
	RecommendationType string    `json:"recommendation_type"` // complement, upgrade, alternative
	ConfidenceScore    float64   `json:"confidence_score"`
	Reason             string    `json:"reason,omitempty"`
	PurchasedTogether  int       `gorm:"default:0" json:"purchased_together"`
	ViewedTogether     int       `gorm:"default:0" json:"viewed_together"`
	ConversionRate     float64   `json:"conversion_rate"`
	IsActive           bool      `gorm:"default:true" json:"is_active"`
	IsFeatured         bool      `gorm:"default:false" json:"is_featured"`
	Priority           int       `gorm:"default:0" json:"priority"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessoryInsurance insurance coverage for high-value accessories
type AccessoryInsurance struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID     uuid.UUID       `gorm:"type:uuid;not null" json:"accessory_id"`
	PolicyNumber    string          `gorm:"uniqueIndex;not null" json:"policy_number"`
	InsuranceType   string          `gorm:"not null" json:"insurance_type"`
	Provider        string          `json:"provider"`
	Coverage        decimal.Decimal `sql:"type:decimal(10,2)" json:"coverage"`
	Premium         decimal.Decimal `sql:"type:decimal(10,2)" json:"premium"`
	Deductible      decimal.Decimal `sql:"type:decimal(10,2)" json:"deductible"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
	Status          string          `gorm:"default:'active'" json:"status"`
	ClaimsMade      int             `gorm:"default:0" json:"claims_made"`
	ClaimsLimit     int             `json:"claims_limit"`
	CoverageDetails string          `json:"coverage_details,omitempty"` // JSON object
	Exclusions      string          `json:"exclusions,omitempty"`       // JSON array
	IsTransferable  bool            `gorm:"default:false" json:"is_transferable"`
	RequiresProof   bool            `gorm:"default:true" json:"requires_proof"`
	DocumentURL     string          `json:"document_url,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	CancelledAt     *time.Time      `json:"cancelled_at,omitempty"`
}

// AccessoryAuthenticationLog tracks authenticity verification attempts
type AccessoryAuthenticationLog struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID      uuid.UUID  `gorm:"type:uuid;not null" json:"accessory_id"`
	SKU              string     `json:"sku"`
	SerialNumber     string     `json:"serial_number,omitempty"`
	VerificationDate time.Time  `json:"verification_date"`
	Method           string     `json:"method"` // qr_code, serial, hologram, online
	Result           string     `json:"result"` // genuine, counterfeit, suspicious
	ConfidenceScore  float64    `json:"confidence_score"`
	UserID           *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	IPAddress        string     `json:"ip_address,omitempty"`
	Location         string     `json:"location,omitempty"`
	DeviceInfo       string     `json:"device_info,omitempty"`
	AlertRaised      bool       `gorm:"default:false" json:"alert_raised"`
	Notes            string     `json:"notes,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// AccessoryCustomerFeedback tracks customer satisfaction
type AccessoryCustomerFeedback struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID       uuid.UUID  `gorm:"type:uuid;not null" json:"accessory_id"`
	CustomerID        uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	OrderID           uuid.UUID  `gorm:"type:uuid" json:"order_id"`
	SatisfactionScore int        `json:"satisfaction_score"` // 1-10
	ValueForMoney     int        `json:"value_for_money"`    // 1-5
	QualityRating     int        `json:"quality_rating"`     // 1-5
	RecommendToOthers bool       `json:"recommend_to_others"`
	RepurchaseIntent  bool       `json:"repurchase_intent"`
	FeedbackText      string     `json:"feedback_text,omitempty"`
	ImprovementAreas  string     `json:"improvement_areas,omitempty"`  // JSON array
	UsageFrequency    string     `json:"usage_frequency"`              // daily, weekly, monthly, rarely
	UsageScenarios    string     `json:"usage_scenarios,omitempty"`    // JSON array
	ComparedToBrands  string     `json:"compared_to_brands,omitempty"` // JSON array
	ResponseProvided  bool       `gorm:"default:false" json:"response_provided"`
	ResponseText      string     `json:"response_text,omitempty"`
	ResponseDate      *time.Time `json:"response_date,omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessoryMarketAnalytics tracks market performance
type AccessoryMarketAnalytics struct {
	ID                   uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID          uuid.UUID       `gorm:"type:uuid;not null" json:"accessory_id"`
	AnalysisPeriod       string          `json:"analysis_period"`
	StartDate            time.Time       `json:"start_date"`
	EndDate              time.Time       `json:"end_date"`
	TotalViews           int             `json:"total_views"`
	UniqueViewers        int             `json:"unique_viewers"`
	AddToCartCount       int             `json:"add_to_cart_count"`
	PurchaseCount        int             `json:"purchase_count"`
	ConversionRate       float64         `json:"conversion_rate"`
	AbandonmentRate      float64         `json:"abandonment_rate"`
	ReturnRate           float64         `json:"return_rate"`
	Revenue              decimal.Decimal `sql:"type:decimal(10,2)" json:"revenue"`
	Profit               decimal.Decimal `sql:"type:decimal(10,2)" json:"profit"`
	MarketRank           int             `json:"market_rank"`
	CategoryRank         int             `json:"category_rank"`
	SearchRank           int             `json:"search_rank"`
	CompetitorGrowth     float64         `json:"competitor_growth"`
	PricePosition        string          `json:"price_position"`                  // premium, competitive, budget
	CustomerDemographics string          `json:"customer_demographics,omitempty"` // JSON object
	TopSearchTerms       string          `json:"top_search_terms,omitempty"`      // JSON array
	CrossSellSuccess     float64         `json:"cross_sell_success"`
	UpSellSuccess        float64         `json:"up_sell_success"`
	CreatedAt            time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// AccessorySubscription handles subscription-based accessories
type AccessorySubscription struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	AccessoryID       uuid.UUID       `gorm:"type:uuid;not null" json:"accessory_id"`
	CustomerID        uuid.UUID       `gorm:"type:uuid;not null" json:"customer_id"`
	SubscriptionPlan  string          `json:"subscription_plan"` // monthly, quarterly, annual
	PlanPrice         decimal.Decimal `sql:"type:decimal(10,2)" json:"plan_price"`
	Discount          float64         `json:"discount"`
	StartDate         time.Time       `json:"start_date"`
	NextDelivery      time.Time       `json:"next_delivery"`
	LastDelivery      *time.Time      `json:"last_delivery,omitempty"`
	DeliveryFrequency string          `json:"delivery_frequency"` // monthly, bi-monthly, quarterly
	Quantity          int             `json:"quantity"`
	Status            string          `gorm:"default:'active'" json:"status"`
	PauseStart        *time.Time      `json:"pause_start,omitempty"`
	PauseEnd          *time.Time      `json:"pause_end,omitempty"`
	CancelReason      string          `json:"cancel_reason,omitempty"`
	TotalDeliveries   int             `json:"total_deliveries"`
	TotalSpent        decimal.Decimal `sql:"type:decimal(10,2)" json:"total_spent"`
	PaymentMethod     string          `json:"payment_method"`
	AutoRenew         bool            `gorm:"default:true" json:"auto_renew"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	CancelledAt       *time.Time      `json:"cancelled_at,omitempty"`
}
