package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// WhiteLabelService defines the interface for white label tenant business logic operations
type WhiteLabelService interface {
	// Tenant organization management
	CreateTenantOrganization(ctx context.Context, tenant *shared.TenantOrganization) error
	GetTenantOrganizationByID(ctx context.Context, id uuid.UUID) (*shared.TenantOrganization, error)
	GetTenantOrganizationByTenantCode(ctx context.Context, tenantCode string) (*shared.TenantOrganization, error)
	UpdateTenantOrganization(ctx context.Context, tenant *shared.TenantOrganization) error
	SuspendTenantOrganization(ctx context.Context, tenantID uuid.UUID, reason string) error
	ActivateTenantOrganization(ctx context.Context, tenantID uuid.UUID) error
	ListTenantOrganizations(ctx context.Context, status string, limit, offset int) ([]*shared.TenantOrganization, int64, error)

	// Tenant configuration management
	CreateTenantConfiguration(ctx context.Context, config *shared.TenantConfiguration) error
	GetTenantConfigurationByID(ctx context.Context, id uuid.UUID) (*shared.TenantConfiguration, error)
	GetTenantConfigurations(ctx context.Context, tenantID uuid.UUID, category string) ([]*shared.TenantConfiguration, error)
	UpdateTenantConfiguration(ctx context.Context, config *shared.TenantConfiguration) error
	DeleteTenantConfiguration(ctx context.Context, id uuid.UUID) error

	// Tenant user management
	AddTenantUser(ctx context.Context, tenantUser *shared.TenantUser) error
	GetTenantUserByID(ctx context.Context, id uuid.UUID) (*shared.TenantUser, error)
	GetTenantUsers(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*shared.TenantUser, int64, error)
	UpdateTenantUser(ctx context.Context, tenantUser *shared.TenantUser) error
	RemoveTenantUser(ctx context.Context, id uuid.UUID) error

	// Partner integration management
	CreatePartnerIntegration(ctx context.Context, integration *shared.PartnerIntegration) error
	GetPartnerIntegrationByID(ctx context.Context, id uuid.UUID) (*shared.PartnerIntegration, error)
	GetPartnerIntegrationsByType(ctx context.Context, integrationType string, limit, offset int) ([]*shared.PartnerIntegration, int64, error)
	UpdatePartnerIntegration(ctx context.Context, integration *shared.PartnerIntegration) error
	DeletePartnerIntegration(ctx context.Context, id uuid.UUID) error
	TestPartnerIntegration(ctx context.Context, integrationID uuid.UUID) (bool, error)

	// API key management
	CreateTenantAPIKey(ctx context.Context, apiKey *shared.TenantAPIKey) error
	GetTenantAPIKeyByID(ctx context.Context, id uuid.UUID) (*shared.TenantAPIKey, error)
	GetTenantAPIKeys(ctx context.Context, tenantID uuid.UUID) ([]*shared.TenantAPIKey, error)
	UpdateTenantAPIKey(ctx context.Context, apiKey *shared.TenantAPIKey) error
	RevokeTenantAPIKey(ctx context.Context, id uuid.UUID) error
	ValidateAPIKey(ctx context.Context, apiKey string) (*shared.TenantAPIKey, error)

	// Tenant billing management
	CreateTenantBilling(ctx context.Context, billing *shared.TenantBilling) error
	GetTenantBillingByID(ctx context.Context, id uuid.UUID) (*shared.TenantBilling, error)
	GetTenantBillingHistory(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*shared.TenantBilling, int64, error)
	UpdateTenantBilling(ctx context.Context, billing *shared.TenantBilling) error
	ProcessTenantBilling(ctx context.Context, tenantID uuid.UUID, periodStart, periodEnd string) error

	// White label statistics
	GetWhiteLabelStatistics(ctx context.Context, tenantID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetTenantUsageStatistics(ctx context.Context, tenantID uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
}
