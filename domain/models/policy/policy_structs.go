package policy

import (
	"time"

	"github.com/google/uuid"
)

// ============================================
// EMBEDDED STRUCT TYPES FOR ORGANIZATION
// ============================================

// PolicyIdentification groups all identification fields
type PolicyIdentification struct {
	PolicyNumber       string     `gorm:"uniqueIndex;not null" json:"policy_number" validate:"required"`
	PolicyVersion      int        `gorm:"default:1" json:"policy_version" validate:"min=1"`
	ContractNumber     string     `gorm:"index" json:"contract_number"`
	QuoteID            *uuid.UUID `gorm:"type:uuid;index" json:"quote_id"`
	ParentPolicyID     *uuid.UUID `gorm:"type:uuid;index" json:"parent_policy_id"`     // For renewals
	BundleID           *uuid.UUID `gorm:"type:uuid;index" json:"bundle_id"`            // For bundled policies
	CorporateAccountID *uuid.UUID `gorm:"type:uuid;index" json:"corporate_account_id"` // For corporate policies
}

// PolicyClassification groups all classification fields
type PolicyClassification struct {
	PolicyType     PolicyType  `gorm:"type:varchar(50);not null;index" json:"policy_type" validate:"required"`
	PolicyCategory string      `gorm:"type:varchar(50)" json:"policy_category" validate:"omitempty,oneof=individual family corporate group"`
	CoverageType   string      `gorm:"type:varchar(50);not null" json:"coverage_type" validate:"required"`
	CoverageLevel  string      `gorm:"type:varchar(20)" json:"coverage_level" validate:"omitempty,oneof=basic standard premium platinum enterprise"`
	BusinessLine   string      `gorm:"type:varchar(50)" json:"business_line" validate:"omitempty,oneof=personal commercial"`
	Channel        ChannelType `gorm:"type:varchar(50);index" json:"channel" validate:"required"`
	PartnerCode    string      `gorm:"type:varchar(20);index" json:"partner_code"`
	ProductLine    string      `gorm:"type:varchar(50)" json:"product_line"`
}

// PolicyCoverageDetails groups all coverage-related fields
type PolicyCoverageDetails struct {
	CoverageAmount     Money          `gorm:"embedded;embeddedPrefix:coverage_" json:"coverage_amount" validate:"required"`
	Deductible         Money          `gorm:"embedded;embeddedPrefix:deductible_" json:"deductible"`
	DeductibleType     DeductibleType `gorm:"type:varchar(20)" json:"deductible_type"`
	ExcessAmount       Money          `gorm:"embedded;embeddedPrefix:excess_" json:"excess_amount"`
	OutOfPocketMax     Money          `gorm:"embedded;embeddedPrefix:oop_max_" json:"out_of_pocket_max"`
	CoinsurancePercent float64        `gorm:"type:decimal(5,2)" json:"coinsurance_percent" validate:"min=0,max=100"`
	AggregateLimit     Money          `gorm:"embedded;embeddedPrefix:aggregate_" json:"aggregate_limit"`
	RemainingLimit     Money          `gorm:"embedded;embeddedPrefix:remaining_" json:"remaining_limit"`
	CoverageLimits     CoverageLimits `gorm:"type:jsonb" json:"coverage_limits"`
	CoverageRegions    []Region       `gorm:"type:jsonb" json:"coverage_regions"`
	Exclusions         []Exclusion    `gorm:"type:jsonb" json:"exclusions"`
}

