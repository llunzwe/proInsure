// Package services defines the service interfaces following hexagonal architecture
package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/models/claim"
)

// ClaimRequest represents a request to create or update a claim
type ClaimRequest struct {
	PolicyID        uuid.UUID              `json:"policy_id"`
	DeviceID        uuid.UUID              `json:"device_id"`
	ClaimType       string                 `json:"claim_type"`
	IncidentDate    time.Time              `json:"incident_date"`
	Description     string                 `json:"description"`
	ClaimedAmount   float64                `json:"claimed_amount"`
	Location        string                 `json:"location"`
	Evidence        map[string]interface{} `json:"evidence"`
	WitnessContacts []string               `json:"witness_contacts"`
	PoliceReport    *string                `json:"police_report"`
}

// ClaimAssessment represents the result of claim assessment
type ClaimAssessment struct {
	ClaimID               uuid.UUID `json:"claim_id"`
	Eligible              bool      `json:"eligible"`
	EstimatedPayout       float64   `json:"estimated_payout"`
	DeductibleAmount      float64   `json:"deductible_amount"`
	FraudRiskScore        float64   `json:"fraud_risk_score"`
	RequiresInvestigation bool      `json:"requires_investigation"`
	RequiresDocuments     []string  `json:"requires_documents"`
	RecommendedAction     string    `json:"recommended_action"`
	Reasons               []string  `json:"reasons"`
}

// ClaimProcessingRequest represents a request to process a claim
type ClaimProcessingRequest struct {
	ClaimID        uuid.UUID              `json:"claim_id"`
	Action         string                 `json:"action"` // approve, deny, investigate, request_info
	ApprovedAmount *float64               `json:"approved_amount"`
	DenialReason   *string                `json:"denial_reason"`
	ProcessorNotes string                 `json:"processor_notes"`
	AdditionalData map[string]interface{} `json:"additional_data"`
}

