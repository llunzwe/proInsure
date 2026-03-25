-- SmartSure PostgreSQL Database Schema - Core Tables
-- Enterprise-grade smartphone insurance blockchain system
-- Generated from domain entities in internal/domain/entities/

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable JSONB for better performance with JSON data
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =============================================================================
-- CORE TABLES
-- =============================================================================

-- Customers table
CREATE TABLE customers (
    customer_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_number VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    
    -- Personal Information
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20),
    nationality VARCHAR(100),
    marital_status VARCHAR(50),
    occupation VARCHAR(200),
    employer VARCHAR(200),
    annual_income DECIMAL(15,2),
    income_currency VARCHAR(3) DEFAULT 'USD',
    education VARCHAR(200),
    
    -- Emergency Contact
    emergency_contact_name VARCHAR(200),
    emergency_contact_relationship VARCHAR(100),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_email VARCHAR(255),
    emergency_contact_address TEXT,
    
    -- Address Information
    address_type VARCHAR(50) DEFAULT 'primary',
    street_address TEXT NOT NULL,
    street_address2 TEXT,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) NOT NULL,
    is_primary_address BOOLEAN DEFAULT true,
    is_verified_address BOOLEAN DEFAULT false,
    address_verification_date TIMESTAMP,
    
    -- Identity Information
    identity_type VARCHAR(50),
    identity_number VARCHAR(100),
    issuing_country VARCHAR(100),
    issuing_authority VARCHAR(200),
    identity_issue_date DATE,
    identity_expiry_date DATE,
    is_identity_verified BOOLEAN DEFAULT false,
    identity_verification_date TIMESTAMP,
    identity_document_images JSONB,
    
    -- Biometric Data
    fingerprint_hash VARCHAR(255),
    face_recognition_id VARCHAR(255),
    iris_scan_hash VARCHAR(255),
    voice_print_hash VARCHAR(255),
    biometric_enrollment_date TIMESTAMP,
    biometric_last_updated TIMESTAMP,
    biometric_quality DECIMAL(3,2),
    
    -- Account Information
    account_status VARCHAR(50) DEFAULT 'pending',
    account_type VARCHAR(50) DEFAULT 'individual',
    risk_level VARCHAR(20) DEFAULT 'medium',
    risk_score DECIMAL(5,4) DEFAULT 0.5,
    risk_assessment_date TIMESTAMP,
    next_risk_assessment TIMESTAMP,
    risk_category VARCHAR(50),
    risk_mitigation_plan TEXT,
    
    -- Credit Information
    credit_score INTEGER,
    credit_score_range VARCHAR(20),
    credit_bureau VARCHAR(100),
    credit_last_updated TIMESTAMP,
    credit_history JSONB,
    credit_factors JSONB,
    credit_limitations TEXT[],
    
    -- Security Information
    password_hash VARCHAR(255),
    password_salt VARCHAR(255),
    password_last_changed TIMESTAMP,
    password_expiry TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0,
    last_failed_login TIMESTAMP,
    account_locked BOOLEAN DEFAULT false,
    lock_reason TEXT,
    lock_expiry TIMESTAMP,
    security_questions JSONB,
    two_factor_enabled BOOLEAN DEFAULT false,
    two_factor_method VARCHAR(50),
    
    -- Authentication Information
    last_login_method VARCHAR(50),
    last_login_ip VARCHAR(45),
    last_login_location JSONB,
    last_login_device VARCHAR(200),
    session_timeout INTERVAL,
    max_sessions INTEGER DEFAULT 5,
    active_sessions JSONB,
    login_history JSONB,
    
    -- Access Control
    roles JSONB,
    permissions JSONB,
    access_level VARCHAR(50),
    restricted_features TEXT[],
    access_history JSONB,
    last_access_review TIMESTAMP,
    next_access_review TIMESTAMP,
    
    -- Compliance Information
    gdpr_compliant BOOLEAN DEFAULT false,
    ccpa_compliant BOOLEAN DEFAULT false,
    sox_compliant BOOLEAN DEFAULT false,
    pci_compliant BOOLEAN DEFAULT false,
    last_audit_date TIMESTAMP,
    next_audit_date TIMESTAMP,
    compliance_score DECIMAL(5,4),
    compliance_violations JSONB,
    remediation_plan TEXT,
    
    -- Regulatory Flags
    regulatory_flags JSONB,
    
    -- Data Protection
    encryption_level VARCHAR(50),
    data_retention INTERVAL,
    data_residency VARCHAR(100),
    consent_status JSONB,
    data_subject_rights TEXT[],
    last_consent_update TIMESTAMP,
    data_breach_plan TEXT,
    privacy_impact_assessment BOOLEAN DEFAULT false,
    
    -- KYC Information
    kyc_status VARCHAR(50) DEFAULT 'pending',
    kyc_level VARCHAR(20),
    kyc_verification_date TIMESTAMP,
    kyc_expiry_date TIMESTAMP,
    kyc_documents JSONB,
    kyc_verification_method VARCHAR(100),
    kyc_verifier_id UUID,
    kyc_notes TEXT,
    
    -- Fraud Detection
    fraud_risk_score DECIMAL(5,4),
    fraud_risk_level VARCHAR(20),
    fraud_detection_model VARCHAR(100),
    fraud_last_assessment TIMESTAMP,
    fraud_alerts JSONB,
    fraud_false_positives INTEGER DEFAULT 0,
    fraud_true_positives INTEGER DEFAULT 0,
    fraud_model_accuracy DECIMAL(5,4),
    
    -- Behavioral Analysis
    behavioral_score DECIMAL(5,4),
    behavioral_risk_indicators JSONB,
    behavioral_pattern_analysis JSONB,
    behavioral_last_analysis TIMESTAMP,
    behavioral_next_analysis TIMESTAMP,
    behavioral_anomalies JSONB,
    
    -- Business Rules
    business_rules JSONB,
    restrictions JSONB,
    
    -- Customer Preferences
    language VARCHAR(10) DEFAULT 'en',
    currency VARCHAR(3) DEFAULT 'USD',
    timezone VARCHAR(50),
    notification_preferences JSONB,
    privacy_settings JSONB,
    marketing_preferences JSONB,
    
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
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_account_status CHECK (account_status IN ('active', 'inactive', 'suspended', 'pending', 'verified', 'unverified', 'blocked', 'under_review')),
    CONSTRAINT chk_account_type CHECK (account_type IN ('individual', 'business', 'corporate', 'premium', 'vip')),
    CONSTRAINT chk_risk_level CHECK (risk_level IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT chk_kyc_status CHECK (kyc_status IN ('pending', 'verified', 'rejected', 'expired')),
    CONSTRAINT chk_fraud_risk_level CHECK (fraud_risk_level IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT chk_credit_score CHECK (credit_score >= 300 AND credit_score <= 850),
    CONSTRAINT chk_annual_income CHECK (annual_income >= 0),
    CONSTRAINT chk_date_of_birth CHECK (date_of_birth <= CURRENT_DATE - INTERVAL '18 years')
);

-- Devices table
CREATE TABLE devices (
    device_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    imei VARCHAR(15) UNIQUE NOT NULL,
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    model VARCHAR(200) NOT NULL,
    manufacturer VARCHAR(200) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    
    -- Device Specifications
    processor VARCHAR(200),
    ram INTEGER,
    storage INTEGER,
    screen_size DECIMAL(4,2),
    resolution VARCHAR(50),
    battery_capacity INTEGER,
    camera_mp INTEGER,
    operating_system VARCHAR(100),
    os_version VARCHAR(50),
    network_type VARCHAR(50),
    color VARCHAR(50),
    weight DECIMAL(6,2),
    dimensions JSONB,
    
    -- Device Condition
    overall_condition VARCHAR(50),
    physical_condition VARCHAR(50),
    functional_condition VARCHAR(50),
    cosmetic_condition VARCHAR(50),
    battery_health DECIMAL(3,2),
    screen_condition VARCHAR(50),
    body_condition VARCHAR(50),
    water_damage BOOLEAN DEFAULT false,
    repair_history JSONB,
    last_assessment TIMESTAMP,
    assessed_by UUID,
    grade VARCHAR(5),
    grade_score DECIMAL(3,2),
    
    -- Device Status
    status VARCHAR(50) DEFAULT 'active',
    
    -- Ownership Information
    current_owner_id UUID NOT NULL REFERENCES customers(customer_id),
    previous_owners JSONB,
    ownership_type VARCHAR(50),
    purchase_price DECIMAL(15,2),
    purchase_location VARCHAR(200),
    transfer_history JSONB,
    lien_holder UUID,
    lien_amount DECIMAL(15,2),
    
    -- Value Information
    original_value DECIMAL(15,2) NOT NULL,
    current_value DECIMAL(15,2),
    depreciation_rate DECIMAL(5,4),
    insurance_value DECIMAL(15,2),
    market_value DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- IoT Integration
    iot_enabled BOOLEAN DEFAULT false,
    iot_device_sensors TEXT[],
    iot_monitoring_frequency INTERVAL,
    iot_alert_thresholds JSONB,
    iot_last_data_sync TIMESTAMP,
    iot_data_quality DECIMAL(3,2),
    iot_integration_status VARCHAR(50),
    iot_provider VARCHAR(100),
    iot_data_points JSONB,
    
    -- Security Features
    biometric_auth BOOLEAN DEFAULT false,
    fingerprint BOOLEAN DEFAULT false,
    face_recognition BOOLEAN DEFAULT false,
    iris_scan BOOLEAN DEFAULT false,
    encryption BOOLEAN DEFAULT false,
    encryption_level VARCHAR(50),
    remote_wipe BOOLEAN DEFAULT false,
    find_my_device BOOLEAN DEFAULT false,
    two_factor_auth BOOLEAN DEFAULT false,
    last_security_update TIMESTAMP,
    security_score DECIMAL(3,2),
    
    -- Tracking Information
    last_known_location JSONB,
    last_seen TIMESTAMP,
    tracking_enabled BOOLEAN DEFAULT false,
    tracking_provider VARCHAR(100),
    location_history JSONB,
    movement_pattern VARCHAR(50),
    geofencing JSONB,
    alert_settings JSONB,
    
    -- Recovery Information
    recovery_plan JSONB,
    recovery_status VARCHAR(50),
    recovery_attempts JSONB,
    recovery_reward_amount DECIMAL(15,2),
    recovery_contact_info JSONB,
    last_recovery_attempt TIMESTAMP,
    recovery_probability DECIMAL(3,2),
    
    -- Warranty Information
    warranty_info JSONB,
    
    -- Business Rules
    business_rules JSONB,
    restrictions JSONB,
    
    -- Compliance Information
    gdpr_compliant BOOLEAN DEFAULT false,
    ccpa_compliant BOOLEAN DEFAULT false,
    sox_compliant BOOLEAN DEFAULT false,
    pci_compliant BOOLEAN DEFAULT false,
    last_audit_date TIMESTAMP,
    next_audit_date TIMESTAMP,
    compliance_score DECIMAL(5,4),
    compliance_violations JSONB,
    remediation_plan TEXT,
    
    -- Regulatory Flags
    regulatory_flags JSONB,
    
    -- Data Protection
    encryption_level_data VARCHAR(50),
    data_retention INTERVAL,
    data_residency VARCHAR(100),
    consent_status JSONB,
    data_subject_rights TEXT[],
    last_consent_update TIMESTAMP,
    data_breach_plan TEXT,
    privacy_impact_assessment BOOLEAN DEFAULT false,
    
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
    manufacture_date DATE,
    purchase_date DATE,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_device_status CHECK (status IN ('active', 'inactive', 'lost', 'stolen', 'damaged', 'repairing', 'recovered', 'disposed', 'under_review')),
    CONSTRAINT chk_imei_length CHECK (LENGTH(imei) = 15),
    CONSTRAINT chk_original_value CHECK (original_value > 0),
    CONSTRAINT chk_battery_health CHECK (battery_health >= 0 AND battery_health <= 1),
    CONSTRAINT chk_grade_score CHECK (grade_score >= 0 AND grade_score <= 1),
    CONSTRAINT chk_screen_size CHECK (screen_size > 0),
    CONSTRAINT chk_weight CHECK (weight > 0)
); 