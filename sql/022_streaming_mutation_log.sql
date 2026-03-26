-- 022_streaming_mutation_log.sql

-- Kafka-compatible streaming log
CREATE TABLE streaming_mutation_log (
    id BIGSERIAL,
    mutation_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Event details
    mutation_type VARCHAR(100) NOT NULL, -- PolicyIssued, ClaimApproved, etc.
    source_table VARCHAR(100) NOT NULL,
    source_record_id UUID NOT NULL,
    
    -- Payload
    payload JSONB NOT NULL,
    payload_hash VARCHAR(64) NOT NULL,
    schema_version VARCHAR(20) DEFAULT '1.0',
    
    -- Temporal tracking
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    sequence_number BIGINT GENERATED ALWAYS AS IDENTITY, -- Strict ordering
    
    -- Delivery tracking
    delivery_status VARCHAR(20) DEFAULT 'pending' CHECK (delivery_status IN ('pending', 'delivered', 'failed', 'retrying')),
    delivery_attempts INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    
    -- Subscriber tracking (array of successful deliveries)
    delivered_to UUID[],
    
    -- Partitioning key for TimescaleDB
    partition_key INTEGER GENERATED ALWAYS AS (EXTRACT(YEAR FROM created_at) * 100 + EXTRACT(MONTH FROM created_at)) STORED
) PARTITION BY RANGE (created_at);

-- Create monthly partitions (example for 2026)
CREATE TABLE streaming_mutation_log_2026_01 PARTITION OF streaming_mutation_log
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE streaming_mutation_log_2026_02 PARTITION OF streaming_mutation_log
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
-- ... continue for all months

-- Convert to hypertable for better time-series performance
-- Note: If using native partitioning above, TimescaleDB hypertable might conflict, choose one approach
-- SELECT create_hypertable('streaming_mutation_log', 'created_at', chunk_time_interval => INTERVAL '1 day');

-- Streaming subscribers
CREATE TABLE streaming_subscribers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    subscriber_type VARCHAR(20) CHECK (subscriber_type IN ('KAFKA', 'WEBHOOK', 'SQS', 'PUBSUB')),
    
    -- Subscription pattern (which events)
    event_pattern VARCHAR(200), -- Regex or specific event types, e.g., "Policy.*|Claim.*"
    tables_of_interest VARCHAR(100)[], -- {'policies', 'claims'}
    
    -- Delivery config
    delivery_config JSONB NOT NULL, -- URL for webhook, topic for Kafka, etc.
    headers JSONB, -- Auth headers, etc.
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'error')),
    last_processed_sequence BIGINT DEFAULT 0,
    last_processed_at TIMESTAMPTZ,
    
    -- Error handling
    error_count INTEGER DEFAULT 0,
    last_error TEXT,
    circuit_breaker_open BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Delivery attempts log
CREATE TABLE streaming_delivery_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mutation_id UUID REFERENCES streaming_mutation_log(mutation_id),
    subscriber_id UUID REFERENCES streaming_subscribers(id),
    attempted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN,
    http_status INTEGER,
    response_body TEXT,
    error_message TEXT
);

CREATE INDEX idx_mutation_log_status ON streaming_mutation_log(delivery_status);
CREATE INDEX idx_mutation_log_sequence ON streaming_mutation_log(sequence_number);
CREATE INDEX idx_mutation_log_source ON streaming_mutation_log(source_table, source_record_id);
CREATE INDEX idx_mutation_log_created ON streaming_mutation_log(created_at);
