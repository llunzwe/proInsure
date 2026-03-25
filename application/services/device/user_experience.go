package device

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// GetUserSatisfactionScore calculates user satisfaction based on device issues
func (d *Device) GetUserSatisfactionScore() float64 {
	score := 100.0

	// Deduct for repairs
	repairCount := d.GetRepairCount()
	if repairCount > 0 {
		score -= float64(repairCount) * 5 // 5 points per repair
	}

	// Deduct for claims
	claimCount := len(d.Claims)
	if claimCount > 0 {
		score -= float64(claimCount) * 8 // 8 points per claim
	}

	// Deduct for poor condition
	conditionPenalties := map[string]float64{
		"poor":      30,
		"fair":      15,
		"good":      5,
		"excellent": 0,
	}

	if penalty, exists := conditionPenalties[d.Condition]; exists {
		score -= penalty
	}

	// Deduct for component issues
	if d.ScreenCondition == "broken" {
		score -= 20
	} else if d.ScreenCondition == "cracked" {
		score -= 10
	}

	if !d.TouchScreenResponsive {
		score -= 25
	}

	if d.BatteryHealth > 0 && d.BatteryHealth < 70 {
		score -= 15
	}

	if d.ChargingPortCondition != "working" {
		score -= 10
	}

	// Performance impact
	performanceScore := d.GetPerformanceScore()
	if performanceScore < 50 {
		score -= 20
	} else if performanceScore < 70 {
		score -= 10
	}

	// Bonus for good maintenance
	if d.LastInspection != nil && time.Since(*d.LastInspection) < 180*24*time.Hour {
		score += 5
	}

	// Bonus for security features enabled
	if d.FindMyDeviceEnabled && d.RemoteLockEnabled && d.EncryptionEnabled {
		score += 5
	}

	// Ensure score is between 0 and 100
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	return score
}

// GenerateSupportTicket creates a support ticket for an issue
func (d *Device) GenerateSupportTicket(issue string) string {
	// Generate ticket ID
	ticketID := fmt.Sprintf("TKT-%s-%d",
		strings.ToUpper(d.ID.String()[:8]),
		time.Now().Unix())

	// Determine priority based on issue
	priority := d.determineSupportPriority(issue)

	// Create ticket details (in production, this would create an actual ticket)
	ticket := map[string]interface{}{
		"ticket_id":    ticketID,
		"device_id":    d.ID,
		"device_model": fmt.Sprintf("%s %s", d.Brand, d.Model),
		"imei":         d.IMEI,
		"issue":        issue,
		"priority":     priority,
		"created_at":   time.Now(),
		"status":       "open",
	}

	// Add relevant device information
	if strings.Contains(strings.ToLower(issue), "battery") {
		ticket["battery_health"] = d.BatteryHealth
		ticket["battery_cycles"] = d.BatteryCycles
	}

	if strings.Contains(strings.ToLower(issue), "screen") {
		ticket["screen_condition"] = d.ScreenCondition
		ticket["touch_responsive"] = d.TouchScreenResponsive
	}

	if strings.Contains(strings.ToLower(issue), "water") {
		ticket["water_damage"] = d.WaterDamageIndicator
		ticket["water_resistance"] = d.WaterResistance
	}

	return ticketID
}

// determineSupportPriority determines the priority of a support issue
func (d *Device) determineSupportPriority(issue string) string {
	issueLower := strings.ToLower(issue)

	// Critical issues
	criticalKeywords := []string{"stolen", "lost", "fraud", "blacklist", "emergency", "fire", "smoke"}
	for _, keyword := range criticalKeywords {
		if strings.Contains(issueLower, keyword) {
			return "critical"
		}
	}

	// High priority issues
	highKeywords := []string{"not working", "broken", "dead", "won't turn on", "no power", "water damage"}
	for _, keyword := range highKeywords {
		if strings.Contains(issueLower, keyword) {
			return "high"
		}
	}

	// Medium priority issues
	mediumKeywords := []string{"slow", "intermittent", "sometimes", "cracked", "battery", "charging"}
	for _, keyword := range mediumKeywords {
		if strings.Contains(issueLower, keyword) {
			return "medium"
		}
	}

	// Default to low priority
	return "low"
}

