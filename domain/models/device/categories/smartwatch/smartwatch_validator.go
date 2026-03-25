package smartwatch

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	
	"smartsure/internal/domain/models/device/categories/base"
)

// SmartwatchValidator provides smartwatch-specific validation
// Works with existing device validation infrastructure
type SmartwatchValidator struct {
	base.BaseValidator
}

// NewSmartwatchValidator creates a new smartwatch validator
func NewSmartwatchValidator() *SmartwatchValidator {
	return &SmartwatchValidator{
		BaseValidator: base.BaseValidator{},
	}
}

// ValidateSpec validates smartwatch specifications
func (v *SmartwatchValidator) ValidateSpec(spec base.CategorySpec) error {
	swSpec, ok := spec.(*SmartwatchSpec)
	if !ok {
		return fmt.Errorf("invalid spec type: expected SmartwatchSpec")
	}

	// Basic validation
	if err := swSpec.Validate(); err != nil {
		return err
	}

	// Additional smartwatch-specific validations
	if err := v.validateDisplay(swSpec); err != nil {
		return err
	}

	if err := v.validateHealthSensors(swSpec); err != nil {
		return err
	}

	if err := v.validateBattery(swSpec); err != nil {
		return err
	}

	if err := v.validateConnectivity(swSpec); err != nil {
		return err
	}

	if err := v.validateCompatibility(swSpec); err != nil {
		return err
	}

	return nil
}

// ValidateIMEI validates IMEI for cellular-enabled smartwatches
func (v *SmartwatchValidator) ValidateIMEI(imei string) error {
	// For smartwatches, IMEI is optional (only cellular models have it)
	if imei == "" {
		return nil // Empty IMEI is valid for non-cellular models
	}

	// For cellular models, use base IMEI validation
	return v.BaseValidator.ValidateIMEI(imei)
}

// ValidateSerialNumber validates smartwatch serial number
func (v *SmartwatchValidator) ValidateSerialNumber(serial string) error {
	// Try to detect manufacturer from serial format
	if len(serial) == 12 && strings.HasPrefix(serial, "F") {
		return v.validateAppleWatchSerial(serial)
	} else if len(serial) >= 11 && len(serial) <= 15 {
		return v.validateSamsungWatchSerial(serial)
	} else if len(serial) >= 9 && len(serial) <= 10 {
		return v.validateGarminSerial(serial)
	}

	// Generic serial validation
	return v.BaseValidator.ValidateSerialNumber(serial)
}

// ValidateModel validates model compatibility with manufacturer
func (v *SmartwatchValidator) ValidateModel(manufacturer, model string) error {
	// Check known manufacturer-model combinations
	validModels := map[string][]string{
		"apple": {
			"watch se", "watch series 3", "watch series 4", "watch series 5",
			"watch series 6", "watch series 7", "watch series 8", "watch series 9",
			"watch ultra", "watch ultra 2",
		},
		"samsung": {
			"galaxy watch", "galaxy watch 3", "galaxy watch 4", "galaxy watch 5",
			"galaxy watch 6", "gear s3", "gear sport", "galaxy watch active",
		},
		"garmin": {
			"fenix 6", "fenix 7", "forerunner 245", "forerunner 945",
			"vivoactive 4", "venu 2", "instinct 2", "epix",
		},
		"fitbit": {
			"versa 2", "versa 3", "versa 4", "sense", "sense 2",
		},
	}

	lowerManufacturer := strings.ToLower(manufacturer)
	lowerModel := strings.ToLower(model)

	if models, exists := validModels[lowerManufacturer]; exists {
		for _, validModel := range models {
			if strings.Contains(lowerModel, validModel) {
				return nil
			}
		}
		// If we have the manufacturer but model doesn't match
		// Just issue a warning, don't fail
	}

	// For unknown manufacturers or models, just check basic format
	if len(model) < 2 || len(model) > 50 {
		return fmt.Errorf("model name must be between 2 and 50 characters")
	}

	return nil
}

// ValidateAge validates device age for insurance eligibility
func (v *SmartwatchValidator) ValidateAge(releaseDate time.Time) error {
	age := time.Since(releaseDate)
	maxAge := 3 * 365 * 24 * time.Hour // 3 years

	if age > maxAge {
		return fmt.Errorf("smartwatch is too old for insurance (max 3 years)")
	}

	if age < 0 {
		return fmt.Errorf("release date cannot be in the future")
	}

	return nil
}

