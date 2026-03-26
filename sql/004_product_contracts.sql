-- 004_product_contracts.sql

-- Immutable insurance product definitions
CREATE TABLE product_contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_code VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Coverage rules (flexible JSONB structure)
    coverage_rules JSONB NOT NULL DEFAULT '{}'::jsonb,
    -- Example: {"theft": {"limit": 1000, "deductible": 50, "waiting_days": 14}, 
    --           "damage": {"limit": 800, "deductible": 25, "excludes": ["water"]}}
    
    -- Pricing
    premium_rate DECIMAL(10,6) NOT NULL, -- Percentage of sum insured
    premium_frequency VARCHAR(20) DEFAULT 'monthly' CHECK (premium_frequency IN ('monthly', 'quarterly', 'annual', 'single')),
    min_premium_amount DECIMAL(28,8),
    max_premium_amount DECIMAL(28,8),
    
    -- Sum insured limits
    min_sum_insured DECIMAL(28,8) NOT NULL,
    max_sum_insured DECIMAL(28,8) NOT NULL,
    max_claim_amount DECIMAL(28,8), -- Per claim limit
    
    -- Temporal validity
    valid_from DATE NOT NULL,
    valid_to DATE NOT NULL DEFAULT 'infinity',
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Immutable identification
    contract_hash UUID NOT NULL, -- Deterministic UUID from rules content
    immutable_hash VARCHAR(64) NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    
    -- Constraints
    CONSTRAINT unique_contract_version UNIQUE (product_code, version),
    CONSTRAINT valid_dates CHECK (valid_from < valid_to),
    CONSTRAINT valid_limits CHECK (min_sum_insured < max_sum_insured)
);

-- Product exclusions and conditions
CREATE TABLE product_conditions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_contract_id UUID NOT NULL REFERENCES product_contracts(id),
    condition_type VARCHAR(50) NOT NULL, -- exclusion, waiting_period, requirement
    description TEXT NOT NULL,
    applies_to_claim_types claim_type[],
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to generate deterministic contract hash
CREATE OR REPLACE FUNCTION generate_contract_hash() RETURNS TRIGGER AS $$
BEGIN
    -- Create deterministic hash based on coverage rules and limits
    NEW.contract_hash := uuid_generate_v5(
        uuid_ns_oid(), 
        NEW.product_code || ':' || NEW.version || ':' || NEW.coverage_rules::text
    );
    NEW.immutable_hash := encode(digest(
        NEW.product_code || NEW.coverage_rules::text || NEW.premium_rate::text,
        'sha256'
    ), 'hex');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_product_contracts_hash BEFORE INSERT ON product_contracts
    FOR EACH ROW EXECUTE FUNCTION generate_contract_hash();

CREATE INDEX idx_product_contracts_code ON product_contracts(product_code);
CREATE INDEX idx_product_contracts_current ON product_contracts(is_current) WHERE is_current = TRUE;
CREATE INDEX idx_product_contracts_dates ON product_contracts(valid_from, valid_to);
