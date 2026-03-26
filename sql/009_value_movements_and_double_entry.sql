-- 009_value_movements_and_double_entry.sql

-- Double-entry transactions (movements)
CREATE TABLE value_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(100) NOT NULL UNIQUE,
    type movement_type NOT NULL,
    description TEXT,
    
    -- Dates
    entry_date DATE NOT NULL DEFAULT CURRENT_DATE,
    value_date DATE NOT NULL DEFAULT CURRENT_DATE, -- When value actually changes hands
    posted_at TIMESTAMPTZ,
    
    -- Status workflow
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'pending', 'posted', 'reversing', 'reversed', 'cancelled')),
    
    -- Totals (must balance)
    total_debits DECIMAL(28,8) NOT NULL DEFAULT 0,
    total_credits DECIMAL(28,8) NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Cross-currency if needed
    exchange_rate DECIMAL(18,8) DEFAULT 1.0,
    exchange_rate_id UUID, -- Reference to exchange_rates table
    foreign_currency CHAR(3),
    foreign_amount DECIMAL(28,8),
    
    -- ISO 20022 payment messaging (optional)
    uetr UUID, -- Unique End-to-End Transaction Reference
    end_to_end_id VARCHAR(35), -- End-to-End Identification
    payment_instruction_id VARCHAR(35), -- Payment Instruction ID
    
    -- Link to business entity
    correlation_id UUID, -- Generic link to policy, claim, etc.
    correlation_type VARCHAR(50), -- 'policy', 'claim', 'premium_payment'
    
    -- Reversal tracking
    reversed_by UUID REFERENCES value_movements(id),
    reversal_reason TEXT,
    
    -- Bitemporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    approved_by UUID REFERENCES parties(id), -- 4-eyes
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    -- Critical constraint: debits must equal credits
    CONSTRAINT balanced_movement CHECK (total_debits = total_credits),
    CONSTRAINT positive_amounts CHECK (total_debits >= 0 AND total_credits >= 0)
);

-- Configuration table for approval thresholds
CREATE TABLE approval_thresholds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    config_name VARCHAR(50) NOT NULL UNIQUE,
    movement_type movement_type,
    threshold_amount DECIMAL(28,8) NOT NULL DEFAULT 5000,
    currency CHAR(3) DEFAULT 'USD',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Insert default thresholds
INSERT INTO approval_thresholds (config_name, movement_type, threshold_amount) VALUES
('default', NULL, 5000),
('premium_refund', 'refund', 1000),
('claim_payout', 'claim_payout', 5000),
('commission_payment', 'commission', 2500);

-- Trigger to check 4-eyes approval for large amounts
CREATE OR REPLACE FUNCTION check_movement_approval() RETURNS TRIGGER AS $$
DECLARE
    v_threshold DECIMAL(28,8);
BEGIN
    -- Get threshold for this movement type
    SELECT threshold_amount INTO v_threshold
    FROM approval_thresholds
    WHERE (movement_type IS NULL OR movement_type = NEW.type)
    AND is_active = TRUE
    ORDER BY movement_type NULLS LAST
    LIMIT 1;
    
    -- Default if no config found
    IF v_threshold IS NULL THEN
        v_threshold := 5000;
    END IF;
    
    IF NEW.total_debits > v_threshold AND NEW.approved_by IS NULL THEN
        RAISE EXCEPTION 'Movements over % require secondary approval (4-eyes principle)', v_threshold;
    END IF;
    IF NEW.status = 'posted' AND NEW.posted_at IS NULL THEN
        NEW.posted_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_movement_approval BEFORE INSERT OR UPDATE ON value_movements
    FOR EACH ROW EXECUTE FUNCTION check_movement_approval();

CREATE INDEX idx_movements_reference ON value_movements(reference);
CREATE INDEX idx_movements_type ON value_movements(type);
CREATE INDEX idx_movements_status ON value_movements(status);
CREATE INDEX idx_movements_dates ON value_movements(entry_date, value_date);
CREATE INDEX idx_movements_correlation ON value_movements(correlation_id, correlation_type);
CREATE INDEX idx_movements_posted ON value_movements(posted_at) WHERE status = 'posted';

-- Trigger for hash
CREATE TRIGGER trg_movements_hash BEFORE INSERT OR UPDATE ON value_movements
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

-- Exchange rates table for multi-currency transactions
CREATE TABLE exchange_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_currency CHAR(3) NOT NULL,
    to_currency CHAR(3) NOT NULL,
    rate DECIMAL(18,8) NOT NULL,
    rate_date DATE NOT NULL,
    rate_type VARCHAR(20) DEFAULT 'spot' CHECK (rate_type IN ('spot', 'forward', 'average', 'month_end')),
    source VARCHAR(50), -- e.g., 'ECB', 'Reuters', 'CentralBank'
    source_reference VARCHAR(100),
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(from_currency, to_currency, rate_date, rate_type)
);

CREATE INDEX idx_exchange_rates_lookup ON exchange_rates(from_currency, to_currency, rate_date);
CREATE INDEX idx_exchange_rates_current ON exchange_rates(from_currency, to_currency) WHERE valid_to = 'infinity';

-- Immutable financial event log for quick regulatory reports
CREATE TABLE financial_event_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_type VARCHAR(50) NOT NULL, -- balance_change, movement_posted, etc.
    
    -- References
    container_id UUID REFERENCES value_containers(id),
    movement_id UUID REFERENCES value_movements(id),
    
    -- Balance information
    previous_balance DECIMAL(28,8),
    new_balance DECIMAL(28,8),
    change_amount DECIMAL(28,8),
    currency CHAR(3),
    
    -- Context
    correlation_id UUID,
    correlation_type VARCHAR(50),
    
    -- Integrity
    event_hash VARCHAR(64) NOT NULL,
    previous_event_hash VARCHAR(64), -- Chain of events
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fin_event_time ON financial_event_log(event_time);
CREATE INDEX idx_fin_event_container ON financial_event_log(container_id);
CREATE INDEX idx_fin_event_type ON financial_event_log(event_type);

-- Trigger to populate financial event log when movement is posted
CREATE OR REPLACE FUNCTION log_financial_event() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'posted' AND OLD.status != 'posted' THEN
        -- Get the previous event hash for chaining
        INSERT INTO financial_event_log (
            event_time, event_type, movement_id, correlation_id, correlation_type,
            change_amount, currency, event_hash, previous_event_hash
        )
        SELECT 
            NEW.posted_at,
            'movement_posted',
            NEW.id,
            NEW.correlation_id,
            NEW.correlation_type,
            NEW.total_debits,
            NEW.currency,
            encode(digest(NEW::text || CURRENT_TIMESTAMP::text, 'sha256'), 'hex'),
            (SELECT event_hash FROM financial_event_log ORDER BY event_time DESC LIMIT 1)
        FROM value_movements vm
        WHERE vm.id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_log_financial_event AFTER UPDATE ON value_movements
    FOR EACH ROW EXECUTE FUNCTION log_financial_event();
