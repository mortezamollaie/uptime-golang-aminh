package main

import (
	"fmt"
	"net/url"
	"strings"
)

// IsValidURL validates if a string is a valid HTTP/HTTPS URL
func IsValidURL(urlStr string) bool {
	if strings.TrimSpace(urlStr) == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

// SanitizeURL cleans and validates a URL string
func SanitizeURL(urlStr string) (string, error) {
	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Add http:// if no scheme is provided
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + urlStr
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("URL must use HTTP or HTTPS protocol")
	}

	return parsedURL.String(), nil
}
