package commands

import (
	"fmt"

	"github.com/marcusziade/pomme/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Pomme configuration",
	Long:  `Commands for managing Pomme configuration settings.`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Creates a new configuration file with default values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("path")
		
		if err := config.InitConfig(configPath); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		
		fmt.Println("Configuration initialized successfully!")
		return nil
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Long:  `Displays the current configuration settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		fmt.Println("Current configuration:")
		fmt.Println("API:")
		fmt.Printf("  Base URL: %s\n", cfg.API.BaseURL)
		fmt.Printf("  Timeout: %d seconds\n", cfg.API.Timeout)
		
		fmt.Println("Auth:")
		fmt.Printf("  Key ID: %s\n", maskString(cfg.Auth.KeyID))
		fmt.Printf("  Issuer ID: %s\n", maskString(cfg.Auth.IssuerID))
		fmt.Printf("  Private Key Path: %s\n", cfg.Auth.PrivateKeyPath)
		
		fmt.Println("Defaults:")
		fmt.Printf("  Output Format: %s\n", cfg.Defaults.OutputFormat)
		fmt.Printf("  Vendor Number: %s\n", maskString(cfg.Defaults.VendorNumber))
		
		return nil
	},
}

// maskString returns a masked version of a string for secure display
func maskString(s string) string {
	if s == "" {
		return "<not set>"
	}
	
	if len(s) <= 4 {
		return "****"
	}
	
	// Show first 2 and last 2 characters
	return fmt.Sprintf("%s****%s", s[:2], s[len(s)-2:])
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configViewCmd)
	
	configInitCmd.Flags().String("path", "", "Path to create the config file (default: user config dir)")
}
