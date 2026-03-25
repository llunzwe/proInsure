package device

import "strings"

// CheckNetworkCompatibility verifies carrier compatibility
func (d *Device) CheckNetworkCompatibility(carrier string) bool {
	// Check if device is locked to a specific network
	if d.NetworkStatus == "locked" && d.NetworkOperator != carrier {
		return false
	}

	// Check band compatibility based on carrier
	supportedBands := d.GetSupportedBands()

	// Carrier band requirements (simplified)
	carrierBands := map[string][]string{
		"Verizon":  {"B2", "B4", "B5", "B13", "B66"},
		"AT&T":     {"B2", "B4", "B5", "B12", "B17", "B29", "B30"},
		"T-Mobile": {"B2", "B4", "B12", "B66", "B71"},
		"Sprint":   {"B25", "B26", "B41"},
		// International carriers
		"Vodafone": {"B1", "B3", "B7", "B8", "B20"},
		"Orange":   {"B1", "B3", "B7", "B20", "B28"},
		"O2":       {"B1", "B3", "B7", "B8", "B20"},
	}

	requiredBands, exists := carrierBands[carrier]
	if !exists {
		// Unknown carrier, assume compatible
		return true
	}

	// Check if device supports at least the primary bands
	compatibleBands := 0
	for _, required := range requiredBands {
		for _, supported := range supportedBands {
			if supported == required {
				compatibleBands++
				break
			}
		}
	}

	// Need at least 60% band compatibility
	requiredCompatibility := float64(len(requiredBands)) * 0.6
	return float64(compatibleBands) >= requiredCompatibility
}

// IsUnlocked checks if device is carrier unlocked
func (d *Device) IsUnlocked() bool {
	return d.NetworkStatus == "unlocked" || d.NetworkStatus == ""
}

// GetSupportedBands returns supported network bands
func (d *Device) GetSupportedBands() []string {
	bands := []string{}

	// Basic bands supported by all modern smartphones
	basicBands := []string{"B1", "B3", "B7", "B20", "B28"}
	bands = append(bands, basicBands...)

	// Add bands based on device segment and capabilities
	if d.DeviceSegment == "flagship" || d.DeviceSegment == "premium" {
		// Flagship devices support more bands
		premiumBands := []string{"B2", "B4", "B5", "B8", "B12", "B13", "B17", "B25", "B26", "B29", "B30", "B38", "B39", "B40", "B41", "B66", "B71"}
		bands = append(bands, premiumBands...)
	} else if d.DeviceSegment == "mid_range" {
		// Mid-range devices support common bands
		midBands := []string{"B2", "B4", "B5", "B8", "B12", "B17", "B66"}
		bands = append(bands, midBands...)
	}

	// Add 5G bands if capable
	if d.Is5GCapable {
		bands5G := []string{"n1", "n3", "n5", "n7", "n8", "n20", "n28", "n38", "n41", "n66", "n71", "n77", "n78", "n79"}
		bands = append(bands, bands5G...)
	}

	// Region-specific bands based on country code
	if d.CountryCode != "" {
		regionBands := d.getRegionSpecificBands()
		bands = append(bands, regionBands...)
	}

	// Remove duplicates
	uniqueBands := make(map[string]bool)
	result := []string{}
	for _, band := range bands {
		if !uniqueBands[band] {
			uniqueBands[band] = true
			result = append(result, band)
		}
	}

	return result
}

// getRegionSpecificBands returns bands specific to device region
func (d *Device) getRegionSpecificBands() []string {
	regionBands := map[string][]string{
		"US": {"B12", "B13", "B14", "B17", "B29", "B30", "B66", "B71"},
		"EU": {"B8", "B20", "B28", "B32"},
		"CN": {"B34", "B38", "B39", "B40", "B41"},
		"JP": {"B11", "B18", "B19", "B21", "B42"},
		"IN": {"B3", "B5", "B8", "B40", "B41"},
		"KR": {"B3", "B5", "B7", "B8"},
		"AU": {"B28", "B32", "B40"},
	}

	// Get first two characters of country code for region
	region := strings.ToUpper(d.CountryCode[:2])
	if bands, exists := regionBands[region]; exists {
		return bands
	}

	// Default to EU bands for unknown regions
	return regionBands["EU"]
}

// CheckRoamingEligibility checks international roaming capability
func (d *Device) CheckRoamingEligibility(country string) bool {
	// Check if device is unlocked
	if !d.IsUnlocked() {
		return false
	}

	// Check if device is region locked
	if d.RegionLocked {
		// Check if target country is in the same region
		if !d.isSameRegion(country) {
			return false
		}
	}

	// Check band compatibility for common international bands
	internationalBands := []string{"B1", "B3", "B7", "B8", "B20"}
	supportedBands := d.GetSupportedBands()

	compatibleCount := 0
	for _, intBand := range internationalBands {
		for _, supported := range supportedBands {
			if supported == intBand {
				compatibleCount++
				break
			}
		}
	}

	// Need at least 3 international bands for good roaming
	return compatibleCount >= 3
}

