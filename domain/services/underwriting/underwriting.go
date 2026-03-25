package underwriting

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// RiskCalculator provides risk calculation services
type RiskCalculator struct {
	config *UnderwritingConfig
}

// PremiumCalculator provides premium calculation services
type PremiumCalculator struct {
	config *UnderwritingConfig
}

// PolicyValidator provides policy validation services
type PolicyValidator struct {
	config *UnderwritingConfig
}

// UnderwritingConfig holds underwriting configuration
type UnderwritingConfig struct {
	EnableRiskCalculation    bool
	EnablePremiumCalculation bool
	EnablePolicyValidation   bool
	BasePremiumRate          float64
	RiskMultipliers          map[string]float64
	MaxCoverageAmount        float64
	MinCoverageAmount        float64
}

// RiskAssessment represents a risk assessment result
type RiskAssessment struct {
	ID              uuid.UUID `json:"id"`
	CustomerID      uuid.UUID `json:"customerId"`
	DeviceID        uuid.UUID `json:"deviceId"`
	RiskScore       float64   `json:"riskScore"`
	RiskLevel       string    `json:"riskLevel"`
	RiskFactors     []string  `json:"riskFactors"`
	Recommendations []string  `json:"recommendations"`
	Timestamp       time.Time `json:"timestamp"`
}

// PremiumCalculation represents a premium calculation result
type PremiumCalculation struct {
	ID               uuid.UUID `json:"id"`
	BaseAmount       float64   `json:"baseAmount"`
	RiskAdjustment   float64   `json:"riskAdjustment"`
	Discounts        float64   `json:"discounts"`
	Taxes            float64   `json:"taxes"`
	TotalPremium     float64   `json:"totalPremium"`
	PaymentFrequency string    `json:"paymentFrequency"`
	Timestamp        time.Time `json:"timestamp"`
}

// PolicyValidation represents a policy validation result
type PolicyValidation struct {
	ID        uuid.UUID `json:"id"`
	IsValid   bool      `json:"isValid"`
	Errors    []string  `json:"errors"`
	Warnings  []string  `json:"warnings"`
	Timestamp time.Time `json:"timestamp"`
}

// NewRiskCalculator creates a new risk calculator
func NewRiskCalculator(config *UnderwritingConfig) *RiskCalculator {
	return &RiskCalculator{
		config: config,
	}
}

// NewPremiumCalculator creates a new premium calculator
func NewPremiumCalculator(config *UnderwritingConfig) *PremiumCalculator {
	return &PremiumCalculator{
		config: config,
	}
}

// NewPolicyValidator creates a new policy validator
func NewPolicyValidator(config *UnderwritingConfig) *PolicyValidator {
	return &PolicyValidator{
		config: config,
	}
}

// CalculateRisk calculates risk for a customer and device combination
func (r *RiskCalculator) CalculateRisk(ctx context.Context, customerData, deviceData map[string]interface{}) (*RiskAssessment, error) {
	if !r.config.EnableRiskCalculation {
		return &RiskAssessment{
			ID:        uuid.New(),
			RiskScore: 0.5,
			RiskLevel: "MEDIUM",
			Timestamp: time.Now(),
		}, nil
	}

	riskScore := 0.0
	riskFactors := []string{}
	recommendations := []string{}

	// Customer risk factors
	if age, ok := customerData["age"].(float64); ok {
		if age < 25 {
			riskScore += 0.2
			riskFactors = append(riskFactors, "Young driver profile")
		} else if age > 65 {
			riskScore += 0.1
			riskFactors = append(riskFactors, "Senior customer profile")
		}
	}

	if claimHistory, ok := customerData["previous_claims"].(float64); ok {
		if claimHistory > 2 {
			riskScore += 0.3
			riskFactors = append(riskFactors, "High claim history")
			recommendations = append(recommendations, "Consider higher deductible")
		} else if claimHistory > 0 {
			riskScore += 0.1
			riskFactors = append(riskFactors, "Previous claims")
		}
	}

	// Device risk factors
	if deviceAge, ok := deviceData["age_months"].(float64); ok {
		if deviceAge > 24 {
			riskScore += 0.2
			riskFactors = append(riskFactors, "Older device")
		}
	}

	if deviceValue, ok := deviceData["value"].(float64); ok {
		if deviceValue > 1000 {
			riskScore += 0.1
			riskFactors = append(riskFactors, "High-value device")
		}
	}

	// Determine risk level
	riskLevel := "LOW"
	if riskScore > 0.7 {
		riskLevel = "HIGH"
		recommendations = append(recommendations, "Consider additional security measures")
	} else if riskScore > 0.3 {
		riskLevel = "MEDIUM"
	}

	return &RiskAssessment{
		ID:              uuid.New(),
		RiskScore:       math.Min(riskScore, 1.0),
		RiskLevel:       riskLevel,
		RiskFactors:     riskFactors,
		Recommendations: recommendations,
		Timestamp:       time.Now(),
	}, nil
}

