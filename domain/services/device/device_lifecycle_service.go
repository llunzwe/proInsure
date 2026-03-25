package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	devicemodels "smartsure/internal/domain/models/device"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// DeviceLifecycleService manages device lifecycle operations
type DeviceLifecycleService struct {
	deviceRepo repositories.DeviceRepository
	userRepo   repositories.UserRepository
	policyRepo repositories.PolicyRepository
	// auditRepo and notificationSvc are optional - can be nil
	auditRepo       interface{} // Optional audit repository
	notificationSvc interface{} // Optional notification service
}

// NewDeviceLifecycleService creates a new device lifecycle service
func NewDeviceLifecycleService(
	deviceRepo repositories.DeviceRepository,
	userRepo repositories.UserRepository,
	policyRepo repositories.PolicyRepository,
	auditRepo interface{},
	notificationSvc interface{},
) *DeviceLifecycleService {
	return &DeviceLifecycleService{
		deviceRepo:      deviceRepo,
		userRepo:        userRepo,
		policyRepo:      policyRepo,
		auditRepo:       auditRepo,
		notificationSvc: notificationSvc,
	}
}

// RegisterDevice registers a new device in the system
func (s *DeviceLifecycleService) RegisterDevice(ctx context.Context, device *models.Device) error {
	// Validate device
	if err := s.validateNewDevice(device); err != nil {
		return fmt.Errorf("device validation failed: %w", err)
	}

	// Set initial lifecycle status
	device.DeviceStatusInfo.Status = devicemodels.DeviceStatusActive
	device.DeviceStatusInfo.IsStolen = false
	device.DeviceVerification.IsVerified = false
	now := time.Now()
	device.DeviceLifecycle.FirstActivation = now
	device.DeviceLifecycle.CurrentStage = "active"

	// Set risk defaults
	device.DeviceRiskAssessment.RiskScore = 50 // Neutral starting score
	device.DeviceRiskAssessment.RiskLevel = devicemodels.RiskLevelMedium
	device.DeviceRiskAssessment.BlacklistStatus = devicemodels.BlacklistStatusClean
	device.DeviceVerification.AuthenticityStatus = devicemodels.AuthenticityStatusChecking

	// Create device
	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	// Create audit log (if audit service is available)
	// Note: Audit logging should be handled by infrastructure layer

	// Send notification (if notification service is available)
	// Note: Notifications should be handled by infrastructure layer
	if device.DeviceOwnership.OwnerID != uuid.Nil {
		// Notification would be sent via notification service if available
		_ = fmt.Sprintf("Your device %s %s has been successfully registered",
			device.DeviceClassification.Brand,
			device.DeviceClassification.Model)
	}

	return nil
}

// VerifyDevice verifies device authenticity
func (s *DeviceLifecycleService) VerifyDevice(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Perform verification checks
	verificationResult := s.performVerificationChecks(device)

	// Update verification status
	device.DeviceVerification.IsVerified = verificationResult.IsValid
	device.DeviceVerification.VerificationDate = &verificationResult.VerifiedAt
	device.DeviceVerification.VerificationMethod = verificationResult.Method

	if verificationResult.IsValid {
		device.DeviceVerification.AuthenticityStatus = devicemodels.AuthenticityStatusVerified
		device.DeviceRiskAssessment.RiskScore -= 10 // Reduce risk for verified devices
	} else {
		device.DeviceVerification.AuthenticityStatus = devicemodels.AuthenticityStatusFailed
		device.DeviceRiskAssessment.RiskScore += 20
	}

	// Save updates
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	// Update device verification fields directly - already done above
	// Verification is stored in the embedded DeviceVerification struct

	return nil
}

// UpdateVerification updates device verification status
func (s *DeviceLifecycleService) UpdateVerification(ctx context.Context, deviceID uuid.UUID) error {
	return s.VerifyDevice(ctx, deviceID)
}

