-- 007_insurance_chart_of_accounts.sql

-- Chart of accounts specific to insurance (IFRS 17 compliant)
CREATE TABLE insurance_chart_of_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_code VARCHAR(50) NOT NULL UNIQUE,
    account_name VARCHAR(100) NOT NULL,
    account_type account_type NOT NULL,
    
    -- Insurance specific classification
    insurance_category VARCHAR(50) CHECK (insurance_category IN (
        'premium_income', 'claim_expense', 'acquisition_costs', 'maintenance_costs',
        'reserve_reported', 'reserve_ibnr', 'reserve_uer', 'reserve_uneamed',
        'reinsurance_premium', 'reinsurance_recovery', 'commission_expense',
        'salvage_recovery', 'subrogation_recovery', 'replacement_device_cost',
        'client_money', 'bank', 'cash', 'suspense'
    )),
    
    parent_account_id UUID REFERENCES insurance_chart_of_accounts(id),
    normal_balance VARCHAR(6) CHECK (normal_balance IN ('debit', 'credit')),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- For statutory reporting
    ifrs17_portfolio VARCHAR(50),
    ifrs17_group VARCHAR(50),
    tax_code VARCHAR(20),
    
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Insert standard insurance COA
INSERT INTO insurance_chart_of_accounts (account_code, account_name, account_type, insurance_category, normal_balance, immutable_hash) VALUES
-- Assets
('1000', 'Cash and Bank', 'ASSET', 'cash', 'debit', 'init'),
('1100', 'Premium Receivable', 'ASSET', 'premium_income', 'debit', 'init'),
('1200', 'Claims Recoverable from Reinsurers', 'ASSET', 'reinsurance_recovery', 'debit', 'init'),
('1300', 'Salvage Inventory', 'ASSET', 'salvage_recovery', 'debit', 'init'),
('1400', 'Replacement Device Inventory', 'ASSET', 'replacement_device_cost', 'debit', 'init'),

-- Liabilities
('2000', 'Claim Reserves - Reported', 'LIABILITY', 'reserve_reported', 'credit', 'init'),
('2100', 'Claim Reserves - IBNR', 'LIABILITY', 'reserve_ibnr', 'credit', 'init'),
('2200', 'Unearned Premium Reserve', 'LIABILITY', 'reserve_uer', 'credit', 'init'),
('2300', 'Outstanding Claims Payable', 'LIABILITY', 'claim_expense', 'credit', 'init'),
('2400', 'Client Money Held', 'LIABILITY', 'client_money', 'credit', 'init'),

-- Income
('4000', 'Gross Premium Income', 'INCOME', 'premium_income', 'credit', 'init'),
('4100', 'Reinsurance Commission Income', 'INCOME', 'commission_expense', 'credit', 'init'),
('4200', 'Salvage Recovery Income', 'INCOME', 'salvage_recovery', 'credit', 'init'),

-- Expenses
('5000', 'Gross Claims Paid', 'EXPENSE', 'claim_expense', 'debit', 'init'),
('5100', 'Acquisition Costs', 'EXPENSE', 'acquisition_costs', 'debit', 'init'),
('5200', 'Maintenance Costs', 'EXPENSE', 'maintenance_costs', 'debit', 'init'),
('5300', 'Reinsurance Premium Ceded', 'EXPENSE', 'reinsurance_premium', 'debit', 'init'),
('5400', 'Replacement Device Cost', 'EXPENSE', 'replacement_device_cost', 'debit', 'init');

CREATE INDEX idx_coa_code ON insurance_chart_of_accounts(account_code);
CREATE INDEX idx_coa_type ON insurance_chart_of_accounts(account_type);
CREATE INDEX idx_coa_category ON insurance_chart_of_accounts(insurance_category);
CREATE INDEX idx_coa_current ON insurance_chart_of_accounts(is_active) WHERE is_active = TRUE;
