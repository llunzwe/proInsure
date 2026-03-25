package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/device"
)

// ComplianceSecurityService defines the interface for compliance and security business logic operations
type ComplianceSecurityService interface {
	// Device compliance operations
	CreateDeviceComplianceStatus(ctx context.Context, compliance *device.DeviceComplianceStatus) error
	GetDeviceComplianceStatus(ctx context.Context, deviceID uuid.UUID) (*device.DeviceComplianceStatus, error)
	UpdateDeviceComplianceStatus(ctx context.Context, compliance *device.DeviceComplianceStatus) error
	AssessDeviceCompliance(ctx context.Context, deviceID uuid.UUID) (*device.DeviceComplianceStatus, error)

	// Compliance queries
	GetCompliantDevices(ctx context.Context, limit, offset int) ([]uuid.UUID, int64, error)
	GetNonCompliantDevices(ctx context.Context, limit, offset int) ([]uuid.UUID, int64, error)
	GetDevicesNeedingAssessment(ctx context.Context, assessmentDate string, limit, offset int) ([]uuid.UUID, int64, error)

	// Device legal hold operations
	CreateLegalHold(ctx context.Context, hold *device.DeviceLegalHolds) error
	GetLegalHoldByDeviceID(ctx context.Context, deviceID uuid.UUID) (*device.DeviceLegalHolds, error)
	GetActiveLegalHolds(ctx context.Context, limit, offset int) ([]*device.DeviceLegalHolds, int64, error)
	UpdateLegalHold(ctx context.Context, hold *device.DeviceLegalHolds) error
	ReleaseLegalHold(ctx context.Context, deviceID uuid.UUID, releasedBy uuid.UUID, releaseReason string) error

	// Device security compliance operations
	CreateDeviceSecurityCompliance(ctx context.Context, compliance *device.DeviceSecurityCompliance) error
	GetDeviceSecurityCompliance(ctx context.Context, deviceID uuid.UUID) (*device.DeviceSecurityCompliance, error)
	UpdateDeviceSecurityCompliance(ctx context.Context, compliance *device.DeviceSecurityCompliance) error

	// Regulatory compliance operations
	GetRegulatoryComplianceStatus(ctx context.Context, deviceID uuid.UUID, region string) (map[string]interface{}, error)
	GetDevicesByComplianceRegion(ctx context.Context, region string, limit, offset int) ([]uuid.UUID, int64, error)
	GetDevicesWithViolations(ctx context.Context, severity string, limit, offset int) ([]uuid.UUID, int64, error)

	// Audit and certification operations
	RecordComplianceAudit(ctx context.Context, deviceID uuid.UUID, auditData map[string]interface{}) error
	GetComplianceAuditHistory(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]map[string]interface{}, int64, error)
	GetCertificationStatus(ctx context.Context, deviceID uuid.UUID) (map[string]interface{}, error)

	// Compliance statistics
	GetComplianceStatistics(ctx context.Context, region *string, startDate, endDate string) (map[string]interface{}, error)
	GetComplianceScoreByRegion(ctx context.Context, region string) (float64, error)
}
