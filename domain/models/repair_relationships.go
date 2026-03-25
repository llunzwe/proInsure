package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"smartsure/pkg/database"
)

// RepairPartsUsed tracks spare parts consumed during a repair
type RepairPartsUsed struct {
	database.BaseModel
	RepairBookingID uuid.UUID `gorm:"type:uuid;not null;index" json:"repair_booking_id"`
	SparePartID     uuid.UUID `gorm:"type:uuid;not null;index" json:"spare_part_id"`
	TechnicianID    uuid.UUID `gorm:"type:uuid;not null" json:"technician_id"`

	// Part details at time of use (for historical tracking)
	PartNumber   string `gorm:"not null" json:"part_number"`
	PartName     string `gorm:"not null" json:"part_name"`
	PartType     string `gorm:"not null" json:"part_type"`
	Quality      string `json:"quality"` // oem, genuine, aftermarket
	SerialNumber string `json:"serial_number,omitempty"`
	BatchNumber  string `json:"batch_number,omitempty"`

	// Quantity and pricing
	QuantityUsed   int             `gorm:"not null;default:1" json:"quantity_used"`
	UnitCost       decimal.Decimal `gorm:"type:decimal(10,2)" json:"unit_cost"`
	RetailPrice    decimal.Decimal `gorm:"type:decimal(10,2)" json:"retail_price"`
	ChargedPrice   decimal.Decimal `gorm:"type:decimal(10,2)" json:"charged_price"`
	DiscountAmount decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	TotalCost      decimal.Decimal `gorm:"type:decimal(10,2)" json:"total_cost"`

	// Warranty and tracking
	WarrantyDays    int        `json:"warranty_days"`
	WarrantyExpiry  *time.Time `json:"warranty_expiry,omitempty"`
	IsWarrantyClaim bool       `gorm:"default:false" json:"is_warranty_claim"`
	IsReturnable    bool       `gorm:"default:true" json:"is_returnable"`
	ReturnDeadline  *time.Time `json:"return_deadline,omitempty"`

	// Installation details
	InstalledAt       time.Time `json:"installed_at"`
	InstallationNotes string    `json:"installation_notes,omitempty"`
	TestResults       string    `json:"test_results,omitempty"`
	IsSuccessful      bool      `gorm:"default:true" json:"is_successful"`
	FailureReason     string    `json:"failure_reason,omitempty"`

	// Source tracking
	SourceType      string     `json:"source_type"` // inventory, external, customer_provided
	SourceLocation  string     `json:"source_location,omitempty"`
	PurchaseOrderID *uuid.UUID `gorm:"type:uuid" json:"purchase_order_id,omitempty"`

	// Quality control
	QualityChecked bool    `gorm:"default:false" json:"quality_checked"`
	QualityScore   float64 `json:"quality_score,omitempty"`
	DefectsFound   string  `json:"defects_found,omitempty"` // JSON array

	// Relationships
	RepairBooking RepairBooking `gorm:"foreignKey:RepairBookingID;constraint:OnDelete:CASCADE" json:"repair_booking,omitempty"`
	SparePart     SparePart     `gorm:"foreignKey:SparePartID" json:"spare_part,omitempty"`
	Technician    Technician    `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
}

// RepairAccessoriesReplaced tracks accessories replaced or repaired during service
type RepairAccessoriesReplaced struct {
	database.BaseModel
	RepairBookingID uuid.UUID  `gorm:"type:uuid;not null;index" json:"repair_booking_id"`
	OldAccessoryID  *uuid.UUID `gorm:"type:uuid;index" json:"old_accessory_id,omitempty"`
	NewAccessoryID  *uuid.UUID `gorm:"type:uuid;index" json:"new_accessory_id,omitempty"`
	TechnicianID    uuid.UUID  `gorm:"type:uuid;not null" json:"technician_id"`

	// Action taken
	ActionType   string `gorm:"not null" json:"action_type"` // replaced, repaired, removed, added
	ActionReason string `json:"action_reason"`

	// Accessory details
	AccessoryType string `gorm:"not null" json:"accessory_type"`
	Brand         string `json:"brand"`
	Model         string `json:"model"`
	SerialNumber  string `json:"serial_number,omitempty"`

	// Old accessory condition
	OldCondition      string `json:"old_condition,omitempty"`
	OldConditionNotes string `json:"old_condition_notes,omitempty"`
	OldPhotos         string `json:"old_photos,omitempty"` // JSON array of URLs
	IsOldReturned     bool   `gorm:"default:false" json:"is_old_returned"`
	OldDisposalMethod string `json:"old_disposal_method,omitempty"`

	// New accessory details (if replaced)
	NewCondition      string     `json:"new_condition,omitempty"`
	NewWarrantyDays   int        `json:"new_warranty_days,omitempty"`
	NewWarrantyExpiry *time.Time `json:"new_warranty_expiry,omitempty"`

	// Pricing
	ReplacementCost decimal.Decimal `gorm:"type:decimal(10,2)" json:"replacement_cost"`
	LaborCost       decimal.Decimal `gorm:"type:decimal(10,2)" json:"labor_cost"`
	ChargedPrice    decimal.Decimal `gorm:"type:decimal(10,2)" json:"charged_price"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	TotalCost       decimal.Decimal `gorm:"type:decimal(10,2)" json:"total_cost"`

	// Insurance and warranty
	IsCoveredByInsurance  bool       `gorm:"default:false" json:"is_covered_by_insurance"`
	InsuranceClaimID      *uuid.UUID `gorm:"type:uuid" json:"insurance_claim_id,omitempty"`
	IsWarrantyReplacement bool       `gorm:"default:false" json:"is_warranty_replacement"`
	WarrantyClaimID       *uuid.UUID `gorm:"type:uuid" json:"warranty_claim_id,omitempty"`

	// Quality control
	TestPerformed    bool       `gorm:"default:false" json:"test_performed"`
	TestResults      string     `json:"test_results,omitempty"`
	CustomerApproved bool       `gorm:"default:false" json:"customer_approved"`
	ApprovalDate     *time.Time `json:"approval_date,omitempty"`

	// Timestamps
	ReplacedAt  time.Time  `json:"replaced_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Relationships
	RepairBooking RepairBooking `gorm:"foreignKey:RepairBookingID;constraint:OnDelete:CASCADE" json:"repair_booking,omitempty"`
	OldAccessory  *Accessory    `gorm:"foreignKey:OldAccessoryID" json:"old_accessory,omitempty"`
	NewAccessory  *Accessory    `gorm:"foreignKey:NewAccessoryID" json:"new_accessory,omitempty"`
	Technician    Technician    `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
}

