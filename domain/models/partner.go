package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// Note: policy package import removed to break import cycle
	// Policy relationships are handled via PolicyID UUID only
)

// Partner represents third-party partners in the insurance ecosystem
type Partner struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PartnerCode string    `gorm:"uniqueIndex;not null" json:"partner_code"`
	PartnerType string    `gorm:"not null" json:"partner_type"`    // repair_technician, agent, broker, inspector, adjuster, vendor
	Status      string    `gorm:"default:'pending'" json:"status"` // pending, active, suspended, terminated

	// Personal/Business Information
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	BusinessName   string `json:"business_name"`
	Email          string `gorm:"uniqueIndex;not null" json:"email"`
	Phone          string `gorm:"not null" json:"phone"`
	AlternatePhone string `json:"alternate_phone"`

	// Authentication
	PasswordHash     string     `json:"-"`
	TwoFactorEnabled bool       `gorm:"default:false" json:"two_factor_enabled"`
	LastLoginAt      *time.Time `json:"last_login_at"`

	// Identification
	LicenseNumber string     `json:"license_number"`
	LicenseType   string     `json:"license_type"`
	LicenseExpiry *time.Time `json:"license_expiry"`
	TaxID         string     `json:"tax_id"`
	NationalID    string     `json:"national_id"`

	// Location
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`
	ServiceArea string `gorm:"type:json" json:"service_area"` // JSON array of postal codes or regions

	// Professional Details
	Specializations   string  `gorm:"type:json" json:"specializations"` // JSON array
	Certifications    string  `gorm:"type:json" json:"certifications"`  // JSON array
	YearsOfExperience int     `json:"years_of_experience"`
	Rating            float64 `gorm:"default:0" json:"rating"`
	ReviewCount       int     `gorm:"default:0" json:"review_count"`

	// Operational
	IsAvailable        bool   `gorm:"default:true" json:"is_available"`
	WorkingHours       string `gorm:"type:json" json:"working_hours"` // JSON object
	MaxAssignments     int    `gorm:"default:10" json:"max_assignments"`
	CurrentAssignments int    `gorm:"default:0" json:"current_assignments"`

	// Financial
	CommissionRate float64 `json:"commission_rate"`
	PaymentTerms   string  `json:"payment_terms"` // net30, net60, immediate
	BankAccount    string  `json:"-"`             // Encrypted
	PaymentMethod  string  `json:"payment_method"`

	// Compliance
	BackgroundCheck     bool       `gorm:"default:false" json:"background_check"`
	BackgroundCheckDate *time.Time `json:"background_check_date"`
	InsuranceCoverage   float64    `json:"insurance_coverage"`
	InsuranceExpiry     *time.Time `json:"insurance_expiry"`
	ContractStartDate   *time.Time `json:"contract_start_date"`
	ContractEndDate     *time.Time `json:"contract_end_date"`

	// Performance
	CompletedJobs        int     `gorm:"default:0" json:"completed_jobs"`
	SuccessRate          float64 `gorm:"default:0" json:"success_rate"`
	AverageResponseTime  int     `json:"average_response_time"` // minutes
	CustomerSatisfaction float64 `gorm:"default:0" json:"customer_satisfaction"`

	// Metadata
	Notes    string `json:"notes"`
	Tags     string `gorm:"type:json" json:"tags"` // JSON array
	Metadata string `gorm:"type:json" json:"metadata"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	RepairShopID *uuid.UUID          `gorm:"type:uuid" json:"repair_shop_id"`
	RepairShop   *RepairShop         `gorm:"foreignKey:RepairShopID" json:"repair_shop,omitempty"`
	Assignments  []PartnerAssignment `gorm:"foreignKey:PartnerID" json:"assignments,omitempty"`
	Commissions  []Commission        `gorm:"foreignKey:AgentID" json:"commissions,omitempty"`
}

