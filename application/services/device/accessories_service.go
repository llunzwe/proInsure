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

// AccessoriesService manages device accessories
// This is an APPLICATION SERVICE - handles database and orchestration
type AccessoriesService struct {
	db *gorm.DB
}

// NewAccessoriesService creates a new accessories service
func NewAccessoriesService(db *gorm.DB) *AccessoriesService {
	return &AccessoriesService{db: db}
}

// RegisterAccessory adds a new accessory to a device
func (s *AccessoriesService) RegisterAccessory(ctx context.Context, deviceID string, req *RegisterAccessoryRequest) (*device.DeviceAccessory, error) {
	// Verify device exists
	var device device.Device
	if err := s.db.WithContext(ctx).First(&device, "id = ?", deviceID).Error; err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	accessory := &device.DeviceAccessory{
		DeviceID:             uuid.MustParse(deviceID),
		AccessoryType:        req.AccessoryType,
		Brand:                req.Brand,
		Model:                req.Model,
		SerialNumber:         req.SerialNumber,
		PurchaseDate:         req.PurchaseDate,
		PurchasePrice:        req.PurchasePrice,
		Condition:            req.Condition,
		IsOriginal:           req.IsOriginal,
		Color:                req.Color,
		Material:             req.Material,
		Compatibility:        req.Compatibility,
		IsCoveredByInsurance: req.RequestInsurance,
		PhotoURLs:            req.PhotoURLs,
		ReceiptURL:           req.ReceiptURL,
	}

	// Calculate current value
	accessory.CurrentValue = accessory.CalculateDepreciation()

	// Set warranty if applicable
	if req.WarrantyMonths > 0 && req.PurchaseDate != nil {
		warrantyExpiry := req.PurchaseDate.AddDate(0, req.WarrantyMonths, 0)
		accessory.WarrantyExpiry = &warrantyExpiry
	}

	// Save accessory
	if err := s.db.WithContext(ctx).Create(accessory).Error; err != nil {
		return nil, fmt.Errorf("failed to register accessory: %w", err)
	}

	return accessory, nil
}

// VerifyAccessory verifies an accessory's authenticity and condition
func (s *AccessoriesService) VerifyAccessory(ctx context.Context, accessoryID string) error {
	var accessory device.DeviceAccessory
	if err := s.db.WithContext(ctx).First(&accessory, "id = ?", accessoryID).Error; err != nil {
		return fmt.Errorf("accessory not found: %w", err)
	}

	// Mark as verified
	accessory.IsVerified = true
	now := time.Now()
	accessory.VerificationDate = &now
	accessory.LastInspectionDate = &now

	// Update value based on condition
	accessory.CurrentValue = accessory.CalculateDepreciation() * accessory.GetConditionScore()

	return s.db.WithContext(ctx).Save(&accessory).Error
}

// ReportLostAccessory marks an accessory as lost
func (s *AccessoriesService) ReportLostAccessory(ctx context.Context, accessoryID string, lostDate time.Time) error {
	var accessory device.DeviceAccessory
	if err := s.db.WithContext(ctx).First(&accessory, "id = ?", accessoryID).Error; err != nil {
		return fmt.Errorf("accessory not found: %w", err)
	}

	if accessory.IsLost {
		return errors.New("accessory already reported as lost")
	}

	accessory.MarkAsLost()

	// Create audit log
	s.logAccessoryEvent(ctx, accessoryID, "reported_lost", "Accessory reported as lost")

	return s.db.WithContext(ctx).Save(&accessory).Error
}

// ReportDamagedAccessory marks an accessory as damaged
func (s *AccessoriesService) ReportDamagedAccessory(ctx context.Context, accessoryID string, req *DamageReportRequest) error {
	var accessory device.DeviceAccessory
	if err := s.db.WithContext(ctx).First(&accessory, "id = ?", accessoryID).Error; err != nil {
		return fmt.Errorf("accessory not found: %w", err)
	}

	accessory.MarkAsDamaged(req.Description)
	accessory.Condition = req.NewCondition

	// Recalculate value
	accessory.CurrentValue = accessory.CalculateDepreciation() * accessory.GetConditionScore()

	// Create audit log
	s.logAccessoryEvent(ctx, accessoryID, "reported_damaged", req.Description)

	return s.db.WithContext(ctx).Save(&accessory).Error
}

