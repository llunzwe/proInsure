package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"smartsure/internal/domain/models/accessory"
)

// Accessory represents a device accessory with comprehensive management
type Accessory struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// Embedded structs for organization
	accessory.AccessoryIdentification `gorm:"embedded;embeddedPrefix:id_" json:"identification"`
	accessory.AccessoryBasicInfo      `gorm:"embedded;embeddedPrefix:info_" json:"basic_info"`
	accessory.AccessoryCompatibility  `gorm:"embedded;embeddedPrefix:compat_" json:"compatibility"`
	accessory.AccessoryPricing        `gorm:"embedded;embeddedPrefix:price_" json:"pricing"`
	accessory.AccessoryInventory      `gorm:"embedded;embeddedPrefix:inv_" json:"inventory"`
	accessory.AccessorySupplier       `gorm:"embedded;embeddedPrefix:supp_" json:"supplier"`
	accessory.AccessoryQuality        `gorm:"embedded;embeddedPrefix:qual_" json:"quality"`
	accessory.AccessoryWarranty       `gorm:"embedded;embeddedPrefix:warr_" json:"warranty"`
	accessory.AccessoryDimensions     `gorm:"embedded;embeddedPrefix:dim_" json:"dimensions"`
	accessory.AccessoryMedia          `gorm:"embedded;embeddedPrefix:media_" json:"media"`
	accessory.AccessoryMetrics        `gorm:"embedded;embeddedPrefix:metrics_" json:"metrics"`
	accessory.AccessoryStatus         `gorm:"embedded;embeddedPrefix:status_" json:"status"`
	accessory.AccessoryAuthentication `gorm:"embedded;embeddedPrefix:auth_" json:"authentication"`
	accessory.AccessoryTechnicalSpecs `gorm:"embedded;embeddedPrefix:tech_" json:"technical_specs"`
	accessory.AccessoryMarketData     `gorm:"embedded;embeddedPrefix:market_" json:"market_data"`
	accessory.AccessoryEnvironmental  `gorm:"embedded;embeddedPrefix:env_" json:"environmental"`

	// Timestamps
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - Feature Models
	Reviews            []accessory.AccessoryReview              `gorm:"foreignKey:AccessoryID" json:"reviews,omitempty"`
	Promotions         []accessory.AccessoryPromotion           `gorm:"foreignKey:AccessoryID" json:"promotions,omitempty"`
	Movements          []accessory.AccessoryMovement            `gorm:"foreignKey:AccessoryID" json:"movements,omitempty"`
	Claims             []accessory.AccessoryClaim               `gorm:"foreignKey:AccessoryID" json:"claims,omitempty"`
	PriceHistory       []accessory.AccessoryPriceHistory        `gorm:"foreignKey:AccessoryID" json:"price_history,omitempty"`
	Compatibility      []accessory.AccessoryCompatibilityMatrix `gorm:"foreignKey:AccessoryID" json:"compatibility_matrix,omitempty"`
	Alerts             []accessory.AccessoryAlert               `gorm:"foreignKey:AccessoryID" json:"alerts,omitempty"`
	Recommendations    []accessory.AccessoryRecommendation      `gorm:"foreignKey:AccessoryID" json:"recommendations,omitempty"`
	Insurance          []accessory.AccessoryInsurance           `gorm:"foreignKey:AccessoryID" json:"insurance,omitempty"`
	AuthenticationLogs []accessory.AccessoryAuthenticationLog   `gorm:"foreignKey:AccessoryID" json:"authentication_logs,omitempty"`
	CustomerFeedback   []accessory.AccessoryCustomerFeedback    `gorm:"foreignKey:AccessoryID" json:"customer_feedback,omitempty"`
	MarketAnalytics    []accessory.AccessoryMarketAnalytics     `gorm:"foreignKey:AccessoryID" json:"market_analytics,omitempty"`
	Subscriptions      []accessory.AccessorySubscription        `gorm:"foreignKey:AccessoryID" json:"subscriptions,omitempty"`
}

func (Accessory) TableName() string {
	return "accessories"
}

// BeforeCreate hook
func (a *Accessory) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.AccessoryIdentification.SKU == "" {
		a.AccessoryIdentification.SKU = a.GenerateSKU()
	}
	if a.AccessoryInventory.MinStock == 0 {
		a.AccessoryInventory.MinStock = accessory.DefaultMinStock
	}
	if a.AccessoryWarranty.ReturnDays == 0 {
		a.AccessoryWarranty.ReturnDays = accessory.DefaultReturnDays
	}
	return a.Validate()
}

// ============== Business Logic Methods for Accessory ==============

