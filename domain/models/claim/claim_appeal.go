package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ClaimAppeal represents the appeal process for denied claims
type ClaimAppeal struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Appeal Identification
	AppealNumber   string `gorm:"type:varchar(100);unique" json:"appeal_number"`
	AppealLevel    int    `json:"appeal_level"`                        // 1st, 2nd, 3rd level
	AppealType     string `gorm:"type:varchar(50)" json:"appeal_type"` // internal, external, regulatory
	AppealCategory string `gorm:"type:varchar(50)" json:"appeal_category"`

	// Appeal Reason
	DenialReason    string  `gorm:"type:text" json:"denial_reason"`
	AppealReason    string  `gorm:"type:text" json:"appeal_reason"`
	DisputedAmount  float64 `gorm:"type:decimal(15,2)" json:"disputed_amount"`
	RequestedRelief string  `gorm:"type:text" json:"requested_relief"`

	// Timeline
	FiledDate        time.Time  `json:"filed_date"`
	ReceivedDate     time.Time  `json:"received_date"`
	AcknowledgedDate *time.Time `json:"acknowledged_date,omitempty"`
	ReviewStartDate  *time.Time `json:"review_start_date,omitempty"`
	HearingDate      *time.Time `json:"hearing_date,omitempty"`
	DecisionDate     *time.Time `json:"decision_date,omitempty"`
	DeadlineDate     time.Time  `json:"deadline_date"`

	// Appellant Information
	AppellantName      string         `gorm:"type:varchar(255)" json:"appellant_name"`
	AppellantType      string         `gorm:"type:varchar(50)" json:"appellant_type"` // customer, provider, attorney
	RepresentedBy      string         `gorm:"type:varchar(255)" json:"represented_by"`
	AttorneyInvolved   bool           `gorm:"default:false" json:"attorney_involved"`
	ContactInformation datatypes.JSON `gorm:"type:json" json:"contact_information"`

	// Review Process
	ReviewerID            *uuid.UUID     `gorm:"type:uuid" json:"reviewer_id,omitempty"`
	ReviewerNotes         string         `gorm:"type:text" json:"reviewer_notes"`
	ReviewPanel           datatypes.JSON `gorm:"type:json" json:"review_panel"` // []Reviewer
	ExternalReviewer      bool           `gorm:"default:false" json:"external_reviewer"`
	MedicalReviewRequired bool           `gorm:"default:false" json:"medical_review_required"`

	// Supporting Documentation
	DocumentsSubmitted datatypes.JSON `gorm:"type:json" json:"documents_submitted"`
	AdditionalEvidence datatypes.JSON `gorm:"type:json" json:"additional_evidence"`
	ExpertOpinions     datatypes.JSON `gorm:"type:json" json:"expert_opinions"`
	MedicalRecords     datatypes.JSON `gorm:"type:json" json:"medical_records"`

	// Decision
	DecisionStatus    string         `gorm:"type:varchar(50)" json:"decision_status"` // upheld, overturned, partial
	DecisionSummary   string         `gorm:"type:text" json:"decision_summary"`
	ApprovedAmount    float64        `gorm:"type:decimal(15,2)" json:"approved_amount"`
	DecisionRationale string         `gorm:"type:text" json:"decision_rationale"`
	ConditionsImposed datatypes.JSON `gorm:"type:json" json:"conditions_imposed"`

	// Regulatory Compliance
	RegulatoryBody      string     `gorm:"type:varchar(255)" json:"regulatory_body"`
	RegulatoryDeadline  *time.Time `json:"regulatory_deadline,omitempty"`
	ComplianceStatus    string     `gorm:"type:varchar(50)" json:"compliance_status"`
	ReportedToRegulator bool       `gorm:"default:false" json:"reported_to_regulator"`

	// Communication
	NotificationsSent datatypes.JSON `gorm:"type:json" json:"notifications_sent"`
	CorrespondenceLog datatypes.JSON `gorm:"type:json" json:"correspondence_log"`
	LastContactDate   *time.Time     `json:"last_contact_date,omitempty"`

	// Escalation
	EscalatedToExternal bool       `gorm:"default:false" json:"escalated_to_external"`
	EscalationDate      *time.Time `json:"escalation_date,omitempty"`
	EscalationReason    string     `gorm:"type:text" json:"escalation_reason"`

	// Status
	AppealStatus string    `gorm:"type:varchar(50)" json:"appeal_status"`
	Priority     string    `gorm:"type:varchar(20)" json:"priority"`
	SLAStatus    string    `gorm:"type:varchar(50)" json:"sla_status"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ClaimCommunication represents all communications related to a claim
type ClaimCommunication struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Communication Preferences
	PreferredChannel  string         `gorm:"type:varchar(50)" json:"preferred_channel"` // email, sms, phone, mail
	PreferredLanguage string         `gorm:"type:varchar(50)" json:"preferred_language"`
	PreferredTime     string         `gorm:"type:varchar(50)" json:"preferred_time"`
	DoNotContactFlags datatypes.JSON `gorm:"type:json" json:"do_not_contact_flags"`

	// Contact History
	TotalContacts    int            `json:"total_contacts"`
	InboundContacts  int            `json:"inbound_contacts"`
	OutboundContacts int            `json:"outbound_contacts"`
	LastContactDate  *time.Time     `json:"last_contact_date,omitempty"`
	LastContactType  string         `gorm:"type:varchar(50)" json:"last_contact_type"`
	ContactHistory   datatypes.JSON `gorm:"type:json" json:"contact_history"` // []ContactRecord

	// Email Communications
	EmailsSent        int     `json:"emails_sent"`
	EmailsReceived    int     `json:"emails_received"`
	EmailBounces      int     `json:"email_bounces"`
	EmailOpenRate     float64 `gorm:"type:decimal(5,2)" json:"email_open_rate"`
	EmailClickRate    float64 `gorm:"type:decimal(5,2)" json:"email_click_rate"`
	UnsubscribedEmail bool    `gorm:"default:false" json:"unsubscribed_email"`

	// SMS Communications
	SMSSent         int  `json:"sms_sent"`
	SMSReceived     int  `json:"sms_received"`
	SMSDelivered    int  `json:"sms_delivered"`
	SMSFailed       int  `json:"sms_failed"`
	UnsubscribedSMS bool `gorm:"default:false" json:"unsubscribed_sms"`

	// Phone Communications
	CallsMade           int  `json:"calls_made"`
	CallsReceived       int  `json:"calls_received"`
	AverageCallDuration int  `json:"average_call_duration_seconds"`
	VoicemailsLeft      int  `json:"voicemails_left"`
	CallBackRequested   bool `gorm:"default:false" json:"callback_requested"`

	// Mail Communications
	LettersSent       int  `json:"letters_sent"`
	CertifiedMailSent int  `json:"certified_mail_sent"`
	MailReturned      int  `json:"mail_returned"`
	AddressVerified   bool `gorm:"default:true" json:"address_verified"`

	// Automated Communications
	AutoNotifications    datatypes.JSON `gorm:"type:json" json:"auto_notifications"` // []Notification
	TriggeredEmails      datatypes.JSON `gorm:"type:json" json:"triggered_emails"`
	ScheduledComms       datatypes.JSON `gorm:"type:json" json:"scheduled_comms"`
	NextScheduledContact *time.Time     `json:"next_scheduled_contact,omitempty"`

	// Templates Used
	TemplatesUsed       datatypes.JSON `gorm:"type:json" json:"templates_used"` // []Template
	CustomMessages      datatypes.JSON `gorm:"type:json" json:"custom_messages"`
	PersonalizationData datatypes.JSON `gorm:"type:json" json:"personalization_data"`

	// Response Tracking
	ResponseRequired bool       `gorm:"default:false" json:"response_required"`
	ResponseDeadline *time.Time `json:"response_deadline,omitempty"`
	ResponseReceived bool       `gorm:"default:false" json:"response_received"`
	ResponseDate     *time.Time `json:"response_date,omitempty"`
	FollowUpRequired bool       `gorm:"default:false" json:"follow_up_required"`

	// Satisfaction
	SatisfactionScore  float64        `gorm:"type:decimal(5,2)" json:"satisfaction_score"`
	NPSScore           int            `json:"nps_score"`
	FeedbackReceived   datatypes.JSON `gorm:"type:json" json:"feedback_received"`
	ComplaintsReceived int            `json:"complaints_received"`

	// Compliance
	ConsentObtained bool           `gorm:"default:true" json:"consent_obtained"`
	ConsentDate     *time.Time     `json:"consent_date,omitempty"`
	OptOutRequests  datatypes.JSON `gorm:"type:json" json:"opt_out_requests"`
	ComplianceFlags datatypes.JSON `gorm:"type:json" json:"compliance_flags"`

	// Status
	CommunicationStatus string    `gorm:"type:varchar(50)" json:"communication_status"`
	LastUpdateDate      time.Time `json:"last_update_date"`
	CreatedAt           time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ClaimAnalytics represents analytics and reporting for claims
type ClaimAnalytics struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Performance Metrics
	ProcessingTime         int  `json:"processing_time_hours"`
	CycleTime              int  `json:"cycle_time_days"`
	TouchTime              int  `json:"touch_time_minutes"`
	WaitTime               int  `json:"wait_time_hours"`
	FirstContactResolution bool `gorm:"default:false" json:"first_contact_resolution"`

	// Cost Analysis
	TotalCost         float64 `gorm:"type:decimal(15,2)" json:"total_cost"`
	IndemnityCost     float64 `gorm:"type:decimal(15,2)" json:"indemnity_cost"`
	ExpenseCost       float64 `gorm:"type:decimal(15,2)" json:"expense_cost"`
	LegalCost         float64 `gorm:"type:decimal(15,2)" json:"legal_cost"`
	InvestigationCost float64 `gorm:"type:decimal(15,2)" json:"investigation_cost"`
	CostPercentile    float64 `gorm:"type:decimal(5,2)" json:"cost_percentile"`

	// Predictive Scores
	SettlementPrediction float64 `gorm:"type:decimal(15,2)" json:"settlement_prediction"`
	DurationPrediction   int     `json:"duration_prediction_days"`
	LitigationRisk       float64 `gorm:"type:decimal(5,2)" json:"litigation_risk"`
	CustomerChurnRisk    float64 `gorm:"type:decimal(5,2)" json:"customer_churn_risk"`
	FraudProbability     float64 `gorm:"type:decimal(5,2)" json:"fraud_probability"`

	// Benchmarking
	IndustryBenchmark  datatypes.JSON `gorm:"type:json" json:"industry_benchmark"`
	InternalBenchmark  datatypes.JSON `gorm:"type:json" json:"internal_benchmark"`
	PerformanceRanking int            `json:"performance_ranking"`
	EfficiencyScore    float64        `gorm:"type:decimal(5,2)" json:"efficiency_score"`

	// Trends
	TrendAnalysis     datatypes.JSON `gorm:"type:json" json:"trend_analysis"`
	SeasonalFactors   datatypes.JSON `gorm:"type:json" json:"seasonal_factors"`
	PatternIdentified datatypes.JSON `gorm:"type:json" json:"pattern_identified"`

	// Quality Metrics
	QualityScore float64 `gorm:"type:decimal(5,2)" json:"quality_score"`
	AccuracyRate float64 `gorm:"type:decimal(5,2)" json:"accuracy_rate"`
	ErrorRate    float64 `gorm:"type:decimal(5,2)" json:"error_rate"`
	ReworkRate   float64 `gorm:"type:decimal(5,2)" json:"rework_rate"`

	// Customer Metrics
	CustomerEffort       float64 `gorm:"type:decimal(5,2)" json:"customer_effort"`
	CustomerSatisfaction float64 `gorm:"type:decimal(5,2)" json:"customer_satisfaction"`
	NetPromoterScore     int     `json:"net_promoter_score"`
	CustomerContacts     int     `json:"customer_contacts"`

	// SLA Compliance
	SLACompliant   bool    `gorm:"default:true" json:"sla_compliant"`
	SLABreaches    int     `json:"sla_breaches"`
	SLAPerformance float64 `gorm:"type:decimal(5,2)" json:"sla_performance"`

	// Status
	AnalyticsStatus string    `gorm:"type:varchar(50)" json:"analytics_status"`
	LastCalculated  time.Time `json:"last_calculated"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsOpen checks if appeal is still open
