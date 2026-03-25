package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
)

// PaymentRepository defines the interface for payment persistence operations
type PaymentRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, payment *models.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	GetByReference(ctx context.Context, reference string) (*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error)
	GetByPolicyID(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error)
	GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*models.Payment, error)
	GetByType(ctx context.Context, paymentType string, limit, offset int) ([]*models.Payment, int64, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Payment, int64, error)
	GetByProvider(ctx context.Context, provider string, limit, offset int) ([]*models.Payment, int64, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*models.Payment, int64, error)
	Search(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*models.Payment, int64, error)

	// Payment method operations
	CreatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error
	GetPaymentMethodByID(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error)
	GetPaymentMethodsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, id uuid.UUID) error
	SetDefaultPaymentMethod(ctx context.Context, userID uuid.UUID, methodID uuid.UUID) error

	// Transaction operations
	GetPendingPayments(ctx context.Context, limit int) ([]*models.Payment, error)
	GetScheduledPayments(ctx context.Context, scheduledFor time.Time, limit int) ([]*models.Payment, error)
	UpdateProviderTransactionID(ctx context.Context, paymentID uuid.UUID, providerTxnID string) error

	// Reconciliation
	GetUnreconciledPayments(ctx context.Context, limit, offset int) ([]*models.Payment, int64, error)
	MarkReconciled(ctx context.Context, paymentID uuid.UUID, reconciledAt time.Time) error

	// Statistics
	GetPaymentStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetTotalAmountByStatus(ctx context.Context, status string, startDate, endDate time.Time) (float64, error)
}
