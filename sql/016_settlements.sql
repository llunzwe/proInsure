-- 016_settlements.sql

-- Payout finality tracking
CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    movement_id UUID NOT NULL REFERENCES value_movements(id),
    
    -- Settlement details
    method settlement_method NOT NULL,
    amount DECIMAL(28,8) NOT NULL,
    currency CHAR(3) DEFAULT 'USD',
    fees DECIMAL(28,8) DEFAULT 0, -- Transaction fees
    
    -- Payee information
    payee_party_id UUID REFERENCES parties(id),
    payee_account_details_encrypted BYTEA, -- Bank account, mobile money wallet, etc.
    payee_reference VARCHAR(100), -- Customer reference
    
    -- Status progression
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'provisional', 'final', 'failed', 'recalled')),
    
    -- Timing
    provisional_at TIMESTAMPTZ,
    final_at TIMESTAMPTZ,
    expected_finality_at TIMESTAMPTZ,
    
    -- External tracking
    external_reference VARCHAR(100), -- Bank transaction ID, etc.
    gateway_response JSONB,
    
    -- Proof and reconciliation
    proof_document_id UUID, -- Reference to uploaded proof of payment
    reconciliation_date DATE,
    reconciled_by UUID REFERENCES parties(id),
    
    -- For replacement devices (alternative to cash)
    replacement_device_id UUID, -- Reference to inventory
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL
);

-- Settlement reconciliation items
CREATE TABLE settlement_reconciliation (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    settlement_id UUID NOT NULL REFERENCES settlements(id),
    reconciliation_batch_id UUID,
    expected_amount DECIMAL(28,8),
    actual_amount DECIMAL(28,8),
    difference DECIMAL(28,8),
    status VARCHAR(20), -- matched, unmatched, discrepancy
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_settlements_claim ON settlements(claim_id);
CREATE INDEX idx_settlements_status ON settlements(status);
CREATE INDEX idx_settlements_finality ON settlements(final_at);
