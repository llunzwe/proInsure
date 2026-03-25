package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// ReinsuranceContract represents reinsurance agreements
type ReinsuranceContract struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ContractNumber    string         `gorm:"uniqueIndex;not null" json:"contract_number"`
	ReinsurerName     string         `gorm:"not null" json:"reinsurer_name"`
	ReinsurerCode     string         `gorm:"not null" json:"reinsurer_code"`
	ContractType      string         `gorm:"not null" json:"contract_type"` // quota_share, surplus, excess_of_loss, stop_loss
	CoverageType      string         `gorm:"not null" json:"coverage_type"` // proportional, non_proportional
	EffectiveDate     time.Time      `gorm:"not null" json:"effective_date"`
	ExpiryDate        time.Time      `gorm:"not null" json:"expiry_date"`
	CessionPercentage float64        `json:"cession_percentage"` // for quota share
	RetentionLimit    float64        `json:"retention_limit"`    // for surplus and excess of loss
	CoverageLimit     float64        `json:"coverage_limit"`
	Priority          int            `json:"priority"` // layer priority for excess of loss
	CommissionRate    float64        `json:"commission_rate"`
	ProfitCommission  float64        `json:"profit_commission"`
	LossRatio         float64        `gorm:"default:0" json:"loss_ratio"`
	PremiumCeded      float64        `gorm:"default:0" json:"premium_ceded"`
	ClaimsCeded       float64        `gorm:"default:0" json:"claims_ceded"`
	CommissionEarned  float64        `gorm:"default:0" json:"commission_earned"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	Terms             string         `json:"terms"`      // JSON object with contract terms
	Exclusions        string         `json:"exclusions"` // JSON array of exclusions
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Cessions []ReinsuranceCession `gorm:"foreignKey:ContractID" json:"cessions,omitempty"`
	Claims   []ReinsuranceClaim   `gorm:"foreignKey:ContractID" json:"claims,omitempty"`
}

// ReinsuranceCession represents individual policy cessions to reinsurers
type ReinsuranceCession struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ContractID         uuid.UUID  `gorm:"type:uuid;not null" json:"contract_id"`
	PolicyID           uuid.UUID  `gorm:"type:uuid;not null" json:"policy_id"`
	CessionNumber      string     `gorm:"uniqueIndex;not null" json:"cession_number"`
	CessionDate        time.Time  `gorm:"not null" json:"cession_date"`
	CededPremium       float64    `gorm:"not null" json:"ceded_premium"`
	CededSumInsured    float64    `gorm:"not null" json:"ceded_sum_insured"`
	CessionPercentage  float64    `json:"cession_percentage"`
	RetainedPremium    float64    `json:"retained_premium"`
	RetainedSumInsured float64    `json:"retained_sum_insured"`
	CommissionAmount   float64    `json:"commission_amount"`
	BrokerageAmount    float64    `json:"brokerage_amount"`
	Status             string     `gorm:"not null;default:'active'" json:"status"` // active, cancelled, expired
	CancellationDate   *time.Time `json:"cancellation_date"`
	CancellationReason string     `json:"cancellation_reason"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Contract ReinsuranceContract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	Policy   models.Policy       `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
}

// ReinsuranceClaim represents claims under reinsurance contracts
type ReinsuranceClaim struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ContractID             uuid.UUID  `gorm:"type:uuid;not null" json:"contract_id"`
	OriginalClaimID        uuid.UUID  `gorm:"type:uuid;not null" json:"original_claim_id"`
	ReinsuranceClaimNumber string     `gorm:"uniqueIndex;not null" json:"reinsurance_claim_number"`
	NotificationDate       time.Time  `gorm:"not null" json:"notification_date"`
	LossDate               time.Time  `gorm:"not null" json:"loss_date"`
	GrossClaimAmount       float64    `gorm:"not null" json:"gross_claim_amount"`
	CededClaimAmount       float64    `gorm:"not null" json:"ceded_claim_amount"`
	RetainedClaimAmount    float64    `json:"retained_claim_amount"`
	RecoveryAmount         float64    `gorm:"default:0" json:"recovery_amount"`
	Status                 string     `gorm:"not null;default:'reported'" json:"status"` // reported, acknowledged, settled, disputed
	SettlementDate         *time.Time `json:"settlement_date"`
	PaymentDate            *time.Time `json:"payment_date"`
	DisputeReason          string     `json:"dispute_reason"`
	Documentation          string     `json:"documentation"` // JSON array of document URLs
	CreatedAt              time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Contract      ReinsuranceContract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	OriginalClaim models.Claim        `gorm:"foreignKey:OriginalClaimID" json:"original_claim,omitempty"`
}

// CoinsuranceAgreement represents co-insurance arrangements
type CoinsuranceAgreement struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AgreementNumber string         `gorm:"uniqueIndex;not null" json:"agreement_number"`
	LeadInsurerID   uuid.UUID      `gorm:"type:uuid;not null" json:"lead_insurer_id"`
	PolicyID        uuid.UUID      `gorm:"type:uuid;not null" json:"policy_id"`
	TotalSumInsured float64        `gorm:"not null" json:"total_sum_insured"`
	TotalPremium    float64        `gorm:"not null" json:"total_premium"`
	AgreementType   string         `gorm:"not null" json:"agreement_type"` // quota_share, first_loss, layers
	EffectiveDate   time.Time      `gorm:"not null" json:"effective_date"`
	ExpiryDate      time.Time      `gorm:"not null" json:"expiry_date"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	Terms           string         `json:"terms"` // JSON object
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Policy       models.Policy            `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Participants []CoinsuranceParticipant `gorm:"foreignKey:AgreementID" json:"participants,omitempty"`
	Claims       []CoinsuranceClaim       `gorm:"foreignKey:AgreementID" json:"claims,omitempty"`
}

// CoinsuranceParticipant represents individual insurers in co-insurance
type CoinsuranceParticipant struct {
	ID                      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AgreementID             uuid.UUID `gorm:"type:uuid;not null" json:"agreement_id"`
	InsurerName             string    `gorm:"not null" json:"insurer_name"`
	InsurerCode             string    `gorm:"not null" json:"insurer_code"`
	ParticipationPercentage float64   `gorm:"not null" json:"participation_percentage"`
	ShareOfSumInsured       float64   `gorm:"not null" json:"share_of_sum_insured"`
	ShareOfPremium          float64   `gorm:"not null" json:"share_of_premium"`
	IsLeadInsurer           bool      `gorm:"default:false" json:"is_lead_insurer"`
	ContactInfo             string    `json:"contact_info"`                            // JSON object
	BankingDetails          string    `json:"banking_details"`                         // JSON object
	Status                  string    `gorm:"not null;default:'active'" json:"status"` // active, withdrawn, suspended
	CreatedAt               time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Agreement CoinsuranceAgreement `gorm:"foreignKey:AgreementID" json:"agreement,omitempty"`
}

// CoinsuranceClaim represents claims under co-insurance agreements
type CoinsuranceClaim struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	AgreementID      uuid.UUID  `gorm:"type:uuid;not null" json:"agreement_id"`
	OriginalClaimID  uuid.UUID  `gorm:"type:uuid;not null" json:"original_claim_id"`
	TotalClaimAmount float64    `gorm:"not null" json:"total_claim_amount"`
	OurShare         float64    `gorm:"not null" json:"our_share"`
	OurPercentage    float64    `gorm:"not null" json:"our_percentage"`
	Status           string     `gorm:"not null;default:'notified'" json:"status"` // notified, agreed, settled, disputed
	SettlementDate   *time.Time `json:"settlement_date"`
	DisputeDetails   string     `json:"dispute_details"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Agreement     CoinsuranceAgreement `gorm:"foreignKey:AgreementID" json:"agreement,omitempty"`
	OriginalClaim models.Claim         `gorm:"foreignKey:OriginalClaimID" json:"original_claim,omitempty"`
}

