-- SmartSure PostgreSQL Database Schema - Additional Business Tables
-- Enterprise-grade smartphone insurance blockchain system
-- Supporting remaining business options from the documentation

-- =============================================================================
-- ADDITIONAL BUSINESS TABLES
-- =============================================================================

-- Repair Network table (Business Option 6)
CREATE TABLE repair_network (
    repair_network_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    network_name VARCHAR(200) NOT NULL,
    network_type VARCHAR(50), -- certified, authorized, independent
    
    -- Contact Information
    contact_person VARCHAR(200),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    website VARCHAR(255),
    
    -- Location Information
    address JSONB,
    service_areas JSONB,
    operating_hours JSONB,
    
    -- Certification Information
    certification_status VARCHAR(50), -- certified, pending, expired, suspended
    certification_date TIMESTAMP,
    certification_expiry TIMESTAMP,
    certification_body VARCHAR(200),
    certification_number VARCHAR(100),
    
    -- Service Capabilities
    service_types TEXT[],
    device_brands TEXT[],
    device_models TEXT[],
    warranty_services BOOLEAN DEFAULT false,
    insurance_services BOOLEAN DEFAULT false,
    
    -- Quality Metrics
    quality_rating DECIMAL(3,2),
    customer_satisfaction_score DECIMAL(3,2),
    completion_rate DECIMAL(5,4),
    average_repair_time INTERVAL,
    
    -- Business Information
    business_license VARCHAR(100),
    insurance_coverage BOOLEAN DEFAULT false,
    insurance_provider VARCHAR(200),
    insurance_policy_number VARCHAR(100),
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_repair_network_type CHECK (network_type IN ('certified', 'authorized', 'independent', 'premium')),
    CONSTRAINT chk_certification_status CHECK (certification_status IN ('certified', 'pending', 'expired', 'suspended', 'revoked')),
    CONSTRAINT chk_quality_rating CHECK (quality_rating >= 0 AND quality_rating <= 1),
    CONSTRAINT chk_customer_satisfaction_score CHECK (customer_satisfaction_score >= 0 AND customer_satisfaction_score <= 1)
);

-- Warranty Extensions table (Business Option 7)
CREATE TABLE warranty_extensions (
    warranty_extension_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    original_warranty_id UUID,
    
    -- Extension Details
    extension_type VARCHAR(50), -- manufacturer, third_party, insurance_based
    extension_period INTERVAL,
    extension_start_date TIMESTAMP,
    extension_end_date TIMESTAMP,
    
    -- Coverage Details
    coverage_scope TEXT[],
    coverage_limits JSONB,
    coverage_exclusions TEXT[],
    additional_benefits JSONB,
    
    -- Cost Information
    extension_cost DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(50),
    payment_status VARCHAR(50),
    
    -- Provider Information
    warranty_provider VARCHAR(200),
    provider_contact JSONB,
    provider_terms TEXT,
    
    -- Status
    status VARCHAR(50), -- active, expired, cancelled, transferred
    activation_date TIMESTAMP,
    cancellation_date TIMESTAMP,
    cancellation_reason TEXT,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_warranty_extension_type CHECK (extension_type IN ('manufacturer', 'third_party', 'insurance_based', 'retailer')),
    CONSTRAINT chk_warranty_status CHECK (status IN ('active', 'expired', 'cancelled', 'transferred', 'suspended')),
    CONSTRAINT chk_payment_status CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded'))
);

