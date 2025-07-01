package models

import "time"

// CustomerReview represents a customer review from the App Store
type CustomerReview struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes CustomerReviewAttributes `json:"attributes"`
	Links      map[string]interface{}   `json:"links"`
}

// CustomerReviewAttributes contains review details
type CustomerReviewAttributes struct {
	Rating               int       `json:"rating"`
	Title                string    `json:"title"`
	Body                 string    `json:"body"`
	ReviewerNickname     string    `json:"reviewerNickname"`
	CreatedDate          time.Time `json:"createdDate"`
	Territory            string    `json:"territory"`
}

// CustomerReviewResponse represents a developer response to a review
type CustomerReviewResponse struct {
	ID         string                           `json:"id"`
	Type       string                           `json:"type"`
	Attributes CustomerReviewResponseAttributes `json:"attributes"`
}

// CustomerReviewResponseAttributes contains response details
type CustomerReviewResponseAttributes struct {
	ResponseBody string    `json:"responseBody"`
	ModifiedDate time.Time `json:"modifiedDate"`
	State        string    `json:"state"` // PENDING_PUBLISH, PUBLISHED
}

// ReviewSummary provides aggregated review statistics
type ReviewSummary struct {
	AppID          string             `json:"appId"`
	TotalReviews   int                `json:"totalReviews"`
	AverageRating  float64            `json:"averageRating"`
	RatingCounts   map[int]int        `json:"ratingCounts"`   // 1-5 star counts
	TerritoryStats []TerritoryReviews `json:"territoryStats"`
	RecentReviews  []CustomerReview   `json:"recentReviews"`
}

// TerritoryReviews contains review stats for a specific territory
type TerritoryReviews struct {
	Territory     string  `json:"territory"`
	ReviewCount   int     `json:"reviewCount"`
	AverageRating float64 `json:"averageRating"`
}

// ReviewFilter contains filtering options for reviews
type ReviewFilter struct {
	AppID       string    `json:"appId"`
	Territory   string    `json:"territory,omitempty"`
	Rating      int       `json:"rating,omitempty"`      // 1-5
	StartDate   time.Time `json:"startDate,omitempty"`
	EndDate     time.Time `json:"endDate,omitempty"`
	Limit       int       `json:"limit,omitempty"`
	Sort        string    `json:"sort,omitempty"`        // mostRecent, mostCritical, mostHelpful
}