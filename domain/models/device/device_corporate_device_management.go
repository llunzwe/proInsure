package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// CorporateDepartment represents departments within a corporate account
type CorporateDepartment struct {
	database.BaseModel
	CorporateAccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"corporate_account_id"`
	DepartmentName     string    `gorm:"not null" json:"department_name"`
	DepartmentCode     string    `json:"department_code"`
	CostCenter         string    `json:"cost_center"`

	// Department Details
	ManagerName      string `json:"manager_name"`
	ManagerEmail     string `json:"manager_email"`
	EmployeeCount    int    `json:"employee_count"`
	DeviceAllocation int    `json:"device_allocation"`
	DevicesUsed      int    `json:"devices_used"`

	// Budget & Limits
	MonthlyBudget float64 `json:"monthly_budget"`
	BudgetUsed    float64 `json:"budget_used"`
	BudgetPeriod  string  `json:"budget_period"` // monthly, quarterly, annual
	DeviceLimit   int     `json:"device_limit"`

	// Policies
	ApprovalRequired bool    `json:"approval_required"`
	ApproverEmail    string  `json:"approver_email"`
	AutoApproveLimit float64 `json:"auto_approve_limit"`
	RestrictedModels string  `gorm:"type:json" json:"restricted_models"` // JSON array

	// Status
	Status           string     `json:"status"` // active, inactive, suspended
	CreatedDate      time.Time  `json:"created_date"`
	LastActivityDate *time.Time `json:"last_activity_date"`

	// Relationships
	// CorporateAccount should be loaded via service layer using CorporateAccountID to avoid circular import
	DeviceAssignments []CorporateDeviceAssignment `gorm:"foreignKey:DepartmentID" json:"device_assignments,omitempty"`
}

// CorporateDeviceFleet manages fleet of corporate devices
type CorporateDeviceFleet struct {
	database.BaseModel
	CorporateAccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"corporate_account_id"`
	FleetName          string    `gorm:"not null" json:"fleet_name"`
	FleetType          string    `json:"fleet_type"` // standard, executive, field, sales

	// Fleet Statistics
	TotalDevices    int `json:"total_devices"`
	ActiveDevices   int `json:"active_devices"`
	InactiveDevices int `json:"inactive_devices"`
	DamagedDevices  int `json:"damaged_devices"`
	LostDevices     int `json:"lost_devices"`

	// Device Categories
	DeviceCategories string `gorm:"type:json" json:"device_categories"` // JSON object with counts
	PreferredBrands  string `gorm:"type:json" json:"preferred_brands"`  // JSON array
	StandardModels   string `gorm:"type:json" json:"standard_models"`   // JSON array

	// Management
	RefreshCycle    int        `json:"refresh_cycle"` // months
	LastRefreshDate *time.Time `json:"last_refresh_date"`
	NextRefreshDate *time.Time `json:"next_refresh_date"`
	AutoRenewal     bool       `json:"auto_renewal"`

	// Cost Management
	TotalValue     float64 `json:"total_value"`
	MonthlyExpense float64 `json:"monthly_expense"`
	AnnualBudget   float64 `json:"annual_budget"`
	CostPerDevice  float64 `json:"cost_per_device"`

	// Compliance
	MDMIntegration   bool       `json:"mdm_integration"`
	MDMProvider      string     `json:"mdm_provider"`
	ComplianceRate   float64    `json:"compliance_rate"` // percentage
	SecurityPolicyID *uuid.UUID `gorm:"type:uuid" json:"security_policy_id"`

	// Relationships
	// CorporateAccount should be loaded via service layer using CorporateAccountID to avoid circular import
	DeviceAssignments []CorporateDeviceAssignment `gorm:"foreignKey:FleetID" json:"device_assignments,omitempty"`
	DevicePools       []CorporateDevicePool       `gorm:"foreignKey:FleetID" json:"device_pools,omitempty"`
}

