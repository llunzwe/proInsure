package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	// "smartsure/internal/domain/models/policy" // Unused import
)

// PolicyQuoteRequest represents a request for policy quote
type PolicyQuoteRequest struct {
	CustomerID   uuid.UUID `json:"customer_id"`
	DeviceIMEI   string    `json:"device_imei"`
	CoverageType string    `json:"coverage_type"`
	CoverageAmount float64  `json:"coverage_amount"`
	Deductible   float64    `json:"deductible"`
	PolicyType   string    `json:"policy_type"`
}

// PolicyQuoteResponse represents the response for policy quote
type PolicyQuoteResponse struct {
	PolicyNumber  string    `json:"policy_number"`
	PremiumAmount float64   `json:"premium_amount"`
	CoverageAmount float64  `json:"coverage_amount"`
	Deductible    float64    `json:"deductible"`
	ValidUntil    time.Time `json:"valid_until"`
	QuoteID       uuid.UUID `json:"quote_id"`
}

// PolicyService defines the interface for policy business logic
type PolicyService interface {
	// Lifecycle Management
	CreatePolicy(ctx context.Context, policy *models.Policy) error
	UpdatePolicy(ctx context.Context, policy *models.Policy) error
	GetPolicyByID(ctx context.Context, id uuid.UUID) (*models.Policy, error)
	GetPolicyByNumber(ctx context.Context, policyNumber string) (*models.Policy, error)

	// Policy Status Operations
	ActivatePolicy(ctx context.Context, policyID uuid.UUID) error
	SuspendPolicy(ctx context.Context, policyID uuid.UUID, reason string) error
	CancelPolicy(ctx context.Context, policyID uuid.UUID, reason string, effectiveDate time.Time) error
	RenewPolicy(ctx context.Context, policyID uuid.UUID) (*models.Policy, error)

	// Validation & Eligibility
	ValidatePolicyForClaim(ctx context.Context, policyID uuid.UUID) (bool, string, error)
	ValidatePolicyForRenewal(ctx context.Context, policyID uuid.UUID) (bool, string, error)
	ValidatePolicyForModification(ctx context.Context, policyID uuid.UUID) (bool, error)
	CheckPolicyActive(ctx context.Context, policyID uuid.UUID) (bool, error)
	CheckPolicyExpired(ctx context.Context, policyID uuid.UUID) (bool, error)

	// Premium Calculations
	CalculatePremium(ctx context.Context, policy *models.Policy) (float64, error)
	CalculateTotalPremium(ctx context.Context, policy *models.Policy) (float64, error)
	RecalculatePremiumAfterClaim(ctx context.Context, policyID uuid.UUID, claimAmount float64) error
	ApplyDiscount(ctx context.Context, policyID uuid.UUID, discountCode string) error

	// Coverage Management
	UpdateCoverage(ctx context.Context, policyID uuid.UUID, newCoverage float64) error
	CheckCoverageLimits(ctx context.Context, policyID uuid.UUID, claimAmount float64) (bool, error)
	GetRemainingCoverage(ctx context.Context, policyID uuid.UUID) (float64, error)

	// Payment Management
	ProcessPayment(ctx context.Context, policyID uuid.UUID, amount float64) error
	CheckPaymentStatus(ctx context.Context, policyID uuid.UUID) (string, error)
	GetPaymentHistory(ctx context.Context, policyID uuid.UUID) ([]models.Payment, error)

	// Risk Assessment
	AssessRisk(ctx context.Context, policyID uuid.UUID) (float64, error)
	UpdateRiskScore(ctx context.Context, policyID uuid.UUID, newScore float64) error

	// Reporting
	GetPolicyMetrics(ctx context.Context, policyID uuid.UUID) (map[string]interface{}, error)
	GetExpiringPolicies(ctx context.Context, days int) ([]models.Policy, error)
	GetPoliciesByStatus(ctx context.Context, status string) ([]models.Policy, error)

	// ======================================
	// MISSING IMPLEMENTED METHODS
	// ======================================

	// Analytics & History
	GetPolicyAnalytics(ctx context.Context, policyID uuid.UUID) (map[string]interface{}, error)
	GetPolicyClaimsHistory(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)
	GetPolicyPaymentHistory(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)

	// Documents
	GetPolicyDocuments(ctx context.Context, policyID uuid.UUID) ([]map[string]interface{}, error)
	UploadPolicyDocument(ctx context.Context, policyID uuid.UUID, documentType, fileName, filePath string, fileSize int64, mimeType string, uploadedBy uuid.UUID, description string) (*models.Document, error)
	GetPolicyCertificate(ctx context.Context, policyID uuid.UUID) (map[string]interface{}, error)

	// Riders & Add-ons
	AddPolicyRider(ctx context.Context, policyID uuid.UUID, riderType, description string, additionalPremium, coverageAmount float64, effectiveDate, expiryDate, terms string) (interface{}, error) // PolicyRider type not found
	GetPolicyRiders(ctx context.Context, policyID uuid.UUID) ([]interface{}, error) // PolicyRider type not found
	UpdatePolicyRider(ctx context.Context, policyID, riderID uuid.UUID, description string, additionalPremium, coverageAmount float64, expiryDate, terms, status string) error
	RemovePolicyRider(ctx context.Context, policyID, riderID uuid.UUID) error

	// Policy Management
	ComparePolicies(ctx context.Context, policyIDs []uuid.UUID, compareFields []string) (map[string]interface{}, error)
	GetPolicyRecommendations(ctx context.Context, userID uuid.UUID, deviceID *uuid.UUID, riskTolerance, budget string) ([]map[string]interface{}, error)
	PortPolicy(ctx context.Context, policyID uuid.UUID, targetInsurer, reason, requestedDate, notes string) (interface{}, error) // PolicyPortRequest type not found
	GetPolicyPortStatus(ctx context.Context, portID uuid.UUID) (map[string]interface{}, error)

	// Legacy methods for backward compatibility
	GetPolicyQuote(ctx context.Context, req *PolicyQuoteRequest) (*PolicyQuoteResponse, error)
	GetPoliciesByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.Policy, error)
	GetActivePolicies(ctx context.Context) ([]models.Policy, error)
	GetUserPolicies(ctx context.Context, userID uuid.UUID) ([]*models.Policy, error)
	UpdatePolicyByID(ctx context.Context, policyID uuid.UUID, updates map[string]interface{}) (*models.Policy, error) // Renamed to avoid duplicate
	DeletePolicy(ctx context.Context, policyID uuid.UUID) error
}