// CalculatePremium calculates premium based on risk assessment and coverage
func (p *PremiumCalculator) CalculatePremium(ctx context.Context, riskAssessment *RiskAssessment, coverageAmount float64, policyTerm int) (*PremiumCalculation, error) {
	if !p.config.EnablePremiumCalculation {
		return &PremiumCalculation{
			ID:           uuid.New(),
			BaseAmount:   100.0,
			TotalPremium: 100.0,
			Timestamp:    time.Now(),
		}, nil
	}

	// Base premium calculation
	baseRate := p.config.BasePremiumRate
	if baseRate == 0 {
		baseRate = 0.05 // 5% of device value annually
	}

	baseAmount := coverageAmount * baseRate * float64(policyTerm) / 12.0

	// Risk adjustment
	riskMultiplier := 1.0
	switch riskAssessment.RiskLevel {
	case "HIGH":
		riskMultiplier = 1.5
	case "MEDIUM":
		riskMultiplier = 1.2
	case "LOW":
		riskMultiplier = 0.9
	}

	riskAdjustment := baseAmount * (riskMultiplier - 1.0)

	// Calculate discounts (placeholder logic)
	discounts := 0.0
	if riskAssessment.RiskLevel == "LOW" {
		discounts = baseAmount * 0.1 // 10% discount for low risk
	}

	// Calculate taxes (placeholder - typically varies by jurisdiction)
	subtotal := baseAmount + riskAdjustment - discounts
	taxes := subtotal * 0.08 // 8% tax

	totalPremium := subtotal + taxes

	return &PremiumCalculation{
		ID:               uuid.New(),
		BaseAmount:       baseAmount,
		RiskAdjustment:   riskAdjustment,
		Discounts:        discounts,
		Taxes:            taxes,
		TotalPremium:     totalPremium,
		PaymentFrequency: "MONTHLY",
		Timestamp:        time.Now(),
	}, nil
}

// ValidatePolicy validates a policy before issuance
func (v *PolicyValidator) ValidatePolicy(ctx context.Context, policyData map[string]interface{}) (*PolicyValidation, error) {
	if !v.config.EnablePolicyValidation {
		return &PolicyValidation{
			ID:        uuid.New(),
			IsValid:   true,
			Timestamp: time.Now(),
		}, nil
	}

	errors := []string{}
	warnings := []string{}

	// Validate coverage amount
	if coverageAmount, ok := policyData["coverage_amount"].(float64); ok {
		if coverageAmount > v.config.MaxCoverageAmount {
			errors = append(errors, fmt.Sprintf("Coverage amount exceeds maximum allowed: %.2f", v.config.MaxCoverageAmount))
		}
		if coverageAmount < v.config.MinCoverageAmount {
			errors = append(errors, fmt.Sprintf("Coverage amount below minimum required: %.2f", v.config.MinCoverageAmount))
		}
	} else {
		errors = append(errors, "Coverage amount is required")
	}

	// Validate policy term
	if policyTerm, ok := policyData["policy_term"].(float64); ok {
		if policyTerm < 1 || policyTerm > 36 {
			errors = append(errors, "Policy term must be between 1 and 36 months")
		}
	} else {
		errors = append(errors, "Policy term is required")
	}

	// Validate customer information
	if customerID, ok := policyData["customer_id"].(string); !ok || customerID == "" {
		errors = append(errors, "Customer ID is required")
	}

	// Validate device information
	if deviceID, ok := policyData["device_id"].(string); !ok || deviceID == "" {
		errors = append(errors, "Device ID is required")
	}

	// Add warnings for potential issues
	if deductible, ok := policyData["deductible"].(float64); ok {
		if coverageAmount, ok := policyData["coverage_amount"].(float64); ok {
			if deductible > coverageAmount*0.5 {
				warnings = append(warnings, "Deductible is more than 50% of coverage amount")
			}
		}
	}

	return &PolicyValidation{
		ID:        uuid.New(),
		IsValid:   len(errors) == 0,
		Errors:    errors,
		Warnings:  warnings,
		Timestamp: time.Now(),
	}, nil
}
