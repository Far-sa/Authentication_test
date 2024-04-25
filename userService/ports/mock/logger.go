package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLogger struct {
	mock.Mock
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
	// You might want to assert on message and fields here (optional)
}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

// MockLogger is a mock implementation of Logger for testing
// type MockLogger struct {
// 	InfoFn  func(msg string, fields ...zap.Field)
// 	DebugFn func(msg string, fields ...zap.Field)
// 	ErrorFn func(msg string, fields ...zap.Field)
// 	WarnFn  func(msg string, fields ...zap.Field)
// }

// func (m *MockLogger) Info(msg string, fields ...zap.Field) {
// 	if m.InfoFn != nil {
// 		m.InfoFn(msg, fields...)
// 	}
// }

// func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
// 	if m.DebugFn != nil {
// 		m.DebugFn(msg, fields...)
// 	}
// }

// func (m *MockLogger) Error(msg string, fields ...zap.Field) {
// 	if m.ErrorFn != nil {
// 		m.ErrorFn(msg, fields...)
// 	}
// }

// func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
// 	if m.WarnFn != nil {
// 		m.WarnFn(msg, fields...)
// 	}
// }

// Add mocks for other logging methods as needed...
