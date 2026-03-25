package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceRental represents a device rental agreement
type DeviceRental struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RentalStatus   string    `gorm:"type:varchar(50);default:'reserved'" json:"rental_status"` // reserved, active, overdue, completed, cancelled
	RentalPlanType string    `gorm:"type:varchar(50)" json:"rental_plan_type"`                 // daily, weekly, monthly, long_term, flexible

	// Rental Period
	RentalStartDate   time.Time  `json:"rental_start_date"`
	RentalEndDate     time.Time  `json:"rental_end_date"`
	ActualReturnDate  *time.Time `json:"actual_return_date"`
	ExtensionCount    int        `json:"extension_count"`
	MaxRentalDuration int        `json:"max_rental_duration"` // in days

	// Pricing
	DailyRate         float64 `json:"daily_rate"`
	WeeklyRate        float64 `json:"weekly_rate"`
	MonthlyRate       float64 `json:"monthly_rate"`
	CurrentRate       float64 `json:"current_rate"` // Applied rate
	TotalRentalFee    float64 `json:"total_rental_fee"`
	SecurityDeposit   float64 `json:"security_deposit"`
	DepositStatus     string  `json:"deposit_status"`  // held, partial_returned, full_returned, forfeited
	LateReturnFee     float64 `json:"late_return_fee"` // Per day
	TotalLateFees     float64 `json:"total_late_fees"`
	DamageWaiverFee   float64 `json:"damage_waiver_fee"`
	InsuranceIncluded bool    `gorm:"default:false" json:"insurance_included"`

	// Rent-to-Own
	RentToOwnOption      bool       `gorm:"default:false" json:"rent_to_own_option"`
	RentToOwnThreshold   float64    `json:"rent_to_own_threshold"` // Total amount to own
	RentToOwnProgress    float64    `json:"rent_to_own_progress"`  // Amount paid towards ownership
	OwnershipTransferred bool       `gorm:"default:false" json:"ownership_transferred"`
	OwnershipDate        *time.Time `json:"ownership_date"`

	// Condition Tracking
	InitialCondition string  `json:"initial_condition"` // excellent, good, fair
	ReturnCondition  string  `json:"return_condition"`
	DamageAssessment string  `gorm:"type:json" json:"damage_assessment"` // JSON object
	DamageCharges    float64 `json:"damage_charges"`
	ConditionPhotos  string  `gorm:"type:json" json:"condition_photos"` // JSON array of URLs

	// Contract & Agreement
	ContractID         string     `json:"contract_id"`
	ContractSignedDate *time.Time `json:"contract_signed_date"`
	ContractURL        string     `json:"contract_url"`
	TermsAccepted      bool       `gorm:"default:false" json:"terms_accepted"`
	AutoRenewal        bool       `gorm:"default:false" json:"auto_renewal"`

	// Delivery & Return
	DeliveryMethod       string     `json:"delivery_method"` // pickup, delivery, mail
	DeliveryAddress      string     `gorm:"type:json" json:"delivery_address"`
	DeliveryDate         *time.Time `json:"delivery_date"`
	ReturnMethod         string     `json:"return_method"`
	ReturnTrackingNumber string     `json:"return_tracking_number"`

	// Relationships
	// Device and User should be loaded via service layer using DeviceID and UserID to avoid circular import
}

// RentalPayment tracks individual rental payments
type RentalPayment struct {
	database.BaseModel
	RentalID      uuid.UUID `gorm:"type:uuid;not null" json:"rental_id"`
	PaymentDate   time.Time `json:"payment_date"`
	PaymentAmount float64   `json:"payment_amount"`
	PaymentType   string    `json:"payment_type"`   // rental_fee, late_fee, damage_fee, deposit
	PaymentMethod string    `json:"payment_method"` // card, bank, cash
	PaymentStatus string    `json:"payment_status"` // pending, completed, failed
	TransactionID string    `json:"transaction_id"`

	// Relationships
	Rental DeviceRental `gorm:"foreignKey:RentalID" json:"rental,omitempty"`
}

// RentalHistory tracks rental history for devices
type RentalHistory struct {
	database.BaseModel
	DeviceID         uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RentalID         uuid.UUID `gorm:"type:uuid" json:"rental_id"`
	RentalPeriodDays int       `json:"rental_period_days"`
	TotalPaid        float64   `json:"total_paid"`
	ReturnCondition  string    `json:"return_condition"`
	OnTimeReturn     bool      `json:"on_time_return"`

	// Relationships
	// Device and User should be loaded via service layer using DeviceID and UserID to avoid circular import
	Rental DeviceRental `gorm:"foreignKey:RentalID" json:"rental,omitempty"`
}

// TableName returns the table name
func (t *DeviceRental) TableName() string {
	return "device_rentals"
}

func (t *RentalPayment) TableName() string {
	return "rental_payments"
}

func (t *RentalHistory) TableName() string {
	return "rental_histories"
}

// BeforeCreate handles pre-creation logic
func (dr *DeviceRental) BeforeCreate(tx *gorm.DB) error {
	if err := dr.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// CalculateTotalFee calculates total rental fee based on period
func (dr *DeviceRental) CalculateTotalFee() {
	days := dr.RentalEndDate.Sub(dr.RentalStartDate).Hours() / 24

	switch dr.RentalPlanType {
	case "daily":
		dr.TotalRentalFee = dr.DailyRate * days
	case "weekly":
		weeks := days / 7
		extraDays := int(days) % 7
		dr.TotalRentalFee = (dr.WeeklyRate * weeks) + (dr.DailyRate * float64(extraDays))
	case "monthly":
		months := days / 30
		extraDays := int(days) % 30
		dr.TotalRentalFee = (dr.MonthlyRate * months) + (dr.DailyRate * float64(extraDays))
	default:
		dr.TotalRentalFee = dr.CurrentRate * days
	}
}

// CalculateLateFees calculates late return fees
func (dr *DeviceRental) CalculateLateFees() {
	if dr.ActualReturnDate == nil || dr.ActualReturnDate.Before(dr.RentalEndDate) {
		dr.TotalLateFees = 0
		return
	}

	lateDays := dr.ActualReturnDate.Sub(dr.RentalEndDate).Hours() / 24
	dr.TotalLateFees = dr.LateReturnFee * lateDays
}

// UpdateRentToOwnProgress updates progress towards ownership
func (dr *DeviceRental) UpdateRentToOwnProgress(payment float64) {
	dr.RentToOwnProgress += payment
	if dr.RentToOwnOption && dr.RentToOwnProgress >= dr.RentToOwnThreshold {
		dr.OwnershipTransferred = true
		now := time.Now()
		dr.OwnershipDate = &now
	}
}

// IsOverdue checks if rental is overdue
func (dr *DeviceRental) IsOverdue() bool {
	return time.Now().After(dr.RentalEndDate) && dr.ActualReturnDate == nil
}

// ExtendRental extends the rental period
func (dr *DeviceRental) ExtendRental(days int) {
	dr.RentalEndDate = dr.RentalEndDate.AddDate(0, 0, days)
	dr.ExtensionCount++
	dr.CalculateTotalFee()
}
