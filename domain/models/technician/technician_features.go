package technician

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TechnicianAssignment represents work assignments for technicians
type TechnicianAssignment struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID      uuid.UUID  `gorm:"type:uuid;not null" json:"technician_id"`
	RepairBookingID   *uuid.UUID `gorm:"type:uuid" json:"repair_booking_id,omitempty"`
	ServiceRequestID  *uuid.UUID `gorm:"type:uuid" json:"service_request_id,omitempty"`
	AssignmentNumber  string     `gorm:"uniqueIndex;not null" json:"assignment_number"`
	Type              string     `gorm:"not null" json:"type"`
	Priority          string     `gorm:"not null;default:'normal'" json:"priority"`
	Status            string     `gorm:"not null;default:'pending'" json:"status"`
	Title             string     `gorm:"not null" json:"title"`
	Description       string     `json:"description"`
	Location          string     `json:"location"`
	CustomerName      string     `json:"customer_name"`
	CustomerPhone     string     `json:"customer_phone"`
	ScheduledStart    time.Time  `gorm:"not null" json:"scheduled_start"`
	ScheduledEnd      time.Time  `gorm:"not null" json:"scheduled_end"`
	ActualStart       *time.Time `json:"actual_start,omitempty"`
	ActualEnd         *time.Time `json:"actual_end,omitempty"`
	EstimatedDuration int        `json:"estimated_duration"` // minutes
	ActualDuration    int        `json:"actual_duration"`    // minutes
	TravelTime        int        `json:"travel_time"`        // minutes
	TravelDistance    float64    `json:"travel_distance"`    // km
	CompletionNotes   string     `json:"completion_notes,omitempty"`
	CustomerSignature string     `json:"customer_signature,omitempty"`
	Rating            int        `json:"rating,omitempty"`
	Feedback          string     `json:"feedback,omitempty"`
	IsReassigned      bool       `gorm:"default:false" json:"is_reassigned"`
	ReassignedFrom    *uuid.UUID `gorm:"type:uuid" json:"reassigned_from,omitempty"`
	ReassignReason    string     `json:"reassign_reason,omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TechnicianSchedule represents work schedules
type TechnicianSchedule struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID  uuid.UUID  `gorm:"type:uuid;not null" json:"technician_id"`
	ScheduleDate  time.Time  `gorm:"not null;index" json:"schedule_date"`
	ShiftType     string     `gorm:"not null" json:"shift_type"`
	StartTime     time.Time  `gorm:"not null" json:"start_time"`
	EndTime       time.Time  `gorm:"not null" json:"end_time"`
	BreakStart    *time.Time `json:"break_start,omitempty"`
	BreakEnd      *time.Time `json:"break_end,omitempty"`
	Status        string     `gorm:"default:'scheduled'" json:"status"`
	CheckInTime   *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime  *time.Time `json:"check_out_time,omitempty"`
	ActualHours   float64    `json:"actual_hours"`
	OvertimeHours float64    `json:"overtime_hours"`
	Location      string     `json:"location,omitempty"`
	IsRemote      bool       `gorm:"default:false" json:"is_remote"`
	IsOnCall      bool       `gorm:"default:false" json:"is_on_call"`
	Notes         string     `json:"notes,omitempty"`
	ApprovedBy    *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TechnicianSkill represents individual skills and proficiency
type TechnicianSkill struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID      uuid.UUID  `gorm:"type:uuid;not null" json:"technician_id"`
	SkillName         string     `gorm:"not null" json:"skill_name"`
	SkillCategory     string     `gorm:"not null" json:"skill_category"`
	ProficiencyLevel  string     `gorm:"not null" json:"proficiency_level"`
	YearsExperience   float64    `json:"years_experience"`
	LastUsedDate      *time.Time `json:"last_used_date,omitempty"`
	AssessmentScore   float64    `json:"assessment_score"`
	AssessmentDate    *time.Time `json:"assessment_date,omitempty"`
	CertificationName string     `json:"certification_name,omitempty"`
	CertificationDate *time.Time `json:"certification_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
	TrainingRequired  bool       `gorm:"default:false" json:"training_required"`
	IsVerified        bool       `gorm:"default:false" json:"is_verified"`
	VerifiedBy        *uuid.UUID `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TechnicianTraining represents training records
type TechnicianTraining struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID      uuid.UUID       `gorm:"type:uuid;not null" json:"technician_id"`
	TrainingName      string          `gorm:"not null" json:"training_name"`
	TrainingType      string          `gorm:"not null" json:"training_type"`
	Provider          string          `json:"provider"`
	Instructor        string          `json:"instructor,omitempty"`
	Description       string          `json:"description"`
	Duration          int             `json:"duration"` // hours
	StartDate         time.Time       `json:"start_date"`
	EndDate           time.Time       `json:"end_date"`
	Status            string          `gorm:"default:'scheduled'" json:"status"`
	CompletionRate    float64         `json:"completion_rate"`
	TestScore         float64         `json:"test_score,omitempty"`
	PassingScore      float64         `json:"passing_score"`
	Attempts          int             `gorm:"default:0" json:"attempts"`
	CertificateNumber string          `json:"certificate_number,omitempty"`
	CertificateURL    string          `json:"certificate_url,omitempty"`
	ExpiryDate        *time.Time      `json:"expiry_date,omitempty"`
	Cost              decimal.Decimal `sql:"type:decimal(10,2)" json:"cost"`
	IsMandatory       bool            `gorm:"default:false" json:"is_mandatory"`
	IsRecurring       bool            `gorm:"default:false" json:"is_recurring"`
	RecurrenceMonths  int             `json:"recurrence_months,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TechnicianCertification represents professional certifications
