package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceMarketplace represents a device listing in the marketplace
type DeviceMarketplace struct {
	database.BaseModel
	DeviceID      uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	SellerID      uuid.UUID `gorm:"type:uuid;not null" json:"seller_id"`
	ListingStatus string    `gorm:"type:varchar(50);default:'draft'" json:"listing_status"` // draft, active, pending, sold, expired, removed
	ListingType   string    `gorm:"type:varchar(50)" json:"listing_type"`                   // sale, auction, trade, wanted

	// Listing Details
	Title          string     `gorm:"not null" json:"title"`
	Description    string     `json:"description"`
	Condition      string     `json:"condition"` // new, like_new, excellent, good, fair, poor
	ListingDate    time.Time  `gorm:"autoCreateTime" json:"listing_date"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	ViewCount      int        `json:"view_count"`
	WatchlistCount int        `json:"watchlist_count"`

	// Pricing
	ListingPrice    float64 `json:"listing_price"`
	OriginalPrice   float64 `json:"original_price"`
	MinimumPrice    float64 `json:"minimum_price"` // Reserve price
	CurrentBid      float64 `json:"current_bid"`   // For auctions
	BuyNowPrice     float64 `json:"buy_now_price"`
	ShippingCost    float64 `json:"shipping_cost"`
	Currency        string  `gorm:"default:'USD'" json:"currency"`
	PriceNegotiable bool    `gorm:"default:false" json:"price_negotiable"`

	// Platform & Visibility
	Platform          string     `json:"platform"` // internal, ebay, amazon, facebook, craigslist
	PlatformListingID string     `json:"platform_listing_id"`
	Featured          bool       `gorm:"default:false" json:"featured"`
	Promoted          bool       `gorm:"default:false" json:"promoted"`
	PromotionEndDate  *time.Time `json:"promotion_end_date"`

	// Seller Information
	SellerType     string  `json:"seller_type"`   // individual, business, certified
	SellerRating   float64 `json:"seller_rating"` // 0-5
	SellerVerified bool    `gorm:"default:false" json:"seller_verified"`
	BusinessName   string  `json:"business_name"`
	ResponseTime   int     `json:"response_time"` // in hours

	// Verification & Trust
	DeviceVerified     bool       `gorm:"default:false" json:"device_verified"`
	VerificationDate   *time.Time `json:"verification_date"`
	VerificationReport string     `json:"verification_report"` // URL or document ID
	AuthenticityProof  string     `json:"authenticity_proof"`
	OwnershipProof     string     `json:"ownership_proof"`

	// Media
	PhotoURLs    string `gorm:"type:json" json:"photo_urls"` // JSON array of photo URLs
	VideoURL     string `json:"video_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	PhotoCount   int    `json:"photo_count"`

	// Transaction Details
	SoldDate      *time.Time `json:"sold_date"`
	FinalPrice    float64    `json:"final_price"`
	BuyerID       *uuid.UUID `gorm:"type:uuid" json:"buyer_id"`
	TransactionID string     `json:"transaction_id"`
	PaymentMethod string     `json:"payment_method"`
	PaymentStatus string     `json:"payment_status"` // pending, completed, failed, refunded

	// Shipping & Delivery
	ShippingMethod    string     `json:"shipping_method"` // standard, express, pickup
	ShippingProvider  string     `json:"shipping_provider"`
	TrackingNumber    string     `json:"tracking_number"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	DeliveryAddress   string     `gorm:"type:json" json:"delivery_address"`
	ShippingStatus    string     `json:"shipping_status"` // pending, shipped, in_transit, delivered

	// Returns & Protection
	ReturnPolicy     string `json:"return_policy"` // no_returns, 7_days, 14_days, 30_days
	BuyerProtection  bool   `gorm:"default:true" json:"buyer_protection"`
	WarrantyIncluded bool   `gorm:"default:false" json:"warranty_included"`
	WarrantyDuration int    `json:"warranty_duration"` // in days

	// Commission & Fees
	PlatformCommission float64 `json:"platform_commission"`
	ListingFee         float64 `json:"listing_fee"`
	TransactionFee     float64 `json:"transaction_fee"`
	PromotionFee       float64 `json:"promotion_fee"`
	TotalFees          float64 `json:"total_fees"`
	NetEarnings        float64 `json:"net_earnings"`

	// Analytics
	ImpressionCount   int     `json:"impression_count"`
	ClickThroughRate  float64 `json:"click_through_rate"`
	ConversionRate    float64 `json:"conversion_rate"`
	AverageTimeToSell int     `json:"average_time_to_sell"` // in days

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// Seller and Buyer should be loaded via service layer using SellerID and BuyerID to avoid circular import
}

// MarketplaceOffer represents offers made on listings
type MarketplaceOffer struct {
	database.BaseModel
	ListingID   uuid.UUID  `gorm:"type:uuid;not null" json:"listing_id"`
	BuyerID     uuid.UUID  `gorm:"type:uuid;not null" json:"buyer_id"`
	OfferAmount float64    `json:"offer_amount"`
	OfferStatus string     `gorm:"type:varchar(50);default:'pending'" json:"offer_status"` // pending, accepted, rejected, countered, expired
	OfferDate   time.Time  `gorm:"autoCreateTime" json:"offer_date"`
	ExpiryDate  *time.Time `json:"expiry_date"`
	Message     string     `json:"message"`

	// Counter Offer
	IsCounterOffer bool       `gorm:"default:false" json:"is_counter_offer"`
	CounterAmount  float64    `json:"counter_amount"`
	CounterMessage string     `json:"counter_message"`
	CounterDate    *time.Time `json:"counter_date"`

	// Decision
	DecisionDate   *time.Time `json:"decision_date"`
	DecisionReason string     `json:"decision_reason"`

	// Relationships
	Listing DeviceMarketplace `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
	// Buyer should be loaded via service layer using BuyerID to avoid circular import
}

