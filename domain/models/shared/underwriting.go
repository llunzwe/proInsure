package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/policy"
	"smartsure/internal/domain/types"
)

// UnderwritingDecision represents underwriting decisions for policies
type UnderwritingDecision struct {
	ID       uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	PolicyID *uuid.UUID `gorm:"type:uuid" json:"policy_id"`
	QuoteID  *uuid.UUID `gorm:"type:uuid" json:"quote_id"`
	UserID   uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`

	// Decision Details
	DecisionNumber string    `gorm:"uniqueIndex;not null" json:"decision_number"`
	DecisionType   string    `json:"decision_type"` // auto, manual, hybrid
	Decision       string    `json:"decision"`      // approved, declined, refer, pending
	DecisionDate   time.Time `json:"decision_date"`

	// Risk Assessment
	RiskScore    float64         `json:"risk_score"`
	RiskCategory string          `json:"risk_category"` // low, medium, high, very_high
	RiskFactors  types.JSONArray `gorm:"type:json" json:"risk_factors"`

	// Underwriting Details
	UnderwriterID      *uuid.UUID `gorm:"type:uuid" json:"underwriter_id"`
	UnderwriterNotes   string     `gorm:"type:text" json:"underwriter_notes"`
	AutoDecisionReason string     `json:"auto_decision_reason"`
	ManualOverride     bool       `json:"manual_override"`

	// Pricing Adjustments
	BasePremium    float64 `json:"base_premium"`
	RiskAdjustment float64 `json:"risk_adjustment"`
	FinalPremium   float64 `json:"final_premium"`
	LoadingFactor  float64 `json:"loading_factor"`

	// Conditions
	Conditions        types.JSONArray `gorm:"type:json" json:"conditions"`
	Exclusions        types.JSONArray `gorm:"type:json" json:"exclusions"`
	RequiredDocuments types.JSONArray `gorm:"type:json" json:"required_documents"`

	// Validity
	ValidUntil time.Time `json:"valid_until"`
	IsExpired  bool      `json:"is_expired"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Policy       *models.Policy     `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Quote        *policy.Quote      `gorm:"foreignKey:QuoteID" json:"quote,omitempty"`
	User         *models.User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Underwriter  *models.User       `gorm:"foreignKey:UnderwriterID" json:"underwriter,omitempty"`
	RulesApplied []UnderwritingRule `gorm:"many2many:underwriting_decision_rules;" json:"rules_applied,omitempty"`
}

// UnderwritingRule represents business rules for underwriting
type UnderwritingRule struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RuleCode    string    `gorm:"uniqueIndex;not null" json:"rule_code"`
	RuleName    string    `gorm:"not null" json:"rule_name"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // eligibility, pricing, risk, fraud

	// Rule Configuration
	RuleType    string `json:"rule_type"` // condition, calculation, validation
	Priority    int    `json:"priority"`
	IsActive    bool   `json:"is_active"`
	IsMandatory bool   `json:"is_mandatory"`

	// Rule Logic
	Condition  string `gorm:"type:json" json:"condition"`
	Action     string `gorm:"type:json" json:"action"`
	Parameters string `gorm:"type:json" json:"parameters"`

	// Impact
	RiskImpact     float64 `json:"risk_impact"`
	PremiumImpact  float64 `json:"premium_impact"`
	DecisionImpact string  `json:"decision_impact"` // approve, decline, refer

	// Applicability
	ApplicableProducts types.JSONArray `gorm:"type:json" json:"applicable_products"`
	ApplicableRegions  types.JSONArray `gorm:"type:json" json:"applicable_regions"`
	EffectiveFrom      time.Time       `json:"effective_from"`
	EffectiveTo        *time.Time      `json:"effective_to"`

	// Audit
	CreatedBy  uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	ApprovedBy *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// RiskAssessment represents detailed risk assessment
type RiskAssessment struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AssessmentNumber string    `gorm:"uniqueIndex;not null" json:"assessment_number"`
	EntityType       string    `json:"entity_type"` // user, device, policy, claim
	EntityID         uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`

	// Assessment Details
	AssessmentType string    `json:"assessment_type"` // initial, periodic, triggered
	TriggerReason  string    `json:"trigger_reason"`
	AssessmentDate time.Time `json:"assessment_date"`

	// Risk Scores
	OverallScore    float64 `json:"overall_score"`
	FraudScore      float64 `json:"fraud_score"`
	CreditScore     float64 `json:"credit_score"`
	BehavioralScore float64 `json:"behavioral_score"`
	DeviceScore     float64 `json:"device_score"`

	// Risk Categories
	RiskLevel   string `json:"risk_level"` // low, medium, high, critical
	FraudRisk   string `json:"fraud_risk"`
	PaymentRisk string `json:"payment_risk"`
	ClaimRisk   string `json:"claim_risk"`

	// Risk Factors
	PositiveFactors types.JSONArray `gorm:"type:json" json:"positive_factors"`
	NegativeFactors types.JSONArray `gorm:"type:json" json:"negative_factors"`
	RedFlags        types.JSONArray `gorm:"type:json" json:"red_flags"`

	// Data Sources
	DataSources        types.JSONArray `gorm:"type:json" json:"data_sources"`
	ExternalScores     string          `gorm:"type:json" json:"external_scores"`
	VerificationStatus string          `json:"verification_status"`

	// Actions & Recommendations
	RecommendedAction string          `json:"recommended_action"`
	RequiredActions   types.JSONArray `gorm:"type:json" json:"required_actions"`
	Mitigations       types.JSONArray `gorm:"type:json" json:"mitigations"`

	// Review
	ReviewRequired bool       `json:"review_required"`
	ReviewedBy     *uuid.UUID `gorm:"type:uuid" json:"reviewed_by"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
	ReviewNotes    string     `json:"review_notes"`

	// Validity
	ValidUntil         time.Time  `json:"valid_until"`
	NextAssessmentDate *time.Time `json:"next_assessment_date"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Reviewer *models.User `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

// RiskModel represents actuarial risk models
type RiskModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ModelCode string    `gorm:"uniqueIndex;not null" json:"model_code"`
	ModelName string    `gorm:"not null" json:"model_name"`
	ModelType string    `json:"model_type"` // pricing, fraud, claims, retention
	Version   string    `json:"version"`

	// Model Configuration
	Algorithm  string          `json:"algorithm"`
	Parameters string          `gorm:"type:json" json:"parameters"`
	Features   types.JSONArray `gorm:"type:json" json:"features"`
	Weights    string          `gorm:"type:json" json:"weights"`

	// Performance Metrics
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	F1Score   float64 `json:"f1_score"`

	// Training
	TrainingDataset   string    `json:"training_dataset"`
	TrainingDate      time.Time `json:"training_date"`
	ValidationDataset string    `json:"validation_dataset"`
	TestDataset       string    `json:"test_dataset"`

	// Deployment
	Status     string     `json:"status"` // development, testing, production, deprecated
	DeployedAt *time.Time `json:"deployed_at"`
	DeployedBy *uuid.UUID `gorm:"type:uuid" json:"deployed_by"`

	// Usage
	UsageCount int        `json:"usage_count"`
	LastUsedAt *time.Time `json:"last_used_at"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Deployer *models.User `gorm:"foreignKey:DeployedBy" json:"deployer,omitempty"`
}

// UnderwritingDocument represents documents required for underwriting
type UnderwritingDocument struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DecisionID   uuid.UUID `gorm:"type:uuid;not null" json:"decision_id"`
	DocumentType string    `json:"document_type"` // proof_of_purchase, device_photos, identity
	DocumentName string    `json:"document_name"`

	// Document Details
	DocumentURL  string `json:"document_url"`
	DocumentHash string `json:"document_hash"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`

	// Verification
	IsVerified        bool       `json:"is_verified"`
	VerifiedBy        *uuid.UUID `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt        *time.Time `json:"verified_at"`
	VerificationNotes string     `json:"verification_notes"`

	// Status
	Status          string `json:"status"` // pending, approved, rejected
	RejectionReason string `json:"rejection_reason"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Decision *UnderwritingDecision `gorm:"foreignKey:DecisionID" json:"decision,omitempty"`
	Verifier *models.User          `gorm:"foreignKey:VerifiedBy" json:"verifier,omitempty"`
}
