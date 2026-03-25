package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyClaimLimits represents claim limits and frequencies for a policy
type PolicyClaimLimits struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`

	// Annual Limits
	MaxClaimsPerYear   int       `gorm:"default:2" json:"max_claims_per_year"`
	ClaimsThisYear     int       `json:"claims_this_year"`
	ClaimYearStartDate time.Time `json:"claim_year_start_date"`

	// Claim History
	TotalClaimsFiled    int        `json:"total_claims_filed"`
	TotalClaimsApproved int        `json:"total_claims_approved"`
	TotalClaimsRejected int        `json:"total_claims_rejected"`
	TotalClaimsPaid     float64    `json:"total_claims_paid"`
	LastClaimDate       *time.Time `json:"last_claim_date"`
	LastClaimAmount     float64    `json:"last_claim_amount"`

	// Waiting Periods
	FirstClaimWaiting    int `gorm:"default:30" json:"first_claim_waiting"`    // Days
	BetweenClaimsWaiting int `gorm:"default:30" json:"between_claims_waiting"` // Days
	MajorClaimWaiting    int `gorm:"default:90" json:"major_claim_waiting"`    // Days for claims > 50% of coverage

	// Frequency Limits
	ClaimFrequencyLimit   int `gorm:"default:2" json:"claim_frequency_limit"` // Max claims in consecutive period
	ConsecutivePeriodDays int `gorm:"default:365" json:"consecutive_period_days"`

	// Per-Type Limits
	ScreenClaimsLimit int `gorm:"default:2" json:"screen_claims_limit"`
	WaterDamageLimit  int `gorm:"default:1" json:"water_damage_limit"`
	TheftClaimsLimit  int `gorm:"default:1" json:"theft_claims_limit"`

	// Current Usage
	ScreenClaimsUsed int `json:"screen_claims_used"`
	WaterDamageUsed  int `json:"water_damage_used"`
	TheftClaimsUsed  int `json:"theft_claims_used"`

	// Financial Limits
	MaxClaimAmount     float64 `json:"max_claim_amount"`
	AnnualClaimCap     float64 `json:"annual_claim_cap"`
	LifetimeClaimCap   float64 `json:"lifetime_claim_cap"`
	CurrentYearClaimed float64 `json:"current_year_claimed"`
	LifetimeClaimed    float64 `json:"lifetime_claimed"`

	// Status
	IsActive      bool   `gorm:"default:true" json:"is_active"`
	ClaimsBlocked bool   `gorm:"default:false" json:"claims_blocked"`
	BlockedReason string `json:"blocked_reason"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyClaimLimits) TableName() string {
	return "policy_claim_limits"
}

// CanFileNewClaim checks if a new claim can be filed
func (pcl *PolicyClaimLimits) CanFileNewClaim(claimType string, amount float64) (bool, string) {
	if !pcl.IsActive {
		return false, "Policy claim limits inactive"
	}

	if pcl.ClaimsBlocked {
		return false, pcl.BlockedReason
	}

	// Check annual limit
	if pcl.ClaimsThisYear >= pcl.MaxClaimsPerYear {
		return false, "Annual claim limit reached"
	}

	// Check financial limits
	if amount > pcl.MaxClaimAmount {
		return false, "Claim amount exceeds maximum"
	}

	if pcl.CurrentYearClaimed+amount > pcl.AnnualClaimCap {
		return false, "Would exceed annual claim cap"
	}

	if pcl.LifetimeClaimed+amount > pcl.LifetimeClaimCap {
		return false, "Would exceed lifetime claim cap"
	}

	// Check type-specific limits
	switch claimType {
	case "screen":
		if pcl.ScreenClaimsUsed >= pcl.ScreenClaimsLimit {
			return false, "Screen repair limit reached"
		}
	case "water":
		if pcl.WaterDamageUsed >= pcl.WaterDamageLimit {
			return false, "Water damage claim limit reached"
		}
	case "theft":
		if pcl.TheftClaimsUsed >= pcl.TheftClaimsLimit {
			return false, "Theft claim limit reached"
		}
	}

	// Check waiting periods
	reason := pcl.CheckWaitingPeriod(amount)
	if reason != "" {
		return false, reason
	}

	return true, ""
}

// CheckWaitingPeriod checks if waiting period requirements are met
func (pcl *PolicyClaimLimits) CheckWaitingPeriod(amount float64) string {
	now := time.Now()

	// First claim waiting period
	if pcl.TotalClaimsFiled == 0 {
		daysSinceStart := int(now.Sub(pcl.ClaimYearStartDate).Hours() / 24)
		if daysSinceStart < pcl.FirstClaimWaiting {
			return "First claim waiting period not met"
		}
	}

	// Between claims waiting period
	if pcl.LastClaimDate != nil {
		daysSinceLastClaim := int(now.Sub(*pcl.LastClaimDate).Hours() / 24)
		if daysSinceLastClaim < pcl.BetweenClaimsWaiting {
			return "Minimum time between claims not met"
		}

		// Major claim waiting period
		if amount > pcl.MaxClaimAmount*0.5 && daysSinceLastClaim < pcl.MajorClaimWaiting {
			return "Major claim waiting period not met"
		}
	}

	return ""
}

// GetRemainingClaims returns number of claims remaining this year
func (pcl *PolicyClaimLimits) GetRemainingClaims() int {
	return pcl.MaxClaimsPerYear - pcl.ClaimsThisYear
}

// GetClaimFrequency returns claims per year rate
func (pcl *PolicyClaimLimits) GetClaimFrequency() float64 {
	yearsActive := time.Since(pcl.ClaimYearStartDate).Hours() / 24 / 365
	if yearsActive < 1 {
		yearsActive = 1
	}
	return float64(pcl.TotalClaimsFiled) / yearsActive
}

// RecordClaim records a new claim
func (pcl *PolicyClaimLimits) RecordClaim(claimType string, amount float64, approved bool) {
	now := time.Now()
	pcl.TotalClaimsFiled++
	pcl.ClaimsThisYear++
	pcl.LastClaimDate = &now
	pcl.LastClaimAmount = amount

	if approved {
		pcl.TotalClaimsApproved++
		pcl.TotalClaimsPaid += amount
		pcl.CurrentYearClaimed += amount
		pcl.LifetimeClaimed += amount

		// Update type-specific counters
		switch claimType {
		case "screen":
			pcl.ScreenClaimsUsed++
		case "water":
			pcl.WaterDamageUsed++
		case "theft":
			pcl.TheftClaimsUsed++
		}
	} else {
		pcl.TotalClaimsRejected++
	}
}

// ResetAnnualLimits resets annual claim limits
func (pcl *PolicyClaimLimits) ResetAnnualLimits() {
	pcl.ClaimsThisYear = 0
	pcl.CurrentYearClaimed = 0
	pcl.ScreenClaimsUsed = 0
	pcl.WaterDamageUsed = 0
	pcl.TheftClaimsUsed = 0
	pcl.ClaimYearStartDate = time.Now()
}

// GetClaimApprovalRate returns the claim approval rate
func (pcl *PolicyClaimLimits) GetClaimApprovalRate() float64 {
	if pcl.TotalClaimsFiled == 0 {
		return 0
	}
	return float64(pcl.TotalClaimsApproved) / float64(pcl.TotalClaimsFiled) * 100
}

// IsHighFrequencyClaimant checks if this is a high-frequency claimant
func (pcl *PolicyClaimLimits) IsHighFrequencyClaimant() bool {
	return pcl.GetClaimFrequency() > 3 // More than 3 claims per year
}
