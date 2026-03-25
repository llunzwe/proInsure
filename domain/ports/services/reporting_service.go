package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	// "smartsure/internal/domain/models" // Unused import
	"smartsure/internal/domain/models/claim"
	"smartsure/internal/domain/models/shared"
)

// ReportingService defines the interface for reporting business logic operations
type ReportingService interface {
	// Report generation
	GenerateReport(ctx context.Context, reportType string, parameters map[string]interface{}) (*shared.Report, error)
	CreateReport(ctx context.Context, report *shared.Report) error
	GetReportByID(ctx context.Context, id uuid.UUID) (*shared.Report, error)
	GetReportByReportNumber(ctx context.Context, reportNumber string) (*shared.Report, error)
	UpdateReport(ctx context.Context, report *shared.Report) error
	DeleteReport(ctx context.Context, id uuid.UUID) error

	// Report queries
	GetReportsByType(ctx context.Context, reportType string, limit, offset int) ([]*shared.Report, int64, error)
	GetReportsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.Report, int64, error)
	GetReportsByPeriod(ctx context.Context, periodStart, periodEnd string, limit, offset int) ([]*shared.Report, int64, error)
	GetScheduledReports(ctx context.Context, limit, offset int) ([]*shared.Report, int64, error)

	// Dashboard management
	CreateDashboard(ctx context.Context, dashboard *shared.Dashboard) error
	GetDashboardByID(ctx context.Context, id uuid.UUID) (*shared.Dashboard, error)
	GetDashboardByCode(ctx context.Context, dashboardCode string) (*shared.Dashboard, error)
	GetDashboardsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.Dashboard, int64, error)
	UpdateDashboard(ctx context.Context, dashboard *shared.Dashboard) error
	DeleteDashboard(ctx context.Context, id uuid.UUID) error

	// Dashboard widget management
	CreateDashboardWidget(ctx context.Context, widget interface{}) error // DashboardWidget type not found
	GetDashboardWidgetsByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]interface{}, error) // DashboardWidget type not found
	UpdateDashboardWidget(ctx context.Context, widget interface{}) error // DashboardWidget type not found
	DeleteDashboardWidget(ctx context.Context, id uuid.UUID) error

	// Report scheduling
	CreateReportSchedule(ctx context.Context, schedule *shared.ReportSchedule) error
	GetReportScheduleByID(ctx context.Context, id uuid.UUID) (*shared.ReportSchedule, error)
	GetActiveSchedules(ctx context.Context, limit, offset int) ([]*shared.ReportSchedule, int64, error)
	UpdateReportSchedule(ctx context.Context, schedule *shared.ReportSchedule) error
	DeleteReportSchedule(ctx context.Context, id uuid.UUID) error
	ProcessScheduledReports(ctx context.Context) error

	// Claim reporting
	GenerateClaimReport(ctx context.Context, claimID uuid.UUID, reportType string) (*claim.ClaimReporting, error)
	GetClaimReportingByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimReporting, error)
	UpdateClaimReporting(ctx context.Context, claimReporting *claim.ClaimReporting) error

	// Actuarial data
	CreateActuarialData(ctx context.Context, data *shared.ActuarialData) error
	GetActuarialDataByPeriod(ctx context.Context, periodStart, periodEnd string, dataType string) ([]*shared.ActuarialData, error)
	UpdateActuarialData(ctx context.Context, data *shared.ActuarialData) error

	// Reporting statistics
	GetReportingStatistics(ctx context.Context, startDate, endDate string) (map[string]interface{}, error)
	GetReportGenerationTime(ctx context.Context, reportType string, startDate, endDate string) (int, error)
	GetDashboardUsageStatistics(ctx context.Context, dashboardID uuid.UUID, startDate, endDate string) (map[string]interface{}, error)

	// ======================================
	// MISSING IMPLEMENTED METHODS
	// ======================================

	// Real-time Analytics
	GetRealTimeAnalytics(ctx context.Context, userID uuid.UUID, metrics, timeRange string) (map[string]interface{}, error)
	GetLiveDashboard(ctx context.Context, userID uuid.UUID, dashboardType, refreshInterval string) (map[string]interface{}, error)
	StreamAnalytics(ctx context.Context, writer http.ResponseWriter, userID uuid.UUID, streamType, interval string) error

	// Predictive Analytics
	GetPredictiveClaimsAnalytics(ctx context.Context, userID uuid.UUID, timeframe, confidence string) (map[string]interface{}, error)
	GetPredictiveChurnAnalytics(ctx context.Context, userID uuid.UUID, segment, timeframe string) (map[string]interface{}, error)
	GetPredictivePremiumAnalytics(ctx context.Context, userID uuid.UUID, marketSegment, riskProfile string) (map[string]interface{}, error)

	// Custom Reports
	CreateCustomReport(ctx context.Context, userID uuid.UUID, name, description, reportType string, filters map[string]interface{}, columns []string, schedule *string, isPublic bool) (interface{}, error) // CustomReport type not found
	GetCustomReport(ctx context.Context, reportID, userID uuid.UUID) (interface{}, error) // CustomReport type not found
	UpdateCustomReport(ctx context.Context, reportID, userID uuid.UUID, name, description string, filters map[string]interface{}, columns []string, schedule *string, isPublic *bool) error
	DeleteCustomReport(ctx context.Context, reportID, userID uuid.UUID) error

	// Data Export
	ExportAnalyticsData(ctx context.Context, userID uuid.UUID, dataType string, filters map[string]interface{}, format string, dateRange map[string]string, includeRaw bool) (uuid.UUID, error)
	GetAnalyticsExports(ctx context.Context, userID uuid.UUID, limit, offset int) ([]interface{}, int, error) // AnalyticsExport type not found
	DownloadAnalyticsExport(ctx context.Context, exportID, userID uuid.UUID) ([]byte, string, string, error)

	// Business Intelligence
	GetBIDashboards(ctx context.Context, userID uuid.UUID, category, isPublic string) ([]interface{}, error) // BIDashboard type not found
	CreateBIDashboard(ctx context.Context, ownerID uuid.UUID, name, description, category string, isPublic bool, widgets []map[string]interface{}, layout map[string]interface{}, permissions map[string]interface{}) (interface{}, error) // BIDashboard type not found

	// Legacy methods for backward compatibility - commented out to avoid duplicates
	// GetReport(ctx context.Context, reportID uuid.UUID) (*shared.Report, error)
	// GetReports(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*shared.Report, int64, error)
	// UpdateReport(ctx context.Context, reportID uuid.UUID, updates map[string]interface{}) error
	// DeleteReport(ctx context.Context, reportID uuid.UUID) error
	GetDashboard(ctx context.Context, dashboardType string) (*shared.Dashboard, error)
	CreateScheduledReport(ctx context.Context, schedule *shared.ReportSchedule) error
	GetKPIs(ctx context.Context, category string, dateRange string) ([]*shared.KPI, error)
}
