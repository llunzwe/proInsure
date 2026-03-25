package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PolicyRiskAccumulation represents advanced risk accumulation analysis
type PolicyRiskAccumulation struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID        uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	CalculationDate time.Time `gorm:"type:timestamp;not null" json:"calculation_date"`

	// Geographic Accumulation
	GeographicAccumulation  datatypes.JSON `gorm:"type:json" json:"geographic_accumulation"` // map[string]float64
	GeographicConcentration float64        `gorm:"type:decimal(5,2)" json:"geographic_concentration"`
	HighRiskZones           datatypes.JSON `gorm:"type:json" json:"high_risk_zones"`
	ExposureByRegion        datatypes.JSON `gorm:"type:json" json:"exposure_by_region"`

	// Peril Accumulation
	PerilAccumulation  datatypes.JSON `gorm:"type:json" json:"peril_accumulation"` // map[string]float64
	DominantPeril      string         `gorm:"type:varchar(100)" json:"dominant_peril"`
	PerilConcentration float64        `gorm:"type:decimal(5,2)" json:"peril_concentration"`
	MultiPerilExposure Money          `gorm:"embedded;embeddedPrefix:multi_peril_exposure_" json:"multi_peril_exposure"`

	// Catastrophe Exposure
	CATExposure        Money          `gorm:"embedded;embeddedPrefix:cat_exposure_" json:"cat_exposure"`
	CATModelResults    datatypes.JSON `gorm:"type:json" json:"cat_model_results"`
	NATCATExposure     Money          `gorm:"embedded;embeddedPrefix:natcat_exposure_" json:"natcat_exposure"`
	ManMadeCATexposure Money          `gorm:"embedded;embeddedPrefix:manmade_cat_exposure_" json:"manmade_cat_exposure"`

	// PML (Probable Maximum Loss)
	PMLCalculation Money          `gorm:"embedded;embeddedPrefix:pml_" json:"pml_calculation"`
	PMLPercentage  float64        `gorm:"type:decimal(5,2)" json:"pml_percentage"`
	PMLScenarios   datatypes.JSON `gorm:"type:json" json:"pml_scenarios"`
	ReturnPeriods  datatypes.JSON `gorm:"type:json" json:"return_periods"` // map[int]float64

	// VaR & Risk Metrics
	VaRCalculation     Money   `gorm:"embedded;embeddedPrefix:var_" json:"var_calculation"`
	VaRConfidenceLevel float64 `gorm:"type:decimal(5,2)" json:"var_confidence_level"`
	TailVaR            Money   `gorm:"embedded;embeddedPrefix:tail_var_" json:"tail_var"`
	ExpectedShortfall  Money   `gorm:"embedded;embeddedPrefix:expected_shortfall_" json:"expected_shortfall"`

	// Correlation & Dependencies
	CorrelationMatrix   datatypes.JSON `gorm:"type:json" json:"correlation_matrix"` // [][]float64
	DependencyStructure datatypes.JSON `gorm:"type:json" json:"dependency_structure"`
	SystemicRiskFactor  float64        `gorm:"type:decimal(10,4)" json:"systemic_risk_factor"`
	ContagionRisk       float64        `gorm:"type:decimal(5,2)" json:"contagion_risk"`

	// Concentration Risk
	ConcentrationRisk Money   `gorm:"embedded;embeddedPrefix:concentration_risk_" json:"concentration_risk"`
	SingleRiskLimit   Money   `gorm:"embedded;embeddedPrefix:single_risk_limit_" json:"single_risk_limit"`
	AggregateLimit    Money   `gorm:"embedded;embeddedPrefix:aggregate_limit_" json:"aggregate_limit"`
	LimitUtilization  float64 `gorm:"type:decimal(5,2)" json:"limit_utilization"`

	// Diversification
	DiversificationBenefit float64 `gorm:"type:decimal(5,2)" json:"diversification_benefit"`
	DiversificationIndex   float64 `gorm:"type:decimal(10,4)" json:"diversification_index"`
	EffectiveExposure      Money   `gorm:"embedded;embeddedPrefix:effective_exposure_" json:"effective_exposure"`

	// Status
	AccumulationStatus string    `gorm:"type:varchar(50)" json:"accumulation_status"`
	RiskLevel          string    `gorm:"type:varchar(50)" json:"risk_level"` // low, medium, high, critical
	LastReviewDate     time.Time `gorm:"type:timestamp" json:"last_review_date"`
	NextReviewDate     time.Time `gorm:"type:timestamp" json:"next_review_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// PolicyHierarchy represents policy relationships and hierarchy
type PolicyHierarchy struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Master/Sub Policy Structure
	MasterPolicyID *uuid.UUID     `gorm:"type:uuid" json:"master_policy_id,omitempty"`
	IsMasterPolicy bool           `gorm:"type:boolean;default:false" json:"is_master_policy"`
	SubPolicies    datatypes.JSON `gorm:"type:json" json:"sub_policies"` // []uuid.UUID
	SubPolicyCount int            `gorm:"type:int" json:"sub_policy_count"`

	// Linked Policies
	LinkedPolicies datatypes.JSON `gorm:"type:json" json:"linked_policies"`  // []uuid.UUID
	LinkType       string         `gorm:"type:varchar(50)" json:"link_type"` // family, bundle, corporate
	LinkReason     string         `gorm:"type:text" json:"link_reason"`

	// Dependent Policies
	DependentPolicies datatypes.JSON `gorm:"type:json" json:"dependent_policies"` // []uuid.UUID
	DependencyType    string         `gorm:"type:varchar(50)" json:"dependency_type"`
	DependencyRules   datatypes.JSON `gorm:"type:json" json:"dependency_rules"`

	// Cross-Product Links
	CrossProductLinks datatypes.JSON `gorm:"type:json" json:"cross_product_links"` // []CrossProduct
	CrossSellOrigin   *uuid.UUID     `gorm:"type:uuid" json:"cross_sell_origin,omitempty"`
	UpsellOrigin      *uuid.UUID     `gorm:"type:uuid" json:"upsell_origin,omitempty"`

	// Group Hierarchy
	PolicyGroupID  *uuid.UUID `gorm:"type:uuid" json:"policy_group_id,omitempty"`
	GroupName      string     `gorm:"type:varchar(255)" json:"group_name"`
	GroupRole      string     `gorm:"type:varchar(50)" json:"group_role"` // primary, secondary, member
	HierarchyLevel int        `gorm:"type:int" json:"hierarchy_level"`
	ParentLevel    *uuid.UUID `gorm:"type:uuid" json:"parent_level,omitempty"`

	// Inheritance Settings
	InheritsFromMaster   bool           `gorm:"type:boolean;default:false" json:"inherits_from_master"`
	InheritedAttributes  datatypes.JSON `gorm:"type:json" json:"inherited_attributes"`  // []string
	OverriddenAttributes datatypes.JSON `gorm:"type:json" json:"overridden_attributes"` // []string

	// Aggregation
	AggregatedLimits      bool           `gorm:"type:boolean;default:false" json:"aggregated_limits"`
	AggregatedDeductibles bool           `gorm:"type:boolean;default:false" json:"aggregated_deductibles"`
	SharedBenefits        datatypes.JSON `gorm:"type:json" json:"shared_benefits"`

	// Status
	HierarchyStatus string    `gorm:"type:varchar(50)" json:"hierarchy_status"`
	LastSyncDate    time.Time `gorm:"type:timestamp" json:"last_sync_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// PolicyExperienceRating represents experience-based rating adjustments
