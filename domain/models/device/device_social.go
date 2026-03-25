package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// DeviceSocial represents social and community features
type DeviceSocial struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Device Owner Communities
	CommunityMemberships datatypes.JSON `gorm:"type:json" json:"community_memberships"` // []Community
	ActiveCommunities    int            `gorm:"type:int" json:"active_communities"`
	CommunityRole        string         `gorm:"type:varchar(50)" json:"community_role"` // member, moderator, expert
	CommunityPoints      int            `gorm:"type:int" json:"community_points"`
	CommunityBadges      datatypes.JSON `gorm:"type:json" json:"community_badges"` // []Badge

	// Peer Comparisons
	PeerGroup         string         `gorm:"type:varchar(100)" json:"peer_group"`
	PeerRanking       int            `gorm:"type:int" json:"peer_ranking"`
	PeerComparisons   datatypes.JSON `gorm:"type:json" json:"peer_comparisons"` // []Comparison
	BetterThanPercent float64        `gorm:"type:decimal(5,2)" json:"better_than_percent"`
	PeerBenchmarks    datatypes.JSON `gorm:"type:json" json:"peer_benchmarks"` // map[string]float64

	// Device Reviews & Ratings
	UserRating         float64        `gorm:"type:decimal(3,2)" json:"user_rating"` // 1-5 stars
	ReviewText         string         `gorm:"type:text" json:"review_text"`
	ReviewDate         *time.Time     `gorm:"type:timestamp" json:"review_date,omitempty"`
	ReviewHelpfulCount int            `gorm:"type:int" json:"review_helpful_count"`
	ReviewVerified     bool           `gorm:"type:boolean;default:false" json:"review_verified"`
	DeviceRatings      datatypes.JSON `gorm:"type:json" json:"device_ratings"` // []Rating

	// Tips Sharing Platform
	SharedTips       datatypes.JSON `gorm:"type:json" json:"shared_tips"`   // []Tip
	TipsReceived     datatypes.JSON `gorm:"type:json" json:"tips_received"` // []Tip
	TipContributions int            `gorm:"type:int" json:"tip_contributions"`
	TipUpvotes       int            `gorm:"type:int" json:"tip_upvotes"`
	TipExpertStatus  bool           `gorm:"type:boolean;default:false" json:"tip_expert_status"`

	// Expert Consultations
	ExpertConsultations  datatypes.JSON `gorm:"type:json" json:"expert_consultations"` // []Consultation
	ConsultationRequests int            `gorm:"type:int" json:"consultation_requests"`
	ExpertRating         float64        `gorm:"type:decimal(3,2)" json:"expert_rating"`
	ExpertAvailable      bool           `gorm:"type:boolean;default:false" json:"expert_available"`
	ConsultationTopics   datatypes.JSON `gorm:"type:json" json:"consultation_topics"` // []string

	// Social Sharing
	ShareCount      int            `gorm:"type:int" json:"share_count"`
	SharePlatforms  datatypes.JSON `gorm:"type:json" json:"share_platforms"` // map[string]int
	ViralContent    datatypes.JSON `gorm:"type:json" json:"viral_content"`   // []Content
	SocialReach     int            `gorm:"type:int" json:"social_reach"`
	InfluencerScore float64        `gorm:"type:decimal(5,2)" json:"influencer_score"`

	// Forums & Discussions
	ForumPosts      int  `gorm:"type:int" json:"forum_posts"`
	ForumReplies    int  `gorm:"type:int" json:"forum_replies"`
	HelpfulAnswers  int  `gorm:"type:int" json:"helpful_answers"`
	ForumReputation int  `gorm:"type:int" json:"forum_reputation"`
	TopContributor  bool `gorm:"type:boolean;default:false" json:"top_contributor"`

	// User Groups
	UserGroups      datatypes.JSON `gorm:"type:json" json:"user_groups"` // []UserGroup
	GroupOwner      bool           `gorm:"type:boolean;default:false" json:"group_owner"`
	GroupModerator  bool           `gorm:"type:boolean;default:false" json:"group_moderator"`
	GroupActivities datatypes.JSON `gorm:"type:json" json:"group_activities"` // []Activity

	// Following & Followers
	FollowingCount int            `gorm:"type:int" json:"following_count"`
	FollowersCount int            `gorm:"type:int" json:"followers_count"`
	Following      datatypes.JSON `gorm:"type:json" json:"following"` // []UserID
	Followers      datatypes.JSON `gorm:"type:json" json:"followers"` // []UserID

	// Status
	SocialStatus     string    `gorm:"type:varchar(50)" json:"social_status"`
	LastActivityDate time.Time `gorm:"type:timestamp" json:"last_activity_date"`
	CreatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceAutomation represents automation and workflow features
