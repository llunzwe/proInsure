package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
)

// UserRepositoryImpl implements the UserRepository interface using PostgreSQL
type UserRepositoryImpl struct {
	*BaseRepository
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl
func NewUserRepository(db *gorm.DB, logger Logger) repositories.UserRepository {
	return &UserRepositoryImpl{
		BaseRepository: NewBaseRepository(db, logger),
		db:             db,
	}
}

// === Basic CRUD Operations ===

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Create(ctx, user, nil)
}

// GetByID retrieves a user by ID
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.FindByID(ctx, id, &user, &QueryOptions{
		Preload: []string{"Devices", "Policies", "Claims", "PaymentMethods"},
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByPhone retrieves a user by phone number
func (r *UserRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Update(ctx, user, nil)
}

// Delete hard deletes a user
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.HardDelete(ctx, id, &models.User{}, nil)
}

// SoftDelete soft deletes a user
func (r *UserRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.Delete(ctx, id, &models.User{}, nil)
}

// Restore restores a soft-deleted user
func (r *UserRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	db := r.GetDB(ctx, nil)
	result := db.Model(&models.User{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// === Batch Operations ===

// CreateBatch creates multiple users
func (r *UserRepositoryImpl) CreateBatch(ctx context.Context, users []*models.User) error {
	return r.BaseRepository.BulkCreate(ctx, users, nil)
}

// UpdateBatch updates multiple users
func (r *UserRepositoryImpl) UpdateBatch(ctx context.Context, users []*models.User) error {
	tx, err := r.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			r.RollbackTransaction(tx)
		}
	}()

	for _, user := range users {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
	}

	return r.CommitTransaction(tx)
}

// DeleteBatch deletes multiple users
func (r *UserRepositoryImpl) DeleteBatch(ctx context.Context, ids []uuid.UUID) error {
	db := r.GetDB(ctx, nil)
	return db.Where("id IN ?", ids).Delete(&models.User{}).Error
}

// GetByIDs retrieves multiple users by IDs
func (r *UserRepositoryImpl) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// === Search & Filtering ===

// Search searches users based on filters
func (r *UserRepositoryImpl) Search(ctx context.Context, filters repositories.UserSearchFilters) ([]*models.User, int64, error) {
	var users []*models.User
	var totalCount int64

	db := r.GetDB(ctx, nil)
	query := r.buildSearchQuery(db.Model(&models.User{}), filters)

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
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// buildSearchQuery builds the search query based on filters
func (r *UserRepositoryImpl) buildSearchQuery(query *gorm.DB, filters repositories.UserSearchFilters) *gorm.DB {
	// Basic filters
	if filters.Name != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ?",
			"%"+filters.Name+"%", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filters.Email+"%")
	}
	if filters.Phone != "" {
		query = query.Where("phone_number ILIKE ?", "%"+filters.Phone+"%")
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Segment != "" {
		query = query.Where("customer_segment = ?", filters.Segment)
	}

	// Location filters
	if filters.Country != "" {
		query = query.Where("country = ?", filters.Country)
	}
	if filters.City != "" {
		query = query.Where("city = ?", filters.City)
	}

	// Risk filters
	if filters.MinRiskScore > 0 {
		query = query.Where("risk_score >= ?", filters.MinRiskScore)
	}
	if filters.MaxRiskScore > 0 {
		query = query.Where("risk_score <= ?", filters.MaxRiskScore)
	}
	if filters.BlacklistStatus != nil {
		query = query.Where("blacklist_status = ?", *filters.BlacklistStatus)
	}

	// Financial filters
	if filters.MinCreditScore > 0 {
		query = query.Where("credit_score >= ?", filters.MinCreditScore)
	}
	if filters.MaxCreditScore > 0 {
		query = query.Where("credit_score <= ?", filters.MaxCreditScore)
	}
	if filters.HasOutstandingBalance != nil && *filters.HasOutstandingBalance {
		query = query.Where("outstanding_balance > 0")
	}

	// Compliance filters
	if filters.KYCStatus != "" {
		query = query.Where("kyc_status = ?", filters.KYCStatus)
	}
	if filters.EmailVerified != nil {
		query = query.Where("email_verified = ?", *filters.EmailVerified)
	}
	if filters.PhoneVerified != nil {
		query = query.Where("phone_verified = ?", *filters.PhoneVerified)
	}

	// Date filters
	if !filters.CreatedAfter.IsZero() {
		query = query.Where("created_at >= ?", filters.CreatedAfter)
	}
	if !filters.CreatedBefore.IsZero() {
		query = query.Where("created_at <= ?", filters.CreatedBefore)
	}

	// Family/Group filters
	if filters.HouseholdID != nil {
		query = query.Where("household_id = ?", *filters.HouseholdID)
	}
	if filters.CorporateAccountID != nil {
		query = query.Where("corporate_account_id = ?", *filters.CorporateAccountID)
	}

	return query
}

// SearchByName searches users by name
func (r *UserRepositoryImpl) SearchByName(ctx context.Context, name string, limit int) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	query := db.Where("first_name ILIKE ? OR last_name ILIKE ? OR CONCAT(first_name, ' ', last_name) ILIKE ?",
		"%"+name+"%", "%"+name+"%", "%"+name+"%")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// SearchByReferralCode searches for a user by referral code
func (r *UserRepositoryImpl) SearchByReferralCode(ctx context.Context, code string) (*models.User, error) {
	var user models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("referral_code = ?", code).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByHouseholdID gets all users in a household
func (r *UserRepositoryImpl) GetByHouseholdID(ctx context.Context, householdID uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("household_id = ?", householdID).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetByCorporateAccountID gets all users under a corporate account
func (r *UserRepositoryImpl) GetByCorporateAccountID(ctx context.Context, corporateID uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("corporate_account_id = ?", corporateID).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetBySegment gets users by customer segment
func (r *UserRepositoryImpl) GetBySegment(ctx context.Context, segment string) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("customer_segment = ?", segment).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetByStatus gets users by status
func (r *UserRepositoryImpl) GetByStatus(ctx context.Context, status string, limit int) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	query := db.Where("status = ?", status)
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
