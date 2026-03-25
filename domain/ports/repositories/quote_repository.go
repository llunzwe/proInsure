package repositories

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models/policy"
)

// QuoteRepository defines the interface for quote persistence operations
type QuoteRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, quote *policy.Quote) error
	GetByID(ctx context.Context, id uuid.UUID) (*policy.Quote, error)
	GetByQuoteNumber(ctx context.Context, quoteNumber string) (*policy.Quote, error)
	Update(ctx context.Context, quote *policy.Quote) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*policy.Quote, int64, error)
	GetByDeviceID(ctx context.Context, deviceID uuid.UUID, limit, offset int) ([]*policy.Quote, int64, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*policy.Quote, int64, error)
	GetByType(ctx context.Context, quoteType string, limit, offset int) ([]*policy.Quote, int64, error)
	GetExpiredQuotes(ctx context.Context, beforeDate time.Time, limit, offset int) ([]*policy.Quote, int64, error)
	GetValidQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error)
	Search(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*policy.Quote, int64, error)

	// Conversion operations
	GetUnconvertedQuotes(ctx context.Context, limit, offset int) ([]*policy.Quote, int64, error)
	MarkAsConverted(ctx context.Context, quoteID uuid.UUID, policyID uuid.UUID, convertedAt time.Time) error

	// Quote items operations
	CreateQuoteItem(ctx context.Context, item *policy.QuoteItem) error
	GetQuoteItemsByQuoteID(ctx context.Context, quoteID uuid.UUID) ([]*policy.QuoteItem, error)
	UpdateQuoteItem(ctx context.Context, item *policy.QuoteItem) error
	DeleteQuoteItem(ctx context.Context, id uuid.UUID) error

	// Quote history operations
	CreateQuoteHistory(ctx context.Context, history *policy.QuoteHistory) error
	GetQuoteHistoryByQuoteID(ctx context.Context, quoteID uuid.UUID) ([]*policy.QuoteHistory, error)

	// Statistics
	GetQuoteStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetConversionRate(ctx context.Context, startDate, endDate time.Time) (float64, error)
}
