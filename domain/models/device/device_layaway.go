package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceLayaway represents a device layaway/lay-buy agreement
type DeviceLayaway struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	LayawayStatus string    `gorm:"type:varchar(50);default:'active'" json:"layaway_status"` // active, completed, cancelled, defaulted, extended

	// Layaway Terms
	LayawayStartDate time.Time  `json:"layaway_start_date"`
	LayawayEndDate   time.Time  `json:"layaway_end_date"`
	ExtendedEndDate  *time.Time `json:"extended_end_date"`
	CompletionDate   *time.Time `json:"completion_date"`

	// Financial Details
	TotalAmount       float64 `json:"total_amount"`
	InitialDeposit    float64 `json:"initial_deposit"`
	DepositPercentage float64 `json:"deposit_percentage"` // % of total required as deposit
	BalanceRemaining  float64 `json:"balance_remaining"`
	TotalPaid         float64 `json:"total_paid"`

	// Payment Schedule
	PaymentFrequency string  `json:"payment_frequency"` // weekly, bi_weekly, monthly
	PaymentAmount    float64 `json:"payment_amount"`    // Regular payment amount
	TotalPayments    int     `json:"total_payments"`    // Number of payments required
	PaymentsMade     int     `json:"payments_made"`
	MissedPayments   int     `json:"missed_payments"`

	// Storage & Handling
	StorageLocation string  `json:"storage_location"`
	StorageFee      float64 `json:"storage_fee"`
	HandlingFee     float64 `json:"handling_fee"`
	ReservationFee  float64 `json:"reservation_fee"`

	// Cancellation & Default
	CancellationFee    float64    `json:"cancellation_fee"`
	CancellationDate   *time.Time `json:"cancellation_date"`
	CancellationReason string     `json:"cancellation_reason"`
	RefundAmount       float64    `json:"refund_amount"`
	DefaultThreshold   int        `json:"default_threshold"` // Missed payments before default
	DefaultDate        *time.Time `json:"default_date"`

	// Extensions
	ExtensionsAllowed int     `json:"extensions_allowed"`
	ExtensionsUsed    int     `json:"extensions_used"`
	ExtensionFee      float64 `json:"extension_fee"`
	ExtensionDuration int     `json:"extension_duration"` // Days per extension

	// Final Payment & Collection
	FinalPaymentDue    *time.Time `json:"final_payment_due"`
	ReadyForCollection bool       `gorm:"default:false" json:"ready_for_collection"`
	CollectionMethod   string     `json:"collection_method"` // pickup, delivery
	CollectionDate     *time.Time `json:"collection_date"`
	CollectionAddress  string     `gorm:"type:json" json:"collection_address"`

	// Device Details
	ReservedDevice  bool   `gorm:"default:true" json:"reserved_device"`  // Specific device reserved
	DeviceCondition string `json:"device_condition"`                     // Condition when reserved
	PriceProtection bool   `gorm:"default:true" json:"price_protection"` // Price locked at agreement

	// Documentation
	AgreementNumber    string `json:"agreement_number"`
	AgreementURL       string `json:"agreement_url"`
	PaymentScheduleDoc string `json:"payment_schedule_doc"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// LayawayPayment tracks individual layaway payments
type LayawayPayment struct {
	database.BaseModel
	LayawayID     uuid.UUID  `gorm:"type:uuid;not null" json:"layaway_id"`
	PaymentNumber int        `json:"payment_number"`
	DueDate       time.Time  `json:"due_date"`
	PaymentDate   *time.Time `json:"payment_date"`
	DueAmount     float64    `json:"due_amount"`
	PaidAmount    float64    `json:"paid_amount"`
	PaymentStatus string     `json:"payment_status"` // pending, paid, partial, missed, late
	PaymentMethod string     `json:"payment_method"` // cash, card, bank_transfer
	LateFee       float64    `json:"late_fee"`
	TransactionID string     `json:"transaction_id"`
	ReceiptNumber string     `json:"receipt_number"`

	// Relationships
	Layaway DeviceLayaway `gorm:"foreignKey:LayawayID" json:"layaway,omitempty"`
}

// LayawaySchedule represents the payment schedule
type LayawaySchedule struct {
	database.BaseModel
	LayawayID     uuid.UUID  `gorm:"type:uuid;not null" json:"layaway_id"`
	PaymentNumber int        `json:"payment_number"`
	DueDate       time.Time  `json:"due_date"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"` // pending, paid, overdue
	PaidDate      *time.Time `json:"paid_date"`

	// Relationships
	Layaway DeviceLayaway `gorm:"foreignKey:LayawayID" json:"layaway,omitempty"`
}

