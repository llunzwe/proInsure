-- =============================================================================
-- FILENAME: 016_settlements.sql
-- DESCRIPTION: Settlement finality, reinsurance, commissions, and suspense accounts
-- VERSION: 1.0.0
-- DEPENDENCIES: 000_extensions.sql, 008_value_containers.sql, 013_claims.sql
-- =============================================================================
-- SECURITY CLASSIFICATION: CONFIDENTIAL
-- DATA SENSITIVITY: Financial settlements - SOX and PCI DSS scope
-- =============================================================================
-- ISO/IEC COMPLIANCE:
--   - ISO/IEC 27001:2013 - Financial controls
--   - ISO 20022:2023 - Financial messaging (payments)
--   - PCI DSS v4.0 - Payment card data security
--   - FCA CASS 5 - Client money rules
--   - IFRS 17 - Insurance contract settlements
-- =============================================================================
-- CHANGE LOG:
--   v1.0.0 (2026-03-26) - Initial release with full compliance
-- =============================================================================

-- -----------------------------------------------------------------------------
-- SECTION 1: SETTLEMENTS (PAYOUT FINALITY)
-- PURPOSE: Track claim settlements with payment finality
-- STANDARD: ISO 20022 payment messages (pacs.008, camt.054)
-- -----------------------------------------------------------------------------

/**
 * TABLE: settlements
 * DESCRIPTION: Final settlement tracking for insurance claims
 * 
 * COMPLIANCE:
 *   - ISO 20022: Payment initiation (pain) and clearing (pacs) messages
 *   - PCI DSS: Cardholder data protection scope
 *   - SOX 404: Settlement approval and audit trail
 * 
 * PAYMENT FINALITY: Tracks progression from pending → provisional → final
 */
CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- References
    claim_id UUID NOT NULL REFERENCES claims(id),
    movement_id UUID NOT NULL REFERENCES value_movements(id),
    
    -- Settlement details (ISO 20022 aligned)
    method settlement_method NOT NULL,
    amount DECIMAL(28,8) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',  -- ISO 4217
    fees DECIMAL(28,8) DEFAULT 0,              -- Transaction fees
    net_amount DECIMAL(28,8) GENERATED ALWAYS AS (amount - fees) STORED,
    
    -- ISO 20022 payment identifiers
    uetr UUID,                                  -- Unique End-to-End Transaction Reference
    end_to_end_id VARCHAR(35),                  -- End-to-End Identification
    payment_instruction_id VARCHAR(35),         -- Payment Instruction ID
    
    -- Payee information (PCI DSS scope if card payment)
    payee_party_id UUID REFERENCES parties(id),
    payee_account_details_encrypted BYTEA,      -- Bank account, mobile wallet (AES-256)
    payee_account_token VARCHAR(100),           -- Tokenized reference (PCI DSS)
    payee_reference VARCHAR(100),               -- Customer reference
    
    -- Status progression (payment finality)
    status VARCHAR(20) DEFAULT 'pending' 
        CHECK (status IN ('pending', 'provisional', 'final', 'failed', 'recalled', 'reversed')),
    
    -- Finality timestamps
    provisional_at TIMESTAMPTZ,                 -- Funds reserved
    final_at TIMESTAMPTZ,                       -- Irrevocable settlement
    expected_finality_at TIMESTAMPTZ,           -- Predicted finality
    
    -- External tracking
    external_reference VARCHAR(100),            -- Bank transaction ID
    gateway_reference VARCHAR(100),             -- Payment gateway reference
    gateway_response JSONB,                     -- Full gateway response
    
    -- Proof and reconciliation
    proof_document_id UUID,                     -- Proof of payment document
    settlement_proof_encrypted BYTEA,           -- Encrypted proof document
    reconciliation_date DATE,
    reconciled_by UUID REFERENCES parties(id),
    
    -- Replacement device (in-kind settlement)
    settlement_type VARCHAR(20) DEFAULT 'cash' 
        CHECK (settlement_type IN ('cash', 'replacement_device', 'repair_voucher', 'store_credit')),
    replacement_device_id UUID,                 -- Reference to inventory
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL
);

COMMENT ON TABLE settlements IS 
'Claim settlements with ISO 20022 payment identifiers and finality tracking. PCI DSS scope for card payments.';

COMMENT ON COLUMN settlements.uetr IS 
'ISO 20022 Unique End-to-End Transaction Reference (UUID format).';

COMMENT ON COLUMN settlements.payee_account_details_encrypted IS 
'PCI DSS: AES-256 encrypted account details. Tokenization preferred.';

-- -----------------------------------------------------------------------------
-- SECTION 2: SETTLEMENT RECONCILIATION
-- PURPOSE: Match settlements to bank statements
-- -----------------------------------------------------------------------------

