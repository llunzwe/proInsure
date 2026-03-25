package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceRecommendations manages personalized recommendations
type DeviceRecommendations struct {
	database.BaseModel
	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID      uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	GeneratedAt time.Time `json:"generated_at"`

	// Personalized Insurance Recommendations
	InsuranceRecommendations string  `gorm:"type:json" json:"insurance_recommendations"` // JSON array
	CurrentCoverageGaps      string  `gorm:"type:json" json:"current_coverage_gaps"`     // JSON array
	RecommendedAddOns        string  `gorm:"type:json" json:"recommended_add_ons"`       // JSON array
	PersonalizationScore     float64 `json:"personalization_score"`                      // 0-100

	// Coverage Optimization Suggestions
	OptimalCoverageLevel  string  `json:"optimal_coverage_level"`
	CoverageAdjustments   string  `gorm:"type:json" json:"coverage_adjustments"` // JSON array
	OverInsuredAreas      string  `gorm:"type:json" json:"over_insured_areas"`   // JSON array
	UnderInsuredAreas     string  `gorm:"type:json" json:"under_insured_areas"`  // JSON array
	OptimizationPotential float64 `json:"optimization_potential"`                // savings percentage

	// Cost-saving Recommendations
	MonthlySavingsPotential float64 `json:"monthly_savings_potential"`
	AnnualSavingsPotential  float64 `json:"annual_savings_potential"`
	DiscountOpportunities   string  `gorm:"type:json" json:"discount_opportunities"` // JSON array
	BundleOptions           string  `gorm:"type:json" json:"bundle_options"`         // JSON array
	CostReductionTips       string  `gorm:"type:json" json:"cost_reduction_tips"`    // JSON array

	// Feature Recommendations
	RecommendedFeatures   string  `gorm:"type:json" json:"recommended_features"` // JSON array
	UnusedFeatures        string  `gorm:"type:json" json:"unused_features"`      // JSON array
	FeatureUpgrades       string  `gorm:"type:json" json:"feature_upgrades"`     // JSON array
	FeatureRelevanceScore float64 `json:"feature_relevance_score"`               // 0-100

	// Service Recommendations
	ServiceUpgrades     string `gorm:"type:json" json:"service_upgrades"`     // JSON array
	SupportPlanOptions  string `gorm:"type:json" json:"support_plan_options"` // JSON array
	PreferredChannels   string `gorm:"type:json" json:"preferred_channels"`   // JSON array
	ServiceOptimization string `gorm:"type:json" json:"service_optimization"` // JSON array

	// Upgrade Recommendations
	DeviceUpgradeOptions string  `gorm:"type:json" json:"device_upgrade_options"` // JSON array
	UpgradeTimeline      string  `json:"upgrade_timeline"`
	TradeInValue         float64 `json:"trade_in_value"`
	UpgradeBenefits      string  `gorm:"type:json" json:"upgrade_benefits"` // JSON array

	// Maintenance Recommendations
	MaintenanceSchedule   string  `gorm:"type:json" json:"maintenance_schedule"`   // JSON array
	PreventiveCare        string  `gorm:"type:json" json:"preventive_care"`        // JSON array
	RepairRecommendations string  `gorm:"type:json" json:"repair_recommendations"` // JSON array
	MaintenanceSavings    float64 `json:"maintenance_savings"`

	// Security Recommendations
	SecurityImprovements string  `gorm:"type:json" json:"security_improvements"` // JSON array
	PrivacyEnhancements  string  `gorm:"type:json" json:"privacy_enhancements"`  // JSON array
	SecurityScore        float64 `json:"security_score"`                         // 0-100
	RiskMitigationSteps  string  `gorm:"type:json" json:"risk_mitigation_steps"` // JSON array

	// Usage Optimization Tips
	UsageOptimization   string `gorm:"type:json" json:"usage_optimization"`   // JSON array
	PerformanceTips     string `gorm:"type:json" json:"performance_tips"`     // JSON array
	BatteryOptimization string `gorm:"type:json" json:"battery_optimization"` // JSON array
	StorageOptimization string `gorm:"type:json" json:"storage_optimization"` // JSON array

	// Personalization Effectiveness
	AcceptanceRate     float64 `json:"acceptance_rate"`     // percentage
	ImplementationRate float64 `json:"implementation_rate"` // percentage
	SatisfactionImpact float64 `json:"satisfaction_impact"` // percentage improvement
	RevenueImpact      float64 `json:"revenue_impact"`
	EngagementImpact   float64 `json:"engagement_impact"` // percentage

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// DeviceLoyaltyProgram manages loyalty rewards and benefits
type DeviceLoyaltyProgram struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`

	// Points and Rewards Tracking
	CurrentPoints  int        `json:"current_points"`
	LifetimePoints int        `json:"lifetime_points"`
	PointsExpiring int        `json:"points_expiring"`
	ExpirationDate *time.Time `json:"expiration_date"`
	PointValue     float64    `json:"point_value"` // monetary value per point

	// Loyalty Tier Status
	CurrentTier      string    `json:"current_tier"` // bronze, silver, gold, platinum
	TierStartDate    time.Time `json:"tier_start_date"`
	NextTier         string    `json:"next_tier"`
	PointsToNextTier int       `json:"points_to_next_tier"`
	TierBenefits     string    `gorm:"type:json" json:"tier_benefits"` // JSON array

	// Benefits Utilization
	AvailableBenefits  string  `gorm:"type:json" json:"available_benefits"` // JSON array
	UsedBenefits       string  `gorm:"type:json" json:"used_benefits"`      // JSON array
	BenefitUtilization float64 `json:"benefit_utilization"`                 // percentage
	SavedAmount        float64 `json:"saved_amount"`

	// Reward Redemption History
	TotalRedemptions  int        `json:"total_redemptions"`
	LastRedemption    *time.Time `json:"last_redemption"`
	RedemptionHistory string     `gorm:"type:json" json:"redemption_history"` // JSON array
	PreferredRewards  string     `gorm:"type:json" json:"preferred_rewards"`  // JSON array

	// Point Earning Patterns
	MonthlyEarning    int     `json:"monthly_earning"`
	AverageEarning    float64 `json:"average_earning"`
	EarningStreaks    int     `json:"earning_streaks"` // consecutive months
	BonusPoints       int     `json:"bonus_points"`
	EarningMultiplier float64 `json:"earning_multiplier"`

	// Tier Upgrade Progress
	QualifyingPoints     int        `json:"qualifying_points"`
	QualifyingPeriod     string     `json:"qualifying_period"`
	ProgressPercentage   float64    `json:"progress_percentage"`
	EstimatedUpgradeDate *time.Time `json:"estimated_upgrade_date"`

	// Special Offers Tracking
	ActiveOffers    string  `gorm:"type:json" json:"active_offers"`    // JSON array
	ClaimedOffers   string  `gorm:"type:json" json:"claimed_offers"`   // JSON array
	ExclusiveOffers string  `gorm:"type:json" json:"exclusive_offers"` // JSON array
	OfferConversion float64 `json:"offer_conversion"`                  // percentage

	// Referral Rewards
	ReferralCode    string  `json:"referral_code"`
	TotalReferrals  int     `json:"total_referrals"`
	ActiveReferrals int     `json:"active_referrals"`
	ReferralPoints  int     `json:"referral_points"`
	ReferralBonuses float64 `json:"referral_bonuses"`

	// Anniversary Benefits
	AccountAnniversary  time.Time `json:"account_anniversary"`
	AnniversaryBenefits string    `gorm:"type:json" json:"anniversary_benefits"` // JSON array
	MilestoneBenefits   string    `gorm:"type:json" json:"milestone_benefits"`   // JSON array
	YearsActive         int       `json:"years_active"`

	// Loyalty Program ROI
	CustomerValue   float64 `json:"customer_value"`
	RetentionImpact float64 `json:"retention_impact"` // percentage
	SpendIncrease   float64 `json:"spend_increase"`   // percentage
	EngagementScore float64 `json:"engagement_score"` // 0-100
	ProgramROI      float64 `json:"program_roi"`      // percentage

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// DeviceFeedback collects and tracks user feedback
type DeviceFeedback struct {
	database.BaseModel
	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	SubmittedAt time.Time `json:"submitted_at"`

	// User Feedback Collection
	FeedbackType     string `json:"feedback_type"` // rating, review, suggestion, complaint, compliment
	FeedbackCategory string `json:"feedback_category"`
	FeedbackChannel  string `json:"feedback_channel"` // app, web, email, phone, social
	FeedbackContent  string `gorm:"type:text" json:"feedback_content"`

	// Rating History
	Rating          float64 `json:"rating"`                            // 0-5
	PreviousRatings string  `gorm:"type:json" json:"previous_ratings"` // JSON array
	AverageRating   float64 `json:"average_rating"`
	RatingTrend     string  `json:"rating_trend"` // improving, stable, declining

	// Review Submissions
	ReviewTitle      string `json:"review_title"`
	ReviewText       string `gorm:"type:text" json:"review_text"`
	ReviewPlatform   string `json:"review_platform"`
	ReviewVisibility string `json:"review_visibility"` // public, private
	ReviewHelpful    int    `json:"review_helpful"`    // helpful votes

	// Feature Requests
	RequestedFeatures string `gorm:"type:json" json:"requested_features"` // JSON array
	FeaturePriority   string `json:"feature_priority"`                    // low, medium, high, critical
	FeatureStatus     string `json:"feature_status"`                      // submitted, reviewing, planned, implemented
	FeatureVotes      int    `json:"feature_votes"`

	// Bug Reports
	BugDescription    string `gorm:"type:text" json:"bug_description"`
	BugSeverity       string `json:"bug_severity"`                       // minor, major, critical
	BugStatus         string `json:"bug_status"`                         // reported, confirmed, fixing, resolved
	BugResolutionTime *int   `json:"bug_resolution_time"`                // hours
	AffectedFeatures  string `gorm:"type:json" json:"affected_features"` // JSON array

	// Improvement Suggestions
	Suggestions          string  `gorm:"type:json" json:"suggestions"` // JSON array
	SuggestionImpact     string  `json:"suggestion_impact"`            // low, medium, high
	ImplementationStatus string  `json:"implementation_status"`
	EstimatedValue       float64 `json:"estimated_value"`

	// Complaint Tracking
	ComplaintSeverity   string  `json:"complaint_severity"` // low, medium, high
	ComplaintStatus     string  `json:"complaint_status"`   // open, investigating, resolved
	ResolutionTime      *int    `json:"resolution_time"`    // hours
	CompensationOffered float64 `json:"compensation_offered"`
	CustomerSatisfied   *bool   `json:"customer_satisfied"`

	// Compliment Tracking
	ComplimentCategory  string `json:"compliment_category"`
	ComplimentedFeature string `json:"complimented_feature"`
	ComplimentedStaff   string `json:"complimented_staff"`
	SharedInternally    bool   `json:"shared_internally"`

	// Feedback Response Tracking
	ResponseProvided     bool     `json:"response_provided"`
	ResponseTime         *int     `json:"response_time"` // hours
	ResponseChannel      string   `json:"response_channel"`
	ResponseSatisfaction *float64 `json:"response_satisfaction"` // 0-5
	FollowUpRequired     bool     `json:"follow_up_required"`

	// Feedback Action Items
	ActionItems    string     `gorm:"type:json" json:"action_items"` // JSON array
	ActionPriority string     `json:"action_priority"`
	ActionOwner    string     `json:"action_owner"`
	ActionDeadline *time.Time `json:"action_deadline"`
	ActionStatus   string     `json:"action_status"`

	// Sentiment Analysis
	SentimentScore    float64 `json:"sentiment_score"`    // -1 to 1
	SentimentCategory string  `json:"sentiment_category"` // positive, neutral, negative
	EmotionDetected   string  `json:"emotion_detected"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// DeviceUserJourney tracks customer journey and experience
type DeviceUserJourney struct {
	database.BaseModel
	DeviceID         uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	JourneyStartDate time.Time `json:"journey_start_date"`

	// Customer Journey Mapping
	CurrentStage     string `json:"current_stage"`                      // awareness, consideration, purchase, retention, advocacy
	JourneyPhase     string `json:"journey_phase"`                      // onboarding, active, at-risk, churned
	JourneyDuration  int    `json:"journey_duration"`                   // days
	StageTransitions string `gorm:"type:json" json:"stage_transitions"` // JSON array

	// Touchpoint Tracking
	TotalTouchpoints     int     `json:"total_touchpoints"`
	RecentTouchpoints    string  `gorm:"type:json" json:"recent_touchpoints"`    // JSON array
	TouchpointFrequency  float64 `json:"touchpoint_frequency"`                   // per month
	PreferredTouchpoints string  `gorm:"type:json" json:"preferred_touchpoints"` // JSON array
	TouchpointQuality    float64 `json:"touchpoint_quality"`                     // 0-100

	// Journey Stage Identification
	StageProgress    float64 `json:"stage_progress"`     // percentage
	StageHealthScore float64 `json:"stage_health_score"` // 0-100
	NextLikelyStage  string  `json:"next_likely_stage"`
	StageRisks       string  `gorm:"type:json" json:"stage_risks"` // JSON array

	// Pain Point Identification
	IdentifiedPainPoints string  `gorm:"type:json" json:"identified_pain_points"` // JSON array
	PainPointSeverity    string  `gorm:"type:json" json:"pain_point_severity"`    // JSON object
	ResolutionStatus     string  `gorm:"type:json" json:"resolution_status"`      // JSON object
	PainPointImpact      float64 `json:"pain_point_impact"`                       // 0-100

	// Moment of Truth Tracking
	CriticalMoments string  `gorm:"type:json" json:"critical_moments"` // JSON array
	MomentOutcomes  string  `gorm:"type:json" json:"moment_outcomes"`  // JSON object
	PositiveMoments int     `json:"positive_moments"`
	NegativeMoments int     `json:"negative_moments"`
	MomentRecovery  float64 `json:"moment_recovery"` // success rate

	// Channel Preference Tracking
	PreferredChannel     string  `json:"preferred_channel"`
	ChannelUsage         string  `gorm:"type:json" json:"channel_usage"` // JSON object
	CrossChannelJourney  bool    `json:"cross_channel_journey"`
	ChannelSwitchPoints  string  `gorm:"type:json" json:"channel_switch_points"` // JSON array
	ChannelEffectiveness float64 `json:"channel_effectiveness"`                  // 0-100

	// Journey Optimization Opportunities
	OptimizationAreas   string  `gorm:"type:json" json:"optimization_areas"`   // JSON array
	AutomationPotential float64 `json:"automation_potential"`                  // percentage
	PersonalizationGaps string  `gorm:"type:json" json:"personalization_gaps"` // JSON array
	ExperienceGaps      string  `gorm:"type:json" json:"experience_gaps"`      // JSON array

	// Conversion Funnel Analysis
	FunnelStage           string  `json:"funnel_stage"`
	ConversionProbability float64 `json:"conversion_probability"`               // 0-100
	FunnelDropoffs        string  `gorm:"type:json" json:"funnel_dropoffs"`     // JSON array
	ConversionBarriers    string  `gorm:"type:json" json:"conversion_barriers"` // JSON array
	ConversionDrivers     string  `gorm:"type:json" json:"conversion_drivers"`  // JSON array

	// Abandonment Point Tracking
	AbandonmentRisk       float64 `json:"abandonment_risk"`                        // 0-100
	AbandonmentIndicators string  `gorm:"type:json" json:"abandonment_indicators"` // JSON array
	RecoveryAttempts      int     `json:"recovery_attempts"`
	RecoverySuccess       bool    `json:"recovery_success"`

	// Journey Personalization
	PersonalizationLevel  float64 `json:"personalization_level"`                  // 0-100
	PersonalizedElements  string  `gorm:"type:json" json:"personalized_elements"` // JSON array
	PersonalizationImpact float64 `json:"personalization_impact"`
	NextBestAction        string  `json:"next_best_action"`

	// Journey Metrics
	JourneyScore      float64 `json:"journey_score"`      // 0-100
	EngagementLevel   float64 `json:"engagement_level"`   // 0-100
	SatisfactionLevel float64 `json:"satisfaction_level"` // 0-100
	LoyaltyIndicator  float64 `json:"loyalty_indicator"`  // 0-100

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}
