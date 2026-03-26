-- 013_claims.sql

-- Insurance claims
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_number VARCHAR(50) NOT NULL UNIQUE,
    
    -- Relationships
    policy_id UUID NOT NULL REFERENCES policies(id),
    device_id UUID REFERENCES devices(id), -- Denormalized for performance
    
    -- Claim details
    claim_type claim_type NOT NULL,
    incident_date DATE NOT NULL,
    reported_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT,
    incident_location VARCHAR(255),
    
    -- Amounts
    claimed_amount DECIMAL(28,8) NOT NULL,
    approved_amount DECIMAL(28,8),
    deductible_applied DECIMAL(28,8) DEFAULT 0,
    excess_amount DECIMAL(28,8) DEFAULT 0, -- Non-covered portion
    
    -- Damage categorization taxonomy
    damage_severity VARCHAR(20) CHECK (damage_severity IN ('minor', 'moderate', 'severe', 'total_loss')),
    damage_type_detailed VARCHAR(50), -- hairline_crack, shattered_functional, lcd_bleeding, touch_unresponsive, water_lci_only, water_corrosion, water_non_functional
    is_cosmetic_only BOOLEAN DEFAULT FALSE,
    is_functional_damage BOOLEAN DEFAULT TRUE,
    
    -- Automated triage
    triage_score INTEGER CHECK (triage_score BETWEEN 0 AND 100), -- AI confidence score
    triage_decision VARCHAR(20) CHECK (triage_decision IN ('instant_approve', 'fast_track', 'standard_review', 'manual_review', 'fraud_flag')),
    instant_approval_eligible BOOLEAN DEFAULT FALSE,
    instant_approval_threshold DECIMAL(28,8) DEFAULT 300, -- Auto-approve under this amount
    
    -- Fraud indicators (automated)
    fraud_flags_triggered JSONB, -- Array of triggered rules
    imei_mismatch BOOLEAN DEFAULT FALSE,
    claim_within_hours_of_policy_start INTEGER, -- Calculated field
    multiple_devices_same_incident BOOLEAN DEFAULT FALSE,
    
    -- Customer experience
    fnol_submitted_via VARCHAR(20) DEFAULT 'app' CHECK (fnol_submitted_via IN ('app', 'web', 'phone', 'agent')),
    preferred_repair_channel VARCHAR(20) CHECK (preferred_repair_channel IN ('mail_in', 'walk_in', 'on_site', 'express_replacement')),
    push_notifications_enabled BOOLEAN DEFAULT TRUE,
    
    -- Third party liability
    liable_third_party_id UUID REFERENCES parties(id),
    police_report_number VARCHAR(50),
    police_report_filed_at TIMESTAMPTZ
    
    -- Status workflow
    status claim_status DEFAULT 'submitted',
    
    -- Timeline
    submitted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    assigned_at TIMESTAMPTZ,
    assessed_at TIMESTAMPTZ,
    decision_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    reopened_at TIMESTAMPTZ,
    
    -- Assignment
    assessor_party_id UUID REFERENCES parties(id),
    adjuster_party_id UUID REFERENCES parties(id),
    handled_by UUID REFERENCES parties(id), -- Internal handler
    
    -- Fraud and risk
    fraud_score DECIMAL(5,2) CHECK (fraud_score >= 0 AND fraud_score <= 100),
    fraud_flags JSONB, -- Array of triggered rules
    investigation_required BOOLEAN DEFAULT FALSE,
    investigation_notes TEXT,
    
    -- Payment
    payout_movement_id UUID REFERENCES value_movements(id),
    settlement_method settlement_method,
    settlement_amount DECIMAL(28,8),
    
    -- Real-time authorization link
    real_time_posting_id UUID, -- Link to real_time_postings table
    
    -- Bitemporal
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity',
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Integrity
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    CONSTRAINT valid_incident_date CHECK (incident_date <= CURRENT_DATE),
    CONSTRAINT valid_amounts CHECK (approved_amount IS NULL OR approved_amount <= claimed_amount)
);

