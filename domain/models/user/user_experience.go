package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserJourneyMapping tracks complete user journey through the system
type UserJourneyMapping struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Journey Overview
	JourneyStage       string                   `gorm:"type:varchar(50)" json:"journey_stage"`
	CurrentStep        string                   `gorm:"type:varchar(100)" json:"current_step"`
	JourneyProgress    float64                  `gorm:"default:0" json:"journey_progress"`
	StageHistory       []map[string]interface{} `gorm:"type:json" json:"stage_history"`
	TimeInCurrentStage int                      `json:"time_in_current_stage_hours"`

	// Touchpoints
	TotalTouchpoints        int                      `gorm:"default:0" json:"total_touchpoints"`
	TouchpointHistory       []map[string]interface{} `gorm:"type:json" json:"touchpoint_history"`
	PreferredTouchpoints    []string                 `gorm:"type:json" json:"preferred_touchpoints"`
	TouchpointEffectiveness map[string]float64       `gorm:"type:json" json:"touchpoint_effectiveness"`

	// Pain Points
	IdentifiedPainPoints []string                 `gorm:"type:json" json:"identified_pain_points"`
	FrictionPoints       []map[string]interface{} `gorm:"type:json" json:"friction_points"`
	AbandonmentPoints    []string                 `gorm:"type:json" json:"abandonment_points"`
	ResolutionAttempts   map[string]int           `gorm:"type:json" json:"resolution_attempts"`

	// Moments of Truth
	CriticalMoments []map[string]interface{} `gorm:"type:json" json:"critical_moments"`
	DelightMoments  []map[string]interface{} `gorm:"type:json" json:"delight_moments"`
	FailureMoments  []map[string]interface{} `gorm:"type:json" json:"failure_moments"`
	RecoveryActions []map[string]interface{} `gorm:"type:json" json:"recovery_actions"`

	// Experience Metrics
	JourneyScore   float64 `gorm:"default:0" json:"journey_score"`
	FrictionScore  float64 `gorm:"default:0" json:"friction_score"`
	EffortScore    float64 `gorm:"default:0" json:"effort_score"`
	CompletionRate float64 `gorm:"default:0" json:"completion_rate"`
	TimeToValue    int     `json:"time_to_value_hours"`

	// Personalization
	PersonalizationLevel float64                  `gorm:"default:0" json:"personalization_level"`
	PersonalizedElements []string                 `gorm:"type:json" json:"personalized_elements"`
	Recommendations      []map[string]interface{} `gorm:"type:json" json:"recommendations"`
	NextBestActions      []string                 `gorm:"type:json" json:"next_best_actions"`
}

// UserOnboarding manages user onboarding process
type UserOnboarding struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Onboarding Status
	OnboardingStatus     string     `gorm:"type:varchar(20)" json:"onboarding_status"`
	OnboardingStartDate  time.Time  `json:"onboarding_start_date"`
	OnboardingEndDate    *time.Time `json:"onboarding_end_date"`
	CompletionPercentage float64    `gorm:"default:0" json:"completion_percentage"`

	// Steps & Milestones
	CompletedSteps   []string                 `gorm:"type:json" json:"completed_steps"`
	PendingSteps     []string                 `gorm:"type:json" json:"pending_steps"`
	SkippedSteps     []string                 `gorm:"type:json" json:"skipped_steps"`
	Milestones       []map[string]interface{} `gorm:"type:json" json:"milestones"`
	CurrentMilestone string                   `gorm:"type:varchar(100)" json:"current_milestone"`

	// Guidance & Support
	TutorialsViewed      []string `gorm:"type:json" json:"tutorials_viewed"`
	HelpArticlesAccessed []string `gorm:"type:json" json:"help_articles_accessed"`
	SupportContactCount  int      `gorm:"default:0" json:"support_contact_count"`
	GuidedTourCompleted  bool     `gorm:"default:false" json:"guided_tour_completed"`

	// Time Metrics
	TimeToFirstAction   int `json:"time_to_first_action_minutes"`
	TimeToSetup         int `json:"time_to_setup_minutes"`
	TimeToFirstPurchase int `json:"time_to_first_purchase_hours"`
	AverageStepTime     int `json:"average_step_time_minutes"`

	// Success Indicators
	ActivationStatus   bool       `gorm:"default:false" json:"activation_status"`
	ActivationDate     *time.Time `json:"activation_date"`
	FirstValueRealized bool       `gorm:"default:false" json:"first_value_realized"`
	OnboardingScore    float64    `gorm:"default:0" json:"onboarding_score"`

	// Drop-off Analysis
	DropOffPoints        []string `gorm:"type:json" json:"drop_off_points"`
	ReengagementAttempts int      `gorm:"default:0" json:"reengagement_attempts"`
	ReengagementSuccess  bool     `gorm:"default:false" json:"reengagement_success"`
}

