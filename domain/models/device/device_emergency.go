package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceEmergencyContacts manages ICE contacts and emergency information
type DeviceEmergencyContacts struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID   uuid.UUID `gorm:"type:uuid;index" json:"user_id"`

	// Primary Emergency Contact
	PrimaryContactName  string `json:"primary_contact_name"`
	PrimaryContactPhone string `json:"primary_contact_phone"`
	PrimaryContactEmail string `json:"primary_contact_email"`
	PrimaryRelationship string `json:"primary_relationship"` // spouse, parent, sibling, friend, doctor
	PrimaryVerified     bool   `json:"primary_verified"`

	// Secondary Contacts
	SecondaryContacts     string `gorm:"type:json" json:"secondary_contacts"` // JSON array
	EmergencyContactCount int    `json:"emergency_contact_count"`
	ContactPriority       string `gorm:"type:json" json:"contact_priority"` // JSON array ordered by priority

	// Medical Information
	BloodType             string `json:"blood_type"`
	Allergies             string `gorm:"type:json" json:"allergies"`          // JSON array
	MedicalConditions     string `gorm:"type:json" json:"medical_conditions"` // JSON array
	Medications           string `gorm:"type:json" json:"medications"`        // JSON array
	DoctorName            string `json:"doctor_name"`
	DoctorPhone           string `json:"doctor_phone"`
	HospitalPreference    string `json:"hospital_preference"`
	InsuranceProvider     string `json:"insurance_provider"`
	InsurancePolicyNumber string `json:"insurance_policy_number"`

	// Contact Verification
	VerificationStatus   string     `json:"verification_status"` // verified, pending, failed
	LastVerificationDate *time.Time `json:"last_verification_date"`
	VerificationMethod   string     `json:"verification_method"` // sms, email, call
	VerificationAttempts int        `json:"verification_attempts"`

	// Availability & Preferences
	ContactAvailability    string `gorm:"type:json" json:"contact_availability"` // JSON object with hours
	PreferredContactMethod string `json:"preferred_contact_method"`              // call, sms, email
	LanguagePreference     string `json:"language_preference"`
	TimeZone               string `json:"time_zone"`

	// International Emergency Numbers
	LocalEmergencyNumber string `json:"local_emergency_number"` // e.g., 911, 999, 112
	PoliceNumber         string `json:"police_number"`
	FireNumber           string `json:"fire_number"`
	MedicalNumber        string `json:"medical_number"`
	CountryCode          string `json:"country_code"`

	// Additional Information
	SpecialInstructions string     `gorm:"type:text" json:"special_instructions"`
	LastUpdated         time.Time  `json:"last_updated"`
	UpdateReminder      bool       `json:"update_reminder"`
	NextReviewDate      *time.Time `json:"next_review_date"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// DeviceBackupStatus tracks backup status and health
type DeviceBackupStatus struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Backup Status
	BackupEnabled       bool       `json:"backup_enabled"`
	LastBackupDate      *time.Time `json:"last_backup_date"`
	NextScheduledBackup *time.Time `json:"next_scheduled_backup"`
	BackupStatus        string     `json:"backup_status"` // success, failed, in_progress, pending

	// Backup Frequency
	BackupFrequency   string `json:"backup_frequency"` // hourly, daily, weekly, monthly
	AutoBackupEnabled bool   `json:"auto_backup_enabled"`
	BackupSchedule    string `gorm:"type:json" json:"backup_schedule"` // JSON object with schedule
	BackupWindow      string `json:"backup_window"`                    // time window for backups

	// Backup Size & Storage
	LastBackupSize     int64   `json:"last_backup_size"`  // bytes
	TotalBackupSize    int64   `json:"total_backup_size"` // bytes
	StorageUsed        int64   `json:"storage_used"`      // bytes
	StorageLimit       int64   `json:"storage_limit"`     // bytes
	CompressionEnabled bool    `json:"compression_enabled"`
	CompressionRatio   float64 `json:"compression_ratio"`

	// Cloud Backup
	CloudBackupEnabled bool       `json:"cloud_backup_enabled"`
	CloudProvider      string     `json:"cloud_provider"`     // iCloud, Google Drive, OneDrive, AWS
	CloudStorageUsed   int64      `json:"cloud_storage_used"` // bytes
	CloudSyncStatus    string     `json:"cloud_sync_status"`  // synced, syncing, error
	LastCloudSync      *time.Time `json:"last_cloud_sync"`

	// Local Backup
	LocalBackupEnabled bool   `json:"local_backup_enabled"`
	LocalBackupPath    string `json:"local_backup_path"`
	LocalStorageUsed   int64  `json:"local_storage_used"` // bytes
	ExternalStorage    bool   `json:"external_storage"`

	// Backup Integrity
	IntegrityCheckDate *time.Time `json:"integrity_check_date"`
	IntegrityStatus    string     `json:"integrity_status"` // verified, corrupted, unknown
	ChecksumVerified   bool       `json:"checksum_verified"`
	BackupVersion      string     `json:"backup_version"`

	// Restore Points
	RestorePointCount int        `json:"restore_point_count"`
	RestorePoints     string     `gorm:"type:json" json:"restore_points"` // JSON array
	LastRestoreDate   *time.Time `json:"last_restore_date"`
	RestoreSuccess    bool       `json:"restore_success"`

	// Encryption
	EncryptionEnabled bool   `json:"encryption_enabled"`
	EncryptionType    string `json:"encryption_type"` // AES256, RSA
	KeyManagement     string `json:"key_management"`

	// Incremental Backup
	IncrementalEnabled bool       `json:"incremental_enabled"`
	LastFullBackup     *time.Time `json:"last_full_backup"`
	IncrementalCount   int        `json:"incremental_count"`
	DifferentialBackup bool       `json:"differential_backup"`

	// Failure Tracking
	BackupFailures  int        `json:"backup_failures"`
	LastFailureDate *time.Time `json:"last_failure_date"`
	FailureReasons  string     `gorm:"type:json" json:"failure_reasons"` // JSON array
	AlertsEnabled   bool       `json:"alerts_enabled"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceRecoveryOptions manages data recovery capabilities and options
