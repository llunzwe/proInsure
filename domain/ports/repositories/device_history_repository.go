package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DeviceHistoryRepository defines the interface for device history operations
type DeviceHistoryRepository interface {
	// === History Recording ===

	// RecordDeviceEvent records a device event
	RecordDeviceEvent(ctx context.Context, event *DeviceEvent) error

	// RecordOwnershipChange records ownership transfer
	RecordOwnershipChange(ctx context.Context, change *OwnershipChange) error

	// RecordValueChange records value change
	RecordValueChange(ctx context.Context, change *ValueChange) error

	// RecordConditionChange records condition change
	RecordConditionChange(ctx context.Context, change *ConditionChange) error

	// RecordLocationChange records location change
	RecordLocationChange(ctx context.Context, change *LocationChange) error

	// RecordRepairHistory records repair event
	RecordRepairHistory(ctx context.Context, repair *RepairHistory) error

	// RecordClaimHistory records claim event
	RecordClaimHistory(ctx context.Context, claim *ClaimHistory) error

	// RecordMaintenanceEvent records maintenance activity
	RecordMaintenanceEvent(ctx context.Context, maintenance *MaintenanceHistory) error

	// === History Retrieval ===

	// GetDeviceHistory gets complete device history
	GetDeviceHistory(ctx context.Context, deviceID uuid.UUID) ([]*DeviceEvent, error)

	// GetDeviceHistoryByType gets history filtered by event type
	GetDeviceHistoryByType(ctx context.Context, deviceID uuid.UUID, eventType string) ([]*DeviceEvent, error)

	// GetDeviceHistoryInPeriod gets history within time period
	GetDeviceHistoryInPeriod(ctx context.Context, deviceID uuid.UUID, startDate, endDate time.Time) ([]*DeviceEvent, error)

	// GetOwnershipHistory gets ownership transfer history
	GetOwnershipHistory(ctx context.Context, deviceID uuid.UUID) ([]*OwnershipChange, error)

	// GetValueHistory gets value change history
	GetValueHistory(ctx context.Context, deviceID uuid.UUID) ([]*ValueChange, error)

	// GetConditionHistory gets condition change history
	GetConditionHistory(ctx context.Context, deviceID uuid.UUID) ([]*ConditionChange, error)

	// GetLocationHistory gets location history
	GetLocationHistory(ctx context.Context, deviceID uuid.UUID) ([]*LocationChange, error)

	// GetRepairHistory gets repair history
	GetRepairHistory(ctx context.Context, deviceID uuid.UUID) ([]*RepairHistory, error)

	// GetClaimHistory gets claim history
	GetClaimHistory(ctx context.Context, deviceID uuid.UUID) ([]*ClaimHistory, error)

	// GetMaintenanceHistory gets maintenance history
	GetMaintenanceHistory(ctx context.Context, deviceID uuid.UUID) ([]*MaintenanceHistory, error)

	// === Usage Analytics ===

	// GetUsageStatistics gets device usage statistics
	GetUsageStatistics(ctx context.Context, deviceID uuid.UUID, period int) (*UsageStatistics, error)

	// GetUsagePattern analyzes usage patterns
	GetUsagePattern(ctx context.Context, deviceID uuid.UUID) (*UsagePattern, error)

	// GetBatteryHistory gets battery health history
	GetBatteryHistory(ctx context.Context, deviceID uuid.UUID) ([]*BatteryHealth, error)

	// GetPerformanceHistory gets performance metrics history
	GetPerformanceHistory(ctx context.Context, deviceID uuid.UUID) ([]*PerformanceMetric, error)

	// GetNetworkUsageHistory gets network usage history
	GetNetworkUsageHistory(ctx context.Context, deviceID uuid.UUID) ([]*NetworkUsage, error)

	// === Lifecycle Analytics ===

	// GetDeviceLifecycle gets device lifecycle analysis
	GetDeviceLifecycle(ctx context.Context, deviceID uuid.UUID) (*DeviceLifecycleAnalysis, error)

	// GetDepreciationHistory tracks value depreciation over time
	GetDepreciationHistory(ctx context.Context, deviceID uuid.UUID) ([]*DepreciationRecord, error)

	// GetWarrantyHistory gets warranty status changes
	GetWarrantyHistory(ctx context.Context, deviceID uuid.UUID) ([]*WarrantyEvent, error)

	// GetInsuranceHistory gets insurance coverage history
	GetInsuranceHistory(ctx context.Context, deviceID uuid.UUID) ([]*InsuranceEvent, error)

	// === Risk History ===

	// GetRiskScoreHistory gets risk score changes over time
	GetRiskScoreHistory(ctx context.Context, deviceID uuid.UUID) ([]*RiskScoreChange, error)

	// GetIncidentHistory gets security/damage incidents
	GetIncidentHistory(ctx context.Context, deviceID uuid.UUID) ([]*IncidentRecord, error)

	// GetComplianceHistory gets compliance status changes
	GetComplianceHistory(ctx context.Context, deviceID uuid.UUID) ([]*ComplianceEvent, error)

	// === Aggregated Analytics ===

	// GetDeviceHealthScore calculates overall device health
	GetDeviceHealthScore(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// GetReliabilityScore calculates device reliability
	GetReliabilityScore(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// GetTotalCostOfOwnership calculates TCO
	GetTotalCostOfOwnership(ctx context.Context, deviceID uuid.UUID) (decimal.Decimal, error)

	// PredictFailure predicts potential failure
	PredictFailure(ctx context.Context, deviceID uuid.UUID) (*FailurePrediction, error)

	// GetUpgradeRecommendation generates upgrade recommendation
	GetUpgradeRecommendation(ctx context.Context, deviceID uuid.UUID) (*UpgradeRecommendation, error)

	// === Comparative Analytics ===

	// CompareWithSimilarDevices compares with similar devices
	CompareWithSimilarDevices(ctx context.Context, deviceID uuid.UUID) (*ComparativeAnalysis, error)

	// GetCategoryBenchmarks gets category benchmarks
	GetCategoryBenchmarks(ctx context.Context, category string) (*CategoryBenchmarks, error)

	// === Data Management ===

	// ArchiveHistory archives old history data
	ArchiveHistory(ctx context.Context, deviceID uuid.UUID, beforeDate time.Time) error

	// PurgeHistory permanently deletes history
	PurgeHistory(ctx context.Context, deviceID uuid.UUID, beforeDate time.Time) error

	// ExportHistory exports device history
	ExportHistory(ctx context.Context, deviceID uuid.UUID, format string) ([]byte, error)
}

// === History Models ===

// DeviceEvent represents any device event
type DeviceEvent struct {
	ID          uuid.UUID              `json:"id"`
	DeviceID    uuid.UUID              `json:"device_id"`
	EventType   string                 `json:"event_type"`
	EventData   map[string]interface{} `json:"event_data"`
	Severity    string                 `json:"severity"`
	UserID      uuid.UUID              `json:"user_id"`
	Source      string                 `json:"source"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
}

// OwnershipChange represents ownership transfer
type OwnershipChange struct {
	ID             uuid.UUID       `json:"id"`
	DeviceID       uuid.UUID       `json:"device_id"`
	PreviousOwner  uuid.UUID       `json:"previous_owner_id"`
	NewOwner       uuid.UUID       `json:"new_owner_id"`
	TransferType   string          `json:"transfer_type"` // sale, gift, corporate_assignment
	TransferPrice  decimal.Decimal `json:"transfer_price"`
	TransferDate   time.Time       `json:"transfer_date"`
	VerificationID string          `json:"verification_id"`
	Notes          string          `json:"notes"`
}

// ValueChange represents device value change
type ValueChange struct {
	ID            uuid.UUID       `json:"id"`
	DeviceID      uuid.UUID       `json:"device_id"`
	PreviousValue decimal.Decimal `json:"previous_value"`
	NewValue      decimal.Decimal `json:"new_value"`
	ChangeReason  string          `json:"change_reason"`
	ChangePercent float64         `json:"change_percent"`
	Source        string          `json:"source"`
	RecordedAt    time.Time       `json:"recorded_at"`
}

// ConditionChange represents condition change
type ConditionChange struct {
	ID                uuid.UUID `json:"id"`
	DeviceID          uuid.UUID `json:"device_id"`
	PreviousGrade     string    `json:"previous_grade"`
	NewGrade          string    `json:"new_grade"`
	PreviousCondition string    `json:"previous_condition"`
	NewCondition      string    `json:"new_condition"`
	ChangeReason      string    `json:"change_reason"`
	InspectionID      string    `json:"inspection_id"`
	InspectorID       uuid.UUID `json:"inspector_id"`
	Photos            []string  `json:"photos"`
	RecordedAt        time.Time `json:"recorded_at"`
}

// LocationChange represents location change
type LocationChange struct {
	ID               uuid.UUID `json:"id"`
	DeviceID         uuid.UUID `json:"device_id"`
	PreviousLocation string    `json:"previous_location"`
	NewLocation      string    `json:"new_location"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	Country          string    `json:"country"`
	Region           string    `json:"region"`
	City             string    `json:"city"`
	RecordedAt       time.Time `json:"recorded_at"`
}

// RepairHistory represents repair event
type RepairHistory struct {
	ID            uuid.UUID       `json:"id"`
	DeviceID      uuid.UUID       `json:"device_id"`
	RepairID      uuid.UUID       `json:"repair_id"`
	RepairType    string          `json:"repair_type"`
	Components    []string        `json:"components_replaced"`
	RepairCost    decimal.Decimal `json:"repair_cost"`
	TechnicianID  uuid.UUID       `json:"technician_id"`
	RepairShopID  uuid.UUID       `json:"repair_shop_id"`
	WarrantyClaim bool            `json:"warranty_claim"`
	RepairDate    time.Time       `json:"repair_date"`
	CompletedDate time.Time       `json:"completed_date"`
	Quality       string          `json:"quality_rating"`
}

// ClaimHistory represents claim event
type ClaimHistory struct {
	ID             uuid.UUID       `json:"id"`
	DeviceID       uuid.UUID       `json:"device_id"`
	ClaimID        uuid.UUID       `json:"claim_id"`
	ClaimType      string          `json:"claim_type"`
	ClaimAmount    decimal.Decimal `json:"claim_amount"`
	ApprovedAmount decimal.Decimal `json:"approved_amount"`
	Status         string          `json:"status"`
	ClaimDate      time.Time       `json:"claim_date"`
	ResolvedDate   *time.Time      `json:"resolved_date"`
	DenialReason   string          `json:"denial_reason"`
}

// MaintenanceHistory represents maintenance activity
type MaintenanceHistory struct {
	ID               uuid.UUID       `json:"id"`
	DeviceID         uuid.UUID       `json:"device_id"`
	MaintenanceType  string          `json:"maintenance_type"`
	Description      string          `json:"description"`
	PerformedBy      uuid.UUID       `json:"performed_by"`
	Cost             decimal.Decimal `json:"cost"`
	NextScheduled    *time.Time      `json:"next_scheduled"`
	PreventiveAction bool            `json:"preventive_action"`
	PerformedAt      time.Time       `json:"performed_at"`
}

// UsageStatistics represents usage stats
type UsageStatistics struct {
	DeviceID         uuid.UUID          `json:"device_id"`
	PeriodDays       int                `json:"period_days"`
	TotalActiveHours float64            `json:"total_active_hours"`
	DailyAvgHours    float64            `json:"daily_avg_hours"`
	ChargesCycles    int                `json:"charge_cycles"`
	DataUsageGB      float64            `json:"data_usage_gb"`
	AppUsageHours    map[string]float64 `json:"app_usage_hours"`
	PeakUsageTime    string             `json:"peak_usage_time"`
	IdleTimePercent  float64            `json:"idle_time_percent"`
}

// UsagePattern represents usage patterns
type UsagePattern struct {
	DeviceID      uuid.UUID          `json:"device_id"`
	PatternType   string             `json:"pattern_type"` // heavy, moderate, light
	WeekdayUsage  float64            `json:"weekday_avg_hours"`
	WeekendUsage  float64            `json:"weekend_avg_hours"`
	PeakHours     []int              `json:"peak_hours"`
	TopActivities []string           `json:"top_activities"`
	UsageScore    float64            `json:"usage_score"`
	Anomalies     []string           `json:"anomalies"`
	Predictions   map[string]float64 `json:"predictions"`
}

// BatteryHealth represents battery health data
type BatteryHealth struct {
	ID              uuid.UUID `json:"id"`
	DeviceID        uuid.UUID `json:"device_id"`
	HealthPercent   float64   `json:"health_percent"`
	CycleCount      int       `json:"cycle_count"`
	MaxCapacityMAH  int       `json:"max_capacity_mah"`
	CurrentCapacity int       `json:"current_capacity_mah"`
	Temperature     float64   `json:"temperature_celsius"`
	RecordedAt      time.Time `json:"recorded_at"`
}

// PerformanceMetric represents performance data
type PerformanceMetric struct {
	ID             uuid.UUID `json:"id"`
	DeviceID       uuid.UUID `json:"device_id"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    float64   `json:"memory_usage"`
	StorageUsed    float64   `json:"storage_used_gb"`
	AppCrashes     int       `json:"app_crashes"`
	SystemRestarts int       `json:"system_restarts"`
	ResponseTime   float64   `json:"response_time_ms"`
	BenchmarkScore int       `json:"benchmark_score"`
	RecordedAt     time.Time `json:"recorded_at"`
}

// NetworkUsage represents network usage data
type NetworkUsage struct {
	ID             uuid.UUID       `json:"id"`
	DeviceID       uuid.UUID       `json:"device_id"`
	WiFiDataGB     float64         `json:"wifi_data_gb"`
	CellularDataGB float64         `json:"cellular_data_gb"`
	RoamingDataGB  float64         `json:"roaming_data_gb"`
	HotspotDataGB  float64         `json:"hotspot_data_gb"`
	DataCostUSD    decimal.Decimal `json:"data_cost_usd"`
	PeriodStart    time.Time       `json:"period_start"`
	PeriodEnd      time.Time       `json:"period_end"`
}

// DeviceLifecycleAnalysis represents lifecycle analysis
type DeviceLifecycleAnalysis struct {
	DeviceID            uuid.UUID       `json:"device_id"`
	AgeMonths           int             `json:"age_months"`
	ExpectedLifeMonths  int             `json:"expected_life_months"`
	RemainingLifeMonths int             `json:"remaining_life_months"`
	LifecycleStage      string          `json:"lifecycle_stage"` // new, mature, aging, end_of_life
	TotalRepairs        int             `json:"total_repairs"`
	TotalClaims         int             `json:"total_claims"`
	TotalCostOwnership  decimal.Decimal `json:"total_cost_ownership"`
	ReplacementScore    float64         `json:"replacement_score"`
	OptimalReplaceDate  time.Time       `json:"optimal_replace_date"`
}

// DepreciationRecord represents depreciation data
type DepreciationRecord struct {
	ID                uuid.UUID       `json:"id"`
	DeviceID          uuid.UUID       `json:"device_id"`
	OriginalValue     decimal.Decimal `json:"original_value"`
	CurrentValue      decimal.Decimal `json:"current_value"`
	DepreciatedAmount decimal.Decimal `json:"depreciated_amount"`
	DepreciationRate  float64         `json:"depreciation_rate"`
	Method            string          `json:"method"` // linear, accelerated, market_based
	RecordedAt        time.Time       `json:"recorded_at"`
}

// WarrantyEvent represents warranty event
type WarrantyEvent struct {
	ID         uuid.UUID `json:"id"`
	DeviceID   uuid.UUID `json:"device_id"`
	EventType  string    `json:"event_type"` // activated, extended, expired, claimed
	WarrantyID string    `json:"warranty_id"`
	Provider   string    `json:"provider"`
	Coverage   string    `json:"coverage_details"`
	ValidUntil time.Time `json:"valid_until"`
	ClaimCount int       `json:"claim_count"`
	RecordedAt time.Time `json:"recorded_at"`
}

// InsuranceEvent represents insurance event
type InsuranceEvent struct {
	ID         uuid.UUID       `json:"id"`
	DeviceID   uuid.UUID       `json:"device_id"`
	EventType  string          `json:"event_type"` // enrolled, renewed, cancelled, claimed
	PolicyID   uuid.UUID       `json:"policy_id"`
	Premium    decimal.Decimal `json:"premium"`
	Coverage   string          `json:"coverage_type"`
	Deductible decimal.Decimal `json:"deductible"`
	ValidFrom  time.Time       `json:"valid_from"`
	ValidUntil time.Time       `json:"valid_until"`
	RecordedAt time.Time       `json:"recorded_at"`
}

// RiskScoreChange represents risk score change
type RiskScoreChange struct {
	ID            uuid.UUID          `json:"id"`
	DeviceID      uuid.UUID          `json:"device_id"`
	PreviousScore float64            `json:"previous_score"`
	NewScore      float64            `json:"new_score"`
	ChangeReason  string             `json:"change_reason"`
	RiskFactors   map[string]float64 `json:"risk_factors"`
	RecordedAt    time.Time          `json:"recorded_at"`
}

// IncidentRecord represents incident data
type IncidentRecord struct {
	ID           uuid.UUID  `json:"id"`
	DeviceID     uuid.UUID  `json:"device_id"`
	IncidentType string     `json:"incident_type"` // theft, damage, malfunction, security_breach
	Severity     string     `json:"severity"`
	Description  string     `json:"description"`
	Resolution   string     `json:"resolution"`
	Impact       string     `json:"impact"`
	ReportedBy   uuid.UUID  `json:"reported_by"`
	ReportedAt   time.Time  `json:"reported_at"`
	ResolvedAt   *time.Time `json:"resolved_at"`
}

// ComplianceEvent represents compliance event
type ComplianceEvent struct {
	ID              uuid.UUID  `json:"id"`
	DeviceID        uuid.UUID  `json:"device_id"`
	ComplianceType  string     `json:"compliance_type"`
	Status          string     `json:"status"` // compliant, non_compliant, pending
	Requirements    []string   `json:"requirements"`
	Violations      []string   `json:"violations"`
	CertificationID string     `json:"certification_id"`
	ValidUntil      *time.Time `json:"valid_until"`
	RecordedAt      time.Time  `json:"recorded_at"`
}

// FailurePrediction represents failure prediction
type FailurePrediction struct {
	DeviceID           uuid.UUID       `json:"device_id"`
	PredictionDate     time.Time       `json:"prediction_date"`
	FailureProbability float64         `json:"failure_probability"`
	PredictedComponent string          `json:"predicted_component"`
	TimeToFailureDays  int             `json:"time_to_failure_days"`
	Confidence         float64         `json:"confidence"`
	RecommendedAction  string          `json:"recommended_action"`
	PreventiveCost     decimal.Decimal `json:"preventive_cost"`
}

// UpgradeRecommendation represents upgrade recommendation
type UpgradeRecommendation struct {
	DeviceID           uuid.UUID       `json:"device_id"`
	RecommendationType string          `json:"recommendation_type"` // immediate, planned, optional
	Reason             string          `json:"reason"`
	CurrentValue       decimal.Decimal `json:"current_value"`
	UpgradeCost        decimal.Decimal `json:"upgrade_cost"`
	RecommendedModels  []string        `json:"recommended_models"`
	OptimalUpgradeDate time.Time       `json:"optimal_upgrade_date"`
	SavingsIfUpgrade   decimal.Decimal `json:"savings_if_upgrade"`
}

// ComparativeAnalysis represents comparative analysis
type ComparativeAnalysis struct {
	DeviceID        uuid.UUID          `json:"device_id"`
	ComparisonGroup string             `json:"comparison_group"`
	DeviceCount     int                `json:"device_count"`
	PerformanceRank int                `json:"performance_rank"`
	ReliabilityRank int                `json:"reliability_rank"`
	ValueRank       int                `json:"value_rank"`
	Percentile      float64            `json:"percentile"`
	Strengths       []string           `json:"strengths"`
	Weaknesses      []string           `json:"weaknesses"`
	Metrics         map[string]float64 `json:"comparative_metrics"`
}

// CategoryBenchmarks represents category benchmarks
type CategoryBenchmarks struct {
	Category          string          `json:"category"`
	AvgLifespanMonths int             `json:"avg_lifespan_months"`
	AvgRepairCount    float64         `json:"avg_repair_count"`
	AvgClaimCount     float64         `json:"avg_claim_count"`
	AvgTCO            decimal.Decimal `json:"avg_total_cost_ownership"`
	AvgDepreciation   float64         `json:"avg_depreciation_rate"`
	TopFailures       []string        `json:"top_failures"`
	BestPractices     []string        `json:"best_practices"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
