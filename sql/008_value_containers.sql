-- =============================================================================
-- FILENAME: 008_value_containers.sql
-- DESCRIPTION: Double-entry ledger accounts (value containers) with FCA CASS support
-- VERSION: 1.0.0
-- DEPENDENCIES: 000_extensions.sql, 007_insurance_chart_of_accounts.sql
-- =============================================================================
-- SECURITY CLASSIFICATION: CONFIDENTIAL
-- DATA SENSITIVITY: Financial balances - SOX and FCA regulated
-- =============================================================================
-- ISO/IEC COMPLIANCE:
--   - ISO/IEC 27001:2013 - Financial data protection
--   - FCA CASS 5 - Client Money rules
--   - FCA SYSC - Systems and controls
--   - SOX 404 - Internal controls over financial reporting
--   - IFRS 17 - Insurance contract accounting
-- =============================================================================
-- CHANGE LOG:
--   v1.0.0 (2026-03-26) - Initial release with FCA CASS compliance
-- =============================================================================

-- -----------------------------------------------------------------------------
-- SECTION 1: VALUE CONTAINERS (LEDGER ACCOUNTS)
-- PURPOSE: Double-entry bookkeeping accounts with balance tracking
-- STANDARD: Luca Pacioli double-entry method (1494)
-- -----------------------------------------------------------------------------

/**
 * TABLE: value_containers
 * DESCRIPTION: Individual ledger accounts for double-entry accounting
 * 
 * COMPLIANCE:
 *   - FCA CASS 5.5: Client money segregation
 *   - SOX 404: Balance verification and audit trail
 *   - IFRS 17: Contractual service margin tracking
 * 
 * SECURITY: All balance changes logged to immutable event store
 * ARCHITECTURE: Each container maps to one COA account, many containers per COA
 */
CREATE TABLE value_containers (
    -- Primary identification
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,   -- Human-readable account code
    name VARCHAR(100) NOT NULL,         -- Account description
    
    -- Account classification
    type account_type NOT NULL,         -- ASSET, LIABILITY, EQUITY, INCOME, EXPENSE
    normal_balance VARCHAR(6) NOT NULL 
        CHECK (normal_balance IN ('debit', 'credit')),
    
    -- Chart of Accounts mapping
    coa_account_id UUID REFERENCES insurance_chart_of_accounts(id),
    
    -- Current balance (derived but stored for performance)
    current_balance DECIMAL(28,8) DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'USD',  -- ISO 4217
    
    -- Multi-currency support
    functional_currency CHAR(3) DEFAULT 'USD',  -- Reporting currency
    foreign_balance DECIMAL(28,8),              -- Balance in foreign currency
    
    -- Account status
    status VARCHAR(20) DEFAULT 'active' 
        CHECK (status IN ('active', 'closed', 'frozen', 'suspended')),
    frozen_at TIMESTAMPTZ,
    frozen_reason TEXT,
    closed_at TIMESTAMPTZ,
    
    -- FCA CASS Client Money segregation
    is_client_money_account BOOLEAN DEFAULT FALSE,
    client_money_type VARCHAR(50) CHECK (client_money_type IN (
        'premium', 'claim', 'return', 'mixed'
    )),
    related_party_id UUID REFERENCES parties(id),  -- For segregated client accounts
    ring_fence_group_id UUID,                       -- Groups related client accounts
    
    -- Trust account linkage (external bank account)
    trust_account_id UUID REFERENCES trust_accounts(id),
    
    -- Bitemporal versioning
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Audit and integrity (SOX 404)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL,
    
    -- Constraints
    CONSTRAINT valid_balance CHECK (
        -- Assets normally have debit balances
        (type != 'ASSET' OR current_balance >= 0 OR status != 'active')
    )
);

-- Table documentation
COMMENT ON TABLE value_containers IS 
'Double-entry ledger accounts with FCA CASS client money segregation. SOX 404 controlled.';

COMMENT ON COLUMN value_containers.is_client_money_account IS 
'FCA CASS 5.5: TRUE if account holds client money requiring segregation.';

