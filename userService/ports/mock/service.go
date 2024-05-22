package mocks

import (
	"context"
	"user-svc/internal/service/param"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func NewMockService() *MockService {
	return &MockService{}
}

func (m *MockService) Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(param.RegisterResponse), args.Error(1)
}

func (m *MockService) GetUserProfile(ctx context.Context, userID uint) (param.UserInfo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(param.UserInfo), args.Error(1)
}
