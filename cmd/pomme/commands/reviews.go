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
	"github.com/marcusziade/pomme/internal/models"
	"github.com/marcusziade/pomme/internal/services/reviews"
	"github.com/spf13/cobra"
)

var reviewsCmd = &cobra.Command{
	Use:     "reviews",
	Short:   "Manage customer reviews and ratings",
	Long:    "View, analyze, and respond to customer reviews from the App Store",
	Aliases: []string{"review", "ratings"},
}

var reviewsListCmd = &cobra.Command{
	Use:     "list <app-id>",
	Short:   "List customer reviews",
	Long:    "Display customer reviews with filtering options",
	Args:    cobra.ExactArgs(1),
	RunE:    runReviewsList,
}

var reviewsSummaryCmd = &cobra.Command{
	Use:     "summary <app-id>",
	Short:   "Show review summary and statistics",
	Long:    "Display aggregated review statistics including ratings distribution and territory breakdown",
	Args:    cobra.ExactArgs(1),
	RunE:    runReviewsSummary,
}

var reviewsRespondCmd = &cobra.Command{
	Use:     "respond <review-id> <response-text>",
	Short:   "Respond to a customer review",
	Long:    "Create or update a response to a customer review",
	Args:    cobra.ExactArgs(2),
	RunE:    runReviewsRespond,
}

var (
	// Flags
	reviewsRating    int
	reviewsLimit     int
	reviewsSort      string
	reviewsVerbose   bool
)

func init() {
	// Add subcommands
	reviewsCmd.AddCommand(reviewsListCmd)
	reviewsCmd.AddCommand(reviewsSummaryCmd)
	reviewsCmd.AddCommand(reviewsRespondCmd)
	
	// Add flags for list command
	reviewsListCmd.Flags().IntVar(&reviewsRating, "rating", 0, "Filter by rating (1-5)")
	reviewsListCmd.Flags().IntVar(&reviewsLimit, "limit", 20, "Number of reviews to display")
	reviewsListCmd.Flags().StringVar(&reviewsSort, "sort", "recent", "Sort order (recent, critical, helpful)")
	reviewsListCmd.Flags().BoolVar(&reviewsVerbose, "verbose", false, "Show full review content")
}

func runReviewsList(cmd *cobra.Command, args []string) error {
	appID := args[0]
	
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create client and service
	apiClient, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	
	svc := reviews.NewService(apiClient)
	
	fmt.Printf("📱 Fetching reviews for app %s...\n\n", appID)
	
	// Create filter
	filter := models.ReviewFilter{
		AppID:     appID,
		Rating:    reviewsRating,
		Limit:     reviewsLimit,
		Sort:      reviewsSort,
	}
	
	// Fetch reviews
	ctx := context.Background()
	reviewList, err := svc.GetReviews(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to fetch reviews: %w", err)
	}
	
	// Display reviews
	displayReviews(reviewList, reviewsVerbose)
	
	return nil
}

func runReviewsSummary(cmd *cobra.Command, args []string) error {
	appID := args[0]
	
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create client and service
	apiClient, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	
	svc := reviews.NewService(apiClient)
	
	fmt.Printf("📊 Generating review summary for app %s...\n\n", appID)
	
	// Fetch summary
	ctx := context.Background()
	summary, err := svc.GetReviewSummary(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to fetch summary: %w", err)
	}
	
	// Display summary
	displayReviewSummary(summary)
	
	return nil
}

func runReviewsRespond(cmd *cobra.Command, args []string) error {
	reviewID := args[0]
	responseText := args[1]
	
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create client and service
	apiClient, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	
	svc := reviews.NewService(apiClient)
	
	fmt.Printf("💬 Responding to review %s...\n", reviewID)
	
	// Send response
	ctx := context.Background()
	if err := svc.RespondToReview(ctx, reviewID, responseText); err != nil {
		return fmt.Errorf("failed to respond to review: %w", err)
	}
	
	fmt.Println("✅ Response submitted successfully!")
	
	return nil
}


