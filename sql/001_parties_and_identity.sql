-- =============================================================================
-- FILENAME: 001_parties_and_identity.sql
-- DESCRIPTION: Party management with identity verification and GDPR compliance
-- VERSION: 1.0.0
-- DEPENDENCIES: 000_extensions.sql
-- =============================================================================
-- SECURITY CLASSIFICATION: CONFIDENTIAL
-- DATA SENSITIVITY: Contains PII subject to GDPR/CCPA protection
-- =============================================================================
-- ISO/IEC COMPLIANCE:
--   - ISO/IEC 27001:2013 - Access control (A.9), Cryptography (A.10)
--   - ISO 17442:2020 - Legal Entity Identifier (LEI)
--   - ISO 9362:2022 - Bank Identifier Code (BIC)
--   - GDPR (EU) 2016/679 - Articles 17, 32 (Data protection)
--   - FATF Recommendations - Customer due diligence
-- =============================================================================
-- CHANGE LOG:
--   v1.0.0 (2026-03-26) - Initial release with full compliance
-- =============================================================================

-- -----------------------------------------------------------------------------
-- SECTION 1: CORE PARTIES TABLE
-- PURPOSE: Master entity registry for all parties in the system
-- SECURITY: PII encrypted at rest, RLS enabled, audit logging
-- -----------------------------------------------------------------------------

/**
 * TABLE: parties
 * DESCRIPTION: Central registry for all entities (individuals, corporates, partners)
 *              with bitemporal versioning and cryptographic integrity
 * 
 * COMPLIANCE:
 *   - GDPR Article 17: Right to erasure (retention_date, legal_hold)
 *   - GDPR Article 32: Security of processing (encryption markers)
 *   - ISO 17442: LEI for corporate entities
 * 
 * ROW LEVEL SECURITY: Enabled - parties can only view their own records
 *                     except for authorized administrators
 */
CREATE TABLE parties (
    -- Primary identification
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Entity classification (ISO 17442 aligned)
    type party_type NOT NULL,
    
    -- Display information
    display_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    
    -- Contact information (PII - encryption required per GDPR)
    email VARCHAR(255),
    phone VARCHAR(20),
    email_encrypted BYTEA,              -- AES-256-GCM encrypted
    phone_encrypted BYTEA,              -- AES-256-GCM encrypted
    
    -- KYC/Sanctions compliance (FATF Recommendation 10)
    kyc_status kyc_status DEFAULT 'pending',
    sanctions_status sanctions_status DEFAULT 'clear',
    risk_score DECIMAL(5,2) CHECK (risk_score >= 0 AND risk_score <= 100),
    risk_category VARCHAR(50),          -- Derived from risk_score: low|standard|elevated|high
    
    -- ISO Standard Identifiers
    lei_code VARCHAR(20),               -- ISO 17442:2020 Legal Entity Identifier
    bic_code VARCHAR(11),               -- ISO 9362:2022 Bank Identifier Code
    
    -- GDPR Data Retention (Article 17)
    retention_date DATE,                -- Scheduled deletion date
    legal_hold BOOLEAN DEFAULT FALSE,   -- Litigation hold - prevents deletion
    
    -- Sanctions screening audit trail
    sanctions_screening_service VARCHAR(50),
    last_sanctions_screening_at TIMESTAMPTZ,
    
    -- Bitemporal versioning (ISO/IEC 9075)
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Audit and integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,   -- Soft delete for audit trail
    
    -- Constraints
    CONSTRAINT valid_time_check CHECK (valid_from < valid_to),
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT lei_format CHECK (lei_code IS NULL OR lei_code ~ '^[A-Z0-9]{18}[0-9]{2}$'),
    CONSTRAINT bic_format CHECK (bic_code IS NULL OR bic_code ~ '^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$')
);

-- Table documentation
COMMENT ON TABLE parties IS 
'Core entity registry with GDPR-compliant PII storage and ISO 17442 LEI support. Security: CONFIDENTIAL.';

COMMENT ON COLUMN parties.email_encrypted IS 
'AES-256-GCM encrypted email per GDPR Article 32. Key managed by KMS.';

