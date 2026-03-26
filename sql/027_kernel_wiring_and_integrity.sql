-- =============================================================================
-- FILENAME: 027_kernel_wiring_and_integrity.sql
-- DESCRIPTION: System integrity controls, RLS policies, and audit infrastructure
-- VERSION: 1.0.0
-- DEPENDENCIES: All preceding schema files
-- =============================================================================
-- SECURITY CLASSIFICATION: CONFIDENTIAL
-- DATA SENSITIVITY: System-level security configuration
-- =============================================================================
-- ISO/IEC COMPLIANCE:
--   - ISO/IEC 27001:2013 - Access control (A.9), Audit logging (A.12.4)
--   - ISO/IEC 27002:2022 - Security controls implementation
--   - SOX 404 - Internal controls over financial reporting
--   - GDPR Article 32 - Security of processing
-- =============================================================================
-- CHANGE LOG:
--   v1.0.0 (2026-03-26) - Initial release with comprehensive security controls
-- =============================================================================

-- -----------------------------------------------------------------------------
-- SECTION 1: SCHEMA AUDIT LOG
-- PURPOSE: Track all DDL changes for compliance and rollback
-- COMPLIANCE: SOX 404 - Change management, ISO 27001 A.12.1.2
-- -----------------------------------------------------------------------------

/**
 * TABLE: schema_audit_log
 * DESCRIPTION: Immutable audit trail for all database schema changes
 * 
 * SECURITY: Append-only, no UPDATE or DELETE allowed
 * RETENTION: 7 years (SOX requirement)
 */
CREATE TABLE schema_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Change classification
    action VARCHAR(50) NOT NULL,        -- CREATE, ALTER, DROP, COMMENT
    object_type VARCHAR(50) NOT NULL,   -- TABLE, INDEX, FUNCTION, POLICY, etc.
    object_name VARCHAR(200) NOT NULL,
    schema_name VARCHAR(100) DEFAULT 'public',
    
    -- DDL tracking
    ddl_command TEXT NOT NULL,          -- Full SQL command (masked for sensitive data)
    ddl_hash VARCHAR(64),               -- SHA-256 of command for integrity
    
    -- Context
    performed_by VARCHAR(100) NOT NULL DEFAULT CURRENT_USER,
    performed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_ip INET,
    application_name VARCHAR(100),
    
    -- Version control
    application_version VARCHAR(20),
    migration_script_version INTEGER REFERENCES schema_metadata(version),
    git_commit_hash VARCHAR(40),
    
    -- Approval (SOX 404)
    change_ticket_reference VARCHAR(50),    -- Jira/ServiceNow ticket
    approved_by VARCHAR(100),
    approval_timestamp TIMESTAMPTZ,
    
    -- Impact assessment
    affected_tables VARCHAR(100)[],
    estimated_impact VARCHAR(20) CHECK (estimated_impact IN ('low', 'medium', 'high', 'critical'))
);

COMMENT ON TABLE schema_audit_log IS 
'Immutable DDL audit trail per SOX 404 and ISO 27001 A.12.1.2. Retention: 7 years.';

-- Index for efficient querying
CREATE INDEX idx_schema_audit_time ON schema_audit_log(performed_at);
CREATE INDEX idx_schema_audit_object ON schema_audit_log(object_type, object_name);
CREATE INDEX idx_schema_audit_action ON schema_audit_log(action);

-- Prevent modifications to audit log
CREATE OR REPLACE FUNCTION prevent_audit_log_modification()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Schema audit log is immutable per SOX 404. Modifications not permitted.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_protect_schema_audit
    BEFORE UPDATE OR DELETE ON schema_audit_log
    FOR EACH ROW EXECUTE FUNCTION prevent_audit_log_modification();

-- -----------------------------------------------------------------------------
-- SECTION 2: DDL EVENT TRIGGER
-- PURPOSE: Automatic logging of all schema changes
-- SECURITY: Requires superuser or rds_superuser role
-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION log_ddl_event()
RETURNS event_trigger 
SECURITY DEFINER
AS $$
DECLARE
    v_object_name TEXT;
    v_schema_name TEXT;
BEGIN
    -- Extract object information from command tags
    SELECT object_identity INTO v_object_name
    FROM pg_event_trigger_ddl_commands()
    WHERE command_tag = tg_tag
    LIMIT 1;
    
    INSERT INTO schema_audit_log (
        action,
        object_type,
        object_name,
        schema_name,
        ddl_command,
        ddl_hash,
        performed_by,
        performed_at
    ) VALUES (
        tg_tag,
        tg_event,
        COALESCE(v_object_name, 'unknown'),
        current_schema(),
        current_query(),
        encode(digest(current_query(), 'sha256'), 'hex'),
        current_user,
        CURRENT_TIMESTAMP
    );
