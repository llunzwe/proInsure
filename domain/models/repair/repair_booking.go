package repair

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RepairBooking represents a repair booking/appointment
// This is kept in the repair package to maintain domain boundaries
type RepairBooking struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	BookingNumber        string         `gorm:"uniqueIndex;not null" json:"booking_number"`
	ClaimID              *uuid.UUID     `gorm:"type:uuid" json:"claim_id"`
	UserID               uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	DeviceID             uuid.UUID      `gorm:"type:uuid;not null" json:"device_id"`
	RepairShopID         uuid.UUID      `gorm:"type:uuid;not null" json:"repair_shop_id"`
	RepairType           string         `gorm:"not null" json:"repair_type"`                     // screen, battery, camera, water_damage, etc.
	ServiceType          string         `gorm:"not null;default:'in_store'" json:"service_type"` // in_store, mobile, pickup, mail_in, etc.
	IssueDescription     string         `gorm:"not null" json:"issue_description"`
	EstimatedCost        float64        `json:"estimated_cost"`
	ActualCost           float64        `json:"actual_cost"`
	DiagnosticFee        float64        `json:"diagnostic_fee"`     // Fee charged for diagnosis
	DiscountAmount       float64        `json:"discount_amount"`    // Discount applied to the repair
	EstimatedDuration    int            `json:"estimated_duration"` // hours
	ActualDuration       int            `json:"actual_duration"`    // hours
	ScheduledDate        time.Time      `gorm:"not null" json:"scheduled_date"`
	CompletedDate        *time.Time     `json:"completed_date"`
	ActualCompletionTime *time.Time     `json:"actual_completion_time"`                     // Actual time when repair was completed
	Status               string         `gorm:"not null;default:'scheduled'" json:"status"` // scheduled, in_progress, completed, cancelled, delayed
	Priority             string         `gorm:"default:'normal'" json:"priority"`           // low, normal, high, urgent
	PartsRequired        string         `json:"parts_required"`                             // JSON array
	PartsAvailability    string         `json:"parts_availability"`                         // JSON object
	TechnicianID         *uuid.UUID     `gorm:"type:uuid" json:"technician_id"`
	TechnicianNotes      string         `json:"technician_notes"`
	CustomerNotes        string         `json:"customer_notes"`
	QualityCheckPassed   bool           `gorm:"default:false" json:"quality_check_passed"`
	WarrantyStartDate    *time.Time     `json:"warranty_start_date"`
	WarrantyEndDate      *time.Time     `json:"warranty_end_date"`
	IsUnderWarranty      bool           `gorm:"default:false" json:"is_under_warranty"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Note: Relationships are defined in the parent models package to avoid circular dependencies
	// These include: Claim, User, Device, RepairShop, PartsUsed, AccessoriesReplaced, TechnicianAssignments
}

// RepairStatusUpdate represents status updates for repair bookings
type RepairStatusUpdate struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RepairBookingID     uuid.UUID  `gorm:"type:uuid;not null" json:"repair_booking_id"`
	Status              string     `gorm:"not null" json:"status"`
	Message             string     `json:"message"`
	UpdatedBy           uuid.UUID  `gorm:"type:uuid;not null" json:"updated_by"`
	UpdatedByRole       string     `json:"updated_by_role"` // customer, technician, admin
	EstimatedCompletion *time.Time `json:"estimated_completion"`
	Photos              string     `json:"photos"` // JSON array of photo URLs
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Note: RepairBooking relationship is defined in parent models package
}

// TableName methods
func (RepairBooking) TableName() string {
	return "repair_bookings"
}

func (RepairStatusUpdate) TableName() string {
	return "repair_status_updates"
}
