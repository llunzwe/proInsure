package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSecurityProfile manages comprehensive security settings and history
type UserSecurityProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Security Scoring
	SecurityScore      float64    `gorm:"default:50" json:"security_score"`
	SecurityLevel      string     `gorm:"type:varchar(20)" json:"security_level"` // low/medium/high/critical
	ThreatLevel        string     `gorm:"type:varchar(20)" json:"threat_level"`
	VulnerabilityScore float64    `gorm:"default:0" json:"vulnerability_score"`
	CompromiseRisk     float64    `gorm:"default:0" json:"compromise_risk"`
	LastSecurityAudit  *time.Time `json:"last_security_audit"`
	NextSecurityReview *time.Time `json:"next_security_review"`

	// Authentication Security
	AuthenticationStrength float64    `gorm:"default:0" json:"authentication_strength"`
	PasswordStrength       float64    `gorm:"default:0" json:"password_strength"`
	PasswordAge            int        `json:"password_age_days"`
	PasswordHistory        []string   `gorm:"type:json" json:"-"` // Hashed passwords
	LastPasswordChange     time.Time  `json:"last_password_change"`
	PasswordExpiryDate     *time.Time `json:"password_expiry_date"`
	RequirePasswordChange  bool       `gorm:"default:false" json:"require_password_change"`

	// Multi-Factor Authentication
	MFAEnabled           bool       `gorm:"default:false" json:"mfa_enabled"`
	MFAMethods           []string   `gorm:"type:json" json:"mfa_methods"`
	PrimaryMFAMethod     string     `gorm:"type:varchar(50)" json:"primary_mfa_method"`
	BackupCodes          []string   `gorm:"type:json" json:"-"` // Encrypted
	BackupCodesGenerated *time.Time `json:"backup_codes_generated"`
	MFARecoveryEmail     string     `gorm:"type:varchar(255)" json:"mfa_recovery_email"`
	MFARecoveryPhone     string     `gorm:"type:varchar(20)" json:"mfa_recovery_phone"`

	// Device Security
	TrustedDevices       []map[string]interface{} `gorm:"type:json" json:"trusted_devices"`
	DeviceFingerprints   []string                 `gorm:"type:json" json:"device_fingerprints"`
	MaxDevices           int                      `gorm:"default:5" json:"max_devices"`
	DeviceRotationPolicy string                   `gorm:"type:varchar(50)" json:"device_rotation_policy"`
	LastDeviceCleanup    *time.Time               `json:"last_device_cleanup"`

	// Access Control
	IPRestrictions     []string               `gorm:"type:json" json:"ip_restrictions"`
	GeofencingEnabled  bool                   `gorm:"default:false" json:"geofencing_enabled"`
	AllowedLocations   []map[string]float64   `gorm:"type:json" json:"allowed_locations"`
	TimeBasedAccess    map[string]string      `gorm:"type:json" json:"time_based_access"`
	SessionManagement  map[string]interface{} `gorm:"type:json" json:"session_management"`
	ConcurrentSessions int                    `gorm:"default:1" json:"concurrent_sessions"`

	// Security Incidents
	IncidentCount      int                      `gorm:"default:0" json:"incident_count"`
	LastIncidentDate   *time.Time               `json:"last_incident_date"`
	ActiveThreats      []string                 `gorm:"type:json" json:"active_threats"`
	MitigationMeasures []string                 `gorm:"type:json" json:"mitigation_measures"`
	SecurityAlerts     []map[string]interface{} `gorm:"type:json" json:"security_alerts"`

	// Breach History
	BreachCount          int        `gorm:"default:0" json:"breach_count"`
	LastBreachDate       *time.Time `json:"last_breach_date"`
	BreachImpact         string     `gorm:"type:text" json:"breach_impact"`
	RecoveryActions      []string   `gorm:"type:json" json:"recovery_actions"`
	CompensationProvided bool       `gorm:"default:false" json:"compensation_provided"`
}

