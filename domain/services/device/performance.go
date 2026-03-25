package device

import (
	"fmt"
	"math"
	"time"

	"smartsure/internal/domain/models"
)

// MaintenanceTask represents a maintenance recommendation
type MaintenanceTask struct {
	Task        string    `json:"task"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	Description string    `json:"description"`
}

// ComponentStatus represents status of a hardware component
type ComponentStatus struct {
	Component   string    `json:"component"`
	Status      string    `json:"status"`
	Health      float64   `json:"health"`
	LastChecked time.Time `json:"last_checked"`
}

// GetPerformanceScore calculates overall performance rating
func (d *models.Device) GetPerformanceScore() float64 {
	score := 100.0

	// Battery performance impact
	if d.BatteryHealth > 0 {
		batteryImpact := (100 - float64(d.BatteryHealth)) * 0.3 // Max 30% impact
		score -= batteryImpact
	}

	// Age-based degradation
	ageYears := float64(d.GetDeviceAge()) / 365
	ageDegradation := ageYears * 5 // 5% per year
	if ageDegradation > 30 {
		ageDegradation = 30 // Cap at 30%
	}
	score -= ageDegradation

	// Storage impact (assumed based on age and segment)
	storageImpact := 0.0
	if d.StorageCapacity > 0 {
		if d.StorageCapacity <= 32 {
			storageImpact = 10 // Low storage devices degrade faster
		} else if d.StorageCapacity <= 64 {
			storageImpact = 5
		}
	}
	score -= storageImpact

	// RAM impact
	ramImpact := 0.0
	if d.RAM > 0 {
		if d.RAM <= 2 {
			ramImpact = 15 // Low RAM significantly impacts performance
		} else if d.RAM <= 4 {
			ramImpact = 8
		} else if d.RAM <= 6 {
			ramImpact = 3
		}
	}
	score -= ramImpact

	// Component failures impact
	componentImpacts := map[string]float64{
		"broken":       15,
		"damaged":      10,
		"intermittent": 5,
		"distorted":    5,
		"not_working":  15,
	}

	// Check various components
	if impact, exists := componentImpacts[d.ChargingPortCondition]; exists {
		score -= impact * 0.5 // Charging port has less impact on performance
	}
	if impact, exists := componentImpacts[d.SpeakerCondition]; exists {
		score -= impact * 0.3
	}
	if impact, exists := componentImpacts[d.MicrophoneCondition]; exists {
		score -= impact * 0.3
	}

	// Screen responsiveness impact
	if !d.TouchScreenResponsive {
		score -= 20
	}
	if d.DeadPixels {
		score -= 5
	}

	// Software/OS impact
	if d.OSVersion != "" {
		// Older OS versions might have performance issues
		if ageYears > 3 {
			score -= 5 // Outdated OS impact
		}
	}

	// Ensure score is between 0 and 100
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	return score
}

// PredictBatteryReplacement predicts when battery needs replacement
func (d *models.Device) PredictBatteryReplacement() *time.Time {
	if d.BatteryHealth <= 0 {
		// No battery health data
		return nil
	}

	// Battery is already critical
	if d.BatteryHealth < 60 {
		now := time.Now()
		return &now
	}

	// Calculate degradation rate
	// Typical battery degrades ~20% per year with normal use
	ageYears := float64(d.GetDeviceAge()) / 365
	if ageYears < 0.5 {
		ageYears = 0.5 // Minimum for calculation
	}

	currentDegradation := 100 - float64(d.BatteryHealth)
	degradationRate := currentDegradation / ageYears // Percent per year

	// Adjust degradation rate based on battery cycles
	if d.BatteryCycles > 0 {
		// Typical battery rated for 500-1000 cycles
		cycleBasedRate := float64(d.BatteryCycles) / 500 * 20 // 20% per 500 cycles
		// Use the higher degradation rate
		if cycleBasedRate > degradationRate {
			degradationRate = cycleBasedRate
		}
	}

	// Calculate when battery will reach 60% health (replacement threshold)
	remainingHealth := float64(d.BatteryHealth) - 60
	if remainingHealth <= 0 {
		now := time.Now()
		return &now
	}

	yearsToReplacement := remainingHealth / degradationRate
	replacementDate := time.Now().AddDate(0, 0, int(yearsToReplacement*365))

	return &replacementDate
}

// GetDegradationRate calculates component degradation rate
func (d *models.Device) GetDegradationRate() float64 {
	// Base degradation rate
	baseRate := 10.0 // 10% per year baseline

	// Adjust based on device segment
	segmentMultipliers := map[string]float64{
		"flagship":  0.8, // Better components
		"premium":   0.9,
		"mid_range": 1.0,
		"budget":    1.2, // Faster degradation
	}

	if multiplier, exists := segmentMultipliers[d.DeviceSegment]; exists {
		baseRate *= multiplier
	}

	// Adjust based on usage patterns (simulated)
	// Heavy repairs indicate heavy usage
	repairCount := float64(d.GetRepairCount())
	if repairCount > 3 {
		baseRate *= 1.5
	} else if repairCount > 1 {
		baseRate *= 1.2
	}

	// Adjust based on condition
	conditionMultipliers := map[string]float64{
		"poor":      2.0,
		"fair":      1.5,
		"good":      1.0,
		"excellent": 0.8,
	}

	if multiplier, exists := conditionMultipliers[d.Condition]; exists {
		baseRate *= multiplier
	}

	// Environmental factors
	if d.WaterDamageIndicator != "white" {
		baseRate *= 2.0 // Water damage accelerates degradation
	}

	// Battery health impacts overall degradation
	if d.BatteryHealth > 0 && d.BatteryHealth < 70 {
		baseRate *= 1.3
	}

	return baseRate
}

// CheckHardwareIntegrity verifies hardware components
func (d *models.Device) CheckHardwareIntegrity() error {
	issues := []string{}

	// Check critical components
	if d.ScreenCondition == "broken" {
		issues = append(issues, "Screen broken")
	}

	if !d.TouchScreenResponsive {
		issues = append(issues, "Touch screen not responsive")
	}

	if d.ChargingPortCondition == "damaged" || d.ChargingPortCondition == "not_working" {
		issues = append(issues, "Charging port damaged")
	}

	if d.BatteryHealth > 0 && d.BatteryHealth < 50 {
		issues = append(issues, "Critical battery degradation")
	}

	if d.CameraCondition == "both_issue" {
		issues = append(issues, "Both cameras malfunctioning")
	}

	if d.SpeakerCondition == "not_working" {
		issues = append(issues, "Speaker not working")
	}

	if d.MicrophoneCondition == "not_working" {
		issues = append(issues, "Microphone not working")
	}

	if d.WaterDamageIndicator == "red" || d.WaterDamageIndicator == "pink" {
		issues = append(issues, "Water damage detected")
	}

	if len(issues) > 0 {
		return fmt.Errorf("hardware integrity issues: %v", issues)
	}

	return nil
}

// GetOptimalMaintenanceSchedule returns maintenance recommendations
func (d *models.Device) GetOptimalMaintenanceSchedule() []MaintenanceTask {
	tasks := []MaintenanceTask{}
	now := time.Now()

	// Battery maintenance
	if d.BatteryHealth > 0 {
		if d.BatteryHealth < 60 {
			tasks = append(tasks, MaintenanceTask{
				Task:        "Battery Replacement",
				Priority:    "High",
				DueDate:     now,
				Description: "Battery health critical, immediate replacement recommended",
			})
		} else if d.BatteryHealth < 70 {
			tasks = append(tasks, MaintenanceTask{
				Task:        "Battery Check",
				Priority:    "Medium",
				DueDate:     now.AddDate(0, 1, 0),
				Description: "Battery showing degradation, monitor closely",
			})
		}
	}

	// Screen maintenance
	if d.ScreenCondition == "cracked" {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Screen Repair",
			Priority:    "Medium",
			DueDate:     now.AddDate(0, 0, 14),
			Description: "Screen cracked, repair to prevent further damage",
		})
	} else if d.ScreenCondition == "minor_scratches" {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Screen Protection",
			Priority:    "Low",
			DueDate:     now.AddDate(0, 1, 0),
			Description: "Apply screen protector to prevent further scratches",
		})
	}

	// Charging port maintenance
	if d.ChargingPortCondition == "intermittent" {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Charging Port Cleaning",
			Priority:    "Medium",
			DueDate:     now.AddDate(0, 0, 7),
			Description: "Clean charging port to improve connection",
		})
	}

	// Software maintenance
	if d.GetDeviceAge() > 180 && (d.LastInspection == nil || time.Since(*d.LastInspection) > 180*24*time.Hour) {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Software Update Check",
			Priority:    "Medium",
			DueDate:     now.AddDate(0, 0, 7),
			Description: "Check for OS and security updates",
		})
	}

	// General inspection
	if d.LastInspection == nil || time.Since(*d.LastInspection) > 365*24*time.Hour {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Annual Inspection",
			Priority:    "Medium",
			DueDate:     now.AddDate(0, 1, 0),
			Description: "Comprehensive device inspection recommended",
		})
	}

	// Storage cleanup for low storage devices
	if d.StorageCapacity > 0 && d.StorageCapacity <= 32 {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Storage Cleanup",
			Priority:    "Low",
			DueDate:     now.AddDate(0, 3, 0),
			Description: "Clean up storage to maintain performance",
		})
	}

	// Water damage check
	if d.WaterResistance == "none" && d.GetDeviceAge() > 365 {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Moisture Check",
			Priority:    "Low",
			DueDate:     now.AddDate(0, 6, 0),
			Description: "Check water damage indicators",
		})
	}

	// Camera maintenance
	if d.CameraCondition == "front_issue" || d.CameraCondition == "rear_issue" {
		tasks = append(tasks, MaintenanceTask{
			Task:        "Camera Inspection",
			Priority:    "Low",
			DueDate:     now.AddDate(0, 2, 0),
			Description: "Inspect and clean camera lenses",
		})
	}

	return tasks
}

// GetComponentHealthStatus returns detailed component health
func (d *models.Device) GetComponentHealthStatus() []ComponentStatus {
	components := []ComponentStatus{}
	now := time.Now()

	// Battery status
	if d.BatteryHealth > 0 {
		status := "healthy"
		health := float64(d.BatteryHealth)
		if d.BatteryHealth < 60 {
			status = "critical"
		} else if d.BatteryHealth < 70 {
			status = "degraded"
		} else if d.BatteryHealth < 80 {
			status = "fair"
		}

		components = append(components, ComponentStatus{
			Component:   "Battery",
			Status:      status,
			Health:      health,
			LastChecked: now,
		})
	}

	// Screen status
	screenHealth := 100.0
	screenStatus := "healthy"
	switch d.ScreenCondition {
	case "broken":
		screenHealth = 0
		screenStatus = "failed"
	case "cracked":
		screenHealth = 40
		screenStatus = "damaged"
	case "minor_scratches":
		screenHealth = 80
		screenStatus = "fair"
	case "perfect":
		screenHealth = 100
		screenStatus = "healthy"
	}

	components = append(components, ComponentStatus{
		Component:   "Screen",
		Status:      screenStatus,
		Health:      screenHealth,
		LastChecked: now,
	})

	// Charging port status
	portHealth := 100.0
	portStatus := "healthy"
	switch d.ChargingPortCondition {
	case "damaged", "not_working":
		portHealth = 0
		portStatus = "failed"
	case "intermittent":
		portHealth = 50
		portStatus = "degraded"
	case "working":
		portHealth = 100
		portStatus = "healthy"
	}

	components = append(components, ComponentStatus{
		Component:   "Charging Port",
		Status:      portStatus,
		Health:      portHealth,
		LastChecked: now,
	})

	// Camera status
	cameraHealth := 100.0
	cameraStatus := "healthy"
	switch d.CameraCondition {
	case "both_issue":
		cameraHealth = 0
		cameraStatus = "failed"
	case "front_issue", "rear_issue":
		cameraHealth = 50
		cameraStatus = "partial_failure"
	case "all_working":
		cameraHealth = 100
		cameraStatus = "healthy"
	}

	components = append(components, ComponentStatus{
		Component:   "Camera System",
		Status:      cameraStatus,
		Health:      cameraHealth,
		LastChecked: now,
	})

	// Speaker status
	speakerHealth := 100.0
	speakerStatus := "healthy"
	switch d.SpeakerCondition {
	case "not_working":
		speakerHealth = 0
		speakerStatus = "failed"
	case "distorted":
		speakerHealth = 50
		speakerStatus = "degraded"
	case "working":
		speakerHealth = 100
		speakerStatus = "healthy"
	}

	components = append(components, ComponentStatus{
		Component:   "Speaker",
		Status:      speakerStatus,
		Health:      speakerHealth,
		LastChecked: now,
	})

	// Microphone status
	micHealth := 100.0
	micStatus := "healthy"
	switch d.MicrophoneCondition {
	case "not_working":
		micHealth = 0
		micStatus = "failed"
	case "intermittent":
		micHealth = 50
		micStatus = "degraded"
	case "working":
		micHealth = 100
		micStatus = "healthy"
	}

	components = append(components, ComponentStatus{
		Component:   "Microphone",
		Status:      micStatus,
		Health:      micHealth,
		LastChecked: now,
	})

	// Water damage indicator
	waterHealth := 100.0
	waterStatus := "healthy"
	switch d.WaterDamageIndicator {
	case "red", "pink":
		waterHealth = 0
		waterStatus = "water_damaged"
	case "white":
		waterHealth = 100
		waterStatus = "healthy"
	}

	components = append(components, ComponentStatus{
		Component:   "Water Damage Indicator",
		Status:      waterStatus,
		Health:      waterHealth,
		LastChecked: now,
	})

	// Touch screen
	touchHealth := 100.0
	touchStatus := "healthy"
	if !d.TouchScreenResponsive {
		touchHealth = 0
		touchStatus = "failed"
	}
	if d.DeadPixels {
		touchHealth = math.Max(touchHealth-20, 0)
		if touchStatus == "healthy" {
			touchStatus = "degraded"
		}
	}

	components = append(components, ComponentStatus{
		Component:   "Touch Screen",
		Status:      touchStatus,
		Health:      touchHealth,
		LastChecked: now,
	})

	return components
}

// CalculatePerformanceDegradation estimates performance loss over time
func (d *models.Device) CalculatePerformanceDegradation() float64 {
	// Start with current performance
	currentPerformance := d.GetPerformanceScore()

	// Estimate original performance (when new)
	originalPerformance := 100.0

	// Calculate degradation
	degradation := originalPerformance - currentPerformance

	// Express as percentage
	degradationPercentage := (degradation / originalPerformance) * 100

	return degradationPercentage
}
