package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// UserIdentification contains core identification fields
type UserIdentification struct {
	Email          string     `gorm:"uniqueIndex;not null" json:"email"`
	PhoneNumber    string     `gorm:"uniqueIndex" json:"phone_number"`
	FirstName      string     `gorm:"not null" json:"first_name"`
	LastName       string     `gorm:"not null" json:"last_name"`
	MiddleName     string     `json:"middle_name"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Gender         string     `gorm:"type:varchar(10)" json:"gender"`
	Nationality    string     `json:"nationality"`
	PassportNumber string     `gorm:"index" json:"passport_number"`
	NationalID     string     `gorm:"index" json:"national_id"`
	SocialSecurity string     `gorm:"index" json:"social_security"`
	TaxID          string     `gorm:"index" json:"tax_id"`
	DriversLicense string     `gorm:"index" json:"drivers_license"`
}

// UserAuthentication contains authentication and security fields
type UserAuthentication struct {
	PasswordHash        string     `gorm:"not null" json:"-"`
	TwoFactorEnabled    bool       `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret     string     `json:"-"`
	BiometricEnabled    bool       `gorm:"default:false" json:"biometric_enabled"`
	BiometricData       string     `gorm:"type:json" json:"-"`
	LastLoginAt         *time.Time `json:"last_login_at"`
	LastPasswordChange  *time.Time `json:"last_password_change"`
	PasswordResetToken  string     `json:"-"`
	PasswordResetExpiry *time.Time `json:"-"`
	LoginAttempts       int        `gorm:"default:0" json:"login_attempts"`
	LockedUntil         *time.Time `json:"locked_until"`
	IPWhitelist         string     `gorm:"type:json" json:"ip_whitelist"`
	DeviceFingerprintID string     `json:"device_fingerprint_id"`
	SecurityQuestions   string     `gorm:"type:json" json:"-"`
	PINCode             string     `json:"-"`
	SessionToken        string     `json:"-"`
	SessionExpiry       *time.Time `json:"-"`
}

