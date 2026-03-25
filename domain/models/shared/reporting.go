package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/policy"
	"smartsure/internal/domain/types"
)

// Report represents generated reports
type Report struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ReportNumber string    `gorm:"uniqueIndex;not null" json:"report_number"`
	ReportType   string    `gorm:"not null" json:"report_type"` // financial, operational, claims, risk
	ReportName   string    `gorm:"not null" json:"report_name"`
	Description  string    `json:"description"`

	// Report Configuration
	Template  string `json:"template"`
	Format    string `json:"format"`    // pdf, excel, csv, json
	Frequency string `json:"frequency"` // daily, weekly, monthly, quarterly, annual, adhoc

	// Time Period
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`

	// Generation Details
	GeneratedBy    uuid.UUID `gorm:"type:uuid" json:"generated_by"`
	GeneratedAt    time.Time `json:"generated_at"`
	GenerationTime int       `json:"generation_time"` // seconds

	// File Details
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
	FileHash string `json:"file_hash"`

	// Distribution
	Recipients     types.JSONArray `gorm:"type:json" json:"recipients"`
	DeliveryMethod string           `json:"delivery_method"` // email, download, api
	DeliveredAt    *time.Time       `json:"delivered_at"`

	// Parameters
	Parameters string `gorm:"type:json" json:"parameters"`
	Filters    string `gorm:"type:json" json:"filters"`

	// Status
	Status       string `json:"status"` // pending, generating, completed, failed
	ErrorMessage string `json:"error_message"`
	RetryCount   int    `json:"retry_count"`

	// Metadata
	Tags        types.JSONArray `gorm:"type:json" json:"tags"`
	IsScheduled bool             `json:"is_scheduled"`
	ScheduleID  *uuid.UUID       `gorm:"type:uuid" json:"schedule_id"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - temporarily commented for migration
	// Generator         *User          `gorm:"foreignKey:GeneratedBy" json:"generator,omitempty"`
	// Schedule          *ReportSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
}

// ReportSchedule represents scheduled reports
type ReportSchedule struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ScheduleName string    `gorm:"not null" json:"schedule_name"`
	ReportType   string    `json:"report_type"`

	// Schedule Configuration
	Frequency      string `json:"frequency"` // hourly, daily, weekly, monthly
	CronExpression string `json:"cron_expression"`
	TimeZone       string `json:"timezone"`

	// Report Configuration
	Template   string `json:"template"`
	Parameters string `gorm:"type:json" json:"parameters"`
	Format     string `json:"format"`

	// Distribution
	Recipients     types.JSONArray `gorm:"type:json" json:"recipients"`
	DeliveryMethod string           `json:"delivery_method"`

	// Execution
	NextRunAt     *time.Time `json:"next_run_at"`
	LastRunAt     *time.Time `json:"last_run_at"`
	LastRunStatus string     `json:"last_run_status"`
	RunCount      int        `json:"run_count"`

	// Status
	IsActive bool `json:"is_active"`
	IsPaused bool `json:"is_paused"`

	// Ownership
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - temporarily commented for migration
	// Creator           *User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	// Reports           []Report       `gorm:"foreignKey:ScheduleID" json:"reports,omitempty"`
}

// Dashboard represents analytical dashboards
type Dashboard struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DashboardCode string    `gorm:"uniqueIndex;not null" json:"dashboard_code"`
	DashboardName string    `gorm:"not null" json:"dashboard_name"`
	Description   string    `json:"description"`

	// Dashboard Type
	Type     string `json:"type"`     // executive, operational, analytical
	Category string `json:"category"` // sales, claims, finance, risk

	// Layout Configuration
	Layout  string           `gorm:"type:json" json:"layout"`
	Widgets types.JSONArray `gorm:"type:json" json:"widgets"`
	Theme   string           `json:"theme"`

	// Access Control
	IsPublic     bool             `json:"is_public"`
	AccessLevel  string           `json:"access_level"`
	AllowedRoles types.JSONArray `gorm:"type:json" json:"allowed_roles"`

	// Ownership
	OwnerID    uuid.UUID        `gorm:"type:uuid" json:"owner_id"`
	IsShared   bool             `json:"is_shared"`
	SharedWith types.JSONArray `gorm:"type:json" json:"shared_with"`

	// Refresh Settings
	AutoRefresh     bool       `json:"auto_refresh"`
	RefreshInterval int        `json:"refresh_interval"` // seconds
	LastRefreshedAt *time.Time `json:"last_refreshed_at"`

	// Usage
	ViewCount    int        `json:"view_count"`
	LastViewedAt *time.Time `json:"last_viewed_at"`
	IsFavorite   bool       `json:"is_favorite"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Owner *models.User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

// KPI represents key performance indicators
type KPI struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	KPICode     string    `gorm:"uniqueIndex;not null" json:"kpi_code"`
	KPIName     string    `gorm:"not null" json:"kpi_name"`
	Description string    `json:"description"`

	// KPI Configuration
	Category string `json:"category"`
	Type     string `json:"type"` // metric, ratio, percentage, count
	Unit     string `json:"unit"`

	// Calculation
	Formula         string `json:"formula"`
	DataSource      string `json:"data_source"`
	AggregationType string `json:"aggregation_type"` // sum, avg, count, min, max

	// Current Values
	CurrentValue  float64 `json:"current_value"`
	PreviousValue float64 `json:"previous_value"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Trend         string  `json:"trend"` // up, down, stable

	// Targets
	TargetValue       float64 `json:"target_value"`
	MinValue          float64 `json:"min_value"`
	MaxValue          float64 `json:"max_value"`
	ThresholdWarning  float64 `json:"threshold_warning"`
	ThresholdCritical float64 `json:"threshold_critical"`

	// Period
	Period      string    `json:"period"` // daily, weekly, monthly, quarterly
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`

	// Status
	Status   string `json:"status"` // on_target, warning, critical
	IsActive bool   `json:"is_active"`

	// Update Information
	LastCalculatedAt  time.Time  `json:"last_calculated_at"`
	NextCalculationAt *time.Time `json:"next_calculation_at"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// FinancialTransaction represents financial transactions
