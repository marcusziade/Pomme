package sales

import (
	"fmt"
	"time"

	"github.com/marcusziade/pomme/internal/models"
)

// ReportOptions configures a sales report request
type ReportOptions struct {
	Period         models.ReportFrequency
	Date           time.Time
	ReportType     models.ReportType
	VendorNumber   string
	NoCache        bool
	IncludeAnalysis bool
	PreviousPeriod *models.SalesReport // For trend analysis
}

// FormatDate formats the date according to the period
func (o ReportOptions) FormatDate() string {
	switch o.Period {
	case models.ReportFrequencyDaily:
		return o.Date.Format("2006-01-02")
	case models.ReportFrequencyWeekly:
		// Apple weeks end on Saturday
		return o.Date.Format("2006-01-02")
	case models.ReportFrequencyMonthly:
		return o.Date.Format("2006-01")
	case models.ReportFrequencyYearly:
		return o.Date.Format("2006")
	default:
		return o.Date.Format("2006-01-02")
	}
}

// CacheKey generates a unique cache key for the report
func (o ReportOptions) CacheKey() string {
	return fmt.Sprintf("sales:%s:%s:%s:%s",
		o.Period,
		o.FormatDate(),
		o.ReportType,
		o.VendorNumber,
	)
}

// TrendOptions configures trend analysis
type TrendOptions struct {
	Frequency    models.ReportFrequency
	EndDate      time.Time
	Periods      int // Number of periods to analyze
	ReportType   models.ReportType
	VendorNumber string
	GroupBy      string // "app", "country", "platform"
}

// Comparison represents a comparison between two reports
type Comparison struct {
	Current         *models.SalesReport
	Previous        *models.SalesReport
	UnitsChange     float64 // Percentage
	ProceedsChange  map[string]float64 // Currency -> Percentage
	NewApps         []string
	RemovedApps     []string
	TopGainers      []AppChange
	TopLosers       []AppChange
	CountryChanges  map[string]CountryChange
}

// AppChange represents change in app performance
type AppChange struct {
	AppID          string
	AppName        string
	UnitsChange    int
	UnitsPercent   float64
	ProceedsChange map[string]float64
}

// CountryChange represents change in country performance
type CountryChange struct {
	Country        string
	UnitsChange    int
	UnitsPercent   float64
	ProceedsChange map[string]float64
}

// TrendReport represents trends over multiple periods
type TrendReport struct {
	Periods        []time.Time
	Frequency      models.ReportFrequency
	TotalUnits     []int
	TotalProceeds  map[string][]float64 // Currency -> Values per period
	AppTrends      map[string]*AppTrend
	CountryTrends  map[string]*CountryTrend
	Insights       []Insight
}

// AppTrend represents an app's performance over time
type AppTrend struct {
	AppID     string
	AppName   string
	Units     []int
	Proceeds  map[string][]float64
	Growth    float64 // Overall growth rate
	Stability float64 // Variance measure
}

// CountryTrend represents a country's performance over time
type CountryTrend struct {
	Country   string
	Units     []int
	Proceeds  map[string][]float64
	Growth    float64
}

// Insight represents an automatically generated insight
type Insight struct {
	Type        InsightType
	Severity    InsightSeverity
	Title       string
	Description string
	Data        map[string]interface{}
}

// InsightType categorizes insights
type InsightType string

const (
	InsightGrowth      InsightType = "growth"
	InsightDecline     InsightType = "decline"
	InsightAnomaly     InsightType = "anomaly"
	InsightOpportunity InsightType = "opportunity"
	InsightWarning     InsightType = "warning"
)

// InsightSeverity indicates the importance of an insight
type InsightSeverity string

const (
	InsightSeverityInfo     InsightSeverity = "info"
	InsightSeveritySuccess  InsightSeverity = "success"
	InsightSeverityWarning  InsightSeverity = "warning"
	InsightSeverityCritical InsightSeverity = "critical"
)

// ExportOptions configures report export
type ExportOptions struct {
	Format       ExportFormat
	IncludeCharts bool
	GroupBy      string
	Currencies   []string // Filter specific currencies
}

// ExportFormat specifies the export format
type ExportFormat string

const (
	ExportCSV   ExportFormat = "csv"
	ExportJSON  ExportFormat = "json"
	ExportExcel ExportFormat = "excel"
	ExportPDF   ExportFormat = "pdf"
)

// FilterOptions allows filtering sales data
type FilterOptions struct {
	Apps         []string
	Countries    []string
	MinUnits     int
	MinProceeds  float64
	Currency     string
	ProductTypes []string
	Platforms    []string
	DateRange    *DateRange
}

// DateRange represents a date range
type DateRange struct {
	Start time.Time
	End   time.Time
}

// Validate checks if the date range is valid
func (d *DateRange) Validate() error {
	if d.Start.After(d.End) {
		return fmt.Errorf("start date must be before end date")
	}
	if d.End.After(time.Now()) {
		return fmt.Errorf("end date cannot be in the future")
	}
	return nil
}