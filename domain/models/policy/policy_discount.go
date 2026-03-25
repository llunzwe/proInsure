package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyDiscount represents discounts applied to a policy
type PolicyDiscount struct {
	database.BaseModel
	PolicyID     uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	DiscountCode string    `gorm:"not null;index" json:"discount_code"`
	DiscountType string    `gorm:"type:varchar(50)" json:"discount_type"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`

	// Amount
	DiscountPercent float64 `json:"discount_percent"`
	DiscountAmount  float64 `json:"discount_amount"`
	MaxDiscount     float64 `json:"max_discount"`

	// Eligibility
	EligibilityCriteria string `gorm:"type:json" json:"eligibility_criteria"`
	ValidationRules     string `gorm:"type:json" json:"validation_rules"`

	// Duration
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	IsRecurring       bool      `gorm:"default:false" json:"is_recurring"`
	RecurrencePattern string    `json:"recurrence_pattern"`

	// Application
	AppliedDate time.Time `json:"applied_date"`
	AppliedBy   uuid.UUID `gorm:"type:uuid" json:"applied_by"`
	PromoCode   string    `json:"promo_code"`
	Source      string    `json:"source"` // system, promo, loyalty, referral

	// Conditions
	MinPremium       float64    `json:"min_premium"`
	MinTenure        int        `json:"min_tenure"` // Months
	RequiresApproval bool       `gorm:"default:false" json:"requires_approval"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyDiscount) TableName() string {
	return "policy_discounts"
}

// IsValid checks if discount is currently valid
func (pd *PolicyDiscount) IsValid() bool {
	now := time.Now()
	return now.After(pd.StartDate) && now.Before(pd.EndDate)
}

// CalculateDiscountAmount calculates the actual discount amount
func (pd *PolicyDiscount) CalculateDiscountAmount(premium float64) float64 {
	amount := 0.0

	if pd.DiscountPercent > 0 {
		amount = premium * (pd.DiscountPercent / 100)
	} else {
		amount = pd.DiscountAmount
	}

	if pd.MaxDiscount > 0 && amount > pd.MaxDiscount {
		amount = pd.MaxDiscount
	}

	return amount
}

// IsEligible checks if policy meets eligibility criteria
func (pd *PolicyDiscount) IsEligible(premium float64, tenureMonths int) bool {
	if premium < pd.MinPremium {
		return false
	}

	if tenureMonths < pd.MinTenure {
		return false
	}

	return true
}

// IsApproved checks if discount is approved
func (pd *PolicyDiscount) IsApproved() bool {
	if !pd.RequiresApproval {
		return true
	}
	return pd.ApprovedBy != nil
}

// GetRemainingDays returns remaining days for the discount
func (pd *PolicyDiscount) GetRemainingDays() int {
	if time.Now().After(pd.EndDate) {
		return 0
	}
	return int(time.Until(pd.EndDate).Hours() / 24)
}

// IsExpiring checks if discount is expiring soon
func (pd *PolicyDiscount) IsExpiring() bool {
	daysRemaining := pd.GetRemainingDays()
	return daysRemaining > 0 && daysRemaining <= 7
}

// IsPromotional checks if this is a promotional discount
func (pd *PolicyDiscount) IsPromotional() bool {
	return pd.Source == "promo" && pd.PromoCode != ""
}

// IsLoyalty checks if this is a loyalty discount
func (pd *PolicyDiscount) IsLoyalty() bool {
	return pd.Source == "loyalty" || pd.DiscountType == "loyalty"
}

// IsReferral checks if this is a referral discount
func (pd *PolicyDiscount) IsReferral() bool {
	return pd.Source == "referral" || pd.DiscountType == "referral"
}

// IsSystemGenerated checks if discount was system-generated
func (pd *PolicyDiscount) IsSystemGenerated() bool {
	return pd.Source == "system"
}

// CanStack checks if discount can be stacked with others
func (pd *PolicyDiscount) CanStack() bool {
	// Promotional discounts typically don't stack
	if pd.IsPromotional() {
		return false
	}

	// Loyalty and system discounts can usually stack
	return pd.IsLoyalty() || pd.IsSystemGenerated()
}

// GetEffectivePercent returns the effective discount percentage
func (pd *PolicyDiscount) GetEffectivePercent(premium float64) float64 {
	if pd.DiscountPercent > 0 {
		return pd.DiscountPercent
	}

	if pd.DiscountAmount > 0 && premium > 0 {
		return (pd.DiscountAmount / premium) * 100
	}

	return 0
}

// ShouldRenew checks if discount should be renewed
func (pd *PolicyDiscount) ShouldRenew() bool {
	return pd.IsRecurring && pd.IsExpiring()
}
