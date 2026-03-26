-- 023_control_and_batch_processing.sql

-- Batch processing control
CREATE TABLE control_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_reference VARCHAR(100) NOT NULL UNIQUE,
    
    batch_type VARCHAR(50) NOT NULL CHECK (batch_type IN (
        'premium_billing', 'renewal_processing', 'reserve_calculation',
        'claims_escalation', 'eod', 'eom', 'eoy', 'data_archive'
    )),
    
    -- Status workflow
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled')),
    
    -- Control totals (double-entry conservation)
    expected_count INTEGER,
    processed_count INTEGER DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    
    expected_amount DECIMAL(28,8),
    actual_amount DECIMAL(28,8),
    
    -- Timing
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    
    -- Configuration
    parameters JSONB, -- Input parameters for the batch
    
    -- Error handling
    error_log JSONB, -- Array of error details
    warning_log JSONB,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id)
);

-- Individual batch entries
CREATE TABLE batch_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES control_batches(id) ON DELETE CASCADE,
    
    sequence_number INTEGER NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    
    -- Processing status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'skipped')),
    
    -- Input/Output
    input_data JSONB,
    output_data JSONB,
    error_message TEXT,
    
    -- Timing
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    duration_ms INTEGER
);

-- Batch dependencies (DAG support)
CREATE TABLE batch_dependencies (
    batch_id UUID REFERENCES control_batches(id),
    depends_on_batch_id UUID REFERENCES control_batches(id),
    PRIMARY KEY (batch_id, depends_on_batch_id)
);

CREATE INDEX idx_batches_status ON control_batches(status);
CREATE INDEX idx_batches_type ON control_batches(batch_type);
CREATE INDEX idx_batches_scheduled ON control_batches(scheduled_at);
CREATE INDEX idx_batch_entries_batch ON batch_entries(batch_id);
CREATE INDEX idx_batch_entries_status ON batch_entries(status);
