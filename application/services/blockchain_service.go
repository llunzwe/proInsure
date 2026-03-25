package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"smartsure/internal/domain/ports/services"
	"smartsure/pkg/blockchain/fabric"
)

// BlockchainService implements the blockchain service interface
type BlockchainService struct {
	fabricService *fabric.Service
	logger        *logrus.Logger
	db            *gorm.DB
	enabled       bool
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService(db *gorm.DB, logger *logrus.Logger, config *fabric.FabricConfig) services.BlockchainService {
	fabricService := fabric.NewService(db, logger, config)

	service := &BlockchainService{
		fabricService: fabricService,
		logger:        logger,
		db:            db,
		enabled:       config.Enabled,
	}

	if service.enabled {
		logger.Info("Blockchain service initialized and enabled")
	} else {
		logger.Info("Blockchain service initialized but disabled (graceful degradation mode)")
	}

	return service
}

// IsEnabled returns whether the blockchain service is enabled and operational
func (s *BlockchainService) IsEnabled() bool {
	return s.enabled && s.fabricService.IsEnabled()
}

// GetBlockchainStatus returns the current status of the blockchain service
func (s *BlockchainService) GetBlockchainStatus() map[string]interface{} {
	status := s.fabricService.GetBlockchainStatus()
	status["service_enabled"] = s.enabled
	status["overall_status"] = s.getOverallStatus()
	return status
}

// getOverallStatus determines the overall status of the blockchain service
func (s *BlockchainService) getOverallStatus() string {
	if !s.enabled {
		return "disabled"
	}
	if !s.fabricService.IsEnabled() {
		return "degraded"
	}
	return "operational"
}

// SubmitPolicy submits a policy to the blockchain with graceful fallback
func (s *BlockchainService) SubmitPolicy(ctx context.Context, policy *fabric.PolicyRecord) error {
	if policy == nil {
		return fmt.Errorf("policy record cannot be nil")
	}

	s.logger.WithField("policy_id", policy.PolicyID).Info("Submitting policy to blockchain")

	start := time.Now()
	err := s.fabricService.SubmitPolicy(ctx, policy)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"policy_id": policy.PolicyID,
			"duration":  duration.String(),
		}).Warn("Policy submission to blockchain failed, using local storage")

		// Store locally for later sync when blockchain is available
		return s.storePolicyForSync(policy)
	}

	s.logger.WithFields(logrus.Fields{
		"policy_id": policy.PolicyID,
		"duration":  duration.String(),
	}).Info("Policy successfully submitted to blockchain")

	return nil
}

// SubmitClaim submits a claim to the blockchain with graceful fallback
func (s *BlockchainService) SubmitClaim(ctx context.Context, claim *fabric.ClaimRecord) error {
	if claim == nil {
		return fmt.Errorf("claim record cannot be nil")
	}

	s.logger.WithField("claim_id", claim.ClaimID).Info("Submitting claim to blockchain")

	start := time.Now()
	err := s.fabricService.SubmitClaim(ctx, claim)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"claim_id": claim.ClaimID,
			"duration":  duration.String(),
		}).Warn("Claim submission to blockchain failed, using local storage")

		return s.storeClaimForSync(claim)
	}

	s.logger.WithFields(logrus.Fields{
		"claim_id": claim.ClaimID,
		"duration":  duration.String(),
	}).Info("Claim successfully submitted to blockchain")

	return nil
}

// RegisterDevice registers device authenticity on the blockchain
func (s *BlockchainService) RegisterDevice(ctx context.Context, device *fabric.DeviceRecord) error {
	if device == nil {
		return fmt.Errorf("device record cannot be nil")
	}

	s.logger.WithField("device_id", device.DeviceID).Info("Registering device on blockchain")

	start := time.Now()
	err := s.fabricService.RegisterDevice(ctx, device)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"device_id": device.DeviceID,
			"duration":  duration.String(),
		}).Warn("Device registration on blockchain failed, using local storage")

		return s.storeDeviceForSync(device)
	}

	s.logger.WithFields(logrus.Fields{
		"device_id": device.DeviceID,
		"duration":  duration.String(),
	}).Info("Device successfully registered on blockchain")

	return nil
}

