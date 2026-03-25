package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models/claim"
	"smartsure/pkg/database"
)

// Claim represents a comprehensive insurance claim with all enterprise features
type Claim struct {
	database.AuditableModel

	// Embedded structs for logical organization
	ClaimIdentification claim.ClaimIdentification `gorm:"embedded;embeddedPrefix:id_" json:"identification"`
	ClaimFinancial      claim.ClaimFinancial      `gorm:"embedded;embeddedPrefix:fin_" json:"financial"`
	ClaimLifecycle      claim.ClaimLifecycle      `gorm:"embedded;embeddedPrefix:lc_" json:"lifecycle"`
	ClaimInvestigation  claim.ClaimInvestigation  `gorm:"embedded;embeddedPrefix:inv_" json:"investigation"`
	ClaimSettlement     claim.ClaimSettlement     `gorm:"embedded;embeddedPrefix:set_" json:"settlement"`
	ClaimDocumentation  claim.ClaimDocumentation  `gorm:"embedded;embeddedPrefix:doc_" json:"documentation"`
	ClaimAssignment     claim.ClaimAssignment     `gorm:"embedded;embeddedPrefix:asg_" json:"assignment"`
	ClaimCompliance     claim.ClaimCompliance     `gorm:"embedded;embeddedPrefix:cmp_" json:"compliance"`
	ClaimMetrics        claim.ClaimMetrics        `gorm:"embedded;embeddedPrefix:met_" json:"metrics"`

	// Core Foreign Keys
	PolicyID   uuid.UUID `gorm:"type:uuid;not null" json:"policy_id"`
	CustomerID uuid.UUID `gorm:"type:uuid;not null" json:"customer_id"`
	DeviceID   uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`

	// Core Relationships
	Policy   *Policy `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Customer *User   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Device   *Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`

	// Feature Model Relationships (One-to-One)
	Investigation    *claim.ClaimInvestigationDetail `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"investigation_detail,omitempty"`
	FraudDetection   *claim.ClaimFraudDetection      `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"fraud_detection,omitempty"`
	Workflow         *claim.ClaimWorkflow            `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"workflow,omitempty"`
	SettlementDetail *claim.ClaimSettlementDetail    `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"settlement_detail,omitempty"`
	PaymentInfo      *claim.ClaimPayment             `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"payment_info,omitempty"`
	Reserve          *claim.ClaimReserve             `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"reserve,omitempty"`
	Subrogation      *claim.ClaimSubrogation         `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"subrogation,omitempty"`
	Appeal           *claim.ClaimAppeal              `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"appeal,omitempty"`
	Communication    *claim.ClaimCommunication       `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"communication,omitempty"`
	Analytics        *claim.ClaimAnalytics           `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"analytics,omitempty"`
	Litigation       *claim.ClaimLitigation          `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"litigation,omitempty"`
	Arbitration      *claim.ClaimArbitration         `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"arbitration,omitempty"`
	Reporting        *claim.ClaimReporting           `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"reporting,omitempty"`

	// One-to-Many Relationships
	Documents []Document           `gorm:"foreignKey:ClaimID" json:"documents,omitempty"`
	Payments  []claim.ClaimPayment `gorm:"foreignKey:ClaimID" json:"payments,omitempty"`

	// Smartphone-Specific Feature Models (One-to-One)
	DeviceDiagnostics  *claim.ClaimDeviceDiagnostics  `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"device_diagnostics,omitempty"`
	RepairNetwork      *claim.ClaimRepairNetwork      `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"repair_network,omitempty"`
	ReplacementDevice  *claim.ClaimReplacementDevice  `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"replacement_device,omitempty"`
	DigitalAssets      *claim.ClaimDigitalAssets      `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"digital_assets,omitempty"`
	Accessories        *claim.ClaimAccessories        `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"accessories,omitempty"`
	Geolocation        *claim.ClaimGeolocation        `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"geolocation,omitempty"`
	PreventiveMeasures *claim.ClaimPreventiveMeasures `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"preventive_measures,omitempty"`
	SelfService        *claim.ClaimSelfService        `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"self_service,omitempty"`

	// Advanced Feature Models (One-to-One)
	BiometricVerification *claim.ClaimBiometricVerification `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"biometric_verification,omitempty"`
	IoT5G                 *claim.Claim5GAndIoT              `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"iot_5g,omitempty"`
	AugmentedReality      *claim.ClaimAugmentedReality      `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"augmented_reality,omitempty"`
	Cryptocurrency        *claim.ClaimCryptocurrency        `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"cryptocurrency,omitempty"`
	SubscriptionServices  *claim.ClaimSubscriptionServices  `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"subscription_services,omitempty"`
	EnvironmentalImpact   *claim.ClaimEnvironmentalImpact   `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"environmental_impact,omitempty"`
	FoldableDevice        *claim.ClaimFoldableFlexible      `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"foldable_device,omitempty"`
	HealthWellness        *claim.ClaimHealthAndWellness     `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"health_wellness,omitempty"`
	BusinessContinuity    *claim.ClaimBusinessContinuity    `gorm:"foreignKey:ClaimID;constraint:OnDelete:CASCADE" json:"business_continuity,omitempty"`

	// Additional fields for advanced features
	MicroInsuranceID   *uuid.UUID `gorm:"type:uuid" json:"micro_insurance_id"`
	RepairBookingID    *uuid.UUID `gorm:"type:uuid" json:"repair_booking_id"`
	ReplacementOrderID *uuid.UUID `gorm:"type:uuid" json:"replacement_order_id"`
	TheftReportID      *uuid.UUID `gorm:"type:uuid" json:"theft_report_id"`
}

// ClaimDocument represents a document attached to a claim
type ClaimDocument struct {
	database.BaseModel
	ClaimID      uuid.UUID `gorm:"type:uuid;not null" json:"claim_id"`
	DocumentType string    `gorm:"type:varchar(50);not null" json:"document_type"`
	FileName     string    `gorm:"not null" json:"file_name"`
	FilePath     string    `gorm:"not null" json:"file_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	UploadedBy   uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
}

