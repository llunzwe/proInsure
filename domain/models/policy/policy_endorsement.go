package policy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/pkg/database"
)

// PolicyEndorsement represents mid-term policy changes
type PolicyEndorsement struct {
	database.BaseModel
	PolicyID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"policy_id"`
	EndorsementNumber string     `gorm:"uniqueIndex;not null" json:"endorsement_number"`
	EndorsementType   string     `gorm:"not null" json:"endorsement_type"` // coverage_change, sum_insured_change, beneficiary_change, device_replacement
	RequestDate       time.Time  `gorm:"not null" json:"request_date"`
	EffectiveDate     time.Time  `gorm:"not null" json:"effective_date"`
	Status            string     `gorm:"not null;default:'pending'" json:"status"` // pending, approved, rejected, implemented
	RequestedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"requested_by"`
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovalDate      *time.Time `json:"approval_date"`
	RejectionReason   string     `json:"rejection_reason"`

	// Financial adjustments
	PremiumAdjustment    float64 `gorm:"default:0" json:"premium_adjustment"`
	SumInsuredAdjustment float64 `gorm:"default:0" json:"sum_insured_adjustment"`
	ProcessingFee        float64 `gorm:"default:0" json:"processing_fee"`

	// Change details
	ChangeDetails  string `gorm:"type:json" json:"change_details"`  // JSON object with specific changes
	PreviousValues string `gorm:"type:json" json:"previous_values"` // JSON object with old values
	NewValues      string `gorm:"type:json" json:"new_values"`      // JSON object with new values

	// Documentation
	DocumentsRequired string `gorm:"type:json" json:"documents_required"` // JSON array
	DocumentsReceived string `gorm:"type:json" json:"documents_received"` // JSON array
	DigitalSignature  string `json:"digital_signature"`

	// Implementation tracking
	ImplementationDate *time.Time `json:"implementation_date"`
	ImplementedBy      *uuid.UUID `gorm:"type:uuid" json:"implemented_by"`
	ValidationStatus   string     `json:"validation_status"`
	ValidationErrors   string     `gorm:"type:json" json:"validation_errors"` // JSON array

	// Metadata
	Notes      string `gorm:"type:text" json:"notes"`
	AuditTrail string `gorm:"type:json" json:"audit_trail"` // JSON array of actions
	ClientIP   string `json:"client_ip"`
	UserAgent  string `json:"user_agent"`

	// Relationships (loaded via service layer to avoid circular imports)
	// Policy is loaded via service layer using PolicyID
	// RequestedByUser is loaded via service layer using RequestedBy
	// ApprovedByUser is loaded via service layer using ApprovedBy
}

// BeforeCreate handles pre-creation logic (UUID generation handled by BaseModel)
func (pe *PolicyEndorsement) BeforeCreate(tx *gorm.DB) error {
	if err := pe.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Generate endorsement number if not provided
	if pe.EndorsementNumber == "" {
		pe.EndorsementNumber = pe.generateEndorsementNumber()
	}

	// Set default status
	if pe.Status == "" {
		pe.Status = "pending"
	}

	// Set request date if not provided
	if pe.RequestDate.IsZero() {
		pe.RequestDate = time.Now()
	}

	return nil
}

// generateEndorsementNumber generates a unique endorsement number
func (pe *PolicyEndorsement) generateEndorsementNumber() string {
	timestamp := time.Now().Format("20060102150405")
	randomStr := pe.ID.String()[:8]
	return "END-" + timestamp + "-" + randomStr
}

// IsValid checks if the endorsement is valid
func (pe *PolicyEndorsement) IsValid() bool {
	return pe.Status == "approved" || pe.Status == "implemented"
}

// IsPending checks if the endorsement is pending
func (pe *PolicyEndorsement) IsPending() bool {
	return pe.Status == "pending"
}

// IsImplemented checks if the endorsement has been implemented
func (pe *PolicyEndorsement) IsImplemented() bool {
	return pe.Status == "implemented" && pe.ImplementationDate != nil
}

// CanBeApproved checks if the endorsement can be approved
func (pe *PolicyEndorsement) CanBeApproved() bool {
	return pe.Status == "pending" &&
		pe.EffectiveDate.After(time.Now()) &&
		pe.DocumentsReceived != "" // Has required documents
}

// CanBeImplemented checks if the endorsement can be implemented
func (pe *PolicyEndorsement) CanBeImplemented() bool {
	return pe.Status == "approved" &&
		!pe.EffectiveDate.After(time.Now()) // Effective date has passed
}

// GetNetPremiumChange calculates the net premium change
func (pe *PolicyEndorsement) GetNetPremiumChange() float64 {
	return pe.PremiumAdjustment - pe.ProcessingFee
}

// TableName returns the table name
func (pe *PolicyEndorsement) TableName() string {
	return "policy_endorsements"
}
