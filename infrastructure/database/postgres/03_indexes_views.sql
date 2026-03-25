-- SmartSure PostgreSQL Database Schema - Indexes, Views, and Functions
-- Enterprise-grade smartphone insurance blockchain system

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
CREATE INDEX idx_customers_date_of_birth ON customers(date_of_birth);
CREATE INDEX idx_customers_annual_income ON customers(annual_income);

-- Devices indexes
CREATE INDEX idx_devices_imei ON devices(imei);
CREATE INDEX idx_devices_serial_number ON devices(serial_number);
CREATE INDEX idx_devices_current_owner_id ON devices(current_owner_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_manufacturer ON devices(manufacturer);
CREATE INDEX idx_devices_model ON devices(model);
CREATE INDEX idx_devices_brand ON devices(brand);
CREATE INDEX idx_devices_created_at ON devices(created_at);
CREATE INDEX idx_devices_purchase_date ON devices(purchase_date);
CREATE INDEX idx_devices_manufacture_date ON devices(manufacture_date);
CREATE INDEX idx_devices_current_value ON devices(current_value);
CREATE INDEX idx_devices_original_value ON devices(original_value);

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
CREATE INDEX idx_policies_final_premium ON policies(final_premium);
CREATE INDEX idx_policies_risk_score ON policies(risk_score);

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
CREATE INDEX idx_claims_claimed_amount ON claims(claimed_amount);
CREATE INDEX idx_claims_approved_amount ON claims(approved_amount);
CREATE INDEX idx_claims_severity ON claims(severity);
CREATE INDEX idx_claims_priority ON claims(priority);

-- Analytics Events indexes
CREATE INDEX idx_analytics_events_event_type ON analytics_events(event_type);
CREATE INDEX idx_analytics_events_source ON analytics_events(source);
CREATE INDEX idx_analytics_events_user_id ON analytics_events(user_id);
CREATE INDEX idx_analytics_events_device_id ON analytics_events(device_id);
CREATE INDEX idx_analytics_events_timestamp ON analytics_events(timestamp);
CREATE INDEX idx_analytics_events_created_at ON analytics_events(created_at);
CREATE INDEX idx_analytics_events_session_id ON analytics_events(session_id);

-- Composite indexes for common queries
CREATE INDEX idx_policies_customer_status ON policies(customer_id, status);
CREATE INDEX idx_policies_device_status ON policies(device_id, status);
CREATE INDEX idx_claims_policy_status ON claims(policy_id, status);
CREATE INDEX idx_claims_customer_status ON claims(customer_id, status);
CREATE INDEX idx_claims_device_status ON claims(device_id, status);
CREATE INDEX idx_devices_owner_status ON devices(current_owner_id, status);
CREATE INDEX idx_analytics_events_user_timestamp ON analytics_events(user_id, timestamp);
CREATE INDEX idx_analytics_events_device_timestamp ON analytics_events(device_id, timestamp);
CREATE INDEX idx_analytics_events_type_timestamp ON analytics_events(event_type, timestamp);

-- Partial indexes for active records
CREATE INDEX idx_policies_active ON policies(customer_id, device_id) WHERE status = 'active';
CREATE INDEX idx_claims_open ON claims(policy_id, status) WHERE status IN ('filed', 'under_review', 'investigating', 'pending');
CREATE INDEX idx_devices_active ON devices(current_owner_id, status) WHERE status = 'active';
CREATE INDEX idx_customers_active ON customers(customer_id, account_status) WHERE account_status = 'active';

-- JSONB indexes for better performance
CREATE INDEX idx_customers_personal_info_gin ON customers USING GIN ((personal_info::jsonb));
CREATE INDEX idx_customers_address_gin ON customers USING GIN ((address::jsonb));
CREATE INDEX idx_devices_specifications_gin ON devices USING GIN ((specifications::jsonb));
CREATE INDEX idx_policies_coverage_limits_gin ON policies USING GIN ((coverage_limits::jsonb));
CREATE INDEX idx_claims_evidence_gin ON claims USING GIN ((evidence::jsonb));
CREATE INDEX idx_analytics_events_payload_gin ON analytics_events USING GIN (payload);

-- Text search indexes
CREATE INDEX idx_customers_name_search ON customers USING gin(to_tsvector('english', first_name || ' ' || last_name));
CREATE INDEX idx_devices_model_search ON devices USING gin(to_tsvector('english', model || ' ' || manufacturer));
CREATE INDEX idx_claims_description_search ON claims USING gin(to_tsvector('english', incident_description));

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

-- Function to update device current value based on depreciation
CREATE OR REPLACE FUNCTION update_device_current_value()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.original_value IS NOT NULL AND NEW.manufacture_date IS NOT NULL THEN
        NEW.current_value = calculate_device_depreciation(
            NEW.original_value, 
            NEW.manufacture_date, 
            COALESCE(NEW.depreciation_rate, 0.15)
        );
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically calculate device current value
CREATE TRIGGER update_device_value BEFORE INSERT OR UPDATE ON devices FOR EACH ROW EXECUTE FUNCTION update_device_current_value();

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
    c.email as customer_email,
    c.phone as customer_phone,
    d.imei,
    d.model as device_model,
    d.manufacturer as device_manufacturer,
    d.brand as device_brand,
    p.policy_type,
    p.coverage_type,
    p.final_premium,
    p.currency,
    p.effective_date,
    p.expiration_date,
    p.status,
    p.next_payment_date,
    p.risk_score,
    p.fraud_risk_score
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
    c.email as customer_email,
    d.imei,
    d.model as device_model,
    cl.claim_type,
    cl.status,
    cl.severity,
    cl.priority,
    cl.claimed_amount,
    cl.approved_amount,
    cl.paid_amount,
    cl.currency,
    cl.incident_date,
    cl.reported_date,
    cl.created_at,
    cl.assigned_to,
    cl.processing_time,
    cl.resolution_date
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
    c.email as customer_email,
    c.phone as customer_phone,
    d.device_id,
    d.imei,
    d.serial_number,
    d.model,
    d.manufacturer,
    d.brand,
    d.status as device_status,
    d.original_value,
    d.current_value,
    d.currency,
    d.purchase_date,
    d.registration_date,
    d.overall_condition,
    d.battery_health,
    d.grade,
    d.grade_score
FROM customers c
JOIN devices d ON c.customer_id = d.current_owner_id;

-- Risk assessment view
CREATE VIEW risk_assessment_summary AS
SELECT 
    c.customer_id,
    c.customer_number,
    c.first_name || ' ' || c.last_name as customer_name,
    c.email as customer_email,
    c.risk_level,
    c.risk_score,
    c.fraud_risk_level,
    c.fraud_risk_score,
    c.kyc_status,
    c.account_status,
    c.credit_score,
    c.annual_income,
    COUNT(p.policy_id) as active_policies,
    COUNT(cl.claim_id) as total_claims,
    COUNT(CASE WHEN cl.status = 'approved' THEN 1 END) as approved_claims,
    COUNT(CASE WHEN cl.status = 'denied' THEN 1 END) as denied_claims,
    SUM(CASE WHEN cl.status = 'approved' THEN cl.approved_amount ELSE 0 END) as total_approved_amount,
    AVG(CASE WHEN cl.status = 'approved' THEN cl.approved_amount END) as avg_approved_amount
FROM customers c
LEFT JOIN policies p ON c.customer_id = p.customer_id AND p.status = 'active'
LEFT JOIN claims cl ON c.customer_id = cl.customer_id
GROUP BY c.customer_id, c.customer_number, c.first_name, c.last_name, c.email, 
         c.risk_level, c.risk_score, c.fraud_risk_level, c.fraud_risk_score, 
         c.kyc_status, c.account_status, c.credit_score, c.annual_income;

-- Device analytics view
CREATE VIEW device_analytics AS
SELECT 
    d.device_id,
    d.imei,
    d.model,
    d.manufacturer,
    d.brand,
    d.current_owner_id,
    c.first_name || ' ' || c.last_name as current_owner,
    d.status,
    d.original_value,
    d.current_value,
    d.depreciation_rate,
    d.overall_condition,
    d.battery_health,
    d.grade,
    d.grade_score,
    d.purchase_date,
    d.registration_date,
    COUNT(p.policy_id) as total_policies,
    COUNT(CASE WHEN p.status = 'active' THEN 1 END) as active_policies,
    COUNT(cl.claim_id) as total_claims,
    COUNT(CASE WHEN cl.status = 'approved' THEN 1 END) as approved_claims,
    SUM(CASE WHEN cl.status = 'approved' THEN cl.approved_amount ELSE 0 END) as total_claim_amount,
    AVG(CASE WHEN cl.status = 'approved' THEN cl.approved_amount END) as avg_claim_amount
FROM devices d
LEFT JOIN customers c ON d.current_owner_id = c.customer_id
LEFT JOIN policies p ON d.device_id = p.device_id
LEFT JOIN claims cl ON d.device_id = cl.device_id
GROUP BY d.device_id, d.imei, d.model, d.manufacturer, d.brand, d.current_owner_id,
         c.first_name, c.last_name, d.status, d.original_value, d.current_value,
         d.depreciation_rate, d.overall_condition, d.battery_health, d.grade,
         d.grade_score, d.purchase_date, d.registration_date;

-- Policy analytics view
CREATE VIEW policy_analytics AS
SELECT 
    p.policy_id,
    p.policy_number,
    p.policy_type,
    p.coverage_type,
    p.status,
    p.final_premium,
    p.currency,
    p.risk_score,
    p.fraud_risk_score,
    p.effective_date,
    p.expiration_date,
    c.customer_number,
    c.first_name || ' ' || c.last_name as customer_name,
    c.risk_level as customer_risk_level,
    c.risk_score as customer_risk_score,
    d.imei,
    d.model as device_model,
    d.current_value as device_current_value,
    COUNT(cl.claim_id) as total_claims,
    COUNT(CASE WHEN cl.status = 'approved' THEN 1 END) as approved_claims,
    SUM(CASE WHEN cl.status = 'approved' THEN cl.approved_amount ELSE 0 END) as total_claim_amount,
    AVG(CASE WHEN cl.status = 'approved' THEN cl.approved_amount END) as avg_claim_amount,
    CASE 
        WHEN COUNT(cl.claim_id) > 0 THEN 
            SUM(CASE WHEN cl.status = 'approved' THEN cl.approved_amount ELSE 0 END) / p.final_premium
        ELSE 0 
    END as loss_ratio
FROM policies p
JOIN customers c ON p.customer_id = c.customer_id
JOIN devices d ON p.device_id = d.device_id
LEFT JOIN claims cl ON p.policy_id = cl.policy_id
GROUP BY p.policy_id, p.policy_number, p.policy_type, p.coverage_type, p.status,
         p.final_premium, p.currency, p.risk_score, p.fraud_risk_score,
         p.effective_date, p.expiration_date, c.customer_number, c.first_name,
         c.last_name, c.risk_level, c.risk_score, d.imei, d.model, d.current_value;

-- Analytics events summary view
CREATE VIEW analytics_events_summary AS
SELECT 
    event_type,
    source,
    DATE(timestamp) as event_date,
    COUNT(*) as event_count,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT device_id) as unique_devices,
    COUNT(DISTINCT session_id) as unique_sessions,
    AVG(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - timestamp))) as avg_age_seconds
