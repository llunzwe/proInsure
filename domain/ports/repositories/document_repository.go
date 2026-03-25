package repositories

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/document"
)

// DocumentRepository defines the interface for document persistence operations
type DocumentRepository interface {
	// === Core CRUD Operations ===
	Create(ctx context.Context, doc *models.Document) error
	Update(ctx context.Context, doc *models.Document) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error)
	GetByDocumentNumber(ctx context.Context, number string) (*models.Document, error)
	List(ctx context.Context, criteria DocumentSearchCriteria) ([]*models.Document, int64, error)

	// === Batch Operations ===
	BulkCreate(ctx context.Context, docs []*models.Document) error
	BulkUpdate(ctx context.Context, docs []*models.Document) error
	BulkDelete(ctx context.Context, ids []uuid.UUID) error

	// === Search & Filter ===
	Search(ctx context.Context, query string, filters map[string]interface{}) ([]*models.Document, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.Document, error)
	GetByPolicy(ctx context.Context, policyID uuid.UUID) ([]*models.Document, error)
	GetByClaim(ctx context.Context, claimID uuid.UUID) ([]*models.Document, error)
	GetByType(ctx context.Context, docType string) ([]*models.Document, error)
	GetByCategory(ctx context.Context, category string) ([]*models.Document, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Document, error)
	GetExpiring(ctx context.Context, days int) ([]*models.Document, error)

	// === Document Lifecycle ===
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, userID uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	Expire(ctx context.Context, id uuid.UUID) error
	SetRetention(ctx context.Context, id uuid.UUID, retentionPeriod int) error
	SetLegalHold(ctx context.Context, id uuid.UUID, hold bool, reason string) error

	// === Version Management ===
	CreateVersion(ctx context.Context, version *models.Document) error
	GetVersions(ctx context.Context, documentID uuid.UUID) ([]*models.Document, error)
	GetVersion(ctx context.Context, documentID uuid.UUID, versionNumber string) (*models.Document, error)
	RestoreVersion(ctx context.Context, documentID uuid.UUID, versionNumber string) error
	CompareVersions(ctx context.Context, documentID uuid.UUID, v1, v2 string) (*VersionComparison, error)

	// === Signatures ===
	AddSignature(ctx context.Context, documentID uuid.UUID, signature *models.ESignature) error
	GetSignatures(ctx context.Context, documentID uuid.UUID) ([]*models.ESignature, error)
	VerifySignature(ctx context.Context, documentID uuid.UUID, signatureID uuid.UUID) (bool, error)
	RequestSignature(ctx context.Context, documentID uuid.UUID, signerID uuid.UUID) error

	// === Sharing & Access ===
	CreateShare(ctx context.Context, share *models.DocumentAccess) error
	GetShares(ctx context.Context, documentID uuid.UUID) ([]*models.DocumentAccess, error)
	RevokeShare(ctx context.Context, shareID uuid.UUID) error
	UpdatePermissions(ctx context.Context, documentID uuid.UUID, userID uuid.UUID, permissions []string) error
	GetAccessLog(ctx context.Context, documentID uuid.UUID) ([]*models.DocumentAccess, error)

	// === Workflow ===
	CreateWorkflow(ctx context.Context, workflow *document.DocumentWorkflow) error
	GetWorkflow(ctx context.Context, documentID uuid.UUID) (*document.DocumentWorkflow, error)
	UpdateWorkflowStage(ctx context.Context, workflowID uuid.UUID, stage string, userID uuid.UUID) error
	ApproveDocument(ctx context.Context, documentID uuid.UUID, approverID uuid.UUID, comments string) error
	RejectDocument(ctx context.Context, documentID uuid.UUID, approverID uuid.UUID, reason string) error

	// === OCR & Processing ===
	CreateOCRRecord(ctx context.Context, documentID uuid.UUID, ocrData map[string]interface{}) error
	GetOCRRecord(ctx context.Context, documentID uuid.UUID) (map[string]interface{}, error)
	UpdateOCRStatus(ctx context.Context, documentID uuid.UUID, status string, result interface{}) error
	GetDocumentsForOCR(ctx context.Context) ([]*models.Document, error)

	// === Classification ===
	CreateClassification(ctx context.Context, documentID uuid.UUID, docType string, category string) error
	GetClassification(ctx context.Context, documentID uuid.UUID) (map[string]interface{}, error)
	UpdateTags(ctx context.Context, documentID uuid.UUID, tags []string) error
	GetByTags(ctx context.Context, tags []string) ([]*models.Document, error)

	// === Compliance ===
	CreateComplianceRecord(ctx context.Context, documentID uuid.UUID, standard string, status string, notes string) error
	GetComplianceRecords(ctx context.Context, documentID uuid.UUID) (map[string]interface{}, error)
	CheckCompliance(ctx context.Context, documentID uuid.UUID, standard string) (bool, error)
	GetNonCompliantDocuments(ctx context.Context) ([]*models.Document, error)

	// === Analytics ===
	CreateAnalytics(ctx context.Context, analytics *document.DocumentAnalytics) error
	GetAnalytics(ctx context.Context, documentID uuid.UUID) (*document.DocumentAnalytics, error)
	IncrementViewCount(ctx context.Context, documentID uuid.UUID) error
	IncrementDownloadCount(ctx context.Context, documentID uuid.UUID) error
	GetPopularDocuments(ctx context.Context, limit int) ([]*models.Document, error)

	// === Retention & Disposal ===
	GetDocumentsForRetention(ctx context.Context) ([]*models.Document, error)
	GetDocumentsForDisposal(ctx context.Context) ([]*models.Document, error)
	MarkForDisposal(ctx context.Context, documentID uuid.UUID, reason string) error
	ExecuteDisposal(ctx context.Context, documentID uuid.UUID) error

	// === Audit ===
	CreateAuditLog(ctx context.Context, documentID uuid.UUID, action string, userID uuid.UUID, details map[string]interface{}) error
	GetAuditLogs(ctx context.Context, documentID uuid.UUID, from, to time.Time) ([]map[string]interface{}, error)
	GetUserActivity(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]map[string]interface{}, error)

	// === Integration ===
	CreateIntegration(ctx context.Context, documentID uuid.UUID, provider string, externalID string) error
	GetIntegrations(ctx context.Context, documentID uuid.UUID) ([]map[string]interface{}, error)
	SyncWithExternalSystem(ctx context.Context, documentID uuid.UUID, systemID string) error

	// === Statistics ===
	GetStatistics(ctx context.Context) (*DocumentStatistics, error)
	GetStorageUsage(ctx context.Context, userID *uuid.UUID) (*StorageUsage, error)
	GetTypeDistribution(ctx context.Context) (map[string]int64, error)
}

