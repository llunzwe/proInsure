-- 000_extensions.sql
-- Foundation extensions and smartphone insurance specific setups

-- Core PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "ltree";
CREATE EXTENSION IF NOT EXISTS "timescaledb";

-- Custom types for insurance domain
DO $$ BEGIN
    CREATE TYPE party_type AS ENUM ('individual', 'corporate', 'repair_shop', 'agent', 'insurer', 'assessor', 'oem');
    CREATE TYPE kyc_status AS ENUM ('pending', 'verified', 'rejected', 'expired');
    CREATE TYPE sanctions_status AS ENUM ('clear', 'potential_match', 'confirmed');
    CREATE TYPE device_category AS ENUM ('smartphone', 'tablet', 'laptop', 'smartwatch', 'wearable', 'accessory');
    CREATE TYPE device_status AS ENUM ('active', 'stolen', 'damaged', 'decommissioned', 'blacklisted');
    CREATE TYPE policy_status AS ENUM ('active', 'lapsed', 'cancelled', 'expired', 'pending');
    CREATE TYPE claim_type AS ENUM ('theft', 'damage', 'loss', 'repair', 'water_damage', 'screen_damage');
    CREATE TYPE claim_status AS ENUM ('submitted', 'under_review', 'assessment', 'approved', 'rejected', 'paid', 'closed', 'appealed');
    CREATE TYPE account_type AS ENUM ('ASSET', 'LIABILITY', 'EQUITY', 'INCOME', 'EXPENSE');
    CREATE TYPE reserve_type AS ENUM ('reported', 'ibnr', 'uer', 'unearned_premium');
    CREATE TYPE movement_type AS ENUM ('premium_payment', 'claim_payout', 'reserve_creation', 'reserve_release', 'refund', 'commission', 'salvage_recovery', 'replacement_cost');
    CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'overdue', 'written_off', 'failed');
    CREATE TYPE settlement_method AS ENUM ('bank_transfer', 'mobile_money', 'replacement_device', 'repair_voucher', 'cash');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- IMEI validation function
CREATE OR REPLACE FUNCTION validate_imei(imei TEXT) RETURNS BOOLEAN AS $$
DECLARE
    sum INTEGER := 0;
    i INTEGER;
    digit INTEGER;
    check_digit INTEGER;
    imei_clean TEXT;
BEGIN
    -- Remove any dashes or spaces
    imei_clean := regexp_replace(imei, '[-\s]', '', 'g');
    
    -- Check length (15 digits standard)
    IF length(imei_clean) != 15 OR imei_clean !~ '^\d{15}$' THEN
        RETURN FALSE;
    END IF;
    
    -- Luhn algorithm validation
    FOR i IN 1..14 LOOP
        digit := substring(imei_clean from i for 1)::INTEGER;
        IF i % 2 = 0 THEN
            digit := digit * 2;
            IF digit > 9 THEN
                digit := digit - 9;
            END IF;
        END IF;
        sum := sum + digit;
    END LOOP;
    
    check_digit := (10 - (sum % 10)) % 10;
    RETURN check_digit = substring(imei_clean from 15 for 1)::INTEGER;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Cryptographic hash generation for immutability
CREATE OR REPLACE FUNCTION generate_row_hash() RETURNS TRIGGER AS $$
BEGIN
    NEW.immutable_hash := encode(digest(
        NEW::text || CURRENT_TIMESTAMP::text || txid_current()::text, 
        'sha256'
    ), 'hex');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Bitemporal constraint trigger
CREATE OR REPLACE FUNCTION check_bitemporal_overlap() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_class 
        WHERE oid = TG_RELID
    ) THEN
        -- Check for overlapping valid_time periods for same entity
        IF EXISTS (
            SELECT 1 FROM pg_class c
            JOIN pg_attribute a ON c.oid = a.attrelid
            WHERE c.oid = TG_RELID
            AND a.attname = 'valid_to'
        ) THEN
            EXECUTE format(
                'SELECT 1 FROM %I.%I 
                 WHERE id = $1.id 
                 AND valid_from < $1.valid_to 
                 AND valid_to > $1.valid_from 
                 AND system_time = $1.system_time
                 AND id != COALESCE($1.id, ''00000000-0000-0000-0000-000000000000''::uuid)',
                TG_TABLE_SCHEMA, TG_TABLE_NAME
            ) USING NEW;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
