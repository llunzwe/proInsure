package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ClaimInvestigationDetail represents detailed investigation management
type ClaimInvestigationDetail struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Investigation Team
	LeadInvestigatorID   uuid.UUID      `gorm:"type:uuid" json:"lead_investigator_id"`
	TeamMembers          datatypes.JSON `gorm:"type:json" json:"team_members"` // []uuid.UUID
	ExternalInvestigator bool           `gorm:"default:false" json:"external_investigator"`
	VendorID             *uuid.UUID     `gorm:"type:uuid" json:"vendor_id,omitempty"`

	// Investigation Details
	CaseNumber         string         `gorm:"type:varchar(100);unique" json:"case_number"`
	Priority           string         `gorm:"type:varchar(20)" json:"priority"`
	Complexity         string         `gorm:"type:varchar(20)" json:"complexity"` // low, medium, high, critical
	InvestigationType  string         `gorm:"type:varchar(50)" json:"investigation_type"`
	InvestigationScope datatypes.JSON `gorm:"type:json" json:"investigation_scope"`
	Methodology        datatypes.JSON `gorm:"type:json" json:"methodology"`

	// Timeline
	RequestedDate       time.Time  `json:"requested_date"`
	StartedDate         *time.Time `json:"started_date,omitempty"`
	EstimatedCompletion time.Time  `json:"estimated_completion"`
	CompletedDate       *time.Time `json:"completed_date,omitempty"`
	DaysElapsed         int        `json:"days_elapsed"`
	Extensions          int        `json:"extensions"`

	// Evidence Collection
	EvidenceInventory datatypes.JSON `gorm:"type:json" json:"evidence_inventory"`
	ChainOfCustody    datatypes.JSON `gorm:"type:json" json:"chain_of_custody"`
	ForensicReports   datatypes.JSON `gorm:"type:json" json:"forensic_reports"`
	LabResults        datatypes.JSON `gorm:"type:json" json:"lab_results"`
	ExpertTestimonies datatypes.JSON `gorm:"type:json" json:"expert_testimonies"`

	// Interviews & Statements
	InterviewsConducted int            `json:"interviews_conducted"`
	InterviewRecords    datatypes.JSON `gorm:"type:json" json:"interview_records"`
	StatementsTaken     int            `json:"statements_taken"`
	StatementRecords    datatypes.JSON `gorm:"type:json" json:"statement_records"`
	Contradictions      datatypes.JSON `gorm:"type:json" json:"contradictions"`

	// Site Investigation
	SiteVisitRequired    bool           `gorm:"default:false" json:"site_visit_required"`
	SiteVisitDate        *time.Time     `json:"site_visit_date,omitempty"`
	SiteInspectionReport datatypes.JSON `gorm:"type:json" json:"site_inspection_report"`
	PhotoDocumentation   datatypes.JSON `gorm:"type:json" json:"photo_documentation"`
	VideoDocumentation   datatypes.JSON `gorm:"type:json" json:"video_documentation"`

	// Surveillance
	SurveillanceRequired bool           `gorm:"default:false" json:"surveillance_required"`
	SurveillancePeriod   datatypes.JSON `gorm:"type:json" json:"surveillance_period"` // DateRange
	SurveillanceReports  datatypes.JSON `gorm:"type:json" json:"surveillance_reports"`
	SurveillanceFindings datatypes.JSON `gorm:"type:json" json:"surveillance_findings"`

	// Background Checks
	BackgroundChecksDone bool           `gorm:"default:false" json:"background_checks_done"`
	CriminalHistory      datatypes.JSON `gorm:"type:json" json:"criminal_history"`
	ClaimHistory         datatypes.JSON `gorm:"type:json" json:"claim_history"`
	FinancialStatus      datatypes.JSON `gorm:"type:json" json:"financial_status"`
	SocialMediaAnalysis  datatypes.JSON `gorm:"type:json" json:"social_media_analysis"`

	// Findings & Recommendations
	KeyFindings       datatypes.JSON `gorm:"type:json" json:"key_findings"`
	RedFlags          datatypes.JSON `gorm:"type:json" json:"red_flags"`
	GreenFlags        datatypes.JSON `gorm:"type:json" json:"green_flags"`
	Recommendations   datatypes.JSON `gorm:"type:json" json:"recommendations"`
	ConclusionSummary string         `gorm:"type:text" json:"conclusion_summary"`

	// Outcome
	InvestigationResult    string  `gorm:"type:varchar(50)" json:"investigation_result"` // legitimate, fraudulent, inconclusive
	ConfidenceLevel        float64 `gorm:"type:decimal(5,2)" json:"confidence_level"`
	FraudAmount            float64 `gorm:"type:decimal(15,2)" json:"fraud_amount"`
	RecoveryPotential      float64 `gorm:"type:decimal(15,2)" json:"recovery_potential"`
	ProsecutionRecommended bool    `gorm:"default:false" json:"prosecution_recommended"`

	// Quality & Review
	QualityReviewDone bool       `gorm:"default:false" json:"quality_review_done"`
	ReviewedBy        *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewDate        *time.Time `json:"review_date,omitempty"`
	ReviewComments    string     `gorm:"type:text" json:"review_comments"`

	// Cost Management
	InvestigationCost float64 `gorm:"type:decimal(15,2)" json:"investigation_cost"`
	BudgetAllocated   float64 `gorm:"type:decimal(15,2)" json:"budget_allocated"`
	CostJustification string  `gorm:"type:text" json:"cost_justification"`
	ROI               float64 `gorm:"type:decimal(10,2)" json:"roi"`

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

