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

// ThesaurusRecommendation represents a single recommendation entry returned by the Thesaurus API
type ThesaurusRecommendation struct {
	Word        string  `json:"word"`
	Topic       string  `json:"topic"`
	Subtopic    string  `json:"subtopic"`
	Subsubtopic string  `json:"subsubtopic"`
	CEFRLevel   string  `json:"CEFR_level"`
	Score       float64 `json:"score"`
}

// thesaurusRecommendRequest is the payload for the /api/recommend endpoint
// It is defined unexported because callers should use the public client method.
type thesaurusRecommendRequest struct {
	Words []string `json:"words"`
}

// ThesaurusError represents an error response from the Thesaurus API
// Currently the service returns standard FastAPI errors so we capture the common fields.
// This struct can be extended if the API evolves.
// See https://fastapi.tiangolo.com/tutorial/handling-errors/ for the default schema.
// Example JSON:
// {"detail": "..."}
type ThesaurusError struct {
	Detail string `json:"detail"`
}

func (e ThesaurusError) Error() string {
	return fmt.Sprintf("thesaurus API error: %s", e.Detail)
}

// ThesaurusClient provides an interface to the internal Thesaurus API.
// It mirrors the style of other utility clients (DictionaryClient, DistractorClient).
// The baseURL should point to the root of the service (e.g. http://localhost:8002).
// All requests will be made relative to this base.
type ThesaurusClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// ThesaurusClientConfig holds configuration for the ThesaurusClient.
// If BaseURL is empty, it will be populated from the THESAURUS_API_URL env var,
// falling back to "http://localhost:8002".
// If Timeout is zero, a sensible default (10s) will be used.
type ThesaurusClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

// NewThesaurusClient creates a new Thesaurus client with the provided configuration.
// It applies sane defaults for missing config values.
func NewThesaurusClient(config ThesaurusClientConfig) *ThesaurusClient {
	if config.BaseURL == "" {
		// Prefer environment variable if provided
		if envURL := os.Getenv("THESAURUS_API_URL"); envURL != "" {
			config.BaseURL = envURL
		} else {
			config.BaseURL = "http://localhost:8002"
		}
	}

	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	return &ThesaurusClient{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		timeout: config.Timeout,
	}
}

// Recommend fetches vocabulary recommendations for the provided list of known words.
// It returns a slice of ThesaurusRecommendation entries.
func (c *ThesaurusClient) Recommend(ctx context.Context, knownWords []string) ([]ThesaurusRecommendation, error) {
	if len(knownWords) == 0 {
		return nil, fmt.Errorf("knownWords cannot be empty")
	}

	// Prepare request payload
	payload := thesaurusRecommendRequest{Words: knownWords}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/recommend", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		var apiErr ThesaurusError
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.Detail != "" {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response
	var recommendations []ThesaurusRecommendation
	if err := json.Unmarshal(respBody, &recommendations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return recommendations, nil
}

// HealthCheck performs a basic health check by pinging the /health endpoint.
// It returns true if the service responds with HTTP 200.
func (c *ThesaurusClient) HealthCheck(ctx context.Context) (bool, error) {
	// Prepare request body {"ping":"test"}
	body, _ := json.Marshal(map[string]string{"ping": "test"})

	url := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return true, nil
}

// RecommendWithDefaults is a convenience wrapper that creates a client with default settings
// and immediately calls Recommend.
func RecommendWithDefaults(ctx context.Context, knownWords []string) ([]ThesaurusRecommendation, error) {
	client := NewThesaurusClient(ThesaurusClientConfig{})
	return client.Recommend(ctx, knownWords)
}
