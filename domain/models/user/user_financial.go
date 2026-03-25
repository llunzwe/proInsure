package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserFinancialProfile contains detailed financial assessment and history
type UserFinancialProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Income Information
	AnnualIncome           decimal.Decimal `gorm:"type:decimal(15,2)" json:"annual_income"`
	MonthlyIncome          decimal.Decimal `gorm:"type:decimal(15,2)" json:"monthly_income"`
	IncomeSource           string          `gorm:"type:varchar(50)" json:"income_source"`
	IncomeVerified         bool            `gorm:"default:false" json:"income_verified"`
	IncomeVerificationDate *time.Time      `json:"income_verification_date"`
	EmploymentStatus       string          `gorm:"type:varchar(50)" json:"employment_status"`
	Employer               string          `gorm:"type:varchar(100)" json:"employer"`
	JobTitle               string          `gorm:"type:varchar(100)" json:"job_title"`
	YearsEmployed          float64         `gorm:"default:0" json:"years_employed"`

	// Assets & Liabilities
	TotalAssets        decimal.Decimal            `gorm:"type:decimal(15,2)" json:"total_assets"`
	LiquidAssets       decimal.Decimal            `gorm:"type:decimal(15,2)" json:"liquid_assets"`
	TotalLiabilities   decimal.Decimal            `gorm:"type:decimal(15,2)" json:"total_liabilities"`
	NetWorth           decimal.Decimal            `gorm:"type:decimal(15,2)" json:"net_worth"`
	DebtToIncomeRatio  float64                    `gorm:"default:0" json:"debt_to_income_ratio"`
	AssetBreakdown     map[string]decimal.Decimal `gorm:"type:json" json:"asset_breakdown"`
	LiabilityBreakdown map[string]decimal.Decimal `gorm:"type:json" json:"liability_breakdown"`

	// Banking Information
	PrimaryBank              string          `gorm:"type:varchar(100)" json:"primary_bank"`
	BankAccountTypes         []string        `gorm:"type:json" json:"bank_account_types"`
	BankingRelationshipYears float64         `gorm:"default:0" json:"banking_relationship_years"`
	AverageBalance           decimal.Decimal `gorm:"type:decimal(15,2)" json:"average_balance"`
	OverdraftHistory         int             `gorm:"default:0" json:"overdraft_history"`

	// Investment Profile
	InvestmentExperience string                     `gorm:"type:varchar(50)" json:"investment_experience"`
	RiskTolerance        string                     `gorm:"type:varchar(20)" json:"risk_tolerance"`
	InvestmentPortfolio  map[string]decimal.Decimal `gorm:"type:json" json:"investment_portfolio"`
	PortfolioValue       decimal.Decimal            `gorm:"type:decimal(15,2)" json:"portfolio_value"`
	InvestmentGoals      []string                   `gorm:"type:json" json:"investment_goals"`

	// Financial Health Metrics
	FinancialHealthScore float64 `gorm:"default:0" json:"financial_health_score"`
	CashFlowScore        float64 `gorm:"default:0" json:"cash_flow_score"`
	SavingsRate          float64 `gorm:"default:0" json:"savings_rate"`
	EmergencyFundMonths  float64 `gorm:"default:0" json:"emergency_fund_months"`
	FinancialStressLevel string  `gorm:"type:varchar(20)" json:"financial_stress_level"`

	// Financial Goals
	ShortTermGoals      []map[string]interface{} `gorm:"type:json" json:"short_term_goals"`
	LongTermGoals       []map[string]interface{} `gorm:"type:json" json:"long_term_goals"`
	RetirementPlanning  map[string]interface{}   `gorm:"type:json" json:"retirement_planning"`
	FinancialMilestones []map[string]interface{} `gorm:"type:json" json:"financial_milestones"`
}

