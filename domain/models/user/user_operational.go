package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserPerformanceMetrics tracks operational performance metrics
type UserPerformanceMetrics struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Response Metrics
	AverageResponseTime    int     `json:"average_response_time_minutes"`
	FirstResponseTime      int     `json:"first_response_time_minutes"`
	ResolutionTime         int     `json:"resolution_time_hours"`
	EscalationRate         float64 `gorm:"default:0" json:"escalation_rate"`
	FirstContactResolution float64 `gorm:"default:0" json:"first_contact_resolution"`

	// Satisfaction Metrics
	CustomerSatisfaction float64 `gorm:"default:0" json:"customer_satisfaction"`
	NetPromoterScore     float64 `gorm:"default:0" json:"net_promoter_score"`
	CustomerEffortScore  float64 `gorm:"default:0" json:"customer_effort_score"`
	QualityScore         float64 `gorm:"default:0" json:"quality_score"`

	// Efficiency Metrics
	ProcessingEfficiency float64 `gorm:"default:0" json:"processing_efficiency"`
	AutomationRate       float64 `gorm:"default:0" json:"automation_rate"`
	ErrorRate            float64 `gorm:"default:0" json:"error_rate"`
	ReworkRate           float64 `gorm:"default:0" json:"rework_rate"`

	// Volume Metrics
	TransactionVolume   int `gorm:"default:0" json:"transaction_volume"`
	InteractionCount    int `gorm:"default:0" json:"interaction_count"`
	ServiceRequestCount int `gorm:"default:0" json:"service_request_count"`
	ComplaintCount      int `gorm:"default:0" json:"complaint_count"`

	// Cost Metrics
	ServiceCost     decimal.Decimal `gorm:"type:decimal(10,2)" json:"service_cost"`
	AcquisitionCost decimal.Decimal `gorm:"type:decimal(10,2)" json:"acquisition_cost"`
	RetentionCost   decimal.Decimal `gorm:"type:decimal(10,2)" json:"retention_cost"`
	SupportCost     decimal.Decimal `gorm:"type:decimal(10,2)" json:"support_cost"`

	// Benchmarking
	IndustryBenchmark map[string]float64 `gorm:"type:json" json:"industry_benchmark"`
	PeerComparison    map[string]float64 `gorm:"type:json" json:"peer_comparison"`
	PerformanceTrend  string             `gorm:"type:varchar(20)" json:"performance_trend"`
	ImprovementAreas  []string           `gorm:"type:json" json:"improvement_areas"`
}

// UserQualityAssurance manages quality and service standards
type UserQualityAssurance struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Quality Scores
	OverallQualityScore float64 `gorm:"default:0" json:"overall_quality_score"`
	ServiceQualityScore float64 `gorm:"default:0" json:"service_quality_score"`
	DataQualityScore    float64 `gorm:"default:0" json:"data_quality_score"`
	ProcessQualityScore float64 `gorm:"default:0" json:"process_quality_score"`

	// Quality Checks
	QualityChecksPassed int                      `gorm:"default:0" json:"quality_checks_passed"`
	QualityChecksFailed int                      `gorm:"default:0" json:"quality_checks_failed"`
	LastQualityCheck    *time.Time               `json:"last_quality_check"`
	QualityIssues       []map[string]interface{} `gorm:"type:json" json:"quality_issues"`

	// Standards Compliance
	StandardsMet        []string        `gorm:"type:json" json:"standards_met"`
	StandardsViolated   []string        `gorm:"type:json" json:"standards_violated"`
	ComplianceLevel     float64         `gorm:"default:0" json:"compliance_level"`
	CertificationStatus map[string]bool `gorm:"type:json" json:"certification_status"`

	// Audits
	AuditHistory  []map[string]interface{} `gorm:"type:json" json:"audit_history"`
	LastAuditDate *time.Time               `json:"last_audit_date"`
	AuditScore    float64                  `gorm:"default:0" json:"audit_score"`
	AuditFindings []string                 `gorm:"type:json" json:"audit_findings"`

	// Improvements
	ImprovementPlan   []map[string]interface{} `gorm:"type:json" json:"improvement_plan"`
	CorrectiveActions []map[string]interface{} `gorm:"type:json" json:"corrective_actions"`
	PreventiveActions []map[string]interface{} `gorm:"type:json" json:"preventive_actions"`
	TrainingRequired  []string                 `gorm:"type:json" json:"training_required"`
}

