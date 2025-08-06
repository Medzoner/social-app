package like

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"social-app/api/like"
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

// LikePost godoc
// @Summary Like a post
// @Description Like a post by ID
// @Tags Likes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param like body like.PostInput true "Like Input"
// @Success 200 {object} models.Like "Like created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 409 {object} map[string]string "Like already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id}/like [post]
// @Security ApiKeyAuth
func (h Handler) LikePost(c *middleware.Context) {
	var l like.PostInput
	if err := c.BindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input binding"})
		return
	}
	l.UserID = c.User.ID
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	if postID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID must be a positive integer"})
		return
	}
	l.PostID = uint64(postID)

	r, err := h.usecase.LikePost(c.Request.Context(), l)
	if err != nil {
		if err.Error() == "like already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "You have already liked this post"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	c.JSON(http.StatusOK, r)
}

// UnlikePost godoc
// @Summary Unlike a post
// @Description Unlike a post by ID
// @Tags Likes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]string "Post unliked successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Like not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id}/unlike [post]
// @Security BearerAuth
func (h Handler) UnlikePost(c *middleware.Context) {
	var l like.PostInput
	l.UserID = c.User.ID

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	if postID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID must be a positive integer"})
		return
	}
	l.PostID = uint64(postID)

	_, err = h.usecase.UnlikePost(c.Request.Context(), l)
	if err != nil {
		if err.Error() == "like not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "You have not liked this post"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

// GetLikes godoc
// @Summary Get likes for a post
// @Description Retrieve all likes for a specific post by ID
// @Tags Likes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} models.Like "List of likes"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id}/likes [get]
// @Security BearerAuth
func (h Handler) GetLikes(c *middleware.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil || postIDInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	l, err := h.usecase.GetLikes(c.Request.Context(), uint64(postIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve likes"})
		return
	}

	c.JSON(http.StatusOK, l)
}

// GetLikeStats godoc
// @Summary Get like statistics for a post
// @Description Retrieve like statistics for a specific post by ID
// @Tags Likes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.LikeStats "Like statistics"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id}/like-stats [get]
// @Security BearerAuth
func (h Handler) GetLikeStats(c *middleware.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil || postIDInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	stats, err := h.usecase.GetLikeStats(c.Request.Context(), []uint64{uint64(postIDInt)}, c.User.ID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve like stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h Handler) GetStats(c *middleware.Context) {
	postIDs := make([]uint64, 0)
	for _, idStr := range c.QueryArray("ids") {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID: " + idStr})
			return
		}
		postIDs = append(postIDs, id)
	}

	likes, err := h.usecase.GetStats(c.Request.Context(), postIDs, c.User.ID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve likes"})
		return
	}

	c.JSON(http.StatusOK, likes)
}
