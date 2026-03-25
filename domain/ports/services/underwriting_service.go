package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// UnderwritingService defines the interface for underwriting business logic operations
type UnderwritingService interface {
	// Underwriting decision operations
	CreateUnderwritingDecision(ctx context.Context, decision *shared.UnderwritingDecision) error
	GetUnderwritingDecisionByID(ctx context.Context, id uuid.UUID) (*shared.UnderwritingDecision, error)
	GetUnderwritingDecisionByDecisionNumber(ctx context.Context, decisionNumber string) (*shared.UnderwritingDecision, error)
	UpdateUnderwritingDecision(ctx context.Context, decision *shared.UnderwritingDecision) error
	ReviewUnderwritingDecision(ctx context.Context, decisionID uuid.UUID, underwriterID uuid.UUID, decision string, notes string) error

	// Underwriting queries
	GetUnderwritingDecisionByPolicyID(ctx context.Context, policyID uuid.UUID) (*shared.UnderwritingDecision, error)
	GetUnderwritingDecisionByQuoteID(ctx context.Context, quoteID uuid.UUID) (*shared.UnderwritingDecision, error)
	GetUnderwritingDecisionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	GetUnderwritingDecisionsByDecision(ctx context.Context, decision string, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)
	GetPendingUnderwritingDecisions(ctx context.Context, limit, offset int) ([]*shared.UnderwritingDecision, int64, error)

	// Underwriting rule management
	CreateUnderwritingRule(ctx context.Context, rule *shared.UnderwritingRule) error
	GetUnderwritingRuleByID(ctx context.Context, id uuid.UUID) (*shared.UnderwritingRule, error)
	GetUnderwritingRuleByCode(ctx context.Context, ruleCode string) (*shared.UnderwritingRule, error)
	GetActiveUnderwritingRules(ctx context.Context, category string, limit, offset int) ([]*shared.UnderwritingRule, int64, error)
	UpdateUnderwritingRule(ctx context.Context, rule *shared.UnderwritingRule) error
	DeleteUnderwritingRule(ctx context.Context, id uuid.UUID) error

	// Risk assessment and decision
	AssessRisk(ctx context.Context, userID, deviceID uuid.UUID, quoteID *uuid.UUID) (*shared.UnderwritingDecision, error)
	AutoUnderwrite(ctx context.Context, quoteID uuid.UUID) (*shared.UnderwritingDecision, error)
	ManualUnderwrite(ctx context.Context, quoteID uuid.UUID, underwriterID uuid.UUID, decision string, notes string, conditions []string) (*shared.UnderwritingDecision, error)

	// Rule application
	GetApplicableRules(ctx context.Context, category string, conditions map[string]interface{}) ([]*shared.UnderwritingRule, error)
	EvaluateRules(ctx context.Context, quoteID uuid.UUID, rules []*shared.UnderwritingRule) (map[string]interface{}, error)

	// Pricing adjustments
	CalculateRiskAdjustment(ctx context.Context, basePremium float64, riskScore float64, riskCategory string) (float64, error)
	ApplyLoadingFactor(ctx context.Context, basePremium float64, riskFactors []string) (float64, error)

	// Underwriting statistics
	GetUnderwritingStatistics(ctx context.Context, underwriterID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetApprovalRate(ctx context.Context, startDate, endDate string) (float64, error)
	GetAverageRiskScore(ctx context.Context, decision string, startDate, endDate string) (float64, error)
	GetDecisionDistribution(ctx context.Context, startDate, endDate string) (map[string]int64, error)
}