// SubrogationCase represents subrogation recovery cases
type SubrogationCase struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CaseNumber          string         `gorm:"uniqueIndex;not null" json:"case_number"`
	ClaimID             uuid.UUID      `gorm:"type:uuid;not null" json:"claim_id"`
	ThirdPartyName      string         `gorm:"not null" json:"third_party_name"`
	ThirdPartyInsurer   string         `json:"third_party_insurer"`
	ThirdPartyPolicy    string         `json:"third_party_policy"`
	LiabilityPercentage float64        `json:"liability_percentage"`
	ClaimAmount         float64        `gorm:"not null" json:"claim_amount"`
	RecoveryAmount      float64        `gorm:"default:0" json:"recovery_amount"`
	ExpectedRecovery    float64        `json:"expected_recovery"`
	Status              string         `gorm:"not null;default:'initiated'" json:"status"` // initiated, investigating, negotiating, legal_action, settled, closed
	Priority            string         `gorm:"default:'medium'" json:"priority"`           // low, medium, high, urgent
	AssignedLawyer      string         `json:"assigned_lawyer"`
	LegalFirm           string         `json:"legal_firm"`
	CourtCaseNumber     string         `json:"court_case_number"`
	SettlementDate      *time.Time     `json:"settlement_date"`
	ClosureDate         *time.Time     `json:"closure_date"`
	ClosureReason       string         `json:"closure_reason"`
	LegalCosts          float64        `gorm:"default:0" json:"legal_costs"`
	NetRecovery         float64        `gorm:"default:0" json:"net_recovery"`
	Notes               string         `json:"notes"`
	Documents           string         `json:"documents"` // JSON array of document URLs
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Claim      models.Claim          `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	Activities []SubrogationActivity `gorm:"foreignKey:CaseID" json:"activities,omitempty"`
}

// SubrogationActivity represents activities in subrogation cases
type SubrogationActivity struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CaseID         uuid.UUID  `gorm:"type:uuid;not null" json:"case_id"`
	ActivityType   string     `gorm:"not null" json:"activity_type"` // investigation, negotiation, legal_filing, settlement, payment
	Description    string     `gorm:"not null" json:"description"`
	ActivityDate   time.Time  `gorm:"not null" json:"activity_date"`
	PerformedBy    string     `json:"performed_by"`
	Cost           float64    `gorm:"default:0" json:"cost"`
	Outcome        string     `json:"outcome"`
	NextAction     string     `json:"next_action"`
	NextActionDate *time.Time `json:"next_action_date"`
	Documents      string     `json:"documents"` // JSON array
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Case SubrogationCase `gorm:"foreignKey:CaseID" json:"case,omitempty"`
}

// PolicyEndorsement represents mid-term policy changes
type PolicyEndorsement struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	PolicyID             uuid.UUID  `gorm:"type:uuid;not null" json:"policy_id"`
	EndorsementNumber    string     `gorm:"uniqueIndex;not null" json:"endorsement_number"`
	EndorsementType      string     `gorm:"not null" json:"endorsement_type"` // coverage_change, sum_insured_change, beneficiary_change, device_replacement
	RequestDate          time.Time  `gorm:"not null" json:"request_date"`
	EffectiveDate        time.Time  `gorm:"not null" json:"effective_date"`
	Status               string     `gorm:"not null;default:'pending'" json:"status"` // pending, approved, rejected, implemented
	RequestedBy          uuid.UUID  `gorm:"type:uuid;not null" json:"requested_by"`
	ApprovedBy           *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovalDate         *time.Time `json:"approval_date"`
	RejectionReason      string     `json:"rejection_reason"`
	PremiumAdjustment    float64    `gorm:"default:0" json:"premium_adjustment"`
	SumInsuredAdjustment float64    `gorm:"default:0" json:"sum_insured_adjustment"`
	ChangeDetails        string     `json:"change_details"`     // JSON object with specific changes
	PreviousValues       string     `json:"previous_values"`    // JSON object with old values
	NewValues            string     `json:"new_values"`         // JSON object with new values
	DocumentsRequired    string     `json:"documents_required"` // JSON array
	DocumentsReceived    string     `json:"documents_received"` // JSON array
	ImplementationDate   *time.Time `json:"implementation_date"`
	Notes                string     `json:"notes"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Policy          models.Policy `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	RequestedByUser models.User   `gorm:"foreignKey:RequestedBy" json:"requested_by_user,omitempty"`
	ApprovedByUser  *models.User  `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
}

