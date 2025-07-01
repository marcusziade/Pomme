package sales

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/marcusziade/pomme/internal/models"
)

// Analyzer provides analysis and insights for sales data
type Analyzer struct{}

// NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeTrends analyzes trends between current and previous reports
func (a *Analyzer) AnalyzeTrends(current *models.SalesReport, previous *models.SalesReport) *models.TrendAnalysis {
	if previous == nil {
		return nil
	}

	trends := &models.TrendAnalysis{
		NewCountries:  []string{},
		LostCountries: []string{},
	}

	// Calculate units trend
	prevUnits := previous.Summary.TotalUnits
	currUnits := current.Summary.TotalUnits
	
	if prevUnits > 0 {
		trends.UnitsChange = ((float64(currUnits) - float64(prevUnits)) / float64(prevUnits)) * 100
		trends.UnitsTrend = a.getTrendDirection(trends.UnitsChange)
	}

	// Calculate proceeds trend (use primary currency)
	var primaryCurrency string
	var maxProceeds float64
	
	for currency, amount := range current.Summary.TotalProceeds {
		if amount > maxProceeds {
			primaryCurrency = currency
			maxProceeds = amount
		}
	}
	
	if primaryCurrency != "" && previous.Summary.TotalProceeds[primaryCurrency] > 0 {
		prevProceeds := previous.Summary.TotalProceeds[primaryCurrency]
		currProceeds := current.Summary.TotalProceeds[primaryCurrency]
		trends.ProceedsChange = ((currProceeds - prevProceeds) / prevProceeds) * 100
		trends.ProceedsTrend = a.getTrendDirection(trends.ProceedsChange)
	}

	// Find new and lost countries
	currCountries := a.getCountrySet(current)
	prevCountries := a.getCountrySet(previous)
	
	for country := range currCountries {
		if !prevCountries[country] {
			trends.NewCountries = append(trends.NewCountries, country)
		}
	}
	
	for country := range prevCountries {
		if !currCountries[country] {
			trends.LostCountries = append(trends.LostCountries, country)
		}
	}

	// Find best and worst performers
	appPerformance := a.calculateAppPerformance(current, previous)
	
	if len(appPerformance) > 0 {
		// Sort by units change
		sort.Slice(appPerformance, func(i, j int) bool {
			return appPerformance[i].Change > appPerformance[j].Change
		})
		
		trends.BestPerformer = &appPerformance[0]
		if len(appPerformance) > 1 {
			trends.WorstPerformer = &appPerformance[len(appPerformance)-1]
		}
	}

	return trends
}

// Compare creates a detailed comparison between two reports
func (a *Analyzer) Compare(current, previous *models.SalesReport) *Comparison {
	comp := &Comparison{
		Current:        current,
		Previous:       previous,
		ProceedsChange: make(map[string]float64),
		CountryChanges: make(map[string]CountryChange),
	}

	// Calculate overall changes
	if previous.Summary.TotalUnits > 0 {
		comp.UnitsChange = ((float64(current.Summary.TotalUnits) - float64(previous.Summary.TotalUnits)) / 
			float64(previous.Summary.TotalUnits)) * 100
	}

	// Calculate proceeds changes by currency
	for currency, currAmount := range current.Summary.TotalProceeds {
		if prevAmount, ok := previous.Summary.TotalProceeds[currency]; ok && prevAmount > 0 {
			comp.ProceedsChange[currency] = ((currAmount - prevAmount) / prevAmount) * 100
		}
	}

	// Find new and removed apps
	currApps := make(map[string]bool)
	prevApps := make(map[string]bool)
	
	for _, app := range current.Apps {
		currApps[app.AppID] = true
	}
	
	for _, app := range previous.Apps {
		prevApps[app.AppID] = true
		if !currApps[app.AppID] {
			comp.RemovedApps = append(comp.RemovedApps, app.AppName)
		}
	}
	
	for _, app := range current.Apps {
		if !prevApps[app.AppID] {
			comp.NewApps = append(comp.NewApps, app.AppName)
		}
	}

	// Calculate app performance changes
	appChanges := a.calculateAppChanges(current, previous)
	
	// Sort by units change
	sort.Slice(appChanges, func(i, j int) bool {
		return appChanges[i].UnitsChange > appChanges[j].UnitsChange
	})
	
	// Get top gainers and losers
	for i, change := range appChanges {
		if i < 5 && change.UnitsChange > 0 {
			comp.TopGainers = append(comp.TopGainers, change)
		}
		if i >= len(appChanges)-5 && change.UnitsChange < 0 {
			comp.TopLosers = append(comp.TopLosers, change)
		}
	}

	// Calculate country changes
	comp.CountryChanges = a.calculateCountryChanges(current, previous)

	return comp
}

