package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ClaimLitigation represents legal proceedings related to a claim
type ClaimLitigation struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Case Information
	CaseNumber   string `gorm:"type:varchar(100);unique" json:"case_number"`
	CaseType     string `gorm:"type:varchar(50)" json:"case_type"` // civil, arbitration, mediation
	Jurisdiction string `gorm:"type:varchar(100)" json:"jurisdiction"`
	Court        string `gorm:"type:varchar(255)" json:"court"`
	Judge        string `gorm:"type:varchar(255)" json:"judge"`
	Docket       string `gorm:"type:varchar(100)" json:"docket"`

	// Parties
	Plaintiff    string         `gorm:"type:varchar(255)" json:"plaintiff"`
	Defendant    string         `gorm:"type:varchar(255)" json:"defendant"`
	ThirdParties datatypes.JSON `gorm:"type:json" json:"third_parties"` // []Party
	CoDefendants datatypes.JSON `gorm:"type:json" json:"co_defendants"`

	// Legal Representation
	OurAttorney     string         `gorm:"type:varchar(255)" json:"our_attorney"`
	OurLawFirm      string         `gorm:"type:varchar(255)" json:"our_law_firm"`
	OpposingCounsel string         `gorm:"type:varchar(255)" json:"opposing_counsel"`
	OpposingLawFirm string         `gorm:"type:varchar(255)" json:"opposing_law_firm"`
	AttorneyContact datatypes.JSON `gorm:"type:json" json:"attorney_contact"`

	// Case Timeline
	FilingDate         time.Time  `json:"filing_date"`
	ServiceDate        *time.Time `json:"service_date,omitempty"`
	AnswerDueDate      *time.Time `json:"answer_due_date,omitempty"`
	AnswerFiledDate    *time.Time `json:"answer_filed_date,omitempty"`
	DiscoveryStartDate *time.Time `json:"discovery_start_date,omitempty"`
	DiscoveryEndDate   *time.Time `json:"discovery_end_date,omitempty"`
	MediationDate      *time.Time `json:"mediation_date,omitempty"`
	TrialDate          *time.Time `json:"trial_date,omitempty"`
	SettlementDate     *time.Time `json:"settlement_date,omitempty"`
	JudgmentDate       *time.Time `json:"judgment_date,omitempty"`
	AppealDeadline     *time.Time `json:"appeal_deadline,omitempty"`

	// Financial Exposure
	ClaimedDamages    float64 `gorm:"type:decimal(15,2)" json:"claimed_damages"`
	PunitiveDamages   float64 `gorm:"type:decimal(15,2)" json:"punitive_damages"`
	AttorneyFees      float64 `gorm:"type:decimal(15,2)" json:"attorney_fees"`
	CourtCosts        float64 `gorm:"type:decimal(15,2)" json:"court_costs"`
	ExpertWitnessFees float64 `gorm:"type:decimal(15,2)" json:"expert_witness_fees"`
	DiscoveryCosts    float64 `gorm:"type:decimal(15,2)" json:"discovery_costs"`
	TotalExposure     float64 `gorm:"type:decimal(15,2)" json:"total_exposure"`
	ReserveAmount     float64 `gorm:"type:decimal(15,2)" json:"reserve_amount"`

	// Settlement
	SettlementOffered     bool           `gorm:"default:false" json:"settlement_offered"`
	SettlementAmount      float64        `gorm:"type:decimal(15,2)" json:"settlement_amount"`
	SettlementTerms       datatypes.JSON `gorm:"type:json" json:"settlement_terms"`
	SettlementApproved    bool           `gorm:"default:false" json:"settlement_approved"`
	ConfidentialityClause bool           `gorm:"default:false" json:"confidentiality_clause"`

	// Discovery
	InterrogatoriesSent int `json:"interrogatories_sent"`
	InterrogatoriesRcvd int `json:"interrogatories_received"`
	DepositionsTaken    int `json:"depositions_taken"`
	DepositionsDefended int `json:"depositions_defended"`
	DocumentRequests    int `json:"document_requests"`
	DocumentsProduced   int `json:"documents_produced"`
	AdmissionRequests   int `json:"admission_requests"`

	// Motions
	MotionsFiled   datatypes.JSON `gorm:"type:json" json:"motions_filed"` // []Motion
	MotionsPending int            `json:"motions_pending"`
	MotionsWon     int            `json:"motions_won"`
	MotionsLost    int            `json:"motions_lost"`

	// Expert Witnesses
	ExpertWitnesses   datatypes.JSON `gorm:"type:json" json:"expert_witnesses"` // []Expert
	ExpertReports     datatypes.JSON `gorm:"type:json" json:"expert_reports"`
	DaubertChallenges datatypes.JSON `gorm:"type:json" json:"daubert_challenges"`

	// Trial
	JuryTrial        bool           `gorm:"default:false" json:"jury_trial"`
	BenchTrial       bool           `gorm:"default:false" json:"bench_trial"`
	TrialDuration    int            `json:"trial_duration_days"`
	Verdict          string         `gorm:"type:varchar(50)" json:"verdict"`
	JudgmentAmount   float64        `gorm:"type:decimal(15,2)" json:"judgment_amount"`
	PostTrialMotions datatypes.JSON `gorm:"type:json" json:"post_trial_motions"`

	// Appeals
	AppealFiled    bool       `gorm:"default:false" json:"appeal_filed"`
	AppealDate     *time.Time `json:"appeal_date,omitempty"`
	AppellateCoart string     `gorm:"type:varchar(255)" json:"appellate_court"`
	AppealStatus   string     `gorm:"type:varchar(50)" json:"appeal_status"`
	AppealDecision string     `gorm:"type:text" json:"appeal_decision"`

	// Risk Assessment
	LiabilityAssessment string  `gorm:"type:varchar(50)" json:"liability_assessment"` // strong, moderate, weak
	SuccessProbability  float64 `gorm:"type:decimal(5,2)" json:"success_probability"`
	RiskLevel           string  `gorm:"type:varchar(20)" json:"risk_level"`
	RecommendedStrategy string  `gorm:"type:text" json:"recommended_strategy"`

	// Documentation
	PleadingsFile      string         `gorm:"type:varchar(500)" json:"pleadings_file"`
	DiscoveryFile      string         `gorm:"type:varchar(500)" json:"discovery_file"`
	CorrespondenceFile string         `gorm:"type:varchar(500)" json:"correspondence_file"`
	EvidenceInventory  datatypes.JSON `gorm:"type:json" json:"evidence_inventory"`

	// Status
	LitigationStatus string     `gorm:"type:varchar(50)" json:"litigation_status"`
	Priority         string     `gorm:"type:varchar(20)" json:"priority"`
	NextDeadline     *time.Time `json:"next_deadline,omitempty"`
	NextAction       string     `gorm:"type:text" json:"next_action"`
	CreatedAt        time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	ManagedBy        uuid.UUID  `gorm:"type:uuid" json:"managed_by"`
}

