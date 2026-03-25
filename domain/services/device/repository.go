package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smartsure/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository implements DeviceRepository interface
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new device repository
func NewRepository(db *gorm.DB) DeviceRepository {
	return &Repository{
		db: db,
	}
}

// Create creates a new device
func (r *Repository) Create(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

// GetByID retrieves a device by ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("CorporateAccount").
		Preload("Policies").
		Preload("Claims").
		First(&device, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("device not found: %w", err)
		}
		return nil, err
	}

	return &device, nil
}

// GetByIMEI retrieves a device by IMEI
func (r *Repository) GetByIMEI(ctx context.Context, imei string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("CorporateAccount").
		Where("imei = ?", imei).
		First(&device).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("device not found with IMEI %s: %w", imei, err)
		}
		return nil, err
	}

	return &device, nil
}

// GetBySerialNumber retrieves a device by serial number
func (r *Repository) GetBySerialNumber(ctx context.Context, serialNumber string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("CorporateAccount").
		Where("serial_number = ?", serialNumber).
		First(&device).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("device not found with serial number %s: %w", serialNumber, err)
		}
		return nil, err
	}

	return &device, nil
}

// Update updates an existing device
func (r *Repository) Update(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

// Delete soft deletes a device
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Device{}, "id = ?", id).Error
}

// List lists devices with filters and pagination
func (r *Repository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*models.Device, int64, error) {
	var devices []*models.Device
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Device{})

	// Apply filters
	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "owner_id":
			query = query.Where("owner_id = ?", value)
		case "corporate_account_id":
			query = query.Where("corporate_account_id = ?", value)
		case "brand":
			query = query.Where("brand = ?", value)
		case "model":
			query = query.Where("model LIKE ?", fmt.Sprintf("%%%v%%", value))
		case "category":
			query = query.Where("category = ?", value)
		case "condition":
			query = query.Where("condition = ?", value)
		case "is_stolen":
			query = query.Where("is_stolen = ?", value)
		case "is_verified":
			query = query.Where("is_verified = ?", value)
		case "risk_score_min":
			query = query.Where("risk_score >= ?", value)
		case "risk_score_max":
			query = query.Where("risk_score <= ?", value)
		case "created_after":
			query = query.Where("created_at >= ?", value)
		case "created_before":
			query = query.Where("created_at <= ?", value)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch
	err := query.
		Preload("Owner").
		Preload("CorporateAccount").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&devices).Error

	return devices, total, err
}

// Search searches devices by query string
func (r *Repository) Search(ctx context.Context, query string, limit, offset int) ([]*models.Device, int64, error) {
	var devices []*models.Device
	var total int64

	searchQuery := r.db.WithContext(ctx).Model(&models.Device{}).
		Where("imei LIKE ? OR serial_number LIKE ? OR model LIKE ? OR brand LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	// Count total
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch
	err := searchQuery.
		Preload("Owner").
		Preload("CorporateAccount").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&devices).Error

	return devices, total, err
}

// GetByOwnerID retrieves all devices for a specific owner
func (r *Repository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*models.Device, error) {
	var devices []*models.Device
	err := r.db.WithContext(ctx).
		Preload("Policies").
		Preload("Claims").
		Where("owner_id = ?", ownerID).
		Find(&devices).Error

	return devices, err
}

// GetByCorporateAccountID retrieves all devices for a corporate account
func (r *Repository) GetByCorporateAccountID(ctx context.Context, corporateID uuid.UUID) ([]*models.Device, error) {
	var devices []*models.Device
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("CorporateAssignment").
		Preload("BYODRegistration").
		Where("corporate_account_id = ?", corporateID).
		Find(&devices).Error

	return devices, err
}

// BulkCreate creates multiple devices
func (r *Repository) BulkCreate(ctx context.Context, devices []*models.Device) error {
	return r.db.WithContext(ctx).CreateInBatches(devices, 100).Error
}

