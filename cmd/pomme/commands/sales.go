package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/marcus/pomme/internal/config"
	"github.com/marcus/pomme/internal/models"
	"github.com/marcus/pomme/pkg/pomme"
	"github.com/spf13/cobra"
)

var salesCmd = &cobra.Command{
	Use:   "sales",
	Short: "Sales and financial reports",
	Long:  `Commands for retrieving and analyzing sales and financial reports.`,
}

var salesReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Get sales reports",
	Long:  `Retrieves sales reports for your apps.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get report parameters
		periodStr, _ := cmd.Flags().GetString("period")
		dateStr, _ := cmd.Flags().GetString("date")
		vendorNumber, _ := cmd.Flags().GetString("vendor")
		typeStr, _ := cmd.Flags().GetString("type")
		_, _ = cmd.Flags().GetBool("summary") // We'll use this later when implementing summary display
		
		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		// Use default vendor number if not provided
		if vendorNumber == "" {
			vendorNumber = cfg.Defaults.VendorNumber
			if vendorNumber == "" {
				return fmt.Errorf("vendor number is required, please provide with --vendor or set in config")
			}
		}
		
		// Convert period string to report frequency
		var period models.ReportFrequency
		switch strings.ToUpper(periodStr) {
		case "DAILY":
			period = models.ReportFrequencyDaily
		case "WEEKLY":
			period = models.ReportFrequencyWeekly
		case "MONTHLY":
			period = models.ReportFrequencyMonthly
		case "YEARLY":
			period = models.ReportFrequencyYearly
		default:
			return fmt.Errorf("invalid period: %s (must be DAILY, WEEKLY, MONTHLY, or YEARLY)", periodStr)
		}
		
		// Handle "latest" date
		if dateStr == "latest" {
			// Use yesterday for daily reports, last week for weekly, etc.
			now := time.Now().AddDate(0, 0, -1)
			switch period {
			case models.ReportFrequencyDaily:
				dateStr = now.Format("2006-01-02")
			case models.ReportFrequencyWeekly:
				// Find the previous Saturday (Apple's week end)
				daysSinceSaturday := (int(now.Weekday()) + 1) % 7
				lastSaturday := now.AddDate(0, 0, -daysSinceSaturday)
				dateStr = lastSaturday.Format("2006-01-02")
			case models.ReportFrequencyMonthly:
				// Previous month
				lastMonth := now.AddDate(0, -1, 0)
				dateStr = lastMonth.Format("2006-01")
			case models.ReportFrequencyYearly:
				// Previous year
				lastYear := now.AddDate(-1, 0, 0)
				dateStr = lastYear.Format("2006")
			}
		}
		
		// Convert report type string to ReportType
		var reportType models.ReportType
		switch strings.ToUpper(typeStr) {
		case "SALES":
			reportType = models.ReportTypeSales
		case "SUBSCRIPTION":
			reportType = models.ReportTypeSubscription
		case "SUBSCRIPTION_EVENT":
			reportType = models.ReportTypeSubscriptionEvent
		default:
			return fmt.Errorf("invalid report type: %s (must be SALES, SUBSCRIPTION, or SUBSCRIPTION_EVENT)", typeStr)
		}
		
		// Validate config
		if cfg.Auth.KeyID == "" {
			return fmt.Errorf("key ID not set in config")
		}
		if cfg.Auth.IssuerID == "" {
			return fmt.Errorf("issuer ID not set in config")
		}
		if cfg.Auth.PrivateKeyPath == "" {
			return fmt.Errorf("private key path not set in config")
		}
		
		// Read private key
		privateKeyData, err := os.ReadFile(cfg.Auth.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read private key: %w", err)
		}
		
		// Create client
		client := pomme.NewClient(
			cfg.Auth.KeyID,
			cfg.Auth.IssuerID,
			string(privateKeyData),
		)
		
		// Fetch sales report
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		
		fmt.Printf("Fetching %s %s report for %s...\n", 
			strings.ToLower(string(period)), 
			strings.ToLower(string(reportType)), 
			dateStr)
		
		reportData, err := client.GetSalesReport(ctx, period, dateStr, reportType, vendorNumber)
		if err != nil {
			return fmt.Errorf("failed to get sales report: %w", err)
		}
		
		// Write report data to stdout or file
		if reportData != nil {
			fmt.Println("Report successfully retrieved!")
			fmt.Printf("Report size: %d bytes\n", len(reportData))
			
			// TODO: Parse and display report data based on format/summary flags
		} else {
			fmt.Println("No report data available for the specified period.")
		}
		
		return nil
	},
}

var salesTrendsCmd = &cobra.Command{
	Use:   "trends",
	Short: "Analyze sales trends",
	Long:  `Analyzes trends in your sales data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement sales trends analysis
		fmt.Println("Sales trends analysis will be implemented in a future version.")
		return nil
	},
}

func init() {
	salesCmd.AddCommand(salesReportCmd)
	salesCmd.AddCommand(salesTrendsCmd)
	
	// Add flags for sales report command
	salesReportCmd.Flags().String("period", "DAILY", "Report period (DAILY, WEEKLY, MONTHLY, YEARLY)")
	salesReportCmd.Flags().String("date", "latest", "Report date (YYYY-MM-DD or 'latest')")
	salesReportCmd.Flags().String("vendor", "", "Vendor number (default: from config)")
	salesReportCmd.Flags().String("type", "SALES", "Report type (SALES, SUBSCRIPTION, SUBSCRIPTION_EVENT)")
	salesReportCmd.Flags().Bool("summary", false, "Show summary only")
}