package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// CommunicationRepository defines the interface for communication persistence operations
type CommunicationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, communication *shared.Communication) error
	GetByID(ctx context.Context, id uuid.UUID) (*shared.Communication, error)
	Update(ctx context.Context, communication *shared.Communication) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByRecipient(ctx context.Context, recipientType string, recipientID uuid.UUID, limit, offset int) ([]*shared.Communication, int64, error)
	GetByType(ctx context.Context, commType string, limit, offset int) ([]*shared.Communication, int64, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.Communication, int64, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.Communication, int64, error)
	GetByTrigger(ctx context.Context, triggerType string, triggerID uuid.UUID) ([]*shared.Communication, error)
	Search(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*shared.Communication, int64, error)

	// Template operations
	CreateTemplate(ctx context.Context, template *shared.CommunicationTemplate) error
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*shared.CommunicationTemplate, error)
	GetTemplateByCode(ctx context.Context, code string) (*shared.CommunicationTemplate, error)
	UpdateTemplate(ctx context.Context, template *shared.CommunicationTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	ListTemplates(ctx context.Context, templateType, category string) ([]*shared.CommunicationTemplate, error)

	// Event operations
	CreateEvent(ctx context.Context, event *shared.CommunicationEvent) error
	GetEventsByCommunication(ctx context.Context, commID uuid.UUID) ([]*shared.CommunicationEvent, error)
	GetEventsByType(ctx context.Context, eventType string, startDate, endDate time.Time) ([]*shared.CommunicationEvent, error)

	// Preference operations
	GetPreferences(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.NotificationPreference, error)
	UpdatePreferences(ctx context.Context, prefs *shared.NotificationPreference) error

	// Campaign operations
	CreateCampaign(ctx context.Context, campaign *shared.CommunicationCampaign) error
	GetCampaignByID(ctx context.Context, id uuid.UUID) (*shared.CommunicationCampaign, error)
	GetCampaignByCode(ctx context.Context, code string) (*shared.CommunicationCampaign, error)
	UpdateCampaign(ctx context.Context, campaign *shared.CommunicationCampaign) error
	ListCampaigns(ctx context.Context, status string, limit, offset int) ([]*shared.CommunicationCampaign, int64, error)

	// Queue operations
	CreateQueueItem(ctx context.Context, item *shared.CommunicationQueue) error
	GetPendingQueueItems(ctx context.Context, limit int) ([]*shared.CommunicationQueue, error)
	UpdateQueueItem(ctx context.Context, item *shared.CommunicationQueue) error
	DeleteQueueItem(ctx context.Context, id uuid.UUID) error
}