// PolicyPricing groups all pricing and premium fields
type PolicyPricing struct {
	BasePremium         Money        `gorm:"embedded;embeddedPrefix:base_premium_" json:"base_premium" validate:"required"`
	RiskAdjustedPremium Money        `gorm:"embedded;embeddedPrefix:risk_premium_" json:"risk_adjusted_premium"`
	FinalPremium        Money        `gorm:"embedded;embeddedPrefix:final_premium_" json:"final_premium" validate:"required"`
	PremiumTax          Money        `gorm:"embedded;embeddedPrefix:premium_tax_" json:"premium_tax"`
	AdminFee            Money        `gorm:"embedded;embeddedPrefix:admin_fee_" json:"admin_fee"`
	ProcessingFee       Money        `gorm:"embedded;embeddedPrefix:processing_fee_" json:"processing_fee"`
	TotalAmount         Money        `gorm:"embedded;embeddedPrefix:total_" json:"total_amount"`
	Currency            CurrencyCode `gorm:"type:varchar(3);default:'USD';index" json:"currency" validate:"required,len=3"`
	PricingTier         string       `gorm:"type:varchar(20)" json:"pricing_tier"`
	PricingModel        string       `gorm:"type:varchar(30)" json:"pricing_model"` // fixed, usage-based, tiered
	LastPriceUpdate     *time.Time   `json:"last_price_update"`
}

// PolicyDiscounts groups all discount-related fields
type PolicyDiscounts struct {
	DiscountType       string  `gorm:"type:varchar(30)" json:"discount_type"`
	DiscountAmount     Money   `gorm:"embedded;embeddedPrefix:discount_" json:"discount_amount"`
	DiscountPercentage float64 `gorm:"type:decimal(5,2)" json:"discount_percentage" validate:"min=0,max=100"`
	DiscountReason     string  `gorm:"type:text" json:"discount_reason"`
	LoyaltyDiscount    Money   `gorm:"embedded;embeddedPrefix:loyalty_disc_" json:"loyalty_discount"`
	BundleDiscount     Money   `gorm:"embedded;embeddedPrefix:bundle_disc_" json:"bundle_discount"`
	NoClaimsBonus      Money   `gorm:"embedded;embeddedPrefix:ncb_" json:"no_claims_bonus"`
	CorporateDiscount  Money   `gorm:"embedded;embeddedPrefix:corp_disc_" json:"corporate_discount"`
	PromoCode          string  `gorm:"type:varchar(20);index" json:"promo_code"`
	PromoDiscount      Money   `gorm:"embedded;embeddedPrefix:promo_disc_" json:"promo_discount"`
	EarlyPayDiscount   Money   `gorm:"embedded;embeddedPrefix:early_pay_" json:"early_pay_discount"`
	TotalDiscounts     Money   `gorm:"embedded;embeddedPrefix:total_disc_" json:"total_discounts"`
}

// PolicyLoadings groups all loading-related fields
type PolicyLoadings struct {
	LoadingAmount     Money   `gorm:"embedded;embeddedPrefix:loading_" json:"loading_amount"`
	LoadingPercentage float64 `gorm:"type:decimal(5,2)" json:"loading_percentage" validate:"min=0,max=200"`
	LoadingReason     string  `gorm:"type:text" json:"loading_reason"`
	LoadingType       string  `gorm:"type:varchar(30)" json:"loading_type"`
	RiskLoading       Money   `gorm:"embedded;embeddedPrefix:risk_loading_" json:"risk_loading"`
	AdminLoading      Money   `gorm:"embedded;embeddedPrefix:admin_loading_" json:"admin_loading"`
	TotalLoadings     Money   `gorm:"embedded;embeddedPrefix:total_loading_" json:"total_loadings"`
}

