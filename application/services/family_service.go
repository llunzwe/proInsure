package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/policy"
	"smartsure/pkg/logger"
)

// FamilyService handles family plan and shared coverage management
type FamilyService struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewFamilyService creates a new FamilyService instance
func NewFamilyService(db *gorm.DB, logger *logger.Logger) *FamilyService {
	return &FamilyService{
		db:     db,
		logger: logger,
	}
}

// FamilyPlanRequest represents a request to create a family plan
type FamilyPlanRequest struct {
	PrimaryUserID    uuid.UUID `json:"primary_user_id" validate:"required"`
	PlanName         string    `json:"plan_name" validate:"required"`
	MaxMembers       int       `json:"max_members" validate:"required"`
	CoverageType     string    `json:"coverage_type" validate:"required"`
	SharedDeductible float64   `json:"shared_deductible"`
	BillingCycle     string    `json:"billing_cycle"`
}

// FamilyMemberRequest represents a request to add a family member
type FamilyMemberRequest struct {
	FamilyPlanID uuid.UUID `json:"family_plan_id" validate:"required"`
	UserID       uuid.UUID `json:"user_id" validate:"required"`
	Relationship string    `json:"relationship" validate:"required"`
	Role         string    `json:"role"`
	DateOfBirth  time.Time `json:"date_of_birth"`
}

// SharedCoverageRequest represents a request to set up shared coverage
type SharedCoverageRequest struct {
	FamilyPlanID     uuid.UUID `json:"family_plan_id" validate:"required"`
	CoverageType     string    `json:"coverage_type" validate:"required"`
	SharedLimit      float64   `json:"shared_limit" validate:"required"`
	IndividualLimit  float64   `json:"individual_limit"`
	DeductibleType   string    `json:"deductible_type"`
	SharedDeductible float64   `json:"shared_deductible"`
}

// GroupDiscountRequest represents a request to apply group discounts
type GroupDiscountRequest struct {
	FamilyPlanID   uuid.UUID `json:"family_plan_id" validate:"required"`
	DiscountType   string    `json:"discount_type" validate:"required"`
	DiscountReason string    `json:"discount_reason"`
	ValidUntil     time.Time `json:"valid_until"`
}

// CreateFamilyPlan creates a new family insurance plan
func (s *FamilyService) CreateFamilyPlan(ctx context.Context, req *FamilyPlanRequest) (*models.FamilyPlan, error) {
	s.logger.Info("Creating family plan", "primary_user_id", req.PrimaryUserID, "plan_name", req.PlanName)

	// Verify primary user exists
	var primaryUser models.User
	if err := s.db.WithContext(ctx).First(&primaryUser, "id = ?", req.PrimaryUserID).Error; err != nil {
		return nil, errors.New("primary user not found")
	}

	// Check if user already has an active family plan
	var existingPlan models.FamilyPlan
	if err := s.db.WithContext(ctx).Where("primary_user_id = ? AND is_active = ?", req.PrimaryUserID, true).First(&existingPlan).Error; err == nil {
		return nil, errors.New("user already has an active family plan")
	}

	// Generate plan number (for future use)
	_ = s.generateFamilyPlanNumber()

	// Calculate base premium for family plan
	basePremium := s.calculateFamilyPlanPremium(req)

	// Create family plan
	plan := &models.FamilyPlan{
		PrimaryUserID:        req.PrimaryUserID,
		PlanName:             req.PlanName,
		MaxMembers:           req.MaxMembers,
		CurrentMembers:       1, // Primary user
		MaxDevicesPerMember:  3,
		BasePremium:          basePremium,
		CurrentPremium:       basePremium,
		SharedDeductible:     req.SharedDeductible,
		IndividualDeductible: req.SharedDeductible / 2,
		TotalCoverageLimit:   float64(req.MaxMembers) * 5000.0,
		PerMemberLimit:       5000.0,
		GroupDiscount:        0.15,
		BillingCycle:         req.BillingCycle,
		Status:               "active",
		IsAutoRenew:          true,
		StartDate:            time.Now(),
	}

	// Set plan benefits
	benefits := map[string]interface{}{
		"shared_coverage":       true,
		"family_deductible":     req.SharedDeductible,
		"multi_device_discount": 0.15,
		"family_claim_sharing":  true,
		"emergency_assistance":  true,
		"24_7_support":          true,
	}
	benefitsJSON, _ := json.Marshal(benefits)
	plan.SharedBenefits = string(benefitsJSON)

	if err := s.db.WithContext(ctx).Create(plan).Error; err != nil {
		return nil, fmt.Errorf("failed to create family plan: %w", err)
	}

	// Add primary user as first family member
	_, err := s.addFamilyMember(ctx, plan.ID, req.PrimaryUserID, "primary", "admin", time.Time{})
	if err != nil {
		s.logger.Error("Failed to add primary user as family member", "error", err)
	}

	s.logger.Info("Family plan created successfully", "plan_id", plan.ID)
	return plan, nil
}

