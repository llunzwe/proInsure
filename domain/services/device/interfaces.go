package device

import (
	"context"
	"time"

	"smartsure/internal/domain/models"
	devicemodels "smartsure/internal/domain/models/device"

	"github.com/google/uuid"
)

// Common types used across device services
type (
	// ValuationBreakdown represents detailed device valuation
	ValuationBreakdown struct {
		PurchasePrice       float64            `json:"purchase_price"`
		CurrentMarketValue  float64            `json:"current_market_value"`
		DepreciationRate    float64            `json:"depreciation_rate"`
		ConditionAdjustment float64            `json:"condition_adjustment"`
		SegmentMultiplier   float64            `json:"segment_multiplier"`
		ComponentValues     map[string]float64 `json:"component_values"`
		FinalValue          float64            `json:"final_value"`
		CalculatedAt        time.Time          `json:"calculated_at"`
	}

	// DeviceValidationError represents a validation error
	DeviceValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
		Code    string `json:"code"`
	}

	// DeviceComplianceReport represents device compliance check results
	DeviceComplianceReport struct {
		DeviceID      uuid.UUID `json:"device_id"`
		Compliant     bool      `json:"compliant"`
		Restrictions  []string  `json:"restrictions"`
		Requirements  []string  `json:"requirements"`
		CheckedAt     time.Time `json:"checked_at"`
		NextCheckDate time.Time `json:"next_check_date"`
	}

	// DeviceRiskAssessment represents device risk evaluation results
	DeviceRiskAssessment struct {
		DeviceID    uuid.UUID          `json:"device_id"`
		RiskScore   float64            `json:"risk_score"`
		RiskLevel   string             `json:"risk_level"`
		RiskFactors map[string]float64 `json:"risk_factors"`
		TheftRisk   float64            `json:"theft_risk"`
		FraudRisk   float64            `json:"fraud_risk"`
		AssessedAt  time.Time          `json:"assessed_at"`
		ValidUntil  time.Time          `json:"valid_until"`
	}

	// DeviceHealthReport represents device health assessment
	DeviceHealthReport struct {
		DeviceID           uuid.UUID         `json:"device_id"`
		OverallHealth      float64           `json:"overall_health"`
		BatteryHealth      int               `json:"battery_health"`
		ScreenHealth       string            `json:"screen_health"`
		ComponentHealth    map[string]string `json:"component_health"`
		FunctionalIssues   []string          `json:"functional_issues"`
		RecommendedActions []string          `json:"recommended_actions"`
		AssessmentDate     time.Time         `json:"assessment_date"`
	}

	// DeviceVerificationResult represents device verification outcome
	DeviceVerificationResult struct {
		DeviceID           uuid.UUID `json:"device_id"`
		IMEI               string    `json:"imei"`
		Verified           bool      `json:"verified"`
		AuthenticityStatus string    `json:"authenticity_status"`
		BlacklistStatus    string    `json:"blacklist_status"`
		GreyMarket         bool      `json:"grey_market"`
		VerificationMethod string    `json:"verification_method"`
		VerifiedAt         time.Time `json:"verified_at"`
		VerifiedBy         uuid.UUID `json:"verified_by"`
		Notes              string    `json:"notes"`
	}

	// DeviceEligibility represents device insurance eligibility
	DeviceEligibility struct {
		DeviceID           uuid.UUID `json:"device_id"`
		Eligible           bool      `json:"eligible"`
		Reasons            []string  `json:"reasons"`
		Restrictions       []string  `json:"restrictions"`
		AvailableCoverages []string  `json:"available_coverages"`
		MaxCoverageAmount  float64   `json:"max_coverage_amount"`
		RecommendedPremium float64   `json:"recommended_premium"`
		EvaluatedAt        time.Time `json:"evaluated_at"`
	}

	// DeviceInspectionReport represents device inspection results
	DeviceInspectionReport struct {
		InspectionID      uuid.UUID         `json:"inspection_id"`
		DeviceID          uuid.UUID         `json:"device_id"`
		InspectorID       uuid.UUID         `json:"inspector_id"`
		InspectionType    string            `json:"inspection_type"`
		PhysicalCondition map[string]string `json:"physical_condition"`
		FunctionalTests   map[string]bool   `json:"functional_tests"`
		PhotoEvidence     []string          `json:"photo_evidence"`
		Grade             string            `json:"grade"`
		PassedInspection  bool              `json:"passed_inspection"`
		Issues            []string          `json:"issues"`
		InspectionDate    time.Time         `json:"inspection_date"`
	}

	// DeviceTrackingInfo represents device location and tracking data
	DeviceTrackingInfo struct {
		DeviceID          uuid.UUID `json:"device_id"`
		CurrentLocation   Location  `json:"current_location"`
		LastKnownLocation Location  `json:"last_known_location"`
		TrackingEnabled   bool      `json:"tracking_enabled"`
		RemoteLockStatus  string    `json:"remote_lock_status"`
		FindMyEnabled     bool      `json:"find_my_enabled"`
		LastUpdated       time.Time `json:"last_updated"`
	}

	// Location represents geographic location
	Location struct {
		Latitude  float64   `json:"latitude"`
		Longitude float64   `json:"longitude"`
		Accuracy  float64   `json:"accuracy"`
		Address   string    `json:"address"`
		Timestamp time.Time `json:"timestamp"`
	}

	// DeviceMarketAnalysis represents market analysis for device
	DeviceMarketAnalysis struct {
		DeviceID     uuid.UUID    `json:"device_id"`
		Brand        string       `json:"brand"`
		Model        string       `json:"model"`
		MarketValue  float64      `json:"market_value"`
		ResaleValue  float64      `json:"resale_value"`
		TradeInValue float64      `json:"trade_in_value"`
		DemandLevel  string       `json:"demand_level"`
		PriceHistory []PricePoint `json:"price_history"`
		MarketTrend  string       `json:"market_trend"`
		AnalyzedAt   time.Time    `json:"analyzed_at"`
	}

	// PricePoint represents a price at a point in time
	PricePoint struct {
		Price  float64   `json:"price"`
		Date   time.Time `json:"date"`
		Source string    `json:"source"`
	}

	// DeviceServiceHistory represents device service and repair history
	DeviceServiceHistory struct {
		DeviceID         uuid.UUID         `json:"device_id"`
		TotalRepairs     int               `json:"total_repairs"`
		TotalCost        float64           `json:"total_cost"`
		LastServiceDate  *time.Time        `json:"last_service_date"`
		Repairs          []RepairRecord    `json:"repairs"`
		WarrantyServices []WarrantyService `json:"warranty_services"`
	}

	// RepairRecord represents a single repair event
	RepairRecord struct {
		RepairID     uuid.UUID `json:"repair_id"`
		Type         string    `json:"type"`
		Description  string    `json:"description"`
		Cost         float64   `json:"cost"`
		TechnicianID uuid.UUID `json:"technician_id"`
		RepairDate   time.Time `json:"repair_date"`
		Warranty     bool      `json:"warranty"`
	}

	// WarrantyService represents warranty service record
	WarrantyService struct {
		ServiceID   uuid.UUID `json:"service_id"`
		ServiceType string    `json:"service_type"`
		Provider    string    `json:"provider"`
		ServiceDate time.Time `json:"service_date"`
		Description string    `json:"description"`
	}

	// DeviceNetworkInfo represents device network and connectivity information
	DeviceNetworkInfo struct {
		DeviceID        uuid.UUID `json:"device_id"`
		NetworkOperator string    `json:"network_operator"`
		NetworkStatus   string    `json:"network_status"`
		SignalStrength  int       `json:"signal_strength"`
		DataUsage       DataUsage `json:"data_usage"`
		IsRoaming       bool      `json:"is_roaming"`
		LastUpdated     time.Time `json:"last_updated"`
	}

	// DataUsage represents device data usage statistics
	DataUsage struct {
		TotalUsed     int64     `json:"total_used"`
		MonthlyLimit  int64     `json:"monthly_limit"`
		CurrentPeriod string    `json:"current_period"`
		LastReset     time.Time `json:"last_reset"`
	}

	// DevicePerformanceMetrics represents device performance data
	DevicePerformanceMetrics struct {
		DeviceID      uuid.UUID `json:"device_id"`
		CPUUsage      float64   `json:"cpu_usage"`
		MemoryUsage   float64   `json:"memory_usage"`
		StorageUsed   int64     `json:"storage_used"`
		StorageTotal  int64     `json:"storage_total"`
		BatteryLevel  int       `json:"battery_level"`
		Temperature   float64   `json:"temperature"`
		AppCrashCount int       `json:"app_crash_count"`
		SystemUptime  int64     `json:"system_uptime"`
		CollectedAt   time.Time `json:"collected_at"`
	}
)

