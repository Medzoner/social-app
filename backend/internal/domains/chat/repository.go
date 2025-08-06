package chat

import (
	"context"
	"fmt"

	"social-app/internal/connector"
	"social-app/internal/models"
	"social-app/pkg/pagination"
)

const limit = 10

type ChatRepository interface {
	CreateMessage(ctx context.Context, msg models.Message) (models.Message, error)
	GetChatList(ctx context.Context, senderID uint64) ([]models.Message, error)
	GetMessages(ctx context.Context, cursor string, userID, otherID uint64) (models.MessageList, error)
	MarkMessagesAsRead(ctx context.Context, userID, otherID uint64) (int64, error)
}

type Repository struct {
	conn connector.DBConn
}

func NewRepository(con connector.DBConn) Repository {
	return Repository{
		conn: con,
	}
}

func (r Repository) CreateMessage(ctx context.Context, c models.Message) (models.Message, error) {
	r.conn.DB.Create(&c)
	if r.conn.DB.Error != nil {
		return models.Message{}, fmt.Errorf("error creating message: %w", r.conn.DB.Error)
	}

	return c, nil
}

func (r Repository) GetChatList(ctx context.Context, userID uint64) ([]models.Message, error) {
	var messages []models.Message

	subQuery := r.conn.DB.
		Table("messages").
		Select("MAX(id)").
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Group("LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id)")

	if err := r.conn.DB.
		Preload("User").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch chat list: %w", err)
	}

	for i := range messages {
		var otherUserID uint64
		if messages[i].SenderID == userID {
			otherUserID = messages[i].ReceiverID
		} else {
			otherUserID = messages[i].SenderID
		}

		var otherUser models.User
		if err := r.conn.DB.Select("id", "username", "avatar").First(&otherUser, otherUserID).Error; err == nil {
			messages[i].User = otherUser
		}
	}

	return messages, nil
}

func (r Repository) GetMessages(ctx context.Context, cursor string, userID, otherID uint64) (models.MessageList, error) {
	var messages []models.Message
	db := r.conn.DB.
		Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID,
			otherID,
			otherID,
			userID).
		Limit(limit + 1)

	db = pagination.CursorFilter[models.Message](cursor, db, "desc")

	if err := db.Find(&messages).Error; err != nil {
		return models.MessageList{}, fmt.Errorf("error fetching messages: %w", err)
	}

	nextCursor, hasMore, items := pagination.NextCursor(messages, limit)

	return models.MessageList{
		Messages:   items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r Repository) MarkMessagesAsRead(ctx context.Context, userID, otherID uint64) (int64, error) {
	result := r.conn.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND read = false", otherID, userID).
		Update("read", true)

	if result.Error != nil {
		return 0, fmt.Errorf("error marking messages as read: %w", result.Error)
	}

	result = r.conn.DB.Model(&models.Notification{}).
		Where("type = 'message' AND user_id = ? AND type = ?", userID, models.NotificationTypeMessage).
		Update("is_read", true)

	if result.Error != nil {
		return 0, fmt.Errorf("error marking notifications as read: %w", result.Error)
	}

	return result.RowsAffected, nil
}
