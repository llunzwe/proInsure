-- 001_parties_and_identity.sql

-- Core parties table (individuals, corporates, repair shops, etc.)
CREATE TABLE parties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type party_type NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    email_encrypted BYTEA, -- Encrypted storage
    phone_encrypted BYTEA,
    kyc_status kyc_status DEFAULT 'pending',
    sanctions_status sanctions_status DEFAULT 'clear',
    risk_score DECIMAL(5,2) CHECK (risk_score >= 0 AND risk_score <= 100),
    risk_category VARCHAR(50), -- Derived from score
    
    -- Bitemporal fields
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Audit and integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    -- Constraints
    CONSTRAINT valid_time_check CHECK (valid_from < valid_to),
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- Party identifiers (national ID, passport, tax ID)
CREATE TABLE party_identifiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- national_id, passport, tax_id, company_reg, driving_license
    value_encrypted BYTEA NOT NULL, -- Always encrypted
    value_hash VARCHAR(64), -- For duplicate detection without decryption
    verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    verification_method VARCHAR(100),
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(party_id, type, valid_from)
);

-- Corporate accounts extension
CREATE TABLE corporate_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    corporate_name VARCHAR(255) NOT NULL,
    registration_number VARCHAR(100),
    tax_id_encrypted BYTEA,
    industry_sector VARCHAR(100),
    fleet_policy_limit DECIMAL(28,8) DEFAULT 0, -- Max aggregate insured value
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'closed')),
    credit_limit DECIMAL(28,8),
    payment_terms_days INTEGER DEFAULT 30,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    immutable_hash VARCHAR(64) NOT NULL,
    UNIQUE(party_id, valid_from)
);

-- Corporate employees linkage
CREATE TABLE corporate_employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    corporate_id UUID NOT NULL REFERENCES corporate_accounts(id) ON DELETE CASCADE,
    party_id UUID NOT NULL REFERENCES parties(id), -- The employee
    employee_id VARCHAR(100), -- Internal ID
    department VARCHAR(100),
    role VARCHAR(50) CHECK (role IN ('admin', 'user', 'approver', 'viewer')),
    is_primary_contact BOOLEAN DEFAULT FALSE,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(corporate_id, party_id, valid_from)
);

-- Indexes for performance
CREATE INDEX idx_parties_type ON parties(type);
CREATE INDEX idx_parties_kyc ON parties(kyc_status) WHERE kyc_status != 'verified';
CREATE INDEX idx_parties_sanctions ON parties(sanctions_status) WHERE sanctions_status != 'clear';
CREATE INDEX idx_parties_current ON parties(is_current) WHERE is_current = TRUE;
CREATE INDEX idx_parties_valid_time ON parties(valid_from, valid_to);
CREATE INDEX idx_party_identifiers_party ON party_identifiers(party_id);
CREATE INDEX idx_corporate_accounts_party ON corporate_accounts(party_id);
CREATE INDEX idx_corporate_employees_corp ON corporate_employees(corporate_id);
CREATE INDEX idx_corporate_employees_party ON corporate_employees(party_id);

-- Triggers for immutability
CREATE TRIGGER trg_parties_hash BEFORE INSERT OR UPDATE ON parties
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();
CREATE TRIGGER trg_corp_accounts_hash BEFORE INSERT OR UPDATE ON corporate_accounts
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();
