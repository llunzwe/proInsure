-- 006_policy_discounts_and_limits.sql

-- Discounts applied to policies
CREATE TABLE policy_discounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    discount_type VARCHAR(50) NOT NULL CHECK (discount_type IN ('loyalty', 'bundle', 'no_claims', 'corporate', 'promotional', 'referral')),
    discount_code VARCHAR(50),
    percentage DECIMAL(5,2) CHECK (percentage >= 0 AND percentage <= 100),
    fixed_amount DECIMAL(28,8),
    currency CHAR(3) DEFAULT 'USD',
    reason TEXT,
    applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    applied_by UUID REFERENCES parties(id),
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    is_active BOOLEAN DEFAULT TRUE,
    
    CONSTRAINT discount_value_check CHECK (
        (percentage IS NOT NULL AND fixed_amount IS NULL) OR 
        (percentage IS NULL AND fixed_amount IS NOT NULL)
    )
);

-- Claim limits tracking (annual, lifetime, per-claim)
CREATE TABLE policy_claim_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    limit_type VARCHAR(50) NOT NULL CHECK (limit_type IN ('annual', 'lifetime', 'per_claim', 'consecutive_claims')),
    limit_amount DECIMAL(28,8) NOT NULL,
    used_amount DECIMAL(28,8) DEFAULT 0,
    remaining_amount DECIMAL(28,8) GENERATED ALWAYS AS (limit_amount - used_amount) STORED,
    currency CHAR(3) DEFAULT 'USD',
    claim_count_limit INTEGER, -- For count-based limits (e.g., max 2 claims per year)
    claim_count_used INTEGER DEFAULT 0,
    period_start DATE, -- For annual limits
    period_end DATE,
    reset_date DATE, -- When annual limits reset
    last_updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT limit_positive CHECK (limit_amount > 0),
    CONSTRAINT usage_not_negative CHECK (used_amount >= 0)
);

-- Trigger to update last_updated
CREATE OR REPLACE FUNCTION update_limit_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_limit_timestamp BEFORE UPDATE ON policy_claim_limits
    FOR EACH ROW EXECUTE FUNCTION update_limit_timestamp();

CREATE INDEX idx_discounts_policy ON policy_discounts(policy_id);
CREATE INDEX idx_discounts_active ON policy_discounts(is_current) WHERE is_active = TRUE;
CREATE INDEX idx_claim_limits_policy ON policy_claim_limits(policy_id);
CREATE INDEX idx_claim_limits_type ON policy_claim_limits(limit_type);


-- Claims-free bonus structure (no-claims discount accumulation)
CREATE TABLE no_claims_bonus (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    
    -- Current status
    consecutive_claim_free_years INTEGER DEFAULT 0,
    current_discount_percentage DECIMAL(5,2) DEFAULT 0,
    max_discount_percentage DECIMAL(5,2) DEFAULT 50.00, -- Industry standard cap
    
    -- Tier structure
    tier_1_threshold INTEGER DEFAULT 1, -- 1 year = 10%
    tier_1_discount DECIMAL(5,2) DEFAULT 10.00,
    tier_2_threshold INTEGER DEFAULT 2, -- 2 years = 20%
    tier_2_discount DECIMAL(5,2) DEFAULT 20.00,
    tier_3_threshold INTEGER DEFAULT 3, -- 3+ years = 30-50%
    tier_3_discount DECIMAL(5,2) DEFAULT 30.00,
    tier_4_threshold INTEGER DEFAULT 4,
    tier_4_discount DECIMAL(5,2) DEFAULT 40.00,
    tier_5_threshold INTEGER DEFAULT 5,
    tier_5_discount DECIMAL(5,2) DEFAULT 50.00,
    
    -- Claim forgiveness
    claim_forgiveness_used BOOLEAN DEFAULT FALSE,
    claim_forgiveness_available BOOLEAN DEFAULT TRUE,
    claim_forgiveness_reset_date DATE,
    
    -- History
    last_claim_date TIMESTAMPTZ,
    discount_reset_after_claim BOOLEAN DEFAULT FALSE, -- TRUE = reset to 0%, FALSE = step down one tier
    
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(policy_id)
);

