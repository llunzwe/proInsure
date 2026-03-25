package device

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"smartsure/internal/domain/models/device"
	domainServices "smartsure/internal/domain/services/device"
)

// DeviceEcosystemService orchestrates all device ecosystem operations
// This is the main APPLICATION SERVICE for device management
type DeviceEcosystemService struct {
	db                    *gorm.DB
	fraudService          *domainServices.FraudDetectionService
	valuationService      *domainServices.ValuationService
	qualityService        *domainServices.QualityService
	performanceService    *domainServices.PerformanceService
	sustainabilityService *domainServices.SustainabilityService
	biService             *domainServices.BusinessIntelligenceService
	accessoriesService    *AccessoriesService
	sparePartsService     *SparePartsService
}

// NewDeviceEcosystemService creates a new device ecosystem service
func NewDeviceEcosystemService(db *gorm.DB) *DeviceEcosystemService {
	return &DeviceEcosystemService{
		db:                    db,
		fraudService:          domainServices.NewFraudDetectionService(),
		valuationService:      domainServices.NewValuationService(),
		qualityService:        domainServices.NewQualityService(),
		performanceService:    domainServices.NewPerformanceService(),
		sustainabilityService: domainServices.NewSustainabilityService(),
		biService:             domainServices.NewBusinessIntelligenceService(),
		accessoriesService:    NewAccessoriesService(db),
		sparePartsService:     NewSparePartsService(db),
	}
}

// ===== Core Device Management =====

// RegisterDevice registers a new device in the system
func (s *DeviceEcosystemService) RegisterDevice(ctx context.Context, req *RegisterDeviceRequest) (*device.Device, error) {
	// Validate IMEI using domain service
	if !s.fraudService.ValidateIMEI(req.IMEI) {
		return nil, errors.New("invalid IMEI")
	}
	
	// Check if device already exists
	var existingDevice device.Device
	if err := s.db.WithContext(ctx).Where("imei = ?", req.IMEI).First(&existingDevice).Error; err == nil {
		return nil, errors.New("device with this IMEI already exists")
	}
	
	// Create new device
	device := &device.Device{
		IMEI:            req.IMEI,
		SerialNumber:    req.SerialNumber,
		Model:           req.Model,
		Brand:           req.Brand,
		Manufacturer:    req.Manufacturer,
		OperatingSystem: req.OperatingSystem,
		OSVersion:       req.OSVersion,
		StorageCapacity: req.StorageCapacity,
		RAM:             req.RAM,
		Color:           req.Color,
		PurchaseDate:    req.PurchaseDate,
		PurchasePrice:   req.PurchasePrice,
		OwnerID:         req.OwnerID,
		Status:          "active",
		Condition:       "good",
		IsVerified:      false,
	}
	
	// Calculate initial values
	device.CurrentValue = device.CalculateDepreciation()
	device.UpdateMarketValue()
	
	// Save to database
	if err := s.db.WithContext(ctx).Create(device).Error; err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}
	
	// Initialize lifecycle tracking
	lifecycle := &device.DeviceLifecycle{
		DeviceID:        device.ID,
		Stage:           "registered",
		AcquisitionDate: time.Now(),
		AcquisitionCost: req.PurchasePrice,
	}
	s.db.WithContext(ctx).Create(lifecycle)
	
	return device, nil
}

// GetDevice retrieves a device with all related data
func (s *DeviceEcosystemService) GetDevice(ctx context.Context, deviceID string) (*device.Device, error) {
	var device device.Device
	err := s.db.WithContext(ctx).
		Preload("Swaps").
		Preload("TradeIns").
		Preload("Refurbishments").
		Preload("Rentals").
		Preload("Financings").
		Preload("Layaways").
		Preload("Repairs").
		Preload("Lifecycle").
		Preload("Subscriptions").
		Preload("MarketListings").
		Preload("Accessories").
		Preload("SpareParts").
		First(&device, "id = ?", deviceID).Error
	
	if err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}
	
	return &device, nil
}

// VerifyDevice performs device verification and authentication
func (s *DeviceEcosystemService) VerifyDevice(ctx context.Context, deviceID string, req *VerifyDeviceRequest) error {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return err
	}
	
	// Use domain service for fraud checks
	riskScore := s.fraudService.CalculateRiskScore(device)
	if riskScore > 70 {
		device.AuthenticityStatus = "suspicious"
		device.FraudRiskScore = riskScore
		s.db.WithContext(ctx).Save(device)
		return errors.New("device failed fraud verification")
	}
	
	// Update verification status
	device.UpdateVerification()
	device.AuthenticityStatus = "verified"
	device.ProofOfOwnershipVerified = req.ProofOfOwnership
	
	if req.VideoVerificationURL != "" {
		device.VideoVerificationURL = req.VideoVerificationURL
	}
	
	return s.db.WithContext(ctx).Save(device).Error
}

