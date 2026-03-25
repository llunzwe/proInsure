package device

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// FraudService handles fraud detection and prevention for devices
type FraudService struct {
	db *gorm.DB
}

// NewFraudService creates a new fraud service
func NewFraudService(db *gorm.DB) *FraudService {
	return &FraudService{
		db: db,
	}
}

// DetectFraud performs fraud detection on a device
func (s *FraudService) DetectFraud(ctx context.Context, deviceID string) (float64, error) {
	var device models.Device
	if err := s.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error; err != nil {
		return 0, fmt.Errorf("failed to load device: %w", err)
	}

	riskScore := device.DeviceRiskAssessment.RiskScore
	if device.DeviceStatusInfo.IsStolen {
		riskScore = 100.0
	}

	return riskScore, nil
}
