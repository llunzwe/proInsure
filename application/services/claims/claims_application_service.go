package device

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/device"
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
		eligibilityService: domainServices.NewClaimsEligibilityService(),
	}
}

// ProcessClaimEligibility checks eligibility and provides quote for a claim
func (s *ClaimsApplicationService) ProcessClaimEligibility(ctx context.Context, deviceID string, claimType string) (*ClaimEligibilityResult, error) {
	// 1. Load device with related data from database
	var device device.Device
	if err := s.db.WithContext(ctx).
		Preload("Claims").
		Preload("Repairs").
		First(&device, "id = ?", deviceID).Error; err != nil {
		return nil, fmt.Errorf("failed to load device: %w", err)
	}

	// 2. Get claim history count (infrastructure)
	claimHistory := s.getClaimHistoryCount(ctx, deviceID)

	// 3. Check for active claims (infrastructure)
	hasActiveClaim := s.hasActiveClaim(ctx, deviceID)
	if hasActiveClaim {
		return &ClaimEligibilityResult{
			DeviceID:   deviceID,
			ClaimType:  claimType,
			IsEligible: false,
			Reason:     "Device has an active claim in progress",
		}, nil
	}

	// 4. Use domain service for eligibility check
	isEligible := s.eligibilityService.CheckClaimEligibility(&device, claimType)

	result := &ClaimEligibilityResult{
		DeviceID:   deviceID,
		ClaimType:  claimType,
		IsEligible: isEligible,
		CheckedAt:  time.Now(),
	}

	if isEligible {
		// 5. Calculate financial details using domain service
		result.Deductible = s.eligibilityService.CalculateDeductible(&device, claimType, claimHistory)
		result.CoverageLimit = s.eligibilityService.CalculateCoverageLimit(&device)
		result.EstimatedPayout = s.eligibilityService.CalculateClaimPayout(&device, claimType, claimHistory)
		result.ClaimFrequency = s.calculateClaimFrequency(ctx, deviceID)
	} else {
		result.Reason = s.determineIneligibilityReason(&device, claimType)
	}

	// 6. Save eligibility check result for audit trail
	s.saveEligibilityCheck(ctx, result)

	return result, nil
}

// FileClaim initiates a new claim for a device
func (s *ClaimsApplicationService) FileClaim(ctx context.Context, request *FileClaimRequest) (*models.Claim, error) {
	// 1. Check eligibility first
	eligibility, err := s.ProcessClaimEligibility(ctx, request.DeviceID, request.ClaimType)
	if err != nil {
		return nil, fmt.Errorf("failed to check eligibility: %w", err)
	}

	if !eligibility.IsEligible {
		return nil, fmt.Errorf("device not eligible for claim: %s", eligibility.Reason)
	}

	// 2. Create claim record
	claim := &models.Claim{
		ID:             uuid.New(),
		DeviceID:       uuid.MustParse(request.DeviceID),
		Type:           request.ClaimType,
		Status:         "pending",
		ClaimAmount:    eligibility.EstimatedPayout,
		ApprovedAmount: 0,
		Description:    request.Description,
		IncidentDate:   &request.IncidentDate,
	}

	// 3. Save claim to database
	if err := s.db.WithContext(ctx).Create(claim).Error; err != nil {
		return nil, fmt.Errorf("failed to create claim: %w", err)
	}

	// 4. Create claim processing workflow
	s.initiateClaimWorkflow(ctx, claim)

	return claim, nil
}

// GetClaimHistory retrieves claim history for a device (infrastructure)
func (s *ClaimsApplicationService) GetClaimHistory(ctx context.Context, deviceID string) ([]ClaimSummary, error) {
	var claims []models.Claim

	if err := s.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		Find(&claims).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve claims: %w", err)
	}

	summaries := make([]ClaimSummary, len(claims))
	for i, claim := range claims {
		summaries[i] = ClaimSummary{
			ClaimID:      claim.ID.String(),
			ClaimType:    claim.Type,
			Status:       claim.Status,
			Amount:       claim.ApprovedAmount,
			FiledDate:    claim.CreatedAt,
			ResolvedDate: s.getResolvedDate(&claim),
		}
	}

	return summaries, nil
}

