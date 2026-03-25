package device

import (
	"math"
	"time"

	"smartsure/internal/domain/models"
)

// GetEnvironmentalImpact returns environmental metrics for the device
func (d *models.Device) GetEnvironmentalImpact() map[string]interface{} {
	impact := map[string]interface{}{
		"carbon_footprint_kg":    d.CalculateCarbonFootprint(),
		"recyclability_score":    d.CalculateRecyclabilityScore(),
		"repairability_score":    d.CalculateRepairabilityScore(),
		"lifecycle_extension":    d.CalculateLifecycleExtension(),
		"ewaste_prevention_kg":   d.CalculateEWastePrevention(),
		"resource_efficiency":    d.CalculateResourceEfficiency(),
		"environmental_grade":    d.GetEnvironmentalGrade(),
		"sustainability_actions": d.GetSustainabilityActions(),
	}

	return impact
}

// IsEndOfLife checks if device has reached end of lifecycle
func (d *models.Device) IsEndOfLife() bool {
	// Check age
	if d.GetDeviceAge() > 2555 { // 7+ years
		return true
	}

	// Check condition
	if d.Condition == "poor" && d.Grade == "F" {
		return true
	}

	// Check repair cost vs value
	estimatedRepairCost := d.EstimateRepairCost("motherboard")
	if estimatedRepairCost > d.CurrentValue*0.7 {
		return true
	}

	// Check battery health
	if d.BatteryHealth > 0 && d.BatteryHealth < 50 {
		return true
	}

	// Check OS support
	if d.OSVersion != "" {
		// Simulate checking if OS is still supported
		// In production, this would check against manufacturer data
		deviceAge := d.GetDeviceAge() / 365
		if deviceAge > 5 {
			return true
		}
	}

	// Check parts availability
	if d.PartsAvailability == "discontinued" {
		return true
	}

	return false
}

// GetRecyclingValue estimates the value for recycling/parts recovery
func (d *models.Device) GetRecyclingValue() float64 {
	// Base recycling value for materials
	materialValue := 5.0 // Base value for plastics and common metals

	// Add precious metals value based on device type
	switch d.DeviceSegment {
	case "flagship":
		materialValue += 15 // More gold, silver, palladium
	case "premium":
		materialValue += 10
	case "mid_range":
		materialValue += 7
	case "budget":
		materialValue += 3
	}

	// Add value for rare earth elements
	if d.Is5GCapable {
		materialValue += 5 // More rare earth elements in 5G devices
	}

	// Add battery recycling value
	if d.BatteryCapacity > 0 {
		batteryValue := float64(d.BatteryCapacity) * 0.002 // $0.002 per mAh
		materialValue += batteryValue
	}

	// Add component salvage value if components are working
	if !d.IsEndOfLife() {
		salvageValue := d.CalculateSalvageValue()
		materialValue += salvageValue * 0.3 // 30% of salvage value for recycling
	}

	// Reduce value for damaged devices
	if d.WaterDamageIndicator != "white" {
		materialValue *= 0.5
	}

	// Environmental bonus for proper recycling
	if d.ProofOfOwnershipVerified {
		materialValue *= 1.1 // 10% bonus for verified devices
	}

	return materialValue
}

// CalculateCarbonFootprint estimates carbon footprint in kg CO2
func (d *models.Device) CalculateCarbonFootprint() float64 {
	// Base manufacturing footprint
	baseFootprint := 70.0 // kg CO2 average for smartphone manufacturing

	// Adjust based on device segment
	switch d.DeviceSegment {
	case "flagship":
		baseFootprint = 85 // More complex manufacturing
	case "premium":
		baseFootprint = 75
	case "mid_range":
		baseFootprint = 70
	case "budget":
		baseFootprint = 60
	}

	// Add usage footprint (2kg CO2 per year of use)
	yearsUsed := float64(d.GetDeviceAge()) / 365
	usageFootprint := yearsUsed * 2

	// Add repair footprint
	repairFootprint := float64(d.GetRepairCount()) * 5 // 5kg CO2 per repair

	// Calculate total
	totalFootprint := baseFootprint + usageFootprint + repairFootprint

	// Reduce if device is being reused/refurbished
	if len(d.Refurbishments) > 0 {
		totalFootprint *= 0.5 // Footprint shared with next owner
	}

	return totalFootprint
}

