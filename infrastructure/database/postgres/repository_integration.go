package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string

	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // in seconds
	ConnMaxIdleTime int // in seconds

	// GORM settings
	LogLevel      logger.LogLevel
	SkipDefaultTx bool
	PrepareStmt   bool
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "smartsure",
		SSLMode:  "disable",

		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 300, // 5 minutes
		ConnMaxIdleTime: 60,  // 1 minute

		LogLevel:      logger.Info,
		SkipDefaultTx: true,
		PrepareStmt:   true,
	}
}

// BuildDSN builds PostgreSQL DSN from config
func (c *DatabaseConfig) BuildDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// NewDatabase creates a new database connection
func NewDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := config.BuildDSN()

	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(config.LogLevel),
		SkipDefaultTransaction: config.SkipDefaultTx,
		PrepareStmt:            config.PrepareStmt,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RepositoryService provides a complete repository service with all dependencies
type RepositoryService struct {
	db      *gorm.DB
	logger  Logger
	manager *RepositoryManager
}

// NewRepositoryService creates a new repository service
func NewRepositoryService(config *DatabaseConfig, logger Logger) (*RepositoryService, error) {
	// Create database connection
	db, err := NewDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create repository manager
	manager := NewRepositoryManager(db, logger)

	service := &RepositoryService{
		db:      db,
		logger:  logger,
		manager: manager,
	}

	return service, nil
}

// GetManager returns the repository manager
func (s *RepositoryService) GetManager() *RepositoryManager {
	return s.manager
}

// GetDB returns the database connection
func (s *RepositoryService) GetDB() *gorm.DB {
	return s.db
}

// Initialize initializes the database with indexes and migrations
func (s *RepositoryService) Initialize() error {
	// Create indexes
	if err := CreateAllIndexes(s.db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	// Note: In production, use proper migration tools
	// if err := MigrateAllModels(s.db); err != nil {
	//     return fmt.Errorf("failed to migrate models: %w", err)
	// }

	s.logger.Info("Database initialized successfully")
	return nil
}

// Close closes the database connection
func (s *RepositoryService) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	s.logger.Info("Database connection closed")
	return nil
}

// HealthCheck performs a health check on the database
func (s *RepositoryService) HealthCheck(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Check if we can execute a simple query
	var result int
	if err := s.db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	return nil
}

// === Usage Example ===

/*
Example usage of the repository service:

func main() {
    // Create logger (implement your logger that satisfies the Logger interface)
    logger := NewAppLogger()

    // Create database config
    config := &DatabaseConfig{
        Host:     os.Getenv("DB_HOST"),
        Port:     5432,
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        SSLMode:  "require",
    }

    // Create repository service
    service, err := NewRepositoryService(config, logger)
    if err != nil {
        log.Fatal("Failed to create repository service:", err)
    }
    defer service.Close()

    // Initialize database
    if err := service.Initialize(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }

    // Get repository manager
    repos := service.GetManager()

    // Use repositories
    ctx := context.Background()

    // Example: Create a user
    user := &models.User{
        Email:     "user@example.com",
        FirstName: "John",
        LastName:  "Doe",
    }
    if err := repos.Users().Create(ctx, user); err != nil {
        log.Error("Failed to create user:", err)
    }

    // Example: Transaction
    err = repos.RunInTransaction(func(txRepos *RepositoryManager) error {
        // Create device
        device := &models.Device{
            IMEI:    "123456789",
            OwnerID: user.ID,
        }
        if err := txRepos.Devices().Create(ctx, device); err != nil {
            return err
        }

        // Create policy
        policy := &models.Policy{
            CustomerID: user.ID,
            DeviceID:   device.ID,
        }
        if err := txRepos.Policies().Create(ctx, policy); err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        log.Error("Transaction failed:", err)
    }
}
*/

// === Testing Helpers ===

// NewTestDatabase creates a test database connection
func NewTestDatabase() (*gorm.DB, error) {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "smartsure_test",
		SSLMode:  "disable",
		LogLevel: logger.Silent, // Silent for tests
	}

	return NewDatabase(config)
}

// SetupTestRepositories creates repositories for testing
func SetupTestRepositories(t testing.TB) (*RepositoryManager, func()) {
	db, err := NewTestDatabase()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations for test
	if err := MigrateAllModels(db); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	logger := NewTestLogger(t)
	manager := NewRepositoryManager(db, logger)

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return manager, cleanup
}

// TestLogger implements Logger interface for testing
type TestLogger struct {
	t testing.TB
}

func NewTestLogger(t testing.TB) *TestLogger {
	return &TestLogger{t: t}
}

func (l *TestLogger) Info(msg string, fields ...interface{}) {
	l.t.Logf("INFO: %s %v", msg, fields)
}

func (l *TestLogger) Error(msg string, err error, fields ...interface{}) {
	l.t.Logf("ERROR: %s - %v %v", msg, err, fields)
}

func (l *TestLogger) Debug(msg string, fields ...interface{}) {
	l.t.Logf("DEBUG: %s %v", msg, fields)
}

func (l *TestLogger) Warn(msg string, fields ...interface{}) {
	l.t.Logf("WARN: %s %v", msg, fields)
}
