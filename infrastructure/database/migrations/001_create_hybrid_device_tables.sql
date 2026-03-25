-- Migration: Create Hybrid Device Tables
-- Version: 001
-- Description: Creates tables for hybrid polymorphic device storage system

-- Create devices_hybrid table with JSONB for flexible category-specific data
CREATE TABLE IF NOT EXISTS devices_hybrid (
    -- Base fields
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Core device fields
    imei VARCHAR(15) NOT NULL,
    serial_number VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    owner_id UUID NOT NULL,
    
    -- JSONB fields for flexibility
    specifications JSONB NOT NULL,
    insurance_data JSONB,
    risk_data JSONB,
    
    -- Common fields
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    condition VARCHAR(20) NOT NULL DEFAULT 'good',
    grade VARCHAR(20) NOT NULL DEFAULT 'B',
    
    -- Indexed fields for performance
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    market_value DECIMAL(10,2),
    risk_score DECIMAL(5,2),
    
    -- Constraints
    CONSTRAINT unique_imei UNIQUE (imei),
    CONSTRAINT unique_serial_number UNIQUE (serial_number),
    CONSTRAINT valid_market_value CHECK (market_value >= 0),
    CONSTRAINT valid_risk_score CHECK (risk_score >= 0 AND risk_score <= 100)
);

