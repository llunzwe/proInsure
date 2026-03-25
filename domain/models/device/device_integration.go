package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceWarrantyProviders manages multiple warranty provider relationships
type DeviceWarrantyProviders struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ProviderID string    `gorm:"index" json:"provider_id"`

	// Provider Information
	ProviderName     string    `json:"provider_name"`
	ProviderType     string    `json:"provider_type"`   // manufacturer, extended, third_party
	ProviderStatus   string    `json:"provider_status"` // active, inactive, suspended
	ContractNumber   string    `json:"contract_number"`
	RegistrationDate time.Time `json:"registration_date"`

	// Warranty Terms
	CoverageType       string    `json:"coverage_type"` // comprehensive, limited, accidental
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	CoverageAmount     float64   `json:"coverage_amount"`
	Deductible         float64   `json:"deductible"`
	TermsAndConditions string    `gorm:"type:text" json:"terms_and_conditions"`

	// Performance Metrics
	ClaimsProcessed    int     `json:"claims_processed"`
	ClaimsApproved     int     `json:"claims_approved"`
	ClaimsDenied       int     `json:"claims_denied"`
	AverageProcessTime float64 `json:"average_process_time"` // days
	ApprovalRate       float64 `json:"approval_rate"`        // percentage
	CustomerRating     float64 `json:"customer_rating"`      // 0-5

	// Coverage Management
	CoverageOverlap      bool   `json:"coverage_overlap"`
	OverlappingProviders string `gorm:"type:json" json:"overlapping_providers"` // JSON array
	PrimaryCoverage      bool   `json:"primary_coverage"`
	CoverageGaps         string `gorm:"type:json" json:"coverage_gaps"` // JSON array

	// Provider History
	PreviousProviders string     `gorm:"type:json" json:"previous_providers"` // JSON array
	SwitchCount       int        `json:"switch_count"`
	LastSwitchDate    *time.Time `json:"last_switch_date"`
	SwitchReasons     string     `gorm:"type:json" json:"switch_reasons"` // JSON array

	// Warranty Stacking
	StackingEnabled   bool    `json:"stacking_enabled"`
	StackedWarranties string  `gorm:"type:json" json:"stacked_warranties"` // JSON array
	TotalCoverage     float64 `json:"total_coverage"`
	EffectiveCoverage float64 `json:"effective_coverage"`

	// Dispute Resolution
	DisputesRaised   int    `json:"disputes_raised"`
	DisputesResolved int    `json:"disputes_resolved"`
	OpenDisputes     int    `json:"open_disputes"`
	DisputeHistory   string `gorm:"type:json" json:"dispute_history"` // JSON array

	// Cost Analysis
	PremiumCost      float64 `json:"premium_cost"`
	TotalClaimsPaid  float64 `json:"total_claims_paid"`
	CostBenefitRatio float64 `json:"cost_benefit_ratio"`
	SavedAmount      float64 `json:"saved_amount"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceServiceProviders manages authorized service center network
type DeviceServiceProviders struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ProviderID string    `gorm:"index" json:"provider_id"`

	// Provider Details
	ProviderName        string `json:"provider_name"`
	ServiceType         string `json:"service_type"`         // repair, maintenance, inspection
	AuthorizationStatus string `json:"authorization_status"` // authorized, preferred, standard
	CertificationLevel  string `json:"certification_level"`  // gold, silver, bronze
	ProviderCategory    string `json:"provider_category"`    // official, partner, third_party

	// Service Metrics
	ServicesCompleted int     `json:"services_completed"`
	AverageRating     float64 `json:"average_rating"`      // 0-5
	QualityScore      float64 `json:"quality_score"`       // 0-100
	OnTimeDelivery    float64 `json:"on_time_delivery"`    // percentage
	RepairSuccessRate float64 `json:"repair_success_rate"` // percentage

	// Turnaround Time
	AverageTurnaround float64 `json:"average_turnaround"` // hours
	FastestService    float64 `json:"fastest_service"`    // hours
	SlowestService    float64 `json:"slowest_service"`    // hours
	CurrentQueueTime  float64 `json:"current_queue_time"` // hours

	// Cost Comparison
	AverageCost          float64 `json:"average_cost"`
	MinimumCost          float64 `json:"minimum_cost"`
	MaximumCost          float64 `json:"maximum_cost"`
	PriceCompetitiveness string  `json:"price_competitiveness"` // high, medium, low
	DiscountOffered      float64 `json:"discount_offered"`      // percentage

	// Preferred Programs
	IsPreferred       bool       `json:"is_preferred"`
	PreferredSince    *time.Time `json:"preferred_since"`
	PreferredBenefits string     `gorm:"type:json" json:"preferred_benefits"` // JSON array
	LoyaltyPoints     int        `json:"loyalty_points"`

	// Service Level Agreement
	SLAActive     bool    `json:"sla_active"`
	SLATerms      string  `gorm:"type:json" json:"sla_terms"` // JSON object
	SLACompliance float64 `json:"sla_compliance"`             // percentage
	SLAViolations int     `json:"sla_violations"`

	// Certification Status
	Certifications      string     `gorm:"type:json" json:"certifications"` // JSON array
	CertificationExpiry *time.Time `json:"certification_expiry"`
	TrainingCompleted   bool       `json:"training_completed"`

	// Customer Satisfaction
	SatisfactionScore float64 `json:"satisfaction_score"` // 0-100
	TotalReviews      int     `json:"total_reviews"`
	PositiveReviews   int     `json:"positive_reviews"`
	Complaints        int     `json:"complaints"`

	// Geographic Coverage
	ServiceLocation string  `json:"service_location"`
	CoverageArea    string  `gorm:"type:json" json:"coverage_area"` // JSON object
	Distance        float64 `json:"distance"`                       // km
	HomeService     bool    `json:"home_service"`
	ServiceRadius   float64 `json:"service_radius"` // km

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceInsurancePartners manages other insurance coverage coordination
type DeviceInsurancePartners struct {
	database.BaseModel
	DeviceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	PartnerID string    `gorm:"index" json:"partner_id"`

	// Partner Information
	PartnerName    string `json:"partner_name"`
	InsuranceType  string `json:"insurance_type"` // property, auto, health, travel
	PolicyNumber   string `json:"policy_number"`
	CoverageStatus string `json:"coverage_status"` // active, expired, cancelled

	// Coverage Details
	CoverageAmount float64   `json:"coverage_amount"`
	Deductible     float64   `json:"deductible"`
	CoverageStart  time.Time `json:"coverage_start"`
	CoverageEnd    time.Time `json:"coverage_end"`
	CoverageScope  string    `gorm:"type:json" json:"coverage_scope"` // JSON array

	// Coverage Coordination
	IsPrimary          bool   `json:"is_primary"`
	CoordinationStatus string `json:"coordination_status"`                 // coordinated, pending, none
	CoordinationRules  string `gorm:"type:json" json:"coordination_rules"` // JSON object
	BenefitOrder       int    `json:"benefit_order"`

	// Coverage Gap Analysis
	GapsIdentified      bool   `json:"gaps_identified"`
	CoverageGaps        string `gorm:"type:json" json:"coverage_gaps"`        // JSON array
	OverlappingCoverage string `gorm:"type:json" json:"overlapping_coverage"` // JSON array
	RecommendedActions  string `gorm:"type:json" json:"recommended_actions"`  // JSON array

	// Claims Coordination
	ClaimsCoordinated   int        `json:"claims_coordinated"`
	PendingCoordination int        `json:"pending_coordination"`
	CoordinationIssues  string     `gorm:"type:json" json:"coordination_issues"` // JSON array
	LastCoordination    *time.Time `json:"last_coordination"`

	// Subrogation Management
	SubrogationActive bool    `json:"subrogation_active"`
	SubrogationAmount float64 `json:"subrogation_amount"`
	SubrogationStatus string  `json:"subrogation_status"` // pending, in_progress, completed
	RecoveredAmount   float64 `json:"recovered_amount"`

	// Co-insurance
	CoInsuranceEnabled   bool    `json:"co_insurance_enabled"`
	CoInsurancePercent   float64 `json:"co_insurance_percent"`
	SharedResponsibility float64 `json:"shared_responsibility"`

	// Performance Metrics
	ResponseTime      float64 `json:"response_time"`      // hours
	SettlementTime    float64 `json:"settlement_time"`    // days
	DisputeRate       float64 `json:"dispute_rate"`       // percentage
	SatisfactionScore float64 `json:"satisfaction_score"` // 0-100

	// Integration Status
	APIIntegrated     bool       `json:"api_integrated"`
	IntegrationStatus string     `json:"integration_status"` // connected, disconnected, error
	LastSyncTime      *time.Time `json:"last_sync_time"`
	SyncErrors        int        `json:"sync_errors"`

	// Settlement Reconciliation
	TotalSettlements     float64    `json:"total_settlements"`
	PendingSettlements   float64    `json:"pending_settlements"`
	ReconciliationStatus string     `json:"reconciliation_status"` // balanced, discrepancy, pending
	LastReconciliation   *time.Time `json:"last_reconciliation"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceFinancialInstitutions manages financial provider integrations
