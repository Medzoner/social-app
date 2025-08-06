package chat

import (
	"context"
	"fmt"

	"social-app/api/chat"
	"social-app/internal/models"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(r Repository) UseCase {
	return UseCase{
		repo: r,
	}
}

func (u UseCase) CreateMessage(ctx context.Context, msg chat.CreateMessageInput) (models.Message, error) {
	newMsg := models.Message{
		UserID:     msg.UserID,
		SenderID:   msg.UserID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
	}

	m, err := u.repo.CreateMessage(ctx, newMsg)
	if err != nil {
		return models.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	return m, nil
}

func (u UseCase) GetChatList(ctx context.Context, userID uint64) ([]models.Message, error) {
	m, err := u.repo.GetChatList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat list: %w", err)
	}

	return m, nil
}

func (u UseCase) GetMessages(ctx context.Context, cursor string, userID, otherID uint64) (models.MessageList, error) {
	m, err := u.repo.GetMessages(ctx, cursor, userID, otherID)
	if err != nil {
		return models.MessageList{}, fmt.Errorf("failed to get messages: %w", err)
	}
	return m, nil
}

func (u UseCase) MarkMessagesAsRead(ctx context.Context, userID, otherID uint64) (int64, error) {
	result, err := u.repo.MarkMessagesAsRead(ctx, userID, otherID)
	if err != nil {
		return 0, fmt.Errorf("failed to mark messages as read: %w", err)
	}
	return result, nil
}
