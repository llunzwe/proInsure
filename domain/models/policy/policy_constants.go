package policy

// Coverage Tier Constants - Used across Product, Policy, and Device models
const (
	// Tier Codes (string)
	TierBasic    = "basic"
	TierStandard = "standard"
	TierPremium  = "premium"
	TierPlatinum = "platinum"

	// Tier Levels (int) - For sorting and comparison
	TierLevelBasic    = 1
	TierLevelStandard = 2
	TierLevelPremium  = 3
	TierLevelPlatinum = 4
)

// Premium Type Constants
const (
	PremiumTypeBase         = "base"
	PremiumTypeRiskAdjusted = "risk_adjusted"
	PremiumTypeFinal        = "final"
	PremiumTypeMonthly      = "monthly"
	PremiumTypeQuarterly    = "quarterly"
	PremiumTypeAnnual       = "annual"
)

// NOTE: Deductible Type constants are defined in policy_types.go

// Coverage Limit Type Constants
const (
	LimitTypePerClaim  = "per_claim"
	LimitTypePerYear   = "per_year"
	LimitTypeLifetime  = "lifetime"
	LimitTypeUnlimited = "unlimited"
)

// NOTE: PolicyStatus constants are defined in policy_types.go

// GetTierLevel returns the numeric level for a tier code
func GetTierLevel(tierCode string) int {
	switch tierCode {
	case TierBasic:
		return TierLevelBasic
	case TierStandard:
		return TierLevelStandard
	case TierPremium:
		return TierLevelPremium
	case TierPlatinum:
		return TierLevelPlatinum
	default:
		return 0
	}
}

// GetTierCode returns the string code for a tier level
func GetTierCode(tierLevel int) string {
	switch tierLevel {
	case TierLevelBasic:
		return TierBasic
	case TierLevelStandard:
		return TierStandard
	case TierLevelPremium:
		return TierPremium
	case TierLevelPlatinum:
		return TierPlatinum
	default:
		return ""
	}
}

// IsPremiumTier checks if the tier is premium or platinum
func IsPremiumTier(tierCode string) bool {
	return tierCode == TierPremium || tierCode == TierPlatinum
}

// CompareTiers returns -1 if tier1 < tier2, 0 if equal, 1 if tier1 > tier2
func CompareTiers(tier1, tier2 string) int {
	level1 := GetTierLevel(tier1)
	level2 := GetTierLevel(tier2)

	if level1 < level2 {
		return -1
	} else if level1 > level2 {
		return 1
	}
	return 0
}