// UserCapacityManagement manages service limits and usage
type UserCapacityManagement struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Service Limits
	ServiceTier      string `gorm:"type:varchar(20)" json:"service_tier"`
	TransactionLimit int    `gorm:"default:0" json:"transaction_limit"`
	StorageLimit     int    `gorm:"default:0" json:"storage_limit_gb"`
	BandwidthLimit   int    `gorm:"default:0" json:"bandwidth_limit_gb"`
	APICallLimit     int    `gorm:"default:0" json:"api_call_limit"`
	DeviceLimit      int    `gorm:"default:0" json:"device_limit"`
	PolicyLimit      int    `gorm:"default:0" json:"policy_limit"`

	// Current Usage
	TransactionUsage int     `gorm:"default:0" json:"transaction_usage"`
	StorageUsage     float64 `gorm:"default:0" json:"storage_usage_gb"`
	BandwidthUsage   float64 `gorm:"default:0" json:"bandwidth_usage_gb"`
	APICallUsage     int     `gorm:"default:0" json:"api_call_usage"`
	DeviceUsage      int     `gorm:"default:0" json:"device_usage"`
	PolicyUsage      int     `gorm:"default:0" json:"policy_usage"`

	// Usage Patterns
	PeakUsageTimes      []string               `gorm:"type:json" json:"peak_usage_times"`
	UsageTrends         map[string]interface{} `gorm:"type:json" json:"usage_trends"`
	ResourceUtilization float64                `gorm:"default:0" json:"resource_utilization"`
	CapacityWarnings    []string               `gorm:"type:json" json:"capacity_warnings"`

	// Quotas
	QuotaResets          map[string]time.Time `gorm:"type:json" json:"quota_resets"`
	QuotaOverages        map[string]int       `gorm:"type:json" json:"quota_overages"`
	OverageCharges       decimal.Decimal      `gorm:"type:decimal(10,2)" json:"overage_charges"`
	QuotaUpgradeEligible bool                 `gorm:"default:false" json:"quota_upgrade_eligible"`
}

// UserScheduling manages appointments and availability
type UserScheduling struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Availability
	AvailabilitySchedule map[string]interface{} `gorm:"type:json" json:"availability_schedule"`
	TimeZone             string                 `gorm:"type:varchar(50)" json:"time_zone"`
	WorkingHours         map[string]string      `gorm:"type:json" json:"working_hours"`
	BlackoutDates        []time.Time            `gorm:"type:json" json:"blackout_dates"`
	PreferredTimes       []string               `gorm:"type:json" json:"preferred_times"`

	// Appointments
	ScheduledAppointments []map[string]interface{} `gorm:"type:json" json:"scheduled_appointments"`
	RecurringAppointments []map[string]interface{} `gorm:"type:json" json:"recurring_appointments"`
	AppointmentHistory    []map[string]interface{} `gorm:"type:json" json:"appointment_history"`
	CancelledAppointments []map[string]interface{} `gorm:"type:json" json:"cancelled_appointments"`
	NoShowCount           int                      `gorm:"default:0" json:"no_show_count"`

	// Reminders
	ReminderSettings map[string]interface{} `gorm:"type:json" json:"reminder_settings"`
	ReminderChannels []string               `gorm:"type:json" json:"reminder_channels"`
	ReminderTiming   map[string]int         `gorm:"type:json" json:"reminder_timing_hours"`

	// Calendar Integration
	CalendarSync      bool              `gorm:"default:false" json:"calendar_sync"`
	CalendarProviders []string          `gorm:"type:json" json:"calendar_providers"`
	SyncedCalendars   map[string]string `gorm:"type:json" json:"synced_calendars"`
	LastSyncTime      *time.Time        `json:"last_sync_time"`

	// Booking Preferences
	BookingBuffer      int    `json:"booking_buffer_minutes"`
	MaxAdvanceBooking  int    `json:"max_advance_booking_days"`
	MinAdvanceBooking  int    `json:"min_advance_booking_hours"`
	AllowRescheduling  bool   `gorm:"default:true" json:"allow_rescheduling"`
	ReschedulingPolicy string `gorm:"type:text" json:"rescheduling_policy"`
}

