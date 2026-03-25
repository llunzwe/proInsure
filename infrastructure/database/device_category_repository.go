package database

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/application/services"
	"smartsure/internal/domain/models/device/categories/base"
)

// DeviceCategoryRepository implements the repository interface for device categories
type DeviceCategoryRepository struct {
	db *gorm.DB
}

// NewDeviceCategoryRepository creates a new device category repository
func NewDeviceCategoryRepository(db *gorm.DB) *DeviceCategoryRepository {
	return &DeviceCategoryRepository{
		db: db,
	}
}

// DeviceCategoryModel represents the database model for devices with categories
type DeviceCategoryModel struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key"`
	Category         string     `gorm:"type:varchar(50);not null;index"`
	Specifications   JSONB      `gorm:"type:jsonb;not null"`
	IMEI             string     `gorm:"type:varchar(15);unique;index"`
	SerialNumber     string     `gorm:"type:varchar(50);unique;index"`
	OwnerID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	Status           string     `gorm:"type:varchar(20);not null;default:'active'"`
	InsuranceProfile *JSONB     `gorm:"type:jsonb"`
	RiskAssessment   *JSONB     `gorm:"type:jsonb"`
	CreatedAt        time.Time  `gorm:"not null"`
	UpdatedAt        time.Time  `gorm:"not null"`
	DeletedAt        *time.Time `gorm:"index"`

	// Category-specific fields (stored in JSONB but indexed for performance)
	Manufacturer      string  `gorm:"type:varchar(100);index"`
	Model             string  `gorm:"type:varchar(100);index"`
	MarketValue       float64 `gorm:"type:decimal(10,2);index"`
	RiskScore         float64 `gorm:"type:decimal(5,2);index"`
	EligibilityStatus string  `gorm:"type:varchar(20);index"`
}

// TableName specifies the table name for the model
func (DeviceCategoryModel) TableName() string {
	return "device_categories"
}

// JSONB represents a JSONB database type
type JSONB json.RawMessage

// Scan implements the Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB("null")
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = JSONB(v)
	case string:
		*j = JSONB(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// Value implements the driver Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}

// Create creates a new device in the database
func (r *DeviceCategoryRepository) Create(ctx context.Context, device *services.DeviceEntity) error {
	model := r.entityToModel(device)

	// Extract indexed fields from specifications for better query performance
	r.extractIndexedFields(device, model)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	return nil
}

// Update updates an existing device in the database
func (r *DeviceCategoryRepository) Update(ctx context.Context, device *services.DeviceEntity) error {
	model := r.entityToModel(device)
	r.extractIndexedFields(device, model)

	result := r.db.WithContext(ctx).Model(&DeviceCategoryModel{}).
		Where("id = ?", device.ID).
		Updates(model)

	if result.Error != nil {
		return fmt.Errorf("failed to update device: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

// FindByID finds a device by its ID
func (r *DeviceCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*services.DeviceEntity, error) {
	var model DeviceCategoryModel

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device not found")
		}
		return nil, fmt.Errorf("failed to find device: %w", err)
	}

	return r.modelToEntity(&model), nil
}

// FindByIMEI finds a device by its IMEI
func (r *DeviceCategoryRepository) FindByIMEI(ctx context.Context, imei string) (*services.DeviceEntity, error) {
	var model DeviceCategoryModel

	if err := r.db.WithContext(ctx).
		Where("imei = ? AND deleted_at IS NULL", imei).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device not found")
		}
		return nil, fmt.Errorf("failed to find device: %w", err)
	}

	return r.modelToEntity(&model), nil
}