FROM analytics_events
GROUP BY event_type, source, DATE(timestamp);

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

-- Function to generate customer number
CREATE OR REPLACE FUNCTION generate_customer_number()
RETURNS VARCHAR(50) AS $$
BEGIN
    RETURN 'CUST-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
           LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0');
END;
$$ LANGUAGE plpgsql;

-- Function to generate policy number
CREATE OR REPLACE FUNCTION generate_policy_number()
RETURNS VARCHAR(50) AS $$
BEGIN
    RETURN 'POL-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
           LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0');
END;
$$ LANGUAGE plpgsql;

-- Function to generate claim number
CREATE OR REPLACE FUNCTION generate_claim_number()
RETURNS VARCHAR(50) AS $$
BEGIN
    RETURN 'CLM-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
           LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0');
END;
$$ LANGUAGE plpgsql;

-- Function to calculate customer risk score
CREATE OR REPLACE FUNCTION calculate_customer_risk_score(
    customer_uuid UUID
)
RETURNS DECIMAL AS $$
DECLARE
    risk_score DECIMAL := 0.5; -- Base risk score
    income_factor DECIMAL;
    claim_factor DECIMAL;
    credit_factor DECIMAL;
    customer_record RECORD;
BEGIN
    -- Get customer information
    SELECT annual_income, credit_score INTO customer_record
    FROM customers 
    WHERE customer_id = customer_uuid;
    
    -- Income factor (lower income = higher risk)
    IF customer_record.annual_income < 30000 THEN
        income_factor := 0.8;
    ELSIF customer_record.annual_income < 60000 THEN
        income_factor := 0.6;
    ELSIF customer_record.annual_income < 100000 THEN
        income_factor := 0.4;
    ELSE
        income_factor := 0.2;
    END IF;
    
    -- Credit factor (lower credit score = higher risk)
    IF customer_record.credit_score < 500 THEN
        credit_factor := 0.9;
    ELSIF customer_record.credit_score < 600 THEN
        credit_factor := 0.7;
    ELSIF customer_record.credit_score < 700 THEN
        credit_factor := 0.5;
    ELSIF customer_record.credit_score < 800 THEN
        credit_factor := 0.3;
    ELSE
        credit_factor := 0.1;
    END IF;
    
    -- Claim history factor
    SELECT COUNT(*) INTO claim_factor
    FROM claims c
    JOIN policies p ON c.policy_id = p.policy_id
    WHERE p.customer_id = customer_uuid AND c.status = 'approved';
    
    claim_factor := LEAST(1.0, claim_factor * 0.2);
    
    -- Calculate final risk score
    risk_score := (income_factor * 0.3 + credit_factor * 0.3 + claim_factor * 0.4);
    
    RETURN LEAST(1.0, GREATEST(0.0, risk_score));
