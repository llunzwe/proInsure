-- 014_claim_evidence.sql

-- Cryptographic evidence storage
CREATE TABLE claim_evidence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    document_id UUID NOT NULL, -- Reference to documents table
    
    -- Evidence specifics
    evidence_type VARCHAR(50) NOT NULL CHECK (evidence_type IN (
        'photo_device', 'photo_damage', 'photo_serial', 'photo_imei',
        'police_report', 'repair_estimate', 'purchase_receipt', 
        'proof_of_ownership', 'id_document', 'video', 'audio'
    )),
    
    -- Verification
    is_verified BOOLEAN DEFAULT FALSE,
    verified_by UUID REFERENCES parties(id),
    verified_at TIMESTAMPTZ,
    verification_method VARCHAR(100), -- manual, ai_vision, blockchain_notarization
    
    -- For photos: EXIF data extraction
    exif_data JSONB,
    geolocation_lat DECIMAL(10,8),
    geolocation_lon DECIMAL(11,8),
    captured_at TIMESTAMPTZ,
    
    -- Cryptographic integrity
    content_hash VARCHAR(64) NOT NULL, -- SHA-256 of file
    blockchain_anchor VARCHAR(256), -- Transaction hash if anchored
    
    -- AI analysis results
    ai_analysis_result JSONB, -- Object detection, damage assessment, etc.
    ai_confidence_score DECIMAL(5,2),
    
    uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    uploaded_by UUID REFERENCES parties(id),
    
    -- Retention
    retention_until DATE,
    is_deleted BOOLEAN DEFAULT FALSE
);

-- Evidence chain of custody
CREATE TABLE evidence_custody_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    evidence_id UUID NOT NULL REFERENCES claim_evidence(id),
    action VARCHAR(50) NOT NULL, -- uploaded, viewed, downloaded, verified, deleted
    performed_by UUID REFERENCES parties(id),
    performed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    digital_signature BYTEA
);

CREATE INDEX idx_evidence_claim ON claim_evidence(claim_id);
CREATE INDEX idx_evidence_type ON claim_evidence(evidence_type);
CREATE INDEX idx_evidence_hash ON claim_evidence(content_hash);