// MarketplaceReview represents reviews for marketplace transactions
type MarketplaceReview struct {
	database.BaseModel
	ListingID  uuid.UUID `gorm:"type:uuid;not null" json:"listing_id"`
	ReviewerID uuid.UUID `gorm:"type:uuid;not null" json:"reviewer_id"`
	RevieweeID uuid.UUID `gorm:"type:uuid;not null" json:"reviewee_id"`
	ReviewType string    `json:"review_type"` // seller_review, buyer_review, product_review
	Rating     float64   `json:"rating"`      // 0-5
	Title      string    `json:"title"`
	Comment    string    `json:"comment"`

	// Review Aspects
	CommunicationRating float64 `json:"communication_rating"`
	ShippingRating      float64 `json:"shipping_rating"`
	AccuracyRating      float64 `json:"accuracy_rating"`
	ValueRating         float64 `json:"value_rating"`

	// Verification
	VerifiedPurchase bool      `gorm:"default:false" json:"verified_purchase"`
	ReviewDate       time.Time `gorm:"autoCreateTime" json:"review_date"`
	Helpful          int       `json:"helpful"`
	NotHelpful       int       `json:"not_helpful"`

	// Response
	ResponseText string     `json:"response_text"`
	ResponseDate *time.Time `json:"response_date"`

	// Relationships
	Listing DeviceMarketplace `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
	// Reviewer and Reviewee should be loaded via service layer using ReviewerID and RevieweeID to avoid circular import
}

// MarketplaceWatchlist represents users watching listings
type MarketplaceWatchlist struct {
	database.BaseModel
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ListingID      uuid.UUID `gorm:"type:uuid;not null" json:"listing_id"`
	AddedDate      time.Time `gorm:"autoCreateTime" json:"added_date"`
	NotifyOnPrice  bool      `gorm:"default:true" json:"notify_on_price"`
	NotifyOnStatus bool      `gorm:"default:true" json:"notify_on_status"`
	PriceThreshold float64   `json:"price_threshold"` // Notify if price drops below

	// Relationships
	// User should be loaded via service layer using UserID to avoid circular import
	Listing DeviceMarketplace `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}

