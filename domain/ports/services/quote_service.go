package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/policy"
)

// QuoteService defines the interface for quote business logic operations
type QuoteService interface {
	// Quote creation and management
	CreateQuote(ctx context.Context, quote *policy.Quote) error
	GetQuoteByID(ctx context.Context, id uuid.UUID) (*policy.Quote, error)
	GetQuoteByQuoteNumber(ctx context.Context, quoteNumber string) (*policy.Quote, error)
	UpdateQuote(ctx context.Context, quote *policy.Quote) error
	DeleteQuote(ctx context.Context, id uuid.UUID) error

	// Quote queries
	GetQuotesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*policy.Quote, int64, error)
	GetQuotesByDeviceID(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]*policy.Quote, int64, error)
	GetQuotesByStatus(ctx context.Context, status string, limit, offset int) ([]*policy.Quote, int64, error)
	GetValidQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error)
	GetExpiredQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error)

	// Quote conversion
	ConvertQuoteToPolicy(ctx context.Context, quoteID uuid.UUID) (uuid.UUID, error)
	GetUnconvertedQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error)

	// Quote calculation
	CalculateQuote(ctx context.Context, userID, deviceID uuid.UUID, coverageType string, coverageAmount float64) (*policy.Quote, error)
	RecalculateQuote(ctx context.Context, quoteID uuid.UUID, adjustments map[string]interface{}) (*policy.Quote, error)

	// Quote validation
	ValidateQuote(ctx context.Context, quoteID uuid.UUID) (bool, string, error)
	CheckQuoteValidity(ctx context.Context, quoteID uuid.UUID) (bool, error)

	// Quote statistics
	GetQuoteStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetConversionRate(ctx context.Context, startDate, endDate string) (float64, error)
}
