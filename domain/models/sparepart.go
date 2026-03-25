package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"smartsure/internal/domain/models/sparepart"
)

// SparePart represents a spare part for device repairs
type SparePart struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// Embedded structs for organization
	sparepart.SparePartIdentification        `gorm:"embedded;embeddedPrefix:id_" json:"identification"`
	sparepart.SparePartBasicInfo             `gorm:"embedded;embeddedPrefix:info_" json:"basic_info"`
	sparepart.SparePartCompatibility         `gorm:"embedded;embeddedPrefix:compat_" json:"compatibility"`
	sparepart.SparePartPricing               `gorm:"embedded;embeddedPrefix:price_" json:"pricing"`
	sparepart.SparePartInventory             `gorm:"embedded;embeddedPrefix:inv_" json:"inventory"`
	sparepart.SparePartSupplier              `gorm:"embedded;embeddedPrefix:supp_" json:"supplier"`
	sparepart.SparePartQuality               `gorm:"embedded;embeddedPrefix:qual_" json:"quality"`
	sparepart.SparePartWarranty              `gorm:"embedded;embeddedPrefix:warr_" json:"warranty"`
	sparepart.SparePartLifecycle             `gorm:"embedded;embeddedPrefix:life_" json:"lifecycle"`
	sparepart.SparePartUsage                 `gorm:"embedded;embeddedPrefix:usage_" json:"usage"`
	sparepart.SparePartCompliance            `gorm:"embedded;embeddedPrefix:compl_" json:"compliance"`
	sparepart.SparePartMetrics               `gorm:"embedded;embeddedPrefix:metrics_" json:"metrics"`
	sparepart.SparePartAuthentication        `gorm:"embedded;embeddedPrefix:auth_" json:"authentication"`
	sparepart.SparePartTechnicalRequirements `gorm:"embedded;embeddedPrefix:tech_" json:"technical_requirements"`
	sparepart.SparePartEnvironmentalImpact   `gorm:"embedded;embeddedPrefix:env_" json:"environmental_impact"`
	sparepart.SparePartB2BPricing            `gorm:"embedded;embeddedPrefix:b2b_" json:"b2b_pricing"`

	// Timestamps
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - Feature Models
	Requests            []sparepart.SparePartRequest            `gorm:"foreignKey:SparePartID" json:"requests,omitempty"`
	Movements           []sparepart.SparePartMovement           `gorm:"foreignKey:SparePartID" json:"movements,omitempty"`
	Allocations         []sparepart.SparePartAllocation         `gorm:"foreignKey:SparePartID" json:"allocations,omitempty"`
	PurchaseOrders      []sparepart.SparePartPurchaseOrder      `gorm:"foreignKey:SparePartID" json:"purchase_orders,omitempty"`
	Returns             []sparepart.SparePartReturn             `gorm:"foreignKey:SparePartID" json:"returns,omitempty"`
	QualityChecks       []sparepart.SparePartQualityCheck       `gorm:"foreignKey:SparePartID" json:"quality_checks,omitempty"`
	UsageHistory        []sparepart.SparePartUsageHistory       `gorm:"foreignKey:SparePartID" json:"usage_history,omitempty"`
	Forecasts           []sparepart.SparePartForecast           `gorm:"foreignKey:SparePartID" json:"forecasts,omitempty"`
	Alerts              []sparepart.SparePartAlert              `gorm:"foreignKey:SparePartID" json:"alerts,omitempty"`
	CostAnalyses        []sparepart.SparePartCostAnalysis       `gorm:"foreignKey:SparePartID" json:"cost_analyses,omitempty"`
	AuthenticationLogs  []sparepart.SparePartAuthenticationLog  `gorm:"foreignKey:SparePartID" json:"authentication_logs,omitempty"`
	InstallationRecords []sparepart.SparePartInstallationRecord `gorm:"foreignKey:SparePartID" json:"installation_records,omitempty"`
	MarketAnalytics     []sparepart.SparePartMarketAnalytics    `gorm:"foreignKey:SparePartID" json:"market_analytics,omitempty"`
	TradeIns            []sparepart.SparePartTradeIn            `gorm:"foreignKey:OldPartID" json:"trade_ins,omitempty"`
}

func (SparePart) TableName() string {
	return "spare_parts"
}

