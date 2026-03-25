package services

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// PartnerService defines the interface for partner business logic operations
type PartnerService interface {
	// Partner management
	CreatePartner(ctx context.Context, partner *models.Partner) error
	GetPartnerByID(ctx context.Context, id uuid.UUID) (*models.Partner, error)
	GetPartnerByCode(ctx context.Context, partnerCode string) (*models.Partner, error)
	UpdatePartner(ctx context.Context, partner *models.Partner) error
	DeactivatePartner(ctx context.Context, id uuid.UUID, reason string) error

	// Partner search and listing
	ListPartners(ctx context.Context, partnerType, status string, limit, offset int) ([]*models.Partner, int64, error)
	SearchPartners(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*models.Partner, int64, error)
	GetAvailablePartners(ctx context.Context, partnerType string, limit, offset int) ([]*models.Partner, int64, error)

	// Partner assignment management
	CreateAssignment(ctx context.Context, assignment *models.PartnerAssignment) error
	GetAssignmentByID(ctx context.Context, id uuid.UUID) (*models.PartnerAssignment, error)
	GetAssignmentsByPartnerID(ctx context.Context, partnerID uuid.UUID, limit, offset int) ([]*models.PartnerAssignment, int64, error)
	UpdateAssignment(ctx context.Context, assignment *models.PartnerAssignment) error
	CompleteAssignment(ctx context.Context, assignmentID uuid.UUID, completionNotes string) error

	// Partner performance
	UpdatePartnerRating(ctx context.Context, partnerID uuid.UUID, rating float64, reviewCount int) error
	GetPartnerPerformance(ctx context.Context, partnerID uuid.UUID) (map[string]interface{}, error)
	GetTopRatedPartners(ctx context.Context, partnerType string, limit int) ([]*models.Partner, error)

	// Partner availability
	UpdatePartnerAvailability(ctx context.Context, partnerID uuid.UUID, isAvailable bool) error
	CheckPartnerAvailability(ctx context.Context, partnerID uuid.UUID, date string) (bool, error)
}
