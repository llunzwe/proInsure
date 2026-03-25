package device

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	
	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// DeviceInsuranceService handles device insurance operations
type DeviceInsuranceService struct {
	deviceRepo repositories.DeviceRepository
	policyRepo repositories.PolicyRepository
	claimRepo  repositories.ClaimRepository
}

// NewDeviceInsuranceService creates a new device insurance service
func NewDeviceInsuranceService(
	deviceRepo repositories.DeviceRepository,
	policyRepo repositories.PolicyRepository,
	claimRepo repositories.ClaimRepository,
) *DeviceInsuranceService {
	return &DeviceInsuranceService{
		deviceRepo: deviceRepo,
		policyRepo: policyRepo,
		claimRepo:  claimRepo,
	}
}

// CheckInsuranceEligibility checks if a device is eligible for insurance
func (s *DeviceInsuranceService) CheckInsuranceEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return false, "", fmt.Errorf("failed to get device: %w", err)
	}

	// Check device status
	if device.DeviceLifecycle.Status == "stolen" || device.DeviceLifecycle.Status == "lost" {
		return false, "Device is reported as stolen or lost", nil
	}

	// Check device condition
	conditionScore := s.calculateConditionScore(device)
	if conditionScore < 40 {
		return false, "Device condition is too poor for insurance", nil
	}

	// Check device age
	age := s.calculateDeviceAge(device)
	if age > 5*365 { // 5 years in days
		return false, "Device is too old (>5 years) for insurance", nil
	}

	// Check device value
	currentValue := s.calculateCurrentValue(device)
	if currentValue.LessThan(decimal.NewFromFloat(100)) {
		return false, "Device value is below minimum threshold ($100)", nil
	}

	// Check existing claims
	activeClaims, err := s.claimRepo.GetActiveClaimsByDevice(ctx, deviceID)
	if err == nil && len(activeClaims) > 0 {
		return false, "Device has active claims pending", nil
	}

	// Check blacklist status
	if device.DeviceRiskAssessment.BlacklistStatus == "blacklisted" {
		return false, "Device is blacklisted", nil
	}

	return true, "Device is eligible for insurance", nil
}

// CalculatePremium calculates insurance premium for a device
func (s *DeviceInsuranceService) CalculatePremium(ctx context.Context, deviceID uuid.UUID) (*services.Premium, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// Base premium calculation
	currentValue := s.calculateCurrentValue(device)
	basePremium := currentValue.Mul(decimal.NewFromFloat(0.08)) // 8% of device value

	var adjustments []services.Adjustment

	// Risk adjustment
	riskScore := device.DeviceRiskAssessment.RiskScore
	if riskScore > 70 {
		adjustment := basePremium.Mul(decimal.NewFromFloat(0.25)) // 25% increase
		adjustments = append(adjustments, services.Adjustment{
			Type:   "high_risk",
			Amount: adjustment.InexactFloat64(),
			Reason: "High risk device",
		})
		basePremium = basePremium.Add(adjustment)
	} else if riskScore < 30 {
		adjustment := basePremium.Mul(decimal.NewFromFloat(0.10)) // 10% discount
		adjustments = append(adjustments, services.Adjustment{
			Type:   "low_risk",
			Amount: -adjustment.InexactFloat64(),
			Reason: "Low risk device",
		})
		basePremium = basePremium.Sub(adjustment)
	}

	// Condition adjustment
	conditionScore := s.calculateConditionScore(device)
	if conditionScore >= 80 {
		adjustment := basePremium.Mul(decimal.NewFromFloat(0.05)) // 5% discount
		adjustments = append(adjustments, services.Adjustment{
			Type:   "excellent_condition",
			Amount: -adjustment.InexactFloat64(),
			Reason: "Excellent device condition",
		})
		basePremium = basePremium.Sub(adjustment)
	}

	// Corporate discount
	if device.DeviceOwnership.CorporateAccountID != nil {
		adjustment := basePremium.Mul(decimal.NewFromFloat(0.15)) // 15% corporate discount
		adjustments = append(adjustments, services.Adjustment{
			Type:   "corporate_discount",
			Amount: -adjustment.InexactFloat64(),
			Reason: "Corporate account discount",
		})
		basePremium = basePremium.Sub(adjustment)
	}

	return &services.Premium{
		DeviceID:    deviceID,
		BaseAmount:  currentValue.InexactFloat64() * 0.08,
		Adjustments: adjustments,
		FinalAmount: basePremium.InexactFloat64(),
		ValidUntil:  time.Now().Add(30 * 24 * time.Hour), // Valid for 30 days
	}, nil
}

