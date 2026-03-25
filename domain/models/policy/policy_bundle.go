package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyBundle represents a bundle of multiple policies
type PolicyBundle struct {
	database.BaseModel
	BundleNumber       string     `gorm:"uniqueIndex;not null" json:"bundle_number"`
	BundleName         string     `gorm:"not null" json:"bundle_name"`
	BundleType         string     `gorm:"type:varchar(50)" json:"bundle_type"` // family, corporate, multi-device
	CustomerID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"customer_id"`
	CorporateAccountID *uuid.UUID `gorm:"type:uuid;index" json:"corporate_account_id"`
	Status             string     `gorm:"type:varchar(20);default:'active'" json:"status"`
	TotalPolicies      int        `json:"total_policies"`
	ActivePolicies     int        `json:"active_policies"`

	// Pricing
	BasePremium     float64 `json:"base_premium"`
	BundleDiscount  float64 `json:"bundle_discount"`
	DiscountPercent float64 `json:"discount_percent"`
	TotalPremium    float64 `json:"total_premium"`
	Currency        string  `gorm:"default:'USD'" json:"currency"`

	// Dates
	EffectiveDate  time.Time  `json:"effective_date"`
	ExpirationDate time.Time  `json:"expiration_date"`
	RenewalDate    *time.Time `json:"renewal_date"`

	// Configuration
	MaxPolicies       int  `gorm:"default:10" json:"max_policies"`
	MinPolicies       int  `gorm:"default:2" json:"min_policies"`
	AllowMixedTypes   bool `gorm:"default:true" json:"allow_mixed_types"`
	SharedDeductible  bool `gorm:"default:false" json:"shared_deductible"`
	SharedLimit       bool `gorm:"default:false" json:"shared_limit"`
	AutoAddNewDevices bool `gorm:"default:false" json:"auto_add_new_devices"`

	// Note: Policies relationship is handled in the main models package to avoid import cycles
	// Customer relationship should be loaded via service layer using CustomerID to avoid circular import
}

// TableName returns the table name
func (PolicyBundle) TableName() string {
	return "policy_bundles"
}

// CalculateBundleDiscount calculates discount based on number of policies
func (pb *PolicyBundle) CalculateBundleDiscount() float64 {
	discountRate := 0.0
	switch {
	case pb.ActivePolicies >= 5:
		discountRate = 0.20 // 20% for 5+ policies
	case pb.ActivePolicies >= 3:
		discountRate = 0.15 // 15% for 3-4 policies
	case pb.ActivePolicies >= 2:
		discountRate = 0.10 // 10% for 2 policies
	}
	return pb.BasePremium * discountRate
}

// CanAddPolicy checks if more policies can be added
func (pb *PolicyBundle) CanAddPolicy() bool {
	return pb.Status == "active" && pb.TotalPolicies < pb.MaxPolicies
}

// IsEligibleForBundle checks if bundle meets minimum requirements
func (pb *PolicyBundle) IsEligibleForBundle() bool {
	return pb.ActivePolicies >= pb.MinPolicies
}

// GetBundleUtilization returns the percentage of bundle capacity used
func (pb *PolicyBundle) GetBundleUtilization() float64 {
	if pb.MaxPolicies <= 0 {
		return 0
	}
	return (float64(pb.TotalPolicies) / float64(pb.MaxPolicies)) * 100
}

// IsActive checks if the bundle is currently active
func (pb *PolicyBundle) IsActive() bool {
	now := time.Now()
	return pb.Status == "active" &&
		now.After(pb.EffectiveDate) &&
		now.Before(pb.ExpirationDate)
}

// NeedsRenewal checks if the bundle needs renewal
func (pb *PolicyBundle) NeedsRenewal() bool {
	if pb.Status != "active" {
		return false
	}
	daysUntilExpiry := int(time.Until(pb.ExpirationDate).Hours() / 24)
	return daysUntilExpiry > 0 && daysUntilExpiry <= 30
}
