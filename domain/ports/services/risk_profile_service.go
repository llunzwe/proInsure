package services

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models/shared"
)

// RiskProfileService defines the interface for risk profile business logic operations
type RiskProfileService interface {
	// Risk profile management
	CreateRiskProfile(ctx context.Context, profile *shared.RiskProfile) error
	GetRiskProfileByID(ctx context.Context, id uuid.UUID) (*shared.RiskProfile, error)
	GetRiskProfileByEntity(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.RiskProfile, error)
	UpdateRiskProfile(ctx context.Context, profile *shared.RiskProfile) error
	DeleteRiskProfile(ctx context.Context, id uuid.UUID) error

	// Risk assessment operations
	AssessRisk(ctx context.Context, entityType string, entityID uuid.UUID, assessmentData map[string]interface{}) (*shared.RiskProfile, error)
	RecalculateRiskScore(ctx context.Context, profileID uuid.UUID) (float64, error)
	UpdateRiskScore(ctx context.Context, profileID uuid.UUID, score float64, riskLevel string) error

	// Risk score updates
	UpdateFraudRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateCreditRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateClaimRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateBehavioralRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateDeviceRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateGeographicRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error

	// Risk factor management
	AddRiskFactor(ctx context.Context, profileID uuid.UUID, factorType string, factors []string) error
	GetRiskFactors(ctx context.Context, profileID uuid.UUID) (map[string][]string, error)
	RemoveRiskFactor(ctx context.Context, profileID uuid.UUID, factorType string, factor string) error

	// Risk profile queries
	GetRiskProfilesByEntityType(ctx context.Context, entityType string, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetRiskProfilesByRiskLevel(ctx context.Context, riskLevel string, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetRiskProfilesByRiskCategory(ctx context.Context, riskCategory string, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetHighRiskProfiles(ctx context.Context, threshold float64, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetLowRiskProfiles(ctx context.Context, threshold float64, limit, offset int) ([]*shared.RiskProfile, int64, error)

	// Risk assessment history
	GetAssessmentHistory(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetLatestAssessment(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.RiskProfile, error)
	CompareAssessments(ctx context.Context, profileID uuid.UUID, previousProfileID uuid.UUID) (map[string]interface{}, error)

	// Risk statistics
	GetRiskDistribution(ctx context.Context, entityType string) (map[string]int64, error)
	GetAverageRiskScore(ctx context.Context, entityType string) (float64, error)
	GetRiskTrends(ctx context.Context, entityType string, entityID uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
}
