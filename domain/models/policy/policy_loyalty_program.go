package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyLoyaltyProgram represents loyalty and rewards program for a policy
type PolicyLoyaltyProgram struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`

	// Points System
	LoyaltyPoints  int        `json:"loyalty_points"`
	LifetimePoints int        `json:"lifetime_points"`
	PointsExpiry   *time.Time `json:"points_expiry"`
	PendingPoints  int        `json:"pending_points"`

	// Tier System
	LoyaltyTier      string    `gorm:"type:varchar(20)" json:"loyalty_tier"` // bronze, silver, gold, platinum, diamond
	TierStartDate    time.Time `json:"tier_start_date"`
	TierExpiryDate   time.Time `json:"tier_expiry_date"`
	NextTierProgress float64   `json:"next_tier_progress"` // Percentage to next tier
	TierBenefitValue float64   `json:"tier_benefit_value"` // Current tier discount percentage

	// Rewards
	RewardPoints      int     `json:"reward_points"`
	CashbackEarned    float64 `json:"cashback_earned"`
	CashbackAvailable float64 `json:"cashback_available"`
	CashbackRedeemed  float64 `json:"cashback_redeemed"`
	VoucherCredits    float64 `json:"voucher_credits"`

	// Referral Program
	ReferralCode         string  `gorm:"uniqueIndex" json:"referral_code"`
	ReferralCount        int     `json:"referral_count"`
	SuccessfulReferrals  int     `json:"successful_referrals"`
	ReferralBonus        float64 `json:"referral_bonus"`
	ReferralBonusPending float64 `json:"referral_bonus_pending"`

	// Activity Tracking
	LastActivityDate  *time.Time `json:"last_activity_date"`
	ConsecutiveMonths int        `json:"consecutive_months"`                  // Months with activity
	MilestoneAchieved string     `gorm:"type:json" json:"milestone_achieved"` // JSON array of milestones
	BadgesEarned      string     `gorm:"type:json" json:"badges_earned"`      // JSON array of badges

	// Engagement Metrics
	EngagementScore   float64 `json:"engagement_score"`
	AppUsageFrequency int     `json:"app_usage_frequency"` // Times per month
	ReviewsSubmitted  int     `json:"reviews_submitted"`
	SurveysCompleted  int     `json:"surveys_completed"`

	// Special Programs
	VIPMember    bool `gorm:"default:false" json:"vip_member"`
	EarlyAdopter bool `gorm:"default:false" json:"early_adopter"`
	BetaTester   bool `gorm:"default:false" json:"beta_tester"`

	// Status
	IsActive         bool   `gorm:"default:true" json:"is_active"`
	IsSuspended      bool   `gorm:"default:false" json:"is_suspended"`
	SuspensionReason string `json:"suspension_reason"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
	// Customer accessed via Policy.Customer relationship - CustomerID field removed to avoid duplication
}

// TableName returns the table name
func (PolicyLoyaltyProgram) TableName() string {
	return "policy_loyalty_programs"
}

// GetTierBenefits returns benefits based on tier
func (plp *PolicyLoyaltyProgram) GetTierBenefits() float64 {
	switch plp.LoyaltyTier {
	case "diamond":
		return 0.25 // 25% discount
	case "platinum":
		return 0.20 // 20% discount
	case "gold":
		return 0.15 // 15% discount
	case "silver":
		return 0.10 // 10% discount
	case "bronze":
		return 0.05 // 5% discount
	default:
		return 0
	}
}

// CalculatePointsForActivity calculates points for an activity
func (plp *PolicyLoyaltyProgram) CalculatePointsForActivity(activity string, value float64) int {
	basePoints := 0

	switch activity {
	case "premium_payment":
		basePoints = int(value * 0.1) // 1 point per $10
	case "claim_free_month":
		basePoints = 50
	case "referral":
		basePoints = 100
	case "review":
		basePoints = 25
	case "survey":
		basePoints = 15
	case "app_login":
		basePoints = 5
	case "policy_renewal":
		basePoints = 200
	default:
		basePoints = 10
	}

	// Apply tier multiplier
	multiplier := 1.0
	switch plp.LoyaltyTier {
	case "diamond":
		multiplier = 2.0
	case "platinum":
		multiplier = 1.75
	case "gold":
		multiplier = 1.5
	case "silver":
		multiplier = 1.25
	}

	return int(float64(basePoints) * multiplier)
}

