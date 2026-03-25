package claim

import (
	"time"
	
	"github.com/google/uuid"
	"gorm.io/datatypes"
	
	"smartsure/pkg/database"
)

// ClaimDeviceDiagnostics - Remote device diagnostics and health check
type ClaimDeviceDiagnostics struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Remote Diagnostics
	DiagnosticRunDate  time.Time `json:"diagnostic_run_date"`
	DeviceIMEI         string    `json:"device_imei"`
	DeviceSerialNumber string    `json:"device_serial_number"`
	OSVersion          string    `json:"os_version"`
	AppVersion         string    `json:"app_version"`

	// Hardware Health
	BatteryHealth      float64        `json:"battery_health"` // Percentage
	BatteryCycles      int            `json:"battery_cycles"`
	ScreenHealth       string         `json:"screen_health"` // good, cracked, dead_pixels, touch_issues
	SpeakerStatus      string         `json:"speaker_status"`
	MicrophoneStatus   string         `json:"microphone_status"`
	CameraStatus       datatypes.JSON `gorm:"type:json" json:"camera_status"` // front/back camera status
	ChargingPortStatus string         `json:"charging_port_status"`
	ButtonsStatus      datatypes.JSON `gorm:"type:json" json:"buttons_status"` // power, volume, home

	// Sensor Status
	AccelerometerWorking      bool `json:"accelerometer_working"`
	GyroscopeWorking          bool `json:"gyroscope_working"`
	ProximitySensorWorking    bool `json:"proximity_sensor_working"`
	AmbientLightSensorWorking bool `json:"ambient_light_sensor_working"`
	FingerprintSensorWorking  bool `json:"fingerprint_sensor_working"`
	FaceIDWorking             bool `json:"face_id_working"`

	// Connectivity Tests
	WiFiWorking      bool `json:"wifi_working"`
	BluetoothWorking bool `json:"bluetooth_working"`
	CellularWorking  bool `json:"cellular_working"`
	GPSWorking       bool `json:"gps_working"`
	NFCWorking       bool `json:"nfc_working"`

	// Performance Metrics
	StorageUsed    int64      `json:"storage_used_gb"`
	StorageTotal   int64      `json:"storage_total_gb"`
	RAMUsage       float64    `json:"ram_usage_percent"`
	CPUTemperature float64    `json:"cpu_temperature"`
	LastBackupDate *time.Time `json:"last_backup_date"`

	// Damage Detection
	WaterDamageIndicator bool `json:"water_damage_indicator"`
	DropDetectedCount    int  `json:"drop_detected_count"`
	ThermalEventCount    int  `json:"thermal_event_count"`

	// Diagnostic Results
	DiagnosticScore      float64        `json:"diagnostic_score"`
	IssuesDetected       datatypes.JSON `gorm:"type:json" json:"issues_detected"`
	RepairRecommendation string         `gorm:"type:text" json:"repair_recommendation"`
	EstimatedRepairTime  int            `json:"estimated_repair_time_hours"`
}

// ClaimRepairNetwork - Integration with authorized repair network
type ClaimRepairNetwork struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Repair Center Selection
	PreferredRepairCenterID *uuid.UUID     `gorm:"type:uuid" json:"preferred_repair_center_id"`
	AssignedRepairCenterID  *uuid.UUID     `gorm:"type:uuid" json:"assigned_repair_center_id"`
	RepairCenterName        string         `json:"repair_center_name"`
	RepairCenterAddress     datatypes.JSON `gorm:"type:json" json:"repair_center_address"`
	RepairCenterRating      float64        `json:"repair_center_rating"`
	Distance                float64        `json:"distance_km"`

	// Repair Booking
	BookingReference string     `json:"booking_reference"`
	BookingDate      *time.Time `json:"booking_date"`
	AppointmentDate  *time.Time `json:"appointment_date"`
	TimeSlot         string     `json:"time_slot"`

	// Service Options
	ServiceType          string     `json:"service_type"` // walk_in, mail_in, home_service, pickup
	ExpressService       bool       `json:"express_service"`
	LoanerDeviceRequired bool       `json:"loaner_device_required"`
	LoanerDeviceID       *uuid.UUID `gorm:"type:uuid" json:"loaner_device_id"`
	DataTransferRequired bool       `json:"data_transfer_required"`

	// Parts Management
	PartsRequired     datatypes.JSON `gorm:"type:json" json:"parts_required"`
	PartsAvailability string         `json:"parts_availability"` // in_stock, ordered, backordered
	GenuinePartsUsed  bool           `json:"genuine_parts_used"`
	PartsWarranty     int            `json:"parts_warranty_months"`

	// Repair Tracking
	RepairStatus         string     `json:"repair_status"`
	RepairStartDate      *time.Time `json:"repair_start_date"`
	RepairCompletionDate *time.Time `json:"repair_completion_date"`
	QualityCheckPassed   *bool      `json:"quality_check_passed"`
	CustomerPickupDate   *time.Time `json:"customer_pickup_date"`

	// Costs
	LaborCost       float64 `json:"labor_cost"`
	PartsCost       float64 `json:"parts_cost"`
	ShippingCost    float64 `json:"shipping_cost"`
	TotalRepairCost float64 `json:"total_repair_cost"`
	CustomerCopay   float64 `json:"customer_copay"`
}

