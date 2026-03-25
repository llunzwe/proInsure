package technician

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TechnicianIdentification contains basic identification info
type TechnicianIdentification struct {
	EmployeeID       string     `gorm:"uniqueIndex;not null" json:"employee_id"`
	BadgeNumber      string     `gorm:"uniqueIndex" json:"badge_number"`
	SSN              string     `json:"-"` // Encrypted, hidden in JSON
	TaxID            string     `json:"-"` // Encrypted, hidden in JSON
	LicenseNumber    string     `json:"license_number,omitempty"`
	LicenseType      string     `json:"license_type,omitempty"`
	LicenseExpiry    time.Time  `json:"license_expiry,omitempty"`
	LicenseState     string     `json:"license_state,omitempty"`
	PassportNumber   string     `json:"-"` // Hidden in JSON
	VisaStatus       string     `json:"visa_status,omitempty"`
	VisaExpiry       *time.Time `json:"visa_expiry,omitempty"`
	WorkPermitNumber string     `json:"work_permit_number,omitempty"`
	WorkPermitExpiry *time.Time `json:"work_permit_expiry,omitempty"`
}

// TechnicianPersonal contains personal information
type TechnicianPersonal struct {
	FirstName                string    `gorm:"not null" json:"first_name"`
	LastName                 string    `gorm:"not null" json:"last_name"`
	MiddleName               string    `json:"middle_name,omitempty"`
	PreferredName            string    `json:"preferred_name,omitempty"`
	DateOfBirth              time.Time `json:"date_of_birth"`
	Gender                   string    `json:"gender,omitempty"`
	MaritalStatus            string    `json:"marital_status,omitempty"`
	Nationality              string    `json:"nationality,omitempty"`
	Languages                string    `json:"languages"` // JSON array
	ProfilePhotoURL          string    `json:"profile_photo_url,omitempty"`
	EmergencyContactName     string    `json:"emergency_contact_name"`
	EmergencyContactPhone    string    `json:"emergency_contact_phone"`
	EmergencyContactRelation string    `json:"emergency_contact_relation"`
	BloodType                string    `json:"blood_type,omitempty"`
	MedicalConditions        string    `json:"medical_conditions,omitempty"` // JSON array, encrypted
	Allergies                string    `json:"allergies,omitempty"`          // JSON array
}

// TechnicianContact contains contact information
type TechnicianContact struct {
	Email             string     `gorm:"uniqueIndex;not null" json:"email"`
	Phone             string     `gorm:"not null" json:"phone"`
	AlternatePhone    string     `json:"alternate_phone,omitempty"`
	Address           string     `gorm:"not null" json:"address"`
	AddressLine2      string     `json:"address_line2,omitempty"`
	City              string     `gorm:"not null" json:"city"`
	State             string     `gorm:"not null" json:"state"`
	Country           string     `gorm:"not null" json:"country"`
	PostalCode        string     `json:"postal_code"`
	PermanentAddress  string     `json:"permanent_address,omitempty"`
	MailingAddress    string     `json:"mailing_address,omitempty"`
	PreferredContact  string     `gorm:"default:'email'" json:"preferred_contact"`
	NotificationPrefs string     `json:"notification_prefs"` // JSON object
	DoNotDisturb      bool       `gorm:"default:false" json:"do_not_disturb"`
	DoNotDisturbStart *time.Time `json:"do_not_disturb_start,omitempty"`
	DoNotDisturbEnd   *time.Time `json:"do_not_disturb_end,omitempty"`
}

