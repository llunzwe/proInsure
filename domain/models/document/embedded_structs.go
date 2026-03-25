package document

import (
	"time"
	
	"github.com/google/uuid"
)

// DocumentMetadata contains core document metadata
type DocumentMetadata struct {
	Type        string `gorm:"not null" json:"type"` // policy, claim, invoice, etc.
	Category    string `json:"category"`             // legal, financial, evidence, etc.
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	MimeType    string `gorm:"not null" json:"mime_type"`
	Size        int64  `json:"size"`           // bytes
	Hash        string `json:"hash,omitempty"` // SHA256 hash for integrity
	Version     int    `gorm:"default:1" json:"version"`
	Tags        string `json:"tags,omitempty"`     // JSON array of tags
	Metadata    string `json:"metadata,omitempty"` // JSON for additional data
}

// DocumentStorage contains storage-related information
type DocumentStorage struct {
	StorageProvider string `gorm:"not null" json:"storage_provider"` // s3, azure, gcs, local
	StoragePath     string `gorm:"not null" json:"storage_path"`
	StorageBucket   string `json:"storage_bucket,omitempty"`
	URL             string `json:"url,omitempty"`
	ThumbnailURL    string `json:"thumbnail_url,omitempty"`
	EncryptionKey   string `json:"-"` // Encrypted document key
}

// DocumentSecurity contains security and verification information
type DocumentSecurity struct {
	IsPublic           bool       `gorm:"default:false" json:"is_public"`
	SecurityLevel      string     `json:"security_level"` // public, internal, confidential, restricted
	IsVerified         bool       `gorm:"default:false" json:"is_verified"`
	VerifiedBy         *uuid.UUID `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	VerificationStatus string     `json:"verification_status"` // not_required, pending, completed, failed
	VerificationNotes  string     `json:"verification_notes,omitempty"`
}

// DocumentLifecycle contains lifecycle management information
type DocumentLifecycle struct {
	Status          string     `json:"status"` // draft, pending, active, archived, expired
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	RetentionPeriod int        `json:"retention_period"` // days
	ArchivedAt      *time.Time `json:"archived_at,omitempty"`
	ArchivedBy      *uuid.UUID `gorm:"type:uuid" json:"archived_by,omitempty"`
	DisposalDate    *time.Time `json:"disposal_date,omitempty"`
	LegalHold       bool       `gorm:"default:false" json:"legal_hold"`
	LegalHoldReason string     `json:"legal_hold_reason,omitempty"`
}

// DocumentRelationships contains references to related entities
type DocumentRelationships struct {
	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	PolicyID         *uuid.UUID `gorm:"type:uuid" json:"policy_id,omitempty"`
	ClaimID          *uuid.UUID `gorm:"type:uuid" json:"claim_id,omitempty"`
	PaymentID        *uuid.UUID `gorm:"type:uuid" json:"payment_id,omitempty"`
	DeviceID         *uuid.UUID `gorm:"type:uuid" json:"device_id,omitempty"`
	RepairID         *uuid.UUID `gorm:"type:uuid" json:"repair_id,omitempty"`
	VendorID         *uuid.UUID `gorm:"type:uuid" json:"vendor_id,omitempty"`
	TechnicianID     *uuid.UUID `gorm:"type:uuid" json:"technician_id,omitempty"`
	CorporateID      *uuid.UUID `gorm:"type:uuid" json:"corporate_id,omitempty"`
	ParentDocumentID *uuid.UUID `gorm:"type:uuid" json:"parent_document_id,omitempty"`
}

// DocumentCompliance contains compliance and regulatory information
type DocumentCompliance struct {
	ComplianceChecked  bool       `gorm:"default:false" json:"compliance_checked"`
	ComplianceStatus   string     `json:"compliance_status"` // compliant, non_compliant, pending_review
	ComplianceNotes    string     `json:"compliance_notes,omitempty"`
	RegulatoryStandard string     `json:"regulatory_standard,omitempty"` // GDPR, HIPAA, SOC2, etc.
	AuditTrail         string     `json:"audit_trail,omitempty"`         // JSON audit log
	LastAuditDate      *time.Time `json:"last_audit_date,omitempty"`
	NextAuditDate      *time.Time `json:"next_audit_date,omitempty"`
}

// DocumentProcessing contains document processing information
type DocumentProcessing struct {
	OCRProcessed     bool       `gorm:"default:false" json:"ocr_processed"`
	OCRText          string     `json:"ocr_text,omitempty"` // Extracted text
	OCRConfidence    float64    `json:"ocr_confidence"`     // OCR confidence score
	OCRLanguage      string     `json:"ocr_language,omitempty"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
	ProcessingStatus string     `json:"processing_status"`           // pending, processing, completed, failed
	ProcessingErrors string     `json:"processing_errors,omitempty"` // JSON array of errors
	AIClassification string     `json:"ai_classification,omitempty"` // AI-determined document type
	AIConfidence     float64    `json:"ai_confidence"`               // AI confidence score
	ExtractedData    string     `json:"extracted_data,omitempty"`    // JSON of extracted fields
}

