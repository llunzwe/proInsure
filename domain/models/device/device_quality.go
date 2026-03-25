package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceQualityMetrics tracks overall quality assessments
type DeviceQualityMetrics struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AssessmentDate time.Time `json:"assessment_date"`

	// Build Quality
	BuildQualityScore   float64 `json:"build_quality_score"`             // 0-100
	StructuralIntegrity float64 `json:"structural_integrity"`            // 0-100
	AssemblyPrecision   float64 `json:"assembly_precision"`              // 0-100
	GapTolerances       string  `gorm:"type:json" json:"gap_tolerances"` // JSON object
	FitAndFinish        float64 `json:"fit_and_finish"`                  // 0-100

	// Material Quality
	MaterialGrade       string  `json:"material_grade"`       // premium, standard, budget
	MaterialDurability  float64 `json:"material_durability"`  // 0-100
	MaterialConsistency float64 `json:"material_consistency"` // 0-100
	SurfaceQuality      float64 `json:"surface_quality"`      // 0-100
	CoatingQuality      float64 `json:"coating_quality"`      // 0-100

	// Component Quality
	ComponentRatings string  `gorm:"type:json" json:"component_ratings"` // JSON object
	DisplayQuality   float64 `json:"display_quality"`                    // 0-100
	CameraQuality    float64 `json:"camera_quality"`                     // 0-100
	BatteryQuality   float64 `json:"battery_quality"`                    // 0-100
	ProcessorQuality float64 `json:"processor_quality"`                  // 0-100

	// Finish Quality
	PaintQuality  float64 `json:"paint_quality"`  // 0-100
	ScreenCoating float64 `json:"screen_coating"` // 0-100
	ButtonFeel    float64 `json:"button_feel"`    // 0-100
	PortAlignment float64 `json:"port_alignment"` // 0-100

	// Durability Testing
	DurabilityScore   float64 `json:"durability_score"`                   // 0-100
	DropTestResults   string  `gorm:"type:json" json:"drop_test_results"` // JSON object
	ScratchResistance float64 `json:"scratch_resistance"`                 // 0-100
	BendTestResults   string  `gorm:"type:json" json:"bend_test_results"` // JSON object
	WearResistance    float64 `json:"wear_resistance"`                    // 0-100

	// Quality Control
	QCPassRate    float64 `json:"qc_pass_rate"` // percentage
	QCInspections int     `json:"qc_inspections"`
	QCFailures    int     `json:"qc_failures"`
	QCCategories  string  `gorm:"type:json" json:"qc_categories"` // JSON array
	RetestCount   int     `json:"retest_count"`

	// Customer Perception
	CustomerRating     float64 `json:"customer_rating"`   // 0-5
	PerceivedQuality   float64 `json:"perceived_quality"` // 0-100
	CustomerComplaints int     `json:"customer_complaints"`
	ReturnRate         float64 `json:"return_rate"` // percentage

	// Quality Improvement
	ImprovementAreas string  `gorm:"type:json" json:"improvement_areas"` // JSON array
	QualityTrend     string  `json:"quality_trend"`                      // improving, stable, declining
	ImprovementScore float64 `json:"improvement_score"`                  // 0-100
	ActionItems      string  `gorm:"type:json" json:"action_items"`      // JSON array

	// Benchmarks
	IndustryBenchmark    float64 `json:"industry_benchmark"`                     // 0-100
	CompetitorComparison string  `gorm:"type:json" json:"competitor_comparison"` // JSON object
	CategoryRanking      int     `json:"category_ranking"`
	OverallQualityScore  float64 `json:"overall_quality_score"` // 0-100

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceDefectHistory tracks manufacturing defects and issues
type DeviceDefectHistory struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	DefectID string    `gorm:"uniqueIndex" json:"defect_id"`

	// Defect Information
	DefectType        string `json:"defect_type"` // manufacturing, design, component, software
	DefectCategory    string `json:"defect_category"`
	DefectDescription string `gorm:"type:text" json:"defect_description"`
	DefectLocation    string `json:"defect_location"`

	// Discovery Timeline
	ManufactureDate time.Time `json:"manufacture_date"`
	DiscoveryDate   time.Time `json:"discovery_date"`
	ReportedDate    time.Time `json:"reported_date"`
	DaysToDiscovery int       `json:"days_to_discovery"`
	DiscoveryMethod string    `json:"discovery_method"` // inspection, customer, testing

	// Severity Classification
	SeverityLevel        string `json:"severity_level"` // critical, major, minor, cosmetic
	SafetyImpact         bool   `json:"safety_impact"`
	FunctionalImpact     string `json:"functional_impact"`      // none, partial, complete
	UserExperienceImpact string `json:"user_experience_impact"` // high, medium, low

	// Batch Correlation
	BatchNumber     string  `json:"batch_number"`
	BatchDefectRate float64 `json:"batch_defect_rate"` // percentage
	AffectedUnits   int     `json:"affected_units"`
	BatchDateRange  string  `gorm:"type:json" json:"batch_date_range"` // JSON object

	// Serial Number Mapping
	SerialNumber       string `json:"serial_number"`
	ProductionLine     string `json:"production_line"`
	ManufacturingPlant string `json:"manufacturing_plant"`
	QCInspector        string `json:"qc_inspector"`

	// Resolution Tracking
	ResolutionStatus string     `json:"resolution_status"` // pending, in_progress, resolved
	ResolutionDate   *time.Time `json:"resolution_date"`
	ResolutionMethod string     `json:"resolution_method"` // repair, replace, refund
	ResolutionCost   float64    `json:"resolution_cost"`
	ResolutionTime   int        `json:"resolution_time"` // days

	// Warranty Claims
	UnderWarranty       bool    `json:"under_warranty"`
	WarrantyClaimNumber string  `json:"warranty_claim_number"`
	WarrantyCovered     bool    `json:"warranty_covered"`
	CustomerCost        float64 `json:"customer_cost"`

	// Known Issues Database
	KnownIssue    bool      `json:"known_issue"`
	IssueID       string    `json:"issue_id"`
	FirstReported time.Time `json:"first_reported"`
	TotalReports  int       `json:"total_reports"`

	// Cost Impact
	DirectCost              float64 `json:"direct_cost"`
	IndirectCost            float64 `json:"indirect_cost"`
	ReputationImpact        string  `json:"reputation_impact"`         // high, medium, low
	CustomerRetentionImpact float64 `json:"customer_retention_impact"` // percentage

	// Recall Risk
	RecallRisk         string `json:"recall_risk"` // high, medium, low, none
	RecallRecommended  bool   `json:"recall_recommended"`
	RegulatoryNotified bool   `json:"regulatory_notified"`
	PublicDisclosure   bool   `json:"public_disclosure"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceRecallStatus manages product recalls
type DeviceRecallStatus struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecallID string    `gorm:"uniqueIndex" json:"recall_id"`

	// Recall Information
	RecallType     string `json:"recall_type"` // mandatory, voluntary, safety
	RecallReason   string `gorm:"type:text" json:"recall_reason"`
	RecallSeverity string `json:"recall_severity"` // critical, major, minor
	RecallScope    string `json:"recall_scope"`    // global, regional, national

	// Notification Status
	NotificationSent     bool       `json:"notification_sent"`
	NotificationDate     *time.Time `json:"notification_date"`
	NotificationMethod   string     `gorm:"type:json" json:"notification_method"` // JSON array
	CustomerAcknowledged bool       `json:"customer_acknowledged"`
	AcknowledgeDate      *time.Time `json:"acknowledge_date"`

	// Compliance Tracking
	ComplianceStatus   string    `json:"compliance_status"` // pending, compliant, non_compliant
	ComplianceDeadline time.Time `json:"compliance_deadline"`
	ComplianceActions  string    `gorm:"type:json" json:"compliance_actions"` // JSON array
	RegulatoryBody     string    `json:"regulatory_body"`

	// Repair/Replacement
	ResolutionOption string     `json:"resolution_option"` // repair, replace, refund
	ResolutionStatus string     `json:"resolution_status"` // pending, scheduled, completed
	ScheduledDate    *time.Time `json:"scheduled_date"`
	CompletedDate    *time.Time `json:"completed_date"`
	ServiceLocation  string     `json:"service_location"`

	// Cost Tracking
	RecallCost           float64 `json:"recall_cost"`
	LaborCost            float64 `json:"labor_cost"`
	PartsCost            float64 `json:"parts_cost"`
	LogisticsCost        float64 `json:"logistics_cost"`
	CompensationProvided float64 `json:"compensation_provided"`

	// Safety Documentation
	SafetyHazard   string `gorm:"type:text" json:"safety_hazard"`
	InjuryReports  int    `json:"injury_reports"`
	PropertyDamage int    `json:"property_damage"`
	RiskMitigation string `gorm:"type:json" json:"risk_mitigation"` // JSON array

	// Regulatory Requirements
	RegulatoryFiling    bool       `json:"regulatory_filing"`
	FilingDate          *time.Time `json:"filing_date"`
	RegulatoryApproval  bool       `json:"regulatory_approval"`
	ComplianceDocuments string     `gorm:"type:json" json:"compliance_documents"` // JSON array

	// Voluntary Participation
	VoluntaryRecall   bool    `json:"voluntary_recall"`
	ParticipationRate float64 `json:"participation_rate"`                  // percentage
	IncentivesOffered string  `gorm:"type:json" json:"incentives_offered"` // JSON array

	// Effectiveness Metrics
	RecallEffectiveness float64 `json:"recall_effectiveness"` // percentage
	UnitsRecalled       int     `json:"units_recalled"`
	UnitsOutstanding    int     `json:"units_outstanding"`
	CompletionRate      float64 `json:"completion_rate"` // percentage

	// Legal Liability
	LegalAction      bool    `json:"legal_action"`
	LiabilityAmount  float64 `json:"liability_amount"`
	InsuranceClaim   bool    `json:"insurance_claim"`
	SettlementAmount float64 `json:"settlement_amount"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceTestResults documents quality testing outcomes
