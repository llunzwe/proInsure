-- 003_device_specifications_and_iot.sql

-- Technical specifications
CREATE TABLE device_specifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    storage_gb INTEGER,
    ram_gb INTEGER,
    os VARCHAR(50), -- iOS, Android, etc.
    os_version VARCHAR(50),
    screen_size DECIMAL(4,2), -- Inches
    screen_resolution VARCHAR(20), -- e.g., "1170x2532"
    color VARCHAR(50),
    condition_grade VARCHAR(20) CHECK (condition_grade IN ('new', 'like_new', 'good', 'fair', 'poor', 'damaged')),
    warranty_expiry DATE,
    carrier_locked BOOLEAN,
    carrier VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(device_id)
);

-- IoT telemetry data (TimescaleDB hypertable)
CREATE TABLE device_iot_telemetry (
    id UUID DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Device health
    battery_level INTEGER CHECK (battery_level BETWEEN 0 AND 100),
    battery_health INTEGER CHECK (battery_health BETWEEN 0 AND 100), -- Percentage of original capacity
    storage_used_gb INTEGER,
    storage_available_gb INTEGER,
    
    -- Location (encrypted or hashed for privacy)
    location_lat DECIMAL(10,8),
    location_lon DECIMAL(11,8),
    location_accuracy DECIMAL(8,2), -- Meters
    location_city VARCHAR(100), -- Derived/approximate
    
    -- Security status
    security_status VARCHAR(50) CHECK (security_status IN ('locked', 'unlocked', 'compromised', 'unknown')),
    find_my_device_enabled BOOLEAN,
    last_backup_at TIMESTAMPTZ,
    
    -- Health scoring
    health_score INTEGER CHECK (health_score BETWEEN 0 AND 100), -- Predictive failure risk
    risk_factors JSONB, -- Array of detected risks
    
    -- Raw telemetry payload
    raw_payload JSONB,
    
    PRIMARY KEY (id, recorded_at)
);

-- Convert to TimescaleDB hypertable for time-series optimization
SELECT create_hypertable('device_iot_telemetry', 'recorded_at', 
                         chunk_time_interval => INTERVAL '7 days',
                         if_not_exists => TRUE);

-- Device analytics (materialized view or projection table)
CREATE TABLE device_analytics (
    device_id UUID PRIMARY KEY REFERENCES devices(id),
    claim_count INTEGER DEFAULT 0,
    total_claimed_amount DECIMAL(28,8) DEFAULT 0,
    average_repair_cost DECIMAL(28,8),
    last_claim_date TIMESTAMPTZ,
    fraud_risk_score DECIMAL(5,2),
    reliability_score INTEGER CHECK (reliability_score BETWEEN 0 AND 100),
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_iot_device_time ON device_iot_telemetry(device_id, recorded_at DESC);
CREATE INDEX idx_iot_security ON device_iot_telemetry(security_status) WHERE security_status != 'locked';
CREATE INDEX idx_iot_health ON device_iot_telemetry(health_score) WHERE health_score < 50;
CREATE INDEX idx_iot_location ON device_iot_telemetry USING gist(
    ll_to_earth(location_lat, location_lon)
) WHERE location_lat IS NOT NULL;
