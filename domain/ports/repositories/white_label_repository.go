package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// WhiteLabelRepository defines the interface for white label tenant persistence operations
type WhiteLabelRepository interface {
	// Tenant organization operations
	CreateTenantOrganization(ctx context.Context, tenant *shared.TenantOrganization) error
	GetTenantOrganizationByID(ctx context.Context, id uuid.UUID) (*shared.TenantOrganization, error)
	GetTenantOrganizationByTenantCode(ctx context.Context, tenantCode string) (*shared.TenantOrganization, error)
	GetTenantOrganizationByBusinessLicense(ctx context.Context, businessLicense string) (*shared.TenantOrganization, error)
	UpdateTenantOrganization(ctx context.Context, tenant *shared.TenantOrganization) error
	DeleteTenantOrganization(ctx context.Context, id uuid.UUID) error
	ListTenantOrganizations(ctx context.Context, status string, limit, offset int) ([]*shared.TenantOrganization, int64, error)
	SearchTenantOrganizations(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*shared.TenantOrganization, int64, error)

	// Tenant configuration operations
	CreateTenantConfiguration(ctx context.Context, config *shared.TenantConfiguration) error
	GetTenantConfigurationByID(ctx context.Context, id uuid.UUID) (*shared.TenantConfiguration, error)
	GetTenantConfigurationsByTenantID(ctx context.Context, tenantID uuid.UUID, category string) ([]*shared.TenantConfiguration, error)
	UpdateTenantConfiguration(ctx context.Context, config *shared.TenantConfiguration) error
	DeleteTenantConfiguration(ctx context.Context, id uuid.UUID) error

	// Tenant user operations
	CreateTenantUser(ctx context.Context, tenantUser *shared.TenantUser) error
	GetTenantUserByID(ctx context.Context, id uuid.UUID) (*shared.TenantUser, error)
	GetTenantUsersByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*shared.TenantUser, int64, error)
	GetTenantUserByUserID(ctx context.Context, userID uuid.UUID) (*shared.TenantUser, error)
	UpdateTenantUser(ctx context.Context, tenantUser *shared.TenantUser) error
	DeleteTenantUser(ctx context.Context, id uuid.UUID) error

	// Partner integration operations
	CreatePartnerIntegration(ctx context.Context, integration *shared.PartnerIntegration) error
	GetPartnerIntegrationByID(ctx context.Context, id uuid.UUID) (*shared.PartnerIntegration, error)
	GetPartnerIntegrationByPartnerCode(ctx context.Context, partnerCode string) (*shared.PartnerIntegration, error)
	GetPartnerIntegrationsByType(ctx context.Context, integrationType string, limit, offset int) ([]*shared.PartnerIntegration, int64, error)
	GetActivePartnerIntegrations(ctx context.Context, limit, offset int) ([]*shared.PartnerIntegration, int64, error)
	UpdatePartnerIntegration(ctx context.Context, integration *shared.PartnerIntegration) error
	DeletePartnerIntegration(ctx context.Context, id uuid.UUID) error

	// API key operations
	CreateTenantAPIKey(ctx context.Context, apiKey *shared.TenantAPIKey) error
	GetTenantAPIKeyByID(ctx context.Context, id uuid.UUID) (*shared.TenantAPIKey, error)
	GetTenantAPIKeyByKey(ctx context.Context, key string) (*shared.TenantAPIKey, error)
	GetTenantAPIKeysByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*shared.TenantAPIKey, error)
	UpdateTenantAPIKey(ctx context.Context, apiKey *shared.TenantAPIKey) error
	RevokeTenantAPIKey(ctx context.Context, id uuid.UUID) error

	// Tenant billing operations
	CreateTenantBilling(ctx context.Context, billing *shared.TenantBilling) error
	GetTenantBillingByID(ctx context.Context, id uuid.UUID) (*shared.TenantBilling, error)
	GetTenantBillingByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*shared.TenantBilling, int64, error)
	UpdateTenantBilling(ctx context.Context, billing *shared.TenantBilling) error

	// Statistics
	GetWhiteLabelStatistics(ctx context.Context, tenantID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetTenantUsageStatistics(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
}