// GetRecommendedUpgrades suggests upgrade options based on device state
func (d *Device) GetRecommendedUpgrades() []string {
	upgrades := []string{}

	// Age-based recommendations
	deviceAge := d.GetDeviceAge()
	if deviceAge > 1095 { // Older than 3 years
		upgrades = append(upgrades, "Consider upgrading to a newer model for better performance and features")
	}

	// Storage recommendations
	if d.StorageCapacity > 0 && d.StorageCapacity <= 32 {
		upgrades = append(upgrades, "Upgrade to a device with more storage (64GB or higher)")
	}

	// RAM recommendations
	if d.RAM > 0 && d.RAM <= 3 {
		upgrades = append(upgrades, "Upgrade to a device with at least 4GB RAM for smoother performance")
	}

	// 5G recommendation
	if !d.Is5GCapable && d.DeviceSegment != "budget" {
		upgrades = append(upgrades, "Consider a 5G-capable device for future-proof connectivity")
	}

	// Battery life recommendation
	if d.BatteryHealth > 0 && d.BatteryHealth < 70 {
		upgrades = append(upgrades, "Upgrade to a new device with fresh battery for all-day usage")
	}

	// Screen recommendation
	if d.ScreenCondition == "broken" || d.ScreenCondition == "cracked" {
		upgrades = append(upgrades, "Upgrade to a new device instead of expensive screen repair")
	}

	// Water resistance recommendation
	if d.WaterResistance == "none" || d.WaterResistance == "" {
		upgrades = append(upgrades, "Consider a water-resistant device (IP67/IP68) for better protection")
	}

	// Segment-based recommendations
	switch d.DeviceSegment {
	case "budget":
		upgrades = append(upgrades, "Consider a mid-range device for significantly better performance")
	case "mid_range":
		if deviceAge > 730 { // Older than 2 years
			upgrades = append(upgrades, "Premium devices offer better longevity and features")
		}
	}

	// Security recommendations
	if d.BiometricType == "none" || d.BiometricType == "" {
		upgrades = append(upgrades, "Upgrade to a device with biometric security (fingerprint/face)")
	}

	// Camera recommendations
	if d.CameraCondition != "all_working" {
		upgrades = append(upgrades, "Upgrade to a device with advanced camera features")
	}

	// Trade-in opportunity
	if d.IsEligibleForTradeIn() {
		tradeInValue := d.CalculateTradeInValue()
		if tradeInValue > 100 {
			upgrades = append(upgrades, fmt.Sprintf("Trade-in eligible: Get $%.2f towards a new device", tradeInValue))
		}
	}

	// Swap program eligibility
	if deviceAge > 365 && d.Grade != "F" {
		upgrades = append(upgrades, "Eligible for device swap program - upgrade with loyalty benefits")
	}

	return upgrades
}

// IsEligibleForLoaner checks if user can get a loaner device
func (d *Device) IsEligibleForLoaner() bool {
	// Check basic eligibility
	if !d.LoanerDeviceEligible {
		return false
	}

	// Must be verified
	if !d.IsVerified {
		return false
	}

	// Must not be stolen or blacklisted
	if d.IsStolen || d.BlacklistStatus == "blocked" {
		return false
	}

	// Check if under repair
	hasActiveRepair := false
	for _, repair := range d.Repairs {
		if repair.RepairStatus == "in_progress" || repair.RepairStatus == "pending" {
			hasActiveRepair = true
			break
		}
	}

	// Must have active repair or claim
	hasActiveClaim := d.HasActiveClaim()
	if !hasActiveRepair && !hasActiveClaim {
		return false
	}

	// Check device value (loaner for valuable devices)
	if d.CurrentValue < 200 {
		return false
	}

	// Check customer standing (no excessive claims)
	if d.GetClaimFrequency() > 3 {
		return false
	}

	// Premium segment devices get priority
	if d.DeviceSegment == "flagship" || d.DeviceSegment == "premium" {
		return true
	}

	// Check repair duration for other segments
	estimatedRepairDays := 3 // Default estimate
	if d.WaterDamageIndicator != "white" {
		estimatedRepairDays = 7
	}
	if d.ScreenCondition == "broken" {
		estimatedRepairDays = 5
	}

	// Loaner only for longer repairs
	return estimatedRepairDays >= 3
}

// GetWarrantyClaimOptions returns available warranty claim types
func (d *Device) GetWarrantyClaimOptions() []string {
	options := []string{}

	// Check if warranty is active
	if !d.IsWarrantyActive() {
		return options
	}

	// Manufacturing defects
	if d.GetDeviceAge() < 365 { // Within first year
		options = append(options, "Manufacturing defect")

		// Battery issues within first year
		if d.BatteryHealth > 0 && d.BatteryHealth < 80 {
			options = append(options, "Battery defect")
		}

		// Screen issues not caused by damage
		if d.ScreenCondition != "perfect" && d.ScreenCondition != "minor_scratches" {
			options = append(options, "Display defect")
		}

		// Component failures
		if d.ChargingPortCondition != "working" {
			options = append(options, "Charging port defect")
		}
		if d.SpeakerCondition != "working" {
			options = append(options, "Speaker defect")
		}
		if d.MicrophoneCondition != "working" {
			options = append(options, "Microphone defect")
		}
		if d.CameraCondition != "all_working" {
			options = append(options, "Camera defect")
		}
		if !d.TouchScreenResponsive {
			options = append(options, "Touch screen defect")
		}
	}

	// Software issues
	if d.OSVersion != "" {
		options = append(options, "Software malfunction")
	}

	// Dead pixels (usually covered under warranty)
	if d.DeadPixels {
		options = append(options, "Dead pixel warranty")
	}

	// Extended warranty options
	if d.WarrantyExpiry != nil && time.Until(*d.WarrantyExpiry) > 365*24*time.Hour {
		// Has extended warranty
		options = append(options, "Extended warranty claim")

		// Additional coverage under extended warranty
		if d.ScreenCondition == "cracked" {
			options = append(options, "Accidental damage (extended)")
		}
	}

	return options
}

