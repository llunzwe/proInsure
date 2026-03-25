package models

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"smartsure/internal/domain/models/technician"
)

// Technician represents a repair technician with comprehensive management
type Technician struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// Core Relationships
	RepairShopID *uuid.UUID `gorm:"type:uuid;index" json:"repair_shop_id,omitempty"` // Primary repair shop association
	UserID       *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`        // Link to user account if applicable

	// Embedded structs for organization
	technician.TechnicianIdentification `gorm:"embedded;embeddedPrefix:id_" json:"identification"`
	technician.TechnicianPersonal       `gorm:"embedded;embeddedPrefix:pers_" json:"personal"`
	technician.TechnicianContact        `gorm:"embedded;embeddedPrefix:cont_" json:"contact"`
	technician.TechnicianProfessional   `gorm:"embedded;embeddedPrefix:prof_" json:"professional"`
	technician.TechnicianAvailability   `gorm:"embedded;embeddedPrefix:avail_" json:"availability"`
	technician.TechnicianPerformance    `gorm:"embedded;embeddedPrefix:perf_" json:"performance"`
	technician.TechnicianCompensation   `gorm:"embedded;embeddedPrefix:comp_" json:"compensation"`
	technician.TechnicianCompliance     `gorm:"embedded;embeddedPrefix:compl_" json:"compliance"`
	technician.TechnicianLocation       `gorm:"embedded;embeddedPrefix:loc_" json:"location"`
	technician.TechnicianStatus         `gorm:"embedded;embeddedPrefix:status_" json:"status"`
	technician.TechnicianMetrics        `gorm:"embedded;embeddedPrefix:metrics_" json:"metrics"`

	// Timestamps
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships - Core Models
	RepairShop         *RepairShop                  `gorm:"foreignKey:RepairShopID" json:"repair_shop,omitempty"`
	User               *User                        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RepairAssignments  []RepairTechnicianAssignment `gorm:"foreignKey:TechnicianID" json:"repair_assignments,omitempty"`
	PartsUsed          []RepairPartsUsed            `gorm:"foreignKey:TechnicianID" json:"parts_used,omitempty"`
	AccessoriesHandled []RepairAccessoriesReplaced  `gorm:"foreignKey:TechnicianID" json:"accessories_handled,omitempty"`

	// Relationships - Feature Models
	Assignments        []technician.TechnicianAssignment        `gorm:"foreignKey:TechnicianID" json:"assignments,omitempty"`
	Schedules          []technician.TechnicianSchedule          `gorm:"foreignKey:TechnicianID" json:"schedules,omitempty"`
	Skills             []technician.TechnicianSkill             `gorm:"foreignKey:TechnicianID" json:"skills,omitempty"`
	Trainings          []technician.TechnicianTraining          `gorm:"foreignKey:TechnicianID" json:"trainings,omitempty"`
	Certifications     []technician.TechnicianCertification     `gorm:"foreignKey:TechnicianID" json:"certifications,omitempty"`
	PerformanceReviews []technician.TechnicianPerformanceReview `gorm:"foreignKey:TechnicianID" json:"performance_reviews,omitempty"`
	Leaves             []technician.TechnicianLeave             `gorm:"foreignKey:TechnicianID" json:"leaves,omitempty"`
	Timesheets         []technician.TechnicianTimesheet         `gorm:"foreignKey:TechnicianID" json:"timesheets,omitempty"`
	Equipment          []technician.TechnicianEquipment         `gorm:"foreignKey:TechnicianID" json:"equipment,omitempty"`
	Documents          []technician.TechnicianDocument          `gorm:"foreignKey:TechnicianID" json:"documents,omitempty"`
	Incidents          []technician.TechnicianIncident          `gorm:"foreignKey:TechnicianID" json:"incidents,omitempty"`
}

func (Technician) TableName() string {
	return "technicians"
}

// BeforeCreate hook
func (t *Technician) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.TechnicianIdentification.EmployeeID == "" {
		t.TechnicianIdentification.EmployeeID = "EMP-" + uuid.New().String()[:8]
	}
	if t.TechnicianProfessional.JoinDate.IsZero() {
		t.TechnicianProfessional.JoinDate = time.Now()
	}
	if t.TechnicianProfessional.Type == "" {
		t.TechnicianProfessional.Type = technician.TechnicianTypeEmployee
	}
	if t.TechnicianProfessional.Level == "" {
		t.TechnicianProfessional.Level = technician.CertificationJunior
	}
	return t.Validate()
}

// ============== Business Logic Methods for Technician ==============

// Validation Methods
func (t *Technician) Validate() error {
	if t.TechnicianPersonal.FirstName == "" || t.TechnicianPersonal.LastName == "" {
		return errors.New("technician name is required")
	}
	if t.TechnicianContact.Email == "" {
		return errors.New("email is required")
	}
	if t.TechnicianContact.Phone == "" {
		return errors.New("phone number is required")
	}
	if t.TechnicianContact.Address == "" || t.TechnicianContact.City == "" {
		return errors.New("address is incomplete")
	}
	if t.TechnicianProfessional.Position == "" {
		return errors.New("position is required")
	}
	if t.TechnicianPersonal.DateOfBirth.IsZero() {
		return errors.New("date of birth is required")
	}
	// Age validation
	age := t.GetAge()
	if age < 18 {
		return errors.New("technician must be at least 18 years old")
	}
	if age > 70 {
		return errors.New("technician age exceeds maximum allowed")
	}
	return nil
}

// Status Methods
func (t *Technician) IsActive() bool {
	return t.TechnicianStatus.IsActive &&
		t.TechnicianStatus.EmploymentStatus == technician.StatusActive &&
		!t.TechnicianStatus.OnLeave
}

func (t *Technician) IsAvailable() bool {
	return t.IsActive() &&
		t.TechnicianAvailability.Status == technician.AvailabilityAvailable &&
		t.TechnicianAvailability.CurrentAssignments < t.TechnicianAvailability.MaxAssignments
}

func (t *Technician) IsOnLeave() bool {
	return t.TechnicianStatus.OnLeave &&
		t.TechnicianStatus.LeaveStartDate != nil &&
		t.TechnicianStatus.LeaveEndDate != nil &&
		time.Now().After(*t.TechnicianStatus.LeaveStartDate) &&
		time.Now().Before(*t.TechnicianStatus.LeaveEndDate)
}

func (t *Technician) IsInProbation() bool {
	if t.TechnicianProfessional.ProbationEndDate == nil {
		return false
	}
	return time.Now().Before(*t.TechnicianProfessional.ProbationEndDate)
}

func (t *Technician) IsCertified() bool {
	return t.TechnicianProfessional.CertificationLevel != technician.CertificationNone &&
		t.TechnicianProfessional.CertificationLevel != technician.CertificationJunior
}

func (t *Technician) IsEligibleForPromotion() bool {
	if !t.IsActive() || t.IsInProbation() {
		return false
	}

	// Check performance
	if t.TechnicianPerformance.Rating < technician.MinimumPerformanceScore {
		return false
	}

	// Check experience
	yearsOfService := t.GetYearsOfService()
	switch t.TechnicianProfessional.Level {
	case technician.CertificationJunior:
		return yearsOfService >= 1 && t.TechnicianPerformance.CompletedRepairs >= 100
	case technician.CertificationIntermediate:
		return yearsOfService >= 2 && t.TechnicianPerformance.CompletedRepairs >= 500
	case technician.CertificationSenior:
		return yearsOfService >= 5 && t.TechnicianPerformance.CompletedRepairs >= 1000
	default:
		return false
	}
}

// Personal Methods
func (t *Technician) GetFullName() string {
	if t.TechnicianPersonal.PreferredName != "" {
		return t.TechnicianPersonal.PreferredName + " " + t.TechnicianPersonal.LastName
	}
	middle := ""
	if t.TechnicianPersonal.MiddleName != "" {
		middle = " " + t.TechnicianPersonal.MiddleName
	}
	return t.TechnicianPersonal.FirstName + middle + " " + t.TechnicianPersonal.LastName
}

func (t *Technician) GetAge() int {
	if t.TechnicianPersonal.DateOfBirth.IsZero() {
		return 0
	}
	return int(time.Since(t.TechnicianPersonal.DateOfBirth).Hours() / 24 / 365)
}

func (t *Technician) GetYearsOfService() float64 {
	if t.TechnicianProfessional.JoinDate.IsZero() {
		return 0
	}
	return time.Since(t.TechnicianProfessional.JoinDate).Hours() / 24 / 365
}

// Availability Methods
func (t *Technician) CanAcceptAssignment() bool {
	if !t.IsAvailable() {
		return false
	}
	return t.TechnicianAvailability.CurrentAssignments < t.TechnicianAvailability.MaxAssignments
}

func (t *Technician) GetAvailableSlots() int {
	if !t.IsAvailable() {
		return 0
	}
	return t.TechnicianAvailability.MaxAssignments - t.TechnicianAvailability.CurrentAssignments
}

func (t *Technician) IsWorkingToday() bool {
	if !t.IsActive() {
		return false
	}

	// Check if today is in working days
	today := time.Now().Weekday().String()
	if t.TechnicianAvailability.WorkingDays != "" {
		return strings.Contains(strings.ToLower(t.TechnicianAvailability.WorkingDays), strings.ToLower(today))
	}

	// Check schedule
	for _, schedule := range t.Schedules {
		if schedule.ScheduleDate.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			return schedule.Status == "scheduled" || schedule.Status == "checked_in"
		}
	}

	return false
}

func (t *Technician) IsOnShift() bool {
	if !t.IsWorkingToday() {
		return false
	}

	if t.TechnicianAvailability.ShiftStartTime != nil && t.TechnicianAvailability.ShiftEndTime != nil {
		now := time.Now()
		return now.After(*t.TechnicianAvailability.ShiftStartTime) && now.Before(*t.TechnicianAvailability.ShiftEndTime)
	}

	return false
}

func (t *Technician) IsOnBreak() bool {
	if !t.IsOnShift() {
		return false
	}

	if t.TechnicianAvailability.BreakStartTime != nil && t.TechnicianAvailability.BreakEndTime != nil {
		now := time.Now()
		return now.After(*t.TechnicianAvailability.BreakStartTime) && now.Before(*t.TechnicianAvailability.BreakEndTime)
	}

	return false
}

// Skill Methods
func (t *Technician) HasSkill(skillName string) bool {
	for _, skill := range t.Skills {
		if strings.EqualFold(skill.SkillName, skillName) {
			return skill.IsVerified
		}
	}
	return false
}

func (t *Technician) GetSkillLevel(skillName string) string {
	for _, skill := range t.Skills {
		if strings.EqualFold(skill.SkillName, skillName) {
			return skill.ProficiencyLevel
		}
	}
	return technician.SkillBeginner
}

func (t *Technician) HasSpecialization(specialization string) bool {
	if t.TechnicianProfessional.Specializations == "" {
		return false
	}
	return strings.Contains(strings.ToLower(t.TechnicianProfessional.Specializations), strings.ToLower(specialization))
}

func (t *Technician) CanPerformRepair(repairType string) bool {
	// Check if has the skill
	switch repairType {
	case "screen_repair":
		return t.HasSkill("screen_replacement") || t.HasSpecialization(technician.SpecializationScreen)
	case "battery_replacement":
		return t.HasSkill("battery_replacement") || t.HasSpecialization(technician.SpecializationBattery)
	case "water_damage":
		return t.HasSkill("water_damage_repair") || t.HasSpecialization(technician.SpecializationWaterDamage)
	case "motherboard":
		return t.HasSkill("board_level_repair") || t.HasSpecialization(technician.SpecializationMotherboard)
	default:
		return t.TechnicianProfessional.Level != technician.CertificationJunior
	}
}

// Performance Methods
func (t *Technician) GetPerformanceScore() float64 {
	weights := map[string]float64{
		"rating":          0.20,
		"customer_rating": 0.25,
		"success_rate":    0.20,
		"first_time_fix":  0.15,
		"quality_score":   0.10,
		"productivity":    0.10,
	}

	score := t.TechnicianPerformance.Rating*weights["rating"] +
		t.TechnicianPerformance.CustomerRating*weights["customer_rating"] +
		t.TechnicianPerformance.SuccessRate*weights["success_rate"] +
		t.TechnicianPerformance.FirstTimeFixRate*weights["first_time_fix"] +
		t.TechnicianPerformance.QualityScore*weights["quality_score"] +
		t.TechnicianPerformance.ProductivityScore*weights["productivity"]

	return math.Min(score, 100)
}

func (t *Technician) GetPerformanceRating() string {
	score := t.GetPerformanceScore()
	switch {
	case score >= 90:
		return technician.PerformanceOutstanding
	case score >= 80:
		return technician.PerformanceExcellent
	case score >= 70:
		return technician.PerformanceGood
	case score >= 60:
		return technician.PerformanceSatisfactory
	case score >= 50:
		return technician.PerformanceNeedsImprove
	default:
		return technician.PerformancePoor
	}
}

func (t *Technician) IsHighPerformer() bool {
	return t.GetPerformanceScore() >= 80 &&
		t.TechnicianPerformance.CustomerRating >= technician.MinimumCustomerRating &&
		t.TechnicianPerformance.FirstTimeFixRate >= technician.MinimumFirstTimeFixRate &&
		t.TechnicianPerformance.ErrorRate <= technician.MaximumErrorRate
}

func (t *Technician) NeedsPerformanceReview() bool {
	if t.TechnicianPerformance.LastPerformanceReview == nil {
		return true
	}
	daysSinceReview := time.Since(*t.TechnicianPerformance.LastPerformanceReview).Hours() / 24
	return daysSinceReview >= float64(technician.DefaultPerformanceReviewDays)
}

// Compensation Methods
func (t *Technician) GetHourlyRate() decimal.Decimal {
	if t.TechnicianCompensation.CompensationType == technician.CompensationHourly {
		return t.TechnicianCompensation.HourlyRate
	}
	// Calculate from salary if salaried
	if t.TechnicianCompensation.CompensationType == technician.CompensationSalary {
		// Assuming 2080 working hours per year (40 hours/week * 52 weeks)
		return t.TechnicianCompensation.BaseSalary.Div(decimal.NewFromInt(2080))
	}
	return decimal.Zero
}

func (t *Technician) CalculateOvertimeRate() decimal.Decimal {
	if t.TechnicianCompensation.OvertimeRate.GreaterThan(decimal.Zero) {
		return t.TechnicianCompensation.OvertimeRate
	}
	return t.GetHourlyRate().Mul(decimal.NewFromFloat(technician.DefaultOvertimeMultiplier))
}

func (t *Technician) GetVacationDaysRemaining() int {
	return t.TechnicianCompensation.VacationDays - t.TechnicianCompensation.UsedVacationDays
}

func (t *Technician) GetSickDaysRemaining() int {
	return t.TechnicianCompensation.SickDays - t.TechnicianCompensation.UsedSickDays
}

func (t *Technician) IsEligibleForBonus() bool {
	return t.TechnicianCompensation.BonusEligible && t.IsHighPerformer()
}

// Certification Methods
func (t *Technician) HasValidCertification(certName string) bool {
	for _, cert := range t.Certifications {
		if strings.EqualFold(cert.CertificationName, certName) {
			if cert.Status == "active" && cert.IsVerified {
				if cert.ExpiryDate == nil {
					return true
				}
				return time.Now().Before(*cert.ExpiryDate)
			}
		}
	}
	return false
}

func (t *Technician) GetActiveCertifications() []technician.TechnicianCertification {
	var active []technician.TechnicianCertification
	for _, cert := range t.Certifications {
		if cert.Status == "active" && cert.IsVerified {
			if cert.ExpiryDate == nil || time.Now().Before(*cert.ExpiryDate) {
				active = append(active, cert)
			}
		}
	}
	return active
}

func (t *Technician) NeedsCertificationRenewal() bool {
	for _, cert := range t.Certifications {
		if cert.RenewalRequired && cert.RenewalDate != nil {
			if time.Now().AddDate(0, 1, 0).After(*cert.RenewalDate) {
				return true
			}
		}
	}
	return false
}

// Training Methods
func (t *Technician) HasCompletedTraining(trainingName string) bool {
	for _, training := range t.Trainings {
		if strings.EqualFold(training.TrainingName, trainingName) {
			return training.Status == technician.TrainingCompleted && training.TestScore >= training.PassingScore
		}
	}
	return false
}

func (t *Technician) GetPendingTrainings() []technician.TechnicianTraining {
	var pending []technician.TechnicianTraining
	for _, training := range t.Trainings {
		if training.IsMandatory && (training.Status == technician.TrainingScheduled || training.Status == technician.TrainingPending) {
			pending = append(pending, training)
		}
	}
	return pending
}

func (t *Technician) NeedsTrainingRefresh() bool {
	for _, training := range t.Trainings {
		if training.IsRecurring && training.ExpiryDate != nil {
			if time.Now().After(*training.ExpiryDate) {
				return true
			}
		}
	}
	return false
}

// Compliance Methods
func (t *Technician) IsCompliant() bool {
	return t.TechnicianCompliance.ComplianceStatus == "compliant" &&
		t.TechnicianCompliance.BackgroundCheckStatus == technician.BackgroundClear &&
		t.TechnicianCompliance.CriminalRecordClear &&
		!t.TechnicianCompliance.UnderInvestigation
}

func (t *Technician) RequiresBackgroundCheck() bool {
	if t.TechnicianCompliance.BackgroundCheckDate == nil {
		return true
	}
	if t.TechnicianCompliance.BackgroundCheckExpiry != nil {
		return time.Now().After(*t.TechnicianCompliance.BackgroundCheckExpiry)
	}
	// Default to check every 2 years
	return time.Since(*t.TechnicianCompliance.BackgroundCheckDate).Hours()/24/365 >= float64(technician.DefaultBackgroundCheckYears)
}

func (t *Technician) HasSignedRequiredDocuments() bool {
	return t.TechnicianCompliance.NDASigned && t.TechnicianCompliance.ContractSigned
}

// Location Methods
func (t *Technician) GetDistanceFromHome(lat, lng float64) float64 {
	// Haversine formula
	const earthRadius = 6371 // km

	lat1 := t.TechnicianLocation.HomeLatitude * math.Pi / 180
	lat2 := lat * math.Pi / 180
	deltaLat := (lat - t.TechnicianLocation.HomeLatitude) * math.Pi / 180
	deltaLon := (lng - t.TechnicianLocation.HomeLongitude) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (t *Technician) IsWithinServiceRadius(lat, lng float64) bool {
	if !t.TechnicianAvailability.TravelWilling {
		return false
	}
	distance := t.GetDistanceFromHome(lat, lng)
	return distance <= float64(t.TechnicianAvailability.MaxTravelDistance)
}

func (t *Technician) UpdateCurrentLocation(lat, lng float64) {
	t.TechnicianLocation.CurrentLatitude = lat
	t.TechnicianLocation.CurrentLongitude = lng
	now := time.Now()
	t.TechnicianLocation.LocationUpdatedAt = &now
}

// Assignment Methods
func (t *Technician) GetActiveAssignments() []technician.TechnicianAssignment {
	var active []technician.TechnicianAssignment
	for _, assignment := range t.Assignments {
		if assignment.Status == technician.AssignmentAccepted ||
			assignment.Status == technician.AssignmentInProgress {
			active = append(active, assignment)
		}
	}
	return active
}

func (t *Technician) GetTodaysAssignments() []technician.TechnicianAssignment {
	var today []technician.TechnicianAssignment
	todayStr := time.Now().Format("2006-01-02")
	for _, assignment := range t.Assignments {
		if assignment.ScheduledStart.Format("2006-01-02") == todayStr {
			today = append(today, assignment)
		}
	}
	return today
}

func (t *Technician) CanAcceptUrgentAssignment() bool {
	if !t.IsAvailable() {
		return false
	}
	// Check if has capacity for urgent work
	activeAssignments := t.GetActiveAssignments()
	for _, assignment := range activeAssignments {
		if assignment.Priority == technician.PriorityEmergency || assignment.Priority == technician.PriorityCritical {
			return false // Already on urgent assignment
		}
	}
	return true
}

// Equipment Methods
func (t *Technician) GetAssignedEquipment() []technician.TechnicianEquipment {
	var assigned []technician.TechnicianEquipment
	for _, equipment := range t.Equipment {
		if equipment.Status == technician.EquipmentAssigned {
			assigned = append(assigned, equipment)
		}
	}
	return assigned
}

func (t *Technician) HasEquipmentOverdue() bool {
	for _, equipment := range t.Equipment {
		if equipment.MaintenanceDue != nil && time.Now().After(*equipment.MaintenanceDue) {
			return true
		}
	}
	return false
}

// Incident Methods
func (t *Technician) GetSafetyScore() float64 {
	if len(t.Incidents) == 0 {
		return 100.0
	}

	// Calculate based on incident severity and frequency
	score := 100.0
	for _, incident := range t.Incidents {
		if incident.IncidentDate.After(time.Now().AddDate(-1, 0, 0)) { // Last year
			switch incident.Severity {
			case "critical":
				score -= 20
			case "major":
				score -= 10
			case "minor":
				score -= 5
			}
		}
	}

	return math.Max(score, 0)
}

func (t *Technician) HasRecentIncident() bool {
	for _, incident := range t.Incidents {
		if incident.IncidentDate.After(time.Now().AddDate(0, -3, 0)) { // Last 3 months
			return true
		}
	}
	return false
}

// Summary Methods
func (t *Technician) GetSummary() string {
	return fmt.Sprintf("%s (%s) - %s %s. Status: %s, Performance: %.1f/100, Completed: %d repairs",
		t.GetFullName(),
		t.TechnicianIdentification.EmployeeID,
		t.TechnicianProfessional.Level,
		t.TechnicianProfessional.Position,
		t.TechnicianStatus.EmploymentStatus,
		t.GetPerformanceScore(),
		t.TechnicianPerformance.CompletedRepairs,
	)
}

func (t *Technician) GetDisplayName() string {
	return fmt.Sprintf("%s (%s)", t.GetFullName(), t.TechnicianIdentification.EmployeeID)
}

func (t *Technician) GetBadgeInfo() string {
	certLevel := t.TechnicianProfessional.CertificationLevel
	if certLevel == "" {
		certLevel = technician.CertificationNone
	}
	return fmt.Sprintf("%s - %s", t.TechnicianIdentification.BadgeNumber, certLevel)
}
