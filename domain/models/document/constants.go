package document

// Document Types
const (
	TypePolicy         = "policy"
	TypeClaim          = "claim"
	TypeInvoice        = "invoice"
	TypeReceipt        = "receipt"
	TypeContract       = "contract"
	TypeIDProof        = "id_proof"
	TypeDamagePhoto    = "damage_photo"
	TypeRepairReport   = "repair_report"
	TypeValuation      = "valuation_report"
	TypeInspection     = "inspection_report"
	TypeLegalDocument  = "legal_document"
	TypeCertificate    = "certificate"
	TypeCorrespondence = "correspondence"
)

// Document Categories
const (
	CategoryLegal          = "legal"
	CategoryFinancial      = "financial"
	CategoryEvidence       = "evidence"
	CategoryIdentification = "identification"
	CategoryClaim          = "claim"
	CategoryPolicy         = "policy"
	CategoryCompliance     = "compliance"
	CategoryOperational    = "operational"
)

// Storage Providers
const (
	StorageS3    = "s3"
	StorageAzure = "azure"
	StorageGCS   = "gcs"
	StorageLocal = "local"
)

// Signature Types
const (
	SignatureDrawn              = "drawn"
	SignatureTyped              = "typed"
	SignatureUploaded           = "uploaded"
	SignatureDigitalCertificate = "digital_certificate"
	SignatureBiometric          = "biometric"
)

// Signer Roles
const (
	SignerPolicyholder = "policyholder"
	SignerAgent        = "agent"
	SignerWitness      = "witness"
	SignerNotary       = "notary"
	SignerManager      = "manager"
	SignerTechnician   = "technician"
	SignerClaimant     = "claimant"
)

// Signature Status
const (
	SignaturePending  = "pending"
	SignatureSigned   = "signed"
	SignatureDeclined = "declined"
	SignatureExpired  = "expired"
	SignatureRevoked  = "revoked"
)

// Verification Methods
const (
	VerificationEmail     = "email"
	VerificationSMS       = "sms"
	VerificationBiometric = "biometric"
	VerificationPassword  = "password"
	VerificationOTP       = "otp"
)

// Document Generation Types
const (
	GenerationInstant   = "instant"
	GenerationScheduled = "scheduled"
	GenerationBatch     = "batch"
)

// Document Generation Status
const (
	GenerationPending    = "pending"
	GenerationProcessing = "processing"
	GenerationCompleted  = "completed"
	GenerationFailed     = "failed"
)

// Output Formats
const (
	FormatPDF  = "pdf"
	FormatHTML = "html"
	FormatDOCX = "docx"
	FormatJSON = "json"
	FormatXML  = "xml"
)

// Delivery Methods
const (
	DeliveryEmail    = "email"
	DeliveryDownload = "download"
	DeliveryAPI      = "api"
	DeliverySMS      = "sms"
	DeliveryWebhook  = "webhook"
)

// Access Types
const (
	AccessView     = "view"
	AccessDownload = "download"
	AccessEdit     = "edit"
	AccessSign     = "sign"
	AccessShare    = "share"
	AccessDelete   = "delete"
	AccessPrint    = "print"
)

// Access Permissions
const (
	PermissionOwner  = "owner"
	PermissionEditor = "editor"
	PermissionViewer = "viewer"
	PermissionSigner = "signer"
	PermissionAdmin  = "admin"
)

// Content Types
const (
	ContentHTML     = "html"
	ContentMarkdown = "markdown"
	ContentPDF      = "pdf"
	ContentPlain    = "plain"
)

// Template Categories
const (
	TemplateCategoryPolicy    = "policy"
	TemplateCategoryClaim     = "claim"
	TemplateCategoryFinancial = "financial"
	TemplateCategoryLegal     = "legal"
	TemplateCategoryNotice    = "notice"
)

// Document Security Levels
const (
	SecurityPublic       = "public"
	SecurityInternal     = "internal"
	SecurityConfidential = "confidential"
	SecurityRestricted   = "restricted"
)

// Document Status
const (
	StatusDraft      = "draft"
	StatusPending    = "pending"
	StatusActive     = "active"
	StatusArchived   = "archived"
	StatusExpired    = "expired"
	StatusSuperseded = "superseded"
)

// Verification Status
const (
	VerificationNotRequired = "not_required"
	VerificationPending     = "pending"
	VerificationInProgress  = "in_progress"
	VerificationCompleted   = "completed"
	VerificationFailed      = "failed"
)

// Signature Providers
const (
	ProviderDocuSign  = "docusign"
	ProviderHelloSign = "hellosign"
	ProviderAdobeSign = "adobe_sign"
	ProviderPandaDoc  = "pandadoc"
	ProviderInternal  = "internal"
)

// Document Retention Periods (days)
const (
	RetentionTemporary  = 30
	RetentionShortTerm  = 90
	RetentionMediumTerm = 365
	RetentionLongTerm   = 2555  // 7 years
	RetentionPermanent  = 36500 // 100 years
)

// File Size Limits (bytes)
const (
	MaxFileSize      = 52428800 // 50MB
	MaxImageSize     = 10485760 // 10MB
	MaxThumbnailSize = 1048576  // 1MB
	MaxSignatureSize = 524288   // 512KB
)

// Default Values
const (
	DefaultLanguage        = "en"
	DefaultVersion         = 1
	DefaultMaxVersions     = 10
	DefaultAccessExpiry    = 30 // days
	DefaultSignatureExpiry = 7  // days
)
