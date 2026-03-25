package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentMethod represents stored payment methods for users
type PaymentMethod struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Type          string         `gorm:"not null" json:"type"` // card, bank_transfer, mobile_money, direct_debit
	Provider      string         `json:"provider"`             // stripe, paypal, mpesa, bank_name
	Last4         string         `json:"last4,omitempty"`      // Last 4 digits of card/account
	Brand         string         `json:"brand,omitempty"`      // visa, mastercard, amex
	ExpiryMonth   int            `json:"expiry_month,omitempty"`
	ExpiryYear    int            `json:"expiry_year,omitempty"`
	HolderName    string         `json:"holder_name"`
	IsDefault     bool           `gorm:"default:false" json:"is_default"`
	ProviderToken string         `json:"-"`                              // Encrypted token from payment provider
	Metadata      string         `json:"metadata,omitempty"`             // JSON for additional data
	Status        string         `gorm:"default:'active'" json:"status"` // active, expired, suspended
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Payments []Payment `gorm:"foreignKey:PaymentMethodID" json:"payments,omitempty"`
}

// Payment represents a payment transaction
type Payment struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	PolicyID         *uuid.UUID `gorm:"type:uuid" json:"policy_id,omitempty"`
	ClaimID          *uuid.UUID `gorm:"type:uuid" json:"claim_id,omitempty"`
	PaymentMethodID  *uuid.UUID `gorm:"type:uuid" json:"payment_method_id,omitempty"`
	Type             string     `gorm:"not null" json:"type"` // premium, deductible, refund, payout, commission
	Amount           float64    `gorm:"not null" json:"amount"`
	Currency         string     `gorm:"not null;default:'USD'" json:"currency"`
	Status           string     `gorm:"not null;default:'pending'" json:"status"` // pending, processing, completed, failed, refunded
	Provider         string     `json:"provider"`                                 // stripe, paypal, bank_transfer
	ProviderTxnID    string     `json:"provider_txn_id,omitempty"`
	ProviderResponse string     `json:"provider_response,omitempty"` // JSON response from provider
	Reference        string     `gorm:"uniqueIndex" json:"reference"`
	Description      string     `json:"description"`
	FailureReason    string     `json:"failure_reason,omitempty"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
	Metadata         string     `json:"metadata,omitempty"` // JSON for additional data

	// Additional Payment Fields
	Fee              float64    `gorm:"default:0" json:"fee"`
	Tax              float64    `gorm:"default:0" json:"tax"`
	NetAmount        float64    `json:"net_amount"`
	RetryCount       int        `gorm:"default:0" json:"retry_count"`
	ScheduledFor     *time.Time `json:"scheduled_for"`
	Reconciled       bool       `gorm:"default:false" json:"reconciled"`
	ReconciledAt     *time.Time `json:"reconciled_at"`
	RefundedAmount   float64    `gorm:"default:0" json:"refunded_amount"`
	ChargebackStatus string     `json:"chargeback_status"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Policy        *Policy        `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Claim         *Claim         `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	PaymentMethod *PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
}

// Subscription represents recurring billing subscriptions
type Subscription struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	PolicyID           uuid.UUID      `gorm:"type:uuid;not null" json:"policy_id"`
	PaymentMethodID    uuid.UUID      `gorm:"type:uuid;not null" json:"payment_method_id"`
	Status             string         `gorm:"not null;default:'active'" json:"status"` // active, paused, cancelled, expired
	BillingCycle       string         `gorm:"not null" json:"billing_cycle"`           // monthly, quarterly, annual
	Amount             float64        `gorm:"not null" json:"amount"`
	Currency           string         `gorm:"not null;default:'USD'" json:"currency"`
	NextBillingDate    time.Time      `json:"next_billing_date"`
	LastBillingDate    *time.Time     `json:"last_billing_date,omitempty"`
	TrialEndDate       *time.Time     `json:"trial_end_date,omitempty"`
	CancelledAt        *time.Time     `json:"cancelled_at,omitempty"`
	CancellationReason string         `json:"cancellation_reason,omitempty"`
	ProviderSubID      string         `json:"provider_sub_id,omitempty"` // Stripe/PayPal subscription ID
	RetryCount         int            `gorm:"default:0" json:"retry_count"`
	FailedPayments     int            `gorm:"default:0" json:"failed_payments"`
	GracePeriodEndDate *time.Time     `json:"grace_period_end_date,omitempty"`
	AutoRenew          bool           `gorm:"default:true" json:"auto_renew"`
	Metadata           string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User           *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Policy         *Policy         `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	PaymentMethod  *PaymentMethod   `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	BillingHistory []BillingHistory `gorm:"foreignKey:SubscriptionID" json:"billing_history,omitempty"`
}

