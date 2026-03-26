-- 024_reference_data.sql

-- Static reference data with temporal validity
CREATE TABLE reference_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    category VARCHAR(50) NOT NULL CHECK (category IN (
        'country', 'currency', 'language', 'claim_type', 'device_brand', 'device_model',
        'bank_code', 'payment_method', 'id_document_type', 'repair_type', 'rejection_reason'
    )),
    
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Hierarchical
    parent_code VARCHAR(50) REFERENCES reference_data(code),
    hierarchy_path ltree, -- For tree structures
    
    -- Metadata
    metadata JSONB, -- Flexible attributes per category
    
    -- Temporal validity
    valid_from DATE NOT NULL DEFAULT CURRENT_DATE,
    valid_to DATE NOT NULL DEFAULT '9999-12-31',
    is_current BOOLEAN DEFAULT TRUE,
    
    -- Internationalization
    locale VARCHAR(10) DEFAULT 'en',
    localized_names JSONB, -- {"es": "Nombre", "fr": "Nom"}
    
    -- Ordering
    sort_order INTEGER DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(category, code, valid_from, locale)
);

-- Insert core reference data
INSERT INTO reference_data (category, code, name, description) VALUES
('claim_type', 'THEFT', 'Theft', 'Device stolen'),
('claim_type', 'DAMAGE', 'Accidental Damage', 'Physical damage to device'),
('claim_type', 'LOSS', 'Loss', 'Device lost'),
('claim_type', 'WATER', 'Water Damage', 'Liquid damage'),
('claim_type', 'SCREEN', 'Screen Damage', 'Broken screen only'),
('currency', 'USD', 'US Dollar', 'United States Dollar'),
('currency', 'EUR', 'Euro', 'European Euro'),
('currency', 'GBP', 'British Pound', 'UK Pound Sterling');

-- Create indexes for ltree if using hierarchical data
CREATE INDEX idx_reference_path ON reference_data USING gist(hierarchy_path);
CREATE INDEX idx_reference_category ON reference_data(category, is_current);
