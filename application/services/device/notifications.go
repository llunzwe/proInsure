package device

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Alert represents a device alert
type Alert struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Priority     string    `json:"priority"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
	Acknowledged bool      `json:"acknowledged"`
}

// Reminder represents a device-related reminder
type Reminder struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	DueDate   time.Time `json:"due_date"`
	Message   string    `json:"message"`
	Recurring bool      `json:"recurring"`
	Frequency string    `json:"frequency"` // daily, weekly, monthly, yearly
}

// GetPendingAlerts returns active alerts for device
func (d *Device) GetPendingAlerts() []Alert {
	alerts := []Alert{}
	now := time.Now()

	// Critical alerts

	// Stolen device alert
	if d.IsStolen {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-STOLEN-%s", d.ID),
			Type:         "security",
			Priority:     "critical",
			Message:      "Device reported as stolen - immediate action required",
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Blacklist alert
	if d.BlacklistStatus == "blocked" {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-BLACKLIST-%s", d.ID),
			Type:         "security",
			Priority:     "critical",
			Message:      "Device is blacklisted - service restricted",
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// High priority alerts

	// Battery critical
	if d.BatteryHealth > 0 && d.BatteryHealth < 50 {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-BATTERY-%s", d.ID),
			Type:         "maintenance",
			Priority:     "high",
			Message:      fmt.Sprintf("Battery health critical at %d%% - replacement needed", d.BatteryHealth),
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Water damage detected
	if d.WaterDamageIndicator == "red" || d.WaterDamageIndicator == "pink" {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-WATER-%s", d.ID),
			Type:         "damage",
			Priority:     "high",
			Message:      "Water damage detected - device may malfunction",
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Fraud risk alert
	if d.FraudRiskScore > 75 {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-FRAUD-%s", d.ID),
			Type:         "security",
			Priority:     "high",
			Message:      fmt.Sprintf("High fraud risk detected (%.0f/100)", d.FraudRiskScore),
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Medium priority alerts

	// Warranty expiring
	if d.WarrantyExpiry != nil && time.Until(*d.WarrantyExpiry) < 30*24*time.Hour && time.Until(*d.WarrantyExpiry) > 0 {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-WARRANTY-%s", d.ID),
			Type:         "expiry",
			Priority:     "medium",
			Message:      fmt.Sprintf("Warranty expiring in %d days", int(time.Until(*d.WarrantyExpiry).Hours()/24)),
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Screen damage
	if d.ScreenCondition == "cracked" || d.ScreenCondition == "broken" {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-SCREEN-%s", d.ID),
			Type:         "damage",
			Priority:     "medium",
			Message:      fmt.Sprintf("Screen %s - repair recommended", d.ScreenCondition),
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Inspection overdue
	if d.RequiresInspection() {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-INSPECTION-%s", d.ID),
			Type:         "maintenance",
			Priority:     "medium",
			Message:      "Device inspection overdue",
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Low priority alerts

	// Device age alert
	if d.GetDeviceAge() > 1095 { // 3 years
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-AGE-%s", d.ID),
			Type:         "information",
			Priority:     "low",
			Message:      "Device is over 3 years old - consider upgrade options",
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Security features not enabled
	if !d.FindMyDeviceEnabled || !d.RemoteLockEnabled {
		alerts = append(alerts, Alert{
			ID:           fmt.Sprintf("ALERT-SECURITY-%s", d.ID),
			Type:         "security",
			Priority:     "low",
			Message:      "Security features not fully enabled",
			Timestamp:    now,
			Acknowledged: false,
		})
	}

	// Active claim alerts
	for _, claim := range d.Claims {
		if claim.Status == "pending" || claim.Status == "processing" {
			alerts = append(alerts, Alert{
				ID:           fmt.Sprintf("ALERT-CLAIM-%s", claim.ID),
				Type:         "claim",
				Priority:     "medium",
				Message:      fmt.Sprintf("Claim %s is %s", claim.ID.String()[:8], claim.Status),
				Timestamp:    now,
				Acknowledged: false,
			})
		}
	}

	// Active repair alerts
	for _, repair := range d.Repairs {
		if repair.RepairStatus == "in_progress" {
			alerts = append(alerts, Alert{
				ID:           fmt.Sprintf("ALERT-REPAIR-%s", repair.ID),
				Type:         "repair",
				Priority:     "medium",
				Message:      fmt.Sprintf("Repair in progress - estimated completion: %s", repair.EstimatedCompletion.Format("Jan 2")),
				Timestamp:    now,
				Acknowledged: false,
			})
		}
	}

	return alerts
}

// ShouldNotifyForMaintenance checks if maintenance reminder needed
func (d *Device) ShouldNotifyForMaintenance() bool {
	// Check if inspection is due
	if d.LastInspection == nil {
		return true // Never inspected
	}

	daysSinceInspection := int(time.Since(*d.LastInspection).Hours() / 24)

	// Annual inspection recommended
	if daysSinceInspection > 365 {
		return true
	}

	// Check battery health
	if d.BatteryHealth > 0 && d.BatteryHealth < 70 {
		return true
	}

	// Check for damage
	if d.ScreenCondition == "cracked" || d.ScreenCondition == "broken" {
		return true
	}

	// Check charging port issues
	if d.ChargingPortCondition == "intermittent" || d.ChargingPortCondition == "damaged" {
		return true
	}

	// Check device age
	if d.GetDeviceAge() > 730 && daysSinceInspection > 180 { // 2+ year old device, 6 months since inspection
		return true
	}

	// Check performance
	if d.GetPerformanceScore() < 60 {
		return true
	}

	return false
}

// GetRenewalReminders returns renewal/expiry reminders
func (d *Device) GetRenewalReminders() []Reminder {
	reminders := []Reminder{}
	now := time.Now()

	// Warranty expiry reminder
	if d.WarrantyExpiry != nil {
		daysUntilExpiry := int(time.Until(*d.WarrantyExpiry).Hours() / 24)

		if daysUntilExpiry > 0 && daysUntilExpiry <= 90 {
			reminders = append(reminders, Reminder{
				ID:        fmt.Sprintf("REM-WARRANTY-%s", d.ID),
				Type:      "warranty_expiry",
				DueDate:   *d.WarrantyExpiry,
				Message:   fmt.Sprintf("Warranty expires in %d days - consider extended warranty", daysUntilExpiry),
				Recurring: false,
				Frequency: "",
			})
		}
	}

	// Insurance renewal reminders
	for _, policy := range d.Policies {
		if policy.Status == "active" && policy.EndDate != nil {
			daysUntilRenewal := int(time.Until(*policy.EndDate).Hours() / 24)

			if daysUntilRenewal > 0 && daysUntilRenewal <= 30 {
				reminders = append(reminders, Reminder{
					ID:        fmt.Sprintf("REM-POLICY-%s", policy.ID),
					Type:      "insurance_renewal",
					DueDate:   *policy.EndDate,
					Message:   fmt.Sprintf("Insurance policy expires in %d days", daysUntilRenewal),
					Recurring: false,
					Frequency: "",
				})
			}
		}
	}

	// Subscription renewals
	for _, subscription := range d.Subscriptions {
		if subscription.IsActive() && subscription.NextBillingDate != nil {
			daysUntilBilling := int(time.Until(*subscription.NextBillingDate).Hours() / 24)

			if daysUntilBilling > 0 && daysUntilBilling <= 7 {
				reminders = append(reminders, Reminder{
					ID:        fmt.Sprintf("REM-SUB-%s", subscription.ID),
					Type:      "subscription_renewal",
					DueDate:   *subscription.NextBillingDate,
					Message:   fmt.Sprintf("%s subscription renews in %d days", subscription.PlanType, daysUntilBilling),
					Recurring: true,
					Frequency: "monthly",
				})
			}
		}
	}

	// Financing payment reminders
	for _, financing := range d.Financings {
		if financing.FinanceStatus == "active" && financing.NextPaymentDate != nil {
			daysUntilPayment := int(time.Until(*financing.NextPaymentDate).Hours() / 24)

			if daysUntilPayment > 0 && daysUntilPayment <= 5 {
				reminders = append(reminders, Reminder{
					ID:        fmt.Sprintf("REM-FINANCE-%s", financing.ID),
					Type:      "payment_due",
					DueDate:   *financing.NextPaymentDate,
					Message:   fmt.Sprintf("Finance payment of %.2f due in %d days", financing.MonthlyPayment, daysUntilPayment),
					Recurring: true,
					Frequency: "monthly",
				})
			}
		}
	}

	// Rental return reminders
	for _, rental := range d.Rentals {
		if rental.RentalStatus == "active" && rental.EndDate != nil {
			daysUntilReturn := int(time.Until(*rental.EndDate).Hours() / 24)

			if daysUntilReturn > 0 && daysUntilReturn <= 7 {
				reminders = append(reminders, Reminder{
					ID:        fmt.Sprintf("REM-RENTAL-%s", rental.ID),
					Type:      "rental_return",
					DueDate:   *rental.EndDate,
					Message:   fmt.Sprintf("Rental period ends in %d days", daysUntilReturn),
					Recurring: false,
					Frequency: "",
				})
			}
		}
	}

	// Maintenance reminders
	if d.ShouldNotifyForMaintenance() {
		maintenanceDue := now.AddDate(0, 0, 7) // Due in 7 days
		if d.LastInspection != nil && time.Since(*d.LastInspection) > 365*24*time.Hour {
			maintenanceDue = now // Overdue
		}

		reminders = append(reminders, Reminder{
			ID:        fmt.Sprintf("REM-MAINT-%s", d.ID),
			Type:      "maintenance",
			DueDate:   maintenanceDue,
			Message:   "Device maintenance recommended",
			Recurring: true,
			Frequency: "yearly",
		})
	}

	// Trade-in opportunity reminder
	if d.IsEligibleForTradeIn() && d.GetDeviceAge() > 730 { // 2+ years old
		reminders = append(reminders, Reminder{
			ID:        fmt.Sprintf("REM-TRADEIN-%s", d.ID),
			Type:      "opportunity",
			DueDate:   now.AddDate(0, 1, 0), // Check monthly
			Message:   fmt.Sprintf("Trade-in value: $%.2f - Consider upgrading", d.CalculateTradeInValue()),
			Recurring: true,
			Frequency: "monthly",
		})
	}

	// Battery replacement reminder
	batteryReplacement := d.PredictBatteryReplacement()
	if batteryReplacement != nil && time.Until(*batteryReplacement) < 90*24*time.Hour {
		reminders = append(reminders, Reminder{
			ID:        fmt.Sprintf("REM-BATTERY-%s", d.ID),
			Type:      "maintenance",
			DueDate:   *batteryReplacement,
			Message:   "Battery replacement recommended soon",
			Recurring: false,
			Frequency: "",
		})
	}

	return reminders
}

// TriggerSecurityAlert triggers security notification
func (d *Device) TriggerSecurityAlert(reason string) error {
	// Validate reason
	if reason == "" {
		return fmt.Errorf("security alert reason cannot be empty")
	}

	// Determine alert severity
	severity := "medium"
	securityActions := []string{}

	reasonLower := strings.ToLower(reason)

	// Check for critical security issues
	if strings.Contains(reasonLower, "stolen") || strings.Contains(reasonLower, "lost") {
		severity = "critical"
		d.IsStolen = true
		now := time.Now()
		d.StolenDate = &now
		d.StolenReason = reason
		securityActions = append(securityActions, "Device marked as stolen")

		// Enable security features if available
		if !d.RemoteLockEnabled {
			securityActions = append(securityActions, "Remote lock recommended")
		}
		if !d.FindMyDeviceEnabled {
			securityActions = append(securityActions, "Find My Device should be enabled")
		}
	}

	// Check for fraud
	if strings.Contains(reasonLower, "fraud") || strings.Contains(reasonLower, "suspicious") {
		severity = "high"
		d.FraudRiskScore = math.Min(d.FraudRiskScore+20, 100)
		securityActions = append(securityActions, "Fraud risk score increased")
	}

	// Check for unauthorized access
	if strings.Contains(reasonLower, "unauthorized") || strings.Contains(reasonLower, "breach") {
		severity = "high"
		securityActions = append(securityActions, "Change all passwords immediately")
		securityActions = append(securityActions, "Enable two-factor authentication")
	}

	// Check for blacklist
	if strings.Contains(reasonLower, "blacklist") || strings.Contains(reasonLower, "blocked") {
		severity = "critical"
		d.BlacklistStatus = "blocked"
		securityActions = append(securityActions, "Device blacklisted")
	}

	// Create security alert record (in production, this would be saved to database)
	alert := map[string]interface{}{
		"device_id":        d.ID,
		"alert_type":       "security",
		"severity":         severity,
		"reason":           reason,
		"triggered_at":     time.Now(),
		"device_status":    d.Status,
		"security_actions": securityActions,
		"imei":             d.IMEI,
		"location":         d.LastKnownLocation,
	}

	// Log the alert
	fmt.Printf("Security Alert Triggered: %+v\n", alert)

	// Update device status if critical
	if severity == "critical" {
		d.Status = "compromised"
	}

	// Update last known location timestamp
	if d.LastKnownLocation != "" {
		now := time.Now()
		d.LastLocationUpdate = &now
	}

	return nil
}

// GetNotificationPreferences returns notification settings for device
func (d *Device) GetNotificationPreferences() map[string]interface{} {
	// Default preferences (in production, these would be user-configurable)
	preferences := map[string]interface{}{
		"security_alerts": map[string]bool{
			"theft_alerts":      true,
			"fraud_alerts":      true,
			"blacklist_alerts":  true,
			"location_tracking": d.FindMyDeviceEnabled,
		},
		"maintenance_alerts": map[string]bool{
			"inspection_reminders": true,
			"battery_alerts":       true,
			"damage_alerts":        true,
			"warranty_expiry":      true,
		},
		"service_alerts": map[string]bool{
			"claim_updates":         true,
			"repair_updates":        true,
			"subscription_renewals": true,
			"payment_reminders":     true,
		},
		"marketing_alerts": map[string]bool{
			"upgrade_offers":     false,
			"trade_in_offers":    true,
			"new_features":       false,
			"promotional_offers": false,
		},
		"notification_channels": map[string]bool{
			"email":  true,
			"sms":    true,
			"push":   true,
			"in_app": true,
		},
		"quiet_hours": map[string]interface{}{
			"enabled":    false,
			"start_time": "22:00",
			"end_time":   "08:00",
		},
	}

	return preferences
}
