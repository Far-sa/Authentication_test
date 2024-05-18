package mysql_test

import (
	"context"
	"errors"
	"testing"
	"user-svc/adapters/repository/db"
	"user-svc/adapters/repository/mysql"
	"user-svc/internal/entity"
	mocks "user-svc/ports/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMysqlDB_CreateUser_Success(t *testing.T) {
	// Setup mocks
	// mockMetrics := mocks.NewMockDatabaseMetrics()
	mockConfig := mocks.NewMockConfig()
	mockLogger := mocks.NewMockLogger()

	type testCase struct {
		name          string
		mockExecError error
		expectedUser  entity.User
	}

	cases := []testCase{
		{
			name:          "Successful user creation",
			mockExecError: nil,
			expectedUser: entity.User{
				ID:          1,
				PhoneNumber: "123456789",
				Email:       "test@example.com",
				Password:    "password",
			},
		},
		{
			name:          "Error executing command",
			mockExecError: errors.New("error executing command"),
			expectedUser:  entity.User{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.MockDB)
			mockMetrics := mocks.NewMockDatabaseMetrics()

			// TODO: connection
			// mysqlDB := mysql.MysqlDB{db: mockDB}
			dbPool, _ := db.GetConnectionPool(mockConfig)
			mysqlDB := mysql.New(dbPool, mockLogger, mockMetrics)

			user := entity.User{
				PhoneNumber: tt.expectedUser.PhoneNumber,
				Email:       tt.expectedUser.Email,
				Password:    tt.expectedUser.Password,
			}

			ctx := context.Background()

			//? explain
			mockResult := &MockResult{LastInsertIdValue: 1, RowsAffectedValue: 1}
			mockDB.On("ExecContext", mock.Anything, mock.Anything, user.PhoneNumber, user.Email, user.Password).
				Return(mockResult, tt.mockExecError)

			//? explain
			mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Once()

			createdUser, err := mysqlDB.CreateUser(ctx, user)

			mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Once()

			if tt.mockExecError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, createdUser)
			}

			mockDB.AssertExpectations(t)
		})
	}

}

type MockResult struct {
	LastInsertIdValue int64
	RowsAffectedValue int64
}

func (m *MockResult) LastInsertId() (int64, error) {
	return m.LastInsertIdValue, nil
}

func (m *MockResult) RowsAffected() (int64, error) {
	return m.RowsAffectedValue, nil
}
