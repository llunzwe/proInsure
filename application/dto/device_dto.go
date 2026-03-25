// Package dto provides data transfer objects for application layer
package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DeviceDTO represents device data for API responses
type DeviceDTO struct {
	ID          uuid.UUID       `json:"id"`
	IMEI        string          `json:"imei"`
	SerialNumber string         `json:"serial_number,omitempty"`
	Brand       string          `json:"brand"`
	Model       string          `json:"model"`
	Category    string          `json:"category"`
	Status      string          `json:"status"`
	Condition   string          `json:"condition"`
	CurrentValue decimal.Decimal `json:"current_value"`
	MarketValue  decimal.Decimal `json:"market_value,omitempty"`
	PurchasePrice decimal.Decimal `json:"purchase_price,omitempty"`
	PurchaseDate *time.Time      `json:"purchase_date,omitempty"`
	OwnerID     uuid.UUID       `json:"owner_id"`
	IsInsured   bool            `json:"is_insured"`
	IsStolen    bool            `json:"is_stolen"`
	IsBlacklisted bool          `json:"is_blacklisted"`
	RiskScore   float64         `json:"risk_score"`
	WarrantyEndDate *time.Time  `json:"warranty_end_date,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// DeviceCreateRequest represents a request to create a device
type DeviceCreateRequest struct {
	IMEI           string          `json:"imei" binding:"required,len=15,numeric"`
	SerialNumber   string          `json:"serial_number,omitempty"`
	Brand          string          `json:"brand" binding:"required,max=100"`
	Model          string          `json:"model" binding:"required,max=100"`
	Category       string          `json:"category" binding:"required,oneof=smartphone smartwatch tablet laptop wearable"`
	PurchasePrice  decimal.Decimal `json:"purchase_price" binding:"required,gt=0"`
	PurchaseDate   *time.Time      `json:"purchase_date,omitempty"`
	WarrantyMonths int             `json:"warranty_months,omitempty"`
	Condition      string          `json:"condition" binding:"omitempty,oneof=excellent good fair poor"`
	OwnerID        uuid.UUID       `json:"owner_id" binding:"required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// DeviceUpdateRequest represents a request to update a device
type DeviceUpdateRequest struct {
	SerialNumber   *string          `json:"serial_number,omitempty"`
	Condition      *string          `json:"condition,omitempty" binding:"omitempty,oneof=excellent good fair poor"`
	CurrentValue   *decimal.Decimal `json:"current_value,omitempty"`
	Status         *string          `json:"status,omitempty" binding:"omitempty,oneof=active inactive lost stolen damaged"`
	WarrantyEndDate *time.Time      `json:"warranty_end_date,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// DeviceSearchRequest represents search parameters for devices
type DeviceSearchRequest struct {
	Brand           string    `form:"brand"`
	Model           string    `form:"model"`
	Category        string    `form:"category"`
	Status          string    `form:"status"`
	OwnerID         uuid.UUID `form:"owner_id"`
	IsInsured       *bool     `form:"is_insured"`
	IsStolen        *bool     `form:"is_stolen"`
	MinValue        decimal.Decimal `form:"min_value"`
	MaxValue        decimal.Decimal `form:"max_value"`
	MinRiskScore    float64   `form:"min_risk_score"`
	MaxRiskScore    float64   `form:"max_risk_score"`
	PurchasedAfter  time.Time `form:"purchased_after"`
	PurchasedBefore time.Time `form:"purchased_before"`
	SortBy          string    `form:"sort_by" binding:"omitempty,oneof=created_at updated_at brand model current_value risk_score"`
	SortOrder       string    `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page            int       `form:"page" binding:"omitempty,min=1"`
	PageSize        int       `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// DeviceInsuranceRequest represents a request to insure a device
type DeviceInsuranceRequest struct {
	DeviceID        uuid.UUID       `json:"device_id" binding:"required"`
	UserID          uuid.UUID       `json:"user_id" binding:"required"`
	CoverageType    string          `json:"coverage_type" binding:"required,oneof=basic standard premium"`
	Deductible      decimal.Decimal `json:"deductible" binding:"required,gt=0"`
	PaymentMethod   string          `json:"payment_method" binding:"required,oneof=credit_card debit_card bank_transfer"`
	BillingCycle    string          `json:"billing_cycle" binding:"required,oneof=monthly quarterly yearly"`
}

// DeviceInsuranceResponse represents the response for insuring a device
type DeviceInsuranceResponse struct {
	PolicyID       uuid.UUID       `json:"policy_id"`
	DeviceID       uuid.UUID       `json:"device_id"`
	PolicyNumber   string          `json:"policy_number"`
	Premium        decimal.Decimal `json:"premium"`
	CoverageAmount decimal.Decimal `json:"coverage_amount"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	Status         string          `json:"status"`
	Message        string          `json:"message"`
}

// DeviceValuationRequest represents a request to value a device
type DeviceValuationRequest struct {
	DeviceID      uuid.UUID `json:"device_id" binding:"required"`
	Condition     string    `json:"condition" binding:"required,oneof=excellent good fair poor"`
	ScreenCondition string  `json:"screen_condition" binding:"omitempty,oneof=perfect minor_scratches scratched cracked"`
	BodyCondition string    `json:"body_condition" binding:"omitempty,oneof=perfect minor_wear scratched dented damaged"`
	BatteryHealth *int      `json:"battery_health,omitempty" binding:"omitempty,min=0,max=100"`
}

// DeviceValuationResponse represents a device valuation result
type DeviceValuationResponse struct {
	DeviceID         uuid.UUID       `json:"device_id"`
	EstimatedValue   decimal.Decimal `json:"estimated_value"`
	DepreciationRate float64         `json:"depreciation_rate"`
	MarketTrend      string          `json:"market_trend"`
	ConfidenceScore  float64         `json:"confidence_score"`
	ValuationDate    time.Time       `json:"valuation_date"`
	Factors          []ValuationFactor `json:"factors"`
}

// ValuationFactor represents a factor in device valuation
type ValuationFactor struct {
	Factor    string  `json:"factor"`
	Impact    string  `json:"impact"` // positive, negative, neutral
	Weight    float64 `json:"weight"`
	Description string `json:"description"`
}

// DeviceListResponse represents a paginated list of devices
type DeviceListResponse struct {
	Devices    []DeviceDTO `json:"devices"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// DeviceTransferRequest represents a request to transfer device ownership
type DeviceTransferRequest struct {
	DeviceID      uuid.UUID `json:"device_id" binding:"required"`
	CurrentOwnerID uuid.UUID `json:"current_owner_id" binding:"required"`
	NewOwnerID    uuid.UUID `json:"new_owner_id" binding:"required"`
	TransferReason string   `json:"transfer_reason,omitempty"`
}

// DeviceEligibilityResponse represents device insurance eligibility
type DeviceEligibilityResponse struct {
	DeviceID    uuid.UUID `json:"device_id"`
	Eligible    bool      `json:"eligible"`
	Reason      string    `json:"reason,omitempty"`
	Suggestions []string  `json:"suggestions,omitempty"`
	Conditions  []string  `json:"conditions,omitempty"`
}
