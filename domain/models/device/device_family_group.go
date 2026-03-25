package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceFamilyGroup manages family plan and multi-device groups
type DeviceFamilyGroup struct {
	database.BaseModel
	GroupID     uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"group_id"`
	GroupName   string    `json:"group_name"`
	CreatedDate time.Time `json:"created_date"`

	// Family Plan Management
	PlanType        string    `json:"plan_type"`     // family, business, enterprise
	PlanStatus      string    `json:"plan_status"`   // active, suspended, cancelled
	BillingCycle    string    `json:"billing_cycle"` // monthly, quarterly, annual
	NextBillingDate time.Time `json:"next_billing_date"`

	// Family Member Devices
	TotalDevices      int    `json:"total_devices"`
	ActiveDevices     int    `json:"active_devices"`
	MaxDevicesAllowed int    `json:"max_devices_allowed"`
	DeviceList        string `gorm:"type:json" json:"device_list"` // JSON array of device IDs

	// Shared Benefits Tracking
	SharedBenefits      string  `gorm:"type:json" json:"shared_benefits"` // JSON array
	BenefitPool         float64 `json:"benefit_pool"`
	BenefitUtilization  float64 `json:"benefit_utilization"`                   // percentage
	BenefitDistribution string  `gorm:"type:json" json:"benefit_distribution"` // JSON object

	// Family Discount Application
	DiscountPercentage float64 `json:"discount_percentage"`
	DiscountAmount     float64 `json:"discount_amount"`
	TotalSavings       float64 `json:"total_savings"`
	DiscountTier       string  `json:"discount_tier"`

	// Primary Account Holder
	PrimaryUserID      uuid.UUID `gorm:"type:uuid" json:"primary_user_id"`
	PrimaryDeviceID    uuid.UUID `gorm:"type:uuid" json:"primary_device_id"`
	AccountPermissions string    `gorm:"type:json" json:"account_permissions"` // JSON object
	PaymentMethodID    uuid.UUID `gorm:"type:uuid" json:"payment_method_id"`

	// Dependent Device Management
	DependentDevices   string `gorm:"type:json" json:"dependent_devices"`   // JSON array
	DependentUsers     string `gorm:"type:json" json:"dependent_users"`     // JSON array
	AccessRestrictions string `gorm:"type:json" json:"access_restrictions"` // JSON object
	AgeRestrictions    string `gorm:"type:json" json:"age_restrictions"`    // JSON object

	// Family Usage Aggregation
	TotalDataUsage     float64 `json:"total_data_usage"` // GB
	TotalMinutesUsed   int     `json:"total_minutes_used"`
	TotalMessagessSent int     `json:"total_messages_sent"`
	SharedDataPool     float64 `json:"shared_data_pool"`                   // GB
	DataDistribution   string  `gorm:"type:json" json:"data_distribution"` // JSON object

	// Family Billing Consolidation
	ConsolidatedInvoice bool    `json:"consolidated_invoice"`
	TotalMonthlyCharge  float64 `json:"total_monthly_charge"`
	IndividualCharges   string  `gorm:"type:json" json:"individual_charges"` // JSON object
	PaymentSplitOption  string  `json:"payment_split_option"`
	AutoPayEnabled      bool    `json:"auto_pay_enabled"`

	// Cross-device Permissions
	DeviceSharing    bool   `json:"device_sharing"`
	ContentSharing   bool   `json:"content_sharing"`
	LocationSharing  bool   `json:"location_sharing"`
	BackupSharing    bool   `json:"backup_sharing"`
	PermissionMatrix string `gorm:"type:json" json:"permission_matrix"` // JSON object

	// Family Safety Features
	ParentalControls   bool   `json:"parental_controls"`
	ScreenTimeTracking bool   `json:"screen_time_tracking"`
	LocationTracking   bool   `json:"location_tracking"`
	EmergencyContacts  string `gorm:"type:json" json:"emergency_contacts"` // JSON array
	SafetyAlerts       string `gorm:"type:json" json:"safety_alerts"`      // JSON array

	// Group Statistics
	GroupAge           int     `json:"group_age"`           // days
	MemberSatisfaction float64 `json:"member_satisfaction"` // 0-100
	RetentionRate      float64 `json:"retention_rate"`      // percentage
	UpgradeRate        float64 `json:"upgrade_rate"`        // percentage

	// Relationships
	// PrimaryUser should be loaded via service layer using PrimaryUserID to avoid circular import
	// PrimaryDevice should be loaded via service layer using PrimaryDeviceID to avoid circular import
	// PaymentMethod should be loaded via service layer using PaymentMethodID to avoid circular import
}

// DeviceGroupDiscounts manages multi-device and group discounts
type DeviceGroupDiscounts struct {
	database.BaseModel
	GroupID      uuid.UUID  `gorm:"type:uuid;index" json:"group_id"`
	DiscountCode string     `gorm:"uniqueIndex" json:"discount_code"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`

	// Multi-device Discount Tracking
	DeviceCount         int     `json:"device_count"`
	DiscountPerDevice   float64 `json:"discount_per_device"`
	TotalDiscountAmount float64 `json:"total_discount_amount"`
	MinimumDevices      int     `json:"minimum_devices"`
	MaximumDevices      int     `json:"maximum_devices"`

	// Bundle Pricing
	BundleType         string  `json:"bundle_type"`
	BundlePrice        float64 `json:"bundle_price"`
	BundleSavings      float64 `json:"bundle_savings"`
	BundleComponents   string  `gorm:"type:json" json:"bundle_components"`   // JSON array
	BundleRestrictions string  `gorm:"type:json" json:"bundle_restrictions"` // JSON array

	// Volume Discounts
	VolumeThresholds   string  `gorm:"type:json" json:"volume_thresholds"` // JSON array
	CurrentVolumeTier  string  `json:"current_volume_tier"`
	NextVolumeTier     string  `json:"next_volume_tier"`
	UnitsToNextTier    int     `json:"units_to_next_tier"`
	VolumeDiscountRate float64 `json:"volume_discount_rate"` // percentage

	// Corporate Discounts
	CorporatePlan      bool    `json:"corporate_plan"`
	CompanyName        string  `json:"company_name"`
	EmployeeCount      int     `json:"employee_count"`
	CorporateRate      float64 `json:"corporate_rate"` // percentage
	CorporateAgreement string  `json:"corporate_agreement"`

	// Student Discounts
	StudentPlan         bool       `json:"student_plan"`
	InstitutionName     string     `json:"institution_name"`
	StudentVerification string     `json:"student_verification"`
	StudentDiscountRate float64    `json:"student_discount_rate"` // percentage
	VerificationExpiry  *time.Time `json:"verification_expiry"`

	// Senior Discounts
	SeniorPlan         bool    `json:"senior_plan"`
	MinimumAge         int     `json:"minimum_age"`
	SeniorDiscountRate float64 `json:"senior_discount_rate"` // percentage
	AgeVerified        bool    `json:"age_verified"`

	// Military Discounts
	MilitaryPlan         bool    `json:"military_plan"`
	ServiceBranch        string  `json:"service_branch"`
	MilitaryStatus       string  `json:"military_status"`        // active, veteran, family
	MilitaryDiscountRate float64 `json:"military_discount_rate"` // percentage
	MilitaryVerification string  `json:"military_verification"`

	// Loyalty Discounts
	LoyaltyTier         string  `json:"loyalty_tier"`
	YearsAsCustomer     int     `json:"years_as_customer"`
	LoyaltyDiscountRate float64 `json:"loyalty_discount_rate"` // percentage
	LoyaltyPoints       int     `json:"loyalty_points"`
	PointsMultiplier    float64 `json:"points_multiplier"`

	// Promotional Discounts
	PromoCode         string  `json:"promo_code"`
	PromoType         string  `json:"promo_type"`
	PromoDiscountRate float64 `json:"promo_discount_rate"` // percentage
	PromoMaxUsage     int     `json:"promo_max_usage"`
	PromoUsageCount   int     `json:"promo_usage_count"`

	// Discount Eligibility Tracking
	EligibilityStatus    string     `json:"eligibility_status"`                    // eligible, pending, expired
	EligibilityCriteria  string     `gorm:"type:json" json:"eligibility_criteria"` // JSON object
	VerificationRequired bool       `json:"verification_required"`
	LastVerification     *time.Time `json:"last_verification"`
	NextReviewDate       *time.Time `json:"next_review_date"`

	// Discount Stacking
	StackingAllowed      bool    `json:"stacking_allowed"`
	StackedDiscounts     string  `gorm:"type:json" json:"stacked_discounts"` // JSON array
	MaxStackCount        int     `json:"max_stack_count"`
	TotalStackedDiscount float64 `json:"total_stacked_discount"`

	// Discount Impact
	TotalSaved      float64 `json:"total_saved"`
	MonthlyImpact   float64 `json:"monthly_impact"`
	AnnualImpact    float64 `json:"annual_impact"`
	CustomerValue   float64 `json:"customer_value"`
	RetentionImpact float64 `json:"retention_impact"` // percentage

	// Relationships
	FamilyGroup *DeviceFamilyGroup `gorm:"foreignKey:GroupID" json:"family_group,omitempty"`
}

// Methods for DeviceSatisfactionScore
func (dss *DeviceSatisfactionScore) IsHighlySatisfied() bool {
	return dss.OverallSatisfaction >= 80 && dss.NPSCategory == "promoter"
}

func (dss *DeviceSatisfactionScore) IsAtRisk() bool {
	return dss.ChurnRisk > 70 || dss.NPSCategory == "detractor" ||
		dss.OverallSatisfaction < 50
}

func (dss *DeviceSatisfactionScore) GetSatisfactionGrade() string {
	if dss.OverallSatisfaction >= 90 {
		return "Excellent"
	} else if dss.OverallSatisfaction >= 75 {
		return "Good"
	} else if dss.OverallSatisfaction >= 60 {
		return "Fair"
	} else if dss.OverallSatisfaction >= 40 {
		return "Poor"
	}
	return "Very Poor"
}

func (dss *DeviceSatisfactionScore) NeedsAttention() bool {
	return dss.SatisfactionTrend == "decreasing" || dss.ChurnRisk > 50 ||
		dss.SupportSatisfaction < 60
}

func (dss *DeviceSatisfactionScore) IsPromoter() bool {
	return dss.NPSCategory == "promoter" && dss.LikelihoodToRecommend >= 9
}

// Methods for DeviceRecommendations
func (dr *DeviceRecommendations) HasHighSavingsPotential() bool {
	return dr.MonthlySavingsPotential > 20 || dr.OptimizationPotential > 30
}

func (dr *DeviceRecommendations) IsHighlyPersonalized() bool {
	return dr.PersonalizationScore >= 80 && dr.AcceptanceRate > 70
}

func (dr *DeviceRecommendations) HasSecurityConcerns() bool {
	return dr.SecurityScore < 60 // Security score below 60 indicates concerns
}

func (dr *DeviceRecommendations) NeedsUpgrade() bool {
	return dr.TradeInValue > 0 && dr.DeviceUpgradeOptions != ""
}

func (dr *DeviceRecommendations) GetRecommendationPriority() string {
	if dr.SecurityScore < 50 {
		return "Critical"
	} else if dr.MonthlySavingsPotential > 50 {
		return "High"
	} else if dr.FeatureRelevanceScore > 70 {
		return "Medium"
	}
	return "Low"
}

// Methods for DeviceLoyaltyProgram
func (dlp *DeviceLoyaltyProgram) IsTopTier() bool {
	return dlp.CurrentTier == "platinum" || dlp.CurrentTier == "gold"
}

func (dlp *DeviceLoyaltyProgram) HasExpiringSoon() bool {
	if dlp.ExpirationDate == nil {
		return false
	}
	return time.Until(*dlp.ExpirationDate) < 30*24*time.Hour // Less than 30 days
}

func (dlp *DeviceLoyaltyProgram) IsHighlyEngaged() bool {
	return dlp.EngagementScore >= 70 && dlp.BenefitUtilization >= 60
}

func (dlp *DeviceLoyaltyProgram) GetLoyaltyLevel() string {
	if dlp.LifetimePoints >= 100000 {
		return "Elite"
	} else if dlp.LifetimePoints >= 50000 {
		return "Premium"
	} else if dlp.LifetimePoints >= 20000 {
		return "Advanced"
	} else if dlp.LifetimePoints >= 5000 {
		return "Member"
	}
	return "Basic"
}

func (dlp *DeviceLoyaltyProgram) IsProfitable() bool {
	return dlp.ProgramROI > 0 && dlp.CustomerValue > dlp.SavedAmount
}

// Methods for DeviceFeedback
func (df *DeviceFeedback) IsPositiveFeedback() bool {
	return df.SentimentScore > 0.3 || df.Rating >= 4 ||
		df.FeedbackType == "compliment"
}

func (df *DeviceFeedback) RequiresUrgentResponse() bool {
	return df.ComplaintSeverity == "high" || df.BugSeverity == "critical" ||
		(df.FeedbackType == "complaint" && !df.ResponseProvided)
}

func (df *DeviceFeedback) HasActionableInsights() bool {
	return df.FeatureStatus == "submitted" || df.SuggestionImpact == "high" ||
		df.BugStatus == "confirmed"
}

func (df *DeviceFeedback) GetResponsePriority() string {
	if df.BugSeverity == "critical" || df.ComplaintSeverity == "high" {
		return "Urgent"
	} else if df.FeaturePriority == "high" || df.SuggestionImpact == "high" {
		return "High"
	} else if df.Rating < 3 {
		return "Medium"
	}
	return "Low"
}

func (df *DeviceFeedback) IsResolved() bool {
	return (df.ComplaintStatus == "resolved" && df.CustomerSatisfied != nil && *df.CustomerSatisfied) ||
		df.BugStatus == "resolved" || df.FeatureStatus == "implemented"
}

// Methods for DeviceUserJourney
func (duj *DeviceUserJourney) IsHealthyJourney() bool {
	return duj.JourneyScore >= 70 && duj.AbandonmentRisk < 30 &&
		duj.PositiveMoments > duj.NegativeMoments
}

func (duj *DeviceUserJourney) IsChurnRisk() bool {
	return duj.JourneyPhase == "at-risk" || duj.AbandonmentRisk > 70 ||
		duj.CurrentStage == "churned"
}

func (duj *DeviceUserJourney) HasConversionPotential() bool {
	return duj.ConversionProbability > 60 && duj.FunnelStage != ""
}

func (duj *DeviceUserJourney) GetJourneyHealth() string {
	if duj.JourneyScore >= 80 {
		return "Excellent"
	} else if duj.JourneyScore >= 60 {
		return "Good"
	} else if duj.JourneyScore >= 40 {
		return "Fair"
	}
	return "Poor"
}

func (duj *DeviceUserJourney) NeedsIntervention() bool {
	return duj.AbandonmentRisk > 50 || duj.NegativeMoments > 3 ||
		duj.StageHealthScore < 50
}

// Methods for DeviceFamilyGroup
func (dfg *DeviceFamilyGroup) IsFullCapacity() bool {
	return dfg.ActiveDevices >= dfg.MaxDevicesAllowed
}

func (dfg *DeviceFamilyGroup) HasHighUtilization() bool {
	return dfg.BenefitUtilization > 75 && dfg.DataDistribution != ""
}

func (dfg *DeviceFamilyGroup) IsEligibleForUpgrade() bool {
	return dfg.ActiveDevices >= dfg.MaxDevicesAllowed*80/100 && dfg.RetentionRate > 80
}

func (dfg *DeviceFamilyGroup) GetPlanValue() string {
	savingsPerDevice := dfg.TotalSavings / float64(dfg.ActiveDevices)
	if savingsPerDevice > 50 {
		return "Excellent"
	} else if savingsPerDevice > 30 {
		return "Good"
	} else if savingsPerDevice > 10 {
		return "Fair"
	}
	return "Basic"
}

func (dfg *DeviceFamilyGroup) HasActiveSafety() bool {
	return dfg.ParentalControls || dfg.LocationTracking || dfg.ScreenTimeTracking
}

// Methods for DeviceGroupDiscounts
func (dgd *DeviceGroupDiscounts) IsActive() bool {
	return dgd.EligibilityStatus == "eligible" &&
		(dgd.EndDate == nil || time.Now().Before(*dgd.EndDate))
}

func (dgd *DeviceGroupDiscounts) GetTotalDiscount() float64 {
	if dgd.StackingAllowed {
		return dgd.TotalStackedDiscount
	}
	return dgd.TotalDiscountAmount
}

func (dgd *DeviceGroupDiscounts) IsMaximized() bool {
	return dgd.DeviceCount >= dgd.MaximumDevices ||
		(dgd.StackingAllowed && dgd.StackedDiscounts != "" &&
			len(dgd.StackedDiscounts) >= dgd.MaxStackCount)
}

func (dgd *DeviceGroupDiscounts) RequiresVerification() bool {
	return dgd.VerificationRequired &&
		(dgd.LastVerification == nil ||
			time.Since(*dgd.LastVerification) > 365*24*time.Hour)
}

func (dgd *DeviceGroupDiscounts) GetDiscountType() string {
	if dgd.CorporatePlan {
		return "Corporate"
	} else if dgd.StudentPlan {
		return "Student"
	} else if dgd.SeniorPlan {
		return "Senior"
	} else if dgd.MilitaryPlan {
		return "Military"
	} else if dgd.DeviceCount > 1 {
		return "Multi-Device"
	}
	return "Individual"
}
