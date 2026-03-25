package models

import (
	"reflect"
)

// Helper functions to access policy fields using reflection
// This avoids importing the models package and breaking the import cycle

// getFieldValue gets a field value from the policy using reflection
func (pc *PolicyCalculator) getFieldValue(fieldPath ...string) reflect.Value {
	val := reflect.ValueOf(pc.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, field := range fieldPath {
		val = val.FieldByName(field)
		if !val.IsValid() {
			return reflect.Value{}
		}
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return reflect.Value{}
			}
			val = val.Elem()
		}
	}
	return val
}

// getFloat64Field gets a float64 field value
func (pc *PolicyCalculator) getFloat64Field(fieldPath ...string) float64 {
	val := pc.getFieldValue(fieldPath...)
	if !val.IsValid() || val.Kind() != reflect.Float64 {
		return 0.0
	}
	return val.Float()
}

// getStringField gets a string field value
func (pc *PolicyCalculator) getStringField(fieldPath ...string) string {
	val := pc.getFieldValue(fieldPath...)
	if !val.IsValid() || val.Kind() != reflect.String {
		return ""
	}
	return val.String()
}

// getIntField gets an int field value
func (pc *PolicyCalculator) getIntField(fieldPath ...string) int {
	val := pc.getFieldValue(fieldPath...)
	if !val.IsValid() {
		return 0
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(val.Uint())
	}
	return 0
}

// Helper to get embedded struct field (e.g., PolicyPricing.Currency)
func (pc *PolicyCalculator) getEmbeddedField(embeddedStruct, field string) reflect.Value {
	val := reflect.ValueOf(pc.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	embeddedVal := val.FieldByName(embeddedStruct)
	if !embeddedVal.IsValid() {
		return reflect.Value{}
	}

	fieldVal := embeddedVal.FieldByName(field)
	return fieldVal
}

// getPolicyPricingCurrency gets currency from PolicyPricing
func (pc *PolicyCalculator) getPolicyPricingCurrency() string {
	val := pc.getEmbeddedField("PolicyPricing", "Currency")
	if !val.IsValid() {
		return "USD"
	}
	if val.Kind() == reflect.String {
		return val.String()
	}
	// If Currency is a custom type, get its underlying string value
	if val.Kind() == reflect.Uint8 || val.Kind() == reflect.Uint16 || val.Kind() == reflect.Uint32 {
		// Try to get string representation
		return val.String()
	}
	return "USD"
}

// getPolicyCoverageAmount gets coverage amount from PolicyCoverageDetails
func (pc *PolicyCalculator) getPolicyCoverageAmount() float64 {
	// Access PolicyCoverageDetails.CoverageAmount.Amount
	coverageDetails := pc.getFieldValue("PolicyCoverageDetails")
	if !coverageDetails.IsValid() {
		return 1000.0
	}

	coverageAmount := coverageDetails.FieldByName("CoverageAmount")
	if !coverageAmount.IsValid() {
		return 1000.0
	}

	amount := coverageAmount.FieldByName("Amount")
	if !amount.IsValid() || amount.Kind() != reflect.Float64 {
		return 1000.0
	}
	return amount.Float()
}

// getPolicyBasePremium gets base premium from PolicyPricing
func (pc *PolicyCalculator) getPolicyBasePremium() float64 {
	pricing := pc.getFieldValue("PolicyPricing")
	if !pricing.IsValid() {
		return 0.0
	}

	basePremium := pricing.FieldByName("BasePremium")
	if !basePremium.IsValid() {
		return 0.0
	}

	amount := basePremium.FieldByName("Amount")
	if !amount.IsValid() || amount.Kind() != reflect.Float64 {
		return 0.0
	}
	return amount.Float()
}

// getPolicyRiskAdjustedPremium gets risk adjusted premium
func (pc *PolicyCalculator) getPolicyRiskAdjustedPremium() float64 {
	pricing := pc.getFieldValue("PolicyPricing")
	if !pricing.IsValid() {
		return 0.0
	}

	riskPremium := pricing.FieldByName("RiskAdjustedPremium")
	if !riskPremium.IsValid() {
		return 0.0
	}

	amount := riskPremium.FieldByName("Amount")
	if !amount.IsValid() || amount.Kind() != reflect.Float64 {
		return 0.0
	}
	return amount.Float()
}