// PartnerAssignment represents assignments given to partners
type PartnerAssignment struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AssignmentNumber string    `gorm:"uniqueIndex;not null" json:"assignment_number"`
	PartnerID        uuid.UUID `gorm:"type:uuid;not null" json:"partner_id"`
	AssignmentType   string    `gorm:"not null" json:"assignment_type"` // repair, inspection, sales, claim_assessment

	// Related Entities
	ClaimID         *uuid.UUID `gorm:"type:uuid" json:"claim_id"`
	PolicyID        *uuid.UUID `gorm:"type:uuid" json:"policy_id"`
	DeviceID        *uuid.UUID `gorm:"type:uuid" json:"device_id"`
	UserID          *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	RepairBookingID *uuid.UUID `gorm:"type:uuid" json:"repair_booking_id"`

	// Assignment Details
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Priority    string `gorm:"default:'normal'" json:"priority"` // low, normal, high, urgent
	Status      string `gorm:"default:'assigned'" json:"status"` // assigned, accepted, in_progress, completed, cancelled

	// Timeline
	AssignedAt  time.Time  `gorm:"not null" json:"assigned_at"`
	AcceptedAt  *time.Time `json:"accepted_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Deadline    *time.Time `json:"deadline"`

	// Location
	ServiceLocation string  `json:"service_location"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`

	// Financial
	EstimatedCost float64    `json:"estimated_cost"`
	ActualCost    float64    `json:"actual_cost"`
	PartnerFee    float64    `json:"partner_fee"`
	IsPaid        bool       `gorm:"default:false" json:"is_paid"`
	PaidAt        *time.Time `json:"paid_at"`

	// Documentation
	PreAssessment  string `gorm:"type:json" json:"pre_assessment"`
	PostAssessment string `gorm:"type:json" json:"post_assessment"`
	Photos         string `gorm:"type:json" json:"photos"`    // JSON array of photo URLs
	Documents      string `gorm:"type:json" json:"documents"` // JSON array of document IDs
	PartnerNotes   string `json:"partner_notes"`
	InternalNotes  string `json:"internal_notes"`

	// Quality
	QualityScore     float64 `json:"quality_score"`
	CustomerRating   float64 `json:"customer_rating"`
	CustomerFeedback string  `json:"customer_feedback"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Partner Partner `gorm:"foreignKey:PartnerID" json:"partner,omitempty"`
	Claim   *Claim  `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	// Policy relationship removed to break import cycle - use PolicyID to load separately
	Device        *Device        `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RepairBooking *RepairBooking `gorm:"foreignKey:RepairBookingID" json:"repair_booking,omitempty"`
}

