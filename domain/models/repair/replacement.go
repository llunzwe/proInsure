package repair

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReplacementDevice represents replacement devices available in inventory
type ReplacementDevice struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Model             string         `gorm:"not null" json:"model"`
	Brand             string         `gorm:"not null" json:"brand"`
	Manufacturer      string         `gorm:"not null" json:"manufacturer"`
	StorageCapacity   int            `json:"storage_capacity"`
	Color             string         `json:"color"`
	Condition         string         `gorm:"not null" json:"condition"` // new, refurbished, like_new
	Grade             string         `json:"grade"`                     // A+, A, B+, B, C
	MarketValue       float64        `gorm:"not null" json:"market_value"`
	ReplacementCost   float64        `gorm:"not null" json:"replacement_cost"`
	AvailableQuantity int            `gorm:"default:0" json:"available_quantity"`
	ReservedQuantity  int            `gorm:"default:0" json:"reserved_quantity"`
	Location          string         `json:"location"`                          // warehouse location
	WarrantyPeriod    int            `gorm:"default:90" json:"warranty_period"` // days
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Note: ReplacementOrders relationship is defined in parent models package
}

// ReplacementOrder represents device replacement orders
type ReplacementOrder struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNumber          string         `gorm:"uniqueIndex;not null" json:"order_number"`
	ClaimID              uuid.UUID      `gorm:"type:uuid;not null" json:"claim_id"`
	UserID               uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	OriginalDeviceID     uuid.UUID      `gorm:"type:uuid;not null" json:"original_device_id"`
	ReplacementDeviceID  uuid.UUID      `gorm:"type:uuid;not null" json:"replacement_device_id"`
	ReplacementType      string         `gorm:"not null" json:"replacement_type"`   // like_for_like, upgrade, downgrade
	ReplacementReason    string         `gorm:"not null" json:"replacement_reason"` // theft, damage, loss, defect
	DeliveryMethod       string         `gorm:"not null" json:"delivery_method"`    // express, standard, pickup
	DeliveryAddress      string         `json:"delivery_address"`
	DeliveryInstructions string         `json:"delivery_instructions"`
	EstimatedDelivery    time.Time      `json:"estimated_delivery"`
	ActualDelivery       *time.Time     `json:"actual_delivery"`
	TrackingNumber       string         `json:"tracking_number"`
	Status               string         `gorm:"not null;default:'pending'" json:"status"` // pending, processing, shipped, delivered, cancelled
	IsExpress            bool           `gorm:"default:false" json:"is_express"`
	ExpressFee           float64        `gorm:"default:0" json:"express_fee"`
	UpgradeFee           float64        `gorm:"default:0" json:"upgrade_fee"`
	TotalCost            float64        `json:"total_cost"`
	ReturnRequired       bool           `gorm:"default:true" json:"return_required"`
	ReturnDeadline       *time.Time     `json:"return_deadline"`
	ReturnTrackingNumber string         `json:"return_tracking_number"`
	IsReturned           bool           `gorm:"default:false" json:"is_returned"`
	ReturnedAt           *time.Time     `json:"returned_at"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Note: Relationships are defined in parent models package to avoid circular dependencies
	// These include: Claim, User, OriginalDevice, ReplacementDevice, StatusUpdates
}

// ReplacementStatusUpdate represents status updates for replacement orders
type ReplacementStatusUpdate struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ReplacementOrderID uuid.UUID  `gorm:"type:uuid;not null" json:"replacement_order_id"`
	Status             string     `gorm:"not null" json:"status"`
	Message            string     `json:"message"`
	Location           string     `json:"location"`
	EstimatedDelivery  *time.Time `json:"estimated_delivery"`
	UpdatedBy          uuid.UUID  `gorm:"type:uuid;not null" json:"updated_by"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Note: ReplacementOrder relationship is defined in parent models package
}

// TableName methods
func (ReplacementDevice) TableName() string {
	return "replacement_devices"
}

func (ReplacementOrder) TableName() string {
	return "replacement_orders"
}

func (ReplacementStatusUpdate) TableName() string {
	return "replacement_status_updates"
}
