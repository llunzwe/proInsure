package device

import (
	"context"
	"time"

	"smartsure/internal/domain/models/device"

	"gorm.io/gorm"
)

// DeviceControlService is the main service that coordinates all device control operations
type DeviceControlService struct {
	db               *gorm.DB
	ecosystemService *DeviceEcosystemService
	fraudAppService  *FraudApplicationService
}

// NewDeviceControlService creates a new device control service
func NewDeviceControlService(db *gorm.DB) *DeviceControlService {
	return &DeviceControlService{
		db:               db,
		ecosystemService: NewDeviceEcosystemService(db),
		fraudAppService:  NewFraudApplicationService(db),
	}
}

// Fraud Detection Methods (delegated to fraud application service)
func (s *DeviceControlService) CheckDeviceForFraud(ctx context.Context, deviceID string) (*FraudCheckResult, error) {
	return s.fraudAppService.CheckDeviceForFraud(ctx, deviceID)
}

func (s *DeviceControlService) ReportSuspiciousDevice(ctx context.Context, deviceID string, reason string) error {
	return s.fraudAppService.ReportSuspiciousDevice(ctx, deviceID, reason)
}

func (s *DeviceControlService) GetFraudHistory(ctx context.Context, deviceID string) ([]FraudCheckResult, error) {
	return s.fraudAppService.GetFraudHistory(ctx, deviceID)
}

// GetDevice retrieves a device by ID using the ecosystem service
func (s *DeviceControlService) GetDevice(ctx context.Context, deviceID string) (*device.Device, error) {
	return s.ecosystemService.GetDevice(ctx, deviceID)
}

// RegisterDevice registers a new device
func (s *DeviceControlService) RegisterDevice(ctx context.Context, req *RegisterDeviceRequest) (*device.Device, error) {
	return s.ecosystemService.RegisterDevice(ctx, req)
}

// VerifyDevice performs device verification
func (s *DeviceControlService) VerifyDevice(ctx context.Context, deviceID string, req *VerifyDeviceRequest) error {
	return s.ecosystemService.VerifyDevice(ctx, deviceID, req)
}

// Helper function to get repair count (example of accessing related data)
func GetRepairCount(d *device.Device) int {
	if d.Repairs == nil {
		return 0
	}
	return len(d.Repairs)
}

// Helper function to get device age
func GetDeviceAge(d *device.Device) int {
	if d.PurchaseDate != nil {
		return int(time.Since(*d.PurchaseDate).Hours() / 24)
	}
	return int(time.Since(d.RegistrationDate).Hours() / 24)
}

// Helper function to check warranty status
func IsWarrantyActive(d *device.Device) bool {
	if d.WarrantyExpiry == nil {
		return false
	}
	return time.Now().Before(*d.WarrantyExpiry)
}

// Helper function to check if device can be insured
func CanBeInsured(d *device.Device) bool {
	// Basic eligibility check
	if d.Status != "active" || d.IsStolen || !d.IsVerified {
		return false
	}

	// Age check
	if GetDeviceAge(d) > 1460 { // 4 years
		return false
	}

	// Value check
	if d.CurrentValue < 100 {
		return false
	}

	// Blacklist check
	if d.BlacklistStatus == "blocked" {
		return false
	}

	return true
}

// Helper function to calculate condition score
func GetConditionScore(d *device.Device) float64 {
	scores := map[string]float64{
		"excellent": 1.0,
		"good":      0.85,
		"fair":      0.7,
		"poor":      0.5,
	}

	score, exists := scores[d.Condition]
	if !exists {
		score = 0.7
	}

	// Adjust based on screen condition
	if d.ScreenCondition == "cracked" {
		score *= 0.7
	} else if d.ScreenCondition == "broken" {
		score *= 0.5
	}

	return score
}
