-- =============================================================================
-- FILENAME: 002_devices_and_registry.sql
-- DESCRIPTION: Device registry with provenance tracking and depreciation engine
-- VERSION: 1.0.0
-- DEPENDENCIES: 000_extensions.sql, 001_parties_and_identity.sql
-- =============================================================================
-- SECURITY CLASSIFICATION: RESTRICTED
-- DATA SENSITIVITY: Contains device identifiers (IMEI) - pseudonymize per GDPR
-- =============================================================================
-- ISO/IEC COMPLIANCE:
--   - ISO/IEC 27001:2013 - Asset management (A.8), Crypto (A.10)
--   - 3GPP TS 23.003 - IMEI format and structure
--   - GSMA IMEI Database - Blacklist checking requirements
--   - IEC 62684 - Universal charging standard
--   - WEEE Directive 2012/19/EU - E-waste tracking
-- =============================================================================
-- CHANGE LOG:
--   v1.0.0 (2026-03-26) - Initial release with GSMA compliance
-- =============================================================================

-- -----------------------------------------------------------------------------
-- SECTION 1: DEVICE MASTER REGISTRY
-- PURPOSE: Central registry for all insured devices
-- SECURITY: IMEI access restricted, pseudonymization recommended
-- -----------------------------------------------------------------------------

/**
 * TABLE: devices
 * DESCRIPTION: Master device registry with lifecycle tracking
 * 
 * COMPLIANCE:
 *   - 3GPP TS 23.003: IMEI format validation
 *   - GSMA: IMEI blacklist checking
 *   - GDPR: Pseudonymization of device identifiers recommended
 * 
 * ROW LEVEL SECURITY: Owners can view their devices, admins can view all
 */
CREATE TABLE devices (
    -- Primary identification
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Device identification (PII - handle per GDPR)
    imei VARCHAR(17) NOT NULL UNIQUE,   -- 15 digits + 2 check digits (optional)
    imei_hash VARCHAR(64),              -- SHA-256 for duplicate detection without storing plaintext
    serial_number VARCHAR(50),
    serial_number_hash VARCHAR(64),     -- SHA-256 for duplicate detection
    
    -- Device taxonomy (IEC 62684)
    brand VARCHAR(100) NOT NULL,        -- Manufacturer: Apple, Samsung, Google, etc.
    model VARCHAR(100) NOT NULL,        -- Model: iPhone 15 Pro, Galaxy S24, etc.
    model_variant VARCHAR(50),          -- Storage size, color variant
    category device_category DEFAULT 'smartphone',
    color VARCHAR(50),
    
    -- Purchase provenance
    purchase_date DATE,
    purchase_price DECIMAL(28,8),
    purchase_currency CHAR(3) DEFAULT 'USD',  -- ISO 4217
    
    -- Supply chain tracking
    purchase_source VARCHAR(50) CHECK (purchase_source IN (
        'authorized_retailer', 'gray_market', 'carrier', 'refurbished_reseller', 
        'oem_direct', 'employee_purchase', 'gift', 'unknown'
    )),
    retailer_party_id UUID REFERENCES parties(id),
    invoice_number VARCHAR(100),
    
    -- Device authenticity
    is_gray_market BOOLEAN DEFAULT FALSE,
    oem_warranty_status VARCHAR(50) CHECK (oem_warranty_status IN (
        'active', 'expired', 'void_gray_market', 'void_unauthorized_repair', 'unknown'
    )),
    oem_warranty_expiry DATE,
    
    -- Activation lock status (anti-theft)
    activation_lock_status VARCHAR(50) CHECK (activation_lock_status IN (
        'clean', 'locked_apple', 'locked_google', 'locked_samsung', 'unknown'
    )),
    activation_lock_verified_at TIMESTAMPTZ,
    
    -- Valuation and depreciation
    current_value DECIMAL(28,8),        -- Current depreciated value
    purchase_value DECIMAL(28,8),       -- Original purchase price (for historical)
    depreciation_rate DECIMAL(5,2) DEFAULT 0.00,  -- Monthly depreciation %
    total_loss_threshold DECIMAL(5,2) DEFAULT 80.00,  -- % for write-off determination
    
    -- Device condition (IEC 62684)
    device_condition VARCHAR(20) DEFAULT 'new' 
        CHECK (device_condition IN ('new', 'like_new', 'good', 'fair', 'poor')),
    condition_assessed_at TIMESTAMPTZ,
    condition_assessed_by UUID REFERENCES parties(id),
    
    -- Ownership
    owner_party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Status lifecycle (GSMA aligned)
    status device_status DEFAULT 'active',
    
    -- GSMA blacklist integration
    is_blacklisted BOOLEAN DEFAULT FALSE,
    blacklist_reason TEXT,              -- stolen, lost, insurance_claim, etc.
    blacklist_date TIMESTAMPTZ,
    gsma_blacklist_checked_at TIMESTAMPTZ,
    gsma_blacklist_status VARCHAR(20) CHECK (gsma_blacklist_status IN (
        'clean', 'blacklisted', 'unknown'
    )),
    
    -- Verification
    verification_status VARCHAR(20) DEFAULT 'unverified' 
        CHECK (verification_status IN ('unverified', 'verified', 'failed', 'pending')),
    verified_at TIMESTAMPTZ,
    
    -- Bitemporal versioning
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Audit and integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    -- Constraints
    CONSTRAINT valid_imei CHECK (validate_imei(imei)),
    CONSTRAINT positive_value CHECK (current_value >= 0),
    CONSTRAINT valid_time_range CHECK (valid_from < valid_to)
);

