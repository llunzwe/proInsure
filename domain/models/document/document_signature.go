package document

import (
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ESignature represents an electronic signature on a document
type ESignature struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DocumentID uuid.UUID `gorm:"type:uuid;not null" json:"document_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Embedded structs for organization
	SignatureMetadata
	SignatureVerification
	SignatureLifecycle

	AuditTrail string    `json:"audit_trail,omitempty"` // JSON audit log
	Metadata   string    `json:"metadata,omitempty"`    // JSON for additional data
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Note: Relationships are defined in parent models package
	// These include: Document, User
}

// TableName specifies the table name
func (ESignature) TableName() string {
	return "e_signatures"
}

// BeforeCreate hook for UUID generation
func (e *ESignature) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.Status == "" {
		e.Status = SignaturePending
	}
	return nil
}
