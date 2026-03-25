package smartphone

import (
	"fmt"
	"regexp"
	"strings"
	
	"smartsure/internal/domain/models/device/categories/base"
)

// SmartphoneValidator validates smartphone-specific data
type SmartphoneValidator struct {
	base.BaseValidator
	supportedManufacturers map[string]bool
	supportedNetworks      map[string]bool
	minSpecs               SmartphoneMinSpecs
}

// SmartphoneMinSpecs defines minimum specifications for insurance
type SmartphoneMinSpecs struct {
	MinRAM        int
	MinStorage    int
	MinBattery    int
	MinScreenSize float64
	MaxScreenSize float64
	MinOSVersion  map[string]float64 // iOS: 14.0, Android: 10.0, etc.
}

// NewSmartphoneValidator creates a new smartphone validator
func NewSmartphoneValidator() *SmartphoneValidator {
	validator := &SmartphoneValidator{
		BaseValidator: *base.NewBaseValidator(),
		supportedManufacturers: map[string]bool{
			"Apple":    true,
			"Samsung":  true,
			"Google":   true,
			"OnePlus":  true,
			"Xiaomi":   true,
			"Huawei":   true,
			"Sony":     true,
			"LG":       true,
			"Motorola": true,
			"Nokia":    true,
			"Asus":     true,
			"Oppo":     true,
			"Vivo":     true,
			"Realme":   true,
		},
		supportedNetworks: map[string]bool{
			"2G":  true,
			"3G":  true,
			"4G":  true,
			"LTE": true,
			"5G":  true,
		},
		minSpecs: SmartphoneMinSpecs{
			MinRAM:        1,
			MinStorage:    8,
			MinBattery:    1000,
			MinScreenSize: 3.0,
			MaxScreenSize: 8.0,
			MinOSVersion: map[string]float64{
				"iOS":     12.0,
				"Android": 8.0,
			},
		},
	}

	// Add smartphone-specific validation rules
	validator.setupValidationRules()

	return validator
}

// ValidateSpec validates smartphone specifications
func (v *SmartphoneValidator) ValidateSpec(spec base.CategorySpec) error {
	phoneSpec, ok := spec.(*SmartphoneSpec)
	if !ok {
		return fmt.Errorf("invalid specification type: expected SmartphoneSpec")
	}

	// Validate basic spec requirements
	if err := phoneSpec.Validate(); err != nil {
		return err
	}

	// Validate manufacturer
	if !v.supportedManufacturers[phoneSpec.Manufacturer] {
		return fmt.Errorf("unsupported manufacturer: %s", phoneSpec.Manufacturer)
	}

	// Validate specifications against minimums
	if phoneSpec.RAM < v.minSpecs.MinRAM {
		return fmt.Errorf("RAM below minimum requirement: %dGB < %dGB",
			phoneSpec.RAM, v.minSpecs.MinRAM)
	}

	if phoneSpec.StorageCapacity < v.minSpecs.MinStorage {
		return fmt.Errorf("storage below minimum requirement: %dGB < %dGB",
			phoneSpec.StorageCapacity, v.minSpecs.MinStorage)
	}

	if phoneSpec.BatteryCapacity < v.minSpecs.MinBattery {
		return fmt.Errorf("battery capacity below minimum: %dmAh < %dmAh",
			phoneSpec.BatteryCapacity, v.minSpecs.MinBattery)
	}

	// Validate screen size
	if phoneSpec.ScreenSize < v.minSpecs.MinScreenSize ||
		phoneSpec.ScreenSize > v.minSpecs.MaxScreenSize {
		return fmt.Errorf("screen size out of range: %.1f inches (valid: %.1f-%.1f)",
			phoneSpec.ScreenSize, v.minSpecs.MinScreenSize, v.minSpecs.MaxScreenSize)
	}

	// Validate network support
	if err := v.validateNetworkSupport(phoneSpec.NetworkSupport); err != nil {
		return err
	}

	// Validate cameras
	if err := v.validateCameras(phoneSpec); err != nil {
		return err
	}

	// Validate OS version
	if err := v.validateOSVersion(phoneSpec.OSVersion); err != nil {
		return err
	}

	// Validate special features
	if err := v.validateSpecialFeatures(phoneSpec); err != nil {
		return err
	}

	return nil
}

// ValidateIMEI validates smartphone IMEI with manufacturer checks
func (v *SmartphoneValidator) ValidateIMEI(imei string) error {
	// First, use base IMEI validation
	if err := v.BaseValidator.ValidateIMEI(imei); err != nil {
		return err
	}

	// Check TAC (Type Allocation Code) - first 8 digits
	if len(imei) >= 8 {
		tac := imei[:8]
		if err := v.validateTAC(tac); err != nil {
			return fmt.Errorf("invalid TAC: %w", err)
		}
	}

	return nil
}