// AddFamilyMember adds a new member to a family plan
func (s *FamilyService) AddFamilyMember(ctx context.Context, req *FamilyMemberRequest) (*models.FamilyMember, error) {
	s.logger.Info("Adding family member", "plan_id", req.FamilyPlanID, "user_id", req.UserID)

	// Verify family plan exists and is active
	var plan models.FamilyPlan
	if err := s.db.WithContext(ctx).First(&plan, "id = ? AND is_active = ?", req.FamilyPlanID, true).Error; err != nil {
		return nil, errors.New("family plan not found or inactive")
	}

	// Check if plan has reached member limit
	if plan.CurrentMembers >= plan.MaxMembers {
		return nil, errors.New("family plan has reached maximum member limit")
	}

	// Verify user exists
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", req.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is already a member of any family plan
	var existingMember models.FamilyMember
	if err := s.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", req.UserID, true).First(&existingMember).Error; err == nil {
		return nil, errors.New("user is already a member of a family plan")
	}

	return s.addFamilyMember(ctx, req.FamilyPlanID, req.UserID, req.Relationship, req.Role, req.DateOfBirth)
}

// SetupSharedCoverage sets up shared coverage for a family plan
func (s *FamilyService) SetupSharedCoverage(ctx context.Context, req *SharedCoverageRequest) (*policy.SharedCoverage, error) {
	s.logger.Info("Setting up shared coverage", "plan_id", req.FamilyPlanID, "coverage_type", req.CoverageType)

	// Verify family plan exists
	var plan models.FamilyPlan
	if err := s.db.WithContext(ctx).First(&plan, "id = ? AND is_active = ?", req.FamilyPlanID, true).Error; err != nil {
		return nil, errors.New("family plan not found or inactive")
	}

	// Check if shared coverage already exists for this type
	var existingCoverage policy.SharedCoverage
	if err := s.db.WithContext(ctx).Where("family_plan_id = ? AND coverage_type = ?", req.FamilyPlanID, req.CoverageType).First(&existingCoverage).Error; err == nil {
		return nil, errors.New("shared coverage already exists for this type")
	}

	// Create shared coverage
	coverage := &policy.SharedCoverage{
		PlanType:        "family",
		PlanID:          req.FamilyPlanID,
		CoverageType:    "deductible",
		SharedAmount:    req.SharedLimit,
		UsedAmount:      0.0,
		RemainingAmount: req.SharedLimit,
		ResetPeriod:     "annual",
		LastReset:       time.Now(),
		NextReset:       time.Now().AddDate(1, 0, 0),
		IsActive:        true,
	}

	// Set coverage rules
	rules := map[string]interface{}{
		"max_claims_per_member":    3,
		"max_claims_per_year":      10,
		"requires_approval_above":  1000.0,
		"family_approval_required": false,
		"emergency_override":       true,
	}
	// Rules are stored in the coverage configuration
	_ = rules // Rules will be implemented in future version

	if err := s.db.WithContext(ctx).Create(coverage).Error; err != nil {
		return nil, fmt.Errorf("failed to create shared coverage: %w", err)
	}

	s.logger.Info("Shared coverage created successfully", "coverage_id", coverage.ID)
	return coverage, nil
}

// ApplyGroupDiscount applies a group discount to a family plan
func (s *FamilyService) ApplyGroupDiscount(ctx context.Context, req *GroupDiscountRequest) (*policy.GroupDiscount, error) {
	s.logger.Info("Applying group discount", "plan_id", req.FamilyPlanID, "discount_type", req.DiscountType)

	// Verify family plan exists
	var plan models.FamilyPlan
	if err := s.db.WithContext(ctx).First(&plan, "id = ? AND is_active = ?", req.FamilyPlanID, true).Error; err != nil {
		return nil, errors.New("family plan not found or inactive")
	}

	// Calculate discount percentage and amount
	discountPercentage := s.calculateDiscountPercentage(req.DiscountType, plan.CurrentMembers)
	discountAmount := plan.CurrentPremium * discountPercentage

	// Create group discount
	discount := &policy.GroupDiscount{
		DiscountName:   "Family Plan Discount",
		DiscountType:   req.DiscountType,
		MinimumMembers: 2,
		MaximumMembers: plan.MaxMembers,
		DiscountRate:   discountPercentage,
	}

	// Set discount conditions
	conditions := map[string]interface{}{
		"minimum_members":       2,
		"minimum_premium":       100.0,
		"requires_auto_renewal": true,
		"stackable":             false,
		"applies_to":            "premium",
	}
	// Conditions are stored in the discount configuration
	_ = conditions // Conditions will be implemented in future version

	if err := s.db.WithContext(ctx).Create(discount).Error; err != nil {
		return nil, fmt.Errorf("failed to create group discount: %w", err)
	}

	// Update family plan premium
	plan.CurrentPremium = plan.BasePremium - discountAmount
	if err := s.db.WithContext(ctx).Save(&plan).Error; err != nil {
		s.logger.Error("Failed to update plan premium", "error", err)
	}

	s.logger.Info("Group discount applied successfully", "discount_id", discount.ID)
	return discount, nil
}