// ValidateCondition validates device condition
func (v *SmartwatchValidator) ValidateCondition(condition string) error {
	validConditions := []string{
		"excellent", "mint", "like_new",
		"good", "fair", "poor",
		"broken", "for_parts",
	}

	lowerCondition := strings.ToLower(condition)
	for _, valid := range validConditions {
		if lowerCondition == valid {
			// Only excellent, good conditions eligible for new insurance
			if lowerCondition == "poor" || lowerCondition == "broken" || lowerCondition == "for_parts" {
				return fmt.Errorf("device must be in good or excellent condition for insurance")
			}
			return nil
		}
	}

	return fmt.Errorf("invalid condition: %s", condition)
}

// Private validation methods

func (v *SmartwatchValidator) validateDisplay(spec *SmartwatchSpec) error {
	// Validate display type
	validDisplayTypes := []string{"oled", "amoled", "lcd", "e-ink", "retina", "super amoled"}
	displayLower := strings.ToLower(spec.DisplayType)
	valid := false
	for _, dt := range validDisplayTypes {
		if displayLower == dt {
			valid = true
			break
		}
	}
	if !valid && spec.DisplayType != "" {
		return fmt.Errorf("invalid display type: %s", spec.DisplayType)
	}

	// Validate resolution format
	if spec.Resolution != "" {
		resPattern := regexp.MustCompile(`^\d+x\d+$`)
		if !resPattern.MatchString(spec.Resolution) {
			return fmt.Errorf("invalid resolution format: %s (expected format: WIDTHxHEIGHT)", spec.Resolution)
		}
	}

	return nil
}

func (v *SmartwatchValidator) validateHealthSensors(spec *SmartwatchSpec) error {
	// ECG requires heart rate sensor
	if spec.ECGSensor && !spec.HeartRateSensor {
		return fmt.Errorf("ECG sensor requires heart rate sensor")
	}

	// Validate health features
	validHealthFeatures := []string{
		"sleep_tracking", "stress_monitoring", "breathing_exercises",
		"menstrual_tracking", "blood_glucose", "body_composition",
		"fall_detection", "emergency_sos", "medication_reminders",
	}

	for _, feature := range spec.HealthFeatures {
		valid := false
		featureLower := strings.ToLower(feature)
		for _, vf := range validHealthFeatures {
			if featureLower == vf {
				valid = true
				break
			}
		}
		if !valid {
			// Don't fail on unknown features, just skip
			continue
		}
	}

	return nil
}

func (v *SmartwatchValidator) validateBattery(spec *SmartwatchSpec) error {
	// Validate charging method
	validChargingMethods := []string{"magnetic", "wireless", "pins", "pogo", "usb-c", "proprietary"}
	chargingLower := strings.ToLower(spec.ChargingMethod)
	valid := false
	for _, cm := range validChargingMethods {
		if chargingLower == cm {
			valid = true
			break
		}
	}
	if !valid && spec.ChargingMethod != "" {
		return fmt.Errorf("invalid charging method: %s", spec.ChargingMethod)
	}

	// Validate charging time
	if spec.ChargingTime < 0 || spec.ChargingTime > 300 {
		return fmt.Errorf("charging time must be between 0 and 300 minutes")
	}

	return nil
}

func (v *SmartwatchValidator) validateConnectivity(spec *SmartwatchSpec) error {
	// Validate Bluetooth version
	if spec.BluetoothVersion != "" {
		btPattern := regexp.MustCompile(`^(4\.[0-2]|5\.[0-3])$`)
		if !btPattern.MatchString(spec.BluetoothVersion) {
			return fmt.Errorf("invalid Bluetooth version: %s", spec.BluetoothVersion)
		}
	}

	// Cellular requires certain features
	if spec.Cellular && !spec.GPSBuiltIn {
		// Warning: Cellular usually comes with GPS
	}

	return nil
}

func (v *SmartwatchValidator) validateCompatibility(spec *SmartwatchSpec) error {
	// Validate OS
	validOS := []string{"watchos", "wear os", "tizen", "garmin os", "fitbit os", "harmony os", "proprietary"}
	osLower := strings.ToLower(spec.OS)
	valid := false
	for _, os := range validOS {
		if osLower == os {
			valid = true
			break
		}
	}
	if !valid && spec.OS != "" {
		return fmt.Errorf("invalid operating system: %s", spec.OS)
	}

	// Check OS compatibility with manufacturer
	osCompatibility := map[string][]string{
		"apple":   {"watchos"},
		"samsung": {"wear os", "tizen"},
		"google":  {"wear os"},
		"garmin":  {"garmin os"},
		"fitbit":  {"fitbit os"},
		"huawei":  {"harmony os"},
	}

	if compatibleOS, exists := osCompatibility[strings.ToLower(spec.Manufacturer)]; exists {
		compatible := false
		for _, cos := range compatibleOS {
			if osLower == cos {
				compatible = true
				break
			}
		}
		if !compatible && spec.OS != "" {
			// Just a warning, not a hard failure
		}
	}

	return nil
}

