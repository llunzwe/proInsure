package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smartsure/internal/domain/models/device"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SparePartsService manages device spare parts inventory
// This is an APPLICATION SERVICE - handles database and orchestration
type SparePartsService struct {
	db *gorm.DB
}

// NewSparePartsService creates a new spare parts service
func NewSparePartsService(db *gorm.DB) *SparePartsService {
	return &SparePartsService{db: db}
}

// AddSparePart adds a new spare part to inventory
func (s *SparePartsService) AddSparePart(ctx context.Context, deviceID string, req *AddSparePartRequest) (*device.DeviceSparePart, error) {
	// Verify device exists
	var device device.Device
	if err := s.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error; err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	part := &device.DeviceSparePart{
		DeviceID:           uuid.MustParse(deviceID),
		PartType:           req.PartType,
		PartNumber:         req.PartNumber,
		PartName:           req.PartName,
		Manufacturer:       req.Manufacturer,
		SupplierName:       req.SupplierName,
		SupplierPartNumber: req.SupplierPartNumber,
		Quality:            req.Quality,
		IsOriginalPart:     req.IsOriginalPart,
		IsCertified:        req.IsCertified,
		CertificationBody:  req.CertificationBody,
		QualityGrade:       req.QualityGrade,
		QuantityInStock:    req.InitialQuantity,
		MinimumStock:       req.MinimumStock,
		UnitCost:           req.UnitCost,
		RetailPrice:        req.RetailPrice,
		WholesalePrice:     req.WholesalePrice,
		WarrantyDays:       req.WarrantyDays,
		WarehouseLocation:  req.WarehouseLocation,
		ShelfNumber:        req.ShelfNumber,
		CompatibleModels:   req.CompatibleModels,
		Specifications:     req.Specifications,
		InstallDifficulty:  req.InstallDifficulty,
		InsuranceApproved:  req.InsuranceApproved,
		MaxClaimAmount:     req.MaxClaimAmount,
		LeadTimeDays:       req.LeadTimeDays,
	}

	// Set reorder point if not specified
	if part.ReorderPoint == 0 && part.MinimumStock > 0 {
		part.ReorderPoint = part.MinimumStock * 2
	}

	// Save part
	if err := s.db.WithContext(ctx).Create(part).Error; err != nil {
		return nil, fmt.Errorf("failed to add spare part: %w", err)
	}

	// Log inventory addition
	s.logInventoryTransaction(ctx, part.ID.String(), "added", req.InitialQuantity, "Initial stock")

	return part, nil
}

// UpdateStock updates spare part stock quantity
func (s *SparePartsService) UpdateStock(ctx context.Context, partID string, req *UpdateStockRequest) error {
	var part device.DeviceSparePart
	if err := s.db.WithContext(ctx).First(&part, "id = ?", partID).Error; err != nil {
		return fmt.Errorf("spare part not found: %w", err)
	}

	oldQuantity := part.QuantityInStock
	part.UpdateStock(req.Quantity, req.Operation)

	// Update last received date if adding stock
	if req.Operation == "add" {
		now := time.Now()
		part.LastReceivedDate = &now
	}

	if err := s.db.WithContext(ctx).Save(&part).Error; err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Log transaction
	s.logInventoryTransaction(ctx, partID, req.Operation, req.Quantity, req.Reason)

	// Check if reorder needed
	if part.NeedsReorder() {
		s.createReorderAlert(ctx, &part)
	}

	// Alert if out of stock
	if part.QuantityInStock == 0 {
		s.createOutOfStockAlert(ctx, &part)
	}

	return nil
}

// UseSparePart records usage of a spare part in a repair
func (s *SparePartsService) UseSparePart(ctx context.Context, partID string, repairID string, quantity int) error {
	var part device.DeviceSparePart
	if err := s.db.WithContext(ctx).First(&part, "id = ?", partID).Error; err != nil {
		return fmt.Errorf("spare part not found: %w", err)
	}

	// Check stock availability
	if part.QuantityInStock < quantity {
		return fmt.Errorf("insufficient stock: available=%d, requested=%d", part.QuantityInStock, quantity)
	}

	// Check if part is approved for insurance if repair is insurance claim
	repair := s.getRepairInfo(ctx, repairID)
	if repair != nil && repair.InsuranceClaim && !part.IsApprovedForInsurance() {
		return errors.New("part not approved for insurance claims")
	}

	// Record usage
	for i := 0; i < quantity; i++ {
		part.RecordUsage(uuid.MustParse(repairID))
	}

	if err := s.db.WithContext(ctx).Save(&part).Error; err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}

	// Log usage
	s.logInventoryTransaction(ctx, partID, "used", quantity, fmt.Sprintf("Used in repair %s", repairID))

	// Check if reorder needed
	if part.NeedsReorder() {
		s.createReorderAlert(ctx, &part)
	}

	return nil
}

