package repositories

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// CorporateRepository defines the interface for corporate account persistence operations
type CorporateRepository interface {
	// Corporate account operations
	CreateCorporateAccount(ctx context.Context, account *models.CorporateAccount) error
	GetCorporateAccountByID(ctx context.Context, id uuid.UUID) (*models.CorporateAccount, error)
	GetCorporateAccountByRegistration(ctx context.Context, registration string) (*models.CorporateAccount, error)
	GetCorporateAccountByTaxID(ctx context.Context, taxID string) (*models.CorporateAccount, error)
	UpdateCorporateAccount(ctx context.Context, account *models.CorporateAccount) error
	DeleteCorporateAccount(ctx context.Context, id uuid.UUID) error
	ListCorporateAccounts(ctx context.Context, status string, limit, offset int) ([]*models.CorporateAccount, int64, error)
	SearchCorporateAccounts(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*models.CorporateAccount, int64, error)

	// Corporate employee operations
	CreateCorporateEmployee(ctx context.Context, employee *models.CorporateEmployee) error
	GetCorporateEmployeeByID(ctx context.Context, id uuid.UUID) (*models.CorporateEmployee, error)
	GetCorporateEmployeesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.CorporateEmployee, int64, error)
	GetCorporateEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*models.CorporateEmployee, error)
	GetCorporateEmployeeByEmployeeID(ctx context.Context, accountID uuid.UUID, employeeID string) (*models.CorporateEmployee, error)
	UpdateCorporateEmployee(ctx context.Context, employee *models.CorporateEmployee) error
	DeleteCorporateEmployee(ctx context.Context, id uuid.UUID) error
	GetActiveEmployeesByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.CorporateEmployee, error)

	// Corporate policy operations
	CreateCorporatePolicy(ctx context.Context, policy *models.CorporatePolicy) error
	GetCorporatePolicyByID(ctx context.Context, id uuid.UUID) (*models.CorporatePolicy, error)
	GetCorporatePolicyByNumber(ctx context.Context, policyNumber string) (*models.CorporatePolicy, error)
	GetCorporatePoliciesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.CorporatePolicy, int64, error)
	GetActiveCorporatePoliciesByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.CorporatePolicy, error)
	UpdateCorporatePolicy(ctx context.Context, policy *models.CorporatePolicy) error
	DeleteCorporatePolicy(ctx context.Context, id uuid.UUID) error

	// Fleet device operations
	CreateFleetDevice(ctx context.Context, device *models.FleetDevice) error
	GetFleetDeviceByID(ctx context.Context, id uuid.UUID) (*models.FleetDevice, error)
	GetFleetDevicesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error)
	GetFleetDevicesByPolicyID(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error)
	GetFleetDevicesByEmployeeID(ctx context.Context, employeeID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error)
	GetFleetDevicesByStatus(ctx context.Context, accountID uuid.UUID, status string, limit, offset int) ([]*models.FleetDevice, int64, error)
	GetAvailableFleetDevices(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error)
	UpdateFleetDevice(ctx context.Context, device *models.FleetDevice) error
	DeleteFleetDevice(ctx context.Context, id uuid.UUID) error

	// Statistics
	GetCorporateAccountStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error)
	GetEmployeeStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error)
	GetFleetStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error)
}
