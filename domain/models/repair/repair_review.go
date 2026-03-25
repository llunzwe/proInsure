package repair

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RepairReview represents customer reviews for repair shops
type RepairReview struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	RepairShopID     uuid.UUID      `gorm:"type:uuid;not null" json:"repair_shop_id"`
	RepairBookingID  uuid.UUID      `gorm:"type:uuid;not null" json:"repair_booking_id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Rating           int            `gorm:"not null" json:"rating"` // 1-5
	QualityRating    int            `json:"quality_rating"`         // 1-5
	ServiceRating    int            `json:"service_rating"`         // 1-5
	TimelinessRating int            `json:"timeliness_rating"`      // 1-5
	Title            string         `json:"title"`
	Comment          string         `json:"comment"`
	Pros             string         `json:"pros"`
	Cons             string         `json:"cons"`
	WouldRecommend   bool           `gorm:"default:true" json:"would_recommend"`
	IsVerified       bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Note: Relationships are defined in parent models package to avoid circular dependencies
	// These include: RepairShop, RepairBooking, User
}

// TableName returns the table name for RepairReview model
func (RepairReview) TableName() string {
	return "repair_reviews"
}
