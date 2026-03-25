package communication

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// CommunicationService implements domain-level communication business logic
type CommunicationService struct {
	repo repositories.CommunicationRepository
}

// NewCommunicationService creates a new communication service
func NewCommunicationService(repo repositories.CommunicationRepository) services.CommunicationService {
	return &CommunicationService{
		repo: repo,
	}
}

// SendCommunication sends a communication
func (s *CommunicationService) SendCommunication(ctx context.Context, communication *shared.Communication) error {
	// Validate communication
	if err := s.validateCommunication(communication); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set defaults
	if communication.Status == "" {
		communication.Status = "queued"
	}
	if communication.CreatedAt.IsZero() {
		communication.CreatedAt = time.Now()
	}

	// Create communication record
	return s.repo.Create(ctx, communication)
}

// SendBulkCommunication sends multiple communications
func (s *CommunicationService) SendBulkCommunication(ctx context.Context, communications []*shared.Communication) error {
	for _, comm := range communications {
		if err := s.SendCommunication(ctx, comm); err != nil {
			return fmt.Errorf("failed to send communication %s: %w", comm.CommunicationCode, err)
		}
	}
	return nil
}

// ScheduleCommunication schedules a communication for later
func (s *CommunicationService) ScheduleCommunication(ctx context.Context, communication *shared.Communication, scheduledFor time.Time) error {
	if scheduledFor.Before(time.Now()) {
		return fmt.Errorf("scheduled time must be in the future")
	}

	communication.ScheduledFor = &scheduledFor
	communication.Status = "scheduled"

	return s.repo.Create(ctx, communication)
}

// GetCommunication retrieves a communication by ID
func (s *CommunicationService) GetCommunication(ctx context.Context, id uuid.UUID) (*shared.Communication, error) {
	return s.repo.GetByID(ctx, id)
}

// GetCommunicationsByRecipient gets communications for a recipient
func (s *CommunicationService) GetCommunicationsByRecipient(ctx context.Context, recipientType string, recipientID uuid.UUID, limit, offset int) ([]*shared.Communication, int64, error) {
	return s.repo.GetByRecipient(ctx, recipientType, recipientID, limit, offset)
}

// GetCommunicationsByType gets communications by type
func (s *CommunicationService) GetCommunicationsByType(ctx context.Context, commType string, limit, offset int) ([]*shared.Communication, int64, error) {
	return s.repo.GetByType(ctx, commType, limit, offset)
}

// Template operations

// CreateTemplate creates a new communication template
func (s *CommunicationService) CreateTemplate(ctx context.Context, template *shared.CommunicationTemplate) error {
	if template.TemplateCode == "" {
		return fmt.Errorf("template code is required")
	}
	if template.Body == "" {
		return fmt.Errorf("template body is required")
	}

	return s.repo.CreateTemplate(ctx, template)
}

// GetTemplate retrieves a template by ID
func (s *CommunicationService) GetTemplate(ctx context.Context, id uuid.UUID) (*shared.CommunicationTemplate, error) {
	return s.repo.GetTemplateByID(ctx, id)
}

// GetTemplateByCode retrieves a template by code
func (s *CommunicationService) GetTemplateByCode(ctx context.Context, code string) (*shared.CommunicationTemplate, error) {
	return s.repo.GetTemplateByCode(ctx, code)
}

// UpdateTemplate updates a template
func (s *CommunicationService) UpdateTemplate(ctx context.Context, template *shared.CommunicationTemplate) error {
	return s.repo.UpdateTemplate(ctx, template)
}

// DeleteTemplate deletes a template
func (s *CommunicationService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, id)
}

// ListTemplates lists templates by type and category
func (s *CommunicationService) ListTemplates(ctx context.Context, templateType, category string) ([]*shared.CommunicationTemplate, error) {
	return s.repo.ListTemplates(ctx, templateType, category)
}

// Event tracking

// TrackEvent tracks a communication event
func (s *CommunicationService) TrackEvent(ctx context.Context, commID uuid.UUID, eventType string, eventData map[string]interface{}) error {
	event := &shared.CommunicationEvent{
		CommunicationID: commID,
		EventType:       eventType,
		CreatedAt:       time.Now(),
	}

	// Serialize event data if provided
	if eventData != nil {
		// Note: In production, properly serialize to JSON string
		event.EventData = fmt.Sprintf("%v", eventData)
	}

	return s.repo.CreateEvent(ctx, event)
}

// GetCommunicationEvents gets events for a communication
func (s *CommunicationService) GetCommunicationEvents(ctx context.Context, commID uuid.UUID) ([]*shared.CommunicationEvent, error) {
	return s.repo.GetEventsByCommunication(ctx, commID)
}

// Preference operations

// GetPreferences retrieves notification preferences
func (s *CommunicationService) GetPreferences(ctx context.Context, entityType string, entityID uuid.UUID) (*shared.NotificationPreference, error) {
	return s.repo.GetPreferences(ctx, entityType, entityID)
}

