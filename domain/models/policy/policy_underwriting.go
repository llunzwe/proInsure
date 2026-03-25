package policy

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyUnderwriting represents underwriting details for a policy
type PolicyUnderwriting struct {
	database.BaseModel
	PolicyID       uuid.UUID `gorm:"type:uuid;not null;index" json:"policy_id"`
	UnderwritingID string    `gorm:"uniqueIndex;not null" json:"underwriting_id"`
	Status         string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Method         string    `gorm:"type:varchar(50)" json:"method"` // auto, manual, hybrid

	// Assessment
	RiskScore    float64 `json:"risk_score"`
	RiskCategory string  `json:"risk_category"`
	RiskFactors  string  `gorm:"type:json" json:"risk_factors"`

	// Device Assessment
	DeviceScore        float64    `json:"device_score"`
	DeviceCondition    string     `json:"device_condition"`
	DeviceValue        float64    `json:"device_value"`
	InspectionRequired bool       `gorm:"default:false" json:"inspection_required"`
	InspectionStatus   string     `json:"inspection_status"`
	InspectionDate     *time.Time `json:"inspection_date"`

	// User Assessment
	UserScore        float64 `json:"user_score"`
	CreditScore      int     `json:"credit_score"`
	ClaimsHistory    string  `gorm:"type:json" json:"claims_history"`
	PreviousPolicies int     `json:"previous_policies"`

	// Location Assessment
	LocationRisk float64 `json:"location_risk"`
	TheftRate    float64 `json:"theft_rate"`
	DisasterRisk float64 `json:"disaster_risk"`

	// Decision
	Decision       string    `gorm:"type:varchar(20)" json:"decision"` // approved, declined, refer
	DecisionDate   time.Time `json:"decision_date"`
	UnderwriterID  uuid.UUID `gorm:"type:uuid" json:"underwriter_id"`
	DecisionReason string    `json:"decision_reason"`

	// Conditions
	Conditions string `gorm:"type:json" json:"conditions"`
	Loadings   string `gorm:"type:json" json:"loadings"`
	Discounts  string `gorm:"type:json" json:"discounts"`

	// Premium Calculation
	BasePremium     float64 `json:"base_premium"`
	AdjustedPremium float64 `json:"adjusted_premium"`
	FinalPremium    float64 `json:"final_premium"`

	// Documentation
	Documents string `gorm:"type:json" json:"documents"`
	Notes     string `json:"notes"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
	// Underwriter relationship should be loaded via service layer using UnderwriterID to avoid circular import
}

// TableName returns the table name
func (PolicyUnderwriting) TableName() string {
	return "policy_underwritings"
}

// RequiresManualReview checks if manual review is needed
func (pu *PolicyUnderwriting) RequiresManualReview() bool {
	return pu.Method == "manual" ||
		pu.RiskScore > 75 ||
		pu.Decision == "refer"
}

// IsHighRisk checks if underwriting indicates high risk
func (pu *PolicyUnderwriting) IsHighRisk() bool {
	return pu.RiskCategory == "high" || pu.RiskCategory == "very_high"
}

// IsApproved checks if underwriting is approved
func (pu *PolicyUnderwriting) IsApproved() bool {
	return pu.Decision == "approved" && pu.Status == "completed"
}

// IsDeclined checks if underwriting is declined
func (pu *PolicyUnderwriting) IsDeclined() bool {
	return pu.Decision == "declined"
}

// CanAutoApprove checks if policy can be auto-approved
func (pu *PolicyUnderwriting) CanAutoApprove() bool {
	return pu.Method == "auto" &&
		pu.RiskScore < 50 &&
		pu.CreditScore > 650 &&
		!pu.InspectionRequired
}

// CalculatePremiumAdjustment calculates premium adjustment based on risk
func (pu *PolicyUnderwriting) CalculatePremiumAdjustment() float64 {
	adjustment := 1.0

	// Risk score adjustment
	if pu.RiskScore > 75 {
		adjustment += 0.30 // 30% increase for very high risk
	} else if pu.RiskScore > 50 {
		adjustment += 0.15 // 15% increase for high risk
	} else if pu.RiskScore < 25 {
		adjustment -= 0.10 // 10% discount for low risk
	}

	// Credit score adjustment
	if pu.CreditScore > 750 {
		adjustment -= 0.05 // 5% discount for excellent credit
	} else if pu.CreditScore < 600 {
		adjustment += 0.10 // 10% increase for poor credit
	}

	// Claims history adjustment
	if pu.PreviousPolicies > 0 && pu.ClaimsHistory == "" {
		adjustment -= 0.05 // 5% discount for claim-free history
	}

	// Location risk adjustment
	if pu.LocationRisk > 70 {
		adjustment += 0.10 // 10% increase for high-risk location
	}

	return adjustment
}

// GetProcessingTime returns time taken for underwriting
func (pu *PolicyUnderwriting) GetProcessingTime() time.Duration {
	if pu.Status != "completed" {
		return 0
	}
	return pu.DecisionDate.Sub(pu.CreatedAt)
}

// NeedsInspection checks if inspection is required
func (pu *PolicyUnderwriting) NeedsInspection() bool {
	return pu.InspectionRequired && pu.InspectionStatus != "completed"
}

// IsExpired checks if underwriting decision has expired
func (pu *PolicyUnderwriting) IsExpired() bool {
	// Underwriting decisions expire after 30 days
	expiryDate := pu.DecisionDate.AddDate(0, 0, 30)
	return time.Now().After(expiryDate)
}

// GetRiskLevel returns the risk level category
func (pu *PolicyUnderwriting) GetRiskLevel() string {
	switch {
	case pu.RiskScore < 25:
		return "low"
	case pu.RiskScore < 50:
		return "medium"
	case pu.RiskScore < 75:
		return "high"
	default:
		return "very_high"
	}
}

// ShouldRefer checks if case should be referred to senior underwriter
func (pu *PolicyUnderwriting) ShouldRefer() bool {
	return pu.RiskScore > 70 ||
		pu.DeviceValue > 5000 ||
		pu.CreditScore < 550 ||
		pu.Method == "hybrid"
}
