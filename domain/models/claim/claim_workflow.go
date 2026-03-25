package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ClaimWorkflow represents the workflow and state machine for claims
type ClaimWorkflow struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Workflow Definition
	WorkflowType    string         `gorm:"type:varchar(50)" json:"workflow_type"`
	WorkflowVersion string         `gorm:"type:varchar(20)" json:"workflow_version"`
	CurrentState    string         `gorm:"type:varchar(50)" json:"current_state"`
	PreviousState   string         `gorm:"type:varchar(50)" json:"previous_state"`
	StateHistory    datatypes.JSON `gorm:"type:json" json:"state_history"` // []StateTransition

	// Process Steps
	TotalSteps     int            `json:"total_steps"`
	CompletedSteps int            `json:"completed_steps"`
	CurrentStep    string         `gorm:"type:varchar(100)" json:"current_step"`
	StepProgress   datatypes.JSON `gorm:"type:json" json:"step_progress"`   // map[string]StepStatus
	PendingActions datatypes.JSON `gorm:"type:json" json:"pending_actions"` // []Action

	// SLA Management
	SLAType              string     `gorm:"type:varchar(50)" json:"sla_type"`
	SLAStartTime         time.Time  `json:"sla_start_time"`
	SLATargetTime        time.Time  `json:"sla_target_time"`
	SLABreached          bool       `gorm:"default:false" json:"sla_breached"`
	SLABreachTime        *time.Time `json:"sla_breach_time,omitempty"`
	TimeRemaining        int        `json:"time_remaining_hours"`
	BusinessHoursElapsed int        `json:"business_hours_elapsed"`

	// Approvals
	RequiresApproval     bool           `gorm:"default:false" json:"requires_approval"`
	ApprovalLevel        int            `json:"approval_level"`
	ApprovalChain        datatypes.JSON `gorm:"type:json" json:"approval_chain"` // []Approver
	CurrentApprover      *uuid.UUID     `gorm:"type:uuid" json:"current_approver,omitempty"`
	ApprovalHistory      datatypes.JSON `gorm:"type:json" json:"approval_history"` // []ApprovalRecord
	AutoApprovalEligible bool           `gorm:"default:false" json:"auto_approval_eligible"`

	// Escalations
	EscalationLevel    int            `json:"escalation_level"`
	EscalationTriggers datatypes.JSON `gorm:"type:json" json:"escalation_triggers"` // []Trigger
	EscalationHistory  datatypes.JSON `gorm:"type:json" json:"escalation_history"`  // []Escalation
	NextEscalationTime *time.Time     `json:"next_escalation_time,omitempty"`

	// Rules & Conditions
	BusinessRules  datatypes.JSON `gorm:"type:json" json:"business_rules"`  // []Rule
	RuleViolations datatypes.JSON `gorm:"type:json" json:"rule_violations"` // []Violation
	Preconditions  datatypes.JSON `gorm:"type:json" json:"preconditions"`
	Postconditions datatypes.JSON `gorm:"type:json" json:"postconditions"`

	// Automation
	AutomationEnabled bool           `gorm:"default:true" json:"automation_enabled"`
	AutomatedActions  datatypes.JSON `gorm:"type:json" json:"automated_actions"` // []AutoAction
	ManualOverrides   datatypes.JSON `gorm:"type:json" json:"manual_overrides"`  // []Override
	ScriptExecutions  datatypes.JSON `gorm:"type:json" json:"script_executions"` // []ScriptRun

	// Notifications
	NotificationsSent    int            `json:"notifications_sent"`
	NotificationQueue    datatypes.JSON `gorm:"type:json" json:"notification_queue"` // []Notification
	LastNotificationTime *time.Time     `json:"last_notification_time,omitempty"`

	// Metrics
	CycleTime    int `json:"cycle_time_hours"`
	WaitTime     int `json:"wait_time_hours"`
	WorkTime     int `json:"work_time_hours"`
	TouchTime    int `json:"touch_time_minutes"`
	HandoffCount int `json:"handoff_count"`
	ReworkCount  int `json:"rework_count"`

	// Status
	WorkflowStatus string    `gorm:"type:varchar(50)" json:"workflow_status"`
	BlockedReason  string    `gorm:"type:text" json:"blocked_reason"`
	Priority       int       `json:"priority"`
	CreatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ClaimSettlementDetail represents detailed settlement information