// CorporateDeviceAssignment tracks device assignments to employees
type CorporateDeviceAssignment struct {
	database.BaseModel
	DeviceID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id"`
	CorporateAccountID uuid.UUID  `gorm:"type:uuid;not null;index" json:"corporate_account_id"`
	EmployeeID         uuid.UUID  `gorm:"type:uuid;index" json:"employee_id"`
	DepartmentID       *uuid.UUID `gorm:"type:uuid;index" json:"department_id"`
	FleetID            *uuid.UUID `gorm:"type:uuid;index" json:"fleet_id"`

	// Assignment Details
	AssignmentType     string     `json:"assignment_type"` // permanent, temporary, loaner
	AssignmentDate     time.Time  `json:"assignment_date"`
	ReturnDate         *time.Time `json:"return_date"`
	ExpectedReturnDate *time.Time `json:"expected_return_date"`
	AssignmentStatus   string     `json:"assignment_status"` // active, returned, lost, damaged

	// Employee Information
	EmployeeName     string `json:"employee_name"`
	EmployeeEmail    string `json:"employee_email"`
	EmployeePhone    string `json:"employee_phone"`
	EmployeeIDNumber string `json:"employee_id_number"`
	JobTitle         string `json:"job_title"`

	// Device Ownership
	OwnershipType       string  `json:"ownership_type"` // corporate, byod, hybrid
	PurchaseType        string  `json:"purchase_type"`  // company_purchased, employee_purchased, leased
	ReimbursementAmount float64 `json:"reimbursement_amount"`

	// Usage Policies
	UsagePolicyID      *uuid.UUID `gorm:"type:uuid" json:"usage_policy_id"`
	PersonalUseAllowed bool       `json:"personal_use_allowed"`
	InternationalUse   bool       `json:"international_use_allowed"`
	DataLimit          int        `json:"data_limit"` // GB per month

	// Condition Tracking
	InitialCondition  string `json:"initial_condition"`
	CurrentCondition  string `json:"current_condition"`
	ReturnCondition   string `json:"return_condition"`
	DamageReported    bool   `json:"damage_reported"`
	DamageDescription string `json:"damage_description"`

	// Compliance
	AgreementSigned     bool       `json:"agreement_signed"`
	AgreementDate       *time.Time `json:"agreement_date"`
	TrainingCompleted   bool       `json:"training_completed"`
	LastComplianceCheck *time.Time `json:"last_compliance_check"`

	// Relationships
	// Device, CorporateAccount, Employee should be loaded via service layer using their respective IDs to avoid circular import
	Department  *CorporateDepartment  `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Fleet       *CorporateDeviceFleet `gorm:"foreignKey:FleetID" json:"fleet,omitempty"`
	UsagePolicy *CorporateUsagePolicy `gorm:"foreignKey:UsagePolicyID" json:"usage_policy,omitempty"`
}

// CorporateBilling handles corporate billing and invoicing
type CorporateBilling struct {
	database.BaseModel
	CorporateAccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"corporate_account_id"`
	BillingCycle       string    `json:"billing_cycle"` // monthly, quarterly, annual

	// Billing Details
	BillingMethod   string    `json:"billing_method"` // invoice, direct_debit, credit_card
	PaymentTerms    int       `json:"payment_terms"`  // days
	InvoicePrefix   string    `json:"invoice_prefix"`
	NextInvoiceDate time.Time `json:"next_invoice_date"`

	// Financial Summary
	CurrentBalance    float64    `json:"current_balance"`
	OutstandingAmount float64    `json:"outstanding_amount"`
	LastPaymentAmount float64    `json:"last_payment_amount"`
	LastPaymentDate   *time.Time `json:"last_payment_date"`
	TotalPaidToDate   float64    `json:"total_paid_to_date"`

	// Cost Breakdown
	DeviceCharges     float64 `json:"device_charges"`
	ServiceCharges    float64 `json:"service_charges"`
	ClaimDeductibles  float64 `json:"claim_deductibles"`
	AdditionalCharges float64 `json:"additional_charges"`
	TaxAmount         float64 `json:"tax_amount"`

	// Billing Contact
	BillingContactName string `json:"billing_contact_name"`
	BillingEmail       string `json:"billing_email"`
	BillingPhone       string `json:"billing_phone"`
	PONumber           string `json:"po_number"`

	// Payment Status
	PaymentStatus string `json:"payment_status"` // current, overdue, suspended
	DaysOverdue   int    `json:"days_overdue"`
	CreditStatus  string `json:"credit_status"` // good, warning, hold

	// Relationships
	// CorporateAccount should be loaded via service layer using CorporateAccountID to avoid circular import
}

