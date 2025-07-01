package sales

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/marcusziade/pomme/internal/api"
	"github.com/marcusziade/pomme/internal/models"
	"github.com/marcusziade/pomme/internal/services/cache"
)

// Service handles all sales-related operations
type Service struct {
	client      *api.Client
	cache       cache.Cache
	parser      *Parser
	analyzer    *Analyzer
	concurrency int
}

// NewService creates a new sales service
func NewService(client *api.Client, cacheService cache.Cache) *Service {
	return &Service{
		client:      client,
		cache:       cacheService,
		parser:      NewParser(),
		analyzer:    NewAnalyzer(),
		concurrency: 4, // Default concurrent operations
	}
}

// GetReport fetches and processes a sales report
func (s *Service) GetReport(ctx context.Context, options ReportOptions) (*models.SalesReport, error) {
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
	rawData, err := s.fetchReport(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch report: %w", err)
	}

	// Parse the report concurrently
	report, err := s.parseReport(ctx, rawData, options)
	if err != nil {
		return nil, fmt.Errorf("failed to parse report: %w", err)
	}

	// Analyze the data
	if options.IncludeAnalysis {
		report.Summary.Trends = s.analyzer.AnalyzeTrends(report, options.PreviousPeriod)
	}

	// Cache the result
	if s.cache != nil && !options.NoCache {
		s.cache.Set(cacheKey, report, 24*time.Hour)
	}

	return report, nil
}

// GetMultipleReports fetches multiple reports concurrently
func (s *Service) GetMultipleReports(ctx context.Context, requests []ReportOptions) ([]*models.SalesReport, error) {
	results := make([]*models.SalesReport, len(requests))
	errors := make([]error, len(requests))
	
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.concurrency)
	
	for i, req := range requests {
		wg.Add(1)
		go func(idx int, options ReportOptions) {
			defer wg.Done()
			
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			report, err := s.GetReport(ctx, options)
			results[idx] = report
			errors[idx] = err
		}(i, req)
	}
	
	wg.Wait()
	
	// Check for errors
	var errs []string
	for i, err := range errors {
		if err != nil {
			errs = append(errs, fmt.Sprintf("report %d: %v", i, err))
		}
	}
	
	if len(errs) > 0 {
		return results, fmt.Errorf("multiple errors: %s", strings.Join(errs, "; "))
	}
	
	return results, nil
}

// GetComparison compares two periods
func (s *Service) GetComparison(ctx context.Context, current, previous ReportOptions) (*Comparison, error) {
	// Fetch both reports concurrently
	var currentReport, previousReport *models.SalesReport
	var currentErr, previousErr error
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	go func() {
		defer wg.Done()
		currentReport, currentErr = s.GetReport(ctx, current)
	}()
	
	go func() {
		defer wg.Done()
		previousReport, previousErr = s.GetReport(ctx, previous)
	}()
	
	wg.Wait()
	
	if currentErr != nil {
		return nil, fmt.Errorf("failed to get current report: %w", currentErr)
	}
	if previousErr != nil {
		return nil, fmt.Errorf("failed to get previous report: %w", previousErr)
	}
	
	return s.analyzer.Compare(currentReport, previousReport), nil
}

// GetTrends analyzes trends over multiple periods
func (s *Service) GetTrends(ctx context.Context, options TrendOptions) (*TrendReport, error) {
	// Generate report options for each period
	requests := s.generateTrendRequests(options)
	
	// Fetch all reports concurrently
	reports, err := s.GetMultipleReports(ctx, requests)
	if err != nil {
		return nil, err
	}
	
	// Analyze trends
	return s.analyzer.AnalyzeTrendSeries(reports, options), nil
}

// fetchReport retrieves raw report data from the API
func (s *Service) fetchReport(ctx context.Context, options ReportOptions) ([]byte, error) {
	// This would be implemented when we have the proper API client integration
	return nil, fmt.Errorf("not implemented")
}

