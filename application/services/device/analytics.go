package device

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// InsuranceReport represents a comprehensive insurance report
type InsuranceReport struct {
	DeviceID         string                 `json:"device_id"`
	ReportDate       time.Time              `json:"report_date"`
	DeviceInfo       map[string]interface{} `json:"device_info"`
	RiskAssessment   map[string]interface{} `json:"risk_assessment"`
	Coverage         map[string]interface{} `json:"coverage"`
	Claims           []ClaimSummary         `json:"claims"`
	Recommendations  []string               `json:"recommendations"`
	FinancialSummary map[string]float64     `json:"financial_summary"`
}

// AuditEntry represents an audit trail entry
type AuditEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Action      string    `json:"action"`
	Details     string    `json:"details"`
	PerformedBy string    `json:"performed_by"`
}

// GenerateInsuranceReport creates comprehensive insurance report
func (d *Device) GenerateInsuranceReport() InsuranceReport {
	report := InsuranceReport{
		DeviceID:   d.ID.String(),
		ReportDate: time.Now(),
	}

	// Device Information
	report.DeviceInfo = map[string]interface{}{
		"brand":           d.Brand,
		"model":           d.Model,
		"imei":            d.IMEI,
		"serial_number":   d.SerialNumber,
		"purchase_date":   d.PurchaseDate,
		"purchase_price":  d.PurchasePrice,
		"current_value":   d.CurrentValue,
		"market_value":    d.MarketValue,
		"age_days":        d.GetDeviceAge(),
		"condition":       d.Condition,
		"grade":           d.Grade,
		"segment":         d.DeviceSegment,
		"warranty_active": d.IsWarrantyActive(),
		"warranty_expiry": d.WarrantyExpiry,
	}

	// Risk Assessment
	report.RiskAssessment = map[string]interface{}{
		"risk_score":          d.CalculateRiskScore(),
		"theft_risk_level":    d.TheftRiskLevel,
		"fraud_risk_score":    d.FraudRiskScore,
		"suspicious_activity": d.GetSuspiciousActivityScore(),
		"blacklist_status":    d.BlacklistStatus,
		"authenticity_status": d.AuthenticityStatus,
		"security_score":      d.GetSecurityScore(),
		"fraud_patterns":      d.DetectFraudPatterns(),
	}

	// Coverage Information
	report.Coverage = map[string]interface{}{
		"insurable":            d.CanBeInsured(),
		"monthly_premium":      d.CalculateInsurancePremium(),
		"coverage_limit":       d.GetCoverageLimit(),
		"eligible_services":    d.GetEligibleServices(),
		"active_subscriptions": d.HasActiveSubscription(),
		"active_financing":     d.HasActiveFinancing(),
		"active_rental":        d.HasActiveRental(),
	}

	// Claims History
	report.Claims = d.GetClaimHistory()

	// Financial Summary
	report.FinancialSummary = map[string]float64{
		"total_ownership_cost": d.CalculateTotalOwnershipCost(),
		"depreciation":         d.CalculateDepreciation(),
		"replacement_cost":     d.CalculateReplacementCost(),
		"trade_in_value":       d.CalculateTradeInValue(),
		"salvage_value":        d.CalculateSalvageValue(),
		"total_repair_cost":    d.GetTotalRepairCost(),
		"claim_frequency":      d.GetClaimFrequency(),
	}

	// Recommendations
	recommendations := []string{}

	// Insurance recommendations
	if d.CanBeInsured() && !d.HasActiveSubscription() {
		recommendations = append(recommendations, "Consider comprehensive insurance coverage")
	}

	// Maintenance recommendations
	if d.RequiresInspection() {
		recommendations = append(recommendations, "Schedule device inspection")
	}

	// Security recommendations
	if !d.FindMyDeviceEnabled {
		recommendations = append(recommendations, "Enable Find My Device for theft protection")
	}

	// Upgrade recommendations
	if d.GetDeviceAge() > 1095 && d.Condition != "excellent" {
		recommendations = append(recommendations, "Consider device upgrade or trade-in")
	}

	// Risk mitigation
	if d.RiskScore > 70 {
		recommendations = append(recommendations, "High risk device - consider additional coverage")
	}

	report.Recommendations = recommendations

	return report
}

