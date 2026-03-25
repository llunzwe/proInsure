package device

import (
	"fmt"
	"strings"
	"time"
)

// MeetsRegulatoryRequirements checks regional compliance
func (d *Device) MeetsRegulatoryRequirements(region string) bool {
	// Check basic device information
	if d.IMEI == "" || d.Model == "" || d.Brand == "" {
		return false
	}

	// Check IMEI validity
	if !d.ValidateIMEI() {
		return false
	}

	// Region-specific requirements
	switch strings.ToUpper(region) {
	case "US":
		return d.meetsFCCRequirements()
	case "EU":
		return d.meetsCERequirements()
	case "UK":
		return d.meetsUKCARequirements()
	case "CN":
		return d.meetsCCCRequirements()
	case "IN":
		return d.meetsBISRequirements()
	case "AU":
		return d.meetsACMARequirements()
	case "JP":
		return d.meetsJATERequirements()
	default:
		// Unknown region, check basic safety
		return d.CheckSARCompliance()
	}
}

// meetsFCCRequirements checks US FCC compliance
func (d *Device) meetsFCCRequirements() bool {
	// Check SAR limits (US: 1.6 W/kg)
	if !d.CheckSARCompliance() {
		return false
	}

	// Check if device is grey market
	if d.GreyMarket && !d.ImportedDevice {
		return false
	}

	// Check band compatibility for US
	usBands := []string{"B2", "B4", "B12", "B13", "B66", "B71"}
	supportedBands := d.GetSupportedBands()

	compatibleCount := 0
	for _, usBand := range usBands {
		for _, supported := range supportedBands {
			if supported == usBand {
				compatibleCount++
				break
			}
		}
	}

	// Need at least 3 US bands
	return compatibleCount >= 3
}

// meetsCERequirements checks EU CE marking requirements
func (d *Device) meetsCERequirements() bool {
	// Check SAR limits (EU: 2.0 W/kg body, 4.0 W/kg limbs)
	if !d.CheckSARCompliance() {
		return false
	}

	// Check RoHS compliance (assumed based on age and brand)
	if d.GetDeviceAge() > 3650 { // Devices older than 10 years might not comply
		return false
	}

	// Check band compatibility for EU
	euBands := []string{"B1", "B3", "B7", "B8", "B20", "B28"}
	supportedBands := d.GetSupportedBands()

	compatibleCount := 0
	for _, euBand := range euBands {
		for _, supported := range supportedBands {
			if supported == euBand {
				compatibleCount++
				break
			}
		}
	}

	return compatibleCount >= 4
}

// meetsUKCARequirements checks UK UKCA marking requirements
func (d *Device) meetsUKCARequirements() bool {
	// Similar to CE requirements for now
	return d.meetsCERequirements()
}

// meetsCCCRequirements checks China CCC requirements
func (d *Device) meetsCCCRequirements() bool {
	// Check if device is approved for Chinese market
	if d.CountryCode == "CN" || strings.Contains(d.CountryCode, "CN") {
		return true
	}

	// Check band compatibility for China
	cnBands := []string{"B1", "B3", "B38", "B39", "B40", "B41"}
	supportedBands := d.GetSupportedBands()

	compatibleCount := 0
	for _, cnBand := range cnBands {
		for _, supported := range supportedBands {
			if supported == cnBand {
				compatibleCount++
				break
			}
		}
	}

	return compatibleCount >= 3
}

// meetsBISRequirements checks India BIS requirements
func (d *Device) meetsBISRequirements() bool {
	// Check SAR limits (India: 1.6 W/kg)
	if !d.CheckSARCompliance() {
		return false
	}

	// Check if device is imported with proper documentation
	if d.ImportedDevice && !d.ProofOfOwnershipVerified {
		return false
	}

	return true
}

// meetsACMARequirements checks Australia ACMA requirements
func (d *Device) meetsACMARequirements() bool {
	// Check band compatibility for Australia
	auBands := []string{"B1", "B3", "B7", "B28", "B32", "B40"}
	supportedBands := d.GetSupportedBands()

	compatibleCount := 0
	for _, auBand := range auBands {
		for _, supported := range supportedBands {
			if supported == auBand {
				compatibleCount++
				break
			}
		}
	}

	return compatibleCount >= 3
}