// UpdateClaimStatus updates the status of a claim (infrastructure)
func (s *ClaimsApplicationService) UpdateClaimStatus(ctx context.Context, claimID string, status string, notes string) error {
	var claim models.Claim
	if err := s.db.WithContext(ctx).First(&claim, "id = ?", claimID).Error; err != nil {
		return fmt.Errorf("claim not found: %w", err)
	}

	// Update status
	claim.Status = status
	claim.InvestigationNotes = notes

	// Set appropriate timestamps
	now := time.Now()
	switch status {
	case "approved", "rejected":
		claim.ResolvedDate = &now
	case "paid":
		claim.PayoutDate = &now
	}

	return s.db.WithContext(ctx).Save(&claim).Error
}

// Private helper methods

func (s *ClaimsApplicationService) getClaimHistoryCount(ctx context.Context, deviceID string) int {
	var count int64
	s.db.WithContext(ctx).
		Model(&models.Claim{}).
		Where("device_id = ? AND status != ?", deviceID, "draft").
		Count(&count)
	return int(count)
}

func (s *ClaimsApplicationService) hasActiveClaim(ctx context.Context, deviceID string) bool {
	var count int64
	s.db.WithContext(ctx).
		Model(&models.Claim{}).
		Where("device_id = ? AND status IN ?", deviceID, []string{"pending", "processing", "under_review"}).
		Count(&count)
	return count > 0
}

func (s *ClaimsApplicationService) calculateClaimFrequency(ctx context.Context, deviceID string) float64 {
	var device device.Device
	if err := s.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error; err != nil {
		return 0
	}

	var claimCount int64
	s.db.WithContext(ctx).
		Model(&models.Claim{}).
		Where("device_id = ? AND status NOT IN ?", deviceID, []string{"draft", "cancelled"}).
		Count(&claimCount)

	age := time.Since(device.RegistrationDate).Hours() / (24 * 365)
	if age < 0.5 {
		age = 0.5
	}

	return float64(claimCount) / age
}

func (s *ClaimsApplicationService) determineIneligibilityReason(Device *models.Device, claimType string) string {
	if device.Status != "active" {
		return "Device is not active"
	}
	if device.IsStolen {
		return "Device is reported as stolen"
	}
	if !device.IsVerified {
		return "Device is not verified"
	}
	if device.BlacklistStatus == "blocked" {
		return "Device is blacklisted"
	}

	// Type-specific reasons
	switch claimType {
	case "theft":
		if !device.FindMyDeviceEnabled {
			return "Find My Device must be enabled for theft claims"
		}
	case "water_damage":
		if device.WaterDamageIndicator != "white" {
			return "Water damage already detected"
		}
	case "battery_replacement":
		if device.BatteryHealth >= 80 {
			return "Battery health is above 80%"
		}
	}

	return "Device does not meet eligibility criteria for " + claimType
}

func (s *ClaimsApplicationService) getResolvedDate(claim *models.Claim) *time.Time {
	if claim.Status == "approved" || claim.Status == "rejected" || claim.Status == "paid" {
		return &claim.UpdatedAt
	}
	return nil
}

func (s *ClaimsApplicationService) saveEligibilityCheck(ctx context.Context, result *ClaimEligibilityResult) {
	// Save to audit table in production
	s.db.WithContext(ctx).Table("claim_eligibility_checks").Create(result)
}

func (s *ClaimsApplicationService) initiateClaimWorkflow(ctx context.Context, claim *models.Claim) {
	// In production, this would trigger workflow engine
	// For now, just log
	fmt.Printf("Claim workflow initiated for claim %s\n", claim.ID)
}

// Data structures

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

type FileClaimRequest struct {
	DeviceID     string    `json:"device_id"`
	ClaimType    string    `json:"claim_type"`
	Description  string    `json:"description"`
	IncidentDate time.Time `json:"incident_date"`
}

type ClaimSummary struct {
	ClaimID      string     `json:"claim_id"`
	ClaimType    string     `json:"claim_type"`
	Status       string     `json:"status"`
	Amount       float64    `json:"amount"`
	FiledDate    time.Time  `json:"filed_date"`
	ResolvedDate *time.Time `json:"resolved_date,omitempty"`
}