// GracePeriod represents premium payment grace periods
type GracePeriod struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	PolicyID             uuid.UUID  `gorm:"type:uuid;not null" json:"policy_id"`
	DueDate              time.Time  `gorm:"not null" json:"due_date"`
	GraceStartDate       time.Time  `gorm:"not null" json:"grace_start_date"`
	GraceEndDate         time.Time  `gorm:"not null" json:"grace_end_date"`
	GraceDays            int        `gorm:"not null" json:"grace_days"`
	OutstandingAmount    float64    `gorm:"not null" json:"outstanding_amount"`
	Status               string     `gorm:"not null;default:'active'" json:"status"` // active, paid, lapsed, reinstated
	NotificationsSent    int        `gorm:"default:0" json:"notifications_sent"`
	LastNotificationDate *time.Time `json:"last_notification_date"`
	PaymentDate          *time.Time `json:"payment_date"`
	LapseDate            *time.Time `json:"lapse_date"`
	ReinstatementDate    *time.Time `json:"reinstatement_date"`
	ReinstatementFee     float64    `gorm:"default:0" json:"reinstatement_fee"`
	Notes                string     `json:"notes"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Policy models.Policy `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
}

// TableName methods
func (ReinsuranceContract) TableName() string {
	return "reinsurance_contracts"
}