// DocumentWorkflow contains workflow and approval information
type DocumentWorkflow struct {
	WorkflowID       *uuid.UUID `gorm:"type:uuid" json:"workflow_id,omitempty"`
	WorkflowStatus   string     `json:"workflow_status"` // not_started, in_progress, completed, rejected
	CurrentStep      string     `json:"current_step,omitempty"`
	TotalSteps       int        `json:"total_steps"`
	ApprovalRequired bool       `gorm:"default:false" json:"approval_required"`
	ApprovalStatus   string     `json:"approval_status"` // pending, approved, rejected
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	RejectedBy       *uuid.UUID `gorm:"type:uuid" json:"rejected_by,omitempty"`
	RejectedAt       *time.Time `json:"rejected_at,omitempty"`
	RejectionReason  string     `json:"rejection_reason,omitempty"`
	WorkflowMetadata string     `json:"workflow_metadata,omitempty"` // JSON for workflow data
}

// DocumentSharing contains sharing and collaboration information
type DocumentSharing struct {
	IsShared           bool       `gorm:"default:false" json:"is_shared"`
	ShareCount         int        `gorm:"default:0" json:"share_count"`
	SharedWith         string     `json:"shared_with,omitempty"` // JSON array of user IDs
	ShareExpiresAt     *time.Time `json:"share_expires_at,omitempty"`
	PublicShareLink    string     `json:"public_share_link,omitempty"`
	SharePassword      string     `json:"-"` // Encrypted share password
	DownloadCount      int        `gorm:"default:0" json:"download_count"`
	ViewCount          int        `gorm:"default:0" json:"view_count"`
	LastViewedAt       *time.Time `json:"last_viewed_at,omitempty"`
	LastDownloadedAt   *time.Time `json:"last_downloaded_at,omitempty"`
	CollaborationNotes string     `json:"collaboration_notes,omitempty"` // JSON array of notes
}

// DocumentAnalytics contains analytics and metrics
type DocumentAnalytics struct {
	AccessCount         int        `gorm:"default:0" json:"access_count"`
	UniqueViewers       int        `gorm:"default:0" json:"unique_viewers"`
	AverageViewTime     int        `json:"average_view_time"` // seconds
	SearchAppearances   int        `gorm:"default:0" json:"search_appearances"`
	CitationCount       int        `gorm:"default:0" json:"citation_count"`
	RelatedDocuments    string     `json:"related_documents,omitempty"` // JSON array of related doc IDs
	PopularityScore     float64    `json:"popularity_score"`
	LastAnalyticsUpdate *time.Time `json:"last_analytics_update,omitempty"`
}

// SignatureMetadata contains electronic signature metadata
type SignatureMetadata struct {
	SignerName    string `gorm:"not null" json:"signer_name"`
	SignerEmail   string `gorm:"not null" json:"signer_email"`
	SignerRole    string `json:"signer_role"`                    // policyholder, agent, witness, notary
	SignatureType string `gorm:"not null" json:"signature_type"` // drawn, typed, uploaded, digital_certificate
	SignatureData string `json:"signature_data,omitempty"`       // Base64 encoded signature
	CertificateID string `json:"certificate_id,omitempty"`       // Digital certificate ID
	Provider      string `json:"provider,omitempty"`             // docusign, hellosign, adobe_sign, internal
	ProviderTxnID string `json:"provider_txn_id,omitempty"`
}

// SignatureVerification contains signature verification details
type SignatureVerification struct {
	VerificationCode   string `json:"-"`                             // For email/SMS verification
	VerificationMethod string `json:"verification_method,omitempty"` // email, sms, biometric, password
	IPAddress          string `json:"ip_address"`
	UserAgent          string `json:"user_agent"`
	Location           string `json:"location,omitempty"` // Geolocation
	BiometricData      string `json:"-"`                  // Encrypted biometric data
}

// SignatureLifecycle contains signature status and timestamps
type SignatureLifecycle struct {
	Status        string     `gorm:"not null;default:'pending'" json:"status"` // pending, signed, declined, expired, revoked
	SignedAt      *time.Time `json:"signed_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	DeclinedAt    *time.Time `json:"declined_at,omitempty"`
	DeclineReason string     `json:"decline_reason,omitempty"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	RevokeReason  string     `json:"revoke_reason,omitempty"`
}

// TemplateConfiguration contains template configuration
type TemplateConfiguration struct {
	Code        string `gorm:"uniqueIndex;not null" json:"code"` // template identifier
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Category    string `gorm:"not null" json:"category"` // policy, claim, financial, legal
	Version     string `gorm:"not null" json:"version"`
	Language    string `gorm:"default:'en'" json:"language"`
	ContentType string `gorm:"not null" json:"content_type"` // html, markdown, pdf
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// TemplateContent contains template content and requirements
type TemplateContent struct {
	Template           string     `gorm:"type:text" json:"template"` // Template content with placeholders
	Schema             string     `json:"schema,omitempty"`          // JSON schema for template variables
	RequiredFields     string     `json:"required_fields,omitempty"` // JSON array of required fields
	SignatureRequired  bool       `gorm:"default:false" json:"signature_required"`
	SignaturePositions string     `json:"signature_positions,omitempty"` // JSON array of signature positions
	UsageCount         int        `gorm:"default:0" json:"usage_count"`
	LastUsedAt         *time.Time `json:"last_used_at,omitempty"`
}