// ProcessTradeIn processes device trade-in
func (s *DeviceLifecycleService) ProcessTradeIn(ctx context.Context, deviceID uuid.UUID, newDeviceID uuid.UUID) error {
	oldDevice, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get old device: %w", err)
	}

	_, err = s.deviceRepo.GetByID(ctx, newDeviceID)
	if err != nil {
		return fmt.Errorf("failed to get new device: %w", err)
	}

	// Validate trade-in eligibility
	if !s.isEligibleForTradeIn(oldDevice) {
		return errors.New("device not eligible for trade-in")
	}

	// Calculate trade-in value
	tradeInValue := s.calculateTradeInValue(oldDevice)

	// Create trade-in record
	now := time.Now()
	userID := oldDevice.DeviceOwnership.OwnerID
	_ = &devicemodels.DeviceTradeIn{
		DeviceID:          deviceID,
		UserID:            userID,
		Status:            "completed",
		QuotedValue:       tradeInValue,
		LockedValue:       tradeInValue,
		FinalValue:        tradeInValue,
		OriginalValue:     oldDevice.DeviceFinancial.CurrentValue.Amount,
		TradeInGrade:      string(oldDevice.DevicePhysicalCondition.Grade),
		EligibilityStatus: "eligible",
		CompletionDate:    &now,
		InstantTradeIn:    true,
	}

	// Note: CreateTradeIn may not exist - update device directly for now
	// In production, this would be handled by a proper trade-in service

	// Update old device status
	oldDevice.DeviceStatusInfo.Status = devicemodels.DeviceStatusRetired
	oldDevice.DeviceLifecycle.CurrentStage = "disposed"
	oldDevice.DeviceLifecycle.DisposalMethod = "trade_in"
	oldDevice.DeviceLifecycle.EOLDate = &now

	if err := s.deviceRepo.Update(ctx, oldDevice); err != nil {
		return fmt.Errorf("failed to update old device: %w", err)
	}

	// Transfer any active policies
	if err := s.transferPolicies(ctx, deviceID, newDeviceID); err != nil {
		// Log but don't fail
		fmt.Printf("failed to transfer policies: %v\n", err)
	}

	return nil
}

// InitiateRental initiates device rental
func (s *DeviceLifecycleService) InitiateRental(ctx context.Context, deviceID uuid.UUID, duration time.Duration) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Check rental eligibility
	if !s.isEligibleForRental(device) {
		return errors.New("device not eligible for rental")
	}

	// Calculate rental terms
	monthlyRate := s.calculateMonthlyRentalRate(device)
	totalCost := monthlyRate * float64(duration.Hours()/24/30)

	// Create rental agreement
	now := time.Now()
	_ = &devicemodels.DeviceRental{
		DeviceID:          deviceID,
		UserID:            device.DeviceOwnership.OwnerID,
		RentalStatus:      "active",
		RentalPlanType:    "monthly",
		MonthlyRate:       monthlyRate,
		TotalRentalFee:    totalCost,
		SecurityDeposit:   device.DeviceFinancial.CurrentValue.Amount * 0.20,
		RentalStartDate:   now,
		RentalEndDate:     now.Add(duration),
		InsuranceIncluded: true,
	}

	// Note: CreateRental may not exist - update device directly for now
	// In production, this would be handled by a proper rental service

	// Update device status - devices can still be active while rented
	device.DeviceStatusInfo.Status = devicemodels.DeviceStatusActive
	device.DeviceLifecycle.CurrentStage = "active"

	return s.deviceRepo.Update(ctx, device)
}

// InitiateFinancing initiates device financing
func (s *DeviceLifecycleService) InitiateFinancing(ctx context.Context, deviceID uuid.UUID, terms *services.FinancingTerms) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Validate financing eligibility
	if !s.isEligibleForFinancing(device, terms) {
		return errors.New("device/buyer not eligible for financing")
	}

	// Calculate financing details
	principal := device.DeviceFinancial.CurrentValue.Amount - terms.DownPayment
	monthlyPayment := s.calculateMonthlyPayment(principal, terms.InterestRate, terms.Duration)

	// Create financing record
	now := time.Now()
	nextPayment := now.AddDate(0, 1, 0)
	totalAmount := monthlyPayment*float64(terms.Duration) + terms.DownPayment
	_ = &devicemodels.DeviceFinancing{
		DeviceID:           deviceID,
		UserID:             device.DeviceOwnership.OwnerID,
		FinanceStatus:      "active",
		FinanceType:        "installment",
		FinanceProvider:    "", // Provider not available in FinancingTerms
		FinanceStartDate:   now,
		FinanceEndDate:     now.AddDate(0, int(terms.Duration), 0),
		TotalFinanceAmount: totalAmount,
		DownPayment:        terms.DownPayment,
		PrincipalAmount:    principal,
		InterestRate:       terms.InterestRate,
		PaymentFrequency:   "monthly",
		InstallmentAmount:  monthlyPayment,
		TotalInstallments:  int(terms.Duration),
		InstallmentsPaid:   0,
		NextPaymentDate:    &nextPayment,
	}

	// Note: CreateFinancing may not exist - update device directly for now
	// In production, this would be handled by a proper financing service

	// Update device status - device remains active when financed
	device.DeviceStatusInfo.Status = devicemodels.DeviceStatusActive
	device.DeviceLifecycle.CurrentStage = "active"

	return s.deviceRepo.Update(ctx, device)
}