// TechnicianProfessional contains professional qualifications
type TechnicianProfessional struct {
	Type               string     `gorm:"not null;default:'employee'" json:"type"`
	Department         string     `json:"department,omitempty"`
	Position           string     `gorm:"not null" json:"position"`
	Level              string     `gorm:"not null;default:'junior'" json:"level"`
	Specializations    string     `json:"specializations"` // JSON array
	Skills             string     `json:"skills"`          // JSON array of skill objects
	CertificationLevel string     `gorm:"default:'junior'" json:"certification_level"`
	Certifications     string     `json:"certifications"` // JSON array of certification objects
	YearsExperience    int        `json:"years_experience"`
	PreviousEmployers  string     `json:"previous_employers,omitempty"` // JSON array
	Education          string     `json:"education,omitempty"`          // JSON array
	TrainingCompleted  string     `json:"training_completed"`           // JSON array
	JoinDate           time.Time  `gorm:"not null" json:"join_date"`
	ConfirmationDate   *time.Time `json:"confirmation_date,omitempty"`
	ProbationEndDate   *time.Time `json:"probation_end_date,omitempty"`
	LastPromotionDate  *time.Time `json:"last_promotion_date,omitempty"`
	NextReviewDate     *time.Time `json:"next_review_date,omitempty"`
}

// TechnicianAvailability contains work availability info
type TechnicianAvailability struct {
	Status             string     `gorm:"default:'available'" json:"status"`
	CurrentAssignments int        `gorm:"default:0" json:"current_assignments"`
	MaxAssignments     int        `gorm:"default:5" json:"max_assignments"`
	ShiftPreference    string     `json:"shift_preference,omitempty"`
	CurrentShift       string     `json:"current_shift,omitempty"`
	ShiftStartTime     *time.Time `json:"shift_start_time,omitempty"`
	ShiftEndTime       *time.Time `json:"shift_end_time,omitempty"`
	BreakStartTime     *time.Time `json:"break_start_time,omitempty"`
	BreakEndTime       *time.Time `json:"break_end_time,omitempty"`
	WorkingDays        string     `json:"working_days"`  // JSON array of days
	WorkingHours       string     `json:"working_hours"` // JSON object
	OvertimeAvailable  bool       `gorm:"default:false" json:"overtime_available"`
	WeekendAvailable   bool       `gorm:"default:false" json:"weekend_available"`
	HolidayAvailable   bool       `gorm:"default:false" json:"holiday_available"`
	TravelWilling      bool       `gorm:"default:false" json:"travel_willing"`
	MaxTravelDistance  int        `json:"max_travel_distance"` // km
	RemoteCapable      bool       `gorm:"default:false" json:"remote_capable"`
	OnCallAvailable    bool       `gorm:"default:false" json:"on_call_available"`
}

// TechnicianPerformance contains performance metrics
type TechnicianPerformance struct {
	Rating                float64    `gorm:"default:0" json:"rating"`
	CustomerRating        float64    `gorm:"default:0" json:"customer_rating"`
	PeerRating            float64    `gorm:"default:0" json:"peer_rating"`
	SupervisorRating      float64    `gorm:"default:0" json:"supervisor_rating"`
	CompletedRepairs      int        `gorm:"default:0" json:"completed_repairs"`
	SuccessRate           float64    `gorm:"default:0" json:"success_rate"`
	FirstTimeFixRate      float64    `gorm:"default:0" json:"first_time_fix_rate"`
	ReworkRate            float64    `gorm:"default:0" json:"rework_rate"`
	AverageRepairTime     int        `json:"average_repair_time"` // minutes
	FastestRepairTime     int        `json:"fastest_repair_time"` // minutes
	QualityScore          float64    `gorm:"default:0" json:"quality_score"`
	ProductivityScore     float64    `gorm:"default:0" json:"productivity_score"`
	EfficiencyRate        float64    `gorm:"default:0" json:"efficiency_rate"`
	ErrorRate             float64    `gorm:"default:0" json:"error_rate"`
	ComplaintRate         float64    `gorm:"default:0" json:"complaint_rate"`
	ComplimentCount       int        `gorm:"default:0" json:"compliment_count"`
	AttendanceRate        float64    `gorm:"default:0" json:"attendance_rate"`
	PunctualityScore      float64    `gorm:"default:0" json:"punctuality_score"`
	TrainingScore         float64    `gorm:"default:0" json:"training_score"`
	SafetyScore           float64    `gorm:"default:0" json:"safety_score"`
	LastPerformanceReview *time.Time `json:"last_performance_review,omitempty"`
}

