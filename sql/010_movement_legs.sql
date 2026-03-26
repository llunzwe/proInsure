-- 010_movement_legs.sql

-- Individual double-entry legs
CREATE TABLE movement_legs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movement_id UUID NOT NULL REFERENCES value_movements(id) ON DELETE CASCADE,
    container_id UUID NOT NULL REFERENCES value_containers(id),
    
    -- Direction and amount
    direction VARCHAR(6) NOT NULL CHECK (direction IN ('debit', 'credit')),
    amount DECIMAL(28,8) NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Running balance snapshot (for audit trail)
    container_balance_after DECIMAL(28,8),
    
    -- Description
    description TEXT,
    line_number INTEGER NOT NULL, -- For ordering
    
    -- Dimensions for reporting
    cost_center VARCHAR(50),
    project_code VARCHAR(50),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(movement_id, line_number)
);

-- Function to update container balances when movement is posted
CREATE OR REPLACE FUNCTION update_container_balance() RETURNS TRIGGER AS $$
DECLARE
    v_movement_status VARCHAR(20);
BEGIN
    -- Only process if parent movement is posted
    SELECT status INTO v_movement_status FROM value_movements WHERE id = NEW.movement_id;
    
    IF v_movement_status = 'posted' THEN
        IF NEW.direction = 'debit' THEN
            UPDATE value_containers 
            SET current_balance = current_balance + NEW.amount
            WHERE id = NEW.container_id AND normal_balance = 'debit';
            
            UPDATE value_containers 
            SET current_balance = current_balance - NEW.amount
            WHERE id = NEW.container_id AND normal_balance = 'credit';
        ELSE -- credit
            UPDATE value_containers 
            SET current_balance = current_balance - NEW.amount
            WHERE id = NEW.container_id AND normal_balance = 'debit';
            
            UPDATE value_containers 
            SET current_balance = current_balance + NEW.amount
            WHERE id = NEW.container_id AND normal_balance = 'credit';
        END IF;
        
        -- Update the balance snapshot on the leg
        SELECT current_balance INTO NEW.container_balance_after 
        FROM value_containers WHERE id = NEW.container_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Actually, better to handle this via a trigger on value_movements status change
-- Simplified approach: legs are immutable once movement is posted

CREATE INDEX idx_legs_movement ON movement_legs(movement_id);
CREATE INDEX idx_legs_container ON movement_legs(container_id);
CREATE INDEX idx_legs_direction ON movement_legs(direction);
