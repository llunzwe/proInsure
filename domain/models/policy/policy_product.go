package policy

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/types"
)

// Product represents an insurance product
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // basic, premium, enterprise
	Type        string    `json:"type"`     // damage, theft, warranty_extension, comprehensive
	Status      string    `json:"status"`   // active, inactive, deprecated
	Version     string    `json:"version"`

	// Coverage Details
	CoverageAmount    float64 `json:"coverage_amount"`
	DeductibleAmount  float64 `json:"deductible_amount"`
	MaxClaimsPerYear  int     `json:"max_claims_per_year"`
	WaitingPeriodDays int     `json:"waiting_period_days"`

	// Pricing
	BasePremium      float64 `json:"base_premium"`
	PremiumFrequency string  `json:"premium_frequency"` // monthly, quarterly, annual
	Currency         string  `json:"currency"`

	// Features & Benefits
	Features        types.JSONArray `gorm:"type:json" json:"features"`
	Benefits        types.JSONArray `gorm:"type:json" json:"benefits"`
	Exclusions      types.JSONArray `gorm:"type:json" json:"exclusions"`
	TermsConditions string          `gorm:"type:text" json:"terms_conditions"`

	// Eligibility
	MinDeviceAge   int             `json:"min_device_age"` // in days
	MaxDeviceAge   int             `json:"max_device_age"` // in days
	MinDeviceValue float64         `json:"min_device_value"`
	MaxDeviceValue float64         `json:"max_device_value"`
	EligibleBrands types.JSONArray `gorm:"type:json" json:"eligible_brands"`
	EligibleModels types.JSONArray `gorm:"type:json" json:"eligible_models"`

	// Marketing
	MarketingName    string `json:"marketing_name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	BrochureURL      string `json:"brochure_url"`

	// Customer Trust & Social Proof
	AverageRating     float64         `json:"average_rating" gorm:"default:0"`      // 4.5 out of 5
	TotalRatings      int             `json:"total_ratings" gorm:"default:0"`       // Number of reviews
	TotalPoliciesSold int             `json:"total_policies_sold" gorm:"default:0"` // Sales count
	PopularityRank    int             `json:"popularity_rank" gorm:"default:999"`   // Ranking position
	RecommendedFor    types.JSONArray `gorm:"type:json" json:"recommended_for"`     // ["Students", "Professionals"]
	TrustBadges       types.JSONArray `gorm:"type:json" json:"trust_badges"`        // ["BBB Accredited", "ISO Certified"]
	UnderwriterName   string          `json:"underwriter_name"`                     // Insurance company backing
	UnderwriterLogo   string          `json:"underwriter_logo"`                     // Underwriter logo URL

	// Service Level & Claims
	AvgClaimApprovalTime int             `json:"avg_claim_approval_time" gorm:"default:24"`  // Hours
	AvgPayoutTime        int             `json:"avg_payout_time" gorm:"default:5"`           // Days
	ClaimSuccessRate     float64         `json:"claim_success_rate" gorm:"default:0"`        // Percentage (95.5)
	SupportAvailability  string          `json:"support_availability" gorm:"default:'24/7'"` // "24/7" or "Business Hours"
	SupportChannels      types.JSONArray `gorm:"type:json" json:"support_channels"`          // ["Phone", "Chat", "Email"]

	// Geographic & Regulatory
	AvailableRegions     types.JSONArray `gorm:"type:json" json:"available_regions"`    // ["US", "CA", "UK"]
	RestrictedStates     types.JSONArray `gorm:"type:json" json:"restricted_states"`    // ["NY", "CA"] for US
	RequiresLicense      bool            `json:"requires_license" gorm:"default:false"` // Special licensing required
	RegulatoryDisclaimer string          `json:"regulatory_disclaimer"`                 // Legal compliance text

	// Promotional & Discounts
	CurrentPromotion       string          `json:"current_promotion"` // "20% off first 3 months"
	PromoCode              string          `json:"promo_code"`        // "SAVE20"
	PromoValidUntil        *time.Time      `json:"promo_valid_until"`
	FirstTimeBuyerDiscount float64         `json:"first_time_buyer_discount" gorm:"default:0"` // Percentage
	ReferralBonus          float64         `json:"referral_bonus" gorm:"default:0"`            // Dollar amount
	BulkDiscountRates      types.JSONArray `gorm:"type:json" json:"bulk_discount_rates"`       // Volume pricing tiers
	AutoPayDiscount        float64         `json:"autopay_discount" gorm:"default:0"`          // Percentage off

	// Payment Options
	AcceptedPaymentMethods types.JSONArray `gorm:"type:json" json:"accepted_payment_methods"` // ["Credit", "Debit", "PayPal"]
	SupportsInstallments   bool            `json:"supports_installments" gorm:"default:false"`
	MinimumDownPayment     float64         `json:"minimum_down_payment" gorm:"default:0"`
	FreeTrialDays          int             `json:"free_trial_days" gorm:"default:0"`

	// Comparison & Differentiation
	ComparisonHighlights types.JSONArray `gorm:"type:json" json:"comparison_highlights"` // Key selling points
	UniqueSellingPoints  types.JSONArray `gorm:"type:json" json:"unique_selling_points"` // What makes this unique
	BestForScenarios     types.JSONArray `gorm:"type:json" json:"best_for_scenarios"`    // Ideal customer scenarios
	CompetitorComparison types.JSONArray `gorm:"type:json" json:"competitor_comparison"` // vs. other products

	// Add-ons & Upgrades
	AvailableAddons types.JSONArray `gorm:"type:json" json:"available_addons"` // Additional purchases
	UpgradePaths    types.JSONArray `gorm:"type:json" json:"upgrade_paths"`    // Upgrade options
	BundleOptions   types.JSONArray `gorm:"type:json" json:"bundle_options"`   // Compatible bundles

	// Educational Content
	FAQs                 types.JSONArray `gorm:"type:json" json:"faqs"`                   // Common questions
	VideoTutorialURL     string          `json:"video_tutorial_url"`                      // How-to video
	SampleClaimScenarios types.JSONArray `gorm:"type:json" json:"sample_claim_scenarios"` // Real examples
	RequiredDocuments    types.JSONArray `gorm:"type:json" json:"required_documents"`     // Purchase requirements

	// Dynamic Pricing
	DynamicPricingEnabled bool            `json:"dynamic_pricing_enabled" gorm:"default:false"` // Variable pricing
	PriceFactors          types.JSONArray `gorm:"type:json" json:"price_factors"`               // Pricing variables
	PriceRange            string          `json:"price_range"`                                  // "$9.99 - $29.99"
	PriceMatchGuarantee   bool            `json:"price_match_guarantee" gorm:"default:false"`

	// Sustainability & CSR
	GreenCertified      bool    `json:"green_certified" gorm:"default:false"` // Eco-friendly
	CarbonNeutral       bool    `json:"carbon_neutral" gorm:"default:false"`
	CharityContribution float64 `json:"charity_contribution" gorm:"default:0"`  // Percentage donated
	RecyclingProgram    bool    `json:"recycling_program" gorm:"default:false"` // Device recycling included

	// Display Configuration
	ShowPriceComparison bool   `json:"show_price_comparison" gorm:"default:false"` // Compare with competitors
	ShowSavingsAmount   bool   `json:"show_savings_amount" gorm:"default:false"`   // Display savings
	HighlightColor      string `json:"highlight_color"`                            // Promotional banner color
	BadgeText           string `json:"badge_text"`                                 // "LIMITED TIME", "NEW"
	DisplayPriority     int    `json:"display_priority" gorm:"default:0"`          // Special sorting priority

	// Metadata
	LaunchDate   time.Time       `json:"launch_date"`
	EndDate      *time.Time      `json:"end_date"`
	IsPromoted   bool            `json:"is_promoted"`
	DisplayOrder int             `json:"display_order"`
	Tags         types.JSONArray `gorm:"type:json" json:"tags"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ProductTiers []ProductTier `gorm:"foreignKey:ProductID" json:"product_tiers,omitempty"`
	// Removed Policies to prevent circular reference - use Product.GetPolicies() method instead
}

