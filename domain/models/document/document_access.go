package document

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentAccess tracks document access permissions and logs
type DocumentAccess struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	DocumentID     uuid.UUID  `gorm:"type:uuid;not null" json:"document_id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	AccessType     string     `gorm:"not null" json:"access_type"` // view, download, edit, sign, share
	Permission     string     `gorm:"not null" json:"permission"`  // owner, editor, viewer, signer
	GrantedBy      uuid.UUID  `gorm:"type:uuid" json:"granted_by"`
	GrantedAt      time.Time  `json:"granted_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	AccessCount    int        `gorm:"default:0" json:"access_count"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
	IPAddress      string     `json:"ip_address,omitempty"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty"`
	RevokedBy      *uuid.UUID `gorm:"type:uuid" json:"revoked_by,omitempty"`
	RevokeReason   string     `json:"revoke_reason,omitempty"`
	Metadata       string     `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Note: Relationships are defined in parent models package
	// These include: Document, User, GrantedByUser, RevokedByUser
}

// TableName specifies the table name
func (DocumentAccess) TableName() string {
	return "document_accesses"
}

// BeforeCreate hook for UUID generation
func (d *DocumentAccess) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	if d.GrantedAt.IsZero() {
		d.GrantedAt = time.Now()
	}
	return nil
}
