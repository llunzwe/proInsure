package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceSwap represents a device swap/upgrade transaction
type DeviceSwap struct {
	database.BaseModel
	DeviceID       uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	SwapType       string     `gorm:"type:varchar(50);not null" json:"swap_type"`            // upgrade, replacement, exchange
	SwapStatus     string     `gorm:"type:varchar(50);default:'pending'" json:"swap_status"` // pending, approved, in_progress, completed, cancelled
	SwapTier       string     `gorm:"type:varchar(20)" json:"swap_tier"`                     // premium, standard, basic
	RequestDate    time.Time  `gorm:"autoCreateTime" json:"request_date"`
	ApprovalDate   *time.Time `json:"approval_date"`
	CompletionDate *time.Time `json:"completion_date"`

	// Eligibility & Pricing
	EligibilityStatus string     `gorm:"type:varchar(50)" json:"eligibility_status"` // eligible, not_eligible, pending_review
	NextUpgradeDate   *time.Time `json:"next_upgrade_date"`
	SwapFee           float64    `json:"swap_fee"`
	EarlySwapFee      float64    `json:"early_swap_fee"`
	SwapCredit        float64    `json:"swap_credit"`
	TotalCost         float64    `json:"total_cost"`

	// Device Details
	CurrentDeviceID   uuid.UUID  `gorm:"type:uuid;not null" json:"current_device_id"`
	TargetDeviceModel string     `json:"target_device_model"`
	TargetDeviceID    *uuid.UUID `gorm:"type:uuid" json:"target_device_id"`
	ReturnCondition   string     `json:"return_condition"` // excellent, good, fair, poor
	ConditionNotes    string     `json:"condition_notes"`

	// Program Details
	LoyaltyPointsUsed   int    `json:"loyalty_points_used"`
	LoyaltyPointsEarned int    `json:"loyalty_points_earned"`
	UpgradePathOptions  string `gorm:"type:json" json:"upgrade_path_options"` // JSON array of available models
	SwapReason          string `json:"swap_reason"`

	// Tracking
	TrackingNumber  string `json:"tracking_number"`
	ShippingAddress string `gorm:"type:json" json:"shipping_address"`
	ReturnLabel     string `json:"return_label"`

	// Relationships
	// Device, User, CurrentDevice, TargetDevice should be loaded via service layer using their respective IDs to avoid circular import
}

// SwapHistory represents historical swap transactions for a device
type SwapHistory struct {
	database.BaseModel
	DeviceID         uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	PreviousDeviceID uuid.UUID `gorm:"type:uuid" json:"previous_device_id"`
	SwapDate         time.Time `json:"swap_date"`
	SwapReason       string    `json:"swap_reason"`
	SwapCost         float64   `json:"swap_cost"`
	SwapType         string    `json:"swap_type"`

	// Relationships
	// Device and PreviousDevice should be loaded via service layer using DeviceID and PreviousDeviceID to avoid circular import
}

// TableName returns the table name
func (t *DeviceSwap) TableName() string {
	return "device_swaps"
}

func (t *SwapHistory) TableName() string {
	return "swap_histories"
}

// BeforeCreate handles pre-creation logic
func (ds *DeviceSwap) BeforeCreate(tx *gorm.DB) error {
	if err := ds.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// CalculateTotalCost calculates the total swap cost
func (ds *DeviceSwap) CalculateTotalCost() {
	ds.TotalCost = ds.SwapFee + ds.EarlySwapFee - ds.SwapCredit
	if ds.TotalCost < 0 {
		ds.TotalCost = 0
	}
}

// IsEligible checks if swap is eligible
func (ds *DeviceSwap) IsEligible() bool {
	return ds.EligibilityStatus == "eligible" &&
		(ds.NextUpgradeDate == nil || time.Now().After(*ds.NextUpgradeDate))
}

// Complete marks the swap as completed
func (ds *DeviceSwap) Complete() {
	now := time.Now()
	ds.CompletionDate = &now
	ds.SwapStatus = "completed"
}
