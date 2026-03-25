package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/policy"
)

// PolicyRepository defines the comprehensive contract for policy persistence
// This interface covers all aspects of policy management including lifecycle,
// coverage, claims, underwriting, compliance, and corporate features
type PolicyRepository interface {
	// ======================================
	// BASIC CRUD OPERATIONS
	// ======================================

	// Core CRUD operations
	Create(ctx context.Context, policy *models.Policy) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Policy, error)
	GetByIDWithRelations(ctx context.Context, id uuid.UUID, relations []string) (*models.Policy, error)
	GetByPolicyNumber(ctx context.Context, policyNumber string) (*models.Policy, error)
	Update(ctx context.Context, policy *models.Policy) error
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error

	// Batch operations
	CreateBatch(ctx context.Context, policies []*models.Policy) error
	UpdateBatch(ctx context.Context, policies []*models.Policy) error
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error

	// ======================================
	// POLICY IDENTIFICATION & SEARCH
	// ======================================

	// Search operations
	Search(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*models.Policy, int64, error)
	FindByContractNumber(ctx context.Context, contractNumber string) (*models.Policy, error)
	FindByExternalReference(ctx context.Context, reference string) (*models.Policy, error)
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]*models.Policy, error)
	FindByPolicyCategory(ctx context.Context, category string) ([]*models.Policy, error)
	FindByBusinessLine(ctx context.Context, businessLine string) ([]*models.Policy, error)

	// ======================================
	// CUSTOMER & OWNERSHIP MANAGEMENT
	// ======================================

	// Customer related queries
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*models.Policy, error)
	FindByCustomerIDWithPagination(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*models.Policy, int64, error)
	FindActiveByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*models.Policy, error)
	TransferPolicyOwner(ctx context.Context, policyID, fromCustomerID, toCustomerID uuid.UUID) error
	GetCustomerPolicyHistory(ctx context.Context, customerID uuid.UUID) ([]*models.Policy, error)

	// Device related queries
	FindByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*models.Policy, error)
	FindActiveByDeviceID(ctx context.Context, deviceID uuid.UUID) (*models.Policy, error)
	GetDevicePolicyHistory(ctx context.Context, deviceID uuid.UUID) ([]*models.Policy, error)

	// ======================================
	// POLICY STATUS & LIFECYCLE
	// ======================================

	// Status operations
	FindByStatus(ctx context.Context, status string) ([]*models.Policy, error)
	FindByMultipleStatuses(ctx context.Context, statuses []string) ([]*models.Policy, error)
	UpdateStatus(ctx context.Context, policyID uuid.UUID, status string) error
	UpdateStatusBatch(ctx context.Context, policyIDs []uuid.UUID, status string) error

	// Lifecycle queries
	FindActive(ctx context.Context) ([]*models.Policy, error)
	FindExpired(ctx context.Context) ([]*models.Policy, error)
	FindCancelled(ctx context.Context) ([]*models.Policy, error)
	FindSuspended(ctx context.Context) ([]*models.Policy, error)
	FindPendingActivation(ctx context.Context) ([]*models.Policy, error)
	FindExpiringWithinDays(ctx context.Context, days int) ([]*models.Policy, error)

	// Policy actions
	ActivatePolicy(ctx context.Context, policyID uuid.UUID) error
	SuspendPolicy(ctx context.Context, policyID uuid.UUID, reason string) error
	ReinstatePolicy(ctx context.Context, policyID uuid.UUID) error
	CancelPolicy(ctx context.Context, policyID uuid.UUID, reason string, effectiveDate time.Time) error
	ExpirePolicy(ctx context.Context, policyID uuid.UUID) error

	// ======================================
	// RENEWAL MANAGEMENT
	// ======================================

	// Renewal operations
	FindPendingRenewal(ctx context.Context) ([]*models.Policy, error)
	FindRenewalsDueWithinDays(ctx context.Context, days int) ([]*models.Policy, error)
	CreateRenewal(ctx context.Context, originalPolicyID uuid.UUID, renewal *models.Policy) error
	GetRenewalHistory(ctx context.Context, policyID uuid.UUID) ([]*models.Policy, error)
	MarkAsRenewed(ctx context.Context, policyID uuid.UUID, newPolicyID uuid.UUID) error
	UpdateAutoRenewalStatus(ctx context.Context, policyID uuid.UUID, enabled bool) error

	// ======================================
	// COVERAGE & BENEFITS MANAGEMENT
	// ======================================

	// Coverage queries
	FindByCoverageType(ctx context.Context, coverageType string) ([]*models.Policy, error)
	FindByCoverageLevel(ctx context.Context, level string) ([]*models.Policy, error)
	FindWithSpecificCoverage(ctx context.Context, coverageCode string) ([]*models.Policy, error)
	UpdateCoverage(ctx context.Context, policyID uuid.UUID, coverage *policy.PolicyCoverage) error
	GetCoverageHistory(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)

	// Benefits and limits
	GetPolicyBenefits(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyBenefit, error)
	GetPolicyLimits(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyLimit, error)
	GetPolicyExclusions(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyExclusion, error)
	UpdateBenefits(ctx context.Context, policyID uuid.UUID, benefits []*policy.PolicyBenefit) error

	// ======================================
	// PRICING & FINANCIAL
	// ======================================

	// Premium operations
	GetTotalPremium(ctx context.Context, policyID uuid.UUID) (float64, error)
	UpdatePremium(ctx context.Context, policyID uuid.UUID, premium float64) error
	CalculatePremiumAdjustment(ctx context.Context, policyID uuid.UUID) (float64, error)
	GetPremiumHistory(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)

	// Discount management
	ApplyDiscount(ctx context.Context, policyID uuid.UUID, discount *policy.PolicyDiscount) error
	RemoveDiscount(ctx context.Context, policyID uuid.UUID, discountID uuid.UUID) error
	GetActiveDiscounts(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyDiscount, error)
	CalculateTotalDiscount(ctx context.Context, policyID uuid.UUID) (float64, error)

	// Payment related
	GetPaymentSchedule(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyPaymentSchedule, error)
	UpdatePaymentFrequency(ctx context.Context, policyID uuid.UUID, frequency string) error
	GetOutstandingAmount(ctx context.Context, policyID uuid.UUID) (float64, error)
	RecordPayment(ctx context.Context, policyID uuid.UUID, payment *models.Payment) error

	// ======================================
	// CLAIMS MANAGEMENT
	// ======================================

	// Claims operations
	GetClaims(ctx context.Context, policyID uuid.UUID) ([]*models.Claim, error)
	GetActiveClaims(ctx context.Context, policyID uuid.UUID) ([]*models.Claim, error)
	GetClaimStatistics(ctx context.Context, policyID uuid.UUID) (map[string]interface{}, error)
	CountClaims(ctx context.Context, policyID uuid.UUID) (int64, error)
	GetTotalClaimAmount(ctx context.Context, policyID uuid.UUID) (float64, error)
	CheckClaimLimits(ctx context.Context, policyID uuid.UUID, claimType string) (bool, string, error)
	UpdateClaimCount(ctx context.Context, policyID uuid.UUID) error

	// ======================================
	// UNDERWRITING & RISK ASSESSMENT
	// ======================================

	// Underwriting operations
	FindPendingUnderwriting(ctx context.Context) ([]*models.Policy, error)
	UpdateUnderwritingStatus(ctx context.Context, policyID uuid.UUID, status string, notes string) error
	AssignUnderwriter(ctx context.Context, policyID uuid.UUID, underwriterID uuid.UUID) error
	GetUnderwritingHistory(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)
	ApproveUnderwriting(ctx context.Context, policyID uuid.UUID, underwriterID uuid.UUID) error
	RejectUnderwriting(ctx context.Context, policyID uuid.UUID, underwriterID uuid.UUID, reason string) error

	// Risk assessment
	UpdateRiskScore(ctx context.Context, policyID uuid.UUID, scores map[string]float64) error
	GetRiskProfile(ctx context.Context, policyID uuid.UUID) (map[string]interface{}, error)
	FindHighRiskPolicies(ctx context.Context, threshold float64) ([]*models.Policy, error)
	FindLowRiskPolicies(ctx context.Context, threshold float64) ([]*models.Policy, error)
	CalculateLossRatio(ctx context.Context, policyID uuid.UUID) (float64, error)

	// ======================================
	// CORPORATE & ENTERPRISE FEATURES
	// ======================================

	// Corporate account management
	FindByCorporateAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.Policy, error)
	FindByCorporateDepartment(ctx context.Context, departmentID uuid.UUID) ([]*models.Policy, error)
	GetCorporateFleetPolicies(ctx context.Context, corporateID uuid.UUID) ([]*models.Policy, error)
	UpdateCorporateAssociation(ctx context.Context, policyID, corporateID uuid.UUID) error

	// Bundle management
	FindByBundleID(ctx context.Context, bundleID uuid.UUID) ([]*models.Policy, error)
	CreateBundle(ctx context.Context, policies []*models.Policy, bundle *policy.PolicyBundle) error
	RemoveFromBundle(ctx context.Context, policyID uuid.UUID) error
	GetBundleDiscount(ctx context.Context, bundleID uuid.UUID) (float64, error)

	// Group policies
	FindFamilyPolicies(ctx context.Context, familyID uuid.UUID) ([]*models.Policy, error)
	FindGroupPolicies(ctx context.Context, groupID uuid.UUID) ([]*models.Policy, error)
	AddToFamilyPlan(ctx context.Context, policyID, familyPlanID uuid.UUID) error

	// ======================================
	// COMPLIANCE & REGULATORY
	// ======================================

	// Compliance operations
	FindNonCompliantPolicies(ctx context.Context) ([]*models.Policy, error)
	CheckComplianceStatus(ctx context.Context, policyID uuid.UUID) (bool, []string, error)
	UpdateComplianceStatus(ctx context.Context, policyID uuid.UUID, status string, issues []string) error
	GetComplianceAudits(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)

	// Regulatory reporting
	GetPoliciesForRegion(ctx context.Context, region string) ([]*models.Policy, error)
	GetPoliciesRequiringReporting(ctx context.Context) ([]*models.Policy, error)
	MarkAsReported(ctx context.Context, policyID uuid.UUID, reportType string) error

	// ======================================
	// DOCUMENT & COMMUNICATION MANAGEMENT
	// ======================================

	// Document operations
	GetPolicyDocuments(ctx context.Context, policyID uuid.UUID) ([]*models.Document, error)
	AttachDocument(ctx context.Context, policyID uuid.UUID, document *models.Document) error
	GeneratePolicyDocument(ctx context.Context, policyID uuid.UUID, documentType string) (*models.Document, error)

	// Communication preferences
	GetCommunicationPreferences(ctx context.Context, policyID uuid.UUID) (*policy.PolicyCommunicationPreference, error)
	UpdateCommunicationPreferences(ctx context.Context, policyID uuid.UUID, prefs *policy.PolicyCommunicationPreference) error
	GetCommunicationHistory(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// MODIFICATIONS & ENDORSEMENTS
	// ======================================

	// Modification tracking
	GetModificationHistory(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyModification, error)
	CreateModification(ctx context.Context, modification *policy.PolicyModification) error
	ApproveModification(ctx context.Context, modificationID uuid.UUID, approverID uuid.UUID) error
	RejectModification(ctx context.Context, modificationID uuid.UUID, reason string) error

	// Endorsements and riders
	GetEndorsements(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyEndorsement, error)
	AddEndorsement(ctx context.Context, policyID uuid.UUID, endorsement *policy.PolicyEndorsement) error
	GetRiders(ctx context.Context, policyID uuid.UUID) ([]*policy.PolicyRider, error)
	AddRider(ctx context.Context, policyID uuid.UUID, rider *policy.PolicyRider) error

	// ======================================
	// ANALYTICS & REPORTING
	// ======================================

	// Statistics and metrics
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountByProduct(ctx context.Context, productID uuid.UUID) (int64, error)
	CountByCustomer(ctx context.Context, customerID uuid.UUID) (int64, error)
	GetStatisticsByDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)

	// Performance metrics
	GetPolicyPerformance(ctx context.Context, policyID uuid.UUID) (map[string]interface{}, error)
	GetProfitabilityScore(ctx context.Context, policyID uuid.UUID) (float64, error)
	GetRetentionScore(ctx context.Context, policyID uuid.UUID) (float64, error)
	GetCustomerLifetimeValue(ctx context.Context, customerID uuid.UUID) (float64, error)

	// Analytics
	GetChurnRiskPolicies(ctx context.Context, threshold float64) ([]*models.Policy, error)
	GetHighValuePolicies(ctx context.Context, minPremium float64) ([]*models.Policy, error)
	GetLowPerformingPolicies(ctx context.Context) ([]*models.Policy, error)

	// ======================================
	// QUOTE & CONVERSION
	// ======================================

	// Quote operations
	ConvertQuoteToPolicy(ctx context.Context, quoteID uuid.UUID) (*models.Policy, error)
	GetQuoteConversionRate(ctx context.Context, startDate, endDate time.Time) (float64, error)
	FindByQuoteID(ctx context.Context, quoteID uuid.UUID) (*models.Policy, error)

	// ======================================
	// AGENT & BROKER MANAGEMENT
	// ======================================

	// Agent/Broker operations
	FindByAgentID(ctx context.Context, agentID uuid.UUID) ([]*models.Policy, error)
	FindByBrokerID(ctx context.Context, brokerID uuid.UUID) ([]*models.Policy, error)
	AssignAgent(ctx context.Context, policyID, agentID uuid.UUID) error
	GetAgentPerformance(ctx context.Context, agentID uuid.UUID) (map[string]interface{}, error)

	// ======================================
	// LOYALTY & REWARDS
	// ======================================

	// Loyalty program operations
	GetLoyaltyPoints(ctx context.Context, policyID uuid.UUID) (int, error)
	AddLoyaltyPoints(ctx context.Context, policyID uuid.UUID, points int, reason string) error
	RedeemLoyaltyPoints(ctx context.Context, policyID uuid.UUID, points int, redemptionType string) error
	GetLoyaltyHistory(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)

	// ======================================
	// SMART FEATURES & IOT
	// ======================================

	// Smart features management
	EnableSmartFeatures(ctx context.Context, policyID uuid.UUID, features map[string]bool) error
	GetSmartFeaturesStatus(ctx context.Context, policyID uuid.UUID) (*policy.PolicySmartFeatures, error)
	UpdateDeviceHealth(ctx context.Context, policyID uuid.UUID, healthScore float64) error
	CheckAutoClaimEligibility(ctx context.Context, policyID uuid.UUID, claimData map[string]interface{}) (bool, error)

	// ======================================
	// INTERNATIONAL COVERAGE
	// ======================================

	// International operations
	GetInternationalCoverage(ctx context.Context, policyID uuid.UUID) (*policy.PolicyInternationalCoverage, error)
	UpdateInternationalCoverage(ctx context.Context, policyID uuid.UUID, coverage *policy.PolicyInternationalCoverage) error
	CheckCountryCoverage(ctx context.Context, policyID uuid.UUID, countryCode string) (bool, error)
	RecordInternationalTravel(ctx context.Context, policyID uuid.UUID, countryCode string, startDate, endDate time.Time) error

	// ======================================
	// ENVIRONMENTAL & SUSTAINABILITY
	// ======================================

	// Environmental features
	GetEnvironmentalScore(ctx context.Context, policyID uuid.UUID) (float64, error)
	UpdateEnvironmentalFeatures(ctx context.Context, policyID uuid.UUID, features *policy.PolicyEnvironmental) error
	FindGreenPolicies(ctx context.Context) ([]*models.Policy, error)
	CalculateCarbonOffset(ctx context.Context, policyID uuid.UUID) (float64, error)

	// ======================================
	// AUDIT & HISTORY
	// ======================================

	// Audit operations
	GetAuditTrail(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)
	LogAuditEvent(ctx context.Context, policyID uuid.UUID, event map[string]interface{}) error
	GetVersionHistory(ctx context.Context, policyID uuid.UUID) ([]*models.Policy, error)
	RestoreVersion(ctx context.Context, policyID uuid.UUID, version int) error

	// ======================================
	// ADVANCED SEARCH & FILTERING
	// ======================================

	// Complex queries
	FindByComplexCriteria(ctx context.Context, criteria map[string]interface{}) ([]*models.Policy, error)
	FindByPremiumRange(ctx context.Context, minPremium, maxPremium float64) ([]*models.Policy, error)
	FindByEffectiveDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.Policy, error)
	FindByDeductibleRange(ctx context.Context, minDeductible, maxDeductible float64) ([]*models.Policy, error)
	SearchFullText(ctx context.Context, searchTerm string, limit, offset int) ([]*models.Policy, int64, error)
}
