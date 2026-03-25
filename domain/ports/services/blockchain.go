package services

import (
	"context"
	"smartsure/pkg/blockchain/fabric"
)

// BlockchainService defines the interface for blockchain operations
type BlockchainService interface {
	// Core blockchain operations
	IsEnabled() bool
	GetBlockchainStatus() map[string]interface{}

	// Policy operations
	SubmitPolicy(ctx context.Context, policy *fabric.PolicyRecord) error

	// Claims operations
	SubmitClaim(ctx context.Context, claim *fabric.ClaimRecord) error

	// Device operations
	RegisterDevice(ctx context.Context, device *fabric.DeviceRecord) error

	// User operations
	RegisterUser(ctx context.Context, user *fabric.UserRecord) error

	// Payment operations
	RecordPayment(ctx context.Context, payment *fabric.PaymentRecord) error

	// Compliance operations
	RecordCompliance(ctx context.Context, compliance *fabric.ComplianceRecord) error

	// Operations recording
	RecordOperation(ctx context.Context, operation *fabric.OperationsRecord) error

	// Audit operations
	SubmitAudit(ctx context.Context, audit *fabric.AuditRecord) error

	// Query operations
	QueryBlockchainRecord(ctx context.Context, recordID string) (*fabric.BlockchainRecord, error)
	VerifyRecord(ctx context.Context, recordID string) (bool, error)
}
