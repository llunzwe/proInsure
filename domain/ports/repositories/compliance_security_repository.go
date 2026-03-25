package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/device"
)

// ComplianceSecurityRepository defines the interface for compliance and security persistence operations
type ComplianceSecurityRepository interface {
	// Device compliance operations
	CreateDeviceComplianceStatus(ctx context.Context, compliance *device.DeviceComplianceStatus) error
	GetDeviceComplianceStatusByDeviceID(ctx context.Context, deviceID uuid.UUID) (*device.DeviceComplianceStatus, error)
	UpdateDeviceComplianceStatus(ctx context.Context, compliance *device.DeviceComplianceStatus) error
	GetCompliantDevices(ctx context.Context, limit, offset int) ([]uuid.UUID, int64, error)
	GetNonCompliantDevices(ctx context.Context, limit, offset int) ([]uuid.UUID, int64, error)
	GetDevicesNeedingAssessment(ctx context.Context, assessmentDate time.Time, limit, offset int) ([]uuid.UUID, int64, error)

	// Device legal holds operations
	CreateDeviceLegalHold(ctx context.Context, hold *device.DeviceLegalHolds) error
	GetDeviceLegalHoldByDeviceID(ctx context.Context, deviceID uuid.UUID) (*device.DeviceLegalHolds, error)
	GetActiveLegalHolds(ctx context.Context, limit, offset int) ([]*device.DeviceLegalHolds, int64, error)
	UpdateDeviceLegalHold(ctx context.Context, hold *device.DeviceLegalHolds) error
	ReleaseLegalHold(ctx context.Context, deviceID uuid.UUID, releasedBy uuid.UUID, releaseReason string) error

	// Device security compliance operations
	CreateDeviceSecurityCompliance(ctx context.Context, compliance *device.DeviceSecurityCompliance) error
	GetDeviceSecurityComplianceByDeviceID(ctx context.Context, deviceID uuid.UUID) (*device.DeviceSecurityCompliance, error)
	UpdateDeviceSecurityCompliance(ctx context.Context, compliance *device.DeviceSecurityCompliance) error

	// Regulatory compliance operations
	GetRegulatoryComplianceStatus(ctx context.Context, deviceID uuid.UUID, region string) (map[string]interface{}, error)
	GetDevicesByComplianceRegion(ctx context.Context, region string, limit, offset int) ([]uuid.UUID, int64, error)
	GetDevicesWithViolations(ctx context.Context, severity string, limit, offset int) ([]uuid.UUID, int64, error)

	// Audit and certification operations
	RecordComplianceAudit(ctx context.Context, deviceID uuid.UUID, auditData map[string]interface{}) error
	GetComplianceAuditHistory(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]map[string]interface{}, int64, error)
	GetCertificationStatus(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)

	// Statistics
	GetComplianceStatistics(ctx context.Context, region *string, startDate, endDate time.Time) (map[string]interface{}, error)
	GetComplianceScoreByRegion(ctx context.Context, region string) (float64, error)
}