type PolicyExperienceRating struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Experience Modifier
	ExperienceModifier float64 `gorm:"type:decimal(10,4)" json:"experience_modifier"`
	BaseModifier       float64 `gorm:"type:decimal(10,4)" json:"base_modifier"`
	AdjustedModifier   float64 `gorm:"type:decimal(10,4)" json:"adjusted_modifier"`
	ModifierTrend      string  `gorm:"type:varchar(50)" json:"modifier_trend"` // improving, stable, deteriorating

	// Credibility
	CredibilityFactor       float64 `gorm:"type:decimal(5,4)" json:"credibility_factor"`
	CredibilityWeight       float64 `gorm:"type:decimal(5,4)" json:"credibility_weight"`
	MinimumCredibility      float64 `gorm:"type:decimal(5,4)" json:"minimum_credibility"`
	FullCredibilityStandard int     `gorm:"type:int" json:"full_credibility_standard"`

	// Loss History
	LossHistoryPeriod int            `gorm:"type:int" json:"loss_history_period_years"`
	HistoricalLosses  datatypes.JSON `gorm:"type:json" json:"historical_losses"` // []LossRecord
	IncurredLosses    Money          `gorm:"embedded;embeddedPrefix:incurred_losses_" json:"incurred_losses"`
	PaidLosses        Money          `gorm:"embedded;embeddedPrefix:paid_losses_" json:"paid_losses"`
	LossFrequency     float64        `gorm:"type:decimal(10,2)" json:"loss_frequency"`
	LossSeverity      Money          `gorm:"embedded;embeddedPrefix:loss_severity_" json:"loss_severity"`

	// Industry Comparison
	IndustryBenchmark   float64 `gorm:"type:decimal(10,4)" json:"industry_benchmark"`
	IndustryCode        string  `gorm:"type:varchar(50)" json:"industry_code"`
	PeerGroupComparison float64 `gorm:"type:decimal(5,2)" json:"peer_group_comparison"`
	RelativePerformance string  `gorm:"type:varchar(50)" json:"relative_performance"` // better, average, worse

	// Schedule Rating
	ScheduleRating        datatypes.JSON `gorm:"type:json" json:"schedule_rating"` // []RatingFactor
	ScheduleCredits       float64        `gorm:"type:decimal(5,2)" json:"schedule_credits"`
	ScheduleDebits        float64        `gorm:"type:decimal(5,2)" json:"schedule_debits"`
	NetScheduleAdjustment float64        `gorm:"type:decimal(5,2)" json:"net_schedule_adjustment"`

	// Retroactive Adjustments
	RetroactiveDate       *time.Time `gorm:"type:timestamp" json:"retroactive_date,omitempty"`
	RetroactiveAdjustment float64    `gorm:"type:decimal(10,2)" json:"retroactive_adjustment"`
	SunsetDate            *time.Time `gorm:"type:timestamp" json:"sunset_date,omitempty"`

	// Merit Rating
	MeritRating         float64        `gorm:"type:decimal(10,4)" json:"merit_rating"`
	MeritFactors        datatypes.JSON `gorm:"type:json" json:"merit_factors"`
	SafetyPrograms      datatypes.JSON `gorm:"type:json" json:"safety_programs"`
	LossControlMeasures datatypes.JSON `gorm:"type:json" json:"loss_control_measures"`

	// Calculation Details
	CalculationMethod string    `gorm:"type:varchar(100)" json:"calculation_method"`
	CalculationDate   time.Time `gorm:"type:timestamp" json:"calculation_date"`
	EffectiveDate     time.Time `gorm:"type:timestamp" json:"effective_date"`
	ExpirationDate    time.Time `gorm:"type:timestamp" json:"expiration_date"`

	// Status
	RatingStatus     string     `gorm:"type:varchar(50)" json:"rating_status"`
	ApprovalRequired bool       `gorm:"type:boolean;default:false" json:"approval_required"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovalDate     *time.Time `gorm:"type:timestamp" json:"approval_date,omitempty"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// PolicyOpportunities represents cross-sell and upsell opportunities
