package policy

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/types"
)

// Quote represents an insurance quote
type Quote struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuoteNumber string    `gorm:"uniqueIndex;not null" json:"quote_number"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	DeviceID    uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`

	// Quote Details
	Status  string `json:"status"`  // draft, pending, approved, expired, converted
	Type    string `json:"type"`    // new, renewal, modification
	Channel string `json:"channel"` // web, mobile, agent, partner

	// Device Information
	DeviceMake  string  `json:"device_make"`
	DeviceModel string  `json:"device_model"`
	DeviceValue float64 `json:"device_value"`
	DeviceAge   int     `json:"device_age"` // in days
	IMEI        string  `json:"imei"`

	// Coverage
	CoverageType     string  `json:"coverage_type"`
	CoverageAmount   float64 `json:"coverage_amount"`
	DeductibleAmount float64 `json:"deductible_amount"`
	PolicyTerm       int     `json:"policy_term"` // in months

	// Pricing
	BasePremium      float64 `json:"base_premium"`
	TaxAmount        float64 `json:"tax_amount"`
	DiscountAmount   float64 `json:"discount_amount"`
	TotalPremium     float64 `json:"total_premium"`
	Currency         string  `json:"currency"`
	PaymentFrequency string  `json:"payment_frequency"`

	// Risk Assessment
	RiskScore         float64 `json:"risk_score"`
	RiskCategory      string  `json:"risk_category"`
	UnderwritingNotes string  `json:"underwriting_notes"`

	// Validity
	ValidFrom         time.Time  `json:"valid_from"`
	ValidUntil        time.Time  `json:"valid_until"`
	ConvertedAt       *time.Time `json:"converted_at"`
	ConvertedPolicyID *uuid.UUID `gorm:"type:uuid" json:"converted_policy_id"`

	// Sales Information
	AgentID      *uuid.UUID `gorm:"type:uuid" json:"agent_id"`
	ReferralCode string     `json:"referral_code"`
	PromoCode    string     `json:"promo_code"`
	Source       string     `json:"source"`
	Campaign     string     `json:"campaign"`

	// Metadata
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	SessionID string `json:"session_id"`
	Metadata  string `gorm:"type:json" json:"metadata"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	// User should be loaded via service layer using UserID to avoid circular import
	// Device should be loaded via service layer using DeviceID to avoid circular import
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	// Agent should be loaded via service layer using AgentID to avoid circular import
	// ConvertedPolicy should be loaded via service layer using ConvertedPolicyID to avoid circular import
	QuoteItems   []QuoteItem    `gorm:"foreignKey:QuoteID" json:"quote_items,omitempty"`
	QuoteHistory []QuoteHistory `gorm:"foreignKey:QuoteID" json:"quote_history,omitempty"`
}

// ============================================
// QUOTE BUSINESS METHODS
// ============================================

// IsValid checks if the quote is still valid
func (q *Quote) IsValid() bool {
	now := time.Now()
	return q.Status != "expired" &&
		q.Status != "converted" &&
		now.After(q.ValidFrom) &&
		now.Before(q.ValidUntil)
}

// CanConvert checks if quote can be converted to policy
func (q *Quote) CanConvert() bool {
	return q.IsValid() && q.Status == "approved"
}

// DaysUntilExpiry returns days until quote expires
func (q *Quote) DaysUntilExpiry() int {
	if q.Status == "expired" || q.Status == "converted" {
		return 0
	}
	duration := time.Until(q.ValidUntil)
	return int(duration.Hours() / 24)
}

// IsExpired checks if quote has expired
func (q *Quote) IsExpired() bool {
	return q.Status == "expired" || time.Now().After(q.ValidUntil)
}

// MarkAsExpired updates quote status to expired
func (q *Quote) MarkAsExpired() {
	q.Status = "expired"
}

// ConvertToPolicy marks quote as converted and sets policy reference
func (q *Quote) ConvertToPolicy(policyID uuid.UUID) error {
	if !q.CanConvert() {
		return fmt.Errorf("quote cannot be converted: invalid status or expired")
	}

	now := time.Now()
	q.Status = "converted"
	q.ConvertedAt = &now
	q.ConvertedPolicyID = &policyID
	return nil
}

// GetConvertedPolicyID returns the ID of the converted policy if it exists
// The actual Policy entity should be loaded via the service layer to avoid circular import
func (q *Quote) GetConvertedPolicyID() *uuid.UUID {
	return q.ConvertedPolicyID
}

// CalculateTotalPrice calculates total from all selected quote items
func (q *Quote) CalculateTotalPrice() float64 {
	total := q.BasePremium

	for _, item := range q.QuoteItems {
		if item.IsSelected || item.IsMandatory {
			total += item.FinalPrice * float64(item.Quantity)
		}
	}

	total += q.TaxAmount
	total -= q.DiscountAmount

	return total
}

// GetEffectivePremium returns the final premium amount
func (q *Quote) GetEffectivePremium() float64 {
	return q.TotalPremium
}

// IsHighRisk checks if quote is high risk
func (q *Quote) IsHighRisk() bool {
	return q.RiskScore > 75 || q.RiskCategory == "high" || q.RiskCategory == "very_high"
}

// RequiresApproval checks if quote needs manual approval
func (q *Quote) RequiresApproval() bool {
	return q.IsHighRisk() ||
		q.CoverageAmount > 5000 ||
		q.Status == "pending"
}

// Approve approves the quote for conversion
func (q *Quote) Approve(approverID uuid.UUID, notes string) error {
	if q.Status != "pending" {
		return fmt.Errorf("only pending quotes can be approved")
	}

	q.Status = "approved"
	q.UnderwritingNotes = notes
	return nil
}

// AddQuoteItem adds an item to the quote
func (q *Quote) AddQuoteItem(item QuoteItem) {
	item.QuoteID = q.ID
	q.QuoteItems = append(q.QuoteItems, item)
	q.TotalPremium = q.CalculateTotalPrice()
}

// ApplyDiscount applies a discount to the quote
func (q *Quote) ApplyDiscount(discount Discount) error {
	if !discount.IsValidForQuote(q) {
		return fmt.Errorf("discount not valid for this quote")
	}

	if discount.Type == "percentage" {
		q.DiscountAmount = q.BasePremium * (discount.DiscountPercent / 100)
	} else if discount.Type == "fixed" {
		q.DiscountAmount = discount.DiscountAmount
	}

	// Apply max discount limit
	if discount.MaxDiscountAmount > 0 && q.DiscountAmount > discount.MaxDiscountAmount {
		q.DiscountAmount = discount.MaxDiscountAmount
	}

	q.PromoCode = discount.Code
	q.TotalPremium = q.CalculateTotalPrice()
	return nil
}

// GetQuoteAge returns the age of the quote in days
func (q *Quote) GetQuoteAge() int {
	duration := time.Since(q.CreatedAt)
	return int(duration.Hours() / 24)
}

// NeedsFollowUp checks if quote needs follow-up
func (q *Quote) NeedsFollowUp() bool {
	if q.Status == "converted" || q.Status == "expired" {
		return false
	}

	// Follow up if quote is older than 3 days and still valid
	return q.GetQuoteAge() > 3 && q.IsValid()
}

// GetConversionRate helper for analytics
func GetQuoteConversionRate(db *gorm.DB, startDate, endDate time.Time) (float64, error) {
	var totalQuotes int64
	var convertedQuotes int64

	db.Model(&Quote{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalQuotes)

	db.Model(&Quote{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", "converted").
		Count(&convertedQuotes)

	if totalQuotes == 0 {
		return 0, nil
	}

	return float64(convertedQuotes) / float64(totalQuotes) * 100, nil
}

// QuoteItem represents individual items in a quote
type QuoteItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuoteID     uuid.UUID `gorm:"type:uuid;not null" json:"quote_id"`
	ItemType    string    `json:"item_type"` // coverage, addon, fee
	ItemCode    string    `json:"item_code"`
	ItemName    string    `json:"item_name"`
	Description string    `json:"description"`

	// Pricing
	BasePrice  float64 `json:"base_price"`
	Discount   float64 `json:"discount"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`

	Quantity    int  `json:"quantity"`
	IsMandatory bool `json:"is_mandatory"`
	IsSelected  bool `json:"is_selected"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Quote *Quote `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
}

// GetTotalPrice calculates the total price for the quote item
func (qi *QuoteItem) GetTotalPrice() float64 {
	return qi.FinalPrice * float64(qi.Quantity)
}

// IsOptional checks if the item is optional
func (qi *QuoteItem) IsOptional() bool {
	return !qi.IsMandatory
}

// QuoteHistory tracks changes to quotes
type QuoteHistory struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuoteID        uuid.UUID `gorm:"type:uuid;not null" json:"quote_id"`
	Action         string    `json:"action"` // created, modified, approved, rejected, expired, converted
	PreviousStatus string    `json:"previous_status"`
	NewStatus      string    `json:"new_status"`
	Changes        string    `gorm:"type:json" json:"changes"`
	PerformedBy    uuid.UUID `gorm:"type:uuid" json:"performed_by"`
	Notes          string    `json:"notes"`
	Timestamp      time.Time `json:"timestamp"`

	// Relationships
	Quote *Quote `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	// User should be loaded via service layer using PerformedBy to avoid circular import
}

