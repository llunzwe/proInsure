package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/document"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/pkg/errors"
)

// documentRepositoryImpl implements DocumentRepository interface
type documentRepositoryImpl struct {
	*BaseRepository
	db     *gorm.DB
	logger Logger
}

// NewDocumentRepository creates a new document repository instance
func NewDocumentRepository(db *gorm.DB, logger Logger) repositories.DocumentRepository {
	return &documentRepositoryImpl{
		BaseRepository: NewBaseRepository(db, logger),
		db:             db,
		logger:         logger,
	}
}

// === Core CRUD Operations ===

func (r *documentRepositoryImpl) Create(ctx context.Context, doc *models.Document) error {
	if err := r.db.WithContext(ctx).Create(doc).Error; err != nil {
		return r.HandleError(err, "create document")
	}
	r.logger.Info("Document created", "id", doc.ID, "number", doc.IDDocumentNumber)
	return nil
}

func (r *documentRepositoryImpl) Update(ctx context.Context, doc *models.Document) error {
	result := r.db.WithContext(ctx).Model(doc).Updates(doc)
	if result.Error != nil {
		return r.HandleError(result.Error, "update document")
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("document not found")
	}
	return nil
}

func (r *documentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Document{}, "id = ?", id)
	if result.Error != nil {
		return r.HandleError(result.Error, "delete document")
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("document not found")
	}
	return nil
}

func (r *documentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error) {
	var doc models.Document
	err := r.db.WithContext(ctx).
		Preload("DocumentOCR").
		Preload("DocumentWorkflow").
		Preload("DocumentVersions").
		Preload("DocumentShares").
		Preload("DocumentAnalytics").
		Preload("DocumentRetention").
		Preload("DocumentClassification").
		First(&doc, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("document not found")
		}
		return nil, r.HandleError(err, "get document by ID")
	}
	return &doc, nil
}

func (r *documentRepositoryImpl) GetByDocumentNumber(ctx context.Context, number string) (*models.Document, error) {
	var doc models.Document
	err := r.db.WithContext(ctx).
		Where("id_document_number = ?", number).
		First(&doc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("document not found")
		}
		return nil, r.HandleError(err, "get document by number")
	}
	return &doc, nil
}

func (r *documentRepositoryImpl) List(ctx context.Context, criteria repositories.DocumentSearchCriteria) ([]*models.Document, int64, error) {
	var documents []*models.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Document{})

	// Apply filters
	if criteria.UserID != nil {
		query = query.Where("rel_user_id = ?", *criteria.UserID)
	}
	if criteria.PolicyID != nil {
		query = query.Where("rel_policy_id = ?", *criteria.PolicyID)
	}
	if criteria.ClaimID != nil {
		query = query.Where("rel_claim_id = ?", *criteria.ClaimID)
	}
	if criteria.Type != "" {
		query = query.Where("id_type = ?", criteria.Type)
	}
	if criteria.Category != "" {
		query = query.Where("id_category = ?", criteria.Category)
	}
	if len(criteria.Status) > 0 {
		query = query.Where("lc_status IN ?", criteria.Status)
	}
	if criteria.SecurityLevel != "" {
		query = query.Where("sec_security_level = ?", criteria.SecurityLevel)
	}
	if criteria.DateFrom != nil {
		query = query.Where("created_at >= ?", *criteria.DateFrom)
	}
	if criteria.DateTo != nil {
		query = query.Where("created_at <= ?", *criteria.DateTo)
	}
	if criteria.IsExpired != nil && *criteria.IsExpired {
		query = query.Where("lc_expires_at < ?", time.Now())
	}
	if criteria.IsVerified != nil {
		query = query.Where("ver_is_verified = ?", *criteria.IsVerified)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, r.HandleError(err, "count documents")
	}

	// Apply sorting
	if criteria.SortBy != "" {
		order := "ASC"
		if criteria.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", criteria.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if criteria.Limit > 0 {
		query = query.Limit(criteria.Limit)
	}
	if criteria.Offset > 0 {
		query = query.Offset(criteria.Offset)
	}

	// Execute query
	if err := query.Find(&documents).Error; err != nil {
		return nil, 0, r.HandleError(err, "list documents")
	}

	return documents, total, nil
}

