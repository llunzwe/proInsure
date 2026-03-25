package device

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	domainServices "smartsure/internal/domain/services/device"
)

// ClaimsApplicationService handles claims with infrastructure dependencies
// This is an APPLICATION SERVICE - orchestrates domain logic and database operations
type ClaimsApplicationService struct {
	db                 *gorm.DB
	eligibilityService *domainServices.ClaimsEligibilityService
}

// NewClaimsApplicationService creates a new claims application service
func NewClaimsApplicationService(db *gorm.DB) *ClaimsApplicationService {
	return &ClaimsApplicationService{
		db:                 db,
		eligibilityService: domainServices.NewClaimsEligibilityService(nil, nil, nil),
	}
}

// ClaimEligibilityResult represents the result of eligibility check
type ClaimEligibilityResult struct {
	DeviceID        string    `json:"device_id"`
	ClaimType       string    `json:"claim_type"`
	IsEligible      bool      `json:"is_eligible"`
	Reason          string    `json:"reason,omitempty"`
	Deductible      float64   `json:"deductible,omitempty"`
	CoverageLimit   float64   `json:"coverage_limit,omitempty"`
	EstimatedPayout float64   `json:"estimated_payout,omitempty"`
	ClaimFrequency  float64   `json:"claim_frequency,omitempty"`
	CheckedAt       time.Time `json:"checked_at"`
}

// ProcessClaimEligibility checks eligibility and provides quote for a claim
func (s *ClaimsApplicationService) ProcessClaimEligibility(ctx context.Context, deviceID string, claimType string) (*ClaimEligibilityResult, error) {
	var dev models.Device
	if err := s.db.WithContext(ctx).
		Preload("Claims").
		Preload("Repairs").
		First(&dev, "id = ?", deviceID).Error; err != nil {
		return nil, fmt.Errorf("failed to load device: %w", err)
	}

	hasActiveClaim := s.hasActiveClaim(ctx, deviceID)
	if hasActiveClaim {
		return &ClaimEligibilityResult{
			DeviceID:   deviceID,
			ClaimType:  claimType,
			IsEligible: false,
			Reason:     "Device has an active claim in progress",
			CheckedAt:  time.Now(),
		}, nil
	}

	result := &ClaimEligibilityResult{
		DeviceID:   deviceID,
		ClaimType:  claimType,
		IsEligible: true,
		CheckedAt:  time.Now(),
	}

	return result, nil
}

// hasActiveClaim checks if device has active claims
func (s *ClaimsApplicationService) hasActiveClaim(ctx context.Context, deviceID string) bool {
	var count int64
	s.db.WithContext(ctx).
		Model(&models.Claim{}).
		Where("device_id = ? AND lc_status IN ('submitted', 'under_review', 'approved')", deviceID).
		Count(&count)
	return count > 0
}
