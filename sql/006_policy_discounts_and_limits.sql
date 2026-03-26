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
