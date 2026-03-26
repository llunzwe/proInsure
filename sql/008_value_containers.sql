-- 008_value_containers.sql

-- Ledger accounts (value containers) - double entry accounting
CREATE TABLE value_containers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    type account_type NOT NULL,
    normal_balance VARCHAR(6) NOT NULL CHECK (normal_balance IN ('debit', 'credit')),
    
    -- Reference to COA
    coa_account_id UUID REFERENCES insurance_chart_of_accounts(id),
    
    -- Balance (derived, but stored for performance)
    current_balance DECIMAL(28,8) DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'closed', 'frozen')),
    frozen_at TIMESTAMPTZ,
    frozen_reason TEXT,
    
    -- For client money ring-fencing
    is_client_money_account BOOLEAN DEFAULT FALSE,
    related_party_id UUID REFERENCES parties(id), -- For segregated client accounts
    
    -- Bitemporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    
    CONSTRAINT valid_balance CHECK (type != 'ASSET' OR current_balance >= 0) -- Assets can't be negative (configurable)
);

-- Container balances history (monthly snapshots)
CREATE TABLE container_balance_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    container_id UUID NOT NULL REFERENCES value_containers(id),
    snapshot_date DATE NOT NULL,
    opening_balance DECIMAL(28,8),
    closing_balance DECIMAL(28,8),
    total_debits DECIMAL(28,8) DEFAULT 0,
    total_credits DECIMAL(28,8) DEFAULT 0,
    movement_count INTEGER DEFAULT 0,
    UNIQUE(container_id, snapshot_date)
);

-- Indexes
CREATE INDEX idx_containers_code ON value_containers(code);
CREATE INDEX idx_containers_coa ON value_containers(coa_account_id);
CREATE INDEX idx_containers_party ON value_containers(related_party_id) WHERE is_client_money_account = TRUE;
CREATE INDEX idx_containers_current ON value_containers(is_current) WHERE is_current = TRUE;

-- Trigger
CREATE TRIGGER trg_containers_hash BEFORE INSERT OR UPDATE ON value_containers
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();