// AnalyzeTrendSeries analyzes trends across multiple periods
func (a *Analyzer) AnalyzeTrendSeries(reports []*models.SalesReport, options TrendOptions) *TrendReport {
	if len(reports) == 0 {
		return nil
	}

	trend := &TrendReport{
		Periods:       make([]time.Time, len(reports)),
		Frequency:     options.Frequency,
		TotalUnits:    make([]int, len(reports)),
		TotalProceeds: make(map[string][]float64),
		AppTrends:     make(map[string]*AppTrend),
		CountryTrends: make(map[string]*CountryTrend),
		Insights:      []Insight{},
	}

	// Initialize currency arrays
	currencies := a.getAllCurrencies(reports)
	for currency := range currencies {
		trend.TotalProceeds[currency] = make([]float64, len(reports))
	}

	// Process each report
	for i, report := range reports {
		if report == nil {
			continue
		}
		
		trend.Periods[i] = report.Date
		trend.TotalUnits[i] = report.Summary.TotalUnits
		
		// Track proceeds by currency
		for currency, amount := range report.Summary.TotalProceeds {
			trend.TotalProceeds[currency][i] = amount
		}
		
		// Track app trends
		for _, app := range report.Apps {
			if _, exists := trend.AppTrends[app.AppID]; !exists {
				trend.AppTrends[app.AppID] = &AppTrend{
					AppID:    app.AppID,
					AppName:  app.AppName,
					Units:    make([]int, len(reports)),
					Proceeds: make(map[string][]float64),
				}
				
				// Initialize proceeds arrays
				for currency := range currencies {
					trend.AppTrends[app.AppID].Proceeds[currency] = make([]float64, len(reports))
				}
			}
			
			trend.AppTrends[app.AppID].Units[i] = app.Summary.TotalUnits
			
			for currency, amount := range app.Summary.TotalProceeds {
				trend.AppTrends[app.AppID].Proceeds[currency][i] = amount
			}
		}
	}

	// Calculate growth rates and stability
	for _, appTrend := range trend.AppTrends {
		appTrend.Growth = a.calculateGrowthRate(appTrend.Units)
		appTrend.Stability = a.calculateStability(appTrend.Units)
	}

	// Generate insights
	trend.Insights = a.generateInsights(trend)

	return trend
}

// getTrendDirection determines the trend direction based on percentage change
func (a *Analyzer) getTrendDirection(change float64) models.TrendDirection {
	if change > 5 {
		return models.TrendUp
	} else if change < -5 {
		return models.TrendDown
	}
	return models.TrendFlat
}

// getCountrySet returns a set of countries from a report
func (a *Analyzer) getCountrySet(report *models.SalesReport) map[string]bool {
	countries := make(map[string]bool)
	
	for _, app := range report.Apps {
		for _, sale := range app.Sales {
			countries[sale.Country] = true
		}
	}
	
	return countries
}

// calculateAppPerformance calculates performance changes for each app
func (a *Analyzer) calculateAppPerformance(current, previous *models.SalesReport) []models.AppRanking {
	prevAppMap := make(map[string]*models.AppSales)
	for _, app := range previous.Apps {
		prevAppMap[app.AppID] = &app
	}

	var rankings []models.AppRanking
	
	for i, app := range current.Apps {
		ranking := models.AppRanking{
			AppID:    app.AppID,
			AppName:  app.AppName,
			Units:    app.Summary.TotalUnits,
			Proceeds: app.Summary.TotalProceeds,
			Rank:     i + 1,
		}
		
		if prevApp, ok := prevAppMap[app.AppID]; ok {
			// Find previous rank
			for j, pApp := range previous.Apps {
				if pApp.AppID == app.AppID {
					ranking.Change = (j + 1) - ranking.Rank
					break
				}
			}
			
			// Calculate units change
			if prevApp.Summary.TotalUnits > 0 {
				change := float64(app.Summary.TotalUnits - prevApp.Summary.TotalUnits)
				ranking.Change = int(change)
			}
		}
		
		rankings = append(rankings, ranking)
	}
	
	return rankings
}

