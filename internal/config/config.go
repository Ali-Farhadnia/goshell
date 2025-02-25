package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Shell    ShellConfig    `mapstructure:"shell"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ShellConfig holds shell-specific configuration
type ShellConfig struct {
	Verbose     bool `mapstructure:"verbose"`
	HistorySize int  `mapstructure:"historySize"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver      string `mapstructure:"driver"`
	DSN         string `mapstructure:"dsn"`
	AutoMigrate bool   `mapstructure:"autoMigrate"`
	LogLevel    string `mapstructure:"logLevel"`
}

// Load loads the configuration from a file
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("shell.verbose", false)
	viper.SetDefault("shell.historySize", 1000)

	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.dsn", "host=localhost user=goshell password=password dbname=goshell port=5432 sslmode=disable")
	viper.SetDefault("database.autoMigrate", true)
	viper.SetDefault("database.logLevel", "silent")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
