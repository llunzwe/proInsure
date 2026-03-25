package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyPaymentSchedule represents scheduled payments for a policy
type PolicyPaymentSchedule struct {
	database.BaseModel
	PolicyID       uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	ScheduleNumber string    `gorm:"uniqueIndex;not null" json:"schedule_number"`
	PaymentNumber  int       `json:"payment_number"`
	TotalPayments  int       `json:"total_payments"`

	// Payment Details
	DueDate         time.Time `json:"due_date"`
	Amount          float64   `json:"amount"`
	PrincipalAmount float64   `json:"principal_amount"`
	FeeAmount       float64   `json:"fee_amount"`
	TaxAmount       float64   `json:"tax_amount"`
	Status          string    `gorm:"type:varchar(20);default:'pending'" json:"status"`

	// Payment Processing
	PaymentMethod    string     `json:"payment_method"`
	ProcessedDate    *time.Time `json:"processed_date"`
	TransactionID    string     `json:"transaction_id"`
	PaymentReference string     `json:"payment_reference"`

	// Late Payment
	GracePeriodEnd time.Time `json:"grace_period_end"`
	IsLate         bool      `gorm:"default:false" json:"is_late"`
	LateFee        float64   `json:"late_fee"`
	DaysOverdue    int       `json:"days_overdue"`

	// Retry Logic
	RetryCount    int        `json:"retry_count"`
	LastRetryDate *time.Time `json:"last_retry_date"`
	NextRetryDate *time.Time `json:"next_retry_date"`
	MaxRetries    int        `gorm:"default:3" json:"max_retries"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyPaymentSchedule) TableName() string {
	return "policy_payment_schedules"
}

// IsOverdue checks if payment is overdue
func (pps *PolicyPaymentSchedule) IsOverdue() bool {
	return pps.Status == "pending" && time.Now().After(pps.GracePeriodEnd)
}

// CanRetry checks if payment can be retried
func (pps *PolicyPaymentSchedule) CanRetry() bool {
	return pps.Status == "failed" && pps.RetryCount < pps.MaxRetries
}

// GetDaysUntilDue returns days until payment is due
func (pps *PolicyPaymentSchedule) GetDaysUntilDue() int {
	if time.Now().After(pps.DueDate) {
		return 0
	}
	return int(time.Until(pps.DueDate).Hours() / 24)
}

// GetDaysOverdue returns days overdue
func (pps *PolicyPaymentSchedule) GetDaysOverdue() int {
	if pps.Status != "pending" || time.Now().Before(pps.DueDate) {
		return 0
	}
	return int(time.Since(pps.DueDate).Hours() / 24)
}

// IsInGracePeriod checks if payment is in grace period
func (pps *PolicyPaymentSchedule) IsInGracePeriod() bool {
	now := time.Now()
	return pps.Status == "pending" &&
		now.After(pps.DueDate) &&
		now.Before(pps.GracePeriodEnd)
}

// CalculateLateFee calculates the late fee
func (pps *PolicyPaymentSchedule) CalculateLateFee() float64 {
	if !pps.IsOverdue() {
		return 0
	}

	daysOverdue := pps.GetDaysOverdue()

	// Progressive late fee structure
	var fee float64
	switch {
	case daysOverdue <= 7:
		fee = pps.Amount * 0.02 // 2% for first week
	case daysOverdue <= 14:
		fee = pps.Amount * 0.05 // 5% for second week
	case daysOverdue <= 30:
		fee = pps.Amount * 0.10 // 10% for first month
	default:
		fee = pps.Amount * 0.15 // 15% after a month
	}

	// Cap at maximum late fee
	maxFee := 50.0
	if fee > maxFee {
		fee = maxFee
	}

	return fee
}

// GetTotalAmountDue returns total amount due including late fees
func (pps *PolicyPaymentSchedule) GetTotalAmountDue() float64 {
	total := pps.Amount

	if pps.IsOverdue() {
		total += pps.CalculateLateFee()
	}

	return total
}

// ShouldSendReminder checks if a payment reminder should be sent
func (pps *PolicyPaymentSchedule) ShouldSendReminder() bool {
	if pps.Status != "pending" {
		return false
	}

	daysUntilDue := pps.GetDaysUntilDue()

	// Send reminders at 7, 3, and 1 day before due date
	return daysUntilDue == 7 || daysUntilDue == 3 || daysUntilDue == 1
}

// GetNextRetryTime calculates next retry time with exponential backoff
func (pps *PolicyPaymentSchedule) GetNextRetryTime() time.Time {
	// Exponential backoff: 1h, 4h, 12h, 24h
	hoursToWait := 1
	for i := 0; i < pps.RetryCount; i++ {
		hoursToWait *= 4
		if hoursToWait > 24 {
			hoursToWait = 24
		}
	}

	return time.Now().Add(time.Duration(hoursToWait) * time.Hour)
}

// IsLastPayment checks if this is the last payment
func (pps *PolicyPaymentSchedule) IsLastPayment() bool {
	return pps.PaymentNumber == pps.TotalPayments
}