// GetRefurbishmentPotential assesses refurbishment viability
func (d *models.Device) GetRefurbishmentPotential() string {
	score := 100.0

	// Deduct for age
	ageYears := float64(d.GetDeviceAge()) / 365
	if ageYears > 4 {
		score -= 40
	} else if ageYears > 3 {
		score -= 25
	} else if ageYears > 2 {
		score -= 15
	} else if ageYears > 1 {
		score -= 5
	}

	// Deduct for condition
	conditionScore := d.GetConditionScore()
	score -= (1 - conditionScore) * 30

	// Deduct for repairs needed
	if d.ScreenCondition == "broken" {
		score -= 20
	} else if d.ScreenCondition == "cracked" {
		score -= 10
	}

	if d.BatteryHealth < 70 && d.BatteryHealth > 0 {
		score -= 15
	}

	// Deduct for water damage
	if d.WaterDamageIndicator != "white" {
		score -= 30
	}

	// Deduct for missing features
	if !d.OriginalBox {
		score -= 5
	}

	// Categorize potential
	if score >= 80 {
		return "excellent"
	} else if score >= 65 {
		return "good"
	} else if score >= 50 {
		return "moderate"
	} else if score >= 30 {
		return "low"
	} else {
		return "not_viable"
	}
}

// PredictRemainingLifespan predicts remaining useful days
func (d *models.Device) PredictRemainingLifespan() int {
	// Base lifespan in days
	baseLifespan := 2555 // 7 years

	// Adjust based on segment
	switch d.DeviceSegment {
	case "flagship":
		baseLifespan = 2920 // 8 years
	case "premium":
		baseLifespan = 2555 // 7 years
	case "mid_range":
		baseLifespan = 2190 // 6 years
	case "budget":
		baseLifespan = 1825 // 5 years
	}

	// Current age
	currentAge := d.GetDeviceAge()

	// Calculate remaining base lifespan
	remaining := baseLifespan - currentAge

	// Adjust for condition
	conditionMultiplier := d.GetConditionScore()
	remaining = int(float64(remaining) * conditionMultiplier)

	// Adjust for battery health
	if d.BatteryHealth > 0 && d.BatteryHealth < 80 {
		batteryMultiplier := float64(d.BatteryHealth) / 100
		remaining = int(float64(remaining) * batteryMultiplier)
	}

	// Adjust for repair history
	repairPenalty := d.GetRepairCount() * 90 // 90 days per major repair
	remaining -= repairPenalty

	// Adjust for water damage
	if d.WaterDamageIndicator != "white" {
		remaining = int(float64(remaining) * 0.5)
	}

	// Minimum lifespan if device is still functional
	if remaining < 180 && d.Status == "active" && !d.IsEndOfLife() {
		remaining = 180 // At least 6 months
	}

	// Can't be negative
	if remaining < 0 {
		remaining = 0
	}

	return remaining
}

