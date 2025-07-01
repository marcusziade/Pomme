package pomme

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/marcusziade/pomme/internal/api"
	"github.com/marcusziade/pomme/internal/auth"
	"github.com/marcusziade/pomme/internal/models"
)

// Client is a high-level client for the App Store Connect API
type Client struct {
	apiClient *api.Client
	keyID     string
	issuerID  string
}

// NewClient creates a new App Store Connect API client
func NewClient(keyID, issuerID, privateKey string) *Client {
	authConfig := auth.JWTConfig{
		KeyID:         keyID,
		IssuerID:      issuerID,
		PrivateKeyPEM: privateKey,
		Expiration:    20 * time.Minute,
	}
	
	apiClient := api.NewClient("https://api.appstoreconnect.apple.com", authConfig)
	
	return &Client{
		apiClient: apiClient,
		keyID:     keyID,
		issuerID:  issuerID,
	}
}

// GetSalesReport fetches a sales report from the App Store Connect API
func (c *Client) GetSalesReport(ctx context.Context, frequency models.ReportFrequency, reportDate string, reportType models.ReportType, vendorNumber string) ([]byte, error) {
	// Construct the report API URL
	url := fmt.Sprintf("/v1/salesReports?filter[frequency]=%s&filter[reportDate]=%s&filter[reportSubType]=%s&filter[reportType]=%s&filter[vendorNumber]=%s",
		frequency, reportDate, "SUMMARY", reportType, vendorNumber)
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiClient.BaseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Get auth token
	token, err := c.apiClient.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/a-gzip")
	
	// Execute the request
	resp, err := c.apiClient.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		// Try to parse as API error
		var errorResp struct {
			Errors []struct {
				Code   string `json:"code"`
				Status string `json:"status"`
				Title  string `json:"title"`
				Detail string `json:"detail"`
			} `json:"errors"`
		}
		
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil && len(errorResp.Errors) > 0 {
			// Special case for "no sales data"
			if resp.StatusCode == 404 && errorResp.Errors[0].Code == "NOT_FOUND" &&
				strings.Contains(errorResp.Errors[0].Detail, "no sales") {
				// Return nil to indicate no data available
				return nil, nil
			}
			return nil, fmt.Errorf("API error: %s - %s", errorResp.Errors[0].Code, errorResp.Errors[0].Detail)
		}
		
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}
	
	// Handle the response based on content type
	contentType := resp.Header.Get("Content-Type")
	
	switch contentType {
	case "application/a-gzip":
		// Decompress gzipped data
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.Close()
		
		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read gzipped data: %w", err)
		}
		
		return data, nil
		
	case "application/json":
		// This might be an error response
		var errorResp struct {
			Errors []struct {
				Code   string `json:"code"`
				Status string `json:"status"`
				Title  string `json:"title"`
				Detail string `json:"detail"`
			} `json:"errors"`
		}
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil && len(errorResp.Errors) > 0 {
			return nil, fmt.Errorf("API error: %s - %s", errorResp.Errors[0].Code, errorResp.Errors[0].Detail)
		}
		
		return nil, fmt.Errorf("unexpected JSON response for sales report")
		
	default:
		// Try to read raw data
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		
		// If it's text/plain, it might be uncompressed CSV
		if strings.HasPrefix(contentType, "text/") || len(data) > 0 {
			return data, nil
		}
		
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}
}

// GetFinancialReport fetches a financial report from the App Store Connect API
func (c *Client) GetFinancialReport(ctx context.Context, regionCode, fiscalYear, fiscalPeriod, vendorNumber string) ([]byte, error) {
	// Construct the financial report API URL
	url := fmt.Sprintf("/v1/financeReports?filter[regionCode]=%s&filter[reportDate]=%s-%s&filter[reportType]=FINANCIAL&filter[vendorNumber]=%s",
		regionCode, fiscalYear, fiscalPeriod, vendorNumber)
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiClient.BaseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Get auth token
	token, err := c.apiClient.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/a-gzip")
	
	// Execute the request
	resp, err := c.apiClient.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	// Handle gzipped response
	if resp.Header.Get("Content-Type") == "application/a-gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.Close()
		
		return io.ReadAll(reader)
	}
	
	// Otherwise, read raw response
	return io.ReadAll(resp.Body)
}

// ListApps retrieves a list of apps
func (c *Client) ListApps(ctx context.Context) ([]models.App, error) {
	url := "/v1/apps"
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiClient.BaseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Get auth token
	token, err := c.apiClient.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.apiClient.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var result struct {
		Data []models.App `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Data, nil
}

// GetApps retrieves a list of apps (stub)
func (c *Client) GetApps(ctx context.Context) ([]models.App, error) {
	return c.ListApps(ctx)
}

// GetApp retrieves a single app (stub)
func (c *Client) GetApp(ctx context.Context, appID string) (*models.App, error) {
	// TODO: Implement single app fetch
	return nil, fmt.Errorf("not implemented")
}

// GetReviews retrieves reviews for an app (stub)
func (c *Client) GetReviews(ctx context.Context, appID string) ([]models.Review, error) {
	// TODO: Implement reviews fetch
	return nil, fmt.Errorf("not implemented")
}