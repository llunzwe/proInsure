package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyEnvironmental represents environmental and sustainability features for a policy
type PolicyEnvironmental struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`
	// Removed CustomerID - access via Policy.CustomerID to avoid duplication

	// Green Policy Features
	GreenPolicy         bool       `gorm:"default:false" json:"green_policy"`
	GreenPolicyTier     string     `gorm:"type:varchar(20)" json:"green_policy_tier"` // basic, silver, gold, platinum
	CertificationNumber string     `json:"certification_number"`
	CertifiedBy         string     `json:"certified_by"`
	CertificationDate   *time.Time `json:"certification_date"`

	// Paperless Options
	PaperlessDiscount bool       `gorm:"default:false" json:"paperless_discount"`
	PaperlessEnrolled bool       `gorm:"default:false" json:"paperless_enrolled"`
	EnrollmentDate    *time.Time `json:"enrollment_date"`
	DocumentsSaved    int        `json:"documents_saved"` // Number of paper documents saved
	TreesSaved        float64    `json:"trees_saved"`     // Estimated trees saved

	// Recycling Program
	RecyclingReward      bool       `gorm:"default:false" json:"recycling_reward"`
	RecyclingPartner     string     `json:"recycling_partner"`
	DevicesRecycled      int        `json:"devices_recycled"`
	RecyclingCredits     float64    `json:"recycling_credits"`
	LastRecyclingDate    *time.Time `json:"last_recycling_date"`
	RecyclingCertificate string     `json:"recycling_certificate"`

	// Carbon Offset
	CarbonOffset        bool    `gorm:"default:false" json:"carbon_offset"`
	CarbonCredits       float64 `json:"carbon_credits"`
	OffsetProvider      string  `json:"offset_provider"`
	OffsetProjectType   string  `json:"offset_project_type"` // renewable, reforestation, etc.
	OffsetCertification string  `json:"offset_certification"`

	// Sustainability Metrics
	SustainabilityScore float64 `json:"sustainability_score"`
	CarbonFootprint     float64 `json:"carbon_footprint"`  // kg CO2
	EWasteReduction     float64 `json:"ewaste_reduction"`  // kg
	EnergyEfficiency    float64 `json:"energy_efficiency"` // Percentage

	// Environmental Impact
	PlasticReduction  float64 `json:"plastic_reduction"`  // kg
	WaterConservation float64 `json:"water_conservation"` // liters
	RenewableEnergy   float64 `json:"renewable_energy"`   // Percentage of operations

	// Eco-Friendly Services
	EcoRepairNetwork       bool `gorm:"default:false" json:"eco_repair_network"`
	GreenShipping          bool `gorm:"default:false" json:"green_shipping"`
	BiodegradablePackaging bool `gorm:"default:false" json:"biodegradable_packaging"`
	LocalSourcing          bool `gorm:"default:false" json:"local_sourcing"`

	// Rewards & Incentives
	EcoPoints    int     `json:"eco_points"`
	EcoRewards   float64 `json:"eco_rewards"` // Monetary value
	TreesPlanted int     `json:"trees_planted"`
	GreenBadges  string  `gorm:"type:json" json:"green_badges"` // JSON array of badges

	// Participation Programs
	ClimateActionMember bool   `gorm:"default:false" json:"climate_action_member"`
	CircularEconomy     bool   `gorm:"default:false" json:"circular_economy"`
	ZeroWasteProgram    bool   `gorm:"default:false" json:"zero_waste_program"`
	SustainableGoals    string `gorm:"type:json" json:"sustainable_goals"` // JSON array of SDG goals

	// Tracking
	LastAssessmentDate  *time.Time `json:"last_assessment_date"`
	NextAssessmentDate  *time.Time `json:"next_assessment_date"`
	ImprovementTarget   float64    `json:"improvement_target"`   // Percentage
	ImprovementProgress float64    `json:"improvement_progress"` // Percentage

	// Status
	IsActive           bool   `gorm:"default:true" json:"is_active"`
	VerificationStatus string `gorm:"type:varchar(20)" json:"verification_status"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
	// Customer accessed via Policy.Customer relationship
}

// TableName returns the table name
func (PolicyEnvironmental) TableName() string {
	return "policy_environmentals"
}

// IsGreenPolicy checks if this is an eco-friendly policy
func (pe *PolicyEnvironmental) IsGreenPolicy() bool {
	return pe.IsActive && (pe.GreenPolicy ||
		pe.PaperlessDiscount ||
		pe.RecyclingReward ||
		pe.CarbonOffset)
}

