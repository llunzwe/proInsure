package device

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	
	"smartsure/internal/domain/models"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// deviceService implements the DeviceService interface
type deviceService struct {
	repo repositories.DeviceRepository
}

// NewDeviceService creates a new device service
func NewDeviceService(repo repositories.DeviceRepository) services.DeviceService {
	return &deviceService{
		repo: repo,
	}
}

// ============================================
// INSURANCE & RISK METHODS
// ============================================

// CheckInsuranceEligibility checks if device is eligible for insurance
func (s *deviceService) CheckInsuranceEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return false, "", err
	}

	if device.DeviceStatusInfo.Status != "active" {
		return false, "Device is not active", nil
	}

	if device.DeviceStatusInfo.IsStolen {
		return false, "Device is reported stolen", nil
	}

	if string(device.DevicePhysicalCondition.Grade) == "F" {
		return false, "Device condition is too poor", nil
	}

	// Check if device has too many claims
	if len(device.Claims) > 3 {
		return false, "Device has too many previous claims", nil
	}

	// Check blacklist status
	if device.DeviceRiskAssessment.BlacklistStatus == "blocked" {
		return false, "Device is blacklisted", nil
	}

	return true, "Eligible for insurance", nil
}

// CalculateDepreciation calculates current value based on age and condition
func (s *deviceService) CalculateDepreciation(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	if device.DeviceFinancial.PurchaseDate == nil || device.DeviceFinancial.PurchasePrice.Amount == 0 {
		return device.DeviceFinancial.CurrentValue.Amount, nil
	}

	age := time.Since(*device.DeviceFinancial.PurchaseDate).Hours() / (24 * 365) // years

	// Different depreciation rates based on device segment
	depreciationRate := 0.2 // Default 20% per year
	switch device.DeviceClassification.Segment {
	case "flagship":
		depreciationRate = 0.18 // 18% per year
	case "premium":
		depreciationRate = 0.2 // 20% per year
	case "mid_range":
		depreciationRate = 0.25 // 25% per year
	case "budget":
		depreciationRate = 0.3 // 30% per year
	}

	conditionMultiplier := map[string]float64{
		"excellent": 1.0,
		"good":      0.85,
		"fair":      0.7,
		"poor":      0.5,
	}

	multiplier, exists := conditionMultiplier[string(device.DevicePhysicalCondition.Condition)]
	if !exists {
		multiplier = 0.8
	}

	depreciated := device.DeviceFinancial.PurchasePrice.Amount * (1 - (depreciationRate * age)) * multiplier

	// Minimum value based on segment
	minPercentage := 0.1 // 10% default
	if device.DeviceClassification.Segment == "flagship" {
		minPercentage = 0.15 // 15% for flagship
	}

	if depreciated < device.DeviceFinancial.PurchasePrice.Amount*minPercentage {
		depreciated = device.DeviceFinancial.PurchasePrice.Amount * minPercentage
	}

	return depreciated, nil
}

// AssessRisk performs risk assessment on the device
func (s *deviceService) AssessRisk(ctx context.Context, deviceID uuid.UUID) (*services.RiskAssessment, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	riskScore := device.DeviceRiskAssessment.RiskScore
	if riskScore == 0 {
		// Calculate risk score based on various factors
		riskScore = s.calculateRiskScore(device)
	}

	riskLevel := "low"
	switch {
	case riskScore > 80:
		riskLevel = "very_high"
	case riskScore > 60:
		riskLevel = "high"
	case riskScore > 40:
		riskLevel = "medium"
	default:
		riskLevel = "low"
	}

	factors := []string{}
	if device.DeviceStatusInfo.IsStolen {
		factors = append(factors, "device_stolen")
	}
	if len(device.Claims) > 2 {
		factors = append(factors, "high_claim_frequency")
	}
	if string(device.DevicePhysicalCondition.Grade) == "D" || string(device.DevicePhysicalCondition.Grade) == "F" {
		factors = append(factors, "poor_condition")
	}
	if device.DeviceRiskAssessment.BlacklistStatus != "" {
		factors = append(factors, "blacklist_concern")
	}

	return &services.RiskAssessment{
		DeviceID:  deviceID,
		RiskScore: riskScore,
		RiskLevel: riskLevel,
		Factors:   factors,
		Timestamp: time.Now(),
	}, nil
}

