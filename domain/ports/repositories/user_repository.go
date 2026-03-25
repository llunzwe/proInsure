package repositories

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	
	"smartsure/internal/domain/models"
)

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	// === Basic CRUD Operations ===
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error

	// === Batch Operations ===
	CreateBatch(ctx context.Context, users []*models.User) error
	UpdateBatch(ctx context.Context, users []*models.User) error
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)

	// === Search & Filtering ===
	Search(ctx context.Context, filters UserSearchFilters) ([]*models.User, int64, error)
	SearchByName(ctx context.Context, name string, limit int) ([]*models.User, error)
	SearchByReferralCode(ctx context.Context, code string) (*models.User, error)
	GetByHouseholdID(ctx context.Context, householdID uuid.UUID) ([]*models.User, error)
	GetByCorporateAccountID(ctx context.Context, corporateID uuid.UUID) ([]*models.User, error)
	GetBySegment(ctx context.Context, segment string) ([]*models.User, error)
	GetByStatus(ctx context.Context, status string, limit int) ([]*models.User, error)

	// === Authentication & Security ===
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	IncrementLoginAttempts(ctx context.Context, userID uuid.UUID) (int, error)
	ResetLoginAttempts(ctx context.Context, userID uuid.UUID) error
	LockAccount(ctx context.Context, userID uuid.UUID, until time.Time) error
	UnlockAccount(ctx context.Context, userID uuid.UUID) error
	Enable2FA(ctx context.Context, userID uuid.UUID, secret string) error
	Disable2FA(ctx context.Context, userID uuid.UUID) error
	UpdateSecurityQuestions(ctx context.Context, userID uuid.UUID, questions map[string]string) error
	LogSecurityEvent(ctx context.Context, event *models.UserSecurityEvent) error
	GetSecurityEvents(ctx context.Context, userID uuid.UUID, limit int) ([]*models.UserSecurityEvent, error)

	// === Verification & Compliance ===
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
	VerifyPhone(ctx context.Context, userID uuid.UUID) error
	UpdateKYCStatus(ctx context.Context, userID uuid.UUID, status, level string) error
	UpdateAMLStatus(ctx context.Context, userID uuid.UUID, status string) error
	UpdatePEPStatus(ctx context.Context, userID uuid.UUID, isPEP bool) error
	UpdateSanctionsScreening(ctx context.Context, userID uuid.UUID, screened bool, date time.Time) error
	UpdateGDPRConsent(ctx context.Context, userID uuid.UUID, consent bool) error
	UpdateTermsAcceptance(ctx context.Context, userID uuid.UUID, version string) error
	GetComplianceStatus(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	GetUsersRequiringCompliance(ctx context.Context, checkType string) ([]*models.User, error)

	// === Risk & Fraud Management ===
	UpdateRiskScore(ctx context.Context, userID uuid.UUID, score float64) error
	UpdateFraudScore(ctx context.Context, userID uuid.UUID, score float64) error
	UpdateCreditScore(ctx context.Context, userID uuid.UUID, score int) error
	BlacklistUser(ctx context.Context, userID uuid.UUID, reason string) error
	UnblacklistUser(ctx context.Context, userID uuid.UUID) error
	GetBlacklistedUsers(ctx context.Context) ([]*models.User, error)
	IncrementSuspiciousActivity(ctx context.Context, userID uuid.UUID) (int, error)
	GetHighRiskUsers(ctx context.Context, threshold float64) ([]*models.User, error)
	GetFraudAlerts(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)

	// === Financial Operations ===
	UpdateOutstandingBalance(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error
	UpdateCreditLimit(ctx context.Context, userID uuid.UUID, limit decimal.Decimal) error
	RecordPaymentFailure(ctx context.Context, userID uuid.UUID) (int, error)
	RecordPaymentSuccess(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error
	UpdateTotalPremiumPaid(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error
	GetFinancialSummary(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	GetUsersWithOutstandingBalance(ctx context.Context) ([]*models.User, error)
	UpdatePreferredPaymentMethod(ctx context.Context, userID uuid.UUID, method string) error

	// === Insurance Operations ===
	UpdateInsuranceProfile(ctx context.Context, userID uuid.UUID, profile map[string]interface{}) error
	UpdatePolicyCount(ctx context.Context, userID uuid.UUID, active, total int) error
	UpdateClaimStats(ctx context.Context, userID uuid.UUID, total, approved int, avgAmount decimal.Decimal) error
	RecordClaim(ctx context.Context, userID uuid.UUID, claimID uuid.UUID) error
	UpdateDeviceCount(ctx context.Context, userID uuid.UUID, count int) error
	GetInsuranceHistory(ctx context.Context, userID uuid.UUID) ([]*models.UserInsuranceHistory, error)
	CreateInsuranceHistory(ctx context.Context, history *models.UserInsuranceHistory) error
	GetUsersWithExpiringPolicies(ctx context.Context, days int) ([]*models.User, error)
	GetUsersEligibleForUpgrade(ctx context.Context) ([]*models.User, error)

	// === Customer Engagement ===
	UpdateEngagementScore(ctx context.Context, userID uuid.UUID, score float64) error
	UpdateLastEngagement(ctx context.Context, userID uuid.UUID) error
	UpdateNPS(ctx context.Context, userID uuid.UUID, score int) error
	UpdateCustomerSatisfaction(ctx context.Context, userID uuid.UUID, rating float64) error
	UpdateCommunicationPreferences(ctx context.Context, userID uuid.UUID, prefs map[string]bool) error
	GetEngagedUsers(ctx context.Context, minScore float64) ([]*models.User, error)
	GetInactiveUsers(ctx context.Context, days int) ([]*models.User, error)
	TrackActivity(ctx context.Context, activity *models.UserActivity) error
	GetUserActivities(ctx context.Context, userID uuid.UUID, limit int) ([]*models.UserActivity, error)

	// === Loyalty & Rewards ===
	UpdateLoyaltyTier(ctx context.Context, userID uuid.UUID, tier string) error
	UpdateReferralCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetReferrals(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
	CreateReward(ctx context.Context, reward *models.UserReward) error
	GetUserRewards(ctx context.Context, userID uuid.UUID) ([]*models.UserReward, error)
	RedeemReward(ctx context.Context, rewardID uuid.UUID) error
	GetTopReferrers(ctx context.Context, limit int) ([]*models.User, error)
	CalculateLifetimeValue(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error)

	// === Support & Communication ===
	CreateSupportTicket(ctx context.Context, ticket *models.UserSupportTicket) error
	GetSupportTickets(ctx context.Context, userID uuid.UUID) ([]*models.UserSupportTicket, error)
	UpdateSupportTier(ctx context.Context, userID uuid.UUID, tier string) error
	LogNotification(ctx context.Context, notification *models.UserNotificationLog) error
	GetNotificationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.UserNotificationLog, error)
	UpdateNotificationStatus(ctx context.Context, notificationID uuid.UUID, status string) error
	GetUndeliveredNotifications(ctx context.Context, userID uuid.UUID) ([]*models.UserNotificationLog, error)

	// === Family & Household ===
	CreateHousehold(ctx context.Context, headUserID uuid.UUID) (uuid.UUID, error)
	AddToHousehold(ctx context.Context, userID, householdID uuid.UUID) error
	RemoveFromHousehold(ctx context.Context, userID uuid.UUID) error
	SetHouseholdHead(ctx context.Context, userID, householdID uuid.UUID) error
	GetHouseholdMembers(ctx context.Context, householdID uuid.UUID) ([]*models.User, error)
	UpdateFamilyDiscount(ctx context.Context, householdID uuid.UUID, rate float64) error
	GetFamilyStats(ctx context.Context, householdID uuid.UUID) (map[string]interface{}, error)

	// === Document Management ===
	CreateDocument(ctx context.Context, document *models.UserDocument) error
	GetUserDocuments(ctx context.Context, userID uuid.UUID) ([]*models.UserDocument, error)
	UpdateDocumentStatus(ctx context.Context, documentID uuid.UUID, status string) error
	DeleteDocument(ctx context.Context, documentID uuid.UUID) error
	GetExpiredDocuments(ctx context.Context) ([]*models.UserDocument, error)
	VerifyDocument(ctx context.Context, documentID uuid.UUID, verifierID uuid.UUID) error

	// === Preferences & Settings ===
	GetUserPreferences(ctx context.Context, userID uuid.UUID) ([]*models.UserPreference, error)
	SetPreference(ctx context.Context, pref *models.UserPreference) error
	DeletePreference(ctx context.Context, userID uuid.UUID, key string) error
	GetPreferenceByKey(ctx context.Context, userID uuid.UUID, key string) (*models.UserPreference, error)
	UpdateLanguagePreference(ctx context.Context, userID uuid.UUID, language string) error
	UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string) error

	// === Payment Methods ===
	AddPaymentMethod(ctx context.Context, method *models.UserPaymentMethod) error
	GetPaymentMethods(ctx context.Context, userID uuid.UUID) ([]*models.UserPaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, method *models.UserPaymentMethod) error
	DeletePaymentMethod(ctx context.Context, methodID uuid.UUID) error
	SetDefaultPaymentMethod(ctx context.Context, userID, methodID uuid.UUID) error
	GetDefaultPaymentMethod(ctx context.Context, userID uuid.UUID) (*models.UserPaymentMethod, error)

	// === Audit & Compliance ===
	LogAuditEvent(ctx context.Context, audit *models.UserAuditLog) error
	GetAuditHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.UserAuditLog, error)
	GetAuditByDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.UserAuditLog, error)
	GetSuspiciousActivities(ctx context.Context, threshold int) ([]*models.User, error)
	GenerateComplianceReport(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)

	// === Analytics & Reporting ===
	GetUserStatistics(ctx context.Context, userID uuid.UUID) (UserStatistics, error)
	GetSegmentAnalytics(ctx context.Context, segment string) (map[string]interface{}, error)
	GetChurnPrediction(ctx context.Context, userID uuid.UUID) (float64, error)
	GetRetentionMetrics(ctx context.Context) (map[string]interface{}, error)
	GetUsersByRiskCategory(ctx context.Context) (map[string][]*models.User, error)
	GetRevenueBySegment(ctx context.Context) (map[string]decimal.Decimal, error)
	GetConversionMetrics(ctx context.Context, period string) (map[string]interface{}, error)

	// === Suspension & Account Management ===
	SuspendUser(ctx context.Context, userID uuid.UUID, reason string) error
	ReactivateUser(ctx context.Context, userID uuid.UUID) error
	GetSuspendedUsers(ctx context.Context) ([]*models.User, error)
	ScheduleDeletion(ctx context.Context, userID uuid.UUID, date time.Time) error
	CancelDeletion(ctx context.Context, userID uuid.UUID) error
	AnonymizeUser(ctx context.Context, userID uuid.UUID) error
	ExportUserData(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)

	// === Performance & Optimization ===
	GetUserWithRelations(ctx context.Context, userID uuid.UUID, relations []string) (*models.User, error)
	PreloadUserData(ctx context.Context, userID uuid.UUID) error
	RefreshUserCache(ctx context.Context, userID uuid.UUID) error
	BulkUpdateRiskScores(ctx context.Context, scores map[uuid.UUID]float64) error
	OptimizeUserQueries(ctx context.Context) error
}

// UserSearchFilters defines search criteria for users
type UserSearchFilters struct {
	// Basic filters
	Name             string
	Email            string
	Phone            string
	Status           string
	Segment          string
	SubscriptionTier string

	// Location filters
	Country       string
	StateProvince string
	City          string
	CoverageZone  string

	// Risk filters
	MinRiskScore    float64
	MaxRiskScore    float64
	MinFraudScore   float64
	MaxFraudScore   float64
	BlacklistStatus *bool

	// Financial filters
	MinCreditScore        int
	MaxCreditScore        int
	HasOutstandingBalance *bool
	MinLifetimeValue      decimal.Decimal
	MaxLifetimeValue      decimal.Decimal

	// Insurance filters
	HasActivePolicies *bool
	MinPoliciesCount  int
	MaxPoliciesCount  int
	MinClaimsCount    int
	MaxClaimsCount    int

	// Engagement filters
	MinEngagementScore float64
	MaxEngagementScore float64
	LastEngagementDays int
	IsActive           *bool

	// Compliance filters
	KYCStatus     string
	AMLStatus     string
	PEPStatus     *bool
	GDPRConsent   *bool
	EmailVerified *bool
	PhoneVerified *bool

	// Date filters
	CreatedAfter  time.Time
	CreatedBefore time.Time
	UpdatedAfter  time.Time
	UpdatedBefore time.Time

	// Family/Group filters
	HouseholdID        *uuid.UUID
	IsHouseholdHead    *bool
	CorporateAccountID *uuid.UUID

	// Pagination
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
}

// UserStatistics represents aggregated user statistics
type UserStatistics struct {
	UserID                uuid.UUID
	TotalPolicies         int
	ActivePolicies        int
	TotalClaims           int
	ApprovedClaims        int
	ClaimApprovalRate     float64
	AverageClaimAmount    decimal.Decimal
	TotalPremiumPaid      decimal.Decimal
	OutstandingBalance    decimal.Decimal
	PaymentFailureRate    float64
	DeviceCount           int
	ReferralCount         int
	EngagementScore       float64
	LifetimeValue         decimal.Decimal
	ChurnProbability      float64
	LastActivityDate      time.Time
	DaysSinceLastActivity int
	RiskCategory          string
	CustomerSegment       string
	LoyaltyTier           string
}
