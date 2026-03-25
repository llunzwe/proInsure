package policy

import (
	"database/sql/driver"
	"errors"
	"time"
)

// ============================================
// POLICY TYPE DEFINITIONS
// ============================================

// PolicyStatus represents the lifecycle status of a policy
type PolicyStatus string

const (
	PolicyStatusDraft      PolicyStatus = "draft"
	PolicyStatusPending    PolicyStatus = "pending"
	PolicyStatusActive     PolicyStatus = "active"
	PolicyStatusSuspended  PolicyStatus = "suspended"
	PolicyStatusCancelled  PolicyStatus = "cancelled"
	PolicyStatusExpired    PolicyStatus = "expired"
	PolicyStatusRenewing   PolicyStatus = "renewing"
	PolicyStatusTerminated PolicyStatus = "terminated"
)

// Scan implements the Scanner interface for database operations
func (ps *PolicyStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		*ps = PolicyStatus(v)
	case []byte:
		*ps = PolicyStatus(v)
	default:
		return errors.New("cannot scan PolicyStatus")
	}
	return nil
}

// Value implements the driver Valuer interface
func (ps PolicyStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// IsValid checks if the policy status is valid
func (ps PolicyStatus) IsValid() bool {
	validStatuses := []PolicyStatus{
		PolicyStatusDraft, PolicyStatusPending, PolicyStatusActive,
		PolicyStatusSuspended, PolicyStatusCancelled, PolicyStatusExpired,
		PolicyStatusRenewing, PolicyStatusTerminated,
	}
	for _, status := range validStatuses {
		if ps == status {
			return true
		}
	}
	return false
}

// PolicyType represents the type of insurance policy
type PolicyType string

const (
	PolicyTypeBasic         PolicyType = "basic"
	PolicyTypeStandard      PolicyType = "standard"
	PolicyTypeComprehensive PolicyType = "comprehensive"
	PolicyTypePremium       PolicyType = "premium"
	PolicyTypePlatinum      PolicyType = "platinum"
	PolicyTypeEnterprise    PolicyType = "enterprise"
)

// PaymentFrequency represents how often payments are made
type PaymentFrequency string

const (
	PaymentFrequencyMonthly    PaymentFrequency = "monthly"
	PaymentFrequencyQuarterly  PaymentFrequency = "quarterly"
	PaymentFrequencySemiAnnual PaymentFrequency = "semi-annual"
	PaymentFrequencyAnnual     PaymentFrequency = "annual"
	PaymentFrequencyOneTime    PaymentFrequency = "one-time"
	PaymentFrequencyCustom     PaymentFrequency = "custom"
)

// PaymentStatus represents the status of payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusPaid       PaymentStatus = "paid"
	PaymentStatusOverdue    PaymentStatus = "overdue"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusDisputed   PaymentStatus = "disputed"
)

// UnderwritingStatus represents the underwriting status
type UnderwritingStatus string

const (
	UnderwritingStatusPending     UnderwritingStatus = "pending"
	UnderwritingStatusInReview    UnderwritingStatus = "in_review"
	UnderwritingStatusApproved    UnderwritingStatus = "approved"
	UnderwritingStatusConditional UnderwritingStatus = "conditional"
	UnderwritingStatusReferred    UnderwritingStatus = "referred"
	UnderwritingStatusRejected    UnderwritingStatus = "rejected"
	UnderwritingStatusExpired     UnderwritingStatus = "expired"
)

// RiskCategory represents the risk categorization
type RiskCategory string

const (
	RiskCategoryVeryLow  RiskCategory = "very_low"
	RiskCategoryLow      RiskCategory = "low"
	RiskCategoryMedium   RiskCategory = "medium"
	RiskCategoryHigh     RiskCategory = "high"
	RiskCategoryVeryHigh RiskCategory = "very_high"
	RiskCategoryCritical RiskCategory = "critical"
)

// CurrencyCode represents ISO 4217 currency codes
type CurrencyCode string