// DeviceRepository defines the interface for device data operations
type DeviceRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error)
	GetByIMEI(ctx context.Context, imei string) (*models.Device, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (*models.Device, error)
	Update(ctx context.Context, device *models.Device) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Listing and search
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*models.Device, int64, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*models.Device, int64, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*models.Device, error)
	GetByCorporateAccountID(ctx context.Context, corporateID uuid.UUID) ([]*models.Device, error)

	// Bulk operations
	BulkCreate(ctx context.Context, devices []*models.Device) error
	BulkUpdate(ctx context.Context, devices []*models.Device) error
	BulkUpdateStatus(ctx context.Context, deviceIDs []uuid.UUID, status string) error

	// Specialized queries
	GetExpiredWarrantyDevices(ctx context.Context, days int) ([]*models.Device, error)
	GetDevicesNeedingInspection(ctx context.Context) ([]*models.Device, error)
	GetStolenDevices(ctx context.Context) ([]*models.Device, error)
	GetBlacklistedDevices(ctx context.Context) ([]*models.Device, error)
	GetDevicesByRiskLevel(ctx context.Context, riskLevel string) ([]*models.Device, error)
	GetUnverifiedDevices(ctx context.Context, days int) ([]*models.Device, error)
}

// DeviceService defines the interface for device business logic
type DeviceService interface {
	// Registration and management
	RegisterDevice(ctx context.Context, device *models.Device) error
	UpdateDevice(ctx context.Context, device *models.Device) error
	DeactivateDevice(ctx context.Context, deviceID uuid.UUID, reason string) error
	TransferOwnership(ctx context.Context, deviceID uuid.UUID, newOwnerID uuid.UUID) error

	// Verification and validation
	VerifyDevice(ctx context.Context, deviceID uuid.UUID) (*DeviceVerificationResult, error)
	ValidateIMEI(ctx context.Context, imei string) error
	CheckBlacklistStatus(ctx context.Context, imei string) (string, error)
	VerifyAuthenticity(ctx context.Context, deviceID uuid.UUID) error

	// Risk and compliance
	AssessDeviceRisk(ctx context.Context, deviceID uuid.UUID) (*DeviceRiskAssessment, error)
	CheckCompliance(ctx context.Context, deviceID uuid.UUID) (*DeviceComplianceReport, error)
	EvaluateEligibility(ctx context.Context, deviceID uuid.UUID) (*DeviceEligibility, error)

	// Valuation and depreciation
	CalculateCurrentValue(ctx context.Context, deviceID uuid.UUID) (*ValuationBreakdown, error)
	GetMarketAnalysis(ctx context.Context, deviceID uuid.UUID) (*DeviceMarketAnalysis, error)
	CalculateTradeInValue(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// Health and inspection
	AssessDeviceHealth(ctx context.Context, deviceID uuid.UUID) (*DeviceHealthReport, error)
	ScheduleInspection(ctx context.Context, deviceID uuid.UUID, inspectorID uuid.UUID) error
	SubmitInspectionReport(ctx context.Context, report *DeviceInspectionReport) error

	// Tracking and security
	UpdateLocation(ctx context.Context, deviceID uuid.UUID, location Location) error
	GetTrackingInfo(ctx context.Context, deviceID uuid.UUID) (*DeviceTrackingInfo, error)
	RemoteLock(ctx context.Context, deviceID uuid.UUID) error
	RemoteWipe(ctx context.Context, deviceID uuid.UUID) error
	ReportStolen(ctx context.Context, deviceID uuid.UUID, details map[string]interface{}) error

	// Service history
	GetServiceHistory(ctx context.Context, deviceID uuid.UUID) (*DeviceServiceHistory, error)
	RecordRepair(ctx context.Context, deviceID uuid.UUID, repair RepairRecord) error

	// Network and performance
	GetNetworkInfo(ctx context.Context, deviceID uuid.UUID) (*DeviceNetworkInfo, error)
	GetPerformanceMetrics(ctx context.Context, deviceID uuid.UUID) (*DevicePerformanceMetrics, error)
	UpdatePerformanceMetrics(ctx context.Context, metrics *DevicePerformanceMetrics) error

	// Bulk operations
	BulkRegister(ctx context.Context, devices []*models.Device) error
	BulkVerify(ctx context.Context, deviceIDs []uuid.UUID) (map[uuid.UUID]*DeviceVerificationResult, error)
	BulkAssessRisk(ctx context.Context, deviceIDs []uuid.UUID) (map[uuid.UUID]*DeviceRiskAssessment, error)
}

// DeviceValidationService defines the interface for device validation
type DeviceValidationService interface {
	ValidateDevice(ctx context.Context, device *models.Device) []DeviceValidationError
	ValidateIMEI(ctx context.Context, imei string) error
	ValidateSerialNumber(ctx context.Context, serialNumber string) error
	ValidateOwnership(ctx context.Context, deviceID uuid.UUID, ownerID uuid.UUID) error
	ValidateEligibility(ctx context.Context, device *models.Device) error
	ValidateTransfer(ctx context.Context, deviceID uuid.UUID, fromOwnerID, toOwnerID uuid.UUID) error
}

// DeviceEventService defines the interface for device event management
type DeviceEventService interface {
	RecordEvent(ctx context.Context, deviceID uuid.UUID, eventType string, details map[string]interface{}) error
	GetDeviceEvents(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]devicemodels.DeviceHistory, error)
	GetEventsByType(ctx context.Context, deviceID uuid.UUID, eventType string) ([]devicemodels.DeviceHistory, error)
	GetRecentEvents(ctx context.Context, deviceID uuid.UUID, hours int) ([]devicemodels.DeviceHistory, error)
}

