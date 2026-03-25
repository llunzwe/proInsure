package repair

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ShopIdentification contains basic identification info for repair shops
type ShopIdentification struct {
	ShopNumber       string     `gorm:"uniqueIndex;not null" json:"shop_number"`
	BusinessLicense  string     `gorm:"uniqueIndex;not null" json:"business_license"`
	TaxID            string     `gorm:"uniqueIndex" json:"tax_id"`
	RegistrationDate time.Time  `json:"registration_date"`
	ShopType         string     `gorm:"not null;default:'partner'" json:"shop_type"`
	ChainID          *uuid.UUID `gorm:"type:uuid" json:"chain_id,omitempty"`
	FranchiseID      *uuid.UUID `gorm:"type:uuid" json:"franchise_id,omitempty"`
	ParentShopID     *uuid.UUID `gorm:"type:uuid" json:"parent_shop_id,omitempty"`
	ExternalID       string     `json:"external_id,omitempty"`
}

// ShopLocation contains location and service area info
type ShopLocation struct {
	Address          string  `gorm:"not null" json:"address"`
	AddressLine2     string  `json:"address_line_2,omitempty"`
	City             string  `gorm:"not null" json:"city"`
	State            string  `gorm:"not null" json:"state"`
	Country          string  `gorm:"not null" json:"country"`
	PostalCode       string  `json:"postal_code"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	GeohashCode      string  `gorm:"index" json:"geohash_code,omitempty"`
	ServiceRadius    int     `gorm:"default:20" json:"service_radius_km"`
	ServiceAreas     string  `json:"service_areas"` // JSON array of postal codes
	IsMobileService  bool    `gorm:"default:false" json:"is_mobile_service"`
	HasStorefront    bool    `gorm:"default:true" json:"has_storefront"`
	ParkingAvailable bool    `gorm:"default:false" json:"parking_available"`
	WheelchairAccess bool    `gorm:"default:false" json:"wheelchair_access"`
	PublicTransport  string  `json:"public_transport,omitempty"`
}

// ShopContact contains contact information
type ShopContact struct {
	PrimaryPhone     string `gorm:"not null" json:"primary_phone"`
	SecondaryPhone   string `json:"secondary_phone,omitempty"`
	EmergencyPhone   string `json:"emergency_phone,omitempty"`
	Email            string `gorm:"not null" json:"email"`
	SupportEmail     string `json:"support_email,omitempty"`
	Website          string `json:"website,omitempty"`
	BookingURL       string `json:"booking_url,omitempty"`
	ContactPerson    string `json:"contact_person"`
	ContactTitle     string `json:"contact_title,omitempty"`
	PreferredContact string `gorm:"default:'email'" json:"preferred_contact"`
	Languages        string `json:"languages"` // JSON array
}

// ShopOperations contains operational details
type ShopOperations struct {
	OperatingHours   string     `json:"operating_hours"`  // JSON object
	HolidaySchedule  string     `json:"holiday_schedule"` // JSON array
	TemporaryClosure bool       `gorm:"default:false" json:"temporary_closure"`
	ClosureReason    string     `json:"closure_reason,omitempty"`
	ReopenDate       *time.Time `json:"reopen_date,omitempty"`
	CapacityPerDay   int        `gorm:"default:10" json:"capacity_per_day"`
	CurrentCapacity  int        `json:"current_capacity"`
	QueueLength      int        `json:"queue_length"`
	AverageWaitTime  int        `json:"average_wait_time_minutes"`
	ExpressService   bool       `gorm:"default:false" json:"express_service"`
	AppointmentOnly  bool       `gorm:"default:false" json:"appointment_only"`
	WalkInAccepted   bool       `gorm:"default:true" json:"walk_in_accepted"`
	OnlineBooking    bool       `gorm:"default:true" json:"online_booking"`
	HomeService      bool       `gorm:"default:false" json:"home_service"`
}

// ShopCapabilities contains service capabilities
type ShopCapabilities struct {
	Services          string          `json:"services"`          // JSON array of services
	Specializations   string          `json:"specializations"`   // JSON array of brands/models
	SupportedBrands   string          `json:"supported_brands"`  // JSON array
	CertifiedFor      string          `json:"certified_for"`     // JSON array of certifications
	PartsInventory    string          `json:"parts_inventory"`   // JSON object
	DiagnosticTools   string          `json:"diagnostic_tools"`  // JSON array
	RepairTechniques  string          `json:"repair_techniques"` // JSON array
	MaxDeviceValue    decimal.Decimal `sql:"type:decimal(10,2)" json:"max_device_value"`
	MinDeviceValue    decimal.Decimal `sql:"type:decimal(10,2)" json:"min_device_value"`
	SameDayRepair     bool            `gorm:"default:false" json:"same_day_repair"`
	DataRecovery      bool            `gorm:"default:false" json:"data_recovery"`
	BoardLevelRepair  bool            `gorm:"default:false" json:"board_level_repair"`
	WaterDamageRepair bool            `gorm:"default:false" json:"water_damage_repair"`
}

// ShopCertification contains certification and quality info
type ShopCertification struct {
	CertificationLevel  string     `gorm:"not null;default:'silver'" json:"certification_level"`
	CertificationDate   time.Time  `json:"certification_date"`
	CertificationExpiry time.Time  `json:"certification_expiry"`
	QualityScore        float64    `gorm:"default:0" json:"quality_score"`
	ComplianceScore     float64    `gorm:"default:0" json:"compliance_score"`
	SafetyScore         float64    `gorm:"default:0" json:"safety_score"`
	EnvironmentalScore  float64    `gorm:"default:0" json:"environmental_score"`
	ISOCertifications   string     `json:"iso_certifications"`   // JSON array
	IndustryMemberships string     `json:"industry_memberships"` // JSON array
	Awards              string     `json:"awards"`               // JSON array
	LastAuditDate       *time.Time `json:"last_audit_date,omitempty"`
	NextAuditDate       *time.Time `json:"next_audit_date,omitempty"`
	AuditResults        string     `json:"audit_results,omitempty"` // JSON object
}

// ShopPerformance contains performance metrics
type ShopPerformance struct {
	Rating               float64 `gorm:"default:0" json:"rating"`
	ReviewCount          int     `gorm:"default:0" json:"review_count"`
	CompletedRepairs     int     `gorm:"default:0" json:"completed_repairs"`
	SuccessRate          float64 `gorm:"default:0" json:"success_rate"`
	FirstTimeFixRate     float64 `gorm:"default:0" json:"first_time_fix_rate"`
	RepeatRepairRate     float64 `gorm:"default:0" json:"repeat_repair_rate"`
	AverageRepairTime    int     `json:"average_repair_time_hours"`
	OnTimePercentage     float64 `gorm:"default:0" json:"on_time_percentage"`
	CustomerSatisfaction float64 `gorm:"default:0" json:"customer_satisfaction"`
	NetPromoterScore     int     `json:"net_promoter_score"`
	ComplaintRate        float64 `gorm:"default:0" json:"complaint_rate"`
	ResponseTime         int     `json:"response_time_minutes"`
	ResolutionTime       int     `json:"resolution_time_hours"`
	EscalationRate       float64 `gorm:"default:0" json:"escalation_rate"`
}

// ShopFinancial contains financial information
type ShopFinancial struct {
	PricingTier        string          `gorm:"default:'standard'" json:"pricing_tier"`
	LaborRate          decimal.Decimal `sql:"type:decimal(10,2)" json:"labor_rate"`
	DiagnosticFee      decimal.Decimal `sql:"type:decimal(10,2)" json:"diagnostic_fee"`
	MinimumCharge      decimal.Decimal `sql:"type:decimal(10,2)" json:"minimum_charge"`
	PaymentTerms       string          `json:"payment_terms"`     // JSON object
	AcceptedPayments   string          `json:"accepted_payments"` // JSON array
	InsuranceAccepted  bool            `gorm:"default:true" json:"insurance_accepted"`
	DirectBilling      bool            `gorm:"default:false" json:"direct_billing"`
	WarrantyPeriod     int             `gorm:"default:90" json:"warranty_period_days"`
	RefundPolicy       string          `json:"refund_policy"`
	MonthlyRevenue     decimal.Decimal `sql:"type:decimal(10,2)" json:"monthly_revenue"`
	OutstandingBalance decimal.Decimal `sql:"type:decimal(10,2)" json:"outstanding_balance"`
	CreditLimit        decimal.Decimal `sql:"type:decimal(10,2)" json:"credit_limit"`
	PaymentStatus      string          `gorm:"default:'current'" json:"payment_status"`
	LastPaymentDate    *time.Time      `json:"last_payment_date,omitempty"`
	CommissionRate     float64         `gorm:"default:0.15" json:"commission_rate"`
}

// ShopCompliance contains compliance and legal info
type ShopCompliance struct {
	ComplianceStatus        string          `gorm:"default:'compliant'" json:"compliance_status"`
	LicenseStatus           string          `gorm:"default:'active'" json:"license_status"`
	LicenseExpiry           time.Time       `json:"license_expiry"`
	InsurancePolicy         string          `json:"insurance_policy"`
	InsuranceProvider       string          `json:"insurance_provider"`
	InsuranceExpiry         time.Time       `json:"insurance_expiry"`
	LiabilityCoverage       decimal.Decimal `sql:"type:decimal(10,2)" json:"liability_coverage"`
	BondNumber              string          `json:"bond_number,omitempty"`
	BondAmount              decimal.Decimal `sql:"type:decimal(10,2)" json:"bond_amount"`
	RegulatoryViolations    int             `gorm:"default:0" json:"regulatory_violations"`
	LastViolationDate       *time.Time      `json:"last_violation_date,omitempty"`
	BackgroundCheckDate     *time.Time      `json:"background_check_date,omitempty"`
	BackgroundCheckStatus   string          `json:"background_check_status,omitempty"`
	DataProtectionCompliant bool            `gorm:"default:false" json:"data_protection_compliant"`
}

// ShopIntegration contains third-party integration info
type ShopIntegration struct {
	ManagementSystem      string     `json:"management_system,omitempty"`
	InventorySystem       string     `json:"inventory_system,omitempty"`
	POSSystem             string     `json:"pos_system,omitempty"`
	APIEnabled            bool       `gorm:"default:false" json:"api_enabled"`
	APIKey                string     `json:"-"` // Hidden in JSON
	WebhookURL            string     `json:"webhook_url,omitempty"`
	IntegrationStatus     string     `gorm:"default:'pending'" json:"integration_status"`
	LastSyncDate          *time.Time `json:"last_sync_date,omitempty"`
	SyncFrequency         string     `json:"sync_frequency,omitempty"`
	AutoUpdateInventory   bool       `gorm:"default:false" json:"auto_update_inventory"`
	AutoConfirmBookings   bool       `gorm:"default:false" json:"auto_confirm_bookings"`
	AutoSendNotifications bool       `gorm:"default:true" json:"auto_send_notifications"`
}

// ShopStatus contains status and lifecycle info
type ShopStatus struct {
	IsActive              bool       `gorm:"default:true;index" json:"is_active"`
	Status                string     `gorm:"default:'active';index" json:"status"`
	StatusReason          string     `json:"status_reason,omitempty"`
	StatusChangedAt       *time.Time `json:"status_changed_at,omitempty"`
	StatusChangedBy       *uuid.UUID `gorm:"type:uuid" json:"status_changed_by,omitempty"`
	OnboardingCompleted   bool       `gorm:"default:false" json:"onboarding_completed"`
	OnboardingCompletedAt *time.Time `json:"onboarding_completed_at,omitempty"`
	VerificationStatus    string     `gorm:"default:'pending'" json:"verification_status"`
	VerificationDate      *time.Time `json:"verification_date,omitempty"`
	LastActivityDate      *time.Time `json:"last_activity_date,omitempty"`
	DeactivationDate      *time.Time `json:"deactivation_date,omitempty"`
	ReactivationDate      *time.Time `json:"reactivation_date,omitempty"`
}

// ShopMetrics contains business metrics
type ShopMetrics struct {
	TotalDevicesRepaired  int             `gorm:"default:0" json:"total_devices_repaired"`
	TotalRevenue          decimal.Decimal `sql:"type:decimal(12,2)" json:"total_revenue"`
	AverageTicketValue    decimal.Decimal `sql:"type:decimal(10,2)" json:"average_ticket_value"`
	RepeatCustomerRate    float64         `gorm:"default:0" json:"repeat_customer_rate"`
	ReferralRate          float64         `gorm:"default:0" json:"referral_rate"`
	ConversionRate        float64         `gorm:"default:0" json:"conversion_rate"`
	CancellationRate      float64         `gorm:"default:0" json:"cancellation_rate"`
	WarrantyClaimRate     float64         `gorm:"default:0" json:"warranty_claim_rate"`
	PartsWastageRate      float64         `gorm:"default:0" json:"parts_wastage_rate"`
	TechnicianUtilization float64         `gorm:"default:0" json:"technician_utilization"`
	AveragePartsCost      decimal.Decimal `sql:"type:decimal(10,2)" json:"average_parts_cost"`
	ProfitMargin          float64         `gorm:"default:0" json:"profit_margin"`
	YearOverYearGrowth    float64         `gorm:"default:0" json:"year_over_year_growth"`
	MarketShare           float64         `gorm:"default:0" json:"market_share"`
}
