package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/marcusziade/pomme/internal/models"
	"github.com/marcusziade/pomme/internal/services/sales"
	"github.com/spf13/cobra"
)

const (
	// ANSI color codes
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// displayMonthlyReport shows a beautiful monthly sales report
func displayMonthlyReport(cmd *cobra.Command, report *models.SalesReport) {
	// Header
	fmt.Printf("\n%s📊 Sales Report for %s%s\n", colorBold, report.Date.Format("January 2006"), colorReset)
	fmt.Println(strings.Repeat("─", 60))

	// Summary metrics
	displaySummaryMetrics(report.Summary)

	// App performance
	if len(report.Apps) > 0 {
		fmt.Printf("\n%s📱 App Performance%s\n", colorBold, colorReset)
		fmt.Println(strings.Repeat("─", 60))
		
		displayAppTable(report.Apps)
	}

	// Country breakdown if requested
	if mustGetBool(cmd, "by-country") {
		displayCountryBreakdown(report)
	}

	// Trends and insights
	if report.Summary.Trends != nil {
		displayTrends(report.Summary.Trends)
	}

	// Footer
	fmt.Printf("\n%sGenerated: %s%s\n", colorGray, report.GeneratedAt.Format("2006-01-02 15:04:05"), colorReset)
}

// displaySummaryMetrics shows key metrics in a card-like format
func displaySummaryMetrics(summary models.ReportSummary) {
	// Total units card
	fmt.Printf("\n  %s📦 Total Units%s\n", colorBold, colorReset)
	fmt.Printf("  %s%s%s\n", colorCyan, formatNumber(summary.TotalUnits), colorReset)

	// Revenue cards by currency
	if len(summary.TotalProceeds) > 0 {
		fmt.Printf("\n  %s💰 Revenue%s\n", colorBold, colorReset)
		
		// Sort currencies for consistent display
		currencies := make([]string, 0, len(summary.TotalProceeds))
		for currency := range summary.TotalProceeds {
			currencies = append(currencies, currency)
		}
		sort.Strings(currencies)

		for _, currency := range currencies {
			amount := summary.TotalProceeds[currency]
			if amount > 0 {
				fmt.Printf("  %s%s %.2f%s", colorGreen, currency, amount, colorReset)
				
				// Add original currency if available
				if origAmount, ok := summary.TotalProceeds[currency+"_ORIG"]; ok {
					fmt.Printf(" %s(%.2f local)%s", colorGray, origAmount, colorReset)
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Printf("  %sNo revenue (free apps only)%s\n", colorGray, colorReset)
	}

	// Countries
	fmt.Printf("\n  %s🌍 Countries%s\n", colorBold, colorReset)
	fmt.Printf("  %s%d markets%s\n", colorCyan, summary.TotalCountries, colorReset)
}

// displayAppTable shows apps in a formatted table
func displayAppTable(apps []models.AppSales) {
	// Calculate column widths
	maxAppName := 20
	for _, app := range apps {
		if len(app.AppName) > maxAppName {
			maxAppName = len(app.AppName)
		}
	}

	// Header
	fmt.Printf("\n  %-*s  %8s  %12s  %s\n", 
		maxAppName, "App", "Units", "Revenue", "Top Markets")
	fmt.Printf("  %s  %s  %s  %s\n",
		strings.Repeat("─", maxAppName),
		strings.Repeat("─", 8),
		strings.Repeat("─", 12),
		strings.Repeat("─", 20))

	// Rows
	for _, app := range apps {
		// App name (truncate if needed)
		appName := app.AppName
		if len(appName) > maxAppName {
			appName = appName[:maxAppName-3] + "..."
		}

		// Format revenue
		revenueStr := formatRevenue(app.Summary.TotalProceeds)

		// Top markets
		topMarkets := ""
		for i, country := range app.Summary.TopCountries {
			if i > 2 {
				break
			}
			if i > 0 {
				topMarkets += ", "
			}
			topMarkets += country.Country
		}

		fmt.Printf("  %-*s  %8s  %12s  %s\n",
			maxAppName, appName,
			formatNumber(app.Summary.TotalUnits),
			revenueStr,
			topMarkets)
	}
}

// displayCountryBreakdown shows sales by country
func displayCountryBreakdown(report *models.SalesReport) {
	fmt.Printf("\n%s🌍 Country Breakdown%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 60))

	// Aggregate by country
	countryData := make(map[string]*models.CountrySales)
	
	for _, app := range report.Apps {
		for _, sale := range app.Sales {
			if _, exists := countryData[sale.Country]; !exists {
				countryData[sale.Country] = &models.CountrySales{
					Country:     sale.Country,
					CountryName: sale.CountryName,
					Proceeds:    make(map[string]float64),
				}
			}
			
			countryData[sale.Country].Units += sale.Units
			
			if sale.DeveloperProceeds.Amount > 0 {
				countryData[sale.Country].Proceeds[sale.DeveloperProceeds.Currency] += sale.DeveloperProceeds.Amount
			}
		}
	}

	// Convert to slice and sort
	countries := make([]models.CountrySales, 0, len(countryData))
	for _, country := range countryData {
		countries = append(countries, *country)
	}
	
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Units > countries[j].Units
	})

	// Display top countries
	fmt.Printf("\n  %-20s  %8s  %s\n", "Country", "Units", "Revenue")
	fmt.Printf("  %s  %s  %s\n",
		strings.Repeat("─", 20),
		strings.Repeat("─", 8),
		strings.Repeat("─", 20))

	for i, country := range countries {
		if i >= 10 {
			break
		}
		
		countryName := country.CountryName
		if countryName == "" {
			countryName = country.Country
		}
		if len(countryName) > 20 {
			countryName = countryName[:17] + "..."
		}

		fmt.Printf("  %-20s  %8s  %s\n",
			countryName,
			formatNumber(country.Units),
			formatRevenue(country.Proceeds))
	}

	if len(countries) > 10 {
		fmt.Printf("  %s... and %d more countries%s\n", colorGray, len(countries)-10, colorReset)
	}
}

// displayTrends shows trend analysis
func displayTrends(trends *models.TrendAnalysis) {
	fmt.Printf("\n%s📈 Trends & Insights%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 60))

	// Units trend
	trendIcon := trends.UnitsTrend.GetIcon()
	trendColor := getTrendColor(trends.UnitsTrend)
	
	fmt.Printf("\n  Units: %s%s %.1f%%%s",
		trendColor, trendIcon, trends.UnitsChange, colorReset)
	
	if trends.UnitsTrend == models.TrendUp {
		fmt.Printf(" %s📈 Growth!%s", colorGreen, colorReset)
	} else if trends.UnitsTrend == models.TrendDown {
		fmt.Printf(" %s📉 Decline%s", colorRed, colorReset)
	}
	fmt.Println()

	// Revenue trend
	if trends.ProceedsChange != 0 {
		trendIcon = trends.ProceedsTrend.GetIcon()
		trendColor = getTrendColor(trends.ProceedsTrend)
		
		fmt.Printf("  Revenue: %s%s %.1f%%%s\n",
			trendColor, trendIcon, trends.ProceedsChange, colorReset)
	}

	// New markets
	if len(trends.NewCountries) > 0 {
		fmt.Printf("\n  %s🆕 New Markets:%s %s\n",
			colorGreen, colorReset,
			strings.Join(trends.NewCountries, ", "))
	}

	// Best performer
	if trends.BestPerformer != nil {
		fmt.Printf("\n  %s🏆 Top Performer:%s %s (+%d units)\n",
			colorGreen, colorReset,
			trends.BestPerformer.AppName,
			trends.BestPerformer.Change)
	}
}

// displayComparison shows a comparison between two periods
func displayComparison(cmd *cobra.Command, comp *sales.Comparison) {
	fmt.Printf("\n%s📊 Sales Comparison%s\n", colorBold, colorReset)
	fmt.Printf("%s vs %s\n",
		comp.Previous.Date.Format("January 2006"),
		comp.Current.Date.Format("January 2006"))
	fmt.Println(strings.Repeat("═", 60))

	// Overall metrics
	fmt.Printf("\n%s📈 Overall Performance%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 40))

	// Units comparison
	prevUnits := comp.Previous.Summary.TotalUnits
	currUnits := comp.Current.Summary.TotalUnits
	unitsChange := currUnits - prevUnits
	
	fmt.Printf("  Units: %s → %s ",
		formatNumber(prevUnits),
		formatNumber(currUnits))
	
	if unitsChange > 0 {
		fmt.Printf("%s(+%s, +%.1f%%)%s\n",
			colorGreen, formatNumber(unitsChange), comp.UnitsChange, colorReset)
	} else if unitsChange < 0 {
		fmt.Printf("%s(%s, %.1f%%)%s\n",
			colorRed, formatNumber(unitsChange), comp.UnitsChange, colorReset)
	} else {
		fmt.Printf("%s(no change)%s\n", colorGray, colorReset)
	}

	// Revenue comparison by currency
	if len(comp.ProceedsChange) > 0 {
		fmt.Printf("\n  Revenue:\n")
		for currency, change := range comp.ProceedsChange {
			prevAmount := comp.Previous.Summary.TotalProceeds[currency]
			currAmount := comp.Current.Summary.TotalProceeds[currency]
			
			fmt.Printf("    %s: %.2f → %.2f ",
				currency, prevAmount, currAmount)
			
			if change > 0 {
				fmt.Printf("%s(+%.1f%%)%s\n", colorGreen, change, colorReset)
			} else if change < 0 {
				fmt.Printf("%s(%.1f%%)%s\n", colorRed, change, colorReset)
			} else {
				fmt.Printf("%s(no change)%s\n", colorGray, colorReset)
			}
		}
	}

	// App changes
	if len(comp.TopGainers) > 0 || len(comp.TopLosers) > 0 {
		fmt.Printf("\n%s📱 App Performance%s\n", colorBold, colorReset)
		fmt.Println(strings.Repeat("─", 40))

		if len(comp.TopGainers) > 0 {
			fmt.Printf("\n  %s🚀 Top Gainers:%s\n", colorGreen, colorReset)
			for _, app := range comp.TopGainers {
				fmt.Printf("    • %s: +%s units (+%.1f%%)\n",
					app.AppName,
					formatNumber(app.UnitsChange),
					app.UnitsPercent)
			}
		}

		if len(comp.TopLosers) > 0 {
			fmt.Printf("\n  %s📉 Biggest Declines:%s\n", colorRed, colorReset)
			for _, app := range comp.TopLosers {
				fmt.Printf("    • %s: %s units (%.1f%%)\n",
					app.AppName,
					formatNumber(app.UnitsChange),
					app.UnitsPercent)
			}
		}
	}

	// New and removed apps
	if len(comp.NewApps) > 0 {
		fmt.Printf("\n  %s✨ New Apps:%s %s\n",
			colorGreen, colorReset,
			strings.Join(comp.NewApps, ", "))
	}

	if len(comp.RemovedApps) > 0 {
		fmt.Printf("  %s❌ Removed Apps:%s %s\n",
			colorRed, colorReset,
			strings.Join(comp.RemovedApps, ", "))
	}
}

// displayTrendsReport shows trend analysis
func displayTrendsReport(cmd *cobra.Command, trends *sales.TrendReport) {
	fmt.Printf("\n%s📈 Sales Trends Analysis%s\n", colorBold, colorReset)
	fmt.Printf("Period: %s to %s (%d periods)\n",
		trends.Periods[0].Format("Jan 2006"),
		trends.Periods[len(trends.Periods)-1].Format("Jan 2006"),
		len(trends.Periods))
	fmt.Println(strings.Repeat("═", 60))

	// Show chart if requested
	if mustGetBool(cmd, "chart") {
		displayASCIIChart(trends)
	}

	// Overall trends
	fmt.Printf("\n%s📊 Overall Trends%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 40))

	// Calculate growth
	firstUnits := trends.TotalUnits[0]
	lastUnits := trends.TotalUnits[len(trends.TotalUnits)-1]
	
	if firstUnits > 0 {
		growth := ((float64(lastUnits) - float64(firstUnits)) / float64(firstUnits)) * 100
		fmt.Printf("  Total Growth: ")
		
		if growth > 0 {
			fmt.Printf("%s+%.1f%%%s\n", colorGreen, growth, colorReset)
		} else {
			fmt.Printf("%s%.1f%%%s\n", colorRed, growth, colorReset)
		}
	}

	// Period summary
	fmt.Printf("\n  Period Summary:\n")
	for i, period := range trends.Periods {
		fmt.Printf("    %s: %s units",
			period.Format("Jan 2006"),
			formatNumber(trends.TotalUnits[i]))
		
		// Show revenue for main currency
		if len(trends.TotalProceeds) > 0 {
			// Find primary currency
			var primaryCurrency string
			var maxAmount float64
			for currency, amounts := range trends.TotalProceeds {
				if amounts[i] > maxAmount {
					primaryCurrency = currency
					maxAmount = amounts[i]
				}
			}
			
			if primaryCurrency != "" && maxAmount > 0 {
				fmt.Printf(" (%s %.2f)", primaryCurrency, maxAmount)
			}
		}
		
		fmt.Println()
	}

	// Insights
	if len(trends.Insights) > 0 {
		fmt.Printf("\n%s💡 Insights%s\n", colorBold, colorReset)
		fmt.Println(strings.Repeat("─", 40))

		for _, insight := range trends.Insights {
			icon := getInsightIcon(insight.Type)
			color := getInsightColor(insight.Severity)
			
			fmt.Printf("\n  %s%s %s%s\n", color, icon, insight.Title, colorReset)
			fmt.Printf("  %s\n", insight.Description)
		}
	}
}

// displayASCIIChart shows a simple ASCII chart of trends
func displayASCIIChart(trends *sales.TrendReport) {
	fmt.Printf("\n%s📊 Units Trend Chart%s\n", colorBold, colorReset)
	
	// Find max value for scaling
	maxUnits := 0
	for _, units := range trends.TotalUnits {
		if units > maxUnits {
			maxUnits = units
		}
	}
	
	if maxUnits == 0 {
		return
	}
	
	// Chart height
	height := 10
	
	// Draw chart
	for h := height; h >= 0; h-- {
		threshold := float64(h) / float64(height) * float64(maxUnits)
		
		// Y-axis label
		if h == height {
			fmt.Printf("%6d │", maxUnits)
		} else if h == height/2 {
			fmt.Printf("%6d │", maxUnits/2)
		} else if h == 0 {
			fmt.Printf("%6d │", 0)
		} else {
			fmt.Printf("       │")
		}
		
		// Bars
		for _, units := range trends.TotalUnits {
			if float64(units) >= threshold {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
	
	// X-axis
	fmt.Print("       └")
	for range trends.Periods {
		fmt.Print("──")
	}
	fmt.Println()
	
	// Period labels
	fmt.Print("        ")
	for _, period := range trends.Periods {
		fmt.Printf("%-2s", period.Format("1"))
	}
	fmt.Println()
}

// displayDetailedReport shows a detailed report with all information
func displayDetailedReport(report *models.SalesReport) {
	// This would show more detailed information including:
	// - Individual transactions
	// - Device breakdowns
	// - Platform analysis
	// - Detailed country metrics
	// etc.
	
	displayMonthlyReport(nil, report)
}

// Helper functions

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	
	// Add thousands separators
	str := fmt.Sprintf("%d", n)
	result := ""
	
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	
	return result
}

func formatRevenue(amounts map[string]float64) string {
	if len(amounts) == 0 {
		return "-"
	}
	
	// Sort currencies
	currencies := make([]string, 0, len(amounts))
	for currency := range amounts {
		if !strings.HasSuffix(currency, "_ORIG") {
			currencies = append(currencies, currency)
		}
	}
	sort.Strings(currencies)
	
	// Format each currency
	parts := make([]string, 0, len(currencies))
	for _, currency := range currencies {
		if amount := amounts[currency]; amount > 0 {
			parts = append(parts, fmt.Sprintf("%s %.2f", currency, amount))
		}
	}
	
	if len(parts) == 0 {
		return "-"
	}
	
	return strings.Join(parts, ", ")
}

func getTrendColor(trend models.TrendDirection) string {
	switch trend {
	case models.TrendUp:
		return colorGreen
	case models.TrendDown:
		return colorRed
	default:
		return colorYellow
	}
}

func getInsightIcon(insightType sales.InsightType) string {
	switch insightType {
	case sales.InsightGrowth:
		return "📈"
	case sales.InsightDecline:
		return "📉"
	case sales.InsightAnomaly:
		return "⚠️"
	case sales.InsightOpportunity:
		return "💡"
	case sales.InsightWarning:
		return "⚠️"
	default:
		return "ℹ️"
	}
}

func getInsightColor(severity sales.InsightSeverity) string {
	switch severity {
	case sales.InsightSeveritySuccess:
		return colorGreen
	case sales.InsightSeverityWarning:
		return colorYellow
	case sales.InsightSeverityCritical:
		return colorRed
	default:
		return colorBlue
	}
}