package smartwatch

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartwatchAccessory represents smartwatch-specific accessories
type SmartwatchAccessory struct {
	database.BaseModel

	DeviceID             uuid.UUID  `gorm:"type:uuid;not null" json:"device_id"`
	AccessoryType        string     `gorm:"type:varchar(50);not null" json:"accessory_type"`
	Brand                string     `json:"brand"`
	Model                string     `json:"model"`
	SerialNumber         string     `json:"serial_number"`
	PurchaseDate         *time.Time `json:"purchase_date"`
	PurchasePrice        float64    `json:"purchase_price"`
	CurrentValue         float64    `json:"current_value"`
	Condition            string     `gorm:"type:varchar(20);default:'good'" json:"condition"`
	IsOriginal           bool       `gorm:"default:false" json:"is_original"`
	IsCoveredByInsurance bool       `gorm:"default:false" json:"is_covered_by_insurance"`
	IsIncludedInPolicy   bool       `gorm:"default:false" json:"is_included_in_policy"`
	WarrantyExpiry       *time.Time `json:"warranty_expiry"`

	// Watch Band/Strap specifics
	BandType       string `json:"band_type"`     // sport, leather, metal, milanese, link, solo_loop
	BandMaterial   string `json:"band_material"` // silicone, leather, stainless_steel, titanium, nylon
	BandSize       string `json:"band_size"`     // S, M, L, or mm measurements (38mm, 40mm, 42mm, 44mm, 45mm, 49mm)
	BandColor      string `json:"band_color"`
	BandWidth      int    `json:"band_width"` // mm
	QuickRelease   bool   `gorm:"default:false" json:"quick_release"`
	Adjustable     bool   `gorm:"default:true" json:"adjustable"`
	WaterResistant bool   `gorm:"default:false" json:"water_resistant"`
	SweatResistant bool   `gorm:"default:false" json:"sweat_resistant"`
	Hypoallergenic bool   `gorm:"default:false" json:"hypoallergenic"`
	MagneticClasp  bool   `gorm:"default:false" json:"magnetic_clasp"`

	// Watch Case/Protector specifics
	CaseType              string `json:"case_type"`             // bumper, full_cover, screen_only
	CaseMaterial          string `json:"case_material"`         // TPU, PC, aluminum, carbon_fiber
	ScreenProtectorType   string `json:"screen_protector_type"` // tempered_glass, film, liquid
	ScreenProtectorCurved bool   `gorm:"default:false" json:"screen_protector_curved"`
	ImpactResistance      string `json:"impact_resistance"`  // basic, military_grade
	ScratchResistance     int    `json:"scratch_resistance"` // Mohs scale rating

	// Charging Accessory specifics
	ChargerType         string  `json:"charger_type"`     // magnetic, dock, stand, portable, cable
	ChargingSpeed       string  `json:"charging_speed"`   // standard, fast
	ChargingWattage     float64 `json:"charging_wattage"` // watts
	WirelessCharging    bool    `gorm:"default:false" json:"wireless_charging"`
	MagneticAlignment   bool    `gorm:"default:false" json:"magnetic_alignment"`
	ChargingPorts       int     `json:"charging_ports"`        // for multi-device chargers
	PortableCapacity    int     `json:"portable_capacity"`     // mAh for portable chargers
	ChargingCableLength float64 `json:"charging_cable_length"` // meters

	// Stand/Dock specifics
	StandType          string `json:"stand_type"`     // desktop, nightstand, travel
	StandMaterial      string `json:"stand_material"` // aluminum, wood, plastic
	AdjustableAngle    bool   `gorm:"default:false" json:"adjustable_angle"`
	ChargingIntegrated bool   `gorm:"default:false" json:"charging_integrated"`
	PhoneStandIncluded bool   `gorm:"default:false" json:"phone_stand_included"`
	CableManagement    bool   `gorm:"default:false" json:"cable_management"`

	// Fitness Accessory specifics
	HeartRateStrap    bool   `gorm:"default:false" json:"heart_rate_strap"`
	HeartRateAccuracy string `json:"heart_rate_accuracy"` // medical_grade, fitness_grade
	ArmBand           bool   `gorm:"default:false" json:"arm_band"`
	SportClip         bool   `gorm:"default:false" json:"sport_clip"`

	// Storage/Travel specifics
	CaseCapacity     int    `json:"case_capacity"` // number of watches
	TravelPouch      bool   `gorm:"default:false" json:"travel_pouch"`
	HardShell        bool   `gorm:"default:false" json:"hard_shell"`
	WaterproofRating string `json:"waterproof_rating"` // IPX rating

	// Compatibility
	CompatibleModels string `gorm:"type:json" json:"compatible_models"` // JSON array of watch models
	CompatibleSizes  string `gorm:"type:json" json:"compatible_sizes"`  // JSON array of sizes
	UniversalFit     bool   `gorm:"default:false" json:"universal_fit"`

	// Special Features
	SpecialEdition bool   `gorm:"default:false" json:"special_edition"`
	LimitedEdition bool   `gorm:"default:false" json:"limited_edition"`
	DesignerBrand  string `json:"designer_brand"`                      // Hermès, Nike, etc.
	SmartFeatures  bool   `gorm:"default:false" json:"smart_features"` // LED, sensors in band

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
}

