package policy

import (
	"gorm.io/gorm"
)

// ============================================
// POLICY DATABASE MIGRATION
// ============================================

// MigratePolicyTables creates all policy-related tables with proper indexes and constraints
func MigratePolicyTables(db *gorm.DB) error {
	// Create main Policy table with all embedded structs
	// Note: Migration handled by main models package to avoid import cycles
	// if err := db.AutoMigrate(&Policy{}); err != nil {
	// 	return err
	// }

	// Create PolicyModification table
	if err := db.AutoMigrate(&PolicyModification{}); err != nil {
		return err
	}

	// Create supporting policy tables
	supportingTables := []interface{}{
		&PolicyBundle{},
		// PolicyEndorsement should be migrated from the models package separately to avoid circular import
		&PolicyRider{},
		&PolicyRenewal{},
		&PolicyPaymentSchedule{},
		&PolicyLimit{},
		&PolicyExclusion{},
		&PolicyUnderwriting{},
		&PolicyBenefit{},
		&PolicyDiscount{},
		&PolicyCommunicationPreference{},

		// Smartphone-specific features
		&PolicyCoverage{},
		&PolicyServiceOptions{},
		&PolicyInternationalCoverage{},
		&PolicyClaimLimits{},
		&PolicyLoyaltyProgram{},
		&PolicySmartFeatures{},
		&PolicyFamilyGroup{},
		&PolicyEnvironmental{},
		&PolicyCorporate{},
	}

	for _, table := range supportingTables {
		if err := db.AutoMigrate(table); err != nil {
			return err
		}
	}

	// Add custom indexes
	if err := AddPolicyIndexes(db); err != nil {
		return err
	}

	// Add constraints
	if err := AddPolicyConstraints(db); err != nil {
		return err
	}

	return nil
}

