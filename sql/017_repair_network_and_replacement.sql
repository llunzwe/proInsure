-- 017_repair_network_and_replacement.sql

-- Repair shops
CREATE TABLE repair_shops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    party_id UUID NOT NULL REFERENCES parties(id),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    
    -- Authorization with certification tiers
    is_authorized BOOLEAN DEFAULT FALSE,
    authorization_level VARCHAR(50), -- Tier 1 (OEM Authorized), Tier 2 (Premium Independent), Tier 3 (Economy), Mobile Unit
    certification_tier VARCHAR(20) CHECK (certification_tier IN ('tier1_oem', 'tier2_premium', 'tier3_economy', 'mobile_unit')),
    
    -- OEM certifications
    apple_irp_certified BOOLEAN DEFAULT FALSE,
    samsung_asp_certified BOOLEAN DEFAULT FALSE,
    google_authorized BOOLEAN DEFAULT FALSE,
    other_oem_certifications JSONB,
    
    -- Quality metrics (ISO 9001)
    first_time_fix_rate DECIMAL(5,2), -- Percentage
    customer_satisfaction_score DECIMAL(3,2), -- Out of 5
    average_repair_time_hours DECIMAL(6,2),
    warranty_return_rate DECIMAL(5,2), -- Percentage of repairs requiring rework
    
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

-- Repair parts inventory tracking
CREATE TABLE repair_parts_inventory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repair_shop_id UUID NOT NULL REFERENCES repair_shops(id),
    
    -- Part details
    part_number VARCHAR(50) NOT NULL,
    part_name VARCHAR(100) NOT NULL,
    oem_part BOOLEAN DEFAULT TRUE, -- Genuine vs aftermarket
    compatible_brands VARCHAR(100)[],
    compatible_models VARCHAR(100)[],
    
    -- Serialized tracking
    serial_number VARCHAR(100),
    batch_number VARCHAR(50),
    
    -- Inventory
    quantity_on_hand INTEGER DEFAULT 0,
    quantity_reserved INTEGER DEFAULT 0,
    reorder_point INTEGER DEFAULT 5,
    
    -- Costs
    unit_cost DECIMAL(28,8),
    retail_price DECIMAL(28,8),
    
    -- Supplier
    supplier_party_id UUID REFERENCES parties(id),
    supplier_part_number VARCHAR(50),
    
    -- Status
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'low_stock', 'out_of_stock', 'discontinued')),
    
    last_restocked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(repair_shop_id, part_number, serial_number)
);

CREATE INDEX idx_parts_inventory_shop ON repair_parts_inventory(repair_shop_id);
CREATE INDEX idx_parts_inventory_status ON repair_parts_inventory(status);
CREATE INDEX idx_parts_inventory_part ON repair_parts_inventory(part_number);

-- Parts usage per repair
CREATE TABLE repair_parts_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repair_booking_id UUID NOT NULL REFERENCES repair_bookings(id),
    part_id UUID REFERENCES repair_parts_inventory(id),
    
    part_number VARCHAR(50) NOT NULL,
    part_name VARCHAR(100),
    serial_number VARCHAR(100), -- If serialized part
    
    quantity_used INTEGER NOT NULL DEFAULT 1,
    unit_cost DECIMAL(28,8),
    unit_price_charged DECIMAL(28,8),
    
    -- Core exchange
    old_part_returned BOOLEAN DEFAULT FALSE,
    old_part_serial_number VARCHAR(100),
    core_exchange_credit DECIMAL(28,8) DEFAULT 0,
    
    installed_at TIMESTAMPTZ,
    installed_by UUID REFERENCES parties(id)
);

CREATE INDEX idx_parts_usage_repair ON repair_parts_usage(repair_booking_id);

-- Repair warranty tracking
CREATE TABLE repair_warranties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repair_booking_id UUID NOT NULL REFERENCES repair_bookings(id),
    
    warranty_type VARCHAR(20) DEFAULT 'standard' CHECK (warranty_type IN ('standard', 'extended', 'oem_parts', 'aftermarket')),
    warranty_days INTEGER NOT NULL,
    warranty_start_date DATE NOT NULL,
    warranty_end_date DATE NOT NULL,
    
    -- Coverage
    parts_covered BOOLEAN DEFAULT TRUE,
    labor_covered BOOLEAN DEFAULT TRUE,
    
    -- Claims
    warranty_claim_made BOOLEAN DEFAULT FALSE,
    warranty_claim_id UUID REFERENCES claims(id),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_repair_warranties_booking ON repair_warranties(repair_booking_id);

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
