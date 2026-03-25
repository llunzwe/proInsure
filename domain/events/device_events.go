package events

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of device event
type EventType string

const (
	// Device Lifecycle Events
	DeviceRegisteredEvent  EventType = "device.registered"
	DeviceVerifiedEvent    EventType = "device.verified"
	DeviceTransferredEvent EventType = "device.transferred"
	DeviceStolenEvent      EventType = "device.stolen"
	DeviceRecoveredEvent   EventType = "device.recovered"
	DeviceLostEvent        EventType = "device.lost"
	DeviceDeactivatedEvent EventType = "device.deactivated"
	DeviceReactivatedEvent EventType = "device.reactivated"
	DeviceRetiredEvent     EventType = "device.retired"

	// Insurance Events
	DeviceInsuredEvent           EventType = "device.insured"
	DeviceUninsuredEvent         EventType = "device.uninsured"
	DevicePremiumCalculatedEvent EventType = "device.premium_calculated"
	DeviceClaimFiledEvent        EventType = "device.claim_filed"
	DeviceClaimApprovedEvent     EventType = "device.claim_approved"
	DeviceClaimRejectedEvent     EventType = "device.claim_rejected"

	// Risk Events
	DeviceRiskChangedEvent     EventType = "device.risk_changed"
	DeviceBlacklistedEvent     EventType = "device.blacklisted"
	DeviceWhitelistedEvent     EventType = "device.whitelisted"
	DeviceFraudDetectedEvent   EventType = "device.fraud_detected"
	DeviceAnomalyDetectedEvent EventType = "device.anomaly_detected"

	// Trade & Finance Events
	DeviceTradedInEvent           EventType = "device.traded_in"
	DeviceRentedEvent             EventType = "device.rented"
	DeviceRentalEndedEvent        EventType = "device.rental_ended"
	DeviceFinancedEvent           EventType = "device.financed"
	DeviceFinancingCompletedEvent EventType = "device.financing_completed"
	DeviceListedForSaleEvent      EventType = "device.listed_for_sale"
	DeviceSoldEvent               EventType = "device.sold"

	// Maintenance Events
	DeviceRepairScheduledEvent   EventType = "device.repair_scheduled"
	DeviceRepairCompletedEvent   EventType = "device.repair_completed"
	DeviceMaintenanceNeededEvent EventType = "device.maintenance_needed"
	DeviceRefurbishedEvent       EventType = "device.refurbished"

	// Condition Events
	DeviceConditionChangedEvent EventType = "device.condition_changed"
	DeviceDamagedEvent          EventType = "device.damaged"
	DeviceInspectedEvent        EventType = "device.inspected"
	DeviceGradedEvent           EventType = "device.graded"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	ID            uuid.UUID              `json:"id"`
	Type          EventType              `json:"type"`
	AggregateID   uuid.UUID              `json:"aggregate_id"` // Device ID
	Timestamp     time.Time              `json:"timestamp"`
	UserID        uuid.UUID              `json:"user_id,omitempty"`
	CorrelationID uuid.UUID              `json:"correlation_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// DeviceRegistered event
type DeviceRegistered struct {
	BaseEvent
	IMEI         string    `json:"imei"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	OwnerID      uuid.UUID `json:"owner_id"`
	RegisteredAt time.Time `json:"registered_at"`
}

// DeviceVerified event
type DeviceVerified struct {
	BaseEvent
	VerificationMethod  string                 `json:"verification_method"`
	VerifiedBy          uuid.UUID              `json:"verified_by"`
	IsAuthentic         bool                   `json:"is_authentic"`
	VerificationDetails map[string]interface{} `json:"verification_details"`
}

// DeviceTransferred event
type DeviceTransferred struct {
	BaseEvent
	FromOwnerID   uuid.UUID `json:"from_owner_id"`
	ToOwnerID     uuid.UUID `json:"to_owner_id"`
	TransferType  string    `json:"transfer_type"` // sale, gift, inheritance
	TransferredAt time.Time `json:"transferred_at"`
	TransferValue float64   `json:"transfer_value,omitempty"`
}

// DeviceStolen event
type DeviceStolen struct {
	BaseEvent
	ReportedBy         uuid.UUID `json:"reported_by"`
	ReportedAt         time.Time `json:"reported_at"`
	PoliceReportNumber string    `json:"police_report_number,omitempty"`
	Location           string    `json:"location,omitempty"`
	Description        string    `json:"description,omitempty"`
}

