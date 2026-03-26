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
    payment_method_type payment_method_type,
    payment_method_id UUID, -- Reference to payment_methods table
    payment_gateway VARCHAR(50), -- Stripe, M-Pesa, Paystack, etc.
    gateway_transaction_id VARCHAR(100),
    payment_reference VARCHAR(100), -- External transaction ID
    payment_gateway_response JSONB,
    payment_gateway_webhook_log JSONB,
    
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
CREATE INDEX idx_premium_gateway ON premium_payments(payment_gateway) WHERE payment_gateway IS NOT NULL;

-- Payment methods table for tokenized payment instrument details
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Method details
    method_type payment_method_type NOT NULL,
    display_name VARCHAR(100), -- e.g., "Visa ending in 4242"
    
    -- Tokenized details (PCI DSS compliance)
    token_encrypted BYTEA, -- Encrypted token from payment gateway
    token_reference VARCHAR(100), -- Gateway's token reference
    last_four_digits VARCHAR(4),
    expiry_month INTEGER,
    expiry_year INTEGER,
    
    -- Gateway info
    payment_gateway VARCHAR(50) NOT NULL,
    gateway_customer_id VARCHAR(100),
    
    -- Status
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'revoked', 'failed')),
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMPTZ,
    verified_at TIMESTAMPTZ,
    
    -- PCI DSS
    pci_compliance_verified BOOLEAN DEFAULT FALSE,
    pci_validation_date DATE
);

CREATE INDEX idx_payment_methods_party ON payment_methods(party_id);
CREATE INDEX idx_payment_methods_active ON payment_methods(party_id, is_active) WHERE is_active = TRUE;

-- Tax rates table for VAT/GST/Insurance Premium Tax
CREATE TABLE tax_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tax_code VARCHAR(20) NOT NULL UNIQUE,
    tax_name VARCHAR(100) NOT NULL,
    tax_type VARCHAR(50) NOT NULL, -- vat, gst, ipt, sales_tax
    rate DECIMAL(10,6) NOT NULL, -- e.g., 0.20 for 20%
    
    -- Jurisdiction
    country_code CHAR(2), -- ISO 3166
    region_code VARCHAR(10),
    
    -- Validity
    valid_from DATE NOT NULL,
    valid_to DATE NOT NULL DEFAULT 'infinity',
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    -- GL mapping
    liability_container_id UUID REFERENCES value_containers(id), -- For tax remittance
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tax_rates_lookup ON tax_rates(country_code, tax_type, valid_from) WHERE is_active = TRUE;

-- Link tax to premium payments
ALTER TABLE premium_payments ADD COLUMN IF NOT EXISTS tax_rate_id UUID REFERENCES tax_rates(id);
ALTER TABLE premium_payments ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(28,8) DEFAULT 0;
ALTER TABLE premium_payments ADD COLUMN IF NOT EXISTS gross_amount DECIMAL(28,8);

-- Billing schedule table for recurring payments
CREATE TABLE billing_schedule (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    
    -- Schedule details
    schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('initial', 'renewal', 'adjustment')),
    billing_frequency VARCHAR(20) NOT NULL, -- monthly, quarterly, annual
    
    -- Period
    period_number INTEGER NOT NULL,
    period_start_date DATE NOT NULL,
    period_end_date DATE NOT NULL,
    
    -- Amounts
    amount DECIMAL(28,8) NOT NULL,
    currency CHAR(3) DEFAULT 'USD',
    tax_amount DECIMAL(28,8) DEFAULT 0,
    total_amount DECIMAL(28,8) NOT NULL,
    
    -- Status
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'invoiced', 'paid', 'failed', 'cancelled')),
    
    -- Payment
    payment_id UUID REFERENCES premium_payments(id),
    payment_method_id UUID REFERENCES payment_methods(id),
    
    -- Due date
    due_date DATE NOT NULL,
    invoice_generated_at TIMESTAMPTZ,
    
    -- Retry logic
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    next_retry_date TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(policy_id, period_number)
);

CREATE INDEX idx_billing_schedule_policy ON billing_schedule(policy_id);
CREATE INDEX idx_billing_schedule_due ON billing_schedule(due_date, status) WHERE status IN ('scheduled', 'failed');
CREATE INDEX idx_billing_schedule_retry ON billing_schedule(next_retry_date) WHERE status = 'failed';

-- Retry configuration table
CREATE TABLE retry_configurations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    config_name VARCHAR(50) NOT NULL UNIQUE,
    retry_intervals INTEGER[] DEFAULT ARRAY[1, 3, 7], -- Days between retries
    max_retries INTEGER DEFAULT 3,
    escalate_after_retries BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE
);

-- Insert default retry configuration
INSERT INTO retry_configurations (config_name, retry_intervals, max_retries) VALUES 
('default_premium_retry', ARRAY[1, 3, 7], 3);

-- Premium refunds table (for cancellations, overpayments)
CREATE TABLE premium_refunds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    refund_reference VARCHAR(50) NOT NULL UNIQUE,
    
    -- Original payment reference
    original_payment_id UUID NOT NULL REFERENCES premium_payments(id),
    policy_id UUID NOT NULL REFERENCES policies(id),
    
    -- Refund details
    refund_type VARCHAR(50) NOT NULL CHECK (refund_type IN ('cooling_off', 'cancellation', 'overpayment', 'adjustment', 'claim_consequential')),
    refund_reason TEXT,
    
    -- Amounts
    original_amount DECIMAL(28,8) NOT NULL,
    refund_amount DECIMAL(28,8) NOT NULL,
    cancellation_fee DECIMAL(28,8) DEFAULT 0, -- Fee deducted from refund
    net_refund DECIMAL(28,8) NOT NULL, -- refund_amount - cancellation_fee
    currency CHAR(3) DEFAULT 'USD',
    
    -- Pro-rata calculation (for cancellations)
    days_used INTEGER,
    days_total INTEGER,
    pro_rata_percentage DECIMAL(5,2),
    
    -- Status workflow
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'processing', 'completed', 'failed')),
    
    -- Approval (4-eyes)
    requested_by UUID REFERENCES parties(id),
    requested_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    approved_by UUID REFERENCES parties(id),
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,
    
    -- Payment
    movement_id UUID REFERENCES value_movements(id),
    refunded_at TIMESTAMPTZ,
    refund_method_type payment_method_type,
    refund_destination_reference VARCHAR(100), -- Account, card reference, etc.
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL
);

CREATE TRIGGER trg_premium_refunds_hash BEFORE INSERT OR UPDATE ON premium_refunds
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE INDEX idx_premium_refunds_payment ON premium_refunds(original_payment_id);
CREATE INDEX idx_premium_refunds_policy ON premium_refunds(policy_id);
CREATE INDEX idx_premium_refunds_status ON premium_refunds(status);
CREATE INDEX idx_premium_refunds_type ON premium_refunds(refund_type);