// PolicyPaymentInfo groups all payment-related fields
type PolicyPaymentInfo struct {
	PaymentFrequency   PaymentFrequency `gorm:"type:varchar(20);default:'monthly';index" json:"payment_frequency" validate:"required"`
	PaymentMethod      string           `gorm:"type:varchar(30)" json:"payment_method" validate:"required"`
	PaymentStatus      PaymentStatus    `gorm:"type:varchar(20);default:'pending';index" json:"payment_status"`
	AutoRenewal        bool             `gorm:"default:false;index" json:"auto_renewal"`
	AutoPayment        bool             `gorm:"default:false" json:"auto_payment"`
	BillingCycle       int              `json:"billing_cycle" validate:"min=0"`
	NextBillingDate    *time.Time       `gorm:"index" json:"next_billing_date"`
	LastPaymentDate    *time.Time       `json:"last_payment_date"`
	PaymentGracePeriod int              `json:"payment_grace_period" validate:"min=0,max=60"` // Days
	OutstandingAmount  Money            `gorm:"embedded;embeddedPrefix:outstanding_" json:"outstanding_amount"`
	InstallmentPlan    *InstallmentPlan `gorm:"type:jsonb" json:"installment_plan"`
	PaymentRetries     int              `gorm:"default:0" json:"payment_retries"`
	LastRetryDate      *time.Time       `json:"last_retry_date"`
	PaymentProcessor   string           `gorm:"type:varchar(50)" json:"payment_processor"`
	PaymentReference   string           `gorm:"type:varchar(100);index" json:"payment_reference"`
}

// PolicyLifecycle groups all lifecycle-related fields
type PolicyLifecycle struct {
	Status             PolicyStatus       `gorm:"type:varchar(20);not null;default:'pending';index" json:"status" validate:"required"`
	UnderwritingStatus UnderwritingStatus `gorm:"type:varchar(20);default:'pending';index" json:"underwriting_status"`
	ApprovalStatus     string             `gorm:"type:varchar(20);index" json:"approval_status"`
	ActivationDate     *time.Time         `gorm:"index" json:"activation_date"`
	EffectiveDate      time.Time          `gorm:"not null;index" json:"effective_date" validate:"required"`
	ExpirationDate     time.Time          `gorm:"not null;index" json:"expiration_date" validate:"required,gtfield=EffectiveDate"`
	RenewalDate        *time.Time         `gorm:"index" json:"renewal_date"`
	CancellationDate   *time.Time         `json:"cancellation_date"`
	CancellationReason string             `gorm:"type:text" json:"cancellation_reason"`
	SuspensionDate     *time.Time         `json:"suspension_date"`
	SuspensionReason   string             `gorm:"type:text" json:"suspension_reason"`
	ReinstateDate      *time.Time         `json:"reinstate_date"`
	LastModifiedDate   *time.Time         `json:"last_modified_date"`
	LastReviewDate     *time.Time         `json:"last_review_date"`
	NextReviewDate     *time.Time         `gorm:"index" json:"next_review_date"`
	StateTransitions   []StateTransition  `gorm:"type:jsonb" json:"state_transitions"`
}

// PolicyRiskAssessment groups all risk-related fields
type PolicyRiskAssessment struct {
	RiskScore         float64      `gorm:"type:decimal(5,2);default:50.0;index" json:"risk_score" validate:"min=0,max=100"`
	RiskCategory      RiskCategory `gorm:"type:varchar(20);index" json:"risk_category"`
	RiskFactors       []RiskFactor `gorm:"type:jsonb" json:"risk_factors"`
	DeviceRiskScore   float64      `gorm:"type:decimal(5,2)" json:"device_risk_score" validate:"min=0,max=100"`
	UserRiskScore     float64      `gorm:"type:decimal(5,2)" json:"user_risk_score" validate:"min=0,max=100"`
	LocationRiskScore float64      `gorm:"type:decimal(5,2)" json:"location_risk_score" validate:"min=0,max=100"`
	BehaviorRiskScore float64      `gorm:"type:decimal(5,2)" json:"behavior_risk_score" validate:"min=0,max=100"`
	FraudRiskScore    float64      `gorm:"type:decimal(5,2)" json:"fraud_risk_score" validate:"min=0,max=100"`
	ClaimFrequency    int          `json:"claim_frequency" validate:"min=0"`
	ClaimSeverity     float64      `gorm:"type:decimal(10,2)" json:"claim_severity" validate:"min=0"`
	LossRatio         float64      `gorm:"type:decimal(5,2)" json:"loss_ratio" validate:"min=0"`
	RiskTrend         string       `gorm:"type:varchar(20)" json:"risk_trend"` // increasing, stable, decreasing
	LastRiskUpdate    *time.Time   `json:"last_risk_update"`
	RiskMatrix        RiskMatrix   `gorm:"type:jsonb" json:"risk_matrix"`
}