// UserPersonalization manages personalized experiences
type UserPersonalization struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Preferences
	UIPreferences      map[string]interface{} `gorm:"type:json" json:"ui_preferences"`
	ContentPreferences map[string]interface{} `gorm:"type:json" json:"content_preferences"`
	FeaturePreferences map[string]bool        `gorm:"type:json" json:"feature_preferences"`
	DisplayPreferences map[string]interface{} `gorm:"type:json" json:"display_preferences"`

	// Recommendations
	ProductRecommendations []map[string]interface{} `gorm:"type:json" json:"product_recommendations"`
	ContentRecommendations []map[string]interface{} `gorm:"type:json" json:"content_recommendations"`
	ActionRecommendations  []map[string]interface{} `gorm:"type:json" json:"action_recommendations"`
	RecommendationHistory  []map[string]interface{} `gorm:"type:json" json:"recommendation_history"`

	// Behavioral Preferences
	BrowsingPatterns    map[string]interface{} `gorm:"type:json" json:"browsing_patterns"`
	InteractionPatterns map[string]interface{} `gorm:"type:json" json:"interaction_patterns"`
	TimePreferences     map[string]string      `gorm:"type:json" json:"time_preferences"`
	DevicePreferences   map[string]interface{} `gorm:"type:json" json:"device_preferences"`

	// Dynamic Content
	PersonalizedContent []map[string]interface{}   `gorm:"type:json" json:"personalized_content"`
	DynamicPricing      map[string]decimal.Decimal `gorm:"type:json" json:"dynamic_pricing"`
	TargetedOffers      []map[string]interface{}   `gorm:"type:json" json:"targeted_offers"`
	CustomizedFeatures  []string                   `gorm:"type:json" json:"customized_features"`

	// AI/ML Settings
	AIAssistanceEnabled  bool    `gorm:"default:true" json:"ai_assistance_enabled"`
	PredictiveText       bool    `gorm:"default:true" json:"predictive_text"`
	SmartSuggestions     bool    `gorm:"default:true" json:"smart_suggestions"`
	AutoPersonalization  bool    `gorm:"default:true" json:"auto_personalization"`
	PersonalizationScore float64 `gorm:"default:0" json:"personalization_score"`
}

