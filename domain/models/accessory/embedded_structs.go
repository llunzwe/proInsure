package accessory

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AccessoryIdentification contains unique identifiers
type AccessoryIdentification struct {
	SKU            string `json:"sku"`
	Barcode        string `json:"barcode,omitempty"`
	UPC            string `json:"upc,omitempty"`
	EAN            string `json:"ean,omitempty"`
	ISBN           string `json:"isbn,omitempty"`
	PartNumber     string `json:"part_number,omitempty"`
	ManufacturerID string `json:"manufacturer_id,omitempty"`
	SupplierCode   string `json:"supplier_code,omitempty"`
}

// AccessoryBasicInfo contains basic product information
type AccessoryBasicInfo struct {
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model,omitempty"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	SubCategory string  `json:"sub_category,omitempty"`
	Description string  `json:"description"`
	ShortDesc   string  `json:"short_description,omitempty"`
	Features    string  `json:"features,omitempty"` // JSON array
	Color       string  `json:"color,omitempty"`
	Material    string  `json:"material,omitempty"`
	Weight      float64 `json:"weight,omitempty"` // grams
	WeightUnit  string  `json:"weight_unit,omitempty"`
}

// AccessoryCompatibility contains device compatibility info
type AccessoryCompatibility struct {
	CompatibilityType   string `json:"compatibility_type"`
	UniversalCompatible bool   `json:"universal_compatible"`
	BrandCompatible     string `json:"brand_compatible,omitempty"` // JSON array
	ModelCompatible     string `json:"model_compatible,omitempty"` // JSON array
	OSCompatible        string `json:"os_compatible,omitempty"`    // JSON array
	MinOSVersion        string `json:"min_os_version,omitempty"`
	MaxOSVersion        string `json:"max_os_version,omitempty"`
	ConnectionType      string `json:"connection_type,omitempty"`
	ConnectorType       string `json:"connector_type,omitempty"`
	WirelessProtocol    string `json:"wireless_protocol,omitempty"`
	TestedDevices       string `json:"tested_devices,omitempty"`    // JSON array
	IncompatibleWith    string `json:"incompatible_with,omitempty"` // JSON array
}

// AccessoryPricing contains pricing information
type AccessoryPricing struct {
	CostPrice       decimal.Decimal `json:"cost_price"`
	RetailPrice     decimal.Decimal `json:"retail_price"`
	SalePrice       decimal.Decimal `json:"sale_price,omitempty"`
	WholesalePrice  decimal.Decimal `json:"wholesale_price,omitempty"`
	MinPrice        decimal.Decimal `json:"min_price,omitempty"`
	MaxDiscount     float64         `json:"max_discount"`
	CurrentDiscount float64         `json:"current_discount"`
	TaxRate         float64         `json:"tax_rate"`
	Currency        string          `json:"currency"`
	PricingType     string          `json:"pricing_type"`
	TierPricing     string          `json:"tier_pricing,omitempty"` // JSON object
	BundlePrice     decimal.Decimal `json:"bundle_price,omitempty"`
	PromoPrice      decimal.Decimal `json:"promo_price,omitempty"`
	PromoStartDate  *time.Time      `json:"promo_start_date,omitempty"`
	PromoEndDate    *time.Time      `json:"promo_end_date,omitempty"`
}

// AccessoryInventory contains stock and inventory data
type AccessoryInventory struct {
	CurrentStock    int             `json:"current_stock"`
	AvailableStock  int             `json:"available_stock"`
	ReservedStock   int             `json:"reserved_stock"`
	IncomingStock   int             `json:"incoming_stock"`
	MinStock        int             `json:"min_stock"`
	MaxStock        int             `json:"max_stock"`
	ReorderPoint    int             `json:"reorder_point"`
	ReorderQuantity int             `json:"reorder_quantity"`
	StockStatus     string          `json:"stock_status"`
	Location        string          `json:"location,omitempty"`
	WarehouseID     uuid.UUID       `json:"warehouse_id,omitempty"`
	BinLocation     string          `json:"bin_location,omitempty"`
	LastStockCheck  *time.Time      `json:"last_stock_check,omitempty"`
	NextStockCheck  *time.Time      `json:"next_stock_check,omitempty"`
	StockValue      decimal.Decimal `json:"stock_value"`
}

