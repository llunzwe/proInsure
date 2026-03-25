package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models/shared"
)

// SupportRepository defines the interface for support ticket persistence operations
type SupportRepository interface {
	// Basic CRUD operations
	CreateTicket(ctx context.Context, ticket *shared.SupportTicket) error
	GetTicketByID(ctx context.Context, id uuid.UUID) (*shared.SupportTicket, error)
	GetTicketByTicketNumber(ctx context.Context, ticketNumber string) (*shared.SupportTicket, error)
	UpdateTicket(ctx context.Context, ticket *shared.SupportTicket) error
	DeleteTicket(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetTicketsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByStatus(ctx context.Context, status string, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByPriority(ctx context.Context, priority string, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByAssignedTo(ctx context.Context, agentID uuid.UUID, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetOpenTickets(ctx context.Context, limit, offset int) ([]*shared.SupportTicket, int64, error)
	GetTicketsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*shared.SupportTicket, int64, error)
	SearchTickets(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*shared.SupportTicket, int64, error)

	// Ticket message operations
	CreateTicketMessage(ctx context.Context, message *shared.TicketMessage) error
	GetTicketMessagesByTicketID(ctx context.Context, ticketID uuid.UUID, limit, offset int) ([]*shared.TicketMessage, int64, error)
	UpdateTicketMessage(ctx context.Context, message *shared.TicketMessage) error

	// Ticket attachment operations
	CreateTicketAttachment(ctx context.Context, attachment *shared.TicketAttachment) error
	GetTicketAttachmentsByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*shared.TicketAttachment, error)
	GetTicketAttachmentByID(ctx context.Context, id uuid.UUID) (*shared.TicketAttachment, error)
	DeleteTicketAttachment(ctx context.Context, id uuid.UUID) error

	// FAQ operations
	CreateFAQ(ctx context.Context, faq *shared.FAQ) error
	GetFAQByID(ctx context.Context, id uuid.UUID) (*shared.FAQ, error)
	GetFAQsByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.FAQ, int64, error)
	SearchFAQs(ctx context.Context, query string, limit, offset int) ([]*shared.FAQ, int64, error)
	UpdateFAQ(ctx context.Context, faq *shared.FAQ) error
	DeleteFAQ(ctx context.Context, id uuid.UUID) error

	// Knowledge base operations
	CreateKnowledgeBaseArticle(ctx context.Context, article *shared.KnowledgeBase) error
	GetKnowledgeBaseArticleByID(ctx context.Context, id uuid.UUID) (*shared.KnowledgeBase, error)
	GetKnowledgeBaseArticleBySlug(ctx context.Context, slug string) (*shared.KnowledgeBase, error)
	GetKnowledgeBaseArticlesByCategory(ctx context.Context, category string, limit, offset int) ([]*shared.KnowledgeBase, int64, error)
	SearchKnowledgeBase(ctx context.Context, query string, limit, offset int) ([]*shared.KnowledgeBase, int64, error)
	UpdateKnowledgeBaseArticle(ctx context.Context, article *shared.KnowledgeBase) error
	DeleteKnowledgeBaseArticle(ctx context.Context, id uuid.UUID) error

	// Live chat operations
	CreateLiveChat(ctx context.Context, chat *shared.LiveChat) error
	GetLiveChatByID(ctx context.Context, id uuid.UUID) (*shared.LiveChat, error)
	GetLiveChatBySessionID(ctx context.Context, sessionID string) (*shared.LiveChat, error)
	GetLiveChatsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*shared.LiveChat, int64, error)
	GetActiveLiveChats(ctx context.Context, limit int) ([]*shared.LiveChat, error)
	UpdateLiveChat(ctx context.Context, chat *shared.LiveChat) error

	// Chat message operations
	CreateChatMessage(ctx context.Context, message *shared.ChatMessage) error
	GetChatMessagesByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*shared.ChatMessage, int64, error)

	// Statistics
	GetSupportStatistics(ctx context.Context, agentID *uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error)
	GetTicketResponseTime(ctx context.Context, ticketID uuid.UUID) (time.Duration, error)
	GetSatisfactionRating(ctx context.Context, agentID *uuid.UUID, startDate, endDate time.Time) (float64, error)
}
