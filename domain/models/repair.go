package models

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"smartsure/internal/domain/models/repair"
)

// RepairShop represents authorized repair shops in the network with comprehensive management
type RepairShop struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	// Embedded structs for organization
	repair.ShopIdentification `gorm:"embedded;embeddedPrefix:shop_" json:"identification"`
	repair.ShopLocation       `gorm:"embedded;embeddedPrefix:location_" json:"location"`
	repair.ShopContact        `gorm:"embedded;embeddedPrefix:contact_" json:"contact"`
	repair.ShopOperations     `gorm:"embedded;embeddedPrefix:ops_" json:"operations"`
	repair.ShopCapabilities   `gorm:"embedded;embeddedPrefix:cap_" json:"capabilities"`
	repair.ShopCertification  `gorm:"embedded;embeddedPrefix:cert_" json:"certification"`
	repair.ShopPerformance    `gorm:"embedded;embeddedPrefix:perf_" json:"performance"`
	repair.ShopFinancial      `gorm:"embedded;embeddedPrefix:fin_" json:"financial"`
	repair.ShopCompliance     `gorm:"embedded;embeddedPrefix:comp_" json:"compliance"`
	repair.ShopIntegration    `gorm:"embedded;embeddedPrefix:int_" json:"integration"`
	repair.ShopStatus         `gorm:"embedded;embeddedPrefix:status_" json:"status"`
	repair.ShopMetrics        `gorm:"embedded;embeddedPrefix:metrics_" json:"metrics"`

	// Timestamps
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - Core
	RepairBookings []RepairBooking `gorm:"foreignKey:RepairShopID" json:"repair_bookings,omitempty"`
	Reviews        []RepairReview  `gorm:"foreignKey:RepairShopID" json:"reviews,omitempty"`

	// Relationships - Feature Models
	Technicians        []repair.RepairTechnician         `gorm:"foreignKey:RepairShopID" json:"technicians,omitempty"`
	PartsInventory     []repair.RepairPartsInventory     `gorm:"foreignKey:RepairShopID" json:"parts_inventory,omitempty"`
	SLAs               []repair.RepairSLA                `gorm:"foreignKey:RepairShopID" json:"slas,omitempty"`
	PerformanceMetrics []repair.RepairPerformanceMetrics `gorm:"foreignKey:RepairShopID" json:"performance_metrics,omitempty"`
	TrainingPrograms   []repair.RepairTraining           `gorm:"foreignKey:RepairShopID" json:"training_programs,omitempty"`
}

// RepairBooking represents a repair booking/appointment with full relationships
type RepairBooking struct {
	repair.RepairBooking // Embed the repair package struct

	// Relationships
	Claim                 *Claim                       `gorm:"foreignKey:ClaimID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"claim,omitempty"`
	User                  User                         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device                Device                       `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	RepairShop            RepairShop                   `gorm:"foreignKey:RepairShopID" json:"repair_shop,omitempty"`
	PartsUsed             []RepairPartsUsed            `gorm:"foreignKey:RepairBookingID" json:"parts_used,omitempty"`
	AccessoriesReplaced   []RepairAccessoriesReplaced  `gorm:"foreignKey:RepairBookingID" json:"accessories_replaced,omitempty"`
	TechnicianAssignments []RepairTechnicianAssignment `gorm:"foreignKey:RepairBookingID" json:"technician_assignments,omitempty"`
	StatusUpdates         []RepairStatusUpdate         `gorm:"foreignKey:RepairBookingID" json:"status_updates,omitempty"`
}

// RepairStatusUpdate represents status updates for repair bookings with full relationships
type RepairStatusUpdate struct {
	repair.RepairStatusUpdate // Embed the repair package struct

	// Relationships
	RepairBooking RepairBooking `gorm:"foreignKey:RepairBookingID" json:"repair_booking,omitempty"`
}

// RepairReview represents customer reviews for repair shops with full relationships
type RepairReview struct {
	repair.RepairReview // Embed the repair package struct

	// Relationships
	RepairShop    RepairShop    `gorm:"foreignKey:RepairShopID" json:"repair_shop,omitempty"`
	RepairBooking RepairBooking `gorm:"foreignKey:RepairBookingID" json:"repair_booking,omitempty"`
	User          User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ReplacementDevice represents replacement devices available in inventory with full relationships
type ReplacementDevice struct {
	repair.ReplacementDevice // Embed the repair package struct

	// Relationships
	ReplacementOrders []ReplacementOrder `gorm:"foreignKey:ReplacementDeviceID" json:"replacement_orders,omitempty"`
}

// ReplacementOrder represents device replacement orders with full relationships
type ReplacementOrder struct {
	repair.ReplacementOrder // Embed the repair package struct

	// Relationships
	Claim             Claim                     `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	User              User                      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OriginalDevice    Device                    `gorm:"foreignKey:OriginalDeviceID" json:"original_device,omitempty"`
	ReplacementDevice ReplacementDevice         `gorm:"foreignKey:ReplacementDeviceID" json:"replacement_device,omitempty"`
	StatusUpdates     []ReplacementStatusUpdate `gorm:"foreignKey:ReplacementOrderID" json:"status_updates,omitempty"`
}

// ReplacementStatusUpdate represents status updates for replacement orders with full relationships
type ReplacementStatusUpdate struct {
	repair.ReplacementStatusUpdate // Embed the repair package struct

	// Relationships
	ReplacementOrder ReplacementOrder `gorm:"foreignKey:ReplacementOrderID" json:"replacement_order,omitempty"`
}

// TemporaryDevice represents temporary devices for lending with full relationships
type TemporaryDevice struct {
	repair.TemporaryDevice // Embed the repair package struct

	// Relationships
	CurrentUser *User        `gorm:"foreignKey:CurrentUserID" json:"current_user,omitempty"`
	LoanHistory []DeviceLoan `gorm:"foreignKey:TemporaryDeviceID" json:"loan_history,omitempty"`
}

