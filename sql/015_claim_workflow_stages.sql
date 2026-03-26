-- 015_claim_workflow_stages.sql

-- Detailed claim processing workflow
CREATE TABLE claim_workflow_stages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    
    stage VARCHAR(50) NOT NULL CHECK (stage IN (
        'submission', 'triage', 'document_review', 'investigation',
        'assessment', 'approval_pending', 'approved', 'payment_processing',
        'settlement', 'closure', 'appeal', 'reopened',
        -- Repair-centric stages
        'repair_authorization', 'parts_availability_check', 'repair_in_progress',
        'quality_assurance', 'customer_acceptance', 'device_returned',
        -- Total loss stages
        'salvage_collection', 'replacement_selection', 'data_migration_verification'
    )),
    
    stage_order INTEGER NOT NULL, -- Sequence number
    
    -- Timing
    entered_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    exited_at TIMESTAMPTZ,
    duration_seconds INTEGER GENERATED ALWAYS AS (
        EXTRACT(EPOCH FROM (COALESCE(exited_at, CURRENT_TIMESTAMP) - entered_at))
    ) STORED,
    
    -- Assignment
    assigned_to UUID REFERENCES parties(id),
    assigned_role VARCHAR(50),
    
    -- SLA tracking
    sla_deadline TIMESTAMPTZ, -- Service level agreement deadline
    sla_breached BOOLEAN DEFAULT FALSE,
    
    -- Stage details
    notes TEXT,
    checklist_completed JSONB, -- Key-value of checklist items
    
    -- 4-eyes approval
    requires_approval BOOLEAN DEFAULT FALSE,
    approved_by UUID REFERENCES parties(id),
    approved_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Workflow stage history (if reprocessing)
CREATE TABLE claim_stage_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_stage_id UUID REFERENCES claim_workflow_stages(id),
    previous_status VARCHAR(50),
    new_status VARCHAR(50),
    changed_by UUID REFERENCES parties(id),
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    reason TEXT
);

-- Indexes
CREATE INDEX idx_workflow_claim ON claim_workflow_stages(claim_id);
CREATE INDEX idx_workflow_stage ON claim_workflow_stages(stage);
CREATE INDEX idx_workflow_active ON claim_workflow_stages(claim_id, exited_at) WHERE exited_at IS NULL;
CREATE INDEX idx_workflow_sla ON claim_workflow_stages(sla_deadline) WHERE sla_breached = FALSE;

-- SLA configuration by severity tier
CREATE TABLE claim_sla_configurations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Classification
    severity_tier VARCHAR(20) NOT NULL CHECK (severity_tier IN ('critical', 'standard', 'economy')),
    claim_type claim_type,
    device_category device_category,
    
    -- SLA targets (hours)
    acknowledgment_sla_hours INTEGER DEFAULT 4,
    triage_sla_hours INTEGER DEFAULT 24,
    assessment_sla_hours INTEGER DEFAULT 48,
    approval_sla_hours INTEGER DEFAULT 24,
    repair_sla_hours INTEGER DEFAULT 72,
    replacement_sla_hours INTEGER DEFAULT 48,
    settlement_sla_hours INTEGER DEFAULT 24,
    
    -- Critical business device settings
    is_business_device BOOLEAN DEFAULT FALSE,
    express_processing_available BOOLEAN DEFAULT FALSE,
    
    -- Escalation
    escalation_threshold_hours INTEGER,
    escalation_notify_roles VARCHAR(50)[],
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Insert default SLA configurations
INSERT INTO claim_sla_configurations (severity_tier, is_business_device, triage_sla_hours, repair_sla_hours, replacement_sla_hours) VALUES
('critical', TRUE, 1, 4, 4),     -- Critical business device: 4-hour replacement
('standard', FALSE, 24, 72, 48), -- Standard: 24h triage, 72h repair
('economy', FALSE, 48, 168, 120); -- Economy: 48h triage, 7 day repair

-- Stage escalation tracking
CREATE TABLE claim_stage_escalations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_stage_id UUID NOT NULL REFERENCES claim_workflow_stages(id),
    
    escalation_level INTEGER DEFAULT 1,
    escalated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    escalated_by UUID REFERENCES parties(id),
    
    reason TEXT,
    action_taken VARCHAR(100),
    resolved_at TIMESTAMPTZ,
    resolved_by UUID REFERENCES parties(id),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_escalations_stage ON claim_stage_escalations(workflow_stage_id);
