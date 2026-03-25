package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/device/categories"
	"smartsure/internal/domain/models/device/categories/base"
)

// DeviceCategoryService handles device category-specific operations
type DeviceCategoryService struct {
	categoryService *categories.CategoryService
	repository      DeviceRepository
	eventPublisher  EventPublisher
}

// DeviceRepository interface for device persistence
type DeviceRepository interface {
	Create(ctx context.Context, device *DeviceEntity) error
	Update(ctx context.Context, device *DeviceEntity) error
	FindByID(ctx context.Context, id uuid.UUID) (*DeviceEntity, error)
	FindByIMEI(ctx context.Context, imei string) (*DeviceEntity, error)
	List(ctx context.Context, filters map[string]interface{}) ([]*DeviceEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EventPublisher interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}

// DeviceEntity represents a device in the repository
type DeviceEntity struct {
	ID               uuid.UUID              `json:"id"`
	Category         base.CategoryType      `json:"category"`
	Specifications   json.RawMessage        `json:"specifications"`
	IMEI             string                 `json:"imei"`
	SerialNumber     string                 `json:"serial_number"`
	OwnerID          uuid.UUID              `json:"owner_id"`
	Status           string                 `json:"status"`
	InsuranceProfile *base.InsuranceProfile `json:"insurance_profile,omitempty"`
	RiskAssessment   *base.RiskAssessment   `json:"risk_assessment,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// NewDeviceCategoryService creates a new device category service
func NewDeviceCategoryService(repo DeviceRepository, publisher EventPublisher) *DeviceCategoryService {
	return &DeviceCategoryService{
		categoryService: categories.NewCategoryService(),
		repository:      repo,
		eventPublisher:  publisher,
	}
}

// RegisterDevice registers a new device with category-specific validation
func (s *DeviceCategoryService) RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*DeviceEntity, error) {
	// Check if category is supported
	if !s.categoryService.IsSupported(req.Category) {
		return nil, fmt.Errorf("unsupported device category: %s", req.Category)
	}

	// Create specification from JSON
	spec, err := s.categoryService.CreateSpec(req.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to create specification: %w", err)
	}

	// Populate specification from request data
	if err := spec.FromJSON(req.Specifications); err != nil {
		return nil, fmt.Errorf("failed to parse specifications: %w", err)
	}

	// Validate device specifications
	if err := s.categoryService.ValidateDevice(req.Category, spec); err != nil {
		return nil, fmt.Errorf("device validation failed: %w", err)
	}

	// Check device eligibility for insurance
	eligible, issues, err := s.categoryService.ValidateEligibility(spec)
	if err != nil {
		return nil, fmt.Errorf("eligibility check failed: %w", err)
	}

	if !eligible {
		return nil, fmt.Errorf("device not eligible for insurance: %v", issues)
	}

	// Create device entity
	device := &DeviceEntity{
		ID:             uuid.New(),
		Category:       req.Category,
		Specifications: req.Specifications,
		IMEI:           req.IMEI,
		SerialNumber:   req.SerialNumber,
		OwnerID:        req.OwnerID,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Generate insurance profile
	profile, err := s.generateInsuranceProfile(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to generate insurance profile: %w", err)
	}
	device.InsuranceProfile = profile

	// Perform initial risk assessment
	assessment, err := s.performRiskAssessment(ctx, spec, base.ClaimHistory{})
	if err != nil {
		return nil, fmt.Errorf("failed to perform risk assessment: %w", err)
	}
	device.RiskAssessment = &assessment

	// Save to repository
	if err := s.repository.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to save device: %w", err)
	}

	// Publish device registered event
	event := DeviceRegisteredEvent{
		DeviceID:  device.ID,
		Category:  device.Category,
		OwnerID:   device.OwnerID,
		Timestamp: time.Now(),
	}
	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish event: %v\n", err)
	}

	return device, nil
}

// CalculatePremium calculates insurance premium for a device
func (s *DeviceCategoryService) CalculatePremium(ctx context.Context, deviceID uuid.UUID, coverage base.CoverageLevel) (*PremiumCalculation, error) {
	// Fetch device
	device, err := s.repository.FindByID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	// Create specification from stored data
	spec, err := s.categoryService.CreateSpec(device.Category)
	if err != nil {
		return nil, err
	}

	if err := spec.FromJSON(device.Specifications); err != nil {
		return nil, err
	}

	// Calculate premium
	premium, err := s.categoryService.CalculatePremium(spec, coverage)
	if err != nil {
		return nil, fmt.Errorf("premium calculation failed: %w", err)
	}

	// Get coverage options
	options, err := s.categoryService.GetCoverageOptions(spec)
	if err != nil {
		return nil, err
	}

	// Find selected option
	var selectedOption *base.CoverageOption
	for _, opt := range options {
		if opt.Level == coverage {
			selectedOption = &opt
			break
		}
	}

	return &PremiumCalculation{
		DeviceID:       deviceID,
		Coverage:       coverage,
		MonthlyPremium: premium,
		AnnualPremium:  premium * 12,
		CoverageOption: selectedOption,
		CalculatedAt:   time.Now(),
	}, nil
}

// UpdateRiskAssessment updates risk assessment for a device
func (s *DeviceCategoryService) UpdateRiskAssessment(ctx context.Context, deviceID uuid.UUID, claimHistory base.ClaimHistory) error {
	// Fetch device
	device, err := s.repository.FindByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("device not found: %w", err)
	}

	// Create specification
	spec, err := s.categoryService.CreateSpec(device.Category)
	if err != nil {
		return err
	}

	if err := spec.FromJSON(device.Specifications); err != nil {
		return err
	}

	// Perform risk assessment
	assessment, err := s.performRiskAssessment(ctx, spec, claimHistory)
	if err != nil {
		return fmt.Errorf("risk assessment failed: %w", err)
	}

	// Update device
	device.RiskAssessment = &assessment
	device.UpdatedAt = time.Now()

	if err := s.repository.Update(ctx, device); err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	// Publish risk assessment updated event
	event := RiskAssessmentUpdatedEvent{
		DeviceID:       deviceID,
		RiskScore:      assessment.Score,
		RiskLevel:      assessment.Level,
		RequiresReview: assessment.RequiresReview,
		Timestamp:      time.Now(),
	}
	_ = s.eventPublisher.Publish(ctx, event)

	return nil
}

// GetDeviceValuation gets current valuation for a device
func (s *DeviceCategoryService) GetDeviceValuation(ctx context.Context, deviceID uuid.UUID) (*DeviceValuation, error) {
	// Fetch device
	device, err := s.repository.FindByID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	// Create specification
	spec, err := s.categoryService.CreateSpec(device.Category)
	if err != nil {
		return nil, err
	}

	if err := spec.FromJSON(device.Specifications); err != nil {
		return nil, err
	}

	// Calculate various values
	currentValue, _ := s.categoryService.CalculateDepreciation(spec, device.CreatedAt.Format(time.RFC3339))
	resaleValue, _ := s.categoryService.CalculateResaleValue(spec, "good")

	// Get max coverage from insurance calculator
	factory := base.GetGlobalFactory()
	calculator, _ := factory.CreateInsuranceCalculator(device.Category)
	maxCoverage := calculator.CalculateMaxCoverage(spec)

	return &DeviceValuation{
		DeviceID:         deviceID,
		MarketValue:      spec.GetMarketValue(),
		DepreciatedValue: currentValue,
		ResaleValue:      resaleValue,
		MaxCoverage:      maxCoverage,
		ValuationDate:    time.Now(),
	}, nil
}

// Helper methods

func (s *DeviceCategoryService) generateInsuranceProfile(ctx context.Context, spec base.CategorySpec) (*base.InsuranceProfile, error) {
	// Get coverage options
	options, err := s.categoryService.GetCoverageOptions(spec)
	if err != nil {
		return nil, err
	}

	// Perform initial risk assessment
	assessment, err := s.categoryService.AssessRisk(spec, base.ClaimHistory{})
	if err != nil {
		return nil, err
	}

	// Determine recommended coverage based on value and risk
	recommendedCoverage := base.CoverageStandard
	if spec.GetMarketValue() > 1000 {
		recommendedCoverage = base.CoveragePremium
	}
	if assessment.Score > 70 {
		recommendedCoverage = base.CoverageComprehensive
	}

	profile := &base.InsuranceProfile{
		DeviceID:            uuid.New(),
		CategoryType:        spec.GetCategory(),
		EligibilityStatus:   "eligible",
		RiskAssessment:      assessment,
		RecommendedCoverage: recommendedCoverage,
		AvailableOptions:    options,
		LastUpdated:         time.Now(),
	}

	return profile, nil
}

func (s *DeviceCategoryService) performRiskAssessment(ctx context.Context, spec base.CategorySpec, history base.ClaimHistory) (base.RiskAssessment, error) {
	return s.categoryService.AssessRisk(spec, history)
}

// Request and response types

// RegisterDeviceRequest represents a device registration request
type RegisterDeviceRequest struct {
	Category       base.CategoryType `json:"category"`
	Specifications json.RawMessage   `json:"specifications"`
	IMEI           string            `json:"imei"`
	SerialNumber   string            `json:"serial_number"`
	OwnerID        uuid.UUID         `json:"owner_id"`
}

// PremiumCalculation represents a premium calculation result
type PremiumCalculation struct {
	DeviceID       uuid.UUID            `json:"device_id"`
	Coverage       base.CoverageLevel   `json:"coverage"`
	MonthlyPremium float64              `json:"monthly_premium"`
	AnnualPremium  float64              `json:"annual_premium"`
	CoverageOption *base.CoverageOption `json:"coverage_option,omitempty"`
	CalculatedAt   time.Time            `json:"calculated_at"`
}

// DeviceValuation represents device valuation data
type DeviceValuation struct {
	DeviceID         uuid.UUID `json:"device_id"`
	MarketValue      float64   `json:"market_value"`
	DepreciatedValue float64   `json:"depreciated_value"`
	ResaleValue      float64   `json:"resale_value"`
	MaxCoverage      float64   `json:"max_coverage"`
	ValuationDate    time.Time `json:"valuation_date"`
}

// Domain events

// DeviceRegisteredEvent is published when a device is registered
type DeviceRegisteredEvent struct {
	DeviceID  uuid.UUID         `json:"device_id"`
	Category  base.CategoryType `json:"category"`
	OwnerID   uuid.UUID         `json:"owner_id"`
	Timestamp time.Time         `json:"timestamp"`
}

// RiskAssessmentUpdatedEvent is published when risk assessment is updated
type RiskAssessmentUpdatedEvent struct {
	DeviceID       uuid.UUID `json:"device_id"`
	RiskScore      float64   `json:"risk_score"`
	RiskLevel      string    `json:"risk_level"`
	RequiresReview bool      `json:"requires_review"`
	Timestamp      time.Time `json:"timestamp"`
}
