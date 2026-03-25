package smartphone

import (
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// SmartphoneNetwork tracks smartphone-specific network capabilities and metrics
type SmartphoneNetwork struct {
	database.BaseModel
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	RecordDate time.Time `json:"record_date"`

	// Cellular Capabilities
	Supports5G       bool   `gorm:"default:false" json:"supports_5g"`
	Supports5GSA     bool   `gorm:"default:false" json:"supports_5g_sa"`     // Standalone 5G
	Supports5GNSA    bool   `gorm:"default:false" json:"supports_5g_nsa"`    // Non-standalone 5G
	Supports5GmmWave bool   `gorm:"default:false" json:"supports_5g_mmwave"` // Millimeter wave
	FiveGBands       string `gorm:"type:json" json:"five_g_bands"`           // JSON array of supported bands

	SupportsLTE    bool   `gorm:"default:true" json:"supports_lte"`
	LTECategory    int    `json:"lte_category"`                         // CAT 4, 6, 12, 16, 20, etc
	LTEBands       string `gorm:"type:json" json:"lte_bands"`           // JSON array
	SupportsVoLTE  bool   `gorm:"default:false" json:"supports_volte"`  // Voice over LTE
	SupportsVoWiFi bool   `gorm:"default:false" json:"supports_vowifi"` // Voice over WiFi
	SupportsViLTE  bool   `gorm:"default:false" json:"supports_vilte"`  // Video over LTE

	// SIM Configuration
	SIMType            string `json:"sim_type"` // physical, eSIM, hybrid
	DualSIMSupport     bool   `gorm:"default:false" json:"dual_sim_support"`
	DualSIMType        string `json:"dual_sim_type"` // dual_physical, physical_esim
	TripleSIMSupport   bool   `gorm:"default:false" json:"triple_sim_support"`
	ESIMSupport        bool   `gorm:"default:false" json:"esim_support"`
	ActiveSIMCount     int    `json:"active_sim_count"`
	DualSIMDualStandby bool   `gorm:"default:false" json:"dual_sim_dual_standby"` // DSDS
	DualSIMDualActive  bool   `gorm:"default:false" json:"dual_sim_dual_active"`  // DSDA

	// Carrier Information
	PrimaryCarrier     string `json:"primary_carrier"`
	SecondaryCarrier   string `json:"secondary_carrier"`
	CarrierAggregation bool   `gorm:"default:false" json:"carrier_aggregation"` // LTE-A
	MaxDownloadSpeed   int    `json:"max_download_speed"`                       // Mbps theoretical
	MaxUploadSpeed     int    `json:"max_upload_speed"`                         // Mbps theoretical

	// Current Network Status
	CurrentNetworkType string  `json:"current_network_type"` // 5G, LTE, 3G, 2G
	CurrentBand        string  `json:"current_band"`
	SignalStrength     int     `json:"signal_strength"` // dBm
	SignalQuality      string  `json:"signal_quality"`  // excellent, good, fair, poor
	NetworkLatency     float64 `json:"network_latency"` // ms
	PacketLoss         float64 `json:"packet_loss"`     // percentage

	// Data Usage Metrics
	CellularDataUsed float64 `json:"cellular_data_used"` // GB this period
	DataLimit        float64 `json:"data_limit"`         // GB
	DataPlanType     string  `json:"data_plan_type"`     // unlimited, limited, prepaid
	RoamingEnabled   bool    `gorm:"default:false" json:"roaming_enabled"`
	RoamingDataUsed  float64 `json:"roaming_data_used"` // GB
	DataSaverEnabled bool    `gorm:"default:false" json:"data_saver_enabled"`

	// WiFi Capabilities
	WiFiGeneration string `json:"wifi_generation"`                      // WiFi 6E, WiFi 6, WiFi 5
	WiFi6ESupport  bool   `gorm:"default:false" json:"wifi_6e_support"` // 6GHz band
	WiFi6Support   bool   `gorm:"default:false" json:"wifi_6_support"`  // 802.11ax
	WiFi5Support   bool   `gorm:"default:true" json:"wifi_5_support"`   // 802.11ac
	DualBandWiFi   bool   `gorm:"default:true" json:"dual_band_wifi"`   // 2.4GHz + 5GHz
	MIMOSupport    string `json:"mimo_support"`                         // 2x2, 4x4, 8x8

	// Current WiFi Status
	WiFiConnected      bool   `gorm:"default:false" json:"wifi_connected"`
	WiFiSSID           string `json:"wifi_ssid"`
	WiFiFrequency      string `json:"wifi_frequency"` // 2.4GHz, 5GHz, 6GHz
	WiFiChannel        int    `json:"wifi_channel"`
	WiFiLinkSpeed      int    `json:"wifi_link_speed"`      // Mbps
	WiFiSignalStrength int    `json:"wifi_signal_strength"` // dBm
	WiFiSecurity       string `json:"wifi_security"`        // WPA3, WPA2, Open

	// Bluetooth Capabilities
	BluetoothVersion string `json:"bluetooth_version"`                // 5.3, 5.2, 5.1, 5.0
	BluetoothLE      bool   `gorm:"default:true" json:"bluetooth_le"` // Low Energy
	BluetoothAptX    bool   `gorm:"default:false" json:"bluetooth_aptx"`
	BluetoothLDAC    bool   `gorm:"default:false" json:"bluetooth_ldac"`
	BluetoothRange   int    `json:"bluetooth_range"` // meters
	ConnectedDevices int    `json:"connected_devices"`
	PairedDevices    string `gorm:"type:json" json:"paired_devices"` // JSON array

	// NFC Capabilities
	NFCSupport        bool `gorm:"default:false" json:"nfc_support"`
	NFCEnabled        bool `gorm:"default:false" json:"nfc_enabled"`
	NFCPaymentEnabled bool `gorm:"default:false" json:"nfc_payment_enabled"`
	NFCTagsRead       int  `json:"nfc_tags_read"`

	// Advanced Features
	WiFiDirect         bool    `gorm:"default:false" json:"wifi_direct"`
	WiFiHotspotEnabled bool    `gorm:"default:false" json:"wifi_hotspot_enabled"`
	HotspotClients     int     `json:"hotspot_clients"`
	HotspotDataUsed    float64 `json:"hotspot_data_used"` // GB
	USBTethering       bool    `gorm:"default:false" json:"usb_tethering"`
	BluetoothTethering bool    `gorm:"default:false" json:"bluetooth_tethering"`

	// Network Performance
	DownloadSpeed      float64 `json:"download_speed"`    // Mbps actual
	UploadSpeed        float64 `json:"upload_speed"`      // Mbps actual
	Jitter             float64 `json:"jitter"`            // ms
	DNSResponseTime    float64 `json:"dns_response_time"` // ms
	TCPRetransmissions int     `json:"tcp_retransmissions"`

	// Network Security
	VPNActive       bool   `gorm:"default:false" json:"vpn_active"`
	VPNProtocol     string `json:"vpn_protocol"` // OpenVPN, WireGuard, IKEv2
	PrivateDNS      bool   `gorm:"default:false" json:"private_dns"`
	DNSOverHTTPS    bool   `gorm:"default:false" json:"dns_over_https"`
	FirewallEnabled bool   `gorm:"default:false" json:"firewall_enabled"`

	// Network Quality Score
	NetworkQualityScore float64 `json:"network_quality_score"` // 0-100
	ConnectivityScore   float64 `json:"connectivity_score"`    // 0-100
	ReliabilityScore    float64 `json:"reliability_score"`     // 0-100
}

// Business logic methods

// Is5GCapable checks if device supports any 5G
func (sn *SmartphoneNetwork) Is5GCapable() bool {
	return sn.Supports5G || sn.Supports5GSA || sn.Supports5GNSA
}

// Is5GUltraCapable checks if device supports mmWave 5G
func (sn *SmartphoneNetwork) Is5GUltraCapable() bool {
	return sn.Supports5GmmWave
}

// IsDualSIMActive checks if dual SIM is actively used
func (sn *SmartphoneNetwork) IsDualSIMActive() bool {
	return sn.DualSIMSupport && sn.ActiveSIMCount >= 2
}

// HasESIM checks if device has eSIM capability
func (sn *SmartphoneNetwork) HasESIM() bool {
	return sn.ESIMSupport || sn.DualSIMType == "physical_esim"
}

// IsWiFi6Ready checks for WiFi 6 or 6E support
func (sn *SmartphoneNetwork) IsWiFi6Ready() bool {
	return sn.WiFi6Support || sn.WiFi6ESupport
}

// GetNetworkGeneration returns the best network generation supported
func (sn *SmartphoneNetwork) GetNetworkGeneration() string {
	if sn.Is5GCapable() {
		if sn.Supports5GmmWave {
			return "5G Ultra"
		}
		return "5G"
	}
	if sn.SupportsLTE {
		return "4G LTE"
	}
	return "3G"
}

// HasHighSpeedConnectivity checks for high-speed connectivity
func (sn *SmartphoneNetwork) HasHighSpeedConnectivity() bool {
	// Check cellular
	if sn.CurrentNetworkType == "5G" && sn.DownloadSpeed > 100 {
		return true
	}
	// Check WiFi
	if sn.WiFiConnected && sn.WiFiLinkSpeed > 100 {
		return true
	}
	return false
}

// IsRoaming checks if device is currently roaming
func (sn *SmartphoneNetwork) IsRoaming() bool {
	return sn.RoamingEnabled && sn.RoamingDataUsed > 0
}

// GetSignalQuality categorizes signal strength
func (sn *SmartphoneNetwork) GetSignalQuality() string {
	// For cellular (dBm)
	if !sn.WiFiConnected {
		switch {
		case sn.SignalStrength >= -70:
			return "excellent"
		case sn.SignalStrength >= -85:
			return "good"
		case sn.SignalStrength >= -100:
			return "fair"
		default:
			return "poor"
		}
	}

	// For WiFi (dBm)
	switch {
	case sn.WiFiSignalStrength >= -50:
		return "excellent"
	case sn.WiFiSignalStrength >= -60:
		return "good"
	case sn.WiFiSignalStrength >= -70:
		return "fair"
	default:
		return "poor"
	}
}

// HasDataLimitRisk checks if approaching data limit
func (sn *SmartphoneNetwork) HasDataLimitRisk() bool {
	if sn.DataPlanType == "unlimited" {
		return false
	}
	if sn.DataLimit > 0 {
		usagePercentage := (sn.CellularDataUsed / sn.DataLimit) * 100
		return usagePercentage > 80
	}
	return false
}

// GetDataUsageRate calculates data usage rate
func (sn *SmartphoneNetwork) GetDataUsageRate() float64 {
	// GB per day
	daysSinceRecord := time.Since(sn.RecordDate).Hours() / 24
	if daysSinceRecord > 0 {
		return sn.CellularDataUsed / daysSinceRecord
	}
	return 0
}

// IsHotspotActive checks if hotspot is active
func (sn *SmartphoneNetwork) IsHotspotActive() bool {
	return sn.WiFiHotspotEnabled && sn.HotspotClients > 0
}

// HasSecureConnection checks connection security
func (sn *SmartphoneNetwork) HasSecureConnection() bool {
	// Check VPN
	if sn.VPNActive {
		return true
	}
	// Check WiFi security
	if sn.WiFiConnected && sn.WiFiSecurity == "WPA3" {
		return true
	}
	// Check DNS security
	if sn.PrivateDNS || sn.DNSOverHTTPS {
		return true
	}
	return false
}

// GetBluetoothGeneration returns Bluetooth generation
func (sn *SmartphoneNetwork) GetBluetoothGeneration() float64 {
	switch sn.BluetoothVersion {
	case "5.3":
		return 5.3
	case "5.2":
		return 5.2
	case "5.1":
		return 5.1
	case "5.0":
		return 5.0
	case "4.2":
		return 4.2
	default:
		return 4.0
	}
}

// HasHighQualityAudioCodec checks for high-quality audio codecs
func (sn *SmartphoneNetwork) HasHighQualityAudioCodec() bool {
	return sn.BluetoothAptX || sn.BluetoothLDAC
}

// CalculateNetworkScore calculates overall network score
func (sn *SmartphoneNetwork) CalculateNetworkScore() float64 {
	score := 50.0 // Base score

	// Network generation scoring
	if sn.Is5GCapable() {
		score += 20
		if sn.Is5GUltraCapable() {
			score += 5
		}
	} else if sn.SupportsLTE {
		score += 10
	}

	// WiFi scoring
	if sn.WiFi6ESupport {
		score += 10
	} else if sn.WiFi6Support {
		score += 7
	} else if sn.WiFi5Support {
		score += 5
	}

	// Speed scoring
	if sn.DownloadSpeed > 100 {
		score += 10
	} else if sn.DownloadSpeed > 50 {
		score += 5
	}

	// Feature scoring
	if sn.HasESIM() {
		score += 3
	}
	if sn.NFCSupport {
		score += 2
	}
	if sn.DualSIMSupport {
		score += 3
	}

	// Quality adjustments
	signalQuality := sn.GetSignalQuality()
	switch signalQuality {
	case "excellent":
		score += 5
	case "good":
		score += 3
	case "fair":
		score += 1
	case "poor":
		score -= 5
	}

	// Latency penalty
	if sn.NetworkLatency > 100 {
		score -= 10
	} else if sn.NetworkLatency > 50 {
		score -= 5
	}

	// Packet loss penalty
	if sn.PacketLoss > 5 {
		score -= 10
	} else if sn.PacketLoss > 2 {
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

// NeedsNetworkOptimization checks if network needs optimization
func (sn *SmartphoneNetwork) NeedsNetworkOptimization() bool {
	// Poor signal quality
	if sn.GetSignalQuality() == "poor" {
		return true
	}
	// High packet loss
	if sn.PacketLoss > 5 {
		return true
	}
	// High latency
	if sn.NetworkLatency > 100 {
		return true
	}
	// Low speeds
	if sn.DownloadSpeed < 10 {
		return true
	}
	// Poor network score
	if sn.CalculateNetworkScore() < 50 {
		return true
	}
	return false
}

// IsDataIntensive checks if usage pattern is data-intensive
func (sn *SmartphoneNetwork) IsDataIntensive() bool {
	// More than 1GB per day
	if sn.GetDataUsageRate() > 1.0 {
		return true
	}
	// Heavy hotspot usage
	if sn.HotspotDataUsed > 10 {
		return true
	}
	// Total usage high
	if sn.CellularDataUsed > 50 {
		return true
	}
	return false
}
