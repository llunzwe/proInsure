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
    
    -- Claims validation telemetry
    water_damage_sensor_status BOOLEAN, -- Liquid Contact Indicator
    water_damage_sensor_raw INTEGER, -- Raw sensor reading
    drop_detected BOOLEAN, -- Impact event flag
    drop_g_force DECIMAL(8,2), -- Impact severity
    drop_timestamp TIMESTAMPTZ,
    
    -- Battery diagnostics
    battery_cycle_count INTEGER,
    battery_maximum_capacity_percentage DECIMAL(5,2), -- % of original capacity
    battery_temperature_celsius DECIMAL(5,2),
    
    -- Screen time & usage (for wear-and-tear analysis)
    daily_screen_time_minutes INTEGER,
    app_usage_patterns JSONB, -- Categorized usage data
    
    -- Pre-claim diagnostic data
    diagnostic_timestamp TIMESTAMPTZ,
    camera_functional BOOLEAN,
    camera_test_result JSONB,
    biometric_sensor_status VARCHAR(50), -- face_id, touch_id, working, failed
    
    -- Geofencing & location
    geofence_violation BOOLEAN,
    country_code VARCHAR(2), -- ISO 3166 for coverage validation
    
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

-- Remote diagnostic API integration logs
CREATE TABLE remote_diagnostics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id),
    
    diagnostic_type VARCHAR(50) NOT NULL, -- apple_gsx, samsung_remote, oem_api
    api_endpoint VARCHAR(255),
    request_payload JSONB,
    response_payload JSONB,
    
    -- Diagnostic results
    overall_status VARCHAR(20), -- pass, fail, warning
    hardware_issues JSONB, -- Array of detected hardware problems
    software_issues JSONB,
    
    -- Timestamps
    requested_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    
    created_by UUID REFERENCES parties(id)
);

CREATE INDEX idx_remote_diag_device ON remote_diagnostics(device_id);
CREATE INDEX idx_remote_diag_type ON remote_diagnostics(diagnostic_type);

-- IoT data retention policies
CREATE TABLE iot_data_retention_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_name VARCHAR(100) NOT NULL UNIQUE,
    
    -- Retention rules
    hot_storage_days INTEGER DEFAULT 30, -- In TimescaleDB
    warm_storage_days INTEGER DEFAULT 90, -- Compressed
    cold_storage_days INTEGER DEFAULT 365, -- Archive
    deletion_after_days INTEGER DEFAULT 2555, -- 7 years
    
    -- Data classification
    data_type VARCHAR(50), -- location, health, diagnostics, usage
    jurisdiction VARCHAR(50), -- GDPR, CCPA, etc.
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Insert default retention policy
INSERT INTO iot_data_retention_policies (policy_name, data_type, hot_storage_days, warm_storage_days, cold_storage_days) VALUES
('default_location', 'location', 7, 30, 90),
('default_health', 'health', 30, 90, 365),
('default_diagnostics', 'diagnostics', 90, 365, 2555),
('default_usage', 'usage', 30, 90, 365);

-- Device end-of-support tracking
CREATE TABLE device_end_of_support (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    
    -- Support dates
    os_launch_date DATE,
    os_support_end_date DATE,
    security_update_end_date DATE,
    
    -- Insurance implications
    insurance_value_impact DECIMAL(5,2) DEFAULT 0.00, -- % reduction at support end
    repair_parts_availability VARCHAR(20) DEFAULT 'available' CHECK (repair_parts_availability IN ('available', 'limited', 'unavailable')),
    
    notification_sent_at TIMESTAMPTZ,
    
    UNIQUE(brand, model)
);

CREATE INDEX idx_eos_dates ON device_end_of_support(os_support_end_date) WHERE os_support_end_date > CURRENT_DATE;