// GetUsageStatistics returns device usage statistics
func (d *Device) GetUsageStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"device_id":         d.ID,
		"registration_date": d.RegistrationDate,
		"days_owned":        d.GetDeviceAge(),
		"years_owned":       float64(d.GetDeviceAge()) / 365,
	}

	// Service usage
	serviceStats := map[string]int{
		"total_repairs":        d.GetRepairCount(),
		"total_claims":         len(d.Claims),
		"active_subscriptions": 0,
		"completed_swaps":      0,
		"trade_ins":            0,
		"marketplace_listings": 0,
	}

	// Count active subscriptions
	for _, sub := range d.Subscriptions {
		if sub.IsActive() {
			serviceStats["active_subscriptions"]++
		}
	}

	// Count swaps
	for _, swap := range d.Swaps {
		if swap.SwapStatus == "completed" {
			serviceStats["completed_swaps"]++
		}
	}

	// Count trade-ins
	serviceStats["trade_ins"] = len(d.TradeIns)

	// Count marketplace listings
	serviceStats["marketplace_listings"] = len(d.MarketListings)

	stats["service_usage"] = serviceStats

	// Condition tracking
	conditionStats := map[string]interface{}{
		"current_condition": d.Condition,
		"condition_score":   d.GetConditionScore(),
		"screen_condition":  d.ScreenCondition,
		"battery_health":    d.BatteryHealth,
		"water_damage":      d.WaterDamageIndicator,
		"performance_score": d.GetPerformanceScore(),
	}
	stats["condition_metrics"] = conditionStats

	// Financial metrics
	financialStats := map[string]float64{
		"total_spent":          d.CalculateTotalOwnershipCost(),
		"current_value":        d.CurrentValue,
		"value_retention_rate": (d.CurrentValue / d.PurchasePrice) * 100,
		"monthly_cost":         d.CalculateTotalOwnershipCost() / math.Max(float64(d.GetDeviceAge())/30, 1),
	}
	stats["financial_metrics"] = financialStats

	// Claim statistics
	if len(d.Claims) > 0 {
		approvedClaims := 0
		totalClaimAmount := 0.0
		for _, claim := range d.Claims {
			if claim.Status == "approved" || claim.Status == "paid" {
				approvedClaims++
				totalClaimAmount += claim.ApprovedAmount
			}
		}

		claimStats := map[string]interface{}{
			"total_claims":       len(d.Claims),
			"approved_claims":    approvedClaims,
			"approval_rate":      float64(approvedClaims) / float64(len(d.Claims)) * 100,
			"total_claim_amount": totalClaimAmount,
			"avg_claim_amount":   totalClaimAmount / math.Max(float64(approvedClaims), 1),
		}
		stats["claim_statistics"] = claimStats
	}

	return stats
}

// CalculateROI calculates return on investment
func (d *Device) CalculateROI() float64 {
	// Calculate total investment
	totalInvestment := d.PurchasePrice

	// Add insurance premiums
	monthsOwned := math.Max(float64(d.GetDeviceAge())/30, 1)
	monthlyPremium := d.CalculateInsurancePremium()
	totalInvestment += monthlyPremium * monthsOwned

	// Add repair costs not covered by insurance
	totalInvestment += d.GetTotalRepairCost() * 0.2 // Assume 80% covered by insurance

	// Calculate returns/value
	totalValue := d.CurrentValue

	// Add value from claims paid out
	for _, claim := range d.Claims {
		if claim.Status == "paid" {
			totalValue += claim.ApprovedAmount
		}
	}

	// Add saved replacement costs due to repairs
	savedReplacementCost := 0.0
	if d.GetRepairCount() > 0 {
		// Each repair potentially saved a device replacement
		savedReplacementCost = d.CalculateReplacementCost() * 0.3 * float64(d.GetRepairCount())
	}
	totalValue += savedReplacementCost

	// Calculate ROI
	if totalInvestment == 0 {
		return 0
	}

	roi := ((totalValue - totalInvestment) / totalInvestment) * 100

	return roi
}

