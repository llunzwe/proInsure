package grading

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DeviceGrading provides device grading services
type DeviceGrading struct {
	config *GradingConfig
}

// ValueAssessor provides device value assessment services
type ValueAssessor struct {
	config *GradingConfig
}

// RiskAssessor provides risk assessment services
type RiskAssessor struct {
	config *GradingConfig
}

// GradingConfig holds grading configuration
type GradingConfig struct {
	EnablePhysicalGrading bool
	EnableValueAssessment bool
	EnableRiskAssessment  bool
	MaxAssessmentTime     time.Duration
}

// GradingResult represents the result of device grading
type GradingResult struct {
	ID             uuid.UUID              `json:"id"`
	DeviceID       uuid.UUID              `json:"deviceId"`
	Grade          string                 `json:"grade"`
	Score          float64                `json:"score"`
	EstimatedValue float64                `json:"estimatedValue"`
	RiskLevel      string                 `json:"riskLevel"`
	Condition      string                 `json:"condition"`
	Timestamp      time.Time              `json:"timestamp"`
	Details        map[string]interface{} `json:"details"`
}

// NewDeviceGrading creates a new device grading service
func NewDeviceGrading(config *GradingConfig) *DeviceGrading {
	return &DeviceGrading{
		config: config,
	}
}

// NewValueAssessor creates a new value assessor
func NewValueAssessor(config *GradingConfig) *ValueAssessor {
	return &ValueAssessor{
		config: config,
	}
}

// NewRiskAssessor creates a new risk assessor
func NewRiskAssessor(config *GradingConfig) *RiskAssessor {
	return &RiskAssessor{
		config: config,
	}
}

// GradeDevice grades a device based on its condition and specifications
func (g *DeviceGrading) GradeDevice(ctx context.Context, deviceData map[string]interface{}) (*GradingResult, error) {
	if !g.config.EnablePhysicalGrading {
		return &GradingResult{
			ID:        uuid.New(),
			Grade:     "UNKNOWN",
			Score:     0.5,
			Condition: "NOT_ASSESSED",
			Timestamp: time.Now(),
		}, nil
	}

	// Basic grading logic
	score := 1.0
	grade := "EXCELLENT"
	condition := "NEW"

	// Check device age
	if age, ok := deviceData["age_months"].(float64); ok {
		if age > 24 {
			score -= 0.3
			grade = "GOOD"
			condition = "USED"
		} else if age > 12 {
			score -= 0.1
			grade = "VERY_GOOD"
			condition = "LIGHTLY_USED"
		}
	}

	// Check physical condition
	if physicalCondition, ok := deviceData["physical_condition"].(string); ok {
		switch physicalCondition {
		case "DAMAGED":
			score -= 0.5
			grade = "POOR"
			condition = "DAMAGED"
		case "WORN":
			score -= 0.2
			grade = "FAIR"
			condition = "WORN"
		}
	}

	return &GradingResult{
		ID:        uuid.New(),
		Grade:     grade,
		Score:     score,
		Condition: condition,
		Timestamp: time.Now(),
		Details:   deviceData,
	}, nil
}

// AssessValue assesses the market value of a device
func (v *ValueAssessor) AssessValue(ctx context.Context, deviceData map[string]interface{}) (float64, error) {
	if !v.config.EnableValueAssessment {
		return 0.0, nil
	}

	baseValue := 500.0 // Default base value

	// Get original price if available
	if originalPrice, ok := deviceData["original_price"].(float64); ok {
		baseValue = originalPrice
	}

	// Apply depreciation based on age
	if age, ok := deviceData["age_months"].(float64); ok {
		depreciationRate := 0.05 // 5% per month
		depreciation := age * depreciationRate
		if depreciation > 0.8 {
			depreciation = 0.8 // Max 80% depreciation
		}
		baseValue = baseValue * (1 - depreciation)
	}

	// Apply condition adjustments
	if condition, ok := deviceData["condition"].(string); ok {
		switch condition {
		case "DAMAGED":
			baseValue *= 0.3
		case "WORN":
			baseValue *= 0.6
		case "USED":
			baseValue *= 0.8
		case "LIGHTLY_USED":
			baseValue *= 0.9
		}
	}

	return baseValue, nil
}

// AssessRisk assesses the risk level of insuring a device
func (r *RiskAssessor) AssessRisk(ctx context.Context, deviceData map[string]interface{}, userProfile map[string]interface{}) (string, float64, error) {
	if !r.config.EnableRiskAssessment {
		return "MEDIUM", 0.5, nil
	}

	riskScore := 0.0

	// Device-based risk factors
	if age, ok := deviceData["age_months"].(float64); ok {
		if age > 24 {
			riskScore += 0.3
		} else if age > 12 {
			riskScore += 0.1
		}
	}

	if condition, ok := deviceData["condition"].(string); ok {
		switch condition {
		case "DAMAGED":
			riskScore += 0.5
		case "WORN":
			riskScore += 0.2
		}
	}

	// User-based risk factors
	if claimHistory, ok := userProfile["previous_claims"].(float64); ok {
		if claimHistory > 2 {
			riskScore += 0.3
		} else if claimHistory > 0 {
			riskScore += 0.1
		}
	}

	// Determine risk level
	riskLevel := "LOW"
	if riskScore > 0.7 {
		riskLevel = "HIGH"
	} else if riskScore > 0.3 {
		riskLevel = "MEDIUM"
	}

	return riskLevel, riskScore, nil
}
