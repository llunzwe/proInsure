package device

import (
	"context"

	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// FraudService handles fraud-related operations for devices
type FraudService struct {
	db *gorm.DB
}

// NewFraudService creates a new fraud service
func NewFraudService(db *gorm.DB) *FraudService {
	return &FraudService{
		db: db,
	}
}

// CheckFraudStatus checks fraud status for a device
func (s *FraudService) CheckFraudStatus(ctx context.Context, deviceID string) (bool, error) {
	var device models.Device
	if err := s.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error; err != nil {
		return false, err
	}
	return device.DeviceStatusInfo.IsStolen, nil
}
