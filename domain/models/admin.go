package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/types"
)

// Role represents user roles in the system
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoleCode    string    `gorm:"uniqueIndex;not null" json:"role_code"`
	RoleName    string    `gorm:"not null" json:"role_name"`
	Description string    `json:"description"`

	// Role Details
	Type       string `json:"type"`  // system, custom
	Level      int    `json:"level"` // hierarchy level
	Department string `json:"department"`

	// Status
	IsActive     bool `json:"is_active"`
	IsDefault    bool `json:"is_default"`
	IsSuperAdmin bool `json:"is_super_admin"`

	// Permissions (JSON array of permission codes)
	Permissions types.JSONArray `gorm:"type:json" json:"permissions"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	UserRoles []UserRole `gorm:"foreignKey:RoleID" json:"user_roles,omitempty"`
}

// UserRole represents user-role assignments
type UserRole struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RoleID uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`

	// Assignment Details
	AssignedBy uuid.UUID  `gorm:"type:uuid" json:"assigned_by"`
	AssignedAt time.Time  `json:"assigned_at"`
	ExpiresAt  *time.Time `json:"expires_at"`

	// Status
	IsActive  bool `json:"is_active"`
	IsPrimary bool `json:"is_primary"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User     *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role     *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Assigner *User `gorm:"foreignKey:AssignedBy" json:"assigner,omitempty"`
}

// Permission represents system permissions
type Permission struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PermissionCode string    `gorm:"uniqueIndex;not null" json:"permission_code"`
	PermissionName string    `gorm:"not null" json:"permission_name"`
	Description    string    `json:"description"`

	// Permission Details
	Resource string `json:"resource"` // entity being accessed
	Action   string `json:"action"`   // create, read, update, delete, execute
	Scope    string `json:"scope"`    // own, team, department, all

	// Categorization
	Category string `json:"category"`
	Module   string `json:"module"`

	// Status
	IsActive    bool `json:"is_active"`
	IsCritical  bool `json:"is_critical"`
	RequiresMFA bool `json:"requires_mfa"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// AuditLog represents system audit logs
type AuditLog struct {
	ID     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID *uuid.UUID `gorm:"type:uuid" json:"user_id"`

	// Action Details
	Action     string `gorm:"not null" json:"action"`
	Resource   string `json:"resource"`
	ResourceID string `json:"resource_id"`
	Module     string `json:"module"`

	// Request Details
	Method    string `json:"method"` // GET, POST, PUT, DELETE
	Endpoint  string `json:"endpoint"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	SessionID string `json:"session_id"`

	// Data
	OldData string `gorm:"type:json" json:"old_data"`
	NewData string `gorm:"type:json" json:"new_data"`
	Changes string `gorm:"type:json" json:"changes"`

	// Response
	Status       string `json:"status"` // success, failure
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`

	// Performance
	ResponseTime int `json:"response_time"` // milliseconds

	// Metadata
	Tags           types.JSONArray `gorm:"type:json" json:"tags"`
	AdditionalData string          `gorm:"type:json" json:"additional_data"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// SystemConfiguration represents system settings
type SystemConfiguration struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ConfigKey   string    `gorm:"uniqueIndex;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	ConfigType  string    `json:"config_type"` // string, number, boolean, json

	// Configuration Details
	Category    string `json:"category"`
	Module      string `json:"module"`
	Description string `json:"description"`

	// Validation
	ValidationRules string          `gorm:"type:json" json:"validation_rules"`
	DefaultValue    string          `json:"default_value"`
	PossibleValues  types.JSONArray `gorm:"type:json" json:"possible_values"`

	// Security
	IsEncrypted     bool `json:"is_encrypted"`
	IsSensitive     bool `json:"is_sensitive"`
	RequiresRestart bool `json:"requires_restart"`

	// Change Management
	ModifiedBy       *uuid.UUID `gorm:"type:uuid" json:"modified_by"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovalRequired bool       `json:"approval_required"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Modifier *User `gorm:"foreignKey:ModifiedBy" json:"modifier,omitempty"`
	Approver *User `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

// Department represents organizational departments
type Department struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DepartmentCode string    `gorm:"uniqueIndex;not null" json:"department_code"`
	DepartmentName string    `gorm:"not null" json:"department_name"`
	Description    string    `json:"description"`

	// Hierarchy
	ParentDepartmentID *uuid.UUID `gorm:"type:uuid" json:"parent_department_id"`
	Level              int        `json:"level"`
	Path               string     `json:"path"` // hierarchical path

	// Management
	ManagerID *uuid.UUID `gorm:"type:uuid" json:"manager_id"`
	HeadCount int        `json:"head_count"`
	Budget    float64    `json:"budget"`

	// Contact
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`

	// Status
	IsActive bool `json:"is_active"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ParentDepartment *Department `gorm:"foreignKey:ParentDepartmentID" json:"parent_department,omitempty"`
	Manager          *User       `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Teams            []Team      `gorm:"foreignKey:DepartmentID" json:"teams,omitempty"`
}

