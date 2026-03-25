package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AsyncOperation represents an asynchronous operation being tracked
type AsyncOperation struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OperationType string         `gorm:"type:varchar(100);not null;index" json:"operation_type"` // "blockchain_submit", "claim_assessment", "fraud_check"
	EntityType    string         `gorm:"type:varchar(50);not null;index" json:"entity_type"` // "policy", "claim", "device"
	EntityID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"entity_id"`
	Status        string         `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"` // "pending", "processing", "completed", "failed"
	Progress      int            `gorm:"type:integer;default:0" json:"progress"` // 0-100
	Error         string         `gorm:"type:text" json:"error"`
	Result        string         `gorm:"type:text" json:"result"` // JSON result
	StartedAt     *time.Time     `gorm:"type:timestamp" json:"started_at"`
	CompletedAt   *time.Time     `gorm:"type:timestamp" json:"completed_at"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name
func (AsyncOperation) TableName() string {
	return "async_operations"
}

// MarkProcessing marks operation as processing
func (o *AsyncOperation) MarkProcessing() {
	o.Status = "processing"
	now := time.Now()
	o.StartedAt = &now
	o.Progress = 10
}

// MarkCompleted marks operation as completed
func (o *AsyncOperation) MarkCompleted(result string) {
	o.Status = "completed"
	o.Progress = 100
	now := time.Now()
	o.CompletedAt = &now
	o.Result = result
	o.Error = ""
}

// MarkFailed marks operation as failed
func (o *AsyncOperation) MarkFailed(err error) {
	o.Status = "failed"
	now := time.Now()
	o.CompletedAt = &now
	if err != nil {
		o.Error = err.Error()
	}
}

// UpdateProgress updates operation progress
func (o *AsyncOperation) UpdateProgress(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	o.Progress = progress
}

