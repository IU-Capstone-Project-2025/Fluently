package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTimeFormat parses time string in various formats and returns HH:MM format
// Supports formats like: "9 30", "9:30", "09:30", "9.30", "930"
func ParseTimeFormat(timeStr string) (string, error) {
	// Remove all spaces and normalize
	normalized := strings.ReplaceAll(timeStr, " ", "")

	// Try different formats
	formats := []string{
		"15:04", // Standard HH:MM
		"15.04", // HH.MM
		"1504",  // HHMM
		"3:04",  // H:MM
		"3.04",  // H.MM
		"304",   // HMM
		"3:4",   // H:M
		"3.4",   // H.M
		"34",    // HM
	}

	for _, format := range formats {
		if parsedTime, err := time.Parse(format, normalized); err == nil {
			// Validate hours and minutes
			hour := parsedTime.Hour()
			minute := parsedTime.Minute()

			if hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59 {
				return parsedTime.Format("15:04"), nil
			}
		}
	}

	// Special handling for 4-digit format like "0930", "2130"
	if len(normalized) == 4 {
		hourStr := normalized[:2]
		minuteStr := normalized[2:]

		hour, err1 := strconv.Atoi(hourStr)
		minute, err2 := strconv.Atoi(minuteStr)

		if err1 == nil && err2 == nil && hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59 {
			return fmt.Sprintf("%02d:%02d", hour, minute), nil
		}
	}

	// Special handling for 3-digit format like "930", "2130" (when hour is single digit)
	if len(normalized) == 3 {
		hourStr := normalized[:1]
		minuteStr := normalized[1:]

		hour, err1 := strconv.Atoi(hourStr)
		minute, err2 := strconv.Atoi(minuteStr)

		if err1 == nil && err2 == nil && hour >= 0 && hour <= 9 && minute >= 0 && minute <= 59 {
			return fmt.Sprintf("%02d:%02d", hour, minute), nil
		}
	}

	// If no format worked, try to parse manually
	// Handle cases like "9 30", "9:30", "09:30"
	parts := strings.FieldsFunc(timeStr, func(r rune) bool {
		return r == ':' || r == '.' || r == ' '
	})

	if len(parts) == 2 {
		hourStr := strings.TrimSpace(parts[0])
		minuteStr := strings.TrimSpace(parts[1])

		hour, err1 := strconv.Atoi(hourStr)
		minute, err2 := strconv.Atoi(minuteStr)

		if err1 == nil && err2 == nil && hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59 {
			return fmt.Sprintf("%02d:%02d", hour, minute), nil
		}
	}

	return "", fmt.Errorf("invalid time format: %s", timeStr)
}

// IsValidTimeFormat checks if the time string can be parsed into valid HH:MM format
func IsValidTimeFormat(timeStr string) bool {
	_, err := ParseTimeFormat(timeStr)
	return err == nil
}

// ParseTimeToTime parses time string and returns time.Time object
// This is useful for backend API calls that expect time.Time
func ParseTimeToTime(timeStr string) (*time.Time, error) {
	parsedTime, err := ParseTimeFormat(timeStr)
	if err != nil {
		return nil, err
	}

	// Parse the standardized format
	t, err := time.Parse("15:04", parsedTime)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