// GetFamilyPlanAnalytics generates analytics for a family plan
func (s *FamilyService) GetFamilyPlanAnalytics(ctx context.Context, planID uuid.UUID) (map[string]interface{}, error) {
	s.logger.Info("Generating family plan analytics", "plan_id", planID)

	analytics := make(map[string]interface{})

	// Get family plan details
	var plan models.FamilyPlan
	if err := s.db.WithContext(ctx).Preload("FamilyMembers").First(&plan, "id = ?", planID).Error; err != nil {
		return nil, fmt.Errorf("family plan not found: %w", err)
	}

	analytics["plan_id"] = plan.ID
	analytics["plan_name"] = plan.PlanName
	analytics["current_members"] = plan.CurrentMembers
	analytics["max_members"] = plan.MaxMembers
	analytics["utilization_rate"] = float64(plan.CurrentMembers) / float64(plan.MaxMembers)

	// Premium analytics
	analytics["base_premium"] = plan.BasePremium
	analytics["current_premium"] = plan.CurrentPremium
	analytics["savings_amount"] = plan.BasePremium - plan.CurrentPremium
	analytics["savings_percentage"] = (plan.BasePremium - plan.CurrentPremium) / plan.BasePremium * 100

	// Member analytics - query members separately to avoid circular dependency
	var members []models.FamilyMember
	if err := s.db.Where("family_plan_id = ?", plan.ID).Find(&members).Error; err != nil {
		// If we can't get members, set defaults
		analytics["member_count"] = 0
		analytics["member_relationships"] = make(map[string]int)
		return analytics, nil
	}

	memberRelationships := make(map[string]int)
	for _, member := range members {
		// Age calculation would require User relationship to get DateOfBirth
		// For now, just count relationships
		memberRelationships[member.Relationship]++
	}

	// Member ages would be calculated from User relationship
	analytics["member_count"] = len(members)
	analytics["member_relationships"] = memberRelationships

	// Coverage analytics
	var sharedCoverages []policy.SharedCoverage
	s.db.WithContext(ctx).Where("family_plan_id = ?", planID).Find(&sharedCoverages)

	totalSharedLimit := 0.0
	totalUsedAmount := 0.0
	coverageTypes := []string{}

	for _, coverage := range sharedCoverages {
		totalSharedLimit += coverage.SharedAmount
		totalUsedAmount += coverage.UsedAmount
		coverageTypes = append(coverageTypes, coverage.CoverageType)
	}

	analytics["total_shared_limit"] = totalSharedLimit
	analytics["total_used_amount"] = totalUsedAmount
	analytics["coverage_utilization"] = totalUsedAmount / totalSharedLimit * 100
	analytics["coverage_types"] = coverageTypes

	// Discount analytics
	var activeDiscounts []policy.GroupDiscount
	s.db.WithContext(ctx).Where("family_plan_id = ? AND is_active = ? AND valid_until > ?", planID, true, time.Now()).Find(&activeDiscounts)

	totalDiscountAmount := 0.0
	discountTypes := []string{}

	for _, discount := range activeDiscounts {
		totalDiscountAmount += discount.DiscountRate
		discountTypes = append(discountTypes, discount.DiscountType)
	}

	analytics["active_discounts"] = len(activeDiscounts)
	analytics["total_discount_amount"] = totalDiscountAmount
	analytics["discount_types"] = discountTypes

	return analytics, nil
}

// Helper methods