// UserAddress contains location and address information
type UserAddress struct {
	Country            string  `json:"country"`
	StateProvince      string  `json:"state_province"`
	City               string  `json:"city"`
	District           string  `json:"district"`
	Address            string  `json:"address"`
	AddressLine2       string  `json:"address_line2"`
	PostalCode         string  `json:"postal_code"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	TimeZone           string  `gorm:"default:'UTC'" json:"timezone"`
	BillingAddress     string  `gorm:"type:json" json:"billing_address"`
	ShippingAddress    string  `gorm:"type:json" json:"shipping_address"`
	PrimaryResidence   string  `json:"primary_residence"`
	SecondaryResidence string  `json:"secondary_residence"`
	CoverageZone       string  `json:"coverage_zone"`
	IsInternational    bool    `gorm:"default:false" json:"is_international"`
}

// UserInsuranceProfile contains insurance-specific information
type UserInsuranceProfile struct {
	CustomerSegment     string          `gorm:"type:varchar(30)" json:"customer_segment"`  // individual/family/corporate/student/senior
	SubscriptionTier    string          `gorm:"type:varchar(20)" json:"subscription_tier"` // basic/standard/premium/platinum
	LoyaltyTier         string          `gorm:"type:varchar(20)" json:"loyalty_tier"`      // bronze/silver/gold/platinum/diamond
	InsuranceExperience int             `json:"insurance_experience_years"`
	PreviousInsurer     string          `json:"previous_insurer"`
	FirstPolicyDate     *time.Time      `json:"first_policy_date"`
	TotalPoliciesCount  int             `gorm:"default:0" json:"total_policies_count"`
	ActivePoliciesCount int             `gorm:"default:0" json:"active_policies_count"`
	TotalClaimsCount    int             `gorm:"default:0" json:"total_claims_count"`
	ApprovedClaimsCount int             `gorm:"default:0" json:"approved_claims_count"`
	ClaimApprovalRate   float64         `gorm:"default:0" json:"claim_approval_rate"`
	AverageClaimAmount  decimal.Decimal `gorm:"type:decimal(15,2)" json:"average_claim_amount"`
	LastClaimDate       *time.Time      `json:"last_claim_date"`
	LastPolicyReview    *time.Time      `json:"last_policy_review"`
	UnderwritingScore   float64         `gorm:"default:50.0" json:"underwriting_score"`
	MaxDevicesAllowed   int             `gorm:"default:5" json:"max_devices_allowed"`
	CurrentDeviceCount  int             `gorm:"default:0" json:"current_device_count"`
	PreferredDeductible decimal.Decimal `gorm:"type:decimal(10,2)" json:"preferred_deductible"`
	PreferredCoverage   string          `gorm:"type:json" json:"preferred_coverage"`
}

// UserFinancial contains financial and payment information
type UserFinancial struct {
	CreditScore           *int            `json:"credit_score"`
	CreditCheckDate       *time.Time      `json:"credit_check_date"`
	CreditLimit           decimal.Decimal `gorm:"type:decimal(15,2)" json:"credit_limit"`
	OutstandingBalance    decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"outstanding_balance"`
	TotalPremiumPaid      decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"total_premium_paid"`
	PaymentFailureCount   int             `gorm:"default:0" json:"payment_failure_count"`
	LastPaymentDate       *time.Time      `json:"last_payment_date"`
	NextPaymentDate       *time.Time      `json:"next_payment_date"`
	PreferredPaymentDay   int             `json:"preferred_payment_day"`
	PreferredBillingCycle string          `gorm:"type:varchar(20)" json:"preferred_billing_cycle"` // monthly/quarterly/annual
	AutoRenewalEnabled    bool            `gorm:"default:true" json:"auto_renewal_enabled"`
	PaymentMethod         string          `gorm:"type:varchar(20)" json:"payment_method"` // card/bank/wallet/crypto
	BankAccountVerified   bool            `gorm:"default:false" json:"bank_account_verified"`
	TaxExempt             bool            `gorm:"default:false" json:"tax_exempt"`
	TaxExemptionNumber    string          `json:"tax_exemption_number"`
	DiscountCodes         string          `gorm:"type:json" json:"discount_codes"`
	TotalDiscountReceived decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"total_discount_received"`
	StripeCustomerID      string          `gorm:"index" json:"stripe_customer_id"` // For Stripe integration
}

// UserRiskAssessment contains risk and fraud assessment data
type UserRiskAssessment struct {
	RiskScore               float64    `gorm:"default:50.0" json:"risk_score"`
	FraudScore              float64    `gorm:"default:0" json:"fraud_score"`
	BlacklistStatus         bool       `gorm:"default:false" json:"blacklist_status"`
	BlacklistReason         string     `json:"blacklist_reason"`
	SuspiciousActivityCount int        `gorm:"default:0" json:"suspicious_activity_count"`
	LastRiskAssessment      *time.Time `json:"last_risk_assessment"`
	RiskFactors             string     `gorm:"type:json" json:"risk_factors"`
	ClaimFrequency          float64    `json:"claim_frequency"`
	LossRatio               float64    `json:"loss_ratio"`
	ChurnRisk               float64    `gorm:"default:0" json:"churn_risk"`
	DefaultRisk             float64    `gorm:"default:0" json:"default_risk"`
	FraudAlerts             string     `gorm:"type:json" json:"fraud_alerts"`
	InvestigationStatus     string     `json:"investigation_status"`
	RiskMitigationActions   string     `gorm:"type:json" json:"risk_mitigation_actions"`
}