type DeviceTestResults struct {
	database.BaseModel
	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	TestBatchID string    `gorm:"index" json:"test_batch_id"`
	TestDate    time.Time `json:"test_date"`

	// Test Documentation
	TestProtocol   string `json:"test_protocol"`
	TestStandard   string `json:"test_standard"` // ISO, IEC, etc.
	TestLab        string `json:"test_lab"`
	TestTechnician string `json:"test_technician"`
	TestReport     string `gorm:"type:text" json:"test_report"`

	// Stress Testing
	StressTested       bool   `json:"stress_tested"`
	StressTestDuration int    `json:"stress_test_duration"` // hours
	StressTestCycles   int    `json:"stress_test_cycles"`
	StressTestPassed   bool   `json:"stress_test_passed"`
	FailurePoint       string `json:"failure_point"`

	// Drop Testing
	DropTestHeight      float64 `json:"drop_test_height"` // meters
	DropTestSurface     string  `json:"drop_test_surface"`
	DropTestRepetitions int     `json:"drop_test_repetitions"`
	DropDamage          string  `gorm:"type:json" json:"drop_damage"` // JSON array
	DropTestPassed      bool    `json:"drop_test_passed"`

	// Water Resistance
	WaterResistanceRating string  `json:"water_resistance_rating"` // IP rating
	SubmersionDepth       float64 `json:"submersion_depth"`        // meters
	SubmersionDuration    int     `json:"submersion_duration"`     // minutes
	WaterIngressDetected  bool    `json:"water_ingress_detected"`
	SealIntegrity         bool    `json:"seal_integrity"`

	// Temperature Testing
	MinTemperature    float64 `json:"min_temperature"` // celsius
	MaxTemperature    float64 `json:"max_temperature"` // celsius
	ThermalCycles     int     `json:"thermal_cycles"`
	ThermalShock      bool    `json:"thermal_shock"`
	TemperaturePassed bool    `json:"temperature_passed"`

	// Battery Safety
	BatteryTested       bool `json:"battery_tested"`
	OverchargeTesting   bool `json:"overcharge_testing"`
	ShortCircuitTesting bool `json:"short_circuit_testing"`
	PunctureTesting     bool `json:"puncture_testing"`
	BatterySafetyPassed bool `json:"battery_safety_passed"`

	// EMC Compliance
	EMCTested              bool    `json:"emc_tested"`
	RadiatedEmissions      float64 `json:"radiated_emissions"`  // dBμV/m
	ConductedEmissions     float64 `json:"conducted_emissions"` // dBμV
	ElectrostaticDischarge bool    `json:"electrostatic_discharge"`
	EMCCompliant           bool    `json:"emc_compliant"`

	// Performance Benchmarks
	BenchmarkScore int     `json:"benchmark_score"`
	CPUPerformance float64 `json:"cpu_performance"`
	GPUPerformance float64 `json:"gpu_performance"`
	MemorySpeed    float64 `json:"memory_speed"`
	StorageSpeed   float64 `json:"storage_speed"`

	// Reliability Testing
	MTBF             int     `json:"mtbf"`              // hours - Mean Time Between Failures
	MTTF             int     `json:"mttf"`              // hours - Mean Time To Failure
	FailureRate      float64 `json:"failure_rate"`      // percentage
	ReliabilityScore float64 `json:"reliability_score"` // 0-100

	// Accelerated Aging
	AgingTestDuration int     `json:"aging_test_duration"` // hours
	SimulatedYears    float64 `json:"simulated_years"`
	AgingDefects      string  `gorm:"type:json" json:"aging_defects"` // JSON array
	ExpectedLifespan  int     `json:"expected_lifespan"`              // months

	// Overall Results
	OverallTestResult    string  `json:"overall_test_result"` // pass, fail, conditional
	TestScore            float64 `json:"test_score"`          // 0-100
	CertificationGranted bool    `json:"certification_granted"`
	RetestRequired       bool    `json:"retest_required"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceQualityIncidents tracks quality issues and complaints
type DeviceQualityIncidents struct {
	database.BaseModel
	DeviceID     uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	IncidentID   string    `gorm:"uniqueIndex" json:"incident_id"`
	IncidentDate time.Time `json:"incident_date"`

	// Issue Reporting
	IssueType        string `json:"issue_type"` // defect, malfunction, safety, cosmetic
	IssueCategory    string `json:"issue_category"`
	IssueDescription string `gorm:"type:text" json:"issue_description"`
	ReportedBy       string `json:"reported_by"` // customer, qc, service_center
	ReportChannel    string `json:"report_channel"`

	// Customer Complaints
	ComplaintNumber     string    `json:"complaint_number"`
	CustomerID          uuid.UUID `json:"customer_id"`
	ComplaintSeverity   string    `json:"complaint_severity"` // high, medium, low
	CustomerImpact      string    `gorm:"type:text" json:"customer_impact"`
	CustomerExpectation string    `json:"customer_expectation"`

	// Investigation Status
	InvestigationStatus  string     `json:"investigation_status"` // pending, ongoing, completed
	InvestigatorAssigned string     `json:"investigator_assigned"`
	InvestigationStart   *time.Time `json:"investigation_start"`
	InvestigationEnd     *time.Time `json:"investigation_end"`
	FindingsReport       string     `gorm:"type:text" json:"findings_report"`

	// Root Cause Analysis
	RootCauseIdentified bool   `json:"root_cause_identified"`
	RootCause           string `gorm:"type:text" json:"root_cause"`
	ContributingFactors string `gorm:"type:json" json:"contributing_factors"` // JSON array
	SystematicIssue     bool   `json:"systematic_issue"`

	// Corrective Actions
	CorrectiveActions     string     `gorm:"type:json" json:"corrective_actions"` // JSON array
	ActionPlan            string     `gorm:"type:text" json:"action_plan"`
	ActionDeadline        *time.Time `json:"action_deadline"`
	ActionCompleted       bool       `json:"action_completed"`
	EffectivenessVerified bool       `json:"effectiveness_verified"`

	// Preventive Measures
	PreventiveMeasures string `gorm:"type:json" json:"preventive_measures"` // JSON array
	ProcessImprovement string `gorm:"type:text" json:"process_improvement"`
	TrainingRequired   bool   `json:"training_required"`
	ProcedureUpdated   bool   `json:"procedure_updated"`

	// Quality Trends
	RecurrenceCount  int    `json:"recurrence_count"`
	TrendIdentified  string `json:"trend_identified"`
	RelatedIncidents string `gorm:"type:json" json:"related_incidents"` // JSON array
	PatternAnalysis  string `gorm:"type:text" json:"pattern_analysis"`

	// Supplier Issues
	SupplierRelated   bool   `json:"supplier_related"`
	SupplierName      string `json:"supplier_name"`
	ComponentAffected string `json:"component_affected"`
	SupplierNotified  bool   `json:"supplier_notified"`
	SupplierResponse  string `gorm:"type:text" json:"supplier_response"`

	// Cost Tracking
	InvestigationCost float64 `json:"investigation_cost"`
	ResolutionCost    float64 `json:"resolution_cost"`
	CompensationCost  float64 `json:"compensation_cost"`
	TotalQualityCost  float64 `json:"total_quality_cost"`

	// Continuous Improvement
	LessonLearned        string `gorm:"type:text" json:"lesson_learned"`
	BestPracticeShared   bool   `json:"best_practice_shared"`
	KnowledgeBaseUpdated bool   `json:"knowledge_base_updated"`
	ImprovementMetrics   string `gorm:"type:json" json:"improvement_metrics"` // JSON object

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceCertificationStatus tracks compliance and certifications
type DeviceCertificationStatus struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Safety Certifications
	SafetyCertifications string `gorm:"type:json" json:"safety_certifications"` // JSON array
	ULCertified          bool   `json:"ul_certified"`
	CEMarking            bool   `json:"ce_marking"`
	FCCCompliant         bool   `json:"fcc_compliant"`
	SafetyStandards      string `gorm:"type:json" json:"safety_standards"` // JSON array

	// Quality Certifications
	QualityCertifications string `gorm:"type:json" json:"quality_certifications"` // JSON array
	ISO9001Certified      bool   `json:"iso9001_certified"`
	ISO14001Certified     bool   `json:"iso14001_certified"`
	QualityStandards      string `gorm:"type:json" json:"quality_standards"` // JSON array

	// Environmental Certifications
	EnvironmentalCerts  string `gorm:"type:json" json:"environmental_certs"` // JSON array
	RoHSCompliant       bool   `json:"rohs_compliant"`
	REACHCompliant      bool   `json:"reach_compliant"`
	EnergyStarCertified bool   `json:"energy_star_certified"`
	EPEATRegistered     string `json:"epeat_registered"` // gold, silver, bronze

	// Industry Certifications
	IndustryCerts   string `gorm:"type:json" json:"industry_certs"` // JSON array
	MILSTDCompliant bool   `json:"milstd_compliant"`
	MedicalGrade    bool   `json:"medical_grade"`
	AutomotiveGrade bool   `json:"automotive_grade"`
	AerospaceGrade  bool   `json:"aerospace_grade"`

	// Regional Certifications
	RegionalCerts        string `gorm:"type:json" json:"regional_certs"`         // JSON array
	CountrySpecificCerts string `gorm:"type:json" json:"country_specific_certs"` // JSON object
	ImportApprovals      string `gorm:"type:json" json:"import_approvals"`       // JSON array
	ExportLicenses       string `gorm:"type:json" json:"export_licenses"`        // JSON array

	// Renewal Tracking
	CertificationExpiry string     `gorm:"type:json" json:"certification_expiry"` // JSON object
	RenewalSchedule     string     `gorm:"type:json" json:"renewal_schedule"`     // JSON array
	RenewalInProgress   bool       `json:"renewal_in_progress"`
	LastRenewalDate     *time.Time `json:"last_renewal_date"`
	NextRenewalDate     *time.Time `json:"next_renewal_date"`

	// Audit Results
	LastAuditDate     *time.Time `json:"last_audit_date"`
	AuditResult       string     `json:"audit_result"`                        // pass, fail, conditional
	AuditFindings     string     `gorm:"type:json" json:"audit_findings"`     // JSON array
	CorrectiveActions string     `gorm:"type:json" json:"corrective_actions"` // JSON array
	NextAuditDate     *time.Time `json:"next_audit_date"`

	// Non-Conformance
	NonConformances     int    `json:"non_conformances"`
	OpenNonConformances int    `json:"open_non_conformances"`
	NCRList             string `gorm:"type:json" json:"ncr_list"` // JSON array - Non-Conformance Reports
	CAPAStatus          string `json:"capa_status"`               // Corrective and Preventive Action

	// Certification Costs
	InitialCertCost float64 `json:"initial_cert_cost"`
	RenewalCost     float64 `json:"renewal_cost"`
	AuditCost       float64 `json:"audit_cost"`
	MaintenanceCost float64 `json:"maintenance_cost"`
	TotalCertCost   float64 `json:"total_cert_cost"`

	// Compliance Verification
	ComplianceStatus string    `json:"compliance_status"`                // compliant, partial, non_compliant
	ComplianceScore  float64   `json:"compliance_score"`                 // 0-100
	ComplianceGaps   string    `gorm:"type:json" json:"compliance_gaps"` // JSON array
	VerificationDate time.Time `json:"verification_date"`
	VerifiedBy       string    `json:"verified_by"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceQualityMetrics
func (dqm *DeviceQualityMetrics) IsHighQuality() bool {
	return dqm.OverallQualityScore >= 80 && dqm.BuildQualityScore >= 75 &&
		dqm.QCPassRate >= 95
}

func (dqm *DeviceQualityMetrics) HasQualityIssues() bool {
	return dqm.QualityTrend == "declining" || dqm.CustomerComplaints > 5 ||
		dqm.ReturnRate > 5 || dqm.QCFailures > 3
}

func (dqm *DeviceQualityMetrics) GetQualityGrade() string {
	if dqm.OverallQualityScore >= 90 {
		return "A"
	} else if dqm.OverallQualityScore >= 80 {
		return "B"
	} else if dqm.OverallQualityScore >= 70 {
		return "C"
	} else if dqm.OverallQualityScore >= 60 {
		return "D"
	}
	return "F"
}

func (dqm *DeviceQualityMetrics) MeetsBenchmark() bool {
	return dqm.OverallQualityScore >= dqm.IndustryBenchmark
}

func (dqm *DeviceQualityMetrics) GetWeakestComponent() string {
	components := map[string]float64{
		"display":   dqm.DisplayQuality,
		"camera":    dqm.CameraQuality,
		"battery":   dqm.BatteryQuality,
		"processor": dqm.ProcessorQuality,
	}

	minScore := 100.0
	weakest := ""
	for name, score := range components {
		if score < minScore {
			minScore = score
			weakest = name
		}
	}
	return weakest
}

// Methods for DeviceDefectHistory
func (ddh *DeviceDefectHistory) IsCriticalDefect() bool {
	return ddh.SeverityLevel == "critical" || ddh.SafetyImpact ||
		ddh.RecallRisk == "high"
}

func (ddh *DeviceDefectHistory) IsResolved() bool {
	return ddh.ResolutionStatus == "resolved" && ddh.ResolutionDate != nil
}

func (ddh *DeviceDefectHistory) IsSystematic() bool {
	return ddh.KnownIssue || ddh.BatchDefectRate > 5 ||
		ddh.TotalReports > 10
}

func (ddh *DeviceDefectHistory) NeedsRecall() bool {
	return ddh.RecallRecommended || ddh.RecallRisk == "high" ||
		(ddh.SafetyImpact && ddh.AffectedUnits > 100)
}

func (ddh *DeviceDefectHistory) GetTotalCost() float64 {
	return ddh.DirectCost + ddh.IndirectCost + ddh.ResolutionCost + ddh.CustomerCost
}

// Methods for DeviceRecallStatus
func (drs *DeviceRecallStatus) IsActive() bool {
	return drs.ResolutionStatus != "completed" && drs.ComplianceStatus != "compliant"
}

func (drs *DeviceRecallStatus) IsCompliant() bool {
	return drs.ComplianceStatus == "compliant" &&
		time.Now().Before(drs.ComplianceDeadline)
}

func (drs *DeviceRecallStatus) IsSafetyRelated() bool {
	return drs.RecallType == "safety" || drs.InjuryReports > 0 ||
		drs.RecallSeverity == "critical"
}

func (drs *DeviceRecallStatus) GetTotalCost() float64 {
	return drs.RecallCost + drs.LaborCost + drs.PartsCost +
		drs.LogisticsCost + drs.CompensationProvided +
		drs.LiabilityAmount + drs.SettlementAmount
}

func (drs *DeviceRecallStatus) GetEffectiveness() float64 {
	if drs.UnitsRecalled+drs.UnitsOutstanding > 0 {
		return float64(drs.UnitsRecalled) / float64(drs.UnitsRecalled+drs.UnitsOutstanding) * 100
	}
	return 0
}

// Methods for DeviceTestResults
func (dtr *DeviceTestResults) PassedAllTests() bool {
	return dtr.OverallTestResult == "pass" && dtr.CertificationGranted &&
		!dtr.RetestRequired
}

func (dtr *DeviceTestResults) HasSafetyIssues() bool {
	return !dtr.BatterySafetyPassed || dtr.WaterIngressDetected ||
		!dtr.TemperaturePassed || !dtr.EMCCompliant
}

func (dtr *DeviceTestResults) IsDurable() bool {
	return dtr.DropTestPassed && dtr.StressTestPassed &&
		dtr.ReliabilityScore >= 80 && dtr.MTBF > 10000
}

func (dtr *DeviceTestResults) GetQualityScore() float64 {
	if dtr.TestScore > 0 {
		return dtr.TestScore
	}
	// Calculate based on test results
	score := 100.0
	if !dtr.StressTestPassed {
		score -= 15
	}
	if !dtr.DropTestPassed {
		score -= 15
	}
	if !dtr.BatterySafetyPassed {
		score -= 20
	}
	if dtr.WaterIngressDetected {
		score -= 10
	}
	if !dtr.TemperaturePassed {
		score -= 10
	}
	if !dtr.EMCCompliant {
		score -= 10
	}
	if score < 0 {
		score = 0
	}
	return score
}

func (dtr *DeviceTestResults) MeetsStandards() bool {
	return dtr.EMCCompliant && dtr.BatterySafetyPassed &&
		dtr.TestStandard != "" && dtr.CertificationGranted
}

// Methods for DeviceQualityIncidents
func (dqi *DeviceQualityIncidents) IsSerious() bool {
	return dqi.ComplaintSeverity == "high" || dqi.SystematicIssue ||
		dqi.RecurrenceCount > 3
}

func (dqi *DeviceQualityIncidents) IsInvestigated() bool {
	return dqi.InvestigationStatus == "completed" && dqi.RootCauseIdentified &&
		dqi.ActionCompleted
}

func (dqi *DeviceQualityIncidents) HasSupplierIssue() bool {
	return dqi.SupplierRelated && dqi.ComponentAffected != ""
}

func (dqi *DeviceQualityIncidents) RequiresTraining() bool {
	return dqi.TrainingRequired || dqi.RootCause == "human_error" ||
		dqi.RecurrenceCount > 2
}

func (dqi *DeviceQualityIncidents) GetResolutionProgress() float64 {
	progress := 0.0
	if dqi.InvestigationStatus == "completed" {
		progress += 25
	}
	if dqi.RootCauseIdentified {
		progress += 25
	}
	if dqi.ActionCompleted {
		progress += 25
	}
	if dqi.EffectivenessVerified {
		progress += 25
	}
	return progress
}

// Methods for DeviceCertificationStatus
func (dcs *DeviceCertificationStatus) IsFullyCertified() bool {
	return dcs.ComplianceStatus == "compliant" && dcs.ComplianceScore >= 95 &&
		dcs.OpenNonConformances == 0
}

func (dcs *DeviceCertificationStatus) HasSafetyCertification() bool {
	return dcs.ULCertified || dcs.CEMarking || dcs.FCCCompliant
}

func (dcs *DeviceCertificationStatus) IsEnvironmentallyCompliant() bool {
	return dcs.RoHSCompliant && dcs.REACHCompliant &&
		(dcs.EnergyStarCertified || dcs.EPEATRegistered != "")
}

func (dcs *DeviceCertificationStatus) NeedsRenewal() bool {
	if dcs.NextRenewalDate != nil {
		return time.Now().After(dcs.NextRenewalDate.AddDate(0, -3, 0)) // 3 months before expiry
	}
	return dcs.RenewalInProgress
}

func (dcs *DeviceCertificationStatus) GetCertificationLevel() string {
	score := 0
	if dcs.ULCertified {
		score++
	}
	if dcs.CEMarking {
		score++
	}
	if dcs.FCCCompliant {
		score++
	}
	if dcs.ISO9001Certified {
		score++
	}
	if dcs.ISO14001Certified {
		score++
	}
	if dcs.RoHSCompliant {
		score++
	}
	if dcs.REACHCompliant {
		score++
	}
	if dcs.EnergyStarCertified {
		score++
	}

	if score >= 7 {
		return "premium"
	} else if score >= 5 {
		return "standard"
	} else if score >= 3 {
		return "basic"
	}
	return "minimal"
}
