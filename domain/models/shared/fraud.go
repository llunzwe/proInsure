package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
)

// FraudDetection represents fraud detection analysis for claims and users
type FraudDetection struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ClaimID             *uuid.UUID     `gorm:"type:uuid" json:"claim_id"`
	UserID              *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	DeviceID            *uuid.UUID     `gorm:"type:uuid" json:"device_id"`
	FraudScore          float64        `gorm:"not null" json:"fraud_score"`   // 0-100
	RiskLevel           string         `gorm:"not null" json:"risk_level"`    // low, medium, high, critical
	AnalysisType        string         `gorm:"not null" json:"analysis_type"` // claim, user_behavior, device_pattern
	Indicators          string         `json:"indicators"`                    // JSON array of fraud indicators
	BehavioralPatterns  string         `json:"behavioral_patterns"`           // JSON object
	SocialMediaFindings string         `json:"social_media_findings"`         // JSON object
	NetworkAnalysis     string         `json:"network_analysis"`              // JSON object
	MLModelVersion      string         `json:"ml_model_version"`
	ConfidenceLevel     float64        `json:"confidence_level"`
	IsManualReview      bool           `gorm:"default:false" json:"is_manual_review"`
	ReviewedBy          *uuid.UUID     `gorm:"type:uuid" json:"reviewed_by"`
	ReviewedAt          *time.Time     `json:"reviewed_at"`
	ReviewNotes         string         `json:"review_notes"`
	Status              string         `gorm:"not null;default:'pending'" json:"status"` // pending, cleared, flagged, blocked
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Claim           *models.Claim    `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	User            *models.User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device          *models.Device   `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	FraudIndicators []FraudIndicator `gorm:"foreignKey:FraudDetectionID" json:"fraud_indicators,omitempty"`
}

// FraudIndicator represents specific fraud indicators found during analysis
type FraudIndicator struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	FraudDetectionID uuid.UUID `gorm:"type:uuid;not null" json:"fraud_detection_id"`
	IndicatorType    string    `gorm:"not null" json:"indicator_type"` // behavioral, temporal, network, device, claim_pattern
	IndicatorName    string    `gorm:"not null" json:"indicator_name"`
	Description      string    `json:"description"`
	Severity         string    `gorm:"not null" json:"severity"` // low, medium, high, critical
	Weight           float64   `json:"weight"`                   // contribution to overall fraud score
	Evidence         string    `json:"evidence"`                 // JSON object with evidence
	IsConfirmed      bool      `gorm:"default:false" json:"is_confirmed"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	FraudDetection FraudDetection `gorm:"foreignKey:FraudDetectionID" json:"fraud_detection,omitempty"`
}

// BlacklistedDevice represents devices that are blacklisted for fraud
type BlacklistedDevice struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	IMEI          string         `gorm:"uniqueIndex;not null" json:"imei"`
	SerialNumber  *string        `gorm:"uniqueIndex" json:"serial_number"`
	Reason        string         `gorm:"not null" json:"reason"`
	Source        string         `gorm:"not null" json:"source"` // internal, law_enforcement, industry_database
	BlacklistedBy uuid.UUID      `gorm:"type:uuid;not null" json:"blacklisted_by"`
	BlacklistedAt time.Time      `gorm:"not null" json:"blacklisted_at"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	ExpiresAt     *time.Time     `json:"expires_at"`
	Notes         string         `json:"notes"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserBehaviorProfile represents behavioral patterns for fraud detection
type UserBehaviorProfile struct {
	ID                    uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID                uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	ClaimFrequency        float64        `json:"claim_frequency"` // claims per year
	AverageClaimValue     float64        `json:"average_claim_value"`
	ClaimPatterns         string         `json:"claim_patterns"` // JSON object
	DeviceChangeFrequency float64        `json:"device_change_frequency"`
	LoginPatterns         string         `json:"login_patterns"`     // JSON object
	LocationPatterns      string         `json:"location_patterns"`  // JSON object
	PaymentPatterns       string         `json:"payment_patterns"`   // JSON object
	SocialConnections     string         `json:"social_connections"` // JSON object
	RiskScore             float64        `gorm:"default:0" json:"risk_score"`
	LastAnalyzed          time.Time      `json:"last_analyzed"`
	CreatedAt             time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// FraudNetwork represents networks of potentially fraudulent users/devices
type FraudNetwork struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	NetworkName     string         `gorm:"not null" json:"network_name"`
	NetworkType     string         `gorm:"not null" json:"network_type"` // device_sharing, location_cluster, payment_method, social_network
	Description     string         `json:"description"`
	RiskLevel       string         `gorm:"not null" json:"risk_level"` // low, medium, high, critical
	MemberCount     int            `gorm:"default:0" json:"member_count"`
	TotalClaimValue float64        `gorm:"default:0" json:"total_claim_value"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	DetectedAt      time.Time      `gorm:"not null" json:"detected_at"`
	LastActivity    time.Time      `json:"last_activity"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	NetworkMembers []FraudNetworkMember `gorm:"foreignKey:NetworkID" json:"network_members,omitempty"`
}

// FraudNetworkMember represents members of fraud networks
type FraudNetworkMember struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	NetworkID          uuid.UUID  `gorm:"type:uuid;not null" json:"network_id"`
	UserID             *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	DeviceID           *uuid.UUID `gorm:"type:uuid" json:"device_id"`
	MemberRole         string     `json:"member_role"`         // primary, secondary, peripheral
	ConnectionStrength float64    `json:"connection_strength"` // 0-1
	JoinedAt           time.Time  `gorm:"not null" json:"joined_at"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Network FraudNetwork   `gorm:"foreignKey:NetworkID" json:"network,omitempty"`
	User    *models.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device  *models.Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// TableName methods
