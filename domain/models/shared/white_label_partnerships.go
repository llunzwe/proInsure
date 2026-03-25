package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// TenantOrganization represents white-label tenant organizations
type TenantOrganization struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TenantCode        string         `gorm:"uniqueIndex;not null" json:"tenant_code"`
	OrganizationName  string         `gorm:"not null" json:"organization_name"`
	LegalName         string         `gorm:"not null" json:"legal_name"`
	OrganizationType  string         `gorm:"not null" json:"organization_type"` // telco, bank, oem, broker, direct
	BusinessLicense   string         `gorm:"uniqueIndex;not null" json:"business_license"`
	TaxIdentifier     string         `gorm:"uniqueIndex;not null" json:"tax_identifier"`
	Country           string         `gorm:"not null" json:"country"`
	Region            string         `json:"region"`
	PrimaryContact    string         `gorm:"not null" json:"primary_contact"`
	ContactEmail      string         `gorm:"not null" json:"contact_email"`
	ContactPhone      string         `gorm:"not null" json:"contact_phone"`
	BillingAddress    string         `json:"billing_address"` // JSON object
	TechnicalContact  string         `json:"technical_contact"`
	TechnicalEmail    string         `json:"technical_email"`
	ContractStartDate time.Time      `gorm:"not null" json:"contract_start_date"`
	ContractEndDate   time.Time      `gorm:"not null" json:"contract_end_date"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	IsSuspended       bool           `gorm:"default:false" json:"is_suspended"`
	SuspensionReason  string         `json:"suspension_reason"`
	TierLevel         string         `gorm:"default:'basic'" json:"tier_level"` // basic, premium, enterprise
	MaxUsers          int            `gorm:"default:1000" json:"max_users"`
	MaxPolicies       int            `gorm:"default:10000" json:"max_policies"`
	MaxClaims         int            `gorm:"default:1000" json:"max_claims"`
	CurrentUsers      int            `gorm:"default:0" json:"current_users"`
	CurrentPolicies   int            `gorm:"default:0" json:"current_policies"`
	CurrentClaims     int            `gorm:"default:0" json:"current_claims"`
	MonthlyRevenue    float64        `gorm:"default:0" json:"monthly_revenue"`
	CommissionRate    float64        `gorm:"default:0.15" json:"commission_rate"` // 15% default
	CustomDomain      string         `json:"custom_domain"`
	BrandingConfig    string         `json:"branding_config"`                         // JSON object with colors, logos, etc.
	FeatureFlags      string         `json:"feature_flags"`                           // JSON object with enabled features
	APILimits         string         `json:"api_limits"`                              // JSON object with rate limits
	DataRetentionDays int            `gorm:"default:2555" json:"data_retention_days"` // 7 years default
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Configurations []TenantConfiguration `gorm:"foreignKey:TenantID" json:"configurations,omitempty"`
	Users          []TenantUser          `gorm:"foreignKey:TenantID" json:"users,omitempty"`
	APIKeys        []TenantAPIKey        `gorm:"foreignKey:TenantID" json:"api_keys,omitempty"`
	Billing        []TenantBilling       `gorm:"foreignKey:TenantID" json:"billing,omitempty"`
}

// TenantConfiguration represents tenant-specific configurations
type TenantConfiguration struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	TenantID       uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	ConfigCategory string    `gorm:"not null" json:"config_category"` // branding, features, limits, integrations
	ConfigKey      string    `gorm:"not null" json:"config_key"`
	ConfigValue    string    `json:"config_value"`
	DataType       string    `gorm:"default:'string'" json:"data_type"` // string, number, boolean, json, array
	IsEncrypted    bool      `gorm:"default:false" json:"is_encrypted"`
	Description    string    `json:"description"`
	IsEditable     bool      `gorm:"default:true" json:"is_editable"`
	ValidValues    string    `json:"valid_values"` // JSON array of allowed values
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Tenant TenantOrganization `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// TenantUser represents users within tenant organizations
type TenantUser struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TenantID    uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	TenantRole  string     `gorm:"not null" json:"tenant_role"` // admin, manager, agent, viewer
	Permissions string     `json:"permissions"`                 // JSON array of specific permissions
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	JoinedAt    time.Time  `gorm:"not null" json:"joined_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Tenant TenantOrganization `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   models.User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TenantAPIKey represents API keys for tenant organizations
type TenantAPIKey struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TenantID    uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	KeyName     string     `gorm:"not null" json:"key_name"`
	APIKey      string     `gorm:"uniqueIndex;not null" json:"api_key"`
	APISecret   string     `json:"api_secret"`                     // encrypted
	KeyType     string     `gorm:"not null" json:"key_type"`       // production, sandbox, test
	Permissions string     `json:"permissions"`                    // JSON array of API permissions
	RateLimit   int        `gorm:"default:1000" json:"rate_limit"` // requests per hour
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	UsageCount  int64      `gorm:"default:0" json:"usage_count"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships - temporarily commented for migration
	// Tenant        TenantOrganization `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	// CreatedByUser User               `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

