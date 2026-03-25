package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"smartsure/internal/domain/models/user"
	"smartsure/pkg/database"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// User represents a comprehensive insurance customer with all necessary features
type User struct {
	// Base fields
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy      *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID     `gorm:"type:uuid" json:"updated_by,omitempty"`
	ProfilePicture string         `json:"profile_picture"`

	// Embedded structs for organization
	user.UserIdentification   `gorm:"embedded"`
	user.UserAuthentication   `gorm:"embedded"`
	user.UserAddress          `gorm:"embedded"`
	user.UserInsuranceProfile `gorm:"embedded"`
	user.UserFinancial        `gorm:"embedded"`
	user.UserRiskAssessment   `gorm:"embedded"`
	user.UserCompliance       `gorm:"embedded"`
	user.UserEngagement       `gorm:"embedded"`
	user.UserSupport          `gorm:"embedded"`
	user.UserDigital          `gorm:"embedded"`
	user.UserMarketing        `gorm:"embedded"`
	user.UserHousehold        `gorm:"embedded"`
	user.UserMetadata         `gorm:"embedded"`

	// Core Entity Relationships
	Devices  []Device  `gorm:"foreignKey:OwnerID" json:"devices,omitempty"`
	Policies []Policy  `gorm:"foreignKey:CustomerID" json:"policies,omitempty"`
	Claims   []Claim   `gorm:"foreignKey:CustomerID" json:"claims,omitempty"`
	Payments []Payment `gorm:"foreignKey:UserID" json:"payments,omitempty"`

	// Analytics & Intelligence Models (One-to-One)
	Analytics           *user.UserAnalytics           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"analytics,omitempty"`
	BehaviorAnalytics   *user.UserBehaviorAnalytics   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"behavior_analytics,omitempty"`
	PredictiveModeling  *user.UserPredictiveModeling  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"predictive_modeling,omitempty"`
	CohortAnalysis      *user.UserCohortAnalysis      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"cohort_analysis,omitempty"`
	AttributionModeling *user.UserAttributionModeling `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"attribution_modeling,omitempty"`
	LifecycleAnalytics  *user.UserLifecycleAnalytics  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"lifecycle_analytics,omitempty"`

	// Fraud & Security Models (One-to-One)
	FraudDetection     *user.UserFraudDetection     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"fraud_detection,omitempty"`
	SecurityProfile    *user.UserSecurityProfile    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"security_profile,omitempty"`
	PrivacySettings    *user.UserPrivacySettings    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"privacy_settings,omitempty"`
	AccessControl      *user.UserAccessControl      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"access_control,omitempty"`
	Biometrics         *user.UserBiometrics         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"biometrics,omitempty"`
	Encryption         *user.UserEncryption         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"encryption,omitempty"`
	ThreatIntelligence *user.UserThreatIntelligence `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"threat_intelligence,omitempty"`

	// Workflow & Operations Models (One-to-One)
	Workflow           *user.UserWorkflow           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"workflow,omitempty"`
	ProcessAutomation  *user.UserProcessAutomation  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"process_automation,omitempty"`
	SLA                *user.UserSLA                `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"sla,omitempty"`
	PerformanceMetrics *user.UserPerformanceMetrics `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"performance_metrics,omitempty"`
	QualityAssurance   *user.UserQualityAssurance   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"quality_assurance,omitempty"`
	CapacityManagement *user.UserCapacityManagement `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"capacity_management,omitempty"`
	Scheduling         *user.UserScheduling         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"scheduling,omitempty"`
	ServiceCatalog     *user.UserServiceCatalog     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"service_catalog,omitempty"`

	// Communication Models (One-to-One)
	CommunicationHistory *user.UserCommunicationHistory    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"communication_history,omitempty"`
	NotificationPrefs    *user.UserNotificationPreferences `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"notification_preferences,omitempty"`
	Sentiment            *user.UserSentiment               `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"sentiment,omitempty"`

	// Financial Models (One-to-One)
	FinancialProfile  *user.UserFinancialProfile  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"financial_profile,omitempty"`
	CreditProfile     *user.UserCreditProfile     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"credit_profile,omitempty"`
	PaymentBehavior   *user.UserPaymentBehavior   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"payment_behavior,omitempty"`
	BillingProfile    *user.UserBillingProfile    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"billing_profile,omitempty"`
	TaxProfile        *user.UserTaxProfile        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"tax_profile,omitempty"`
	InvestmentProfile *user.UserInvestmentProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"investment_profile,omitempty"`

	// Ecosystem Models (One-to-One)
	DeviceEcosystem  *user.UserDeviceEcosystem  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"device_ecosystem,omitempty"`
	PolicyPortfolio  *user.UserPolicyPortfolio  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"policy_portfolio,omitempty"`
	ClaimPatterns    *user.UserClaimPatterns    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"claim_patterns,omitempty"`
	PaymentEcosystem *user.UserPaymentEcosystem `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"payment_ecosystem,omitempty"`
	ReferralNetwork  *user.UserReferralNetwork  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"referral_network,omitempty"`
	LoyaltyJourney   *user.UserLoyaltyJourney   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"loyalty_journey,omitempty"`

	// Compliance Models (One-to-One)
	ComplianceProfile     *user.UserComplianceProfile     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"compliance_profile,omitempty"`
	CrossBorderCompliance *user.UserCrossBorderCompliance `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"cross_border_compliance,omitempty"`
	DataRetention         *user.UserDataRetention         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"data_retention,omitempty"`
	ComplianceReporting   *user.UserComplianceReporting   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"compliance_reporting,omitempty"`

	// Experience Models (One-to-One)
	JourneyMapping  *user.UserJourneyMapping  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"journey_mapping,omitempty"`
	Onboarding      *user.UserOnboarding      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"onboarding,omitempty"`
	Personalization *user.UserPersonalization `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"personalization,omitempty"`
	Retention       *user.UserRetention       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"retention,omitempty"`
	Segmentation    *user.UserSegmentation    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"segmentation,omitempty"`
	PredictiveNeeds *user.UserPredictiveNeeds `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"predictive_needs,omitempty"`

	// Gamification & Social Models (One-to-One)
	Gamification   *user.UserGamification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"gamification,omitempty"`
	Education      *user.UserEducation    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"education,omitempty"`
	Social         *user.UserSocial       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"social,omitempty"`
	RewardsProgram *user.UserRewards      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"rewards_program,omitempty"`

	// Integration Models (One-to-One)
	Integrations        *user.UserIntegrations        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"integrations,omitempty"`
	APIUsage            *user.UserAPIUsage            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"api_usage,omitempty"`
	ThirdPartyData      *user.UserThirdPartyData      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"third_party_data,omitempty"`
	DataExchange        *user.UserDataExchange        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"data_exchange,omitempty"`
	PartnerIntegrations *user.UserPartnerIntegrations `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"partner_integrations,omitempty"`
	AutomationRules     *user.UserAutomationRules     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"automation_rules,omitempty"`

	// Underwriting Models (One-to-One)
	Underwriting       *user.UserUnderwriting       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"underwriting,omitempty"`
	RiskEvolution      *user.UserRiskEvolution      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"risk_evolution,omitempty"`
	Reserves           *user.UserReserves           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"reserves,omitempty"`
	UnderwritingModels *user.UserUnderwritingModels `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"underwriting_models,omitempty"`
	PricingModel       *user.UserPricingModel       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"pricing_model,omitempty"`

	// Operational Models (One-to-One)
	MultiCurrency       *user.UserMultiCurrency       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"multi_currency,omitempty"`
	EnvironmentalImpact *user.UserEnvironmentalImpact `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"environmental_impact,omitempty"`

	// Existing One-to-Many Relationships
	PaymentMethods     []UserPaymentMethod    `gorm:"foreignKey:UserID" json:"payment_methods,omitempty"`
	InsuranceHistories []UserInsuranceHistory `gorm:"foreignKey:UserID" json:"insurance_histories,omitempty"`
	SupportTickets     []UserSupportTicket    `gorm:"foreignKey:UserID" json:"support_tickets,omitempty"`
	NotificationLogs   []UserNotificationLog  `gorm:"foreignKey:UserID" json:"notification_logs,omitempty"`
	AuditLogs          []user.UserAuditLog    `gorm:"foreignKey:UserID" json:"audit_logs,omitempty"`
	SecurityEvents     []UserSecurityEvent    `gorm:"foreignKey:UserID" json:"security_events,omitempty"`
	Preferences        []UserPreference       `gorm:"foreignKey:UserID" json:"preferences,omitempty"`
	Documents          []UserDocument         `gorm:"foreignKey:UserID" json:"documents,omitempty"`
	Activities         []UserActivity         `gorm:"foreignKey:UserID" json:"activities,omitempty"`
	Rewards            []UserReward           `gorm:"foreignKey:UserID" json:"rewards,omitempty"`

	// Workflow & Task Management (One-to-Many)
	ApprovalRequests []user.UserApprovalRequest `gorm:"foreignKey:UserID" json:"approval_requests,omitempty"`
	ApprovalHistory  []user.UserApprovalHistory `gorm:"foreignKey:UserID" json:"approval_history,omitempty"`
	Escalations      []user.UserEscalation      `gorm:"foreignKey:UserID" json:"escalations,omitempty"`
	TaskManagement   []user.UserTaskManagement  `gorm:"foreignKey:UserID" json:"task_management,omitempty"`

	// Communication (One-to-Many)
	CommunicationLogs   []user.UserCommunicationLog   `gorm:"foreignKey:UserID" json:"communication_logs,omitempty"`
	CampaignEngagements []user.UserCampaignEngagement `gorm:"foreignKey:UserID" json:"campaign_engagements,omitempty"`
	Conversations       []user.UserConversation       `gorm:"foreignKey:UserID" json:"conversations,omitempty"`
	MessageTemplates    []user.UserMessageTemplate    `gorm:"foreignKey:UserID" json:"message_templates,omitempty"`

	// Compliance (One-to-Many)
	RegulatoryFilings []user.UserRegulatoryFiling `gorm:"foreignKey:UserID" json:"regulatory_filings,omitempty"`
	LegalHolds        []user.UserLegalHolds       `gorm:"foreignKey:UserID" json:"legal_holds,omitempty"`

	// Fraud & Security (One-to-Many)
	AnomalyDetections     []user.UserAnomalyDetection     `gorm:"foreignKey:UserID" json:"anomaly_detections,omitempty"`
	IdentityVerifications []user.UserIdentityVerification `gorm:"foreignKey:UserID" json:"identity_verifications,omitempty"`
	FraudInvestigations   []user.UserFraudInvestigation   `gorm:"foreignKey:UserID" json:"fraud_investigations,omitempty"`
	BlacklistManagement   *user.UserBlacklistManagement   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"blacklist_management,omitempty"`
	FraudPrevention       *user.UserFraudPrevention       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"fraud_prevention,omitempty"`

	// Experience & Feedback (One-to-Many)
	Feedback   []user.UserFeedback   `gorm:"foreignKey:UserID" json:"feedback,omitempty"`
	Challenges []user.UserChallenges `gorm:"foreignKey:UserID" json:"challenges,omitempty"`
	Missions   *user.UserMissions    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"missions,omitempty"`

	// Underwriting History (One-to-Many)
	UnderwritingHistory []user.UserUnderwritingHistory `gorm:"foreignKey:UserID" json:"underwriting_history,omitempty"`

	// Incidents (One-to-Many)
	IncidentReports []user.UserIncidentManagement `gorm:"foreignKey:UserID" json:"incident_reports,omitempty"`

	// Family relationships
	FamilyMembers []User `gorm:"foreignKey:HouseholdID" json:"family_members,omitempty"`
	ReferralsMade []User `gorm:"foreignKey:ReferredBy" json:"referrals_made,omitempty"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate handles pre-creation logic
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Set created_by from context if available
	if userID := database.GetUserIDFromContext(tx.Statement.Context); userID != uuid.Nil {
		u.CreatedBy = &userID
	}

	// Generate referral code if not set
	if u.ReferralCode == "" {
		u.ReferralCode = u.GenerateReferralCode()
	}

	// Set default values
	if u.CustomerSegment == "" {
		u.CustomerSegment = "individual"
	}
	if u.SubscriptionTier == "" {
		u.SubscriptionTier = "basic"
	}
	if u.PreferredBillingCycle == "" {
		u.PreferredBillingCycle = "monthly"
	}

	return nil
}

// BeforeUpdate handles pre-update logic
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Set updated_by from context if available
	if userID := database.GetUserIDFromContext(tx.Statement.Context); userID != uuid.Nil {
		u.UpdatedBy = &userID
	}

	// Update risk scores if significant changes
	if tx.Statement.Changed("CreditScore", "ClaimFrequency", "PaymentFailureCount") {
		u.RecalculateRiskScore()
	}

	return nil
}

// === Core Identification Methods ===

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	if u.MiddleName != "" {
		return fmt.Sprintf("%s %s %s", u.FirstName, u.MiddleName, u.LastName)
	}
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// GetAge returns the age of the user
func (u *User) GetAge() int {
	if u.DateOfBirth == nil {
		return 0
	}
	now := time.Now()
	years := now.Year() - u.DateOfBirth.Year()
	if now.YearDay() < u.DateOfBirth.YearDay() {
		years--
	}
	return years
}

// === Verification & Compliance Methods ===

// IsVerified checks if the user is fully verified
func (u *User) IsVerified() bool {
	return u.EmailVerified && u.PhoneVerified && u.DocumentsVerified && u.KYCStatus == "approved"
}

// HasCompletedKYC checks if KYC is complete
func (u *User) HasCompletedKYC() bool {
	return u.KYCStatus == "approved" && u.KYCLevel != "basic"
}

// RequiresAdditionalVerification checks if more verification is needed
func (u *User) RequiresAdditionalVerification() bool {
	// High-risk users or high-value customers need enhanced verification
	return u.RiskScore > 70 || u.VIPStatus || u.TotalPremiumPaid.GreaterThan(decimal.NewFromInt(10000))
}

// IsCompliant checks overall compliance status
func (u *User) IsCompliant() bool {
	return u.KYCStatus == "approved" &&
		u.AMLStatus != "failed" &&
		!u.PEPStatus &&
		u.SanctionsScreened &&
		u.GDPRConsent
}

// === Insurance Eligibility Methods ===

// CanPurchaseInsurance checks if user can purchase insurance
func (u *User) CanPurchaseInsurance() bool {
	return u.Status == "active" &&
		u.IsVerified() &&
		!u.BlacklistStatus &&
		u.RiskScore < 85 &&
		!u.IsAccountLocked() &&
		!u.IsSuspended()
}

// GetMaxDevicesAllowed returns max devices user can insure
func (u *User) GetMaxDevicesAllowed() int {
	base := u.MaxDevicesAllowed
	if u.VIPStatus {
		base += 5
	}
	if u.CustomerSegment == "family" {
		base = base * u.FamilyMemberCount
	}
	if u.CustomerSegment == "corporate" {
		base = 999 // Unlimited for corporate
	}
	return base
}

// CanAddMoreDevices checks if user can add more devices
func (u *User) CanAddMoreDevices() bool {
	return u.CurrentDeviceCount < u.GetMaxDevicesAllowed()
}

// === Risk Assessment Methods ===

// IsHighRisk determines if user is high risk
func (u *User) IsHighRisk() bool {
	return u.RiskScore > 75 ||
		u.FraudScore > 60 ||
		u.ClaimFrequency > 3 ||
		u.DefaultRisk > 0.5
}

// GetRiskCategory returns risk category
func (u *User) GetRiskCategory() string {
	switch {
	case u.RiskScore < 30:
		return "low"
	case u.RiskScore < 60:
		return "medium"
	case u.RiskScore < 80:
		return "high"
	default:
		return "very_high"
	}
}

// RecalculateRiskScore recalculates risk score based on various factors
func (u *User) RecalculateRiskScore() float64 {
	score := 50.0

	// Credit score impact
	if u.CreditScore != nil {
		if *u.CreditScore < 500 {
			score += 20
		} else if *u.CreditScore > 700 {
			score -= 10
		}
	}

	// Claims history impact
	score += u.ClaimFrequency * 5
	if u.ClaimApprovalRate > 0.8 {
		score += 10
	}

	// Payment history impact
	score += float64(u.PaymentFailureCount) * 3
	if u.PaymentFailureCount == 0 && u.TotalPremiumPaid.GreaterThan(decimal.NewFromInt(5000)) {
		score -= 5
	}

	// Fraud indicators
	score += u.FraudScore * 0.5
	score += float64(u.SuspiciousActivityCount) * 2

	// Compliance impact
	if !u.IsVerified() {
		score += 15
	}
	if u.BlacklistStatus {
		score += 30
	}

	// Loyalty discount
	if u.FirstPolicyDate != nil {
		yearsAsCustomer := time.Since(*u.FirstPolicyDate).Hours() / 24 / 365
		score -= min(yearsAsCustomer*2, 10)
	}

	// Normalize score
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	u.RiskScore = score
	return score
}

// GetFraudRisk returns fraud risk level
func (u *User) GetFraudRisk() string {
	if u.FraudScore > 70 || u.BlacklistStatus {
		return "high"
	} else if u.FraudScore > 40 || u.SuspiciousActivityCount > 3 {
		return "medium"
	}
	return "low"
}

// === Financial Methods ===

// GetCreditLimit returns user's credit limit for premium financing
func (u *User) GetCreditLimit() decimal.Decimal {
	if u.CreditLimit.IsPositive() {
		return u.CreditLimit
	}

	// Calculate based on credit score and history
	baseLimit := decimal.NewFromInt(1000)
	if u.CreditScore != nil {
		multiplier := decimal.NewFromFloat(float64(*u.CreditScore) / 100)
		baseLimit = baseLimit.Mul(multiplier)
	}

	// Adjust based on payment history
	if u.PaymentFailureCount == 0 {
		baseLimit = baseLimit.Mul(decimal.NewFromFloat(1.5))
	}

	// VIP bonus
	if u.VIPStatus {
		baseLimit = baseLimit.Mul(decimal.NewFromInt(2))
	}

	return baseLimit
}

// HasOutstandingBalance checks if user has unpaid balance
func (u *User) HasOutstandingBalance() bool {
	return u.OutstandingBalance.IsPositive()
}

// IsInGoodStanding checks financial standing
func (u *User) IsInGoodStanding() bool {
	return !u.HasOutstandingBalance() &&
		u.PaymentFailureCount < 3 &&
		u.DefaultRisk < 0.3
}

// === Customer Segmentation Methods ===

// GetCustomerSegment returns customer segment
func (u *User) GetCustomerSegment() string {
	if u.CustomerSegment != "" {
		return u.CustomerSegment
	}

	// Auto-determine segment
	if u.CorporateAccountID != nil {
		return "corporate"
	}
	if u.FamilyMemberCount > 1 {
		return "family"
	}
	if u.StudentVerified {
		return "student"
	}
	if u.SeniorCitizen {
		return "senior"
	}
	return "individual"
}

// CalculateLoyaltyTier calculates loyalty tier based on history
func (u *User) CalculateLoyaltyTier() string {
	if u.FirstPolicyDate == nil {
		return "bronze"
	}

	years := time.Since(*u.FirstPolicyDate).Hours() / 24 / 365
	totalValue := u.TotalPremiumPaid.InexactFloat64()

	if years > 5 && totalValue > 50000 {
		return "diamond"
	} else if years > 3 && totalValue > 25000 {
		return "platinum"
	} else if years > 2 && totalValue > 10000 {
		return "gold"
	} else if years > 1 && totalValue > 5000 {
		return "silver"
	}
	return "bronze"
}

// IsVIP checks if user is a VIP customer
func (u *User) IsVIP() bool {
	return u.VIPStatus ||
		u.CalculateLoyaltyTier() == "diamond" ||
		u.CustomerLifetimeValue.GreaterThan(decimal.NewFromInt(100000))
}

// === Claims Methods ===

// CanFileClaimToday checks if user can file a new claim
func (u *User) CanFileClaimToday() bool {
	if u.LastClaimDate != nil {
		daysSinceLastClaim := time.Since(*u.LastClaimDate).Hours() / 24
		if daysSinceLastClaim < 30 && u.ClaimFrequency > 2 {
			return false // Too many claims recently
		}
	}

	return u.CanPurchaseInsurance() &&
		u.ActivePoliciesCount > 0 &&
		!u.IsHighRisk()
}

// HasOutstandingClaims checks for pending claims
func (u *User) HasOutstandingClaims(db *gorm.DB) bool {
	var count int64
	db.Model(&Claim{}).Where("customer_id = ? AND status IN ?", u.ID, []string{"pending", "processing", "investigating"}).Count(&count)
	return count > 0
}

// GetClaimHistory returns claim statistics
func (u *User) GetClaimHistory() map[string]interface{} {
	return map[string]interface{}{
		"total_claims":    u.TotalClaimsCount,
		"approved_claims": u.ApprovedClaimsCount,
		"approval_rate":   u.ClaimApprovalRate,
		"average_amount":  u.AverageClaimAmount,
		"last_claim_date": u.LastClaimDate,
		"claim_frequency": u.ClaimFrequency,
	}
}

// === Discount & Rewards Methods ===

// GetDiscountPercentage calculates total discount
func (u *User) GetDiscountPercentage() float64 {
	discount := 0.0

	// Loyalty discount
	switch u.CalculateLoyaltyTier() {
	case "diamond":
		discount += 20
	case "platinum":
		discount += 15
	case "gold":
		discount += 10
	case "silver":
		discount += 5
	}

	// Family discount
	if u.CustomerSegment == "family" && u.FamilyMemberCount > 2 {
		discount += u.FamilyDiscountRate
	}

	// Student/Senior discount
	if u.StudentVerified {
		discount += 10
	}
	if u.SeniorCitizen {
		discount += 15
	}

	// Military discount
	if u.MilitaryService {
		discount += 10
	}

	// Good standing bonus
	if u.IsInGoodStanding() && u.PaymentFailureCount == 0 {
		discount += 5
	}

	// Cap at 40% max discount
	if discount > 40 {
		discount = 40
	}

	return discount
}

// GetRenewalDiscount calculates renewal discount
func (u *User) GetRenewalDiscount() float64 {
	baseDiscount := u.GetDiscountPercentage()

	// Additional renewal loyalty
	if u.AutoRenewalEnabled {
		baseDiscount += 5
	}

	// No claims bonus
	if u.LastClaimDate == nil || time.Since(*u.LastClaimDate).Hours()/24 > 365 {
		baseDiscount += 10
	}

	return min(baseDiscount, 45) // Cap at 45% for renewals
}

// === Service Level Methods ===

// GetServiceLevel returns customer service tier
func (u *User) GetServiceLevel() string {
	if u.ServiceTier != "" {
		return u.ServiceTier
	}

	if u.IsVIP() || u.CustomerSegment == "corporate" {
		return "vip"
	}
	if u.CalculateLoyaltyTier() == "gold" || u.CalculateLoyaltyTier() == "platinum" {
		return "priority"
	}
	return "standard"
}

// IsEligibleForUpgrade checks upgrade eligibility
func (u *User) IsEligibleForUpgrade() bool {
	return u.UpgradeEligibility ||
		(u.SubscriptionTier != "platinum" &&
			u.IsInGoodStanding() &&
			u.EngagementScore > 70)
}

// === Retention & Analytics Methods ===

// GetRetentionRisk returns retention risk level
func (u *User) GetRetentionRisk() string {
	if u.ChurnRisk > 0.7 {
		return "high"
	} else if u.ChurnRisk > 0.4 {
		return "medium"
	}
	return "low"
}

// CalculateLifetimeValue calculates CLV
func (u *User) CalculateLifetimeValue() decimal.Decimal {
	if !u.CustomerLifetimeValue.IsZero() {
		return u.CustomerLifetimeValue
	}

	// Simple CLV calculation
	avgPremiumPerYear := u.TotalPremiumPaid.Div(decimal.NewFromFloat(max(1, time.Since(u.CreatedAt).Hours()/24/365)))
	expectedLifespan := decimal.NewFromFloat(5 - u.ChurnRisk*4) // 1-5 years based on churn risk

	clv := avgPremiumPerYear.Mul(expectedLifespan)

	// Adjust for referrals
	referralValue := decimal.NewFromInt(int64(u.ReferralCount * 500))
	clv = clv.Add(referralValue)

	u.CustomerLifetimeValue = clv
	return clv
}

// === Communication Methods ===

// GetPreferredCommunicationChannel returns preferred channel
func (u *User) GetPreferredCommunicationChannel() string {
	if u.PreferredContactMethod != "" {
		return u.PreferredContactMethod
	}

	// Default based on opt-ins
	if u.OptInWhatsApp {
		return "whatsapp"
	} else if u.OptInSMS {
		return "sms"
	} else if u.OptInPush && u.PushToken != "" {
		return "push"
	}
	return "email"
}

// ShouldReceiveMarketing checks marketing eligibility
func (u *User) ShouldReceiveMarketing() bool {
	return u.MarketingConsent &&
		!u.DoNotDisturb &&
		u.Status == "active" &&
		!u.IsSuspended()
}

// === Authentication Methods ===

// IsAccountLocked checks if the account is currently locked
func (u *User) IsAccountLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// IsSuspended checks if the account is suspended
func (u *User) IsSuspended() bool {
	return u.SuspendedAt != nil
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.LoginAttempts = 0 // Reset login attempts on successful login
	u.LockedUntil = nil // Clear any account lock
}

// IncrementLoginAttempts increments failed login attempts
func (u *User) IncrementLoginAttempts() {
	u.LoginAttempts++
	if u.LoginAttempts >= 5 {
		// Lock account for 30 minutes after 5 failed attempts
		lockUntil := time.Now().Add(30 * time.Minute)
		u.LockedUntil = &lockUntil
	}
}

// GenerateReferralCode generates a unique referral code
func (u *User) GenerateReferralCode() string {
	// Generate a random 8-character code
	b := make([]byte, 4)
	rand.Read(b)
	randomPart := hex.EncodeToString(b)

	// Create code with user initials and random part
	firstInitial := strings.ToUpper(string(u.FirstName[0]))
	lastInitial := strings.ToUpper(string(u.LastName[0]))

	return fmt.Sprintf("%s%s%s", firstInitial, lastInitial, strings.ToUpper(randomPart))
}

// === Validation Methods ===

// Validate performs comprehensive validation
func (u *User) Validate() error {
	// Required fields
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.FirstName == "" || u.LastName == "" {
		return fmt.Errorf("first and last name are required")
	}
	if u.PhoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}

	// Age validation for insurance
	age := u.GetAge()
	if age > 0 && age < 18 && !u.IsHouseholdHead {
		return fmt.Errorf("must be 18 or older to purchase insurance independently")
	}

	// Risk validation
	if u.RiskScore > 95 {
		return fmt.Errorf("risk score too high for insurance eligibility")
	}

	// Blacklist check
	if u.BlacklistStatus {
		return fmt.Errorf("user is blacklisted: %s", u.BlacklistReason)
	}

	return nil
}

// Helper function for min
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Helper function for max
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