// Validation Methods
func (a *Accessory) Validate() error {
	if a.AccessoryIdentification.SKU == "" {
		return errors.New("SKU is required")
	}
	if a.AccessoryBasicInfo.Name == "" {
		return errors.New("accessory name is required")
	}
	if a.AccessoryBasicInfo.Brand == "" {
		return errors.New("brand is required")
	}
	if a.AccessoryBasicInfo.Type == "" {
		return errors.New("accessory type is required")
	}
	if a.AccessoryPricing.RetailPrice.LessThanOrEqual(decimal.Zero) {
		return errors.New(accessory.ErrInvalidPrice)
	}
	if a.AccessoryInventory.CurrentStock < 0 {
		return errors.New(accessory.ErrInvalidStock)
	}
	if a.AccessoryPricing.CostPrice.GreaterThan(a.AccessoryPricing.RetailPrice) {
		return errors.New("cost price cannot exceed retail price")
	}
	if a.AccessoryPricing.MaxDiscount > accessory.DefaultDiscountMax {
		return errors.New(accessory.ErrExceededMaxDiscount)
	}
	return nil
}

// Identification Methods
func (a *Accessory) GenerateSKU() string {
	prefix := strings.ToUpper(a.AccessoryBasicInfo.Type[:3])
	brand := strings.ToUpper(a.AccessoryBasicInfo.Brand[:3])
	timestamp := time.Now().Unix()
	return fmt.Sprintf("ACC-%s-%s-%d", prefix, brand, timestamp%100000)
}

func (a *Accessory) GetDisplayName() string {
	return fmt.Sprintf("%s %s %s", a.AccessoryBasicInfo.Brand, a.AccessoryBasicInfo.Model, a.AccessoryBasicInfo.Name)
}

func (a *Accessory) GetFullSKU() string {
	if a.AccessoryIdentification.PartNumber != "" {
		return fmt.Sprintf("%s-%s", a.AccessoryIdentification.SKU, a.AccessoryIdentification.PartNumber)
	}
	return a.AccessoryIdentification.SKU
}

// Stock Management Methods
func (a *Accessory) IsInStock() bool {
	return a.AccessoryInventory.AvailableStock > 0 &&
		a.AccessoryInventory.StockStatus == accessory.StockInStock
}

func (a *Accessory) IsLowStock() bool {
	return a.AccessoryInventory.CurrentStock <= accessory.LowStockThreshold &&
		a.AccessoryInventory.CurrentStock > accessory.CriticalStockThreshold
}

func (a *Accessory) IsCriticalStock() bool {
	return a.AccessoryInventory.CurrentStock <= accessory.CriticalStockThreshold
}

func (a *Accessory) NeedsReorder() bool {
	return a.AccessoryInventory.CurrentStock <= a.AccessoryInventory.ReorderPoint
}

func (a *Accessory) GetAvailableQuantity() int {
	return a.AccessoryInventory.AvailableStock
}

func (a *Accessory) ReserveStock(quantity int) error {
	if quantity > a.AccessoryInventory.AvailableStock {
		return errors.New(accessory.ErrInsufficientStock)
	}
	a.AccessoryInventory.AvailableStock -= quantity
	a.AccessoryInventory.ReservedStock += quantity
	return nil
}

func (a *Accessory) ReleaseStock(quantity int) {
	a.AccessoryInventory.AvailableStock += quantity
	a.AccessoryInventory.ReservedStock -= quantity
	if a.AccessoryInventory.ReservedStock < 0 {
		a.AccessoryInventory.ReservedStock = 0
	}
}

func (a *Accessory) UpdateStock(quantity int, movementType string) {
	switch movementType {
	case accessory.MovementInbound:
		a.AccessoryInventory.CurrentStock += quantity
		a.AccessoryInventory.AvailableStock += quantity
	case accessory.MovementOutbound, accessory.MovementSale:
		a.AccessoryInventory.CurrentStock -= quantity
		a.AccessoryInventory.AvailableStock -= quantity
	case accessory.MovementDamage, accessory.MovementLoss:
		a.AccessoryInventory.CurrentStock -= quantity
		a.AccessoryInventory.AvailableStock -= quantity
	}
	a.UpdateStockStatus()
}

func (a *Accessory) UpdateStockStatus() {
	switch {
	case a.AccessoryInventory.CurrentStock <= 0:
		a.AccessoryInventory.StockStatus = accessory.StockOutOfStock
	case a.AccessoryInventory.CurrentStock <= accessory.CriticalStockThreshold:
		a.AccessoryInventory.StockStatus = accessory.StockLowStock
	case a.AccessoryStatus.IsDiscontinued:
		a.AccessoryInventory.StockStatus = accessory.StockDiscontinued
	default:
		a.AccessoryInventory.StockStatus = accessory.StockInStock
	}
}

