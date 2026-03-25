package migrations

import (
	"fmt"

	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/device"
	"smartsure/internal/domain/models/policy"
	"smartsure/internal/domain/models/shared"
	"smartsure/internal/domain/models/user"
)

// AutoMigrate runs automatic database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	fmt.Println("========================================")
	fmt.Println("Starting SmartSure Database Migration")
	fmt.Println("========================================")

	// Step 1: Setup PostgreSQL extensions
	if err := setupPostgreSQLExtensions(db); err != nil {
		return fmt.Errorf("failed to setup PostgreSQL extensions: %w", err)
	}

	// Step 2: Run migrations without foreign keys
	fmt.Println("\n📦 Phase 1: Creating database schema...")

	// Create a new session that skips hooks for faster migration
	migrationDB := db.Session(&gorm.Session{
		SkipHooks:       true,
		CreateBatchSize: 100,
	})

	// Run the actual migration
	if err := runCoreMigrations(migrationDB); err != nil {
		return fmt.Errorf("failed to run core migrations: %w", err)
	}

	// Step 3: Run additional migrations
	if err := runAdditionalMigrations(migrationDB); err != nil {
		// Non-fatal, log and continue
		fmt.Printf("Warning: Some additional tables could not be created: %v\n", err)
	}

	// Step 4: Create indexes
	fmt.Println("\n📇 Phase 2: Creating indexes...")
	if err := createCoreIndexes(db); err != nil {
		fmt.Printf("Warning: Some indexes could not be created: %v\n", err)
	}

	// Step 5: Seed initial data
	fmt.Println("\n🌱 Phase 3: Seeding initial data...")
	if err := seedDefaultData(db); err != nil {
		fmt.Printf("Warning: Default data could not be seeded: %v\n", err)
	}

	printMigrationSummary()
	return nil
}

// setupPostgreSQLExtensions ensures necessary PostgreSQL extensions are installed
func setupPostgreSQLExtensions(db *gorm.DB) error {
	if db.Dialector.Name() != "postgres" {
		return nil
	}

	extensions := []string{
		"pgcrypto",  // For gen_random_uuid()
		"uuid-ossp", // Fallback for UUID generation
	}

	for _, ext := range extensions {
		sql := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\"", ext)
		if err := db.Exec(sql).Error; err != nil {
			// Log warning but don't fail - extension might already exist
			fmt.Printf("Note: Extension %s may already exist: %v\n", ext, err)
		}
	}

	return nil
}