// TableName returns the table name
func (t *DeviceLayaway) TableName() string {
	return "device_layaways"
}

func (t *LayawayPayment) TableName() string {
	return "layaway_payments"
}

func (t *LayawaySchedule) TableName() string {
	return "layaway_schedules"
}

// BeforeCreate handles pre-creation logic
func (dl *DeviceLayaway) BeforeCreate(tx *gorm.DB) error {
	if err := dl.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// CalculatePaymentAmount calculates regular payment amount
func (dl *DeviceLayaway) CalculatePaymentAmount() {
	remainingAmount := dl.TotalAmount - dl.InitialDeposit
	if dl.TotalPayments > 0 {
		dl.PaymentAmount = remainingAmount / float64(dl.TotalPayments)
	}
	dl.BalanceRemaining = remainingAmount
}

// ProcessPayment processes a layaway payment
func (dl *DeviceLayaway) ProcessPayment(amount float64) {
	dl.TotalPaid += amount
	dl.BalanceRemaining = dl.TotalAmount - dl.TotalPaid
	dl.PaymentsMade++

	if dl.BalanceRemaining <= 0 {
		dl.LayawayStatus = "completed"
		now := time.Now()
		dl.CompletionDate = &now
		dl.ReadyForCollection = true
	}
}

// CheckDefault checks if layaway should be marked as defaulted
func (dl *DeviceLayaway) CheckDefault() {
	if dl.MissedPayments >= dl.DefaultThreshold {
		dl.LayawayStatus = "defaulted"
		if dl.DefaultDate == nil {
			now := time.Now()
			dl.DefaultDate = &now
		}
	}
}

// ExtendLayaway extends the layaway period
func (dl *DeviceLayaway) ExtendLayaway(days int) error {
	if dl.ExtensionsUsed >= dl.ExtensionsAllowed {
		return gorm.ErrInvalidData
	}

	newEndDate := dl.LayawayEndDate.AddDate(0, 0, days)
	dl.ExtendedEndDate = &newEndDate
	dl.ExtensionsUsed++
	dl.LayawayStatus = "extended"

	// Recalculate payment schedule
	dl.CalculatePaymentAmount()

	return nil
}

// CancelLayaway cancels the layaway agreement
func (dl *DeviceLayaway) CancelLayaway(reason string) {
	dl.LayawayStatus = "cancelled"
	now := time.Now()
	dl.CancellationDate = &now
	dl.CancellationReason = reason

	// Calculate refund (total paid minus fees)
	dl.RefundAmount = dl.TotalPaid - dl.CancellationFee - dl.StorageFee - dl.HandlingFee
	if dl.RefundAmount < 0 {
		dl.RefundAmount = 0
	}
}

// IsPaymentOverdue checks if current payment is overdue
func (dl *DeviceLayaway) IsPaymentOverdue() bool {
	// Logic to check against payment schedule
	return dl.LayawayStatus == "active" && dl.MissedPayments > 0
}

// CanCollect checks if device can be collected
func (dl *DeviceLayaway) CanCollect() bool {
	return dl.LayawayStatus == "completed" && dl.BalanceRemaining <= 0
}