// BeforeCreate hook
func (s *SparePart) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.SparePartIdentification.PartNumber == "" {
		s.SparePartIdentification.PartNumber = s.GeneratePartNumber()
	}
	if s.SparePartInventory.MinimumStock == 0 {
		s.SparePartInventory.MinimumStock = sparepart.DefaultMinStock
	}
	if s.SparePartInventory.ReorderPoint == 0 {
		s.SparePartInventory.ReorderPoint = sparepart.DefaultReorderPoint
	}
	return s.Validate()
}

// ============== Business Logic Methods for SparePart ==============

// Validation Methods
func (s *SparePart) Validate() error {
	if s.SparePartIdentification.PartNumber == "" {
		return errors.New("part number is required")
	}
	if s.SparePartBasicInfo.Name == "" {
		return errors.New("part name is required")
	}
	if s.SparePartBasicInfo.Type == "" {
		return errors.New("part type is required")
	}
	if s.SparePartBasicInfo.Brand == "" {
		return errors.New("brand is required")
	}
	if s.SparePartPricing.CostPrice.LessThanOrEqual(decimal.Zero) {
		return errors.New("cost price must be positive")
	}
	if s.SparePartInventory.CurrentStock < 0 {
		return errors.New(sparepart.ErrInvalidQuantity)
	}
	if s.SparePartInventory.CurrentStock > s.SparePartInventory.MaximumStock {
		return errors.New(sparepart.ErrExceedsMaxStock)
	}
	return nil
}

// Identification Methods
func (s *SparePart) GeneratePartNumber() string {
	prefix := "SP"
	typeCode := strings.ToUpper(s.SparePartBasicInfo.Type[:3])
	brandCode := strings.ToUpper(s.SparePartBasicInfo.Brand[:3])
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%s-%s-%d", prefix, typeCode, brandCode, timestamp%1000000)
}

func (s *SparePart) GetDisplayName() string {
	return fmt.Sprintf("%s - %s %s", s.SparePartBasicInfo.Name, s.SparePartBasicInfo.Brand, s.SparePartBasicInfo.Model)
}

func (s *SparePart) GetFullPartNumber() string {
	if s.SparePartIdentification.OEMPartNumber != "" {
		return fmt.Sprintf("%s (OEM: %s)", s.SparePartIdentification.PartNumber, s.SparePartIdentification.OEMPartNumber)
	}
	return s.SparePartIdentification.PartNumber
}

// Stock Management Methods
func (s *SparePart) IsAvailable() bool {
	return s.SparePartInventory.AvailableStock > 0 &&
		s.SparePartInventory.StockStatus == sparepart.StockAvailable
}

func (s *SparePart) IsCriticalStock() bool {
	return s.SparePartInventory.CurrentStock <= sparepart.CriticalStockLevel
}

func (s *SparePart) IsLowStock() bool {
	return s.SparePartInventory.CurrentStock <= sparepart.LowStockLevel &&
		s.SparePartInventory.CurrentStock > sparepart.CriticalStockLevel
}

func (s *SparePart) NeedsReorder() bool {
	return s.SparePartInventory.CurrentStock <= s.SparePartInventory.ReorderPoint
}

func (s *SparePart) GetAvailableQuantity() int {
	return s.SparePartInventory.AvailableStock - s.SparePartInventory.ReservedStock - s.SparePartInventory.AllocatedStock
}

func (s *SparePart) ReserveStock(quantity int) error {
	available := s.GetAvailableQuantity()
	if quantity > available {
		return fmt.Errorf("%s: requested %d, available %d", sparepart.ErrInsufficientStock, quantity, available)
	}
	s.SparePartInventory.ReservedStock += quantity
	return nil
}

func (s *SparePart) UpdateStockStatus() {
	switch {
	case s.SparePartInventory.CurrentStock <= 0:
		s.SparePartInventory.StockStatus = sparepart.StockOutOfStock
	case s.SparePartInventory.CurrentStock <= sparepart.EmergencyStockLevel:
		s.SparePartInventory.StockStatus = sparepart.StockCritical
	case s.SparePartInventory.CurrentStock <= sparepart.LowStockLevel:
		s.SparePartInventory.StockStatus = sparepart.StockLow
	case s.SparePartLifecycle.IsObsolete:
		s.SparePartInventory.StockStatus = sparepart.StockObsolete
	default:
		s.SparePartInventory.StockStatus = sparepart.StockAvailable
	}
}

