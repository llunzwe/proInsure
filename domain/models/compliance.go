package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/compliance"
	"smartsure/pkg/database"

	"gorm.io/gorm"
)

// Compliance represents the main compliance management entity
type Compliance struct {
	database.AuditableModel

	// Embedded structs for organization
	compliance.ComplianceIdentification `gorm:"embedded;embeddedPrefix:identification_"`
	compliance.ComplianceStatus         `gorm:"embedded;embeddedPrefix:status_"`
	compliance.ComplianceRiskProfile    `gorm:"embedded;embeddedPrefix:risk_"`
	compliance.ComplianceRegulatory     `gorm:"embedded;embeddedPrefix:regulatory_"`
	compliance.CompliancePrivacy        `gorm:"embedded;embeddedPrefix:privacy_"`
	compliance.ComplianceAudit          `gorm:"embedded;embeddedPrefix:audit_"`
	compliance.ComplianceControl        `gorm:"embedded;embeddedPrefix:control_"`
	compliance.ComplianceEvidence       `gorm:"embedded;embeddedPrefix:evidence_"`
	compliance.ComplianceMonitoring     `gorm:"embedded;embeddedPrefix:monitoring_"`
	compliance.ComplianceReporting      `gorm:"embedded;embeddedPrefix:reporting_"`
	compliance.ComplianceTraining       `gorm:"embedded;embeddedPrefix:training_"`
	compliance.ComplianceRelationships  `gorm:"embedded;embeddedPrefix:relationship_"`
	compliance.ComplianceMetadata       `gorm:"embedded;embeddedPrefix:metadata_"`

	// Feature model relationships
	Incidents                []compliance.ComplianceIncident                `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	Risks                    []compliance.ComplianceRisk                    `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	ConsentRecords           []compliance.ComplianceConsentRecord           `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	PrivacyImpactAssessments []compliance.CompliancePrivacyImpactAssessment `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	QualityMetrics           []compliance.ComplianceQualityMetric           `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	BusinessContinuityPlans  []compliance.ComplianceBusinessContinuityPlan  `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	ControlAssessments       []compliance.ComplianceControlAssessment       `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	VendorAssessments        []compliance.ComplianceVendorAssessment        `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	TrainingRecords          []compliance.ComplianceTrainingRecord          `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	ObligationRegisters      []compliance.ComplianceObligationRegister      `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	DataSubjectRequests      []compliance.ComplianceDataSubjectRequest      `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`
	SecurityControls         []compliance.ComplianceSecurityControl         `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE"`

	// Relationships with other entities
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	PolicyID   *uuid.UUID `gorm:"type:uuid;index" json:"policy_id,omitempty"`
	ClaimID    *uuid.UUID `gorm:"type:uuid;index" json:"claim_id,omitempty"`
	DeviceID   *uuid.UUID `gorm:"type:uuid;index" json:"device_id,omitempty"`
	DocumentID *uuid.UUID `gorm:"type:uuid;index" json:"document_id,omitempty"`

	// Relationships (lazy-loaded)
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Policy   *Policy   `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Claim    *Claim    `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	Device   *Device   `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	Document *Document `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
}

// BeforeCreate hook
func (c *Compliance) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	// Generate compliance number if not set
	if c.ComplianceNumber == "" {
		c.ComplianceNumber = c.GenerateComplianceNumber()
	}

	// Set default status
	if c.Status == "" {
		c.Status = compliance.ComplianceStatusPending
	}

	// Calculate initial risk score
	if c.RiskScore == 0 {
		c.RiskScore = c.CalculateRiskScore()
	}

	// Set next assessment date if not set
	if c.NextAssessmentDate.IsZero() {
		c.NextAssessmentDate = time.Now().AddDate(0, 3, 0) // 3 months default
	}

	return c.Validate()
}

// =====================================
// Validation Methods
// =====================================

// Validate performs comprehensive validation
func (c *Compliance) Validate() error {
	if c.ComplianceType == "" {
		return fmt.Errorf("compliance type is required")
	}
	if c.Title == "" {
		return fmt.Errorf("compliance title is required")
	}
	if c.Category == "" {
		return fmt.Errorf("compliance category is required")
	}

	// Validate risk scores
	if c.RiskScore < 0 || c.RiskScore > 100 {
		return fmt.Errorf("risk score must be between 0 and 100")
	}
	if c.ComplianceScore < 0 || c.ComplianceScore > 100 {
		return fmt.Errorf("compliance score must be between 0 and 100")
	}

	// Validate maturity level
	if c.MaturityLevel < 0 || c.MaturityLevel > 5 {
		return fmt.Errorf("maturity level must be between 0 and 5")
	}

	return nil
}

