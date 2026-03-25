package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceAutomationRule represents automated rules and triggers
type DeviceAutomationRule struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RuleID         string    `gorm:"uniqueIndex;not null" json:"rule_id"`

	// Rule Definition
	RuleName       string    `gorm:"not null" json:"rule_name"`
	Description    string    `gorm:"type:text" json:"description"`
	RuleType       string    `gorm:"type:varchar(30);not null" json:"rule_type"` // event_based, time_based, condition_based, composite

	// Status & Control
	Enabled        bool      `gorm:"default:true;index" json:"enabled"`
	Priority       string    `gorm:"type:varchar(20);default:'normal'" json:"priority"` // low, normal, high, critical
	Status         string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active, paused, disabled, error

	// Trigger Conditions
	TriggerType    string    `gorm:"type:varchar(30);not null" json:"trigger_type"` // sensor, time, event, composite
	TriggerConditions string `gorm:"type:json;not null" json:"trigger_conditions"` // JSON trigger logic
	TriggerParameters string `gorm:"type:json" json:"trigger_parameters,omitempty"` // JSON trigger parameters

	// Action Definition
	Actions        string    `gorm:"type:json;not null" json:"actions"` // JSON action definitions
	ActionParameters string  `gorm:"type:json" json:"action_parameters,omitempty"` // JSON action parameters

	// Scheduling (for time-based rules)
	ScheduleType   string    `gorm:"type:varchar(20)" json:"schedule_type"` // cron, interval, one_time
	ScheduleConfig string    `gorm:"type:json" json:"schedule_config,omitempty"` // JSON schedule configuration
	NextExecution  *time.Time `json:"next_execution,omitempty"`

	// Execution Control
	MaxExecutions  int       `gorm:"default:-1" json:"max_executions"` // -1 for unlimited
	ExecutionCount int       `gorm:"default:0" json:"execution_count"`
	MaxFrequency   int       `gorm:"default:60" json:"max_frequency"` // max executions per hour
	CooldownPeriod int       `gorm:"default:0" json:"cooldown_period"` // seconds between executions

	// Error Handling
	ErrorHandling  string    `gorm:"type:json" json:"error_handling,omitempty"` // JSON error handling rules
	RetryPolicy    string    `gorm:"type:json" json:"retry_policy,omitempty"` // JSON retry configuration
	CircuitBreaker string    `gorm:"type:json" json:"circuit_breaker,omitempty"` // JSON circuit breaker config

	// Performance & Monitoring
	ExecutionTime  int       `json:"execution_time,omitempty"` // average execution time in ms
	SuccessRate    float64   `json:"success_rate,omitempty"` // 0-100
	LastExecution  *time.Time `json:"last_execution,omitempty"`
	LastSuccess    *time.Time `json:"last_success,omitempty"`
	LastFailure    *time.Time `json:"last_failure,omitempty"`

	// Audit & Compliance
	CreatedBy      uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	RequiresApproval bool    `gorm:"default:false" json:"requires_approval"`
	ApprovalStatus string    `gorm:"type:varchar(20)" json:"approval_status"` // pending, approved, rejected
	ApprovedBy     *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	Executions     []DeviceRuleExecution  `gorm:"foreignKey:RuleID;references:RuleID" json:"executions,omitempty"`
}