// Pricing Methods
func (a *Accessory) GetCurrentPrice() decimal.Decimal {
	if a.AccessoryPricing.PromoPrice.GreaterThan(decimal.Zero) && a.IsOnPromotion() {
		return a.AccessoryPricing.PromoPrice
	}
	if a.AccessoryPricing.SalePrice.GreaterThan(decimal.Zero) {
		return a.AccessoryPricing.SalePrice
	}
	return a.AccessoryPricing.RetailPrice
}

func (a *Accessory) CalculateDiscount() decimal.Decimal {
	originalPrice := a.AccessoryPricing.RetailPrice
	currentPrice := a.GetCurrentPrice()

	if originalPrice.Equal(currentPrice) {
		return decimal.Zero
	}

	discount := originalPrice.Sub(currentPrice)
	return discount
}

func (a *Accessory) GetDiscountPercentage() float64 {
	if a.AccessoryPricing.RetailPrice.IsZero() {
		return 0
	}
	discount := a.CalculateDiscount()
	percentage := discount.Div(a.AccessoryPricing.RetailPrice).Mul(decimal.NewFromInt(100))
	result, _ := percentage.Float64()
	return result
}

func (a *Accessory) IsOnPromotion() bool {
	if a.AccessoryPricing.PromoStartDate == nil || a.AccessoryPricing.PromoEndDate == nil {
		return false
	}
	now := time.Now()
	return now.After(*a.AccessoryPricing.PromoStartDate) && now.Before(*a.AccessoryPricing.PromoEndDate)
}

func (a *Accessory) CalculateProfit() decimal.Decimal {
	return a.GetCurrentPrice().Sub(a.AccessoryPricing.CostPrice)
}

func (a *Accessory) GetProfitMargin() float64 {
	if a.GetCurrentPrice().IsZero() {
		return 0
	}
	profit := a.CalculateProfit()
	margin := profit.Div(a.GetCurrentPrice()).Mul(decimal.NewFromInt(100))
	result, _ := margin.Float64()
	return result
}

func (a *Accessory) ApplyBulkDiscount(quantity int) decimal.Decimal {
	price := a.GetCurrentPrice()
	if quantity >= accessory.DefaultBulkMinQuantity {
		discount := price.Mul(decimal.NewFromFloat(accessory.DefaultBulkDiscountRate))
		return price.Sub(discount)
	}
	return price
}

// Compatibility Methods
func (a *Accessory) IsCompatibleWith(brand, model string) bool {
	if a.AccessoryCompatibility.UniversalCompatible {
		return true
	}

	// Check brand compatibility
	if a.AccessoryCompatibility.BrandCompatible != "" {
		var brands []string
		json.Unmarshal([]byte(a.AccessoryCompatibility.BrandCompatible), &brands)
		for _, b := range brands {
			if strings.EqualFold(b, brand) {
				// Check model if brand matches
				if a.AccessoryCompatibility.ModelCompatible != "" {
					var models []string
					json.Unmarshal([]byte(a.AccessoryCompatibility.ModelCompatible), &models)
					for _, m := range models {
						if strings.Contains(strings.ToLower(m), strings.ToLower(model)) {
							return true
						}
					}
					return false
				}
				return true
			}
		}
	}

	// Check compatibility matrix
	for _, compat := range a.Compatibility {
		if strings.EqualFold(compat.DeviceBrand, brand) &&
			strings.Contains(strings.ToLower(compat.DeviceModel), strings.ToLower(model)) {
			return compat.CompatibilityLevel != "none"
		}
	}

	return false
}

func (a *Accessory) GetCompatibilityLevel(brand, model string) string {
	for _, compat := range a.Compatibility {
		if strings.EqualFold(compat.DeviceBrand, brand) &&
			strings.Contains(strings.ToLower(compat.DeviceModel), strings.ToLower(model)) {
			return compat.CompatibilityLevel
		}
	}

	if a.AccessoryCompatibility.UniversalCompatible {
		return "universal"
	}

	return "unknown"
}

// Quality Methods
func (a *Accessory) IsHighQuality() bool {
	return a.AccessoryQuality.QualityGrade == accessory.QualityOriginal ||
		a.AccessoryQuality.QualityGrade == accessory.QualityOEM ||
		a.AccessoryQuality.QualityGrade == accessory.QualityPremium
}

func (a *Accessory) RequiresInspection() bool {
	if a.AccessoryQuality.InspectionRequired {
		return true
	}

	if a.AccessoryQuality.LastInspection != nil {
		daysSinceInspection := time.Since(*a.AccessoryQuality.LastInspection).Hours() / 24
		return daysSinceInspection > 180 // 6 months
	}

	return a.AccessoryQuality.DefectRate > accessory.DefectRateWarning
}

