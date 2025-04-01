package models

// SalesReportResponse represents the response for a sales report request
type SalesReportResponse struct {
	ReportData []byte  `json:"data"`
	Links     Links     `json:"links"`
}

// SalesReportRequest represents a request for a sales report
type SalesReportRequest struct {
	FilterFinancialReport SalesReportFilter `json:"filter[financialReport]"`
	FilterFrequency       string             `json:"filter[frequency]"`
	FilterReportType      string             `json:"filter[reportType]"`
	FilterReportSubType   string             `json:"filter[reportSubType],omitempty"`
	FilterVendorNumber    string             `json:"filter[vendorNumber]"`
	FilterVersion         string             `json:"filter[version],omitempty"`
}

// SalesReportFilter represents filter parameters for a sales report
type SalesReportFilter struct {
	ReportDate string `json:"reportDate"`
	ReportType string `json:"reportType"`
}

// SalesRecord represents a record in a sales report
type SalesRecord struct {
	Provider            string `csv:"Provider"`
	ProviderCountry     string `csv:"Provider Country"`
	SKU                 string `csv:"SKU"`
	Developer           string `csv:"Developer"`
	Title               string `csv:"Title"`
	Version             string `csv:"Version"`
	ProductTypeID       string `csv:"Product Type Identifier"`
	Units               int    `csv:"Units"`
	DeveloperProceeds   float64 `csv:"Developer Proceeds"`
	BeginDate           string `csv:"Begin Date"`
	EndDate             string `csv:"End Date"`
	CustomerCurrency    string `csv:"Customer Currency"`
	CountryCode         string `csv:"Country Code"`
	CurrencyOfProceeds  string `csv:"Currency of Proceeds"`
	AppleID             string `csv:"Apple Identifier"`
	CustomerPrice       float64 `csv:"Customer Price"`
	PromoCode           string `csv:"Promo Code"`
	ParentID            string `csv:"Parent Identifier"`
	Subscription        string `csv:"Subscription"`
	Period              string `csv:"Period"`
	Category            string `csv:"Category"`
	CMB                 string `csv:"CMB"`
	DeviceType          string `csv:"Device"`
	SupportedPlatforms  string `csv:"Supported Platforms"`
	ProceedsReason      string `csv:"Proceeds Reason"`
	PreservedPricing    string `csv:"Preserved Pricing"`
	Client              string `csv:"Client"`
	OrderType           string `csv:"Order Type"`
}
