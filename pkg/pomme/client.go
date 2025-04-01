package pomme

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/marcusziade/pomme/internal/api"
	"github.com/marcusziade/pomme/internal/auth"
	"github.com/marcusziade/pomme/internal/models"
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
	// Only use mock data if explicitly requested for testing
	if os.Getenv("POMME_DEBUG") == "1" || os.Getenv("POMME_MOCK") == "1" {
		fmt.Println("Debug mode: Generating mock sales report data")
		return generateMockSalesReport(string(freq), reportDate, string(reportType), vendorNumber)
	}

	// We don't need to prepare request data for a GET request
	// The parameters are passed in the URL query string
	
	// Map report type to the correct parameters
	var reportSubType string
	switch reportType {
	case models.ReportTypeSales:
		reportSubType = "SUMMARY" // or DETAILED, SUMMARY_TERRITORY, etc.
	case models.ReportTypeSubscription:
		reportSubType = "SUMMARY" 
	case models.ReportTypeSubscriptionEvent:
		reportSubType = "SUMMARY"
	default:
		reportSubType = "SUMMARY"
	}
	
	// Construct the URL with query parameters
	queryParams := fmt.Sprintf("?filter[frequency]=%s&filter[reportType]=%s&filter[reportSubType]=%s&filter[vendorNumber]=%s&filter[reportDate]=%s",
		string(freq), string(reportType), reportSubType, vendorNumber, reportDate)
	
	// For sales reports, we need to use a custom request with specific headers
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiClient.BaseURL+"/salesReports"+queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set required headers for sales reports
	req.Header.Set("Accept", "application/a-gzip")
	
	// Add authorization token
	token, err := c.apiClient.GetAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	authHeader := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", authHeader)
	
	// Make the API request - sales reports need special handling
	resp, err := c.apiClient.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for error responses
	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("API request failed with status: %s (failed to read response body: %v)", resp.Status, err)
		}
		
		fmt.Printf("Debug: Error response body: %s\n", string(bodyBytes))
		
		// Parse as JSON error to get details
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
				// Return an empty report instead of an error
				fmt.Println("No sales data found for the specified period.")
				// Return empty CSV with headers
				return []byte("Provider,Provider Country,SKU,Developer,Title,Version,Product Type Identifier,Units,Developer Proceeds,Begin Date,End Date,Customer Currency,Country Code,Currency of Proceeds,Apple Identifier,Customer Price,Promo Code,Parent Identifier,Subscription,Period,Category,CMB,Device,Supported Platforms,Proceeds Reason,Preserved Pricing,Client,Order Type\n"), nil
			}
			return nil, fmt.Errorf("API error: %s - %s", errorResp.Errors[0].Code, errorResp.Errors[0].Detail)
		}
		
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}
	
	// For sales reports, the response is gzipped raw report data
	// We need to handle it differently than JSON API responses
	
	fmt.Printf("Debug: Response content type: %s\n", resp.Header.Get("Content-Type"))
	
	// Check if we got the expected content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/a-gzip" {
		// We received gzipped data, so decompress it
		var reader io.ReadCloser
		var err error
		
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.Close()
		
		// Read all the data
		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read gzipped data: %w", err)
		}
		
		// Return the decompressed data
		return data, nil
	} else if contentType == "application/json" {
		// This is likely an error response in JSON format
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
		
		return nil, fmt.Errorf("unexpected response format: %s", contentType)
	} else {
		// Unexpected content type, try to read the raw data anyway
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		
		// Log the first part of the response for debugging
		if len(data) > 0 {
			previewSize := 100
			if len(data) < previewSize {
				previewSize = len(data)
			}
			fmt.Printf("Debug: Response data preview: %s\n", data[:previewSize])
		}
		
		return data, nil
	}
}

