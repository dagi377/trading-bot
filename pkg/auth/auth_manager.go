package auth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AuthManager handles authentication for data sources
type AuthManager struct {
	APIKeys     map[string]string
	AccessToken string
	ExpiresAt   time.Time
}

// BaseTokenResponse represents a generic token response
type BaseTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewAuthManager creates a new AuthManager instance
func NewAuthManager() *AuthManager {
	return &AuthManager{
		APIKeys: make(map[string]string),
	}
}

// AddAPIKey adds an API key for a specific service
func (a *AuthManager) AddAPIKey(service, apiKey string) {
	a.APIKeys[service] = apiKey
}

// GetAPIKey retrieves an API key for a specific service
func (a *AuthManager) GetAPIKey(service string) (string, error) {
	if key, exists := a.APIKeys[service]; exists {
		return key, nil
	}
	return "", fmt.Errorf("API key for service %s not found", service)
}

// GetAuthenticatedRequest creates an authenticated HTTP request
func (a *AuthManager) GetAuthenticatedRequest(service, method, endpoint string, body io.Reader) (*http.Request, error) {
	apiKey, err := a.GetAPIKey(service)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication based on service type
	switch service {
	case "alphavantage":
		// Alpha Vantage uses query parameters for authentication
		q := req.URL.Query()
		q.Add("apikey", apiKey)
		req.URL.RawQuery = q.Encode()
	case "finnhub":
		// Finnhub uses an X-Finnhub-Token header
		req.Header.Add("X-Finnhub-Token", apiKey)
	case "financialmodelingprep":
		// FMP uses query parameters for authentication
		q := req.URL.Query()
		q.Add("apikey", apiKey)
		req.URL.RawQuery = q.Encode()
	default:
		// Default to Authorization header
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// ValidateAuth checks if authentication is valid for a service
func (a *AuthManager) ValidateAuth(service string) error {
	_, err := a.GetAPIKey(service)
	if err != nil {
		return errors.New("missing API key for " + service)
	}
	return nil
}
