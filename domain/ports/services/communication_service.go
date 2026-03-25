package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// CommunicationService defines the interface for communication business logic
type CommunicationService interface {
	// Send operations
	SendCommunication(ctx context.Context, communication *shared.Communication) error
	SendBulkCommunication(ctx context.Context, communications []*shared.Communication) error
	ScheduleCommunication(ctx context.Context, communication *shared.Communication, scheduledFor time.Time) error

	// Retrieve operations
	GetCommunication(ctx context.Context, id uuid.UUID) (*shared.Communication, error)
	GetCommunicationsByRecipient(ctx context.Context, recipientType string, recipientID uuid.UUID, limit, offset int) ([]*shared.Communication, int64, error)
	GetCommunicationsByType(ctx context.Context, commType string, limit, offset int) ([]*shared.Communication, int64, error)

	// Template operations
	CreateTemplate(ctx context.Context, template *shared.CommunicationTemplate) error
	GetTemplate(ctx context.Context, id uuid.UUID) (*shared.CommunicationTemplate, error)
	GetTemplateByCode(ctx context.Context, code string) (*shared.CommunicationTemplate, error)
	UpdateTemplate(ctx context.Context, template *shared.CommunicationTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	ListTemplates(ctx context.Context, templateType, category string) ([]*shared.CommunicationTemplate, error)

	// Event tracking
	TrackEvent(ctx context.Context, commID uuid.UUID, eventType string, eventData map[string]interface{}) error
	GetCommunicationEvents(ctx context.Context, commID uuid.UUID) ([]*shared.CommunicationEvent, error)

	// Preference operations
	GetPreferences(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.NotificationPreference, error)
	UpdatePreferences(ctx context.Context, entityType string, entityID uuid.UUID, prefs *shared.NotificationPreference) error

	// Campaign operations
	CreateCampaign(ctx context.Context, campaign *shared.CommunicationCampaign) error
	GetCampaign(ctx context.Context, id uuid.UUID) (*shared.CommunicationCampaign, error)
	StartCampaign(ctx context.Context, campaignID uuid.UUID) error
	StopCampaign(ctx context.Context, campaignID uuid.UUID) error
	ListCampaigns(ctx context.Context, status string, limit, offset int) ([]*shared.CommunicationCampaign, int64, error)

	// Queue operations
	ProcessQueue(ctx context.Context) error
	GetQueueStatus(ctx context.Context) (map[string]interface{}, error)

	// Notification operations
	SendNotification(ctx context.Context, recipientID uuid.UUID, title, message string) error
	SendNotificationWithData(ctx context.Context, recipientID uuid.UUID, title, message string, data map[string]interface{}) error
}
