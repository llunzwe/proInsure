package repair

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// RepairTechnician represents technicians working at repair shops
type RepairTechnician struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	RepairShopID       uuid.UUID      `gorm:"type:uuid;not null" json:"repair_shop_id"`
	EmployeeID         string         `gorm:"uniqueIndex;not null" json:"employee_id"`
	Name               string         `gorm:"not null" json:"name"`
	Email              string         `gorm:"not null" json:"email"`
	Phone              string         `json:"phone"`
	CertificationLevel string         `gorm:"default:'junior'" json:"certification_level"`
	Specializations    string         `json:"specializations"` // JSON array
	YearsExperience    int            `json:"years_experience"`
	TrainingCompleted  string         `json:"training_completed"` // JSON array
	Rating             float64        `gorm:"default:0" json:"rating"`
	CompletedRepairs   int            `gorm:"default:0" json:"completed_repairs"`
	SuccessRate        float64        `gorm:"default:0" json:"success_rate"`
	AverageTime        int            `json:"average_time_minutes"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	JoinedAt           time.Time      `json:"joined_at"`
	LastActiveAt       *time.Time     `json:"last_active_at,omitempty"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// RepairPartsInventory represents parts inventory for repair shops
type RepairPartsInventory struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	RepairShopID     uuid.UUID       `gorm:"type:uuid;not null" json:"repair_shop_id"`
	PartNumber       string          `gorm:"not null;index" json:"part_number"`
	PartName         string          `gorm:"not null" json:"part_name"`
	Category         string          `gorm:"not null" json:"category"`
	Brand            string          `json:"brand"`
	Quality          string          `gorm:"default:'oem'" json:"quality"`
	CompatibleModels string          `json:"compatible_models"` // JSON array
	QuantityInStock  int             `gorm:"default:0" json:"quantity_in_stock"`
	ReservedQuantity int             `gorm:"default:0" json:"reserved_quantity"`
	MinimumStock     int             `gorm:"default:5" json:"minimum_stock"`
	MaximumStock     int             `gorm:"default:50" json:"maximum_stock"`
	ReorderPoint     int             `gorm:"default:10" json:"reorder_point"`
	ReorderQuantity  int             `gorm:"default:20" json:"reorder_quantity"`
	UnitCost         decimal.Decimal `sql:"type:decimal(10,2)" json:"unit_cost"`
	SellingPrice     decimal.Decimal `sql:"type:decimal(10,2)" json:"selling_price"`
	SupplierID       string          `json:"supplier_id,omitempty"`
	SupplierName     string          `json:"supplier_name,omitempty"`
	LastOrderDate    *time.Time      `json:"last_order_date,omitempty"`
	ExpiryDate       *time.Time      `json:"expiry_date,omitempty"`
	Location         string          `json:"location"`
	IsActive         bool            `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}

// RepairWarranty represents warranty information for repairs
type RepairWarranty struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RepairBookingID uuid.UUID  `gorm:"type:uuid;not null" json:"repair_booking_id"`
	WarrantyNumber  string     `gorm:"uniqueIndex;not null" json:"warranty_number"`
	Type            string     `gorm:"not null" json:"type"`
	Coverage        string     `json:"coverage"` // JSON object
	StartDate       time.Time  `gorm:"not null" json:"start_date"`
	EndDate         time.Time  `gorm:"not null" json:"end_date"`
	Terms           string     `json:"terms"`
	Exclusions      string     `json:"exclusions"` // JSON array
	ClaimCount      int        `gorm:"default:0" json:"claim_count"`
	MaxClaims       int        `gorm:"default:1" json:"max_claims"`
	IsTransferable  bool       `gorm:"default:false" json:"is_transferable"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	VoidReason      string     `json:"void_reason,omitempty"`
	VoidDate        *time.Time `json:"void_date,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// RepairQualityCheck represents quality checks for completed repairs
type RepairQualityCheck struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	RepairBookingID  uuid.UUID `gorm:"type:uuid;not null" json:"repair_booking_id"`
	CheckedBy        uuid.UUID `gorm:"type:uuid;not null" json:"checked_by"`
	CheckType        string    `gorm:"not null" json:"check_type"`
	ChecklistItems   string    `json:"checklist_items"` // JSON array
	PassedItems      string    `json:"passed_items"`    // JSON array
	FailedItems      string    `json:"failed_items"`    // JSON array
	OverallResult    string    `gorm:"not null" json:"overall_result"`
	Score            float64   `json:"score"`
	Notes            string    `json:"notes"`
	Photos           string    `json:"photos"` // JSON array of URLs
	RequiresRework   bool      `gorm:"default:false" json:"requires_rework"`
	ReworkReason     string    `json:"rework_reason,omitempty"`
	CustomerApproved bool      `gorm:"default:false" json:"customer_approved"`
	CheckedAt        time.Time `gorm:"not null" json:"checked_at"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// RepairDiagnostic represents diagnostic results for devices
