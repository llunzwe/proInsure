package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ReinsuranceType constants
const (
	ReinsuranceTypeTreaty       = "treaty"
	ReinsuranceTypeFacultative  = "facultative"
	ReinsuranceTypeExcessOfLoss = "excess_of_loss"
	ReinsuranceTypeQuotaShare   = "quota_share"
	ReinsuranceTypeSurplus      = "surplus"
	ReinsuranceTypeStopLoss     = "stop_loss"
	ReinsuranceTypeCatastrophe  = "catastrophe"
)

// PolicyReinsurance represents reinsurance arrangements for a policy
type PolicyReinsurance struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID        uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	ReinsuranceType string    `gorm:"type:varchar(50);not null" json:"reinsurance_type"`
	TreatyID        string    `gorm:"type:varchar(100)" json:"treaty_id"`
	TreatyName      string    `gorm:"type:varchar(255)" json:"treaty_name"`

	// Financial Details
	CededPremium          Money   `gorm:"embedded;embeddedPrefix:ceded_premium_" json:"ceded_premium"`
	CededRisk             float64 `gorm:"type:decimal(5,2)" json:"ceded_risk_percentage"`
	RetentionLimit        Money   `gorm:"embedded;embeddedPrefix:retention_limit_" json:"retention_limit"`
	CessionLimit          Money   `gorm:"embedded;embeddedPrefix:cession_limit_" json:"cession_limit"`
	ReinsuranceCommission float64 `gorm:"type:decimal(5,2)" json:"reinsurance_commission_rate"`
	ProfitCommission      float64 `gorm:"type:decimal(5,2)" json:"profit_commission_rate"`

	// Coverage Details
	EffectiveDate      time.Time `gorm:"type:timestamp" json:"effective_date"`
	ExpirationDate     time.Time `gorm:"type:timestamp" json:"expiration_date"`
	LayerStart         Money     `gorm:"embedded;embeddedPrefix:layer_start_" json:"layer_start"`
	LayerEnd           Money     `gorm:"embedded;embeddedPrefix:layer_end_" json:"layer_end"`
	ReinstatementLimit int       `gorm:"type:int" json:"reinstatement_limit"`
	ReinstatementCost  Money     `gorm:"embedded;embeddedPrefix:reinstatement_cost_" json:"reinstatement_cost"`

	// Recovery & Claims
	RecoverableAmount Money      `gorm:"embedded;embeddedPrefix:recoverable_" json:"recoverable_amount"`
	RecoveredAmount   Money      `gorm:"embedded;embeddedPrefix:recovered_" json:"recovered_amount"`
	PendingRecovery   Money      `gorm:"embedded;embeddedPrefix:pending_recovery_" json:"pending_recovery"`
	ClaimsReported    int        `gorm:"type:int" json:"claims_reported"`
	LastClaimDate     *time.Time `gorm:"type:timestamp" json:"last_claim_date,omitempty"`

	// Partners
	LeadReinsurer    string         `gorm:"type:varchar(255)" json:"lead_reinsurer"`
	ReinsurerShares  datatypes.JSON `gorm:"type:json" json:"reinsurer_shares"` // map[string]float64
	BrokerName       string         `gorm:"type:varchar(255)" json:"broker_name"`
	BrokerCommission float64        `gorm:"type:decimal(5,2)" json:"broker_commission_rate"`

	// Terms & Conditions
	TreatyTerms       datatypes.JSON `gorm:"type:json" json:"treaty_terms"`
	ExclusionsList    datatypes.JSON `gorm:"type:json" json:"exclusions_list"`
	TerritorialScope  datatypes.JSON `gorm:"type:json" json:"territorial_scope"`
	ArbitrationClause string         `gorm:"type:text" json:"arbitration_clause"`

	// Risk Assessment
	CatastropheExposure Money   `gorm:"embedded;embeddedPrefix:cat_exposure_" json:"catastrophe_exposure"`
	AggregateLimit      Money   `gorm:"embedded;embeddedPrefix:aggregate_limit_" json:"aggregate_limit"`
	AggregateDeductible Money   `gorm:"embedded;embeddedPrefix:aggregate_deductible_" json:"aggregate_deductible"`
	LossRatio           float64 `gorm:"type:decimal(10,2)" json:"loss_ratio"`

	// Status
	Status             string     `gorm:"type:varchar(50)" json:"status"`
	RenewalStatus      string     `gorm:"type:varchar(50)" json:"renewal_status"`
	CancellationDate   *time.Time `gorm:"type:timestamp" json:"cancellation_date,omitempty"`
	CancellationReason string     `gorm:"type:text" json:"cancellation_reason"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// PolicyCoinsurance represents co-insurance arrangements
type PolicyCoinsurance struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID    uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	AgreementID string    `gorm:"type:varchar(100)" json:"agreement_id"`

	// Leadership
	IsLeadInsurer   bool      `gorm:"type:boolean;default:false" json:"is_lead_insurer"`
	LeadInsurerID   uuid.UUID `gorm:"type:uuid" json:"lead_insurer_id"`
	LeadInsurerName string    `gorm:"type:varchar(255)" json:"lead_insurer_name"`
	LeadShare       float64   `gorm:"type:decimal(5,2)" json:"lead_share_percentage"`

	// Coinsurers
	CoinsurersCount   int            `gorm:"type:int" json:"coinsurers_count"`
	ShareDistribution datatypes.JSON `gorm:"type:json" json:"share_distribution"` // map[string]CoinsuranceShare
	TotalShares       float64        `gorm:"type:decimal(5,2)" json:"total_shares"`

	// Financial
	GrossPremium       Money   `gorm:"embedded;embeddedPrefix:gross_premium_" json:"gross_premium"`
	OurShare           float64 `gorm:"type:decimal(5,2)" json:"our_share_percentage"`
	OurPremium         Money   `gorm:"embedded;embeddedPrefix:our_premium_" json:"our_premium"`
	OurClaimsLiability Money   `gorm:"embedded;embeddedPrefix:our_claims_liability_" json:"our_claims_liability"`

	// Claims Handling
	ClaimsLeader        string `gorm:"type:varchar(255)" json:"claims_leader"`
	ClaimsCooperation   string `gorm:"type:varchar(50)" json:"claims_cooperation"` // follow_the_leader, unanimous, proportional
	SettlementAuthority Money  `gorm:"embedded;embeddedPrefix:settlement_authority_" json:"settlement_authority"`

	// Terms
	EffectiveDate    time.Time      `gorm:"type:timestamp" json:"effective_date"`
	ExpirationDate   time.Time      `gorm:"type:timestamp" json:"expiration_date"`
	CoinsuranceTerms datatypes.JSON `gorm:"type:json" json:"coinsurance_terms"`

	// Administrative
	AdministrationFee float64 `gorm:"type:decimal(5,2)" json:"administration_fee_rate"`
	LeadFee           float64 `gorm:"type:decimal(5,2)" json:"lead_fee_rate"`
	BillingMethod     string  `gorm:"type:varchar(50)" json:"billing_method"` // direct, lead_administered

	// Status
	Status         string     `gorm:"type:varchar(50)" json:"status"`
	ApprovalStatus string     `gorm:"type:varchar(50)" json:"approval_status"`
	ApprovedBy     *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovalDate   *time.Time `gorm:"type:timestamp" json:"approval_date,omitempty"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// =====================================