// CorporateDevicePool manages spare and replacement devices
type CorporateDevicePool struct {
	database.BaseModel
	CorporateAccountID uuid.UUID  `gorm:"type:uuid;not null;index" json:"corporate_account_id"`
	FleetID            *uuid.UUID `gorm:"type:uuid;index" json:"fleet_id"`
	PoolName           string     `json:"pool_name"`
	PoolType           string     `json:"pool_type"` // spare, replacement, emergency

	// Pool Statistics
	TotalDevices       int `json:"total_devices"`
	AvailableDevices   int `json:"available_devices"`
	AssignedDevices    int `json:"assigned_devices"`
	MaintenanceDevices int `json:"maintenance_devices"`

	// Device Management
	DeviceModels string `gorm:"type:json" json:"device_models"` // JSON object with counts
	MinimumStock int    `json:"minimum_stock"`
	ReorderPoint int    `json:"reorder_point"`
	AutoReorder  bool   `json:"auto_reorder"`

	// Location
	StorageLocation  string `json:"storage_location"`
	LocationAddress  string `gorm:"type:json" json:"location_address"` // JSON object
	AccessRestricted bool   `json:"access_restricted"`

	// Usage Tracking
	CheckoutCount       int        `json:"checkout_count"` // total historical
	AverageCheckoutDays float64    `json:"average_checkout_days"`
	LastCheckoutDate    *time.Time `json:"last_checkout_date"`

	// Status
	PoolStatus    string     `json:"pool_status"` // active, inactive, maintenance
	CreatedDate   time.Time  `json:"created_date"`
	LastAuditDate *time.Time `json:"last_audit_date"`

	// Relationships
	// CorporateAccount should be loaded via service layer using CorporateAccountID to avoid circular import
	Fleet *CorporateDeviceFleet `gorm:"foreignKey:FleetID" json:"fleet,omitempty"`
	// Pool devices should be loaded via service layer through many2many relationship
	PoolDeviceIDs []uuid.UUID `gorm:"type:uuid[]" json:"pool_device_ids,omitempty"`
}

