package models

// ResourceLinks represents common link objects in API responses
type ResourceLinks struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// PagingInformation represents paging information in API responses
type PagingInformation struct {
	Next  string `json:"next,omitempty"`
	Total int    `json:"total,omitempty"`
}

// Links represents links in API responses
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Errors []APIError `json:"errors"`
}

// APIError represents an individual error in an API response
type APIError struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// PlatformType represents app platform types
type PlatformType string

const (
	PlatformIOS    PlatformType = "IOS"
	PlatformMacOS  PlatformType = "MAC_OS"
	PlatformTvOS   PlatformType = "TV_OS"
)

// ReportType represents report types
type ReportType string

const (
	ReportTypeSales              ReportType = "SALES"
	ReportTypeSubscription       ReportType = "SUBSCRIPTION"
	ReportTypeSubscriptionEvent  ReportType = "SUBSCRIPTION_EVENT"
)

// ReportSubType represents report subtypes
type ReportSubType string

const (
	ReportSubTypeSummary    ReportSubType = "SUMMARY"
	ReportSubTypeDetailed   ReportSubType = "DETAILED"
)

// ReportFrequency represents report frequency
type ReportFrequency string

const (
	ReportFrequencyDaily   ReportFrequency = "DAILY"
	ReportFrequencyWeekly  ReportFrequency = "WEEKLY"
	ReportFrequencyMonthly ReportFrequency = "MONTHLY"
	ReportFrequencyYearly  ReportFrequency = "YEARLY"
)

// String returns the string representation of the frequency
func (f ReportFrequency) String() string {
	return string(f)
}