// UserCompliance contains KYC, AML and regulatory compliance data
type UserCompliance struct {
	UserType               string     `gorm:"type:varchar(20);not null;default:'customer'" json:"user_type"`
	Status                 string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	KYCStatus              string     `gorm:"type:varchar(20);not null;default:'not_started'" json:"kyc_status"`
	KYCLevel               string     `gorm:"type:varchar(20);not null;default:'basic'" json:"kyc_level"`
	KYCCompletionDate      *time.Time `json:"kyc_completion_date"`
	AMLStatus              string     `gorm:"type:varchar(20)" json:"aml_status"`
	AMLCheckDate           *time.Time `json:"aml_check_date"`
	PEPStatus              bool       `gorm:"default:false" json:"pep_status"`
	SanctionsScreened      bool       `gorm:"default:false" json:"sanctions_screened"`
	SanctionsCheckDate     *time.Time `json:"sanctions_check_date"`
	EmailVerified          bool       `gorm:"default:false" json:"email_verified"`
	PhoneVerified          bool       `gorm:"default:false" json:"phone_verified"`
	DocumentsVerified      bool       `gorm:"default:false" json:"documents_verified"`
	VerificationMethod     string     `json:"verification_method"`
	VerificationDate       *time.Time `json:"verification_date"`
	GDPRConsent            bool       `gorm:"default:false" json:"gdpr_consent"`
	DataRetentionConsent   bool       `gorm:"default:false" json:"data_retention_consent"`
	MarketingConsent       bool       `gorm:"default:false" json:"marketing_consent"`
	TermsAcceptedVersion   string     `json:"terms_accepted_version"`
	TermsAcceptedAt        *time.Time `json:"terms_accepted_at"`
	RegulatoryRestrictions string     `gorm:"type:json" json:"regulatory_restrictions"`
}

// UserEngagement contains customer engagement and communication preferences
type UserEngagement struct {
	PreferredLanguage      string          `gorm:"default:'en'" json:"language_preference"`
	PreferredContactMethod string          `gorm:"type:varchar(20)" json:"preferred_contact_method"` // email/sms/push/whatsapp
	PreferredContactTime   string          `json:"preferred_contact_time"`
	OptInEmail             bool            `gorm:"default:true" json:"opt_in_email"`
	OptInSMS               bool            `gorm:"default:false" json:"opt_in_sms"`
	OptInPush              bool            `gorm:"default:false" json:"opt_in_push"`
	OptInWhatsApp          bool            `gorm:"default:false" json:"opt_in_whatsapp"`
	DoNotDisturb           bool            `gorm:"default:false" json:"do_not_disturb"`
	NotificationPrefs      string          `gorm:"type:json" json:"notification_preferences"`
	LastEngagementDate     *time.Time      `json:"last_engagement_date"`
	EngagementScore        float64         `gorm:"default:0" json:"engagement_score"`
	NetPromoterScore       *int            `json:"net_promoter_score"`
	CustomerSatisfaction   float64         `gorm:"default:0" json:"customer_satisfaction"`
	LastSurveyDate         *time.Time      `json:"last_survey_date"`
	RenewalReminderDays    int             `gorm:"default:30" json:"renewal_reminder_days"`
	CustomerLifetimeValue  decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"customer_lifetime_value"`
	RetentionProbability   float64         `gorm:"default:0.5" json:"retention_probability"`
	CrossSellScore         float64         `gorm:"default:0" json:"cross_sell_score"`
	UpgradeEligibility     bool            `gorm:"default:false" json:"upgrade_eligibility"`
}

// UserSupport contains customer support related information
type UserSupport struct {
	AssignedAgent           *uuid.UUID `gorm:"type:uuid" json:"assigned_agent"`
	ServiceTier             string     `gorm:"type:varchar(20)" json:"service_tier"` // standard/priority/vip
	OpenTicketsCount        int        `gorm:"default:0" json:"open_tickets_count"`
	TotalTicketsCount       int        `gorm:"default:0" json:"total_tickets_count"`
	AverageResponseTime     int        `json:"average_response_time_hours"`
	LastSupportInteraction  *time.Time `json:"last_support_interaction"`
	SatisfactionRating      float64    `gorm:"default:0" json:"satisfaction_rating"`
	PreferredSupportChannel string     `gorm:"type:varchar(20)" json:"preferred_support_channel"` // chat/phone/email
	SupportNotes            string     `gorm:"type:text" json:"support_notes"`
	EscalationCount         int        `gorm:"default:0" json:"escalation_count"`
	ComplaintCount          int        `gorm:"default:0" json:"complaint_count"`
	ComplimentCount         int        `gorm:"default:0" json:"compliment_count"`
}

// UserDigital contains digital experience and app usage data
type UserDigital struct {
	OnboardingCompleted  bool       `gorm:"default:false" json:"onboarding_completed"`
	OnboardingDate       *time.Time `json:"onboarding_date"`
	AppVersion           string     `json:"app_version"`
	LastAppActivity      *time.Time `json:"last_app_activity"`
	AppUsageFrequency    string     `json:"app_usage_frequency"` // daily/weekly/monthly/rare
	PushToken            string     `json:"push_token"`
	DeviceOS             string     `json:"device_os"` // iOS/Android/Web
	DeviceModel          string     `json:"device_model"`
	BetaTester           bool       `gorm:"default:false" json:"beta_tester"`
	FeatureFlags         string     `gorm:"type:json" json:"feature_flags"`
	AppCrashCount        int        `gorm:"default:0" json:"app_crash_count"`
	LastCrashDate        *time.Time `json:"last_crash_date"`
	DigitalAdoptionScore float64    `gorm:"default:0" json:"digital_adoption_score"`
	SelfServiceUsage     float64    `gorm:"default:0" json:"self_service_usage"`
}

// UserMarketing contains marketing and acquisition data
type UserMarketing struct {
	ReferralCode        string          `gorm:"uniqueIndex" json:"referral_code"`
	ReferredBy          *uuid.UUID      `gorm:"type:uuid" json:"referred_by"`
	ReferralCount       int             `gorm:"default:0" json:"referral_count"`
	ReferralRewards     decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"referral_rewards"`
	AcquisitionChannel  string          `json:"acquisition_channel"`
	AcquisitionCampaign string          `json:"acquisition_campaign"`
	AcquisitionCost     decimal.Decimal `gorm:"type:decimal(10,2)" json:"acquisition_cost"`
	AcquisitionDate     *time.Time      `json:"acquisition_date"`
	PromoCodesUsed      string          `gorm:"type:json" json:"promo_codes_used"`
	CampaignResponses   string          `gorm:"type:json" json:"campaign_responses"`
	MarketingTags       string          `gorm:"type:json" json:"marketing_tags"`
	ABTestGroups        string          `gorm:"type:json" json:"ab_test_groups"`
}

