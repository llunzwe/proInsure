package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceTradeIn represents a device trade-in transaction
type DeviceTradeIn struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Status   string    `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, quoted, locked, received, inspected, completed, rejected

	// Valuation
	QuotedValue   float64    `json:"quoted_value"`
	LockedValue   float64    `json:"locked_value"`   // Guaranteed value
	FinalValue    float64    `json:"final_value"`    // After inspection
	OriginalValue float64    `json:"original_value"` // Before deductions
	QuoteExpiry   *time.Time `json:"quote_expiry"`

	// Condition Assessment
	TradeInGrade      string  `gorm:"type:varchar(10)" json:"trade_in_grade"` // A+, A, B, C, D
	ConditionScore    float64 `json:"condition_score"`                        // 0-100
	EligibilityStatus string  `json:"eligibility_status"`                     // eligible, not_eligible, conditional

	// Deductions
	ScreenDeduction     float64 `json:"screen_deduction"`
	BodyDeduction       float64 `json:"body_deduction"`
	FunctionalDeduction float64 `json:"functional_deduction"`
	AccessoryDeduction  float64 `json:"accessory_deduction"`
	TotalDeductions     float64 `json:"total_deductions"`
	DeductionDetails    string  `gorm:"type:json" json:"deduction_details"` // JSON object with detailed deductions

	// Processing
	ReceivedDate    *time.Time `json:"received_date"`
	InspectionDate  *time.Time `json:"inspection_date"`
	CompletionDate  *time.Time `json:"completion_date"`
	InspectorID     *uuid.UUID `gorm:"type:uuid" json:"inspector_id"`
	InspectionNotes string     `json:"inspection_notes"`

	// Payment
	PaymentMethod string     `json:"payment_method"` // cash, credit, voucher, account_credit
	VoucherCode   string     `json:"voucher_code"`
	VoucherExpiry *time.Time `json:"voucher_expiry"`
	PaymentStatus string     `json:"payment_status"` // pending, processed, failed
	PaymentDate   *time.Time `json:"payment_date"`

	// Trade-In Options
	InstantTradeIn bool   `gorm:"default:false" json:"instant_trade_in"`
	PartnerNetwork string `json:"partner_network"` // Partner accepting trade-in
	ShippingMethod string `json:"shipping_method"` // mail_in, drop_off, pickup
	ShippingLabel  string `json:"shipping_label"`
	TrackingNumber string `json:"tracking_number"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User and Inspector should be loaded via service layer using UserID and InspectorID to avoid circular import
}

// TradeInHistory represents historical trade-in values for tracking
type TradeInHistory struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TradeInID       uuid.UUID `gorm:"type:uuid" json:"trade_in_id"`
	TradeInDate     time.Time `json:"trade_in_date"`
	TradeInValue    float64   `json:"trade_in_value"`
	DeviceCondition string    `json:"device_condition"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
	TradeIn DeviceTradeIn `gorm:"foreignKey:TradeInID" json:"trade_in,omitempty"`
}

// TableName returns the table name
func (t *DeviceTradeIn) TableName() string {
	return "device_trade_ins"
}

func (t *TradeInHistory) TableName() string {
	return "trade_in_histories"
}

// BeforeCreate handles pre-creation logic
func (dt *DeviceTradeIn) BeforeCreate(tx *gorm.DB) error {
	if err := dt.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// CalculateFinalValue calculates value after all deductions
func (dt *DeviceTradeIn) CalculateFinalValue() {
	dt.TotalDeductions = dt.ScreenDeduction + dt.BodyDeduction +
		dt.FunctionalDeduction + dt.AccessoryDeduction
	dt.FinalValue = dt.OriginalValue - dt.TotalDeductions
	if dt.FinalValue < 0 {
		dt.FinalValue = 0
	}
}

// LockValue locks the trade-in value for guarantee period
func (dt *DeviceTradeIn) LockValue(days int) {
	dt.LockedValue = dt.QuotedValue
	expiry := time.Now().AddDate(0, 0, days)
	dt.QuoteExpiry = &expiry
	dt.Status = "locked"
}

// IsQuoteValid checks if quote is still valid
func (dt *DeviceTradeIn) IsQuoteValid() bool {
	if dt.QuoteExpiry == nil {
		return false
	}
	return time.Now().Before(*dt.QuoteExpiry)
}

// Complete marks trade-in as completed
func (dt *DeviceTradeIn) Complete() {
	now := time.Now()
	dt.CompletionDate = &now
	dt.Status = "completed"
}