// =====================================
// Status and Lifecycle Methods
// =====================================

// IsCompliant checks if compliance is achieved
func (c *Compliance) IsCompliant() bool {
	return c.Status == compliance.ComplianceStatusCompliant
}

// IsNonCompliant checks if non-compliant
func (c *Compliance) IsNonCompliant() bool {
	return c.Status == compliance.ComplianceStatusNonCompliant
}

// IsPartiallyCompliant checks if partially compliant
func (c *Compliance) IsPartiallyCompliant() bool {
	return c.Status == compliance.ComplianceStatusPartial
}

// IsPendingAssessment checks if pending assessment
func (c *Compliance) IsPendingAssessment() bool {
	return c.Status == compliance.ComplianceStatusPending
}

// IsExpired checks if certification has expired
func (c *Compliance) IsExpired() bool {
	return c.Status == compliance.ComplianceStatusExpired ||
		(!c.CertificationExpiry.IsZero() && time.Now().After(c.CertificationExpiry))
}

// IsUnderRemediation checks if under remediation
func (c *Compliance) IsUnderRemediation() bool {
	return c.Status == compliance.ComplianceStatusRemediation
}

// IsActive checks if compliance is active
func (c *Compliance) IsActive() bool {
	return c.ComplianceStatus.IsActive && !c.IsExpired()
}

// IsCritical checks if compliance is critical
func (c *Compliance) IsCritical() bool {
	return c.ComplianceStatus.IsCritical || c.RiskLevel == compliance.RiskLevelCritical
}

// NeedsAssessment checks if assessment is due
func (c *Compliance) NeedsAssessment() bool {
	return time.Now().After(c.NextAssessmentDate) || c.IsPendingAssessment()
}

// NeedsAttestation checks if attestation is required
func (c *Compliance) NeedsAttestation() bool {
	return c.RequiresAttestation && !c.AttestationCompleted
}

// DaysUntilAssessment returns days until next assessment
func (c *Compliance) DaysUntilAssessment() int {
	if c.NextAssessmentDate.IsZero() {
		return -1
	}
	duration := time.Until(c.NextAssessmentDate)
	return int(duration.Hours() / 24)
}

// DaysUntilExpiry returns days until certification expiry
func (c *Compliance) DaysUntilExpiry() int {
	if c.CertificationExpiry.IsZero() {
		return -1
	}
	duration := time.Until(c.CertificationExpiry)
	return int(duration.Hours() / 24)
}

// SetStatus updates status with validation
func (c *Compliance) SetStatus(newStatus string) error {
	validStatuses := []string{
		compliance.ComplianceStatusCompliant,
		compliance.ComplianceStatusNonCompliant,
		compliance.ComplianceStatusPartial,
		compliance.ComplianceStatusPending,
		compliance.ComplianceStatusExpired,
		compliance.ComplianceStatusRemediation,
		compliance.ComplianceStatusException,
	}

	isValid := false
	for _, status := range validStatuses {
		if status == newStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid compliance status: %s", newStatus)
	}

	c.Status = newStatus
	return nil
}

// =====================================
// Risk Assessment Methods
// =====================================

// CalculateRiskScore calculates overall risk score
func (c *Compliance) CalculateRiskScore() float64 {
	score := 0.0
	factors := 0

	// Base risk score
	if c.InherentRisk > 0 {
		score += c.InherentRisk
		factors++
	}

	// Add risk level factor
	switch c.RiskLevel {
	case compliance.RiskLevelCritical:
		score += 90
		factors++
	case compliance.RiskLevelHigh:
		score += 70
		factors++
	case compliance.RiskLevelMedium:
		score += 50
		factors++
	case compliance.RiskLevelLow:
		score += 30
		factors++
	case compliance.RiskLevelMinimal:
		score += 10
		factors++
	}

	// Add control effectiveness factor
	switch c.ControlEffectiveness {
	case compliance.ControlEffectivenessIneffective:
		score += 30
		factors++
	case compliance.ControlEffectivenessPartial:
		score += 20
		factors++
	case compliance.ControlEffectivenessEffective:
		score += 10
		factors++
	}

	// Add compliance status factor
	if c.IsNonCompliant() {
		score += 20
		factors++
	} else if c.IsPartiallyCompliant() {
		score += 10
		factors++
	}

	if factors == 0 {
		return 50.0 // Default medium risk
	}

	return score / float64(factors)
}

