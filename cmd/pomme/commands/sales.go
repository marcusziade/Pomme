package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/marcusziade/pomme/internal/config"
	"github.com/marcusziade/pomme/internal/models"
	"github.com/marcusziade/pomme/internal/output"
	"github.com/marcusziade/pomme/internal/services/cache"
	"github.com/marcusziade/pomme/internal/services/sales"
	"github.com/marcusziade/pomme/pkg/pomme"
	"github.com/spf13/cobra"
)

var salesCmd = &cobra.Command{
	Use:   "sales",
	Short: "Sales and financial reports",
	Long: `Powerful sales reporting and analytics for your App Store apps.

Examples:
  # Show sales for the latest available month
  pomme sales

  # Show sales for a specific month
  pomme sales --month 2025-03

  # Compare two months
  pomme sales compare --current 2025-03 --previous 2025-02

  # Show trends over the last 6 months
  pomme sales trends --months 6

  # Export sales data
  pomme sales export --month 2025-03 --format csv`,
}

// Main sales command - shows latest monthly report by default
var salesReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Get sales reports (default: latest monthly)",
	Long: `Fetches and displays sales reports with automatic currency handling and insights.

The command automatically handles Apple's 5-day report availability delay.`,
	RunE:  runSalesReport,
}

// Quick monthly overview
var salesMonthlyCmd = &cobra.Command{
	Use:     "monthly [YYYY-MM]",
	Aliases: []string{"month", "m"},
	Short:   "Quick monthly sales overview",
	Example: `  pomme sales monthly           # Latest available month
  pomme sales monthly 2025-03   # Specific month`,
	RunE: runMonthlyReport,
}

// Compare periods
var salesCompareCmd = &cobra.Command{
	Use:     "compare",
	Aliases: []string{"comp", "vs"},
	Short:   "Compare sales between two periods",
	Example: `  pomme sales compare --months 2   # Compare last 2 months
  pomme sales compare --current 2025-03 --previous 2025-02`,
	RunE: runCompare,
}

// Trend analysis
var salesTrendsCmd = &cobra.Command{
	Use:     "trends",
	Aliases: []string{"trend", "t"},
	Short:   "Analyze sales trends over time",
	Example: `  pomme sales trends --months 6    # Last 6 months
  pomme sales trends --quarters 4  # Last 4 quarters
  pomme sales trends --years 2     # Last 2 years`,
	RunE: runTrends,
}

// Export data
var salesExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export sales data in various formats",
	Example: `  pomme sales export --month 2025-03 --format csv
  pomme sales export --last 3 --format excel
  pomme sales export --year 2025 --format json`,
	RunE: runExport,
}

// Watch for updates
var salesWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for new sales data",
	Long:  "Continuously monitors for new sales data and notifies when available.",
	RunE:  runWatch,
}

