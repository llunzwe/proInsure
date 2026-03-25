// Package repositories defines the repository interfaces following hexagonal architecture
package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/claim"
)

// ClaimSearchFilters defines comprehensive search criteria for claims
type ClaimSearchFilters struct {
	// Basic filters
	Status     []string
	ClaimType  string
	Priority   string
	CustomerID *uuid.UUID
	PolicyID   *uuid.UUID
	DeviceID   *uuid.UUID

	// Date range filters
	IncidentDateFrom   *time.Time
	IncidentDateTo     *time.Time
	ReportedDateFrom   *time.Time
	ReportedDateTo     *time.Time
	SettlementDateFrom *time.Time
	SettlementDateTo   *time.Time

	// Financial filters
	ClaimedAmountMin  float64
	ClaimedAmountMax  float64
	ApprovedAmountMin float64
	ApprovedAmountMax float64

	// Investigation filters
	RequiresInvestigation *bool
	FraudScoreMin         float64
	FraudScoreMax         float64
	FraudRiskLevel        string

	// Assignment filters
	AdjusterID     *uuid.UUID
	InvestigatorID *uuid.UUID
	WorkflowStage  string

	// Advanced filters
	HasLitigation  *bool
	HasAppeal      *bool
	IsHighPriority *bool
	SLABreached    *bool

	// Pagination
	Limit  int
	Offset int
	SortBy string
	Order  string
}

// ClaimStatistics represents aggregated claim statistics
type ClaimStatistics struct {
	TotalClaims           int64
	OpenClaims            int64
	ClosedClaims          int64
	ApprovedClaims        int64
	DeniedClaims          int64
	TotalClaimedAmount    float64
	TotalApprovedAmount   float64
	TotalSettledAmount    float64
	AverageFraudScore     float64
	AverageProcessingDays int
}