// CorporateUsagePolicy defines usage policies for corporate devices
type CorporateUsagePolicy struct {
	database.BaseModel
	CorporateAccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"corporate_account_id"`
	PolicyName         string    `gorm:"not null" json:"policy_name"`
	PolicyVersion      string    `json:"policy_version"`

	// Usage Rules
	PersonalUseAllowed bool `json:"personal_use_allowed"`
	AppInstallAllowed  bool `json:"app_install_allowed"`
	CameraAllowed      bool `json:"camera_allowed"`
	ScreenshotAllowed  bool `json:"screenshot_allowed"`
	USBAllowed         bool `json:"usb_allowed"`

	// Network & Data
	WiFiOnly           bool `json:"wifi_only"`
	VPNRequired        bool `json:"vpn_required"`
	DataRoamingAllowed bool `json:"data_roaming_allowed"`
	HotspotAllowed     bool `json:"hotspot_allowed"`
	MonthlyDataLimit   int  `json:"monthly_data_limit"` // GB

	// Security Requirements
	PasswordRequired   bool `json:"password_required"`
	BiometricRequired  bool `json:"biometric_required"`
	EncryptionRequired bool `json:"encryption_required"`
	RemoteWipeEnabled  bool `json:"remote_wipe_enabled"`
	LocationTracking   bool `json:"location_tracking"`

	// App Restrictions
	BlacklistedApps string `gorm:"type:json" json:"blacklisted_apps"` // JSON array
	WhitelistedApps string `gorm:"type:json" json:"whitelisted_apps"` // JSON array
	RequiredApps    string `gorm:"type:json" json:"required_apps"`    // JSON array

	// Compliance
	ComplianceChecks string `gorm:"type:json" json:"compliance_checks"` // JSON array
	ViolationActions string `gorm:"type:json" json:"violation_actions"` // JSON object
	GracePeriodDays  int    `json:"grace_period_days"`

	// Status
	PolicyStatus  string     `json:"policy_status"` // active, draft, archived
	EffectiveDate time.Time  `json:"effective_date"`
	ExpiryDate    *time.Time `json:"expiry_date"`

	// Relationships
	// CorporateAccount should be loaded via service layer using CorporateAccountID to avoid circular import
	DeviceAssignments []CorporateDeviceAssignment `gorm:"foreignKey:UsagePolicyID" json:"device_assignments,omitempty"`
}

// Methods for CorporateDepartment
func (cd *CorporateDepartment) HasBudget() bool {
	return cd.BudgetUsed < cd.MonthlyBudget
}

func (cd *CorporateDepartment) GetBudgetUtilization() float64 {
	if cd.MonthlyBudget == 0 {
		return 0
	}
	return cd.BudgetUsed / cd.MonthlyBudget * 100
}

func (cd *CorporateDepartment) CanAddDevice() bool {
	return cd.Status == "active" && cd.DevicesUsed < cd.DeviceLimit
}

func (cd *CorporateDepartment) GetRemainingBudget() float64 {
	return cd.MonthlyBudget - cd.BudgetUsed
}

// Methods for CorporateDeviceFleet
func (cdf *CorporateDeviceFleet) NeedsRefresh() bool {
	if cdf.NextRefreshDate == nil {
		return false
	}
	return time.Now().After(*cdf.NextRefreshDate)
}

func (cdf *CorporateDeviceFleet) GetHealthScore() float64 {
	total := float64(cdf.TotalDevices)
	if total == 0 {
		return 100
	}
	healthy := float64(cdf.ActiveDevices)
	return (healthy / total) * 100
}

func (cdf *CorporateDeviceFleet) IsCompliant() bool {
	return cdf.MDMIntegration && cdf.ComplianceRate >= 90
}

func (cdf *CorporateDeviceFleet) GetLossRate() float64 {
	if cdf.TotalDevices == 0 {
		return 0
	}
	return float64(cdf.LostDevices+cdf.DamagedDevices) / float64(cdf.TotalDevices) * 100
}

// Methods for CorporateDeviceAssignment
func (cda *CorporateDeviceAssignment) IsActive() bool {
	return cda.AssignmentStatus == "active" && cda.ReturnDate == nil
}

func (cda *CorporateDeviceAssignment) IsOverdue() bool {
	if cda.ExpectedReturnDate == nil || cda.ReturnDate != nil {
		return false
	}
	return time.Now().After(*cda.ExpectedReturnDate)
}

func (cda *CorporateDeviceAssignment) IsCompliant() bool {
	return cda.AgreementSigned && cda.TrainingCompleted
}

func (cda *CorporateDeviceAssignment) IsCorporateOwned() bool {
	return cda.OwnershipType == "corporate"
}

func (cda *CorporateDeviceAssignment) GetAssignmentDuration() int {
	if cda.ReturnDate != nil {
		return int(cda.ReturnDate.Sub(cda.AssignmentDate).Hours() / 24)
	}
	return int(time.Since(cda.AssignmentDate).Hours() / 24)
}

// Methods for CorporateBilling
func (cb *CorporateBilling) IsOverdue() bool {
	return cb.PaymentStatus == "overdue" || cb.DaysOverdue > 0
}

func (cb *CorporateBilling) GetTotalCharges() float64 {
	return cb.DeviceCharges + cb.ServiceCharges + cb.ClaimDeductibles +
		cb.AdditionalCharges + cb.TaxAmount
}

func (cb *CorporateBilling) HasCreditAvailable() bool {
	return cb.CreditStatus == "good" && cb.OutstandingAmount == 0
}

func (cb *CorporateBilling) GetPaymentHealth() string {
	if cb.DaysOverdue > 60 {
		return "Critical"
	} else if cb.DaysOverdue > 30 {
		return "Warning"
	} else if cb.DaysOverdue > 0 {
		return "Late"
	}
	return "Good"
}

// Methods for CorporateDevicePool
func (cdp *CorporateDevicePool) NeedsReorder() bool {
	return cdp.AvailableDevices <= cdp.ReorderPoint && cdp.AutoReorder
}

func (cdp *CorporateDevicePool) GetUtilization() float64 {
	if cdp.TotalDevices == 0 {
		return 0
	}
	return float64(cdp.AssignedDevices) / float64(cdp.TotalDevices) * 100
}

func (cdp *CorporateDevicePool) HasAvailability() bool {
	return cdp.AvailableDevices > 0 && cdp.PoolStatus == "active"
}

func (cdp *CorporateDevicePool) IsCritical() bool {
	return cdp.AvailableDevices < cdp.MinimumStock
}

// Methods for CorporateUsagePolicy
func (cup *CorporateUsagePolicy) IsActive() bool {
	return cup.PolicyStatus == "active" && time.Now().After(cup.EffectiveDate)
}

func (cup *CorporateUsagePolicy) IsRestrictive() bool {
	return !cup.PersonalUseAllowed && !cup.AppInstallAllowed &&
		cup.VPNRequired && !cup.DataRoamingAllowed
}

func (cup *CorporateUsagePolicy) RequiresHighSecurity() bool {
	return cup.PasswordRequired && cup.BiometricRequired &&
		cup.EncryptionRequired && cup.RemoteWipeEnabled
}

func (cup *CorporateUsagePolicy) IsExpired() bool {
	if cup.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*cup.ExpiryDate)
}
