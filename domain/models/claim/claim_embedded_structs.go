package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ClaimIdentification contains core identification fields
type ClaimIdentification struct {
	ClaimNumber     string     `gorm:"uniqueIndex;not null" json:"claim_number"`
	ReferenceNumber string     `gorm:"type:varchar(100)" json:"reference_number"`
	ClaimType       string     `gorm:"type:varchar(50);not null" json:"claim_type"`
	SubType         string     `gorm:"type:varchar(50)" json:"sub_type"`
	Priority        string     `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	Category        string     `gorm:"type:varchar(50)" json:"category"`
	Source          string     `gorm:"type:varchar(50)" json:"source"` // web, mobile, phone, email
	BatchID         string     `gorm:"type:varchar(100)" json:"batch_id"`
	GroupClaimID    *uuid.UUID `gorm:"type:uuid" json:"group_claim_id,omitempty"`
}

// ClaimFinancial contains all financial-related fields
type ClaimFinancial struct {
	ClaimedAmount       float64        `gorm:"not null" json:"claimed_amount"`
	ApprovedAmount      float64        `json:"approved_amount"`
	SettledAmount       float64        `json:"settled_amount"`
	DeductibleAmount    float64        `json:"deductible_amount"`
	CoPaymentAmount     float64        `json:"co_payment_amount"`
	DepreciationApplied float64        `json:"depreciation_applied"` // Added for device value depreciation
	Currency            string         `gorm:"default:'USD'" json:"currency"`
	ExchangeRate        float64        `gorm:"type:decimal(10,6)" json:"exchange_rate"`
	ReserveAmount       float64        `json:"reserve_amount"`
	RecoveryAmount      float64        `json:"recovery_amount"`
	SalvageValue        float64        `json:"salvage_value"`
	SubrogationAmount   float64        `json:"subrogation_amount"`
	TotalPayout         float64        `json:"total_payout"`
	PaymentMethod       string         `json:"payment_method"`
	PaymentReference    string         `json:"payment_reference"`
	TaxAmount           float64        `json:"tax_amount"`
	InterestAmount      float64        `json:"interest_amount"`
	PenaltyAmount       float64        `json:"penalty_amount"`
	CostContainment     datatypes.JSON `gorm:"type:json" json:"cost_containment"`
}

// ClaimLifecycle contains status and timeline information
type ClaimLifecycle struct {
	Status                 string         `gorm:"type:varchar(20);not null;default:'submitted'" json:"status"`
	PreviousStatus         string         `gorm:"type:varchar(20)" json:"previous_status"`
	StatusHistory          datatypes.JSON `gorm:"type:json" json:"status_history"`
	IncidentDate           time.Time      `gorm:"not null" json:"incident_date"`
	ReportedDate           time.Time      `gorm:"autoCreateTime" json:"reported_date"`
	NotificationDate       *time.Time     `json:"notification_date,omitempty"`
	AssignmentDate         *time.Time     `json:"assignment_date,omitempty"`
	FirstContactDate       *time.Time     `json:"first_contact_date,omitempty"`
	InvestigationStartDate *time.Time     `json:"investigation_start_date,omitempty"`
	InvestigationEndDate   *time.Time     `json:"investigation_end_date,omitempty"`
	ApprovalDate           *time.Time     `json:"approval_date,omitempty"`
	DenialDate             *time.Time     `json:"denial_date,omitempty"`
	SettlementDate         *time.Time     `json:"settlement_date,omitempty"`
	ClosedDate             *time.Time     `json:"closed_date,omitempty"`
	ReopenedDate           *time.Time     `json:"reopened_date,omitempty"`
	EstimatedSettlement    time.Time      `json:"estimated_settlement"`
	SLADeadline            time.Time      `json:"sla_deadline"`
	LastUpdateDate         time.Time      `json:"last_update_date"`
}

// ClaimInvestigation contains investigation and fraud-related fields
type ClaimInvestigation struct {
	RequiresInvestigation bool           `gorm:"default:false" json:"requires_investigation"`
	InvestigationStatus   string         `json:"investigation_status"`
	InvestigationType     string         `json:"investigation_type"`
	InvestigatorID        *uuid.UUID     `gorm:"type:uuid" json:"investigator_id,omitempty"`
	InvestigationNotes    string         `gorm:"type:text" json:"investigation_notes"`
	FraudScore            float64        `gorm:"default:0.0" json:"fraud_score"`
	FraudRiskLevel        string         `json:"fraud_risk_level"`
	FraudIndicators       datatypes.JSON `gorm:"type:json" json:"fraud_indicators"`
	RedFlags              datatypes.JSON `gorm:"type:json" json:"red_flags"`
	EvidenceCollected     datatypes.JSON `gorm:"type:json" json:"evidence_collected"`
	WitnessStatements     datatypes.JSON `gorm:"type:json" json:"witness_statements"`
	ExpertOpinions        datatypes.JSON `gorm:"type:json" json:"expert_opinions"`
	ForensicAnalysis      datatypes.JSON `gorm:"type:json" json:"forensic_analysis"`
	SurveillanceReports   datatypes.JSON `gorm:"type:json" json:"surveillance_reports"`
	BackgroundChecks      datatypes.JSON `gorm:"type:json" json:"background_checks"`
}

// ClaimSettlement contains settlement and payment details
type ClaimSettlement struct {
	SettlementType       string         `json:"settlement_type"`
	SettlementStatus     string         `json:"settlement_status"`
	SettledAt            *time.Time     `json:"settled_at,omitempty"`
	PayoutDate           *time.Time     `json:"payout_date,omitempty"`
	PayoutMethod         string         `json:"payout_method"`
	PayoutStatus         string         `json:"payout_status"`
	BankAccountDetails   datatypes.JSON `gorm:"type:json" json:"bank_account_details"`
	CheckNumber          string         `json:"check_number"`
	TransactionID        string         `json:"transaction_id"`
	RecoveryStatus       string         `json:"recovery_status"`
	RecoveryDetails      datatypes.JSON `gorm:"type:json" json:"recovery_details"`
	SubrogationStatus    string         `json:"subrogation_status"`
	SubrogationDetails   datatypes.JSON `gorm:"type:json" json:"subrogation_details"`
	VendorPayments       datatypes.JSON `gorm:"type:json" json:"vendor_payments"`
	ReimbursementDetails datatypes.JSON `gorm:"type:json" json:"reimbursement_details"`
	SettlementAgreement  datatypes.JSON `gorm:"type:json" json:"settlement_agreement"`
	ReleaseForm          datatypes.JSON `gorm:"type:json" json:"release_form"`
}

// ClaimDocumentation contains document and evidence management
type ClaimDocumentation struct {
	Description         string         `gorm:"not null" json:"description"`
	DetailedDescription string         `gorm:"type:text" json:"detailed_description"`
	IncidentLocation    string         `json:"incident_location"`
	LocationCoordinates datatypes.JSON `gorm:"type:json" json:"location_coordinates"`
	PoliceReportNumber  string         `json:"police_report_number"`
	PoliceStation       string         `json:"police_station"`
	WitnessContact      string         `json:"witness_contact"`
	WitnessCount        int            `json:"witness_count"`
	DocumentsReceived   datatypes.JSON `gorm:"type:json" json:"documents_received"`
	DocumentsPending    datatypes.JSON `gorm:"type:json" json:"documents_pending"`
	PhotoEvidence       datatypes.JSON `gorm:"type:json" json:"photo_evidence"`
	VideoEvidence       datatypes.JSON `gorm:"type:json" json:"video_evidence"`
	AudioRecordings     datatypes.JSON `gorm:"type:json" json:"audio_recordings"`
	DamageAssessment    datatypes.JSON `gorm:"type:json" json:"damage_assessment"`
	RepairEstimates     datatypes.JSON `gorm:"type:json" json:"repair_estimates"`
	MedicalReports      datatypes.JSON `gorm:"type:json" json:"medical_reports"`
	ThirdPartyReports   datatypes.JSON `gorm:"type:json" json:"third_party_reports"`
}

// ClaimAssignment contains assignment and ownership details
type ClaimAssignment struct {
	AssignedTo        *uuid.UUID     `gorm:"type:uuid" json:"assigned_to,omitempty"`
	AssignedTeam      string         `json:"assigned_team"`
	AssignmentQueue   string         `json:"assignment_queue"`
	Adjuster          *uuid.UUID     `gorm:"type:uuid" json:"adjuster_id,omitempty"`
	AdjusterNotes     string         `gorm:"type:text" json:"adjuster_notes"`
	Supervisor        *uuid.UUID     `gorm:"type:uuid" json:"supervisor_id,omitempty"`
	EscalatedTo       *uuid.UUID     `gorm:"type:uuid" json:"escalated_to,omitempty"`
	EscalationReason  string         `json:"escalation_reason"`
	WorkloadScore     float64        `json:"workload_score"`
	AssignmentHistory datatypes.JSON `gorm:"type:json" json:"assignment_history"`
	HandoffNotes      datatypes.JSON `gorm:"type:json" json:"handoff_notes"`
}

// ClaimCompliance contains regulatory and compliance information
type ClaimCompliance struct {
	ComplianceStatus      string         `json:"compliance_status"`
	RegulatoryRegion      string         `json:"regulatory_region"`
	ReportingRequired     bool           `gorm:"default:false" json:"reporting_required"`
	ReportedToRegulator   bool           `gorm:"default:false" json:"reported_to_regulator"`
	RegulatoryDeadline    *time.Time     `json:"regulatory_deadline,omitempty"`
	ComplianceChecks      datatypes.JSON `gorm:"type:json" json:"compliance_checks"`
	AuditTrail            datatypes.JSON `gorm:"type:json" json:"audit_trail"`
	DataRetentionDate     time.Time      `json:"data_retention_date"`
	PrivacyCompliance     datatypes.JSON `gorm:"type:json" json:"privacy_compliance"`
	LegalHold             bool           `gorm:"default:false" json:"legal_hold"`
	LitigationStatus      string         `json:"litigation_status"`
	LegalRepresentation   datatypes.JSON `gorm:"type:json" json:"legal_representation"`
	RegulatoryFines       float64        `json:"regulatory_fines"`
	ComplianceCertificate string         `json:"compliance_certificate"`
}

// ClaimMetrics contains performance and analytics metrics
type ClaimMetrics struct {
	ProcessingTime        int            `json:"processing_time_hours"`
	TouchPoints           int            `json:"touch_points"`
	CustomerInteractions  int            `json:"customer_interactions"`
	DocumentsProcessed    int            `json:"documents_processed"`
	CycleTime             int            `json:"cycle_time_days"`
	ResponseTime          int            `json:"response_time_hours"`
	ResolutionTime        int            `json:"resolution_time_days"`
	CustomerSatisfaction  float64        `json:"customer_satisfaction"`
	NPSScore              int            `json:"nps_score"`
	QualityScore          float64        `json:"quality_score"`
	EfficiencyScore       float64        `json:"efficiency_score"`
	CostEffectiveness     float64        `json:"cost_effectiveness"`
	SLACompliance         bool           `gorm:"default:true" json:"sla_compliance"`
	PerformanceIndicators datatypes.JSON `gorm:"type:json" json:"performance_indicators"`
}