// ClaimArbitration represents arbitration proceedings
type ClaimArbitration struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Arbitration Details
	ArbitrationNumber   string `gorm:"type:varchar(100);unique" json:"arbitration_number"`
	ArbitrationForum    string `gorm:"type:varchar(255)" json:"arbitration_forum"` // AAA, JAMS, etc.
	ArbitrationRules    string `gorm:"type:varchar(255)" json:"arbitration_rules"`
	ArbitrationLocation string `gorm:"type:varchar(255)" json:"arbitration_location"`

	// Arbitrators
	SingleArbitrator bool           `gorm:"default:true" json:"single_arbitrator"`
	ArbitratorName   string         `gorm:"type:varchar(255)" json:"arbitrator_name"`
	ArbitratorPanel  datatypes.JSON `gorm:"type:json" json:"arbitrator_panel"` // []Arbitrator
	ArbitratorFees   float64        `gorm:"type:decimal(15,2)" json:"arbitrator_fees"`

	// Process
	DemandDate          time.Time      `json:"demand_date"`
	ResponseDate        *time.Time     `json:"response_date,omitempty"`
	ArbitratorAppointed *time.Time     `json:"arbitrator_appointed,omitempty"`
	PreliminaryHearing  *time.Time     `json:"preliminary_hearing,omitempty"`
	HearingDates        datatypes.JSON `gorm:"type:json" json:"hearing_dates"` // []Date
	BriefsDue           *time.Time     `json:"briefs_due,omitempty"`

	// Claims & Counterclaims
	ClaimsAmount        float64        `gorm:"type:decimal(15,2)" json:"claims_amount"`
	CounterclaimsAmount float64        `gorm:"type:decimal(15,2)" json:"counterclaims_amount"`
	ClaimsDescription   string         `gorm:"type:text" json:"claims_description"`
	DefensesRaised      datatypes.JSON `gorm:"type:json" json:"defenses_raised"`

	// Discovery
	DiscoveryPermitted bool       `gorm:"default:false" json:"discovery_permitted"`
	DiscoveryDeadline  *time.Time `json:"discovery_deadline,omitempty"`
	DocumentsExchanged int        `json:"documents_exchanged"`
	WitnessesDisclosed int        `json:"witnesses_disclosed"`

	// Award
	AwardDate       *time.Time     `json:"award_date,omitempty"`
	AwardAmount     float64        `gorm:"type:decimal(15,2)" json:"award_amount"`
	AwardType       string         `gorm:"type:varchar(50)" json:"award_type"` // final, interim, partial
	AwardDetails    datatypes.JSON `gorm:"type:json" json:"award_details"`
	InterestAwarded float64        `gorm:"type:decimal(15,2)" json:"interest_awarded"`
	CostsAwarded    float64        `gorm:"type:decimal(15,2)" json:"costs_awarded"`

	// Enforcement
	AwardConfirmed    bool       `gorm:"default:false" json:"award_confirmed"`
	ConfirmationDate  *time.Time `json:"confirmation_date,omitempty"`
	ConfirmationCourt string     `gorm:"type:varchar(255)" json:"confirmation_court"`
	VacaturMotion     bool       `gorm:"default:false" json:"vacatur_motion"`
	AppealFiled       bool       `gorm:"default:false" json:"appeal_filed"`

	// Costs
	FilingFees         float64 `gorm:"type:decimal(15,2)" json:"filing_fees"`
	AdministrativeFees float64 `gorm:"type:decimal(15,2)" json:"administrative_fees"`
	AttorneyFees       float64 `gorm:"type:decimal(15,2)" json:"attorney_fees"`
	ExpertFees         float64 `gorm:"type:decimal(15,2)" json:"expert_fees"`
	TotalCosts         float64 `gorm:"type:decimal(15,2)" json:"total_costs"`

	// Settlement
	SettlementAttempted bool       `gorm:"default:false" json:"settlement_attempted"`
	SettlementReached   bool       `gorm:"default:false" json:"settlement_reached"`
	SettlementAmount    float64    `gorm:"type:decimal(15,2)" json:"settlement_amount"`
	SettlementDate      *time.Time `json:"settlement_date,omitempty"`

	// Status
	ArbitrationStatus string    `gorm:"type:varchar(50)" json:"arbitration_status"`
	Confidential      bool      `gorm:"default:true" json:"confidential"`
	Expedited         bool      `gorm:"default:false" json:"expedited"`
	CreatedAt         time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsHighRisk checks if litigation is high risk
