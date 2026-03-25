package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PolicyPredictiveAnalytics represents AI/ML-driven analytics for a policy
type PolicyPredictiveAnalytics struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID     uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	AnalysisDate time.Time `gorm:"type:timestamp;not null" json:"analysis_date"`

	// Churn Prediction
	ChurnPredictionScore     float64        `gorm:"type:decimal(5,4)" json:"churn_prediction_score"`
	ChurnProbability         float64        `gorm:"type:decimal(5,4)" json:"churn_probability"`
	ChurnRiskLevel           string         `gorm:"type:varchar(50)" json:"churn_risk_level"` // low, medium, high, critical
	ChurnDrivers             datatypes.JSON `gorm:"type:json" json:"churn_drivers"`           // []ChurnFactor
	RetentionRecommendations datatypes.JSON `gorm:"type:json" json:"retention_recommendations"`

	// Claim Prediction
	ClaimPredictionScore float64        `gorm:"type:decimal(5,4)" json:"claim_prediction_score"`
	ClaimProbability     float64        `gorm:"type:decimal(5,4)" json:"claim_probability"`
	ExpectedClaimAmount  Money          `gorm:"embedded;embeddedPrefix:expected_claim_" json:"expected_claim_amount"`
	ClaimTimelineDays    int            `gorm:"type:int" json:"claim_timeline_days"`
	ClaimTypesPredicted  datatypes.JSON `gorm:"type:json" json:"claim_types_predicted"`

	// Fraud Detection
	FraudProbability   float64        `gorm:"type:decimal(5,4)" json:"fraud_probability"`
	FraudRiskScore     float64        `gorm:"type:decimal(10,2)" json:"fraud_risk_score"`
	FraudIndicators    datatypes.JSON `gorm:"type:json" json:"fraud_indicators"` // []FraudIndicator
	AnomalyScore       float64        `gorm:"type:decimal(10,2)" json:"anomaly_score"`
	SuspiciousPatterns datatypes.JSON `gorm:"type:json" json:"suspicious_patterns"`

	// Lifetime Value
	LifetimeValuePrediction Money          `gorm:"embedded;embeddedPrefix:ltv_prediction_" json:"lifetime_value_prediction"`
	LTVConfidenceInterval   datatypes.JSON `gorm:"type:json" json:"ltv_confidence_interval"`
	ExpectedTenureMonths    int            `gorm:"type:int" json:"expected_tenure_months"`
	CrossSellPotential      float64        `gorm:"type:decimal(5,2)" json:"cross_sell_potential"`
	UpsellPotential         float64        `gorm:"type:decimal(5,2)" json:"upsell_potential"`

	// Next Best Action
	NextBestAction      string         `gorm:"type:varchar(255)" json:"next_best_action"`
	ActionPriority      string         `gorm:"type:varchar(50)" json:"action_priority"` // urgent, high, medium, low
	ActionReason        string         `gorm:"type:text" json:"action_reason"`
	RecommendedActions  datatypes.JSON `gorm:"type:json" json:"recommended_actions"` // []RecommendedAction
	ActionEffectiveness float64        `gorm:"type:decimal(5,2)" json:"action_effectiveness"`

	// Risk Scoring (ML-based)
	RiskScoreML    float64        `gorm:"type:decimal(10,2)" json:"risk_score_ml"`
	RiskFactors    datatypes.JSON `gorm:"type:json" json:"risk_factors"`      // []RiskFactor
	RiskTrend      string         `gorm:"type:varchar(50)" json:"risk_trend"` // increasing, stable, decreasing
	RiskPercentile float64        `gorm:"type:decimal(5,2)" json:"risk_percentile"`
	PeerComparison float64        `gorm:"type:decimal(5,2)" json:"peer_comparison"`

	// Behavioral Segmentation
	BehavioralSegment      string         `gorm:"type:varchar(100)" json:"behavioral_segment"`
	SegmentCharacteristics datatypes.JSON `gorm:"type:json" json:"segment_characteristics"`
	SegmentMigrationProb   float64        `gorm:"type:decimal(5,4)" json:"segment_migration_probability"`
	TargetSegment          string         `gorm:"type:varchar(100)" json:"target_segment"`

	// Propensity Scores
	PropensityScores   datatypes.JSON `gorm:"type:json" json:"propensity_scores"` // map[string]float64
	PaymentDefaultProb float64        `gorm:"type:decimal(5,4)" json:"payment_default_probability"`
	RenewalProbability float64        `gorm:"type:decimal(5,4)" json:"renewal_probability"`
	LapseProbability   float64        `gorm:"type:decimal(5,4)" json:"lapse_probability"`

	// Model Information
	ModelsUsed       datatypes.JSON `gorm:"type:json" json:"models_used"` // []ModelInfo
	ModelVersion     string         `gorm:"type:varchar(50)" json:"model_version"`
	ModelConfidence  float64        `gorm:"type:decimal(5,4)" json:"model_confidence"`
	DataQualityScore float64        `gorm:"type:decimal(5,2)" json:"data_quality_score"`

	// Recommendations
	PricingOptimization     datatypes.JSON `gorm:"type:json" json:"pricing_optimization"`
	CoverageRecommendations datatypes.JSON `gorm:"type:json" json:"coverage_recommendations"`
	ServiceRecommendations  datatypes.JSON `gorm:"type:json" json:"service_recommendations"`

	// Status
	AnalysisStatus   string    `gorm:"type:varchar(50)" json:"analysis_status"`
	LastModelUpdate  time.Time `gorm:"type:timestamp" json:"last_model_update"`
	NextAnalysisDate time.Time `gorm:"type:timestamp" json:"next_analysis_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// PolicyCustomerJourney represents customer experience tracking
type PolicyCustomerJourney struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Journey Tracking
	JourneyStage     string         `gorm:"type:varchar(100)" json:"journey_stage"`
	JourneyStartDate time.Time      `gorm:"type:timestamp" json:"journey_start_date"`
	CurrentMilestone string         `gorm:"type:varchar(255)" json:"current_milestone"`
	JourneyMap       datatypes.JSON `gorm:"type:json" json:"journey_map"` // []JourneyStep

	// TouchPoints
	TouchPoints         datatypes.JSON `gorm:"type:json" json:"touch_points"` // []TouchPoint
	TotalTouchPoints    int            `gorm:"type:int" json:"total_touch_points"`
	LastTouchPoint      string         `gorm:"type:varchar(255)" json:"last_touch_point"`
	LastTouchPointDate  time.Time      `gorm:"type:timestamp" json:"last_touch_point_date"`
	PreferredTouchPoint string         `gorm:"type:varchar(100)" json:"preferred_touch_point"`

	// Interactions
	InteractionHistory   datatypes.JSON `gorm:"type:json" json:"interaction_history"` // []Interaction
	TotalInteractions    int            `gorm:"type:int" json:"total_interactions"`
	PositiveInteractions int            `gorm:"type:int" json:"positive_interactions"`
	NegativeInteractions int            `gorm:"type:int" json:"negative_interactions"`
	LastInteractionDate  time.Time      `gorm:"type:timestamp" json:"last_interaction_date"`

	// Channel Preferences
	PreferredChannels datatypes.JSON `gorm:"type:json" json:"preferred_channels"`  // []string
	ChannelUsageStats datatypes.JSON `gorm:"type:json" json:"channel_usage_stats"` // map[string]ChannelStats
	OmniChannelScore  float64        `gorm:"type:decimal(5,2)" json:"omni_channel_score"`
	DigitalEngagement float64        `gorm:"type:decimal(5,2)" json:"digital_engagement_score"`

	// Customer Effort
	CustomerEffortScore    float64        `gorm:"type:decimal(5,2)" json:"customer_effort_score"`
	FrictionPoints         datatypes.JSON `gorm:"type:json" json:"friction_points"` // []FrictionPoint
	ResolutionTime         int            `gorm:"type:int" json:"avg_resolution_time_hours"`
	FirstContactResolution float64        `gorm:"type:decimal(5,2)" json:"first_contact_resolution_rate"`

	// Satisfaction & NPS
	NPSScore            int            `gorm:"type:int" json:"nps_score"`
	NPSHistory          datatypes.JSON `gorm:"type:json" json:"nps_history"` // []NPSRecord
	CSATScore           float64        `gorm:"type:decimal(5,2)" json:"csat_score"`
	CSATHistory         datatypes.JSON `gorm:"type:json" json:"csat_history"` // []CSATRecord
	SatisfactionDrivers datatypes.JSON `gorm:"type:json" json:"satisfaction_drivers"`

	// Complaints & Feedback
	ComplaintHistory        datatypes.JSON `gorm:"type:json" json:"complaint_history"` // []Complaint
	ComplaintCount          int            `gorm:"type:int" json:"complaint_count"`
	ResolvedComplaints      int            `gorm:"type:int" json:"resolved_complaints"`
	ComplaintResolutionTime int            `gorm:"type:int" json:"avg_complaint_resolution_days"`
	FeedbackHistory         datatypes.JSON `gorm:"type:json" json:"feedback_history"` // []Feedback

	// Surveys
	SurveyResponses    datatypes.JSON `gorm:"type:json" json:"survey_responses"` // []SurveyResponse
	SurveyResponseRate float64        `gorm:"type:decimal(5,2)" json:"survey_response_rate"`
	LastSurveyDate     *time.Time     `gorm:"type:timestamp" json:"last_survey_date,omitempty"`
	NextSurveyDate     *time.Time     `gorm:"type:timestamp" json:"next_survey_date,omitempty"`

	// Sentiment Analysis
	SentimentScore   float64        `gorm:"type:decimal(5,2)" json:"sentiment_score"` // -100 to +100
	SentimentTrend   string         `gorm:"type:varchar(50)" json:"sentiment_trend"`  // improving, stable, declining
	SentimentHistory datatypes.JSON `gorm:"type:json" json:"sentiment_history"`       // []SentimentRecord
	EmotionalDrivers datatypes.JSON `gorm:"type:json" json:"emotional_drivers"`

	// Engagement Metrics
	EngagementScore      float64   `gorm:"type:decimal(5,2)" json:"engagement_score"`
	LastEngagementDate   time.Time `gorm:"type:timestamp" json:"last_engagement_date"`
	EngagementFrequency  float64   `gorm:"type:decimal(10,2)" json:"engagement_frequency_per_month"`
	ActiveEngagementDays int       `gorm:"type:int" json:"active_engagement_days"`

	// Personalization
	PersonalizationScore float64        `gorm:"type:decimal(5,2)" json:"personalization_score"`
	PersonalPreferences  datatypes.JSON `gorm:"type:json" json:"personal_preferences"`
	CommunicationPrefs   datatypes.JSON `gorm:"type:json" json:"communication_preferences"`
	ContentPreferences   datatypes.JSON `gorm:"type:json" json:"content_preferences"`

	// Loyalty & Advocacy
	LoyaltyScore  float64 `gorm:"type:decimal(5,2)" json:"loyalty_score"`
	AdvocacyScore float64 `gorm:"type:decimal(5,2)" json:"advocacy_score"`
	ReferralCount int     `gorm:"type:int" json:"referral_count"`
	SocialShares  int     `gorm:"type:int" json:"social_shares"`
	ReviewsPosted int     `gorm:"type:int" json:"reviews_posted"`

	// Status
	JourneyStatus  string    `gorm:"type:varchar(50)" json:"journey_status"`
	RiskStatus     string    `gorm:"type:varchar(50)" json:"risk_status"` // at_risk, stable, champion
	LastUpdateDate time.Time `gorm:"type:timestamp" json:"last_update_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// =====================================
