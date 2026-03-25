-- SmartSure PostgreSQL Database Schema - Policies and Claims Tables
-- Enterprise-grade smartphone insurance blockchain system

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