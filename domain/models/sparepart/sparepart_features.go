package sparepart

import (
	"time"
	
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SparePartRequest represents a request for spare parts
type SparePartRequest struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID      uuid.UUID  `gorm:"type:uuid;not null" json:"spare_part_id"`
	RequestNumber    string     `gorm:"uniqueIndex;not null" json:"request_number"`
	RequestType      string     `gorm:"not null" json:"request_type"`
	Priority         string     `gorm:"default:'normal'" json:"priority"`
	Status           string     `gorm:"default:'pending'" json:"status"`
	RequestedBy      uuid.UUID  `gorm:"type:uuid;not null" json:"requested_by"`
	RequestedFor     string     `json:"requested_for,omitempty"` // repair_id, technician_id, etc.
	Quantity         int        `gorm:"not null" json:"quantity"`
	QuantityApproved int        `json:"quantity_approved"`
	QuantityIssued   int        `json:"quantity_issued"`
	RequiredDate     time.Time  `json:"required_date"`
	Purpose          string     `json:"purpose"`
	JobReference     string     `json:"job_reference,omitempty"`
	CustomerID       *uuid.UUID `gorm:"type:uuid" json:"customer_id,omitempty"`
	DeviceID         *uuid.UUID `gorm:"type:uuid" json:"device_id,omitempty"`
	RepairID         *uuid.UUID `gorm:"type:uuid" json:"repair_id,omitempty"`
	CostCenter       string     `json:"cost_center,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	IssuedBy         *uuid.UUID `gorm:"type:uuid" json:"issued_by,omitempty"`
	IssuedAt         *time.Time `json:"issued_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	RejectionReason  string     `json:"rejection_reason,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartMovement tracks all part movements
type SparePartMovement struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID    uuid.UUID       `gorm:"type:uuid;not null" json:"spare_part_id"`
	MovementNumber string          `gorm:"uniqueIndex;not null" json:"movement_number"`
	MovementType   string          `gorm:"not null" json:"movement_type"`
	ReferenceType  string          `json:"reference_type,omitempty"`
	ReferenceID    string          `json:"reference_id,omitempty"`
	Quantity       int             `gorm:"not null" json:"quantity"`
	FromLocation   string          `json:"from_location,omitempty"`
	ToLocation     string          `json:"to_location,omitempty"`
	FromStatus     string          `json:"from_status,omitempty"`
	ToStatus       string          `json:"to_status,omitempty"`
	UnitCost       decimal.Decimal `sql:"type:decimal(10,2)" json:"unit_cost"`
	TotalCost      decimal.Decimal `sql:"type:decimal(10,2)" json:"total_cost"`
	StockBefore    int             `json:"stock_before"`
	StockAfter     int             `json:"stock_after"`
	Reason         string          `json:"reason"`
	BatchNumber    string          `json:"batch_number,omitempty"`
	SerialNumbers  string          `json:"serial_numbers,omitempty"` // JSON array
	ExpiryDate     *time.Time      `json:"expiry_date,omitempty"`
	QualityChecked bool            `gorm:"default:false" json:"quality_checked"`
	PerformedBy    uuid.UUID       `gorm:"type:uuid;not null" json:"performed_by"`
	AuthorizedBy   *uuid.UUID      `gorm:"type:uuid" json:"authorized_by,omitempty"`
	Notes          string          `json:"notes,omitempty"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// SparePartAllocation tracks part allocations to repairs/jobs
type SparePartAllocation struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID      uuid.UUID       `gorm:"type:uuid;not null" json:"spare_part_id"`
	AllocationNumber string          `gorm:"uniqueIndex;not null" json:"allocation_number"`
	AllocatedTo      string          `gorm:"not null" json:"allocated_to"` // repair, technician, job
	AllocatedToID    uuid.UUID       `gorm:"type:uuid;not null" json:"allocated_to_id"`
	Quantity         int             `gorm:"not null" json:"quantity"`
	QuantityUsed     int             `gorm:"default:0" json:"quantity_used"`
	QuantityReturned int             `gorm:"default:0" json:"quantity_returned"`
	Status           string          `gorm:"default:'allocated'" json:"status"`
	AllocatedBy      uuid.UUID       `gorm:"type:uuid;not null" json:"allocated_by"`
	AllocatedAt      time.Time       `json:"allocated_at"`
	ExpectedReturn   *time.Time      `json:"expected_return,omitempty"`
	ActualReturn     *time.Time      `json:"actual_return,omitempty"`
	UsageNotes       string          `json:"usage_notes,omitempty"`
	ReturnCondition  string          `json:"return_condition,omitempty"`
	ChargeType       string          `json:"charge_type,omitempty"` // customer, warranty, internal
	ChargedAmount    decimal.Decimal `sql:"type:decimal(10,2)" json:"charged_amount"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartPurchaseOrder represents purchase orders
type SparePartPurchaseOrder struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID      uuid.UUID       `gorm:"type:uuid;not null" json:"spare_part_id"`
	PONumber         string          `gorm:"uniqueIndex;not null" json:"po_number"`
	SupplierID       uuid.UUID       `gorm:"type:uuid;not null" json:"supplier_id"`
	SupplierName     string          `json:"supplier_name"`
	Quantity         int             `gorm:"not null" json:"quantity"`
	UnitPrice        decimal.Decimal `sql:"type:decimal(10,2)" json:"unit_price"`
	TotalAmount      decimal.Decimal `sql:"type:decimal(10,2)" json:"total_amount"`
	Tax              decimal.Decimal `sql:"type:decimal(10,2)" json:"tax"`
	ShippingCost     decimal.Decimal `sql:"type:decimal(10,2)" json:"shipping_cost"`
	GrandTotal       decimal.Decimal `sql:"type:decimal(10,2)" json:"grand_total"`
	Currency         string          `json:"currency"`
	OrderDate        time.Time       `json:"order_date"`
	ExpectedDelivery time.Time       `json:"expected_delivery"`
	ActualDelivery   *time.Time      `json:"actual_delivery,omitempty"`
	Status           string          `gorm:"default:'pending'" json:"status"`
	PaymentTerms     string          `json:"payment_terms,omitempty"`
	PaymentStatus    string          `json:"payment_status,omitempty"`
	InvoiceNumber    string          `json:"invoice_number,omitempty"`
	TrackingNumber   string          `json:"tracking_number,omitempty"`
	ReceivedQuantity int             `json:"received_quantity"`
	AcceptedQuantity int             `json:"accepted_quantity"`
	RejectedQuantity int             `json:"rejected_quantity"`
	RejectionReason  string          `json:"rejection_reason,omitempty"`
	Notes            string          `json:"notes,omitempty"`
	OrderedBy        uuid.UUID       `gorm:"type:uuid" json:"ordered_by"`
	ApprovedBy       *uuid.UUID      `gorm:"type:uuid" json:"approved_by,omitempty"`
	ReceivedBy       *uuid.UUID      `gorm:"type:uuid" json:"received_by,omitempty"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartReturn handles part returns
type SparePartReturn struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID       uuid.UUID       `gorm:"type:uuid;not null" json:"spare_part_id"`
	ReturnNumber      string          `gorm:"uniqueIndex;not null" json:"return_number"`
	ReturnType        string          `gorm:"not null" json:"return_type"`
	ReturnReason      string          `gorm:"not null" json:"return_reason"`
	Quantity          int             `gorm:"not null" json:"quantity"`
	Condition         string          `json:"condition"`
	BatchNumber       string          `json:"batch_number,omitempty"`
	SerialNumbers     string          `json:"serial_numbers,omitempty"` // JSON array
	OriginalReference string          `json:"original_reference,omitempty"`
	ReturnedBy        uuid.UUID       `gorm:"type:uuid;not null" json:"returned_by"`
	ReturnedFrom      string          `json:"returned_from,omitempty"`
	InspectionStatus  string          `gorm:"default:'pending'" json:"inspection_status"`
	InspectionNotes   string          `json:"inspection_notes,omitempty"`
	InspectedBy       *uuid.UUID      `gorm:"type:uuid" json:"inspected_by,omitempty"`
	InspectedAt       *time.Time      `json:"inspected_at,omitempty"`
	Disposition       string          `json:"disposition,omitempty"` // restock, scrap, return_to_vendor
	RestockQuantity   int             `json:"restock_quantity"`
	ScrapQuantity     int             `json:"scrap_quantity"`
	RefundAmount      decimal.Decimal `sql:"type:decimal(10,2)" json:"refund_amount"`
	RefundStatus      string          `json:"refund_status,omitempty"`
	ProcessedBy       *uuid.UUID      `gorm:"type:uuid" json:"processed_by,omitempty"`
	ProcessedAt       *time.Time      `json:"processed_at,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartQualityCheck tracks quality inspections
type SparePartQualityCheck struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID       uuid.UUID  `gorm:"type:uuid;not null" json:"spare_part_id"`
	CheckNumber       string     `gorm:"uniqueIndex;not null" json:"check_number"`
	CheckType         string     `gorm:"not null" json:"check_type"`
	BatchNumber       string     `json:"batch_number,omitempty"`
	SampleSize        int        `json:"sample_size"`
	TotalQuantity     int        `json:"total_quantity"`
	PassedQuantity    int        `json:"passed_quantity"`
	FailedQuantity    int        `json:"failed_quantity"`
	DefectTypes       string     `json:"defect_types,omitempty"` // JSON array
	TestResults       string     `json:"test_results,omitempty"` // JSON object
	OverallStatus     string     `json:"overall_status"`
	CheckedBy         uuid.UUID  `gorm:"type:uuid;not null" json:"checked_by"`
	CheckedAt         time.Time  `json:"checked_at"`
	Standards         string     `json:"standards,omitempty"` // JSON array
	Equipment         string     `json:"equipment,omitempty"`
	Temperature       float64    `json:"temperature,omitempty"`
	Humidity          float64    `json:"humidity,omitempty"`
	Images            string     `json:"images,omitempty"` // JSON array
	CertificateNumber string     `json:"certificate_number,omitempty"`
	CertificateURL    string     `json:"certificate_url,omitempty"`
	NextCheckDate     *time.Time `json:"next_check_date,omitempty"`
	Notes             string     `json:"notes,omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartUsageHistory tracks detailed usage
type SparePartUsageHistory struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID     uuid.UUID       `gorm:"type:uuid;not null" json:"spare_part_id"`
	UsageDate       time.Time       `gorm:"not null" json:"usage_date"`
	UsedBy          uuid.UUID       `gorm:"type:uuid;not null" json:"used_by"`
	UsedFor         string          `gorm:"not null" json:"used_for"`
	ReferenceType   string          `json:"reference_type"`
	ReferenceID     uuid.UUID       `gorm:"type:uuid" json:"reference_id"`
	Quantity        int             `gorm:"not null" json:"quantity"`
	UnitCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"unit_cost"`
	ChargedCost     decimal.Decimal `sql:"type:decimal(10,2)" json:"charged_cost"`
	DeviceID        *uuid.UUID      `gorm:"type:uuid" json:"device_id,omitempty"`
	CustomerID      *uuid.UUID      `gorm:"type:uuid" json:"customer_id,omitempty"`
	RepairID        *uuid.UUID      `gorm:"type:uuid" json:"repair_id,omitempty"`
	WorkOrderNumber string          `json:"work_order_number,omitempty"`
	Success         bool            `gorm:"default:true" json:"success"`
	FailureReason   string          `json:"failure_reason,omitempty"`
	InstallTime     int             `json:"install_time,omitempty"` // minutes
	Notes           string          `json:"notes,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// SparePartForecast for demand planning
type SparePartForecast struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID      uuid.UUID  `gorm:"type:uuid;not null" json:"spare_part_id"`
	ForecastPeriod   string     `gorm:"not null" json:"forecast_period"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          time.Time  `json:"end_date"`
	PredictedDemand  int        `json:"predicted_demand"`
	MinDemand        int        `json:"min_demand"`
	MaxDemand        int        `json:"max_demand"`
	ConfidenceLevel  float64    `json:"confidence_level"`
	Method           string     `json:"method"` // moving_average, exponential_smoothing, etc.
	SeasonalFactor   float64    `json:"seasonal_factor"`
	TrendFactor      float64    `json:"trend_factor"`
	HistoricalData   string     `json:"historical_data,omitempty"` // JSON object
	Assumptions      string     `json:"assumptions,omitempty"`
	RecommendedOrder int        `json:"recommended_order"`
	RecommendedDate  time.Time  `json:"recommended_date"`
	ActualDemand     int        `json:"actual_demand,omitempty"`
	Accuracy         float64    `json:"accuracy,omitempty"`
	GeneratedBy      string     `json:"generated_by"`
	ReviewedBy       *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartAlert for inventory notifications
type SparePartAlert struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID      uuid.UUID  `gorm:"type:uuid;not null" json:"spare_part_id"`
	AlertType        string     `gorm:"not null" json:"alert_type"`
	AlertLevel       string     `gorm:"default:'info'" json:"alert_level"`
	Title            string     `gorm:"not null" json:"title"`
	Message          string     `json:"message"`
	TriggerValue     string     `json:"trigger_value,omitempty"`
	CurrentValue     string     `json:"current_value,omitempty"`
	Threshold        string     `json:"threshold,omitempty"`
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	IsAcknowledged   bool       `gorm:"default:false" json:"is_acknowledged"`
	AcknowledgedBy   *uuid.UUID `gorm:"type:uuid" json:"acknowledged_by,omitempty"`
	AcknowledgedAt   *time.Time `json:"acknowledged_at,omitempty"`
	ActionRequired   string     `json:"action_required,omitempty"`
	ActionTaken      string     `json:"action_taken,omitempty"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	AutoResolve      bool       `gorm:"default:false" json:"auto_resolve"`
	NotificationSent bool       `gorm:"default:false" json:"notification_sent"`
	Recipients       string     `json:"recipients,omitempty"` // JSON array
	EscalationLevel  int        `gorm:"default:0" json:"escalation_level"`
	NextEscalation   *time.Time `json:"next_escalation,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartCostAnalysis for financial tracking
type SparePartCostAnalysis struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID         uuid.UUID       `gorm:"type:uuid;not null" json:"spare_part_id"`
	AnalysisPeriod      string          `json:"analysis_period"`
	StartDate           time.Time       `json:"start_date"`
	EndDate             time.Time       `json:"end_date"`
	TotalQuantityUsed   int             `json:"total_quantity_used"`
	TotalCost           decimal.Decimal `sql:"type:decimal(10,2)" json:"total_cost"`
	TotalRevenue        decimal.Decimal `sql:"type:decimal(10,2)" json:"total_revenue"`
	GrossProfit         decimal.Decimal `sql:"type:decimal(10,2)" json:"gross_profit"`
	ProfitMargin        float64         `json:"profit_margin"`
	AverageCostPerUnit  decimal.Decimal `sql:"type:decimal(10,2)" json:"average_cost_per_unit"`
	AveragePricePerUnit decimal.Decimal `sql:"type:decimal(10,2)" json:"average_price_per_unit"`
	CarryingCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"carrying_cost"`
	OrderingCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"ordering_cost"`
	StockoutCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"stockout_cost"`
	ObsolescenceCost    decimal.Decimal `sql:"type:decimal(10,2)" json:"obsolescence_cost"`
	TotalCostOwnership  decimal.Decimal `sql:"type:decimal(10,2)" json:"total_cost_ownership"`
	CostTrend           string          `json:"cost_trend,omitempty"`
	Recommendations     string          `json:"recommendations,omitempty"`
	CreatedAt           time.Time       `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy           uuid.UUID       `gorm:"type:uuid" json:"created_by"`
}