COMMENT ON COLUMN value_containers.ring_fence_group_id IS 
'Groups client money accounts for regulatory reporting and reconciliation.';

-- -----------------------------------------------------------------------------
-- SECTION 2: CONTAINER BALANCE HISTORY
-- PURPOSE: Point-in-time balance snapshots for reporting
-- -----------------------------------------------------------------------------

CREATE TABLE container_balance_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    container_id UUID NOT NULL REFERENCES value_containers(id),
    
    -- Snapshot details
    snapshot_date DATE NOT NULL,
    snapshot_type VARCHAR(20) DEFAULT 'daily' 
        CHECK (snapshot_type IN ('daily', 'month_end', 'quarter_end', 'year_end', 'ad_hoc')),
    
    -- Balance information
    opening_balance DECIMAL(28,8),
    closing_balance DECIMAL(28,8),
    total_debits DECIMAL(28,8) DEFAULT 0,
    total_credits DECIMAL(28,8) DEFAULT 0,
    movement_count INTEGER DEFAULT 0,
    
    -- Foreign currency
    foreign_opening_balance DECIMAL(28,8),
    foreign_closing_balance DECIMAL(28,8),
    exchange_rate DECIMAL(18,8),
    
    -- Reconciliation
    reconciled BOOLEAN DEFAULT FALSE,
    reconciled_at TIMESTAMPTZ,
    reconciled_by UUID REFERENCES parties(id),
    reconciliation_difference DECIMAL(28,8),
    
    UNIQUE(container_id, snapshot_date, snapshot_type)
);

COMMENT ON TABLE container_balance_history IS 
'Point-in-time balance snapshots for financial reporting and audit trails.';

-- -----------------------------------------------------------------------------
-- SECTION 3: TRUST ACCOUNTS (FCA CASS)
-- PURPOSE: External bank accounts for client money segregation
-- COMPLIANCE: FCA CASS 5.5 - Statutory trust
-- -----------------------------------------------------------------------------

CREATE TABLE trust_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Account identification
    account_name VARCHAR(100) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    bank_branch VARCHAR(100),
    bank_address TEXT,
    
    -- ISO identifiers
    bic_code VARCHAR(11),               -- ISO 9362:2022
    iban VARCHAR(34),                   -- ISO 13616
    sort_code VARCHAR(20),              -- UK-specific
    
    -- Currency (ISO 4217)
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Account classification
    account_type VARCHAR(50) NOT NULL 
        CHECK (account_type IN ('client_money', 'operational', 'reserve', 'reinsurance')),
    
    -- FCA CASS specific
    cass_designation VARCHAR(50),       -- Statutory trust designation
    fca_reference VARCHAR(100),         -- FCA registration reference
    
    -- Balance tracking
    current_balance DECIMAL(28,8) DEFAULT 0,
    last_reconciled_at TIMESTAMPTZ,
    last_reconciled_balance DECIMAL(28,8),
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' 
        CHECK (status IN ('active', 'suspended', 'closed', 'frozen')),
    opened_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMPTZ,
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL
);

COMMENT ON TABLE trust_accounts IS 
'FCA CASS 5.5 trust accounts for client money segregation. Statutory trust status.';

COMMENT ON COLUMN trust_accounts.cass_designation IS 
'FCA CASS statutory trust designation for client money accounts.';

-- -----------------------------------------------------------------------------
-- SECTION 4: CLIENT MONEY CALCULATIONS
-- PURPOSE: Daily client money reconciliation (FCA CASS 5.6)
-- COMPLIANCE: Daily calculation and reconciliation required by FCA
-- -----------------------------------------------------------------------------

