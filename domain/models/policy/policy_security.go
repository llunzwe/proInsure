package policy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================
// POLICY SECURITY SERVICE
// ============================================

// PolicySecurityService provides security features for policies
type PolicySecurityService struct {
	db            *gorm.DB
	encryptionKey []byte
}

// NewPolicySecurityService creates a new security service
func NewPolicySecurityService(db *gorm.DB, encryptionKey string) *PolicySecurityService {
	// Create 32-byte key from provided string
	hash := sha256.Sum256([]byte(encryptionKey))
	return &PolicySecurityService{
		db:            db,
		encryptionKey: hash[:],
	}
}

// ============================================
// ENCRYPTION METHODS
// ============================================

// EncryptSensitiveFields encrypts sensitive policy data
func (pss *PolicySecurityService) EncryptSensitiveFields(policy interface{}) error {
	// TODO: Implement proper type assertion and encryption logic
	// Placeholder implementation to avoid compilation errors
	return nil
}

// DecryptSensitiveFields decrypts sensitive policy data
func (pss *PolicySecurityService) DecryptSensitiveFields(policy interface{}) error {
	// TODO: Implement proper type assertion and decryption logic
	// Placeholder implementation to avoid compilation errors
	return nil
}

// decrypt decrypts a base64 encoded ciphertext string
func (pss *PolicySecurityService) decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(pss.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertextBytes := data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// isEncrypted checks if a string is already encrypted
func (pss *PolicySecurityService) isEncrypted(s string) bool {
	// Check if string is valid base64
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil && len(s) > 32 // Encrypted strings are longer
}

// ============================================
// AUDIT METHODS
// ============================================

// PolicyAuditLog represents a policy audit trail entry
type PolicyAuditLog struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	PolicyID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"policy_id"`
	Action    string          `gorm:"type:varchar(50);not null" json:"action"`
	UserID    uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	UserName  string          `gorm:"type:varchar(100)" json:"user_name"`
	UserRole  string          `gorm:"type:varchar(50)" json:"user_role"`
	Timestamp time.Time       `gorm:"not null;index" json:"timestamp"`
	IPAddress string          `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string          `gorm:"type:text" json:"user_agent"`
	SessionID string          `gorm:"type:varchar(100)" json:"session_id"`
	Changes   json.RawMessage `gorm:"type:jsonb" json:"changes"`
	OldValues json.RawMessage `gorm:"type:jsonb" json:"old_values"`
	NewValues json.RawMessage `gorm:"type:jsonb" json:"new_values"`
	Result    string          `gorm:"type:varchar(20)" json:"result"` // success, failure, partial
	ErrorMsg  string          `gorm:"type:text" json:"error_message"`
	Hash      string          `gorm:"type:varchar(64);unique" json:"hash"` // SHA256 hash for integrity
	PrevHash  string          `gorm:"type:varchar(64)" json:"prev_hash"`   // Previous entry hash (blockchain-style)
}

// LogPolicyAction creates an audit log entry for a policy action
func (pss *PolicySecurityService) LogPolicyAction(
	policyID uuid.UUID,
	action string,
	userID uuid.UUID,
	userName, userRole string,
	ipAddress, userAgent, sessionID string,
	oldValues, newValues interface{},
	result string,
	errorMsg string,
) error {
	// Convert values to JSON
	var oldJSON, newJSON json.RawMessage
	if oldValues != nil {
		data, err := json.Marshal(oldValues)
		if err != nil {
			return err
		}
		oldJSON = data
	}
	if newValues != nil {
		data, err := json.Marshal(newValues)
		if err != nil {
			return err
		}
		newJSON = data
	}

	// Calculate changes
	changes := pss.calculateChanges(oldJSON, newJSON)

	// Get previous hash for blockchain-style chaining
	var prevHash string
	var lastEntry PolicyAuditLog
	if err := pss.db.Order("timestamp DESC").First(&lastEntry).Error; err == nil {
		prevHash = lastEntry.Hash
	}

	// Create audit log entry
	auditLog := PolicyAuditLog{
		ID:        uuid.New(),
		PolicyID:  policyID,
		Action:    action,
		UserID:    userID,
		UserName:  userName,
		UserRole:  userRole,
		Timestamp: time.Now(),
		IPAddress: ipAddress,
		UserAgent: userAgent,
		SessionID: sessionID,
		Changes:   changes,
		OldValues: oldJSON,
		NewValues: newJSON,
		Result:    result,
		ErrorMsg:  errorMsg,
		PrevHash:  prevHash,
	}

	// Calculate hash for this entry
	auditLog.Hash = pss.calculateHash(auditLog)

	// Save to database
	return pss.db.Create(&auditLog).Error
}

