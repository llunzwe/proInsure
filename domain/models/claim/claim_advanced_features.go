package claim

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"smartsure/pkg/database"
)

// ClaimBiometricVerification - Biometric verification for claim authentication
type ClaimBiometricVerification struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Verification Methods
	FaceIDVerified       bool `json:"face_id_verified"`
	FingerprintVerified  bool `json:"fingerprint_verified"`
	VoiceVerified        bool `json:"voice_verified"`
	BehavioralBiometrics bool `json:"behavioral_biometrics"`

	// Verification Details
	VerificationTimestamp     time.Time `json:"verification_timestamp"`
	VerificationScore         float64   `json:"verification_score"`
	FailedAttempts            int       `json:"failed_attempts"`
	DeviceUsedForVerification string    `json:"device_used_for_verification"`

	// Liveness Detection
	LivenessCheckPassed   bool    `json:"liveness_check_passed"`
	PhotoMatchScore       float64 `json:"photo_match_score"`
	VideoVerificationDone bool    `json:"video_verification_done"`

	// Security
	SuspiciousActivity    bool   `json:"suspicious_activity"`
	VerificationBypass    bool   `json:"verification_bypass"`
	BypassReason          string `json:"bypass_reason"`
	SecondaryVerification bool   `json:"secondary_verification"`
}

// Claim5GAndIoT - 5G and IoT device specific features
type Claim5GAndIoT struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// 5G Features
	Is5GDevice            bool `json:"is_5g_device"`
	Network5GAvailable    bool `json:"network_5g_available"`
	MillimeterWaveSupport bool `json:"millimeter_wave_support"`
	NetworkSlicingActive  bool `json:"network_slicing_active"`

	// Connected Devices
	ConnectedWearables    datatypes.JSON `gorm:"type:json" json:"connected_wearables"`
	SmartHomeDevices      datatypes.JSON `gorm:"type:json" json:"smart_home_devices"`
	AutomotiveIntegration bool           `json:"automotive_integration"`
	ConnectedDevicesCount int            `json:"connected_devices_count"`

	// IoT Ecosystem Impact
	EcosystemDisruption      bool           `json:"ecosystem_disruption"`
	SecondaryDevicesAffected datatypes.JSON `gorm:"type:json" json:"secondary_devices_affected"`
	DataSyncLoss             bool           `json:"data_sync_loss"`
	AutomationRulesLost      datatypes.JSON `gorm:"type:json" json:"automation_rules_lost"`

	// Edge Computing
	EdgeComputingEnabled bool `json:"edge_computing_enabled"`
	EdgeDataLoss         bool `json:"edge_data_loss"`
	LocalMLModelsLost    bool `json:"local_ml_models_lost"`
}

// ClaimAugmentedReality - AR/VR related claim features
type ClaimAugmentedReality struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// AR Capabilities
	ARKitSupported        bool `json:"arkit_supported"`
	ARCoreSupported       bool `json:"arcore_supported"`
	LiDARScannerPresent   bool `json:"lidar_scanner_present"`
	DepthCameraFunctional bool `json:"depth_camera_functional"`

	// AR Content Loss
	ARAppsAffected       datatypes.JSON `gorm:"type:json" json:"ar_apps_affected"`
	AR3DModelsLost       int            `json:"ar_3d_models_lost"`
	ARMeasurementsLost   bool           `json:"ar_measurements_lost"`
	ARNavigationDataLost bool           `json:"ar_navigation_data_lost"`

	// Professional Use
	ProfessionalARUse   bool    `json:"professional_ar_use"`
	BusinessImpact      string  `gorm:"type:text" json:"business_impact"`
	RevenueLossEstimate float64 `json:"revenue_loss_estimate"`

	// VR Integration
	VRHeadsetPaired       bool           `json:"vr_headset_paired"`
	VRContentLost         bool           `json:"vr_content_lost"`
	VRAccessoriesAffected datatypes.JSON `gorm:"type:json" json:"vr_accessories_affected"`
}

// ClaimCryptocurrency - Crypto wallet and digital asset claims
type ClaimCryptocurrency struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Wallet Information
	CryptoWalletPresent   bool           `json:"crypto_wallet_present"`
	WalletTypes           datatypes.JSON `gorm:"type:json" json:"wallet_types"` // hot, cold, hardware
	WalletBackupAvailable bool           `json:"wallet_backup_available"`
	SeedPhraseSecured     bool           `json:"seed_phrase_secured"`

	// Asset Details
	CryptoAssetsValue float64        `json:"crypto_assets_value_usd"`
	NFTsPresent       bool           `json:"nfts_present"`
	NFTsValue         float64        `json:"nfts_value_usd"`
	DeFiPositions     datatypes.JSON `gorm:"type:json" json:"defi_positions"`

	// Recovery
	RecoveryAttempted  bool    `json:"recovery_attempted"`
	RecoverySuccessful *bool   `json:"recovery_successful"`
	AssetsRecovered    float64 `json:"assets_recovered_usd"`

	// Security
	TwoFactorEnabled     bool `json:"two_factor_enabled"`
	HardwareWalletLinked bool `json:"hardware_wallet_linked"`
	MultiSigRequired     bool `json:"multisig_required"`

	// Coverage
	CryptoCovered  bool           `json:"crypto_covered"`
	CoverageLimit  float64        `json:"coverage_limit_usd"`
	ExcludedAssets datatypes.JSON `gorm:"type:json" json:"excluded_assets"`
}

