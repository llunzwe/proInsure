package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PolicyIntegrations represents third-party integrations for a policy
type PolicyIntegrations struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// TPA (Third Party Administrator)
	TPAEnabled  bool           `gorm:"type:boolean;default:false" json:"tpa_enabled"`
	TPAName     string         `gorm:"type:varchar(255)" json:"tpa_name"`
	TPAID       string         `gorm:"type:varchar(100)" json:"tpa_id"`
	TPAServices datatypes.JSON `gorm:"type:json" json:"tpa_services"` // []string
	TPAFees     Money          `gorm:"embedded;embeddedPrefix:tpa_fees_" json:"tpa_fees"`

	// MGA (Managing General Agent)
	MGAEnabled    bool           `gorm:"type:boolean;default:false" json:"mga_enabled"`
	MGAName       string         `gorm:"type:varchar(255)" json:"mga_name"`
	MGAID         string         `gorm:"type:varchar(100)" json:"mga_id"`
	MGAAuthority  datatypes.JSON `gorm:"type:json" json:"mga_authority"`
	MGACommission float64        `gorm:"type:decimal(5,2)" json:"mga_commission_rate"`

	// External System References
	ExternalSystemRefs datatypes.JSON `gorm:"type:json" json:"external_system_refs"` // map[string]string
	LegacySystemID     string         `gorm:"type:varchar(100)" json:"legacy_system_id"`
	PartnerSystemID    string         `gorm:"type:varchar(100)" json:"partner_system_id"`

	// API Integration
	APIUsageMetrics datatypes.JSON `gorm:"type:json" json:"api_usage_metrics"` // []APIMetric
	APICallsCount   int            `gorm:"type:int" json:"api_calls_count"`
	APIErrorRate    float64        `gorm:"type:decimal(5,2)" json:"api_error_rate"`
	LastAPICallDate *time.Time     `gorm:"type:timestamp" json:"last_api_call_date,omitempty"`
	APIRateLimit    int            `gorm:"type:int" json:"api_rate_limit"`

	// Webhooks
	WebhookSubscriptions datatypes.JSON `gorm:"type:json" json:"webhook_subscriptions"` // []Webhook
	WebhookEndpoint      string         `gorm:"type:varchar(500)" json:"webhook_endpoint"`
	WebhookSecret        string         `gorm:"type:varchar(255)" json:"webhook_secret"`
	WebhookRetryCount    int            `gorm:"type:int" json:"webhook_retry_count"`
	LastWebhookDate      *time.Time     `gorm:"type:timestamp" json:"last_webhook_date,omitempty"`

	// Data Synchronization
	DataSyncStatus datatypes.JSON `gorm:"type:json" json:"data_sync_status"` // map[string]SyncStatus
	LastSyncDate   time.Time      `gorm:"type:timestamp" json:"last_sync_date"`
	NextSyncDate   time.Time      `gorm:"type:timestamp" json:"next_sync_date"`
	SyncFrequency  string         `gorm:"type:varchar(50)" json:"sync_frequency"` // realtime, hourly, daily
	SyncErrors     datatypes.JSON `gorm:"type:json" json:"sync_errors"`

	// Integration Logs
	IntegrationLogs datatypes.JSON `gorm:"type:json" json:"integration_logs"` // []IntegrationLog
	ErrorLogs       datatypes.JSON `gorm:"type:json" json:"error_logs"`
	SuccessRate     float64        `gorm:"type:decimal(5,2)" json:"success_rate"`

	// Status
	IntegrationStatus   string    `gorm:"type:varchar(50)" json:"integration_status"`
	HealthCheckStatus   string    `gorm:"type:varchar(50)" json:"health_check_status"`
	LastHealthCheckDate time.Time `gorm:"type:timestamp" json:"last_health_check_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// PolicyTelematics represents telematics and IoT data for usage-based insurance
type PolicyTelematics struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Connected Devices
	ConnectedDevices datatypes.JSON `gorm:"type:json" json:"connected_devices"` // []IoTDevice
	DeviceCount      int            `gorm:"type:int" json:"device_count"`
	PrimaryDeviceID  string         `gorm:"type:varchar(100)" json:"primary_device_id"`
	DeviceType       string         `gorm:"type:varchar(50)" json:"device_type"` // smartphone, obd, blackbox

	// Usage-Based Insurance (UBI) Metrics
	UsageBasedPricing bool    `gorm:"type:boolean;default:false" json:"usage_based_pricing"`
	MonthlyUsageHours float64 `gorm:"type:decimal(10,2)" json:"monthly_usage_hours"`
	TotalDistanceKm   float64 `gorm:"type:decimal(10,2)" json:"total_distance_km"`
	AverageSpeedKmh   float64 `gorm:"type:decimal(10,2)" json:"average_speed_kmh"`
	NightUsagePercent float64 `gorm:"type:decimal(5,2)" json:"night_usage_percent"`

	// Driving/Usage Behavior (for auto/device insurance)
	DrivingScore           float64 `gorm:"type:decimal(5,2)" json:"driving_score"`
	AccelerationEvents     int     `gorm:"type:int" json:"harsh_acceleration_events"`
	BrakingEvents          int     `gorm:"type:int" json:"harsh_braking_events"`
	CorneringEvents        int     `gorm:"type:int" json:"harsh_cornering_events"`
	SpeedingEvents         int     `gorm:"type:int" json:"speeding_events"`
	PhoneUsageWhileDriving int     `gorm:"type:int" json:"phone_usage_while_driving"`

	// Safety Scores
	SafetyScore         float64 `gorm:"type:decimal(5,2)" json:"safety_score"`
	RiskScore           float64 `gorm:"type:decimal(5,2)" json:"risk_score"`
	ComparisonToAverage float64 `gorm:"type:decimal(5,2)" json:"comparison_to_average"`
	SafetyImprovement   float64 `gorm:"type:decimal(5,2)" json:"safety_improvement"`

	// Real-Time Monitoring
	RealTimeTracking   bool       `gorm:"type:boolean;default:false" json:"real_time_tracking"`
	LastLocationLat    float64    `gorm:"type:decimal(10,8)" json:"last_location_lat"`
	LastLocationLng    float64    `gorm:"type:decimal(11,8)" json:"last_location_lng"`
	LastLocationTime   *time.Time `gorm:"type:timestamp" json:"last_location_time,omitempty"`
	GeofenceViolations int        `gorm:"type:int" json:"geofence_violations"`

	// Alerts & Notifications
	RealTimeAlerts      datatypes.JSON `gorm:"type:json" json:"real_time_alerts"` // []Alert
	AlertThresholds     datatypes.JSON `gorm:"type:json" json:"alert_thresholds"`
	AlertsCount         int            `gorm:"type:int" json:"alerts_count"`
	CriticalAlertsCount int            `gorm:"type:int" json:"critical_alerts_count"`
	LastAlertDate       *time.Time     `gorm:"type:timestamp" json:"last_alert_date,omitempty"`

	// Provider Information
	TelematicsProvider    string     `gorm:"type:varchar(255)" json:"telematics_provider"`
	ProviderAccountID     string     `gorm:"type:varchar(100)" json:"provider_account_id"`
	DataCollectionConsent bool       `gorm:"type:boolean;default:false" json:"data_collection_consent"`
	ConsentDate           *time.Time `gorm:"type:timestamp" json:"consent_date,omitempty"`

	// Behavior Modifiers
	BehaviorModifiers datatypes.JSON `gorm:"type:json" json:"behavior_modifiers"` // []Modifier
	DiscountEarned    float64        `gorm:"type:decimal(5,2)" json:"discount_earned"`
	PenaltyApplied    float64        `gorm:"type:decimal(5,2)" json:"penalty_applied"`

	// Data Quality
	DataCompleteness float64   `gorm:"type:decimal(5,2)" json:"data_completeness"`
	LastDataReceived time.Time `gorm:"type:timestamp" json:"last_data_received"`
	DataGapDays      int       `gorm:"type:int" json:"data_gap_days"`

	// Status
	TelematicsStatus    string     `gorm:"type:varchar(50)" json:"telematics_status"`
	CalibrationStatus   string     `gorm:"type:varchar(50)" json:"calibration_status"`
	LastCalibrationDate *time.Time `gorm:"type:timestamp" json:"last_calibration_date,omitempty"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// PolicyMultiCurrency represents multi-currency support for international policies