func (ReinsuranceCession) TableName() string {
	return "reinsurance_cessions"
}

func (ReinsuranceClaim) TableName() string {
	return "reinsurance_claims"
}

func (CoinsuranceAgreement) TableName() string {
	return "coinsurance_agreements"
}

func (CoinsuranceParticipant) TableName() string {
	return "coinsurance_participants"
}

func (CoinsuranceClaim) TableName() string {
	return "coinsurance_claims"
}

func (SubrogationCase) TableName() string {
	return "subrogation_cases"
}

func (SubrogationActivity) TableName() string {
	return "subrogation_activities"
}

func (PolicyEndorsement) TableName() string {
	return "policy_endorsements"
}

func (GracePeriod) TableName() string {
	return "grace_periods"
}

// BeforeCreate hooks
func (rc *ReinsuranceContract) BeforeCreate(tx *gorm.DB) error {
	if rc.ID == uuid.Nil {
		rc.ID = uuid.New()
	}
	if rc.ContractNumber == "" {
		rc.ContractNumber = "RC-" + uuid.New().String()[:8]
	}
	return nil
}

func (rces *ReinsuranceCession) BeforeCreate(tx *gorm.DB) error {
	if rces.ID == uuid.Nil {
		rces.ID = uuid.New()
	}
	if rces.CessionNumber == "" {
		rces.CessionNumber = "CESS-" + uuid.New().String()[:8]
	}
	return nil
}

func (rcl *ReinsuranceClaim) BeforeCreate(tx *gorm.DB) error {
	if rcl.ID == uuid.Nil {
		rcl.ID = uuid.New()
	}
	if rcl.ReinsuranceClaimNumber == "" {
		rcl.ReinsuranceClaimNumber = "RCL-" + uuid.New().String()[:8]
	}
	return nil
}

func (ca *CoinsuranceAgreement) BeforeCreate(tx *gorm.DB) error {
	if ca.ID == uuid.Nil {
		ca.ID = uuid.New()
	}
	if ca.AgreementNumber == "" {
		ca.AgreementNumber = "COINS-" + uuid.New().String()[:8]
	}
	return nil
}

func (cp *CoinsuranceParticipant) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == uuid.Nil {
		cp.ID = uuid.New()
	}
	return nil
}

func (cc *CoinsuranceClaim) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	return nil
}

