// config/config.go

package config

import (
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Internal   InternalConfig
	Delivery   DeliveryConfig
	Repository RepositoryConfig
	// Add other configuration sections as needed
}

// InternalConfig represents configuration settings for the internal service layer
type InternalConfig struct {
	// Define internal configuration settings here
}

// DeliveryConfig represents configuration settings for the delivery layer adapters
type DeliveryConfig struct {
	Port int // Example delivery layer configuration
}

// RepositoryConfig represents configuration settings for the repository layer adapters
type RepositoryConfig struct {
	DatabaseURL string // Example repository layer configuration
}

// Load loads configuration values from various sources
func Load() (*Config, error) {
	viper.SetEnvPrefix("MYAPP") // Prefix for environment variables (optional)
	viper.AutomaticEnv()        // Automatically read from environment variables

	// Load configuration values
	config := &Config{
		Internal: InternalConfig{
			// Load internal configuration settings here
		},
		Delivery: DeliveryConfig{
			Port: viper.GetInt("DELIVERY_PORT"),
		},
		Repository: RepositoryConfig{
			DatabaseURL: viper.GetString("DATABASE_URL"),
		},
		// Load other configuration sections as needed
	}

	// Perform validation or additional configuration loading if necessary

	return config, nil
}
