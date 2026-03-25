package repair

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TemporaryDevice represents temporary devices for lending
type TemporaryDevice struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Model           string         `gorm:"not null" json:"model"`
	Brand           string         `gorm:"not null" json:"brand"`
	IMEI            string         `gorm:"uniqueIndex;not null" json:"imei"`
	SerialNumber    string         `gorm:"uniqueIndex" json:"serial_number"`
	Condition       string         `gorm:"not null" json:"condition"`
	IsAvailable     bool           `gorm:"default:true" json:"is_available"`
	CurrentUserID   *uuid.UUID     `gorm:"type:uuid" json:"current_user_id"`
	LoanStartDate   *time.Time     `json:"loan_start_date"`
	LoanEndDate     *time.Time     `json:"loan_end_date"`
	MaxLoanDuration int            `gorm:"default:14" json:"max_loan_duration"` // days
	DailyRate       float64        `gorm:"default:0" json:"daily_rate"`
	SecurityDeposit float64        `gorm:"default:100" json:"security_deposit"`
	Location        string         `json:"location"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Note: Relationships are defined in parent models package
	// These include: CurrentUser, LoanHistory
}

// DeviceLoan represents temporary device loans
type DeviceLoan struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	LoanNumber        string         `gorm:"uniqueIndex;not null" json:"loan_number"`
	TemporaryDeviceID uuid.UUID      `gorm:"type:uuid;not null" json:"temporary_device_id"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	ClaimID           *uuid.UUID     `gorm:"type:uuid" json:"claim_id"`
	StartDate         time.Time      `gorm:"not null" json:"start_date"`
	EndDate           time.Time      `gorm:"not null" json:"end_date"`
	ActualReturnDate  *time.Time     `json:"actual_return_date"`
	Status            string         `gorm:"not null;default:'active'" json:"status"` // active, returned, overdue, lost
	DailyRate         float64        `json:"daily_rate"`
	SecurityDeposit   float64        `json:"security_deposit"`
	TotalCost         float64        `json:"total_cost"`
	IsReturned        bool           `gorm:"default:false" json:"is_returned"`
	ReturnCondition   string         `json:"return_condition"`
	DamageAssessment  string         `json:"damage_assessment"`
	AdditionalCharges float64        `gorm:"default:0" json:"additional_charges"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Note: Relationships are defined in parent models package
	// These include: TemporaryDevice, User, Claim
}

// TableName methods
func (TemporaryDevice) TableName() string {
	return "temporary_devices"
}

func (DeviceLoan) TableName() string {
	return "device_loans"
}
