package mysql_test

import (
	"database/sql"
	"testing"
	"user-svc/adapters/repository/db"
	"user-svc/adapters/repository/mysql"
	"user-svc/ports"
	mocks "user-svc/ports/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMysqlDB(t *testing.T) {

	// Mock dependencies
	mockConfig := mocks.NewMockConfig()
	//mockMetrics := mocks.NewMockDatabaseMetrics()
	mockLogger := mocks.NewMockLogger()

	// Set expectations for mock config
	expectedDbConfig := ports.DatabaseConfig{
		User:     "test_user",
		Password: "test_password",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "test_database",
	}
	_ = expectedDbConfig

	// Set expectations for the mock objects
	mockConfig.On("GetDatabaseConfig").Return(expectedDbConfig)

	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Once()

	// Call the function under test
	dbPool, _ := db.GetConnectionPool(mockConfig)
	mysqlDB := mysql.New(dbPool, mockLogger)

	// Assertions
	assert.NotNil(t, mysqlDB, "MysqlDB instance should not be nil")
	assert.NotNil(t, mysqlDB.IsDBConnected(), "Database connection should not be nil")
}

// Integration test to test database connection
func TestDatabaseConnection(t *testing.T) {
	// Connect to the database
	db, err := sql.Open("mysql", "testuser:testpass@tcp(localhost:3306)/testdb")
	assert.NoError(t, err, "Failed to open database connection")
	assert.NotNil(t, db, "Database connection should not be nil")

	// Ping the database
	err = db.Ping()
	assert.NoError(t, err, "Failed to ping database")
}
