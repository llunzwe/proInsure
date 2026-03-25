package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
)

// CorporateService defines the interface for corporate account business logic operations
type CorporateService interface {
	// Corporate account management
	CreateCorporateAccount(ctx context.Context, account *models.CorporateAccount) error
	GetCorporateAccountByID(ctx context.Context, id uuid.UUID) (*models.CorporateAccount, error)
	GetCorporateAccountByRegistration(ctx context.Context, registration string) (*models.CorporateAccount, error)
	UpdateCorporateAccount(ctx context.Context, account *models.CorporateAccount) error
	DeactivateCorporateAccount(ctx context.Context, id uuid.UUID, reason string) error
	ListCorporateAccounts(ctx context.Context, status string, limit, offset int) ([]*models.CorporateAccount, int64, error)

	// Corporate employee management
	AddEmployee(ctx context.Context, accountID uuid.UUID, employee *models.CorporateEmployee) error
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (*models.CorporateEmployee, error)
	GetEmployeesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.CorporateEmployee, int64, error)
	UpdateEmployee(ctx context.Context, employee *models.CorporateEmployee) error
	RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error
	GetActiveEmployeesByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.CorporateEmployee, error)

	// Corporate policy management
	CreateCorporatePolicy(ctx context.Context, policy *models.CorporatePolicy) error
	GetCorporatePolicyByID(ctx context.Context, id uuid.UUID) (*models.CorporatePolicy, error)
	GetCorporatePoliciesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.CorporatePolicy, int64, error)
	UpdateCorporatePolicy(ctx context.Context, policy *models.CorporatePolicy) error
	GetActiveCorporatePoliciesByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.CorporatePolicy, error)

	// Fleet device management
	AddFleetDevice(ctx context.Context, device *models.FleetDevice) error
	GetFleetDeviceByID(ctx context.Context, id uuid.UUID) (*models.FleetDevice, error)
	GetFleetDevicesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error)
	AssignFleetDevice(ctx context.Context, deviceID uuid.UUID, employeeID uuid.UUID) error
	ReturnFleetDevice(ctx context.Context, deviceID uuid.UUID) error
	GetAvailableFleetDevices(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error)
	UpdateFleetDevice(ctx context.Context, device *models.FleetDevice) error

	// Corporate billing and financials
	CalculateCorporateDiscount(ctx context.Context, accountID uuid.UUID, baseAmount float64) (float64, error)
	GetCorporateBillingSummary(ctx context.Context, accountID uuid.UUID, startDate, endDate string) (map[string]interface{}, error)

	// Corporate statistics
	GetCorporateAccountStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error)
	GetEmployeeStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error)
	GetFleetStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error)
}