// GetAvailableParts gets all available spare parts for a device model
func (s *SparePartsService) GetAvailableParts(ctx context.Context, deviceID string, partType string) ([]device.DeviceSparePart, error) {
	query := s.db.WithContext(ctx).Where("device_id = ? AND status = ? AND is_active = ?", deviceID, "available", true)

	if partType != "" {
		query = query.Where("part_type = ?", partType)
	}

	var parts []device.DeviceSparePart
	if err := query.Find(&parts).Error; err != nil {
		return nil, fmt.Errorf("failed to get spare parts: %w", err)
	}

	// Filter for in-stock parts only
	inStockParts := []device.DeviceSparePart{}
	for _, part := range parts {
		if part.IsInStock() {
			inStockParts = append(inStockParts, part)
		}
	}

	return inStockParts, nil
}

// GetPartsNeedingReorder gets parts that need reordering
func (s *SparePartsService) GetPartsNeedingReorder(ctx context.Context) ([]device.DeviceSparePart, error) {
	var parts []device.DeviceSparePart

	err := s.db.WithContext(ctx).
		Where("quantity_in_stock <= reorder_point AND status = ? AND is_active = ?", "available", true).
		Find(&parts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get parts needing reorder: %w", err)
	}

	return parts, nil
}

// CreateReorderRequest creates a reorder request for low stock parts
func (s *SparePartsService) CreateReorderRequest(ctx context.Context, partID string) (*ReorderRequest, error) {
	var part device.DeviceSparePart
	if err := s.db.WithContext(ctx).First(&part, "id = ?", partID).Error; err != nil {
		return nil, fmt.Errorf("spare part not found: %w", err)
	}

	reorderQty := part.CalculateReorderQuantity()

	request := &ReorderRequest{
		PartID:            partID,
		PartNumber:        part.PartNumber,
		PartName:          part.PartName,
		SupplierName:      part.SupplierName,
		CurrentStock:      part.QuantityInStock,
		ReorderPoint:      part.ReorderPoint,
		ReorderQuantity:   reorderQty,
		UnitCost:          part.UnitCost,
		TotalCost:         float64(reorderQty) * part.UnitCost,
		LeadTimeDays:      part.LeadTimeDays,
		Priority:          part.GetReplacementPriority(),
		RequestedAt:       time.Now(),
		EstimatedDelivery: time.Now().AddDate(0, 0, part.LeadTimeDays),
	}

	// Save reorder request to database
	if err := s.saveReorderRequest(ctx, request); err != nil {
		return nil, err
	}

	// Update last order date
	now := time.Now()
	part.LastOrderDate = &now
	s.db.WithContext(ctx).Save(&part)

	return request, nil
}

// GetInventoryValue calculates total inventory value
func (s *SparePartsService) GetInventoryValue(ctx context.Context, deviceID string) (*InventoryValue, error) {
	var parts []device.DeviceSparePart

	query := s.db.WithContext(ctx).Where("is_active = ?", true)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.Find(&parts).Error; err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	result := &InventoryValue{
		CalculatedAt: time.Now(),
		DeviceID:     deviceID,
	}

	for _, part := range parts {
		partValue := float64(part.QuantityInStock) * part.UnitCost
		result.TotalValue += partValue
		result.TotalQuantity += part.QuantityInStock
		result.UniquePartsCount++

		if part.QuantityInStock == 0 {
			result.OutOfStockCount++
		} else if part.NeedsReorder() {
			result.LowStockCount++
		}

		// Calculate potential revenue
		result.PotentialRevenue += float64(part.QuantityInStock) * part.RetailPrice

		// Track by category
		if result.ValueByCategory == nil {
			result.ValueByCategory = make(map[string]float64)
		}
		result.ValueByCategory[part.PartType] += partValue
	}

	// Calculate profit margin
	if result.PotentialRevenue > 0 {
		result.ProfitMargin = ((result.PotentialRevenue - result.TotalValue) / result.PotentialRevenue) * 100
	}

	return result, nil
}