CREATE TABLE settlement_reconciliation (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    settlement_id UUID NOT NULL REFERENCES settlements(id),
    
    -- Batch processing
    reconciliation_batch_id UUID,
    statement_date DATE,
    
    -- Matching
    expected_amount DECIMAL(28,8),
    actual_amount DECIMAL(28,8),
    difference DECIMAL(28,8),
    
    -- Status
    status VARCHAR(20) CHECK (status IN ('matched', 'unmatched', 'discrepancy', 'resolved')),
    match_confidence DECIMAL(5,2),              -- AI matching score
    
    -- Resolution
    resolution_notes TEXT,
    resolved_by UUID REFERENCES parties(id),
    resolved_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE settlement_reconciliation IS 
'Automated settlement reconciliation against bank statements.';

-- -----------------------------------------------------------------------------
-- SECTION 3: REINSURANCE TREATIES
-- PURPOSE: Reinsurance contract management
-- STANDARD: CEIA/ACORD reinsurance standards
-- -----------------------------------------------------------------------------

CREATE TABLE reinsurance_treaties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Treaty identification
    treaty_code VARCHAR(50) NOT NULL UNIQUE,
    treaty_name VARCHAR(100) NOT NULL,
    
    -- Parties
    reinsurer_party_id UUID NOT NULL REFERENCES parties(id),
    reinsurer_lei VARCHAR(20),                  -- ISO 17442
    
    -- Treaty structure (CEIA standard types)
    treaty_type reinsurance_treaty_type NOT NULL,
    
    -- Coverage period
    coverage_start_date DATE NOT NULL,
    coverage_end_date DATE NOT NULL,
    
    -- Financial terms
    cession_percentage DECIMAL(5,2),            -- For proportional (e.g., 30.00 = 30%)
    excess_amount DECIMAL(28,8),                -- For XOL: attachment point
    limit_amount DECIMAL(28,8),                 -- Maximum reinsurer liability
    aggregate_limit DECIMAL(28,8),              -- Annual aggregate
    
    -- Commission structure
    commission_percentage DECIMAL(5,2),         -- Ceding commission
    profit_commission_percentage DECIMAL(5,2),  -- Profit-based
    no_claims_bonus_percentage DECIMAL(5,2),    -- NCB on ceded premium
    
    -- Premium and claims
    minimum_premium DECIMAL(28,8),
    deposit_premium DECIMAL(28,8),
    
    -- Accounting
    premium_ceded_container_id UUID REFERENCES value_containers(id),
    recovery_receivable_container_id UUID REFERENCES value_containers(id),
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' 
        CHECK (status IN ('active', 'expired', 'cancelled', 'in_run_off')),
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL
);

COMMENT ON TABLE reinsurance_treaties IS 
'Reinsurance treaty management per CEIA/ACORD standards. Supports proportional and non-proportional structures.';

COMMENT ON COLUMN reinsurance_treaties.treaty_type IS 
'CEIA standard: proportional_quota, proportional_surplus, non_proportional_xol, facultative.';

-- -----------------------------------------------------------------------------
-- SECTION 4: REINSURANCE CEDED PREMIUMS
-- PURPOSE: Track premiums ceded to reinsurers
-- -----------------------------------------------------------------------------

CREATE TABLE reinsurance_ceded_premiums (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    treaty_id UUID NOT NULL REFERENCES reinsurance_treaties(id),
    policy_id UUID NOT NULL REFERENCES policies(id),
    
    -- Premium calculation
    original_premium DECIMAL(28,8) NOT NULL,
    ceded_percentage DECIMAL(5,2) NOT NULL,
    ceded_amount DECIMAL(28,8) NOT NULL,
    commission_amount DECIMAL(28,8),
    net_ceded DECIMAL(28,8),                    -- ceded_amount - commission_amount
    currency CHAR(3) DEFAULT 'USD',
    
    -- Accounting period
    accounting_period VARCHAR(20) NOT NULL,     -- e.g., "2026-Q1"
    
    -- Settlement
    movement_id UUID REFERENCES value_movements(id),
    status VARCHAR(20) DEFAULT 'accrual' 
        CHECK (status IN ('accrual', 'settled', 'recovered')),
    settled_at TIMESTAMPTZ,
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL
);

-- -----------------------------------------------------------------------------
-- SECTION 5: REINSURANCE RECOVERIES
-- PURPOSE: Track claim recoveries from reinsurers
-- -----------------------------------------------------------------------------