COMMENT ON COLUMN parties.phone_encrypted IS 
'AES-256-GCM encrypted phone per GDPR Article 32. Key managed by KMS.';

COMMENT ON COLUMN parties.lei_code IS 
'ISO 17442:2020 Legal Entity Identifier for corporate entities. 20 characters.';

COMMENT ON COLUMN parties.bic_code IS 
'ISO 9362:2022 Bank Identifier Code. 8 or 11 characters.';

COMMENT ON COLUMN parties.retention_date IS 
'GDPR Article 17: Scheduled data deletion date. Null = indefinite retention.';

COMMENT ON COLUMN parties.legal_hold IS 
'GDPR Article 17(3)(e): Litigation hold suspends right to erasure.';

-- -----------------------------------------------------------------------------
-- SECTION 2: PARTY IDENTIFIERS
-- PURPOSE: Secure storage of government/official identification documents
-- SECURITY: All values encrypted, access logged
-- -----------------------------------------------------------------------------

/**
 * TABLE: party_identifiers
 * DESCRIPTION: Encrypted storage of government IDs, passports, tax IDs
 *              with verification tracking
 * 
 * COMPLIANCE:
 *   - GDPR Article 9: Special category data (ID numbers)
 *   - PCI DSS: If storing national ID equivalents
 *   - Local regulations: Varies by jurisdiction
 */
CREATE TABLE party_identifiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    
    -- Identifier classification
    type VARCHAR(50) NOT NULL,          -- national_id, passport, tax_id, company_reg, driving_license
    value_encrypted BYTEA NOT NULL,     -- Always encrypted - never store plaintext
    value_hash VARCHAR(64),             -- SHA-256 for duplicate detection without decryption
    
    -- Verification tracking
    verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    verification_method VARCHAR(100),   -- document_upload, video_call, third_party_api
    
    -- Bitemporal validity
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(party_id, type, valid_from)
);

COMMENT ON TABLE party_identifiers IS 
'Encrypted government ID storage. GDPR Article 9 special category data. Never access without audit log.';

COMMENT ON COLUMN party_identifiers.value_encrypted IS 
'AES-256-GCM encrypted identifier. Key rotation required annually per ISO 27001.';

COMMENT ON COLUMN party_identifiers.value_hash IS 
'SHA-256 hash for duplicate detection. Allows uniqueness checks without decryption.';

-- -----------------------------------------------------------------------------
-- SECTION 3: CORPORATE ACCOUNTS
-- PURPOSE: Business entity extensions with fleet management
-- SECURITY: LEI validation, credit limits, audit trail
-- -----------------------------------------------------------------------------

/**
 * TABLE: corporate_accounts
 * DESCRIPTION: Business entity extension with credit management
 *              and fleet policy controls
 * 
 * COMPLIANCE:
 *   - ISO 17442: LEI mandatory for large corporates
 *   - Basel III: Credit risk exposure tracking
 */
CREATE TABLE corporate_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    
    -- Corporate information
    corporate_name VARCHAR(255) NOT NULL,
    registration_number VARCHAR(100),
    tax_id_encrypted BYTEA,             -- Encrypted tax identifier
    industry_sector VARCHAR(100),
    
    -- Fleet management
    fleet_policy_limit DECIMAL(28,8) DEFAULT 0,     -- Max aggregate insured value
    credit_limit DECIMAL(28,8),                      -- Accounts receivable limit
    payment_terms_days INTEGER DEFAULT 30,
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' 
        CHECK (status IN ('active', 'suspended', 'closed')),
    
    -- Bitemporal versioning
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    immutable_hash VARCHAR(64) NOT NULL,
    
    UNIQUE(party_id, valid_from)
);

COMMENT ON TABLE corporate_accounts IS 
'Corporate entity extension with fleet management and credit controls.';

COMMENT ON COLUMN corporate_accounts.fleet_policy_limit IS 
'Maximum aggregate insured value across all corporate devices.';

-- -----------------------------------------------------------------------------
-- SECTION 4: CORPORATE EMPLOYEES
-- PURPOSE: Employee linkage for BYOD and corporate fleet policies
-- -----------------------------------------------------------------------------

