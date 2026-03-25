-- SmartSure PostgreSQL Database Schema - Additional Tables
-- Enterprise-grade smartphone insurance blockchain system

-- =============================================================================
-- ADDITIONAL DOMAIN TABLES
-- =============================================================================

-- IoT Devices table
CREATE TABLE iot_devices (
    iot_device_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    sensor_type VARCHAR(100) NOT NULL,
    sensor_id VARCHAR(100) UNIQUE NOT NULL,
    sensor_model VARCHAR(200),
    manufacturer VARCHAR(200),
    firmware_version VARCHAR(50),
    hardware_version VARCHAR(50),
    
    -- Configuration
    enabled BOOLEAN DEFAULT true,
    calibration_status VARCHAR(50) DEFAULT 'pending',
    last_calibration_date TIMESTAMP,
    next_calibration_date TIMESTAMP,
    accuracy DECIMAL(5,4),
    precision DECIMAL(5,4),
    
    -- Data Collection
    sampling_rate DECIMAL(10,2), -- samples per second
    data_format VARCHAR(50),
    unit_of_measurement VARCHAR(20),
    min_value DECIMAL(15,4),
    max_value DECIMAL(15,4),
    threshold_high DECIMAL(15,4),
    threshold_low DECIMAL(15,4),
    
    -- Connectivity
    connection_type VARCHAR(50), -- wifi, bluetooth, cellular, etc.
    connection_status VARCHAR(50),
    last_connection TIMESTAMP,
    ip_address VARCHAR(45),
    mac_address VARCHAR(17),
    
    -- Power Management
    battery_level DECIMAL(3,2),
    power_source VARCHAR(50),
    power_consumption DECIMAL(10,4),
    last_power_check TIMESTAMP,
    
    -- Location
    installation_location JSONB,
    installation_date TIMESTAMP,
    installation_notes TEXT,
    
    -- Maintenance
    maintenance_schedule JSONB,
    last_maintenance_date TIMESTAMP,
    next_maintenance_date TIMESTAMP,
    maintenance_history JSONB,
    
    -- Status
    status VARCHAR(50) DEFAULT 'active',
    health_score DECIMAL(3,2),
    error_count INTEGER DEFAULT 0,
    last_error_date TIMESTAMP,
    error_log JSONB,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_chain_id VARCHAR(100),
    blockchain_network_id VARCHAR(100),
    blockchain_smart_contract VARCHAR(200),
    blockchain_gas_used BIGINT,
    blockchain_gas_price BIGINT,
    blockchain_status VARCHAR(50),
    blockchain_confirmations INTEGER,
    
    -- Audit Information
    audit_trail JSONB,
    provenance_log JSONB,
    
    -- Version Control
    version INTEGER DEFAULT 1,
    previous_version_id UUID,
    change_reason TEXT,
    
    -- Metadata
    tags JSONB,
    custom_fields JSONB,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_iot_device_status CHECK (status IN ('active', 'inactive', 'maintenance', 'error', 'offline')),
    CONSTRAINT chk_iot_device_accuracy CHECK (accuracy >= 0 AND accuracy <= 1),
    CONSTRAINT chk_iot_device_battery_level CHECK (battery_level >= 0 AND battery_level <= 1),
    CONSTRAINT chk_iot_device_health_score CHECK (health_score >= 0 AND health_score <= 1)
);