-- Trade-in Offers table (Business Option 8)
CREATE TABLE trade_in_offers (
    trade_in_offer_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    
    -- Offer Details
    offer_type VARCHAR(50), -- trade_in, buyback, upgrade
    offer_amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    offer_currency VARCHAR(3) DEFAULT 'USD',
    
    -- Device Assessment
    device_condition VARCHAR(50),
    device_grade VARCHAR(5),
    device_grade_score DECIMAL(3,2),
    assessment_notes TEXT,
    assessment_photos TEXT[],
    
    -- Offer Terms
    offer_valid_until TIMESTAMP,
    acceptance_deadline TIMESTAMP,
    terms_and_conditions TEXT,
    special_conditions JSONB,
    
    -- Trade-in Process
    trade_in_status VARCHAR(50), -- offered, accepted, declined, expired, completed
    acceptance_date TIMESTAMP,
    completion_date TIMESTAMP,
    trade_in_location JSONB,
    
    -- Payment Information
    payment_method VARCHAR(50),
    payment_status VARCHAR(50),
    payment_date TIMESTAMP,
    payment_transaction_id VARCHAR(255),
    
    -- New Device Information
    new_device_id UUID REFERENCES devices(device_id),
    new_device_discount DECIMAL(15,2),
    upgrade_incentive DECIMAL(15,2),
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_trade_in_offer_type CHECK (offer_type IN ('trade_in', 'buyback', 'upgrade', 'exchange')),
    CONSTRAINT chk_trade_in_status CHECK (trade_in_status IN ('offered', 'accepted', 'declined', 'expired', 'completed', 'cancelled')),
    CONSTRAINT chk_device_grade_score CHECK (device_grade_score >= 0 AND device_grade_score <= 1),
    CONSTRAINT chk_offer_amount CHECK (offer_amount > 0)
);

-- Device Grading table (Business Option 9)
CREATE TABLE device_grading (
    grading_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    
    -- Grading Details
    grading_type VARCHAR(50), -- physical, functional, cosmetic, comprehensive
    overall_grade VARCHAR(5), -- A, B, C, D, F
    overall_score DECIMAL(3,2),
    
    -- Component Grades
    physical_grade VARCHAR(5),
    physical_score DECIMAL(3,2),
    functional_grade VARCHAR(5),
    functional_score DECIMAL(3,2),
    cosmetic_grade VARCHAR(5),
    cosmetic_score DECIMAL(3,2),
    battery_grade VARCHAR(5),
    battery_score DECIMAL(3,2),
    screen_grade VARCHAR(5),
    screen_score DECIMAL(3,2),
    
    -- Assessment Details
    assessment_criteria JSONB,
    assessment_notes TEXT,
    assessment_photos TEXT[],
    assessment_videos TEXT[],
    
    -- Grading Factors
    age_factor DECIMAL(3,2),
    usage_factor DECIMAL(3,2),
    damage_factor DECIMAL(3,2),
    repair_factor DECIMAL(3,2),
    
    -- Certification
    certified BOOLEAN DEFAULT false,
    certification_date TIMESTAMP,
    certification_valid_until TIMESTAMP,
    certification_body VARCHAR(200),
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    grading_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_grading_type CHECK (grading_type IN ('physical', 'functional', 'cosmetic', 'comprehensive', 'specialized')),
    CONSTRAINT chk_overall_grade CHECK (overall_grade IN ('A+', 'A', 'A-', 'B+', 'B', 'B-', 'C+', 'C', 'C-', 'D+', 'D', 'D-', 'F')),
    CONSTRAINT chk_overall_score CHECK (overall_score >= 0 AND overall_score <= 1)
);

-- Claim Automation table (Business Option 10)
CREATE TABLE claim_automation (
    automation_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(claim_id),
    
    -- Automation Details
    automation_type VARCHAR(50), -- document_processing, fraud_detection, assessment, approval
    automation_status VARCHAR(50), -- pending, processing, completed, failed, manual_review
    automation_score DECIMAL(5,4),
    
    -- Processing Steps
    processing_steps JSONB,
    current_step VARCHAR(100),
    step_progress DECIMAL(5,2),
    estimated_completion TIMESTAMP,
    
    -- AI/ML Models
    model_version VARCHAR(50),
    model_confidence DECIMAL(5,4),
    model_accuracy DECIMAL(5,4),
    model_parameters JSONB,
    
    -- Decision Information
    automated_decision VARCHAR(50), -- approve, deny, refer, pending
    decision_confidence DECIMAL(5,4),
    decision_reasoning TEXT,
    decision_factors JSONB,
    
    -- Manual Review
    manual_review_required BOOLEAN DEFAULT false,
    review_reason TEXT,
    reviewed_by UUID,
    review_date TIMESTAMP,
    review_notes TEXT,
    
    -- Performance Metrics
    processing_time INTERVAL,
    accuracy_score DECIMAL(5,4),
    efficiency_score DECIMAL(5,4),
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_automation_type CHECK (automation_type IN ('document_processing', 'fraud_detection', 'assessment', 'approval', 'payout')),
    CONSTRAINT chk_automation_status CHECK (automation_status IN ('pending', 'processing', 'completed', 'failed', 'manual_review', 'cancelled')),
    CONSTRAINT chk_automated_decision CHECK (automated_decision IN ('approve', 'deny', 'refer', 'pending', 'escalate'))
);