// Discount represents discount rules and applications
type Discount struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // percentage, fixed, bundle, referral, loyalty

	// Discount Values
	DiscountPercent   float64 `json:"discount_percent"`
	DiscountAmount    float64 `json:"discount_amount"`
	MaxDiscountAmount float64 `json:"max_discount_amount"`

	// Applicability
	ApplicableProducts types.JSONArray `gorm:"type:json" json:"applicable_products"`
	MinPurchaseAmount  float64         `json:"min_purchase_amount"`
	MaxUsageCount      int             `json:"max_usage_count"`
	MaxUsagePerUser    int             `json:"max_usage_per_user"`

	// Validity
	ValidFrom        time.Time  `json:"valid_from"`
	ValidUntil       *time.Time `json:"valid_until"`
	IsActive         bool       `json:"is_active"`
	RequiresApproval bool       `json:"requires_approval"`

	// Usage Tracking
	UsageCount         int     `json:"usage_count"`
	TotalDiscountGiven float64 `json:"total_discount_given"`

	// Rules
	EligibilityCriteria string `gorm:"type:json" json:"eligibility_criteria"`
	ExclusionRules      string `gorm:"type:json" json:"exclusion_rules"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ============================================
// DISCOUNT BUSINESS METHODS
// ============================================

// IsValidForQuote checks if discount is valid for a specific quote
func (d *Discount) IsValidForQuote(quote *Quote) bool {
	now := time.Now()

	// Check if discount is active
	if !d.IsActive {
		return false
	}

	// Check validity period
	if now.Before(d.ValidFrom) || (d.ValidUntil != nil && now.After(*d.ValidUntil)) {
		return false
	}

	// Check minimum purchase amount
	if d.MinPurchaseAmount > 0 && quote.TotalPremium < d.MinPurchaseAmount {
		return false
	}

	// Check product applicability
	if len(d.ApplicableProducts) > 0 {
		productFound := false
		for _, pid := range d.ApplicableProducts {
			if pidStr, ok := pid.(string); ok {
				productID, _ := uuid.Parse(pidStr)
				if productID == quote.ProductID {
					productFound = true
					break
				}
			}
		}
		if !productFound {
			return false
		}
	}

	// Check usage limits
	if d.MaxUsageCount > 0 && d.UsageCount >= d.MaxUsageCount {
		return false
	}

	return true
}

// CalculateDiscountAmount calculates the discount amount for a given base amount
func (d *Discount) CalculateDiscountAmount(baseAmount float64) float64 {
	var discountAmount float64

	if d.Type == "percentage" {
		discountAmount = baseAmount * (d.DiscountPercent / 100)
	} else if d.Type == "fixed" {
		discountAmount = d.DiscountAmount
	}

	// Apply max discount limit
	if d.MaxDiscountAmount > 0 && discountAmount > d.MaxDiscountAmount {
		discountAmount = d.MaxDiscountAmount
	}

	return discountAmount
}

// IncrementUsage increments the usage counter
func (d *Discount) IncrementUsage(discountGiven float64) {
	d.UsageCount++
	d.TotalDiscountGiven += discountGiven
}

// IsExpired checks if discount has expired
func (d *Discount) IsExpired() bool {
	if d.ValidUntil == nil {
		return false
	}
	return time.Now().After(*d.ValidUntil)
}

// GetRemainingUses returns remaining uses for the discount
func (d *Discount) GetRemainingUses() int {
	if d.MaxUsageCount <= 0 {
		return -1 // Unlimited
	}
	return d.MaxUsageCount - d.UsageCount
}

// SalesLead represents potential customers
type SalesLead struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	LeadNumber string    `gorm:"uniqueIndex;not null" json:"lead_number"`

	// Contact Information
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`

	// Lead Details
	Source   string `json:"source"`   // website, campaign, referral, partner
	Channel  string `json:"channel"`  // online, phone, agent, partner
	Status   string `json:"status"`   // new, contacted, qualified, converted, lost
	Priority string `json:"priority"` // low, medium, high, urgent
	Score    int    `json:"score"`

	// Product Interest
	InterestedProducts types.JSONArray `gorm:"type:json" json:"interested_products"`
	DeviceType         string          `json:"device_type"`
	EstimatedValue     float64         `json:"estimated_value"`

	// Assignment
	AssignedTo *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	AssignedAt *time.Time `json:"assigned_at"`

	// Conversion
	ConvertedAt       *time.Time `json:"converted_at"`
	ConvertedToUserID *uuid.UUID `gorm:"type:uuid" json:"converted_to_user_id"`
	ConversionValue   float64    `json:"conversion_value"`

	// Follow-up
	LastContactedAt *time.Time `json:"last_contacted_at"`
	NextFollowUpAt  *time.Time `json:"next_follow_up_at"`
	Notes           string     `gorm:"type:text" json:"notes"`

	// Campaign/Marketing
	CampaignID  string `json:"campaign_id"`
	UTMSource   string `json:"utm_source"`
	UTMCampaign string `json:"utm_campaign"`
	UTMMedium   string `json:"utm_medium"`
	ReferrerURL string `json:"referrer_url"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	// Agent should be loaded via service layer using AssignedTo to avoid circular import
	// ConvertedUser should be loaded via service layer using ConvertedToUserID to avoid circular import
}

// ============================================
// SALES LEAD BUSINESS METHODS
// ============================================

// IsQualified checks if lead is qualified
func (sl *SalesLead) IsQualified() bool {
	return sl.Status == "qualified" && sl.Score >= 70
}

// Convert marks the lead as converted
func (sl *SalesLead) Convert(userID uuid.UUID, value float64) {
	now := time.Now()
	sl.Status = "converted"
	sl.ConvertedAt = &now
	sl.ConvertedToUserID = &userID
	sl.ConversionValue = value
}

// NeedsFollowUp checks if lead needs follow-up
func (sl *SalesLead) NeedsFollowUp() bool {
	if sl.Status == "converted" || sl.Status == "lost" {
		return false
	}

	if sl.NextFollowUpAt != nil && time.Now().After(*sl.NextFollowUpAt) {
		return true
	}

	// Follow up if no contact in last 7 days
	if sl.LastContactedAt != nil {
		daysSinceContact := int(time.Since(*sl.LastContactedAt).Hours() / 24)
		return daysSinceContact > 7
	}

	return true // New lead needs contact
}

// SetNextFollowUp sets the next follow-up date
func (sl *SalesLead) SetNextFollowUp(days int) {
	followUpDate := time.Now().AddDate(0, 0, days)
	sl.NextFollowUpAt = &followUpDate
}

// UpdateLastContact updates the last contact timestamp
func (sl *SalesLead) UpdateLastContact() {
	now := time.Now()
	sl.LastContactedAt = &now
}

// GetLeadAge returns the age of the lead in days
func (sl *SalesLead) GetLeadAge() int {
	duration := time.Since(sl.CreatedAt)
	return int(duration.Hours() / 24)
}

// GetPriorityLevel returns numeric priority level
func (sl *SalesLead) GetPriorityLevel() int {
	switch sl.Priority {
	case "urgent":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}
