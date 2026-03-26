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
    
    -- Mid-term device changes
    original_device_id UUID REFERENCES devices(id), -- For device swaps
    is_device_swap BOOLEAN DEFAULT FALSE,
    swap_date TIMESTAMPTZ,
    swap_reason VARCHAR(50), -- upgrade, replacement, repair_exchange
    pro_rata_adjustment_amount DECIMAL(28,8), -- Premium adjustment for swap
    
    -- Enterprise/BYOD specifics
    is_byod_policy BOOLEAN DEFAULT FALSE,
    mdm_enrolled BOOLEAN DEFAULT FALSE,
    mdm_enrollment_date TIMESTAMPTZ,
    security_policy_compliant BOOLEAN DEFAULT TRUE,
    jailbreak_detected BOOLEAN DEFAULT FALSE,
    
    -- Family plan / corporate allocation
    family_group_id UUID REFERENCES family_groups(id),
    corporate_fleet_id UUID REFERENCES corporate_fleet_policies(id),
    cost_center_code VARCHAR(50), -- For corporate billing allocation
    
    -- Earned premium tracking (daily pro-rata for UPR)
    daily_earned_rate DECIMAL(28,8), -- 1/365th of annual premium
    earned_premium_to_date DECIMAL(28,8) DEFAULT 0,
    unearned_premium_reserve DECIMAL(28,8), -- UPR balance
    
    -- IFRS 17 Fields
    ifrs17_contract_type ifrs17_contract_type DEFAULT 'non_participating',
    discount_rate DECIMAL(10,6),
    risk_margin DECIMAL(28,8),
    fulfilment_cash_flows JSONB, -- Calculated at policy inception
    
    -- Renewal tracking
    renewal_eligible BOOLEAN DEFAULT FALSE,
    renewal_reminder_sent_at TIMESTAMPTZ,
    renewal_invited_at TIMESTAMPTZ,
    
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

-- Temporary device coverage (loaner phones during repair)
CREATE TABLE temporary_device_coverage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    original_policy_id UUID NOT NULL REFERENCES policies(id),
    original_device_id UUID NOT NULL REFERENCES devices(id),
    
    loaner_device_id UUID REFERENCES devices(id),
    loaner_imei VARCHAR(17),
    loaner_model VARCHAR(100),
    
    -- Coverage period
    coverage_start_date DATE NOT NULL,
    coverage_end_date DATE,
    
    -- Claim handling
    max_claim_amount DECIMAL(28,8),
    deductible DECIMAL(28,8),
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'claimed', 'returned')),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_temp_coverage_policy ON temporary_device_coverage(original_policy_id);
CREATE INDEX idx_temp_coverage_status ON temporary_device_coverage(status);

-- Multi-device bundles (family plan aggregates)
CREATE TABLE policy_device_bundles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    master_policy_id UUID NOT NULL REFERENCES policies(id),
    device_id UUID NOT NULL REFERENCES devices(id),
    
    -- Bundle pricing
    is_primary_device BOOLEAN DEFAULT FALSE,
    bundle_discount_percentage DECIMAL(5,2) DEFAULT 0,
    
    -- Aggregate deductibles
    deductible_tier INTEGER DEFAULT 1, -- 1st claim = tier 1 deductible, 2nd = tier 2, etc.
    tier_1_deductible DECIMAL(28,8),
    tier_2_deductible DECIMAL(28,8),
    tier_3_plus_deductible DECIMAL(28,8),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(master_policy_id, device_id)
);

CREATE INDEX idx_device_bundles_policy ON policy_device_bundles(master_policy_id);

-- Daily earned premium recognition (for IFRS/UPR)
CREATE TABLE earned_premium_daily (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    
    recognition_date DATE NOT NULL,
    daily_amount DECIMAL(28,8) NOT NULL,
    cumulative_earned DECIMAL(28,8) NOT NULL,
    unearned_balance DECIMAL(28,8) NOT NULL,
    
    movement_id UUID REFERENCES value_movements(id), -- Accounting entry
    
    UNIQUE(policy_id, recognition_date)
);

CREATE INDEX idx_earned_premium_policy ON earned_premium_daily(policy_id);
CREATE INDEX idx_earned_premium_date ON earned_premium_daily(recognition_date);

-- Premium calculation audit for regulatory review
CREATE TABLE premium_calculation_audit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    calculation_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    -- Inputs
    base_rate DECIMAL(10,6),
    sum_insured DECIMAL(28,8),
    device_age_months INTEGER,
    device_condition VARCHAR(20),
    risk_score DECIMAL(5,2),
    discounts_applied JSONB, -- Array of applied discounts
    loadings_applied JSONB, -- Array of applied loadings
    
    -- Calculation steps
    calculated_premium DECIMAL(28,8),
    tax_amount DECIMAL(28,8),
    total_premium DECIMAL(28,8),
    
    -- Audit trail
    calculated_by UUID REFERENCES parties(id),
    calculation_model_version VARCHAR(20),
    input_hash VARCHAR(64), -- Hash of inputs for integrity
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_premium_calc_policy ON premium_calculation_audit(policy_id);
CREATE INDEX idx_premium_calc_date ON premium_calculation_audit(calculation_date);
