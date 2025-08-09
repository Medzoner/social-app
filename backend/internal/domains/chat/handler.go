package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"social-app/api/chat"
	notif "social-app/internal/domains/notification"
	"social-app/internal/models"
	"social-app/pkg/middleware"
)

type Handler struct {
	usecase UseCase
	notif   notif.UseCase
}

func NewHandler(uc UseCase, n notif.UseCase) Handler {
	return Handler{
		usecase: uc,
		notif:   n,
	}
}

func (h Handler) CreatMessage(c *middleware.Context) {
	msg := chat.CreateMessageInput{UserID: c.User.ID}
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message binding"})
		return
	}

	newMsg, err := h.usecase.CreateMessage(c.Request.Context(), msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	if err := h.notifyNewMessage(c.Request.Context(), newMsg, msg.ReceiverID); err != nil {
		fmt.Println("Error notifying new message:", err)
	}

	c.JSON(http.StatusOK, newMsg)
}

// GetChatList RegisterRoutes godoc
// @Summary      Register chat routes
// @Description  Register chat routes
// @Tags         chat
// @Accept       json
// @Produce      json
// @Router       /chats [get]
// @Success      200 {array}  []models.Message
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
func (h Handler) GetChatList(c *middleware.Context) {
	results, err := h.usecase.GetChatList(c.Request.Context(), c.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve chat list"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetMessages godoc
// @Summary      Get messages
// @Description  Get messages between two users
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        cursor query string false "Cursor for pagination"
// @Router       /messages/{id} [get]
// @Success      200 {array}  models.MessageList
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
func (h Handler) GetMessages(c *middleware.Context) {
	cursor := c.Query("cursor")

	otherID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if otherID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be greater than 0"})
		return
	}
	items, err := h.usecase.GetMessages(c.Request.Context(), cursor, c.User.ID, uint64(otherID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// MarkMessagesAsRead godoc
// @Summary      Mark messages as read
// @Description  Mark messages as read between two users
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Router       /messages/{id}/read [post]
// @Success      200 {object} map[string]bool
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
func (h Handler) MarkMessagesAsRead(c *middleware.Context) {
	otherID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if otherID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be greater than 0"})
		return
	}

	result, err := h.usecase.MarkMessagesAsRead(c.Request.Context(), c.User.ID, uint64(otherID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not mark messages as read"})
		return
	}

	if err := h.notifySeenMessage(c.Request.Context(), c.User.ID, uint64(otherID)); err != nil {
		fmt.Println("Error notifying new message:", err)
	}

	c.JSON(http.StatusOK, gin.H{"updated": result})
}

func (h Handler) notifySeenMessage(ctx context.Context, userID, receiverID uint64) error {
	notification := models.Notification{
		UserID:     userID,
		Type:       models.NotificationTypeSeen,
		Link:       fmt.Sprintf("/chat/%d", userID),
		ReceiverID: receiverID,
	}

	if err := h.notif.SendNotification(ctx, notification); err != nil {
		return fmt.Errorf("error sending notification: %w", err)
	}

	return nil
}

func (h Handler) notifyNewMessage(ctx context.Context, newMsg models.Message, userID uint64) error {
	pl, err := json.Marshal(newMsg)
	if err != nil {
		return fmt.Errorf("error marshalling message payload: %w", err)
	}
	notification := models.Notification{
		UserID:     newMsg.ReceiverID,
		Type:       models.NotificationTypeMessage,
		Content:    newMsg.Content,
		Link:       fmt.Sprintf("/chat/%d", userID),
		Payload:    string(pl),
		TypeID:     newMsg.ID,
		ReceiverID: newMsg.ReceiverID,
	}

	if err = h.notif.SendNotification(ctx, notification); err != nil {
		return fmt.Errorf("error sending notification: %w", err)
	}

	return nil
}