// ClaimReplacementDevice - Device replacement management
type ClaimReplacementDevice struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Replacement Options
	ReplacementType        string         `json:"replacement_type"` // new, refurbished, like_model, upgrade
	OriginalDeviceModel    string         `json:"original_device_model"`
	ReplacementDeviceModel string         `json:"replacement_device_model"`
	ColorPreference        datatypes.JSON `gorm:"type:json" json:"color_preference"`
	StoragePreference      int            `json:"storage_preference_gb"`

	// Replacement Device Details
	ReplacementDeviceID     *uuid.UUID `gorm:"type:uuid" json:"replacement_device_id"`
	ReplacementIMEI         string     `json:"replacement_imei"`
	ReplacementSerialNumber string     `json:"replacement_serial_number"`
	ReplacementValue        float64    `json:"replacement_value"`

	// Delivery Options
	DeliveryMethod         string         `json:"delivery_method"` // store_pickup, home_delivery, express_delivery
	DeliveryAddress        datatypes.JSON `gorm:"type:json" json:"delivery_address"`
	DeliveryTrackingNumber string         `json:"delivery_tracking_number"`
	EstimatedDeliveryDate  *time.Time     `json:"estimated_delivery_date"`
	ActualDeliveryDate     *time.Time     `json:"actual_delivery_date"`
	SignatureRequired      bool           `json:"signature_required"`

	// Data Transfer
	DataTransferMethod    string     `json:"data_transfer_method"` // cloud, cable, wireless
	DataTransferStatus    string     `json:"data_transfer_status"`
	DataTransferCompleted *time.Time `json:"data_transfer_completed"`

	// Old Device Return
	OldDeviceReturnRequired bool       `json:"old_device_return_required"`
	ReturnKitSent           bool       `json:"return_kit_sent"`
	ReturnTrackingNumber    string     `json:"return_tracking_number"`
	OldDeviceReceived       *time.Time `json:"old_device_received"`
	OldDeviceCondition      string     `json:"old_device_condition"`
	DataWipeConfirmed       bool       `json:"data_wipe_confirmed"`

	// Activation
	ActivationStatus string     `json:"activation_status"`
	ActivationDate   *time.Time `json:"activation_date"`
	SIMCardIncluded  bool       `json:"sim_card_included"`
	CarrierActivated bool       `json:"carrier_activated"`
}

// ClaimDigitalAssets - Digital content and app restoration
type ClaimDigitalAssets struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// App Purchases
	AppPurchasesValue   float64        `json:"app_purchases_value"`
	AppsList            datatypes.JSON `gorm:"type:json" json:"apps_list"`
	InAppPurchasesValue float64        `json:"in_app_purchases_value"`
	SubscriptionsValue  float64        `json:"subscriptions_value"`

	// Digital Content
	PhotosCount      int  `json:"photos_count"`
	VideosCount      int  `json:"videos_count"`
	ContactsCount    int  `json:"contacts_count"`
	MessagesBackedUp bool `json:"messages_backed_up"`

	// Cloud Services
	CloudBackupAvailable bool       `json:"cloud_backup_available"`
	LastCloudBackupDate  *time.Time `json:"last_cloud_backup_date"`
	CloudStorageProvider string     `json:"cloud_storage_provider"`
	BackupSize           int64      `json:"backup_size_gb"`

	// Recovery Status
	DataRecoveryAttempted   bool           `json:"data_recovery_attempted"`
	DataRecoverySuccessful  *bool          `json:"data_recovery_successful"`
	RecoveredDataPercentage float64        `json:"recovered_data_percentage"`
	UnrecoverableData       datatypes.JSON `gorm:"type:json" json:"unrecoverable_data"`

	// Compensation
	DigitalContentCovered bool    `json:"digital_content_covered"`
	CompensationAmount    float64 `json:"compensation_amount"`
	AppStoreCredit        float64 `json:"app_store_credit"`
	CloudStorageCredit    int     `json:"cloud_storage_credit_months"`
}