// meetsJATERequirements checks Japan JATE/Telec requirements
func (d *Device) meetsJATERequirements() bool {
	// Check band compatibility for Japan
	jpBands := []string{"B1", "B3", "B8", "B11", "B18", "B19", "B21", "B42"}
	supportedBands := d.GetSupportedBands()

	compatibleCount := 0
	for _, jpBand := range jpBands {
		for _, supported := range supportedBands {
			if supported == jpBand {
				compatibleCount++
				break
			}
		}
	}

	return compatibleCount >= 4
}

// GetCertifications returns safety/regulatory certifications
func (d *Device) GetCertifications() []string {
	certifications := []string{}

	// Basic certifications most devices have
	if d.CheckSARCompliance() {
		certifications = append(certifications, "SAR_Compliant")
	}

	// Brand-specific certifications
	majorBrands := []string{"Apple", "Samsung", "Google", "OnePlus", "Xiaomi", "Oppo", "Vivo"}
	for _, brand := range majorBrands {
		if d.Brand == brand {
			certifications = append(certifications, "RoHS", "WEEE")
			break
		}
	}

	// Water resistance certifications
	if d.WaterResistance != "" && d.WaterResistance != "none" {
		certifications = append(certifications, d.WaterResistance)
	}

	// Regional certifications based on compliance checks
	if d.meetsFCCRequirements() {
		certifications = append(certifications, "FCC")
	}
	if d.meetsCERequirements() {
		certifications = append(certifications, "CE")
	}
	if d.meetsBISRequirements() {
		certifications = append(certifications, "BIS")
	}

	// Network certifications
	if d.Is5GCapable {
		certifications = append(certifications, "5G_Certified")
	}
	if d.ValidateESIMCapability() {
		certifications = append(certifications, "eSIM_Compatible")
	}

	// Security certifications for certain models
	if d.BiometricType == "both" || d.BiometricType == "face_id" {
		certifications = append(certifications, "Biometric_Security")
	}

	return certifications
}

// IsImportCompliant verifies import compliance
func (d *Device) IsImportCompliant() bool {
	// Must be properly documented
	if !d.ProofOfOwnershipVerified {
		return false
	}

	// Check if device is marked as imported
	if !d.ImportedDevice && d.GreyMarket {
		return false // Grey market without proper import flag
	}

	// Check IMEI validity
	if !d.ValidateIMEI() {
		return false
	}

	// Check blacklist status
	if d.BlacklistStatus == "blocked" {
		return false
	}

	// Check if device meets local regulatory requirements
	if d.CountryCode != "" {
		region := strings.ToUpper(d.CountryCode[:2])
		if !d.MeetsRegulatoryRequirements(region) {
			return false
		}
	}

	return true
}

// GenerateComplianceReport creates compliance summary
func (d *Device) GenerateComplianceReport() map[string]interface{} {
	report := map[string]interface{}{
		"device_id":        d.ID,
		"imei":             d.IMEI,
		"model":            d.Model,
		"brand":            d.Brand,
		"report_date":      time.Now(),
		"imei_valid":       d.ValidateIMEI(),
		"import_compliant": d.IsImportCompliant(),
		"certifications":   d.GetCertifications(),
	}

	// Regional compliance checks
	regions := []string{"US", "EU", "UK", "CN", "IN", "AU", "JP"}
	regionalCompliance := make(map[string]bool)
	for _, region := range regions {
		regionalCompliance[region] = d.MeetsRegulatoryRequirements(region)
	}
	report["regional_compliance"] = regionalCompliance

	// Safety compliance
	report["sar_compliant"] = d.CheckSARCompliance()

	// Import/Export status
	report["grey_market"] = d.GreyMarket
	report["imported_device"] = d.ImportedDevice
	report["region_locked"] = d.RegionLocked

	// Blacklist and authenticity
	report["blacklist_status"] = d.BlacklistStatus
	report["authenticity_status"] = d.AuthenticityStatus

	// Network compliance
	report["network_status"] = d.NetworkStatus
	report["is_unlocked"] = d.IsUnlocked()

	// Environmental compliance
	report["environmental_grade"] = d.GetEnvironmentalGrade()
	report["recyclability_score"] = d.CalculateRecyclabilityScore()

	// Compliance issues
	issues := []string{}
	if !d.ValidateIMEI() {
		issues = append(issues, "Invalid IMEI")
	}
	if d.BlacklistStatus == "blocked" {
		issues = append(issues, "Device blacklisted")
	}
	if d.GreyMarket && !d.ImportedDevice {
		issues = append(issues, "Undocumented grey market device")
	}
	if d.AuthenticityStatus == "suspicious" {
		issues = append(issues, "Authenticity concerns")
	}
	report["compliance_issues"] = issues

	// Overall compliance status
	overallCompliant := len(issues) == 0 && d.IsImportCompliant()
	report["overall_compliant"] = overallCompliant

	return report
}

