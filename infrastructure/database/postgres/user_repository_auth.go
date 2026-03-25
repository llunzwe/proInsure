package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// === Authentication & Security Methods ===

// UpdatePassword updates user password
func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"password_hash":       passwordHash,
		"password_changed_at": time.Now(),
	}, &models.User{}, nil)
}

// UpdateLastLogin updates last login time
func (r *UserRepositoryImpl) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"last_login_at":  time.Now(),
		"last_active_at": time.Now(),
	}, &models.User{}, nil)
}

// IncrementLoginAttempts increments login attempts
func (r *UserRepositoryImpl) IncrementLoginAttempts(ctx context.Context, userID uuid.UUID) (int, error) {
	db := r.GetDB(ctx, nil)

	var user models.User
	if err := db.Model(&models.User{}).Where("id = ?", userID).
		Update("failed_login_attempts", gorm.Expr("failed_login_attempts + ?", 1)).
		First(&user).Error; err != nil {
		return 0, err
	}

	return user.FailedLoginAttempts, nil
}

// ResetLoginAttempts resets login attempts
func (r *UserRepositoryImpl) ResetLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"failed_login_attempts": 0,
	}, &models.User{}, nil)
}

// LockAccount locks user account
func (r *UserRepositoryImpl) LockAccount(ctx context.Context, userID uuid.UUID, until time.Time) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"account_locked": true,
		"locked_until":   until,
	}, &models.User{}, nil)
}

// UnlockAccount unlocks user account
func (r *UserRepositoryImpl) UnlockAccount(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"account_locked": false,
		"locked_until":   nil,
	}, &models.User{}, nil)
}

// Enable2FA enables two-factor authentication
func (r *UserRepositoryImpl) Enable2FA(ctx context.Context, userID uuid.UUID, secret string) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"two_factor_enabled": true,
		"two_factor_secret":  secret,
	}, &models.User{}, nil)
}

// Disable2FA disables two-factor authentication
func (r *UserRepositoryImpl) Disable2FA(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"two_factor_enabled": false,
		"two_factor_secret":  nil,
	}, &models.User{}, nil)
}

// UpdateSecurityQuestions updates security questions
func (r *UserRepositoryImpl) UpdateSecurityQuestions(ctx context.Context, userID uuid.UUID, questions map[string]string) error {
	questionsJSON, err := json.Marshal(questions)
	if err != nil {
		return err
	}

	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"security_questions": questionsJSON,
	}, &models.User{}, nil)
}

// LogSecurityEvent logs a security event
func (r *UserRepositoryImpl) LogSecurityEvent(ctx context.Context, event *models.UserSecurityEvent) error {
	return r.Create(ctx, event, nil)
}

// GetSecurityEvents retrieves security events for a user
func (r *UserRepositoryImpl) GetSecurityEvents(ctx context.Context, userID uuid.UUID, limit int) ([]*models.UserSecurityEvent, error) {
	var events []*models.UserSecurityEvent
	db := r.GetDB(ctx, nil)

	query := db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// === Verification & Compliance ===

// VerifyEmail verifies user email
func (r *UserRepositoryImpl) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"email_verified":    true,
		"email_verified_at": time.Now(),
	}, &models.User{}, nil)
}

// VerifyPhone verifies user phone
func (r *UserRepositoryImpl) VerifyPhone(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"phone_verified":    true,
		"phone_verified_at": time.Now(),
	}, &models.User{}, nil)
}

// UpdateKYCStatus updates KYC status
func (r *UserRepositoryImpl) UpdateKYCStatus(ctx context.Context, userID uuid.UUID, status, level string) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"kyc_status":      status,
		"kyc_level":       level,
		"kyc_verified_at": time.Now(),
	}, &models.User{}, nil)
}

// UpdateAMLStatus updates AML status
func (r *UserRepositoryImpl) UpdateAMLStatus(ctx context.Context, userID uuid.UUID, status string) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"aml_status":     status,
		"aml_last_check": time.Now(),
	}, &models.User{}, nil)
}