// GetCustomerServiceHistory returns interaction history
func (d *Device) GetCustomerServiceHistory() map[string]interface{} {
	history := map[string]interface{}{
		"device_id":          d.ID,
		"registration_date":  d.RegistrationDate,
		"total_claims":       len(d.Claims),
		"total_repairs":      d.GetRepairCount(),
		"satisfaction_score": d.GetUserSatisfactionScore(),
	}

	// Calculate response metrics
	totalInteractions := len(d.Claims) + d.GetRepairCount()
	if totalInteractions > 0 {
		// Simulated metrics
		history["avg_resolution_time_days"] = 5
		history["first_contact_resolution_rate"] = 0.75
	}

	// Add recent issues
	recentIssues := []string{}
	if d.ScreenCondition == "cracked" || d.ScreenCondition == "broken" {
		recentIssues = append(recentIssues, "Screen damage")
	}
	if d.BatteryHealth > 0 && d.BatteryHealth < 70 {
		recentIssues = append(recentIssues, "Battery degradation")
	}
	if d.ChargingPortCondition != "working" {
		recentIssues = append(recentIssues, "Charging issues")
	}
	history["recent_issues"] = recentIssues

	// Add service recommendations
	history["recommended_services"] = d.GetEligibleServices()

	// Add warranty status
	history["warranty_active"] = d.IsWarrantyActive()
	history["warranty_options"] = d.GetWarrantyClaimOptions()

	// Add loaner eligibility
	history["loaner_eligible"] = d.IsEligibleForLoaner()

	return history
}

// GenerateCustomerInsights provides insights about device usage
func (d *Device) GenerateCustomerInsights() map[string]interface{} {
	insights := map[string]interface{}{}

	// Usage pattern
	deviceAge := d.GetDeviceAge()
	if deviceAge > 1460 { // 4+ years
		insights["usage_pattern"] = "Long-term user"
		insights["loyalty_level"] = "High"
	} else if deviceAge > 730 { // 2+ years
		insights["usage_pattern"] = "Regular user"
		insights["loyalty_level"] = "Medium"
	} else {
		insights["usage_pattern"] = "New user"
		insights["loyalty_level"] = "Building"
	}

	// Care level
	careScore := 100.0
	if d.Condition == "excellent" && d.ScreenCondition == "perfect" {
		careScore = 100
	} else if d.Condition == "good" {
		careScore = 80
	} else if d.Condition == "fair" {
		careScore = 60
	} else {
		careScore = 40
	}

	// Adjust for repairs
	repairCount := d.GetRepairCount()
	if repairCount > 3 {
		careScore -= 20
	} else if repairCount > 1 {
		careScore -= 10
	}

	insights["device_care_score"] = math.Max(careScore, 0)

	// Risk profile
	riskLevel := "Low"
	if d.GetClaimFrequency() > 2 {
		riskLevel = "High"
	} else if d.GetClaimFrequency() > 1 {
		riskLevel = "Medium"
	}
	insights["risk_profile"] = riskLevel

	// Upgrade likelihood
	upgradeScore := 0.0
	if deviceAge > 1095 { // 3+ years
		upgradeScore += 40
	} else if deviceAge > 730 { // 2+ years
		upgradeScore += 20
	}

	if d.Condition == "poor" || d.Condition == "fair" {
		upgradeScore += 30
	}

	if d.BatteryHealth > 0 && d.BatteryHealth < 70 {
		upgradeScore += 20
	}

	insights["upgrade_likelihood"] = fmt.Sprintf("%.0f%%", math.Min(upgradeScore, 100))

	// Service preferences
	preferences := []string{}
	if d.GetRepairCount() > 0 {
		preferences = append(preferences, "Repair services")
	}
	if len(d.Claims) > 0 {
		preferences = append(preferences, "Insurance claims")
	}
	if len(d.Subscriptions) > 0 {
		preferences = append(preferences, "Subscription services")
	}
	insights["service_preferences"] = preferences

	// Value consciousness
	valueScore := "Standard"
	if d.DeviceSegment == "budget" {
		valueScore = "Price-conscious"
	} else if d.DeviceSegment == "flagship" {
		valueScore = "Premium-oriented"
	}
	insights["value_orientation"] = valueScore

	return insights
}