// UserCreditProfile manages credit scoring and history
type UserCreditProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Credit Scores
	CreditScore         int                      `gorm:"default:0" json:"credit_score"`
	CreditScoreProvider string                   `gorm:"type:varchar(50)" json:"credit_score_provider"`
	ScoreDate           time.Time                `json:"score_date"`
	ScoreRange          string                   `gorm:"type:varchar(20)" json:"score_range"` // poor/fair/good/excellent
	PreviousScores      []map[string]interface{} `gorm:"type:json" json:"previous_scores"`
	ScoreTrend          string                   `gorm:"type:varchar(20)" json:"score_trend"`
	ScoreFactors        []string                 `gorm:"type:json" json:"score_factors"`

	// Credit History
	CreditHistoryLength int        `json:"credit_history_length_months"`
	OldestAccount       *time.Time `json:"oldest_account"`
	NewestAccount       *time.Time `json:"newest_account"`
	TotalAccounts       int        `gorm:"default:0" json:"total_accounts"`
	OpenAccounts        int        `gorm:"default:0" json:"open_accounts"`
	ClosedAccounts      int        `gorm:"default:0" json:"closed_accounts"`

	// Credit Utilization
	TotalCreditLimit decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_credit_limit"`
	TotalCreditUsed  decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_credit_used"`
	UtilizationRatio float64         `gorm:"default:0" json:"utilization_ratio"`

	// Payment History
	OnTimePayments      int        `gorm:"default:0" json:"on_time_payments"`
	LatePayments        int        `gorm:"default:0" json:"late_payments"`
	MissedPayments      int        `gorm:"default:0" json:"missed_payments"`
	PaymentHistoryScore float64    `gorm:"default:0" json:"payment_history_score"`
	LastLatePayment     *time.Time `json:"last_late_payment"`
	ConsecutiveOnTime   int        `gorm:"default:0" json:"consecutive_on_time"`

	// Negative Marks
	Defaults        int        `gorm:"default:0" json:"defaults"`
	Bankruptcies    int        `gorm:"default:0" json:"bankruptcies"`
	Collections     int        `gorm:"default:0" json:"collections"`
	PublicRecords   int        `gorm:"default:0" json:"public_records"`
	HardInquiries   int        `gorm:"default:0" json:"hard_inquiries"`
	LastHardInquiry *time.Time `json:"last_hard_inquiry"`

	// Credit Monitoring
	MonitoringEnabled  bool                     `gorm:"default:false" json:"monitoring_enabled"`
	AlertsEnabled      bool                     `gorm:"default:false" json:"alerts_enabled"`
	LastMonitoringDate *time.Time               `json:"last_monitoring_date"`
	ChangeAlerts       []map[string]interface{} `gorm:"type:json" json:"change_alerts"`
}

// UserPaymentBehavior tracks payment patterns and behavior
type UserPaymentBehavior struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Payment Patterns
	PreferredPaymentMethod string   `gorm:"type:varchar(50)" json:"preferred_payment_method"`
	PaymentMethods         []string `gorm:"type:json" json:"payment_methods"`
	PreferredPaymentDay    int      `json:"preferred_payment_day"`
	PreferredPaymentTime   string   `gorm:"type:varchar(10)" json:"preferred_payment_time"`
	AutoPayEnabled         bool     `gorm:"default:false" json:"auto_pay_enabled"`
	AutoPayMethods         []string `gorm:"type:json" json:"auto_pay_methods"`

	// Payment Metrics
	TotalPayments        int             `gorm:"default:0" json:"total_payments"`
	SuccessfulPayments   int             `gorm:"default:0" json:"successful_payments"`
	FailedPayments       int             `gorm:"default:0" json:"failed_payments"`
	AveragePaymentAmount decimal.Decimal `gorm:"type:decimal(15,2)" json:"average_payment_amount"`
	TotalPaymentAmount   decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_payment_amount"`
	PaymentSuccessRate   float64         `gorm:"default:0" json:"payment_success_rate"`

	// Payment Timing
	EarlyPayments    int     `gorm:"default:0" json:"early_payments"`
	OnTimePayments   int     `gorm:"default:0" json:"on_time_payments"`
	LatePayments     int     `gorm:"default:0" json:"late_payments"`
	AverageDaysEarly float64 `gorm:"default:0" json:"average_days_early"`
	AverageDaysLate  float64 `gorm:"default:0" json:"average_days_late"`
	PaymentVelocity  float64 `gorm:"default:0" json:"payment_velocity"`

	// Retry Behavior
	RetryAttempts     int     `gorm:"default:0" json:"retry_attempts"`
	SuccessfulRetries int     `gorm:"default:0" json:"successful_retries"`
	RetrySuccessRate  float64 `gorm:"default:0" json:"retry_success_rate"`
	AverageRetryDelay int     `json:"average_retry_delay_hours"`

	// Failure Analysis
	FailureReasons map[string]int `gorm:"type:json" json:"failure_reasons"`
	RecoveryRate   float64        `gorm:"default:0" json:"recovery_rate"`
	ChurnRiskScore float64        `gorm:"default:0" json:"churn_risk_score"`
}