func (ca *ClaimAppeal) IsOpen() bool {
	return ca.AppealStatus != "closed" &&
		ca.DecisionDate == nil
}

// IsOverdue checks if appeal is overdue
func (ca *ClaimAppeal) IsOverdue() bool {
	return time.Now().After(ca.DeadlineDate) && ca.IsOpen()
}

// IsSuccessful checks if appeal was successful
func (ca *ClaimAppeal) IsSuccessful() bool {
	return ca.DecisionStatus == "overturned" ||
		(ca.DecisionStatus == "partial" && ca.ApprovedAmount > 0)
}

// NeedsResponse checks if communication needs response
func (cc *ClaimCommunication) NeedsResponse() bool {
	if !cc.ResponseRequired {
		return false
	}
	if cc.ResponseReceived {
		return false
	}
	if cc.ResponseDeadline != nil {
		return time.Now().Before(*cc.ResponseDeadline)
	}
	return true
}

// IsHighEngagement checks for high customer engagement
func (cc *ClaimCommunication) IsHighEngagement() bool {
	return cc.TotalContacts > 10 ||
		cc.InboundContacts > 5 ||
		cc.ComplaintsReceived > 0
}

// GetEngagementScore calculates engagement score
func (cc *ClaimCommunication) GetEngagementScore() float64 {
	openRate := cc.EmailOpenRate
	responseRate := float64(0)
	if cc.ResponseRequired && cc.ResponseReceived {
		responseRate = 100
	}
	return (openRate + responseRate + cc.SatisfactionScore) / 3
}

// IsHighCost checks if claim is high cost
func (ca *ClaimAnalytics) IsHighCost() bool {
	return ca.TotalCost > 10000 || ca.CostPercentile > 90
}

// IsEfficient checks if claim processing was efficient
func (ca *ClaimAnalytics) IsEfficient() bool {
	return ca.EfficiencyScore > 80 &&
		ca.SLACompliant &&
		ca.FirstContactResolution
}

// GetPerformanceScore calculates overall performance
func (ca *ClaimAnalytics) GetPerformanceScore() float64 {
	return (ca.EfficiencyScore + ca.QualityScore +
		ca.CustomerSatisfaction + float64(ca.SLAPerformance)) / 4
}