// AccessorySupplier contains supplier information
type AccessorySupplier struct {
	SupplierID        uuid.UUID       `json:"supplier_id"`
	SupplierName      string          `json:"supplier_name"`
	SupplierType      string          `json:"supplier_type"`
	SupplierSKU       string          `json:"supplier_sku,omitempty"`
	LeadTimeDays      int             `json:"lead_time_days"`
	MinOrderQuantity  int             `json:"min_order_quantity"`
	SupplierCost      decimal.Decimal `json:"supplier_cost"`
	LastOrderDate     *time.Time      `json:"last_order_date,omitempty"`
	NextOrderDate     *time.Time      `json:"next_order_date,omitempty"`
	PreferredSupplier bool            `json:"preferred_supplier"`
	SupplierRating    float64         `json:"supplier_rating"`
	PaymentTerms      string          `json:"payment_terms,omitempty"`
	DeliveryTerms     string          `json:"delivery_terms,omitempty"`
	ContractExpiry    *time.Time      `json:"contract_expiry,omitempty"`
}

// AccessoryQuality contains quality and certification info
type AccessoryQuality struct {
	QualityGrade       string     `json:"quality_grade"`
	Condition          string     `json:"condition"`
	Certifications     string     `json:"certifications,omitempty"`   // JSON array
	SafetyStandards    string     `json:"safety_standards,omitempty"` // JSON array
	TestingCompleted   bool       `json:"testing_completed"`
	TestingDate        *time.Time `json:"testing_date,omitempty"`
	TestingReport      string     `json:"testing_report,omitempty"`
	DefectRate         float64    `json:"defect_rate"`
	ReturnRate         float64    `json:"return_rate"`
	QualityScore       float64    `json:"quality_score"`
	InspectionRequired bool       `json:"inspection_required"`
	LastInspection     *time.Time `json:"last_inspection,omitempty"`
	InspectionNotes    string     `json:"inspection_notes,omitempty"`
	RecallStatus       bool       `json:"recall_status"`
	RecallDate         *time.Time `json:"recall_date,omitempty"`
	RecallReason       string     `json:"recall_reason,omitempty"`
}

// AccessoryWarranty contains warranty information
type AccessoryWarranty struct {
	HasWarranty       bool            `json:"has_warranty"`
	WarrantyMonths    int             `json:"warranty_months"`
	WarrantyType      string          `json:"warranty_type,omitempty"`
	WarrantyProvider  string          `json:"warranty_provider,omitempty"`
	WarrantyStatus    string          `json:"warranty_status"`
	WarrantyStartDate *time.Time      `json:"warranty_start_date,omitempty"`
	WarrantyEndDate   *time.Time      `json:"warranty_end_date,omitempty"`
	ExtendedWarranty  bool            `json:"extended_warranty"`
	ExtendedMonths    int             `json:"extended_months"`
	ExtendedCost      decimal.Decimal `json:"extended_cost"`
	ReturnPolicy      string          `json:"return_policy,omitempty"`
	ReturnDays        int             `json:"return_days"`
	RestockingFee     float64         `json:"restocking_fee"`
	ExchangeAllowed   bool            `json:"exchange_allowed"`
}

// AccessoryDimensions contains physical dimensions
type AccessoryDimensions struct {
	Length          float64 `json:"length,omitempty"`    // cm
	Width           float64 `json:"width,omitempty"`     // cm
	Height          float64 `json:"height,omitempty"`    // cm
	Diameter        float64 `json:"diameter,omitempty"`  // cm
	Thickness       float64 `json:"thickness,omitempty"` // mm
	DimensionUnit   string  `json:"dimension_unit,omitempty"`
	PackageLength   float64 `json:"package_length,omitempty"`
	PackageWidth    float64 `json:"package_width,omitempty"`
	PackageHeight   float64 `json:"package_height,omitempty"`
	PackageWeight   float64 `json:"package_weight,omitempty"`
	PackageType     string  `json:"package_type,omitempty"`
	UnitsPerPackage int     `json:"units_per_package"`
}

