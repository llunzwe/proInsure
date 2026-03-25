package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/device"
)

// DeviceStateService handles device state transitions and status changes
type DeviceStateService struct {
	deviceService *deviceService
}

// NewDeviceStateService creates a new device state service
func NewDeviceStateService(deviceService *deviceService) *DeviceStateService {
	return &DeviceStateService{
		deviceService: deviceService,
	}
}

// MarkAsStolen marks a device as stolen
func (s *DeviceStateService) MarkAsStolen(ctx context.Context, deviceID uuid.UUID, date time.Time, reportID string) error {
	dev, err := s.deviceService.repo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if dev.DeviceStatusInfo.IsStolen {
		return errors.New("device is already marked as stolen")
	}

	// Update device state
	dev.DeviceStatusInfo.IsStolen = true
	dev.DeviceStatusInfo.Status = device.DeviceStatusStolen
	dev.DeviceStatusInfo.StolenDate = &date
	dev.DeviceStatusInfo.StolenReportID = reportID

	// Update risk assessment
	dev.DeviceRiskAssessment.RiskScore = 100
	dev.DeviceRiskAssessment.TheftRiskLevel = device.RiskLevelCritical

	// Note: Insurance eligibility should be tracked via DeviceInsurance relationship
	// DeviceCompliance tracks regulatory compliance, not insurance eligibility
	// Insurance eligibility is determined by business logic, not stored in compliance

	// Update blacklist status
	dev.DeviceRiskAssessment.BlacklistStatus = device.BlacklistStatusBlocked

	// Save changes
	return s.deviceService.repo.Update(ctx, dev)
}

// MarkAsLost marks a device as lost
func (s *DeviceStateService) MarkAsLost(ctx context.Context, deviceID uuid.UUID, date time.Time) error {
	dev, err := s.deviceService.repo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if dev.DeviceStatusInfo.Status == device.DeviceStatusLost {
		return errors.New("device is already marked as lost")
	}

	// Update device state
	dev.DeviceStatusInfo.Status = device.DeviceStatusLost
	dev.DeviceStatusInfo.IsLost = true
	dev.DeviceStatusInfo.LostDate = &date

	// Update risk assessment
	if dev.DeviceRiskAssessment.RiskScore < 70 {
		dev.DeviceRiskAssessment.RiskScore = 70
	}
	dev.DeviceRiskAssessment.TheftRiskLevel = s.calculateRiskLevel(dev.DeviceRiskAssessment.RiskScore)

	// Note: Insurance eligibility should be tracked via DeviceInsurance relationship
	// DeviceCompliance tracks regulatory compliance, not insurance eligibility
	// Insurance eligibility is determined by business logic, not stored in compliance

	// Save changes
	return s.deviceService.repo.Update(ctx, dev)
}

// MarkAsRecovered marks a stolen/lost device as recovered
func (s *DeviceStateService) MarkAsRecovered(ctx context.Context, deviceID uuid.UUID) error {
	dev, err := s.deviceService.repo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if !dev.DeviceStatusInfo.IsStolen && !dev.DeviceStatusInfo.IsLost {
		return errors.New("device is not marked as stolen or lost")
	}

	wasStolen := dev.DeviceStatusInfo.IsStolen

	// Clear stolen/lost status
	dev.DeviceStatusInfo.IsStolen = false
	dev.DeviceStatusInfo.IsLost = false
	dev.DeviceStatusInfo.Status = device.DeviceStatusActive
	dev.DeviceStatusInfo.StolenDate = nil
	dev.DeviceStatusInfo.LostDate = nil

	// Update blacklist status if was stolen
	if wasStolen {
		dev.DeviceRiskAssessment.BlacklistStatus = device.BlacklistStatusClean
	}

	// Recalculate risk score
	dev.DeviceRiskAssessment.RiskScore = s.calculateBaseRiskScore(dev)
	dev.DeviceRiskAssessment.TheftRiskLevel = s.calculateRiskLevel(dev.DeviceRiskAssessment.RiskScore)

	// Note: Insurance eligibility should be tracked via DeviceInsurance relationship
	// DeviceCompliance tracks regulatory compliance, not insurance eligibility
	// Insurance eligibility is determined by business logic, not stored in compliance

	// Save changes
	return s.deviceService.repo.Update(ctx, dev)
}

