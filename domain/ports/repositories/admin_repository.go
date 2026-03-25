package repositories

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// AdminRepository defines the interface for administrative operations
type AdminRepository interface {
	// Role management
	CreateRole(ctx context.Context, role *models.Role) error
	GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	ListRoles(ctx context.Context) ([]*models.Role, error)

	// Permission management
	CreatePermission(ctx context.Context, permission *models.Permission) error
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	GetPermissionByCode(ctx context.Context, code string) (*models.Permission, error)
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error
	ListPermissions(ctx context.Context) ([]*models.Permission, error)

	// User role assignments
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID, assignedBy uuid.UUID) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*models.UserRole, error)

	// Audit logging
	CreateAuditLog(ctx context.Context, log *models.AuditLog) error
	GetAuditLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*models.AuditLog, int64, error)

	// System configuration
	GetConfig(ctx context.Context, key string) (*models.SystemConfiguration, error)
	SetConfig(ctx context.Context, config *models.SystemConfiguration) error
	ListConfigs(ctx context.Context, category string) ([]*models.SystemConfiguration, error)

	// Department management
	CreateDepartment(ctx context.Context, dept *models.Department) error
	GetDepartmentByID(ctx context.Context, id uuid.UUID) (*models.Department, error)
	UpdateDepartment(ctx context.Context, dept *models.Department) error
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	ListDepartments(ctx context.Context) ([]*models.Department, error)

	// Team management
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByID(ctx context.Context, id uuid.UUID) (*models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	ListTeams(ctx context.Context, departmentID *uuid.UUID) ([]*models.Team, error)
}