// RepairTechnicianAssignment tracks which technicians worked on which repairs
type RepairTechnicianAssignment struct {
	database.BaseModel
	RepairBookingID uuid.UUID `gorm:"type:uuid;not null;index" json:"repair_booking_id"`
	TechnicianID    uuid.UUID `gorm:"type:uuid;not null;index" json:"technician_id"`
	RepairShopID    uuid.UUID `gorm:"type:uuid;not null;index" json:"repair_shop_id"`

	// Assignment details
	AssignmentType string     `gorm:"not null" json:"assignment_type"` // primary, secondary, specialist, supervisor
	AssignedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"assigned_by"`
	AssignedAt     time.Time  `gorm:"not null" json:"assigned_at"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`

	// Work performed
	TasksPerformed string          `json:"tasks_performed"` // JSON array
	HoursWorked    decimal.Decimal `gorm:"type:decimal(5,2)" json:"hours_worked"`
	LaborRate      decimal.Decimal `gorm:"type:decimal(10,2)" json:"labor_rate"`
	TotalLaborCost decimal.Decimal `gorm:"type:decimal(10,2)" json:"total_labor_cost"`

	// Performance tracking
	QualityScore    float64 `json:"quality_score,omitempty"`
	EfficiencyScore float64 `json:"efficiency_score,omitempty"`
	CustomerRating  float64 `json:"customer_rating,omitempty"`
	Notes           string  `json:"notes,omitempty"`

	// Status
	Status         string `gorm:"default:'assigned'" json:"status"` // assigned, in_progress, completed, reassigned
	ReassignReason string `json:"reassign_reason,omitempty"`

	// Relationships
	RepairBooking RepairBooking `gorm:"foreignKey:RepairBookingID;constraint:OnDelete:CASCADE" json:"repair_booking,omitempty"`
	Technician    Technician    `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	RepairShop    RepairShop    `gorm:"foreignKey:RepairShopID" json:"repair_shop,omitempty"`
}