// UserPrivacySettings manages privacy preferences and controls
type UserPrivacySettings struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Privacy Levels
	PrivacyLevel      string `gorm:"type:varchar(20)" json:"privacy_level"` // public/friends/private/custom
	DataSharingLevel  string `gorm:"type:varchar(20)" json:"data_sharing_level"`
	ProfileVisibility string `gorm:"type:varchar(20)" json:"profile_visibility"`
	SearchVisibility  bool   `gorm:"default:true" json:"search_visibility"`
	DirectoryListing  bool   `gorm:"default:false" json:"directory_listing"`

	// Data Sharing Preferences
	ShareWithPartners bool            `gorm:"default:false" json:"share_with_partners"`
	ShareForAnalytics bool            `gorm:"default:true" json:"share_for_analytics"`
	ShareForMarketing bool            `gorm:"default:false" json:"share_for_marketing"`
	ShareForResearch  bool            `gorm:"default:false" json:"share_for_research"`
	DataCategories    map[string]bool `gorm:"type:json" json:"data_categories"`

	// Consent Management
	ConsentHistory    []map[string]interface{} `gorm:"type:json" json:"consent_history"`
	ActiveConsents    map[string]interface{}   `gorm:"type:json" json:"active_consents"`
	ConsentVersion    string                   `gorm:"type:varchar(20)" json:"consent_version"`
	LastConsentUpdate *time.Time               `json:"last_consent_update"`
	ConsentChannels   map[string]bool          `gorm:"type:json" json:"consent_channels"`

	// Data Rights
	AccessRequestCount     int        `gorm:"default:0" json:"access_request_count"`
	LastAccessRequest      *time.Time `json:"last_access_request"`
	DeletionRequested      bool       `gorm:"default:false" json:"deletion_requested"`
	DeletionScheduled      *time.Time `json:"deletion_scheduled"`
	PortabilityRequests    int        `gorm:"default:0" json:"portability_requests"`
	LastPortabilityRequest *time.Time `json:"last_portability_request"`

	// Cookie Preferences
	CookiesAccepted   bool            `gorm:"default:false" json:"cookies_accepted"`
	EssentialCookies  bool            `gorm:"default:true" json:"essential_cookies"`
	FunctionalCookies bool            `gorm:"default:false" json:"functional_cookies"`
	AnalyticsCookies  bool            `gorm:"default:false" json:"analytics_cookies"`
	MarketingCookies  bool            `gorm:"default:false" json:"marketing_cookies"`
	CookiePreferences map[string]bool `gorm:"type:json" json:"cookie_preferences"`

	// Tracking Preferences
	DoNotTrack        bool `gorm:"default:false" json:"do_not_track"`
	LocationTracking  bool `gorm:"default:false" json:"location_tracking"`
	BehaviorTracking  bool `gorm:"default:false" json:"behavior_tracking"`
	CrossSiteTracking bool `gorm:"default:false" json:"cross_site_tracking"`
	DeviceTracking    bool `gorm:"default:false" json:"device_tracking"`

	// Communication Privacy
	HideEmail         bool   `gorm:"default:false" json:"hide_email"`
	HidePhone         bool   `gorm:"default:true" json:"hide_phone"`
	HideAddress       bool   `gorm:"default:true" json:"hide_address"`
	UseProxyEmail     bool   `gorm:"default:false" json:"use_proxy_email"`
	ProxyEmailAddress string `gorm:"type:varchar(255)" json:"proxy_email_address"`
}

// UserAccessControl manages role-based access and permissions
type UserAccessControl struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Roles & Permissions
	Roles             []string        `gorm:"type:json" json:"roles"`
	Permissions       []string        `gorm:"type:json" json:"permissions"`
	CustomPermissions map[string]bool `gorm:"type:json" json:"custom_permissions"`
	PermissionGroups  []string        `gorm:"type:json" json:"permission_groups"`
	AccessLevel       string          `gorm:"type:varchar(20)" json:"access_level"`

	// Delegation
	DelegatedAccess       []map[string]interface{} `gorm:"type:json" json:"delegated_access"`
	DelegatedFrom         []uuid.UUID              `gorm:"type:json" json:"delegated_from"`
	DelegatedTo           []uuid.UUID              `gorm:"type:json" json:"delegated_to"`
	PowerOfAttorney       bool                     `gorm:"default:false" json:"power_of_attorney"`
	PowerOfAttorneyHolder *uuid.UUID               `gorm:"type:uuid" json:"power_of_attorney_holder"`

	// Resource Access
	ResourceAccess map[string][]string `gorm:"type:json" json:"resource_access"`
	DataAccess     map[string]bool     `gorm:"type:json" json:"data_access"`
	FeatureAccess  map[string]bool     `gorm:"type:json" json:"feature_access"`
	APIAccess      map[string]bool     `gorm:"type:json" json:"api_access"`

	// Time-Based Access
	TemporaryAccess []map[string]interface{} `gorm:"type:json" json:"temporary_access"`
	AccessSchedule  map[string]interface{}   `gorm:"type:json" json:"access_schedule"`
	AccessExpiry    *time.Time               `json:"access_expiry"`

	// Audit
	LastPermissionChange    *time.Time               `json:"last_permission_change"`
	PermissionChangeHistory []map[string]interface{} `gorm:"type:json" json:"permission_change_history"`
	AccessReviewRequired    bool                     `gorm:"default:false" json:"access_review_required"`
	LastAccessReview        *time.Time               `json:"last_access_review"`
}

