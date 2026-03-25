package sparepart

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SparePartIdentification contains unique identifiers
type SparePartIdentification struct {
	PartNumber     string `json:"part_number"`
	SKU            string `json:"sku"`
	ManufacturerPN string `json:"manufacturer_pn,omitempty"`
	OEMPartNumber  string `json:"oem_part_number,omitempty"`
	AlternativePN  string `json:"alternative_pn,omitempty"` // JSON array
	Barcode        string `json:"barcode,omitempty"`
	QRCode         string `json:"qr_code,omitempty"`
	SerialNumber   string `json:"serial_number,omitempty"`
	BatchNumber    string `json:"batch_number,omitempty"`
	RevisionNumber string `json:"revision_number,omitempty"`
}

// SparePartBasicInfo contains basic part information
type SparePartBasicInfo struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Category        string  `json:"category"`
	SubCategory     string  `json:"sub_category,omitempty"`
	Brand           string  `json:"brand"`
	Manufacturer    string  `json:"manufacturer,omitempty"`
	Model           string  `json:"model,omitempty"`
	Description     string  `json:"description"`
	TechnicalSpecs  string  `json:"technical_specs,omitempty"` // JSON object
	Color           string  `json:"color,omitempty"`
	Material        string  `json:"material,omitempty"`
	Weight          float64 `json:"weight,omitempty"` // grams
	IsCritical      bool    `json:"is_critical"`
	IsConsumable    bool    `json:"is_consumable"`
	RequiresTesting bool    `json:"requires_testing"`
}

// SparePartCompatibility contains device compatibility
type SparePartCompatibility struct {
	CompatibilityLevel  string `json:"compatibility_level"`
	DeviceBrands        string `json:"device_brands,omitempty"`   // JSON array
	DeviceModels        string `json:"device_models,omitempty"`   // JSON array
	DeviceVersions      string `json:"device_versions,omitempty"` // JSON array
	DeviceYears         string `json:"device_years,omitempty"`    // JSON array
	RegionSpecific      bool   `json:"region_specific"`
	RegionCodes         string `json:"region_codes,omitempty"`        // JSON array
	CrossReference      string `json:"cross_reference,omitempty"`     // JSON object
	NotCompatibleWith   string `json:"not_compatible_with,omitempty"` // JSON array
	RequiredTools       string `json:"required_tools,omitempty"`      // JSON array
	InstallationTime    int    `json:"installation_time"`             // minutes
	DifficultyLevel     string `json:"difficulty_level,omitempty"`
	RequiresProgramming bool   `json:"requires_programming"`
	RequiresCalibration bool   `json:"requires_calibration"`
}

// SparePartPricing contains cost and pricing information
type SparePartPricing struct {
	CostPrice         decimal.Decimal `json:"cost_price"`
	StandardPrice     decimal.Decimal `json:"standard_price"`
	RepairPrice       decimal.Decimal `json:"repair_price"`
	WholesalePrice    decimal.Decimal `json:"wholesale_price,omitempty"`
	RetailPrice       decimal.Decimal `json:"retail_price,omitempty"`
	EmergencyPrice    decimal.Decimal `json:"emergency_price,omitempty"`
	LastPurchasePrice decimal.Decimal `json:"last_purchase_price"`
	AverageCost       decimal.Decimal `json:"average_cost"`
	TotalCost         decimal.Decimal `json:"total_cost"`
	Currency          string          `json:"currency"`
	TaxRate           float64         `json:"tax_rate"`
	ImportDuty        float64         `json:"import_duty,omitempty"`
	ShippingCost      decimal.Decimal `json:"shipping_cost,omitempty"`
	HandlingCost      decimal.Decimal `json:"handling_cost,omitempty"`
	MarkupPercentage  float64         `json:"markup_percentage"`
	PriceLastUpdated  *time.Time      `json:"price_last_updated,omitempty"`
}

// SparePartInventory contains stock management data
type SparePartInventory struct {
	CurrentStock        int        `json:"current_stock"`
	AvailableStock      int        `json:"available_stock"`
	ReservedStock       int        `json:"reserved_stock"`
	AllocatedStock      int        `json:"allocated_stock"`
	QuarantineStock     int        `json:"quarantine_stock"`
	MinimumStock        int        `json:"minimum_stock"`
	MaximumStock        int        `json:"maximum_stock"`
	ReorderPoint        int        `json:"reorder_point"`
	ReorderQuantity     int        `json:"reorder_quantity"`
	OptimalStock        int        `json:"optimal_stock"`
	StockStatus         string     `json:"stock_status"`
	Location            string     `json:"location"`
	ShelfNumber         string     `json:"shelf_number,omitempty"`
	BinLocation         string     `json:"bin_location,omitempty"`
	WarehouseID         uuid.UUID  `json:"warehouse_id,omitempty"`
	RepairShopID        uuid.UUID  `json:"repair_shop_id,omitempty"`
	LastStockCheck      *time.Time `json:"last_stock_check,omitempty"`
	NextStockCheck      *time.Time `json:"next_stock_check,omitempty"`
	CycleCountFrequency string     `json:"cycle_count_frequency,omitempty"`
}

