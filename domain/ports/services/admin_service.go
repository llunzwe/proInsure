package services

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// AdminService defines the interface for administrative business logic
type AdminService interface {
	// Role management
	CreateRole(ctx context.Context, role *models.Role) error
	GetRole(ctx context.Context, id uuid.UUID) (*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	ListRoles(ctx context.Context) ([]*models.Role, error)

	// Permission management
	CreatePermission(ctx context.Context, permission *models.Permission) error
	GetPermission(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error
	ListPermissions(ctx context.Context) ([]*models.Permission, error)

	// User role assignments
	AssignRole(ctx context.Context, userID, roleID, assignedBy uuid.UUID) error
	RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*models.UserRole, error)
	HasPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error)

	// Audit logging
	LogAction(ctx context.Context, action, resource, resourceID string, userID *uuid.UUID, metadata map[string]interface{}) error
	GetAuditLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*models.AuditLog, int64, error)

	// System configuration
	GetConfiguration(ctx context.Context, key string) (string, error)
	SetConfiguration(ctx context.Context, key, value string, userID uuid.UUID) error
	ListConfigurations(ctx context.Context, category string) ([]*models.SystemConfiguration, error)

	// Department management
	CreateDepartment(ctx context.Context, dept *models.Department) error
	GetDepartment(ctx context.Context, id uuid.UUID) (*models.Department, error)
	UpdateDepartment(ctx context.Context, dept *models.Department) error
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	ListDepartments(ctx context.Context) ([]*models.Department, error)

	// Team management
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, id uuid.UUID) (*models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	ListTeams(ctx context.Context, departmentID *uuid.UUID) ([]*models.Team, error)
}
