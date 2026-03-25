package device

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// DeviceQueryBuilder provides optimized query building for Device entities
type DeviceQueryBuilder struct {
	db         *gorm.DB
	preloads   []string
	conditions map[string]interface{}
	joins      []string
	selects    []string
	limit      int
	offset     int
	orderBy    string
	tx         *gorm.DB // For transaction support
}

// NewDeviceQueryBuilder creates a new query builder
func NewDeviceQueryBuilder(db *gorm.DB) *DeviceQueryBuilder {
	return &DeviceQueryBuilder{
		db:         db,
		preloads:   []string{},
		conditions: make(map[string]interface{}),
		joins:      []string{},
		selects:    []string{},
	}
}

// WithTransaction uses a transaction for the query
func (qb *DeviceQueryBuilder) WithTransaction(tx *gorm.DB) *DeviceQueryBuilder {
	qb.tx = tx
	return qb
}

// === Core Relationships ===

// WithOwner preloads the device owner
func (qb *DeviceQueryBuilder) WithOwner() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Owner")
	return qb
}

// WithPolicies preloads device policies
func (qb *DeviceQueryBuilder) WithPolicies() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Policies")
	return qb
}

// WithActivePolicies preloads only active policies
func (qb *DeviceQueryBuilder) WithActivePolicies() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Policies.Active")
	return qb
}

// WithClaims preloads device claims
func (qb *DeviceQueryBuilder) WithClaims() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Claims")
	return qb
}

// WithCorporateAccount preloads corporate account if exists
func (qb *DeviceQueryBuilder) WithCorporateAccount() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "CorporateAccount")
	return qb
}

// === Device Ecosystem Relationships ===

// WithRepairs preloads device repairs
func (qb *DeviceQueryBuilder) WithRepairs() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Repairs")
	return qb
}

// WithRentals preloads device rentals
func (qb *DeviceQueryBuilder) WithRentals() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Rentals")
	return qb
}

// WithTradeIns preloads device trade-ins
func (qb *DeviceQueryBuilder) WithTradeIns() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "TradeIns")
	return qb
}

// WithSubscriptions preloads device subscriptions
func (qb *DeviceQueryBuilder) WithSubscriptions() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Subscriptions")
	return qb
}

// WithFinancing preloads financing information
func (qb *DeviceQueryBuilder) WithFinancing() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Financings")
	return qb
}

// WithMarketListings preloads marketplace listings
func (qb *DeviceQueryBuilder) WithMarketListings() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "MarketListings")
	return qb
}

// === Insurance & Risk Relationships ===

// WithRiskProfile preloads risk profile
func (qb *DeviceQueryBuilder) WithRiskProfile() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "RiskProfile")
	return qb
}

// WithClaimHistory preloads claim history
func (qb *DeviceQueryBuilder) WithClaimHistory() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "ClaimHistory")
	return qb
}

// WithPremiumCalculations preloads premium calculations
func (qb *DeviceQueryBuilder) WithPremiumCalculations() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "PremiumCalculations")
	return qb
}

// WithInsuranceAudits preloads insurance audits
func (qb *DeviceQueryBuilder) WithInsuranceAudits() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "InsuranceAudits")
	return qb
}

// === Fraud & Security Relationships ===

// WithFraudPatterns preloads fraud patterns
func (qb *DeviceQueryBuilder) WithFraudPatterns() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "FraudPatterns")
	return qb
}

// WithAnomalyDetections preloads anomaly detections
func (qb *DeviceQueryBuilder) WithAnomalyDetections() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "AnomalyDetections")
	return qb
}

// WithVerifications preloads device verifications
func (qb *DeviceQueryBuilder) WithVerifications() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "Verifications")
	return qb
}

// WithDeviceHistory preloads device history
func (qb *DeviceQueryBuilder) WithDeviceHistory() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "DeviceHistory")
	return qb
}

// === Analytics & Intelligence ===

// WithUsagePatterns preloads usage patterns
func (qb *DeviceQueryBuilder) WithUsagePatterns() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "UsagePatterns")
	return qb
}

// WithBehaviorScore preloads behavior score
func (qb *DeviceQueryBuilder) WithBehaviorScore() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "BehaviorScore")
	return qb
}

// WithFailurePredictions preloads failure predictions
func (qb *DeviceQueryBuilder) WithFailurePredictions() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "FailurePredictions")
	return qb
}

// WithMaintenanceSchedule preloads maintenance schedule
func (qb *DeviceQueryBuilder) WithMaintenanceSchedule() *DeviceQueryBuilder {
	qb.preloads = append(qb.preloads, "MaintenanceSchedule")
	return qb
}

// === Composite Preloads ===