// AddPolicyIndexes adds performance indexes to policy tables
func AddPolicyIndexes(db *gorm.DB) error {
	// Composite indexes for common queries
	compositeIndexes := []struct {
		Table   string
		Name    string
		Columns []string
	}{
		// Policy search indexes
		{
			Table:   "policies",
			Name:    "idx_policy_customer_status",
			Columns: []string{"customer_id", "status"},
		},
		{
			Table:   "policies",
			Name:    "idx_policy_device_status",
			Columns: []string{"device_id", "status"},
		},
		{
			Table:   "policies",
			Name:    "idx_policy_effective_expiration",
			Columns: []string{"effective_date", "expiration_date"},
		},
		{
			Table:   "policies",
			Name:    "idx_policy_payment_status_next_billing",
			Columns: []string{"payment_status", "next_billing_date"},
		},
		{
			Table:   "policies",
			Name:    "idx_policy_underwriting_risk",
			Columns: []string{"underwriting_status", "risk_score"},
		},
		{
			Table:   "policies",
			Name:    "idx_policy_corporate_bundle",
			Columns: []string{"corporate_account_id", "bundle_id"},
		},

		// PolicyModification indexes
		{
			Table:   "policy_modifications",
			Name:    "idx_modification_policy_status",
			Columns: []string{"policy_id", "status"},
		},
		{
			Table:   "policy_modifications",
			Name:    "idx_modification_type_date",
			Columns: []string{"modification_type", "effective_date"},
		},
		{
			Table:   "policy_modifications",
			Name:    "idx_modification_approval_required",
			Columns: []string{"requires_approval", "status"},
		},
	}

	// Create composite indexes
	for _, idx := range compositeIndexes {
		sql := "CREATE INDEX IF NOT EXISTS " + idx.Name + " ON " + idx.Table + " ("
		for i, col := range idx.Columns {
			if i > 0 {
				sql += ", "
			}
			sql += col
		}
		sql += ")"

		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	// Add partial indexes for performance
	partialIndexes := []string{
		// Active policies only
		"CREATE INDEX IF NOT EXISTS idx_policy_active_only ON policies (customer_id, device_id) WHERE status = 'active'",

		// Pending renewals
		"CREATE INDEX IF NOT EXISTS idx_policy_pending_renewal ON policies (expiration_date, auto_renewal) WHERE status = 'active' AND auto_renewal = true",

		// High-risk policies
		"CREATE INDEX IF NOT EXISTS idx_policy_high_risk ON policies (risk_score, underwriting_status) WHERE risk_score > 70",

		// Overdue payments
		"CREATE INDEX IF NOT EXISTS idx_policy_overdue ON policies (payment_status, next_billing_date) WHERE payment_status != 'paid'",

		// Pending modifications
		"CREATE INDEX IF NOT EXISTS idx_modification_pending ON policy_modifications (policy_id, requested_date) WHERE status IN ('pending', 'in_review')",
	}

	for _, sql := range partialIndexes {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	// Add text search indexes (if PostgreSQL)
	if db.Dialector.Name() == "postgres" {
		textSearchIndexes := []string{
			// Full text search on policy number
			"CREATE INDEX IF NOT EXISTS idx_policy_number_search ON policies USING gin(to_tsvector('english', policy_number))",

			// Full text search on notes and metadata
			"CREATE INDEX IF NOT EXISTS idx_policy_notes_search ON policies USING gin(to_tsvector('english', notes || ' ' || internal_notes))",

			// JSONB indexes for efficient queries
			"CREATE INDEX IF NOT EXISTS idx_policy_tags ON policies USING gin(tags)",
			"CREATE INDEX IF NOT EXISTS idx_policy_custom_fields ON policies USING gin(custom_fields)",
			"CREATE INDEX IF NOT EXISTS idx_policy_coverage_limits ON policies USING gin(coverage_limits)",
		}

		for _, sql := range textSearchIndexes {
			if err := db.Exec(sql).Error; err != nil {
				// Skip if error (might not be PostgreSQL)
				continue
			}
		}
	}

	return nil
}

// AddPolicyConstraints adds database constraints for data integrity
func AddPolicyConstraints(db *gorm.DB) error {
	constraints := []string{
		// Policy constraints
		"ALTER TABLE policies ADD CONSTRAINT chk_policy_dates CHECK (expiration_date > effective_date)",
		"ALTER TABLE policies ADD CONSTRAINT chk_policy_coverage_positive CHECK (coverage_amount > 0)",
		"ALTER TABLE policies ADD CONSTRAINT chk_policy_premium_positive CHECK (final_premium_amount > 0)",
		"ALTER TABLE policies ADD CONSTRAINT chk_policy_risk_score CHECK (risk_score >= 0 AND risk_score <= 100)",
		"ALTER TABLE policies ADD CONSTRAINT chk_policy_deductible CHECK (deductible_amount >= 0)",
		"ALTER TABLE policies ADD CONSTRAINT chk_policy_coinsurance CHECK (coinsurance_percent >= 0 AND coinsurance_percent <= 100)",

		// PolicyModification constraints
		"ALTER TABLE policy_modifications ADD CONSTRAINT chk_modification_dates CHECK (effective_date >= requested_date)",
		"ALTER TABLE policy_modifications ADD CONSTRAINT chk_modification_impact CHECK (premium_impact_amount IS NOT NULL OR coverage_impact_amount IS NOT NULL)",

		// Foreign key constraints with CASCADE
		"ALTER TABLE policy_coverage ADD CONSTRAINT fk_policy_coverage_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_service_options ADD CONSTRAINT fk_policy_service_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_international_coverage ADD CONSTRAINT fk_policy_intl_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_claim_limits ADD CONSTRAINT fk_policy_limits_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_loyalty_programs ADD CONSTRAINT fk_policy_loyalty_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_smart_features ADD CONSTRAINT fk_policy_smart_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_family_groups ADD CONSTRAINT fk_policy_family_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_environmental ADD CONSTRAINT fk_policy_env_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
		"ALTER TABLE policy_corporate ADD CONSTRAINT fk_policy_corp_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE",
	}

	// Apply constraints
	for _, constraint := range constraints {
		if err := db.Exec(constraint).Error; err != nil {
			// Skip if constraint already exists or not supported
			continue
		}
	}

	// Add triggers for automatic field updates (if PostgreSQL)
	if db.Dialector.Name() == "postgres" {
		triggers := []string{
			// Update remaining_limit when claims are processed
			`CREATE OR REPLACE FUNCTION update_policy_remaining_limit()
			RETURNS TRIGGER AS $$
			BEGIN
				IF NEW.status = 'approved' AND OLD.status != 'approved' THEN
					UPDATE policies 
					SET remaining_amount = remaining_amount - NEW.approved_amount
					WHERE id = NEW.policy_id;
				END IF;
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql`,

			`CREATE TRIGGER trg_update_remaining_limit
			AFTER UPDATE ON claims
			FOR EACH ROW
			EXECUTE FUNCTION update_policy_remaining_limit()`,

			// Update risk scores when claims are filed
			`CREATE OR REPLACE FUNCTION update_policy_risk_on_claim()
			RETURNS TRIGGER AS $$
			BEGIN
				UPDATE policies 
				SET claim_frequency = claim_frequency + 1,
					loss_ratio = (SELECT SUM(approved_amount) FROM claims WHERE policy_id = NEW.policy_id) / final_premium_amount * 100
				WHERE id = NEW.policy_id;
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql`,

			`CREATE TRIGGER trg_update_risk_on_claim
			AFTER INSERT ON claims
			FOR EACH ROW
			EXECUTE FUNCTION update_policy_risk_on_claim()`,

			// Audit trail for policy modifications
			`CREATE OR REPLACE FUNCTION audit_policy_changes()
			RETURNS TRIGGER AS $$
			BEGIN
				INSERT INTO policy_modifications (
					policy_id, 
					modification_type, 
					field_modified,
					old_value,
					new_value,
					requested_by,
					requested_date,
					effective_date,
					status
				)
				SELECT 
					NEW.id,
					'system_adjustment',
					key,
					to_json(old_value),
					to_json(new_value),
					NEW.modified_by,
					NOW(),
					NOW(),
					'completed'
				FROM jsonb_each(to_jsonb(OLD)) AS old_table(key, old_value)
				JOIN jsonb_each(to_jsonb(NEW)) AS new_table(key, new_value) 
					ON old_table.key = new_table.key
				WHERE old_value IS DISTINCT FROM new_value
					AND key NOT IN ('updated_at', 'modified_by');
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql`,

			`CREATE TRIGGER trg_audit_policy_changes
			AFTER UPDATE ON policies
			FOR EACH ROW
			EXECUTE FUNCTION audit_policy_changes()`,
		}

		for _, trigger := range triggers {
			if err := db.Exec(trigger).Error; err != nil {
				// Skip if trigger creation fails
				continue
			}
		}
	}

	return nil
}

// CreatePolicyViews creates useful database views for reporting
func CreatePolicyViews(db *gorm.DB) error {
	views := []string{
		// Active policies summary view
		`CREATE OR REPLACE VIEW v_active_policies AS
		SELECT 
			p.id,
			p.policy_number,
			p.customer_id,
			u.name as customer_name,
			p.device_id,
			d.brand || ' ' || d.model as device_name,
			p.coverage_amount,
			p.final_premium_amount as premium,
			p.effective_date,
			p.expiration_date,
			p.risk_score,
			p.payment_status,
			p.next_billing_date
		FROM policies p
		LEFT JOIN users u ON p.customer_id = u.id
		LEFT JOIN devices d ON p.device_id = d.id
		WHERE p.status = 'active'
			AND p.deleted_at IS NULL`,

		// Policies needing renewal view
		`CREATE OR REPLACE VIEW v_policies_needing_renewal AS
		SELECT 
			p.*,
			DATE_PART('day', p.expiration_date - CURRENT_DATE) as days_until_expiry
		FROM policies p
		WHERE p.status = 'active'
			AND p.auto_renewal = true
			AND p.expiration_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '60 days'
			AND p.deleted_at IS NULL
		ORDER BY p.expiration_date`,

		// High-risk policies view
		`CREATE OR REPLACE VIEW v_high_risk_policies AS
		SELECT 
			p.id,
			p.policy_number,
			p.risk_score,
			p.fraud_risk_score,
			p.loss_ratio,
			p.claim_frequency,
			p.final_premium_amount as premium,
			p.coverage_amount,
			COUNT(c.id) as total_claims,
			SUM(c.approved_amount) as total_claim_amount
		FROM policies p
		LEFT JOIN claims c ON p.id = c.policy_id AND c.status IN ('approved', 'paid')
		WHERE p.risk_score > 70
			OR p.fraud_risk_score > 60
			OR p.loss_ratio > 80
		GROUP BY p.id`,

		// Policy financial summary view
		`CREATE OR REPLACE VIEW v_policy_financial_summary AS
		SELECT 
			DATE_TRUNC('month', p.created_at) as month,
			COUNT(p.id) as new_policies,
			SUM(p.final_premium_amount) as total_premium,
			SUM(p.coverage_amount) as total_coverage,
			AVG(p.risk_score) as avg_risk_score,
			COUNT(DISTINCT p.customer_id) as unique_customers
		FROM policies p
		WHERE p.status IN ('active', 'expired')
			AND p.deleted_at IS NULL
		GROUP BY DATE_TRUNC('month', p.created_at)
		ORDER BY month DESC`,
	}

	// Create views
	for _, view := range views {
		if err := db.Exec(view).Error; err != nil {
			// Skip if view creation fails
			continue
		}
	}

	return nil
}
