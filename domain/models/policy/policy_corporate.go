package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyCorporate represents corporate/business-specific features for a policy
type PolicyCorporate struct {
	database.BaseModel
	PolicyID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`
	CorporateAccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"corporate_account_id"`

	// Corporate Plan Details
	IsCorporatePolicy   bool   `gorm:"default:false" json:"is_corporate_policy"`
	CorporatePlanType   string `gorm:"type:varchar(50)" json:"corporate_plan_type"` // enterprise, business, small_business, startup
	BillingType         string `gorm:"type:varchar(20)" json:"billing_type"`        // monthly, quarterly, annual, custom
	InvoicingEnabled    bool   `gorm:"default:true" json:"invoicing_enabled"`
	PurchaseOrderNumber string `json:"purchase_order_number"`
	CostCenter          string `json:"cost_center"`

	// Department & Employee Management
	DepartmentID    *uuid.UUID `gorm:"type:uuid;index" json:"department_id"`
	DepartmentName  string     `json:"department_name"`
	EmployeeID      string     `json:"employee_id"`
	EmployeeName    string     `json:"employee_name"`
	EmployeeEmail   string     `json:"employee_email"`
	JobTitle        string     `json:"job_title"`
	ManagerApproval bool       `gorm:"default:false" json:"manager_approval"`
	ManagerID       string     `json:"manager_id"`

	// Device Management
	DeviceOwnership     string `gorm:"type:varchar(20)" json:"device_ownership"` // company_owned, byod, cope (corporate-owned-personally-enabled)
	MDMEnrolled         bool   `gorm:"default:false" json:"mdm_enrolled"`
	MDMProvider         string `json:"mdm_provider"`                                  // microsoft_intune, vmware_workspace, mobile_iron
	ComplianceStatus    string `gorm:"type:varchar(20)" json:"compliance_status"`     // compliant, non_compliant, pending
	SecurityPolicyLevel string `gorm:"type:varchar(20)" json:"security_policy_level"` // basic, standard, high, maximum

	// BYOD Specific
	BYODProgram             bool    `gorm:"default:false" json:"byod_program"`
	BYODReimbursement       float64 `json:"byod_reimbursement"`
	PersonalUseAllowed      bool    `gorm:"default:true" json:"personal_use_allowed"`
	WorkProfileRequired     bool    `gorm:"default:false" json:"work_profile_required"`
	ContainerizationEnabled bool    `gorm:"default:false" json:"containerization_enabled"`

	// Corporate Coverage Options
	BusinessDataProtection bool `gorm:"default:true" json:"business_data_protection"`
	CorporateAppCoverage   bool `gorm:"default:true" json:"corporate_app_coverage"`
	RemoteWipeEnabled      bool `gorm:"default:true" json:"remote_wipe_enabled"`
	DataLeakPrevention     bool `gorm:"default:false" json:"data_leak_prevention"`
	ComplianceMonitoring   bool `gorm:"default:false" json:"compliance_monitoring"`

	// Bulk & Fleet Management
	PolicyBundleID   *uuid.UUID `gorm:"type:uuid;index" json:"policy_bundle_id"` // References PolicyBundle for bulk policies
	FleetSize        int        `json:"fleet_size"`
	DevicesAllocated int        `json:"devices_allocated"`
	DevicesActive    int        `json:"devices_active"`

	// Corporate Discounts & Pricing
	CorporateDiscount  float64 `json:"corporate_discount"` // Percentage
	VolumeDiscount     float64 `json:"volume_discount"`    // Based on fleet size
	LoyaltyDiscount    float64 `json:"loyalty_discount"`   // For long-term contracts
	ContractLength     int     `json:"contract_length"`    // Months
	AutoRenewalEnabled bool    `gorm:"default:true" json:"auto_renewal_enabled"`

	// Budget & Limits
	DepartmentBudget      float64 `json:"department_budget"`
	BudgetUsed            float64 `json:"budget_used"`
	BudgetRemaining       float64 `json:"budget_remaining"`
	MaxDevicesPerEmployee int     `gorm:"default:2" json:"max_devices_per_employee"`
	MaxClaimAmount        float64 `json:"max_claim_amount"`

	// Service Level Agreement
	SLATier                 string `gorm:"type:varchar(20)" json:"sla_tier"` // bronze, silver, gold, platinum
	PrioritySupport         bool   `gorm:"default:false" json:"priority_support"`
	DedicatedAccountManager bool   `gorm:"default:false" json:"dedicated_account_manager"`
	ResponseTimeHours       int    `gorm:"default:24" json:"response_time_hours"`
	ResolutionTimeHours     int    `gorm:"default:72" json:"resolution_time_hours"`

	// Reporting & Analytics
	ReportingEnabled bool   `gorm:"default:true" json:"reporting_enabled"`
	ReportFrequency  string `gorm:"type:varchar(20)" json:"report_frequency"` // weekly, monthly, quarterly
	CustomReporting  bool   `gorm:"default:false" json:"custom_reporting"`
	APIAccessEnabled bool   `gorm:"default:false" json:"api_access_enabled"`

	// Status & Dates
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	ContractStartDate time.Time  `json:"contract_start_date"`
	ContractEndDate   *time.Time `json:"contract_end_date"`
	LastReviewDate    *time.Time `json:"last_review_date"`
	NextReviewDate    *time.Time `json:"next_review_date"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
	// CorporateAccount and Department relationships should be loaded via service layer
	// using CorporateAccountID and DepartmentID to avoid circular import
}

// TableName returns the table name
func (PolicyCorporate) TableName() string {
	return "policy_corporates"
}

// IsCorporate checks if this is a corporate policy
func (pc *PolicyCorporate) IsCorporate() bool {
	return pc.IsActive && pc.IsCorporatePolicy
}

// IsBYOD checks if this is a BYOD policy
func (pc *PolicyCorporate) IsBYOD() bool {
	return pc.DeviceOwnership == "byod" && pc.BYODProgram
}

// IsCompanyOwned checks if the device is company-owned
func (pc *PolicyCorporate) IsCompanyOwned() bool {
	return pc.DeviceOwnership == "company_owned" || pc.DeviceOwnership == "cope"
}

// GetTotalDiscount calculates total corporate discount
func (pc *PolicyCorporate) GetTotalDiscount() float64 {
	totalDiscount := pc.CorporateDiscount + pc.VolumeDiscount + pc.LoyaltyDiscount

	// Cap at maximum corporate discount
	if totalDiscount > 0.40 { // 40% maximum
		totalDiscount = 0.40
	}

	return totalDiscount
}

// CalculateVolumeDiscount calculates discount based on fleet size
func (pc *PolicyCorporate) CalculateVolumeDiscount() float64 {
	switch {
	case pc.FleetSize >= 1000:
		return 0.20 // 20% for 1000+ devices
	case pc.FleetSize >= 500:
		return 0.15 // 15% for 500-999 devices
	case pc.FleetSize >= 100:
		return 0.10 // 10% for 100-499 devices
	case pc.FleetSize >= 50:
		return 0.07 // 7% for 50-99 devices
	case pc.FleetSize >= 20:
		return 0.05 // 5% for 20-49 devices
	case pc.FleetSize >= 10:
		return 0.03 // 3% for 10-19 devices
	default:
		return 0
	}
}

// IsCompliant checks if the device meets compliance requirements
func (pc *PolicyCorporate) IsCompliant() bool {
	return pc.ComplianceStatus == "compliant" && pc.MDMEnrolled
}

// RequiresMDM checks if MDM enrollment is required
func (pc *PolicyCorporate) RequiresMDM() bool {
	return pc.IsCompanyOwned() || pc.SecurityPolicyLevel == "high" ||
		pc.SecurityPolicyLevel == "maximum"
}

// CanAddDevice checks if another device can be added
func (pc *PolicyCorporate) CanAddDevice() bool {
	if !pc.IsActive {
		return false
	}

	// Check budget constraints
	if pc.BudgetRemaining <= 0 {
		return false
	}

	// Check device allocation limits
	if pc.DevicesAllocated >= pc.FleetSize {
		return false
	}

	return true
}

// GetSLAResponseTime returns SLA response time in hours
func (pc *PolicyCorporate) GetSLAResponseTime() int {
	switch pc.SLATier {
	case "platinum":
		return 2
	case "gold":
		return 4
	case "silver":
		return 12
	case "bronze":
		return 24
	default:
		return 48
	}
}

// GetSLAResolutionTime returns SLA resolution time in hours
func (pc *PolicyCorporate) GetSLAResolutionTime() int {
	switch pc.SLATier {
	case "platinum":
		return 24
	case "gold":
		return 48
	case "silver":
		return 72
	case "bronze":
		return 96
	default:
		return 120
	}
}

// HasPremiumSupport checks if premium support is included
func (pc *PolicyCorporate) HasPremiumSupport() bool {
	return pc.PrioritySupport || pc.DedicatedAccountManager ||
		pc.SLATier == "gold" || pc.SLATier == "platinum"
}

// NeedsBudgetApproval checks if budget approval is needed
func (pc *PolicyCorporate) NeedsBudgetApproval(amount float64) bool {
	return amount > pc.MaxClaimAmount ||
		(pc.BudgetUsed+amount) > pc.DepartmentBudget
}

// UpdateBudgetUsage updates the budget usage
func (pc *PolicyCorporate) UpdateBudgetUsage(amount float64) {
	pc.BudgetUsed += amount
	pc.BudgetRemaining = pc.DepartmentBudget - pc.BudgetUsed
}

// IsContractExpiring checks if contract is expiring soon
func (pc *PolicyCorporate) IsContractExpiring() bool {
	if pc.ContractEndDate == nil {
		return false
	}

	daysUntilExpiry := int(time.Until(*pc.ContractEndDate).Hours() / 24)
	return daysUntilExpiry <= 90 // Expiring in 90 days or less
}

// GetCorporatePlanFeatures returns features for the corporate plan type
func (pc *PolicyCorporate) GetCorporatePlanFeatures() map[string]bool {
	features := make(map[string]bool)

	switch pc.CorporatePlanType {
	case "enterprise":
		features["unlimited_devices"] = true
		features["priority_support"] = true
		features["dedicated_manager"] = true
		features["custom_reporting"] = true
		features["api_access"] = true
		features["mdm_integration"] = true
		features["compliance_monitoring"] = true

	case "business":
		features["priority_support"] = true
		features["mdm_integration"] = true
		features["standard_reporting"] = true
		features["bulk_management"] = true

	case "small_business":
		features["standard_support"] = true
		features["basic_reporting"] = true
		features["bulk_management"] = true

	case "startup":
		features["basic_support"] = true
		features["self_service"] = true
	}

	return features
}

// GetSecurityLevel returns the security requirements for the policy
func (pc *PolicyCorporate) GetSecurityLevel() string {
	if pc.SecurityPolicyLevel != "" {
		return pc.SecurityPolicyLevel
	}

	// Determine based on other factors
	if pc.DataLeakPrevention && pc.ComplianceMonitoring {
		return "maximum"
	} else if pc.MDMEnrolled && pc.RemoteWipeEnabled {
		return "high"
	} else if pc.BusinessDataProtection {
		return "standard"
	}

	return "basic"
}