-- Claim reserve link (explicit many-to-many if needed)
CREATE TABLE claim_reserve_links (
    claim_id UUID REFERENCES claims(id),
    reserve_id UUID REFERENCES reserves(id),
    PRIMARY KEY (claim_id, reserve_id)
);

-- Indexes
CREATE INDEX idx_claims_number ON claims(claim_number);
CREATE INDEX idx_claims_policy ON claims(policy_id);
CREATE INDEX idx_claims_status ON claims(status);
CREATE INDEX idx_claims_type ON claims(claim_type);
CREATE INDEX idx_claims_fraud ON claims(fraud_score) WHERE fraud_score > 50;
CREATE INDEX idx_claims_dates ON claims(incident_date, submitted_at);
CREATE INDEX idx_claims_current ON claims(is_current) WHERE is_current = TRUE;

-- Trigger
CREATE TRIGGER trg_claims_hash BEFORE INSERT OR UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

-- Claim appeals table
CREATE TABLE claim_appeals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    appeal_number VARCHAR(50) NOT NULL UNIQUE,
    
    -- Appeal details
    appeal_reason TEXT NOT NULL,
    appellant_party_id UUID REFERENCES parties(id),
    original_decision VARCHAR(20) NOT NULL, -- approved, rejected
    requested_amount DECIMAL(28,8),
    
    -- Status workflow
    status VARCHAR(20) DEFAULT 'submitted' CHECK (status IN ('submitted', 'under_review', 'upheld', 'rejected', 'partially_upheld')),
    
    -- Timeline
    submitted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    
    -- Reviewer
    reviewed_by UUID REFERENCES parties(id),
    review_notes TEXT,
    
    -- Outcome
    final_decision VARCHAR(20),
    final_amount DECIMAL(28,8),
    movement_id UUID REFERENCES value_movements(id), -- For any additional payment
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    immutable_hash VARCHAR(64) NOT NULL
);

CREATE TRIGGER trg_claim_appeals_hash BEFORE INSERT OR UPDATE ON claim_appeals
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE INDEX idx_claim_appeals_claim ON claim_appeals(claim_id);
CREATE INDEX idx_claim_appeals_status ON claim_appeals(status);

-- Salvage recoveries table (recovery from damaged/repaired devices)
CREATE TABLE salvage_recoveries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    
    -- Recovery details
    recovery_reference VARCHAR(50) NOT NULL UNIQUE,
    recovery_type VARCHAR(50) NOT NULL CHECK (recovery_type IN ('device_sale', 'parts_sale', 'scrap_value')),
    
    -- Device info
    device_description TEXT,
    device_condition VARCHAR(20), -- damaged, repaired, parts_only
    
    -- Financial
    estimated_value DECIMAL(28,8),
    actual_recovery_amount DECIMAL(28,8),
    recovery_costs DECIMAL(28,8) DEFAULT 0, -- Auction fees, transport, etc.
    net_recovery DECIMAL(28,8), -- actual_recovery - recovery_costs
    currency CHAR(3) DEFAULT 'USD',
    
    -- Status
    status recovery_status DEFAULT 'pending',
    
    -- Timeline
    identified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    recovered_at TIMESTAMPTZ,
    
    -- Buyer/Recovery agent
    recovered_from_party_id UUID REFERENCES parties(id),
    recovery_method VARCHAR(50), -- auction, direct_sale, trade_in
    
    -- Accounting
    movement_id UUID REFERENCES value_movements(id),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL
);

CREATE TRIGGER trg_salvage_recoveries_hash BEFORE INSERT OR UPDATE ON salvage_recoveries
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE INDEX idx_salvage_claim ON salvage_recoveries(claim_id);
CREATE INDEX idx_salvage_status ON salvage_recoveries(status);

