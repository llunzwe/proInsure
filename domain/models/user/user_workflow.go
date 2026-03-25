package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserWorkflow manages approval workflows and process automation for users
type UserWorkflow struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Active Workflows
	ActiveWorkflows       []map[string]interface{} `gorm:"type:json" json:"active_workflows"`
	PendingApprovals      int                      `gorm:"default:0" json:"pending_approvals"`
	CompletedWorkflows    int                      `gorm:"default:0" json:"completed_workflows"`
	FailedWorkflows       int                      `gorm:"default:0" json:"failed_workflows"`
	AverageCompletionTime int                      `json:"average_completion_time_hours"`
	WorkflowEfficiency    float64                  `gorm:"default:0" json:"workflow_efficiency"`

	// Approval Requirements
	RequiresApproval     map[string]bool            `gorm:"type:json" json:"requires_approval"`
	ApprovalLevels       map[string]int             `gorm:"type:json" json:"approval_levels"`
	ApprovalThresholds   map[string]decimal.Decimal `gorm:"type:json" json:"approval_thresholds"`
	AutoApprovalEnabled  bool                       `gorm:"default:false" json:"auto_approval_enabled"`
	AutoApprovalCriteria map[string]interface{}     `gorm:"type:json" json:"auto_approval_criteria"`
	DelegatedApprovals   []uuid.UUID                `gorm:"type:json" json:"delegated_approvals"`
	EscalationRules      map[string]interface{}     `gorm:"type:json" json:"escalation_rules"`
	SLAConfiguration     map[string]int             `gorm:"type:json" json:"sla_configuration"`
}

// UserApprovalRequest represents a pending approval request
type UserApprovalRequest struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	WorkflowID  uuid.UUID  `gorm:"type:uuid;not null" json:"workflow_id"`
	RequesterID uuid.UUID  `gorm:"type:uuid;not null" json:"requester_id"`
	ApproverID  *uuid.UUID `gorm:"type:uuid" json:"approver_id"`
	DelegatedTo *uuid.UUID `gorm:"type:uuid" json:"delegated_to"`

	// Request Details
	RequestType          string                 `gorm:"type:varchar(50)" json:"request_type"`
	RequestStatus        string                 `gorm:"type:varchar(20)" json:"request_status"` // pending/approved/rejected/expired
	Priority             string                 `gorm:"type:varchar(20)" json:"priority"`
	Subject              string                 `gorm:"not null" json:"subject"`
	Description          string                 `gorm:"type:text" json:"description"`
	RequestData          map[string]interface{} `gorm:"type:json" json:"request_data"`
	Amount               decimal.Decimal        `gorm:"type:decimal(15,2)" json:"amount"`
	Currency             string                 `gorm:"type:varchar(3)" json:"currency"`
	ApprovalLevel        int                    `gorm:"default:1" json:"approval_level"`
	CurrentLevel         int                    `gorm:"default:1" json:"current_level"`
	RequiredApprovals    int                    `gorm:"default:1" json:"required_approvals"`
	ReceivedApprovals    int                    `gorm:"default:0" json:"received_approvals"`
	ExpiryDate           *time.Time             `json:"expiry_date"`
	SLADeadline          *time.Time             `json:"sla_deadline"`
	ResponseRequired     bool                   `gorm:"default:true" json:"response_required"`
	AutoApprovalEligible bool                   `gorm:"default:false" json:"auto_approval_eligible"`
	EscalationLevel      int                    `gorm:"default:0" json:"escalation_level"`
	EscalatedAt          *time.Time             `json:"escalated_at"`
}

// UserApprovalHistory stores approval decision history
type UserApprovalHistory struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID            uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ApprovalRequestID uuid.UUID `gorm:"type:uuid;not null" json:"approval_request_id"`
	ApproverID        uuid.UUID `gorm:"type:uuid;not null" json:"approver_id"`

	// Decision Details
	Decision         string                 `gorm:"type:varchar(20)" json:"decision"` // approved/rejected/delegated/escalated
	DecisionDate     time.Time              `json:"decision_date"`
	DecisionReason   string                 `gorm:"type:text" json:"decision_reason"`
	Comments         string                 `gorm:"type:text" json:"comments"`
	Conditions       []string               `gorm:"type:json" json:"conditions"`
	ResponseTime     int                    `json:"response_time_hours"`
	SLAMet           bool                   `gorm:"default:true" json:"sla_met"`
	DecisionMetadata map[string]interface{} `gorm:"type:json" json:"decision_metadata"`
	SystemGenerated  bool                   `gorm:"default:false" json:"system_generated"`
	OverrideApplied  bool                   `gorm:"default:false" json:"override_applied"`
	OverrideReason   string                 `gorm:"type:text" json:"override_reason"`
}