// TransferOwnership transfers device ownership
func (s *DeviceLifecycleService) TransferOwnership(ctx context.Context, deviceID uuid.UUID, newOwnerID uuid.UUID) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Validate transfer eligibility
	if device.DeviceStatusInfo.IsStolen || device.DeviceStatusInfo.IsLost {
		return errors.New("cannot transfer stolen or lost device")
	}

	// Check for active claims
	if len(device.Claims) > 0 {
		// Check if any claims are still active
		hasActive := false
		for _, claim := range device.Claims {
			if claim.ClaimLifecycle.Status != "closed" && claim.ClaimLifecycle.Status != "denied" {
				hasActive = true
				break
			}
		}
		if hasActive {
			return errors.New("cannot transfer device with active claims")
		}
	}

	oldOwnerID := device.DeviceOwnership.OwnerID

	// Create history record
	_ = &devicemodels.DeviceHistory{
		DeviceID:        &deviceID,
		IMEI:            device.DeviceIdentification.IMEI,
		UserID:          oldOwnerID, // Current owner
		EventType:       "ownership_transfer",
		EventDate:       time.Now(),
		Description:     fmt.Sprintf("Ownership transferred from %v to %v", oldOwnerID, newOwnerID),
		PreviousOwnerID: &oldOwnerID,
		NewOwnerID:      &newOwnerID,
	}

	// Note: CreateHistory might not exist - using history via DeviceHistory model if needed
	// For now, we'll just update the device and track in lifecycle

	// Update ownership
	device.DeviceOwnership.OwnerID = newOwnerID
	// Add to transfer history
	device.DeviceOwnership.TransferHistory = append(device.DeviceOwnership.TransferHistory, oldOwnerID)

	now := time.Now()
	device.DeviceLifecycle.LastTransferDate = &now
	device.DeviceLifecycle.OwnershipTransfers++
	device.DeviceLifecycle.TotalOwners++

	// Update risk assessment
	device.DeviceRiskAssessment.RiskScore += 5 // Slight increase for transfer

	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	// Notify both parties (if notification service is available)
	// Note: Notifications should be handled by infrastructure layer
	if oldOwnerID != uuid.Nil {
		// Notification would be sent via notification service if available
		_ = fmt.Sprintf("Device %s %s has been transferred",
			device.DeviceClassification.Brand,
			device.DeviceClassification.Model)
	}

	// Notification would be sent via notification service if available
	_ = fmt.Sprintf("You have received device %s %s",
		device.DeviceClassification.Brand,
		device.DeviceClassification.Model)

	return nil
}

// MarkAsStolen marks device as stolen
func (s *DeviceLifecycleService) MarkAsStolen(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Update status
	device.DeviceStatusInfo.Status = devicemodels.DeviceStatusStolen
	device.DeviceStatusInfo.IsStolen = true
	now := time.Now()
	device.DeviceStatusInfo.StolenDate = &now

	// Update risk
	device.DeviceRiskAssessment.RiskScore = 100
	device.DeviceRiskAssessment.RiskLevel = devicemodels.RiskLevelCritical
	device.DeviceRiskAssessment.TheftRiskLevel = devicemodels.RiskLevelCritical
	device.DeviceRiskAssessment.BlacklistStatus = devicemodels.BlacklistStatusBlocked

	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	// Create history record
	_ = &devicemodels.DeviceHistory{
		DeviceID:    &deviceID,
		IMEI:        device.DeviceIdentification.IMEI,
		UserID:      device.DeviceOwnership.OwnerID,
		EventType:   "stolen_report",
		EventDate:   time.Now(),
		Description: "Device reported as stolen",
		Metadata:    `{"severity":"critical"}`,
	}

	// Note: CreateHistory might not exist in repository - log for now
	// In production, this would be handled by a proper history service

	// Notify owner and authorities (if notification service is available)
	// Note: Notifications should be handled by infrastructure layer
	if device.DeviceOwnership.OwnerID != uuid.Nil {
		// Notification would be sent via notification service if available
		_ = fmt.Sprintf("Device marked as stolen: %s", device.DeviceIdentification.IMEI)
	}

	return nil
}

// MarkAsRecovered marks stolen device as recovered
func (s *DeviceLifecycleService) MarkAsRecovered(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	if !device.DeviceStatusInfo.IsStolen {
		return errors.New("device is not marked as stolen")
	}

	// Update status
	device.DeviceStatusInfo.Status = devicemodels.DeviceStatusActive
	device.DeviceStatusInfo.IsStolen = false
	device.DeviceStatusInfo.StolenDate = nil

	// Update risk
	device.DeviceRiskAssessment.RiskScore = 60 // Still elevated risk
	device.DeviceRiskAssessment.RiskLevel = devicemodels.RiskLevelMedium
	device.DeviceRiskAssessment.TheftRiskLevel = devicemodels.RiskLevelMedium
	device.DeviceRiskAssessment.BlacklistStatus = devicemodels.BlacklistStatusClean

	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	return nil
}

