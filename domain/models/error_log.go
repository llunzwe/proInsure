package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrorLog represents an error that occurred in the system
type ErrorLog struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ErrorType     string         `gorm:"type:varchar(100);not null;index" json:"error_type"` // "blockchain", "database", "external_service", "validation", "business_logic"
	Severity      string         `gorm:"type:varchar(50);not null;index" json:"severity"` // "low", "medium", "high", "critical"
	Service       string         `gorm:"type:varchar(100);not null;index" json:"service"` // Service name
	Operation     string         `gorm:"type:varchar(200);not null;index" json:"operation"` // Operation name
	ErrorCode     string         `gorm:"type:varchar(50)" json:"error_code"`
	ErrorMessage  string         `gorm:"type:text;not null" json:"error_message"`
	StackTrace    string         `gorm:"type:text" json:"stack_trace"`
	Context       string         `gorm:"type:jsonb" json:"context"` // Additional context as JSON
	UserID        *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	EntityType    string         `gorm:"type:varchar(50);index" json:"entity_type"`
	EntityID      *uuid.UUID     `gorm:"type:uuid;index" json:"entity_id,omitempty"`
	Resolved      bool           `gorm:"type:boolean;default:false;index" json:"resolved"`
	ResolvedAt    *time.Time     `gorm:"type:timestamp" json:"resolved_at"`
	ResolvedBy    *uuid.UUID     `gorm:"type:uuid" json:"resolved_by,omitempty"`
	Resolution    string         `gorm:"type:text" json:"resolution"`
	OccurrenceCount int           `gorm:"type:integer;default:1" json:"occurrence_count"`
	FirstOccurredAt time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"first_occurred_at"`
	LastOccurredAt  time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"last_occurred_at"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name
func (ErrorLog) TableName() string {
	return "error_logs"
}

// IsCritical returns true if error is critical
func (e *ErrorLog) IsCritical() bool {
	return e.Severity == "critical"
}

// IncrementOccurrence increments occurrence count
func (e *ErrorLog) IncrementOccurrence() {
	e.OccurrenceCount++
	e.LastOccurredAt = time.Now()
}