// Inspector represents insurance inspectors for device verification
type Inspector struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PartnerID     uuid.UUID `gorm:"type:uuid;not null" json:"partner_id"`
	InspectorCode string    `gorm:"uniqueIndex;not null" json:"inspector_code"`

	// Certification
	CertificationBody   string    `json:"certification_body"`
	CertificationNumber string    `json:"certification_number"`
	CertificationLevel  string    `json:"certification_level"` // junior, senior, expert
	ValidUntil          time.Time `json:"valid_until"`

	// Specialization
	DeviceBrands    string `gorm:"type:json" json:"device_brands"`    // JSON array of brands
	InspectionTypes string `gorm:"type:json" json:"inspection_types"` // JSON array: pre_insurance, claim, quality_check

	// Performance
	InspectionsCompleted  int     `gorm:"default:0" json:"inspections_completed"`
	AccuracyRate          float64 `gorm:"default:0" json:"accuracy_rate"`
	AverageInspectionTime int     `json:"average_inspection_time"` // minutes
	FraudDetectionRate    float64 `json:"fraud_detection_rate"`

	// Availability
	MaxInspectionsPerDay int        `gorm:"default:5" json:"max_inspections_per_day"`
	CurrentQueue         int        `gorm:"default:0" json:"current_queue"`
	NextAvailable        *time.Time `json:"next_available"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Partner     Partner            `gorm:"foreignKey:PartnerID" json:"partner,omitempty"`
	Inspections []DeviceInspection `gorm:"foreignKey:InspectorID" json:"inspections,omitempty"`
}

// DeviceInspection represents a device inspection conducted by an inspector
type DeviceInspection struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	InspectionNumber string    `gorm:"uniqueIndex;not null" json:"inspection_number"`
	InspectorID      uuid.UUID `gorm:"type:uuid;not null" json:"inspector_id"`
	DeviceID         uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Inspection Type
	Type        string    `gorm:"not null" json:"type"` // pre_insurance, claim_verification, quality_check, renewal
	Purpose     string    `json:"purpose"`
	RequestedBy uuid.UUID `gorm:"type:uuid" json:"requested_by"`

	// Related Entities
	PolicyID *uuid.UUID `gorm:"type:uuid" json:"policy_id"`
	ClaimID  *uuid.UUID `gorm:"type:uuid" json:"claim_id"`

	// Inspection Details
	ScheduledDate  time.Time  `json:"scheduled_date"`
	InspectionDate *time.Time `json:"inspection_date"`
	Location       string     `json:"location"`                          // onsite, remote, shop
	Status         string     `gorm:"default:'scheduled'" json:"status"` // scheduled, in_progress, completed, failed

	// Device Condition Assessment
	PhysicalCondition    string `gorm:"type:json" json:"physical_condition"` // Detailed JSON assessment
	FunctionalTests      string `gorm:"type:json" json:"functional_tests"`   // JSON test results
	IMEIVerified         bool   `gorm:"default:false" json:"imei_verified"`
	AuthenticityVerified bool   `gorm:"default:false" json:"authenticity_verified"`

	// Results
	OverallGrade       string  `json:"overall_grade"` // A, B, C, D, F
	EstimatedValue     float64 `json:"estimated_value"`
	RepairRequired     bool    `gorm:"default:false" json:"repair_required"`
	RepairEstimate     float64 `json:"repair_estimate"`
	IsInsurable        bool    `gorm:"default:true" json:"is_insurable"`
	RecommendedPremium float64 `json:"recommended_premium"`

	// Documentation
	InspectionReport string `gorm:"type:text" json:"inspection_report"`
	Photos           string `gorm:"type:json" json:"photos"` // JSON array of photo URLs
	Videos           string `gorm:"type:json" json:"videos"` // JSON array of video URLs
	Signature        string `json:"signature"`               // Digital signature

	// Fraud Detection
	FraudIndicators string  `gorm:"type:json" json:"fraud_indicators"`
	FraudRiskScore  float64 `json:"fraud_risk_score"`
	RequiresReview  bool    `gorm:"default:false" json:"requires_review"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Inspector Inspector `gorm:"foreignKey:InspectorID" json:"inspector,omitempty"`
	Device    Device    `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// Policy relationship removed to break import cycle - use PolicyID to load separately
	Claim *Claim `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
}

// Vendor represents third-party vendors for parts and services
type Vendor struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	VendorCode string    `gorm:"uniqueIndex;not null" json:"vendor_code"`
	VendorName string    `gorm:"not null" json:"vendor_name"`
	VendorType string    `gorm:"not null" json:"vendor_type"` // parts_supplier, service_provider, logistics

	// Contact
	ContactPerson string `json:"contact_person"`
	Email         string `gorm:"not null" json:"email"`
	Phone         string `gorm:"not null" json:"phone"`
	Website       string `json:"website"`

	// Business Details
	BusinessLicense  string    `json:"business_license"`
	TaxID            string    `json:"tax_id"`
	RegistrationDate time.Time `json:"registration_date"`

	// Address
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`

	// Services
	ServicesOffered string `gorm:"type:json" json:"services_offered"` // JSON array
	ProductCatalog  string `gorm:"type:json" json:"product_catalog"`  // JSON catalog
	BrandsSupported string `gorm:"type:json" json:"brands_supported"` // JSON array

	// Financial
	PaymentTerms string  `json:"payment_terms"`
	CreditLimit  float64 `json:"credit_limit"`
	DiscountRate float64 `json:"discount_rate"`
	Currency     string  `gorm:"default:'USD'" json:"currency"`

	// Performance
	Rating          float64 `gorm:"default:0" json:"rating"`
	DeliveryTime    int     `json:"delivery_time"` // average days
	FulfillmentRate float64 `json:"fulfillment_rate"`
	QualityScore    float64 `json:"quality_score"`

	// Status
	Status            string     `gorm:"default:'active'" json:"status"` // active, suspended, terminated
	IsPreferred       bool       `gorm:"default:false" json:"is_preferred"`
	ContractStartDate *time.Time `json:"contract_start_date"`
	ContractEndDate   *time.Time `json:"contract_end_date"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	PurchaseOrders []VendorPurchaseOrder `gorm:"foreignKey:VendorID" json:"purchase_orders,omitempty"`
}

