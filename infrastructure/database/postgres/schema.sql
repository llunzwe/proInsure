-- SmartSure PostgreSQL Database Schema
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

-- Policies table
CREATE TABLE policies (
    policy_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Policy Details
    policy_type VARCHAR(50) NOT NULL,
    coverage_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    
    -- Premium Information
    base_premium DECIMAL(15,2) NOT NULL,
    risk_premium DECIMAL(15,2),
    loading_premium DECIMAL(15,2),
    discount_amount DECIMAL(15,2),
    final_premium DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_frequency VARCHAR(20),
    payment_method VARCHAR(50),
    next_payment_date TIMESTAMP,
    premium_history JSONB,
    
    -- Coverage Information
    coverage_limits JSONB,
    deductibles JSONB,
    coinsurance JSONB,
    coverage_period INTERVAL,
    geographic_scope JSONB,
    usage_restrictions JSONB,
    
    -- Risk Assessment
    risk_score DECIMAL(5,4),
    risk_factors JSONB,
    underwriting_notes JSONB,
    
    -- IoT Integration
    iot_enabled BOOLEAN DEFAULT false,
    iot_device_sensors TEXT[],
    iot_monitoring_frequency INTERVAL,
    iot_alert_thresholds JSONB,
    iot_last_data_sync TIMESTAMP,
    iot_data_quality DECIMAL(3,2),
    iot_integration_status VARCHAR(50),
    iot_provider VARCHAR(100),
    
    -- Fraud Detection
    fraud_risk_score DECIMAL(5,4),
    fraud_risk_level VARCHAR(20),
    fraud_detection_model VARCHAR(100),
    fraud_last_assessment TIMESTAMP,
    fraud_alerts JSONB,
    fraud_false_positives INTEGER DEFAULT 0,
    fraud_true_positives INTEGER DEFAULT 0,
    fraud_model_accuracy DECIMAL(5,4),
    
    -- Claims History
    claims_history JSONB,
    
    -- Business Rules
    business_rules JSONB,
    exclusions JSONB,
    endorsements JSONB,
    
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
    effective_date TIMESTAMP NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_policy_type CHECK (policy_type IN ('comprehensive', 'accident', 'theft', 'liability', 'extended_warranty', 'micro_insurance', 'peer_to_peer', 'cooperative')),
    CONSTRAINT chk_coverage_type CHECK (coverage_type IN ('physical_damage', 'liquid_damage', 'theft', 'loss', 'liability', 'data_recovery', 'repair', 'replacement')),
    CONSTRAINT chk_policy_status CHECK (status IN ('draft', 'pending', 'active', 'suspended', 'cancelled', 'expired', 'under_review', 'declined')),
    CONSTRAINT chk_payment_frequency CHECK (payment_frequency IN ('monthly', 'quarterly', 'annually', 'one_time')),
    CONSTRAINT chk_payment_method CHECK (payment_method IN ('credit_card', 'debit_card', 'bank_transfer', 'crypto', 'mobile_money')),
    CONSTRAINT chk_fraud_risk_level CHECK (fraud_risk_level IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT chk_effective_date CHECK (effective_date < expiration_date),
    CONSTRAINT chk_final_premium CHECK (final_premium > 0)
);

-- Claims table
CREATE TABLE claims (
    claim_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_number VARCHAR(50) UNIQUE NOT NULL,
    policy_id UUID NOT NULL REFERENCES policies(policy_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Claim Details
    claim_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'filed',
    severity VARCHAR(20),
    priority VARCHAR(20),
    
    -- Incident Information
    incident_date TIMESTAMP NOT NULL,
    reported_date TIMESTAMP NOT NULL,
    incident_location JSONB,
    incident_description TEXT NOT NULL,
    incident_type VARCHAR(100),
    
    -- Financial Information
    claimed_amount DECIMAL(15,2) NOT NULL,
    approved_amount DECIMAL(15,2),
    paid_amount DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    deductible DECIMAL(15,2),
    coinsurance DECIMAL(5,4),
    
    -- Assessment Information
    assessment_id UUID,
    assessor_id UUID,
    assessment_date TIMESTAMP,
    damage_type VARCHAR(100),
    damage_severity VARCHAR(50),
    repair_cost DECIMAL(15,2),
    replace_cost DECIMAL(15,2),
    recommended_action VARCHAR(200),
    assessment_notes TEXT,
    assessment_photos TEXT[],
    assessment_videos TEXT[],
    assessment_documents TEXT[],
    
    -- Investigation Information
    investigation_id UUID,
    investigator_id UUID,
    investigation_start_date TIMESTAMP,
    investigation_end_date TIMESTAMP,
    investigation_status VARCHAR(50),
    investigation_findings TEXT,
    investigation_conclusion TEXT,
    investigation_recommendations TEXT[],
    investigation_evidence JSONB,
    investigation_witness_statements JSONB,
    
    -- Evidence
    evidence JSONB,
    witnesses JSONB,
    
    -- Processing Information
    assigned_to UUID,
    assigned_at TIMESTAMP,
    processing_time INTERVAL,
    resolution_date TIMESTAMP,
    
    -- IoT Data
    iot_device_sensors TEXT[],
    iot_data_points JSONB,
    iot_anomalies JSONB,
    iot_data_quality DECIMAL(3,2),
    iot_last_data_sync TIMESTAMP,
    iot_analysis_complete BOOLEAN DEFAULT false,
    
    -- Digital Forensics
    forensics_id UUID,
    forensic_analyst_id UUID,
    forensics_start_date TIMESTAMP,
    forensics_end_date TIMESTAMP,
    forensics_status VARCHAR(50),
    forensics_findings TEXT,
    forensics_evidence JSONB,
    forensics_tools TEXT[],
    forensics_methodology TEXT,
    
    -- Business Rules
    business_rules JSONB,
    exclusions JSONB,
    approvals JSONB,
    
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
    CONSTRAINT chk_claim_type CHECK (claim_type IN ('physical_damage', 'liquid_damage', 'theft', 'loss', 'liability', 'data_recovery', 'repair', 'replacement')),
    CONSTRAINT chk_claim_status CHECK (status IN ('filed', 'under_review', 'investigating', 'pending', 'approved', 'denied', 'paid', 'closed', 'reopened', 'appealed')),
    CONSTRAINT chk_claim_severity CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT chk_claim_priority CHECK (priority IN ('low', 'normal', 'high', 'urgent', 'emergency')),
    CONSTRAINT chk_claimed_amount CHECK (claimed_amount > 0),
    CONSTRAINT chk_incident_date CHECK (incident_date <= reported_date)
);

-- Analytics Events table
CREATE TABLE analytics_events (
    event_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(50) NOT NULL,
    source VARCHAR(100) NOT NULL,
    user_id UUID REFERENCES customers(customer_id),
    device_id UUID REFERENCES devices(device_id),
    session_id UUID,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payload JSONB NOT NULL,
    tags JSONB,
    custom_fields JSONB,
    
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
    encryption_enabled BOOLEAN DEFAULT false,
    encryption_type VARCHAR(50),
    data_retention_policy VARCHAR(100),
    access_logs JSONB,
    
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
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_analytics_event_type CHECK (event_type IN ('user_action', 'device_event', 'system_event', 'policy_event', 'claim_event', 'custom'))
);

-- =============================================================================
-- INDEXES FOR PERFORMANCE
-- =============================================================================

-- Customers indexes
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_customer_number ON customers(customer_number);
CREATE INDEX idx_customers_account_status ON customers(account_status);
CREATE INDEX idx_customers_kyc_status ON customers(kyc_status);
CREATE INDEX idx_customers_risk_level ON customers(risk_level);
CREATE INDEX idx_customers_created_at ON customers(created_at);
CREATE INDEX idx_customers_last_login_date ON customers(last_login_date);
CREATE INDEX idx_customers_country ON customers(country);
CREATE INDEX idx_customers_city ON customers(city);

-- Devices indexes
CREATE INDEX idx_devices_imei ON devices(imei);
CREATE INDEX idx_devices_serial_number ON devices(serial_number);
CREATE INDEX idx_devices_current_owner_id ON devices(current_owner_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_manufacturer ON devices(manufacturer);
CREATE INDEX idx_devices_model ON devices(model);
CREATE INDEX idx_devices_created_at ON devices(created_at);
CREATE INDEX idx_devices_purchase_date ON devices(purchase_date);
CREATE INDEX idx_devices_manufacture_date ON devices(manufacture_date);

-- Policies indexes
CREATE INDEX idx_policies_policy_number ON policies(policy_number);
CREATE INDEX idx_policies_customer_id ON policies(customer_id);
CREATE INDEX idx_policies_device_id ON policies(device_id);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_policy_type ON policies(policy_type);
CREATE INDEX idx_policies_coverage_type ON policies(coverage_type);
CREATE INDEX idx_policies_effective_date ON policies(effective_date);
CREATE INDEX idx_policies_expiration_date ON policies(expiration_date);
CREATE INDEX idx_policies_created_at ON policies(created_at);
CREATE INDEX idx_policies_next_payment_date ON policies(next_payment_date);

-- Claims indexes
CREATE INDEX idx_claims_claim_number ON claims(claim_number);
CREATE INDEX idx_claims_policy_id ON claims(policy_id);
CREATE INDEX idx_claims_customer_id ON claims(customer_id);
CREATE INDEX idx_claims_device_id ON claims(device_id);
CREATE INDEX idx_claims_status ON claims(status);
CREATE INDEX idx_claims_claim_type ON claims(claim_type);
CREATE INDEX idx_claims_incident_date ON claims(incident_date);
CREATE INDEX idx_claims_reported_date ON claims(reported_date);
CREATE INDEX idx_claims_created_at ON claims(created_at);
CREATE INDEX idx_claims_assigned_to ON claims(assigned_to);

-- Analytics Events indexes
CREATE INDEX idx_analytics_events_event_type ON analytics_events(event_type);
CREATE INDEX idx_analytics_events_source ON analytics_events(source);
CREATE INDEX idx_analytics_events_user_id ON analytics_events(user_id);
CREATE INDEX idx_analytics_events_device_id ON analytics_events(device_id);
CREATE INDEX idx_analytics_events_timestamp ON analytics_events(timestamp);
CREATE INDEX idx_analytics_events_created_at ON analytics_events(created_at);

-- Composite indexes for common queries
CREATE INDEX idx_policies_customer_status ON policies(customer_id, status);
CREATE INDEX idx_policies_device_status ON policies(device_id, status);
CREATE INDEX idx_claims_policy_status ON claims(policy_id, status);
CREATE INDEX idx_claims_customer_status ON claims(customer_id, status);
CREATE INDEX idx_devices_owner_status ON devices(current_owner_id, status);
CREATE INDEX idx_analytics_events_user_timestamp ON analytics_events(user_id, timestamp);
CREATE INDEX idx_analytics_events_device_timestamp ON analytics_events(device_id, timestamp);

-- JSONB indexes for better performance
CREATE INDEX idx_customers_personal_info_gin ON customers USING GIN ((personal_info::jsonb));
CREATE INDEX idx_customers_address_gin ON customers USING GIN ((address::jsonb));
CREATE INDEX idx_devices_specifications_gin ON devices USING GIN ((specifications::jsonb));
CREATE INDEX idx_policies_coverage_limits_gin ON policies USING GIN ((coverage_limits::jsonb));
CREATE INDEX idx_claims_evidence_gin ON claims USING GIN ((evidence::jsonb));
CREATE INDEX idx_analytics_events_payload_gin ON analytics_events USING GIN (payload);

-- =============================================================================
-- TRIGGERS FOR AUTOMATIC UPDATES
-- =============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_devices_updated_at BEFORE UPDATE ON devices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_policies_updated_at BEFORE UPDATE ON policies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_claims_updated_at BEFORE UPDATE ON claims FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_analytics_events_updated_at BEFORE UPDATE ON analytics_events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- VIEWS FOR COMMON QUERIES
-- =============================================================================

-- Active policies view
CREATE VIEW active_policies AS
SELECT 
    p.policy_id,
    p.policy_number,
    c.customer_number,
    c.first_name || ' ' || c.last_name as customer_name,
    d.imei,
    d.model as device_model,
    p.policy_type,
    p.coverage_type,
    p.final_premium,
    p.currency,
    p.effective_date,
    p.expiration_date,
    p.status
FROM policies p
JOIN customers c ON p.customer_id = c.customer_id
JOIN devices d ON p.device_id = d.device_id
WHERE p.status = 'active';

-- Claims summary view
CREATE VIEW claims_summary AS
SELECT 
    cl.claim_id,
    cl.claim_number,
    p.policy_number,
    c.customer_number,
    c.first_name || ' ' || c.last_name as customer_name,
    d.imei,
    cl.claim_type,
    cl.status,
    cl.claimed_amount,
    cl.approved_amount,
    cl.paid_amount,
    cl.currency,
    cl.incident_date,
    cl.reported_date,
    cl.created_at
FROM claims cl
JOIN policies p ON cl.policy_id = p.policy_id
JOIN customers c ON cl.customer_id = c.customer_id
JOIN devices d ON cl.device_id = d.device_id;

-- Customer devices view
CREATE VIEW customer_devices AS
SELECT 
    c.customer_id,
    c.customer_number,
    c.first_name || ' ' || c.last_name as customer_name,
    d.device_id,
    d.imei,
    d.serial_number,
    d.model,
    d.manufacturer,
    d.brand,
    d.status as device_status,
    d.current_value,
    d.currency,
    d.purchase_date,
    d.registration_date
FROM customers c
JOIN devices d ON c.customer_id = d.current_owner_id;

-- Risk assessment view
CREATE VIEW risk_assessment_summary AS
SELECT 
    c.customer_id,
    c.customer_number,
    c.first_name || ' ' || c.last_name as customer_name,
    c.risk_level,
    c.risk_score,
    c.fraud_risk_level,
    c.fraud_risk_score,
    c.kyc_status,
    c.account_status,
    COUNT(p.policy_id) as active_policies,
    COUNT(cl.claim_id) as total_claims,
    SUM(CASE WHEN cl.status = 'approved' THEN cl.approved_amount ELSE 0 END) as total_approved_claims
FROM customers c
LEFT JOIN policies p ON c.customer_id = p.customer_id AND p.status = 'active'
LEFT JOIN claims cl ON c.customer_id = cl.customer_id
GROUP BY c.customer_id, c.customer_number, c.first_name, c.last_name, c.risk_level, c.risk_score, c.fraud_risk_level, c.fraud_risk_score, c.kyc_status, c.account_status;

-- =============================================================================
-- FUNCTIONS FOR BUSINESS LOGIC
-- =============================================================================

-- Function to calculate device depreciation
CREATE OR REPLACE FUNCTION calculate_device_depreciation(
    original_value DECIMAL,
    manufacture_date DATE,
    depreciation_rate DECIMAL DEFAULT 0.15
)
RETURNS DECIMAL AS $$
DECLARE
    age_years DECIMAL;
    depreciation_factor DECIMAL;
BEGIN
    age_years := EXTRACT(YEAR FROM AGE(CURRENT_DATE, manufacture_date));
    depreciation_factor := GREATEST(0.1, 1.0 - (depreciation_rate * age_years));
    RETURN original_value * depreciation_factor;
END;
$$ LANGUAGE plpgsql;

-- Function to check if customer is eligible for new policy
CREATE OR REPLACE FUNCTION check_policy_eligibility(
    customer_uuid UUID,
    device_uuid UUID
)
RETURNS BOOLEAN AS $$
DECLARE
    customer_status VARCHAR(50);
    device_status VARCHAR(50);
    active_policies_count INTEGER;
BEGIN
    -- Check customer status
    SELECT account_status INTO customer_status 
    FROM customers 
    WHERE customer_id = customer_uuid;
    
    IF customer_status != 'active' THEN
        RETURN FALSE;
    END IF;
    
    -- Check device status
    SELECT status INTO device_status 
    FROM devices 
    WHERE device_id = device_uuid;
    
    IF device_status != 'active' THEN
        RETURN FALSE;
    END IF;
    
    -- Check if device already has active policy
    SELECT COUNT(*) INTO active_policies_count 
    FROM policies 
    WHERE device_id = device_uuid AND status = 'active';
    
    RETURN active_policies_count = 0;
END;
$$ LANGUAGE plpgsql;

-- Function to calculate claim risk score
CREATE OR REPLACE FUNCTION calculate_claim_risk_score(
    claim_amount DECIMAL,
    customer_risk_score DECIMAL,
    device_value DECIMAL,
    claim_history_count INTEGER DEFAULT 0
)
RETURNS DECIMAL AS $$
DECLARE
    amount_factor DECIMAL;
    history_factor DECIMAL;
BEGIN
    -- Amount factor (higher amount = higher risk)
    amount_factor := LEAST(1.0, claim_amount / device_value);
    
    -- History factor (more claims = higher risk)
    history_factor := LEAST(1.0, claim_history_count * 0.2);
    
    RETURN (customer_risk_score * 0.4 + amount_factor * 0.4 + history_factor * 0.2);
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- COMMENTS FOR DOCUMENTATION
-- =============================================================================

COMMENT ON TABLE customers IS 'Core customer entity with comprehensive personal, financial, and security information';
COMMENT ON TABLE devices IS 'Smartphone devices with specifications, condition tracking, and ownership history';
COMMENT ON TABLE policies IS 'Insurance policies linking customers to devices with coverage and premium details';
COMMENT ON TABLE claims IS 'Insurance claims with incident details, assessment, and processing information';
COMMENT ON TABLE analytics_events IS 'Analytics events for tracking user behavior, device events, and system activities';

COMMENT ON COLUMN customers.customer_id IS 'Unique identifier for the customer';
COMMENT ON COLUMN customers.customer_number IS 'Business-friendly customer number for external reference';
COMMENT ON COLUMN customers.email IS 'Primary email address for communication';
COMMENT ON COLUMN customers.phone IS 'Primary phone number for communication';
COMMENT ON COLUMN customers.risk_score IS 'Calculated risk score from 0.0 to 1.0';
COMMENT ON COLUMN customers.kyc_status IS 'Know Your Customer verification status';

COMMENT ON COLUMN devices.device_id IS 'Unique identifier for the device';
COMMENT ON COLUMN devices.imei IS 'International Mobile Equipment Identity (15 digits)';
COMMENT ON COLUMN devices.serial_number IS 'Manufacturer serial number';
COMMENT ON COLUMN devices.current_owner_id IS 'Reference to current owner in customers table';
COMMENT ON COLUMN devices.current_value IS 'Current market value of the device';

COMMENT ON COLUMN policies.policy_id IS 'Unique identifier for the policy';
COMMENT ON COLUMN policies.policy_number IS 'Business-friendly policy number for external reference';
COMMENT ON COLUMN policies.customer_id IS 'Reference to customer in customers table';
COMMENT ON COLUMN policies.device_id IS 'Reference to device in devices table';
COMMENT ON COLUMN policies.final_premium IS 'Final calculated premium amount';

COMMENT ON COLUMN claims.claim_id IS 'Unique identifier for the claim';
COMMENT ON COLUMN claims.claim_number IS 'Business-friendly claim number for external reference';
COMMENT ON COLUMN claims.policy_id IS 'Reference to policy in policies table';
COMMENT ON COLUMN claims.claimed_amount IS 'Amount claimed by the customer';
COMMENT ON COLUMN claims.approved_amount IS 'Amount approved by the insurance company';

COMMENT ON COLUMN analytics_events.event_id IS 'Unique identifier for the analytics event';
COMMENT ON COLUMN analytics_events.event_type IS 'Type of analytics event';
COMMENT ON COLUMN analytics_events.source IS 'Source system generating the event';
COMMENT ON COLUMN analytics_events.payload IS 'JSON payload containing event data'; 