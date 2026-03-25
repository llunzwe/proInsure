package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceUsagePattern tracks daily usage patterns and behaviors
type DeviceUsagePattern struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordDate time.Time `json:"record_date"`

	// Screen Time Metrics
	TotalScreenTime int `json:"total_screen_time"` // minutes per day
	ActiveHours     int `json:"active_hours"`      // hours with activity
	UnlockCount     int `json:"unlock_count"`
	LongestSession  int `json:"longest_session"` // minutes
	AverageSession  int `json:"average_session"` // minutes
	PeakUsageHour   int `json:"peak_usage_hour"` // 0-23

	// Charging Patterns
	ChargingSessions      int     `json:"charging_sessions"`
	OvernightCharging     bool    `json:"overnight_charging"`
	FastChargingCount     int     `json:"fast_charging_count"`
	WirelessChargingCount int     `json:"wireless_charging_count"`
	BatteryDrainRate      float64 `json:"battery_drain_rate"` // % per hour

	// App Usage
	AppsInstalled   int    `json:"apps_installed"`
	AppsUninstalled int    `json:"apps_uninstalled"`
	ActiveApps      string `gorm:"type:json" json:"active_apps"`       // JSON array
	HighRiskApps    string `gorm:"type:json" json:"high_risk_apps"`    // gaming, crypto, betting
	AppSessionTimes string `gorm:"type:json" json:"app_session_times"` // JSON object

	// Data Usage
	WiFiDataUsage     float64 `json:"wifi_data_usage"`     // GB
	CellularDataUsage float64 `json:"cellular_data_usage"` // GB
	TotalDataUsage    float64 `json:"total_data_usage"`    // GB
	DataRoaming       bool    `json:"data_roaming"`
	HotspotUsage      float64 `json:"hotspot_usage"` // GB

	// Usage Context
	WeekdayUsage   bool `json:"weekday_usage"`
	WeekendUsage   bool `json:"weekend_usage"`
	NightUsage     bool `json:"night_usage"`      // 10PM - 6AM
	WorkHoursUsage bool `json:"work_hours_usage"` // 9AM - 5PM

	// Risk Indicators
	AnomalyScore float64 `json:"anomaly_score"`
	RiskLevel    string  `json:"risk_level"` // low, medium, high

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceBehaviorScore calculates risk scores based on usage patterns
type DeviceBehaviorScore struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	CalculationDate time.Time `json:"calculation_date"`

	// Overall Scores
	OverallRiskScore float64 `json:"overall_risk_score"` // 0-100
	BehaviorScore    float64 `json:"behavior_score"`     // 0-100
	ConsistencyScore float64 `json:"consistency_score"`  // 0-100
	TrustScore       float64 `json:"trust_score"`        // 0-100

	// Anomaly Detection
	AnomalyDetected   bool       `json:"anomaly_detected"`
	AnomalyTypes      string     `gorm:"type:json" json:"anomaly_types"` // JSON array
	AnomalyConfidence float64    `json:"anomaly_confidence"`
	LastAnomalyDate   *time.Time `json:"last_anomaly_date"`

	// Behavior Changes
	SuddenChanges        int `json:"sudden_changes"`
	LocationChanges      int `json:"location_changes"`
	NetworkChanges       int `json:"network_changes"`
	AppPermissionChanges int `json:"app_permission_changes"`
	SecurityChanges      int `json:"security_changes"`

	// Pattern Analysis
	UsagePattern      string `json:"usage_pattern"`                       // heavy, moderate, light
	BehaviorTrend     string `json:"behavior_trend"`                      // improving, stable, degrading
	PredictedBehavior string `gorm:"type:json" json:"predicted_behavior"` // JSON object

	// Risk Factors
	HighRiskActivities    string `gorm:"type:json" json:"high_risk_activities"`   // JSON array
	RiskFactors           string `gorm:"type:json" json:"risk_factors"`           // JSON object with scores
	MitigationSuggestions string `gorm:"type:json" json:"mitigation_suggestions"` // JSON array

	// ML Model Data
	ModelVersion      string  `json:"model_version"`
	ModelConfidence   float64 `json:"model_confidence"`
	FeatureImportance string  `gorm:"type:json" json:"feature_importance"` // JSON object

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceLocationHistory tracks movement patterns and location analysis
type DeviceLocationHistory struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordDate time.Time `json:"record_date"`

	// Location Data
	CurrentLocation string `json:"current_location"` // lat,lng
	CurrentCountry  string `json:"current_country"`
	CurrentCity     string `json:"current_city"`
	LocationType    string `json:"location_type"` // home, work, travel, unknown

	// Movement Patterns
	DistanceTraveled float64 `json:"distance_traveled"` // km per day
	LocationsVisited int     `json:"locations_visited"`
	MovementPattern  string  `json:"movement_pattern"` // stationary, local, traveler
	Velocity         float64 `json:"velocity"`         // average km/h

	// Location Analysis
	HighRiskLocation    bool   `json:"high_risk_location"`
	RiskZones           string `gorm:"type:json" json:"risk_zones"` // JSON array
	InternationalTravel bool   `json:"international_travel"`
	CountriesVisited    string `gorm:"type:json" json:"countries_visited"` // JSON array

	// Stability Metrics
	HomeLocationStable bool   `json:"home_location_stable"`
	WorkLocationStable bool   `json:"work_location_stable"`
	UnusualVisit       bool   `json:"unusual_visit"`
	UnusualLocations   string `gorm:"type:json" json:"unusual_locations"` // JSON array

	// Geofencing
	GeofenceViolations int        `json:"geofence_violations"`
	GeofenceZones      string     `gorm:"type:json" json:"geofence_zones"` // JSON array
	LastViolation      *time.Time `json:"last_violation"`

	// Accuracy & Confidence
	LocationAccuracy    float64 `json:"location_accuracy"` // meters
	GPSEnabled          bool    `json:"gps_enabled"`
	WiFiLocationEnabled bool    `json:"wifi_location_enabled"`
	ConfidenceScore     float64 `json:"confidence_score"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceNetworkActivity monitors network usage and security
type DeviceNetworkActivity struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordDate time.Time `json:"record_date"`

	// Data Consumption
	DailyDataUsage   float64 `json:"daily_data_usage"`   // GB
	MonthlyDataUsage float64 `json:"monthly_data_usage"` // GB
	PeakDataHour     int     `json:"peak_data_hour"`     // 0-23
	DataTrend        string  `json:"data_trend"`         // increasing, stable, decreasing

	// Network History
	WiFiNetworks      string `gorm:"type:json" json:"wifi_networks"`    // JSON array
	TrustedNetworks   string `gorm:"type:json" json:"trusted_networks"` // JSON array
	PublicNetworks    int    `json:"public_networks"`
	UnsecuredNetworks int    `json:"unsecured_networks"`

	// Cellular Usage
	CellularCarrier string `json:"cellular_carrier"`
	SignalStrength  int    `json:"signal_strength"` // dBm
	NetworkType     string `json:"network_type"`    // 5G, 4G, 3G
	RoamingActive   bool   `json:"roaming_active"`
	RoamingDays     int    `json:"roaming_days"`

	// Hotspot & VPN
	HotspotEnabled     bool    `json:"hotspot_enabled"`
	HotspotConnections int     `json:"hotspot_connections"`
	VPNActive          bool    `json:"vpn_active"`
	VPNProvider        string  `json:"vpn_provider"`
	VPNUsageHours      float64 `json:"vpn_usage_hours"`

	// Security Assessment
	NetworkRiskScore   float64    `json:"network_risk_score"` // 0-100
	SuspiciousActivity bool       `json:"suspicious_activity"`
	SecurityIncidents  int        `json:"security_incidents"`
	LastIncident       *time.Time `json:"last_incident"`

	// Bandwidth Analysis
	AverageSpeed float64 `json:"average_speed"` // Mbps
	PeakSpeed    float64 `json:"peak_speed"`    // Mbps
	Latency      float64 `json:"latency"`       // ms
	PacketLoss   float64 `json:"packet_loss"`   // percentage

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}





// DeviceCustomerSegmentation categorizes customers for targeted strategies
type DeviceCustomerSegmentation struct {
	database.BaseModel
	DeviceID         uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID           uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	SegmentationDate time.Time `json:"segmentation_date"`

	// Primary Segmentation
	CustomerSegment string `json:"customer_segment"` // premium, standard, basic
	UsageSegment    string `json:"usage_segment"`    // power, regular, light
	RiskSegment     string `json:"risk_segment"`     // low, medium, high
	ValueSegment    string `json:"value_segment"`    // high-value, mid-value, low-value

	// Behavioral Segments
	TechSavviness    string `json:"tech_savviness"`    // expert, intermediate, beginner
	PriceSensitivity string `json:"price_sensitivity"` // high, medium, low
	BrandLoyalty     string `json:"brand_loyalty"`     // loyal, neutral, switcher
	EngagementLevel  string `json:"engagement_level"`  // active, moderate, passive

	// Premium Classification
	PremiumTier   string `json:"premium_tier"`   // platinum, gold, silver, bronze
	LoyaltyStatus string `json:"loyalty_status"` // vip, loyal, regular, new
	LoyaltyPoints int    `json:"loyalty_points"`
	TierBenefits  string `gorm:"type:json" json:"tier_benefits"` // JSON array

	// Risk Grouping
	ClaimFrequency string `json:"claim_frequency"` // none, low, medium, high
	PaymentHistory string `json:"payment_history"` // excellent, good, fair, poor
	FraudRisk      string `json:"fraud_risk"`      // minimal, low, medium, high
	ChurnRisk      string `json:"churn_risk"`      // minimal, low, medium, high

	// Value Metrics
	LifetimeValue      float64 `json:"lifetime_value"`
	AverageRevenue     float64 `json:"average_revenue"`
	ProfitContribution float64 `json:"profit_contribution"`
	CostToServe        float64 `json:"cost_to_serve"`

	// Targeting & Personalization
	MarketingPreferences   string  `gorm:"type:json" json:"marketing_preferences"`   // JSON object
	ProductRecommendations string  `gorm:"type:json" json:"product_recommendations"` // JSON array
	OptimalChannels        string  `gorm:"type:json" json:"optimal_channels"`        // JSON array
	PersonalizationScore   float64 `json:"personalization_score"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// DeviceProfitability tracks device-level P&L and financial metrics
