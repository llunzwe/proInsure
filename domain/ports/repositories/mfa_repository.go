package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// MFARepository defines the interface for MFA operations
type MFARepository interface {
	CreateMFAMethod(ctx context.Context, method *models.MFAMethod) error
	GetMFAMethod(ctx context.Context, methodID uuid.UUID) (*models.MFAMethod, error)
	GetUserMFAMethods(ctx context.Context, userID uuid.UUID) ([]*models.MFAMethod, error)
	UpdateMFAMethod(ctx context.Context, method *models.MFAMethod) error
	DeleteMFAMethod(ctx context.Context, methodID uuid.UUID) error

	CreateMFAAttempt(ctx context.Context, attempt *models.MFAAttempt) error
	GetMFAAttempts(ctx context.Context, userID uuid.UUID, limit int) ([]*models.MFAAttempt, error)
}

// SessionRepository defines the interface for session operations
type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.UserSession) error
	GetSessionByID(ctx context.Context, sessionID string) (*models.UserSession, error)
	GetUserActiveSessions(ctx context.Context, userID uuid.UUID) ([]*models.UserSession, error)
	UpdateSession(ctx context.Context, session *models.UserSession) error
	UpdateSessionActivity(ctx context.Context, sessionID string) error
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID, exceptSessionID string) error
	DeleteExpiredSessions(ctx context.Context) error
}

// mfaRepository implements MFARepository
type mfaRepository struct {
	db *gorm.DB
}

// NewMFARepository creates a new MFA repository
func NewMFARepository(db *gorm.DB) MFARepository {
	return &mfaRepository{db: db}
}

// CreateMFAMethod creates a new MFA method
func (r *mfaRepository) CreateMFAMethod(ctx context.Context, method *models.MFAMethod) error {
	return r.db.WithContext(ctx).Create(method).Error
}

// GetMFAMethod gets an MFA method by ID
func (r *mfaRepository) GetMFAMethod(ctx context.Context, methodID uuid.UUID) (*models.MFAMethod, error) {
	var method models.MFAMethod
	err := r.db.WithContext(ctx).Where("id = ?", methodID).First(&method).Error
	return &method, err
}

// GetUserMFAMethods gets all MFA methods for a user
func (r *mfaRepository) GetUserMFAMethods(ctx context.Context, userID uuid.UUID) ([]*models.MFAMethod, error) {
	var methods []*models.MFAMethod
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&methods).Error
	return methods, err
}

// UpdateMFAMethod updates an MFA method
func (r *mfaRepository) UpdateMFAMethod(ctx context.Context, method *models.MFAMethod) error {
	return r.db.WithContext(ctx).Save(method).Error
}

// DeleteMFAMethod deletes an MFA method
func (r *mfaRepository) DeleteMFAMethod(ctx context.Context, methodID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.MFAMethod{}, "id = ?", methodID).Error
}

// CreateMFAAttempt creates an MFA attempt record
func (r *mfaRepository) CreateMFAAttempt(ctx context.Context, attempt *models.MFAAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

// GetMFAAttempts gets MFA attempts for a user
func (r *mfaRepository) GetMFAAttempts(ctx context.Context, userID uuid.UUID, limit int) ([]*models.MFAAttempt, error) {
	var attempts []*models.MFAAttempt
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&attempts).Error
	return attempts, err
}

// sessionRepository implements SessionRepository
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

// CreateSession creates a new user session
func (r *sessionRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByID gets a session by ID
func (r *sessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error
	return &session, err
}

// GetUserActiveSessions gets all active sessions for a user
func (r *sessionRepository) GetUserActiveSessions(ctx context.Context, userID uuid.UUID) ([]*models.UserSession, error) {
	var sessions []*models.UserSession
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true AND expires_at > ?",
		userID, ctx.Value("now")).Find(&sessions).Error
	return sessions, err
}

// UpdateSession updates a session
func (r *sessionRepository) UpdateSession(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// UpdateSessionActivity updates session last activity
func (r *sessionRepository) UpdateSessionActivity(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("session_id = ?", sessionID).
		Update("last_activity", ctx.Value("now")).Error
}

// RevokeSession revokes a session
func (r *sessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	now := ctx.Value("now")
	return r.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"is_active":  false,
			"revoked_at": now,
		}).Error
}

// RevokeAllUserSessions revokes all sessions for a user except specified
func (r *sessionRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID, exceptSessionID string) error {
	now := ctx.Value("now")
	query := r.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("user_id = ? AND is_active = true", userID)

	if exceptSessionID != "" {
		query = query.Where("session_id != ?", exceptSessionID)
	}

	return query.Updates(map[string]interface{}{
		"is_active":  false,
		"revoked_at": now,
	}).Error
}

// DeleteExpiredSessions deletes expired sessions
func (r *sessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", ctx.Value("now")).
		Delete(&models.UserSession{}).Error
}

// KYCRRepository defines the interface for KYC operations
type KYCRRepository interface {
	CreateKYCRequest(ctx context.Context, request *models.KYCRequest) error
	GetKYCRequest(ctx context.Context, requestID string) (*models.KYCRequest, error)
	GetUserKYCRequests(ctx context.Context, userID uuid.UUID) ([]*models.KYCRequest, error)
	UpdateKYCRequest(ctx context.Context, request *models.KYCRequest) error
	GetPendingKYCRequests(ctx context.Context, limit, offset int) ([]*models.KYCRequest, error)
}

// kycRepository implements KYCRRepository
type kycRepository struct {
	db *gorm.DB
}

// NewKYCRepository creates a new KYC repository
func NewKYCRepository(db *gorm.DB) KYCRRepository {
	return &kycRepository{db: db}
}

// CreateKYCRequest creates a new KYC request
func (r *kycRepository) CreateKYCRequest(ctx context.Context, request *models.KYCRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

// GetKYCRequest gets a KYC request by ID
func (r *kycRepository) GetKYCRequest(ctx context.Context, requestID string) (*models.KYCRequest, error) {
	var request models.KYCRequest
	err := r.db.WithContext(ctx).Where("request_id = ?", requestID).First(&request).Error
	return &request, err
}

// GetUserKYCRequests gets all KYC requests for a user
func (r *kycRepository) GetUserKYCRequests(ctx context.Context, userID uuid.UUID) ([]*models.KYCRequest, error) {
	var requests []*models.KYCRequest
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&requests).Error
	return requests, err
}

// UpdateKYCRequest updates a KYC request
func (r *kycRepository) UpdateKYCRequest(ctx context.Context, request *models.KYCRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

// GetPendingKYCRequests gets pending KYC requests for review
func (r *kycRepository) GetPendingKYCRequests(ctx context.Context, limit, offset int) ([]*models.KYCRequest, error) {
	var requests []*models.KYCRequest
	err := r.db.WithContext(ctx).Where("status = ?", "pending").
		Order("created_at ASC").Limit(limit).Offset(offset).Find(&requests).Error
	return requests, err
}