type RepairDiagnostic struct {
	ID                 uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	RepairBookingID    uuid.UUID       `gorm:"type:uuid;not null" json:"repair_booking_id"`
	TechnicianID       uuid.UUID       `gorm:"type:uuid;not null" json:"technician_id"`
	DiagnosticType     string          `gorm:"not null" json:"diagnostic_type"`
	TestsPerformed     string          `json:"tests_performed"`     // JSON array
	TestResults        string          `json:"test_results"`        // JSON object
	IssuesFound        string          `json:"issues_found"`        // JSON array
	RecommendedRepairs string          `json:"recommended_repairs"` // JSON array
	EstimatedCost      decimal.Decimal `sql:"type:decimal(10,2)" json:"estimated_cost"`
	EstimatedTime      int             `json:"estimated_time_minutes"`
	Severity           string          `gorm:"not null" json:"severity"`
	CanBeRepaired      bool            `gorm:"default:true" json:"can_be_repaired"`
	RequiresParts      bool            `gorm:"default:false" json:"requires_parts"`
	PartsNeeded        string          `json:"parts_needed,omitempty"` // JSON array
	CustomerApproval   string          `gorm:"default:'pending'" json:"customer_approval"`
	ApprovedAt         *time.Time      `json:"approved_at,omitempty"`
	DiagnosticFee      decimal.Decimal `sql:"type:decimal(10,2)" json:"diagnostic_fee"`
	FeeWaived          bool            `gorm:"default:false" json:"fee_waived"`
	CompletedAt        time.Time       `json:"completed_at"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// RepairCommunication represents communication logs for repairs
type RepairCommunication struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RepairBookingID uuid.UUID  `gorm:"type:uuid;not null" json:"repair_booking_id"`
	Type            string     `gorm:"not null" json:"type"`      // email, sms, call, in_app
	Direction       string     `gorm:"not null" json:"direction"` // inbound, outbound
	Subject         string     `json:"subject,omitempty"`
	Message         string     `gorm:"not null" json:"message"`
	Recipient       string     `gorm:"not null" json:"recipient"`
	Sender          string     `gorm:"not null" json:"sender"`
	Status          string     `gorm:"not null" json:"status"`
	DeliveredAt     *time.Time `json:"delivered_at,omitempty"`
	ReadAt          *time.Time `json:"read_at,omitempty"`
	ResponseTo      *uuid.UUID `gorm:"type:uuid" json:"response_to,omitempty"`
	Attachments     string     `json:"attachments,omitempty"` // JSON array
	Metadata        string     `json:"metadata,omitempty"`    // JSON object
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// RepairCostBreakdown represents detailed cost breakdown for repairs
type RepairCostBreakdown struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	RepairBookingID  uuid.UUID       `gorm:"type:uuid;not null" json:"repair_booking_id"`
	LaborCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"labor_cost"`
	PartsCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"parts_cost"`
	DiagnosticCost   decimal.Decimal `sql:"type:decimal(10,2)" json:"diagnostic_cost"`
	ShippingCost     decimal.Decimal `sql:"type:decimal(10,2)" json:"shipping_cost"`
	ExpressFee       decimal.Decimal `sql:"type:decimal(10,2)" json:"express_fee"`
	ServiceFee       decimal.Decimal `sql:"type:decimal(10,2)" json:"service_fee"`
	Tax              decimal.Decimal `sql:"type:decimal(10,2)" json:"tax"`
	Discount         decimal.Decimal `sql:"type:decimal(10,2)" json:"discount"`
	DiscountReason   string          `json:"discount_reason,omitempty"`
	InsuranceCovered decimal.Decimal `sql:"type:decimal(10,2)" json:"insurance_covered"`
	CustomerPays     decimal.Decimal `sql:"type:decimal(10,2)" json:"customer_pays"`
	TotalCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"total_cost"`
	Currency         string          `gorm:"default:'USD'" json:"currency"`
	PaymentStatus    string          `gorm:"default:'pending'" json:"payment_status"`
	PaymentMethod    string          `json:"payment_method,omitempty"`
	InvoiceNumber    string          `json:"invoice_number,omitempty"`
	InvoiceDate      *time.Time      `json:"invoice_date,omitempty"`
	PaidAt           *time.Time      `json:"paid_at,omitempty"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// RepairSLA represents SLA agreements for repairs
type RepairSLA struct {
	ID                     uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	RepairShopID           uuid.UUID       `gorm:"type:uuid;not null" json:"repair_shop_id"`
	SLALevel               string          `gorm:"not null" json:"sla_level"`
	ResponseTime           int             `gorm:"not null" json:"response_time_minutes"`
	DiagnosticTime         int             `gorm:"not null" json:"diagnostic_time_hours"`
	RepairTime             int             `gorm:"not null" json:"repair_time_hours"`
	ResolutionTime         int             `gorm:"not null" json:"resolution_time_hours"`
	UpdateFrequency        int             `gorm:"not null" json:"update_frequency_hours"`
	CompensationPerHour    decimal.Decimal `sql:"type:decimal(10,2)" json:"compensation_per_hour"`
	MaxCompensation        decimal.Decimal `sql:"type:decimal(10,2)" json:"max_compensation"`
	EscalationThreshold    int             `json:"escalation_threshold_hours"`
	PrioritySupport        bool            `gorm:"default:false" json:"priority_support"`
	DedicatedTechnician    bool            `gorm:"default:false" json:"dedicated_technician"`
	GuaranteedAvailability bool            `gorm:"default:false" json:"guaranteed_availability"`
	IsActive               bool            `gorm:"default:true" json:"is_active"`
	CreatedAt              time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// RepairPerformanceMetrics represents performance metrics for repair shops
type RepairPerformanceMetrics struct {
	ID                    uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	RepairShopID          uuid.UUID       `gorm:"type:uuid;not null" json:"repair_shop_id"`
	Period                string          `gorm:"not null" json:"period"` // daily, weekly, monthly
	PeriodStart           time.Time       `gorm:"not null" json:"period_start"`
	PeriodEnd             time.Time       `gorm:"not null" json:"period_end"`
	TotalRepairs          int             `json:"total_repairs"`
	CompletedRepairs      int             `json:"completed_repairs"`
	CancelledRepairs      int             `json:"cancelled_repairs"`
	AverageRepairTime     int             `json:"average_repair_time_hours"`
	FirstTimeFixRate      float64         `json:"first_time_fix_rate"`
	CustomerSatisfaction  float64         `json:"customer_satisfaction"`
	Revenue               decimal.Decimal `sql:"type:decimal(10,2)" json:"revenue"`
	PartsCost             decimal.Decimal `sql:"type:decimal(10,2)" json:"parts_cost"`
	LaborCost             decimal.Decimal `sql:"type:decimal(10,2)" json:"labor_cost"`
	ProfitMargin          float64         `json:"profit_margin"`
	TechnicianUtilization float64         `json:"technician_utilization"`
	SLACompliance         float64         `json:"sla_compliance"`
	QualityScore          float64         `json:"quality_score"`
	CreatedAt             time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

// RepairTraining represents training programs for repair shops
type RepairTraining struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	RepairShopID      *uuid.UUID      `gorm:"type:uuid" json:"repair_shop_id,omitempty"`
	TechnicianID      *uuid.UUID      `gorm:"type:uuid" json:"technician_id,omitempty"`
	TrainingName      string          `gorm:"not null" json:"training_name"`
	TrainingType      string          `gorm:"not null" json:"training_type"`
	Provider          string          `json:"provider"`
	Description       string          `json:"description"`
	Duration          int             `json:"duration_hours"`
	StartDate         time.Time       `json:"start_date"`
	EndDate           time.Time       `json:"end_date"`
	Status            string          `gorm:"default:'scheduled'" json:"status"`
	CompletionRate    float64         `json:"completion_rate"`
	TestScore         float64         `json:"test_score,omitempty"`
	CertificateNumber string          `json:"certificate_number,omitempty"`
	CertificateURL    string          `json:"certificate_url,omitempty"`
	ExpiryDate        *time.Time      `json:"expiry_date,omitempty"`
	Cost              decimal.Decimal `sql:"type:decimal(10,2)" json:"cost"`
	IsMandatory       bool            `gorm:"default:false" json:"is_mandatory"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// RepairSupplier represents parts suppliers for repair shops
type RepairSupplier struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	SupplierCode  string          `gorm:"uniqueIndex;not null" json:"supplier_code"`
	Name          string          `gorm:"not null" json:"name"`
	Type          string          `gorm:"not null" json:"type"`
	Categories    string          `json:"categories"` // JSON array
	Brands        string          `json:"brands"`     // JSON array
	ContactName   string          `json:"contact_name"`
	ContactEmail  string          `json:"contact_email"`
	ContactPhone  string          `json:"contact_phone"`
	Address       string          `json:"address"`
	PaymentTerms  string          `json:"payment_terms"`
	MinimumOrder  decimal.Decimal `sql:"type:decimal(10,2)" json:"minimum_order"`
	LeadTime      int             `json:"lead_time_days"`
	Rating        float64         `gorm:"default:0" json:"rating"`
	QualityScore  float64         `gorm:"default:0" json:"quality_score"`
	DeliveryScore float64         `gorm:"default:0" json:"delivery_score"`
	PricingScore  float64         `gorm:"default:0" json:"pricing_score"`
	IsPreferred   bool            `gorm:"default:false" json:"is_preferred"`
	IsActive      bool            `gorm:"default:true" json:"is_active"`
	ContractStart *time.Time      `json:"contract_start,omitempty"`
	ContractEnd   *time.Time      `json:"contract_end,omitempty"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

// RepairFeedback represents customer feedback for repairs
type RepairFeedback struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RepairBookingID   uuid.UUID  `gorm:"type:uuid;not null" json:"repair_booking_id"`
	UserID            uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	FeedbackType      string     `gorm:"not null" json:"feedback_type"`
	SatisfactionLevel string     `gorm:"not null" json:"satisfaction_level"`
	QualityRating     int        `json:"quality_rating"`
	SpeedRating       int        `json:"speed_rating"`
	ServiceRating     int        `json:"service_rating"`
	ValueRating       int        `json:"value_rating"`
	LikelyToRecommend int        `json:"likely_to_recommend"`
	Comments          string     `json:"comments"`
	Improvements      string     `json:"improvements"`
	ComplaintRaised   bool       `gorm:"default:false" json:"complaint_raised"`
	ComplaintDetails  string     `json:"complaint_details,omitempty"`
	ResponseProvided  bool       `gorm:"default:false" json:"response_provided"`
	ResponseText      string     `json:"response_text,omitempty"`
	RespondedAt       *time.Time `json:"responded_at,omitempty"`
	IsFeatured        bool       `gorm:"default:false" json:"is_featured"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