type ClaimSettlementDetail struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Settlement Calculation
	ClaimedAmount       float64 `gorm:"type:decimal(15,2)" json:"claimed_amount"`
	AssessedAmount      float64 `gorm:"type:decimal(15,2)" json:"assessed_amount"`
	ApprovedAmount      float64 `gorm:"type:decimal(15,2)" json:"approved_amount"`
	DeductibleApplied   float64 `gorm:"type:decimal(15,2)" json:"deductible_applied"`
	CopaymentApplied    float64 `gorm:"type:decimal(15,2)" json:"copayment_applied"`
	DepreciationApplied float64 `gorm:"type:decimal(15,2)" json:"depreciation_applied"`
	LimitAdjustment     float64 `gorm:"type:decimal(15,2)" json:"limit_adjustment"`
	FinalSettlement     float64 `gorm:"type:decimal(15,2)" json:"final_settlement"`

	// Payment Details
	PaymentMethod        string         `gorm:"type:varchar(50)" json:"payment_method"`
	PaymentFrequency     string         `gorm:"type:varchar(50)" json:"payment_frequency"` // lump_sum, installments
	NumberOfInstallments int            `json:"number_of_installments"`
	InstallmentAmount    float64        `gorm:"type:decimal(15,2)" json:"installment_amount"`
	PaymentSchedule      datatypes.JSON `gorm:"type:json" json:"payment_schedule"` // []PaymentSchedule

	// Banking Information
	BankName      string `gorm:"type:varchar(255)" json:"bank_name"`
	AccountNumber string `gorm:"type:varchar(100)" json:"account_number"`
	RoutingNumber string `gorm:"type:varchar(50)" json:"routing_number"`
	IBAN          string `gorm:"type:varchar(50)" json:"iban"`
	SWIFTCode     string `gorm:"type:varchar(20)" json:"swift_code"`

	// Tax Implications
	TaxableAmount float64 `gorm:"type:decimal(15,2)" json:"taxable_amount"`
	TaxRate       float64 `gorm:"type:decimal(5,2)" json:"tax_rate"`
	TaxWithheld   float64 `gorm:"type:decimal(15,2)" json:"tax_withheld"`
	TaxFormSent   bool    `gorm:"default:false" json:"tax_form_sent"`

	// Settlement Agreement
	AgreementNumber       string         `gorm:"type:varchar(100)" json:"agreement_number"`
	AgreementDate         time.Time      `json:"agreement_date"`
	AgreementTerms        datatypes.JSON `gorm:"type:json" json:"agreement_terms"`
	ReleaseObtained       bool           `gorm:"default:false" json:"release_obtained"`
	ReleaseDate           *time.Time     `json:"release_date,omitempty"`
	ConfidentialityClause bool           `gorm:"default:false" json:"confidentiality_clause"`

	// Vendor Payments
	VendorPayments    datatypes.JSON `gorm:"type:json" json:"vendor_payments"` // []VendorPayment
	RepairShopPayment float64        `gorm:"type:decimal(15,2)" json:"repair_shop_payment"`
	MedicalPayment    float64        `gorm:"type:decimal(15,2)" json:"medical_payment"`
	LegalFees         float64        `gorm:"type:decimal(15,2)" json:"legal_fees"`
	OtherExpenses     float64        `gorm:"type:decimal(15,2)" json:"other_expenses"`

	// Recovery & Subrogation
	SubrogationPotential float64    `gorm:"type:decimal(15,2)" json:"subrogation_potential"`
	SubrogationInitiated bool       `gorm:"default:false" json:"subrogation_initiated"`
	RecoveredAmount      float64    `gorm:"type:decimal(15,2)" json:"recovered_amount"`
	RecoveryDate         *time.Time `json:"recovery_date,omitempty"`
	RecoveryCost         float64    `gorm:"type:decimal(15,2)" json:"recovery_cost"`
	NetRecovery          float64    `gorm:"type:decimal(15,2)" json:"net_recovery"`

	// Salvage
	SalvageValue     float64    `gorm:"type:decimal(15,2)" json:"salvage_value"`
	SalvageRecovered bool       `gorm:"default:false" json:"salvage_recovered"`
	SalvageVendor    string     `gorm:"type:varchar(255)" json:"salvage_vendor"`
	SalvageDate      *time.Time `json:"salvage_date,omitempty"`

	// Settlement Status
	SettlementStatus     string     `gorm:"type:varchar(50)" json:"settlement_status"`
	SettlementDate       time.Time  `json:"settlement_date"`
	PaymentProcessedDate *time.Time `json:"payment_processed_date,omitempty"`
	PaymentClearedDate   *time.Time `json:"payment_cleared_date,omitempty"`

	// Quality & Audit
	AuditRequired  bool           `gorm:"default:false" json:"audit_required"`
	AuditCompleted bool           `gorm:"default:false" json:"audit_completed"`
	AuditFindings  datatypes.JSON `gorm:"type:json" json:"audit_findings"`
	QualityScore   float64        `gorm:"type:decimal(5,2)" json:"quality_score"`

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	SettledBy uuid.UUID `gorm:"type:uuid" json:"settled_by"`
}

// =====================================
// METHODS
// =====================================

// IsSLABreached checks if SLA has been breached
func (cw *ClaimWorkflow) IsSLABreached() bool {
	return cw.SLABreached || time.Now().After(cw.SLATargetTime)
}

// NeedsEscalation checks if workflow needs escalation
func (cw *ClaimWorkflow) NeedsEscalation() bool {
	if cw.NextEscalationTime != nil {
		return time.Now().After(*cw.NextEscalationTime)
	}
	return cw.IsSLABreached()
}

// GetCompletionPercentage calculates workflow completion
func (cw *ClaimWorkflow) GetCompletionPercentage() float64 {
	if cw.TotalSteps == 0 {
		return 0
	}
	return float64(cw.CompletedSteps) / float64(cw.TotalSteps) * 100
}

// IsBlocked checks if workflow is blocked
func (cw *ClaimWorkflow) IsBlocked() bool {
	return cw.WorkflowStatus == "blocked" || cw.BlockedReason != ""
}

// IsFullySettled checks if claim is fully settled
func (cs *ClaimSettlementDetail) IsFullySettled() bool {
	return cs.SettlementStatus == "completed" &&
		cs.PaymentClearedDate != nil &&
		cs.ReleaseObtained
}

// HasRecoveryPotential checks if there's recovery potential
func (cs *ClaimSettlementDetail) HasRecoveryPotential() bool {
	return cs.SubrogationPotential > 0 || cs.SalvageValue > 0
}

// GetNetPayment calculates net payment after recoveries
func (cs *ClaimSettlementDetail) GetNetPayment() float64 {
	return cs.FinalSettlement - cs.RecoveredAmount - cs.SalvageValue
}

// RequiresAudit checks if settlement requires audit
func (cs *ClaimSettlementDetail) RequiresAudit() bool {
	return cs.AuditRequired || cs.FinalSettlement > 50000 ||
		cs.SubrogationPotential > 10000
}
