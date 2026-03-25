package repositories

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// PartnerRepository defines the interface for partner persistence operations
type PartnerRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, partner *models.Partner) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Partner, error)
	GetByPartnerCode(ctx context.Context, partnerCode string) (*models.Partner, error)
	Update(ctx context.Context, partner *models.Partner) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByType(ctx context.Context, partnerType string, limit, offset int) ([]*models.Partner, int64, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Partner, int64, error)
	GetByEmail(ctx context.Context, email string) (*models.Partner, error)
	GetByServiceArea(ctx context.Context, area string, limit, offset int) ([]*models.Partner, int64, error)
	GetAvailablePartners(ctx context.Context, partnerType string, limit, offset int) ([]*models.Partner, int64, error)
	Search(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*models.Partner, int64, error)

	// Assignment operations
	GetAssignmentsByPartner(ctx context.Context, partnerID uuid.UUID, limit, offset int) ([]*models.PartnerAssignment, int64, error)
	GetActiveAssignments(ctx context.Context, partnerID uuid.UUID) ([]*models.PartnerAssignment, error)
	CreateAssignment(ctx context.Context, assignment *models.PartnerAssignment) error
	UpdateAssignment(ctx context.Context, assignment *models.PartnerAssignment) error

	// Performance metrics
	UpdatePerformanceMetrics(ctx context.Context, partnerID uuid.UUID, metrics map[string]interface{}) error
	GetTopRatedPartners(ctx context.Context, partnerType string, limit int) ([]*models.Partner, error)

	// Availability management
	UpdateAvailability(ctx context.Context, partnerID uuid.UUID, isAvailable bool) error
	UpdateCurrentAssignments(ctx context.Context, partnerID uuid.UUID, count int) error
}
