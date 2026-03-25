package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/pkg/logger"
)

// CorporateService handles corporate and fleet management
type CorporateService struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewCorporateService creates a new CorporateService instance
func NewCorporateService(db *gorm.DB, logger *logger.Logger) *CorporateService {
	return &CorporateService{
		db:     db,
		logger: logger,
	}
}

// GetDB returns the database connection
func (s *CorporateService) GetDB() *gorm.DB {
	return s.db
}

// CorporateAccountRequest represents a request to create a corporate account
type CorporateAccountRequest struct {
	CompanyName    string  `json:"company_name" validate:"required"`
	TaxID          string  `json:"tax_id" validate:"required"`
	Industry       string  `json:"industry" validate:"required"`
	CompanySize    string  `json:"company_size" validate:"required"`
	ContactName    string  `json:"contact_name" validate:"required"`
	ContactEmail   string  `json:"contact_email" validate:"required,email"`
	ContactPhone   string  `json:"contact_phone" validate:"required"`
	BillingAddress string  `json:"billing_address" validate:"required"`
	PaymentTerms   string  `json:"payment_terms"`
	CreditLimit    float64 `json:"credit_limit"`
}

// EmployeeRegistrationRequest represents a request to register an employee
type EmployeeRegistrationRequest struct {
	CorporateAccountID uuid.UUID  `json:"corporate_account_id" validate:"required"`
	EmployeeID         string     `json:"employee_id" validate:"required"`
	FirstName          string     `json:"first_name" validate:"required"`
	LastName           string     `json:"last_name" validate:"required"`
	Email              string     `json:"email" validate:"required,email"`
	Phone              string     `json:"phone" validate:"required"`
	Department         string     `json:"department"`
	JobTitle           string     `json:"job_title"`
	ManagerID          *uuid.UUID `json:"manager_id"`
	StartDate          time.Time  `json:"start_date"`
}

// CorporatePolicyRequest represents a request to create a corporate policy
type CorporatePolicyRequest struct {
	CorporateAccountID uuid.UUID `json:"corporate_account_id" validate:"required"`
	PolicyName         string    `json:"policy_name" validate:"required"`
	CoverageType       string    `json:"coverage_type" validate:"required"`
	EmployeeLimit      int       `json:"employee_limit"`
	DeviceLimit        int       `json:"device_limit"`
	CoverageAmount     float64   `json:"coverage_amount" validate:"required"`
	Deductible         float64   `json:"deductible"`
	PremiumStructure   string    `json:"premium_structure"`
}

// FleetDeviceRequest represents a request to add a device to fleet
type FleetDeviceRequest struct {
	CorporateAccountID uuid.UUID `json:"corporate_account_id" validate:"required"`
	EmployeeID         uuid.UUID `json:"employee_id" validate:"required"`
	DeviceType         string    `json:"device_type" validate:"required"`
	Brand              string    `json:"brand" validate:"required"`
	Model              string    `json:"model" validate:"required"`
	SerialNumber       string    `json:"serial_number" validate:"required"`
	IMEI               string    `json:"imei"`
	PurchaseDate       time.Time `json:"purchase_date" validate:"required"`
	PurchasePrice      float64   `json:"purchase_price" validate:"required"`
	AssignmentType     string    `json:"assignment_type" validate:"required"`
}