// ============== Table Names ==============

func (r *RepairPartsUsed) TableName() string {
	return "repair_parts_used"
}

func (r *RepairAccessoriesReplaced) TableName() string {
	return "repair_accessories_replaced"
}

func (r *RepairTechnicianAssignment) TableName() string {
	return "repair_technician_assignments"
}

// ============== Hooks ==============

// BeforeCreate hook for RepairPartsUsed
func (r *RepairPartsUsed) BeforeCreate(tx *gorm.DB) error {
	if err := r.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Calculate total cost
	r.TotalCost = r.ChargedPrice.Mul(decimal.NewFromInt(int64(r.QuantityUsed))).Sub(r.DiscountAmount)

	// Set warranty expiry if warranty days provided
	if r.WarrantyDays > 0 {
		expiry := r.InstalledAt.AddDate(0, 0, r.WarrantyDays)
		r.WarrantyExpiry = &expiry
	}

	// Set return deadline (typically 30 days)
	if r.IsReturnable {
		deadline := r.InstalledAt.AddDate(0, 0, 30)
		r.ReturnDeadline = &deadline
	}

	return nil
}

// BeforeCreate hook for RepairAccessoriesReplaced
func (r *RepairAccessoriesReplaced) BeforeCreate(tx *gorm.DB) error {
	if err := r.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Calculate total cost
	r.TotalCost = r.ReplacementCost.Add(r.LaborCost).Sub(r.DiscountAmount)

	// Set warranty expiry for new accessories
	if r.NewWarrantyDays > 0 {
		expiry := r.ReplacedAt.AddDate(0, 0, r.NewWarrantyDays)
		r.NewWarrantyExpiry = &expiry
	}

	return nil
}

// BeforeCreate hook for RepairTechnicianAssignment
func (r *RepairTechnicianAssignment) BeforeCreate(tx *gorm.DB) error {
	if err := r.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Calculate total labor cost
	r.TotalLaborCost = r.LaborRate.Mul(r.HoursWorked)

	return nil
}

// ============== Business Logic Methods for RepairPartsUsed ==============

// Validate validates the repair parts used record
func (r *RepairPartsUsed) Validate() error {
	if r.RepairBookingID == uuid.Nil {
		return errors.New("repair booking ID is required")
	}
	if r.SparePartID == uuid.Nil {
		return errors.New("spare part ID is required")
	}
	if r.TechnicianID == uuid.Nil {
		return errors.New("technician ID is required")
	}
	if r.QuantityUsed <= 0 {
		return errors.New("quantity used must be greater than 0")
	}
	if r.UnitCost.IsNegative() || r.RetailPrice.IsNegative() || r.ChargedPrice.IsNegative() {
		return errors.New("prices cannot be negative")
	}
	return nil
}

// CalculateMargin calculates the profit margin on the part
func (r *RepairPartsUsed) CalculateMargin() decimal.Decimal {
	if r.UnitCost.IsZero() {
		return decimal.Zero
	}
	margin := r.ChargedPrice.Sub(r.UnitCost).Div(r.ChargedPrice).Mul(decimal.NewFromInt(100))
	return margin
}

// IsUnderWarranty checks if the part is still under warranty
func (r *RepairPartsUsed) IsUnderWarranty() bool {
	if r.WarrantyExpiry == nil {
		return false
	}
	return time.Now().Before(*r.WarrantyExpiry)
}