// UpdatePEPStatus updates PEP status
func (r *UserRepositoryImpl) UpdatePEPStatus(ctx context.Context, userID uuid.UUID, isPEP bool) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"pep_status":     isPEP,
		"pep_checked_at": time.Now(),
	}, &models.User{}, nil)
}

// UpdateSanctionsScreening updates sanctions screening
func (r *UserRepositoryImpl) UpdateSanctionsScreening(ctx context.Context, userID uuid.UUID, screened bool, date time.Time) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"sanctions_screened":       screened,
		"sanctions_screening_date": date,
	}, &models.User{}, nil)
}

// UpdateGDPRConsent updates GDPR consent
func (r *UserRepositoryImpl) UpdateGDPRConsent(ctx context.Context, userID uuid.UUID, consent bool) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"gdpr_consent":      consent,
		"gdpr_consent_date": time.Now(),
	}, &models.User{}, nil)
}

// UpdateTermsAcceptance updates terms acceptance
func (r *UserRepositoryImpl) UpdateTermsAcceptance(ctx context.Context, userID uuid.UUID, version string) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"terms_accepted":    true,
		"terms_version":     version,
		"terms_accepted_at": time.Now(),
	}, &models.User{}, nil)
}

// GetComplianceStatus gets compliance status
func (r *UserRepositoryImpl) GetComplianceStatus(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	var user models.User
	db := r.GetDB(ctx, nil)

	if err := db.Select("kyc_status", "kyc_level", "aml_status", "pep_status",
		"sanctions_screened", "gdpr_consent", "terms_accepted").
		Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"kyc_status":         user.KYCStatus,
		"kyc_level":          user.KYCLevel,
		"aml_status":         user.AMLStatus,
		"pep_status":         user.PEPStatus,
		"sanctions_screened": user.SanctionsScreened,
		"gdpr_consent":       user.GDPRConsent,
		"terms_accepted":     user.TermsAccepted,
	}, nil
}

// GetUsersRequiringCompliance gets users requiring compliance checks
func (r *UserRepositoryImpl) GetUsersRequiringCompliance(ctx context.Context, checkType string) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	query := db.Model(&models.User{})

	switch checkType {
	case "kyc":
		query = query.Where("kyc_status IN ? OR kyc_verified_at < ?",
			[]string{"pending", "expired"}, time.Now().AddDate(0, -6, 0))
	case "aml":
		query = query.Where("aml_status = ? OR aml_last_check < ?",
			"pending", time.Now().AddDate(0, -1, 0))
	case "pep":
		query = query.Where("pep_checked_at IS NULL OR pep_checked_at < ?",
			time.Now().AddDate(0, -3, 0))
	case "sanctions":
		query = query.Where("sanctions_screening_date IS NULL OR sanctions_screening_date < ?",
			time.Now().AddDate(0, -1, 0))
	default:
		return nil, fmt.Errorf("unknown check type: %s", checkType)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// === Risk & Fraud Management ===

// UpdateRiskScore updates risk score
func (r *UserRepositoryImpl) UpdateRiskScore(ctx context.Context, userID uuid.UUID, score float64) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"risk_score":           score,
		"risk_assessment_date": time.Now(),
	}, &models.User{}, nil)
}

// UpdateFraudScore updates fraud score
func (r *UserRepositoryImpl) UpdateFraudScore(ctx context.Context, userID uuid.UUID, score float64) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"fraud_score":      score,
		"fraud_check_date": time.Now(),
	}, &models.User{}, nil)
}

// UpdateCreditScore updates credit score
func (r *UserRepositoryImpl) UpdateCreditScore(ctx context.Context, userID uuid.UUID, score int) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"credit_score":      score,
		"credit_check_date": time.Now(),
	}, &models.User{}, nil)
}