CREATE TABLE client_money_calculations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    calculation_date DATE NOT NULL,
    trust_account_id UUID NOT NULL REFERENCES trust_accounts(id),
    
    -- Client money calculation (FCA CASS 5.6.2)
    calculated_client_money DECIMAL(28,8) NOT NULL,  -- From system records
    actual_trust_balance DECIMAL(28,8) NOT NULL,     -- From bank statement
    
    -- Difference analysis
    difference DECIMAL(28,8) NOT NULL,
    difference_explanation TEXT,
    
    -- Breakdown
    premium_holdings DECIMAL(28,8),
    claim_holdings DECIMAL(28,8),
    return_premium_holdings DECIMAL(28,8),
    
    -- Reconciliation status
    status VARCHAR(20) NOT NULL 
        CHECK (status IN ('balanced', 'shortfall', 'surplus', 'reconciled')),
    
    -- Reconciliation details
    reconciled_at TIMESTAMPTZ,
    reconciled_by UUID REFERENCES parties(id),
    reconciliation_notes TEXT,
    supporting_document_ids UUID[],
    
    -- Compliance (FCA CASS 5.6.6)
    compliance_status VARCHAR(20) DEFAULT 'compliant' 
        CHECK (compliance_status IN ('compliant', 'breach', 'investigating', 'resolved')),
    breach_amount DECIMAL(28,8),
    breach_notified_at TIMESTAMPTZ,     -- FCA notification timestamp
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    
    UNIQUE(calculation_date, trust_account_id)
);

COMMENT ON TABLE client_money_calculations IS 
'FCA CASS 5.6 daily client money calculations and reconciliation. Breach notification tracked.';

COMMENT ON COLUMN client_money_calculations.calculated_client_money IS 
'System-calculated client money per FCA CASS 5.6.2.';

COMMENT ON COLUMN client_money_calculations.breach_notified_at IS 
'Timestamp of FCA notification in case of CASS rule breach.';

-- -----------------------------------------------------------------------------
-- SECTION 5: INDEXES
-- -----------------------------------------------------------------------------

-- Container indexes
CREATE INDEX idx_containers_code ON value_containers(code);
CREATE INDEX idx_containers_coa ON value_containers(coa_account_id);
CREATE INDEX idx_containers_party ON value_containers(related_party_id) 
    WHERE is_client_money_account = TRUE;
CREATE INDEX idx_containers_trust ON value_containers(trust_account_id) 
    WHERE trust_account_id IS NOT NULL;
CREATE INDEX idx_containers_current ON value_containers(is_current) 
    WHERE is_current = TRUE;

-- Balance history indexes
CREATE INDEX idx_balance_history_container ON container_balance_history(container_id);
CREATE INDEX idx_balance_history_date ON container_balance_history(snapshot_date);
CREATE INDEX idx_balance_history_type ON container_balance_history(snapshot_type);

-- Trust account indexes
CREATE INDEX idx_trust_accounts_status ON trust_accounts(status);
CREATE INDEX idx_trust_accounts_type ON trust_accounts(account_type);

-- Client money indexes
CREATE INDEX idx_client_money_calc_date ON client_money_calculations(calculation_date);
CREATE INDEX idx_client_money_calc_trust ON client_money_calculations(trust_account_id);
CREATE INDEX idx_client_money_calc_status ON client_money_calculations(status);
CREATE INDEX idx_client_money_compliance ON client_money_calculations(compliance_status) 
    WHERE compliance_status != 'compliant';

-- -----------------------------------------------------------------------------
-- SECTION 6: TRIGGERS
-- -----------------------------------------------------------------------------

CREATE TRIGGER trg_containers_hash 
    BEFORE INSERT OR UPDATE ON value_containers
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_trust_accounts_hash 
    BEFORE INSERT OR UPDATE ON trust_accounts
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

-- Prevent deletion of client money accounts with balances (FCA CASS)
CREATE OR REPLACE FUNCTION prevent_client_money_deletion()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.is_client_money_account AND OLD.current_balance != 0 THEN
        RAISE EXCEPTION 'FCA CASS: Cannot delete client money account with non-zero balance. Transfer funds first.';
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_cm_deletion
    BEFORE DELETE ON value_containers
    FOR EACH ROW EXECUTE FUNCTION prevent_client_money_deletion();

-- =============================================================================
-- END OF FILE: 008_value_containers.sql
-- =============================================================================