// DeviceNotificationService defines the interface for device-related notifications
type DeviceNotificationService interface {
	SendVerificationReminder(ctx context.Context, deviceID uuid.UUID) error
	SendInspectionDue(ctx context.Context, deviceID uuid.UUID) error
	SendWarrantyExpiring(ctx context.Context, deviceID uuid.UUID) error
	SendRiskAlert(ctx context.Context, deviceID uuid.UUID, riskLevel string) error
	SendMaintenanceRecommendation(ctx context.Context, deviceID uuid.UUID, recommendations []string) error
}

// DeviceAnalyticsService defines the interface for device analytics
type DeviceAnalyticsService interface {
	GetDeviceStatistics(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	GetFleetAnalytics(ctx context.Context, corporateID uuid.UUID) (map[string]interface{}, error)
	GetRiskTrends(ctx context.Context, deviceID uuid.UUID, days int) ([]map[string]interface{}, error)
	GetValueTrends(ctx context.Context, deviceID uuid.UUID, days int) ([]PricePoint, error)
	GetFailurePrediction(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	GetUsagePatterns(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
}

// DeviceIntegrationService defines the interface for third-party integrations
type DeviceIntegrationService interface {
	SyncWithManufacturer(ctx context.Context, deviceID uuid.UUID) error
	SyncWithCarrier(ctx context.Context, deviceID uuid.UUID) error
	VerifyWithGSMA(ctx context.Context, imei string) (*DeviceVerificationResult, error)
	CheckWithCheckMEND(ctx context.Context, imei string) (map[string]interface{}, error)
	GetManufacturerWarranty(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
}
