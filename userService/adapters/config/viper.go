package config

import (
	"fmt"
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

func NewViperAdapter() (*ViperAdapter, error) {

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".") // Use provided filepath directly
	v.AutomaticEnv()
	//v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return &ViperAdapter{viper: v}, nil

}

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

func (va *ViperAdapter) GetBrokerConfig() ports.BrokerConfig {
	var brokerConfig ports.BrokerConfig
	if err := va.viper.UnmarshalKey("rabbitmq", &brokerConfig); err != nil {
		log.Printf("failed to unmarshal broker config: %v", err)
	}
	log.Printf("Broker Config: %+v", brokerConfig) // Debug log

	return brokerConfig
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
