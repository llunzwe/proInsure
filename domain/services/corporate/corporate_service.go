package corporate

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// corporateService implements the CorporateService interface
type corporateService struct {
	corporateRepo repositories.CorporateRepository
}

// NewCorporateService creates a new corporate service
func NewCorporateService(
	corporateRepo repositories.CorporateRepository,
) services.CorporateService {
	return &corporateService{
		corporateRepo: corporateRepo,
	}
}

// CreateCorporateAccount creates a new corporate account
func (s *corporateService) CreateCorporateAccount(ctx context.Context, account *models.CorporateAccount) error {
	if account == nil {
		return errors.New("corporate account cannot be nil")
	}
	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}
	if account.Status == "" {
		account.Status = "active"
	}
	return s.corporateRepo.CreateCorporateAccount(ctx, account)
}

// GetCorporateAccountByID retrieves a corporate account by ID
func (s *corporateService) GetCorporateAccountByID(ctx context.Context, id uuid.UUID) (*models.CorporateAccount, error) {
	return s.corporateRepo.GetCorporateAccountByID(ctx, id)
}

// GetCorporateAccountByRegistration retrieves a corporate account by registration number
func (s *corporateService) GetCorporateAccountByRegistration(ctx context.Context, registration string) (*models.CorporateAccount, error) {
	return s.corporateRepo.GetCorporateAccountByRegistration(ctx, registration)
}

// UpdateCorporateAccount updates an existing corporate account
func (s *corporateService) UpdateCorporateAccount(ctx context.Context, account *models.CorporateAccount) error {
	if account == nil {
		return errors.New("corporate account cannot be nil")
	}
	return s.corporateRepo.UpdateCorporateAccount(ctx, account)
}

// DeactivateCorporateAccount deactivates a corporate account
func (s *corporateService) DeactivateCorporateAccount(ctx context.Context, id uuid.UUID, reason string) error {
	account, err := s.corporateRepo.GetCorporateAccountByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get corporate account: %w", err)
	}
	if account == nil {
		return errors.New("corporate account not found")
	}
	account.Status = "suspended"
	return s.corporateRepo.UpdateCorporateAccount(ctx, account)
}

// ListCorporateAccounts lists corporate accounts with filters
func (s *corporateService) ListCorporateAccounts(ctx context.Context, status string, limit, offset int) ([]*models.CorporateAccount, int64, error) {
	return s.corporateRepo.ListCorporateAccounts(ctx, status, limit, offset)
}

// AddEmployee adds an employee to a corporate account
func (s *corporateService) AddEmployee(ctx context.Context, accountID uuid.UUID, employee *models.CorporateEmployee) error {
	if employee == nil {
		return errors.New("employee cannot be nil")
	}
	if employee.ID == uuid.Nil {
		employee.ID = uuid.New()
	}
	employee.CorporateAccountID = accountID
	return s.corporateRepo.CreateCorporateEmployee(ctx, employee)
}

// GetEmployeeByID retrieves an employee by ID
func (s *corporateService) GetEmployeeByID(ctx context.Context, id uuid.UUID) (*models.CorporateEmployee, error) {
	return s.corporateRepo.GetCorporateEmployeeByID(ctx, id)
}

// GetEmployeesByAccountID gets employees for a corporate account
func (s *corporateService) GetEmployeesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.CorporateEmployee, int64, error) {
	return s.corporateRepo.GetCorporateEmployeesByAccountID(ctx, accountID, limit, offset)
}

// UpdateEmployee updates an employee
func (s *corporateService) UpdateEmployee(ctx context.Context, employee *models.CorporateEmployee) error {
	if employee == nil {
		return errors.New("employee cannot be nil")
	}
	return s.corporateRepo.UpdateCorporateEmployee(ctx, employee)
}

// RemoveEmployee removes an employee from a corporate account
func (s *corporateService) RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error {
	return s.corporateRepo.DeleteCorporateEmployee(ctx, employeeID)
}

// GetActiveEmployeesByAccountID gets active employees for a corporate account
func (s *corporateService) GetActiveEmployeesByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.CorporateEmployee, error) {
	return s.corporateRepo.GetActiveEmployeesByAccountID(ctx, accountID)
}

// CreateCorporatePolicy creates a corporate policy
func (s *corporateService) CreateCorporatePolicy(ctx context.Context, policy *models.CorporatePolicy) error {
	if policy == nil {
		return errors.New("corporate policy cannot be nil")
	}
	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}
	return s.corporateRepo.CreateCorporatePolicy(ctx, policy)
}

// GetCorporatePolicyByID retrieves a corporate policy by ID
func (s *corporateService) GetCorporatePolicyByID(ctx context.Context, id uuid.UUID) (*models.CorporatePolicy, error) {
	return s.corporateRepo.GetCorporatePolicyByID(ctx, id)
}

