package family

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// familyService implements the FamilyService interface
type familyService struct {
	familyRepo repositories.FamilyRepository
}

// NewFamilyService creates a new family service
func NewFamilyService(
	familyRepo repositories.FamilyRepository,
) services.FamilyService {
	return &familyService{
		familyRepo: familyRepo,
	}
}

// CreateFamilyPlan creates a new family plan
func (s *familyService) CreateFamilyPlan(ctx context.Context, plan *models.FamilyPlan) error {
	if plan == nil {
		return errors.New("family plan cannot be nil")
	}
	if plan.ID == uuid.Nil {
		plan.ID = uuid.New()
	}
	return s.familyRepo.CreateFamilyPlan(ctx, plan)
}

// GetFamilyPlanByID retrieves a family plan by ID
func (s *familyService) GetFamilyPlanByID(ctx context.Context, id uuid.UUID) (*models.FamilyPlan, error) {
	return s.familyRepo.GetFamilyPlanByID(ctx, id)
}

// GetFamilyPlansByPrimaryUserID gets family plans for a primary user
func (s *familyService) GetFamilyPlansByPrimaryUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.FamilyPlan, int64, error) {
	return s.familyRepo.GetFamilyPlansByPrimaryUserID(ctx, userID, limit, offset)
}

// UpdateFamilyPlan updates an existing family plan
func (s *familyService) UpdateFamilyPlan(ctx context.Context, plan *models.FamilyPlan) error {
	if plan == nil {
		return errors.New("family plan cannot be nil")
	}
	return s.familyRepo.UpdateFamilyPlan(ctx, plan)
}

// CancelFamilyPlan cancels a family plan
func (s *familyService) CancelFamilyPlan(ctx context.Context, planID uuid.UUID, reason string) error {
	plan, err := s.familyRepo.GetFamilyPlanByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("failed to get family plan: %w", err)
	}
	if plan == nil {
		return errors.New("family plan not found")
	}
	plan.Status = "cancelled"
	return s.familyRepo.UpdateFamilyPlan(ctx, plan)
}

// ListFamilyPlans lists family plans with filters
func (s *familyService) ListFamilyPlans(ctx context.Context, status string, limit, offset int) ([]*models.FamilyPlan, int64, error) {
	return s.familyRepo.ListFamilyPlans(ctx, status, limit, offset)
}

// AddFamilyMember adds a member to a family plan
func (s *familyService) AddFamilyMember(ctx context.Context, member *models.FamilyMember) error {
	if member == nil {
		return errors.New("family member cannot be nil")
	}
	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	return s.familyRepo.CreateFamilyMember(ctx, member)
}

// GetFamilyMemberByID retrieves a family member by ID
func (s *familyService) GetFamilyMemberByID(ctx context.Context, id uuid.UUID) (*models.FamilyMember, error) {
	return s.familyRepo.GetFamilyMemberByID(ctx, id)
}

// GetFamilyMembersByFamilyPlanID gets all members for a family plan
func (s *familyService) GetFamilyMembersByFamilyPlanID(ctx context.Context, familyPlanID uuid.UUID) ([]*models.FamilyMember, error) {
	return s.familyRepo.GetFamilyMembersByFamilyPlanID(ctx, familyPlanID)
}

// GetActiveFamilyMembersByFamilyPlanID gets active members for a family plan
func (s *familyService) GetActiveFamilyMembersByFamilyPlanID(ctx context.Context, familyPlanID uuid.UUID) ([]*models.FamilyMember, error) {
	return s.familyRepo.GetActiveFamilyMembersByFamilyPlanID(ctx, familyPlanID)
}

// UpdateFamilyMember updates a family member
func (s *familyService) UpdateFamilyMember(ctx context.Context, member *models.FamilyMember) error {
	if member == nil {
		return errors.New("family member cannot be nil")
	}
	return s.familyRepo.UpdateFamilyMember(ctx, member)
}

// RemoveFamilyMember removes a member from a family plan
func (s *familyService) RemoveFamilyMember(ctx context.Context, memberID uuid.UUID) error {
	return s.familyRepo.DeleteFamilyMember(ctx, memberID)
}

// CalculateFamilyDiscount calculates discount for family plan
func (s *familyService) CalculateFamilyDiscount(ctx context.Context, planID uuid.UUID, basePremium float64) (float64, error) {
	plan, err := s.familyRepo.GetFamilyPlanByID(ctx, planID)
	if err != nil {
		return 0, fmt.Errorf("failed to get family plan: %w", err)
	}
	if plan == nil {
		return 0, errors.New("family plan not found")
	}
	// Apply group discount
	return basePremium * (1 - plan.GroupDiscount), nil
}

// UpdateFamilyPlanPremium updates the premium for a family plan
func (s *familyService) UpdateFamilyPlanPremium(ctx context.Context, planID uuid.UUID) error {
	plan, err := s.familyRepo.GetFamilyPlanByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("failed to get family plan: %w", err)
	}
	if plan == nil {
		return errors.New("family plan not found")
	}
	// Recalculate premium based on members and devices
	// Simplified calculation
	plan.CurrentPremium = plan.BasePremium * float64(plan.CurrentMembers)
	return s.familyRepo.UpdateFamilyPlan(ctx, plan)
}

// CheckMemberLimit checks if member limit has been reached
func (s *familyService) CheckMemberLimit(ctx context.Context, planID uuid.UUID) (bool, error) {
	plan, err := s.familyRepo.GetFamilyPlanByID(ctx, planID)
	if err != nil {
		return false, fmt.Errorf("failed to get family plan: %w", err)
	}
	if plan == nil {
		return false, errors.New("family plan not found")
	}
	return plan.CurrentMembers < plan.MaxMembers, nil
}

// CheckDeviceLimit checks if device limit has been reached for a member
func (s *familyService) CheckDeviceLimit(ctx context.Context, planID uuid.UUID, memberID uuid.UUID) (bool, error) {
	member, err := s.familyRepo.GetFamilyMemberByID(ctx, memberID)
	if err != nil {
		return false, fmt.Errorf("failed to get family member: %w", err)
	}
	if member == nil {
		return false, errors.New("family member not found")
	}
	return member.CurrentDevices < member.DeviceLimit, nil
}

// GetFamilyPlanStatistics gets statistics for a family plan
func (s *familyService) GetFamilyPlanStatistics(ctx context.Context, familyPlanID uuid.UUID) (map[string]interface{}, error) {
	return s.familyRepo.GetFamilyPlanStatistics(ctx, familyPlanID)
}

// GetFamilyMemberCount gets the count of family members
func (s *familyService) GetFamilyMemberCount(ctx context.Context, familyPlanID uuid.UUID) (int64, error) {
	return s.familyRepo.GetFamilyMemberCount(ctx, familyPlanID)
}

// GetActiveMemberCount gets the count of active family members
func (s *familyService) GetActiveMemberCount(ctx context.Context, familyPlanID uuid.UUID) (int64, error) {
	return s.familyRepo.GetActiveMemberCount(ctx, familyPlanID)
}