// UserIncidentManagement tracks incidents and resolutions
type UserIncidentManagement struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Incident Details
	IncidentID     string     `gorm:"type:varchar(100);uniqueIndex" json:"incident_id"`
	IncidentType   string     `gorm:"type:varchar(50)" json:"incident_type"`
	Severity       string     `gorm:"type:varchar(20)" json:"severity"`
	Status         string     `gorm:"type:varchar(20)" json:"status"`
	ReportedDate   time.Time  `json:"reported_date"`
	ResolvedDate   *time.Time `json:"resolved_date"`
	ResolutionTime int        `json:"resolution_time_hours"`

	// Impact
	ImpactLevel      string          `gorm:"type:varchar(20)" json:"impact_level"`
	AffectedServices []string        `gorm:"type:json" json:"affected_services"`
	UserImpact       string          `gorm:"type:text" json:"user_impact"`
	BusinessImpact   string          `gorm:"type:text" json:"business_impact"`
	FinancialImpact  decimal.Decimal `gorm:"type:decimal(10,2)" json:"financial_impact"`

	// Resolution
	ResolutionSteps    []map[string]interface{} `gorm:"type:json" json:"resolution_steps"`
	RootCause          string                   `gorm:"type:text" json:"root_cause"`
	CorrectiveActions  []string                 `gorm:"type:json" json:"corrective_actions"`
	PreventiveMeasures []string                 `gorm:"type:json" json:"preventive_measures"`

	// Communication
	StakeholdersNotified []uuid.UUID              `gorm:"type:json" json:"stakeholders_notified"`
	CommunicationLog     []map[string]interface{} `gorm:"type:json" json:"communication_log"`
	PublicStatement      string                   `gorm:"type:text" json:"public_statement"`

	// Follow-up
	PostIncidentReview bool                     `gorm:"default:false" json:"post_incident_review"`
	ReviewDate         *time.Time               `json:"review_date"`
	LessonsLearned     []string                 `gorm:"type:json" json:"lessons_learned"`
	ActionItems        []map[string]interface{} `gorm:"type:json" json:"action_items"`
}

// UserServiceCatalog manages available services for users
type UserServiceCatalog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Available Services
	AvailableServices []map[string]interface{} `gorm:"type:json" json:"available_services"`
	EnabledServices   []string                 `gorm:"type:json" json:"enabled_services"`
	DisabledServices  []string                 `gorm:"type:json" json:"disabled_services"`
	PendingServices   []string                 `gorm:"type:json" json:"pending_services"`

	// Service Subscriptions
	ActiveSubscriptions  []map[string]interface{} `gorm:"type:json" json:"active_subscriptions"`
	ExpiredSubscriptions []map[string]interface{} `gorm:"type:json" json:"expired_subscriptions"`
	SubscriptionTiers    map[string]string        `gorm:"type:json" json:"subscription_tiers"`
	RenewalDates         map[string]time.Time     `gorm:"type:json" json:"renewal_dates"`

	// Add-ons
	ActiveAddons    []map[string]interface{}   `gorm:"type:json" json:"active_addons"`
	AvailableAddons []map[string]interface{}   `gorm:"type:json" json:"available_addons"`
	AddonUsage      map[string]int             `gorm:"type:json" json:"addon_usage"`
	AddonCosts      map[string]decimal.Decimal `gorm:"type:json" json:"addon_costs"`

	// Service Requests
	ServiceRequests  []map[string]interface{} `gorm:"type:json" json:"service_requests"`
	PendingRequests  []map[string]interface{} `gorm:"type:json" json:"pending_requests"`
	ApprovedRequests []map[string]interface{} `gorm:"type:json" json:"approved_requests"`
	DeniedRequests   []map[string]interface{} `gorm:"type:json" json:"denied_requests"`

	// Custom Services
	CustomServices        []map[string]interface{} `gorm:"type:json" json:"custom_services"`
	ServiceConfigurations map[string]interface{}   `gorm:"type:json" json:"service_configurations"`
	ServiceDependencies   map[string][]string      `gorm:"type:json" json:"service_dependencies"`

	// Pricing
	ServicePricing     map[string]decimal.Decimal `gorm:"type:json" json:"service_pricing"`
	DiscountedServices map[string]decimal.Decimal `gorm:"type:json" json:"discounted_services"`
	BundledServices    []map[string]interface{}   `gorm:"type:json" json:"bundled_services"`
	TotalServiceCost   decimal.Decimal            `gorm:"type:decimal(10,2)" json:"total_service_cost"`
}

