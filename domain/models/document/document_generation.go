package document

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentGeneration tracks document generation requests from templates
type DocumentGeneration struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TemplateID          uuid.UUID  `gorm:"type:uuid;not null" json:"template_id"`
	DocumentID          *uuid.UUID `gorm:"type:uuid" json:"document_id,omitempty"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	RequestType         string     `gorm:"not null" json:"request_type"`             // instant, scheduled, batch
	Status              string     `gorm:"not null;default:'pending'" json:"status"` // pending, processing, completed, failed
	Variables           string     `gorm:"type:text" json:"variables"`               // JSON data for template variables
	OutputFormat        string     `gorm:"not null" json:"output_format"`            // pdf, html, docx
	DeliveryMethod      string     `json:"delivery_method,omitempty"`                // email, download, api
	DeliveryAddress     string     `json:"delivery_address,omitempty"`               // Email or API endpoint
	ProcessingStarted   *time.Time `json:"processing_started,omitempty"`
	ProcessingCompleted *time.Time `json:"processing_completed,omitempty"`
	ErrorMessage        string     `json:"error_message,omitempty"`
	RetryCount          int        `gorm:"default:0" json:"retry_count"`
	Metadata            string     `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Note: Relationships are defined in parent models package
	// These include: Template, Document, User
}

// TableName specifies the table name
func (DocumentGeneration) TableName() string {
	return "document_generations"
}

// BeforeCreate hook for UUID generation
func (d *DocumentGeneration) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	if d.Status == "" {
		d.Status = GenerationPending
	}
	if d.RequestType == "" {
		d.RequestType = GenerationInstant
	}
	return nil
}