// ValidateModel validates if model exists and is valid for insurance
func (v *SmartphoneValidator) ValidateModel(manufacturer, model string) error {
	if manufacturer == "" {
		return fmt.Errorf("manufacturer cannot be empty")
	}

	if model == "" {
		return fmt.Errorf("model cannot be empty")
	}

	// Check if manufacturer is supported
	if !v.supportedManufacturers[manufacturer] {
		return fmt.Errorf("manufacturer %s is not supported for insurance", manufacturer)
	}

	// Validate model format based on manufacturer
	if err := v.validateModelFormat(manufacturer, model); err != nil {
		return err
	}

	// Check for blacklisted models
	if v.isModelBlacklisted(manufacturer, model) {
		return fmt.Errorf("model %s %s is not eligible for insurance", manufacturer, model)
	}

	return nil
}

// ValidateCondition validates device condition for insurance eligibility
func (v *SmartphoneValidator) ValidateCondition(condition string) error {
	validConditions := map[string]bool{
		"mint":      true,
		"excellent": true,
		"good":      true,
		"fair":      true,
		"poor":      false, // Not eligible
		"broken":    false, // Not eligible
	}

	eligible, exists := validConditions[strings.ToLower(condition)]
	if !exists {
		return fmt.Errorf("invalid condition: %s", condition)
	}

	if !eligible {
		return fmt.Errorf("devices in %s condition are not eligible for insurance", condition)
	}

	return nil
}

// ValidateRepairCost validates if repair cost is reasonable
func (v *SmartphoneValidator) ValidateRepairCost(component string, cost float64) error {
	maxCosts := map[string]float64{
		"screen":        500.0,
		"battery":       150.0,
		"camera":        300.0,
		"charging_port": 100.0,
		"speaker":       80.0,
		"microphone":    60.0,
		"motherboard":   800.0,
		"back_glass":    200.0,
		"buttons":       50.0,
	}

	maxCost, exists := maxCosts[component]
	if !exists {
		// Unknown component, use general validation
		if cost > 1000.0 {
			return fmt.Errorf("repair cost exceeds maximum limit for %s", component)
		}
		return nil
	}

	if cost > maxCost {
		return fmt.Errorf("repair cost for %s exceeds maximum: $%.2f > $%.2f",
			component, cost, maxCost)
	}

	if cost < 0 {
		return fmt.Errorf("repair cost cannot be negative")
	}

	return nil
}

// IsCompatible checks if accessories/parts are compatible
func (v *SmartphoneValidator) IsCompatible(itemType, itemModel string) bool {
	// Check general smartphone compatibility
	compatibleTypes := map[string]bool{
		"case":             true,
		"screen_protector": true,
		"charger":          true,
		"cable":            true,
		"headphones":       true,
		"power_bank":       true,
		"car_mount":        true,
		"wireless_charger": true,
	}

	return compatibleTypes[itemType]
}

// Helper methods

func (v *SmartphoneValidator) setupValidationRules() {
	// IMEI validation rule
	v.Rules["imei"] = base.ValidationRule{
		Field:    "imei",
		Type:     "string",
		Required: true,
		Pattern:  `^\d{15}$`,
		ErrorMsg: "IMEI must be exactly 15 digits",
	}

	// Serial number validation rule
	v.Rules["serial_number"] = base.ValidationRule{
		Field:    "serial_number",
		Type:     "string",
		Required: true,
		Pattern:  `^[A-Z0-9]{8,20}$`,
		ErrorMsg: "Serial number must be 8-20 alphanumeric characters",
	}

	// Model validation rule
	v.Rules["model"] = base.ValidationRule{
		Field:    "model",
		Type:     "string",
		Required: true,
		ErrorMsg: "Model is required",
	}

	// Screen size validation
	v.Rules["screen_size"] = base.ValidationRule{
		Field:    "screen_size",
		Type:     "float",
		Required: true,
		Min:      3.0,
		Max:      8.0,
		ErrorMsg: "Screen size must be between 3.0 and 8.0 inches",
	}

	// RAM validation
	v.Rules["ram"] = base.ValidationRule{
		Field:    "ram",
		Type:     "int",
		Required: true,
		Min:      1,
		Max:      32,
		ErrorMsg: "RAM must be between 1 and 32 GB",
	}

	// Storage validation
	v.Rules["storage"] = base.ValidationRule{
		Field:    "storage",
		Type:     "int",
		Required: true,
		Min:      8,
		Max:      2048,
		ErrorMsg: "Storage must be between 8 and 2048 GB",
	}
}