type DeviceAutomation struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Automated Claim Triggers
	AutoClaimEnabled    bool           `gorm:"type:boolean;default:false" json:"auto_claim_enabled"`
	ClaimTriggers       datatypes.JSON `gorm:"type:json" json:"claim_triggers"` // []Trigger
	AutoClaimsSubmitted int            `gorm:"type:int" json:"auto_claims_submitted"`
	LastAutoClaimDate   *time.Time     `gorm:"type:timestamp" json:"last_auto_claim_date,omitempty"`
	AutoClaimConditions datatypes.JSON `gorm:"type:json" json:"auto_claim_conditions"` // []Condition

	// Scheduled Maintenance
	MaintenanceReminders datatypes.JSON `gorm:"type:json" json:"maintenance_reminders"` // []Reminder
	NextMaintenanceAlert time.Time      `gorm:"type:timestamp" json:"next_maintenance_alert"`
	ReminderFrequency    string         `gorm:"type:varchar(50)" json:"reminder_frequency"`
	SnoozeCount          int            `gorm:"type:int" json:"snooze_count"`
	AutoScheduling       bool           `gorm:"type:boolean;default:false" json:"auto_scheduling"`

	// Renewal Automation
	AutoRenewalEnabled   bool           `gorm:"type:boolean;default:false" json:"auto_renewal_enabled"`
	RenewalSettings      datatypes.JSON `gorm:"type:json" json:"renewal_settings"` // RenewalConfig
	RenewalNotifications int            `gorm:"type:int" json:"renewal_notifications"`
	DaysBeforeRenewal    int            `gorm:"type:int" json:"days_before_renewal"`
	AutoRenewalHistory   datatypes.JSON `gorm:"type:json" json:"auto_renewal_history"` // []Renewal

	// Price Monitoring
	PriceAlerts         datatypes.JSON `gorm:"type:json" json:"price_alerts"` // []PriceAlert
	PriceDropThreshold  float64        `gorm:"type:decimal(5,2)" json:"price_drop_threshold"`
	PriceIncreaseAlert  bool           `gorm:"type:boolean;default:false" json:"price_increase_alert"`
	MarketPriceTracking bool           `gorm:"type:boolean;default:false" json:"market_price_tracking"`
	LastPriceCheck      time.Time      `gorm:"type:timestamp" json:"last_price_check"`

	// Upgrade Eligibility
	UpgradeAlerts      datatypes.JSON `gorm:"type:json" json:"upgrade_alerts"`     // []Alert
	EligibilityChecks  datatypes.JSON `gorm:"type:json" json:"eligibility_checks"` // []Check
	NextUpgradeDate    *time.Time     `gorm:"type:timestamp" json:"next_upgrade_date,omitempty"`
	UpgradeRecommended bool           `gorm:"type:boolean;default:false" json:"upgrade_recommended"`
	AutoUpgradeCheck   bool           `gorm:"type:boolean;default:false" json:"auto_upgrade_check"`

	// Workflow Automation
	ActiveWorkflows    datatypes.JSON `gorm:"type:json" json:"active_workflows"`   // []Workflow
	WorkflowTemplates  datatypes.JSON `gorm:"type:json" json:"workflow_templates"` // []Template
	CustomWorkflows    datatypes.JSON `gorm:"type:json" json:"custom_workflows"`   // []CustomWorkflow
	WorkflowExecutions int            `gorm:"type:int" json:"workflow_executions"`
	LastWorkflowRun    *time.Time     `gorm:"type:timestamp" json:"last_workflow_run,omitempty"`

	// Rule Engine
	AutomationRules datatypes.JSON `gorm:"type:json" json:"automation_rules"` // []Rule
	ActiveRules     int            `gorm:"type:int" json:"active_rules"`
	RuleExecutions  int            `gorm:"type:int" json:"rule_executions"`
	RuleViolations  datatypes.JSON `gorm:"type:json" json:"rule_violations"` // []Violation

	// Smart Notifications
	SmartNotifications   bool           `gorm:"type:boolean;default:false" json:"smart_notifications"`
	NotificationRules    datatypes.JSON `gorm:"type:json" json:"notification_rules"` // []NotificationRule
	QuietHours           datatypes.JSON `gorm:"type:json" json:"quiet_hours"`        // TimeRange
	NotificationBatching bool           `gorm:"type:boolean;default:false" json:"notification_batching"`

	// Integration Hooks
	WebhookSubscriptions datatypes.JSON `gorm:"type:json" json:"webhook_subscriptions"` // []Webhook
	APITriggers          datatypes.JSON `gorm:"type:json" json:"api_triggers"`          // []APITrigger
	EventSubscriptions   datatypes.JSON `gorm:"type:json" json:"event_subscriptions"`   // []Event

	// Status
	AutomationStatus  string    `gorm:"type:varchar(50)" json:"automation_status"`
	LastExecutionDate time.Time `gorm:"type:timestamp" json:"last_execution_date"`
	CreatedAt         time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceAdvancedAnalytics represents advanced analytics features