// ProcessAccessoryReplacement handles accessory replacement
func (s *AccessoriesService) ProcessAccessoryReplacement(ctx context.Context, oldAccessoryID string) (*device.DeviceAccessory, error) {
	// Get old accessory
	var oldAccessory device.DeviceAccessory
	if err := s.db.WithContext(ctx).First(&oldAccessory, "id = ?", oldAccessoryID).Error; err != nil {
		return nil, fmt.Errorf("accessory not found: %w", err)
	}

	// Check eligibility
	if !oldAccessory.IsEligibleForClaim() {
		return nil, errors.New("accessory not eligible for replacement")
	}

	// Create replacement accessory
	replacementAccessory := &device.DeviceAccessory{
		DeviceID:             oldAccessory.DeviceID,
		AccessoryType:        oldAccessory.AccessoryType,
		Brand:                oldAccessory.Brand,
		Model:                oldAccessory.Model,
		Condition:            "excellent",
		IsOriginal:           oldAccessory.IsOriginal,
		IsCoveredByInsurance: oldAccessory.IsCoveredByInsurance,
		IsIncludedInPolicy:   oldAccessory.IsIncludedInPolicy,
		PurchasePrice:        oldAccessory.CalculateReplacementCost(),
		CurrentValue:         oldAccessory.CalculateReplacementCost(),
		Color:                oldAccessory.Color,
		Material:             oldAccessory.Material,
		Compatibility:        oldAccessory.Compatibility,
		IsVerified:           true,
	}

	// Set purchase date to today
	now := time.Now()
	replacementAccessory.PurchaseDate = &now
	replacementAccessory.VerificationDate = &now

	// Save replacement
	if err := s.db.WithContext(ctx).Create(replacementAccessory).Error; err != nil {
		return nil, fmt.Errorf("failed to create replacement: %w", err)
	}

	// Mark old accessory as replaced
	oldAccessory.MarkAsReplaced(&replacementAccessory.ID)
	s.db.WithContext(ctx).Save(&oldAccessory)

	// Create audit log
	s.logAccessoryEvent(ctx, oldAccessoryID, "replaced", fmt.Sprintf("Replaced with %s", replacementAccessory.ID))

	return replacementAccessory, nil
}

// GetDeviceAccessories retrieves all accessories for a device
func (s *AccessoriesService) GetDeviceAccessories(ctx context.Context, deviceID string, includeReplaced bool) ([]device.DeviceAccessory, error) {
	query := s.db.WithContext(ctx).Where("device_id = ?", deviceID)

	if !includeReplaced {
		query = query.Where("is_replaced = ?", false)
	}

	var accessories []device.DeviceAccessory
	if err := query.Find(&accessories).Error; err != nil {
		return nil, fmt.Errorf("failed to get accessories: %w", err)
	}

	return accessories, nil
}

// CalculateAccessoriesInsuranceValue calculates total insurance value
func (s *AccessoriesService) CalculateAccessoriesInsuranceValue(ctx context.Context, deviceID string) (*AccessoriesInsuranceValue, error) {
	accessories, err := s.GetDeviceAccessories(ctx, deviceID, false)
	if err != nil {
		return nil, err
	}

	result := &AccessoriesInsuranceValue{
		DeviceID:     deviceID,
		CalculatedAt: time.Now(),
	}

	for _, acc := range accessories {
		if acc.IsCoveredByInsurance && !acc.IsLost && !acc.IsReplaced {
			result.TotalValue += acc.CurrentValue
			result.ReplacementCost += acc.CalculateReplacementCost()
			result.CoveredCount++

			result.Items = append(result.Items, AccessoryInsuranceItem{
				AccessoryID:     acc.ID.String(),
				Type:            acc.AccessoryType,
				Brand:           acc.Brand,
				Model:           acc.Model,
				CurrentValue:    acc.CurrentValue,
				ReplacementCost: acc.CalculateReplacementCost(),
				IsVerified:      acc.IsVerified,
				Condition:       acc.Condition,
			})
		}
	}

	// Calculate monthly premium (2% of total value)
	result.MonthlyPremium = result.TotalValue * 0.02

	return result, nil
}

// UpdateAccessoryCondition updates accessory condition after inspection
func (s *AccessoriesService) UpdateAccessoryCondition(ctx context.Context, accessoryID string, condition string, notes string) error {
	var accessory device.DeviceAccessory
	if err := s.db.WithContext(ctx).First(&accessory, "id = ?", accessoryID).Error; err != nil {
		return fmt.Errorf("accessory not found: %w", err)
	}

	accessory.Condition = condition
	now := time.Now()
	accessory.LastInspectionDate = &now

	// Recalculate value
	accessory.CurrentValue = accessory.CalculateDepreciation() * accessory.GetConditionScore()

	// Log the inspection
	s.logAccessoryEvent(ctx, accessoryID, "inspected", notes)

	return s.db.WithContext(ctx).Save(&accessory).Error
}