END;
$$ LANGUAGE plpgsql;

-- Note: Enable with superuser privileges
-- CREATE EVENT TRIGGER ddl_audit_trigger 
--     ON ddl_command_end 
--     EXECUTE FUNCTION log_ddl_event();

-- -----------------------------------------------------------------------------
-- SECTION 3: DATA INTEGRITY MONITORING
-- PURPOSE: Automated integrity checks and alerts
-- -----------------------------------------------------------------------------

CREATE TABLE integrity_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Check classification
    check_name VARCHAR(100) NOT NULL UNIQUE,
    check_type VARCHAR(50) NOT NULL CHECK (check_type IN ('balance', 'reconciliation', 'constraint', 'business_rule')),
    check_description TEXT,
    
    -- SQL to execute
    check_query TEXT NOT NULL,
    expected_result TEXT,               -- Expected result or pattern
    
    -- Schedule
    frequency VARCHAR(20) CHECK (frequency IN ('realtime', 'hourly', 'daily', 'weekly', 'monthly')),
    last_run_at TIMESTAMPTZ,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE integrity_checks IS 
'Automated data integrity check configuration.';

-- Integrity check results
CREATE TABLE integrity_check_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    check_id UUID NOT NULL REFERENCES integrity_checks(id),
    
    -- Execution details
    run_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    run_duration_ms INTEGER,
    
    -- Results
    status VARCHAR(20) CHECK (status IN ('passed', 'failed', 'error')),
    actual_result TEXT,
    records_checked INTEGER,
    violations_found INTEGER,
    
    -- Alerting
    alert_sent BOOLEAN DEFAULT FALSE,
    alert_sent_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    resolved_by UUID REFERENCES parties(id)
);

CREATE INDEX idx_integrity_results_check ON integrity_check_results(check_id);
CREATE INDEX idx_integrity_results_time ON integrity_check_results(run_at);
CREATE INDEX idx_integrity_results_status ON integrity_check_results(status) WHERE status != 'passed';

-- Sample integrity checks
INSERT INTO integrity_checks (check_name, check_type, check_description, check_query, frequency) VALUES
('double_entry_balance', 'balance', 'Verify all value_movements have balanced debits=credits',
 'SELECT COUNT(*) FROM value_movements WHERE total_debits != total_credits AND status = ''posted''', 'hourly'),

('client_money_reconciliation', 'reconciliation', 'Verify client money calculations match trust accounts',
 'SELECT COUNT(*) FROM client_money_calculations WHERE status != ''reconciled'' AND calculation_date = CURRENT_DATE', 'daily'),

('negative_balances', 'constraint', 'Check for negative balances in asset accounts',
 'SELECT COUNT(*) FROM value_containers WHERE type = ''ASSET'' AND current_balance < 0 AND status = ''active''', 'hourly'),

('unreconciled_settlements', 'reconciliation', 'Count settlements pending reconciliation > 7 days',
 'SELECT COUNT(*) FROM settlements WHERE status = ''final'' AND reconciliation_date IS NULL AND final_at < CURRENT_TIMESTAMP - INTERVAL ''7 days''', 'daily'),

('suspense_items_aging', 'business_rule', 'Flag suspense items > 30 days old',
 'SELECT COUNT(*) FROM suspense_items WHERE status = ''unmatched'' AND received_at < CURRENT_TIMESTAMP - INTERVAL ''30 days''', 'daily');

-- -----------------------------------------------------------------------------
-- SECTION 4: PERFORMANCE MONITORING
-- PURPOSE: Query performance and system health tracking
-- -----------------------------------------------------------------------------

CREATE TABLE performance_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Query identification
    query_hash VARCHAR(64),             -- Normalized query hash
    query_text TEXT,                    -- Truncated query text
    
    -- Execution metrics
    execution_time_ms INTEGER,
    rows_affected BIGINT,
    rows_returned BIGINT,
    
    -- Resource usage
    shared_blks_hit BIGINT,
    shared_blks_read BIGINT,
    temp_blks_written BIGINT,
    
    -- Context
    executed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    executed_by VARCHAR(100),
    application_name VARCHAR(100),
    client_ip INET
);

