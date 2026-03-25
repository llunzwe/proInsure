package verification

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// IMEIVerifier provides IMEI verification services
type IMEIVerifier struct {
	config *VerificationConfig
}

// IdentityVerifier provides identity verification services
type IdentityVerifier struct {
	config *VerificationConfig
}

// EmailVerifier provides email verification services
type EmailVerifier struct {
	config *VerificationConfig
}

// PhoneVerifier provides phone verification services
type PhoneVerifier struct {
	config *VerificationConfig
}

// VerificationConfig holds verification configuration
type VerificationConfig struct {
	IMEIValidationEnabled  bool
	IdentityCheckEnabled   bool
	EmailValidationEnabled bool
	PhoneValidationEnabled bool
	MaxRetries             int
	TimeoutDuration        time.Duration
}

// VerificationResult represents the result of a verification
type VerificationResult struct {
	ID        uuid.UUID              `json:"id"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Score     float64                `json:"score"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
}

// NewIMEIVerifier creates a new IMEI verifier
func NewIMEIVerifier(config *VerificationConfig) *IMEIVerifier {
	return &IMEIVerifier{
		config: config,
	}
}

// NewIdentityVerifier creates a new identity verifier
func NewIdentityVerifier(config *VerificationConfig) *IdentityVerifier {
	return &IdentityVerifier{
		config: config,
	}
}

// NewEmailVerifier creates a new email verifier
func NewEmailVerifier(config *VerificationConfig) *EmailVerifier {
	return &EmailVerifier{
		config: config,
	}
}

// NewPhoneVerifier creates a new phone verifier
func NewPhoneVerifier(config *VerificationConfig) *PhoneVerifier {
	return &PhoneVerifier{
		config: config,
	}
}

// VerifyIMEI verifies an IMEI number
func (v *IMEIVerifier) VerifyIMEI(ctx context.Context, imei string) (*VerificationResult, error) {
	if !v.config.IMEIValidationEnabled {
		return &VerificationResult{
			ID:        uuid.New(),
			Type:      "IMEI",
			Status:    "SKIPPED",
			Score:     1.0,
			Message:   "IMEI validation disabled",
			Timestamp: time.Now(),
		}, nil
	}

	// Basic IMEI validation
	if len(imei) != 15 {
		return &VerificationResult{
			ID:        uuid.New(),
			Type:      "IMEI",
			Status:    "FAILED",
			Score:     0.0,
			Message:   "Invalid IMEI length",
			Timestamp: time.Now(),
		}, fmt.Errorf("invalid IMEI length: %d", len(imei))
	}

	// Check if IMEI contains only digits
	if matched, _ := regexp.MatchString(`^\d{15}$`, imei); !matched {
		return &VerificationResult{
			ID:        uuid.New(),
			Type:      "IMEI",
			Status:    "FAILED",
			Score:     0.0,
			Message:   "IMEI contains non-digit characters",
			Timestamp: time.Now(),
		}, fmt.Errorf("invalid IMEI format")
	}

	return &VerificationResult{
		ID:        uuid.New(),
		Type:      "IMEI",
		Status:    "VERIFIED",
		Score:     1.0,
		Message:   "IMEI verification successful",
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"imei": imei,
		},
	}, nil
}

// VerifyIdentity verifies identity information
func (v *IdentityVerifier) VerifyIdentity(ctx context.Context, identityData map[string]interface{}) (*VerificationResult, error) {
	if !v.config.IdentityCheckEnabled {
		return &VerificationResult{
			ID:        uuid.New(),
			Type:      "IDENTITY",
			Status:    "SKIPPED",
			Score:     1.0,
			Message:   "Identity verification disabled",
			Timestamp: time.Now(),
		}, nil
	}

	// Basic identity validation
	score := 0.0
	checks := 0

	if name, ok := identityData["name"].(string); ok && name != "" {
		score += 0.3
		checks++
	}

	if dob, ok := identityData["date_of_birth"].(string); ok && dob != "" {
		score += 0.3
		checks++
	}

	if id, ok := identityData["id_number"].(string); ok && id != "" {
		score += 0.4
		checks++
	}

	if checks == 0 {
		return &VerificationResult{
			ID:        uuid.New(),
			Type:      "IDENTITY",
			Status:    "FAILED",
			Score:     0.0,
			Message:   "No identity data provided",
			Timestamp: time.Now(),
		}, fmt.Errorf("no identity data provided")
	}

	status := "VERIFIED"
	if score < 0.7 {
		status = "PARTIAL"
	}

	return &VerificationResult{
		ID:        uuid.New(),
		Type:      "IDENTITY",
		Status:    status,
		Score:     score,
		Message:   "Identity verification completed",
		Timestamp: time.Now(),
		Details:   identityData,
	}, nil
}

// VerifyEmail verifies an email address
func (v *EmailVerifier) VerifyEmail(ctx context.Context, email string) error {
	if !v.config.EmailValidationEnabled {
		return nil
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// VerifyPhone verifies a phone number
func (v *PhoneVerifier) VerifyPhone(ctx context.Context, phone string) error {
	if !v.config.PhoneValidationEnabled {
		return nil
	}

	// Basic phone validation - remove spaces and special characters
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "+", "")

	if len(cleanPhone) < 10 || len(cleanPhone) > 15 {
		return fmt.Errorf("invalid phone number length: %s", phone)
	}

	// Check if phone contains only digits
	if matched, _ := regexp.MatchString(`^\d+$`, cleanPhone); !matched {
		return fmt.Errorf("invalid phone number format: %s", phone)
	}

	return nil
}
