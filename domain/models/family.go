package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FamilyPlan represents family insurance plans that can be used across the system
type FamilyPlan struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PlanName             string         `gorm:"not null" json:"plan_name"`
	PrimaryUserID        uuid.UUID      `gorm:"type:uuid;not null" json:"primary_user_id"`
	MaxMembers           int            `gorm:"default:6" json:"max_members"`
	CurrentMembers       int            `gorm:"default:1" json:"current_members"`
	MaxDevicesPerMember  int            `gorm:"default:3" json:"max_devices_per_member"`
	SharedDeductible     float64        `gorm:"not null" json:"shared_deductible"`
	IndividualDeductible float64        `json:"individual_deductible"`
	TotalCoverageLimit   float64        `gorm:"not null" json:"total_coverage_limit"`
	PerMemberLimit       float64        `json:"per_member_limit"`
	GroupDiscount        float64        `gorm:"default:0.15" json:"group_discount"` // 15% default
	BasePremium          float64        `gorm:"not null" json:"base_premium"`
	CurrentPremium       float64        `gorm:"not null" json:"current_premium"`
	BillingCycle         string         `gorm:"default:'monthly'" json:"billing_cycle"`
	Status               string         `gorm:"not null;default:'active'" json:"status"` // active, suspended, cancelled
	StartDate            time.Time      `gorm:"not null" json:"start_date"`
	EndDate              *time.Time     `json:"end_date"`
	IsAutoRenew          bool           `gorm:"default:true" json:"is_auto_renew"`
	ParentalControls     string         `json:"parental_controls"` // JSON object
	SharedBenefits       string         `json:"shared_benefits"`   // JSON object
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	PrimaryUser User           `gorm:"foreignKey:PrimaryUserID" json:"primary_user,omitempty"`
	Members     []FamilyMember `gorm:"foreignKey:FamilyPlanID" json:"members,omitempty"`
}

// FamilyMember represents individual members of family insurance plans
type FamilyMember struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	FamilyPlanID      uuid.UUID      `gorm:"type:uuid;not null" json:"family_plan_id"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Relationship      string         `gorm:"not null" json:"relationship"` // spouse, child, parent, sibling, other
	JoinDate          time.Time      `gorm:"not null" json:"join_date"`
	LeaveDate         *time.Time     `json:"leave_date"`
	Status            string         `gorm:"not null;default:'active'" json:"status"` // active, inactive, removed
	IsMinor           bool           `gorm:"default:false" json:"is_minor"`
	GuardianUserID    *uuid.UUID     `gorm:"type:uuid" json:"guardian_user_id"`
	DeviceLimit       int            `gorm:"default:3" json:"device_limit"`
	CurrentDevices    int            `gorm:"default:0" json:"current_devices"`
	IndividualPremium float64        `json:"individual_premium"`
	CoverageLevel     string         `gorm:"default:'standard'" json:"coverage_level"` // basic, standard, premium
	Restrictions      string         `json:"restrictions"`                             // JSON object for parental controls
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	FamilyPlan FamilyPlan `gorm:"foreignKey:FamilyPlanID" json:"family_plan,omitempty"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Guardian   *User      `gorm:"foreignKey:GuardianUserID" json:"guardian,omitempty"`
}

// TableName methods
func (FamilyPlan) TableName() string {
	return "family_plans"
}

func (FamilyMember) TableName() string {
	return "family_members"
}

// BeforeCreate hooks
func (fp *FamilyPlan) BeforeCreate(tx *gorm.DB) error {
	if fp.ID == uuid.Nil {
		fp.ID = uuid.New()
	}
	return nil
}

func (fm *FamilyMember) BeforeCreate(tx *gorm.DB) error {
	if fm.ID == uuid.Nil {
		fm.ID = uuid.New()
	}
	return nil
}

// Business logic methods for FamilyPlan
func (fp *FamilyPlan) IsActive() bool {
	return fp.Status == "active"
}

func (fp *FamilyPlan) CanAddMember() bool {
	return fp.CurrentMembers < fp.MaxMembers
}

func (fp *FamilyPlan) AddMember() {
	if fp.CanAddMember() {
		fp.CurrentMembers++
		fp.UpdatePremium()
	}
}

func (fp *FamilyPlan) RemoveMember() {
	if fp.CurrentMembers > 1 {
		fp.CurrentMembers--
		fp.UpdatePremium()
	}
}

func (fp *FamilyPlan) UpdatePremium() {
	fp.CurrentPremium = fp.CalculatePremium()
}

func (fp *FamilyPlan) CalculatePremium() float64 {
	return fp.BasePremium * float64(fp.CurrentMembers) * (1 - fp.GroupDiscount)
}

func (fp *FamilyPlan) GetRemainingSlots() int {
	return fp.MaxMembers - fp.CurrentMembers
}

func (fp *FamilyPlan) HasParentalControls() bool {
	return fp.ParentalControls != "" && fp.ParentalControls != "{}"
}

// Business logic methods for FamilyMember
func (fm *FamilyMember) IsActive() bool {
	return fm.Status == "active" && fm.LeaveDate == nil
}

func (fm *FamilyMember) CanAddDevice() bool {
	return fm.CurrentDevices < fm.DeviceLimit
}

func (fm *FamilyMember) AddDevice() bool {
	if fm.CanAddDevice() {
		fm.CurrentDevices++
		return true
	}
	return false
}

func (fm *FamilyMember) RemoveDevice() bool {
	if fm.CurrentDevices > 0 {
		fm.CurrentDevices--
		return true
	}
	return false
}

func (fm *FamilyMember) RequiresGuardianApproval() bool {
	return fm.IsMinor && fm.GuardianUserID != nil
}

func (fm *FamilyMember) GetAvailableDeviceSlots() int {
	return fm.DeviceLimit - fm.CurrentDevices
}

func (fm *FamilyMember) IsAdult() bool {
	return !fm.IsMinor
}
