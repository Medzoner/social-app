package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"social-app/api/comment"
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

// CreateComment godoc
// @Summary Create a comment
// @Description Create a comment for a post
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param comment body comment.CreateCommentInput true "Comment data"
// @Success 200 {object} models.Comment "Comment created successfully"
// @Failure 400 {object} map[string]string "Invalid post ID or comment binding"
// @Failure 500 {object} map[string]string "Failed to create comment"
// @Router /posts/{id}/comments [post]
// @Security BearerAuth
func (h Handler) CreateComment(c *middleware.Context) {
	postID, err := c.GetUint64("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var com comment.CreateCommentInput
	if err := c.BindJSON(&c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment binding"})
	}
	com.UserID = c.User.ID
	com.PostID = postID

	r, err := h.usecase.CreateComment(c.Request.Context(), com)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, r)
}

// GetComments godoc
// @Summary Get comments for a post
// @Description Get comments for a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param cursor query string false "Cursor for pagination"
// @Success 200 {array} models.Comment "List of comments"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 500 {object} map[string]string "Failed to get comments"
// @Router /posts/{id}/comments [get]
// @Security BearerAuth
func (h Handler) GetComments(c *middleware.Context) {
	postID, err := c.GetUint64("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cursor := c.Query("cursor")
	cm, err := h.usecase.GetComments(c.Request.Context(), cursor, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, cm)
}

// CountByPosts godoc
// @Summary Count comments for posts
// @Description Count the number of comments for specific posts
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]int "Count of comments"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 500 {object} map[string]string "Failed to count comments"
// @Router /posts/{id}/comments/count [get]
// @Security BearerAuth
func (h Handler) CountByPosts(c *middleware.Context) {
	postIDs := make([]uint64, 0)
	for _, idStr := range c.QueryArray("ids") {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID: " + idStr})
			return
		}
		postIDs = append(postIDs, id)
	}

	cc, err := h.usecase.CountByPosts(c.Request.Context(), postIDs, c.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"counts": cc,
	})
}
