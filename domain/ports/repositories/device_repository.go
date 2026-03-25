package repositories

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// DeviceSearchFilters defines filters for device search
type DeviceSearchFilters struct {
	Brand            string
	Model            string
	Status           string
	OwnerID          *uuid.UUID
	MinValue         float64
	MaxValue         float64
	MinRiskScore     float64
	MaxRiskScore     float64
	IsInsured        *bool
	IsStolen         *bool
	IsBlacklisted    *bool
	IsCorporateOwned *bool
	PurchasedAfter   *time.Time
	PurchasedBefore  *time.Time
	Limit            int
	Offset           int
	SortBy           string
	SortOrder        string
}

// DeviceRepository defines the comprehensive contract for device persistence
// This interface covers all aspects of device management including lifecycle,
// insurance, corporate management, analytics, and ecosystem integrations
type DeviceRepository interface {
	// ======================================
	// BASIC CRUD OPERATIONS
	// ======================================

	// Core CRUD operations
	Create(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error)
	GetByIDWithRelations(ctx context.Context, id uuid.UUID, relations []string) (*models.Device, error)
	GetByIMEI(ctx context.Context, imei string) (*models.Device, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (*models.Device, error)
	Update(ctx context.Context, device *models.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error

	// Batch operations
	CreateBatch(ctx context.Context, devices []*models.Device) error
	UpdateBatch(ctx context.Context, devices []*models.Device) error
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error

	// ======================================
	// DEVICE IDENTIFICATION & SEARCH
	// ======================================

	// Search by device attributes
	FindByBrandAndModel(ctx context.Context, brand, model string) ([]*models.Device, error)
	FindByManufacturer(ctx context.Context, manufacturer string) ([]*models.Device, error)
	FindByCategory(ctx context.Context, category string) ([]*models.Device, error)
	FindByOS(ctx context.Context, os string) ([]*models.Device, error)
	FindByModelYear(ctx context.Context, year int) ([]*models.Device, error)
	Search(ctx context.Context, filters DeviceSearchFilters) ([]*models.Device, int64, error)

	// ======================================
	// OWNERSHIP & USER MANAGEMENT
	// ======================================

	// Ownership queries
	FindByOwner(ctx context.Context, ownerID uuid.UUID) ([]*models.Device, error)
	FindByOwnerWithPagination(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]*models.Device, int64, error)
	TransferOwnership(ctx context.Context, deviceID, fromOwnerID, toOwnerID uuid.UUID) error
	GetOwnershipHistory(ctx context.Context, deviceID uuid.UUID) ([]*models.Device, error)

	// ======================================
	// LIFECYCLE & STATUS MANAGEMENT
	// ======================================

	// Status operations
	FindByStatus(ctx context.Context, status string) ([]*models.Device, error)
	FindByMultipleStatuses(ctx context.Context, statuses []string) ([]*models.Device, error)
	UpdateStatus(ctx context.Context, deviceID uuid.UUID, status string) error
	UpdateStatusBatch(ctx context.Context, deviceIDs []uuid.UUID, status string) error
	MarkAsStolen(ctx context.Context, deviceID uuid.UUID, reportDate time.Time) error
	MarkAsRecovered(ctx context.Context, deviceID uuid.UUID) error
	MarkAsLost(ctx context.Context, deviceID uuid.UUID) error

	// Verification operations
	MarkAsVerified(ctx context.Context, deviceID uuid.UUID, verifierID uuid.UUID) error
	FindUnverified(ctx context.Context) ([]*models.Device, error)
	FindPendingVerification(ctx context.Context) ([]*models.Device, error)

	// ======================================
	// CONDITION & QUALITY ASSESSMENT
	// ======================================

	// Condition queries
	FindByCondition(ctx context.Context, condition string) ([]*models.Device, error)
	FindByGrade(ctx context.Context, grade string) ([]*models.Device, error)
	FindByConditionRange(ctx context.Context, minCondition, maxCondition string) ([]*models.Device, error)
	UpdateCondition(ctx context.Context, deviceID uuid.UUID, condition, grade string) error
	GetConditionHistory(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// INSURANCE & POLICY MANAGEMENT
	// ======================================

	// Insurance eligibility
	FindEligibleForInsurance(ctx context.Context) ([]*models.Device, error)
	CheckInsuranceEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error)
	FindInsured(ctx context.Context) ([]*models.Device, error)
	FindUninsured(ctx context.Context) ([]*models.Device, error)
	FindByPolicyStatus(ctx context.Context, policyStatus string) ([]*models.Device, error)

	// Claims related
	FindWithActiveClaims(ctx context.Context) ([]*models.Device, error)
	FindWithPendingClaims(ctx context.Context) ([]*models.Device, error)
	GetClaimHistory(ctx context.Context, deviceID uuid.UUID) ([]*models.Claim, error)
	GetTotalClaimAmount(ctx context.Context, deviceID uuid.UUID) (float64, error)
	CountClaims(ctx context.Context, deviceID uuid.UUID) (int64, error)

	// ======================================
	// RISK ASSESSMENT & FRAUD DETECTION
	// ======================================

	// Risk assessment
	FindHighRiskDevices(ctx context.Context, threshold float64) ([]*models.Device, error)
	FindLowRiskDevices(ctx context.Context, threshold float64) ([]*models.Device, error)
	UpdateRiskScore(ctx context.Context, deviceID uuid.UUID, score float64) error
	CalculateRiskScore(ctx context.Context, deviceID uuid.UUID) (float64, error)
	GetRiskHistory(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// Fraud detection
	FindBlacklisted(ctx context.Context) ([]*models.Device, error)
	AddToBlacklist(ctx context.Context, deviceID uuid.UUID, reason string) error
	RemoveFromBlacklist(ctx context.Context, deviceID uuid.UUID) error
	CheckBlacklistStatus(ctx context.Context, imei string) (bool, string, error)
	FindSuspiciousDevices(ctx context.Context) ([]*models.Device, error)
	GetFraudInvestigations(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// CORPORATE & ENTERPRISE MANAGEMENT
	// ======================================

	// Corporate account management
	FindByCorporateAccount(ctx context.Context, accountID uuid.UUID) ([]*models.Device, error)
	FindByCorporateDepartment(ctx context.Context, departmentID uuid.UUID) ([]*models.Device, error)
	AssignToCorporateAccount(ctx context.Context, deviceID, accountID uuid.UUID) error
	UnassignFromCorporateAccount(ctx context.Context, deviceID uuid.UUID) error

	// BYOD management
	FindBYODDevices(ctx context.Context) ([]*models.Device, error)
	FindCorporateOwnedDevices(ctx context.Context) ([]*models.Device, error)
	RegisterBYOD(ctx context.Context, deviceID, employeeID uuid.UUID) error
	UnregisterBYOD(ctx context.Context, deviceID uuid.UUID) error

	// Fleet management
	FindFleetDevices(ctx context.Context, fleetID uuid.UUID) ([]*models.Device, error)
	GetFleetStatistics(ctx context.Context, fleetID uuid.UUID) (map[string]interface{}, error)

	// ======================================
	// VALUATION & FINANCIAL
	// ======================================

	// Valuation operations
	GetTotalValue(ctx context.Context, ownerID uuid.UUID) (float64, error)
	GetCurrentMarketValue(ctx context.Context, deviceID uuid.UUID) (float64, error)
	UpdateMarketValue(ctx context.Context, deviceID uuid.UUID, value float64) error
	GetDepreciationValue(ctx context.Context, deviceID uuid.UUID) (float64, error)
	GetTradeInValue(ctx context.Context, deviceID uuid.UUID) (float64, error)
	CalculateTotalCostOfOwnership(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// ======================================
	// DEVICE ECOSYSTEM SERVICES
	// ======================================

	// Repair management
	FindDevicesInRepair(ctx context.Context) ([]*models.Device, error)
	GetRepairHistory(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)
	GetRepairCosts(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// Rental management
	FindRentalDevices(ctx context.Context) ([]*models.Device, error)
	FindAvailableForRental(ctx context.Context) ([]*models.Device, error)
	GetActiveRentals(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// Trade-in management
	FindEligibleForTradeIn(ctx context.Context) ([]*models.Device, error)
	GetTradeInHistory(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// Subscription management
	FindDevicesWithActiveSubscriptions(ctx context.Context) ([]*models.Device, error)
	GetSubscriptionDetails(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// NETWORK & CONNECTIVITY
	// ======================================

	// Network operations
	FindByNetworkOperator(ctx context.Context, operator string) ([]*models.Device, error)
	FindByNetworkStatus(ctx context.Context, status string) ([]*models.Device, error)
	UpdateNetworkInfo(ctx context.Context, deviceID uuid.UUID, operator, status string) error
	GetDataUsageHistory(ctx context.Context, deviceID uuid.UUID, startDate, endDate time.Time) ([]map[string]interface{}, error)

	// ======================================
	// COMPLIANCE & REGULATORY
	// ======================================

	// Compliance checks
	FindNonCompliantDevices(ctx context.Context) ([]*models.Device, error)
	CheckComplianceStatus(ctx context.Context, deviceID uuid.UUID) (bool, []string, error)
	UpdateComplianceStatus(ctx context.Context, deviceID uuid.UUID, status string, issues []string) error
	GetComplianceHistory(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// ANALYTICS & REPORTING
	// ======================================

	// Statistics and aggregations
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountByOwner(ctx context.Context, ownerID uuid.UUID) (int64, error)
	CountByBrand(ctx context.Context, brand string) (int64, error)
	CountByCategory(ctx context.Context, category string) (int64, error)
	GetStatisticsByDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)

	// Performance metrics
	GetPerformanceMetrics(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	GetBatteryHealth(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	GetStorageHealth(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)

	// Usage analytics
	GetUsagePatterns(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)
	GetBehaviorScore(ctx context.Context, deviceID uuid.UUID) (float64, error)

	// ======================================
	// WARRANTY & SERVICE MANAGEMENT
	// ======================================

	// Warranty operations
	FindUnderWarranty(ctx context.Context) ([]*models.Device, error)
	FindWarrantyExpiring(ctx context.Context, days int) ([]*models.Device, error)
	UpdateWarrantyStatus(ctx context.Context, deviceID uuid.UUID, status string, expiryDate time.Time) error
	GetWarrantyHistory(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// QUALITY & CERTIFICATION
	// ======================================

	// Quality management
	GetQualityMetrics(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	FindDefectiveDevices(ctx context.Context) ([]*models.Device, error)
	GetRecallStatus(ctx context.Context, deviceID uuid.UUID) (bool, map[string]interface{}, error)
	GetCertificationStatus(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// EMERGENCY & RECOVERY
	// ======================================

	// Emergency operations
	FindDevicesInEmergencyMode(ctx context.Context) ([]*models.Device, error)
	ActivateEmergencyMode(ctx context.Context, deviceID uuid.UUID) error
	DeactivateEmergencyMode(ctx context.Context, deviceID uuid.UUID) error
	GetEmergencyContacts(ctx context.Context, deviceID uuid.UUID) ([]map[string]interface{}, error)

	// Backup and recovery
	GetBackupStatus(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	UpdateBackupStatus(ctx context.Context, deviceID uuid.UUID, status map[string]interface{}) error

	// ======================================
	// ENVIRONMENTAL & SUSTAINABILITY
	// ======================================

	// Sustainability metrics
	GetCarbonFootprint(ctx context.Context, deviceID uuid.UUID) (float64, error)
	GetRecyclingScore(ctx context.Context, deviceID uuid.UUID) (float64, error)
	GetSustainabilityMetrics(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)
	FindEcoFriendlyDevices(ctx context.Context) ([]*models.Device, error)

	// ======================================
	// ADVANCED SEARCH & FILTERING
	// ======================================

	// Complex queries with multiple criteria
	FindByComplexCriteria(ctx context.Context, criteria map[string]interface{}) ([]*models.Device, error)
	GetDevicesNearLocation(ctx context.Context, latitude, longitude float64, radiusKm float64) ([]*models.Device, error)
	FindByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*models.Device, error)
	FindByAgeRange(ctx context.Context, minMonths, maxMonths int) ([]*models.Device, error)
}
