package chat

type CreateMessageInput struct {
	Content    string `json:"content"`
	ReceiverID uint64 `json:"to"`
	UserID     uint64 `json:"user_id"`
}