// UserBiometrics manages biometric authentication data
type UserBiometrics struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Biometric Types
	FingerprintEnabled bool `gorm:"default:false" json:"fingerprint_enabled"`
	FaceIDEnabled      bool `gorm:"default:false" json:"face_id_enabled"`
	VoiceEnabled       bool `gorm:"default:false" json:"voice_enabled"`
	IrisEnabled        bool `gorm:"default:false" json:"iris_enabled"`
	BehavioralEnabled  bool `gorm:"default:false" json:"behavioral_enabled"`

	// Enrollment Status
	EnrollmentStatus  map[string]string    `gorm:"type:json" json:"enrollment_status"`
	EnrollmentDates   map[string]time.Time `gorm:"type:json" json:"enrollment_dates"`
	EnrollmentQuality map[string]float64   `gorm:"type:json" json:"enrollment_quality"`

	// Templates (encrypted references only)
	TemplateReferences map[string]string `gorm:"type:json" json:"-"`
	TemplateVersions   map[string]int    `gorm:"type:json" json:"template_versions"`

	// Authentication History
	LastBiometricAuth   *time.Time         `json:"last_biometric_auth"`
	AuthenticationCount map[string]int     `gorm:"type:json" json:"authentication_count"`
	FailedAttempts      map[string]int     `gorm:"type:json" json:"failed_attempts"`
	SuccessRate         map[string]float64 `gorm:"type:json" json:"success_rate"`

	// Device Bindings
	RegisteredDevices  []map[string]interface{} `gorm:"type:json" json:"registered_devices"`
	DeviceLimit        int                      `gorm:"default:3" json:"device_limit"`
	CrossDeviceEnabled bool                     `gorm:"default:false" json:"cross_device_enabled"`

	// Fallback Methods
	FallbackEnabled  bool     `gorm:"default:true" json:"fallback_enabled"`
	FallbackMethods  []string `gorm:"type:json" json:"fallback_methods"`
	RequireBiometric bool     `gorm:"default:false" json:"require_biometric"`

	// Security Settings
	AntiSpoofingEnabled bool    `gorm:"default:true" json:"anti_spoofing_enabled"`
	LivenessDetection   bool    `gorm:"default:true" json:"liveness_detection"`
	MatchThreshold      float64 `gorm:"default:0.95" json:"match_threshold"`
	MaxRetries          int     `gorm:"default:3" json:"max_retries"`
}