// parseReport parses raw CSV data into a structured report
func (s *Service) parseReport(ctx context.Context, data []byte, options ReportOptions) (*models.SalesReport, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty report data")
	}
	
	records, err := s.parser.ParseCSV(data)
	if err != nil {
		return nil, err
	}
	
	// Process records concurrently by app
	appMap := s.groupRecordsByApp(records)
	appSales := make([]models.AppSales, 0, len(appMap))
	
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.concurrency)
	
	for appID, appRecords := range appMap {
		wg.Add(1)
		go func(id string, records []models.SalesRecord) {
			defer wg.Done()
			
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			sales := s.processAppSales(id, records)
			
			mu.Lock()
			appSales = append(appSales, sales)
			mu.Unlock()
		}(appID, appRecords)
	}
	
	wg.Wait()
	
	// Sort apps by total units
	sort.Slice(appSales, func(i, j int) bool {
		return appSales[i].Summary.TotalUnits > appSales[j].Summary.TotalUnits
	})
	
	// Build report
	report := &models.SalesReport{
		Period:      options.Period,
		Date:        options.Date,
		VendorID:    options.VendorNumber,
		Apps:        appSales,
		GeneratedAt: time.Now(),
	}
	
	// Calculate summary
	report.Summary = s.calculateReportSummary(report)
	
	return report, nil
}

// groupRecordsByApp groups sales records by app ID
func (s *Service) groupRecordsByApp(records []models.SalesRecord) map[string][]models.SalesRecord {
	appMap := make(map[string][]models.SalesRecord)
	
	for _, record := range records {
		appMap[record.AppleID] = append(appMap[record.AppleID], record)
	}
	
	return appMap
}

// processAppSales processes sales records for a single app
func (s *Service) processAppSales(appID string, records []models.SalesRecord) models.AppSales {
	if len(records) == 0 {
		return models.AppSales{AppID: appID}
	}
	
	// Use first record for app metadata
	first := records[0]
	appSales := models.AppSales{
		AppID:   appID,
		AppName: first.Title,
		SKU:     first.SKU,
		Sales:   make([]models.Sale, 0, len(records)),
	}
	
	// Convert records to sales
	summary := models.AppSummary{
		TotalProceeds: make(map[string]float64),
		AvgPrice:      make(map[string]float64),
		PlatformSplit: make(map[string]int),
		DeviceSplit:   make(map[string]int),
	}
	
	countryMap := make(map[string]*models.CountrySales)
	priceSum := make(map[string]float64)
	priceCount := make(map[string]int)
	
	for _, record := range records {
		sale := models.Sale{
			Date:    s.parser.ParseDate(record.BeginDate),
			Country: record.CountryCode,
			Units:   record.Units,
			CustomerPrice: models.Money{
				Amount:   record.CustomerPrice,
				Currency: record.CustomerCurrency,
			},
			DeveloperProceeds: models.Money{
				Amount:   record.DeveloperProceeds,
				Currency: record.CurrencyOfProceeds,
			},
			ProductType: record.ProductTypeID,
			Platform:    record.SupportedPlatforms,
			Device:      record.DeviceType,
			PromoCode:   record.PromoCode,
			ParentID:    record.ParentID,
			Category:    record.Category,
		}
		
		appSales.Sales = append(appSales.Sales, sale)
		
		// Update summary
		summary.TotalUnits += sale.Units
		
		if sale.DeveloperProceeds.Amount > 0 && sale.DeveloperProceeds.Currency != "" {
			summary.TotalProceeds[sale.DeveloperProceeds.Currency] += sale.DeveloperProceeds.Amount
		}
		
		if sale.CustomerPrice.Amount > 0 && sale.CustomerPrice.Currency != "" {
			priceSum[sale.CustomerPrice.Currency] += sale.CustomerPrice.Amount * float64(sale.Units)
			priceCount[sale.CustomerPrice.Currency] += sale.Units
		}
		
		// Platform and device splits
		if sale.Platform != "" {
			summary.PlatformSplit[sale.Platform] += sale.Units
		}
		if sale.Device != "" {
			summary.DeviceSplit[sale.Device] += sale.Units
		}
		
		// Country aggregation
		if _, exists := countryMap[sale.Country]; !exists {
			countryMap[sale.Country] = &models.CountrySales{
				Country:   sale.Country,
				CountryName: s.getCountryName(sale.Country),
				Proceeds:  make(map[string]float64),
			}
		}
		countryMap[sale.Country].Units += sale.Units
		if sale.DeveloperProceeds.Amount > 0 && sale.DeveloperProceeds.Currency != "" {
			countryMap[sale.Country].Proceeds[sale.DeveloperProceeds.Currency] += sale.DeveloperProceeds.Amount
		}
	}
	
	// Calculate averages
	for currency, sum := range priceSum {
		if count := priceCount[currency]; count > 0 {
			summary.AvgPrice[currency] = sum / float64(count)
		}
	}
	
	// Convert country map to slice and sort
	summary.Countries = len(countryMap)
	topCountries := make([]models.CountrySales, 0, len(countryMap))
	for _, country := range countryMap {
		topCountries = append(topCountries, *country)
	}
	
	sort.Slice(topCountries, func(i, j int) bool {
		return topCountries[i].Units > topCountries[j].Units
	})
	
	// Keep top 5 countries
	if len(topCountries) > 5 {
		summary.TopCountries = topCountries[:5]
	} else {
		summary.TopCountries = topCountries
	}
	
	appSales.Summary = summary
	return appSales
}

