package services

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// FamilyService defines the interface for family plan business logic operations
type FamilyService interface {
	// Family plan management
	CreateFamilyPlan(ctx context.Context, plan *models.FamilyPlan) error
	GetFamilyPlanByID(ctx context.Context, id uuid.UUID) (*models.FamilyPlan, error)
	GetFamilyPlansByPrimaryUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.FamilyPlan, int64, error)
	UpdateFamilyPlan(ctx context.Context, plan *models.FamilyPlan) error
	CancelFamilyPlan(ctx context.Context, planID uuid.UUID, reason string) error
	ListFamilyPlans(ctx context.Context, status string, limit, offset int) ([]*models.FamilyPlan, int64, error)

	// Family member management
	AddFamilyMember(ctx context.Context, member *models.FamilyMember) error
	GetFamilyMemberByID(ctx context.Context, id uuid.UUID) (*models.FamilyMember, error)
	GetFamilyMembersByFamilyPlanID(ctx context.Context, familyPlanID uuid.UUID) ([]*models.FamilyMember, error)
	GetActiveFamilyMembersByFamilyPlanID(ctx context.Context, familyPlanID uuid.UUID) ([]*models.FamilyMember, error)
	UpdateFamilyMember(ctx context.Context, member *models.FamilyMember) error
	RemoveFamilyMember(ctx context.Context, memberID uuid.UUID) error

	// Family plan operations
	CalculateFamilyDiscount(ctx context.Context, planID uuid.UUID, basePremium float64) (float64, error)
	UpdateFamilyPlanPremium(ctx context.Context, planID uuid.UUID) error
	CheckMemberLimit(ctx context.Context, planID uuid.UUID) (bool, error)
	CheckDeviceLimit(ctx context.Context, planID uuid.UUID, memberID uuid.UUID) (bool, error)

	// Family plan statistics
	GetFamilyPlanStatistics(ctx context.Context, familyPlanID uuid.UUID) (map[string]interface{}, error)
	GetFamilyMemberCount(ctx context.Context, familyPlanID uuid.UUID) (int64, error)
	GetActiveMemberCount(ctx context.Context, familyPlanID uuid.UUID) (int64, error)
}
