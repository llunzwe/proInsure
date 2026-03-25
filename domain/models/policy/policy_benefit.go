package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyBenefit represents benefits included in a policy
type PolicyBenefit struct {
	database.BaseModel
	PolicyID    uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	BenefitCode string    `gorm:"not null;index" json:"benefit_code"`
	BenefitType string    `gorm:"type:varchar(50)" json:"benefit_type"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`

	// Coverage
	CoverageAmount float64 `json:"coverage_amount"`
	MaxClaims      int     `json:"max_claims"`
	ClaimsUsed     int     `json:"claims_used"`

	// Availability
	EffectiveDate  time.Time `json:"effective_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	WaitingPeriod  int       `json:"waiting_period"` // Days

	// Terms
	Terms      string `gorm:"type:json" json:"terms"`
	Conditions string `gorm:"type:json" json:"conditions"`

	// Usage
	LastUsedDate    *time.Time `json:"last_used_date"`
	TotalAmountUsed float64    `json:"total_amount_used"`

	// Flags
	IsActive           bool       `gorm:"default:true" json:"is_active"`
	IsOptional         bool       `gorm:"default:false" json:"is_optional"`
	RequiresActivation bool       `gorm:"default:false" json:"requires_activation"`
	ActivationDate     *time.Time `json:"activation_date"`

	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyBenefit) TableName() string {
	return "policy_benefits"
}

// IsAvailable checks if benefit is currently available
func (pb *PolicyBenefit) IsAvailable() bool {
	now := time.Now()
	waitingPeriodEnd := pb.EffectiveDate.AddDate(0, 0, pb.WaitingPeriod)

	// Check activation requirement
	if pb.RequiresActivation && pb.ActivationDate == nil {
		return false
	}

	return pb.IsActive &&
		now.After(waitingPeriodEnd) &&
		now.Before(pb.ExpirationDate) &&
		pb.ClaimsUsed < pb.MaxClaims
}

// GetRemainingClaims returns remaining claims for this benefit
func (pb *PolicyBenefit) GetRemainingClaims() int {
	return pb.MaxClaims - pb.ClaimsUsed
}

// HasRemainingCoverage checks if there's remaining coverage
func (pb *PolicyBenefit) HasRemainingCoverage() bool {
	if pb.CoverageAmount <= 0 {
		return true // Unlimited coverage
	}
	return pb.TotalAmountUsed < pb.CoverageAmount
}

// GetRemainingCoverage returns remaining coverage amount
func (pb *PolicyBenefit) GetRemainingCoverage() float64 {
	if pb.CoverageAmount <= 0 {
		return -1 // Unlimited
	}
	remaining := pb.CoverageAmount - pb.TotalAmountUsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsInWaitingPeriod checks if benefit is in waiting period
func (pb *PolicyBenefit) IsInWaitingPeriod() bool {
	if pb.WaitingPeriod <= 0 {
		return false
	}
	waitingPeriodEnd := pb.EffectiveDate.AddDate(0, 0, pb.WaitingPeriod)
	return time.Now().Before(waitingPeriodEnd)
}

// GetUtilization returns utilization percentage
func (pb *PolicyBenefit) GetUtilization() float64 {
	if pb.MaxClaims <= 0 && pb.CoverageAmount <= 0 {
		return 0 // No limits
	}

	// Calculate based on claims if limited
	if pb.MaxClaims > 0 {
		return (float64(pb.ClaimsUsed) / float64(pb.MaxClaims)) * 100
	}

	// Calculate based on coverage amount
	if pb.CoverageAmount > 0 {
		return (pb.TotalAmountUsed / pb.CoverageAmount) * 100
	}

	return 0
}

// UseBenefit records usage of the benefit
func (pb *PolicyBenefit) UseBenefit(amount float64) {
	pb.ClaimsUsed++
	pb.TotalAmountUsed += amount
	now := time.Now()
	pb.LastUsedDate = &now
}

// NeedsActivation checks if benefit needs activation
func (pb *PolicyBenefit) NeedsActivation() bool {
	return pb.RequiresActivation && pb.ActivationDate == nil
}

// Activate activates the benefit
func (pb *PolicyBenefit) Activate() {
	if pb.RequiresActivation {
		now := time.Now()
		pb.ActivationDate = &now
		pb.IsActive = true
	}
}

// GetDaysUntilExpiry returns days until benefit expires
func (pb *PolicyBenefit) GetDaysUntilExpiry() int {
	if time.Now().After(pb.ExpirationDate) {
		return 0
	}
	return int(time.Until(pb.ExpirationDate).Hours() / 24)
}

// GetDaysInWaitingPeriod returns remaining days in waiting period
func (pb *PolicyBenefit) GetDaysInWaitingPeriod() int {
	if pb.WaitingPeriod <= 0 {
		return 0
	}
	waitingPeriodEnd := pb.EffectiveDate.AddDate(0, 0, pb.WaitingPeriod)
	if time.Now().After(waitingPeriodEnd) {
		return 0
	}
	return int(time.Until(waitingPeriodEnd).Hours() / 24)
}
