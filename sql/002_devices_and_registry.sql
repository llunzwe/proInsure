-- 002_devices_and_registry.sql

-- Central device registry
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    imei VARCHAR(17) NOT NULL UNIQUE,
    serial_number VARCHAR(50),
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    category device_category DEFAULT 'smartphone',
    color VARCHAR(50),
    
    -- Purchase info
    purchase_date DATE,
    purchase_price DECIMAL(28,8),
    purchase_currency CHAR(3) DEFAULT 'USD',
    current_value DECIMAL(28,8), -- Depreciated value
    depreciation_rate DECIMAL(5,2) DEFAULT 0.00, -- Monthly depreciation %
    
    -- Ownership
    owner_party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Status
    status device_status DEFAULT 'active',
    is_blacklisted BOOLEAN DEFAULT FALSE,
    blacklist_reason TEXT,
    blacklist_date TIMESTAMPTZ,
    verification_status VARCHAR(20) DEFAULT 'unverified' CHECK (verification_status IN ('unverified', 'verified', 'failed', 'pending')),
    verified_at TIMESTAMPTZ,
    
    -- Bitemporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    -- Constraints
    CONSTRAINT valid_imei CHECK (validate_imei(imei)),
    CONSTRAINT positive_value CHECK (current_value >= 0),
    CONSTRAINT valid_time_range CHECK (valid_from < valid_to)
);

-- Device ownership history (explicit tracking)
CREATE TABLE device_ownership_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    owner_party_id UUID NOT NULL REFERENCES parties(id),
    acquired_at TIMESTAMPTZ NOT NULL,
    transferred_at TIMESTAMPTZ,
    transfer_reason VARCHAR(100), -- sale, gift, corporate_assignment, etc.
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity'
);

-- Device verification logs
CREATE TABLE device_verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    verification_type VARCHAR(50) NOT NULL, -- imei_check, photo_verification, document_check
    status VARCHAR(20) NOT NULL,
    verified_by UUID REFERENCES parties(id), -- Assessor or system
    metadata JSONB, -- Flexibility for different verification types
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_devices_imei ON devices(imei);
CREATE INDEX idx_devices_owner ON devices(owner_party_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_blacklist ON devices(is_blacklisted) WHERE is_blacklisted = TRUE;
CREATE INDEX idx_devices_current ON devices(is_current) WHERE is_current = TRUE;
CREATE INDEX idx_devices_valid_time ON devices(valid_from, valid_to);
CREATE INDEX idx_ownership_device ON device_ownership_history(device_id);
CREATE INDEX idx_ownership_owner ON device_ownership_history(owner_party_id);

-- Trigger for IMEI validation and hash
CREATE TRIGGER trg_devices_hash BEFORE INSERT OR UPDATE ON devices
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();
