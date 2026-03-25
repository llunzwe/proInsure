package models

import (
	"time"

	"github.com/google/uuid"
	"smartsure/pkg/database"
)

// MFAProvider represents MFA provider types
type MFAProvider string

const (
	MFAProviderTOTP    MFAProvider = "totp"
	MFAProviderSMS     MFAProvider = "sms"
	MFAProviderEmail   MFAProvider = "email"
	MFAProviderPush    MFAProvider = "push"
	MFAProviderHardware MFAProvider = "hardware"
)

// MFAStatus represents MFA status
type MFAStatus string

const (
	MFAStatusActive   MFAStatus = "active"
	MFAStatusInactive MFAStatus = "inactive"
	MFAStatusPending  MFAStatus = "pending"
	MFAStatusBlocked  MFAStatus = "blocked"
)

// MFAMethod represents a user's MFA method
type MFAMethod struct {
	database.BaseModel

	UserID       uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider     MFAProvider `gorm:"type:varchar(20);not null" json:"provider"`
	Status       MFAStatus   `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	IsDefault    bool        `gorm:"default:false" json:"is_default"`

	// Provider-specific data
	Secret       string     `json:"-"` // Encrypted TOTP secret
	PhoneNumber  string     `gorm:"type:varchar(20)" json:"phone_number,omitempty"`
	Email        string     `gorm:"type:varchar(255)" json:"email,omitempty"`
	DeviceToken  string     `json:"device_token,omitempty"`
	HardwareKey  string     `json:"hardware_key,omitempty"`

	// Verification data
	LastUsed     *time.Time `json:"last_used,omitempty"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	FailureCount int        `gorm:"default:0" json:"failure_count"`

	// Backup codes
	BackupCodes  []string   `gorm:"type:jsonb" json:"-"` // Encrypted backup codes

	// Metadata
	UserAgent    string     `json:"user_agent,omitempty"`
	IPAddress    string     `json:"ip_address,omitempty"`

	// Relationships
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// MFAAttempt represents an MFA verification attempt
type MFAAttempt struct {
	database.BaseModel

	UserID       uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	MethodID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"method_id"`
	Provider     MFAProvider `gorm:"type:varchar(20);not null" json:"provider"`

	Code         string      `json:"-"` // Encrypted verification code
	CodeHash     string      `gorm:"not null" json:"-"` // Hashed for verification

	IPAddress    string      `gorm:"not null" json:"ip_address"`
	UserAgent    string      `json:"user_agent,omitempty"`

	Status       string      `gorm:"type:varchar(20);not null" json:"status"` // success, failed, expired
	ExpiresAt    time.Time   `gorm:"not null" json:"expires_at"`

	// Relationships
	User         User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Method       MFAMethod   `gorm:"foreignKey:MethodID" json:"method,omitempty"`
}

// UserSession represents an authenticated user session
type UserSession struct {
	database.BaseModel

	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	SessionID    string    `gorm:"uniqueIndex;not null" json:"session_id"`

	IPAddress    string    `gorm:"not null" json:"ip_address"`
	UserAgent    string    `json:"user_agent,omitempty"`
	DeviceFingerprint string `json:"device_fingerprint,omitempty"`

	AccessToken  string    `gorm:"uniqueIndex" json:"access_token,omitempty"`
	RefreshToken string    `gorm:"uniqueIndex" json:"refresh_token,omitempty"`

	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	LastActivity time.Time `gorm:"not null" json:"last_activity"`

	IsActive     bool      `gorm:"default:true;index" json:"is_active"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	RevokedBy    *uuid.UUID `gorm:"type:uuid" json:"revoked_by,omitempty"`

	// Security tracking
	LoginAttempts int       `gorm:"default:0" json:"login_attempts"`
	SuspiciousActivity bool `gorm:"default:false" json:"suspicious_activity"`

	// Metadata
	Location     string    `json:"location,omitempty"`
	Timezone     string    `json:"timezone,omitempty"`

	// Relationships
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// KYCRequest represents a Know Your Customer verification request
type KYCRequest struct {
	database.BaseModel

	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RequestID    string    `gorm:"uniqueIndex;not null" json:"request_id"`

	Status       string    `gorm:"type:varchar(30);not null;default:'pending'" json:"status"`
	Provider     string    `gorm:"type:varchar(50)" json:"provider,omitempty"` // jumio, onfido, etc.

	// Personal information
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	Nationality  string    `json:"nationality,omitempty"`

	// Document information
	DocumentType string    `gorm:"type:varchar(30)" json:"document_type"` // passport, id_card, drivers_license
	DocumentNumber string  `json:"document_number,omitempty"`

	// File references
	DocumentFrontURL string `json:"document_front_url,omitempty"`
	DocumentBackURL  string `json:"document_back_url,omitempty"`
	SelfieURL        string `json:"selfie_url,omitempty"`

	// Verification results
	VerificationScore float64   `json:"verification_score,omitempty"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	RejectedReason    string    `json:"rejected_reason,omitempty"`

	// Provider response
	ProviderResponse string    `gorm:"type:jsonb" json:"provider_response,omitempty"`
	ProviderReference string   `json:"provider_reference,omitempty"`

	// Audit
	SubmittedBy  uuid.UUID `gorm:"type:uuid;not null" json:"submitted_by"`
	ReviewedBy   *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`

	// Relationships
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for MFAMethod
func (MFAMethod) TableName() string {
	return "mfa_methods"
}

// TableName returns the table name for MFAAttempt
func (MFAAttempt) TableName() string {
	return "mfa_attempts"
}

// TableName returns the table name for UserSession
func (UserSession) TableName() string {
	return "user_sessions"
}

// TableName returns the table name for KYCRequest
func (KYCRequest) TableName() string {
	return "kyc_requests"
}
