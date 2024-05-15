package postgres_test

import (
	"auth-svc/adapters/repository/postgres"
	"auth-svc/internal/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAuthRepository(db)

	var userID uint = 1
	token := "test_token"
	expiration := time.Now().Add(24 * time.Hour)

	mock.ExpectExec("INSERT INTO tokens").
		WithArgs(userID, token, expiration).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.StoreToken(userID, token, expiration)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRetrieveToken(t *testing.T) {
	// Create a new sqlmock database
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Define test data
	var userID uint = 1
	expectedToken := &entity.Token{
		ID:         1,
		UserID:     userID,
		TokenValue: "mocked_token",
		Expiration: time.Now(),
	}

	// Set up mock expectations
	query := "SELECT id, user_id, token, expiration FROM tokens WHERE user_id = ?"
	mock.ExpectQuery(query).WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "expiration"}).
			AddRow(expectedToken.ID, expectedToken.UserID, expectedToken.TokenValue, expectedToken.Expiration))

	// Create a repository with the mock database
	repo := postgres.NewAuthRepository(mockDB)

	// Call the function under test
	token, err := repo.RetrieveToken(userID)

	// Verify the results
	assert.NoError(t, err, "Error retrieving token")
	assert.NotNil(t, token, "Expected a token, but got nil")
	assert.Equal(t, expectedToken, token, "Retrieved token does not match expected token")

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Unfulfilled expectations")
}

func TestTestRetrieveTokenTableDriven(t *testing.T) {
	type testCase struct {
		name          string
		expectedQuery string
		expectedArgs  []interface{}
		expectedError string
		expectedToken *entity.Token
	}

	cases := []testCase{
		{
			name:          "Successful token retrieval",
			expectedQuery: "SELECT id, user_id, token, expiration FROM tokens WHERE user_id = $1",
			expectedArgs:  []interface{}{1},
			expectedToken: &entity.Token{ID: 1, UserID: 1, TokenValue: "test_token", Expiration: time.Now().Add(time.Hour)}, // Replace with expected data
		},
		{
			name:          "User not found",
			expectedQuery: "SELECT id, user_id, token, expiration FROM tokens WHERE user_id = $1",
			expectedArgs:  []interface{}{1},
			expectedError: "sql: no rows in result set", // Adjust error message if needed
		},
		{
			name:          "Error during query execution",
			expectedQuery: "SELECT id, user_id, token, expiration FROM tokens WHERE user_id = $1",
			expectedArgs:  []interface{}{1},
			expectedError: "some_database_error", // Replace with expected error message
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)

			mock.ExpectQuery(c.expectedQuery).WithArgs(c.expectedArgs).
				WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "expiration"}).
					AddRow(c.expectedToken.ID, c.expectedToken.UserID, c.expectedToken.TokenValue, c.expectedToken.Expiration))

			db := postgres.NewAuthRepository(mockDB)
			token, err := db.RetrieveToken(c.expectedToken.UserID)

			// Assert that the token
			if c.expectedError != "" {
				require.EqualError(t, err, c.expectedError) // Verify expected error
			} else {
				require.NoError(t, err) // No error expected
			}

			// Verify if a token was returned (if expected)
			if c.expectedToken != nil {
				require.NotNil(t, token)                   // Token should not be nil
				require.Equal(t, *c.expectedToken, *token) // Compare token details
			}

			// Ensure expected SQL statement was executed
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestStoreTokenTableDriven(t *testing.T) {

	type testCase struct {
		name          string
		expectedQuery string
		expectedArgs  []interface{}
		expectedError string
	}

	cases := []testCase{
		{
			name:          "Successful token storage",
			expectedQuery: "INSERT INTO tokens (user_id, token, expiration) VALUES ($1, $2, $3)",
			expectedArgs:  []interface{}{1, "test_token", time.Now().Add(24 * time.Hour)},
			expectedError: "",
		},
		{
			name:          "Error on database execution",
			expectedQuery: "INSERT INTO tokens (user_id, token, expiration) VALUES ($1, $2, $3)",
			expectedArgs:  []interface{}{1, "test_token", time.Now().Add(24 * time.Hour)},
			expectedError: "some_database_error", // Replace with expected error message
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// Create a mock database connection
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			// Configure expected SQL statement and arguments based on the test case
			mock.ExpectExec(c.expectedQuery).WithArgs(c.expectedArgs).
				WillReturnResult(sqlmock.NewResult(1, 1))

			db := postgres.NewAuthRepository(mockDB)
			err = db.StoreToken(1, "test_token", time.Now().Add(time.Hour))

			// Assertions
			if c.expectedError != "" {
				require.EqualError(t, err, c.expectedError) // Verify expected error
			} else {
				require.NoError(t, err) // No error expected
			}
			assert.NoError(t, mock.ExpectationsWereMet())

		})

	}

}