func (a *Accessory) HasHighReturnRate() bool {
	return a.AccessoryQuality.ReturnRate > accessory.ReturnRateWarning
}

func (a *Accessory) IsCertified() bool {
	if a.AccessoryQuality.Certifications == "" {
		return false
	}

	var certs []string
	json.Unmarshal([]byte(a.AccessoryQuality.Certifications), &certs)
	return len(certs) > 0
}

// Warranty Methods
func (a *Accessory) HasWarranty() bool {
	return a.AccessoryWarranty.HasWarranty && a.AccessoryWarranty.WarrantyMonths > 0
}

func (a *Accessory) IsUnderWarranty(purchaseDate time.Time) bool {
	if !a.HasWarranty() {
		return false
	}

	warrantyEnd := purchaseDate.AddDate(0, a.AccessoryWarranty.WarrantyMonths, 0)
	return time.Now().Before(warrantyEnd)
}

func (a *Accessory) GetWarrantyExpiry(purchaseDate time.Time) time.Time {
	return purchaseDate.AddDate(0, a.AccessoryWarranty.WarrantyMonths, 0)
}

func (a *Accessory) IsReturnable(purchaseDate time.Time) bool {
	returnWindow := purchaseDate.AddDate(0, 0, a.AccessoryWarranty.ReturnDays)
	return time.Now().Before(returnWindow)
}

func (a *Accessory) CalculateRestockingFee(returnAmount decimal.Decimal) decimal.Decimal {
	if a.AccessoryWarranty.RestockingFee == 0 {
		return decimal.Zero
	}
	fee := returnAmount.Mul(decimal.NewFromFloat(a.AccessoryWarranty.RestockingFee))
	return fee
}

// Status Methods
func (a *Accessory) IsActive() bool {
	return a.AccessoryStatus.IsActive && !a.AccessoryStatus.IsDiscontinued
}

func (a *Accessory) IsAvailable() bool {
	return a.IsActive() && a.IsInStock()
}

func (a *Accessory) IsTrending() bool {
	return a.AccessoryStatus.IsBestseller || a.AccessoryMetrics.TrendingStatus
}

func (a *Accessory) IsNewArrival() bool {
	if a.AccessoryStatus.LaunchDate == nil {
		return false
	}
	daysSinceLaunch := time.Since(*a.AccessoryStatus.LaunchDate).Hours() / 24
	return daysSinceLaunch <= 30 && a.AccessoryStatus.IsNewArrival
}

func (a *Accessory) RequiresShipping() bool {
	return a.AccessoryStatus.RequiresShipping && !a.AccessoryStatus.OnlineOnly
}

func (a *Accessory) CanBeGifted() bool {
	return a.AccessoryStatus.IsGiftable && a.IsAvailable()
}

// Metrics Methods
func (a *Accessory) GetPopularityScore() float64 {
	score := 0.0

	// Rating component (40%)
	if a.AccessoryMetrics.ReviewCount >= accessory.MinReviewCount {
		score += (a.AccessoryMetrics.AverageRating / 5.0) * 40
	}

	// Sales component (30%)
	if a.AccessoryMetrics.SalesCount > 0 {
		salesScore := math.Min(float64(a.AccessoryMetrics.SalesCount)/1000.0, 1.0) * 30
		score += salesScore
	}

	// Conversion rate component (20%)
	score += a.AccessoryMetrics.ConversionRate * 20

	// Repurchase rate component (10%)
	score += a.AccessoryMetrics.RepurchaseRate * 10

	return score
}

func (a *Accessory) IsSlowMoving() bool {
	if a.AccessoryMetrics.LastSoldDate == nil {
		return a.AccessoryMetrics.DaysInInventory > accessory.SlowMovingDays
	}
	daysSinceLastSale := time.Since(*a.AccessoryMetrics.LastSoldDate).Hours() / 24
	return daysSinceLastSale > float64(accessory.SlowMovingDays)
}

func (a *Accessory) IsDeadStock() bool {
	if a.AccessoryMetrics.LastSoldDate == nil {
		return a.AccessoryMetrics.DaysInInventory > accessory.DeadStockDays
	}
	daysSinceLastSale := time.Since(*a.AccessoryMetrics.LastSoldDate).Hours() / 24
	return daysSinceLastSale > float64(accessory.DeadStockDays)
}

func (a *Accessory) CalculateTurnoverRate() float64 {
	if a.AccessoryInventory.CurrentStock == 0 {
		return 0
	}
	if a.AccessoryMetrics.DaysInInventory == 0 {
		return 0
	}

	// Annualized turnover
	return (365.0 / float64(a.AccessoryMetrics.DaysInInventory)) *
		float64(a.AccessoryMetrics.SalesCount) / float64(a.AccessoryInventory.CurrentStock)
}