func (s *FamilyService) addFamilyMember(ctx context.Context, planID, userID uuid.UUID, relationship, role string, dateOfBirth time.Time) (*models.FamilyMember, error) {
	member := &models.FamilyMember{
		FamilyPlanID:      planID,
		UserID:            userID,
		Relationship:      relationship,
		JoinDate:          time.Now(),
		Status:            "active",
		IsMinor:           false,
		DeviceLimit:       3,
		CurrentDevices:    0,
		IndividualPremium: 0.0,
		CoverageLevel:     "standard",
	}

	// Set member restrictions based on relationship
	restrictions := map[string]interface{}{
		"can_make_claims":   true,
		"requires_approval": relationship == "child",
		"spending_limit":    1000.0,
	}
	restrictionsJSON, _ := json.Marshal(restrictions)
	member.Restrictions = string(restrictionsJSON)

	if err := s.db.WithContext(ctx).Create(member).Error; err != nil {
		return nil, fmt.Errorf("failed to add family member: %w", err)
	}

	// Update family plan member count
	if err := s.db.WithContext(ctx).Model(&models.FamilyPlan{}).Where("id = ?", planID).UpdateColumn("current_members", gorm.Expr("current_members + ?", 1)).Error; err != nil {
		s.logger.Error("Failed to update member count", "error", err)
	}

	// Recalculate family plan premium
	s.recalculateFamilyPlanPremium(ctx, planID)

	return member, nil
}

func (s *FamilyService) generateFamilyPlanNumber() string {
	return fmt.Sprintf("FAM-%d-%s", time.Now().Unix(), uuid.New().String()[:8])
}

func (s *FamilyService) calculateFamilyPlanPremium(req *FamilyPlanRequest) float64 {
	basePremium := 100.0 // Base premium for primary member

	// Coverage type multiplier
	coverageMultipliers := map[string]float64{
		"basic":         1.0,
		"standard":      1.5,
		"comprehensive": 2.0,
		"premium":       2.5,
	}

	multiplier := coverageMultipliers[req.CoverageType]
	if multiplier == 0 {
		multiplier = 1.0
	}

	// Family size factor (discount for larger families)
	familyDiscount := 1.0
	if req.MaxMembers >= 6 {
		familyDiscount = 0.80 // 20% discount for large families
	} else if req.MaxMembers >= 4 {
		familyDiscount = 0.90 // 10% discount for medium families
	} else if req.MaxMembers >= 3 {
		familyDiscount = 0.95 // 5% discount for small families
	}

	return basePremium * multiplier * familyDiscount
}

func (s *FamilyService) calculateDiscountPercentage(discountType string, memberCount int) float64 {
	switch discountType {
	case "multi_member":
		if memberCount >= 5 {
			return 0.20 // 20% discount for 5+ members
		} else if memberCount >= 3 {
			return 0.15 // 15% discount for 3-4 members
		} else if memberCount >= 2 {
			return 0.10 // 10% discount for 2 members
		}
	case "loyalty":
		return 0.05 // 5% loyalty discount
	case "promotional":
		return 0.25 // 25% promotional discount
	case "referral":
		return 0.10 // 10% referral discount
	}
	return 0.0
}

func (s *FamilyService) getMemberPermissions(role, relationship string) map[string]interface{} {
	permissions := map[string]interface{}{
		"can_view_plan":       true,
		"can_file_claims":     true,
		"can_add_devices":     false,
		"can_modify_plan":     false,
		"can_invite_members":  false,
		"can_remove_members":  false,
		"can_view_all_claims": false,
		"requires_approval":   true,
	}

	switch role {
	case "admin":
		permissions["can_add_devices"] = true
		permissions["can_modify_plan"] = true
		permissions["can_invite_members"] = true
		permissions["can_remove_members"] = true
		permissions["can_view_all_claims"] = true
		permissions["requires_approval"] = false
	case "manager":
		permissions["can_add_devices"] = true
		permissions["can_invite_members"] = true
		permissions["can_view_all_claims"] = true
		permissions["requires_approval"] = false
	case "member":
		// Default permissions already set
	}

	// Relationship-based adjustments
	if relationship == "primary" {
		permissions["can_modify_plan"] = true
		permissions["can_remove_members"] = true
		permissions["requires_approval"] = false
	}

	return permissions
}

func (s *FamilyService) recalculateFamilyPlanPremium(ctx context.Context, planID uuid.UUID) {
	var plan models.FamilyPlan
	if err := s.db.WithContext(ctx).First(&plan, "id = ?", planID).Error; err != nil {
		return
	}

	// Recalculate premium based on current member count
	memberMultiplier := 1.0 + (float64(plan.CurrentMembers-1) * 0.3) // 30% per additional member
	newPremium := plan.BasePremium * memberMultiplier

	// Apply family discount
	if plan.CurrentMembers >= 4 {
		newPremium *= 0.85 // 15% family discount
	} else if plan.CurrentMembers >= 3 {
		newPremium *= 0.90 // 10% family discount
	}

	plan.CurrentPremium = newPremium
	s.db.WithContext(ctx).Save(&plan)
}
