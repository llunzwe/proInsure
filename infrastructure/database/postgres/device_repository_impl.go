package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
)

// DeviceRepositoryImpl implements the DeviceRepository interface using PostgreSQL
type DeviceRepositoryImpl struct {
	*BaseRepository
	db *gorm.DB
}

// NewDeviceRepository creates a new instance of DeviceRepositoryImpl
func NewDeviceRepository(db *gorm.DB, logger Logger) repositories.DeviceRepository {
	return &DeviceRepositoryImpl{
		BaseRepository: NewBaseRepository(db, logger),
		db:             db,
	}
}

// === Basic CRUD Operations ===

// Create creates a new device
func (r *DeviceRepositoryImpl) Create(ctx context.Context, device *models.Device) error {
	return r.BaseRepository.Create(ctx, device, nil)
}

// GetByID retrieves a device by ID
func (r *DeviceRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := r.FindByID(ctx, id, &device, &QueryOptions{
		Preload: []string{"Owner", "Policies", "Claims", "Repairs"},
	})
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// GetByIMEI retrieves a device by IMEI
func (r *DeviceRepositoryImpl) GetByIMEI(ctx context.Context, imei string) (*models.Device, error) {
	var device models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Where("imei = ?", imei).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &device, nil
}

// GetBySerialNumber retrieves a device by serial number
func (r *DeviceRepositoryImpl) GetBySerialNumber(ctx context.Context, serial string) (*models.Device, error) {
	var device models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Where("serial_number = ?", serial).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &device, nil
}

// Update updates a device
func (r *DeviceRepositoryImpl) Update(ctx context.Context, device *models.Device) error {
	return r.BaseRepository.Update(ctx, device, nil)
}

// Delete deletes a device
func (r *DeviceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.Delete(ctx, id, &models.Device{}, nil)
}

// SoftDelete soft deletes a device
func (r *DeviceRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.Delete(ctx, id, &models.Device{}, nil)
}

// Restore restores a soft-deleted device
func (r *DeviceRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	db := r.GetDB(ctx, nil)
	result := db.Model(&models.Device{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// === Batch Operations ===

// CreateBatch creates multiple devices
func (r *DeviceRepositoryImpl) CreateBatch(ctx context.Context, devices []*models.Device) error {
	return r.BaseRepository.BulkCreate(ctx, devices, nil)
}

// UpdateBatch updates multiple devices
func (r *DeviceRepositoryImpl) UpdateBatch(ctx context.Context, devices []*models.Device) error {
	tx, err := r.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			r.RollbackTransaction(tx)
		}
	}()

	for _, device := range devices {
		if err := tx.Save(device).Error; err != nil {
			return err
		}
	}

	return r.CommitTransaction(tx)
}

// DeleteBatch deletes multiple devices
func (r *DeviceRepositoryImpl) DeleteBatch(ctx context.Context, ids []uuid.UUID) error {
	db := r.GetDB(ctx, nil)
	return db.Where("id IN ?", ids).Delete(&models.Device{}).Error
}

// === Search & Filtering ===

// Search searches devices based on filters
func (r *DeviceRepositoryImpl) Search(ctx context.Context, filters repositories.DeviceSearchFilters) ([]*models.Device, int64, error) {
	var devices []*models.Device
	var totalCount int64

	db := r.GetDB(ctx, nil)
	query := r.buildSearchQuery(db.Model(&models.Device{}), filters)

	// Count total before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Apply sorting
	if filters.SortBy != "" {
		order := "ASC"
		if filters.SortOrder != "" {
			order = filters.SortOrder
		}
		query = query.Order(fmt.Sprintf("%s %s", filters.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// Execute query
	if err := query.Find(&devices).Error; err != nil {
		return nil, 0, err
	}

	return devices, totalCount, nil
}

// buildSearchQuery builds the search query based on filters
func (r *DeviceRepositoryImpl) buildSearchQuery(query *gorm.DB, filters repositories.DeviceSearchFilters) *gorm.DB {
	// Basic filters
	if filters.Brand != "" {
		query = query.Where("brand ILIKE ?", "%"+filters.Brand+"%")
	}
	if filters.Model != "" {
		query = query.Where("model ILIKE ?", "%"+filters.Model+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.OwnerID != nil {
		query = query.Where("owner_id = ?", *filters.OwnerID)
	}

	// Value filters
	if filters.MinValue.GreaterThan(decimal.Zero) {
		query = query.Where("current_value >= ?", filters.MinValue)
	}
	if filters.MaxValue.GreaterThan(decimal.Zero) {
		query = query.Where("current_value <= ?", filters.MaxValue)
	}

	// Risk filters
	if filters.MinRiskScore > 0 {
		query = query.Where("risk_score >= ?", filters.MinRiskScore)
	}
	if filters.MaxRiskScore > 0 {
		query = query.Where("risk_score <= ?", filters.MaxRiskScore)
	}

	// Boolean filters
	if filters.IsInsured != nil {
		query = query.Where("is_insured = ?", *filters.IsInsured)
	}
	if filters.IsStolen != nil {
		query = query.Where("is_stolen = ?", *filters.IsStolen)
	}
	if filters.IsBlacklisted != nil {
		query = query.Where("blacklist_status = ?", *filters.IsBlacklisted)
	}
	if filters.IsCorporateOwned != nil {
		query = query.Where("is_corporate_owned = ?", *filters.IsCorporateOwned)
	}

	// Date filters
	if !filters.PurchasedAfter.IsZero() {
		query = query.Where("purchase_date >= ?", filters.PurchasedAfter)
	}
	if !filters.PurchasedBefore.IsZero() {
		query = query.Where("purchase_date <= ?", filters.PurchasedBefore)
	}
	if !filters.WarrantyExpiresAfter.IsZero() {
		query = query.Where("warranty_end_date >= ?", filters.WarrantyExpiresAfter)
	}
	if !filters.WarrantyExpiresBefore.IsZero() {
		query = query.Where("warranty_end_date <= ?", filters.WarrantyExpiresBefore)
	}

	return query
}

// GetByOwnerID gets devices by owner ID
func (r *DeviceRepositoryImpl) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Where("owner_id = ?", ownerID).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// GetByBrandAndModel gets devices by brand and model
func (r *DeviceRepositoryImpl) GetByBrandAndModel(ctx context.Context, brand, model string) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Where("brand = ? AND model = ?", brand, model).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// GetByStatus gets devices by status
func (r *DeviceRepositoryImpl) GetByStatus(ctx context.Context, status string, limit int) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	query := db.Where("status = ?", status)
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// === Device Status Management ===

// UpdateStatus updates device status
func (r *DeviceRepositoryImpl) UpdateStatus(ctx context.Context, deviceID uuid.UUID, status string) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"status":            status,
		"status_updated_at": time.Now(),
	}, &models.Device{}, nil)
}

// MarkAsStolen marks device as stolen
func (r *DeviceRepositoryImpl) MarkAsStolen(ctx context.Context, deviceID uuid.UUID, stolenDate time.Time) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"is_stolen":   true,
		"stolen_date": stolenDate,
		"status":      "stolen",
	}, &models.Device{}, nil)
}

// MarkAsRecovered marks device as recovered
func (r *DeviceRepositoryImpl) MarkAsRecovered(ctx context.Context, deviceID uuid.UUID) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"is_stolen":      false,
		"recovered_date": time.Now(),
		"status":         "active",
	}, &models.Device{}, nil)
}