type PolicyMultiCurrency struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Currency Configuration
	OriginalCurrency   string `gorm:"type:varchar(3);not null" json:"original_currency"`
	BillingCurrency    string `gorm:"type:varchar(3);not null" json:"billing_currency"`
	ReportingCurrency  string `gorm:"type:varchar(3)" json:"reporting_currency"`
	SettlementCurrency string `gorm:"type:varchar(3)" json:"settlement_currency"`

	// Exchange Rates
	ExchangeRates  datatypes.JSON `gorm:"type:json" json:"exchange_rates"` // map[string]ExchangeRate
	RateSource     string         `gorm:"type:varchar(100)" json:"rate_source"`
	RateDate       time.Time      `gorm:"type:timestamp" json:"rate_date"`
	FixedRate      bool           `gorm:"type:boolean;default:false" json:"fixed_rate"`
	FixedRateValue float64        `gorm:"type:decimal(15,6)" json:"fixed_rate_value"`

	// Currency Conversions
	ConversionHistory datatypes.JSON `gorm:"type:json" json:"conversion_history"` // []Conversion
	TotalConversions  int            `gorm:"type:int" json:"total_conversions"`
	ConversionFees    Money          `gorm:"embedded;embeddedPrefix:conversion_fees_" json:"conversion_fees"`

	// Hedging
	HedgingEnabled       bool           `gorm:"type:boolean;default:false" json:"hedging_enabled"`
	HedgingInstruments   datatypes.JSON `gorm:"type:json" json:"hedging_instruments"` // []Hedge
	HedgingCost          Money          `gorm:"embedded;embeddedPrefix:hedging_cost_" json:"hedging_cost"`
	HedgingEffectiveness float64        `gorm:"type:decimal(5,2)" json:"hedging_effectiveness"`

	// FX Risk
	FXRiskExposure Money   `gorm:"embedded;embeddedPrefix:fx_risk_exposure_" json:"fx_risk_exposure"`
	FXRiskLimit    Money   `gorm:"embedded;embeddedPrefix:fx_risk_limit_" json:"fx_risk_limit"`
	FXVolatility   float64 `gorm:"type:decimal(10,4)" json:"fx_volatility"`
	VaRAmount      Money   `gorm:"embedded;embeddedPrefix:var_amount_" json:"var_amount"`

	// International Taxes
	InternationalTaxes datatypes.JSON `gorm:"type:json" json:"international_taxes"` // map[string]Tax
	WithholdingTax     Money          `gorm:"embedded;embeddedPrefix:withholding_tax_" json:"withholding_tax"`
	VATAmount          Money          `gorm:"embedded;embeddedPrefix:vat_amount_" json:"vat_amount"`
	CustomsDuty        Money          `gorm:"embedded;embeddedPrefix:customs_duty_" json:"customs_duty"`

	// Multi-Currency Amounts
	PremiumOriginal Money `gorm:"embedded;embeddedPrefix:premium_original_" json:"premium_original"`
	PremiumBilling  Money `gorm:"embedded;embeddedPrefix:premium_billing_" json:"premium_billing"`
	ClaimsOriginal  Money `gorm:"embedded;embeddedPrefix:claims_original_" json:"claims_original"`
	ClaimsBilling   Money `gorm:"embedded;embeddedPrefix:claims_billing_" json:"claims_billing"`

	// Settlement
	SettlementExchangeRate float64    `gorm:"type:decimal(15,6)" json:"settlement_exchange_rate"`
	SettlementDate         *time.Time `gorm:"type:timestamp" json:"settlement_date,omitempty"`
	SettlementAmount       Money      `gorm:"embedded;embeddedPrefix:settlement_amount_" json:"settlement_amount"`

	// Status
	CurrencyStatus string    `gorm:"type:varchar(50)" json:"currency_status"`
	LastUpdateDate time.Time `gorm:"type:timestamp" json:"last_update_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// PolicyAutomation represents automation and batch processing for a policy
type PolicyAutomation struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Batch Processing
	BatchProcessingID string `gorm:"type:varchar(100)" json:"batch_processing_id"`
	BatchGroup        string `gorm:"type:varchar(100)" json:"batch_group"`
	BatchPriority     int    `gorm:"type:int" json:"batch_priority"`
	BatchStatus       string `gorm:"type:varchar(50)" json:"batch_status"`

	// Automation Rules
	AutomationRules   datatypes.JSON `gorm:"type:json" json:"automation_rules"` // []Rule
	ActiveRules       int            `gorm:"type:int" json:"active_rules_count"`
	RulesExecuted     int            `gorm:"type:int" json:"rules_executed"`
	LastExecutionDate *time.Time     `gorm:"type:timestamp" json:"last_execution_date,omitempty"`

	// Scheduled Tasks
	ScheduledTasks datatypes.JSON `gorm:"type:json" json:"scheduled_tasks"` // []Task
	NextTaskDate   *time.Time     `gorm:"type:timestamp" json:"next_task_date,omitempty"`
	RecurringTasks int            `gorm:"type:int" json:"recurring_tasks_count"`
	CompletedTasks int            `gorm:"type:int" json:"completed_tasks"`

	// Workflow
	WorkflowStatus      string         `gorm:"type:varchar(50)" json:"workflow_status"`
	CurrentWorkflowStep string         `gorm:"type:varchar(100)" json:"current_workflow_step"`
	WorkflowProgress    float64        `gorm:"type:decimal(5,2)" json:"workflow_progress"`
	WorkflowHistory     datatypes.JSON `gorm:"type:json" json:"workflow_history"`

	// Triggers
	TriggerConditions datatypes.JSON `gorm:"type:json" json:"trigger_conditions"` // []Trigger
	TriggersActivated int            `gorm:"type:int" json:"triggers_activated"`
	LastTriggerDate   *time.Time     `gorm:"type:timestamp" json:"last_trigger_date,omitempty"`

	// Auto-Approval
	AutoApprovalEnabled  bool           `gorm:"type:boolean;default:false" json:"auto_approval_enabled"`
	AutoApprovalCriteria datatypes.JSON `gorm:"type:json" json:"auto_approval_criteria"` // []Criteria
	AutoApprovedCount    int            `gorm:"type:int" json:"auto_approved_count"`
	AutoRejectedCount    int            `gorm:"type:int" json:"auto_rejected_count"`

	// RPA (Robotic Process Automation)
	RPAEnabled     bool           `gorm:"type:boolean;default:false" json:"rpa_enabled"`
	RPAProcesses   datatypes.JSON `gorm:"type:json" json:"rpa_processes"` // []RPAProcess
	RPAExecutions  int            `gorm:"type:int" json:"rpa_executions"`
	RPASuccessRate float64        `gorm:"type:decimal(5,2)" json:"rpa_success_rate"`

	// Performance
	AutomationEfficiency float64 `gorm:"type:decimal(5,2)" json:"automation_efficiency"`
	TimeSavedHours       float64 `gorm:"type:decimal(10,2)" json:"time_saved_hours"`
	ErrorRate            float64 `gorm:"type:decimal(5,2)" json:"error_rate"`

	// Status
	AutomationStatus string    `gorm:"type:varchar(50)" json:"automation_status"`
	LastRunDate      time.Time `gorm:"type:timestamp" json:"last_run_date"`
	NextRunDate      time.Time `gorm:"type:timestamp" json:"next_run_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// =====================================