// CalculateRecyclabilityScore calculates how recyclable the device is
func (d *models.Device) CalculateRecyclabilityScore() float64 {
	score := 70.0 // Base recyclability for modern smartphones

	// Add points for valuable materials
	if d.DeviceSegment == "flagship" || d.DeviceSegment == "premium" {
		score += 10 // More precious metals
	}

	// Add points for battery condition
	if d.BatteryHealth > 70 {
		score += 5
	}

	// Add points for modular design (simplified assessment)
	if d.Brand == "Fairphone" { // Example of modular design
		score += 15
	}

	// Deduct for hazardous damage
	if d.WaterDamageIndicator != "white" {
		score -= 10 // Potential chemical contamination
	}

	// Deduct for mixed materials
	if d.WaterResistance == "IP68" {
		score -= 5 // More adhesives and seals
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// CalculateRepairabilityScore calculates device repairability
func (d *models.Device) CalculateRepairabilityScore() float64 {
	score := 50.0 // Base score

	// Add points for parts availability
	switch d.PartsAvailability {
	case "readily_available":
		score += 20
	case "limited":
		score += 10
	case "scarce":
		score += 5
	}

	// Add points for repair network
	if d.RepairNetworkCompatible {
		score += 10
	}

	// Deduct for complex repairs
	if d.WaterResistance == "IP68" {
		score -= 10 // Harder to open and reseal
	}

	// Add points based on successful repairs
	if d.GetRepairCount() > 0 {
		successfulRepairs := 0
		for _, repair := range d.Repairs {
			if repair.RepairStatus == "completed" {
				successfulRepairs++
			}
		}
		if successfulRepairs > 0 {
			score += math.Min(float64(successfulRepairs*5), 20) // Max 20 points
		}
	}

	// Brand-specific adjustments
	repairableBrands := map[string]float64{
		"Fairphone": 20,
		"Google":    10,
		"Samsung":   5,
	}

	if bonus, exists := repairableBrands[d.Brand]; exists {
		score += bonus
	}

	// Cap at 100
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

// CalculateLifecycleExtension calculates how much the device lifecycle has been extended
func (d *models.Device) CalculateLifecycleExtension() float64 {
	// Base expected lifespan
	expectedLifespan := 1095.0 // 3 years default

	switch d.DeviceSegment {
	case "flagship":
		expectedLifespan = 1460 // 4 years
	case "premium":
		expectedLifespan = 1277 // 3.5 years
	case "budget":
		expectedLifespan = 730 // 2 years
	}

	// Current age
	currentAge := float64(d.GetDeviceAge())

	// Extension percentage
	extension := ((currentAge - expectedLifespan) / expectedLifespan) * 100

	// Only positive extensions count
	if extension < 0 {
		extension = 0
	}

	return extension
}

// CalculateEWastePrevention calculates e-waste prevented in kg
func (d *models.Device) CalculateEWastePrevention() float64 {
	// Average smartphone weight
	deviceWeight := 0.2 // 200g average

	eWastePrevented := 0.0

	// Each repair prevents potential device replacement
	eWastePrevented += float64(d.GetRepairCount()) * deviceWeight

	// Lifecycle extension prevents waste
	extensionPercentage := d.CalculateLifecycleExtension()
	if extensionPercentage > 0 {
		eWastePrevented += deviceWeight * (extensionPercentage / 100)
	}

	// Refurbishment prevents waste
	if len(d.Refurbishments) > 0 {
		eWastePrevented += deviceWeight
	}

	// Trade-in for reuse prevents waste
	if len(d.TradeIns) > 0 {
		for _, tradeIn := range d.TradeIns {
			if tradeIn.TradeInStatus == "completed" {
				eWastePrevented += deviceWeight
			}
		}
	}

	return eWastePrevented
}

// CalculateResourceEfficiency calculates resource efficiency score
func (d *models.Device) CalculateResourceEfficiency() float64 {
	score := 50.0 // Base score

	// Points for longevity
	ageYears := float64(d.GetDeviceAge()) / 365
	if ageYears > 4 {
		score += 20
	} else if ageYears > 3 {
		score += 15
	} else if ageYears > 2 {
		score += 10
	} else if ageYears > 1 {
		score += 5
	}

	// Points for repairs over replacement
	score += math.Min(float64(d.GetRepairCount()*5), 15)

	// Points for good maintenance
	if d.BatteryHealth > 80 {
		score += 10
	}

	// Points for reuse/refurbishment
	if len(d.Refurbishments) > 0 {
		score += 15
	}

	// Deduct for premature issues
	if d.IsEndOfLife() && ageYears < 3 {
		score -= 20
	}

	// Cap at 100
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

// GetEnvironmentalGrade returns an environmental grade
func (d *models.Device) GetEnvironmentalGrade() string {
	// Calculate composite score
	recyclability := d.CalculateRecyclabilityScore()
	repairability := d.CalculateRepairabilityScore()
	efficiency := d.CalculateResourceEfficiency()

	compositeScore := (recyclability + repairability + efficiency) / 3

	if compositeScore >= 85 {
		return "A+"
	} else if compositeScore >= 75 {
		return "A"
	} else if compositeScore >= 65 {
		return "B"
	} else if compositeScore >= 55 {
		return "C"
	} else if compositeScore >= 45 {
		return "D"
	} else {
		return "F"
	}
}

// GetSustainabilityActions returns recommended sustainability actions
func (d *models.Device) GetSustainabilityActions() []string {
	actions := []string{}

	// Battery maintenance
	if d.BatteryHealth < 80 && d.BatteryHealth > 60 {
		actions = append(actions, "Consider battery replacement to extend device life")
	}

	// Screen repair
	if d.ScreenCondition == "cracked" {
		actions = append(actions, "Repair screen to prevent further damage")
	}

	// Refurbishment potential
	potential := d.GetRefurbishmentPotential()
	if potential == "excellent" || potential == "good" {
		actions = append(actions, "Device is a good candidate for refurbishment program")
	}

	// Trade-in suggestion
	if d.IsEligibleForTradeIn() && d.GetDeviceAge() > 730 {
		actions = append(actions, "Consider trade-in for device recycling/reuse")
	}

	// Recycling recommendation
	if d.IsEndOfLife() {
		actions = append(actions, "Device ready for responsible recycling")
	}

	// Maintenance reminder
	if d.LastInspection == nil || time.Since(*d.LastInspection) > 365*24*time.Hour {
		actions = append(actions, "Schedule device inspection for preventive maintenance")
	}

	// Upgrade efficiency
	if d.GetDeviceAge() > 1460 && d.Condition == "excellent" {
		actions = append(actions, "Consider donating device for reuse instead of recycling")
	}

	return actions
}
