package categories

import (
	"encoding/json"
	"fmt"

	"smartsure/internal/domain/models/device/categories/base"
	"smartsure/internal/domain/models/device/categories/smartphone"
	"smartsure/internal/domain/models/device/categories/smartwatch"
	"smartsure/internal/domain/models/device/categories/tablet"
)

// SpecHelper provides utilities for working with category-specific specifications
type SpecHelper struct{}

// NewSpecHelper creates a new specification helper
func NewSpecHelper() *SpecHelper {
	return &SpecHelper{}
}

// GetSpecForCategory returns the appropriate spec type based on device category
func (h *SpecHelper) GetSpecForCategory(category string, specData json.RawMessage) (base.CategorySpec, error) {
	switch category {
	case "smartphone":
		spec := &smartphone.SmartphoneSpec{}
		if err := json.Unmarshal(specData, spec); err != nil {
			return nil, fmt.Errorf("failed to unmarshal smartphone spec: %w", err)
		}
		return spec, nil

	case "smartwatch":
		spec := &smartwatch.SmartwatchSpec{}
		if err := json.Unmarshal(specData, spec); err != nil {
			return nil, fmt.Errorf("failed to unmarshal smartwatch spec: %w", err)
		}
		return spec, nil

	case "tablet":
		// Import and use tablet specifications
		tabletSpec := &tablet.TabletSpec{}
		return tabletSpec, nil

	default:
		return nil, fmt.Errorf("unsupported device category: %s", category)
	}
}

// ExtractCommonSpecs extracts common fields that exist across categories
// This is useful for backward compatibility and generic operations
func (h *SpecHelper) ExtractCommonSpecs(spec base.CategorySpec) map[string]interface{} {
	common := make(map[string]interface{})

	// Extract base fields that all categories have
	common["manufacturer"] = spec.GetManufacturer()
	common["model"] = spec.GetModel()
	common["market_value"] = spec.GetMarketValue()
	common["category"] = spec.GetCategory()

	// Extract category-specific common fields
	switch s := spec.(type) {
	case *smartphone.SmartphoneSpec:
		common["storage_capacity"] = s.StorageCapacity
		common["ram"] = s.RAM
		common["screen_size"] = s.ScreenSize
		common["screen_type"] = s.ScreenType
		common["battery_capacity"] = s.BatteryCapacity
		common["charging_wattage"] = s.ChargingWattage
		common["wireless_charging"] = s.WirelessCharging
		common["water_resistance"] = s.WaterResistance
		common["biometric_types"] = s.BiometricTypes
		common["is_5g_capable"] = s.Is5GCapable
		common["color"] = s.Color

	case *smartwatch.SmartwatchSpec:
		common["storage_capacity"] = s.StorageCapacity
		common["ram"] = s.RAM / 1024.0              // Convert MB to GB for consistency
		common["screen_size"] = s.ScreenSize / 25.4 // Convert mm to inches
		common["screen_type"] = s.ScreenType
		common["battery_capacity"] = s.BatteryCapacity
		common["charging_wattage"] = int(s.ChargingWattage)
		common["wireless_charging"] = s.WirelessCharging
		common["water_resistance"] = s.WaterResistance
		common["biometric_types"] = s.BiometricTypes
		common["is_5g_capable"] = s.Cellular // Cellular watches can have 5G
		common["color"] = s.Color
	}

	return common
}

// CompareSpecs compares specifications between two devices
func (h *SpecHelper) CompareSpecs(spec1, spec2 base.CategorySpec) map[string]interface{} {
	comparison := make(map[string]interface{})

	// Compare categories first
	if spec1.GetCategory() != spec2.GetCategory() {
		comparison["category_mismatch"] = true
		comparison["category1"] = spec1.GetCategory()
		comparison["category2"] = spec2.GetCategory()
		return comparison
	}

	// Extract common specs for comparison
	specs1 := h.ExtractCommonSpecs(spec1)
	specs2 := h.ExtractCommonSpecs(spec2)

	// Compare each field
	for key, value1 := range specs1 {
		if value2, exists := specs2[key]; exists {
			comparison[key] = map[string]interface{}{
				"device1": value1,
				"device2": value2,
				"same":    value1 == value2,
			}
		}
	}

	return comparison
}

