package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Color constants for help display
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBox    = "\033[94m"
)

var (
	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "pomme",
		Short: "🍎 Beautiful App Store Connect CLI for sales & reviews",
		Long:  getLongDescription(),
		Example: getExamples(),
	}
)

func getLongDescription() string {
	return fmt.Sprintf(`%s
%s╔═══════════════════════════════════════════════════════════════════════╗%s
%s║                                                                       ║%s
%s║      🍎 Pomme - Beautiful App Store Connect CLI                      ║%s
%s║                                                                       ║%s
%s║      Track sales, monitor reviews, and analyze your App Store        ║%s
%s║      performance with style and ease.                                 ║%s
%s║                                                                       ║%s
%s╚═══════════════════════════════════════════════════════════════════════╝%s

%s📚 QUICK START%s
────────────────────────────────────────────────────────────────────────

  1. %sConfigure your API credentials:%s
     $ pomme config init

  2. %sView your latest sales:%s
     $ pomme sales

  3. %sMonitor customer reviews:%s
     $ pomme reviews list <app-id>

%s💡 KEY FEATURES%s
────────────────────────────────────────────────────────────────────────

  %s📊 Sales Reports%s
  • View monthly sales with multi-currency support
  • Compare sales periods and track trends
  • Export data in JSON, CSV, or beautiful tables
  • Automatic currency grouping and USD conversion

  %s⭐ Review Management%s
  • List and analyze customer reviews
  • View rating distributions and summaries
  • Respond to customer feedback
  • Monitor reviews across territories

  %s🚀 Smart Features%s
  • Beautiful colored terminal output
  • Intelligent caching for better performance
  • Concurrent data fetching with Go routines
  • Clean, focused commands that do one thing well

%s🎯 COMMON WORKFLOWS%s
────────────────────────────────────────────────────────────────────────

  %sDaily Sales Check:%s
    $ pomme sales                    # Latest month
    $ pomme sales monthly --details  # With app breakdown

  %sMonthly Comparison:%s
    $ pomme sales compare --current 2025-03 --previous 2025-02

  %sReview Monitoring:%s
    $ pomme reviews summary <app-id>
    $ pomme reviews list <app-id> --rating 1  # Critical reviews

  %sData Export:%s
    $ pomme sales monthly --output csv > sales.csv
    $ pomme sales monthly --output json | jq '.revenue'

%s⚙️  CONFIGURATION%s
────────────────────────────────────────────────────────────────────────

  Config file: ~/.config/pomme/pomme.yaml
  
  Environment variables:
    POMME_AUTH_KEY_ID         - Your API key ID
    POMME_AUTH_ISSUER_ID      - Your issuer ID
    POMME_AUTH_PRIVATE_KEY_PATH - Path to .p8 key file

%s📖 MORE HELP%s
────────────────────────────────────────────────────────────────────────

  • Use 'pomme <command> --help' for detailed command info
  • Visit https://github.com/marcusziade/pomme for full docs
  • Report issues at https://github.com/marcusziade/pomme/issues

Built with ❤️  using Go  |  App Store Connect API v1`,
		colorBold, colorBox, colorReset,
		colorBox, colorReset,
		colorBox, colorReset,
		colorBox, colorReset,
		colorBox, colorReset,
		colorBox, colorReset,
		colorBox, colorReset,
		colorBox, colorReset,
		colorBold, colorReset,
		colorCyan, colorReset,
		colorCyan, colorReset,
		colorCyan, colorReset,
		colorBold, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorGreen, colorReset,
		colorBold, colorReset,
		colorYellow, colorReset,
		colorYellow, colorReset,
		colorYellow, colorReset,
		colorYellow, colorReset,
		colorBold, colorReset,
		colorBold, colorReset,
		colorBold, colorReset,
	)
}

func getExamples() string {
	return `  # First time setup
  pomme config init

  # View latest sales
  pomme sales

  # View specific month sales
  pomme sales monthly 2025-03

  # Compare sales periods
  pomme sales compare --current 2025-03 --previous 2025-02

  # List app reviews
  pomme reviews list <app-id>

  # Get review summary
  pomme reviews summary <app-id>

  # Export sales to CSV
  pomme sales monthly --output csv > march_sales.csv`
}

func init() {
	// Add global flags here
	RootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (json, csv, table)")

	// Add subcommands
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(authCmd)
	RootCmd.AddCommand(salesCmd)
	RootCmd.AddCommand(appsCmd)
	RootCmd.AddCommand(reviewsCmd)
}
