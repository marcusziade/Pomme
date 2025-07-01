package commands

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/marcusziade/pomme/internal/models"
	"github.com/marcusziade/pomme/internal/services/cache"
	"github.com/marcusziade/pomme/internal/services/sales"
	"github.com/marcusziade/pomme/pkg/pomme"
)

// salesService is a temporary adapter until we fully refactor
type salesService struct {
	client *pomme.Client
	cache  cache.Cache
}

// GetReport adapts the existing client to the new service interface
func (s *salesService) GetReport(ctx context.Context, options sales.ReportOptions) (*models.SalesReport, error) {
	// Check cache first
	cacheKey := options.CacheKey()
	if s.cache != nil && !options.NoCache {
		if cached, err := s.cache.Get(cacheKey); err == nil {
			if report, ok := cached.(*models.SalesReport); ok {
				return report, nil
			}
		}
	}

	// Fetch from API
	rawData, err := s.client.GetSalesReport(
		ctx,
		options.Period,
		options.FormatDate(),
		options.ReportType,
		options.VendorNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch report: %w", err)
	}

	if rawData == nil || len(rawData) == 0 {
		// No data available
		return nil, nil
	}

	// Parse the CSV data
	report, err := s.parseReport(rawData, options)
	if err != nil {
		return nil, fmt.Errorf("failed to parse report: %w", err)
	}

	// Cache the result
	if s.cache != nil && !options.NoCache {
		s.cache.Set(cacheKey, report, 24*time.Hour)
	}

	return report, nil
}

// GetComparison fetches and compares two reports
func (s *salesService) GetComparison(ctx context.Context, current, previous sales.ReportOptions) (*sales.Comparison, error) {
	// Fetch both reports
	currentReport, err := s.GetReport(ctx, current)
	if err != nil {
		return nil, fmt.Errorf("failed to get current report: %w", err)
	}

	previousReport, err := s.GetReport(ctx, previous)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous report: %w", err)
	}

	// Simple comparison
	comp := &sales.Comparison{
		Current:        currentReport,
		Previous:       previousReport,
		ProceedsChange: make(map[string]float64),
	}

	if currentReport != nil && previousReport != nil {
		// Calculate changes
		if previousReport.Summary.TotalUnits > 0 {
			comp.UnitsChange = ((float64(currentReport.Summary.TotalUnits) - float64(previousReport.Summary.TotalUnits)) / 
				float64(previousReport.Summary.TotalUnits)) * 100
		}

		// Calculate proceeds changes
		for currency, currAmount := range currentReport.Summary.TotalProceeds {
			if prevAmount, ok := previousReport.Summary.TotalProceeds[currency]; ok && prevAmount > 0 {
				comp.ProceedsChange[currency] = ((currAmount - prevAmount) / prevAmount) * 100
			}
		}
	}

	return comp, nil
}

// GetTrends is not implemented in the adapter
func (s *salesService) GetTrends(ctx context.Context, options sales.TrendOptions) (*sales.TrendReport, error) {
	return nil, fmt.Errorf("trends analysis not yet implemented")
}

// parseReport parses CSV data into a sales report
func (s *salesService) parseReport(data []byte, options sales.ReportOptions) (*models.SalesReport, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	reader.Comma = '\t'
	reader.LazyQuotes = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create field map
	fieldMap := make(map[string]int)
	for i, field := range header {
		fieldMap[strings.TrimSpace(field)] = i
	}

	// Process records
	recordMap := make(map[string][]models.SalesRecord)
	
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip bad rows
		}

		record := parseRecord(row, fieldMap)
		if record.AppleID != "" {
			recordMap[record.AppleID] = append(recordMap[record.AppleID], record)
		}
	}

	// Build report
	report := &models.SalesReport{
		Period:      options.Period,
		Date:        options.Date,
		VendorID:    options.VendorNumber,
		GeneratedAt: time.Now(),
		Apps:        make([]models.AppSales, 0, len(recordMap)),
		Summary: models.ReportSummary{
			TotalProceeds: make(map[string]float64),
		},
	}

	// Process each app
	for appID, records := range recordMap {
		if len(records) == 0 {
			continue
		}

		appSales := models.AppSales{
			AppID:   appID,
			AppName: records[0].Title,
			SKU:     records[0].SKU,
			Summary: models.AppSummary{
				TotalProceeds: make(map[string]float64),
			},
		}

		// Aggregate data
		for _, record := range records {
			appSales.Summary.TotalUnits += record.Units
			
			if record.DeveloperProceeds > 0 && record.CurrencyOfProceeds != "" {
				appSales.Summary.TotalProceeds[record.CurrencyOfProceeds] += record.DeveloperProceeds
			}
		}

		report.Apps = append(report.Apps, appSales)
		
		// Update report summary
		report.Summary.TotalUnits += appSales.Summary.TotalUnits
		for currency, amount := range appSales.Summary.TotalProceeds {
			report.Summary.TotalProceeds[currency] += amount
		}
	}

	report.Summary.TotalApps = len(report.Apps)
	report.Summary.Period = options.Period.String()

	return report, nil
}

// parseRecord parses a CSV row into a SalesRecord
func parseRecord(row []string, fieldMap map[string]int) models.SalesRecord {
	getValue := func(field string) string {
		if idx, ok := fieldMap[field]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	record := models.SalesRecord{
		Provider:           getValue("Provider"),
		ProviderCountry:    getValue("Provider Country"),
		SKU:                getValue("SKU"),
		Developer:          getValue("Developer"),
		Title:              getValue("Title"),
		Version:            getValue("Version"),
		ProductTypeID:      getValue("Product Type Identifier"),
		AppleID:            getValue("Apple Identifier"),
		CountryCode:        getValue("Country Code"),
		CustomerCurrency:   getValue("Customer Currency"),
		CurrencyOfProceeds: getValue("Currency of Proceeds"),
		BeginDate:          getValue("Begin Date"),
		EndDate:            getValue("End Date"),
	}

	// Parse numeric fields
	if units := getValue("Units"); units != "" {
		record.Units, _ = strconv.Atoi(units)
	}

	if price := getValue("Customer Price"); price != "" {
		record.CustomerPrice, _ = strconv.ParseFloat(price, 64)
	}

	if proceeds := getValue("Developer Proceeds"); proceeds != "" {
		record.DeveloperProceeds, _ = strconv.ParseFloat(proceeds, 64)
	}

	return record
}