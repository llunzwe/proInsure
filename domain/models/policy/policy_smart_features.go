package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicySmartFeatures represents IoT and smart monitoring features for a policy
type PolicySmartFeatures struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`

	// Smart Monitoring
	SmartMonitoring    bool   `gorm:"default:false" json:"smart_monitoring"`
	MonitoringLevel    string `gorm:"type:varchar(20)" json:"monitoring_level"` // basic, advanced, premium
	DeviceHealthAlerts bool   `gorm:"default:false" json:"device_health_alerts"`
	AlertThreshold     string `gorm:"type:json" json:"alert_threshold"` // JSON with thresholds

	// Preventive Maintenance
	PreventiveMaintenance bool       `gorm:"default:false" json:"preventive_maintenance"`
	MaintenanceSchedule   string     `gorm:"type:json" json:"maintenance_schedule"`
	LastMaintenanceDate   *time.Time `json:"last_maintenance_date"`
	NextMaintenanceDate   *time.Time `json:"next_maintenance_date"`
	MaintenanceReminders  bool       `gorm:"default:true" json:"maintenance_reminders"`

	// Auto Detection
	AutoClaimDetection bool    `gorm:"default:false" json:"auto_claim_detection"`
	AutoClaimThreshold float64 `json:"auto_claim_threshold"` // Damage percentage
	AutoClaimApproval  bool    `gorm:"default:false" json:"auto_claim_approval"`
	DetectionAccuracy  float64 `json:"detection_accuracy"` // Historical accuracy

	// Remote Diagnostics
	RemoteDiagnostics   bool       `gorm:"default:false" json:"remote_diagnostics"`
	DiagnosticFrequency string     `gorm:"type:varchar(20)" json:"diagnostic_frequency"` // daily, weekly, monthly
	LastDiagnosticDate  *time.Time `json:"last_diagnostic_date"`
	DiagnosticHistory   string     `gorm:"type:json" json:"diagnostic_history"`

	// Usage Analytics
	UsageBasedPricing  bool    `gorm:"default:false" json:"usage_based_pricing"`
	MonthlyUsageHours  float64 `json:"monthly_usage_hours"`
	AverageScreenTime  float64 `json:"average_screen_time"` // Hours per day
	AppUsageMonitoring bool    `gorm:"default:false" json:"app_usage_monitoring"`
	LocationTracking   bool    `gorm:"default:false" json:"location_tracking"`

	// Safety & Security
	SafetyScore       float64 `json:"safety_score"`
	SafetyFactors     string  `gorm:"type:json" json:"safety_factors"` // JSON with factors
	SecurityScore     float64 `json:"security_score"`
	AntiTheftEnabled  bool    `gorm:"default:false" json:"anti_theft_enabled"`
	RemoteLockEnabled bool    `gorm:"default:false" json:"remote_lock_enabled"`
	RemoteWipeEnabled bool    `gorm:"default:false" json:"remote_wipe_enabled"`

	// AI Features
	AIAssistant        bool       `gorm:"default:false" json:"ai_assistant"`
	PredictiveAnalysis bool       `gorm:"default:false" json:"predictive_analysis"`
	RiskPrediction     float64    `json:"risk_prediction"`
	FailurePrediction  *time.Time `json:"failure_prediction"`

	// Data Collection
	DataCollectionConsent bool `gorm:"default:false" json:"data_collection_consent"`
	DataSharingConsent    bool `gorm:"default:false" json:"data_sharing_consent"`
	AnonymizedData        bool `gorm:"default:true" json:"anonymized_data"`
	DataRetentionDays     int  `gorm:"default:90" json:"data_retention_days"`

	// Status
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	LastSyncDate     *time.Time `json:"last_sync_date"`
	ConnectionStatus string     `gorm:"type:varchar(20)" json:"connection_status"` // connected, disconnected, intermittent

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicySmartFeatures) TableName() string {
	return "policy_smart_features"
}

// HasSmartFeatures checks if smart features are enabled
func (psf *PolicySmartFeatures) HasSmartFeatures() bool {
	return psf.IsActive && (psf.SmartMonitoring ||
		psf.DeviceHealthAlerts ||
		psf.PreventiveMaintenance ||
		psf.AutoClaimDetection ||
		psf.RemoteDiagnostics)
}

// IsEligibleForPreventiveMaintenance checks eligibility
func (psf *PolicySmartFeatures) IsEligibleForPreventiveMaintenance() bool {
	if !psf.IsActive || !psf.PreventiveMaintenance {
		return false
	}

	if psf.NextMaintenanceDate != nil {
		return time.Now().After(*psf.NextMaintenanceDate)
	}

	return true
}

// GetSafetyScoreDiscount returns discount based on safety score
func (psf *PolicySmartFeatures) GetSafetyScoreDiscount() float64 {
	if !psf.UsageBasedPricing {
		return 0
	}

	if psf.SafetyScore >= 90 {
		return 0.15 // 15% discount for excellent safety
	} else if psf.SafetyScore >= 75 {
		return 0.10 // 10% discount for good safety
	} else if psf.SafetyScore >= 60 {
		return 0.05 // 5% discount for average safety
	}

	return 0
}

// UpdateSafetyScore updates the safety score based on usage patterns
func (psf *PolicySmartFeatures) UpdateSafetyScore(factors map[string]float64) {
	score := 50.0 // Base score

	// Factor in various safety indicators
	if dropRate, ok := factors["drop_rate"]; ok {
		score -= dropRate * 10
	}
	if chargePattern, ok := factors["charge_pattern"]; ok {
		score += chargePattern * 5
	}
	if appSecurity, ok := factors["app_security"]; ok {
		score += appSecurity * 10
	}
	if locationRisk, ok := factors["location_risk"]; ok {
		score -= locationRisk * 5
	}

	// Ensure score is within bounds
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	psf.SafetyScore = score
}

// ShouldSendMaintenanceReminder checks if maintenance reminder should be sent
func (psf *PolicySmartFeatures) ShouldSendMaintenanceReminder() bool {
	if !psf.IsActive || !psf.PreventiveMaintenance || !psf.MaintenanceReminders {
		return false
	}

	if psf.NextMaintenanceDate == nil {
		return false
	}

	daysUntilMaintenance := int(time.Until(*psf.NextMaintenanceDate).Hours() / 24)
	return daysUntilMaintenance <= 7 && daysUntilMaintenance > 0
}

// CanAutoApproveClaim checks if claim can be auto-approved
func (psf *PolicySmartFeatures) CanAutoApproveClaim(damageScore float64) bool {
	return psf.IsActive &&
		psf.AutoClaimDetection &&
		psf.AutoClaimApproval &&
		damageScore >= psf.AutoClaimThreshold &&
		psf.DetectionAccuracy >= 0.85 // 85% accuracy threshold
}

// IsConnected checks if device is connected
func (psf *PolicySmartFeatures) IsConnected() bool {
	return psf.ConnectionStatus == "connected"
}

// GetMonitoringIntensity returns monitoring intensity level
func (psf *PolicySmartFeatures) GetMonitoringIntensity() string {
	if !psf.SmartMonitoring {
		return "none"
	}

	activeFeatures := 0
	if psf.DeviceHealthAlerts {
		activeFeatures++
	}
	if psf.PreventiveMaintenance {
		activeFeatures++
	}
	if psf.AutoClaimDetection {
		activeFeatures++
	}
	if psf.RemoteDiagnostics {
		activeFeatures++
	}
	if psf.LocationTracking {
		activeFeatures++
	}

	if activeFeatures >= 4 {
		return "high"
	} else if activeFeatures >= 2 {
		return "medium"
	}
	return "low"
}

// HasSecurityFeatures checks if security features are enabled
func (psf *PolicySmartFeatures) HasSecurityFeatures() bool {
	return psf.AntiTheftEnabled ||
		psf.RemoteLockEnabled ||
		psf.RemoteWipeEnabled
}

// GetDataPrivacyLevel returns data privacy level
func (psf *PolicySmartFeatures) GetDataPrivacyLevel() string {
	if !psf.DataCollectionConsent {
		return "no_collection"
	}

	if psf.AnonymizedData && !psf.DataSharingConsent {
		return "high"
	} else if psf.AnonymizedData && psf.DataSharingConsent {
		return "medium"
	}
	return "low"
}
