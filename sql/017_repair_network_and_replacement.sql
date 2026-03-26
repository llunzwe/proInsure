-- 017_repair_network_and_replacement.sql

-- Repair shops
CREATE TABLE repair_shops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    
    -- Authorization
    is_authorized BOOLEAN DEFAULT FALSE,
    authorization_level VARCHAR(50), -- basic, premium, oem_certified
    authorized_at TIMESTAMPTZ,
    authorized_until TIMESTAMPTZ,
    
    -- Capabilities
    supported_brands VARCHAR(100)[],
    supported_categories device_category[],
    max_claim_value DECIMAL(28,8), -- Authorization limit
    
    -- Pricing
    service_fee_schedule JSONB, -- Rates per repair type
    parts_markup_percentage DECIMAL(5,2) DEFAULT 0,
    
    -- Performance
    average_repair_days INTEGER,
    rating DECIMAL(3,2),
    active_repair_count INTEGER DEFAULT 0,
    
    valid_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMPTZ NOT NULL DEFAULT 'infinity'
);

-- Repair bookings
CREATE TABLE repair_bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    repair_shop_id UUID NOT NULL REFERENCES repair_shops(id),
    device_id UUID NOT NULL REFERENCES devices(id),
    
    -- Booking details
    booking_reference VARCHAR(50) UNIQUE,
    scheduled_date DATE,
    completed_date DATE,
    
    -- Status
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'in_progress', 'waiting_parts', 'completed', 'cancelled', 'returned')),
    
    -- Costs
    estimated_cost DECIMAL(28,8),
    final_cost DECIMAL(28,8),
    parts_cost DECIMAL(28,8),
    labor_cost DECIMAL(28,8),
    
    -- Details
    damage_description TEXT,
    repair_description TEXT,
    warranty_days INTEGER DEFAULT 90,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Replacement device inventory
CREATE TABLE replacement_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    imei VARCHAR(17) NOT NULL UNIQUE,
    model VARCHAR(100) NOT NULL,
    brand VARCHAR(100),
    color VARCHAR(50),
    storage_gb INTEGER,
    condition VARCHAR(20) CHECK (condition IN ('new', 'refurbished', 'like_new')),
    
    -- Costing
    cost_price DECIMAL(28,8),
    retail_value DECIMAL(28,8),
    currency CHAR(3) DEFAULT 'USD',
    
    -- Status
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'reserved', 'assigned', 'delivered', 'returned', 'written_off')),
    
    -- Assignment
    assigned_to_claim_id UUID REFERENCES claims(id),
    assigned_at TIMESTAMPTZ,
    
    -- Supplier
    supplier_party_id UUID REFERENCES parties(id),
    purchase_date DATE,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_replacement_imei CHECK (validate_imei(imei))
);

-- Temporary loaner devices
CREATE TABLE temporary_loans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id), -- The loaner device
    original_device_id UUID NOT NULL REFERENCES devices(id), -- Customer's device being repaired
    borrower_party_id UUID NOT NULL REFERENCES parties(id),
    claim_id UUID NOT NULL REFERENCES claims(id),
    
    loan_start TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expected_return_date DATE,
    actual_return_date TIMESTAMPTZ,
    
    condition_before TEXT,
    condition_after TEXT,
    
    deposit_amount DECIMAL(28,8) DEFAULT 0,
    deposit_returned BOOLEAN DEFAULT FALSE,
    
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'returned', 'overdue', 'damaged'))
);

-- Indexes
CREATE INDEX idx_repair_shops_party ON repair_shops(party_id);
CREATE INDEX idx_repair_bookings_claim ON repair_bookings(claim_id);
CREATE INDEX idx_repair_bookings_shop ON repair_bookings(repair_shop_id);
CREATE INDEX idx_replacement_status ON replacement_devices(status);
CREATE INDEX idx_loans_claim ON temporary_loans(claim_id);
CREATE INDEX idx_loans_active ON temporary_loans(status) WHERE status = 'active';
