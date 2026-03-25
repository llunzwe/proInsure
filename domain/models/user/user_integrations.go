package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserIntegrations manages third-party integrations and API access
type UserIntegrations struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// API Access
	APIEnabled     bool           `gorm:"default:false" json:"api_enabled"`
	APIKey         string         `gorm:"type:varchar(255)" json:"-"`
	APISecret      string         `gorm:"type:varchar(255)" json:"-"`
	APIKeyID       string         `gorm:"type:varchar(100);uniqueIndex" json:"api_key_id"`
	APIVersion     string         `gorm:"type:varchar(20)" json:"api_version"`
	RateLimits     map[string]int `gorm:"type:json" json:"rate_limits"`
	APIPermissions []string       `gorm:"type:json" json:"api_permissions"`
	LastAPICall    *time.Time     `json:"last_api_call"`
	APICallCount   int            `gorm:"default:0" json:"api_call_count"`

	// Connected Services
	ConnectedServices  []map[string]interface{} `gorm:"type:json" json:"connected_services"`
	ServiceCredentials map[string]string        `gorm:"type:json" json:"-"` // Encrypted
	ServiceTokens      map[string]interface{}   `gorm:"type:json" json:"-"` // Encrypted
	ServicePermissions map[string][]string      `gorm:"type:json" json:"service_permissions"`
	LastSyncTimes      map[string]time.Time     `gorm:"type:json" json:"last_sync_times"`
	SyncStatus         map[string]string        `gorm:"type:json" json:"sync_status"`

	// OAuth Connections
	OAuthProviders   []map[string]interface{} `gorm:"type:json" json:"oauth_providers"`
	OAuthTokens      map[string]interface{}   `gorm:"type:json" json:"-"` // Encrypted
	OAuthScopes      map[string][]string      `gorm:"type:json" json:"oauth_scopes"`
	TokenExpiryDates map[string]time.Time     `gorm:"type:json" json:"token_expiry_dates"`
	RefreshTokens    map[string]string        `gorm:"type:json" json:"-"` // Encrypted

	// Webhooks
	WebhookEndpoints    []map[string]interface{} `gorm:"type:json" json:"webhook_endpoints"`
	WebhookEvents       []string                 `gorm:"type:json" json:"webhook_events"`
	WebhookSecret       string                   `gorm:"type:varchar(255)" json:"-"`
	WebhookRetries      int                      `gorm:"default:3" json:"webhook_retries"`
	WebhookStatus       map[string]string        `gorm:"type:json" json:"webhook_status"`
	LastWebhookDelivery *time.Time               `json:"last_webhook_delivery"`

	// Integration Status
	IntegrationHealth  map[string]string        `gorm:"type:json" json:"integration_health"`
	FailedIntegrations []string                 `gorm:"type:json" json:"failed_integrations"`
	IntegrationErrors  []map[string]interface{} `gorm:"type:json" json:"integration_errors"`
	RecoveryAttempts   map[string]int           `gorm:"type:json" json:"recovery_attempts"`

	// Data Sync
	DataSyncEnabled    bool                   `gorm:"default:false" json:"data_sync_enabled"`
	SyncFrequency      string                 `gorm:"type:varchar(20)" json:"sync_frequency"`
	SyncDirection      string                 `gorm:"type:varchar(20)" json:"sync_direction"` // bidirectional/inbound/outbound
	DataMappings       map[string]interface{} `gorm:"type:json" json:"data_mappings"`
	ConflictResolution string                 `gorm:"type:varchar(50)" json:"conflict_resolution"`
	LastFullSync       *time.Time             `json:"last_full_sync"`
}

// UserAPIUsage tracks API usage and metrics
type UserAPIUsage struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Usage Metrics
	TotalRequests         int     `gorm:"default:0" json:"total_requests"`
	SuccessfulRequests    int     `gorm:"default:0" json:"successful_requests"`
	FailedRequests        int     `gorm:"default:0" json:"failed_requests"`
	AverageResponseTime   float64 `gorm:"default:0" json:"average_response_time_ms"`
	PeakRequestsPerMinute int     `gorm:"default:0" json:"peak_requests_per_minute"`

	// Endpoint Usage
	EndpointUsage          map[string]int `gorm:"type:json" json:"endpoint_usage"`
	MethodDistribution     map[string]int `gorm:"type:json" json:"method_distribution"`
	StatusCodeDistribution map[string]int `gorm:"type:json" json:"status_code_distribution"`
	ErrorTypes             map[string]int `gorm:"type:json" json:"error_types"`

	// Rate Limiting
	RateLimitHits int                  `gorm:"default:0" json:"rate_limit_hits"`
	QuotaUsage    map[string]int       `gorm:"type:json" json:"quota_usage"`
	QuotaResets   map[string]time.Time `gorm:"type:json" json:"quota_resets"`
	BurstCapacity int                  `gorm:"default:0" json:"burst_capacity"`

	// Performance
	AverageLatency map[string]float64 `gorm:"type:json" json:"average_latency"`
	P95Latency     map[string]float64 `gorm:"type:json" json:"p95_latency"`
	P99Latency     map[string]float64 `gorm:"type:json" json:"p99_latency"`
	Bandwidth      map[string]int     `gorm:"type:json" json:"bandwidth_bytes"`

	// Time-based Metrics
	HourlyUsage    map[string]int `gorm:"type:json" json:"hourly_usage"`
	DailyUsage     map[string]int `gorm:"type:json" json:"daily_usage"`
	MonthlyUsage   map[string]int `gorm:"type:json" json:"monthly_usage"`
	PeakUsageTimes []string       `gorm:"type:json" json:"peak_usage_times"`
}