// TechnicianCompensation contains compensation and benefits info
type TechnicianCompensation struct {
	CompensationType   string          `gorm:"default:'hourly'" json:"compensation_type"`
	BaseSalary         decimal.Decimal `sql:"type:decimal(10,2)" json:"base_salary"`
	HourlyRate         decimal.Decimal `sql:"type:decimal(10,2)" json:"hourly_rate"`
	OvertimeRate       decimal.Decimal `sql:"type:decimal(10,2)" json:"overtime_rate"`
	CommissionRate     float64         `gorm:"default:0" json:"commission_rate"`
	BonusEligible      bool            `gorm:"default:false" json:"bonus_eligible"`
	LastBonusAmount    decimal.Decimal `sql:"type:decimal(10,2)" json:"last_bonus_amount"`
	LastBonusDate      *time.Time      `json:"last_bonus_date,omitempty"`
	TotalEarnings      decimal.Decimal `sql:"type:decimal(10,2)" json:"total_earnings"`
	Benefits           string          `json:"benefits,omitempty"` // JSON array
	InsuranceEnrolled  bool            `gorm:"default:false" json:"insurance_enrolled"`
	RetirementEnrolled bool            `gorm:"default:false" json:"retirement_enrolled"`
	VacationDays       int             `gorm:"default:0" json:"vacation_days"`
	SickDays           int             `gorm:"default:0" json:"sick_days"`
	UsedVacationDays   int             `gorm:"default:0" json:"used_vacation_days"`
	UsedSickDays       int             `gorm:"default:0" json:"used_sick_days"`
	PayrollID          string          `json:"payroll_id,omitempty"`
	BankAccountNumber  string          `json:"-"` // Encrypted, hidden
	BankRoutingNumber  string          `json:"-"` // Encrypted, hidden
	PaymentMethod      string          `gorm:"default:'direct_deposit'" json:"payment_method"`
}

// TechnicianCompliance contains compliance and legal info
type TechnicianCompliance struct {
	BackgroundCheckStatus   string     `gorm:"default:'pending'" json:"background_check_status"`
	BackgroundCheckDate     *time.Time `json:"background_check_date,omitempty"`
	BackgroundCheckExpiry   *time.Time `json:"background_check_expiry,omitempty"`
	DrugTestStatus          string     `json:"drug_test_status,omitempty"`
	DrugTestDate            *time.Time `json:"drug_test_date,omitempty"`
	DrivingRecordStatus     string     `json:"driving_record_status,omitempty"`
	DrivingRecordDate       *time.Time `json:"driving_record_date,omitempty"`
	CriminalRecordClear     bool       `gorm:"default:false" json:"criminal_record_clear"`
	SecurityClearance       string     `json:"security_clearance,omitempty"`
	SecurityClearanceExpiry *time.Time `json:"security_clearance_expiry,omitempty"`
	NDASigned               bool       `gorm:"default:false" json:"nda_signed"`
	NDASignedDate           *time.Time `json:"nda_signed_date,omitempty"`
	ContractSigned          bool       `gorm:"default:false" json:"contract_signed"`
	ContractSignedDate      *time.Time `json:"contract_signed_date,omitempty"`
	ContractEndDate         *time.Time `json:"contract_end_date,omitempty"`
	ComplianceTraining      string     `json:"compliance_training,omitempty"` // JSON array
	ComplianceStatus        string     `gorm:"default:'pending'" json:"compliance_status"`
	Violations              int        `gorm:"default:0" json:"violations"`
	LastViolationDate       *time.Time `json:"last_violation_date,omitempty"`
	UnderInvestigation      bool       `gorm:"default:false" json:"under_investigation"`
}