func init() {
	salesCmd.AddCommand(salesReportCmd)
	salesCmd.AddCommand(salesMonthlyCmd)
	salesCmd.AddCommand(salesCompareCmd)
	salesCmd.AddCommand(salesTrendsCmd)
	salesCmd.AddCommand(salesExportCmd)
	salesCmd.AddCommand(salesWatchCmd)

	// Set monthly as the default when no subcommand is specified
	salesCmd.RunE = runMonthlyReport

	// Global flags
	salesCmd.PersistentFlags().String("vendor", "", "Vendor number (default: from config)")
	salesCmd.PersistentFlags().Bool("no-cache", false, "Skip cache and fetch fresh data")
	salesCmd.PersistentFlags().Bool("json", false, "Output raw JSON")

	// Report command flags
	salesReportCmd.Flags().String("period", "MONTHLY", "Report period (DAILY, WEEKLY, MONTHLY, YEARLY)")
	salesReportCmd.Flags().String("date", "latest", "Report date (YYYY-MM-DD or 'latest')")
	salesReportCmd.Flags().String("type", "SALES", "Report type (SALES, SUBSCRIPTION, SUBSCRIPTION_EVENT)")

	// Monthly command flags
	salesMonthlyCmd.Flags().Bool("details", false, "Show detailed breakdown")
	salesMonthlyCmd.Flags().Bool("by-country", false, "Group by country")
	salesMonthlyCmd.Flags().Bool("by-app", false, "Group by app")

	// Compare command flags
	salesCompareCmd.Flags().String("current", "", "Current period (YYYY-MM)")
	salesCompareCmd.Flags().String("previous", "", "Previous period (YYYY-MM)")
	salesCompareCmd.Flags().Int("months", 0, "Compare last N months")
	salesCompareCmd.Flags().Bool("percentage", false, "Show percentage changes")

	// Trends command flags
	salesTrendsCmd.Flags().Int("months", 0, "Analyze last N months")
	salesTrendsCmd.Flags().Int("quarters", 0, "Analyze last N quarters")
	salesTrendsCmd.Flags().Int("years", 0, "Analyze last N years")
	salesTrendsCmd.Flags().String("group", "total", "Group by (total, app, country, platform)")
	salesTrendsCmd.Flags().Bool("chart", false, "Display ASCII chart")

	// Export command flags
	salesExportCmd.Flags().String("month", "", "Specific month (YYYY-MM)")
	salesExportCmd.Flags().Int("last", 0, "Last N months")
	salesExportCmd.Flags().String("year", "", "Full year (YYYY)")
	salesExportCmd.Flags().String("format", "csv", "Export format (csv, json, excel)")
	salesExportCmd.Flags().String("output", "", "Output file (default: stdout)")
	salesExportCmd.Flags().Bool("detailed", false, "Include all transaction details")

	// Watch command flags
	salesWatchCmd.Flags().Duration("interval", 1*time.Hour, "Check interval")
	salesWatchCmd.Flags().Bool("notify", false, "Send desktop notification")
}

func runMonthlyReport(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse month from args or flags
	var targetMonth time.Time
	if len(args) > 0 {
		parsed, err := time.Parse("2006-01", args[0])
		if err != nil {
			return fmt.Errorf("invalid month format, use YYYY-MM: %w", err)
		}
		targetMonth = parsed
	} else {
		// Calculate the latest available month
		targetMonth = calculateLatestAvailableMonth()
	}

	// Get configuration
	cfg, service, err := setupSalesService(cmd)
	if err != nil {
		return err
	}

	// Create report options
	options := sales.ReportOptions{
		Period:       models.ReportFrequencyMonthly,
		Date:         targetMonth,
		ReportType:   models.ReportTypeSales,
		VendorNumber: cfg.Defaults.VendorNumber,
		NoCache:      mustGetBool(cmd, "no-cache"),
		IncludeAnalysis: true,
	}

	// Show what we're fetching
	fmt.Printf("📊 Fetching sales report for %s...\n", targetMonth.Format("January 2006"))

	// Fetch the report
	report, err := service.GetReport(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to fetch report: %w", err)
	}

	if report == nil || len(report.Apps) == 0 {
		fmt.Printf("\n❌ No sales data available for %s\n", targetMonth.Format("January 2006"))
		
		// Suggest checking other months
		fmt.Println("\n💡 Try checking recent months:")
		for i := 1; i <= 3; i++ {
			suggestedMonth := targetMonth.AddDate(0, -i, 0)
			fmt.Printf("   pomme sales monthly %s\n", suggestedMonth.Format("2006-01"))
		}
		return nil
	}

	// Display the report
	if mustGetBool(cmd, "json") {
		return output.JSON(report)
	}

	displayMonthlyReport(cmd, report)
	return nil
}

func runSalesReport(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get configuration
	cfg, service, err := setupSalesService(cmd)
	if err != nil {
		return err
	}

	// Parse report options
	period, err := parseReportPeriod(mustGetString(cmd, "period"))
	if err != nil {
		return err
	}

	date, err := parseReportDate(mustGetString(cmd, "date"), period)
	if err != nil {
		return err
	}

	reportType, err := parseReportType(mustGetString(cmd, "type"))
	if err != nil {
		return err
	}

	// Create report options
	options := sales.ReportOptions{
		Period:       period,
		Date:         date,
		ReportType:   reportType,
		VendorNumber: getVendorNumber(cmd, cfg),
		NoCache:      mustGetBool(cmd, "no-cache"),
	}

	// Fetch the report
	report, err := service.GetReport(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to fetch report: %w", err)
	}

	// Display the report
	if mustGetBool(cmd, "json") {
		return output.JSON(report)
	}

	displayDetailedReport(report)
	return nil
}