CREATE TABLE reinsurance_recoveries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    treaty_id UUID NOT NULL REFERENCES reinsurance_treaties(id),
    claim_id UUID NOT NULL REFERENCES claims(id),
    
    -- Recovery calculation
    original_claim_amount DECIMAL(28,8) NOT NULL,
    recoverable_percentage DECIMAL(5,2) NOT NULL,
    recoverable_amount DECIMAL(28,8) NOT NULL,
    deductible_applied DECIMAL(28,8) DEFAULT 0,
    net_recovery DECIMAL(28,8),
    currency CHAR(3) DEFAULT 'USD',
    
    -- Status
    status recovery_status DEFAULT 'pending',
    
    -- Timeline
    claimed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    proof_of_loss_submitted_at TIMESTAMPTZ,
    recovered_at TIMESTAMPTZ,
    
    -- Accounting
    movement_id UUID REFERENCES value_movements(id),
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL
);

-- -----------------------------------------------------------------------------
-- SECTION 6: COMMISSIONS
-- PURPOSE: Agent and broker commission tracking
-- STANDARD: IFRS 15 revenue recognition for commissions
-- -----------------------------------------------------------------------------

CREATE TABLE commissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    commission_reference VARCHAR(50) NOT NULL UNIQUE,
    
    -- Relationships
    agent_party_id UUID NOT NULL REFERENCES parties(id),
    policy_id UUID REFERENCES policies(id),
    premium_payment_id UUID REFERENCES premium_payments(id),
    
    -- Commission details (IFRS 15)
    commission_type commission_type NOT NULL,
    base_amount DECIMAL(28,8) NOT NULL,         -- Premium amount commission is based on
    commission_rate DECIMAL(10,6) NOT NULL,     -- e.g., 0.15 = 15%
    commission_amount DECIMAL(28,8) NOT NULL,
    currency CHAR(3) DEFAULT 'USD',
    
    -- Clawback provisions
    clawback_eligible BOOLEAN DEFAULT TRUE,
    clawback_period_months INTEGER DEFAULT 12,
    clawback_applied BOOLEAN DEFAULT FALSE,
    clawback_amount DECIMAL(28,8) DEFAULT 0,
    
    -- Accounting period
    accounting_period VARCHAR(20) NOT NULL,
    
    -- Status
    status VARCHAR(20) DEFAULT 'accrual' 
        CHECK (status IN ('accrual', 'approved', 'paid', 'reversed', 'clawed_back')),
    
    -- Payment
    movement_id UUID REFERENCES value_movements(id),
    paid_at TIMESTAMPTZ,
    paid_to_account_encrypted BYTEA,            -- Encrypted bank details
    
    -- Approval (4-eyes)
    approved_by UUID REFERENCES parties(id),
    approved_at TIMESTAMPTZ,
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL
);

COMMENT ON TABLE commissions IS 
'Agent/broker commission tracking per IFRS 15. Supports clawback provisions.';

-- -----------------------------------------------------------------------------
-- SECTION 7: SUSPENSE ITEMS
-- PURPOSE: Unallocated payments requiring manual reconciliation
-- STANDARD: FCA CASS 5.6 - Suspense account management
-- -----------------------------------------------------------------------------

CREATE TABLE suspense_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    suspense_reference VARCHAR(50) NOT NULL UNIQUE,
    
    -- Receipt details
    amount DECIMAL(28,8) NOT NULL,
    currency CHAR(3) NOT NULL,
    received_at TIMESTAMPTZ NOT NULL,
    
    -- Source information
    payment_method_type payment_method_type,
    payment_gateway VARCHAR(50),
    gateway_transaction_id VARCHAR(100),
    bank_reference VARCHAR(100),
    
    -- Remitter information (for bank transfers)
    remitter_name TEXT,
    remitter_account VARCHAR(100),
    remitter_bank VARCHAR(100),
    received_in_container_id UUID REFERENCES value_containers(id),
    
    -- Matching attempts
    attempted_policy_id UUID REFERENCES policies(id),
    attempted_party_id UUID REFERENCES parties(id),
    matching_confidence DECIMAL(5,2),           -- AI matching score
    matching_algorithm VARCHAR(50),
    
    -- Status workflow
    status suspense_status DEFAULT 'unmatched',
    
    -- Allocation
    allocated_to_policy_id UUID REFERENCES policies(id),
    allocated_to_premium_payment_id UUID REFERENCES premium_payments(id),
    allocated_amount DECIMAL(28,8),
    allocated_at TIMESTAMPTZ,
    allocated_by UUID REFERENCES parties(id),
    allocation_movement_id UUID REFERENCES value_movements(id),
    
    -- Refund
    refunded_amount DECIMAL(28,8),
    refunded_at TIMESTAMPTZ,
    refund_reason TEXT,
    refund_movement_id UUID REFERENCES value_movements(id),
    
    -- Investigation
    notes TEXT,
    investigation_notes TEXT,
    escalated BOOLEAN DEFAULT FALSE,
    escalated_at TIMESTAMPTZ,
    
    -- FCA CASS
    client_money_calculation_id UUID REFERENCES client_money_calculations(id),
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL
);