// Pricing Methods
func (s *SparePart) GetStandardPrice(priority string) decimal.Decimal {
	basePrice := s.SparePartPricing.StandardPrice

	switch priority {
	case sparepart.PriorityEmergency:
		if s.SparePartPricing.EmergencyPrice.GreaterThan(decimal.Zero) {
			return s.SparePartPricing.EmergencyPrice
		}
		return basePrice.Mul(decimal.NewFromFloat(1.5)) // 50% markup for emergency
	case sparepart.PriorityCritical:
		return basePrice.Mul(decimal.NewFromFloat(1.25)) // 25% markup
	case sparepart.PriorityHigh:
		return basePrice.Mul(decimal.NewFromFloat(1.1)) // 10% markup
	default:
		return basePrice
	}
}

func (s *SparePart) CalculateProfit() decimal.Decimal {
	return s.SparePartPricing.StandardPrice.Sub(s.SparePartPricing.CostPrice)
}

func (s *SparePart) GetProfitMargin() float64 {
	if s.SparePartPricing.StandardPrice.IsZero() {
		return 0
	}
	profit := s.CalculateProfit()
	margin := profit.Div(s.SparePartPricing.StandardPrice).Mul(decimal.NewFromInt(100))
	result, _ := margin.Float64()
	return result
}

// Compatibility Methods
func (s *SparePart) IsCompatibleWith(brand, model string) bool {
	if s.SparePartCompatibility.CompatibilityLevel == sparepart.CompatibilityUniversal {
		return true
	}

	// Check brand compatibility
	if s.SparePartCompatibility.DeviceBrands != "" {
		var brands []string
		json.Unmarshal([]byte(s.SparePartCompatibility.DeviceBrands), &brands)
		for _, b := range brands {
			if strings.EqualFold(b, brand) {
				// Check model
				if s.SparePartCompatibility.DeviceModels != "" {
					var models []string
					json.Unmarshal([]byte(s.SparePartCompatibility.DeviceModels), &models)
					for _, m := range models {
						if strings.Contains(strings.ToLower(m), strings.ToLower(model)) {
							return true
						}
					}
				}
				return true
			}
		}
	}

	return false
}

func (s *SparePart) RequiresProgramming() bool {
	return s.SparePartCompatibility.RequiresProgramming || s.SparePartCompatibility.RequiresCalibration
}

// Quality Methods
func (s *SparePart) IsHighQuality() bool {
	return s.SparePartQuality.QualityGrade == sparepart.QualityOEM ||
		s.SparePartQuality.QualityGrade == sparepart.QualityOriginal ||
		s.SparePartQuality.QualityGrade == sparepart.QualityAAA
}

func (s *SparePart) GetQualityScore() float64 {
	if s.SparePartQuality.QualityScore > 0 {
		return s.SparePartQuality.QualityScore
	}

	// Calculate based on defect and failure rates
	score := 100.0
	score -= s.SparePartQuality.DefectRate * 100
	score -= s.SparePartQuality.FailureRate * 50
	score -= s.SparePartQuality.ReturnRate * 25

	return math.Max(score, 0)
}

// Warranty Methods
func (s *SparePart) HasWarranty() bool {
	return s.SparePartWarranty.HasWarranty && s.SparePartWarranty.WarrantyDays > 0
}

func (s *SparePart) IsUnderWarranty(installDate time.Time) bool {
	if !s.HasWarranty() {
		return false
	}

	warrantyEnd := installDate.AddDate(0, 0, s.SparePartWarranty.WarrantyDays)
	return time.Now().Before(warrantyEnd)
}

// Lifecycle Methods
func (s *SparePart) IsExpired() bool {
	if s.SparePartLifecycle.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*s.SparePartLifecycle.ExpiryDate)
}

func (s *SparePart) IsObsolete() bool {
	if s.SparePartLifecycle.IsObsolete {
		return true
	}

	// Check if part hasn't been used for obsolete period
	if s.SparePartUsage.LastUsedDate != nil {
		daysSinceUse := time.Since(*s.SparePartUsage.LastUsedDate).Hours() / 24
		return daysSinceUse > float64(sparepart.ObsoleteDays)
	}

	return false
}