// ClaimSubscriptionServices - Subscription and recurring services
type ClaimSubscriptionServices struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Active Subscriptions
	ActiveSubscriptions      datatypes.JSON `gorm:"type:json" json:"active_subscriptions"`
	TotalMonthlyValue        float64        `json:"total_monthly_value"`
	AnnualSubscriptionsValue float64        `json:"annual_subscriptions_value"`

	// Service Categories
	StreamingServices    datatypes.JSON `gorm:"type:json" json:"streaming_services"`
	CloudStorageServices datatypes.JSON `gorm:"type:json" json:"cloud_storage_services"`
	ProductivityApps     datatypes.JSON `gorm:"type:json" json:"productivity_apps"`
	GamingSubscriptions  datatypes.JSON `gorm:"type:json" json:"gaming_subscriptions"`
	FitnessApps          datatypes.JSON `gorm:"type:json" json:"fitness_apps"`

	// Transfer Support
	TransferSupported bool           `json:"transfer_supported"`
	TransferCompleted bool           `json:"transfer_completed"`
	FailedTransfers   datatypes.JSON `gorm:"type:json" json:"failed_transfers"`

	// Compensation
	ProRataRefundDue      float64 `json:"pro_rata_refund_due"`
	ServiceCreditsOffered float64 `json:"service_credits_offered"`
	FreeMonthsOffered     int     `json:"free_months_offered"`
}

// ClaimEnvironmentalImpact - Environmental and sustainability tracking
type ClaimEnvironmentalImpact struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Device Lifecycle
	DeviceAge             int  `json:"device_age_months"`
	PreviousRepairs       int  `json:"previous_repairs"`
	RecyclingEligible     bool `json:"recycling_eligible"`
	RefurbishmentPossible bool `json:"refurbishment_possible"`

	// Environmental Scores
	CarbonFootprintSaved float64 `json:"carbon_footprint_saved_kg"`
	EWasteReduced        float64 `json:"e_waste_reduced_kg"`
	RepairVsReplaceScore float64 `json:"repair_vs_replace_score"`
	SustainabilityRating string  `json:"sustainability_rating"`

	// Circular Economy
	TradeInOffered           bool           `json:"trade_in_offered"`
	RecyclingProgramEnrolled bool           `json:"recycling_program_enrolled"`
	CertifiedRecycler        bool           `json:"certified_recycler"`
	MaterialsRecovered       datatypes.JSON `gorm:"type:json" json:"materials_recovered"`

	// Green Options
	GreenRepairChosen      bool `json:"green_repair_chosen"`
	RefurbishedReplacement bool `json:"refurbished_replacement"`
	CarbonOffsetPurchased  bool `json:"carbon_offset_purchased"`
	TreesPlantedCount      int  `json:"trees_planted_count"`
}

// ClaimFoldableFlexible - Foldable and flexible display claims
type ClaimFoldableFlexible struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Device Type
	IsFoldableDevice  bool   `json:"is_foldable_device"`
	IsFlexibleDisplay bool   `json:"is_flexible_display"`
	IsDualScreen      bool   `json:"is_dual_screen"`
	FormFactor        string `json:"form_factor"`

	// Damage Specifics
	HingeDamage       bool `json:"hinge_damage"`
	FoldCreaseDamage  bool `json:"fold_crease_damage"`
	FlexPointFailure  bool `json:"flex_point_failure"`
	ScreenDelamintion bool `json:"screen_delamination"`

	// Usage Metrics
	FoldCycles   int     `json:"fold_cycles"`
	DailyFolds   int     `json:"daily_folds_average"`
	MaxFoldAngle float64 `json:"max_fold_angle"`

	// Repair Complexity
	SpecializedRepairNeeded bool    `json:"specialized_repair_needed"`
	OEMRepairRequired       bool    `json:"oem_repair_required"`
	RepairComplexityScore   float64 `json:"repair_complexity_score"`
	EstimatedRepairCost     float64 `json:"estimated_repair_cost"`
}