// SparePartAuthenticationLog tracks authenticity verification attempts
type SparePartAuthenticationLog struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID      uuid.UUID  `gorm:"type:uuid;not null" json:"spare_part_id"`
	PartNumber       string     `json:"part_number"`
	SerialNumber     string     `json:"serial_number,omitempty"`
	VerificationDate time.Time  `json:"verification_date"`
	VerificationCode string     `json:"verification_code,omitempty"`
	Method           string     `json:"method"` // qr_code, serial, blockchain, hologram
	Result           string     `json:"result"` // genuine, counterfeit, suspicious, inconclusive
	ConfidenceScore  float64    `json:"confidence_score"`
	DeviceIMEI       string     `json:"device_imei,omitempty"`
	DeviceID         *uuid.UUID `gorm:"type:uuid" json:"device_id,omitempty"`
	TechnicianID     *uuid.UUID `gorm:"type:uuid" json:"technician_id,omitempty"`
	RepairID         *uuid.UUID `gorm:"type:uuid" json:"repair_id,omitempty"`
	Location         string     `json:"location,omitempty"`
	IPAddress        string     `json:"ip_address,omitempty"`
	UserAgent        string     `json:"user_agent,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	AlertRaised      bool       `gorm:"default:false" json:"alert_raised"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// SparePartInstallationRecord tracks part installations