func runCompare(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	cfg, service, err := setupSalesService(cmd)
	if err != nil {
		return err
	}

	// Determine periods to compare
	var currentOpt, previousOpt sales.ReportOptions

	if months := mustGetInt(cmd, "months"); months > 0 {
		// Compare last N months
		current := calculateLatestAvailableMonth()
		previous := current.AddDate(0, -months, 0)

		currentOpt = sales.ReportOptions{
			Period:       models.ReportFrequencyMonthly,
			Date:         current,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
		}

		previousOpt = sales.ReportOptions{
			Period:       models.ReportFrequencyMonthly,
			Date:         previous,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
		}
	} else {
		// Use specific months
		currentStr := mustGetString(cmd, "current")
		previousStr := mustGetString(cmd, "previous")

		if currentStr == "" || previousStr == "" {
			return fmt.Errorf("specify either --months or both --current and --previous")
		}

		current, err := time.Parse("2006-01", currentStr)
		if err != nil {
			return fmt.Errorf("invalid current month: %w", err)
		}

		previous, err := time.Parse("2006-01", previousStr)
		if err != nil {
			return fmt.Errorf("invalid previous month: %w", err)
		}

		currentOpt = sales.ReportOptions{
			Period:       models.ReportFrequencyMonthly,
			Date:         current,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
		}

		previousOpt = sales.ReportOptions{
			Period:       models.ReportFrequencyMonthly,
			Date:         previous,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
		}
	}

	fmt.Printf("📊 Comparing %s vs %s...\n", 
		currentOpt.Date.Format("January 2006"),
		previousOpt.Date.Format("January 2006"))

	// Fetch comparison
	comparison, err := service.GetComparison(ctx, currentOpt, previousOpt)
	if err != nil {
		return fmt.Errorf("failed to get comparison: %w", err)
	}

	// Display comparison
	if mustGetBool(cmd, "json") {
		return output.JSON(comparison)
	}

	displayComparison(cmd, comparison)
	return nil
}

func runTrends(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	cfg, service, err := setupSalesService(cmd)
	if err != nil {
		return err
	}

	// Determine trend period
	var trendOpt sales.TrendOptions

	if months := mustGetInt(cmd, "months"); months > 0 {
		trendOpt = sales.TrendOptions{
			Frequency:    models.ReportFrequencyMonthly,
			EndDate:      calculateLatestAvailableMonth(),
			Periods:      months,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
			GroupBy:      mustGetString(cmd, "group"),
		}
	} else if quarters := mustGetInt(cmd, "quarters"); quarters > 0 {
		// For quarters, we'll use monthly reports and aggregate
		trendOpt = sales.TrendOptions{
			Frequency:    models.ReportFrequencyMonthly,
			EndDate:      calculateLatestAvailableMonth(),
			Periods:      quarters * 3,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
			GroupBy:      mustGetString(cmd, "group"),
		}
	} else if years := mustGetInt(cmd, "years"); years > 0 {
		trendOpt = sales.TrendOptions{
			Frequency:    models.ReportFrequencyYearly,
			EndDate:      time.Now().AddDate(-1, 0, 0), // Last complete year
			Periods:      years,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
			GroupBy:      mustGetString(cmd, "group"),
		}
	} else {
		// Default to last 6 months
		trendOpt = sales.TrendOptions{
			Frequency:    models.ReportFrequencyMonthly,
			EndDate:      calculateLatestAvailableMonth(),
			Periods:      6,
			ReportType:   models.ReportTypeSales,
			VendorNumber: getVendorNumber(cmd, cfg),
			GroupBy:      mustGetString(cmd, "group"),
		}
	}

	fmt.Printf("📈 Analyzing trends over %d %s...\n", trendOpt.Periods, trendOpt.Frequency)

	// Fetch trends
	trends, err := service.GetTrends(ctx, trendOpt)
	if err != nil {
		return fmt.Errorf("failed to analyze trends: %w", err)
	}

	// Display trends
	if mustGetBool(cmd, "json") {
		return output.JSON(trends)
	}

	displayTrendsReport(cmd, trends)
	return nil
}

