package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// DeviceSupplyChain represents supply chain and provenance tracking
type DeviceSupplyChain struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Manufacturing Information
	ManufacturerID     string    `gorm:"type:varchar(100)" json:"manufacturer_id"`
	ManufacturerName   string    `gorm:"type:varchar(255)" json:"manufacturer_name"`
	ManufacturingDate  time.Time `gorm:"type:timestamp" json:"manufacturing_date"`
	ManufacturingPlant string    `gorm:"type:varchar(255)" json:"manufacturing_plant"`
	BatchNumber        string    `gorm:"type:varchar(100);index" json:"batch_number"`
	LotNumber          string    `gorm:"type:varchar(100)" json:"lot_number"`
	ProductionLineID   string    `gorm:"type:varchar(100)" json:"production_line_id"`
	QualityControlID   string    `gorm:"type:varchar(100)" json:"quality_control_id"`

	// Component Tracking
	ComponentSources   datatypes.JSON `gorm:"type:json" json:"component_sources"`   // []ComponentSource
	CriticalComponents datatypes.JSON `gorm:"type:json" json:"critical_components"` // []Component
	ComponentBatch     datatypes.JSON `gorm:"type:json" json:"component_batch"`     // map[string]string
	AssemblyDate       time.Time      `gorm:"type:timestamp" json:"assembly_date"`

	// Authenticity Verification
	AuthenticityStatus    string         `gorm:"type:varchar(50)" json:"authenticity_status"` // genuine, counterfeit, unknown
	AuthenticityScore     float64        `gorm:"type:decimal(5,2)" json:"authenticity_score"`
	VerificationMethod    string         `gorm:"type:varchar(100)" json:"verification_method"`
	VerificationDate      time.Time      `gorm:"type:timestamp" json:"verification_date"`
	VerifiedBy            string         `gorm:"type:varchar(255)" json:"verified_by"`
	CounterfeitIndicators datatypes.JSON `gorm:"type:json" json:"counterfeit_indicators"` // []string

	// Gray Market Detection
	MarketType           string `gorm:"type:varchar(50)" json:"market_type"` // official, gray, parallel, black
	ImportChannel        string `gorm:"type:varchar(100)" json:"import_channel"`
	ParallelImportStatus bool   `gorm:"type:boolean;default:false" json:"parallel_import_status"`
	AuthorizedDealer     bool   `gorm:"type:boolean;default:true" json:"authorized_dealer"`
	DealerID             string `gorm:"type:varchar(100)" json:"dealer_id"`
	DealerName           string `gorm:"type:varchar(255)" json:"dealer_name"`

	// Supply Chain Path
	SupplyChainSteps   datatypes.JSON `gorm:"type:json" json:"supply_chain_steps"` // []SupplyChainStep
	OriginCountry      string         `gorm:"type:varchar(100)" json:"origin_country"`
	DestinationCountry string         `gorm:"type:varchar(100)" json:"destination_country"`
	TransitCountries   datatypes.JSON `gorm:"type:json" json:"transit_countries"` // []string
	ImportDate         *time.Time     `gorm:"type:timestamp" json:"import_date,omitempty"`
	CustomsClearance   bool           `gorm:"type:boolean;default:false" json:"customs_clearance"`
	CustomsDocumentID  string         `gorm:"type:varchar(100)" json:"customs_document_id"`

	// Vulnerability Assessment
	VulnerabilityScore float64        `gorm:"type:decimal(5,2)" json:"vulnerability_score"`
	RiskFactors        datatypes.JSON `gorm:"type:json" json:"risk_factors"`     // []RiskFactor
	SecurityThreats    datatypes.JSON `gorm:"type:json" json:"security_threats"` // []Threat
	LastAssessmentDate time.Time      `gorm:"type:timestamp" json:"last_assessment_date"`

	// Certifications
	Certifications      datatypes.JSON `gorm:"type:json" json:"certifications"` // []Certification
	ComplianceStatus    string         `gorm:"type:varchar(50)" json:"compliance_status"`
	RegulatoryApprovals datatypes.JSON `gorm:"type:json" json:"regulatory_approvals"` // []Approval

	// Audit Trail
	AuditHistory  datatypes.JSON `gorm:"type:json" json:"audit_history"` // []AuditEntry
	LastAuditDate time.Time      `gorm:"type:timestamp" json:"last_audit_date"`
	NextAuditDate time.Time      `gorm:"type:timestamp" json:"next_audit_date"`

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceMultiCurrency represents multi-currency and international features
type DeviceMultiCurrency struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Currency Information
	PurchaseCurrency string    `gorm:"type:varchar(3)" json:"purchase_currency"`
	PurchaseAmount   float64   `gorm:"type:decimal(15,2)" json:"purchase_amount"`
	LocalCurrency    string    `gorm:"type:varchar(3)" json:"local_currency"`
	LocalAmount      float64   `gorm:"type:decimal(15,2)" json:"local_amount"`
	ExchangeRate     float64   `gorm:"type:decimal(15,6)" json:"exchange_rate"`
	ExchangeRateDate time.Time `gorm:"type:timestamp" json:"exchange_rate_date"`

	// International Warranty
	InternationalWarranty bool           `gorm:"type:boolean;default:false" json:"international_warranty"`
	WarrantyRegions       datatypes.JSON `gorm:"type:json" json:"warranty_regions"` // []string
	WarrantyRestrictions  datatypes.JSON `gorm:"type:json" json:"warranty_restrictions"`
	GlobalServiceProgram  bool           `gorm:"type:boolean;default:false" json:"global_service_program"`

	// Cross-Border Insurance
	CrossBorderCoverage bool           `gorm:"type:boolean;default:false" json:"cross_border_coverage"`
	CoveredCountries    datatypes.JSON `gorm:"type:json" json:"covered_countries"`  // []string
	ExcludedCountries   datatypes.JSON `gorm:"type:json" json:"excluded_countries"` // []string
	TerritorialLimits   datatypes.JSON `gorm:"type:json" json:"territorial_limits"`

	// Import/Export Information
	ImportDuty    float64 `gorm:"type:decimal(15,2)" json:"import_duty"`
	ExportDuty    float64 `gorm:"type:decimal(15,2)" json:"export_duty"`
	CustomsValue  float64 `gorm:"type:decimal(15,2)" json:"customs_value"`
	HSCode        string  `gorm:"type:varchar(20)" json:"hs_code"`
	ImportLicense string  `gorm:"type:varchar(100)" json:"import_license"`
	ExportLicense string  `gorm:"type:varchar(100)" json:"export_license"`

	// Regional Pricing
	RegionalPrices       datatypes.JSON `gorm:"type:json" json:"regional_prices"` // map[string]Price
	PriceVariations      datatypes.JSON `gorm:"type:json" json:"price_variations"`
	MarketSpecificModels datatypes.JSON `gorm:"type:json" json:"market_specific_models"`

	// Currency Conversion History
	ConversionHistory  datatypes.JSON `gorm:"type:json" json:"conversion_history"` // []ConversionRecord
	LastConversionDate *time.Time     `gorm:"type:timestamp" json:"last_conversion_date,omitempty"`

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceDocumentManagement represents document and receipt management (excluding media)
type DeviceDocumentManagement struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Ownership Documents
	ProofOfPurchase      datatypes.JSON `gorm:"type:json" json:"proof_of_purchase"` // Document
	Receipt              datatypes.JSON `gorm:"type:json" json:"receipt"`           // Document
	Invoice              datatypes.JSON `gorm:"type:json" json:"invoice"`           // Document
	OwnershipCertificate datatypes.JSON `gorm:"type:json" json:"ownership_certificate"`
	TransferDocuments    datatypes.JSON `gorm:"type:json" json:"transfer_documents"` // []Document

	// Inspection Documents
	InspectionReports datatypes.JSON `gorm:"type:json" json:"inspection_reports"` // []Report
	DamageAssessments datatypes.JSON `gorm:"type:json" json:"damage_assessments"` // []Assessment
	RepairEstimates   datatypes.JSON `gorm:"type:json" json:"repair_estimates"`   // []Estimate
	ConditionReports  datatypes.JSON `gorm:"type:json" json:"condition_reports"`  // []Report

	// Insurance Documents
	PolicyDocuments datatypes.JSON `gorm:"type:json" json:"policy_documents"` // []Document
	ClaimDocuments  datatypes.JSON `gorm:"type:json" json:"claim_documents"`  // []Document
	CoverageProofs  datatypes.JSON `gorm:"type:json" json:"coverage_proofs"`  // []Document

	// Service Documents
	ServiceHistory     datatypes.JSON `gorm:"type:json" json:"service_history"`     // []ServiceRecord
	MaintenanceRecords datatypes.JSON `gorm:"type:json" json:"maintenance_records"` // []Record
	RepairInvoices     datatypes.JSON `gorm:"type:json" json:"repair_invoices"`     // []Invoice

	// Warranty Documents
	WarrantyCard     datatypes.JSON `gorm:"type:json" json:"warranty_card"`
	ExtendedWarranty datatypes.JSON `gorm:"type:json" json:"extended_warranty"`
	WarrantyClaims   datatypes.JSON `gorm:"type:json" json:"warranty_claims"` // []Claim

	// Legal Documents
	LegalNotices        datatypes.JSON `gorm:"type:json" json:"legal_notices"`        // []Notice
	ComplianceDocuments datatypes.JSON `gorm:"type:json" json:"compliance_documents"` // []Document
	RegulatoryFilings   datatypes.JSON `gorm:"type:json" json:"regulatory_filings"`   // []Filing

	// Document Metadata
	TotalDocuments    int            `gorm:"type:int" json:"total_documents"`
	LastDocumentAdded *time.Time     `gorm:"type:timestamp" json:"last_document_added,omitempty"`
	DocumentRetention datatypes.JSON `gorm:"type:json" json:"document_retention"` // RetentionPolicy

	// Status
	Status    string    `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsAuthentic checks if device is verified as authentic
func (dsc *DeviceSupplyChain) IsAuthentic() bool {
	return dsc.AuthenticityStatus == "genuine" && dsc.AuthenticityScore >= 90
}

// IsGrayMarket checks if device is from gray market
func (dsc *DeviceSupplyChain) IsGrayMarket() bool {
	return dsc.MarketType == "gray" || dsc.ParallelImportStatus
}

// IsHighRisk checks if supply chain has high vulnerability
func (dsc *DeviceSupplyChain) IsHighRisk() bool {
	return dsc.VulnerabilityScore > 70 || dsc.AuthenticityStatus == "counterfeit"
}

// HasInternationalCoverage checks if device has international coverage
func (dmc *DeviceMultiCurrency) HasInternationalCoverage() bool {
	return dmc.InternationalWarranty || dmc.CrossBorderCoverage
}

// GetConversionRate calculates current conversion rate
func (dmc *DeviceMultiCurrency) GetConversionRate(targetCurrency string) float64 {
	if dmc.LocalCurrency == targetCurrency {
		return 1.0
	}
	// This would typically call an external service
	return dmc.ExchangeRate
}

// HasCompleteDocumentation checks if all required documents are present
func (ddm *DeviceDocumentManagement) HasCompleteDocumentation() bool {
	return ddm.ProofOfPurchase != nil && ddm.Receipt != nil &&
		ddm.OwnershipCertificate != nil && ddm.WarrantyCard != nil
}
