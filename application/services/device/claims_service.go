package device

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// ClaimsService handles device-related claims processing
type ClaimsService struct {
	db *gorm.DB
}

// NewClaimsService creates a new claims service
func NewClaimsService(db *gorm.DB) *ClaimsService {
	return &ClaimsService{
		db: db,
	}
}

// FileClaimRequest represents a request to file a claim
type FileClaimRequest struct {
	DeviceID     string    `json:"device_id"`
	ClaimType    string    `json:"claim_type"`
	Description  string    `json:"description"`
	IncidentDate time.Time `json:"incident_date"`
}

// FileClaim initiates a new claim for a device
func (s *ClaimsService) FileClaim(ctx context.Context, request *FileClaimRequest) (*models.Claim, error) {
	deviceUUID, err := uuid.Parse(request.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("invalid device ID: %w", err)
	}

	claim := &models.Claim{
		DeviceID: deviceUUID,
	}
	claim.ClaimIdentification.ClaimType = request.ClaimType
	claim.ClaimDocumentation.Description = request.Description
	claim.ClaimLifecycle.IncidentDate = request.IncidentDate
	claim.ClaimLifecycle.Status = "submitted"

	if err := s.db.WithContext(ctx).Create(claim).Error; err != nil {
		return nil, fmt.Errorf("failed to create claim: %w", err)
	}

	return claim, nil
}
