package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceVerificationRecord represents a device verification record
type DeviceVerificationRecord struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	DeviceID         *uuid.UUID `gorm:"type:uuid" json:"device_id,omitempty"`
	IMEI             string     `gorm:"not null;index" json:"imei"`
	IMEI2            string     `json:"imei2,omitempty"` // For dual-SIM devices
	SerialNumber     string     `json:"serial_number,omitempty"`
	VerificationType string     `gorm:"not null" json:"verification_type"` // imei_check, blacklist, manufacturer, carrier
	Status           string     `gorm:"not null" json:"status"`            // valid, invalid, stolen, blacklisted, unknown
	Provider         string     `json:"provider"`                          // gsma, checkmend, apple, samsung
	ProviderResponse string     `json:"provider_response"`                 // JSON response from provider

	// Device Details from Verification
	Manufacturer   string     `json:"manufacturer,omitempty"`
	Model          string     `json:"model,omitempty"`
	ModelName      string     `json:"model_name,omitempty"`
	DeviceType     string     `json:"device_type,omitempty"` // smartphone, tablet, wearable
	ReleaseDate    *time.Time `json:"release_date,omitempty"`
	TechnicalSpecs string     `json:"technical_specs,omitempty"` // JSON object

	// Blacklist & Theft Information
	IsBlacklisted      bool       `gorm:"default:false" json:"is_blacklisted"`
	BlacklistDate      *time.Time `json:"blacklist_date,omitempty"`
	BlacklistReason    string     `json:"blacklist_reason,omitempty"`
	BlacklistCountry   string     `json:"blacklist_country,omitempty"`
	IsStolen           bool       `gorm:"default:false" json:"is_stolen"`
	StolenDate         *time.Time `json:"stolen_date,omitempty"`
	StolenLocation     string     `json:"stolen_location,omitempty"`
	PoliceReportNumber string     `json:"police_report_number,omitempty"`

	// Carrier Information
	CarrierLocked    bool       `gorm:"default:false" json:"carrier_locked"`
	CurrentCarrier   string     `json:"current_carrier,omitempty"`
	OriginalCarrier  string     `json:"original_carrier,omitempty"`
	ActivationStatus string     `json:"activation_status,omitempty"` // activated, not_activated
	ActivationDate   *time.Time `json:"activation_date,omitempty"`

	// Warranty Information
	WarrantyStatus  string     `json:"warranty_status,omitempty"` // active, expired, void
	WarrantyExpiry  *time.Time `json:"warranty_expiry,omitempty"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	PurchaseCountry string     `json:"purchase_country,omitempty"`

	// Verification Results
	VerificationScore   float64 `gorm:"default:0" json:"verification_score"` // 0-100
	RiskLevel           string  `json:"risk_level,omitempty"`                // low, medium, high
	VerificationNotes   string  `json:"verification_notes,omitempty"`
	IsEligible          bool    `gorm:"default:false" json:"is_eligible"`
	IneligibilityReason string  `json:"ineligibility_reason,omitempty"`

	// Metadata
	RequestIP  string         `json:"request_ip,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
	VerifiedBy *uuid.UUID     `gorm:"type:uuid" json:"verified_by,omitempty"`
	Metadata   string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	// Device should be loaded via service layer using DeviceID to avoid circular import
}

// IMEIDatabase represents cached IMEI database entries
type IMEIDatabase struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TAC                string     `gorm:"uniqueIndex;not null" json:"tac"` // Type Allocation Code (first 8 digits)
	Manufacturer       string     `gorm:"not null" json:"manufacturer"`
	Model              string     `gorm:"not null" json:"model"`
	ModelName          string     `json:"model_name"`
	DeviceType         string     `json:"device_type"`
	Capabilities       string     `json:"capabilities"` // JSON array
	LaunchDate         *time.Time `json:"launch_date,omitempty"`
	DiscontinuedDate   *time.Time `json:"discontinued_date,omitempty"`
	AveragePrice       float64    `json:"average_price"`
	RepairCostEstimate float64    `json:"repair_cost_estimate"`
	PopularityScore    int        `json:"popularity_score"`
	IsActive           bool       `gorm:"default:true" json:"is_active"`
	LastUpdated        time.Time  `json:"last_updated"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// BlacklistEntry represents blacklisted devices
type BlacklistEntry struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	IMEI            string         `gorm:"uniqueIndex;not null" json:"imei"`
	DeviceID        *uuid.UUID     `gorm:"type:uuid" json:"device_id,omitempty"`
	Reason          string         `gorm:"not null" json:"reason"` // stolen, lost, unpaid, fraud
	ReportedDate    time.Time      `gorm:"not null" json:"reported_date"`
	ReportedBy      string         `json:"reported_by"` // carrier, police, insurance, user
	Country         string         `json:"country"`
	Region          string         `json:"region,omitempty"`
	CaseNumber      string         `json:"case_number,omitempty"`
	ContactInfo     string         `json:"contact_info,omitempty"`
	Status          string         `gorm:"default:'active'" json:"status"` // active, resolved, pending
	ResolvedDate    *time.Time     `json:"resolved_date,omitempty"`
	ResolutionNotes string         `json:"resolution_notes,omitempty"`
	Source          string         `json:"source"`             // gsma, local_db, carrier, police
	Metadata        string         `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// DeviceHistory tracks device ownership and insurance history
type DeviceHistory struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	IMEI            string     `gorm:"not null;index" json:"imei"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	DeviceID        *uuid.UUID `gorm:"type:uuid" json:"device_id,omitempty"`
	EventType       string     `gorm:"not null" json:"event_type"` // registered, transferred, claimed, repaired, replaced
	EventDate       time.Time  `gorm:"not null" json:"event_date"`
	Description     string     `json:"description"`
	PreviousOwnerID *uuid.UUID `gorm:"type:uuid" json:"previous_owner_id,omitempty"`
	NewOwnerID      *uuid.UUID `gorm:"type:uuid" json:"new_owner_id,omitempty"`
	PolicyID        *uuid.UUID `gorm:"type:uuid" json:"policy_id,omitempty"`
	ClaimID         *uuid.UUID `gorm:"type:uuid" json:"claim_id,omitempty"`
	Location        string     `json:"location,omitempty"`
	IPAddress       string     `json:"ip_address,omitempty"`
	Metadata        string     `json:"metadata,omitempty"` // JSON for additional data
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hooks for UUID generation
func (d *DeviceVerificationRecord) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (i *IMEIDatabase) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (b *BlacklistEntry) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (d *DeviceHistory) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