// ===== Device Swap Management =====

// InitiateDeviceSwap starts a device swap/upgrade process
func (s *DeviceEcosystemService) InitiateDeviceSwap(ctx context.Context, oldDeviceID, newDeviceID string) (*device.DeviceSwap, error) {
	// Load both devices
	oldDevice, err := s.GetDevice(ctx, oldDeviceID)
	if err != nil {
		return nil, fmt.Errorf("old device not found: %w", err)
	}
	
	newDevice, err := s.GetDevice(ctx, newDeviceID)
	if err != nil {
		return nil, fmt.Errorf("new device not found: %w", err)
	}
	
	// Check eligibility
	if !oldDevice.IsEligibleForTradeIn() {
		return nil, errors.New("old device not eligible for swap")
	}
	
	// Create swap record
	swap := &device.DeviceSwap{
		DeviceID:         uuid.MustParse(oldDeviceID),
		NewDeviceID:      uuid.MustParse(newDeviceID),
		SwapDate:         time.Now(),
		SwapValue:        oldDevice.MarketValue,
		SwapReason:       "upgrade",
		SwapStatus:       "initiated",
		EligibilityScore: oldDevice.GetConditionScore() * 100,
		LoyaltyPoints:    int(oldDevice.MarketValue / 10), // 1 point per $10 value
		SwapFee:          0,                               // Can be calculated based on business rules
		IsUpgrade:        newDevice.PurchasePrice > oldDevice.PurchasePrice,
	}
	
	if err := s.db.WithContext(ctx).Create(swap).Error; err != nil {
		return nil, fmt.Errorf("failed to create swap: %w", err)
	}
	
	return swap, nil
}

// ===== Trade-In Management =====

// ProcessTradeIn handles device trade-in valuation and processing
func (s *DeviceEcosystemService) ProcessTradeIn(ctx context.Context, deviceID string) (*device.DeviceTradeIn, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	if !device.IsEligibleForTradeIn() {
		return nil, errors.New("device not eligible for trade-in")
	}
	
	// Use domain service for valuation
	tradeInValue := s.valuationService.CalculateTradeInValue(device)
	
	tradeIn := &device.DeviceTradeIn{
		DeviceID:       device.ID,
		TradeInDate:    time.Now(),
		OfferedValue:   tradeInValue,
		AcceptedValue:  tradeInValue,
		Status:         "pending",
		ConditionGrade: device.Grade,
		ConditionNotes: device.InspectionNotes,
		TradeInPartner: "SmartSure",
		InstantTradeIn: false,
		TradeInMethod:  "mail_in",
	}
	
	// Apply deductions based on condition
	deductions := s.calculateTradeInDeductions(device)
	tradeIn.Deductions = deductions
	tradeIn.AcceptedValue -= deductions
	
	if err := s.db.WithContext(ctx).Create(tradeIn).Error; err != nil {
		return nil, fmt.Errorf("failed to process trade-in: %w", err)
	}
	
	return tradeIn, nil
}

// ===== Rental Management =====

// CreateRentalAgreement creates a new device rental agreement
func (s *DeviceEcosystemService) CreateRentalAgreement(ctx context.Context, deviceID string, req *RentalRequest) (*device.DeviceRental, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	if !device.IsEligibleForRental() {
		return nil, errors.New("device not eligible for rental")
	}
	
	rental := &device.DeviceRental{
		DeviceID:        device.ID,
		RenterID:        req.RenterID,
		RentalStartDate: req.StartDate,
		RentalEndDate:   req.EndDate,
		MonthlyRate:     s.calculateRentalRate(device),
		DepositAmount:   device.MarketValue * 0.2, // 20% deposit
		RentalStatus:    "active",
		RentalType:      req.RentalType,
		ContractTerms:   req.ContractTerms,
	}
	
	// Set rent-to-own details if applicable
	if req.RentalType == "rent_to_own" {
		rental.RentToOwn = true
		rental.RentToOwnPrice = device.MarketValue * 1.3 // 30% markup for rent-to-own
		rental.RentToOwnMonths = 24                      // Standard 24-month plan
	}
	
	if err := s.db.WithContext(ctx).Create(rental).Error; err != nil {
		return nil, fmt.Errorf("failed to create rental: %w", err)
	}
	
	// Update device status
	device.Status = "rented"
	s.db.WithContext(ctx).Save(device)
	
	return rental, nil
}