// GetRiskCategory returns risk category based on score
func (c *Compliance) GetRiskCategory() string {
	score := c.RiskScore
	if score >= 80 {
		return compliance.RiskCategoryCyber
	} else if score >= 70 {
		return compliance.RiskCategoryCompliance
	} else if score >= 60 {
		return compliance.RiskCategoryOperational
	} else if score >= 50 {
		return compliance.RiskCategoryFinancial
	} else if score >= 40 {
		return compliance.RiskCategoryReputational
	} else if score >= 30 {
		return compliance.RiskCategoryStrategic
	}
	return compliance.RiskCategoryLegal
}

// IsHighRisk checks if high risk
func (c *Compliance) IsHighRisk() bool {
	return c.RiskScore >= 70 || c.RiskLevel == compliance.RiskLevelHigh || c.RiskLevel == compliance.RiskLevelCritical
}

// IsAcceptableRisk checks if risk is within appetite
func (c *Compliance) IsAcceptableRisk() bool {
	return c.ResidualRisk <= c.RiskAppetite
}

// NeedsRiskTreatment checks if risk treatment is needed
func (c *Compliance) NeedsRiskTreatment() bool {
	return c.ResidualRisk > c.RiskTolerance || c.IsHighRisk()
}

// =====================================
// Privacy and Data Protection Methods
// =====================================

// HasPIIData checks if PII is present
func (c *Compliance) HasPIIData() bool {
	return c.PIIPresent
}

// HasSensitiveData checks if sensitive data is present
func (c *Compliance) HasSensitiveData() bool {
	return c.SensitiveDataPresent
}

// RequiresConsent checks if consent is required
func (c *Compliance) RequiresConsent() bool {
	return c.ConsentRequired
}

// HasValidConsent checks if valid consent exists
func (c *Compliance) HasValidConsent() bool {
	return c.ConsentRequired && c.ConsentObtained
}

// RequiresPIA checks if Privacy Impact Assessment is needed
func (c *Compliance) RequiresPIA() bool {
	return c.HasPIIData() || c.HasSensitiveData() || c.CrossBorderTransfer
}

// NeedsPIA checks if PIA is due
func (c *Compliance) NeedsPIA() bool {
	return time.Now().After(c.NextPIADate) || c.RequiresPIA()
}

// HasDPOApproval checks if DPO has approved
func (c *Compliance) HasDPOApproval() bool {
	return c.DPOApproval && c.DPOApprovedBy != uuid.Nil
}

// GetDataRetentionDays returns data retention period
func (c *Compliance) GetDataRetentionDays() int {
	if c.DataRetentionPeriod > 0 {
		return c.DataRetentionPeriod
	}
	return 365 * 7 // Default 7 years
}

// =====================================
// Audit and Assessment Methods
// =====================================

// IsAuditDue checks if audit is due
func (c *Compliance) IsAuditDue() bool {
	return time.Now().After(c.NextAuditDate)
}

// GetAuditStatus returns current audit status
func (c *Compliance) GetAuditStatus() string {
	if c.AuditStatus != "" {
		return c.AuditStatus
	}
	if c.IsAuditDue() {
		return compliance.AuditStatusScheduled
	}
	return compliance.AuditStatusCompleted
}

// HasCriticalFindings checks for critical audit findings
func (c *Compliance) HasCriticalFindings() bool {
	return c.CriticalFindings > 0
}

// HasNonConformities checks for non-conformities
func (c *Compliance) HasNonConformities() bool {
	return c.NonConformities > 0
}

// GetTotalFindings returns total audit findings
func (c *Compliance) GetTotalFindings() int {
	return c.CriticalFindings + c.MajorFindings + c.MinorFindings + c.Observations
}

// NeedsCorrectiveAction checks if corrective action is needed
func (c *Compliance) NeedsCorrectiveAction() bool {
	return c.HasCriticalFindings() || c.HasNonConformities() || c.CorrectiveActions > 0
}

// GetEffectivenessRating returns control effectiveness rating
func (c *Compliance) GetEffectivenessRating() string {
	if c.EffectivenessRating != "" {
		return c.EffectivenessRating
	}

	// Calculate based on compliance score
	if c.ComplianceScore >= 90 {
		return compliance.ControlEffectivenessEffective
	} else if c.ComplianceScore >= 70 {
		return compliance.ControlEffectivenessPartial
	} else if c.ComplianceScore >= 50 {
		return compliance.ControlEffectivenessPartial
	}
	return compliance.ControlEffectivenessIneffective
}

// =====================================
// Monitoring and Reporting Methods
// =====================================

// IsMonitoringEnabled checks if monitoring is enabled
func (c *Compliance) IsMonitoringEnabled() bool {
	return c.MonitoringEnabled
}