const (
	CurrencyUSD CurrencyCode = "USD"
	CurrencyEUR CurrencyCode = "EUR"
	CurrencyGBP CurrencyCode = "GBP"
	CurrencyJPY CurrencyCode = "JPY"
	CurrencyCAD CurrencyCode = "CAD"
	CurrencyAUD CurrencyCode = "AUD"
	CurrencyCHF CurrencyCode = "CHF"
)

// ChannelType represents the sales channel
type ChannelType string

const (
	ChannelDirect  ChannelType = "direct"
	ChannelAgent   ChannelType = "agent"
	ChannelBroker  ChannelType = "broker"
	ChannelPartner ChannelType = "partner"
	ChannelOnline  ChannelType = "online"
	ChannelMobile  ChannelType = "mobile"
	ChannelAPI     ChannelType = "api"
)

// DeductibleType represents how deductibles are calculated
type DeductibleType string

const (
	DeductibleTypeFixed      DeductibleType = "fixed"
	DeductibleTypePercentage DeductibleType = "percentage"
	DeductibleTypeVariable   DeductibleType = "variable"
	DeductibleTypeNone       DeductibleType = "none"
)

// ============================================
// SUPPORTING STRUCT TYPES
// ============================================

// Money represents monetary values with currency
type Money struct {
	Amount   float64      `gorm:"type:decimal(15,2)" json:"amount" validate:"min=0"`
	Currency CurrencyCode `gorm:"type:varchar(3)" json:"currency"`
}

// NewMoney creates a new Money instance
func NewMoney(amount float64, currency CurrencyCode) Money {
	return Money{Amount: amount, Currency: currency}
}

// Add adds two Money values (must have same currency)
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

// Subtract subtracts two Money values (must have same currency)
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}
	return Money{Amount: m.Amount - other.Amount, Currency: m.Currency}, nil
}

// CoverageLimits represents per-peril coverage limits
type CoverageLimits struct {
	PerOccurrence Money            `json:"per_occurrence"`
	Annual        Money            `json:"annual"`
	Lifetime      Money            `json:"lifetime"`
	PerPeril      map[string]Money `json:"per_peril,omitempty"`
	SubLimits     map[string]Money `json:"sub_limits,omitempty"`
}

// Region represents a geographic region for coverage
type Region struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	Country      string   `json:"country"`
	States       []string `json:"states,omitempty"`
	Cities       []string `json:"cities,omitempty"`
	PostalCodes  []string `json:"postal_codes,omitempty"`
	IsActive     bool     `json:"is_active"`
	Restrictions string   `json:"restrictions,omitempty"`
}

