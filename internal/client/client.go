package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/marcusziade/pomme/internal/api"
	"github.com/marcusziade/pomme/internal/auth"
	"github.com/marcusziade/pomme/internal/config"
)

// Client wraps the API client with convenience methods
type Client struct {
	apiClient *api.Client
}

// New creates a new client from config
func New(cfg *config.Config) (*Client, error) {
	// Read private key from file
	privateKeyData, err := os.ReadFile(cfg.Auth.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	
	authConfig := auth.JWTConfig{
		KeyID:          cfg.Auth.KeyID,
		IssuerID:       cfg.Auth.IssuerID,
		PrivateKeyPEM:  string(privateKeyData),
		Expiration:     20 * time.Minute,
	}
	
	apiClient := api.NewClient(cfg.API.BaseURL, authConfig)
	
	return &Client{
		apiClient: apiClient,
	}, nil
}

// NewRequest creates a new HTTP request with JSON body
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var bodyReader *bytes.Reader
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}
	
	// Create the full URL
	url := fmt.Sprintf("%s%s", c.apiClient.BaseURL, path)
	
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	// Add auth token
	token, err := c.apiClient.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	return req, nil
}

// Do executes an HTTP request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.apiClient.HTTPClient.Do(req)
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.apiClient.Get(ctx, path)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	return c.apiClient.Post(ctx, path, bytes.NewReader(jsonBody))
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	return c.apiClient.Patch(ctx, path, bytes.NewReader(jsonBody))
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.apiClient.Delete(ctx, path)
}