// DeviceRecovered event
type DeviceRecovered struct {
	BaseEvent
	RecoveredBy      uuid.UUID `json:"recovered_by"`
	RecoveredAt      time.Time `json:"recovered_at"`
	RecoveryLocation string    `json:"recovery_location,omitempty"`
	Condition        string    `json:"condition"`
}

// DeviceRiskChanged event
type DeviceRiskChanged struct {
	BaseEvent
	OldRiskScore float64  `json:"old_risk_score"`
	NewRiskScore float64  `json:"new_risk_score"`
	OldRiskLevel string   `json:"old_risk_level"`
	NewRiskLevel string   `json:"new_risk_level"`
	ChangeReason string   `json:"change_reason"`
	RiskFactors  []string `json:"risk_factors"`
}

// DeviceFraudDetected event
type DeviceFraudDetected struct {
	BaseEvent
	FraudType       string                 `json:"fraud_type"`
	FraudScore      float64                `json:"fraud_score"`
	Confidence      float64                `json:"confidence"`
	Indicators      []string               `json:"indicators"`
	DetectionMethod string                 `json:"detection_method"`
	ActionTaken     string                 `json:"action_taken"`
	Details         map[string]interface{} `json:"details"`
}

// DeviceConditionChanged event
type DeviceConditionChanged struct {
	BaseEvent
	OldCondition   string    `json:"old_condition"`
	NewCondition   string    `json:"new_condition"`
	OldGrade       string    `json:"old_grade"`
	NewGrade       string    `json:"new_grade"`
	AssessedBy     uuid.UUID `json:"assessed_by"`
	AssessmentDate time.Time `json:"assessment_date"`
	Reason         string    `json:"reason"`
}

// DeviceTradedIn event
type DeviceTradedIn struct {
	BaseEvent
	NewDeviceID   uuid.UUID `json:"new_device_id"`
	TradeInValue  float64   `json:"trade_in_value"`
	OriginalValue float64   `json:"original_value"`
	Condition     string    `json:"condition"`
	Grade         string    `json:"grade"`
	ProcessedAt   time.Time `json:"processed_at"`
}

