package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/repair"
)

// RepairService defines the interface for repair business logic operations
type RepairService interface {
	// Repair booking management
	CreateRepairBooking(ctx context.Context, booking *repair.RepairBooking) error
	GetRepairBookingByID(ctx context.Context, id uuid.UUID) (*repair.RepairBooking, error)
	GetRepairBookingByBookingNumber(ctx context.Context, bookingNumber string) (*repair.RepairBooking, error)
	UpdateRepairBooking(ctx context.Context, booking *repair.RepairBooking) error
	CancelRepairBooking(ctx context.Context, bookingID uuid.UUID, reason string) error

	// Repair booking queries
	GetRepairBookingsByDeviceID(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]*repair.RepairBooking, int64, error)
	GetRepairBookingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*repair.RepairBooking, int64, error)
	GetRepairBookingsByStatus(ctx context.Context, status string, limit, offset int) ([]*repair.RepairBooking, int64, error)
	GetScheduledRepairs(ctx context.Context, scheduledDate string, limit int) ([]*repair.RepairBooking, error)

	// Repair shop management
	CreateRepairShop(ctx context.Context, shop *models.RepairShop) error
	GetRepairShopByID(ctx context.Context, id uuid.UUID) (*models.RepairShop, error)
	GetRepairShopsByLocation(ctx context.Context, location string, limit, offset int) ([]*models.RepairShop, int64, error)
	UpdateRepairShop(ctx context.Context, shop *models.RepairShop) error
	ListRepairShops(ctx context.Context, limit, offset int) ([]*models.RepairShop, int64, error)

	// Repair status management
	UpdateRepairStatus(ctx context.Context, bookingID uuid.UUID, status string, notes string) error
	AddStatusUpdate(ctx context.Context, statusUpdate *repair.RepairStatusUpdate) error
	GetStatusUpdatesByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*repair.RepairStatusUpdate, error)

	// Repair review management
	CreateRepairReview(ctx context.Context, review *repair.RepairReview) error
	GetRepairReviewByID(ctx context.Context, id uuid.UUID) (*repair.RepairReview, error)
	GetRepairReviewsByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*repair.RepairReview, error)
	GetRepairReviewsByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*repair.RepairReview, int64, error)

	// Repair technician management
	AssignTechnician(ctx context.Context, bookingID uuid.UUID, technicianID uuid.UUID) error
	GetTechniciansByShopID(ctx context.Context, shopID uuid.UUID) ([]*repair.RepairTechnician, error)

	// Repair statistics
	GetRepairStatistics(ctx context.Context, shopID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetAverageRepairTime(ctx context.Context, shopID *uuid.UUID, startDate, endDate string) (int, error)
}