// MarkAsLost marks device as lost
func (r *DeviceRepositoryImpl) MarkAsLost(ctx context.Context, deviceID uuid.UUID, lostDate time.Time) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"is_lost":   true,
		"lost_date": lostDate,
		"status":    "lost",
	}, &models.Device{}, nil)
}

// MarkAsFound marks device as found
func (r *DeviceRepositoryImpl) MarkAsFound(ctx context.Context, deviceID uuid.UUID) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"is_lost":    false,
		"found_date": time.Now(),
		"status":     "active",
	}, &models.Device{}, nil)
}

// === Ownership Management ===

// TransferOwnership transfers device ownership
func (r *DeviceRepositoryImpl) TransferOwnership(ctx context.Context, deviceID, newOwnerID uuid.UUID) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"owner_id":                newOwnerID,
		"ownership_transfer_date": time.Now(),
	}, &models.Device{}, nil)
}

// UpdateOwnershipHistory updates ownership history
func (r *DeviceRepositoryImpl) UpdateOwnershipHistory(ctx context.Context, deviceID uuid.UUID, history interface{}) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"ownership_history": history,
	}, &models.Device{}, nil)
}

// GetOwnershipHistory gets ownership history
func (r *DeviceRepositoryImpl) GetOwnershipHistory(ctx context.Context, deviceID uuid.UUID) ([]interface{}, error) {
	var device models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Select("ownership_history").Where("id = ?", deviceID).First(&device).Error; err != nil {
		return nil, err
	}

	// Parse JSON field if needed
	return []interface{}{}, nil
}

// === Value & Depreciation ===

// UpdateCurrentValue updates current value
func (r *DeviceRepositoryImpl) UpdateCurrentValue(ctx context.Context, deviceID uuid.UUID, value decimal.Decimal) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"current_value":       value,
		"last_valuation_date": time.Now(),
	}, &models.Device{}, nil)
}

// UpdateMarketValue updates market value
func (r *DeviceRepositoryImpl) UpdateMarketValue(ctx context.Context, deviceID uuid.UUID, value decimal.Decimal) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"market_value": value,
	}, &models.Device{}, nil)
}