// DeviceWorkflowInstance represents running workflow instances
type DeviceWorkflowInstance struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	WorkflowID     string    `gorm:"uniqueIndex;not null" json:"workflow_id"`
	TemplateID     string    `gorm:"index" json:"template_id,omitempty"` // reference to workflow template

	// Workflow Details
	WorkflowName   string    `gorm:"not null" json:"workflow_name"`
	Description    string    `gorm:"type:text" json:"description"`
	WorkflowType   string    `gorm:"type:varchar(30);not null" json:"workflow_type"` // maintenance, incident, upgrade, custom

	// Status & Progress
	Status         string    `gorm:"type:varchar(20);not null;default:'running'" json:"status"` // running, paused, completed, failed, cancelled
	Progress       float64   `gorm:"default:0" json:"progress"` // 0-100
	CurrentStep    string    `json:"current_step,omitempty"`
	NextStep       string    `json:"next_step,omitempty"`

	// Timing
	StartedAt      time.Time `gorm:"not null" json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	EstimatedDuration int    `json:"estimated_duration,omitempty"` // minutes
	ActualDuration int       `json:"actual_duration,omitempty"` // minutes

	// Execution Context
	ExecutionContext string  `gorm:"type:json" json:"execution_context,omitempty"` // JSON execution context
	InputParameters  string  `gorm:"type:json" json:"input_parameters,omitempty"` // JSON input parameters
	OutputResults    string  `gorm:"type:json" json:"output_results,omitempty"` // JSON output results

	// Step Execution
	WorkflowSteps   string  `gorm:"type:json;not null" json:"workflow_steps"` // JSON workflow definition
	ExecutedSteps   string  `gorm:"type:json" json:"executed_steps,omitempty"` // JSON execution history
	PendingSteps    string  `gorm:"type:json" json:"pending_steps,omitempty"` // JSON remaining steps

	// Error Handling
	ErrorCount      int     `gorm:"default:0" json:"error_count"`
	LastError       string  `json:"last_error,omitempty"`
	ErrorDetails    string  `gorm:"type:json" json:"error_details,omitempty"` // JSON error information

	// Rollback & Recovery
	RollbackPlan    string  `gorm:"type:json" json:"rollback_plan,omitempty"` // JSON rollback procedures
	RollbackExecuted bool   `gorm:"default:false" json:"rollback_executed"`
	RecoveryActions string  `gorm:"type:json" json:"recovery_actions,omitempty"` // JSON recovery steps taken

	// Performance Metrics
	StepExecutionTimes string `gorm:"type:json" json:"step_execution_times,omitempty"` // JSON timing data
	ResourceUsage    string  `gorm:"type:json" json:"resource_usage,omitempty"` // JSON resource consumption

	// Audit
	InitiatedBy     uuid.UUID `gorm:"type:uuid;not null" json:"initiated_by"`
	Priority        string    `gorm:"type:varchar(20);default:'normal'" json:"priority"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	StepExecutions []DeviceWorkflowStep   `gorm:"foreignKey:WorkflowID;references:WorkflowID" json:"step_executions,omitempty"`
}

// DeviceWorkflowStep represents individual workflow step executions
type DeviceWorkflowStep struct {
	database.BaseModel

	WorkflowID     string    `gorm:"not null;index" json:"workflow_id"`
	StepID         string    `gorm:"not null" json:"step_id"`
	StepName       string    `gorm:"not null" json:"step_name"`
	StepType       string    `gorm:"type:varchar(30);not null" json:"step_type"` // action, condition, parallel, subworkflow

	// Execution Details
	Status         string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, running, completed, failed, skipped
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	Duration       int       `json:"duration,omitempty"` // seconds

	// Step Configuration
	StepConfig     string    `gorm:"type:json" json:"step_config,omitempty"` // JSON step configuration
	InputData      string    `gorm:"type:json" json:"input_data,omitempty"` // JSON input data
	OutputData     string    `gorm:"type:json" json:"output_data,omitempty"` // JSON output data

	// Dependencies & Flow Control
	Dependencies   string    `gorm:"type:json" json:"dependencies,omitempty"` // JSON step dependencies
	Condition      string    `gorm:"type:json" json:"condition,omitempty"` // JSON execution condition
	ParallelGroup  string    `json:"parallel_group,omitempty"` // group for parallel execution

	// Error Handling
	ErrorHandling  string    `gorm:"type:json" json:"error_handling,omitempty"` // JSON error handling
	RetryCount     int       `gorm:"default:0" json:"retry_count"`
	MaxRetries     int       `gorm:"default:3" json:"max_retries"`

	// Assignment & Responsibility
	AssignedTo     *uuid.UUID `gorm:"type:uuid" json:"assigned_to,omitempty"`
	AssignedRole   string    `json:"assigned_role,omitempty"`

	// Relationships
	Workflow       DeviceWorkflowInstance `gorm:"foreignKey:WorkflowID;references:WorkflowID" json:"workflow,omitempty"`
}