-- Blockchain Verification table (Business Option 11)
CREATE TABLE blockchain_verifications (
    verification_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(claim_id),
    
    -- Verification Details
    verification_type VARCHAR(50), -- claim_verification, policy_verification, device_verification
    verification_status VARCHAR(50), -- pending, verified, failed, disputed
    verification_score DECIMAL(5,4),
    
    -- Blockchain Information
    transaction_hash VARCHAR(255),
    block_number BIGINT,
    block_timestamp TIMESTAMP,
    chain_id VARCHAR(100),
    network_id VARCHAR(100),
    smart_contract VARCHAR(200),
    gas_used BIGINT,
    gas_price BIGINT,
    confirmations INTEGER,
    
    -- Verification Results
    verification_result JSONB,
    verification_evidence JSONB,
    verification_notes TEXT,
    
    -- Timestamps
    verification_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_verification_type CHECK (verification_type IN ('claim_verification', 'policy_verification', 'device_verification', 'customer_verification')),
    CONSTRAINT chk_verification_status CHECK (verification_status IN ('pending', 'verified', 'failed', 'disputed', 'expired'))
);

-- Fraud Detection table (Business Option 12)
CREATE TABLE fraud_detection (
    fraud_detection_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(claim_id),
    
    -- Detection Details
    detection_type VARCHAR(50), -- claim_fraud, policy_fraud, identity_fraud
    detection_status VARCHAR(50), -- pending, detected, cleared, confirmed
    fraud_score DECIMAL(5,4),
    fraud_level VARCHAR(20), -- low, medium, high, critical
    
    -- AI/ML Models
    model_version VARCHAR(50),
    model_confidence DECIMAL(5,4),
    model_accuracy DECIMAL(5,4),
    detection_algorithm VARCHAR(100),
    
    -- Fraud Indicators
    fraud_indicators JSONB,
    risk_factors JSONB,
    suspicious_patterns JSONB,
    anomaly_scores JSONB,
    
    -- Investigation
    investigation_required BOOLEAN DEFAULT false,
    investigation_status VARCHAR(50),
    investigation_notes TEXT,
    investigation_findings JSONB,
    
    -- Actions
    recommended_action VARCHAR(100),
    action_taken VARCHAR(100),
    action_date TIMESTAMP,
    action_by UUID,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    detection_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_fraud_detection_type CHECK (detection_type IN ('claim_fraud', 'policy_fraud', 'identity_fraud', 'device_fraud')),
    CONSTRAINT chk_fraud_detection_status CHECK (detection_status IN ('pending', 'detected', 'cleared', 'confirmed', 'escalated')),
    CONSTRAINT chk_fraud_level CHECK (fraud_level IN ('low', 'medium', 'high', 'critical'))
);