// ===== Financing Management =====

// CreateFinancingPlan creates a device financing/credit plan
func (s *DeviceEcosystemService) CreateFinancingPlan(ctx context.Context, deviceID string, req *FinancingRequest) (*device.DeviceFinancing, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	// Calculate financing terms
	interestRate := s.calculateInterestRate(req.CreditScore)
	monthlyPayment := s.calculateMonthlyPayment(device.PurchasePrice, req.DownPayment, interestRate, req.TermMonths)
	
	financing := &device.DeviceFinancing{
		DeviceID:          device.ID,
		CustomerID:        req.CustomerID,
		FinanceProvider:   req.Provider,
		PrincipalAmount:   device.PurchasePrice - req.DownPayment,
		InterestRate:      interestRate,
		TermMonths:        req.TermMonths,
		MonthlyPayment:    monthlyPayment,
		DownPayment:       req.DownPayment,
		TotalPayable:      monthlyPayment*float64(req.TermMonths) + req.DownPayment,
		FinanceStatus:     "active",
		CreditScore:       req.CreditScore,
		InsuranceRequired: true, // Require insurance for financed devices
	}
	
	if err := s.db.WithContext(ctx).Create(financing).Error; err != nil {
		return nil, fmt.Errorf("failed to create financing: %w", err)
	}
	
	return financing, nil
}

// ===== Repair Management =====

// CreateRepairOrder creates a new repair order for a device
func (s *DeviceEcosystemService) CreateRepairOrder(ctx context.Context, deviceID string, req *RepairRequest) (*device.DeviceRepair, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	// Estimate repair cost using domain service
	estimatedCost := s.valuationService.EstimateRepairCost(device, req.RepairType)
	
	repair := &device.DeviceRepair{
		DeviceID:         device.ID,
		RepairType:       req.RepairType,
		IssueDescription: req.Description,
		RepairStatus:     "pending",
		EstimatedCost:    estimatedCost,
		RepairProvider:   req.Provider,
		WarrantyRepair:   device.IsWarrantyActive(),
		InsuranceClaim:   req.UseInsurance,
	}
	
	// Check if loaner device needed
	if req.NeedLoanerDevice && device.LoanerDeviceEligible {
		repair.LoanerDeviceID = s.assignLoanerDevice(ctx)
	}
	
	if err := s.db.WithContext(ctx).Create(repair).Error; err != nil {
		return nil, fmt.Errorf("failed to create repair order: %w", err)
	}
	
	return repair, nil
}

// ===== Subscription Management =====

// CreateDeviceSubscription creates a new device subscription
func (s *DeviceEcosystemService) CreateDeviceSubscription(ctx context.Context, deviceID string, req *SubscriptionRequest) (*device.DeviceSubscription, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	subscription := &device.DeviceSubscription{
		DeviceID:        device.ID,
		SubscriberID:    req.SubscriberID,
		PlanType:        req.PlanType,
		PlanName:        req.PlanName,
		MonthlyFee:      req.MonthlyFee,
		StartDate:       time.Now(),
		Status:          "active",
		BillingCycle:    "monthly",
		AutoRenew:       true,
		BundledServices: req.BundledServices,
	}
	
	// Set end date based on plan type
	switch req.PlanType {
	case "monthly":
		endDate := time.Now().AddDate(0, 1, 0)
		subscription.EndDate = &endDate
	case "annual":
		endDate := time.Now().AddDate(1, 0, 0)
		subscription.EndDate = &endDate
	}
	
	if err := s.db.WithContext(ctx).Create(subscription).Error; err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}
	
	return subscription, nil
}

// ===== Marketplace Management =====

// ListDeviceForSale lists a device in the marketplace
func (s *DeviceEcosystemService) ListDeviceForSale(ctx context.Context, deviceID string, req *ListingRequest) (*device.DeviceMarketplace, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	// Check if device can be listed
	if device.IsListedForSale() {
		return nil, errors.New("device already listed for sale")
	}
	
	if device.HasActiveRental() || device.HasActiveFinancing() {
		return nil, errors.New("device has active obligations")
	}
	
	listing := &device.DeviceMarketplace{
		DeviceID:     device.ID,
		SellerID:     device.OwnerID,
		ListingPrice: req.Price,
		ListingDate:  time.Now(),
		Status:       "active",
		Platform:     req.Platform,
		ListingTitle: fmt.Sprintf("%s %s - %s", device.Brand, device.Model, device.Condition),
		Description:  req.Description,
		Photos:       req.Photos,
		ShippingType: req.ShippingType,
	}
	
	if err := s.db.WithContext(ctx).Create(listing).Error; err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}
	
	return listing, nil
}

