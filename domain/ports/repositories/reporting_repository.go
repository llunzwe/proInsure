package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/claim"
	"smartsure/internal/domain/models/shared"
)

// ReportingRepository defines the interface for reporting persistence operations
type ReportingRepository interface {
	// Report operations
	CreateReport(ctx context.Context, report *shared.Report) error
	GetReportByID(ctx context.Context, id uuid.UUID) (*shared.Report, error)
	GetReportByReportNumber(ctx context.Context, reportNumber string) (*shared.Report, error)
	UpdateReport(ctx context.Context, report *shared.Report) error
	DeleteReport(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetReportsByType(ctx context.Context, reportType string, limit, offset int) ([]*shared.Report, int64, error)
	GetReportsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.Report, int64, error)
	GetReportsByPeriod(ctx context.Context, periodStart, periodEnd time.Time, limit, offset int) ([]*shared.Report, int64, error)
	GetReportsByGeneratedBy(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.Report, int64, error)
	GetScheduledReports(ctx context.Context, limit, offset int) ([]*shared.Report, int64, error)
	GetFailedReports(ctx context.Context, limit, offset int) ([]*shared.Report, int64, error)
	SearchReports(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*shared.Report, int64, error)

	// Dashboard operations
	CreateDashboard(ctx context.Context, dashboard *shared.Dashboard) error
	GetDashboardByID(ctx context.Context, id uuid.UUID) (*shared.Dashboard, error)
	GetDashboardByCode(ctx context.Context, dashboardCode string) (*shared.Dashboard, error)
	GetDashboardsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.Dashboard, int64, error)
	GetPublicDashboards(ctx context.Context, limit, offset int) ([]*shared.Dashboard, int64, error)
	UpdateDashboard(ctx context.Context, dashboard *shared.Dashboard) error
	DeleteDashboard(ctx context.Context, id uuid.UUID) error

	// Dashboard widget operations - commented out as DashboardWidget type doesn't exist
	// CreateDashboardWidget(ctx context.Context, widget *shared.DashboardWidget) error
	// GetDashboardWidgetByID(ctx context.Context, id uuid.UUID) (*shared.DashboardWidget, error)
	// GetDashboardWidgetsByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]*shared.DashboardWidget, error)
	// UpdateDashboardWidget(ctx context.Context, widget *shared.DashboardWidget) error
	DeleteDashboardWidget(ctx context.Context, id uuid.UUID) error

	// Report schedule operations
	CreateReportSchedule(ctx context.Context, schedule *shared.ReportSchedule) error
	GetReportScheduleByID(ctx context.Context, id uuid.UUID) (*shared.ReportSchedule, error)
	GetReportSchedulesByReportID(ctx context.Context, reportID uuid.UUID) ([]*shared.ReportSchedule, error)
	GetActiveSchedules(ctx context.Context, limit, offset int) ([]*shared.ReportSchedule, int64, error)
	UpdateReportSchedule(ctx context.Context, schedule *shared.ReportSchedule) error
	DeleteReportSchedule(ctx context.Context, id uuid.UUID) error

	// Claim reporting operations
	CreateClaimReporting(ctx context.Context, claimReporting *claim.ClaimReporting) error
	GetClaimReportingByID(ctx context.Context, id uuid.UUID) (*claim.ClaimReporting, error)
	GetClaimReportingByClaimID(ctx context.Context, claimID uuid.UUID) (*claim.ClaimReporting, error)
	GetClaimReportingByReportNumber(ctx context.Context, reportNumber string) (*claim.ClaimReporting, error)
	UpdateClaimReporting(ctx context.Context, claimReporting *claim.ClaimReporting) error

	// Claim report schedule operations
	CreateClaimReportSchedule(ctx context.Context, schedule *claim.ClaimReportSchedule) error
	GetClaimReportScheduleByID(ctx context.Context, id uuid.UUID) (*claim.ClaimReportSchedule, error)
	GetActiveClaimReportSchedules(ctx context.Context, limit, offset int) ([]*claim.ClaimReportSchedule, int64, error)
	UpdateClaimReportSchedule(ctx context.Context, schedule *claim.ClaimReportSchedule) error

	// Actuarial data operations
	CreateActuarialData(ctx context.Context, data *shared.ActuarialData) error
	GetActuarialDataByID(ctx context.Context, id uuid.UUID) (*shared.ActuarialData, error)
	GetActuarialDataByPeriod(ctx context.Context, periodStart, periodEnd time.Time, dataType string) ([]*shared.ActuarialData, error)
	UpdateActuarialData(ctx context.Context, data *shared.ActuarialData) error
	DeleteActuarialData(ctx context.Context, id uuid.UUID) error

	// Statistics
	GetReportingStatistics(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)
	GetReportGenerationTime(ctx context.Context, reportType string, startDate, endDate time.Time) (time.Duration, error)
	GetDashboardUsageStatistics(ctx context.Context, dashboardID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
}