// IsRealtimeMonitoring checks if realtime monitoring
func (c *Compliance) IsRealtimeMonitoring() bool {
	return c.RealtimeMonitoring && c.MonitoringEnabled
}

// NeedsMonitoring checks if monitoring is due
func (c *Compliance) NeedsMonitoring() bool {
	return time.Now().After(c.NextMonitoringDate) || c.IsHighRisk()
}

// IsReportingRequired checks if reporting is required
func (c *Compliance) IsReportingRequired() bool {
	return c.ReportingRequired
}

// IsRegulatoryReporting checks if regulatory reporting needed
func (c *Compliance) IsRegulatoryReporting() bool {
	return c.RegulatoryReporting
}

// NeedsReporting checks if report is due
func (c *Compliance) NeedsReporting() bool {
	return time.Now().After(c.NextReportDate) && c.IsReportingRequired()
}

// GetReportingFrequency returns reporting frequency
func (c *Compliance) GetReportingFrequency() string {
	if c.ReportingFrequency != "" {
		return c.ReportingFrequency
	}
	if c.IsCritical() {
		return "monthly"
	} else if c.IsHighRisk() {
		return "quarterly"
	}
	return "annually"
}

// =====================================
// Training Methods
// =====================================

// IsTrainingRequired checks if training is required
func (c *Compliance) IsTrainingRequired() bool {
	return c.TrainingRequired
}

// IsTrainingDue checks if training is due
func (c *Compliance) IsTrainingDue() bool {
	return time.Now().After(c.NextTrainingDate) && c.IsTrainingRequired()
}

// GetTrainingCompletionRate returns training completion percentage
func (c *Compliance) GetTrainingCompletionRate() float64 {
	if c.CompletionRate > 0 {
		return c.CompletionRate
	}
	return 0
}

// IsTrainingCompliant checks if training compliance met
func (c *Compliance) IsTrainingCompliant() bool {
	return c.CompletionRate >= 95.0 && c.PassRate >= 80.0
}

// NeedsRefresherTraining checks if refresher training needed
func (c *Compliance) NeedsRefresherTraining() bool {
	return c.RefresherRequired && time.Now().After(c.NextTrainingDate)
}

// =====================================
// Regulatory and Legal Methods
// =====================================

// GetApplicableRegulations returns applicable regulations
func (c *Compliance) GetApplicableRegulations() []string {
	return c.RegulatoryBodies
}

// IsGDPRApplicable checks if GDPR applies
func (c *Compliance) IsGDPRApplicable() bool {
	for _, reg := range c.RegulatoryBodies {
		if strings.Contains(strings.ToUpper(reg), "GDPR") {
			return true
		}
	}
	return false
}

// IsCCPAApplicable checks if CCPA applies
func (c *Compliance) IsCCPAApplicable() bool {
	for _, reg := range c.RegulatoryBodies {
		if strings.Contains(strings.ToUpper(reg), "CCPA") {
			return true
		}
	}
	return false
}

// IsHIPAAApplicable checks if HIPAA applies
func (c *Compliance) IsHIPAAApplicable() bool {
	for _, reg := range c.RegulatoryBodies {
		if strings.Contains(strings.ToUpper(reg), "HIPAA") {
			return true
		}
	}
	return false
}

// HasRegulatoryDeadline checks if regulatory deadline exists
func (c *Compliance) HasRegulatoryDeadline() bool {
	return len(c.RegulatoryDeadlines) > 0
}

// GetUpcomingDeadlines returns deadlines within 30 days
func (c *Compliance) GetUpcomingDeadlines() map[string]time.Time {
	upcoming := make(map[string]time.Time)
	thirtyDays := time.Now().AddDate(0, 0, 30)

	for key, deadline := range c.RegulatoryDeadlines {
		if deadline.After(time.Now()) && deadline.Before(thirtyDays) {
			upcoming[key] = deadline
		}
	}

	return upcoming
}

// GetMaxPenalty returns maximum penalty amount
func (c *Compliance) GetMaxPenalty() float64 {
	return c.PenaltyAmount
}

// HasEnforcementAction checks if enforcement actions exist
func (c *Compliance) HasEnforcementAction() bool {
	return c.EnforcementActions > 0
}

// =====================================
// Incident Management Methods
// =====================================

// HasActiveIncident checks for active incidents
func (c *Compliance) HasActiveIncident() bool {
	for _, incident := range c.Incidents {
		if incident.Status == compliance.IncidentStatusOpen ||
			incident.Status == compliance.IncidentStatusInvestigating {
			return true
		}
	}
	return false
}