// UserThirdPartyData manages data from external sources
type UserThirdPartyData struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Data Sources
	DataProviders   []map[string]interface{} `gorm:"type:json" json:"data_providers"`
	DataSourceTypes []string                 `gorm:"type:json" json:"data_source_types"`
	DataCategories  map[string][]string      `gorm:"type:json" json:"data_categories"`
	LastDataUpdate  map[string]time.Time     `gorm:"type:json" json:"last_data_update"`

	// Social Media Data
	SocialProfiles    map[string]interface{} `gorm:"type:json" json:"social_profiles"`
	SocialMetrics     map[string]interface{} `gorm:"type:json" json:"social_metrics"`
	SocialConnections map[string]int         `gorm:"type:json" json:"social_connections"`
	SocialActivity    map[string]interface{} `gorm:"type:json" json:"social_activity"`

	// Financial Data
	BankingConnections []map[string]interface{} `gorm:"type:json" json:"banking_connections"`
	TransactionData    map[string]interface{}   `gorm:"type:json" json:"transaction_data"`
	CreditBureauData   map[string]interface{}   `gorm:"type:json" json:"credit_bureau_data"`
	InvestmentAccounts []map[string]interface{} `gorm:"type:json" json:"investment_accounts"`

	// Identity Data
	IdentityProviders  []string               `gorm:"type:json" json:"identity_providers"`
	VerifiedIdentities map[string]bool        `gorm:"type:json" json:"verified_identities"`
	IdentityScores     map[string]float64     `gorm:"type:json" json:"identity_scores"`
	IdentityAttributes map[string]interface{} `gorm:"type:json" json:"identity_attributes"`

	// Behavioral Data
	BehavioralProviders []string               `gorm:"type:json" json:"behavioral_providers"`
	BehavioralScores    map[string]float64     `gorm:"type:json" json:"behavioral_scores"`
	BehavioralInsights  map[string]interface{} `gorm:"type:json" json:"behavioral_insights"`
	PredictiveModels    map[string]interface{} `gorm:"type:json" json:"predictive_models"`

	// Data Quality
	DataCompleteness map[string]float64 `gorm:"type:json" json:"data_completeness"`
	DataAccuracy     map[string]float64 `gorm:"type:json" json:"data_accuracy"`
	DataFreshness    map[string]string  `gorm:"type:json" json:"data_freshness"`
	ValidationStatus map[string]string  `gorm:"type:json" json:"validation_status"`
}

// UserDataExchange manages data import/export operations
type UserDataExchange struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Export Operations
	ExportRequests     []map[string]interface{} `gorm:"type:json" json:"export_requests"`
	LastExportDate     *time.Time               `json:"last_export_date"`
	ExportFormats      []string                 `gorm:"type:json" json:"export_formats"`
	ExportDestinations []string                 `gorm:"type:json" json:"export_destinations"`
	ScheduledExports   []map[string]interface{} `gorm:"type:json" json:"scheduled_exports"`

	// Import Operations
	ImportHistory    []map[string]interface{} `gorm:"type:json" json:"import_history"`
	LastImportDate   *time.Time               `json:"last_import_date"`
	ImportSources    []string                 `gorm:"type:json" json:"import_sources"`
	ImportMappings   map[string]interface{}   `gorm:"type:json" json:"import_mappings"`
	ImportValidation map[string]interface{}   `gorm:"type:json" json:"import_validation"`

	// Data Portability
	PortabilityRequests []map[string]interface{} `gorm:"type:json" json:"portability_requests"`
	DataPackages        []map[string]interface{} `gorm:"type:json" json:"data_packages"`
	TransferProtocols   []string                 `gorm:"type:json" json:"transfer_protocols"`
	DataFormats         map[string]string        `gorm:"type:json" json:"data_formats"`

	// Compliance
	DataResidency      string          `gorm:"type:varchar(50)" json:"data_residency"`
	ExportRestrictions []string        `gorm:"type:json" json:"export_restrictions"`
	ImportRestrictions []string        `gorm:"type:json" json:"import_restrictions"`
	ComplianceChecks   map[string]bool `gorm:"type:json" json:"compliance_checks"`
}

