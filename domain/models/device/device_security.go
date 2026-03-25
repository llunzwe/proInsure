package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceSecurityAudit tracks comprehensive security assessments
type DeviceSecurityAudit struct {
	database.BaseModel
	DeviceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	AuditDate time.Time `json:"audit_date"`

	// Vulnerability Assessment
	VulnerabilitiesFound int    `json:"vulnerabilities_found"`
	CriticalVulns        int    `json:"critical_vulns"`
	HighVulns            int    `json:"high_vulns"`
	MediumVulns          int    `json:"medium_vulns"`
	LowVulns             int    `json:"low_vulns"`
	VulnDetails          string `gorm:"type:json" json:"vuln_details"` // JSON array

	// Security Patch Status
	PatchLevel             string     `json:"patch_level"`
	LastPatchDate          *time.Time `json:"last_patch_date"`
	PendingPatches         int        `json:"pending_patches"`
	CriticalPatchesMissing int        `json:"critical_patches_missing"`
	AutoUpdateEnabled      bool       `json:"auto_update_enabled"`
	PatchCompliance        float64    `json:"patch_compliance"` // percentage

	// Scan Results
	LastScanDate    time.Time `json:"last_scan_date"`
	ScanType        string    `json:"scan_type"`     // full, quick, custom
	ScanDuration    int       `json:"scan_duration"` // minutes
	ThreatsDetected int       `json:"threats_detected"`
	ThreatsResolved int       `json:"threats_resolved"`
	ScanHistory     string    `gorm:"type:json" json:"scan_history"` // JSON array

	// Configuration Audit
	ConfigCompliance      float64 `json:"config_compliance"` // percentage
	MisconfigurationCount int     `json:"misconfiguration_count"`
	ConfigViolations      string  `gorm:"type:json" json:"config_violations"` // JSON array
	HardeningScore        float64 `json:"hardening_score"`                    // 0-100
	BestPracticesFollowed int     `json:"best_practices_followed"`

	// Threat Detection
	ThreatDetectionEnabled bool       `json:"threat_detection_enabled"`
	ActiveThreats          int        `json:"active_threats"`
	ThreatMitigationStatus string     `json:"threat_mitigation_status"`
	LastThreatDetected     *time.Time `json:"last_threat_detected"`
	ThreatResponseTime     int        `json:"threat_response_time"` // minutes

	// Security Incidents
	IncidentCount       int        `json:"incident_count"`
	UnresolvedIncidents int        `json:"unresolved_incidents"`
	LastIncident        *time.Time `json:"last_incident"`
	IncidentHistory     string     `gorm:"type:json" json:"incident_history"` // JSON array
	IncidentSeverity    string     `json:"incident_severity"`                 // low, medium, high, critical

	// Compliance Results
	ComplianceFrameworks string  `gorm:"type:json" json:"compliance_frameworks"` // JSON array
	ComplianceStatus     string  `json:"compliance_status"`                      // compliant, partial, non-compliant
	ComplianceScore      float64 `json:"compliance_score"`                       // 0-100
	ComplianceGaps       string  `gorm:"type:json" json:"compliance_gaps"`       // JSON array

	// Security Score
	OverallSecurityScore float64 `json:"overall_security_score"`           // 0-100
	ScoreBreakdown       string  `gorm:"type:json" json:"score_breakdown"` // JSON object
	ScoreTrend           string  `json:"score_trend"`                      // improving, stable, declining
	IndustryBenchmark    float64 `json:"industry_benchmark"`

	// Risk Mitigation
	RiskLevel              string  `json:"risk_level"`                       // low, medium, high, critical
	MitigationPlan         string  `gorm:"type:json" json:"mitigation_plan"` // JSON array
	MitigationProgress     float64 `json:"mitigation_progress"`              // percentage
	EstimatedTimeToResolve int     `json:"estimated_time_to_resolve"`        // hours

	// Recommendations
	RecommendationsCount int    `json:"recommendations_count"`
	ImplementedRecs      int    `json:"implemented_recs"`
	PendingRecs          string `gorm:"type:json" json:"pending_recs"`  // JSON array
	CriticalRecs         string `gorm:"type:json" json:"critical_recs"` // JSON array

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceBiometricData manages biometric authentication settings
type DeviceBiometricData struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UserID   uuid.UUID `gorm:"type:uuid;index" json:"user_id"`

	// Authentication History
	TotalAuthentications int        `json:"total_authentications"`
	SuccessfulAuths      int        `json:"successful_auths"`
	FailedAuths          int        `json:"failed_auths"`
	LastAuthTime         *time.Time `json:"last_auth_time"`
	AuthHistory          string     `gorm:"type:json" json:"auth_history"` // JSON array

	// Fingerprint Data
	FingerprintEnabled   bool       `json:"fingerprint_enabled"`
	FingerprintsEnrolled int        `json:"fingerprints_enrolled"`
	FingerprintSensor    string     `json:"fingerprint_sensor"` // optical, ultrasonic, capacitive
	LastFingerprintAuth  *time.Time `json:"last_fingerprint_auth"`
	FingerprintFailures  int        `json:"fingerprint_failures"`

	// Face Recognition
	FaceIDEnabled    bool       `json:"face_id_enabled"`
	FaceIDEnrolled   bool       `json:"face_id_enrolled"`
	FaceIDVersion    string     `json:"face_id_version"`
	LastFaceAuth     *time.Time `json:"last_face_auth"`
	FaceAuthFailures int        `json:"face_auth_failures"`
	RequireAttention bool       `json:"require_attention"`

	// Voice Recognition
	VoiceAuthEnabled    bool       `json:"voice_auth_enabled"`
	VoiceProfileTrained bool       `json:"voice_profile_trained"`
	VoiceSamplesCount   int        `json:"voice_samples_count"`
	LastVoiceAuth       *time.Time `json:"last_voice_auth"`
	VoiceAuthFailures   int        `json:"voice_auth_failures"`

	// Behavioral Biometrics
	BehavioralEnabled bool    `json:"behavioral_enabled"`
	TypingPatternSet  bool    `json:"typing_pattern_set"`
	SwipePatternSet   bool    `json:"swipe_pattern_set"`
	GaitPatternSet    bool    `json:"gait_pattern_set"`
	BehavioralScore   float64 `json:"behavioral_score"` // 0-100

	// Failure Tracking
	ConsecutiveFailures int        `json:"consecutive_failures"`
	LockoutActive       bool       `json:"lockout_active"`
	LockoutEndTime      *time.Time `json:"lockout_end_time"`
	MaxFailuresAllowed  int        `json:"max_failures_allowed"`
	FailureResetTime    int        `json:"failure_reset_time"` // minutes

	// Fallback Authentication
	FallbackEnabled    bool       `json:"fallback_enabled"`
	FallbackMethod     string     `json:"fallback_method"` // pin, pattern, password
	FallbackUsageCount int        `json:"fallback_usage_count"`
	LastFallbackUsed   *time.Time `json:"last_fallback_used"`

	// Security Incidents
	SecurityIncidents     int        `json:"security_incidents"`
	SpoofingAttempts      int        `json:"spoofing_attempts"`
	CompromisedBiometrics bool       `json:"compromised_biometrics"`
	LastIncident          *time.Time `json:"last_incident"`

	// Multi-Factor Authentication
	MFAEnabled    bool   `json:"mfa_enabled"`
	MFAMethods    string `gorm:"type:json" json:"mfa_methods"` // JSON array
	MFAStrength   string `json:"mfa_strength"`                 // weak, medium, strong
	BackupCodes   int    `json:"backup_codes"`
	RecoveryEmail string `json:"recovery_email"`

	// Data Integrity
	BiometricHash     string    `json:"biometric_hash"`
	LastUpdate        time.Time `json:"last_update"`
	DataEncrypted     bool      `json:"data_encrypted"`
	IntegrityVerified bool      `json:"integrity_verified"`
	TemplateVersion   string    `json:"template_version"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
	// User should be loaded via service layer using UserID to avoid circular import
}

// DeviceEncryptionStatus tracks device encryption configuration
type DeviceEncryptionStatus struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`

	// Device Encryption Level
	FullDiskEncryption bool      `json:"full_disk_encryption"`
	EncryptionLevel    string    `json:"encryption_level"` // none, partial, full
	EncryptionType     string    `json:"encryption_type"`  // software, hardware, hybrid
	EncryptionEnabled  time.Time `json:"encryption_enabled"`

	// Data-at-Rest Encryption
	StorageEncrypted     bool   `json:"storage_encrypted"`
	FileSystemEncryption string `json:"file_system_encryption"` // APFS, ext4, NTFS
	DatabaseEncryption   bool   `json:"database_encryption"`
	BackupEncryption     bool   `json:"backup_encryption"`
	CacheEncryption      bool   `json:"cache_encryption"`

	// Data-in-Transit Encryption
	NetworkEncryption   bool   `json:"network_encryption"`
	TLSVersion          string `json:"tls_version"` // 1.2, 1.3
	VPNEncryption       bool   `json:"vpn_encryption"`
	MessagingEncryption bool   `json:"messaging_encryption"`
	EmailEncryption     bool   `json:"email_encryption"`

	// Key Management
	KeyAlgorithm       string     `json:"key_algorithm"` // AES, RSA, ECC
	KeyLength          int        `json:"key_length"`    // bits
	KeyRotationEnabled bool       `json:"key_rotation_enabled"`
	LastKeyRotation    *time.Time `json:"last_key_rotation"`
	KeyStorageMethod   string     `json:"key_storage_method"` // TPM, secure_enclave, software

	// Algorithm Details
	EncryptionAlgorithm string `json:"encryption_algorithm"` // AES-256, AES-128, ChaCha20
	HashingAlgorithm    string `json:"hashing_algorithm"`    // SHA-256, SHA-512
	CipherMode          string `json:"cipher_mode"`          // CBC, GCM, CTR
	PaddingScheme       string `json:"padding_scheme"`

	// Partition Status
	SystemPartitionEncrypted   bool `json:"system_partition_encrypted"`
	DataPartitionEncrypted     bool `json:"data_partition_encrypted"`
	RecoveryPartitionEncrypted bool `json:"recovery_partition_encrypted"`
	ExternalStorageEncrypted   bool `json:"external_storage_encrypted"`

	// Secure Boot
	SecureBootEnabled bool   `json:"secure_boot_enabled"`
	BootLoaderLocked  bool   `json:"boot_loader_locked"`
	VerifiedBoot      string `json:"verified_boot"` // disabled, orange, yellow, green
	TrustedBootChain  bool   `json:"trusted_boot_chain"`

	// Hardware Encryption
	HardwareEncryption   bool   `json:"hardware_encryption"`
	CryptoProcessor      string `json:"crypto_processor"` // TPM, secure_element
	HardwareAcceleration bool   `json:"hardware_acceleration"`
	FIPSCompliant        bool   `json:"fips_compliant"`

	// Performance Impact
	EncryptionOverhead float64 `json:"encryption_overhead"` // percentage
	CPUUsageImpact     float64 `json:"cpu_usage_impact"`    // percentage
	BatteryImpact      float64 `json:"battery_impact"`      // percentage
	StorageImpact      float64 `json:"storage_impact"`      // percentage

	// Compliance Status
	ComplianceStandards  string    `gorm:"type:json" json:"compliance_standards"` // JSON array
	EncryptionCompliant  bool      `json:"encryption_compliant"`
	LastComplianceCheck  time.Time `json:"last_compliance_check"`
	ComplianceViolations string    `gorm:"type:json" json:"compliance_violations"` // JSON array

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceAntivirusStatus tracks malware protection status
type DeviceAntivirusStatus struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Protection Status
	AntivirusInstalled bool       `json:"antivirus_installed"`
	AntivirusSoftware  string     `json:"antivirus_software"`
	ProtectionActive   bool       `json:"protection_active"`
	LicenseValid       bool       `json:"license_valid"`
	LicenseExpiry      *time.Time `json:"license_expiry"`

	// Definition Updates
	DefinitionsUpToDate  bool      `json:"definitions_up_to_date"`
	LastDefinitionUpdate time.Time `json:"last_definition_update"`
	DefinitionVersion    string    `json:"definition_version"`
	AutoUpdateEnabled    bool      `json:"auto_update_enabled"`
	UpdateFrequency      string    `json:"update_frequency"` // hourly, daily, weekly

	// Scan History
	LastFullScan         *time.Time `json:"last_full_scan"`
	LastQuickScan        *time.Time `json:"last_quick_scan"`
	TotalScansPerformed  int        `json:"total_scans_performed"`
	ScheduledScanEnabled bool       `json:"scheduled_scan_enabled"`
	ScanSchedule         string     `gorm:"type:json" json:"scan_schedule"` // JSON object

	// Threat Detection
	TotalThreatsDetected int        `json:"total_threats_detected"`
	ActiveThreats        int        `json:"active_threats"`
	ThreatsRemoved       int        `json:"threats_removed"`
	LastThreatDetected   *time.Time `json:"last_threat_detected"`
	ThreatHistory        string     `gorm:"type:json" json:"threat_history"` // JSON array

	// Quarantine Management
	QuarantineEnabled    bool       `json:"quarantine_enabled"`
	ItemsInQuarantine    int        `json:"items_in_quarantine"`
	QuarantineSize       int64      `json:"quarantine_size"` // bytes
	AutoCleanEnabled     bool       `json:"auto_clean_enabled"`
	LastQuarantineReview *time.Time `json:"last_quarantine_review"`

	// Real-Time Protection
	RealTimeProtection bool `json:"real_time_protection"`
	FileMonitoring     bool `json:"file_monitoring"`
	BehaviorMonitoring bool `json:"behavior_monitoring"`
	HeuristicScanning  bool `json:"heuristic_scanning"`
	CloudProtection    bool `json:"cloud_protection"`

	// Web Protection
	WebProtectionEnabled    bool `json:"web_protection_enabled"`
	MaliciousSitesBlocked   int  `json:"malicious_sites_blocked"`
	PhishingAttemptsBlocked int  `json:"phishing_attempts_blocked"`
	SafeBrowsing            bool `json:"safe_browsing"`
	HTTPSScanning           bool `json:"https_scanning"`

	// Email Protection
	EmailProtection     bool `json:"email_protection"`
	SpamFilterEnabled   bool `json:"spam_filter_enabled"`
	AttachmentScanning  bool `json:"attachment_scanning"`
	EmailThreatsBlocked int  `json:"email_threats_blocked"`

	// App Scanning
	AppScanningEnabled   bool `json:"app_scanning_enabled"`
	AppsScanned          int  `json:"apps_scanned"`
	MaliciousAppsFound   int  `json:"malicious_apps_found"`
	AppPermissionMonitor bool `json:"app_permission_monitor"`

	// Performance
	SystemImpact         float64   `json:"system_impact"`  // percentage
	ScanSpeed            string    `json:"scan_speed"`     // fast, normal, thorough
	ResourceUsage        float64   `json:"resource_usage"` // percentage
	LastPerformanceCheck time.Time `json:"last_performance_check"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceAccessControl manages user permissions and access controls
type DeviceAccessControl struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// User Access Permissions
	UserCount       int    `json:"user_count"`
	AdminUsers      int    `json:"admin_users"`
	StandardUsers   int    `json:"standard_users"`
	GuestUsers      int    `json:"guest_users"`
	UserPermissions string `gorm:"type:json" json:"user_permissions"` // JSON object

	// App Permissions
	TotalApps            int        `json:"total_apps"`
	HighRiskPermissions  int        `json:"high_risk_permissions"`
	AppPermissionAudit   string     `gorm:"type:json" json:"app_permission_audit"` // JSON array
	LastPermissionReview *time.Time `json:"last_permission_review"`
	OverPrivilegedApps   int        `json:"over_privileged_apps"`

	// Administrator Privileges
	AdminAccessEnabled      bool       `json:"admin_access_enabled"`
	RootAccess              bool       `json:"root_access"`
	SudoersCount            int        `json:"sudoers_count"`
	ElevatedPrivileges      string     `gorm:"type:json" json:"elevated_privileges"` // JSON array
	LastPrivilegeEscalation *time.Time `json:"last_privilege_escalation"`

	// Guest Account
	GuestAccountEnabled bool   `json:"guest_account_enabled"`
	GuestRestrictions   string `gorm:"type:json" json:"guest_restrictions"` // JSON object
	GuestSessionActive  bool   `json:"guest_session_active"`
	GuestDataIsolation  bool   `json:"guest_data_isolation"`

	// Remote Access
	RemoteAccessEnabled bool       `json:"remote_access_enabled"`
	RemoteProtocols     string     `gorm:"type:json" json:"remote_protocols"` // JSON array (SSH, RDP, VNC)
	RemoteConnections   int        `json:"remote_connections"`
	LastRemoteAccess    *time.Time `json:"last_remote_access"`
	RemoteAccessLog     string     `gorm:"type:json" json:"remote_access_log"` // JSON array

	// VPN Configuration
	VPNEnabled        bool   `json:"vpn_enabled"`
	VPNProvider       string `json:"vpn_provider"`
	VPNProtocol       string `json:"vpn_protocol"` // OpenVPN, WireGuard, IPSec
	AlwaysOnVPN       bool   `json:"always_on_vpn"`
	VPNConnectionTime int    `json:"vpn_connection_time"` // hours

	// Network Access Controls
	FirewallEnabled     bool   `json:"firewall_enabled"`
	InboundRules        int    `json:"inbound_rules"`
	OutboundRules       int    `json:"outbound_rules"`
	PortsOpen           string `gorm:"type:json" json:"ports_open"` // JSON array
	NetworkSegmentation bool   `json:"network_segmentation"`

	// File System Permissions
	FilePermissionAudit string `gorm:"type:json" json:"file_permission_audit"` // JSON array
	SensitiveFiles      int    `json:"sensitive_files"`
	WorldReadableFiles  int    `json:"world_readable_files"`
	SetuidFiles         int    `json:"setuid_files"`

	// Privacy Settings
	LocationAccess    string `json:"location_access"` // always, while_using, never
	CameraAccess      string `json:"camera_access"`
	MicrophoneAccess  string `json:"microphone_access"`
	ContactsAccess    string `json:"contacts_access"`
	PrivacyViolations int    `json:"privacy_violations"`

	// Access Violations
	ViolationAttempts  int        `json:"violation_attempts"`
	UnauthorizedAccess int        `json:"unauthorized_access"`
	LastViolation      *time.Time `json:"last_violation"`
	ViolationLog       string     `gorm:"type:json" json:"violation_log"` // JSON array
	SecurityAlerts     int        `json:"security_alerts"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceThreatIntelligence provides threat landscape monitoring
type DeviceThreatIntelligence struct {
	database.BaseModel
	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	LastUpdated time.Time `json:"last_updated"`

	// Threat Landscape
	ThreatLevel     string `json:"threat_level"` // low, medium, high, critical
	ActiveThreats   int    `json:"active_threats"`
	EmergingThreats int    `json:"emerging_threats"`
	ThreatTrends    string `gorm:"type:json" json:"threat_trends"`    // JSON array
	RegionalThreats string `gorm:"type:json" json:"regional_threats"` // JSON array

	// Device-Specific Threats
	TargetedThreats       int     `json:"targeted_threats"`
	DeviceVulnerabilities int     `json:"device_vulnerabilities"`
	OSSpecificThreats     int     `json:"os_specific_threats"`
	AppSpecificThreats    int     `json:"app_specific_threats"`
	ThreatRelevanceScore  float64 `json:"threat_relevance_score"` // 0-100

	// Zero-Day Tracking
	ZeroDayVulns        int    `json:"zero_day_vulns"`
	ActiveExploits      int    `json:"active_exploits"`
	PatchAvailable      string `gorm:"type:json" json:"patch_available"`      // JSON array
	MitigationAvailable string `gorm:"type:json" json:"mitigation_available"` // JSON array
	ExploitComplexity   string `json:"exploit_complexity"`                    // low, medium, high

	// Threat Actors
	ThreatActorsTargeting int    `json:"threat_actors_targeting"`
	APTGroups             string `gorm:"type:json" json:"apt_groups"`           // JSON array
	CybercriminalGroups   string `gorm:"type:json" json:"cybercriminal_groups"` // JSON array
	ThreatActorTTPs       string `gorm:"type:json" json:"threat_actor_ttps"`    // JSON array

	// Attack Vectors
	CommonAttackVectors string  `gorm:"type:json" json:"common_attack_vectors"` // JSON array
	ExploitedVectors    string  `gorm:"type:json" json:"exploited_vectors"`     // JSON array
	VectorProbability   string  `gorm:"type:json" json:"vector_probability"`    // JSON object
	AttackSurfaceScore  float64 `json:"attack_surface_score"`                   // 0-100

	// Threat Mitigation
	MitigationStatus        string  `json:"mitigation_status"` // none, partial, complete
	MitigatedThreats        int     `json:"mitigated_threats"`
	PendingMitigations      int     `json:"pending_mitigations"`
	MitigationEffectiveness float64 `json:"mitigation_effectiveness"`  // percentage
	EstimatedMitigationTime int     `json:"estimated_mitigation_time"` // hours

	// Security Advisories
	ActiveAdvisories   int       `json:"active_advisories"`
	CriticalAdvisories int       `json:"critical_advisories"`
	AdvisoryList       string    `gorm:"type:json" json:"advisory_list"` // JSON array
	LastAdvisoryCheck  time.Time `json:"last_advisory_check"`

	// Intelligence Feeds
	FeedsSubscribed  int       `json:"feeds_subscribed"`
	FeedSources      string    `gorm:"type:json" json:"feed_sources"` // JSON array
	LastFeedUpdate   time.Time `json:"last_feed_update"`
	FeedQualityScore float64   `json:"feed_quality_score"` // 0-100

	// Predictive Analysis
	PredictedThreats  string  `gorm:"type:json" json:"predicted_threats"` // JSON array
	RiskForecast      string  `gorm:"type:json" json:"risk_forecast"`     // JSON object
	ThreatProbability float64 `json:"threat_probability"`                 // percentage
	TimeToCompromise  int     `json:"time_to_compromise"`                 // hours

	// Response Readiness
	ResponsePlan       bool       `json:"response_plan"`
	IncidentResponse   string     `json:"incident_response"`   // automatic, manual, hybrid
	ResponseTime       int        `json:"response_time"`       // minutes
	RecoveryCapability string     `json:"recovery_capability"` // low, medium, high
	ResponseDrillDate  *time.Time `json:"response_drill_date"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceSecurityAudit
func (dsa *DeviceSecurityAudit) IsSecure() bool {
	return dsa.OverallSecurityScore >= 80 && dsa.CriticalVulns == 0 && dsa.ActiveThreats == 0
}

func (dsa *DeviceSecurityAudit) HasCriticalVulnerabilities() bool {
	return dsa.CriticalVulns > 0 || dsa.CriticalPatchesMissing > 0
}

func (dsa *DeviceSecurityAudit) IsCompliant() bool {
	return dsa.ComplianceStatus == "compliant" && dsa.ComplianceScore >= 80
}

func (dsa *DeviceSecurityAudit) NeedsImmediateAction() bool {
	return dsa.RiskLevel == "critical" || dsa.UnresolvedIncidents > 0 || dsa.ActiveThreats > 0
}

func (dsa *DeviceSecurityAudit) GetRiskScore() float64 {
	// Higher risk = lower score
	riskFactors := float64(dsa.CriticalVulns*10 + dsa.HighVulns*5 + dsa.MediumVulns*2 + dsa.LowVulns)
	baseScore := 100 - riskFactors
	if baseScore < 0 {
		baseScore = 0
	}
	// Adjust for compliance
	return baseScore * (dsa.ComplianceScore / 100)
}

// Methods for DeviceBiometricData
func (dbd *DeviceBiometricData) HasBiometrics() bool {
	return dbd.FingerprintEnabled || dbd.FaceIDEnabled || dbd.VoiceAuthEnabled || dbd.BehavioralEnabled
}

func (dbd *DeviceBiometricData) GetAuthSuccessRate() float64 {
	if dbd.TotalAuthentications > 0 {
		return float64(dbd.SuccessfulAuths) / float64(dbd.TotalAuthentications) * 100
	}
	return 0
}

func (dbd *DeviceBiometricData) IsLocked() bool {
	return dbd.LockoutActive && (dbd.LockoutEndTime == nil || time.Now().Before(*dbd.LockoutEndTime))
}

func (dbd *DeviceBiometricData) HasStrongAuthentication() bool {
	return dbd.MFAEnabled && dbd.MFAStrength == "strong" && dbd.HasBiometrics()
}

func (dbd *DeviceBiometricData) IsCompromised() bool {
	return dbd.CompromisedBiometrics || dbd.SpoofingAttempts > 5 || dbd.SecurityIncidents > 0
}

// Methods for DeviceEncryptionStatus
func (des *DeviceEncryptionStatus) IsFullyEncrypted() bool {
	return des.FullDiskEncryption && des.StorageEncrypted && des.NetworkEncryption
}

func (des *DeviceEncryptionStatus) HasStrongEncryption() bool {
	return des.KeyLength >= 256 && (des.EncryptionAlgorithm == "AES-256" || des.EncryptionAlgorithm == "ChaCha20")
}

func (des *DeviceEncryptionStatus) IsBootSecure() bool {
	return des.SecureBootEnabled && des.BootLoaderLocked && des.TrustedBootChain
}

func (des *DeviceEncryptionStatus) NeedsKeyRotation() bool {
	if !des.KeyRotationEnabled {
		return true
	}
	if des.LastKeyRotation != nil {
		// Rotate keys every 90 days
		return time.Since(*des.LastKeyRotation) > 90*24*time.Hour
	}
	return true
}

func (des *DeviceEncryptionStatus) GetEncryptionScore() float64 {
	score := 0.0
	if des.FullDiskEncryption {
		score += 25
	}
	if des.StorageEncrypted {
		score += 25
	}
	if des.NetworkEncryption {
		score += 25
	}
	if des.HasStrongEncryption() {
		score += 25
	}
	// Reduce score for performance impact
	score -= des.EncryptionOverhead
	if score < 0 {
		score = 0
	}
	return score
}

// Methods for DeviceAntivirusStatus
func (das *DeviceAntivirusStatus) IsProtected() bool {
	return das.AntivirusInstalled && das.ProtectionActive && das.DefinitionsUpToDate && das.LicenseValid
}

func (das *DeviceAntivirusStatus) HasActiveThreats() bool {
	return das.ActiveThreats > 0 || das.ItemsInQuarantine > 0
}

func (das *DeviceAntivirusStatus) NeedsUpdate() bool {
	return !das.DefinitionsUpToDate || time.Since(das.LastDefinitionUpdate) > 7*24*time.Hour
}

func (das *DeviceAntivirusStatus) NeedsFullScan() bool {
	if das.LastFullScan == nil {
		return true
	}
	return time.Since(*das.LastFullScan) > 30*24*time.Hour // Monthly full scan
}

func (das *DeviceAntivirusStatus) GetProtectionLevel() string {
	if !das.IsProtected() {
		return "unprotected"
	}
	if das.RealTimeProtection && das.WebProtectionEnabled && das.EmailProtection {
		return "maximum"
	}
	if das.RealTimeProtection {
		return "standard"
	}
	return "minimal"
}

// Methods for DeviceAccessControl
func (dac *DeviceAccessControl) HasElevatedPrivileges() bool {
	return dac.AdminAccessEnabled || dac.RootAccess || dac.SudoersCount > 0
}

func (dac *DeviceAccessControl) IsRemoteAccessSecure() bool {
	return !dac.RemoteAccessEnabled || (dac.VPNEnabled && dac.AlwaysOnVPN)
}

func (dac *DeviceAccessControl) HasPrivacyRisks() bool {
	return dac.HighRiskPermissions > 5 || dac.OverPrivilegedApps > 0 || dac.PrivacyViolations > 0
}

func (dac *DeviceAccessControl) GetAccessRiskScore() float64 {
	riskScore := 0.0
	if dac.GuestAccountEnabled {
		riskScore += 10
	}
	if dac.RemoteAccessEnabled && !dac.VPNEnabled {
		riskScore += 20
	}
	if dac.RootAccess {
		riskScore += 30
	}
	riskScore += float64(dac.ViolationAttempts * 2)
	riskScore += float64(dac.UnauthorizedAccess * 5)

	if riskScore > 100 {
		riskScore = 100
	}
	return riskScore
}

func (dac *DeviceAccessControl) NeedsSecurityReview() bool {
	return dac.ViolationAttempts > 10 || dac.UnauthorizedAccess > 0 ||
		dac.SecurityAlerts > 5 || dac.HighRiskPermissions > 10
}

// Methods for DeviceThreatIntelligence
func (dti *DeviceThreatIntelligence) IsCritical() bool {
	return dti.ThreatLevel == "critical" || dti.ZeroDayVulns > 0 || dti.ActiveExploits > 0
}

func (dti *DeviceThreatIntelligence) IsTargeted() bool {
	return dti.TargetedThreats > 0 || dti.ThreatActorsTargeting > 0
}

func (dti *DeviceThreatIntelligence) GetThreatScore() float64 {
	score := float64(dti.ActiveThreats*10 + dti.EmergingThreats*5 + dti.ZeroDayVulns*20)
	score += dti.ThreatRelevanceScore
	score += dti.AttackSurfaceScore

	// Reduce score based on mitigation
	if dti.MitigationStatus == "complete" {
		score *= 0.3
	} else if dti.MitigationStatus == "partial" {
		score *= 0.6
	}

	if score > 100 {
		score = 100
	}
	return score
}

func (dti *DeviceThreatIntelligence) NeedsMitigation() bool {
	return dti.PendingMitigations > 0 || dti.MitigationStatus == "none" || dti.ActiveThreats > 0
}

func (dti *DeviceThreatIntelligence) HasResponsePlan() bool {
	return dti.ResponsePlan && dti.RecoveryCapability != "low" && dti.ResponseTime < 60
}
