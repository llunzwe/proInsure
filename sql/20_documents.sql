-- 020_documents.sql

-- Document vault with cryptographic integrity
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Classification
    document_type VARCHAR(50) NOT NULL CHECK (document_type IN (
        'policy_pdf', 'policy_schedule', 'claim_photo', 'damage_photo',
        'police_report', 'repair_invoice', 'medical_report', 'id_document',
        'proof_of_address', 'bank_statement', 'contract', 'correspondence',
        'assessment_report', 'settlement_proof'
    )),
    
    -- Storage
    storage_location TEXT NOT NULL, -- S3 URI, IPFS hash, file path
    storage_backend VARCHAR(20) DEFAULT 's3' CHECK (storage_backend IN ('s3', 'ipfs', 'filesystem', 'azure', 'gcs')),
    
    -- Cryptographic integrity
    content_hash VARCHAR(64) NOT NULL, -- SHA-256 of file content
    encryption_status VARCHAR(20) DEFAULT 'unencrypted' CHECK (encryption_status IN ('unencrypted', 'aes256', 'pgp')),
    encryption_key_id VARCHAR(100), -- Reference to KMS key
    
    -- Metadata
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    page_count INTEGER,
    
    -- Linking
    linked_entity_type VARCHAR(50) CHECK (linked_entity_type IN ('policy', 'claim', 'device', 'party', 'payment')),
    linked_entity_id UUID,
    
    -- Upload tracking
    uploaded_by_party_id UUID REFERENCES parties(id),
    uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    upload_ip_address INET,
    
    -- Retention management
    retention_date DATE, -- When to delete
    legal_hold BOOLEAN DEFAULT FALSE, -- Override retention
    legal_hold_reason TEXT,
    
    -- OCR and processing
    ocr_text TEXT,
    ocr_confidence DECIMAL(5,2),
    processed_metadata JSONB, -- Extracted structured data
    
    -- Immutable hash
    immutable_hash VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES parties(id)
);

-- Document access log
CREATE TABLE document_access_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id),
    accessed_by UUID REFERENCES parties(id),
    accessed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    access_type VARCHAR(20) CHECK (access_type IN ('view', 'download', 'share', 'print')),
    ip_address INET,
    success BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_documents_entity ON documents(linked_entity_type, linked_entity_id);
CREATE INDEX idx_documents_type ON documents(document_type);
CREATE INDEX idx_documents_hash ON documents(content_hash);
CREATE INDEX idx_documents_retention ON documents(retention_date) WHERE legal_hold = FALSE;
CREATE INDEX idx_documents_uploaded ON documents(uploaded_at);
