package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyRider represents additional coverage options for a policy
type PolicyRider struct {
	database.BaseModel
	PolicyID    uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	RiderCode   string    `gorm:"uniqueIndex;not null" json:"rider_code"`
	RiderType   string    `gorm:"type:varchar(50)" json:"rider_type"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Status      string    `gorm:"type:varchar(20);default:'active'" json:"status"`

	// Coverage
	CoverageAmount float64 `json:"coverage_amount"`
	Deductible     float64 `json:"deductible"`
	WaitingPeriod  int     `json:"waiting_period"` // Days
	BenefitPeriod  int     `json:"benefit_period"` // Days

	// Pricing
	Premium          float64 `json:"premium"`
	PremiumFrequency string  `gorm:"type:varchar(20)" json:"premium_frequency"`

	// Terms
	EffectiveDate  time.Time `json:"effective_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Terms          string    `gorm:"type:json" json:"terms"`
	Exclusions     string    `gorm:"type:json" json:"exclusions"`

	// Flags
	IsMandatory bool `gorm:"default:false" json:"is_mandatory"`
	IsWaivable  bool `gorm:"default:true" json:"is_waivable"`
	AutoRenew   bool `gorm:"default:true" json:"auto_renew"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyRider) TableName() string {
	return "policy_riders"
}

// IsActive checks if the rider is currently active
func (pr *PolicyRider) IsActive() bool {
	now := time.Now()
	return pr.Status == "active" &&
		now.After(pr.EffectiveDate) &&
		now.Before(pr.ExpirationDate)
}

// HasWaitingPeriodExpired checks if the waiting period has expired
func (pr *PolicyRider) HasWaitingPeriodExpired() bool {
	if pr.WaitingPeriod <= 0 {
		return true
	}
	waitingPeriodEnd := pr.EffectiveDate.AddDate(0, 0, pr.WaitingPeriod)
	return time.Now().After(waitingPeriodEnd)
}

// IsClaimable checks if the rider can be claimed
func (pr *PolicyRider) IsClaimable() bool {
	return pr.IsActive() && pr.HasWaitingPeriodExpired()
}

// GetRemainingBenefitDays returns remaining benefit days
func (pr *PolicyRider) GetRemainingBenefitDays() int {
	if pr.BenefitPeriod <= 0 {
		return -1 // Unlimited
	}

	daysSinceEffective := int(time.Since(pr.EffectiveDate).Hours() / 24)
	if daysSinceEffective >= pr.BenefitPeriod {
		return 0
	}

	return pr.BenefitPeriod - daysSinceEffective
}

// CalculateProRatedPremium calculates pro-rated premium for partial period
func (pr *PolicyRider) CalculateProRatedPremium(days int) float64 {
	if pr.Premium <= 0 || days <= 0 {
		return 0
	}

	var periodDays float64
	switch pr.PremiumFrequency {
	case "monthly":
		periodDays = 30
	case "quarterly":
		periodDays = 90
	case "semi-annual":
		periodDays = 180
	case "annual":
		periodDays = 365
	default:
		periodDays = 30
	}

	return (pr.Premium / periodDays) * float64(days)
}

// ShouldAutoRenew checks if the rider should auto-renew
func (pr *PolicyRider) ShouldAutoRenew() bool {
	if !pr.AutoRenew || pr.Status != "active" {
		return false
	}

	daysUntilExpiry := int(time.Until(pr.ExpirationDate).Hours() / 24)
	return daysUntilExpiry > 0 && daysUntilExpiry <= 30
}