type PolicyOpportunities struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Cross-Sell
	CrossSellOpportunities datatypes.JSON `gorm:"type:json" json:"cross_sell_opportunities"` // []Opportunity
	CrossSellScore         float64        `gorm:"type:decimal(5,2)" json:"cross_sell_score"`
	CrossSellPriority      string         `gorm:"type:varchar(50)" json:"cross_sell_priority"`
	CrossSellReadiness     float64        `gorm:"type:decimal(5,2)" json:"cross_sell_readiness"`

	// Upsell
	UpsellOpportunities  datatypes.JSON `gorm:"type:json" json:"upsell_opportunities"` // []Opportunity
	UpsellScore          float64        `gorm:"type:decimal(5,2)" json:"upsell_score"`
	UpsellPriority       string         `gorm:"type:varchar(50)" json:"upsell_priority"`
	UpsellPotentialValue Money          `gorm:"embedded;embeddedPrefix:upsell_value_" json:"upsell_potential_value"`

	// Product Recommendations
	RecommendedProducts  datatypes.JSON `gorm:"type:json" json:"recommended_products"` // []ProductRecommendation
	RecommendationScore  float64        `gorm:"type:decimal(5,2)" json:"recommendation_score"`
	PersonalizationLevel float64        `gorm:"type:decimal(5,2)" json:"personalization_level"`

	// Customer Propensity
	CustomerPropensity datatypes.JSON `gorm:"type:json" json:"customer_propensity"` // map[string]float64
	PurchaseLikelihood float64        `gorm:"type:decimal(5,2)" json:"purchase_likelihood"`
	TimeToNextPurchase int            `gorm:"type:int" json:"time_to_next_purchase_days"`

	// Next Best Offer
	NextBestOffer       datatypes.JSON `gorm:"type:json" json:"next_best_offer"` // Offer
	OfferValidUntil     *time.Time     `gorm:"type:timestamp" json:"offer_valid_until,omitempty"`
	OfferAcceptanceProb float64        `gorm:"type:decimal(5,2)" json:"offer_acceptance_probability"`

	// Bundle Recommendations
	BundleRecommendations datatypes.JSON `gorm:"type:json" json:"bundle_recommendations"` // []Bundle
	BundleSavings         Money          `gorm:"embedded;embeddedPrefix:bundle_savings_" json:"bundle_savings"`
	BundleDiscount        float64        `gorm:"type:decimal(5,2)" json:"bundle_discount_percentage"`

	// Retention Offers
	RetentionOffers datatypes.JSON `gorm:"type:json" json:"retention_offers"` // []RetentionOffer
	RetentionRisk   float64        `gorm:"type:decimal(5,2)" json:"retention_risk"`
	RetentionValue  Money          `gorm:"embedded;embeddedPrefix:retention_value_" json:"retention_value"`

	// Campaign Tracking
	CampaignHistory      datatypes.JSON `gorm:"type:json" json:"campaign_history"`
	LastCampaignDate     *time.Time     `gorm:"type:timestamp" json:"last_campaign_date,omitempty"`
	CampaignResponseRate float64        `gorm:"type:decimal(5,2)" json:"campaign_response_rate"`

	// Status
	OpportunityStatus string    `gorm:"type:varchar(50)" json:"opportunity_status"`
	LastAnalysisDate  time.Time `gorm:"type:timestamp" json:"last_analysis_date"`
	NextReviewDate    time.Time `gorm:"type:timestamp" json:"next_review_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// PolicyPrivacy represents privacy and data protection settings
type PolicyPrivacy struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`

	// Consent Management
	ConsentRecords        datatypes.JSON `gorm:"type:json" json:"consent_records"` // []Consent
	DataProcessingConsent bool           `gorm:"type:boolean;default:false" json:"data_processing_consent"`
	MarketingConsent      bool           `gorm:"type:boolean;default:false" json:"marketing_consent"`
	ThirdPartyConsent     bool           `gorm:"type:boolean;default:false" json:"third_party_consent"`
	ConsentDate           time.Time      `gorm:"type:timestamp" json:"consent_date"`
	ConsentVersion        string         `gorm:"type:varchar(50)" json:"consent_version"`

	// Data Retention
	DataRetentionPolicy datatypes.JSON `gorm:"type:json" json:"data_retention_policy"` // RetentionPolicy
	RetentionPeriodDays int            `gorm:"type:int" json:"retention_period_days"`
	DataExpirationDate  *time.Time     `gorm:"type:timestamp" json:"data_expiration_date,omitempty"`
	AutoDeleteEnabled   bool           `gorm:"type:boolean;default:false" json:"auto_delete_enabled"`

	// Data Subject Rights
	DeletionRequests      datatypes.JSON `gorm:"type:json" json:"deletion_requests"`      // []DeletionRequest
	AccessRequests        datatypes.JSON `gorm:"type:json" json:"access_requests"`        // []AccessRequest
	PortabilityRequests   datatypes.JSON `gorm:"type:json" json:"portability_requests"`   // []PortabilityRequest
	RectificationRequests datatypes.JSON `gorm:"type:json" json:"rectification_requests"` // []RectificationRequest
	ObjectionRequests     datatypes.JSON `gorm:"type:json" json:"objection_requests"`     // []ObjectionRequest

	// Privacy Preferences
	PrivacyPreferences     datatypes.JSON `gorm:"type:json" json:"privacy_preferences"` // []Preference
	CommunicationOptOuts   datatypes.JSON `gorm:"type:json" json:"communication_opt_outs"`
	DataSharingPreferences datatypes.JSON `gorm:"type:json" json:"data_sharing_preferences"`

	// Anonymization
	AnonymizationStatus string     `gorm:"type:varchar(50)" json:"anonymization_status"`
	AnonymizationDate   *time.Time `gorm:"type:timestamp" json:"anonymization_date,omitempty"`
	AnonymizationMethod string     `gorm:"type:varchar(100)" json:"anonymization_method"`
	PseudonymizationKey string     `gorm:"type:varchar(255)" json:"pseudonymization_key"`

	// Data Breach
	DataBreachNotifications datatypes.JSON `gorm:"type:json" json:"data_breach_notifications"`
	LastBreachDate          *time.Time     `gorm:"type:timestamp" json:"last_breach_date,omitempty"`
	BreachNotificationSent  bool           `gorm:"type:boolean;default:false" json:"breach_notification_sent"`

	// Compliance
	GDPRCompliant       bool           `gorm:"type:boolean;default:false" json:"gdpr_compliant"`
	CCPACompliant       bool           `gorm:"type:boolean;default:false" json:"ccpa_compliant"`
	PrivacyRegulations  datatypes.JSON `gorm:"type:json" json:"privacy_regulations"` // []string
	LastComplianceCheck time.Time      `gorm:"type:timestamp" json:"last_compliance_check"`

	// Status
	PrivacyStatus  string    `gorm:"type:varchar(50)" json:"privacy_status"`
	LastUpdateDate time.Time `gorm:"type:timestamp" json:"last_update_date"`

	// Audit
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