type DeviceFinancialInstitutions struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	InstitutionID string    `gorm:"index" json:"institution_id"`

	// Institution Details
	InstitutionName   string    `json:"institution_name"`
	InstitutionType   string    `json:"institution_type"` // bank, credit_union, fintech, payment_processor
	AccountStatus     string    `json:"account_status"`   // active, pending, suspended, closed
	RelationshipSince time.Time `json:"relationship_since"`

	// Financing Integration
	FinancingEnabled bool    `json:"financing_enabled"`
	CreditLimit      float64 `json:"credit_limit"`
	AvailableCredit  float64 `json:"available_credit"`
	InterestRate     float64 `json:"interest_rate"`
	PaymentTerms     string  `json:"payment_terms"`

	// Payment Gateway
	GatewayActive    bool    `json:"gateway_active"`
	GatewayProvider  string  `json:"gateway_provider"`
	MerchantID       string  `json:"merchant_id"`
	TransactionFee   float64 `json:"transaction_fee"`   // percentage
	SettlementPeriod int     `json:"settlement_period"` // days

	// Bank Account Linking
	AccountLinked       bool   `json:"account_linked"`
	AccountType         string `json:"account_type"` // checking, savings, credit
	AccountLast4        string `json:"account_last4"`
	VerificationStatus  string `json:"verification_status"` // verified, pending, failed
	MicroDepositsStatus string `json:"micro_deposits_status"`

	// Credit Check Integration
	CreditCheckEnabled bool       `json:"credit_check_enabled"`
	LastCreditCheck    *time.Time `json:"last_credit_check"`
	CreditScore        int        `json:"credit_score"`
	CreditBureau       string     `json:"credit_bureau"`
	SoftPullOnly       bool       `json:"soft_pull_only"`

	// Payment Processing
	ProcessorStatus       string  `json:"processor_status"` // active, maintenance, error
	TransactionsProcessed int     `json:"transactions_processed"`
	TotalVolume           float64 `json:"total_volume"`
	SuccessRate           float64 `json:"success_rate"` // percentage
	AverageTicket         float64 `json:"average_ticket"`

	// Financial Reporting
	ReportingEnabled   bool       `json:"reporting_enabled"`
	ReportingFrequency string     `json:"reporting_frequency"` // daily, weekly, monthly
	LastReportDate     *time.Time `json:"last_report_date"`
	ReportFormats      string     `gorm:"type:json" json:"report_formats"` // JSON array

	// Automated Payments
	AutoPayEnabled bool   `json:"auto_pay_enabled"`
	AutoPayDay     int    `json:"auto_pay_day"`                     // day of month
	PaymentMethods string `gorm:"type:json" json:"payment_methods"` // JSON array
	RetryAttempts  int    `json:"retry_attempts"`

	// Reconciliation
	ReconciliationStatus string     `json:"reconciliation_status"` // balanced, discrepancy, pending
	LastReconciliation   *time.Time `json:"last_reconciliation"`
	Discrepancies        string     `gorm:"type:json" json:"discrepancies"` // JSON array

	// Fraud Detection
	FraudMonitoring     bool    `json:"fraud_monitoring"`
	FraudAlerts         int     `json:"fraud_alerts"`
	BlockedTransactions int     `json:"blocked_transactions"`
	FraudScore          float64 `json:"fraud_score"`

	// Chargeback Management
	ChargebacksReceived int     `json:"chargebacks_received"`
	ChargebacksDisputed int     `json:"chargebacks_disputed"`
	ChargebacksWon      int     `json:"chargebacks_won"`
	ChargebackRate      float64 `json:"chargeback_rate"` // percentage

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceEcosystemIntegration manages smart device ecosystem connections
type DeviceEcosystemIntegration struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Smart Home Integration
	SmartHomeEnabled  bool       `json:"smart_home_enabled"`
	SmartHomePlatform string     `json:"smart_home_platform"` // alexa, google_home, homekit, smartthings
	ConnectedDevices  int        `json:"connected_devices"`
	HomeAutomations   string     `gorm:"type:json" json:"home_automations"` // JSON array
	LastHomeSync      *time.Time `json:"last_home_sync"`

	// Wearable Connections
	WearablesConnected int        `json:"wearables_connected"`
	WearableTypes      string     `gorm:"type:json" json:"wearable_types"` // JSON array
	HealthDataSharing  bool       `json:"health_data_sharing"`
	FitnessIntegration bool       `json:"fitness_integration"`
	LastWearableSync   *time.Time `json:"last_wearable_sync"`

	// Automotive Integration
	VehicleConnected   bool   `json:"vehicle_connected"`
	VehicleSystem      string `json:"vehicle_system"` // carplay, android_auto, custom
	VehicleMake        string `json:"vehicle_make"`
	VehicleModel       string `json:"vehicle_model"`
	DrivingDataEnabled bool   `json:"driving_data_enabled"`

	// Health Device Connectivity
	HealthDevices      string `gorm:"type:json" json:"health_devices"` // JSON array
	MedicalDataSync    bool   `json:"medical_data_sync"`
	EmergencyDataShare bool   `json:"emergency_data_share"`
	HealthPlatform     string `json:"health_platform"`

	// IoT Interactions
	IoTDevicesLinked   int    `json:"iot_devices_linked"`
	IoTPlatforms       string `gorm:"type:json" json:"iot_platforms"` // JSON array
	IoTProtocols       string `gorm:"type:json" json:"iot_protocols"` // JSON array
	MeshNetworkEnabled bool   `json:"mesh_network_enabled"`

	// Cloud Services
	CloudServicesLinked string `gorm:"type:json" json:"cloud_services_linked"` // JSON array
	CloudStorageUsed    int64  `json:"cloud_storage_used"`                     // bytes
	CloudSyncEnabled    bool   `json:"cloud_sync_enabled"`
	MultiCloudEnabled   bool   `json:"multi_cloud_enabled"`

	// Productivity Apps
	ProductivityApps string `gorm:"type:json" json:"productivity_apps"` // JSON array
	CalendarSync     bool   `json:"calendar_sync"`
	ContactSync      bool   `json:"contact_sync"`
	DocumentSync     bool   `json:"document_sync"`

	// Entertainment Systems
	EntertainmentLinks string `gorm:"type:json" json:"entertainment_links"` // JSON array
	StreamingServices  string `gorm:"type:json" json:"streaming_services"`  // JSON array
	GamingPlatforms    string `gorm:"type:json" json:"gaming_platforms"`    // JSON array
	MediaSharing       bool   `json:"media_sharing"`

	// Security Integration
	SecuritySystems     string `gorm:"type:json" json:"security_systems"`  // JSON array
	BiometricDevices    string `gorm:"type:json" json:"biometric_devices"` // JSON array
	AccessControlLinked bool   `json:"access_control_linked"`
	SecurityAlertSync   bool   `json:"security_alert_sync"`

	// Cross-Device Sync
	SyncEnabled     bool       `json:"sync_enabled"`
	SyncedDevices   string     `gorm:"type:json" json:"synced_devices"` // JSON array
	LastSyncTime    *time.Time `json:"last_sync_time"`
	SyncConflicts   int        `json:"sync_conflicts"`
	DataConsistency float64    `json:"data_consistency"` // percentage

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceAPIIntegrations manages third-party API connections
type DeviceAPIIntegrations struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	IntegrationID string    `gorm:"index" json:"integration_id"`

	// API Connection
	APIName          string `json:"api_name"`
	APIProvider      string `json:"api_provider"`
	APIVersion       string `json:"api_version"`
	EndpointURL      string `json:"endpoint_url"`
	ConnectionStatus string `json:"connection_status"` // connected, disconnected, error

	// API Usage Metrics
	TotalRequests      int64      `json:"total_requests"`
	SuccessfulRequests int64      `json:"successful_requests"`
	FailedRequests     int64      `json:"failed_requests"`
	AverageLatency     float64    `json:"average_latency"` // ms
	LastRequestTime    *time.Time `json:"last_request_time"`

	// Rate Limiting
	RateLimitEnabled  bool       `json:"rate_limit_enabled"`
	RequestsPerMinute int        `json:"requests_per_minute"`
	RequestsPerHour   int        `json:"requests_per_hour"`
	RequestsPerDay    int        `json:"requests_per_day"`
	CurrentUsage      int        `json:"current_usage"`
	LimitResetTime    *time.Time `json:"limit_reset_time"`

	// Error Tracking
	ErrorRate            float64    `json:"error_rate"` // percentage
	LastError            *time.Time `json:"last_error"`
	CommonErrors         string     `gorm:"type:json" json:"common_errors"` // JSON array
	ErrorRetryCount      int        `json:"error_retry_count"`
	CircuitBreakerStatus string     `json:"circuit_breaker_status"` // closed, open, half_open

	// Health Monitoring
	HealthStatus        string     `json:"health_status"` // healthy, degraded, unhealthy
	Uptime              float64    `json:"uptime"`        // percentage
	LastHealthCheck     *time.Time `json:"last_health_check"`
	HealthCheckInterval int        `json:"health_check_interval"` // minutes
	AlertsEnabled       bool       `json:"alerts_enabled"`

	// Version Management
	CurrentVersion      string `json:"current_version"`
	LatestVersion       string `json:"latest_version"`
	UpdateAvailable     bool   `json:"update_available"`
	AutoUpdate          bool   `json:"auto_update"`
	DeprecationWarnings string `gorm:"type:json" json:"deprecation_warnings"` // JSON array

	// Webhook Configuration
	WebhooksEnabled     bool       `json:"webhooks_enabled"`
	WebhookURL          string     `json:"webhook_url"`
	WebhookEvents       string     `gorm:"type:json" json:"webhook_events"` // JSON array
	WebhookSecret       string     `json:"webhook_secret"`                  // encrypted
	LastWebhookReceived *time.Time `json:"last_webhook_received"`

	// OAuth Management
	OAuthEnabled  bool       `json:"oauth_enabled"`
	OAuthProvider string     `json:"oauth_provider"`
	TokenExpiry   *time.Time `json:"token_expiry"`
	RefreshToken  string     `json:"refresh_token"`           // encrypted
	Scopes        string     `gorm:"type:json" json:"scopes"` // JSON array

	// Cost Tracking
	APITier        string  `json:"api_tier"` // free, basic, premium, enterprise
	MonthlyCost    float64 `json:"monthly_cost"`
	CostPerRequest float64 `json:"cost_per_request"`
	BillingCycle   string  `json:"billing_cycle"`
	UsageVsQuota   float64 `json:"usage_vs_quota"` // percentage

	// Performance Metrics
	ResponseTime95th float64 `json:"response_time_95th"` // ms
	ResponseTime99th float64 `json:"response_time_99th"` // ms
	Throughput       float64 `json:"throughput"`         // requests per second
	DataTransferred  int64   `json:"data_transferred"`   // bytes
	CacheHitRate     float64 `json:"cache_hit_rate"`     // percentage

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceWarrantyProviders
func (dwp *DeviceWarrantyProviders) IsActive() bool {
	return dwp.ProviderStatus == "active" && time.Now().Before(dwp.EndDate)
}

