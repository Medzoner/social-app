package llm

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h Handler) Prompt(c *middleware.Context) {
	var input LLMRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}

	reply, err := h.usecase.AskLLM(c.Request.Context(), input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LLMResponse{Reply: reply})
}