// WithEssentials preloads essential relationships
func (qb *DeviceQueryBuilder) WithEssentials() *DeviceQueryBuilder {
	return qb.WithOwner().WithPolicies().WithClaims()
}

// WithInsuranceData preloads all insurance-related data
func (qb *DeviceQueryBuilder) WithInsuranceData() *DeviceQueryBuilder {
	return qb.WithPolicies().
		WithClaims().
		WithRiskProfile().
		WithClaimHistory().
		WithPremiumCalculations()
}

// WithFullProfile preloads complete device profile
func (qb *DeviceQueryBuilder) WithFullProfile() *DeviceQueryBuilder {
	return qb.WithOwner().
		WithPolicies().
		WithClaims().
		WithRiskProfile().
		WithRepairs().
		WithVerifications().
		WithDeviceHistory()
}

// === Conditions ===

// WhereID filters by device ID
func (qb *DeviceQueryBuilder) WhereID(id uuid.UUID) *DeviceQueryBuilder {
	qb.conditions["id"] = id
	return qb
}

// WhereIMEI filters by IMEI
func (qb *DeviceQueryBuilder) WhereIMEI(imei string) *DeviceQueryBuilder {
	qb.conditions["device_identification.imei"] = imei
	return qb
}

// WhereOwner filters by owner ID
func (qb *DeviceQueryBuilder) WhereOwner(ownerID uuid.UUID) *DeviceQueryBuilder {
	qb.conditions["device_ownership.owner_id"] = ownerID
	return qb
}

// WhereStatus filters by device status
func (qb *DeviceQueryBuilder) WhereStatus(status string) *DeviceQueryBuilder {
	qb.conditions["device_lifecycle.status"] = status
	return qb
}

// WhereStatusIn filters by multiple statuses
func (qb *DeviceQueryBuilder) WhereStatusIn(statuses []string) *DeviceQueryBuilder {
	qb.conditions["device_lifecycle.status IN ?"] = statuses
	return qb
}

// WhereBrand filters by brand
func (qb *DeviceQueryBuilder) WhereBrand(brand string) *DeviceQueryBuilder {
	qb.conditions["device_classification.brand"] = brand
	return qb
}

// WhereModel filters by model
func (qb *DeviceQueryBuilder) WhereModel(model string) *DeviceQueryBuilder {
	qb.conditions["device_classification.model"] = model
	return qb
}

// WhereCondition filters by condition
func (qb *DeviceQueryBuilder) WhereCondition(condition string) *DeviceQueryBuilder {
	qb.conditions["device_physical_condition.condition"] = condition
	return qb
}

// WhereGrade filters by grade
func (qb *DeviceQueryBuilder) WhereGrade(grade string) *DeviceQueryBuilder {
	qb.conditions["device_physical_condition.grade"] = grade
	return qb
}

// WhereRiskLevel filters by risk level
func (qb *DeviceQueryBuilder) WhereRiskLevel(level string) *DeviceQueryBuilder {
	qb.conditions["device_risk_assessment.risk_level"] = level
	return qb
}

// WhereRiskScoreGreaterThan filters by risk score threshold
func (qb *DeviceQueryBuilder) WhereRiskScoreGreaterThan(score float64) *DeviceQueryBuilder {
	qb.conditions["device_risk_assessment.risk_score > ?"] = score
	return qb
}

// WhereBlacklisted filters blacklisted devices
func (qb *DeviceQueryBuilder) WhereBlacklisted() *DeviceQueryBuilder {
	qb.conditions["device_risk_assessment.blacklist_status"] = "blacklisted"
	return qb
}

// WhereStolen filters stolen devices
func (qb *DeviceQueryBuilder) WhereStolen() *DeviceQueryBuilder {
	qb.conditions["device_lifecycle.is_stolen"] = true
	return qb
}

// WhereVerified filters verified devices
func (qb *DeviceQueryBuilder) WhereVerified() *DeviceQueryBuilder {
	qb.conditions["device_lifecycle.is_verified"] = true
	return qb
}

// WhereCorporate filters corporate devices
func (qb *DeviceQueryBuilder) WhereCorporate() *DeviceQueryBuilder {
	qb.conditions["device_ownership.corporate_account_id IS NOT NULL"] = nil
	return qb
}

// === Sorting ===

// OrderBy sets the order clause
func (qb *DeviceQueryBuilder) OrderBy(field string) *DeviceQueryBuilder {
	qb.orderBy = field
	return qb
}

// OrderByCreatedDesc orders by creation date descending
func (qb *DeviceQueryBuilder) OrderByCreatedDesc() *DeviceQueryBuilder {
	qb.orderBy = "created_at DESC"
	return qb
}