func (dwp *DeviceWarrantyProviders) GetApprovalRate() float64 {
	if dwp.ClaimsProcessed > 0 {
		dwp.ApprovalRate = float64(dwp.ClaimsApproved) / float64(dwp.ClaimsProcessed) * 100
	}
	return dwp.ApprovalRate
}

func (dwp *DeviceWarrantyProviders) HasCoverageOverlap() bool {
	return dwp.CoverageOverlap && dwp.OverlappingProviders != "[]"
}

func (dwp *DeviceWarrantyProviders) IsProfitable() bool {
	return dwp.CostBenefitRatio > 1.0 && dwp.SavedAmount > dwp.PremiumCost
}

func (dwp *DeviceWarrantyProviders) HasOpenDisputes() bool {
	return dwp.OpenDisputes > 0
}

// Methods for DeviceServiceProviders
func (dsp *DeviceServiceProviders) IsAuthorized() bool {
	return dsp.AuthorizationStatus == "authorized" || dsp.AuthorizationStatus == "preferred"
}

func (dsp *DeviceServiceProviders) IsHighQuality() bool {
	return dsp.QualityScore >= 80 && dsp.AverageRating >= 4.0
}

func (dsp *DeviceServiceProviders) IsSLACompliant() bool {
	return dsp.SLAActive && dsp.SLACompliance >= 95
}