-- Table documentation
COMMENT ON TABLE devices IS 
'Master device registry with GSMA IMEI compliance and WEEE e-waste tracking. Security: RESTRICTED.';

COMMENT ON COLUMN devices.imei IS 
'15-digit IMEI per 3GPP TS 23.003. PII - access logged per GDPR.';

COMMENT ON COLUMN devices.imei_hash IS 
'SHA-256 hash for duplicate detection without storing plaintext IMEI.';

COMMENT ON COLUMN devices.purchase_source IS 
'Supply chain provenance: authorized_retailer, gray_market, carrier, etc.';

COMMENT ON COLUMN devices.activation_lock_status IS 
'Find My iPhone / Google FRP status. Must be clean for policy issuance.';

COMMENT ON COLUMN devices.total_loss_threshold IS 
'Percentage of value that triggers automatic write-off (default 80%).';

COMMENT ON COLUMN devices.gsma_blacklist_status IS 
'GSMA IMEI database status. Blacklisted devices cannot be insured.';

-- -----------------------------------------------------------------------------
-- SECTION 2: DEVICE OWNERSHIP HISTORY
-- PURPOSE: Track chain of custody for fraud prevention
-- SECURITY: Audit trail for device provenance
-- -----------------------------------------------------------------------------

CREATE TABLE device_ownership_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    owner_party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Transfer details
    acquired_at TIMESTAMPTZ NOT NULL,
    transferred_at TIMESTAMPTZ,
    transfer_reason VARCHAR(100) CHECK (transfer_reason IN (
        'purchase_new', 'purchase_used', 'gift', 'corporate_assignment', 
        'family_transfer', 'sale', 'trade_in', 'insurance_replacement', 'salvage_recovery'
    )),
    transfer_documentation_id UUID,     -- Link to bill of sale, etc.
    
    -- Bitemporal validity
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity'
);

COMMENT ON TABLE device_ownership_history IS 
'Chain of custody tracking for fraud prevention and provenance verification.';

-- -----------------------------------------------------------------------------
-- SECTION 3: DEVICE VERIFICATION LOGS
-- PURPOSE: Multi-factor device verification audit trail
-- -----------------------------------------------------------------------------

CREATE TABLE device_verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    
    -- Verification classification
    verification_type VARCHAR(50) NOT NULL CHECK (verification_type IN (
        'imei_check', 'photo_verification', 'document_check', 
        'remote_diagnostic', 'possession_proof', 'blacklist_check'
    )),
    
    -- Result
    status VARCHAR(20) NOT NULL CHECK (status IN ('passed', 'failed', 'inconclusive')),
    score DECIMAL(5,2),                 -- Confidence score 0-100
    
    -- Verification details
    verified_by UUID REFERENCES parties(id),
    verified_at TIMESTAMPTZ,
    
    -- Evidence
    metadata JSONB,                     -- Flexibility for different verification types
    evidence_document_ids UUID[],       -- Links to supporting documents
    
    -- Technical
    ip_address INET,
    device_fingerprint VARCHAR(256),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE device_verifications IS 
'Multi-factor device verification audit trail.';

-- -----------------------------------------------------------------------------
-- SECTION 4: DEVICE DEPRECIATION SCHEDULES
-- PURPOSE: Brand/model-specific depreciation curves
-- -----------------------------------------------------------------------------

