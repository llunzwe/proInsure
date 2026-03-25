package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/repair"
)

// RepairRepository defines the interface for repair persistence operations
type RepairRepository interface {
	// Basic CRUD operations for repair bookings
	CreateRepairBooking(ctx context.Context, booking *repair.RepairBooking) error
	GetRepairBookingByID(ctx context.Context, id uuid.UUID) (*repair.RepairBooking, error)
	GetRepairBookingByBookingNumber(ctx context.Context, bookingNumber string) (*repair.RepairBooking, error)
	UpdateRepairBooking(ctx context.Context, booking *repair.RepairBooking) error
	DeleteRepairBooking(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetRepairBookingsByDeviceID(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]*repair.RepairBooking, int64, error)
	GetRepairBookingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*repair.RepairBooking, int64, error)
	GetRepairBookingsByStatus(ctx context.Context, status string, limit, offset int) ([]*repair.RepairBooking, int64, error)
	GetRepairBookingsByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*repair.RepairBooking, int64, error)
	GetScheduledRepairs(ctx context.Context, scheduledDate time.Time, limit int) ([]*repair.RepairBooking, error)
	SearchRepairBookings(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*repair.RepairBooking, int64, error)

	// Repair shop operations
	CreateRepairShop(ctx context.Context, shop *models.RepairShop) error
	GetRepairShopByID(ctx context.Context, id uuid.UUID) (*models.RepairShop, error)
	GetRepairShopsByLocation(ctx context.Context, location string, limit, offset int) ([]*models.RepairShop, int64, error)
	UpdateRepairShop(ctx context.Context, shop *models.RepairShop) error
	ListRepairShops(ctx context.Context, limit, offset int) ([]*models.RepairShop, int64, error)

	// Repair status operations
	CreateStatusUpdate(ctx context.Context, status *repair.RepairStatusUpdate) error
	GetStatusUpdatesByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*repair.RepairStatusUpdate, error)

	// Repair review operations
	CreateRepairReview(ctx context.Context, review *repair.RepairReview) error
	GetRepairReviewByID(ctx context.Context, id uuid.UUID) (*repair.RepairReview, error)
	GetRepairReviewsByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*repair.RepairReview, error)
	GetRepairReviewsByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*repair.RepairReview, int64, error)

	// Repair technician operations
	CreateTechnician(ctx context.Context, technician *repair.RepairTechnician) error
	GetTechnicianByID(ctx context.Context, id uuid.UUID) (*repair.RepairTechnician, error)
	GetTechniciansByShopID(ctx context.Context, shopID uuid.UUID) ([]*repair.RepairTechnician, error)

	// Repair parts and inventory operations
	CreatePartInventory(ctx context.Context, inventory *repair.RepairPartsInventory) error
	GetPartsInventoryByShopID(ctx context.Context, shopID uuid.UUID) ([]*repair.RepairPartsInventory, error)
	UpdatePartsInventory(ctx context.Context, inventory *repair.RepairPartsInventory) error

	// Statistics
	GetRepairStatistics(ctx context.Context, shopID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetAverageRepairTime(ctx context.Context, shopID *uuid.UUID, startDate, endDate time.Time) (time.Duration, error)
}