// BulkUpdate updates multiple devices
func (r *Repository) BulkUpdate(ctx context.Context, devices []*models.Device) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, device := range devices {
			if err := tx.Save(device).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkUpdateStatus updates status for multiple devices
func (r *Repository) BulkUpdateStatus(ctx context.Context, deviceIDs []uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("id IN ?", deviceIDs).
		Update("status", status).Error
}

// GetExpiredWarrantyDevices retrieves devices with expired or expiring warranty
func (r *Repository) GetExpiredWarrantyDevices(ctx context.Context, days int) ([]*models.Device, error) {
	var devices []*models.Device

	expiryDate := time.Now().AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("warranty_expiry <= ? AND warranty_expiry IS NOT NULL", expiryDate).
		Find(&devices).Error

	return devices, err
}

// GetDevicesNeedingInspection retrieves devices that need inspection
func (r *Repository) GetDevicesNeedingInspection(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("last_inspection IS NULL OR last_inspection < ?", thirtyDaysAgo).
		Where("status = ?", "active").
		Find(&devices).Error

	return devices, err
}

// GetStolenDevices retrieves all stolen devices
func (r *Repository) GetStolenDevices(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("is_stolen = ?", true).
		Order("stolen_date DESC").
		Find(&devices).Error

	return devices, err
}

// GetBlacklistedDevices retrieves all blacklisted devices
func (r *Repository) GetBlacklistedDevices(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("blacklist_status IN ?", []string{"blocked", "blacklisted"}).
		Find(&devices).Error

	return devices, err
}

// GetDevicesByRiskLevel retrieves devices by risk level
func (r *Repository) GetDevicesByRiskLevel(ctx context.Context, riskLevel string) ([]*models.Device, error) {
	var devices []*models.Device

	// Define risk score ranges
	var minScore, maxScore float64
	switch riskLevel {
	case "low":
		minScore, maxScore = 0, 25
	case "medium":
		minScore, maxScore = 25, 50
	case "high":
		minScore, maxScore = 50, 75
	case "very_high":
		minScore, maxScore = 75, 100
	default:
		return nil, fmt.Errorf("invalid risk level: %s", riskLevel)
	}

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("RiskProfile").
		Where("risk_score >= ? AND risk_score <= ?", minScore, maxScore).
		Find(&devices).Error

	return devices, err
}

// GetUnverifiedDevices retrieves devices that haven't been verified in specified days
func (r *Repository) GetUnverifiedDevices(ctx context.Context, days int) ([]*models.Device, error) {
	var devices []*models.Device

	cutoffDate := time.Now().AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("is_verified = ? OR (last_verified_at IS NULL OR last_verified_at < ?)", false, cutoffDate).
		Find(&devices).Error

	return devices, err
}

// Additional helper methods

// GetActiveDevicesCount returns the count of active devices
func (r *Repository) GetActiveDevicesCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("status = ?", "active").
		Count(&count).Error

	return count, err
}

// GetDevicesByCondition retrieves devices by condition
func (r *Repository) GetDevicesByCondition(ctx context.Context, condition string) ([]*models.Device, error) {
	var devices []*models.Device

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("condition = ?", condition).
		Find(&devices).Error

	return devices, err
}

// GetHighValueDevices retrieves devices above a certain value threshold
func (r *Repository) GetHighValueDevices(ctx context.Context, threshold float64) ([]*models.Device, error) {
	var devices []*models.Device

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("CorporateAccount").
		Where("current_value >= ?", threshold).
		Order("current_value DESC").
		Find(&devices).Error

	return devices, err
}

// GetRecentlyRegisteredDevices retrieves devices registered in the last N days
func (r *Repository) GetRecentlyRegisteredDevices(ctx context.Context, days int) ([]*models.Device, error) {
	var devices []*models.Device

	sinceDate := time.Now().AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("registration_date >= ?", sinceDate).
		Order("registration_date DESC").
		Find(&devices).Error

	return devices, err
}

// GetDevicesWithActivePolicies retrieves devices with active insurance policies
func (r *Repository) GetDevicesWithActivePolicies(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("Policies", "status = ?", "active").
		Joins("JOIN policies ON policies.device_id = devices.id").
		Where("policies.status = ?", "active").
		Distinct().
		Find(&devices).Error

	return devices, err
}

// GetDevicesWithPendingClaims retrieves devices with pending claims
func (r *Repository) GetDevicesWithPendingClaims(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device

	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("Claims", "status = ?", "pending").
		Joins("JOIN claims ON claims.device_id = devices.id").
		Where("claims.status = ?", "pending").
		Distinct().
		Find(&devices).Error

	return devices, err
}
