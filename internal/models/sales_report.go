package models

import (
	"fmt"
	"time"
)

// SalesReport represents a complete sales report with metadata
type SalesReport struct {
	Period      ReportFrequency
	Date        time.Time
	VendorID    string
	Apps        []AppSales
	Summary     ReportSummary
	GeneratedAt time.Time
}

// AppSales represents aggregated sales data for a single app
type AppSales struct {
	AppID       string
	AppName     string
	SKU         string
	Icon        string // Optional app icon URL
	Sales       []Sale
	Summary     AppSummary
}

// Sale represents a single sales transaction
type Sale struct {
	Date              time.Time
	Country           string
	CountryName       string
	Units             int
	CustomerPrice     Money
	DeveloperProceeds Money
	ProductType       string
	Platform          string
	Device            string
	PromoCode         string
	ParentID          string
	Category          string
}

// Money represents a monetary value with currency
type Money struct {
	Amount   float64
	Currency string
}

// AppSummary provides aggregated metrics for an app
type AppSummary struct {
	TotalUnits     int
	TotalProceeds  map[string]float64 // Currency -> Amount
	Countries      int
	AvgPrice       map[string]float64 // Currency -> Average Price
	TopCountries   []CountrySales
	PlatformSplit  map[string]int // Platform -> Units
	DeviceSplit    map[string]int // Device -> Units
}

// CountrySales represents sales data for a specific country
type CountrySales struct {
	Country   string
	CountryName string
	Units     int
	Proceeds  map[string]float64 // Currency -> Amount
}

// ReportSummary provides overall report statistics
type ReportSummary struct {
	TotalApps     int
	TotalUnits    int
	TotalProceeds map[string]float64 // Currency -> Amount
	TotalCountries int
	Period        string
	TopApps       []AppRanking
	TopCountries  []CountrySales
	Trends        *TrendAnalysis
}

// AppRanking represents an app's ranking in the report
type AppRanking struct {
	AppID     string
	AppName   string
	Units     int
	Proceeds  map[string]float64
	Rank      int
	Change    int // Position change from previous period
}

// TrendAnalysis provides trend insights
type TrendAnalysis struct {
	UnitsTrend      TrendDirection
	ProceedsTrend   TrendDirection
	UnitsChange     float64 // Percentage
	ProceedsChange  float64 // Percentage
	NewCountries    []string
	LostCountries   []string
	BestPerformer   *AppRanking
	WorstPerformer  *AppRanking
}

// TrendDirection indicates the direction of a trend
type TrendDirection int

const (
	TrendDown TrendDirection = -1
	TrendFlat TrendDirection = 0
	TrendUp   TrendDirection = 1
)

// String returns a formatted string representation of Money
func (m Money) String() string {
	if m.Currency == "" {
		return fmt.Sprintf("%.2f", m.Amount)
	}
	return fmt.Sprintf("%s %.2f", m.Currency, m.Amount)
}

// IsZero checks if the money amount is zero
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// Add adds two money values of the same currency
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency && m.Currency != "" && other.Currency != "" {
		return Money{}, fmt.Errorf("cannot add different currencies: %s and %s", m.Currency, other.Currency)
	}
	
	currency := m.Currency
	if currency == "" {
		currency = other.Currency
	}
	
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: currency,
	}, nil
}

// FormatCurrency formats the currency map for display
func FormatCurrency(amounts map[string]float64) string {
	if len(amounts) == 0 {
		return "No proceeds"
	}
	
	result := ""
	for currency, amount := range amounts {
		if result != "" {
			result += ", "
		}
		result += fmt.Sprintf("%s %.2f", currency, amount)
	}
	return result
}

// GetTrendIcon returns an icon for the trend direction
func (t TrendDirection) GetIcon() string {
	switch t {
	case TrendUp:
		return "↑"
	case TrendDown:
		return "↓"
	default:
		return "→"
	}
}

// GetTrendColor returns a color code for the trend
func (t TrendDirection) GetColor() string {
	switch t {
	case TrendUp:
		return "green"
	case TrendDown:
		return "red"
	default:
		return "yellow"
	}
}