// CheckSARCompliance verifies radiation safety compliance
func (d *Device) CheckSARCompliance() bool {
	// SAR (Specific Absorption Rate) compliance
	// Different regions have different limits:
	// US/Canada: 1.6 W/kg averaged over 1 gram
	// EU: 2.0 W/kg averaged over 10 grams (head/body), 4.0 W/kg (limbs)
	// India: 1.6 W/kg
	// Australia: 2.0 W/kg

	// Major brands typically comply with strictest standards
	trustedBrands := []string{
		"Apple", "Samsung", "Google", "OnePlus", "Xiaomi",
		"Oppo", "Vivo", "Sony", "LG", "Motorola", "Nokia",
		"Huawei", "Honor", "Realme", "Asus", "HTC",
	}

	brandCompliant := false
	for _, brand := range trustedBrands {
		if strings.EqualFold(d.Brand, brand) {
			brandCompliant = true
			break
		}
	}

	if !brandCompliant {
		// Unknown brand, check other factors
		if d.GreyMarket {
			return false // Grey market devices might not comply
		}
		if d.AuthenticityStatus == "suspicious" {
			return false // Suspicious devices might be counterfeit
		}
	}

	// Check device age - very old devices might not have been tested
	if d.GetDeviceAge() > 5475 { // Older than 15 years
		return false
	}

	// Device segment check - all modern legitimate devices should comply
	legitimateSegments := []string{"flagship", "premium", "mid_range", "budget"}
	segmentValid := false
	for _, segment := range legitimateSegments {
		if d.DeviceSegment == segment {
			segmentValid = true
			break
		}
	}

	return brandCompliant && segmentValid
}

// CheckForRecalls verifies against recall databases
func (d *Device) CheckForRecalls() []Recall {
	recalls := []Recall{}

	// Simulated recall database check
	// In production, this would query actual recall databases

	// Samsung Note 7 battery recall example
	if d.Brand == "Samsung" && strings.Contains(d.Model, "Note 7") {
		recalls = append(recalls, Recall{
			ID:          "CPSC-2016-001",
			Title:       "Battery Fire Hazard",
			Description: "Device battery may overheat and cause fire",
			Date:        time.Date(2016, 9, 15, 0, 0, 0, 0, time.UTC),
			Severity:    "Critical",
			Action:      "Stop use immediately and return for refund",
		})
	}

	// Check for battery-related recalls based on battery health
	if d.BatteryHealth < 50 && d.BatteryHealth > 0 && d.GetDeviceAge() > 1095 {
		// Old device with poor battery might have recall
		recalls = append(recalls, Recall{
			ID:          fmt.Sprintf("BATTERY-%s", d.ID),
			Title:       "Battery Degradation Warning",
			Description: "Battery showing signs of significant degradation",
			Date:        time.Now(),
			Severity:    "Warning",
			Action:      "Consider battery replacement",
		})
	}

	return recalls
}

// Recall represents a product recall
type Recall struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Severity    string    `json:"severity"`
	Action      string    `json:"action"`
}