// CreateCorporateAccount creates a new corporate account
func (s *CorporateService) CreateCorporateAccount(ctx context.Context, req *CorporateAccountRequest) (*models.CorporateAccount, error) {
	s.logger.Info("Creating corporate account", "company_name", req.CompanyName)

	// Check if account already exists
	var existingAccount models.CorporateAccount
	if err := s.db.WithContext(ctx).Where("tax_id = ? OR company_name = ?", req.TaxID, req.CompanyName).First(&existingAccount).Error; err == nil {
		return nil, errors.New("corporate account already exists")
	}

	// Generate account number
	accountNumber := s.generateAccountNumber()

	// Create corporate account
	account := &models.CorporateAccount{
		BusinessRegistration: accountNumber,
		CompanyName:          req.CompanyName,
		TaxID:                req.TaxID,
		Industry:             req.Industry,
		CompanySize:          req.CompanySize,
		Address:              req.BillingAddress,
		City:                 "Unknown",
		State:                "Unknown",
		Country:              "Unknown",
		Phone:                req.ContactPhone,
		Email:                req.ContactEmail,
		ContactPersonID:      uuid.New(),
		PaymentTerms:         req.PaymentTerms,
		CreditLimit:          req.CreditLimit,
		Status:               "pending_verification",
		ContractStartDate:    time.Now(),
	}

	// Set account settings in billing contact field (as JSON)
	settings := map[string]interface{}{
		"auto_enroll_employees":    true,
		"device_replacement_limit": 2,
		"claim_approval_required":  false,
		"bulk_billing_enabled":     true,
		"reporting_frequency":      "monthly",
	}
	settingsJSON, _ := json.Marshal(settings)
	account.BillingContact = string(settingsJSON)

	if err := s.db.WithContext(ctx).Create(account).Error; err != nil {
		return nil, fmt.Errorf("failed to create corporate account: %w", err)
	}

	s.logger.Info("Corporate account created successfully", "account_id", account.ID)
	return account, nil
}

// RegisterEmployee registers a new employee under a corporate account
func (s *CorporateService) RegisterEmployee(ctx context.Context, req *EmployeeRegistrationRequest) (*models.CorporateEmployee, error) {
	s.logger.Info("Registering employee", "employee_id", req.EmployeeID, "account_id", req.CorporateAccountID)

	// Verify corporate account exists and is active
	var account models.CorporateAccount
	if err := s.db.WithContext(ctx).First(&account, "id = ? AND is_active = ?", req.CorporateAccountID, true).Error; err != nil {
		return nil, errors.New("corporate account not found or inactive")
	}

	// Check if employee already exists
	var existingEmployee models.CorporateEmployee
	if err := s.db.WithContext(ctx).Where("corporate_account_id = ? AND (employee_id = ? OR email = ?)",
		req.CorporateAccountID, req.EmployeeID, req.Email).First(&existingEmployee).Error; err == nil {
		return nil, errors.New("employee already exists")
	}

	// Create employee record
	employee := &models.CorporateEmployee{
		CorporateAccountID: req.CorporateAccountID,
		UserID:             uuid.New(), // Create new user ID
		EmployeeID:         req.EmployeeID,
		Department:         req.Department,
		JobTitle:           req.JobTitle,
		Manager:            "TBD", // Convert UUID to string representation later
		StartDate:          req.StartDate,
		Status:             "active",
		CoverageLevel:      "standard",
	}

	// Set employee permissions (stored separately if needed)
	// Note: Permissions field doesn't exist in CorporateEmployee model

	if err := s.db.WithContext(ctx).Create(employee).Error; err != nil {
		return nil, fmt.Errorf("failed to register employee: %w", err)
	}

	s.logger.Info("Employee registered successfully", "employee_id", employee.ID)
	return employee, nil
}

// CreateCorporatePolicy creates a corporate insurance policy
func (s *CorporateService) CreateCorporatePolicy(ctx context.Context, req *CorporatePolicyRequest) (*models.CorporatePolicy, error) {
	s.logger.Info("Creating corporate policy", "account_id", req.CorporateAccountID, "policy_name", req.PolicyName)

	// Verify corporate account
	var account models.CorporateAccount
	if err := s.db.WithContext(ctx).First(&account, "id = ? AND is_active = ?", req.CorporateAccountID, true).Error; err != nil {
		return nil, errors.New("corporate account not found or inactive")
	}

	// Generate policy number
	policyNumber := s.generatePolicyNumber()

	// Calculate premium based on coverage and limits
	premium := s.calculateCorporatePremium(req)

	// Create corporate policy
	policy := &models.CorporatePolicy{
		CorporateAccountID: req.CorporateAccountID,
		PolicyNumber:       policyNumber,
		PolicyType:         "corporate",
		CoverageType:       req.CoverageType,
		MaxDevicesPerUser:  2,
		TotalDeviceLimit:   req.DeviceLimit,
		CoverageAmount:     req.CoverageAmount,
		DeductibleAmount:   req.Deductible,
		PremiumPerDevice:   premium,
		Status:             "active",
		StartDate:          time.Now(),
		EndDate:            time.Now().AddDate(1, 0, 0), // 1 year
	}

	// Set policy terms
	terms := map[string]interface{}{
		"coverage_type":        req.CoverageType,
		"employee_limit":       req.EmployeeLimit,
		"device_limit":         req.DeviceLimit,
		"coverage_amount":      req.CoverageAmount,
		"deductible":           req.Deductible,
		"replacement_limit":    2,
		"claim_limit_per_year": 3,
	}
	termsJSON, _ := json.Marshal(terms)
	policy.Terms = string(termsJSON)

	if err := s.db.WithContext(ctx).Create(policy).Error; err != nil {
		return nil, fmt.Errorf("failed to create corporate policy: %w", err)
	}

	s.logger.Info("Corporate policy created successfully", "policy_id", policy.ID)
	return policy, nil
}