// VendorPurchaseOrder represents purchase orders to vendors
type VendorPurchaseOrder struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNumber string    `gorm:"uniqueIndex;not null" json:"order_number"`
	VendorID    uuid.UUID `gorm:"type:uuid;not null" json:"vendor_id"`

	// Related to
	RepairBookingID *uuid.UUID `gorm:"type:uuid" json:"repair_booking_id"`
	ClaimID         *uuid.UUID `gorm:"type:uuid" json:"claim_id"`

	// Order Details
	OrderType    string  `gorm:"not null" json:"order_type"`      // parts, service, replacement_device
	Items        string  `gorm:"type:json;not null" json:"items"` // JSON array of items
	Quantity     int     `gorm:"not null" json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalAmount  float64 `gorm:"not null" json:"total_amount"`
	Tax          float64 `json:"tax"`
	ShippingCost float64 `json:"shipping_cost"`

	// Status
	Status           string     `gorm:"default:'pending'" json:"status"` // pending, approved, shipped, delivered, cancelled
	OrderDate        time.Time  `gorm:"not null" json:"order_date"`
	ExpectedDelivery time.Time  `json:"expected_delivery"`
	ActualDelivery   *time.Time `json:"actual_delivery"`

	// Tracking
	TrackingNumber  string `json:"tracking_number"`
	ShippingCarrier string `json:"shipping_carrier"`

	// Payment
	PaymentStatus string     `gorm:"default:'pending'" json:"payment_status"` // pending, paid, overdue
	PaymentDate   *time.Time `json:"payment_date"`
	InvoiceNumber string     `json:"invoice_number"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Vendor        Vendor         `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	RepairBooking *RepairBooking `gorm:"foreignKey:RepairBookingID" json:"repair_booking,omitempty"`
	Claim         *Claim         `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
}

// Table names
func (Partner) TableName() string             { return "partners" }
func (PartnerAssignment) TableName() string   { return "partner_assignments" }
func (Inspector) TableName() string           { return "inspectors" }
func (DeviceInspection) TableName() string    { return "device_inspections" }
func (Vendor) TableName() string              { return "vendors" }
func (VendorPurchaseOrder) TableName() string { return "vendor_purchase_orders" }

// BeforeCreate hooks for UUID generation
func (p *Partner) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.PartnerCode == "" {
		p.PartnerCode = "PTR-" + uuid.New().String()[:8]
	}
	return nil
}

func (pa *PartnerAssignment) BeforeCreate(tx *gorm.DB) error {
	if pa.ID == uuid.Nil {
		pa.ID = uuid.New()
	}
	if pa.AssignmentNumber == "" {
		pa.AssignmentNumber = "ASG-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	pa.AssignedAt = time.Now()
	return nil
}

func (i *Inspector) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.InspectorCode == "" {
		i.InspectorCode = "INS-" + uuid.New().String()[:8]
	}
	return nil
}

func (di *DeviceInspection) BeforeCreate(tx *gorm.DB) error {
	if di.ID == uuid.Nil {
		di.ID = uuid.New()
	}
	if di.InspectionNumber == "" {
		di.InspectionNumber = "INSP-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	return nil
}

func (v *Vendor) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	if v.VendorCode == "" {
		v.VendorCode = "VND-" + uuid.New().String()[:8]
	}
	return nil
}