type DeviceAdvancedAnalytics struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Cohort Analysis
	CohortID          string         `gorm:"type:varchar(100)" json:"cohort_id"`
	CohortName        string         `gorm:"type:varchar(255)" json:"cohort_name"`
	CohortMetrics     datatypes.JSON `gorm:"type:json" json:"cohort_metrics"`    // map[string]float64
	CohortComparison  datatypes.JSON `gorm:"type:json" json:"cohort_comparison"` // Comparison
	CohortPerformance float64        `gorm:"type:decimal(5,2)" json:"cohort_performance"`

	// A/B Testing
	ABTestGroups       datatypes.JSON `gorm:"type:json" json:"ab_test_groups"` // []ABTestGroup
	ActiveExperiments  int            `gorm:"type:int" json:"active_experiments"`
	ExperimentResults  datatypes.JSON `gorm:"type:json" json:"experiment_results"` // map[string]Result
	ControlGroup       bool           `gorm:"type:boolean;default:false" json:"control_group"`
	VariantAssignments datatypes.JSON `gorm:"type:json" json:"variant_assignments"` // map[string]string

	// Conversion Funnel
	FunnelStages     datatypes.JSON `gorm:"type:json" json:"funnel_stages"` // []FunnelStage
	ConversionRate   float64        `gorm:"type:decimal(5,2)" json:"conversion_rate"`
	DropoffPoints    datatypes.JSON `gorm:"type:json" json:"dropoff_points"` // []DropoffPoint
	FunnelCompletion float64        `gorm:"type:decimal(5,2)" json:"funnel_completion"`
	TimeToConversion int            `gorm:"type:int" json:"time_to_conversion_hours"`

	// Attribution Modeling
	AttributionModel    string         `gorm:"type:varchar(100)" json:"attribution_model"` // first-touch, last-touch, linear
	TouchPoints         datatypes.JSON `gorm:"type:json" json:"touch_points"`              // []TouchPoint
	AttributionScores   datatypes.JSON `gorm:"type:json" json:"attribution_scores"`        // map[string]float64
	ChannelContribution datatypes.JSON `gorm:"type:json" json:"channel_contribution"`      // map[string]float64
	ConversionPath      datatypes.JSON `gorm:"type:json" json:"conversion_path"`           // []Step

	// Customer Journey Mapping
	JourneyStages       datatypes.JSON `gorm:"type:json" json:"journey_stages"` // []JourneyStage
	CurrentJourneyStage string         `gorm:"type:varchar(100)" json:"current_journey_stage"`
	JourneyDuration     int            `gorm:"type:int" json:"journey_duration_days"`
	JourneyScore        float64        `gorm:"type:decimal(5,2)" json:"journey_score"`
	NextBestAction      string         `gorm:"type:varchar(255)" json:"next_best_action"`

	// Predictive Analytics
	PredictiveModels   datatypes.JSON `gorm:"type:json" json:"predictive_models"` // []Model
	ChurnProbability   float64        `gorm:"type:decimal(5,2)" json:"churn_probability"`
	LTVPrediction      float64        `gorm:"type:decimal(15,2)" json:"ltv_prediction"`
	ClaimLikelihood    float64        `gorm:"type:decimal(5,2)" json:"claim_likelihood"`
	RenewalProbability float64        `gorm:"type:decimal(5,2)" json:"renewal_probability"`

	// Segmentation Analysis
	Segments         datatypes.JSON `gorm:"type:json" json:"segments"` // []Segment
	PrimarySegment   string         `gorm:"type:varchar(100)" json:"primary_segment"`
	SegmentScores    datatypes.JSON `gorm:"type:json" json:"segment_scores"`    // map[string]float64
	SegmentMigration datatypes.JSON `gorm:"type:json" json:"segment_migration"` // Migration

	// Behavioral Analytics
	BehaviorPatterns  datatypes.JSON `gorm:"type:json" json:"behavior_patterns"` // []Pattern
	EngagementScore   float64        `gorm:"type:decimal(5,2)" json:"engagement_score"`
	ActivityFrequency float64        `gorm:"type:decimal(10,2)" json:"activity_frequency"`
	FeatureAdoption   datatypes.JSON `gorm:"type:json" json:"feature_adoption"` // map[string]float64

	// Performance Metrics
	KPITracking      datatypes.JSON `gorm:"type:json" json:"kpi_tracking"`      // map[string]float64
	MetricTrends     datatypes.JSON `gorm:"type:json" json:"metric_trends"`     // map[string][]DataPoint
	AnomalyDetection datatypes.JSON `gorm:"type:json" json:"anomaly_detection"` // []Anomaly
	PerformanceIndex float64        `gorm:"type:decimal(5,2)" json:"performance_index"`

	// Status
	AnalyticsStatus  string    `gorm:"type:varchar(50)" json:"analytics_status"`
	LastAnalysisDate time.Time `gorm:"type:timestamp" json:"last_analysis_date"`
	CreatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsActiveCommunityMember checks if user is active in communities
