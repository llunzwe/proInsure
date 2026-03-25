package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PolicyReserves represents financial reserves for a policy
type PolicyReserves struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID        uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	CalculationDate time.Time `gorm:"type:timestamp;not null" json:"calculation_date"`

	// Premium Reserves
	UnearnedPremiumReserve   Money `gorm:"embedded;embeddedPrefix:upr_" json:"unearned_premium_reserve"`
	DeferredAcquisitionCost  Money `gorm:"embedded;embeddedPrefix:dac_" json:"deferred_acquisition_cost"`
	PremiumDeficiencyReserve Money `gorm:"embedded;embeddedPrefix:pdr_" json:"premium_deficiency_reserve"`

	// Loss Reserves
	LossReserve  Money `gorm:"embedded;embeddedPrefix:loss_reserve_" json:"loss_reserve"`
	IBNR         Money `gorm:"embedded;embeddedPrefix:ibnr_" json:"ibnr"`   // Incurred But Not Reported
	IBNER        Money `gorm:"embedded;embeddedPrefix:ibner_" json:"ibner"` // Incurred But Not Enough Reported
	CaseReserves Money `gorm:"embedded;embeddedPrefix:case_reserves_" json:"case_reserves"`
	ALAE         Money `gorm:"embedded;embeddedPrefix:alae_" json:"alae"` // Allocated Loss Adjustment Expense
	ULAE         Money `gorm:"embedded;embeddedPrefix:ulae_" json:"ulae"` // Unallocated Loss Adjustment Expense

	// Technical Provisions
	TechnicalProvisions Money   `gorm:"embedded;embeddedPrefix:technical_provisions_" json:"technical_provisions"`
	BestEstimate        Money   `gorm:"embedded;embeddedPrefix:best_estimate_" json:"best_estimate"`
	RiskMargin          Money   `gorm:"embedded;embeddedPrefix:risk_margin_" json:"risk_margin"`
	DiscountRate        float64 `gorm:"type:decimal(5,4)" json:"discount_rate"`

	// Statutory Reserves
	StatutoryReserves  Money `gorm:"embedded;embeddedPrefix:statutory_reserves_" json:"statutory_reserves"`
	MinimumReserves    Money `gorm:"embedded;embeddedPrefix:minimum_reserves_" json:"minimum_reserves"`
	AdditionalReserves Money `gorm:"embedded;embeddedPrefix:additional_reserves_" json:"additional_reserves"`

	// Reserve Adequacy
	ReserveAdequacyTest bool       `gorm:"type:boolean" json:"reserve_adequacy_test"`
	AdequacyTestDate    *time.Time `gorm:"type:timestamp" json:"adequacy_test_date,omitempty"`
	AdequacyRatio       float64    `gorm:"type:decimal(10,4)" json:"adequacy_ratio"`
	DevelopmentFactor   float64    `gorm:"type:decimal(10,4)" json:"development_factor"`

	// Confidence Levels
	ConfidenceLevel float64 `gorm:"type:decimal(5,2)" json:"confidence_level"`
	VaRAmount       Money   `gorm:"embedded;embeddedPrefix:var_" json:"var_amount"`   // Value at Risk
	TVaRAmount      Money   `gorm:"embedded;embeddedPrefix:tvar_" json:"tvar_amount"` // Tail Value at Risk

	// Reconciliation
	PriorPeriodReserves Money `gorm:"embedded;embeddedPrefix:prior_reserves_" json:"prior_period_reserves"`
	ReserveMovement     Money `gorm:"embedded;embeddedPrefix:reserve_movement_" json:"reserve_movement"`
	ReleasedReserves    Money `gorm:"embedded;embeddedPrefix:released_reserves_" json:"released_reserves"`

	// Actuarial
	ActuarialMethod  string    `gorm:"type:varchar(100)" json:"actuarial_method"`
	ActuaryName      string    `gorm:"type:varchar(255)" json:"actuary_name"`
	ActuarialOpinion string    `gorm:"type:text" json:"actuarial_opinion"`
	LastReviewDate   time.Time `gorm:"type:timestamp" json:"last_review_date"`

	// Metadata
	Notes       string         `gorm:"type:text" json:"notes"`
	Assumptions datatypes.JSON `gorm:"type:json" json:"assumptions"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid" json:"created_by"`
}

