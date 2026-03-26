-- 025_immutable_event_store.sql

-- Core immutable ledger - Datomic-style datoms
CREATE TABLE immutable_events (
    -- Physical ordering
    event_id BIGSERIAL,
    event_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Event identification
    event_type VARCHAR(100) NOT NULL,
    event_version VARCHAR(10) DEFAULT '1.0',
    
    -- Payload
    payload JSONB NOT NULL,
    payload_hash VARCHAR(64) NOT NULL,
    
    -- Hash chaining for cryptographic integrity
    previous_hash VARCHAR(64),
    event_hash VARCHAR(64) NOT NULL, -- SHA-256 of (type + payload_hash + previous_hash)
    
    -- Datomic datom fields (E-A-V-Tx-Op)
    datom_entity_id UUID NOT NULL, -- The entity (policy_id, claim_id, etc.)
    datom_attribute VARCHAR(200) NOT NULL, -- e.g., "policy.status", "claim.amount"
    datom_value JSONB NOT NULL, -- New value
    datom_operation CHAR(1) NOT NULL CHECK (datom_operation IN ('+', '-')), -- + = assert, - = retract
    datom_valid_time TIMESTAMPTZ NOT NULL, -- Business time when fact is true
    
    -- Causality and correlation
    correlation_id UUID, -- Business transaction ID
    causation_id UUID, -- Parent event that caused this (event_id of previous)
    session_id UUID, -- User session
    
    -- Merkle tree for batch anchoring to blockchain
    merkle_batch_id UUID,
    merkle_tree_path ltree,
    
    -- Source
    source_system VARCHAR(100) DEFAULT 'proinsure-ledger',
    source_ip INET,
    
    -- Partitioning (by event_time)
    PRIMARY KEY (event_id, event_time)
) PARTITION BY RANGE (event_time);

-- Create monthly partitions
CREATE TABLE immutable_events_2026_03 PARTITION OF immutable_events
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
-- ... add more partitions as needed

-- Convert to hypertable for time-series optimization
-- SELECT create_hypertable('immutable_events', 'event_time', chunk_time_interval => INTERVAL '1 day');

-- EAVT Index (Entity-Attribute-Value-Transaction) - Primary Datomic access pattern
CREATE INDEX idx_events_eavt ON immutable_events(datom_entity_id, datom_attribute, datom_valid_time, event_time);

-- AVET Index (Attribute-Value-Entity-Transaction) - For value lookups
CREATE INDEX idx_events_avet ON immutable_events(datom_attribute, datom_value, datom_entity_id, event_time);

-- AEVT Index (Attribute-Entity-Value-Transaction)
CREATE INDEX idx_events_aevt ON immutable_events(datom_attribute, datom_entity_id, datom_valid_time);

-- VAET Index (Value-Attribute-Entity-Transaction) - For reverse lookups
CREATE INDEX idx_events_vaet ON immutable_events(datom_value, datom_attribute, datom_entity_id);

-- Correlation and causation
CREATE INDEX idx_events_correlation ON immutable_events(correlation_id);
CREATE INDEX idx_events_causation ON immutable_events(causation_id);
CREATE INDEX idx_events_type ON immutable_events(event_type, event_time);
CREATE INDEX idx_events_merkle ON immutable_events(merkle_batch_id);

-- Function to calculate event hash with chaining
CREATE OR REPLACE FUNCTION calculate_event_hash() RETURNS TRIGGER AS $$
DECLARE
    v_prev_hash VARCHAR(64);
BEGIN
    -- Get previous hash
    SELECT event_hash INTO v_prev_hash
    FROM immutable_events
    ORDER BY event_id DESC
    LIMIT 1;
    
    IF v_prev_hash IS NULL THEN
        v_prev_hash := '0' * 64; -- Genesis
    END IF;
    
    NEW.previous_hash := v_prev_hash;
    NEW.event_hash := encode(digest(
        NEW.event_type || NEW.payload_hash || v_prev_hash,
        'sha256'
    ), 'hex');
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_event_hash BEFORE INSERT ON immutable_events
    FOR EACH ROW EXECUTE FUNCTION calculate_event_hash();

-- Materialized view for current state (optional optimization)
-- CREATE MATERIALIZED VIEW current_entity_state AS ...