// Usage Methods
func (s *SparePart) IsCriticalPart() bool {
	return s.SparePartBasicInfo.IsCritical ||
		s.SparePartMetrics.CriticalityScore > 80
}

func (s *SparePart) IsHighUsage() bool {
	return s.SparePartUsage.AverageUsagePerMonth > float64(sparepart.HighUsageThreshold)
}

func (s *SparePart) GetUsageTrend() string {
	if s.SparePartUsage.UsageTrend != "" {
		return s.SparePartUsage.UsageTrend
	}

	if s.SparePartUsage.CurrentMonthUsage > s.SparePartUsage.LastMonthUsage {
		return "increasing"
	} else if s.SparePartUsage.CurrentMonthUsage < s.SparePartUsage.LastMonthUsage {
		return "decreasing"
	}
	return "stable"
}

// Supplier Methods
func (s *SparePart) HasReliableSupplier() bool {
	return s.SparePartSupplier.PreferredSupplier &&
		s.SparePartSupplier.ReliabilityScore >= 80 &&
		s.SparePartSupplier.QualityScore >= sparepart.QualityScoreThreshold
}

func (s *SparePart) GetLeadTime() int {
	if s.SparePartSupplier.LeadTimeDays > 0 {
		return s.SparePartSupplier.LeadTimeDays
	}
	return sparepart.DefaultLeadTimeDays
}

func (s *SparePart) ShouldReorderNow() bool {
	if !s.NeedsReorder() {
		return false
	}

	leadTime := s.GetLeadTime()
	daysOfStock := s.CalculateDaysOfStock()

	return daysOfStock <= leadTime
}

func (s *SparePart) CalculateReorderQuantity() int {
	if s.SparePartInventory.ReorderQuantity > 0 {
		return s.SparePartInventory.ReorderQuantity
	}

	// Economic Order Quantity simplified
	monthlyUsage := s.SparePartUsage.AverageUsagePerMonth
	leadTime := float64(s.GetLeadTime()) / 30.0 // Convert to months
	safetyStock := monthlyUsage * 0.5           // 50% safety stock

	reorderQty := (monthlyUsage * (1 + leadTime)) + safetyStock

	// Round to order multiple if specified
	if s.SparePartSupplier.OrderMultiple > 0 {
		multiple := float64(s.SparePartSupplier.OrderMultiple)
		reorderQty = math.Ceil(reorderQty/multiple) * multiple
	}

	// Ensure minimum order quantity
	if s.SparePartSupplier.MinOrderQuantity > 0 && int(reorderQty) < s.SparePartSupplier.MinOrderQuantity {
		return s.SparePartSupplier.MinOrderQuantity
	}

	return int(reorderQty)
}

// Compliance Methods
func (s *SparePart) IsCompliant() bool {
	return s.SparePartCompliance.ComplianceStatus == sparepart.StatusActive &&
		s.SparePartCompliance.RoHSCompliant &&
		(!s.SparePartCompliance.IsRegulated || s.SparePartCompliance.REACHCompliant)
}

func (s *SparePart) RequiresSpecialHandling() bool {
	return s.SparePartCompliance.HazardousClass != "" ||
		s.SparePartCompliance.HandlingPrecautions != "" ||
		s.SparePartCompliance.StorageRequirements != ""
}

// Metrics Methods
func (s *SparePart) CalculateTurnoverRate() float64 {
	if s.SparePartInventory.CurrentStock == 0 {
		return 0
	}
	annualUsage := s.SparePartUsage.AverageUsagePerMonth * 12
	return annualUsage / float64(s.SparePartInventory.CurrentStock)
}

func (s *SparePart) CalculateDaysOfStock() int {
	if s.SparePartUsage.AverageUsagePerMonth == 0 {
		return 999 // Infinite
	}
	dailyUsage := s.SparePartUsage.AverageUsagePerMonth / 30
	if dailyUsage == 0 {
		return 999
	}
	return int(float64(s.SparePartInventory.CurrentStock) / dailyUsage)
}