// Review Methods
func (a *Accessory) GetAverageRating() float64 {
	if len(a.Reviews) == 0 {
		return 0
	}

	total := 0
	verified := 0
	for _, review := range a.Reviews {
		if review.Status == accessory.ReviewApproved {
			total += review.Rating
			if review.IsVerified {
				verified++
			}
		}
	}

	if verified == 0 {
		return 0
	}

	return float64(total) / float64(verified)
}

func (a *Accessory) HasGoodReviews() bool {
	return a.AccessoryMetrics.AverageRating >= accessory.PopularityThreshold &&
		a.AccessoryMetrics.ReviewCount >= accessory.MinReviewCount
}

// Alert Methods
func (a *Accessory) GenerateAlerts() []string {
	var alerts []string

	if a.IsCriticalStock() {
		alerts = append(alerts, "Critical stock level")
	} else if a.IsLowStock() {
		alerts = append(alerts, "Low stock warning")
	}

	if a.HasHighReturnRate() {
		alerts = append(alerts, fmt.Sprintf("High return rate: %.1f%%", a.AccessoryQuality.ReturnRate*100))
	}

	if a.AccessoryQuality.DefectRate > accessory.DefectRateWarning {
		alerts = append(alerts, fmt.Sprintf("High defect rate: %.1f%%", a.AccessoryQuality.DefectRate*100))
	}

	if a.IsSlowMoving() {
		alerts = append(alerts, "Slow moving inventory")
	}

	if a.IsDeadStock() {
		alerts = append(alerts, "Dead stock - consider clearance")
	}

	return alerts
}

// Bundle Methods
func (a *Accessory) CanBeBundled() bool {
	return a.IsActive() && a.AccessoryInventory.CurrentStock > accessory.DefaultBulkMinQuantity
}

func (a *Accessory) GetRecommendedBundles() []uuid.UUID {
	var bundleIDs []uuid.UUID
	for _, rec := range a.Recommendations {
		if rec.RecommendationType == "complement" && rec.IsActive {
			bundleIDs = append(bundleIDs, rec.RecommendedID)
		}
	}
	return bundleIDs
}

// Insurance Methods
func (a *Accessory) IsHighValue() bool {
	return a.AccessoryPricing.RetailPrice.GreaterThan(decimal.NewFromFloat(accessory.HighValueThreshold))
}

func (a *Accessory) RequiresInsurance() bool {
	return a.IsHighValue() && a.AccessoryQuality.QualityGrade == accessory.QualityPremium
}

func (a *Accessory) HasActiveInsurance() bool {
	for _, ins := range a.Insurance {
		if ins.Status == accessory.StatusActive &&
			time.Now().After(ins.StartDate) &&
			time.Now().Before(ins.EndDate) {
			return true
		}
	}
	return false
}

// Summary Methods
func (a *Accessory) GetSummary() string {
	return fmt.Sprintf("%s (%s) - %s, Price: %s, Stock: %d, Rating: %.1f/5",
		a.GetDisplayName(),
		a.AccessoryIdentification.SKU,
		a.AccessoryQuality.QualityGrade,
		a.GetCurrentPrice().String(),
		a.AccessoryInventory.CurrentStock,
		a.AccessoryMetrics.AverageRating,
	)
}

func (a *Accessory) GetStockValue() decimal.Decimal {
	return a.AccessoryPricing.CostPrice.Mul(decimal.NewFromInt(int64(a.AccessoryInventory.CurrentStock)))
}

func (a *Accessory) GetRetailValue() decimal.Decimal {
	return a.AccessoryPricing.RetailPrice.Mul(decimal.NewFromInt(int64(a.AccessoryInventory.CurrentStock)))
}

// ============== Authentication Methods ==============

// VerifyAuthenticity checks if accessory is genuine
func (a *Accessory) VerifyAuthenticity(verificationCode string) bool {
	if a.AccessoryAuthentication.VerificationCode == "" {
		return false
	}
	return a.AccessoryAuthentication.VerificationCode == verificationCode &&
		a.AccessoryAuthentication.IsGenuine &&
		!a.AccessoryAuthentication.IsCounterfeit
}

// IsCounterfeit checks if accessory is counterfeit
func (a *Accessory) IsCounterfeit() bool {
	return a.AccessoryAuthentication.IsCounterfeit ||
		a.AccessoryAuthentication.CounterfeitRiskScore > 70
}

// IsFromAuthorizedSource checks if from authorized dealer/reseller
func (a *Accessory) IsFromAuthorizedSource() bool {
	return a.AccessoryAuthentication.AuthorizedDealer ||
		a.AccessoryAuthentication.AuthorizedReseller
}