// UserRetention manages retention strategies
type UserRetention struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Retention Status
	RetentionStatus     string     `gorm:"type:varchar(20)" json:"retention_status"`
	RetentionScore      float64    `gorm:"default:0" json:"retention_score"`
	ChurnRisk           float64    `gorm:"default:0" json:"churn_risk"`
	ChurnPredictionDate *time.Time `json:"churn_prediction_date"`

	// Engagement Tracking
	LastActiveDate    *time.Time `json:"last_active_date"`
	ActivityFrequency float64    `gorm:"default:0" json:"activity_frequency"`
	EngagementTrend   string     `gorm:"type:varchar(20)" json:"engagement_trend"`
	InactivityPeriod  int        `json:"inactivity_period_days"`

	// Retention Strategies
	ActiveStrategies      []string                 `gorm:"type:json" json:"active_strategies"`
	StrategyEffectiveness map[string]float64       `gorm:"type:json" json:"strategy_effectiveness"`
	InterventionHistory   []map[string]interface{} `gorm:"type:json" json:"intervention_history"`
	WinBackCampaigns      []map[string]interface{} `gorm:"type:json" json:"win_back_campaigns"`

	// Incentives
	RetentionIncentives []map[string]interface{} `gorm:"type:json" json:"retention_incentives"`
	IncentivesRedeemed  []map[string]interface{} `gorm:"type:json" json:"incentives_redeemed"`
	IncentiveValue      decimal.Decimal          `gorm:"type:decimal(10,2)" json:"incentive_value"`

	// Re-engagement
	ReengagementAttempts int        `gorm:"default:0" json:"reengagement_attempts"`
	LastReengagementDate *time.Time `json:"last_reengagement_date"`
	ReengagementChannels []string   `gorm:"type:json" json:"reengagement_channels"`
	ReengagementSuccess  bool       `gorm:"default:false" json:"reengagement_success"`

	// Exit Analysis
	ExitIntentDetected  bool     `gorm:"default:false" json:"exit_intent_detected"`
	ExitReasons         []string `gorm:"type:json" json:"exit_reasons"`
	ExitSurveyCompleted bool     `gorm:"default:false" json:"exit_survey_completed"`
	ExitFeedback        string   `gorm:"type:text" json:"exit_feedback"`
}

// UserSegmentation manages user segmentation data
type UserSegmentation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Segment Membership
	PrimarySegment    string   `gorm:"type:varchar(50)" json:"primary_segment"`
	SecondarySegments []string `gorm:"type:json" json:"secondary_segments"`
	MicroSegments     []string `gorm:"type:json" json:"micro_segments"`
	DynamicSegments   []string `gorm:"type:json" json:"dynamic_segments"`

	// Segment Scores
	SegmentScores        map[string]float64       `gorm:"type:json" json:"segment_scores"`
	SegmentProbabilities map[string]float64       `gorm:"type:json" json:"segment_probabilities"`
	SegmentTransitions   []map[string]interface{} `gorm:"type:json" json:"segment_transitions"`

	// Behavioral Segments
	BehaviorCluster    string `gorm:"type:varchar(50)" json:"behavior_cluster"`
	PurchaseBehavior   string `gorm:"type:varchar(50)" json:"purchase_behavior"`
	EngagementLevel    string `gorm:"type:varchar(20)" json:"engagement_level"`
	TechnologyAdoption string `gorm:"type:varchar(20)" json:"technology_adoption"`

	// Value Segments
	ValueTier            string          `gorm:"type:varchar(20)" json:"value_tier"`
	RevenueContribution  decimal.Decimal `gorm:"type:decimal(15,2)" json:"revenue_contribution"`
	ProfitabilitySegment string          `gorm:"type:varchar(20)" json:"profitability_segment"`
	GrowthPotential      string          `gorm:"type:varchar(20)" json:"growth_potential"`

	// Demographic Segments
	AgeGroup          string `gorm:"type:varchar(20)" json:"age_group"`
	IncomeSegment     string `gorm:"type:varchar(20)" json:"income_segment"`
	GeographicSegment string `gorm:"type:varchar(50)" json:"geographic_segment"`
	LifestyleSegment  string `gorm:"type:varchar(50)" json:"lifestyle_segment"`

	// Psychographic Segments
	PersonalityType string   `gorm:"type:varchar(20)" json:"personality_type"`
	Values          []string `gorm:"type:json" json:"values"`
	Interests       []string `gorm:"type:json" json:"interests"`
	Motivations     []string `gorm:"type:json" json:"motivations"`
}

