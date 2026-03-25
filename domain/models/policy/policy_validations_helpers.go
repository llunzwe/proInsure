package policy

import (
	"reflect"

	"github.com/google/uuid"
)

// Helper methods for PolicyValidator to access policy fields using reflection
// This avoids importing the models package and breaking the import cycle

// getPolicyField gets a field value from the policy using reflection
func (pv *PolicyValidator) getPolicyField(fieldPath ...string) reflect.Value {
	val := reflect.ValueOf(pv.policy)
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

// getStringField gets a string field value
func (pv *PolicyValidator) getStringField(fieldPath ...string) string {
	val := pv.getPolicyField(fieldPath...)
	if !val.IsValid() || val.Kind() != reflect.String {
		return ""
	}
	return val.String()
}

// getFloat64Field gets a float64 field value
func (pv *PolicyValidator) getFloat64Field(fieldPath ...string) float64 {
	val := pv.getPolicyField(fieldPath...)
	if !val.IsValid() {
		return 0.0
	}
	switch val.Kind() {
	case reflect.Float64:
		return val.Float()
	case reflect.Float32:
		return float64(val.Float())
	}
	return 0.0
}

// getUUIDField gets a UUID field value
func (pv *PolicyValidator) getUUIDField(fieldPath ...string) uuid.UUID {
	val := pv.getPolicyField(fieldPath...)
	if !val.IsValid() {
		return uuid.Nil
	}
	if val.Kind() == reflect.Array && val.Type().String() == "[16]uint8" {
		// UUID is represented as [16]byte
		var id uuid.UUID
		reflect.Copy(reflect.ValueOf(id[:]), val)
		return id
	}
	// Try to convert to UUID
	if val.CanInterface() {
		if id, ok := val.Interface().(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

// getBoolField gets a bool field value
func (pv *PolicyValidator) getBoolField(fieldPath ...string) bool {
	val := pv.getPolicyField(fieldPath...)
	if !val.IsValid() || val.Kind() != reflect.Bool {
		return false
	}
	return val.Bool()
}

// getIntField gets an int field value
func (pv *PolicyValidator) getIntField(fieldPath ...string) int {
	val := pv.getPolicyField(fieldPath...)
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

// getNestedTimeField gets a nested time.Time field
func (pv *PolicyValidator) getNestedTimeField(embeddedStruct, field string) reflect.Value {
	val := reflect.ValueOf(pv.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	embeddedVal := val.FieldByName(embeddedStruct)
	if !embeddedVal.IsValid() {
		return reflect.Value{}
	}

	fieldVal := embeddedVal.FieldByName(field)
	if !fieldVal.IsValid() {
		return reflect.Value{}
	}

	return fieldVal
}

// getNestedStringField gets a nested string field
func (pv *PolicyValidator) getNestedStringField(embeddedStruct, field string) string {
	val := reflect.ValueOf(pv.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	embeddedVal := val.FieldByName(embeddedStruct)
	if !embeddedVal.IsValid() {
		return ""
	}

	fieldVal := embeddedVal.FieldByName(field)
	if !fieldVal.IsValid() || fieldVal.Kind() != reflect.String {
		return ""
	}
	return fieldVal.String()
}

// getNestedUUIDField gets a nested UUID field
func (pv *PolicyValidator) getNestedUUIDField(embeddedStruct, field string) uuid.UUID {
	val := reflect.ValueOf(pv.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	embeddedVal := val.FieldByName(embeddedStruct)
	if !embeddedVal.IsValid() {
		return uuid.Nil
	}

	fieldVal := embeddedVal.FieldByName(field)
	if !fieldVal.IsValid() {
		return uuid.Nil
	}

	if fieldVal.CanInterface() {
		if id, ok := fieldVal.Interface().(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

// getNestedBoolField gets a nested bool field
func (pv *PolicyValidator) getNestedBoolField(embeddedStruct, field string) bool {
	val := reflect.ValueOf(pv.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	embeddedVal := val.FieldByName(embeddedStruct)
	if !embeddedVal.IsValid() {
		return false
	}

	fieldVal := embeddedVal.FieldByName(field)
	if !fieldVal.IsValid() || fieldVal.Kind() != reflect.Bool {
		return false
	}
	return fieldVal.Bool()
}

// getNestedIntField gets a nested int field
func (pv *PolicyValidator) getNestedIntField(embeddedStruct, field string) int {
	val := reflect.ValueOf(pv.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	embeddedVal := val.FieldByName(embeddedStruct)
	if !embeddedVal.IsValid() {
		return 0
	}

	fieldVal := embeddedVal.FieldByName(field)
	if !fieldVal.IsValid() {
		return 0
	}
	switch fieldVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(fieldVal.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(fieldVal.Uint())
	}
	return 0
}

// getNestedFloat64Field gets a nested float64 field (e.g., PolicyCoverageDetails.CoinsurancePercent)
// If structField is empty, it accesses the field directly from embeddedStruct
func (pv *PolicyValidator) getNestedFloat64Field(embeddedStruct, structField, field string) float64 {
	val := reflect.ValueOf(pv.policy)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	embeddedVal := val.FieldByName(embeddedStruct)
	if !embeddedVal.IsValid() {
		return 0.0
	}

	var fieldVal reflect.Value
	if structField == "" {
		// Direct field access
		fieldVal = embeddedVal.FieldByName(field)
	} else {
		// Nested field access
		nestedVal := embeddedVal.FieldByName(structField)
		if !nestedVal.IsValid() {
			return 0.0
		}
		fieldVal = nestedVal.FieldByName(field)
	}

	if !fieldVal.IsValid() {
		return 0.0
	}

	switch fieldVal.Kind() {
	case reflect.Float64:
		return fieldVal.Float()
	case reflect.Float32:
		return float64(fieldVal.Float())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(fieldVal.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(fieldVal.Uint())
	}
	return 0.0
}
