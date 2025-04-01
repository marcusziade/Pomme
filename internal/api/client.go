package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/marcus/pomme/internal/auth"
)

// Client represents an App Store Connect API client
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	AuthConfig  auth.JWTConfig
	jwtToken    string
	tokenExpiry time.Time
}

// NewClient creates a new App Store Connect API client
func NewClient(baseURL string, authConfig auth.JWTConfig) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		AuthConfig: authConfig,
	}
}

// getAuthToken returns a valid JWT token, generating a new one if needed
func (c *Client) getAuthToken() (string, error) {
	// Check if we already have a valid token
	if c.jwtToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.jwtToken, nil
	}

	// Generate a new token
	token, err := auth.GenerateToken(c.AuthConfig)
	if err != nil {
		return "", fmt.Errorf("failed to generate auth token: %w", err)
	}

	// Update the client's token and expiry time
	c.jwtToken = token
	c.tokenExpiry = time.Now().Add(15 * time.Minute) // Token expires in 20 min, but we refresh after 15

	return token, nil
}

// Request makes an HTTP request to the App Store Connect API
func (c *Client) Request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Construct the full URL
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Add authorization token
	token, err := c.getAuthToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for error responses
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		
		var errResp struct {
			Errors []struct {
				Status string `json:"status"`
				Code   string `json:"code"`
				Title  string `json:"title"`
				Detail string `json:"detail"`
			} `json:"errors"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && len(errResp.Errors) > 0 {
			return nil, fmt.Errorf("API error: %s - %s", errResp.Errors[0].Code, errResp.Errors[0].Detail)
		}
		
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return resp, nil
}

// Get performs a GET request to the API
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil)
}

// Post performs a POST request to the API
func (c *Client) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, path, body)
}

// Put performs a PUT request to the API
func (c *Client) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, path, body)
}

// Patch performs a PATCH request to the API
func (c *Client) Patch(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.Request(ctx, http.MethodPatch, path, body)
}

// Delete performs a DELETE request to the API
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil)
}