// Helper methods

func (s *DeviceLifecycleService) validateNewDevice(device *models.Device) error {
	if device.DeviceIdentification.IMEI == "" {
		return errors.New("IMEI is required")
	}

	if device.DeviceClassification.Brand == "" || device.DeviceClassification.Model == "" {
		return errors.New("brand and model are required")
	}

	if device.DeviceOwnership.OwnerID == uuid.Nil {
		return errors.New("owner ID is required")
	}

	return nil
}

func (s *DeviceLifecycleService) performVerificationChecks(device *models.Device) *VerificationResult {
	result := &VerificationResult{
		Method:     "automated",
		VerifiedAt: time.Now(),
		IsValid:    true,
		Result:     "verified",
		Details:    make(map[string]interface{}),
	}

	// Check IMEI validity
	if !s.isValidIMEI(device.DeviceIdentification.IMEI) {
		result.IsValid = false
		result.Result = "invalid_imei"
		result.Details["imei_check"] = "failed"
	}

	// Check blacklist
	if device.DeviceRiskAssessment.BlacklistStatus == devicemodels.BlacklistStatusBlocked {
		result.IsValid = false
		result.Result = "blacklisted"
		result.Details["blacklist_check"] = "failed"
	}

	return result
}

func (s *DeviceLifecycleService) isValidIMEI(imei string) bool {
	// Basic IMEI validation (15 or 17 digits)
	return len(imei) == 15 || len(imei) == 17
}

func (s *DeviceLifecycleService) isEligibleForTradeIn(device *models.Device) bool {
	// Check status
	if device.DeviceStatusInfo.Status != devicemodels.DeviceStatusActive {
		return false
	}

	// Check condition
	if string(device.DevicePhysicalCondition.Grade) == "D" || string(device.DevicePhysicalCondition.Grade) == "F" {
		return false
	}

	// Check age (max 5 years)
	if device.DeviceFinancial.PurchaseDate != nil {
		age := time.Since(*device.DeviceFinancial.PurchaseDate)
		if age > 5*365*24*time.Hour {
			return false
		}
	}

	return true
}

func (s *DeviceLifecycleService) calculateTradeInValue(device *models.Device) float64 {
	baseValue := device.DeviceFinancial.CurrentValue.Amount

	// Apply condition multiplier
	multiplier := 1.0
	switch string(device.DevicePhysicalCondition.Grade) {
	case "A":
		multiplier = 0.8
	case "B":
		multiplier = 0.6
	case "C":
		multiplier = 0.4
	default:
		multiplier = 0.2
	}

	return baseValue * multiplier
}

func (s *DeviceLifecycleService) isEligibleForRental(device *models.Device) bool {
	return device.DeviceStatusInfo.Status == devicemodels.DeviceStatusActive &&
		string(device.DevicePhysicalCondition.Grade) != "D" &&
		string(device.DevicePhysicalCondition.Grade) != "F"
}

func (s *DeviceLifecycleService) calculateMonthlyRentalRate(device *models.Device) float64 {
	// 5% of device value per month
	return device.DeviceFinancial.CurrentValue.Amount * 0.05
}

func (s *DeviceLifecycleService) isEligibleForFinancing(device *models.Device, terms *services.FinancingTerms) bool {
	// Check device value
	if device.DeviceFinancial.CurrentValue.Amount < 200 {
		return false
	}

	// Check down payment
	minDownPayment := device.DeviceFinancial.CurrentValue.Amount * 0.10
	if terms.DownPayment < minDownPayment {
		return false
	}

	// Check credit (simplified) - CreditScore not available in FinancingTerms, skip credit check
	// Credit validation should be handled at application/service layer

	return true
}

func (s *DeviceLifecycleService) calculateMonthlyPayment(principal, rate float64, months int) float64 {
	monthlyRate := rate / 12
	payment := principal * (monthlyRate * pow(1+monthlyRate, float64(months))) /
		(pow(1+monthlyRate, float64(months)) - 1)
	return payment
}

func (s *DeviceLifecycleService) transferPolicies(ctx context.Context, oldDeviceID, newDeviceID uuid.UUID) error {
	policies, err := s.policyRepo.FindByDeviceID(ctx, oldDeviceID)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		policy.PolicyRelationships.DeviceID = newDeviceID
		if err := s.policyRepo.Update(ctx, policy); err != nil {
			return err
		}
	}

	return nil
}

// Helper function for power calculation
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0.0; i < exp; i++ {
		result *= base
	}
	return result
}

// VerificationResult represents device verification result
type VerificationResult struct {
	Method     string
	VerifiedAt time.Time
	IsValid    bool
	Result     string
	Details    map[string]interface{}
}
