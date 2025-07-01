package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config contains all the configuration settings for the application
type Config struct {
	Auth     AuthConfig     `json:"auth"`
	API      APIConfig      `json:"api"`
	Defaults DefaultsConfig `json:"defaults"`
}

type AuthConfig struct {
	KeyID          string `json:"key_id"`
	IssuerID       string `json:"issuer_id"`
	PrivateKeyPath string `json:"private_key_path"`
}

type APIConfig struct {
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout"`
}

type DefaultsConfig struct {
	OutputFormat string `json:"output_format"`
	VendorNumber string `json:"vendor_number"`
}

// Load reads the config file and returns a Config struct
func Load() (*Config, error) {
	// Default config
	config := &Config{
		API: APIConfig{
			BaseURL: "https://api.appstoreconnect.apple.com/v1",
			Timeout: 30,
		},
		Defaults: DefaultsConfig{
			OutputFormat: "table",
		},
	}

	// Find config file
	configPath := findConfigFile()
	if configPath == "" {
		// No config file found, check environment variables
		loadFromEnv(config)
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configPath, err)
	}

	// Parse YAML manually (simple format)
	if err := parseYAML(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file %s: %w", configPath, err)
	}

	// Override with environment variables
	loadFromEnv(config)

	return config, nil
}

// findConfigFile looks for the config file in standard locations
func findConfigFile() string {
	// Check locations in order
	locations := []string{
		"pomme.yaml",
		"pomme.yml",
		".pomme.yaml",
		".pomme.yml",
	}

	// Add user config directory
	if configHome, err := os.UserConfigDir(); err == nil {
		for _, name := range []string{"pomme.yaml", "pomme.yml"} {
			locations = append(locations, filepath.Join(configHome, "pomme", name))
		}
	}

	// Add home directory
	if home, err := os.UserHomeDir(); err == nil {
		for _, name := range []string{".pomme.yaml", ".pomme.yml", ".config/pomme/pomme.yaml"} {
			locations = append(locations, filepath.Join(home, name))
		}
	}

	// Check each location
	for _, loc := range locations {
		if info, err := os.Stat(loc); err == nil && !info.IsDir() {
			return loc
		}
	}

	return ""
}

// parseYAML parses a simple YAML format into the config struct
func parseYAML(data []byte, config *Config) error {
	lines := strings.Split(string(data), "\n")
	section := ""
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section headers
		if !strings.HasPrefix(line, " ") && strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		// Assign values based on section and key
		switch section {
		case "auth":
			switch key {
			case "key_id":
				config.Auth.KeyID = value
			case "issuer_id":
				config.Auth.IssuerID = value
			case "private_key_path":
				config.Auth.PrivateKeyPath = value
			}
		case "api":
			switch key {
			case "base_url":
				config.API.BaseURL = value
			case "timeout":
				fmt.Sscanf(value, "%d", &config.API.Timeout)
			}
		case "defaults":
			switch key {
			case "output_format":
				config.Defaults.OutputFormat = value
			case "vendor_number":
				config.Defaults.VendorNumber = value
			}
		}
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *Config) {
	// Auth settings
	if v := os.Getenv("POMME_AUTH_KEY_ID"); v != "" {
		config.Auth.KeyID = v
	}
	if v := os.Getenv("POMME_AUTH_ISSUER_ID"); v != "" {
		config.Auth.IssuerID = v
	}
	if v := os.Getenv("POMME_AUTH_PRIVATE_KEY_PATH"); v != "" {
		config.Auth.PrivateKeyPath = v
	}

	// API settings
	if v := os.Getenv("POMME_API_BASE_URL"); v != "" {
		config.API.BaseURL = v
	}
	if v := os.Getenv("POMME_API_TIMEOUT"); v != "" {
		fmt.Sscanf(v, "%d", &config.API.Timeout)
	}

	// Defaults
	if v := os.Getenv("POMME_DEFAULTS_OUTPUT_FORMAT"); v != "" {
		config.Defaults.OutputFormat = v
	}
	if v := os.Getenv("POMME_DEFAULTS_VENDOR_NUMBER"); v != "" {
		config.Defaults.VendorNumber = v
	}
}

// InitConfig creates a new config file with default values
func InitConfig(configPath string) error {
	if configPath == "" {
		configHome, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("could not find user config directory: %w", err)
		}
		configPath = filepath.Join(configHome, "pomme")
	}

	// Ensure the directory exists
	if err := os.MkdirAll(configPath, 0o755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	// Create a default config
	configFile := filepath.Join(configPath, "pomme.yaml")
	
	defaultConfig := `api:
  base_url: https://api.appstoreconnect.apple.com/v1
  timeout: 30
defaults:
  output_format: table
  vendor_number: ""
auth:
  key_id: ""
  issuer_id: ""
  private_key_path: ""
`

	if err := os.WriteFile(configFile, []byte(defaultConfig), 0o644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	fmt.Printf("Config file created at: %s\n", configFile)
	return nil
}

// Save writes the config to the default location
func Save(cfg *Config) (string, error) {
	configHome, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not find user config directory: %w", err)
	}
	
	configPath := filepath.Join(configHome, "pomme")
	if err := os.MkdirAll(configPath, 0o755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	
	configFile := filepath.Join(configPath, "pomme.yaml")
	
	// Format the config as YAML
	content := fmt.Sprintf(`api:
  base_url: %s
  timeout: %d
defaults:
  output_format: %s
  vendor_number: "%s"
auth:
  key_id: %s
  issuer_id: %s
  private_key_path: %s
`,
		cfg.API.BaseURL,
		cfg.API.Timeout,
		cfg.Defaults.OutputFormat,
		cfg.Defaults.VendorNumber,
		cfg.Auth.KeyID,
		cfg.Auth.IssuerID,
		cfg.Auth.PrivateKeyPath,
	)
	
	if err := os.WriteFile(configFile, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("could not write config file: %w", err)
	}
	
	return configFile, nil
}

// GetConfigPath returns the path to the config file if it exists
func GetConfigPath() string {
	path := findConfigFile()
	if path == "" {
		// Return default location even if file doesn't exist
		if configHome, err := os.UserConfigDir(); err == nil {
			return filepath.Join(configHome, "pomme", "pomme.yaml")
		}
	}
	return path
}