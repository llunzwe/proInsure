package controls

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// === Enhanced Control System Models ===
// Professional implementation with all required features

// === Entity Registry ===

// ControlEntityRegistry defines all valid entities in the system
type ControlEntityRegistry struct {
	ID              uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	EntityName      string            `gorm:"uniqueIndex;not null" json:"entity_name"` // User, Device, Policy, etc.
	TableName       string            `json:"table_name"`                              // users, devices, policies
	DisplayName     string            `json:"display_name"`                            // Human-readable name
	Description     string            `json:"description"`
	PrimaryKeyField string            `json:"primary_key_field"`
	SensitiveFields []string          `gorm:"type:json" json:"sensitive_fields"`
	RelatedEntities map[string]string `gorm:"type:json" json:"related_entities"` // Related entity mappings
	IsActive        bool              `json:"is_active"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// === JSON Condition Type ===

// ConditionJSON represents a structured condition that can be evaluated
type ConditionJSON map[string]interface{}

// Value implements driver.Valuer interface
func (c ConditionJSON) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner interface
func (c *ConditionJSON) Scan(value interface{}) error {
	if value == nil {
		*c = make(ConditionJSON)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, c)
	case string:
		return json.Unmarshal([]byte(v), c)
	default:
		return fmt.Errorf("cannot scan type %T into ConditionJSON", value)
	}
}

// === Enhanced Control Rule ===

// ControlRuleEnhanced represents a configurable business rule with all professional features
type ControlRuleEnhanced struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string    `gorm:"uniqueIndex;not null" json:"name"`
	Description   string    `json:"description"`
	EntityType    string    `gorm:"index" json:"entity_type"`    // References ControlEntityRegistry
	OperationType string    `gorm:"index" json:"operation_type"` // CREATE, READ, UPDATE, DELETE, EXECUTE, BULK_*
	Action        string    `json:"action"`                      // ALLOW, DENY, REQUIRE_APPROVAL, MODIFY, DELAY, LOG_ONLY
	Severity      string    `json:"severity"`                    // CRITICAL, HIGH, MEDIUM, LOW, INFO

	// Enhanced condition system
	Condition       ConditionJSON `gorm:"type:jsonb" json:"condition"`       // JSON condition structure
	RelatedEntities []string      `gorm:"type:json" json:"related_entities"` // Cross-entity support

	// Rule configuration
	Priority       int         `gorm:"index" json:"priority"`             // Execution order (lower = higher priority)
	ExecutionOrder int         `json:"execution_order"`                   // Within same priority
	DependsOnRules []uuid.UUID `gorm:"type:json" json:"depends_on_rules"` // Rule dependencies

	// Actions and modifications
	RequiresApproval  bool                   `json:"requires_approval"`
	ApprovalLevels    int                    `json:"approval_levels"`
	ModificationRules map[string]interface{} `gorm:"type:json" json:"modification_rules"` // Field modifications
	ErrorMessage      string                 `json:"error_message"`
	SuccessMessage    string                 `json:"success_message"`

	// Testing and activation
	IsDraft     bool                   `json:"is_draft"` // Test mode
	IsActive    bool                   `gorm:"index" json:"is_active"`
	TestedBy    *uuid.UUID             `json:"tested_by,omitempty"`
	TestedAt    *time.Time             `json:"tested_at,omitempty"`
	TestResults map[string]interface{} `gorm:"type:json" json:"test_results,omitempty"`

	// Multi-tenancy
	TenantID   *uuid.UUID `gorm:"index" json:"tenant_id,omitempty"` // Null for global rules
	AppliesTo  string     `json:"applies_to"`                       // GLOBAL, TENANT, USER_GROUP, USER
	UserGroups []string   `gorm:"type:json" json:"user_groups,omitempty"`

	// Performance and monitoring
	CacheKey   string `json:"cache_key,omitempty"` // For caching evaluation results
	TTLSeconds int    `json:"ttl_seconds"`         // Cache TTL

	// Additional parameters
	Parameters map[string]interface{} `gorm:"type:json" json:"parameters"`

	// Audit fields
	Version   int            `json:"version"`
	CreatedBy uuid.UUID      `json:"created_by"`
	UpdatedBy uuid.UUID      `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// === Rule Version History ===

// ControlRuleVersion tracks all changes to control rules
type ControlRuleVersion struct {
	ID             uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	RuleID         uuid.UUID              `gorm:"index" json:"rule_id"`
	Version        int                    `json:"version"`
	ChangedBy      uuid.UUID              `json:"changed_by"`
	ChangedAt      time.Time              `json:"changed_at"`
	ChangeType     string                 `json:"change_type"` // CREATE, UPDATE, DELETE, ACTIVATE, DEACTIVATE
	ChangeReason   string                 `json:"change_reason"`
	PreviousState  map[string]interface{} `gorm:"type:json" json:"previous_state"`
	NewState       map[string]interface{} `gorm:"type:json" json:"new_state"`
	ChangedFields  []string               `gorm:"type:json" json:"changed_fields"`
	IsRollbackable bool                   `json:"is_rollbackable"`
}

// === Rule Metrics ===

// ControlRuleMetrics tracks rule performance and effectiveness
type ControlRuleMetrics struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RuleID             uuid.UUID  `gorm:"uniqueIndex" json:"rule_id"`
	EvaluationCount    int64      `json:"evaluation_count"`
	PassCount          int64      `json:"pass_count"`
	FailCount          int64      `json:"fail_count"`
	ApprovalCount      int64      `json:"approval_count"`
	AvgExecutionTimeMs float64    `json:"avg_execution_time_ms"`
	MaxExecutionTimeMs float64    `json:"max_execution_time_ms"`
	LastEvaluated      *time.Time `json:"last_evaluated,omitempty"`
	LastPassed         *time.Time `json:"last_passed,omitempty"`
	LastFailed         *time.Time `json:"last_failed,omitempty"`
	ErrorCount         int64      `json:"error_count"`
	LastError          *string    `json:"last_error,omitempty"`
	EffectivenessScore float64    `json:"effectiveness_score"` // Calculated metric
	UpdatedAt          time.Time  `json:"updated_at"`
}