// Team represents teams within departments
type Team struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TeamCode    string    `gorm:"uniqueIndex;not null" json:"team_code"`
	TeamName    string    `gorm:"not null" json:"team_name"`
	Description string    `json:"description"`

	// Organization
	DepartmentID uuid.UUID  `gorm:"type:uuid;not null" json:"department_id"`
	TeamLeadID   *uuid.UUID `gorm:"type:uuid" json:"team_lead_id"`

	// Team Details
	TeamSize       int    `json:"team_size"`
	TeamType       string `json:"team_type"` // operational, support, sales, technical
	Specialization string `json:"specialization"`

	// Performance
	TargetKPIs  string `gorm:"type:json" json:"target_kpis"`
	CurrentKPIs string `gorm:"type:json" json:"current_kpis"`

	// Status
	IsActive bool `json:"is_active"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Department  *Department  `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	TeamLead    *User        `gorm:"foreignKey:TeamLeadID" json:"team_lead,omitempty"`
	TeamMembers []TeamMember `gorm:"foreignKey:TeamID" json:"team_members,omitempty"`
}

// TeamMember represents team membership
type TeamMember struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TeamID uuid.UUID `gorm:"type:uuid;not null" json:"team_id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Membership Details
	Role     string     `json:"role"` // member, lead, supervisor
	JoinedAt time.Time  `json:"joined_at"`
	LeftAt   *time.Time `json:"left_at"`

	// Status
	IsActive  bool `json:"is_active"`
	IsPrimary bool `json:"is_primary"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Team *Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Task represents operational tasks
type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TaskNumber  string    `gorm:"uniqueIndex;not null" json:"task_number"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`

	// Task Details
	Type     string `json:"type"` // review, approval, processing, investigation
	Category string `json:"category"`
	Priority string `json:"priority"` // low, medium, high, urgent

	// Assignment
	AssignedTo *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	AssignedBy uuid.UUID  `gorm:"type:uuid" json:"assigned_by"`
	TeamID     *uuid.UUID `gorm:"type:uuid" json:"team_id"`

	// Related Entities
	RelatedEntity   string `json:"related_entity"` // policy, claim, user, etc
	RelatedEntityID string `json:"related_entity_id"`

	// Status & Progress
	Status      string     `json:"status"`   // pending, in_progress, completed, cancelled
	Progress    int        `json:"progress"` // 0-100
	CompletedAt *time.Time `json:"completed_at"`

	// Deadlines
	DueDate      *time.Time `json:"due_date"`
	ReminderDate *time.Time `json:"reminder_date"`

	// Metadata
	Tags        types.JSONArray `gorm:"type:json" json:"tags"`
	Attachments types.JSONArray `gorm:"type:json" json:"attachments"`
	Comments    string          `gorm:"type:json" json:"comments"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Assignee *User `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	Assigner *User `gorm:"foreignKey:AssignedBy" json:"assigner,omitempty"`
	Team     *Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

// Notification represents system notifications
type Notification struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Notification Details
	Type     string `json:"type"`     // info, warning, error, success
	Category string `json:"category"` // system, policy, claim, payment
	Title    string `json:"title"`
	Message  string `gorm:"type:text" json:"message"`

	// Delivery
	Channel  string `json:"channel"`  // in_app, email, sms, push
	Priority string `json:"priority"` // low, medium, high

	// Status
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	IsDelivered bool       `json:"is_delivered"`
	DeliveredAt *time.Time `json:"delivered_at"`

	// Action
	ActionRequired bool       `json:"action_required"`
	ActionURL      string     `json:"action_url"`
	ActionTaken    bool       `json:"action_taken"`
	ActionTakenAt  *time.Time `json:"action_taken_at"`

	// Metadata
	Data      string     `gorm:"type:json" json:"data"`
	ExpiresAt *time.Time `json:"expires_at"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
