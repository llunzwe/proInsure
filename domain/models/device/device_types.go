package device

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ============================================
// COMMON TYPE DEFINITIONS
// ============================================

// Money represents a monetary amount with currency
type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency" gorm:"default:'USD'"`
}

// ============================================
// DEVICE TYPE DEFINITIONS
// ============================================

// DeviceStatus represents the operational status of a device
type DeviceStatus string

// DeviceCondition represents the physical condition of a device
type DeviceCondition string

// DeviceGrade represents quality grading of a device
type DeviceGrade string

// ScreenCondition represents the condition of device screen
type ScreenCondition string

// BodyCondition represents the condition of device body/housing
type BodyCondition string

// NetworkStatus represents network lock status
type NetworkStatus string

// BlacklistStatus represents device blacklist checking status
type BlacklistStatus string

// DeviceCategory represents the type/category of device
type DeviceCategory string

const (
	DeviceCategorySmartphone DeviceCategory = "smartphone"
	DeviceCategoryTablet     DeviceCategory = "tablet"
	DeviceCategoryLaptop     DeviceCategory = "laptop"
	DeviceCategorySmartwatch DeviceCategory = "smartwatch"
	DeviceCategoryEarbuds    DeviceCategory = "earbuds"
	DeviceCategoryCamera     DeviceCategory = "camera"
	DeviceCategoryGaming     DeviceCategory = "gaming_console"
	DeviceCategoryOther      DeviceCategory = "other"
)

// RiskLevel represents risk assessment levels
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelVeryHigh RiskLevel = "very_high"
	RiskLevelCritical RiskLevel = "critical"
)

// AuthenticityStatus represents device authenticity verification
type AuthenticityStatus string

// OwnershipType represents how the device is owned
type OwnershipType string

const (
	OwnershipPersonal  OwnershipType = "personal"
	OwnershipCorporate OwnershipType = "corporate"
	OwnershipBYOD      OwnershipType = "byod"
	OwnershipCOPE      OwnershipType = "cope" // Corporate Owned, Personally Enabled
	OwnershipCOBO      OwnershipType = "cobo" // Corporate Owned, Business Only
	OwnershipLeased    OwnershipType = "leased"
	OwnershipRented    OwnershipType = "rented"
)

// WaterResistance represents water resistance rating
type WaterResistance string

// BiometricType represents biometric authentication types
type BiometricType string

const (
	BiometricNone        BiometricType = "none"
	BiometricFingerprint BiometricType = "fingerprint"
	BiometricFaceID      BiometricType = "face_id"
	BiometricIris        BiometricType = "iris"
	BiometricMultiple    BiometricType = "multiple"
)

// MDMProvider represents mobile device management providers
type MDMProvider string

const (
	MDMProviderNone       MDMProvider = "none"
	MDMProviderIntune     MDMProvider = "microsoft_intune"
	MDMProviderWorkspace  MDMProvider = "vmware_workspace"
	MDMProviderMobileIron MDMProvider = "mobileiron"
	MDMProviderJamf       MDMProvider = "jamf"
	MDMProviderAirWatch   MDMProvider = "airwatch"
	MDMProviderMaaS360    MDMProvider = "maas360"
	MDMProviderCustom     MDMProvider = "custom"
)

// ============================================
// DEVICE CONSTANTS
// ============================================

const (
	// Risk thresholds
	MinRiskScore          float64 = 0.0
	MaxRiskScore          float64 = 100.0
	LowRiskThreshold      float64 = 30.0
	MediumRiskThreshold   float64 = 50.0
	HighRiskThreshold     float64 = 70.0
	CriticalRiskThreshold float64 = 90.0

	// Battery health thresholds
	BatteryHealthExcellent int = 90
	BatteryHealthGood      int = 80
	BatteryHealthFair      int = 70
	BatteryHealthPoor      int = 60
	BatteryHealthCritical  int = 50

	// Device age thresholds (in days)
	DeviceAgeNew      int = 30
	DeviceAgeRecent   int = 180
	DeviceAgeStandard int = 365
	DeviceAgeOld      int = 730
	DeviceAgeVeryOld  int = 1095

	// Depreciation rates (annual)
	DepreciationRateSmartphone float64 = 0.30
	DepreciationRateTablet     float64 = 0.25
	DepreciationRateLaptop     float64 = 0.20
	DepreciationRateWearable   float64 = 0.35

	// Value thresholds
	HighValueThreshold    float64 = 2000.0
	PremiumValueThreshold float64 = 5000.0
	LuxuryValueThreshold  float64 = 10000.0

	// Insurance eligibility
	MaxInsurableAge   int     = 1825 // 5 years in days
	MinInsurableValue float64 = 100.0
	MaxInsurableValue float64 = 50000.0

	// Inspection intervals (in days)
	InspectionIntervalStandard  int = 365
	InspectionIntervalHighRisk  int = 180
	InspectionIntervalCorporate int = 90

	// Maximum values
	MaxStorageCapacityGB int     = 8192
	MaxRAMGB             int     = 256
	MaxBatteryCycles     int     = 2000
	MaxScreenSizeInches  float64 = 20.0
)

// ============================================
// CUSTOM TYPE METHODS
// ============================================

// NOTE: DeviceStatus methods are defined where the type is declared

// NOTE: DeviceCondition methods are defined where the type is declared

// NOTE: DeviceGrade methods are defined in device_constants.go

