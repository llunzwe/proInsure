package support

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
	"smartsure/internal/domain/ports/repositories"
	"smartsure/internal/domain/ports/services"
)

// supportService implements the SupportService interface
type supportService struct {
	supportRepo repositories.SupportRepository
}

// NewSupportService creates a new support service
func NewSupportService(
	supportRepo repositories.SupportRepository,
) services.SupportService {
	return &supportService{
		supportRepo: supportRepo,
	}
}

// CreateTicket creates a new support ticket
func (s *supportService) CreateTicket(ctx context.Context, ticket *shared.SupportTicket) error {
	if ticket == nil {
		return errors.New("ticket cannot be nil")
	}
	if ticket.ID == uuid.Nil {
		ticket.ID = uuid.New()
	}
	if ticket.TicketNumber == "" {
		ticket.TicketNumber = fmt.Sprintf("TKT-%s", uuid.New().String()[:8])
	}
	if ticket.Status == "" {
		ticket.Status = "open"
	}
	return s.supportRepo.CreateTicket(ctx, ticket)
}

// GetTicketByID retrieves a ticket by ID
func (s *supportService) GetTicketByID(ctx context.Context, id uuid.UUID) (*shared.SupportTicket, error) {
	return s.supportRepo.GetTicketByID(ctx, id)
}

// GetTicketByTicketNumber retrieves a ticket by ticket number
func (s *supportService) GetTicketByTicketNumber(ctx context.Context, ticketNumber string) (*shared.SupportTicket, error) {
	return s.supportRepo.GetTicketByTicketNumber(ctx, ticketNumber)
}

// UpdateTicket updates an existing ticket
func (s *supportService) UpdateTicket(ctx context.Context, ticket *shared.SupportTicket) error {
	if ticket == nil {
		return errors.New("ticket cannot be nil")
	}
	return s.supportRepo.UpdateTicket(ctx, ticket)
}

// CloseTicket closes a ticket
func (s *supportService) CloseTicket(ctx context.Context, ticketID uuid.UUID, resolutionNotes string) error {
	ticket, err := s.supportRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}
	now := time.Now()
	ticket.Status = "closed"
	ticket.ResolutionNotes = resolutionNotes
	ticket.ResolvedAt = &now
	return s.supportRepo.UpdateTicket(ctx, ticket)
}

// GetTicketsByUserID gets tickets for a user
func (s *supportService) GetTicketsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.SupportTicket, int64, error) {
	return s.supportRepo.GetTicketsByUserID(ctx, userID, limit, offset)
}

// GetTicketsByStatus gets tickets by status
func (s *supportService) GetTicketsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.SupportTicket, int64, error) {
	return s.supportRepo.GetTicketsByStatus(ctx, status, limit, offset)
}

// GetTicketsByPriority gets tickets by priority
func (s *supportService) GetTicketsByPriority(ctx context.Context, priority string, limit, offset int) ([]*shared.SupportTicket, int64, error) {
	return s.supportRepo.GetTicketsByPriority(ctx, priority, limit, offset)
}

// GetTicketsByCategory gets tickets by category
func (s *supportService) GetTicketsByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.SupportTicket, int64, error) {
	return s.supportRepo.GetTicketsByCategory(ctx, category, limit, offset)
}

// GetTicketsByAssignedTo gets tickets assigned to an agent
func (s *supportService) GetTicketsByAssignedTo(ctx context.Context, agentID uuid.UUID, limit, offset int) ([]*shared.SupportTicket, int64, error) {
	return s.supportRepo.GetTicketsByAssignedTo(ctx, agentID, limit, offset)
}

// GetOpenTickets gets open tickets
func (s *supportService) GetOpenTickets(ctx context.Context, limit, offset int) ([]*shared.SupportTicket, int64, error) {
	return s.supportRepo.GetOpenTickets(ctx, limit, offset)
}

// SearchTickets searches for tickets
func (s *supportService) SearchTickets(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*shared.SupportTicket, int64, error) {
	return s.supportRepo.SearchTickets(ctx, query, filters, limit, offset)
}

// AssignTicket assigns a ticket to an agent
func (s *supportService) AssignTicket(ctx context.Context, ticketID uuid.UUID, agentID uuid.UUID) error {
	ticket, err := s.supportRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}
	now := time.Now()
	ticket.AssignedTo = &agentID
	ticket.AssignedAt = &now
	return s.supportRepo.UpdateTicket(ctx, ticket)
}