// ClaimHealthAndWellness - Health and medical app data
type ClaimHealthAndWellness struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Health Data
	HealthAppDataPresent bool    `json:"health_app_data_present"`
	FitnessDataYears     float64 `json:"fitness_data_years"`
	MedicalRecordsStored bool    `json:"medical_records_stored"`
	ECGDataPresent       bool    `json:"ecg_data_present"`
	BloodOxygenData      bool    `json:"blood_oxygen_data_present"`

	// Medical Device Integration
	MedicalDevicesPaired   datatypes.JSON `gorm:"type:json" json:"medical_devices_paired"`
	CriticalHealthAlerts   bool           `json:"critical_health_alerts"`
	EmergencyContactsSetup bool           `json:"emergency_contacts_setup"`
	MedicalIDConfigured    bool           `json:"medical_id_configured"`

	// Data Recovery Priority
	HealthDataPriority    string `json:"health_data_priority"` // critical, high, normal
	DataRecoveryUrgent    bool   `json:"data_recovery_urgent"`
	PhysicianNotification bool   `json:"physician_notification"`

	// Wellness Impact
	WorkoutStreakLost   int  `json:"workout_streak_lost_days"`
	MindfulnessDataLost bool `json:"mindfulness_data_lost"`
	SleepDataLost       bool `json:"sleep_data_lost"`
	NutritionDataLost   bool `json:"nutrition_data_lost"`
}

// ClaimBusinessContinuity - Business impact for professional users
type ClaimBusinessContinuity struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Business Use
	BusinessDevice       bool           `json:"business_device"`
	CompanyOwned         bool           `json:"company_owned"`
	MDMEnrolled          bool           `json:"mdm_enrolled"`
	CompanyAppsInstalled datatypes.JSON `gorm:"type:json" json:"company_apps_installed"`

	// Impact Assessment
	BusinessImpactLevel string  `json:"business_impact_level"` // critical, high, medium, low
	DowntimeHours       float64 `json:"downtime_hours"`
	RevenueLoss         float64 `json:"revenue_loss"`
	ProductivityLoss    float64 `json:"productivity_loss_hours"`

	// Critical Data
	CorporateDataAtRisk  bool `json:"corporate_data_at_risk"`
	CustomerDataExposed  bool `json:"customer_data_exposed"`
	IntellectualProperty bool `json:"intellectual_property_affected"`
	ComplianceImpact     bool `json:"compliance_impact"`

	// Business Continuity
	BackupDeviceProvided bool   `json:"backup_device_provided"`
	RemoteWipeInitiated  bool   `json:"remote_wipe_initiated"`
	DataRecoveryPriority string `json:"data_recovery_priority"`
	ITSupportInvolved    bool   `json:"it_support_involved"`

	// SLA Requirements
	SLAApplicable     bool `json:"sla_applicable"`
	SLAResponseTime   int  `json:"sla_response_hours"`
	SLABreached       bool `json:"sla_breached"`
	PenaltyApplicable bool `json:"penalty_applicable"`
}

// Business Logic Methods

func (b *ClaimBiometricVerification) IsHighConfidence() bool {
	return b.VerificationScore > 0.95 &&
		b.LivenessCheckPassed &&
		b.FailedAttempts == 0
}

func (iot *Claim5GAndIoT) HasEcosystemImpact() bool {
	return iot.EcosystemDisruption ||
		iot.ConnectedDevicesCount > 5 ||
		iot.DataSyncLoss
}

func (ar *ClaimAugmentedReality) RequiresSpecializedSupport() bool {
	return ar.ProfessionalARUse ||
		ar.RevenueLossEstimate > 1000 ||
		ar.LiDARScannerPresent
}

func (c *ClaimCryptocurrency) IsHighValueClaim() bool {
	totalValue := c.CryptoAssetsValue + c.NFTsValue
	return totalValue > 10000 || c.HardwareWalletLinked
}

func (e *ClaimEnvironmentalImpact) GetSustainabilityScore() float64 {
	score := 50.0
	if e.RecyclingEligible {
		score += 10
	}
	if e.RefurbishmentPossible {
		score += 15
	}
	if e.GreenRepairChosen {
		score += 10
	}
	if e.RefurbishedReplacement {
		score += 10
	}
	if e.CarbonOffsetPurchased {
		score += 5
	}
	return score
}

func (f *ClaimFoldableFlexible) IsHighComplexityRepair() bool {
	return f.HingeDamage ||
		f.FoldCreaseDamage ||
		f.RepairComplexityScore > 0.8 ||
		f.OEMRepairRequired
}

func (h *ClaimHealthAndWellness) RequiresUrgentProcessing() bool {
	return h.HealthDataPriority == "critical" ||
		h.CriticalHealthAlerts ||
		h.DataRecoveryUrgent
}

func (b *ClaimBusinessContinuity) RequiresPriorityHandling() bool {
	return b.BusinessImpactLevel == "critical" ||
		b.RevenueLoss > 10000 ||
		b.ComplianceImpact ||
		b.SLAApplicable
}