func (dsp *DeviceServiceProviders) GetSatisfactionRate() float64 {
	if dsp.TotalReviews > 0 {
		return float64(dsp.PositiveReviews) / float64(dsp.TotalReviews) * 100
	}
	return 0
}

func (dsp *DeviceServiceProviders) CanProvideHomeService() bool {
	return dsp.HomeService && dsp.ServiceRadius > 0
}

// Methods for DeviceInsurancePartners
func (dip *DeviceInsurancePartners) HasActiveCoverage() bool {
	return dip.CoverageStatus == "active" && time.Now().Before(dip.CoverageEnd)
}

func (dip *DeviceInsurancePartners) HasGapsIdentified() bool {
	return dip.GapsIdentified && dip.CoverageGaps != "[]"
}

func (dip *DeviceInsurancePartners) IsAPIConnected() bool {
	return dip.APIIntegrated && dip.IntegrationStatus == "connected"
}

func (dip *DeviceInsurancePartners) HasSubrogation() bool {
	return dip.SubrogationActive && dip.SubrogationAmount > 0
}

func (dip *DeviceInsurancePartners) IsReconciled() bool {
	return dip.ReconciliationStatus == "balanced"
}

// Methods for DeviceFinancialInstitutions
func (dfi *DeviceFinancialInstitutions) IsAccountActive() bool {
	return dfi.AccountStatus == "active" && dfi.AccountLinked
}

