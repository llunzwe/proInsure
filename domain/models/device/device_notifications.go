package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// DeviceNotificationManagement represents notification and communication management
type DeviceNotificationManagement struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Push Notifications
	PushEnabled      bool           `gorm:"type:boolean;default:true" json:"push_enabled"`
	PushTokens       datatypes.JSON `gorm:"type:json" json:"push_tokens"`      // []PushToken
	PushPreferences  datatypes.JSON `gorm:"type:json" json:"push_preferences"` // map[string]bool
	PushHistory      datatypes.JSON `gorm:"type:json" json:"push_history"`     // []Notification
	LastPushDate     *time.Time     `gorm:"type:timestamp" json:"last_push_date,omitempty"`
	PushDeliveryRate float64        `gorm:"type:decimal(5,2)" json:"push_delivery_rate"`
	PushClickRate    float64        `gorm:"type:decimal(5,2)" json:"push_click_rate"`

	// Email Campaigns
	EmailEnabled      bool           `gorm:"type:boolean;default:true" json:"email_enabled"`
	EmailPreferences  datatypes.JSON `gorm:"type:json" json:"email_preferences"` // map[string]bool
	EmailCampaigns    datatypes.JSON `gorm:"type:json" json:"email_campaigns"`   // []Campaign
	EmailOpenRate     float64        `gorm:"type:decimal(5,2)" json:"email_open_rate"`
	EmailClickRate    float64        `gorm:"type:decimal(5,2)" json:"email_click_rate"`
	UnsubscribeStatus bool           `gorm:"type:boolean;default:false" json:"unsubscribe_status"`
	LastEmailDate     *time.Time     `gorm:"type:timestamp" json:"last_email_date,omitempty"`

	// SMS Alerts
	SMSEnabled      bool           `gorm:"type:boolean;default:true" json:"sms_enabled"`
	SMSNumber       string         `gorm:"type:varchar(20)" json:"sms_number"`
	SMSPreferences  datatypes.JSON `gorm:"type:json" json:"sms_preferences"` // map[string]bool
	SMSHistory      datatypes.JSON `gorm:"type:json" json:"sms_history"`     // []SMS
	SMSDeliveryRate float64        `gorm:"type:decimal(5,2)" json:"sms_delivery_rate"`
	LastSMSDate     *time.Time     `gorm:"type:timestamp" json:"last_sms_date,omitempty"`
	SMSOptOut       bool           `gorm:"type:boolean;default:false" json:"sms_opt_out"`

	// In-App Messages
	InAppEnabled    bool           `gorm:"type:boolean;default:true" json:"in_app_enabled"`
	InAppMessages   datatypes.JSON `gorm:"type:json" json:"in_app_messages"` // []Message
	UnreadMessages  int            `gorm:"type:int" json:"unread_messages"`
	InAppEngagement float64        `gorm:"type:decimal(5,2)" json:"in_app_engagement"`
	MessagePriority datatypes.JSON `gorm:"type:json" json:"message_priority"` // map[string]int

	// Communication Effectiveness
	PreferredChannel     string         `gorm:"type:varchar(50)" json:"preferred_channel"` // push, email, sms, in_app
	ChannelEffectiveness datatypes.JSON `gorm:"type:json" json:"channel_effectiveness"`    // map[string]float64
	OptimalSendTime      string         `gorm:"type:varchar(50)" json:"optimal_send_time"`
	ResponseRate         float64        `gorm:"type:decimal(5,2)" json:"response_rate"`
	EngagementScore      float64        `gorm:"type:decimal(5,2)" json:"engagement_score"`

	// Notification Schedule
	DoNotDisturbStart *time.Time     `gorm:"type:time" json:"do_not_disturb_start,omitempty"`
	DoNotDisturbEnd   *time.Time     `gorm:"type:time" json:"do_not_disturb_end,omitempty"`
	TimeZone          string         `gorm:"type:varchar(50)" json:"time_zone"`
	PreferredDays     datatypes.JSON `gorm:"type:json" json:"preferred_days"` // []string
	FrequencyCap      int            `gorm:"type:int" json:"frequency_cap_per_day"`

	// Campaign Tracking
	CampaignSubscriptions datatypes.JSON `gorm:"type:json" json:"campaign_subscriptions"` // []Campaign
	CampaignHistory       datatypes.JSON `gorm:"type:json" json:"campaign_history"`       // []CampaignResult
	ABTestParticipation   datatypes.JSON `gorm:"type:json" json:"ab_test_participation"`  // []ABTest
	ConversionTracking    datatypes.JSON `gorm:"type:json" json:"conversion_tracking"`    // []Conversion

	// Status
	NotificationStatus string    `gorm:"type:varchar(50)" json:"notification_status"`
	LastContactDate    time.Time `gorm:"type:timestamp" json:"last_contact_date"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceCalendarScheduling represents calendar and scheduling features
type DeviceCalendarScheduling struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Service Appointments
	ServiceAppointments  datatypes.JSON `gorm:"type:json" json:"service_appointments"` // []Appointment
	UpcomingAppointments int            `gorm:"type:int" json:"upcoming_appointments"`
	NextAppointment      *time.Time     `gorm:"type:timestamp" json:"next_appointment,omitempty"`
	AppointmentReminders datatypes.JSON `gorm:"type:json" json:"appointment_reminders"` // []Reminder
	PreferredTimeSlots   datatypes.JSON `gorm:"type:json" json:"preferred_time_slots"`  // []TimeSlot

	// Inspection Schedules
	InspectionSchedules  datatypes.JSON `gorm:"type:json" json:"inspection_schedules"` // []Schedule
	NextInspectionDate   *time.Time     `gorm:"type:timestamp" json:"next_inspection_date,omitempty"`
	InspectionFrequency  string         `gorm:"type:varchar(50)" json:"inspection_frequency"` // monthly, quarterly, annual
	MissedInspections    int            `gorm:"type:int" json:"missed_inspections"`
	InspectionCompliance float64        `gorm:"type:decimal(5,2)" json:"inspection_compliance"`

	// Warranty Reminders
	WarrantyExpiration   *time.Time     `gorm:"type:timestamp" json:"warranty_expiration,omitempty"`
	WarrantyReminders    datatypes.JSON `gorm:"type:json" json:"warranty_reminders"` // []Reminder
	ExtendedWarrantyDate *time.Time     `gorm:"type:timestamp" json:"extended_warranty_date,omitempty"`
	WarrantyAlertDays    int            `gorm:"type:int" json:"warranty_alert_days"`

	// Payment Schedules
	PaymentDueDates  datatypes.JSON `gorm:"type:json" json:"payment_due_dates"` // []DueDate
	NextPaymentDate  *time.Time     `gorm:"type:timestamp" json:"next_payment_date,omitempty"`
	PaymentReminders datatypes.JSON `gorm:"type:json" json:"payment_reminders"` // []Reminder
	AutoPaySchedule  datatypes.JSON `gorm:"type:json" json:"auto_pay_schedule"`
	PaymentHistory   datatypes.JSON `gorm:"type:json" json:"payment_history"` // []Payment

	// Coverage Renewals
	CoverageRenewalDate *time.Time     `gorm:"type:timestamp" json:"coverage_renewal_date,omitempty"`
	RenewalReminders    datatypes.JSON `gorm:"type:json" json:"renewal_reminders"` // []Reminder
	RenewalWindow       datatypes.JSON `gorm:"type:json" json:"renewal_window"`    // DateRange
	GracePeriodEnd      *time.Time     `gorm:"type:timestamp" json:"grace_period_end,omitempty"`

	// Maintenance Calendar
	MaintenanceDates      datatypes.JSON `gorm:"type:json" json:"maintenance_dates"`      // []MaintenanceDate
	PreventiveMaintenance datatypes.JSON `gorm:"type:json" json:"preventive_maintenance"` // []Schedule
	LastMaintenanceDate   *time.Time     `gorm:"type:timestamp" json:"last_maintenance_date,omitempty"`
	MaintenanceCompliance float64        `gorm:"type:decimal(5,2)" json:"maintenance_compliance"`

	// Blackout Dates
	BlackoutDates        datatypes.JSON `gorm:"type:json" json:"blackout_dates"`   // []DateRange
	HolidaySchedule      datatypes.JSON `gorm:"type:json" json:"holiday_schedule"` // []Holiday
	AvailabilityCalendar datatypes.JSON `gorm:"type:json" json:"availability_calendar"`

	// Recurring Events
	RecurringEvents    datatypes.JSON `gorm:"type:json" json:"recurring_events"` // []RecurringEvent
	EventFrequency     datatypes.JSON `gorm:"type:json" json:"event_frequency"`  // map[string]string
	NextRecurringEvent *time.Time     `gorm:"type:timestamp" json:"next_recurring_event,omitempty"`

	// Status
	SchedulingStatus string    `gorm:"type:varchar(50)" json:"scheduling_status"`
	LastUpdatedDate  time.Time `gorm:"type:timestamp" json:"last_updated_date"`
	CreatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceComparison represents comparison and benchmarking features
type DeviceComparison struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Device Comparison Tools
	ComparisonHistory  datatypes.JSON `gorm:"type:json" json:"comparison_history"`  // []Comparison
	ComparedDevices    datatypes.JSON `gorm:"type:json" json:"compared_devices"`    // []Device
	ComparisonCriteria datatypes.JSON `gorm:"type:json" json:"comparison_criteria"` // []Criteria
	ComparisonScore    float64        `gorm:"type:decimal(5,2)" json:"comparison_score"`
	BetterAlternatives datatypes.JSON `gorm:"type:json" json:"better_alternatives"` // []Alternative

	// Market Value Benchmarking
	MarketValue       float64        `gorm:"type:decimal(15,2)" json:"market_value"`
	MarketPosition    string         `gorm:"type:varchar(50)" json:"market_position"` // below, average, above
	ValuePercentile   float64        `gorm:"type:decimal(5,2)" json:"value_percentile"`
	PriceTrend        string         `gorm:"type:varchar(50)" json:"price_trend"` // increasing, stable, decreasing
	MarketComparisons datatypes.JSON `gorm:"type:json" json:"market_comparisons"` // []MarketComp

	// Feature Comparison
	FeatureMatrix     datatypes.JSON `gorm:"type:json" json:"feature_matrix"` // map[string]map[string]interface{}
	FeatureScore      float64        `gorm:"type:decimal(5,2)" json:"feature_score"`
	MissingFeatures   datatypes.JSON `gorm:"type:json" json:"missing_features"`   // []string
	UniqueFeatures    datatypes.JSON `gorm:"type:json" json:"unique_features"`    // []string
	FeatureComparison datatypes.JSON `gorm:"type:json" json:"feature_comparison"` // []FeatureComp

	// Performance Benchmarking
	PerformanceScore    float64        `gorm:"type:decimal(5,2)" json:"performance_score"`
	BenchmarkResults    datatypes.JSON `gorm:"type:json" json:"benchmark_results"` // map[string]float64
	PerformanceRanking  int            `gorm:"type:int" json:"performance_ranking"`
	PerformanceCategory string         `gorm:"type:varchar(50)" json:"performance_category"`
	BenchmarkHistory    datatypes.JSON `gorm:"type:json" json:"benchmark_history"` // []Benchmark

	// Cost-Benefit Analysis
	TotalCostOwnership float64 `gorm:"type:decimal(15,2)" json:"total_cost_ownership"`
	BenefitScore       float64 `gorm:"type:decimal(5,2)" json:"benefit_score"`
	CostBenefitRatio   float64 `gorm:"type:decimal(10,2)" json:"cost_benefit_ratio"`
	ROICalculation     float64 `gorm:"type:decimal(10,2)" json:"roi_calculation"`
	PaybackPeriod      int     `gorm:"type:int" json:"payback_period_months"`

	// Insurance Comparison
	CoverageComparison   datatypes.JSON `gorm:"type:json" json:"coverage_comparison"` // []CoverageComp
	PremiumComparison    datatypes.JSON `gorm:"type:json" json:"premium_comparison"`  // []PremiumComp
	DeductibleComparison datatypes.JSON `gorm:"type:json" json:"deductible_comparison"`
	BestInsuranceOption  datatypes.JSON `gorm:"type:json" json:"best_insurance_option"`
	InsuranceSavings     float64        `gorm:"type:decimal(15,2)" json:"insurance_savings"`

	// Competitive Analysis
	CompetitorDevices    datatypes.JSON `gorm:"type:json" json:"competitor_devices"` // []Device
	CompetitivePosition  string         `gorm:"type:varchar(50)" json:"competitive_position"`
	MarketShare          float64        `gorm:"type:decimal(5,2)" json:"market_share"`
	CompetitiveAdvantage datatypes.JSON `gorm:"type:json" json:"competitive_advantage"` // []Advantage

	// Recommendations
	UpgradeRecommendations datatypes.JSON `gorm:"type:json" json:"upgrade_recommendations"` // []Recommendation
	AlternativeOptions     datatypes.JSON `gorm:"type:json" json:"alternative_options"`     // []Option
	OptimizationTips       datatypes.JSON `gorm:"type:json" json:"optimization_tips"`       // []Tip

	// Status
	ComparisonStatus   string    `gorm:"type:varchar(50)" json:"comparison_status"`
	LastComparisonDate time.Time `gorm:"type:timestamp" json:"last_comparison_date"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsEngaged checks if user is engaged with notifications