func (v *SmartphoneValidator) validateTAC(tac string) error {
	// TAC validation (Type Allocation Code)
	// In a real system, this would check against a TAC database
	if len(tac) != 8 {
		return fmt.Errorf("TAC must be 8 digits")
	}

	// Check if all digits
	for _, c := range tac {
		if c < '0' || c > '9' {
			return fmt.Errorf("TAC must contain only digits")
		}
	}

	return nil
}

func (v *SmartphoneValidator) validateModelFormat(manufacturer, model string) error {
	switch manufacturer {
	case "Apple":
		// iPhone model format: iPhone [number/name]
		if !strings.HasPrefix(model, "iPhone") {
			return fmt.Errorf("Apple model must start with 'iPhone'")
		}
	case "Samsung":
		// Samsung format: Galaxy [Series] [Model]
		if !strings.Contains(model, "Galaxy") && !strings.HasPrefix(model, "SM-") {
			return fmt.Errorf("Samsung model format invalid")
		}
	case "Google":
		// Google format: Pixel [number/name]
		if !strings.Contains(model, "Pixel") {
			return fmt.Errorf("Google model must contain 'Pixel'")
		}
	}

	return nil
}

func (v *SmartphoneValidator) isModelBlacklisted(manufacturer, model string) bool {
	// Check blacklisted models (e.g., prototype devices, development models)
	blacklist := map[string][]string{
		"Apple": {
			"iPhone Developer",
			"iPhone Prototype",
		},
		"Samsung": {
			"Galaxy Test",
			"SM-TEST",
		},
	}

	if blacklistedModels, exists := blacklist[manufacturer]; exists {
		for _, blacklisted := range blacklistedModels {
			if strings.Contains(model, blacklisted) {
				return true
			}
		}
	}

	return false
}

func (v *SmartphoneValidator) validateNetworkSupport(networks []string) error {
	if len(networks) == 0 {
		return fmt.Errorf("at least one network type must be supported")
	}

	hasValidNetwork := false
	for _, network := range networks {
		if v.supportedNetworks[network] {
			hasValidNetwork = true
			break
		}
	}

	if !hasValidNetwork {
		return fmt.Errorf("device must support at least one valid network type")
	}

	return nil
}

func (v *SmartphoneValidator) validateCameras(spec *SmartphoneSpec) error {
	// Validate rear cameras
	if len(spec.RearCameras) == 0 {
		return fmt.Errorf("at least one rear camera is required")
	}

	for i, camera := range spec.RearCameras {
		if camera.Megapixels <= 0 {
			return fmt.Errorf("rear camera %d: invalid megapixel count", i+1)
		}
		if camera.Megapixels > 200 {
			return fmt.Errorf("rear camera %d: unrealistic megapixel count", i+1)
		}
	}

	// Validate front cameras
	if len(spec.FrontCameras) == 0 {
		return fmt.Errorf("at least one front camera is required")
	}

	for i, camera := range spec.FrontCameras {
		if camera.Megapixels <= 0 {
			return fmt.Errorf("front camera %d: invalid megapixel count", i+1)
		}
		if camera.Megapixels > 100 {
			return fmt.Errorf("front camera %d: unrealistic megapixel count", i+1)
		}
	}

	return nil
}

func (v *SmartphoneValidator) validateOSVersion(osVersion string) error {
	if osVersion == "" {
		return fmt.Errorf("OS version is required")
	}

	// Extract OS name and version
	parts := strings.Fields(osVersion)
	if len(parts) < 2 {
		return fmt.Errorf("invalid OS version format")
	}

	osName := parts[0]

	// Check minimum OS versions
	if minVersion, exists := v.minSpecs.MinOSVersion[osName]; exists {
		// Parse version number
		versionPattern := regexp.MustCompile(`\d+\.?\d*`)
		matches := versionPattern.FindString(osVersion)
		if matches == "" {
			return fmt.Errorf("could not parse OS version number")
		}

		var version float64
		fmt.Sscanf(matches, "%f", &version)

		if version < minVersion {
			return fmt.Errorf("%s version %.1f is below minimum requirement %.1f",
				osName, version, minVersion)
		}
	}

	return nil
}

func (v *SmartphoneValidator) validateSpecialFeatures(spec *SmartphoneSpec) error {
	// Validate foldable devices have additional requirements
	if spec.IsFoldable() {
		if spec.ScreenProtection == "" {
			return fmt.Errorf("foldable devices must have screen protection specified")
		}

		// Foldables should have minimum specs
		if spec.RAM < 6 {
			return fmt.Errorf("foldable devices require minimum 6GB RAM")
		}

		if spec.StorageCapacity < 128 {
			return fmt.Errorf("foldable devices require minimum 128GB storage")
		}
	}

	return nil
}

// GetValidationRules returns all validation rules
func (v *SmartphoneValidator) GetValidationRules() map[string]base.ValidationRule {
	return v.Rules
}
