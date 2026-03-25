package smartwatch

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartwatchNetwork tracks smartwatch-specific network capabilities and metrics
type SmartwatchNetwork struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordDate time.Time `json:"record_date"`

	// Connectivity Model
	ConnectivityType  string `json:"connectivity_type"` // gps_only, gps_cellular, wifi_only
	IsCellularModel   bool   `gorm:"default:false" json:"is_cellular_model"`
	IsGPSModel        bool   `gorm:"default:true" json:"is_gps_model"`
	RequiresPhone     bool   `gorm:"default:true" json:"requires_phone"`
	StandaloneCapable bool   `gorm:"default:false" json:"standalone_capable"`

	// Cellular Capabilities (for LTE models)
	SupportsLTE     bool    `gorm:"default:false" json:"supports_lte"`
	SupportsUMTS    bool    `gorm:"default:false" json:"supports_umts"`
	ESIMEnabled     bool    `gorm:"default:false" json:"esim_enabled"`
	PhysicalSIM     bool    `gorm:"default:false" json:"physical_sim"`
	CellularCarrier string  `json:"cellular_carrier"`
	CellularPlan    string  `json:"cellular_plan"` // standalone, numbershare, family
	DataLimit       float64 `json:"data_limit"`    // GB
	DataUsed        float64 `json:"data_used"`     // GB this period
	RoamingEnabled  bool    `gorm:"default:false" json:"roaming_enabled"`

	// Current Cellular Status (if applicable)
	CellularActive    bool    `gorm:"default:false" json:"cellular_active"`
	SignalStrength    int     `json:"signal_strength"`     // dBm
	SignalBars        int     `json:"signal_bars"`         // 0-5
	NetworkType       string  `json:"network_type"`        // LTE, UMTS, none
	CellularDataSpeed float64 `json:"cellular_data_speed"` // Mbps

	// Bluetooth (Primary Connection)
	BluetoothVersion string `json:"bluetooth_version"` // 5.0, 5.2, 5.3
	BluetoothLE      bool   `gorm:"default:true" json:"bluetooth_le"`
	BluetoothClassic bool   `gorm:"default:false" json:"bluetooth_classic"`
	BluetoothRange   int    `json:"bluetooth_range"` // meters

	// Phone Connection Status
	PhoneConnected     bool       `gorm:"default:false" json:"phone_connected"`
	PhoneName          string     `json:"phone_name"`
	PhoneModel         string     `json:"phone_model"`
	PhoneOS            string     `json:"phone_os"`            // iOS, Android
	ConnectionQuality  string     `json:"connection_quality"`  // excellent, good, poor
	ConnectionDistance float64    `json:"connection_distance"` // meters
	LastDisconnect     *time.Time `json:"last_disconnect"`
	DisconnectCount    int        `json:"disconnect_count"` // daily count

	// WiFi Capabilities
	WiFiSupport    bool   `gorm:"default:true" json:"wifi_support"`
	WiFiGeneration string `json:"wifi_generation"` // WiFi 4, WiFi 5
	WiFiBands      string `json:"wifi_bands"`      // 2.4GHz, 5GHz, both

	// Current WiFi Status
	WiFiConnected      bool   `gorm:"default:false" json:"wifi_connected"`
	WiFiSSID           string `json:"wifi_ssid"`
	WiFiSignalStrength int    `json:"wifi_signal_strength"` // dBm
	WiFiSpeed          int    `json:"wifi_speed"`           // Mbps
	WiFiSecurity       string `json:"wifi_security"`        // WPA3, WPA2, Open
	KnownNetworks      int    `json:"known_networks"`       // count of saved networks

	// Data Sync Performance
	SyncEnabled      bool       `gorm:"default:true" json:"sync_enabled"`
	SyncFrequency    string     `json:"sync_frequency"` // continuous, hourly, daily
	LastSyncTime     *time.Time `json:"last_sync_time"`
	SyncLatency      float64    `json:"sync_latency"`     // seconds
	SyncDataVolume   float64    `json:"sync_data_volume"` // MB per day
	PendingSyncItems int        `json:"pending_sync_items"`
	SyncErrors       int        `json:"sync_errors"`

	// Notification Sync
	NotificationSync        bool    `gorm:"default:true" json:"notification_sync"`
	NotificationLatency     float64 `json:"notification_latency"` // seconds
	MissedNotifications     int     `json:"missed_notifications"`
	NotificationReliability float64 `json:"notification_reliability"` // percentage

	// GPS Performance
	GPSEnabled  bool    `gorm:"default:true" json:"gps_enabled"`
	GPSAccuracy float64 `json:"gps_accuracy"` // meters
	GPSFixTime  float64 `json:"gps_fix_time"` // seconds
	AssistedGPS bool    `gorm:"default:true" json:"assisted_gps"`
	GLONASS     bool    `gorm:"default:false" json:"glonass"`
	Galileo     bool    `gorm:"default:false" json:"galileo"`
	BeiDou      bool    `gorm:"default:false" json:"beidou"`

	// NFC Capabilities
	NFCSupport         bool `gorm:"default:false" json:"nfc_support"`
	NFCEnabled         bool `gorm:"default:false" json:"nfc_enabled"`
	ApplePayEnabled    bool `gorm:"default:false" json:"apple_pay_enabled"`
	GooglePayEnabled   bool `gorm:"default:false" json:"google_pay_enabled"`
	TransitCardSupport bool `gorm:"default:false" json:"transit_card_support"`
	NFCTransactions    int  `json:"nfc_transactions"`

	// Music & Media Streaming
	MusicStreaming      bool    `gorm:"default:false" json:"music_streaming"`
	StreamingService    string  `json:"streaming_service"`     // Apple Music, Spotify, etc
	OfflineMusicStorage int     `json:"offline_music_storage"` // MB
	StreamingQuality    string  `json:"streaming_quality"`     // high, medium, low
	StreamingDataUsed   float64 `json:"streaming_data_used"`   // MB

	// App Data Usage
	AppDataUsage          float64 `json:"app_data_usage"` // MB per day
	BackgroundDataEnabled bool    `gorm:"default:true" json:"background_data_enabled"`
	DataSaverMode         bool    `gorm:"default:false" json:"data_saver_mode"`

	// Emergency Services
	EmergencySOSEnabled    bool `gorm:"default:false" json:"emergency_sos_enabled"`
	InternationalEmergency bool `gorm:"default:false" json:"international_emergency"`
	EmergencyContacts      int  `json:"emergency_contacts"`
	LocationSharing        bool `gorm:"default:false" json:"location_sharing"`

	// Power Management
	AirplaneModeEnabled bool `gorm:"default:false" json:"airplane_mode_enabled"`
	PowerReserveMode    bool `gorm:"default:false" json:"power_reserve_mode"`
	LowPowerMode        bool `gorm:"default:false" json:"low_power_mode"`

	// Network Quality Metrics
	PacketLoss         float64 `json:"packet_loss"`         // percentage
	Latency            float64 `json:"latency"`             // ms
	Jitter             float64 `json:"jitter"`              // ms
	NetworkReliability float64 `json:"network_reliability"` // percentage uptime

	// Family Setup (for kids/elderly)
	FamilySetupEnabled      bool `gorm:"default:false" json:"family_setup_enabled"`
	SchoolTimeEnabled       bool `gorm:"default:false" json:"schooltime_enabled"`
	LocationTrackingEnabled bool `gorm:"default:false" json:"location_tracking_enabled"`
	ContactsRestricted      bool `gorm:"default:false" json:"contacts_restricted"`

	// Network Score
	NetworkScore      float64 `json:"network_score"`      // 0-100
	ConnectivityScore float64 `json:"connectivity_score"` // 0-100
}