// DeviceInsured event
type DeviceInsured struct {
	BaseEvent
	PolicyID       uuid.UUID `json:"policy_id"`
	PolicyNumber   string    `json:"policy_number"`
	InsuranceType  string    `json:"insurance_type"`
	CoverageAmount float64   `json:"coverage_amount"`
	Premium        float64   `json:"premium"`
	Deductible     float64   `json:"deductible"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

// DeviceClaimFiled event
type DeviceClaimFiled struct {
	BaseEvent
	ClaimID      uuid.UUID `json:"claim_id"`
	ClaimType    string    `json:"claim_type"`
	ClaimAmount  float64   `json:"claim_amount"`
	IncidentDate time.Time `json:"incident_date"`
	Description  string    `json:"description"`
	PolicyID     uuid.UUID `json:"policy_id"`
}

// DeviceClaimApproved event
type DeviceClaimApproved struct {
	BaseEvent
	ClaimID        uuid.UUID `json:"claim_id"`
	ApprovedAmount float64   `json:"approved_amount"`
	ApprovedBy     uuid.UUID `json:"approved_by"`
	ApprovalDate   time.Time `json:"approval_date"`
	PayoutMethod   string    `json:"payout_method"`
	Notes          string    `json:"notes,omitempty"`
}

// DeviceRepairScheduled event
type DeviceRepairScheduled struct {
	BaseEvent
	RepairID          uuid.UUID `json:"repair_id"`
	RepairType        string    `json:"repair_type"`
	ScheduledDate     time.Time `json:"scheduled_date"`
	RepairShopID      uuid.UUID `json:"repair_shop_id"`
	EstimatedCost     float64   `json:"estimated_cost"`
	EstimatedDuration string    `json:"estimated_duration"`
	IsWarrantyClaim   bool      `json:"is_warranty_claim"`
	IsInsuranceClaim  bool      `json:"is_insurance_claim"`
}

// DeviceRepairCompleted event
type DeviceRepairCompleted struct {
	BaseEvent
	RepairID         uuid.UUID `json:"repair_id"`
	CompletedDate    time.Time `json:"completed_date"`
	ActualCost       float64   `json:"actual_cost"`
	PartsReplaced    []string  `json:"parts_replaced"`
	RepairQuality    string    `json:"repair_quality"`
	WarrantyProvided bool      `json:"warranty_provided"`
	TechnicianID     uuid.UUID `json:"technician_id"`
}

// DevicePremiumCalculated event
type DevicePremiumCalculated struct {
	BaseEvent
	BaseAmount  float64                `json:"base_amount"`
	Adjustments map[string]float64     `json:"adjustments"`
	FinalAmount float64                `json:"final_amount"`
	RiskFactors map[string]interface{} `json:"risk_factors"`
	ValidUntil  time.Time              `json:"valid_until"`
	Currency    string                 `json:"currency"`
}

// === Event Handler Interface ===

// EventHandler defines the interface for handling events
type EventHandler interface {
	HandleEvent(event interface{}) error
	CanHandle(eventType EventType) bool
}

// === Event Bus Interface ===

// EventBus defines the interface for publishing and subscribing to events
type EventBus interface {
	Publish(event interface{}) error
	Subscribe(eventType EventType, handler EventHandler) error
	Unsubscribe(eventType EventType, handler EventHandler) error
}

// === Event Store Interface ===

// EventStore defines the interface for storing and retrieving events
type EventStore interface {
	Save(event interface{}) error
	GetEvents(aggregateID uuid.UUID) ([]interface{}, error)
	GetEventsByType(eventType EventType, limit int) ([]interface{}, error)
	GetEventsInRange(from, to time.Time) ([]interface{}, error)
}

// === Helper Functions ===

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType EventType, aggregateID uuid.UUID) BaseEvent {
	return BaseEvent{
		ID:          uuid.New(),
		Type:        eventType,
		AggregateID: aggregateID,
		Timestamp:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
}

// WithUser adds user ID to the event
func (e *BaseEvent) WithUser(userID uuid.UUID) *BaseEvent {
	e.UserID = userID
	return e
}

// WithCorrelation adds correlation ID to the event
func (e *BaseEvent) WithCorrelation(correlationID uuid.UUID) *BaseEvent {
	e.CorrelationID = correlationID
	return e
}

// WithMetadata adds metadata to the event
func (e *BaseEvent) WithMetadata(key string, value interface{}) *BaseEvent {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// === Event Factory ===

// DeviceEventFactory provides methods to create device events
type DeviceEventFactory struct{}

// NewDeviceEventFactory creates a new device event factory
func NewDeviceEventFactory() *DeviceEventFactory {
	return &DeviceEventFactory{}
}

// CreateDeviceRegistered creates a device registered event
func (f *DeviceEventFactory) CreateDeviceRegistered(deviceID uuid.UUID, imei, brand, model string, ownerID uuid.UUID) *DeviceRegistered {
	return &DeviceRegistered{
		BaseEvent:    NewBaseEvent(DeviceRegisteredEvent, deviceID),
		IMEI:         imei,
		Brand:        brand,
		Model:        model,
		OwnerID:      ownerID,
		RegisteredAt: time.Now(),
	}
}

// CreateDeviceStolen creates a device stolen event
func (f *DeviceEventFactory) CreateDeviceStolen(deviceID, reportedBy uuid.UUID, location, description string) *DeviceStolen {
	return &DeviceStolen{
		BaseEvent:   NewBaseEvent(DeviceStolenEvent, deviceID),
		ReportedBy:  reportedBy,
		ReportedAt:  time.Now(),
		Location:    location,
		Description: description,
	}
}

// CreateDeviceRiskChanged creates a device risk changed event
func (f *DeviceEventFactory) CreateDeviceRiskChanged(deviceID uuid.UUID, oldScore, newScore float64, oldLevel, newLevel, reason string, factors []string) *DeviceRiskChanged {
	return &DeviceRiskChanged{
		BaseEvent:    NewBaseEvent(DeviceRiskChangedEvent, deviceID),
		OldRiskScore: oldScore,
		NewRiskScore: newScore,
		OldRiskLevel: oldLevel,
		NewRiskLevel: newLevel,
		ChangeReason: reason,
		RiskFactors:  factors,
	}
}

// CreateDeviceFraudDetected creates a device fraud detected event
func (f *DeviceEventFactory) CreateDeviceFraudDetected(deviceID uuid.UUID, fraudType string, score, confidence float64, indicators []string) *DeviceFraudDetected {
	return &DeviceFraudDetected{
		BaseEvent:       NewBaseEvent(DeviceFraudDetectedEvent, deviceID),
		FraudType:       fraudType,
		FraudScore:      score,
		Confidence:      confidence,
		Indicators:      indicators,
		DetectionMethod: "automated",
		Details:         make(map[string]interface{}),
	}
}
