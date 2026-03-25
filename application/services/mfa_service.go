package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/pkg/cache"
	"smartsure/pkg/notification"
)

// MFAService handles MFA operations
type MFAService struct {
	db              *gorm.DB
	mfaRepo         repositories.MFARepository
	sessionRepo     repositories.SessionRepository
	notificationSvc *notification.Service
	cache           *cache.RedisCache
}

// NewMFAService creates a new MFA service
func NewMFAService(db *gorm.DB, notificationSvc *notification.Service, cache *cache.RedisCache) *MFAService {
	return &MFAService{
		db:              db,
		mfaRepo:         repositories.NewMFARepository(db),
		sessionRepo:     repositories.NewSessionRepository(db),
		notificationSvc: notificationSvc,
		cache:           cache,
	}
}

// MFASetupRequest represents MFA setup request
type MFASetupRequest struct {
	UserID   uuid.UUID          `json:"user_id" binding:"required"`
	Provider models.MFAProvider `json:"provider" binding:"required"`
	Data     map[string]string  `json:"data"`
}

// MFASetupResponse represents MFA setup response
type MFASetupResponse struct {
	MethodID    uuid.UUID `json:"method_id"`
	Provider    string    `json:"provider"`
	Secret      string    `json:"secret,omitempty"`
	QRCodeURL   string    `json:"qr_code_url,omitempty"`
	BackupCodes []string  `json:"backup_codes,omitempty"`
}

// MFAVerifyRequest represents MFA verification request
type MFAVerifyRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	MethodID uuid.UUID `json:"method_id" binding:"required"`
	Code     string    `json:"code" binding:"required"`
}

// SessionInfo represents session information
type SessionInfo struct {
	SessionID         string    `json:"session_id"`
	UserID            uuid.UUID `json:"user_id"`
	IPAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	LastActivity      time.Time `json:"last_activity"`
	ExpiresAt         time.Time `json:"expires_at"`
	IsActive          bool      `json:"is_active"`
	DeviceFingerprint string    `json:"device_fingerprint"`
}

// SetupMFA sets up MFA for a user
func (s *MFAService) SetupMFA(ctx context.Context, req *MFASetupRequest) (*MFASetupResponse, error) {
	// Check if user already has MFA enabled
	existingMethods, err := s.mfaRepo.GetUserMFAMethods(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing MFA methods: %w", err)
	}

	// Limit to 3 MFA methods per user
	if len(existingMethods) >= 3 {
		return nil, fmt.Errorf("maximum MFA methods reached")
	}

	method := &models.MFAMethod{
		UserID:   req.UserID,
		Provider: req.Provider,
		Status:   models.MFAStatusPending,
	}

	switch req.Provider {
	case models.MFAProviderTOTP:
		return s.setupTOTP(ctx, method)
	case models.MFAProviderSMS:
		return s.setupSMS(ctx, method, req.Data)
	case models.MFAProviderEmail:
		return s.setupEmail(ctx, method, req.Data)
	case models.MFAProviderPush:
		return s.setupPush(ctx, method, req.Data)
	case models.MFAProviderHardware:
		return s.setupHardware(ctx, method, req.Data)
	default:
		return nil, fmt.Errorf("unsupported MFA provider: %s", req.Provider)
	}
}

