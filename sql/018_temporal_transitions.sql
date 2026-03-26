-- 018_temporal_transitions.sql

-- 4D bitemporal state transitions
CREATE TABLE temporal_transitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Entity reference (polymorphic)
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('policy', 'claim', 'device', 'party', 'premium_payment', 'reserve')),
    entity_id UUID NOT NULL,
    
    -- State change
    from_state VARCHAR(50) NOT NULL,
    to_state VARCHAR(50) NOT NULL,
    transition_event VARCHAR(100) NOT NULL, -- The event that caused transition
    
    -- 4D temporal
    system_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, -- When recorded
    valid_time_start TIMESTAMPTZ NOT NULL, -- When effective in business time
    valid_time_end TIMESTAMPTZ, -- Until when (null if current)
    
    -- Authorization
    triggered_by UUID REFERENCES parties(id),
    approved_by UUID REFERENCES parties(id), -- For 4-eyes
    
    -- Context
    reason TEXT,
    correlation_id UUID, -- Link to event or transaction
    metadata JSONB,
    
    -- Integrity
    immutable_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Helper function to record transitions
CREATE OR REPLACE FUNCTION record_transition(
    p_entity_type VARCHAR,
    p_entity_id UUID,
    p_from_state VARCHAR,
    p_to_state VARCHAR,
    p_event VARCHAR,
    p_triggered_by UUID,
    p_valid_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    p_reason TEXT DEFAULT NULL
) RETURNS UUID AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO temporal_transitions (
        entity_type, entity_id, from_state, to_state, transition_event,
        triggered_by, valid_time_start, reason, immutable_hash
    ) VALUES (
        p_entity_type, p_entity_id, p_from_state, p_to_state, p_event,
        p_triggered_by, p_valid_time, p_reason,
        encode(digest(random()::text, 'sha256'), 'hex') -- Simplified hash
    ) RETURNING id INTO v_id;
    
    RETURN v_id;
END;
$$ LANGUAGE plpgsql;

CREATE INDEX idx_transitions_entity ON temporal_transitions(entity_type, entity_id);
CREATE INDEX idx_transitions_time ON temporal_transitions(valid_time_start, valid_time_end);
CREATE INDEX idx_transitions_system ON temporal_transitions(system_time);