// UserBillingProfile manages billing information and history
type UserBillingProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Billing Details
	BillingCycle    string     `gorm:"type:varchar(20)" json:"billing_cycle"`
	BillingDay      int        `json:"billing_day"`
	NextBillingDate *time.Time `json:"next_billing_date"`
	LastBillingDate *time.Time `json:"last_billing_date"`
	BillingMethod   string     `gorm:"type:varchar(50)" json:"billing_method"`
	Currency        string     `gorm:"type:varchar(3)" json:"currency"`

	// Billing Addresses
	PrimaryBillingAddress     map[string]string   `gorm:"type:json" json:"primary_billing_address"`
	AlternateBillingAddresses []map[string]string `gorm:"type:json" json:"alternate_billing_addresses"`
	TaxJurisdiction           string              `gorm:"type:varchar(100)" json:"tax_jurisdiction"`
	TaxExempt                 bool                `gorm:"default:false" json:"tax_exempt"`
	TaxExemptionNumber        string              `gorm:"type:varchar(50)" json:"tax_exemption_number"`

	// Outstanding Balances
	CurrentBalance    decimal.Decimal `gorm:"type:decimal(15,2)" json:"current_balance"`
	PastDueAmount     decimal.Decimal `gorm:"type:decimal(15,2)" json:"past_due_amount"`
	TotalOutstanding  decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_outstanding"`
	LastPaymentAmount decimal.Decimal `gorm:"type:decimal(15,2)" json:"last_payment_amount"`
	LastPaymentDate   *time.Time      `json:"last_payment_date"`

	// Billing History
	TotalBilled   decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_billed"`
	TotalPaid     decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_paid"`
	TotalRefunded decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_refunded"`
	TotalCredits  decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_credits"`
	TotalDebits   decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_debits"`

	// Collections
	InCollections       bool            `gorm:"default:false" json:"in_collections"`
	CollectionAgency    string          `gorm:"type:varchar(100)" json:"collection_agency"`
	CollectionStartDate *time.Time      `json:"collection_start_date"`
	CollectionAmount    decimal.Decimal `gorm:"type:decimal(15,2)" json:"collection_amount"`
	CollectionStatus    string          `gorm:"type:varchar(50)" json:"collection_status"`
}

