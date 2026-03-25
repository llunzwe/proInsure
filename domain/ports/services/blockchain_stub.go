package services

import (
	"context"
)

// BlockchainService defines the interface for blockchain operations (stub version without fabric dependency)
type BlockchainService interface {
	// Core blockchain operations
	IsEnabled() bool
	GetBlockchainStatus() map[string]interface{}

	// Policy operations
	SubmitPolicy(ctx context.Context, policy interface{}) error

	// Claims operations
	SubmitClaim(ctx context.Context, claim interface{}) error

	// Device operations
	RegisterDevice(ctx context.Context, device interface{}) error

	// User operations
	RegisterUser(ctx context.Context, user interface{}) error

	// Payment operations
	RecordPayment(ctx context.Context, payment interface{}) error

	// Compliance operations
	RecordCompliance(ctx context.Context, compliance interface{}) error

	// Operations recording
	RecordOperation(ctx context.Context, operation interface{}) error

	// Audit operations
	SubmitAudit(ctx context.Context, audit interface{}) error

	// Query operations
	QueryBlockchainRecord(ctx context.Context, recordID string) (interface{}, error)
	VerifyRecord(ctx context.Context, recordID string) (bool, error)

	// Record creation helpers
	CreateUserRecord(userID, email, fullName string) interface{}
	CreatePaymentRecord(paymentID, userID string, amount float64, paymentType string) interface{}
	CreateComplianceRecord(complianceID, entityType, entityID, checkType string) interface{}
	CreateOperationsRecord(operationID, operationType, entityType, entityID, userID, action string) interface{}
	CreatePolicyRecord(policyID, customerID, deviceID string, amount float64, effectiveDate, expiryDate interface{}) interface{}
	CreateClaimRecord(claimID, policyID, customerID, deviceID string, amount float64, incidentDate interface{}) interface{}
	CreateDeviceRecord(deviceID, imei, brand, model string) interface{}
	CreateAuditRecord(entityType, entityID, action, userID string, details map[string]interface{}) interface{}

	// Sync operations
	SyncPendingRecords(ctx context.Context) error
}
