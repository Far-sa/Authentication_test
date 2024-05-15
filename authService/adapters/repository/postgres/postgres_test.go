package postgres_test

import (
	"auth-svc/adapters/repository/postgres"
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

	userID := 1
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