type TechnicianCertification struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID        uuid.UUID       `gorm:"type:uuid;not null" json:"technician_id"`
	CertificationName   string          `gorm:"not null" json:"certification_name"`
	CertificationBody   string          `gorm:"not null" json:"certification_body"`
	CertificationNumber string          `gorm:"uniqueIndex;not null" json:"certification_number"`
	Level               string          `json:"level,omitempty"`
	IssueDate           time.Time       `gorm:"not null" json:"issue_date"`
	ExpiryDate          *time.Time      `json:"expiry_date,omitempty"`
	Status              string          `gorm:"default:'active'" json:"status"`
	VerificationURL     string          `json:"verification_url,omitempty"`
	DocumentURL         string          `json:"document_url,omitempty"`
	RenewalRequired     bool            `gorm:"default:false" json:"renewal_required"`
	RenewalDate         *time.Time      `json:"renewal_date,omitempty"`
	RenewalCost         decimal.Decimal `sql:"type:decimal(10,2)" json:"renewal_cost"`
	TrainingHours       int             `json:"training_hours"`
	ExamScore           float64         `json:"exam_score,omitempty"`
	IsVerified          bool            `gorm:"default:false" json:"is_verified"`
	VerifiedAt          *time.Time      `json:"verified_at,omitempty"`
	CreatedAt           time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TechnicianPerformanceReview represents performance evaluations