// GetCurrentInsurableValue returns the current insurable value
func (s *deviceService) GetCurrentInsurableValue(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	// Use depreciation calculation for insurable value
	return s.CalculateDepreciation(ctx, deviceID)
}

// ============================================
// CLAIMS PROCESSING
// ============================================

// GetActiveClaimCount returns the number of active claims
func (s *deviceService) GetActiveClaimCount(ctx context.Context, deviceID uuid.UUID) (int, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, claim := range device.Claims {
		status := claim.ClaimLifecycle.Status
		if status == "active" || status == "pending" || status == "processing" {
			count++
		}
	}
	return count, nil
}

// GetTotalClaimAmount calculates total claim amount for this device
func (s *deviceService) GetTotalClaimAmount(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, claim := range device.Claims {
		status := claim.ClaimLifecycle.Status
		if status == "approved" || status == "paid" {
			total += claim.ClaimFinancial.ClaimedAmount
		}
	}
	return total, nil
}

// ============================================
// TRADE-IN AND VALUATION
// ============================================

// CalculateTradeInValue calculates estimated trade-in value
func (s *deviceService) CalculateTradeInValue(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Get base value from depreciation
	baseValue, err := s.CalculateDepreciation(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Apply condition multiplier
	conditionMultiplier := map[string]float64{
		"excellent": 0.9,
		"good":      0.75,
		"fair":      0.5,
		"poor":      0.25,
	}

	multiplier, exists := conditionMultiplier[string(device.DevicePhysicalCondition.Condition)]
	if !exists {
		multiplier = 0.5
	}

	tradeInValue := baseValue * multiplier

	// Reduce value if has many repairs
	repairCount := len(device.Repairs)
	if repairCount > 3 {
		tradeInValue *= 0.8
	}

	// Reduce if blacklisted
	if device.DeviceRiskAssessment.BlacklistStatus == "checking" {
		tradeInValue *= 0.7
	}

	return tradeInValue, nil
}

// CheckTradeInEligibility checks trade-in eligibility
func (s *deviceService) CheckTradeInEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return false, "", err
	}

	if device.DeviceStatusInfo.IsStolen {
		return false, "Device is reported stolen", nil
	}

	if device.DeviceRiskAssessment.BlacklistStatus == "blocked" {
		return false, "Device is blacklisted", nil
	}

	if string(device.DevicePhysicalCondition.Condition) == "poor" || string(device.DevicePhysicalCondition.Grade) == "F" {
		return false, "Device condition is too poor", nil
	}

	// Check if not in active rental or financing
	hasActiveRental := len(device.Rentals) > 0
	hasActiveFinancing := len(device.Financings) > 0

	if hasActiveRental || hasActiveFinancing {
		return false, "Device has active rental or financing", nil
	}

	return true, "Eligible for trade-in", nil
}

// ============================================
// RENTAL AND FINANCING
// ============================================

// CheckRentalEligibility checks rental eligibility
func (s *deviceService) CheckRentalEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return false, "", err
	}

	if device.DeviceStatusInfo.Status != "active" {
		return false, "Device is not active", nil
	}

	if device.DeviceStatusInfo.IsStolen {
		return false, "Device is reported stolen", nil
	}

	hasActiveRental := len(device.Rentals) > 0
	if hasActiveRental {
		return false, "Device already has active rental", nil
	}

	return true, "Eligible for rental", nil
}

