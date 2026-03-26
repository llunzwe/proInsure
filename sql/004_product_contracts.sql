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
    premium_frequency VARCHAR(20) DEFAULT 'monthly' CHECK (premium_frequency IN ('monthly', 'quarterly', 'annual', 'single', 'daily', 'pay_per_day')),
    min_premium_amount DECIMAL(28,8),
    max_premium_amount DECIMAL(28,8),
    
    -- Coverage subtype granularity
    is_screen_only_plan BOOLEAN DEFAULT FALSE,
    covers_theft BOOLEAN DEFAULT TRUE,
    covers_loss BOOLEAN DEFAULT TRUE,
    theft_deductible DECIMAL(28,8), -- Different deductible for theft vs damage
    loss_deductible DECIMAL(28,8),
    
    -- International coverage zones
    coverage_zone VARCHAR(20) DEFAULT 'domestic' CHECK (coverage_zone IN ('domestic', 'regional', 'global')),
    excluded_countries CHAR(2)[], -- ISO 3166 country codes
    
    -- Accessory coverage
    max_accessories_value DECIMAL(28,8),
    accessories_included BOOLEAN DEFAULT FALSE,
    
    -- Sum insured limits
    min_sum_insured DECIMAL(28,8) NOT NULL,
    max_sum_insured DECIMAL(28,8) NOT NULL,
    max_claim_amount DECIMAL(28,8), -- Per claim limit
    
    -- IFRS 17 Compliance
    ifrs17_contract_type ifrs17_contract_type DEFAULT 'non_participating',
    ifrs17_contract_boundary VARCHAR(50), -- boundary definition
    discount_rate DECIMAL(10,6), -- For fulfillment cash flows
    risk_margin DECIMAL(28,8), -- Risk adjustment
    fulfilment_cash_flows JSONB, -- Expected future premiums, claims, expenses
    
    -- Reporting
    reporting_standard VARCHAR(20) DEFAULT 'IFRS17',
    
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
CREATE INDEX idx_product_contracts_ifrs17 ON product_contracts(ifrs17_contract_type);

-- Normalized coverage items table (alternative to JSONB coverage_rules)
CREATE TABLE coverage_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_contract_id UUID NOT NULL REFERENCES product_contracts(id) ON DELETE CASCADE,
    coverage_type claim_type NOT NULL,
    coverage_name VARCHAR(100) NOT NULL,
    description TEXT,
    sum_insured_limit DECIMAL(28,8),
    deductible DECIMAL(28,8) DEFAULT 0,
    waiting_period_days INTEGER DEFAULT 0,
    exclusions TEXT[],
    conditions TEXT[],
    is_active BOOLEAN DEFAULT TRUE,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coverage_items_contract ON coverage_items(product_contract_id);
CREATE INDEX idx_coverage_items_type ON coverage_items(coverage_type);


-- Warranty coordination clauses
CREATE TABLE product_warranty_coordination (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_contract_id UUID NOT NULL REFERENCES product_contracts(id) ON DELETE CASCADE,
    
    -- OEM warranty handling
    oem_warranty_exhaustion_required BOOLEAN DEFAULT TRUE, -- Insurance pays only after OEM warranty expires
    oem_warranty_exception_accidental BOOLEAN DEFAULT TRUE, -- Except for accidental damage
    
    -- Overlapping coverage programs
    handles_applecare BOOLEAN DEFAULT FALSE,
    handles_samsung_care BOOLEAN DEFAULT FALSE,
    handles_retailer_extended_warranty BOOLEAN DEFAULT FALSE,
    
    -- Pro-rata refund calculations
    overlapping_coverage_refund_method VARCHAR(50) DEFAULT 'pro_rata' CHECK (overlapping_coverage_refund_method IN ('pro_rata', 'full_premium', 'none')),
    
    -- Coordination rules
    is_primary_payer BOOLEAN DEFAULT FALSE, -- TRUE = insurance pays first, FALSE = OEM pays first
    coordination_delay_days INTEGER DEFAULT 0, -- Wait period for OEM to respond
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_warranty_coord_contract ON product_warranty_coordination(product_contract_id);

-- Dynamic pricing factors configuration
CREATE TABLE pricing_factors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_contract_id UUID NOT NULL REFERENCES product_contracts(id) ON DELETE CASCADE,
    
    factor_name VARCHAR(50) NOT NULL CHECK (factor_name IN (
        'device_age', 'device_condition', 'usage_pattern', 'payment_method_risk', 
        'ubi_score', 'customer_tenure', 'claim_history'
    )),
    
    factor_type VARCHAR(20) NOT NULL CHECK (factor_type IN ('multiplier', 'addition', 'exclusion')),
    
    -- Factor values
    min_value DECIMAL(10,6),
    max_value DECIMAL(10,6),
    default_value DECIMAL(10,6),
    
    -- Configuration
    applies_to_new_business BOOLEAN DEFAULT TRUE,
    applies_to_renewals BOOLEAN DEFAULT TRUE,
    applies_to_mid_term_changes BOOLEAN DEFAULT FALSE,
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pricing_factors_contract ON pricing_factors(product_contract_id);
CREATE INDEX idx_pricing_factors_name ON pricing_factors(factor_name) WHERE is_active = TRUE;

-- Payment method risk scoring
CREATE TABLE payment_method_risk_scores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_method VARCHAR(50) NOT NULL UNIQUE,
    
    risk_multiplier DECIMAL(5,2) DEFAULT 1.00, -- e.g., 1.10 = 10% higher premium
    default_failure_rate DECIMAL(5,2), -- Historical failure %
    
    -- Collection characteristics
    collection_days INTEGER, -- Days to collect
    settlement_cycle_days INTEGER, -- Partner settlement cycle
    
    is_active BOOLEAN DEFAULT TRUE
);

-- Insert default payment method risk scores
INSERT INTO payment_method_risk_scores (payment_method, risk_multiplier, collection_days) VALUES
('credit_card', 1.00, 1),
('debit_card', 1.00, 1),
('bank_transfer', 0.95, 3),
('mobile_money', 1.05, 1),
('carrier_billing', 1.10, 90),
('airtime_deduction', 1.15, 0),
('cash', 1.20, 7);
