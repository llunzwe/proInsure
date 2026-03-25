package repositories

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
)

// FamilyRepository defines the interface for family plan persistence operations
type FamilyRepository interface {
	// Family plan operations
	CreateFamilyPlan(ctx context.Context, plan *models.FamilyPlan) error
	GetFamilyPlanByID(ctx context.Context, id uuid.UUID) (*models.FamilyPlan, error)
	GetFamilyPlansByPrimaryUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.FamilyPlan, int64, error)
	UpdateFamilyPlan(ctx context.Context, plan *models.FamilyPlan) error
	DeleteFamilyPlan(ctx context.Context, id uuid.UUID) error
	ListFamilyPlans(ctx context.Context, status string, limit, offset int) ([]*models.FamilyPlan, int64, error)

	// Family member operations
	CreateFamilyMember(ctx context.Context, member *models.FamilyMember) error
	GetFamilyMemberByID(ctx context.Context, id uuid.UUID) (*models.FamilyMember, error)
	GetFamilyMembersByFamilyPlanID(ctx context.Context, familyPlanID uuid.UUID) ([]*models.FamilyMember, error)
	GetFamilyMemberByUserID(ctx context.Context, userID uuid.UUID) (*models.FamilyMember, error)
	GetActiveFamilyMembersByFamilyPlanID(ctx context.Context, familyPlanID uuid.UUID) ([]*models.FamilyMember, error)
	UpdateFamilyMember(ctx context.Context, member *models.FamilyMember) error
	DeleteFamilyMember(ctx context.Context, id uuid.UUID) error

	// Statistics
	GetFamilyPlanStatistics(ctx context.Context, familyPlanID uuid.UUID) (map[string]interface{}, error)
	GetFamilyMemberCount(ctx context.Context, familyPlanID uuid.UUID) (int64, error)
	GetActiveMemberCount(ctx context.Context, familyPlanID uuid.UUID) (int64, error)
}