// CanReturn checks if the part can still be returned
func (r *RepairPartsUsed) CanReturn() bool {
	if !r.IsReturnable || r.ReturnDeadline == nil {
		return false
	}
	return time.Now().Before(*r.ReturnDeadline)
}

// MarkDefective marks the part as defective
func (r *RepairPartsUsed) MarkDefective(reason string) {
	r.IsSuccessful = false
	r.FailureReason = reason
	r.QualityScore = 0
}

// ============== Business Logic Methods for RepairAccessoriesReplaced ==============

// Validate validates the repair accessories replaced record
func (r *RepairAccessoriesReplaced) Validate() error {
	if r.RepairBookingID == uuid.Nil {
		return errors.New("repair booking ID is required")
	}
	if r.TechnicianID == uuid.Nil {
		return errors.New("technician ID is required")
	}
	if r.ActionType == "" {
		return errors.New("action type is required")
	}
	if r.AccessoryType == "" {
		return errors.New("accessory type is required")
	}

	// Validate action type
	validActions := map[string]bool{
		"replaced": true,
		"repaired": true,
		"removed":  true,
		"added":    true,
	}
	if !validActions[r.ActionType] {
		return errors.New("invalid action type")
	}

	// If replacing, new accessory ID should be provided
	if r.ActionType == "replaced" && r.NewAccessoryID == nil {
		return errors.New("new accessory ID required for replacement")
	}

	return nil
}

// CalculateTotalCost recalculates the total cost
func (r *RepairAccessoriesReplaced) CalculateTotalCost() {
	r.TotalCost = r.ReplacementCost.Add(r.LaborCost).Sub(r.DiscountAmount)
}

// IsNewUnderWarranty checks if the new accessory is under warranty
func (r *RepairAccessoriesReplaced) IsNewUnderWarranty() bool {
	if r.NewWarrantyExpiry == nil {
		return false
	}
	return time.Now().Before(*r.NewWarrantyExpiry)
}

// ApproveByCustomer marks the replacement as customer approved
func (r *RepairAccessoriesReplaced) ApproveByCustomer() {
	r.CustomerApproved = true
	now := time.Now()
	r.ApprovalDate = &now
}

// Complete marks the accessory replacement as completed
func (r *RepairAccessoriesReplaced) Complete() {
	now := time.Now()
	r.CompletedAt = &now
}

// ============== Business Logic Methods for RepairTechnicianAssignment ==============

// Validate validates the technician assignment
func (r *RepairTechnicianAssignment) Validate() error {
	if r.RepairBookingID == uuid.Nil {
		return errors.New("repair booking ID is required")
	}
	if r.TechnicianID == uuid.Nil {
		return errors.New("technician ID is required")
	}
	if r.RepairShopID == uuid.Nil {
		return errors.New("repair shop ID is required")
	}
	if r.AssignmentType == "" {
		return errors.New("assignment type is required")
	}

	// Validate assignment type
	validTypes := map[string]bool{
		"primary":    true,
		"secondary":  true,
		"specialist": true,
		"supervisor": true,
	}
	if !validTypes[r.AssignmentType] {
		return errors.New("invalid assignment type")
	}

	return nil
}

// Start marks the assignment as started
func (r *RepairTechnicianAssignment) Start() {
	now := time.Now()
	r.StartedAt = &now
	r.Status = "in_progress"
}

// Complete marks the assignment as completed
func (r *RepairTechnicianAssignment) Complete(hoursWorked decimal.Decimal) {
	now := time.Now()
	r.CompletedAt = &now
	r.HoursWorked = hoursWorked
	r.Status = "completed"
	r.TotalLaborCost = r.LaborRate.Mul(hoursWorked)
}

// CalculateEfficiency calculates efficiency based on expected vs actual time
func (r *RepairTechnicianAssignment) CalculateEfficiency(expectedHours decimal.Decimal) float64 {
	if r.HoursWorked.IsZero() || expectedHours.IsZero() {
		return 0
	}
	efficiency := expectedHours.Div(r.HoursWorked).Mul(decimal.NewFromInt(100))
	efficiencyFloat, _ := efficiency.Float64()
	return efficiencyFloat
}