// PolicyUnderwritingInfo groups all underwriting fields
type PolicyUnderwritingInfo struct {
	UnderwritingNotes  string             `gorm:"type:text" json:"underwriting_notes"`
	UnderwriterID      *uuid.UUID         `gorm:"type:uuid;index" json:"underwriter_id"`
	UnderwritingDate   *time.Time         `json:"underwriting_date"`
	UnderwritingMethod string             `gorm:"type:varchar(20)" json:"underwriting_method"` // auto, manual, hybrid
	RequiresInspection bool               `gorm:"default:false" json:"requires_inspection"`
	InspectionStatus   string             `gorm:"type:varchar(20)" json:"inspection_status"`
	InspectionDate     *time.Time         `json:"inspection_date"`
	InspectionNotes    string             `gorm:"type:text" json:"inspection_notes"`
	InspectorID        *uuid.UUID         `gorm:"type:uuid" json:"inspector_id"`
	UnderwritingRules  []UnderwritingRule `gorm:"type:jsonb" json:"underwriting_rules"`
	ReferralReason     string             `gorm:"type:text" json:"referral_reason"`
	Conditions         []string           `gorm:"type:jsonb" json:"conditions"`
}

// PolicyCompliance groups all compliance and regulatory fields
type PolicyCompliance struct {
	ComplianceStatus   string            `gorm:"type:varchar(20);index" json:"compliance_status"`
	RegulatoryRegion   string            `gorm:"type:varchar(50)" json:"regulatory_region"`
	TaxJurisdiction    string            `gorm:"type:varchar(50)" json:"tax_jurisdiction"`
	ComplianceChecks   []ComplianceCheck `gorm:"type:jsonb" json:"compliance_checks"`
	KYCStatus          string            `gorm:"type:varchar(20)" json:"kyc_status"`
	AMLStatus          string            `gorm:"type:varchar(20)" json:"aml_status"`
	SanctionsCheck     bool              `gorm:"default:false" json:"sanctions_check"`
	DataPrivacyConsent bool              `gorm:"default:false" json:"data_privacy_consent"`
	GDPRCompliant      bool              `gorm:"default:false" json:"gdpr_compliant"`
	LastComplianceDate *time.Time        `json:"last_compliance_date"`
	ComplianceOfficer  *uuid.UUID        `gorm:"type:uuid" json:"compliance_officer"`
	ComplianceNotes    string            `gorm:"type:text" json:"compliance_notes"`
}

// PolicyDocumentation groups all document-related fields
type PolicyDocumentation struct {
	PolicyDocumentURL  string               `json:"policy_document_url"`
	TermsDocumentURL   string               `json:"terms_document_url"`
	CertificateURL     string               `json:"certificate_url"`
	PreferredLanguage  string               `gorm:"type:varchar(5);default:'en'" json:"preferred_language"`
	DocumentVersion    string               `gorm:"type:varchar(20)" json:"document_version"`
	DocumentGenerated  *time.Time           `json:"document_generated"`
	DocumentsSent      []DocumentSentRecord `gorm:"type:jsonb" json:"documents_sent"`
	DigitalSignatureID string               `gorm:"type:varchar(100)" json:"digital_signature_id"`
	SignedDate         *time.Time           `json:"signed_date"`
}

