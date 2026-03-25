package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyCoverage represents device-specific coverage features for a policy
type PolicyCoverage struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`

	// Screen & Damage Protection
	ScreenProtection      bool `gorm:"default:true" json:"screen_protection"`
	WaterDamageProtection bool `gorm:"default:false" json:"water_damage_protection"`
	AccidentalDamage      bool `gorm:"default:true" json:"accidental_damage"`
	TheftProtection       bool `gorm:"default:true" json:"theft_protection"`
	LossProtection        bool `gorm:"default:false" json:"loss_protection"`
	MechanicalBreakdown   bool `gorm:"default:false" json:"mechanical_breakdown"`
	BatteryProtection     bool `gorm:"default:false" json:"battery_protection"`

	// Environmental & Physical Damage
	FireDamage           bool `gorm:"default:false" json:"fire_damage"`
	PowerSurgeProtection bool `gorm:"default:false" json:"power_surge_protection"`
	NaturalDisasterCover bool `gorm:"default:false" json:"natural_disaster_cover"`
	ExtremeTempDamage    bool `gorm:"default:false" json:"extreme_temp_damage"`
	VandalismCoverage    bool `gorm:"default:false" json:"vandalism_coverage"`

	// Component-Specific Coverage
	CameraDamage         bool `gorm:"default:false" json:"camera_damage"`
	PortConnectorDamage  bool `gorm:"default:false" json:"port_connector_damage"`
	AudioComponentDamage bool `gorm:"default:false" json:"audio_component_damage"`
	CosmeticDamage       bool `gorm:"default:false" json:"cosmetic_damage"`
	TouchscreenFailure   bool `gorm:"default:false" json:"touchscreen_failure"`

	// Software & Security Coverage
	SoftwareMalfunction     bool `gorm:"default:false" json:"software_malfunction"`
	IdentityTheftProtection bool `gorm:"default:false" json:"identity_theft_protection"`
	UnauthorizedUsage       bool `gorm:"default:false" json:"unauthorized_usage"`
	CloudBackupService      bool `gorm:"default:false" json:"cloud_backup_service"`
	AppPurchaseProtection   bool `gorm:"default:false" json:"app_purchase_protection"`

	// Accessories & Additional Services
	AccessoriesCoverage   bool    `gorm:"default:false" json:"accessories_coverage"`
	AccessoriesLimit      float64 `json:"accessories_limit"`
	DataRecoveryService   bool    `gorm:"default:false" json:"data_recovery_service"`
	CyberSecurityCoverage bool    `gorm:"default:false" json:"cyber_security_coverage"`

	// Extended & Special Coverage
	PairedDeviceCoverage    bool `gorm:"default:false" json:"paired_device_coverage"`
	BusinessUseCoverage     bool `gorm:"default:false" json:"business_use_coverage"`
	MysteriousDisappearance bool `gorm:"default:false" json:"mysterious_disappearance"`
	WorldwideCoverage       bool `gorm:"default:false" json:"worldwide_coverage"`
	WarrantyExtension       bool `gorm:"default:false" json:"warranty_extension"`

	// Coverage Tier & Features
	CoverageTier      string `gorm:"type:varchar(20);default:'standard'" json:"coverage_tier"` // basic, standard, premium, platinum
	UnlimitedClaims   bool   `gorm:"default:false" json:"unlimited_claims"`
	ZeroDeductible    bool   `gorm:"default:false" json:"zero_deductible"`
	UpgradeProtection bool   `gorm:"default:false" json:"upgrade_protection"`
	BuybackGuarantee  bool   `gorm:"default:false" json:"buyback_guarantee"`

	// Coverage Limits
	ScreenRepairLimit    int     `gorm:"default:2" json:"screen_repair_limit"` // Per year
	WaterDamageLimit     int     `gorm:"default:1" json:"water_damage_limit"`  // Per year
	TotalClaimLimit      float64 `json:"total_claim_limit"`
	PerClaimLimit        float64 `json:"per_claim_limit"`
	ComponentRepairLimit int     `gorm:"default:3" json:"component_repair_limit"` // Per year

	// Status
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy uuid.UUID `gorm:"type:uuid" json:"modified_by"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyCoverage) TableName() string {
	return "policy_coverages"
}

