package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// SupportService defines the interface for support ticket business logic operations
type SupportService interface {
	// Ticket management
	CreateTicket(ctx context.Context, ticket *shared.SupportTicket) error
	GetTicketByID(ctx context.Context, id uuid.UUID) (*shared.SupportTicket, error)
	GetTicketByTicketNumber(ctx context.Context, ticketNumber string) (*shared.SupportTicket, error)
	UpdateTicket(ctx context.Context, ticket *shared.SupportTicket) error
	CloseTicket(ctx context.Context, ticketID uuid.UUID, resolutionNotes string) error

	// Ticket queries
	GetTicketsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByPriority(ctx context.Context, priority string, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByAssignedTo(ctx context.Context, agentID uuid.UUID, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetOpenTickets(ctx context.Context, limit, offset int) ([]*shared.SupportTicket, int64, error)
	SearchTickets(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*shared.SupportTicket, int64, error)

	// Ticket assignment
	AssignTicket(ctx context.Context, ticketID uuid.UUID, agentID uuid.UUID) error
	ReassignTicket(ctx context.Context, ticketID uuid.UUID, newAgentID uuid.UUID, reason string) error
	UnassignTicket(ctx context.Context, ticketID uuid.UUID) error

	// Ticket messages
	AddTicketMessage(ctx context.Context, message *shared.TicketMessage) error
	GetTicketMessages(ctx context.Context, ticketID uuid.UUID, limit, offset int) ([]*shared.TicketMessage, int64, error)

	// Ticket attachments
	AddTicketAttachment(ctx context.Context, attachment *shared.TicketAttachment) error
	GetTicketAttachments(ctx context.Context, ticketID uuid.UUID) ([]*shared.TicketAttachment, error)

	// FAQ management
	CreateFAQ(ctx context.Context, faq *shared.FAQ) error
	GetFAQByID(ctx context.Context, id uuid.UUID) (*shared.FAQ, error)
	GetFAQsByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.FAQ, int64, error)
	SearchFAQs(ctx context.Context, query string, limit, offset int) ([]*shared.FAQ, int64, error)
	UpdateFAQ(ctx context.Context, faq *shared.FAQ) error
	DeleteFAQ(ctx context.Context, id uuid.UUID) error

	// Knowledge base management
	CreateKnowledgeBaseArticle(ctx context.Context, article *shared.KnowledgeBase) error
	GetKnowledgeBaseArticleByID(ctx context.Context, id uuid.UUID) (*shared.KnowledgeBase, error)
	GetKnowledgeBaseArticleBySlug(ctx context.Context, slug string) (*shared.KnowledgeBase, error)
	GetKnowledgeBaseArticlesByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.KnowledgeBase, int64, error)
	SearchKnowledgeBase(ctx context.Context, query string, limit, offset int) ([]*shared.KnowledgeBase, int64, error)
	UpdateKnowledgeBaseArticle(ctx context.Context, article *shared.KnowledgeBase) error
	DeleteKnowledgeBaseArticle(ctx context.Context, id uuid.UUID) error

	// Live chat management
	CreateLiveChat(ctx context.Context, chat *shared.LiveChat) error
	GetLiveChatByID(ctx context.Context, id uuid.UUID) (*shared.LiveChat, error)
	GetLiveChatsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.LiveChat, int64, error)
	AssignChatToAgent(ctx context.Context, chatID uuid.UUID, agentID uuid.UUID) error
	EndLiveChat(ctx context.Context, chatID uuid.UUID, feedback string, rating int) error

	// Chat messages
	AddChatMessage(ctx context.Context, message *shared.ChatMessage) error
	GetChatMessages(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*shared.ChatMessage, int64, error)

	// Support statistics
	GetSupportStatistics(ctx context.Context, agentID *uuid.UUID, startDate, endDate string) (map[string]interface{}, error)
	GetTicketResponseTime(ctx context.Context, ticketID uuid.UUID) (int, error)
	GetSatisfactionRating(ctx context.Context, agentID *uuid.UUID, startDate, endDate string) (float64, error)
}
