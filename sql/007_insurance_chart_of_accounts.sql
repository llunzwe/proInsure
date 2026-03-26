-- =============================================================================
-- FILENAME: 007_insurance_chart_of_accounts.sql
-- DESCRIPTION: Chart of Accounts per IFRS 17 Insurance Contracts Standard
-- VERSION: 1.0.0
-- DEPENDENCIES: 000_extensions.sql
-- =============================================================================
-- SECURITY CLASSIFICATION: INTERNAL
-- DATA SENSITIVITY: Financial reporting data - SOX relevant
-- =============================================================================
-- ISO/IEC COMPLIANCE:
--   - IFRS 17:2020 - Insurance Contracts (primary standard)
--   - IFRS 9:2014 - Financial Instruments (credit loss)
--   - IFRS 15:2014 - Revenue from Contracts with Customers
--   - ISO 4217:2015 - Currency codes
--   - SOX 404 - Internal controls over financial reporting
-- =============================================================================
-- ACCOUNTING STANDARDS:
--   - Double-entry bookkeeping (Luca Pacioli method, 1494)
--   - Accrual basis accounting
--   - Bitemporal financial reporting
-- =============================================================================
-- CHANGE LOG:
--   v1.0.0 (2026-03-26) - IFRS 17 compliant COA structure
-- =============================================================================

-- -----------------------------------------------------------------------------
-- SECTION 1: CHART OF ACCOUNTS
-- PURPOSE: Hierarchical account structure for insurance financial reporting
-- STANDARD: IFRS 17 Insurance Contracts
-- -----------------------------------------------------------------------------

/**
 * TABLE: insurance_chart_of_accounts
 * DESCRIPTION: Master chart of accounts aligned with IFRS 17 portfolio groups
 *              and insurance contract classification
 * 
 * COMPLIANCE:
 *   - IFRS 17.19: Grouping of insurance contracts
 *   - IFRS 17.32: Separation of components (insurance, investment, service)
 *   - SOX 404: Account-level audit trail and approval
 * 
 * HIERARCHY: Account codes follow logical grouping:
 *   1xxx - Assets (Current and non-current)
 *   2xxx - Liabilities (Claims reserves, UPR, payables)
 *   3xxx - Equity (Retained earnings, reserves)
 *   4xxx - Income (Premium, investment, commission)
 *   5xxx - Expenses (Claims, acquisition, maintenance)
 */
CREATE TABLE insurance_chart_of_accounts (
    -- Primary identification
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Account classification (SOX 404 control point)
    account_code VARCHAR(50) NOT NULL UNIQUE,
    account_name VARCHAR(100) NOT NULL,
    account_type account_type NOT NULL,
    
    -- Insurance-specific classification (IFRS 17.19)
    insurance_category VARCHAR(50) CHECK (insurance_category IN (
        -- Revenue categories
        'premium_income',           -- Gross written premium
        'reinsurance_commission',   -- Commission from reinsurers
        
        -- Expense categories
        'claim_expense',            -- Incurred claims (paid + reserves)
        'acquisition_costs',        -- Deferred acquisition costs (DAC)
        'maintenance_costs',        -- Policy administration
        'reinsurance_premium',      -- Ceded to reinsurers
        
        -- Reserve categories (IFRS 17.19)
        'reserve_reported',         -- Case reserves (known claims)
        'reserve_ibnr',             -- Incurred But Not Reported
        'reserve_uer',              -- Unearned Premium Reserve
        'reserve_uneamed',          -- Alternative UPR designation
        'reserve_lae',              -- Loss Adjustment Expenses
        
        -- Recovery categories
        'reinsurance_recovery',     -- Reinsurance recoverables
        'salvage_recovery',         -- Salvage/subrogation income
        'subrogation_recovery',     -- Third-party recoveries
        
        -- Operational categories
        'replacement_device_cost',  -- Device replacement inventory
        'commission_expense',       -- Agent/broker commissions
        'client_money',             -- FCA CASS ring-fenced funds
        
        -- Balance sheet
        'bank',                     -- Bank accounts
        'cash',                     -- Cash on hand
        'suspense'                  -- Unallocated items
    )),
    
    -- Hierarchical structure
    parent_account_id UUID REFERENCES insurance_chart_of_accounts(id),
    account_level INTEGER DEFAULT 1,    -- 1=Summary, 2=Intermediate, 3=Detail
    
    -- Double-entry configuration
    normal_balance VARCHAR(6) CHECK (normal_balance IN ('debit', 'credit')),
    
    -- Documentation
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- IFRS 17 reporting (IFRS 17.106-109)
    ifrs17_portfolio VARCHAR(50),       -- Portfolio grouping
    ifrs17_group VARCHAR(50),           -- Profitability group
    ifrs17_measurement_model VARCHAR(20) CHECK (ifrs17_measurement_model IN (
        'GMM',      -- General Measurement Model
        'PAA',      -- Premium Allocation Approach
        'VFA'       -- Variable Fee Approach
    )) DEFAULT 'PAA',
    
    -- Tax and regulatory reporting
    tax_code VARCHAR(20),               -- Tax authority code
    reporting_standard VARCHAR(20) DEFAULT 'IFRS17' 
        CHECK (reporting_standard IN ('IFRS17', 'IFRS9', 'GAAP', 'STAT')),
    statutory_account_code VARCHAR(50), -- Local regulatory code
    
    -- Bitemporal validity
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    -- Audit and integrity (SOX 404)
    immutable_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    approved_by UUID REFERENCES parties(id),  -- SOX control: Account creation approval
    
    -- Constraints
    CONSTRAINT valid_account_hierarchy CHECK (
        (parent_account_id IS NULL AND account_level = 1) OR
        (parent_account_id IS NOT NULL AND account_level > 1)
    )
);

