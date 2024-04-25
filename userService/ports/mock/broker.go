package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockEventPublisher struct {
	mock.Mock
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

func (m *MockEventPublisher) PublishUserRegisteredEvent(ctx context.Context, data []byte) error {
	args := m.Called(ctx, data)
	return args.Error(1)
}