func (ds *DeviceSocial) IsActiveCommunityMember() bool {
	return ds.ActiveCommunities > 0 && ds.CommunityPoints > 100
}

// IsInfluencer checks if user is an influencer
func (ds *DeviceSocial) IsInfluencer() bool {
	return ds.InfluencerScore > 70 || ds.FollowersCount > 1000 ||
		ds.TopContributor
}

// HasSocialEngagement checks for social engagement
func (ds *DeviceSocial) HasSocialEngagement() bool {
	return ds.ForumPosts > 0 || ds.SharedTips != nil ||
		ds.ReviewText != ""
}

// HasActiveAutomation checks if automation is active
func (da *DeviceAutomation) HasActiveAutomation() bool {
	return da.AutoClaimEnabled || da.AutoRenewalEnabled ||
		da.ActiveRules > 0 || da.ActiveWorkflows != nil
}

// NeedsRenewal checks if device needs renewal
func (da *DeviceAutomation) NeedsRenewal() bool {
	return da.RenewalNotifications > 0 && !da.AutoRenewalEnabled
}

// HasPriceAlert checks for price alerts
func (da *DeviceAutomation) HasPriceAlert() bool {
	return da.PriceAlerts != nil && da.PriceDropThreshold > 0
}

// IsInExperiment checks if device is in A/B test
func (daa *DeviceAdvancedAnalytics) IsInExperiment() bool {
	return daa.ActiveExperiments > 0 || daa.ABTestGroups != nil
}

// IsHighValue checks if device owner is high value
func (daa *DeviceAdvancedAnalytics) IsHighValue() bool {
	return daa.LTVPrediction > 10000 || daa.EngagementScore > 80
}

// IsAtRisk checks if device owner is at risk
func (daa *DeviceAdvancedAnalytics) IsAtRisk() bool {
	return daa.ChurnProbability > 0.7 || daa.RenewalProbability < 0.3
}
