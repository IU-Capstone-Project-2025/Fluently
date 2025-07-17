package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// LLMMessage represents a single chat message sent to the AI service.
type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// llmChatRequest mirrors the request payload expected by /chat endpoint.
// Field names are lowercase to keep the wrapper internal.
// Optional fields are represented with pointer types so zero values are omitted.
type llmChatRequest struct {
	Messages    []LLMMessage `json:"messages"`
	ModelType   string       `json:"model_type,omitempty"` // "fast" or "balanced"
	MaxTokens   *int         `json:"max_tokens,omitempty"`
	Temperature *float64     `json:"temperature,omitempty"`
}

// LLMChatResponse represents the successful response schema from the /chat endpoint.
type LLMChatResponse struct {
	Response  string `json:"response"`
	ModelUsed string `json:"model_used,omitempty"`
}

// LLMError captures the default FastAPI error response shape.
type LLMError struct {
	Detail string `json:"detail"`
}

func (e LLMError) Error() string {
	return fmt.Sprintf("llm API error: %s", e.Detail)
}

// LLMClientConfig holds configuration for the client.
// If BaseURL is empty, THESAURUS_API_URL is used? Wait that's for Thesaurus. We'll use LLM_API_URL env variable.
type LLMClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

// LLMClient is a lightweight wrapper around the Fluently LLM API.
type LLMClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewLLMClient constructs a new LLMClient with sane defaults.
// BaseURL defaults to http://localhost:8003 or LLM_API_URL env var.
func NewLLMClient(config LLMClientConfig) *LLMClient {
	if config.BaseURL == "" {
		if envURL := os.Getenv("LLM_API_URL"); envURL != "" {
			config.BaseURL = envURL
		} else {
			config.BaseURL = "http://localhost:8003"
		}
	}

	if config.Timeout == 0 {
		config.Timeout = 15 * time.Second
	}

	return &LLMClient{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		timeout: config.Timeout,
	}
}

// Chat sends a list of messages to the LLM service and returns the AI-generated reply.
// modelType can be "fast" or "balanced". If empty, "balanced" is used.
// maxTokens or temperature can be nil to omit the field.
func (c *LLMClient) Chat(ctx context.Context, messages []LLMMessage, modelType string, maxTokens *int, temperature *float64) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("messages cannot be empty")
	}

	if modelType == "" {
		modelType = "balanced"
	}

	reqPayload := llmChatRequest{
		Messages:    messages,
		ModelType:   modelType,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr LLMError
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.Detail != "" {
			return "", apiErr
		}
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp LLMChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return chatResp.Response, nil
}

// ChatSimple is a helper that wraps the simple endpoint /chat/simple.
// This requires only a single user message and returns the model response.
func (c *LLMClient) ChatSimple(ctx context.Context, message string, modelType string) (string, error) {
	if message == "" {
		return "", fmt.Errorf("message cannot be empty")
	}

	if modelType == "" {
		modelType = "balanced"
	}

	url := fmt.Sprintf("%s/chat/simple?message=%s&model_type=%s", c.baseURL, urlQueryEscape(message), urlQueryEscape(modelType))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr LLMError
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.Detail != "" {
			return "", apiErr
		}
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// The simple endpoint returns {"response":"..."}
	var data struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data.Response, nil
}

// HealthCheck pings GET /health and returns true if status is healthy.
func (c *LLMClient) HealthCheck(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// ChatWithDefaults is a convenience wrapper that builds a default client and calls Chat.
func ChatWithDefaults(ctx context.Context, messages []LLMMessage) (string, error) {
	client := NewLLMClient(LLMClientConfig{})
	return client.Chat(ctx, messages, "balanced", nil, nil)
}

// urlQueryEscape is a tiny helper to safely escape query parameters.
func urlQueryEscape(s string) string {
	return (&url.URL{Path: s}).EscapedPath()
}