func (s *SparePart) GetServiceLevel() float64 {
	if s.SparePartMetrics.ServiceLevel > 0 {
		return s.SparePartMetrics.ServiceLevel
	}

	// Calculate based on stockouts
	if s.SparePartMetrics.StockoutDays == 0 {
		return 100.0
	}

	if s.SparePartLifecycle.ReceiptDate != nil {
		totalDays := time.Since(*s.SparePartLifecycle.ReceiptDate).Hours() / 24
		if totalDays > 0 {
			return (1 - float64(s.SparePartMetrics.StockoutDays)/totalDays) * 100
		}
	}

	return 0
}

// Alert Methods
func (s *SparePart) GenerateAlerts() []string {
	var alerts []string

	if s.IsCriticalStock() {
		alerts = append(alerts, fmt.Sprintf("Critical stock level (%d units)", s.SparePartInventory.CurrentStock))
	} else if s.IsLowStock() {
		alerts = append(alerts, "Low stock warning")
	}

	if s.ShouldReorderNow() {
		alerts = append(alerts, fmt.Sprintf("Reorder now - Lead time: %d days", s.GetLeadTime()))
	}

	if s.IsExpired() {
		alerts = append(alerts, "Part has expired")
	}

	if s.SparePartQuality.DefectRate > sparepart.DefectRateWarning {
		alerts = append(alerts, fmt.Sprintf("High defect rate: %.1f%%", s.SparePartQuality.DefectRate*100))
	}

	if s.IsObsolete() {
		alerts = append(alerts, "Part is obsolete - consider replacement")
	}

	if !s.IsCompliant() {
		alerts = append(alerts, "Compliance check required")
	}

	return alerts
}

// Cost Analysis Methods
func (s *SparePart) GetTotalCostOwnership() decimal.Decimal {
	if s.SparePartMetrics.TotalCostOwnership.GreaterThan(decimal.Zero) {
		return s.SparePartMetrics.TotalCostOwnership
	}

	// Calculate TCO
	purchaseCost := s.SparePartPricing.TotalCost
	carryingCost := s.SparePartMetrics.CarryingCost
	orderingCost := s.SparePartMetrics.OrderingCost

	return purchaseCost.Add(carryingCost).Add(orderingCost)
}

func (s *SparePart) GetCostPerUse() decimal.Decimal {
	if s.SparePartLifecycle.TotalIssued == 0 {
		return s.SparePartPricing.StandardPrice
	}

	totalCost := s.GetTotalCostOwnership()
	return totalCost.Div(decimal.NewFromInt(int64(s.SparePartLifecycle.TotalIssued)))
}

// Summary Methods
func (s *SparePart) GetSummary() string {
	return fmt.Sprintf("%s (%s) - %s, Stock: %d, Cost: %s, Quality: %s",
		s.GetDisplayName(),
		s.SparePartIdentification.PartNumber,
		s.SparePartQuality.QualityGrade,
		s.SparePartInventory.CurrentStock,
		s.SparePartPricing.CostPrice.String(),
		s.SparePartQuality.Condition,
	)
}

func (s *SparePart) GetInventoryValue() decimal.Decimal {
	return s.SparePartPricing.CostPrice.Mul(decimal.NewFromInt(int64(s.SparePartInventory.CurrentStock)))
}

func (s *SparePart) GetLocation() string {
	location := s.SparePartInventory.Location
	if s.SparePartInventory.ShelfNumber != "" {
		location += " / " + s.SparePartInventory.ShelfNumber
	}
	if s.SparePartInventory.BinLocation != "" {
		location += " / " + s.SparePartInventory.BinLocation
	}
	return location
}

// ============== Authentication Methods ==============

// VerifyAuthenticity verifies if the part is genuine
func (s *SparePart) VerifyAuthenticity(verificationCode string) bool {
	if s.SparePartAuthentication.VerificationCode == "" {
		return false
	}
	return s.SparePartAuthentication.VerificationCode == verificationCode &&
		s.SparePartAuthentication.IsGenuine &&
		!s.SparePartAuthentication.IsCounterfeit
}

// GenerateAuthenticationToken generates a unique authentication token
func (s *SparePart) GenerateAuthenticationToken() string {
	data := fmt.Sprintf("%s-%s-%s-%d",
		s.SparePartIdentification.PartNumber,
		s.SparePartIdentification.SerialNumber,
		s.SparePartAuthentication.ManufactureDate,
		time.Now().Unix())
	// In production, use proper cryptographic hashing
	return fmt.Sprintf("AUTH-%s-%s", s.ID.String()[:8], data[:16])
}

