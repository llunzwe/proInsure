package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceFinancing represents device financing/credit purchase agreements
type DeviceFinancing struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	FinanceStatus string    `gorm:"type:varchar(50);default:'pending'" json:"finance_status"` // pending, approved, active, completed, defaulted, cancelled
	FinanceType   string    `gorm:"type:varchar(50)" json:"finance_type"`                     // installment, lease, loan, buy_now_pay_later

	// Finance Provider
	FinanceProvider string     `json:"finance_provider"`
	ProviderID      *uuid.UUID `gorm:"type:uuid" json:"provider_id"`
	ApplicationID   string     `json:"application_id"`
	ApprovalDate    *time.Time `json:"approval_date"`
	ContractNumber  string     `json:"contract_number"`

	// Finance Terms
	FinanceStartDate   time.Time `json:"finance_start_date"`
	FinanceEndDate     time.Time `json:"finance_end_date"`
	TotalFinanceAmount float64   `json:"total_finance_amount"`
	DownPayment        float64   `json:"down_payment"`
	PrincipalAmount    float64   `json:"principal_amount"`
	InterestRate       float64   `json:"interest_rate"` // Annual percentage rate
	InterestAmount     float64   `json:"interest_amount"`
	ProcessingFee      float64   `json:"processing_fee"`

	// Payment Schedule
	PaymentFrequency  string     `json:"payment_frequency"` // weekly, bi_weekly, monthly
	InstallmentAmount float64    `json:"installment_amount"`
	TotalInstallments int        `json:"total_installments"`
	InstallmentsPaid  int        `json:"installments_paid"`
	PaymentDayOfMonth int        `json:"payment_day_of_month"`
	NextPaymentDate   *time.Time `json:"next_payment_date"`

	// Outstanding & Status
	OutstandingBalance float64    `json:"outstanding_balance"`
	TotalPaidAmount    float64    `json:"total_paid_amount"`
	ArrearsAmount      float64    `json:"arrears_amount"`
	DefaultStatus      bool       `gorm:"default:false" json:"default_status"`
	DefaultDate        *time.Time `json:"default_date"`
	DaysInArrears      int        `json:"days_in_arrears"`

	// Early Settlement
	EarlySettlementAllowed bool       `gorm:"default:true" json:"early_settlement_allowed"`
	EarlySettlementFee     float64    `json:"early_settlement_fee"`
	EarlySettlementAmount  float64    `json:"early_settlement_amount"`
	SettlementDate         *time.Time `json:"settlement_date"`

	// Credit Check
	CreditCheckStatus string     `json:"credit_check_status"` // passed, failed, manual_review
	CreditScore       int        `json:"credit_score"`
	CreditCheckDate   *time.Time `json:"credit_check_date"`
	RiskCategory      string     `json:"risk_category"` // low, medium, high

	// Co-signer/Guarantor
	RequiresCoSigner bool       `gorm:"default:false" json:"requires_co_signer"`
	CoSignerID       *uuid.UUID `gorm:"type:uuid" json:"co_signer_id"`
	CoSignerName     string     `json:"co_signer_name"`
	CoSignerRelation string     `json:"co_signer_relation"`

	// Insurance Requirement
	InsuranceRequired bool       `gorm:"default:true" json:"insurance_required"`
	InsurancePolicyID *uuid.UUID `gorm:"type:uuid" json:"insurance_policy_id"`
	InsuranceStatus   string     `json:"insurance_status"` // active, lapsed, cancelled

	// Documentation
	ContractURL        string `json:"contract_url"`
	PaymentScheduleDoc string `json:"payment_schedule_doc"`
	CreditAgreementDoc string `json:"credit_agreement_doc"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User, Provider, CoSigner should be loaded via service layer using UserID, ProviderID, CoSignerID to avoid circular import
	// InsurancePolicy should be loaded via service layer using InsurancePolicyID to avoid circular import
}

// FinancePayment tracks individual finance payments
type FinancePayment struct {
	database.BaseModel
	FinancingID      uuid.UUID  `gorm:"type:uuid;not null" json:"financing_id"`
	PaymentNumber    int        `json:"payment_number"` // 1st, 2nd, etc.
	DueDate          time.Time  `json:"due_date"`
	PaymentDate      *time.Time `json:"payment_date"`
	DueAmount        float64    `json:"due_amount"`
	PaidAmount       float64    `json:"paid_amount"`
	PrincipalPortion float64    `json:"principal_portion"`
	InterestPortion  float64    `json:"interest_portion"`
	LateFee          float64    `json:"late_fee"`
	PaymentStatus    string     `json:"payment_status"` // pending, paid, partial, missed, late
	PaymentMethod    string     `json:"payment_method"` // auto_debit, card, bank_transfer
	TransactionID    string     `json:"transaction_id"`
	ReceiptURL       string     `json:"receipt_url"`

	// Relationships
	Financing DeviceFinancing `gorm:"foreignKey:FinancingID" json:"financing,omitempty"`
}

// TableName returns the table name
func (t *DeviceFinancing) TableName() string {
	return "device_financings"
}

func (t *FinancePayment) TableName() string {
	return "finance_payments"
}

// BeforeCreate handles pre-creation logic
func (df *DeviceFinancing) BeforeCreate(tx *gorm.DB) error {
	if err := df.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// CalculateInterest calculates total interest amount
func (df *DeviceFinancing) CalculateInterest() {
	monthlyRate := df.InterestRate / 12 / 100
	months := df.TotalInstallments
	principal := df.PrincipalAmount

	// Simple interest calculation
	df.InterestAmount = principal * monthlyRate * float64(months)
	df.TotalFinanceAmount = principal + df.InterestAmount + df.ProcessingFee
}

// CalculateInstallment calculates monthly installment amount
func (df *DeviceFinancing) CalculateInstallment() {
	if df.TotalInstallments == 0 {
		return
	}
	df.InstallmentAmount = df.TotalFinanceAmount / float64(df.TotalInstallments)
}

// UpdatePaymentStatus updates payment tracking
func (df *DeviceFinancing) UpdatePaymentStatus(paidAmount float64) {
	df.InstallmentsPaid++
	df.TotalPaidAmount += paidAmount
	df.OutstandingBalance = df.TotalFinanceAmount - df.TotalPaidAmount

	if df.OutstandingBalance <= 0 {
		df.FinanceStatus = "completed"
		now := time.Now()
		df.SettlementDate = &now
	}
}

// CalculateEarlySettlement calculates early settlement amount
func (df *DeviceFinancing) CalculateEarlySettlement() {
	remainingPrincipal := df.OutstandingBalance -
		(df.InterestAmount * (1 - float64(df.InstallmentsPaid)/float64(df.TotalInstallments)))
	df.EarlySettlementAmount = remainingPrincipal + df.EarlySettlementFee
}

// CheckDefault checks if account is in default
func (df *DeviceFinancing) CheckDefault(daysTolerance int) {
	if df.DaysInArrears > daysTolerance {
		df.DefaultStatus = true
		if df.DefaultDate == nil {
			now := time.Now()
			df.DefaultDate = &now
		}
	}
}

// IsEligibleForFinancing checks basic eligibility
func (df *DeviceFinancing) IsEligibleForFinancing() bool {
	return df.CreditCheckStatus == "passed" &&
		df.CreditScore >= 600 &&
		!df.DefaultStatus
}
