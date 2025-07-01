package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcusziade/pomme/internal/client"
	"github.com/marcusziade/pomme/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Pomme configuration",
	Long:  `Set up and manage your App Store Connect credentials and preferences.`,
	Example: `  # Interactive setup wizard (recommended for first-time users)
  pomme config init

  # View current configuration
  pomme config show

  # Validate your configuration
  pomme config validate

  # Show detailed setup instructions
  pomme config help`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup wizard for App Store Connect credentials",
	Long: `Guides you through setting up Pomme with your App Store Connect API credentials.

This wizard will help you:
  1. Create an API key in App Store Connect
  2. Download your private key file
  3. Configure Pomme with your credentials
  4. Validate that everything works`,
	RunE: runConfigInit,
}

var configShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"view"},
	Short:   "Display current configuration",
	Long:    `Shows your current configuration with sensitive values partially masked for security.`,
	RunE:    runConfigShow,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your configuration and API credentials",
	Long:  `Tests your configuration by attempting to connect to the App Store Connect API.`,
	RunE:  runConfigValidate,
}

var configHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show detailed setup instructions",
	Long:  `Displays step-by-step instructions for getting your App Store Connect API credentials.`,
	RunE:  runConfigHelp,
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configHelpCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	fmt.Println(colorBold + "🍎 Welcome to Pomme Setup Wizard" + colorReset)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
	fmt.Println("This wizard will help you set up Pomme to work with your")
	fmt.Println("App Store Connect account. You'll need:")
	fmt.Println()
	fmt.Println("  • An Apple Developer account")
	fmt.Println("  • Admin or Account Holder role in App Store Connect")
	fmt.Println("  • About 5 minutes to complete the setup")
	fmt.Println()
	
	if !askYesNo("Ready to begin?", true) {
		fmt.Println("\nSetup cancelled. Run 'pomme config init' when you're ready.")
		return nil
	}

	fmt.Println("\n" + colorBold + "📋 Step 1: Create an API Key" + colorReset)
	fmt.Println(strings.Repeat("─", 40))
	fmt.Println()
	fmt.Println("1. Open " + colorCyan + "https://appstoreconnect.apple.com/access/api" + colorReset)
	fmt.Println("2. Click the " + colorGreen + "+" + colorReset + " button to create a new key")
	fmt.Println("3. Give it a name like " + colorYellow + "\"Pomme CLI\"" + colorReset)
	fmt.Println("4. Choose " + colorYellow + "\"Admin\"" + colorReset + " access (or appropriate role)")
	fmt.Println("5. Click " + colorGreen + "\"Generate\"" + colorReset)
	fmt.Println()
	
	if !askYesNo("Have you created the API key?", false) {
		fmt.Println("\nPlease create the API key first, then run 'pomme config init' again.")
		return nil
	}

	fmt.Println("\n" + colorBold + "📥 Step 2: Download and Save Your Private Key" + colorReset)
	fmt.Println(strings.Repeat("─", 40))
	fmt.Println()
	fmt.Println("1. Click " + colorGreen + "\"Download API Key\"" + colorReset + " (you can only download it once!)")
	fmt.Println("2. Save the .p8 file somewhere secure, like:")
	fmt.Println("   • " + colorCyan + "~/Downloads/AuthKey_XXXXXXXXXX.p8" + colorReset)
	fmt.Println("   • " + colorCyan + "~/.config/pomme/AuthKey_XXXXXXXXXX.p8" + colorReset)
	fmt.Println()
	
	privateKeyPath := askString("Enter the path to your .p8 file", "")
	if privateKeyPath == "" {
		fmt.Println("\nSetup cancelled. Private key path is required.")
		return nil
	}

	// Expand ~ to home directory
	if strings.HasPrefix(privateKeyPath, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			privateKeyPath = filepath.Join(home, privateKeyPath[2:])
		}
	}

	// Verify the file exists
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		fmt.Printf("\n" + colorRed + "Error: File not found: %s" + colorReset + "\n", privateKeyPath)
		return fmt.Errorf("private key file not found")
	}

	fmt.Println("\n" + colorBold + "🔑 Step 3: Enter Your Credentials" + colorReset)
	fmt.Println(strings.Repeat("─", 40))
	fmt.Println()
	fmt.Println("From the API Keys page, you'll need:")
	fmt.Println()
	
	keyID := askString("Key ID (e.g., 73TT63DP5R)", "")
	if keyID == "" {
		fmt.Println("\nSetup cancelled. Key ID is required.")
		return nil
	}

	issuerID := askString("Issuer ID (e.g., a5ebdab5-0ceb-463c-8151-195b902f117b)", "")
	if issuerID == "" {
		fmt.Println("\nSetup cancelled. Issuer ID is required.")
		return nil
	}

	vendorNumber := askString("Vendor Number (optional, e.g., 93036463)", "")

	// Create config
	cfg := &config.Config{
		API: config.APIConfig{
			BaseURL: "https://api.appstoreconnect.apple.com/v1",
			Timeout: 30,
		},
		Auth: config.AuthConfig{
			KeyID:          keyID,
			IssuerID:       issuerID,
			PrivateKeyPath: privateKeyPath,
		},
		Defaults: config.DefaultsConfig{
			OutputFormat: "table",
			VendorNumber: vendorNumber,
		},
	}

	// Save config
	configPath, err := config.Save(cfg)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("\n" + colorGreen + "✅ Configuration saved to: " + configPath + colorReset)

	// Offer to validate
	fmt.Println("\n" + colorBold + "🧪 Step 4: Validate Configuration" + colorReset)
	fmt.Println(strings.Repeat("─", 40))
	
	if askYesNo("Would you like to validate your configuration now?", true) {
		fmt.Println()
		return runConfigValidate(cmd, args)
	}

	fmt.Println("\n" + colorGreen + "✨ Setup complete!" + colorReset)
	fmt.Println("\nYou can now use Pomme commands like:")
	fmt.Println("  • " + colorCyan + "pomme sales" + colorReset + " - View your latest sales")
	fmt.Println("  • " + colorCyan + "pomme apps list" + colorReset + " - List your apps")
	fmt.Println("  • " + colorCyan + "pomme reviews list <app-id>" + colorReset + " - View app reviews")
	fmt.Println("\nRun " + colorYellow + "pomme --help" + colorReset + " for more commands.")

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(colorRed + "No configuration found." + colorReset)
			fmt.Println("\nRun " + colorYellow + "pomme config init" + colorReset + " to set up Pomme.")
			return nil
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath := config.GetConfigPath()
	
	fmt.Println(colorBold + "🔧 Current Configuration" + colorReset)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("\nConfig file: " + colorCyan + "%s" + colorReset + "\n", configPath)
	
	fmt.Println("\n" + colorBold + "API Settings:" + colorReset)
	fmt.Printf("  Base URL: %s\n", cfg.API.BaseURL)
	fmt.Printf("  Timeout: %d seconds\n", cfg.API.Timeout)
	
	fmt.Println("\n" + colorBold + "Authentication:" + colorReset)
	fmt.Printf("  Key ID: %s\n", maskString(cfg.Auth.KeyID))
	fmt.Printf("  Issuer ID: %s\n", maskString(cfg.Auth.IssuerID))
	fmt.Printf("  Private Key: %s\n", cfg.Auth.PrivateKeyPath)
	
	// Check if private key exists
	if _, err := os.Stat(cfg.Auth.PrivateKeyPath); os.IsNotExist(err) {
		fmt.Printf("               " + colorRed + "⚠️  File not found!" + colorReset + "\n")
	} else {
		fmt.Printf("               " + colorGreen + "✓ File exists" + colorReset + "\n")
	}
	
	fmt.Println("\n" + colorBold + "Defaults:" + colorReset)
	fmt.Printf("  Output Format: %s\n", cfg.Defaults.OutputFormat)
	if cfg.Defaults.VendorNumber != "" {
		fmt.Printf("  Vendor Number: %s\n", maskString(cfg.Defaults.VendorNumber))
	}
	
	fmt.Println("\n" + colorGray + "Run 'pomme config validate' to test your configuration." + colorReset)
	
	return nil
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	fmt.Println(colorBold + "🧪 Validating Configuration" + colorReset)
	fmt.Println(strings.Repeat("─", 40))
	
	// Check config exists
	fmt.Print("\nChecking config file... ")
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(colorRed + "✗" + colorReset)
		if os.IsNotExist(err) {
			fmt.Println("\nNo configuration found. Run " + colorYellow + "pomme config init" + colorReset + " first.")
		} else {
			fmt.Printf("\nError: %v\n", err)
		}
		return err
	}
	fmt.Println(colorGreen + "✓" + colorReset)

	// Check required fields
	fmt.Print("Checking credentials... ")
	missing := []string{}
	if cfg.Auth.KeyID == "" {
		missing = append(missing, "Key ID")
	}
	if cfg.Auth.IssuerID == "" {
		missing = append(missing, "Issuer ID")
	}
	if cfg.Auth.PrivateKeyPath == "" {
		missing = append(missing, "Private Key Path")
	}
	
	if len(missing) > 0 {
		fmt.Println(colorRed + "✗" + colorReset)
		fmt.Printf("\nMissing required fields: %s\n", strings.Join(missing, ", "))
		fmt.Println("\nRun " + colorYellow + "pomme config init" + colorReset + " to complete setup.")
		return fmt.Errorf("incomplete configuration")
	}
	fmt.Println(colorGreen + "✓" + colorReset)

	// Check private key file
	fmt.Print("Checking private key file... ")
	if _, err := os.Stat(cfg.Auth.PrivateKeyPath); os.IsNotExist(err) {
		fmt.Println(colorRed + "✗" + colorReset)
		fmt.Printf("\nPrivate key not found: %s\n", cfg.Auth.PrivateKeyPath)
		return fmt.Errorf("private key file not found")
	}
	fmt.Println(colorGreen + "✓" + colorReset)

	// Test API connection
	fmt.Print("Testing API connection... ")
	apiClient, err := client.New(cfg)
	if err != nil {
		fmt.Println(colorRed + "✗" + colorReset)
		fmt.Printf("\nFailed to create client: %v\n", err)
		return err
	}

	// Try to list apps as a simple test
	ctx := context.Background()
	resp, err := apiClient.Get(ctx, "/apps?limit=1")
	if err != nil {
		fmt.Println(colorRed + "✗" + colorReset)
		fmt.Printf("\nAPI request failed: %v\n", err)
		return err
	}
	resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Println(colorRed + "✗" + colorReset)
		fmt.Println("\nAuthentication failed. Please check your credentials.")
		fmt.Println("\nMake sure:")
		fmt.Println("  • Your API key has not been revoked")
		fmt.Println("  • The Key ID and Issuer ID are correct")
		fmt.Println("  • The private key file matches the Key ID")
		return fmt.Errorf("authentication failed")
	}

	if resp.StatusCode >= 400 {
		fmt.Println(colorRed + "✗" + colorReset)
		fmt.Printf("\nAPI returned error: %d\n", resp.StatusCode)
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	fmt.Println(colorGreen + "✓" + colorReset)

	fmt.Println("\n" + colorGreen + "✅ Configuration is valid!" + colorReset)
	fmt.Println("\nYou're all set to use Pomme. Try these commands:")
	fmt.Println("  • " + colorCyan + "pomme sales" + colorReset + " - View your sales reports")
	fmt.Println("  • " + colorCyan + "pomme apps list" + colorReset + " - List your apps")
	
	return nil
}

