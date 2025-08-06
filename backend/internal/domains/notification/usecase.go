package notification

import (
	"context"
	"fmt"

	"social-app/api/notification"
	"social-app/internal/models"
	"social-app/pkg/ws"
)

type UseCase struct {
	repo Repository
	conn *ws.Connector
}

func NewUseCase(r Repository, c *ws.Connector) UseCase {
	return UseCase{
		repo: r,
		conn: c,
	}
}

func (u UseCase) SendNotification(ctx context.Context, notif models.Notification) error {
	err := u.repo.SendNotification(ctx, notif)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	if err := u.conn.Send(notif.ReceiverID, notif); err != nil {
		return fmt.Errorf("failed to write notification to Broadcaster: %w", err)
	}

	return nil
}

func (u UseCase) GetNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	notifications, err := u.repo.GetNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	return notifications, nil
}

func (u UseCase) MarkAllRead(ctx context.Context, input notification.MarkReadInput) error {
	if err := u.repo.MarkAllRead(ctx, input.UserID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

func (u UseCase) MarkRead(ctx context.Context, input notification.MarkReadInput) error {
	if err := u.repo.MarkRead(ctx, input.ID, input.UserID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}
