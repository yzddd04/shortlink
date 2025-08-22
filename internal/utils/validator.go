package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// GenerateShortCode generates a random short code
func GenerateShortCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// ValidateURL validates if the given string is a valid URL
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Add scheme if missing
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	return nil
}

// ValidateShortCode validates if the short code meets requirements
func ValidateShortCode(code string) error {
	if len(code) < 3 || len(code) > 20 {
		return fmt.Errorf("short code must be between 3 and 20 characters")
	}

	// Only allow alphanumeric characters and hyphens
	matched, err := regexp.MatchString(`^[a-zA-Z0-9-]+$`, code)
	if err != nil {
		return fmt.Errorf("error validating short code: %w", err)
	}

	if !matched {
		return fmt.Errorf("short code can only contain letters, numbers, and hyphens")
	}

	return nil
}

// SanitizeURL cleans and normalizes URLs
func SanitizeURL(urlStr string) string {
	// Remove leading/trailing whitespace
	urlStr = strings.TrimSpace(urlStr)

	// Add scheme if missing
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	return urlStr
}