END;
$$ LANGUAGE plpgsql;

-- Function to get customer statistics
CREATE OR REPLACE FUNCTION get_customer_statistics(
    customer_uuid UUID
)
RETURNS JSON AS $$
DECLARE
    result JSON;
BEGIN
    SELECT json_build_object(
        'customer_id', c.customer_id,
        'customer_number', c.customer_number,
        'customer_name', c.first_name || ' ' || c.last_name,
        'account_status', c.account_status,
        'risk_level', c.risk_level,
        'risk_score', c.risk_score,
        'kyc_status', c.kyc_status,
        'total_devices', COUNT(DISTINCT d.device_id),
        'active_devices', COUNT(DISTINCT CASE WHEN d.status = 'active' THEN d.device_id END),
        'total_policies', COUNT(DISTINCT p.policy_id),
        'active_policies', COUNT(DISTINCT CASE WHEN p.status = 'active' THEN p.policy_id END),
        'total_claims', COUNT(DISTINCT cl.claim_id),
        'approved_claims', COUNT(DISTINCT CASE WHEN cl.status = 'approved' THEN cl.claim_id END),
        'total_claim_amount', SUM(CASE WHEN cl.status = 'approved' THEN cl.approved_amount ELSE 0 END),
        'avg_claim_amount', AVG(CASE WHEN cl.status = 'approved' THEN cl.approved_amount END),
        'registration_date', c.registration_date,
        'last_login_date', c.last_login_date
    ) INTO result
    FROM customers c
    LEFT JOIN devices d ON c.customer_id = d.current_owner_id
    LEFT JOIN policies p ON c.customer_id = p.customer_id
    LEFT JOIN claims cl ON c.customer_id = cl.customer_id
    WHERE c.customer_id = customer_uuid
    GROUP BY c.customer_id, c.customer_number, c.first_name, c.last_name,
             c.account_status, c.risk_level, c.risk_score, c.kyc_status,
             c.registration_date, c.last_login_date;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql; 