-- Table documentation
COMMENT ON TABLE insurance_chart_of_accounts IS 
'IFRS 17 compliant Chart of Accounts with portfolio grouping. SOX 404 controlled.';

COMMENT ON COLUMN insurance_chart_of_accounts.insurance_category IS 
'IFRS 17 insurance contract classification for financial statement presentation.';

COMMENT ON COLUMN insurance_chart_of_accounts.ifrs17_measurement_model IS 
'IFRS 17 measurement model: GMM (General), PAA (Premium Allocation), VFA (Variable Fee).';

COMMENT ON COLUMN insurance_chart_of_accounts.reporting_standard IS 
'Financial reporting framework: IFRS17, IFRS9, GAAP (US), STAT (statutory).';

-- -----------------------------------------------------------------------------
-- SECTION 2: STANDARD INSURANCE COA
-- PURPOSE: Pre-populated industry-standard account structure
-- -----------------------------------------------------------------------------

INSERT INTO insurance_chart_of_accounts (
    account_code, account_name, account_type, insurance_category, 
    normal_balance, ifrs17_measurement_model, immutable_hash
) VALUES
-- ============================================================================
-- ASSETS (1xxx) - Economic resources controlled by the entity
-- ============================================================================

-- Current Assets
('1000', 'Cash and Bank', 'ASSET', 'bank', 'debit', 'PAA', 'init'),
('1010', 'Cash on Hand', 'ASSET', 'cash', 'debit', 'PAA', 'init'),
('1020', 'Bank - Operating Account', 'ASSET', 'bank', 'debit', 'PAA', 'init'),
('1030', 'Bank - Client Money Account (FCA CASS)', 'ASSET', 'client_money', 'debit', 'PAA', 'init'),

-- Receivables
('1100', 'Premium Receivable - Direct', 'ASSET', 'premium_income', 'debit', 'PAA', 'init'),
('1110', 'Premium Receivable - Agents', 'ASSET', 'premium_income', 'debit', 'PAA', 'init'),
('1120', 'Allowance for Doubtful Premiums (ECL)', 'ASSET', 'premium_income', 'credit', 'PAA', 'init'),

-- Reinsurance Assets
('1200', 'Reinsurance Recoverables - Reported Claims', 'ASSET', 'reinsurance_recovery', 'debit', 'GMM', 'init'),
('1210', 'Reinsurance Recoverables - IBNR', 'ASSET', 'reinsurance_recovery', 'debit', 'GMM', 'init'),
('1220', 'Reinsurance Commission Receivable', 'ASSET', 'reinsurance_commission', 'debit', 'PAA', 'init'),

-- Inventory and Recoveries
('1300', 'Salvage Inventory - Devices Pending Sale', 'ASSET', 'salvage_recovery', 'debit', 'PAA', 'init'),
('1310', 'Subrogation Receivables', 'ASSET', 'subrogation_recovery', 'debit', 'PAA', 'init'),
('1400', 'Replacement Device Inventory', 'ASSET', 'replacement_device_cost', 'debit', 'PAA', 'init'),