// GetActiveIncidentCount returns count of active incidents
func (c *Compliance) GetActiveIncidentCount() int {
	count := 0
	for _, incident := range c.Incidents {
		if incident.Status != compliance.IncidentStatusClosed {
			count++
		}
	}
	return count
}

// GetCriticalIncidentCount returns count of critical incidents
func (c *Compliance) GetCriticalIncidentCount() int {
	count := 0
	for _, incident := range c.Incidents {
		if incident.Severity == compliance.IncidentSeverityCritical {
			count++
		}
	}
	return count
}

// =====================================
// Helper Methods
// =====================================

// GenerateComplianceNumber generates unique compliance number
func (c *Compliance) GenerateComplianceNumber() string {
	prefix := "COMP"
	if c.ComplianceType != "" {
		switch c.ComplianceType {
		case compliance.ISO27001:
			prefix = "ISO27001"
		case compliance.ISO27701:
			prefix = "ISO27701"
		case compliance.ISO31000:
			prefix = "ISO31000"
		case compliance.ISO9001:
			prefix = "ISO9001"
		default:
			prefix = "COMP"
		}
	}

	timestamp := time.Now().Format("20060102")
	random := uuid.New().String()[:8]
	return fmt.Sprintf("%s-%s-%s", prefix, timestamp, strings.ToUpper(random))
}

// GetDisplayName returns formatted display name
func (c *Compliance) GetDisplayName() string {
	if c.Title != "" {
		return c.Title
	}
	return fmt.Sprintf("%s - %s", c.ComplianceType, c.ComplianceNumber)
}

// GetMaturityLevel returns maturity level
func (c *Compliance) GetMaturityLevel() int {
	if c.MaturityLevel > 0 {
		return c.MaturityLevel
	}

	// Calculate based on compliance score
	if c.ComplianceScore >= 90 {
		return 5
	} else if c.ComplianceScore >= 75 {
		return 4
	} else if c.ComplianceScore >= 60 {
		return 3
	} else if c.ComplianceScore >= 40 {
		return 2
	}
	return 1
}

// GetMaturityDescription returns maturity level description
func (c *Compliance) GetMaturityDescription() string {
	switch c.GetMaturityLevel() {
	case 5:
		return "Optimized"
	case 4:
		return "Managed"
	case 3:
		return "Defined"
	case 2:
		return "Repeatable"
	case 1:
		return "Initial"
	default:
		return "Unknown"
	}
}

// GetPriority returns compliance priority
func (c *Compliance) GetPriority() int {
	if c.Priority > 0 {
		return c.Priority
	}

	// Calculate based on risk and criticality
	if c.IsCritical() {
		return 1
	} else if c.IsHighRisk() {
		return 2
	} else if c.IsNonCompliant() {
		return 3
	} else if c.IsPartiallyCompliant() {
		return 4
	}
	return 5
}

// GetOwnerName returns compliance owner name
func (c *Compliance) GetOwnerName() string {
	if c.ComplianceOfficer != uuid.Nil {
		// This would typically fetch from User model
		return "Compliance Officer"
	}
	return "Unassigned"
}

// GetCompletionPercentage returns overall completion percentage
func (c *Compliance) GetCompletionPercentage() float64 {
	if c.ComplianceScore > 0 {
		return c.ComplianceScore
	}

	// Calculate based on various factors
	total := 0.0
	count := 0

	if c.IsCompliant() {
		total += 100
		count++
	} else if c.IsPartiallyCompliant() {
		total += 50
		count++
	}

	if c.TrainingRequired && c.CompletionRate > 0 {
		total += c.CompletionRate
		count++
	}

	if c.AuditFindings == 0 && c.NonConformities == 0 {
		total += 100
		count++
	}

	if count == 0 {
		return 0
	}

	return total / float64(count)
}

// GetSummary returns compliance summary
func (c *Compliance) GetSummary() string {
	parts := []string{
		fmt.Sprintf("Type: %s", c.ComplianceType),
		fmt.Sprintf("Status: %s", c.Status),
		fmt.Sprintf("Risk: %s", c.RiskLevel),
		fmt.Sprintf("Score: %.1f%%", c.ComplianceScore),
		fmt.Sprintf("Maturity: Level %d", c.GetMaturityLevel()),
	}

	if c.HasActiveIncident() {
		parts = append(parts, fmt.Sprintf("Active Incidents: %d", c.GetActiveIncidentCount()))
	}

	if c.HasCriticalFindings() {
		parts = append(parts, fmt.Sprintf("Critical Findings: %d", c.CriticalFindings))
	}

	return strings.Join(parts, ", ")
}

// TableName specifies the table name for GORM
func (Compliance) TableName() string {
	return "compliances"
}
