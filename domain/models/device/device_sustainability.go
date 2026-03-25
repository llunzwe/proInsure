package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceCarbonFootprint tracks carbon emissions throughout device lifecycle
type DeviceCarbonFootprint struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	CalculationDate time.Time `json:"calculation_date"`

	// Manufacturing Emissions
	ManufacturingEmissions float64 `json:"manufacturing_emissions"` // kg CO2e
	MaterialExtraction     float64 `json:"material_extraction"`     // kg CO2e
	ComponentProduction    float64 `json:"component_production"`    // kg CO2e
	AssemblyEmissions      float64 `json:"assembly_emissions"`      // kg CO2e
	PackagingEmissions     float64 `json:"packaging_emissions"`     // kg CO2e

	// Transportation Emissions
	TransportationEmissions float64 `json:"transportation_emissions"` // kg CO2e
	ShippingDistance        float64 `json:"shipping_distance"`        // km
	TransportMode           string  `json:"transport_mode"`           // air, sea, road, rail
	DistributionEmissions   float64 `json:"distribution_emissions"`   // kg CO2e
	LastMileDelivery        float64 `json:"last_mile_delivery"`       // kg CO2e

	// Usage Phase Emissions
	UsageEmissions     float64 `json:"usage_emissions"`      // kg CO2e
	DailyEnergyUsage   float64 `json:"daily_energy_usage"`   // kWh
	ElectricityGridMix string  `json:"electricity_grid_mix"` // renewable percentage
	ProjectedLifespan  int     `json:"projected_lifespan"`   // months
	AnnualEmissions    float64 `json:"annual_emissions"`     // kg CO2e per year

	// End-of-Life Emissions
	EndOfLifeEmissions    float64 `json:"end_of_life_emissions"`  // kg CO2e
	RecyclingEmissions    float64 `json:"recycling_emissions"`    // kg CO2e
	LandfillEmissions     float64 `json:"landfill_emissions"`     // kg CO2e
	IncinerationEmissions float64 `json:"incineration_emissions"` // kg CO2e
	RefurbishmentSavings  float64 `json:"refurbishment_savings"`  // kg CO2e saved

	// Carbon Offset Tracking
	CarbonOffsetsApplied float64 `json:"carbon_offsets_applied"`           // kg CO2e
	OffsetProjects       string  `gorm:"type:json" json:"offset_projects"` // JSON array
	OffsetCertification  string  `json:"offset_certification"`
	OffsetVerified       bool    `json:"offset_verified"`
	OffsetCost           float64 `json:"offset_cost"`

	// Renewable Energy Usage
	RenewableEnergyUsed float64 `json:"renewable_energy_used"` // percentage
	SolarPowerUsage     float64 `json:"solar_power_usage"`     // kWh
	WindPowerUsage      float64 `json:"wind_power_usage"`      // kWh
	HydroPowerUsage     float64 `json:"hydro_power_usage"`     // kWh
	GreenEnergyCredits  float64 `json:"green_energy_credits"`

	// Carbon Neutrality Progress
	CarbonNeutralTarget   *time.Time `json:"carbon_neutral_target"`
	CurrentNeutralityRate float64    `json:"current_neutrality_rate"` // percentage
	NeutralityGap         float64    `json:"neutrality_gap"`          // kg CO2e
	ImprovementRate       float64    `json:"improvement_rate"`        // percentage per year

	// Emission Reduction Targets
	ReductionTarget  float64 `json:"reduction_target"` // percentage
	BaselineYear     int     `json:"baseline_year"`
	TargetYear       int     `json:"target_year"`
	CurrentReduction float64 `json:"current_reduction"` // percentage achieved

	// Carbon Credit Integration
	CarbonCreditsOwned float64 `json:"carbon_credits_owned"` // tonnes CO2e
	CreditMarketValue  float64 `json:"credit_market_value"`
	CreditRegistry     string  `json:"credit_registry"`
	CreditVerification string  `json:"credit_verification"`

	// Environmental Impact Score
	TotalCarbon        float64 `json:"total_carbon"`        // kg CO2e
	CarbonIntensity    float64 `json:"carbon_intensity"`    // kg CO2e per year
	EnvironmentalScore float64 `json:"environmental_score"` // 0-100
	IndustryAverage    float64 `json:"industry_average"`    // kg CO2e
	PerformanceRating  string  `json:"performance_rating"`  // A-F

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceRecyclingScore evaluates device recyclability and circular economy contribution
type DeviceRecyclingScore struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AssessmentDate time.Time `json:"assessment_date"`

	// Material Recyclability
	RecyclabilityPercentage float64 `json:"recyclability_percentage"`
	RecyclableMaterials     string  `gorm:"type:json" json:"recyclable_materials"` // JSON object
	NonRecyclableParts      string  `gorm:"type:json" json:"non_recyclable_parts"` // JSON array
	MaterialComplexity      int     `json:"material_complexity"`                   // number of different materials
	SeparationDifficulty    string  `json:"separation_difficulty"`                 // easy, medium, hard

	// Hazardous Materials
	HazardousContent   float64 `json:"hazardous_content"`                    // percentage
	HazardousMaterials string  `gorm:"type:json" json:"hazardous_materials"` // JSON array
	LeadContent        float64 `json:"lead_content"`                         // mg
	MercuryContent     float64 `json:"mercury_content"`                      // mg
	CadmiumContent     float64 `json:"cadmium_content"`                      // mg

	// Precious Metal Recovery
	PreciousMetalValue float64 `json:"precious_metal_value"`
	GoldContent        float64 `json:"gold_content"`                         // mg
	SilverContent      float64 `json:"silver_content"`                       // mg
	PlatinumContent    float64 `json:"platinum_content"`                     // mg
	RareEarthElements  string  `gorm:"type:json" json:"rare_earth_elements"` // JSON object

	// Plastic Content Breakdown
	TotalPlasticContent float64 `json:"total_plastic_content"`          // percentage
	PlasticTypes        string  `gorm:"type:json" json:"plastic_types"` // JSON object
	RecyclablePlastic   float64 `json:"recyclable_plastic"`             // percentage
	BioplasticContent   float64 `json:"bioplastic_content"`             // percentage

	// Biodegradable Components
	BiodegradableContent float64 `json:"biodegradable_content"`                // percentage
	BiodegradableParts   string  `gorm:"type:json" json:"biodegradable_parts"` // JSON array
	CompostableElements  float64 `json:"compostable_elements"`                 // percentage
	DegradationTime      int     `json:"degradation_time"`                     // months

	// Recycling Facility Compatibility
	FacilityCompatible bool    `json:"facility_compatible"`
	SpecialHandling    bool    `json:"special_handling"`
	ProcessingCost     float64 `json:"processing_cost"`
	NearestFacility    string  `json:"nearest_facility"`
	TransportDistance  float64 `json:"transport_distance"` // km

	// Take-back Program
	TakeBackEligible    bool    `json:"take_back_eligible"`
	ManufacturerProgram bool    `json:"manufacturer_program"`
	RetailerProgram     bool    `json:"retailer_program"`
	ProgramIncentive    float64 `json:"program_incentive"`
	PrepaidShipping     bool    `json:"prepaid_shipping"`

	// Refurbishment Potential
	RefurbishmentScore    float64 `json:"refurbishment_score"` // 0-100
	RepairableComponents  int     `json:"repairable_components"`
	UpgradeableComponents int     `json:"upgradeable_components"`
	EstimatedRefurbCost   float64 `json:"estimated_refurb_cost"`
	RefurbValueRetention  float64 `json:"refurb_value_retention"` // percentage

	// Component Reusability
	ReusabilityScore   float64 `json:"reusability_score"`                    // 0-100
	ReusableComponents string  `gorm:"type:json" json:"reusable_components"` // JSON array
	ComponentValue     float64 `json:"component_value"`
	DisassemblyTime    int     `json:"disassembly_time"` // minutes

	// Circular Economy Contribution
	CircularityIndex     float64 `json:"circularity_index"`      // 0-100
	MaterialRecoveryRate float64 `json:"material_recovery_rate"` // percentage
	WasteReduction       float64 `json:"waste_reduction"`        // kg
	ResourceSavings      float64 `json:"resource_savings"`       // percentage
	EconomicValue        float64 `json:"economic_value"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceSustainabilityMetrics tracks environmental performance metrics
type DeviceSustainabilityMetrics struct {
	database.BaseModel
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	MeasurementDate time.Time `json:"measurement_date"`

	// Energy Consumption
	EnergyConsumption      float64 `json:"energy_consumption"`       // kWh per year
	IdlePowerDraw          float64 `json:"idle_power_draw"`          // watts
	ActivePowerDraw        float64 `json:"active_power_draw"`        // watts
	PeakPowerDraw          float64 `json:"peak_power_draw"`          // watts
	EnergyEfficiencyRating string  `json:"energy_efficiency_rating"` // A++ to G

	// Battery Efficiency
	BatteryEfficiency  float64 `json:"battery_efficiency"` // percentage
	ChargeCycles       int     `json:"charge_cycles"`
	BatteryDegradation float64 `json:"battery_degradation"` // percentage
	OptimalChargeRate  float64 `json:"optimal_charge_rate"` // watts
	EnergyRetention    float64 `json:"energy_retention"`    // percentage

	// Standby Power Consumption
	StandbyPower   float64 `json:"standby_power"`    // watts
	SleepModePower float64 `json:"sleep_mode_power"` // watts
	OffModePower   float64 `json:"off_mode_power"`   // watts
	WakeUpEnergy   float64 `json:"wake_up_energy"`   // joules

	// Charging Efficiency
	ChargingEfficiency   float64 `json:"charging_efficiency"`    // percentage
	WirelessChargingLoss float64 `json:"wireless_charging_loss"` // percentage
	FastChargingImpact   float64 `json:"fast_charging_impact"`   // efficiency loss
	OptimalChargingTemp  float64 `json:"optimal_charging_temp"`  // celsius

	// Heat Generation
	HeatGeneration     float64 `json:"heat_generation"`     // watts
	ThermalEfficiency  float64 `json:"thermal_efficiency"`  // percentage
	CoolingRequirement float64 `json:"cooling_requirement"` // watts
	HeatRecovery       float64 `json:"heat_recovery"`       // percentage

	// Resource Depletion
	ResourceDepletionScore float64 `json:"resource_depletion_score"`            // 0-100
	CriticalMaterials      string  `gorm:"type:json" json:"critical_materials"` // JSON array
	ScarceMinerals         string  `gorm:"type:json" json:"scarce_minerals"`    // JSON array
	ResourceEfficiency     float64 `json:"resource_efficiency"`                 // percentage

	// Water Usage
	ManufacturingWater  float64 `json:"manufacturing_water"`   // liters
	WaterFootprint      float64 `json:"water_footprint"`       // liters total
	WaterRecyclingRate  float64 `json:"water_recycling_rate"`  // percentage
	WaterPollutionIndex float64 `json:"water_pollution_index"` // 0-100

	// Renewable Materials
	RenewableMaterials  float64 `json:"renewable_materials"`  // percentage
	RecycledContent     float64 `json:"recycled_content"`     // percentage
	BiomassContent      float64 `json:"biomass_content"`      // percentage
	SustainableSourcing float64 `json:"sustainable_sourcing"` // percentage

	// Sustainability Certifications
	CertificationStatus   string `gorm:"type:json" json:"certification_status"`   // JSON array
	CertificationDates    string `gorm:"type:json" json:"certification_dates"`    // JSON object
	CertificationScores   string `gorm:"type:json" json:"certification_scores"`   // JSON object
	PendingCertifications string `gorm:"type:json" json:"pending_certifications"` // JSON array

	// Environmental Compliance
	ComplianceStatus string     `json:"compliance_status"`                  // compliant, partial, non-compliant
	ComplianceIssues string     `gorm:"type:json" json:"compliance_issues"` // JSON array
	ComplianceScore  float64    `json:"compliance_score"`                   // 0-100
	LastAuditDate    *time.Time `json:"last_audit_date"`

	// Overall Metrics
	SustainabilityScore  float64 `json:"sustainability_score"`  // 0-100
	EnvironmentalGrade   string  `json:"environmental_grade"`   // A-F
	ImprovementPotential float64 `json:"improvement_potential"` // percentage
	IndustryBenchmark    float64 `json:"industry_benchmark"`    // 0-100

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceEcoLabel tracks environmental certifications and labels
type DeviceEcoLabel struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Energy Star Rating
	EnergyStarCertified bool    `json:"energy_star_certified"`
	EnergyStarRating    float64 `json:"energy_star_rating"`
	EnergyStarYear      int     `json:"energy_star_year"`
	EnergyStarCategory  string  `json:"energy_star_category"`

	// EPEAT Certification
	EPEATCertified bool       `json:"epeat_certified"`
	EPEATLevel     string     `json:"epeat_level"` // bronze, silver, gold
	EPEATRegistry  string     `json:"epeat_registry"`
	EPEATExpiry    *time.Time `json:"epeat_expiry"`

	// EU Energy Label
	EUEnergyLabel string  `json:"eu_energy_label"` // A+++ to G
	EUEnergyScore float64 `json:"eu_energy_score"`
	EUEnergyClass string  `json:"eu_energy_class"`
	EULabelYear   int     `json:"eu_label_year"`

	// Green Electronics Certification
	GreenElectronics bool       `json:"green_electronics"`
	GreenCertifier   string     `json:"green_certifier"`
	GreenScore       float64    `json:"green_score"` // 0-100
	GreenExpiry      *time.Time `json:"green_expiry"`

	// RoHS Compliance
	RoHSCompliant   bool   `json:"rohs_compliant"`
	RoHSVersion     string `json:"rohs_version"`                     // RoHS 2, RoHS 3
	RoHSExemptions  string `gorm:"type:json" json:"rohs_exemptions"` // JSON array
	RoHSCertificate string `json:"rohs_certificate"`

	// REACH Compliance
	REACHCompliant    bool   `json:"reach_compliant"`
	SVHCContent       string `gorm:"type:json" json:"svhc_content"` // JSON array - Substances of Very High Concern
	REACHRegistration string `json:"reach_registration"`
	REACHDeclaration  string `json:"reach_declaration"`

	// Conflict Minerals
	ConflictFree     bool   `json:"conflict_free"`
	ConflictMinerals string `gorm:"type:json" json:"conflict_minerals"` // JSON object (tin, tungsten, tantalum, gold)
	SupplyChainAudit bool   `json:"supply_chain_audit"`
	SMELTERCompliant bool   `json:"smelter_compliant"` // Conflict-Free Smelter Program

	// Fair Trade Certification
	FairTradeCertified  bool    `json:"fair_trade_certified"`
	FairTradePercentage float64 `json:"fair_trade_percentage"`
	FairTradeComponents string  `gorm:"type:json" json:"fair_trade_components"` // JSON array
	WorkerWelfare       float64 `json:"worker_welfare"`                         // score 0-100

	// Sustainable Packaging
	PackagingScore       float64 `json:"packaging_score"`      // 0-100
	RecyclablePackaging  float64 `json:"recyclable_packaging"` // percentage
	PlasticFreePackaging bool    `json:"plastic_free_packaging"`
	MinimalPackaging     bool    `json:"minimal_packaging"`

	// Eco-design Compliance
	EcoDesignCompliant bool    `json:"eco_design_compliant"`
	DesignForRecycling float64 `json:"design_for_recycling"` // score 0-100
	DesignForRepair    float64 `json:"design_for_repair"`    // score 0-100
	ModularDesign      bool    `json:"modular_design"`

	// Overall Eco Score
	OverallEcoScore    float64    `json:"overall_eco_score"` // 0-100
	EcoRanking         int        `json:"eco_ranking"`       // category ranking
	CertificationCount int        `json:"certification_count"`
	NextReviewDate     *time.Time `json:"next_review_date"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceLifecycleAssessment evaluates environmental impact throughout device lifecycle