-- IoT Data Points table
CREATE TABLE iot_data_points (
    data_point_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    iot_device_id UUID NOT NULL REFERENCES iot_devices(iot_device_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Data Values
    sensor_value DECIMAL(15,4) NOT NULL,
    unit VARCHAR(20),
    data_type VARCHAR(50),
    data_quality DECIMAL(3,2),
    
    -- Location and Context
    location JSONB,
    accuracy DECIMAL(5,4),
    altitude DECIMAL(10,2),
    speed DECIMAL(10,2),
    direction DECIMAL(5,2),
    
    -- Environmental Context
    temperature DECIMAL(5,2),
    humidity DECIMAL(5,2),
    pressure DECIMAL(10,2),
    ambient_light DECIMAL(10,2),
    
    -- Processing
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMP,
    processing_notes TEXT,
    anomaly_detected BOOLEAN DEFAULT false,
    anomaly_score DECIMAL(5,4),
    anomaly_type VARCHAR(100),
    
    -- Metadata
    raw_data JSONB,
    processed_data JSONB,
    flags JSONB,
    
    -- Timestamps
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_iot_data_quality CHECK (data_quality >= 0 AND data_quality <= 1),
    CONSTRAINT chk_iot_anomaly_score CHECK (anomaly_score >= 0 AND anomaly_score <= 1)
);

-- Marketplace Listings table
CREATE TABLE marketplace_listings (
    listing_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    seller_id UUID NOT NULL REFERENCES customers(customer_id),
    
    -- Listing Details
    listing_type VARCHAR(50) NOT NULL, -- sale, trade, auction, etc.
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    subcategory VARCHAR(100),
    
    -- Pricing
    price DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    original_price DECIMAL(15,2),
    discount_percentage DECIMAL(5,2),
    negotiable BOOLEAN DEFAULT false,
    reserve_price DECIMAL(15,2),
    
    -- Auction Details (if applicable)
    auction_start_date TIMESTAMP,
    auction_end_date TIMESTAMP,
    minimum_bid DECIMAL(15,2),
    current_bid DECIMAL(15,2),
    bid_count INTEGER DEFAULT 0,
    winner_id UUID REFERENCES customers(customer_id),
    
    -- Condition and Details
    condition_rating DECIMAL(3,2),
    condition_description TEXT,
    included_items TEXT[],
    excluded_items TEXT[],
    warranty_included BOOLEAN DEFAULT false,
    warranty_details TEXT,
    
    -- Location and Shipping
    location JSONB,
    shipping_available BOOLEAN DEFAULT true,
    shipping_cost DECIMAL(10,2),
    shipping_methods TEXT[],
    pickup_available BOOLEAN DEFAULT true,
    pickup_location JSONB,
    
    -- Status and Visibility
    status VARCHAR(50) DEFAULT 'active',
    visibility VARCHAR(50) DEFAULT 'public',
    featured BOOLEAN DEFAULT false,
    promoted BOOLEAN DEFAULT false,
    
    -- Metrics
    view_count INTEGER DEFAULT 0,
    favorite_count INTEGER DEFAULT 0,
    inquiry_count INTEGER DEFAULT 0,
    last_viewed TIMESTAMP,
    
    -- Images and Media
    images TEXT[],
    videos TEXT[],
    documents TEXT[],
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_chain_id VARCHAR(100),
    blockchain_network_id VARCHAR(100),
    blockchain_smart_contract VARCHAR(200),
    blockchain_gas_used BIGINT,
    blockchain_gas_price BIGINT,
    blockchain_status VARCHAR(50),
    blockchain_confirmations INTEGER,
    
    -- Audit Information
    audit_trail JSONB,
    provenance_log JSONB,
    
    -- Version Control
    version INTEGER DEFAULT 1,
    previous_version_id UUID,
    change_reason TEXT,
    
    -- Metadata
    tags JSONB,
    custom_fields JSONB,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_marketplace_listing_type CHECK (listing_type IN ('sale', 'trade', 'auction', 'rental', 'lease')),
    CONSTRAINT chk_marketplace_status CHECK (status IN ('draft', 'active', 'sold', 'expired', 'cancelled', 'suspended')),
    CONSTRAINT chk_marketplace_visibility CHECK (visibility IN ('public', 'private', 'invite_only')),
    CONSTRAINT chk_marketplace_price CHECK (price > 0),
    CONSTRAINT chk_marketplace_condition_rating CHECK (condition_rating >= 0 AND condition_rating <= 1)
);

-- Repair Orders table
CREATE TABLE repair_orders (
    repair_order_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    claim_id UUID REFERENCES claims(claim_id),
    
    -- Repair Details
    repair_type VARCHAR(100) NOT NULL,
    repair_category VARCHAR(100),
    priority VARCHAR(20) DEFAULT 'normal',
    urgency VARCHAR(20) DEFAULT 'standard',
    
    -- Problem Description
    problem_description TEXT NOT NULL,
    symptoms TEXT[],
    error_codes TEXT[],
    diagnostic_results JSONB,
    
    -- Service Provider
    service_provider_id UUID,
    service_provider_name VARCHAR(200),
    service_provider_location JSONB,
    technician_id UUID,
    technician_name VARCHAR(200),
    technician_contact VARCHAR(100),
    
    -- Scheduling
    scheduled_date TIMESTAMP,
    estimated_duration INTERVAL,
    actual_start_date TIMESTAMP,
    actual_end_date TIMESTAMP,
    actual_duration INTERVAL,
    
    -- Status Tracking
    status VARCHAR(50) DEFAULT 'created',
    current_step VARCHAR(100),
    progress_percentage DECIMAL(5,2) DEFAULT 0,
    estimated_completion TIMESTAMP,
    
    -- Parts and Materials
    parts_required JSONB,
    parts_ordered JSONB,
    parts_received JSONB,
    parts_installed JSONB,
    materials_used JSONB,
    
    -- Cost Information
    estimated_cost DECIMAL(15,2),
    actual_cost DECIMAL(15,2),
    labor_cost DECIMAL(15,2),
    parts_cost DECIMAL(15,2),
    tax_amount DECIMAL(15,2),
    total_cost DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- Warranty and Insurance
    warranty_coverage BOOLEAN DEFAULT false,
    warranty_type VARCHAR(100),
    warranty_expiry_date DATE,
    insurance_coverage BOOLEAN DEFAULT false,
    insurance_policy_number VARCHAR(100),
    deductible_amount DECIMAL(15,2),
    
    -- Quality Assurance
    quality_check_performed BOOLEAN DEFAULT false,
    quality_check_date TIMESTAMP,
    quality_check_by UUID,
    quality_score DECIMAL(3,2),
    quality_notes TEXT,
    
    -- Customer Communication
    customer_notifications JSONB,
    customer_approvals JSONB,
    customer_feedback JSONB,
    customer_satisfaction_score DECIMAL(3,2),
    
    -- Documentation
    work_orders JSONB,
    invoices JSONB,
    receipts JSONB,
    warranties JSONB,
    photos TEXT[],
    videos TEXT[],
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_chain_id VARCHAR(100),
    blockchain_network_id VARCHAR(100),
    blockchain_smart_contract VARCHAR(200),
    blockchain_gas_used BIGINT,
    blockchain_gas_price BIGINT,
    blockchain_status VARCHAR(50),
    blockchain_confirmations INTEGER,
    
    -- Audit Information
    audit_trail JSONB,
    provenance_log JSONB,
    
    -- Version Control
    version INTEGER DEFAULT 1,
    previous_version_id UUID,
    change_reason TEXT,
    
    -- Metadata
    tags JSONB,
    custom_fields JSONB,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_repair_priority CHECK (priority IN ('low', 'normal', 'high', 'urgent', 'emergency')),
    CONSTRAINT chk_repair_urgency CHECK (urgency IN ('standard', 'express', 'same_day', 'next_day')),
    CONSTRAINT chk_repair_status CHECK (status IN ('created', 'scheduled', 'in_progress', 'parts_ordered', 'parts_received', 'completed', 'cancelled', 'on_hold')),
    CONSTRAINT chk_repair_progress CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    CONSTRAINT chk_repair_quality_score CHECK (quality_score >= 0 AND quality_score <= 1),
    CONSTRAINT chk_repair_customer_satisfaction CHECK (customer_satisfaction_score >= 0 AND customer_satisfaction_score <= 1)
);

-- Repair History table (for tracking individual repairs)
CREATE TABLE repair_history (
    repair_history_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repair_order_id UUID NOT NULL REFERENCES repair_orders(repair_order_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Repair Details
    repair_type VARCHAR(100) NOT NULL,
    repair_description TEXT,
    parts_replaced TEXT[],
    labor_hours DECIMAL(5,2),
    cost DECIMAL(15,2),
    
    -- Service Provider
    service_provider_name VARCHAR(200),
    technician_name VARCHAR(200),
    service_location JSONB,
    
    -- Dates
    repair_date TIMESTAMP NOT NULL,
    warranty_expiry_date DATE,
    
    -- Quality
    quality_rating DECIMAL(3,2),
    customer_satisfaction DECIMAL(3,2),
    notes TEXT,
    
    -- Documentation
    work_order_number VARCHAR(100),
    invoice_number VARCHAR(100),
    warranty_document TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_repair_history_quality_rating CHECK (quality_rating >= 0 AND quality_rating <= 1),
    CONSTRAINT chk_repair_history_customer_satisfaction CHECK (customer_satisfaction >= 0 AND customer_satisfaction <= 1)
);

-- =============================================================================
-- INDEXES FOR ADDITIONAL TABLES
-- =============================================================================

-- IoT Devices indexes
CREATE INDEX idx_iot_devices_device_id ON iot_devices(device_id);
CREATE INDEX idx_iot_devices_sensor_type ON iot_devices(sensor_type);
CREATE INDEX idx_iot_devices_status ON iot_devices(status);
CREATE INDEX idx_iot_devices_enabled ON iot_devices(enabled);
CREATE INDEX idx_iot_devices_last_connection ON iot_devices(last_connection);
CREATE INDEX idx_iot_devices_health_score ON iot_devices(health_score);

-- IoT Data Points indexes
CREATE INDEX idx_iot_data_points_iot_device_id ON iot_data_points(iot_device_id);
CREATE INDEX idx_iot_data_points_device_id ON iot_data_points(device_id);
CREATE INDEX idx_iot_data_points_timestamp ON iot_data_points(timestamp);
CREATE INDEX idx_iot_data_points_processed ON iot_data_points(processed);
CREATE INDEX idx_iot_data_points_anomaly_detected ON iot_data_points(anomaly_detected);
CREATE INDEX idx_iot_data_points_sensor_value ON iot_data_points(sensor_value);

-- Marketplace Listings indexes
CREATE INDEX idx_marketplace_listings_device_id ON marketplace_listings(device_id);
CREATE INDEX idx_marketplace_listings_seller_id ON marketplace_listings(seller_id);
CREATE INDEX idx_marketplace_listings_listing_type ON marketplace_listings(listing_type);
CREATE INDEX idx_marketplace_listings_status ON marketplace_listings(status);
CREATE INDEX idx_marketplace_listings_price ON marketplace_listings(price);
CREATE INDEX idx_marketplace_listings_category ON marketplace_listings(category);
CREATE INDEX idx_marketplace_listings_created_at ON marketplace_listings(created_at);
CREATE INDEX idx_marketplace_listings_expires_at ON marketplace_listings(expires_at);
CREATE INDEX idx_marketplace_listings_featured ON marketplace_listings(featured);

-- Repair Orders indexes
CREATE INDEX idx_repair_orders_device_id ON repair_orders(device_id);
CREATE INDEX idx_repair_orders_customer_id ON repair_orders(customer_id);
CREATE INDEX idx_repair_orders_claim_id ON repair_orders(claim_id);
CREATE INDEX idx_repair_orders_status ON repair_orders(status);
CREATE INDEX idx_repair_orders_repair_type ON repair_orders(repair_type);
CREATE INDEX idx_repair_orders_priority ON repair_orders(priority);
CREATE INDEX idx_repair_orders_scheduled_date ON repair_orders(scheduled_date);
CREATE INDEX idx_repair_orders_created_at ON repair_orders(created_at);
CREATE INDEX idx_repair_orders_service_provider_id ON repair_orders(service_provider_id);

-- Repair History indexes
CREATE INDEX idx_repair_history_repair_order_id ON repair_history(repair_order_id);
CREATE INDEX idx_repair_history_device_id ON repair_history(device_id);
CREATE INDEX idx_repair_history_repair_date ON repair_history(repair_date);
CREATE INDEX idx_repair_history_repair_type ON repair_history(repair_type);

-- Composite indexes
CREATE INDEX idx_iot_data_points_device_timestamp ON iot_data_points(device_id, timestamp);
CREATE INDEX idx_marketplace_listings_type_status ON marketplace_listings(listing_type, status);
CREATE INDEX idx_repair_orders_device_status ON repair_orders(device_id, status);

-- Partial indexes
CREATE INDEX idx_iot_devices_active ON iot_devices(device_id, sensor_type) WHERE enabled = true AND status = 'active';
CREATE INDEX idx_marketplace_listings_active ON marketplace_listings(category, price) WHERE status = 'active';
CREATE INDEX idx_repair_orders_open ON repair_orders(device_id, priority) WHERE status IN ('created', 'scheduled', 'in_progress');

-- =============================================================================
-- VIEWS FOR ADDITIONAL TABLES
-- =============================================================================

-- IoT Device Analytics view
CREATE VIEW iot_device_analytics AS
SELECT 
    id.iot_device_id,
    id.device_id,
    d.imei,
    d.model as device_model,
    id.sensor_type,
    id.sensor_id,
    id.status as iot_status,
    id.enabled,
    id.health_score,
    id.accuracy,
    id.last_connection,
    id.battery_level,
    COUNT(idp.data_point_id) as total_data_points,
    COUNT(CASE WHEN idp.processed = true THEN 1 END) as processed_data_points,
    COUNT(CASE WHEN idp.anomaly_detected = true THEN 1 END) as anomaly_count,
    AVG(idp.sensor_value) as avg_sensor_value,
    MIN(idp.sensor_value) as min_sensor_value,
    MAX(idp.sensor_value) as max_sensor_value,
    AVG(idp.data_quality) as avg_data_quality,
    AVG(idp.anomaly_score) as avg_anomaly_score,
    MAX(idp.timestamp) as last_data_point
FROM iot_devices id
JOIN devices d ON id.device_id = d.device_id
LEFT JOIN iot_data_points idp ON id.iot_device_id = idp.iot_device_id
GROUP BY id.iot_device_id, id.device_id, d.imei, d.model, id.sensor_type, id.sensor_id,
         id.status, id.enabled, id.health_score, id.accuracy, id.last_connection, id.battery_level;

-- Marketplace Analytics view
CREATE VIEW marketplace_analytics AS
SELECT 
    ml.listing_id,
    ml.device_id,
    d.imei,
    d.model as device_model,
    d.manufacturer,
    ml.seller_id,
    c.first_name || ' ' || c.last_name as seller_name,
    ml.listing_type,
    ml.category,
    ml.status,
    ml.price,
    ml.currency,
    ml.condition_rating,
    ml.view_count,
    ml.favorite_count,
    ml.inquiry_count,
    ml.featured,
    ml.promoted,
    ml.created_at,
    ml.expires_at,
    CASE 
        WHEN ml.expires_at < CURRENT_TIMESTAMP THEN 'expired'
        WHEN ml.status = 'active' THEN 'active'
        ELSE ml.status
    END as effective_status
FROM marketplace_listings ml
JOIN devices d ON ml.device_id = d.device_id
JOIN customers c ON ml.seller_id = c.customer_id;

-- Repair Analytics view
CREATE VIEW repair_analytics AS
SELECT 
    ro.repair_order_id,
    ro.device_id,
    d.imei,
    d.model as device_model,
    ro.customer_id,
    c.first_name || ' ' || c.last_name as customer_name,
    ro.claim_id,
    ro.repair_type,
    ro.repair_category,
    ro.priority,
    ro.status,
    ro.estimated_cost,
    ro.actual_cost,
    ro.total_cost,
    ro.currency,
    ro.estimated_duration,
    ro.actual_duration,
    ro.progress_percentage,
    ro.quality_score,
    ro.customer_satisfaction_score,
    ro.scheduled_date,
    ro.actual_start_date,
    ro.actual_end_date,
    ro.created_at,
    ro.updated_at,
    CASE 
        WHEN ro.actual_duration IS NOT NULL AND ro.estimated_duration IS NOT NULL THEN
            EXTRACT(EPOCH FROM (ro.actual_duration - ro.estimated_duration)) / 3600
        ELSE NULL
    END as duration_variance_hours
FROM repair_orders ro
JOIN devices d ON ro.device_id = d.device_id
JOIN customers c ON ro.customer_id = c.customer_id;

-- =============================================================================
-- FUNCTIONS FOR ADDITIONAL TABLES
-- =============================================================================

-- Function to calculate IoT device health score
CREATE OR REPLACE FUNCTION calculate_iot_device_health(
    iot_device_uuid UUID
)
RETURNS DECIMAL AS $$
DECLARE
    health_score DECIMAL := 1.0;
    connection_factor DECIMAL;
    battery_factor DECIMAL;
    error_factor DECIMAL;
    data_quality_factor DECIMAL;
    device_record RECORD;
BEGIN
    -- Get device information
    SELECT 
        last_connection,
        battery_level,
        error_count,
        accuracy
    INTO device_record
    FROM iot_devices 
    WHERE iot_device_id = iot_device_uuid;
    
    -- Connection factor (recent connection = better health)
    IF device_record.last_connection IS NULL THEN
        connection_factor := 0.0;
    ELSIF device_record.last_connection > CURRENT_TIMESTAMP - INTERVAL '1 hour' THEN
        connection_factor := 1.0;
    ELSIF device_record.last_connection > CURRENT_TIMESTAMP - INTERVAL '24 hours' THEN
        connection_factor := 0.8;
    ELSIF device_record.last_connection > CURRENT_TIMESTAMP - INTERVAL '7 days' THEN
        connection_factor := 0.5;
    ELSE
        connection_factor := 0.2;
    END IF;
    
    -- Battery factor
    IF device_record.battery_level IS NULL THEN
        battery_factor := 0.5;
    ELSE
        battery_factor := device_record.battery_level;
    END IF;
    
    -- Error factor (fewer errors = better health)
    IF device_record.error_count = 0 THEN
        error_factor := 1.0;
    ELSIF device_record.error_count <= 5 THEN
        error_factor := 0.8;
    ELSIF device_record.error_count <= 10 THEN
        error_factor := 0.6;
    ELSIF device_record.error_count <= 20 THEN
        error_factor := 0.4;
    ELSE
        error_factor := 0.2;
    END IF;
    
    -- Data quality factor
    IF device_record.accuracy IS NULL THEN
        data_quality_factor := 0.5;
    ELSE
        data_quality_factor := device_record.accuracy;
    END IF;
    
    -- Calculate overall health score
    health_score := (connection_factor * 0.3 + battery_factor * 0.2 + 
                    error_factor * 0.3 + data_quality_factor * 0.2);
    
    RETURN LEAST(1.0, GREATEST(0.0, health_score));
END;
$$ LANGUAGE plpgsql;

-- Function to detect IoT data anomalies
CREATE OR REPLACE FUNCTION detect_iot_anomaly(
    sensor_value DECIMAL,
    iot_device_uuid UUID
)
RETURNS BOOLEAN AS $$
DECLARE
    threshold_high DECIMAL;
    threshold_low DECIMAL;
    avg_value DECIMAL;
    std_dev DECIMAL;
    anomaly_threshold DECIMAL := 2.0; -- 2 standard deviations
BEGIN
    -- Get device thresholds
    SELECT threshold_high, threshold_low 
    INTO threshold_high, threshold_low
    FROM iot_devices 
    WHERE iot_device_id = iot_device_uuid;
    
    -- Check against device-specific thresholds
    IF threshold_high IS NOT NULL AND sensor_value > threshold_high THEN
        RETURN TRUE;
    END IF;
    
    IF threshold_low IS NOT NULL AND sensor_value < threshold_low THEN
        RETURN TRUE;
    END IF;
    
    -- Check against statistical thresholds (last 100 data points)
    SELECT AVG(sensor_value), STDDEV(sensor_value)
    INTO avg_value, std_dev
    FROM iot_data_points
    WHERE iot_device_id = iot_device_uuid
    AND timestamp > CURRENT_TIMESTAMP - INTERVAL '24 hours'
    LIMIT 100;
    
    IF avg_value IS NOT NULL AND std_dev IS NOT NULL THEN
        IF ABS(sensor_value - avg_value) > (std_dev * anomaly_threshold) THEN
            RETURN TRUE;
        END IF;
    END IF;
    
    RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

-- Function to calculate repair cost estimate
CREATE OR REPLACE FUNCTION estimate_repair_cost(
    device_uuid UUID,
    repair_type VARCHAR,
    repair_category VARCHAR
)
RETURNS DECIMAL AS $$
DECLARE
    base_cost DECIMAL := 0;
    device_value DECIMAL;
    complexity_factor DECIMAL := 1.0;
    labor_rate DECIMAL := 75.0; -- per hour
    estimated_hours DECIMAL := 2.0;
BEGIN
    -- Get device value
    SELECT current_value INTO device_value
    FROM devices
    WHERE device_id = device_uuid;
    
    -- Base cost by repair type
    CASE repair_type
        WHEN 'screen_replacement' THEN
            base_cost := device_value * 0.3;
            estimated_hours := 1.5;
        WHEN 'battery_replacement' THEN
            base_cost := 50.0;
            estimated_hours := 1.0;
        WHEN 'water_damage' THEN
            base_cost := device_value * 0.5;
            estimated_hours := 3.0;
            complexity_factor := 1.5;
        WHEN 'motherboard_repair' THEN
            base_cost := device_value * 0.7;
            estimated_hours := 4.0;
            complexity_factor := 2.0;
        WHEN 'software_repair' THEN
            base_cost := 25.0;
            estimated_hours := 0.5;
        ELSE
            base_cost := device_value * 0.2;
            estimated_hours := 2.0;
    END CASE;
    
    -- Calculate total cost
    RETURN base_cost + (labor_rate * estimated_hours * complexity_factor);
END;
$$ LANGUAGE plpgsql; 