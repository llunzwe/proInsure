package repositories

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models/shared"
)

// FraudRepository defines the interface for fraud detection persistence operations
type FraudRepository interface {
	// Basic CRUD operations
	CreateFraudDetection(ctx context.Context, detection *shared.FraudDetection) error
	GetFraudDetectionByID(ctx context.Context, id uuid.UUID) (*shared.FraudDetection, error)
	UpdateFraudDetection(ctx context.Context, detection *shared.FraudDetection) error
	DeleteFraudDetection(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetFraudDetectionsByClaimID(ctx context.Context, claimID uuid.UUID) ([]*shared.FraudDetection, error)
	GetFraudDetectionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetFraudDetectionsByDeviceID(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetFraudDetectionsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetFraudDetectionsByRiskLevel(ctx context.Context, riskLevel string, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetHighRiskDetections(ctx context.Context, threshold float64, limit, offset int) ([]*shared.FraudDetection, int64, error)
	GetPendingReviews(ctx context.Context, limit, offset int) ([]*shared.FraudDetection, int64, error)
	SearchFraudDetections(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*shared.FraudDetection, int64, error)

	// Fraud indicator operations
	CreateFraudIndicator(ctx context.Context, indicator *shared.FraudIndicator) error
	GetFraudIndicatorsByDetectionID(ctx context.Context, detectionID uuid.UUID) ([]*shared.FraudIndicator, error)
	UpdateFraudIndicator(ctx context.Context, indicator *shared.FraudIndicator) error

	// Blacklisted device operations
	CreateBlacklistedDevice(ctx context.Context, device *shared.BlacklistedDevice) error
	GetBlacklistedDeviceByIMEI(ctx context.Context, imei string) (*shared.BlacklistedDevice, error)
	GetBlacklistedDeviceBySerialNumber(ctx context.Context, serialNumber string) (*shared.BlacklistedDevice, error)
	GetActiveBlacklistedDevices(ctx context.Context, limit, offset int) ([]*shared.BlacklistedDevice, int64, error)
	UpdateBlacklistedDevice(ctx context.Context, device *shared.BlacklistedDevice) error
	DeactivateBlacklistedDevice(ctx context.Context, imei string) error

	// User behavior profile operations
	CreateUserBehaviorProfile(ctx context.Context, profile *shared.UserBehaviorProfile) error
	GetUserBehaviorProfileByUserID(ctx context.Context, userID uuid.UUID) (*shared.UserBehaviorProfile, error)
	UpdateUserBehaviorProfile(ctx context.Context, profile *shared.UserBehaviorProfile) error

	// Fraud network operations
	CreateFraudNetwork(ctx context.Context, network *shared.FraudNetwork) error
	GetFraudNetworkByID(ctx context.Context, id uuid.UUID) (*shared.FraudNetwork, error)
	GetFraudNetworksByRiskLevel(ctx context.Context, riskLevel string, limit, offset int) ([]*shared.FraudNetwork, int64, error)
	AddNetworkMember(ctx context.Context, networkID uuid.UUID, member *shared.FraudNetworkMember) error
	GetNetworkMembers(ctx context.Context, networkID uuid.UUID) ([]*shared.FraudNetworkMember, error)

	// Statistics and analytics
	GetFraudStatistics(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)
	GetFraudTrends(ctx context.Context, startDate, endDate time.Time, interval string) (map[string]interface{}, error)
	GetTopFraudIndicators(ctx context.Context, limit int, startDate, endDate time.Time) ([]map[string]interface{}, error)
}