// GetComparativeAnalysis compares with another device
func (d *Device) GetComparativeAnalysis(other *Device) map[string]interface{} {
	if other == nil {
		return map[string]interface{}{
			"error": "No comparison device provided",
		}
	}

	analysis := map[string]interface{}{
		"device_1": map[string]interface{}{
			"id":    d.ID,
			"model": fmt.Sprintf("%s %s", d.Brand, d.Model),
		},
		"device_2": map[string]interface{}{
			"id":    other.ID,
			"model": fmt.Sprintf("%s %s", other.Brand, other.Model),
		},
	}

	// Value comparison
	valueComparison := map[string]interface{}{
		"device_1_value":   d.CurrentValue,
		"device_2_value":   other.CurrentValue,
		"value_difference": d.CurrentValue - other.CurrentValue,
		"value_winner":     d.getWinner(d.CurrentValue, other.CurrentValue, true),
	}
	analysis["value_comparison"] = valueComparison

	// Condition comparison
	conditionComparison := map[string]interface{}{
		"device_1_condition":       d.Condition,
		"device_2_condition":       other.Condition,
		"device_1_condition_score": d.GetConditionScore(),
		"device_2_condition_score": other.GetConditionScore(),
		"condition_winner":         d.getWinner(d.GetConditionScore(), other.GetConditionScore(), true),
	}
	analysis["condition_comparison"] = conditionComparison

	// Performance comparison
	performanceComparison := map[string]interface{}{
		"device_1_performance": d.GetPerformanceScore(),
		"device_2_performance": other.GetPerformanceScore(),
		"performance_winner":   d.getWinner(d.GetPerformanceScore(), other.GetPerformanceScore(), true),
	}
	analysis["performance_comparison"] = performanceComparison

	// Risk comparison
	riskComparison := map[string]interface{}{
		"device_1_risk": d.CalculateRiskScore(),
		"device_2_risk": other.CalculateRiskScore(),
		"risk_winner":   d.getWinner(d.CalculateRiskScore(), other.CalculateRiskScore(), false), // Lower risk is better
	}
	analysis["risk_comparison"] = riskComparison

	// Insurance comparison
	insuranceComparison := map[string]interface{}{
		"device_1_premium":   d.CalculateInsurancePremium(),
		"device_2_premium":   other.CalculateInsurancePremium(),
		"premium_difference": d.CalculateInsurancePremium() - other.CalculateInsurancePremium(),
		"device_1_insurable": d.CanBeInsured(),
		"device_2_insurable": other.CanBeInsured(),
	}
	analysis["insurance_comparison"] = insuranceComparison

	// Age comparison
	ageComparison := map[string]interface{}{
		"device_1_age_days": d.GetDeviceAge(),
		"device_2_age_days": other.GetDeviceAge(),
		"age_difference":    d.GetDeviceAge() - other.GetDeviceAge(),
	}
	analysis["age_comparison"] = ageComparison

	// Feature comparison
	featureComparison := map[string]interface{}{
		"5g_capability": map[string]bool{
			"device_1": d.Is5GCapable,
			"device_2": other.Is5GCapable,
		},
		"water_resistance": map[string]string{
			"device_1": d.WaterResistance,
			"device_2": other.WaterResistance,
		},
		"storage": map[string]int{
			"device_1": d.StorageCapacity,
			"device_2": other.StorageCapacity,
		},
		"ram": map[string]int{
			"device_1": d.RAM,
			"device_2": other.RAM,
		},
	}
	analysis["feature_comparison"] = featureComparison

	// Overall recommendation
	score1 := d.GetConditionScore()*0.3 + d.GetPerformanceScore()*0.3 + (100-d.CalculateRiskScore())*0.2 + (d.CurrentValue/d.PurchasePrice)*20
	score2 := other.GetConditionScore()*0.3 + other.GetPerformanceScore()*0.3 + (100-other.CalculateRiskScore())*0.2 + (other.CurrentValue/other.PurchasePrice)*20

	if score1 > score2 {
		analysis["recommendation"] = "Device 1 is the better choice overall"
	} else if score2 > score1 {
		analysis["recommendation"] = "Device 2 is the better choice overall"
	} else {
		analysis["recommendation"] = "Both devices are equally matched"
	}

	return analysis
}

