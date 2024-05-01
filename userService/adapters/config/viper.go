package config

import (
	"log"
	"user-svc/ports"

	"github.com/spf13/viper"
)

type ViperAdapter struct {
	viper *viper.Viper
}

// func NewViperAdapter(configPaths ...string) *ViperAdapter {
// 	va := &ViperAdapter{
// 		viper: viper.New(),
// 	}

// 	for _, path := range configPaths {
// 		va.viper.AddConfigPath(path)
// 	}

// 	return va
// }

// func NewViperAdapter(configFile string) (*ViperAdapter, error) {
// 	v := viper.New()
// 	v.SetConfigFile(configFile)
// 	if err := v.ReadInConfig(); err != nil {
// 		return nil, err
// 	}
// 	return &ViperAdapter{v}, nil
// }

func NewViperAdapter() *ViperAdapter {
	return &ViperAdapter{
		viper: viper.New(),
	}
}

// LoadConfig loads configuration from a YAML file.
func (va *ViperAdapter) LoadConfig(filepath string) error {
	va.viper.SetConfigFile(filepath)
	va.viper.SetConfigType("yaml")

	return va.viper.ReadInConfig()
}

// func (va *ViperAdapter) LoadConfig() error {
// 	va.viper.SetConfigName("config")
// 	va.viper.SetConfigType("yaml")

// 	if err := va.viper.ReadInConfig(); err != nil {
// 		return err
// 	}

// 	//! load env
// 	va.viper.AutomaticEnv()
// 	// err := godotenv.Load()
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

func (va *ViperAdapter) GetDatabaseConfig() ports.DatabaseConfig {
	var dbConfig ports.DatabaseConfig
	va.viper.UnmarshalKey("database", &dbConfig)
	return dbConfig
}

func (va *ViperAdapter) GetHTTPConfig() ports.HTTPConfig {
	var httpConfig ports.HTTPConfig
	if err := va.viper.UnmarshalKey("http_server", &httpConfig); err != nil {
		log.Printf("failed to unmarshal http_server config: %v", err)
	}
	log.Printf("HTTP Config: %+v", httpConfig) // Debug log

	return httpConfig
}

func (va *ViperAdapter) GetConstants() ports.Constants {
	var constants ports.Constants
	va.viper.UnmarshalKey("constants", &constants)
	return constants
}

func (va *ViperAdapter) GetStatics() ports.Statics {
	var statics ports.Statics
	va.viper.UnmarshalKey("statics", &statics)
	return statics
}

func (va *ViperAdapter) GetLoggerConfig() ports.LoggerConfig {
	var loggerConfig ports.LoggerConfig
	va.viper.UnmarshalKey("logger", &loggerConfig)
	return loggerConfig
}