// Business logic methods

// IsWatchBand checks if this is a watch band/strap
func (sa *SmartwatchAccessory) IsWatchBand() bool {
	bandTypes := []string{"band", "strap", "bracelet", "loop", "link"}
	for _, t := range bandTypes {
		if sa.AccessoryType == t {
			return true
		}
	}
	return sa.BandType != ""
}

// IsProtectiveAccessory checks if this is protective
func (sa *SmartwatchAccessory) IsProtectiveAccessory() bool {
	protectiveTypes := []string{"case", "bumper", "screen_protector", "cover"}
	for _, t := range protectiveTypes {
		if sa.AccessoryType == t {
			return true
		}
	}
	return sa.CaseType != "" || sa.ScreenProtectorType != ""
}

// IsChargingAccessory checks if this is a charging accessory
func (sa *SmartwatchAccessory) IsChargingAccessory() bool {
	chargingTypes := []string{"charger", "charging_dock", "charging_stand", "charging_cable", "power_bank"}
	for _, t := range chargingTypes {
		if sa.AccessoryType == t {
			return true
		}
	}
	return sa.ChargerType != ""
}

// IsPremiumBand checks if band is premium/luxury
func (sa *SmartwatchAccessory) IsPremiumBand() bool {
	// Check material
	premiumMaterials := []string{"leather", "stainless_steel", "titanium", "ceramic", "link"}
	for _, m := range premiumMaterials {
		if sa.BandMaterial == m {
			return true
		}
	}
	// Check brand
	if sa.DesignerBrand != "" {
		return true
	}
	// Check price
	if sa.PurchasePrice > 100 {
		return true
	}
	// Check edition
	return sa.SpecialEdition || sa.LimitedEdition
}

// IsSportBand checks if band is sport/fitness oriented
func (sa *SmartwatchAccessory) IsSportBand() bool {
	sportTypes := []string{"sport", "sport_loop", "sport_band", "nike"}
	for _, t := range sportTypes {
		if sa.BandType == t {
			return true
		}
	}
	sportMaterials := []string{"silicone", "fluoroelastomer", "nylon"}
	for _, m := range sportMaterials {
		if sa.BandMaterial == m && (sa.SweatResistant || sa.WaterResistant) {
			return true
		}
	}
	return false
}

// IsFitnessAccessory checks if this is fitness-related
func (sa *SmartwatchAccessory) IsFitnessAccessory() bool {
	return sa.HeartRateStrap || sa.ArmBand || sa.SportClip
}

// NeedsReplacement checks if accessory needs replacement
func (sa *SmartwatchAccessory) NeedsReplacement() bool {
	if sa.IsDamaged && sa.Condition == "broken" {
		return true
	}
	if sa.IsLost {
		return true
	}

	// Check age for bands (they wear out)
	if sa.IsWatchBand() && sa.PurchaseDate != nil {
		age := time.Since(*sa.PurchaseDate)
		// Sport bands typically last 6-12 months with daily use
		if sa.IsSportBand() && age > 12*30*24*time.Hour {
			return true
		}
		// Leather bands last 1-2 years
		if sa.BandMaterial == "leather" && age > 24*30*24*time.Hour {
			return true
		}
	}

	return false
}

