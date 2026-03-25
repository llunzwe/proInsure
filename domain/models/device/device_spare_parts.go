package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceSparePart represents spare parts inventory for device repairs
type DeviceSparePart struct {
	database.BaseModel
	DeviceID           uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	PartType           string    `gorm:"type:varchar(50);not null" json:"part_type"` // screen, battery, camera, etc
	PartNumber         string    `gorm:"uniqueIndex" json:"part_number"`
	PartName           string    `json:"part_name"`
	Manufacturer       string    `json:"manufacturer"`
	SupplierName       string    `json:"supplier_name"`
	SupplierPartNumber string    `json:"supplier_part_number"`

	// Compatibility
	CompatibleModels   string `gorm:"type:json" json:"compatible_models"` // JSON array of device models
	CompatibilityNotes string `json:"compatibility_notes"`
	IsUniversal        bool   `gorm:"default:false" json:"is_universal"`

	// Quality & authenticity
	Quality           string `gorm:"type:varchar(20);default:'oem'" json:"quality"` // oem, genuine, aftermarket, refurbished
	IsOriginalPart    bool   `gorm:"default:false" json:"is_original_part"`
	IsCertified       bool   `gorm:"default:false" json:"is_certified"`
	CertificationBody string `json:"certification_body"`
	QualityGrade      string `gorm:"type:varchar(5)" json:"quality_grade"` // A, B, C

	// Inventory & pricing
	QuantityInStock int     `json:"quantity_in_stock"`
	MinimumStock    int     `json:"minimum_stock"`
	ReorderPoint    int     `json:"reorder_point"`
	UnitCost        float64 `json:"unit_cost"`
	RetailPrice     float64 `json:"retail_price"`
	WholesalePrice  float64 `json:"wholesale_price"`
	Currency        string  `gorm:"default:'USD'" json:"currency"`

	// Warranty & lifecycle
	WarrantyDays     int        `json:"warranty_days"`
	ExpectedLifespan int        `json:"expected_lifespan_days"`
	ManufactureDate  *time.Time `json:"manufacture_date"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	BatchNumber      string     `json:"batch_number"`

	// Location & storage
	WarehouseLocation       string `json:"warehouse_location"`
	ShelfNumber             string `json:"shelf_number"`
	StorageConditions       string `json:"storage_conditions"` // temperature, humidity requirements
	IsFragile               bool   `gorm:"default:false" json:"is_fragile"`
	RequiresSpecialHandling bool   `gorm:"default:false" json:"requires_special_handling"`

	// Usage tracking
	TimesUsed           int        `json:"times_used"`
	LastUsedDate        *time.Time `json:"last_used_date"`
	LastUsedForRepairID *uuid.UUID `gorm:"type:uuid" json:"last_used_for_repair_id"`
	AverageRepairTime   float64    `json:"average_repair_time_hours"`
	SuccessRate         float64    `json:"success_rate"` // Percentage of successful repairs

	// Technical specifications
	Specifications string  `gorm:"type:json" json:"specifications"` // JSON object of technical specs
	Weight         float64 `json:"weight_grams"`
	Dimensions     string  `json:"dimensions"` // LxWxH in mm
	Color          string  `json:"color"`
	Material       string  `json:"material"`

	// For electronic parts
	Voltage       string `json:"voltage"`
	Capacity      string `json:"capacity"` // For batteries: mAh
	PowerRating   string `json:"power_rating"`
	ConnectorType string `json:"connector_type"`

	// Documentation
	DatasheetURL         string `json:"datasheet_url"`
	InstallationGuideURL string `json:"installation_guide_url"`
	PhotoURLs            string `gorm:"type:json" json:"photo_urls"`
	VideoGuideURL        string `json:"video_guide_url"`

	// Status & availability
	Status       string     `gorm:"type:varchar(20);default:'available'" json:"status"` // available, out_of_stock, discontinued, recalled
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	IsRecalled   bool       `gorm:"default:false" json:"is_recalled"`
	RecallDate   *time.Time `json:"recall_date"`
	RecallReason string     `json:"recall_reason"`

	// Supplier & ordering
	LeadTimeDays     int        `json:"lead_time_days"`
	LastOrderDate    *time.Time `json:"last_order_date"`
	LastReceivedDate *time.Time `json:"last_received_date"`
	SupplierRating   float64    `json:"supplier_rating"`

	// Insurance & claims
	InsuranceApproved   bool    `gorm:"default:false" json:"insurance_approved"`
	MaxClaimAmount      float64 `json:"max_claim_amount"`
	RequiresPreApproval bool    `gorm:"default:false" json:"requires_pre_approval"`
	ApprovalNotes       string  `json:"approval_notes"`

	// Performance metrics
	DefectRate           float64 `json:"defect_rate"`                                // Percentage
	ReturnRate           float64 `json:"return_rate"`                                // Percentage
	CustomerSatisfaction float64 `json:"customer_satisfaction"`                      // 0-5 rating
	InstallDifficulty    string  `gorm:"type:varchar(20)" json:"install_difficulty"` // easy, medium, hard, expert

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	UsedInRepairs []DeviceRepair `gorm:"many2many:repair_spare_parts;" json:"used_in_repairs,omitempty"`
}

// TableName returns the table name
func (t *DeviceSparePart) TableName() string {
	return "device_spare_parts"
}

// BeforeCreate handles pre-creation logic
func (dsp *DeviceSparePart) BeforeCreate(tx *gorm.DB) error {
	if err := dsp.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Set default reorder point if not specified
	if dsp.ReorderPoint == 0 && dsp.MinimumStock > 0 {
		dsp.ReorderPoint = dsp.MinimumStock * 2
	}

	return nil
}

// IsInStock checks if part is available in stock
func (dsp *DeviceSparePart) IsInStock() bool {
	return dsp.QuantityInStock > 0 && dsp.Status == "available"
}

// NeedsReorder checks if part needs to be reordered
func (dsp *DeviceSparePart) NeedsReorder() bool {
	return dsp.QuantityInStock <= dsp.ReorderPoint
}

// IsCompatibleWith checks if part is compatible with a specific device model
func (dsp *DeviceSparePart) IsCompatibleWith(deviceModel string) bool {
	// If universal part, compatible with all
	if dsp.IsUniversal {
		return true
	}

	// Check compatible models list
	// In production, parse JSON and check
	// For now, simplified check
	return true
}

// CalculateRepairCost calculates total cost including labor
func (dsp *DeviceSparePart) CalculateRepairCost(laborHours float64, laborRate float64) float64 {
	partCost := dsp.RetailPrice
	laborCost := laborHours * laborRate

	// Add markup based on difficulty
	markup := 1.0
	switch dsp.InstallDifficulty {
	case "hard":
		markup = 1.2
	case "expert":
		markup = 1.3
	}

	return (partCost + laborCost) * markup
}

// IsHighQuality checks if part meets high quality standards
func (dsp *DeviceSparePart) IsHighQuality() bool {
	return (dsp.Quality == "oem" || dsp.Quality == "genuine") &&
		(dsp.QualityGrade == "A" || dsp.QualityGrade == "B")
}

// GetWarrantyEndDate calculates warranty end date for a specific repair
func (dsp *DeviceSparePart) GetWarrantyEndDate(repairDate time.Time) time.Time {
	return repairDate.AddDate(0, 0, dsp.WarrantyDays)
}

// IsExpired checks if part has expired
func (dsp *DeviceSparePart) IsExpired() bool {
	if dsp.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*dsp.ExpiryDate)
}

// UpdateStock updates the stock quantity
func (dsp *DeviceSparePart) UpdateStock(quantity int, operation string) {
	switch operation {
	case "add":
		dsp.QuantityInStock += quantity
	case "remove":
		dsp.QuantityInStock -= quantity
		if dsp.QuantityInStock < 0 {
			dsp.QuantityInStock = 0
		}
	case "set":
		dsp.QuantityInStock = quantity
	}

	// Update status based on stock
	if dsp.QuantityInStock == 0 {
		dsp.Status = "out_of_stock"
	} else if dsp.Status == "out_of_stock" {
		dsp.Status = "available"
	}
}

// RecordUsage records part usage in a repair
func (dsp *DeviceSparePart) RecordUsage(repairID uuid.UUID) {
	dsp.TimesUsed++
	now := time.Now()
	dsp.LastUsedDate = &now
	dsp.LastUsedForRepairID = &repairID
	dsp.UpdateStock(1, "remove")
}

// CalculateReorderQuantity calculates optimal reorder quantity
func (dsp *DeviceSparePart) CalculateReorderQuantity() int {
	// Economic Order Quantity (EOQ) simplified
	// Based on usage rate and lead time

	if dsp.TimesUsed == 0 {
		return dsp.MinimumStock * 2
	}

	// Calculate average daily usage
	daysSinceFirstUse := 365 // Default to 1 year if not tracked
	avgDailyUsage := float64(dsp.TimesUsed) / float64(daysSinceFirstUse)

	// Factor in lead time
	bufferStock := avgDailyUsage * float64(dsp.LeadTimeDays) * 1.5 // 50% safety buffer

	// Reorder quantity
	reorderQty := int(bufferStock) + dsp.MinimumStock

	// Apply maximum limits based on value
	if dsp.UnitCost > 500 {
		// Expensive parts - order less
		if reorderQty > 10 {
			reorderQty = 10
		}
	} else if dsp.UnitCost > 100 {
		if reorderQty > 50 {
			reorderQty = 50
		}
	}

	return reorderQty
}

// GetProfitMargin calculates profit margin
func (dsp *DeviceSparePart) GetProfitMargin() float64 {
	if dsp.RetailPrice == 0 {
		return 0
	}
	return ((dsp.RetailPrice - dsp.UnitCost) / dsp.RetailPrice) * 100
}

// IsApprovedForInsurance checks if part is approved for insurance claims
func (dsp *DeviceSparePart) IsApprovedForInsurance() bool {
	return dsp.InsuranceApproved &&
		dsp.IsHighQuality() &&
		!dsp.IsRecalled &&
		!dsp.IsExpired()
}

// GetReplacementPriority returns priority for replacement orders
func (dsp *DeviceSparePart) GetReplacementPriority() string {
	// Critical parts get high priority
	criticalParts := []string{"screen", "battery", "motherboard", "charging_port"}

	for _, part := range criticalParts {
		if dsp.PartType == part {
			if dsp.QuantityInStock == 0 {
				return "critical"
			}
			if dsp.NeedsReorder() {
				return "high"
			}
		}
	}

	if dsp.NeedsReorder() {
		return "medium"
	}

	return "low"
}