// Reassign reassigns the task to another technician
func (r *RepairTechnicianAssignment) Reassign(newTechnicianID uuid.UUID, reason string) {
	r.Status = "reassigned"
	r.ReassignReason = reason
	// Note: A new assignment record should be created for the new technician
}

// GetPerformanceScore calculates overall performance score
func (r *RepairTechnicianAssignment) GetPerformanceScore() float64 {
	scores := []float64{}
	count := 0

	if r.QualityScore > 0 {
		scores = append(scores, r.QualityScore)
		count++
	}
	if r.EfficiencyScore > 0 {
		scores = append(scores, r.EfficiencyScore)
		count++
	}
	if r.CustomerRating > 0 {
		scores = append(scores, r.CustomerRating)
		count++
	}

	if count == 0 {
		return 0
	}

	total := 0.0
	for _, score := range scores {
		total += score
	}

	return total / float64(count)
}

// ============== Data Integrity Validation Methods ==============

// ValidateDataIntegrity performs comprehensive data integrity checks
type DataIntegrityValidator struct {
	RepairBooking *RepairBooking
	RepairShop    *RepairShop
	Technician    *Technician
}

// NewDataIntegrityValidator creates a new validator instance
func NewDataIntegrityValidator(booking *RepairBooking, shop *RepairShop, tech *Technician) *DataIntegrityValidator {
	return &DataIntegrityValidator{
		RepairBooking: booking,
		RepairShop:    shop,
		Technician:    tech,
	}
}

// ValidatePartUsage validates that a part can be used in a repair
func (v *DataIntegrityValidator) ValidatePartUsage(partUsed *RepairPartsUsed) error {
	// Check if repair booking exists and is active
	if v.RepairBooking == nil {
		return errors.New("repair booking not found")
	}
	if !v.RepairBooking.IsActive() {
		return errors.New("repair booking is not active")
	}

	// Check if part belongs to the repair booking
	if partUsed.RepairBookingID != v.RepairBooking.ID {
		return errors.New("part does not belong to this repair booking")
	}

	// Check if technician is assigned to this repair
	validTechnician := false
	for _, assignment := range v.RepairBooking.TechnicianAssignments {
		if assignment.TechnicianID == partUsed.TechnicianID {
			validTechnician = true
			break
		}
	}
	if !validTechnician {
		return errors.New("technician is not assigned to this repair")
	}

	// Verify part is in shop inventory
	if v.RepairShop != nil {
		partFound := false
		for _, part := range v.RepairShop.PartsInventory {
			if part.PartNumber == partUsed.PartNumber {
				partFound = true
				// Check if sufficient quantity exists
				if part.QuantityInStock < partUsed.QuantityUsed {
					return errors.New("insufficient part quantity in inventory")
				}
				break
			}
		}
		if !partFound {
			return errors.New("part not found in shop inventory")
		}
	}

	return nil
}

// ValidateAccessoryReplacement validates accessory replacement integrity
func (v *DataIntegrityValidator) ValidateAccessoryReplacement(accessory *RepairAccessoriesReplaced) error {
	// Check if repair booking exists and is active
	if v.RepairBooking == nil {
		return errors.New("repair booking not found")
	}
	if !v.RepairBooking.IsActive() {
		return errors.New("repair booking is not active")
	}

	// Check if accessory replacement belongs to the repair booking
	if accessory.RepairBookingID != v.RepairBooking.ID {
		return errors.New("accessory replacement does not belong to this repair booking")
	}

	// Check if technician is assigned to this repair
	validTechnician := false
	for _, assignment := range v.RepairBooking.TechnicianAssignments {
		if assignment.TechnicianID == accessory.TechnicianID {
			validTechnician = true
			break
		}
	}
	if !validTechnician {
		return errors.New("technician is not assigned to this repair")
	}

	// Validate action type
	switch accessory.ActionType {
	case "replaced":
		if accessory.NewAccessoryID == nil {
			return errors.New("new accessory ID required for replacement")
		}
		if accessory.OldAccessoryID == nil {
			return errors.New("old accessory ID required for replacement")
		}
	case "added":
		if accessory.NewAccessoryID == nil {
			return errors.New("new accessory ID required when adding")
		}
	case "removed":
		if accessory.OldAccessoryID == nil {
			return errors.New("old accessory ID required when removing")
		}
	}

	return nil
}