CREATE INDEX idx_perf_logs_time ON performance_logs(executed_at);
CREATE INDEX idx_perf_logs_duration ON performance_logs(execution_time_ms) WHERE execution_time_ms > 1000;
CREATE INDEX idx_perf_logs_hash ON performance_logs(query_hash);

-- -----------------------------------------------------------------------------
-- SECTION 5: DATA MASKING POLICIES
-- PURPOSE: Column-level data masking for non-privileged users
-- SECURITY: GDPR Article 32 - Pseudonymization
-- -----------------------------------------------------------------------------

CREATE TABLE data_masking_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Target
    table_name VARCHAR(100) NOT NULL,
    column_name VARCHAR(100) NOT NULL,
    
    -- Masking rule
    mask_type VARCHAR(50) NOT NULL CHECK (mask_type IN ('full', 'partial', 'regex', 'null', 'random')),
    mask_pattern VARCHAR(200),          -- e.g., 'XXX-XXX-####' for partial SSN
    
    -- Access control
    exempt_roles VARCHAR(100)[],        -- Roles that see unmasked data
    exempt_users VARCHAR(100)[],        -- Users that see unmasked data
    
    -- Conditions
    condition_column VARCHAR(100),      -- e.g., 'role'
    condition_value VARCHAR(100),       -- e.g., 'admin'
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(table_name, column_name)
);

COMMENT ON TABLE data_masking_policies IS 
'GDPR Article 32: Column-level data masking policies for pseudonymization.';

-- Sample masking policies
INSERT INTO data_masking_policies (table_name, column_name, mask_type, mask_pattern, exempt_roles) VALUES
('parties', 'email', 'partial', '***@***.com', ARRAY['admin', 'identity_verifier']),
('parties', 'phone', 'partial', '***-***-####', ARRAY['admin', 'identity_verifier']),
('party_identifiers', 'value_encrypted', 'full', NULL, ARRAY['identity_verifier']),
('devices', 'imei', 'partial', '##############**', ARRAY['device_admin', 'fraud_investigator']),
('settlements', 'payee_account_details_encrypted', 'full', NULL, ARRAY['finance_admin']);

-- -----------------------------------------------------------------------------
-- SECTION 6: PARTITION MANAGEMENT
-- PURPOSE: Automated partition creation for time-series tables
-- -----------------------------------------------------------------------------