// METHODS
// =====================================

// IsIntegrationHealthy checks if integrations are healthy
func (pi *PolicyIntegrations) IsIntegrationHealthy() bool {
	return pi.IntegrationStatus == "active" &&
		pi.HealthCheckStatus == "healthy" &&
		pi.SuccessRate > 95
}

// RequiresSync checks if data sync is needed
func (pi *PolicyIntegrations) RequiresSync() bool {
	return time.Now().After(pi.NextSyncDate) ||
		pi.DataSyncStatus == nil
}

// IsUBIActive checks if usage-based insurance is active
func (pt *PolicyTelematics) IsUBIActive() bool {
	return pt.UsageBasedPricing &&
		pt.TelematicsStatus == "active" &&
		pt.DataCollectionConsent
}

// IsSafeDriver checks if driver is considered safe
func (pt *PolicyTelematics) IsSafeDriver() bool {
	return pt.SafetyScore > 80 &&
		pt.DrivingScore > 80 &&
		pt.SpeedingEvents < 5
}

// HasCriticalAlerts checks for critical alerts
func (pt *PolicyTelematics) HasCriticalAlerts() bool {
	return pt.CriticalAlertsCount > 0
}

// RequiresCurrencyConversion checks if currency conversion is needed
func (pmc *PolicyMultiCurrency) RequiresCurrencyConversion() bool {
	return pmc.OriginalCurrency != pmc.BillingCurrency
}

// IsFXRiskHigh checks if FX risk is high
func (pmc *PolicyMultiCurrency) IsFXRiskHigh() bool {
	return pmc.FXRiskExposure.Amount > pmc.FXRiskLimit.Amount ||
		pmc.FXVolatility > 20
}

// IsAutomationActive checks if automation is active
func (pa *PolicyAutomation) IsAutomationActive() bool {
	return pa.AutomationStatus == "active" &&
		pa.ActiveRules > 0
}

// IsWorkflowComplete checks if workflow is complete
func (pa *PolicyAutomation) IsWorkflowComplete() bool {
	return pa.WorkflowProgress >= 100 ||
		pa.WorkflowStatus == "completed"
}

// GetEfficiency returns automation efficiency
func (pa *PolicyAutomation) GetEfficiency() float64 {
	if pa.RPAEnabled {
		return (pa.AutomationEfficiency + pa.RPASuccessRate) / 2
	}
	return pa.AutomationEfficiency
}
