package smartphone

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartphoneAccessory represents smartphone-specific accessories
type SmartphoneAccessory struct {
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

	// Smartphone-specific accessory details
	CaseType             string `json:"case_type"`              // wallet, folio, bumper, armor, clear
	ProtectionLevel      string `json:"protection_level"`       // minimal, standard, military_grade
	ScreenProtectorType  string `json:"screen_protector_type"`  // tempered_glass, film, liquid, privacy
	ChargerWattage       int    `json:"charger_wattage"`        // 5W, 18W, 20W, 45W, 65W, 100W+
	ChargingProtocol     string `json:"charging_protocol"`      // USB-PD, QC3.0, QC4.0, SuperVOOC, Warp
	CableType            string `json:"cable_type"`             // USB-C, Lightning, Micro-USB
	CableLength          string `json:"cable_length"`           // 0.5m, 1m, 2m, 3m
	WirelessChargingType string `json:"wireless_charging_type"` // Qi, MagSafe, proprietary

	// Mount & Holder specifics
	MountType        string `json:"mount_type"`       // magnetic, suction, clamp, adhesive
	MountLocation    string `json:"mount_location"`   // dashboard, windshield, air_vent, cd_slot
	RotationDegrees  int    `json:"rotation_degrees"` // 360, 180, fixed
	AdjustableAngles bool   `gorm:"default:false" json:"adjustable_angles"`

	// Audio accessory specifics
	HeadphoneType     string `json:"headphone_type"`     // wired, bluetooth, true_wireless
	AudioCodec        string `json:"audio_codec"`        // SBC, AAC, aptX, LDAC, LHDC
	NoiseCancellation string `json:"noise_cancellation"` // none, passive, ANC
	BluetoothVersion  string `json:"bluetooth_version"`  // 4.2, 5.0, 5.1, 5.2, 5.3

	// Camera accessory specifics
	LensType          string `json:"lens_type"`          // wide, macro, telephoto, fisheye
	LensMagnification string `json:"lens_magnification"` // 2x, 10x, 20x, etc
	TripodHeight      string `json:"tripod_height"`      // mini, tabletop, full_size
	GimbalAxes        int    `json:"gimbal_axes"`        // 1, 2, 3 axis

	// Gaming accessory specifics
	ControllerType string `json:"controller_type"` // clip_on, bluetooth, cooling_fan
	CoolingType    string `json:"cooling_type"`    // passive, active_fan, peltier
	TriggerButtons int    `json:"trigger_buttons"` // 2, 4, 6

	// Power bank specifics
	PowerBankCapacity   int  `json:"powerbank_capacity"` // mAh
	PowerBankPorts      int  `json:"powerbank_ports"`    // number of output ports
	PassThroughCharging bool `gorm:"default:false" json:"pass_through_charging"`

	// Compatibility and features
	CompatibleModels string `gorm:"type:json" json:"compatible_models"` // JSON array of compatible phone models
	SpecialFeatures  string `gorm:"type:json" json:"special_features"`  // JSON array of features
	Color            string `json:"color"`
	Material         string `json:"material"`      // plastic, silicone, leather, metal, carbon_fiber
	IPRating         string `json:"ip_rating"`     // IP68, IPX4, etc
	MILSTDRating     string `json:"milstd_rating"` // MIL-STD-810G/H

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

// IsProtectiveAccessory checks if this is a protective accessory
func (sa *SmartphoneAccessory) IsProtectiveAccessory() bool {
	protectiveTypes := []string{"case", "screen_protector", "camera_lens_protector", "bumper", "armor_case"}
	for _, t := range protectiveTypes {
		if sa.AccessoryType == t {
			return true
		}
	}
	return false
}

// IsFastCharger checks if charger supports fast charging
func (sa *SmartphoneAccessory) IsFastCharger() bool {
	return sa.ChargerWattage >= 18
}

// IsSuperFastCharger checks if charger supports super fast charging
func (sa *SmartphoneAccessory) IsSuperFastCharger() bool {
	return sa.ChargerWattage >= 45
}

// IsWirelessCharger checks if this is a wireless charging accessory
func (sa *SmartphoneAccessory) IsWirelessCharger() bool {
	return sa.WirelessChargingType != "" && sa.AccessoryType == "wireless_charger"
}

// IsMagSafeCompatible checks if accessory is MagSafe compatible
func (sa *SmartphoneAccessory) IsMagSafeCompatible() bool {
	return sa.WirelessChargingType == "MagSafe" || sa.MountType == "magnetic"
}

// IsGamingAccessory checks if this is a gaming-related accessory
func (sa *SmartphoneAccessory) IsGamingAccessory() bool {
	gamingTypes := []string{"game_controller", "cooling_fan", "trigger_button", "gaming_grip"}
	for _, t := range gamingTypes {
		if sa.AccessoryType == t {
			return true
		}
	}
	return sa.ControllerType != "" || sa.CoolingType != "" || sa.TriggerButtons > 0
}

// IsAudioAccessory checks if this is an audio accessory
func (sa *SmartphoneAccessory) IsAudioAccessory() bool {
	audioTypes := []string{"headphones", "earbuds", "bluetooth_speaker", "audio_adapter"}
	for _, t := range audioTypes {
		if sa.AccessoryType == t {
			return true
		}
	}
	return sa.HeadphoneType != "" || sa.AudioCodec != ""
}

// IsCameraAccessory checks if this is a camera-related accessory
func (sa *SmartphoneAccessory) IsCameraAccessory() bool {
	cameraTypes := []string{"lens", "tripod", "gimbal", "ring_light", "selfie_stick"}
	for _, t := range cameraTypes {
		if sa.AccessoryType == t {
			return true
		}
	}
	return sa.LensType != "" || sa.GimbalAxes > 0
}

// IsPremiumAccessory checks if this is a premium/high-value accessory
func (sa *SmartphoneAccessory) IsPremiumAccessory() bool {
	// Premium based on price
	if sa.PurchasePrice > 100 {
		return true
	}
	// Premium based on features
	if sa.GimbalAxes == 3 || sa.ChargerWattage >= 65 || sa.AudioCodec == "LDAC" {
		return true
	}
	// Premium materials
	premiumMaterials := []string{"leather", "carbon_fiber", "titanium", "ceramic"}
	for _, m := range premiumMaterials {
		if sa.Material == m {
			return true
		}
	}
	return false
}

// HasMilitaryGradeProtection checks for military-grade protection
func (sa *SmartphoneAccessory) HasMilitaryGradeProtection() bool {
	return sa.MILSTDRating != "" || sa.ProtectionLevel == "military_grade"
}

// IsWaterproof checks if accessory is waterproof
func (sa *SmartphoneAccessory) IsWaterproof() bool {
	// Check IP rating for water resistance (IPX4 and above)
	if len(sa.IPRating) >= 4 {
		if sa.IPRating[0:2] == "IP" {
			waterRating := sa.IPRating[3:4]
			return waterRating >= "4"
		}
	}
	return false
}

// CalculateDepreciation calculates accessory depreciation
func (sa *SmartphoneAccessory) CalculateDepreciation() float64 {
	if sa.PurchaseDate == nil || sa.PurchasePrice == 0 {
		return 0
	}

	monthsOld := int(time.Since(*sa.PurchaseDate).Hours() / 24 / 30)

	// Different depreciation rates for different accessory types
	var depreciationRate float64
	if sa.IsProtectiveAccessory() {
		depreciationRate = 0.05 // 5% per month for cases
	} else if sa.IsFastCharger() || sa.IsWirelessCharger() {
		depreciationRate = 0.03 // 3% per month for chargers
	} else if sa.IsAudioAccessory() {
		depreciationRate = 0.04 // 4% per month for audio
	} else if sa.IsCameraAccessory() {
		depreciationRate = 0.035 // 3.5% per month for camera accessories
	} else {
		depreciationRate = 0.045 // 4.5% per month for others
	}

	depreciation := float64(monthsOld) * depreciationRate * sa.PurchasePrice
	currentValue := sa.PurchasePrice - depreciation

	// Minimum value is 10% of purchase price
	minValue := sa.PurchasePrice * 0.1
	if currentValue < minValue {
		currentValue = minValue
	}

	return currentValue
}

// IsEligibleForClaim checks if accessory can be claimed
func (sa *SmartphoneAccessory) IsEligibleForClaim() bool {
	if !sa.IsCoveredByInsurance || !sa.IsIncludedInPolicy {
		return false
	}
	if sa.IsReplaced {
		return false
	}
	if sa.CurrentValue < 10 { // Minimum claim value
		return false
	}
	return true
}

// GetReplacementCost calculates replacement cost
func (sa *SmartphoneAccessory) GetReplacementCost() float64 {
	if sa.ReplacementCost > 0 {
		return sa.ReplacementCost
	}

	// Estimate based on original price and type
	baseCost := sa.PurchasePrice

	// Adjust for premium accessories
	if sa.IsPremiumAccessory() {
		baseCost *= 1.1 // 10% premium
	}

	// Adjust for discontinued/rare items
	if sa.PurchaseDate != nil && time.Since(*sa.PurchaseDate) > 2*365*24*time.Hour {
		baseCost *= 1.2 // 20% markup for older items
	}

	return baseCost
}

// NeedsReplacement checks if accessory needs replacement
func (sa *SmartphoneAccessory) NeedsReplacement() bool {
	if sa.IsDamaged && sa.Condition == "broken" {
		return true
	}
	if sa.IsLost {
		return true
	}
	// Check age for certain accessories
	if sa.PurchaseDate != nil {
		age := time.Since(*sa.PurchaseDate)
		if sa.IsProtectiveAccessory() && age > 12*30*24*time.Hour { // 1 year for cases
			return true
		}
		if sa.AccessoryType == "battery_case" && age > 18*30*24*time.Hour { // 1.5 years for battery cases
			return true
		}
	}
	return false
}

// GetWarrantyStatus returns warranty status
func (sa *SmartphoneAccessory) GetWarrantyStatus() string {
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
