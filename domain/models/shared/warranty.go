package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// Warranty represents device warranties (manufacturer or extended)
type Warranty struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WarrantyCode string    `gorm:"uniqueIndex;not null" json:"warranty_code"`

	// Device Information
	DeviceID uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Warranty Type
	Type         string `gorm:"not null" json:"type"` // manufacturer, extended, third_party, insurance_add_on
	Provider     string `gorm:"not null" json:"provider"`
	ProviderType string `json:"provider_type"` // manufacturer, retailer, insurance, third_party

	// Coverage Details
	CoverageName        string `gorm:"not null" json:"coverage_name"`
	CoverageDescription string `json:"coverage_description"`
	CoverageComponents  string `gorm:"type:json" json:"coverage_components"` // JSON array of covered components
	Exclusions          string `gorm:"type:json" json:"exclusions"`          // JSON array of exclusions

	// Duration
	StartDate        time.Time `gorm:"not null" json:"start_date"`
	EndDate          time.Time `gorm:"not null" json:"end_date"`
	OriginalDuration int       `json:"original_duration"` // months
	ExtendedDuration int       `json:"extended_duration"` // additional months if extended

	// Registration
	RegistrationNumber string    `json:"registration_number"`
	RegistrationDate   time.Time `json:"registration_date"`
	PurchaseDate       time.Time `json:"purchase_date"`
	PurchasePrice      float64   `json:"purchase_price"`
	ReceiptNumber      string    `json:"receipt_number"`

	// Status
	Status         string `gorm:"default:'active'" json:"status"` // active, expired, cancelled, claimed, suspended
	IsTransferable bool   `gorm:"default:false" json:"is_transferable"`
	TransferCount  int    `gorm:"default:0" json:"transfer_count"`
	MaxTransfers   int    `gorm:"default:1" json:"max_transfers"`

	// Claims
	ClaimsAllowed    int     `gorm:"default:3" json:"claims_allowed"`
	ClaimsUsed       int     `gorm:"default:0" json:"claims_used"`
	TotalClaimAmount float64 `gorm:"default:0" json:"total_claim_amount"`
	ClaimLimit       float64 `json:"claim_limit"`  // Max claim amount per incident
	AnnualLimit      float64 `json:"annual_limit"` // Max total claims per year

	// Financial
	PremiumAmount    float64 `json:"premium_amount"`
	PaymentFrequency string  `json:"payment_frequency"` // one_time, monthly, annual
	Deductible       float64 `json:"deductible"`
	IsPaid           bool    `gorm:"default:false" json:"is_paid"`

	// Terms & Conditions
	TermsVersion       string     `json:"terms_version"`
	TermsAcceptedAt    time.Time  `json:"terms_accepted_at"`
	RequiresInspection bool       `gorm:"default:false" json:"requires_inspection"`
	LastInspectionDate *time.Time `json:"last_inspection_date"`

	// Contact Information
	ServicePhone   string `json:"service_phone"`
	ServiceEmail   string `json:"service_email"`
	ServiceWebsite string `json:"service_website"`

	// Integration
	ExternalID string     `json:"external_id"` // ID in provider's system
	SyncedAt   *time.Time `json:"synced_at"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Device            models.Device      `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User              models.User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	WarrantyClaims    []WarrantyClaim    `gorm:"foreignKey:WarrantyID" json:"warranty_claims,omitempty"`
	ExtendedCoverages []ExtendedCoverage `gorm:"foreignKey:BaseWarrantyID" json:"extended_coverages,omitempty"`
}