// AddFleetDevice adds a device to the corporate fleet
func (s *CorporateService) AddFleetDevice(ctx context.Context, req *FleetDeviceRequest) (*models.FleetDevice, error) {
	s.logger.Info("Adding fleet device", "account_id", req.CorporateAccountID, "serial_number", req.SerialNumber)

	// Verify corporate account and employee
	var account models.CorporateAccount
	if err := s.db.WithContext(ctx).First(&account, "id = ? AND is_active = ?", req.CorporateAccountID, true).Error; err != nil {
		return nil, errors.New("corporate account not found or inactive")
	}

	var employee models.CorporateEmployee
	if err := s.db.WithContext(ctx).First(&employee, "id = ? AND corporate_account_id = ?", req.EmployeeID, req.CorporateAccountID).Error; err != nil {
		return nil, errors.New("employee not found")
	}

	// Check if device already exists
	var existingDevice models.FleetDevice
	if err := s.db.WithContext(ctx).Where("serial_number = ? OR imei = ?", req.SerialNumber, req.IMEI).First(&existingDevice).Error; err == nil {
		return nil, errors.New("device already exists in fleet")
	}

	// Generate asset tag
	assetTag := s.generateAssetTag(req.CorporateAccountID)

	// Create fleet device
	device := &models.FleetDevice{
		CorporateAccountID: req.CorporateAccountID,
		DeviceID:           uuid.New(), // Create new device ID
		AssignedEmployeeID: &req.EmployeeID,
		AssetTag:           assetTag,
		PurchaseOrder:      "PO-" + assetTag,
		Department:         "IT",
		DeviceType:         req.DeviceType,
		AssignmentDate:     &req.PurchaseDate,
		Status:             "available",
		Condition:          "good",
	}

	// Set device configuration
	config := map[string]interface{}{
		"mdm_enrolled":      true,
		"security_policies": []string{"encryption", "remote_wipe", "app_restrictions"},
		"backup_enabled":    true,
		"location_tracking": true,
		"usage_monitoring":  true,
	}
	configJSON, _ := json.Marshal(config)
	device.SecurityPolicies = string(configJSON)

	if err := s.db.WithContext(ctx).Create(device).Error; err != nil {
		return nil, fmt.Errorf("failed to add fleet device: %w", err)
	}

	s.logger.Info("Fleet device added successfully", "device_id", device.ID)
	return device, nil
}