// PolicyCommunication groups all communication-related fields
type PolicyCommunication struct {
	CommunicationPrefs   CommunicationPreferences `gorm:"type:jsonb" json:"communication_prefs"`
	LastCommunication    *time.Time               `json:"last_communication"`
	NextCommunication    *time.Time               `gorm:"index" json:"next_communication"`
	CommunicationHistory []CommunicationRecord    `gorm:"type:jsonb" json:"communication_history"`
	PreferredChannel     string                   `gorm:"type:varchar(20)" json:"preferred_channel"`
	OptInMarketing       bool                     `gorm:"default:false" json:"opt_in_marketing"`
	OptInSurveys         bool                     `gorm:"default:false" json:"opt_in_surveys"`
}

// PolicyAnalytics groups all analytics and metrics fields
type PolicyAnalytics struct {
	ProfitabilityScore    float64            `gorm:"type:decimal(5,2)" json:"profitability_score" validate:"min=0,max=100"`
	CustomerLifetimeValue Money              `gorm:"embedded;embeddedPrefix:clv_" json:"customer_lifetime_value"`
	RetentionScore        float64            `gorm:"type:decimal(5,2)" json:"retention_score" validate:"min=0,max=100"`
	ChurnRisk             float64            `gorm:"type:decimal(5,2)" json:"churn_risk" validate:"min=0,max=100"`
	NetPromoterScore      int                `json:"net_promoter_score" validate:"min=-100,max=100"`
	CustomerSatisfaction  float64            `gorm:"type:decimal(3,2)" json:"customer_satisfaction" validate:"min=0,max=5"`
	EngagementScore       float64            `gorm:"type:decimal(5,2)" json:"engagement_score" validate:"min=0,max=100"`
	CrossSellPotential    float64            `gorm:"type:decimal(5,2)" json:"cross_sell_potential" validate:"min=0,max=100"`
	UpsellPotential       float64            `gorm:"type:decimal(5,2)" json:"upsell_potential" validate:"min=0,max=100"`
	LastAnalyticsUpdate   *time.Time         `json:"last_analytics_update"`
	PerformanceMetrics    map[string]float64 `gorm:"type:jsonb" json:"performance_metrics"`
}

// PolicyTermsConditions groups terms and conditions
type PolicyTermsConditions struct {
	Terms             string   `gorm:"type:text" json:"terms"`
	Conditions        string   `gorm:"type:text" json:"conditions"`
	SpecialConditions string   `gorm:"type:text" json:"special_conditions"`
	Endorsements      []string `gorm:"type:jsonb" json:"endorsements"`
	Riders            []string `gorm:"type:jsonb" json:"riders"`
	Waivers           []string `gorm:"type:jsonb" json:"waivers"`
	Warranties        []string `gorm:"type:jsonb" json:"warranties"`
	Clauses           []string `gorm:"type:jsonb" json:"clauses"`
}

// PolicyMetadata groups all metadata fields
type PolicyMetadata struct {
	Source        string            `gorm:"type:varchar(30)" json:"source"` // web, mobile, agent, partner
	CampaignCode  string            `gorm:"type:varchar(30);index" json:"campaign_code"`
	ReferralCode  string            `gorm:"type:varchar(30);index" json:"referral_code"`
	Tags          []string          `gorm:"type:jsonb" json:"tags"`
	CustomFields  map[string]string `gorm:"type:jsonb" json:"custom_fields"`
	Notes         string            `gorm:"type:text" json:"notes"`
	InternalNotes string            `gorm:"type:text" json:"internal_notes"`
	Version       int               `gorm:"default:1" json:"version"` // For optimistic locking
	IsArchived    bool              `gorm:"default:false;index" json:"is_archived"`
	ArchivedAt    *time.Time        `json:"archived_at"`
	ArchivedBy    *uuid.UUID        `gorm:"type:uuid" json:"archived_by"`
	ImportID      string            `gorm:"type:varchar(100);index" json:"import_id"`
	ExternalRef   string            `gorm:"type:varchar(100);index" json:"external_ref"`
}

