package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// DeviceContract represents contract and agreement management
type DeviceContract struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Lease Agreements
	LeaseAgreement     datatypes.JSON `gorm:"type:json" json:"lease_agreement"` // LeaseContract
	LeaseStartDate     *time.Time     `gorm:"type:timestamp" json:"lease_start_date,omitempty"`
	LeaseEndDate       *time.Time     `gorm:"type:timestamp" json:"lease_end_date,omitempty"`
	MonthlyLeaseAmount float64        `gorm:"type:decimal(15,2)" json:"monthly_lease_amount"`
	LeaseStatus        string         `gorm:"type:varchar(50)" json:"lease_status"`
	BuyoutOption       bool           `gorm:"type:boolean;default:false" json:"buyout_option"`
	BuyoutAmount       float64        `gorm:"type:decimal(15,2)" json:"buyout_amount"`

	// Service Contracts
	ServiceContracts     datatypes.JSON `gorm:"type:json" json:"service_contracts"` // []ServiceContract
	ActiveServicePlans   int            `gorm:"type:int" json:"active_service_plans"`
	ServiceLevel         string         `gorm:"type:varchar(50)" json:"service_level"` // basic, standard, premium
	ServiceProvider      string         `gorm:"type:varchar(255)" json:"service_provider"`
	ServiceContractValue float64        `gorm:"type:decimal(15,2)" json:"service_contract_value"`

	// Extended Warranty
	ExtendedWarranty   datatypes.JSON `gorm:"type:json" json:"extended_warranty"` // WarrantyContract
	WarrantyProvider   string         `gorm:"type:varchar(255)" json:"warranty_provider"`
	WarrantyStartDate  *time.Time     `gorm:"type:timestamp" json:"warranty_start_date,omitempty"`
	WarrantyEndDate    *time.Time     `gorm:"type:timestamp" json:"warranty_end_date,omitempty"`
	WarrantyCost       float64        `gorm:"type:decimal(15,2)" json:"warranty_cost"`
	WarrantyClaimLimit float64        `gorm:"type:decimal(15,2)" json:"warranty_claim_limit"`

	// Insurance Policies
	InsurancePolicies   datatypes.JSON `gorm:"type:json" json:"insurance_policies"` // []InsurancePolicy
	PrimaryInsurer      string         `gorm:"type:varchar(255)" json:"primary_insurer"`
	TotalCoverageAmount float64        `gorm:"type:decimal(15,2)" json:"total_coverage_amount"`
	ActivePoliciesCount int            `gorm:"type:int" json:"active_policies_count"`

	// Terms of Service
	TermsOfService  datatypes.JSON `gorm:"type:json" json:"terms_of_service"` // []TOS
	TOSVersion      string         `gorm:"type:varchar(50)" json:"tos_version"`
	TOSAcceptedDate *time.Time     `gorm:"type:timestamp" json:"tos_accepted_date,omitempty"`
	PrivacyPolicy   datatypes.JSON `gorm:"type:json" json:"privacy_policy"`
	DataAgreements  datatypes.JSON `gorm:"type:json" json:"data_agreements"` // []Agreement

	// Maintenance Agreements
	MaintenanceContracts  datatypes.JSON `gorm:"type:json" json:"maintenance_contracts"` // []MaintenanceContract
	PreventiveMaintenance bool           `gorm:"type:boolean;default:false" json:"preventive_maintenance"`
	MaintenanceSchedule   datatypes.JSON `gorm:"type:json" json:"maintenance_schedule"`
	NextMaintenanceDate   *time.Time     `gorm:"type:timestamp" json:"next_maintenance_date,omitempty"`

	// Financial Agreements
	FinancingAgreement datatypes.JSON `gorm:"type:json" json:"financing_agreement"`
	PaymentPlan        datatypes.JSON `gorm:"type:json" json:"payment_plan"`
	InstallmentAmount  float64        `gorm:"type:decimal(15,2)" json:"installment_amount"`
	RemainingBalance   float64        `gorm:"type:decimal(15,2)" json:"remaining_balance"`

	// Contract Management
	ContractRenewals   datatypes.JSON `gorm:"type:json" json:"contract_renewals"` // []Renewal
	UpcomingRenewals   int            `gorm:"type:int" json:"upcoming_renewals"`
	AutoRenewal        bool           `gorm:"type:boolean;default:false" json:"auto_renewal"`
	ContractViolations datatypes.JSON `gorm:"type:json" json:"contract_violations"` // []Violation

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceShipping represents shipping and logistics management
type DeviceShipping struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Shipping History
	ShippingHistory  datatypes.JSON `gorm:"type:json" json:"shipping_history"` // []ShipmentRecord
	TotalShipments   int            `gorm:"type:int" json:"total_shipments"`
	LastShipmentDate *time.Time     `gorm:"type:timestamp" json:"last_shipment_date,omitempty"`

	// Active Shipments
	ActiveShipments   datatypes.JSON `gorm:"type:json" json:"active_shipments"` // []Shipment
	TrackingNumbers   datatypes.JSON `gorm:"type:json" json:"tracking_numbers"` // []string
	CurrentLocation   string         `gorm:"type:varchar(255)" json:"current_location"`
	EstimatedDelivery *time.Time     `gorm:"type:timestamp" json:"estimated_delivery,omitempty"`

	// Courier Services
	PreferredCourier   string         `gorm:"type:varchar(255)" json:"preferred_courier"`
	CourierAccounts    datatypes.JSON `gorm:"type:json" json:"courier_accounts"`    // map[string]Account
	ShippingRates      datatypes.JSON `gorm:"type:json" json:"shipping_rates"`      // map[string]Rate
	CourierPerformance datatypes.JSON `gorm:"type:json" json:"courier_performance"` // map[string]Metrics

	// Delivery Management
	DeliveryConfirmations datatypes.JSON `gorm:"type:json" json:"delivery_confirmations"` // []Confirmation
	SignatureRequired     bool           `gorm:"type:boolean;default:false" json:"signature_required"`
	DeliveryInstructions  string         `gorm:"type:text" json:"delivery_instructions"`
	DeliveryAttempts      int            `gorm:"type:int" json:"delivery_attempts"`

	// Return Merchandise Authorization
	RMAHistory   datatypes.JSON `gorm:"type:json" json:"rma_history"` // []RMA
	ActiveRMA    datatypes.JSON `gorm:"type:json" json:"active_rma"`
	RMANumber    string         `gorm:"type:varchar(100)" json:"rma_number"`
	ReturnLabel  string         `gorm:"type:varchar(500)" json:"return_label"`
	ReturnStatus string         `gorm:"type:varchar(50)" json:"return_status"`
	RefundStatus string         `gorm:"type:varchar(50)" json:"refund_status"`

	// Shipping Costs
	TotalShippingCosts  float64 `gorm:"type:decimal(15,2)" json:"total_shipping_costs"`
	AverageShippingCost float64 `gorm:"type:decimal(15,2)" json:"average_shipping_cost"`
	ShippingInsurance   float64 `gorm:"type:decimal(15,2)" json:"shipping_insurance"`
	CustomsFees         float64 `gorm:"type:decimal(15,2)" json:"customs_fees"`

	// International Shipping
	InternationalShipping bool           `gorm:"type:boolean;default:false" json:"international_shipping"`
	ExportDocuments       datatypes.JSON `gorm:"type:json" json:"export_documents"`
	ImportDocuments       datatypes.JSON `gorm:"type:json" json:"import_documents"`
	CustomsDeclaration    datatypes.JSON `gorm:"type:json" json:"customs_declaration"`

	// Logistics Partners
	LogisticsPartners datatypes.JSON `gorm:"type:json" json:"logistics_partners"` // []Partner
	FulfillmentCenter string         `gorm:"type:varchar(255)" json:"fulfillment_center"`
	WarehouseLocation string         `gorm:"type:varchar(255)" json:"warehouse_location"`

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceTaxAccounting represents tax and accounting features
type DeviceTaxAccounting struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Tax Depreciation
	DepreciationMethod      string         `gorm:"type:varchar(50)" json:"depreciation_method"` // straight-line, declining, units
	DepreciationSchedule    datatypes.JSON `gorm:"type:json" json:"depreciation_schedule"`      // []DepreciationEntry
	AnnualDepreciation      float64        `gorm:"type:decimal(15,2)" json:"annual_depreciation"`
	AccumulatedDepreciation float64        `gorm:"type:decimal(15,2)" json:"accumulated_depreciation"`
	BookValue               float64        `gorm:"type:decimal(15,2)" json:"book_value"`
	SalvageValue            float64        `gorm:"type:decimal(15,2)" json:"salvage_value"`
	UsefulLifeYears         int            `gorm:"type:int" json:"useful_life_years"`

	// Tax Information
	PurchaseTax     float64 `gorm:"type:decimal(15,2)" json:"purchase_tax"`
	SalesTax        float64 `gorm:"type:decimal(15,2)" json:"sales_tax"`
	GSTAmount       float64 `gorm:"type:decimal(15,2)" json:"gst_amount"`
	VATAmount       float64 `gorm:"type:decimal(15,2)" json:"vat_amount"`
	TaxJurisdiction string  `gorm:"type:varchar(100)" json:"tax_jurisdiction"`
	TaxID           string  `gorm:"type:varchar(100)" json:"tax_id"`

	// Insurance Tax
	InsuranceTax     float64 `gorm:"type:decimal(15,2)" json:"insurance_tax"`
	PremiumTax       float64 `gorm:"type:decimal(15,2)" json:"premium_tax"`
	TaxDeductible    bool    `gorm:"type:boolean;default:false" json:"tax_deductible"`
	DeductibleAmount float64 `gorm:"type:decimal(15,2)" json:"deductible_amount"`

	// Capital Gains/Loss
	PurchasePrice       float64 `gorm:"type:decimal(15,2)" json:"purchase_price"`
	SalePrice           float64 `gorm:"type:decimal(15,2)" json:"sale_price"`
	CapitalGain         float64 `gorm:"type:decimal(15,2)" json:"capital_gain"`
	CapitalLoss         float64 `gorm:"type:decimal(15,2)" json:"capital_loss"`
	HoldingPeriodDays   int     `gorm:"type:int" json:"holding_period_days"`
	LongTermCapitalGain bool    `gorm:"type:boolean;default:false" json:"long_term_capital_gain"`

	// Asset Categorization
	AssetCategory string `gorm:"type:varchar(100)" json:"asset_category"` // fixed, current, intangible
	AssetClass    string `gorm:"type:varchar(100)" json:"asset_class"`
	AssetCode     string `gorm:"type:varchar(100)" json:"asset_code"`
	CostCenter    string `gorm:"type:varchar(100)" json:"cost_center"`
	ProfitCenter  string `gorm:"type:varchar(100)" json:"profit_center"`
	GLAccount     string `gorm:"type:varchar(100)" json:"gl_account"`

	// Accounting Entries
	JournalEntries   datatypes.JSON `gorm:"type:json" json:"journal_entries"` // []JournalEntry
	LedgerAccounts   datatypes.JSON `gorm:"type:json" json:"ledger_accounts"` // []LedgerAccount
	FiscalYear       string         `gorm:"type:varchar(20)" json:"fiscal_year"`
	AccountingPeriod string         `gorm:"type:varchar(20)" json:"accounting_period"`

	// Write-offs
	WriteOffStatus bool       `gorm:"type:boolean;default:false" json:"write_off_status"`
	WriteOffAmount float64    `gorm:"type:decimal(15,2)" json:"write_off_amount"`
	WriteOffDate   *time.Time `gorm:"type:timestamp" json:"write_off_date,omitempty"`
	WriteOffReason string     `gorm:"type:text" json:"write_off_reason"`

	// Compliance
	TaxComplianceStatus string     `gorm:"type:varchar(50)" json:"tax_compliance_status"`
	AuditStatus         string     `gorm:"type:varchar(50)" json:"audit_status"`
	LastAuditDate       *time.Time `gorm:"type:timestamp" json:"last_audit_date,omitempty"`
	TaxFilingStatus     string     `gorm:"type:varchar(50)" json:"tax_filing_status"`

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// HasActiveContract checks if device has active contracts
func (dc *DeviceContract) HasActiveContract() bool {
	return dc.ActiveServicePlans > 0 || dc.LeaseStatus == "active" ||
		dc.ActivePoliciesCount > 0
}

