package types

import (
	"database/sql/driver"
	"encoding/json"
)

// JSONArray custom type for JSON array storage in database
// This is used across multiple models for storing array data in JSON columns
type JSONArray []interface{}

// Value implements the driver.Valuer interface for database storage
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// JSONMap custom type for JSON object storage in database
// This is used for storing key-value pairs in JSON columns
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
