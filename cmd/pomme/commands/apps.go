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

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "App management",
	Long:  `Commands for managing your iOS/macOS/tvOS apps.`,
}

var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your apps",
	Long:  `Lists all apps you have access to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		
		// Get apps
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		apps, err := client.GetApps(ctx)
		if err != nil {
			return fmt.Errorf("failed to get apps: %w", err)
		}
		
		// Format and output apps
		return formatter.Format(apps.Data)
	},
}

var appsInfoCmd = &cobra.Command{
	Use:   "info [app-id]",
	Short: "Get app information",
	Long:  `Retrieves detailed information about a specific app.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID := args[0]
		
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
		
		// Get app info
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		app, err := client.GetApp(ctx, appID)
		if err != nil {
			return fmt.Errorf("failed to get app info: %w", err)
		}
		
		// Format and output app info
		return formatter.Format(app.Data)
	},
}

func init() {
	appsCmd.AddCommand(appsListCmd)
	appsCmd.AddCommand(appsInfoCmd)
	
	// Add flags for apps list command
	appsListCmd.Flags().Bool("include-removed", false, "Include removed apps")
	appsListCmd.Flags().String("platform", "", "Filter by platform (IOS, MAC_OS, TV_OS)")
	
	// Add flags for apps info command
	appsInfoCmd.Flags().Bool("include-versions", false, "Include app versions")
	appsInfoCmd.Flags().Bool("include-builds", false, "Include app builds")
}