// TenantBilling represents billing for tenant organizations
type TenantBilling struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TenantID           uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	BillingPeriodStart time.Time  `gorm:"not null" json:"billing_period_start"`
	BillingPeriodEnd   time.Time  `gorm:"not null" json:"billing_period_end"`
	BaseFee            float64    `gorm:"not null" json:"base_fee"`
	UsageFees          float64    `gorm:"default:0" json:"usage_fees"`
	CommissionEarned   float64    `gorm:"default:0" json:"commission_earned"`
	TotalAmount        float64    `gorm:"not null" json:"total_amount"`
	Currency           string     `gorm:"default:'USD'" json:"currency"`
	Status             string     `gorm:"not null;default:'pending'" json:"status"` // pending, sent, paid, overdue, cancelled
	InvoiceNumber      string     `gorm:"uniqueIndex" json:"invoice_number"`
	InvoiceDate        time.Time  `json:"invoice_date"`
	DueDate            time.Time  `json:"due_date"`
	PaidDate           *time.Time `json:"paid_date"`
	PaymentMethod      string     `json:"payment_method"`
	PaymentReference   string     `json:"payment_reference"`
	UsageMetrics       string     `json:"usage_metrics"` // JSON object with detailed usage
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Tenant TenantOrganization `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// PartnerIntegration represents external partner integrations
type PartnerIntegration struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PartnerCode         string         `gorm:"uniqueIndex;not null" json:"partner_code"`
	PartnerName         string         `gorm:"not null" json:"partner_name"`
	PartnerType         string         `gorm:"not null" json:"partner_type"`     // payment_gateway, repair_network, logistics, data_provider
	IntegrationType     string         `gorm:"not null" json:"integration_type"` // api, webhook, file_transfer, manual
	ContactInfo         string         `json:"contact_info"`                     // JSON object
	ContractDetails     string         `json:"contract_details"`                 // JSON object
	SLARequirements     string         `json:"sla_requirements"`                 // JSON object
	APIEndpoint         string         `json:"api_endpoint"`
	AuthMethod          string         `json:"auth_method"` // api_key, oauth, basic_auth, certificate
	Credentials         string         `json:"credentials"` // encrypted credentials
	IsActive            bool           `gorm:"default:true" json:"is_active"`
	IsTestMode          bool           `gorm:"default:false" json:"is_test_mode"`
	LastHealthCheck     *time.Time     `json:"last_health_check"`
	HealthStatus        string         `gorm:"default:'unknown'" json:"health_status"` // healthy, degraded, down, unknown
	ErrorCount          int            `gorm:"default:0" json:"error_count"`
	SuccessRate         float64        `gorm:"default:0" json:"success_rate"`          // percentage
	AverageResponseTime int            `gorm:"default:0" json:"average_response_time"` // milliseconds
	MonthlyVolume       int64          `gorm:"default:0" json:"monthly_volume"`
	RevenuShare         float64        `json:"revenue_share"` // percentage
	CostPerTransaction  float64        `json:"cost_per_transaction"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Transactions []PartnerTransaction `gorm:"foreignKey:PartnerID" json:"transactions,omitempty"`
	HealthChecks []PartnerHealthCheck `gorm:"foreignKey:PartnerID" json:"health_checks,omitempty"`
}

