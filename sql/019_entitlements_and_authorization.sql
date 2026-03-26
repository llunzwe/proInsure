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

-- Segregation of duties - prohibited combinations
CREATE TABLE segregation_of_duties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    
    -- First action
    first_resource_type VARCHAR(50) NOT NULL,
    first_action VARCHAR(50) NOT NULL,
    
    -- Second action (cannot be done by same person)
    second_resource_type VARCHAR(50) NOT NULL,
    second_action VARCHAR(50) NOT NULL,
    
    -- Time window for checking (NULL = any time)
    time_window_hours INTEGER,
    
    -- Exemptions
    allow_with_approval BOOLEAN DEFAULT TRUE,
    approver_role_id UUID REFERENCES entitlement_roles(id),
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Check segregation of duties violations
CREATE TABLE sod_violations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    segregation_rule_id UUID NOT NULL REFERENCES segregation_of_duties(id),
    party_id UUID NOT NULL REFERENCES parties(id),
    
    first_action_at TIMESTAMPTZ NOT NULL,
    first_action_resource_id UUID,
    
    second_action_at TIMESTAMPTZ,
    second_action_resource_id UUID,
    
    status VARCHAR(20) DEFAULT 'detected' CHECK (status IN ('detected', 'approved', 'rejected', 'escalated')),
    approved_by UUID REFERENCES parties(id),
    approved_at TIMESTAMPTZ,
    notes TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Role hierarchy for composite roles
CREATE TABLE role_hierarchy (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_role_id UUID NOT NULL REFERENCES entitlement_roles(id),
    child_role_id UUID NOT NULL REFERENCES entitlement_roles(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(parent_role_id, child_role_id)
);

-- Session management
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    
    -- Session tokens
    access_token_hash VARCHAR(64) NOT NULL, -- SHA-256 of token
    refresh_token_hash VARCHAR(64),
    
    -- Session details
    started_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    last_activity_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMPTZ,
    
    -- Context
    ip_address INET,
    user_agent TEXT,
    device_fingerprint VARCHAR(256),
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    terminated_reason VARCHAR(50), -- logout, expired, revoked, security
    
    -- MFA
    mfa_verified BOOLEAN DEFAULT FALSE,
    mfa_method VARCHAR(20) -- totp, sms, email, biometric
);

CREATE INDEX idx_sessions_party ON sessions(party_id);
CREATE INDEX idx_sessions_active ON sessions(party_id, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_sessions_expires ON sessions(expires_at) WHERE is_active = TRUE;

-- Session blacklist for revoked tokens before expiry
CREATE TABLE session_blacklist (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    revoked_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    revoked_by UUID REFERENCES parties(id),
    reason VARCHAR(100),
    expires_at TIMESTAMPTZ NOT NULL -- When token would have expired naturally
);

CREATE INDEX idx_session_blacklist_hash ON session_blacklist(token_hash);
CREATE INDEX idx_session_blacklist_expires ON session_blacklist(expires_at);