// ReassignTicket reassigns a ticket to a different agent
func (s *supportService) ReassignTicket(ctx context.Context, ticketID uuid.UUID, newAgentID uuid.UUID, reason string) error {
	return s.AssignTicket(ctx, ticketID, newAgentID)
}

// UnassignTicket unassigns a ticket
func (s *supportService) UnassignTicket(ctx context.Context, ticketID uuid.UUID) error {
	ticket, err := s.supportRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}
	if ticket == nil {
		return errors.New("ticket not found")
	}
	ticket.AssignedTo = nil
	ticket.AssignedAt = nil
	return s.supportRepo.UpdateTicket(ctx, ticket)
}

// AddTicketMessage adds a message to a ticket
func (s *supportService) AddTicketMessage(ctx context.Context, message *shared.TicketMessage) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}
	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}
	return s.supportRepo.CreateTicketMessage(ctx, message)
}

// GetTicketMessages gets messages for a ticket
func (s *supportService) GetTicketMessages(ctx context.Context, ticketID uuid.UUID, limit, offset int) ([]*shared.TicketMessage, int64, error) {
	return s.supportRepo.GetTicketMessagesByTicketID(ctx, ticketID, limit, offset)
}

// AddTicketAttachment adds an attachment to a ticket
func (s *supportService) AddTicketAttachment(ctx context.Context, attachment *shared.TicketAttachment) error {
	if attachment == nil {
		return errors.New("attachment cannot be nil")
	}
	if attachment.ID == uuid.Nil {
		attachment.ID = uuid.New()
	}
	attachment.UploadedAt = time.Now()
	return s.supportRepo.CreateTicketAttachment(ctx, attachment)
}

// GetTicketAttachments gets attachments for a ticket
func (s *supportService) GetTicketAttachments(ctx context.Context, ticketID uuid.UUID) ([]*shared.TicketAttachment, error) {
	return s.supportRepo.GetTicketAttachmentsByTicketID(ctx, ticketID)
}

// CreateFAQ creates a new FAQ
func (s *supportService) CreateFAQ(ctx context.Context, faq *shared.FAQ) error {
	if faq == nil {
		return errors.New("FAQ cannot be nil")
	}
	if faq.ID == uuid.Nil {
		faq.ID = uuid.New()
	}
	return s.supportRepo.CreateFAQ(ctx, faq)
}

// GetFAQByID retrieves an FAQ by ID
func (s *supportService) GetFAQByID(ctx context.Context, id uuid.UUID) (*shared.FAQ, error) {
	return s.supportRepo.GetFAQByID(ctx, id)
}

// GetFAQsByCategory gets FAQs by category
func (s *supportService) GetFAQsByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.FAQ, int64, error) {
	return s.supportRepo.GetFAQsByCategory(ctx, category, limit, offset)
}

// SearchFAQs searches for FAQs
func (s *supportService) SearchFAQs(ctx context.Context, query string, limit, offset int) ([]*shared.FAQ, int64, error) {
	return s.supportRepo.SearchFAQs(ctx, query, limit, offset)
}

// UpdateFAQ updates an FAQ
func (s *supportService) UpdateFAQ(ctx context.Context, faq *shared.FAQ) error {
	if faq == nil {
		return errors.New("FAQ cannot be nil")
	}
	return s.supportRepo.UpdateFAQ(ctx, faq)
}

// DeleteFAQ deletes an FAQ
func (s *supportService) DeleteFAQ(ctx context.Context, id uuid.UUID) error {
	return s.supportRepo.DeleteFAQ(ctx, id)
}

// CreateKnowledgeBaseArticle creates a knowledge base article
func (s *supportService) CreateKnowledgeBaseArticle(ctx context.Context, article *shared.KnowledgeBase) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}
	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}
	return s.supportRepo.CreateKnowledgeBaseArticle(ctx, article)
}

// GetKnowledgeBaseArticleByID retrieves an article by ID
func (s *supportService) GetKnowledgeBaseArticleByID(ctx context.Context, id uuid.UUID) (*shared.KnowledgeBase, error) {
	return s.supportRepo.GetKnowledgeBaseArticleByID(ctx, id)
}

// GetKnowledgeBaseArticleBySlug retrieves an article by slug
func (s *supportService) GetKnowledgeBaseArticleBySlug(ctx context.Context, slug string) (*shared.KnowledgeBase, error) {
	return s.supportRepo.GetKnowledgeBaseArticleBySlug(ctx, slug)
}

