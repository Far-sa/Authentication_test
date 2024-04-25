package mocks

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

// ExecContext mocks the ExecContext method of sqlx.DB
func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// Retrieve the expected result and error based on the input query
	result, err := m.Called(ctx, query, args).Get(0).(sql.Result), m.Called(ctx, query, args).Error(1)
	return result, err
}

// Close mocks the Close method of sqlx.DB
func (m *MockDB) Close() error {
	// You can implement this method if your code under test calls Close on the DB
	// For now, we don't need to mock it
	return nil
}
