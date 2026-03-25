package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// WarrantyRepository defines the interface for warranty persistence operations
type WarrantyRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, warranty *shared.Warranty) error
	GetByID(ctx context.Context, id uuid.UUID) (*shared.Warranty, error)
	GetByWarrantyNumber(ctx context.Context, warrantyNumber string) (*shared.Warranty, error)
	Update(ctx context.Context, warranty *shared.Warranty) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.Warranty, int64, error)
	GetByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*shared.Warranty, error)
	GetByPolicyID(ctx context.Context, policyID uuid.UUID) ([]*shared.Warranty, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.Warranty, int64, error)
	GetActiveWarranties(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*shared.Warranty, int64, error)
	GetExpiredWarranties(ctx context.Context, beforeDate time.Time, limit, offset int) ([]*shared.Warranty, int64, error)
	GetExpiringWarranties(ctx context.Context, days int, limit, offset int) ([]*shared.Warranty, int64, error)
	Search(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*shared.Warranty, int64, error)

	// Warranty claim operations
	CreateWarrantyClaim(ctx context.Context, claim *shared.WarrantyClaim) error
	GetWarrantyClaimByID(ctx context.Context, id uuid.UUID) (*shared.WarrantyClaim, error)
	GetWarrantyClaimsByWarrantyID(ctx context.Context, warrantyID uuid.UUID, limit, offset int) ([]*shared.WarrantyClaim, int64, error)
	GetWarrantyClaimsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.WarrantyClaim, int64, error)
	UpdateWarrantyClaim(ctx context.Context, claim *shared.WarrantyClaim) error

	// Warranty provider operations
	CreateWarrantyProvider(ctx context.Context, provider *shared.WarrantyProvider) error
	GetWarrantyProviderByID(ctx context.Context, id uuid.UUID) (*shared.WarrantyProvider, error)
	GetWarrantyProviders(ctx context.Context, limit, offset int) ([]*shared.WarrantyProvider, int64, error)
	UpdateWarrantyProvider(ctx context.Context, provider *shared.WarrantyProvider) error

	// Statistics
	GetWarrantyStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetActiveWarrantyCount(ctx context.Context, userID *uuid.UUID) (int64, error)
	GetClaimRate(ctx context.Context, warrantyID uuid.UUID) (float64, error)
}