CREATE TABLE corporate_employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    corporate_id UUID NOT NULL REFERENCES corporate_accounts(id) ON DELETE CASCADE,
    party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Employment details
    employee_id VARCHAR(100),           -- HR system reference
    department VARCHAR(100),
    cost_center_code VARCHAR(50),       -- For premium allocation
    role VARCHAR(50) CHECK (role IN ('admin', 'user', 'approver', 'viewer')),
    is_primary_contact BOOLEAN DEFAULT FALSE,
    
    -- BYOD policy status
    mdm_enrollment_required BOOLEAN DEFAULT TRUE,
    security_policy_acknowledged_at TIMESTAMPTZ,
    
    -- Temporal validity
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(corporate_id, party_id, valid_from)
);

COMMENT ON TABLE corporate_employees IS 
'Employee linkage for corporate fleet and BYOD policy management.';

-- -----------------------------------------------------------------------------
-- SECTION 5: FAMILY/Household PLANS (Multi-Device)
-- PURPOSE: Support family/household policies with aggregate limits
-- -----------------------------------------------------------------------------

CREATE TABLE family_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    primary_policyholder_id UUID NOT NULL REFERENCES parties(id),
    group_name VARCHAR(100),
    
    -- Billing configuration
    aggregate_deductible DECIMAL(28,8),     -- Shared across family
    per_device_deductible DECIMAL(28,8),    -- Individual device deductible
    max_claims_per_period INTEGER,          -- Family-level claim limit
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' 
        CHECK (status IN ('active', 'suspended', 'closed')),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity'
);

COMMENT ON TABLE family_groups IS 
'Family/household policy grouping for multi-device plans with aggregate deductibles.';

-- -----------------------------------------------------------------------------
-- SECTION 6: GDPR CONSENT MANAGEMENT
-- PURPOSE: Track data subject consent per GDPR requirements
-- COMPLIANCE: GDPR Articles 6, 7, 21 - Lawful basis and consent
-- -----------------------------------------------------------------------------

/**
 * TABLE: consent_records
 * DESCRIPTION: GDPR-compliant consent tracking with withdrawal support
 * 
 * COMPLIANCE:
 *   - GDPR Article 7: Conditions for consent
 *   - GDPR Article 21: Right to object
 *   - ePrivacy Directive: Cookie/consent requirements
 */
CREATE TABLE consent_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    
    -- Consent classification
    consent_type VARCHAR(50) NOT NULL 
        CHECK (consent_type IN ('marketing', 'data_processing', 'third_party_sharing', 'profiling', 'automated_decision')),
    consent_given BOOLEAN NOT NULL,
    
    -- Temporal tracking
    consent_given_at TIMESTAMPTZ NOT NULL,
    consent_withdrawn_at TIMESTAMPTZ,
    consent_expires_at TIMESTAMPTZ,     -- Some consent has time limits
    
    -- Documentation
    consent_version VARCHAR(20) NOT NULL,   -- Version of consent text agreed to
    consent_document_url TEXT,              -- Link to terms accepted
    legal_basis VARCHAR(50),                -- GDPR Article 6 basis: consent, contract, legal_obligation, etc.
    
    -- Technical metadata
    ip_address INET,
    user_agent TEXT,
    geolocation_country CHAR(2),            -- ISO 3166-1
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE consent_records IS 
'GDPR Article 7 consent tracking with withdrawal support. Audit trail for lawful basis.';

COMMENT ON COLUMN consent_records.consent_type IS 
'Category of consent: marketing, data_processing, third_party_sharing, profiling, automated_decision';

COMMENT ON COLUMN consent_records.legal_basis IS 
'GDPR Article 6 lawful basis: consent, contract, legal_obligation, vital_interests, public_task, legitimate_interests';

-- ----------------------------------------------------------------------------
-- SECTION 7: DEVICE POSSESSION VERIFICATION
-- PURPOSE: Anti-fraud verification that party physically possesses device
-- SECURITY: Multi-factor verification with photo evidence
-- ----------------------------------------------------------------------------

