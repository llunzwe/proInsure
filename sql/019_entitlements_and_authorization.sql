-- 019_entitlements_and_authorization.sql

-- Access control entitlements
CREATE TABLE entitlements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Resource specification
    resource_type VARCHAR(50) NOT NULL CHECK (resource_type IN ('policy', 'claim', 'device', 'payment', 'report', 'admin')),
    resource_id UUID, -- NULL = all resources of type
    
    -- Action
    action VARCHAR(50) NOT NULL CHECK (action IN (
        'view', 'create', 'update', 'delete', 'approve', 
        'file_claim', 'approve_payout', 'cancel_policy', 'issue_refund',
        'assess_claim', 'repair_authorization', 'admin_access'
    )),
    
    -- Constraints
    max_amount DECIMAL(28,8), -- For financial actions
    requires_2fa BOOLEAN DEFAULT FALSE,
    requires_approval BOOLEAN DEFAULT FALSE, -- 4-eyes
    approval_threshold DECIMAL(28,8), -- When approval needed
    
    -- Temporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    
    UNIQUE(party_id, resource_type, resource_id, action, valid_from)
);

-- Authorization attempts log
CREATE TABLE authorization_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failure', 'mfa_required', 'approval_required')),
    failure_reason TEXT,
    
    -- Security context
    authentication_method VARCHAR(50), -- password, biometric, mfa, certificate
    digital_signature BYTEA,
    ip_address INET,
    device_fingerprint VARCHAR(256),
    geo_location_lat DECIMAL(10,8),
    geo_location_lon DECIMAL(11,8),
    
    -- Request details
    request_payload JSONB,
    response_payload JSONB
);

-- Role-based entitlements template
CREATE TABLE entitlement_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    permissions JSONB NOT NULL, -- Array of {resource, action, constraints}
    is_active BOOLEAN DEFAULT TRUE
);

-- Party-role assignments
CREATE TABLE party_role_assignments (
    party_id UUID REFERENCES parties(id),
    role_id UUID REFERENCES entitlement_roles(id),
    assigned_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES parties(id),
    valid_until TIMESTAMPTZ,
    PRIMARY KEY (party_id, role_id)
);

CREATE INDEX idx_entitlements_party ON entitlements(party_id);
CREATE INDEX idx_entitlements_resource ON entitlements(resource_type, resource_id);
CREATE INDEX idx_auth_logs_party ON authorization_logs(party_id);
CREATE INDEX idx_auth_logs_time ON authorization_logs(timestamp);
CREATE INDEX idx_auth_logs_failed ON authorization_logs(status, timestamp) WHERE status = 'failure';