// AccessoryMedia contains images and documentation
type AccessoryMedia struct {
	ThumbnailURL    string     `json:"thumbnail_url,omitempty"`
	ImageURLs       string     `json:"image_urls,omitempty"` // JSON array
	VideoURL        string     `json:"video_url,omitempty"`
	ManualURL       string     `json:"manual_url,omitempty"`
	DatasheetURL    string     `json:"datasheet_url,omitempty"`
	Gallery360URL   string     `json:"gallery_360_url,omitempty"`
	ARModelURL      string     `json:"ar_model_url,omitempty"` // Augmented reality
	InstallGuideURL string     `json:"install_guide_url,omitempty"`
	LastMediaUpdate *time.Time `json:"last_media_update,omitempty"`
}

// AccessoryMetrics contains performance metrics
type AccessoryMetrics struct {
	ViewCount       int             `json:"view_count"`
	SalesCount      int             `json:"sales_count"`
	Revenue         decimal.Decimal `json:"revenue"`
	Profit          decimal.Decimal `json:"profit"`
	ProfitMargin    float64         `json:"profit_margin"`
	AverageRating   float64         `json:"average_rating"`
	ReviewCount     int             `json:"review_count"`
	RecommendRate   float64         `json:"recommend_rate"`
	RepurchaseRate  float64         `json:"repurchase_rate"`
	ConversionRate  float64         `json:"conversion_rate"`
	CartAbandonRate float64         `json:"cart_abandon_rate"`
	DaysInInventory int             `json:"days_in_inventory"`
	TurnoverRate    float64         `json:"turnover_rate"`
	PopularityScore float64         `json:"popularity_score"`
	TrendingStatus  bool            `json:"trending_status"`
	SeasonalDemand  string          `json:"seasonal_demand,omitempty"` // JSON object
	ForecastDemand  int             `json:"forecast_demand"`
	LastSoldDate    *time.Time      `json:"last_sold_date,omitempty"`
	LastRestockDate *time.Time      `json:"last_restock_date,omitempty"`
	StockoutDays    int             `json:"stockout_days"`
}

// AccessoryStatus contains various status flags
type AccessoryStatus struct {
	IsActive         bool       `json:"is_active"`
	IsFeatured       bool       `json:"is_featured"`
	IsPromoted       bool       `json:"is_promoted"`
	IsDiscontinued   bool       `json:"is_discontinued"`
	IsBestseller     bool       `json:"is_bestseller"`
	IsNewArrival     bool       `json:"is_new_arrival"`
	IsClearance      bool       `json:"is_clearance"`
	IsExclusive      bool       `json:"is_exclusive"`
	IsLimitedEdition bool       `json:"is_limited_edition"`
	IsPreorder       bool       `json:"is_preorder"`
	IsBackorder      bool       `json:"is_backorder"`
	IsGiftable       bool       `json:"is_giftable"`
	IsCustomizable   bool       `json:"is_customizable"`
	RequiresShipping bool       `json:"requires_shipping"`
	OnlineOnly       bool       `json:"online_only"`
	InStoreOnly      bool       `json:"in_store_only"`
	LaunchDate       *time.Time `json:"launch_date,omitempty"`
	EndOfLife        *time.Time `json:"end_of_life,omitempty"`
}

// AccessoryAuthentication contains authenticity verification
type AccessoryAuthentication struct {
	IsGenuine            bool       `json:"is_genuine"`
	VerificationCode     string     `json:"verification_code,omitempty"`
	SerialNumber         string     `json:"serial_number,omitempty"`
	ManufactureDate      *time.Time `json:"manufacture_date,omitempty"`
	AuthenticationMethod string     `json:"authentication_method,omitempty"`
	QRCode               string     `json:"qr_code,omitempty"`
	HologramSerial       string     `json:"hologram_serial,omitempty"`
	IsCounterfeit        bool       `json:"is_counterfeit"`
	CounterfeitRiskScore float64    `json:"counterfeit_risk_score"`
	LastVerification     *time.Time `json:"last_verification,omitempty"`
	AuthorizedDealer     bool       `json:"authorized_dealer"`
	AuthorizedReseller   bool       `json:"authorized_reseller"`
}