// ValidateSpecsForCategory validates that specs match the expected category
func (h *SpecHelper) ValidateSpecsForCategory(category string, spec base.CategorySpec) error {
	if spec.GetCategory() != base.CategoryType(category) {
		return fmt.Errorf("spec category mismatch: expected %s, got %s", category, spec.GetCategory())
	}

	// Perform category-specific validation
	switch category {
	case "smartphone":
		if _, ok := spec.(*smartphone.SmartphoneSpec); !ok {
			return fmt.Errorf("invalid spec type for smartphone")
		}

	case "smartwatch":
		if _, ok := spec.(*smartwatch.SmartwatchSpec); !ok {
			return fmt.Errorf("invalid spec type for smartwatch")
		}
	}

	// Run the spec's own validation
	if validator, ok := spec.(interface{ Validate() error }); ok {
		return validator.Validate()
	}

	return nil
}

// MigrateFromGenericSpecs helps migrate from old generic specs to category-specific specs
func (h *SpecHelper) MigrateFromGenericSpecs(
	category string,
	genericSpecs map[string]interface{},
) (base.CategorySpec, error) {

	switch category {
	case "smartphone":
		spec := smartphone.NewSmartphoneSpec()

		// Map generic fields to smartphone-specific fields
		if v, ok := genericSpecs["storage_capacity"].(int); ok {
			spec.StorageCapacity = v
		}
		if v, ok := genericSpecs["ram"].(int); ok {
			spec.RAM = v
		}
		if v, ok := genericSpecs["screen_size"].(float64); ok {
			spec.ScreenSize = v
		}
		if v, ok := genericSpecs["screen_type"].(string); ok {
			spec.ScreenType = v
		}
		if v, ok := genericSpecs["battery_capacity"].(int); ok {
			spec.BatteryCapacity = v
		}
		if v, ok := genericSpecs["charging_wattage"].(int); ok {
			spec.ChargingWattage = v
		}
		if v, ok := genericSpecs["wireless_charging"].(bool); ok {
			spec.WirelessCharging = v
		}
		if v, ok := genericSpecs["water_resistance"].(string); ok {
			spec.WaterResistance = v
		}
		if v, ok := genericSpecs["biometric_type"].(string); ok {
			// Convert single biometric type to array
			spec.BiometricTypes = []string{v}
		}
		if v, ok := genericSpecs["is_5g_capable"].(bool); ok {
			spec.Is5GCapable = v
		}
		if v, ok := genericSpecs["color"].(string); ok {
			spec.Color = v
		}

		return spec, nil

	case "smartwatch":
		spec := smartwatch.NewSmartwatchSpec()

		// Map generic fields to smartwatch-specific fields
		// Note: Some conversions needed (inches to mm, GB to MB, etc.)
		if v, ok := genericSpecs["storage_capacity"].(int); ok {
			spec.StorageCapacity = v
		}
		if v, ok := genericSpecs["ram"].(int); ok {
			spec.RAM = v * 1024 // Convert GB to MB
		}
		if v, ok := genericSpecs["screen_size"].(float64); ok {
			spec.ScreenSize = v * 25.4 // Convert inches to mm
		}
		if v, ok := genericSpecs["screen_type"].(string); ok {
			spec.ScreenType = v
		}
		if v, ok := genericSpecs["battery_capacity"].(int); ok {
			spec.BatteryCapacity = v
		}
		if v, ok := genericSpecs["charging_wattage"].(int); ok {
			spec.ChargingWattage = float64(v)
		}
		if v, ok := genericSpecs["wireless_charging"].(bool); ok {
			spec.WirelessCharging = v
		}
		if v, ok := genericSpecs["water_resistance"].(string); ok {
			spec.WaterResistance = v
		}
		if v, ok := genericSpecs["biometric_type"].(string); ok {
			// Convert single biometric type to array
			spec.BiometricTypes = []string{v}
		}
		if v, ok := genericSpecs["is_5g_capable"].(bool); ok {
			spec.Cellular = v // 5G capability implies cellular
		}
		if v, ok := genericSpecs["color"].(string); ok {
			spec.Color = v
		}

		return spec, nil

	default:
		return nil, fmt.Errorf("unsupported category for migration: %s", category)
	}
}
