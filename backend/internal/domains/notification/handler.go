package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-app/api/notification"
	"social-app/pkg/middleware"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(uc UseCase) Handler {
	return Handler{
		usecase: uc,
	}
}

// func (h Handler) NotificationSocket(c *middleware.Context) {
//	conn, err := ws.ConnectClient(c.Request, c.Writer, c.User.ID)
//	if err != nil {
//		fmt.Println("Broadcaster upgrade failed:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to Broadcaster"})
//		return
//	}
//
//	defer func() {
//		err := ws.CloseClient(c.User.ID)
//		if err != nil {
//			fmt.Println("Error closing Broadcaster connection:", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close Broadcaster connection"})
//		}
//	}()
//
//	for {
//		t, p, err := conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("Unexpected close error: %v", err)
//			} else {
//				log.Printf("WebSocket closed normally: %v", err)
//			}
//			break
//		}
//		fmt.Println("Received message:", string(p), "Type:", t)
//	}
//}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Fetch notifications for the authenticated user
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200 {array} models.Notification
// @Failure 500 {object} map[string]string
// @Router /notifications [get]
// @Security BearerAuth
func (h Handler) GetNotifications(c *middleware.Context) {
	notifications, err := h.usecase.GetNotifications(c.Request.Context(), c.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// MarkAllRead godoc
// @Summary Mark all notifications as read
// @Description Mark all notifications for the authenticated user as read
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/mark-all-read [post]
// @Security BearerAuth
func (h Handler) MarkAllRead(c *middleware.Context) {
	err := h.usecase.MarkAllRead(c.Request.Context(), notification.MarkReadInput{
		UserID: c.User.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// MarkRead godoc
// @Summary Mark specific notifications as read
// @Description Mark specific notifications for the authenticated user as read
// @Tags Notifications
// @Accept json
// @Produce json
// @Param input body notification.MarkReadInput true "Input for marking notifications as read"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/mark-read [post]
// @Security BearerAuth
func (h Handler) MarkRead(c *middleware.Context) {
	input := notification.MarkReadInput{
		UserID: c.User.ID,
	}
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.usecase.MarkRead(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
