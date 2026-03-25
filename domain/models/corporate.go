package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CorporateAccount represents corporate insurance accounts
type CorporateAccount struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CompanyName          string         `gorm:"not null" json:"company_name"`
	BusinessRegistration string         `gorm:"uniqueIndex;not null" json:"business_registration"`
	TaxID                string         `gorm:"uniqueIndex" json:"tax_id"`
	Industry             string         `json:"industry"`
	CompanySize          string         `json:"company_size"` // small, medium, large, enterprise
	Address              string         `gorm:"not null" json:"address"`
	City                 string         `gorm:"not null" json:"city"`
	State                string         `gorm:"not null" json:"state"`
	Country              string         `gorm:"not null" json:"country"`
	PostalCode           string         `json:"postal_code"`
	Phone                string         `gorm:"not null" json:"phone"`
	Email                string         `gorm:"not null" json:"email"`
	Website              string         `json:"website"`
	ContactPersonID      uuid.UUID      `gorm:"type:uuid;not null" json:"contact_person_id"`
	AccountManagerID     *uuid.UUID     `gorm:"type:uuid" json:"account_manager_id"`
	BillingContact       string         `json:"billing_contact"`                       // JSON object
	PaymentTerms         string         `gorm:"default:'net_30'" json:"payment_terms"` // net_30, net_60, prepaid
	CreditLimit          float64        `gorm:"default:0" json:"credit_limit"`
	DiscountRate         float64        `gorm:"default:0" json:"discount_rate"`
	Status               string         `gorm:"not null;default:'active'" json:"status"` // active, suspended, cancelled
	ContractStartDate    time.Time      `gorm:"not null" json:"contract_start_date"`
	ContractEndDate      time.Time      `json:"contract_end_date"`
	IsAutoRenew          bool           `gorm:"default:true" json:"is_auto_renew"`
	TotalEmployees       int            `gorm:"default:0" json:"total_employees"`
	CoveredEmployees     int            `gorm:"default:0" json:"covered_employees"`
	TotalDevices         int            `gorm:"default:0" json:"total_devices"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ContactPerson  User                `gorm:"foreignKey:ContactPersonID" json:"contact_person,omitempty"`
	AccountManager *User               `gorm:"foreignKey:AccountManagerID" json:"account_manager,omitempty"`
	Employees      []CorporateEmployee `gorm:"foreignKey:CorporateAccountID" json:"employees,omitempty"`
	Policies       []CorporatePolicy   `gorm:"foreignKey:CorporateAccountID" json:"policies,omitempty"`
	DeviceFleet    []FleetDevice       `gorm:"foreignKey:CorporateAccountID" json:"device_fleet,omitempty"`
}

// CorporateEmployee represents employees under corporate insurance
type CorporateEmployee struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CorporateAccountID uuid.UUID      `gorm:"type:uuid;not null" json:"corporate_account_id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	EmployeeID         string         `gorm:"not null" json:"employee_id"`
	Department         string         `json:"department"`
	JobTitle           string         `json:"job_title"`
	Manager            string         `json:"manager"`
	StartDate          time.Time      `gorm:"not null" json:"start_date"`
	EndDate            *time.Time     `json:"end_date"`
	Status             string         `gorm:"not null;default:'active'" json:"status"` // active, inactive, terminated
	DeviceAllowance    float64        `gorm:"default:0" json:"device_allowance"`
	CoverageLevel      string         `gorm:"default:'standard'" json:"coverage_level"` // basic, standard, premium
	IsManager          bool           `gorm:"default:false" json:"is_manager"`
	CanApproveDevices  bool           `gorm:"default:false" json:"can_approve_devices"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	CorporateAccount CorporateAccount `gorm:"foreignKey:CorporateAccountID" json:"corporate_account,omitempty"`
	User             User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AssignedDevices  []FleetDevice    `gorm:"foreignKey:AssignedEmployeeID" json:"assigned_devices,omitempty"`
}

// CorporatePolicy represents corporate insurance policies
type CorporatePolicy struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CorporateAccountID uuid.UUID      `gorm:"type:uuid;not null" json:"corporate_account_id"`
	PolicyNumber       string         `gorm:"uniqueIndex;not null" json:"policy_number"`
	PolicyType         string         `gorm:"not null" json:"policy_type"`   // fleet, byod, hybrid
	CoverageType       string         `gorm:"not null" json:"coverage_type"` // basic, comprehensive, premium
	MaxDevicesPerUser  int            `gorm:"default:2" json:"max_devices_per_user"`
	TotalDeviceLimit   int            `json:"total_device_limit"`
	DeductibleAmount   float64        `gorm:"not null" json:"deductible_amount"`
	CoverageAmount     float64        `gorm:"not null" json:"coverage_amount"`
	PremiumPerDevice   float64        `gorm:"not null" json:"premium_per_device"`
	BulkDiscount       float64        `gorm:"default:0" json:"bulk_discount"`
	StartDate          time.Time      `gorm:"not null" json:"start_date"`
	EndDate            time.Time      `gorm:"not null" json:"end_date"`
	Status             string         `gorm:"not null;default:'active'" json:"status"` // active, suspended, cancelled, expired
	Terms              string         `json:"terms"`                                   // JSON object
	Exclusions         string         `json:"exclusions"`                              // JSON array
	ApprovalWorkflow   string         `json:"approval_workflow"`                       // JSON object
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	CorporateAccount CorporateAccount `gorm:"foreignKey:CorporateAccountID" json:"corporate_account,omitempty"`
	CoveredDevices   []FleetDevice    `gorm:"foreignKey:CorporatePolicyID" json:"covered_devices,omitempty"`
}

// FleetDevice represents devices in corporate fleet management
type FleetDevice struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CorporateAccountID uuid.UUID      `gorm:"type:uuid;not null" json:"corporate_account_id"`
	CorporatePolicyID  *uuid.UUID     `gorm:"type:uuid" json:"corporate_policy_id"`
	DeviceID           uuid.UUID      `gorm:"type:uuid;not null" json:"device_id"`
	AssignedEmployeeID *uuid.UUID     `gorm:"type:uuid" json:"assigned_employee_id"`
	AssetTag           string         `gorm:"uniqueIndex" json:"asset_tag"`
	PurchaseOrder      string         `json:"purchase_order"`
	CostCenter         string         `json:"cost_center"`
	Department         string         `json:"department"`
	DeviceType         string         `gorm:"not null" json:"device_type"` // corporate_owned, byod, leased
	AssignmentDate     *time.Time     `json:"assignment_date"`
	ReturnDate         *time.Time     `json:"return_date"`
	Status             string         `gorm:"not null;default:'available'" json:"status"` // available, assigned, in_use, maintenance, retired
	Condition          string         `gorm:"default:'good'" json:"condition"`
	Location           string         `json:"location"`
	LastAuditDate      *time.Time     `json:"last_audit_date"`
	NextAuditDate      *time.Time     `json:"next_audit_date"`
	ComplianceStatus   string         `gorm:"default:'compliant'" json:"compliance_status"` // compliant, non_compliant, pending
	MDMEnrolled        bool           `gorm:"default:false" json:"mdm_enrolled"`
	MDMProfile         string         `json:"mdm_profile"`
	SecurityPolicies   string         `json:"security_policies"` // JSON array
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	CorporateAccount CorporateAccount   `gorm:"foreignKey:CorporateAccountID" json:"corporate_account,omitempty"`
	CorporatePolicy  *CorporatePolicy   `gorm:"foreignKey:CorporatePolicyID" json:"corporate_policy,omitempty"`
	Device           Device             `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	AssignedEmployee *CorporateEmployee `gorm:"foreignKey:AssignedEmployeeID" json:"assigned_employee,omitempty"`
}