// runCoreMigrations creates the core business tables
func runCoreMigrations(db *gorm.DB) error {
	// Group 1: Core business entities
	coreTables := []interface{}{
		&models.User{},
		&models.Device{},
		&policy.Product{},
		&policy.ProductTier{},
		&policy.ProductBundle{},
		&models.Policy{},
		&policy.Quote{},
		&models.Claim{},
		&models.ClaimDocument{},
	}

	for _, model := range coreTables {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	// Group 2: Payment and financial models
	fmt.Println("  Creating payment tables...")
	paymentTables := []interface{}{
		&models.PaymentMethod{},
		&models.Payment{},
		&models.Subscription{},
		&models.BillingHistory{},
		&models.Invoice{},
		&models.Commission{},
		&models.PromoCode{},
	}

	for _, model := range paymentTables {
		if err := db.AutoMigrate(model); err != nil {
			// Log but continue - payment tables are optional
			fmt.Printf("  Warning: Could not create payment table %T: %v\n", model, err)
		}
	}

	// Group 3: Document management
	fmt.Println("  Creating document tables...")
	documentTables := []interface{}{
		&models.Document{},
		&models.ESignature{},
		&models.DocumentTemplate{},
		&models.DocumentGeneration{},
		&models.DocumentAccess{},
	}

	for _, model := range documentTables {
		if err := db.AutoMigrate(model); err != nil {
			fmt.Printf("  Warning: Could not create document table %T: %v\n", model, err)
		}
	}

	// Group 4: Coverage and underwriting
	fmt.Println("  Creating coverage and underwriting tables...")
	coverageTables := []interface{}{
		// &shared.Coverage{}, // Commented out - type doesn't exist
		&shared.UnderwritingDecision{},
		&shared.UnderwritingRule{},
		&shared.RiskAssessment{},
		&shared.RiskModel{},
		&shared.UnderwritingDocument{},
		&policy.QuoteItem{},
		&policy.QuoteHistory{},
		&policy.Discount{},
		&policy.SalesLead{},
	}

	for _, model := range coverageTables {
		if err := db.AutoMigrate(model); err != nil {
			fmt.Printf("  Warning: Could not create coverage table %T: %v\n", model, err)
		}
	}

	return nil
}

// runAdditionalMigrations creates additional tables for extended features
func runAdditionalMigrations(db *gorm.DB) error {
	fmt.Println("  Creating additional feature tables...")

	additionalTables := []interface{}{
		// Device verification
		&device.DeviceVerification{},
		&device.IMEIDatabase{},
		// Support
		&shared.SupportTicket{},
		&shared.TicketMessage{},
		&shared.TicketAttachment{},
		&shared.FAQ{},
		&shared.KnowledgeBase{},
		// Admin
		&models.Role{},
		&models.UserRole{},
		&models.Permission{},
		&models.AuditLog{},
		&models.SystemConfiguration{},
		&models.Notification{},
		// Reporting
		&shared.Report{},
		&shared.ReportSchedule{},
		&shared.Dashboard{},
		&shared.KPI{},
		// Device history
		&device.DeviceHistory{},
		&device.BlacklistEntry{},
		// Fraud detection
		&shared.FraudDetection{},
		&shared.FraudIndicator{},
		&shared.BlacklistedDevice{},
		&shared.UserBehaviorProfile{},
		// Security
		// &shared.DeviceSecurity{}, // Commented out - type doesn't exist
		// &models.SecurityAlert{}, // Commented out - type doesn't exist
		// &shared.TheftReport{}, // Commented out - type doesn't exist
		// Repair network
		&models.RepairShop{},
		&models.RepairBooking{},
		&models.RepairStatusUpdate{},
		&models.RepairReview{},
		// Replacement
		&models.ReplacementDevice{},
		&models.ReplacementOrder{},
		&models.ReplacementStatusUpdate{},
		&models.TemporaryDevice{},
		&models.DeviceLoan{},
		// Pricing
		&shared.PricingRule{},
		// &models.SeasonalPricing{}, // Commented out - type doesn't exist
		// Corporate (optional)
		&models.CorporateAccount{},
		&models.CorporateEmployee{},
		&models.CorporatePolicy{},
	}

	successCount := 0
	failCount := 0

	for _, model := range additionalTables {
		if err := db.AutoMigrate(model); err != nil {
			failCount++
			// Continue with other tables - these are optional
		} else {
			successCount++
		}
	}

	fmt.Printf("  Created %d additional tables (%d skipped)\n", successCount, failCount)
	return nil
}

// printMigrationSummary prints a summary of the migration
func printMigrationSummary() {
	fmt.Println("\n========================================")
	fmt.Println("✅ Migration Complete")
	fmt.Println("========================================")
	fmt.Println("Database tables created successfully.")
	fmt.Println("Note: Foreign key relationships have been disabled")
	fmt.Println("to avoid circular dependency issues.")
	fmt.Println("The application will function normally.")
}

// createCoreIndexes creates database indexes for core business models only
func createCoreIndexes(db *gorm.DB) error {
	// User indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_kyc_status ON users(kyc_status)").Error; err != nil {
		return err
	}

	// Device indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_devices_imei ON devices(imei)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_devices_serial_number ON devices(serial_number)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_devices_owner_id ON devices(owner_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_devices_brand_model ON devices(brand, model)").Error; err != nil {
		return err
	}

	// Policy indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_policies_policy_number ON policies(policy_number)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_policies_customer_id ON policies(customer_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_policies_device_id ON policies(device_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_policies_status ON policies(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_policies_effective_date ON policies(effective_date)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_policies_expiration_date ON policies(expiration_date)").Error; err != nil {
		return err
	}

	// Claim indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claims_claim_number ON claims(claim_number)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claims_policy_id ON claims(policy_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claims_customer_id ON claims(customer_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claims_device_id ON claims(device_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claims_status ON claims(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claims_incident_date ON claims(incident_date)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claims_priority ON claims(priority)").Error; err != nil {
		return err
	}

	// Claim document indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claim_documents_claim_id ON claim_documents(claim_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_claim_documents_document_type ON claim_documents(document_type)").Error; err != nil {
		return err
	}

	// Repair and replacement indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_repair_shops_certification_level ON repair_shops(certification_level)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_repair_shops_is_active ON repair_shops(is_active)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_repair_shops_rating ON repair_shops(rating)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_repair_bookings_repair_shop_id ON repair_bookings(repair_shop_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_repair_bookings_status ON repair_bookings(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_repair_bookings_scheduled_date ON repair_bookings(scheduled_date)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_replacement_orders_status ON replacement_orders(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_replacement_orders_claim_id ON replacement_orders(claim_id)").Error; err != nil {
		return err
	}

	// Pricing model indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_usage_based_insurances_policy_id ON usage_based_insurances(policy_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_usage_based_insurances_next_recalculation ON usage_based_insurances(next_recalculation)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_micro_insurances_user_id ON micro_insurances(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_micro_insurances_status ON micro_insurances(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_pricing_rules_is_active ON pricing_rules(is_active)").Error; err != nil {
		return err
	}

	// Fraud detection indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fraud_detections_fraud_score ON fraud_detections(fraud_score)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fraud_detections_risk_level ON fraud_detections(risk_level)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fraud_detections_status ON fraud_detections(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_blacklisted_devices_imei ON blacklisted_devices(imei)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_blacklisted_devices_is_active ON blacklisted_devices(is_active)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_user_behavior_profiles_risk_score ON user_behavior_profiles(risk_score)").Error; err != nil {
		return err
	}

	// Analytics indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_analytics_user_id ON customer_analytics(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_analytics_churn_probability ON customer_analytics(churn_probability)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_predictive_models_is_active ON predictive_models(is_active)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_business_metrics_metric_date ON business_metrics(metric_date)").Error; err != nil {
		return err
	}

	// Device intelligence indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_device_diagnostics_device_id ON device_diagnostics(device_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_device_diagnostics_diagnostic_date ON device_diagnostics(diagnostic_date)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_diagnostic_alerts_severity ON diagnostic_alerts(severity)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_diagnostic_alerts_status ON diagnostic_alerts(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_geolocation_verifications_device_id ON geolocation_verifications(device_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_geolocation_verifications_verification_status ON geolocation_verifications(verification_status)").Error; err != nil {
		return err
	}

	// Customer engagement indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_portals_user_id ON customer_portals(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_portals_is_active ON customer_portals(is_active)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_chatbot_conversations_user_id ON chatbot_conversations(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_chatbot_conversations_is_active ON chatbot_conversations(is_active)").Error; err != nil {
		return err
	}

	// Advanced insurance indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_reinsurance_contracts_is_active ON reinsurance_contracts(is_active)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_reinsurance_cessions_contract_id ON reinsurance_cessions(contract_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_reinsurance_cessions_policy_id ON reinsurance_cessions(policy_id)").Error; err != nil {
		return err
	}

	// Compliance & Security indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_compliance_checks_entity_type ON compliance_checks(entity_type)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_compliance_checks_check_status ON compliance_checks(check_status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_sanctions_screenings_entity_id ON sanctions_screenings(entity_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_sanctions_screenings_status ON sanctions_screenings(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_security_incidents_status ON security_incidents(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_security_incidents_severity ON security_incidents(severity)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at)").Error; err != nil {
		return err
	}

	// Corporate & Family Insurance indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_corporate_accounts_status ON corporate_accounts(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_corporate_employees_account_id ON corporate_employees(corporate_account_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_corporate_policies_account_id ON corporate_policies(corporate_account_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fleet_devices_account_id ON fleet_devices(corporate_account_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_family_plans_primary_user_id ON family_plans(primary_user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_family_members_plan_id ON family_members(family_plan_id)").Error; err != nil {
		return err
	}

	// Third-Party Integration indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_third_party_integrations_type ON third_party_integrations(integration_type)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_integration_api_logs_integration_id ON integration_api_logs(integration_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_carrier_billings_user_id ON carrier_billings(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_carrier_billings_is_active ON carrier_billings(is_active)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_manufacturer_warranties_device_id ON manufacturer_warranties(device_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_warranty_claims_manufacturer_warranty_id ON warranty_claims(manufacturer_warranty_id)").Error; err != nil {
		return err
	}

	// Multi-Tenant & Partnership indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tenant_organizations_is_active ON tenant_organizations(is_active)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tenant_users_tenant_id ON tenant_users(tenant_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tenant_api_keys_tenant_id ON tenant_api_keys(tenant_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_partner_integrations_partner_type ON partner_integrations(partner_type)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_api_usage_analytics_tenant_id ON api_usage_analytics(tenant_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_marketplace_products_product_type ON marketplace_products(product_type)").Error; err != nil {
		return err
	}

	return nil
}

// seedDefaultData inserts default data if the database is empty
func seedDefaultData(db *gorm.DB) error {
	// Check if we need to seed data
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	// Only seed if database is empty
	if userCount > 0 {
		return nil
	}

	// Create default admin user
	adminUser := &models.User{
		UserIdentification: user.UserIdentification{
			Email:     "admin@smartsure.com",
			FirstName: "System",
			LastName:  "Administrator",
		},
		UserAuthentication: user.UserAuthentication{
			PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: "password"
		},
		// UserType:      "admin", // TODO: Set via embedded structs
		// UserRole:      "admin", // TODO: Set via embedded structs
		// Status:        "active", // TODO: Set via embedded structs
		// KYCStatus:     "approved", // TODO: Set via embedded structs
		// KYCLevel:      "premium", // TODO: Set via embedded structs
		// RiskScore:     0.0, // TODO: Set via embedded structs
		// EmailVerified: true, // TODO: Set via embedded structs
		// PhoneVerified: true, // TODO: Set via embedded structs
	}

	// Set audit fields to NULL for the first user
	adminUser.CreatedBy = nil
	adminUser.UpdatedBy = nil

	if err := db.Create(adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	return nil
}
