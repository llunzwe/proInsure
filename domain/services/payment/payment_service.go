package payment

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// paymentService implements the PaymentService interface
type paymentService struct {
	paymentRepo repositories.PaymentRepository
	// auditRepo   repositories.AuditRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo repositories.PaymentRepository,
	// auditRepo repositories.AuditRepository,
) services.PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		// auditRepo:   auditRepo,
	}
}

// ProcessPayment processes a payment transaction
func (s *paymentService) ProcessPayment(ctx context.Context, payment *models.Payment) error {
	if payment == nil {
		return errors.New("payment cannot be nil")
	}
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}
	if payment.Status == "" {
		payment.Status = "pending"
	}
	if payment.Reference == "" {
		payment.Reference = fmt.Sprintf("PAY-%s", uuid.New().String()[:12])
	}
	// Validate payment amount
	if payment.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}
	return s.paymentRepo.Create(ctx, payment)
}

// GetPaymentByID retrieves a payment by ID
func (s *paymentService) GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

// GetPaymentByReference retrieves a payment by reference
func (s *paymentService) GetPaymentByReference(ctx context.Context, reference string) (*models.Payment, error) {
	return s.paymentRepo.GetByReference(ctx, reference)
}

// UpdatePaymentStatus updates payment status
func (s *paymentService) UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status string, providerResponse map[string]interface{}) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return errors.New("payment not found")
	}
	payment.Status = status
	now := time.Now()
	payment.ProcessedAt = &now
	// Store provider response as JSON string (would need JSON marshaling in production)
	if providerResponse != nil {
		payment.ProviderResponse = fmt.Sprintf("%v", providerResponse)
	}
	return s.paymentRepo.Update(ctx, payment)
}

// GetPaymentsByUserID gets payments for a user
func (s *paymentService) GetPaymentsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error) {
	return s.paymentRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetPaymentsByPolicyID gets payments for a policy
func (s *paymentService) GetPaymentsByPolicyID(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error) {
	return s.paymentRepo.GetByPolicyID(ctx, policyID, limit, offset)
}

// GetPaymentsByStatus gets payments by status
func (s *paymentService) GetPaymentsByStatus(ctx context.Context, status string, limit, offset int) ([]*models.Payment, int64, error) {
	return s.paymentRepo.GetByStatus(ctx, status, limit, offset)
}

// GetPendingPayments gets pending payments
func (s *paymentService) GetPendingPayments(ctx context.Context, limit int) ([]*models.Payment, error) {
	return s.paymentRepo.GetPendingPayments(ctx, limit)
}

// CreatePaymentMethod creates a payment method
func (s *paymentService) CreatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error {
	if method == nil {
		return errors.New("payment method cannot be nil")
	}
	if method.ID == uuid.Nil {
		method.ID = uuid.New()
	}
	return s.paymentRepo.CreatePaymentMethod(ctx, method)
}

// GetPaymentMethodByID retrieves a payment method by ID
func (s *paymentService) GetPaymentMethodByID(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error) {
	return s.paymentRepo.GetPaymentMethodByID(ctx, id)
}

// GetPaymentMethodsByUserID gets payment methods for a user
func (s *paymentService) GetPaymentMethodsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PaymentMethod, error) {
	return s.paymentRepo.GetPaymentMethodsByUserID(ctx, userID)
}

// UpdatePaymentMethod updates a payment method
func (s *paymentService) UpdatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error {
	if method == nil {
		return errors.New("payment method cannot be nil")
	}
	return s.paymentRepo.UpdatePaymentMethod(ctx, method)
}

// DeletePaymentMethod deletes a payment method
func (s *paymentService) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	return s.paymentRepo.DeletePaymentMethod(ctx, id)
}

// SetDefaultPaymentMethod sets the default payment method for a user
func (s *paymentService) SetDefaultPaymentMethod(ctx context.Context, userID uuid.UUID, methodID uuid.UUID) error {
	return s.paymentRepo.SetDefaultPaymentMethod(ctx, userID, methodID)
}

// RefundPayment processes a payment refund
func (s *paymentService) RefundPayment(ctx context.Context, paymentID uuid.UUID, amount float64, reason string) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return errors.New("payment not found")
	}
	if payment.Status != "completed" {
		return errors.New("only completed payments can be refunded")
	}
	if amount > payment.Amount {
		return errors.New("refund amount cannot exceed payment amount")
	}
	payment.Status = "refunded"
	payment.RefundedAmount = amount
	payment.FailureReason = reason
	return s.paymentRepo.Update(ctx, payment)
}

// RetryFailedPayment retries a failed payment
func (s *paymentService) RetryFailedPayment(ctx context.Context, paymentID uuid.UUID) error {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return errors.New("payment not found")
	}
	if payment.Status != "failed" {
		return errors.New("only failed payments can be retried")
	}
	payment.Status = "pending"
	payment.RetryCount++
	payment.ProcessedAt = nil
	return s.paymentRepo.Update(ctx, payment)
}

// ReconcilePayment reconciles a payment
func (s *paymentService) ReconcilePayment(ctx context.Context, paymentID uuid.UUID, reconciledAt string) error {
	reconciledTime, err := time.Parse(time.RFC3339, reconciledAt)
	if err != nil {
		return fmt.Errorf("invalid reconciled date: %w", err)
	}
	return s.paymentRepo.MarkReconciled(ctx, paymentID, reconciledTime)
}

// SchedulePayment schedules a payment for later processing
func (s *paymentService) SchedulePayment(ctx context.Context, payment *models.Payment, scheduledFor string) error {
	if payment == nil {
		return errors.New("payment cannot be nil")
	}
	scheduledTime, err := time.Parse(time.RFC3339, scheduledFor)
	if err != nil {
		return fmt.Errorf("invalid scheduled date: %w", err)
	}
	payment.ScheduledFor = &scheduledTime
	payment.Status = "pending"
	return s.paymentRepo.Create(ctx, payment)
}

// GetScheduledPayments gets scheduled payments for a date
func (s *paymentService) GetScheduledPayments(ctx context.Context, scheduledFor string, limit int) ([]*models.Payment, error) {
	scheduledTime, err := time.Parse(time.RFC3339, scheduledFor)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled date: %w", err)
	}
	return s.paymentRepo.GetScheduledPayments(ctx, scheduledTime, limit)
}

// ProcessScheduledPayments processes all scheduled payments for a date
func (s *paymentService) ProcessScheduledPayments(ctx context.Context, scheduledFor string) error {
	payments, err := s.GetScheduledPayments(ctx, scheduledFor, 100)
	if err != nil {
		return fmt.Errorf("failed to get scheduled payments: %w", err)
	}
	for _, payment := range payments {
		payment.Status = "processing"
		if err := s.paymentRepo.Update(ctx, payment); err != nil {
			return fmt.Errorf("failed to process payment %s: %w", payment.ID, err)
		}
	}
	return nil
}

// GetPaymentStatistics gets payment statistics
func (s *paymentService) GetPaymentStatistics(ctx context.Context, userID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}
	return s.paymentRepo.GetPaymentStatistics(ctx, userID, start, end)
}

// GetTotalAmountByStatus gets total amount by status
func (s *paymentService) GetTotalAmountByStatus(ctx context.Context, status string, startDate, endDate string) (float64, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date: %w", err)
	}
	return s.paymentRepo.GetTotalAmountByStatus(ctx, status, start, end)
}