// GetPartPerformanceMetrics gets performance metrics for a spare part
func (s *SparePartsService) GetPartPerformanceMetrics(ctx context.Context, partID string) (*PartPerformanceMetrics, error) {
	var part device.DeviceSparePart
	if err := s.db.WithContext(ctx).First(&part, "id = ?", partID).Error; err != nil {
		return nil, fmt.Errorf("spare part not found: %w", err)
	}

	metrics := &PartPerformanceMetrics{
		PartID:               partID,
		PartNumber:           part.PartNumber,
		PartName:             part.PartName,
		TimesUsed:            part.TimesUsed,
		SuccessRate:          part.SuccessRate,
		DefectRate:           part.DefectRate,
		ReturnRate:           part.ReturnRate,
		CustomerSatisfaction: part.CustomerSatisfaction,
		AverageRepairTime:    part.AverageRepairTime,
		ProfitMargin:         part.GetProfitMargin(),
		IsHighQuality:        part.IsHighQuality(),
		InsuranceApproved:    part.IsApprovedForInsurance(),
		CalculatedAt:         time.Now(),
	}

	// Calculate turnover rate (usage per month)
	if part.TimesUsed > 0 && part.CreatedAt.Before(time.Now().AddDate(0, -1, 0)) {
		monthsSinceCreation := time.Since(part.CreatedAt).Hours() / (24 * 30)
		metrics.TurnoverRate = float64(part.TimesUsed) / monthsSinceCreation
	}

	return metrics, nil
}

// RecallSparePart marks a spare part as recalled
func (s *SparePartsService) RecallSparePart(ctx context.Context, partID string, reason string) error {
	var part device.DeviceSparePart
	if err := s.db.WithContext(ctx).First(&part, "id = ?", partID).Error; err != nil {
		return fmt.Errorf("spare part not found: %w", err)
	}

	part.IsRecalled = true
	now := time.Now()
	part.RecallDate = &now
	part.RecallReason = reason
	part.Status = "recalled"
	part.IsActive = false

	if err := s.db.WithContext(ctx).Save(&part).Error; err != nil {
		return fmt.Errorf("failed to recall part: %w", err)
	}

	// Create alert for recall
	s.createRecallAlert(ctx, &part, reason)

	// Log the recall
	s.logInventoryTransaction(ctx, partID, "recalled", part.QuantityInStock, reason)

	return nil
}

// Helper methods

func (s *SparePartsService) getRepairInfo(ctx context.Context, repairID string) *device.DeviceRepair {
	var repair device.DeviceRepair
	if err := s.db.WithContext(ctx).First(&repair, "id = ?", repairID).Error; err != nil {
		return nil
	}
	return &repair
}

func (s *SparePartsService) logInventoryTransaction(ctx context.Context, partID string, transactionType string, quantity int, notes string) {
	// In production, write to inventory_transactions table
	fmt.Printf("Inventory Transaction: PartID=%s, Type=%s, Quantity=%d, Notes=%s, Time=%s\n",
		partID, transactionType, quantity, notes, time.Now().Format(time.RFC3339))
}

func (s *SparePartsService) createReorderAlert(ctx context.Context, part *device.DeviceSparePart) {
	// In production, create alert in notifications system
	fmt.Printf("REORDER ALERT: Part %s (%s) needs reordering. Current stock: %d, Reorder point: %d\n",
		part.PartNumber, part.PartName, part.QuantityInStock, part.ReorderPoint)
}

func (s *SparePartsService) createOutOfStockAlert(ctx context.Context, part *device.DeviceSparePart) {
	// In production, create critical alert
	fmt.Printf("OUT OF STOCK ALERT: Part %s (%s) is out of stock!\n", part.PartNumber, part.PartName)
}

func (s *SparePartsService) createRecallAlert(ctx context.Context, part *device.DeviceSparePart, reason string) {
	// In production, create critical recall alert
	fmt.Printf("RECALL ALERT: Part %s (%s) recalled. Reason: %s\n", part.PartNumber, part.PartName, reason)
}

