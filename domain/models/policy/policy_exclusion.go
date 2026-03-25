package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyExclusion represents specific exclusions for a policy
type PolicyExclusion struct {
	database.BaseModel
	PolicyID      uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	ExclusionCode string    `gorm:"not null;index" json:"exclusion_code"`
	ExclusionType string    `gorm:"type:varchar(50)" json:"exclusion_type"`
	Category      string    `gorm:"type:varchar(50)" json:"category"`

	// Details
	Title         string `gorm:"not null" json:"title"`
	Description   string `json:"description"`
	DetailedTerms string `gorm:"type:json" json:"detailed_terms"`

	// Applicability
	EffectiveDate    time.Time  `json:"effective_date"`
	EndDate          *time.Time `json:"end_date"`
	IsWaivable       bool       `gorm:"default:false" json:"is_waivable"`
	WaiverConditions string     `gorm:"type:json" json:"waiver_conditions"`

	// Override
	CanOverride       bool       `gorm:"default:false" json:"can_override"`
	OverrideAuthLevel string     `json:"override_auth_level"`
	OverriddenBy      *uuid.UUID `gorm:"type:uuid" json:"overridden_by"`
	OverrideDate      *time.Time `json:"override_date"`
	OverrideReason    string     `json:"override_reason"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyExclusion) TableName() string {
	return "policy_exclusions"
}

// IsActive checks if the exclusion is currently active
func (pe *PolicyExclusion) IsActive() bool {
	now := time.Now()

	// Check if overridden
	if pe.OverriddenBy != nil {
		return false
	}

	// Check effective dates
	if now.Before(pe.EffectiveDate) {
		return false
	}

	if pe.EndDate != nil && now.After(*pe.EndDate) {
		return false
	}

	return true
}

// CanBeWaived checks if the exclusion can be waived
func (pe *PolicyExclusion) CanBeWaived() bool {
	return pe.IsWaivable && pe.WaiverConditions != ""
}

// CanBeOverridden checks if the exclusion can be overridden
func (pe *PolicyExclusion) CanBeOverridden() bool {
	return pe.CanOverride && pe.OverriddenBy == nil
}

// IsOverridden checks if the exclusion has been overridden
func (pe *PolicyExclusion) IsOverridden() bool {
	return pe.OverriddenBy != nil && pe.OverrideDate != nil
}

// IsTemporary checks if the exclusion is temporary
func (pe *PolicyExclusion) IsTemporary() bool {
	return pe.EndDate != nil
}

// GetRemainingDays returns remaining days for temporary exclusion
func (pe *PolicyExclusion) GetRemainingDays() int {
	if pe.EndDate == nil {
		return -1 // Permanent exclusion
	}

	if time.Now().After(*pe.EndDate) {
		return 0
	}

	return int(time.Until(*pe.EndDate).Hours() / 24)
}

// RequiresHighAuthorization checks if high-level authorization is required
func (pe *PolicyExclusion) RequiresHighAuthorization() bool {
	return pe.OverrideAuthLevel == "director" ||
		pe.OverrideAuthLevel == "executive" ||
		pe.OverrideAuthLevel == "senior_manager"
}

// AppliesToClaim checks if exclusion applies to a specific claim type
func (pe *PolicyExclusion) AppliesToClaim(claimType string) bool {
	if !pe.IsActive() {
		return false
	}

	// Check if exclusion type matches claim type
	return pe.ExclusionType == claimType || pe.Category == claimType
}

// Override marks the exclusion as overridden
func (pe *PolicyExclusion) Override(userID uuid.UUID, reason string) {
	now := time.Now()
	pe.OverriddenBy = &userID
	pe.OverrideDate = &now
	pe.OverrideReason = reason
}

// GetDaysActive returns how many days the exclusion has been active
func (pe *PolicyExclusion) GetDaysActive() int {
	if !pe.IsActive() {
		return 0
	}
	return int(time.Since(pe.EffectiveDate).Hours() / 24)
}

// IsCritical checks if this is a critical exclusion
func (pe *PolicyExclusion) IsCritical() bool {
	// Critical exclusions are not waivable and cannot be overridden
	return !pe.IsWaivable && !pe.CanOverride
}