// METHODS
// =====================================

// IsHighChurnRisk checks if customer has high churn risk
func (ppa *PolicyPredictiveAnalytics) IsHighChurnRisk() bool {
	return ppa.ChurnProbability > 0.7 || ppa.ChurnRiskLevel == "high" ||
		ppa.ChurnRiskLevel == "critical"
}

// IsHighFraudRisk checks if there's high fraud risk
func (ppa *PolicyPredictiveAnalytics) IsHighFraudRisk() bool {
	return ppa.FraudProbability > 0.6 || ppa.FraudRiskScore > 80
}

// ShouldTakeAction checks if immediate action is needed
func (ppa *PolicyPredictiveAnalytics) ShouldTakeAction() bool {
	return ppa.ActionPriority == "urgent" || ppa.ActionPriority == "high"
}

// GetRiskLevel returns the overall risk level
func (ppa *PolicyPredictiveAnalytics) GetRiskLevel() string {
	if ppa.RiskScoreML > 80 {
		return "critical"
	} else if ppa.RiskScoreML > 60 {
		return "high"
	} else if ppa.RiskScoreML > 40 {
		return "medium"
	}
	return "low"
}

// IsHighValueCustomer checks if customer is high value
func (ppa *PolicyPredictiveAnalytics) IsHighValueCustomer() bool {
	return ppa.LifetimeValuePrediction.Amount > 10000 ||
		ppa.CrossSellPotential > 70 ||
		ppa.UpsellPotential > 70
}