// Business logic methods

// IsCellularCapable checks if watch has cellular
func (sn *SmartwatchNetwork) IsCellularCapable() bool {
	return sn.IsCellularModel || sn.SupportsLTE
}

// IsStandalone checks if can work without phone
func (sn *SmartwatchNetwork) IsStandalone() bool {
	return sn.StandaloneCapable || (sn.IsCellularCapable() && sn.WiFiSupport)
}

// HasActiveConnectivity checks any active connection
func (sn *SmartwatchNetwork) HasActiveConnectivity() bool {
	return sn.PhoneConnected || sn.WiFiConnected || sn.CellularActive
}

// GetPrimaryConnection returns main connection type
func (sn *SmartwatchNetwork) GetPrimaryConnection() string {
	if sn.PhoneConnected {
		return "bluetooth"
	}
	if sn.WiFiConnected {
		return "wifi"
	}
	if sn.CellularActive {
		return "cellular"
	}
	return "none"
}

// IsPhoneDependent checks phone dependency
func (sn *SmartwatchNetwork) IsPhoneDependent() bool {
	return sn.RequiresPhone && !sn.StandaloneCapable
}

// GetConnectionQuality rates connection quality
func (sn *SmartwatchNetwork) GetConnectionQuality() string {
	if sn.PhoneConnected {
		// Check Bluetooth connection quality
		if sn.DisconnectCount > 5 {
			return "poor"
		}
		if sn.DisconnectCount > 2 {
			return "fair"
		}
		if sn.ConnectionDistance < 10 {
			return "excellent"
		}
		return "good"
	}

	if sn.WiFiConnected {
		// Check WiFi signal
		if sn.WiFiSignalStrength >= -50 {
			return "excellent"
		}
		if sn.WiFiSignalStrength >= -60 {
			return "good"
		}
		if sn.WiFiSignalStrength >= -70 {
			return "fair"
		}
		return "poor"
	}

	if sn.CellularActive {
		// Check cellular signal
		if sn.SignalBars >= 4 {
			return "excellent"
		}
		if sn.SignalBars >= 3 {
			return "good"
		}
		if sn.SignalBars >= 2 {
			return "fair"
		}
		return "poor"
	}

	return "none"
}