// IsCounterfeit checks if part is identified as counterfeit
func (s *SparePart) IsCounterfeit() bool {
	return s.SparePartAuthentication.IsCounterfeit ||
		s.SparePartAuthentication.CounterfeitRiskScore > 80
}

// GetAuthenticationHistory returns authentication verification logs
func (s *SparePart) GetAuthenticationHistory() []sparepart.SparePartAuthenticationLog {
	return s.AuthenticationLogs
}

// NeedsReauthentication checks if part needs verification
func (s *SparePart) NeedsReauthentication() bool {
	if s.SparePartAuthentication.LastVerification == nil {
		return true
	}
	daysSinceVerification := time.Since(*s.SparePartAuthentication.LastVerification).Hours() / 24
	return daysSinceVerification > 90 // Re-verify every 3 months
}

// ============== Installation Support Methods ==============

// GetInstallationGuide returns installation instructions
func (s *SparePart) GetInstallationGuide() string {
	if s.SparePartTechnicalRequirements.VideoGuideURL != "" {
		return s.SparePartTechnicalRequirements.VideoGuideURL
	}
	return s.SparePartTechnicalRequirements.ManualURL
}

// CanCustomerInstall checks if customer can install the part
func (s *SparePart) CanCustomerInstall() bool {
	return s.SparePartTechnicalRequirements.CustomerInstallable &&
		!s.SparePartTechnicalRequirements.ProfessionalRequired &&
		s.SparePartTechnicalRequirements.RepairDifficulty <= 2
}

// GetRequiredToolsList returns list of required tools
func (s *SparePart) GetRequiredToolsList() []string {
	if s.SparePartTechnicalRequirements.RequiredTools == "" {
		return []string{}
	}
	var tools []string
	json.Unmarshal([]byte(s.SparePartTechnicalRequirements.RequiredTools), &tools)
	return tools
}

// EstimateInstallationCost calculates estimated installation cost
func (s *SparePart) EstimateInstallationCost() decimal.Decimal {
	if s.CanCustomerInstall() {
		return decimal.Zero
	}

	// Base rate per hour (example: $50)
	hourlyRate := decimal.NewFromInt(50)
	hours := decimal.NewFromFloat(float64(s.SparePartTechnicalRequirements.EstimatedTime) / 60.0)

	baseCost := hourlyRate.Mul(hours)

	// Add complexity multiplier
	if s.SparePartTechnicalRequirements.RepairDifficulty >= 4 {
		baseCost = baseCost.Mul(decimal.NewFromFloat(1.5))
	}

	return baseCost
}

// GetDifficultyLevel returns human-readable difficulty
func (s *SparePart) GetDifficultyLevel() string {
	switch s.SparePartTechnicalRequirements.RepairDifficulty {
	case 1:
		return "very_easy"
	case 2:
		return "easy"
	case 3:
		return "moderate"
	case 4:
		return "difficult"
	case 5:
		return "expert_only"
	default:
		return "unknown"
	}
}

// ============== Market Intelligence Methods ==============

// GetMarketDemand returns current market demand score
func (s *SparePart) GetMarketDemand() float64 {
	if len(s.MarketAnalytics) == 0 {
		return 0
	}
	// Get latest analytics
	latest := s.MarketAnalytics[len(s.MarketAnalytics)-1]
	return latest.DemandScore
}

// PredictSeasonalDemand predicts demand for a specific month
func (s *SparePart) PredictSeasonalDemand(month int) int {
	baseUsage := s.SparePartUsage.AverageUsagePerMonth

	// Simple seasonal adjustment
	seasonalFactors := map[int]float64{
		1:  0.9,  // January
		2:  0.95, // February
		3:  1.0,  // March
		4:  1.05, // April
		5:  1.1,  // May
		6:  1.15, // June
		7:  1.2,  // July
		8:  1.2,  // August
		9:  1.1,  // September
		10: 1.05, // October
		11: 1.15, // November (Black Friday)
		12: 1.25, // December (Holiday)
	}

	factor := seasonalFactors[month]
	return int(baseUsage * factor)
}

