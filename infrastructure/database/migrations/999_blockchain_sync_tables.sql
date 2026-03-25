-- Blockchain Sync Records table (moved to pkg/blockchain/models.go)
CREATE TABLE IF NOT EXISTS blockchain_sync_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    record_type VARCHAR(50) NOT NULL,
    record_id VARCHAR(255) NOT NULL,
    record_data TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 5,
    last_error TEXT,
    last_attempt_at TIMESTAMP,
    synced_at TIMESTAMP,
    blockchain_hash VARCHAR(255),
    blockchain_tx_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_blockchain_sync_records_type ON blockchain_sync_records(record_type);
CREATE INDEX IF NOT EXISTS idx_blockchain_sync_records_id ON blockchain_sync_records(record_id);
CREATE INDEX IF NOT EXISTS idx_blockchain_sync_records_status ON blockchain_sync_records(status);
CREATE INDEX IF NOT EXISTS idx_blockchain_sync_records_deleted_at ON blockchain_sync_records(deleted_at);

-- Async Operations table
CREATE TABLE IF NOT EXISTS async_operations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    operation_type VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    progress INTEGER DEFAULT 0,
    error TEXT,
    result TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_async_operations_type ON async_operations(operation_type);
CREATE INDEX IF NOT EXISTS idx_async_operations_entity_type ON async_operations(entity_type);
CREATE INDEX IF NOT EXISTS idx_async_operations_entity_id ON async_operations(entity_id);
CREATE INDEX IF NOT EXISTS idx_async_operations_status ON async_operations(status);
CREATE INDEX IF NOT EXISTS idx_async_operations_deleted_at ON async_operations(deleted_at);

-- Error Logs table
CREATE TABLE IF NOT EXISTS error_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    error_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    service VARCHAR(100) NOT NULL,
    operation VARCHAR(200) NOT NULL,
    error_code VARCHAR(50),
    error_message TEXT NOT NULL,
    stack_trace TEXT,
    context JSONB,
    user_id UUID,
    entity_type VARCHAR(50),
    entity_id UUID,
    resolved BOOLEAN DEFAULT false,
    resolved_at TIMESTAMP,
    resolved_by UUID,
    resolution TEXT,
    occurrence_count INTEGER DEFAULT 1,
    first_occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_error_logs_type ON error_logs(error_type);
CREATE INDEX IF NOT EXISTS idx_error_logs_severity ON error_logs(severity);
CREATE INDEX IF NOT EXISTS idx_error_logs_service ON error_logs(service);
CREATE INDEX IF NOT EXISTS idx_error_logs_operation ON error_logs(operation);
CREATE INDEX IF NOT EXISTS idx_error_logs_user_id ON error_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_error_logs_entity_type ON error_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_error_logs_entity_id ON error_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_error_logs_resolved ON error_logs(resolved);
CREATE INDEX IF NOT EXISTS idx_error_logs_deleted_at ON error_logs(deleted_at);

-- Idempotency Keys table
CREATE TABLE IF NOT EXISTS idempotency_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(255) UNIQUE NOT NULL,
    operation VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    request_hash VARCHAR(64) NOT NULL,
    response TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_idempotency_keys_key ON idempotency_keys(key);
CREATE INDEX IF NOT EXISTS idx_idempotency_keys_operation ON idempotency_keys(operation);
CREATE INDEX IF NOT EXISTS idx_idempotency_keys_entity_type ON idempotency_keys(entity_type);
CREATE INDEX IF NOT EXISTS idx_idempotency_keys_entity_id ON idempotency_keys(entity_id);
CREATE INDEX IF NOT EXISTS idx_idempotency_keys_status ON idempotency_keys(status);
CREATE INDEX IF NOT EXISTS idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);

-- Blockchain Records table (for audit trail)
CREATE TABLE IF NOT EXISTS blockchain_records (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    data JSONB NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    hash VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_blockchain_records_type ON blockchain_records(type);
CREATE INDEX IF NOT EXISTS idx_blockchain_records_status ON blockchain_records(status);
CREATE INDEX IF NOT EXISTS idx_blockchain_records_timestamp ON blockchain_records(timestamp);

