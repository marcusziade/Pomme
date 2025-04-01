package commands

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/marcus/pomme/internal/config"
	"github.com/marcus/pomme/internal/models"
	"github.com/marcus/pomme/internal/output"
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

var salesMonthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Show sales for the past month",
	Long:  `Shows a summary of sales for the past month, with a focus on apps that generated revenue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		// Calculate last month
		now := time.Now()
		lastMonth := now.AddDate(0, -1, 0)
		monthStr := lastMonth.Format("2006-01")
		
		// Use default vendor number if available
		vendorNumber := cfg.Defaults.VendorNumber
		if vendorNumber == "" {
			return fmt.Errorf("vendor number is required in config file")
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
		
		fmt.Printf("Fetching monthly sales report for %s...\n", monthStr)
		
		reportData, err := client.GetSalesReport(ctx, models.ReportFrequencyMonthly, monthStr, models.ReportTypeSales, vendorNumber)
		if err != nil {
			return fmt.Errorf("failed to get sales report: %w", err)
		}
		
		// Process report data
		if reportData != nil && len(reportData) > 0 {
			// Parse the data as TSV
			reader := csv.NewReader(bytes.NewReader(reportData))
			reader.Comma = '\t' // Use tab as delimiter
			
			// Read header row
			header, err := reader.Read()
			if err != nil {
				return fmt.Errorf("failed to read CSV header: %w", err)
			}
			
			// Process records
			var records []models.SalesRecord
			
			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("failed to read CSV row: %w", err)
				}
				
				// Create a map to easily access fields by header name
				fields := make(map[string]string)
				for i, h := range header {
					if i < len(row) {
						fields[h] = row[i]
					}
				}
				
				// Parse the record
				var record models.SalesRecord
				record.Provider = fields["Provider"]
				record.ProviderCountry = fields["Provider Country"]
				record.SKU = fields["SKU"]
				record.Developer = fields["Developer"] 
				record.Title = fields["Title"]
				record.Version = fields["Version"]
				record.ProductTypeID = fields["Product Type Identifier"]
				record.AppleID = fields["Apple Identifier"]
				record.CountryCode = fields["Country Code"]
				record.CustomerCurrency = fields["Customer Currency"]
				record.CurrencyOfProceeds = fields["Currency of Proceeds"]
				record.BeginDate = fields["Begin Date"]
				record.EndDate = fields["End Date"]
				
				// Parse numeric fields
				units, _ := strconv.Atoi(fields["Units"])
				record.Units = units
				
				customerPrice, _ := strconv.ParseFloat(fields["Customer Price"], 64)
				record.CustomerPrice = customerPrice
				
				proceeds, _ := strconv.ParseFloat(fields["Developer Proceeds"], 64)
				record.DeveloperProceeds = proceeds
				
				// Add additional fields
				record.PromoCode = fields["Promo Code"]
				record.ParentID = fields["Parent Identifier"]
				record.Subscription = fields["Subscription"]
				record.Period = fields["Period"]
				record.Category = fields["Category"]
				record.CMB = fields["CMB"]
				record.DeviceType = fields["Device"]
				record.SupportedPlatforms = fields["Supported Platforms"]
				record.ProceedsReason = fields["Proceeds Reason"]
				record.PreservedPricing = fields["Preserved Pricing"]
				record.Client = fields["Client"]
				record.OrderType = fields["Order Type"]
				
				// Keep records regardless of proceeds
				records = append(records, record)
			}
			
			// Create summary by aggregating data
			summary := createSalesSummary(records)
			
			// Filter to only show apps with meaningful data
			var meaningfulSummary []SalesSummary
			for _, s := range summary {
				// Include if there are units or proceeds
				if s.Units > 0 {
					meaningfulSummary = append(meaningfulSummary, s)
				}
			}
			
			if len(meaningfulSummary) > 0 {
				// Sort by proceeds (highest first) before display
				formatter := output.NewFormatter(output.FormatTable, os.Stdout)
				if err := formatter.Format(meaningfulSummary); err != nil {
					return fmt.Errorf("failed to format sales summary: %w", err)
				}
				
				// Print a summary message
				totalProceeds := 0.0
				totalUnits := 0
				for _, s := range meaningfulSummary {
					totalProceeds += s.TotalProceeds
					totalUnits += s.Units
				}
				
				fmt.Printf("\nTotal: %d units, %.2f %s\n", totalUnits, totalProceeds, meaningfulSummary[0].Currency)
			} else {
				fmt.Println("No sales or downloads found for the past month.")
			}
		} else {
			fmt.Println("No report data available for the specified period.")
		}
		
		return nil
	},
}

func init() {
	salesCmd.AddCommand(salesReportCmd)
	salesCmd.AddCommand(salesTrendsCmd)
	salesCmd.AddCommand(salesMonthlyCmd)
	
	// Add flags for sales report command
	salesReportCmd.Flags().String("period", "DAILY", "Report period (DAILY, WEEKLY, MONTHLY, YEARLY)")
	salesReportCmd.Flags().String("date", "latest", "Report date (YYYY-MM-DD or 'latest')")
	salesReportCmd.Flags().String("vendor", "", "Vendor number (default: from config)")
	salesReportCmd.Flags().String("type", "SALES", "Report type (SALES, SUBSCRIPTION, SUBSCRIPTION_EVENT)")
	salesReportCmd.Flags().Bool("summary", false, "Show summary only")
}