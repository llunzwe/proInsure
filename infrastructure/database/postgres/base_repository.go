package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BaseRepository provides common CRUD operations for all entities
type BaseRepository struct {
	db     *gorm.DB
	logger Logger
}

// Logger interface for logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository(db *gorm.DB, logger Logger) *BaseRepository {
	return &BaseRepository{
		db:     db,
		logger: logger,
	}
}

// QueryOptions defines common query options
type QueryOptions struct {
	Preload        []string
	Select         []string
	Where          map[string]interface{}
	Order          string
	Limit          int
	Offset         int
	IncludeDeleted bool
	ForUpdate      bool
	Transaction    *gorm.DB
}

// GetDB returns the appropriate database connection (transaction or main)
func (r *BaseRepository) GetDB(ctx context.Context, opts *QueryOptions) *gorm.DB {
	var db *gorm.DB

	// Use transaction if provided
	if opts != nil && opts.Transaction != nil {
		db = opts.Transaction
	} else {
		db = r.db
	}

	// Add context
	db = db.WithContext(ctx)

	// Apply query options
	if opts != nil {
		if opts.IncludeDeleted {
			db = db.Unscoped()
		}

		if opts.ForUpdate {
			db = db.Clauses(clause.Locking{Strength: "UPDATE"})
		}

		for _, preload := range opts.Preload {
			db = db.Preload(preload)
		}

		if len(opts.Select) > 0 {
			db = db.Select(opts.Select)
		}

		for key, value := range opts.Where {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}

		if opts.Order != "" {
			db = db.Order(opts.Order)
		}

		if opts.Limit > 0 {
			db = db.Limit(opts.Limit)
		}

		if opts.Offset > 0 {
			db = db.Offset(opts.Offset)
		}
	}

	return db
}

// Create creates a new entity
func (r *BaseRepository) Create(ctx context.Context, entity interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	// Set created_at if the entity has this field
	r.setTimestamp(entity, "CreatedAt", time.Now())

	// Set ID if it's a UUID and not already set
	r.setUUIDIfEmpty(entity, "ID")

	if err := db.Create(entity).Error; err != nil {
		r.logger.Error("Failed to create entity", err, "entity", reflect.TypeOf(entity).String())
		return r.translateError(err)
	}

	r.logger.Info("Entity created successfully", "entity", reflect.TypeOf(entity).String())
	return nil
}

