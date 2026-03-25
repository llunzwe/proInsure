package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceRefurbishment represents device refurbishment records
type DeviceRefurbishment struct {
	database.BaseModel
	DeviceID            uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	RefurbishmentStatus string    `gorm:"type:varchar(50);default:'pending'" json:"refurbishment_status"` // pending, in_progress, completed, failed
	RefurbishmentType   string    `gorm:"type:varchar(50)" json:"refurbishment_type"`                     // full, partial, cosmetic, functional

	// Refurbishment Details
	RefurbishmentDate   *time.Time `json:"refurbishment_date"`
	RefurbishmentVendor string     `json:"refurbishment_vendor"`
	VendorID            *uuid.UUID `gorm:"type:uuid" json:"vendor_id"`
	TechnicianID        *uuid.UUID `gorm:"type:uuid" json:"technician_id"`
	TechnicianName      string     `json:"technician_name"`

	// Quality & Grading
	PreRefurbGrade       string  `gorm:"type:varchar(10)" json:"pre_refurb_grade"` // A+, A, B, C, D
	PostRefurbGrade      string  `gorm:"type:varchar(10)" json:"post_refurb_grade"`
	CosmeticGrade        string  `gorm:"type:varchar(10)" json:"cosmetic_grade"`
	FunctionalGrade      string  `gorm:"type:varchar(10)" json:"functional_grade"`
	QualityScore         float64 `json:"quality_score"` // 0-100
	QualityCertificateID string  `json:"quality_certificate_id"`

	// Components & Parts
	ComponentsReplaced string `gorm:"type:json" json:"components_replaced"` // JSON array of replaced components
	PartsUsedType      string `json:"parts_used_type"`                      // original, oem, aftermarket
	TotalPartsCount    int    `json:"total_parts_count"`

	// Costs
	LaborCost         float64 `json:"labor_cost"`
	PartsCost         float64 `json:"parts_cost"`
	TotalRefurbCost   float64 `json:"total_refurb_cost"`
	MarketValueBefore float64 `json:"market_value_before"`
	MarketValueAfter  float64 `json:"market_value_after"`

	// Warranty
	RefurbWarrantyMonths int        `json:"refurb_warranty_months"`
	WarrantyStartDate    *time.Time `json:"warranty_start_date"`
	WarrantyEndDate      *time.Time `json:"warranty_end_date"`
	WarrantyTerms        string     `gorm:"type:json" json:"warranty_terms"`

	// History & Documentation
	PreviousOwnerCount   int        `json:"previous_owner_count"`
	OriginalPurchaseDate *time.Time `json:"original_purchase_date"`
	RefurbCertificate    string     `json:"refurb_certificate"`             // URL or document ID
	InspectionReport     string     `json:"inspection_report"`              // URL or document ID
	PhotosBefore         string     `gorm:"type:json" json:"photos_before"` // JSON array of URLs
	PhotosAfter          string     `gorm:"type:json" json:"photos_after"`  // JSON array of URLs

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// Vendor should be loaded via service layer using VendorID to avoid circular import
	// Technician should be loaded via service layer using TechnicianID to avoid circular import
}

// RefurbishmentComponent represents individual component replacements
type RefurbishmentComponent struct {
	database.BaseModel
	RefurbishmentID uuid.UUID `gorm:"type:uuid;not null" json:"refurbishment_id"`
	ComponentName   string    `gorm:"not null" json:"component_name"`
	ComponentType   string    `json:"component_type"` // screen, battery, camera, etc.
	PartNumber      string    `json:"part_number"`
	SerialNumber    string    `json:"serial_number"`
	IsNewPart       bool      `gorm:"default:true" json:"is_new_part"`
	PartCondition   string    `json:"part_condition"` // new, refurbished, used
	PartCost        float64   `json:"part_cost"`
	LaborCost       float64   `json:"labor_cost"`
	WarrantyMonths  int       `json:"warranty_months"`
	ReplacementDate time.Time `json:"replacement_date"`
	ReplacedBy      string    `json:"replaced_by"` // Technician name

	// Relationships
	Refurbishment DeviceRefurbishment `gorm:"foreignKey:RefurbishmentID" json:"refurbishment,omitempty"`
}

// TableName returns the table name
func (t *DeviceRefurbishment) TableName() string {
	return "device_refurbishments"
}

func (t *RefurbishmentComponent) TableName() string {
	return "refurbishment_components"
}

// BeforeCreate handles pre-creation logic
func (dr *DeviceRefurbishment) BeforeCreate(tx *gorm.DB) error {
	if err := dr.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// CalculateTotalCost calculates total refurbishment cost
func (dr *DeviceRefurbishment) CalculateTotalCost() {
	dr.TotalRefurbCost = dr.LaborCost + dr.PartsCost
}

// CalculateROI calculates return on investment
func (dr *DeviceRefurbishment) CalculateROI() float64 {
	if dr.TotalRefurbCost == 0 {
		return 0
	}
	return ((dr.MarketValueAfter - dr.MarketValueBefore - dr.TotalRefurbCost) / dr.TotalRefurbCost) * 100
}

// SetWarrantyPeriod sets warranty start and end dates
func (dr *DeviceRefurbishment) SetWarrantyPeriod(months int) {
	now := time.Now()
	dr.WarrantyStartDate = &now
	end := now.AddDate(0, months, 0)
	dr.WarrantyEndDate = &end
	dr.RefurbWarrantyMonths = months
}

// IsUnderWarranty checks if refurbishment is still under warranty
func (dr *DeviceRefurbishment) IsUnderWarranty() bool {
	if dr.WarrantyEndDate == nil {
		return false
	}
	return time.Now().Before(*dr.WarrantyEndDate)
}

// Complete marks refurbishment as completed
func (dr *DeviceRefurbishment) Complete(grade string) {
	now := time.Now()
	dr.RefurbishmentDate = &now
	dr.RefurbishmentStatus = "completed"
	dr.PostRefurbGrade = grade
}
