-- 026_transaction_entity.sql

-- First-class Transaction Entity for chain-of-custody
CREATE TABLE transaction_entities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_reference VARCHAR(100) NOT NULL UNIQUE,
    
    -- Classification
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN (
        'policy_issuance', 'policy_renewal', 'policy_cancellation', 'policy_modification',
        'premium_payment', 'premium_refund',
        'claim_submission', 'claim_approval', 'claim_rejection', 'claim_payment',
        'reserve_creation', 'reserve_release',
        'device_transfer', 'ownership_change'
    )),
    
    -- Status
    status VARCHAR(20) DEFAULT 'initiated' CHECK (status IN ('initiated', 'in_progress', 'completed', 'failed', 'reversed')),
    
    -- Chain of custody
    initiated_by UUID REFERENCES parties(id),
    initiated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    
    -- Financial summary
    total_amount DECIMAL(28,8),
    currency CHAR(3),
    movement_id UUID REFERENCES value_movements(id),
    
    -- Links to business entities
    primary_entity_type VARCHAR(50),
    primary_entity_id UUID,
    related_entities JSONB, -- Array of {type, id, role}
    
    -- Compliance
    compliance_flags JSONB,
    sanctions_checked BOOLEAN DEFAULT FALSE,
    sanctions_check_passed BOOLEAN,
    
    -- Reversal tracking
    reversed_by UUID REFERENCES transaction_entities(id),
    reversal_reason TEXT,
    
    -- Event sourcing link
    first_event_id BIGINT,
    last_event_id BIGINT,
    
    -- Immutable hash
    immutable_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Transaction steps (detailed audit trail)
CREATE TABLE transaction_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transaction_entities(id) ON DELETE CASCADE,
    
    step_number INTEGER NOT NULL,
    step_name VARCHAR(100) NOT NULL,
    step_type VARCHAR(50), -- validation, authorization, processing, notification
    
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed', 'skipped')),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    
    input_data JSONB,
    output_data JSONB,
    error_details JSONB,
    
    performed_by UUID REFERENCES parties(id),
    event_id BIGINT REFERENCES immutable_events(event_id),
    
    UNIQUE(transaction_id, step_number)
);

-- Transaction entity links (many-to-many with business entities)
CREATE TABLE transaction_entity_links (
    transaction_id UUID REFERENCES transaction_entities(id),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    link_type VARCHAR(50) DEFAULT 'primary', -- primary, affected, referenced
    PRIMARY KEY (transaction_id, entity_type, entity_id)
);

CREATE INDEX idx_transactions_ref ON transaction_entities(transaction_reference);
CREATE INDEX idx_transactions_type ON transaction_entities(transaction_type);
CREATE INDEX idx_transactions_status ON transaction_entities(status);
CREATE INDEX idx_transactions_entity ON transaction_entities(primary_entity_type, primary_entity_id);
CREATE INDEX idx_transaction_steps_txn ON transaction_steps(transaction_id);