// UpdatePreferences updates notification preferences
func (s *CommunicationService) UpdatePreferences(ctx context.Context, entityType string, entityID uuid.UUID, prefs *shared.NotificationPreference) error {
	prefs.EntityType = entityType
	prefs.EntityID = entityID
	return s.repo.UpdatePreferences(ctx, prefs)
}

// Campaign operations

// CreateCampaign creates a new communication campaign
func (s *CommunicationService) CreateCampaign(ctx context.Context, campaign *shared.CommunicationCampaign) error {
	if campaign.CampaignCode == "" {
		return fmt.Errorf("campaign code is required")
	}
	if campaign.Name == "" {
		return fmt.Errorf("campaign name is required")
	}

	return s.repo.CreateCampaign(ctx, campaign)
}

// GetCampaign retrieves a campaign by ID
func (s *CommunicationService) GetCampaign(ctx context.Context, id uuid.UUID) (*shared.CommunicationCampaign, error) {
	return s.repo.GetCampaignByID(ctx, id)
}

// StartCampaign starts a campaign
func (s *CommunicationService) StartCampaign(ctx context.Context, campaignID uuid.UUID) error {
	campaign, err := s.repo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	now := time.Now()
	campaign.Status = "running"
	campaign.StartedAt = &now

	return s.repo.UpdateCampaign(ctx, campaign)
}

// StopCampaign stops a campaign
func (s *CommunicationService) StopCampaign(ctx context.Context, campaignID uuid.UUID) error {
	campaign, err := s.repo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	now := time.Now()
	campaign.Status = "completed"
	campaign.CompletedAt = &now

	return s.repo.UpdateCampaign(ctx, campaign)
}

// ListCampaigns lists campaigns by status
func (s *CommunicationService) ListCampaigns(ctx context.Context, status string, limit, offset int) ([]*shared.CommunicationCampaign, int64, error) {
	return s.repo.ListCampaigns(ctx, status, limit, offset)
}

// Queue operations

// ProcessQueue processes pending communications in queue
func (s *CommunicationService) ProcessQueue(ctx context.Context) error {
	// Get pending queue items
	items, err := s.repo.GetPendingQueueItems(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to get queue items: %w", err)
	}

	// Process each item
	for _, item := range items {
		// In production, this would trigger actual sending
		// For now, just update status
		item.Status = "processing"
		if err := s.repo.UpdateQueueItem(ctx, item); err != nil {
			// Log error but continue processing
			continue
		}
	}

	return nil
}

// GetQueueStatus gets queue status
func (s *CommunicationService) GetQueueStatus(ctx context.Context) (map[string]interface{}, error) {
	pending, err := s.repo.GetPendingQueueItems(ctx, 1000)
	if err != nil {
		return nil, err
	}

	status := map[string]interface{}{
		"pending_count": len(pending),
		"status":        "operational",
	}

	return status, nil
}

// Notification operations

// SendNotification sends a notification
func (s *CommunicationService) SendNotification(ctx context.Context, recipientID uuid.UUID, title, message string) error {
	comm := &shared.Communication{
		RecipientType:     "user",
		RecipientID:       recipientID,
		Type:              "in_app",
		Category:          "notification",
		Priority:          "normal",
		Subject:           title,
		Body:              message,
		Status:            "queued",
		CommunicationCode: fmt.Sprintf("notif_%s_%d", recipientID.String(), time.Now().Unix()),
	}

	return s.SendCommunication(ctx, comm)
}

// SendNotificationWithData sends a notification with additional data
func (s *CommunicationService) SendNotificationWithData(ctx context.Context, recipientID uuid.UUID, title, message string, data map[string]interface{}) error {
	comm := &shared.Communication{
		RecipientType:     "user",
		RecipientID:       recipientID,
		Type:              "in_app",
		Category:          "notification",
		Priority:          "normal",
		Subject:           title,
		Body:              message,
		Status:            "queued",
		CommunicationCode: fmt.Sprintf("notif_%s_%d", recipientID.String(), time.Now().Unix()),
	}

	// Note: In production, properly serialize data to JSON
	if data != nil {
		comm.Metadata = fmt.Sprintf("%v", data)
	}

	return s.SendCommunication(ctx, comm)
}

// Helper methods

func (s *CommunicationService) validateCommunication(comm *shared.Communication) error {
	if comm.RecipientType == "" {
		return fmt.Errorf("recipient type is required")
	}
	if comm.RecipientID == uuid.Nil {
		return fmt.Errorf("recipient ID is required")
	}
	if comm.Type == "" {
		return fmt.Errorf("communication type is required")
	}
	if comm.Category == "" {
		return fmt.Errorf("communication category is required")
	}
	if comm.CommunicationCode == "" {
		return fmt.Errorf("communication code is required")
	}
	return nil
}