func (dnm *DeviceNotificationManagement) IsEngaged() bool {
	return dnm.EngagementScore > 50 && dnm.ResponseRate > 30
}

// HasOptedOut checks if user has opted out
func (dnm *DeviceNotificationManagement) HasOptedOut() bool {
	return dnm.UnsubscribeStatus || dnm.SMSOptOut ||
		(!dnm.PushEnabled && !dnm.EmailEnabled && !dnm.SMSEnabled)
}

// GetBestChannel returns the best communication channel
func (dnm *DeviceNotificationManagement) GetBestChannel() string {
	if dnm.PreferredChannel != "" {
		return dnm.PreferredChannel
	}
	// Default to push if enabled
	if dnm.PushEnabled {
		return "push"
	}
	return "email"
}

// HasUpcomingEvents checks for upcoming scheduled events
func (dcs *DeviceCalendarScheduling) HasUpcomingEvents() bool {
	return dcs.UpcomingAppointments > 0 ||
		dcs.NextAppointment != nil ||
		dcs.NextPaymentDate != nil
}

// IsCompliant checks scheduling compliance
func (dcs *DeviceCalendarScheduling) IsCompliant() bool {
	return dcs.InspectionCompliance >= 80 &&
		dcs.MaintenanceCompliance >= 80 &&
		dcs.MissedInspections == 0
}

// NeedsRenewal checks if coverage needs renewal
func (dcs *DeviceCalendarScheduling) NeedsRenewal() bool {
	if dcs.CoverageRenewalDate == nil {
		return false
	}
	daysUntilRenewal := time.Until(*dcs.CoverageRenewalDate).Hours() / 24
	return daysUntilRenewal <= 30
}

// IsBetterValue checks if device is better value than market
func (dc *DeviceComparison) IsBetterValue() bool {
	return dc.MarketPosition == "above" && dc.ValuePercentile > 70
}

// HasCompetitiveAdvantage checks for competitive advantage
func (dc *DeviceComparison) HasCompetitiveAdvantage() bool {
	return dc.CompetitivePosition == "leader" ||
		dc.UniqueFeatures != nil ||
		dc.FeatureScore > 80
}

// ShouldUpgrade checks if upgrade is recommended
func (dc *DeviceComparison) ShouldUpgrade() bool {
	return dc.UpgradeRecommendations != nil &&
		dc.BetterAlternatives != nil &&
		dc.CostBenefitRatio > 1.5
}