// === Batch Operations ===

func (r *documentRepositoryImpl) BulkCreate(ctx context.Context, docs []*models.Document) error {
	if len(docs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, doc := range docs {
			if err := tx.Create(doc).Error; err != nil {
				return r.HandleError(err, "bulk create document")
			}
		}
		return nil
	})
}

func (r *documentRepositoryImpl) BulkUpdate(ctx context.Context, docs []*models.Document) error {
	if len(docs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, doc := range docs {
			if err := tx.Model(doc).Updates(doc).Error; err != nil {
				return r.HandleError(err, "bulk update document")
			}
		}
		return nil
	})
}

func (r *documentRepositoryImpl) BulkDelete(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Delete(&models.Document{}, "id IN ?", ids)
	if result.Error != nil {
		return r.HandleError(result.Error, "bulk delete documents")
	}
	return nil
}

// === Search & Filter ===

func (r *documentRepositoryImpl) Search(ctx context.Context, query string, filters map[string]interface{}) ([]*models.Document, error) {
	var documents []*models.Document

	dbQuery := r.db.WithContext(ctx).Model(&models.Document{})

	// Text search
	if query != "" {
		dbQuery = dbQuery.Where(
			"id_document_number ILIKE ? OR id_title ILIKE ? OR md_description ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%",
		)
	}

	// Apply filters
	for key, value := range filters {
		dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if err := dbQuery.Find(&documents).Error; err != nil {
		return nil, r.HandleError(err, "search documents")
	}

	return documents, nil
}

func (r *documentRepositoryImpl) GetByUser(ctx context.Context, userID uuid.UUID) ([]*models.Document, error) {
	var documents []*models.Document
	err := r.db.WithContext(ctx).
		Where("rel_user_id = ?", userID).
		Find(&documents).Error

	if err != nil {
		return nil, r.HandleError(err, "get documents by user")
	}
	return documents, nil
}

func (r *documentRepositoryImpl) GetByPolicy(ctx context.Context, policyID uuid.UUID) ([]*models.Document, error) {
	var documents []*models.Document
	err := r.db.WithContext(ctx).
		Where("rel_policy_id = ?", policyID).
		Find(&documents).Error

	if err != nil {
		return nil, r.HandleError(err, "get documents by policy")
	}
	return documents, nil
}

func (r *documentRepositoryImpl) GetByClaim(ctx context.Context, claimID uuid.UUID) ([]*models.Document, error) {
	var documents []*models.Document
	err := r.db.WithContext(ctx).
		Where("rel_claim_id = ?", claimID).
		Find(&documents).Error

	if err != nil {
		return nil, r.HandleError(err, "get documents by claim")
	}
	return documents, nil
}

func (r *documentRepositoryImpl) GetByType(ctx context.Context, docType string) ([]*models.Document, error) {
	var documents []*models.Document
	err := r.db.WithContext(ctx).
		Where("id_type = ?", docType).
		Find(&documents).Error

	if err != nil {
		return nil, r.HandleError(err, "get documents by type")
	}
	return documents, nil
}

func (r *documentRepositoryImpl) GetByCategory(ctx context.Context, category string) ([]*models.Document, error) {
	var documents []*models.Document
	err := r.db.WithContext(ctx).
		Where("id_category = ?", category).
		Find(&documents).Error

	if err != nil {
		return nil, r.HandleError(err, "get documents by category")
	}
	return documents, nil
}

func (r *documentRepositoryImpl) GetByStatus(ctx context.Context, status string) ([]*models.Document, error) {
	var documents []*models.Document
	err := r.db.WithContext(ctx).
		Where("lc_status = ?", status).
		Find(&documents).Error

	if err != nil {
		return nil, r.HandleError(err, "get documents by status")
	}
	return documents, nil
}

func (r *documentRepositoryImpl) GetExpiring(ctx context.Context, days int) ([]*models.Document, error) {
	var documents []*models.Document
	expiryDate := time.Now().AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Where("lc_expires_at IS NOT NULL AND lc_expires_at <= ? AND lc_status = ?",
			expiryDate, document.StatusActive).
		Find(&documents).Error

	if err != nil {
		return nil, r.HandleError(err, "get expiring documents")
	}
	return documents, nil
}
