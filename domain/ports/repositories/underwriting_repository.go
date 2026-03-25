package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// UnderwritingRepository defines the interface for underwriting persistence operations
type UnderwritingRepository interface {
	// Underwriting decision operations
	CreateUnderwritingDecision(ctx context.Context, decision *shared.UnderwritingDecision) error
	GetUnderwritingDecisionByID(ctx context.Context, id uuid.UUID) (*shared.UnderwritingDecision, error)
	GetUnderwritingDecisionByDecisionNumber(ctx context.Context, decisionNumber string) (*shared.UnderwritingDecision, error)
	UpdateUnderwritingDecision(ctx context.Context, decision *shared.UnderwritingDecision) error
	DeleteUnderwritingDecision(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetUnderwritingDecisionByPolicyID(ctx context.Context, policyID uuid.UUID) (*shared.UnderwritingDecision, error)
	GetUnderwritingDecisionByQuoteID(ctx context.Context, quoteID uuid.UUID) (*shared.UnderwritingDecision, error)
	GetUnderwritingDecisionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	GetUnderwritingDecisionsByDecision(ctx context.Context, decision string, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	GetUnderwritingDecisionsByRiskCategory(ctx context.Context, riskCategory string, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	GetUnderwritingDecisionsByUnderwriterID(ctx context.Context, underwriterID uuid.UUID, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	GetPendingUnderwritingDecisions(ctx context.Context, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	GetExpiredUnderwritingDecisions(ctx context.Context, beforeDate time.Time, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	SearchUnderwritingDecisions(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)

	// Underwriting rule operations
	CreateUnderwritingRule(ctx context.Context, rule *shared.UnderwritingRule) error
	GetUnderwritingRuleByID(ctx context.Context, id uuid.UUID) (*shared.UnderwritingRule, error)
	GetUnderwritingRuleByCode(ctx context.Context, ruleCode string) (*shared.UnderwritingRule, error)
	GetActiveUnderwritingRules(ctx context.Context, category string, limit, offset int) ([]*shared.UnderwritingRule, int64, error)
	UpdateUnderwritingRule(ctx context.Context, rule *shared.UnderwritingRule) error
	DeleteUnderwritingRule(ctx context.Context, id uuid.UUID) error

	// Rule application operations
	GetApplicableRules(ctx context.Context, category string, conditions map[string]interface{}) ([]*shared.UnderwritingRule, error)
	GetRulesByPriority(ctx context.Context, category string) ([]*shared.UnderwritingRule, error)

	// Statistics
	GetUnderwritingStatistics(ctx context.Context, underwriterID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetApprovalRate(ctx context.Context, startDate, endDate time.Time) (float64, error)
	GetAverageRiskScore(ctx context.Context, decision string, startDate, endDate time.Time) (float64, error)
	GetDecisionDistribution(ctx context.Context, startDate, endDate time.Time) (map[string]int64, error)
}