// setupTOTP sets up TOTP MFA
func (s *MFAService) setupTOTP(ctx context.Context, method *models.MFAMethod) (*MFASetupResponse, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SmartSure",
		AccountName: "user@" + method.UserID.String(),
		SecretSize:  32,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	method.Secret = key.Secret()

	if err := s.mfaRepo.CreateMFAMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to create MFA method: %w", err)
	}

	// Generate backup codes
	backupCodes := s.generateBackupCodes()
	method.BackupCodes = backupCodes

	if err := s.mfaRepo.UpdateMFAMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to update MFA method with backup codes: %w", err)
	}

	return &MFASetupResponse{
		MethodID:    method.ID,
		Provider:    string(method.Provider),
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

// setupSMS sets up SMS MFA
func (s *MFAService) setupSMS(ctx context.Context, method *models.MFAMethod, data map[string]string) (*MFASetupResponse, error) {
	phoneNumber, ok := data["phone_number"]
	if !ok {
		return nil, fmt.Errorf("phone_number required for SMS MFA")
	}

	method.PhoneNumber = phoneNumber

	if err := s.mfaRepo.CreateMFAMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to create MFA method: %w", err)
	}

	// Send verification code
	if err := s.sendSMSVerification(method); err != nil {
		return nil, fmt.Errorf("failed to send SMS verification: %w", err)
	}

	return &MFASetupResponse{
		MethodID: method.ID,
		Provider: string(method.Provider),
	}, nil
}

// setupEmail sets up Email MFA
func (s *MFAService) setupEmail(ctx context.Context, method *models.MFAMethod, data map[string]string) (*MFASetupResponse, error) {
	email, ok := data["email"]
	if !ok {
		return nil, fmt.Errorf("email required for email MFA")
	}

	method.Email = email

	if err := s.mfaRepo.CreateMFAMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to create MFA method: %w", err)
	}

	// Send verification code
	if err := s.sendEmailVerification(method); err != nil {
		return nil, fmt.Errorf("failed to send email verification: %w", err)
	}

	return &MFASetupResponse{
		MethodID: method.ID,
		Provider: string(method.Provider),
	}, nil
}

// setupPush sets up Push MFA
func (s *MFAService) setupPush(ctx context.Context, method *models.MFAMethod, data map[string]string) (*MFASetupResponse, error) {
	deviceToken, ok := data["device_token"]
	if !ok {
		return nil, fmt.Errorf("device_token required for push MFA")
	}

	method.DeviceToken = deviceToken

	if err := s.mfaRepo.CreateMFAMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to create MFA method: %w", err)
	}

	return &MFASetupResponse{
		MethodID: method.ID,
		Provider: string(method.Provider),
	}, nil
}

// setupHardware sets up Hardware MFA
func (s *MFAService) setupHardware(ctx context.Context, method *models.MFAMethod, data map[string]string) (*MFASetupResponse, error) {
	hardwareKey, ok := data["hardware_key"]
	if !ok {
		return nil, fmt.Errorf("hardware_key required for hardware MFA")
	}

	method.HardwareKey = hardwareKey

	if err := s.mfaRepo.CreateMFAMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("failed to create MFA method: %w", err)
	}

	return &MFASetupResponse{
		MethodID: method.ID,
		Provider: string(method.Provider),
	}, nil
}