func (dfi *DeviceFinancialInstitutions) HasCreditAvailable() bool {
	return dfi.FinancingEnabled && dfi.AvailableCredit > 0
}

func (dfi *DeviceFinancialInstitutions) IsPaymentGatewayActive() bool {
	return dfi.GatewayActive && dfi.ProcessorStatus == "active"
}

func (dfi *DeviceFinancialInstitutions) GetChargebackRate() float64 {
	if dfi.TransactionsProcessed > 0 {
		dfi.ChargebackRate = float64(dfi.ChargebacksReceived) / float64(dfi.TransactionsProcessed) * 100
	}
	return dfi.ChargebackRate
}

func (dfi *DeviceFinancialInstitutions) HasFraudRisk() bool {
	return dfi.FraudScore > 70 || dfi.FraudAlerts > 5
}

// Methods for DeviceEcosystemIntegration
func (dei *DeviceEcosystemIntegration) IsSmartHomeConnected() bool {
	return dei.SmartHomeEnabled && dei.ConnectedDevices > 0
}

func (dei *DeviceEcosystemIntegration) HasHealthIntegration() bool {
	return dei.HealthDataSharing || dei.MedicalDataSync || dei.FitnessIntegration
}

func (dei *DeviceEcosystemIntegration) IsVehicleConnected() bool {
	return dei.VehicleConnected && dei.VehicleSystem != ""
}