type TechnicianPerformanceReview struct {
	ID                        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID              uuid.UUID  `gorm:"type:uuid;not null" json:"technician_id"`
	ReviewerID                uuid.UUID  `gorm:"type:uuid;not null" json:"reviewer_id"`
	ReviewPeriod              string     `gorm:"not null" json:"review_period"`
	ReviewDate                time.Time  `gorm:"not null" json:"review_date"`
	OverallRating             string     `gorm:"not null" json:"overall_rating"`
	TechnicalSkills           float64    `json:"technical_skills"`
	CustomerService           float64    `json:"customer_service"`
	Teamwork                  float64    `json:"teamwork"`
	Communication             float64    `json:"communication"`
	ProblemSolving            float64    `json:"problem_solving"`
	Reliability               float64    `json:"reliability"`
	Initiative                float64    `json:"initiative"`
	QualityOfWork             float64    `json:"quality_of_work"`
	Productivity              float64    `json:"productivity"`
	SafetyCompliance          float64    `json:"safety_compliance"`
	Strengths                 string     `json:"strengths"`
	AreasForImprovement       string     `json:"areas_for_improvement"`
	Goals                     string     `json:"goals"`
	ActionPlan                string     `json:"action_plan,omitempty"`
	EmployeeComments          string     `json:"employee_comments,omitempty"`
	RecommendedPromotion      bool       `gorm:"default:false" json:"recommended_promotion"`
	RecommendedSalaryIncrease bool       `gorm:"default:false" json:"recommended_salary_increase"`
	IncreasePercentage        float64    `json:"increase_percentage,omitempty"`
	NextReviewDate            time.Time  `json:"next_review_date"`
	IsAcknowledged            bool       `gorm:"default:false" json:"is_acknowledged"`
	AcknowledgedAt            *time.Time `json:"acknowledged_at,omitempty"`
	CreatedAt                 time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                 time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TechnicianLeave represents leave requests and history
type TechnicianLeave struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID     uuid.UUID  `gorm:"type:uuid;not null" json:"technician_id"`
	LeaveType        string     `gorm:"not null" json:"leave_type"`
	StartDate        time.Time  `gorm:"not null" json:"start_date"`
	EndDate          time.Time  `gorm:"not null" json:"end_date"`
	Days             float64    `json:"days"`
	Reason           string     `gorm:"not null" json:"reason"`
	Status           string     `gorm:"default:'pending'" json:"status"`
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	RejectionReason  string     `json:"rejection_reason,omitempty"`
	DocumentURL      string     `json:"document_url,omitempty"`
	IsPaid           bool       `gorm:"default:false" json:"is_paid"`
	IsEmergency      bool       `gorm:"default:false" json:"is_emergency"`
	HandoverNotes    string     `json:"handover_notes,omitempty"`
	BackupTechnician *uuid.UUID `gorm:"type:uuid" json:"backup_technician,omitempty"`
	ReturnDate       *time.Time `json:"return_date,omitempty"`
	ExtensionDays    float64    `json:"extension_days,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TechnicianTimesheet represents time tracking
type TechnicianTimesheet struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID      uuid.UUID  `gorm:"type:uuid;not null" json:"technician_id"`
	Date              time.Time  `gorm:"not null;index" json:"date"`
	ClockIn           time.Time  `gorm:"not null" json:"clock_in"`
	ClockOut          *time.Time `json:"clock_out,omitempty"`
	BreakStart        *time.Time `json:"break_start,omitempty"`
	BreakEnd          *time.Time `json:"break_end,omitempty"`
	TotalHours        float64    `json:"total_hours"`
	RegularHours      float64    `json:"regular_hours"`
	OvertimeHours     float64    `json:"overtime_hours"`
	BreakHours        float64    `json:"break_hours"`
	Location          string     `json:"location,omitempty"`
	LocationLatitude  float64    `json:"location_latitude,omitempty"`
	LocationLongitude float64    `json:"location_longitude,omitempty"`
	IsRemote          bool       `gorm:"default:false" json:"is_remote"`
	IsHoliday         bool       `gorm:"default:false" json:"is_holiday"`
	IsWeekend         bool       `gorm:"default:false" json:"is_weekend"`
	ShiftDifferential float64    `json:"shift_differential,omitempty"`
	Notes             string     `json:"notes,omitempty"`
	ApprovedBy        *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	Status            string     `gorm:"default:'pending'" json:"status"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TechnicianEquipment represents equipment assigned to technicians
type TechnicianEquipment struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID     uuid.UUID       `gorm:"type:uuid;not null" json:"technician_id"`
	EquipmentName    string          `gorm:"not null" json:"equipment_name"`
	EquipmentType    string          `gorm:"not null" json:"equipment_type"`
	SerialNumber     string          `gorm:"uniqueIndex;not null" json:"serial_number"`
	AssetTag         string          `gorm:"uniqueIndex" json:"asset_tag"`
	Brand            string          `json:"brand,omitempty"`
	Model            string          `json:"model,omitempty"`
	Status           string          `gorm:"default:'assigned'" json:"status"`
	AssignedDate     time.Time       `gorm:"not null" json:"assigned_date"`
	ReturnDate       *time.Time      `json:"return_date,omitempty"`
	Condition        string          `gorm:"not null" json:"condition"`
	Notes            string          `json:"notes,omitempty"`
	Value            decimal.Decimal `sql:"type:decimal(10,2)" json:"value"`
	DepreciatedValue decimal.Decimal `sql:"type:decimal(10,2)" json:"depreciated_value"`
	MaintenanceDue   *time.Time      `json:"maintenance_due,omitempty"`
	LastMaintenance  *time.Time      `json:"last_maintenance,omitempty"`
	IsCompanyOwned   bool            `gorm:"default:true" json:"is_company_owned"`
	RequiresReturn   bool            `gorm:"default:true" json:"requires_return"`
	LostDamageFee    decimal.Decimal `sql:"type:decimal(10,2)" json:"lost_damage_fee"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TechnicianDocument represents documents associated with technicians
type TechnicianDocument struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID   uuid.UUID      `gorm:"type:uuid;not null" json:"technician_id"`
	DocumentType   string         `gorm:"not null" json:"document_type"`
	DocumentName   string         `gorm:"not null" json:"document_name"`
	DocumentNumber string         `json:"document_number,omitempty"`
	DocumentURL    string         `gorm:"not null" json:"document_url"`
	IssueDate      *time.Time     `json:"issue_date,omitempty"`
	ExpiryDate     *time.Time     `json:"expiry_date,omitempty"`
	Status         string         `gorm:"default:'active'" json:"status"`
	IsVerified     bool           `gorm:"default:false" json:"is_verified"`
	VerifiedBy     *uuid.UUID     `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt     *time.Time     `json:"verified_at,omitempty"`
	IsConfidential bool           `gorm:"default:false" json:"is_confidential"`
	Notes          string         `json:"notes,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TechnicianIncident represents incidents involving technicians
type TechnicianIncident struct {
	ID                uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	TechnicianID      uuid.UUID       `gorm:"type:uuid;not null" json:"technician_id"`
	IncidentNumber    string          `gorm:"uniqueIndex;not null" json:"incident_number"`
	IncidentType      string          `gorm:"not null" json:"incident_type"`
	IncidentDate      time.Time       `gorm:"not null" json:"incident_date"`
	Location          string          `json:"location"`
	Description       string          `gorm:"not null" json:"description"`
	Severity          string          `gorm:"not null" json:"severity"`
	ReportedBy        uuid.UUID       `gorm:"type:uuid;not null" json:"reported_by"`
	Witnesses         string          `json:"witnesses,omitempty"` // JSON array
	InjuryOccurred    bool            `gorm:"default:false" json:"injury_occurred"`
	InjuryDescription string          `json:"injury_description,omitempty"`
	MedicalAttention  bool            `gorm:"default:false" json:"medical_attention"`
	PropertyDamage    bool            `gorm:"default:false" json:"property_damage"`
	DamageDescription string          `json:"damage_description,omitempty"`
	DamageCost        decimal.Decimal `sql:"type:decimal(10,2)" json:"damage_cost"`
	Investigation     string          `json:"investigation,omitempty"`
	RootCause         string          `json:"root_cause,omitempty"`
	CorrectiveAction  string          `json:"corrective_action,omitempty"`
	PreventiveAction  string          `json:"preventive_action,omitempty"`
	Status            string          `gorm:"default:'open'" json:"status"`
	Resolution        string          `json:"resolution,omitempty"`
	ResolvedAt        *time.Time      `json:"resolved_at,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}