// isSameRegion checks if country is in the same region as device
func (d *Device) isSameRegion(country string) bool {
	// Define regions
	regions := map[string][]string{
		"EU": {"AT", "BE", "BG", "HR", "CY", "CZ", "DK", "EE", "FI", "FR", "DE", "GR", "HU", "IE", "IT", "LV", "LT", "LU", "MT", "NL", "PL", "PT", "RO", "SK", "SI", "ES", "SE"},
		"NA": {"US", "CA", "MX"},
		"SA": {"AR", "BR", "CL", "CO", "PE", "VE"},
		"AS": {"CN", "JP", "KR", "IN", "SG", "MY", "TH", "VN", "ID", "PH"},
		"OC": {"AU", "NZ"},
		"AF": {"ZA", "NG", "EG", "KE", "GH"},
	}

	// Find device region
	deviceRegion := ""
	deviceCountry := strings.ToUpper(d.CountryCode[:2])
	for region, countries := range regions {
		for _, c := range countries {
			if c == deviceCountry {
				deviceRegion = region
				break
			}
		}
		if deviceRegion != "" {
			break
		}
	}

	// Find target country region
	targetCountry := strings.ToUpper(country[:2])
	for region, countries := range regions {
		for _, c := range countries {
			if c == targetCountry {
				return region == deviceRegion
			}
		}
	}

	return false
}

// ValidateESIMCapability verifies eSIM support
func (d *Device) ValidateESIMCapability() bool {
	// Check dual SIM type
	if d.DualSIMType == "eSIM" || d.DualSIMType == "hybrid" {
		return true
	}

	// Check by device model and brand (simplified list)
	eSIMModels := map[string][]string{
		"Apple": {
			"iPhone XS", "iPhone XS Max", "iPhone XR", "iPhone 11",
			"iPhone 11 Pro", "iPhone 11 Pro Max", "iPhone SE",
			"iPhone 12", "iPhone 12 Mini", "iPhone 12 Pro", "iPhone 12 Pro Max",
			"iPhone 13", "iPhone 13 Mini", "iPhone 13 Pro", "iPhone 13 Pro Max",
			"iPhone 14", "iPhone 14 Plus", "iPhone 14 Pro", "iPhone 14 Pro Max",
			"iPhone 15", "iPhone 15 Plus", "iPhone 15 Pro", "iPhone 15 Pro Max",
		},
		"Google": {
			"Pixel 3", "Pixel 3 XL", "Pixel 3a", "Pixel 3a XL",
			"Pixel 4", "Pixel 4 XL", "Pixel 4a", "Pixel 5",
			"Pixel 6", "Pixel 6 Pro", "Pixel 6a",
			"Pixel 7", "Pixel 7 Pro", "Pixel 7a",
			"Pixel 8", "Pixel 8 Pro",
		},
		"Samsung": {
			"Galaxy S20", "Galaxy S20+", "Galaxy S20 Ultra",
			"Galaxy S21", "Galaxy S21+", "Galaxy S21 Ultra",
			"Galaxy S22", "Galaxy S22+", "Galaxy S22 Ultra",
			"Galaxy S23", "Galaxy S23+", "Galaxy S23 Ultra",
			"Galaxy Z Fold", "Galaxy Z Fold2", "Galaxy Z Fold3", "Galaxy Z Fold4", "Galaxy Z Fold5",
			"Galaxy Z Flip", "Galaxy Z Flip3", "Galaxy Z Flip4", "Galaxy Z Flip5",
		},
	}

	// Check if device model supports eSIM
	if models, exists := eSIMModels[d.Brand]; exists {
		for _, model := range models {
			if strings.Contains(d.Model, model) {
				return true
			}
		}
	}

	// Check by device segment and age
	if d.DeviceSegment == "flagship" || d.DeviceSegment == "premium" {
		// Most flagship devices from 2018 onwards support eSIM
		deviceAge := d.GetDeviceAge() / 365
		if deviceAge < 6 {
			// Likely supports eSIM if it's a recent flagship
			return true
		}
	}

	return false
}

// GetNetworkCapabilities returns detailed network capabilities
func (d *Device) GetNetworkCapabilities() map[string]interface{} {
	capabilities := map[string]interface{}{
		"is_unlocked":      d.IsUnlocked(),
		"network_operator": d.NetworkOperator,
		"network_status":   d.NetworkStatus,
		"is_5g_capable":    d.Is5GCapable,
		"dual_sim_type":    d.DualSIMType,
		"esim_capable":     d.ValidateESIMCapability(),
		"supported_bands":  d.GetSupportedBands(),
		"region_locked":    d.RegionLocked,
		"country_code":     d.CountryCode,
	}

	// Add carrier compatibility for major carriers
	carriers := []string{"Verizon", "AT&T", "T-Mobile", "Vodafone", "Orange"}
	compatibility := make(map[string]bool)
	for _, carrier := range carriers {
		compatibility[carrier] = d.CheckNetworkCompatibility(carrier)
	}
	capabilities["carrier_compatibility"] = compatibility

	// Add roaming capabilities for major regions
	roamingRegions := []string{"US", "EU", "UK", "CN", "JP", "AU"}
	roaming := make(map[string]bool)
	for _, region := range roamingRegions {
		roaming[region] = d.CheckRoamingEligibility(region)
	}
	capabilities["roaming_eligibility"] = roaming

	return capabilities
}
