package document

import (
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentTemplate represents reusable document templates
type DocumentTemplate struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`

	// Embedded structs for organization
	TemplateConfiguration
	TemplateContent

	Metadata  string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Note: Relationships are defined in parent models package
	// These include: DocumentGenerations, CreatedByUser
}

// TableName specifies the table name
func (DocumentTemplate) TableName() string {
	return "document_templates"
}

// BeforeCreate hook for UUID generation
func (d *DocumentTemplate) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	if d.Language == "" {
		d.Language = DefaultLanguage
	}
	return nil
}
