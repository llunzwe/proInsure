package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// BYODProgram manages Bring Your Own Device programs
type BYODProgram struct {
	database.BaseModel
	CorporateAccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"corporate_account_id"`
	ProgramName        string    `json:"program_name"`
	ProgramStatus      string    `json:"program_status"` // active, suspended, ended

	// Program Details
	EligibilityCriteria string `gorm:"type:json" json:"eligibility_criteria"` // JSON object
	SupportedDevices    string `gorm:"type:json" json:"supported_devices"`    // JSON array
	MinimumOSVersion    string `gorm:"type:json" json:"minimum_os_version"`   // JSON object per OS

	// Reimbursement
	ReimbursementType   string  `json:"reimbursement_type"` // fixed, percentage, tiered
	ReimbursementAmount float64 `json:"reimbursement_amount"`
	ReimbursementCap    float64 `json:"reimbursement_cap"`
	ReimbursementPeriod string  `json:"reimbursement_period"` // monthly, quarterly, one-time

	// Security Requirements
	MDMEnrollment            bool `json:"mdm_enrollment"`
	ContainerizationRequired bool `json:"containerization_required"`
	WorkProfileRequired      bool `json:"work_profile_required"`
	SecurityAuditRequired    bool `json:"security_audit_required"`

	// Support
	TechSupportIncluded bool `json:"tech_support_included"`
	RepairCoverage      bool `json:"repair_coverage"`
	ReplacementCoverage bool `json:"replacement_coverage"`

	// Enrollment
	TotalEnrolled      int        `json:"total_enrolled"`
	ActiveDevices      int        `json:"active_devices"`
	EnrollmentOpen     bool       `json:"enrollment_open"`
	EnrollmentDeadline *time.Time `json:"enrollment_deadline"`

	// Relationships
	// CorporateAccount should be loaded via service layer using CorporateAccountID to avoid circular import
	BYODDevices []BYODDevice `gorm:"foreignKey:ProgramID" json:"byod_devices,omitempty"`
}

// BYODDevice tracks individual BYOD devices
type BYODDevice struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ProgramID  uuid.UUID `gorm:"type:uuid;not null;index" json:"program_id"`
	EmployeeID uuid.UUID `gorm:"type:uuid;not null;index" json:"employee_id"`

	// Enrollment Details
	EnrollmentDate   time.Time  `json:"enrollment_date"`
	EnrollmentStatus string     `json:"enrollment_status"` // enrolled, pending, rejected
	ApprovalStatus   string     `json:"approval_status"`   // approved, pending, denied
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovalDate     *time.Time `json:"approval_date"`

	// Device Verification
	OwnershipVerified     bool    `json:"ownership_verified"`
	PurchaseProofProvided bool    `json:"purchase_proof_provided"`
	DeviceValue           float64 `json:"device_value"`

	// Work Profile
	WorkProfileInstalled bool       `json:"work_profile_installed"`
	MDMEnrolled          bool       `json:"mdm_enrolled"`
	ComplianceStatus     string     `json:"compliance_status"`
	LastComplianceCheck  *time.Time `json:"last_compliance_check"`

	// Reimbursement
	ReimbursementAmount   float64    `json:"reimbursement_amount"`
	LastReimbursementDate *time.Time `json:"last_reimbursement_date"`
	TotalReimbursed       float64    `json:"total_reimbursed"`

	// Usage
	DataUsageMB       int        `json:"data_usage_mb"`
	WorkAppsInstalled int        `json:"work_apps_installed"`
	LastActiveDate    *time.Time `json:"last_active_date"`

	// Relationships
	// Device and Employee should be loaded via service layer using DeviceID and EmployeeID to avoid circular import
	Program *BYODProgram `gorm:"foreignKey:ProgramID" json:"program,omitempty"`
}

// Methods for BYODProgram
func (bp *BYODProgram) IsActive() bool {
	return bp.ProgramStatus == "active" && bp.EnrollmentOpen
}

func (bp *BYODProgram) HasCapacity() bool {
	// Assuming there might be a limit, but not explicitly defined
	// This could be enhanced with a MaxEnrollment field if needed
	return true
}

func (bp *BYODProgram) GetUtilization() float64 {
	if bp.TotalEnrolled == 0 {
		return 0
	}
	return float64(bp.ActiveDevices) / float64(bp.TotalEnrolled) * 100
}

func (bp *BYODProgram) RequiresMDM() bool {
	return bp.MDMEnrollment || bp.ContainerizationRequired || bp.WorkProfileRequired
}

func (bp *BYODProgram) ProvidesSupport() bool {
	return bp.TechSupportIncluded || bp.RepairCoverage || bp.ReplacementCoverage
}

func (bp *BYODProgram) IsEnrollmentOpen() bool {
	if !bp.EnrollmentOpen {
		return false
	}
	if bp.EnrollmentDeadline == nil {
		return true
	}
	return time.Now().Before(*bp.EnrollmentDeadline)
}

// Methods for BYODDevice
func (bd *BYODDevice) IsActive() bool {
	return bd.EnrollmentStatus == "enrolled" && bd.ApprovalStatus == "approved"
}

func (bd *BYODDevice) IsCompliant() bool {
	return bd.WorkProfileInstalled && bd.MDMEnrolled &&
		bd.ComplianceStatus == "compliant"
}

func (bd *BYODDevice) NeedsReimbursement() bool {
	if bd.LastReimbursementDate == nil {
		return true
	}
	// Assuming monthly reimbursement
	return time.Since(*bd.LastReimbursementDate) > 30*24*time.Hour
}

func (bd *BYODDevice) IsVerified() bool {
	return bd.OwnershipVerified && bd.PurchaseProofProvided
}

func (bd *BYODDevice) GetEnrollmentDuration() int {
	return int(time.Since(bd.EnrollmentDate).Hours() / 24)
}

func (bd *BYODDevice) RequiresComplianceCheck() bool {
	if bd.LastComplianceCheck == nil {
		return true
	}
	// Check compliance every 30 days
	return time.Since(*bd.LastComplianceCheck) > 30*24*time.Hour
}

func (bd *BYODDevice) GetReimbursementStatus() string {
	if !bd.IsActive() {
		return "Ineligible"
	}
	if bd.NeedsReimbursement() {
		return "Pending"
	}
	return "Current"
}