// ClaimFraudDetection represents fraud detection and scoring
type ClaimFraudDetection struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Fraud Scoring
	OverallScore    float64 `gorm:"type:decimal(5,4)" json:"overall_score"`
	MLScore         float64 `gorm:"type:decimal(5,4)" json:"ml_score"`
	RuleBasedScore  float64 `gorm:"type:decimal(5,4)" json:"rule_based_score"`
	BehavioralScore float64 `gorm:"type:decimal(5,4)" json:"behavioral_score"`
	NetworkScore    float64 `gorm:"type:decimal(5,4)" json:"network_score"`
	ScoreConfidence float64 `gorm:"type:decimal(5,2)" json:"score_confidence"`

	// Risk Indicators
	RiskLevel         string         `gorm:"type:varchar(20)" json:"risk_level"`
	RiskFactors       datatypes.JSON `gorm:"type:json" json:"risk_factors"`
	AnomaliesDetected datatypes.JSON `gorm:"type:json" json:"anomalies_detected"`
	PatternMatches    datatypes.JSON `gorm:"type:json" json:"pattern_matches"`

	// Fraud Patterns
	KnownPatterns       datatypes.JSON `gorm:"type:json" json:"known_patterns"`
	SuspiciousPatterns  datatypes.JSON `gorm:"type:json" json:"suspicious_patterns"`
	BehavioralAnomalies datatypes.JSON `gorm:"type:json" json:"behavioral_anomalies"`
	VelocityChecks      datatypes.JSON `gorm:"type:json" json:"velocity_checks"`

	// ML Model Details
	ModelVersion          string         `gorm:"type:varchar(50)" json:"model_version"`
	ModelFeatures         datatypes.JSON `gorm:"type:json" json:"model_features"`
	FeatureImportance     datatypes.JSON `gorm:"type:json" json:"feature_importance"`
	PredictionExplanation datatypes.JSON `gorm:"type:json" json:"prediction_explanation"`

	// Network Analysis
	NetworkConnections datatypes.JSON `gorm:"type:json" json:"network_connections"`
	LinkedClaims       datatypes.JSON `gorm:"type:json" json:"linked_claims"`
	SharedAttributes   datatypes.JSON `gorm:"type:json" json:"shared_attributes"`
	ClusterID          string         `gorm:"type:varchar(100)" json:"cluster_id"`

	// Historical Analysis
	PreviousClaims     int     `json:"previous_claims"`
	ClaimFrequency     float64 `gorm:"type:decimal(10,2)" json:"claim_frequency"`
	AverageClaimAmount float64 `gorm:"type:decimal(15,2)" json:"average_claim_amount"`
	TimesinceLastClaim int     `json:"time_since_last_claim_days"`

	// Alert Management
	AlertsGenerated int            `json:"alerts_generated"`
	AlertDetails    datatypes.JSON `gorm:"type:json" json:"alert_details"`
	AlertPriority   string         `gorm:"type:varchar(20)" json:"alert_priority"`
	AlertAssignee   *uuid.UUID     `gorm:"type:uuid" json:"alert_assignee,omitempty"`
	AlertResolved   bool           `gorm:"default:false" json:"alert_resolved"`

	// Investigation Trigger
	InvestigationTriggered bool    `gorm:"default:false" json:"investigation_triggered"`
	TriggerReason          string  `gorm:"type:text" json:"trigger_reason"`
	TriggerThreshold       float64 `gorm:"type:decimal(5,2)" json:"trigger_threshold"`
	AutoBlocked            bool    `gorm:"default:false" json:"auto_blocked"`

	// Review & Feedback
	ManualReviewRequired bool       `gorm:"default:false" json:"manual_review_required"`
	ReviewedBy           *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewOutcome        string     `gorm:"type:varchar(50)" json:"review_outcome"`
	FeedbackProvided     bool       `gorm:"default:false" json:"feedback_provided"`
	ModelRetrained       bool       `gorm:"default:false" json:"model_retrained"`

	// Status
	DetectionStatus string    `gorm:"type:varchar(50)" json:"detection_status"`
	LastScoredAt    time.Time `json:"last_scored_at"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsHighRisk checks if investigation indicates high risk
func (ci *ClaimInvestigationDetail) IsHighRisk() bool {
	return ci.InvestigationResult == "fraudulent" ||
		ci.ProsecutionRecommended ||
		ci.FraudAmount > 10000
}

// IsComplete checks if investigation is complete
func (ci *ClaimInvestigationDetail) IsComplete() bool {
	return ci.CompletedDate != nil &&
		ci.Status == "completed" &&
		ci.QualityReviewDone
}

// GetROI calculates return on investigation investment
func (ci *ClaimInvestigationDetail) GetROI() float64 {
	if ci.InvestigationCost == 0 {
		return 0
	}
	return (ci.FraudAmount - ci.InvestigationCost) / ci.InvestigationCost * 100
}

// IsFraudulent checks if claim is fraudulent based on scoring
func (cf *ClaimFraudDetection) IsFraudulent() bool {
	return cf.OverallScore > 0.7 ||
		cf.RiskLevel == "critical" ||
		cf.AutoBlocked
}

// NeedsManualReview checks if manual review is required
func (cf *ClaimFraudDetection) NeedsManualReview() bool {
	return cf.ManualReviewRequired ||
		(cf.OverallScore > 0.5 && cf.OverallScore < 0.8) ||
		cf.ScoreConfidence < 0.6
}

// GetRiskCategory returns the risk category based on score
func (cf *ClaimFraudDetection) GetRiskCategory() string {
	if cf.OverallScore < 0.3 {
		return "low"
	} else if cf.OverallScore < 0.5 {
		return "medium"
	} else if cf.OverallScore < 0.7 {
		return "high"
	}
	return "critical"
}
