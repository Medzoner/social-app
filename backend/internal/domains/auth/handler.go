package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-app/api/auth"
	"social-app/api/profile"
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

// Register godoc
// @Summary      Register a user
// @Description  Register a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body auth.RegisterInput true "User registration input"
// @Success      200 {object} map[string]string "registration successful"
// @Failure      400 {object} map[string]string "Invalid input binding"
// @Failure      500 {object} map[string]string "Failed to register user"
// @Router       /register [post]
func (h Handler) Register(c *gin.Context) {
	input := auth.RegisterInput{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input binding"})
	}

	err := h.usecase.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successfully"})
}

// Login godoc
// @Summary      Login a user
// @Description  Login an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body auth.LoginInput true "User login input"
// @Success      200 {object} map[string]string "JWT token"
// @Failure      400 {object} map[string]string "Invalid input binding"
// @Failure      500 {object} map[string]string "Failed to login user"
// @Router       /login [post]
func (h Handler) Login(c *gin.Context) {
	input := auth.LoginInput{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input binding"})
	}

	jwt, err := h.usecase.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	if !jwt.VerifyResponse.Verified {
		c.JSON(http.StatusAccepted, gin.H{
			"verified": false,
			"id":       jwt.VerifyResponse.ID,
			"username": jwt.VerifyResponse.Username,
		})
		return
	}

	c.JSON(http.StatusOK, jwt)
}

// Logout godoc
// @Summary      Logout a user
// @Description  Logout the current user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "logged out successfully"
// @Failure      500 {object} map[string]string "Failed to logout user"
// @Router       /logout [post]
// @Security BearerAuth
func (h Handler) Logout(c *middleware.Context) {
	err := h.usecase.Logout(c.Request.Context(), c.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h Handler) RefreshToken(c *middleware.Context) {
	jwt, err := h.usecase.RefreshToken(c.Request.Context(), c.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, jwt)
}

func (h Handler) Verify(c *middleware.Context) {
	var input profile.VerifyCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code invalide"})
		return
	}

	err := h.usecase.Verify(c.Request.Context(), input.UserID, input.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Échec de la vérification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Succesfully verified"})
}

func (h Handler) RequestVerification(c *middleware.Context) {
	var input profile.ResendCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalid"})
		return
	}

	if input.Type == "email" {
		err := h.usecase.SendEmailVerification(c.Request.Context(), input.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l’envoi du code email"})
			return
		}
	}
	if input.Type == "phone" {
		err := h.usecase.SendPhoneVerification(c.Request.Context(), input.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l’envoi du code SMS"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Code email envoyé"})
}

func (h Handler) OauthLogin(context *gin.Context) {
	resp := h.usecase.OauthLogin(context.Request.Context())
	context.Redirect(http.StatusFound, resp)
}

func (h Handler) OauthCallback(context *gin.Context) {
	var input auth.OauthInput
	if err := context.ShouldBindQuery(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth input"})
		return
	}

	jwt, err := h.usecase.OauthCallback(context.Request.Context(), input)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete OAuth login"})
		return
	}

	if !jwt.VerifyResponse.IsEmpty() {
		context.JSON(http.StatusAccepted, gin.H{
			"verified": false,
			"id":       jwt.VerifyResponse.ID,
			"username": jwt.VerifyResponse.Username,
		})
		return
	}

	context.Redirect(http.StatusFound, "http://localhost:5173/oauth-callback?access_token="+jwt.AccessToken+"&refresh_token="+jwt.RefreshToken+"&id_token="+jwt.IDToken)
}

func (h Handler) IsUserOnline(context *middleware.Context) {
	isOnline, err := h.usecase.IsUserOnline(context.Request.Context(), context.User.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user online status"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"is_online": isOnline})
}
