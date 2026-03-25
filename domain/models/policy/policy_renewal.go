package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyRenewal represents a policy renewal process
type PolicyRenewal struct {
	database.BaseModel
	OriginalPolicyID uuid.UUID  `gorm:"type:uuid;not null;index" json:"original_policy_id"`
	NewPolicyID      *uuid.UUID `gorm:"type:uuid;index" json:"new_policy_id"`
	RenewalNumber    string     `gorm:"uniqueIndex;not null" json:"renewal_number"`
	Status           string     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	RenewalType      string     `gorm:"type:varchar(50)" json:"renewal_type"` // auto, manual, forced

	// Dates
	RenewalDate    time.Time  `json:"renewal_date"`
	EffectiveDate  time.Time  `json:"effective_date"`
	ExpirationDate time.Time  `json:"expiration_date"`
	ProcessedDate  *time.Time `json:"processed_date"`

	// Pricing
	OriginalPremium float64 `json:"original_premium"`
	ProposedPremium float64 `json:"proposed_premium"`
	FinalPremium    float64 `json:"final_premium"`
	PremiumChange   float64 `json:"premium_change"`
	ChangePercent   float64 `json:"change_percent"`

	// Risk Assessment
	RiskScoreChange float64 `json:"risk_score_change"`
	ClaimsInPeriod  int     `json:"claims_in_period"`
	LossRatio       float64 `json:"loss_ratio"`

	// Underwriting
	RequiresReview    bool       `gorm:"default:false" json:"requires_review"`
	ReviewReason      string     `json:"review_reason"`
	UnderwriterID     *uuid.UUID `gorm:"type:uuid" json:"underwriter_id"`
	UnderwritingNotes string     `json:"underwriting_notes"`

	// Communication
	NotificationsSent int        `json:"notifications_sent"`
	LastNotification  *time.Time `json:"last_notification"`
	CustomerResponse  string     `json:"customer_response"`
	ResponseDate      *time.Time `json:"response_date"`

	// Retention
	RetentionOffer string `gorm:"type:json" json:"retention_offer"`
	OfferAccepted  *bool  `json:"offer_accepted"`
	ChurnReason    string `json:"churn_reason"`

	// Relationships
	// Note: OriginalPolicy relationship is handled through embedding in the main Policy struct
	NewPolicy      interface{} `gorm:"foreignKey:NewPolicyID" json:"new_policy,omitempty"` // *models.Policy
}

// TableName returns the table name
func (PolicyRenewal) TableName() string {
	return "policy_renewals"
}

// NeedsCustomerAction checks if renewal needs customer action
func (pr *PolicyRenewal) NeedsCustomerAction() bool {
	return pr.Status == "pending" && pr.RenewalType == "manual"
}

// IsOverdue checks if renewal is overdue
func (pr *PolicyRenewal) IsOverdue() bool {
	return pr.Status == "pending" && time.Now().After(pr.RenewalDate)
}

// GetPremiumChangePercent calculates the percentage change in premium
func (pr *PolicyRenewal) GetPremiumChangePercent() float64 {
	if pr.OriginalPremium <= 0 {
		return 0
	}
	return ((pr.ProposedPremium - pr.OriginalPremium) / pr.OriginalPremium) * 100
}

// IsAutoRenewalEligible checks if the policy can be auto-renewed
func (pr *PolicyRenewal) IsAutoRenewalEligible() bool {
	return pr.RenewalType == "auto" &&
		!pr.RequiresReview &&
		pr.LossRatio < 100 &&
		pr.ClaimsInPeriod < 3
}

// ShouldSendReminder checks if a reminder should be sent
func (pr *PolicyRenewal) ShouldSendReminder() bool {
	if pr.Status != "pending" || pr.CustomerResponse != "" {
		return false
	}

	// Send reminders at specific intervals
	daysSinceLastNotification := 0
	if pr.LastNotification != nil {
		daysSinceLastNotification = int(time.Since(*pr.LastNotification).Hours() / 24)
	}

	// Send at 30, 15, 7, and 3 days before renewal
	daysUntilRenewal := int(time.Until(pr.RenewalDate).Hours() / 24)

	switch daysUntilRenewal {
	case 30, 15, 7, 3:
		return daysSinceLastNotification >= 3 // Don't send more than once every 3 days
	default:
		return false
	}
}

// HasSignificantPriceIncrease checks for significant premium increase
func (pr *PolicyRenewal) HasSignificantPriceIncrease() bool {
	changePercent := pr.GetPremiumChangePercent()
	return changePercent > 20 // More than 20% increase
}

// CanBeProcessed checks if the renewal can be processed
func (pr *PolicyRenewal) CanBeProcessed() bool {
	return pr.Status == "pending" &&
		(pr.RenewalType == "auto" || pr.CustomerResponse == "accepted") &&
		!pr.RequiresReview
}

// GetDaysUntilRenewal returns days until renewal date
func (pr *PolicyRenewal) GetDaysUntilRenewal() int {
	if time.Now().After(pr.RenewalDate) {
		return 0
	}
	return int(time.Until(pr.RenewalDate).Hours() / 24)
}
