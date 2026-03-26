-- 005_policies.sql

-- Insurance policies issued
CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_number VARCHAR(50) NOT NULL UNIQUE,
    
    -- Relationships
    product_contract_id UUID NOT NULL REFERENCES product_contracts(id),
    policyholder_party_id UUID NOT NULL REFERENCES parties(id),
    corporate_employee_id UUID REFERENCES corporate_employees(id), -- Optional corporate
    device_id UUID NOT NULL REFERENCES devices(id),
    
    -- Coverage details
    sum_insured DECIMAL(28,8) NOT NULL,
    deductible DECIMAL(28,8) NOT NULL DEFAULT 0,
    premium_amount DECIMAL(28,8) NOT NULL, -- Per billing period
    premium_currency CHAR(3) DEFAULT 'USD',
    billing_frequency VARCHAR(20) DEFAULT 'monthly',
    
    -- Dates
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    cooling_off_end DATE, -- Cancellation period
    renewal_policy_id UUID REFERENCES policies(id), -- Linked renewal
    
    -- Status
    status policy_status DEFAULT 'pending',
    cancellation_reason VARCHAR(100),
    cancelled_at TIMESTAMPTZ,
    cancellation_initiated_by UUID REFERENCES parties(id),
    
    -- Bitemporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    CONSTRAINT valid_coverage_dates CHECK (start_date < end_date),
    CONSTRAINT valid_premium CHECK (premium_amount > 0),
    CONSTRAINT valid_sum_insured CHECK (sum_insured > 0)
);

-- Policy version history (explicit versioning for audits)
CREATE TABLE policy_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    version_number INTEGER NOT NULL,
    change_reason VARCHAR(255),
    changed_by UUID REFERENCES parties(id),
    previous_values JSONB, -- Snapshot of changed fields
    valid_from TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(policy_id, version_number)
);

-- Indexes
CREATE INDEX idx_policies_number ON policies(policy_number);
CREATE INDEX idx_policies_holder ON policies(policyholder_party_id);
CREATE INDEX idx_policies_device ON policies(device_id);
CREATE INDEX idx_policies_contract ON policies(product_contract_id);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_dates ON policies(start_date, end_date);
CREATE INDEX idx_policies_current ON policies(is_current) WHERE is_current = TRUE;
CREATE INDEX idx_policies_valid_time ON policies(valid_from, valid_to);

-- Trigger
CREATE TRIGGER trg_policies_hash BEFORE INSERT OR UPDATE ON policies
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();