type DeviceRecoveryOptions struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Recovery Capabilities
	RecoveryEnabled    bool   `json:"recovery_enabled"`
	RecoveryMethods    string `gorm:"type:json" json:"recovery_methods"`  // JSON array
	SupportedFormats   string `gorm:"type:json" json:"supported_formats"` // JSON array
	MaxRecoverableSize int64  `json:"max_recoverable_size"`               // bytes

	// Recovery Method Availability
	CloudRecovery        bool `json:"cloud_recovery"`
	LocalRecovery        bool `json:"local_recovery"`
	NetworkRecovery      bool `json:"network_recovery"`
	ProfessionalRecovery bool `json:"professional_recovery"`

	// Recovery Time Estimates
	EstimatedRecoveryTime int    `json:"estimated_recovery_time"` // minutes
	FastRecoveryAvailable bool   `json:"fast_recovery_available"`
	RecoverySpeed         string `json:"recovery_speed"`         // fast, normal, slow
	LastRecoveryDuration  int    `json:"last_recovery_duration"` // minutes

	// Recovery Success Metrics
	RecoveryAttempts     int        `json:"recovery_attempts"`
	SuccessfulRecoveries int        `json:"successful_recoveries"`
	FailedRecoveries     int        `json:"failed_recoveries"`
	SuccessRate          float64    `json:"success_rate"` // percentage
	LastRecoveryDate     *time.Time `json:"last_recovery_date"`
	LastRecoveryStatus   string     `json:"last_recovery_status"`

	// Professional Services
	ProfessionalVendor string `json:"professional_vendor"`
	ServiceLevel       string `json:"service_level"` // standard, priority, emergency
	SLAGuaranteed      bool   `json:"sla_guaranteed"`
	ResponseTime       int    `json:"response_time"` // hours

	// DIY Recovery
	DIYToolsAvailable   bool   `json:"diy_tools_available"`
	RecoverySoftware    string `gorm:"type:json" json:"recovery_software"` // JSON array
	UserGuidesAvailable bool   `json:"user_guides_available"`
	TechnicalSupport    bool   `json:"technical_support"`

	// Cost Estimates
	RecoveryCostEstimate float64 `json:"recovery_cost_estimate"`
	ProfessionalCost     float64 `json:"professional_cost"`
	DIYCost              float64 `json:"diy_cost"`
	EmergencySurcharge   float64 `json:"emergency_surcharge"`

	// Data Priority
	CriticalDataTypes  string `gorm:"type:json" json:"critical_data_types"`  // JSON array
	DataPriorityLevels string `gorm:"type:json" json:"data_priority_levels"` // JSON object
	SelectiveRecovery  bool   `json:"selective_recovery"`

	// Partial Recovery
	PartialRecoverySupported bool   `json:"partial_recovery_supported"`
	RecoverableDataTypes     string `gorm:"type:json" json:"recoverable_data_types"` // JSON array
	UnrecoverableData        string `gorm:"type:json" json:"unrecoverable_data"`     // JSON array

	// Insurance Coverage
	RecoveryInsured   bool    `json:"recovery_insured"`
	InsuranceCoverage float64 `json:"insurance_coverage"`
	Deductible        float64 `json:"deductible"`
	ClaimProcess      string  `json:"claim_process"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceEmergencyLocation manages emergency location sharing and tracking
type DeviceEmergencyLocation struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Emergency Location Sharing
	LocationSharingEnabled bool      `json:"location_sharing_enabled"`
	CurrentLocation        string    `json:"current_location"`  // lat,lng
	LocationAccuracy       float64   `json:"location_accuracy"` // meters
	LastLocationUpdate     time.Time `json:"last_location_update"`

	// Location Accuracy
	GPSEnabled          bool   `json:"gps_enabled"`
	WiFiLocationEnabled bool   `json:"wifi_location_enabled"`
	CellularLocation    bool   `json:"cellular_location"`
	IndoorPositioning   bool   `json:"indoor_positioning"`
	AccuracyLevel       string `json:"accuracy_level"` // high, medium, low

	// Offline Capabilities
	OfflineLocationEnabled bool       `json:"offline_location_enabled"`
	LastKnownLocation      string     `json:"last_known_location"`
	OfflineMapsCached      bool       `json:"offline_maps_cached"`
	CacheUpdateDate        *time.Time `json:"cache_update_date"`

	// Emergency Beacon
	BeaconEnabled        bool       `json:"beacon_enabled"`
	BeaconActivated      bool       `json:"beacon_activated"`
	BeaconActivationTime *time.Time `json:"beacon_activation_time"`
	BeaconBatteryLevel   int        `json:"beacon_battery_level"` // percentage
	BeaconRange          float64    `json:"beacon_range"`         // meters

	// Location History
	LocationHistoryEnabled bool   `json:"location_history_enabled"`
	HistoryRetentionDays   int    `json:"history_retention_days"`
	LocationHistory        string `gorm:"type:json" json:"location_history"` // JSON array
	MovementPattern        string `json:"movement_pattern"`

	// Geofencing
	GeofencingEnabled bool       `json:"geofencing_enabled"`
	GeofenceZones     string     `gorm:"type:json" json:"geofence_zones"`  // JSON array
	GeofenceAlerts    string     `gorm:"type:json" json:"geofence_alerts"` // JSON array
	LastGeofenceEvent *time.Time `json:"last_geofence_event"`

	// Safe Zones
	SafeZones             string `gorm:"type:json" json:"safe_zones"` // JSON array
	CurrentZoneStatus     string `json:"current_zone_status"`         // safe, unsafe, unknown
	SafeZoneNotifications bool   `json:"safe_zone_notifications"`

	// Emergency Routing
	EmergencyRouting     bool   `json:"emergency_routing"`
	NearestHospital      string `json:"nearest_hospital"`
	NearestPoliceStation string `json:"nearest_police_station"`
	EvacuationRoute      string `gorm:"type:json" json:"evacuation_route"` // JSON object
	EstimatedArrival     int    `json:"estimated_arrival"`                 // minutes

	// Sharing Permissions
	SharingWith      string     `gorm:"type:json" json:"sharing_with"` // JSON array of contact IDs
	TemporarySharing bool       `json:"temporary_sharing"`
	SharingExpiry    *time.Time `json:"sharing_expiry"`
	ShareLevel       string     `json:"share_level"` // exact, approximate, city

	// Emergency Services Integration
	EmergencyServicesLinked bool       `json:"emergency_services_linked"`
	DispatchEnabled         bool       `json:"dispatch_enabled"`
	AutomaticDispatch       bool       `json:"automatic_dispatch"`
	LastDispatchTime        *time.Time `json:"last_dispatch_time"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DevicePanicMode manages panic button and SOS features