// HasReliableSync checks sync reliability
func (sn *SmartwatchNetwork) HasReliableSync() bool {
	// Check sync errors
	if sn.SyncErrors > 5 {
		return false
	}
	// Check pending items
	if sn.PendingSyncItems > 50 {
		return false
	}
	// Check notification reliability
	if sn.NotificationReliability < 90 {
		return false
	}
	// Check latency
	if sn.SyncLatency > 5 || sn.NotificationLatency > 3 {
		return false
	}
	return true
}

// IsHighDataUser checks data usage
func (sn *SmartwatchNetwork) IsHighDataUser() bool {
	totalData := sn.DataUsed + (sn.AppDataUsage*30)/1000 +
		(sn.StreamingDataUsed*30)/1000
	return totalData > 1 // More than 1GB per month
}

// HasPaymentCapability checks NFC payments
func (sn *SmartwatchNetwork) HasPaymentCapability() bool {
	return sn.NFCSupport && (sn.ApplePayEnabled || sn.GooglePayEnabled || sn.TransitCardSupport)
}

// GetGPSCapability rates GPS features
func (sn *SmartwatchNetwork) GetGPSCapability() string {
	satelliteSystems := 0
	if sn.GPSEnabled {
		satelliteSystems++
	}
	if sn.GLONASS {
		satelliteSystems++
	}
	if sn.Galileo {
		satelliteSystems++
	}
	if sn.BeiDou {
		satelliteSystems++
	}

	if satelliteSystems >= 3 {
		return "advanced"
	}
	if satelliteSystems >= 2 {
		return "enhanced"
	}
	if satelliteSystems >= 1 {
		return "basic"
	}
	return "none"
}

// CanStreamMusic checks music streaming capability
func (sn *SmartwatchNetwork) CanStreamMusic() bool {
	// Needs WiFi or cellular for streaming
	return sn.MusicStreaming && (sn.WiFiSupport || sn.IsCellularCapable())
}

// IsFamilyWatch checks if family setup watch
func (sn *SmartwatchNetwork) IsFamilyWatch() bool {
	return sn.FamilySetupEnabled || sn.SchoolTimeEnabled || sn.ContactsRestricted
}

// GetNetworkScore calculates overall network score
func (sn *SmartwatchNetwork) GetNetworkScore() float64 {
	score := 40.0 // Base score

	// Connection type scoring
	if sn.IsCellularCapable() {
		score += 15
	}
	if sn.WiFiSupport {
		score += 10
	}

	// Connection quality
	quality := sn.GetConnectionQuality()
	switch quality {
	case "excellent":
		score += 15
	case "good":
		score += 10
	case "fair":
		score += 5
	case "poor":
		score -= 5
	}

	// Sync reliability
	if sn.HasReliableSync() {
		score += 10
	}

	// GPS capability
	gpsCapability := sn.GetGPSCapability()
	switch gpsCapability {
	case "advanced":
		score += 10
	case "enhanced":
		score += 7
	case "basic":
		score += 4
	}

	// NFC payments
	if sn.HasPaymentCapability() {
		score += 5
	}

	// Standalone capability
	if sn.IsStandalone() {
		score += 10
	}

	// Penalties
	if sn.DisconnectCount > 10 {
		score -= 10
	}
	if sn.PacketLoss > 5 {
		score -= 5
	}
	if sn.Latency > 100 {
		score -= 5
	}

	// Cap at 100
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

// NeedsNetworkOptimization checks if optimization needed
func (sn *SmartwatchNetwork) NeedsNetworkOptimization() bool {
	// Poor connection quality
	if sn.GetConnectionQuality() == "poor" {
		return true
	}
	// Unreliable sync
	if !sn.HasReliableSync() {
		return true
	}
	// Frequent disconnects
	if sn.DisconnectCount > 10 {
		return true
	}
	// High packet loss
	if sn.PacketLoss > 5 {
		return true
	}
	// Poor network score
	if sn.GetNetworkScore() < 50 {
		return true
	}
	return false
}

// GetDataUsageCategory categorizes data usage
func (sn *SmartwatchNetwork) GetDataUsageCategory() string {
	totalDataMB := (sn.DataUsed * 1000) + sn.AppDataUsage*30 + sn.StreamingDataUsed*30

	if totalDataMB < 100 {
		return "minimal"
	}
	if totalDataMB < 500 {
		return "light"
	}
	if totalDataMB < 1000 {
		return "moderate"
	}
	return "heavy"
}

// IsEmergencyReady checks emergency features
func (sn *SmartwatchNetwork) IsEmergencyReady() bool {
	// Has emergency SOS
	if !sn.EmergencySOSEnabled {
		return false
	}
	// Has connectivity for emergency
	if !sn.HasActiveConnectivity() && !sn.IsCellularCapable() {
		return false
	}
	// Has emergency contacts
	if sn.EmergencyContacts == 0 {
		return false
	}
	return true
}
