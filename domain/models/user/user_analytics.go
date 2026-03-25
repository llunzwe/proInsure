package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserAnalytics provides comprehensive analytics and predictive modeling for users
type UserAnalytics struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Predictive Analytics
	ChurnProbability     float64         `gorm:"default:0" json:"churn_probability"`
	ChurnRiskCategory    string          `gorm:"type:varchar(20)" json:"churn_risk_category"` // low/medium/high/critical
	NextChurnAssessment  *time.Time      `json:"next_churn_assessment"`
	CLVPrediction        decimal.Decimal `gorm:"type:decimal(15,2)" json:"clv_prediction"`
	CLVConfidenceScore   float64         `gorm:"default:0" json:"clv_confidence_score"`
	Upgradelikelihood    float64         `gorm:"default:0" json:"upgrade_likelihood"`
	CrossSellProbability float64         `gorm:"default:0" json:"cross_sell_probability"`
	DefaultProbability   float64         `gorm:"default:0" json:"default_probability"`
	FraudProbability     float64         `gorm:"default:0" json:"fraud_probability"`
	RenewalProbability   float64         `gorm:"default:0" json:"renewal_probability"`
	RecommendationScore  float64         `gorm:"default:0" json:"recommendation_score"`
	NextBestAction       string          `gorm:"type:varchar(100)" json:"next_best_action"`
	NextBestProduct      string          `gorm:"type:varchar(100)" json:"next_best_product"`
	PredictedLTVMonths   int             `json:"predicted_ltv_months"`

	// Behavior Analytics
	SessionCount           int             `gorm:"default:0" json:"session_count"`
	AverageSessionDuration int             `json:"average_session_duration_seconds"`
	PageViewsPerSession    float64         `gorm:"default:0" json:"page_views_per_session"`
	BounceRate             float64         `gorm:"default:0" json:"bounce_rate"`
	FeatureAdoption        map[string]bool `gorm:"type:json" json:"feature_adoption"`
	LastActiveFeatures     []string        `gorm:"type:json" json:"last_active_features"`
	EngagementTrend        string          `gorm:"type:varchar(20)" json:"engagement_trend"` // increasing/stable/decreasing
	ActivityHeatmap        map[string]int  `gorm:"type:json" json:"activity_heatmap"`
	PreferredDeviceType    string          `gorm:"type:varchar(20)" json:"preferred_device_type"`
	PreferredAccessTime    string          `gorm:"type:varchar(20)" json:"preferred_access_time"`

	// Cohort Analysis
	CohortID              string          `gorm:"type:varchar(50);index" json:"cohort_id"`
	CohortName            string          `gorm:"type:varchar(100)" json:"cohort_name"`
	CohortPerformance     string          `gorm:"type:varchar(20)" json:"cohort_performance"` // above/average/below
	CohortRetentionRate   float64         `gorm:"default:0" json:"cohort_retention_rate"`
	CohortAverageRevenue  decimal.Decimal `gorm:"type:decimal(15,2)" json:"cohort_average_revenue"`
	CohortComparisonScore float64         `gorm:"default:0" json:"cohort_comparison_score"`

	// Segmentation
	PrimarySegment       string   `gorm:"type:varchar(50)" json:"primary_segment"`
	SecondarySegments    []string `gorm:"type:json" json:"secondary_segments"`
	MicroSegments        []string `gorm:"type:json" json:"micro_segments"`
	BehavioralSegment    string   `gorm:"type:varchar(50)" json:"behavioral_segment"`
	ValueSegment         string   `gorm:"type:varchar(50)" json:"value_segment"`
	LifecycleSegment     string   `gorm:"type:varchar(50)" json:"lifecycle_segment"`
	PsychographicSegment string   `gorm:"type:varchar(50)" json:"psychographic_segment"`

	// Attribution & Journey
	AcquisitionChannel     string   `gorm:"type:varchar(50)" json:"acquisition_channel"`
	AttributionPath        []string `gorm:"type:json" json:"attribution_path"`
	FirstTouchChannel      string   `gorm:"type:varchar(50)" json:"first_touch_channel"`
	LastTouchChannel       string   `gorm:"type:varchar(50)" json:"last_touch_channel"`
	ConversionPath         []string `gorm:"type:json" json:"conversion_path"`
	TimeToConversion       int      `json:"time_to_conversion_days"`
	TouchpointCount        int      `gorm:"default:0" json:"touchpoint_count"`
	CampaignResponsiveness float64  `gorm:"default:0" json:"campaign_responsiveness"`

	// Scoring & Metrics
	HealthScore       float64 `gorm:"default:50" json:"health_score"`
	ActivityScore     float64 `gorm:"default:0" json:"activity_score"`
	AdvocacyScore     float64 `gorm:"default:0" json:"advocacy_score"`
	InfluenceScore    float64 `gorm:"default:0" json:"influence_score"`
	TrustScore        float64 `gorm:"default:50" json:"trust_score"`
	SatisfactionTrend string  `gorm:"type:varchar(20)" json:"satisfaction_trend"`
	EffortScore       float64 `gorm:"default:0" json:"effort_score"`
	QualityScore      float64 `gorm:"default:0" json:"quality_score"`
	ComplianceScore   float64 `gorm:"default:0" json:"compliance_score"`
	SecurityScore     float64 `gorm:"default:0" json:"security_score"`

	// ML Model Metadata
	ModelVersion        string             `gorm:"type:varchar(20)" json:"model_version"`
	LastModelUpdate     *time.Time         `json:"last_model_update"`
	ModelConfidence     float64            `gorm:"default:0" json:"model_confidence"`
	FeaturesUsed        []string           `gorm:"type:json" json:"features_used"`
	PredictionTimestamp *time.Time         `json:"prediction_timestamp"`
	ModelPerformance    map[string]float64 `gorm:"type:json" json:"model_performance"`

	// A/B Testing
	ActiveExperiments []string               `gorm:"type:json" json:"active_experiments"`
	ExperimentGroups  map[string]string      `gorm:"type:json" json:"experiment_groups"`
	ExperimentResults map[string]interface{} `gorm:"type:json" json:"experiment_results"`
}

