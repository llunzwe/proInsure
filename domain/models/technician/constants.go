package technician

// Technician Types
const (
	TechnicianTypeEmployee   = "employee"
	TechnicianTypeContractor = "contractor"
	TechnicianTypeFreelance  = "freelance"
	TechnicianTypeIntern     = "intern"
	TechnicianTypeApprentice = "apprentice"
	TechnicianTypePartTime   = "part_time"
	TechnicianTypeFullTime   = "full_time"
	TechnicianTypeOnCall     = "on_call"
	TechnicianTypeRemote     = "remote"
)

// Certification Levels
const (
	CertificationNone         = "none"
	CertificationJunior       = "junior"
	CertificationIntermediate = "intermediate"
	CertificationSenior       = "senior"
	CertificationExpert       = "expert"
	CertificationMaster       = "master"
	CertificationSpecialist   = "specialist"
	CertificationLead         = "lead"
)

// Skill Levels
const (
	SkillBeginner     = "beginner"
	SkillIntermediate = "intermediate"
	SkillAdvanced     = "advanced"
	SkillExpert       = "expert"
	SkillMaster       = "master"
)

// Specialization Areas
const (
	SpecializationScreen       = "screen_repair"
	SpecializationBattery      = "battery_replacement"
	SpecializationMotherboard  = "motherboard_repair"
	SpecializationWaterDamage  = "water_damage"
	SpecializationSoftware     = "software_issues"
	SpecializationCamera       = "camera_repair"
	SpecializationAudio        = "audio_components"
	SpecializationConnectivity = "connectivity_issues"
	SpecializationDiagnostics  = "diagnostics"
	SpecializationDataRecovery = "data_recovery"
	SpecializationMicroSolder  = "micro_soldering"
	SpecializationBoardLevel   = "board_level_repair"
)

// Status Types
const (
	StatusActive     = "active"
	StatusInactive   = "inactive"
	StatusOnLeave    = "on_leave"
	StatusSuspended  = "suspended"
	StatusTerminated = "terminated"
	StatusTraining   = "training"
	StatusProbation  = "probation"
	StatusOnboarding = "onboarding"
	StatusAvailable  = "available"
	StatusBusy       = "busy"
	StatusOnBreak    = "on_break"
	StatusOffDuty    = "off_duty"
)

// Shift Types
const (
	ShiftMorning   = "morning"
	ShiftAfternoon = "afternoon"
	ShiftEvening   = "evening"
	ShiftNight     = "night"
	ShiftWeekend   = "weekend"
	ShiftHoliday   = "holiday"
	ShiftFlexible  = "flexible"
	ShiftRotating  = "rotating"
	ShiftOnCall    = "on_call"
)

// Performance Ratings
const (
	PerformanceExcellent    = "excellent"
	PerformanceGood         = "good"
	PerformanceSatisfactory = "satisfactory"
	PerformanceNeedsImprove = "needs_improvement"
	PerformancePoor         = "poor"
	PerformanceOutstanding  = "outstanding"
)

// Training Types
const (
	TrainingOnboarding      = "onboarding"
	TrainingTechnical       = "technical"
	TrainingSafety          = "safety"
	TrainingCompliance      = "compliance"
	TrainingCustomerService = "customer_service"
	TrainingLeadership      = "leadership"
	TrainingCertification   = "certification"
	TrainingRefresher       = "refresher"
	TrainingAdvanced        = "advanced"
	TrainingSpecialized     = "specialized"
)

// Training Status
const (
	TrainingScheduled  = "scheduled"
	TrainingInProgress = "in_progress"
	TrainingCompleted  = "completed"
	TrainingFailed     = "failed"
	TrainingCancelled  = "cancelled"
	TrainingExpired    = "expired"
	TrainingPending    = "pending"
)

// Assignment Status
const (
	AssignmentPending    = "pending"
	AssignmentAccepted   = "accepted"
	AssignmentInProgress = "in_progress"
	AssignmentCompleted  = "completed"
	AssignmentCancelled  = "cancelled"
	AssignmentReassigned = "reassigned"
	AssignmentPaused     = "paused"
	AssignmentFailed     = "failed"
)

// Assignment Priority
const (
	PriorityLow       = "low"
	PriorityNormal    = "normal"
	PriorityHigh      = "high"
	PriorityUrgent    = "urgent"
	PriorityCritical  = "critical"
	PriorityEmergency = "emergency"
)

