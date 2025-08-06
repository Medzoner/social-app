package models

type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeMessage NotificationType = "message"
	NotificationTypePost    NotificationType = "post"
	NotificationTypeTyping  NotificationType = "typing"
	NotificationTypeSeen    NotificationType = "seen"
)

type Notification struct {
	Content string           `json:"content,omitempty"`
	Type    NotificationType `json:"type"`
	Link    string           `json:"link,omitempty"`
	Payload string           `json:"payload,omitempty"`
	Model
	User       User   `json:"user,omitempty"`
	UserID     uint64 `gorm:"index"          json:"user_id"`
	TypeID     uint64 `json:"type_id"`
	ReceiverID uint64 `json:"receiver_id"`
	IsRead     bool   `json:"is_read"`
}

type NotificationHistory struct {
	Notification
}