type SparePartInstallationRecord struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID      uuid.UUID `gorm:"type:uuid;not null" json:"spare_part_id"`
	DeviceID         uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	DeviceIMEI       string    `json:"device_imei"`
	TechnicianID     uuid.UUID `gorm:"type:uuid;not null" json:"technician_id"`
	RepairID         uuid.UUID `gorm:"type:uuid;not null" json:"repair_id"`
	InstallDate      time.Time `json:"install_date"`
	Duration         int       `json:"duration_minutes"`
	DiagnosticsDone  bool      `gorm:"default:false" json:"diagnostics_done"`
	CalibrationDone  bool      `gorm:"default:false" json:"calibration_done"`
	ProgrammingDone  bool      `gorm:"default:false" json:"programming_done"`
	TestsPassed      bool      `gorm:"default:false" json:"tests_passed"`
	TestResults      string    `json:"test_results,omitempty"` // JSON object
	Success          bool      `gorm:"default:false" json:"success"`
	FailureReason    string    `json:"failure_reason,omitempty"`
	CustomerApproval bool      `gorm:"default:false" json:"customer_approval"`
	WarrantyStart    time.Time `json:"warranty_start"`
	WarrantyEnd      time.Time `json:"warranty_end"`
	Notes            string    `json:"notes,omitempty"`
	PhotosBefore     string    `json:"photos_before,omitempty"` // JSON array
	PhotosAfter      string    `json:"photos_after,omitempty"`  // JSON array
	VideoURL         string    `json:"video_url,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartMarketAnalytics tracks market intelligence
type SparePartMarketAnalytics struct {
	ID                 uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SparePartID        uuid.UUID       `gorm:"type:uuid;not null" json:"spare_part_id"`
	AnalysisPeriod     string          `json:"analysis_period"`
	DemandScore        float64         `json:"demand_score"`
	DemandTrend        string          `json:"demand_trend"` // increasing, stable, decreasing
	SeasonalFactor     float64         `json:"seasonal_factor"`
	PeakSeason         string          `json:"peak_season,omitempty"`
	CompetitorCount    int             `json:"competitor_count"`
	CompetitorMinPrice decimal.Decimal `sql:"type:decimal(10,2)" json:"competitor_min_price"`
	CompetitorMaxPrice decimal.Decimal `sql:"type:decimal(10,2)" json:"competitor_max_price"`
	CompetitorAvgPrice decimal.Decimal `sql:"type:decimal(10,2)" json:"competitor_avg_price"`
	MarketShare        float64         `json:"market_share"`
	CustomerRating     float64         `json:"customer_rating"`
	PopularityRank     int             `json:"popularity_rank"`
	SearchVolume       int             `json:"search_volume"`
	ConversionRate     float64         `json:"conversion_rate"`
	ReturnRate         float64         `json:"return_rate"`
	GeographicDemand   string          `json:"geographic_demand,omitempty"` // JSON object by region
	CustomerSegments   string          `json:"customer_segments,omitempty"` // JSON object
	PriceElasticity    float64         `json:"price_elasticity"`
	RecommendedPrice   decimal.Decimal `sql:"type:decimal(10,2)" json:"recommended_price"`
	ForecastAccuracy   float64         `json:"forecast_accuracy"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// SparePartTradeIn handles part exchange programs
type SparePartTradeIn struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	OldPartID       uuid.UUID       `gorm:"type:uuid;not null" json:"old_part_id"`
	NewPartID       uuid.UUID       `gorm:"type:uuid;not null" json:"new_part_id"`
	CustomerID      uuid.UUID       `gorm:"type:uuid;not null" json:"customer_id"`
	TradeInValue    decimal.Decimal `sql:"type:decimal(10,2)" json:"trade_in_value"`
	Condition       string          `json:"condition"`
	InspectionNotes string          `json:"inspection_notes,omitempty"`
	CoreCharge      decimal.Decimal `sql:"type:decimal(10,2)" json:"core_charge"`
	RefundAmount    decimal.Decimal `sql:"type:decimal(10,2)" json:"refund_amount"`
	Status          string          `gorm:"default:'pending'" json:"status"`
	ProcessedBy     *uuid.UUID      `gorm:"type:uuid" json:"processed_by,omitempty"`
	ProcessedAt     *time.Time      `json:"processed_at,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}
