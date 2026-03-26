-- 021_real_time_claim_authorization.sql

-- Low-latency pre-authorization (<10ms target)
CREATE TABLE real_time_postings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Request context
    request_id VARCHAR(100) NOT NULL UNIQUE,
    claim_id UUID REFERENCES claims(id), -- May be null if pre-auth before claim creation
    
    -- Decision
    auth_status VARCHAR(20) NOT NULL CHECK (auth_status IN ('APPROVED', 'DECLINED', 'PENDING_REVIEW')),
    auth_type VARCHAR(20) NOT NULL CHECK (auth_type IN ('REAL_TIME', 'PRE_AUTH', 'EMERGENCY')),
    
    -- Financials
    requested_amount DECIMAL(28,8) NOT NULL,
    approved_amount DECIMAL(28,8),
    currency CHAR(3) DEFAULT 'USD',
    
    -- Risk checks
    velocity_check_passed BOOLEAN,
    imei_validation_passed BOOLEAN,
    photo_validation_passed BOOLEAN,
    fraud_score DECIMAL(5,2),
    risk_factors_triggered JSONB,
    
    -- JIT Funding
    jit_funding_required BOOLEAN DEFAULT FALSE,
    jit_funding_source_container_id UUID REFERENCES value_containers(id),
    jit_funding_movement_id UUID REFERENCES value_movements(id),
    ring_fence_applied BOOLEAN DEFAULT FALSE,
    ring_fence_container_id UUID REFERENCES value_containers(id),
    
    -- Override
    commando_override BOOLEAN DEFAULT FALSE, -- Emergency bypass
    override_reason TEXT,
    override_approved_by UUID REFERENCES parties(id),
    
    -- Performance metrics
    decision_time_ms INTEGER, -- Target <10ms
    queue_time_ms INTEGER,
    total_processing_ms INTEGER,
    
    -- Product reference
    product_contract_hash UUID, -- Snapshot of rules at decision time
    
    -- Timestamps
    requested_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    decided_at TIMESTAMPTZ,
    valid_until TIMESTAMPTZ, -- Pre-auth expiry
    
    -- Integrity
    immutable_hash VARCHAR(64) NOT NULL
);

-- Real-time rules engine log
CREATE TABLE real_time_rules_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    posting_id UUID NOT NULL REFERENCES real_time_postings(id),
    rule_name VARCHAR(100) NOT NULL,
    rule_version VARCHAR(20),
    result VARCHAR(20) CHECK (result IN ('pass', 'fail', 'review')),
    execution_time_microseconds INTEGER,
    input_data JSONB,
    output_data JSONB,
    executed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_realtime_claim ON real_time_postings(claim_id);
CREATE INDEX idx_realtime_status ON real_time_postings(auth_status);
CREATE INDEX idx_realtime_time ON real_time_postings(decided_at);
CREATE INDEX idx_realtime_request ON real_time_postings(request_id);