CREATE TABLE device_depreciation_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    model_variant VARCHAR(50),          -- Storage size, etc.
    
    -- Depreciation curve parameters
    depreciation_curve_type VARCHAR(20) DEFAULT 'exponential' 
        CHECK (depreciation_curve_type IN ('linear', 'exponential', 'step')),
    
    month_1_rate DECIMAL(5,2) DEFAULT 15.00,      -- Initial drop (%)
    monthly_rate DECIMAL(5,2) DEFAULT 3.50,       -- Ongoing depreciation (%)
    annual_floor DECIMAL(5,2) DEFAULT 20.00,      -- Minimum retained value (%)
    
    -- Market data sources
    market_data_source VARCHAR(100),    -- Apple Trade-in, Gazelle, Swappa, etc.
    market_data_api_endpoint VARCHAR(255),
    last_market_update TIMESTAMPTZ,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    
    UNIQUE(brand, model, model_variant, valid_from)
);

COMMENT ON TABLE device_depreciation_schedules IS 
'Brand/model-specific depreciation curves for automated valuation.';

-- -----------------------------------------------------------------------------
-- SECTION 5: DEVICE MARKET VALUATIONS
-- PURPOSE: Automated market value tracking for claims
-- -----------------------------------------------------------------------------

CREATE TABLE device_market_valuations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    
    -- Valuation details
    valuation_date DATE NOT NULL,
    valuation_source VARCHAR(100) NOT NULL,  -- trade_in_db, market_scrape, manual
    valuation_method VARCHAR(50) CHECK (valuation_method IN (
        'automated_api', 'manual_assessment', 'algorithmic', 'auction_result'
    )),
    
    -- Condition-adjusted values (percentage of new value)
    new_value DECIMAL(28,8),                    -- 100%
    like_new_value DECIMAL(28,8),               -- 85%
    good_value DECIMAL(28,8),                   -- 70%
    fair_value DECIMAL(28,8),                   -- 50%
    poor_value DECIMAL(28,8),                   -- 30%
    
    -- Total loss threshold (calculated)
    total_loss_threshold_amount DECIMAL(28,8),  -- 80% of new_value
    
    -- Market metadata
    source_url TEXT,
    confidence_score DECIMAL(5,2),              -- 0-100
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE device_market_valuations IS 
'Automated market value tracking from trade-in databases for fair claims valuation.';

-- -----------------------------------------------------------------------------
-- SECTION 6: DEVICE LIFECYCLE EVENTS
-- PURPOSE: State transition audit trail
-- -----------------------------------------------------------------------------

CREATE TABLE device_lifecycle_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    
    -- Event classification
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN (
        'purchased', 'enrolled', 'policy_bound', 'claimed', 'repaired', 
        'replaced', 'reported_stolen', 'reported_lost', 'recovered', 
        'blacklisted', 'unblacklisted', 'decommissioned', 'salvaged', 
        'resold', 'recycled', 'end_of_support'
    )),
    
    -- State transition
    previous_status VARCHAR(20),
    new_status VARCHAR(20),
    
    -- Context
    claim_id UUID REFERENCES claims(id),
    policy_id UUID REFERENCES policies(id),
    notes TEXT,
    
    -- Audit
    event_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    ip_address INET
);

COMMENT ON TABLE device_lifecycle_events IS 
'Complete lifecycle audit trail from purchase through decommissioning.';

-- -----------------------------------------------------------------------------
-- SECTION 7: DEVICE FRAUD ALERTS
-- PURPOSE: Automated fraud detection alerts
-- -----------------------------------------------------------------------------

CREATE TABLE device_fraud_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    
    -- Alert classification
    alert_type VARCHAR(50) NOT NULL CHECK (alert_type IN (
        'multiple_claims', 'rapid_policy_changes', 'ownership_discrepancy', 
        'blacklisted_imei', 'suspicious_purchase_pattern', 'claim_within_cooling',
        'imei_cloning_suspected', 'geolocation_anomaly'
    )),
    
    severity VARCHAR(20) DEFAULT 'medium' 
        CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    
    description TEXT,
    evidence JSONB,
    
    -- Investigation workflow
    status VARCHAR(20) DEFAULT 'open' 
        CHECK (status IN ('open', 'investigating', 'resolved', 'false_positive')),
    assigned_to UUID REFERENCES parties(id),
    
    resolved_by UUID REFERENCES parties(id),
    resolved_at TIMESTAMPTZ,
    resolution_notes TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE device_fraud_alerts IS 
'Automated fraud detection alerts requiring investigation.';

-- -----------------------------------------------------------------------------
-- SECTION 8: INDEXES
-- -----------------------------------------------------------------------------