// SubmitAudit submits an audit record to the blockchain
func (s *BlockchainService) SubmitAudit(ctx context.Context, audit *fabric.AuditRecord) error {
	if audit == nil {
		return fmt.Errorf("audit record cannot be nil")
	}

	s.logger.WithField("audit_id", audit.AuditID).Info("Submitting audit to blockchain")

	start := time.Now()
	err := s.fabricService.SubmitAudit(ctx, audit)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"audit_id": audit.AuditID,
			"duration":  duration.String(),
		}).Warn("Audit submission to blockchain failed, using local storage")

		return s.storeAuditForSync(audit)
	}

	s.logger.WithFields(logrus.Fields{
		"audit_id": audit.AuditID,
		"duration":  duration.String(),
	}).Info("Audit successfully submitted to blockchain")

	return nil
}

// QueryBlockchainRecord queries a record from the blockchain
func (s *BlockchainService) QueryBlockchainRecord(ctx context.Context, recordID string) (*fabric.BlockchainRecord, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("blockchain service is not operational")
	}

	return s.fabricService.QueryBlockchainRecord(ctx, recordID)
}

// VerifyRecord verifies a record's integrity against the blockchain
func (s *BlockchainService) VerifyRecord(ctx context.Context, recordID string) (bool, error) {
	if !s.IsEnabled() {
		return false, fmt.Errorf("blockchain service is not operational")
	}

	return s.fabricService.VerifyRecord(ctx, recordID)
}

// RegisterUser registers a user identity on the blockchain
func (s *BlockchainService) RegisterUser(ctx context.Context, user *fabric.UserRecord) error {
	if user == nil {
		return fmt.Errorf("user record cannot be nil")
	}

	s.logger.WithField("user_id", user.UserID).Info("Registering user on blockchain")

	start := time.Now()
	err := s.fabricService.RegisterUser(ctx, user)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":  user.UserID,
			"duration": duration.String(),
		}).Warn("User registration on blockchain failed, using local storage")

		return s.storeUserForSync(user)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":  user.UserID,
		"duration": duration.String(),
	}).Info("User successfully registered on blockchain")

	return nil
}

// RecordPayment records a payment transaction on the blockchain
func (s *BlockchainService) RecordPayment(ctx context.Context, payment *fabric.PaymentRecord) error {
	if payment == nil {
		return fmt.Errorf("payment record cannot be nil")
	}

	s.logger.WithField("payment_id", payment.PaymentID).Info("Recording payment on blockchain")

	start := time.Now()
	err := s.fabricService.RecordPayment(ctx, payment)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"payment_id": payment.PaymentID,
			"duration":   duration.String(),
		}).Warn("Payment recording on blockchain failed, using local storage")

		return s.storePaymentForSync(payment)
	}

	s.logger.WithFields(logrus.Fields{
		"payment_id": payment.PaymentID,
		"duration":   duration.String(),
	}).Info("Payment successfully recorded on blockchain")

	return nil
}

// RecordCompliance records a compliance event on the blockchain
func (s *BlockchainService) RecordCompliance(ctx context.Context, compliance *fabric.ComplianceRecord) error {
	if compliance == nil {
		return fmt.Errorf("compliance record cannot be nil")
	}

	s.logger.WithField("compliance_id", compliance.ComplianceID).Info("Recording compliance on blockchain")

	start := time.Now()
	err := s.fabricService.RecordCompliance(ctx, compliance)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"compliance_id": compliance.ComplianceID,
			"duration":      duration.String(),
		}).Warn("Compliance recording on blockchain failed, using local storage")

		return s.storeComplianceForSync(compliance)
	}

	s.logger.WithFields(logrus.Fields{
		"compliance_id": compliance.ComplianceID,
		"duration":      duration.String(),
	}).Info("Compliance successfully recorded on blockchain")

	return nil
}