// VerifyMFA verifies an MFA code
func (s *MFAService) VerifyMFA(ctx context.Context, req *MFAVerifyRequest) error {
	method, err := s.mfaRepo.GetMFAMethod(ctx, req.MethodID)
	if err != nil {
		return fmt.Errorf("failed to get MFA method: %w", err)
	}

	if method.UserID != req.UserID {
		return fmt.Errorf("MFA method does not belong to user")
	}

	if method.Status != models.MFAStatusActive {
		return fmt.Errorf("MFA method is not active")
	}

	// Create attempt record
	attempt := &models.MFAAttempt{
		UserID:    req.UserID,
		MethodID:  req.MethodID,
		Provider:  method.Provider,
		CodeHash:  hashCode(req.Code),
		Status:    "pending",
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5 minutes
	}

	// Get client info from context (would be set by middleware)
	if httpReq, ok := ctx.Value("http_request").(*http.Request); ok {
		attempt.IPAddress = getClientIP(httpReq)
		attempt.UserAgent = httpReq.UserAgent()
	}

	valid := false
	switch method.Provider {
	case models.MFAProviderTOTP:
		valid = totp.Validate(req.Code, method.Secret)
	case models.MFAProviderSMS, models.MFAProviderEmail:
		// Check cached verification code
		cacheKey := fmt.Sprintf("mfa_code:%s:%s", req.MethodID, req.Code)
		exists, err := s.cache.Exists(ctx, cacheKey)
		if err != nil {
			return fmt.Errorf("failed to check verification code: %w", err)
		}
		valid = exists
		if valid {
			s.cache.Delete(ctx, cacheKey) // One-time use
		}
	case models.MFAProviderPush:
		// Push notification verification would be handled asynchronously
		valid = true // Placeholder
	case models.MFAProviderHardware:
		// Hardware token verification
		valid = s.verifyHardwareToken(req.Code, method.HardwareKey)
	}

	if valid {
		attempt.Status = "success"
		now := time.Now()
		method.LastUsed = &now
		method.FailureCount = 0
		method.VerifiedAt = &now

		if err := s.mfaRepo.UpdateMFAMethod(ctx, method); err != nil {
			return fmt.Errorf("failed to update MFA method: %w", err)
		}
	} else {
		attempt.Status = "failed"
		method.FailureCount++

		// Block method after 5 failures
		if method.FailureCount >= 5 {
			method.Status = models.MFAStatusBlocked
		}

		if err := s.mfaRepo.UpdateMFAMethod(ctx, method); err != nil {
			return fmt.Errorf("failed to update MFA method: %w", err)
		}

		return fmt.Errorf("invalid MFA code")
	}

	if err := s.mfaRepo.CreateMFAAttempt(ctx, attempt); err != nil {
		return fmt.Errorf("failed to create MFA attempt: %w", err)
	}

	return nil
}

// DisableMFA disables an MFA method
func (s *MFAService) DisableMFA(ctx context.Context, userID, methodID uuid.UUID) error {
	method, err := s.mfaRepo.GetMFAMethod(ctx, methodID)
	if err != nil {
		return fmt.Errorf("failed to get MFA method: %w", err)
	}

	if method.UserID != userID {
		return fmt.Errorf("MFA method does not belong to user")
	}

	method.Status = models.MFAStatusInactive
	return s.mfaRepo.UpdateMFAMethod(ctx, method)
}

// GetUserMFAMethods gets all MFA methods for a user
func (s *MFAService) GetUserMFAMethods(ctx context.Context, userID uuid.UUID) ([]*models.MFAMethod, error) {
	return s.mfaRepo.GetUserMFAMethods(ctx, userID)
}

// CreateSession creates a new user session
func (s *MFAService) CreateSession(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string, expiresAt time.Time, ipAddress, userAgent string) (*models.UserSession, error) {
	sessionID := generateSessionID()

	session := &models.UserSession{
		UserID:       userID,
		SessionID:    sessionID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		LastActivity: time.Now(),
		IsActive:     true,
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetActiveSessions gets all active sessions for a user
func (s *MFAService) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]*SessionInfo, error) {
	sessions, err := s.sessionRepo.GetUserActiveSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	var sessionInfos []*SessionInfo
	for _, session := range sessions {
		sessionInfos = append(sessionInfos, &SessionInfo{
			SessionID:         session.SessionID,
			UserID:            session.UserID,
			IPAddress:         session.IPAddress,
			UserAgent:         session.UserAgent,
			LastActivity:      session.LastActivity,
			ExpiresAt:         session.ExpiresAt,
			IsActive:          session.IsActive,
			DeviceFingerprint: session.DeviceFingerprint,
		})
	}

	return sessionInfos, nil
}

// RevokeSession revokes a specific session
func (s *MFAService) RevokeSession(ctx context.Context, userID uuid.UUID, sessionID string) error {
	session, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session.UserID != userID {
		return fmt.Errorf("session does not belong to user")
	}

	now := time.Now()
	session.RevokedAt = &now
	session.IsActive = false

	return s.sessionRepo.UpdateSession(ctx, session)
}

// RevokeAllSessions revokes all sessions for a user except current
func (s *MFAService) RevokeAllSessions(ctx context.Context, userID uuid.UUID, currentSessionID string) error {
	return s.sessionRepo.RevokeAllUserSessions(ctx, userID, currentSessionID)
}