// CalculateDepreciation calculates current value
func (sa *SmartwatchAccessory) CalculateDepreciation() float64 {
	if sa.PurchaseDate == nil || sa.PurchasePrice == 0 {
		return 0
	}

	monthsOld := int(time.Since(*sa.PurchaseDate).Hours() / 24 / 30)

	// Different depreciation rates
	var depreciationRate float64
	if sa.IsPremiumBand() {
		depreciationRate = 0.02 // 2% per month for premium
	} else if sa.IsWatchBand() {
		depreciationRate = 0.08 // 8% per month for regular bands
	} else if sa.IsChargingAccessory() {
		depreciationRate = 0.03 // 3% per month for chargers
	} else if sa.IsProtectiveAccessory() {
		depreciationRate = 0.05 // 5% per month for cases
	} else {
		depreciationRate = 0.04 // 4% per month for others
	}

	depreciation := float64(monthsOld) * depreciationRate * sa.PurchasePrice
	currentValue := sa.PurchasePrice - depreciation

	// Minimum value is 15% of purchase price for premium, 10% for others
	minValue := sa.PurchasePrice * 0.10
	if sa.IsPremiumBand() {
		minValue = sa.PurchasePrice * 0.15
	}

	if currentValue < minValue {
		currentValue = minValue
	}

	return currentValue
}

// IsCompatibleWithSize checks band size compatibility
func (sa *SmartwatchAccessory) IsCompatibleWithSize(size string) bool {
	// Universal fit bands work with all sizes
	if sa.UniversalFit {
		return true
	}
	// Check if size is in compatible sizes
	// This would parse JSON array, simplified here
	return sa.BandSize == size || sa.Adjustable
}

// GetWaterResistanceLevel returns water resistance level
func (sa *SmartwatchAccessory) GetWaterResistanceLevel() string {
	if sa.WaterproofRating != "" {
		return sa.WaterproofRating
	}
	if sa.WaterResistant && sa.SweatResistant {
		return "high"
	}
	if sa.WaterResistant || sa.SweatResistant {
		return "medium"
	}
	return "none"
}

// IsFastCharger checks if charger supports fast charging
func (sa *SmartwatchAccessory) IsFastCharger() bool {
	return sa.ChargingSpeed == "fast" || sa.ChargingWattage >= 5
}

// IsPortableCharger checks if charger is portable
func (sa *SmartwatchAccessory) IsPortableCharger() bool {
	return sa.ChargerType == "portable" || sa.PortableCapacity > 0
}

// GetQualityScore rates accessory quality
func (sa *SmartwatchAccessory) GetQualityScore() float64 {
	score := 50.0 // Base score

	// Original/authentic gets bonus
	if sa.IsOriginal {
		score += 20
	}

	// Material quality
	if sa.IsPremiumBand() {
		score += 15
	}

	// Features
	if sa.QuickRelease {
		score += 5
	}
	if sa.WaterResistant {
		score += 5
	}
	if sa.Hypoallergenic {
		score += 5
	}

	// Condition affects score
	switch sa.Condition {
	case "new":
		score += 10
	case "excellent":
		score += 5
	case "good":
		// no change
	case "fair":
		score -= 10
	case "poor":
		score -= 20
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// IsEligibleForClaim checks claim eligibility
func (sa *SmartwatchAccessory) IsEligibleForClaim() bool {
	if !sa.IsCoveredByInsurance || !sa.IsIncludedInPolicy {
		return false
	}
	if sa.IsReplaced {
		return false
	}
	// Minimum claim value
	if sa.CurrentValue < 20 {
		return false
	}
	// Premium bands always eligible if covered
	if sa.IsPremiumBand() {
		return true
	}
	return true
}

// GetWarrantyStatus returns warranty status
func (sa *SmartwatchAccessory) GetWarrantyStatus() string {
	if sa.WarrantyExpiry == nil {
		return "no_warranty"
	}
	if time.Now().After(*sa.WarrantyExpiry) {
		return "expired"
	}
	daysRemaining := int(time.Until(*sa.WarrantyExpiry).Hours() / 24)
	if daysRemaining <= 30 {
		return "expiring_soon"
	}
	return "active"
}