// DocumentSearchCriteria defines search parameters for documents
type DocumentSearchCriteria struct {
	Query          string
	UserID         *uuid.UUID
	PolicyID       *uuid.UUID
	ClaimID        *uuid.UUID
	Type           string
	Category       string
	Status         []string
	SecurityLevel  string
	Tags           []string
	DateFrom       *time.Time
	DateTo         *time.Time
	HasSignatures  *bool
	IsVerified     *bool
	IsExpired      *bool
	RequiresAction *bool
	SortBy         string
	SortOrder      string
	Limit          int
	Offset         int
}

// DocumentStatistics represents document system statistics
type DocumentStatistics struct {
	TotalDocuments    int64
	ActiveDocuments   int64
	ArchivedDocuments int64
	ExpiredDocuments  int64
	TotalSize         int64
	AverageSize       int64
	DocumentsByType   map[string]int64
	DocumentsByStatus map[string]int64
	SignedDocuments   int64
	VerifiedDocuments int64
	SharedDocuments   int64
	ComplianceRate    float64
}

// StorageUsage represents storage usage information
type StorageUsage struct {
	UserID        uuid.UUID
	TotalSize     int64
	DocumentCount int64
	QuotaLimit    int64
	QuotaUsed     float64
	LargestFile   string
	OldestFile    string
	LastUpload    *time.Time
}

// VersionComparison represents comparison between two document versions
type VersionComparison struct {
	Version1       string
	Version2       string
	Changes        []VersionChange
	AddedFields    []string
	RemovedFields  []string
	ModifiedFields []string
	ChangeCount    int
}

// VersionChange represents a single change between versions
type VersionChange struct {
	Field    string
	OldValue interface{}
	NewValue interface{}
	Type     string // Added, Removed, Modified
}