// DeviceLoan represents temporary device loans with full relationships
type DeviceLoan struct {
	repair.DeviceLoan // Embed the repair package struct

	// Relationships
	TemporaryDevice TemporaryDevice `gorm:"foreignKey:TemporaryDeviceID" json:"temporary_device,omitempty"`
	User            User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Claim           *Claim          `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
}

// TableName methods
func (RepairShop) TableName() string {
	return "repair_shops"
}

func (RepairBooking) TableName() string {
	return "repair_bookings"
}

func (RepairStatusUpdate) TableName() string {
	return "repair_status_updates"
}

func (RepairReview) TableName() string {
	return "repair_reviews"
}

func (ReplacementDevice) TableName() string {
	return "replacement_devices"
}

func (ReplacementOrder) TableName() string {
	return "replacement_orders"
}

func (ReplacementStatusUpdate) TableName() string {
	return "replacement_status_updates"
}

func (TemporaryDevice) TableName() string {
	return "temporary_devices"
}

func (DeviceLoan) TableName() string {
	return "device_loans"
}

// BeforeCreate hooks
func (rs *RepairShop) BeforeCreate(tx *gorm.DB) error {
	if rs.ID == uuid.Nil {
		rs.ID = uuid.New()
	}
	if rs.ShopIdentification.ShopNumber == "" {
		rs.ShopIdentification.ShopNumber = "RS-" + uuid.New().String()[:8]
	}
	if rs.ShopIdentification.RegistrationDate.IsZero() {
		rs.ShopIdentification.RegistrationDate = time.Now()
	}
	if rs.ShopCertification.CertificationDate.IsZero() {
		rs.ShopCertification.CertificationDate = time.Now()
		rs.ShopCertification.CertificationExpiry = time.Now().AddDate(1, 0, 0)
	}
	return rs.Validate()
}

func (rb *RepairBooking) BeforeCreate(tx *gorm.DB) error {
	if rb.ID == uuid.Nil {
		rb.ID = uuid.New()
	}
	if rb.BookingNumber == "" {
		rb.BookingNumber = "RB-" + uuid.New().String()[:8]
	}
	return nil
}

func (rsu *RepairStatusUpdate) BeforeCreate(tx *gorm.DB) error {
	if rsu.ID == uuid.Nil {
		rsu.ID = uuid.New()
	}
	return nil
}

func (rr *RepairReview) BeforeCreate(tx *gorm.DB) error {
	if rr.ID == uuid.Nil {
		rr.ID = uuid.New()
	}
	return nil
}

func (rd *ReplacementDevice) BeforeCreate(tx *gorm.DB) error {
	if rd.ID == uuid.Nil {
		rd.ID = uuid.New()
	}
	return nil
}

func (ro *ReplacementOrder) BeforeCreate(tx *gorm.DB) error {
	if ro.ID == uuid.Nil {
		ro.ID = uuid.New()
	}
	if ro.OrderNumber == "" {
		ro.OrderNumber = "RO-" + uuid.New().String()[:8]
	}
	return nil
}

func (rsu *ReplacementStatusUpdate) BeforeCreate(tx *gorm.DB) error {
	if rsu.ID == uuid.Nil {
		rsu.ID = uuid.New()
	}
	return nil
}

func (td *TemporaryDevice) BeforeCreate(tx *gorm.DB) error {
	if td.ID == uuid.Nil {
		td.ID = uuid.New()
	}
	return nil
}

func (dl *DeviceLoan) BeforeCreate(tx *gorm.DB) error {
	if dl.ID == uuid.Nil {
		dl.ID = uuid.New()
	}
	if dl.LoanNumber == "" {
		dl.LoanNumber = "DL-" + uuid.New().String()[:8]
	}
	return nil
}

// ============== Business Logic Methods for RepairShop ==============

// Validation Methods
func (rs *RepairShop) Validate() error {
	if rs.Name == "" {
		return errors.New("shop name is required")
	}
	if rs.ShopIdentification.BusinessLicense == "" {
		return errors.New("business license is required")
	}
	if rs.ShopLocation.Address == "" || rs.ShopLocation.City == "" {
		return errors.New("shop address is incomplete")
	}
	if rs.ShopContact.PrimaryPhone == "" || rs.ShopContact.Email == "" {
		return errors.New("contact information is incomplete")
	}
	return nil
}

// Status Methods
func (rs *RepairShop) IsActive() bool {
	return rs.ShopStatus.IsActive &&
		rs.ShopStatus.Status == repair.ShopStatusActive &&
		!rs.ShopOperations.TemporaryClosure
}

func (rs *RepairShop) IsOperational() bool {
	return rs.IsActive() &&
		rs.ShopCompliance.ComplianceStatus == "compliant" &&
		rs.ShopCompliance.LicenseStatus == "active" &&
		time.Now().Before(rs.ShopCompliance.LicenseExpiry)
}

func (rs *RepairShop) IsVerified() bool {
	return rs.ShopStatus.VerificationStatus == "verified" &&
		rs.ShopStatus.OnboardingCompleted
}

func (rs *RepairShop) IsPremiumPartner() bool {
	return rs.ShopCertification.CertificationLevel == repair.CertificationPlatinum ||
		rs.ShopCertification.CertificationLevel == repair.CertificationDiamond
}

func (rs *RepairShop) IsHighPerformer() bool {
	return rs.ShopPerformance.Rating >= repair.DefaultMinimumRating &&
		rs.ShopPerformance.OnTimePercentage >= repair.DefaultOnTimeThreshold &&
		rs.ShopPerformance.SuccessRate >= 0.95
}

// Certification Methods
func (rs *RepairShop) GetCertificationScore() float64 {
	return (rs.ShopCertification.QualityScore +
		rs.ShopCertification.ComplianceScore +
		rs.ShopCertification.SafetyScore +
		rs.ShopCertification.EnvironmentalScore) / 4.0
}

func (rs *RepairShop) NeedsRecertification() bool {
	return time.Now().After(rs.ShopCertification.CertificationExpiry) ||
		time.Now().AddDate(0, 1, 0).After(rs.ShopCertification.CertificationExpiry)
}

func (rs *RepairShop) CanUpgradeCertification() bool {
	score := rs.GetCertificationScore()
	switch rs.ShopCertification.CertificationLevel {
	case repair.CertificationBronze:
		return score >= 70 && rs.ShopPerformance.CompletedRepairs >= 100
	case repair.CertificationSilver:
		return score >= 80 && rs.ShopPerformance.CompletedRepairs >= 500
	case repair.CertificationGold:
		return score >= 90 && rs.ShopPerformance.CompletedRepairs >= 1000
	case repair.CertificationPlatinum:
		return score >= 95 && rs.ShopPerformance.CompletedRepairs >= 5000
	default:
		return false
	}
}

// Performance Methods
func (rs *RepairShop) GetEfficiencyScore() float64 {
	if rs.ShopOperations.CapacityPerDay == 0 {
		return 0
	}
	utilization := float64(rs.ShopOperations.CurrentCapacity) / float64(rs.ShopOperations.CapacityPerDay)
	return math.Min(utilization*rs.ShopPerformance.FirstTimeFixRate, 1.0)
}

func (rs *RepairShop) GetQualityRating() string {
	if rs.ShopPerformance.Rating >= 4.5 {
		return "excellent"
	} else if rs.ShopPerformance.Rating >= 4.0 {
		return "good"
	} else if rs.ShopPerformance.Rating >= 3.0 {
		return "average"
	}
	return "poor"
}

func (rs *RepairShop) HasCapacity(date time.Time) bool {
	if !rs.IsOperational() {
		return false
	}
	return rs.ShopOperations.CurrentCapacity < rs.ShopOperations.CapacityPerDay
}

func (rs *RepairShop) GetAvailableSlots(date time.Time) int {
	if !rs.HasCapacity(date) {
		return 0
	}
	return rs.ShopOperations.CapacityPerDay - rs.ShopOperations.CurrentCapacity
}

// Financial Methods
func (rs *RepairShop) CalculateFee(repairType string, isExpress bool) decimal.Decimal {
	baseFee := rs.ShopFinancial.LaborRate

	// Apply express multiplier
	if isExpress {
		baseFee = baseFee.Mul(decimal.NewFromFloat(repair.ExpressFeeMultiplier))
	}

	// Apply weekend/holiday multipliers if applicable
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		baseFee = baseFee.Mul(decimal.NewFromFloat(repair.WeekendFeeMultiplier))
	}

	return baseFee
}

func (rs *RepairShop) GetCommissionAmount(repairCost decimal.Decimal) decimal.Decimal {
	return repairCost.Mul(decimal.NewFromFloat(rs.ShopFinancial.CommissionRate))
}

func (rs *RepairShop) IsPaymentCurrent() bool {
	return rs.ShopFinancial.PaymentStatus == "current" &&
		rs.ShopFinancial.OutstandingBalance.LessThanOrEqual(decimal.Zero)
}

func (rs *RepairShop) HasCreditAvailable(amount decimal.Decimal) bool {
	available := rs.ShopFinancial.CreditLimit.Sub(rs.ShopFinancial.OutstandingBalance)
	return available.GreaterThanOrEqual(amount)
}

// Location Methods
func (rs *RepairShop) CalculateDistance(lat, lng float64) float64 {
	// Haversine formula for distance calculation
	const earthRadius = 6371 // km

	lat1 := rs.ShopLocation.Latitude * math.Pi / 180
	lat2 := lat * math.Pi / 180
	deltaLat := (lat - rs.ShopLocation.Latitude) * math.Pi / 180
	deltaLon := (lng - rs.ShopLocation.Longitude) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (rs *RepairShop) IsWithinServiceArea(lat, lng float64) bool {
	if !rs.ShopLocation.IsMobileService {
		return true // In-store service available regardless of distance
	}
	distance := rs.CalculateDistance(lat, lng)
	return distance <= float64(rs.ShopLocation.ServiceRadius)
}

func (rs *RepairShop) GetServiceType(customerLat, customerLng float64) string {
	if !rs.ShopLocation.IsMobileService {
		return repair.ServiceInStore
	}
	if rs.IsWithinServiceArea(customerLat, customerLng) {
		return repair.ServiceMobile
	}
	return repair.ServiceMailIn
}

// Capability Methods
func (rs *RepairShop) CanRepairDevice(brand, model string) bool {
	if rs.ShopCapabilities.SupportedBrands == "" {
		return true // Accepts all brands
	}
	return strings.Contains(strings.ToLower(rs.ShopCapabilities.SupportedBrands), strings.ToLower(brand))
}

func (rs *RepairShop) CanPerformRepair(repairType string) bool {
	if rs.ShopCapabilities.Services == "" {
		return false
	}
	return strings.Contains(strings.ToLower(rs.ShopCapabilities.Services), strings.ToLower(repairType))
}

func (rs *RepairShop) HasPartInStock(partNumber string) bool {
	for _, part := range rs.PartsInventory {
		if part.PartNumber == partNumber && part.QuantityInStock > part.ReservedQuantity {
			return true
		}
	}
	return false
}

func (rs *RepairShop) SupportsSameDayService() bool {
	return rs.ShopCapabilities.SameDayRepair && rs.ShopOperations.ExpressService
}

// SLA Methods
func (rs *RepairShop) GetSLALevel() string {
	if len(rs.SLAs) == 0 {
		return repair.SLAStandard
	}
	// Return the highest active SLA level
	for _, sla := range rs.SLAs {
		if sla.IsActive {
			return sla.SLALevel
		}
	}
	return repair.SLAStandard
}

func (rs *RepairShop) GetResponseTimeMinutes(priority string) int {
	baseTime := repair.DefaultResponseTimeMinutes

	switch priority {
	case repair.PriorityUrgent, repair.PriorityCritical:
		return baseTime / 2
	case repair.PriorityHigh:
		return int(float64(baseTime) * 0.75)
	case repair.PriorityVIP:
		return baseTime / 3
	default:
		return baseTime
	}
}

// Compliance Methods
func (rs *RepairShop) IsCompliant() bool {
	return rs.ShopCompliance.ComplianceStatus == "compliant" &&
		rs.ShopCompliance.LicenseStatus == "active" &&
		time.Now().Before(rs.ShopCompliance.InsuranceExpiry) &&
		rs.ShopCompliance.RegulatoryViolations == 0
}

func (rs *RepairShop) RequiresLicenseRenewal() bool {
	return time.Now().AddDate(0, 1, 0).After(rs.ShopCompliance.LicenseExpiry)
}

func (rs *RepairShop) RequiresInsuranceRenewal() bool {
	return time.Now().AddDate(0, 1, 0).After(rs.ShopCompliance.InsuranceExpiry)
}

// Integration Methods
func (rs *RepairShop) IsIntegrated() bool {
	return rs.ShopIntegration.IntegrationStatus == "active" &&
		rs.ShopIntegration.APIEnabled
}

func (rs *RepairShop) RequiresSync() bool {
	if rs.ShopIntegration.LastSyncDate == nil {
		return true
	}
	// Check based on sync frequency
	switch rs.ShopIntegration.SyncFrequency {
	case "hourly":
		return time.Since(*rs.ShopIntegration.LastSyncDate).Hours() >= 1
	case "daily":
		return time.Since(*rs.ShopIntegration.LastSyncDate).Hours() >= 24
	case "weekly":
		return time.Since(*rs.ShopIntegration.LastSyncDate).Hours() >= 168
	default:
		return false
	}
}

// Summary Methods
func (rs *RepairShop) GetSummary() string {
	return fmt.Sprintf("%s - %s certified shop in %s, %s. Rating: %.1f/5 (%d reviews). Capacity: %d/%d",
		rs.Name,
		rs.ShopCertification.CertificationLevel,
		rs.ShopLocation.City,
		rs.ShopLocation.State,
		rs.ShopPerformance.Rating,
		rs.ShopPerformance.ReviewCount,
		rs.ShopOperations.CurrentCapacity,
		rs.ShopOperations.CapacityPerDay,
	)
}

func (rs *RepairShop) GetDisplayName() string {
	if rs.ShopIdentification.ChainID != nil {
		return fmt.Sprintf("%s (%s)", rs.Name, rs.ShopIdentification.ShopNumber)
	}
	return rs.Name
}

// ============== Business Logic Methods for RepairShop - Inventory & Technician Management ==============

// GetAvailableTechnicians returns technicians available for assignment
func (rs *RepairShop) GetAvailableTechnicians() []repair.RepairTechnician {
	var available []repair.RepairTechnician
	for _, tech := range rs.Technicians {
		if tech.IsActive {
			available = append(available, tech)
		}
	}
	return available
}

// GetTechnicianBySpecialization returns technicians with specific specialization
func (rs *RepairShop) GetTechnicianBySpecialization(specialization string) []repair.RepairTechnician {
	var specialized []repair.RepairTechnician
	for _, tech := range rs.Technicians {
		if tech.IsActive && strings.Contains(tech.Specializations, specialization) {
			specialized = append(specialized, tech)
		}
	}
	return specialized
}

// GetBestTechnician returns the best available technician based on rating and success rate
func (rs *RepairShop) GetBestTechnician() *repair.RepairTechnician {
	if len(rs.Technicians) == 0 {
		return nil
	}

	var best *repair.RepairTechnician
	var bestScore float64

	for i, tech := range rs.Technicians {
		if tech.IsActive {
			score := (tech.Rating * 0.6) + (tech.SuccessRate * 0.4)
			if score > bestScore {
				bestScore = score
				best = &rs.Technicians[i]
			}
		}
	}

	return best
}

// CheckPartsAvailability checks if parts are available in inventory
func (rs *RepairShop) CheckPartsAvailability(partNumbers []string) map[string]bool {
	availability := make(map[string]bool)

	for _, partNum := range partNumbers {
		available := false
		for _, part := range rs.PartsInventory {
			if part.PartNumber == partNum && part.QuantityInStock > part.ReservedQuantity {
				available = true
				break
			}
		}
		availability[partNum] = available
	}

	return availability
}

// GetPartsInStock returns all parts currently in stock
func (rs *RepairShop) GetPartsInStock() []repair.RepairPartsInventory {
	var inStock []repair.RepairPartsInventory
	for _, part := range rs.PartsInventory {
		if part.QuantityInStock > 0 && part.IsActive {
			inStock = append(inStock, part)
		}
	}
	return inStock
}

// GetLowStockParts returns parts that need reordering
func (rs *RepairShop) GetLowStockParts() []repair.RepairPartsInventory {
	var lowStock []repair.RepairPartsInventory
	for _, part := range rs.PartsInventory {
		if part.QuantityInStock <= part.ReorderPoint && part.IsActive {
			lowStock = append(lowStock, part)
		}
	}
	return lowStock
}

// ReservePart reserves a part for a repair
func (rs *RepairShop) ReservePart(partNumber string, quantity int) error {
	for i, part := range rs.PartsInventory {
		if part.PartNumber == partNumber {
			available := part.QuantityInStock - part.ReservedQuantity
			if available < quantity {
				return errors.New("insufficient stock")
			}
			rs.PartsInventory[i].ReservedQuantity += quantity
			return nil
		}
	}
	return errors.New("part not found in inventory")
}

// ConsumePart marks a part as consumed from inventory
func (rs *RepairShop) ConsumePart(partNumber string, quantity int) error {
	for i, part := range rs.PartsInventory {
		if part.PartNumber == partNumber {
			if part.QuantityInStock < quantity {
				return errors.New("insufficient stock")
			}
			rs.PartsInventory[i].QuantityInStock -= quantity
			if rs.PartsInventory[i].ReservedQuantity >= quantity {
				rs.PartsInventory[i].ReservedQuantity -= quantity
			}
			return nil
		}
	}
	return errors.New("part not found in inventory")
}

// CalculateInventoryValue calculates total value of parts inventory
func (rs *RepairShop) CalculateInventoryValue() decimal.Decimal {
	total := decimal.Zero
	for _, part := range rs.PartsInventory {
		value := part.UnitCost.Mul(decimal.NewFromInt(int64(part.QuantityInStock)))
		total = total.Add(value)
	}
	return total
}

// GetTechnicianUtilization calculates technician utilization rate
func (rs *RepairShop) GetTechnicianUtilization() float64 {
	if len(rs.Technicians) == 0 {
		return 0
	}

	var totalUtilization float64
	activeTechs := 0

	for _, tech := range rs.Technicians {
		if tech.IsActive {
			// Simplified calculation based on completed repairs vs capacity
			utilization := float64(tech.CompletedRepairs) / float64(rs.ShopOperations.CapacityPerDay*30) // Monthly basis
			if utilization > 1 {
				utilization = 1
			}
			totalUtilization += utilization
			activeTechs++
		}
	}

	if activeTechs == 0 {
		return 0
	}

	return (totalUtilization / float64(activeTechs)) * 100
}

// GetTechnicianCount returns number of active technicians
func (rs *RepairShop) GetTechnicianCount() int {
	count := 0
	for _, tech := range rs.Technicians {
		if tech.IsActive {
			count++
		}
	}
	return count
}

// GetTechnicianCapacity returns total repair capacity based on technicians
func (rs *RepairShop) GetTechnicianCapacity() int {
	return rs.GetTechnicianCount() * repair.DefaultTechnicianDailyCapacity
}

// HasSufficientTechnicians checks if shop has enough technicians for current workload
func (rs *RepairShop) HasSufficientTechnicians() bool {
	if rs.ShopOperations.CurrentCapacity == 0 {
		return true
	}
	requiredTechs := (rs.ShopOperations.CurrentCapacity + repair.DefaultTechnicianDailyCapacity - 1) / repair.DefaultTechnicianDailyCapacity
	return rs.GetTechnicianCount() >= requiredTechs
}

// GetPartsReorderList returns list of parts that need to be reordered
func (rs *RepairShop) GetPartsReorderList() []map[string]interface{} {
	var reorderList []map[string]interface{}

	for _, part := range rs.PartsInventory {
		if part.QuantityInStock <= part.ReorderPoint && part.IsActive {
			reorderItem := map[string]interface{}{
				"part_number":      part.PartNumber,
				"part_name":        part.PartName,
				"current_stock":    part.QuantityInStock,
				"reorder_point":    part.ReorderPoint,
				"reorder_quantity": part.ReorderQuantity,
				"supplier":         part.SupplierID,
				"urgency":          rs.calculateReorderUrgency(part),
			}
			reorderList = append(reorderList, reorderItem)
		}
	}

	return reorderList
}

// calculateReorderUrgency determines how urgently a part needs to be reordered
func (rs *RepairShop) calculateReorderUrgency(part repair.RepairPartsInventory) string {
	stockPercentage := float64(part.QuantityInStock) / float64(part.ReorderPoint)
	if stockPercentage == 0 {
		return "critical"
	} else if stockPercentage < 0.25 {
		return "urgent"
	} else if stockPercentage < 0.5 {
		return "high"
	} else if stockPercentage < 0.75 {
		return "medium"
	}
	return "low"
}

// CanHandleRepairType checks if shop has capability for specific repair type
func (rs *RepairShop) CanHandleRepairType(repairType string, deviceValue float64) bool {
	// Check if repair type is in services
	if !strings.Contains(rs.ShopCapabilities.Services, repairType) {
		return false
	}

	// Check device value limits
	maxValue, _ := rs.ShopCapabilities.MaxDeviceValue.Float64()
	minValue, _ := rs.ShopCapabilities.MinDeviceValue.Float64()

	if deviceValue > maxValue || deviceValue < minValue {
		return false
	}

	// Check special repair capabilities
	switch repairType {
	case "water_damage":
		return rs.ShopCapabilities.WaterDamageRepair
	case "board_level":
		return rs.ShopCapabilities.BoardLevelRepair
	case "data_recovery":
		return rs.ShopCapabilities.DataRecovery
	}

	return true
}

// GetInventoryTurnoverRate calculates inventory turnover rate
func (rs *RepairShop) GetInventoryTurnoverRate() float64 {
	if len(rs.PartsInventory) == 0 {
		return 0
	}

	inventoryValue := rs.CalculateInventoryValue()
	if inventoryValue.IsZero() {
		return 0
	}

	// Calculate based on monthly parts cost
	monthlyPartsCost, _ := rs.ShopMetrics.AveragePartsCost.Float64()
	monthlyPartsCost = monthlyPartsCost * float64(rs.ShopMetrics.TotalDevicesRepaired/12) // Approximate monthly repairs

	inventoryValueFloat, _ := inventoryValue.Float64()
	if inventoryValueFloat == 0 {
		return 0
	}

	return (monthlyPartsCost * 12) / inventoryValueFloat // Annual turnover
}

// ============== Business Logic Methods for RepairBooking ==============

// Validation Methods
func (rb *RepairBooking) Validate() error {
	if rb.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if rb.DeviceID == uuid.Nil {
		return errors.New("device ID is required")
	}
	if rb.RepairShopID == uuid.Nil {
		return errors.New("repair shop ID is required")
	}
	if rb.RepairType == "" {
		return errors.New("repair type is required")
	}
	if rb.IssueDescription == "" {
		return errors.New("issue description is required")
	}
	return nil
}

// Status Methods
func (rb *RepairBooking) IsActive() bool {
	return rb.Status != repair.RepairStatusCompleted &&
		rb.Status != repair.RepairStatusCancelled &&
		rb.Status != repair.RepairStatusFailed
}

func (rb *RepairBooking) IsOverdue() bool {
	if !rb.IsActive() {
		return false
	}
	if rb.EstimatedDuration > 0 {
		estimatedComplete := rb.ScheduledDate.Add(time.Hour * time.Duration(rb.EstimatedDuration))
		return time.Now().After(estimatedComplete)
	}
	return time.Now().After(rb.ScheduledDate.AddDate(0, 0, 1))
}

func (rb *RepairBooking) IsUrgent() bool {
	return rb.Priority == repair.PriorityUrgent ||
		rb.Priority == repair.PriorityCritical ||
		rb.Priority == repair.PriorityEmergency
}

func (rb *RepairBooking) RequiresEscalation() bool {
	if rb.Status == repair.RepairStatusEscalated {
		return false
	}
	// Escalate if overdue by more than 24 hours
	if rb.IsOverdue() {
		overdueDuration := time.Since(rb.ScheduledDate.Add(time.Hour * time.Duration(rb.EstimatedDuration)))
		return overdueDuration.Hours() > 24
	}
	return rb.Priority == repair.PriorityEmergency && rb.Status == repair.RepairStatusDelayed
}

// Progress Methods
func (rb *RepairBooking) GetProgressPercentage() int {
	statusOrder := map[string]int{
		repair.RepairStatusScheduled:      10,
		repair.RepairStatusCheckedIn:      20,
		repair.RepairStatusDiagnosing:     30,
		repair.RepairStatusWaitingParts:   40,
		repair.RepairStatusInProgress:     60,
		repair.RepairStatusQualityCheck:   80,
		repair.RepairStatusCompleted:      90,
		repair.RepairStatusReadyForPickup: 95,
		repair.RepairStatusDelivered:      100,
	}
	if progress, ok := statusOrder[rb.Status]; ok {
		return progress
	}
	return 0
}

func (rb *RepairBooking) CanTransitionTo(newStatus string) bool {
	validTransitions := map[string][]string{
		repair.RepairStatusScheduled:      {repair.RepairStatusCheckedIn, repair.RepairStatusCancelled},
		repair.RepairStatusCheckedIn:      {repair.RepairStatusDiagnosing, repair.RepairStatusCancelled},
		repair.RepairStatusDiagnosing:     {repair.RepairStatusWaitingParts, repair.RepairStatusInProgress, repair.RepairStatusFailed},
		repair.RepairStatusWaitingParts:   {repair.RepairStatusInProgress, repair.RepairStatusDelayed, repair.RepairStatusCancelled},
		repair.RepairStatusInProgress:     {repair.RepairStatusQualityCheck, repair.RepairStatusOnHold, repair.RepairStatusFailed},
		repair.RepairStatusQualityCheck:   {repair.RepairStatusCompleted, repair.RepairStatusInProgress},
		repair.RepairStatusCompleted:      {repair.RepairStatusReadyForPickup},
		repair.RepairStatusReadyForPickup: {repair.RepairStatusDelivered},
	}

	allowedStatuses, exists := validTransitions[rb.Status]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}
	return false
}

func (rb *RepairBooking) UpdateStatus(newStatus string) error {
	if !rb.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition from %s to %s", rb.Status, newStatus)
	}
	rb.Status = newStatus

	// Update completion date if completed
	if newStatus == repair.RepairStatusCompleted {
		now := time.Now()
		rb.CompletedDate = &now
		if rb.EstimatedDuration > 0 {
			rb.ActualDuration = int(now.Sub(rb.ScheduledDate).Hours())
		}
	}

	return nil
}

// Duration Methods
func (rb *RepairBooking) CalculateActualDuration() int {
	if rb.CompletedDate == nil {
		return 0
	}
	return int(rb.CompletedDate.Sub(rb.ScheduledDate).Hours())
}

func (rb *RepairBooking) GetElapsedTime() time.Duration {
	if rb.CompletedDate != nil {
		return rb.CompletedDate.Sub(rb.ScheduledDate)
	}
	return time.Since(rb.ScheduledDate)
}

func (rb *RepairBooking) IsOnTime() bool {
	if rb.CompletedDate == nil {
		return !rb.IsOverdue()
	}
	return rb.ActualDuration <= rb.EstimatedDuration
}

// Cost Methods
func (rb *RepairBooking) GetCostVariance() float64 {
	if rb.EstimatedCost == 0 {
		return 0
	}
	return ((rb.ActualCost - rb.EstimatedCost) / rb.EstimatedCost) * 100
}

func (rb *RepairBooking) IsOverBudget() bool {
	return rb.ActualCost > rb.EstimatedCost*1.1 // 10% tolerance
}

func (rb *RepairBooking) RequiresApproval() bool {
	return rb.EstimatedCost > repair.HighValueThreshold ||
		rb.GetCostVariance() > 20 // More than 20% over estimate
}

// Insurance Methods
func (rb *RepairBooking) IsClaimed() bool {
	return rb.ClaimID != nil
}

func (rb *RepairBooking) GetInsuranceCoverageAmount() float64 {
	if !rb.IsClaimed() {
		return 0
	}
	// This would typically check the claim for coverage amount
	return rb.EstimatedCost * 0.8 // Example: 80% coverage
}

func (rb *RepairBooking) GetCustomerLiability() float64 {
	return rb.ActualCost - rb.GetInsuranceCoverageAmount()
}

// Service Methods
func (rb *RepairBooking) IsExpress() bool {
	return rb.ServiceType == repair.ServiceExpress ||
		rb.ServiceType == repair.ServiceWhileYouWait
}

func (rb *RepairBooking) RequiresPickup() bool {
	return rb.ServiceType == repair.ServicePickup ||
		rb.ServiceType == repair.ServiceMobile
}

func (rb *RepairBooking) RequiresShipping() bool {
	return rb.ServiceType == repair.ServiceMailIn
}

// Summary Methods
func (rb *RepairBooking) GetSummary() string {
	return fmt.Sprintf("Booking %s: %s repair for device %s. Status: %s, Progress: %d%%",
		rb.BookingNumber,
		rb.RepairType,
		rb.DeviceID,
		rb.Status,
		rb.GetProgressPercentage(),
	)
}

func (rb *RepairBooking) GetETAString() string {
	if rb.CompletedDate != nil {
		return "Completed"
	}
	if rb.EstimatedDuration > 0 {
		eta := rb.ScheduledDate.Add(time.Hour * time.Duration(rb.EstimatedDuration))
		if time.Now().After(eta) {
			return "Overdue"
		}
		return eta.Format("2006-01-02 15:04")
	}
	return "TBD"
}

// ============== Business Logic Methods for RepairBooking - Parts & Accessories Management ==============

// CalculatePartsCost calculates total cost of all parts used in the repair
func (rb *RepairBooking) CalculatePartsCost() decimal.Decimal {
	total := decimal.Zero
	for _, part := range rb.PartsUsed {
		total = total.Add(part.TotalCost)
	}
	return total
}

// CalculateAccessoriesCost calculates total cost of all accessories replaced
func (rb *RepairBooking) CalculateAccessoriesCost() decimal.Decimal {
	total := decimal.Zero
	for _, accessory := range rb.AccessoriesReplaced {
		total = total.Add(accessory.TotalCost)
	}
	return total
}

// CalculateLaborCost calculates total labor cost from all technician assignments
func (rb *RepairBooking) CalculateLaborCost() decimal.Decimal {
	total := decimal.Zero
	for _, assignment := range rb.TechnicianAssignments {
		total = total.Add(assignment.TotalLaborCost)
	}
	return total
}

// GetTotalCost calculates the complete repair cost
func (rb *RepairBooking) GetTotalCost() decimal.Decimal {
	partsCost := rb.CalculatePartsCost()
	accessoriesCost := rb.CalculateAccessoriesCost()
	laborCost := rb.CalculateLaborCost()
	diagnosticFee := decimal.NewFromFloat(rb.DiagnosticFee)

	total := partsCost.Add(accessoriesCost).Add(laborCost).Add(diagnosticFee)
	return total.Sub(decimal.NewFromFloat(rb.DiscountAmount))
}

// HasSufficientParts checks if all required parts are available
func (rb *RepairBooking) HasSufficientParts() bool {
	if rb.PartsRequired == "" {
		return true
	}
	// This would typically check against inventory
	return len(rb.PartsUsed) > 0
}

// AssignTechnician assigns a technician to the repair
func (rb *RepairBooking) AssignTechnician(technicianID uuid.UUID, assignmentType string) error {
	for _, assignment := range rb.TechnicianAssignments {
		if assignment.TechnicianID == technicianID && assignment.Status != "reassigned" {
			return errors.New("technician already assigned")
		}
	}

	if assignmentType == "primary" {
		// Check if there's already a primary technician
		for _, assignment := range rb.TechnicianAssignments {
			if assignment.AssignmentType == "primary" && assignment.Status != "reassigned" {
				return errors.New("primary technician already assigned")
			}
		}
		rb.TechnicianID = &technicianID
	}

	return nil
}

// GetPrimaryTechnician returns the primary technician assignment
func (rb *RepairBooking) GetPrimaryTechnician() *RepairTechnicianAssignment {
	for _, assignment := range rb.TechnicianAssignments {
		if assignment.AssignmentType == "primary" && assignment.Status != "reassigned" {
			return &assignment
		}
	}
	return nil
}

// GetActiveTechnicians returns all active technician assignments
func (rb *RepairBooking) GetActiveTechnicians() []RepairTechnicianAssignment {
	var active []RepairTechnicianAssignment
	for _, assignment := range rb.TechnicianAssignments {
		if assignment.Status == "assigned" || assignment.Status == "in_progress" {
			active = append(active, assignment)
		}
	}
	return active
}

// MarkPartAsUsed records a part as used in this repair
func (rb *RepairBooking) MarkPartAsUsed(partID uuid.UUID, quantity int, technicianID uuid.UUID) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Check if part already recorded
	for _, part := range rb.PartsUsed {
		if part.SparePartID == partID && part.TechnicianID == technicianID {
			return errors.New("part already recorded for this technician")
		}
	}

	return nil
}

// CanComplete checks if the repair can be marked as completed
func (rb *RepairBooking) CanComplete() bool {
	if rb.Status == "completed" || rb.Status == "cancelled" {
		return false
	}

	// Check if quality check passed
	if !rb.QualityCheckPassed {
		return false
	}

	// Check if all technician assignments are completed
	for _, assignment := range rb.TechnicianAssignments {
		if assignment.Status == "assigned" || assignment.Status == "in_progress" {
			return false
		}
	}

	return true
}

// CompleteRepair marks the repair as completed
func (rb *RepairBooking) CompleteRepair() error {
	if !rb.CanComplete() {
		return errors.New("repair cannot be completed in current state")
	}

	now := time.Now()
	rb.Status = "completed"
	rb.ActualCompletionTime = &now

	// Set warranty dates if applicable
	if rb.IsUnderWarranty && rb.WarrantyEndDate == nil {
		warrantyEnd := now.AddDate(0, 3, 0) // 3 months warranty
		rb.WarrantyEndDate = &warrantyEnd
		rb.WarrantyStartDate = &now
	}

	return nil
}

// GetPartsByType returns all parts used of a specific type
func (rb *RepairBooking) GetPartsByType(partType string) []RepairPartsUsed {
	var parts []RepairPartsUsed
	for _, part := range rb.PartsUsed {
		if part.PartType == partType {
			parts = append(parts, part)
		}
	}
	return parts
}

// GetAccessoriesByType returns all accessories replaced of a specific type
func (rb *RepairBooking) GetAccessoriesByType(accessoryType string) []RepairAccessoriesReplaced {
	var accessories []RepairAccessoriesReplaced
	for _, accessory := range rb.AccessoriesReplaced {
		if accessory.AccessoryType == accessoryType {
			accessories = append(accessories, accessory)
		}
	}
	return accessories
}

// HasWarrantyParts checks if any parts used have warranty
func (rb *RepairBooking) HasWarrantyParts() bool {
	for _, part := range rb.PartsUsed {
		if part.IsWarrantyClaim {
			return true
		}
	}
	return false
}

// GetWarrantyCoverage calculates total warranty coverage amount
func (rb *RepairBooking) GetWarrantyCoverage() decimal.Decimal {
	total := decimal.Zero
	for _, part := range rb.PartsUsed {
		if part.IsWarrantyClaim {
			total = total.Add(part.TotalCost)
		}
	}
	for _, accessory := range rb.AccessoriesReplaced {
		if accessory.IsWarrantyReplacement {
			total = total.Add(accessory.TotalCost)
		}
	}
	return total
}

// GetCustomerCost calculates the cost customer needs to pay
func (rb *RepairBooking) GetCustomerCost() decimal.Decimal {
	totalCost := rb.GetTotalCost()
	warrantyCoverage := rb.GetWarrantyCoverage()
	insuranceCoverage := decimal.NewFromFloat(rb.GetInsuranceCoverageAmount())

	customerCost := totalCost.Sub(warrantyCoverage).Sub(insuranceCoverage)
	if customerCost.IsNegative() {
		return decimal.Zero
	}
	return customerCost
}

// ============== Business Logic Methods for ReplacementOrder ==============

// Validation Methods
func (ro *ReplacementOrder) Validate() error {
	if ro.ClaimID == uuid.Nil {
		return errors.New("claim ID is required")
	}
	if ro.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if ro.OriginalDeviceID == uuid.Nil {
		return errors.New("original device ID is required")
	}
	if ro.ReplacementDeviceID == uuid.Nil {
		return errors.New("replacement device ID is required")
	}
	if ro.ReplacementType == "" {
		return errors.New("replacement type is required")
	}
	if ro.ReplacementReason == "" {
		return errors.New("replacement reason is required")
	}
	return nil
}

// Status Methods
func (ro *ReplacementOrder) IsActive() bool {
	return ro.Status != "delivered" && ro.Status != "cancelled"
}

func (ro *ReplacementOrder) IsInTransit() bool {
	return ro.Status == "shipped"
}

func (ro *ReplacementOrder) IsDelivered() bool {
	return ro.Status == "delivered" && ro.ActualDelivery != nil
}

func (ro *ReplacementOrder) IsOverdue() bool {
	if ro.Status == "delivered" || ro.Status == "cancelled" {
		return false
	}
	return time.Now().After(ro.EstimatedDelivery)
}

// Delivery Methods
func (ro *ReplacementOrder) IsExpressDelivery() bool {
	return ro.IsExpress || ro.DeliveryMethod == repair.DeliveryExpress ||
		ro.DeliveryMethod == repair.DeliverySameDay
}

func (ro *ReplacementOrder) GetDeliveryDays() int {
	switch ro.DeliveryMethod {
	case repair.DeliverySameDay:
		return 0
	case repair.DeliveryOvernight:
		return 1
	case repair.DeliveryTwoDay:
		return 2
	case repair.DeliveryExpress:
		return 1
	case repair.DeliveryStandard:
		return repair.DefaultStandardDays
	default:
		return 5
	}
}

func (ro *ReplacementOrder) UpdateStatus(newStatus string) error {
	validTransitions := map[string][]string{
		"pending":    {"processing", "cancelled"},
		"processing": {"shipped", "cancelled"},
		"shipped":    {"delivered", "cancelled"},
		"delivered":  {}, // Terminal state
		"cancelled":  {}, // Terminal state
	}

	allowedStatuses, exists := validTransitions[ro.Status]
	if !exists {
		return fmt.Errorf("invalid current status: %s", ro.Status)
	}

	canTransition := false
	for _, status := range allowedStatuses {
		if status == newStatus {
			canTransition = true
			break
		}
	}

	if !canTransition {
		return fmt.Errorf("cannot transition from %s to %s", ro.Status, newStatus)
	}

	ro.Status = newStatus

	// Update delivery timestamp if delivered
	if newStatus == "delivered" {
		now := time.Now()
		ro.ActualDelivery = &now
	}

	return nil
}

func (ro *ReplacementOrder) MarkDelivered() error {
	return ro.UpdateStatus("delivered")
}

func (ro *ReplacementOrder) Cancel() error {
	return ro.UpdateStatus("cancelled")
}

// Return Methods
func (ro *ReplacementOrder) RequiresReturn() bool {
	return ro.ReturnRequired && !ro.IsReturned
}

func (ro *ReplacementOrder) IsReturnOverdue() bool {
	if !ro.RequiresReturn() || ro.ReturnDeadline == nil {
		return false
	}
	return time.Now().After(*ro.ReturnDeadline)
}

func (ro *ReplacementOrder) MarkReturned() {
	ro.IsReturned = true
	now := time.Now()
	ro.ReturnedAt = &now
}

func (ro *ReplacementOrder) GetReturnDaysRemaining() int {
	if !ro.RequiresReturn() || ro.ReturnDeadline == nil {
		return 0
	}
	daysRemaining := int(time.Until(*ro.ReturnDeadline).Hours() / 24)
	if daysRemaining < 0 {
		return 0
	}
	return daysRemaining
}

// Cost Methods
func (ro *ReplacementOrder) CalculateTotalCost() float64 {
	total := float64(0)
	if ro.IsExpress {
		total += ro.ExpressFee
	}
	if ro.ReplacementType == repair.ReplacementUpgrade {
		total += ro.UpgradeFee
	}
	return total
}

func (ro *ReplacementOrder) IsUpgrade() bool {
	return ro.ReplacementType == repair.ReplacementUpgrade
}

func (ro *ReplacementOrder) IsDowngrade() bool {
	return ro.ReplacementType == repair.ReplacementDowngrade
}

// Tracking Methods
func (ro *ReplacementOrder) HasTracking() bool {
	return ro.TrackingNumber != ""
}

func (ro *ReplacementOrder) HasReturnTracking() bool {
	return ro.ReturnTrackingNumber != ""
}

// Summary Methods
func (ro *ReplacementOrder) GetSummary() string {
	return fmt.Sprintf("Order %s: %s replacement due to %s. Status: %s, Delivery: %s",
		ro.OrderNumber,
		ro.ReplacementType,
		ro.ReplacementReason,
		ro.Status,
		ro.DeliveryMethod,
	)
}

func (ro *ReplacementOrder) GetDeliveryStatus() string {
	if ro.IsDelivered() {
		return fmt.Sprintf("Delivered on %s", ro.ActualDelivery.Format("2006-01-02"))
	}
	if ro.IsOverdue() {
		return "Overdue"
	}
	if ro.IsInTransit() {
		return fmt.Sprintf("In transit, expected %s", ro.EstimatedDelivery.Format("2006-01-02"))
	}
	return "Processing"
}

// ============== Business Logic Methods for TemporaryDevice ==============

// Validation Methods
func (td *TemporaryDevice) Validate() error {
	if td.Model == "" {
		return errors.New("device model is required")
	}
	if td.Brand == "" {
		return errors.New("device brand is required")
	}
	if td.IMEI == "" {
		return errors.New("IMEI is required")
	}
	if td.Condition == "" {
		return errors.New("device condition is required")
	}
	return nil
}

// Availability Methods
func (td *TemporaryDevice) IsAvailableForLoan() bool {
	return td.IsAvailable && td.CurrentUserID == nil && td.Condition != repair.ConditionBroken
}

func (td *TemporaryDevice) IsOnLoan() bool {
	return td.CurrentUserID != nil && !td.IsAvailable
}

func (td *TemporaryDevice) IsOverdue() bool {
	if !td.IsOnLoan() || td.LoanEndDate == nil {
		return false
	}
	return time.Now().After(*td.LoanEndDate)
}

func (td *TemporaryDevice) GetDaysRemaining() int {
	if !td.IsOnLoan() || td.LoanEndDate == nil {
		return 0
	}
	daysRemaining := int(time.Until(*td.LoanEndDate).Hours() / 24)
	if daysRemaining < 0 {
		return 0
	}
	return daysRemaining
}

// Loan Management Methods
func (td *TemporaryDevice) LoanTo(userID uuid.UUID, days int) error {
	if !td.IsAvailableForLoan() {
		return errors.New("device is not available for loan")
	}
	if days > td.MaxLoanDuration {
		return fmt.Errorf("loan duration %d exceeds maximum %d days", days, td.MaxLoanDuration)
	}

	td.IsAvailable = false
	td.CurrentUserID = &userID
	now := time.Now()
	td.LoanStartDate = &now
	endDate := now.AddDate(0, 0, days)
	td.LoanEndDate = &endDate

	return nil
}

func (td *TemporaryDevice) ExtendLoan(additionalDays int) error {
	if !td.IsOnLoan() {
		return errors.New("device is not currently on loan")
	}
	if td.LoanEndDate == nil {
		return errors.New("loan end date not set")
	}

	currentDuration := int(td.LoanEndDate.Sub(*td.LoanStartDate).Hours() / 24)
	if currentDuration+additionalDays > td.MaxLoanDuration {
		return fmt.Errorf("extended duration would exceed maximum %d days", td.MaxLoanDuration)
	}

	newEndDate := td.LoanEndDate.AddDate(0, 0, additionalDays)
	td.LoanEndDate = &newEndDate

	return nil
}

func (td *TemporaryDevice) Return() error {
	if !td.IsOnLoan() {
		return errors.New("device is not on loan")
	}

	td.IsAvailable = true
	td.CurrentUserID = nil
	td.LoanStartDate = nil
	td.LoanEndDate = nil

	return nil
}

// Cost Methods
func (td *TemporaryDevice) CalculateLoanCost(days int) float64 {
	return float64(days) * td.DailyRate
}

func (td *TemporaryDevice) CalculateLateFee(overdueDays int) float64 {
	return float64(overdueDays) * td.DailyRate * 1.5 // 50% surcharge for late returns
}

func (td *TemporaryDevice) GetTotalDeposit() float64 {
	return td.SecurityDeposit
}

// Condition Methods
func (td *TemporaryDevice) RequiresMaintenance() bool {
	return td.Condition == repair.ConditionFair || td.Condition == repair.ConditionPoor
}

func (td *TemporaryDevice) CanBeLoaned() bool {
	return td.Condition != repair.ConditionBroken && td.Condition != repair.ConditionSalvage
}

func (td *TemporaryDevice) UpdateCondition(newCondition string) error {
	if td.IsOnLoan() {
		return errors.New("cannot update condition while device is on loan")
	}
	td.Condition = newCondition
	if !td.CanBeLoaned() {
		td.IsAvailable = false
	}
	return nil
}

// Summary Methods
func (td *TemporaryDevice) GetSummary() string {
	status := "Available"
	if td.IsOnLoan() {
		status = fmt.Sprintf("On loan to user %s", *td.CurrentUserID)
	}
	return fmt.Sprintf("%s %s (IMEI: %s) - Condition: %s, Status: %s",
		td.Brand,
		td.Model,
		td.IMEI,
		td.Condition,
		status,
	)
}

// ============== Business Logic Methods for DeviceLoan ==============

// Validation Methods
func (dl *DeviceLoan) Validate() error {
	if dl.TemporaryDeviceID == uuid.Nil {
		return errors.New("temporary device ID is required")
	}
	if dl.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if dl.StartDate.After(dl.EndDate) {
		return errors.New("start date must be before end date")
	}
	return nil
}

// Status Methods
func (dl *DeviceLoan) IsActive() bool {
	return dl.Status == repair.LoanStatusActive
}

func (dl *DeviceLoan) IsOverdue() bool {
	return !dl.IsReturned && time.Now().After(dl.EndDate)
}

func (dl *DeviceLoan) IsLost() bool {
	return dl.Status == repair.LoanStatusLost
}

func (dl *DeviceLoan) IsDamaged() bool {
	return dl.Status == repair.LoanStatusDamaged
}

func (dl *DeviceLoan) GetOverdueDays() int {
	if !dl.IsOverdue() {
		return 0
	}
	return int(time.Since(dl.EndDate).Hours() / 24)
}

// Duration Methods
func (dl *DeviceLoan) GetLoanDuration() int {
	return int(dl.EndDate.Sub(dl.StartDate).Hours() / 24)
}

func (dl *DeviceLoan) GetActualDuration() int {
	if dl.ActualReturnDate == nil {
		return int(time.Since(dl.StartDate).Hours() / 24)
	}
	return int(dl.ActualReturnDate.Sub(dl.StartDate).Hours() / 24)
}

func (dl *DeviceLoan) WasReturnedOnTime() bool {
	if dl.ActualReturnDate == nil {
		return false
	}
	return dl.ActualReturnDate.Before(dl.EndDate) || dl.ActualReturnDate.Equal(dl.EndDate)
}

// Cost Methods
func (dl *DeviceLoan) CalculateTotalCost() float64 {
	days := dl.GetLoanDuration()
	baseCost := float64(days) * dl.DailyRate

	// Add late fees if overdue
	if dl.IsOverdue() && dl.ActualReturnDate != nil {
		overdueDays := int(dl.ActualReturnDate.Sub(dl.EndDate).Hours() / 24)
		lateFee := float64(overdueDays) * dl.DailyRate * 1.5
		baseCost += lateFee
	}

	return baseCost + dl.AdditionalCharges
}

func (dl *DeviceLoan) CalculateLateFees() float64 {
	if !dl.IsOverdue() {
		return 0
	}
	overdueDays := dl.GetOverdueDays()
	return float64(overdueDays) * dl.DailyRate * 1.5
}

func (dl *DeviceLoan) GetDepositRefund() float64 {
	if dl.Status == repair.LoanStatusLost {
		return 0
	}
	if dl.Status == repair.LoanStatusDamaged {
		return math.Max(0, dl.SecurityDeposit-dl.AdditionalCharges)
	}
	return dl.SecurityDeposit
}

// Return Methods
func (dl *DeviceLoan) Return(condition string, damageAssessment string) error {
	if dl.IsReturned {
		return errors.New("device already returned")
	}

	dl.IsReturned = true
	dl.ReturnCondition = condition
	dl.DamageAssessment = damageAssessment
	now := time.Now()
	dl.ActualReturnDate = &now

	// Update status based on condition
	switch condition {
	case repair.ConditionBroken, repair.ConditionSalvage:
		dl.Status = repair.LoanStatusDamaged
	case repair.ConditionGood, repair.ConditionExcellent, repair.ConditionLikeNew:
		dl.Status = repair.LoanStatusReturned
	default:
		dl.Status = repair.LoanStatusReturned
	}

	// Calculate additional charges if damaged
	if dl.Status == repair.LoanStatusDamaged {
		dl.AdditionalCharges = dl.SecurityDeposit * 0.5 // Example: 50% of deposit for damage
	}

	return nil
}

func (dl *DeviceLoan) MarkAsLost() error {
	if dl.IsReturned {
		return errors.New("cannot mark returned device as lost")
	}

	dl.Status = repair.LoanStatusLost
	dl.AdditionalCharges = dl.SecurityDeposit // Full deposit forfeited

	return nil
}

func (dl *DeviceLoan) Extend(additionalDays int) error {
	if dl.IsReturned {
		return errors.New("cannot extend returned loan")
	}
	if dl.Status != repair.LoanStatusActive {
		return errors.New("can only extend active loans")
	}

	dl.EndDate = dl.EndDate.AddDate(0, 0, additionalDays)
	dl.Status = repair.LoanStatusExtended
	dl.TotalCost = dl.CalculateTotalCost()

	return nil
}

// Summary Methods
func (dl *DeviceLoan) GetSummary() string {
	return fmt.Sprintf("Loan %s: Device %s to user %s from %s to %s. Status: %s, Cost: $%.2f",
		dl.LoanNumber,
		dl.TemporaryDeviceID,
		dl.UserID,
		dl.StartDate.Format("2006-01-02"),
		dl.EndDate.Format("2006-01-02"),
		dl.Status,
		dl.TotalCost,
	)
}
