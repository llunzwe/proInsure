-- 015_claim_workflow_stages.sql

-- Detailed claim processing workflow
CREATE TABLE claim_workflow_stages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    
    stage VARCHAR(50) NOT NULL CHECK (stage IN (
        'submission', 'triage', 'document_review', 'investigation',
        'assessment', 'approval_pending', 'approved', 'payment_processing',
        'settlement', 'closure', 'appeal', 'reopened'
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
