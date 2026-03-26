-- 012_reserves_and_provisioning.sql

-- IFRS 17 Claim reserves
CREATE TABLE reserves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Link to policy or claim
    policy_id UUID REFERENCES policies(id),
    claim_id UUID REFERENCES claims(id),
    
    -- Reserve classification
    reserve_type reserve_type NOT NULL,
    calculation_method VARCHAR(50), -- actuarial, case_estimate, percentage
    amount DECIMAL(28,8) NOT NULL,
    currency CHAR(3) DEFAULT 'USD',
    
    -- Dates
    calculation_date DATE NOT NULL,
    accounting_period VARCHAR(20) NOT NULL, -- e.g., "2026-Q1"
    
    -- Status workflow
    status VARCHAR(20) DEFAULT 'proposed' CHECK (status IN ('proposed', 'booked', 'adjusted', 'released')),
    
    -- Link to accounting movement
    creation_movement_id UUID REFERENCES value_movements(id),
    release_movement_id UUID REFERENCES value_movements(id),
    
    -- Actuarial details
    confidence_level DECIMAL(5,2), -- e.g., 75.00 for 75th percentile
    discount_rate DECIMAL(10,6),
    discount_rate_source VARCHAR(50), -- e.g., 'risk_free_curve', 'company_portfolio'
    risk_margin DECIMAL(28,8), -- IFRS 17 risk adjustment
    risk_margin_method VARCHAR(50), -- e.g., 'cost_of_capital', 'confidence_level'
    projected_settlement_date DATE,
    
    -- Actuarial model tracking
    actuarial_model_id VARCHAR(50),
    actuarial_model_version VARCHAR(20),
    
    -- Reinsurance recoverable
    reinsurance_recoverable DECIMAL(28,8) DEFAULT 0,
    reinsurance_treaty_id UUID, -- Link to reinsurance treaty
    net_reserve_after_reinsurance DECIMAL(28,8), -- Calculated field
    
    -- Regulatory reporting
    solvency_ii_classification VARCHAR(50),
    local_regulatory_reference VARCHAR(100),
    
    -- Bitemporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    approved_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL
);

-- Reserve adjustments history
CREATE TABLE reserve_adjustments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reserve_id UUID NOT NULL REFERENCES reserves(id),
    previous_amount DECIMAL(28,8),
    new_amount DECIMAL(28,8),
    adjustment_reason VARCHAR(255),
    adjusted_by UUID REFERENCES parties(id),
    adjusted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    movement_id UUID REFERENCES value_movements(id)
);

-- Indexes
CREATE INDEX idx_reserves_policy ON reserves(policy_id);
CREATE INDEX idx_reserves_claim ON reserves(claim_id) WHERE claim_id IS NOT NULL;
CREATE INDEX idx_reserves_type ON reserves(reserve_type);
CREATE INDEX idx_reserves_period ON reserves(accounting_period);
CREATE INDEX idx_reserves_status ON reserves(status);