// GetCompetitivePricing returns competitive market price
func (s *SparePart) GetCompetitivePricing() decimal.Decimal {
	if len(s.MarketAnalytics) == 0 {
		return s.SparePartPricing.StandardPrice
	}

	latest := s.MarketAnalytics[len(s.MarketAnalytics)-1]
	if latest.RecommendedPrice.GreaterThan(decimal.Zero) {
		return latest.RecommendedPrice
	}

	return latest.CompetitorAvgPrice
}

// ShouldDiscontinue checks if part should be discontinued
func (s *SparePart) ShouldDiscontinue() bool {
	// Check if obsolete
	if s.IsObsolete() {
		return true
	}

	// Check market demand
	if s.GetMarketDemand() < 20 {
		return true
	}

	// Check profitability
	if s.GetProfitMargin() < 5 {
		return true
	}

	// Check quality issues
	if s.SparePartQuality.DefectRate > 0.1 {
		return true
	}

	return false
}

// GetMarketTrend returns the market trend direction
func (s *SparePart) GetMarketTrend() string {
	if len(s.MarketAnalytics) < 2 {
		return "insufficient_data"
	}

	recent := s.MarketAnalytics[len(s.MarketAnalytics)-1]
	return recent.DemandTrend
}

// ============== Environmental Methods ==============

// GetEnvironmentalScore returns environmental impact score
func (s *SparePart) GetEnvironmentalScore() float64 {
	if s.SparePartEnvironmentalImpact.EnvironmentalScore > 0 {
		return s.SparePartEnvironmentalImpact.EnvironmentalScore
	}

	score := 100.0

	// Deduct for carbon footprint
	score -= s.SparePartEnvironmentalImpact.CarbonFootprint * 2

	// Add for recyclability
	score += s.SparePartEnvironmentalImpact.RecyclablePercent * 0.3

	// Deduct for hazardous waste
	if s.SparePartEnvironmentalImpact.HazardousWaste {
		score -= 20
	}

	// Deduct for conflict minerals
	if s.SparePartEnvironmentalImpact.ConflictMinerals {
		score -= 15
	}

	// Bonus for certifications
	if s.SparePartEnvironmentalImpact.CarbonNeutral {
		score += 10
	}

	return math.Max(0, math.Min(100, score))
}

// GetRecyclingInstructions returns recycling instructions
func (s *SparePart) GetRecyclingInstructions() string {
	if s.SparePartEnvironmentalImpact.RecyclingInstructions != "" {
		return s.SparePartEnvironmentalImpact.RecyclingInstructions
	}

	// Default instructions based on part type
	switch s.SparePartBasicInfo.Type {
	case sparepart.PartTypeBattery:
		return "Take to certified battery recycling center. Do not dispose in regular trash."
	case sparepart.PartTypeScreen, sparepart.PartTypeLCD, sparepart.PartTypeOLED:
		return "E-waste recycling required. Contains hazardous materials."
	case sparepart.PartTypeMotherboard, sparepart.PartTypeLogicBoard:
		return "Contains precious metals. Take to e-waste recycling facility."
	default:
		return "Check local regulations for electronic waste disposal."
	}
}

// CalculateCarbonOffset calculates carbon offset cost
func (s *SparePart) CalculateCarbonOffset() decimal.Decimal {
	// Average carbon credit price per ton (example: $15)
	carbonCreditPrice := decimal.NewFromFloat(15)

	// Convert kg to tons
	carbonTons := decimal.NewFromFloat(s.SparePartEnvironmentalImpact.CarbonFootprint / 1000)

	return carbonTons.Mul(carbonCreditPrice)
}

// IsEnvironmentallyFriendly checks if part meets environmental standards
func (s *SparePart) IsEnvironmentallyFriendly() bool {
	return s.GetEnvironmentalScore() >= 70 &&
		!s.SparePartEnvironmentalImpact.HazardousWaste &&
		!s.SparePartEnvironmentalImpact.ConflictMinerals &&
		s.SparePartEnvironmentalImpact.RecyclablePercent >= 50
}

// ============== B2B Methods ==============