// List lists devices based on filters
func (r *DeviceCategoryRepository) List(ctx context.Context, filters map[string]interface{}) ([]*services.DeviceEntity, error) {
	query := r.db.WithContext(ctx).Model(&DeviceCategoryModel{}).
		Where("deleted_at IS NULL")

	// Apply filters
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	if ownerID, ok := filters["owner_id"].(uuid.UUID); ok {
		query = query.Where("owner_id = ?", ownerID)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if manufacturer, ok := filters["manufacturer"].(string); ok && manufacturer != "" {
		query = query.Where("manufacturer = ?", manufacturer)
	}

	if minValue, ok := filters["min_value"].(float64); ok {
		query = query.Where("market_value >= ?", minValue)
	}

	if maxValue, ok := filters["max_value"].(float64); ok {
		query = query.Where("market_value <= ?", maxValue)
	}

	if minRisk, ok := filters["min_risk_score"].(float64); ok {
		query = query.Where("risk_score >= ?", minRisk)
	}

	if maxRisk, ok := filters["max_risk_score"].(float64); ok {
		query = query.Where("risk_score <= ?", maxRisk)
	}

	// Apply pagination if provided
	if limit, ok := filters["limit"].(int); ok && limit > 0 {
		query = query.Limit(limit)
	}

	if offset, ok := filters["offset"].(int); ok && offset > 0 {
		query = query.Offset(offset)
	}

	// Apply sorting
	if sortBy, ok := filters["sort_by"].(string); ok && sortBy != "" {
		sortOrder := "asc"
		if order, ok := filters["sort_order"].(string); ok && order != "" {
			sortOrder = order
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	} else {
		query = query.Order("created_at DESC")
	}

	var models []DeviceCategoryModel
	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	entities := make([]*services.DeviceEntity, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}

	return entities, nil
}

// Delete soft deletes a device
func (r *DeviceCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&DeviceCategoryModel{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to delete device: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

// Helper methods

func (r *DeviceCategoryRepository) entityToModel(entity *services.DeviceEntity) *DeviceCategoryModel {
	model := &DeviceCategoryModel{
		ID:             entity.ID,
		Category:       string(entity.Category),
		Specifications: JSONB(entity.Specifications),
		IMEI:           entity.IMEI,
		SerialNumber:   entity.SerialNumber,
		OwnerID:        entity.OwnerID,
		Status:         entity.Status,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}

	if entity.InsuranceProfile != nil {
		profileJSON, _ := json.Marshal(entity.InsuranceProfile)
		profile := JSONB(profileJSON)
		model.InsuranceProfile = &profile
	}

	if entity.RiskAssessment != nil {
		assessmentJSON, _ := json.Marshal(entity.RiskAssessment)
		assessment := JSONB(assessmentJSON)
		model.RiskAssessment = &assessment
		model.RiskScore = entity.RiskAssessment.Score
	}

	return model
}

func (r *DeviceCategoryRepository) modelToEntity(model *DeviceCategoryModel) *services.DeviceEntity {
	entity := &services.DeviceEntity{
		ID:             model.ID,
		Category:       base.CategoryType(model.Category),
		Specifications: json.RawMessage(model.Specifications),
		IMEI:           model.IMEI,
		SerialNumber:   model.SerialNumber,
		OwnerID:        model.OwnerID,
		Status:         model.Status,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}

	if model.InsuranceProfile != nil {
		var profile base.InsuranceProfile
		if err := json.Unmarshal([]byte(*model.InsuranceProfile), &profile); err == nil {
			entity.InsuranceProfile = &profile
		}
	}

	if model.RiskAssessment != nil {
		var assessment base.RiskAssessment
		if err := json.Unmarshal([]byte(*model.RiskAssessment), &assessment); err == nil {
			entity.RiskAssessment = &assessment
		}
	}

	return entity
}

// extractIndexedFields extracts frequently queried fields from JSONB for indexing
func (r *DeviceCategoryRepository) extractIndexedFields(entity *services.DeviceEntity, model *DeviceCategoryModel) {
	// Parse specifications to extract indexed fields
	var specData map[string]interface{}
	if err := json.Unmarshal(entity.Specifications, &specData); err == nil {
		if manufacturer, ok := specData["manufacturer"].(string); ok {
			model.Manufacturer = manufacturer
		}
		if modelName, ok := specData["model"].(string); ok {
			model.Model = modelName
		}
		if marketValue, ok := specData["market_value"].(float64); ok {
			model.MarketValue = marketValue
		}
	}

	// Extract from insurance profile
	if entity.InsuranceProfile != nil {
		model.EligibilityStatus = entity.InsuranceProfile.EligibilityStatus
	}

	// Extract from risk assessment
	if entity.RiskAssessment != nil {
		model.RiskScore = entity.RiskAssessment.Score
	}
}

// Migration creates the device_categories table
func (r *DeviceCategoryRepository) Migrate() error {
	if err := r.db.AutoMigrate(&DeviceCategoryModel{}); err != nil {
		return fmt.Errorf("failed to migrate device_categories table: %w", err)
	}

	// Create additional indexes for better performance
	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_device_categories_category_status 
		ON device_categories(category, status) WHERE deleted_at IS NULL`)

	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_device_categories_owner_category 
		ON device_categories(owner_id, category) WHERE deleted_at IS NULL`)

	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_device_categories_risk_score 
		ON device_categories(risk_score) WHERE deleted_at IS NULL AND risk_score IS NOT NULL`)

	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_device_categories_market_value 
		ON device_categories(market_value) WHERE deleted_at IS NULL AND market_value IS NOT NULL`)

	// Create GIN index for JSONB columns for better query performance
	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_device_categories_specifications 
		ON device_categories USING gin(specifications)`)

	return nil
}