// BillingHistory tracks all billing events
type BillingHistory struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SubscriptionID uuid.UUID      `gorm:"type:uuid;not null" json:"subscription_id"`
	PaymentID      *uuid.UUID     `gorm:"type:uuid" json:"payment_id,omitempty"`
	Type           string         `gorm:"not null" json:"type"` // charge, refund, adjustment, credit
	Amount         float64        `gorm:"not null" json:"amount"`
	Currency       string         `gorm:"not null;default:'USD'" json:"currency"`
	Status         string         `gorm:"not null" json:"status"` // success, failed, pending
	Description    string         `json:"description"`
	InvoiceNumber  string         `json:"invoice_number,omitempty"`
	DueDate        time.Time      `json:"due_date"`
	PaidAt         *time.Time     `json:"paid_at,omitempty"`
	AttemptCount   int            `gorm:"default:1" json:"attempt_count"`
	FailureReason  string         `json:"failure_reason,omitempty"`
	Metadata       string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
	Payment      *Payment      `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}

// Invoice represents generated invoices
type Invoice struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	SubscriptionID *uuid.UUID     `gorm:"type:uuid" json:"subscription_id,omitempty"`
	PaymentID      *uuid.UUID     `gorm:"type:uuid" json:"payment_id,omitempty"`
	InvoiceNumber  string         `gorm:"uniqueIndex;not null" json:"invoice_number"`
	Status         string         `gorm:"not null;default:'draft'" json:"status"` // draft, sent, paid, overdue, cancelled
	Type           string         `gorm:"not null" json:"type"`                   // subscription, one_time, refund
	Subtotal       float64        `gorm:"not null" json:"subtotal"`
	Tax            float64        `gorm:"default:0" json:"tax"`
	Discount       float64        `gorm:"default:0" json:"discount"`
	Total          float64        `gorm:"not null" json:"total"`
	Currency       string         `gorm:"not null;default:'USD'" json:"currency"`
	DueDate        time.Time      `json:"due_date"`
	PaidAt         *time.Time     `json:"paid_at,omitempty"`
	SentAt         *time.Time     `json:"sent_at,omitempty"`
	LineItems      string         `json:"line_items"` // JSON array of items
	Notes          string         `json:"notes,omitempty"`
	Terms          string         `json:"terms,omitempty"`
	PDFUrl         string         `json:"pdf_url,omitempty"`
	Metadata       string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
	Payment      *Payment      `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}

// Commission represents agent/partner commissions
type Commission struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AgentID          uuid.UUID      `gorm:"type:uuid;not null" json:"agent_id"`
	PolicyID         *uuid.UUID     `gorm:"type:uuid" json:"policy_id,omitempty"`
	PaymentID        *uuid.UUID     `gorm:"type:uuid" json:"payment_id,omitempty"`
	Type             string         `gorm:"not null" json:"type"` // new_policy, renewal, referral
	Rate             float64        `gorm:"not null" json:"rate"` // Commission percentage
	BaseAmount       float64        `gorm:"not null" json:"base_amount"`
	CommissionAmount float64        `gorm:"not null" json:"commission_amount"`
	Currency         string         `gorm:"not null;default:'USD'" json:"currency"`
	Status           string         `gorm:"not null;default:'pending'" json:"status"` // pending, approved, paid, cancelled
	ApprovedBy       *uuid.UUID     `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time     `json:"approved_at,omitempty"`
	PaidAt           *time.Time     `json:"paid_at,omitempty"`
	PayoutReference  string         `json:"payout_reference,omitempty"`
	Notes            string         `json:"notes,omitempty"`
	Metadata         string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Agent   *User          `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Policy  *Policy `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Payment *Payment       `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}

