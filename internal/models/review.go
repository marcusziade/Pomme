package models

// ReviewResponse represents the response for a customer review resource
type ReviewResponse struct {
	Data  Review   `json:"data"`
	Links Links     `json:"links"`
}

// ReviewsResponse represents the response for multiple customer review resources
type ReviewsResponse struct {
	Data  []Review `json:"data"`
	Links Links     `json:"links"`
	Meta  struct {
		Paging PagingInformation `json:"paging"`
	} `json:"meta,omitempty"`
}

// Review represents a customer review resource
type Review struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes struct {
		Rating               int    `json:"rating"`
		Title                string `json:"title"`
		Body                 string `json:"body"`
		ReviewerNickname     string `json:"reviewerNickname"`
		CreatedDate          string `json:"createdDate"`
		Territory            string `json:"territory"`
		PublishedResponse    *DeveloperResponse `json:"publishedResponse,omitempty"`
	} `json:"attributes"`
	Relationships struct {
		App struct {
			Links ResourceLinks `json:"links"`
		} `json:"app"`
	} `json:"relationships"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// DeveloperResponse represents a developer response to a customer review
type DeveloperResponse struct {
	ID             string `json:"id,omitempty"`
	ResponseBody   string `json:"responseBody"`
	LastModifiedDate string `json:"lastModifiedDate"`
	State          string `json:"state"`
}
