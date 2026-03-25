package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceSatisfactionScore tracks user satisfaction metrics
type DeviceSatisfactionScore struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	MeasurementDate time.Time `json:"measurement_date"`

	// User Satisfaction Tracking
	OverallSatisfaction  float64 `json:"overall_satisfaction"`   // 0-100
	SatisfactionTrend    string  `json:"satisfaction_trend"`     // increasing, stable, decreasing
	LastInteractionScore float64 `json:"last_interaction_score"` // 0-100
	AverageSatisfaction  float64 `json:"average_satisfaction"`   // lifetime average

	// NPS (Net Promoter Score)
	NPSScore              int       `json:"nps_score"`               // -100 to 100
	NPSCategory           string    `json:"nps_category"`            // promoter, passive, detractor
	LikelihoodToRecommend int       `json:"likelihood_to_recommend"` // 0-10
	NPSLastUpdated        time.Time `json:"nps_last_updated"`
	NPSTrend              string    `json:"nps_trend"` // improving, stable, declining

	// CSAT Scores
	CSATScore        float64    `json:"csat_score"` // 0-100
	CSATResponses    int        `json:"csat_responses"`
	CSATResponseRate float64    `json:"csat_response_rate"` // percentage
	LastCSATSurvey   *time.Time `json:"last_csat_survey"`
	CSATTrend        string     `json:"csat_trend"`

	// Feature Satisfaction Ratings
	FeatureSatisfaction   string  `gorm:"type:json" json:"feature_satisfaction"` // JSON object
	MostSatisfiedFeature  string  `json:"most_satisfied_feature"`
	LeastSatisfiedFeature string  `json:"least_satisfied_feature"`
	FeatureUsageScore     float64 `json:"feature_usage_score"` // 0-100

	// Service Satisfaction Metrics
	ServiceSatisfaction   float64 `json:"service_satisfaction"`   // 0-100
	ResponseTimeRating    float64 `json:"response_time_rating"`   // 0-5
	ResolutionRating      float64 `json:"resolution_rating"`      // 0-5
	ProfessionalismRating float64 `json:"professionalism_rating"` // 0-5
	ServiceInteractions   int     `json:"service_interactions"`

	// Support Experience Ratings
	SupportSatisfaction      float64 `json:"support_satisfaction"` // 0-100
	SupportTickets           int     `json:"support_tickets"`
	AverageResolutionTime    int     `json:"average_resolution_time"`  // hours
	FirstContactResolution   float64 `json:"first_contact_resolution"` // percentage
	SupportChannelPreference string  `json:"support_channel_preference"`

	// Claim Experience Feedback
	ClaimSatisfaction     float64 `json:"claim_satisfaction"`    // 0-100
	ClaimProcessRating    float64 `json:"claim_process_rating"`  // 0-5
	ClaimSpeedRating      float64 `json:"claim_speed_rating"`    // 0-5
	ClaimFairnessRating   float64 `json:"claim_fairness_rating"` // 0-5
	TotalClaimsExperience int     `json:"total_claims_experience"`

	// App Experience Ratings
	AppSatisfaction      float64 `json:"app_satisfaction"`       // 0-100
	AppUsabilityScore    float64 `json:"app_usability_score"`    // 0-100
	AppPerformanceRating float64 `json:"app_performance_rating"` // 0-5
	AppFeatureRating     float64 `json:"app_feature_rating"`     // 0-5
	AppStoreRating       float64 `json:"app_store_rating"`       // 0-5

	// Overall Experience Score
	ExperienceScore     float64 `json:"experience_score"`      // 0-100
	CustomerEffortScore float64 `json:"customer_effort_score"` // 0-100
	EmotionalConnection float64 `json:"emotional_connection"`  // 0-100
	BrandPerception     float64 `json:"brand_perception"`      // 0-100

	// Satisfaction Trend Analysis
	TrendPeriod           string  `json:"trend_period"`    // monthly, quarterly, yearly
	TrendDirection        string  `json:"trend_direction"` // up, down, stable
	TrendMagnitude        float64 `json:"trend_magnitude"` // percentage change
	PredictedSatisfaction float64 `json:"predicted_satisfaction"`
	ChurnRisk             float64 `json:"churn_risk"` // 0-100

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}