// UserProcessAutomation manages automated processes for users
type UserProcessAutomation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Automation Configuration
	AutomationEnabled    bool                   `gorm:"default:false" json:"automation_enabled"`
	AutomatedProcesses   []string               `gorm:"type:json" json:"automated_processes"`
	AutomationRules      map[string]interface{} `gorm:"type:json" json:"automation_rules"`
	Triggers             map[string]interface{} `gorm:"type:json" json:"triggers"`
	Actions              map[string]interface{} `gorm:"type:json" json:"actions"`
	Schedules            map[string]string      `gorm:"type:json" json:"schedules"`
	LastExecutionTime    map[string]time.Time   `gorm:"type:json" json:"last_execution_time"`
	NextScheduledRun     map[string]time.Time   `gorm:"type:json" json:"next_scheduled_run"`
	ExecutionCount       map[string]int         `gorm:"type:json" json:"execution_count"`
	SuccessRate          float64                `gorm:"default:0" json:"success_rate"`
	FailureReasons       []string               `gorm:"type:json" json:"failure_reasons"`
	AutoRetryEnabled     bool                   `gorm:"default:true" json:"auto_retry_enabled"`
	MaxRetryAttempts     int                    `gorm:"default:3" json:"max_retry_attempts"`
	NotificationSettings map[string]bool        `gorm:"type:json" json:"notification_settings"`
	ErrorHandling        map[string]string      `gorm:"type:json" json:"error_handling"`
	DependencyChain      []string               `gorm:"type:json" json:"dependency_chain"`
	ConditionalLogic     map[string]interface{} `gorm:"type:json" json:"conditional_logic"`
	DataTransformations  map[string]interface{} `gorm:"type:json" json:"data_transformations"`
	IntegrationEndpoints map[string]string      `gorm:"type:json" json:"integration_endpoints"`
}

// UserEscalation manages escalation paths and priority handling
type UserEscalation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	EscalatedTo      *uuid.UUID `gorm:"type:uuid" json:"escalated_to"`
	OriginalAssignee *uuid.UUID `gorm:"type:uuid" json:"original_assignee"`

	// Escalation Details
	EscalationType       string                   `gorm:"type:varchar(50)" json:"escalation_type"`
	EscalationReason     string                   `gorm:"type:text" json:"escalation_reason"`
	Priority             string                   `gorm:"type:varchar(20)" json:"priority"`
	Severity             string                   `gorm:"type:varchar(20)" json:"severity"`
	Status               string                   `gorm:"type:varchar(20)" json:"status"`
	EscalationLevel      int                      `gorm:"default:1" json:"escalation_level"`
	MaxEscalationLevel   int                      `gorm:"default:3" json:"max_escalation_level"`
	EscalationDate       time.Time                `json:"escalation_date"`
	ResolutionDeadline   *time.Time               `json:"resolution_deadline"`
	ResolutionDate       *time.Time               `json:"resolution_date"`
	ResolutionTime       int                      `json:"resolution_time_hours"`
	ResolutionDetails    string                   `gorm:"type:text" json:"resolution_details"`
	ImpactAssessment     string                   `gorm:"type:text" json:"impact_assessment"`
	StakeholdersNotified []uuid.UUID              `gorm:"type:json" json:"stakeholders_notified"`
	EscalationPath       []map[string]interface{} `gorm:"type:json" json:"escalation_path"`
	PreventiveMeasures   []string                 `gorm:"type:json" json:"preventive_measures"`
	RecurrencePrevention string                   `gorm:"type:text" json:"recurrence_prevention"`
}

