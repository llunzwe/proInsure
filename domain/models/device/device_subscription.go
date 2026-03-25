package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceSubscription represents device-related subscription services
type DeviceSubscription struct {
	database.BaseModel
	DeviceID           uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID             uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	SubscriptionType   string    `gorm:"type:varchar(50);not null" json:"subscription_type"`           // care_plan, protection, upgrade, bundle
	SubscriptionStatus string    `gorm:"type:varchar(50);default:'active'" json:"subscription_status"` // active, paused, cancelled, expired

	// Subscription Plans
	PlanName string `json:"plan_name"`
	PlanTier string `json:"plan_tier"` // basic, standard, premium, enterprise
	PlanCode string `json:"plan_code"`

	// Subscription Period
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	RenewalDate  *time.Time `json:"renewal_date"`
	BillingCycle string     `json:"billing_cycle"` // monthly, quarterly, annual
	TrialPeriod  bool       `gorm:"default:false" json:"trial_period"`
	TrialEndDate *time.Time `json:"trial_end_date"`
	AutoRenewal  bool       `gorm:"default:true" json:"auto_renewal"`

	// Pricing
	MonthlyPrice       float64    `json:"monthly_price"`
	AnnualPrice        float64    `json:"annual_price"`
	CurrentPrice       float64    `json:"current_price"`
	DiscountPercentage float64    `json:"discount_percentage"`
	DiscountAmount     float64    `json:"discount_amount"`
	TotalPaid          float64    `json:"total_paid"`
	NextPaymentAmount  float64    `json:"next_payment_amount"`
	NextPaymentDate    *time.Time `json:"next_payment_date"`

	// Care Plan Features
	DeviceCarePlan     string `json:"device_care_plan"` // none, basic, premium, enterprise
	AccidentalDamage   bool   `gorm:"default:false" json:"accidental_damage"`
	TheftProtection    bool   `gorm:"default:false" json:"theft_protection"`
	ScreenProtection   bool   `gorm:"default:false" json:"screen_protection"`
	BatteryReplacement bool   `gorm:"default:false" json:"battery_replacement"`
	ExpressRepair      bool   `gorm:"default:false" json:"express_repair"`
	LoanerDevice       bool   `gorm:"default:false" json:"loaner_device"`

	// Coverage Limits
	MaxClaimsPerYear  int     `json:"max_claims_per_year"`
	ClaimsUsed        int     `json:"claims_used"`
	MaxRepairsPerYear int     `json:"max_repairs_per_year"`
	RepairsUsed       int     `json:"repairs_used"`
	CoverageAmount    float64 `json:"coverage_amount"`
	DeductibleAmount  float64 `json:"deductible_amount"`

	// Bundle Services
	CloudStorageGB    int  `json:"cloud_storage_gb"`
	VPNAccess         bool `gorm:"default:false" json:"vpn_access"`
	AntivirusIncluded bool `gorm:"default:false" json:"antivirus_included"`
	TechSupport247    bool `gorm:"default:false" json:"tech_support_24_7"`
	DataRecovery      bool `gorm:"default:false" json:"data_recovery"`

	// Upgrade Program
	AnnualUpgrade       bool       `gorm:"default:false" json:"annual_upgrade"`
	UpgradeEligibleDate *time.Time `json:"upgrade_eligible_date"`
	UpgradeDiscount     float64    `json:"upgrade_discount"`
	TradeInBonus        float64    `json:"trade_in_bonus"`

	// Family/Multi-Device
	FamilyPlan     bool    `gorm:"default:false" json:"family_plan"`
	MaxDevices     int     `json:"max_devices"`
	LinkedDevices  string  `gorm:"type:json" json:"linked_devices"` // JSON array of device IDs
	SharedBenefits bool    `gorm:"default:false" json:"shared_benefits"`
	FamilyDiscount float64 `json:"family_discount"`

	// Add-on Services
	AddOnServices string  `gorm:"type:json" json:"add_on_services"` // JSON array of services
	AddOnCost     float64 `json:"add_on_cost"`

	// Usage & Benefits
	BenefitsUsed    string  `gorm:"type:json" json:"benefits_used"` // JSON tracking of used benefits
	SavingsAchieved float64 `json:"savings_achieved"`
	LoyaltyPoints   int     `json:"loyalty_points"`
	RewardsTier     string  `json:"rewards_tier"`

	// Cancellation
	CancellationDate    *time.Time `json:"cancellation_date"`
	CancellationReason  string     `json:"cancellation_reason"`
	RefundAmount        float64    `json:"refund_amount"`
	EarlyTerminationFee float64    `json:"early_termination_fee"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// SubscriptionBenefit tracks individual benefits usage
type SubscriptionBenefit struct {
	database.BaseModel
	SubscriptionID  uuid.UUID  `gorm:"type:uuid;not null" json:"subscription_id"`
	BenefitType     string     `gorm:"not null" json:"benefit_type"` // repair, replacement, upgrade, support
	BenefitName     string     `json:"benefit_name"`
	UsedDate        time.Time  `json:"used_date"`
	ValueProvided   float64    `json:"value_provided"`
	RelatedEntityID *uuid.UUID `gorm:"type:uuid" json:"related_entity_id"` // repair_id, claim_id, etc.
	Notes           string     `json:"notes"`

	// Relationships
	Subscription DeviceSubscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
}

// SubscriptionPayment tracks subscription payments
type SubscriptionPayment struct {
	database.BaseModel
	SubscriptionID uuid.UUID `gorm:"type:uuid;not null" json:"subscription_id"`
	PaymentDate    time.Time `json:"payment_date"`
	Amount         float64   `json:"amount"`
	PaymentMethod  string    `json:"payment_method"`
	PaymentStatus  string    `json:"payment_status"` // success, failed, pending
	TransactionID  string    `json:"transaction_id"`
	InvoiceNumber  string    `json:"invoice_number"`
	BillingPeriod  string    `json:"billing_period"` // e.g., "Jan 2024"

	// Relationships
	Subscription DeviceSubscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
}

// TableName returns the table name
func (t *DeviceSubscription) TableName() string {
	return "device_subscriptions"
}

func (t *SubscriptionBenefit) TableName() string {
	return "subscription_benefits"
}

func (t *SubscriptionPayment) TableName() string {
	return "subscription_payments"
}

// BeforeCreate handles pre-creation logic
func (ds *DeviceSubscription) BeforeCreate(tx *gorm.DB) error {
	if err := ds.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// IsActive checks if subscription is currently active
func (ds *DeviceSubscription) IsActive() bool {
	if ds.SubscriptionStatus != "active" {
		return false
	}

	if ds.EndDate != nil && time.Now().After(*ds.EndDate) {
		return false
	}

	return true
}

// CanUseBenefit checks if a benefit can be used
func (ds *DeviceSubscription) CanUseBenefit(benefitType string) bool {
	if !ds.IsActive() {
		return false
	}

	switch benefitType {
	case "claim":
		return ds.ClaimsUsed < ds.MaxClaimsPerYear
	case "repair":
		return ds.RepairsUsed < ds.MaxRepairsPerYear
	case "upgrade":
		return ds.AnnualUpgrade && ds.UpgradeEligibleDate != nil &&
			time.Now().After(*ds.UpgradeEligibleDate)
	default:
		return true
	}
}

// UseBenefit records benefit usage
func (ds *DeviceSubscription) UseBenefit(benefitType string, value float64) {
	switch benefitType {
	case "claim":
		ds.ClaimsUsed++
	case "repair":
		ds.RepairsUsed++
	}

	ds.SavingsAchieved += value
}

// CalculateNextPayment calculates next payment amount
func (ds *DeviceSubscription) CalculateNextPayment() {
	baseAmount := ds.CurrentPrice

	if ds.FamilyPlan && ds.FamilyDiscount > 0 {
		baseAmount -= baseAmount * (ds.FamilyDiscount / 100)
	}

	if ds.DiscountPercentage > 0 {
		baseAmount -= baseAmount * (ds.DiscountPercentage / 100)
	} else if ds.DiscountAmount > 0 {
		baseAmount -= ds.DiscountAmount
	}

	baseAmount += ds.AddOnCost

	ds.NextPaymentAmount = baseAmount
}

// SetNextPaymentDate sets the next payment date based on billing cycle
func (ds *DeviceSubscription) SetNextPaymentDate() {
	if ds.NextPaymentDate == nil {
		ds.NextPaymentDate = &ds.StartDate
	}

	var nextDate time.Time
	switch ds.BillingCycle {
	case "monthly":
		nextDate = ds.NextPaymentDate.AddDate(0, 1, 0)
	case "quarterly":
		nextDate = ds.NextPaymentDate.AddDate(0, 3, 0)
	case "annual":
		nextDate = ds.NextPaymentDate.AddDate(1, 0, 0)
	default:
		nextDate = ds.NextPaymentDate.AddDate(0, 1, 0)
	}

	ds.NextPaymentDate = &nextDate
}

// CancelSubscription cancels the subscription
func (ds *DeviceSubscription) CancelSubscription(reason string) {
	ds.SubscriptionStatus = "cancelled"
	now := time.Now()
	ds.CancellationDate = &now
	ds.CancellationReason = reason
	ds.AutoRenewal = false

	// Calculate refund if applicable
	if ds.EndDate != nil && ds.EndDate.After(now) {
		remainingDays := ds.EndDate.Sub(now).Hours() / 24
		totalDays := ds.EndDate.Sub(ds.StartDate).Hours() / 24
		unusedPercentage := remainingDays / totalDays

		ds.RefundAmount = ds.CurrentPrice * unusedPercentage

		// Apply early termination fee if applicable
		if ds.EarlyTerminationFee > 0 {
			ds.RefundAmount -= ds.EarlyTerminationFee
			if ds.RefundAmount < 0 {
				ds.RefundAmount = 0
			}
		}
	}
}

// RenewSubscription renews the subscription for another period
func (ds *DeviceSubscription) RenewSubscription() {
	if !ds.AutoRenewal {
		return
	}

	// Reset usage counters
	ds.ClaimsUsed = 0
	ds.RepairsUsed = 0

	// Update dates
	ds.StartDate = *ds.RenewalDate

	switch ds.BillingCycle {
	case "monthly":
		newEnd := ds.StartDate.AddDate(0, 1, 0)
		ds.EndDate = &newEnd
	case "quarterly":
		newEnd := ds.StartDate.AddDate(0, 3, 0)
		ds.EndDate = &newEnd
	case "annual":
		newEnd := ds.StartDate.AddDate(1, 0, 0)
		ds.EndDate = &newEnd

		// Update upgrade eligibility for annual plans
		if ds.AnnualUpgrade {
			upgradeDate := ds.StartDate.AddDate(0, 11, 0) // Eligible after 11 months
			ds.UpgradeEligibleDate = &upgradeDate
		}
	}

	renewal := ds.EndDate
	ds.RenewalDate = renewal
}