// calculateAppChanges calculates detailed changes for each app
func (a *Analyzer) calculateAppChanges(current, previous *models.SalesReport) []AppChange {
	currAppMap := make(map[string]*models.AppSales)
	for _, app := range current.Apps {
		currAppMap[app.AppID] = &app
	}
	
	prevAppMap := make(map[string]*models.AppSales)
	for _, app := range previous.Apps {
		prevAppMap[app.AppID] = &app
	}

	var changes []AppChange
	
	// Process all apps that appear in either report
	processedApps := make(map[string]bool)
	
	for appID, currApp := range currAppMap {
		change := AppChange{
			AppID:          appID,
			AppName:        currApp.AppName,
			ProceedsChange: make(map[string]float64),
		}
		
		if prevApp, ok := prevAppMap[appID]; ok {
			change.UnitsChange = currApp.Summary.TotalUnits - prevApp.Summary.TotalUnits
			
			if prevApp.Summary.TotalUnits > 0 {
				change.UnitsPercent = (float64(change.UnitsChange) / float64(prevApp.Summary.TotalUnits)) * 100
			}
			
			// Calculate proceeds changes
			for currency, currAmount := range currApp.Summary.TotalProceeds {
				if prevAmount, ok := prevApp.Summary.TotalProceeds[currency]; ok && prevAmount > 0 {
					change.ProceedsChange[currency] = ((currAmount - prevAmount) / prevAmount) * 100
				}
			}
		} else {
			// New app
			change.UnitsChange = currApp.Summary.TotalUnits
			change.UnitsPercent = 100 // New app is 100% growth
		}
		
		changes = append(changes, change)
		processedApps[appID] = true
	}
	
	// Add removed apps
	for appID, prevApp := range prevAppMap {
		if !processedApps[appID] {
			change := AppChange{
				AppID:          appID,
				AppName:        prevApp.AppName,
				UnitsChange:    -prevApp.Summary.TotalUnits,
				UnitsPercent:   -100, // Removed app is -100%
				ProceedsChange: make(map[string]float64),
			}
			
			for currency := range prevApp.Summary.TotalProceeds {
				change.ProceedsChange[currency] = -100
			}
			
			changes = append(changes, change)
		}
	}
	
	return changes
}

// calculateCountryChanges calculates performance changes by country
func (a *Analyzer) calculateCountryChanges(current, previous *models.SalesReport) map[string]CountryChange {
	currCountryData := a.aggregateByCountry(current)
	prevCountryData := a.aggregateByCountry(previous)
	
	changes := make(map[string]CountryChange)
	
	for country, currData := range currCountryData {
		change := CountryChange{
			Country:        country,
			ProceedsChange: make(map[string]float64),
		}
		
		if prevData, ok := prevCountryData[country]; ok {
			change.UnitsChange = currData.Units - prevData.Units
			
			if prevData.Units > 0 {
				change.UnitsPercent = (float64(change.UnitsChange) / float64(prevData.Units)) * 100
			}
			
			// Calculate proceeds changes
			for currency, currAmount := range currData.Proceeds {
				if prevAmount, ok := prevData.Proceeds[currency]; ok && prevAmount > 0 {
					change.ProceedsChange[currency] = ((currAmount - prevAmount) / prevAmount) * 100
				}
			}
		} else {
			// New country
			change.UnitsChange = currData.Units
			change.UnitsPercent = 100
		}
		
		changes[country] = change
	}
	
	return changes
}

// aggregateByCountry aggregates sales data by country
func (a *Analyzer) aggregateByCountry(report *models.SalesReport) map[string]*models.CountrySales {
	countryData := make(map[string]*models.CountrySales)
	
	for _, app := range report.Apps {
		for _, sale := range app.Sales {
			if _, exists := countryData[sale.Country]; !exists {
				countryData[sale.Country] = &models.CountrySales{
					Country:  sale.Country,
					Proceeds: make(map[string]float64),
				}
			}
			
			countryData[sale.Country].Units += sale.Units
			
			if sale.DeveloperProceeds.Amount > 0 && sale.DeveloperProceeds.Currency != "" {
				countryData[sale.Country].Proceeds[sale.DeveloperProceeds.Currency] += sale.DeveloperProceeds.Amount
			}
		}
	}
	
	return countryData
}