// Update updates an existing entity
func (r *BaseRepository) Update(ctx context.Context, entity interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	// Set updated_at if the entity has this field
	r.setTimestamp(entity, "UpdatedAt", time.Now())

	result := db.Save(entity)
	if result.Error != nil {
		r.logger.Error("Failed to update entity", result.Error, "entity", reflect.TypeOf(entity).String())
		return r.translateError(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	r.logger.Info("Entity updated successfully", "entity", reflect.TypeOf(entity).String())
	return nil
}

// UpdateFields updates specific fields of an entity
func (r *BaseRepository) UpdateFields(ctx context.Context, id uuid.UUID, updates map[string]interface{}, model interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	// Add updated_at
	updates["updated_at"] = time.Now()

	result := db.Model(model).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		r.logger.Error("Failed to update fields", result.Error, "id", id)
		return r.translateError(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete performs a soft delete on an entity
func (r *BaseRepository) Delete(ctx context.Context, id uuid.UUID, model interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	result := db.Where("id = ?", id).Delete(model)
	if result.Error != nil {
		r.logger.Error("Failed to delete entity", result.Error, "id", id)
		return r.translateError(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	r.logger.Info("Entity deleted successfully", "id", id)
	return nil
}

// HardDelete performs a hard delete on an entity
func (r *BaseRepository) HardDelete(ctx context.Context, id uuid.UUID, model interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	result := db.Unscoped().Where("id = ?", id).Delete(model)
	if result.Error != nil {
		r.logger.Error("Failed to hard delete entity", result.Error, "id", id)
		return r.translateError(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	r.logger.Info("Entity hard deleted successfully", "id", id)
	return nil
}

// FindByID finds an entity by ID
func (r *BaseRepository) FindByID(ctx context.Context, id uuid.UUID, entity interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	if err := db.First(entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		r.logger.Error("Failed to find entity by ID", err, "id", id)
		return r.translateError(err)
	}

	return nil
}

// FindAll finds all entities matching the criteria
func (r *BaseRepository) FindAll(ctx context.Context, entities interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	if err := db.Find(entities).Error; err != nil {
		r.logger.Error("Failed to find entities", err)
		return r.translateError(err)
	}

	return nil
}

// Count counts entities matching the criteria
func (r *BaseRepository) Count(ctx context.Context, model interface{}, opts *QueryOptions) (int64, error) {
	db := r.GetDB(ctx, opts)

	var count int64
	if err := db.Model(model).Count(&count).Error; err != nil {
		r.logger.Error("Failed to count entities", err)
		return 0, r.translateError(err)
	}

	return count, nil
}

// Exists checks if an entity exists
func (r *BaseRepository) Exists(ctx context.Context, id uuid.UUID, model interface{}, opts *QueryOptions) (bool, error) {
	db := r.GetDB(ctx, opts)

	var count int64
	if err := db.Model(model).Where("id = ?", id).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check existence", err, "id", id)
		return false, r.translateError(err)
	}

	return count > 0, nil
}

// BeginTransaction starts a new transaction
func (r *BaseRepository) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		r.logger.Error("Failed to begin transaction", tx.Error)
		return nil, tx.Error
	}
	return tx, nil
}

// CommitTransaction commits a transaction
func (r *BaseRepository) CommitTransaction(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("Failed to commit transaction", err)
		return err
	}
	return nil
}

// RollbackTransaction rollbacks a transaction
func (r *BaseRepository) RollbackTransaction(tx *gorm.DB) error {
	if err := tx.Rollback().Error; err != nil {
		r.logger.Error("Failed to rollback transaction", err)
		return err
	}
	return nil
}

// Paginate applies pagination to a query
func (r *BaseRepository) Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// BulkCreate creates multiple entities
func (r *BaseRepository) BulkCreate(ctx context.Context, entities interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	// Set timestamps and IDs for each entity
	r.prepareBulkCreate(entities)

	if err := db.CreateInBatches(entities, 100).Error; err != nil {
		r.logger.Error("Failed to bulk create entities", err)
		return r.translateError(err)
	}

	r.logger.Info("Bulk create successful", "count", reflect.ValueOf(entities).Len())
	return nil
}

// BulkUpdate performs bulk updates
func (r *BaseRepository) BulkUpdate(ctx context.Context, ids []uuid.UUID, updates map[string]interface{}, model interface{}, opts *QueryOptions) error {
	db := r.GetDB(ctx, opts)

	// Add updated_at
	updates["updated_at"] = time.Now()

	result := db.Model(model).Where("id IN ?", ids).Updates(updates)
	if result.Error != nil {
		r.logger.Error("Failed to bulk update", result.Error)
		return r.translateError(result.Error)
	}

	r.logger.Info("Bulk update successful", "affected", result.RowsAffected)
	return nil
}

// Helper methods

func (r *BaseRepository) setTimestamp(entity interface{}, fieldName string, value time.Time) {
	v := reflect.ValueOf(entity).Elem()
	field := v.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(value))
	}
}

func (r *BaseRepository) setUUIDIfEmpty(entity interface{}, fieldName string) {
	v := reflect.ValueOf(entity).Elem()
	field := v.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() && field.Type() == reflect.TypeOf(uuid.UUID{}) {
		if field.Interface().(uuid.UUID) == uuid.Nil {
			field.Set(reflect.ValueOf(uuid.New()))
		}
	}
}

func (r *BaseRepository) prepareBulkCreate(entities interface{}) {
	v := reflect.ValueOf(entities)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return
	}

	now := time.Now()
	for i := 0; i < v.Len(); i++ {
		entity := v.Index(i)
		if entity.Kind() == reflect.Ptr {
			entity = entity.Elem()
		}

		// Set timestamps
		if createdAt := entity.FieldByName("CreatedAt"); createdAt.IsValid() && createdAt.CanSet() {
			createdAt.Set(reflect.ValueOf(now))
		}
		if updatedAt := entity.FieldByName("UpdatedAt"); updatedAt.IsValid() && updatedAt.CanSet() {
			updatedAt.Set(reflect.ValueOf(now))
		}

		// Set UUID if empty
		if id := entity.FieldByName("ID"); id.IsValid() && id.CanSet() && id.Type() == reflect.TypeOf(uuid.UUID{}) {
			if id.Interface().(uuid.UUID) == uuid.Nil {
				id.Set(reflect.ValueOf(uuid.New()))
			}
		}
	}
}

func (r *BaseRepository) translateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	errStr := err.Error()
	if strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "UNIQUE constraint") {
		return ErrDuplicateEntry
	}
	if strings.Contains(errStr, "foreign key constraint") {
		return ErrForeignKeyViolation
	}
	if strings.Contains(errStr, "not null constraint") {
		return ErrNotNullViolation
	}

	return err
}

// Common errors
var (
	ErrNotFound            = errors.New("entity not found")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
	ErrNotNullViolation    = errors.New("not null constraint violation")
	ErrInvalidInput        = errors.New("invalid input")
	ErrTransactionFailed   = errors.New("transaction failed")
)

// SearchCriteria defines search parameters
type SearchCriteria struct {
	Query     string
	Fields    []string
	Filters   map[string]interface{}
	DateRange *DateRange
	NumRange  *NumericRange
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

// DateRange defines a date range filter
type DateRange struct {
	From  time.Time
	To    time.Time
	Field string
}

// NumericRange defines a numeric range filter
type NumericRange struct {
	Min   float64
	Max   float64
	Field string
}

// BuildSearchQuery builds a search query based on criteria
func (r *BaseRepository) BuildSearchQuery(db *gorm.DB, criteria *SearchCriteria) *gorm.DB {
	if criteria == nil {
		return db
	}

	// Apply text search
	if criteria.Query != "" && len(criteria.Fields) > 0 {
		var conditions []string
		var args []interface{}
		for _, field := range criteria.Fields {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+criteria.Query+"%")
		}
		db = db.Where(strings.Join(conditions, " OR "), args...)
	}

	// Apply filters
	for key, value := range criteria.Filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Apply date range
	if criteria.DateRange != nil {
		if !criteria.DateRange.From.IsZero() {
			db = db.Where(fmt.Sprintf("%s >= ?", criteria.DateRange.Field), criteria.DateRange.From)
		}
		if !criteria.DateRange.To.IsZero() {
			db = db.Where(fmt.Sprintf("%s <= ?", criteria.DateRange.Field), criteria.DateRange.To)
		}
	}

	// Apply numeric range
	if criteria.NumRange != nil {
		db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", criteria.NumRange.Field),
			criteria.NumRange.Min, criteria.NumRange.Max)
	}

	// Apply sorting
	if criteria.SortBy != "" {
		order := "ASC"
		if criteria.SortOrder != "" {
			order = criteria.SortOrder
		}
		db = db.Order(fmt.Sprintf("%s %s", criteria.SortBy, order))
	}

	// Apply pagination
	if criteria.Page > 0 && criteria.PageSize > 0 {
		db = db.Scopes(r.Paginate(criteria.Page, criteria.PageSize))
	}

	return db
}