// HasProtectionFor checks if specific damage type is covered
func (pc *PolicyCoverage) HasProtectionFor(damageType string) bool {
	if !pc.IsActive {
		return false
	}

	switch damageType {
	case "screen":
		return pc.ScreenProtection
	case "water":
		return pc.WaterDamageProtection
	case "accidental":
		return pc.AccidentalDamage
	case "theft":
		return pc.TheftProtection
	case "loss":
		return pc.LossProtection
	case "mechanical":
		return pc.MechanicalBreakdown
	case "battery":
		return pc.BatteryProtection
	case "fire":
		return pc.FireDamage
	case "power_surge":
		return pc.PowerSurgeProtection
	case "natural_disaster":
		return pc.NaturalDisasterCover
	case "extreme_temp":
		return pc.ExtremeTempDamage
	case "vandalism":
		return pc.VandalismCoverage
	case "camera":
		return pc.CameraDamage
	case "port_connector":
		return pc.PortConnectorDamage
	case "audio":
		return pc.AudioComponentDamage
	case "cosmetic":
		return pc.CosmeticDamage
	case "touchscreen":
		return pc.TouchscreenFailure
	case "software":
		return pc.SoftwareMalfunction
	case "identity_theft":
		return pc.IdentityTheftProtection
	case "unauthorized_usage":
		return pc.UnauthorizedUsage
	case "mysterious_disappearance":
		return pc.MysteriousDisappearance
	default:
		return false
	}
}

// CanClaimAccessories checks if accessories can be claimed
func (pc *PolicyCoverage) CanClaimAccessories(claimAmount float64) bool {
	return pc.IsActive && pc.AccessoriesCoverage && claimAmount <= pc.AccessoriesLimit
}

// GetCoverageScore returns a score based on coverage comprehensiveness
func (pc *PolicyCoverage) GetCoverageScore() int {
	score := 0

	// Core coverage (50 points)
	if pc.ScreenProtection {
		score += 10
	}
	if pc.WaterDamageProtection {
		score += 10
	}
	if pc.AccidentalDamage {
		score += 10
	}
	if pc.TheftProtection {
		score += 10
	}
	if pc.LossProtection {
		score += 10
	}

	// Component coverage (20 points)
	if pc.MechanicalBreakdown {
		score += 5
	}
	if pc.BatteryProtection {
		score += 3
	}
	if pc.CameraDamage {
		score += 4
	}
	if pc.PortConnectorDamage {
		score += 3
	}
	if pc.AudioComponentDamage {
		score += 3
	}
	if pc.TouchscreenFailure {
		score += 2
	}

	// Environmental coverage (15 points)
	if pc.FireDamage {
		score += 5
	}
	if pc.PowerSurgeProtection {
		score += 3
	}
	if pc.NaturalDisasterCover {
		score += 4
	}
	if pc.ExtremeTempDamage {
		score += 2
	}
	if pc.VandalismCoverage {
		score += 1
	}

	// Digital protection (10 points)
	if pc.DataRecoveryService {
		score += 2
	}
	if pc.CyberSecurityCoverage {
		score += 2
	}
	if pc.CloudBackupService {
		score += 2
	}
	if pc.IdentityTheftProtection {
		score += 2
	}
	if pc.SoftwareMalfunction {
		score += 2
	}

	// Premium features (5 points)
	if pc.UnlimitedClaims {
		score += 2
	}
	if pc.ZeroDeductible {
		score += 2
	}
	if pc.WorldwideCoverage {
		score += 1
	}

	return score
}

// IsComprehensive checks if this is comprehensive coverage
func (pc *PolicyCoverage) IsComprehensive() bool {
	return pc.ScreenProtection &&
		pc.WaterDamageProtection &&
		pc.AccidentalDamage &&
		pc.TheftProtection &&
		pc.LossProtection
}