// Scan implements sql.Scanner interface
func (rl *RiskLevel) Scan(value interface{}) error {
	if value == nil {
		*rl = RiskLevelLow
		return nil
	}
	if str, ok := value.(string); ok {
		*rl = RiskLevel(str)
		return nil
	}
	return errors.New("cannot scan RiskLevel")
}

// Value implements driver.Valuer interface
func (rl RiskLevel) Value() (driver.Value, error) {
	return string(rl), nil
}

// GetRiskMultiplier returns a multiplier for insurance calculations
func (rl RiskLevel) GetRiskMultiplier() float64 {
	multipliers := map[RiskLevel]float64{
		RiskLevelLow:      1.0,
		RiskLevelMedium:   1.2,
		RiskLevelHigh:     1.5,
		RiskLevelVeryHigh: 2.0,
		RiskLevelCritical: 3.0,
	}
	if mult, exists := multipliers[rl]; exists {
		return mult
	}
	return 1.5
}

// Scan implements sql.Scanner interface
func (ot *OwnershipType) Scan(value interface{}) error {
	if value == nil {
		*ot = OwnershipPersonal
		return nil
	}
	if str, ok := value.(string); ok {
		*ot = OwnershipType(str)
		return nil
	}
	return errors.New("cannot scan OwnershipType")
}

// Value implements driver.Valuer interface
func (ot OwnershipType) Value() (driver.Value, error) {
	return string(ot), nil
}

// IsCorporate checks if ownership is corporate-related
func (ot OwnershipType) IsCorporate() bool {
	switch ot {
	case OwnershipCorporate, OwnershipBYOD, OwnershipCOPE, OwnershipCOBO:
		return true
	}
	return false
}

// RequiresMDM checks if ownership type requires MDM
func (ot OwnershipType) RequiresMDM() bool {
	switch ot {
	case OwnershipCorporate, OwnershipCOPE, OwnershipCOBO:
		return true
	case OwnershipBYOD:
		return false // Optional for BYOD
	}
	return false
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// GetRiskLevel returns the risk level for a given score
func GetRiskLevel(score float64) RiskLevel {
	if score < LowRiskThreshold {
		return RiskLevelLow
	} else if score < MediumRiskThreshold {
		return RiskLevelMedium
	} else if score < HighRiskThreshold {
		return RiskLevelHigh
	} else if score < CriticalRiskThreshold {
		return RiskLevelVeryHigh
	}
	return RiskLevelCritical
}

// GetDeviceConditionFromScore returns condition based on score
func GetDeviceConditionFromScore(score int) DeviceCondition {
	if score >= 95 {
		return DeviceConditionLikeNew
	} else if score >= 90 {
		return DeviceConditionExcellent
	} else if score >= 75 {
		return DeviceConditionGood
	} else if score >= 60 {
		return DeviceConditionFair
	} else if score >= 40 {
		return DeviceConditionPoor
	}
	return DeviceConditionBroken
}

// GetDeviceGradeFromScore returns grade based on score
func GetDeviceGradeFromScore(score int) DeviceGrade {
	if score >= 95 {
		return DeviceGradeA
	} else if score >= 90 {
		return DeviceGradeAMinus
	} else if score >= 85 {
		return DeviceGradeBPlus
	} else if score >= 80 {
		return DeviceGradeB
	} else if score >= 75 {
		return DeviceGradeBMinus
	} else if score >= 70 {
		return DeviceGradeCPlus
	} else if score >= 65 {
		return DeviceGradeC
	} else if score >= 60 {
		return DeviceGradeD
	}
	return DeviceGradeF
}

// ============================================
// JSON ARRAY TYPES
// ============================================

// FunctionalIssues represents a list of functional issues
type FunctionalIssues []string

// Scan implements sql.Scanner interface
func (fi *FunctionalIssues) Scan(value interface{}) error {
	if value == nil {
		*fi = FunctionalIssues{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, fi)
	case string:
		return json.Unmarshal([]byte(v), fi)
	default:
		return fmt.Errorf("cannot scan FunctionalIssues from %T", value)
	}
}

// Value implements driver.Valuer interface
func (fi FunctionalIssues) Value() (driver.Value, error) {
	if len(fi) == 0 {
		return "[]", nil
	}
	data, err := json.Marshal(fi)
	return string(data), err
}

// HasIssue checks if a specific issue exists
func (fi FunctionalIssues) HasIssue(issue string) bool {
	for _, i := range fi {
		if i == issue {
			return true
		}
	}
	return false
}

// Count returns the number of issues
func (fi FunctionalIssues) Count() int {
	return len(fi)
}

// ============================================
// LOCATION DATA TYPE
// ============================================

// DeviceLocation represents device location information
type DeviceLocation struct {
	Country    string    `json:"country"`
	State      string    `json:"state"`
	City       string    `json:"city"`
	PostalCode string    `json:"postal_code"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Accuracy   float64   `json:"accuracy"` // in meters
	UpdatedAt  time.Time `json:"updated_at"`
}

// Scan implements sql.Scanner interface
func (dl *DeviceLocation) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, dl)
	case string:
		return json.Unmarshal([]byte(v), dl)
	default:
		return fmt.Errorf("cannot scan DeviceLocation from %T", value)
	}
}

// Value implements driver.Valuer interface
func (dl DeviceLocation) Value() (driver.Value, error) {
	data, err := json.Marshal(dl)
	return string(data), err
}