// GetTotalRepairCost calculates total repair costs
func (s *deviceService) GetTotalRepairCost(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Return a calculated cost based on number of repairs
	// In production, this would sum actual repair costs from repair records
	return float64(len(device.Repairs)) * 100.0, nil // Estimate $100 per repair
}

// ============================================
// SUSTAINABILITY METRICS
// ============================================

// CalculateCarbonFootprint returns total carbon footprint
func (s *deviceService) CalculateCarbonFootprint(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Return estimate based on device age
	if device.DeviceFinancial.PurchaseDate != nil {
		age := time.Since(*device.DeviceFinancial.PurchaseDate).Hours() / (24 * 365) // years
		return age * 50.0, nil                                                       // 50kg CO2 per year estimate
	}
	return 50.0, nil
}

// GetRecyclingScore returns recycling score
func (s *deviceService) GetRecyclingScore(ctx context.Context, deviceID uuid.UUID) (int, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Return default score based on device age
	if device.DeviceFinancial.PurchaseDate != nil {
		age := time.Since(*device.DeviceFinancial.PurchaseDate).Hours() / (24 * 365)
		if age > 3 {
			return 60, nil // Older devices get lower scores
		}
	}
	return 80, nil // Newer devices get higher scores
}

// GetRepairabilityScore returns repairability score
func (s *deviceService) GetRepairabilityScore(ctx context.Context, deviceID uuid.UUID) (int, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Return default score based on repair count
	repairCount := len(device.Repairs)
	score := 70 - (repairCount * 10)
	if score < 0 {
		score = 0
	}
	return score, nil
}

// ============================================
// CUSTOMER EXPERIENCE
// ============================================

// GetSatisfactionScore returns average satisfaction score
func (s *deviceService) GetSatisfactionScore(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Return default satisfaction based on condition
	switch string(device.DevicePhysicalCondition.Condition) {
	case "excellent":
		return 90.0, nil
	case "good":
		return 75.0, nil
	case "fair":
		return 60.0, nil
	default:
		return 50.0, nil
	}
}

// ============================================
// HELPER METHODS
// ============================================

// calculateRiskScore calculates a risk score for the device
func (s *deviceService) calculateRiskScore(device *models.Device) float64 {
	riskScore := 0.0

	// Factor in condition
	switch string(device.DevicePhysicalCondition.Grade) {
	case "F":
		riskScore += 30
	case "D":
		riskScore += 20
	case "C":
		riskScore += 10
	}

	// Factor in claims
	riskScore += float64(len(device.Claims) * 10)

	// Factor in stolen status
	if device.DeviceStatusInfo.IsStolen {
		riskScore += 50
	}

	// Factor in blacklist status
	if device.DeviceRiskAssessment.BlacklistStatus != "" {
		riskScore += 20
	}

	// Factor in repair count
	riskScore += float64(len(device.Repairs) * 5)

	if riskScore > 100 {
		riskScore = 100
	}

	return riskScore
}

// Implement remaining interface methods...
// These are stubs that would need full implementation in production

func (s *deviceService) RegisterDevice(ctx context.Context, device *models.Device) error {
	return s.repo.Create(ctx, device)
}

func (s *deviceService) VerifyDevice(ctx context.Context, deviceID uuid.UUID) error {
	// Implementation needed
	return nil
}

func (s *deviceService) UpdateVerification(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return err
	}

	now := time.Now()
	device.DeviceVerification.VerificationDate = &now
	device.DeviceVerification.IsVerified = true

	return s.repo.Update(ctx, device)
}

func (s *deviceService) CalculatePremium(ctx context.Context, deviceID uuid.UUID) (*services.Premium, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) ProcessClaim(ctx context.Context, deviceID uuid.UUID, claimID uuid.UUID) error {
	// Implementation needed
	return nil
}

func (s *deviceService) ValidateClaimEligibility(ctx context.Context, deviceID uuid.UUID) (bool, string, error) {
	// Implementation needed
	return true, "", nil
}

