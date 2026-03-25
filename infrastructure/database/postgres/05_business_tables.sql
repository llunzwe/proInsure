-- SmartSure PostgreSQL Database Schema - Business Tables
-- Enterprise-grade smartphone insurance blockchain system
-- Supporting all 30 business options from the documentation

-- =============================================================================
-- INSURANCE BUSINESS TABLES
-- =============================================================================

-- Underwriting Policies table (Business Option 1)
CREATE TABLE underwriting_policies (
    underwriting_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(policy_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Underwriting Details
    underwriting_type VARCHAR(50) NOT NULL, -- risk_based, automated, manual
    risk_assessment_score DECIMAL(5,4),
    risk_factors JSONB,
    underwriting_decision VARCHAR(50), -- approved, declined, pending, conditional
    underwriting_notes TEXT,
    underwriting_conditions JSONB,
    
    -- Risk Assessment
    device_risk_score DECIMAL(5,4),
    customer_risk_score DECIMAL(5,4),
    environmental_risk_score DECIMAL(5,4),
    behavioral_risk_score DECIMAL(5,4),
    overall_risk_score DECIMAL(5,4),
    
    -- Premium Calculation
    base_premium_rate DECIMAL(5,4),
    risk_adjustment_factor DECIMAL(5,4),
    final_premium_rate DECIMAL(5,4),
    premium_calculation_method VARCHAR(100),
    premium_factors JSONB,
    
    -- Underwriter Information
    underwriter_id UUID,
    underwriter_name VARCHAR(200),
    underwriting_date TIMESTAMP,
    review_date TIMESTAMP,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_underwriting_type CHECK (underwriting_type IN ('risk_based', 'automated', 'manual', 'ai_assisted')),
    CONSTRAINT chk_underwriting_decision CHECK (underwriting_decision IN ('approved', 'declined', 'pending', 'conditional', 'referred'))
);

-- Reseller Policies table (Business Option 2)
CREATE TABLE reseller_policies (
    reseller_policy_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(policy_id),
    reseller_id UUID NOT NULL REFERENCES customers(customer_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Reseller Information
    reseller_type VARCHAR(50), -- b2b, b2c, white_label
    reseller_commission_rate DECIMAL(5,4),
    reseller_commission_amount DECIMAL(15,2),
    reseller_tier VARCHAR(50), -- bronze, silver, gold, platinum
    reseller_agreement_id UUID,
    
    -- Partnership Details
    partnership_type VARCHAR(50),
    revenue_sharing_model VARCHAR(50),
    revenue_split_percentage DECIMAL(5,4),
    white_label_branding JSONB,
    
    -- Sales Information
    sales_channel VARCHAR(100),
    sales_representative_id UUID,
    sales_date TIMESTAMP,
    sales_location JSONB,
    
    -- Commission Tracking
    commission_paid BOOLEAN DEFAULT false,
    commission_payment_date TIMESTAMP,
    commission_payment_method VARCHAR(50),
    commission_transaction_id VARCHAR(255),
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_reseller_type CHECK (reseller_type IN ('b2b', 'b2c', 'white_label', 'affiliate', 'broker')),
    CONSTRAINT chk_reseller_tier CHECK (reseller_tier IN ('bronze', 'silver', 'gold', 'platinum', 'diamond'))
);

-- Protection Plans table (Business Option 3)
CREATE TABLE protection_plans (
    protection_plan_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(policy_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Plan Details
    plan_type VARCHAR(50) NOT NULL, -- basic, standard, premium, ultimate
    plan_name VARCHAR(200),
    plan_description TEXT,
    subscription_frequency VARCHAR(20), -- monthly, quarterly, annually
    auto_renewal BOOLEAN DEFAULT true,
    
    -- Coverage Details
    coverage_limits JSONB,
    coverage_features TEXT[],
    coverage_exclusions TEXT[],
    additional_benefits JSONB,
    
    -- Pricing
    monthly_premium DECIMAL(15,2),
    quarterly_premium DECIMAL(15,2),
    annual_premium DECIMAL(15,2),
    setup_fee DECIMAL(15,2),
    cancellation_fee DECIMAL(15,2),
    
    -- Subscription Management
    subscription_status VARCHAR(50), -- active, paused, cancelled, expired
    next_billing_date TIMESTAMP,
    last_billing_date TIMESTAMP,
    billing_cycle_count INTEGER DEFAULT 0,
    
    -- Usage Tracking
    usage_metrics JSONB,
    usage_limits JSONB,
    usage_alerts JSONB,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_protection_plan_type CHECK (plan_type IN ('basic', 'standard', 'premium', 'ultimate', 'custom')),
    CONSTRAINT chk_subscription_frequency CHECK (subscription_frequency IN ('monthly', 'quarterly', 'annually', 'one_time')),
    CONSTRAINT chk_subscription_status CHECK (subscription_status IN ('active', 'paused', 'cancelled', 'expired', 'suspended'))
);

-- IMEI Validation table (Business Option 4)
CREATE TABLE imei_validations (
    imei_validation_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    imei VARCHAR(15) NOT NULL,
    device_id UUID REFERENCES devices(device_id),
    customer_id UUID REFERENCES customers(customer_id),
    
    -- Validation Details
    validation_status VARCHAR(50) NOT NULL, -- valid, invalid, blacklisted, stolen
    validation_score DECIMAL(5,4),
    validation_method VARCHAR(100), -- gsma, manufacturer, third_party
    validation_source VARCHAR(200),
    
    -- Fraud Detection
    fraud_risk_score DECIMAL(5,4),
    fraud_risk_level VARCHAR(20), -- low, medium, high, critical
    fraud_indicators JSONB,
    fraud_alerts JSONB,
    
    -- Device Information
    device_info JSONB,
    manufacturer_info JSONB,
    model_info JSONB,
    warranty_info JSONB,
    
    -- Blacklist Information
    blacklist_status BOOLEAN DEFAULT false,
    blacklist_reason TEXT,
    blacklist_date TIMESTAMP,
    blacklist_source VARCHAR(200),
    
    -- Theft Information
    theft_status BOOLEAN DEFAULT false,
    theft_report_date TIMESTAMP,
    theft_report_source VARCHAR(200),
    theft_report_details JSONB,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    validation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_imei_validation_status CHECK (validation_status IN ('valid', 'invalid', 'blacklisted', 'stolen', 'unknown')),
    CONSTRAINT chk_fraud_risk_level CHECK (fraud_risk_level IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT chk_imei_length CHECK (LENGTH(imei) = 15)
);

-- Device Tracking table (Business Option 5)
CREATE TABLE device_tracking (
    tracking_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    
    -- Location Information
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    altitude DECIMAL(10,2),
    accuracy DECIMAL(5,2),
    location_source VARCHAR(50), -- gps, wifi, cellular, manual
    
    -- Movement Information
    speed DECIMAL(10,2),
    direction DECIMAL(5,2),
    heading DECIMAL(5,2),
    movement_status VARCHAR(50), -- stationary, moving, unknown
    
    -- Tracking Configuration
    tracking_enabled BOOLEAN DEFAULT true,
    tracking_frequency INTERVAL,
    geofencing_enabled BOOLEAN DEFAULT false,
    geofence_radius DECIMAL(10,2),
    geofence_center JSONB,
    
    -- Alerts and Notifications
    alert_settings JSONB,
    notification_preferences JSONB,
    last_alert_sent TIMESTAMP,
    alert_history JSONB,
    
    -- Recovery Information
    recovery_mode BOOLEAN DEFAULT false,
    recovery_activated TIMESTAMP,
    recovery_contact_info JSONB,
    recovery_instructions TEXT,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    location_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_location_source CHECK (location_source IN ('gps', 'wifi', 'cellular', 'manual', 'bluetooth')),
    CONSTRAINT chk_movement_status CHECK (movement_status IN ('stationary', 'moving', 'unknown', 'error'))
); 