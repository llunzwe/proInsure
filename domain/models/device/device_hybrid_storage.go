package device

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/device/categories/base"
	"smartsure/pkg/database"
)

// DeviceHybrid represents a device with hybrid polymorphic storage
type DeviceHybrid struct {
	database.BaseModel

	// Core fields (strongly typed)
	IMEI         string            `gorm:"type:varchar(15);uniqueIndex;not null" json:"imei"`
	SerialNumber string            `gorm:"type:varchar(50);uniqueIndex;not null" json:"serial_number"`
	Category     base.CategoryType `gorm:"type:varchar(50);not null;index" json:"category"`
	OwnerID      uuid.UUID         `gorm:"type:uuid;not null;index" json:"owner_id"`

	// Flexible category-specific data (JSONB)
	Specifications CategorySpecData `gorm:"type:jsonb;not null" json:"specifications"`

	// Insurance and risk data (JSONB for flexibility)
	InsuranceData *InsuranceData `gorm:"type:jsonb" json:"insurance_data,omitempty"`
	RiskData      *RiskData      `gorm:"type:jsonb" json:"risk_data,omitempty"`

	// Common device fields
	Status    DeviceStatus    `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	Condition DeviceCondition `gorm:"type:varchar(20);not null;default:'good'" json:"condition"`
	Grade     DeviceGrade     `gorm:"type:varchar(20);not null;default:'B'" json:"grade"`

	// Computed fields for indexing and queries
	Manufacturer string  `gorm:"type:varchar(100);index" json:"manufacturer"`
	Model        string  `gorm:"type:varchar(100);index" json:"model"`
	MarketValue  float64 `gorm:"type:decimal(10,2);index" json:"market_value"`
	RiskScore    float64 `gorm:"type:decimal(5,2);index" json:"risk_score"`

	// Relationships (loaded via service layer)
	// Owner should be loaded via service layer using OwnerID
	// Policies should be loaded via service layer
	// Claims should be loaded via service layer
}

// CategorySpecData represents flexible category-specific specifications
type CategorySpecData json.RawMessage

// Scan implements the Scanner interface
func (c *CategorySpecData) Scan(value interface{}) error {
	if value == nil {
		*c = CategorySpecData("null")
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*c = CategorySpecData(v)
	case string:
		*c = CategorySpecData(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// Value implements the driver Valuer interface
func (c CategorySpecData) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return []byte(c), nil
}

// InsuranceData represents insurance-related data
type InsuranceData struct {
	EligibilityStatus   string                `json:"eligibility_status"`
	RecommendedCoverage base.CoverageLevel    `json:"recommended_coverage"`
	AvailableOptions    []base.CoverageOption `json:"available_options"`
	PremiumFactors      map[string]float64    `json:"premium_factors"`
	Restrictions        []string              `json:"restrictions,omitempty"`
	LastAssessment      time.Time             `json:"last_assessment"`
}

// Scan implements the Scanner interface
func (i *InsuranceData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, i)
	case string:
		return json.Unmarshal([]byte(v), i)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

// Value implements the driver Valuer interface
func (i InsuranceData) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// RiskData represents risk assessment data
type RiskData struct {
	Score           float64            `json:"score"`
	Level           string             `json:"level"`
	Factors         map[string]float64 `json:"factors"`
	Recommendations []string           `json:"recommendations"`
	RequiresReview  bool               `json:"requires_review"`
	ValidUntil      time.Time          `json:"valid_until"`
}

// Scan implements the Scanner interface
func (r *RiskData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, r)
	case string:
		return json.Unmarshal([]byte(v), r)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

// Value implements the driver Valuer interface
func (r RiskData) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// TableName specifies the table name
func (DeviceHybrid) TableName() string {
	return "devices_hybrid"
}

// Business Methods

// GetCategorySpec returns the category specification as the appropriate type
func (d *DeviceHybrid) GetCategorySpec() (base.CategorySpec, error) {
	factory := base.GetGlobalFactory()
	spec, err := factory.CreateSpecFromJSON(d.Category, json.RawMessage(d.Specifications))
	if err != nil {
		return nil, fmt.Errorf("failed to create spec: %w", err)
	}
	return spec, nil
}

// UpdateSpecifications updates the device specifications
func (d *DeviceHybrid) UpdateSpecifications(spec base.CategorySpec) error {
	// Validate that the category matches
	if spec.GetCategory() != d.Category {
		return fmt.Errorf("category mismatch: expected %s, got %s", d.Category, spec.GetCategory())
	}

	// Convert spec to JSON
	specJSON, err := spec.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize specifications: %w", err)
	}

	d.Specifications = CategorySpecData(specJSON)

	// Update indexed fields
	d.extractIndexedFields(spec)

	d.UpdatedAt = time.Now()
	return nil
}

// extractIndexedFields extracts frequently queried fields for indexing
func (d *DeviceHybrid) extractIndexedFields(spec base.CategorySpec) {
	d.Manufacturer = spec.GetManufacturer()
	d.Model = spec.GetModel()
	d.MarketValue = spec.GetMarketValue()
}

// CalculatePremium calculates insurance premium for the device
func (d *DeviceHybrid) CalculatePremium(coverage base.CoverageLevel) (float64, error) {
	spec, err := d.GetCategorySpec()
	if err != nil {
		return 0, err
	}

	factory := base.GetGlobalFactory()
	calculator, err := factory.CreateInsuranceCalculator(d.Category)
	if err != nil {
		return 0, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	return calculator.CalculatePremium(spec, coverage)
}

// AssessRisk performs risk assessment for the device
func (d *DeviceHybrid) AssessRisk(claimHistory base.ClaimHistory) (base.RiskAssessment, error) {
	spec, err := d.GetCategorySpec()
	if err != nil {
		return base.RiskAssessment{}, err
	}

	factory := base.GetGlobalFactory()
	calculator, err := factory.CreateInsuranceCalculator(d.Category)
	if err != nil {
		return base.RiskAssessment{}, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	assessment, err := calculator.AssessRisk(spec, claimHistory)
	if err != nil {
		return assessment, err
	}

	// Update risk data
	d.RiskData = &RiskData{
		Score:           assessment.Score,
		Level:           assessment.Level,
		Factors:         assessment.Factors,
		Recommendations: assessment.Recommendations,
		RequiresReview:  assessment.RequiresReview,
		ValidUntil:      assessment.ValidUntil,
	}
	d.RiskScore = assessment.Score

	return assessment, nil
}

// ValidateEligibility checks if device is eligible for insurance
func (d *DeviceHybrid) ValidateEligibility() (bool, []string, error) {
	spec, err := d.GetCategorySpec()
	if err != nil {
		return false, nil, err
	}

	factory := base.GetGlobalFactory()
	calculator, err := factory.CreateInsuranceCalculator(d.Category)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	eligible, reasons := calculator.ValidateEligibility(spec)
	return eligible, reasons, nil
}

// GetCoverageOptions returns available coverage options
func (d *DeviceHybrid) GetCoverageOptions() ([]base.CoverageOption, error) {
	spec, err := d.GetCategorySpec()
	if err != nil {
		return nil, err
	}

	factory := base.GetGlobalFactory()
	calculator, err := factory.CreateInsuranceCalculator(d.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to create insurance calculator: %w", err)
	}

	return calculator.GetCoverageOptions(spec), nil
}

// CalculateDepreciation calculates current depreciated value
func (d *DeviceHybrid) CalculateDepreciation() (float64, error) {
	// Ensure spec exists (validation)
	_, err := d.GetCategorySpec()
	if err != nil {
		return 0, err
	}

	factory := base.GetGlobalFactory()
	calculator, err := factory.CreateCalculator(d.Category)
	if err != nil {
		return 0, fmt.Errorf("failed to create calculator: %w", err)
	}

	return calculator.CalculateDepreciation(d.MarketValue, d.CreatedAt), nil
}

// CalculateResaleValue calculates resale value based on condition
func (d *DeviceHybrid) CalculateResaleValue() (float64, error) {
	spec, err := d.GetCategorySpec()
	if err != nil {
		return 0, err
	}

	factory := base.GetGlobalFactory()
	calculator, err := factory.CreateCalculator(d.Category)
	if err != nil {
		return 0, fmt.Errorf("failed to create calculator: %w", err)
	}

	return calculator.CalculateResaleValue(spec, string(d.Condition)), nil
}

// IsHighValue determines if the device is high value
func (d *DeviceHybrid) IsHighValue() bool {
	return d.MarketValue > 1000
}

// NeedsRiskReassessment checks if risk assessment needs updating
func (d *DeviceHybrid) NeedsRiskReassessment() bool {
	if d.RiskData == nil {
		return true
	}
	return time.Now().After(d.RiskData.ValidUntil)
}

// CanBeInsured checks if the device can be insured
func (d *DeviceHybrid) CanBeInsured() bool {
	// Basic checks
	if d.Status != DeviceStatusActive {
		return false
	}

	// Age check (5 years maximum)
	deviceAge := time.Since(d.CreatedAt)
	if deviceAge > (5 * 365 * 24 * time.Hour) {
		return false
	}

	// Value check
	if d.MarketValue < 100 || d.MarketValue > 10000 {
		return false
	}

	// Condition check
	if d.Condition == DeviceConditionBroken || d.Condition == DeviceConditionPoor {
		return false
	}

	return true
}

// DeviceHybridQueryBuilder provides optimized queries for hybrid devices
type DeviceHybridQueryBuilder struct {
	category     base.CategoryType
	minValue     float64
	maxValue     float64
	minRiskScore float64
	maxRiskScore float64
	status       DeviceStatus
	condition    DeviceCondition
	manufacturer string
	ownerID      *uuid.UUID
}

// NewDeviceHybridQueryBuilder creates a new query builder
func NewDeviceHybridQueryBuilder() *DeviceHybridQueryBuilder {
	return &DeviceHybridQueryBuilder{}
}

// WithCategory filters by category
func (q *DeviceHybridQueryBuilder) WithCategory(category base.CategoryType) *DeviceHybridQueryBuilder {
	q.category = category
	return q
}

// WithValueRange filters by market value range
func (q *DeviceHybridQueryBuilder) WithValueRange(min, max float64) *DeviceHybridQueryBuilder {
	q.minValue = min
	q.maxValue = max
	return q
}

// WithRiskScoreRange filters by risk score range
func (q *DeviceHybridQueryBuilder) WithRiskScoreRange(min, max float64) *DeviceHybridQueryBuilder {
	q.minRiskScore = min
	q.maxRiskScore = max
	return q
}

// WithStatus filters by status
func (q *DeviceHybridQueryBuilder) WithStatus(status DeviceStatus) *DeviceHybridQueryBuilder {
	q.status = status
	return q
}

// WithOwner filters by owner ID
func (q *DeviceHybridQueryBuilder) WithOwner(ownerID uuid.UUID) *DeviceHybridQueryBuilder {
	q.ownerID = &ownerID
	return q
}

// Build builds the query conditions
func (q *DeviceHybridQueryBuilder) Build() map[string]interface{} {
	conditions := make(map[string]interface{})

	if q.category != "" {
		conditions["category"] = q.category
	}
	if q.minValue > 0 {
		conditions["min_value"] = q.minValue
	}
	if q.maxValue > 0 {
		conditions["max_value"] = q.maxValue
	}
	if q.minRiskScore > 0 {
		conditions["min_risk_score"] = q.minRiskScore
	}
	if q.maxRiskScore > 0 {
		conditions["max_risk_score"] = q.maxRiskScore
	}
	if q.status != "" {
		conditions["status"] = q.status
	}
	if q.condition != "" {
		conditions["condition"] = q.condition
	}
	if q.manufacturer != "" {
		conditions["manufacturer"] = q.manufacturer
	}
	if q.ownerID != nil {
		conditions["owner_id"] = *q.ownerID
	}

	return conditions
}
