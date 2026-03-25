package device

import (
	"context"

	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// DeviceFraud handles device-specific fraud operations
type DeviceFraud struct {
	db *gorm.DB
}

// NewDeviceFraud creates a new device fraud handler
func NewDeviceFraud(db *gorm.DB) *DeviceFraud {
	return &DeviceFraud{
		db: db,
	}
}

// IsBlacklisted checks if device is blacklisted
func (d *DeviceFraud) IsBlacklisted(ctx context.Context, deviceID string) (bool, error) {
	var device models.Device
	if err := d.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error; err != nil {
		return false, err
	}
	return device.DeviceVerification.AuthenticityStatus == "counterfeit", nil
}
