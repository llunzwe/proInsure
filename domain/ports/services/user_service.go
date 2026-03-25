package services

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	
	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
)

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	UserType    string `json:"user_type"`
	UserRole    string `json:"user_role"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
}

// UserService defines the interface for user business operations
type UserService interface {
	// === Registration & Onboarding ===
	RegisterUser(ctx context.Context, user *models.User) error
	CompleteOnboarding(ctx context.Context, userID uuid.UUID) error
	VerifyIdentity(ctx context.Context, userID uuid.UUID, documents []models.UserDocument) error
	InitiateKYC(ctx context.Context, userID uuid.UUID) error
	CompleteKYC(ctx context.Context, userID uuid.UUID, kycData map[string]interface{}) error
	SetupAccount(ctx context.Context, userID uuid.UUID, preferences map[string]interface{}) error
	CreateWelcomePackage(ctx context.Context, userID uuid.UUID) error

	// === Authentication & Authorization ===
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, string, error)
	RefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) (string, error)
	Logout(ctx context.Context, userID uuid.UUID, token string) error
	ResetPassword(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	Enable2FA(ctx context.Context, userID uuid.UUID) (string, error)
	Verify2FA(ctx context.Context, userID uuid.UUID, code string) error
	Disable2FA(ctx context.Context, userID uuid.UUID, password string) error
	GenerateAPIKey(ctx context.Context, userID uuid.UUID, name string) (string, error)
	RevokeAPIKey(ctx context.Context, userID uuid.UUID, keyID string) error

	// === Profile Management ===
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
	UploadProfilePicture(ctx context.Context, userID uuid.UUID, imageData []byte) (string, error)
	UpdateContactInfo(ctx context.Context, userID uuid.UUID, contact map[string]string) error
	UpdateAddress(ctx context.Context, userID uuid.UUID, address map[string]string) error
	UpdatePreferences(ctx context.Context, userID uuid.UUID, preferences map[string]interface{}) error
	GetAccountSettings(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	UpdateAccountSettings(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error

	// === Risk Assessment & Underwriting ===
	AssessUserRisk(ctx context.Context, userID uuid.UUID) (float64, error)
	CalculateRiskScore(ctx context.Context, userID uuid.UUID) (float64, map[string]interface{}, error)
	UpdateRiskProfile(ctx context.Context, userID uuid.UUID, factors map[string]interface{}) error
	PerformCreditCheck(ctx context.Context, userID uuid.UUID) (int, error)
	CalculateInsuranceEligibility(ctx context.Context, userID uuid.UUID) (bool, []string, error)
	DetermineUnderwritingDecision(ctx context.Context, userID uuid.UUID) (string, map[string]interface{}, error)
	GetRiskFactors(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	PredictChurnProbability(ctx context.Context, userID uuid.UUID) (float64, error)
	CalculatePremiumRate(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error)

	// === Fraud Detection & Prevention ===
	CheckForFraud(ctx context.Context, userID uuid.UUID, activity string) (bool, float64, error)
	DetectAnomalousActivity(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	FlagSuspiciousUser(ctx context.Context, userID uuid.UUID, reason string) error
	InvestigateFraud(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	BlacklistUser(ctx context.Context, userID uuid.UUID, reason string) error
	WhitelistUser(ctx context.Context, userID uuid.UUID) error
	GetFraudHistory(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	RunBackgroundCheck(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	VerifyIdentityDocuments(ctx context.Context, userID uuid.UUID, documents []models.UserDocument) (bool, error)
	CheckSanctionsList(ctx context.Context, userID uuid.UUID) (bool, error)

	// === Customer Segmentation & Analytics ===
	SegmentUser(ctx context.Context, userID uuid.UUID) (string, error)
	CalculateCustomerLifetimeValue(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error)
	AnalyzeUserBehavior(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	GetUserInsights(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	PredictNextBestAction(ctx context.Context, userID uuid.UUID) (string, float64, error)
	CalculateRetentionScore(ctx context.Context, userID uuid.UUID) (float64, error)
	IdentifyUpsellOpportunities(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	GetUserJourney(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	TrackUserEngagement(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)

	// === Insurance Operations ===
	GetInsuranceProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	GetUserPolicies(ctx context.Context, userID uuid.UUID) ([]*models.Policy, error)
	GetUserClaims(ctx context.Context, userID uuid.UUID) ([]*models.Claim, error)
	GetUserDevices(ctx context.Context, userID uuid.UUID) ([]*models.Device, error)
	CalculateTotalCoverage(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error)
	GetClaimHistory(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	CheckPolicyEligibility(ctx context.Context, userID uuid.UUID, policyType string) (bool, error)
	GetRecommendedPolicies(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	SchedulePolicyReview(ctx context.Context, userID uuid.UUID) error
	ProcessRenewal(ctx context.Context, userID uuid.UUID, policyID uuid.UUID) error

	// === Payment & Billing ===
	ProcessPayment(ctx context.Context, userID uuid.UUID, payment map[string]interface{}) error
	SetupAutoPay(ctx context.Context, userID uuid.UUID, methodID uuid.UUID) error
	UpdatePaymentMethod(ctx context.Context, userID uuid.UUID, method *models.UserPaymentMethod) error
	GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	GetOutstandingBalance(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error)
	ProcessRefund(ctx context.Context, userID uuid.UUID, amount decimal.Decimal, reason string) error
	ApplyDiscount(ctx context.Context, userID uuid.UUID, discountCode string) (decimal.Decimal, error)
	CalculateLoyaltyDiscount(ctx context.Context, userID uuid.UUID) (float64, error)
	SetupPaymentPlan(ctx context.Context, userID uuid.UUID, plan map[string]interface{}) error
	ProcessPremiumFinancing(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error

	// === Family & Household Management ===
	CreateFamilyPlan(ctx context.Context, headUserID uuid.UUID) (uuid.UUID, error)
	AddFamilyMember(ctx context.Context, householdID, memberID uuid.UUID) error
	RemoveFamilyMember(ctx context.Context, householdID, memberID uuid.UUID) error
	GetFamilyMembers(ctx context.Context, householdID uuid.UUID) ([]*models.User, error)
	CalculateFamilyDiscount(ctx context.Context, householdID uuid.UUID) (float64, error)
	ManageParentalControls(ctx context.Context, householdID uuid.UUID, controls map[string]interface{}) error
	ShareBenefits(ctx context.Context, householdID uuid.UUID, benefits map[string]interface{}) error
	GetFamilyUsageStats(ctx context.Context, householdID uuid.UUID) (map[string]interface{}, error)
	SetFamilyBudget(ctx context.Context, householdID uuid.UUID, budget decimal.Decimal) error

	// === Loyalty & Rewards ===
	CalculateLoyaltyPoints(ctx context.Context, userID uuid.UUID) (int, error)
	RedeemLoyaltyPoints(ctx context.Context, userID uuid.UUID, points int, rewardType string) error
	GetAvailableRewards(ctx context.Context, userID uuid.UUID) ([]models.UserReward, error)
	ApplyReferralBonus(ctx context.Context, referrerID, referredID uuid.UUID) error
	GetReferralStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	UpdateLoyaltyTier(ctx context.Context, userID uuid.UUID) (string, error)
	GetLoyaltyBenefits(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)
	CreateSpecialOffer(ctx context.Context, userID uuid.UUID, offer map[string]interface{}) error
	TrackAchievements(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)

	// === Communication & Engagement ===
	SendNotification(ctx context.Context, userID uuid.UUID, notification map[string]interface{}) error
	SendEmail(ctx context.Context, userID uuid.UUID, template string, data map[string]interface{}) error
	SendSMS(ctx context.Context, userID uuid.UUID, message string) error
	SendPushNotification(ctx context.Context, userID uuid.UUID, title, body string) error
	GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (map[string]bool, error)
	UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, prefs map[string]bool) error
	ScheduleCommunication(ctx context.Context, userID uuid.UUID, comm map[string]interface{}) error
	TrackCommunicationEffectiveness(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	ManageCampaignSubscriptions(ctx context.Context, userID uuid.UUID, campaigns []string) error

	// === Support & Service ===
	CreateSupportTicket(ctx context.Context, userID uuid.UUID, ticket *models.UserSupportTicket) error
	GetSupportHistory(ctx context.Context, userID uuid.UUID) ([]*models.UserSupportTicket, error)
	EscalateIssue(ctx context.Context, ticketID uuid.UUID) error
	ProvideSelfService(ctx context.Context, userID uuid.UUID, issue string) (map[string]interface{}, error)
	ScheduleCallback(ctx context.Context, userID uuid.UUID, preferredTime time.Time) error
	GetServiceLevel(ctx context.Context, userID uuid.UUID) (string, error)
	UpgradeServiceTier(ctx context.Context, userID uuid.UUID, tier string) error
	TrackSatisfaction(ctx context.Context, userID uuid.UUID, rating int, feedback string) error
	GetSupportMetrics(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)

	// === Compliance & Privacy ===
	GetGDPRData(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	ExerciseDataRights(ctx context.Context, userID uuid.UUID, right string) error
	RequestDataDeletion(ctx context.Context, userID uuid.UUID) error
	ExportUserData(ctx context.Context, userID uuid.UUID) ([]byte, error)
	UpdateConsents(ctx context.Context, userID uuid.UUID, consents map[string]bool) error
	GetComplianceStatus(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	PerformAMLCheck(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdatePEPStatus(ctx context.Context, userID uuid.UUID, isPEP bool) error
	GenerateComplianceReport(ctx context.Context, userID uuid.UUID) ([]byte, error)
	AuditUserAccess(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)

	// === Account Management ===
	SuspendAccount(ctx context.Context, userID uuid.UUID, reason string) error
	ReactivateAccount(ctx context.Context, userID uuid.UUID) error
	CloseAccount(ctx context.Context, userID uuid.UUID, reason string) error
	MergeAccounts(ctx context.Context, primaryID, secondaryID uuid.UUID) error
	TransferAccount(ctx context.Context, userID uuid.UUID, newOwner map[string]interface{}) error
	BackupAccountData(ctx context.Context, userID uuid.UUID) error
	RestoreAccountData(ctx context.Context, userID uuid.UUID, backupID string) error
	ScheduleAccountDeletion(ctx context.Context, userID uuid.UUID, date time.Time) error
	CancelDeletionRequest(ctx context.Context, userID uuid.UUID) error

	// === Integration & Migration ===
	ImportUserData(ctx context.Context, source string, data map[string]interface{}) (*models.User, error)
	MigrateFromCompetitor(ctx context.Context, competitorData map[string]interface{}) (*models.User, error)
	SyncWithExternalSystem(ctx context.Context, userID uuid.UUID, system string) error
	ExportToPartner(ctx context.Context, userID uuid.UUID, partnerID string) error
	ValidateDataIntegrity(ctx context.Context, userID uuid.UUID) (bool, []string, error)
	ReconcileAccounts(ctx context.Context, userID uuid.UUID) error

	// === Reporting & Analytics ===
	GenerateUserReport(ctx context.Context, userID uuid.UUID, reportType string) ([]byte, error)
	GetUserMetrics(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	TrackUserActivity(ctx context.Context, activity *models.UserActivity) error
	AnalyzeUserTrends(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error)
	GetPerformanceIndicators(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	BenchmarkUser(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	PredictUserBehavior(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	GenerateInsights(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error)

	// === Bulk Operations ===
	BulkImportUsers(ctx context.Context, users []*models.User) (int, []error)
	BulkUpdateUsers(ctx context.Context, updates map[uuid.UUID]map[string]interface{}) (int, []error)
	BulkNotifyUsers(ctx context.Context, userIDs []uuid.UUID, message map[string]interface{}) (int, []error)
	BulkAssessRisk(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]float64, error)
	BulkSegmentUsers(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error)
	BulkGenerateReports(ctx context.Context, userIDs []uuid.UUID, reportType string) (map[uuid.UUID][]byte, error)

	// === Repository Access ===
	GetRepository() repositories.UserRepository

	// ======================================
	// MISSING IMPLEMENTED METHODS
	// ======================================

	// Password Management
	InitiatePasswordReset(ctx context.Context, email string) error
	CompletePasswordReset(ctx context.Context, token, newPassword string) error // Renamed from ResetPassword to avoid duplicate

	// Profile Management
	UpdateUserAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error
	DeleteUserAvatar(ctx context.Context, userID uuid.UUID) error
	GetUserVerificationStatus(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)

	// Account Management
	DeactivateUser(ctx context.Context, userID uuid.UUID, reason, feedback string) error
	ReactivateUser(ctx context.Context, userID uuid.UUID) error

	// API Key Management
	GenerateUserAPIKey(ctx context.Context, userID uuid.UUID, name, description string, expiryDays int) (interface{}, error) // UserAPIKey type not found
	RevokeUserAPIKey(ctx context.Context, userID uuid.UUID, apiKeyID uuid.UUID) error

	// Legacy methods for backward compatibility
	CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*models.User, error)
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
	VerifyPhone(ctx context.Context, userID uuid.UUID) error
	UpdateKYCStatus(ctx context.Context, userID uuid.UUID, status, level string) error
	ValidateToken(tokenString string) (uuid.UUID, error)
}
