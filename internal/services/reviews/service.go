package reviews

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/marcusziade/pomme/internal/client"
	"github.com/marcusziade/pomme/internal/models"
	"github.com/marcusziade/pomme/internal/services/cache"
)

// Service provides customer review functionality
type Service struct {
	client *client.Client
	cache  cache.Cache
	mu     sync.RWMutex
}

// NewService creates a new reviews service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
		cache:  cache.NewMemoryCache(),
	}
}

// GetReviews fetches customer reviews based on filter
func (s *Service) GetReviews(ctx context.Context, filter models.ReviewFilter) ([]models.CustomerReview, error) {
	cacheKey := fmt.Sprintf("reviews_%s_%s_%d", filter.AppID, filter.Territory, filter.Rating)
	
	// Check cache
	if cached, err := s.cache.Get(cacheKey); err == nil {
		if reviews, ok := cached.([]models.CustomerReview); ok {
			return reviews, nil
		}
	}

	// Build request
	endpoint := fmt.Sprintf("/v1/apps/%s/customerReviews", filter.AppID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	if filter.Territory != "" {
		q.Add("filter[territory]", filter.Territory)
	}
	if filter.Rating > 0 {
		q.Add("filter[rating]", strconv.Itoa(filter.Rating))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	} else {
		q.Add("limit", "100")
	}
	if filter.Sort != "" {
		q.Add("sort", s.mapSortField(filter.Sort))
	} else {
		q.Add("sort", "-createdDate")
	}
	req.URL.RawQuery = q.Encode()

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Data  []models.CustomerReview `json:"data"`
		Links map[string]interface{}  `json:"links"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Cache the result
	s.cache.Set(cacheKey, response.Data, 5*time.Minute)

	return response.Data, nil
}

// GetReviewSummary fetches aggregated review statistics
func (s *Service) GetReviewSummary(ctx context.Context, appID string) (*models.ReviewSummary, error) {
	// Fetch reviews for summary
	filter := models.ReviewFilter{
		AppID: appID,
		Limit: 200,
	}
	
	reviews, err := s.GetReviews(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate summary statistics
	summary := &models.ReviewSummary{
		AppID:        appID,
		TotalReviews: len(reviews),
		RatingCounts: make(map[int]int),
	}

	// Territory stats map
	territoryMap := make(map[string]*models.TerritoryReviews)
	
	// Process reviews
	var totalRating float64
	for _, review := range reviews {
		// Update rating counts
		summary.RatingCounts[review.Attributes.Rating]++
		totalRating += float64(review.Attributes.Rating)
		
		// Update territory stats
		territory := review.Attributes.Territory
		if _, exists := territoryMap[territory]; !exists {
			territoryMap[territory] = &models.TerritoryReviews{
				Territory: territory,
			}
		}
		territoryMap[territory].ReviewCount++
		territoryMap[territory].AverageRating += float64(review.Attributes.Rating)
	}

	// Calculate averages
	if summary.TotalReviews > 0 {
		summary.AverageRating = totalRating / float64(summary.TotalReviews)
	}

	// Finalize territory stats
	for _, stats := range territoryMap {
		if stats.ReviewCount > 0 {
			stats.AverageRating /= float64(stats.ReviewCount)
		}
		summary.TerritoryStats = append(summary.TerritoryStats, *stats)
	}

	// Add recent reviews (first 10)
	if len(reviews) > 10 {
		summary.RecentReviews = reviews[:10]
	} else {
		summary.RecentReviews = reviews
	}

	return summary, nil
}

// RespondToReview creates or updates a response to a customer review
func (s *Service) RespondToReview(ctx context.Context, reviewID, responseText string) error {
	// Check if response already exists
	getEndpoint := fmt.Sprintf("/v1/customerReviews/%s/response", reviewID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, getEndpoint, nil)
	if err != nil {
		return fmt.Errorf("creating get request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		// Response exists, update it
		resp.Body.Close()
		return s.updateReviewResponse(ctx, reviewID, responseText)
	}
	if resp != nil {
		resp.Body.Close()
	}

	// Create new response
	return s.createReviewResponse(ctx, reviewID, responseText)
}

// createReviewResponse creates a new response to a review
func (s *Service) createReviewResponse(ctx context.Context, reviewID, responseText string) error {
	endpoint := "/v1/customerReviewResponses"
	
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "customerReviewResponses",
			"attributes": map[string]string{
				"responseBody": responseText,
			},
			"relationships": map[string]interface{}{
				"review": map[string]interface{}{
					"data": map[string]string{
						"type": "customerReviews",
						"id":   reviewID,
					},
				},
			},
		},
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, endpoint, payload)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// updateReviewResponse updates an existing response
func (s *Service) updateReviewResponse(ctx context.Context, reviewID, responseText string) error {
	// First get the response ID
	getEndpoint := fmt.Sprintf("/v1/customerReviews/%s/response", reviewID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, getEndpoint, nil)
	if err != nil {
		return fmt.Errorf("creating get request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("getting response: %w", err)
	}
	defer resp.Body.Close()

	var responseData struct {
		Data models.CustomerReviewResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	// Update the response
	updateEndpoint := fmt.Sprintf("/v1/customerReviewResponses/%s", responseData.Data.ID)
	
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "customerReviewResponses",
			"id":   responseData.Data.ID,
			"attributes": map[string]string{
				"responseBody": responseText,
			},
		},
	}

	updateReq, err := s.client.NewRequest(ctx, http.MethodPatch, updateEndpoint, payload)
	if err != nil {
		return fmt.Errorf("creating update request: %w", err)
	}

	updateResp, err := s.client.Do(updateReq)
	if err != nil {
		return fmt.Errorf("executing update request: %w", err)
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", updateResp.StatusCode)
	}

	return nil
}

// mapSortField maps user-friendly sort names to API fields
func (s *Service) mapSortField(sort string) string {
	switch strings.ToLower(sort) {
	case "mostrecent", "recent":
		return "-createdDate"
	case "mostcritical", "critical":
		return "rating"
	case "mosthelpful", "helpful":
		return "-helpfulCount"
	default:
		return "-createdDate"
	}
}

// GetReviewsByRating fetches reviews filtered by rating
func (s *Service) GetReviewsByRating(ctx context.Context, appID string, rating int) ([]models.CustomerReview, error) {
	filter := models.ReviewFilter{
		AppID:  appID,
		Rating: rating,
		Limit:  100,
	}
	return s.GetReviews(ctx, filter)
}

// GetReviewsByTerritory fetches reviews for a specific territory
func (s *Service) GetReviewsByTerritory(ctx context.Context, appID, territory string) ([]models.CustomerReview, error) {
	filter := models.ReviewFilter{
		AppID:     appID,
		Territory: territory,
		Limit:     100,
	}
	return s.GetReviews(ctx, filter)
}

// SearchReviews searches reviews by keyword
func (s *Service) SearchReviews(ctx context.Context, appID, keyword string) ([]models.CustomerReview, error) {
	// Get all reviews
	reviews, err := s.GetReviews(ctx, models.ReviewFilter{
		AppID: appID,
		Limit: 500,
	})
	if err != nil {
		return nil, err
	}

	// Filter by keyword
	keyword = strings.ToLower(keyword)
	var filtered []models.CustomerReview
	
	for _, review := range reviews {
		title := strings.ToLower(review.Attributes.Title)
		body := strings.ToLower(review.Attributes.Body)
		
		if strings.Contains(title, keyword) || strings.Contains(body, keyword) {
			filtered = append(filtered, review)
		}
	}

	return filtered, nil
}