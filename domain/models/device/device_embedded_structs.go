package device

import (
	"time"

	"github.com/google/uuid"
)

// ============================================
// DEVICE EMBEDDED STRUCTS
// ============================================

// DeviceIdentification contains core device identifiers
type DeviceIdentification struct {
	IMEI         string `gorm:"uniqueIndex;not null" json:"imei"`
	IMEI2        string `json:"imei2"` // For dual-SIM devices
	SerialNumber string `gorm:"uniqueIndex;omitempty" json:"serial_number"`
	ModelNumber  string `json:"model_number"`
}

// DeviceClassification contains device classification information
type DeviceClassification struct {
	Brand        string         `gorm:"not null" json:"brand"`
	Model        string         `gorm:"not null" json:"model"`
	Manufacturer string         `gorm:"not null" json:"manufacturer"`
	Category     DeviceCategory `gorm:"type:varchar(30)" json:"category"`
	Segment      string         `json:"segment"` // flagship, premium, mid-range, budget
}

// DeviceSpecifications contains technical specifications (alias for backward compatibility)
type DeviceSpecifications struct {
	DeviceSpecification
}

// DeviceSpecification contains technical specifications
type DeviceSpecification struct {
	Brand           string          `gorm:"not null" json:"brand"`
	Model           string          `gorm:"not null" json:"model"`
	Manufacturer    string          `gorm:"not null" json:"manufacturer"`
	Category        DeviceCategory  `gorm:"type:varchar(50);default:'smartphone'" json:"category"`
	OperatingSystem string          `json:"operating_system"`
	OSVersion       string          `json:"os_version"`
	StorageCapacity int             `json:"storage_capacity"` // in GB
	RAM             int             `json:"ram"`              // in GB
	Color           string          `json:"color"`
	ScreenSize      float64         `json:"screen_size"` // in inches
	ScreenType      string          `json:"screen_type"` // OLED, AMOLED, LCD
	WaterResistance WaterResistance `json:"water_resistance"`
	BiometricType   BiometricType   `json:"biometric_type"`
}

// DevicePhysicalCondition contains physical condition details
type DevicePhysicalCondition struct {
	Condition        DeviceCondition  `gorm:"type:varchar(20);default:'good'" json:"condition"`
	Grade            DeviceGrade      `gorm:"type:varchar(5)" json:"grade"`
	ScreenCondition  ScreenCondition  `json:"screen_condition"`
	BodyCondition    BodyCondition    `json:"body_condition"`
	BatteryHealth    int              `json:"battery_health"` // percentage
	BatteryCycles    int              `json:"battery_cycles"`
	FunctionalIssues FunctionalIssues `gorm:"type:jsonb" json:"functional_issues"`
	LastInspection   *time.Time       `json:"last_inspection"`
	InspectionNotes  string           `json:"inspection_notes"`
	OriginalBox      bool             `gorm:"default:false" json:"original_box"`
	OriginalReceipt  bool             `gorm:"default:false" json:"original_receipt"`
	Accessories      []string         `gorm:"type:jsonb" json:"accessories"` // List of included accessories
}

// DeviceOwnership contains ownership information
type DeviceOwnership struct {
	OwnerID            uuid.UUID     `gorm:"type:uuid;not null;index" json:"owner_id"`
	OwnershipType      OwnershipType `gorm:"type:varchar(20);default:'personal'" json:"ownership_type"`
	CorporateAccountID *uuid.UUID    `gorm:"type:uuid;index" json:"corporate_account_id"`
	DepartmentID       *uuid.UUID    `gorm:"type:uuid;index" json:"department_id"`
	EmployeeID         *uuid.UUID    `gorm:"type:uuid;index" json:"employee_id"`
	RegistrationDate   time.Time     `gorm:"autoCreateTime" json:"registration_date"`
	TransferHistory    []uuid.UUID   `gorm:"type:jsonb" json:"transfer_history"` // Previous owner IDs
}

