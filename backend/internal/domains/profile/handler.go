package profile

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"social-app/api/profile"
	"social-app/internal/domains/media"
	"social-app/pkg/middleware"
)

type Handler struct {
	usecase UseCase
	mUC     media.UseCase
}

func NewHandler(uc UseCase, m media.UseCase) Handler {
	return Handler{
		usecase: uc,
		mUC:     m,
	}
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get user profile by ID
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
// @Security BearerAuth
func (h Handler) GetProfile(c *middleware.Context) {
	id := c.Param("id")
	user, err := h.usecase.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if user.Avatar == "" {
		c.JSON(http.StatusOK, user)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary      Update user's profile
// @Description  Partially updates bio, avatar or email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body UpdateProfileInput true "Profile patch input"
// @Success      200 {object} models.User
// @Failure      400 {object} map[string]string
// @Router       /users/{id} [patch]
// @Security     Bearer

func (h Handler) UpdateProfile(c *middleware.Context) {
	var input profile.UpdateProfileInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	userID := c.User.ID
	user, err := h.usecase.GetProfile(c.Request.Context(), strconv.FormatUint(userID, 10))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Bio != nil {
		user.Bio = *input.Bio
	}
	if input.Avatar != nil {
		user.Avatar = *input.Avatar
	}
	if input.Email != nil {
		user.Email = *input.Email
	}

	updatedUser, err := h.usecase.UpdateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
