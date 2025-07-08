package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DistractorRequest represents the request payload for the distractor API
type DistractorRequest struct {
	Sentence string `json:"sentence"`
	Word     string `json:"word"`
}

// DistractorResponse represents the response from the distractor API
type DistractorResponse struct {
	PickOptions []string `json:"pick_options"`
}

// DistractorError represents an error response from the distractor API
type DistractorError struct {
	Detail    string `json:"detail"`
	ErrorCode string `json:"error_code,omitempty"`
}

func (e DistractorError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("distractor API error [%s]: %s", e.ErrorCode, e.Detail)
	}
	return fmt.Sprintf("distractor API error: %s", e.Detail)
}

// DistractorClient provides an interface to the FastAPI distractor service
type DistractorClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// DistractorClientConfig holds configuration for the DistractorClient
type DistractorClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

// NewDistractorClient creates a new distractor client with the given configuration
func NewDistractorClient(config DistractorClientConfig) *DistractorClient {
	if config.BaseURL == "" {
		// Check environment variable first
		if envURL := os.Getenv("ML_API_URL"); envURL != "" {
			config.BaseURL = envURL
		} else {
			// Use Docker service name when running in containerized environment
			// Falls back to localhost for local development
			config.BaseURL = "http://localhost:8001"
		}
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &DistractorClient{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		timeout: config.Timeout,
	}
}

// GenerateDistractors calls the distractor API to generate pick options for a given sentence and word
func (c *DistractorClient) GenerateDistractors(ctx context.Context, sentence, word string) ([]string, error) {
	if sentence == "" {
		return nil, fmt.Errorf("sentence cannot be empty")
	}
	if word == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}

	// Prepare request payload
	request := DistractorRequest{
		Sentence: sentence,
		Word:     word,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/v1/generate-distractors", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		var apiError DistractorError
		if jsonErr := json.Unmarshal(respBody, &apiError); jsonErr == nil {
			return nil, apiError
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response
	var response DistractorResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.PickOptions, nil
}

// HealthCheck checks if the distractor service is healthy and ready
func (c *DistractorClient) HealthCheck(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	// Parse health response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read health check response: %w", err)
	}

	var health struct {
		Status      string `json:"status"`
		ModelLoaded bool   `json:"model_loaded"`
	}

	if err := json.Unmarshal(respBody, &health); err != nil {
		return false, fmt.Errorf("failed to parse health check response: %w", err)
	}

	return health.Status == "healthy" && health.ModelLoaded, nil
}

// GeneratePickOptionsWithDefaults is a convenience function that generates distractors with default settings
func GeneratePickOptionsWithDefaults(ctx context.Context, sentence, word string) ([]string, error) {
	client := NewDistractorClient(DistractorClientConfig{})
	return client.GenerateDistractors(ctx, sentence, word)
}