// CalculateDepreciation calculates depreciation
func (r *DeviceRepositoryImpl) CalculateDepreciation(ctx context.Context, deviceID uuid.UUID) (decimal.Decimal, error) {
	var device models.Device
	if err := r.FindByID(ctx, deviceID, &device, nil); err != nil {
		return decimal.Zero, err
	}

	// Calculate depreciation based on age and original price
	ageInMonths := int(time.Since(device.PurchaseDate).Hours() / 24 / 30)
	depreciationRate := decimal.NewFromFloat(0.02) // 2% per month
	depreciation := device.PurchasePrice.Mul(depreciationRate).Mul(decimal.NewFromInt(int64(ageInMonths)))

	currentValue := device.PurchasePrice.Sub(depreciation)
	if currentValue.LessThan(decimal.Zero) {
		currentValue = decimal.Zero
	}

	// Update current value
	r.UpdateCurrentValue(ctx, deviceID, currentValue)

	return currentValue, nil
}

// === Insurance Operations ===

// UpdateInsuranceStatus updates insurance status
func (r *DeviceRepositoryImpl) UpdateInsuranceStatus(ctx context.Context, deviceID uuid.UUID, isInsured bool) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"is_insured": isInsured,
	}, &models.Device{}, nil)
}

// UpdateActivePolicyID updates active policy ID
func (r *DeviceRepositoryImpl) UpdateActivePolicyID(ctx context.Context, deviceID uuid.UUID, policyID *uuid.UUID) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"active_policy_id": policyID,
	}, &models.Device{}, nil)
}

// GetInsuredDevices gets all insured devices
func (r *DeviceRepositoryImpl) GetInsuredDevices(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Where("is_insured = ?", true).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// GetUninsuredDevices gets all uninsured devices
func (r *DeviceRepositoryImpl) GetUninsuredDevices(ctx context.Context, ownerID *uuid.UUID) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	query := db.Where("is_insured = ?", false)
	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}

	if err := query.Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// GetEligibleForInsurance gets devices eligible for insurance
func (r *DeviceRepositoryImpl) GetEligibleForInsurance(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	// Devices that are not insured, not stolen, not blacklisted, and within age limit
	maxAge := time.Now().AddDate(-5, 0, 0) // 5 years old

	if err := db.Where("is_insured = ? AND is_stolen = ? AND blacklist_status = ? AND purchase_date > ?",
		false, false, false, maxAge).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// === Risk Assessment ===

// UpdateRiskScore updates risk score
func (r *DeviceRepositoryImpl) UpdateRiskScore(ctx context.Context, deviceID uuid.UUID, score float64) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"risk_score":           score,
		"risk_assessment_date": time.Now(),
	}, &models.Device{}, nil)
}

// GetHighRiskDevices gets high risk devices
func (r *DeviceRepositoryImpl) GetHighRiskDevices(ctx context.Context, threshold float64) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Where("risk_score >= ?", threshold).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// CalculateRiskScore calculates risk score for a device
func (r *DeviceRepositoryImpl) CalculateRiskScore(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	var device models.Device
	if err := r.FindByID(ctx, deviceID, &device, nil); err != nil {
		return 0, err
	}

	// Simple risk calculation based on various factors
	riskScore := 50.0 // Base score

	// Age factor
	ageInMonths := int(time.Since(device.PurchaseDate).Hours() / 24 / 30)
	riskScore += float64(ageInMonths) * 0.5

	// Condition factor
	if device.Condition == "poor" {
		riskScore += 20
	} else if device.Condition == "fair" {
		riskScore += 10
	}

	// Previous claims
	if device.ClaimCount > 0 {
		riskScore += float64(device.ClaimCount) * 10
	}

	// Update the risk score
	r.UpdateRiskScore(ctx, deviceID, riskScore)

	return riskScore, nil
}

// === Blacklist Management ===

// AddToBlacklist adds device to blacklist
func (r *DeviceRepositoryImpl) AddToBlacklist(ctx context.Context, deviceID uuid.UUID, reason string) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"blacklist_status": true,
		"blacklist_reason": reason,
		"blacklisted_at":   time.Now(),
	}, &models.Device{}, nil)
}

// RemoveFromBlacklist removes device from blacklist
func (r *DeviceRepositoryImpl) RemoveFromBlacklist(ctx context.Context, deviceID uuid.UUID) error {
	return r.UpdateFields(ctx, deviceID, map[string]interface{}{
		"blacklist_status": false,
		"blacklist_reason": nil,
		"blacklisted_at":   nil,
	}, &models.Device{}, nil)
}

// GetBlacklistedDevices gets all blacklisted devices
func (r *DeviceRepositoryImpl) GetBlacklistedDevices(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Where("blacklist_status = ?", true).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// CheckBlacklistStatus checks if device is blacklisted
func (r *DeviceRepositoryImpl) CheckBlacklistStatus(ctx context.Context, imei string) (bool, string, error) {
	var device models.Device
	db := r.GetDB(ctx, nil)

	if err := db.Select("blacklist_status", "blacklist_reason").Where("imei = ?", imei).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, "", nil
		}
		return false, "", err
	}

	return device.BlacklistStatus, device.BlacklistReason, nil
}
