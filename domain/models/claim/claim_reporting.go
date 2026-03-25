package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ClaimReporting represents reporting and analytics dashboards for claims
type ClaimReporting struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ClaimID uuid.UUID `gorm:"type:uuid;not null;index" json:"claim_id"`

	// Report Generation
	ReportNumber  string    `gorm:"type:varchar(100);unique" json:"report_number"`
	ReportType    string    `gorm:"type:varchar(50)" json:"report_type"`
	ReportPeriod  string    `gorm:"type:varchar(50)" json:"report_period"`
	GeneratedDate time.Time `json:"generated_date"`
	GeneratedBy   uuid.UUID `gorm:"type:uuid" json:"generated_by"`

	// Key Performance Indicators
	AverageCycleTime     int     `json:"average_cycle_time_days"`
	FirstCallResolution  float64 `gorm:"type:decimal(5,2)" json:"first_call_resolution"`
	CustomerSatisfaction float64 `gorm:"type:decimal(5,2)" json:"customer_satisfaction"`
	NPSScore             int     `json:"nps_score"`
	ClaimAccuracy        float64 `gorm:"type:decimal(5,2)" json:"claim_accuracy"`
	SLACompliance        float64 `gorm:"type:decimal(5,2)" json:"sla_compliance"`

	// Volume Metrics
	TotalClaims    int `json:"total_claims"`
	OpenClaims     int `json:"open_claims"`
	ClosedClaims   int `json:"closed_claims"`
	ApprovedClaims int `json:"approved_claims"`
	DeniedClaims   int `json:"denied_claims"`
	AppealedClaims int `json:"appealed_claims"`

	// Financial Metrics
	TotalClaimedAmount  float64 `gorm:"type:decimal(15,2)" json:"total_claimed_amount"`
	TotalApprovedAmount float64 `gorm:"type:decimal(15,2)" json:"total_approved_amount"`
	TotalPaidAmount     float64 `gorm:"type:decimal(15,2)" json:"total_paid_amount"`
	AverageClaimAmount  float64 `gorm:"type:decimal(15,2)" json:"average_claim_amount"`
	LossRatio           float64 `gorm:"type:decimal(5,2)" json:"loss_ratio"`
	ExpenseRatio        float64 `gorm:"type:decimal(5,2)" json:"expense_ratio"`
	CombinedRatio       float64 `gorm:"type:decimal(5,2)" json:"combined_ratio"`

	// Fraud Metrics
	FraudDetectionRate float64 `gorm:"type:decimal(5,2)" json:"fraud_detection_rate"`
	FalsePositiveRate  float64 `gorm:"type:decimal(5,2)" json:"false_positive_rate"`
	FraudSavings       float64 `gorm:"type:decimal(15,2)" json:"fraud_savings"`
	InvestigationRate  float64 `gorm:"type:decimal(5,2)" json:"investigation_rate"`

	// Operational Metrics
	AutoApprovalRate float64 `gorm:"type:decimal(5,2)" json:"auto_approval_rate"`
	ManualReviewRate float64 `gorm:"type:decimal(5,2)" json:"manual_review_rate"`
	ReopenRate       float64 `gorm:"type:decimal(5,2)" json:"reopen_rate"`
	EscalationRate   float64 `gorm:"type:decimal(5,2)" json:"escalation_rate"`

	// Productivity Metrics
	ClaimsPerAdjuster float64 `gorm:"type:decimal(10,2)" json:"claims_per_adjuster"`
	AverageHandleTime int     `json:"average_handle_time_minutes"`
	ProductivityScore float64 `gorm:"type:decimal(5,2)" json:"productivity_score"`
	UtilizationRate   float64 `gorm:"type:decimal(5,2)" json:"utilization_rate"`

	// Quality Metrics
	QualityScore    float64 `gorm:"type:decimal(5,2)" json:"quality_score"`
	ErrorRate       float64 `gorm:"type:decimal(5,2)" json:"error_rate"`
	ComplianceScore float64 `gorm:"type:decimal(5,2)" json:"compliance_score"`
	AuditScore      float64 `gorm:"type:decimal(5,2)" json:"audit_score"`

	// Trend Analysis
	ClaimTrend     string         `gorm:"type:varchar(20)" json:"claim_trend"` // increasing, decreasing, stable
	CostTrend      string         `gorm:"type:varchar(20)" json:"cost_trend"`
	FrequencyTrend string         `gorm:"type:varchar(20)" json:"frequency_trend"`
	SeverityTrend  string         `gorm:"type:varchar(20)" json:"severity_trend"`
	TrendData      datatypes.JSON `gorm:"type:json" json:"trend_data"`

	// Comparative Analysis
	IndustryBenchmark    datatypes.JSON `gorm:"type:json" json:"industry_benchmark"`
	RegionalComparison   datatypes.JSON `gorm:"type:json" json:"regional_comparison"`
	HistoricalComparison datatypes.JSON `gorm:"type:json" json:"historical_comparison"`
	PeerGroupRanking     int            `json:"peer_group_ranking"`

	// Dashboard Configuration
	DashboardLayout datatypes.JSON `gorm:"type:json" json:"dashboard_layout"`
	WidgetsConfig   datatypes.JSON `gorm:"type:json" json:"widgets_config"`
	FiltersApplied  datatypes.JSON `gorm:"type:json" json:"filters_applied"`
	DateRange       datatypes.JSON `gorm:"type:json" json:"date_range"`

	// Visualizations
	ChartsGenerated int            `json:"charts_generated"`
	ChartConfigs    datatypes.JSON `gorm:"type:json" json:"chart_configs"`
	ExportFormats   datatypes.JSON `gorm:"type:json" json:"export_formats"`

	// Alerts & Notifications
	AlertsConfigured  int            `json:"alerts_configured"`
	AlertThresholds   datatypes.JSON `gorm:"type:json" json:"alert_thresholds"`
	AlertsTriggered   int            `json:"alerts_triggered"`
	NotificationsSent int            `json:"notifications_sent"`

	// Status
	ReportStatus     string     `gorm:"type:varchar(50)" json:"report_status"`
	LastRefreshed    time.Time  `json:"last_refreshed"`
	NextScheduledRun *time.Time `json:"next_scheduled_run,omitempty"`
	CreatedAt        time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ClaimDashboard represents real-time claim dashboards
