package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/marcusziade/pomme/internal/client"
	"github.com/marcusziade/pomme/internal/config"
	"github.com/marcusziade/pomme/internal/services/analytics"
	"github.com/spf13/cobra"
)

var analyticsCmd = &cobra.Command{
	Use:     "analytics",
	Short:   "View app performance analytics and metrics",
	Long:    "Access detailed performance metrics, launch times, memory usage, and battery consumption data for your apps",
	Aliases: []string{"metrics", "perf"},
}

var analyticsShowCmd = &cobra.Command{
	Use:     "show <app-id>",
	Short:   "Show performance metrics for an app",
	Long:    "Display detailed performance and power metrics including launch times, memory usage, and battery consumption",
	Args:    cobra.ExactArgs(1),
	RunE:    runAnalyticsShow,
}

var analyticsCompareCmd = &cobra.Command{
	Use:     "compare <app-id> --version1 <v1> --version2 <v2>",
	Short:   "Compare metrics between app versions",
	Long:    "Compare performance metrics between two versions of your app to identify improvements or regressions",
	Args:    cobra.ExactArgs(1),
	RunE:    runAnalyticsCompare,
}

var analyticsTrendsCmd = &cobra.Command{
	Use:     "trends <app-id>",
	Short:   "Show performance trends over time",
	Long:    "Display performance trends across multiple app versions to track improvements or regressions",
	Args:    cobra.ExactArgs(1),
	RunE:    runAnalyticsTrends,
}

var analyticsUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Get app usage data",
	Long:  "Retrieves app usage metrics like active devices, sessions, and crash data",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📊 App usage analytics coming soon!")
		return nil
	},
}

var analyticsSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Get subscription metrics",
	Long:  "Retrieves subscription metrics like conversions, renewals, and cancellations",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("💳 Subscription analytics coming soon!")
		return nil
	},
}

var (
	// Flags
	analyticsVersion1  string
	analyticsVersion2  string
	analyticsDevice    string
	analyticsMetric    string
	analyticsShowGoals bool
)

func init() {
	// Add subcommands
	analyticsCmd.AddCommand(analyticsShowCmd)
	analyticsCmd.AddCommand(analyticsCompareCmd)
	analyticsCmd.AddCommand(analyticsTrendsCmd)
	analyticsCmd.AddCommand(analyticsUsageCmd)
	analyticsCmd.AddCommand(analyticsSubscriptionsCmd)
	
	// Add flags for compare
	analyticsCompareCmd.Flags().StringVar(&analyticsVersion1, "version1", "", "First version to compare")
	analyticsCompareCmd.Flags().StringVar(&analyticsVersion2, "version2", "", "Second version to compare")
	analyticsCompareCmd.MarkFlagRequired("version1")
	analyticsCompareCmd.MarkFlagRequired("version2")
	
	// Common flags for show
	analyticsShowCmd.Flags().StringVar(&analyticsDevice, "device", "", "Filter by device type (e.g., iPhone, iPad)")
	analyticsShowCmd.Flags().StringVar(&analyticsMetric, "metric", "", "Show specific metric (launch, memory, battery)")
	analyticsShowCmd.Flags().BoolVar(&analyticsShowGoals, "goals", false, "Show performance goals")
	
	// Flags for trends
	analyticsTrendsCmd.Flags().StringVar(&analyticsDevice, "device", "", "Filter by device type")
	analyticsTrendsCmd.Flags().StringVar(&analyticsMetric, "metric", "", "Show specific metric trend")
	
	// Add flags for analytics usage command
	analyticsUsageCmd.Flags().String("app", "", "App ID (required)")
	analyticsUsageCmd.Flags().String("metric", "activeDevices", "Metric to retrieve (activeDevices, sessions, crashes, etc.)")
	analyticsUsageCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	analyticsUsageCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	analyticsUsageCmd.Flags().String("frequency", "day", "Data frequency (day, week, month)")
	analyticsUsageCmd.MarkFlagRequired("app")
	
	// Add flags for analytics subscriptions command
	analyticsSubscriptionsCmd.Flags().String("app", "", "App ID (required)")
	analyticsSubscriptionsCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	analyticsSubscriptionsCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	analyticsSubscriptionsCmd.Flags().String("subscription-group", "", "Subscription group ID")
	analyticsSubscriptionsCmd.MarkFlagRequired("app")
}

func runAnalyticsShow(cmd *cobra.Command, args []string) error {
	appID := args[0]
	
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create client
	apiClient, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Create analytics service
	svc := analytics.NewService(apiClient)
	
	fmt.Printf("📊 Fetching performance metrics for app %s...\n\n", appID)
	
	// Fetch metrics
	ctx := context.Background()
	report, err := svc.GetAppMetrics(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to fetch metrics: %w", err)
	}

	// Display the report
	displayAnalyticsReport(report)
	
	return nil
}

func runAnalyticsCompare(cmd *cobra.Command, args []string) error {
	appID := args[0]
	
	fmt.Printf("📊 Comparing versions %s and %s for app %s...\n\n", 
		analyticsVersion1, analyticsVersion2, appID)
	
	// This would be implemented when the API supports version comparison
	fmt.Println("Version comparison feature coming soon!")
	
	return nil
}

func runAnalyticsTrends(cmd *cobra.Command, args []string) error {
	appID := args[0]
	
	fmt.Printf("📊 Analyzing performance trends for app %s...\n\n", appID)
	
	// This would show trends over multiple versions
	fmt.Println("Trends analysis feature coming soon!")
	
	return nil
}

// displayAnalyticsReport displays the analytics report in a formatted way
func displayAnalyticsReport(report *analytics.AnalyticsReport) {
	fmt.Printf("%s📊 Performance Analytics Report%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 60))
	
	// Display each metric category
	for _, category := range report.Categories {
		fmt.Printf("\n%s%s%s\n", colorBold, category.DisplayName, colorReset)
		fmt.Println(strings.Repeat("─", 40))
		
		// Create a table for metrics
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "  Metric\tP50\tP90\tP95\tGoal\n")
		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\t%s\n", 
			strings.Repeat("─", 20),
			strings.Repeat("─", 10),
			strings.Repeat("─", 10),
			strings.Repeat("─", 10),
			strings.Repeat("─", 10))
		
		for _, metric := range category.Metrics {
			goalStr := "-"
			if metric.Goal != nil {
				goalStr = fmt.Sprintf("<%g%s", metric.Goal.Value, metric.Unit)
			}
			
			fmt.Fprintf(w, "  %s\t%g%s\t%g%s\t%g%s\t%s\n",
				metric.DisplayName,
				metric.Percentiles.P50.Value, metric.Unit,
				metric.Percentiles.P90.Value, metric.Unit,
				metric.Percentiles.P95.Value, metric.Unit,
				goalStr,
			)
		}
		w.Flush()
	}
	
	// Display insights if any
	if len(report.Insights) > 0 {
		fmt.Printf("\n%s💡 Insights%s\n", colorBold, colorReset)
		fmt.Println(strings.Repeat("─", 40))
		
		for _, insight := range report.Insights {
			icon := "📈"
			color := colorGreen
			if insight.Attributes.Direction == "REGRESSION" {
				icon = "📉"
				color = colorRed
			}
			
			fmt.Printf("  %s %s%s: %.1f%% change%s\n",
				icon,
				color,
				insight.Attributes.InsightType,
				insight.Attributes.Change,
				colorReset,
			)
		}
	}
	
	fmt.Printf("\n%sGenerated: %s%s\n", colorGray, time.Now().Format("2006-01-02 15:04:05"), colorReset)
}