// ValidateTechnicianAssignment validates technician assignment integrity
func (v *DataIntegrityValidator) ValidateTechnicianAssignment(assignment *RepairTechnicianAssignment) error {
	// Check if repair booking exists
	if v.RepairBooking == nil {
		return errors.New("repair booking not found")
	}

	// Check if repair shop exists
	if v.RepairShop == nil {
		return errors.New("repair shop not found")
	}

	// Check if assignment belongs to correct repair booking and shop
	if assignment.RepairBookingID != v.RepairBooking.ID {
		return errors.New("assignment does not belong to this repair booking")
	}
	if assignment.RepairShopID != v.RepairShop.ID {
		return errors.New("assignment does not belong to this repair shop")
	}

	// Verify technician belongs to the repair shop
	if v.Technician != nil {
		if v.Technician.RepairShopID == nil || *v.Technician.RepairShopID != v.RepairShop.ID {
			return errors.New("technician does not belong to this repair shop")
		}
	}

	// Check for duplicate primary assignment
	if assignment.AssignmentType == "primary" {
		for _, existing := range v.RepairBooking.TechnicianAssignments {
			if existing.ID != assignment.ID &&
				existing.AssignmentType == "primary" &&
				existing.Status != "reassigned" {
				return errors.New("primary technician already assigned")
			}
		}
	}

	// Validate status transitions
	validStatuses := map[string]bool{
		"assigned":    true,
		"in_progress": true,
		"completed":   true,
		"reassigned":  true,
	}
	if !validStatuses[assignment.Status] {
		return errors.New("invalid assignment status")
	}

	return nil
}