// PromoCode represents promotional codes for discounts
type PromoCode struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Code               string         `gorm:"uniqueIndex;not null" json:"code"`
	Description        string         `json:"description"`
	Type               string         `gorm:"not null" json:"type"` // percentage, fixed_amount
	Value              float64        `gorm:"not null" json:"value"`
	Currency           string         `json:"currency,omitempty"` // For fixed_amount type
	MinAmount          float64        `json:"min_amount,omitempty"`
	MaxDiscount        float64        `json:"max_discount,omitempty"`
	UsageLimit         int            `json:"usage_limit,omitempty"`
	UsageCount         int            `gorm:"default:0" json:"usage_count"`
	UserLimit          int            `json:"user_limit,omitempty"` // Per user limit
	ValidFrom          time.Time      `json:"valid_from"`
	ValidUntil         time.Time      `json:"valid_until"`
	ApplicableProducts string         `json:"applicable_products,omitempty"` // JSON array of product IDs
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	Metadata           string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hooks for UUID generation
func (p *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.Reference == "" {
		p.Reference = p.GenerateReference()
	}
	p.NetAmount = p.CalculateNetAmount()
	return p.Validate()
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (b *BillingHistory) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.InvoiceNumber == "" {
		i.InvoiceNumber = "INV-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	return nil
}

func (c *Commission) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (p *PromoCode) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Payment validation and business logic methods

// Validate performs validation on the payment
func (p *Payment) Validate() error {
	if p.Amount <= 0 {
		return gorm.ErrInvalidData
	}
	if p.Type == "" {
		return gorm.ErrInvalidData
	}
	validTypes := []string{"premium", "deductible", "refund", "payout", "commission"}
	validType := false
	for _, t := range validTypes {
		if p.Type == t {
			validType = true
			break
		}
	}
	if !validType {
		return gorm.ErrInvalidData
	}
	return nil
}

// GenerateReference generates a unique payment reference
func (p *Payment) GenerateReference() string {
	timestamp := time.Now().Format("20060102150405")
	return "PAY-" + timestamp + "-" + uuid.New().String()[:8]
}

// CalculateNetAmount calculates the net amount after fees and taxes
func (p *Payment) CalculateNetAmount() float64 {
	return p.Amount - p.Fee - p.Tax + p.RefundedAmount
}

// CanRefund checks if the payment can be refunded
func (p *Payment) CanRefund() bool {
	return p.Status == "completed" && p.Type != "refund" && p.RefundedAmount < p.Amount
}

// MarkAsCompleted marks the payment as completed
func (p *Payment) MarkAsCompleted(providerTxnID string) {
	p.Status = "completed"
	p.ProviderTxnID = providerTxnID
	now := time.Now()
	p.ProcessedAt = &now
}

// MarkAsFailed marks the payment as failed
func (p *Payment) MarkAsFailed(reason string) {
	p.Status = "failed"
	p.FailureReason = reason
	p.RetryCount++
}

// ShouldRetry determines if payment should be retried
func (p *Payment) ShouldRetry() bool {
	return p.Status == "failed" && p.RetryCount < 3 && p.Type == "premium"
}

// PaymentMethod validation and business logic methods

// IsExpired checks if the payment method is expired
func (pm *PaymentMethod) IsExpired() bool {
	if pm.Type != "card" {
		return false
	}
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	return pm.ExpiryYear < currentYear || (pm.ExpiryYear == currentYear && pm.ExpiryMonth < currentMonth)
}

// MaskCardNumber masks all but the last 4 digits
func (pm *PaymentMethod) MaskCardNumber() string {
	if pm.Last4 == "" {
		return "****"
	}
	return "**** **** **** " + pm.Last4
}

// Subscription validation and business logic methods

// IsActive checks if subscription is currently active
func (s *Subscription) IsActive() bool {
	return s.Status == "active" && s.NextBillingDate.After(time.Now())
}

// ShouldRenew checks if subscription should be renewed
func (s *Subscription) ShouldRenew() bool {
	return s.AutoRenew && s.IsActive() && time.Now().After(s.NextBillingDate.AddDate(0, 0, -1))
}

// Cancel cancels the subscription
func (s *Subscription) Cancel(reason string) {
	s.Status = "cancelled"
	s.CancellationReason = reason
	now := time.Now()
	s.CancelledAt = &now
	s.AutoRenew = false
}

// Invoice validation and business logic methods

// IsOverdue checks if invoice is overdue
func (i *Invoice) IsOverdue() bool {
	return i.Status != "paid" && i.DueDate.Before(time.Now())
}

// MarkAsPaid marks invoice as paid
func (i *Invoice) MarkAsPaid() {
	i.Status = "paid"
	now := time.Now()
	i.PaidAt = &now
}

// PromoCode validation and business logic methods

// IsValid checks if promo code is valid for use
func (pc *PromoCode) IsValid() bool {
	now := time.Now()
	return pc.IsActive &&
		now.After(pc.ValidFrom) &&
		now.Before(pc.ValidUntil) &&
		(pc.UsageLimit == 0 || pc.UsageCount < pc.UsageLimit)
}

// CalculateDiscount calculates the discount amount
func (pc *PromoCode) CalculateDiscount(amount float64) float64 {
	if !pc.IsValid() {
		return 0
	}

	if amount < pc.MinAmount {
		return 0
	}

	var discount float64
	if pc.Type == "percentage" {
		discount = amount * (pc.Value / 100)
	} else {
		discount = pc.Value
	}

	if pc.MaxDiscount > 0 && discount > pc.MaxDiscount {
		discount = pc.MaxDiscount
	}

	return discount
}