// NeedsReauthentication checks if needs verification
func (a *Accessory) NeedsReauthentication() bool {
	if a.AccessoryAuthentication.LastVerification == nil {
		return true
	}
	daysSinceVerification := time.Since(*a.AccessoryAuthentication.LastVerification).Hours() / 24
	return daysSinceVerification > 180 // 6 months
}

// ============== Technical Specification Methods ==============

// GetChargingPower returns charging power in watts
func (a *Accessory) GetChargingPower() float64 {
	if a.AccessoryTechnicalSpecs.PowerOutput > 0 {
		return a.AccessoryTechnicalSpecs.PowerOutput
	}
	// Calculate from voltage and current
	if a.AccessoryTechnicalSpecs.OutputVoltage != "" && a.AccessoryTechnicalSpecs.OutputCurrent > 0 {
		// Parse voltage string (e.g., "5V", "9V", "20V")
		voltageStr := strings.TrimSuffix(a.AccessoryTechnicalSpecs.OutputVoltage, "V")
		if voltage, err := strconv.ParseFloat(voltageStr, 64); err == nil {
			return voltage * a.AccessoryTechnicalSpecs.OutputCurrent
		}
	}
	return 0
}

// IsFastCharging checks if supports fast charging
func (a *Accessory) IsFastCharging() bool {
	power := a.GetChargingPower()
	return power >= 18 || // 18W+ is generally fast charging
		strings.Contains(strings.ToLower(a.AccessoryTechnicalSpecs.ChargingSpeed), "fast") ||
		strings.Contains(strings.ToLower(a.AccessoryTechnicalSpecs.ChargingSpeed), "quick") ||
		strings.Contains(strings.ToLower(a.AccessoryTechnicalSpecs.ChargingSpeed), "rapid")
}

// IsWirelessCharging checks if supports wireless charging
func (a *Accessory) IsWirelessCharging() bool {
	return a.AccessoryBasicInfo.Type == accessory.AccessoryTypeWirelessCharger ||
		a.AccessoryTechnicalSpecs.WirelessStandard != ""
}

// GetBluetoothVersion returns bluetooth version if applicable
func (a *Accessory) GetBluetoothVersion() string {
	return a.AccessoryTechnicalSpecs.BluetoothVersion
}

// IsWaterproof checks if accessory is waterproof
func (a *Accessory) IsWaterproof() bool {
	if a.AccessoryTechnicalSpecs.IPRating == "" {
		return false
	}
	// IPX4 and above are water-resistant/waterproof
	if strings.HasPrefix(a.AccessoryTechnicalSpecs.IPRating, "IP") && len(a.AccessoryTechnicalSpecs.IPRating) >= 4 {
		waterRating := a.AccessoryTechnicalSpecs.IPRating[3]
		return waterRating >= '4'
	}
	return false
}

// IsMilSpecRated checks if has military standard rating
func (a *Accessory) IsMilSpecRated() bool {
	return a.AccessoryTechnicalSpecs.MILSTDRating != ""
}

// ============== Market Intelligence Methods ==============

// GetMarketDemandScore returns market demand score
func (a *Accessory) GetMarketDemandScore() float64 {
	return a.AccessoryMarketData.DemandScore
}

// GetTrendStatus returns market trend status
func (a *Accessory) GetTrendStatus() string {
	if a.AccessoryMarketData.TrendStatus != "" {
		return a.AccessoryMarketData.TrendStatus
	}
	// Analyze from metrics
	if len(a.MarketAnalytics) > 1 {
		recent := a.MarketAnalytics[len(a.MarketAnalytics)-1]
		previous := a.MarketAnalytics[len(a.MarketAnalytics)-2]
		if recent.Revenue.GreaterThan(previous.Revenue) {
			return "growing"
		} else if recent.Revenue.LessThan(previous.Revenue) {
			return "declining"
		}
	}
	return "stable"
}

// GetCompetitivePricing returns competitive market price
func (a *Accessory) GetCompetitivePricing() decimal.Decimal {
	if a.AccessoryMarketData.OptimalPrice.GreaterThan(decimal.Zero) {
		return a.AccessoryMarketData.OptimalPrice
	}
	if a.AccessoryMarketData.CompetitorAvgPrice.GreaterThan(decimal.Zero) {
		return a.AccessoryMarketData.CompetitorAvgPrice
	}
	return a.GetCurrentPrice()
}

// GetPricePosition returns price position in market
func (a *Accessory) GetPricePosition() string {
	currentPrice := a.GetCurrentPrice()
	competitorMin := a.AccessoryMarketData.CompetitorMinPrice
	competitorMax := a.AccessoryMarketData.CompetitorMaxPrice

	if competitorMin.IsZero() || competitorMax.IsZero() {
		return "unknown"
	}

	if currentPrice.LessThan(competitorMin) {
		return "budget"
	} else if currentPrice.GreaterThan(competitorMax) {
		return "premium"
	}
	return "competitive"
}