// generateMockSalesReport creates mock sales report data for testing
func generateMockSalesReport(freq, reportDate, reportType, vendorNumber string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	// Write header row
	header := []string{
		"Provider", "Provider Country", "SKU", "Developer", "Title", "Version",
		"Product Type Identifier", "Units", "Developer Proceeds", "Begin Date",
		"End Date", "Customer Currency", "Country Code", "Currency of Proceeds",
		"Apple Identifier", "Customer Price", "Promo Code", "Parent Identifier",
		"Subscription", "Period", "Category", "CMB", "Device", "Supported Platforms",
		"Proceeds Reason", "Preserved Pricing", "Client", "Order Type",
	}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}
	
	// Create some mock apps
	apps := []struct {
		title     string
		appleID   string
		sku       string
		price     float64
		countries []string
	}{
		{"Weather Pro", "12345678", "com.example.weather", 4.99, []string{"US", "CA", "GB", "DE", "FR"}},
		{"Fitness Tracker", "23456789", "com.example.fitness", 9.99, []string{"US", "CA", "GB", "AU"}},
		{"Meditation App", "34567890", "com.example.meditation", 2.99, []string{"US", "JP", "KR", "SG"}},
		{"Mobile Scanner", "45678901", "com.example.scanner", 3.99, []string{"US", "GB", "DE", "IT", "ES"}},
	}
	
	// Parse the report date
	var beginDate, endDate time.Time
	var err error
	
	switch freq {
	case "DAILY":
		beginDate, err = time.Parse("2006-01-02", reportDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		endDate = beginDate
	case "WEEKLY":
		beginDate, err = time.Parse("2006-01-02", reportDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		endDate = beginDate.AddDate(0, 0, 6) // Add 6 days for a week
	case "MONTHLY":
		beginDate, err = time.Parse("2006-01", reportDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		// Last day of the month
		endDate = beginDate.AddDate(0, 1, -1)
	case "YEARLY":
		beginDate, err = time.Parse("2006", reportDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		endDate = beginDate.AddDate(1, 0, -1)
	}
	
	// Device types
	devices := []string{"iPhone", "iPad", "Mac", "Apple TV", "Apple Watch"}
	
	// Product types based on report type
	var productTypes []string
	switch reportType {
	case "SALES":
		productTypes = []string{"App", "In-App Purchase", "App Bundle"}
	case "SUBSCRIPTION":
		productTypes = []string{"Auto-Renewable Subscription", "Non-Renewing Subscription"}
	case "SUBSCRIPTION_EVENT":
		productTypes = []string{"Renewal", "Cancellation", "Free Trial"}
	}
	
	// Generate random records
	rand.Seed(time.Now().UnixNano())
	numRecords := 20 + rand.Intn(30) // 20-50 records
	
	for i := 0; i < numRecords; i++ {
		// Select a random app
		app := apps[rand.Intn(len(apps))]
		
		// Select a random country
		country := app.countries[rand.Intn(len(app.countries))]
		
		// Select a random device
		device := devices[rand.Intn(len(devices))]
		
		// Select a random product type
		productType := productTypes[rand.Intn(len(productTypes))]
		
		// Random units and proceeds
		units := 1 + rand.Intn(100)
		customerPrice := app.price
		proceeds := float64(units) * customerPrice * 0.7 // 70% of revenue
		
		// Format dates for Apple's reporting format (MM/DD/YYYY)
		beginDateStr := beginDate.Format("01/02/2006")
		endDateStr := endDate.Format("01/02/2006")
		
		// Generate a record
		record := []string{
			"Apple", "United States", app.sku, "Example Developer", app.title, "1.0",
			productType, fmt.Sprintf("%d", units), fmt.Sprintf("%.2f", proceeds), beginDateStr,
			endDateStr, "USD", country, "USD", app.appleID, fmt.Sprintf("%.2f", customerPrice),
			"", "", "N", "", "Games", "", device, "iOS", "Standard Revenue", "N", "App Store", "New",
		}
		
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}
	
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush CSV writer: %w", err)
	}
	
	return buf.Bytes(), nil
}