// ClaimAccessories - Accessory claims and coverage
type ClaimAccessories struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Claimed Accessories
	AccessoriesClaimed    datatypes.JSON `gorm:"type:json" json:"accessories_claimed"`
	TotalAccessoriesValue float64        `json:"total_accessories_value"`

	// Common Accessories
	CaseIncluded            bool    `json:"case_included"`
	CaseValue               float64 `json:"case_value"`
	ScreenProtectorIncluded bool    `json:"screen_protector_included"`
	ScreenProtectorValue    float64 `json:"screen_protector_value"`
	ChargerIncluded         bool    `json:"charger_included"`
	ChargerValue            float64 `json:"charger_value"`
	EarphonesIncluded       bool    `json:"earphones_included"`
	EarphonesValue          float64 `json:"earphones_value"`

	// Verification
	ReceiptsProvided      bool `json:"receipts_provided"`
	PhotosProvided        bool `json:"photos_provided"`
	SerialNumbersVerified bool `json:"serial_numbers_verified"`

	// Coverage
	AccessoriesCovered bool    `json:"accessories_covered"`
	CoverageLimit      float64 `json:"coverage_limit"`
	ApprovedAmount     float64 `json:"approved_amount"`
	ReplacementOffered bool    `json:"replacement_offered"`
}

// ClaimGeolocation - Location-based claim validation
type ClaimGeolocation struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Incident Location
	IncidentLatitude   float64 `json:"incident_latitude"`
	IncidentLongitude  float64 `json:"incident_longitude"`
	IncidentAddress    string  `gorm:"type:text" json:"incident_address"`
	IncidentCity       string  `json:"incident_city"`
	IncidentState      string  `json:"incident_state"`
	IncidentCountry    string  `json:"incident_country"`
	IncidentPostalCode string  `json:"incident_postal_code"`

	// Location Verification
	LocationVerified        bool           `json:"location_verified"`
	DeviceLastKnownLocation datatypes.JSON `gorm:"type:json" json:"device_last_known_location"`
	LocationMismatch        bool           `json:"location_mismatch"`
	LocationRiskScore       float64        `json:"location_risk_score"`

	// Coverage Validation
	LocationCovered       bool `json:"location_covered"`
	InternationalCoverage bool `json:"international_coverage"`
	TravelDateVerified    bool `json:"travel_date_verified"`

	// Nearby Services
	NearbyRepairCenters       datatypes.JSON `gorm:"type:json" json:"nearby_repair_centers"`
	NearbyStores              datatypes.JSON `gorm:"type:json" json:"nearby_stores"`
	EmergencyServicesNotified bool           `json:"emergency_services_notified"`

	// Weather Data
	WeatherConditions         string `json:"weather_conditions"`
	NaturalDisasterZone       bool   `json:"natural_disaster_zone"`
	DisasterType              string `json:"disaster_type"`
	DisasterDeclarationNumber string `json:"disaster_declaration_number"`
}

// ClaimPreventiveMeasures - Preventive measures and recommendations
type ClaimPreventiveMeasures struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Device Protection Status
	DeviceLockEnabled     bool `json:"device_lock_enabled"`
	FindMyDeviceEnabled   bool `json:"find_my_device_enabled"`
	RemoteWipeAvailable   bool `json:"remote_wipe_available"`
	AntiTheftAppInstalled bool `json:"anti_theft_app_installed"`
	DeviceEncrypted       bool `json:"device_encrypted"`

	// Physical Protection
	ScreenProtectorUsed bool `json:"screen_protector_used"`
	ProtectiveCaseUsed  bool `json:"protective_case_used"`
	WaterproofCaseUsed  bool `json:"waterproof_case_used"`

	// Usage Patterns
	HighRiskActivities    datatypes.JSON `gorm:"type:json" json:"high_risk_activities"`
	PreviousClaimsPattern datatypes.JSON `gorm:"type:json" json:"previous_claims_pattern"`
	RepeatClaimType       bool           `json:"repeat_claim_type"`

	// Recommendations
	PreventiveRecommendations datatypes.JSON `gorm:"type:json" json:"preventive_recommendations"`
	TrainingOffered           bool           `json:"training_offered"`
	SafetyTipsProvided        bool           `json:"safety_tips_provided"`
	DiscountEligible          bool           `json:"discount_eligible"`

	// Risk Reduction
	RiskReductionScore     float64 `json:"risk_reduction_score"`
	ComplianceWithTips     bool    `json:"compliance_with_tips"`
	FutureClaimProbability float64 `json:"future_claim_probability"`
}