// getWinner determines winner in comparison
func (d *Device) getWinner(value1, value2 float64, higherBetter bool) string {
	if higherBetter {
		if value1 > value2 {
			return "device_1"
		} else if value2 > value1 {
			return "device_2"
		}
	} else {
		if value1 < value2 {
			return "device_1"
		} else if value2 < value1 {
			return "device_2"
		}
	}
	return "tie"
}

// GenerateAuditTrail returns complete audit history
func (d *Device) GenerateAuditTrail() []AuditEntry {
	audit := []AuditEntry{}

	// Registration
	audit = append(audit, AuditEntry{
		Timestamp:   d.RegistrationDate,
		Action:      "Device Registered",
		Details:     fmt.Sprintf("Device %s %s registered with IMEI %s", d.Brand, d.Model, d.IMEI),
		PerformedBy: d.OwnerID.String(),
	})

	// Purchase
	if d.PurchaseDate != nil {
		audit = append(audit, AuditEntry{
			Timestamp:   *d.PurchaseDate,
			Action:      "Device Purchased",
			Details:     fmt.Sprintf("Purchased for %s %.2f", d.Currency, d.PurchasePrice),
			PerformedBy: d.OwnerID.String(),
		})
	}

	// Verification
	if d.VerificationDate != nil {
		audit = append(audit, AuditEntry{
			Timestamp:   *d.VerificationDate,
			Action:      "Device Verified",
			Details:     "Device authenticity and ownership verified",
			PerformedBy: "System",
		})
	}

	// Claims
	for _, claim := range d.Claims {
		audit = append(audit, AuditEntry{
			Timestamp:   claim.CreatedAt,
			Action:      "Claim Filed",
			Details:     fmt.Sprintf("Claim %s filed for %s", claim.ID, claim.Type),
			PerformedBy: d.OwnerID.String(),
		})
	}

	// Repairs
	for i, repair := range d.Repairs {
		audit = append(audit, AuditEntry{
			Timestamp:   repair.CreatedAt,
			Action:      "Repair Initiated",
			Details:     fmt.Sprintf("Repair #%d for %s", i+1, repair.IssueDescription),
			PerformedBy: repair.TechnicianID.String(),
		})
	}

	// Status changes
	if d.IsStolen && d.StolenDate != nil {
		audit = append(audit, AuditEntry{
			Timestamp:   *d.StolenDate,
			Action:      "Device Reported Stolen",
			Details:     fmt.Sprintf("Reason: %s", d.StolenReason),
			PerformedBy: d.OwnerID.String(),
		})
	}

	// Last inspection
	if d.LastInspection != nil {
		audit = append(audit, AuditEntry{
			Timestamp:   *d.LastInspection,
			Action:      "Device Inspected",
			Details:     d.InspectionNotes,
			PerformedBy: "Inspector",
		})
	}

	// Subscriptions
	for _, sub := range d.Subscriptions {
		audit = append(audit, AuditEntry{
			Timestamp:   sub.StartDate,
			Action:      "Subscription Started",
			Details:     fmt.Sprintf("%s plan activated", sub.PlanType),
			PerformedBy: d.OwnerID.String(),
		})
	}

	// Financing
	for _, finance := range d.Financings {
		audit = append(audit, AuditEntry{
			Timestamp:   finance.CreatedAt,
			Action:      "Financing Initiated",
			Details:     fmt.Sprintf("Financed %.2f for %d months", finance.FinancedAmount, finance.TermMonths),
			PerformedBy: d.OwnerID.String(),
		})
	}

	// Sort by timestamp (newest first)
	// In production, implement proper sorting

	return audit
}