func (vpo *VendorPurchaseOrder) BeforeCreate(tx *gorm.DB) error {
	if vpo.ID == uuid.Nil {
		vpo.ID = uuid.New()
	}
	if vpo.OrderNumber == "" {
		vpo.OrderNumber = "PO-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	vpo.OrderDate = time.Now()
	return nil
}

// Business logic methods

// IsAvailableForWork checks if partner is available for new assignments
func (p *Partner) IsAvailableForWork() bool {
	return p.Status == "active" && p.IsAvailable && p.CurrentAssignments < p.MaxAssignments
}

// CanTakeAssignment checks if partner can take a specific assignment type
func (p *Partner) CanTakeAssignment(assignmentType string) bool {
	if !p.IsAvailableForWork() {
		return false
	}
	// Check if partner type matches assignment type
	switch assignmentType {
	case "repair":
		return p.PartnerType == "repair_technician"
	case "inspection":
		return p.PartnerType == "inspector"
	case "sales":
		return p.PartnerType == "agent" || p.PartnerType == "broker"
	case "claim_assessment":
		return p.PartnerType == "adjuster"
	default:
		return false
	}
}

// UpdateRating updates partner rating based on new feedback
func (p *Partner) UpdateRating(newRating float64) {
	if p.ReviewCount == 0 {
		p.Rating = newRating
	} else {
		// Calculate weighted average
		p.Rating = ((p.Rating * float64(p.ReviewCount)) + newRating) / float64(p.ReviewCount+1)
	}
	p.ReviewCount++
}

// IsLicenseValid checks if partner's license is valid
func (p *Partner) IsLicenseValid() bool {
	if p.LicenseExpiry == nil {
		return true // No expiry set
	}
	return time.Now().Before(*p.LicenseExpiry)
}

// AcceptAssignment accepts a partner assignment
func (pa *PartnerAssignment) AcceptAssignment() {
	pa.Status = "accepted"
	now := time.Now()
	pa.AcceptedAt = &now
}

// StartAssignment starts working on an assignment
func (pa *PartnerAssignment) StartAssignment() {
	pa.Status = "in_progress"
	now := time.Now()
	pa.StartedAt = &now
}

// CompleteAssignment completes an assignment
func (pa *PartnerAssignment) CompleteAssignment() {
	pa.Status = "completed"
	now := time.Now()
	pa.CompletedAt = &now
}

// IsOverdue checks if assignment is overdue
func (pa *PartnerAssignment) IsOverdue() bool {
	if pa.Deadline == nil || pa.Status == "completed" {
		return false
	}
	return time.Now().After(*pa.Deadline)
}

// CalculateInspectionScore calculates overall inspection score
func (di *DeviceInspection) CalculateInspectionScore() float64 {
	score := 100.0

	// Deduct points based on issues
	if !di.IMEIVerified {
		score -= 20
	}
	if !di.AuthenticityVerified {
		score -= 30
	}
	if di.RepairRequired {
		score -= 15
	}
	if di.FraudRiskScore > 0.5 {
		score -= 25
	}

	// Adjust based on grade
	gradeMultiplier := map[string]float64{
		"A": 1.0,
		"B": 0.85,
		"C": 0.7,
		"D": 0.5,
		"F": 0.2,
	}

	if multiplier, exists := gradeMultiplier[di.OverallGrade]; exists {
		score *= multiplier
	}

	if score < 0 {
		score = 0
	}

	return score
}

// IsHighRisk checks if inspection indicates high risk
func (di *DeviceInspection) IsHighRisk() bool {
	return di.FraudRiskScore > 0.7 || !di.AuthenticityVerified || di.OverallGrade == "F"
}

// CanFulfillOrder checks if vendor can fulfill an order
func (v *Vendor) CanFulfillOrder(orderType string, amount float64) bool {
	if v.Status != "active" {
		return false
	}
	if amount > v.CreditLimit {
		return false
	}
	// Check if vendor offers this service
	// This would need to parse the JSON services_offered field
	return true
}

// IsDeliveryOnTime checks if order was delivered on time
func (vpo *VendorPurchaseOrder) IsDeliveryOnTime() bool {
	if vpo.ActualDelivery == nil {
		return false
	}
	return vpo.ActualDelivery.Before(vpo.ExpectedDelivery) || vpo.ActualDelivery.Equal(vpo.ExpectedDelivery)
}