// WarrantyClaim represents claims made against a warranty
type WarrantyClaim struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ClaimNumber string    `gorm:"uniqueIndex;not null" json:"claim_number"`
	WarrantyID  uuid.UUID `gorm:"type:uuid;not null" json:"warranty_id"`

	// Related Insurance Claim (if applicable)
	InsuranceClaimID *uuid.UUID `gorm:"type:uuid" json:"insurance_claim_id"`

	// Claim Details
	ClaimType          string `gorm:"not null" json:"claim_type"` // repair, replacement, refund
	IssueDescription   string `gorm:"not null" json:"issue_description"`
	IssueCategory      string `json:"issue_category"`                       // hardware, software, accidental_damage
	ComponentsAffected string `gorm:"type:json" json:"components_affected"` // JSON array

	// Status
	Status   string `gorm:"default:'submitted'" json:"status"` // submitted, under_review, approved, rejected, in_repair, completed
	Priority string `gorm:"default:'normal'" json:"priority"`  // low, normal, high, urgent

	// Dates
	SubmittedAt         time.Time  `gorm:"not null" json:"submitted_at"`
	ReviewedAt          *time.Time `json:"reviewed_at"`
	ApprovedAt          *time.Time `json:"approved_at"`
	CompletedAt         *time.Time `json:"completed_at"`
	EstimatedCompletion *time.Time `json:"estimated_completion"`

	// Assessment
	AssessmentNotes    string     `json:"assessment_notes"`
	RequiresInspection bool       `gorm:"default:false" json:"requires_inspection"`
	InspectionDate     *time.Time `json:"inspection_date"`
	InspectionReport   string     `json:"inspection_report"`

	// Resolution
	ResolutionType      string     `json:"resolution_type"` // repair, replace, refund, deny
	ResolutionNotes     string     `json:"resolution_notes"`
	ReplacementDeviceID *uuid.UUID `gorm:"type:uuid" json:"replacement_device_id"`
	RefundAmount        float64    `json:"refund_amount"`

	// Service Details
	ServiceCenterID     *uuid.UUID `gorm:"type:uuid" json:"service_center_id"`
	TechnicianID        *uuid.UUID `gorm:"type:uuid" json:"technician_id"`
	ServiceTicketNumber string     `json:"service_ticket_number"`

	// Financial
	EstimatedCost     float64 `json:"estimated_cost"`
	ActualCost        float64 `json:"actual_cost"`
	DeductibleCharged float64 `json:"deductible_charged"`
	CustomerPaid      float64 `json:"customer_paid"`

	// Documentation
	Documents string `gorm:"type:json" json:"documents"` // JSON array of document IDs
	Photos    string `gorm:"type:json" json:"photos"`    // JSON array of photo URLs

	// Denial (if rejected)
	DenialReason    string     `json:"denial_reason"`
	AppealDeadline  *time.Time `json:"appeal_deadline"`
	AppealSubmitted bool       `gorm:"default:false" json:"appeal_submitted"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Warranty       Warranty      `gorm:"foreignKey:WarrantyID" json:"warranty,omitempty"`
	InsuranceClaim *models.Claim `gorm:"foreignKey:InsuranceClaimID" json:"insurance_claim,omitempty"`
}

// ExtendedCoverage represents additional coverage options
type ExtendedCoverage struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CoverageCode string    `gorm:"uniqueIndex;not null" json:"coverage_code"`

	// Base Warranty
	BaseWarrantyID *uuid.UUID `gorm:"type:uuid" json:"base_warranty_id"`
	PolicyID       *uuid.UUID `gorm:"type:uuid" json:"policy_id"`

	// Coverage Details
	Name        string `gorm:"not null" json:"name"`
	Type        string `gorm:"not null" json:"type"` // accidental_damage, theft, loss, water_damage, screen_protection
	Description string `json:"description"`

	// Coverage Scope
	CoverageItems string `gorm:"type:json" json:"coverage_items"` // JSON array of covered items
	Exclusions    string `gorm:"type:json" json:"exclusions"`     // JSON array
	Territories   string `gorm:"type:json" json:"territories"`    // JSON array of covered territories

	// Duration
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	WaitingPeriod int       `json:"waiting_period"` // days before coverage starts

	// Limits
	IncidentLimit    float64 `json:"incident_limit"`     // Max per incident
	AnnualLimit      float64 `json:"annual_limit"`       // Max per year
	LifetimeLimit    float64 `json:"lifetime_limit"`     // Max total
	IncidentsPerYear int     `json:"incidents_per_year"` // Max number of incidents

	// Financial
	PremiumAmount    float64 `gorm:"not null" json:"premium_amount"`
	PremiumFrequency string  `json:"premium_frequency"` // monthly, quarterly, annual
	Deductible       float64 `json:"deductible"`
	CoPayPercentage  float64 `json:"co_pay_percentage"`

	// Claims
	ClaimsUsed       int     `gorm:"default:0" json:"claims_used"`
	TotalClaimAmount float64 `gorm:"default:0" json:"total_claim_amount"`

	// Status
	Status      string     `gorm:"default:'active'" json:"status"` // active, expired, cancelled, exhausted
	AutoRenew   bool       `gorm:"default:false" json:"auto_renew"`
	RenewalDate *time.Time `json:"renewal_date"`

	// Terms
	TermsVersion    string    `json:"terms_version"`
	TermsAcceptedAt time.Time `json:"terms_accepted_at"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	BaseWarranty   *Warranty       `gorm:"foreignKey:BaseWarrantyID" json:"base_warranty,omitempty"`
	Policy         *models.Policy  `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	CoverageClaims []CoverageClaim `gorm:"foreignKey:ExtendedCoverageID" json:"coverage_claims,omitempty"`
}

// CoverageClaim represents claims against extended coverage
type CoverageClaim struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ClaimNumber        string    `gorm:"uniqueIndex;not null" json:"claim_number"`
	ExtendedCoverageID uuid.UUID `gorm:"type:uuid;not null" json:"extended_coverage_id"`

	// Link to main insurance claim
	InsuranceClaimID *uuid.UUID `gorm:"type:uuid" json:"insurance_claim_id"`

	// Incident Details
	IncidentDate time.Time `gorm:"not null" json:"incident_date"`
	IncidentType string    `gorm:"not null" json:"incident_type"`
	Description  string    `gorm:"not null" json:"description"`
	Location     string    `json:"location"`

	// Claim Amount
	ClaimedAmount     float64 `gorm:"not null" json:"claimed_amount"`
	ApprovedAmount    float64 `json:"approved_amount"`
	DeductibleApplied float64 `json:"deductible_applied"`
	CoPayAmount       float64 `json:"co_pay_amount"`
	PaidAmount        float64 `json:"paid_amount"`

	// Status
	Status      string `gorm:"default:'submitted'" json:"status"` // submitted, under_review, approved, rejected, paid
	ReviewNotes string `json:"review_notes"`

	// Processing Dates
	SubmittedAt time.Time  `json:"submitted_at"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	ApprovedAt  *time.Time `json:"approved_at"`
	PaidAt      *time.Time `json:"paid_at"`

	// Documentation
	RequiredDocuments  string `gorm:"type:json" json:"required_documents"`  // JSON array
	SubmittedDocuments string `gorm:"type:json" json:"submitted_documents"` // JSON array

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ExtendedCoverage ExtendedCoverage `gorm:"foreignKey:ExtendedCoverageID" json:"extended_coverage,omitempty"`
	InsuranceClaim   *models.Claim    `gorm:"foreignKey:InsuranceClaimID" json:"insurance_claim,omitempty"`
}

