package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models/shared"
)

// UpdateCommunication updates a communication
func (s *CommunicationService) UpdateCommunication(ctx context.Context, commID uuid.UUID, updates map[string]interface{}) (*shared.Communication, error) {
	var comm shared.Communication
	if err := s.db.Where("id = ?", commID).First(&comm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("communication not found")
		}
		return nil, fmt.Errorf("failed to find communication: %w", err)
	}

	// Remove sensitive fields
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "sender_id")

	if err := s.db.Model(&comm).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update communication: %w", err)
	}

	return &comm, nil
}

// DeleteCommunication soft deletes a communication
func (s *CommunicationService) DeleteCommunication(ctx context.Context, commID uuid.UUID) error {
	result := s.db.Where("id = ?", commID).Delete(&shared.Communication{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete communication: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("communication not found")
	}
	return nil
}

// GetTemplates retrieves all templates
func (s *CommunicationService) GetTemplates(ctx context.Context, offset, limit int) ([]*shared.CommunicationTemplate, int64, error) {
	var templates []*shared.CommunicationTemplate
	query := s.db.Model(&shared.CommunicationTemplate{})

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get templates: %w", err)
	}

	return templates, total, nil
}

// UpdateTemplate updates a template
func (s *CommunicationService) UpdateTemplate(ctx context.Context, templateID uuid.UUID, updates map[string]interface{}) (*shared.CommunicationTemplate, error) {
	var template shared.CommunicationTemplate
	if err := s.db.Where("id = ?", templateID).First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to find template: %w", err)
	}

	// Remove sensitive fields
	delete(updates, "id")
	delete(updates, "created_at")

	if err := s.db.Model(&template).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return &template, nil
}

// DeleteTemplate soft deletes a template
func (s *CommunicationService) DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	result := s.db.Where("id = ?", templateID).Delete(&shared.CommunicationTemplate{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete template: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("template not found")
	}
	return nil
}

// GetCampaigns retrieves all campaigns
func (s *CommunicationService) GetCampaigns(ctx context.Context, offset, limit int) ([]*shared.CommunicationCampaign, int64, error) {
	var campaigns []*shared.CommunicationCampaign
	query := s.db.Model(&shared.CommunicationCampaign{})

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&campaigns).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get campaigns: %w", err)
	}

	return campaigns, total, nil
}

// GetCampaign retrieves a campaign by ID
func (s *CommunicationService) GetCampaign(ctx context.Context, campaignID uuid.UUID) (*shared.CommunicationCampaign, error) {
	var campaign shared.CommunicationCampaign
	if err := s.db.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("campaign not found")
		}
		return nil, fmt.Errorf("failed to find campaign: %w", err)
	}
	return &campaign, nil
}

// UpdateCampaign updates a campaign
func (s *CommunicationService) UpdateCampaign(ctx context.Context, campaignID uuid.UUID, updates map[string]interface{}) (*shared.CommunicationCampaign, error) {
	var campaign shared.CommunicationCampaign
	if err := s.db.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("campaign not found")
		}
		return nil, fmt.Errorf("failed to find campaign: %w", err)
	}

	// Remove sensitive fields
	delete(updates, "id")
	delete(updates, "created_at")

	if err := s.db.Model(&campaign).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return &campaign, nil
}

// DeleteCampaign soft deletes a campaign
func (s *CommunicationService) DeleteCampaign(ctx context.Context, campaignID uuid.UUID) error {
	result := s.db.Where("id = ?", campaignID).Delete(&shared.CommunicationCampaign{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete campaign: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("campaign not found")
	}
	return nil
}