// UserSLA manages service level agreements and compliance
type UserSLA struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// SLA Configuration
	SLATier               string                   `gorm:"type:varchar(20)" json:"sla_tier"` // bronze/silver/gold/platinum
	ResponseTimeTarget    map[string]int           `gorm:"type:json" json:"response_time_target_hours"`
	ResolutionTimeTarget  map[string]int           `gorm:"type:json" json:"resolution_time_target_hours"`
	AvailabilityTarget    float64                  `gorm:"default:99.9" json:"availability_target"`
	SLAMetrics            map[string]float64       `gorm:"type:json" json:"sla_metrics"`
	CurrentCompliance     float64                  `gorm:"default:100" json:"current_compliance"`
	ComplianceHistory     []map[string]interface{} `gorm:"type:json" json:"compliance_history"`
	Breaches              int                      `gorm:"default:0" json:"breaches"`
	LastBreachDate        *time.Time               `json:"last_breach_date"`
	BreachDetails         []map[string]interface{} `gorm:"type:json" json:"breach_details"`
	Credits               decimal.Decimal          `gorm:"type:decimal(10,2);default:0" json:"credits"`
	CreditHistory         []map[string]interface{} `gorm:"type:json" json:"credit_history"`
	ExclusionPeriods      []map[string]time.Time   `gorm:"type:json" json:"exclusion_periods"`
	CustomTerms           map[string]interface{}   `gorm:"type:json" json:"custom_terms"`
	ReviewDate            *time.Time               `json:"review_date"`
	RenegotiationEligible bool                     `gorm:"default:false" json:"renegotiation_eligible"`
}

// UserTaskManagement manages tasks and activities for users
type UserTaskManagement struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	AssignedBy *uuid.UUID `gorm:"type:uuid" json:"assigned_by"`
	AssignedTo *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`

	// Task Details
	TaskID             string                   `gorm:"type:varchar(50);uniqueIndex" json:"task_id"`
	TaskType           string                   `gorm:"type:varchar(50)" json:"task_type"`
	TaskCategory       string                   `gorm:"type:varchar(50)" json:"task_category"`
	Title              string                   `gorm:"not null" json:"title"`
	Description        string                   `gorm:"type:text" json:"description"`
	Status             string                   `gorm:"type:varchar(20)" json:"status"` // pending/in_progress/completed/cancelled
	Priority           string                   `gorm:"type:varchar(20)" json:"priority"`
	DueDate            *time.Time               `json:"due_date"`
	StartDate          *time.Time               `json:"start_date"`
	CompletionDate     *time.Time               `json:"completion_date"`
	EstimatedHours     float64                  `gorm:"default:0" json:"estimated_hours"`
	ActualHours        float64                  `gorm:"default:0" json:"actual_hours"`
	Progress           float64                  `gorm:"default:0" json:"progress"`
	Dependencies       []string                 `gorm:"type:json" json:"dependencies"`
	Subtasks           []map[string]interface{} `gorm:"type:json" json:"subtasks"`
	Attachments        []string                 `gorm:"type:json" json:"attachments"`
	Comments           []map[string]interface{} `gorm:"type:json" json:"comments"`
	Tags               []string                 `gorm:"type:json" json:"tags"`
	Recurring          bool                     `gorm:"default:false" json:"recurring"`
	RecurrencePattern  string                   `gorm:"type:varchar(50)" json:"recurrence_pattern"`
	NextOccurrence     *time.Time               `json:"next_occurrence"`
	CompletionCriteria []string                 `gorm:"type:json" json:"completion_criteria"`
	QualityScore       float64                  `gorm:"default:0" json:"quality_score"`
	AutomationEligible bool                     `gorm:"default:false" json:"automation_eligible"`
}

// TableName returns the table name
func (UserWorkflow) TableName() string {
	return "user_workflows"
}

// TableName returns the table name
func (UserApprovalRequest) TableName() string {
	return "user_approval_requests"
}

// TableName returns the table name
func (UserApprovalHistory) TableName() string {
	return "user_approval_history"
}

// TableName returns the table name
func (UserProcessAutomation) TableName() string {
	return "user_process_automation"
}

// TableName returns the table name
func (UserEscalation) TableName() string {
	return "user_escalations"
}

// TableName returns the table name
func (UserSLA) TableName() string {
	return "user_sla"
}

// TableName returns the table name
func (UserTaskManagement) TableName() string {
	return "user_task_management"
}