func (cl *ClaimLitigation) IsHighRisk() bool {
	return cl.TotalExposure > 100000 ||
		cl.RiskLevel == "high" ||
		cl.PunitiveDamages > 0 ||
		cl.SuccessProbability < 30
}

// IsActive checks if litigation is active
func (cl *ClaimLitigation) IsActive() bool {
	return cl.LitigationStatus != "closed" &&
		cl.LitigationStatus != "settled" &&
		cl.JudgmentDate == nil
}

// GetTotalCosts calculates total litigation costs
func (cl *ClaimLitigation) GetTotalCosts() float64 {
	return cl.AttorneyFees + cl.CourtCosts +
		cl.ExpertWitnessFees + cl.DiscoveryCosts
}

// ShouldSettle evaluates if settlement is recommended
func (cl *ClaimLitigation) ShouldSettle() bool {
	// Recommend settlement if success probability is low
	// and settlement offer is reasonable
	if cl.SuccessProbability < 40 && cl.SettlementOffered {
		expectedLoss := cl.TotalExposure * (1 - cl.SuccessProbability/100)
		return cl.SettlementAmount < expectedLoss
	}
	return false
}

// HasDeadlineApproaching checks for upcoming deadlines
func (cl *ClaimLitigation) HasDeadlineApproaching() bool {
	if cl.NextDeadline == nil {
		return false
	}
	daysUntilDeadline := time.Until(*cl.NextDeadline).Hours() / 24
	return daysUntilDeadline <= 7
}

// IsFavorable checks if arbitration outcome is favorable
func (ca *ClaimArbitration) IsFavorable() bool {
	if ca.AwardAmount == 0 {
		return false
	}
	// Award is less than half of what was claimed against us
	return ca.AwardAmount < (ca.ClaimsAmount / 2)
}

// IsBinding checks if arbitration is binding
func (ca *ClaimArbitration) IsBinding() bool {
	return ca.AwardDate != nil &&
		ca.AwardType == "final" &&
		!ca.VacaturMotion
}

// GetNetExposure calculates net exposure after award
func (ca *ClaimArbitration) GetNetExposure() float64 {
	totalAward := ca.AwardAmount + ca.InterestAwarded + ca.CostsAwarded
	totalCosts := ca.FilingFees + ca.AdministrativeFees +
		ca.AttorneyFees + ca.ExpertFees + ca.ArbitratorFees
	return totalAward + totalCosts
}

// IsExpedited checks if arbitration is expedited
func (ca *ClaimArbitration) IsExpedited() bool {
	return ca.Expedited || ca.SingleArbitrator
}