// TableName returns the table name
func (t *DeviceMarketplace) TableName() string {
	return "device_marketplace"
}

func (t *MarketplaceOffer) TableName() string {
	return "marketplace_offers"
}

func (t *MarketplaceReview) TableName() string {
	return "marketplace_reviews"
}

func (t *MarketplaceWatchlist) TableName() string {
	return "marketplace_watchlist"
}

// BeforeCreate handles pre-creation logic
func (dm *DeviceMarketplace) BeforeCreate(tx *gorm.DB) error {
	if err := dm.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// CalculateFees calculates all marketplace fees
func (dm *DeviceMarketplace) CalculateFees() {
	// Calculate commission based on final price or listing price
	price := dm.FinalPrice
	if price == 0 {
		price = dm.ListingPrice
	}

	// Platform commission (e.g., 10%)
	dm.PlatformCommission = price * 0.10

	// Add all fees
	dm.TotalFees = dm.PlatformCommission + dm.ListingFee +
		dm.TransactionFee + dm.PromotionFee

	// Calculate net earnings
	dm.NetEarnings = price - dm.TotalFees
}

// IsActive checks if listing is currently active
func (dm *DeviceMarketplace) IsActive() bool {
	if dm.ListingStatus != "active" {
		return false
	}

	if dm.ExpiryDate != nil && time.Now().After(*dm.ExpiryDate) {
		return false
	}

	return true
}

// CompleteSale marks the listing as sold
func (dm *DeviceMarketplace) CompleteSale(buyerID uuid.UUID, finalPrice float64) {
	dm.ListingStatus = "sold"
	now := time.Now()
	dm.SoldDate = &now
	dm.BuyerID = &buyerID
	dm.FinalPrice = finalPrice

	// Calculate fees and net earnings
	dm.CalculateFees()

	// Calculate time to sell
	dm.AverageTimeToSell = int(now.Sub(dm.ListingDate).Hours() / 24)
}

// IncrementView increments the view count
func (dm *DeviceMarketplace) IncrementView() {
	dm.ViewCount++
	dm.ImpressionCount++
}

// AddToWatchlist adds to watchlist count
func (dm *DeviceMarketplace) AddToWatchlist() {
	dm.WatchlistCount++
}

// RemoveFromWatchlist removes from watchlist count
func (dm *DeviceMarketplace) RemoveFromWatchlist() {
	if dm.WatchlistCount > 0 {
		dm.WatchlistCount--
	}
}

// SetFeatured marks listing as featured
func (dm *DeviceMarketplace) SetFeatured(days int) {
	dm.Featured = true
	dm.Promoted = true
	end := time.Now().AddDate(0, 0, days)
	dm.PromotionEndDate = &end
}

// ExpireListing marks the listing as expired
func (dm *DeviceMarketplace) ExpireListing() {
	dm.ListingStatus = "expired"
	now := time.Now()
	dm.ExpiryDate = &now
}

// ValidateOffer validates if an offer is acceptable
func (mo *MarketplaceOffer) ValidateOffer(minimumPrice float64) bool {
	if mo.OfferAmount < minimumPrice {
		return false
	}

	if mo.ExpiryDate != nil && time.Now().After(*mo.ExpiryDate) {
		mo.OfferStatus = "expired"
		return false
	}

	return true
}

// AcceptOffer accepts the offer
func (mo *MarketplaceOffer) AcceptOffer() {
	mo.OfferStatus = "accepted"
	now := time.Now()
	mo.DecisionDate = &now
}

// RejectOffer rejects the offer
func (mo *MarketplaceOffer) RejectOffer(reason string) {
	mo.OfferStatus = "rejected"
	now := time.Now()
	mo.DecisionDate = &now
	mo.DecisionReason = reason
}

// CounterOffer creates a counter offer
func (mo *MarketplaceOffer) CounterOffer(amount float64, message string) {
	mo.OfferStatus = "countered"
	mo.IsCounterOffer = true
	mo.CounterAmount = amount
	mo.CounterMessage = message
	now := time.Now()
	mo.CounterDate = &now
}