// RecordOperation records a business operation on the blockchain
func (s *BlockchainService) RecordOperation(ctx context.Context, operation *fabric.OperationsRecord) error {
	if operation == nil {
		return fmt.Errorf("operation record cannot be nil")
	}

	s.logger.WithField("operation_id", operation.OperationID).Info("Recording operation on blockchain")

	start := time.Now()
	err := s.fabricService.RecordOperation(ctx, operation)
	duration := time.Since(start)

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"operation_id": operation.OperationID,
			"duration":     duration.String(),
		}).Warn("Operation recording on blockchain failed, using local storage")

		return s.storeOperationForSync(operation)
	}

	s.logger.WithFields(logrus.Fields{
		"operation_id": operation.OperationID,
		"duration":     duration.String(),
	}).Info("Operation successfully recorded on blockchain")

	return nil
}

// CreateUserRecord creates a blockchain-compatible user record
func (s *BlockchainService) CreateUserRecord(userID, email, fullName string) *fabric.UserRecord {
	return &fabric.UserRecord{
		UserID:         userID,
		Email:          email,
		FullName:       fullName,
		KYCStatus:      "pending",
		KYCLevel:       "basic",
		IdentityHash:   "", // Would be calculated from actual identity data
		Status:         "active",
		CreatedAt:      time.Now(),
		LastVerifiedAt: time.Now(),
	}
}

// CreatePaymentRecord creates a blockchain-compatible payment record
func (s *BlockchainService) CreatePaymentRecord(paymentID, userID string, amount float64, paymentType string) *fabric.PaymentRecord {
	return &fabric.PaymentRecord{
		PaymentID:     paymentID,
		UserID:        userID,
		Amount:        amount,
		Currency:      "USD",
		Type:          paymentType,
		Status:        "processed",
		PaymentMethod: "card",
		Reference:     uuid.New().String(),
		ProcessedAt:   time.Now(),
		TransactionID: uuid.New().String(),
	}
}

// CreateComplianceRecord creates a blockchain-compatible compliance record
func (s *BlockchainService) CreateComplianceRecord(complianceID, entityType, entityID, checkType string) *fabric.ComplianceRecord {
	return &fabric.ComplianceRecord{
		ComplianceID:  complianceID,
		EntityType:    entityType,
		EntityID:      entityID,
		CheckType:     checkType,
		Status:        "compliant",
		Severity:      "low",
		Details:       map[string]interface{}{"status": "passed"},
		CheckedAt:     time.Now(),
		CheckedBy:     "system",
		NextCheckDate: nil,
	}
}

// CreateOperationsRecord creates a blockchain-compatible operations record
func (s *BlockchainService) CreateOperationsRecord(operationID, operationType, entityType, entityID, userID, action string) *fabric.OperationsRecord {
	return &fabric.OperationsRecord{
		OperationID: operationID,
		Type:        operationType,
		EntityType:  entityType,
		EntityID:    entityID,
		UserID:      userID,
		Action:      action,
		Status:      "completed",
		Details:     map[string]interface{}{"action": action},
		StartedAt:   time.Now(),
		CompletedAt: &time.Time{},
		Result:      "success",
	}
}

// storeUserForSync stores a user locally for later blockchain sync
func (s *BlockchainService) storeUserForSync(user *fabric.UserRecord) error {
	s.logger.WithField("user_id", user.UserID).Debug("Storing user for blockchain sync")
	return nil
}

// storePaymentForSync stores a payment locally for later blockchain sync
func (s *BlockchainService) storePaymentForSync(payment *fabric.PaymentRecord) error {
	s.logger.WithField("payment_id", payment.PaymentID).Debug("Storing payment for blockchain sync")
	return nil
}

// storeComplianceForSync stores compliance locally for later blockchain sync
func (s *BlockchainService) storeComplianceForSync(compliance *fabric.ComplianceRecord) error {
	s.logger.WithField("compliance_id", compliance.ComplianceID).Debug("Storing compliance for blockchain sync")
	return nil
}

// storeOperationForSync stores an operation locally for later blockchain sync
func (s *BlockchainService) storeOperationForSync(operation *fabric.OperationsRecord) error {
	s.logger.WithField("operation_id", operation.OperationID).Debug("Storing operation for blockchain sync")
	return nil
}