CREATE TABLE partition_management (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    table_name VARCHAR(100) NOT NULL,
    partition_column VARCHAR(100) NOT NULL,
    partition_type VARCHAR(20) DEFAULT 'monthly' CHECK (partition_type IN ('daily', 'weekly', 'monthly', 'yearly')),
    
    -- Retention
    retention_period INTERVAL DEFAULT INTERVAL '7 years',
    archive_after INTERVAL DEFAULT INTERVAL '1 year',
    
    -- Status
    last_partition_created DATE,
    partitions_ahead INTEGER DEFAULT 3,  -- Create partitions N periods ahead
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Insert managed tables
INSERT INTO partition_management (table_name, partition_column, partition_type, retention_period) VALUES
('device_iot_telemetry', 'recorded_at', 'daily', INTERVAL '90 days'),
('premium_payments', 'created_at', 'monthly', INTERVAL '7 years'),
('value_movements', 'entry_date', 'monthly', INTERVAL '10 years'),
('financial_event_log', 'event_time', 'monthly', INTERVAL '10 years'),
('authorization_logs', 'timestamp', 'monthly', INTERVAL '7 years');

-- Partition creation function
CREATE OR REPLACE FUNCTION create_monthly_partitions(
    p_table_name TEXT,
    p_months_ahead INTEGER DEFAULT 3
) RETURNS INTEGER AS $$
DECLARE
    v_month DATE;
    v_partition_name TEXT;
    v_start_date DATE;
    v_end_date DATE;
    v_count INTEGER := 0;
BEGIN
    FOR i IN 0..p_months_ahead LOOP
        v_month := DATE_TRUNC('month', CURRENT_DATE + (i || ' months')::INTERVAL);
        v_partition_name := p_table_name || '_' || TO_CHAR(v_month, 'YYYY_MM');
        v_start_date := v_month;
        v_end_date := v_month + INTERVAL '1 month';
        
        -- Check if partition exists
        IF NOT EXISTS (
            SELECT 1 FROM pg_class c 
            JOIN pg_namespace n ON c.relnamespace = n.oid 
            WHERE n.nspname = 'public' AND c.relname = v_partition_name
        ) THEN
            EXECUTE format(
                'CREATE TABLE IF NOT EXISTS %I PARTITION OF %I FOR VALUES FROM (%L) TO (%L)',
                v_partition_name, p_table_name, v_start_date, v_end_date
            );
            v_count := v_count + 1;
        END IF;
    END LOOP;
    
    RETURN v_count;
END;
$$ LANGUAGE plpgsql;

-- -----------------------------------------------------------------------------
-- SECTION 7: ENCRYPTION KEY ROTATION TRACKING
-- PURPOSE: Audit trail for cryptographic key lifecycle
-- -----------------------------------------------------------------------------

CREATE TABLE encryption_key_rotation (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    key_id VARCHAR(100) NOT NULL,
    key_type VARCHAR(50) NOT NULL CHECK (key_type IN ('column_encryption', 'backup_encryption', 'document_encryption', 'tokenization')),
    
    -- Rotation details
    rotated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    rotated_by VARCHAR(100) DEFAULT CURRENT_USER,
    previous_key_id VARCHAR(100),
    new_key_id VARCHAR(100),
    
    -- Scope
    affected_tables VARCHAR(100)[],
    affected_rows BIGINT,
    
    -- Verification
    verification_status VARCHAR(20) DEFAULT 'pending' CHECK (verification_status IN ('pending', 'verified', 'failed')),
    verification_notes TEXT
);

COMMENT ON TABLE encryption_key_rotation IS 
'ISO 27001 A.10.1.2: Cryptographic key rotation audit trail.';

-- -----------------------------------------------------------------------------
-- SECTION 8: SYSTEM HEALTH CHECK VIEW
-- PURPOSE: Comprehensive system status overview
-- -----------------------------------------------------------------------------

CREATE OR REPLACE VIEW v_system_health AS
SELECT
    'schema_version' AS metric,
    MAX(version)::TEXT AS value,
    MAX(applied_at) AS last_updated
FROM schema_metadata
WHERE deployment_status = 'success'

UNION ALL

SELECT
    'unreconciled_settlements' AS metric,
    COUNT(*)::TEXT AS value,
    MAX(created_at) AS last_updated
FROM settlements
WHERE status = 'final' AND reconciliation_date IS NULL

UNION ALL

SELECT
    'open_suspense_items' AS metric,
    COUNT(*)::TEXT AS value,
    MAX(created_at) AS last_updated
FROM suspense_items
WHERE status = 'unmatched'

UNION ALL

SELECT
    'failed_integrity_checks' AS metric,
    COUNT(*)::TEXT AS value,
    MAX(run_at) AS last_updated
FROM integrity_check_results
WHERE status = 'failed' AND resolved_at IS NULL

UNION ALL

SELECT
    'client_money_breaches' AS metric,
    COUNT(*)::TEXT AS value,
    MAX(created_at) AS last_updated
FROM client_money_calculations
WHERE compliance_status = 'breach';

COMMENT ON VIEW v_system_health IS 
'Real-time system health indicators for monitoring dashboards.';

-- -----------------------------------------------------------------------------
-- SECTION 9: MAINTENANCE PROCEDURES
-- PURPOSE: Automated maintenance task configuration
-- -----------------------------------------------------------------------------

CREATE TABLE maintenance_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    task_name VARCHAR(100) NOT NULL UNIQUE,
    task_description TEXT,
    
    -- Schedule (cron format)
    schedule_cron VARCHAR(100) NOT NULL,
    
    -- SQL to execute
    task_sql TEXT NOT NULL,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    last_run_at TIMESTAMPTZ,
    last_run_status VARCHAR(20) CHECK (last_run_status IN ('success', 'failed')),
    last_run_message TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Standard maintenance tasks
INSERT INTO maintenance_tasks (task_name, task_description, schedule_cron, task_sql) VALUES
('update_device_values', 'Recalculate depreciated device values', '0 2 * * *',
 'UPDATE devices SET current_value = calculate_depreciated_value(id) WHERE depreciation_rate > 0'),

('cleanup_old_sessions', 'Remove expired session records', '0 3 * * *',
 'DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP - INTERVAL ''30 days'''),

('archive_old_claims', 'Archive closed claims > 7 years', '0 4 1 * *',
 'SELECT archive_old_records(''claims'', ''closed_at'', INTERVAL ''7 years'')'),

('update_party_analytics', 'Refresh claim velocity scores', '0 */6 * * *',
 'SELECT refresh_party_analytics()'),

('create_future_partitions', 'Create upcoming table partitions', '0 1 * * 0',
 'SELECT create_monthly_partitions(''value_movements'', 3)');

-- =============================================================================
-- END OF FILE: 027_kernel_wiring_and_integrity.sql
-- =============================================================================