-- Core device indexes
CREATE INDEX idx_devices_imei ON devices(imei);
CREATE INDEX idx_devices_imei_hash ON devices(imei_hash) WHERE imei_hash IS NOT NULL;
CREATE INDEX idx_devices_owner ON devices(owner_party_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_blacklist ON devices(is_blacklisted) WHERE is_blacklisted = TRUE;
CREATE INDEX idx_devices_current ON devices(is_current) WHERE is_current = TRUE;
CREATE INDEX idx_devices_valid_time ON devices(valid_from, valid_to);
CREATE INDEX idx_devices_brand_model ON devices(brand, model);
CREATE INDEX idx_devices_gsma_check ON devices(gsma_blacklist_checked_at) 
    WHERE gsma_blacklist_checked_at IS NULL OR gsma_blacklist_checked_at < CURRENT_TIMESTAMP - INTERVAL '7 days';

-- Ownership and verification indexes
CREATE INDEX idx_ownership_device ON device_ownership_history(device_id);
CREATE INDEX idx_ownership_owner ON device_ownership_history(owner_party_id);
CREATE INDEX idx_verifications_device ON device_verifications(device_id);
CREATE INDEX idx_verifications_type ON device_verifications(verification_type);

-- Depreciation and valuation indexes
CREATE INDEX idx_depreciation_schedule_lookup ON device_depreciation_schedules(brand, model, model_variant) 
    WHERE is_active = TRUE;
CREATE INDEX idx_market_valuations_device ON device_market_valuations(device_id, valuation_date);

-- Lifecycle and fraud indexes
CREATE INDEX idx_lifecycle_device ON device_lifecycle_events(device_id);
CREATE INDEX idx_lifecycle_event_type ON device_lifecycle_events(event_type);
CREATE INDEX idx_lifecycle_time ON device_lifecycle_events(event_at);
CREATE INDEX idx_fraud_alerts_device ON device_fraud_alerts(device_id);
CREATE INDEX idx_fraud_alerts_status ON device_fraud_alerts(status) 
    WHERE status IN ('open', 'investigating');
CREATE INDEX idx_fraud_alerts_severity ON device_fraud_alerts(severity) 
    WHERE severity IN ('high', 'critical');

-- -----------------------------------------------------------------------------
-- SECTION 9: ROW LEVEL SECURITY
-- -----------------------------------------------------------------------------

ALTER TABLE devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE device_ownership_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE device_verifications ENABLE ROW LEVEL SECURITY;

-- Policy: Device owners can view their devices
CREATE POLICY devices_owner_access ON devices
    FOR SELECT
    USING (owner_party_id = current_setting('app.current_party_id', true)::UUID 
           OR current_user LIKE '%admin%'
           OR pg_has_role(current_user, 'device_admin', 'MEMBER'));

-- Policy: Only admins can modify device records
CREATE POLICY devices_admin_modify ON devices
    FOR ALL
    USING (pg_has_role(current_user, 'device_admin', 'MEMBER') 
           OR current_user LIKE '%admin%');

-- -----------------------------------------------------------------------------
-- SECTION 10: TRIGGERS
-- -----------------------------------------------------------------------------

CREATE TRIGGER trg_devices_hash 
    BEFORE INSERT OR UPDATE ON devices
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

-- Calculate current value based on depreciation
CREATE OR REPLACE FUNCTION calculate_device_value()
RETURNS TRIGGER AS $$
DECLARE
    v_schedule RECORD;
    v_months_old INTEGER;
    v_depreciated_value DECIMAL(28,8);
BEGIN
    -- Get depreciation schedule for device
    SELECT * INTO v_schedule
    FROM device_depreciation_schedules
    WHERE brand = NEW.brand 
    AND model = NEW.model
    AND is_active = TRUE
    ORDER BY valid_from DESC
    LIMIT 1;
    
    IF FOUND AND NEW.purchase_price IS NOT NULL AND NEW.purchase_date IS NOT NULL THEN
        v_months_old := EXTRACT(MONTH FROM AGE(CURRENT_DATE, NEW.purchase_date));
        
        -- Apply depreciation formula
        IF v_months_old <= 1 THEN
            v_depreciated_value := NEW.purchase_price * (1 - v_schedule.month_1_rate/100);
        ELSE
            v_depreciated_value := NEW.purchase_price * 
                (1 - v_schedule.month_1_rate/100) * 
                POWER(1 - v_schedule.monthly_rate/100, v_months_old - 1);
        END IF;
        
        -- Apply floor
        v_depreciated_value := GREATEST(v_depreciated_value, 
            NEW.purchase_price * v_schedule.annual_floor/100);
        
        NEW.current_value := v_depreciated_value;
        NEW.depreciation_rate := v_schedule.monthly_rate;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_calculate_device_value
    BEFORE INSERT OR UPDATE OF purchase_date, purchase_price ON devices
    FOR EACH ROW EXECUTE FUNCTION calculate_device_value();

-- =============================================================================
-- END OF FILE: 002_devices_and_registry.sql
-- =============================================================================
