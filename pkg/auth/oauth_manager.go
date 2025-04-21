package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthManager handles authentication with Questrade API
type OAuthManager struct {
	ClientID     string
	RefreshToken string
	AccessToken  string
	ApiServer    string
	ExpiresAt    time.Time
}

// TokenResponse represents the response from Questrade token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ApiServer    string `json:"api_server"`
}

// NewOAuthManager creates a new OAuthManager instance
func NewOAuthManager(clientID, refreshToken string) *OAuthManager {
	return &OAuthManager{
		ClientID:     clientID,
		RefreshToken: refreshToken,
	}
}

// RefreshAccessToken refreshes the access token using the refresh token
func (o *OAuthManager) RefreshAccessToken() error {
	if o.RefreshToken == "" {
		return errors.New("refresh token is required")
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", o.RefreshToken)

	req, err := http.NewRequest("POST", "https://login.questrade.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to refresh token, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	o.AccessToken = tokenResp.AccessToken
	o.RefreshToken = tokenResp.RefreshToken
	o.ApiServer = tokenResp.ApiServer
	o.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// GetAuthenticatedRequest creates an authenticated HTTP request
func (o *OAuthManager) GetAuthenticatedRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	if time.Now().After(o.ExpiresAt) {
		if err := o.RefreshAccessToken(); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
	}

	apiURL := fmt.Sprintf("%s%s", o.ApiServer, endpoint)
	req, err := http.NewRequest(method, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