type DeviceProfitability struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	CalculationDate time.Time `json:"calculation_date"`
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`

	// Revenue Streams
	PremiumRevenue    float64 `json:"premium_revenue"`
	DeductibleRevenue float64 `json:"deductible_revenue"`
	ServiceRevenue    float64 `json:"service_revenue"`
	AccessoryRevenue  float64 `json:"accessory_revenue"`
	OtherRevenue      float64 `json:"other_revenue"`
	TotalRevenue      float64 `json:"total_revenue"`

	// Cost Components
	ClaimCosts       float64 `json:"claim_costs"`
	ServiceCosts     float64 `json:"service_costs"`
	AdminCosts       float64 `json:"admin_costs"`
	AcquisitionCosts float64 `json:"acquisition_costs"`
	RetentionCosts   float64 `json:"retention_costs"`
	FraudCosts       float64 `json:"fraud_costs"`
	TotalCosts       float64 `json:"total_costs"`

	// Profitability Metrics
	GrossProfit float64 `json:"gross_profit"`
	GrossMargin float64 `json:"gross_margin"` // percentage
	NetProfit   float64 `json:"net_profit"`
	NetMargin   float64 `json:"net_margin"` // percentage
	EBITDA      float64 `json:"ebitda"`

	// Lifetime Value
	CustomerLifetimeValue float64    `json:"customer_lifetime_value"`
	PaybackPeriod         int        `json:"payback_period"` // months
	BreakEvenDate         *time.Time `json:"break_even_date"`
	ProfitableMonths      int        `json:"profitable_months"`

	// Performance Scoring
	ProfitabilityScore float64 `json:"profitability_score"` // 0-100
	EfficiencyRatio    float64 `json:"efficiency_ratio"`
	LossRatio          float64 `json:"loss_ratio"`
	ExpenseRatio       float64 `json:"expense_ratio"`
	CombinedRatio      float64 `json:"combined_ratio"`

	// Benchmarking
	IndustryAverage float64 `json:"industry_average"`
	PerformanceRank int     `json:"performance_rank"`
	TopPercentile   float64 `json:"top_percentile"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}