func (s *SparePartsService) saveReorderRequest(ctx context.Context, request *ReorderRequest) error {
	// In production, save to reorder_requests table
	fmt.Printf("Reorder Request Created: %+v\n", request)
	return nil
}

// Request/Response structures

type AddSparePartRequest struct {
	PartType           string  `json:"part_type"`
	PartNumber         string  `json:"part_number"`
	PartName           string  `json:"part_name"`
	Manufacturer       string  `json:"manufacturer"`
	SupplierName       string  `json:"supplier_name"`
	SupplierPartNumber string  `json:"supplier_part_number"`
	Quality            string  `json:"quality"`
	IsOriginalPart     bool    `json:"is_original_part"`
	IsCertified        bool    `json:"is_certified"`
	CertificationBody  string  `json:"certification_body"`
	QualityGrade       string  `json:"quality_grade"`
	InitialQuantity    int     `json:"initial_quantity"`
	MinimumStock       int     `json:"minimum_stock"`
	UnitCost           float64 `json:"unit_cost"`
	RetailPrice        float64 `json:"retail_price"`
	WholesalePrice     float64 `json:"wholesale_price"`
	WarrantyDays       int     `json:"warranty_days"`
	WarehouseLocation  string  `json:"warehouse_location"`
	ShelfNumber        string  `json:"shelf_number"`
	CompatibleModels   string  `json:"compatible_models"`
	Specifications     string  `json:"specifications"`
	InstallDifficulty  string  `json:"install_difficulty"`
	InsuranceApproved  bool    `json:"insurance_approved"`
	MaxClaimAmount     float64 `json:"max_claim_amount"`
	LeadTimeDays       int     `json:"lead_time_days"`
}

type UpdateStockRequest struct {
	Quantity  int    `json:"quantity"`
	Operation string `json:"operation"` // add, remove, set
	Reason    string `json:"reason"`
}

type ReorderRequest struct {
	PartID            string    `json:"part_id"`
	PartNumber        string    `json:"part_number"`
	PartName          string    `json:"part_name"`
	SupplierName      string    `json:"supplier_name"`
	CurrentStock      int       `json:"current_stock"`
	ReorderPoint      int       `json:"reorder_point"`
	ReorderQuantity   int       `json:"reorder_quantity"`
	UnitCost          float64   `json:"unit_cost"`
	TotalCost         float64   `json:"total_cost"`
	LeadTimeDays      int       `json:"lead_time_days"`
	Priority          string    `json:"priority"`
	RequestedAt       time.Time `json:"requested_at"`
	EstimatedDelivery time.Time `json:"estimated_delivery"`
	Status            string    `json:"status"`
}

type InventoryValue struct {
	DeviceID         string             `json:"device_id,omitempty"`
	TotalValue       float64            `json:"total_value"`
	TotalQuantity    int                `json:"total_quantity"`
	UniquePartsCount int                `json:"unique_parts_count"`
	OutOfStockCount  int                `json:"out_of_stock_count"`
	LowStockCount    int                `json:"low_stock_count"`
	PotentialRevenue float64            `json:"potential_revenue"`
	ProfitMargin     float64            `json:"profit_margin"`
	ValueByCategory  map[string]float64 `json:"value_by_category"`
	CalculatedAt     time.Time          `json:"calculated_at"`
}

type PartPerformanceMetrics struct {
	PartID               string    `json:"part_id"`
	PartNumber           string    `json:"part_number"`
	PartName             string    `json:"part_name"`
	TimesUsed            int       `json:"times_used"`
	TurnoverRate         float64   `json:"turnover_rate"`
	SuccessRate          float64   `json:"success_rate"`
	DefectRate           float64   `json:"defect_rate"`
	ReturnRate           float64   `json:"return_rate"`
	CustomerSatisfaction float64   `json:"customer_satisfaction"`
	AverageRepairTime    float64   `json:"average_repair_time"`
	ProfitMargin         float64   `json:"profit_margin"`
	IsHighQuality        bool      `json:"is_high_quality"`
	InsuranceApproved    bool      `json:"insurance_approved"`
	CalculatedAt         time.Time `json:"calculated_at"`
}
