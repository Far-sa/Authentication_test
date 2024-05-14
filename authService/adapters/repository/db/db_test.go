package db_test

import (
	"auth-svc/adapters/repository/db"
	"auth-svc/internal/ports"
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockConfig struct {
	mock.Mock
}

func (m *mockConfig) GetDatabaseConfig() ports.DatabaseConfig {
	args := m.Called()
	return args.Get(0).(ports.DatabaseConfig)
}

func (m *mockConfig) GetBrokerConfig() ports.BrokerConfig {
	return m.GetBrokerConfig()
}

func (m *mockConfig) GetConstants() ports.Constants {
	return m.GetConstants()
}

func (m *mockConfig) GetHTTPConfig() ports.HTTPConfig {
	return m.GetHTTPConfig()
}
func TestGetConnectionPoolSuccess(t *testing.T) {
	// Create a mock of the Config interface
	mockCfg := &mockConfig{}

	expectedDbConfig := ports.DatabaseConfig{
		User:     "root",
		Password: "password",
		Host:     "localhost",
		Port:     5432,
		DBName:   "authDb",
	}

	mockCfg.On("GetDatabaseConfig").Return(expectedDbConfig)
	pool, err := db.GetConnectionPool(mockCfg)

	// Assert that no error occurred
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the connection pool is not nil
	if pool == nil {
		t.Errorf("Expected a connection pool object, but got nil")
	}

	// Close the connection pool after the test
	defer pool.Close()
}