-- Prepayments and Deferred Costs
('1500', 'Deferred Acquisition Costs (DAC)', 'ASSET', 'acquisition_costs', 'debit', 'PAA', 'init'),
('1510', 'Prepaid Reinsurance Premiums', 'ASSET', 'reinsurance_premium', 'debit', 'PAA', 'init'),

-- Suspense and Clearing
('1900', 'Suspense Account - Unallocated Receipts', 'ASSET', 'suspense', 'debit', 'PAA', 'init'),

-- ============================================================================
-- LIABILITIES (2xxx) - Present obligations from past events
-- ============================================================================

-- Claim Reserves (IFRS 17.19 - Insurance contract liabilities)
('2000', 'Claim Reserves - Reported (Case)', 'LIABILITY', 'reserve_reported', 'credit', 'GMM', 'init'),
('2010', 'Claim Reserves - IBNR', 'LIABILITY', 'reserve_ibnr', 'credit', 'GMM', 'init'),
('2020', 'Claim Reserves - Loss Adjustment Expenses', 'LIABILITY', 'reserve_lae', 'credit', 'GMM', 'init'),
('2030', 'Claim Reserves - Reinsurance Recovered', 'LIABILITY', 'reinsurance_recovery', 'debit', 'GMM', 'init'),

-- Unearned Premium Reserves
('2100', 'Unearned Premium Reserve (UPR)', 'LIABILITY', 'reserve_uer', 'credit', 'PAA', 'init'),
('2110', 'UPR - Unearned Risk Reserve', 'LIABILITY', 'reserve_uer', 'credit', 'PAA', 'init'),

-- Payables
('2200', 'Outstanding Claims Payable', 'LIABILITY', 'claim_expense', 'credit', 'PAA', 'init'),
('2210', 'Commissions Payable - Agents', 'LIABILITY', 'commission_expense', 'credit', 'PAA', 'init'),
('2220', 'Commissions Payable - Brokers', 'LIABILITY', 'commission_expense', 'credit', 'PAA', 'init'),

-- Client Money (FCA CASS)
('2300', 'Client Money Held - Premiums', 'LIABILITY', 'client_money', 'credit', 'PAA', 'init'),
('2310', 'Client Money Held - Claims', 'LIABILITY', 'client_money', 'credit', 'PAA', 'init'),

-- Other Liabilities
('2400', 'Premium Tax Payable', 'LIABILITY', 'tax_code', 'credit', 'PAA', 'init'),
('2410', 'VAT/GST Payable', 'LIABILITY', 'tax_code', 'credit', 'PAA', 'init'),

-- ============================================================================
-- EQUITY (3xxx) - Residual interest in assets after deducting liabilities
-- ============================================================================
('3000', 'Retained Earnings', 'EQUITY', NULL, 'credit', 'PAA', 'init'),
('3100', 'Insurance Risk Reserve', 'EQUITY', NULL, 'credit', 'PAA', 'init'),

-- ============================================================================
-- INCOME (4xxx) - Increases in economic benefits (IFRS 15/17)
-- ============================================================================

-- Premium Income
('4000', 'Gross Premium Income - Written', 'INCOME', 'premium_income', 'credit', 'PAA', 'init'),
('4010', 'Gross Premium Income - Earned', 'INCOME', 'premium_income', 'credit', 'PAA', 'init'),
('4020', 'Premium Adjustments - Endorsements', 'INCOME', 'premium_income', 'credit', 'PAA', 'init'),

-- Reinsurance Income
('4100', 'Reinsurance Commission Income', 'INCOME', 'reinsurance_commission', 'credit', 'PAA', 'init'),
('4110', 'Reinsurance Profit Commission', 'INCOME', 'reinsurance_commission', 'credit', 'GMM', 'init'),

-- Recovery Income
('4200', 'Salvage Recovery Income', 'INCOME', 'salvage_recovery', 'credit', 'PAA', 'init'),
('4210', 'Subrogation Recovery Income', 'INCOME', 'subrogation_recovery', 'credit', 'PAA', 'init'),
('4220', 'Salvage Parts Income', 'INCOME', 'salvage_recovery', 'credit', 'PAA', 'init'),

-- Other Income
('4300', 'Investment Income', 'INCOME', NULL, 'credit', 'PAA', 'init'),
('4310', 'Foreign Exchange Gains', 'INCOME', NULL, 'credit', 'PAA', 'init'),

-- ============================================================================
-- EXPENSES (5xxx) - Decreases in economic benefits
-- ============================================================================