// ===== Accessories Management =====

// RegisterDeviceAccessory adds an accessory to a device
func (s *DeviceEcosystemService) RegisterDeviceAccessory(ctx context.Context, deviceID string, req *RegisterAccessoryRequest) (*device.DeviceAccessory, error) {
	return s.accessoriesService.RegisterAccessory(ctx, deviceID, req)
}

// GetDeviceAccessories retrieves all accessories for a device
func (s *DeviceEcosystemService) GetDeviceAccessories(ctx context.Context, deviceID string) ([]device.DeviceAccessory, error) {
	return s.accessoriesService.GetDeviceAccessories(ctx, deviceID, false)
}

// ProcessAccessoryReplacement handles accessory replacement claims
func (s *DeviceEcosystemService) ProcessAccessoryReplacement(ctx context.Context, accessoryID string) (*device.DeviceAccessory, error) {
	return s.accessoriesService.ProcessAccessoryReplacement(ctx, accessoryID)
}

// CalculateAccessoriesInsurance calculates insurance value for accessories
func (s *DeviceEcosystemService) CalculateAccessoriesInsurance(ctx context.Context, deviceID string) (*AccessoriesInsuranceValue, error) {
	return s.accessoriesService.CalculateAccessoriesInsuranceValue(ctx, deviceID)
}

// ===== Spare Parts Management =====

// AddSparePart adds a new spare part to inventory
func (s *DeviceEcosystemService) AddSparePart(ctx context.Context, deviceID string, req *AddSparePartRequest) (*device.DeviceSparePart, error) {
	return s.sparePartsService.AddSparePart(ctx, deviceID, req)
}

// GetAvailableSpareParts retrieves available spare parts for a device
func (s *DeviceEcosystemService) GetAvailableSpareParts(ctx context.Context, deviceID string, partType string) ([]device.DeviceSparePart, error) {
	return s.sparePartsService.GetAvailableParts(ctx, deviceID, partType)
}

// UseSparePart records usage of a spare part in a repair
func (s *DeviceEcosystemService) UseSparePart(ctx context.Context, partID string, repairID string, quantity int) error {
	return s.sparePartsService.UseSparePart(ctx, partID, repairID, quantity)
}

// GetSparePartsInventoryValue calculates total inventory value
func (s *DeviceEcosystemService) GetSparePartsInventoryValue(ctx context.Context, deviceID string) (*InventoryValue, error) {
	return s.sparePartsService.GetInventoryValue(ctx, deviceID)
}

// CreateSparePartReorderRequest creates reorder request for low stock
func (s *DeviceEcosystemService) CreateSparePartReorderRequest(ctx context.Context, partID string) (*ReorderRequest, error) {
	return s.sparePartsService.CreateReorderRequest(ctx, partID)
}

// GetComprehensiveDeviceValue calculates total device value including accessories and parts
func (s *DeviceEcosystemService) GetComprehensiveDeviceValue(ctx context.Context, deviceID string) (*ComprehensiveDeviceValue, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	
	// Get accessories value
	accessoriesValue, err := s.CalculateAccessoriesInsurance(ctx, deviceID)
	if err != nil {
		// Log error but don't fail
		fmt.Printf("Failed to get accessories value: %v\n", err)
		accessoriesValue = &AccessoriesInsuranceValue{}
	}
	
	// Get spare parts inventory value
	partsInventory, err := s.GetSparePartsInventoryValue(ctx, deviceID)
	if err != nil {
		// Log error but don't fail
		fmt.Printf("Failed to get spare parts value: %v\n", err)
		partsInventory = &InventoryValue{}
	}
	
	totalValue := &ComprehensiveDeviceValue{
		DeviceID:         deviceID,
		DeviceValue:      device.CurrentValue,
		AccessoriesValue: accessoriesValue.TotalValue,
		SparePartsValue:  partsInventory.TotalValue,
		TotalValue:       device.CurrentValue + accessoriesValue.TotalValue + partsInventory.TotalValue,
		ReplacementCost:  device.CalculateReplacementCost(),
		AccessoriesCount: device.GetInsuredAccessoriesCount(),
		SparePartsCount:  partsInventory.UniquePartsCount,
		HasCriticalParts: device.HasCriticalSpareParts(),
		CalculatedAt:     time.Now(),
	}
	
	return totalValue, nil
}

// ===== Helper Methods =====

