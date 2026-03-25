package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// WarrantyService defines the interface for warranty business logic operations
type WarrantyService interface {
	// Warranty management
	CreateWarranty(ctx context.Context, warranty *shared.Warranty) error
	GetWarrantyByID(ctx context.Context, id uuid.UUID) (*shared.Warranty, error)
	GetWarrantyByWarrantyNumber(ctx context.Context, warrantyNumber string) (*shared.Warranty, error)
	UpdateWarranty(ctx context.Context, warranty *shared.Warranty) error
	CancelWarranty(ctx context.Context, warrantyID uuid.UUID, reason string) error

	// Warranty queries
	GetWarrantiesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.Warranty, int64, error)
	GetWarrantiesByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*shared.Warranty, error)
	GetWarrantiesByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.Warranty, int64, error)
	GetActiveWarranties(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*shared.Warranty, int64, error)
	GetExpiringWarranties(ctx context.Context, days int, limit, offset int) ([]*shared.Warranty, int64, error)

	// Warranty claim management
	CreateWarrantyClaim(ctx context.Context, claim *shared.WarrantyClaim) error
	GetWarrantyClaimByID(ctx context.Context, id uuid.UUID) (*shared.WarrantyClaim, error)
	GetWarrantyClaimsByWarrantyID(ctx context.Context, warrantyID uuid.UUID, limit, offset int) ([]*shared.WarrantyClaim, int64, error)
	UpdateWarrantyClaim(ctx context.Context, claim *shared.WarrantyClaim) error
	ProcessWarrantyClaim(ctx context.Context, claimID uuid.UUID, decision string, notes string) error

	// Warranty validation
	ValidateWarrantyCoverage(ctx context.Context, warrantyID uuid.UUID, claimType string) (bool, string, error)
	CheckWarrantyEligibility(ctx context.Context, deviceID uuid.UUID, warrantyType string) (bool, string, error)

	// Warranty transfer
	TransferWarranty(ctx context.Context, warrantyID uuid.UUID, newOwnerID uuid.UUID) error
	GetWarrantyTransferHistory(ctx context.Context, warrantyID uuid.UUID) ([]map[string]interface{}, error)

	// Warranty provider management
	GetWarrantyProviders(ctx context.Context, limit, offset int) ([]*shared.WarrantyProvider, int64, error)
	GetWarrantyProviderByID(ctx context.Context, id uuid.UUID) (*shared.WarrantyProvider, error)

	// Warranty statistics
	GetWarrantyStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetActiveWarrantyCount(ctx context.Context, userID *uuid.UUID) (int64, error)
	GetClaimRate(ctx context.Context, warrantyID uuid.UUID) (float64, error)
}
