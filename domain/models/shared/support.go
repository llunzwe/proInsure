package shared

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models"
	"smartsure/internal/domain/types"
)

// SupportTicket represents customer support tickets
type SupportTicket struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TicketNumber string    `gorm:"uniqueIndex;not null" json:"ticket_number"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Ticket Details
	Subject     string `gorm:"not null" json:"subject"`
	Description string `gorm:"type:text" json:"description"`
	Category    string `json:"category"` // billing, claim, policy, technical, general
	SubCategory string `json:"subcategory"`
	Type        string `json:"type"` // question, complaint, request, feedback

	// Status & Priority
	Status   string `json:"status"`   // open, in_progress, waiting_customer, resolved, closed
	Priority string `json:"priority"` // low, medium, high, urgent
	Severity string `json:"severity"` // minor, major, critical

	// Assignment
	AssignedTo   *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	AssignedAt   *time.Time `json:"assigned_at"`
	DepartmentID *uuid.UUID `gorm:"type:uuid" json:"department_id"`
	TeamID       *uuid.UUID `gorm:"type:uuid" json:"team_id"`

	// Related Entities
	PolicyID  *uuid.UUID `gorm:"type:uuid" json:"policy_id"`
	ClaimID   *uuid.UUID `gorm:"type:uuid" json:"claim_id"`
	PaymentID *uuid.UUID `gorm:"type:uuid" json:"payment_id"`

	// Channel & Source
	Channel       string `json:"channel"` // email, phone, chat, app, web
	Source        string `json:"source"`  // customer, agent, system
	ContactMethod string `json:"contact_method"`

	// Resolution
	ResolvedBy      *uuid.UUID `gorm:"type:uuid" json:"resolved_by"`
	ResolvedAt      *time.Time `json:"resolved_at"`
	ResolutionNotes string     `gorm:"type:text" json:"resolution_notes"`
	ResolutionCode  string     `json:"resolution_code"`

	// SLA & Metrics
	SLADeadline     *time.Time `json:"sla_deadline"`
	FirstResponseAt *time.Time `json:"first_response_at"`
	ResponseTime    int        `json:"response_time"`   // in minutes
	ResolutionTime  int        `json:"resolution_time"` // in minutes
	ReopenCount     int        `json:"reopen_count"`

	// Customer Satisfaction
	SatisfactionRating  int    `json:"satisfaction_rating"` // 1-5
	SatisfactionComment string `json:"satisfaction_comment"`

	// Metadata
	Tags         types.JSONArray `gorm:"type:json" json:"tags"`
	CustomFields string          `gorm:"type:json" json:"custom_fields"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User        *models.User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Agent       *models.User       `gorm:"foreignKey:AssignedTo" json:"agent,omitempty"`
	Resolver    *models.User       `gorm:"foreignKey:ResolvedBy" json:"resolver,omitempty"`
	Policy      *models.Policy     `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Claim       *models.Claim      `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	Payment     *models.Payment    `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	Messages    []TicketMessage    `gorm:"foreignKey:TicketID" json:"messages,omitempty"`
	Attachments []TicketAttachment `gorm:"foreignKey:TicketID" json:"attachments,omitempty"`
}

// TicketMessage represents messages in support tickets
type TicketMessage struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TicketID uuid.UUID `gorm:"type:uuid;not null" json:"ticket_id"`
	SenderID uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`

	// Message Content
	Message     string `gorm:"type:text;not null" json:"message"`
	MessageType string `json:"message_type"` // text, system, note
	IsInternal  bool   `json:"is_internal"`  // Internal notes not visible to customer

	// Status
	Status string     `json:"status"` // sent, delivered, read
	ReadAt *time.Time `json:"read_at"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Ticket *SupportTicket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	Sender *models.User   `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

// TicketAttachment represents files attached to tickets
type TicketAttachment struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TicketID  uuid.UUID  `gorm:"type:uuid;not null" json:"ticket_id"`
	MessageID *uuid.UUID `gorm:"type:uuid" json:"message_id"`

	// File Details
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`

	UploadedBy uuid.UUID `gorm:"type:uuid" json:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at"`

	// Relationships
	Ticket   *SupportTicket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	Message  *TicketMessage `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	Uploader *models.User   `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}

