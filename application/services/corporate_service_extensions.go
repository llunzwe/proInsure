package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// GetCorporateAccounts retrieves all corporate accounts
func (s *CorporateService) GetCorporateAccounts(ctx context.Context, offset, limit int) ([]*models.CorporateAccount, int64, error) {
	var accounts []*models.CorporateAccount
	query := s.db.Model(&models.CorporateAccount{})

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&accounts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get corporate accounts: %w", err)
	}

	return accounts, total, nil
}

// UpdateCorporateAccount updates a corporate account
func (s *CorporateService) UpdateCorporateAccount(ctx context.Context, accountID uuid.UUID, updates map[string]interface{}) (*models.CorporateAccount, error) {
	var account models.CorporateAccount
	if err := s.db.Where("id = ?", accountID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("corporate account not found")
		}
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	// Remove sensitive fields
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "company_id")

	if err := s.db.Model(&account).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update corporate account: %w", err)
	}

	return &account, nil
}

// DeleteCorporateAccount soft deletes a corporate account
func (s *CorporateService) DeleteCorporateAccount(ctx context.Context, accountID uuid.UUID) error {
	result := s.db.Where("id = ?", accountID).Delete(&models.CorporateAccount{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete corporate account: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("corporate account not found")
	}
	return nil
}

// GetEmployeePolicies retrieves employee policies
func (s *CorporateService) GetEmployeePolicies(ctx context.Context, accountID uuid.UUID, offset, limit int) ([]*models.CorporatePolicy, int64, error) {
	var policies []*models.CorporatePolicy
	query := s.db.Model(&models.CorporatePolicy{}).Where("corporate_account_id = ?", accountID)

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&policies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get employee policies: %w", err)
	}

	return policies, total, nil
}

// GetEmployeePolicy retrieves a specific employee policy
func (s *CorporateService) GetEmployeePolicy(ctx context.Context, policyID uuid.UUID) (*models.CorporatePolicy, error) {
	var policy models.CorporatePolicy
	if err := s.db.Where("id = ?", policyID).First(&policy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("employee policy not found")
		}
		return nil, fmt.Errorf("failed to find policy: %w", err)
	}
	return &policy, nil
}

// DeleteEmployeePolicy soft deletes an employee policy
func (s *CorporateService) DeleteEmployeePolicy(ctx context.Context, policyID uuid.UUID) error {
	result := s.db.Where("id = ?", policyID).Delete(&models.CorporatePolicy{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee policy: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee policy not found")
	}
	return nil
}