// GetMarketShare returns market share percentage
func (a *Accessory) GetMarketShare() float64 {
	return a.AccessoryMarketData.MarketShare
}

// PredictSeasonalDemand predicts demand for specific month
func (a *Accessory) PredictSeasonalDemand(month int) float64 {
	if a.AccessoryMarketData.SeasonalDemand == "" {
		// Default seasonal factors
		factors := map[int]float64{
			1: 0.9, 2: 0.95, 3: 1.0, 4: 1.05, 5: 1.1, 6: 1.15,
			7: 1.2, 8: 1.2, 9: 1.1, 10: 1.05, 11: 1.2, 12: 1.3,
		}
		baseDemand := a.AccessoryMetrics.SalesCount
		return float64(baseDemand) * factors[month]
	}
	// Parse seasonal demand JSON
	var seasonal map[string]float64
	if err := json.Unmarshal([]byte(a.AccessoryMarketData.SeasonalDemand), &seasonal); err == nil {
		monthStr := fmt.Sprintf("month_%d", month)
		if demand, ok := seasonal[monthStr]; ok {
			return demand
		}
	}
	return float64(a.AccessoryMetrics.SalesCount)
}

// ============== Environmental Methods ==============

// GetEnvironmentalScore returns environmental impact score
func (a *Accessory) GetEnvironmentalScore() float64 {
	if a.AccessoryEnvironmental.EnvironmentalScore > 0 {
		return a.AccessoryEnvironmental.EnvironmentalScore
	}

	score := 100.0

	// Deduct for carbon footprint
	score -= a.AccessoryEnvironmental.CarbonFootprint * 2

	// Add for recyclability
	score += a.AccessoryEnvironmental.RecyclablePercent * 0.3

	// Add for biodegradability
	score += a.AccessoryEnvironmental.BiodegradablePercent * 0.2

	// Bonus for eco-friendly packaging
	if a.AccessoryEnvironmental.PackagingEcoFriendly {
		score += 10
	}

	// Deduct for conflict minerals
	if a.AccessoryEnvironmental.ConflictMinerals {
		score -= 20
	}

	// Bonus for carbon neutral
	if a.AccessoryEnvironmental.CarbonNeutral {
		score += 15
	}

	return math.Max(0, math.Min(100, score))
}

// IsEcoFriendly checks if accessory is environmentally friendly
func (a *Accessory) IsEcoFriendly() bool {
	return a.GetEnvironmentalScore() >= 70 &&
		a.AccessoryEnvironmental.RecyclablePercent >= 60 &&
		!a.AccessoryEnvironmental.ConflictMinerals
}

// GetRecyclingCategory returns e-waste category
func (a *Accessory) GetRecyclingCategory() string {
	if a.AccessoryEnvironmental.EWasteCategory != "" {
		return a.AccessoryEnvironmental.EWasteCategory
	}

	// Determine based on type
	switch a.AccessoryBasicInfo.Type {
	case accessory.AccessoryTypeCharger, accessory.AccessoryTypeWirelessCharger,
		accessory.AccessoryTypePowerBank, accessory.AccessoryTypeBatteryCase:
		return "battery_electronics"
	case accessory.AccessoryTypeHeadphones, accessory.AccessoryTypeTrueWirelessEarbuds,
		accessory.AccessoryTypeSpeaker:
		return "audio_electronics"
	case accessory.AccessoryTypeCable, accessory.AccessoryTypeAdapter:
		return "cables_adapters"
	default:
		return "general_electronics"
	}
}

// GetCarbonOffset calculates carbon offset cost
func (a *Accessory) GetCarbonOffset() decimal.Decimal {
	// Average carbon credit price per ton (example: $15)
	carbonCreditPrice := decimal.NewFromFloat(15)
	carbonTons := decimal.NewFromFloat(a.AccessoryEnvironmental.CarbonFootprint / 1000)
	return carbonTons.Mul(carbonCreditPrice)
}

// ============== Customer Satisfaction Methods ==============

// GetCustomerSatisfactionScore calculates satisfaction score
func (a *Accessory) GetCustomerSatisfactionScore() float64 {
	if len(a.CustomerFeedback) == 0 {
		return 0
	}

	totalScore := 0.0
	count := 0

	for _, feedback := range a.CustomerFeedback {
		// Weight different aspects
		score := float64(feedback.SatisfactionScore) * 10 // Convert to 0-100
		score += float64(feedback.ValueForMoney) * 4      // Max 20
		score += float64(feedback.QualityRating) * 4      // Max 20

		if feedback.RecommendToOthers {
			score += 10
		}
		if feedback.RepurchaseIntent {
			score += 10
		}

		totalScore += score / 1.5 // Normalize to 100
		count++
	}

	if count > 0 {
		return totalScore / float64(count)
	}
	return 0
}