// PartnerTransaction represents transactions with partners
type PartnerTransaction struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	PartnerID       uuid.UUID  `gorm:"type:uuid;not null" json:"partner_id"`
	TransactionID   string     `gorm:"uniqueIndex;not null" json:"transaction_id"`
	TransactionType string     `gorm:"not null" json:"transaction_type"` // payment, repair, delivery, data_request
	ReferenceID     *uuid.UUID `gorm:"type:uuid" json:"reference_id"`    // policy, claim, etc.
	ReferenceType   string     `json:"reference_type"`
	Amount          float64    `json:"amount"`
	Currency        string     `gorm:"default:'USD'" json:"currency"`
	Status          string     `gorm:"not null;default:'pending'" json:"status"` // pending, processing, completed, failed, cancelled
	RequestData     string     `json:"request_data"`                             // JSON object
	ResponseData    string     `json:"response_data"`                            // JSON object
	ResponseTime    int        `json:"response_time"`                            // milliseconds
	ErrorMessage    string     `json:"error_message"`
	RetryCount      int        `gorm:"default:0" json:"retry_count"`
	MaxRetries      int        `gorm:"default:3" json:"max_retries"`
	ProcessedAt     *time.Time `json:"processed_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Partner PartnerIntegration `gorm:"foreignKey:PartnerID" json:"partner,omitempty"`
}

// PartnerHealthCheck represents partner integration health monitoring
type PartnerHealthCheck struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PartnerID    uuid.UUID `gorm:"type:uuid;not null" json:"partner_id"`
	CheckDate    time.Time `gorm:"not null" json:"check_date"`
	CheckType    string    `gorm:"not null" json:"check_type"` // ping, api_call, full_test
	ResponseTime int       `json:"response_time"`              // milliseconds
	Status       string    `gorm:"not null" json:"status"`     // healthy, degraded, down
	StatusCode   int       `json:"status_code"`
	ErrorMessage string    `json:"error_message"`
	CheckDetails string    `json:"check_details"` // JSON object
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Partner PartnerIntegration `gorm:"foreignKey:PartnerID" json:"partner,omitempty"`
}

// APIUsageAnalytics represents API usage analytics
type APIUsageAnalytics struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TenantID            *uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	APIKeyID            *uuid.UUID `gorm:"type:uuid" json:"api_key_id"`
	AnalyticsDate       time.Time  `gorm:"not null" json:"analytics_date"`
	Endpoint            string     `gorm:"not null" json:"endpoint"`
	Method              string     `gorm:"not null" json:"method"`
	TotalRequests       int64      `gorm:"default:0" json:"total_requests"`
	SuccessfulRequests  int64      `gorm:"default:0" json:"successful_requests"`
	FailedRequests      int64      `gorm:"default:0" json:"failed_requests"`
	AverageResponseTime int        `gorm:"default:0" json:"average_response_time"` // milliseconds
	MaxResponseTime     int        `gorm:"default:0" json:"max_response_time"`
	MinResponseTime     int        `gorm:"default:0" json:"min_response_time"`
	DataTransferred     int64      `gorm:"default:0" json:"data_transferred"` // bytes
	ErrorCodes          string     `json:"error_codes"`                       // JSON object with error code counts
	TopUserAgents       string     `json:"top_user_agents"`                   // JSON array
	TopIPAddresses      string     `json:"top_ip_addresses"`                  // JSON array
	GeographicData      string     `json:"geographic_data"`                   // JSON object
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Tenant *TenantOrganization `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	APIKey *TenantAPIKey       `gorm:"foreignKey:APIKeyID" json:"api_key,omitempty"`
}

// MarketplaceProduct represents products in the insurance marketplace
type MarketplaceProduct struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ProductCode            string         `gorm:"uniqueIndex;not null" json:"product_code"`
	ProductName            string         `gorm:"not null" json:"product_name"`
	ProductType            string         `gorm:"not null" json:"product_type"` // device_insurance, extended_warranty, cyber_protection
	Category               string         `gorm:"not null" json:"category"`
	Description            string         `json:"description"`
	Features               string         `json:"features"`      // JSON array
	Coverage               string         `json:"coverage"`      // JSON object
	Pricing                string         `json:"pricing"`       // JSON object with pricing tiers
	TargetMarket           string         `json:"target_market"` // JSON array of target markets
	MinimumAge             int            `json:"minimum_age"`
	MaximumAge             int            `json:"maximum_age"`
	GeographicRestrictions string         `json:"geographic_restrictions"` // JSON array
	IsActive               bool           `gorm:"default:true" json:"is_active"`
	IsPromoted             bool           `gorm:"default:false" json:"is_promoted"`
	LaunchDate             time.Time      `json:"launch_date"`
	EndOfSaleDate          *time.Time     `json:"end_of_sale_date"`
	TotalSales             int64          `gorm:"default:0" json:"total_sales"`
	MonthlyRevenue         float64        `gorm:"default:0" json:"monthly_revenue"`
	Rating                 float64        `gorm:"default:0" json:"rating"` // 1-5 stars
	ReviewCount            int            `gorm:"default:0" json:"review_count"`
	CreatedBy              uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - temporarily commented for migration
	// CreatedByUser User              `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	// Sales         []MarketplaceSale `gorm:"foreignKey:ProductID" json:"sales,omitempty"`
	// Reviews       []ProductReview   `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
}

