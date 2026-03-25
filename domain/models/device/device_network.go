package device

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// DeviceNetworkProfile tracks carrier history and network metrics
type DeviceNetworkProfile struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Carrier Information
	CurrentCarrier       string     `json:"current_carrier"`
	CarrierSince         time.Time  `json:"carrier_since"`
	PreviousCarriers     string     `gorm:"type:json" json:"previous_carriers"` // JSON array
	CarrierChangeCount   int        `json:"carrier_change_count"`
	LastCarrierChange    *time.Time `json:"last_carrier_change"`
	CarrierLoyaltyMonths int        `json:"carrier_loyalty_months"`

	// Network Type Preferences
	PreferredNetwork  string  `json:"preferred_network"` // 5G, 4G, 3G, auto
	Network5GEnabled  bool    `json:"network_5g_enabled"`
	Network5GUsage    float64 `json:"network_5g_usage"` // percentage
	Network4GUsage    float64 `json:"network_4g_usage"` // percentage
	Network3GUsage    float64 `json:"network_3g_usage"` // percentage
	NetworkAutoSwitch bool    `json:"network_auto_switch"`

	// Roaming Patterns
	RoamingEnabled           bool    `json:"roaming_enabled"`
	DomesticRoamingDays      int     `json:"domestic_roaming_days"`
	InternationalRoamingDays int     `json:"international_roaming_days"`
	RoamingCountries         string  `gorm:"type:json" json:"roaming_countries"` // JSON array
	DataRoamingEnabled       bool    `json:"data_roaming_enabled"`
	RoamingCost              float64 `json:"roaming_cost"`

	// International Usage
	InternationalUsage   bool    `json:"international_usage"`
	CountriesVisited     string  `gorm:"type:json" json:"countries_visited"` // JSON array
	InternationalMinutes int     `json:"international_minutes"`
	InternationalData    float64 `json:"international_data"` // GB
	InternationalSMS     int     `json:"international_sms"`

	// Network Quality Metrics
	AverageSignalStrength int     `json:"average_signal_strength"` // dBm
	SignalQuality         string  `json:"signal_quality"`          // excellent, good, fair, poor
	NetworkReliability    float64 `json:"network_reliability"`     // percentage uptime
	CallQuality           float64 `json:"call_quality"`            // 0-100
	DataSpeed             float64 `json:"data_speed"`              // Mbps average

	// Signal History
	SignalDrops        int    `json:"signal_drops"`
	DeadZoneEncounters int    `json:"dead_zone_encounters"`
	LowSignalDuration  int    `json:"low_signal_duration"`             // minutes
	SignalHistory      string `gorm:"type:json" json:"signal_history"` // JSON array

	// Speed Measurements
	AverageDownload float64 `json:"average_download"` // Mbps
	AverageUpload   float64 `json:"average_upload"`   // Mbps
	PeakDownload    float64 `json:"peak_download"`    // Mbps
	PeakUpload      float64 `json:"peak_upload"`      // Mbps
	LatencyAverage  float64 `json:"latency_average"`  // ms

	// Network Outages
	OutageCount         int        `json:"outage_count"`
	TotalOutageDuration int        `json:"total_outage_duration"` // minutes
	LastOutage          *time.Time `json:"last_outage"`
	OutageHistory       string     `gorm:"type:json" json:"outage_history"` // JSON array

	// Multi-SIM Usage
	DualSIMEnabled      bool   `json:"dual_sim_enabled"`
	eSIMEnabled         bool   `json:"esim_enabled"`
	ActiveSIMCount      int    `json:"active_sim_count"`
	SIMSwitchFrequency  int    `json:"sim_switch_frequency"`
	PrimarySIMCarrier   string `json:"primary_sim_carrier"`
	SecondarySIMCarrier string `json:"secondary_sim_carrier"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceDataUsage tracks monthly data consumption patterns
type DeviceDataUsage struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	UsageMonth time.Time `json:"usage_month"`

	// Monthly Consumption
	TotalDataUsed    float64 `json:"total_data_used"`    // GB
	CellularDataUsed float64 `json:"cellular_data_used"` // GB
	WiFiDataUsed     float64 `json:"wifi_data_used"`     // GB
	DataLimit        float64 `json:"data_limit"`         // GB
	DataRemaining    float64 `json:"data_remaining"`     // GB
	DaysInCycle      int     `json:"days_in_cycle"`

	// Peak Usage Times
	PeakUsageHour  int     `json:"peak_usage_hour"`  // 0-23
	PeakUsageDay   string  `json:"peak_usage_day"`   // day of week
	WeekdayUsage   float64 `json:"weekday_usage"`    // GB
	WeekendUsage   float64 `json:"weekend_usage"`    // GB
	NightTimeUsage float64 `json:"night_time_usage"` // GB

	// App Data Breakdown
	AppDataBreakdown   string  `gorm:"type:json" json:"app_data_breakdown"` // JSON object
	TopDataApp         string  `json:"top_data_app"`
	TopDataAppUsage    float64 `json:"top_data_app_usage"`   // GB
	SocialMediaData    float64 `json:"social_media_data"`    // GB
	VideoStreamingData float64 `json:"video_streaming_data"` // GB
	MusicStreamingData float64 `json:"music_streaming_data"` // GB
	GamingData         float64 `json:"gaming_data"`          // GB

	// Usage Ratios
	WiFiToCellularRatio float64 `json:"wifi_to_cellular_ratio"`
	ForegroundData      float64 `json:"foreground_data"`       // GB
	BackgroundData      float64 `json:"background_data"`       // GB
	BackgroundDataRatio float64 `json:"background_data_ratio"` // percentage

	// Streaming Services
	StreamingServices string  `gorm:"type:json" json:"streaming_services"` // JSON array
	StreamingQuality  string  `json:"streaming_quality"`                   // auto, high, medium, low
	HDStreamingHours  float64 `json:"hd_streaming_hours"`
	SDStreamingHours  float64 `json:"sd_streaming_hours"`
	OfflineDownloads  float64 `json:"offline_downloads"` // GB

	// Data Saving
	DataSaverEnabled bool    `json:"data_saver_enabled"`
	DataSaverDays    int     `json:"data_saver_days"`
	DataSaved        float64 `json:"data_saved"` // GB
	LowDataMode      bool    `json:"low_data_mode"`

	// Hotspot Sharing
	HotspotDataShared  float64 `json:"hotspot_data_shared"` // GB
	HotspotActivations int     `json:"hotspot_activations"`
	HotspotDuration    int     `json:"hotspot_duration"` // minutes
	DevicesConnected   int     `json:"devices_connected"`

	// Overage Incidents
	DataOverages        int     `json:"data_overages"`
	OverageAmount       float64 `json:"overage_amount"` // GB
	OverageCost         float64 `json:"overage_cost"`
	ThrottlingIncidents int     `json:"throttling_incidents"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceConnectivityIssues tracks network problems and issues