// SparePartSupplier contains supplier information
type SparePartSupplier struct {
	PrimarySupplierID   uuid.UUID       `json:"primary_supplier_id"`
	PrimarySupplierName string          `json:"primary_supplier_name"`
	SupplierType        string          `json:"supplier_type"`
	SupplierPartNumber  string          `json:"supplier_part_number,omitempty"`
	LeadTimeDays        int             `json:"lead_time_days"`
	MinOrderQuantity    int             `json:"min_order_quantity"`
	OrderMultiple       int             `json:"order_multiple"`
	SupplierPrice       decimal.Decimal `json:"supplier_price"`
	LastOrderDate       *time.Time      `json:"last_order_date,omitempty"`
	LastDeliveryDate    *time.Time      `json:"last_delivery_date,omitempty"`
	NextOrderDate       *time.Time      `json:"next_order_date,omitempty"`
	SupplierRating      float64         `json:"supplier_rating"`
	ReliabilityScore    float64         `json:"reliability_score"`
	QualityScore        float64         `json:"quality_score"`
	AlternateSuppliers  string          `json:"alternate_suppliers,omitempty"` // JSON array
	PreferredSupplier   bool            `json:"preferred_supplier"`
	ContractNumber      string          `json:"contract_number,omitempty"`
	ContractExpiry      *time.Time      `json:"contract_expiry,omitempty"`
}

// SparePartQuality contains quality and testing information
type SparePartQuality struct {
	QualityGrade        string     `json:"quality_grade"`
	Condition           string     `json:"condition"`
	TestingStatus       string     `json:"testing_status"`
	TestingDate         *time.Time `json:"testing_date,omitempty"`
	TestingReport       string     `json:"testing_report,omitempty"`
	TestingCriteria     string     `json:"testing_criteria,omitempty"` // JSON object
	QualityScore        float64    `json:"quality_score"`
	DefectRate          float64    `json:"defect_rate"`
	FailureRate         float64    `json:"failure_rate"`
	ReturnRate          float64    `json:"return_rate"`
	Certifications      string     `json:"certifications,omitempty"`       // JSON array
	ComplianceStandards string     `json:"compliance_standards,omitempty"` // JSON array
	InspectionRequired  bool       `json:"inspection_required"`
	LastInspection      *time.Time `json:"last_inspection,omitempty"`
	NextInspection      *time.Time `json:"next_inspection,omitempty"`
	InspectionNotes     string     `json:"inspection_notes,omitempty"`
	ApprovedForUse      bool       `json:"approved_for_use"`
	ApprovalDate        *time.Time `json:"approval_date,omitempty"`
	ApprovedBy          *uuid.UUID `json:"approved_by,omitempty"`
}

// SparePartWarranty contains warranty information
type SparePartWarranty struct {
	HasWarranty          bool    `json:"has_warranty"`
	WarrantyType         string  `json:"warranty_type"`
	WarrantyDays         int     `json:"warranty_days"`
	WarrantyProvider     string  `json:"warranty_provider,omitempty"`
	WarrantyTerms        string  `json:"warranty_terms,omitempty"`
	ClaimProcess         string  `json:"claim_process,omitempty"`
	CoverageDetails      string  `json:"coverage_details,omitempty"`
	ExclusionDetails     string  `json:"exclusion_details,omitempty"`
	RequiresRegistration bool    `json:"requires_registration"`
	ReturnAcceptable     bool    `json:"return_acceptable"`
	ReturnWindow         int     `json:"return_window"` // days
	RestockingFee        float64 `json:"restocking_fee"`
	ExchangeAllowed      bool    `json:"exchange_allowed"`
	RefundAllowed        bool    `json:"refund_allowed"`
}

