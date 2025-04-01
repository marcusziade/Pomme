package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/marcusziade/pomme/internal/config"
	"github.com/marcusziade/pomme/internal/output"
	"github.com/marcusziade/pomme/pkg/pomme"
	"github.com/spf13/cobra"
)

var reviewsCmd = &cobra.Command{
	Use:   "reviews",
	Short: "User reviews management",
	Long:  `Commands for retrieving and analyzing user reviews.`,
}

var reviewsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List reviews",
	Long:  `Lists reviews for your apps.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get app ID
		appID, _ := cmd.Flags().GetString("app")
		if appID == "" {
			return fmt.Errorf("app ID is required, please provide with --app")
		}
		
		// Get limit
		limit, _ := cmd.Flags().GetInt("limit")
		
		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
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
		
		// Get output format
		outputFormat, _ := cmd.Flags().GetString("output")
		if outputFormat == "" {
			outputFormat = cfg.Defaults.OutputFormat
		}
		
		// Create formatter
		formatter := output.NewFormatter(output.Format(outputFormat), os.Stdout)
		
		// Get reviews
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		reviews, err := client.GetReviews(ctx, appID, limit)
		if err != nil {
			return fmt.Errorf("failed to get reviews: %w", err)
		}
		
		// Format and output reviews
		return formatter.Format(reviews.Data)
	},
}

var reviewsWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for new reviews",
	Long:  `Continuously watches for new reviews.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get app ID
		appID, _ := cmd.Flags().GetString("app")
		if appID == "" {
			return fmt.Errorf("app ID is required, please provide with --app")
		}
		
		// Get polling interval
		interval, _ := cmd.Flags().GetInt("interval")
		if interval < 60 {
			fmt.Println("Warning: Polling interval is very low, which might cause rate limiting.")
		}
		
		// Get minimum rating
		minRating, _ := cmd.Flags().GetInt("min-rating")
		
		fmt.Printf("Watching for new reviews for app %s (refresh: %ds, min-rating: %d)\n", 
			appID, interval, minRating)
		fmt.Println("Press Ctrl+C to stop")
		
		// TODO: Implement actual review watching
		// This would involve polling the reviews endpoint at regular intervals
		// and displaying new reviews as they come in
		fmt.Println("Review watching will be implemented in a future version.")
		return nil
	},
}

func init() {
	reviewsCmd.AddCommand(reviewsListCmd)
	reviewsCmd.AddCommand(reviewsWatchCmd)
	
	// Add flags for reviews list command
	reviewsListCmd.Flags().String("app", "", "App ID (required)")
	reviewsListCmd.Flags().Int("limit", 50, "Maximum number of reviews to retrieve")
	reviewsListCmd.Flags().String("sort", "-createdDate", "Sort order (-createdDate, rating, territory)")
	reviewsListCmd.Flags().Int("min-rating", 0, "Filter by minimum rating (1-5)")
	reviewsListCmd.Flags().String("territory", "", "Filter by territory (ISO 3166-1 alpha-2 country code)")
	
	// Add flags for reviews watch command
	reviewsWatchCmd.Flags().String("app", "", "App ID (required)")
	reviewsWatchCmd.Flags().Int("interval", 300, "Polling interval in seconds")
	reviewsWatchCmd.Flags().Int("min-rating", 0, "Filter by minimum rating (1-5)")
	reviewsWatchCmd.MarkFlagRequired("app")
}