type DeviceConnectivityIssues struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Problem History
	TotalIssues    int        `json:"total_issues"`
	ResolvedIssues int        `json:"resolved_issues"`
	OngoingIssues  int        `json:"ongoing_issues"`
	LastIssueDate  *time.Time `json:"last_issue_date"`
	IssueHistory   string     `gorm:"type:json" json:"issue_history"` // JSON array

	// Connection Drops
	ConnectionDrops     int        `json:"connection_drops"`
	DropsPerDay         float64    `json:"drops_per_day"`
	LastDropTime        *time.Time `json:"last_drop_time"`
	AverageDropDuration int        `json:"average_drop_duration"` // seconds

	// Dead Zones
	DeadZoneLocations string `gorm:"type:json" json:"dead_zone_locations"` // JSON array
	DeadZoneFrequency int    `json:"dead_zone_frequency"`
	CommonDeadZones   string `gorm:"type:json" json:"common_dead_zones"` // JSON array

	// Network Switching
	NetworkSwitchFailures int    `json:"network_switch_failures"`
	AutoSwitchIssues      int    `json:"auto_switch_issues"`
	ManualInterventions   int    `json:"manual_interventions"`
	StuckOnNetwork        string `json:"stuck_on_network"`

	// SIM Card Problems
	SIMCardErrors   int        `json:"sim_card_errors"`
	SIMNotDetected  int        `json:"sim_not_detected"`
	SIMCardReplaced bool       `json:"sim_card_replaced"`
	LastSIMError    *time.Time `json:"last_sim_error"`

	// eSIM Issues
	ESIMActivationFailed bool `json:"esim_activation_failed"`
	ESIMTransferIssues   int  `json:"esim_transfer_issues"`
	ESIMProfileErrors    int  `json:"esim_profile_errors"`
	QRCodeScanFailures   int  `json:"qr_code_scan_failures"`

	// Roaming Problems
	RoamingFailures      int    `json:"roaming_failures"`
	RoamingNetworkIssues string `gorm:"type:json" json:"roaming_network_issues"` // JSON array
	RoamingDataBlocked   bool   `json:"roaming_data_blocked"`
	RoamingResolutions   int    `json:"roaming_resolutions"`

	// Compatibility Issues
	NetworkIncompatibility bool   `json:"network_incompatibility"`
	BandSupportIssues      string `gorm:"type:json" json:"band_support_issues"` // JSON array
	FrequencyMismatch      bool   `json:"frequency_mismatch"`
	ProtocolIssues         string `gorm:"type:json" json:"protocol_issues"` // JSON array

	// Antenna Problems
	AntennaIssues       bool `json:"antenna_issues"`
	SignalReceptionWeak bool `json:"signal_reception_weak"`
	AntennaReplaced     bool `json:"antenna_replaced"`
	HardwareDefect      bool `json:"hardware_defect"`

	// Resolution Tracking
	AverageResolutionTime int `json:"average_resolution_time"` // hours
	SelfResolved          int `json:"self_resolved"`
	SupportResolved       int `json:"support_resolved"`
	UnresolvedDuration    int `json:"unresolved_duration"` // days

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceWiFiProfile manages WiFi network connections and preferences
type DeviceWiFiProfile struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Known Networks
	SavedNetworks     int    `json:"saved_networks"`
	KnownNetworkList  string `gorm:"type:json" json:"known_network_list"` // JSON array
	TrustedNetworks   string `gorm:"type:json" json:"trusted_networks"`   // JSON array
	ForgottenNetworks int    `json:"forgotten_networks"`

	// Public WiFi Usage
	PublicWiFiConnections int     `json:"public_wifi_connections"`
	PublicWiFiHours       float64 `json:"public_wifi_hours"`
	CaptivePortalLogins   int     `json:"captive_portal_logins"`
	OpenNetworkUsage      int     `json:"open_network_usage"`

	// Home/Work Networks
	HomeNetworkSSID      string  `json:"home_network_ssid"`
	HomeNetworkStability float64 `json:"home_network_stability"` // percentage
	WorkNetworkSSID      string  `json:"work_network_ssid"`
	WorkNetworkStability float64 `json:"work_network_stability"` // percentage
	PrimaryNetwork       string  `json:"primary_network"`

	// Security Preferences
	SecurityProtocol  string `json:"security_protocol"` // WPA3, WPA2, WEP, Open
	WPA3Enabled       bool   `json:"wpa3_enabled"`
	VPNAutoConnect    bool   `json:"vpn_auto_connect"`
	SecurityWarnings  int    `json:"security_warnings"`
	UnsafeConnections int    `json:"unsafe_connections"`

	// Auto-Connect Settings
	AutoConnectEnabled bool   `json:"auto_connect_enabled"`
	AutoJoinNetworks   string `gorm:"type:json" json:"auto_join_networks"` // JSON array
	AskToJoinNetworks  bool   `json:"ask_to_join_networks"`
	HiddenNetworks     string `gorm:"type:json" json:"hidden_networks"` // JSON array

	// WiFi Calling
	WiFiCallingEnabled bool    `json:"wifi_calling_enabled"`
	WiFiCallMinutes    int     `json:"wifi_call_minutes"`
	WiFiCallQuality    float64 `json:"wifi_call_quality"` // 0-100
	PreferWiFiCalling  bool    `json:"prefer_wifi_calling"`

	// Advanced Features
	MeshNetworkSupport bool    `json:"mesh_network_support"`
	MeshNodesConnected int     `json:"mesh_nodes_connected"`
	WiFi6Support       bool    `json:"wifi6_support"`
	WiFi6ESupport      bool    `json:"wifi6e_support"`
	WiFi6Usage         float64 `json:"wifi6_usage"` // percentage

	// Guest Network Usage
	GuestNetworkCreated   bool    `json:"guest_network_created"`
	GuestDevicesConnected int     `json:"guest_devices_connected"`
	GuestDataShared       float64 `json:"guest_data_shared"` // GB

	// WiFi Sharing
	WiFiPasswordShared int  `json:"wifi_password_shared"`
	QRCodeSharing      bool `json:"qr_code_sharing"`
	NearbySharing      bool `json:"nearby_sharing"`

	// Performance Metrics
	AverageSpeed        float64 `json:"average_speed"`        // Mbps
	ConnectionStability float64 `json:"connection_stability"` // percentage
	PacketLoss          float64 `json:"packet_loss"`          // percentage

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceBluetoothProfile tracks Bluetooth connections and usage
type DeviceBluetoothProfile struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Paired Devices
	TotalPairedDevices  int        `json:"total_paired_devices"`
	ActivePairedDevices int        `json:"active_paired_devices"`
	PairedDeviceList    string     `gorm:"type:json" json:"paired_device_list"` // JSON array
	TrustedDevices      string     `gorm:"type:json" json:"trusted_devices"`    // JSON array
	LastPairingDate     *time.Time `json:"last_pairing_date"`

	// Connection Stability
	ConnectionSuccess     float64 `json:"connection_success"` // percentage
	ConnectionDrops       int     `json:"connection_drops"`
	ReconnectAttempts     int     `json:"reconnect_attempts"`
	AverageConnectionTime int     `json:"average_connection_time"` // seconds

	// Audio Devices
	AudioDevicesCount   int    `json:"audio_devices_count"`
	HeadphonesConnected string `gorm:"type:json" json:"headphones_connected"` // JSON array
	SpeakersConnected   string `gorm:"type:json" json:"speakers_connected"`   // JSON array
	AudioCodecsUsed     string `gorm:"type:json" json:"audio_codecs_used"`    // JSON array

	// Automotive Connectivity
	CarSystemConnected    bool    `json:"car_system_connected"`
	CarMake               string  `json:"car_make"`
	CarModel              string  `json:"car_model"`
	AutoConnectCar        bool    `json:"auto_connect_car"`
	DrivingHoursConnected float64 `json:"driving_hours_connected"`

	// Wearables
	WearableCount        int  `json:"wearable_count"`
	SmartWatchConnected  bool `json:"smartwatch_connected"`
	FitnessBandConnected bool `json:"fitness_band_connected"`
	HealthDataSync       bool `json:"health_data_sync"`

	// Smart Home
	SmartHomeDevices     int    `json:"smart_home_devices"`
	SmartLocks           int    `json:"smart_locks"`
	SmartAppliances      string `gorm:"type:json" json:"smart_appliances"` // JSON array
	HomeAutomationActive bool   `json:"home_automation_active"`

	// Bluetooth Version
	BluetoothVersion string `json:"bluetooth_version"` // 5.3, 5.2, 5.1, 5.0, 4.2
	LESupport        bool   `json:"le_support"`        // Low Energy
	EDRSupport       bool   `json:"edr_support"`       // Enhanced Data Rate
	APTXSupport      bool   `json:"aptx_support"`

	// Connection Range
	AverageRange float64 `json:"average_range"` // meters
	MaxRange     float64 `json:"max_range"`     // meters
	RangeIssues  int     `json:"range_issues"`

	// Security Settings
	Visibility        string `json:"visibility"` // always, contacts, never
	RequiresPairing   bool   `json:"requires_pairing"`
	AutoAcceptFiles   bool   `json:"auto_accept_files"`
	SecurityIncidents int    `json:"security_incidents"`

	// Power Impact
	BluetoothOnTime float64 `json:"bluetooth_on_time"` // hours per day
	BatteryDrain    float64 `json:"battery_drain"`     // percentage
	AlwaysOn        bool    `json:"always_on"`
	PowerOptimized  bool    `json:"power_optimized"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// DeviceHotspotUsage tracks mobile hotspot usage patterns
type DeviceHotspotUsage struct {
	database.BaseModel
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Activation Frequency
	TotalActivations     int        `json:"total_activations"`
	MonthlyActivations   int        `json:"monthly_activations"`
	LastActivation       *time.Time `json:"last_activation"`
	AverageSessionLength int        `json:"average_session_length"` // minutes

	// Data Sharing
	TotalDataShared       float64 `json:"total_data_shared"`        // GB
	MonthlyDataShared     float64 `json:"monthly_data_shared"`      // GB
	AverageDataPerSession float64 `json:"average_data_per_session"` // GB
	PeakDataShared        float64 `json:"peak_data_shared"`         // GB in single session

	// Connected Devices
	MaxDevicesConnected      int    `json:"max_devices_connected"`
	AverageDevicesPerSession int    `json:"average_devices_per_session"`
	TotalUniqueDevices       int    `json:"total_unique_devices"`
	FrequentDevices          string `gorm:"type:json" json:"frequent_devices"` // JSON array

	// Security Settings
	PasswordProtected    bool   `json:"password_protected"`
	SecurityProtocol     string `json:"security_protocol"` // WPA3, WPA2, Open
	PasswordStrength     string `json:"password_strength"` // strong, medium, weak
	PasswordChanges      int    `json:"password_changes"`
	UnauthorizedAttempts int    `json:"unauthorized_attempts"`

	// Usage Duration
	TotalUsageHours float64 `json:"total_usage_hours"`
	LongestSession  int     `json:"longest_session"`  // minutes
	ShortestSession int     `json:"shortest_session"` // minutes
	NightTimeUsage  float64 `json:"night_time_usage"` // hours

	// Battery Impact
	BatteryDrainRate        float64 `json:"battery_drain_rate"` // % per hour
	SessionsEndedLowBattery int     `json:"sessions_ended_low_battery"`
	ChargingWhileHotspot    int     `json:"charging_while_hotspot"`
	PowerBankUsage          bool    `json:"power_bank_usage"`

	// Bandwidth Management
	BandwidthLimit    float64 `json:"bandwidth_limit"` // Mbps
	ThrottlingEnabled bool    `json:"throttling_enabled"`
	QoSEnabled        bool    `json:"qos_enabled"`
	PriorityDevices   string  `gorm:"type:json" json:"priority_devices"` // JSON array

	// Device Types Connected
	LaptopsConnected      int `json:"laptops_connected"`
	TabletsConnected      int `json:"tablets_connected"`
	PhonesConnected       int `json:"phones_connected"`
	GameConsolesConnected int `json:"game_consoles_connected"`
	OtherDevicesConnected int `json:"other_devices_connected"`

	// Location Usage
	HomeUsage        float64 `json:"home_usage"`         // percentage
	WorkUsage        float64 `json:"work_usage"`         // percentage
	TravelUsage      float64 `json:"travel_usage"`       // percentage
	PublicPlaceUsage float64 `json:"public_place_usage"` // percentage

	// Cost Implications
	DataOverageDueToHotspot bool    `json:"data_overage_due_to_hotspot"`
	AdditionalCosts         float64 `json:"additional_costs"`
	HotspotPlanActive       bool    `json:"hotspot_plan_active"`
	HotspotDataAllowance    float64 `json:"hotspot_data_allowance"` // GB

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// Methods for DeviceNetworkProfile
func (dnp *DeviceNetworkProfile) IsLoyal() bool {
	return dnp.CarrierLoyaltyMonths >= 24 && dnp.CarrierChangeCount <= 1
}

func (dnp *DeviceNetworkProfile) Has5GCapability() bool {
	return dnp.Network5GEnabled && dnp.Network5GUsage > 0
}

func (dnp *DeviceNetworkProfile) IsInternationalUser() bool {
	return dnp.InternationalUsage || dnp.InternationalRoamingDays > 0
}

func (dnp *DeviceNetworkProfile) GetSignalQualityScore() float64 {
	score := 100.0
	if dnp.SignalQuality == "poor" {
		score = 25.0
	} else if dnp.SignalQuality == "fair" {
		score = 50.0
	} else if dnp.SignalQuality == "good" {
		score = 75.0
	}
	// Adjust for reliability
	return score * (dnp.NetworkReliability / 100)
}

func (dnp *DeviceNetworkProfile) HasMultiSIM() bool {
	return dnp.DualSIMEnabled || dnp.eSIMEnabled || dnp.ActiveSIMCount > 1
}

// Methods for DeviceDataUsage
func (ddu *DeviceDataUsage) IsOverLimit() bool {
	return ddu.DataLimit > 0 && ddu.TotalDataUsed > ddu.DataLimit
}

func (ddu *DeviceDataUsage) GetUsagePercentage() float64 {
	if ddu.DataLimit > 0 {
		return (ddu.TotalDataUsed / ddu.DataLimit) * 100
	}
	return 0
}

func (ddu *DeviceDataUsage) IsHeavyUser() bool {
	return ddu.TotalDataUsed > 50 // More than 50GB per month
}

func (ddu *DeviceDataUsage) PrefersWiFi() bool {
	return ddu.WiFiToCellularRatio > 2.0 // WiFi usage is twice cellular
}

func (ddu *DeviceDataUsage) HasOverageRisk() bool {
	remainingDays := 30 - ddu.DaysInCycle
	if remainingDays > 0 {
		dailyRate := ddu.TotalDataUsed / float64(ddu.DaysInCycle)
		projectedUsage := ddu.TotalDataUsed + (dailyRate * float64(remainingDays))
		return projectedUsage > ddu.DataLimit
	}
	return false
}

// Methods for DeviceConnectivityIssues
func (dci *DeviceConnectivityIssues) HasCriticalIssues() bool {
	return dci.OngoingIssues > 0 || dci.HardwareDefect || dci.AntennaIssues
}

func (dci *DeviceConnectivityIssues) GetIssueResolutionRate() float64 {
	if dci.TotalIssues > 0 {
		return float64(dci.ResolvedIssues) / float64(dci.TotalIssues) * 100
	}
	return 100
}

func (dci *DeviceConnectivityIssues) HasSIMProblems() bool {
	return dci.SIMCardErrors > 0 || dci.ESIMActivationFailed || dci.SIMCardReplaced
}

func (dci *DeviceConnectivityIssues) IsHighFrequency() bool {
	return dci.ConnectionDrops > 10 || dci.DropsPerDay > 5
}

func (dci *DeviceConnectivityIssues) NeedsSupport() bool {
	return dci.OngoingIssues > 3 || dci.UnresolvedDuration > 7 || dci.HardwareDefect
}

// Methods for DeviceWiFiProfile
func (dwp *DeviceWiFiProfile) UsesPublicWiFi() bool {
	return dwp.PublicWiFiConnections > 0 || dwp.OpenNetworkUsage > 0
}

func (dwp *DeviceWiFiProfile) IsSecurityConscious() bool {
	return dwp.WPA3Enabled && dwp.VPNAutoConnect && dwp.UnsafeConnections == 0
}

func (dwp *DeviceWiFiProfile) HasAdvancedFeatures() bool {
	return dwp.WiFi6Support || dwp.WiFi6ESupport || dwp.MeshNetworkSupport
}

func (dwp *DeviceWiFiProfile) GetNetworkStability() float64 {
	totalStability := dwp.HomeNetworkStability + dwp.WorkNetworkStability
	if totalStability > 0 {
		return totalStability / 2
	}
	return dwp.ConnectionStability
}

func (dwp *DeviceWiFiProfile) SharesWiFi() bool {
	return dwp.WiFiPasswordShared > 0 || dwp.QRCodeSharing || dwp.NearbySharing
}

// Methods for DeviceBluetoothProfile
func (dbp *DeviceBluetoothProfile) IsActiveUser() bool {
	return dbp.ActivePairedDevices > 3 || dbp.BluetoothOnTime > 12
}

func (dbp *DeviceBluetoothProfile) HasCarIntegration() bool {
	return dbp.CarSystemConnected && dbp.AutoConnectCar
}

func (dbp *DeviceBluetoothProfile) HasWearables() bool {
	return dbp.WearableCount > 0 || dbp.SmartWatchConnected || dbp.FitnessBandConnected
}

func (dbp *DeviceBluetoothProfile) GetConnectionReliability() float64 {
	if dbp.ConnectionSuccess > 0 {
		return dbp.ConnectionSuccess
	}
	// Calculate from drops and reconnects
	totalAttempts := float64(dbp.ActivePairedDevices * 30) // Assume 30 connections per device per month
	failures := float64(dbp.ConnectionDrops + dbp.ReconnectAttempts)
	if totalAttempts > 0 {
		return ((totalAttempts - failures) / totalAttempts) * 100
	}
	return 100
}

func (dbp *DeviceBluetoothProfile) IsSmartHomeUser() bool {
	return dbp.SmartHomeDevices > 0 || dbp.SmartLocks > 0 || dbp.HomeAutomationActive
}

// Methods for DeviceHotspotUsage
func (dhu *DeviceHotspotUsage) IsFrequentUser() bool {
	return dhu.MonthlyActivations > 10 || dhu.MonthlyDataShared > 10
}

func (dhu *DeviceHotspotUsage) IsSecure() bool {
	return dhu.PasswordProtected && (dhu.SecurityProtocol == "WPA3" || dhu.SecurityProtocol == "WPA2") &&
		dhu.PasswordStrength == "strong"
}

func (dhu *DeviceHotspotUsage) GetAverageBatteryImpact() float64 {
	if dhu.TotalUsageHours > 0 {
		return dhu.BatteryDrainRate * (dhu.TotalUsageHours / float64(dhu.TotalActivations))
	}
	return 0
}

func (dhu *DeviceHotspotUsage) CausesDataOverage() bool {
	return dhu.DataOverageDueToHotspot ||
		(dhu.HotspotDataAllowance > 0 && dhu.MonthlyDataShared > dhu.HotspotDataAllowance)
}

func (dhu *DeviceHotspotUsage) GetPrimaryUsageLocation() string {
	maxUsage := dhu.HomeUsage
	location := "home"

	if dhu.WorkUsage > maxUsage {
		maxUsage = dhu.WorkUsage
		location = "work"
	}
	if dhu.TravelUsage > maxUsage {
		maxUsage = dhu.TravelUsage
		location = "travel"
	}
	if dhu.PublicPlaceUsage > maxUsage {
		location = "public"
	}

	return location
}