func runConfigHelp(cmd *cobra.Command, args []string) error {
	fmt.Println(colorBold + "📚 App Store Connect API Setup Guide" + colorReset)
	fmt.Println(strings.Repeat("═", 60))
	
	fmt.Println("\n" + colorBold + "Prerequisites:" + colorReset)
	fmt.Println("  • Apple Developer account")
	fmt.Println("  • Admin, Account Holder, or appropriate role in App Store Connect")
	
	fmt.Println("\n" + colorBold + "Step-by-Step Instructions:" + colorReset)
	fmt.Println()
	fmt.Println(colorYellow + "1. Sign in to App Store Connect" + colorReset)
	fmt.Println("   Visit: " + colorCyan + "https://appstoreconnect.apple.com" + colorReset)
	
	fmt.Println("\n" + colorYellow + "2. Navigate to Users and Access" + colorReset)
	fmt.Println("   Click on your account name → Users and Access")
	
	fmt.Println("\n" + colorYellow + "3. Go to the Keys tab" + colorReset)
	fmt.Println("   Select \"Keys\" under \"App Store Connect API\"")
	
	fmt.Println("\n" + colorYellow + "4. Create a new API Key" + colorReset)
	fmt.Println("   • Click the \"+\" button")
	fmt.Println("   • Enter a name (e.g., \"Pomme CLI\")")
	fmt.Println("   • Select access level:")
	fmt.Println("     - " + colorGreen + "Admin" + colorReset + " (recommended for full access)")
	fmt.Println("     - " + colorGreen + "Finance" + colorReset + " (for sales reports only)")
	fmt.Println("     - " + colorGreen + "Sales" + colorReset + " (limited sales access)")
	
	fmt.Println("\n" + colorYellow + "5. Download the Private Key" + colorReset)
	fmt.Println("   • Click \"Generate\"")
	fmt.Println("   • " + colorRed + "IMPORTANT:" + colorReset + " Download the .p8 file immediately")
	fmt.Println("   • You can only download it once!")
	fmt.Println("   • Save it securely, for example:")
	fmt.Println("     " + colorCyan + "~/.config/pomme/AuthKey_XXXXXXXXXX.p8" + colorReset)
	
	fmt.Println("\n" + colorYellow + "6. Note Your Credentials" + colorReset)
	fmt.Println("   From the Keys page, you'll need:")
	fmt.Println("   • " + colorGreen + "Key ID" + colorReset + " (10 characters, e.g., 73TT63DP5R)")
	fmt.Println("   • " + colorGreen + "Issuer ID" + colorReset + " (UUID format)")
	
	fmt.Println("\n" + colorYellow + "7. Find Your Vendor Number (Optional)" + colorReset)
	fmt.Println("   • Go to \"Payments and Financial Reports\"")
	fmt.Println("   • Your vendor number is at the top of the page")
	fmt.Println("   • Format: 8 digits (e.g., 93036463)")
	
	fmt.Println("\n" + colorBold + "Security Tips:" + colorReset)
	fmt.Println("  🔒 Keep your .p8 file secure and never share it")
	fmt.Println("  🔒 Consider restricting the API key's access level")
	fmt.Println("  🔒 Store credentials in environment variables for scripts")
	
	fmt.Println("\n" + colorBold + "Environment Variables:" + colorReset)
	fmt.Println("  You can also use environment variables:")
	fmt.Println("  " + colorCyan + "export POMME_AUTH_KEY_ID=your_key_id" + colorReset)
	fmt.Println("  " + colorCyan + "export POMME_AUTH_ISSUER_ID=your_issuer_id" + colorReset)
	fmt.Println("  " + colorCyan + "export POMME_AUTH_PRIVATE_KEY_PATH=/path/to/key.p8" + colorReset)
	
	fmt.Println("\n" + colorGray + "Ready? Run '" + colorYellow + "pomme config init" + colorGray + "' to start the setup wizard." + colorReset)
	
	return nil
}

// Helper functions
func askYesNo(prompt string, defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)
	defaultStr := "y/N"
	if defaultYes {
		defaultStr = "Y/n"
	}
	
	fmt.Printf("%s [%s]: ", prompt, defaultStr)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	
	if answer == "" {
		return defaultYes
	}
	
	return answer == "y" || answer == "yes"
}

func askString(prompt, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	
	if answer == "" {
		return defaultValue
	}
	
	return answer
}

// maskString returns a masked version of a string for secure display
func maskString(s string) string {
	if s == "" {
		return colorRed + "<not set>" + colorReset
	}
	
	if len(s) <= 4 {
		return "****"
	}
	
	// Show first 2 and last 2 characters
	return fmt.Sprintf("%s****%s", s[:2], s[len(s)-2:])
}