// BlacklistUser blacklists a user
func (r *UserRepositoryImpl) BlacklistUser(ctx context.Context, userID uuid.UUID, reason string) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"blacklist_status": true,
		"blacklist_reason": reason,
		"blacklisted_at":   time.Now(),
	}, &models.User{}, nil)
}

// UnblacklistUser removes user from blacklist
func (r *UserRepositoryImpl) UnblacklistUser(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"blacklist_status": false,
		"blacklist_reason": nil,
		"blacklisted_at":   nil,
	}, &models.User{}, nil)
}

// GetBlacklistedUsers gets all blacklisted users
func (r *UserRepositoryImpl) GetBlacklistedUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("blacklist_status = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// IncrementSuspiciousActivity increments suspicious activity counter
func (r *UserRepositoryImpl) IncrementSuspiciousActivity(ctx context.Context, userID uuid.UUID) (int, error) {
	db := r.GetDB(ctx, nil)

	var user models.User
	if err := db.Model(&models.User{}).Where("id = ?", userID).
		Update("suspicious_activity_count", gorm.Expr("suspicious_activity_count + ?", 1)).
		First(&user).Error; err != nil {
		return 0, err
	}

	return user.SuspiciousActivityCount, nil
}

// GetHighRiskUsers gets users with high risk scores
func (r *UserRepositoryImpl) GetHighRiskUsers(ctx context.Context, threshold float64) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("risk_score >= ?", threshold).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetFraudAlerts gets fraud alerts for a user
func (r *UserRepositoryImpl) GetFraudAlerts(ctx context.Context, userID uuid.UUID) ([]map[string]interface{}, error) {
	var alerts []map[string]interface{}
	db := r.GetDB(ctx, nil)

	// Query fraud-related security events
	var events []models.UserSecurityEvent
	if err := db.Where("user_id = ? AND event_type LIKE ?", userID, "%fraud%").
		Order("created_at DESC").
		Limit(10).
		Find(&events).Error; err != nil {
		return nil, err
	}

	for _, event := range events {
		alerts = append(alerts, map[string]interface{}{
			"id":          event.ID,
			"type":        event.EventType,
			"description": event.Description,
			"severity":    event.Severity,
			"created_at":  event.CreatedAt,
		})
	}

	return alerts, nil
}

// === Financial Operations ===

// UpdateOutstandingBalance updates outstanding balance
func (r *UserRepositoryImpl) UpdateOutstandingBalance(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"outstanding_balance": amount,
	}, &models.User{}, nil)
}

// UpdateCreditLimit updates credit limit
func (r *UserRepositoryImpl) UpdateCreditLimit(ctx context.Context, userID uuid.UUID, limit decimal.Decimal) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"credit_limit": limit,
	}, &models.User{}, nil)
}

// RecordPaymentFailure records a payment failure
func (r *UserRepositoryImpl) RecordPaymentFailure(ctx context.Context, userID uuid.UUID) (int, error) {
	db := r.GetDB(ctx, nil)

	var user models.User
	if err := db.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"payment_failure_count": gorm.Expr("payment_failure_count + ?", 1),
			"last_payment_failure":  time.Now(),
		}).First(&user).Error; err != nil {
		return 0, err
	}

	return user.PaymentFailureCount, nil
}

// RecordPaymentSuccess records a successful payment
func (r *UserRepositoryImpl) RecordPaymentSuccess(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error {
	db := r.GetDB(ctx, nil)

	return db.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_payment_date":   time.Now(),
			"last_payment_amount": amount,
			"total_premium_paid":  gorm.Expr("total_premium_paid + ?", amount),
		}).Error
}

// UpdateTotalPremiumPaid updates total premium paid
func (r *UserRepositoryImpl) UpdateTotalPremiumPaid(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"total_premium_paid": amount,
	}, &models.User{}, nil)
}