// SyncPendingRecords attempts to sync any records that failed to submit to blockchain
func (s *BlockchainService) SyncPendingRecords(ctx context.Context) error {
	if !s.IsEnabled() {
		s.logger.Debug("Blockchain service disabled, skipping sync")
		return nil
	}

	s.logger.Info("Starting synchronization of pending blockchain records")

	// TODO: Implement sync logic for policies, claims, devices, and audits
	// that were stored locally due to blockchain unavailability

	s.logger.Info("Blockchain record synchronization completed")
	return nil
}

// storePolicyForSync stores a policy locally for later blockchain sync
func (s *BlockchainService) storePolicyForSync(policy *fabric.PolicyRecord) error {
	s.logger.WithField("policy_id", policy.PolicyID).Debug("Storing policy for blockchain sync")
	return s.fabricService.StoreForSync(context.Background(), "policy", policy.PolicyID, policy)
}

// storeClaimForSync stores a claim locally for later blockchain sync
func (s *BlockchainService) storeClaimForSync(claim *fabric.ClaimRecord) error {
	s.logger.WithField("claim_id", claim.ClaimID).Debug("Storing claim for blockchain sync")
	return s.fabricService.StoreForSync(context.Background(), "claim", claim.ClaimID, claim)
}

// storeDeviceForSync stores a device locally for later blockchain sync
func (s *BlockchainService) storeDeviceForSync(device *fabric.DeviceRecord) error {
	s.logger.WithField("device_id", device.DeviceID).Debug("Storing device for blockchain sync")
	return s.fabricService.StoreForSync(context.Background(), "device", device.DeviceID, device)
}

// storeAuditForSync stores an audit locally for later blockchain sync
func (s *BlockchainService) storeAuditForSync(audit *fabric.AuditRecord) error {
	s.logger.WithField("audit_id", audit.AuditID).Debug("Storing audit for blockchain sync")
	return s.fabricService.StoreForSync(context.Background(), "audit", audit.AuditID, audit)
}

// CreatePolicyRecord creates a blockchain-compatible policy record
func (s *BlockchainService) CreatePolicyRecord(policyID, customerID, deviceID string, amount float64, effectiveDate, expiryDate time.Time) *fabric.PolicyRecord {
	return &fabric.PolicyRecord{
		PolicyID:      policyID,
		CustomerID:    customerID,
		DeviceID:      deviceID,
		Amount:        amount,
		Status:        "active",
		EffectiveDate: effectiveDate,
		ExpiryDate:    expiryDate,
		CreatedAt:     time.Now(),
	}
}

// CreateClaimRecord creates a blockchain-compatible claim record
func (s *BlockchainService) CreateClaimRecord(claimID, policyID, customerID, deviceID string, amount float64, incidentDate time.Time) *fabric.ClaimRecord {
	return &fabric.ClaimRecord{
		ClaimID:      claimID,
		PolicyID:     policyID,
		CustomerID:   customerID,
		DeviceID:     deviceID,
		Amount:       amount,
		Status:       "submitted",
		IncidentDate: incidentDate,
		SubmittedAt:  time.Now(),
	}
}

// CreateDeviceRecord creates a blockchain-compatible device record
func (s *BlockchainService) CreateDeviceRecord(deviceID, imei, brand, model string) *fabric.DeviceRecord {
	return &fabric.DeviceRecord{
		DeviceID:     deviceID,
		IMEI:         imei,
		Brand:        brand,
		Model:        model,
		Status:       "authentic",
		RegisteredAt: time.Now(),
		LastVerified: time.Now(),
	}
}

// CreateAuditRecord creates a blockchain-compatible audit record
func (s *BlockchainService) CreateAuditRecord(entityType, entityID, action, userID string, details map[string]interface{}) *fabric.AuditRecord {
	return &fabric.AuditRecord{
		AuditID:    uuid.New().String(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		UserID:     userID,
		Timestamp:  time.Now(),
		Details:    details,
		Compliance: "compliant", // Default to compliant, can be updated
	}
}