// ClaimRepository defines the comprehensive contract for claim persistence
// This interface covers all aspects of claim management including lifecycle,
// investigation, fraud detection, workflow, settlement, and reporting
type ClaimRepository interface {
	// ======================================
	// BASIC CRUD OPERATIONS
	// ======================================

	// Core CRUD operations
	Create(ctx context.Context, claim *models.Claim) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Claim, error)
	GetByIDWithRelations(ctx context.Context, id uuid.UUID, relations []string) (*models.Claim, error)
	GetByClaimNumber(ctx context.Context, claimNumber string) (*models.Claim, error)
	Update(ctx context.Context, claim *models.Claim) error
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error

	// Batch operations
	CreateBatch(ctx context.Context, claims []*models.Claim) error
	UpdateBatch(ctx context.Context, claims []*models.Claim) error

	// ======================================
	// CLAIM SEARCH & LISTING
	// ======================================

	// Search and filtering
	Search(ctx context.Context, filters ClaimSearchFilters) ([]*models.Claim, int64, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*models.Claim, error)
	GetByPolicyID(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*models.Claim, error)
	GetByDeviceID(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]*models.Claim, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Claim, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.Claim, error)

	// Advanced queries
	GetOpenClaims(ctx context.Context, limit, offset int) ([]*models.Claim, error)
	GetPendingApprovalClaims(ctx context.Context, limit, offset int) ([]*models.Claim, error)
	GetHighPriorityClaims(ctx context.Context, limit, offset int) ([]*models.Claim, error)
	GetClaimsNearingSLA(ctx context.Context, hoursThreshold int) ([]*models.Claim, error)
	GetSLABreachedClaims(ctx context.Context) ([]*models.Claim, error)

	// ======================================
	// CLAIM LIFECYCLE MANAGEMENT
	// ======================================

	// Status transitions
	UpdateStatus(ctx context.Context, claimID uuid.UUID, status string, updatedBy uuid.UUID) error
	ApproveClaim(ctx context.Context, claimID uuid.UUID, approvedAmount float64, approverID uuid.UUID) error
	DenyClaim(ctx context.Context, claimID uuid.UUID, denialReason string, denierID uuid.UUID) error
	CloseClaim(ctx context.Context, claimID uuid.UUID, closureReason string) error
	ReopenClaim(ctx context.Context, claimID uuid.UUID, reopenReason string, reopenerID uuid.UUID) error
	WithdrawClaim(ctx context.Context, claimID uuid.UUID, withdrawalReason string) error

	// ======================================
	// INVESTIGATION & FRAUD DETECTION
	// ======================================

	// Investigation management
	CreateInvestigation(ctx context.Context, investigation *claim.ClaimInvestigationDetail) error
	GetInvestigationByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimInvestigationDetail, error)
	UpdateInvestigation(ctx context.Context, investigation *claim.ClaimInvestigationDetail) error
	AssignInvestigator(ctx context.Context, claimID, investigatorID uuid.UUID) error

	// Fraud detection
	CreateFraudDetection(ctx context.Context, fraudDetection *claim.ClaimFraudDetection) error
	GetFraudDetectionByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimFraudDetection, error)
	UpdateFraudScore(ctx context.Context, claimID uuid.UUID, fraudScore float64, indicators map[string]interface{}) error
	GetHighFraudRiskClaims(ctx context.Context, scoreThreshold float64) ([]*models.Claim, error)

	// ======================================
	// WORKFLOW & ASSIGNMENT
	// ======================================

	// Workflow management
	CreateWorkflow(ctx context.Context, workflow *claim.ClaimWorkflow) error
	GetWorkflowByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimWorkflow, error)
	UpdateWorkflowStage(ctx context.Context, claimID uuid.UUID, stage string, completedBy uuid.UUID) error
	GetWorkflowHistory(ctx context.Context, claimID uuid.UUID) ([]map[string]interface{}, error)

	// Assignment management
	AssignAdjuster(ctx context.Context, claimID, adjusterID uuid.UUID) error
	GetClaimsByAdjuster(ctx context.Context, adjusterID uuid.UUID) ([]*models.Claim, error)
	GetUnassignedClaims(ctx context.Context) ([]*models.Claim, error)
	ReassignClaim(ctx context.Context, claimID, fromAdjusterID, toAdjusterID uuid.UUID, reason string) error

	// ======================================
	// SETTLEMENT & PAYMENT
	// ======================================

	// Settlement management
	CreateSettlement(ctx context.Context, settlement *claim.ClaimSettlementDetail) error
	GetSettlementByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimSettlementDetail, error)
	UpdateSettlement(ctx context.Context, settlement *claim.ClaimSettlementDetail) error
	GetPendingSettlements(ctx context.Context) ([]*claim.ClaimSettlementDetail, error)

	// Payment management
	CreatePayment(ctx context.Context, payment *claim.ClaimPayment) error
	GetPaymentsByClaimID(ctx context.Context, claimID uuid.UUID) ([]*claim.ClaimPayment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status string) error
	GetPendingPayments(ctx context.Context) ([]*claim.ClaimPayment, error)

	// Reserve management
	CreateReserve(ctx context.Context, reserve *claim.ClaimReserve) error
	GetReserveByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimReserve, error)
	UpdateReserve(ctx context.Context, reserve *claim.ClaimReserve) error
	GetTotalReserves(ctx context.Context, policyID uuid.UUID) (float64, error)

	// ======================================
	// APPEAL & LITIGATION
	// ======================================

	// Appeal management
	CreateAppeal(ctx context.Context, appeal *claim.ClaimAppeal) error
	GetAppealByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimAppeal, error)
	UpdateAppealStatus(ctx context.Context, appealID uuid.UUID, status, decision string) error
	GetPendingAppeals(ctx context.Context) ([]*claim.ClaimAppeal, error)

	// Litigation management
	CreateLitigation(ctx context.Context, litigation *claim.ClaimLitigation) error
	GetLitigationByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimLitigation, error)
	UpdateLitigationStatus(ctx context.Context, litigationID uuid.UUID, status string) error
	GetActiveLitigations(ctx context.Context) ([]*claim.ClaimLitigation, error)

	// ======================================
	// SMARTPHONE-SPECIFIC FEATURES
	// ======================================

	// Device diagnostics
	CreateDeviceDiagnostics(ctx context.Context, diagnostics *claim.ClaimDeviceDiagnostics) error
	GetDeviceDiagnosticsByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimDeviceDiagnostics, error)
	UpdateDiagnosticsResults(ctx context.Context, claimID uuid.UUID, results map[string]interface{}) error

	// Repair network
	CreateRepairNetwork(ctx context.Context, repairNetwork *claim.ClaimRepairNetwork) error
	GetRepairNetworkByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimRepairNetwork, error)
	AssignRepairCenter(ctx context.Context, claimID, repairCenterID uuid.UUID) error

	// Replacement device
	CreateReplacementDevice(ctx context.Context, replacement *claim.ClaimReplacementDevice) error
	GetReplacementDeviceByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimReplacementDevice, error)
	UpdateReplacementStatus(ctx context.Context, claimID uuid.UUID, status string) error

	// ======================================
	// REPORTING & ANALYTICS
	// ======================================

	// Statistics
	GetStatistics(ctx context.Context, filters ClaimSearchFilters) (*ClaimStatistics, error)
	GetStatisticsByPeriod(ctx context.Context, startDate, endDate time.Time) (*ClaimStatistics, error)
	GetStatisticsByCustomer(ctx context.Context, customerID uuid.UUID) (*ClaimStatistics, error)
	GetStatisticsByPolicy(ctx context.Context, policyID uuid.UUID) (*ClaimStatistics, error)

	// Analytics
	GetClaimTrends(ctx context.Context, period string, limit int) ([]map[string]interface{}, error)
	GetFraudTrends(ctx context.Context, period string, limit int) ([]map[string]interface{}, error)
	GetProcessingTimeAnalytics(ctx context.Context) (map[string]interface{}, error)
	GetLossRatioByProduct(ctx context.Context, productID uuid.UUID) (float64, error)

	// Reporting
	GenerateClaimReport(ctx context.Context, claimID uuid.UUID, reportType string) ([]byte, error)
	GetReportingMetrics(ctx context.Context, startDate, endDate time.Time) (*claim.ClaimReporting, error)
	CreateReportSchedule(ctx context.Context, schedule *claim.ClaimReportSchedule) error

	// ======================================
	// COMPLIANCE & AUDIT
	// ======================================

	// Audit trail
	GetAuditTrail(ctx context.Context, claimID uuid.UUID) ([]map[string]interface{}, error)
	LogClaimActivity(ctx context.Context, claimID uuid.UUID, activity, performedBy string, details map[string]interface{}) error

	// Compliance
	GetComplianceStatus(ctx context.Context, claimID uuid.UUID) (map[string]interface{}, error)
	CheckRegulatoryCompliance(ctx context.Context, claimID uuid.UUID) (bool, []string, error)
	GenerateComplianceReport(ctx context.Context, startDate, endDate time.Time) ([]byte, error)

	// ======================================
	// BULK OPERATIONS & OPTIMIZATION
	// ======================================

	// Bulk updates
	BulkUpdateStatus(ctx context.Context, claimIDs []uuid.UUID, status string) error
	BulkAssignAdjuster(ctx context.Context, claimIDs []uuid.UUID, adjusterID uuid.UUID) error
	BulkUpdatePriority(ctx context.Context, claimIDs []uuid.UUID, priority string) error

	// Performance optimization
	PreloadRelations(ctx context.Context, claims []*models.Claim, relations []string) error
	GetClaimsForBatchProcessing(ctx context.Context, batchSize int) ([]*models.Claim, error)
	OptimizeClaimQueries(ctx context.Context) error

	// ======================================
	// COMMUNICATION & NOTIFICATIONS
	// ======================================

	// Communication tracking
	CreateCommunication(ctx context.Context, communication *claim.ClaimCommunication) error
	GetCommunicationsByClaimID(ctx context.Context, claimID uuid.UUID) ([]*claim.ClaimCommunication, error)
	LogCustomerInteraction(ctx context.Context, claimID uuid.UUID, interactionType, channel, notes string) error
	GetCustomerSatisfactionScore(ctx context.Context, claimID uuid.UUID) (float64, error)
}