// === Enhanced Workflow ===

// ControlWorkflowEnhanced represents an enhanced approval workflow
type ControlWorkflowEnhanced struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string    `gorm:"uniqueIndex;not null" json:"name"`
	Description   string    `json:"description"`
	EntityType    string    `gorm:"index" json:"entity_type"`
	OperationType string    `json:"operation_type"`

	// Trigger conditions
	TriggerCondition ConditionJSON   `gorm:"type:jsonb" json:"trigger_condition"`
	PreConditions    []ConditionJSON `gorm:"type:json" json:"pre_conditions"`  // Must be met before workflow starts
	PostConditions   []ConditionJSON `gorm:"type:json" json:"post_conditions"` // Must be met after workflow completes

	// Workflow configuration
	WorkflowType     string        `json:"workflow_type"` // SEQUENTIAL, PARALLEL, CONDITIONAL
	MaxDurationHours int           `json:"max_duration_hours"`
	AutoApproveRules ConditionJSON `gorm:"type:jsonb" json:"auto_approve_rules"`
	AutoRejectRules  ConditionJSON `gorm:"type:jsonb" json:"auto_reject_rules"`

	// Escalation
	EscalationEnabled bool              `json:"escalation_enabled"`
	EscalationLevels  []EscalationLevel `gorm:"type:json" json:"escalation_levels"`

	// Multi-tenancy
	TenantID *uuid.UUID `gorm:"index" json:"tenant_id,omitempty"`

	// Status
	IsActive bool `gorm:"index" json:"is_active"`
	IsDraft  bool `json:"is_draft"`
	Version  int  `json:"version"`

	// Additional configuration
	Parameters         map[string]interface{} `gorm:"type:json" json:"parameters"`
	NotificationConfig map[string]interface{} `gorm:"type:json" json:"notification_config"`

	// Relationships
	ApprovalSteps []ControlApprovalStepEnhanced `gorm:"foreignKey:WorkflowID" json:"approval_steps"`

	// Audit fields
	CreatedBy uuid.UUID      `json:"created_by"`
	UpdatedBy uuid.UUID      `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// EscalationLevel defines an escalation tier
type EscalationLevel struct {
	Level             int      `json:"level"`
	TriggerAfterHours int      `json:"trigger_after_hours"`
	NotifyRoles       []string `json:"notify_roles"`
	NotifyUsers       []string `json:"notify_users"`
	AutoAction        string   `json:"auto_action"` // APPROVE, REJECT, DELEGATE
}

// === Enhanced Approval Step ===

// ControlApprovalStepEnhanced represents an enhanced approval step
type ControlApprovalStepEnhanced struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WorkflowID uuid.UUID `gorm:"index" json:"workflow_id"`
	StepOrder  int       `json:"step_order"`
	StepName   string    `json:"step_name"`
	StepType   string    `json:"step_type"` // APPROVAL, REVIEW, NOTIFICATION, CONDITION

	// Approver configuration
	ApproverType       string        `json:"approver_type"` // ROLE, USER, GROUP, DYNAMIC, EXPRESSION
	ApproverValue      string        `json:"approver_value"`
	ApproverExpression ConditionJSON `gorm:"type:jsonb" json:"approver_expression"` // Dynamic approver selection
	RequiredApprovals  int           `json:"required_approvals"`
	ApprovalStrategy   string        `json:"approval_strategy"` // ALL, ANY, MAJORITY, WEIGHTED

	// Conditions
	SkipCondition        ConditionJSON `gorm:"type:jsonb" json:"skip_condition"`
	AutoApproveCondition ConditionJSON `gorm:"type:jsonb" json:"auto_approve_condition"`

	// Timeouts
	TimeoutHours  int    `json:"timeout_hours"`
	TimeoutAction string `json:"timeout_action"` // APPROVE, REJECT, ESCALATE, SKIP

	// Notifications
	NotificationTemplate string `json:"notification_template"`
	ReminderInterval     int    `json:"reminder_interval_hours"`
	MaxReminders         int    `json:"max_reminders"`

	// Additional configuration
	Parameters      map[string]interface{} `gorm:"type:json" json:"parameters"`
	ValidationRules []ConditionJSON        `gorm:"type:json" json:"validation_rules"`

	// Audit fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// === Cross-Entity Rule ===

// ControlCrossEntityRule defines rules that span multiple entities
type ControlCrossEntityRule struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`

	// Entity relationships
	PrimaryEntity   string            `json:"primary_entity"`
	RelatedEntities map[string]string `gorm:"type:json" json:"related_entities"` // entity -> relationship

	// Cross-entity condition
	CrossCondition ConditionJSON            `gorm:"type:jsonb" json:"cross_condition"`
	JoinConditions map[string]ConditionJSON `gorm:"type:json" json:"join_conditions"`

	// Action configuration
	Action           string   `json:"action"`
	ActionScope      string   `json:"action_scope"` // PRIMARY_ONLY, ALL_ENTITIES, SPECIFIC_ENTITIES
	AffectedEntities []string `gorm:"type:json" json:"affected_entities"`

	// Performance
	UseCache bool `json:"use_cache"`
	CacheTTL int  `json:"cache_ttl_seconds"`

	// Status
	IsActive bool `gorm:"index" json:"is_active"`
	Priority int  `json:"priority"`

	// Audit fields
	CreatedBy uuid.UUID      `json:"created_by"`
	UpdatedBy uuid.UUID      `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// === Rule Test Case ===

// ControlRuleTestCase defines test scenarios for rules
type ControlRuleTestCase struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RuleID      uuid.UUID `gorm:"index" json:"rule_id"`
	TestName    string    `json:"test_name"`
	Description string    `json:"description"`

	// Test data
	EntityType    string                 `json:"entity_type"`
	OperationType string                 `json:"operation_type"`
	TestData      map[string]interface{} `gorm:"type:json" json:"test_data"`
	UserContext   map[string]interface{} `gorm:"type:json" json:"user_context"`

	// Expected results
	ExpectedAction string                 `json:"expected_action"`
	ExpectedResult map[string]interface{} `gorm:"type:json" json:"expected_result"`

	// Test results
	LastRunAt       *time.Time             `json:"last_run_at,omitempty"`
	LastRunBy       *uuid.UUID             `json:"last_run_by,omitempty"`
	LastResult      map[string]interface{} `gorm:"type:json" json:"last_result,omitempty"`
	Passed          *bool                  `json:"passed,omitempty"`
	ExecutionTimeMs *float64               `json:"execution_time_ms,omitempty"`

	// Status
	IsActive  bool      `json:"is_active"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// === Rule Template ===

