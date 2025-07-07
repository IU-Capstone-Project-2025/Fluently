package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DictionaryPhonetic represents a phonetic entry in the API response
type DictionaryPhonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

// DictionaryDefinition represents a definition entry
type DictionaryDefinition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example,omitempty"`
	Synonyms   []string `json:"synonyms,omitempty"`
	Antonyms   []string `json:"antonyms,omitempty"`
}

// DictionaryMeaning represents a meaning entry with part of speech
type DictionaryMeaning struct {
	PartOfSpeech string                 `json:"partOfSpeech"`
	Definitions  []DictionaryDefinition `json:"definitions"`
}

// DictionaryResponse represents the response from dictionaryapi.dev
type DictionaryResponse struct {
	Word      string               `json:"word"`
	Phonetic  string               `json:"phonetic,omitempty"`
	Phonetics []DictionaryPhonetic `json:"phonetics,omitempty"`
	Origin    string               `json:"origin,omitempty"`
	Meanings  []DictionaryMeaning  `json:"meanings"`
}

// DictionaryError represents an error response from the dictionary API
type DictionaryError struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	Resolution string `json:"resolution"`
}

func (e DictionaryError) Error() string {
	return fmt.Sprintf("dictionary API error: %s - %s", e.Title, e.Message)
}

// WordInfo represents processed information from the dictionary API
type WordInfo struct {
	Word         string
	Phonetic     string
	AudioURL     string
	PartOfSpeech string
	Definition   string
	Example      string
	Origin       string
}

// DictionaryClient provides an interface to the dictionaryapi.dev service
type DictionaryClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// DictionaryClientConfig holds configuration for the DictionaryClient
type DictionaryClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

// NewDictionaryClient creates a new dictionary client with the given configuration
func NewDictionaryClient(config DictionaryClientConfig) *DictionaryClient {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.dictionaryapi.dev"
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	return &DictionaryClient{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		timeout: config.Timeout,
	}
}

// GetWordInfo fetches word information from the dictionary API
func (c *DictionaryClient) GetWordInfo(ctx context.Context, word string) (*WordInfo, error) {
	if word == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/v2/entries/en/%s", c.baseURL, word)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

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
		var apiError DictionaryError
		if jsonErr := json.Unmarshal(respBody, &apiError); jsonErr == nil {
			return nil, apiError
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response - API returns an array
	var responses []DictionaryResponse
	if err := json.Unmarshal(respBody, &responses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("no entries found for word: %s", word)
	}

	// Process the first entry
	return c.processResponse(&responses[0])
}

// processResponse extracts relevant information from the API response
func (c *DictionaryClient) processResponse(response *DictionaryResponse) (*WordInfo, error) {
	wordInfo := &WordInfo{
		Word:   response.Word,
		Origin: response.Origin,
	}

	// Extract phonetic information
	if response.Phonetic != "" {
		wordInfo.Phonetic = response.Phonetic
	} else if len(response.Phonetics) > 0 {
		// Use the first phonetic entry if main phonetic is empty
		for _, phonetic := range response.Phonetics {
			if phonetic.Text != "" {
				wordInfo.Phonetic = phonetic.Text
				break
			}
		}
	}

	// Extract audio URL
	for _, phonetic := range response.Phonetics {
		if phonetic.Audio != "" {
			// Ensure the audio URL is complete
			audioURL := phonetic.Audio
			if len(audioURL) >= 2 && audioURL[:2] == "//" {
				audioURL = "https:" + audioURL
			}
			wordInfo.AudioURL = audioURL
			break
		}
	}

	// Extract part of speech and definition from the first meaning
	if len(response.Meanings) > 0 {
		meaning := response.Meanings[0]
		wordInfo.PartOfSpeech = meaning.PartOfSpeech

		if len(meaning.Definitions) > 0 {
			definition := meaning.Definitions[0]
			wordInfo.Definition = definition.Definition
			wordInfo.Example = definition.Example
		}
	}

	return wordInfo, nil
}

// GetAllPartsOfSpeech returns all parts of speech for a word
func (c *DictionaryClient) GetAllPartsOfSpeech(ctx context.Context, word string) ([]string, error) {
	if word == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}

	url := fmt.Sprintf("%s/api/v2/entries/en/%s", c.baseURL, word)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiError DictionaryError
		if jsonErr := json.Unmarshal(respBody, &apiError); jsonErr == nil {
			return nil, apiError
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var responses []DictionaryResponse
	if err := json.Unmarshal(respBody, &responses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("no entries found for word: %s", word)
	}

	// Collect all unique parts of speech
	partsOfSpeech := make(map[string]bool)
	for _, response := range responses {
		for _, meaning := range response.Meanings {
			if meaning.PartOfSpeech != "" {
				partsOfSpeech[meaning.PartOfSpeech] = true
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(partsOfSpeech))
	for pos := range partsOfSpeech {
		result = append(result, pos)
	}

	return result, nil
}

// HealthCheck checks if the dictionary service is accessible
func (c *DictionaryClient) HealthCheck(ctx context.Context) (bool, error) {
	// Test with a simple word to check if the service is working
	_, err := c.GetWordInfo(ctx, "test")
	if err != nil {
		// Check if it's a network error or service unavailable
		var dictionaryErr DictionaryError
		if ok := json.Unmarshal([]byte(err.Error()), &dictionaryErr); ok == nil {
			// If we get a proper API error response, the service is working
			// even if the word might not be found
			return true, dictionaryErr
		}
		return false, fmt.Errorf("dictionary service health check failed: %w", err)
	}
	return true, nil
}

// GetWordInfoWithDefaults is a convenience function that gets word info with default settings
func GetWordInfoWithDefaults(ctx context.Context, word string) (*WordInfo, error) {
	client := NewDictionaryClient(DictionaryClientConfig{})
	return client.GetWordInfo(ctx, word)
}