// MarketplaceSale represents sales through the marketplace
type MarketplaceSale struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ProductID        uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	TenantID         *uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	CustomerID       uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	SaleDate         time.Time  `gorm:"not null" json:"sale_date"`
	SaleAmount       float64    `gorm:"not null" json:"sale_amount"`
	Currency         string     `gorm:"default:'USD'" json:"currency"`
	CommissionRate   float64    `json:"commission_rate"`
	CommissionAmount float64    `json:"commission_amount"`
	PaymentStatus    string     `gorm:"not null;default:'pending'" json:"payment_status"` // pending, paid, failed, refunded
	PaymentMethod    string     `json:"payment_method"`
	ReferralSource   string     `json:"referral_source"`
	SalesChannel     string     `json:"sales_channel"` // web, mobile, api, partner
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Product  MarketplaceProduct  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Tenant   *TenantOrganization `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Customer models.User         `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

// ProductReview represents customer reviews for marketplace products
type ProductReview struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ProductID       uuid.UUID      `gorm:"type:uuid;not null" json:"product_id"`
	CustomerID      uuid.UUID      `gorm:"type:uuid;not null" json:"customer_id"`
	Rating          int            `gorm:"not null" json:"rating"` // 1-5 stars
	Title           string         `json:"title"`
	ReviewText      string         `json:"review_text"`
	Pros            string         `json:"pros"`
	Cons            string         `json:"cons"`
	WouldRecommend  bool           `gorm:"default:true" json:"would_recommend"`
	IsVerified      bool           `gorm:"default:false" json:"is_verified"`
	IsModerated     bool           `gorm:"default:false" json:"is_moderated"`
	ModeratedBy     *uuid.UUID     `gorm:"type:uuid" json:"moderated_by"`
	ModerationNotes string         `json:"moderation_notes"`
	HelpfulVotes    int            `gorm:"default:0" json:"helpful_votes"`
	TotalVotes      int            `gorm:"default:0" json:"total_votes"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Product         MarketplaceProduct `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Customer        models.User        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ModeratedByUser *models.User       `gorm:"foreignKey:ModeratedBy" json:"moderated_by_user,omitempty"`
}

// TableName methods
func (TenantOrganization) TableName() string {
	return "tenant_organizations"
}

func (TenantConfiguration) TableName() string {
	return "tenant_configurations"
}

func (TenantUser) TableName() string {
	return "tenant_users"
}

func (TenantAPIKey) TableName() string {
	return "tenant_api_keys"
}

func (TenantBilling) TableName() string {
	return "tenant_billing"
}

func (PartnerIntegration) TableName() string {
	return "partner_integrations"
}

func (PartnerTransaction) TableName() string {
	return "partner_transactions"
}

func (PartnerHealthCheck) TableName() string {
	return "partner_health_checks"
}

func (APIUsageAnalytics) TableName() string {
	return "api_usage_analytics"
}

func (MarketplaceProduct) TableName() string {
	return "marketplace_products"
}

func (MarketplaceSale) TableName() string {
	return "marketplace_sales"
}

func (ProductReview) TableName() string {
	return "product_reviews"
}

// BeforeCreate hooks
func (to *TenantOrganization) BeforeCreate(tx *gorm.DB) error {
	if to.ID == uuid.Nil {
		to.ID = uuid.New()
	}
	if to.TenantCode == "" {
		to.TenantCode = "TNT-" + uuid.New().String()[:8]
	}
	return nil
}

func (tc *TenantConfiguration) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == uuid.Nil {
		tc.ID = uuid.New()
	}
	return nil
}

func (tu *TenantUser) BeforeCreate(tx *gorm.DB) error {
	if tu.ID == uuid.Nil {
		tu.ID = uuid.New()
	}
	return nil
}

func (tak *TenantAPIKey) BeforeCreate(tx *gorm.DB) error {
	if tak.ID == uuid.Nil {
		tak.ID = uuid.New()
	}
	if tak.APIKey == "" {
		tak.APIKey = "sk_" + uuid.New().String()
	}
	return nil
}

func (tb *TenantBilling) BeforeCreate(tx *gorm.DB) error {
	if tb.ID == uuid.Nil {
		tb.ID = uuid.New()
	}
	if tb.InvoiceNumber == "" {
		tb.InvoiceNumber = "INV-" + uuid.New().String()[:8]
	}
	return nil
}

func (pi *PartnerIntegration) BeforeCreate(tx *gorm.DB) error {
	if pi.ID == uuid.Nil {
		pi.ID = uuid.New()
	}
	return nil
}