CREATE INDEX idx_no_claims_policy ON no_claims_bonus(policy_id);

-- Usage-based insurance (UBI) safe device rewards
CREATE TABLE ubi_rewards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    
    -- Safety score (0-100)
    current_safety_score DECIMAL(5,2),
    score_calculation_date TIMESTAMPTZ,
    
    -- Score components
    drop_risk_score DECIMAL(5,2), -- Based on accelerometer data
    usage_pattern_score DECIMAL(5,2), -- Heavy vs light usage
    care_score DECIMAL(5,2), -- Case usage, screen protector, etc.
    
    -- Reward structure
    reward_percentage DECIMAL(5,2) DEFAULT 0, -- Premium rebate %
    reward_amount DECIMAL(28,8),
    currency CHAR(3) DEFAULT 'USD',
    
    -- Payout
    reward_paid BOOLEAN DEFAULT FALSE,
    paid_at TIMESTAMPTZ,
    movement_id UUID REFERENCES value_movements(id),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ubi_policy ON ubi_rewards(policy_id);

-- Per-device-type claim limits and waiting periods
CREATE TABLE device_type_claim_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_contract_id UUID NOT NULL REFERENCES product_contracts(id),
    
    device_category device_category NOT NULL,
    
    -- Screen claim sub-limits
    max_screen_claims_per_period INTEGER DEFAULT 2,
    screen_claim_period_months INTEGER DEFAULT 12,
    
    -- Theft claim waiting period (anti-fraud)
    theft_waiting_period_days INTEGER DEFAULT 14,
    
    -- Concurrent claim blocking
    concurrent_claims_blocked BOOLEAN DEFAULT TRUE,
    min_days_between_claims INTEGER DEFAULT 7,
    
    -- Seasonal adjustments
    holiday_theft_multiplier DECIMAL(5,2) DEFAULT 1.00, -- Reserve increase Nov-Dec
    back_to_school_damage_multiplier DECIMAL(5,2) DEFAULT 1.00, -- September spike
    
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    is_active BOOLEAN DEFAULT TRUE,
    
    UNIQUE(product_contract_id, device_category, valid_from)
);

CREATE INDEX idx_device_type_limits_contract ON device_type_claim_limits(product_contract_id);
CREATE INDEX idx_device_type_limits_category ON device_type_claim_limits(device_category) WHERE is_active = TRUE;

-- Discount validation rules
CREATE TABLE discount_validation_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_type VARCHAR(50) NOT NULL,
    
    -- Validation criteria
    max_discount_percentage DECIMAL(5,2),
    max_discount_amount DECIMAL(28,8),
    min_policy_term_months INTEGER,
    max_combined_discounts INTEGER,
    
    -- Exclusions
    excluded_product_codes VARCHAR(50)[],
    excluded_customer_types VARCHAR(50)[],
    
    -- Stacking rules
    can_stack_with_others BOOLEAN DEFAULT TRUE,
    stack_priority INTEGER DEFAULT 1, -- Application order
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_discount_rules_type ON discount_validation_rules(discount_type) WHERE is_active = TRUE;

-- Limit reset trigger function
CREATE OR REPLACE FUNCTION reset_annual_limits() RETURNS TRIGGER AS $$
BEGIN
    -- Reset annual limits when reset_date is reached
    UPDATE policy_claim_limits
    SET used_amount = 0,
        claim_count_used = 0,
        period_start = CURRENT_DATE,
        period_end = CURRENT_DATE + INTERVAL '1 year',
        reset_date = CURRENT_DATE + INTERVAL '1 year'
    WHERE limit_type = 'annual'
    AND reset_date <= CURRENT_DATE;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Note: Schedule this function to run daily via pg_cron or similar