type DeviceLifecycleAssessment struct {
	database.BaseModel
	DeviceID       uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AssessmentDate time.Time `json:"assessment_date"`

	// Raw Material Extraction
	RawMaterialImpact float64 `json:"raw_material_impact"` // environmental points
	MiningEmissions   float64 `json:"mining_emissions"`    // kg CO2e
	HabitatDisruption float64 `json:"habitat_disruption"`  // hectares
	WaterConsumption  float64 `json:"water_consumption"`   // liters
	EnergyExtraction  float64 `json:"energy_extraction"`   // kWh

	// Manufacturing Process Impact
	ManufacturingImpact float64 `json:"manufacturing_impact"`            // environmental points
	ProcessEmissions    float64 `json:"process_emissions"`               // kg CO2e
	ChemicalUsage       string  `gorm:"type:json" json:"chemical_usage"` // JSON array
	WasteGeneration     float64 `json:"waste_generation"`                // kg
	EnergyConsumption   float64 `json:"energy_consumption"`              // kWh

	// Distribution Environmental Cost
	DistributionImpact   float64 `json:"distribution_impact"` // environmental points
	TransportEmissions   float64 `json:"transport_emissions"` // kg CO2e
	PackagingWaste       float64 `json:"packaging_waste"`     // kg
	ColdChainRequirement bool    `json:"cold_chain_requirement"`
	LastMileImpact       float64 `json:"last_mile_impact"` // environmental points

	// Use Phase Environmental Impact
	UsePhaseImpact       float64 `json:"use_phase_impact"`                      // environmental points
	OperationalEmissions float64 `json:"operational_emissions"`                 // kg CO2e per year
	ConsumablesRequired  string  `gorm:"type:json" json:"consumables_required"` // JSON array
	MaintenanceImpact    float64 `json:"maintenance_impact"`                    // environmental points
	ExpectedLifespan     int     `json:"expected_lifespan"`                     // months

	// End-of-Life Treatment Options
	RecyclingOption     float64 `json:"recycling_option"`     // environmental benefit
	RefurbishmentOption float64 `json:"refurbishment_option"` // environmental benefit
	LandfillImpact      float64 `json:"landfill_impact"`      // environmental points
	IncinerationImpact  float64 `json:"incineration_impact"`  // environmental points
	RecoveryPotential   float64 `json:"recovery_potential"`   // percentage

	// Total Environmental Footprint
	TotalFootprint      float64 `json:"total_footprint"`      // environmental points
	CarbonFootprint     float64 `json:"carbon_footprint"`     // kg CO2e total
	WaterFootprint      float64 `json:"water_footprint"`      // liters total
	LandFootprint       float64 `json:"land_footprint"`       // hectares
	EcologicalFootprint float64 `json:"ecological_footprint"` // global hectares

	// Lifecycle Cost Analysis
	EnvironmentalCost  float64 `json:"environmental_cost"`
	SocialCost         float64 `json:"social_cost"`
	ExternalitiesCost  float64 `json:"externalities_cost"`
	TotalLifecycleCost float64 `json:"total_lifecycle_cost"`

	// Environmental Payback Period
	PaybackPeriod           int        `json:"payback_period"` // months
	BreakevenPoint          *time.Time `json:"breakeven_point"`
	NetEnvironmentalBenefit float64    `json:"net_environmental_benefit"`
	ROIEnvironmental        float64    `json:"roi_environmental"` // percentage

	// Sustainability Improvement Tracking
	ImprovementAreas   string  `gorm:"type:json" json:"improvement_areas"`   // JSON array
	ReductionPotential float64 `json:"reduction_potential"`                  // percentage
	ImplementedChanges string  `gorm:"type:json" json:"implemented_changes"` // JSON array
	ImpactReduction    float64 `json:"impact_reduction"`                     // percentage achieved

	// Benchmark Comparisons
	IndustryAverage       float64 `json:"industry_average"` // environmental points
	BestInClass           float64 `json:"best_in_class"`    // environmental points
	PerformancePercentile float64 `json:"performance_percentile"`
	CompetitorComparison  string  `gorm:"type:json" json:"competitor_comparison"` // JSON object

	// Overall Assessment
	LCAScore             float64 `json:"lca_score"`                        // 0-100
	EnvironmentalGrade   string  `json:"environmental_grade"`              // A-F
	SustainabilityRating float64 `json:"sustainability_rating"`            // 0-5 stars
	Recommendations      string  `gorm:"type:json" json:"recommendations"` // JSON array

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceRepairability tracks device repair capabilities and right-to-repair compliance
type DeviceRepairability struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Repairability Index Score
	RepairabilityScore  float64 `json:"repairability_score"`   // 0-10
	EaseOfDisassembly   float64 `json:"ease_of_disassembly"`   // 0-10
	AvailabilityOfParts float64 `json:"availability_of_parts"` // 0-10
	PriceOfParts        float64 `json:"price_of_parts"`        // 0-10 (affordability)
	DocumentationScore  float64 `json:"documentation_score"`   // 0-10

	// Spare Parts Availability
	PartsAvailability int    `json:"parts_availability"` // years guaranteed
	OEMPartsAvailable bool   `json:"oem_parts_available"`
	ThirdPartyParts   bool   `json:"third_party_parts"`
	PartsInventory    string `gorm:"type:json" json:"parts_inventory"`  // JSON array
	PartsPriceList    string `gorm:"type:json" json:"parts_price_list"` // JSON object

	// Repair Documentation
	ServiceManualAvailable bool `json:"service_manual_available"`
	RepairGuidesAvailable  bool `json:"repair_guides_available"`
	VideoTutorials         bool `json:"video_tutorials"`
	SchematicsAvailable    bool `json:"schematics_available"`
	DiagnosticTools        bool `json:"diagnostic_tools"`

	// Tool Requirements
	SpecialToolsNeeded bool    `json:"special_tools_needed"`
	ToolsList          string  `gorm:"type:json" json:"tools_list"` // JSON array
	StandardToolsOnly  bool    `json:"standard_tools_only"`
	ProprietaryTools   int     `json:"proprietary_tools"`
	ToolCost           float64 `json:"tool_cost"`

	// Modular Design
	ModularDesignScore     float64 `json:"modular_design_score"` // 0-10
	ReplaceableComponents  int     `json:"replaceable_components"`
	UpgradeableComponents  int     `json:"upgradeable_components"`
	ComponentAccessibility float64 `json:"component_accessibility"` // 0-10

	// DIY Repair Feasibility
	DIYFeasibility       float64 `json:"diy_feasibility"`                        // 0-10
	SkillLevelRequired   string  `json:"skill_level_required"`                   // beginner, intermediate, expert
	SafetyConsiderations string  `gorm:"type:json" json:"safety_considerations"` // JSON array
	WarrantyImpact       string  `json:"warranty_impact"`                        // void, partial, maintained

	// Authorized Repair Network
	AuthorizedCenters   int     `json:"authorized_centers"`
	IndependentRepair   bool    `json:"independent_repair"`
	RepairNetworkSize   int     `json:"repair_network_size"`
	AverageDistance     float64 `json:"average_distance"` // km to nearest center
	OnlineRepairService bool    `json:"online_repair_service"`

	// Average Repair Time
	AverageRepairTime   int    `json:"average_repair_time"`                    // hours
	CommonRepairs       string `gorm:"type:json" json:"common_repairs"`        // JSON array
	RepairTimeBreakdown string `gorm:"type:json" json:"repair_time_breakdown"` // JSON object
	TurnaroundTime      int    `json:"turnaround_time"`                        // days

	// Repair Cost vs Replacement
	AverageRepairCost float64 `json:"average_repair_cost"`
	ReplacementCost   float64 `json:"replacement_cost"`
	RepairCostRatio   float64 `json:"repair_cost_ratio"`  // repair/replacement
	CostEffectiveness float64 `json:"cost_effectiveness"` // 0-100

	// Right-to-Repair Compliance
	RightToRepairScore    float64 `json:"right_to_repair_score"` // 0-100
	LegalCompliance       bool    `json:"legal_compliance"`
	ConsumerRightsSupport bool    `json:"consumer_rights_support"`
	OpenSourceHardware    bool    `json:"open_source_hardware"`
	RepairRestrictions    string  `gorm:"type:json" json:"repair_restrictions"` // JSON array

	// Overall Repairability
	OverallScore           float64 `json:"overall_score"`       // 0-100
	RepairabilityGrade     string  `json:"repairability_grade"` // A-F
	IndustryRanking        int     `json:"industry_ranking"`
	ImprovementSuggestions string  `gorm:"type:json" json:"improvement_suggestions"` // JSON array

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceCarbonFootprint
func (dcf *DeviceCarbonFootprint) IsLowCarbon() bool {
	return dcf.TotalCarbon < dcf.IndustryAverage && dcf.EnvironmentalScore >= 70
}

func (dcf *DeviceCarbonFootprint) IsCarbonNeutral() bool {
	netCarbon := dcf.TotalCarbon - dcf.CarbonOffsetsApplied - dcf.RefurbishmentSavings
	return netCarbon <= 0 || dcf.CurrentNeutralityRate >= 100
}

func (dcf *DeviceCarbonFootprint) GetCarbonGrade() string {
	if dcf.PerformanceRating != "" {
		return dcf.PerformanceRating
	}
	if dcf.EnvironmentalScore >= 90 {
		return "A"
	} else if dcf.EnvironmentalScore >= 80 {
		return "B"
	} else if dcf.EnvironmentalScore >= 70 {
		return "C"
	} else if dcf.EnvironmentalScore >= 60 {
		return "D"
	}
	return "F"
}

func (dcf *DeviceCarbonFootprint) OnTrackForTarget() bool {
	return dcf.CurrentReduction >= dcf.ReductionTarget*0.8 // 80% of target achieved
}

func (dcf *DeviceCarbonFootprint) GetLifecycleEmissions() float64 {
	return dcf.ManufacturingEmissions + dcf.TransportationEmissions +
		dcf.UsageEmissions + dcf.EndOfLifeEmissions
}

// Methods for DeviceRecyclingScore
func (drs *DeviceRecyclingScore) IsHighlyRecyclable() bool {
	return drs.RecyclabilityPercentage >= 80 && drs.CircularityIndex >= 70
}

func (drs *DeviceRecyclingScore) HasHazardousMaterials() bool {
	return drs.HazardousContent > 0 || drs.LeadContent > 0 ||
		drs.MercuryContent > 0 || drs.CadmiumContent > 0
}

func (drs *DeviceRecyclingScore) IsCircularEconomyReady() bool {
	return drs.CircularityIndex >= 70 && drs.RefurbishmentScore >= 70 &&
		drs.MaterialRecoveryRate >= 80
}

func (drs *DeviceRecyclingScore) HasTakeBackProgram() bool {
	return drs.TakeBackEligible && (drs.ManufacturerProgram || drs.RetailerProgram)
}

func (drs *DeviceRecyclingScore) GetRecyclingGrade() string {
	score := (drs.RecyclabilityPercentage + drs.CircularityIndex + drs.ReusabilityScore) / 3
	if score >= 90 {
		return "A+"
	} else if score >= 80 {
		return "A"
	} else if score >= 70 {
		return "B"
	} else if score >= 60 {
		return "C"
	} else if score >= 50 {
		return "D"
	}
	return "F"
}

// Methods for DeviceSustainabilityMetrics
func (dsm *DeviceSustainabilityMetrics) IsEnergyEfficient() bool {
	return dsm.ChargingEfficiency >= 85 && dsm.BatteryEfficiency >= 90 &&
		(dsm.EnergyEfficiencyRating == "A++" || dsm.EnergyEfficiencyRating == "A+")
}

func (dsm *DeviceSustainabilityMetrics) IsSustainable() bool {
	return dsm.SustainabilityScore >= 70 && dsm.ComplianceStatus == "compliant" &&
		dsm.RenewableMaterials >= 30
}

func (dsm *DeviceSustainabilityMetrics) HasLowStandbyPower() bool {
	return dsm.StandbyPower < 0.5 && dsm.OffModePower < 0.1 // watts
}

func (dsm *DeviceSustainabilityMetrics) GetSustainabilityLevel() string {
	if dsm.EnvironmentalGrade != "" {
		return dsm.EnvironmentalGrade
	}
	if dsm.SustainabilityScore >= 80 {
		return "Excellent"
	} else if dsm.SustainabilityScore >= 60 {
		return "Good"
	} else if dsm.SustainabilityScore >= 40 {
		return "Fair"
	}
	return "Poor"
}

func (dsm *DeviceSustainabilityMetrics) NeedsImprovement() bool {
	return dsm.ImprovementPotential > 30 || dsm.SustainabilityScore < dsm.IndustryBenchmark
}

// Methods for DeviceEcoLabel
func (del *DeviceEcoLabel) HasGreenCertifications() bool {
	return del.EnergyStarCertified || del.EPEATCertified || del.GreenElectronics
}

func (del *DeviceEcoLabel) IsFullyCompliant() bool {
	return del.RoHSCompliant && del.REACHCompliant && del.ConflictFree &&
		del.EcoDesignCompliant
}

func (del *DeviceEcoLabel) GetEcoGrade() string {
	if del.OverallEcoScore >= 90 {
		return "Platinum"
	} else if del.OverallEcoScore >= 80 {
		return "Gold"
	} else if del.OverallEcoScore >= 70 {
		return "Silver"
	} else if del.OverallEcoScore >= 60 {
		return "Bronze"
	}
	return "None"
}

func (del *DeviceEcoLabel) HasEthicalSourcing() bool {
	return del.ConflictFree && del.FairTradeCertified && del.WorkerWelfare >= 70
}

func (del *DeviceEcoLabel) GetCertificationStrength() string {
	if del.CertificationCount >= 8 {
		return "Very Strong"
	} else if del.CertificationCount >= 5 {
		return "Strong"
	} else if del.CertificationCount >= 3 {
		return "Moderate"
	}
	return "Weak"
}

// Methods for DeviceLifecycleAssessment
func (dla *DeviceLifecycleAssessment) HasLowImpact() bool {
	return dla.TotalFootprint < dla.IndustryAverage && dla.LCAScore >= 70
}

func (dla *DeviceLifecycleAssessment) GetLCAGrade() string {
	if dla.EnvironmentalGrade != "" {
		return dla.EnvironmentalGrade
	}
	if dla.LCAScore >= 85 {
		return "A"
	} else if dla.LCAScore >= 70 {
		return "B"
	} else if dla.LCAScore >= 55 {
		return "C"
	} else if dla.LCAScore >= 40 {
		return "D"
	}
	return "F"
}

func (dla *DeviceLifecycleAssessment) IsEnvironmentallyPositive() bool {
	return dla.NetEnvironmentalBenefit > 0 && dla.ROIEnvironmental > 0
}

func (dla *DeviceLifecycleAssessment) GetDominantPhase() string {
	maxImpact := dla.RawMaterialImpact
	phase := "raw_material"

	if dla.ManufacturingImpact > maxImpact {
		maxImpact = dla.ManufacturingImpact
		phase = "manufacturing"
	}
	if dla.DistributionImpact > maxImpact {
		maxImpact = dla.DistributionImpact
		phase = "distribution"
	}
	if dla.UsePhaseImpact > maxImpact {
		maxImpact = dla.UsePhaseImpact
		phase = "use"
	}
	totalEndOfLife := dla.LandfillImpact + dla.IncinerationImpact - dla.RecyclingOption - dla.RefurbishmentOption
	if totalEndOfLife > maxImpact {
		phase = "end_of_life"
	}

	return phase
}

func (dla *DeviceLifecycleAssessment) PerformsBetterThanAverage() bool {
	return dla.PerformancePercentile >= 50 && dla.TotalFootprint < dla.IndustryAverage
}

// Methods for DeviceRepairability
func (dr *DeviceRepairability) IsHighlyRepairable() bool {
	return dr.RepairabilityScore >= 7.0 && dr.OverallScore >= 70
}

func (dr *DeviceRepairability) HasDIYSupport() bool {
	return dr.DIYFeasibility >= 6.0 && dr.RepairGuidesAvailable &&
		(dr.StandardToolsOnly || !dr.SpecialToolsNeeded)
}

func (dr *DeviceRepairability) IsRightToRepairCompliant() bool {
	return dr.RightToRepairScore >= 70 && dr.LegalCompliance &&
		dr.ConsumerRightsSupport
}

func (dr *DeviceRepairability) IsCostEffectiveToRepair() bool {
	return dr.RepairCostRatio < 0.5 && dr.CostEffectiveness >= 70 // Repair costs less than 50% of replacement
}

func (dr *DeviceRepairability) GetRepairabilityClass() string {
	if dr.RepairabilityGrade != "" {
		return dr.RepairabilityGrade
	}
	if dr.RepairabilityScore >= 8 {
		return "Excellent"
	} else if dr.RepairabilityScore >= 6 {
		return "Good"
	} else if dr.RepairabilityScore >= 4 {
		return "Fair"
	} else if dr.RepairabilityScore >= 2 {
		return "Poor"
	}
	return "Very Poor"
}
