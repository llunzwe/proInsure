package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ClaimPayment represents payment processing for claims
type ClaimPayment struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Payment Identification
	PaymentNumber   string `gorm:"type:varchar(100);unique" json:"payment_number"`
	PaymentType     string `gorm:"type:varchar(50)" json:"payment_type"`     // initial, supplemental, final
	PaymentCategory string `gorm:"type:varchar(50)" json:"payment_category"` // indemnity, expense, recovery
	ReferenceNumber string `gorm:"type:varchar(100)" json:"reference_number"`

	// Amount Details
	RequestedAmount     float64 `gorm:"type:decimal(15,2)" json:"requested_amount"`
	ApprovedAmount      float64 `gorm:"type:decimal(15,2)" json:"approved_amount"`
	ProcessedAmount     float64 `gorm:"type:decimal(15,2)" json:"processed_amount"`
	Currency            string  `gorm:"type:varchar(3)" json:"currency"`
	ExchangeRate        float64 `gorm:"type:decimal(10,6)" json:"exchange_rate"`
	LocalCurrencyAmount float64 `gorm:"type:decimal(15,2)" json:"local_currency_amount"`

	// Payment Method
	PaymentMethod      string `gorm:"type:varchar(50)" json:"payment_method"`
	ElectronicTransfer bool   `gorm:"default:true" json:"electronic_transfer"`
	CheckNumber        string `gorm:"type:varchar(50)" json:"check_number"`
	WireTransferID     string `gorm:"type:varchar(100)" json:"wire_transfer_id"`
	CreditCardLast4    string `gorm:"type:varchar(4)" json:"credit_card_last4"`
	DigitalWalletID    string `gorm:"type:varchar(100)" json:"digital_wallet_id"`

	// Payee Information
	PayeeType    string         `gorm:"type:varchar(50)" json:"payee_type"` // customer, vendor, provider
	PayeeID      uuid.UUID      `gorm:"type:uuid" json:"payee_id"`
	PayeeName    string         `gorm:"type:varchar(255)" json:"payee_name"`
	PayeeAddress datatypes.JSON `gorm:"type:json" json:"payee_address"`
	PayeeTaxID   string         `gorm:"type:varchar(50)" json:"payee_tax_id"`

	// Banking Details
	BankAccountID       string `gorm:"type:varchar(100)" json:"bank_account_id"`
	BankName            string `gorm:"type:varchar(255)" json:"bank_name"`
	AccountType         string `gorm:"type:varchar(50)" json:"account_type"`
	RoutingNumber       string `gorm:"type:varchar(50)" json:"routing_number"`
	AccountNumberMasked string `gorm:"type:varchar(20)" json:"account_number_masked"`

	// Processing Timeline
	RequestedDate time.Time  `json:"requested_date"`
	ApprovedDate  *time.Time `json:"approved_date,omitempty"`
	ProcessedDate *time.Time `json:"processed_date,omitempty"`
	ScheduledDate *time.Time `json:"scheduled_date,omitempty"`
	SentDate      *time.Time `json:"sent_date,omitempty"`
	ClearedDate   *time.Time `json:"cleared_date,omitempty"`
	VoidedDate    *time.Time `json:"voided_date,omitempty"`

	// Transaction Details
	TransactionID      string `gorm:"type:varchar(100)" json:"transaction_id"`
	BatchID            string `gorm:"type:varchar(100)" json:"batch_id"`
	ProcessorReference string `gorm:"type:varchar(100)" json:"processor_reference"`
	AuthorizationCode  string `gorm:"type:varchar(50)" json:"authorization_code"`
	ConfirmationNumber string `gorm:"type:varchar(100)" json:"confirmation_number"`

	// Status Management
	PaymentStatus        string `gorm:"type:varchar(50)" json:"payment_status"`
	ProcessingStatus     string `gorm:"type:varchar(50)" json:"processing_status"`
	ReconciliationStatus string `gorm:"type:varchar(50)" json:"reconciliation_status"`
	HoldReason           string `gorm:"type:text" json:"hold_reason"`
	RejectionReason      string `gorm:"type:text" json:"rejection_reason"`

	// Fees & Deductions
	ProcessingFee   float64 `gorm:"type:decimal(15,2)" json:"processing_fee"`
	TransactionFee  float64 `gorm:"type:decimal(15,2)" json:"transaction_fee"`
	TaxWithheld     float64 `gorm:"type:decimal(15,2)" json:"tax_withheld"`
	OtherDeductions float64 `gorm:"type:decimal(15,2)" json:"other_deductions"`
	NetAmount       float64 `gorm:"type:decimal(15,2)" json:"net_amount"`

	// Compliance & Reporting
	TaxReportable    bool           `gorm:"default:false" json:"tax_reportable"`
	Form1099Required bool           `gorm:"default:false" json:"form_1099_required"`
	Form1099Sent     bool           `gorm:"default:false" json:"form_1099_sent"`
	ComplianceChecks datatypes.JSON `gorm:"type:json" json:"compliance_checks"`
	AMLVerification  bool           `gorm:"default:false" json:"aml_verification"`

	// Audit Trail
	ApprovedBy          *uuid.UUID     `gorm:"type:uuid" json:"approved_by,omitempty"`
	ProcessedBy         *uuid.UUID     `gorm:"type:uuid" json:"processed_by,omitempty"`
	VoidedBy            *uuid.UUID     `gorm:"type:uuid" json:"voided_by,omitempty"`
	ModificationHistory datatypes.JSON `gorm:"type:json" json:"modification_history"`

	// Status
	IsVoid            bool       `gorm:"default:false" json:"is_void"`
	IsReissued        bool       `gorm:"default:false" json:"is_reissued"`
	ReissuedPaymentID *uuid.UUID `gorm:"type:uuid" json:"reissued_payment_id,omitempty"`
	CreatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ClaimReserve represents claim reserve management