// GetCorporatePoliciesByAccountID gets corporate policies for an account
func (s *corporateService) GetCorporatePoliciesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.CorporatePolicy, int64, error) {
	return s.corporateRepo.GetCorporatePoliciesByAccountID(ctx, accountID, limit, offset)
}

// UpdateCorporatePolicy updates a corporate policy
func (s *corporateService) UpdateCorporatePolicy(ctx context.Context, policy *models.CorporatePolicy) error {
	if policy == nil {
		return errors.New("corporate policy cannot be nil")
	}
	return s.corporateRepo.UpdateCorporatePolicy(ctx, policy)
}

// GetActiveCorporatePoliciesByAccountID gets active corporate policies for an account
func (s *corporateService) GetActiveCorporatePoliciesByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.CorporatePolicy, error) {
	return s.corporateRepo.GetActiveCorporatePoliciesByAccountID(ctx, accountID)
}

// AddFleetDevice adds a device to the corporate fleet
func (s *corporateService) AddFleetDevice(ctx context.Context, device *models.FleetDevice) error {
	if device == nil {
		return errors.New("fleet device cannot be nil")
	}
	if device.ID == uuid.Nil {
		device.ID = uuid.New()
	}
	return s.corporateRepo.CreateFleetDevice(ctx, device)
}

// GetFleetDeviceByID retrieves a fleet device by ID
func (s *corporateService) GetFleetDeviceByID(ctx context.Context, id uuid.UUID) (*models.FleetDevice, error) {
	return s.corporateRepo.GetFleetDeviceByID(ctx, id)
}

// GetFleetDevicesByAccountID gets fleet devices for a corporate account
func (s *corporateService) GetFleetDevicesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error) {
	return s.corporateRepo.GetFleetDevicesByAccountID(ctx, accountID, limit, offset)
}

// AssignFleetDevice assigns a fleet device to an employee
func (s *corporateService) AssignFleetDevice(ctx context.Context, deviceID uuid.UUID, employeeID uuid.UUID) error {
	device, err := s.corporateRepo.GetFleetDeviceByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get fleet device: %w", err)
	}
	if device == nil {
		return errors.New("fleet device not found")
	}
	device.AssignedEmployeeID = &employeeID
	device.Status = "assigned"
	now := time.Now()
	device.AssignmentDate = &now
	return s.corporateRepo.UpdateFleetDevice(ctx, device)
}

// ReturnFleetDevice returns a fleet device from an employee
func (s *corporateService) ReturnFleetDevice(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.corporateRepo.GetFleetDeviceByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get fleet device: %w", err)
	}
	if device == nil {
		return errors.New("fleet device not found")
	}
	device.AssignedEmployeeID = nil
	device.Status = "available"
	now := time.Now()
	device.ReturnDate = &now
	return s.corporateRepo.UpdateFleetDevice(ctx, device)
}

// GetAvailableFleetDevices gets available fleet devices
func (s *corporateService) GetAvailableFleetDevices(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*models.FleetDevice, int64, error) {
	return s.corporateRepo.GetAvailableFleetDevices(ctx, accountID, limit, offset)
}

// UpdateFleetDevice updates a fleet device
func (s *corporateService) UpdateFleetDevice(ctx context.Context, device *models.FleetDevice) error {
	if device == nil {
		return errors.New("fleet device cannot be nil")
	}
	return s.corporateRepo.UpdateFleetDevice(ctx, device)
}

// CalculateCorporateDiscount calculates discount for corporate account
func (s *corporateService) CalculateCorporateDiscount(ctx context.Context, accountID uuid.UUID, baseAmount float64) (float64, error) {
	account, err := s.corporateRepo.GetCorporateAccountByID(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get corporate account: %w", err)
	}
	if account == nil {
		return 0, errors.New("corporate account not found")
	}
	return account.CalculateDiscount(baseAmount), nil
}

// GetCorporateBillingSummary gets billing summary for corporate account
func (s *corporateService) GetCorporateBillingSummary(ctx context.Context, accountID uuid.UUID, startDate, endDate string) (map[string]interface{}, error) {
	// This would aggregate billing data from policies and devices
	// Simplified implementation
	return map[string]interface{}{
		"account_id": accountID,
		"period":     fmt.Sprintf("%s to %s", startDate, endDate),
	}, nil
}

// GetCorporateAccountStatistics gets statistics for a corporate account
func (s *corporateService) GetCorporateAccountStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error) {
	return s.corporateRepo.GetCorporateAccountStatistics(ctx, accountID)
}

// GetEmployeeStatistics gets employee statistics for a corporate account
func (s *corporateService) GetEmployeeStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error) {
	return s.corporateRepo.GetEmployeeStatistics(ctx, accountID)
}

// GetFleetStatistics gets fleet statistics for a corporate account
func (s *corporateService) GetFleetStatistics(ctx context.Context, accountID uuid.UUID) (map[string]interface{}, error) {
	return s.corporateRepo.GetFleetStatistics(ctx, accountID)
}