// GetFinancialSummary gets financial summary for a user
func (r *UserRepositoryImpl) GetFinancialSummary(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	var user models.User
	db := r.GetDB(ctx, nil)

	if err := db.Select("outstanding_balance", "credit_limit", "total_premium_paid",
		"payment_failure_count", "last_payment_date", "last_payment_amount").
		Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"outstanding_balance":   user.OutstandingBalance,
		"credit_limit":          user.CreditLimit,
		"total_premium_paid":    user.TotalPremiumPaid,
		"payment_failure_count": user.PaymentFailureCount,
		"last_payment_date":     user.LastPaymentDate,
		"last_payment_amount":   user.LastPaymentAmount,
	}, nil
}

// GetUsersWithOutstandingBalance gets users with outstanding balance
func (r *UserRepositoryImpl) GetUsersWithOutstandingBalance(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	if err := db.Where("outstanding_balance > ?", 0).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// UpdatePreferredPaymentMethod updates preferred payment method
func (r *UserRepositoryImpl) UpdatePreferredPaymentMethod(ctx context.Context, userID uuid.UUID, method string) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"preferred_payment_method": method,
	}, &models.User{}, nil)
}

// === Insurance Operations ===

// UpdateInsuranceProfile updates insurance profile
func (r *UserRepositoryImpl) UpdateInsuranceProfile(ctx context.Context, userID uuid.UUID, profile map[string]interface{}) error {
	return r.UpdateFields(ctx, userID, profile, &models.User{}, nil)
}

// UpdatePolicyCount updates policy counts
func (r *UserRepositoryImpl) UpdatePolicyCount(ctx context.Context, userID uuid.UUID, active, total int) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"active_policies_count": active,
		"total_policies_count":  total,
	}, &models.User{}, nil)
}

// UpdateClaimStats updates claim statistics
func (r *UserRepositoryImpl) UpdateClaimStats(ctx context.Context, userID uuid.UUID, total, approved int, avgAmount decimal.Decimal) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"total_claims_count":    total,
		"approved_claims_count": approved,
		"average_claim_amount":  avgAmount,
	}, &models.User{}, nil)
}

// RecordClaim records a new claim for a user
func (r *UserRepositoryImpl) RecordClaim(ctx context.Context, userID uuid.UUID, claimID uuid.UUID) error {
	db := r.GetDB(ctx, nil)

	return db.Model(&models.User{}).Where("id = ?", userID).
		Update("total_claims_count", gorm.Expr("total_claims_count + ?", 1)).Error
}

// UpdateDeviceCount updates device count
func (r *UserRepositoryImpl) UpdateDeviceCount(ctx context.Context, userID uuid.UUID, count int) error {
	return r.UpdateFields(ctx, userID, map[string]interface{}{
		"device_count": count,
	}, &models.User{}, nil)
}

// GetInsuranceHistory gets insurance history for a user
func (r *UserRepositoryImpl) GetInsuranceHistory(ctx context.Context, userID uuid.UUID) ([]*models.UserInsuranceHistory, error) {
	var history []*models.UserInsuranceHistory
	db := r.GetDB(ctx, nil)

	if err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// CreateInsuranceHistory creates insurance history record
func (r *UserRepositoryImpl) CreateInsuranceHistory(ctx context.Context, history *models.UserInsuranceHistory) error {
	return r.Create(ctx, history, nil)
}

// GetUsersWithExpiringPolicies gets users with expiring policies
func (r *UserRepositoryImpl) GetUsersWithExpiringPolicies(ctx context.Context, days int) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	expiryDate := time.Now().AddDate(0, 0, days)

	// This query would need to join with policies table
	query := db.Table("users").
		Joins("JOIN policies ON users.id = policies.customer_id").
		Where("policies.status = ? AND policies.end_date <= ?", "active", expiryDate).
		Distinct("users.id")

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersEligibleForUpgrade gets users eligible for upgrade
func (r *UserRepositoryImpl) GetUsersEligibleForUpgrade(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	db := r.GetDB(ctx, nil)

	// Users with high engagement, good payment history, and basic tier
	if err := db.Where("subscription_tier = ? AND engagement_score >= ? AND payment_failure_count < ?",
		"basic", 7.0, 2).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
