package notification

import (
	"context"
	"fmt"

	"social-app/internal/connector"
	"social-app/internal/models"
)

type NotificationRepository interface {
	SendNotification(ctx context.Context, n models.Notification) error
	GetNotifications(ctx context.Context, userID uint64) ([]models.Notification, error)
	MarkAllRead(ctx context.Context, userID uint64) error
	MarkRead(ctx context.Context, id, userID uint64) error
}

type Repository struct {
	conn connector.DBConn
}

func NewRepository(con connector.DBConn) Repository {
	return Repository{
		conn: con,
	}
}

func (r Repository) SendNotification(ctx context.Context, notif models.Notification) error {
	r.conn.DB.Create(&notif)
	if err := r.conn.DB.Error; err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}

func (r Repository) GetNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	var notifications []models.Notification
	r.conn.DB.
		Where("user_id", userID).
		Where("is_read = false").
		Find(&notifications)
	if err := r.conn.DB.Error; err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	if len(notifications) <= 10 {
		return notifications, nil
	}

	// clean notifications
	err := r.clean(notifications)
	if err != nil {
		return nil, fmt.Errorf("failed to clean notifications: %w", err)
	}

	return notifications, nil
}

func (r Repository) clean(notifications []models.Notification) error {
	for i := range notifications {
		if err := r.conn.DB.Unscoped().Delete(&notifications[i]).Error; err != nil {
			return fmt.Errorf("failed to hard delete notification ID %d: %w", notifications[i].ID, err)
		}
	}
	return nil
}

func (r Repository) MarkAllRead(ctx context.Context, userID uint64) error {
	r.conn.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true)
	if err := r.conn.DB.Error; err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

func (r Repository) MarkRead(ctx context.Context, id, userID uint64) error {
	r.conn.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	if err := r.conn.DB.Error; err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}