func (s *deviceService) DetectFraud(ctx context.Context, deviceID uuid.UUID) (*services.FraudAnalysis, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) CheckBlacklistStatus(ctx context.Context, deviceID uuid.UUID) (bool, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return false, err
	}

	return device.DeviceRiskAssessment.BlacklistStatus == "blocked", nil
}

func (s *deviceService) ReportStolen(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return err
	}

	device.DeviceStatusInfo.IsStolen = true
	device.DeviceStatusInfo.Status = "stolen"
	now := time.Now()
	device.DeviceStatusInfo.StolenDate = &now

	return s.repo.Update(ctx, device)
}

func (s *deviceService) ProcessTradeIn(ctx context.Context, deviceID uuid.UUID, newDeviceID uuid.UUID) error {
	// Implementation needed
	return nil
}

func (s *deviceService) InitiateRental(ctx context.Context, deviceID uuid.UUID, duration time.Duration) error {
	// Implementation needed
	return nil
}

func (s *deviceService) InitiateFinancing(ctx context.Context, deviceID uuid.UUID, terms *services.FinancingTerms) error {
	// Implementation needed
	return nil
}

func (s *deviceService) ScheduleRepair(ctx context.Context, deviceID uuid.UUID, repairType string) error {
	// Implementation needed
	return nil
}

func (s *deviceService) GetRepairHistory(ctx context.Context, deviceID uuid.UUID) ([]*services.RepairRecord, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) AnalyzeUsagePatterns(ctx context.Context, deviceID uuid.UUID) (*services.UsageAnalysis, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) PredictFailure(ctx context.Context, deviceID uuid.UUID) (*services.FailurePrediction, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) GenerateUpgradeRecommendation(ctx context.Context, deviceID uuid.UUID) (*services.UpgradeRecommendation, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) CalculateBehaviorScore(ctx context.Context, deviceID uuid.UUID) (float64, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return 0, err
	}

	// Return risk score as proxy for behavior score
	return device.DeviceRiskAssessment.RiskScore, nil
}

func (s *deviceService) CheckCompliance(ctx context.Context, deviceID uuid.UUID) (*services.ComplianceReport, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) ApplyLegalHold(ctx context.Context, deviceID uuid.UUID, reason string) error {
	// Implementation needed
	return nil
}

func (s *deviceService) GenerateRegulatoryReport(ctx context.Context, deviceID uuid.UUID) (*services.RegulatoryReport, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) CheckEcoFriendlyCertification(ctx context.Context, deviceID uuid.UUID) (bool, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return false, err
	}

	return device.EcoLabel != nil, nil
}

func (s *deviceService) GenerateRecommendations(ctx context.Context, deviceID uuid.UUID) ([]*services.Recommendation, error) {
	// Implementation needed
	return nil, nil
}

func (s *deviceService) ProcessFeedback(ctx context.Context, deviceID uuid.UUID, feedback *services.CustomerFeedback) error {
	// Implementation needed
	return nil
}

func (s *deviceService) AssignToCorporateUser(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) error {
	// Implementation needed
	return nil
}

func (s *deviceService) RegisterBYOD(ctx context.Context, deviceID uuid.UUID, programID uuid.UUID) error {
	// Implementation needed
	return nil
}

func (s *deviceService) AddToPool(ctx context.Context, deviceID uuid.UUID, poolID uuid.UUID) error {
	// Implementation needed
	return nil
}

func (s *deviceService) GetCorporateDeviceStatus(ctx context.Context, deviceID uuid.UUID) (*services.CorporateStatus, error) {
	device, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	status := &services.CorporateStatus{
		IsCorporate: device.DeviceOwnership.CorporateAccountID != nil || device.CorporateAssignment != nil,
		IsBYOD:      device.BYODRegistration != nil,
	}

	if device.CorporateAssignment != nil {
		status.AssignedTo = &device.DeviceOwnership.OwnerID
	}

	return status, nil
}
