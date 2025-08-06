package ws

type MessageType string

const (
	TypeTyping  MessageType = "typing"
	TypeSeen    MessageType = "seen"
	TypeMessage MessageType = "message"
)

type WSMessage struct {
	Payload    map[string]any `json:"payload,omitempty"`
	Type       MessageType    `json:"type"`
	UserID     uint64         `json:"user_id"`
	SenderID   uint64         `json:"sender_id"`
	ReceiverID uint64         `json:"receiver_id"`
}
