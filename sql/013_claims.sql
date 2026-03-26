-- 013_claims.sql

-- Insurance claims
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_number VARCHAR(50) NOT NULL UNIQUE,
    
    -- Relationships
    policy_id UUID NOT NULL REFERENCES policies(id),
    device_id UUID REFERENCES devices(id), -- Denormalized for performance
    
    -- Claim details
    claim_type claim_type NOT NULL,
    incident_date DATE NOT NULL,
    reported_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT,
    incident_location VARCHAR(255),
    
    -- Amounts
    claimed_amount DECIMAL(28,8) NOT NULL,
    approved_amount DECIMAL(28,8),
    deductible_applied DECIMAL(28,8) DEFAULT 0,
    excess_amount DECIMAL(28,8) DEFAULT 0, -- Non-covered portion
    
    -- Status workflow
    status claim_status DEFAULT 'submitted',
    
    -- Timeline
    submitted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    assigned_at TIMESTAMPTZ,
    assessed_at TIMESTAMPTZ,
    decision_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    reopened_at TIMESTAMPTZ,
    
    -- Assignment
    assessor_party_id UUID REFERENCES parties(id),
    adjuster_party_id UUID REFERENCES parties(id),
    handled_by UUID REFERENCES parties(id), -- Internal handler
    
    -- Fraud and risk
    fraud_score DECIMAL(5,2) CHECK (fraud_score >= 0 AND fraud_score <= 100),
    fraud_flags JSONB, -- Array of triggered rules
    investigation_required BOOLEAN DEFAULT FALSE,
    investigation_notes TEXT,
    
    -- Payment
    payout_movement_id UUID REFERENCES value_movements(id),
    settlement_method settlement_method,
    settlement_amount DECIMAL(28,8),
    
    -- Real-time authorization link
    real_time_posting_id UUID, -- Link to real_time_postings table
    
    -- Bitemporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    CONSTRAINT valid_incident_date CHECK (incident_date <= CURRENT_DATE),
    CONSTRAINT valid_amounts CHECK (approved_amount IS NULL OR approved_amount <= claimed_amount)
);

-- Claim reserve link (explicit many-to-many if needed)
CREATE TABLE claim_reserve_links (
    claim_id UUID REFERENCES claims(id),
    reserve_id UUID REFERENCES reserves(id),
    PRIMARY KEY (claim_id, reserve_id)
);

-- Indexes
CREATE INDEX idx_claims_number ON claims(claim_number);
CREATE INDEX idx_claims_policy ON claims(policy_id);
CREATE INDEX idx_claims_status ON claims(status);
CREATE INDEX idx_claims_type ON claims(claim_type);
CREATE INDEX idx_claims_fraud ON claims(fraud_score) WHERE fraud_score > 50;
CREATE INDEX idx_claims_dates ON claims(incident_date, submitted_at);
CREATE INDEX idx_claims_current ON claims(is_current) WHERE is_current = TRUE;

-- Trigger
CREATE TRIGGER trg_claims_hash BEFORE INSERT OR UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();