// ExportToJSON exports device data as JSON
func (d *Device) ExportToJSON() ([]byte, error) {
	// Create export structure
	export := map[string]interface{}{
		"device_info": map[string]interface{}{
			"id":               d.ID,
			"imei":             d.IMEI,
			"serial_number":    d.SerialNumber,
			"brand":            d.Brand,
			"model":            d.Model,
			"manufacturer":     d.Manufacturer,
			"operating_system": d.OperatingSystem,
			"os_version":       d.OSVersion,
		},
		"specifications": map[string]interface{}{
			"storage_capacity": d.StorageCapacity,
			"ram":              d.RAM,
			"color":            d.Color,
			"screen_size":      d.ScreenSize,
			"screen_type":      d.ScreenType,
			"battery_capacity": d.BatteryCapacity,
			"water_resistance": d.WaterResistance,
			"5g_capable":       d.Is5GCapable,
			"dual_sim_type":    d.DualSIMType,
		},
		"condition": map[string]interface{}{
			"condition":        d.Condition,
			"grade":            d.Grade,
			"screen_condition": d.ScreenCondition,
			"body_condition":   d.BodyCondition,
			"battery_health":   d.BatteryHealth,
			"water_damage":     d.WaterDamageIndicator,
		},
		"financial": map[string]interface{}{
			"purchase_date":  d.PurchaseDate,
			"purchase_price": d.PurchasePrice,
			"current_value":  d.CurrentValue,
			"market_value":   d.MarketValue,
			"currency":       d.Currency,
		},
		"insurance": map[string]interface{}{
			"risk_score":       d.RiskScore,
			"theft_risk_level": d.TheftRiskLevel,
			"fraud_risk_score": d.FraudRiskScore,
			"monthly_premium":  d.CalculateInsurancePremium(),
			"coverage_limit":   d.GetCoverageLimit(),
			"insurable":        d.CanBeInsured(),
		},
		"status": map[string]interface{}{
			"status":              d.Status,
			"is_stolen":           d.IsStolen,
			"is_locked":           d.IsLocked,
			"is_verified":         d.IsVerified,
			"blacklist_status":    d.BlacklistStatus,
			"authenticity_status": d.AuthenticityStatus,
		},
		"metadata": map[string]interface{}{
			"registration_date": d.RegistrationDate,
			"last_verified_at":  d.LastVerifiedAt,
			"warranty_expiry":   d.WarrantyExpiry,
			"last_inspection":   d.LastInspection,
			"export_date":       time.Now(),
		},
	}

	// Convert to JSON
	return json.MarshalIndent(export, "", "  ")
}

// ImportFromJSON imports device data from JSON
func (d *Device) ImportFromJSON(data []byte) error {
	var imported map[string]interface{}
	if err := json.Unmarshal(data, &imported); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Import device info
	if info, ok := imported["device_info"].(map[string]interface{}); ok {
		if imei, ok := info["imei"].(string); ok {
			d.IMEI = imei
		}
		if serial, ok := info["serial_number"].(string); ok {
			d.SerialNumber = serial
		}
		if brand, ok := info["brand"].(string); ok {
			d.Brand = brand
		}
		if model, ok := info["model"].(string); ok {
			d.Model = model
		}
	}

	// Import specifications
	if specs, ok := imported["specifications"].(map[string]interface{}); ok {
		if storage, ok := specs["storage_capacity"].(float64); ok {
			d.StorageCapacity = int(storage)
		}
		if ram, ok := specs["ram"].(float64); ok {
			d.RAM = int(ram)
		}
	}

	// Import condition
	if condition, ok := imported["condition"].(map[string]interface{}); ok {
		if cond, ok := condition["condition"].(string); ok {
			d.Condition = cond
		}
		if grade, ok := condition["grade"].(string); ok {
			d.Grade = grade
		}
	}

	// Validate imported data
	return d.Validate()
}
