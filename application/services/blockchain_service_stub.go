package services

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"smartsure/internal/domain/ports/services"
)

// StubFabricConfig is a stub config that doesn't require fabric package
type StubFabricConfig struct {
	Enabled bool
}

// BlockchainServiceStub implements the blockchain service interface without fabric dependency
type BlockchainServiceStub struct {
	logger  *logrus.Logger
	db      *gorm.DB
	enabled bool
}

// NewBlockchainServiceStub creates a new stub blockchain service
func NewBlockchainServiceStub(db *gorm.DB, logger *logrus.Logger, enabled bool) services.BlockchainService {
	service := &BlockchainServiceStub{
		logger:  logger,
		db:      db,
		enabled: enabled,
	}

	if service.enabled {
		logger.Warn("Blockchain service requested but fabric SDK unavailable - running in stub mode")
	} else {
		logger.Info("Blockchain service disabled (stub mode)")
	}

	return service
}

// IsEnabled returns whether the blockchain service is enabled
func (s *BlockchainServiceStub) IsEnabled() bool {
	return false // Always disabled in stub mode
}

// GetBlockchainStatus returns the current status of the blockchain service
func (s *BlockchainServiceStub) GetBlockchainStatus() map[string]interface{} {
	return map[string]interface{}{
		"service_enabled": false,
		"overall_status":  "disabled",
		"reason":          "Fabric SDK unavailable - running in stub mode",
	}
}

// SubmitPolicy is a stub implementation
func (s *BlockchainServiceStub) SubmitPolicy(ctx context.Context, policy interface{}) error {
	s.logger.Debug("SubmitPolicy called but blockchain is disabled (stub mode)")
	return nil
}

// SubmitClaim is a stub implementation
func (s *BlockchainServiceStub) SubmitClaim(ctx context.Context, claim interface{}) error {
	s.logger.Debug("SubmitClaim called but blockchain is disabled (stub mode)")
	return nil
}

// RegisterDevice is a stub implementation
func (s *BlockchainServiceStub) RegisterDevice(ctx context.Context, device interface{}) error {
	s.logger.Debug("RegisterDevice called but blockchain is disabled (stub mode)")
	return nil
}

// SubmitAudit is a stub implementation
func (s *BlockchainServiceStub) SubmitAudit(ctx context.Context, audit interface{}) error {
	s.logger.Debug("SubmitAudit called but blockchain is disabled (stub mode)")
	return nil
}

// QueryBlockchainRecord is a stub implementation
func (s *BlockchainServiceStub) QueryBlockchainRecord(ctx context.Context, recordID string) (interface{}, error) {
	return nil, fmt.Errorf("blockchain service is not operational (stub mode)")
}

// VerifyRecord is a stub implementation
func (s *BlockchainServiceStub) VerifyRecord(ctx context.Context, recordID string) (bool, error) {
	return false, fmt.Errorf("blockchain service is not operational (stub mode)")
}

// RegisterUser is a stub implementation
func (s *BlockchainServiceStub) RegisterUser(ctx context.Context, user interface{}) error {
	s.logger.Debug("RegisterUser called but blockchain is disabled (stub mode)")
	return nil
}

// RecordPayment is a stub implementation
func (s *BlockchainServiceStub) RecordPayment(ctx context.Context, payment interface{}) error {
	s.logger.Debug("RecordPayment called but blockchain is disabled (stub mode)")
	return nil
}

// RecordCompliance is a stub implementation
func (s *BlockchainServiceStub) RecordCompliance(ctx context.Context, compliance interface{}) error {
	s.logger.Debug("RecordCompliance called but blockchain is disabled (stub mode)")
	return nil
}

// RecordOperation is a stub implementation
func (s *BlockchainServiceStub) RecordOperation(ctx context.Context, operation interface{}) error {
	s.logger.Debug("RecordOperation called but blockchain is disabled (stub mode)")
	return nil
}

// CreateUserRecord is a stub implementation
func (s *BlockchainServiceStub) CreateUserRecord(userID, email, fullName string) interface{} {
	return map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"name":    fullName,
	}
}

// CreatePaymentRecord is a stub implementation
func (s *BlockchainServiceStub) CreatePaymentRecord(paymentID, userID string, amount float64, paymentType string) interface{} {
	return map[string]interface{}{
		"payment_id": paymentID,
		"user_id":    userID,
		"amount":     amount,
		"type":       paymentType,
	}
}

// CreateComplianceRecord is a stub implementation
func (s *BlockchainServiceStub) CreateComplianceRecord(complianceID, entityType, entityID, checkType string) interface{} {
	return map[string]interface{}{
		"compliance_id": complianceID,
		"entity_type":   entityType,
		"entity_id":     entityID,
		"check_type":    checkType,
	}
}

// CreateOperationsRecord is a stub implementation
func (s *BlockchainServiceStub) CreateOperationsRecord(operationID, operationType, entityType, entityID, userID, action string) interface{} {
	return map[string]interface{}{
		"operation_id": operationID,
		"type":         operationType,
		"entity_type":  entityType,
		"entity_id":    entityID,
		"user_id":      userID,
		"action":       action,
	}
}

// CreatePolicyRecord is a stub implementation
func (s *BlockchainServiceStub) CreatePolicyRecord(policyID, customerID, deviceID string, amount float64, effectiveDate, expiryDate interface{}) interface{} {
	return map[string]interface{}{
		"policy_id":      policyID,
		"customer_id":    customerID,
		"device_id":      deviceID,
		"amount":         amount,
		"effective_date": effectiveDate,
		"expiry_date":    expiryDate,
	}
}

// CreateClaimRecord is a stub implementation
func (s *BlockchainServiceStub) CreateClaimRecord(claimID, policyID, customerID, deviceID string, amount float64, incidentDate interface{}) interface{} {
	return map[string]interface{}{
		"claim_id":      claimID,
		"policy_id":     policyID,
		"customer_id":   customerID,
		"device_id":     deviceID,
		"amount":        amount,
		"incident_date": incidentDate,
	}
}

// CreateDeviceRecord is a stub implementation
func (s *BlockchainServiceStub) CreateDeviceRecord(deviceID, imei, brand, model string) interface{} {
	return map[string]interface{}{
		"device_id": deviceID,
		"imei":      imei,
		"brand":     brand,
		"model":     model,
	}
}

// CreateAuditRecord is a stub implementation
func (s *BlockchainServiceStub) CreateAuditRecord(entityType, entityID, action, userID string, details map[string]interface{}) interface{} {
	return map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
		"action":      action,
		"user_id":     userID,
		"details":     details,
	}
}

// SyncPendingRecords is a stub implementation
func (s *BlockchainServiceStub) SyncPendingRecords(ctx context.Context) error {
	s.logger.Debug("SyncPendingRecords called but blockchain is disabled (stub mode)")
	return nil
}
