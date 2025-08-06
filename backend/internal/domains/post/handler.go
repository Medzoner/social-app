package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-app/api/post"
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

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post
// @Tags Post
// @Accept json
// @Produce json
// @Param post body post.CreatePostInput true "Create Post Input"
// @Success 200 {object} models.Post
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [post]
// @Security BearerAuth
func (h Handler) CreatePost(c *middleware.Context) {
	input := post.CreatePostInput{
		UserID: c.User.ID,
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.usecase.CreatePost(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPosts godoc
// @Summary Get posts by user ID
// @Description Get posts by user ID
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param cursor query string false "Cursor for pagination"
// @Success 200 {array} models.Post
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id}/posts [get]
// @Security BearerAuth
func (h Handler) GetPosts(c *middleware.Context) {
	posts, err := h.usecase.GetPosts(c.Request.Context(), c.Query("cursor"), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