// Availability Status
const (
	AvailabilityAvailable    = "available"
	AvailabilityBusy         = "busy"
	AvailabilityOnAssignment = "on_assignment"
	AvailabilityOnBreak      = "on_break"
	AvailabilityOffDuty      = "off_duty"
	AvailabilityOnLeave      = "on_leave"
	AvailabilityEmergency    = "emergency_only"
)

// Leave Types
const (
	LeaveAnnual    = "annual"
	LeaveSick      = "sick"
	LeavePersonal  = "personal"
	LeaveMaternity = "maternity"
	LeavePaternity = "paternity"
	LeaveEmergency = "emergency"
	LeaveTraining  = "training"
	LeaveUnpaid    = "unpaid"
	LeaveStudy     = "study"
)

// Compensation Types
const (
	CompensationHourly     = "hourly"
	CompensationSalary     = "salary"
	CompensationCommission = "commission"
	CompensationContract   = "contract"
	CompensationPieceRate  = "piece_rate"
	CompensationMixed      = "mixed"
)

// Bonus Types
const (
	BonusPerformance          = "performance"
	BonusProductivity         = "productivity"
	BonusQuality              = "quality"
	BonusReferral             = "referral"
	BonusRetention            = "retention"
	BonusHoliday              = "holiday"
	BonusSpecialProject       = "special_project"
	BonusCustomerSatisfaction = "customer_satisfaction"
)

// Document Types
const (
	DocumentResume      = "resume"
	DocumentCertificate = "certificate"
	DocumentLicense     = "license"
	DocumentID          = "identification"
	DocumentContract    = "contract"
	DocumentNDA         = "nda"
	DocumentTraining    = "training_record"
	DocumentPerformance = "performance_review"
	DocumentBackground  = "background_check"
	DocumentInsurance   = "insurance"
)

// Background Check Status
const (
	BackgroundPending  = "pending"
	BackgroundClear    = "clear"
	BackgroundFlagged  = "flagged"
	BackgroundFailed   = "failed"
	BackgroundExpired  = "expired"
	BackgroundRenewing = "renewing"
)

// Equipment Status
const (
	EquipmentAssigned    = "assigned"
	EquipmentAvailable   = "available"
	EquipmentDamaged     = "damaged"
	EquipmentLost        = "lost"
	EquipmentReturned    = "returned"
	EquipmentMaintenance = "maintenance"
	EquipmentRetired     = "retired"
)

// Notification Preferences
const (
	NotificationEmail = "email"
	NotificationSMS   = "sms"
	NotificationPush  = "push"
	NotificationInApp = "in_app"
	NotificationAll   = "all"
	NotificationNone  = "none"
)

// Quality Metrics
const (
	QualityExcellent = 5
	QualityGood      = 4
	QualityAverage   = 3
	QualityPoor      = 2
	QualityTerrible  = 1
)

// Default Values
const (
	DefaultMaxAssignments        = 5
	DefaultMinRestHours          = 8
	DefaultMaxWorkHoursPerDay    = 10
	DefaultMaxWorkHoursPerWeek   = 40
	DefaultBreakDurationMinutes  = 30
	DefaultProbationDays         = 90
	DefaultNoticePeroidDays      = 14
	DefaultTrainingValidityDays  = 365
	DefaultCertificationYears    = 2
	DefaultPerformanceReviewDays = 90
	DefaultBackgroundCheckYears  = 2
	DefaultMinimumWage           = 15.0
	DefaultOvertimeMultiplier    = 1.5
	DefaultWeekendMultiplier     = 1.25
	DefaultHolidayMultiplier     = 2.0
	DefaultNightShiftMultiplier  = 1.3
	DefaultEmergencyMultiplier   = 2.5
)

// Thresholds
const (
	MinimumSkillScore       = 60
	MinimumPerformanceScore = 70
	MinimumCustomerRating   = 3.5
	MinimumEfficiencyRate   = 0.75
	MaximumErrorRate        = 0.05
	MaximumComplaintRate    = 0.02
	MinimumAttendanceRate   = 0.95
	MaximumReworkRate       = 0.1
	MinimumFirstTimeFixRate = 0.85
	MinimumTrainingScore    = 80
)

// Time Limits (in hours unless specified)
const (
	MaxShiftDuration      = 12
	MinShiftDuration      = 4
	MaxContinuousWork     = 6
	MinBreakAfterWork     = 1
	MaxWeeklyOvertime     = 20
	ResponseTimeUrgent    = 1
	ResponseTimeHigh      = 2
	ResponseTimeNormal    = 4
	MaxAssignmentDuration = 8
	MaxTravelTime         = 2
)