// METHODS
// =====================================

// IsActive checks if reinsurance is currently active
func (pr *PolicyReinsurance) IsActive() bool {
	now := time.Now()
	return pr.Status == "active" &&
		now.After(pr.EffectiveDate) &&
		now.Before(pr.ExpirationDate)
}

// GetNetRetention calculates the net retention after reinsurance
func (pr *PolicyReinsurance) GetNetRetention(grossAmount float64) float64 {
	if pr.CededRisk > 0 {
		return grossAmount * (1 - pr.CededRisk/100)
	}
	if pr.RetentionLimit.Amount > 0 && grossAmount > pr.RetentionLimit.Amount {
		return pr.RetentionLimit.Amount
	}
	return grossAmount
}

// GetRecoveryRate calculates the recovery rate from reinsurance
func (pr *PolicyReinsurance) GetRecoveryRate() float64 {
	if pr.RecoverableAmount.Amount > 0 {
		return (pr.RecoveredAmount.Amount / pr.RecoverableAmount.Amount) * 100
	}
	return 0
}

// CanFileRecovery checks if a recovery claim can be filed
func (pr *PolicyReinsurance) CanFileRecovery(claimAmount float64) bool {
	if !pr.IsActive() {
		return false
	}
	if pr.ReinsuranceType == ReinsuranceTypeExcessOfLoss {
		return claimAmount > pr.LayerStart.Amount
	}
	return pr.CededRisk > 0
}

// GetOurSharePercentage returns our share percentage in coinsurance
func (pc *PolicyCoinsurance) GetOurSharePercentage() float64 {
	if pc.IsLeadInsurer {
		return pc.LeadShare
	}
	return pc.OurShare
}

// IsCoinsuranceActive checks if coinsurance is currently active
func (pc *PolicyCoinsurance) IsCoinsuranceActive() bool {
	now := time.Now()
	return pc.Status == "active" &&
		now.After(pc.EffectiveDate) &&
		now.Before(pc.ExpirationDate)
}

// CalculateOurLiability calculates our liability for a claim amount
func (pc *PolicyCoinsurance) CalculateOurLiability(claimAmount float64) float64 {
	sharePercentage := pc.GetOurSharePercentage()
	return claimAmount * (sharePercentage / 100)
}

// RequiresConsensus checks if consensus is required for decisions
func (pc *PolicyCoinsurance) RequiresConsensus() bool {
	return pc.ClaimsCooperation == "unanimous"
}
