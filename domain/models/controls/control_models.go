package controls

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// === Control System Models ===
// These models allow dynamic configuration of controls without code changes

// ControlRule represents a configurable business rule
type ControlRule struct {
	ID               uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string                 `gorm:"uniqueIndex;not null" json:"name"`
	Description      string                 `json:"description"`
	EntityType       string                 `json:"entity_type"`    // Device, Policy, Claim, Payment, User
	OperationType    string                 `json:"operation_type"` // CREATE, UPDATE, DELETE, EXECUTE
	RuleType         string                 `json:"rule_type"`      // VALIDATION, LIMIT, APPROVAL, RESTRICTION
	Severity         string                 `json:"severity"`       // CRITICAL, HIGH, MEDIUM, LOW
	IsActive         bool                   `json:"is_active"`
	Priority         int                    `json:"priority"`                   // Lower number = higher priority
	Condition        string                 `gorm:"type:text" json:"condition"` // SQL/Expression condition
	Parameters       map[string]interface{} `gorm:"type:json" json:"parameters"`
	ErrorMessage     string                 `json:"error_message"`
	RequiresApproval bool                   `json:"requires_approval"`
	ApprovalLevels   int                    `json:"approval_levels"`
	CreatedBy        uuid.UUID              `json:"created_by"`
	UpdatedBy        uuid.UUID              `json:"updated_by"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	DeletedAt        gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlWorkflow represents a configurable approval workflow
type ControlWorkflow struct {
	ID                uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name              string                 `gorm:"uniqueIndex;not null" json:"name"`
	Description       string                 `json:"description"`
	EntityType        string                 `json:"entity_type"`
	TriggerCondition  string                 `gorm:"type:text" json:"trigger_condition"` // When to trigger workflow
	ApprovalSteps     []ControlApprovalStep  `gorm:"foreignKey:WorkflowID" json:"approval_steps"`
	AutoApproveRules  string                 `gorm:"type:text" json:"auto_approve_rules"` // Conditions for auto-approval
	TimeoutHours      int                    `json:"timeout_hours"`
	EscalationEnabled bool                   `json:"escalation_enabled"`
	EscalationAfter   int                    `json:"escalation_after_hours"`
	IsActive          bool                   `json:"is_active"`
	Parameters        map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy         uuid.UUID              `json:"created_by"`
	UpdatedBy         uuid.UUID              `json:"updated_by"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlApprovalStep represents a step in an approval workflow
type ControlApprovalStep struct {
	ID                   uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	WorkflowID           uuid.UUID              `json:"workflow_id"`
	StepOrder            int                    `json:"step_order"`
	Name                 string                 `json:"name"`
	ApproverType         string                 `json:"approver_type"`                       // ROLE, USER, GROUP, DYNAMIC
	ApproverValue        string                 `json:"approver_value"`                      // Role name, user ID, group ID, or expression
	RequiredApprovals    int                    `json:"required_approvals"`                  // Number of approvals needed at this step
	CanSkipCondition     string                 `gorm:"type:text" json:"can_skip_condition"` // Condition to skip this step
	NotificationTemplate string                 `json:"notification_template"`
	Parameters           map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

// ControlApprovalRequest represents an active approval request
type ControlApprovalRequest struct {
	ID              uuid.UUID                `gorm:"type:uuid;primaryKey" json:"id"`
	WorkflowID      uuid.UUID                `json:"workflow_id"`
	Workflow        *ControlWorkflow         `gorm:"foreignKey:WorkflowID" json:"workflow,omitempty"`
	EntityType      string                   `json:"entity_type"`
	EntityID        uuid.UUID                `json:"entity_id"`
	RequestorID     uuid.UUID                `json:"requestor_id"`
	CurrentStep     int                      `json:"current_step"`
	Status          string                   `json:"status"` // PENDING, APPROVED, REJECTED, EXPIRED
	RequestData     map[string]interface{}   `gorm:"type:json" json:"request_data"`
	ApprovalHistory []ControlApprovalHistory `gorm:"foreignKey:RequestID" json:"approval_history"`
	ExpiresAt       time.Time                `json:"expires_at"`
	CompletedAt     *time.Time               `json:"completed_at,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// ControlApprovalHistory tracks approval decisions
type ControlApprovalHistory struct {
	ID           uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	RequestID    uuid.UUID              `json:"request_id"`
	StepOrder    int                    `json:"step_order"`
	ApproverID   uuid.UUID              `json:"approver_id"`
	Decision     string                 `json:"decision"` // APPROVED, REJECTED
	Comments     string                 `json:"comments"`
	DecisionData map[string]interface{} `gorm:"type:json" json:"decision_data"`
	DecidedAt    time.Time              `json:"decided_at"`
}

// ControlIntegrityRule defines data integrity rules
type ControlIntegrityRule struct {
	ID                   uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name                 string                 `gorm:"uniqueIndex;not null" json:"name"`
	Description          string                 `json:"description"`
	EntityType           string                 `json:"entity_type"`
	IntegrityType        string                 `json:"integrity_type"` // REFERENTIAL, UNIQUE, CASCADE, DEPENDENCY
	ParentEntity         string                 `json:"parent_entity,omitempty"`
	ChildEntity          string                 `json:"child_entity,omitempty"`
	RelationshipType     string                 `json:"relationship_type"` // ONE_TO_ONE, ONE_TO_MANY, MANY_TO_MANY
	OnDeleteAction       string                 `json:"on_delete_action"`  // CASCADE, RESTRICT, SET_NULL, SOFT_DELETE
	OnUpdateAction       string                 `json:"on_update_action"`  // CASCADE, RESTRICT
	ValidationExpression string                 `gorm:"type:text" json:"validation_expression"`
	IsActive             bool                   `json:"is_active"`
	Parameters           map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy            uuid.UUID              `json:"created_by"`
	UpdatedBy            uuid.UUID              `json:"updated_by"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	DeletedAt            gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlAuditRule defines what and how to audit
type ControlAuditRule struct {
	ID                uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name              string                 `gorm:"uniqueIndex;not null" json:"name"`
	Description       string                 `json:"description"`
	EntityType        string                 `json:"entity_type"`
	OperationType     string                 `json:"operation_type"`
	AuditLevel        string                 `json:"audit_level"`                       // FULL, CHANGES_ONLY, SUMMARY
	SensitiveFields   []string               `gorm:"type:json" json:"sensitive_fields"` // Fields to encrypt/mask
	RetentionDays     int                    `json:"retention_days"`
	RequiresSignature bool                   `json:"requires_signature"`
	AlertOnAnomaly    bool                   `json:"alert_on_anomaly"`
	AnomalyRules      string                 `gorm:"type:text" json:"anomaly_rules"`
	IsActive          bool                   `json:"is_active"`
	Parameters        map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy         uuid.UUID              `json:"created_by"`
	UpdatedBy         uuid.UUID              `json:"updated_by"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlEncryptionRule defines field encryption rules
type ControlEncryptionRule struct {
	ID              uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string                 `gorm:"uniqueIndex;not null" json:"name"`
	EntityType      string                 `json:"entity_type"`
	FieldName       string                 `json:"field_name"`
	EncryptionType  string                 `json:"encryption_type"`  // AES256, RSA, HASH
	MaskingStrategy string                 `json:"masking_strategy"` // FULL, PARTIAL, NONE
	MaskingPattern  string                 `json:"masking_pattern"`  // e.g., "****-****-****-{last4}"
	IsActive        bool                   `json:"is_active"`
	Parameters      map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy       uuid.UUID              `json:"created_by"`
	UpdatedBy       uuid.UUID              `json:"updated_by"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlTransactionRule defines transaction control rules
type ControlTransactionRule struct {
	ID             uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string                 `gorm:"uniqueIndex;not null" json:"name"`
	Description    string                 `json:"description"`
	EntityType     string                 `json:"entity_type"`
	OperationType  string                 `json:"operation_type"`
	IsolationLevel string                 `json:"isolation_level"` // READ_UNCOMMITTED, READ_COMMITTED, REPEATABLE_READ, SERIALIZABLE
	MaxRetries     int                    `json:"max_retries"`
	RetryBackoff   string                 `json:"retry_backoff"` // LINEAR, EXPONENTIAL
	TimeoutSeconds int                    `json:"timeout_seconds"`
	RequiresLock   bool                   `json:"requires_lock"`
	LockType       string                 `json:"lock_type"` // PESSIMISTIC, OPTIMISTIC
	BatchSizeLimit int                    `json:"batch_size_limit"`
	IsActive       bool                   `json:"is_active"`
	Parameters     map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy      uuid.UUID              `json:"created_by"`
	UpdatedBy      uuid.UUID              `json:"updated_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlFieldValidation defines field-level validation rules
type ControlFieldValidation struct {
	ID             uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	EntityType     string                 `json:"entity_type"`
	FieldName      string                 `json:"field_name"`
	ValidationType string                 `json:"validation_type"`                  // REQUIRED, FORMAT, RANGE, CUSTOM
	ValidationRule string                 `gorm:"type:text" json:"validation_rule"` // Regex or expression
	MinValue       *float64               `json:"min_value,omitempty"`
	MaxValue       *float64               `json:"max_value,omitempty"`
	MinLength      *int                   `json:"min_length,omitempty"`
	MaxLength      *int                   `json:"max_length,omitempty"`
	AllowedValues  []string               `gorm:"type:json" json:"allowed_values,omitempty"`
	ErrorMessage   string                 `json:"error_message"`
	IsActive       bool                   `json:"is_active"`
	Parameters     map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy      uuid.UUID              `json:"created_by"`
	UpdatedBy      uuid.UUID              `json:"updated_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlLimit defines operational limits
type ControlLimit struct {
	ID            uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string                 `gorm:"uniqueIndex;not null" json:"name"`
	EntityType    string                 `json:"entity_type"`
	LimitType     string                 `json:"limit_type"` // COUNT, AMOUNT, RATE, TIME
	LimitValue    float64                `json:"limit_value"`
	TimeWindow    string                 `json:"time_window"` // DAILY, WEEKLY, MONTHLY, YEARLY
	Scope         string                 `json:"scope"`       // USER, DEVICE, ACCOUNT, GLOBAL
	ScopeID       *uuid.UUID             `json:"scope_id,omitempty"`
	ActionOnLimit string                 `json:"action_on_limit"` // REJECT, APPROVE_REQUIRED, ALERT
	IsActive      bool                   `json:"is_active"`
	Parameters    map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy     uuid.UUID              `json:"created_by"`
	UpdatedBy     uuid.UUID              `json:"updated_by"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	DeletedAt     gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlPermission defines access control rules
type ControlPermission struct {
	ID          uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string                 `gorm:"uniqueIndex;not null" json:"name"`
	EntityType  string                 `json:"entity_type"`
	Operation   string                 `json:"operation"` // CREATE, READ, UPDATE, DELETE, EXECUTE
	Role        string                 `json:"role"`
	Conditions  string                 `gorm:"type:text" json:"conditions"`   // Additional conditions
	FieldAccess map[string]string      `gorm:"type:json" json:"field_access"` // Field-level permissions
	IsActive    bool                   `json:"is_active"`
	Parameters  map[string]interface{} `gorm:"type:json" json:"parameters"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	UpdatedBy   uuid.UUID              `json:"updated_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// ControlTemplate defines reusable control configurations
type ControlTemplate struct {
	ID           uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string                 `gorm:"uniqueIndex;not null" json:"name"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"` // COMPLIANCE, SECURITY, BUSINESS, OPERATIONAL
	TemplateData map[string]interface{} `gorm:"type:json" json:"template_data"`
	IsPublic     bool                   `json:"is_public"`
	Version      string                 `json:"version"`
	CreatedBy    uuid.UUID              `json:"created_by"`
	UpdatedBy    uuid.UUID              `json:"updated_by"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    gorm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
}

// === Helper Methods ===

// IsApplicable checks if a control rule is applicable
func (cr *ControlRule) IsApplicable(entityType, operationType string) bool {
	return cr.IsActive &&
		cr.EntityType == entityType &&
		cr.OperationType == operationType
}

// NeedsApproval checks if workflow needs approval
func (cw *ControlWorkflow) NeedsApproval(data map[string]interface{}) bool {
	if !cw.IsActive {
		return false
	}
	// This would evaluate the trigger condition with the data
	// For now, return true if workflow is active
	return true
}

// IsExpired checks if approval request is expired
func (car *ControlApprovalRequest) IsExpired() bool {
	return time.Now().After(car.ExpiresAt)
}

// GetCurrentStep returns current approval step
func (car *ControlApprovalRequest) GetCurrentStep() *ControlApprovalStep {
	if car.Workflow == nil || len(car.Workflow.ApprovalSteps) == 0 {
		return nil
	}

	for _, step := range car.Workflow.ApprovalSteps {
		if step.StepOrder == car.CurrentStep {
			return &step
		}
	}

	return nil
}

// IsComplete checks if all approval steps are complete
func (car *ControlApprovalRequest) IsComplete() bool {
	return car.Status == "APPROVED" || car.Status == "REJECTED"
}

// GetRetentionDate returns when audit should be deleted
func (car *ControlAuditRule) GetRetentionDate() time.Time {
	return time.Now().AddDate(0, 0, car.RetentionDays)
}

// ShouldEncrypt checks if field should be encrypted
func (cer *ControlEncryptionRule) ShouldEncrypt(fieldName string) bool {
	return cer.IsActive && cer.FieldName == fieldName
}

// GetMaskedValue returns masked value based on pattern
func (cer *ControlEncryptionRule) GetMaskedValue(value string) string {
	if cer.MaskingStrategy == "NONE" {
		return value
	}

	if cer.MaskingStrategy == "FULL" {
		return "********"
	}

	// Partial masking - show last 4 characters
	if len(value) > 4 {
		return "****" + value[len(value)-4:]
	}

	return "****"
}

// IsWithinLimit checks if value is within limit
func (cl *ControlLimit) IsWithinLimit(currentValue float64) bool {
	return currentValue <= cl.LimitValue
}

// HasPermission checks if role has permission
func (cp *ControlPermission) HasPermission(role, operation string) bool {
	return cp.IsActive &&
		cp.Role == role &&
		cp.Operation == operation
}

// ValidateField validates a field value
func (cfv *ControlFieldValidation) ValidateField(value interface{}) error {
	if !cfv.IsActive {
		return nil
	}

	// Implement validation logic based on validation type
	// This is simplified - real implementation would be more complex

	return nil
}
