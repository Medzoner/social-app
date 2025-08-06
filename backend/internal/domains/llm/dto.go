package llm

type LLMRequest struct {
	Message string `binding:"required" json:"message"`
}

type LLMResponse struct {
	Reply string `json:"reply"`
}