// Exclusion represents a policy exclusion
type Exclusion struct {
	Type          string     `json:"type"`
	Description   string     `json:"description"`
	Code          string     `json:"code"`
	Category      string     `json:"category"`
	EffectiveDate time.Time  `json:"effective_date"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	Reason        string     `json:"reason"`
}

// RiskFactor represents a factor contributing to risk assessment
type RiskFactor struct {
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Score       float64   `json:"score"`
	Weight      float64   `json:"weight"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InstallmentPlan represents a payment installment plan
type InstallmentPlan struct {
	PlanID           string           `json:"plan_id"`
	TotalAmount      Money            `json:"total_amount"`
	NumberOfPayments int              `json:"number_of_payments"`
	PaymentAmount    Money            `json:"payment_amount"`
	Frequency        PaymentFrequency `json:"frequency"`
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	Status           string           `json:"status"`
	PaidInstallments int              `json:"paid_installments"`
	NextDueDate      *time.Time       `json:"next_due_date,omitempty"`
}

// ComplianceCheck represents a regulatory compliance check
type ComplianceCheck struct {
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	CheckedAt time.Time  `json:"checked_at"`
	CheckedBy string     `json:"checked_by"`
	Result    string     `json:"result"`
	Notes     string     `json:"notes,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	Timestamp   time.Time         `json:"timestamp"`
	Action      string            `json:"action"`
	UserID      string            `json:"user_id"`
	UserName    string            `json:"user_name"`
	Changes     map[string]Change `json:"changes,omitempty"`
	IPAddress   string            `json:"ip_address"`
	UserAgent   string            `json:"user_agent"`
	Description string            `json:"description"`
}

// Change represents a field change in audit log
type Change struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

// ============================================
// PREMIUM CALCULATION CONSTANTS
// ============================================

const (
	// Discount factors
	NoClaimsDiscountRate      = 0.10 // 10% discount for no claims
	LoyaltyDiscountRate2Years = 0.05 // 5% discount for 2+ years
	LoyaltyDiscountRate5Years = 0.10 // 10% discount for 5+ years
	BundleDiscountRate        = 0.15 // 15% bundle discount
	CorporateDiscountRate     = 0.20 // 20% corporate discount

	// Loading factors
	HighRiskLoadingRate       = 0.30 // 30% loading for high risk
	VeryHighRiskLoadingRate   = 0.50 // 50% loading for very high risk
	FrequentClaimsLoadingRate = 0.20 // 20% loading for frequent claims

	// Premium adjustment factors
	MaxPremiumIncrease = 2.0 // 100% maximum increase
	MaxPremiumDiscount = 0.5 // 50% maximum discount

	// Cancellation penalty
	CancellationPenaltyRate = 0.10 // 10% cancellation penalty

	// Risk thresholds
	HighFraudRiskThreshold   = 80.0
	HighLossRatioThreshold   = 150.0
	HighChurnRiskThreshold   = 70.0
	ReviewLossRatioThreshold = 80.0
	FrequentClaimsThreshold  = 2

	// Grace periods
	DefaultPaymentGracePeriod = 30 // days
	MaxPaymentGracePeriod     = 60 // days

	// Renewal periods
	RenewalNotice60Days = 60
	RenewalNotice30Days = 30
	RenewalNotice15Days = 15
	RenewalNotice7Days  = 7
)

// ============================================
// VALIDATION CONSTANTS
// ============================================

const (
	MinCoverageAmount = 0.0
	MaxCoverageAmount = 10000000.0 // 10 million
	MinDeductible     = 0.0
	MaxDeductible     = 100000.0 // 100k
	MinRiskScore      = 0.0
	MaxRiskScore      = 100.0
	MinCoinsurance    = 0.0
	MaxCoinsurance    = 100.0
	MinNPS            = -100
	MaxNPS            = 100
	MinSatisfaction   = 0.0
	MaxSatisfaction   = 5.0
)

// ============================================
// ERROR DEFINITIONS
// ============================================

var (
	ErrInvalidPolicyStatus     = errors.New("invalid policy status")
	ErrInvalidPaymentFrequency = errors.New("invalid payment frequency")
	ErrInvalidRiskScore        = errors.New("risk score must be between 0 and 100")
	ErrInvalidCoverageAmount   = errors.New("coverage amount must be positive")
	ErrInvalidDeductible       = errors.New("deductible cannot exceed coverage amount")
	ErrInvalidDates            = errors.New("expiration date must be after effective date")
	ErrPolicyNotActive         = errors.New("policy is not active")
	ErrPolicyCancelled         = errors.New("policy has been cancelled")
	ErrPaymentOverdue          = errors.New("policy payment is overdue")
	ErrCoverageLimitExhausted  = errors.New("policy coverage limit exhausted")
	ErrHighFraudRisk           = errors.New("high fraud risk detected")
	ErrHighLossRatio           = errors.New("loss ratio too high")
	ErrOutstandingPayment      = errors.New("outstanding payment required")
	ErrNoActivePolicy          = errors.New("no active policy found")
	ErrDuplicatePolicy         = errors.New("duplicate policy exists")
	ErrInvalidCurrency         = errors.New("invalid currency code")
)