// ProductTier represents different tiers of a product
type ProductTier struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	TierCode  string    `gorm:"uniqueIndex;not null" json:"tier_code"` // basic, standard, premium, platinum
	TierName  string    `gorm:"not null" json:"tier_name"`
	TierLevel int       `json:"tier_level"` // 1=Basic, 2=Standard, 3=Premium, 4=Platinum

	// Coverage Adjustments
	CoverageMultiplier   float64 `json:"coverage_multiplier"`
	AdditionalCoverage   float64 `json:"additional_coverage"`
	DeductibleAdjustment float64 `json:"deductible_adjustment"`

	// Pricing
	PremiumMultiplier float64 `json:"premium_multiplier"`
	FixedPremium      float64 `json:"fixed_premium"`
	DiscountPercent   float64 `json:"discount_percent"`

	// Features
	AdditionalFeatures types.JSONArray `gorm:"type:json" json:"additional_features"`
	AdditionalBenefits types.JSONArray `gorm:"type:json" json:"additional_benefits"`

	IsDefault bool      `json:"is_default"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// ProductBundle represents bundled products
type ProductBundle struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BundleCode  string    `gorm:"uniqueIndex;not null" json:"bundle_code"`
	BundleName  string    `gorm:"not null" json:"bundle_name"`
	Description string    `json:"description"`

	// Bundle Details
	Products        types.JSONArray `gorm:"type:json" json:"products"` // Array of product IDs
	TotalPrice      float64         `json:"total_price"`
	DiscountAmount  float64         `json:"discount_amount"`
	DiscountPercent float64         `json:"discount_percent"`

	// Validity
	ValidFrom   time.Time  `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	MaxQuantity int        `json:"max_quantity"`
	MinQuantity int        `json:"min_quantity"`

	Status    string         `json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// REMOVED: Coverage entity - Use PolicyCoverage instead to avoid duplication
// PolicyCoverage in policy_coverage.go provides all coverage functionality

// JSONArray type has been moved to internal/domain/types/json_types.go to avoid circular dependencies

// ============================================
// PRODUCT BUSINESS METHODS
// ============================================

// IsAvailable checks if the product is currently available for purchase
func (p *Product) IsAvailable() bool {
	now := time.Now()
	return p.Status == "active" &&
		now.After(p.LaunchDate) &&
		(p.EndDate == nil || now.Before(*p.EndDate))
}

// GetEffectivePrice calculates the price based on customer type and applicable discounts
func (p *Product) GetEffectivePrice(customerType string, hasAutoPay bool) float64 {
	price := p.BasePremium

	// Apply first-time buyer discount
	if customerType == "new" && p.FirstTimeBuyerDiscount > 0 {
		price *= (1 - p.FirstTimeBuyerDiscount/100)
	}

	// Apply autopay discount
	if hasAutoPay && p.AutoPayDiscount > 0 {
		price *= (1 - p.AutoPayDiscount/100)
	}

	return price
}

// IsEligibleForDevice checks if a device meets all eligibility criteria
func (p *Product) IsEligibleForDevice(deviceAge int, deviceValue float64, deviceBrand string, deviceModel string) bool {
	// Check age requirements
	if deviceAge < p.MinDeviceAge || deviceAge > p.MaxDeviceAge {
		return false
	}

	// Check value requirements
	if deviceValue < p.MinDeviceValue || deviceValue > p.MaxDeviceValue {
		return false
	}

	// Check brand eligibility if specified
	if len(p.EligibleBrands) > 0 {
		brandFound := false
		for _, brand := range p.EligibleBrands {
			if brandStr, ok := brand.(string); ok && brandStr == deviceBrand {
				brandFound = true
				break
			}
		}
		if !brandFound {
			return false
		}
	}

	// Check model eligibility if specified
	if len(p.EligibleModels) > 0 {
		modelFound := false
		for _, model := range p.EligibleModels {
			if modelStr, ok := model.(string); ok && modelStr == deviceModel {
				modelFound = true
				break
			}
		}
		if !modelFound {
			return false
		}
	}

	return true
}

// GetPopularityBadge returns a badge based on sales volume
func (p *Product) GetPopularityBadge() string {
	if p.TotalPoliciesSold > 10000 {
		return "Best Seller"
	} else if p.TotalPoliciesSold > 5000 {
		return "Popular Choice"
	} else if p.TotalPoliciesSold > 1000 {
		return "Customer Favorite"
	} else if p.TotalPoliciesSold > 100 {
		return "Trending"
	}
	return ""
}

// GetRatingDisplay returns a formatted rating display
func (p *Product) GetRatingDisplay() string {
	if p.TotalRatings == 0 {
		return "No ratings yet"
	}
	return fmt.Sprintf("%.1f ⭐ (%d reviews)", p.AverageRating, p.TotalRatings)
}

// HasActivePromotion checks if there's a current promotion
func (p *Product) HasActivePromotion() bool {
	if p.CurrentPromotion == "" {
		return false
	}
	if p.PromoValidUntil == nil {
		return true // No expiry
	}
	return time.Now().Before(*p.PromoValidUntil)
}

// GetBulkDiscount calculates discount based on quantity
func (p *Product) GetBulkDiscount(quantity int) float64 {
	if len(p.BulkDiscountRates) == 0 {
		return 0
	}

	// Bulk discount rates expected format: [{"min": 2, "max": 5, "discount": 10}, ...]
	for _, rate := range p.BulkDiscountRates {
		if rateMap, ok := rate.(map[string]interface{}); ok {
			min := int(rateMap["min"].(float64))
			max := int(rateMap["max"].(float64))
			discount := rateMap["discount"].(float64)

			if quantity >= min && quantity <= max {
				return discount
			}
		}
	}

	return 0
}

// IsAvailableInRegion checks if product is available in a specific region
func (p *Product) IsAvailableInRegion(region string, state string) bool {
	// Check if region is in available regions
	regionFound := false
	for _, r := range p.AvailableRegions {
		if regionStr, ok := r.(string); ok && regionStr == region {
			regionFound = true
			break
		}
	}

	if !regionFound {
		return false
	}

	// Check if state is restricted (for US)
	if region == "US" && state != "" {
		for _, s := range p.RestrictedStates {
			if stateStr, ok := s.(string); ok && stateStr == state {
				return false // State is restricted
			}
		}
	}

	return true
}

// GetTierByCode returns the ProductTier for a given tier code
func (p *Product) GetTierByCode(tierCode string) *ProductTier {
	for _, tier := range p.ProductTiers {
		if tier.TierCode == tierCode && tier.IsActive {
			return &tier
		}
	}
	return nil
}

// GetDefaultTier returns the default product tier
func (p *Product) GetDefaultTier() *ProductTier {
	for _, tier := range p.ProductTiers {
		if tier.IsDefault && tier.IsActive {
			return &tier
		}
	}
	// If no default, return first active tier
	for _, tier := range p.ProductTiers {
		if tier.IsActive {
			return &tier
		}
	}
	return nil
}

// CalculatePremiumForTier calculates the premium for a specific tier
func (p *Product) CalculatePremiumForTier(tierCode string) float64 {
	tier := p.GetTierByCode(tierCode)
	if tier == nil {
		return p.BasePremium
	}

	premium := p.BasePremium

	// Apply tier multiplier
	if tier.PremiumMultiplier > 0 {
		premium *= tier.PremiumMultiplier
	}

	// Or use fixed premium if set
	if tier.FixedPremium > 0 {
		premium = tier.FixedPremium
	}

	// Apply tier discount
	if tier.DiscountPercent > 0 {
		premium *= (1 - tier.DiscountPercent/100)
	}

	return premium
}

// GetSavingsMessage returns a formatted savings message for marketing
func (p *Product) GetSavingsMessage() string {
	if !p.ShowSavingsAmount {
		return ""
	}

	if p.FirstTimeBuyerDiscount > 0 {
		savings := p.BasePremium * (p.FirstTimeBuyerDiscount / 100)
		return fmt.Sprintf("Save $%.2f as a new customer!", savings)
	}

	if p.HasActivePromotion() {
		return p.CurrentPromotion
	}

	if p.AutoPayDiscount > 0 {
		savings := p.BasePremium * (p.AutoPayDiscount / 100)
		return fmt.Sprintf("Save $%.2f/month with AutoPay", savings)
	}

	return ""
}

// GetDisplayBadge returns the appropriate badge for display
func (p *Product) GetDisplayBadge() string {
	// Priority order for badges
	if p.BadgeText != "" {
		return p.BadgeText
	}

	if p.HasActivePromotion() {
		return "LIMITED OFFER"
	}

	if popularity := p.GetPopularityBadge(); popularity != "" {
		return popularity
	}

	if p.IsPromoted {
		return "FEATURED"
	}

	// Check if new (launched within last 30 days)
	if time.Since(p.LaunchDate).Hours()/24 <= 30 {
		return "NEW"
	}

	if p.GreenCertified {
		return "ECO-FRIENDLY"
	}

	return ""
}

// GetPolicies retrieves all policies for this product (prevents circular reference)
// Returns interface{} to avoid import cycle - cast to []models.Policy in calling code
func (p *Product) GetPolicies(db *gorm.DB) ([]interface{}, error) {
	var policies []interface{} // []models.Policy
	err := db.Where("product_id = ?", p.ID).Find(&policies).Error
	return policies, err
}

// GetActivePoliciesCount returns the count of active policies for this product
func (p *Product) GetActivePoliciesCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Table("policies").
		Where("product_id = ? AND status = ?", p.ID, "active").
		Count(&count).Error
	return count, err
}