// ClaimSelfService - Self-service claim features
type ClaimSelfService struct {
	database.BaseModel
	ClaimID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"claim_id"`

	// Self-Assessment
	SelfAssessmentCompleted bool    `json:"self_assessment_completed"`
	DamagePhotosUploaded    int     `json:"damage_photos_uploaded"`
	VideoUploaded           bool    `json:"video_uploaded"`
	AIAssessmentScore       float64 `json:"ai_assessment_score"`

	// Instant Actions
	InstantApprovalEligible bool    `json:"instant_approval_eligible"`
	InstantSettlementAmount float64 `json:"instant_settlement_amount"`
	VirtualInspectionDone   bool    `json:"virtual_inspection_done"`

	// Digital Documentation
	DigitalSignatureProvided bool `json:"digital_signature_provided"`
	EFormCompleted           bool `json:"e_form_completed"`
	DocumentsAutoExtracted   bool `json:"documents_auto_extracted"`
	OCRProcessed             bool `json:"ocr_processed"`

	// Chatbot Interaction
	ChatbotUsed              bool           `json:"chatbot_used"`
	ChatbotResolutionSuccess bool           `json:"chatbot_resolution_success"`
	ChatTranscript           datatypes.JSON `gorm:"type:json" json:"chat_transcript"`
	EscalatedToHuman         bool           `json:"escalated_to_human"`

	// Self-Service Actions
	AppointmentSelfBooked    bool `json:"appointment_self_booked"`
	RepairCenterSelfSelected bool `json:"repair_center_self_selected"`
	ReplacementSelfOrdered   bool `json:"replacement_self_ordered"`
	TrackingEnabled          bool `json:"tracking_enabled"`
}

// Business Logic Methods

func (d *ClaimDeviceDiagnostics) CalculateHealthScore() float64 {
	score := 100.0

	// Hardware deductions
	if d.ScreenHealth != "good" {
		score -= 20
	}
	if d.BatteryHealth < 80 {
		score -= 15
	}
	if !d.CellularWorking {
		score -= 15
	}
	if !d.WiFiWorking {
		score -= 10
	}
	if d.WaterDamageIndicator {
		score -= 30
	}

	// Sensor deductions
	if !d.AccelerometerWorking {
		score -= 5
	}
	// Camera status check - deduct if camera has issues (JSON field contains problem data)
	if d.CameraStatus != nil && len(d.CameraStatus) > 0 {
		// Assume if CameraStatus contains data, it indicates problems
		score -= 10
	}

	if score < 0 {
		score = 0
	}
	return score
}

func (r *ClaimRepairNetwork) IsExpressEligible() bool {
	return r.PartsAvailability == "in_stock" &&
		r.Distance < 10 &&
		r.RepairCenterRating > 4.0
}

func (r *ClaimReplacementDevice) NeedsOldDeviceReturn() bool {
	return r.OldDeviceReturnRequired &&
		!r.OldDeviceReceived.IsZero() &&
		r.ReplacementValue > 500
}

func (d *ClaimDigitalAssets) CalculateDigitalLoss() float64 {
	return d.AppPurchasesValue +
		d.InAppPurchasesValue +
		d.SubscriptionsValue
}

func (g *ClaimGeolocation) IsHighRiskLocation() bool {
	return g.LocationRiskScore > 0.7 ||
		g.NaturalDisasterZone ||
		!g.LocationCovered
}

func (p *ClaimPreventiveMeasures) CalculatePreventionScore() float64 {
	score := 0.0
	if p.DeviceLockEnabled {
		score += 20
	}
	if p.FindMyDeviceEnabled {
		score += 20
	}
	if p.ScreenProtectorUsed {
		score += 15
	}
	if p.ProtectiveCaseUsed {
		score += 15
	}
	if p.AntiTheftAppInstalled {
		score += 15
	}
	if p.DeviceEncrypted {
		score += 15
	}
	return score
}

func (s *ClaimSelfService) IsFullySelfServiced() bool {
	return s.SelfAssessmentCompleted &&
		s.InstantApprovalEligible &&
		s.DigitalSignatureProvided &&
		!s.EscalatedToHuman
}
