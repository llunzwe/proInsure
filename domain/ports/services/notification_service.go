package services

import (
	"context"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
)

// NotificationService defines the interface for notification business operations
type NotificationService interface {
	// Notification Management
	GetUserNotifications(ctx context.Context, userID uuid.UUID, status, notificationType string, limit, offset int) ([]*models.Notification, int, error)
	GetNotificationByID(ctx context.Context, notificationID, userID uuid.UUID) (*models.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID, userID uuid.UUID) error
	DeleteNotification(ctx context.Context, notificationID, userID uuid.UUID) error

	// Notification Preferences
	GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (interface{}, error) // NotificationPreferences type not found
	UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, emailEnabled, smsEnabled, pushEnabled, inAppEnabled *bool, emailPrefs, smsPrefs, pushPrefs map[string]bool, quietStart, quietEnd, timezone *string) error
	SendTestNotification(ctx context.Context, userID uuid.UUID, channel string) error

	// Notification Templates
	CreateNotificationTemplate(ctx context.Context, name, description, templateType, subject, content string, variables map[string]string, isActive bool) (interface{}, error) // NotificationTemplate type not found
	GetNotificationTemplates(ctx context.Context, templateType, isActive string, limit, offset int) ([]interface{}, int, error) // NotificationTemplate type not found
	UpdateNotificationTemplate(ctx context.Context, templateID uuid.UUID, name, description, templateType, subject, content string, variables map[string]string, isActive *bool) error
	DeleteNotificationTemplate(ctx context.Context, templateID uuid.UUID) error

	// Bulk Notifications
	SendBulkNotifications(ctx context.Context, userIDs []uuid.UUID, templateID uuid.UUID, variables map[string]string, channels []string, scheduleAt *string, priority string) (uuid.UUID, error)
}
