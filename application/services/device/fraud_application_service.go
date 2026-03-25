package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"smartsure/internal/domain/models/device"
	domainServices "smartsure/internal/domain/services/device"
)

// FraudApplicationService handles fraud detection with infrastructure dependencies
// This is an APPLICATION SERVICE - orchestrates domain logic and external systems
type FraudApplicationService struct {
	db           *gorm.DB
	fraudService *domainServices.FraudDetectionService
	// In production, add external service clients here:
	// blacklistAPI     *BlacklistAPIClient
	// carrierAPI       *CarrierAPIClient
}

// NewFraudApplicationService creates a new fraud application service
func NewFraudApplicationService(db *gorm.DB) *FraudApplicationService {
	return &FraudApplicationService{
		db:           db,
		fraudService: domainServices.NewFraudDetectionService(),
	}
}

// CheckDeviceForFraud performs comprehensive fraud check using domain logic and external systems
func (s *FraudApplicationService) CheckDeviceForFraud(ctx context.Context, deviceID string) (*FraudCheckResult, error) {
	// 1. Load device from database
	var device device.Device
	if err := s.db.WithContext(ctx).
		Preload("Claims").
		Preload("Repairs").
		Preload("Lifecycle").
		First(&device, "id = ?", deviceID).Error; err != nil {
		return nil, fmt.Errorf("failed to load device: %w", err)
	}

	result := &FraudCheckResult{
		DeviceID:  deviceID,
		CheckedAt: time.Now(),
	}

	// 2. Use domain service for pure fraud detection logic
	result.RiskScore = s.fraudService.CalculateRiskScore(&device)
	result.Anomalies = s.fraudService.DetectAnomalies(&device)
	result.IsHighRisk = s.fraudService.IsHighRiskDevice(&device)

	// 3. Check external blacklist databases (infrastructure concern)
	blacklistStatus, err := s.checkBlacklist(ctx, device.IMEI)
	if err != nil {
		// Log error but don't fail the entire check
		fmt.Printf("Blacklist check failed: %v\n", err)
	} else {
		result.BlacklistStatus = blacklistStatus
	}

	// 4. Check with carrier for stolen status (infrastructure concern)
	if device.NetworkOperator != "" {
		stolenStatus, err := s.checkCarrierStolenDatabase(ctx, device.IMEI, device.NetworkOperator)
		if err != nil {
			fmt.Printf("Carrier check failed: %v\n", err)
		} else {
			result.CarrierStolenStatus = stolenStatus
		}
	}

	// 5. Save fraud check result to database
	if err := s.saveFraudCheckResult(ctx, result); err != nil {
		fmt.Printf("Failed to save fraud check result: %v\n", err)
	}

	// 6. Update device fraud score in database
	device.FraudRiskScore = int(result.RiskScore)
	if err := s.db.WithContext(ctx).Save(&device).Error; err != nil {
		fmt.Printf("Failed to update device fraud score: %v\n", err)
	}

	return result, nil
}

// CheckBlacklist checks device against external blacklist databases (infrastructure)
func (s *FraudApplicationService) checkBlacklist(ctx context.Context, imei string) (string, error) {
	// In production, this would call external blacklist APIs
	// For now, simulate with database check

	var blacklistEntry struct {
		IMEI   string
		Status string
		Reason string
	}

	err := s.db.WithContext(ctx).
		Table("device_blacklist").
		Where("imei = ?", imei).
		First(&blacklistEntry).Error

	if err == gorm.ErrRecordNotFound {
		return "clean", nil
	}
	if err != nil {
		return "unknown", err
	}

	return blacklistEntry.Status, nil
}

// CheckCarrierStolenDatabase checks with carrier for stolen device status (infrastructure)
func (s *FraudApplicationService) checkCarrierStolenDatabase(ctx context.Context, imei string, carrier string) (bool, error) {
	// In production, this would call carrier API
	// Example: GET https://api.carrier.com/stolen-devices/{imei}

	// For now, simulate with database check
	var stolenRecord struct {
		IMEI       string
		ReportedAt time.Time
	}

	err := s.db.WithContext(ctx).
		Table("carrier_stolen_devices").
		Where("imei = ? AND carrier = ?", imei, carrier).
		First(&stolenRecord).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// ReportSuspiciousDevice reports a device to authorities and blacklists (infrastructure)
func (s *FraudApplicationService) ReportSuspiciousDevice(ctx context.Context, deviceID string, reason string) error {
	// 1. Load device
	var device device.Device
	if err := s.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error; err != nil {
		return fmt.Errorf("failed to load device: %w", err)
	}

	// 2. Use domain service to verify it's actually suspicious
	if !s.fraudService.IsHighRiskDevice(&device) {
		return errors.New("device does not meet high risk criteria for reporting")
	}

	// 3. Report to external systems (infrastructure concerns)
	// In production:
	// - Report to carrier blacklist API
	// - Report to insurance fraud database
	// - Report to law enforcement API if stolen

	// 4. Update local blacklist database
	blacklistEntry := map[string]interface{}{
		"imei":        device.IMEI,
		"device_id":   deviceID,
		"reason":      reason,
		"reported_at": time.Now(),
		"status":      "blocked",
	}

	if err := s.db.WithContext(ctx).Table("device_blacklist").Create(blacklistEntry).Error; err != nil {
		return fmt.Errorf("failed to blacklist device: %w", err)
	}

	// 5. Update device status
	device.BlacklistStatus = "blocked"
	if err := s.db.WithContext(ctx).Save(&device).Error; err != nil {
		return fmt.Errorf("failed to update device status: %w", err)
	}

	return nil
}

// GetFraudHistory retrieves historical fraud checks for a device (infrastructure)
func (s *FraudApplicationService) GetFraudHistory(ctx context.Context, deviceID string) ([]FraudCheckResult, error) {
	var results []FraudCheckResult

	err := s.db.WithContext(ctx).
		Table("fraud_check_results").
		Where("device_id = ?", deviceID).
		Order("checked_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fraud history: %w", err)
	}

	return results, nil
}

// saveFraudCheckResult persists fraud check results (infrastructure)
func (s *FraudApplicationService) saveFraudCheckResult(ctx context.Context, result *FraudCheckResult) error {
	return s.db.WithContext(ctx).Table("fraud_check_results").Create(result).Error
}

// FraudCheckResult represents the outcome of a fraud check
type FraudCheckResult struct {
	ID                  uint      `json:"id"`
	DeviceID            string    `json:"device_id"`
	CheckedAt           time.Time `json:"checked_at"`
	RiskScore           float64   `json:"risk_score"`
	Anomalies           []string  `json:"anomalies" gorm:"serializer:json"`
	IsHighRisk          bool      `json:"is_high_risk"`
	BlacklistStatus     string    `json:"blacklist_status"`
	CarrierStolenStatus bool      `json:"carrier_stolen_status"`
}