// TableName methods
func (CorporateAccount) TableName() string {
	return "corporate_accounts"
}

func (CorporateEmployee) TableName() string {
	return "corporate_employees"
}

func (CorporatePolicy) TableName() string {
	return "corporate_policies"
}

func (FleetDevice) TableName() string {
	return "fleet_devices"
}

// BeforeCreate hooks
func (ca *CorporateAccount) BeforeCreate(tx *gorm.DB) error {
	if ca.ID == uuid.Nil {
		ca.ID = uuid.New()
	}
	return nil
}

func (ce *CorporateEmployee) BeforeCreate(tx *gorm.DB) error {
	if ce.ID == uuid.Nil {
		ce.ID = uuid.New()
	}
	return nil
}

func (cp *CorporatePolicy) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == uuid.Nil {
		cp.ID = uuid.New()
	}
	if cp.PolicyNumber == "" {
		cp.PolicyNumber = "CP-" + uuid.New().String()[:8]
	}
	return nil
}

func (fd *FleetDevice) BeforeCreate(tx *gorm.DB) error {
	if fd.ID == uuid.Nil {
		fd.ID = uuid.New()
	}
	return nil
}

// Business logic methods for CorporateAccount
func (ca *CorporateAccount) IsActive() bool {
	return ca.Status == "active"
}

func (ca *CorporateAccount) CalculateDiscount(baseAmount float64) float64 {
	return baseAmount * (1 - ca.DiscountRate)
}

func (ca *CorporateAccount) HasCreditLimit() bool {
	return ca.CreditLimit > 0
}

// Business logic methods for CorporateEmployee
func (ce *CorporateEmployee) IsActive() bool {
	return ce.Status == "active" && ce.EndDate == nil
}

func (ce *CorporateEmployee) CanManageDevices() bool {
	return ce.IsManager || ce.CanApproveDevices
}

// Business logic methods for FleetDevice
func (fd *FleetDevice) IsAvailable() bool {
	return fd.Status == "available"
}

func (fd *FleetDevice) AssignTo(employeeID uuid.UUID) {
	fd.AssignedEmployeeID = &employeeID
	fd.Status = "assigned"
	now := time.Now()
	fd.AssignmentDate = &now
}

func (fd *FleetDevice) Return() {
	fd.AssignedEmployeeID = nil
	fd.Status = "available"
	now := time.Now()
	fd.ReturnDate = &now
}

func (fd *FleetDevice) IsCompliant() bool {
	return fd.ComplianceStatus == "compliant"
}
