package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyLimit represents specific coverage limits for perils
type PolicyLimit struct {
	database.BaseModel
	PolicyID  uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	LimitCode string    `gorm:"not null;index" json:"limit_code"`
	PerilType string    `gorm:"type:varchar(50)" json:"peril_type"`
	LimitType string    `gorm:"type:varchar(50)" json:"limit_type"` // per_occurrence, aggregate, combined

	// Amounts
	LimitAmount     float64 `json:"limit_amount"`
	UsedAmount      float64 `json:"used_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Deductible      float64 `json:"deductible"`

	// Sub-limits
	SubLimitType   string  `json:"sub_limit_type"`
	SubLimitAmount float64 `json:"sub_limit_amount"`

	// Time-based limits
	LimitPeriod     string    `gorm:"type:varchar(20)" json:"limit_period"` // annual, lifetime, per_claim
	PeriodStartDate time.Time `json:"period_start_date"`
	PeriodEndDate   time.Time `json:"period_end_date"`

	// Conditions
	WaitingPeriod    int     `json:"waiting_period"` // Days
	Reinstatement    bool    `gorm:"default:false" json:"reinstatement"`
	ReinstatementFee float64 `json:"reinstatement_fee"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyLimit) TableName() string {
	return "policy_limits"
}

// GetUtilization returns the percentage of limit used
func (pl *PolicyLimit) GetUtilization() float64 {
	if pl.LimitAmount <= 0 {
		return 0
	}
	return (pl.UsedAmount / pl.LimitAmount) * 100
}

// IsExhausted checks if limit is exhausted
func (pl *PolicyLimit) IsExhausted() bool {
	return pl.RemainingAmount <= 0
}

// CanAccommodateClaim checks if limit can accommodate a claim amount
func (pl *PolicyLimit) CanAccommodateClaim(claimAmount float64) bool {
	effectiveAmount := claimAmount - pl.Deductible
	if effectiveAmount <= 0 {
		return true // Claim is within deductible
	}
	return pl.RemainingAmount >= effectiveAmount
}

// UpdateUsage updates the used and remaining amounts
func (pl *PolicyLimit) UpdateUsage(claimAmount float64) {
	effectiveAmount := claimAmount - pl.Deductible
	if effectiveAmount <= 0 {
		return // No impact on limits
	}

	pl.UsedAmount += effectiveAmount
	pl.RemainingAmount = pl.LimitAmount - pl.UsedAmount

	if pl.RemainingAmount < 0 {
		pl.RemainingAmount = 0
	}
}

// IsInWaitingPeriod checks if limit is in waiting period
func (pl *PolicyLimit) IsInWaitingPeriod() bool {
	if pl.WaitingPeriod <= 0 {
		return false
	}
	waitingPeriodEnd := pl.PeriodStartDate.AddDate(0, 0, pl.WaitingPeriod)
	return time.Now().Before(waitingPeriodEnd)
}

// ShouldResetLimit checks if limit should be reset (for annual limits)
func (pl *PolicyLimit) ShouldResetLimit() bool {
	if pl.LimitPeriod != "annual" {
		return false
	}
	return time.Now().After(pl.PeriodEndDate)
}

// ResetLimit resets the limit for a new period
func (pl *PolicyLimit) ResetLimit() {
	if pl.LimitPeriod == "annual" {
		pl.UsedAmount = 0
		pl.RemainingAmount = pl.LimitAmount
		pl.PeriodStartDate = time.Now()
		pl.PeriodEndDate = time.Now().AddDate(1, 0, 0)
	}
}

// CanReinstate checks if the limit can be reinstated
func (pl *PolicyLimit) CanReinstate() bool {
	return pl.Reinstatement && pl.IsExhausted()
}

// GetEffectiveLimit returns the effective limit considering sub-limits
func (pl *PolicyLimit) GetEffectiveLimit() float64 {
	if pl.SubLimitAmount > 0 && pl.SubLimitAmount < pl.LimitAmount {
		return pl.SubLimitAmount
	}
	return pl.LimitAmount
}

// IsAggregate checks if this is an aggregate limit
func (pl *PolicyLimit) IsAggregate() bool {
	return pl.LimitType == "aggregate"
}

// IsPerOccurrence checks if this is a per-occurrence limit
func (pl *PolicyLimit) IsPerOccurrence() bool {
	return pl.LimitType == "per_occurrence"
}
