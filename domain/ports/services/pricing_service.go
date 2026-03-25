package services

import (
	"context"

	"smartsure/internal/domain/models/shared"

	"github.com/google/uuid"
)

// PricingService defines the interface for pricing business logic operations
type PricingService interface {
	// Pricing calculation
	CalculatePremium(ctx context.Context, userID, deviceID uuid.UUID, coverageType string, coverageAmount float64, term int) (float64, error)
	CalculateDynamicPremium(ctx context.Context, policyID uuid.UUID, adjustments map[string]interface{}) (float64, error)
	RecalculatePremium(ctx context.Context, policyID uuid.UUID, newFactors map[string]interface{}) (float64, error)

	// Pricing rule management
	CreatePricingRule(ctx context.Context, rule *shared.PricingRule) error
	GetPricingRuleByID(ctx context.Context, id uuid.UUID) (*shared.PricingRule, error)
	GetPricingRuleByName(ctx context.Context, ruleName string) (*shared.PricingRule, error)
	UpdatePricingRule(ctx context.Context, rule *shared.PricingRule) error
	DeletePricingRule(ctx context.Context, id uuid.UUID) error
	ListPricingRules(ctx context.Context, ruleType string, isActive *bool, limit, offset int) ([]*shared.PricingRule, int64, error)
	GetApplicableRules(ctx context.Context, ruleType string, conditions map[string]interface{}) ([]*shared.PricingRule, error)

	// Usage-based insurance operations
	CreateUBI(ctx context.Context, ubi *shared.UsageBasedInsurance) error
	GetUBIByID(ctx context.Context, id uuid.UUID) (*shared.UsageBasedInsurance, error)
	GetUBIByPolicyID(ctx context.Context, policyID uuid.UUID) (*shared.UsageBasedInsurance, error)
	UpdateUBI(ctx context.Context, ubi *shared.UsageBasedInsurance) error
	RecalculateUBIPremium(ctx context.Context, ubiID uuid.UUID) (float64, error)
	GetActiveUBIs(ctx context.Context, limit, offset int) ([]*shared.UsageBasedInsurance, int64, error)

	// Usage record operations
	RecordUsage(ctx context.Context, ubiID uuid.UUID, record *shared.UsageRecord) error
	GetUsageRecords(ctx context.Context, ubiID uuid.UUID, startDate, endDate string, limit, offset int) ([]*shared.UsageRecord, int64, error)
	GetUsageStatistics(ctx context.Context, ubiID uuid.UUID, startDate, endDate string) (map[string]interface{}, error)

	// Micro insurance operations
	CreateMicroInsurance(ctx context.Context, micro *shared.MicroInsurance) error
	GetMicroInsuranceByID(ctx context.Context, id uuid.UUID) (*shared.MicroInsurance, error)
	GetMicroInsuranceByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.MicroInsurance, int64, error)
	UpdateMicroInsurance(ctx context.Context, micro *shared.MicroInsurance) error
	GetActiveMicroInsurance(ctx context.Context, limit, offset int) ([]*shared.MicroInsurance, int64, error)

	// Discount and promotion operations
	ApplyDiscount(ctx context.Context, basePremium float64, discountCode string, userID uuid.UUID) (float64, error)
	ValidateDiscountCode(ctx context.Context, discountCode string, userID uuid.UUID) (bool, float64, error)
	CalculateBulkDiscount(ctx context.Context, basePremium float64, deviceCount int) (float64, error)

	// Pricing statistics
	GetPricingStatistics(ctx context.Context, startDate, endDate string) (map[string]interface{}, error)
	GetUBIStatistics(ctx context.Context, policyID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetAveragePremium(ctx context.Context, coverageType string, startDate, endDate string) (float64, error)
}