// FAQ represents frequently asked questions
type FAQ struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Category    string    `gorm:"not null" json:"category"`
	SubCategory string    `json:"subcategory"`
	Question    string    `gorm:"not null" json:"question"`
	Answer      string    `gorm:"type:text;not null" json:"answer"`

	// Content
	ShortAnswer    string `json:"short_answer"`
	DetailedAnswer string `gorm:"type:text" json:"detailed_answer"`
	VideoURL       string `json:"video_url"`
	ImageURL       string `json:"image_url"`

	// Metadata
	Tags        types.JSONArray `gorm:"type:json" json:"tags"`
	Keywords    types.JSONArray `gorm:"type:json" json:"keywords"`
	RelatedFAQs types.JSONArray `gorm:"type:json" json:"related_faqs"`

	// Usage & Feedback
	ViewCount       int `json:"view_count"`
	HelpfulCount    int `json:"helpful_count"`
	NotHelpfulCount int `json:"not_helpful_count"`

	// Status
	IsActive     bool `json:"is_active"`
	IsFeatured   bool `json:"is_featured"`
	DisplayOrder int  `json:"display_order"`

	// Multilingual
	Language     string `json:"language"`
	Translations string `gorm:"type:json" json:"translations"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// KnowledgeBase represents help articles and guides
type KnowledgeBase struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ArticleCode string    `gorm:"uniqueIndex;not null" json:"article_code"`
	Title       string    `gorm:"not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`

	// Content
	Summary       string `json:"summary"`
	Content       string `gorm:"type:text;not null" json:"content"`
	ContentFormat string `json:"content_format"` // markdown, html, plain

	// Categorization
	Category    string `json:"category"`
	SubCategory string `json:"subcategory"`
	Type        string `json:"type"` // guide, tutorial, policy, troubleshooting

	// Metadata
	Author   uuid.UUID       `gorm:"type:uuid" json:"author"`
	Reviewer *uuid.UUID      `gorm:"type:uuid" json:"reviewer"`
	Tags     types.JSONArray `gorm:"type:json" json:"tags"`
	Keywords types.JSONArray `gorm:"type:json" json:"keywords"`

	// Media
	FeaturedImage string          `json:"featured_image"`
	Attachments   types.JSONArray `gorm:"type:json" json:"attachments"`
	Videos        types.JSONArray `gorm:"type:json" json:"videos"`

	// Usage & Feedback
	ViewCount     int     `json:"view_count"`
	LikeCount     int     `json:"like_count"`
	ShareCount    int     `json:"share_count"`
	AverageRating float64 `json:"average_rating"`
	RatingCount   int     `json:"rating_count"`

	// Status
	Status      string     `json:"status"` // draft, published, archived
	IsPublished bool       `json:"is_published"`
	IsFeatured  bool       `json:"is_featured"`
	PublishedAt *time.Time `json:"published_at"`

	// Version Control
	Version        string    `json:"version"`
	LastModifiedBy uuid.UUID `gorm:"type:uuid" json:"last_modified_by"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	AuthorUser   *models.User `gorm:"foreignKey:Author" json:"author_user,omitempty"`
	ReviewerUser *models.User `gorm:"foreignKey:Reviewer" json:"reviewer_user,omitempty"`
	ModifierUser *models.User `gorm:"foreignKey:LastModifiedBy" json:"modifier_user,omitempty"`
}

// LiveChat represents live chat sessions
type LiveChat struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID string    `gorm:"uniqueIndex;not null" json:"session_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Chat Details
	Status     string `json:"status"`  // waiting, connected, on_hold, ended
	Channel    string `json:"channel"` // web, mobile
	Department string `json:"department"`

	// Agent Assignment
	AgentID         *uuid.UUID `gorm:"type:uuid" json:"agent_id"`
	AssignedAt      *time.Time `json:"assigned_at"`
	TransferredFrom *uuid.UUID `gorm:"type:uuid" json:"transferred_from"`
	TransferredTo   *uuid.UUID `gorm:"type:uuid" json:"transferred_to"`

	// Timing
	StartedAt    time.Time  `json:"started_at"`
	ConnectedAt  *time.Time `json:"connected_at"`
	EndedAt      *time.Time `json:"ended_at"`
	WaitTime     int        `json:"wait_time"`     // seconds
	ChatDuration int        `json:"chat_duration"` // seconds

	// Chat Quality
	Rating        int    `json:"rating"`
	Feedback      string `json:"feedback"`
	TranscriptURL string `json:"transcript_url"`

	// Context
	InitialMessage string `json:"initial_message"`
	PageURL        string `json:"page_url"`
	UserAgent      string `json:"user_agent"`
	IPAddress      string `json:"ip_address"`

	// Resolution
	IsResolved       bool   `json:"is_resolved"`
	ResolutionNotes  string `json:"resolution_notes"`
	FollowUpRequired bool   `json:"follow_up_required"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User     *models.User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Agent    *models.User  `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	Messages []ChatMessage `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
}

// ChatMessage represents messages in live chat
type ChatMessage struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ChatID   uuid.UUID `gorm:"type:uuid;not null" json:"chat_id"`
	SenderID uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`

	// Message Content
	Message     string `gorm:"type:text;not null" json:"message"`
	MessageType string `json:"message_type"` // text, image, file, system

	// Delivery
	IsDelivered bool       `json:"is_delivered"`
	IsRead      bool       `json:"is_read"`
	DeliveredAt *time.Time `json:"delivered_at"`
	ReadAt      *time.Time `json:"read_at"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Chat   *LiveChat    `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	Sender *models.User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