// CalculateDepreciation calculates device depreciation
func (s *DeviceInsuranceService) CalculateDepreciation(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, fmt.Errorf("failed to get device: %w", err)
	}

	if device.DeviceFinancial.PurchaseDate == nil {
		return 0, errors.New("purchase date not available")
	}

	age := s.calculateDeviceAge(device)
	purchasePrice := decimal.NewFromFloat(device.DeviceFinancial.PurchasePrice)

	// Depreciation formula: 20% first year, 15% second year, 10% subsequent years
	var depreciationRate decimal.Decimal

	switch {
	case age <= 365:
		depreciationRate = decimal.NewFromFloat(0.20)
	case age <= 730:
		depreciationRate = decimal.NewFromFloat(0.35) // 20% + 15%
	default:
		years := float64(age) / 365
		depreciationRate = decimal.NewFromFloat(0.35 + (years-2)*0.10)
		if depreciationRate.GreaterThan(decimal.NewFromFloat(0.85)) {
			depreciationRate = decimal.NewFromFloat(0.85) // Max 85% depreciation
		}
	}

	depreciation := purchasePrice.Mul(depreciationRate)
	return depreciation.InexactFloat64(), nil
}

// GetCurrentInsurableValue gets the current insurable value of a device
func (s *DeviceInsuranceService) GetCurrentInsurableValue(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, fmt.Errorf("failed to get device: %w", err)
	}

	currentValue := s.calculateCurrentValue(device)

	// Apply condition factor
	conditionScore := s.calculateConditionScore(device)
	conditionFactor := decimal.NewFromFloat(float64(conditionScore) / 100)

	insurableValue := currentValue.Mul(conditionFactor)

	// Apply minimum and maximum limits
	minValue := decimal.NewFromFloat(100)
	maxValue := decimal.NewFromFloat(10000)

	if insurableValue.LessThan(minValue) {
		insurableValue = minValue
	}
	if insurableValue.GreaterThan(maxValue) {
		insurableValue = maxValue
	}

	return insurableValue.InexactFloat64(), nil
}

// ValidateClaimEligibility validates if a device can file a claim
func (s *DeviceInsuranceService) ValidateClaimEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return false, "", fmt.Errorf("failed to get device: %w", err)
	}

	// Check if device has active policy
	activePolicies, err := s.policyRepo.GetActivePoliciesByDevice(ctx, deviceID)
	if err != nil || len(activePolicies) == 0 {
		return false, "No active insurance policy found for device", nil
	}

	// Check claim frequency
	claims, err := s.claimRepo.GetClaimsByDeviceInPeriod(ctx, deviceID, time.Now().AddDate(0, -12, 0), time.Now())
	if err == nil && len(claims) >= 3 {
		return false, "Maximum claim limit (3 per year) reached", nil
	}

	// Check for pending claims
	activeClaims, err := s.claimRepo.GetActiveClaimsByDevice(ctx, deviceID)
	if err == nil && len(activeClaims) > 0 {
		return false, "Device has pending claims that must be resolved first", nil
	}

	// Check waiting period (30 days from policy start)
	policy := activePolicies[0]
	if time.Since(policy.StartDate) < 30*24*time.Hour {
		return false, "Policy is still in waiting period (30 days)", nil
	}

	return true, "Device is eligible to file a claim", nil
}

// Helper methods

func (s *DeviceInsuranceService) calculateConditionScore(device *models.Device) int {
	score := 100

	// Screen condition impact
	switch device.DevicePhysicalCondition.ScreenCondition {
	case "cracked":
		score -= 30
	case "scratched":
		score -= 15
	case "minor_scratches":
		score -= 5
	}

	// Body condition impact
	switch device.DevicePhysicalCondition.BodyCondition {
	case "damaged":
		score -= 25
	case "dented":
		score -= 15
	case "scratched":
		score -= 10
	case "minor_wear":
		score -= 5
	}

	// Grade impact
	switch device.DevicePhysicalCondition.Grade {
	case "A":
		// No reduction
	case "B":
		score -= 10
	case "C":
		score -= 20
	case "D":
		score -= 40
	default:
		score -= 50
	}

	if score < 0 {
		score = 0
	}

	return score
}

func (s *DeviceInsuranceService) calculateDeviceAge(device *models.Device) int {
	if device.DeviceFinancial.PurchaseDate == nil {
		// Use default age of 1 year if purchase date not available
		return 365
	}

	age := time.Since(*device.DeviceFinancial.PurchaseDate)
	return int(age.Hours() / 24) // Return age in days
}

func (s *DeviceInsuranceService) calculateCurrentValue(device *models.Device) decimal.Decimal {
	// Use market value if available
	if device.DeviceFinancial.MarketValue > 0 {
		return decimal.NewFromFloat(device.DeviceFinancial.MarketValue)
	}

	// Otherwise calculate based on purchase price and depreciation
	if device.DeviceFinancial.PurchasePrice > 0 {
		purchasePrice := decimal.NewFromFloat(device.DeviceFinancial.PurchasePrice)
		age := s.calculateDeviceAge(device)

		// Simple depreciation: 20% per year
		years := decimal.NewFromFloat(float64(age) / 365)
		depreciationRate := years.Mul(decimal.NewFromFloat(0.20))
		if depreciationRate.GreaterThan(decimal.NewFromFloat(0.85)) {
			depreciationRate = decimal.NewFromFloat(0.85)
		}

		currentValue := purchasePrice.Mul(decimal.NewFromFloat(1).Sub(depreciationRate))
		return currentValue
	}

	// Default value
	return decimal.NewFromFloat(500)
}