func (FraudDetection) TableName() string {
	return "fraud_detections"
}

func (FraudIndicator) TableName() string {
	return "fraud_indicators"
}

func (BlacklistedDevice) TableName() string {
	return "blacklisted_devices"
}

func (UserBehaviorProfile) TableName() string {
	return "user_behavior_profiles"
}

func (FraudNetwork) TableName() string {
	return "fraud_networks"
}

func (FraudNetworkMember) TableName() string {
	return "fraud_network_members"
}

// BeforeCreate hooks
func (fd *FraudDetection) BeforeCreate(tx *gorm.DB) error {
	if fd.ID == uuid.Nil {
		fd.ID = uuid.New()
	}
	return nil
}

func (fi *FraudIndicator) BeforeCreate(tx *gorm.DB) error {
	if fi.ID == uuid.Nil {
		fi.ID = uuid.New()
	}
	return nil
}

func (bd *BlacklistedDevice) BeforeCreate(tx *gorm.DB) error {
	if bd.ID == uuid.Nil {
		bd.ID = uuid.New()
	}
	return nil
}

func (ubp *UserBehaviorProfile) BeforeCreate(tx *gorm.DB) error {
	if ubp.ID == uuid.Nil {
		ubp.ID = uuid.New()
	}
	return nil
}

func (fn *FraudNetwork) BeforeCreate(tx *gorm.DB) error {
	if fn.ID == uuid.Nil {
		fn.ID = uuid.New()
	}
	return nil
}

func (fnm *FraudNetworkMember) BeforeCreate(tx *gorm.DB) error {
	if fnm.ID == uuid.Nil {
		fnm.ID = uuid.New()
	}
	return nil
}

// Business logic methods for FraudDetection
func (fd *FraudDetection) IsHighRisk() bool {
	return fd.RiskLevel == "high" || fd.RiskLevel == "critical"
}

func (fd *FraudDetection) RequiresManualReview() bool {
	return fd.FraudScore >= 70 || fd.IsHighRisk()
}

func (fd *FraudDetection) Flag() {
	fd.Status = "flagged"
}

func (fd *FraudDetection) Clear() {
	fd.Status = "cleared"
}

func (fd *FraudDetection) Block() {
	fd.Status = "blocked"
}

// Business logic methods for BlacklistedDevice
func (bd *BlacklistedDevice) IsExpired() bool {
	if bd.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*bd.ExpiresAt)
}

func (bd *BlacklistedDevice) Deactivate() {
	bd.IsActive = false
}

// Business logic methods for UserBehaviorProfile
func (ubp *UserBehaviorProfile) UpdateRiskScore() {
	// Simple risk scoring algorithm
	score := 0.0

	// High claim frequency increases risk
	if ubp.ClaimFrequency > 2 {
		score += 30
	} else if ubp.ClaimFrequency > 1 {
		score += 15
	}

	// High average claim value increases risk
	if ubp.AverageClaimValue > 1000 {
		score += 25
	} else if ubp.AverageClaimValue > 500 {
		score += 10
	}

	// Frequent device changes increase risk
	if ubp.DeviceChangeFrequency > 3 {
		score += 20
	} else if ubp.DeviceChangeFrequency > 1 {
		score += 10
	}

	ubp.RiskScore = score
}

func (ubp *UserBehaviorProfile) IsHighRisk() bool {
	return ubp.RiskScore >= 70
}

// Business logic methods for FraudNetwork
func (fn *FraudNetwork) AddMember(userID *uuid.UUID, deviceID *uuid.UUID, role string) *FraudNetworkMember {
	member := &FraudNetworkMember{
		NetworkID:          fn.ID,
		UserID:             userID,
		DeviceID:           deviceID,
		MemberRole:         role,
		ConnectionStrength: 1.0,
		JoinedAt:           time.Now(),
	}
	fn.MemberCount++
	return member
}

func (fn *FraudNetwork) Deactivate() {
	fn.IsActive = false
}
