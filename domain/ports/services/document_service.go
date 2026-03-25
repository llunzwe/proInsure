package services

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/document"
	"smartsure/internal/domain/ports/repositories"
)

// DocumentService defines the interface for document business operations
type DocumentService interface {
	// === Document Management ===
	CreateDocument(ctx context.Context, doc *models.Document, content io.Reader, userID uuid.UUID) (*models.Document, error)
	UpdateDocument(ctx context.Context, id uuid.UUID, updates map[string]interface{}, userID uuid.UUID) (*models.Document, error)
	DeleteDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Document, error)
	ListDocuments(ctx context.Context, criteria repositories.DocumentSearchCriteria) ([]*models.Document, int64, error)
	SearchDocuments(ctx context.Context, query string, userID uuid.UUID) ([]*models.Document, error)

	// === Content Management ===
	UploadContent(ctx context.Context, documentID uuid.UUID, content io.Reader, userID uuid.UUID) error
	DownloadContent(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) (io.ReadCloser, error)
	GetContentURL(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) (string, error)
	GenerateThumbnail(ctx context.Context, documentID uuid.UUID) error

	// === Version Control ===
	CreateVersion(ctx context.Context, documentID uuid.UUID, comment string, userID uuid.UUID) (interface{}, error) // DocumentVersion type not found
	ListVersions(ctx context.Context, documentID uuid.UUID) ([]interface{}, error)                                  // DocumentVersion type not found
	GetVersion(ctx context.Context, documentID uuid.UUID, versionNumber string) (interface{}, error)                // DocumentVersion type not found
	RestoreVersion(ctx context.Context, documentID uuid.UUID, versionNumber string, userID uuid.UUID) error
	CompareVersions(ctx context.Context, documentID uuid.UUID, v1, v2 string) (*repositories.VersionComparison, error)

	// === Lifecycle Management ===
	UpdateStatus(ctx context.Context, documentID uuid.UUID, status string, userID uuid.UUID) error
	ArchiveDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error
	RestoreDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error
	ExpireDocument(ctx context.Context, documentID uuid.UUID) error
	SetRetentionPolicy(ctx context.Context, documentID uuid.UUID, retentionPeriod int, userID uuid.UUID) error
	SetLegalHold(ctx context.Context, documentID uuid.UUID, hold bool, reason string, userID uuid.UUID) error

	// === Signatures & Verification ===
	RequestSignature(ctx context.Context, documentID uuid.UUID, signerID uuid.UUID, requesterID uuid.UUID) error
	SignDocument(ctx context.Context, documentID uuid.UUID, signature []byte, userID uuid.UUID) (*models.ESignature, error)
	VerifySignature(ctx context.Context, documentID uuid.UUID, signatureID uuid.UUID) (bool, error)
	GetSignatures(ctx context.Context, documentID uuid.UUID) ([]*models.ESignature, error)
	VerifyDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) (bool, error)

	// === Sharing & Collaboration ===
	ShareDocument(ctx context.Context, documentID uuid.UUID, targetUserID uuid.UUID, permissions []string, expiresAt *time.Time, userID uuid.UUID) (interface{}, error) // DocumentShare type not found
	RevokeShare(ctx context.Context, shareID uuid.UUID, userID uuid.UUID) error
	GetSharedDocuments(ctx context.Context, userID uuid.UUID) ([]*models.Document, error)
	GetDocumentShares(ctx context.Context, documentID uuid.UUID) ([]*models.DocumentAccess, error)
	UpdateSharePermissions(ctx context.Context, shareID uuid.UUID, permissions []string, userID uuid.UUID) error

	// === Workflow Management ===
	InitiateWorkflow(ctx context.Context, documentID uuid.UUID, workflowType string, userID uuid.UUID) error
	AdvanceWorkflow(ctx context.Context, workflowID uuid.UUID, action string, comments string, userID uuid.UUID) error
	GetWorkflowStatus(ctx context.Context, documentID uuid.UUID) (*document.DocumentWorkflow, error)
	ApproveDocument(ctx context.Context, documentID uuid.UUID, comments string, userID uuid.UUID) error
	RejectDocument(ctx context.Context, documentID uuid.UUID, reason string, userID uuid.UUID) error
	GetPendingApprovals(ctx context.Context, userID uuid.UUID) ([]*models.Document, error)

	// === OCR & Text Processing ===
	ProcessOCR(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error
	GetOCRResult(ctx context.Context, documentID uuid.UUID) (interface{}, error) // DocumentOCR type not found
	ExtractText(ctx context.Context, documentID uuid.UUID) (string, error)
	ExtractMetadata(ctx context.Context, documentID uuid.UUID) (map[string]interface{}, error)
	SearchInContent(ctx context.Context, documentID uuid.UUID, searchText string) ([]string, error)

	// === Classification & Tagging ===
	ClassifyDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) (interface{}, error) // DocumentClassification type not found
	UpdateTags(ctx context.Context, documentID uuid.UUID, tags []string, userID uuid.UUID) error
	GetDocumentsByTag(ctx context.Context, tags []string) ([]*models.Document, error)
	SuggestTags(ctx context.Context, documentID uuid.UUID) ([]string, error)
	AutoCategorize(ctx context.Context, documentID uuid.UUID) error

	// === Compliance & Audit ===
	CheckCompliance(ctx context.Context, documentID uuid.UUID, standards []string) (map[string]interface{}, error)
	GetComplianceReport(ctx context.Context, documentID uuid.UUID) (map[string]interface{}, error)
	MarkAsCompliant(ctx context.Context, documentID uuid.UUID, standard string, userID uuid.UUID) error
	GetAuditTrail(ctx context.Context, documentID uuid.UUID, from, to time.Time) ([]interface{}, error) // DocumentAuditLog type not found
	LogActivity(ctx context.Context, documentID uuid.UUID, action string, details map[string]interface{}, userID uuid.UUID) error

	// === Analytics & Reporting ===
	GetDocumentAnalytics(ctx context.Context, documentID uuid.UUID) (*document.DocumentAnalytics, error)
	GetUserDocumentStats(ctx context.Context, userID uuid.UUID) (*repositories.DocumentStatistics, error)
	GetSystemStats(ctx context.Context) (*repositories.DocumentStatistics, error)
	GenerateUsageReport(ctx context.Context, from, to time.Time) (map[string]interface{}, error)
	GetPopularDocuments(ctx context.Context, limit int) ([]*models.Document, error)

	// === Templates ===
	CreateTemplate(ctx context.Context, template *document.DocumentTemplate, userID uuid.UUID) (*document.DocumentTemplate, error)
	GetTemplate(ctx context.Context, templateID uuid.UUID) (*document.DocumentTemplate, error)
	ListTemplates(ctx context.Context, category string) ([]*document.DocumentTemplate, error)
	GenerateFromTemplate(ctx context.Context, templateID uuid.UUID, data map[string]interface{}, userID uuid.UUID) (*models.Document, error)

	// === Bulk Operations ===
	BulkUpload(ctx context.Context, documents []DocumentUpload, userID uuid.UUID) ([]*models.Document, []error)
	BulkDelete(ctx context.Context, documentIDs []uuid.UUID, userID uuid.UUID) (int, []error)
	BulkArchive(ctx context.Context, documentIDs []uuid.UUID, userID uuid.UUID) (int, []error)
	BulkUpdateStatus(ctx context.Context, documentIDs []uuid.UUID, status string, userID uuid.UUID) (int, []error)
	BulkTag(ctx context.Context, documentIDs []uuid.UUID, tags []string, userID uuid.UUID) (int, []error)

	// === Integration ===
	ExportDocuments(ctx context.Context, criteria repositories.DocumentSearchCriteria, format string, userID uuid.UUID) (io.ReadCloser, error)
	ImportDocuments(ctx context.Context, source io.Reader, format string, userID uuid.UUID) ([]*models.Document, error)
	SyncWithExternalSystem(ctx context.Context, systemID string, userID uuid.UUID) error
	GetIntegrationStatus(ctx context.Context, documentID uuid.UUID, systemID string) (map[string]interface{}, error)

	// === Retention & Disposal ===
	ApplyRetentionPolicies(ctx context.Context) (int, error)
	GetDocumentsForDisposal(ctx context.Context) ([]*models.Document, error)
	DisposeDocuments(ctx context.Context, documentIDs []uuid.UUID, reason string, userID uuid.UUID) error
	GetRetentionSchedule(ctx context.Context) ([]map[string]interface{}, error)

	// === Security ===
	EncryptDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error
	DecryptDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error
	SetAccessControl(ctx context.Context, documentID uuid.UUID, acl map[string][]string, userID uuid.UUID) error
	CheckAccess(ctx context.Context, documentID uuid.UUID, userID uuid.UUID, permission string) bool
	GenerateSecureLink(ctx context.Context, documentID uuid.UUID, expiresIn time.Duration, userID uuid.UUID) (string, error)
}

// DocumentUpload represents a document upload request
type DocumentUpload struct {
	FileName    string
	Content     io.Reader
	ContentType string
	Size        int64
	Metadata    map[string]interface{}
}

// DocumentTemplate represents a document template
type DocumentTemplate struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Content     string                 `json:"content"`
	Variables   []TemplateVariable     `json:"variables"`
	Metadata    map[string]interface{} `json:"metadata"`
	IsActive    bool                   `json:"is_active"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// TemplateVariable represents a variable in a document template
type TemplateVariable struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"default_value"`
	Validation   string      `json:"validation"`
}