// UserMultiCurrency manages multi-currency support
type UserMultiCurrency struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Currency Settings
	PrimaryCurrency     string   `gorm:"type:varchar(3)" json:"primary_currency"`
	SecondaryCurrencies []string `gorm:"type:json" json:"secondary_currencies"`
	PreferredCurrency   string   `gorm:"type:varchar(3)" json:"preferred_currency"`
	DisplayCurrency     string   `gorm:"type:varchar(3)" json:"display_currency"`

	// Exchange Rates
	ExchangeRates  map[string]float64 `gorm:"type:json" json:"exchange_rates"`
	LastRateUpdate *time.Time         `json:"last_rate_update"`
	RateProvider   string             `gorm:"type:varchar(50)" json:"rate_provider"`
	CustomRates    map[string]float64 `gorm:"type:json" json:"custom_rates"`

	// Currency Balances
	CurrencyBalances  map[string]decimal.Decimal `gorm:"type:json" json:"currency_balances"`
	HoldAmounts       map[string]decimal.Decimal `gorm:"type:json" json:"hold_amounts"`
	AvailableBalances map[string]decimal.Decimal `gorm:"type:json" json:"available_balances"`

	// Conversion History
	ConversionHistory []map[string]interface{}   `gorm:"type:json" json:"conversion_history"`
	TotalConverted    map[string]decimal.Decimal `gorm:"type:json" json:"total_converted"`
	ConversionFees    decimal.Decimal            `gorm:"type:decimal(10,2)" json:"conversion_fees"`

	// International Support
	InternationalEnabled bool                   `gorm:"default:false" json:"international_enabled"`
	SupportedCountries   []string               `gorm:"type:json" json:"supported_countries"`
	RestrictedCountries  []string               `gorm:"type:json" json:"restricted_countries"`
	RegionalPricing      map[string]interface{} `gorm:"type:json" json:"regional_pricing"`

	// Settlement
	SettlementCurrency  string          `gorm:"type:varchar(3)" json:"settlement_currency"`
	AutoConversion      bool            `gorm:"default:false" json:"auto_conversion"`
	ConversionThreshold decimal.Decimal `gorm:"type:decimal(10,2)" json:"conversion_threshold"`
	ConversionStrategy  string          `gorm:"type:varchar(50)" json:"conversion_strategy"`
}

// UserEnvironmentalImpact tracks sustainability metrics
type UserEnvironmentalImpact struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Carbon Footprint
	CarbonFootprint       float64 `gorm:"default:0" json:"carbon_footprint_kg"`
	CarbonOffset          float64 `gorm:"default:0" json:"carbon_offset_kg"`
	NetCarbon             float64 `gorm:"default:0" json:"net_carbon_kg"`
	CarbonReductionTarget float64 `gorm:"default:0" json:"carbon_reduction_target"`

	// Sustainability Score
	SustainabilityScore float64 `gorm:"default:0" json:"sustainability_score"`
	EnvironmentalRating string  `gorm:"type:varchar(20)" json:"environmental_rating"`
	ImpactCategory      string  `gorm:"type:varchar(50)" json:"impact_category"`
	ImprovementTrend    string  `gorm:"type:varchar(20)" json:"improvement_trend"`

	// Green Initiatives
	ParticipatingPrograms []string                 `gorm:"type:json" json:"participating_programs"`
	GreenCertifications   []string                 `gorm:"type:json" json:"green_certifications"`
	EcoFriendlyActions    []map[string]interface{} `gorm:"type:json" json:"eco_friendly_actions"`
	TreesPlanted          int                      `gorm:"default:0" json:"trees_planted"`

	// Resource Usage
	PaperlessBilling  bool    `gorm:"default:false" json:"paperless_billing"`
	DigitalDocuments  int     `gorm:"default:0" json:"digital_documents"`
	PaperSaved        float64 `gorm:"default:0" json:"paper_saved_kg"`
	EnergyConsumption float64 `gorm:"default:0" json:"energy_consumption_kwh"`

	// Recycling & Disposal
	RecycledDevices      int     `gorm:"default:0" json:"recycled_devices"`
	ProperDisposal       bool    `gorm:"default:false" json:"proper_disposal"`
	RecyclingCredits     int     `gorm:"default:0" json:"recycling_credits"`
	CircularEconomyScore float64 `gorm:"default:0" json:"circular_economy_score"`

	// Rewards & Incentives
	GreenRewards         []map[string]interface{} `gorm:"type:json" json:"green_rewards"`
	EcoPoints            int                      `gorm:"default:0" json:"eco_points"`
	SustainabilityBadges []string                 `gorm:"type:json" json:"sustainability_badges"`
	CarbonCreditBalance  float64                  `gorm:"default:0" json:"carbon_credit_balance"`
}

// TableName returns the table name
func (UserPerformanceMetrics) TableName() string {
	return "user_performance_metrics"
}

// TableName returns the table name
func (UserQualityAssurance) TableName() string {
	return "user_quality_assurance"
}

// TableName returns the table name
func (UserCapacityManagement) TableName() string {
	return "user_capacity_management"
}

// TableName returns the table name
func (UserScheduling) TableName() string {
	return "user_scheduling"
}

// TableName returns the table name
func (UserIncidentManagement) TableName() string {
	return "user_incident_management"
}

// TableName returns the table name
func (UserServiceCatalog) TableName() string {
	return "user_service_catalog"
}

// TableName returns the table name
func (UserMultiCurrency) TableName() string {
	return "user_multi_currency"
}

// TableName returns the table name
func (UserEnvironmentalImpact) TableName() string {
	return "user_environmental_impact"
}