// DeviceFinancial contains financial information
type DeviceFinancial struct {
	PurchaseDate     *time.Time `json:"purchase_date"`
	PurchasePrice    Money      `gorm:"embedded;embeddedPrefix:purchase_" json:"purchase_price"`
	CurrentValue     Money      `gorm:"embedded;embeddedPrefix:current_" json:"current_value"`
	MarketValue      Money      `gorm:"embedded;embeddedPrefix:market_" json:"market_value"`
	DepreciatedValue Money      `gorm:"embedded;embeddedPrefix:depreciated_" json:"depreciated_value"`
	TradeInValue     Money      `gorm:"embedded;embeddedPrefix:tradein_" json:"trade_in_value"`
	SalvageValue     Money      `gorm:"embedded;embeddedPrefix:salvage_" json:"salvage_value"`
	LastValuation    *time.Time `json:"last_valuation"`
	ValuationMethod  string     `json:"valuation_method"` // market, depreciation, appraisal

	// Settlement References
	SettlementReference   string `json:"settlement_reference,omitempty"`    // Settlement transaction reference
	PaymentReference      string `json:"payment_reference,omitempty"`       // Payment processing reference
	CheckNumber           string `json:"check_number,omitempty"`            // Physical check number
	WireTransferReference string `json:"wire_transfer_reference,omitempty"` // Wire transfer reference

	// Escrow & Trust Accounts
	EscrowAccountNumber string `json:"escrow_account_number,omitempty"` // Escrow account reference
	TrustAccountNumber  string `json:"trust_account_number,omitempty"`  // Trust account reference

	// Banking References
	BankReferenceNumber    string `json:"bank_reference_number,omitempty"`    // Bank processing reference
	ACHReference           string `json:"ach_reference,omitempty"`            // ACH transaction reference
	CreditCardAuthCode     string `json:"credit_card_auth_code,omitempty"`    // Card authorization code
	DigitalWalletReference string `json:"digital_wallet_reference,omitempty"` // Digital wallet reference

	// Financial Institution Identifiers
	BankRoutingNumber string `json:"bank_routing_number,omitempty"` // ABA routing number
	SWIFTCode         string `json:"swift_code,omitempty"`          // SWIFT/BIC code
	IBAN              string `json:"iban,omitempty"`                // International bank account number
	BSBCode           string `json:"bsb_code,omitempty"`            // Australian BSB code
}

// DeviceStatusInfo contains status and state information
type DeviceStatusInfo struct {
	Status         DeviceStatus    `gorm:"type:varchar(30);not null;default:'active';index" json:"status"`
	IsStolen       bool            `gorm:"default:false;index" json:"is_stolen"`
	StolenDate     *time.Time      `json:"stolen_date"`
	StolenLocation *DeviceLocation `gorm:"type:jsonb" json:"stolen_location"`
	StolenReportID string          `json:"stolen_report_id"`
	IsLost         bool            `gorm:"default:false;index" json:"is_lost"`
	LostDate       *time.Time      `json:"lost_date"`
	IsLocked       bool            `gorm:"default:false" json:"is_locked"`
	LockedReason   string          `json:"locked_reason"`
	IsRetired      bool            `gorm:"default:false" json:"is_retired"`
	RetiredDate    *time.Time      `json:"retired_date"`
	RetiredReason  string          `json:"retired_reason"`

	// Court & Legal References
	CourtCaseNumber    string `json:"court_case_number,omitempty"`    // Court case reference
	LegalHoldReference string `json:"legal_hold_reference,omitempty"` // Legal hold tracking number
	SubpoenaReference  string `json:"subpoena_reference,omitempty"`   // Subpoena reference number
	DiscoveryReference string `json:"discovery_reference,omitempty"`  // Discovery reference
}

// DeviceSecurity contains security-related information
type DeviceSecurity struct {
	IsLocked            bool   `gorm:"default:false" json:"is_locked"`
	LockType            string `json:"lock_type"` // pin, pattern, fingerprint, face
	FindMyDeviceEnabled bool   `gorm:"default:false" json:"find_my_device_enabled"`
	RemoteWipeEnabled   bool   `gorm:"default:false" json:"remote_wipe_enabled"`
	EncryptionEnabled   bool   `gorm:"default:false" json:"encryption_enabled"`
	BiometricEnabled    bool   `gorm:"default:false" json:"biometric_enabled"`
	TwoFactorEnabled    bool   `gorm:"default:false" json:"two_factor_enabled"`
}

// DeviceVerification contains verification and authenticity information
type DeviceVerification struct {
	IsVerified         bool               `gorm:"default:false;index" json:"is_verified"`
	VerificationDate   *time.Time         `json:"verification_date"`
	VerifiedBy         *uuid.UUID         `gorm:"type:uuid" json:"verified_by"`
	VerificationMethod string             `json:"verification_method"`
	AuthenticityStatus AuthenticityStatus `json:"authenticity_status"`
	GreyMarket         bool               `gorm:"default:false" json:"grey_market"`
	CounterfeitRisk    float64            `json:"counterfeit_risk"` // 0-100

	// Insurance Verification Codes
	InsuranceVerificationCode string `json:"insurance_verification_code,omitempty"` // Verification PIN/code
	PolicyVerificationToken   string `json:"policy_verification_token,omitempty"`   // Digital verification token
	ClaimVerificationCode     string `json:"claim_verification_code,omitempty"`     // Claim verification code

	// Digital Verification
	BlockchainCertificateID string `json:"blockchain_certificate_id,omitempty"` // Blockchain verification
	DigitalSignatureID      string `json:"digital_signature_id,omitempty"`      // Digital signature reference
	TimestampAuthorityID    string `json:"timestamp_authority_id,omitempty"`    // Trusted timestamp reference
}

