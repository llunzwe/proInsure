package categories

import (
	"fmt"
	"sync"
	"time"
	
	"smartsure/internal/domain/models/device/categories/base"
	"smartsure/internal/domain/models/device/categories/smartphone"
	"smartsure/internal/domain/models/device/categories/smartwatch"
)

// init registers all category implementations
func init() {
	registry := base.GetGlobalFactory()

	// Register smartphone implementations
	registry.RegisterSpec(base.CategorySmartphone, func() base.CategorySpec {
		return smartphone.NewSmartphoneSpec()
	})

	registry.RegisterInsuranceCalculator(base.CategorySmartphone, func() base.InsuranceCalculator {
		return smartphone.NewSmartphoneInsuranceCalculator()
	})

	registry.RegisterValidator(base.CategorySmartphone, func() base.CategoryValidator {
		return smartphone.NewSmartphoneValidator()
	})

	registry.RegisterCalculator(base.CategorySmartphone, func() base.CategoryCalculator {
		return smartphone.NewSmartphoneRiskCalculator()
	})

	// Register smartwatch implementations
	registry.RegisterSpec(base.CategorySmartwatch, func() base.CategorySpec {
		return smartwatch.NewSmartwatchSpec()
	})

	registry.RegisterInsuranceCalculator(base.CategorySmartwatch, func() base.InsuranceCalculator {
		return smartwatch.NewSmartwatchInsuranceCalculator()
	})

	registry.RegisterValidator(base.CategorySmartwatch, func() base.CategoryValidator {
		return smartwatch.NewSmartwatchValidator()
	})

	registry.RegisterCalculator(base.CategorySmartwatch, func() base.CategoryCalculator {
		return smartwatch.NewSmartwatchRiskCalculator()
	})
}

// CategoryService provides high-level category operations
type CategoryService struct {
	factory *base.CategoryFactory
	manager *base.CategoryManager
	mu      sync.RWMutex
}

// NewCategoryService creates a new category service
func NewCategoryService() *CategoryService {
	factory := base.GetGlobalFactory()
	return &CategoryService{
		factory: factory,
		manager: base.NewCategoryManager(factory),
	}
}

// GetSupportedCategories returns all supported device categories
func (s *CategoryService) GetSupportedCategories() []base.CategoryType {
	return s.factory.GetSupportedCategories()
}

// IsSupported checks if a category is supported
func (s *CategoryService) IsSupported(category base.CategoryType) bool {
	return s.factory.IsSupported(category)
}

// CreateSpec creates a specification for a category
func (s *CategoryService) CreateSpec(category base.CategoryType) (base.CategorySpec, error) {
	return s.factory.CreateSpec(category)
}

// ValidateDevice validates a device specification
func (s *CategoryService) ValidateDevice(category base.CategoryType, spec base.CategorySpec) error {
	validator, err := s.factory.CreateValidator(category)
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}

	return validator.ValidateSpec(spec)
}

// CalculatePremium calculates insurance premium for a device
func (s *CategoryService) CalculatePremium(spec base.CategorySpec, coverage base.CoverageLevel) (float64, error) {
	calculator, err := s.factory.CreateInsuranceCalculator(spec.GetCategory())
	if err != nil {
		return 0, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	return calculator.CalculatePremium(spec, coverage)
}

// AssessRisk performs risk assessment for a device
func (s *CategoryService) AssessRisk(spec base.CategorySpec, history base.ClaimHistory) (base.RiskAssessment, error) {
	calculator, err := s.factory.CreateInsuranceCalculator(spec.GetCategory())
	if err != nil {
		return base.RiskAssessment{}, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	return calculator.AssessRisk(spec, history)
}

// GetCoverageOptions returns available coverage options for a device
func (s *CategoryService) GetCoverageOptions(spec base.CategorySpec) ([]base.CoverageOption, error) {
	calculator, err := s.factory.CreateInsuranceCalculator(spec.GetCategory())
	if err != nil {
		return nil, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	return calculator.GetCoverageOptions(spec), nil
}

// CalculateDepreciation calculates depreciated value
func (s *CategoryService) CalculateDepreciation(spec base.CategorySpec, purchaseDate string) (float64, error) {
	calculator, err := s.factory.CreateCalculator(spec.GetCategory())
	if err != nil {
		return 0, fmt.Errorf("failed to create calculator: %w", err)
	}

	// Parse purchase date
	var pd time.Time
	if err := pd.UnmarshalText([]byte(purchaseDate)); err != nil {
		return 0, fmt.Errorf("invalid purchase date: %w", err)
	}

	return calculator.CalculateDepreciation(spec.GetMarketValue(), pd), nil
}

// CalculateResaleValue calculates resale value
func (s *CategoryService) CalculateResaleValue(spec base.CategorySpec, condition string) (float64, error) {
	calculator, err := s.factory.CreateCalculator(spec.GetCategory())
	if err != nil {
		return 0, fmt.Errorf("failed to create calculator: %w", err)
	}

	return calculator.CalculateResaleValue(spec, condition), nil
}

// ValidateEligibility checks if device is eligible for insurance
func (s *CategoryService) ValidateEligibility(spec base.CategorySpec) (bool, []string, error) {
	calculator, err := s.factory.CreateInsuranceCalculator(spec.GetCategory())
	if err != nil {
		return false, nil, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	eligible, issues := calculator.ValidateEligibility(spec)
	return eligible, issues, nil
}

// CategoryStatistics provides category-level statistics
type CategoryStatistics struct {
	Category         base.CategoryType `json:"category"`
	TotalDevices     int               `json:"total_devices"`
	ActivePolicies   int               `json:"active_policies"`
	TotalClaims      int               `json:"total_claims"`
	AverageRiskScore float64           `json:"average_risk_score"`
	AveragePremium   float64           `json:"average_premium"`
	PopularModels    []string          `json:"popular_models"`
}

// GetCategoryStatistics returns statistics for a category
func (s *CategoryService) GetCategoryStatistics(category base.CategoryType) (*CategoryStatistics, error) {
	if !s.IsSupported(category) {
		return nil, fmt.Errorf("unsupported category: %s", category)
	}

	// In a real implementation, this would query the database
	// For now, return mock data
	return &CategoryStatistics{
		Category:         category,
		TotalDevices:     1000,
		ActivePolicies:   850,
		TotalClaims:      120,
		AverageRiskScore: 45.5,
		AveragePremium:   25.99,
		PopularModels: []string{
			"iPhone 14 Pro",
			"Samsung Galaxy S23",
			"Google Pixel 7",
		},
	}, nil
}

// CategoryReport generates a report for a category
type CategoryReport struct {
	Category    base.CategoryType      `json:"category"`
	Period      string                 `json:"period"`
	Statistics  *CategoryStatistics    `json:"statistics"`
	TopRisks    []string               `json:"top_risks"`
	Trends      map[string]interface{} `json:"trends"`
	GeneratedAt time.Time              `json:"generated_at"`
}

// GenerateCategoryReport generates a comprehensive report for a category
func (s *CategoryService) GenerateCategoryReport(category base.CategoryType, period string) (*CategoryReport, error) {
	stats, err := s.GetCategoryStatistics(category)
	if err != nil {
		return nil, err
	}

	report := &CategoryReport{
		Category:    category,
		Period:      period,
		Statistics:  stats,
		GeneratedAt: time.Now(),
		TopRisks: []string{
			"Screen damage",
			"Water damage",
			"Battery degradation",
		},
		Trends: map[string]interface{}{
			"claim_increase": "+5%",
			"premium_change": "-2%",
			"new_policies":   "+12%",
		},
	}

	return report, nil
}