// IsLeased checks if device is currently leased
func (dc *DeviceContract) IsLeased() bool {
	if dc.LeaseStatus != "active" || dc.LeaseEndDate == nil {
		return false
	}
	return time.Now().Before(*dc.LeaseEndDate)
}

// NeedsRenewal checks if any contracts need renewal
func (dc *DeviceContract) NeedsRenewal() bool {
	return dc.UpcomingRenewals > 0 && !dc.AutoRenewal
}

// HasActiveShipment checks if device is currently being shipped
func (ds *DeviceShipping) HasActiveShipment() bool {
	return len(ds.TrackingNumbers) > 0 && ds.ActiveShipments != nil
}

// HasRMA checks if device has active return authorization
func (ds *DeviceShipping) HasRMA() bool {
	return ds.RMANumber != "" && ds.ReturnStatus == "active"
}

// GetCurrentDepreciation calculates current depreciation
func (dta *DeviceTaxAccounting) GetCurrentDepreciation() float64 {
	if dta.UsefulLifeYears == 0 {
		return 0
	}
	return dta.AnnualDepreciation
}

// IsFullyDepreciated checks if asset is fully depreciated
func (dta *DeviceTaxAccounting) IsFullyDepreciated() bool {
	return dta.BookValue <= dta.SalvageValue
}

// HasTaxLiability checks if there's tax liability
func (dta *DeviceTaxAccounting) HasTaxLiability() bool {
	totalTax := dta.PurchaseTax + dta.SalesTax + dta.GSTAmount +
		dta.VATAmount + dta.InsuranceTax + dta.PremiumTax
	return totalTax > 0
}