// WarrantyProvider represents warranty service providers
type WarrantyProvider struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProviderCode string    `gorm:"uniqueIndex;not null" json:"provider_code"`
	Name         string    `gorm:"not null" json:"name"`
	Type         string    `gorm:"not null" json:"type"` // manufacturer, third_party, insurance

	// Contact Information
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`

	// Service Details
	ServicesOffered string `gorm:"type:json" json:"services_offered"` // JSON array
	SupportedBrands string `gorm:"type:json" json:"supported_brands"` // JSON array
	ServiceCenters  string `gorm:"type:json" json:"service_centers"`  // JSON array of locations

	// Integration
	APIEndpoint     string `json:"api_endpoint"`
	APIKey          string `json:"-"`                // Encrypted
	IntegrationType string `json:"integration_type"` // api, manual, email

	// Performance
	Rating               float64 `gorm:"default:0" json:"rating"`
	AverageResponseTime  int     `json:"average_response_time"` // hours
	ClaimApprovalRate    float64 `json:"claim_approval_rate"`
	CustomerSatisfaction float64 `json:"customer_satisfaction"`

	// Contract
	ContractStartDate time.Time  `json:"contract_start_date"`
	ContractEndDate   *time.Time `json:"contract_end_date"`
	CommissionRate    float64    `json:"commission_rate"`

	// Status
	Status      string `gorm:"default:'active'" json:"status"` // active, suspended, terminated
	IsPreferred bool   `gorm:"default:false" json:"is_preferred"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Warranties []Warranty `gorm:"foreignKey:Provider" json:"warranties,omitempty"`
}