type ClaimReserve struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Reserve Amounts
	InitialReserve   float64 `gorm:"type:decimal(15,2)" json:"initial_reserve"`
	CurrentReserve   float64 `gorm:"type:decimal(15,2)" json:"current_reserve"`
	IndemnityReserve float64 `gorm:"type:decimal(15,2)" json:"indemnity_reserve"`
	ExpenseReserve   float64 `gorm:"type:decimal(15,2)" json:"expense_reserve"`
	LegalReserve     float64 `gorm:"type:decimal(15,2)" json:"legal_reserve"`
	RecoveryReserve  float64 `gorm:"type:decimal(15,2)" json:"recovery_reserve"`

	// Reserve History
	ReserveHistory       datatypes.JSON `gorm:"type:json" json:"reserve_history"` // []ReserveChange
	AdjustmentCount      int            `json:"adjustment_count"`
	LastAdjustmentDate   *time.Time     `json:"last_adjustment_date,omitempty"`
	LastAdjustmentAmount float64        `gorm:"type:decimal(15,2)" json:"last_adjustment_amount"`
	LastAdjustmentReason string         `gorm:"type:text" json:"last_adjustment_reason"`

	// Incurred Amounts
	IncurredTotal     float64 `gorm:"type:decimal(15,2)" json:"incurred_total"`
	PaidToDate        float64 `gorm:"type:decimal(15,2)" json:"paid_to_date"`
	OutstandingAmount float64 `gorm:"type:decimal(15,2)" json:"outstanding_amount"`
	UltimateEstimate  float64 `gorm:"type:decimal(15,2)" json:"ultimate_estimate"`

	// Reserve Adequacy
	AdequacyStatus   string  `gorm:"type:varchar(50)" json:"adequacy_status"`
	AdequacyScore    float64 `gorm:"type:decimal(5,2)" json:"adequacy_score"`
	UnderReserved    bool    `gorm:"default:false" json:"under_reserved"`
	OverReserved     bool    `gorm:"default:false" json:"over_reserved"`
	ReserveDeviation float64 `gorm:"type:decimal(10,2)" json:"reserve_deviation"`

	// Approval Requirements
	RequiresApproval  bool       `gorm:"default:false" json:"requires_approval"`
	ApprovalThreshold float64    `gorm:"type:decimal(15,2)" json:"approval_threshold"`
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovalDate      *time.Time `json:"approval_date,omitempty"`
	ApprovalNotes     string     `gorm:"type:text" json:"approval_notes"`

	// IBNR Allocation
	IBNRAmount         float64    `gorm:"type:decimal(15,2)" json:"ibnr_amount"`
	IBNRPercentage     float64    `gorm:"type:decimal(5,2)" json:"ibnr_percentage"`
	IBNRMethodology    string     `gorm:"type:varchar(100)" json:"ibnr_methodology"`
	IBNRLastCalculated *time.Time `json:"ibnr_last_calculated,omitempty"`

	// Status
	ReserveStatus string    `gorm:"type:varchar(50)" json:"reserve_status"`
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy     uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy     uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// ClaimSubrogation represents subrogation and recovery management
type ClaimSubrogation struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Subrogation Details
	SubrogationNumber string `gorm:"type:varchar(100);unique" json:"subrogation_number"`
	SubrogationType   string `gorm:"type:varchar(50)" json:"subrogation_type"`
	LiabilityParty    string `gorm:"type:varchar(255)" json:"liability_party"`
	LiabilityInsurer  string `gorm:"type:varchar(255)" json:"liability_insurer"`
	PolicyNumber      string `gorm:"type:varchar(100)" json:"policy_number"`
	ClaimNumber       string `gorm:"type:varchar(100)" json:"claim_number"`

	// Financial Details
	SubrogationPotential float64 `gorm:"type:decimal(15,2)" json:"subrogation_potential"`
	DemandAmount         float64 `gorm:"type:decimal(15,2)" json:"demand_amount"`
	OfferedAmount        float64 `gorm:"type:decimal(15,2)" json:"offered_amount"`
	SettledAmount        float64 `gorm:"type:decimal(15,2)" json:"settled_amount"`
	RecoveredAmount      float64 `gorm:"type:decimal(15,2)" json:"recovered_amount"`
	LegalCosts           float64 `gorm:"type:decimal(15,2)" json:"legal_costs"`
	NetRecovery          float64 `gorm:"type:decimal(15,2)" json:"net_recovery"`

	// Legal Proceedings
	LegalAction      bool       `gorm:"default:false" json:"legal_action"`
	AttorneyAssigned bool       `gorm:"default:false" json:"attorney_assigned"`
	AttorneyName     string     `gorm:"type:varchar(255)" json:"attorney_name"`
	LawFirm          string     `gorm:"type:varchar(255)" json:"law_firm"`
	CaseNumber       string     `gorm:"type:varchar(100)" json:"case_number"`
	CourtName        string     `gorm:"type:varchar(255)" json:"court_name"`
	FilingDate       *time.Time `json:"filing_date,omitempty"`
	HearingDate      *time.Time `json:"hearing_date,omitempty"`

	// Timeline
	InitiatedDate    time.Time  `json:"initiated_date"`
	DemandLetterDate *time.Time `json:"demand_letter_date,omitempty"`
	ResponseDate     *time.Time `json:"response_date,omitempty"`
	NegotiationStart *time.Time `json:"negotiation_start,omitempty"`
	SettlementDate   *time.Time `json:"settlement_date,omitempty"`
	RecoveryDate     *time.Time `json:"recovery_date,omitempty"`
	ClosedDate       *time.Time `json:"closed_date,omitempty"`

	// Documentation
	DemandLetter        datatypes.JSON `gorm:"type:json" json:"demand_letter"`
	Correspondence      datatypes.JSON `gorm:"type:json" json:"correspondence"`
	SettlementAgreement datatypes.JSON `gorm:"type:json" json:"settlement_agreement"`
	LegalDocuments      datatypes.JSON `gorm:"type:json" json:"legal_documents"`

	// Status
	SubrogationStatus string    `gorm:"type:varchar(50)" json:"subrogation_status"`
	RecoveryStatus    string    `gorm:"type:varchar(50)" json:"recovery_status"`
	DisputeStatus     string    `gorm:"type:varchar(50)" json:"dispute_status"`
	Priority          string    `gorm:"type:varchar(20)" json:"priority"`
	CreatedAt         time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsProcessed checks if payment is processed