// displayReviews shows reviews in a formatted table
func displayReviews(reviews []models.CustomerReview, verbose bool) {
	if len(reviews) == 0 {
		fmt.Println("No reviews found.")
		return
	}
	
	fmt.Printf("%s📱 Customer Reviews%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 80))
	
	for i, review := range reviews {
		// Rating stars
		stars := strings.Repeat("⭐", review.Attributes.Rating)
		emptyStars := strings.Repeat("☆", 5-review.Attributes.Rating)
		
		// Review header
		fmt.Printf("\n%s%s%s %s%s%s  %s%-20s%s  %s\n",
			colorBold,
			stars,
			emptyStars,
			colorCyan,
			review.Attributes.ReviewerNickname,
			colorReset,
			colorGray,
			review.Attributes.Territory,
			colorReset,
			review.Attributes.CreatedDate.Format("2006-01-02"),
		)
		
		// Title
		if review.Attributes.Title != "" {
			fmt.Printf("%s%s%s\n", colorBold, review.Attributes.Title, colorReset)
		}
		
		// Body
		if verbose || len(review.Attributes.Body) <= 200 {
			fmt.Printf("%s\n", review.Attributes.Body)
		} else {
			// Truncate long reviews
			fmt.Printf("%s...\n", review.Attributes.Body[:200])
		}
		
		// Separator between reviews
		if i < len(reviews)-1 {
			fmt.Println(strings.Repeat("─", 40))
		}
	}
	
	fmt.Printf("\n%sShowing %d reviews%s\n", colorGray, len(reviews), colorReset)
}

// displayReviewSummary shows aggregated review statistics
func displayReviewSummary(summary *models.ReviewSummary) {
	fmt.Printf("%s📊 Review Summary%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("═", 60))
	
	// Overall stats
	fmt.Printf("\n%sOverall Statistics%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 40))
	
	// Average rating with visual representation
	avgStars := int(summary.AverageRating + 0.5)
	stars := strings.Repeat("⭐", avgStars)
	emptyStars := strings.Repeat("☆", 5-avgStars)
	
	fmt.Printf("  Average Rating: %s%.1f%s %s%s\n", 
		colorYellow, summary.AverageRating, colorReset, stars, emptyStars)
	fmt.Printf("  Total Reviews: %s%d%s\n", 
		colorCyan, summary.TotalReviews, colorReset)
	
	// Rating distribution
	fmt.Printf("\n%sRating Distribution%s\n", colorBold, colorReset)
	fmt.Println(strings.Repeat("─", 40))
	
	maxCount := 0
	for _, count := range summary.RatingCounts {
		if count > maxCount {
			maxCount = count
		}
	}
	
	// Display rating bars
	for rating := 5; rating >= 1; rating-- {
		count := summary.RatingCounts[rating]
		percentage := float64(count) / float64(summary.TotalReviews) * 100
		
		// Create visual bar
		barLength := 30
		if maxCount > 0 {
			barLength = int(float64(count) / float64(maxCount) * 30)
		}
		bar := strings.Repeat("█", barLength)
		
		fmt.Printf("  %d⭐ %s%-30s%s %4d (%5.1f%%)\n",
			rating,
			colorGreen,
			bar,
			colorReset,
			count,
			percentage,
		)
	}
	
	// Territory breakdown
	if len(summary.TerritoryStats) > 0 {
		fmt.Printf("\n%sTop Territories%s\n", colorBold, colorReset)
		fmt.Println(strings.Repeat("─", 40))
		
		// Sort and show top 10 territories
		limit := 10
		if len(summary.TerritoryStats) < limit {
			limit = len(summary.TerritoryStats)
		}
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "  Territory\tReviews\tAvg Rating\n")
		fmt.Fprintf(w, "  %s\t%s\t%s\n",
			strings.Repeat("─", 10),
			strings.Repeat("─", 7),
			strings.Repeat("─", 10))
		
		for i := 0; i < limit; i++ {
			territory := summary.TerritoryStats[i]
			fmt.Fprintf(w, "  %s\t%d\t%.1f\n",
				territory.Territory,
				territory.ReviewCount,
				territory.AverageRating,
			)
		}
		w.Flush()
	}
	
	// Recent reviews preview
	if len(summary.RecentReviews) > 0 {
		fmt.Printf("\n%sRecent Reviews%s\n", colorBold, colorReset)
		fmt.Println(strings.Repeat("─", 40))
		
		for i, review := range summary.RecentReviews {
			if i >= 3 {
				break // Show only top 3
			}
			
			stars := strings.Repeat("⭐", review.Attributes.Rating)
			fmt.Printf("\n  %s %s%-15s%s %s%s%s\n",
				stars,
				colorCyan,
				review.Attributes.ReviewerNickname,
				colorReset,
				colorGray,
				review.Attributes.CreatedDate.Format("2006-01-02"),
				colorReset,
			)
			
			if review.Attributes.Title != "" {
				fmt.Printf("  %s\n", review.Attributes.Title)
			}
		}
	}
	
	fmt.Printf("\n%sGenerated: %s%s\n", colorGray, time.Now().Format("2006-01-02 15:04:05"), colorReset)
}