-- IoT Damage Detection table (Business Option 13)
CREATE TABLE iot_damage_detection (
    iot_detection_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    iot_device_id UUID REFERENCES iot_devices(iot_device_id),
    
    -- Detection Details
    sensor_type VARCHAR(100), -- accelerometer, gyroscope, moisture, temperature
    detection_type VARCHAR(50), -- impact, water_damage, temperature, vibration
    detection_status VARCHAR(50), -- detected, confirmed, false_alarm, resolved
    
    -- Sensor Data
    sensor_value DECIMAL(15,4),
    threshold_value DECIMAL(15,4),
    unit_of_measurement VARCHAR(20),
    data_quality DECIMAL(3,2),
    
    -- Damage Assessment
    damage_severity VARCHAR(20), -- low, medium, high, critical
    damage_probability DECIMAL(5,4),
    damage_estimate DECIMAL(15,2),
    recommended_action VARCHAR(200),
    
    -- Alert Information
    alert_generated BOOLEAN DEFAULT false,
    alert_sent_to UUID,
    alert_sent_at TIMESTAMP,
    alert_response VARCHAR(50),
    
    -- Investigation
    investigation_required BOOLEAN DEFAULT false,
    investigation_status VARCHAR(50),
    investigation_notes TEXT,
    investigation_findings JSONB,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    detection_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_iot_detection_type CHECK (detection_type IN ('impact', 'water_damage', 'temperature', 'vibration', 'pressure', 'humidity')),
    CONSTRAINT chk_iot_detection_status CHECK (detection_status IN ('detected', 'confirmed', 'false_alarm', 'resolved', 'investigating')),
    CONSTRAINT chk_iot_damage_severity CHECK (damage_severity IN ('low', 'medium', 'high', 'critical'))
);

-- White Label API Services table (Business Option 14)
CREATE TABLE white_label_apis (
    api_service_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_name VARCHAR(200) NOT NULL,
    service_type VARCHAR(50), -- insurance_api, claims_api, device_api
    
    -- API Configuration
    api_key VARCHAR(255) UNIQUE NOT NULL,
    api_secret VARCHAR(255),
    api_version VARCHAR(20),
    base_url VARCHAR(500),
    
    -- Permissions and Access
    permissions JSONB,
    rate_limits JSONB,
    access_controls JSONB,
    allowed_ips TEXT[],
    
    -- Branding
    white_label_branding JSONB,
    custom_domain VARCHAR(255),
    logo_url VARCHAR(500),
    color_scheme JSONB,
    
    -- Usage Tracking
    total_requests BIGINT DEFAULT 0,
    successful_requests BIGINT DEFAULT 0,
    failed_requests BIGINT DEFAULT 0,
    last_request_at TIMESTAMP,
    
    -- Billing
    billing_model VARCHAR(50), -- per_request, monthly, usage_based
    billing_rate DECIMAL(10,4),
    billing_currency VARCHAR(3) DEFAULT 'USD',
    billing_period VARCHAR(20),
    
    -- Status
    status VARCHAR(50), -- active, inactive, suspended, expired
    activation_date TIMESTAMP,
    expiry_date TIMESTAMP,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_white_label_service_type CHECK (service_type IN ('insurance_api', 'claims_api', 'device_api', 'analytics_api', 'full_api')),
    CONSTRAINT chk_white_label_billing_model CHECK (billing_model IN ('per_request', 'monthly', 'usage_based', 'tiered')),
    CONSTRAINT chk_white_label_status CHECK (status IN ('active', 'inactive', 'suspended', 'expired', 'pending'))
);