// PolicyAudit groups audit-related fields
type PolicyAudit struct {
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedByName  string     `gorm:"type:varchar(100)" json:"created_by_name"`
	ModifiedBy     uuid.UUID  `gorm:"type:uuid;index" json:"modified_by"`
	ModifiedByName string     `gorm:"type:varchar(100)" json:"modified_by_name"`
	ApprovedBy     *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovedByName *string    `gorm:"type:varchar(100)" json:"approved_by_name"`
	ApprovalDate   *time.Time `json:"approval_date"`
	LastAuditDate  *time.Time `json:"last_audit_date"`
	// Audit trail should be loaded via service layer using audit log IDs to avoid circular import
	AuditTrailIDs []uuid.UUID `gorm:"type:jsonb" json:"audit_trail_ids"`
	ChangeHistory []Change    `gorm:"type:jsonb" json:"change_history"`
}

// PolicyRelationships groups all relationship fields
type PolicyRelationships struct {
	CustomerID uuid.UUID  `gorm:"type:uuid;not null;index" json:"customer_id" validate:"required"`
	DeviceID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"device_id" validate:"required"`
	ProductID  *uuid.UUID `gorm:"type:uuid;index" json:"product_id"`
	AgentID    *uuid.UUID `gorm:"type:uuid;index" json:"agent_id"`
	BrokerID   *uuid.UUID `gorm:"type:uuid;index" json:"broker_id"`
}

// ============================================
// SUPPORTING TYPES FOR EMBEDDED STRUCTS
// ============================================

// StateTransition represents a state change in the policy lifecycle
type StateTransition struct {
	FromState PolicyStatus `json:"from_state"`
	ToState   PolicyStatus `json:"to_state"`
	Timestamp time.Time    `json:"timestamp"`
	UserID    uuid.UUID    `json:"user_id"`
	Reason    string       `json:"reason"`
}

// RiskMatrix represents a risk assessment matrix
type RiskMatrix struct {
	Likelihood   int          `json:"likelihood"` // 1-5
	Impact       int          `json:"impact"`     // 1-5
	Score        int          `json:"score"`      // Likelihood * Impact
	Category     RiskCategory `json:"category"`
	Mitigations  []string     `json:"mitigations"`
	LastAssessed time.Time    `json:"last_assessed"`
	Assessor     uuid.UUID    `json:"assessor"`
}

// UnderwritingRule represents a rule applied during underwriting
type UnderwritingRule struct {
	RuleID    string    `json:"rule_id"`
	RuleName  string    `json:"rule_name"`
	Applied   bool      `json:"applied"`
	Result    string    `json:"result"`
	Impact    string    `json:"impact"`
	AppliedAt time.Time `json:"applied_at"`
}

// DocumentSentRecord represents a document sending record
type DocumentSentRecord struct {
	DocumentType string    `json:"document_type"`
	SentTo       string    `json:"sent_to"`
	SentVia      string    `json:"sent_via"`
	SentAt       time.Time `json:"sent_at"`
	Status       string    `json:"status"`
	TrackingID   string    `json:"tracking_id"`
}

// CommunicationPreferences represents communication preferences
type CommunicationPreferences struct {
	Email         bool     `json:"email"`
	SMS           bool     `json:"sms"`
	PushNotif     bool     `json:"push_notif"`
	PhoneCall     bool     `json:"phone_call"`
	PostalMail    bool     `json:"postal_mail"`
	FrequencyDays int      `json:"frequency_days"`
	BestTime      string   `json:"best_time"`
	TimeZone      string   `json:"timezone"`
	DoNotDisturb  []string `json:"do_not_disturb"` // Date ranges
}

// CommunicationRecord represents a communication history record
type CommunicationRecord struct {
	Type      string    `json:"type"`
	Channel   string    `json:"channel"`
	Subject   string    `json:"subject"`
	Content   string    `json:"content"`
	SentAt    time.Time `json:"sent_at"`
	Status    string    `json:"status"`
	Response  string    `json:"response"`
	MessageID string    `json:"message_id"`
}