// CheckAccessoryWarranty checks warranty status
func (s *AccessoriesService) CheckAccessoryWarranty(ctx context.Context, accessoryID string) (*WarrantyStatus, error) {
	var accessory device.DeviceAccessory
	if err := s.db.WithContext(ctx).First(&accessory, "id = ?", accessoryID).Error; err != nil {
		return nil, fmt.Errorf("accessory not found: %w", err)
	}

	status := &WarrantyStatus{
		AccessoryID:     accessoryID,
		IsUnderWarranty: accessory.IsUnderWarranty(),
		CheckedAt:       time.Now(),
	}

	if accessory.WarrantyExpiry != nil {
		status.WarrantyExpiry = *accessory.WarrantyExpiry
		status.DaysRemaining = int(time.Until(*accessory.WarrantyExpiry).Hours() / 24)

		if status.DaysRemaining < 0 {
			status.DaysRemaining = 0
		}
	}

	return status, nil
}

// GetAccessoriesRequiringInspection finds accessories needing inspection
func (s *AccessoriesService) GetAccessoriesRequiringInspection(ctx context.Context) ([]device.DeviceAccessory, error) {
	var accessories []device.DeviceAccessory

	// Find accessories that are insured and need inspection
	err := s.db.WithContext(ctx).
		Where("is_covered_by_insurance = ? AND is_replaced = ? AND is_lost = ?", true, false, false).
		Where("last_inspection_date IS NULL OR last_inspection_date < ?", time.Now().AddDate(0, 0, -365)).
		Find(&accessories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get accessories requiring inspection: %w", err)
	}

	return accessories, nil
}

// Helper method to log accessory events
func (s *AccessoriesService) logAccessoryEvent(ctx context.Context, accessoryID string, eventType string, description string) {
	// In production, this would write to an audit log table
	fmt.Printf("Accessory Event: ID=%s, Type=%s, Description=%s, Time=%s\n",
		accessoryID, eventType, description, time.Now().Format(time.RFC3339))
}

// Request/Response structures

type RegisterAccessoryRequest struct {
	AccessoryType    string     `json:"accessory_type"`
	Brand            string     `json:"brand"`
	Model            string     `json:"model"`
	SerialNumber     string     `json:"serial_number"`
	PurchaseDate     *time.Time `json:"purchase_date"`
	PurchasePrice    float64    `json:"purchase_price"`
	Condition        string     `json:"condition"`
	IsOriginal       bool       `json:"is_original"`
	Color            string     `json:"color"`
	Material         string     `json:"material"`
	Compatibility    string     `json:"compatibility"`
	RequestInsurance bool       `json:"request_insurance"`
	WarrantyMonths   int        `json:"warranty_months"`
	PhotoURLs        string     `json:"photo_urls"`
	ReceiptURL       string     `json:"receipt_url"`
}

type DamageReportRequest struct {
	Description  string `json:"description"`
	NewCondition string `json:"new_condition"`
	PhotoURLs    string `json:"photo_urls"`
}

type AccessoriesInsuranceValue struct {
	DeviceID        string                   `json:"device_id"`
	TotalValue      float64                  `json:"total_value"`
	ReplacementCost float64                  `json:"replacement_cost"`
	CoveredCount    int                      `json:"covered_count"`
	MonthlyPremium  float64                  `json:"monthly_premium"`
	Items           []AccessoryInsuranceItem `json:"items"`
	CalculatedAt    time.Time                `json:"calculated_at"`
}

type AccessoryInsuranceItem struct {
	AccessoryID     string  `json:"accessory_id"`
	Type            string  `json:"type"`
	Brand           string  `json:"brand"`
	Model           string  `json:"model"`
	CurrentValue    float64 `json:"current_value"`
	ReplacementCost float64 `json:"replacement_cost"`
	IsVerified      bool    `json:"is_verified"`
	Condition       string  `json:"condition"`
}

type WarrantyStatus struct {
	AccessoryID     string    `json:"accessory_id"`
	IsUnderWarranty bool      `json:"is_under_warranty"`
	WarrantyExpiry  time.Time `json:"warranty_expiry,omitempty"`
	DaysRemaining   int       `json:"days_remaining"`
	CheckedAt       time.Time `json:"checked_at"`
}