-- In-App Insurance Purchases table (Business Option 15)
CREATE TABLE in_app_insurance_purchases (
    purchase_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    retailer_id UUID NOT NULL REFERENCES customers(customer_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Purchase Details
    insurance_type VARCHAR(50), -- comprehensive, accident, theft, basic
    purchase_type VARCHAR(50), -- immediate, deferred, trial
    purchase_amount DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- App Information
    app_name VARCHAR(200),
    app_version VARCHAR(50),
    platform VARCHAR(20), -- ios, android, web
    device_identifier VARCHAR(255),
    
    -- Purchase Flow
    purchase_flow JSONB,
    user_journey JSONB,
    conversion_funnel JSONB,
    abandonment_reason TEXT,
    
    -- Payment Information
    payment_method VARCHAR(50),
    payment_status VARCHAR(50),
    payment_transaction_id VARCHAR(255),
    payment_processor VARCHAR(100),
    
    -- Commission
    retailer_commission_rate DECIMAL(5,4),
    retailer_commission_amount DECIMAL(15,2),
    commission_paid BOOLEAN DEFAULT false,
    commission_payment_date TIMESTAMP,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_in_app_insurance_type CHECK (insurance_type IN ('comprehensive', 'accident', 'theft', 'basic', 'premium')),
    CONSTRAINT chk_in_app_purchase_type CHECK (purchase_type IN ('immediate', 'deferred', 'trial', 'subscription')),
    CONSTRAINT chk_in_app_platform CHECK (platform IN ('ios', 'android', 'web', 'mobile_web')),
    CONSTRAINT chk_in_app_payment_status CHECK (payment_status IN ('pending', 'completed', 'failed', 'refunded'))
);

-- Telco Bundled Insurance table (Business Option 16)
CREATE TABLE telco_bundled_insurance (
    telco_bundle_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telco_id UUID NOT NULL REFERENCES customers(customer_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    
    -- Bundle Details
    bundle_type VARCHAR(50), -- data_insurance, voice_insurance, complete_bundle
    bundle_name VARCHAR(200),
    bundle_description TEXT,
    
    -- Telco Services
    telco_services JSONB,
    data_plan VARCHAR(100),
    voice_plan VARCHAR(100),
    text_plan VARCHAR(100),
    
    -- Insurance Coverage
    insurance_coverage JSONB,
    coverage_limits JSONB,
    coverage_exclusions TEXT[],
    additional_benefits JSONB,
    
    -- Pricing
    bundle_price DECIMAL(15,2),
    insurance_portion DECIMAL(15,2),
    telco_portion DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- Revenue Sharing
    revenue_split_percentage DECIMAL(5,4),
    telco_commission_rate DECIMAL(5,4),
    telco_commission_amount DECIMAL(15,2),
    
    -- Contract Details
    contract_duration INTERVAL,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    auto_renewal BOOLEAN DEFAULT true,
    
    -- Status
    status VARCHAR(50), -- active, suspended, cancelled, expired
    activation_date TIMESTAMP,
    cancellation_date TIMESTAMP,
    cancellation_reason TEXT,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_telco_bundle_type CHECK (bundle_type IN ('data_insurance', 'voice_insurance', 'complete_bundle', 'premium_bundle')),
    CONSTRAINT chk_telco_bundle_status CHECK (status IN ('active', 'suspended', 'cancelled', 'expired', 'pending')),
    CONSTRAINT chk_telco_bundle_price CHECK (bundle_price > 0)
);

-- Device Theft Reporting table (Business Option 17)
CREATE TABLE device_theft_reporting (
    theft_report_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(device_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    
    -- Theft Details
    theft_date TIMESTAMP NOT NULL,
    theft_location JSONB,
    theft_description TEXT,
    theft_type VARCHAR(50), -- stolen, lost, robbery, burglary
    
    -- Reporting Information
    report_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    report_method VARCHAR(50), -- app, web, phone, in_person
    report_source VARCHAR(100),
    report_number VARCHAR(100),
    
    -- Law Enforcement
    police_report_number VARCHAR(100),
    police_department VARCHAR(200),
    police_contact JSONB,
    law_enforcement_status VARCHAR(50),
    
    -- Investigation
    investigation_status VARCHAR(50), -- pending, investigating, closed, reopened
    investigation_notes TEXT,
    investigation_findings JSONB,
    investigation_evidence JSONB,
    
    -- Recovery
    recovery_status VARCHAR(50), -- not_recovered, recovered, partially_recovered
    recovery_date TIMESTAMP,
    recovery_location JSONB,
    recovery_method VARCHAR(100),
    recovery_notes TEXT,
    
    -- Insurance Impact
    insurance_claim_filed BOOLEAN DEFAULT false,
    claim_id UUID REFERENCES claims(claim_id),
    claim_status VARCHAR(50),
    claim_amount DECIMAL(15,2),
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_device_theft_type CHECK (theft_type IN ('stolen', 'lost', 'robbery', 'burglary', 'unknown')),
    CONSTRAINT chk_device_theft_report_method CHECK (report_method IN ('app', 'web', 'phone', 'in_person', 'email')),
    CONSTRAINT chk_device_theft_investigation_status CHECK (investigation_status IN ('pending', 'investigating', 'closed', 'reopened', 'suspended')),
    CONSTRAINT chk_device_theft_recovery_status CHECK (recovery_status IN ('not_recovered', 'recovered', 'partially_recovered', 'in_progress'))
);

-- Insurance Analytics table (Business Option 18)
CREATE TABLE insurance_analytics (
    analytics_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Analytics Details
    analytics_type VARCHAR(50), -- business_intelligence, risk_analytics, performance_metrics
    analytics_period VARCHAR(20), -- daily, weekly, monthly, quarterly, yearly
    analytics_date DATE,
    
    -- Key Metrics
    total_policies INTEGER,
    active_policies INTEGER,
    total_claims INTEGER,
    active_claims INTEGER,
    fraud_rate DECIMAL(5,4),
    average_premium DECIMAL(15,2),
    customer_retention_rate DECIMAL(5,4),
    
    -- Financial Metrics
    total_premium_revenue DECIMAL(15,2),
    total_claim_payouts DECIMAL(15,2),
    net_profit DECIMAL(15,2),
    loss_ratio DECIMAL(5,4),
    combined_ratio DECIMAL(5,4),
    
    -- Risk Metrics
    average_risk_score DECIMAL(5,4),
    high_risk_policies INTEGER,
    risk_distribution JSONB,
    risk_trends JSONB,
    
    -- Customer Metrics
    total_customers INTEGER,
    new_customers INTEGER,
    churned_customers INTEGER,
    customer_satisfaction_score DECIMAL(3,2),
    customer_lifetime_value DECIMAL(15,2),
    
    -- Device Metrics
    total_devices INTEGER,
    device_types_distribution JSONB,
    device_conditions_distribution JSONB,
    average_device_value DECIMAL(15,2),
    
    -- Performance Metrics
    claim_processing_time INTERVAL,
    average_settlement_time INTERVAL,
    customer_service_metrics JSONB,
    operational_efficiency DECIMAL(5,4),
    
    -- Blockchain Metrics
    blockchain_transactions INTEGER,
    blockchain_verification_rate DECIMAL(5,4),
    blockchain_performance_metrics JSONB,
    
    -- Timestamps
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_insurance_analytics_type CHECK (analytics_type IN ('business_intelligence', 'risk_analytics', 'performance_metrics', 'fraud_analytics')),
    CONSTRAINT chk_insurance_analytics_period CHECK (analytics_period IN ('daily', 'weekly', 'monthly', 'quarterly', 'yearly')),
    CONSTRAINT chk_insurance_analytics_fraud_rate CHECK (fraud_rate >= 0 AND fraud_rate <= 1),
    CONSTRAINT chk_insurance_analytics_customer_retention CHECK (customer_retention_rate >= 0 AND customer_retention_rate <= 1)
);

-- Insurance Cooperatives table (Business Option 19)
CREATE TABLE insurance_cooperatives (
    cooperative_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cooperative_name VARCHAR(200) NOT NULL,
    cooperative_type VARCHAR(50), -- peer_to_peer, community, professional
    
    -- Cooperative Details
    description TEXT,
    mission_statement TEXT,
    governance_model VARCHAR(50), -- democratic, representative, hybrid
    membership_criteria JSONB,
    
    -- Membership
    member_count INTEGER DEFAULT 0,
    max_members INTEGER,
    membership_fee DECIMAL(15,2),
    membership_currency VARCHAR(3) DEFAULT 'USD',
    
    -- Risk Pooling
    risk_pool_size DECIMAL(15,2),
    risk_distribution JSONB,
    risk_sharing_model VARCHAR(50), -- equal, proportional, tiered
    risk_assessment_method VARCHAR(100),
    
    -- Financial Information
    total_assets DECIMAL(15,2),
    total_liabilities DECIMAL(15,2),
    net_worth DECIMAL(15,2),
    reserve_ratio DECIMAL(5,4),
    
    -- Governance
    board_members JSONB,
    voting_rights JSONB,
    decision_making_process TEXT,
    meeting_schedule JSONB,
    
    -- Status
    status VARCHAR(50), -- active, inactive, suspended, dissolved
    formation_date TIMESTAMP,
    dissolution_date TIMESTAMP,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_insurance_cooperative_type CHECK (cooperative_type IN ('peer_to_peer', 'community', 'professional', 'industry')),
    CONSTRAINT chk_insurance_cooperative_governance CHECK (governance_model IN ('democratic', 'representative', 'hybrid', 'consensus')),
    CONSTRAINT chk_insurance_cooperative_status CHECK (status IN ('active', 'inactive', 'suspended', 'dissolved', 'forming'))
);

-- Cooperative Members table
CREATE TABLE cooperative_members (
    member_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cooperative_id UUID NOT NULL REFERENCES insurance_cooperatives(cooperative_id),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    
    -- Membership Details
    membership_type VARCHAR(50), -- founding, regular, premium, honorary
    membership_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    membership_status VARCHAR(50), -- active, inactive, suspended, terminated
    
    -- Voting Rights
    voting_power INTEGER DEFAULT 1,
    voting_rights_active BOOLEAN DEFAULT true,
    voting_history JSONB,
    
    -- Financial Participation
    contribution_amount DECIMAL(15,2),
    contribution_frequency VARCHAR(20), -- monthly, quarterly, annually
    total_contributions DECIMAL(15,2),
    last_contribution_date TIMESTAMP,
    
    -- Risk Sharing
    risk_share_percentage DECIMAL(5,4),
    risk_liability DECIMAL(15,2),
    claims_history JSONB,
    
    -- Governance Participation
    board_member BOOLEAN DEFAULT false,
    committee_member BOOLEAN DEFAULT false,
    participation_score DECIMAL(3,2),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_cooperative_membership_type CHECK (membership_type IN ('founding', 'regular', 'premium', 'honorary')),
    CONSTRAINT chk_cooperative_membership_status CHECK (membership_status IN ('active', 'inactive', 'suspended', 'terminated')),
    CONSTRAINT chk_cooperative_contribution_frequency CHECK (contribution_frequency IN ('monthly', 'quarterly', 'annually', 'one_time'))
);

-- Decentralized Marketplaces table (Business Option 20)
CREATE TABLE decentralized_marketplaces (
    marketplace_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    marketplace_name VARCHAR(200) NOT NULL,
    marketplace_type VARCHAR(50), -- insurance_marketplace, device_marketplace, service_marketplace
    
    -- Marketplace Details
    description TEXT,
    marketplace_url VARCHAR(500),
    platform_type VARCHAR(50), -- web, mobile, blockchain_native
    
    -- Governance
    governance_model VARCHAR(50), -- decentralized, hybrid, community_governed
    voting_mechanism VARCHAR(100),
    consensus_algorithm VARCHAR(100),
    
    -- Token Economics
    native_token VARCHAR(100),
    token_supply BIGINT,
    token_distribution JSONB,
    staking_requirements JSONB,
    
    -- Trading Features
    trading_pairs JSONB,
    liquidity_pools JSONB,
    trading_fees DECIMAL(5,4),
    listing_fees DECIMAL(15,2),
    
    -- Smart Contracts
    smart_contract_address VARCHAR(255),
    contract_version VARCHAR(20),
    contract_audit_status VARCHAR(50),
    contract_security_score DECIMAL(3,2),
    
    -- Performance Metrics
    total_volume DECIMAL(20,2),
    total_transactions INTEGER,
    active_users INTEGER,
    average_transaction_value DECIMAL(15,2),
    
    -- Status
    status VARCHAR(50), -- active, inactive, maintenance, deprecated
    launch_date TIMESTAMP,
    last_updated TIMESTAMP,
    
    -- Blockchain Information
    blockchain_transaction_hash VARCHAR(255),
    blockchain_block_number BIGINT,
    blockchain_block_timestamp TIMESTAMP,
    blockchain_status VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_decentralized_marketplace_type CHECK (marketplace_type IN ('insurance_marketplace', 'device_marketplace', 'service_marketplace', 'comprehensive')),
    CONSTRAINT chk_decentralized_platform_type CHECK (platform_type IN ('web', 'mobile', 'blockchain_native', 'hybrid')),
    CONSTRAINT chk_decentralized_governance CHECK (governance_model IN ('decentralized', 'hybrid', 'community_governed', 'dao')),
    CONSTRAINT chk_decentralized_status CHECK (status IN ('active', 'inactive', 'maintenance', 'deprecated', 'beta'))
); 