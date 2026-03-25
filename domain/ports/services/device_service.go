package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	// "smartsure/internal/domain/models/device" // Unused import
)

// DeviceService defines business operations for devices
type DeviceService interface {
	// Registration and Setup
	RegisterDevice(ctx context.Context, device *models.Device) error
	VerifyDevice(ctx context.Context, deviceID uuid.UUID) error
	UpdateVerification(ctx context.Context, deviceID uuid.UUID) error

	// Insurance and Risk Assessment
	CheckInsuranceEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error)
	CalculateDepreciation(ctx context.Context, deviceID uuid.UUID) (float64, error)
	AssessRisk(ctx context.Context, deviceID uuid.UUID) (*RiskAssessment, error)
	CalculatePremium(ctx context.Context, deviceID uuid.UUID) (*Premium, error)
	GetCurrentInsurableValue(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// Claims Processing
	ProcessClaim(ctx context.Context, deviceID uuid.UUID, claimID uuid.UUID) error
	GetActiveClaimCount(ctx context.Context, deviceID uuid.UUID) (int, error)
	GetTotalClaimAmount(ctx context.Context, deviceID uuid.UUID) (float64, error)
	ValidateClaimEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error)

	// Fraud Detection
	DetectFraud(ctx context.Context, deviceID uuid.UUID) (*FraudAnalysis, error)
	CheckBlacklistStatus(ctx context.Context, deviceID uuid.UUID) (bool, error)
	ReportStolen(ctx context.Context, deviceID uuid.UUID) error

	// Trade-In and Valuation
	CalculateTradeInValue(ctx context.Context, deviceID uuid.UUID) (float64, error)
	CheckTradeInEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error)
	ProcessTradeIn(ctx context.Context, deviceID uuid.UUID, newDeviceID uuid.UUID) error

	// Rental and Financing
	CheckRentalEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error)
	InitiateRental(ctx context.Context, deviceID uuid.UUID, duration time.Duration) error
	InitiateFinancing(ctx context.Context, deviceID uuid.UUID, terms *FinancingTerms) error

	// Repairs and Maintenance
	ScheduleRepair(ctx context.Context, deviceID uuid.UUID, repairType string) error
	GetRepairHistory(ctx context.Context, deviceID uuid.UUID) ([]*RepairRecord, error)
	GetTotalRepairCost(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// Analytics and Intelligence
	AnalyzeUsagePatterns(ctx context.Context, deviceID uuid.UUID) (*UsageAnalysis, error)
	PredictFailure(ctx context.Context, deviceID uuid.UUID) (*FailurePrediction, error)
	GenerateUpgradeRecommendation(ctx context.Context, deviceID uuid.UUID) (*UpgradeRecommendation, error)
	CalculateBehaviorScore(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// Compliance and Legal
	CheckCompliance(ctx context.Context, deviceID uuid.UUID) (*ComplianceReport, error)
	ApplyLegalHold(ctx context.Context, deviceID uuid.UUID, reason string) error
	GenerateRegulatoryReport(ctx context.Context, deviceID uuid.UUID) (*RegulatoryReport, error)

	// Sustainability Metrics
	CalculateCarbonFootprint(ctx context.Context, deviceID uuid.UUID) (float64, error)
	GetRecyclingScore(ctx context.Context, deviceID uuid.UUID) (int, error)
	GetRepairabilityScore(ctx context.Context, deviceID uuid.UUID) (int, error)
	CheckEcoFriendlyCertification(ctx context.Context, deviceID uuid.UUID) (bool, error)

	// Customer Experience
	GetSatisfactionScore(ctx context.Context, deviceID uuid.UUID) (float64, error)
	GenerateRecommendations(ctx context.Context, deviceID uuid.UUID) ([]*Recommendation, error)
	ProcessFeedback(ctx context.Context, deviceID uuid.UUID, feedback *CustomerFeedback) error

	// Corporate Management
	AssignToCorporateUser(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) error
	RegisterBYOD(ctx context.Context, deviceID uuid.UUID, programID uuid.UUID) error
	AddToPool(ctx context.Context, deviceID uuid.UUID, poolID uuid.UUID) error
	GetCorporateDeviceStatus(ctx context.Context, deviceID uuid.UUID) (*CorporateStatus, error)

	// ======================================
	// MISSING IMPLEMENTED METHODS
	// ======================================

	// Analytics & Monitoring
	GetDeviceAnalytics(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	GetDeviceUsageHistory(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)
	GetDeviceHealthScore(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)

	// Device Groups
	CreateDeviceGroup(ctx context.Context, ownerID uuid.UUID, groupName, description string, deviceIDs []uuid.UUID, groupType string, isPublic bool, tags []string) (interface{}, error) // DeviceGroup type not found
	GetDeviceGroups(ctx context.Context, ownerID uuid.UUID) ([]interface{}, error)                                                                                                       // DeviceGroup type not found
	UpdateDeviceGroup(ctx context.Context, groupID, ownerID uuid.UUID, groupName, description string, deviceIDs []uuid.UUID, groupType string, isPublic *bool, tags []string) error
	DeleteDeviceGroup(ctx context.Context, groupID, ownerID uuid.UUID) error

	// Device Transfers
	InitiateDeviceTransfer(ctx context.Context, deviceID, senderID uuid.UUID, recipientEmail, transferNote, transferType, expiryDate string) (interface{}, error) // DeviceTransfer type not found
	AcceptDeviceTransfer(ctx context.Context, transferID, recipientID uuid.UUID, note string) error
	RejectDeviceTransfer(ctx context.Context, transferID, recipientID uuid.UUID, reason, note string) error

	// Bulk Operations
	BulkRegisterDevices(ctx context.Context, devices []RegisterDeviceRequest, groupID *uuid.UUID) ([]*BulkOperationResult, error)
	BulkUpdateDevices(ctx context.Context, ownerID uuid.UUID, deviceUpdates []struct {
		DeviceID uuid.UUID
		Updates  map[string]interface{}
	}) ([]*BulkOperationResult, error)

	// Legacy methods for backward compatibility
	GetDeviceByID(ctx context.Context, deviceID uuid.UUID) (*models.Device, error)
	GetDeviceByIMEI(ctx context.Context, imei string) (*models.Device, error)
	UpdateDevice(ctx context.Context, deviceID uuid.UUID, updates map[string]interface{}) (*models.Device, error)
	DeleteDevice(ctx context.Context, deviceID uuid.UUID) error
	GetUserDevices(ctx context.Context, userID uuid.UUID) ([]*models.Device, error)
	GradeDevice(ctx context.Context, req *DeviceGradingRequest) (*DeviceGradingResponse, error)
	GetDevicesByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Device, error)
	GetEligibleDevices(ctx context.Context, ownerID uuid.UUID) ([]models.Device, error)
}

// Supporting types for DeviceService

// RiskAssessment contains risk analysis results
type RiskAssessment struct {
	DeviceID  uuid.UUID `json:"device_id"`
	RiskScore float64   `json:"risk_score"`
	RiskLevel string    `json:"risk_level"`
	Factors   []string  `json:"factors"`
	Timestamp time.Time `json:"timestamp"`
}

// Premium contains premium calculation details
type Premium struct {
	DeviceID    uuid.UUID    `json:"device_id"`
	BaseAmount  float64      `json:"base_amount"`
	Adjustments []Adjustment `json:"adjustments"`
	FinalAmount float64      `json:"final_amount"`
	ValidUntil  time.Time    `json:"valid_until"`
}

// Adjustment represents a premium adjustment
type Adjustment struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
}