// UserPredictiveNeeds identifies and predicts user needs
type UserPredictiveNeeds struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Current Needs
	IdentifiedNeeds []map[string]interface{} `gorm:"type:json" json:"identified_needs"`
	NeedPriorities  map[string]int           `gorm:"type:json" json:"need_priorities"`
	UnmetNeeds      []string                 `gorm:"type:json" json:"unmet_needs"`
	SatisfiedNeeds  []string                 `gorm:"type:json" json:"satisfied_needs"`

	// Future Needs Prediction
	PredictedNeeds       []map[string]interface{} `gorm:"type:json" json:"predicted_needs"`
	NeedProbabilities    map[string]float64       `gorm:"type:json" json:"need_probabilities"`
	NeedTimeline         map[string]time.Time     `gorm:"type:json" json:"need_timeline"`
	LifeEventPredictions []map[string]interface{} `gorm:"type:json" json:"life_event_predictions"`

	// Need Categories
	FinancialNeeds   []string `gorm:"type:json" json:"financial_needs"`
	SecurityNeeds    []string `gorm:"type:json" json:"security_needs"`
	ConvenienceNeeds []string `gorm:"type:json" json:"convenience_needs"`
	SocialNeeds      []string `gorm:"type:json" json:"social_needs"`

	// Solution Matching
	RecommendedSolutions []map[string]interface{} `gorm:"type:json" json:"recommended_solutions"`
	SolutionFitScores    map[string]float64       `gorm:"type:json" json:"solution_fit_scores"`
	AdoptionLikelihood   map[string]float64       `gorm:"type:json" json:"adoption_likelihood"`

	// Proactive Outreach
	OutreachScheduled bool       `gorm:"default:false" json:"outreach_scheduled"`
	OutreachTiming    *time.Time `json:"outreach_timing"`
	OutreachChannel   string     `gorm:"type:varchar(50)" json:"outreach_channel"`
	OutreachMessage   string     `gorm:"type:text" json:"outreach_message"`
}

// UserFeedback manages user feedback and satisfaction
type UserFeedback struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Feedback Entry
	FeedbackID   string `gorm:"type:varchar(100);uniqueIndex" json:"feedback_id"`
	FeedbackType string `gorm:"type:varchar(50)" json:"feedback_type"`
	Category     string `gorm:"type:varchar(50)" json:"category"`
	Subject      string `gorm:"type:varchar(255)" json:"subject"`
	Content      string `gorm:"type:text" json:"content"`
	Rating       *int   `json:"rating"`
	Sentiment    string `gorm:"type:varchar(20)" json:"sentiment"`

	// Context
	RelatedEntity   string `gorm:"type:varchar(50)" json:"related_entity"`
	RelatedEntityID string `gorm:"type:varchar(100)" json:"related_entity_id"`
	Channel         string `gorm:"type:varchar(50)" json:"channel"`
	Platform        string `gorm:"type:varchar(50)" json:"platform"`

	// Response
	ResponseRequired bool       `gorm:"default:false" json:"response_required"`
	ResponseProvided bool       `gorm:"default:false" json:"response_provided"`
	ResponseDate     *time.Time `json:"response_date"`
	ResponseContent  string     `gorm:"type:text" json:"response_content"`
	RespondedBy      *uuid.UUID `gorm:"type:uuid" json:"responded_by"`

	// Follow-up
	FollowUpRequired  bool       `gorm:"default:false" json:"follow_up_required"`
	FollowUpCompleted bool       `gorm:"default:false" json:"follow_up_completed"`
	FollowUpDate      *time.Time `json:"follow_up_date"`
	ActionTaken       string     `gorm:"type:text" json:"action_taken"`

	// Analytics
	ImpactScore float64  `gorm:"default:0" json:"impact_score"`
	Priority    string   `gorm:"type:varchar(20)" json:"priority"`
	Tags        []string `gorm:"type:json" json:"tags"`
	Visibility  string   `gorm:"type:varchar(20)" json:"visibility"`
}

// TableName returns the table name
func (UserJourneyMapping) TableName() string {
	return "user_journey_mapping"
}

// TableName returns the table name
func (UserOnboarding) TableName() string {
	return "user_onboarding"
}

// TableName returns the table name
func (UserPersonalization) TableName() string {
	return "user_personalization"
}

// TableName returns the table name
func (UserRetention) TableName() string {
	return "user_retention"
}

// TableName returns the table name
func (UserSegmentation) TableName() string {
	return "user_segmentation"
}

// TableName returns the table name
func (UserPredictiveNeeds) TableName() string {
	return "user_predictive_needs"
}

// TableName returns the table name
func (UserFeedback) TableName() string {
	return "user_feedback"
}