// GetNetPromoterScore calculates NPS
func (a *Accessory) GetNetPromoterScore() float64 {
	if len(a.CustomerFeedback) == 0 {
		return 0
	}

	promoters := 0
	detractors := 0
	total := 0

	for _, feedback := range a.CustomerFeedback {
		if feedback.SatisfactionScore >= 9 {
			promoters++
		} else if feedback.SatisfactionScore <= 6 {
			detractors++
		}
		total++
	}

	if total > 0 {
		return ((float64(promoters) - float64(detractors)) / float64(total)) * 100
	}
	return 0
}

// GetRepurchaseRate returns percentage likely to repurchase
func (a *Accessory) GetRepurchaseRate() float64 {
	if len(a.CustomerFeedback) == 0 {
		return 0
	}

	repurchaseIntent := 0
	for _, feedback := range a.CustomerFeedback {
		if feedback.RepurchaseIntent {
			repurchaseIntent++
		}
	}

	return (float64(repurchaseIntent) / float64(len(a.CustomerFeedback))) * 100
}

// ============== Subscription Methods ==============

// HasSubscriptionOption checks if available as subscription
func (a *Accessory) HasSubscriptionOption() bool {
	// Check if there are active subscriptions
	for _, sub := range a.Subscriptions {
		if sub.Status == "active" {
			return true
		}
	}

	// Consumable accessories are good for subscriptions
	return a.AccessoryBasicInfo.Type == accessory.AccessoryTypeScreenProtector ||
		a.AccessoryBasicInfo.Type == accessory.AccessoryTypeCleaningKit ||
		a.AccessoryBasicInfo.Type == accessory.AccessoryTypeCable
}

// GetSubscriptionPrice calculates subscription price
func (a *Accessory) GetSubscriptionPrice(plan string) decimal.Decimal {
	basePrice := a.GetCurrentPrice()

	switch plan {
	case "monthly":
		return basePrice.Mul(decimal.NewFromFloat(0.85)) // 15% discount
	case "quarterly":
		return basePrice.Mul(decimal.NewFromFloat(0.80)) // 20% discount
	case "annual":
		return basePrice.Mul(decimal.NewFromFloat(0.70)) // 30% discount
	default:
		return basePrice
	}
}

// GetActiveSubscriptions returns count of active subscriptions
func (a *Accessory) GetActiveSubscriptions() int {
	count := 0
	for _, sub := range a.Subscriptions {
		if sub.Status == "active" {
			count++
		}
	}
	return count
}

// ============== Analytics Enhancement Methods ==============

// GetConversionRate calculates conversion rate
func (a *Accessory) GetConversionRate() float64 {
	if len(a.MarketAnalytics) == 0 {
		return 0
	}

	latest := a.MarketAnalytics[len(a.MarketAnalytics)-1]
	return latest.ConversionRate
}

// GetCrossSellingSuccess returns cross-sell success rate
func (a *Accessory) GetCrossSellingSuccess() float64 {
	if len(a.MarketAnalytics) == 0 {
		return 0
	}

	latest := a.MarketAnalytics[len(a.MarketAnalytics)-1]
	return latest.CrossSellSuccess
}

// GetBestPairingProducts returns best products to pair with
func (a *Accessory) GetBestPairingProducts() []uuid.UUID {
	var pairings []uuid.UUID

	// Get recommendations with high conversion
	for _, rec := range a.Recommendations {
		if rec.RecommendationType == "complement" && rec.ConversionRate > 0.2 {
			pairings = append(pairings, rec.RecommendedID)
		}
	}

	return pairings
}

// IsSmartAccessory checks if accessory has smart features
func (a *Accessory) IsSmartAccessory() bool {
	// Check for smart features
	hasApp := strings.Contains(strings.ToLower(a.AccessoryBasicInfo.Description), "app")
	hasBluetooth := a.AccessoryTechnicalSpecs.BluetoothVersion != ""
	hasWifi := a.AccessoryTechnicalSpecs.WifiStandard != ""

	// Specific smart accessory types
	smartTypes := []string{
		accessory.AccessoryTypeTrueWirelessEarbuds,
		accessory.AccessoryTypeSmartwatchBand,
		accessory.AccessoryTypeGimbalStabilizer,
		accessory.AccessoryTypeVRHeadset,
		accessory.AccessoryTypeGameController,
	}

	for _, smartType := range smartTypes {
		if a.AccessoryBasicInfo.Type == smartType {
			return true
		}
	}

	return hasApp || hasBluetooth || hasWifi
}
