package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyCommunicationPreference represents communication preferences for a policy
type PolicyCommunicationPreference struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	// Removed CustomerID - access via Policy.CustomerID to avoid duplication

	// Channels
	EmailEnabled    bool `gorm:"default:true" json:"email_enabled"`
	SMSEnabled      bool `gorm:"default:false" json:"sms_enabled"`
	PushEnabled     bool `gorm:"default:false" json:"push_enabled"`
	WhatsAppEnabled bool `gorm:"default:false" json:"whatsapp_enabled"`
	PhoneEnabled    bool `gorm:"default:false" json:"phone_enabled"`

	// Preferences
	PreferredChannel  string `gorm:"type:varchar(20)" json:"preferred_channel"`
	PreferredLanguage string `gorm:"default:'en'" json:"preferred_language"`
	PreferredTime     string `json:"preferred_time"` // morning, afternoon, evening
	Timezone          string `json:"timezone"`

	// Notification Types
	RenewalReminders bool `gorm:"default:true" json:"renewal_reminders"`
	PaymentReminders bool `gorm:"default:true" json:"payment_reminders"`
	ClaimUpdates     bool `gorm:"default:true" json:"claim_updates"`
	PolicyChanges    bool `gorm:"default:true" json:"policy_changes"`
	MarketingOffers  bool `gorm:"default:false" json:"marketing_offers"`
	Newsletters      bool `gorm:"default:false" json:"newsletters"`

	// Frequency
	ReminderFrequency string `json:"reminder_frequency"` // daily, weekly, monthly
	ReminderDays      int    `json:"reminder_days"`      // Days before due date

	// Opt-out
	OptOutAll    bool       `gorm:"default:false" json:"opt_out_all"`
	OptOutDate   *time.Time `json:"opt_out_date"`
	OptOutReason string     `json:"opt_out_reason"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
	// Customer accessed via Policy.Customer relationship - CustomerID field removed to avoid duplication
}

// TableName returns the table name
func (PolicyCommunicationPreference) TableName() string {
	return "policy_communication_preferences"
}

// ShouldNotify checks if customer should be notified for a type
func (pcp *PolicyCommunicationPreference) ShouldNotify(notificationType string) bool {
	if pcp.OptOutAll {
		return false
	}

	switch notificationType {
	case "renewal":
		return pcp.RenewalReminders
	case "payment":
		return pcp.PaymentReminders
	case "claim":
		return pcp.ClaimUpdates
	case "policy":
		return pcp.PolicyChanges
	case "marketing":
		return pcp.MarketingOffers
	case "newsletter":
		return pcp.Newsletters
	default:
		return true
	}
}

// GetActiveChannels returns list of active communication channels
func (pcp *PolicyCommunicationPreference) GetActiveChannels() []string {
	channels := []string{}

	if pcp.EmailEnabled {
		channels = append(channels, "email")
	}
	if pcp.SMSEnabled {
		channels = append(channels, "sms")
	}
	if pcp.PushEnabled {
		channels = append(channels, "push")
	}
	if pcp.WhatsAppEnabled {
		channels = append(channels, "whatsapp")
	}
	if pcp.PhoneEnabled {
		channels = append(channels, "phone")
	}

	return channels
}

// HasActiveChannels checks if any communication channel is enabled
func (pcp *PolicyCommunicationPreference) HasActiveChannels() bool {
	return pcp.EmailEnabled ||
		pcp.SMSEnabled ||
		pcp.PushEnabled ||
		pcp.WhatsAppEnabled ||
		pcp.PhoneEnabled
}

// GetBestChannel returns the best available channel
func (pcp *PolicyCommunicationPreference) GetBestChannel() string {
	// Return preferred channel if it's enabled
	switch pcp.PreferredChannel {
	case "email":
		if pcp.EmailEnabled {
			return "email"
		}
	case "sms":
		if pcp.SMSEnabled {
			return "sms"
		}
	case "push":
		if pcp.PushEnabled {
			return "push"
		}
	case "whatsapp":
		if pcp.WhatsAppEnabled {
			return "whatsapp"
		}
	case "phone":
		if pcp.PhoneEnabled {
			return "phone"
		}
	}

	// Fallback to first available channel
	if pcp.EmailEnabled {
		return "email"
	}
	if pcp.SMSEnabled {
		return "sms"
	}
	if pcp.WhatsAppEnabled {
		return "whatsapp"
	}
	if pcp.PushEnabled {
		return "push"
	}
	if pcp.PhoneEnabled {
		return "phone"
	}

	return ""
}

// IsInPreferredTime checks if current time is in preferred time window
func (pcp *PolicyCommunicationPreference) IsInPreferredTime() bool {
	if pcp.PreferredTime == "" {
		return true // No preference
	}

	hour := time.Now().Hour()

	switch pcp.PreferredTime {
	case "morning":
		return hour >= 6 && hour < 12
	case "afternoon":
		return hour >= 12 && hour < 18
	case "evening":
		return hour >= 18 && hour < 22
	default:
		return true
	}
}

// ShouldSendReminder checks if reminder should be sent based on frequency
func (pcp *PolicyCommunicationPreference) ShouldSendReminder(daysUntilDue int) bool {
	if pcp.ReminderDays <= 0 {
		return false
	}

	// Check if we're within the reminder window
	if daysUntilDue > pcp.ReminderDays {
		return false
	}

	// Check frequency
	switch pcp.ReminderFrequency {
	case "daily":
		return true
	case "weekly":
		return daysUntilDue == pcp.ReminderDays ||
			daysUntilDue == 7 ||
			daysUntilDue == 1
	case "monthly":
		return daysUntilDue == pcp.ReminderDays ||
			daysUntilDue == 1
	default:
		return daysUntilDue == pcp.ReminderDays
	}
}

// OptOut opts out of all communications
func (pcp *PolicyCommunicationPreference) OptOut(reason string) {
	pcp.OptOutAll = true
	now := time.Now()
	pcp.OptOutDate = &now
	pcp.OptOutReason = reason
}

// OptIn opts back into communications
func (pcp *PolicyCommunicationPreference) OptIn() {
	pcp.OptOutAll = false
	pcp.OptOutDate = nil
	pcp.OptOutReason = ""
}