// DeviceMarketAnalysis provides market positioning and competitive insights
type DeviceMarketAnalysis struct {
	database.BaseModel
	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ModelName    string    `json:"model_name"`
	Brand        string    `json:"brand"`
	AnalysisDate time.Time `json:"analysis_date"`

	// Market Positioning
	MarketPosition string  `json:"market_position"`                  // leader, challenger, follower, niche
	MarketSegment  string  `json:"market_segment"`                   // premium, mid-range, budget
	TargetAudience string  `gorm:"type:json" json:"target_audience"` // JSON object
	MarketShare    float64 `json:"market_share"`                     // percentage

	// Competitive Analysis
	MainCompetitors       string `gorm:"type:json" json:"main_competitors"`       // JSON array
	CompetitiveAdvantages string `gorm:"type:json" json:"competitive_advantages"` // JSON array
	CompetitiveWeaknesses string `gorm:"type:json" json:"competitive_weaknesses"` // JSON array
	PriceCompetitiveness  string `json:"price_competitiveness"`                   // above, at, below market

	// Pricing Analysis
	AverageMarketPrice float64 `json:"average_market_price"`
	PriceElasticity    float64 `json:"price_elasticity"`
	OptimalPricePoint  float64 `json:"optimal_price_point"`
	PricePositioning   string  `json:"price_positioning"` // premium, competitive, value

	// Market Trends
	MarketGrowthRate float64 `json:"market_growth_rate"`
	TrendDirection   string  `json:"trend_direction"`                    // growing, stable, declining
	EmergingTrends   string  `gorm:"type:json" json:"emerging_trends"`   // JSON array
	ThreatAssessment string  `gorm:"type:json" json:"threat_assessment"` // JSON object

	// Brand Perception
	BrandStrength     float64 `json:"brand_strength"`     // 0-100
	BrandRecognition  float64 `json:"brand_recognition"`  // 0-100
	CustomerSentiment string  `json:"customer_sentiment"` // positive, neutral, negative
	ReputationScore   float64 `json:"reputation_score"`   // 0-100

	// Feature Analysis
	KeyFeatures       string  `gorm:"type:json" json:"key_features"`       // JSON array
	FeatureImportance string  `gorm:"type:json" json:"feature_importance"` // JSON object
	MissingFeatures   string  `gorm:"type:json" json:"missing_features"`   // JSON array
	InnovationScore   float64 `json:"innovation_score"`                    // 0-100

	// Market Opportunities
	GrowthOpportunities string  `gorm:"type:json" json:"growth_opportunities"` // JSON array
	UnderservedSegments string  `gorm:"type:json" json:"underserved_segments"` // JSON array
	ExpansionPotential  float64 `json:"expansion_potential"`                   // 0-100

	// Recommendations
	StrategicRecommendations string `gorm:"type:json" json:"strategic_recommendations"` // JSON array
	PricingStrategy          string `json:"pricing_strategy"`
	MarketingFocus           string `gorm:"type:json" json:"marketing_focus"` // JSON object

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceUsagePattern
func (dup *DeviceUsagePattern) IsHeavyUser() bool {
	return dup.TotalScreenTime > 360 || dup.ActiveHours > 12 // >6 hours screen time or >12 active hours
}

func (dup *DeviceUsagePattern) HasRiskyApps() bool {
	return dup.HighRiskApps != "" && dup.HighRiskApps != "[]"
}

func (dup *DeviceUsagePattern) CalculateRiskLevel() string {
	if dup.AnomalyScore > 70 || dup.HasRiskyApps() {
		dup.RiskLevel = "high"
	} else if dup.AnomalyScore > 40 {
		dup.RiskLevel = "medium"
	} else {
		dup.RiskLevel = "low"
	}
	return dup.RiskLevel
}

// Methods for DeviceBehaviorScore
func (dbs *DeviceBehaviorScore) IsHighRisk() bool {
	return dbs.OverallRiskScore > 70 || dbs.AnomalyDetected
}

func (dbs *DeviceBehaviorScore) NeedsIntervention() bool {
	return dbs.SuddenChanges > 5 || dbs.BehaviorTrend == "degrading"
}

func (dbs *DeviceBehaviorScore) CalculateTrustScore() float64 {
	dbs.TrustScore = (100 - dbs.OverallRiskScore) * (dbs.ConsistencyScore / 100)
	if dbs.AnomalyDetected {
		dbs.TrustScore *= 0.7
	}
	return dbs.TrustScore
}

// Methods for DeviceLocationHistory
func (dlh *DeviceLocationHistory) IsInternationalTraveler() bool {
	return dlh.InternationalTravel || dlh.MovementPattern == "traveler"
}

func (dlh *DeviceLocationHistory) HasGeofenceViolation() bool {
	return dlh.GeofenceViolations > 0
}

func (dlh *DeviceLocationHistory) IsHighRiskLocation() bool {
	return dlh.HighRiskLocation || dlh.UnusualVisit
}

// Methods for DeviceNetworkActivity
func (dna *DeviceNetworkActivity) IsHighDataUser() bool {
	return dna.MonthlyDataUsage > 50 // >50GB per month
}

func (dna *DeviceNetworkActivity) HasSecurityRisk() bool {
	return dna.NetworkRiskScore > 60 || dna.UnsecuredNetworks > 3 || dna.SuspiciousActivity
}

func (dna *DeviceNetworkActivity) IsRoaming() bool {
	return dna.RoamingActive || dna.RoamingDays > 0
}

// DeviceFailurePrediction methods removed - struct moved to device_predictive.go



// Methods for DeviceMaintenanceSchedule



// Methods for DeviceUpgradeRecommendation



// Methods for DeviceValueForecast



// Methods for DeviceCustomerSegmentation



// Methods for DeviceProfitability
func (dp *DeviceProfitability) CalculateProfitMargin() float64 {
	if dp.TotalRevenue > 0 {
		dp.NetMargin = (dp.NetProfit / dp.TotalRevenue) * 100
	}
	return dp.NetMargin
}

func (dp *DeviceProfitability) IsProfitable() bool {
	return dp.NetProfit > 0 && dp.ProfitabilityScore > 50
}

func (dp *DeviceProfitability) CalculateLossRatio() float64 {
	if dp.PremiumRevenue > 0 {
		dp.LossRatio = (dp.ClaimCosts / dp.PremiumRevenue) * 100
	}
	return dp.LossRatio
}

// Methods for DeviceChurnPrediction



// Methods for DeviceMarketAnalysis