func runExport(cmd *cobra.Command, args []string) error {
	// TODO: Implement export functionality
	fmt.Println("Export functionality coming soon!")
	return nil
}

func runWatch(cmd *cobra.Command, args []string) error {
	// TODO: Implement watch functionality
	fmt.Println("Watch functionality coming soon!")
	return nil
}

// Helper functions

func setupSalesService(cmd *cobra.Command) (*config.Config, *salesService, error) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Validate config
	if cfg.Auth.KeyID == "" || cfg.Auth.IssuerID == "" || cfg.Auth.PrivateKeyPath == "" {
		return nil, nil, fmt.Errorf("authentication not configured. Run 'pomme config init' first")
	}

	// Read private key
	privateKeyData, err := os.ReadFile(cfg.Auth.PrivateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Create client
	client := pomme.NewClient(cfg.Auth.KeyID, cfg.Auth.IssuerID, string(privateKeyData))

	// Create cache
	cacheService := cache.NewMemoryCache()

	// Create sales service
	// For now, we'll use the client directly until we refactor the service layer
	service := &salesService{
		client: client,
		cache:  cacheService,
	}

	return cfg, service, nil
}

func calculateLatestAvailableMonth() time.Time {
	now := time.Now()
	
	// If we're within the first 5 days of the month, go back 2 months
	if now.Day() <= 5 {
		return now.AddDate(0, -2, 0)
	}
	
	// Otherwise, last month's data should be available
	return now.AddDate(0, -1, 0)
}

func parseReportPeriod(period string) (models.ReportFrequency, error) {
	switch strings.ToUpper(period) {
	case "DAILY", "D":
		return models.ReportFrequencyDaily, nil
	case "WEEKLY", "W":
		return models.ReportFrequencyWeekly, nil
	case "MONTHLY", "M":
		return models.ReportFrequencyMonthly, nil
	case "YEARLY", "Y":
		return models.ReportFrequencyYearly, nil
	default:
		return "", fmt.Errorf("invalid period: %s", period)
	}
}

func parseReportDate(dateStr string, period models.ReportFrequency) (time.Time, error) {
	if dateStr == "latest" {
		switch period {
		case models.ReportFrequencyMonthly:
			return calculateLatestAvailableMonth(), nil
		case models.ReportFrequencyDaily:
			return time.Now().AddDate(0, 0, -1), nil
		case models.ReportFrequencyWeekly:
			// Find last Saturday
			now := time.Now()
			daysSinceSaturday := (int(now.Weekday()) + 1) % 7
			return now.AddDate(0, 0, -daysSinceSaturday), nil
		case models.ReportFrequencyYearly:
			return time.Now().AddDate(-1, 0, 0), nil
		}
	}

	// Parse based on period
	switch period {
	case models.ReportFrequencyDaily, models.ReportFrequencyWeekly:
		return time.Parse("2006-01-02", dateStr)
	case models.ReportFrequencyMonthly:
		return time.Parse("2006-01", dateStr)
	case models.ReportFrequencyYearly:
		return time.Parse("2006", dateStr)
	default:
		return time.Parse("2006-01-02", dateStr)
	}
}

func parseReportType(typeStr string) (models.ReportType, error) {
	switch strings.ToUpper(typeStr) {
	case "SALES", "S":
		return models.ReportTypeSales, nil
	case "SUBSCRIPTION", "SUB":
		return models.ReportTypeSubscription, nil
	case "SUBSCRIPTION_EVENT", "SUB_EVENT":
		return models.ReportTypeSubscriptionEvent, nil
	default:
		return "", fmt.Errorf("invalid report type: %s", typeStr)
	}
}

func getVendorNumber(cmd *cobra.Command, cfg *config.Config) string {
	if vendor := mustGetString(cmd, "vendor"); vendor != "" {
		return vendor
	}
	return cfg.Defaults.VendorNumber
}

// Flag helper functions
func mustGetString(cmd *cobra.Command, flag string) string {
	val, _ := cmd.Flags().GetString(flag)
	return val
}

func mustGetBool(cmd *cobra.Command, flag string) bool {
	val, _ := cmd.Flags().GetBool(flag)
	return val
}

func mustGetInt(cmd *cobra.Command, flag string) int {
	val, _ := cmd.Flags().GetInt(flag)
	return val
}