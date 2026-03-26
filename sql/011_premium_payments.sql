-- 011_premium_payments.sql

-- Premium billing and collection
CREATE TABLE premium_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    
    -- Billing period
    period VARCHAR(20) NOT NULL, -- e.g., "2026-03"
    period_start_date DATE NOT NULL,
    period_end_date DATE NOT NULL,
    
    -- Amounts
    amount_due DECIMAL(28,8) NOT NULL,
    amount_paid DECIMAL(28,8) DEFAULT 0,
    amount_written_off DECIMAL(28,8) DEFAULT 0,
    currency CHAR(3) DEFAULT 'USD',
    
    -- Dates
    due_date DATE NOT NULL,
    paid_date TIMESTAMPTZ,
    last_reminder_sent TIMESTAMPTZ,
    
    -- Status
    status payment_status DEFAULT 'pending',
    
    -- Payment details
    payment_method VARCHAR(50), -- credit_card, bank_transfer, mobile_money, deduction
    payment_reference VARCHAR(100), -- External transaction ID
    payment_gateway_response JSONB,
    
    -- Link to accounting
    movement_id UUID REFERENCES value_movements(id), -- When paid
    
    -- Retry logic for failed payments
    retry_count INTEGER DEFAULT 0,
    next_retry_date TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT amount_check CHECK (amount_paid + amount_written_off <= amount_due)
);

-- Payment attempts log (for failures and retries)
CREATE TABLE premium_payment_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_id UUID NOT NULL REFERENCES premium_payments(id),
    attempted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    amount_attempted DECIMAL(28,8),
    status VARCHAR(20) NOT NULL, -- success, failed, declined
    failure_reason TEXT,
    gateway_response JSONB,
    ip_address INET,
    device_fingerprint VARCHAR(256)
);

-- Indexes
CREATE INDEX idx_premium_policy ON premium_payments(policy_id);
CREATE INDEX idx_premium_status ON premium_payments(status);
CREATE INDEX idx_premium_due_date ON premium_payments(due_date);
CREATE INDEX idx_premium_overdue ON premium_payments(due_date, status) WHERE status = 'overdue';
CREATE INDEX idx_premium_movement ON premium_payments(movement_id);