func (pt *PartnerTransaction) BeforeCreate(tx *gorm.DB) error {
	if pt.ID == uuid.Nil {
		pt.ID = uuid.New()
	}
	if pt.TransactionID == "" {
		pt.TransactionID = "PTX-" + uuid.New().String()[:8]
	}
	return nil
}

func (phc *PartnerHealthCheck) BeforeCreate(tx *gorm.DB) error {
	if phc.ID == uuid.Nil {
		phc.ID = uuid.New()
	}
	return nil
}

func (aua *APIUsageAnalytics) BeforeCreate(tx *gorm.DB) error {
	if aua.ID == uuid.Nil {
		aua.ID = uuid.New()
	}
	return nil
}

func (mp *MarketplaceProduct) BeforeCreate(tx *gorm.DB) error {
	if mp.ID == uuid.Nil {
		mp.ID = uuid.New()
	}
	if mp.ProductCode == "" {
		mp.ProductCode = "PRD-" + uuid.New().String()[:8]
	}
	return nil
}

func (ms *MarketplaceSale) BeforeCreate(tx *gorm.DB) error {
	if ms.ID == uuid.Nil {
		ms.ID = uuid.New()
	}
	return nil
}

func (pr *ProductReview) BeforeCreate(tx *gorm.DB) error {
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}
	return nil
}

// Business logic methods for TenantOrganization
func (to *TenantOrganization) CheckActive() bool {
	now := time.Now()
	return to.IsActive && !to.IsSuspended && now.Before(to.ContractEndDate)
}

func (to *TenantOrganization) IsNearCapacity() bool {
	userCapacity := float64(to.CurrentUsers) / float64(to.MaxUsers)
	policyCapacity := float64(to.CurrentPolicies) / float64(to.MaxPolicies)
	return userCapacity > 0.8 || policyCapacity > 0.8
}

func (to *TenantOrganization) Suspend(reason string) {
	to.IsSuspended = true
	to.SuspensionReason = reason
}

func (to *TenantOrganization) Reactivate() {
	to.IsSuspended = false
	to.SuspensionReason = ""
}

func (to *TenantOrganization) CalculateMonthlyRevenue() {
	// This would typically calculate based on policies, claims, and commission rates
	// Placeholder implementation
}

// Business logic methods for TenantAPIKey
func (tak *TenantAPIKey) IsValid() bool {
	if !tak.IsActive {
		return false
	}
	if tak.ExpiresAt != nil && time.Now().After(*tak.ExpiresAt) {
		return false
	}
	return true
}

func (tak *TenantAPIKey) RecordUsage() {
	tak.UsageCount++
	now := time.Now()
	tak.LastUsedAt = &now
}

func (tak *TenantAPIKey) IsRateLimited() bool {
	// This would implement rate limiting logic
	// Placeholder implementation
	return false
}

// Business logic methods for PartnerIntegration
func (pi *PartnerIntegration) IsHealthy() bool {
	return pi.HealthStatus == "healthy" && pi.SuccessRate >= 95.0
}

func (pi *PartnerIntegration) UpdateHealthStatus() {
	if pi.SuccessRate >= 99.0 {
		pi.HealthStatus = "healthy"
	} else if pi.SuccessRate >= 95.0 {
		pi.HealthStatus = "degraded"
	} else {
		pi.HealthStatus = "down"
	}
}

func (pi *PartnerIntegration) CalculateSuccessRate() {
	// This would calculate success rate based on recent transactions
	// Placeholder implementation
}

// Business logic methods for MarketplaceProduct
func (mp *MarketplaceProduct) IsAvailable() bool {
	now := time.Now()
	if mp.EndOfSaleDate != nil && now.After(*mp.EndOfSaleDate) {
		return false
	}
	return mp.IsActive && now.After(mp.LaunchDate)
}

func (mp *MarketplaceProduct) UpdateRating() {
	// This would calculate average rating from reviews
	// Placeholder implementation
}

func (mp *MarketplaceProduct) RecordSale(amount float64) {
	mp.TotalSales++
	mp.MonthlyRevenue += amount
}

// Business logic methods for ProductReview
func (pr *ProductReview) CalculateHelpfulnessRatio() float64 {
	if pr.TotalVotes == 0 {
		return 0
	}
	return float64(pr.HelpfulVotes) / float64(pr.TotalVotes)
}

func (pr *ProductReview) IsHelpful() bool {
	return pr.CalculateHelpfulnessRatio() >= 0.6
}

func (pr *ProductReview) Moderate(moderatorID uuid.UUID, notes string) {
	pr.IsModerated = true
	pr.ModeratedBy = &moderatorID
	pr.ModerationNotes = notes
}