// calculateReportSummary calculates the overall report summary
func (s *Service) calculateReportSummary(report *models.SalesReport) models.ReportSummary {
	summary := models.ReportSummary{
		TotalApps:     len(report.Apps),
		TotalProceeds: make(map[string]float64),
		Period:        report.Period.String(),
	}
	
	countrySet := make(map[string]bool)
	
	for _, app := range report.Apps {
		summary.TotalUnits += app.Summary.TotalUnits
		
		// Aggregate proceeds
		for currency, amount := range app.Summary.TotalProceeds {
			summary.TotalProceeds[currency] += amount
		}
		
		// Count unique countries
		for _, sale := range app.Sales {
			countrySet[sale.Country] = true
		}
	}
	
	summary.TotalCountries = len(countrySet)
	
	// Top apps (already sorted)
	for i, app := range report.Apps {
		if i >= 5 {
			break
		}
		summary.TopApps = append(summary.TopApps, models.AppRanking{
			AppID:    app.AppID,
			AppName:  app.AppName,
			Units:    app.Summary.TotalUnits,
			Proceeds: app.Summary.TotalProceeds,
			Rank:     i + 1,
		})
	}
	
	return summary
}

// getCountryName returns the full country name for a country code
func (s *Service) getCountryName(code string) string {
	// This would be populated from a proper country code mapping
	countryNames := map[string]string{
		"US": "United States",
		"GB": "United Kingdom",
		"DE": "Germany",
		"FR": "France",
		"JP": "Japan",
		"CN": "China",
		"CA": "Canada",
		"AU": "Australia",
		// ... add more as needed
	}
	
	if name, ok := countryNames[code]; ok {
		return name
	}
	return code
}

// generateTrendRequests generates report requests for trend analysis
func (s *Service) generateTrendRequests(options TrendOptions) []ReportOptions {
	requests := make([]ReportOptions, options.Periods)
	
	for i := 0; i < options.Periods; i++ {
		date := options.EndDate
		
		switch options.Frequency {
		case models.ReportFrequencyDaily:
			date = date.AddDate(0, 0, -i)
		case models.ReportFrequencyWeekly:
			date = date.AddDate(0, 0, -i*7)
		case models.ReportFrequencyMonthly:
			date = date.AddDate(0, -i, 0)
		case models.ReportFrequencyYearly:
			date = date.AddDate(-i, 0, 0)
		}
		
		requests[i] = ReportOptions{
			Period:       options.Frequency,
			Date:         date,
			ReportType:   options.ReportType,
			VendorNumber: options.VendorNumber,
		}
	}
	
	return requests
}