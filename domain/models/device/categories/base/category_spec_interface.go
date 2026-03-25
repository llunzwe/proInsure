package base

import (
	"encoding/json"
	"time"
	
	"github.com/google/uuid"
)

// CategoryType represents the type of device category
type CategoryType string

const (
	CategorySmartphone CategoryType = "smartphone"
	CategorySmartwatch CategoryType = "smartwatch"
	CategoryTablet     CategoryType = "tablet"
	CategoryLaptop     CategoryType = "laptop"
	CategoryWearable   CategoryType = "wearable"
	CategoryIoT        CategoryType = "iot"
	CategoryAccessory  CategoryType = "accessory"
)

// CategorySpec defines the interface for category-specific specifications
type CategorySpec interface {
	// GetCategory returns the category type
	GetCategory() CategoryType

	// GetSpecID returns the unique identifier for the specification
	GetSpecID() uuid.UUID

	// Validate checks if the specifications are valid
	Validate() error

	// ToJSON converts the spec to JSON for storage
	ToJSON() (json.RawMessage, error)

	// FromJSON populates the spec from JSON
	FromJSON(data json.RawMessage) error

	// GetManufacturer returns the device manufacturer
	GetManufacturer() string

	// GetModel returns the device model
	GetModel() string

	// GetReleaseDate returns when the device was released
	GetReleaseDate() time.Time

	// GetMarketValue returns the current market value
	GetMarketValue() float64

	// GetDepreciationRate returns the annual depreciation rate
	GetDepreciationRate() float64

	// IsHighEnd determines if the device is considered high-end
	IsHighEnd() bool

	// GetRiskFactors returns device-specific risk factors
	GetRiskFactors() map[string]float64

	// GetRepairCosts returns estimated repair costs by component
	GetRepairCosts() map[string]float64

	// SupportsFeature checks if a specific feature is supported
	SupportsFeature(feature string) bool

	// GetTechnicalSpecs returns all technical specifications
	GetTechnicalSpecs() map[string]interface{}

	// GetWarrantyPeriod returns the standard warranty period in months
	GetWarrantyPeriod() int

	// GetCompatibleAccessories returns list of compatible accessory types
	GetCompatibleAccessories() []string
}

// BaseCategorySpec provides common implementation for all category specs
type BaseCategorySpec struct {
	ID           uuid.UUID              `json:"id"`
	Category     CategoryType           `json:"category"`
	Manufacturer string                 `json:"manufacturer"`
	Model        string                 `json:"model"`
	ReleaseDate  time.Time              `json:"release_date"`
	MarketValue  float64                `json:"market_value"`
	Depreciation float64                `json:"depreciation_rate"`
	RiskFactors  map[string]float64     `json:"risk_factors"`
	RepairCosts  map[string]float64     `json:"repair_costs"`
	Features     []string               `json:"features"`
	TechSpecs    map[string]interface{} `json:"technical_specs"`
	Warranty     int                    `json:"warranty_months"`
	Accessories  []string               `json:"compatible_accessories"`
}

// GetCategory returns the category type
func (b *BaseCategorySpec) GetCategory() CategoryType {
	return b.Category
}

// GetSpecID returns the unique identifier
func (b *BaseCategorySpec) GetSpecID() uuid.UUID {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return b.ID
}

// GetManufacturer returns the manufacturer
func (b *BaseCategorySpec) GetManufacturer() string {
	return b.Manufacturer
}

// GetModel returns the model
func (b *BaseCategorySpec) GetModel() string {
	return b.Model
}

// GetReleaseDate returns the release date
func (b *BaseCategorySpec) GetReleaseDate() time.Time {
	return b.ReleaseDate
}

// GetMarketValue returns the market value
func (b *BaseCategorySpec) GetMarketValue() float64 {
	return b.MarketValue
}

// GetDepreciationRate returns the depreciation rate
func (b *BaseCategorySpec) GetDepreciationRate() float64 {
	return b.Depreciation
}

// IsHighEnd checks if device is high-end based on market value
func (b *BaseCategorySpec) IsHighEnd() bool {
	return b.MarketValue > 800.0
}

// GetRiskFactors returns risk factors
func (b *BaseCategorySpec) GetRiskFactors() map[string]float64 {
	if b.RiskFactors == nil {
		b.RiskFactors = make(map[string]float64)
	}
	return b.RiskFactors
}

// GetRepairCosts returns repair costs
func (b *BaseCategorySpec) GetRepairCosts() map[string]float64 {
	if b.RepairCosts == nil {
		b.RepairCosts = make(map[string]float64)
	}
	return b.RepairCosts
}

// SupportsFeature checks feature support
func (b *BaseCategorySpec) SupportsFeature(feature string) bool {
	for _, f := range b.Features {
		if f == feature {
			return true
		}
	}
	return false
}

// GetTechnicalSpecs returns technical specifications
func (b *BaseCategorySpec) GetTechnicalSpecs() map[string]interface{} {
	if b.TechSpecs == nil {
		b.TechSpecs = make(map[string]interface{})
	}
	return b.TechSpecs
}

// GetWarrantyPeriod returns warranty period
func (b *BaseCategorySpec) GetWarrantyPeriod() int {
	if b.Warranty == 0 {
		return 12 // Default 12 months
	}
	return b.Warranty
}

// GetCompatibleAccessories returns compatible accessories
func (b *BaseCategorySpec) GetCompatibleAccessories() []string {
	return b.Accessories
}

// ToJSON converts to JSON
func (b *BaseCategorySpec) ToJSON() (json.RawMessage, error) {
	return json.Marshal(b)
}

// FromJSON populates from JSON
func (b *BaseCategorySpec) FromJSON(data json.RawMessage) error {
	return json.Unmarshal(data, b)
}

// CategoryMetrics represents performance and usage metrics
type CategoryMetrics struct {
	FailureRate     float64            `json:"failure_rate"`
	AverageLifespan int                `json:"average_lifespan_months"`
	CommonIssues    []string           `json:"common_issues"`
	RepairFrequency map[string]float64 `json:"repair_frequency"`
	CustomerRating  float64            `json:"customer_rating"`
	ReturnRate      float64            `json:"return_rate"`
}