// FraudAnalysis contains fraud detection results
type FraudAnalysis struct {
	DeviceID    uuid.UUID `json:"device_id"`
	FraudScore  float64   `json:"fraud_score"`
	RiskLevel   string    `json:"risk_level"`
	Indicators  []string  `json:"indicators"`
	Recommended string    `json:"recommended_action"`
}

// FinancingTerms contains financing parameters
type FinancingTerms struct {
	Duration      int     `json:"duration_months"`
	InterestRate  float64 `json:"interest_rate"`
	DownPayment   float64 `json:"down_payment"`
	MonthlyAmount float64 `json:"monthly_amount"`
}

// RepairRecord contains repair history information
type RepairRecord struct {
	ID          uuid.UUID `json:"id"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	Status      string    `json:"status"`
}

// UsageAnalysis contains usage pattern analysis
type UsageAnalysis struct {
	DeviceID        uuid.UUID `json:"device_id"`
	AverageDaily    float64   `json:"average_daily_hours"`
	TopApps         []string  `json:"top_apps"`
	DataConsumption float64   `json:"data_consumption_gb"`
	BatteryHealth   int       `json:"battery_health"`
}

// FailurePrediction contains failure prediction analysis
type FailurePrediction struct {
	DeviceID    uuid.UUID `json:"device_id"`
	Component   string    `json:"component"`
	Probability float64   `json:"probability"`
	TimeFrame   string    `json:"timeframe"`
	Preventable bool      `json:"preventable"`
}

// UpgradeRecommendation contains upgrade suggestions
type UpgradeRecommendation struct {
	DeviceID         uuid.UUID `json:"device_id"`
	RecommendedModel string    `json:"recommended_model"`
	Reason           string    `json:"reason"`
	EstimatedValue   float64   `json:"estimated_trade_value"`
	Savings          float64   `json:"potential_savings"`
}

// ComplianceReport contains compliance check results
type ComplianceReport struct {
	DeviceID    uuid.UUID `json:"device_id"`
	IsCompliant bool      `json:"is_compliant"`
	Violations  []string  `json:"violations"`
	LastChecked time.Time `json:"last_checked"`
}

// RegulatoryReport contains regulatory reporting information
type RegulatoryReport struct {
	DeviceID    uuid.UUID `json:"device_id"`
	ReportType  string    `json:"report_type"`
	Content     string    `json:"content"`
	GeneratedAt time.Time `json:"generated_at"`
}

// Recommendation contains device recommendations
type Recommendation struct {
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

// CustomerFeedback contains customer feedback data
type CustomerFeedback struct {
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	Category    string    `json:"category"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// CorporateStatus contains corporate device status
type CorporateStatus struct {
	IsCorporate      bool       `json:"is_corporate"`
	IsBYOD           bool       `json:"is_byod"`
	AssignedTo       *uuid.UUID `json:"assigned_to"`
	Department       string     `json:"department"`
	ComplianceStatus string     `json:"compliance_status"`
}

// RegisterDeviceRequest represents a request to register a new device
type RegisterDeviceRequest struct {
	IMEI            string    `json:"imei" binding:"required"`
	SerialNumber    string    `json:"serial_number"`
	Model           string    `json:"model" binding:"required"`
	Brand           string    `json:"brand" binding:"required"`
	Manufacturer    string    `json:"manufacturer" binding:"required"`
	OperatingSystem string    `json:"operating_system"`
	OSVersion       string    `json:"os_version"`
	StorageCapacity int       `json:"storage_capacity"`
	RAM             int       `json:"ram"`
	Color           string    `json:"color"`
	PurchaseDate    string    `json:"purchase_date"`
	PurchasePrice   float64   `json:"purchase_price"`
	OwnerID         uuid.UUID `json:"owner_id"`
}

// DeviceGradingRequest represents a request to grade a device
type DeviceGradingRequest struct {
	DeviceID      uuid.UUID `json:"device_id" binding:"required"`
	Condition     string    `json:"condition"`
	ScreenDamage  bool      `json:"screen_damage"`
	WaterDamage   bool      `json:"water_damage"`
	BatteryHealth int       `json:"battery_health"`
	StorageUsed   int       `json:"storage_used"`
}

// DeviceGradingResponse represents the response from device grading
type DeviceGradingResponse struct {
	DeviceID  uuid.UUID `json:"device_id"`
	Grade     string    `json:"grade"`
	Score     int       `json:"score"`
	Value     float64   `json:"value"`
	Condition string    `json:"condition"`
	Comments  []string  `json:"comments"`
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	Index    int       `json:"index"`
	IMEI     string    `json:"imei,omitempty"`
	DeviceID uuid.UUID `json:"device_id,omitempty"`
	Success  bool      `json:"success"`
	Error    string    `json:"error,omitempty"`
}
