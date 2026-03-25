package services

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models/shared"
)

// FraudService defines the interface for fraud detection business logic operations
type FraudService interface {
	// Fraud detection operations
	AnalyzeFraudRisk(ctx context.Context, entityType string, entityID uuid.UUID, analysisType string) (*shared.FraudDetection, error)
	GetFraudDetectionByID(ctx context.Context, id uuid.UUID) (*shared.FraudDetection, error)
	UpdateFraudDetection(ctx context.Context, detection *shared.FraudDetection) error
	ReviewFraudDetection(ctx context.Context, detectionID uuid.UUID, reviewedBy uuid.UUID, decision string, notes string) error

	// Fraud detection queries
	GetFraudDetectionsByClaimID(ctx context.Context, claimID uuid.UUID) ([]*shared.FraudDetection, error)
	GetFraudDetectionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetFraudDetectionsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetHighRiskDetections(ctx context.Context, threshold float64, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetPendingReviews(ctx context.Context, limit, offset int) ([]*shared.FraudDetection, int64, error)

	// Fraud indicator operations
	AddFraudIndicator(ctx context.Context, detectionID uuid.UUID, indicator *shared.FraudIndicator) error
	GetFraudIndicatorsByDetectionID(ctx context.Context, detectionID uuid.UUID) ([]*shared.FraudIndicator, error)
	ConfirmFraudIndicator(ctx context.Context, indicatorID uuid.UUID, isConfirmed bool) error

	// Blacklisted device operations
	BlacklistDevice(ctx context.Context, imei string, serialNumber *string, reason string, blacklistedBy uuid.UUID, expiresAt *string) error
	GetBlacklistedDeviceByIMEI(ctx context.Context, imei string) (*shared.BlacklistedDevice, error)
	GetBlacklistedDeviceBySerialNumber(ctx context.Context, serialNumber string) (*shared.BlacklistedDevice, error)
	CheckDeviceBlacklistStatus(ctx context.Context, imei string, serialNumber *string) (bool, *shared.BlacklistedDevice, error)
	DeactivateBlacklistedDevice(ctx context.Context, imei string) error

	// User behavior profile operations
	UpdateUserBehaviorProfile(ctx context.Context, userID uuid.UUID, behaviorData map[string]interface{}) error
	GetUserBehaviorProfile(ctx context.Context, userID uuid.UUID) (*shared.UserBehaviorProfile, error)
	AnalyzeUserBehavior(ctx context.Context, userID uuid.UUID) (*shared.FraudDetection, error)

	// Fraud network operations
	CreateFraudNetwork(ctx context.Context, network *shared.FraudNetwork) error
	GetFraudNetworkByID(ctx context.Context, id uuid.UUID) (*shared.FraudNetwork, error)
	AddNetworkMember(ctx context.Context, networkID uuid.UUID, member *shared.FraudNetworkMember) error
	GetNetworkMembers(ctx context.Context, networkID uuid.UUID) ([]*shared.FraudNetworkMember, error)

	// Fraud analysis for claims
	AnalyzeClaimFraudRisk(ctx context.Context, claimID uuid.UUID, claimType string) (*shared.FraudDetection, error)
	AnalyzeDeviceFraudRisk(ctx context.Context, deviceID uuid.UUID) (*shared.FraudDetection, error)
	AnalyzeUserFraudRisk(ctx context.Context, userID uuid.UUID) (*shared.FraudDetection, error)

	// Fraud statistics
	GetFraudStatistics(ctx context.Context, startDate, endDate string) (map[string]interface{}, error)
	GetFraudTrends(ctx context.Context, startDate, endDate string, interval string) (map[string]interface{}, error)
	GetTopFraudIndicators(ctx context.Context, limit int, startDate, endDate string) ([]map[string]interface{}, error)
}