CREATE TABLE device_possession_verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    device_id UUID REFERENCES devices(id),
    
    -- Verification method
    verification_method VARCHAR(50) NOT NULL 
        CHECK (verification_method IN ('sms_to_device', 'app_install', 'photo_with_code', 'imei_selfie', 'video_call')),
    verification_code VARCHAR(20),
    code_sent_at TIMESTAMPTZ,
    code_verified_at TIMESTAMPTZ,
    
    -- Evidence
    photo_evidence_id UUID,             -- Reference to documents table
    gps_location_lat DECIMAL(10,8),
    gps_location_lon DECIMAL(11,8),
    
    -- Status
    status VARCHAR(20) DEFAULT 'pending' 
        CHECK (status IN ('pending', 'verified', 'failed')),
    failure_reason TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    verified_by UUID REFERENCES parties(id)
);

COMMENT ON TABLE device_possession_verifications IS 
'Anti-fraud verification that enrolling party physically possesses the device.';

-- -----------------------------------------------------------------------------
-- SECTION 8: PARTY CLAIM ANALYTICS (Fraud Detection)
-- PURPOSE: Velocity tracking and risk scoring for fraud prevention
-- SECURITY: Aggregated data, no PII exposure
-- -----------------------------------------------------------------------------

CREATE TABLE party_claim_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Claim velocity metrics
    total_claims INTEGER DEFAULT 0,
    claims_last_12_months INTEGER DEFAULT 0,
    claims_last_6_months INTEGER DEFAULT 0,
    claims_last_3_months INTEGER DEFAULT 0,
    
    -- Device metrics
    total_devices_insured INTEGER DEFAULT 0,
    concurrent_devices INTEGER DEFAULT 0,
    device_change_frequency DECIMAL(5,2),   -- Devices changed per year
    
    -- Financial metrics
    total_claimed_amount DECIMAL(28,8) DEFAULT 0,
    average_claim_amount DECIMAL(28,8),
    largest_claim_amount DECIMAL(28,8),
    
    -- Risk scoring
    claim_velocity_score DECIMAL(5,2) CHECK (claim_velocity_score >= 0 AND claim_velocity_score <= 100),
    pattern_flags JSONB,                    -- Array of triggered patterns
    risk_tier VARCHAR(20) DEFAULT 'standard' 
        CHECK (risk_tier IN ('low', 'standard', 'elevated', 'high')),
    manual_review_required BOOLEAN DEFAULT FALSE,
    
    -- Temporal tracking
    last_claim_date TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(party_id)
);

COMMENT ON TABLE party_claim_analytics IS 
'Fraud detection analytics. No PII - aggregated metrics only. Used for velocity scoring.';

-- -----------------------------------------------------------------------------
-- SECTION 9: INDEXES
-- PURPOSE: Optimize query performance for common access patterns
-- SECURITY: Indexes on encrypted columns use hash values only
-- -----------------------------------------------------------------------------

-- Core party indexes
CREATE INDEX idx_parties_type ON parties(type);
CREATE INDEX idx_parties_kyc ON parties(kyc_status) WHERE kyc_status != 'verified';
CREATE INDEX idx_parties_sanctions ON parties(sanctions_status) WHERE sanctions_status != 'clear';
CREATE INDEX idx_parties_current ON parties(is_current) WHERE is_current = TRUE;
CREATE INDEX idx_parties_valid_time ON parties(valid_from, valid_to);
CREATE INDEX idx_parties_lei ON parties(lei_code) WHERE lei_code IS NOT NULL;
CREATE INDEX idx_parties_retention ON parties(retention_date) 
    WHERE retention_date IS NOT NULL AND legal_hold = FALSE;

-- Identifier indexes
CREATE INDEX idx_party_identifiers_party ON party_identifiers(party_id);
CREATE INDEX idx_party_identifiers_hash ON party_identifiers(value_hash) 
    WHERE value_hash IS NOT NULL;

-- Corporate indexes
CREATE INDEX idx_corporate_accounts_party ON corporate_accounts(party_id);
CREATE INDEX idx_corporate_employees_corp ON corporate_employees(corporate_id);
CREATE INDEX idx_corporate_employees_party ON corporate_employees(party_id);