// LockDevice locks a device
func (s *DeviceStateService) LockDevice(ctx context.Context, deviceID uuid.UUID, reason string) error {
	dev, err := s.deviceService.repo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if dev.DeviceStatusInfo.IsLocked {
		return errors.New("device is already locked")
	}

	dev.DeviceStatusInfo.IsLocked = true
	dev.DeviceStatusInfo.LockedReason = reason

	return s.deviceService.repo.Update(ctx, dev)
}

// UnlockDevice unlocks a device
func (s *DeviceStateService) UnlockDevice(ctx context.Context, deviceID uuid.UUID) error {
	dev, err := s.deviceService.repo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if !dev.DeviceStatusInfo.IsLocked {
		return errors.New("device is not locked")
	}

	dev.DeviceStatusInfo.IsLocked = false
	dev.DeviceStatusInfo.LockedReason = ""

	return s.deviceService.repo.Update(ctx, dev)
}

// RetireDevice retires a device
func (s *DeviceStateService) RetireDevice(ctx context.Context, deviceID uuid.UUID, reason string) error {
	dev, err := s.deviceService.repo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if dev.DeviceStatusInfo.Status == device.DeviceStatusRetired {
		return errors.New("device is already retired")
	}

	now := time.Now()
	dev.DeviceStatusInfo.Status = device.DeviceStatusRetired
	dev.DeviceStatusInfo.IsRetired = true
	dev.DeviceStatusInfo.RetiredDate = &now
	dev.DeviceStatusInfo.RetiredReason = reason

	// Note: Insurance eligibility should be tracked via DeviceInsurance relationship
	// DeviceCompliance tracks regulatory compliance, not insurance eligibility
	// Insurance eligibility is determined by business logic, not stored in compliance

	return s.deviceService.repo.Update(ctx, dev)
}

// Helper methods
func (s *DeviceStateService) calculateRiskLevel(score float64) device.RiskLevel {
	switch {
	case score < 25:
		return device.RiskLevelLow
	case score < 50:
		return device.RiskLevelMedium
	case score < 75:
		return device.RiskLevelHigh
	default:
		return device.RiskLevelCritical
	}
}

func (s *DeviceStateService) calculateBaseRiskScore(dev *models.Device) float64 {
	score := 0.0

	// Age factor based on purchase date
	if dev.DeviceFinancial.PurchaseDate != nil {
		age := time.Since(*dev.DeviceFinancial.PurchaseDate).Hours() / 24 / 365
		score += age * 5 // 5 points per year
	}

	// Condition factor
	switch string(dev.DevicePhysicalCondition.Condition) {
	case "poor":
		score += 30
	case "fair":
		score += 20
	case "good":
		score += 10
	case "excellent":
		score += 5
	}

	// Grade factor
	switch string(dev.DevicePhysicalCondition.Grade) {
	case "F":
		score += 30
	case "D":
		score += 20
	case "C":
		score += 15
	case "B":
		score += 10
	case "A":
		score += 5
	}

	// Claims history
	score += float64(len(dev.Claims)) * 10

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

func (s *DeviceStateService) checkInsurability(dev *models.Device) bool {
	// Device must be active
	if dev.DeviceStatusInfo.Status != device.DeviceStatusActive {
		return false
	}

	// Not stolen (lost devices already have status != "active")
	if dev.DeviceStatusInfo.IsStolen {
		return false
	}

	// Not blacklisted
	if dev.DeviceRiskAssessment.BlacklistStatus == device.BlacklistStatusBlocked {
		return false
	}

	// Condition check
	if string(dev.DevicePhysicalCondition.Grade) == "F" {
		return false
	}

	// Risk score check
	if dev.DeviceRiskAssessment.RiskScore >= 80 {
		return false
	}

	return true
}