// UserBehaviorAnalytics tracks detailed user behavior patterns
type UserBehaviorAnalytics struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Session Analytics
	SessionID        string            `gorm:"type:varchar(100)" json:"session_id"`
	SessionStartTime time.Time         `json:"session_start_time"`
	SessionEndTime   *time.Time        `json:"session_end_time"`
	SessionDuration  int               `json:"session_duration_seconds"`
	PageViews        int               `gorm:"default:0" json:"page_views"`
	Actions          []string          `gorm:"type:json" json:"actions"`
	DeviceInfo       map[string]string `gorm:"type:json" json:"device_info"`
	Location         map[string]string `gorm:"type:json" json:"location"`
	ReferrerSource   string            `gorm:"type:varchar(255)" json:"referrer_source"`
	ExitPage         string            `gorm:"type:varchar(255)" json:"exit_page"`
	ConversionEvents []string          `gorm:"type:json" json:"conversion_events"`

	// Interaction Patterns
	ClickStream         []map[string]interface{} `gorm:"type:json" json:"click_stream"`
	ScrollDepth         float64                  `gorm:"default:0" json:"scroll_depth"`
	TimeOnPage          map[string]int           `gorm:"type:json" json:"time_on_page"`
	InteractionRate     float64                  `gorm:"default:0" json:"interaction_rate"`
	FormCompletionRate  float64                  `gorm:"default:0" json:"form_completion_rate"`
	VideoWatchTime      int                      `json:"video_watch_time_seconds"`
	DocumentsDownloaded []string                 `gorm:"type:json" json:"documents_downloaded"`
	SearchQueries       []string                 `gorm:"type:json" json:"search_queries"`
	ErrorEncountered    []string                 `gorm:"type:json" json:"errors_encountered"`
	FeedbackProvided    bool                     `gorm:"default:false" json:"feedback_provided"`
}

// UserPredictiveModeling contains ML-based predictions for user behavior
type UserPredictiveModeling struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Predictions
	NextPurchaseDate        *time.Time             `json:"next_purchase_date"`
	NextPurchaseAmount      decimal.Decimal        `gorm:"type:decimal(10,2)" json:"next_purchase_amount"`
	NextPurchaseCategory    string                 `gorm:"type:varchar(50)" json:"next_purchase_category"`
	NextClaimProbability    float64                `gorm:"default:0" json:"next_claim_probability"`
	NextClaimEstimatedValue decimal.Decimal        `gorm:"type:decimal(10,2)" json:"next_claim_estimated_value"`
	NextContactDate         *time.Time             `json:"next_contact_date"`
	NextContactReason       string                 `gorm:"type:varchar(100)" json:"next_contact_reason"`
	LifetimeEvents          map[string]interface{} `gorm:"type:json" json:"lifetime_events"`
	PropensityScores        map[string]float64     `gorm:"type:json" json:"propensity_scores"`
	TimeSeriesPredictions   map[string]interface{} `gorm:"type:json" json:"time_series_predictions"`
	AnomalyScores           map[string]float64     `gorm:"type:json" json:"anomaly_scores"`
	ClusterAssignment       string                 `gorm:"type:varchar(50)" json:"cluster_assignment"`
	SimilarUserProfiles     []uuid.UUID            `gorm:"type:json" json:"similar_user_profiles"`
	RecommendedActions      []string               `gorm:"type:json" json:"recommended_actions"`
	PredictionExplanation   map[string]interface{} `gorm:"type:json" json:"prediction_explanation"`
}