// GetCorporateAnalytics generates analytics for a corporate account
func (s *CorporateService) GetCorporateAnalytics(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	s.logger.Info("Generating corporate analytics", "account_id", accountID)

	analytics := make(map[string]interface{})

	// Employee metrics
	var totalEmployees, activeEmployees int64
	s.db.WithContext(ctx).Model(&models.CorporateEmployee{}).Where("corporate_account_id = ?", accountID).Count(&totalEmployees)
	s.db.WithContext(ctx).Model(&models.CorporateEmployee{}).Where("corporate_account_id = ? AND is_active = ?", accountID, true).Count(&activeEmployees)

	analytics["total_employees"] = totalEmployees
	analytics["active_employees"] = activeEmployees

	// Device metrics
	var totalDevices, activeDevices int64
	s.db.WithContext(ctx).Model(&models.FleetDevice{}).Where("corporate_account_id = ?", accountID).Count(&totalDevices)
	s.db.WithContext(ctx).Model(&models.FleetDevice{}).Where("corporate_account_id = ? AND is_active = ?", accountID, true).Count(&activeDevices)

	analytics["total_devices"] = totalDevices
	analytics["active_devices"] = activeDevices

	// Policy metrics
	var activePolicies int64
	var totalPremiums float64
	s.db.WithContext(ctx).Model(&models.CorporatePolicy{}).Where("corporate_account_id = ? AND is_active = ?", accountID, true).Count(&activePolicies)
	s.db.WithContext(ctx).Model(&models.CorporatePolicy{}).Where("corporate_account_id = ?", accountID).Select("SUM(premium_amount)").Scan(&totalPremiums)

	analytics["active_policies"] = activePolicies
	analytics["total_premiums"] = totalPremiums

	// Claims metrics for the period
	var totalClaims, approvedClaims int64
	var totalClaimAmount float64

	// Note: This assumes claims are linked to corporate employees through user_id
	var employeeIDs []uuid.UUID
	s.db.WithContext(ctx).Model(&models.CorporateEmployee{}).Where("corporate_account_id = ?", accountID).Pluck("user_id", &employeeIDs)

	if len(employeeIDs) > 0 {
		s.db.WithContext(ctx).Model(&models.Claim{}).
			Where("user_id IN ? AND created_at BETWEEN ? AND ?", employeeIDs, startDate, endDate).
			Count(&totalClaims)
		s.db.WithContext(ctx).Model(&models.Claim{}).
			Where("user_id IN ? AND created_at BETWEEN ? AND ? AND status = ?", employeeIDs, startDate, endDate, "approved").
			Count(&approvedClaims)
		s.db.WithContext(ctx).Model(&models.Claim{}).
			Where("user_id IN ? AND created_at BETWEEN ? AND ? AND status = ?", employeeIDs, startDate, endDate, "approved").
			Select("SUM(claim_amount)").Scan(&totalClaimAmount)
	}

	analytics["total_claims"] = totalClaims
	analytics["approved_claims"] = approvedClaims
	analytics["total_claim_amount"] = totalClaimAmount

	// Calculate ratios
	if totalClaims > 0 {
		analytics["claim_approval_rate"] = float64(approvedClaims) / float64(totalClaims)
	} else {
		analytics["claim_approval_rate"] = 0.0
	}

	if totalPremiums > 0 {
		analytics["loss_ratio"] = totalClaimAmount / totalPremiums
	} else {
		analytics["loss_ratio"] = 0.0
	}

	analytics["devices_per_employee"] = float64(totalDevices) / float64(totalEmployees)

	return analytics, nil
}

// Helper methods

func (s *CorporateService) generateAccountNumber() string {
	return fmt.Sprintf("CORP-%d", time.Now().Unix())
}

func (s *CorporateService) generatePolicyNumber() string {
	return fmt.Sprintf("POL-%d-%s", time.Now().Unix(), uuid.New().String()[:8])
}

func (s *CorporateService) generateAssetTag(accountID uuid.UUID) string {
	return fmt.Sprintf("AST-%s-%d", accountID.String()[:8], time.Now().Unix()%10000)
}

func (s *CorporateService) calculateCorporatePremium(req *CorporatePolicyRequest) float64 {
	basePremium := 50.0 // Base premium per employee per month

	// Coverage type multiplier
	coverageMultipliers := map[string]float64{
		"basic":         1.0,
		"standard":      1.5,
		"comprehensive": 2.0,
		"premium":       2.5,
	}

	multiplier := coverageMultipliers[req.CoverageType]
	if multiplier == 0 {
		multiplier = 1.0
	}

	// Employee count factor
	employeeDiscount := 1.0
	if req.EmployeeLimit > 100 {
		employeeDiscount = 0.85 // 15% discount for large companies
	} else if req.EmployeeLimit > 50 {
		employeeDiscount = 0.90 // 10% discount for medium companies
	} else if req.EmployeeLimit > 20 {
		employeeDiscount = 0.95 // 5% discount for small companies
	}

	// Coverage amount factor
	coverageFactor := req.CoverageAmount / 1000.0
	if coverageFactor > 5.0 {
		coverageFactor = 5.0 // Cap the factor
	}

	totalPremium := basePremium * float64(req.EmployeeLimit) * multiplier * employeeDiscount * coverageFactor

	return totalPremium
}