COMMENT ON TABLE suspense_items IS 
'FCA CASS 5.6 suspense account for unallocated receipts. Daily reconciliation required.';

-- -----------------------------------------------------------------------------
-- SECTION 8: INDEXES
-- -----------------------------------------------------------------------------

-- Settlement indexes
CREATE INDEX idx_settlements_claim ON settlements(claim_id);
CREATE INDEX idx_settlements_status ON settlements(status);
CREATE INDEX idx_settlements_finality ON settlements(final_at);
CREATE INDEX idx_settlements_uetr ON settlements(uetr) WHERE uetr IS NOT NULL;
CREATE INDEX idx_settlements_movement ON settlements(movement_id);

-- Reconciliation indexes
CREATE INDEX idx_settlement_recon_settlement ON settlement_reconciliation(settlement_id);
CREATE INDEX idx_settlement_recon_batch ON settlement_reconciliation(reconciliation_batch_id);

-- Reinsurance indexes
CREATE INDEX idx_reinsurance_treaties_reinsurer ON reinsurance_treaties(reinsurer_party_id);
CREATE INDEX idx_reinsurance_treaties_dates ON reinsurance_treaties(coverage_start_date, coverage_end_date);
CREATE INDEX idx_reinsurance_treaties_status ON reinsurance_treaties(status);

CREATE INDEX idx_reinsurance_ceded_treaty ON reinsurance_ceded_premiums(treaty_id);
CREATE INDEX idx_reinsurance_ceded_policy ON reinsurance_ceded_premiums(policy_id);
CREATE INDEX idx_reinsurance_ceded_period ON reinsurance_ceded_premiums(accounting_period);

CREATE INDEX idx_reinsurance_recoveries_treaty ON reinsurance_recoveries(treaty_id);
CREATE INDEX idx_reinsurance_recoveries_claim ON reinsurance_recoveries(claim_id);
CREATE INDEX idx_reinsurance_recoveries_status ON reinsurance_recoveries(status);

-- Commission indexes
CREATE INDEX idx_commissions_agent ON commissions(agent_party_id);
CREATE INDEX idx_commissions_policy ON commissions(policy_id);
CREATE INDEX idx_commissions_period ON commissions(accounting_period);
CREATE INDEX idx_commissions_status ON commissions(status);
CREATE INDEX idx_commissions_type ON commissions(commission_type);

-- Suspense indexes
CREATE INDEX idx_suspense_status ON suspense_items(status);
CREATE INDEX idx_suspense_received ON suspense_items(received_at);
CREATE INDEX idx_suspense_reference ON suspense_items(gateway_transaction_id, bank_reference);
CREATE INDEX idx_suspense_allocated ON suspense_items(allocated_to_policy_id) WHERE status = 'matched';
CREATE INDEX idx_suspense_escalated ON suspense_items(escalated) WHERE escalated = TRUE;

-- -----------------------------------------------------------------------------
-- SECTION 9: TRIGGERS
-- -----------------------------------------------------------------------------

CREATE TRIGGER trg_settlements_hash 
    BEFORE INSERT OR UPDATE ON settlements
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_reinsurance_treaties_hash 
    BEFORE INSERT OR UPDATE ON reinsurance_treaties
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_reinsurance_ceded_hash 
    BEFORE INSERT OR UPDATE ON reinsurance_ceded_premiums
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_reinsurance_recoveries_hash 
    BEFORE INSERT OR UPDATE ON reinsurance_recoveries
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_commissions_hash 
    BEFORE INSERT OR UPDATE ON commissions
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_suspense_items_hash 
    BEFORE INSERT OR UPDATE ON suspense_items
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

-- -----------------------------------------------------------------------------
-- SECTION 10: PAYMENT FINALITY FUNCTION
-- PURPOSE: Enforce payment finality rules
-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION check_settlement_finality()
RETURNS TRIGGER AS $$
BEGIN
    -- Once final, cannot be modified
    IF OLD.status = 'final' AND NEW.status != 'final' THEN
        RAISE EXCEPTION 'Payment finality: Cannot reverse final settlement. Use reversal movement.';
    END IF;
    
    -- Set finality timestamp
    IF NEW.status = 'final' AND OLD.status != 'final' THEN
        NEW.final_at = CURRENT_TIMESTAMP;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_settlement_finality
    BEFORE UPDATE ON settlements
    FOR EACH ROW EXECUTE FUNCTION check_settlement_finality();

-- =============================================================================
-- END OF FILE: 016_settlements.sql
-- =============================================================================