func (dei *DeviceEcosystemIntegration) HasCloudBackup() bool {
	return dei.CloudSyncEnabled && dei.CloudStorageUsed > 0
}

func (dei *DeviceEcosystemIntegration) GetTotalConnectedDevices() int {
	return dei.ConnectedDevices + dei.WearablesConnected + dei.IoTDevicesLinked
}

func (dei *DeviceEcosystemIntegration) HasSyncIssues() bool {
	return dei.SyncConflicts > 0 || dei.DataConsistency < 95
}

// Methods for DeviceAPIIntegrations
func (dai *DeviceAPIIntegrations) IsConnected() bool {
	return dai.ConnectionStatus == "connected" && dai.HealthStatus == "healthy"
}

func (dai *DeviceAPIIntegrations) GetSuccessRate() float64 {
	if dai.TotalRequests > 0 {
		return float64(dai.SuccessfulRequests) / float64(dai.TotalRequests) * 100
	}
	return 0
}

func (dai *DeviceAPIIntegrations) IsRateLimited() bool {
	return dai.RateLimitEnabled && dai.CurrentUsage >= dai.RequestsPerMinute
}

func (dai *DeviceAPIIntegrations) HasHighErrorRate() bool {
	return dai.ErrorRate > 5 || dai.CircuitBreakerStatus == "open"
}

func (dai *DeviceAPIIntegrations) NeedsUpdate() bool {
	return dai.UpdateAvailable && dai.CurrentVersion != dai.LatestVersion
}

func (dai *DeviceAPIIntegrations) IsTokenExpired() bool {
	if dai.OAuthEnabled && dai.TokenExpiry != nil {
		return time.Since(*dai.TokenExpiry) > 0
	}
	return false
}

func (dai *DeviceAPIIntegrations) IsOverQuota() bool {
	return dai.UsageVsQuota >= 90
}

func (dai *DeviceAPIIntegrations) GetCostEfficiency() float64 {
	if dai.TotalRequests > 0 && dai.MonthlyCost > 0 {
		return float64(dai.TotalRequests) / dai.MonthlyCost
	}
	return 0
}
