package db_test

import (
	"auth-svc/adapters/repository/db"
	"auth-svc/internal/ports"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockConfigAdapter struct {
	mock.Mock
}

func (m *MockConfigAdapter) GetDatabaseConfig() ports.DatabaseConfig {
	args := m.Called()
	return args.Get(0).(ports.DatabaseConfig)
}

func TestGetConnectionPool(t *testing.T) {
	// Create a mock configuration adapter
	mockAdapter := new(MockConfigAdapter)

	// Define expected database configuration
	expectedDbConfig := ports.DatabaseConfig{
		User:     "test_user",
		Password: "test_password",
		Host:     "localhost",
		Port:     5432,
		DBName:   "test_database",
	}

	mockAdapter.On("GetDatabaseConfig").Return(expectedDbConfig)

	pool, err := db.GetConnectionPool(mockAdapter)

	require.NoError(t, err)

	// Verify that pool is not nil
	require.NotNil(t, pool)

	err = pool.Ping()
	require.NoError(t, err)

	// Assertions on mock interactions (optional)
	mockAdapter.AssertExpectations(t)
}

func (m *MockConfigAdapter) GetBrokerConfig() ports.BrokerConfig {
	return m.GetBrokerConfig()
}

func (m *MockConfigAdapter) GetConstants() ports.Constants {
	return m.GetConstants()
}

func (m *MockConfigAdapter) GetHTTPConfig() ports.HTTPConfig {
	return m.GetHTTPConfig()
}

func TestGetConnectionPoolTableDriven(t *testing.T) {
	type testCase struct {
		name             string
		expectedDbConfig ports.DatabaseConfig
		expectedError    string
	}

	cases := []testCase{
		{name: "successful connection",
			expectedDbConfig: ports.DatabaseConfig{
				User:     "test_user", // Replace with your test credentials
				Password: "test_password",
				Host:     "localhost",
				Port:     5432,
				DBName:   "test_database",
			},
			expectedError: "",
		},
		{
			name: "Error: Invalid user in configuration",
			expectedDbConfig: ports.DatabaseConfig{ // Invalid user
				User:     "invalid_user",
				Password: "test_password",
				Host:     "localhost",
				Port:     5432,
				DBName:   "test_database",
			},
			expectedError: "pq: password authentication failed for user \"invalid_user\"", // Expected error message
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			mockAdapter := new(MockConfigAdapter)
			mockAdapter.On("GetDatabaseConfig").Return(c.expectedDbConfig)
			pool, err := db.GetConnectionPool(mockAdapter)

			// Assertions based on the test case
			if c.expectedError != "" {
				require.EqualError(t, err, c.expectedError)
			} else {
				require.NoError(t, err)
				defer pool.Close() // Close the connection pool after the test
			}

		})
	}
}

//!!!!
// var (
// 	testDbUser     string
// 	testDbPassword string
// 	testDbHost     string
// 	// Add other credentials as needed (port, database name, etc.)
// )

// func TestMain(m *testing.M) {
// 	// Load configuration before tests
// 	viper.SetConfigFile("config_test.json") // Replace with your config file path
// 	err := viper.ReadInConfig()
// 	require.NoError(t, err)
// 	testDbUser = viper.GetString("database.user")
// 	testDbPassword = viper.GetString("database.password")
// 	testDbHost = viper.GetString("database.host")
// 	// Add logic to retrieve other credentials from config

// 	// Run tests
// 	m.Run()
// }

// func TestRealDatabaseConnection(t *testing.T) {
// 	// Use the loaded credentials to establish a connection
// 	db, err := your_database_package.OpenConnection(testDbUser, testDbPassword, testDbHost, /* other connection parameters */)
// 	require.NoError(t, err)
// 	defer db.Close() // Close the connection after the test

// 	// Test the connection using Ping
// 	err = db.Ping(context.Background())
// 	require.NoError(t, err)
// }
