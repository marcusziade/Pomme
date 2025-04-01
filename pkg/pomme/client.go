package pomme

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marcus/pomme/internal/api"
	"github.com/marcus/pomme/internal/auth"
	"github.com/marcus/pomme/internal/models"
)

// Client is a high-level client for the App Store Connect API
type Client struct {
	apiClient    *api.Client
	vendorNumber string
}

// NewClient creates a new App Store Connect API client
func NewClient(keyID, issuerID, privateKeyPEM string, vendorNumber ...string) *Client {
	// Create auth config
	authConfig := auth.JWTConfig{
		KeyID:         keyID,
		IssuerID:      issuerID,
		PrivateKeyPEM: privateKeyPEM,
	}

	// Create API client
	apiClient := api.NewClient("https://api.appstoreconnect.apple.com/v1", authConfig)

	// Initialize client
	client := &Client{
		apiClient: apiClient,
	}
	
	// Set vendor number if provided
	if len(vendorNumber) > 0 && vendorNumber[0] != "" {
		client.vendorNumber = vendorNumber[0]
	}

	return client
}

// GetApps retrieves a list of apps
func (c *Client) GetApps(ctx context.Context) (*models.AppsResponse, error) {
	// Print debug info
	fmt.Println("Debug: Using Key ID:", c.apiClient.AuthConfig.KeyID)
	fmt.Println("Debug: Using Issuer ID:", c.apiClient.AuthConfig.IssuerID)
	fmt.Println("Debug: Private Key length:", len(c.apiClient.AuthConfig.PrivateKeyPEM))

	// Use specific query parameters to make sure we're using the API correctly
	query := "?limit=50&fields[apps]=name,bundleId,primaryLocale,sku"
	resp, err := c.apiClient.Get(ctx, "/apps"+query)
	if err != nil {
		return nil, fmt.Errorf("failed to get apps: %w", err)
	}
	defer resp.Body.Close()
	
	// Debug response
	statusCode := resp.StatusCode
	fmt.Printf("Debug: API response status code: %d\n", statusCode)
	fmt.Printf("Debug: Response headers: %v\n", resp.Header)

	var appsResp models.AppsResponse
	if err := json.NewDecoder(resp.Body).Decode(&appsResp); err != nil {
		return nil, fmt.Errorf("failed to decode apps response: %w", err)
	}

	return &appsResp, nil
}

// GetApp retrieves information about a specific app
func (c *Client) GetApp(ctx context.Context, appID string) (*models.AppResponse, error) {
	resp, err := c.apiClient.Get(ctx, fmt.Sprintf("/apps/%s", appID))
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}
	defer resp.Body.Close()

	var appResp models.AppResponse
	if err := json.NewDecoder(resp.Body).Decode(&appResp); err != nil {
		return nil, fmt.Errorf("failed to decode app response: %w", err)
	}

	return &appResp, nil
}

// GetReviews retrieves reviews for a specific app
func (c *Client) GetReviews(ctx context.Context, appID string, limit int) (*models.ReviewsResponse, error) {
	path := fmt.Sprintf("/apps/%s/customerReviews?limit=%d", appID, limit)
	resp, err := c.apiClient.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer resp.Body.Close()

	var reviewsResp models.ReviewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&reviewsResp); err != nil {
		return nil, fmt.Errorf("failed to decode reviews response: %w", err)
	}

	return &reviewsResp, nil
}

// GetSalesReport retrieves a sales report
func (c *Client) GetSalesReport(
	ctx context.Context,
	freq models.ReportFrequency,
	reportDate string,
	reportType models.ReportType,
	vendorNumber string,
) ([]byte, error) {
	// This is a placeholder implementation
	// TODO: Implement the actual sales report API call

	// For now, just simulate a successful API call with demo data
	// In a real implementation, we would make an API call to /salesReports
	time.Sleep(1 * time.Second) // Simulate API latency
	
	// Return empty data for now - in a real implementation, this would be the report data
	return []byte{}, fmt.Errorf("not implemented yet")
}