// DeviceNetwork contains network and connectivity information
type DeviceNetwork struct {
	NetworkOperator string        `json:"network_operator"`
	NetworkStatus   NetworkStatus `json:"network_status"`
	NetworkLocked   bool          `gorm:"default:false" json:"network_locked"`
	SIMCount        int           `json:"sim_count"`
	ESIMSupport     bool          `gorm:"default:false" json:"esim_support"`
	FiveGCapable    bool          `gorm:"default:false" json:"five_g_capable"`
	LastOnline      *time.Time    `json:"last_online"`
	IPAddress       string        `json:"ip_address"`
	MACAddress      string        `json:"mac_address"`
}

// DeviceRiskAssessment contains risk evaluation data
type DeviceRiskAssessment struct {
	RiskScore       float64         `json:"risk_score"` // 0-100
	RiskLevel       RiskLevel       `json:"risk_level"`
	TheftRiskLevel  RiskLevel       `json:"theft_risk_level"`
	DamageRiskLevel RiskLevel       `json:"damage_risk_level"`
	FraudRiskScore  float64         `json:"fraud_risk_score"` // 0-100
	BlacklistStatus BlacklistStatus `json:"blacklist_status"`
	BlacklistDate   *time.Time      `json:"blacklist_date"`
	BlacklistReason string          `json:"blacklist_reason"`
	RiskFactors     []string        `gorm:"type:jsonb" json:"risk_factors"`
	LastRiskUpdate  *time.Time      `json:"last_risk_update"`
}

// DeviceWarranty contains warranty information
type DeviceWarranty struct {
	WarrantyStatus     string     `json:"warranty_status"` // active, expired, void
	WarrantyStartDate  *time.Time `json:"warranty_start_date"`
	WarrantyExpiry     *time.Time `json:"warranty_expiry"`
	WarrantyProvider   string     `json:"warranty_provider"`
	WarrantyType       string     `json:"warranty_type"` // manufacturer, extended, third_party
	ExtendedWarrantyID *uuid.UUID `gorm:"type:uuid" json:"extended_warranty_id"`
	WarrantyClaims     int        `json:"warranty_claims"`
	LastWarrantyClaim  *time.Time `json:"last_warranty_claim"`
}

// DeviceLocation contains location tracking data
type DeviceLocationInfo struct {
	CurrentLocation    *DeviceLocation  `gorm:"type:jsonb" json:"current_location"`
	HomeLocation       *DeviceLocation  `gorm:"type:jsonb" json:"home_location"`
	WorkLocation       *DeviceLocation  `gorm:"type:jsonb" json:"work_location"`
	LocationHistory    []DeviceLocation `gorm:"type:jsonb" json:"location_history"`
	GeofenceEnabled    bool             `gorm:"default:false" json:"geofence_enabled"`
	GeofenceViolations int              `json:"geofence_violations"`
	LastLocationUpdate *time.Time       `json:"last_location_update"`
}

// DeviceCorporateManagement contains corporate/MDM information
type DeviceCorporateManagement struct {
	MDMEnrolled         bool        `gorm:"default:false;index" json:"mdm_enrolled"`
	MDMProvider         MDMProvider `json:"mdm_provider"`
	MDMProfileID        string      `json:"mdm_profile_id"`
	MDMCompliant        bool        `gorm:"default:false" json:"mdm_compliant"`
	ComplianceStatus    string      `json:"compliance_status"`
	LastComplianceCheck *time.Time  `json:"last_compliance_check"`
	SecurityPolicies    []string    `gorm:"type:jsonb" json:"security_policies"`
	AppRestrictions     []string    `gorm:"type:jsonb" json:"app_restrictions"`
	DataContainerized   bool        `gorm:"default:false" json:"data_containerized"`
	RemoteWipeEnabled   bool        `gorm:"default:false" json:"remote_wipe_enabled"`
	LastSync            *time.Time  `json:"last_sync"`
}