// PolicyInvestment represents investment components of a policy
type PolicyInvestment struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Investment Components
	InvestmentComponent Money `gorm:"embedded;embeddedPrefix:investment_component_" json:"investment_component"`
	AccumulatedValue    Money `gorm:"embedded;embeddedPrefix:accumulated_value_" json:"accumulated_value"`
	CashValue           Money `gorm:"embedded;embeddedPrefix:cash_value_" json:"cash_value"`
	AccountValue        Money `gorm:"embedded;embeddedPrefix:account_value_" json:"account_value"`

	// Returns
	GuaranteedReturn float64 `gorm:"type:decimal(5,4)" json:"guaranteed_return_rate"`
	GuaranteedAmount Money   `gorm:"embedded;embeddedPrefix:guaranteed_amount_" json:"guaranteed_amount"`
	ProjectedReturn  float64 `gorm:"type:decimal(5,4)" json:"projected_return_rate"`
	ProjectedAmount  Money   `gorm:"embedded;embeddedPrefix:projected_amount_" json:"projected_amount"`
	ActualReturn     float64 `gorm:"type:decimal(10,4)" json:"actual_return_rate"`
	YearToDateReturn float64 `gorm:"type:decimal(10,4)" json:"ytd_return"`

	// Surrender Values
	SurrenderValue     Money      `gorm:"embedded;embeddedPrefix:surrender_value_" json:"surrender_value"`
	SurrenderCharges   Money      `gorm:"embedded;embeddedPrefix:surrender_charges_" json:"surrender_charges"`
	SurrenderPeriodEnd *time.Time `gorm:"type:timestamp" json:"surrender_period_end,omitempty"`
	FreeLookPeriodEnd  *time.Time `gorm:"type:timestamp" json:"free_look_period_end,omitempty"`

	// Policy Loans
	PolicyLoanAvailable Money   `gorm:"embedded;embeddedPrefix:loan_available_" json:"policy_loan_available"`
	PolicyLoanBalance   Money   `gorm:"embedded;embeddedPrefix:loan_balance_" json:"policy_loan_balance"`
	LoanInterestRate    float64 `gorm:"type:decimal(5,4)" json:"loan_interest_rate"`
	AccruedLoanInterest Money   `gorm:"embedded;embeddedPrefix:accrued_loan_interest_" json:"accrued_loan_interest"`

	// Dividends
	DividendsEarned      Money   `gorm:"embedded;embeddedPrefix:dividends_earned_" json:"dividends_earned"`
	DividendOption       string  `gorm:"type:varchar(50)" json:"dividend_option"` // cash, reinvest, reduce_premium, accumulate
	AccumulatedDividends Money   `gorm:"embedded;embeddedPrefix:accumulated_dividends_" json:"accumulated_dividends"`
	DividendInterestRate float64 `gorm:"type:decimal(5,4)" json:"dividend_interest_rate"`

	// Investment Allocation
	InvestmentStrategy string         `gorm:"type:varchar(100)" json:"investment_strategy"`
	RiskProfile        string         `gorm:"type:varchar(50)" json:"risk_profile"` // conservative, moderate, aggressive
	AssetAllocation    datatypes.JSON `gorm:"type:json" json:"asset_allocation"`    // map[string]float64
	FundOptions        datatypes.JSON `gorm:"type:json" json:"fund_options"`

	// Riders & Benefits
	AcceleratedBenefit Money `gorm:"embedded;embeddedPrefix:accelerated_benefit_" json:"accelerated_benefit"`
	DeathBenefit       Money `gorm:"embedded;embeddedPrefix:death_benefit_" json:"death_benefit"`
	MaturityBenefit    Money `gorm:"embedded;embeddedPrefix:maturity_benefit_" json:"maturity_benefit"`

	// Performance Metrics
	InceptionDate       time.Time `gorm:"type:timestamp" json:"inception_date"`
	TimeWeightedReturn  float64   `gorm:"type:decimal(10,4)" json:"time_weighted_return"`
	MoneyWeightedReturn float64   `gorm:"type:decimal(10,4)" json:"money_weighted_return"`
	BenchmarkReturn     float64   `gorm:"type:decimal(10,4)" json:"benchmark_return"`
	AlphaValue          float64   `gorm:"type:decimal(10,4)" json:"alpha_value"`

	// Fees & Charges
	ManagementFee     float64 `gorm:"type:decimal(5,4)" json:"management_fee_rate"`
	AdministrationFee Money   `gorm:"embedded;embeddedPrefix:admin_fee_" json:"administration_fee"`
	MortalityCharges  Money   `gorm:"embedded;embeddedPrefix:mortality_charges_" json:"mortality_charges"`
	ExpenseCharges    Money   `gorm:"embedded;embeddedPrefix:expense_charges_" json:"expense_charges"`

	// Status
	InvestmentStatus  string    `gorm:"type:varchar(50)" json:"investment_status"`
	LastValuationDate time.Time `gorm:"type:timestamp" json:"last_valuation_date"`
	NextValuationDate time.Time `gorm:"type:timestamp" json:"next_valuation_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// PolicyCommission represents commission structure for a policy
type PolicyCommission struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Commission Structure
	InitialCommission     Money   `gorm:"embedded;embeddedPrefix:initial_commission_" json:"initial_commission"`
	InitialCommissionRate float64 `gorm:"type:decimal(5,2)" json:"initial_commission_rate"`
	RenewalCommission     Money   `gorm:"embedded;embeddedPrefix:renewal_commission_" json:"renewal_commission"`
	RenewalCommissionRate float64 `gorm:"type:decimal(5,2)" json:"renewal_commission_rate"`

	// Override & Bonus
	OverrideCommission   Money   `gorm:"embedded;embeddedPrefix:override_commission_" json:"override_commission"`
	OverrideRate         float64 `gorm:"type:decimal(5,2)" json:"override_rate"`
	BonusCommission      Money   `gorm:"embedded;embeddedPrefix:bonus_commission_" json:"bonus_commission"`
	BonusQualification   string  `gorm:"type:varchar(255)" json:"bonus_qualification"`
	ContingentCommission Money   `gorm:"embedded;embeddedPrefix:contingent_commission_" json:"contingent_commission"`

	// Trail Commission
	TrailCommission     Money   `gorm:"embedded;embeddedPrefix:trail_commission_" json:"trail_commission"`
	TrailCommissionRate float64 `gorm:"type:decimal(5,2)" json:"trail_commission_rate"`
	TrailPeriodMonths   int     `gorm:"type:int" json:"trail_period_months"`

	// Schedule
	CommissionSchedule datatypes.JSON `gorm:"type:json" json:"commission_schedule"` // []CommissionTier
	VestingSchedule    datatypes.JSON `gorm:"type:json" json:"vesting_schedule"`
	PaymentFrequency   string         `gorm:"type:varchar(50)" json:"payment_frequency"` // monthly, quarterly, annual
	NextPaymentDate    *time.Time     `gorm:"type:timestamp" json:"next_payment_date,omitempty"`

	// Clawback
	ClawbackProvisions   datatypes.JSON `gorm:"type:json" json:"clawback_provisions"`
	ClawbackPeriodMonths int            `gorm:"type:int" json:"clawback_period_months"`
	ClawbackAmount       Money          `gorm:"embedded;embeddedPrefix:clawback_amount_" json:"clawback_amount"`
	ClawbackStatus       string         `gorm:"type:varchar(50)" json:"clawback_status"`

	// Payment Details
	TotalPaid         Money      `gorm:"embedded;embeddedPrefix:total_paid_" json:"total_paid"`
	TotalPending      Money      `gorm:"embedded;embeddedPrefix:total_pending_" json:"total_pending"`
	LastPaymentDate   *time.Time `gorm:"type:timestamp" json:"last_payment_date,omitempty"`
	LastPaymentAmount Money      `gorm:"embedded;embeddedPrefix:last_payment_" json:"last_payment_amount"`

	// Payees
	PrimaryAgentID       uuid.UUID      `gorm:"type:uuid" json:"primary_agent_id"`
	PrimaryAgentShare    float64        `gorm:"type:decimal(5,2)" json:"primary_agent_share"`
	SplitCommissions     datatypes.JSON `gorm:"type:json" json:"split_commissions"` // []CommissionSplit
	HierarchyCommissions datatypes.JSON `gorm:"type:json" json:"hierarchy_commissions"`

	// Tax & Compliance
	TaxWithheld      Money   `gorm:"embedded;embeddedPrefix:tax_withheld_" json:"tax_withheld"`
	TaxRate          float64 `gorm:"type:decimal(5,2)" json:"tax_rate"`
	W9OnFile         bool    `gorm:"type:boolean;default:false" json:"w9_on_file"`
	ComplianceStatus string  `gorm:"type:varchar(50)" json:"compliance_status"`

	// Status
	CommissionStatus string     `gorm:"type:varchar(50)" json:"commission_status"`
	PaymentStatus    string     `gorm:"type:varchar(50)" json:"payment_status"`
	ApprovalStatus   string     `gorm:"type:varchar(50)" json:"approval_status"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovalDate     *time.Time `gorm:"type:timestamp" json:"approval_date,omitempty"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// =====================================
// METHODS
// =====================================

// GetTotalReserves calculates total reserves
func (pr *PolicyReserves) GetTotalReserves() float64 {
	return pr.UnearnedPremiumReserve.Amount + pr.LossReserve.Amount +
		pr.IBNR.Amount + pr.ALAE.Amount + pr.ULAE.Amount
}

// IsAdequate checks if reserves are adequate
func (pr *PolicyReserves) IsAdequate() bool {
	return pr.ReserveAdequacyTest && pr.AdequacyRatio >= 1.0
}

// GetRunoffRatio calculates the runoff ratio
func (pr *PolicyReserves) GetRunoffRatio() float64 {
	if pr.PriorPeriodReserves.Amount > 0 {
		return pr.ReleasedReserves.Amount / pr.PriorPeriodReserves.Amount
	}
	return 0
}

// GetNetValue calculates net value after charges
func (pi *PolicyInvestment) GetNetValue() float64 {
	return pi.AccountValue.Amount - pi.PolicyLoanBalance.Amount -
		pi.AccruedLoanInterest.Amount
}

// CanTakeLoan checks if a policy loan can be taken
func (pi *PolicyInvestment) CanTakeLoan() bool {
	return pi.PolicyLoanAvailable.Amount > 0 &&
		pi.InvestmentStatus == "active" &&
		pi.CashValue.Amount > 0
}

// IsInSurrenderPeriod checks if policy is still in surrender period
func (pi *PolicyInvestment) IsInSurrenderPeriod() bool {
	if pi.SurrenderPeriodEnd == nil {
		return false
	}
	return time.Now().Before(*pi.SurrenderPeriodEnd)
}

// GetTotalCommissionDue calculates total commission due
func (pc *PolicyCommission) GetTotalCommissionDue() float64 {
	return pc.InitialCommission.Amount + pc.RenewalCommission.Amount +
		pc.OverrideCommission.Amount + pc.BonusCommission.Amount +
		pc.TrailCommission.Amount - pc.ClawbackAmount.Amount
}

// IsPaymentDue checks if commission payment is due
func (pc *PolicyCommission) IsPaymentDue() bool {
	if pc.NextPaymentDate == nil {
		return false
	}
	return time.Now().After(*pc.NextPaymentDate) && pc.TotalPending.Amount > 0
}

// IsSubjectToClawback checks if commission is subject to clawback
func (pc *PolicyCommission) IsSubjectToClawback(monthsSinceIssue int) bool {
	return monthsSinceIssue <= pc.ClawbackPeriodMonths &&
		pc.ClawbackProvisions != nil
}
