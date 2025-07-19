package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// TTSService handles text-to-speech functionality
type TTSService struct {
	cacheDir string
	logger   *zap.Logger
}

// NewTTSService creates a new TTS service
func NewTTSService(cacheDir string, logger *zap.Logger) *TTSService {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		logger.Error("Failed to create TTS cache directory", zap.Error(err))
	}

	return &TTSService{
		cacheDir: cacheDir,
		logger:   logger,
	}
}

// GenerateVoiceMessage generates a voice message for the given text
func (tts *TTSService) GenerateVoiceMessage(text string, language string) ([]byte, error) {
	// Clean and validate input
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty text provided")
	}

	if language == "" {
		language = "en" // Default to English
	}

	// Check cache first
	cacheKey := tts.getCacheKey(text, language)
	cachedPath := filepath.Join(tts.cacheDir, cacheKey+".mp3")

	if _, err := os.Stat(cachedPath); err == nil {
		// File exists in cache, read and return it
		return os.ReadFile(cachedPath)
	}

	// Generate new voice message
	audioData, err := tts.generateFromGoogleTTS(text, language)
	if err != nil {
		tts.logger.Error("Failed to generate TTS audio",
			zap.String("text", text),
			zap.String("language", language),
			zap.Error(err))
		return nil, err
	}

	// Cache the result
	if err := os.WriteFile(cachedPath, audioData, 0644); err != nil {
		tts.logger.Warn("Failed to cache TTS audio", zap.Error(err))
	}

	return audioData, nil
}

// generateFromGoogleTTS generates audio using Google Translate TTS (unofficial API)
func (tts *TTSService) generateFromGoogleTTS(text string, language string) ([]byte, error) {
	// This uses the unofficial Google Translate TTS API
	// In production, you should use the official Google Cloud TTS API

	baseURL := "https://translate.google.com/translate_tts"
	params := url.Values{}
	params.Set("ie", "UTF-8")
	params.Set("q", text)
	params.Set("tl", language)
	params.Set("client", "tw-ob")
	params.Set("tk", tts.generateToken(text)) // Simple token generation

	requestURL := baseURL + "?" + params.Encode()

	// Create HTTP request with proper headers
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://translate.google.com/")

	// Make HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make TTS request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TTS request failed with status: %d", resp.StatusCode)
	}

	// Read response body
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read TTS response: %w", err)
	}

	return audioData, nil
}

// generateToken generates a simple token for the request
// This is a simplified version - the actual Google TTS token generation is more complex
func (tts *TTSService) generateToken(text string) string {
	// Simple hash-based token for demo purposes
	// In production, implement proper Google TTS token generation
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)[:16]
}

// getCacheKey generates a cache key for the given text and language
func (tts *TTSService) getCacheKey(text string, language string) string {
	key := fmt.Sprintf("%s_%s", language, text)
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

// GenerateWordVoiceMessage generates a voice message specifically for word learning
func (tts *TTSService) GenerateWordVoiceMessage(word string) ([]byte, error) {
	// Format the word for better pronunciation
	formattedText := fmt.Sprintf("%s", word)

	return tts.GenerateVoiceMessage(formattedText, "en")
}

// CreateVoiceMessageFromBytes creates a voice message file that can be sent via Telegram
func (tts *TTSService) CreateVoiceMessageFromBytes(audioData []byte, filename string) (string, error) {
	tempPath := filepath.Join(tts.cacheDir, "temp_"+filename+".mp3")

	err := os.WriteFile(tempPath, audioData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary voice file: %w", err)
	}

	return tempPath, nil
}

// CleanupTempFiles removes temporary voice files older than specified duration
func (tts *TTSService) CleanupTempFiles(maxAge time.Duration) error {
	entries, err := os.ReadDir(tts.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), "temp_") {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if now.Sub(info.ModTime()) > maxAge {
				filePath := filepath.Join(tts.cacheDir, entry.Name())
				if err := os.Remove(filePath); err != nil {
					tts.logger.Warn("Failed to remove old temp file",
						zap.String("file", filePath),
						zap.Error(err))
				}
			}
		}
	}

	return nil
}

// ConvertToOGG converts MP3 audio data to OGG format for better Telegram compatibility
// This is a placeholder function - in production, you'd use a proper audio conversion library
func (tts *TTSService) ConvertToOGG(mp3Data []byte) ([]byte, error) {
	// For now, return the MP3 data as-is
	// Telegram supports MP3 for voice messages, though OGG is preferred
	// To implement proper conversion, you could use:
	// - ffmpeg binary
	// - go-audio libraries
	// - External service

	return mp3Data, nil
}

// ValidateAudioData performs basic validation on audio data
func (tts *TTSService) ValidateAudioData(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty audio data")
	}

	if len(data) < 100 {
		return fmt.Errorf("audio data too small, might be corrupted")
	}

	// Check for MP3 header
	if len(data) >= 3 && bytes.Equal(data[:3], []byte{0xFF, 0xFB, 0x90}) {
		return nil // Valid MP3
	}

	// Check for other audio formats if needed
	// For now, assume it's valid if it's not empty and has reasonable size

	return nil
}