// Table names
func (Warranty) TableName() string         { return "warranties" }
func (WarrantyClaim) TableName() string    { return "warranty_claims" }
func (ExtendedCoverage) TableName() string { return "extended_coverages" }
func (CoverageClaim) TableName() string    { return "coverage_claims" }
func (WarrantyProvider) TableName() string { return "warranty_providers" }

// BeforeCreate hooks
func (w *Warranty) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	if w.WarrantyCode == "" {
		w.WarrantyCode = "WAR-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	w.RegistrationDate = time.Now()
	return nil
}

func (wc *WarrantyClaim) BeforeCreate(tx *gorm.DB) error {
	if wc.ID == uuid.Nil {
		wc.ID = uuid.New()
	}
	if wc.ClaimNumber == "" {
		wc.ClaimNumber = "WC-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	wc.SubmittedAt = time.Now()
	return nil
}

func (ec *ExtendedCoverage) BeforeCreate(tx *gorm.DB) error {
	if ec.ID == uuid.Nil {
		ec.ID = uuid.New()
	}
	if ec.CoverageCode == "" {
		ec.CoverageCode = "EXT-" + uuid.New().String()[:8]
	}
	return nil
}

func (cc *CoverageClaim) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	if cc.ClaimNumber == "" {
		cc.ClaimNumber = "CC-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:6]
	}
	cc.SubmittedAt = time.Now()
	return nil
}

func (wp *WarrantyProvider) BeforeCreate(tx *gorm.DB) error {
	if wp.ID == uuid.Nil {
		wp.ID = uuid.New()
	}
	if wp.ProviderCode == "" {
		wp.ProviderCode = "WP-" + uuid.New().String()[:8]
	}
	return nil
}

// Business Logic Methods

// IsActive checks if warranty is currently active
func (w *Warranty) IsActive() bool {
	now := time.Now()
	return w.Status == "active" && now.After(w.StartDate) && now.Before(w.EndDate)
}

// IsExpired checks if warranty has expired
func (w *Warranty) IsExpired() bool {
	return time.Now().After(w.EndDate) || w.Status == "expired"
}

// CanClaim checks if a claim can be made
func (w *Warranty) CanClaim() bool {
	if !w.IsActive() {
		return false
	}
	if w.ClaimsUsed >= w.ClaimsAllowed {
		return false
	}
	return true
}

// DaysRemaining calculates days remaining in warranty
func (w *Warranty) DaysRemaining() int {
	if w.IsExpired() {
		return 0
	}
	duration := w.EndDate.Sub(time.Now())
	return int(duration.Hours() / 24)
}

// UseClaim increments the claims used counter
func (w *Warranty) UseClaim(amount float64) {
	w.ClaimsUsed++
	w.TotalClaimAmount += amount
}

// Transfer transfers warranty to a new user
func (w *Warranty) Transfer(newUserID uuid.UUID) error {
	if !w.IsTransferable {
		return gorm.ErrInvalidData
	}
	if w.TransferCount >= w.MaxTransfers {
		return gorm.ErrInvalidData
	}
	w.UserID = newUserID
	w.TransferCount++
	return nil
}

