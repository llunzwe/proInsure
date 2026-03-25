package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyFamilyGroup represents family and group features for a policy
type PolicyFamilyGroup struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`
	// Removed CustomerID - access via Policy.CustomerID to avoid duplication

	// Family Plan
	FamilyPlan     bool       `gorm:"default:false" json:"family_plan"`
	FamilyPlanID   *uuid.UUID `gorm:"type:uuid;index" json:"family_plan_id"`
	FamilyMembers  int        `json:"family_members"`
	FamilyHeadID   uuid.UUID  `gorm:"type:uuid" json:"family_head_id"`
	FamilyRelation string     `json:"family_relation"` // head, spouse, child, parent

	// Family Benefits
	SharedDeductible bool    `gorm:"default:false" json:"shared_deductible"`
	SharedLimit      bool    `gorm:"default:false" json:"shared_limit"`
	FamilyDiscount   float64 `json:"family_discount"`
	CrossCoverage    bool    `gorm:"default:false" json:"cross_coverage"` // Members can claim for each other's devices

	// Group Features
	GroupDiscount     bool       `gorm:"default:false" json:"group_discount"`
	GroupID           *uuid.UUID `gorm:"type:uuid;index" json:"group_id"`
	GroupName         string     `json:"group_name"`
	GroupType         string     `gorm:"type:varchar(50)" json:"group_type"` // corporate, association, community
	GroupSize         int        `json:"group_size"`
	GroupDiscountRate float64    `json:"group_discount_rate"`

	// Special Discounts
	StudentDiscount bool       `gorm:"default:false" json:"student_discount"`
	StudentID       string     `json:"student_id"`
	StudentVerified bool       `gorm:"default:false" json:"student_verified"`
	StudentExpiry   *time.Time `json:"student_expiry"`

	SeniorDiscount bool `gorm:"default:false" json:"senior_discount"`
	SeniorAge      int  `json:"senior_age"`
	SeniorVerified bool `gorm:"default:false" json:"senior_verified"`

	MilitaryDiscount bool   `gorm:"default:false" json:"military_discount"`
	MilitaryID       string `json:"military_id"`
	MilitaryBranch   string `json:"military_branch"`
	MilitaryVerified bool   `gorm:"default:false" json:"military_verified"`

	EmployeeDiscount bool   `gorm:"default:false" json:"employee_discount"`
	EmployerName     string `json:"employer_name"`
	EmployeeID       string `json:"employee_id"`
	EmployeeVerified bool   `gorm:"default:false" json:"employee_verified"`

	// Discount Values
	StudentDiscountRate  float64 `gorm:"default:0.10" json:"student_discount_rate"`  // 10%
	SeniorDiscountRate   float64 `gorm:"default:0.15" json:"senior_discount_rate"`   // 15%
	MilitaryDiscountRate float64 `gorm:"default:0.20" json:"military_discount_rate"` // 20%
	EmployeeDiscountRate float64 `gorm:"default:0.12" json:"employee_discount_rate"` // 12%

	// Family Management
	CanAddMembers    bool       `gorm:"default:true" json:"can_add_members"`
	MaxFamilyMembers int        `gorm:"default:6" json:"max_family_members"`
	MemberAddedDate  *time.Time `json:"member_added_date"`

	// Group Management
	GroupJoinDate *time.Time `json:"group_join_date"`
	GroupExitDate *time.Time `json:"group_exit_date"`
	GroupStatus   string     `gorm:"type:varchar(20)" json:"group_status"` // active, pending, inactive

	// Status
	IsActive             bool       `gorm:"default:true" json:"is_active"`
	VerificationStatus   string     `gorm:"type:varchar(20)" json:"verification_status"` // verified, pending, failed
	LastVerificationDate *time.Time `json:"last_verification_date"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
	// Customer accessed via Policy.Customer relationship - CustomerID field removed to avoid duplication
}

// TableName returns the table name
func (PolicyFamilyGroup) TableName() string {
	return "policy_family_groups"
}

// IsFamilyPolicy checks if this is a family policy
func (pfg *PolicyFamilyGroup) IsFamilyPolicy() bool {
	return pfg.IsActive && pfg.FamilyPlan && pfg.FamilyMembers > 1
}

// GetGroupDiscountRate returns the applicable group discount rate
func (pfg *PolicyFamilyGroup) GetGroupDiscountRate() float64 {
	if !pfg.IsActive || !pfg.GroupDiscount || pfg.GroupID == nil {
		return 0
	}
	return pfg.GroupDiscountRate
}