type FinancialTransaction struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TransactionNumber string    `gorm:"uniqueIndex;not null" json:"transaction_number"`

	// Transaction Details
	Type        string `json:"type"` // premium, claim_payment, refund, commission
	Category    string `json:"category"`
	Description string `json:"description"`

	// Financial Details
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate"`
	BaseAmount   float64 `json:"base_amount"` // in base currency

	// Accounting
	DebitAccount  string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	CostCenter    string `json:"cost_center"`
	GLCode        string `json:"gl_code"`

	// Related Entities
	EntityType string     `json:"entity_type"` // policy, claim, payment
	EntityID   string     `json:"entity_id"`
	UserID     *uuid.UUID `gorm:"type:uuid" json:"user_id"`

	// Processing
	Status      string     `json:"status"` // pending, processing, completed, failed, reversed
	ProcessedAt *time.Time `json:"processed_at"`
	ProcessedBy *uuid.UUID `gorm:"type:uuid" json:"processed_by"`

	// Settlement
	SettlementDate *time.Time `json:"settlement_date"`
	SettlementRef  string     `json:"settlement_ref"`
	BankReference  string     `json:"bank_reference"`

	// Reconciliation
	IsReconciled      bool       `json:"is_reconciled"`
	ReconciledAt      *time.Time `json:"reconciled_at"`
	ReconciledBy      *uuid.UUID `gorm:"type:uuid" json:"reconciled_by"`
	ReconciliationRef string     `json:"reconciliation_ref"`

	// Audit
	ReversalReason string     `json:"reversal_reason"`
	ReversedAt     *time.Time `json:"reversed_at"`
	ReversedBy     *uuid.UUID `gorm:"type:uuid" json:"reversed_by"`

	// Metadata
	Notes string           `json:"notes"`
	Tags  types.JSONArray `gorm:"type:json" json:"tags"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User       *models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Processor  *models.User `gorm:"foreignKey:ProcessedBy" json:"processor,omitempty"`
	Reconciler *models.User `gorm:"foreignKey:ReconciledBy" json:"reconciler,omitempty"`
	Reverser   *models.User `gorm:"foreignKey:ReversedBy" json:"reverser,omitempty"`
}

// ActuarialData represents actuarial calculations and data
type ActuarialData struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DataType string    `json:"data_type"` // loss_ratio, frequency, severity, exposure
	Category string    `json:"category"`

	// Period
	PeriodType  string    `json:"period_type"` // monthly, quarterly, annual
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`

	// Metrics
	PolicyCount    int     `json:"policy_count"`
	PremiumEarned  float64 `json:"premium_earned"`
	PremiumWritten float64 `json:"premium_written"`
	ClaimCount     int     `json:"claim_count"`
	ClaimAmount    float64 `json:"claim_amount"`
	ClaimPaid      float64 `json:"claim_paid"`
	ClaimReserve   float64 `json:"claim_reserve"`

	// Ratios & Calculations
	LossRatio      float64 `json:"loss_ratio"`
	ExpenseRatio   float64 `json:"expense_ratio"`
	CombinedRatio  float64 `json:"combined_ratio"`
	ClaimFrequency float64 `json:"claim_frequency"`
	ClaimSeverity  float64 `json:"claim_severity"`

	// Risk Metrics
	PureRiskPremium  float64 `json:"pure_risk_premium"`
	TechnicalPremium float64 `json:"technical_premium"`
	RiskAdjustment   float64 `json:"risk_adjustment"`
	SafetyMargin     float64 `json:"safety_margin"`

	// Reserves
	UnearnedPremium   float64 `json:"unearned_premium"`
	OutstandingClaims float64 `json:"outstanding_claims"`
	IBNR              float64 `json:"ibnr"` // Incurred But Not Reported
	TechnicalReserves float64 `json:"technical_reserves"`

	// Product/Segment
	ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Segment   string     `json:"segment"`
	Region    string     `json:"region"`

	// Confidence & Quality
	ConfidenceLevel float64 `json:"confidence_level"`
	DataQuality     string  `json:"data_quality"`
	SampleSize      int     `json:"sample_size"`

	// Processing
	CalculatedAt time.Time `json:"calculated_at"`
	CalculatedBy string    `json:"calculated_by"` // system, manual

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Product *policy.Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
