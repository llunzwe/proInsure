package repositories

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models/shared"
)

// RiskProfileRepository defines the interface for risk profile persistence operations
type RiskProfileRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, profile *shared.RiskProfile) error
	GetByID(ctx context.Context, id uuid.UUID) (*shared.RiskProfile, error)
	GetByProfileCode(ctx context.Context, profileCode string) (*shared.RiskProfile, error)
	Update(ctx context.Context, profile *shared.RiskProfile) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.RiskProfile, error)
	GetByEntityType(ctx context.Context, entityType string, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetByRiskLevel(ctx context.Context, riskLevel string, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetByRiskCategory(ctx context.Context, riskCategory string, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetHighRiskProfiles(ctx context.Context, threshold float64, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetLowRiskProfiles(ctx context.Context, threshold float64, limit, offset int) ([]*shared.RiskProfile, int64, error)
	Search(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*shared.RiskProfile, int64, error)

	// Risk score operations
	UpdateRiskScore(ctx context.Context, profileID uuid.UUID, score float64, riskLevel string) error
	UpdateFraudRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateCreditRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateClaimRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateBehavioralRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateDeviceRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error
	UpdateGeographicRiskScore(ctx context.Context, profileID uuid.UUID, score float64) error

	// Risk factor operations
	AddRiskFactor(ctx context.Context, profileID uuid.UUID, factorType string, factors []string) error
	GetRiskFactors(ctx context.Context, profileID uuid.UUID) (map[string][]string, error)

	// Assessment history
	GetAssessmentHistory(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]*shared.RiskProfile, int64, error)
	GetLatestAssessment(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.RiskProfile, error)

	// Statistics
	GetRiskDistribution(ctx context.Context, entityType string) (map[string]int64, error)
	GetAverageRiskScore(ctx context.Context, entityType string) (float64, error)
}