type ClaimDashboard struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Dashboard Identity
	DashboardName string         `gorm:"type:varchar(255)" json:"dashboard_name"`
	DashboardType string         `gorm:"type:varchar(50)" json:"dashboard_type"` // executive, operational, analytical
	OwnerID       uuid.UUID      `gorm:"type:uuid" json:"owner_id"`
	SharedWith    datatypes.JSON `gorm:"type:json" json:"shared_with"` // []UserID

	// Real-time Metrics
	LiveClaimsCount  int `json:"live_claims_count"`
	ProcessingQueue  int `json:"processing_queue"`
	PendingApprovals int `json:"pending_approvals"`
	CriticalAlerts   int `json:"critical_alerts"`

	// Today's Performance
	TodayClaimsReceived  int     `json:"today_claims_received"`
	TodayClaimsProcessed int     `json:"today_claims_processed"`
	TodayClaimsApproved  int     `json:"today_claims_approved"`
	TodayClaimsDenied    int     `json:"today_claims_denied"`
	TodayAmount          float64 `gorm:"type:decimal(15,2)" json:"today_amount"`

	// Week Performance
	WeekClaimsTotal   int     `json:"week_claims_total"`
	WeekApprovalRate  float64 `gorm:"type:decimal(5,2)" json:"week_approval_rate"`
	WeekAverageTime   int     `json:"week_average_time_hours"`
	WeekSLACompliance float64 `gorm:"type:decimal(5,2)" json:"week_sla_compliance"`

	// Month Performance
	MonthClaimsTotal   int     `json:"month_claims_total"`
	MonthTotalAmount   float64 `gorm:"type:decimal(15,2)" json:"month_total_amount"`
	MonthLossRatio     float64 `gorm:"type:decimal(5,2)" json:"month_loss_ratio"`
	MonthFraudDetected int     `json:"month_fraud_detected"`

	// Top Metrics
	TopClaimTypes    datatypes.JSON `gorm:"type:json" json:"top_claim_types"`
	TopDenialReasons datatypes.JSON `gorm:"type:json" json:"top_denial_reasons"`
	TopAdjusters     datatypes.JSON `gorm:"type:json" json:"top_adjusters"`
	TopRegions       datatypes.JSON `gorm:"type:json" json:"top_regions"`

	// Heat Maps
	GeographicHeatmap datatypes.JSON `gorm:"type:json" json:"geographic_heatmap"`
	TimeHeatmap       datatypes.JSON `gorm:"type:json" json:"time_heatmap"`
	CategoryHeatmap   datatypes.JSON `gorm:"type:json" json:"category_heatmap"`

	// Predictive Analytics
	ForecastedVolume  int     `json:"forecasted_volume"`
	ForecastedAmount  float64 `gorm:"type:decimal(15,2)" json:"forecasted_amount"`
	RiskScore         float64 `gorm:"type:decimal(5,2)" json:"risk_score"`
	AnomaliesDetected int     `json:"anomalies_detected"`

	// Drill-down Data
	DrilldownEnabled  bool           `gorm:"default:true" json:"drilldown_enabled"`
	DetailedBreakdown datatypes.JSON `gorm:"type:json" json:"detailed_breakdown"`
	FilterOptions     datatypes.JSON `gorm:"type:json" json:"filter_options"`

	// Refresh Settings
	AutoRefresh     bool           `gorm:"default:true" json:"auto_refresh"`
	RefreshInterval int            `json:"refresh_interval_seconds"`
	LastRefresh     time.Time      `json:"last_refresh"`
	DataSources     datatypes.JSON `gorm:"type:json" json:"data_sources"`

	// Status
	DashboardStatus string    `gorm:"type:varchar(50)" json:"dashboard_status"`
	IsDefault       bool      `gorm:"default:false" json:"is_default"`
	IsPublic        bool      `gorm:"default:false" json:"is_public"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ClaimReportSchedule represents scheduled report generation
type ClaimReportSchedule struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Schedule Identity
	ScheduleName string         `gorm:"type:varchar(255)" json:"schedule_name"`
	ReportType   string         `gorm:"type:varchar(50)" json:"report_type"`
	Recipients   datatypes.JSON `gorm:"type:json" json:"recipients"` // []EmailAddress

	// Schedule Configuration
	Frequency  string `gorm:"type:varchar(50)" json:"frequency"`   // daily, weekly, monthly, quarterly
	DayOfWeek  int    `json:"day_of_week"`                         // 0-6 for weekly
	DayOfMonth int    `json:"day_of_month"`                        // 1-31 for monthly
	TimeOfDay  string `gorm:"type:varchar(10)" json:"time_of_day"` // HH:MM
	Timezone   string `gorm:"type:varchar(50)" json:"timezone"`

	// Report Parameters
	ReportParameters datatypes.JSON `gorm:"type:json" json:"report_parameters"`
	DataFilters      datatypes.JSON `gorm:"type:json" json:"data_filters"`
	IncludeCharts    bool           `gorm:"default:true" json:"include_charts"`
	IncludeRawData   bool           `gorm:"default:false" json:"include_raw_data"`

	// Delivery Options
	DeliveryMethod string `gorm:"type:varchar(50)" json:"delivery_method"` // email, sftp, api
	FileFormat     string `gorm:"type:varchar(20)" json:"file_format"`     // pdf, excel, csv, json
	Compression    bool   `gorm:"default:false" json:"compression"`
	Encryption     bool   `gorm:"default:false" json:"encryption"`

	// Execution History
	LastRunTime   *time.Time `json:"last_run_time,omitempty"`
	LastRunStatus string     `gorm:"type:varchar(50)" json:"last_run_status"`
	NextRunTime   time.Time  `json:"next_run_time"`
	RunCount      int        `json:"run_count"`
	FailureCount  int        `json:"failure_count"`

	// Status
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsHealthy checks if KPIs are within acceptable ranges
func (cr *ClaimReporting) IsHealthy() bool {
	return cr.SLACompliance >= 95 &&
		cr.CustomerSatisfaction >= 80 &&
		cr.LossRatio < 70 &&
		cr.ErrorRate < 5
}

// NeedsAttention checks if any metric needs attention
func (cr *ClaimReporting) NeedsAttention() bool {
	return cr.SLACompliance < 90 ||
		cr.CustomerSatisfaction < 70 ||
		cr.FraudDetectionRate < 80 ||
		cr.ErrorRate > 10
}

// GetCombinedRatio calculates the combined ratio
func (cr *ClaimReporting) GetCombinedRatio() float64 {
	return cr.LossRatio + cr.ExpenseRatio
}

// IsProfitable checks if claims operation is profitable
func (cr *ClaimReporting) IsProfitable() bool {
	return cr.GetCombinedRatio() < 100
}

// HasAnomalies checks for anomalies in the dashboard
func (cd *ClaimDashboard) HasAnomalies() bool {
	return cd.AnomaliesDetected > 0 || cd.CriticalAlerts > 0
}

// NeedsRefresh checks if dashboard needs refresh
func (cd *ClaimDashboard) NeedsRefresh() bool {
	if !cd.AutoRefresh {
		return false
	}
	timeSinceRefresh := time.Since(cd.LastRefresh).Seconds()
	return timeSinceRefresh >= float64(cd.RefreshInterval)
}

// GetPerformanceTrend returns performance trend
func (cd *ClaimDashboard) GetPerformanceTrend() string {
	if cd.WeekSLACompliance > 95 && cd.WeekApprovalRate > 80 {
		return "excellent"
	} else if cd.WeekSLACompliance > 90 && cd.WeekApprovalRate > 70 {
		return "good"
	} else if cd.WeekSLACompliance > 85 && cd.WeekApprovalRate > 60 {
		return "fair"
	}
	return "needs_improvement"
}

// ShouldRun checks if scheduled report should run
func (crs *ClaimReportSchedule) ShouldRun() bool {
	if !crs.IsActive {
		return false
	}
	return time.Now().After(crs.NextRunTime)
}

// GetNextRunTime calculates the next run time
func (crs *ClaimReportSchedule) GetNextRunTime() time.Time {
	now := time.Now()

	switch crs.Frequency {
	case "daily":
		return now.AddDate(0, 0, 1)
	case "weekly":
		return now.AddDate(0, 0, 7)
	case "monthly":
		return now.AddDate(0, 1, 0)
	case "quarterly":
		return now.AddDate(0, 3, 0)
	default:
		return now.AddDate(0, 0, 1) // Default to daily
	}
}

// IsOverdue checks if scheduled report is overdue
func (crs *ClaimReportSchedule) IsOverdue() bool {
	return crs.IsActive && time.Now().After(crs.NextRunTime)
}
