package models

import "time"

type Message struct {
	Content string `json:"content"`
	Model
	User       User   `json:"user"`
	UserID     uint64 `gorm:"index"       json:"user_id"`
	SenderID   uint64 `json:"sender_id"`
	ReceiverID uint64 `json:"receiver_id"`
	Read       bool   `json:"read"`
}

func (m Message) GetCursorFields() (createdAt time.Time, id uint64) {
	return m.CreatedAt, m.ID
}

type MessageList struct {
	NextCursor string    `json:"next_cursor,omitempty"`
	Messages   []Message `json:"messages"`
	HasMore    bool      `json:"has_more"`
}