// calculateChanges calculates the differences between old and new values
func (pss *PolicySecurityService) calculateChanges(oldJSON, newJSON json.RawMessage) json.RawMessage {
	if oldJSON == nil || newJSON == nil {
		return nil
	}

	var oldMap, newMap map[string]interface{}
	json.Unmarshal(oldJSON, &oldMap)
	json.Unmarshal(newJSON, &newMap)

	changes := make(map[string]map[string]interface{})

	// Find modified and added fields
	for key, newValue := range newMap {
		oldValue, exists := oldMap[key]
		if !exists {
			changes[key] = map[string]interface{}{
				"action": "added",
				"new":    newValue,
			}
		} else if !pss.valuesEqual(oldValue, newValue) {
			changes[key] = map[string]interface{}{
				"action": "modified",
				"old":    oldValue,
				"new":    newValue,
			}
		}
	}

	// Find removed fields
	for key, oldValue := range oldMap {
		if _, exists := newMap[key]; !exists {
			changes[key] = map[string]interface{}{
				"action": "removed",
				"old":    oldValue,
			}
		}
	}

	result, _ := json.Marshal(changes)
	return result
}

// calculateHash generates a SHA256 hash for audit log integrity
func (pss *PolicySecurityService) calculateHash(log PolicyAuditLog) string {
	// Create a string representation of important fields
	data := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
		log.PolicyID.String(),
		log.Action,
		log.UserID.String(),
		log.Timestamp.Format(time.RFC3339),
		string(log.Changes),
		log.Result,
		log.PrevHash,
		pss.encryptionKey,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// valuesEqual compares two interface{} values for equality
func (pss *PolicySecurityService) valuesEqual(a, b interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

// VerifyAuditTrail verifies the integrity of the audit trail
func (pss *PolicySecurityService) VerifyAuditTrail(policyID uuid.UUID) (bool, error) {
	var logs []PolicyAuditLog
	err := pss.db.Where("policy_id = ?", policyID).
		Order("timestamp ASC").
		Find(&logs).Error
	if err != nil {
		return false, err
	}

	for i, log := range logs {
		// Verify hash
		expectedHash := pss.calculateHash(log)
		if log.Hash != expectedHash {
			return false, fmt.Errorf("hash mismatch at entry %d", i)
		}

		// Verify chain (except for first entry)
		if i > 0 && log.PrevHash != logs[i-1].Hash {
			return false, fmt.Errorf("chain broken at entry %d", i)
		}
	}

	return true, nil
}

// ============================================
// ACCESS CONTROL
// ============================================

// PolicyAccessControl manages access permissions for policies
type PolicyAccessControl struct {
	db *gorm.DB
}

// NewPolicyAccessControl creates a new access control service
func NewPolicyAccessControl(db *gorm.DB) *PolicyAccessControl {
	return &PolicyAccessControl{db: db}
}

// CanRead checks if a user can read a policy
func (pac *PolicyAccessControl) CanRead(policy interface{}, userID uuid.UUID, userRole string) bool { // *models.Policy
	// TODO: Implement proper access control logic
	// Placeholder implementation to avoid compilation errors
	return false
}
func (pac *PolicyAccessControl) CanUpdate(policy interface{}, userID uuid.UUID, userRole string) bool { // *models.Policy
	// Only certain roles can update
	allowedRoles := []string{"admin", "underwriter", "policy_manager"}
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}

	// Agent can update certain fields
	// Note: Using reflection to avoid import cycle - policy parameter is interface{}
	// This would need proper type assertion in actual implementation
	_ = policy
	_ = userID
	return false
}


// ============================================
// DATA MASKING
// ============================================

// MaskSensitiveData masks sensitive information for display
func (pss *PolicySecurityService) MaskSensitiveData(policy interface{}, userRole string) interface{} { // *models.Policy
	// TODO: Implement proper data masking logic
	// Placeholder implementation to avoid compilation errors
	return policy
}

// maskString masks all but the last n characters of a string
func (pss *PolicySecurityService) maskString(s string, showLast int) string {
	if len(s) <= showLast {
		return s
	}
	masked := ""
	for i := 0; i < len(s)-showLast; i++ {
		masked += "*"
	}
	return masked + s[len(s)-showLast:]
}

// maskName masks parts of a name for privacy
func (pss *PolicySecurityService) maskName(name string) string {
	if name == "" {
		return ""
	}
	// Show only first letter and last letter of each word
	words := strings.Fields(name)
	maskedWords := []string{}
	for _, word := range words {
		if len(word) > 2 {
			maskedWords = append(maskedWords, string(word[0])+"***"+string(word[len(word)-1]))
		} else {
			maskedWords = append(maskedWords, "***")
		}
	}
	return strings.Join(maskedWords, " ")
}

// ============================================
// COMPLIANCE MONITORING
// ============================================

// CheckComplianceViolations checks for any compliance violations
func (pss *PolicySecurityService) CheckComplianceViolations(policy interface{}) []string { // *models.Policy
	violations := []string{}

	// TODO: Implement proper compliance checking using reflection
	// This avoids import cycle by using interface{} and reflection
	// In production, would use proper type assertion and field access
	
	_ = policy
	// Placeholder implementation - would check:
	// - KYC compliance for high-value policies (>$10,000)
	// - AML clearance for policies over $50,000
	// - Data privacy consent
	// - GDPR compliance for EU
	// - Sanctions checks
	// - Document requirements
	// - Signature requirements

	return violations
}