// getAllCurrencies returns all currencies present in the reports
func (a *Analyzer) getAllCurrencies(reports []*models.SalesReport) map[string]bool {
	currencies := make(map[string]bool)
	
	for _, report := range reports {
		if report != nil {
			for currency := range report.Summary.TotalProceeds {
				currencies[currency] = true
			}
		}
	}
	
	return currencies
}

// calculateGrowthRate calculates the compound growth rate
func (a *Analyzer) calculateGrowthRate(values []int) float64 {
	if len(values) < 2 {
		return 0
	}
	
	// Find first and last non-zero values
	var firstValue, lastValue float64
	var firstIndex, lastIndex int
	
	for i, v := range values {
		if v > 0 {
			if firstValue == 0 {
				firstValue = float64(v)
				firstIndex = i
			}
			lastValue = float64(v)
			lastIndex = i
		}
	}
	
	if firstValue == 0 || lastValue == 0 || lastIndex <= firstIndex {
		return 0
	}
	
	periods := float64(lastIndex - firstIndex)
	return (math.Pow(lastValue/firstValue, 1/periods) - 1) * 100
}

// calculateStability calculates the coefficient of variation (lower is more stable)
func (a *Analyzer) calculateStability(values []int) float64 {
	if len(values) < 2 {
		return 0
	}
	
	// Calculate mean
	sum := 0.0
	count := 0
	
	for _, v := range values {
		if v > 0 {
			sum += float64(v)
			count++
		}
	}
	
	if count == 0 {
		return 0
	}
	
	mean := sum / float64(count)
	
	// Calculate standard deviation
	variance := 0.0
	for _, v := range values {
		if v > 0 {
			diff := float64(v) - mean
			variance += diff * diff
		}
	}
	
	if count > 1 {
		variance /= float64(count - 1)
	}
	
	stdDev := math.Sqrt(variance)
	
	// Return coefficient of variation
	if mean > 0 {
		return (stdDev / mean) * 100
	}
	
	return 0
}

// generateInsights creates automatic insights from trend data
func (a *Analyzer) generateInsights(trend *TrendReport) []Insight {
	var insights []Insight
	
	// Check overall growth
	if len(trend.TotalUnits) >= 2 {
		firstUnits := trend.TotalUnits[0]
		lastUnits := trend.TotalUnits[len(trend.TotalUnits)-1]
		
		if firstUnits > 0 {
			growth := ((float64(lastUnits) - float64(firstUnits)) / float64(firstUnits)) * 100
			
			if growth > 50 {
				insights = append(insights, Insight{
					Type:     InsightGrowth,
					Severity: InsightSeveritySuccess,
					Title:    "Strong Growth",
					Description: "Your total units have grown by " + formatPercent(growth) + " over the analyzed period.",
					Data: map[string]interface{}{
						"growth": growth,
						"from":   firstUnits,
						"to":     lastUnits,
					},
				})
			} else if growth < -20 {
				insights = append(insights, Insight{
					Type:     InsightDecline,
					Severity: InsightSeverityWarning,
					Title:    "Declining Sales",
					Description: "Your total units have declined by " + formatPercent(-growth) + " over the analyzed period.",
					Data: map[string]interface{}{
						"decline": growth,
						"from":    firstUnits,
						"to":      lastUnits,
					},
				})
			}
		}
	}
	
	// Check for volatile apps
	for appID, appTrend := range trend.AppTrends {
		if appTrend.Stability > 50 {
			insights = append(insights, Insight{
				Type:     InsightAnomaly,
				Severity: InsightSeverityInfo,
				Title:    "Volatile App Performance",
				Description: appTrend.AppName + " shows high volatility in sales performance.",
				Data: map[string]interface{}{
					"appID":      appID,
					"appName":    appTrend.AppName,
					"volatility": appTrend.Stability,
				},
			})
		}
	}
	
	// Check for opportunities
	// Add more insight generation logic here...
	
	return insights
}

// formatPercent formats a percentage value
func formatPercent(value float64) string {
	if value >= 0 {
		return "+" + fmt.Sprintf("%.1f%%", value)
	}
	return fmt.Sprintf("%.1f%%", value)
}