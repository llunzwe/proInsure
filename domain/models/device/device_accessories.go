package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceAccessory represents generic accessories for any device type
// Category-specific accessory details are stored in:
// - smartphone.SmartphoneAccessory for phones
// - smartwatch.SmartwatchAccessory for watches
type DeviceAccessory struct {
	database.BaseModel
	DeviceID             uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	DeviceCategory       string     `gorm:"type:varchar(30);not null" json:"device_category"` // smartphone, smartwatch, tablet, laptop
	AccessoryType        string     `gorm:"type:varchar(50);not null" json:"accessory_type"`  // case, charger, band, etc
	Brand                string     `json:"brand"`
	Model                string     `json:"model"`
	SerialNumber         string     `json:"serial_number"`
	PurchaseDate         *time.Time `json:"purchase_date"`
	PurchasePrice        float64    `json:"purchase_price"`
	CurrentValue         float64    `json:"current_value"`
	Condition            string     `gorm:"type:varchar(20);default:'good'" json:"condition"`
	IsOriginal           bool       `gorm:"default:false" json:"is_original"` // Original manufacturer accessory
	IsCoveredByInsurance bool       `gorm:"default:false" json:"is_covered_by_insurance"`
	IsIncludedInPolicy   bool       `gorm:"default:false" json:"is_included_in_policy"`
	WarrantyExpiry       *time.Time `json:"warranty_expiry"`

	// Common accessory details
	Compatibility string `json:"compatibility"` // Device models compatible with
	Color         string `json:"color"`
	Material      string `json:"material"` // leather, silicone, plastic, metal, etc

	// Category-specific reference
	CategorySpecificID   *uuid.UUID `gorm:"type:uuid" json:"category_specific_id"` // ID in category-specific table
	CategorySpecificType string     `json:"category_specific_type"`                // Table name: smartphone_accessories, smartwatch_accessories

	// Common electronic features (applicable to multiple categories)
	PowerOutput        string `json:"power_output"` // For chargers: 5W, 20W, etc
	WirelessCapability bool   `gorm:"default:false" json:"wireless_capability"`
	BatteryCapacity    int    `json:"battery_capacity"` // For power banks, battery cases

	// Tracking & verification
	PhotoURLs          string     `gorm:"type:json" json:"photo_urls"`
	ReceiptURL         string     `json:"receipt_url"`
	IsVerified         bool       `gorm:"default:false" json:"is_verified"`
	VerificationDate   *time.Time `json:"verification_date"`
	LastInspectionDate *time.Time `json:"last_inspection_date"`

	// Loss/damage tracking
	IsLost            bool       `gorm:"default:false" json:"is_lost"`
	LostDate          *time.Time `json:"lost_date"`
	IsDamaged         bool       `gorm:"default:false" json:"is_damaged"`
	DamageDescription string     `json:"damage_description"`
	DamageDate        *time.Time `json:"damage_date"`

	// Replacement info
	ReplacementCost        float64    `json:"replacement_cost"`
	IsReplaced             bool       `gorm:"default:false" json:"is_replaced"`
	ReplacementDate        *time.Time `json:"replacement_date"`
	ReplacementAccessoryID *uuid.UUID `gorm:"type:uuid" json:"replacement_accessory_id"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// Category-specific accessory should be loaded via service layer using CategorySpecificID
}

// TableName returns the table name
func (t *DeviceAccessory) TableName() string {
	return "device_accessories"
}

// BeforeCreate handles pre-creation logic
func (da *DeviceAccessory) BeforeCreate(tx *gorm.DB) error {
	if err := da.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Calculate current value if not set
	if da.CurrentValue == 0 && da.PurchasePrice > 0 {
		da.CurrentValue = da.CalculateDepreciation()
	}

	return nil
}

// CalculateDepreciation calculates current value based on age
func (da *DeviceAccessory) CalculateDepreciation() float64 {
	if da.PurchasePrice == 0 || da.PurchaseDate == nil {
		return da.CurrentValue
	}

	ageYears := time.Since(*da.PurchaseDate).Hours() / (24 * 365)

	// Use the category and type-specific depreciation rate
	depreciationRate := da.GetDepreciationRate()

	depreciated := da.PurchasePrice * (1 - (depreciationRate * ageYears))

	// Minimum value based on category
	minValue := da.PurchasePrice * 0.1 // Default 10%

	// Premium accessories retain more value
	if da.IsOriginal && (da.DeviceCategory == "smartwatch" || da.DeviceCategory == "smartphone") {
		minValue = da.PurchasePrice * 0.15 // 15% for original accessories
	}

	if depreciated < minValue {
		depreciated = minValue
	}

	return depreciated
}

// IsEligibleForClaim checks if accessory can be claimed
func (da *DeviceAccessory) IsEligibleForClaim() bool {
	// Must be covered by insurance
	if !da.IsCoveredByInsurance || !da.IsIncludedInPolicy {
		return false
	}

	// Must be verified
	if !da.IsVerified {
		return false
	}

	// Value threshold (minimum $20)
	if da.CurrentValue < 20 {
		return false
	}

	// Not already replaced
	if da.IsReplaced {
		return false
	}

	return true
}

// CalculateReplacementCost calculates cost to replace accessory
func (da *DeviceAccessory) CalculateReplacementCost() float64 {
	if da.ReplacementCost > 0 {
		return da.ReplacementCost
	}

	// For original accessories, use purchase price
	if da.IsOriginal {
		return da.PurchasePrice
	}

	// For third-party, use current market value
	replacementCost := da.CurrentValue * 1.1 // Add 10% buffer

	// Apply minimum based on type
	minimums := map[string]float64{
		"charger":          15,
		"cable":            10,
		"case":             20,
		"screen_protector": 10,
		"headphones":       30,
		"power_bank":       25,
		"wireless_charger": 35,
	}

	if min, exists := minimums[da.AccessoryType]; exists {
		if replacementCost < min {
			replacementCost = min
		}
	}

	return replacementCost
}

// IsUnderWarranty checks if accessory is under warranty
func (da *DeviceAccessory) IsUnderWarranty() bool {
	if da.WarrantyExpiry == nil {
		return false
	}
	return time.Now().Before(*da.WarrantyExpiry)
}

// MarkAsLost marks accessory as lost
func (da *DeviceAccessory) MarkAsLost() {
	da.IsLost = true
	now := time.Now()
	da.LostDate = &now
}

// MarkAsDamaged marks accessory as damaged
func (da *DeviceAccessory) MarkAsDamaged(description string) {
	da.IsDamaged = true
	da.DamageDescription = description
	now := time.Now()
	da.DamageDate = &now
}

// MarkAsReplaced marks accessory as replaced
func (da *DeviceAccessory) MarkAsReplaced(replacementID *uuid.UUID) {
	da.IsReplaced = true
	now := time.Now()
	da.ReplacementDate = &now
	da.ReplacementAccessoryID = replacementID
}

// GetConditionScore returns numeric condition score
func (da *DeviceAccessory) GetConditionScore() float64 {
	scores := map[string]float64{
		"excellent": 1.0,
		"good":      0.85,
		"fair":      0.7,
		"poor":      0.5,
	}

	if score, exists := scores[da.Condition]; exists {
		return score
	}

	return 0.7 // Default to fair
}

// RequiresInspection checks if accessory needs inspection
func (da *DeviceAccessory) RequiresInspection() bool {
	if da.LastInspectionDate == nil {
		return true
	}

	// Require inspection every 12 months for covered accessories
	if da.IsCoveredByInsurance {
		daysSinceInspection := time.Since(*da.LastInspectionDate).Hours() / 24
		return daysSinceInspection > 365
	}

	return false
}

// GetCategorySpecificTable returns the table name for category-specific details
func (da *DeviceAccessory) GetCategorySpecificTable() string {
	switch da.DeviceCategory {
	case "smartphone":
		return "smartphone_accessories"
	case "smartwatch":
		return "smartwatch_accessories"
	case "tablet":
		return "tablet_accessories"
	case "laptop":
		return "laptop_accessories"
	default:
		return ""
	}
}

// IsElectronicAccessory checks if this is an electronic accessory
func (da *DeviceAccessory) IsElectronicAccessory() bool {
	// Check for electronic features
	if da.PowerOutput != "" || da.BatteryCapacity > 0 || da.WirelessCapability {
		return true
	}

	// Check common electronic accessory types
	electronicTypes := []string{
		"charger", "wireless_charger", "power_bank", "cable",
		"headphones", "earbuds", "speaker", "adapter",
	}

	for _, t := range electronicTypes {
		if da.AccessoryType == t {
			return true
		}
	}

	return false
}

// IsProtectiveAccessory checks if this is a protective accessory
func (da *DeviceAccessory) IsProtectiveAccessory() bool {
	protectiveTypes := []string{
		"case", "cover", "screen_protector", "bumper",
		"skin", "armor", "pouch", "sleeve",
	}

	for _, t := range protectiveTypes {
		if da.AccessoryType == t {
			return true
		}
	}

	return false
}

// GetDepreciationRate returns depreciation rate based on type and category
func (da *DeviceAccessory) GetDepreciationRate() float64 {
	baseRate := 0.3 // 30% per year default

	// Adjust for original accessories
	if da.IsOriginal {
		baseRate = 0.25
	}

	// Electronic accessories depreciate faster
	if da.IsElectronicAccessory() {
		baseRate = 0.35
	}

	// Protective accessories for premium devices hold value better
	if da.IsProtectiveAccessory() && da.IsOriginal {
		baseRate = 0.20
	}

	// Category-specific adjustments
	switch da.DeviceCategory {
	case "smartwatch":
		// Watch bands can hold value if premium
		if da.Material == "leather" || da.Material == "titanium" || da.Material == "ceramic" {
			baseRate = 0.15
		}
	case "smartphone":
		// Phone accessories generally depreciate at standard rate
		// No adjustment needed
	}

	return baseRate
}