// OrderByRiskScore orders by risk score
func (qb *DeviceQueryBuilder) OrderByRiskScore(desc bool) *DeviceQueryBuilder {
	if desc {
		qb.orderBy = "device_risk_assessment.risk_score DESC"
	} else {
		qb.orderBy = "device_risk_assessment.risk_score ASC"
	}
	return qb
}

// OrderByValue orders by current value
func (qb *DeviceQueryBuilder) OrderByValue(desc bool) *DeviceQueryBuilder {
	if desc {
		qb.orderBy = "device_financial.current_value DESC"
	} else {
		qb.orderBy = "device_financial.current_value ASC"
	}
	return qb
}

// === Pagination ===

// Limit sets the limit
func (qb *DeviceQueryBuilder) Limit(limit int) *DeviceQueryBuilder {
	qb.limit = limit
	return qb
}

// Offset sets the offset
func (qb *DeviceQueryBuilder) Offset(offset int) *DeviceQueryBuilder {
	qb.offset = offset
	return qb
}

// Paginate sets limit and offset for pagination
func (qb *DeviceQueryBuilder) Paginate(page, pageSize int) *DeviceQueryBuilder {
	qb.limit = pageSize
	qb.offset = (page - 1) * pageSize
	return qb
}

// === Selective Loading ===

// SelectBasic selects only basic device fields
func (qb *DeviceQueryBuilder) SelectBasic() *DeviceQueryBuilder {
	qb.selects = append(qb.selects, "id", "device_identification.*",
		"device_classification.*", "device_lifecycle.status")
	return qb
}

// SelectFinancial selects financial fields
func (qb *DeviceQueryBuilder) SelectFinancial() *DeviceQueryBuilder {
	qb.selects = append(qb.selects, "id", "device_financial.*")
	return qb
}

// SelectRisk selects risk-related fields
func (qb *DeviceQueryBuilder) SelectRisk() *DeviceQueryBuilder {
	qb.selects = append(qb.selects, "id", "device_risk_assessment.*")
	return qb
}

// === Execution Methods ===

// Find executes the query and returns multiple devices
func (qb *DeviceQueryBuilder) Find(ctx context.Context) ([]*models.Device, error) {
	query := qb.buildQuery()

	var devices []*models.Device
	if err := query.WithContext(ctx).Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to find devices: %w", err)
	}

	return devices, nil
}

// First executes the query and returns the first device
func (qb *DeviceQueryBuilder) First(ctx context.Context) (*models.Device, error) {
	query := qb.buildQuery()

	var device models.Device
	if err := query.WithContext(ctx).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return &device, nil
}

// Count returns the count of devices matching the conditions
func (qb *DeviceQueryBuilder) Count(ctx context.Context) (int64, error) {
	query := qb.buildBaseQuery()

	var count int64
	if err := query.WithContext(ctx).Model(&models.Device{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count devices: %w", err)
	}

	return count, nil
}

// Exists checks if any device matches the conditions
func (qb *DeviceQueryBuilder) Exists(ctx context.Context) (bool, error) {
	count, err := qb.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// buildQuery builds the complete query with all options
func (qb *DeviceQueryBuilder) buildQuery() *gorm.DB {
	query := qb.buildBaseQuery()

	// Apply preloads
	for _, preload := range qb.preloads {
		query = query.Preload(preload)
	}

	// Apply selects
	if len(qb.selects) > 0 {
		query = query.Select(qb.selects)
	}

	// Apply ordering
	if qb.orderBy != "" {
		query = query.Order(qb.orderBy)
	}

	// Apply limit and offset
	if qb.limit > 0 {
		query = query.Limit(qb.limit)
	}
	if qb.offset > 0 {
		query = query.Offset(qb.offset)
	}

	return query
}

// buildBaseQuery builds the base query with conditions
func (qb *DeviceQueryBuilder) buildBaseQuery() *gorm.DB {
	var query *gorm.DB

	if qb.tx != nil {
		query = qb.tx.Model(&models.Device{})
	} else {
		query = qb.db.Model(&models.Device{})
	}

	// Apply conditions
	for key, value := range qb.conditions {
		if value == nil {
			query = query.Where(key)
		} else {
			query = query.Where(key, value)
		}
	}

	// Apply joins
	for _, join := range qb.joins {
		query = query.Joins(join)
	}

	return query
}

// Reset clears all builder options
func (qb *DeviceQueryBuilder) Reset() *DeviceQueryBuilder {
	qb.preloads = []string{}
	qb.conditions = make(map[string]interface{})
	qb.joins = []string{}
	qb.selects = []string{}
	qb.limit = 0
	qb.offset = 0
	qb.orderBy = ""
	qb.tx = nil
	return qb
}