// UserCohortAnalysis tracks user cohort membership and performance
type UserCohortAnalysis struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Cohort Details
	AcquisitionCohort     string          `gorm:"type:varchar(50);index" json:"acquisition_cohort"`
	BehaviorCohort        string          `gorm:"type:varchar(50)" json:"behavior_cohort"`
	RevenueCohort         string          `gorm:"type:varchar(50)" json:"revenue_cohort"`
	RetentionCohort       string          `gorm:"type:varchar(50)" json:"retention_cohort"`
	CohortJoinDate        time.Time       `json:"cohort_join_date"`
	CohortPosition        int             `json:"cohort_position"`
	CohortPercentile      float64         `gorm:"default:0" json:"cohort_percentile"`
	DaysInCohort          int             `gorm:"default:0" json:"days_in_cohort"`
	CohortGraduationDate  *time.Time      `json:"cohort_graduation_date"`
	ComparedToAverage     float64         `gorm:"default:0" json:"compared_to_average"`
	CohortContribution    decimal.Decimal `gorm:"type:decimal(15,2)" json:"cohort_contribution"`
	CohortRetentionStatus string          `gorm:"type:varchar(20)" json:"cohort_retention_status"`
	CohortEngagementLevel string          `gorm:"type:varchar(20)" json:"cohort_engagement_level"`
}

// UserAttributionModeling tracks marketing attribution and channel effectiveness
type UserAttributionModeling struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Attribution Models
	LinearAttribution        map[string]float64 `gorm:"type:json" json:"linear_attribution"`
	TimeDecayAttribution     map[string]float64 `gorm:"type:json" json:"time_decay_attribution"`
	PositionBasedAttribution map[string]float64 `gorm:"type:json" json:"position_based_attribution"`
	DataDrivenAttribution    map[string]float64 `gorm:"type:json" json:"data_driven_attribution"`
	CustomAttribution        map[string]float64 `gorm:"type:json" json:"custom_attribution"`

	// Channel Performance
	ChannelROI                map[string]float64         `gorm:"type:json" json:"channel_roi"`
	ChannelConversionRate     map[string]float64         `gorm:"type:json" json:"channel_conversion_rate"`
	ChannelEngagement         map[string]float64         `gorm:"type:json" json:"channel_engagement"`
	ChannelCostPerAcquisition map[string]decimal.Decimal `gorm:"type:json" json:"channel_cost_per_acquisition"`
	PreferredChannels         []string                   `gorm:"type:json" json:"preferred_channels"`
	ChannelJourney            []map[string]interface{}   `gorm:"type:json" json:"channel_journey"`
	MultiTouchPoints          int                        `gorm:"default:0" json:"multi_touch_points"`
	AttributionWindow         int                        `json:"attribution_window_days"`
}

// UserLifecycleAnalytics tracks user lifecycle stages and progression
type UserLifecycleAnalytics struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Lifecycle Stages
	CurrentStage          string                   `gorm:"type:varchar(50)" json:"current_stage"`
	StageHistory          []map[string]interface{} `gorm:"type:json" json:"stage_history"`
	StageTransitionDates  map[string]time.Time     `gorm:"type:json" json:"stage_transition_dates"`
	TimeInCurrentStage    int                      `json:"time_in_current_stage_days"`
	NextPredictedStage    string                   `gorm:"type:varchar(50)" json:"next_predicted_stage"`
	StageProgressionScore float64                  `gorm:"default:0" json:"stage_progression_score"`
	StageCompletionRate   float64                  `gorm:"default:0" json:"stage_completion_rate"`
	MilestonesAchieved    []string                 `gorm:"type:json" json:"milestones_achieved"`
	NextMilestone         string                   `gorm:"type:varchar(100)" json:"next_milestone"`
	LifecycleValue        decimal.Decimal          `gorm:"type:decimal(15,2)" json:"lifecycle_value"`
	MaturityLevel         string                   `gorm:"type:varchar(20)" json:"maturity_level"`
	ActivationStatus      string                   `gorm:"type:varchar(20)" json:"activation_status"`
	ActivationDate        *time.Time               `json:"activation_date"`
	ReactivationAttempts  int                      `gorm:"default:0" json:"reactivation_attempts"`
	WinBackEligibility    bool                     `gorm:"default:false" json:"win_back_eligibility"`
}

// TableName returns the table name
func (UserAnalytics) TableName() string {
	return "user_analytics"
}

// TableName returns the table name
func (UserBehaviorAnalytics) TableName() string {
	return "user_behavior_analytics"
}

// TableName returns the table name
func (UserPredictiveModeling) TableName() string {
	return "user_predictive_modeling"
}

// TableName returns the table name
func (UserCohortAnalysis) TableName() string {
	return "user_cohort_analysis"
}

// TableName returns the table name
func (UserAttributionModeling) TableName() string {
	return "user_attribution_modeling"
}

// TableName returns the table name
func (UserLifecycleAnalytics) TableName() string {
	return "user_lifecycle_analytics"
}