// DeviceUsageMetrics contains usage statistics
type DeviceUsageMetrics struct {
	DailyUsageHours  float64    `json:"daily_usage_hours"`
	TotalScreenTime  int        `json:"total_screen_time"` // in hours
	AppInstallCount  int        `json:"app_install_count"`
	DataUsageGB      float64    `json:"data_usage_gb"`
	CallMinutes      int        `json:"call_minutes"`
	MessageCount     int        `json:"message_count"`
	PhotoCount       int        `json:"photo_count"`
	VideoCount       int        `json:"video_count"`
	LastBackup       *time.Time `json:"last_backup"`
	CloudStorageUsed float64    `json:"cloud_storage_used"` // in GB
}

// DeviceLifecycle contains lifecycle management information
type DeviceLifecycleInfo struct {
	LifecycleStage     string     `json:"lifecycle_stage"` // new, active, maintenance, end_of_life
	ManufactureDate    *time.Time `json:"manufacture_date"`
	FirstActivation    *time.Time `json:"first_activation"`
	ExpectedLifespan   int        `json:"expected_lifespan"` // in months
	RefreshEligible    bool       `gorm:"default:false" json:"refresh_eligible"`
	RefreshDate        *time.Time `json:"refresh_date"`
	DisposalDate       *time.Time `json:"disposal_date"`
	DisposalMethod     string     `json:"disposal_method"` // recycled, donated, destroyed
	EnvironmentalScore float64    `json:"environmental_score"`
}

// DeviceInsurance contains insurance-specific information
type DeviceInsurance struct {
	InsuranceEligible  bool       `gorm:"default:true" json:"insurance_eligible"`
	InsuranceValue     Money      `gorm:"embedded;embeddedPrefix:insurance_" json:"insurance_value"`
	MaxCoverage        Money      `gorm:"embedded;embeddedPrefix:max_coverage_" json:"max_coverage"`
	RecommendedPremium Money      `gorm:"embedded;embeddedPrefix:recommended_premium_" json:"recommended_premium"`
	ActivePolicies     int        `json:"active_policies"`
	TotalClaims        int        `json:"total_claims"`
	LastClaimDate      *time.Time `json:"last_claim_date"`
	ClaimFrequency     float64    `json:"claim_frequency"`
	PreExistingDamage  []string   `gorm:"type:jsonb" json:"pre_existing_damage"`
	ExclusionReasons   []string   `gorm:"type:jsonb" json:"exclusion_reasons"`
	SpecialConditions  []string   `gorm:"type:jsonb" json:"special_conditions"`

	// Policy Reference Numbers
	PolicyNumber       string `gorm:"index" json:"policy_number"`      // Primary policy number
	CertificateNumber  string `gorm:"index" json:"certificate_number"` // Individual certificate
	MasterPolicyNumber string `json:"master_policy_number,omitempty"`  // For group/master policies
	BinderNumber       string `json:"binder_number,omitempty"`         // Temporary policy reference
	QuoteNumber        string `gorm:"index" json:"quote_number"`       // Original quote reference

	// Underwriting References
	UnderwritingReference string `json:"underwriting_reference,omitempty"` // Underwriter tracking
	RiskAssessmentID      string `json:"risk_assessment_id,omitempty"`     // Risk assessment reference
	ActuarialReference    string `json:"actuarial_reference,omitempty"`    // Actuarial review reference

	// Regulatory Filings
	RateFilingReference string `json:"rate_filing_reference,omitempty"` // Insurance department filing
	FormFilingReference string `json:"form_filing_reference,omitempty"` // Form approval reference

	// Active Claim References
	ActiveClaimNumber     string `json:"active_claim_number,omitempty"`     // Currently open claim
	LastClaimNumber       string `json:"last_claim_number,omitempty"`       // Most recent claim reference
	ClaimReserveReference string `json:"claim_reserve_reference,omitempty"` // Reserve account reference

	// Incident Reporting
	IncidentReportNumber string `json:"incident_report_number,omitempty"` // Police/insurance report
	FNOLReference        string `json:"fnol_reference,omitempty"`         // First Notice of Loss reference
	LossControlReference string `json:"loss_control_reference,omitempty"` // Loss control case reference
}

