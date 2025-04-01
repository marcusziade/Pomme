package commands

import (
	"fmt"
	"os"

	"github.com/marcus/pomme/internal/auth"
	"github.com/marcus/pomme/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication management",
	Long:  `Commands for managing App Store Connect API authentication.`,
}

var authTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test authentication",
	Long:  `Tests your authentication credentials by generating a JWT token.`,
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
		
		// Create JWT config
		jwtConfig := auth.JWTConfig{
			KeyID:         cfg.Auth.KeyID,
			IssuerID:      cfg.Auth.IssuerID,
			PrivateKeyPEM: string(privateKeyData),
		}
		
		// Generate token
		token, err := auth.GenerateToken(jwtConfig)
		if err != nil {
			return fmt.Errorf("failed to generate JWT token: %w", err)
		}
		
		fmt.Println("Authentication successful!")
		fmt.Println("JWT token generated:")
		fmt.Printf("%s...\n", token[:40])
		return nil
	},
}

var authSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up authentication credentials",
	Long:  `Interactive setup for App Store Connect API credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement interactive credential setup
		fmt.Println("Interactive setup will be implemented in a future version.")
		fmt.Println("For now, please edit your config file manually.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authTestCmd)
	authCmd.AddCommand(authSetupCmd)
}
