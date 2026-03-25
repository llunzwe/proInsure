package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// SimpleAutoMigrate creates tables without foreign keys for initial setup
func SimpleAutoMigrate(db *gorm.DB) error {
	// Enable UUID extension for PostgreSQL
	if db.Dialector.Name() == "postgres" {
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	}

	fmt.Println("Creating tables without foreign key constraints...")

	// Create basic tables using raw SQL to avoid GORM's automatic foreign key detection
	tables := []string{
		// Core tables
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ,
			created_by UUID,
			updated_by UUID,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			phone_number TEXT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			date_of_birth TIMESTAMPTZ,
			gender VARCHAR(10),
			nationality TEXT,
			country TEXT,
			city TEXT,
			address TEXT,
			postal_code TEXT,
			user_type VARCHAR(20) NOT NULL DEFAULT 'customer',
			user_role VARCHAR(20) NOT NULL DEFAULT 'customer',
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			kyc_status VARCHAR(20) NOT NULL DEFAULT 'not_started',
			kyc_level VARCHAR(20) NOT NULL DEFAULT 'basic',
			risk_score DECIMAL DEFAULT 50.0,
			credit_score BIGINT,
			profile_picture TEXT,
			two_factor_enabled BOOLEAN DEFAULT false,
			last_login_at TIMESTAMPTZ,
			email_verified BOOLEAN DEFAULT false,
			phone_verified BOOLEAN DEFAULT false,
			referral_code TEXT,
			referred_by UUID,
			notification_prefs JSON,
			marketing_consent BOOLEAN DEFAULT false,
			language_preference TEXT DEFAULT 'en',
			timezone TEXT DEFAULT 'UTC',
			suspended_at TIMESTAMPTZ,
			suspension_reason TEXT,
			login_attempts INTEGER DEFAULT 0,
			locked_until TIMESTAMPTZ,
			metadata JSON
		)`,

		`CREATE TABLE IF NOT EXISTS devices (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ,
			owner_id UUID,
			imei TEXT NOT NULL UNIQUE,
			imei2 TEXT,
			serial_number TEXT UNIQUE,
			brand TEXT NOT NULL,
			model TEXT NOT NULL,
			color TEXT,
			storage_capacity INTEGER,
			purchase_date TIMESTAMPTZ,
			purchase_price DECIMAL,
			status VARCHAR(20) DEFAULT 'active',
			verification_status VARCHAR(20) DEFAULT 'unverified',
			network_operator TEXT,
			network_status TEXT,
			battery_health INTEGER,
			screen_condition TEXT,
			body_condition TEXT,
			functional_issues TEXT,
			accessories TEXT,
			original_box BOOLEAN DEFAULT false,
			original_receipt BOOLEAN DEFAULT false,
			market_value DECIMAL,
			last_inspection TIMESTAMPTZ,
			inspection_notes TEXT,
			metadata JSON
		)`,

		`CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT,
			category TEXT,
			type TEXT,
			status TEXT,
			version TEXT,
			coverage_amount DECIMAL,
			deductible_amount DECIMAL,
			max_claims_per_year INTEGER,
			waiting_period_days INTEGER,
			base_premium DECIMAL,
			premium_frequency TEXT,
			currency TEXT,
			features JSON,
			benefits JSON,
			exclusions JSON,
			terms_conditions TEXT,
			min_device_age INTEGER,
			max_device_age INTEGER,
			min_device_value DECIMAL,
			max_device_value DECIMAL,
			eligible_brands JSON,
			eligible_models JSON,
			marketing_name TEXT,
			short_description TEXT,
			image_url TEXT,
			brochure_url TEXT,
			launch_date TIMESTAMPTZ,
			end_date TIMESTAMPTZ,
			is_promoted BOOLEAN,
			display_order INTEGER,
			tags JSON,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)`,

		`CREATE TABLE IF NOT EXISTS policies (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ,
			created_by UUID,
			updated_by UUID,
			policy_number TEXT NOT NULL UNIQUE,
			customer_id UUID NOT NULL,
			device_id UUID NOT NULL,
			product_id UUID NOT NULL,
			quote_id UUID,
			status VARCHAR(20) NOT NULL DEFAULT 'draft',
			effective_date TIMESTAMPTZ NOT NULL,
			expiration_date TIMESTAMPTZ NOT NULL,
			renewal_date TIMESTAMPTZ,
			coverage_amount DECIMAL NOT NULL,
			deductible_amount DECIMAL DEFAULT 0,
			premium_amount DECIMAL NOT NULL,
			premium_frequency VARCHAR(20) DEFAULT 'monthly',
			payment_method VARCHAR(50),
			auto_renewal BOOLEAN DEFAULT false,
			cancellation_date TIMESTAMPTZ,
			cancellation_reason TEXT,
			risk_score DECIMAL DEFAULT 0,
			claims_count INTEGER DEFAULT 0,
			last_claim_date TIMESTAMPTZ,
			total_claimed_amount DECIMAL DEFAULT 0,
			notes TEXT,
			terms_accepted BOOLEAN DEFAULT false,
			terms_accepted_at TIMESTAMPTZ,
			signed_document_url TEXT,
			renewal_count INTEGER DEFAULT 0,
			discount_applied DECIMAL DEFAULT 0,
			referral_code TEXT,
			agent_id UUID,
			underwriting_decision_id UUID,
			corporate_policy_id UUID,
			family_plan_id UUID,
			metadata JSON
		)`,

		`CREATE TABLE IF NOT EXISTS quotes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			quote_number TEXT NOT NULL,
			user_id UUID NOT NULL,
			device_id UUID NOT NULL,
			product_id UUID NOT NULL,
			status TEXT,
			type TEXT,
			channel TEXT,
			device_make TEXT,
			device_model TEXT,
			device_value DECIMAL,
			device_age BIGINT,
			imei TEXT,
			coverage_type TEXT,
			coverage_amount DECIMAL,
			deductible_amount DECIMAL,
			policy_term BIGINT,
			base_premium DECIMAL,
			tax_amount DECIMAL,
			discount_amount DECIMAL,
			total_premium DECIMAL,
			currency TEXT,
			payment_frequency TEXT,
			risk_score DECIMAL,
			risk_category TEXT,
			underwriting_notes TEXT,
			valid_from TIMESTAMPTZ,
			valid_until TIMESTAMPTZ,
			converted_at TIMESTAMPTZ,
			converted_policy_id UUID,
			agent_id UUID,
			referral_code TEXT,
			promo_code TEXT,
			source TEXT,
			campaign TEXT,
			ip_address TEXT,
			user_agent TEXT,
			session_id TEXT,
			metadata JSON,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)`,

		`CREATE TABLE IF NOT EXISTS claims (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ,
			created_by UUID,
			updated_by UUID,
			claim_number TEXT NOT NULL UNIQUE,
			policy_id UUID NOT NULL,
			customer_id UUID NOT NULL,
			device_id UUID NOT NULL,
			claim_type VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'submitted',
			priority VARCHAR(20) DEFAULT 'medium',
			incident_date TIMESTAMPTZ NOT NULL,
			reported_date TIMESTAMPTZ,
			description TEXT NOT NULL,
			location TEXT,
			claimed_amount DECIMAL NOT NULL,
			approved_amount DECIMAL,
			deductible_amount DECIMAL,
			currency TEXT DEFAULT 'USD',
			assigned_to UUID,
			processor_notes TEXT,
			fraud_score DECIMAL DEFAULT 0,
			is_urgent BOOLEAN DEFAULT false,
			requires_investigation BOOLEAN DEFAULT false,
			estimated_settlement TIMESTAMPTZ,
			settled_at TIMESTAMPTZ,
			micro_insurance_id UUID,
			repair_booking_id UUID,
			replacement_order_id UUID,
			theft_report_id UUID,
			police_report_number TEXT,
			witness_contact TEXT,
			damage_assessment JSON,
			repair_estimate DECIMAL,
			total_loss BOOLEAN DEFAULT false,
			rejection_reason TEXT,
			appeal_status TEXT,
			appeal_date TIMESTAMPTZ,
			investigation_notes TEXT,
			payout_method TEXT,
			payout_date TIMESTAMPTZ
		)`,
	}

	// Execute table creation
	for _, sql := range tables {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_devices_owner_id ON devices(owner_id)",
		"CREATE INDEX IF NOT EXISTS idx_devices_imei ON devices(imei)",
		"CREATE INDEX IF NOT EXISTS idx_devices_deleted_at ON devices(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_policies_policy_number ON policies(policy_number)",
		"CREATE INDEX IF NOT EXISTS idx_policies_customer_id ON policies(customer_id)",
		"CREATE INDEX IF NOT EXISTS idx_policies_device_id ON policies(device_id)",
		"CREATE INDEX IF NOT EXISTS idx_policies_deleted_at ON policies(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_quotes_quote_number ON quotes(quote_number)",
		"CREATE INDEX IF NOT EXISTS idx_quotes_user_id ON quotes(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_quotes_deleted_at ON quotes(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_claims_claim_number ON claims(claim_number)",
		"CREATE INDEX IF NOT EXISTS idx_claims_policy_id ON claims(policy_id)",
		"CREATE INDEX IF NOT EXISTS idx_claims_customer_id ON claims(customer_id)",
		"CREATE INDEX IF NOT EXISTS idx_claims_deleted_at ON claims(deleted_at)",
	}

	for _, sql := range indexes {
		db.Exec(sql) // Ignore errors for indexes
	}

	fmt.Println("✅ Essential tables created successfully!")
	fmt.Println("Tables created: users, devices, products, policies, quotes, claims")
	fmt.Println("")
	fmt.Println("Note: Foreign key constraints have been omitted to avoid circular dependencies.")
	fmt.Println("You can add foreign keys manually later if needed.")

	return nil
}