// =====================================
// METHODS
// =====================================

// IsHighRisk checks if accumulation risk is high
func (pra *PolicyRiskAccumulation) IsHighRisk() bool {
	return pra.RiskLevel == "high" || pra.RiskLevel == "critical" ||
		pra.ConcentrationRisk.Amount > pra.SingleRiskLimit.Amount ||
		pra.LimitUtilization > 80
}

// GetPMLRatio calculates PML as percentage of total exposure
func (pra *PolicyRiskAccumulation) GetPMLRatio() float64 {
	if pra.CATExposure.Amount > 0 {
		return (pra.PMLCalculation.Amount / pra.CATExposure.Amount) * 100
	}
	return 0
}

// IsMaster checks if this is a master policy
func (ph *PolicyHierarchy) IsMaster() bool {
	return ph.IsMasterPolicy && ph.MasterPolicyID == nil
}

// HasSubPolicies checks if policy has sub-policies
func (ph *PolicyHierarchy) HasSubPolicies() bool {
	return ph.SubPolicyCount > 0
}

// IsLinked checks if policy is linked to others
func (ph *PolicyHierarchy) IsLinked() bool {
	return ph.LinkedPolicies != nil || ph.MasterPolicyID != nil
}

// IsCredible checks if experience rating is credible
func (per *PolicyExperienceRating) IsCredible() bool {
	return per.CredibilityFactor >= per.MinimumCredibility
}

// NeedsApproval checks if rating needs approval
func (per *PolicyExperienceRating) NeedsApproval() bool {
	return per.ApprovalRequired ||
		per.ExperienceModifier > 1.5 ||
		per.ExperienceModifier < 0.5
}

// HasOpportunities checks if there are sales opportunities
func (po *PolicyOpportunities) HasOpportunities() bool {
	return po.CrossSellScore > 50 || po.UpsellScore > 50
}

// IsHighValue checks if customer is high value for opportunities
func (po *PolicyOpportunities) IsHighValue() bool {
	return po.UpsellPotentialValue.Amount > 5000 ||
		po.PurchaseLikelihood > 70
}

// HasPrivacyConsent checks if privacy consent is given
func (pp *PolicyPrivacy) HasPrivacyConsent() bool {
	return pp.DataProcessingConsent && pp.ConsentRecords != nil
}

// RequiresDeletion checks if data deletion is required
func (pp *PolicyPrivacy) RequiresDeletion() bool {
	if pp.DataExpirationDate != nil {
		return time.Now().After(*pp.DataExpirationDate) || pp.DeletionRequests != nil
	}
	return pp.DeletionRequests != nil
}