// DeviceRuleExecution represents automation rule execution records
type DeviceRuleExecution struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RuleID         string    `gorm:"not null;index" json:"rule_id"`
	ExecutionID    string    `gorm:"uniqueIndex;not null" json:"execution_id"`

	// Execution Details
	ExecutedAt     time.Time `gorm:"not null;index" json:"executed_at"`
	Duration       int       `json:"duration,omitempty"` // milliseconds
	Status         string    `gorm:"type:varchar(20);not null" json:"status"` // success, failed, partial, timeout

	// Trigger Information
	TriggerEvent   string    `gorm:"type:json" json:"trigger_event,omitempty"` // JSON trigger event data
	TriggerConditions string `gorm:"type:json" json:"trigger_conditions"` // JSON conditions that triggered execution

	// Execution Results
	ActionsExecuted string   `gorm:"type:json" json:"actions_executed,omitempty"` // JSON actions performed
	ActionResults   string   `gorm:"type:json" json:"action_results,omitempty"` // JSON action results
	OutputData      string   `gorm:"type:json" json:"output_data,omitempty"` // JSON execution output

	// Error Information
	ErrorMessage    string   `json:"error_message,omitempty"`
	ErrorDetails    string   `gorm:"type:json" json:"error_details,omitempty"` // JSON error information
	PartialSuccess  bool     `gorm:"default:false" json:"partial_success"`

	// Performance Metrics
	ResourceUsage   string   `gorm:"type:json" json:"resource_usage,omitempty"` // JSON resource consumption
	PerformanceData string   `gorm:"type:json" json:"performance_data,omitempty"` // JSON performance metrics

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	Rule           DeviceAutomationRule `gorm:"foreignKey:RuleID;references:RuleID" json:"rule,omitempty"`
}

