package base

import (
	"fmt"
	"regexp"
	"time"
)

// CategoryValidator defines validation interface for category-specific data
type CategoryValidator interface {
	// ValidateSpec validates the complete specification
	ValidateSpec(spec CategorySpec) error

	// ValidateIMEI validates IMEI for the specific category
	ValidateIMEI(imei string) error

	// ValidateSerialNumber validates serial number format
	ValidateSerialNumber(serial string) error

	// ValidateModel validates if model exists and is valid
	ValidateModel(manufacturer, model string) error

	// ValidateAge checks if device age is within acceptable limits
	ValidateAge(releaseDate time.Time) error

	// ValidateCondition validates device condition for insurance
	ValidateCondition(condition string) error

	// ValidateMarketValue validates if market value is realistic
	ValidateMarketValue(value float64, model string) error

	// ValidateRepairCost validates repair cost estimates
	ValidateRepairCost(component string, cost float64) error

	// GetValidationRules returns all validation rules for the category
	GetValidationRules() map[string]ValidationRule

	// IsCompatible checks if accessories/parts are compatible
	IsCompatible(itemType, itemModel string) bool
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	Field      string                  `json:"field"`
	Type       string                  `json:"type"`
	Required   bool                    `json:"required"`
	Min        interface{}             `json:"min,omitempty"`
	Max        interface{}             `json:"max,omitempty"`
	Pattern    string                  `json:"pattern,omitempty"`
	Values     []string                `json:"values,omitempty"`
	CustomFunc func(interface{}) error `json:"-"`
	ErrorMsg   string                  `json:"error_message"`
}

// ValidationResult contains validation results
type ValidationResult struct {
	IsValid  bool              `json:"is_valid"`
	Errors   []ValidationError `json:"errors"`
	Warnings []string          `json:"warnings"`
	Score    float64           `json:"validation_score"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// BaseValidator provides common validation functionality
type BaseValidator struct {
	Rules         map[string]ValidationRule
	IMEIPattern   *regexp.Regexp
	SerialPattern *regexp.Regexp
	MaxAge        int // Maximum age in days
	MinValue      float64
	MaxValue      float64
}

// NewBaseValidator creates a new base validator
func NewBaseValidator() *BaseValidator {
	return &BaseValidator{
		Rules:         make(map[string]ValidationRule),
		IMEIPattern:   regexp.MustCompile(`^\d{15}$`),
		SerialPattern: regexp.MustCompile(`^[A-Z0-9]{8,20}$`),
		MaxAge:        1825, // 5 years
		MinValue:      50.0,
		MaxValue:      50000.0,
	}
}

// ValidateIMEI validates IMEI format and checksum
func (b *BaseValidator) ValidateIMEI(imei string) error {
	if !b.IMEIPattern.MatchString(imei) {
		return fmt.Errorf("invalid IMEI format: must be 15 digits")
	}

	// Luhn algorithm validation
	if !b.validateLuhn(imei) {
		return fmt.Errorf("invalid IMEI: checksum validation failed")
	}

	return nil
}

// ValidateSerialNumber validates serial number format
func (b *BaseValidator) ValidateSerialNumber(serial string) error {
	if !b.SerialPattern.MatchString(serial) {
		return fmt.Errorf("invalid serial number format: must be 8-20 alphanumeric characters")
	}
	return nil
}

// ValidateAge validates device age
func (b *BaseValidator) ValidateAge(releaseDate time.Time) error {
	age := time.Since(releaseDate)
	if age.Hours()/24 > float64(b.MaxAge) {
		return fmt.Errorf("device too old: maximum age is %d days", b.MaxAge)
	}
	if age < 0 {
		return fmt.Errorf("invalid release date: cannot be in the future")
	}
	return nil
}

// ValidateMarketValue validates market value
func (b *BaseValidator) ValidateMarketValue(value float64, model string) error {
	if value < b.MinValue {
		return fmt.Errorf("market value too low: minimum is $%.2f", b.MinValue)
	}
	if value > b.MaxValue {
		return fmt.Errorf("market value too high: maximum is $%.2f", b.MaxValue)
	}
	return nil
}

// validateLuhn implements Luhn algorithm for IMEI validation
func (b *BaseValidator) validateLuhn(imei string) bool {
	sum := 0
	double := false

	for i := len(imei) - 1; i >= 0; i-- {
		digit := int(imei[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}

// GetValidationRules returns validation rules
func (b *BaseValidator) GetValidationRules() map[string]ValidationRule {
	return b.Rules
}

// CategoryValidationConfig contains category-specific validation configuration
type CategoryValidationConfig struct {
	CategoryType        CategoryType              `json:"category_type"`
	RequiredFields      []string                  `json:"required_fields"`
	OptionalFields      []string                  `json:"optional_fields"`
	FieldValidators     map[string]ValidationRule `json:"field_validators"`
	BusinessRules       []BusinessRule            `json:"business_rules"`
	CompatibilityMatrix map[string][]string       `json:"compatibility_matrix"`
}

// BusinessRule represents a business validation rule
type BusinessRule struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Condition   func(interface{}) bool  `json:"-"`
	Action      func(interface{}) error `json:"-"`
	Severity    string                  `json:"severity"` // warning, error, critical
}

// ValidatorFactory creates validators for specific categories
type ValidatorFactory interface {
	CreateValidator(category CategoryType) (CategoryValidator, error)
	RegisterValidator(category CategoryType, validator CategoryValidator)
	GetSupportedCategories() []CategoryType
}
