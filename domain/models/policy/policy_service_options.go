package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyServiceOptions represents replacement and repair service options for a policy
type PolicyServiceOptions struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`

	// Replacement Options
	ReplacementType    string `gorm:"type:varchar(20)" json:"replacement_type"` // new, refurbished, repair_first
	LoanerDevice       bool   `gorm:"default:false" json:"loaner_device"`
	ExpressReplacement bool   `gorm:"default:false" json:"express_replacement"`
	SameDayService     bool   `gorm:"default:false" json:"same_day_service"`

	// Service Network
	PreferredRepairNetwork   string `gorm:"type:json" json:"preferred_repair_network"`
	AuthorizedServiceCenters int    `json:"authorized_service_centers"`
	NetworkTier              string `gorm:"type:varchar(20)" json:"network_tier"` // premium, standard, basic

	// Convenience Services
	HomeServiceAvailable bool `gorm:"default:false" json:"home_service_available"`
	PickupDropService    bool `gorm:"default:false" json:"pickup_drop_service"`
	CourierService       bool `gorm:"default:false" json:"courier_service"`

	// Service Levels
	ServiceLevel         string `gorm:"type:varchar(20)" json:"service_level"` // platinum, gold, silver, bronze
	MaxServiceTime       int    `json:"max_service_time"`                      // Hours
	GuaranteedTurnaround bool   `gorm:"default:false" json:"guaranteed_turnaround"`

	// Service Limits
	FreeServicesPerYear  int     `gorm:"default:2" json:"free_services_per_year"`
	ServicesUsed         int     `json:"services_used"`
	ServiceFeeAfterLimit float64 `json:"service_fee_after_limit"`

	// Status
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	LastServiceDate *time.Time `json:"last_service_date"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyServiceOptions) TableName() string {
	return "policy_service_options"
}

// HasExpressService checks if express services are available
func (pso *PolicyServiceOptions) HasExpressService() bool {
	return pso.IsActive && (pso.ExpressReplacement || pso.SameDayService)
}

// HasHomeService checks if home service is available
func (pso *PolicyServiceOptions) HasHomeService() bool {
	return pso.IsActive && (pso.HomeServiceAvailable || pso.PickupDropService)
}

// GetServiceLevel returns the effective service level
func (pso *PolicyServiceOptions) GetServiceLevel() string {
	if !pso.IsActive {
		return "none"
	}

	if pso.SameDayService {
		return "same_day"
	} else if pso.ExpressReplacement {
		return "express"
	} else if pso.HomeServiceAvailable {
		return "premium"
	}
	return "standard"
}

// CanUseFreeService checks if free service is available
func (pso *PolicyServiceOptions) CanUseFreeService() bool {
	return pso.IsActive && pso.ServicesUsed < pso.FreeServicesPerYear
}

// GetServiceFee returns the service fee based on usage
func (pso *PolicyServiceOptions) GetServiceFee() float64 {
	if pso.CanUseFreeService() {
		return 0
	}
	return pso.ServiceFeeAfterLimit
}

// HasLoanerDevice checks if loaner device is available
func (pso *PolicyServiceOptions) HasLoanerDevice() bool {
	return pso.IsActive && pso.LoanerDevice
}

// GetMaxTurnaroundTime returns maximum service time in hours
func (pso *PolicyServiceOptions) GetMaxTurnaroundTime() int {
	if pso.SameDayService {
		return 24
	} else if pso.ExpressReplacement {
		return 48
	} else if pso.MaxServiceTime > 0 {
		return pso.MaxServiceTime
	}
	return 72 // Default 3 days
}

// IsPremiumService checks if this is premium service
func (pso *PolicyServiceOptions) IsPremiumService() bool {
	return pso.ServiceLevel == "platinum" || pso.ServiceLevel == "gold"
}

// UseService records a service usage
func (pso *PolicyServiceOptions) UseService() {
	pso.ServicesUsed++
	now := time.Now()
	pso.LastServiceDate = &now
}