// DeviceSmartAction represents intelligent automated actions
type DeviceSmartAction struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ActionID       string    `gorm:"uniqueIndex;not null" json:"action_id"`

	// Action Definition
	ActionName     string    `gorm:"not null" json:"action_name"`
	ActionType     string    `gorm:"type:varchar(30);not null" json:"action_type"` // preventive, corrective, optimization, security
	Description    string    `gorm:"type:text" json:"description"`

	// AI-Driven Intelligence
	AIModelUsed    string    `json:"ai_model_used,omitempty"`
	ConfidenceScore float64  `json:"confidence_score,omitempty"` // 0-100
	DecisionLogic  string    `gorm:"type:json" json:"decision_logic,omitempty"` // JSON decision criteria

	// Action Configuration
	ActionConfig   string    `gorm:"type:json;not null" json:"action_config"` // JSON action configuration
	Prerequisites  string    `gorm:"type:json" json:"prerequisites,omitempty"` // JSON requirements
	ResourceRequirements string `gorm:"type:json" json:"resource_requirements,omitempty"` // JSON resource needs

	// Execution Control
	AutoExecute     bool      `gorm:"default:false" json:"auto_execute"`
	RequiresApproval bool     `gorm:"default:true" json:"requires_approval"`
	ApprovalTimeout int       `gorm:"default:3600" json:"approval_timeout"` // seconds
	EscalationTime  int       `gorm:"default:7200" json:"escalation_time"` // seconds

	// Status & Scheduling
	Status         string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, executing, completed, failed
	ScheduledFor   *time.Time `json:"scheduled_for,omitempty"`
	ExecutedAt     *time.Time `json:"executed_at,omitempty"`

	// Execution Results
	ExecutionResult string   `gorm:"type:json" json:"execution_result,omitempty"` // JSON execution results
	SuccessMetrics  string   `gorm:"type:json" json:"success_metrics,omitempty"` // JSON success measurements
	ImpactAssessment string  `gorm:"type:json" json:"impact_assessment,omitempty"` // JSON impact analysis

	// Learning & Adaptation
	Effectiveness   float64  `json:"effectiveness,omitempty"` // 0-100, measured effectiveness
	AdaptationData  string   `gorm:"type:json" json:"adaptation_data,omitempty"` // JSON learning data for model improvement

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceConditionalAction represents actions triggered by conditions
type DeviceConditionalAction struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ActionID       string    `gorm:"uniqueIndex;not null" json:"action_id"`

	// Condition Definition
	ConditionName  string    `gorm:"not null" json:"condition_name"`
	ConditionLogic string    `gorm:"type:json;not null" json:"condition_logic"` // JSON condition logic
	ConditionType  string    `gorm:"type:varchar(30);not null" json:"condition_type"` // threshold, pattern, anomaly, composite

	// Action Definition
	ActionDefinition string  `gorm:"type:json;not null" json:"action_definition"` // JSON action to execute
	ActionType      string   `gorm:"type:varchar(30);not null" json:"action_type"` // notification, command, workflow, alert

	// Trigger Control
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	Priority       string    `gorm:"type:varchar(20);default:'normal'" json:"priority"`
	CooldownPeriod int       `gorm:"default:300" json:"cooldown_period"` // seconds between triggers

	// Monitoring
	TriggerCount   int       `gorm:"default:0" json:"trigger_count"`
	LastTriggered  *time.Time `json:"last_triggered,omitempty"`
	SuccessRate    float64   `json:"success_rate,omitempty"` // 0-100

	// Performance Tuning
	FalsePositiveRate float64 `json:"false_positive_rate,omitempty"`
	TuningParameters string   `gorm:"type:json" json:"tuning_parameters,omitempty"` // JSON tuning parameters

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceTriggerAction represents event-driven trigger actions
type DeviceTriggerAction struct {
	database.BaseModel

	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	ActionID       string    `gorm:"uniqueIndex;not null" json:"action_id"`

	// Trigger Definition
	TriggerName    string    `gorm:"not null" json:"trigger_name"`
	TriggerType    string    `gorm:"type:varchar(30);not null" json:"trigger_type"` // event, sensor, time, external
	TriggerConfig  string    `gorm:"type:json;not null" json:"trigger_config"` // JSON trigger configuration

	// Action Chain
	ActionChain    string    `gorm:"type:json;not null" json:"action_chain"` // JSON sequence of actions
	FallbackActions string   `gorm:"type:json" json:"fallback_actions,omitempty"` // JSON fallback actions

	// Execution Control
	MaxConcurrency int       `gorm:"default:1" json:"max_concurrency"` // max concurrent executions
	ExecutionTimeout int     `gorm:"default:300" json:"execution_timeout"` // seconds
	RetryPolicy    string    `gorm:"type:json" json:"retry_policy,omitempty"` // JSON retry configuration

	// Monitoring & Analytics
	ExecutionHistory string  `gorm:"type:json" json:"execution_history,omitempty"` // JSON execution history
	PerformanceMetrics string `gorm:"type:json" json:"performance_metrics,omitempty"` // JSON performance data

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceWorkflowTemplate represents reusable workflow templates
type DeviceWorkflowTemplate struct {
	database.BaseModel

	TemplateID     string    `gorm:"uniqueIndex;not null" json:"template_id"`
	Name           string    `gorm:"not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	Category       string    `gorm:"type:varchar(30);not null" json:"category"` // maintenance, incident, upgrade, custom

	// Template Definition
	TemplateDefinition string `gorm:"type:json;not null" json:"template_definition"` // JSON workflow template
	DefaultParameters  string `gorm:"type:json" json:"default_parameters,omitempty"` // JSON default parameters
	ParameterSchema    string `gorm:"type:json" json:"parameter_schema,omitempty"` // JSON parameter validation schema

	// Usage Tracking
	UsageCount      int       `gorm:"default:0" json:"usage_count"`
	LastUsed        *time.Time `json:"last_used,omitempty"`
	SuccessRate     float64   `json:"success_rate,omitempty"` // 0-100
	AverageDuration int       `json:"average_duration,omitempty"` // minutes

	// Version Control
	Version         string    `gorm:"default:'1.0.0'" json:"version"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	DeprecatedAt    *time.Time `json:"deprecated_at,omitempty"`

	// Audit
	CreatedBy       uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	IsPublic        bool      `gorm:"default:false" json:"is_public"`
}

// TableName returns the table name for DeviceAutomationRule
func (DeviceAutomationRule) TableName() string {
	return "device_automation_rules"
}

// TableName returns the table name for DeviceWorkflowInstance
func (DeviceWorkflowInstance) TableName() string {
	return "device_workflow_instances"
}

// TableName returns the table name for DeviceWorkflowStep
func (DeviceWorkflowStep) TableName() string {
	return "device_workflow_steps"
}

// TableName returns the table name for DeviceRuleExecution
func (DeviceRuleExecution) TableName() string {
	return "device_rule_executions"
}

// TableName returns the table name for DeviceSmartAction
func (DeviceSmartAction) TableName() string {
	return "device_smart_actions"
}

// TableName returns the table name for DeviceConditionalAction
func (DeviceConditionalAction) TableName() string {
	return "device_conditional_actions"
}

// TableName returns the table name for DeviceTriggerAction
func (DeviceTriggerAction) TableName() string {
	return "device_trigger_actions"
}

// TableName returns the table name for DeviceWorkflowTemplate
func (DeviceWorkflowTemplate) TableName() string {
	return "device_workflow_templates"
}