-- Claims and Benefits
('5000', 'Gross Claims Paid', 'EXPENSE', 'claim_expense', 'debit', 'PAA', 'init'),
('5010', 'Claim Reserves - Movement', 'EXPENSE', 'claim_expense', 'debit', 'GMM', 'init'),
('5020', 'Loss Adjustment Expenses (ALAE)', 'EXPENSE', 'claim_expense', 'debit', 'PAA', 'init'),
('5030', 'Defense and Cost Containment (DCC)', 'EXPENSE', 'claim_expense', 'debit', 'PAA', 'init'),

-- Acquisition Costs
('5100', 'Commission Expense - Initial', 'EXPENSE', 'commission_expense', 'debit', 'PAA', 'init'),
('5110', 'Commission Expense - Renewal', 'EXPENSE', 'commission_expense', 'debit', 'PAA', 'init'),
('5120', 'Commission Expense - Override', 'EXPENSE', 'commission_expense', 'debit', 'PAA', 'init'),
('5130', 'Brokerage Expense', 'EXPENSE', 'commission_expense', 'debit', 'PAA', 'init'),
('5140', 'Underwriting Costs - Other', 'EXPENSE', 'acquisition_costs', 'debit', 'PAA', 'init'),

-- Maintenance and Operations
('5200', 'Policy Maintenance Costs', 'EXPENSE', 'maintenance_costs', 'debit', 'PAA', 'init'),
('5210', 'Customer Service Costs', 'EXPENSE', 'maintenance_costs', 'debit', 'PAA', 'init'),
('5220', 'Technology and Systems', 'EXPENSE', 'maintenance_costs', 'debit', 'PAA', 'init'),

-- Reinsurance Costs
('5300', 'Reinsurance Premium Ceded', 'EXPENSE', 'reinsurance_premium', 'debit', 'GMM', 'init'),
('5310', 'Reinsurance Commission - Return', 'EXPENSE', 'reinsurance_premium', 'debit', 'GMM', 'init'),

-- Device Operations
('5400', 'Replacement Device Cost', 'EXPENSE', 'replacement_device_cost', 'debit', 'PAA', 'init'),
('5410', 'Repair Network Costs', 'EXPENSE', 'replacement_device_cost', 'debit', 'PAA', 'init'),
('5420', 'Shipping and Logistics', 'EXPENSE', 'replacement_device_cost', 'debit', 'PAA', 'init'),
('5430', 'Salvage Disposal Costs', 'EXPENSE', 'replacement_device_cost', 'debit', 'PAA', 'init'),

-- Other Expenses
('5500', 'Foreign Exchange Losses', 'EXPENSE', NULL, 'debit', 'PAA', 'init'),
('5510', 'Bad Debt Expense (ECL)', 'EXPENSE', NULL, 'debit', 'IFRS9', 'init'),
('5520', 'Suspense Write-offs', 'EXPENSE', NULL, 'debit', 'PAA', 'init');

-- -----------------------------------------------------------------------------
-- SECTION 3: INDEXES
-- -----------------------------------------------------------------------------

CREATE INDEX idx_coa_code ON insurance_chart_of_accounts(account_code);
CREATE INDEX idx_coa_type ON insurance_chart_of_accounts(account_type);
CREATE INDEX idx_coa_category ON insurance_chart_of_accounts(insurance_category);
CREATE INDEX idx_coa_current ON insurance_chart_of_accounts(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_coa_ifrs17 ON insurance_chart_of_accounts(ifrs17_measurement_model);
CREATE INDEX idx_coa_parent ON insurance_chart_of_accounts(parent_account_id);

-- -----------------------------------------------------------------------------
-- SECTION 4: TRIGGERS
-- -----------------------------------------------------------------------------

CREATE TRIGGER trg_coa_hash 
    BEFORE INSERT OR UPDATE ON insurance_chart_of_accounts
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

-- SOX 404: Require approval for new account creation
CREATE OR REPLACE FUNCTION check_coa_approval()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.created_by IS NOT NULL AND NEW.approved_by IS NULL THEN
        RAISE EXCEPTION 'SOX 404: New COA accounts require approval by authorized personnel';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_coa_approval
    BEFORE INSERT ON insurance_chart_of_accounts
    FOR EACH ROW EXECUTE FUNCTION check_coa_approval();

-- =============================================================================
-- END OF FILE: 007_insurance_chart_of_accounts.sql
-- =============================================================================
