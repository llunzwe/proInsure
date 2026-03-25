package base

import (
	"encoding/json"
	"fmt"
	"sync"
)

// CategoryFactory creates category-specific instances
type CategoryFactory struct {
	specRegistry       map[CategoryType]func() CategorySpec
	insuranceRegistry  map[CategoryType]func() InsuranceCalculator
	validatorRegistry  map[CategoryType]func() CategoryValidator
	calculatorRegistry map[CategoryType]func() CategoryCalculator
	mu                 sync.RWMutex
}

// NewCategoryFactory creates a new category factory
func NewCategoryFactory() *CategoryFactory {
	return &CategoryFactory{
		specRegistry:       make(map[CategoryType]func() CategorySpec),
		insuranceRegistry:  make(map[CategoryType]func() InsuranceCalculator),
		validatorRegistry:  make(map[CategoryType]func() CategoryValidator),
		calculatorRegistry: make(map[CategoryType]func() CategoryCalculator),
	}
}

// RegisterSpec registers a specification factory for a category
func (f *CategoryFactory) RegisterSpec(category CategoryType, factory func() CategorySpec) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.specRegistry[category] = factory
}

// RegisterInsuranceCalculator registers an insurance calculator factory
func (f *CategoryFactory) RegisterInsuranceCalculator(category CategoryType, factory func() InsuranceCalculator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.insuranceRegistry[category] = factory
}

// RegisterValidator registers a validator factory
func (f *CategoryFactory) RegisterValidator(category CategoryType, factory func() CategoryValidator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.validatorRegistry[category] = factory
}

// RegisterCalculator registers a calculator factory
func (f *CategoryFactory) RegisterCalculator(category CategoryType, factory func() CategoryCalculator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calculatorRegistry[category] = factory
}

// CreateSpec creates a category specification instance
func (f *CategoryFactory) CreateSpec(category CategoryType) (CategorySpec, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	factory, exists := f.specRegistry[category]
	if !exists {
		return nil, fmt.Errorf("no specification registered for category: %s", category)
	}

	return factory(), nil
}

// CreateInsuranceCalculator creates an insurance calculator instance
func (f *CategoryFactory) CreateInsuranceCalculator(category CategoryType) (InsuranceCalculator, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	factory, exists := f.insuranceRegistry[category]
	if !exists {
		return nil, fmt.Errorf("no insurance calculator registered for category: %s", category)
	}

	return factory(), nil
}

// CreateValidator creates a validator instance
func (f *CategoryFactory) CreateValidator(category CategoryType) (CategoryValidator, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	factory, exists := f.validatorRegistry[category]
	if !exists {
		return nil, fmt.Errorf("no validator registered for category: %s", category)
	}

	return factory(), nil
}

// CreateCalculator creates a calculator instance
func (f *CategoryFactory) CreateCalculator(category CategoryType) (CategoryCalculator, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	factory, exists := f.calculatorRegistry[category]
	if !exists {
		return nil, fmt.Errorf("no calculator registered for category: %s", category)
	}

	return factory(), nil
}

// CreateSpecFromJSON creates a spec from JSON data
func (f *CategoryFactory) CreateSpecFromJSON(category CategoryType, data json.RawMessage) (CategorySpec, error) {
	spec, err := f.CreateSpec(category)
	if err != nil {
		return nil, err
	}

	if err := spec.FromJSON(data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}

	return spec, nil
}

// GetSupportedCategories returns all supported categories
func (f *CategoryFactory) GetSupportedCategories() []CategoryType {
	f.mu.RLock()
	defer f.mu.RUnlock()

	categories := make([]CategoryType, 0, len(f.specRegistry))
	for category := range f.specRegistry {
		categories = append(categories, category)
	}
	return categories
}

// IsSupported checks if a category is supported
func (f *CategoryFactory) IsSupported(category CategoryType) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	_, exists := f.specRegistry[category]
	return exists
}

// CategoryManager manages category-specific operations
type CategoryManager struct {
	factory    *CategoryFactory
	specs      map[string]CategorySpec
	validators map[CategoryType]CategoryValidator
	mu         sync.RWMutex
}

// NewCategoryManager creates a new category manager
func NewCategoryManager(factory *CategoryFactory) *CategoryManager {
	return &CategoryManager{
		factory:    factory,
		specs:      make(map[string]CategorySpec),
		validators: make(map[CategoryType]CategoryValidator),
	}
}

// ProcessDevice processes a device based on its category
func (m *CategoryManager) ProcessDevice(deviceID string, category CategoryType, specData json.RawMessage) error {
	// Create specification
	spec, err := m.factory.CreateSpecFromJSON(category, specData)
	if err != nil {
		return fmt.Errorf("failed to create spec: %w", err)
	}

	// Get or create validator
	validator, err := m.getOrCreateValidator(category)
	if err != nil {
		return fmt.Errorf("failed to get validator: %w", err)
	}

	// Validate specification
	if err := validator.ValidateSpec(spec); err != nil {
		return fmt.Errorf("specification validation failed: %w", err)
	}

	// Store specification
	m.mu.Lock()
	m.specs[deviceID] = spec
	m.mu.Unlock()

	return nil
}

// GetSpec retrieves a stored specification
func (m *CategoryManager) GetSpec(deviceID string) (CategorySpec, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	spec, exists := m.specs[deviceID]
	return spec, exists
}

// CalculatePremium calculates insurance premium for a device
func (m *CategoryManager) CalculatePremium(deviceID string, coverage CoverageLevel) (float64, error) {
	spec, exists := m.GetSpec(deviceID)
	if !exists {
		return 0, fmt.Errorf("device specification not found: %s", deviceID)
	}

	calculator, err := m.factory.CreateInsuranceCalculator(spec.GetCategory())
	if err != nil {
		return 0, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	return calculator.CalculatePremium(spec, coverage)
}

// AssessRisk performs risk assessment for a device
func (m *CategoryManager) AssessRisk(deviceID string, history ClaimHistory) (RiskAssessment, error) {
	spec, exists := m.GetSpec(deviceID)
	if !exists {
		return RiskAssessment{}, fmt.Errorf("device specification not found: %s", deviceID)
	}

	calculator, err := m.factory.CreateInsuranceCalculator(spec.GetCategory())
	if err != nil {
		return RiskAssessment{}, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	return calculator.AssessRisk(spec, history)
}

// getOrCreateValidator gets or creates a validator for a category
func (m *CategoryManager) getOrCreateValidator(category CategoryType) (CategoryValidator, error) {
	m.mu.RLock()
	validator, exists := m.validators[category]
	m.mu.RUnlock()

	if exists {
		return validator, nil
	}

	// Create new validator
	validator, err := m.factory.CreateValidator(category)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.validators[category] = validator
	m.mu.Unlock()

	return validator, nil
}

// CategoryRegistry global registry for category implementations
var globalRegistry *CategoryFactory
var registryOnce sync.Once

// GetGlobalFactory returns the global category factory
func GetGlobalFactory() *CategoryFactory {
	registryOnce.Do(func() {
		globalRegistry = NewCategoryFactory()
	})
	return globalRegistry
}