// TableName returns the table name for Claim model
func (Claim) TableName() string {
	return "claims"
}

// TableName returns the table name for ClaimDocument model
func (ClaimDocument) TableName() string {
	return "claim_documents"
}

// BeforeCreate handles pre-creation logic with business identifier generation
func (c *Claim) BeforeCreate(tx *gorm.DB) error {
	// Call parent BeforeCreate to handle UUID generation and audit fields
	if err := c.AuditableModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Generate claim number if not provided
	if c.ClaimIdentification.ClaimNumber == "" {
		generator := database.NewBusinessIdentifierGenerator()
		c.ClaimIdentification.ClaimNumber = generator.GenerateClaimNumber()
	}

	// Set initial status
	if c.ClaimLifecycle.Status == "" {
		c.ClaimLifecycle.Status = claim.StatusSubmitted
	}

	// Initialize metrics
	c.ClaimMetrics.TouchPoints = 1
	c.ClaimMetrics.CustomerInteractions = 1

	// Calculate initial reserve
	if c.ClaimFinancial.ClaimedAmount > 0 {
		c.ClaimFinancial.ReserveAmount = c.ClaimFinancial.ClaimedAmount * claim.ReserveInitialPercent
	}

	return nil
}

// BeforeCreate handles pre-creation logic (UUID generation handled by BaseModel)
func (cd *ClaimDocument) BeforeCreate(tx *gorm.DB) error {
	// Call parent BeforeCreate to handle UUID generation
	if err := cd.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// IsOpen checks if the claim is still open for processing
func (c *Claim) IsOpen() bool {
	openStatuses := []string{
		claim.StatusSubmitted, claim.StatusUnderReview, claim.StatusInvestigating,
		claim.StatusPendingDocuments, claim.StatusPendingApproval,
	}
	for _, status := range openStatuses {
		if c.ClaimLifecycle.Status == status {
			return true
		}
	}
	return false
}

// IsClosed checks if the claim has been closed
func (c *Claim) IsClosed() bool {
	closedStatuses := []string{
		claim.StatusApproved, claim.StatusDenied, claim.StatusSettled,
		claim.StatusWithdrawn, claim.StatusExpired, claim.StatusClosed,
	}
	for _, status := range closedStatuses {
		if c.ClaimLifecycle.Status == status {
			return true
		}
	}
	return false
}

// IsSettled checks if the claim has been settled
func (c *Claim) IsSettled() bool {
	return c.ClaimLifecycle.Status == claim.StatusSettled &&
		c.ClaimLifecycle.SettlementDate != nil
}

// DaysOpen returns the number of days the claim has been open
func (c *Claim) DaysOpen() int {
	if c.IsClosed() && c.ClaimLifecycle.ClosedDate != nil {
		return int(c.ClaimLifecycle.ClosedDate.Sub(c.ClaimLifecycle.ReportedDate).Hours() / 24)
	}
	return int(time.Since(c.ClaimLifecycle.ReportedDate).Hours() / 24)
}

// CanApprove checks if the claim can be approved
func (c *Claim) CanApprove() bool {
	return c.IsOpen() &&
		c.ClaimLifecycle.Status != claim.StatusDenied &&
		c.ClaimInvestigation.FraudScore < claim.ThresholdFraudScore
}

// generateClaimNumber is deprecated - use BusinessIdentifierGenerator instead
// Kept for backward compatibility
func generateClaimNumber() string {
	generator := database.NewBusinessIdentifierGenerator()
	return generator.GenerateClaimNumber()
}

// Validate performs validation on the claim
func (c *Claim) Validate() error {
	if c.ClaimFinancial.ClaimedAmount <= 0 {
		return gorm.ErrInvalidData
	}
	if c.ClaimLifecycle.IncidentDate.After(time.Now()) {
		return gorm.ErrInvalidData
	}
	if c.ClaimIdentification.ClaimType == "" {
		return gorm.ErrInvalidData
	}
	return nil
}

// CalculatePayout calculates the final payout amount
func (c *Claim) CalculatePayout(policyDeductible float64, coverageLimit float64) float64 {
	payout := c.ClaimFinancial.ClaimedAmount - policyDeductible

	// Apply deductible
	if payout < 0 {
		payout = 0
	}

	// Apply coverage limit
	if payout > coverageLimit {
		payout = coverageLimit
	}

	// Apply approved amount if set
	if c.ClaimFinancial.ApprovedAmount > 0 && c.ClaimFinancial.ApprovedAmount < payout {
		payout = c.ClaimFinancial.ApprovedAmount
	}

	// Apply depreciation if applicable
	if c.ClaimFinancial.DepreciationApplied > 0 {
		payout -= c.ClaimFinancial.DepreciationApplied
	}

	c.ClaimFinancial.TotalPayout = payout
	return payout
}

// SetStatus updates the claim status with validation
func (c *Claim) SetStatus(newStatus string) error {
	validStatuses := []string{
		claim.StatusSubmitted, claim.StatusUnderReview, claim.StatusInvestigating,
		claim.StatusPendingDocuments, claim.StatusPendingApproval, claim.StatusApproved,
		claim.StatusDenied, claim.StatusSettled, claim.StatusWithdrawn, claim.StatusExpired,
		claim.StatusAppeal, claim.StatusClosed,
	}

	validTransition := false
	for _, status := range validStatuses {
		if newStatus == status {
			validTransition = true
			break
		}
	}

	if !validTransition {
		return gorm.ErrInvalidData
	}

	// Validate status transitions
	if c.IsClosed() && newStatus != claim.StatusWithdrawn && newStatus != claim.StatusAppeal {
		return gorm.ErrInvalidData // Cannot reopen closed claims except for appeal
	}

	c.ClaimLifecycle.PreviousStatus = c.ClaimLifecycle.Status
	c.ClaimLifecycle.Status = newStatus

	// Update relevant dates
	now := time.Now()
	switch newStatus {
	case claim.StatusApproved:
		c.ClaimLifecycle.ApprovalDate = &now
	case claim.StatusDenied:
		c.ClaimLifecycle.DenialDate = &now
	case claim.StatusSettled:
		c.ClaimLifecycle.SettlementDate = &now
	case claim.StatusClosed:
		c.ClaimLifecycle.ClosedDate = &now
	}

	return nil
}

// IsHighPriority determines if claim should be prioritized
func (c *Claim) IsHighPriority() bool {
	return c.ClaimIdentification.Priority == claim.PriorityHigh ||
		c.ClaimIdentification.Priority == claim.PriorityUrgent ||
		c.ClaimIdentification.Priority == claim.PriorityCritical ||
		c.ClaimFinancial.ClaimedAmount > claim.ThresholdHighValue ||
		c.ClaimIdentification.ClaimType == claim.TypeTheft
}

// RequiresReview checks if claim needs manual review
func (c *Claim) RequiresReview() bool {
	return c.ClaimInvestigation.FraudScore > 0.5 ||
		c.ClaimFinancial.ClaimedAmount > claim.ThresholdInvestigation ||
		c.ClaimInvestigation.RequiresInvestigation
}

// GetProcessingTime returns expected processing time in days
func (c *Claim) GetProcessingTime() int {
	if c.IsHighPriority() {
		return 2
	}
	if c.RequiresReview() {
		return 7
	}
	return 5
}

// GetSLAHours returns the SLA deadline in hours based on priority
func (c *Claim) GetSLAHours() int {
	switch c.ClaimIdentification.Priority {
	case claim.PriorityUrgent, claim.PriorityCritical:
		return claim.SLAUrgent
	case claim.PriorityHigh:
		return claim.SLAHigh
	case claim.PriorityMedium:
		return claim.SLAMedium
	case claim.PriorityLow:
		return claim.SLALow
	default:
		return claim.SLAStandard
	}
}

// CanAutoApprove checks if claim can be auto-approved
func (c *Claim) CanAutoApprove() bool {
	return c.IsOpen() &&
		c.ClaimFinancial.ClaimedAmount <= claim.ThresholdAutoApproval &&
		c.ClaimInvestigation.FraudScore < 0.3 &&
		!c.ClaimInvestigation.RequiresInvestigation
}

// AssessFraudRisk performs fraud risk assessment
func (c *Claim) AssessFraudRisk() string {
	score := c.ClaimInvestigation.FraudScore

	if score < 0.3 {
		c.ClaimInvestigation.FraudRiskLevel = claim.FraudRiskLow
	} else if score < 0.5 {
		c.ClaimInvestigation.FraudRiskLevel = claim.FraudRiskMedium
	} else if score < 0.7 {
		c.ClaimInvestigation.FraudRiskLevel = claim.FraudRiskHigh
	} else {
		c.ClaimInvestigation.FraudRiskLevel = claim.FraudRiskCritical
	}

	return c.ClaimInvestigation.FraudRiskLevel
}

// Smartphone-Specific Methods

// IsSmartphoneClaim checks if this is a smartphone/mobile device claim
func (c *Claim) IsSmartphoneClaim() bool {
	if c.Device == nil {
		return false
	}
	category := string(c.Device.DeviceClassification.Category)
	return category == "smartphone" || category == "mobile" || category == "tablet"
}

// RequiresDeviceDiagnostics determines if remote diagnostics should be run
func (c *Claim) RequiresDeviceDiagnostics() bool {
	return c.ClaimIdentification.ClaimType == claim.TypeMalfunction ||
		c.ClaimIdentification.ClaimType == claim.TypeBatteryIssue ||
		c.ClaimIdentification.ClaimType == claim.TypeSoftwareIssue
}

// IsEligibleForInstantReplacement checks instant replacement eligibility
func (c *Claim) IsEligibleForInstantReplacement() bool {
	if c.ReplacementDevice == nil {
		return false
	}
	return c.ClaimFinancial.ClaimedAmount < 1000 &&
		c.ClaimInvestigation.FraudScore < 0.3 &&
		c.GetClaimHistory() < 3 // Less than 3 claims in past year
}

// GetDigitalAssetsValue calculates total digital assets value
func (c *Claim) GetDigitalAssetsValue() float64 {
	if c.DigitalAssets != nil {
		return c.DigitalAssets.CalculateDigitalLoss()
	}
	return 0.0
}

// RequiresSpecializedRepair checks if specialized repair is needed
func (c *Claim) RequiresSpecializedRepair() bool {
	if c.FoldableDevice != nil && c.FoldableDevice.IsHighComplexityRepair() {
		return true
	}
	if c.AugmentedReality != nil && c.AugmentedReality.RequiresSpecializedSupport() {
		return true
	}
	return false
}

// GetEnvironmentalScore returns sustainability score
func (c *Claim) GetEnvironmentalScore() float64 {
	if c.EnvironmentalImpact != nil {
		return c.EnvironmentalImpact.GetSustainabilityScore()
	}
	return 50.0 // Default neutral score
}

// IsBusinessCritical checks if claim is business critical
func (c *Claim) IsBusinessCritical() bool {
	if c.BusinessContinuity != nil {
		return c.BusinessContinuity.RequiresPriorityHandling()
	}
	return false
}

// GetRepairVsReplaceRecommendation suggests repair or replacement
func (c *Claim) GetRepairVsReplaceRecommendation() string {
	if c.DeviceDiagnostics != nil {
		healthScore := c.DeviceDiagnostics.CalculateHealthScore()
		if healthScore < 30 {
			return "replace"
		} else if healthScore < 60 {
			return "repair_or_replace"
		}
		return "repair"
	}
	return "requires_assessment"
}

// IsSelfServiceEligible checks if claim can be self-serviced
func (c *Claim) IsSelfServiceEligible() bool {
	return c.ClaimFinancial.ClaimedAmount < 500 &&
		c.ClaimIdentification.ClaimType != claim.TypeTheft &&
		c.ClaimInvestigation.FraudScore < 0.2
}

// GetClaimHistory returns number of previous claims (mock implementation)
func (c *Claim) GetClaimHistory() int {
	// This would query the database for previous claims by the same customer
	// Mock implementation for now
	return 0
}
