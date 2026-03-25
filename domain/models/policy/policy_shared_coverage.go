package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SharedCoverage represents shared coverage benefits for family/corporate plans
type SharedCoverage struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PlanType        string    `gorm:"not null" json:"plan_type"` // family, corporate
	PlanID          uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`
	CoverageType    string    `gorm:"not null" json:"coverage_type"` // deductible, limit, benefit
	SharedAmount    float64   `gorm:"not null" json:"shared_amount"`
	UsedAmount      float64   `gorm:"default:0" json:"used_amount"`
	RemainingAmount float64   `json:"remaining_amount"`
	ResetPeriod     string    `gorm:"default:'annual'" json:"reset_period"` // monthly, quarterly, annual
	LastReset       time.Time `json:"last_reset"`
	NextReset       time.Time `json:"next_reset"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// GroupDiscount represents group discounts for family/corporate plans
type GroupDiscount struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	DiscountName       string         `gorm:"not null" json:"discount_name"`
	DiscountType       string         `gorm:"not null" json:"discount_type"` // family, corporate, volume, loyalty
	MinimumMembers     int            `gorm:"default:2" json:"minimum_members"`
	MaximumMembers     int            `json:"maximum_members"`
	DiscountRate       float64        `gorm:"not null" json:"discount_rate"` // 0.0 to 1.0
	IsPercentage       bool           `gorm:"default:true" json:"is_percentage"`
	FixedAmount        float64        `json:"fixed_amount"`
	ApplicableProducts string         `json:"applicable_products"` // JSON array
	ValidFrom          time.Time      `gorm:"not null" json:"valid_from"`
	ValidTo            *time.Time     `json:"valid_to"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName methods
func (SharedCoverage) TableName() string {
	return "shared_coverages"
}

func (GroupDiscount) TableName() string {
	return "group_discounts"
}

// BeforeCreate hooks
func (sc *SharedCoverage) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}

func (gd *GroupDiscount) BeforeCreate(tx *gorm.DB) error {
	if gd.ID == uuid.Nil {
		gd.ID = uuid.New()
	}
	return nil
}

// Business logic methods for SharedCoverage
func (sc *SharedCoverage) UpdateRemainingAmount() {
	sc.RemainingAmount = sc.SharedAmount - sc.UsedAmount
}

func (sc *SharedCoverage) UseAmount(amount float64) bool {
	if sc.IsActive && sc.RemainingAmount >= amount {
		sc.UsedAmount += amount
		sc.UpdateRemainingAmount()
		return true
	}
	return false
}

func (sc *SharedCoverage) IsResetDue() bool {
	return time.Now().After(sc.NextReset)
}

func (sc *SharedCoverage) Reset() {
	sc.UsedAmount = 0
	sc.UpdateRemainingAmount()
	sc.LastReset = time.Now()

	// Calculate next reset based on period
	switch sc.ResetPeriod {
	case "monthly":
		sc.NextReset = sc.LastReset.AddDate(0, 1, 0)
	case "quarterly":
		sc.NextReset = sc.LastReset.AddDate(0, 3, 0)
	case "annual":
		sc.NextReset = sc.LastReset.AddDate(1, 0, 0)
	default:
		sc.NextReset = sc.LastReset.AddDate(1, 0, 0)
	}
}

func (sc *SharedCoverage) GetUtilizationPercentage() float64 {
	if sc.SharedAmount == 0 {
		return 0
	}
	return (sc.UsedAmount / sc.SharedAmount) * 100
}

// Business logic methods for GroupDiscount
func (gd *GroupDiscount) IsValid() bool {
	now := time.Now()
	if gd.ValidTo == nil {
		return gd.IsActive && now.After(gd.ValidFrom)
	}
	return gd.IsActive && now.After(gd.ValidFrom) && now.Before(*gd.ValidTo)
}

func (gd *GroupDiscount) CalculateDiscount(baseAmount float64, memberCount int) float64 {
	if !gd.IsValid() || memberCount < gd.MinimumMembers {
		return baseAmount
	}

	if gd.MaximumMembers > 0 && memberCount > gd.MaximumMembers {
		return baseAmount
	}

	if gd.IsPercentage {
		return baseAmount * (1 - gd.DiscountRate)
	}

	return baseAmount - gd.FixedAmount
}

func (gd *GroupDiscount) IsApplicableForType(discountType string) bool {
	return gd.IsActive && gd.DiscountType == discountType
}

func (gd *GroupDiscount) IsExpired() bool {
	if gd.ValidTo == nil {
		return false
	}
	return time.Now().After(*gd.ValidTo)
}

func (gd *GroupDiscount) GetDiscountAmount(baseAmount float64) float64 {
	if gd.IsPercentage {
		return baseAmount * gd.DiscountRate
	}
	return gd.FixedAmount
}