func (s *DeviceEcosystemService) calculateTradeInDeductions(Device *models.Device) float64 {
	deductions := 0.0
	
	if device.ScreenCondition == "cracked" {
		deductions += 50
	} else if device.ScreenCondition == "broken" {
		deductions += 100
	}
	
	if device.BodyCondition == "damaged" {
		deductions += 30
	}
	
	if device.BatteryHealth < 80 && device.BatteryHealth > 0 {
		deductions += 20
	}
	
	if device.WaterDamageIndicator != "white" {
		deductions += 75
	}
	
	return deductions
}

func (s *DeviceEcosystemService) calculateRentalRate(Device *models.Device) float64 {
	// Base rate is 5% of market value per month
	baseRate := device.MarketValue * 0.05
	
	// Adjust for device segment
	switch device.DeviceSegment {
	case "flagship":
		baseRate *= 1.2
	case "premium":
		baseRate *= 1.1
	case "budget":
		baseRate *= 0.9
	}
	
	return baseRate
}

func (s *DeviceEcosystemService) calculateInterestRate(creditScore int) float64 {
	if creditScore >= 750 {
		return 0.05 // 5% APR
	} else if creditScore >= 650 {
		return 0.08 // 8% APR
	} else if creditScore >= 550 {
		return 0.12 // 12% APR
	}
	return 0.15 // 15% APR
}

func (s *DeviceEcosystemService) calculateMonthlyPayment(principal, downPayment, interestRate float64, termMonths int) float64 {
	loanAmount := principal - downPayment
	monthlyRate := interestRate / 12
	
	if monthlyRate == 0 {
		return loanAmount / float64(termMonths)
	}
	
	// PMT formula
	payment := loanAmount * (monthlyRate*(1+monthlyRate) ^ float64(termMonths)) /
		((1 + monthlyRate) ^ float64(termMonths) - 1)
	
	return payment
}

func (s *DeviceEcosystemService) assignLoanerDevice(ctx context.Context) *uuid.UUID {
	// In production, this would find an available loaner device
	// For now, return nil
	return nil
}

// Request structures

type RegisterDeviceRequest struct {
	IMEI            string     `json:"imei"`
	SerialNumber    string     `json:"serial_number"`
	Model           string     `json:"model"`
	Brand           string     `json:"brand"`
	Manufacturer    string     `json:"manufacturer"`
	OperatingSystem string     `json:"operating_system"`
	OSVersion       string     `json:"os_version"`
	StorageCapacity int        `json:"storage_capacity"`
	RAM             int        `json:"ram"`
	Color           string     `json:"color"`
	PurchaseDate    *time.Time `json:"purchase_date"`
	PurchasePrice   float64    `json:"purchase_price"`
	OwnerID         uuid.UUID  `json:"owner_id"`
}

type VerifyDeviceRequest struct {
	ProofOfOwnership     bool   `json:"proof_of_ownership"`
	VideoVerificationURL string `json:"video_verification_url"`
}

type RentalRequest struct {
	RenterID      uuid.UUID `json:"renter_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	RentalType    string    `json:"rental_type"`
	ContractTerms string    `json:"contract_terms"`
}

type FinancingRequest struct {
	CustomerID  uuid.UUID `json:"customer_id"`
	Provider    string    `json:"provider"`
	DownPayment float64   `json:"down_payment"`
	TermMonths  int       `json:"term_months"`
	CreditScore int       `json:"credit_score"`
}

type RepairRequest struct {
	RepairType       string `json:"repair_type"`
	Description      string `json:"description"`
	Provider         string `json:"provider"`
	UseInsurance     bool   `json:"use_insurance"`
	NeedLoanerDevice bool   `json:"need_loaner_device"`
}

type SubscriptionRequest struct {
	SubscriberID    uuid.UUID `json:"subscriber_id"`
	PlanType        string    `json:"plan_type"`
	PlanName        string    `json:"plan_name"`
	MonthlyFee      float64   `json:"monthly_fee"`
	BundledServices string    `json:"bundled_services"`
}

type ListingRequest struct {
	Price        float64 `json:"price"`
	Platform     string  `json:"platform"`
	Description  string  `json:"description"`
	Photos       string  `json:"photos"`
	ShippingType string  `json:"shipping_type"`
}

type ComprehensiveDeviceValue struct {
	DeviceID         string    `json:"device_id"`
	DeviceValue      float64   `json:"device_value"`
	AccessoriesValue float64   `json:"accessories_value"`
	SparePartsValue  float64   `json:"spare_parts_value"`
	TotalValue       float64   `json:"total_value"`
	ReplacementCost  float64   `json:"replacement_cost"`
	AccessoriesCount int       `json:"accessories_count"`
	SparePartsCount  int       `json:"spare_parts_count"`
	HasCriticalParts bool      `json:"has_critical_parts"`
	CalculatedAt     time.Time `json:"calculated_at"`
}
