package quote

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/policy"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// quoteService implements the QuoteService interface
type quoteService struct {
	quoteRepo  repositories.QuoteRepository
	pricingSvc services.PricingService
}

// NewQuoteService creates a new quote service
func NewQuoteService(
	quoteRepo repositories.QuoteRepository,
	pricingSvc services.PricingService,
) services.QuoteService {
	return &quoteService{
		quoteRepo:  quoteRepo,
		pricingSvc: pricingSvc,
	}
}

// CreateQuote creates a new quote
func (s *quoteService) CreateQuote(ctx context.Context, quote *policy.Quote) error {
	if quote == nil {
		return errors.New("quote cannot be nil")
	}
	if quote.ID == uuid.Nil {
		quote.ID = uuid.New()
	}
	if quote.QuoteNumber == "" {
		quote.QuoteNumber = fmt.Sprintf("QTE-%s", uuid.New().String()[:8])
	}
	if quote.Status == "" {
		quote.Status = "pending"
	}
	return s.quoteRepo.Create(ctx, quote)
}

// GetQuoteByID retrieves a quote by ID
func (s *quoteService) GetQuoteByID(ctx context.Context, id uuid.UUID) (*policy.Quote, error) {
	return s.quoteRepo.GetByID(ctx, id)
}

// GetQuoteByQuoteNumber retrieves a quote by quote number
func (s *quoteService) GetQuoteByQuoteNumber(ctx context.Context, quoteNumber string) (*policy.Quote, error) {
	return s.quoteRepo.GetByQuoteNumber(ctx, quoteNumber)
}

// UpdateQuote updates an existing quote
func (s *quoteService) UpdateQuote(ctx context.Context, quote *policy.Quote) error {
	if quote == nil {
		return errors.New("quote cannot be nil")
	}
	return s.quoteRepo.Update(ctx, quote)
}

// DeleteQuote deletes a quote
func (s *quoteService) DeleteQuote(ctx context.Context, id uuid.UUID) error {
	return s.quoteRepo.Delete(ctx, id)
}

// GetQuotesByUserID gets quotes for a user
func (s *quoteService) GetQuotesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*policy.Quote, int64, error) {
	return s.quoteRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetQuotesByDeviceID gets quotes for a device
func (s *quoteService) GetQuotesByDeviceID(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]*policy.Quote, int64, error) {
	return s.quoteRepo.GetByDeviceID(ctx, deviceID, limit, offset)
}

// GetQuotesByStatus gets quotes by status
func (s *quoteService) GetQuotesByStatus(ctx context.Context, status string, limit, offset int) ([]*policy.Quote, int64, error) {
	return s.quoteRepo.GetByStatus(ctx, status, limit, offset)
}

// GetValidQuotes gets valid (non-expired) quotes
func (s *quoteService) GetValidQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error) {
	return s.quoteRepo.GetValidQuotes(ctx, limit, offset)
}

// GetExpiredQuotes gets expired quotes
func (s *quoteService) GetExpiredQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error) {
	now := time.Now()
	return s.quoteRepo.GetExpiredQuotes(ctx, now, limit, offset)
}

// ConvertQuoteToPolicy converts a quote to a policy
func (s *quoteService) ConvertQuoteToPolicy(ctx context.Context, quoteID uuid.UUID) (uuid.UUID, error) {
	quote, err := s.quoteRepo.GetByID(ctx, quoteID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get quote: %w", err)
	}
	if quote == nil {
		return uuid.Nil, errors.New("quote not found")
	}
	if quote.Status != "approved" {
		return uuid.Nil, errors.New("only approved quotes can be converted to policies")
	}
	if time.Now().After(quote.ValidUntil) {
		return uuid.Nil, errors.New("quote has expired")
	}
	// Generate new policy ID (actual policy creation would be done by policy service)
	newPolicyID := uuid.New()
	now := time.Now()
	if err := s.quoteRepo.MarkAsConverted(ctx, quoteID, newPolicyID, now); err != nil {
		return uuid.Nil, fmt.Errorf("failed to mark quote as converted: %w", err)
	}
	return newPolicyID, nil
}