// TechnicianLocation contains location and assignment info
type TechnicianLocation struct {
	PrimaryShopID       *uuid.UUID `gorm:"type:uuid" json:"primary_shop_id,omitempty"`
	CurrentShopID       *uuid.UUID `gorm:"type:uuid" json:"current_shop_id,omitempty"`
	AssignedRegion      string     `json:"assigned_region,omitempty"`
	AssignedTerritory   string     `json:"assigned_territory,omitempty"`
	ServiceAreas        string     `json:"service_areas,omitempty"` // JSON array
	CurrentLocation     string     `json:"current_location,omitempty"`
	CurrentLatitude     float64    `json:"current_latitude,omitempty"`
	CurrentLongitude    float64    `json:"current_longitude,omitempty"`
	LocationUpdatedAt   *time.Time `json:"location_updated_at,omitempty"`
	HomeLatitude        float64    `json:"home_latitude,omitempty"`
	HomeLongitude       float64    `json:"home_longitude,omitempty"`
	CommuteDistance     float64    `json:"commute_distance,omitempty"` // km
	CommuteTime         int        `json:"commute_time,omitempty"`     // minutes
	TransportMode       string     `json:"transport_mode,omitempty"`
	CompanyVehicle      bool       `gorm:"default:false" json:"company_vehicle"`
	VehicleLicensePlate string     `json:"vehicle_license_plate,omitempty"`
}

// TechnicianStatus contains status and lifecycle info
type TechnicianStatus struct {
	IsActive           bool       `gorm:"default:true;index" json:"is_active"`
	EmploymentStatus   string     `gorm:"default:'active'" json:"employment_status"`
	StatusReason       string     `json:"status_reason,omitempty"`
	StatusChangedAt    *time.Time `json:"status_changed_at,omitempty"`
	StatusChangedBy    *uuid.UUID `gorm:"type:uuid" json:"status_changed_by,omitempty"`
	OnLeave            bool       `gorm:"default:false" json:"on_leave"`
	LeaveType          string     `json:"leave_type,omitempty"`
	LeaveStartDate     *time.Time `json:"leave_start_date,omitempty"`
	LeaveEndDate       *time.Time `json:"leave_end_date,omitempty"`
	LeaveReason        string     `json:"leave_reason,omitempty"`
	LeaveApprovedBy    *uuid.UUID `gorm:"type:uuid" json:"leave_approved_by,omitempty"`
	ReturnDate         *time.Time `json:"return_date,omitempty"`
	TerminationDate    *time.Time `json:"termination_date,omitempty"`
	TerminationReason  string     `json:"termination_reason,omitempty"`
	EligibleForRehire  bool       `gorm:"default:true" json:"eligible_for_rehire"`
	LastWorkingDay     *time.Time `json:"last_working_day,omitempty"`
	ExitInterviewDate  *time.Time `json:"exit_interview_date,omitempty"`
	ExitInterviewNotes string     `json:"exit_interview_notes,omitempty"`
}

// TechnicianMetrics contains business metrics
type TechnicianMetrics struct {
	TotalRepairsCompleted  int             `gorm:"default:0" json:"total_repairs_completed"`
	TotalHoursWorked       float64         `gorm:"default:0" json:"total_hours_worked"`
	TotalOvertimeHours     float64         `gorm:"default:0" json:"total_overtime_hours"`
	TotalRevenue           decimal.Decimal `sql:"type:decimal(12,2)" json:"total_revenue"`
	AverageTicketValue     decimal.Decimal `sql:"type:decimal(10,2)" json:"average_ticket_value"`
	CustomerSatisfaction   float64         `gorm:"default:0" json:"customer_satisfaction"`
	RepeatCustomerRate     float64         `gorm:"default:0" json:"repeat_customer_rate"`
	ReferralRate           float64         `gorm:"default:0" json:"referral_rate"`
	TrainingHoursCompleted float64         `gorm:"default:0" json:"training_hours_completed"`
	CertificationsEarned   int             `gorm:"default:0" json:"certifications_earned"`
	AwardsReceived         int             `gorm:"default:0" json:"awards_received"`
	SkillsAcquired         int             `gorm:"default:0" json:"skills_acquired"`
	ProjectsCompleted      int             `gorm:"default:0" json:"projects_completed"`
	TeamContribution       float64         `gorm:"default:0" json:"team_contribution"`
	InnovationScore        float64         `gorm:"default:0" json:"innovation_score"`
	CostSavings            decimal.Decimal `sql:"type:decimal(10,2)" json:"cost_savings"`
	WarrantyClaimRate      float64         `gorm:"default:0" json:"warranty_claim_rate"`
	RetentionScore         float64         `gorm:"default:0" json:"retention_score"`
}