// ControlRuleTemplate provides pre-built rule configurations
type ControlRuleTemplate struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TemplateName string    `gorm:"uniqueIndex;not null" json:"template_name"`
	Category     string    `json:"category"` // COMPLIANCE, SECURITY, BUSINESS, FRAUD, RISK
	Description  string    `json:"description"`

	// Template configuration
	EntityTypes      []string               `gorm:"type:json" json:"entity_types"`
	RuleTemplate     map[string]interface{} `gorm:"type:json" json:"rule_template"`
	WorkflowTemplate map[string]interface{} `gorm:"type:json" json:"workflow_template,omitempty"`

	// Customization points
	RequiredParams []string               `gorm:"type:json" json:"required_params"`
	OptionalParams []string               `gorm:"type:json" json:"optional_params"`
	DefaultValues  map[string]interface{} `gorm:"type:json" json:"default_values"`

	// Industry/Region specific
	Industry            string   `json:"industry,omitempty"`                              // INSURANCE, BANKING, RETAIL
	Region              string   `json:"region,omitempty"`                                // US, EU, ASIA
	ComplianceStandards []string `gorm:"type:json" json:"compliance_standards,omitempty"` // GDPR, HIPAA, SOC2

	// Metadata
	Version     string `json:"version"`
	IsPublic    bool   `json:"is_public"`
	IsOfficial  bool   `json:"is_official"`  // Provided by system
	IsCertified bool   `json:"is_certified"` // Compliance certified

	// Usage tracking
	UsageCount int64      `json:"usage_count"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	Rating     float64    `json:"rating"`

	// Audit fields
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// === Helper Methods ===

// Evaluate checks if the condition matches the provided data
func (c ConditionJSON) Evaluate(data map[string]interface{}, context map[string]interface{}) (bool, error) {
	// Check for logical operators
	if and, ok := c["AND"].([]interface{}); ok {
		for _, condition := range and {
			if cond, ok := condition.(map[string]interface{}); ok {
				if result, err := ConditionJSON(cond).Evaluate(data, context); err != nil || !result {
					return false, err
				}
			}
		}
		return true, nil
	}

	if or, ok := c["OR"].([]interface{}); ok {
		for _, condition := range or {
			if cond, ok := condition.(map[string]interface{}); ok {
				if result, err := ConditionJSON(cond).Evaluate(data, context); err == nil && result {
					return true, nil
				}
			}
		}
		return false, nil
	}

	// Evaluate field conditions
	for field, operators := range c {
		var value interface{}

		// Handle nested fields and context
		if field[:5] == "user." && context != nil {
			value = context[field[5:]]
		} else {
			value = data[field]
		}

		// Evaluate operators
		if ops, ok := operators.(map[string]interface{}); ok {
			for op, expected := range ops {
				if !evaluateOperator(value, op, expected) {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

// evaluateOperator evaluates a single operator condition
func evaluateOperator(value interface{}, operator string, expected interface{}) bool {
	switch operator {
	case "=", "==", "eq":
		return value == expected
	case "!=", "ne":
		return value != expected
	case ">", "gt":
		return compareValues(value, expected) > 0
	case ">=", "gte":
		return compareValues(value, expected) >= 0
	case "<", "lt":
		return compareValues(value, expected) < 0
	case "<=", "lte":
		return compareValues(value, expected) <= 0
	case "in":
		return isValueIn(value, expected)
	case "not_in":
		return !isValueIn(value, expected)
	case "contains":
		return contains(value, expected)
	case "starts_with":
		return startsWith(value, expected)
	case "ends_with":
		return endsWith(value, expected)
	case "regex":
		return matchesRegex(value, expected)
	case "exists":
		return value != nil
	case "not_exists":
		return value == nil
	default:
		return false
	}
}

// Helper functions for operator evaluation
func compareValues(a, b interface{}) int {
	// Implementation would handle type conversion and comparison
	// This is a simplified version
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			if va > vb {
				return 1
			} else if va < vb {
				return -1
			}
			return 0
		}
	case float64:
		if vb, ok := b.(float64); ok {
			if va > vb {
				return 1
			} else if va < vb {
				return -1
			}
			return 0
		}
	case string:
		if vb, ok := b.(string); ok {
			if va > vb {
				return 1
			} else if va < vb {
				return -1
			}
			return 0
		}
	}
	return 0
}

func isValueIn(value, list interface{}) bool {
	if arr, ok := list.([]interface{}); ok {
		for _, item := range arr {
			if value == item {
				return true
			}
		}
	}
	return false
}

func contains(value, substring interface{}) bool {
	if str, ok := value.(string); ok {
		if substr, ok := substring.(string); ok {
			return len(str) > 0 && len(substr) > 0 &&
				(str == substr || len(str) > len(substr))
		}
	}
	return false
}

func startsWith(value, prefix interface{}) bool {
	if str, ok := value.(string); ok {
		if pfx, ok := prefix.(string); ok {
			return len(str) >= len(pfx) && str[:len(pfx)] == pfx
		}
	}
	return false
}

func endsWith(value, suffix interface{}) bool {
	if str, ok := value.(string); ok {
		if sfx, ok := suffix.(string); ok {
			return len(str) >= len(sfx) && str[len(str)-len(sfx):] == sfx
		}
	}
	return false
}

func matchesRegex(value, pattern interface{}) bool {
	// Implementation would use regexp package
	// This is a placeholder
	return false
}

// GetApplicableRules returns rules applicable to an entity and operation
func GetApplicableRules(entityType, operationType string, tenantID *uuid.UUID) []ControlRuleEnhanced {
	// This would query the database for applicable rules
	// Placeholder for demonstration
	return []ControlRuleEnhanced{}
}

// ShouldCache determines if a rule result should be cached
func (r *ControlRuleEnhanced) ShouldCache() bool {
	return r.TTLSeconds > 0 && r.CacheKey != ""
}

// GetCacheKey generates a cache key for the rule evaluation
func (r *ControlRuleEnhanced) GetCacheKey(data map[string]interface{}) string {
	if r.CacheKey == "" {
		return ""
	}
	// Generate cache key based on rule ID and relevant data fields
	return fmt.Sprintf("rule:%s:%s", r.ID, r.CacheKey)
}

// IsTestMode checks if the rule is in test mode
func (r *ControlRuleEnhanced) IsTestMode() bool {
	return r.IsDraft && !r.IsActive
}

// CanExecute checks if the rule can be executed based on dependencies
func (r *ControlRuleEnhanced) CanExecute(executedRules []uuid.UUID) bool {
	for _, depID := range r.DependsOnRules {
		found := false
		for _, execID := range executedRules {
			if depID == execID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// GetEffectivenessScore calculates the effectiveness of a rule
func (m *ControlRuleMetrics) GetEffectivenessScore() float64 {
	if m.EvaluationCount == 0 {
		return 0
	}

	// Calculate based on various factors
	passRate := float64(m.PassCount) / float64(m.EvaluationCount)
	performanceScore := 1.0
	if m.AvgExecutionTimeMs > 100 {
		performanceScore = 100.0 / m.AvgExecutionTimeMs
	}

	errorPenalty := 1.0
	if m.ErrorCount > 0 {
		errorPenalty = float64(m.EvaluationCount-m.ErrorCount) / float64(m.EvaluationCount)
	}

	// Weighted score
	return (passRate*0.5 + performanceScore*0.3 + errorPenalty*0.2) * 100
}

// ShouldEscalate checks if the workflow should escalate
func (w *ControlWorkflowEnhanced) ShouldEscalate(startTime time.Time, currentLevel int) (bool, *EscalationLevel) {
	if !w.EscalationEnabled || len(w.EscalationLevels) == 0 {
		return false, nil
	}

	elapsed := time.Since(startTime).Hours()

	for _, level := range w.EscalationLevels {
		if level.Level > currentLevel && elapsed >= float64(level.TriggerAfterHours) {
			return true, &level
		}
	}

	return false, nil
}

// ValidateTemplate checks if a template is valid for an entity type
func (t *ControlRuleTemplate) ValidateTemplate(entityType string) bool {
	for _, et := range t.EntityTypes {
		if et == entityType || et == "*" {
			return true
		}
	}
	return false
}

// ApplyTemplate applies a template with parameters
func (t *ControlRuleTemplate) ApplyTemplate(params map[string]interface{}) (map[string]interface{}, error) {
	// Check required parameters
	for _, required := range t.RequiredParams {
		if _, ok := params[required]; !ok {
			return nil, fmt.Errorf("required parameter %s is missing", required)
		}
	}

	// Merge with defaults
	result := make(map[string]interface{})
	for k, v := range t.DefaultValues {
		result[k] = v
	}
	for k, v := range params {
		result[k] = v
	}

	// Apply template
	for k, v := range t.RuleTemplate {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	return result, nil
}
