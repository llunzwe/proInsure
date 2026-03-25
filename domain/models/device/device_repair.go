package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceRepair represents a device repair service record
type DeviceRepair struct {
	database.BaseModel
	DeviceID     uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TicketNumber string    `gorm:"uniqueIndex;not null" json:"ticket_number"`
	RepairStatus string    `gorm:"type:varchar(50);default:'pending'" json:"repair_status"` // pending, diagnosed, quoted, approved, in_progress, completed, cancelled
	Priority     string    `json:"priority"`                                                // low, normal, high, urgent

	// Issue & Diagnosis
	IssueReported   string     `json:"issue_reported"`
	IssueCategory   string     `json:"issue_category"` // hardware, software, physical_damage, water_damage, other
	DiagnosisDate   *time.Time `json:"diagnosis_date"`
	DiagnosisResult string     `json:"diagnosis_result"`
	DiagnosticFee   float64    `json:"diagnostic_fee"`
	TechnicianNotes string     `json:"technician_notes"`

	// Repair Details
	RepairType        string     `json:"repair_type"` // warranty, insurance, paid, recall
	RepairStartDate   *time.Time `json:"repair_start_date"`
	RepairEndDate     *time.Time `json:"repair_end_date"`
	EstimatedDuration int        `json:"estimated_duration"` // in hours
	ActualDuration    int        `json:"actual_duration"`    // in hours

	// Components & Parts
	ComponentsRepaired string `gorm:"type:json" json:"components_repaired"` // JSON array
	PartsUsed          string `gorm:"type:json" json:"parts_used"`          // JSON array with part details
	PartsSource        string `json:"parts_source"`                         // original, oem, third_party

	// Costs & Payment
	LaborCost         float64 `json:"labor_cost"`
	PartsCost         float64 `json:"parts_cost"`
	TotalCost         float64 `json:"total_cost"`
	QuotedAmount      float64 `json:"quoted_amount"`
	DiscountAmount    float64 `json:"discount_amount"`
	InsuranceCoverage float64 `json:"insurance_coverage"`
	CustomerPayment   float64 `json:"customer_payment"`
	Deductible        float64 `json:"deductible"`
	PaymentStatus     string  `json:"payment_status"` // pending, partial, paid, waived
	PaymentMethod     string  `json:"payment_method"`

	// Service Provider
	RepairVendorID   *uuid.UUID `gorm:"type:uuid" json:"repair_vendor_id"`
	RepairVendorName string     `json:"repair_vendor_name"`
	TechnicianID     *uuid.UUID `gorm:"type:uuid" json:"technician_id"`
	TechnicianName   string     `json:"technician_name"`
	ServiceLocation  string     `json:"service_location"` // in_store, on_site, mail_in
	ServiceCenterID  string     `json:"service_center_id"`

	// Warranty Information
	UnderWarranty      bool       `gorm:"default:false" json:"under_warranty"`
	WarrantyType       string     `json:"warranty_type"` // manufacturer, extended, repair
	WarrantyClaimID    *uuid.UUID `gorm:"type:uuid" json:"warranty_claim_id"`
	RepairWarrantyDays int        `json:"repair_warranty_days"`
	WarrantyExpiryDate *time.Time `json:"warranty_expiry_date"`

	// Insurance Claim
	HasInsuranceClaim bool       `gorm:"default:false" json:"has_insurance_claim"`
	InsuranceClaimID  *uuid.UUID `gorm:"type:uuid" json:"insurance_claim_id"`
	InsurancePolicyID *uuid.UUID `gorm:"type:uuid" json:"insurance_policy_id"`
	PreAuthRequired   bool       `gorm:"default:false" json:"pre_auth_required"`
	PreAuthCode       string     `json:"pre_auth_code"`

	// Quality & Satisfaction
	QualityCheck         bool    `gorm:"default:false" json:"quality_check"`
	QualityScore         float64 `json:"quality_score"`         // 0-100
	CustomerSatisfaction float64 `json:"customer_satisfaction"` // 0-5
	CustomerFeedback     string  `json:"customer_feedback"`
	RepairSuccessful     bool    `gorm:"default:false" json:"repair_successful"`
	ReworkRequired       bool    `gorm:"default:false" json:"rework_required"`
	ReworkCount          int     `json:"rework_count"`

	// Loaner Device
	LoanerProvided   bool       `gorm:"default:false" json:"loaner_provided"`
	LoanerDeviceID   *uuid.UUID `gorm:"type:uuid" json:"loaner_device_id"`
	LoanerReturnDate *time.Time `json:"loaner_return_date"`
	LoanerCondition  string     `json:"loaner_condition"`

	// Pickup & Delivery
	ServiceMethod   string     `json:"service_method"` // walk_in, pickup, mail_in
	PickupDate      *time.Time `json:"pickup_date"`
	PickupAddress   string     `gorm:"type:json" json:"pickup_address"`
	DeliveryDate    *time.Time `json:"delivery_date"`
	DeliveryAddress string     `gorm:"type:json" json:"delivery_address"`
	TrackingNumber  string     `json:"tracking_number"`
	ShippingCost    float64    `json:"shipping_cost"`

	// Documentation
	RepairInvoiceURL string `json:"repair_invoice_url"`
	ServiceReportURL string `json:"service_report_url"`
	PhotosBefore     string `gorm:"type:json" json:"photos_before"` // JSON array of URLs
	PhotosAfter      string `gorm:"type:json" json:"photos_after"`  // JSON array of URLs

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User, RepairVendor, Technician should be loaded via service layer using UserID, RepairVendorID, TechnicianID to avoid circular import
	// LoanerDevice should be loaded via service layer using LoanerDeviceID to avoid circular import
	// InsuranceClaim should be loaded via service layer using InsuranceClaimID to avoid circular import
	// InsurancePolicy should be loaded via service layer using InsurancePolicyID to avoid circular import
}