// SparePartLifecycle contains lifecycle management data
type SparePartLifecycle struct {
	ManufactureDate   *time.Time `json:"manufacture_date,omitempty"`
	ReceiptDate       *time.Time `json:"receipt_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
	ShelfLife         int        `json:"shelf_life"` // months
	UsageLife         int        `json:"usage_life"` // months
	InstallCount      int        `json:"install_count"`
	FailureCount      int        `json:"failure_count"`
	ReturnCount       int        `json:"return_count"`
	ScrapCount        int        `json:"scrap_count"`
	TotalIssued       int        `json:"total_issued"`
	TotalConsumed     int        `json:"total_consumed"`
	IsObsolete        bool       `json:"is_obsolete"`
	ObsoleteDate      *time.Time `json:"obsolete_date,omitempty"`
	ReplacementPartID *uuid.UUID `json:"replacement_part_id,omitempty"`
	PhaseOutDate      *time.Time `json:"phase_out_date,omitempty"`
	EndOfLife         *time.Time `json:"end_of_life,omitempty"`
	DisposalMethod    string     `json:"disposal_method,omitempty"`
}

// SparePartUsage contains usage tracking information
type SparePartUsage struct {
	AverageUsagePerMonth float64    `json:"average_usage_per_month"`
	PeakUsagePerMonth    float64    `json:"peak_usage_per_month"`
	CurrentMonthUsage    int        `json:"current_month_usage"`
	LastMonthUsage       int        `json:"last_month_usage"`
	YearToDateUsage      int        `json:"year_to_date_usage"`
	UsageTrend           string     `json:"usage_trend,omitempty"`      // increasing, stable, decreasing
	SeasonalPattern      string     `json:"seasonal_pattern,omitempty"` // JSON object
	ForecastNextMonth    int        `json:"forecast_next_month"`
	RepairTypeUsage      string     `json:"repair_type_usage,omitempty"`   // JSON object
	TechnicianUsage      string     `json:"technician_usage,omitempty"`    // JSON object
	CustomerTypeUsage    string     `json:"customer_type_usage,omitempty"` // JSON object
	HighUsagePeriods     string     `json:"high_usage_periods,omitempty"`  // JSON array
	LastUsedDate         *time.Time `json:"last_used_date,omitempty"`
	LastIssuedTo         *uuid.UUID `json:"last_issued_to,omitempty"`
	MostUsedFor          string     `json:"most_used_for,omitempty"`
}

// SparePartCompliance contains regulatory compliance data
type SparePartCompliance struct {
	IsRegulated         bool       `json:"is_regulated"`
	RegulatoryCategory  string     `json:"regulatory_category,omitempty"`
	ComplianceStatus    string     `json:"compliance_status"`
	HSCode              string     `json:"hs_code,omitempty"`
	CountryOfOrigin     string     `json:"country_of_origin,omitempty"`
	ImportRestrictions  string     `json:"import_restrictions,omitempty"` // JSON array
	ExportRestrictions  string     `json:"export_restrictions,omitempty"` // JSON array
	HazardousClass      string     `json:"hazardous_class,omitempty"`
	SafetyDataSheet     string     `json:"safety_data_sheet,omitempty"`
	HandlingPrecautions string     `json:"handling_precautions,omitempty"`
	StorageRequirements string     `json:"storage_requirements,omitempty"`
	DisposalGuidelines  string     `json:"disposal_guidelines,omitempty"`
	EnvironmentalImpact string     `json:"environmental_impact,omitempty"`
	RecyclingCode       string     `json:"recycling_code,omitempty"`
	REACHCompliant      bool       `json:"reach_compliant"`
	RoHSCompliant       bool       `json:"rohs_compliant"`
	LastComplianceCheck *time.Time `json:"last_compliance_check,omitempty"`
	NextComplianceCheck *time.Time `json:"next_compliance_check,omitempty"`
}

// SparePartAuthentication contains authenticity verification data
type SparePartAuthentication struct {
	IsGenuine            bool       `json:"is_genuine"`
	VerificationCode     string     `json:"verification_code,omitempty"`
	SerialNumber         string     `json:"serial_number,omitempty"`
	IMEI                 string     `json:"imei,omitempty"`
	ManufactureDate      *time.Time `json:"manufacture_date,omitempty"`
	AuthenticationMethod string     `json:"authentication_method,omitempty"`
	BlockchainHash       string     `json:"blockchain_hash,omitempty"`
	QRCode               string     `json:"qr_code,omitempty"`
	HologramSerial       string     `json:"hologram_serial,omitempty"`
	IsCounterfeit        bool       `json:"is_counterfeit"`
	CounterfeitRiskScore float64    `json:"counterfeit_risk_score"`
	LastVerification     *time.Time `json:"last_verification,omitempty"`
	VerificationHistory  string     `json:"verification_history,omitempty"` // JSON array
	AuthorizedDealer     bool       `json:"authorized_dealer"`
}

// SparePartTechnicalRequirements contains installation and repair requirements
type SparePartTechnicalRequirements struct {
	RepairDifficulty     int    `json:"repair_difficulty"` // 1-5 scale
	EstimatedTime        int    `json:"estimated_time_minutes"`
	RequiredTools        string `json:"required_tools,omitempty"` // JSON array
	RequiredSkillLevel   string `json:"required_skill_level"`
	VideoGuideURL        string `json:"video_guide_url,omitempty"`
	ManualURL            string `json:"manual_url,omitempty"`
	DiagnosticStepsURL   string `json:"diagnostic_steps_url,omitempty"`
	CalibrationRequired  bool   `json:"calibration_required"`
	CalibrationProcedure string `json:"calibration_procedure,omitempty"`
	ProgrammingRequired  bool   `json:"programming_required"`
	ProgrammingSteps     string `json:"programming_steps,omitempty"`
	SpecialInstructions  string `json:"special_instructions,omitempty"`
	SafetyPrecautions    string `json:"safety_precautions,omitempty"`
	CustomerInstallable  bool   `json:"customer_installable"`
	ProfessionalRequired bool   `json:"professional_required"`
	TrainingRequired     bool   `json:"training_required"`
	TrainingCourseID     string `json:"training_course_id,omitempty"`
}

// SparePartEnvironmentalImpact contains environmental and sustainability data
type SparePartEnvironmentalImpact struct {
	CarbonFootprint       float64 `json:"carbon_footprint_kg"`
	RecyclablePercent     float64 `json:"recyclable_percent"`
	BiodegradablePercent  float64 `json:"biodegradable_percent"`
	HazardousWaste        bool    `json:"hazardous_waste"`
	EWasteCategory        string  `json:"ewaste_category,omitempty"`
	DisposalMethod        string  `json:"disposal_method,omitempty"`
	RecyclingInstructions string  `json:"recycling_instructions,omitempty"`
	ConflictMinerals      bool    `json:"conflict_minerals"`
	ConflictMineralsList  string  `json:"conflict_minerals_list,omitempty"` // JSON array
	EnvironmentalScore    float64 `json:"environmental_score"`
	SustainabilityCert    string  `json:"sustainability_cert,omitempty"`
	WEEECompliant         bool    `json:"weee_compliant"`
	EnergyStarRated       bool    `json:"energy_star_rated"`
	CarbonNeutral         bool    `json:"carbon_neutral"`
	CarbonOffsetProgram   string  `json:"carbon_offset_program,omitempty"`
}

// SparePartB2BPricing contains business-to-business pricing tiers
type SparePartB2BPricing struct {
	TierName       string          `json:"tier_name"`
	MinQuantity    int             `json:"min_quantity"`
	MaxQuantity    int             `json:"max_quantity"`
	UnitPrice      decimal.Decimal `sql:"type:decimal(10,2)" json:"unit_price"`
	VolumeDiscount float64         `json:"volume_discount"`
	PaymentTerms   string          `json:"payment_terms"`
	CreditLimit    decimal.Decimal `sql:"type:decimal(10,2)" json:"credit_limit"`
	ContractPrice  bool            `json:"contract_price"`
	ContractID     string          `json:"contract_id,omitempty"`
	SpecialTerms   string          `json:"special_terms,omitempty"`
	ValidFrom      *time.Time      `json:"valid_from,omitempty"`
	ValidUntil     *time.Time      `json:"valid_until,omitempty"`
}

// SparePartMetrics contains performance metrics
type SparePartMetrics struct {
	TurnoverRate        float64         `json:"turnover_rate"`
	StockoutFrequency   int             `json:"stockout_frequency"`
	StockoutDays        int             `json:"stockout_days"`
	ServiceLevel        float64         `json:"service_level"`
	FillRate            float64         `json:"fill_rate"`
	SuccessRate         float64         `json:"success_rate"`
	CostPerUse          decimal.Decimal `json:"cost_per_use"`
	RevenueGenerated    decimal.Decimal `json:"revenue_generated"`
	ProfitMargin        float64         `json:"profit_margin"`
	CarryingCost        decimal.Decimal `json:"carrying_cost"`
	OrderingCost        decimal.Decimal `json:"ordering_cost"`
	TotalCostOwnership  decimal.Decimal `json:"total_cost_ownership"`
	DaysOfSupply        int             `json:"days_of_supply"`
	LeadTimeVariability float64         `json:"lead_time_variability"`
	DemandVariability   float64         `json:"demand_variability"`
	CriticalityScore    float64         `json:"criticality_score"`
	PerformanceScore    float64         `json:"performance_score"`
}