func (sc *SubrogationCase) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	if sc.CaseNumber == "" {
		sc.CaseNumber = "SUB-" + uuid.New().String()[:8]
	}
	return nil
}

func (sa *SubrogationActivity) BeforeCreate(tx *gorm.DB) error {
	if sa.ID == uuid.Nil {
		sa.ID = uuid.New()
	}
	return nil
}

func (pe *PolicyEndorsement) BeforeCreate(tx *gorm.DB) error {
	if pe.ID == uuid.Nil {
		pe.ID = uuid.New()
	}
	if pe.EndorsementNumber == "" {
		pe.EndorsementNumber = "END-" + uuid.New().String()[:8]
	}
	return nil
}

func (gp *GracePeriod) BeforeCreate(tx *gorm.DB) error {
	if gp.ID == uuid.Nil {
		gp.ID = uuid.New()
	}
	return nil
}

// Business logic methods for ReinsuranceContract
func (rc *ReinsuranceContract) CheckActive() bool {
	now := time.Now()
	return rc.IsActive && now.After(rc.EffectiveDate) && now.Before(rc.ExpiryDate)
}

func (rc *ReinsuranceContract) CalculateCession(grossPremium, grossSumInsured float64) (cededPremium, cededSumInsured float64) {
	switch rc.ContractType {
	case "quota_share":
		cededPremium = grossPremium * (rc.CessionPercentage / 100)
		cededSumInsured = grossSumInsured * (rc.CessionPercentage / 100)
	case "surplus":
		if grossSumInsured > rc.RetentionLimit {
			excess := grossSumInsured - rc.RetentionLimit
			cededSumInsured = excess
			cededPremium = grossPremium * (excess / grossSumInsured)
		}
	case "excess_of_loss":
		// Excess of loss is claim-based, not premium-based
		cededPremium = 0
		cededSumInsured = 0
	}
	return
}

func (rc *ReinsuranceContract) UpdateLossRatio() {
	if rc.PremiumCeded > 0 {
		rc.LossRatio = rc.ClaimsCeded / rc.PremiumCeded
	}
}

// Business logic methods for SubrogationCase
func (sc *SubrogationCase) CalculateNetRecovery() {
	sc.NetRecovery = sc.RecoveryAmount - sc.LegalCosts
}

func (sc *SubrogationCase) IsEconomicallyViable() bool {
	expectedNet := sc.ExpectedRecovery - sc.LegalCosts
	return expectedNet > (sc.ClaimAmount * 0.1) // At least 10% of claim amount
}

func (sc *SubrogationCase) Close(reason string) {
	sc.Status = "closed"
	now := time.Now()
	sc.ClosureDate = &now
	sc.ClosureReason = reason
	sc.CalculateNetRecovery()
}

// Business logic methods for PolicyEndorsement
func (pe *PolicyEndorsement) Approve(approverID uuid.UUID) {
	pe.Status = "approved"
	pe.ApprovedBy = &approverID
	now := time.Now()
	pe.ApprovalDate = &now
}

func (pe *PolicyEndorsement) Reject(reason string) {
	pe.Status = "rejected"
	pe.RejectionReason = reason
}

func (pe *PolicyEndorsement) Implement() {
	pe.Status = "implemented"
	now := time.Now()
	pe.ImplementationDate = &now
}

// Business logic methods for GracePeriod
func (gp *GracePeriod) IsActive() bool {
	now := time.Now()
	return gp.Status == "active" && now.Before(gp.GraceEndDate)
}

func (gp *GracePeriod) IsExpired() bool {
	return time.Now().After(gp.GraceEndDate) && gp.Status == "active"
}

func (gp *GracePeriod) MarkPaid() {
	gp.Status = "paid"
	now := time.Now()
	gp.PaymentDate = &now
}

func (gp *GracePeriod) MarkLapsed() {
	gp.Status = "lapsed"
	now := time.Now()
	gp.LapseDate = &now
}

func (gp *GracePeriod) Reinstate(fee float64) {
	gp.Status = "reinstated"
	gp.ReinstatementFee = fee
	now := time.Now()
	gp.ReinstatementDate = &now
}