// UserPartnerIntegrations manages partner service integrations
type UserPartnerIntegrations struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Partner Services
	ActivePartners  []map[string]interface{} `gorm:"type:json" json:"active_partners"`
	PartnerAccounts map[string]string        `gorm:"type:json" json:"partner_accounts"`
	PartnerStatus   map[string]string        `gorm:"type:json" json:"partner_status"`
	PartnerBenefits map[string][]string      `gorm:"type:json" json:"partner_benefits"`

	// Service Usage
	ServiceUsage    map[string]interface{} `gorm:"type:json" json:"service_usage"`
	ServiceQuotas   map[string]int         `gorm:"type:json" json:"service_quotas"`
	ServiceExpiry   map[string]time.Time   `gorm:"type:json" json:"service_expiry"`
	ServiceRenewals map[string]bool        `gorm:"type:json" json:"service_renewals"`

	// Rewards & Benefits
	PartnerRewards  map[string]interface{}   `gorm:"type:json" json:"partner_rewards"`
	CrossPromotions []map[string]interface{} `gorm:"type:json" json:"cross_promotions"`
	BundleOffers    []map[string]interface{} `gorm:"type:json" json:"bundle_offers"`
	ExclusiveDeals  []map[string]interface{} `gorm:"type:json" json:"exclusive_deals"`

	// Data Sharing
	SharedDataTypes    map[string][]string      `gorm:"type:json" json:"shared_data_types"`
	DataSharingConsent map[string]bool          `gorm:"type:json" json:"data_sharing_consent"`
	DataSharingHistory []map[string]interface{} `gorm:"type:json" json:"data_sharing_history"`
	PrivacySettings    map[string]interface{}   `gorm:"type:json" json:"privacy_settings"`
}

// UserAutomationRules manages user-specific automation rules
type UserAutomationRules struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Automation Rules
	ActiveRules    []map[string]interface{} `gorm:"type:json" json:"active_rules"`
	RuleTriggers   map[string]interface{}   `gorm:"type:json" json:"rule_triggers"`
	RuleConditions map[string]interface{}   `gorm:"type:json" json:"rule_conditions"`
	RuleActions    map[string]interface{}   `gorm:"type:json" json:"rule_actions"`
	RuleSchedules  map[string]interface{}   `gorm:"type:json" json:"rule_schedules"`

	// Execution History
	ExecutionHistory  []map[string]interface{} `gorm:"type:json" json:"execution_history"`
	LastExecutionTime map[string]time.Time     `gorm:"type:json" json:"last_execution_time"`
	ExecutionCount    map[string]int           `gorm:"type:json" json:"execution_count"`
	SuccessRate       map[string]float64       `gorm:"type:json" json:"success_rate"`

	// Templates
	RuleTemplates   []map[string]interface{} `gorm:"type:json" json:"rule_templates"`
	CustomTemplates []map[string]interface{} `gorm:"type:json" json:"custom_templates"`
	SharedTemplates []map[string]interface{} `gorm:"type:json" json:"shared_templates"`

	// Notifications
	NotificationRules   map[string]interface{} `gorm:"type:json" json:"notification_rules"`
	AlertConfigurations map[string]interface{} `gorm:"type:json" json:"alert_configurations"`
	EscalationRules     map[string]interface{} `gorm:"type:json" json:"escalation_rules"`

	// Performance
	RulePerformance         map[string]interface{} `gorm:"type:json" json:"rule_performance"`
	OptimizationSuggestions []string               `gorm:"type:json" json:"optimization_suggestions"`
	ResourceUsage           map[string]interface{} `gorm:"type:json" json:"resource_usage"`
}

// TableName returns the table name
func (UserIntegrations) TableName() string {
	return "user_integrations"
}

// TableName returns the table name
func (UserAPIUsage) TableName() string {
	return "user_api_usage"
}

// TableName returns the table name
func (UserThirdPartyData) TableName() string {
	return "user_third_party_data"
}

// TableName returns the table name
func (UserDataExchange) TableName() string {
	return "user_data_exchange"
}

// TableName returns the table name
func (UserPartnerIntegrations) TableName() string {
	return "user_partner_integrations"
}

// TableName returns the table name
func (UserAutomationRules) TableName() string {
	return "user_automation_rules"
}