// UserTaxProfile manages tax-related information
type UserTaxProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Tax Identification
	TaxID             string `gorm:"type:varchar(50)" json:"tax_id"`
	TaxIDType         string `gorm:"type:varchar(20)" json:"tax_id_type"`
	TaxCountry        string `gorm:"type:varchar(2)" json:"tax_country"`
	TaxResidency      string `gorm:"type:varchar(2)" json:"tax_residency"`
	TaxClassification string `gorm:"type:varchar(50)" json:"tax_classification"`

	// Tax Documents
	W9Provided          bool       `gorm:"default:false" json:"w9_provided"`
	W8Provided          bool       `gorm:"default:false" json:"w8_provided"`
	TaxFormsOnFile      []string   `gorm:"type:json" json:"tax_forms_on_file"`
	LastTaxDocumentDate *time.Time `json:"last_tax_document_date"`

	// Tax Reporting
	ReportingYear         int                      `json:"reporting_year"`
	TotalTaxablePremiums  decimal.Decimal          `gorm:"type:decimal(15,2)" json:"total_taxable_premiums"`
	TotalDeductions       decimal.Decimal          `gorm:"type:decimal(15,2)" json:"total_deductions"`
	TaxWithheld           decimal.Decimal          `gorm:"type:decimal(15,2)" json:"tax_withheld"`
	EstimatedTaxLiability decimal.Decimal          `gorm:"type:decimal(15,2)" json:"estimated_tax_liability"`
	TaxDocumentsSent      []map[string]interface{} `gorm:"type:json" json:"tax_documents_sent"`

	// International Tax
	FATCAReporting   bool               `gorm:"default:false" json:"fatca_reporting"`
	CRSReporting     bool               `gorm:"default:false" json:"crs_reporting"`
	TaxTreaties      []string           `gorm:"type:json" json:"tax_treaties"`
	WithholdingRates map[string]float64 `gorm:"type:json" json:"withholding_rates"`
}

// UserInvestmentProfile tracks investment products and returns
type UserInvestmentProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Investment Products
	InvestmentPolicies   []uuid.UUID                `gorm:"type:json" json:"investment_policies"`
	PolicyCashValues     map[string]decimal.Decimal `gorm:"type:json" json:"policy_cash_values"`
	PolicyLoans          map[string]decimal.Decimal `gorm:"type:json" json:"policy_loans"`
	DividendOptions      map[string]string          `gorm:"type:json" json:"dividend_options"`
	AccumulatedDividends decimal.Decimal            `gorm:"type:decimal(15,2)" json:"accumulated_dividends"`

	// Performance
	TotalInvested    decimal.Decimal `gorm:"type:decimal(15,2)" json:"total_invested"`
	CurrentValue     decimal.Decimal `gorm:"type:decimal(15,2)" json:"current_value"`
	UnrealizedGains  decimal.Decimal `gorm:"type:decimal(15,2)" json:"unrealized_gains"`
	RealizedGains    decimal.Decimal `gorm:"type:decimal(15,2)" json:"realized_gains"`
	YearToDateReturn float64         `gorm:"default:0" json:"year_to_date_return"`
	LifetimeReturn   float64         `gorm:"default:0" json:"lifetime_return"`
	AnnualizedReturn float64         `gorm:"default:0" json:"annualized_return"`

	// Risk Profile
	RiskAssessmentDate *time.Time `json:"risk_assessment_date"`
	RiskScore          float64    `gorm:"default:0" json:"risk_score"`
	RiskCapacity       string     `gorm:"type:varchar(20)" json:"risk_capacity"`
	TimeHorizon        string     `gorm:"type:varchar(20)" json:"time_horizon"`
	LiquidityNeeds     string     `gorm:"type:varchar(20)" json:"liquidity_needs"`

	// Allocations
	AssetAllocation   map[string]float64 `gorm:"type:json" json:"asset_allocation"`
	TargetAllocation  map[string]float64 `gorm:"type:json" json:"target_allocation"`
	RebalancingNeeded bool               `gorm:"default:false" json:"rebalancing_needed"`
	LastRebalanceDate *time.Time         `json:"last_rebalance_date"`
}

// TableName returns the table name
func (UserFinancialProfile) TableName() string {
	return "user_financial_profiles"
}

// TableName returns the table name
func (UserCreditProfile) TableName() string {
	return "user_credit_profiles"
}

// TableName returns the table name
func (UserPaymentBehavior) TableName() string {
	return "user_payment_behavior"
}

// TableName returns the table name
func (UserBillingProfile) TableName() string {
	return "user_billing_profiles"
}

// TableName returns the table name
func (UserTaxProfile) TableName() string {
	return "user_tax_profiles"
}

// TableName returns the table name
func (UserInvestmentProfile) TableName() string {
	return "user_investment_profiles"
}