// GetKnowledgeBaseArticlesByCategory gets articles by category
func (s *supportService) GetKnowledgeBaseArticlesByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.KnowledgeBase, int64, error) {
	return s.supportRepo.GetKnowledgeBaseArticlesByCategory(ctx, category, limit, offset)
}

// SearchKnowledgeBase searches knowledge base
func (s *supportService) SearchKnowledgeBase(ctx context.Context, query string, limit, offset int) ([]*shared.KnowledgeBase, int64, error) {
	return s.supportRepo.SearchKnowledgeBase(ctx, query, limit, offset)
}

// UpdateKnowledgeBaseArticle updates a knowledge base article
func (s *supportService) UpdateKnowledgeBaseArticle(ctx context.Context, article *shared.KnowledgeBase) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}
	return s.supportRepo.UpdateKnowledgeBaseArticle(ctx, article)
}

// DeleteKnowledgeBaseArticle deletes a knowledge base article
func (s *supportService) DeleteKnowledgeBaseArticle(ctx context.Context, id uuid.UUID) error {
	return s.supportRepo.DeleteKnowledgeBaseArticle(ctx, id)
}

// CreateLiveChat creates a live chat session
func (s *supportService) CreateLiveChat(ctx context.Context, chat *shared.LiveChat) error {
	if chat == nil {
		return errors.New("live chat cannot be nil")
	}
	if chat.ID == uuid.Nil {
		chat.ID = uuid.New()
	}
	if chat.SessionID == "" {
		chat.SessionID = fmt.Sprintf("CHAT-%s", uuid.New().String()[:12])
	}
	chat.StartedAt = time.Now()
	return s.supportRepo.CreateLiveChat(ctx, chat)
}

// GetLiveChatByID retrieves a live chat by ID
func (s *supportService) GetLiveChatByID(ctx context.Context, id uuid.UUID) (*shared.LiveChat, error) {
	return s.supportRepo.GetLiveChatByID(ctx, id)
}

// GetLiveChatsByUserID gets live chats for a user
func (s *supportService) GetLiveChatsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.LiveChat, int64, error) {
	return s.supportRepo.GetLiveChatsByUserID(ctx, userID, limit, offset)
}

// AssignChatToAgent assigns a chat to an agent
func (s *supportService) AssignChatToAgent(ctx context.Context, chatID uuid.UUID, agentID uuid.UUID) error {
	chat, err := s.supportRepo.GetLiveChatByID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to get live chat: %w", err)
	}
	if chat == nil {
		return errors.New("live chat not found")
	}
	now := time.Now()
	chat.AgentID = &agentID
	chat.AssignedAt = &now
	chat.Status = "connected"
	if chat.ConnectedAt == nil {
		chat.ConnectedAt = &now
	}
	return s.supportRepo.UpdateLiveChat(ctx, chat)
}

// EndLiveChat ends a live chat session
func (s *supportService) EndLiveChat(ctx context.Context, chatID uuid.UUID, feedback string, rating int) error {
	chat, err := s.supportRepo.GetLiveChatByID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to get live chat: %w", err)
	}
	if chat == nil {
		return errors.New("live chat not found")
	}
	now := time.Now()
	chat.Status = "ended"
	chat.EndedAt = &now
	chat.Rating = rating
	chat.Feedback = feedback
	return s.supportRepo.UpdateLiveChat(ctx, chat)
}

// AddChatMessage adds a message to a live chat
func (s *supportService) AddChatMessage(ctx context.Context, message *shared.ChatMessage) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}
	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}
	return s.supportRepo.CreateChatMessage(ctx, message)
}

// GetChatMessages gets messages for a live chat
func (s *supportService) GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*shared.ChatMessage, int64, error) {
	return s.supportRepo.GetChatMessagesByChatID(ctx, chatID, limit, offset)
}

// GetSupportStatistics gets support statistics
func (s *supportService) GetSupportStatistics(ctx context.Context, agentID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}
	return s.supportRepo.GetSupportStatistics(ctx, agentID, start, end)
}

// GetTicketResponseTime gets ticket response time
func (s *supportService) GetTicketResponseTime(ctx context.Context, ticketID uuid.UUID) (int, error) {
	duration, err := s.supportRepo.GetTicketResponseTime(ctx, ticketID)
	if err != nil {
		return 0, err
	}
	return int(duration.Minutes()), nil
}

// GetSatisfactionRating gets satisfaction rating
func (s *supportService) GetSatisfactionRating(ctx context.Context, agentID *uuid.UUID, startDate, endDate string) (float64, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date: %w", err)
	}
	return s.supportRepo.GetSatisfactionRating(ctx, agentID, start, end)
}