// UpdateSessionActivity updates session last activity
func (s *MFAService) UpdateSessionActivity(ctx context.Context, sessionID string) error {
	return s.sessionRepo.UpdateSessionActivity(ctx, sessionID)
}

// ValidateSession validates if a session is active and not expired
func (s *MFAService) ValidateSession(ctx context.Context, sessionID string) (*models.UserSession, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if !session.IsActive || session.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("session is invalid or expired")
	}

	return session, nil
}

// GetMFAMethod gets a single MFA method
func (s *MFAService) GetMFAMethod(ctx context.Context, methodID uuid.UUID) (*models.MFAMethod, error) {
	return s.mfaRepo.GetMFAMethod(ctx, methodID)
}

// UpdateMFAMethod updates an MFA method
func (s *MFAService) UpdateMFAMethod(ctx context.Context, method *models.MFAMethod) error {
	return s.mfaRepo.UpdateMFAMethod(ctx, method)
}

// UnsetDefaultMFAMethods unsets all default MFA methods for a user
func (s *MFAService) UnsetDefaultMFAMethods(ctx context.Context, userID uuid.UUID) error {
	return s.db.WithContext(ctx).Model(&models.MFAMethod{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

// GetMFAAttempts gets MFA attempts for a user
func (s *MFAService) GetMFAAttempts(ctx context.Context, userID uuid.UUID, limit int) ([]*models.MFAAttempt, error) {
	return s.mfaRepo.GetMFAAttempts(ctx, userID, limit)
}

// UpdateSession updates a session
func (s *MFAService) UpdateSession(ctx context.Context, session *models.UserSession) error {
	return s.sessionRepo.UpdateSession(ctx, session)
}

// generateBackupCodes generates 10 backup codes
func (s *MFAService) generateBackupCodes() []string {
	var codes []string
	for i := 0; i < 10; i++ {
		code, _ := generateRandomCode(8)
		codes = append(codes, code)
	}
	return codes
}

// sendSMSVerification sends SMS verification code
func (s *MFAService) sendSMSVerification(method *models.MFAMethod) error {
	code, _ := generateRandomCode(6)

	// Cache the code for 5 minutes
	cacheKey := fmt.Sprintf("mfa_code:%s:%s", method.ID, code)
	ctx := context.Background()
	s.cache.Set(ctx, cacheKey, true, 5*time.Minute)

	// Send SMS (placeholder - integrate with SMS provider)
	message := fmt.Sprintf("Your SmartSure verification code is: %s", code)
	return s.notificationSvc.SendSMS(method.PhoneNumber, message)
}

// sendEmailVerification sends email verification code
func (s *MFAService) sendEmailVerification(method *models.MFAMethod) error {
	code, _ := generateRandomCode(6)

	// Cache the code for 5 minutes
	cacheKey := fmt.Sprintf("mfa_code:%s:%s", method.ID, code)
	ctx := context.Background()
	s.cache.Set(ctx, cacheKey, true, 5*time.Minute)

	// Send email (placeholder - integrate with email provider)
	subject := "SmartSure MFA Verification"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 5 minutes.", code)
	return s.notificationSvc.SendEmail(method.Email, subject, body)
}

// verifyHardwareToken verifies hardware token (placeholder implementation)
func (s *MFAService) verifyHardwareToken(code, hardwareKey string) bool {
	// Placeholder - integrate with hardware token provider
	return len(code) == 6 && len(hardwareKey) > 0
}

// generateRandomCode generates a random numeric code
func generateRandomCode(length int) (string, error) {
	var code strings.Builder
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code.WriteString(strconv.Itoa(int(num.Int64())))
	}
	return code.String(), nil
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// hashCode creates a SHA256 hash of the code
func hashCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

// getClientIP gets the client IP address from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if multiple
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip, _, _ = strings.Cut(ip, ":")
	}
	return ip
}