// UserEncryption manages encryption keys and settings
type UserEncryption struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Encryption Status
	EncryptionEnabled   bool   `gorm:"default:true" json:"encryption_enabled"`
	EncryptionLevel     string `gorm:"type:varchar(20)" json:"encryption_level"`
	EncryptionAlgorithm string `gorm:"type:varchar(50)" json:"encryption_algorithm"`
	KeyLength           int    `gorm:"default:256" json:"key_length"`

	// Key Management
	PublicKeyID          string     `gorm:"type:varchar(100)" json:"public_key_id"`
	KeyRotationEnabled   bool       `gorm:"default:true" json:"key_rotation_enabled"`
	KeyRotationFrequency int        `json:"key_rotation_frequency_days"`
	LastKeyRotation      *time.Time `json:"last_key_rotation"`
	NextKeyRotation      *time.Time `json:"next_key_rotation"`
	KeyVersion           int        `gorm:"default:1" json:"key_version"`

	// Data Encryption
	EncryptedFields     []string        `gorm:"type:json" json:"encrypted_fields"`
	EncryptionScope     map[string]bool `gorm:"type:json" json:"encryption_scope"`
	EndToEndEnabled     bool            `gorm:"default:false" json:"end_to_end_enabled"`
	AtRestEncryption    bool            `gorm:"default:true" json:"at_rest_encryption"`
	InTransitEncryption bool            `gorm:"default:true" json:"in_transit_encryption"`

	// Recovery
	RecoveryKeyGenerated bool     `gorm:"default:false" json:"recovery_key_generated"`
	RecoveryKeyID        string   `gorm:"type:varchar(100)" json:"recovery_key_id"`
	RecoveryMethods      []string `gorm:"type:json" json:"recovery_methods"`
	EscrowEnabled        bool     `gorm:"default:false" json:"escrow_enabled"`

	// Compliance
	ComplianceStandards []string                 `gorm:"type:json" json:"compliance_standards"`
	CertificationStatus map[string]bool          `gorm:"type:json" json:"certification_status"`
	AuditTrail          []map[string]interface{} `gorm:"type:json" json:"audit_trail"`
	LastAuditDate       *time.Time               `json:"last_audit_date"`
}

// UserThreatIntelligence tracks security threats and assessments
type UserThreatIntelligence struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Threat Assessment
	ThreatScore      float64                  `gorm:"default:0" json:"threat_score"`
	ThreatLevel      string                   `gorm:"type:varchar(20)" json:"threat_level"`
	ActiveThreats    []map[string]interface{} `gorm:"type:json" json:"active_threats"`
	ThreatCategories []string                 `gorm:"type:json" json:"threat_categories"`
	ThreatVectors    []string                 `gorm:"type:json" json:"threat_vectors"`

	// Vulnerability Assessment
	Vulnerabilities []map[string]interface{} `gorm:"type:json" json:"vulnerabilities"`
	PatchStatus     map[string]string        `gorm:"type:json" json:"patch_status"`
	ExposureLevel   float64                  `gorm:"default:0" json:"exposure_level"`
	AttackSurface   map[string]interface{}   `gorm:"type:json" json:"attack_surface"`

	// Incident Response
	IncidentResponsePlan map[string]interface{} `gorm:"type:json" json:"incident_response_plan"`
	ResponseTeamContacts []map[string]string    `gorm:"type:json" json:"response_team_contacts"`
	EscalationProcedure  map[string]interface{} `gorm:"type:json" json:"escalation_procedure"`
	RecoveryProcedure    map[string]interface{} `gorm:"type:json" json:"recovery_procedure"`

	// Threat Intelligence Feed
	ThreatFeeds            []string   `gorm:"type:json" json:"threat_feeds"`
	LastFeedUpdate         *time.Time `json:"last_feed_update"`
	IndicatorsOfCompromise []string   `gorm:"type:json" json:"indicators_of_compromise"`
	BlockedIPs             []string   `gorm:"type:json" json:"blocked_ips"`
	BlockedDomains         []string   `gorm:"type:json" json:"blocked_domains"`

	// Monitoring
	MonitoringEnabled    bool               `gorm:"default:true" json:"monitoring_enabled"`
	AlertThresholds      map[string]float64 `gorm:"type:json" json:"alert_thresholds"`
	NotificationSettings map[string]bool    `gorm:"type:json" json:"notification_settings"`
	LastSecurityScan     *time.Time         `json:"last_security_scan"`
	NextScheduledScan    *time.Time         `json:"next_scheduled_scan"`
}

// TableName returns the table name
func (UserSecurityProfile) TableName() string {
	return "user_security_profiles"
}

// TableName returns the table name
func (UserPrivacySettings) TableName() string {
	return "user_privacy_settings"
}

// TableName returns the table name
func (UserAccessControl) TableName() string {
	return "user_access_control"
}

// TableName returns the table name
func (UserBiometrics) TableName() string {
	return "user_biometrics"
}

// TableName returns the table name
func (UserEncryption) TableName() string {
	return "user_encryption"
}

// TableName returns the table name
func (UserThreatIntelligence) TableName() string {
	return "user_threat_intelligence"
}