// CreateClaimRequest represents a request to create a new claim
type CreateClaimRequest struct {
	PolicyID      uuid.UUID `json:"policy_id"`
	PolicyNumber  string    `json:"policy_number"`
	IncidentType  string    `json:"incident_type" binding:"required"`
	IncidentDate  string    `json:"incident_date" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	Location      string    `json:"location"`
	EstimatedCost float64   `json:"estimated_cost"`
}

// ClaimAssessmentResponse represents the result of claim assessment
type ClaimAssessmentResponse struct {
	ClaimID          uuid.UUID `json:"claim_id"`
	Eligible         bool      `json:"eligible"`
	EstimatedAmount  float64   `json:"estimated_amount"`
	AssessmentStatus string    `json:"assessment_status"`
	FraudScore       float64   `json:"fraud_score"`
	RiskLevel        string    `json:"risk_level"`
	AssessmentNotes  string    `json:"assessment_notes"`
}

// ClaimSettlementRequest represents a settlement request
type ClaimSettlementRequest struct {
	ClaimID          uuid.UUID         `json:"claim_id"`
	SettlementAmount float64           `json:"settlement_amount"`
	PaymentMethod    string            `json:"payment_method"`
	BankDetails      map[string]string `json:"bank_details"`
	TaxDeductions    float64           `json:"tax_deductions"`
}

// ClaimMetrics represents claim performance metrics
type ClaimMetrics struct {
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	ApprovalRate          float64       `json:"approval_rate"`
	AveragePayout         float64       `json:"average_payout"`
	FraudDetectionRate    float64       `json:"fraud_detection_rate"`
	CustomerSatisfaction  float64       `json:"customer_satisfaction"`
	SLACompliance         float64       `json:"sla_compliance"`
}

// ClaimService defines the comprehensive business operations for claims
// This interface represents the application layer for claim management
type ClaimService interface {
	// ======================================
	// CLAIM CREATION & SUBMISSION
	// ======================================

	// Claim submission
	SubmitClaim(ctx context.Context, customerID uuid.UUID, request *ClaimRequest) (*models.Claim, error)
	SaveDraftClaim(ctx context.Context, customerID uuid.UUID, draft *ClaimRequest) (uuid.UUID, error)
	ValidateClaim(ctx context.Context, claim *ClaimRequest) (bool, []string, error)
	CheckEligibility(ctx context.Context, policyID, deviceID uuid.UUID, claimType string) (bool, []string, error)

	// Duplicate detection
	CheckForDuplicates(ctx context.Context, deviceID uuid.UUID, incidentDate time.Time, claimType string) ([]*models.Claim, error)

	// ======================================
	// CLAIM ASSESSMENT & EVALUATION
	// ======================================

	// Assessment
	AssessClaim(ctx context.Context, claimID uuid.UUID) (*ClaimAssessment, error)
	CalculateEstimatedPayout(ctx context.Context, claimID uuid.UUID) (float64, error)
	EvaluateFraudRisk(ctx context.Context, claimID uuid.UUID) (float64, map[string]interface{}, error)
	DetermineRequiredDocuments(ctx context.Context, claimID uuid.UUID) ([]string, error)

	// Risk scoring
	CalculateClaimRiskScore(ctx context.Context, claimID uuid.UUID) (float64, map[string]float64, error)
	GetHistoricalClaimPatterns(ctx context.Context, customerID uuid.UUID) (map[string]interface{}, error)

	// ======================================
	// CLAIM PROCESSING & WORKFLOW
	// ======================================

	// Processing
	ProcessClaim(ctx context.Context, request *ClaimProcessingRequest, processorID uuid.UUID) error
	ApproveClaim(ctx context.Context, claimID uuid.UUID, approvedAmount float64, approverID uuid.UUID) error
	DenyClaim(ctx context.Context, claimID uuid.UUID, reason string, denierID uuid.UUID) error
	RequestAdditionalInformation(ctx context.Context, claimID uuid.UUID, requiredInfo []string) error

	// Workflow management
	InitiateWorkflow(ctx context.Context, claimID uuid.UUID, workflowType string) (*claim.ClaimWorkflow, error)
	AdvanceWorkflowStage(ctx context.Context, claimID uuid.UUID, action string, performedBy uuid.UUID) error
	GetWorkflowStatus(ctx context.Context, claimID uuid.UUID) (*claim.ClaimWorkflow, error)
	EscalateClaimToSupervisor(ctx context.Context, claimID uuid.UUID, reason string) error

	// Assignment
	AssignClaimToAdjuster(ctx context.Context, claimID, adjusterID uuid.UUID) error
	AutoAssignClaim(ctx context.Context, claimID uuid.UUID) (uuid.UUID, error)
	ReassignClaim(ctx context.Context, claimID, newAdjusterID uuid.UUID, reason string) error
	BalanceAdjusterWorkload(ctx context.Context) error

	// ======================================
	// INVESTIGATION & FRAUD DETECTION
	// ======================================

	// Investigation
	InitiateInvestigation(ctx context.Context, claimID uuid.UUID, reason string) (*claim.ClaimInvestigationDetail, error)
	AssignInvestigator(ctx context.Context, claimID, investigatorID uuid.UUID) error
	RecordInvestigationFindings(ctx context.Context, claimID uuid.UUID, findings map[string]interface{}) error
	CompleteInvestigation(ctx context.Context, claimID uuid.UUID, conclusion string, recommendedAction string) error

	// Fraud detection
	RunFraudDetection(ctx context.Context, claimID uuid.UUID) (*claim.ClaimFraudDetection, error)
	AnalyzeFraudPatterns(ctx context.Context, claimID uuid.UUID) ([]string, float64, error)
	FlagSuspiciousClaim(ctx context.Context, claimID uuid.UUID, indicators []string) error
	CreateFraudAlert(ctx context.Context, claimID uuid.UUID, severity string, details map[string]interface{}) error

	// Network analysis
	DetectFraudNetwork(ctx context.Context, claimID uuid.UUID) ([]uuid.UUID, map[string]interface{}, error)
	AnalyzeClaimVelocity(ctx context.Context, customerID uuid.UUID, period time.Duration) (int, bool, error)

	// ======================================
	// SETTLEMENT & PAYMENT
	// ======================================

	// Settlement
	InitiateSettlement(ctx context.Context, claimID uuid.UUID) (*claim.ClaimSettlementDetail, error)
	CalculateFinalSettlement(ctx context.Context, claimID uuid.UUID) (float64, map[string]float64, error)
	ProcessSettlement(ctx context.Context, request *ClaimSettlementRequest) error
	SchedulePayment(ctx context.Context, claimID uuid.UUID, paymentDate time.Time) error

	// Payment processing
	ProcessClaimPayment(ctx context.Context, claimID uuid.UUID) (*claim.ClaimPayment, error)
	VerifyPaymentDetails(ctx context.Context, claimID uuid.UUID, bankDetails map[string]string) (bool, error)
	RetryFailedPayment(ctx context.Context, paymentID uuid.UUID) error
	ReconcilePayment(ctx context.Context, claimID uuid.UUID, transactionID string) error

	// Reserve management
	SetInitialReserve(ctx context.Context, claimID uuid.UUID, amount float64) error
	UpdateReserve(ctx context.Context, claimID uuid.UUID, newAmount float64, reason string) error
	CalculateReserveAdequacy(ctx context.Context, claimID uuid.UUID) (float64, bool, error)

	// ======================================
	// SMARTPHONE-SPECIFIC FEATURES
	// ======================================

	// Device diagnostics
	RunRemoteDiagnostics(ctx context.Context, claimID uuid.UUID) (*claim.ClaimDeviceDiagnostics, error)
	AnalyzeDiagnosticResults(ctx context.Context, claimID uuid.UUID) (map[string]interface{}, error)
	ValidateDeviceCondition(ctx context.Context, claimID uuid.UUID, diagnosticData map[string]interface{}) (bool, []string, error)

	// Repair management
	FindNearestRepairCenter(ctx context.Context, claimID uuid.UUID, location string) ([]map[string]interface{}, error)
	ScheduleRepair(ctx context.Context, claimID uuid.UUID, repairCenterID uuid.UUID, scheduledDate time.Time) error
	TrackRepairStatus(ctx context.Context, claimID uuid.UUID) (string, map[string]interface{}, error)
	ApproveRepairQuote(ctx context.Context, claimID uuid.UUID, quoteID uuid.UUID) error

	// Replacement device
	CheckReplacementEligibility(ctx context.Context, claimID uuid.UUID) (bool, []string, error)
	InitiateDeviceReplacement(ctx context.Context, claimID uuid.UUID, replacementType string) (*claim.ClaimReplacementDevice, error)
	TrackReplacementShipment(ctx context.Context, claimID uuid.UUID) (string, map[string]interface{}, error)
	ConfirmDeviceReceived(ctx context.Context, claimID uuid.UUID, serialNumber string) error

	// Digital assets
	AssessDigitalLoss(ctx context.Context, claimID uuid.UUID) (float64, []string, error)
	InitiateDataRecovery(ctx context.Context, claimID uuid.UUID) (bool, map[string]interface{}, error)

	// ======================================
	// APPEAL & DISPUTE RESOLUTION
	// ======================================

	// Appeal management
	SubmitAppeal(ctx context.Context, claimID uuid.UUID, appealReason string, supportingDocs []uuid.UUID) (*claim.ClaimAppeal, error)
	ReviewAppeal(ctx context.Context, appealID uuid.UUID, reviewerID uuid.UUID) error
	MakeAppealDecision(ctx context.Context, appealID uuid.UUID, decision string, rationale string) error
	EscalateAppeal(ctx context.Context, appealID uuid.UUID, escalationLevel int) error

	// Dispute resolution
	InitiateDispute(ctx context.Context, claimID uuid.UUID, disputeType string, details map[string]interface{}) error
	ScheduleArbitration(ctx context.Context, claimID uuid.UUID, arbitratorID uuid.UUID, scheduledDate time.Time) error
	RecordArbitrationDecision(ctx context.Context, arbitrationID uuid.UUID, decision string, award float64) error

	// ======================================
	// LITIGATION SUPPORT
	// ======================================

	// Litigation management
	InitiateLitigation(ctx context.Context, claimID uuid.UUID, legalGrounds string) (*claim.ClaimLitigation, error)
	UpdateLitigationStatus(ctx context.Context, litigationID uuid.UUID, status string, updates map[string]interface{}) error
	RecordCourtDecision(ctx context.Context, litigationID uuid.UUID, decision string, awardAmount float64) error
	CalculateLitigationCosts(ctx context.Context, claimID uuid.UUID) (float64, map[string]float64, error)

	// ======================================
	// COMMUNICATION & CUSTOMER SERVICE
	// ======================================

	// Communication
	SendClaimNotification(ctx context.Context, claimID uuid.UUID, notificationType string, channel string) error
	RecordCustomerInteraction(ctx context.Context, claimID uuid.UUID, interactionType string, notes string) error
	GetCommunicationHistory(ctx context.Context, claimID uuid.UUID) ([]*claim.ClaimCommunication, error)
	ScheduleFollowUp(ctx context.Context, claimID uuid.UUID, followUpDate time.Time, reason string) error

	// Customer satisfaction
	RequestFeedback(ctx context.Context, claimID uuid.UUID) error
	RecordSatisfactionScore(ctx context.Context, claimID uuid.UUID, score int, feedback string) error
	AnalyzeSentiment(ctx context.Context, claimID uuid.UUID) (float64, string, error)

	// ======================================
	// REPORTING & ANALYTICS
	// ======================================

	// Reporting
	GenerateClaimReport(ctx context.Context, claimID uuid.UUID, reportType string) ([]byte, error)
	GenerateSettlementLetter(ctx context.Context, claimID uuid.UUID) ([]byte, error)
	CreateManagementDashboard(ctx context.Context, filters map[string]interface{}) (map[string]interface{}, error)

	// Analytics
	GetClaimMetrics(ctx context.Context, startDate, endDate time.Time) (*ClaimMetrics, error)
	AnalyzeClaimTrends(ctx context.Context, period string) ([]map[string]interface{}, error)
	PredictClaimOutcome(ctx context.Context, claimID uuid.UUID) (string, float64, error)
	CalculateLossRatio(ctx context.Context, productID uuid.UUID, period time.Duration) (float64, error)

	// Predictive analytics
	PredictProcessingTime(ctx context.Context, claimID uuid.UUID) (time.Duration, error)
	PredictSettlementAmount(ctx context.Context, claimID uuid.UUID) (float64, float64, error)
	IdentifyHighRiskClaims(ctx context.Context, threshold float64) ([]*models.Claim, error)

	// ======================================
	// COMPLIANCE & AUDIT
	// ======================================

	// Compliance
	CheckRegulatoryCompliance(ctx context.Context, claimID uuid.UUID) (bool, []string, error)
	GenerateComplianceReport(ctx context.Context, startDate, endDate time.Time) ([]byte, error)
	EnsureDataPrivacy(ctx context.Context, claimID uuid.UUID) error
	HandleGDPRRequest(ctx context.Context, claimID uuid.UUID, requestType string) error

	// Audit
	GetAuditTrail(ctx context.Context, claimID uuid.UUID) ([]map[string]interface{}, error)
	LogClaimActivity(ctx context.Context, claimID uuid.UUID, activity string, details map[string]interface{}) error
	GenerateAuditReport(ctx context.Context, claimID uuid.UUID) ([]byte, error)

	// ======================================
	// INTEGRATION & AUTOMATION
	// ======================================

	// Third-party integration
	SyncWithExternalSystem(ctx context.Context, claimID uuid.UUID, systemName string) error
	FetchExternalData(ctx context.Context, claimID uuid.UUID, dataSource string) (map[string]interface{}, error)

	// Automation
	ApplyBusinessRules(ctx context.Context, claimID uuid.UUID) ([]string, error)
	ExecuteAutomatedWorkflow(ctx context.Context, claimID uuid.UUID, workflowName string) error
	OptimizeClaimRouting(ctx context.Context, claimID uuid.UUID) (string, uuid.UUID, error)

	// ======================================
	// QUALITY ASSURANCE
	// ======================================

	// Quality control
	PerformQualityCheck(ctx context.Context, claimID uuid.UUID) (bool, []string, error)
	ReviewAdjusterDecision(ctx context.Context, claimID uuid.UUID, reviewerID uuid.UUID) (bool, string, error)
	CalculateQualityScore(ctx context.Context, claimID uuid.UUID) (float64, map[string]float64, error)

	// ======================================
	// MISSING IMPLEMENTED METHODS
	// ======================================

	// Assessment Management
	GetClaimAssessmentDetails(ctx context.Context, claimID uuid.UUID) (map[string]interface{}, error)
	RequestClaimReassessment(ctx context.Context, claimID uuid.UUID, reason, notes string) error
	OverrideClaimAssessment(ctx context.Context, claimID uuid.UUID, newAmount float64, newStatus, overrideReason, approvedBy, notes string) error

	// Payment Management
	InitiateClaimPayment(ctx context.Context, claimID uuid.UUID, paymentMethodID uuid.UUID, notes string) (uuid.UUID, error)
	GetClaimPaymentStatus(ctx context.Context, claimID uuid.UUID) (map[string]interface{}, error)
	UpdateClaimPaymentDetails(ctx context.Context, claimID uuid.UUID, paymentMethodID uuid.UUID, paymentDate, referenceNumber, notes string) error

	// Appeals Management
	CreateClaimAppeal(ctx context.Context, claimID uuid.UUID, appealReason, description string, evidence []string, notes string) (*claim.ClaimAppeal, error)
	GetClaimAppeals(ctx context.Context, claimID uuid.UUID) ([]*claim.ClaimAppeal, error)
	UpdateClaimAppeal(ctx context.Context, claimID uuid.UUID, appealID uuid.UUID, status, decision, notes, reviewedBy string) error

	// Fraud Detection
	GetClaimFraudScore(ctx context.Context, claimID uuid.UUID) (map[string]interface{}, error)
	StartClaimFraudInvestigation(ctx context.Context, claimID uuid.UUID, investigationReason, assignedInvestigator, priority, notes string, evidence []string) (interface{}, error) // TODO: Fix type
	UpdateClaimFraudStatus(ctx context.Context, claimID uuid.UUID, fraudStatus string, confidence float64, findings, investigator, notes string) error

	// Analytics & Reporting
	GetClaimAnalytics(ctx context.Context, claimID uuid.UUID) (map[string]interface{}, error)
	GetSettlementTrends(ctx context.Context, startDate, endDate, category, region string) (map[string]interface{}, error)
	GetFraudStatistics(ctx context.Context, startDate, endDate, category, severity string) (map[string]interface{}, error)

	// Communication
	SendClaimUpdate(ctx context.Context, claimID uuid.UUID, updateType, message, priority string, includeDocs, sendEmail, sendSMS bool) error
	GetClaimCommunicationLog(ctx context.Context, claimID uuid.UUID, limit, offset int) ([]*claim.ClaimCommunication, int, error)

	// Legacy methods for backward compatibility (already implemented)
	CreateClaim(ctx context.Context, req *CreateClaimRequest) (*models.Claim, error)
	GetClaimByID(ctx context.Context, claimID uuid.UUID) (*models.Claim, error)
	GetClaimByNumber(ctx context.Context, claimNumber string) (*models.Claim, error)
	// ProcessClaim(ctx context.Context, req *ClaimProcessingRequest) (*models.Claim, error) // Duplicate - removed, see line 127
	// AssessClaim(ctx context.Context, claimID uuid.UUID) (*ClaimAssessmentResponse, error) // Duplicate - removed, see line 113
	GetClaimsByCustomer(ctx context.Context, customerID uuid.UUID) ([]models.Claim, error)
	GetClaimsByStatus(ctx context.Context, status string) ([]models.Claim, error)
	GetPendingClaims(ctx context.Context) ([]models.Claim, error)
	AddClaimDocument(ctx context.Context, claimID uuid.UUID, docType, fileName, filePath string, fileSize int64, mimeType string, uploadedBy uuid.UUID) (*models.ClaimDocument, error)
	GetUserClaims(ctx context.Context, userID uuid.UUID) ([]*models.Claim, error)
	UpdateClaim(ctx context.Context, claimID uuid.UUID, updates map[string]interface{}) (*models.Claim, error)
	DeleteClaim(ctx context.Context, claimID uuid.UUID) error
}