-- Create indexes for better query performance
CREATE INDEX idx_devices_hybrid_category ON devices_hybrid(category) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_owner_id ON devices_hybrid(owner_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_status ON devices_hybrid(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_manufacturer ON devices_hybrid(manufacturer) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_model ON devices_hybrid(model) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_market_value ON devices_hybrid(market_value) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_risk_score ON devices_hybrid(risk_score) WHERE deleted_at IS NULL AND risk_score IS NOT NULL;
CREATE INDEX idx_devices_hybrid_deleted_at ON devices_hybrid(deleted_at);

-- Create GIN indexes for JSONB columns for efficient queries
CREATE INDEX idx_devices_hybrid_specifications_gin ON devices_hybrid USING gin(specifications);
CREATE INDEX idx_devices_hybrid_insurance_data_gin ON devices_hybrid USING gin(insurance_data) WHERE insurance_data IS NOT NULL;
CREATE INDEX idx_devices_hybrid_risk_data_gin ON devices_hybrid USING gin(risk_data) WHERE risk_data IS NOT NULL;

-- Create composite indexes for common queries
CREATE INDEX idx_devices_hybrid_category_status ON devices_hybrid(category, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_owner_category ON devices_hybrid(owner_id, category) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_hybrid_category_value_range ON devices_hybrid(category, market_value) WHERE deleted_at IS NULL;

-- Create device_categories table (already created in previous implementation)
CREATE TABLE IF NOT EXISTS device_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL,
    specifications JSONB NOT NULL,
    imei VARCHAR(15) UNIQUE NOT NULL,
    serial_number VARCHAR(50) UNIQUE NOT NULL,
    owner_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    insurance_profile JSONB,
    risk_assessment JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Indexed fields
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    market_value DECIMAL(10,2),
    risk_score DECIMAL(5,2),
    eligibility_status VARCHAR(20)
);

-- Create indexes for device_categories if not exists
CREATE INDEX IF NOT EXISTS idx_device_categories_category ON device_categories(category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_device_categories_owner_id ON device_categories(owner_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_device_categories_status ON device_categories(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_device_categories_specifications_gin ON device_categories USING gin(specifications);

-- Create trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for automatic timestamp updates
DROP TRIGGER IF EXISTS update_devices_hybrid_updated_at ON devices_hybrid;
CREATE TRIGGER update_devices_hybrid_updated_at 
    BEFORE UPDATE ON devices_hybrid 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_device_categories_updated_at ON device_categories;
CREATE TRIGGER update_device_categories_updated_at 
    BEFORE UPDATE ON device_categories 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to extract and update indexed fields from JSONB
CREATE OR REPLACE FUNCTION update_device_indexed_fields()
RETURNS TRIGGER AS $$
BEGIN
    -- Extract manufacturer from specifications
    NEW.manufacturer = NEW.specifications->>'manufacturer';
    
    -- Extract model from specifications
    NEW.model = NEW.specifications->>'model';
    
    -- Extract market_value from specifications
    IF NEW.specifications->>'market_value' IS NOT NULL THEN
        NEW.market_value = (NEW.specifications->>'market_value')::DECIMAL(10,2);
    END IF;
    
    -- Extract risk_score from risk_data if exists
    IF NEW.risk_data IS NOT NULL AND NEW.risk_data->>'score' IS NOT NULL THEN
        NEW.risk_score = (NEW.risk_data->>'score')::DECIMAL(5,2);
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically extract indexed fields
DROP TRIGGER IF EXISTS extract_devices_hybrid_indexed_fields ON devices_hybrid;
CREATE TRIGGER extract_devices_hybrid_indexed_fields 
    BEFORE INSERT OR UPDATE ON devices_hybrid 
    FOR EACH ROW 
    EXECUTE FUNCTION update_device_indexed_fields();

-- Create view for easy querying of device insurance status
CREATE OR REPLACE VIEW device_insurance_status AS
SELECT 
    d.id,
    d.category,
    d.imei,
    d.serial_number,
    d.manufacturer,
    d.model,
    d.market_value,
    d.risk_score,
    d.status,
    d.condition,
    d.insurance_data->>'eligibility_status' as eligibility_status,
    d.insurance_data->>'recommended_coverage' as recommended_coverage,
    d.risk_data->>'level' as risk_level,
    d.risk_data->>'requires_review' as requires_review,
    d.owner_id,
    d.created_at,
    d.updated_at
FROM devices_hybrid d
WHERE d.deleted_at IS NULL;

-- Create view for high-risk devices
CREATE OR REPLACE VIEW high_risk_devices AS
SELECT 
    d.*,
    d.risk_data->>'level' as risk_level,
    d.risk_data->'recommendations' as risk_recommendations
FROM devices_hybrid d
WHERE d.deleted_at IS NULL 
    AND d.risk_score > 70
    OR (d.risk_data->>'requires_review')::boolean = true;

-- Create view for devices needing risk reassessment
CREATE OR REPLACE VIEW devices_needing_reassessment AS
SELECT 
    d.id,
    d.category,
    d.imei,
    d.manufacturer,
    d.model,
    d.risk_data->>'valid_until' as risk_valid_until,
    d.owner_id
FROM devices_hybrid d
WHERE d.deleted_at IS NULL
    AND (
        d.risk_data IS NULL 
        OR (d.risk_data->>'valid_until')::timestamp < CURRENT_TIMESTAMP
    );

-- Create materialized view for device statistics by category
CREATE MATERIALIZED VIEW IF NOT EXISTS device_category_stats AS
SELECT 
    category,
    COUNT(*) as total_devices,
    COUNT(*) FILTER (WHERE status = 'active') as active_devices,
    COUNT(*) FILTER (WHERE insurance_data IS NOT NULL) as insured_devices,
    AVG(market_value) as avg_market_value,
    AVG(risk_score) as avg_risk_score,
    MIN(market_value) as min_market_value,
    MAX(market_value) as max_market_value,
    COUNT(DISTINCT manufacturer) as manufacturer_count,
    COUNT(DISTINCT model) as model_count
FROM devices_hybrid
WHERE deleted_at IS NULL
GROUP BY category;

-- Create index on materialized view
CREATE UNIQUE INDEX idx_device_category_stats_category 
    ON device_category_stats(category);

-- Create function to refresh materialized view
CREATE OR REPLACE FUNCTION refresh_device_category_stats()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY device_category_stats;
END;
$$ language 'plpgsql';

-- Grant permissions (adjust based on your user setup)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON devices_hybrid TO smartsure_app;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON device_categories TO smartsure_app;
-- GRANT SELECT ON device_insurance_status TO smartsure_app;
-- GRANT SELECT ON high_risk_devices TO smartsure_app;
-- GRANT SELECT ON devices_needing_reassessment TO smartsure_app;
-- GRANT SELECT ON device_category_stats TO smartsure_app;
-- GRANT EXECUTE ON FUNCTION refresh_device_category_stats() TO smartsure_app;
