package device

import (
	"context"

	"github.com/google/uuid"

	devicemodels "smartsure/internal/domain/models/device"
	"smartsure/internal/domain/ports/repositories"
)

// ClaimsEligibilityService handles device claims eligibility checks
type ClaimsEligibilityService struct {
	deviceRepo repositories.DeviceRepository
	policyRepo repositories.PolicyRepository
	claimRepo  repositories.ClaimRepository
}

// NewClaimsEligibilityService creates a new claims eligibility service
func NewClaimsEligibilityService(
	deviceRepo repositories.DeviceRepository,
	policyRepo repositories.PolicyRepository,
	claimRepo repositories.ClaimRepository,
) *ClaimsEligibilityService {
	return &ClaimsEligibilityService{
		deviceRepo: deviceRepo,
		policyRepo: policyRepo,
		claimRepo:  claimRepo,
	}
}

// CheckClaimEligibility checks if a device is eligible for a claim
func (s *ClaimsEligibilityService) CheckClaimEligibility(ctx context.Context, deviceID uuid.UUID, claimType string) (bool, string, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return false, "", err
	}

	// Check if device has active policy
	policy, err := s.policyRepo.FindActiveByDeviceID(ctx, deviceID)
	if err != nil {
		return false, "", err
	}

	if policy == nil {
		return false, "Device has no active insurance policy", nil
	}

	// Check device status
	if device.DeviceStatusInfo.Status != devicemodels.DeviceStatusActive && !device.DeviceStatusInfo.IsStolen && !device.DeviceStatusInfo.IsLost {
		return false, "Device is not in a valid state for claims", nil
	}

	// Check claim history
	claims, err := s.claimRepo.GetByDeviceID(ctx, deviceID, 100, 0)
	if err != nil {
		return false, "", err
	}

	// Check if claim limit reached
	activeClaims := 0
	for _, claim := range claims {
		if claim.ClaimLifecycle.Status == "pending" || claim.ClaimLifecycle.Status == "approved" {
			activeClaims++
		}
	}

	if activeClaims > 0 {
		return false, "Device has active claims pending resolution", nil
	}

	// Check claim frequency limits (business rule: max 3 claims per year)
	recentClaims := 0
	for _, claim := range claims {
		// Count claims from last 12 months
		// Simplified check - in production, would check actual dates
		if claim.ClaimLifecycle.Status == "approved" {
			recentClaims++
		}
	}

	if recentClaims >= 3 {
		return false, "Claim limit reached (maximum 3 claims per year)", nil
	}

	return true, "Device is eligible for claim", nil
}

// ValidateClaimRequest validates a claim request before processing
func (s *ClaimsEligibilityService) ValidateClaimRequest(ctx context.Context, deviceID uuid.UUID, claimType string, claimAmount float64) (bool, []string, error) {
	eligible, reason, err := s.CheckClaimEligibility(ctx, deviceID, claimType)
	if err != nil {
		return false, nil, err
	}

	var errors []string
	if !eligible {
		errors = append(errors, reason)
	}

	// Additional validations
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return false, nil, err
	}

	// Validate claim amount doesn't exceed device value
	if claimAmount > device.DeviceFinancial.CurrentValue.Amount {
		errors = append(errors, "Claim amount exceeds device current value")
	}

	// Validate claim type
	validClaimTypes := map[string]bool{
		"theft":       true,
		"loss":        true,
		"damage":      true,
		"malfunction": true,
	}

	if !validClaimTypes[claimType] {
		errors = append(errors, "Invalid claim type")
	}

	return len(errors) == 0, errors, nil
}

// GetClaimEligibilityDetails returns detailed eligibility information
func (s *ClaimsEligibilityService) GetClaimEligibilityDetails(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	details := make(map[string]interface{})

	// Check policy
	policy, err := s.policyRepo.FindActiveByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	details["has_active_policy"] = policy != nil
	if policy != nil {
		details["policy_id"] = policy.ID
		details["policy_type"] = string(policy.PolicyClassification.PolicyType)
	}

	// Check device status
	details["device_status"] = string(device.DeviceStatusInfo.Status)
	details["is_stolen"] = device.DeviceStatusInfo.IsStolen
	details["is_lost"] = device.DeviceStatusInfo.IsLost

	// Check claim history
	claims, err := s.claimRepo.GetByDeviceID(ctx, deviceID, 100, 0)
	if err != nil {
		return nil, err
	}

	details["total_claims"] = len(claims)
	details["active_claims"] = 0
	details["approved_claims"] = 0

	for _, claim := range claims {
		if claim.ClaimLifecycle.Status == "pending" || claim.ClaimLifecycle.Status == "approved" {
			details["active_claims"] = details["active_claims"].(int) + 1
		}
		if claim.ClaimLifecycle.Status == "approved" {
			details["approved_claims"] = details["approved_claims"].(int) + 1
		}
	}

	// Eligibility result
	eligible, reason, _ := s.CheckClaimEligibility(ctx, deviceID, "damage")
	details["eligible"] = eligible
	details["eligibility_reason"] = reason

	return details, nil
}
