-- 027_kernel_wiring_and_integrity.sql

-- Audit log for all DDL changes (schema versioning)
CREATE TABLE schema_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action VARCHAR(50) NOT NULL, -- CREATE, ALTER, DROP
    object_type VARCHAR(50), -- TABLE, INDEX, FUNCTION, etc.
    object_name VARCHAR(200),
    ddl_command TEXT,
    performed_by VARCHAR(100) DEFAULT CURRENT_USER,
    performed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    application_version VARCHAR(20)
);

-- Function to log DDL
CREATE OR REPLACE FUNCTION log_ddl() RETURNS event_trigger AS $$
BEGIN
    INSERT INTO schema_audit_log (action, object_type, object_name, ddl_command)
    VALUES (
        tg_tag,
        tg_event,
        'unknown', -- Object name extraction requires parsing
        current_query()
    );
END;
$$ LANGUAGE plpgsql;

-- Create event trigger for DDL logging (optional - requires superuser)
-- CREATE EVENT TRIGGER ddl_trigger ON ddl_command_end EXECUTE FUNCTION log_ddl();

-- Row Level Security (RLS) Policies

-- Enable RLS on all tables
ALTER TABLE parties ENABLE ROW LEVEL SECURITY;
ALTER TABLE policies ENABLE ROW LEVEL SECURITY;
ALTER TABLE claims ENABLE ROW LEVEL SECURITY;
ALTER TABLE value_movements ENABLE ROW LEVEL SECURITY;

-- Example RLS policy: Parties can only see their own data
CREATE POLICY party_isolation ON parties
    FOR ALL
    USING (id = current_setting('app.current_party_id')::UUID);

-- Example RLS policy: Assessors can see claims assigned to them
CREATE POLICY claim_access ON claims
    FOR SELECT
    USING (
        assessor_party_id = current_setting('app.current_party_id')::UUID
        OR EXISTS (
            SELECT 1 FROM entitlements e
            WHERE e.party_id = current_setting('app.current_party_id')::UUID
            AND e.resource_type = 'claim'
            AND (e.resource_id = claims.id OR e.resource_id IS NULL)
            AND e.action = 'view'
        )
    );

-- Peer caching metadata
CREATE TABLE peer_cache_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_name VARCHAR(100) NOT NULL,
    record_id UUID NOT NULL,
    cache_version INTEGER DEFAULT 1,
    last_modified TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    etag VARCHAR(64),
    expires_at TIMESTAMPTZ,
    cached_by_peers UUID[], -- Array of peer node IDs
    invalidation_required BOOLEAN DEFAULT FALSE
);

-- Columnar archival (for cold data)
CREATE TABLE archival_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_table VARCHAR(100) NOT NULL,
    archive_table VARCHAR(100) NOT NULL,
    date_from DATE NOT NULL,
    date_to DATE NOT NULL,
    record_count BIGINT,
    archive_location TEXT, -- S3 path or file path
    compressed_size_bytes BIGINT,
    original_size_bytes BIGINT,
    archived_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    restored_at TIMESTAMPTZ,
    verification_hash VARCHAR(64)
);

-- Data integrity checks (scheduled)
CREATE TABLE integrity_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    check_name VARCHAR(100) NOT NULL,
    check_type VARCHAR(50), -- balance_check, hash_verification, orphan_check
    last_run_at TIMESTAMPTZ,
    last_status VARCHAR(20), -- passed, failed
    last_error TEXT,
    check_query TEXT,
    schedule VARCHAR(50) -- cron expression
);

-- Insert standard integrity checks
INSERT INTO integrity_checks (check_name, check_type, check_query, schedule) VALUES
('Double Entry Balance', 'balance_check', 
 'SELECT movement_id FROM movement_legs GROUP BY movement_id HAVING SUM(CASE WHEN direction=''debit'' THEN amount ELSE -amount END) != 0', 
 '0 2 * * *'), -- Daily at 2 AM
('Immutable Hash Chain', 'hash_verification',
 'SELECT event_id FROM immutable_events WHERE event_hash != encode(digest(event_type || payload_hash || previous_hash, ''sha256''), ''hex'')',
 '0 3 * * *');

-- Function to verify entire ledger integrity
CREATE OR REPLACE FUNCTION verify_ledger_integrity() RETURNS TABLE (
    check_name VARCHAR,
    status VARCHAR,
    details TEXT
) AS $$
DECLARE
    v_check RECORD;
    v_result BOOLEAN;
    v_count INTEGER;
BEGIN
    FOR v_check IN SELECT * FROM integrity_checks WHERE check_type = 'balance_check' LOOP
        EXECUTE 'SELECT COUNT(*) FROM (' || v_check.check_query || ') t' INTO v_count;
        
        check_name := v_check.check_name;
        IF v_count = 0 THEN
            status := 'PASSED';
            details := 'No discrepancies found';
        ELSE
            status := 'FAILED';
            details := v_count || ' records failed validation';
        END IF;
        
        -- Update last run
        UPDATE integrity_checks 
        SET last_run_at = CURRENT_TIMESTAMP, 
            last_status = status,
            last_error = CASE WHEN status = 'FAILED' THEN details ELSE NULL END
        WHERE id = v_check.id;
        
        RETURN NEXT;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Performance monitoring
CREATE TABLE query_performance_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    query_hash VARCHAR(64),
    query_text TEXT,
    execution_time_ms INTEGER,
    rows_affected BIGINT,
    user_id UUID,
    ip_address INET,
    executed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Partition management automation
CREATE OR REPLACE FUNCTION create_monthly_partitions(
    p_table_name TEXT,
    p_start_date DATE,
    p_months INTEGER
) RETURNS VOID AS $$
DECLARE
    v_month_start DATE;
    v_month_end DATE;
    v_partition_name TEXT;
BEGIN
    FOR i IN 0..p_months-1 LOOP
        v_month_start := p_start_date + (i || ' months')::interval;
        v_month_end := v_month_start + '1 month'::interval;
        v_partition_name := p_table_name || '_' || to_char(v_month_start, 'YYYY_MM');
        
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES FROM (%L) TO (%L)',
            v_partition_name, p_table_name, v_month_start, v_month_end
        );
    END LOOP;
END;
$$ LANGUAGE plpgsql;
