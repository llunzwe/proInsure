package postgres

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"smartsure/internal/domain/ports/repositories"
)

// RepositoryFactory creates and manages repository instances
type RepositoryFactory struct {
	db     *gorm.DB
	logger Logger
	mu     sync.RWMutex

	// Repository instances (singleton pattern)
	userRepo    repositories.UserRepository
	deviceRepo  repositories.DeviceRepository
	policyRepo  repositories.PolicyRepository
	claimRepo   repositories.ClaimRepository
	paymentRepo repositories.PaymentRepository

	// Add more repository instances as needed
	// corporateRepo repositories.CorporateRepository
	// repairRepo    repositories.RepairRepository
	documentRepo repositories.DocumentRepository
	// adminRepo     repositories.AdminRepository
	// partnerRepo   repositories.PartnerRepository
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB, logger Logger) *RepositoryFactory {
	return &RepositoryFactory{
		db:     db,
		logger: logger,
	}
}

// GetUserRepository returns the user repository instance
func (f *RepositoryFactory) GetUserRepository() repositories.UserRepository {
	f.mu.RLock()
	if f.userRepo != nil {
		f.mu.RUnlock()
		return f.userRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check after acquiring write lock
	if f.userRepo == nil {
		f.userRepo = NewUserRepository(f.db, f.logger)
	}

	return f.userRepo
}

// GetDeviceRepository returns the device repository instance
func (f *RepositoryFactory) GetDeviceRepository() repositories.DeviceRepository {
	f.mu.RLock()
	if f.deviceRepo != nil {
		f.mu.RUnlock()
		return f.deviceRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.deviceRepo == nil {
		f.deviceRepo = NewDeviceRepository(f.db, f.logger)
	}

	return f.deviceRepo
}

// GetPolicyRepository returns the policy repository instance
func (f *RepositoryFactory) GetPolicyRepository() repositories.PolicyRepository {
	f.mu.RLock()
	if f.policyRepo != nil {
		f.mu.RUnlock()
		return f.policyRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.policyRepo == nil {
		f.policyRepo = NewPolicyRepository(f.db, f.logger)
	}

	return f.policyRepo
}

// GetClaimRepository returns the claim repository instance
func (f *RepositoryFactory) GetClaimRepository() repositories.ClaimRepository {
	f.mu.RLock()
	if f.claimRepo != nil {
		f.mu.RUnlock()
		return f.claimRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.claimRepo == nil {
		f.claimRepo = NewClaimRepository(f.db, f.logger)
	}

	return f.claimRepo
}

// GetPaymentRepository returns the payment repository instance
func (f *RepositoryFactory) GetPaymentRepository() repositories.PaymentRepository {
	f.mu.RLock()
	if f.paymentRepo != nil {
		f.mu.RUnlock()
		return f.paymentRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.paymentRepo == nil {
		f.paymentRepo = NewPaymentRepository(f.db, f.logger)
	}

	return f.paymentRepo
}

// GetDocumentRepository returns the document repository instance
func (f *RepositoryFactory) GetDocumentRepository() repositories.DocumentRepository {
	f.mu.RLock()
	if f.documentRepo != nil {
		f.mu.RUnlock()
		return f.documentRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.documentRepo == nil {
		f.documentRepo = NewDocumentRepository(f.db, f.logger)
	}

	return f.documentRepo
}

// === Transaction Support ===

// WithTransaction executes a function within a database transaction
func (f *RepositoryFactory) WithTransaction(fn func(*RepositoryFactory) error) error {
	tx := f.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create a new factory with the transaction
	txFactory := &RepositoryFactory{
		db:     tx,
		logger: f.logger,
	}

	// Execute the function
	if err := fn(txFactory); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// === Repository Registry ===

// RepositoryRegistry manages all repositories in a centralized way
type RepositoryRegistry struct {
	factory *RepositoryFactory

	// Cached repository references for quick access
	User     repositories.UserRepository
	Device   repositories.DeviceRepository
	Policy   repositories.PolicyRepository
	Claim    repositories.ClaimRepository
	Payment  repositories.PaymentRepository
	Document repositories.DocumentRepository
}

// NewRepositoryRegistry creates a new repository registry
func NewRepositoryRegistry(db *gorm.DB, logger Logger) *RepositoryRegistry {
	factory := NewRepositoryFactory(db, logger)

	return &RepositoryRegistry{
		factory:  factory,
		User:     factory.GetUserRepository(),
		Device:   factory.GetDeviceRepository(),
		Policy:   factory.GetPolicyRepository(),
		Claim:    factory.GetClaimRepository(),
		Payment:  factory.GetPaymentRepository(),
		Document: factory.GetDocumentRepository(),
	}
}

// GetFactory returns the underlying repository factory
func (r *RepositoryRegistry) GetFactory() *RepositoryFactory {
	return r.factory
}

// WithTransaction executes a function within a database transaction
func (r *RepositoryRegistry) WithTransaction(fn func(*RepositoryRegistry) error) error {
	return r.factory.WithTransaction(func(txFactory *RepositoryFactory) error {
		txRegistry := &RepositoryRegistry{
			factory:  txFactory,
			User:     txFactory.GetUserRepository(),
			Device:   txFactory.GetDeviceRepository(),
			Policy:   txFactory.GetPolicyRepository(),
			Claim:    txFactory.GetClaimRepository(),
			Payment:  txFactory.GetPaymentRepository(),
			Document: txFactory.GetDocumentRepository(),
		}
		return fn(txRegistry)
	})
}

// === Repository Manager ===

// RepositoryManager provides a high-level interface for repository operations
type RepositoryManager struct {
	registry *RepositoryRegistry
	logger   Logger
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB, logger Logger) *RepositoryManager {
	return &RepositoryManager{
		registry: NewRepositoryRegistry(db, logger),
		logger:   logger,
	}
}

// GetRegistry returns the repository registry
func (m *RepositoryManager) GetRegistry() *RepositoryRegistry {
	return m.registry
}

// Users returns the user repository
func (m *RepositoryManager) Users() repositories.UserRepository {
	return m.registry.User
}

// Devices returns the device repository
func (m *RepositoryManager) Devices() repositories.DeviceRepository {
	return m.registry.Device
}

// Policies returns the policy repository
func (m *RepositoryManager) Policies() repositories.PolicyRepository {
	return m.registry.Policy
}

// Claims returns the claim repository
func (m *RepositoryManager) Claims() repositories.ClaimRepository {
	return m.registry.Claim
}

// Payments returns the payment repository
func (m *RepositoryManager) Payments() repositories.PaymentRepository {
	return m.registry.Payment
}

// Documents returns the document repository
func (m *RepositoryManager) Documents() repositories.DocumentRepository {
	return m.registry.Document
}

// RunInTransaction executes a function within a database transaction
func (m *RepositoryManager) RunInTransaction(fn func(*RepositoryManager) error) error {
	return m.registry.WithTransaction(func(txRegistry *RepositoryRegistry) error {
		txManager := &RepositoryManager{
			registry: txRegistry,
			logger:   m.logger,
		}
		return fn(txManager)
	})
}

// === Helper Functions ===

// CreateAllIndexes creates database indexes for all entities
func CreateAllIndexes(db *gorm.DB) error {
	// User indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone_number);
		CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
		CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
	`).Error; err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	// Device indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_devices_imei ON devices(imei);
		CREATE INDEX IF NOT EXISTS idx_devices_serial ON devices(serial_number);
		CREATE INDEX IF NOT EXISTS idx_devices_owner ON devices(owner_id);
		CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
	`).Error; err != nil {
		return fmt.Errorf("failed to create device indexes: %w", err)
	}

	// Policy indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_policies_number ON policies(policy_number);
		CREATE INDEX IF NOT EXISTS idx_policies_customer ON policies(customer_id);
		CREATE INDEX IF NOT EXISTS idx_policies_device ON policies(device_id);
		CREATE INDEX IF NOT EXISTS idx_policies_status ON policies(status);
		CREATE INDEX IF NOT EXISTS idx_policies_dates ON policies(start_date, end_date);
	`).Error; err != nil {
		return fmt.Errorf("failed to create policy indexes: %w", err)
	}

	// Claim indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_claims_number ON claims(claim_number);
		CREATE INDEX IF NOT EXISTS idx_claims_policy ON claims(policy_id);
		CREATE INDEX IF NOT EXISTS idx_claims_customer ON claims(customer_id);
		CREATE INDEX IF NOT EXISTS idx_claims_status ON claims(status);
		CREATE INDEX IF NOT EXISTS idx_claims_filed ON claims(filed_date);
	`).Error; err != nil {
		return fmt.Errorf("failed to create claim indexes: %w", err)
	}

	// Payment indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_payments_transaction ON payments(transaction_id);
		CREATE INDEX IF NOT EXISTS idx_payments_user ON payments(user_id);
		CREATE INDEX IF NOT EXISTS idx_payments_policy ON payments(policy_id);
		CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
		CREATE INDEX IF NOT EXISTS idx_payments_created ON payments(created_at);
	`).Error; err != nil {
		return fmt.Errorf("failed to create payment indexes: %w", err)
	}

	// Document indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_documents_number ON documents(id_document_number);
		CREATE INDEX IF NOT EXISTS idx_documents_type ON documents(id_type);
		CREATE INDEX IF NOT EXISTS idx_documents_category ON documents(id_category);
		CREATE INDEX IF NOT EXISTS idx_documents_status ON documents(lc_status);
		CREATE INDEX IF NOT EXISTS idx_documents_user ON documents(rel_user_id);
		CREATE INDEX IF NOT EXISTS idx_documents_policy ON documents(rel_policy_id);
		CREATE INDEX IF NOT EXISTS idx_documents_claim ON documents(rel_claim_id);
		CREATE INDEX IF NOT EXISTS idx_documents_created ON documents(created_at);
		CREATE INDEX IF NOT EXISTS idx_documents_expires ON documents(lc_expires_at);
		CREATE INDEX IF NOT EXISTS idx_documents_security ON documents(sec_security_level);
	`).Error; err != nil {
		return fmt.Errorf("failed to create document indexes: %w", err)
	}

	return nil
}

// MigrateAllModels runs auto-migration for all models
func MigrateAllModels(db *gorm.DB) error {
	// Note: In production, use proper migration tools like golang-migrate
	// This is just for development/testing purposes

	models := []interface{}{
		// Add all your model structs here
		// &models.User{},
		// &models.Device{},
		// &models.Policy{},
		// &models.Claim{},
		// &models.Payment{},
		// &models.PaymentMethod{},
		// &models.Invoice{},
		// &models.Subscription{},
		// ... add all other models
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	return nil
}