// GetB2BPrice returns B2B price based on quantity
func (s *SparePart) GetB2BPrice(quantity int) decimal.Decimal {
	basePrice := s.SparePartPricing.StandardPrice

	// Check B2B pricing tiers
	if quantity >= s.SparePartB2BPricing.MinQuantity && quantity <= s.SparePartB2BPricing.MaxQuantity {
		if s.SparePartB2BPricing.UnitPrice.GreaterThan(decimal.Zero) {
			return s.SparePartB2BPricing.UnitPrice
		}

		// Apply volume discount
		discount := basePrice.Mul(decimal.NewFromFloat(s.SparePartB2BPricing.VolumeDiscount))
		return basePrice.Sub(discount)
	}

	// Default bulk discount
	if quantity >= 100 {
		return basePrice.Mul(decimal.NewFromFloat(0.7)) // 30% discount
	} else if quantity >= 50 {
		return basePrice.Mul(decimal.NewFromFloat(0.8)) // 20% discount
	} else if quantity >= 10 {
		return basePrice.Mul(decimal.NewFromFloat(0.9)) // 10% discount
	}

	return basePrice
}

// IsContractPricing checks if contract pricing applies
func (s *SparePart) IsContractPricing() bool {
	return s.SparePartB2BPricing.ContractPrice &&
		s.SparePartB2BPricing.ValidUntil != nil &&
		time.Now().Before(*s.SparePartB2BPricing.ValidUntil)
}

// ============== Advanced Analytics Methods ==============

// GetFailurePrediction predicts failure rate
func (s *SparePart) GetFailurePrediction() float64 {
	baseRate := s.SparePartQuality.FailureRate

	// Adjust based on age
	if s.SparePartLifecycle.ManufactureDate != nil {
		ageMonths := time.Since(*s.SparePartLifecycle.ManufactureDate).Hours() / 24 / 30
		ageMultiplier := 1 + (ageMonths/24)*0.1 // 10% increase every 2 years
		baseRate *= ageMultiplier
	}

	return math.Min(baseRate, 1.0)
}

// GetLifespanEstimate estimates part lifespan in months
func (s *SparePart) GetLifespanEstimate() int {
	if s.SparePartLifecycle.UsageLife > 0 {
		return s.SparePartLifecycle.UsageLife
	}

	// Estimate based on quality and type
	baseLifespan := 24 // 2 years default

	// Adjust for quality
	switch s.SparePartQuality.QualityGrade {
	case sparepart.QualityOEM, sparepart.QualityOriginal:
		baseLifespan = 36
	case sparepart.QualityAAA:
		baseLifespan = 30
	case sparepart.QualityAftermarket:
		baseLifespan = 24
	case sparepart.QualityRefurbished:
		baseLifespan = 18
	case sparepart.QualityUsed:
		baseLifespan = 12
	}

	// Adjust for failure rate
	if s.SparePartQuality.FailureRate > 0 {
		baseLifespan = int(float64(baseLifespan) * (1 - s.SparePartQuality.FailureRate))
	}

	return baseLifespan
}

// GetQualityTrend analyzes quality trend over time
func (s *SparePart) GetQualityTrend() string {
	if len(s.QualityChecks) < 2 {
		return "insufficient_data"
	}

	// Compare recent quality checks
	recent := s.QualityChecks[len(s.QualityChecks)-1]
	previous := s.QualityChecks[len(s.QualityChecks)-2]

	recentFailRate := float64(recent.FailedQuantity) / float64(recent.TotalQuantity)
	previousFailRate := float64(previous.FailedQuantity) / float64(previous.TotalQuantity)

	if recentFailRate > previousFailRate*1.1 {
		return "declining"
	} else if recentFailRate < previousFailRate*0.9 {
		return "improving"
	}
	return "stable"
}

// GetCustomerSatisfactionScore calculates satisfaction score
func (s *SparePart) GetCustomerSatisfactionScore() float64 {
	score := 100.0

	// Factor in return rate
	score -= s.SparePartQuality.ReturnRate * 100

	// Factor in failure rate
	score -= s.SparePartQuality.FailureRate * 50

	// Factor in success rate from installations
	successfulInstalls := 0
	totalInstalls := len(s.InstallationRecords)

	for _, install := range s.InstallationRecords {
		if install.Success && install.CustomerApproval {
			successfulInstalls++
		}
	}

	if totalInstalls > 0 {
		successRate := float64(successfulInstalls) / float64(totalInstalls)
		score = score * successRate
	}

	return math.Max(0, math.Min(100, score))
}
