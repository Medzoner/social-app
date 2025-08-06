package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"social-app/internal/models"
	"social-app/pkg/middleware"
)

type Handler struct {
	conn *Connector
}

func NewHandler(c *Connector) Handler {
	return Handler{
		conn: c,
	}
}

func (h Handler) HandleWebSocket(c *middleware.Context) {
	userID := c.User.ID

	conn, _, err := h.conn.ConnectClient(c.Request, c.Writer, userID)
	if err != nil {
		log.Printf("[WS] Failed to connect user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket failed"})
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	defer func() {
		err := h.conn.CloseClient(c.Request.Context(), userID)
		if err != nil {
			log.Printf("[WS] Error closing connection for user %d: %v", userID, err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[WS] Context done, closing connection for user %d", userID)
			return
		default:
			if conn == nil {
				log.Printf("[WS] Connection for user %d is nil, reconnecting...", userID)
				var sID string
				conn, sID, err = h.conn.ConnectClient(c.Request, c.Writer, userID)
				if err != nil {
					log.Printf("[WS] Reconnect error for user %d: %v", userID, err)
					break
				}
				log.Printf("[WS] User %d reconnected with socket ID %s", userID, sID)
			}

			if err := h.readAndDispatch(userID); err != nil {
				log.Printf("[WS] Read error for user %d: %v", userID, err)
				return
			}
		}
	}
}

func (h Handler) readAndDispatch(senderID uint64) error {
	_, raw, err := h.conn.Clients[senderID].ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
			return nil
		}

		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return fmt.Errorf("connection closed by user %d: %w", senderID, err)
		}

		return fmt.Errorf("read error: %w", err)
	}

	var msg WSMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		log.Printf("[WS] Invalid JSON from user %d: %s", senderID, string(raw))
		return fmt.Errorf("invalid JSON: %w", err)
	}

	msg.SenderID = senderID

	return h.dispatchMessage(msg)
}

func (h Handler) dispatchMessage(msg WSMessage) error {
	switch msg.Type {
	case TypeTyping, TypeSeen, TypeMessage:
		return h.handleSignalMessage(msg)
	default:
		log.Printf("[WS] Unhandled message type '%s' from user %d", msg.Type, msg.SenderID)
		return nil
	}
}

func (h Handler) handleSignalMessage(msg WSMessage) error {
	if msg.ReceiverID == 0 {
		log.Printf("[WS] Missing receiver_id in '%s' message from user %d", msg.Type, msg.SenderID)
		return nil
	}

	notification := models.Notification{
		UserID:     msg.SenderID,
		Type:       models.NotificationType(msg.Type),
		ReceiverID: msg.ReceiverID,
		Content:    fmt.Sprintf("User %d sent signal: %s", msg.SenderID, msg.Type),
	}

	return h.conn.Send(msg.ReceiverID, notification)
}
