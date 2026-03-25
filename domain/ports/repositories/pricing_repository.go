package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// PricingRepository defines the interface for pricing persistence operations
type PricingRepository interface {
	// Pricing rule operations
	CreatePricingRule(ctx context.Context, rule *shared.PricingRule) error
	GetPricingRuleByID(ctx context.Context, id uuid.UUID) (*shared.PricingRule, error)
	GetPricingRuleByName(ctx context.Context, ruleName string) (*shared.PricingRule, error)
	UpdatePricingRule(ctx context.Context, rule *shared.PricingRule) error
	DeletePricingRule(ctx context.Context, id uuid.UUID) error
	ListPricingRules(ctx context.Context, ruleType string, isActive *bool, limit, offset int) ([]*shared.PricingRule, int64, error)

	// Usage-based insurance operations
	CreateUBI(ctx context.Context, ubi *shared.UsageBasedInsurance) error
	GetUBIByID(ctx context.Context, id uuid.UUID) (*shared.UsageBasedInsurance, error)
	GetUBIByPolicyID(ctx context.Context, policyID uuid.UUID) (*shared.UsageBasedInsurance, error)
	GetUBIByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.UsageBasedInsurance, int64, error)
	UpdateUBI(ctx context.Context, ubi *shared.UsageBasedInsurance) error
	GetActiveUBIs(ctx context.Context, limit, offset int) ([]*shared.UsageBasedInsurance, int64, error)

	// Usage record operations
	CreateUsageRecord(ctx context.Context, record *shared.UsageRecord) error
	GetUsageRecordsByUBIID(ctx context.Context, ubiID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*shared.UsageRecord, int64, error)
	GetUsageRecordsByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]*shared.UsageRecord, int64, error)

	// Micro insurance operations
	CreateMicroInsurance(ctx context.Context, micro *shared.MicroInsurance) error
	GetMicroInsuranceByID(ctx context.Context, id uuid.UUID) (*shared.MicroInsurance, error)
	GetMicroInsuranceByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.MicroInsurance, int64, error)
	GetActiveMicroInsurance(ctx context.Context, limit, offset int) ([]*shared.MicroInsurance, int64, error)
	UpdateMicroInsurance(ctx context.Context, micro *shared.MicroInsurance) error

	// Pricing calculation operations
	GetApplicableRules(ctx context.Context, ruleType string, conditions map[string]interface{}) ([]*shared.PricingRule, error)
	GetRulesForRecalculation(ctx context.Context, recalculationTime time.Time) ([]*shared.UsageBasedInsurance, error)

	// Statistics
	GetPricingStatistics(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)
	GetUBIStatistics(ctx context.Context, policyID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
}