func (cp *ClaimPayment) IsProcessed() bool {
	return cp.PaymentStatus == "processed" && cp.ProcessedDate != nil
}

// IsCleared checks if payment has cleared
func (cp *ClaimPayment) IsCleared() bool {
	return cp.PaymentStatus == "cleared" && cp.ClearedDate != nil
}

// RequiresCompliance checks if payment requires compliance checks
func (cp *ClaimPayment) RequiresCompliance() bool {
	return cp.ProcessedAmount > 10000 || cp.TaxReportable || cp.Form1099Required
}

// IsAdequate checks if reserve is adequate
func (cr *ClaimReserve) IsAdequate() bool {
	return cr.AdequacyStatus == "adequate" && !cr.UnderReserved
}

// NeedsAdjustment checks if reserve needs adjustment
func (cr *ClaimReserve) NeedsAdjustment() bool {
	deviation := cr.ReserveDeviation
	return deviation > 20 || deviation < -20 || cr.UnderReserved || cr.OverReserved
}

// GetReserveUtilization calculates reserve utilization
func (cr *ClaimReserve) GetReserveUtilization() float64 {
	if cr.CurrentReserve == 0 {
		return 0
	}
	return (cr.PaidToDate / cr.CurrentReserve) * 100
}

// IsRecoveryViable checks if recovery is viable
func (cs *ClaimSubrogation) IsRecoveryViable() bool {
	return cs.SubrogationPotential > cs.LegalCosts &&
		cs.SubrogationStatus != "abandoned"
}

// GetRecoveryRate calculates recovery rate
func (cs *ClaimSubrogation) GetRecoveryRate() float64 {
	if cs.SubrogationPotential == 0 {
		return 0
	}
	return (cs.RecoveredAmount / cs.SubrogationPotential) * 100
}

// IsLegalActionRequired checks if legal action is needed
func (cs *ClaimSubrogation) IsLegalActionRequired() bool {
	return cs.SubrogationPotential > 50000 ||
		cs.DisputeStatus == "contested" ||
		cs.LegalAction
}