// GetUnconvertedQuotes gets quotes that haven't been converted
func (s *quoteService) GetUnconvertedQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error) {
	return s.quoteRepo.GetUnconvertedQuotes(ctx, limit, offset)
}

// CalculateQuote calculates a new quote based on parameters
func (s *quoteService) CalculateQuote(ctx context.Context, userID, deviceID uuid.UUID, coverageType string, coverageAmount float64) (*policy.Quote, error) {
	// Calculate premium using pricing service
	premium, err := s.pricingSvc.CalculatePremium(ctx, userID, deviceID, coverageType, coverageAmount, 12)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate premium: %w", err)
	}
	// Create quote with calculated values
	quote := &policy.Quote{
		ID:             uuid.New(),
		QuoteNumber:    fmt.Sprintf("QTE-%s", uuid.New().String()[:8]),
		UserID:         userID,
		DeviceID:       deviceID,
		Status:         "pending",
		Type:           "new",
		CoverageType:   coverageType,
		CoverageAmount: coverageAmount,
		BasePremium:    premium,
		TotalPremium:   premium,
		ValidFrom:      time.Now(),
		ValidUntil:     time.Now().AddDate(0, 0, 30), // 30 days validity
	}
	// Calculate tax (simplified - would use actual tax calculation)
	quote.TaxAmount = premium * 0.1
	quote.TotalPremium = premium + quote.TaxAmount
	return quote, nil
}

// RecalculateQuote recalculates an existing quote with adjustments
func (s *quoteService) RecalculateQuote(ctx context.Context, quoteID uuid.UUID, adjustments map[string]interface{}) (*policy.Quote, error) {
	quote, err := s.quoteRepo.GetByID(ctx, quoteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}
	if quote == nil {
		return nil, errors.New("quote not found")
	}
	// Apply adjustments
	if newCoverage, ok := adjustments["coverage_amount"].(float64); ok {
		quote.CoverageAmount = newCoverage
	}
	// Recalculate premium
	premium, err := s.pricingSvc.CalculatePremium(ctx, quote.UserID, quote.DeviceID, quote.CoverageType, quote.CoverageAmount, 12)
	if err != nil {
		return nil, fmt.Errorf("failed to recalculate premium: %w", err)
	}
	quote.BasePremium = premium
	quote.TaxAmount = premium * 0.1
	quote.TotalPremium = premium + quote.TaxAmount
	return quote, nil
}

// ValidateQuote validates a quote
func (s *quoteService) ValidateQuote(ctx context.Context, quoteID uuid.UUID) (bool, string, error) {
	quote, err := s.quoteRepo.GetByID(ctx, quoteID)
	if err != nil {
		return false, "", fmt.Errorf("failed to get quote: %w", err)
	}
	if quote == nil {
		return false, "Quote not found", nil
	}
	if quote.Status != "approved" {
		return false, "Quote is not approved", nil
	}
	if time.Now().After(quote.ValidUntil) {
		return false, "Quote has expired", nil
	}
	if quote.ConvertedAt != nil {
		return false, "Quote has already been converted", nil
	}
	return true, "Quote is valid", nil
}

// CheckQuoteValidity checks if quote is still valid
func (s *quoteService) CheckQuoteValidity(ctx context.Context, quoteID uuid.UUID) (bool, error) {
	valid, _, err := s.ValidateQuote(ctx, quoteID)
	return valid, err
}

// GetQuoteStatistics gets quote statistics
func (s *quoteService) GetQuoteStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}
	return s.quoteRepo.GetQuoteStatistics(ctx, userID, start, end)
}

// GetConversionRate gets quote conversion rate
func (s *quoteService) GetConversionRate(ctx context.Context, startDate, endDate string) (float64, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date: %w", err)
	}
	return s.quoteRepo.GetConversionRate(ctx, start, end)
}