type DevicePanicMode struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Panic Button Configuration
	PanicButtonEnabled bool       `json:"panic_button_enabled"`
	ActivationMethod   string     `json:"activation_method"` // single, double, triple, long_press
	ActivationCount    int        `json:"activation_count"`
	LastActivation     *time.Time `json:"last_activation"`
	FalseAlarmCount    int        `json:"false_alarm_count"`

	// Emergency SOS
	SOSEnabled        bool   `json:"sos_enabled"`
	SOSMessage        string `json:"sos_message"`
	SOSContacts       string `gorm:"type:json" json:"sos_contacts"` // JSON array
	AutoDialEmergency bool   `json:"auto_dial_emergency"`
	CountdownTimer    int    `json:"countdown_timer"` // seconds before activation

	// Silent Alarm
	SilentAlarmEnabled bool   `json:"silent_alarm_enabled"`
	SilentNotification string `gorm:"type:json" json:"silent_notification"` // JSON array of contacts
	DiscreetMode       bool   `json:"discreet_mode"`
	HiddenActivation   bool   `json:"hidden_activation"`

	// Automatic Calls
	AutoCallEnabled bool   `json:"auto_call_enabled"`
	CallPriority    string `gorm:"type:json" json:"call_priority"` // JSON array ordered
	CallAttempts    int    `json:"call_attempts"`
	RetryInterval   int    `json:"retry_interval"` // seconds

	// Emergency Recording
	RecordingEnabled  bool   `json:"recording_enabled"`
	AudioRecording    bool   `json:"audio_recording"`
	VideoRecording    bool   `json:"video_recording"`
	RecordingDuration int    `json:"recording_duration"` // seconds
	AutoUpload        bool   `json:"auto_upload"`
	RecordingStorage  string `json:"recording_storage"` // local, cloud, both

	// Distress Signal
	DistressSignalEnabled bool    `json:"distress_signal_enabled"`
	SignalType            string  `json:"signal_type"`      // sms, email, app, all
	BroadcastRadius       float64 `json:"broadcast_radius"` // km
	SignalFrequency       int     `json:"signal_frequency"` // minutes

	// Emergency Lockdown
	LockdownEnabled   bool `json:"lockdown_enabled"`
	LockdownActivated bool `json:"lockdown_activated"`
	DataWipeEnabled   bool `json:"data_wipe_enabled"`
	RemoteLockEnabled bool `json:"remote_lock_enabled"`
	LocationOnlyMode  bool `json:"location_only_mode"`

	// Data Protection
	EncryptOnPanic      bool `json:"encrypt_on_panic"`
	HideSensitiveData   bool `json:"hide_sensitive_data"`
	BackupBeforeWipe    bool `json:"backup_before_wipe"`
	SecureDeleteEnabled bool `json:"secure_delete_enabled"`

	// Notification Cascade
	CascadeEnabled    bool   `json:"cascade_enabled"`
	NotificationOrder string `gorm:"type:json" json:"notification_order"` // JSON array
	EscalationLevels  string `gorm:"type:json" json:"escalation_levels"`  // JSON object
	MaxNotifications  int    `json:"max_notifications"`

	// Post-Emergency
	PostEmergencyProtocol string `gorm:"type:json" json:"post_emergency_protocol"` // JSON object
	SafeWordRequired      bool   `json:"safe_word_required"`
	DeactivationCode      string `json:"deactivation_code"` // encrypted
	FollowUpRequired      bool   `json:"follow_up_required"`
	IncidentReporting     bool   `json:"incident_reporting"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceDisasterRecovery manages disaster recovery planning and execution
type DeviceDisasterRecovery struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Recovery Planning
	DRPlanActive   bool      `json:"dr_plan_active"`
	PlanVersion    string    `json:"plan_version"`
	LastPlanUpdate time.Time `json:"last_plan_update"`
	NextReviewDate time.Time `json:"next_review_date"`
	PlanCompliance float64   `json:"plan_compliance"` // percentage

	// Critical Data
	CriticalDataIdentified bool   `json:"critical_data_identified"`
	CriticalDataTypes      string `gorm:"type:json" json:"critical_data_types"` // JSON array
	DataClassification     string `gorm:"type:json" json:"data_classification"` // JSON object
	TotalCriticalData      int64  `json:"total_critical_data"`                  // bytes

	// Recovery Priority
	RecoveryTiers string `gorm:"type:json" json:"recovery_tiers"` // JSON object
	Tier1Systems  string `gorm:"type:json" json:"tier1_systems"`  // JSON array
	Tier2Systems  string `gorm:"type:json" json:"tier2_systems"`  // JSON array
	Tier3Systems  string `gorm:"type:json" json:"tier3_systems"`  // JSON array

	// Alternative Device
	AlternativeDevice   bool   `json:"alternative_device"`
	DeviceProvisionTime int    `json:"device_provision_time"` // hours
	PreConfigured       bool   `json:"pre_configured"`
	HotStandby          bool   `json:"hot_standby"`
	DeviceLocation      string `json:"device_location"`

	// Emergency Communication
	EmergencyChannels string `gorm:"type:json" json:"emergency_channels"` // JSON array
	PrimaryChannel    string `json:"primary_channel"`
	BackupChannels    string `gorm:"type:json" json:"backup_channels"`    // JSON array
	CommunicationTree string `gorm:"type:json" json:"communication_tree"` // JSON object

	// Business Continuity
	BCPIntegrated        bool   `json:"bcp_integrated"`
	BusinessImpactLevel  string `json:"business_impact_level"`  // critical, high, medium, low
	MaxTolerableDowntime int    `json:"max_tolerable_downtime"` // hours
	WorkAroundAvailable  bool   `json:"work_around_available"`

	// DR Testing
	LastDRTest      *time.Time `json:"last_dr_test"`
	TestFrequency   string     `json:"test_frequency"`                  // monthly, quarterly, annual
	TestScenarios   string     `gorm:"type:json" json:"test_scenarios"` // JSON array
	TestResults     string     `gorm:"type:json" json:"test_results"`   // JSON array
	TestSuccessRate float64    `json:"test_success_rate"`

	// Recovery Objectives
	RTO                int `json:"rto"`                  // Recovery Time Objective (hours)
	RPO                int `json:"rpo"`                  // Recovery Point Objective (hours)
	ActualRecoveryTime int `json:"actual_recovery_time"` // hours
	DataLossAcceptable int `json:"data_loss_acceptable"` // hours

	// Cost Management
	DRBudget          float64 `json:"dr_budget"`
	EstimatedDRCost   float64 `json:"estimated_dr_cost"`
	ActualDRCost      float64 `json:"actual_dr_cost"`
	CostBenefitRatio  float64 `json:"cost_benefit_ratio"`
	InsuranceCoverage float64 `json:"insurance_coverage"`

	// Recovery Status
	LastDisasterEvent  *time.Time `json:"last_disaster_event"`
	RecoveryInProgress bool       `json:"recovery_in_progress"`
	RecoveryStartTime  *time.Time `json:"recovery_start_time"`
	RecoveryEndTime    *time.Time `json:"recovery_end_time"`
	RecoveryStatus     string     `json:"recovery_status"` // not_started, in_progress, completed, failed

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceEmergencyContacts
func (dec *DeviceEmergencyContacts) HasPrimaryContact() bool {
	return dec.PrimaryContactName != "" && dec.PrimaryContactPhone != ""
}

func (dec *DeviceEmergencyContacts) IsVerified() bool {
	return dec.VerificationStatus == "verified" && dec.PrimaryVerified
}

func (dec *DeviceEmergencyContacts) NeedsUpdate() bool {
	if dec.NextReviewDate != nil {
		return time.Since(*dec.NextReviewDate) > 0
	}
	return time.Since(dec.LastUpdated) > 365*24*time.Hour // Update annually
}

func (dec *DeviceEmergencyContacts) GetEmergencyNumber() string {
	if dec.LocalEmergencyNumber != "" {
		return dec.LocalEmergencyNumber
	}
	// Default emergency numbers by region
	switch dec.CountryCode {
	case "US", "CA":
		return "911"
	case "GB":
		return "999"
	case "EU":
		return "112"
	default:
		return "112" // International standard
	}
}

func (dec *DeviceEmergencyContacts) HasMedicalInfo() bool {
	return dec.BloodType != "" || dec.Allergies != "[]" || dec.MedicalConditions != "[]"
}

// Methods for DeviceBackupStatus
func (dbs *DeviceBackupStatus) IsBackupCurrent() bool {
	if dbs.LastBackupDate == nil {
		return false
	}
	// Check if backup is within frequency window
	switch dbs.BackupFrequency {
	case "hourly":
		return time.Since(*dbs.LastBackupDate) < 2*time.Hour
	case "daily":
		return time.Since(*dbs.LastBackupDate) < 48*time.Hour
	case "weekly":
		return time.Since(*dbs.LastBackupDate) < 14*24*time.Hour
	default:
		return time.Since(*dbs.LastBackupDate) < 30*24*time.Hour
	}
}

func (dbs *DeviceBackupStatus) GetStorageUsagePercent() float64 {
	if dbs.StorageLimit > 0 {
		return float64(dbs.StorageUsed) / float64(dbs.StorageLimit) * 100
	}
	return 0
}

func (dbs *DeviceBackupStatus) IsStorageFull() bool {
	return dbs.GetStorageUsagePercent() >= 95
}

func (dbs *DeviceBackupStatus) HasFailures() bool {
	return dbs.BackupFailures > 0 || dbs.BackupStatus == "failed"
}

func (dbs *DeviceBackupStatus) IsIntegrityValid() bool {
	return dbs.IntegrityStatus == "verified" && dbs.ChecksumVerified
}

// Methods for DeviceRecoveryOptions
func (dro *DeviceRecoveryOptions) CanRecover() bool {
	return dro.RecoveryEnabled && (dro.CloudRecovery || dro.LocalRecovery || dro.NetworkRecovery)
}

func (dro *DeviceRecoveryOptions) GetSuccessRate() float64 {
	if dro.RecoveryAttempts > 0 {
		dro.SuccessRate = float64(dro.SuccessfulRecoveries) / float64(dro.RecoveryAttempts) * 100
	}
	return dro.SuccessRate
}

func (dro *DeviceRecoveryOptions) IsFastRecovery() bool {
	return dro.FastRecoveryAvailable && dro.RecoverySpeed == "fast"
}

func (dro *DeviceRecoveryOptions) GetTotalCost() float64 {
	if dro.ServiceLevel == "emergency" {
		return dro.ProfessionalCost + dro.EmergencySurcharge
	}
	return dro.RecoveryCostEstimate
}

func (dro *DeviceRecoveryOptions) IsInsured() bool {
	return dro.RecoveryInsured && dro.InsuranceCoverage > 0
}

// Methods for DeviceEmergencyLocation
func (del *DeviceEmergencyLocation) IsLocationShared() bool {
	return del.LocationSharingEnabled && del.CurrentLocation != ""
}

func (del *DeviceEmergencyLocation) IsHighAccuracy() bool {
	return del.AccuracyLevel == "high" && del.LocationAccuracy < 10 // within 10 meters
}

func (del *DeviceEmergencyLocation) IsBeaconActive() bool {
	return del.BeaconEnabled && del.BeaconActivated && del.BeaconBatteryLevel > 10
}

func (del *DeviceEmergencyLocation) IsInSafeZone() bool {
	return del.CurrentZoneStatus == "safe"
}

func (del *DeviceEmergencyLocation) CanDispatchEmergency() bool {
	return del.EmergencyServicesLinked && del.DispatchEnabled && del.IsLocationShared()
}

func (del *DeviceEmergencyLocation) HasOfflineCapability() bool {
	return del.OfflineLocationEnabled && del.OfflineMapsCached
}

// Methods for DevicePanicMode
func (dpm *DevicePanicMode) IsPanicActive() bool {
	return dpm.PanicButtonEnabled && dpm.LastActivation != nil &&
		time.Since(*dpm.LastActivation) < 30*time.Minute
}

func (dpm *DevicePanicMode) CanActivateSOS() bool {
	return dpm.SOSEnabled && dpm.SOSContacts != "[]"
}

func (dpm *DevicePanicMode) IsRecording() bool {
	return dpm.RecordingEnabled && (dpm.AudioRecording || dpm.VideoRecording)
}

func (dpm *DevicePanicMode) IsLockdownActive() bool {
	return dpm.LockdownEnabled && dpm.LockdownActivated
}

func (dpm *DevicePanicMode) HasFalseAlarms() bool {
	return dpm.FalseAlarmCount > 3 // More than 3 false alarms
}

func (dpm *DevicePanicMode) GetActivationSensitivity() string {
	switch dpm.ActivationMethod {
	case "single":
		return "high"
	case "double":
		return "medium"
	case "triple", "long_press":
		return "low"
	default:
		return "medium"
	}
}

// Methods for DeviceDisasterRecovery
func (ddr *DeviceDisasterRecovery) IsDRPlanCurrent() bool {
	return ddr.DRPlanActive && time.Since(ddr.LastPlanUpdate) < 90*24*time.Hour // Within 3 months
}

func (ddr *DeviceDisasterRecovery) MeetsRTO() bool {
	return ddr.ActualRecoveryTime <= ddr.RTO
}

func (ddr *DeviceDisasterRecovery) MeetsRPO() bool {
	return ddr.DataLossAcceptable <= ddr.RPO
}

func (ddr *DeviceDisasterRecovery) NeedsTesting() bool {
	if ddr.LastDRTest == nil {
		return true
	}
	switch ddr.TestFrequency {
	case "monthly":
		return time.Since(*ddr.LastDRTest) > 30*24*time.Hour
	case "quarterly":
		return time.Since(*ddr.LastDRTest) > 90*24*time.Hour
	case "annual":
		return time.Since(*ddr.LastDRTest) > 365*24*time.Hour
	default:
		return time.Since(*ddr.LastDRTest) > 180*24*time.Hour
	}
}

func (ddr *DeviceDisasterRecovery) GetRecoveryEfficiency() float64 {
	if ddr.EstimatedDRCost > 0 {
		return (ddr.EstimatedDRCost - ddr.ActualDRCost) / ddr.EstimatedDRCost * 100
	}
	return 0
}

func (ddr *DeviceDisasterRecovery) IsCritical() bool {
	return ddr.BusinessImpactLevel == "critical" && ddr.MaxTolerableDowntime < 4
}

func (ddr *DeviceDisasterRecovery) HasAlternativeDevice() bool {
	return ddr.AlternativeDevice && (ddr.PreConfigured || ddr.HotStandby)
}