// ValidateInventoryConsistency validates inventory consistency after parts usage
func (v *DataIntegrityValidator) ValidateInventoryConsistency(partsUsed []RepairPartsUsed) error {
	if v.RepairShop == nil {
		return errors.New("repair shop not found")
	}

	// Create a map of parts consumption
	consumption := make(map[string]int)
	for _, part := range partsUsed {
		consumption[part.PartNumber] += part.QuantityUsed
	}

	// Validate against inventory
	for partNumber, quantityUsed := range consumption {
		found := false
		for _, inventoryPart := range v.RepairShop.PartsInventory {
			if inventoryPart.PartNumber == partNumber {
				found = true
				available := inventoryPart.QuantityInStock - inventoryPart.ReservedQuantity
				if available < quantityUsed {
					return fmt.Errorf("insufficient inventory for part %s: available %d, required %d",
						partNumber, available, quantityUsed)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("part %s not found in inventory", partNumber)
		}
	}

	return nil
}

// ValidateRepairCompletion validates if repair can be marked as complete
func (v *DataIntegrityValidator) ValidateRepairCompletion() error {
	if v.RepairBooking == nil {
		return errors.New("repair booking not found")
	}

	// Check if quality check passed
	if !v.RepairBooking.QualityCheckPassed {
		return errors.New("quality check not passed")
	}

	// Check if all technician assignments are completed
	for _, assignment := range v.RepairBooking.TechnicianAssignments {
		if assignment.Status == "assigned" || assignment.Status == "in_progress" {
			return fmt.Errorf("technician assignment %s is not completed", assignment.ID)
		}
	}

	// Check if all parts used are marked as successful
	for _, part := range v.RepairBooking.PartsUsed {
		if !part.IsSuccessful {
			return fmt.Errorf("part installation %s failed: %s", part.ID, part.FailureReason)
		}
	}

	// Check if all accessory replacements are approved
	for _, accessory := range v.RepairBooking.AccessoriesReplaced {
		if !accessory.CustomerApproved && accessory.ActionType != "removed" {
			return fmt.Errorf("accessory replacement %s not approved by customer", accessory.ID)
		}
	}

	// Validate costs are calculated
	totalCost := v.RepairBooking.GetTotalCost()
	if totalCost.IsZero() && len(v.RepairBooking.PartsUsed) > 0 {
		return errors.New("repair costs not calculated")
	}

	return nil
}

// ValidateCostIntegrity validates cost calculations across all repair components
func (v *DataIntegrityValidator) ValidateCostIntegrity() error {
	if v.RepairBooking == nil {
		return errors.New("repair booking not found")
	}

	// Calculate expected total from components
	partsCost := v.RepairBooking.CalculatePartsCost()
	accessoriesCost := v.RepairBooking.CalculateAccessoriesCost()
	laborCost := v.RepairBooking.CalculateLaborCost()
	diagnosticFee := decimal.NewFromFloat(v.RepairBooking.DiagnosticFee)
	discountAmount := decimal.NewFromFloat(v.RepairBooking.DiscountAmount)

	expectedTotal := partsCost.Add(accessoriesCost).Add(laborCost).Add(diagnosticFee).Sub(discountAmount)
	actualTotal := decimal.NewFromFloat(v.RepairBooking.ActualCost)

	// Allow for small rounding differences (0.01)
	difference := expectedTotal.Sub(actualTotal).Abs()
	tolerance := decimal.NewFromFloat(0.01)

	if difference.GreaterThan(tolerance) {
		return fmt.Errorf("cost integrity check failed: expected %s, actual %s",
			expectedTotal.String(), actualTotal.String())
	}

	return nil
}

// ValidateTechnicianWorkload validates if technician can take on more work
func (v *DataIntegrityValidator) ValidateTechnicianWorkload(technicianID uuid.UUID) error {
	if v.Technician == nil || v.Technician.ID != technicianID {
		return errors.New("technician not found")
	}

	// Count active assignments
	activeCount := 0
	for _, assignment := range v.Technician.RepairAssignments {
		if assignment.Status == "assigned" || assignment.Status == "in_progress" {
			activeCount++
		}
	}

	// Check against maximum capacity (e.g., 5 concurrent repairs)
	const maxConcurrentRepairs = 5
	if activeCount >= maxConcurrentRepairs {
		return fmt.Errorf("technician at maximum capacity: %d active repairs", activeCount)
	}

	return nil
}

// ValidatePartCompatibility validates if a part is compatible with the device
func (v *DataIntegrityValidator) ValidatePartCompatibility(partUsed *RepairPartsUsed, deviceModel string) error {
	// This would typically check against a compatibility matrix
	// For now, we'll do basic validation

	if partUsed.PartNumber == "" {
		return errors.New("part number is required")
	}

	if deviceModel == "" {
		return errors.New("device model is required for compatibility check")
	}

	// Check if part quality meets minimum standards
	acceptableQualities := map[string]bool{
		"oem":         true,
		"genuine":     true,
		"aftermarket": true,
	}

	if partUsed.Quality != "" && !acceptableQualities[partUsed.Quality] {
		return fmt.Errorf("part quality %s not acceptable", partUsed.Quality)
	}

	return nil
}

// ValidateWarrantyEligibility validates if repair is eligible for warranty
func (v *DataIntegrityValidator) ValidateWarrantyEligibility() error {
	if v.RepairBooking == nil {
		return errors.New("repair booking not found")
	}

	// Check if any non-genuine parts were used
	for _, part := range v.RepairBooking.PartsUsed {
		if part.Quality != "oem" && part.Quality != "genuine" {
			if v.RepairBooking.IsUnderWarranty {
				return errors.New("non-genuine parts void warranty")
			}
		}
	}

	// Check if repair shop is authorized for warranty repairs
	if v.RepairBooking.IsUnderWarranty && v.RepairShop != nil {
		if v.RepairShop.ShopCertification.CertificationLevel != "platinum" &&
			v.RepairShop.ShopCertification.CertificationLevel != "gold" {
			return errors.New("shop not authorized for warranty repairs")
		}
	}

	return nil
}