// GetGreenDiscountRate returns environmental discount rate
func (pe *PolicyEnvironmental) GetGreenDiscountRate() float64 {
	if !pe.IsActive {
		return 0
	}

	discount := 0.0

	// Paperless discount
	if pe.PaperlessDiscount && pe.PaperlessEnrolled {
		discount += 0.02 // 2%
	}

	// Recycling reward
	if pe.RecyclingReward && pe.DevicesRecycled > 0 {
		discount += 0.03 // 3%
	}

	// Carbon offset
	if pe.CarbonOffset && pe.CarbonCredits > 0 {
		discount += 0.01 // 1%
	}

	// Green policy tier bonuses
	switch pe.GreenPolicyTier {
	case "platinum":
		discount += 0.05 // 5%
	case "gold":
		discount += 0.03 // 3%
	case "silver":
		discount += 0.02 // 2%
	case "basic":
		discount += 0.01 // 1%
	}

	// Cap at maximum green discount
	if discount > 0.15 { // 15% maximum
		discount = 0.15
	}

	return discount
}

// CalculateSustainabilityScore calculates overall sustainability score
func (pe *PolicyEnvironmental) CalculateSustainabilityScore() float64 {
	score := 50.0 // Base score

	// Policy features
	if pe.PaperlessEnrolled {
		score += 10
	}
	if pe.RecyclingReward {
		score += 15
	}
	if pe.CarbonOffset {
		score += 10
	}
	if pe.GreenPolicy {
		score += 5
	}

	// Eco services
	if pe.EcoRepairNetwork {
		score += 5
	}
	if pe.GreenShipping {
		score += 3
	}
	if pe.BiodegradablePackaging {
		score += 2
	}

	// Participation
	if pe.ClimateActionMember {
		score += 5
	}
	if pe.CircularEconomy {
		score += 5
	}
	if pe.ZeroWasteProgram {
		score += 5
	}

	// Actual impact
	if pe.DevicesRecycled > 0 {
		score += float64(pe.DevicesRecycled) * 2
	}
	if pe.CarbonCredits > 0 {
		score += pe.CarbonCredits * 0.1
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	pe.SustainabilityScore = score
	return score
}

// AddRecycledDevice records a recycled device
func (pe *PolicyEnvironmental) AddRecycledDevice() {
	pe.DevicesRecycled++
	pe.RecyclingCredits += 50 // $50 credit per device
	pe.EcoPoints += 100
	now := time.Now()
	pe.LastRecyclingDate = &now

	// Environmental impact
	pe.EWasteReduction += 0.5 // 0.5 kg average device weight
	pe.CarbonFootprint -= 50  // 50 kg CO2 saved
}

// GetEnvironmentalImpact returns the total environmental impact
func (pe *PolicyEnvironmental) GetEnvironmentalImpact() map[string]float64 {
	return map[string]float64{
		"carbon_saved":         -pe.CarbonFootprint, // Negative is good
		"ewaste_reduced":       pe.EWasteReduction,
		"plastic_reduced":      pe.PlasticReduction,
		"water_conserved":      pe.WaterConservation,
		"trees_saved":          pe.TreesSaved,
		"renewable_energy":     pe.RenewableEnergy,
		"sustainability_score": pe.SustainabilityScore,
	}
}

// IsEligibleForEcoReward checks eco reward eligibility
func (pe *PolicyEnvironmental) IsEligibleForEcoReward() bool {
	return pe.IsActive && pe.EcoPoints >= 500
}

// RedeemEcoPoints redeems eco points for rewards
func (pe *PolicyEnvironmental) RedeemEcoPoints(points int) bool {
	if pe.EcoPoints < points {
		return false
	}

	pe.EcoPoints -= points
	// Convert to monetary reward (1 point = $0.10)
	pe.EcoRewards += float64(points) * 0.10

	// Plant a tree for every 1000 points redeemed
	if points >= 1000 {
		pe.TreesPlanted += points / 1000
	}

	return true
}

// GetGreenTier returns the green policy tier based on score
func (pe *PolicyEnvironmental) GetGreenTier() string {
	score := pe.SustainabilityScore

	switch {
	case score >= 90:
		return "platinum"
	case score >= 75:
		return "gold"
	case score >= 60:
		return "silver"
	case score >= 40:
		return "basic"
	default:
		return "none"
	}
}

// NeedsAssessment checks if environmental assessment is due
func (pe *PolicyEnvironmental) NeedsAssessment() bool {
	if pe.NextAssessmentDate == nil {
		return true
	}
	return time.Now().After(*pe.NextAssessmentDate)
}

// UpdateCarbonFootprint updates carbon footprint based on activities
func (pe *PolicyEnvironmental) UpdateCarbonFootprint(activity string, value float64) {
	switch activity {
	case "device_purchase":
		pe.CarbonFootprint += value * 70 // 70 kg CO2 per device
	case "repair":
		pe.CarbonFootprint += value * 5 // 5 kg CO2 per repair
	case "shipping":
		if pe.GreenShipping {
			pe.CarbonFootprint += value * 0.5 // Reduced emissions
		} else {
			pe.CarbonFootprint += value * 2
		}
	case "recycling":
		pe.CarbonFootprint -= value * 50 // Saved emissions
	case "offset":
		pe.CarbonFootprint -= value
	}
}

// GetEcoProgress returns progress towards improvement target
func (pe *PolicyEnvironmental) GetEcoProgress() float64 {
	if pe.ImprovementTarget <= 0 {
		return 100
	}
	return (pe.ImprovementProgress / pe.ImprovementTarget) * 100
}