// Extend extends the warranty period
func (w *Warranty) Extend(months int) {
	w.ExtendedDuration += months
	w.EndDate = w.EndDate.AddDate(0, months, 0)
}

// ApproveClaim approves a warranty claim
func (wc *WarrantyClaim) ApproveClaim() {
	wc.Status = "approved"
	now := time.Now()
	wc.ApprovedAt = &now
}

// RejectClaim rejects a warranty claim
func (wc *WarrantyClaim) RejectClaim(reason string) {
	wc.Status = "rejected"
	wc.DenialReason = reason
	deadline := time.Now().AddDate(0, 0, 30) // 30 days to appeal
	wc.AppealDeadline = &deadline
}

// CompleteClaim marks claim as completed
func (wc *WarrantyClaim) CompleteClaim() {
	wc.Status = "completed"
	now := time.Now()
	wc.CompletedAt = &now
}

// CalculateCustomerCost calculates what customer needs to pay
func (wc *WarrantyClaim) CalculateCustomerCost() float64 {
	return wc.DeductibleCharged + (wc.ActualCost - wc.RefundAmount)
}

// IsActive checks if extended coverage is active
func (ec *ExtendedCoverage) IsActive() bool {
	now := time.Now()
	waitingPeriodEnd := ec.StartDate.AddDate(0, 0, ec.WaitingPeriod)
	return ec.Status == "active" && now.After(waitingPeriodEnd) && now.Before(ec.EndDate)
}

// CanClaim checks if a claim can be made against coverage
func (ec *ExtendedCoverage) CanClaim() bool {
	if !ec.IsActive() {
		return false
	}
	if ec.IncidentsPerYear > 0 && ec.ClaimsUsed >= ec.IncidentsPerYear {
		return false
	}
	if ec.LifetimeLimit > 0 && ec.TotalClaimAmount >= ec.LifetimeLimit {
		return false
	}
	return true
}

// CalculatePayout calculates the payout for a coverage claim
func (ec *ExtendedCoverage) CalculatePayout(claimedAmount float64) (payout float64, deductible float64, coPay float64) {
	deductible = ec.Deductible
	remainingAmount := claimedAmount - deductible

	if remainingAmount <= 0 {
		return 0, deductible, 0
	}

	// Apply incident limit
	if ec.IncidentLimit > 0 && remainingAmount > ec.IncidentLimit {
		remainingAmount = ec.IncidentLimit
	}

	// Apply co-pay
	coPay = remainingAmount * (ec.CoPayPercentage / 100)
	payout = remainingAmount - coPay

	// Check against annual limit
	if ec.AnnualLimit > 0 {
		yearClaims := ec.TotalClaimAmount // This would need to filter by current year
		if yearClaims+payout > ec.AnnualLimit {
			payout = ec.AnnualLimit - yearClaims
		}
	}

	return payout, deductible, coPay
}

// ShouldRenew checks if coverage should be renewed
func (ec *ExtendedCoverage) ShouldRenew() bool {
	if !ec.AutoRenew || ec.Status != "active" {
		return false
	}
	daysUntilExpiry := int(ec.EndDate.Sub(time.Now()).Hours() / 24)
	return daysUntilExpiry <= 30 // Renew if 30 days or less until expiry
}

// ApproveClaim approves a coverage claim
func (cc *CoverageClaim) ApproveClaim(amount float64) {
	cc.Status = "approved"
	cc.ApprovedAmount = amount
	now := time.Now()
	cc.ApprovedAt = &now
}

// PayClaim marks claim as paid
func (cc *CoverageClaim) PayClaim() {
	cc.Status = "paid"
	now := time.Now()
	cc.PaidAt = &now
}

// IsOverdue checks if provider response is overdue
func (wp *WarrantyProvider) IsOverdue(hours int) bool {
	return wp.AverageResponseTime > hours
}