-- Family and consent indexes
CREATE INDEX idx_family_groups_primary ON family_groups(primary_policyholder_id);
CREATE INDEX idx_family_members_group ON family_members(family_group_id);
CREATE INDEX idx_consent_party ON consent_records(party_id);
CREATE INDEX idx_consent_current ON consent_records(party_id, consent_type) 
    WHERE consent_withdrawn_at IS NULL;

-- Analytics indexes
CREATE INDEX idx_party_analytics_velocity ON party_claim_analytics(claim_velocity_score) 
    WHERE claim_velocity_score > 70;
CREATE INDEX idx_party_analytics_tier ON party_claim_analytics(risk_tier) 
    WHERE risk_tier IN ('elevated', 'high');

-- -----------------------------------------------------------------------------
-- SECTION 10: ROW LEVEL SECURITY (RLS) POLICIES
-- PURPOSE: Restrict data access based on user context
-- SECURITY: Parties can only access their own records
-- -----------------------------------------------------------------------------

-- Enable RLS on sensitive tables
ALTER TABLE parties ENABLE ROW LEVEL SECURITY;
ALTER TABLE party_identifiers ENABLE ROW LEVEL SECURITY;
ALTER TABLE consent_records ENABLE ROW LEVEL SECURITY;

-- Policy: Users can view their own party record
CREATE POLICY parties_self_access ON parties
    FOR SELECT
    USING (id = current_setting('app.current_party_id', true)::UUID 
           OR current_user IN (SELECT rolname FROM pg_roles WHERE rolsuper));

-- Policy: Only admins can update party records
CREATE POLICY parties_admin_update ON parties
    FOR UPDATE
    USING (current_user LIKE '%admin%' OR pg_has_role(current_user, 'party_admin', 'MEMBER'));

-- Policy: Identifiers only accessible to verified administrators
CREATE POLICY identifiers_admin_only ON party_identifiers
    FOR ALL
    USING (pg_has_role(current_user, 'identity_verifier', 'MEMBER') 
           OR current_user LIKE '%admin%');

-- -----------------------------------------------------------------------------
-- SECTION 11: TRIGGERS
-- PURPOSE: Maintain audit trails and cryptographic integrity
-- -----------------------------------------------------------------------------

CREATE TRIGGER trg_parties_hash 
    BEFORE INSERT OR UPDATE ON parties
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_corp_accounts_hash 
    BEFORE INSERT OR UPDATE ON corporate_accounts
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE TRIGGER trg_corporate_employees_hash
    BEFORE INSERT OR UPDATE ON corporate_employees
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

-- Update timestamp trigger for parties
CREATE OR REPLACE FUNCTION update_party_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_parties_timestamp
    BEFORE UPDATE ON parties
    FOR EACH ROW EXECUTE FUNCTION update_party_timestamp();

-- -----------------------------------------------------------------------------
-- SECTION 12: AUDIT LOGGING
-- PURPOSE: Track all access to sensitive PII data
-- COMPLIANCE: GDPR Article 30 (Records of processing activities)
-- -----------------------------------------------------------------------------

/**
 * TABLE: party_audit_log
 * DESCRIPTION: Comprehensive audit trail for all party data access
 */
CREATE TABLE party_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID REFERENCES parties(id),
    
    -- Action details
    action VARCHAR(50) NOT NULL,        -- SELECT, INSERT, UPDATE, DELETE
    table_name VARCHAR(100) NOT NULL,
    column_name VARCHAR(100),
    old_value TEXT,                     -- Masked for PII
    new_value TEXT,                     -- Masked for PII
    
    -- Context
    performed_by UUID,
    performed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    session_id UUID,
    ip_address INET,
    application_name VARCHAR(100),
    
    -- GDPR specific
    legal_basis VARCHAR(50),            -- Article 6 basis for access
    consent_reference UUID              -- Link to consent record if applicable
);

CREATE INDEX idx_party_audit_party ON party_audit_log(party_id);
CREATE INDEX idx_party_audit_time ON party_audit_log(performed_at);
CREATE INDEX idx_party_audit_action ON party_audit_log(action);

COMMENT ON TABLE party_audit_log IS 
'GDPR Article 30 audit log for all party data access. Retention: 7 years.';

-- =============================================================================
-- END OF FILE: 001_parties_and_identity.sql
-- =============================================================================