-- Subrogation recoveries table (recovery from third parties)
CREATE TABLE subrogation_recoveries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    
    -- Subrogation details
    subrogation_reference VARCHAR(50) NOT NULL UNIQUE,
    subrogation_type VARCHAR(50) NOT NULL CHECK (subrogation_type IN ('third_party_liability', 'theft_recovery', 'warranty_claim', 'other_insurer')),
    
    -- Third party info
    liable_party_id UUID REFERENCES parties(id),
    liable_party_name VARCHAR(255),
    liable_party_contact TEXT,
    liable_party_insurer VARCHAR(100),
    liable_party_policy_ref VARCHAR(100),
    
    -- Claim details
    claim_amount DECIMAL(28,8) NOT NULL,
    recovery_amount DECIMAL(28,8),
    recovery_costs DECIMAL(28,8) DEFAULT 0, -- Legal fees, etc.
    net_recovery DECIMAL(28,8),
    currency CHAR(3) DEFAULT 'USD',
    
    -- Status
    status recovery_status DEFAULT 'pending',
    
    -- Timeline
    identified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    demand_sent_at TIMESTAMPTZ,
    recovered_at TIMESTAMPTZ,
    written_off_at TIMESTAMPTZ,
    
    -- Recovery method
    recovery_method VARCHAR(50), -- direct_demand, legal_action, arbitration
    legal_reference VARCHAR(100),
    
    -- Notes
    notes TEXT,
    
    -- Accounting
    movement_id UUID REFERENCES value_movements(id),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES parties(id),
    immutable_hash VARCHAR(64) NOT NULL
);

CREATE TRIGGER trg_subrogation_recoveries_hash BEFORE INSERT OR UPDATE ON subrogation_recoveries
    FOR EACH ROW EXECUTE FUNCTION generate_row_hash();

CREATE INDEX idx_subrogation_claim ON subrogation_recoveries(claim_id);
CREATE INDEX idx_subrogation_status ON subrogation_recoveries(status);
CREATE INDEX idx_subrogation_liable ON subrogation_recoveries(liable_party_id);


-- Claim communication log (correspondence with policyholder)
CREATE TABLE claim_communications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    
    communication_type VARCHAR(20) NOT NULL CHECK (communication_type IN ('email', 'sms', 'push', 'phone', 'letter')),
    direction VARCHAR(10) NOT NULL CHECK (direction IN ('inbound', 'outbound')),
    
    -- Content
    subject VARCHAR(255),
    body TEXT,
    attachment_ids UUID[], -- References to documents
    
    -- Metadata
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    opened_at TIMESTAMPTZ,
    
    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'delivered', 'failed', 'bounced')),
    
    -- Sender/Recipient
    from_party_id UUID REFERENCES parties(id),
    to_party_id UUID REFERENCES parties(id),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_claim_communications_claim ON claim_communications(claim_id);
CREATE INDEX idx_claim_communications_type ON claim_communications(communication_type);

-- Real-time claim authorization decisions
CREATE TABLE claim_authorization_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    
    -- Decision details
    decision_type VARCHAR(20) NOT NULL CHECK (decision_type IN ('instant', 'fast_track', 'manual', 'declined')),
    decision VARCHAR(20) NOT NULL CHECK (decision IN ('approved', 'declined', 'pending_review', 'information_required')),
    
    -- AI/Rule analysis
    ai_confidence_score DECIMAL(5,2),
    rule_engine_results JSONB, -- Which rules fired
    image_recognition_results JSONB, -- Damage classification
    ocr_extraction_results JSONB, -- Police report parsing
    
    -- Amounts
    recommended_payout_amount DECIMAL(28,8),
    recommended_deductible DECIMAL(28,8),
    
    -- Customer self-service
    settlement_options_presented JSONB, -- cash, repair, replacement options
    customer_selected_option VARCHAR(20),
    
    -- Timing
    decision_made_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    decision_latency_ms INTEGER, -- Performance tracking
    
    -- Override
    overridden_by UUID REFERENCES parties(id),
    override_reason TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_authorization_claim ON claim_authorization_decisions(claim_id);
CREATE INDEX idx_authorization_decision ON claim_authorization_decisions(decision);
