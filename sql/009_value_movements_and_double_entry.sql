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
    foreign_currency CHAR(3),
    foreign_amount DECIMAL(28,8),
    
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

-- Trigger to check 4-eyes approval for large amounts
CREATE OR REPLACE FUNCTION check_movement_approval() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.total_debits > 10000 AND NEW.approved_by IS NULL THEN
        RAISE EXCEPTION 'Movements over 10,000 require secondary approval';
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