// UserHousehold contains family and group membership information
type UserHousehold struct {
	HouseholdID        *uuid.UUID      `gorm:"type:uuid" json:"household_id"`
	IsHouseholdHead    bool            `gorm:"default:false" json:"is_household_head"`
	FamilyMemberCount  int             `gorm:"default:1" json:"family_member_count"`
	FamilyRole         string          `json:"family_role"` // parent/spouse/child/other
	SharedDeductible   decimal.Decimal `gorm:"type:decimal(10,2)" json:"shared_deductible"`
	FamilyDiscountRate float64         `gorm:"default:0" json:"family_discount_rate"`
	GroupMemberships   string          `gorm:"type:json" json:"group_memberships"`
	CorporateAccountID *uuid.UUID      `gorm:"type:uuid" json:"corporate_account_id"`
	EmployeeID         string          `json:"employee_id"`
	DepartmentCode     string          `json:"department_code"`
	StudentVerified    bool            `gorm:"default:false" json:"student_verified"`
	StudentInstitution string          `json:"student_institution"`
	SeniorCitizen      bool            `gorm:"default:false" json:"senior_citizen"`
	MilitaryService    bool            `gorm:"default:false" json:"military_service"`
	VIPStatus          bool            `gorm:"default:false" json:"vip_status"`
	EarlyAdopter       bool            `gorm:"default:false" json:"early_adopter"`
}

// UserMetadata contains additional flexible data
type UserMetadata struct {
	SuspendedAt      *time.Time `json:"suspended_at"`
	SuspensionReason string     `json:"suspension_reason"`
	Notes            string     `gorm:"type:text" json:"notes"`
	Tags             string     `gorm:"type:json" json:"tags"`
	CustomFields     string     `gorm:"type:json" json:"custom_fields"`
	Metadata         string     `gorm:"type:json" json:"metadata"`
	Source           string     `json:"source"`
	ImportID         string     `json:"import_id"`
	LegacyID         string     `json:"legacy_id"`
}
