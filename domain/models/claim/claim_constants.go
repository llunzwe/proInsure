package claim

// Claim Status Constants
const (
	StatusSubmitted        = "submitted"
	StatusUnderReview      = "under_review"
	StatusInvestigating    = "investigating"
	StatusPendingDocuments = "pending_documents"
	StatusPendingApproval  = "pending_approval"
	StatusApproved         = "approved"
	StatusDenied           = "denied"
	StatusSettled          = "settled"
	StatusWithdrawn        = "withdrawn"
	StatusExpired          = "expired"
	StatusAppeal           = "appeal"
	StatusClosed           = "closed"
)

// Claim Type Constants
const (
	TypeDamage           = "damage"
	TypeTheft            = "theft"
	TypeLoss             = "loss"
	TypeMalfunction      = "malfunction"
	TypeWaterDamage      = "water_damage"
	TypeScreenDamage     = "screen_damage"
	TypeBatteryIssue     = "battery_issue"
	TypePhysicalDamage   = "physical_damage"
	TypeSoftwareIssue    = "software_issue"
	TypeAccessoryDamage  = "accessory_damage"
	TypeWarrantyClaim    = "warranty_claim"
	TypeAccidentalDamage = "accidental_damage"
	TypeNaturalDisaster  = "natural_disaster"
)

// Priority Levels
const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityUrgent   = "urgent"
	PriorityCritical = "critical"
)

// Settlement Methods
const (
	SettlementCash         = "cash"
	SettlementBankTransfer = "bank_transfer"
	SettlementCheck        = "check"
	SettlementRepair       = "repair"
	SettlementReplacement  = "replacement"
	SettlementVoucher      = "voucher"
	SettlementCredit       = "credit"
)

// Investigation Status
const (
	InvestigationNotRequired = "not_required"
	InvestigationPending     = "pending"
	InvestigationInProgress  = "in_progress"
	InvestigationCompleted   = "completed"
	InvestigationEscalated   = "escalated"
)

// Fraud Risk Levels
const (
	FraudRiskNone     = "none"
	FraudRiskLow      = "low"
	FraudRiskMedium   = "medium"
	FraudRiskHigh     = "high"
	FraudRiskCritical = "critical"
)

// Document Types
const (
	DocTypeProofOfPurchase  = "proof_of_purchase"
	DocTypePoliceReport     = "police_report"
	DocTypeDamagePhoto      = "damage_photo"
	DocTypeRepairEstimate   = "repair_estimate"
	DocTypeInvoice          = "invoice"
	DocTypeMedicalReport    = "medical_report"
	DocTypeWitnessStatement = "witness_statement"
	DocTypeIncidentReport   = "incident_report"
	DocTypeIdentification   = "identification"
	DocTypePolicyDocument   = "policy_document"
)

// Appeal Status
const (
	AppealNotFiled  = "not_filed"
	AppealSubmitted = "submitted"
	AppealReviewing = "reviewing"
	AppealApproved  = "approved"
	AppealDenied    = "denied"
	AppealWithdrawn = "withdrawn"
)

// Workflow States
const (
	WorkflowInitiated         = "initiated"
	WorkflowDocumentReview    = "document_review"
	WorkflowTechnicalReview   = "technical_review"
	WorkflowFraudCheck        = "fraud_check"
	WorkflowApprovalPending   = "approval_pending"
	WorkflowPaymentProcessing = "payment_processing"
	WorkflowCompleted         = "completed"
)

// SLA Durations (in hours)
const (
	SLAUrgent   = 24
	SLAHigh     = 48
	SLAMedium   = 72
	SLAStandard = 120
	SLALow      = 168
)

// Threshold Constants
const (
	ThresholdAutoApproval     = 500.0
	ThresholdManagerApproval  = 5000.0
	ThresholdDirectorApproval = 20000.0
	ThresholdFraudScore       = 0.7
	ThresholdHighValue        = 10000.0
	ThresholdInvestigation    = 3000.0
)

// Reserve Percentages
const (
	ReserveInitialPercent = 0.8 // 80% of claimed amount
	ReserveMinimum        = 100.0
	ReserveMaximum        = 1000000.0
)