// IsPromoter checks if customer is a promoter (NPS >= 9)
func (pcj *PolicyCustomerJourney) IsPromoter() bool {
	return pcj.NPSScore >= 9
}

// IsDetractor checks if customer is a detractor (NPS <= 6)
func (pcj *PolicyCustomerJourney) IsDetractor() bool {
	return pcj.NPSScore <= 6
}

// IsAtRisk checks if customer is at risk of churning
func (pcj *PolicyCustomerJourney) IsAtRisk() bool {
	return pcj.RiskStatus == "at_risk" ||
		pcj.SentimentScore < -20 ||
		pcj.ComplaintCount > 2 ||
		pcj.NPSScore <= 6
}

// GetSentimentCategory returns sentiment category
func (pcj *PolicyCustomerJourney) GetSentimentCategory() string {
	if pcj.SentimentScore >= 50 {
		return "very_positive"
	} else if pcj.SentimentScore >= 20 {
		return "positive"
	} else if pcj.SentimentScore >= -20 {
		return "neutral"
	} else if pcj.SentimentScore >= -50 {
		return "negative"
	}
	return "very_negative"
}

// IsHighlyEngaged checks if customer is highly engaged
func (pcj *PolicyCustomerJourney) IsHighlyEngaged() bool {
	return pcj.EngagementScore > 70 ||
		pcj.EngagementFrequency > 5 ||
		pcj.DigitalEngagement > 80
}

// HasUnresolvedComplaints checks for unresolved complaints
func (pcj *PolicyCustomerJourney) HasUnresolvedComplaints() bool {
	return pcj.ComplaintCount > pcj.ResolvedComplaints
}