// DeviceMaintenanceHistory contains maintenance and service records
type DeviceMaintenanceHistory struct {
	LastService         *time.Time `json:"last_service"`
	NextServiceDue      *time.Time `json:"next_service_due"`
	ServiceCount        int        `json:"service_count"`
	RepairCount         int        `json:"repair_count"`
	TotalRepairCost     Money      `gorm:"embedded;embeddedPrefix:total_repair_" json:"total_repair_cost"`
	LastRepairDate      *time.Time `json:"last_repair_date"`
	RepairHistory       []string   `gorm:"type:jsonb" json:"repair_history"`
	MaintenanceSchedule string     `json:"maintenance_schedule"` // monthly, quarterly, annual
	UnderMaintenance    bool       `gorm:"default:false" json:"under_maintenance"`
	MaintenanceContract *uuid.UUID `gorm:"type:uuid" json:"maintenance_contract_id"`
}

// DeviceDocumentation contains documentation and evidence information
type DeviceDocumentation struct {
	PhotoURLs          []string   `gorm:"type:jsonb" json:"photo_urls"`
	InvoiceURL         string     `json:"invoice_url"`
	ReceiptURL         string     `json:"receipt_url"`
	ProofOfPurchaseURL string     `json:"proof_of_purchase_url"`
	InspectionReports  []string   `gorm:"type:jsonb" json:"inspection_reports"`
	CertificationDocs  []string   `gorm:"type:jsonb" json:"certification_docs"`
	WarrantyDocs       []string   `gorm:"type:jsonb" json:"warranty_docs"`
	LastInspectionDate *time.Time `json:"last_inspection_date"`

	// Legal Document References
	AffidavitReference    string `json:"affidavit_reference,omitempty"`     // Sworn statement reference
	DepositionReference   string `json:"deposition_reference,omitempty"`    // Deposition reference
	ExpertReportReference string `json:"expert_report_reference,omitempty"` // Expert report reference

	// Regulatory Filings
	SECReference          string `json:"sec_reference,omitempty"`            // SEC filing reference
	FINRAReference        string `json:"finra_reference,omitempty"`          // FINRA reference
	StateInsuranceDeptRef string `json:"state_insurance_dept_ref,omitempty"` // State department reference
}

// DeviceCompliance contains regulatory compliance information
type DeviceCompliance struct {
	RegulatoryCompliant bool       `gorm:"default:true" json:"regulatory_compliant"`
	ComplianceRegion    string     `json:"compliance_region"`
	CertificationIDs    []string   `gorm:"type:jsonb" json:"certification_ids"`
	ImportCompliant     bool       `gorm:"default:true" json:"import_compliant"`
	ExportRestricted    bool       `gorm:"default:false" json:"export_restricted"`
	DataPrivacyLevel    string     `json:"data_privacy_level"` // low, medium, high, critical
	EncryptionEnabled   bool       `gorm:"default:true" json:"encryption_enabled"`
	EncryptionType      string     `json:"encryption_type"`
	SecureBootEnabled   bool       `gorm:"default:false" json:"secure_boot_enabled"`
	LastSecurityPatch   *time.Time `json:"last_security_patch"`

	// Regulatory Filings
	ComplianceReference  string `json:"compliance_reference,omitempty"`   // Compliance case reference
	AuditReference       string `json:"audit_reference,omitempty"`        // Audit tracking reference
	RegulatoryCaseNumber string `json:"regulatory_case_number,omitempty"` // Regulatory filing number
}

// DeviceMetadata contains additional metadata
type DeviceMetadata struct {
	Source          string                 `json:"source"` // manual, import, api, migration
	ImportBatch     string                 `json:"import_batch"`
	ExternalID      string                 `json:"external_id"`
	Tags            []string               `gorm:"type:jsonb" json:"tags"`
	CustomFields    map[string]interface{} `gorm:"type:jsonb" json:"custom_fields"`
	Notes           string                 `json:"notes"`
	InternalNotes   string                 `json:"internal_notes"`
	LastModifiedBy  *uuid.UUID             `gorm:"type:uuid" json:"last_modified_by"`
	ModificationLog []string               `gorm:"type:jsonb" json:"modification_log"`
}

// DeviceAudit contains audit trail information
type DeviceAudit struct {
	CreatedBy      uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	CreatedByName  string     `json:"created_by_name"`
	ModifiedBy     *uuid.UUID `gorm:"type:uuid" json:"modified_by"`
	ModifiedByName string     `json:"modified_by_name"`
	ApprovedBy     *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedByName string     `json:"approved_by_name"`
	ApprovalDate   *time.Time `json:"approval_date"`
	AuditTrail     []string   `gorm:"type:jsonb" json:"audit_trail"`
	Version        int        `gorm:"default:1" json:"version"`
	ChangeHistory  []string   `gorm:"type:jsonb" json:"change_history"`
}
