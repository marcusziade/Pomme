package commands

import (
	"github.com/spf13/cobra"
)

var (
	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "pomme",
		Short: "App Store Connect CLI tool",
		Long: `Pomme is a powerful App Store Connect CLI tool that allows you to
interact with the App Store Connect API to manage your apps, view sales reports,
monitor reviews, and access analytics data.`,
	}
)

func init() {
	// Add global flags here
	RootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (json, csv, table)")

	// Add subcommands
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(authCmd)
	RootCmd.AddCommand(salesCmd)
	RootCmd.AddCommand(appsCmd)
	RootCmd.AddCommand(reviewsCmd)
	RootCmd.AddCommand(analyticsCmd)
}
