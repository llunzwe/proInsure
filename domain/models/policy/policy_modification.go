package policy

import (
	"encoding/json"
	"errors"
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================
// POLICY MODIFICATION MODELS
// ============================================

// ModificationType represents the type of modification
type ModificationType string

const (
	ModificationTypeCoverageIncrease  ModificationType = "coverage_increase"
	ModificationTypeCoverageDecrease  ModificationType = "coverage_decrease"
	ModificationTypeDeductibleChange  ModificationType = "deductible_change"
	ModificationTypeAddRider          ModificationType = "add_rider"
	ModificationTypeRemoveRider       ModificationType = "remove_rider"
	ModificationTypeAddEndorsement    ModificationType = "add_endorsement"
	ModificationTypeRemoveEndorsement ModificationType = "remove_endorsement"
	ModificationTypePaymentChange     ModificationType = "payment_change"
	ModificationTypeBeneficiaryChange ModificationType = "beneficiary_change"
	ModificationTypeAddressChange     ModificationType = "address_change"
	ModificationTypeDeviceChange      ModificationType = "device_change"
	ModificationTypeOwnershipTransfer ModificationType = "ownership_transfer"
	ModificationTypeTermExtension     ModificationType = "term_extension"
	ModificationTypeTermReduction     ModificationType = "term_reduction"
	ModificationTypePremiumAdjustment ModificationType = "premium_adjustment"
	ModificationTypeDiscountAdd       ModificationType = "discount_add"
	ModificationTypeDiscountRemove    ModificationType = "discount_remove"
	ModificationTypeSuspension        ModificationType = "suspension"
	ModificationTypeReinstatement     ModificationType = "reinstatement"
	ModificationTypeCancellation      ModificationType = "cancellation"
	ModificationTypeCorrection        ModificationType = "correction"
	ModificationTypeSystemAdjustment  ModificationType = "system_adjustment"
)

// ModificationStatus represents the status of a modification request
type ModificationStatus string

const (
	ModificationStatusDraft      ModificationStatus = "draft"
	ModificationStatusPending    ModificationStatus = "pending"
	ModificationStatusInReview   ModificationStatus = "in_review"
	ModificationStatusApproved   ModificationStatus = "approved"
	ModificationStatusRejected   ModificationStatus = "rejected"
	ModificationStatusProcessing ModificationStatus = "processing"
	ModificationStatusCompleted  ModificationStatus = "completed"
	ModificationStatusFailed     ModificationStatus = "failed"
	ModificationStatusCancelled  ModificationStatus = "cancelled"
	ModificationStatusReverted   ModificationStatus = "reverted"
)

// PolicyModification represents a change made to a policy with enterprise features
type PolicyModification struct {
	database.BaseModel

	// Core Fields
	PolicyID         uuid.UUID          `gorm:"type:uuid;not null;index" json:"policy_id" validate:"required"`
	ModificationType ModificationType   `gorm:"type:varchar(50);not null;index" json:"modification_type" validate:"required"`
	Status           ModificationStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	Priority         string             `gorm:"type:varchar(20);default:'normal'" json:"priority" validate:"oneof=low normal high critical"`

	// Change Details
	FieldModified string          `gorm:"type:varchar(100);index" json:"field_modified"`
	OldValue      json.RawMessage `gorm:"type:jsonb" json:"old_value"`
	NewValue      json.RawMessage `gorm:"type:jsonb" json:"new_value"`
	ChangeSet     []FieldChange   `gorm:"type:jsonb" json:"change_set"` // For multiple field changes

	// Business Context
	Reason                string `gorm:"type:text;not null" json:"reason" validate:"required,min=10"`
	BusinessJustification string `gorm:"type:text" json:"business_justification"`
	CustomerRequested     bool   `gorm:"default:false" json:"customer_requested"`
	SystemGenerated       bool   `gorm:"default:false" json:"system_generated"`
	ComplianceRelated     bool   `gorm:"default:false" json:"compliance_related"`

	// Approval Workflow
	RequiresApproval   bool       `gorm:"default:false;index" json:"requires_approval"`
	ApprovalLevel      int        `gorm:"default:0" json:"approval_level"` // 0=none, 1=supervisor, 2=manager, 3=director
	ApprovalWorkflowID *uuid.UUID `gorm:"type:uuid" json:"approval_workflow_id"`
	ApprovalNotes      string     `gorm:"type:text" json:"approval_notes"`
	RejectionReason    string     `gorm:"type:text" json:"rejection_reason"`

	// Financial Impact
	ImpactOnPremium    Money `gorm:"embedded;embeddedPrefix:premium_impact_" json:"impact_on_premium"`
	ImpactOnCoverage   Money `gorm:"embedded;embeddedPrefix:coverage_impact_" json:"impact_on_coverage"`
	RefundAmount       Money `gorm:"embedded;embeddedPrefix:refund_" json:"refund_amount"`
	AdditionalCharge   Money `gorm:"embedded;embeddedPrefix:charge_" json:"additional_charge"`
	ProRatedAdjustment bool  `gorm:"default:true" json:"pro_rated_adjustment"`

	// Dates and Timing
	RequestedDate  time.Time  `gorm:"not null;index" json:"requested_date"`
	EffectiveDate  time.Time  `gorm:"not null;index" json:"effective_date" validate:"required"`
	ExpirationDate *time.Time `json:"expiration_date"` // For temporary modifications
	ProcessedDate  *time.Time `json:"processed_date"`
	ApprovedDate   *time.Time `json:"approved_date"`
	CompletedDate  *time.Time `json:"completed_date"`
	RevertedDate   *time.Time `json:"reverted_date"`

	// User Tracking
	RequestedBy     uuid.UUID  `gorm:"type:uuid;not null;index" json:"requested_by" validate:"required"`
	RequestedByType string     `gorm:"type:varchar(20)" json:"requested_by_type"` // customer, agent, system, admin
	ApprovedBy      *uuid.UUID `gorm:"type:uuid;index" json:"approved_by"`
	ProcessedBy     *uuid.UUID `gorm:"type:uuid" json:"processed_by"`
	RevertedBy      *uuid.UUID `gorm:"type:uuid" json:"reverted_by"`

	// Documentation
	DocumentationURLs []string            `gorm:"type:jsonb" json:"documentation_urls"`
	SupportingDocs    []DocumentReference `gorm:"type:jsonb" json:"supporting_docs"`
	CustomerConsent   bool                `gorm:"default:false" json:"customer_consent"`
	ConsentTimestamp  *time.Time          `json:"consent_timestamp"`
	ConsentMethod     string              `gorm:"type:varchar(50)" json:"consent_method"` // email, sms, phone, in_person, online

	// Validation and Compliance
	ValidationRules     []ValidationResult `gorm:"type:jsonb" json:"validation_rules"`
	ComplianceChecks    []ComplianceCheck  `gorm:"type:jsonb" json:"compliance_checks"`
	RegulatoryImpact    string             `gorm:"type:text" json:"regulatory_impact"`
	RequiresReporting   bool               `gorm:"default:false" json:"requires_reporting"`
	ReportedToRegulator bool               `gorm:"default:false" json:"reported_to_regulator"`

	// Communication
	CustomerNotified   bool       `gorm:"default:false" json:"customer_notified"`
	NotificationDate   *time.Time `json:"notification_date"`
	NotificationMethod string     `gorm:"type:varchar(50)" json:"notification_method"`
	ConfirmationCode   string     `gorm:"type:varchar(50);index" json:"confirmation_code"`

	// Audit and Versioning
	Version                int                 `gorm:"default:1" json:"version"`
	PreviousModificationID *uuid.UUID          `gorm:"type:uuid" json:"previous_modification_id"`
	IsRevertible           bool                `gorm:"default:true" json:"is_revertible"`
	RevertedFrom           *uuid.UUID          `gorm:"type:uuid" json:"reverted_from"`
	AuditTrail             []ModificationAudit `gorm:"type:jsonb" json:"audit_trail"`

	// Error Handling
	ErrorMessage  string          `gorm:"type:text" json:"error_message"`
	ErrorDetails  json.RawMessage `gorm:"type:jsonb" json:"error_details"`
	RetryCount    int             `gorm:"default:0" json:"retry_count"`
	MaxRetries    int             `gorm:"default:3" json:"max_retries"`
	LastRetryDate *time.Time      `json:"last_retry_date"`

	// Performance Metrics
	ProcessingTime int `json:"processing_time_ms"` // milliseconds
	ValidationTime int `json:"validation_time_ms"`
	ApprovalTime   int `json:"approval_time_hours"` // hours from request to approval

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
	// User relationships (Requester, Approver, Processor) should be loaded via service layer
	// using RequestedBy, ApprovedBy, ProcessedBy UUIDs to avoid circular import
	DependentMods []PolicyModification `gorm:"foreignKey:PreviousModificationID" json:"dependent_modifications,omitempty"`
}

// FieldChange represents a single field change
type FieldChange struct {
	Field          string      `json:"field"`
	Path           string      `json:"path"` // JSON path for nested fields
	OldValue       interface{} `json:"old_value"`
	NewValue       interface{} `json:"new_value"`
	DataType       string      `json:"data_type"`
	IsRequired     bool        `json:"is_required"`
	IsEncrypted    bool        `json:"is_encrypted"`
	ValidationRule string      `json:"validation_rule"`
}

// DocumentReference represents a supporting document
type DocumentReference struct {
	DocumentID   uuid.UUID  `json:"document_id"`
	DocumentType string     `json:"document_type"`
	FileName     string     `json:"file_name"`
	URL          string     `json:"url"`
	UploadedAt   time.Time  `json:"uploaded_at"`
	UploadedBy   uuid.UUID  `json:"uploaded_by"`
	IsVerified   bool       `json:"is_verified"`
	VerifiedBy   *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
}

// ValidationResult represents the result of a validation rule
type ValidationResult struct {
	RuleID      string    `json:"rule_id"`
	RuleName    string    `json:"rule_name"`
	RuleType    string    `json:"rule_type"`
	Passed      bool      `json:"passed"`
	Message     string    `json:"message"`
	Severity    string    `json:"severity"` // info, warning, error, critical
	ValidatedAt time.Time `json:"validated_at"`
	ValidatedBy string    `json:"validated_by"` // system, manual, external
}

// ModificationAudit represents an audit entry for modification workflow
type ModificationAudit struct {
	Timestamp  time.Time          `json:"timestamp"`
	Action     string             `json:"action"`
	Actor      uuid.UUID          `json:"actor"`
	ActorName  string             `json:"actor_name"`
	ActorRole  string             `json:"actor_role"`
	FromStatus ModificationStatus `json:"from_status"`
	ToStatus   ModificationStatus `json:"to_status"`
	Comments   string             `json:"comments"`
	IPAddress  string             `json:"ip_address"`
	UserAgent  string             `json:"user_agent"`
	SessionID  string             `json:"session_id"`
}

// TableName returns the table name for PolicyModification model
func (PolicyModification) TableName() string {
	return "policy_modifications"
}

// BeforeCreate handles pre-creation logic
func (pm *PolicyModification) BeforeCreate(tx *gorm.DB) error {
	// Call parent BeforeCreate to handle UUID generation
	if err := pm.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Set requested date if not set
	if pm.RequestedDate.IsZero() {
		pm.RequestedDate = time.Now()
	}

	// Generate confirmation code
	if pm.ConfirmationCode == "" {
		pm.ConfirmationCode = generateConfirmationCode()
	}

	// Determine if approval is required based on type and impact
	pm.DetermineApprovalRequirement()

	// Add initial audit entry
	pm.AddAuditEntry("created", uuid.Nil, "", "", ModificationStatus(""), pm.Status, "Modification request created")

	return nil
}

// ============================================
// BUSINESS LOGIC METHODS
// ============================================

// IsApproved checks if the modification has been approved
func (pm *PolicyModification) IsApproved() bool {
	return pm.Status == ModificationStatusApproved && pm.ApprovedBy != nil
}

// IsCompleted checks if the modification is completed
func (pm *PolicyModification) IsCompleted() bool {
	return pm.Status == ModificationStatusCompleted && pm.CompletedDate != nil
}

// CanProcess checks if the modification can be processed
func (pm *PolicyModification) CanProcess() bool {
	if pm.RequiresApproval && !pm.IsApproved() {
		return false
	}
	return pm.Status == ModificationStatusPending ||
		pm.Status == ModificationStatusApproved ||
		pm.Status == ModificationStatusProcessing
}

// CanRevert checks if the modification can be reverted
func (pm *PolicyModification) CanRevert() bool {
	return pm.IsRevertible &&
		pm.Status == ModificationStatusCompleted &&
		pm.RevertedDate == nil &&
		time.Since(pm.EffectiveDate) < 90*24*time.Hour // Within 90 days
}

// RequiresHighLevelApproval checks if high-level approval is needed
func (pm *PolicyModification) RequiresHighLevelApproval() bool {
	// Check financial impact
	if pm.ImpactOnPremium.Amount > 1000 || pm.ImpactOnCoverage.Amount > 10000 {
		return true
	}

	// Check modification type
	switch pm.ModificationType {
	case ModificationTypeOwnershipTransfer,
		ModificationTypeCancellation,
		ModificationTypePremiumAdjustment:
		return true
	}

	// Check if compliance related
	if pm.ComplianceRelated {
		return true
	}

	return false
}

// DetermineApprovalRequirement determines if approval is required
func (pm *PolicyModification) DetermineApprovalRequirement() {
	// System adjustments typically don't need approval
	if pm.SystemGenerated && pm.ModificationType == ModificationTypeSystemAdjustment {
		pm.RequiresApproval = false
		pm.ApprovalLevel = 0
		return
	}

	// Address changes typically don't need approval
	if pm.ModificationType == ModificationTypeAddressChange {
		pm.RequiresApproval = false
		pm.ApprovalLevel = 0
		return
	}

	// Coverage increases need approval
	if pm.ModificationType == ModificationTypeCoverageIncrease {
		pm.RequiresApproval = true
		pm.ApprovalLevel = 1
	}

	// High impact changes need manager approval
	if pm.RequiresHighLevelApproval() {
		pm.RequiresApproval = true
		pm.ApprovalLevel = 2

		// Critical changes need director approval
		if pm.Priority == "critical" {
			pm.ApprovalLevel = 3
		}
	}
}

// AddAuditEntry adds an audit trail entry
func (pm *PolicyModification) AddAuditEntry(action string, actor uuid.UUID, actorName, actorRole string,
	fromStatus, toStatus ModificationStatus, comments string) {
	audit := ModificationAudit{
		Timestamp:  time.Now(),
		Action:     action,
		Actor:      actor,
		ActorName:  actorName,
		ActorRole:  actorRole,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		Comments:   comments,
	}
	pm.AuditTrail = append(pm.AuditTrail, audit)
}

// SetStatus updates the modification status with audit
func (pm *PolicyModification) SetStatus(newStatus ModificationStatus, actor uuid.UUID, actorName, actorRole, comments string) error {
	oldStatus := pm.Status

	// Validate status transition
	if !pm.IsValidStatusTransition(newStatus) {
		return errors.New("invalid status transition")
	}

	pm.Status = newStatus
	pm.AddAuditEntry("status_change", actor, actorName, actorRole, oldStatus, newStatus, comments)

	// Update relevant dates
	now := time.Now()
	switch newStatus {
	case ModificationStatusApproved:
		pm.ApprovedDate = &now
		pm.ApprovedBy = &actor
	case ModificationStatusProcessing:
		pm.ProcessedDate = &now
		pm.ProcessedBy = &actor
	case ModificationStatusCompleted:
		pm.CompletedDate = &now
	case ModificationStatusReverted:
		pm.RevertedDate = &now
		pm.RevertedBy = &actor
	}

	return nil
}

// IsValidStatusTransition checks if a status transition is valid
func (pm *PolicyModification) IsValidStatusTransition(newStatus ModificationStatus) bool {
	validTransitions := map[ModificationStatus][]ModificationStatus{
		ModificationStatusDraft:      {ModificationStatusPending, ModificationStatusCancelled},
		ModificationStatusPending:    {ModificationStatusInReview, ModificationStatusApproved, ModificationStatusRejected, ModificationStatusCancelled},
		ModificationStatusInReview:   {ModificationStatusApproved, ModificationStatusRejected, ModificationStatusPending},
		ModificationStatusApproved:   {ModificationStatusProcessing, ModificationStatusCancelled},
		ModificationStatusProcessing: {ModificationStatusCompleted, ModificationStatusFailed},
		ModificationStatusCompleted:  {ModificationStatusReverted},
		ModificationStatusFailed:     {ModificationStatusPending, ModificationStatusCancelled}, // Can retry
	}

	allowedStatuses, exists := validTransitions[pm.Status]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}

// Validate performs comprehensive validation
func (pm *PolicyModification) Validate() error {
	// Required fields
	if pm.PolicyID == uuid.Nil {
		return errors.New("policy ID is required")
	}
	if pm.ModificationType == "" {
		return errors.New("modification type is required")
	}
	if pm.Reason == "" || len(pm.Reason) < 10 {
		return errors.New("reason must be at least 10 characters")
	}
	if pm.RequestedBy == uuid.Nil {
		return errors.New("requester ID is required")
	}

	// Date validations
	if pm.EffectiveDate.IsZero() {
		return errors.New("effective date is required")
	}
	if pm.EffectiveDate.Before(pm.RequestedDate) {
		return errors.New("effective date cannot be before requested date")
	}

	// Expiration date validation for temporary modifications
	if pm.ExpirationDate != nil && pm.ExpirationDate.Before(pm.EffectiveDate) {
		return errors.New("expiration date must be after effective date")
	}

	// Financial validation
	if pm.RefundAmount.Amount < 0 || pm.AdditionalCharge.Amount < 0 {
		return errors.New("financial amounts cannot be negative")
	}

	// Customer consent validation
	if pm.CustomerRequested && !pm.CustomerConsent {
		return errors.New("customer consent is required for customer-requested modifications")
	}

	return nil
}

// CalculateFinancialImpact calculates the total financial impact
func (pm *PolicyModification) CalculateFinancialImpact() Money {
	// Start with premium impact
	totalImpact := pm.ImpactOnPremium

	// Add additional charges
	if pm.AdditionalCharge.Amount > 0 {
		totalImpact, _ = totalImpact.Add(pm.AdditionalCharge)
	}

	// Subtract refunds
	if pm.RefundAmount.Amount > 0 {
		totalImpact, _ = totalImpact.Subtract(pm.RefundAmount)
	}

	return totalImpact
}

// GetProcessingDuration returns the duration from request to completion
func (pm *PolicyModification) GetProcessingDuration() time.Duration {
	if pm.CompletedDate == nil {
		return time.Since(pm.RequestedDate)
	}
	return pm.CompletedDate.Sub(pm.RequestedDate)
}

// IsExpired checks if a temporary modification has expired
func (pm *PolicyModification) IsExpired() bool {
	if pm.ExpirationDate == nil {
		return false
	}
	return time.Now().After(*pm.ExpirationDate)
}

// ShouldRetry determines if a failed modification should be retried
func (pm *PolicyModification) ShouldRetry() bool {
	return pm.Status == ModificationStatusFailed &&
		pm.RetryCount < pm.MaxRetries &&
		!pm.SystemGenerated // Don't retry system-generated modifications
}

// GetApprovalLevelDescription returns a description of the approval level
func (pm *PolicyModification) GetApprovalLevelDescription() string {
	switch pm.ApprovalLevel {
	case 0:
		return "No approval required"
	case 1:
		return "Supervisor approval required"
	case 2:
		return "Manager approval required"
	case 3:
		return "Director approval required"
	default:
		return "Unknown approval level"
	}
}

// IsHighPriority checks if the modification is high priority
func (pm *PolicyModification) IsHighPriority() bool {
	return pm.Priority == "high" || pm.Priority == "critical"
}

// RequiresCustomerNotification checks if customer should be notified
func (pm *PolicyModification) RequiresCustomerNotification() bool {
	// Always notify for customer-requested changes
	if pm.CustomerRequested {
		return true
	}

	// Notify for significant changes
	switch pm.ModificationType {
	case ModificationTypeCoverageIncrease,
		ModificationTypeCoverageDecrease,
		ModificationTypePremiumAdjustment,
		ModificationTypeCancellation,
		ModificationTypeSuspension,
		ModificationTypeReinstatement:
		return true
	}

	// Notify if there's a financial impact
	if pm.ImpactOnPremium.Amount != 0 ||
		pm.RefundAmount.Amount > 0 ||
		pm.AdditionalCharge.Amount > 0 {
		return true
	}

	return false
}

// GetValidationSeverity returns the highest severity from validation results
func (pm *PolicyModification) GetValidationSeverity() string {
	severityOrder := map[string]int{
		"info":     1,
		"warning":  2,
		"error":    3,
		"critical": 4,
	}

	highestSeverity := "info"
	highestOrder := 1

	for _, validation := range pm.ValidationRules {
		if order, exists := severityOrder[validation.Severity]; exists && order > highestOrder {
			highestSeverity = validation.Severity
			highestOrder = order
		}
	}

	return highestSeverity
}

// HasValidationErrors checks if there are any validation errors
func (pm *PolicyModification) HasValidationErrors() bool {
	for _, validation := range pm.ValidationRules {
		if !validation.Passed && (validation.Severity == "error" || validation.Severity == "critical") {
			return true
		}
	}
	return false
}

// GetEstimatedProcessingTime returns estimated processing time in hours
func (pm *PolicyModification) GetEstimatedProcessingTime() int {
	baseTime := 1 // 1 hour base

	// Add time for approval
	if pm.RequiresApproval {
		baseTime += pm.ApprovalLevel * 24 // 24 hours per approval level
	}

	// Add time for complex changes
	switch pm.ModificationType {
	case ModificationTypeOwnershipTransfer:
		baseTime += 48
	case ModificationTypeCancellation:
		baseTime += 24
	case ModificationTypePremiumAdjustment:
		baseTime += 12
	}

	// High priority gets expedited
	if pm.Priority == "critical" {
		baseTime = baseTime / 2
	}

	return baseTime
}

// generateConfirmationCode generates a random 8-character confirmation code
func generateConfirmationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