// HasDataServices checks if data-related services are covered
func (pc *PolicyCoverage) HasDataServices() bool {
	return pc.DataRecoveryService || pc.CyberSecurityCoverage ||
		pc.CloudBackupService || pc.IdentityTheftProtection
}

// GetCoverageTierFeatures returns the features included in the coverage tier
func (pc *PolicyCoverage) GetCoverageTierFeatures() map[string]bool {
	features := make(map[string]bool)

	switch pc.CoverageTier {
	case "basic":
		features["screen"] = true
		features["accidental_damage"] = true

	case "standard":
		features["screen"] = true
		features["accidental_damage"] = true
		features["water_damage"] = true
		features["theft"] = true
		features["battery"] = true

	case "premium":
		features["screen"] = true
		features["accidental_damage"] = true
		features["water_damage"] = true
		features["theft"] = true
		features["loss"] = true
		features["battery"] = true
		features["mechanical_breakdown"] = true
		features["accessories"] = true
		features["camera_damage"] = true

	case "platinum":
		features["all_physical_damage"] = true
		features["all_components"] = true
		features["all_environmental"] = true
		features["all_digital"] = true
		features["unlimited_claims"] = true
		features["zero_deductible"] = true
		features["worldwide"] = true
	}

	return features
}

// IsPremiumTier checks if this is a premium or platinum tier
func (pc *PolicyCoverage) IsPremiumTier() bool {
	return pc.CoverageTier == "premium" || pc.CoverageTier == "platinum"
}

// HasEnvironmentalProtection checks if environmental damage is covered
func (pc *PolicyCoverage) HasEnvironmentalProtection() bool {
	return pc.FireDamage || pc.PowerSurgeProtection ||
		pc.NaturalDisasterCover || pc.ExtremeTempDamage
}

// HasComponentProtection checks if component-specific damage is covered
func (pc *PolicyCoverage) HasComponentProtection() bool {
	return pc.CameraDamage || pc.PortConnectorDamage ||
		pc.AudioComponentDamage || pc.TouchscreenFailure
}

// HasDigitalProtection checks if digital/software protection is included
func (pc *PolicyCoverage) HasDigitalProtection() bool {
	return pc.SoftwareMalfunction || pc.IdentityTheftProtection ||
		pc.UnauthorizedUsage || pc.CloudBackupService ||
		pc.AppPurchaseProtection || pc.CyberSecurityCoverage
}

// HasSpecialCoverage checks for special coverage options
func (pc *PolicyCoverage) HasSpecialCoverage() bool {
	return pc.PairedDeviceCoverage || pc.BusinessUseCoverage ||
		pc.MysteriousDisappearance || pc.WorldwideCoverage ||
		pc.WarrantyExtension
}

// GetClaimLimitForType returns the claim limit for a specific damage type
func (pc *PolicyCoverage) GetClaimLimitForType(damageType string) int {
	if pc.UnlimitedClaims {
		return 999 // Effectively unlimited
	}

	switch damageType {
	case "screen":
		return pc.ScreenRepairLimit
	case "water":
		return pc.WaterDamageLimit
	case "component":
		return pc.ComponentRepairLimit
	default:
		return 2 // Default limit
	}
}

// ValidateClaimAmount checks if claim amount is within limits
func (pc *PolicyCoverage) ValidateClaimAmount(amount float64) bool {
	if pc.PerClaimLimit > 0 && amount > pc.PerClaimLimit {
		return false
	}

	if pc.TotalClaimLimit > 0 && amount > pc.TotalClaimLimit {
		return false
	}

	return true
}

// UpgradeTier upgrades the coverage tier
func (pc *PolicyCoverage) UpgradeTier() string {
	switch pc.CoverageTier {
	case "basic":
		pc.CoverageTier = "standard"
	case "standard":
		pc.CoverageTier = "premium"
	case "premium":
		pc.CoverageTier = "platinum"
	}
	return pc.CoverageTier
}