// AccessoryTechnicalSpecs contains detailed technical specifications
type AccessoryTechnicalSpecs struct {
	// Power specifications
	InputVoltage     string  `json:"input_voltage,omitempty"`
	OutputVoltage    string  `json:"output_voltage,omitempty"`
	InputCurrent     float64 `json:"input_current,omitempty"`
	OutputCurrent    float64 `json:"output_current,omitempty"`
	PowerOutput      float64 `json:"power_output,omitempty"` // Watts
	ChargingSpeed    string  `json:"charging_speed,omitempty"`
	WirelessStandard string  `json:"wireless_standard,omitempty"` // Qi, MagSafe

	// Audio specifications
	FrequencyResponse string  `json:"frequency_response,omitempty"`
	Impedance         float64 `json:"impedance,omitempty"`
	Sensitivity       float64 `json:"sensitivity,omitempty"`
	DriverSize        float64 `json:"driver_size,omitempty"`
	NoiseReduction    string  `json:"noise_reduction,omitempty"`

	// Connectivity
	BluetoothVersion  string  `json:"bluetooth_version,omitempty"`
	WifiStandard      string  `json:"wifi_standard,omitempty"`
	CableLength       float64 `json:"cable_length,omitempty"`
	ConnectorTypes    string  `json:"connector_types,omitempty"` // JSON array
	DataTransferSpeed string  `json:"data_transfer_speed,omitempty"`

	// Protection ratings
	IPRating       string  `json:"ip_rating,omitempty"`
	MILSTDRating   string  `json:"milstd_rating,omitempty"`
	DropTestHeight float64 `json:"drop_test_height,omitempty"`
	OperatingTemp  string  `json:"operating_temp,omitempty"`
	StorageTemp    string  `json:"storage_temp,omitempty"`
}

// AccessoryMarketData contains market intelligence
type AccessoryMarketData struct {
	DemandScore        float64         `json:"demand_score"`
	TrendStatus        string          `json:"trend_status"`
	SeasonalDemand     string          `json:"seasonal_demand,omitempty"` // JSON object
	CompetitorCount    int             `json:"competitor_count"`
	CompetitorMinPrice decimal.Decimal `json:"competitor_min_price"`
	CompetitorMaxPrice decimal.Decimal `json:"competitor_max_price"`
	CompetitorAvgPrice decimal.Decimal `json:"competitor_avg_price"`
	MarketShare        float64         `json:"market_share"`
	CustomerSegments   string          `json:"customer_segments,omitempty"` // JSON object
	GeographicDemand   string          `json:"geographic_demand,omitempty"` // JSON object
	PriceElasticity    float64         `json:"price_elasticity"`
	OptimalPrice       decimal.Decimal `json:"optimal_price"`
}

// AccessoryEnvironmental contains environmental impact data
type AccessoryEnvironmental struct {
	CarbonFootprint      float64 `json:"carbon_footprint_kg"`
	RecyclablePercent    float64 `json:"recyclable_percent"`
	BiodegradablePercent float64 `json:"biodegradable_percent"`
	PackagingEcoFriendly bool    `json:"packaging_eco_friendly"`
	EWasteCategory       string  `json:"ewaste_category,omitempty"`
	DisposalMethod       string  `json:"disposal_method,omitempty"`
	ConflictMinerals     bool    `json:"conflict_minerals"`
	EnvironmentalScore   float64 `json:"environmental_score"`
	CarbonNeutral        bool    `json:"carbon_neutral"`
	SustainabilityCert   string  `json:"sustainability_cert,omitempty"`
}