// HasSpecialDiscount checks for any special discount
func (pfg *PolicyFamilyGroup) HasSpecialDiscount() bool {
	return pfg.StudentDiscount ||
		pfg.SeniorDiscount ||
		pfg.MilitaryDiscount ||
		pfg.EmployeeDiscount ||
		pfg.GroupDiscount
}

// GetTotalDiscountRate calculates total discount rate from all sources
func (pfg *PolicyFamilyGroup) GetTotalDiscountRate() float64 {
	if !pfg.IsActive {
		return 0
	}

	discount := 0.0

	// Add verified discounts
	if pfg.StudentDiscount && pfg.StudentVerified {
		discount += pfg.StudentDiscountRate
	}
	if pfg.SeniorDiscount && pfg.SeniorVerified {
		discount += pfg.SeniorDiscountRate
	}
	if pfg.MilitaryDiscount && pfg.MilitaryVerified {
		discount += pfg.MilitaryDiscountRate
	}
	if pfg.EmployeeDiscount && pfg.EmployeeVerified {
		discount += pfg.EmployeeDiscountRate
	}
	if pfg.GroupDiscount && pfg.GroupID != nil {
		discount += pfg.GroupDiscountRate
	}
	if pfg.FamilyPlan {
		discount += pfg.FamilyDiscount
	}

	// Cap at maximum discount
	if discount > 0.50 { // 50% maximum
		discount = 0.50
	}

	return discount
}

// CanAddFamilyMember checks if a new family member can be added
func (pfg *PolicyFamilyGroup) CanAddFamilyMember() bool {
	return pfg.IsActive &&
		pfg.FamilyPlan &&
		pfg.CanAddMembers &&
		pfg.FamilyMembers < pfg.MaxFamilyMembers
}

// IsEligibleForFamilyDiscount checks family discount eligibility
func (pfg *PolicyFamilyGroup) IsEligibleForFamilyDiscount() bool {
	return pfg.IsActive && pfg.FamilyPlan && pfg.FamilyMembers >= 2
}

// GetFamilyDiscountRate returns family discount based on member count
func (pfg *PolicyFamilyGroup) GetFamilyDiscountRate() float64 {
	if !pfg.IsEligibleForFamilyDiscount() {
		return 0
	}

	switch {
	case pfg.FamilyMembers >= 5:
		return 0.20 // 20% for 5+ members
	case pfg.FamilyMembers >= 4:
		return 0.15 // 15% for 4 members
	case pfg.FamilyMembers >= 3:
		return 0.10 // 10% for 3 members
	case pfg.FamilyMembers >= 2:
		return 0.05 // 5% for 2 members
	default:
		return 0
	}
}

// NeedsVerification checks if any discount needs verification
func (pfg *PolicyFamilyGroup) NeedsVerification() bool {
	// Student status needs annual verification
	if pfg.StudentDiscount && !pfg.StudentVerified {
		return true
	}
	if pfg.StudentDiscount && pfg.StudentExpiry != nil && time.Now().After(*pfg.StudentExpiry) {
		return true
	}

	// Military status needs verification
	if pfg.MilitaryDiscount && !pfg.MilitaryVerified {
		return true
	}

	// Employee status needs verification
	if pfg.EmployeeDiscount && !pfg.EmployeeVerified {
		return true
	}

	// Senior status needs one-time verification
	if pfg.SeniorDiscount && !pfg.SeniorVerified {
		return true
	}

	return false
}

// IsGroupActive checks if group membership is active
func (pfg *PolicyFamilyGroup) IsGroupActive() bool {
	return pfg.GroupDiscount &&
		pfg.GroupID != nil &&
		pfg.GroupStatus == "active"
}

// GetMemberRole returns the family member's role for a given customer
func (pfg *PolicyFamilyGroup) GetMemberRole(customerID uuid.UUID) string {
	if !pfg.FamilyPlan {
		return "individual"
	}

	if pfg.FamilyHeadID == customerID {
		return "family_head"
	}

	return pfg.FamilyRelation
}

// CalculateFamilyPremiumShare calculates this member's share of family premium
func (pfg *PolicyFamilyGroup) CalculateFamilyPremiumShare(totalPremium float64, customerID uuid.UUID) float64 {
	if !pfg.FamilyPlan || pfg.FamilyMembers <= 0 {
		return totalPremium
	}

	// Family head pays 40%, others split remaining 60%
	if pfg.GetMemberRole(customerID) == "family_head" {
		return totalPremium * 0.40
	}

	// Other members split remaining
	otherMembers := pfg.FamilyMembers - 1
	if otherMembers > 0 {
		return (totalPremium * 0.60) / float64(otherMembers)
	}

	return totalPremium
}