// Apple Watch serial validation
func (v *SmartwatchValidator) validateAppleWatchSerial(serial string) error {
	// Apple Watch serials are typically 12 characters
	if len(serial) != 12 {
		return fmt.Errorf("Apple Watch serial must be 12 characters")
	}

	// Should be alphanumeric
	alphaNum := regexp.MustCompile(`^[A-Z0-9]+$`)
	if !alphaNum.MatchString(strings.ToUpper(serial)) {
		return fmt.Errorf("Apple Watch serial must be alphanumeric")
	}

	return nil
}

// Samsung Watch serial validation
func (v *SmartwatchValidator) validateSamsungWatchSerial(serial string) error {
	// Samsung serials are typically 11-15 characters
	if len(serial) < 11 || len(serial) > 15 {
		return fmt.Errorf("Samsung Watch serial must be 11-15 characters")
	}

	return nil
}

// ValidateRepairCost validates repair cost estimates for smartwatch components
func (v *SmartwatchValidator) ValidateRepairCost(component string, cost float64) error {
	component = strings.ToLower(component)

	// Define reasonable cost ranges for smartwatch components
	costRanges := map[string][2]float64{
		"screen":        {50, 300},   // Screen repair
		"battery":       {40, 150},   // Battery replacement
		"charging_port": {30, 100},   // Charging port repair
		"case":          {60, 200},   // Case replacement
		"strap":         {15, 80},    // Strap replacement
		"sensors":       {25, 120},   // Sensor repair/replacement
		"buttons":       {20, 70},    // Button repair
		"water_sealing": {40, 180},   // Water sealing repair
	}

	if range_, exists := costRanges[component]; exists {
		if cost < range_[0] {
			return fmt.Errorf("repair cost %.2f for %s is below minimum expected cost %.2f", cost, component, range_[0])
		}
		if cost > range_[1] {
			return fmt.Errorf("repair cost %.2f for %s exceeds maximum expected cost %.2f", cost, component, range_[1])
		}
	}

	// General validation for unknown components
	if cost <= 0 {
		return fmt.Errorf("repair cost must be greater than zero")
	}

	if cost > 500 {
		return fmt.Errorf("repair cost %.2f seems unreasonably high for a smartwatch component", cost)
	}

	return nil
}

// IsCompatible checks if accessories/parts are compatible with smartwatch
func (v *SmartwatchValidator) IsCompatible(itemType, itemModel string) bool {
	itemType = strings.ToLower(itemType)
	itemModel = strings.ToLower(itemModel)

	switch itemType {
	case "band", "strap", "wristband":
		// Check for common band sizes: 20mm, 22mm, 24mm, 26mm
		return strings.Contains(itemModel, "20mm") ||
			   strings.Contains(itemModel, "22mm") ||
			   strings.Contains(itemModel, "24mm") ||
			   strings.Contains(itemModel, "26mm") ||
			   strings.Contains(itemModel, "standard") ||
			   strings.Contains(itemModel, "universal")

	case "charger", "charging_cable", "dock":
		// Smartwatches typically use proprietary charging
		return strings.Contains(itemModel, "magnetic") ||
			   strings.Contains(itemModel, "proprietary") ||
			   strings.Contains(itemModel, "oem") ||
			   strings.Contains(itemModel, "original")

	case "screen_protector", "screen_guard":
		// Screen protectors for curved/flat screens
		return strings.Contains(itemModel, "curved") ||
			   strings.Contains(itemModel, "flat") ||
			   strings.Contains(itemModel, "tempered") ||
			   strings.Contains(itemModel, "ceramic")

	case "case", "cover":
		// Smartwatch cases
		return strings.Contains(itemModel, "silicone") ||
			   strings.Contains(itemModel, "plastic") ||
			   strings.Contains(itemModel, "metal") ||
			   strings.Contains(itemModel, "leather")

	default:
		// Generic compatibility check
		return strings.Contains(itemModel, "smartwatch") ||
			   strings.Contains(itemModel, "wearable") ||
			   strings.Contains(itemModel, "fitness")
	}
}

// Garmin serial validation
func (v *SmartwatchValidator) validateGarminSerial(serial string) error {
	// Garmin serials are typically 9-10 characters
	if len(serial) < 9 || len(serial) > 10 {
		return fmt.Errorf("Garmin serial must be 9-10 characters")
	}

	return nil
}
