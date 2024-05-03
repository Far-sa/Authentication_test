package mocks

import (
	"user-svc/ports"

	"github.com/stretchr/testify/mock"
)

type MockConfig struct {
	mock.Mock
}

func NewMockConfig() *MockConfig {
	return &MockConfig{}
}

// func (m *MockConfig) LoadConfig(filePath string) error {
// 	args := m.Called()
// 	return args.Error(0)
// }

func (m *MockConfig) GetDatabaseConfig() ports.DatabaseConfig {
	args := m.Called()
	return args.Get(0).(ports.DatabaseConfig)
}

func (m *MockConfig) GetHTTPConfig() ports.HTTPConfig {
	args := m.Called()
	return args.Get(0).(ports.HTTPConfig)
}

func (m *MockConfig) GetBrokerConfig() ports.BrokerConfig {
	args := m.Called()
	return args.Get(0).(ports.BrokerConfig)
}

func (m *MockConfig) GetConstants() ports.Constants {
	args := m.Called()
	return args.Get(0).(ports.Constants)
}

func (m *MockConfig) GetStatics() ports.Statics {
	args := m.Called()
	return args.Get(0).(ports.Statics)
}

func (m *MockConfig) GetLoggerConfig() ports.LoggerConfig {
	args := m.Called()
	return args.Get(0).(ports.LoggerConfig)
}

// MockConfig is a mock implementation of Config for testing
// type MockConfig struct {
// 	LoadConfigFn        func() error
// 	GetDatabaseConfigFn func() ports.DatabaseConfig
// 	GetHTTPConfigFn     func() ports.HTTPConfig
// 	GetConstantsFn      func() ports.Constants
// 	GetStaticsFn        func() ports.Statics
// 	GetLoggerConfigFn   func() ports.LoggerConfig
// }

// func (m *MockConfig) LoadConfig() error {
// 	if m.LoadConfigFn != nil {
// 		return m.LoadConfigFn()
// 	}
// 	return nil
// }

// func (m *MockConfig) GetDatabaseConfig() ports.DatabaseConfig {
// 	if m.GetDatabaseConfigFn != nil {
// 		return m.GetDatabaseConfigFn()
// 	}
// 	return ports.DatabaseConfig{}
// }

// func (m *MockConfig) GetHTTPConfig() ports.HTTPConfig {
// 	if m.GetHTTPConfigFn != nil {
// 		return m.GetHTTPConfigFn()
// 	}
// 	return ports.HTTPConfig{}
// }

// func (m *MockConfig) GetConstants() ports.Constants {
// 	if m.GetConstantsFn != nil {
// 		return m.GetConstantsFn()
// 	}
// 	return ports.Constants{}
// }

// func (m *MockConfig) GetStatics() ports.Statics {
// 	if m.GetStaticsFn != nil {
// 		return m.GetStaticsFn()
// 	}
// 	return ports.Statics{}
// }

// func (m *MockConfig) GetLoggerConfig() ports.LoggerConfig {
// 	if m.GetLoggerConfigFn != nil {
// 		return m.GetLoggerConfigFn()
// 	}
// 	return ports.LoggerConfig{}
// }
