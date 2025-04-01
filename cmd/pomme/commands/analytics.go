package commands

import (
	"github.com/spf13/cobra"
)

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "App analytics",
	Long:  `Commands for retrieving and analyzing app metrics.`,
}

var analyticsUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Get app usage data",
	Long:  `Retrieves app usage metrics like active devices, sessions, and crash data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement analytics usage retrieval
		return nil
	},
}

var analyticsSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Get subscription metrics",
	Long:  `Retrieves subscription metrics like conversions, renewals, and cancellations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement subscription analytics
		return nil
	},
}

func init() {
	analyticsCmd.AddCommand(analyticsUsageCmd)
	analyticsCmd.AddCommand(analyticsSubscriptionsCmd)
	
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
