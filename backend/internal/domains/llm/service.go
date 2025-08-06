package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"social-app/internal/config"
	http2 "social-app/pkg/http"
)

type Service struct {
	ollamaURL string
}

func NewService(cfg config.LLM) Service {
	return Service{
		ollamaURL: cfg.URL,
	}
}

type ollamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatPayload struct {
	Model    string              `json:"model"`
	Messages []ollamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type ollamaChatResponse struct {
	Message ollamaChatMessage `json:"message"`
}

func (s Service) CallOllama(ctx context.Context, prompt string) (string, error) {
	messages := []ollamaChatMessage{
		{
			Role:    "system",
			Content: "Tu es un assistant qui répond toujours en HTML structuré. Utilise <p>, <ol>, <li>, <strong>, etc. Pour les extraits de code, mets dans la balise code <code> avec le style (sans les quotes obliques). Ne mets pas de <html> ou <body>.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	payload := ollamaChatPayload{
		Model:    "mistral",
		Messages: messages,
		Stream:   false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/chat", s.ollamaURL), bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", http2.ContentTypeJSON)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error: %s", raw)
	}

	var result ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}

	return result.Message.Content, nil
}
