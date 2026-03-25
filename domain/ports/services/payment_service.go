package services

import (
	"context"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
)

// PaymentService defines the interface for payment business logic operations
type PaymentService interface {
	// Payment processing
	ProcessPayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	GetPaymentByReference(ctx context.Context, reference string) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status string, providerResponse map[string]interface{}) error

	// Payment queries
	GetPaymentsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error)
	GetPaymentsByPolicyID(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error)
	GetPaymentsByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Payment, int64, error)
	GetPendingPayments(ctx context.Context, limit int) ([]*models.Payment, error)

	// Payment method management
	CreatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error
	GetPaymentMethodByID(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error)
	GetPaymentMethodsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, id uuid.UUID) error
	SetDefaultPaymentMethod(ctx context.Context, userID uuid.UUID, methodID uuid.UUID) error

	// Payment operations
	RefundPayment(ctx context.Context, paymentID uuid.UUID, amount float64, reason string) error
	RetryFailedPayment(ctx context.Context, paymentID uuid.UUID) error
	ReconcilePayment(ctx context.Context, paymentID uuid.UUID, reconciledAt string) error

	// Payment scheduling
	SchedulePayment(ctx context.Context, payment *models.Payment, scheduledFor string) error
	GetScheduledPayments(ctx context.Context, scheduledFor string, limit int) ([]*models.Payment, error)
	ProcessScheduledPayments(ctx context.Context, scheduledFor string) error

	// Payment statistics
	GetPaymentStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetTotalAmountByStatus(ctx context.Context, status string, startDate, endDate string) (float64, error)
}
