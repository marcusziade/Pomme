package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config contains all the configuration settings for the application
type Config struct {
	Auth struct {
		KeyID          string `mapstructure:"key_id"`
		IssuerID       string `mapstructure:"issuer_id"`
		PrivateKeyPath string `mapstructure:"private_key_path"`
	}
	API struct {
		BaseURL string `mapstructure:"base_url"`
		Timeout int    `mapstructure:"timeout"`
	}
	Defaults struct {
		OutputFormat string `mapstructure:"output_format"`
		VendorNumber string `mapstructure:"vendor_number"`
	}
}

// Load reads in config file and ENV variables if set
func Load() (*Config, error) {
	// Set up config file search paths
	configHome, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config directory: %w", err)
	}

	configName := "pomme"
	configType := "yaml"
	configPaths := []string{
		".",                                   // Current directory
		filepath.Join(configHome, "pomme"),  // User config directory
	}

	// Configure viper
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	// Set default values
	v.SetDefault("api.base_url", "https://api.appstoreconnect.apple.com/v1")
	v.SetDefault("api.timeout", 30) // 30 seconds
	v.SetDefault("defaults.output_format", "table")

	// Read environment variables
	v.SetEnvPrefix("POMME")
	v.AutomaticEnv()

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		// It's ok if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
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
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	// Create a default config
	v := viper.New()
	v.SetConfigName("pomme")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// Set default values
	v.Set("api.base_url", "https://api.appstoreconnect.apple.com/v1")
	v.Set("api.timeout", 30)
	v.Set("defaults.output_format", "table")

	// Write the config file
	if err := v.SafeWriteConfig(); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