// AddPoints adds loyalty points
func (plp *PolicyLoyaltyProgram) AddPoints(points int, pending bool) {
	if pending {
		plp.PendingPoints += points
	} else {
		plp.LoyaltyPoints += points
		plp.LifetimePoints += points
		plp.UpdateTier()
	}

	now := time.Now()
	plp.LastActivityDate = &now
}

// UpdateTier updates the loyalty tier based on lifetime points
func (plp *PolicyLoyaltyProgram) UpdateTier() {
	oldTier := plp.LoyaltyTier

	switch {
	case plp.LifetimePoints >= 10000:
		plp.LoyaltyTier = "diamond"
		plp.NextTierProgress = 100
	case plp.LifetimePoints >= 5000:
		plp.LoyaltyTier = "platinum"
		plp.NextTierProgress = float64(plp.LifetimePoints-5000) / 5000 * 100
	case plp.LifetimePoints >= 2500:
		plp.LoyaltyTier = "gold"
		plp.NextTierProgress = float64(plp.LifetimePoints-2500) / 2500 * 100
	case plp.LifetimePoints >= 1000:
		plp.LoyaltyTier = "silver"
		plp.NextTierProgress = float64(plp.LifetimePoints-1000) / 1500 * 100
	case plp.LifetimePoints >= 500:
		plp.LoyaltyTier = "bronze"
		plp.NextTierProgress = float64(plp.LifetimePoints-500) / 500 * 100
	default:
		plp.LoyaltyTier = "basic"
		plp.NextTierProgress = float64(plp.LifetimePoints) / 500 * 100
	}

	// Update tier dates if tier changed
	if oldTier != plp.LoyaltyTier {
		plp.TierStartDate = time.Now()
		plp.TierExpiryDate = time.Now().AddDate(1, 0, 0) // 1 year
		plp.TierBenefitValue = plp.GetTierBenefits()
	}
}

// RedeemPoints redeems points for rewards
func (plp *PolicyLoyaltyProgram) RedeemPoints(points int) bool {
	if plp.LoyaltyPoints < points {
		return false
	}

	plp.LoyaltyPoints -= points
	// Convert to cashback (1 point = $0.01)
	cashback := float64(points) * 0.01
	plp.CashbackAvailable += cashback

	return true
}

// HasEarnedReferralBonus checks if referral bonus is earned
func (plp *PolicyLoyaltyProgram) HasEarnedReferralBonus() bool {
	return plp.SuccessfulReferrals > 0 && plp.ReferralBonus > 0
}

// ProcessReferral processes a successful referral
func (plp *PolicyLoyaltyProgram) ProcessReferral() {
	plp.ReferralCount++
	plp.SuccessfulReferrals++
	bonus := 50.0 // $50 per successful referral
	plp.ReferralBonusPending += bonus

	// Add loyalty points
	plp.AddPoints(100, false)
}

// IsEligibleForVIP checks VIP eligibility
func (plp *PolicyLoyaltyProgram) IsEligibleForVIP() bool {
	return plp.LoyaltyTier == "diamond" ||
		plp.LoyaltyTier == "platinum" ||
		plp.LifetimePoints >= 7500
}

// GetEngagementLevel returns engagement level
func (plp *PolicyLoyaltyProgram) GetEngagementLevel() string {
	if plp.EngagementScore >= 80 {
		return "highly_engaged"
	} else if plp.EngagementScore >= 60 {
		return "engaged"
	} else if plp.EngagementScore >= 40 {
		return "moderate"
	} else if plp.EngagementScore >= 20 {
		return "low"
	}
	return "inactive"
}