// RepairHistory tracks historical repairs for a device
type RepairHistory struct {
	database.BaseModel
	DeviceID          uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	RepairID          uuid.UUID `gorm:"type:uuid;not null" json:"repair_id"`
	RepairDate        time.Time `json:"repair_date"`
	RepairType        string    `json:"repair_type"`
	ComponentRepaired string    `json:"component_repaired"`
	RepairCost        float64   `json:"repair_cost"`
	WarrantyRepair    bool      `json:"warranty_repair"`
	RepairDuration    int       `json:"repair_duration"` // in hours

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	Repair DeviceRepair `gorm:"foreignKey:RepairID" json:"repair,omitempty"`
}

// RepairComponent represents individual component repairs
type RepairComponent struct {
	database.BaseModel
	RepairID         uuid.UUID `gorm:"type:uuid;not null" json:"repair_id"`
	ComponentName    string    `gorm:"not null" json:"component_name"`
	ComponentType    string    `json:"component_type"`
	PartNumber       string    `json:"part_number"`
	SerialNumber     string    `json:"serial_number"`
	IssueDescription string    `json:"issue_description"`
	RepairAction     string    `json:"repair_action"` // replaced, repaired, cleaned, adjusted
	PartCost         float64   `json:"part_cost"`
	LaborCost        float64   `json:"labor_cost"`
	WarrantyPeriod   int       `json:"warranty_period"` // in days

	// Relationships
	Repair DeviceRepair `gorm:"foreignKey:RepairID" json:"repair,omitempty"`
}

// TableName returns the table name
func (t *DeviceRepair) TableName() string {
	return "device_repairs"
}

func (t *RepairHistory) TableName() string {
	return "repair_histories"
}

func (t *RepairComponent) TableName() string {
	return "repair_components"
}

// BeforeCreate handles pre-creation logic
func (dr *DeviceRepair) BeforeCreate(tx *gorm.DB) error {
	if err := dr.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Generate ticket number if not provided
	if dr.TicketNumber == "" {
		dr.TicketNumber = dr.GenerateTicketNumber()
	}

	return nil
}

// GenerateTicketNumber generates a unique ticket number
func (dr *DeviceRepair) GenerateTicketNumber() string {
	return "REP-" + time.Now().Format("20060102") + "-" + dr.ID.String()[:8]
}

// CalculateTotalCost calculates the total repair cost
func (dr *DeviceRepair) CalculateTotalCost() {
	dr.TotalCost = dr.LaborCost + dr.PartsCost + dr.DiagnosticFee + dr.ShippingCost
	dr.CustomerPayment = dr.TotalCost - dr.InsuranceCoverage - dr.DiscountAmount
	if dr.CustomerPayment < 0 {
		dr.CustomerPayment = 0
	}
}

// SetPriority sets repair priority based on criteria
func (dr *DeviceRepair) SetPriority() {
	if dr.RepairType == "recall" {
		dr.Priority = "urgent"
	} else if dr.UnderWarranty {
		dr.Priority = "high"
	} else if dr.HasInsuranceClaim {
		dr.Priority = "normal"
	} else {
		dr.Priority = "normal"
	}
}

// CompleteRepair marks repair as completed
func (dr *DeviceRepair) CompleteRepair(successful bool) {
	dr.RepairStatus = "completed"
	now := time.Now()
	dr.RepairEndDate = &now
	dr.RepairSuccessful = successful

	if dr.RepairStartDate != nil {
		dr.ActualDuration = int(now.Sub(*dr.RepairStartDate).Hours())
	}

	// Set repair warranty expiry
	if successful && dr.RepairWarrantyDays > 0 {
		expiry := now.AddDate(0, 0, dr.RepairWarrantyDays)
		dr.WarrantyExpiryDate = &expiry
	}
}

// IsUnderWarranty checks if repair is under warranty
func (dr *DeviceRepair) IsUnderWarranty() bool {
	if dr.WarrantyExpiryDate == nil {
		return false
	}
	return time.Now().Before(*dr.WarrantyExpiryDate)
}

// RequiresApproval checks if repair needs customer approval
func (dr *DeviceRepair) RequiresApproval() bool {
	return dr.TotalCost > dr.QuotedAmount || dr.PreAuthRequired
}

// EstimateTurnaround estimates repair completion time
func (dr *DeviceRepair) EstimateTurnaround() time.Time {
	if dr.RepairStartDate == nil {
		return time.Now().Add(time.Duration(dr.EstimatedDuration) * time.Hour)
	}
	return dr.RepairStartDate.Add(time.Duration(dr.EstimatedDuration) * time.Hour)
}
