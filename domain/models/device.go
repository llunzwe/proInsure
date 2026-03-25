// Package models contains pure domain entities following hexagonal architecture principles.
// This file contains ONLY:
// - Data structure (entity) definition
// - GORM-specific methods (TableName, BeforeCreate, BeforeUpdate hooks)
// - Simple getter methods that return field values
//
// Business logic has been moved to:
// - /internal/domain/services/device/device_service.go (service implementation)
// - /internal/domain/ports/services/device_service.go (service interface)
// - /internal/domain/ports/repositories/device_repository.go (repository interface)
package models

import (
	"errors"
	"time"

	"smartsure/internal/domain/models/device"
	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Device validation errors
var (
	ErrInvalidIMEI         = errors.New("IMEI cannot be empty")
	ErrInvalidIMEILength   = errors.New("IMEI must be 15 or 17 digits")
	ErrInvalidModel        = errors.New("model cannot be empty")
	ErrInvalidBrand        = errors.New("brand cannot be empty")
	ErrInvalidManufacturer = errors.New("manufacturer cannot be empty")
	ErrInvalidOwnerID      = errors.New("owner ID cannot be nil")
)

// Device represents a device in the database
type Device struct {
	database.BaseModel

	// Embedded structs for organization - Core device attributes
	DeviceIdentification    device.DeviceIdentification    `gorm:"embedded" json:"identification"`
	DeviceClassification    device.DeviceClassification    `gorm:"embedded" json:"classification"`
	DeviceSpecifications    device.DeviceSpecifications    `gorm:"embedded" json:"specifications"`
	DevicePhysicalCondition device.DevicePhysicalCondition `gorm:"embedded" json:"physical_condition"`
	DeviceFinancial         device.DeviceFinancial         `gorm:"embedded" json:"financial"`
	DeviceRiskAssessment    device.DeviceRiskAssessment    `gorm:"embedded" json:"risk_assessment"`
	DeviceSecurity          device.DeviceSecurity          `gorm:"embedded" json:"security"`
	DeviceNetwork           device.DeviceNetwork           `gorm:"embedded" json:"network"`
	DeviceStatusInfo        device.DeviceStatusInfo        `gorm:"embedded" json:"status_info"`
	DeviceLifecycle         device.DeviceLifecycle         `gorm:"embedded" json:"lifecycle"`
	DeviceVerification      device.DeviceVerification      `gorm:"embedded" json:"verification"`
	DeviceCompliance        device.DeviceCompliance        `gorm:"embedded" json:"compliance"`
	DeviceDocumentation     device.DeviceDocumentation     `gorm:"embedded" json:"documentation"`
	DeviceOwnership         device.DeviceOwnership         `gorm:"embedded" json:"ownership"`
	DeviceWarranty          device.DeviceWarranty          `gorm:"embedded" json:"warranty"`

	// Core Relationships
	Owner            User              `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	CorporateAccount *CorporateAccount `gorm:"foreignKey:CorporateAccountID" json:"corporate_account,omitempty"`
	Policies         []Policy          `gorm:"foreignKey:DeviceID" json:"policies,omitempty"`
	Claims           []Claim           `gorm:"foreignKey:DeviceID" json:"claims,omitempty"`

	// Device Ecosystem Relationships
	Swaps          []device.DeviceSwap          `gorm:"foreignKey:DeviceID" json:"swaps,omitempty"`
	TradeIns       []device.DeviceTradeIn       `gorm:"foreignKey:DeviceID" json:"trade_ins,omitempty"`
	Refurbishments []device.DeviceRefurbishment `gorm:"foreignKey:DeviceID" json:"refurbishments,omitempty"`
	Rentals        []device.DeviceRental        `gorm:"foreignKey:DeviceID" json:"rentals,omitempty"`
	Financings     []device.DeviceFinancing     `gorm:"foreignKey:DeviceID" json:"financings,omitempty"`
	Layaways       []device.DeviceLayaway       `gorm:"foreignKey:DeviceID" json:"layaways,omitempty"`
	Repairs        []device.DeviceRepair        `gorm:"foreignKey:DeviceID" json:"repairs,omitempty"`
	Subscriptions  []device.DeviceSubscription  `gorm:"foreignKey:DeviceID" json:"subscriptions,omitempty"`
	MarketListings []device.DeviceMarketplace   `gorm:"foreignKey:DeviceID" json:"market_listings,omitempty"`
	Accessories    []device.DeviceAccessory     `gorm:"foreignKey:DeviceID" json:"accessories,omitempty"`
	SpareParts     []device.DeviceSparePart     `gorm:"foreignKey:DeviceID" json:"spare_parts,omitempty"`

	// Insurance & Claims Relationships
	ClaimHistory        []device.DeviceClaimHistory       `gorm:"foreignKey:DeviceID" json:"claim_history,omitempty"`
	PremiumCalculations []device.DevicePremiumCalculation `gorm:"foreignKey:DeviceID" json:"premium_calculations,omitempty"`
	RiskProfile         *device.DeviceRiskProfile         `gorm:"foreignKey:DeviceID" json:"risk_profile,omitempty"`
	CoverageOptions     []device.DeviceCoverageOptions    `gorm:"foreignKey:DeviceID" json:"coverage_options,omitempty"`
	Deductibles         []device.DeviceDeductibles        `gorm:"foreignKey:DeviceID" json:"deductibles,omitempty"`
	InsuranceAudits     []device.DeviceInsuranceAudit     `gorm:"foreignKey:DeviceID" json:"insurance_audits,omitempty"`
	ClaimValidations    []device.DeviceClaimValidation    `gorm:"foreignKey:DeviceID" json:"claim_validations,omitempty"`

	// Valuation & Financial Relationships
	DepreciationCurve  *device.DeviceDepreciationCurve  `gorm:"foreignKey:DeviceID" json:"depreciation_curve,omitempty"`
	InsurableValue     *device.DeviceInsurableValue     `gorm:"foreignKey:DeviceID" json:"insurable_value,omitempty"`
	TotalCostOwnership *device.DeviceTotalCostOwnership `gorm:"foreignKey:DeviceID" json:"total_cost_ownership,omitempty"`
	ResaleValues       []device.DeviceResaleValue       `gorm:"foreignKey:DeviceID" json:"resale_values,omitempty"`
	FinancialRisk      *device.DeviceFinancialRisk      `gorm:"foreignKey:DeviceID" json:"financial_risk,omitempty"`

	// Fraud Detection Relationships
	FraudPatterns         []device.DeviceFraudPatterns        `gorm:"foreignKey:DeviceID" json:"fraud_patterns,omitempty"`
	AnomalyDetections     []device.DeviceAnomalyDetection     `gorm:"foreignKey:DeviceID" json:"anomaly_detections,omitempty"`
	IdentityVerifications []device.DeviceIdentityVerification `gorm:"foreignKey:DeviceID" json:"identity_verifications,omitempty"`
	FraudInvestigations   []device.DeviceFraudInvestigation   `gorm:"foreignKey:DeviceID" json:"fraud_investigations,omitempty"`
	BlacklistManagement   *device.DeviceBlacklistManagement   `gorm:"foreignKey:DeviceID" json:"blacklist_management,omitempty"`
	FraudPrevention       *device.DeviceFraudPrevention       `gorm:"foreignKey:DeviceID" json:"fraud_prevention,omitempty"`

	// Verification & History Relationships
	Verifications []device.DeviceVerification `gorm:"foreignKey:DeviceID" json:"verifications,omitempty"`
	DeviceHistory []device.DeviceHistory      `gorm:"foreignKey:IMEI;references:IMEI" json:"device_history,omitempty"`

	// Corporate & Enterprise Relationships
	CorporateAssignment *device.CorporateDeviceAssignment `gorm:"foreignKey:DeviceID" json:"corporate_assignment,omitempty"`
	BYODRegistration    *device.BYODDevice                `gorm:"foreignKey:DeviceID" json:"byod_registration,omitempty"`
	PoolAssignments     []device.CorporateDevicePool      `gorm:"many2many:pool_devices;" json:"pool_assignments,omitempty"`

	// IoT & Real-Time Monitoring Relationships
	IoTConnections         []device.IoTConnection         `gorm:"foreignKey:DeviceID" json:"iot_connections,omitempty"`
	ActiveIoTConnection    *device.IoTConnection          `gorm:"foreignKey:DeviceID;where:status='connected';order:last_heartbeat desc;limit:1" json:"active_iot_connection,omitempty"`
	IoTSensorData          []device.IoTSensorData         `gorm:"foreignKey:DeviceID" json:"iot_sensor_data,omitempty"`
	IoTCommands            []device.IoTCommand            `gorm:"foreignKey:DeviceID" json:"iot_commands,omitempty"`
	PendingIoTCommands     []device.IoTCommand            `gorm:"foreignKey:DeviceID;where:status IN ('pending','sent')" json:"pending_iot_commands,omitempty"`
	IoTDeviceConfiguration *device.IoTDeviceConfiguration `gorm:"foreignKey:DeviceID" json:"iot_device_configuration,omitempty"`

	// Predictive Maintenance & Health Relationships
	PredictiveAlerts       []device.PredictiveMaintenanceAlert  `gorm:"foreignKey:DeviceID" json:"predictive_alerts,omitempty"`
	ActivePredictiveAlerts []device.PredictiveMaintenanceAlert  `gorm:"foreignKey:DeviceID;where:status='active'" json:"active_predictive_alerts,omitempty"`
	DeviceHealthScores     []device.DeviceHealthScore           `gorm:"foreignKey:DeviceID;order:last_calculated desc" json:"device_health_scores,omitempty"`
	CurrentHealthScore     *device.DeviceHealthScore            `gorm:"foreignKey:DeviceID;order:last_calculated desc;limit:1" json:"current_health_score,omitempty"`
	FailurePredictions     []device.DeviceFailurePrediction     `gorm:"foreignKey:DeviceID" json:"failure_predictions,omitempty"`
	MaintenanceSchedule    *device.DeviceMaintenanceSchedule    `gorm:"foreignKey:DeviceID" json:"maintenance_schedule,omitempty"`
	UpgradeRecommendations []device.DeviceUpgradeRecommendation `gorm:"foreignKey:DeviceID" json:"upgrade_recommendations,omitempty"`
	ValueForecasts         []device.DeviceValueForecast         `gorm:"foreignKey:DeviceID" json:"value_forecasts,omitempty"`

	// Real-Time Monitoring & Alerting Relationships
	MonitoringSessions      []device.DeviceMonitoringSession  `gorm:"foreignKey:DeviceID" json:"monitoring_sessions,omitempty"`
	ActiveMonitoringSession *device.DeviceMonitoringSession   `gorm:"foreignKey:DeviceID;where:status='active';order:started_at desc;limit:1" json:"active_monitoring_session,omitempty"`
	RealTimeMetrics         []device.DeviceRealTimeMetric     `gorm:"foreignKey:DeviceID" json:"real_time_metrics,omitempty"`
	LiveAlerts              []device.DeviceLiveAlert          `gorm:"foreignKey:DeviceID;where:status='active'" json:"live_alerts,omitempty"`
	AlertConfigurations     []device.DeviceAlertConfiguration `gorm:"foreignKey:DeviceID" json:"alert_configurations,omitempty"`
	AlertHistory            []device.DeviceAlertHistory       `gorm:"foreignKey:DeviceID;order:triggered_at desc" json:"alert_history,omitempty"`
	EscalationRules         []device.DeviceEscalationRule     `gorm:"foreignKey:DeviceID" json:"escalation_rules,omitempty"`

	// Business Intelligence & Analytics Relationships
	CustomerLifetimeValue *device.DeviceCustomerLifetimeValue `gorm:"foreignKey:DeviceID" json:"customer_lifetime_value,omitempty"`
	ChurnPredictions      []device.DeviceChurnPrediction      `gorm:"foreignKey:DeviceID;order:prediction_date desc" json:"churn_predictions,omitempty"`
	PredictiveAnalytics   *device.DevicePredictiveAnalytics   `gorm:"foreignKey:DeviceID" json:"predictive_analytics,omitempty"`
	ProfitabilityMetrics  []device.DeviceProfitabilityMetrics `gorm:"foreignKey:DeviceID;order:calculation_date desc" json:"profitability_metrics,omitempty"`
	MarketIntelligence    *device.DeviceMarketIntelligence    `gorm:"foreignKey:DeviceID" json:"market_intelligence,omitempty"`
	RiskTrendAnalysis     *device.DeviceRiskTrendAnalysis     `gorm:"foreignKey:DeviceID" json:"risk_trend_analysis,omitempty"`
	// Analytics & Intelligence Relationships
	UsagePatterns        []device.DeviceUsagePattern        `gorm:"foreignKey:DeviceID" json:"usage_patterns,omitempty"`
	BehaviorScore        *device.DeviceBehaviorScore        `gorm:"foreignKey:DeviceID" json:"behavior_score,omitempty"`
	LocationHistory      []device.DeviceLocationHistory     `gorm:"foreignKey:DeviceID" json:"location_history,omitempty"`
	NetworkActivity      []device.DeviceNetworkActivity     `gorm:"foreignKey:DeviceID" json:"network_activity,omitempty"`
	CustomerSegmentation *device.DeviceCustomerSegmentation `gorm:"foreignKey:DeviceID" json:"customer_segmentation,omitempty"`
	Profitability        *device.DeviceProfitability        `gorm:"foreignKey:DeviceID" json:"profitability,omitempty"`
	MarketAnalyses       []device.DeviceMarketAnalysis      `gorm:"foreignKey:DeviceID" json:"market_analyses,omitempty"`

	// Compliance & Legal Relationships
	ComplianceStatus    *device.DeviceComplianceStatus    `gorm:"foreignKey:DeviceID" json:"compliance_status,omitempty"`
	LegalHolds          []device.DeviceLegalHolds         `gorm:"foreignKey:DeviceID" json:"legal_holds,omitempty"`
	ExportControls      *device.DeviceExportControls      `gorm:"foreignKey:DeviceID" json:"export_controls,omitempty"`
	DataPrivacy         *device.DeviceDataPrivacy         `gorm:"foreignKey:DeviceID" json:"data_privacy,omitempty"`
	RegulatoryReporting *device.DeviceRegulatoryReporting `gorm:"foreignKey:DeviceID" json:"regulatory_reporting,omitempty"`
	SecurityCompliance  *device.DeviceSecurityCompliance  `gorm:"foreignKey:DeviceID" json:"security_compliance,omitempty"`

	// Emergency & Recovery Relationships
	EmergencyContacts *device.DeviceEmergencyContacts `gorm:"foreignKey:DeviceID" json:"emergency_contacts,omitempty"`
	BackupStatus      *device.DeviceBackupStatus      `gorm:"foreignKey:DeviceID" json:"backup_status,omitempty"`
	RecoveryOptions   *device.DeviceRecoveryOptions   `gorm:"foreignKey:DeviceID" json:"recovery_options,omitempty"`
	EmergencyLocation *device.DeviceEmergencyLocation `gorm:"foreignKey:DeviceID" json:"emergency_location,omitempty"`
	PanicMode         *device.DevicePanicMode         `gorm:"foreignKey:DeviceID" json:"panic_mode,omitempty"`
	DisasterRecovery  *device.DeviceDisasterRecovery  `gorm:"foreignKey:DeviceID" json:"disaster_recovery,omitempty"`

	// Third-Party Integration Relationships
	WarrantyProviders     []device.DeviceWarrantyProviders     `gorm:"foreignKey:DeviceID" json:"warranty_providers,omitempty"`
	ServiceProviders      []device.DeviceServiceProviders      `gorm:"foreignKey:DeviceID" json:"service_providers,omitempty"`
	InsurancePartners     []device.DeviceInsurancePartners     `gorm:"foreignKey:DeviceID" json:"insurance_partners,omitempty"`
	FinancialInstitutions []device.DeviceFinancialInstitutions `gorm:"foreignKey:DeviceID" json:"financial_institutions,omitempty"`
	EcosystemIntegration  *device.DeviceEcosystemIntegration   `gorm:"foreignKey:DeviceID" json:"ecosystem_integration,omitempty"`
	APIIntegrations       []device.DeviceAPIIntegrations       `gorm:"foreignKey:DeviceID" json:"api_integrations,omitempty"`

	// Network & Connectivity Relationships
	NetworkProfile     *device.DeviceNetworkProfile     `gorm:"foreignKey:DeviceID" json:"network_profile,omitempty"`
	DataUsage          []device.DeviceDataUsage         `gorm:"foreignKey:DeviceID" json:"data_usage,omitempty"`
	ConnectivityIssues *device.DeviceConnectivityIssues `gorm:"foreignKey:DeviceID" json:"connectivity_issues,omitempty"`
	WiFiProfile        *device.DeviceWiFiProfile        `gorm:"foreignKey:DeviceID" json:"wifi_profile,omitempty"`
	BluetoothProfile   *device.DeviceBluetoothProfile   `gorm:"foreignKey:DeviceID" json:"bluetooth_profile,omitempty"`
	HotspotUsage       *device.DeviceHotspotUsage       `gorm:"foreignKey:DeviceID" json:"hotspot_usage,omitempty"`

	// Advanced Security & Protection Relationships
	SecurityAudits     []device.DeviceSecurityAudit     `gorm:"foreignKey:DeviceID" json:"security_audits,omitempty"`
	BiometricData      *device.DeviceBiometricData      `gorm:"foreignKey:DeviceID" json:"biometric_data,omitempty"`
	EncryptionStatus   *device.DeviceEncryptionStatus   `gorm:"foreignKey:DeviceID" json:"encryption_status,omitempty"`
	AntivirusStatus    *device.DeviceAntivirusStatus    `gorm:"foreignKey:DeviceID" json:"antivirus_status,omitempty"`
	AccessControl      *device.DeviceAccessControl      `gorm:"foreignKey:DeviceID" json:"access_control,omitempty"`
	ThreatIntelligence *device.DeviceThreatIntelligence `gorm:"foreignKey:DeviceID" json:"threat_intelligence,omitempty"`

	// Performance Monitoring Relationships
	PerformanceMetrics []device.DevicePerformanceMetrics `gorm:"foreignKey:DeviceID" json:"performance_metrics,omitempty"`
	BatteryAnalytics   []device.DeviceBatteryAnalytics   `gorm:"foreignKey:DeviceID" json:"battery_analytics,omitempty"`
	TemperatureHistory []device.DeviceTemperatureHistory `gorm:"foreignKey:DeviceID" json:"temperature_history,omitempty"`
	StorageHealth      *device.DeviceStorageHealth       `gorm:"foreignKey:DeviceID" json:"storage_health,omitempty"`
	MemoryManagement   *device.DeviceMemoryManagement    `gorm:"foreignKey:DeviceID" json:"memory_management,omitempty"`
	ProcessorHealth    *device.DeviceProcessorHealth     `gorm:"foreignKey:DeviceID" json:"processor_health,omitempty"`

	// Quality Assurance Relationships
	QualityMetrics      []device.DeviceQualityMetrics     `gorm:"foreignKey:DeviceID" json:"quality_metrics,omitempty"`
	DefectHistory       []device.DeviceDefectHistory      `gorm:"foreignKey:DeviceID" json:"defect_history,omitempty"`
	RecallStatus        []device.DeviceRecallStatus       `gorm:"foreignKey:DeviceID" json:"recall_status,omitempty"`
	TestResults         []device.DeviceTestResults        `gorm:"foreignKey:DeviceID" json:"test_results,omitempty"`
	QualityIncidents    []device.DeviceQualityIncidents   `gorm:"foreignKey:DeviceID" json:"quality_incidents,omitempty"`
	CertificationStatus *device.DeviceCertificationStatus `gorm:"foreignKey:DeviceID" json:"certification_status,omitempty"`

	// Environmental & Sustainability Relationships
	CarbonFootprint       []device.DeviceCarbonFootprint       `gorm:"foreignKey:DeviceID" json:"carbon_footprint,omitempty"`
	RecyclingScore        []device.DeviceRecyclingScore        `gorm:"foreignKey:DeviceID" json:"recycling_score,omitempty"`
	SustainabilityMetrics []device.DeviceSustainabilityMetrics `gorm:"foreignKey:DeviceID" json:"sustainability_metrics,omitempty"`
	EcoLabel              *device.DeviceEcoLabel               `gorm:"foreignKey:DeviceID" json:"eco_label,omitempty"`
	LifecycleAssessment   []device.DeviceLifecycleAssessment   `gorm:"foreignKey:DeviceID" json:"lifecycle_assessment,omitempty"`
	Repairability         *device.DeviceRepairability          `gorm:"foreignKey:DeviceID" json:"repairability,omitempty"`

	// Customer Experience Relationships
	SatisfactionScores []device.DeviceSatisfactionScore `gorm:"foreignKey:DeviceID" json:"satisfaction_scores,omitempty"`
	Recommendations    []device.DeviceRecommendations   `gorm:"foreignKey:DeviceID" json:"recommendations,omitempty"`
	LoyaltyPrograms    []device.DeviceLoyaltyProgram    `gorm:"foreignKey:DeviceID" json:"loyalty_programs,omitempty"`
	Feedbacks          []device.DeviceFeedback          `gorm:"foreignKey:DeviceID" json:"feedbacks,omitempty"`
	UserJourneys       []device.DeviceUserJourney       `gorm:"foreignKey:DeviceID" json:"user_journeys,omitempty"`

	// Supply Chain & International Relationships
	SupplyChain        *device.DeviceSupplyChain        `gorm:"foreignKey:DeviceID" json:"supply_chain,omitempty"`
	MultiCurrency      *device.DeviceMultiCurrency      `gorm:"foreignKey:DeviceID" json:"multi_currency,omitempty"`
	DocumentManagement *device.DeviceDocumentManagement `gorm:"foreignKey:DeviceID" json:"document_management,omitempty"`

	// Contract & Logistics Relationships
	Contract      *device.DeviceContract      `gorm:"foreignKey:DeviceID" json:"contract,omitempty"`
	Shipping      *device.DeviceShipping      `gorm:"foreignKey:DeviceID" json:"shipping,omitempty"`
	TaxAccounting *device.DeviceTaxAccounting `gorm:"foreignKey:DeviceID" json:"tax_accounting,omitempty"`

	// Cross-Device & Engagement Relationships
	CrossRelationship *device.DeviceCrossRelationship `gorm:"foreignKey:DeviceID" json:"cross_relationship,omitempty"`
	Gamification      *device.DeviceGamification      `gorm:"foreignKey:DeviceID" json:"gamification,omitempty"`
	Education         *device.DeviceEducation         `gorm:"foreignKey:DeviceID" json:"education,omitempty"`

	// Social & Automation Relationships
	Social            *device.DeviceSocial            `gorm:"foreignKey:DeviceID" json:"social,omitempty"`
	Automation        *device.DeviceAutomation        `gorm:"foreignKey:DeviceID" json:"automation,omitempty"`
	AdvancedAnalytics *device.DeviceAdvancedAnalytics `gorm:"foreignKey:DeviceID" json:"advanced_analytics,omitempty"`

	// Communication & Scheduling Relationships
	NotificationManagement *device.DeviceNotificationManagement `gorm:"foreignKey:DeviceID" json:"notification_management,omitempty"`
	CalendarScheduling     *device.DeviceCalendarScheduling     `gorm:"foreignKey:DeviceID" json:"calendar_scheduling,omitempty"`
	Comparison             *device.DeviceComparison             `gorm:"foreignKey:DeviceID" json:"comparison,omitempty"`

	// Privacy & Pricing Relationships
	Privacy         *device.DevicePrivacy         `gorm:"foreignKey:DeviceID" json:"privacy,omitempty"`
	AdvancedPricing *device.DeviceAdvancedPricing `gorm:"foreignKey:DeviceID" json:"advanced_pricing,omitempty"`

	// Advanced Automation & Workflow Relationships
	AutomationRules    []device.DeviceAutomationRule    `gorm:"foreignKey:DeviceID" json:"automation_rules,omitempty"`
	WorkflowInstances  []device.DeviceWorkflowInstance  `gorm:"foreignKey:DeviceID" json:"workflow_instances,omitempty"`
	SmartActions       []device.DeviceSmartAction       `gorm:"foreignKey:DeviceID" json:"smart_actions,omitempty"`
	ConditionalActions []device.DeviceConditionalAction `gorm:"foreignKey:DeviceID" json:"conditional_actions,omitempty"`
	TriggerActions     []device.DeviceTriggerAction     `gorm:"foreignKey:DeviceID" json:"trigger_actions,omitempty"`

	// Advanced Analytics & ML Relationships
	MLModels           []device.DeviceMLModel            `gorm:"foreignKey:DeviceID" json:"ml_models,omitempty"`
	ABTests            []device.DeviceABTest             `gorm:"foreignKey:DeviceID" json:"ab_tests,omitempty"`
	ExperimentTracking []device.DeviceExperimentTracking `gorm:"foreignKey:DeviceID" json:"experiment_tracking,omitempty"`
	FeatureFlags       []device.DeviceFeatureFlag        `gorm:"foreignKey:DeviceID" json:"feature_flags,omitempty"`
	ModelValidations   []device.DeviceModelValidation    `gorm:"foreignKey:DeviceID" json:"model_validations,omitempty"`
}

// TableName returns the table name for Device model
func (d *Device) TableName() string {
	return "devices"
}

// BeforeCreate handles pre-creation logic (UUID generation handled by BaseModel)
func (d *Device) BeforeCreate(tx *gorm.DB) error {
	// Call parent BeforeCreate to handle UUID generation
	if err := d.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// ============================================
// SIMPLE GETTERS - Return field values only
// ============================================

// GetID returns the device ID
func (d *Device) GetID() uuid.UUID {
	return d.ID
}

// GetIMEI returns the device IMEI
func (d *Device) GetIMEI() string {
	return d.DeviceIdentification.IMEI
}

// GetStatus returns the device status
func (d *Device) GetStatus() string {
	return string(d.DeviceStatusInfo.Status)
}

// IsActive checks if device status is active
func (d *Device) IsActive() bool {
	return d.DeviceStatusInfo.Status == device.DeviceStatusActive
}

// IsStolen checks if device is reported stolen
func (d *Device) IsStolen() bool {
	return d.DeviceStatusInfo.IsStolen
}

// GetCondition returns the device condition
func (d *Device) GetCondition() string {
	return string(d.DevicePhysicalCondition.Condition)
}

// GetGrade returns the device grade
func (d *Device) GetGrade() string {
	return string(d.DevicePhysicalCondition.Grade)
}

// GetOwnerID returns the owner ID
func (d *Device) GetOwnerID() uuid.UUID {
	return d.DeviceOwnership.OwnerID
}

// GetPurchasePrice returns the purchase price
func (d *Device) GetPurchasePrice() float64 {
	return d.DeviceFinancial.PurchasePrice.Amount
}

// GetCurrentValue returns the current value
func (d *Device) GetCurrentValue() float64 {
	return d.DeviceFinancial.CurrentValue.Amount
}

// GetMarketValue returns the market value
func (d *Device) GetMarketValue() float64 {
	return d.DeviceFinancial.MarketValue.Amount
}

// GetBlacklistStatus returns the blacklist status
func (d *Device) GetBlacklistStatus() string {
	return string(d.DeviceRiskAssessment.BlacklistStatus)
}

// GetRiskScore returns the risk score
func (d *Device) GetRiskScore() float64 {
	return d.DeviceRiskAssessment.RiskScore
}

// IsVerified checks if device is verified
func (d *Device) IsVerified() bool {
	return d.DeviceVerification.IsVerified
}

// GetVerificationDate returns the verification date
func (d *Device) GetVerificationDate() *time.Time {
	return d.DeviceVerification.VerificationDate
}

// GetModel returns the device model
func (d *Device) GetModel() string {
	return d.DeviceClassification.Model
}

// GetBrand returns the device brand
func (d *Device) GetBrand() string {
	return d.DeviceClassification.Brand
}

// GetManufacturer returns the device manufacturer
func (d *Device) GetManufacturer() string {
	return d.DeviceClassification.Manufacturer
}

// BeforeUpdate handles pre-update logic
func (d *Device) BeforeUpdate(tx *gorm.DB) error {
	// Update the UpdatedAt timestamp
	d.UpdatedAt = time.Now()
	return nil
}

// BeforeSave is a GORM hook for basic field validation (both create and update)
func (d *Device) BeforeSave(tx *gorm.DB) error {
	// Basic field validation for data integrity
	if d.DeviceIdentification.IMEI == "" {
		return ErrInvalidIMEI
	}

	if len(d.DeviceIdentification.IMEI) != 15 && len(d.DeviceIdentification.IMEI) != 17 {
		return ErrInvalidIMEILength
	}

	if d.DeviceClassification.Model == "" {
		return ErrInvalidModel
	}

	if d.DeviceClassification.Brand == "" {
		return ErrInvalidBrand
	}

	if d.DeviceClassification.Manufacturer == "" {
		return ErrInvalidManufacturer
	}

	if d.DeviceOwnership.OwnerID == uuid.Nil {
		return ErrInvalidOwnerID
	}

	// Set defaults if not provided
	if d.DeviceStatusInfo.Status == "" {
		d.DeviceStatusInfo.Status = device.DeviceStatusActive
	}

	if d.DevicePhysicalCondition.Condition == "" {
		d.DevicePhysicalCondition.Condition = "good"
	}

	// Ensure non-negative values
	if d.DeviceFinancial.PurchasePrice.Amount < 0 {
		d.DeviceFinancial.PurchasePrice.Amount = 0
	}

	if d.DeviceFinancial.CurrentValue.Amount < 0 {
		d.DeviceFinancial.CurrentValue.Amount = 0
	}

	if d.DeviceFinancial.MarketValue.Amount < 0 {
		d.DeviceFinancial.MarketValue.Amount = 0
	}

	return nil
}

// ============================================
// RELATIONSHIP CHECKS - Simple boolean checks
// ============================================

// HasPolicies checks if device has policies
func (d *Device) HasPolicies() bool {
	return len(d.Policies) > 0
}

// HasClaims checks if device has claims
func (d *Device) HasClaims() bool {
	return len(d.Claims) > 0
}

// HasRentals checks if device has rentals
func (d *Device) HasRentals() bool {
	return len(d.Rentals) > 0
}

// HasFinancings checks if device has financings
func (d *Device) HasFinancings() bool {
	return len(d.Financings) > 0
}

// HasSubscriptions checks if device has subscriptions
func (d *Device) HasSubscriptions() bool {
	return len(d.Subscriptions) > 0
}

// HasRepairs checks if device has repairs
func (d *Device) HasRepairs() bool {
	return len(d.Repairs) > 0
}

// HasMarketListings checks if device has market listings
func (d *Device) HasMarketListings() bool {
	return len(d.MarketListings) > 0
}

// GetPurchaseDate returns the purchase date
func (d *Device) GetPurchaseDate() *time.Time {
	return d.DeviceFinancial.PurchaseDate
}

// GetWarrantyExpiry returns the warranty expiry date
func (d *Device) GetWarrantyExpiry() *time.Time {
	return d.DeviceWarranty.WarrantyExpiry
}

// IsUnderWarranty checks if device is currently under warranty
func (d *Device) IsUnderWarranty() bool {
	if d.DeviceWarranty.WarrantyExpiry == nil {
		return false
	}
	return time.Now().Before(*d.DeviceWarranty.WarrantyExpiry) && d.DeviceWarranty.WarrantyStatus == "active"
}

// GetWarrantyStatus returns the warranty status
func (d *Device) GetWarrantyStatus() string {
	return d.DeviceWarranty.WarrantyStatus
}

// GetSerialNumber returns the serial number
func (d *Device) GetSerialNumber() string {
	return d.DeviceIdentification.SerialNumber
}

// IsCorporateDevice checks if device is corporate owned
func (d *Device) IsCorporateDevice() bool {
	return d.DeviceOwnership.CorporateAccountID != nil || d.CorporateAssignment != nil
}

// IsBYODDevice checks if device is BYOD
func (d *Device) IsBYODDevice() bool {
	return d.BYODRegistration != nil
}

// GetCorporateAccountID returns the corporate account ID
func (d *Device) GetCorporateAccountID() *uuid.UUID {
	return d.DeviceOwnership.CorporateAccountID
}
