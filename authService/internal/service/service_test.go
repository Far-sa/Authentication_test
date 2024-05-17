package service_test

import (
	"auth-svc/internal/entity"
	"auth-svc/internal/param"
	"auth-svc/internal/ports"
	"auth-svc/internal/service"
	"context"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) ConsumeMessages() (<-chan amqp.Delivery, error) {
	args := m.Called()
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)
}

// func (s mockAuthService) consumeMessages() (chan interface{}, error) {
//     userChan := make(chan interface{})
//     go func() {
//         user := User{ID: 1, Email: "test@example.com", Password: "password123"}
//         userChan <- user
//     }()
//     return userChan, nil
// }

func (m *mockAuthService) ComparePassword(hashPass, reqPass string) (bool, error) {
	args := m.Called(hashPass, reqPass)
	return args.Get(0).(bool), args.Error(1)
}

// CreateAccessToken mocks the createAccessToken method
func (m *mockAuthService) CreateAccessToken(user service.User) (string, error) {
	args := m.Called(user)
	return args.Get(0).(string), args.Error(1)
}

// RefreshAccessToken mocks the refreshAccessToken method
func (m *mockAuthService) RefreshAccessToken(user service.User) (string, error) {
	args := m.Called(user)
	return args.Get(0).(string), args.Error(1)
}

// StoreToken mocks the StoreToken method
func (m *mockAuthService) StoreToken(userID uint, accessToken string, expiry time.Time) error {
	args := m.Called(userID, accessToken, expiry)
	return args.Error(0)
}

func TestLoginServiceLogic(t *testing.T) {
	// Arrange
	mockAuthSvc := new(mockAuthService)
	expectedPassword := "correctPassword"
	expectedUser := service.User{ID: 1, Email: "test@example.com", Password: "hashedPassword"}

	userChan := make(chan interface{})
	userChan <- expectedUser

	mockAuthSvc.On("ConsumeMessage").Return(userChan, nil)
	mockAuthSvc.On("ComparePassword", expectedUser.Password, expectedPassword).Return(true, nil)
	mockAuthSvc.On("CreateAccessToken", expectedUser).Return("valid_access_token", nil)
	mockAuthSvc.On("RefreshAccessToken", expectedUser).Return("valid_refresh_token", nil)
	mockAuthSvc.On("StoreToken", expectedUser.ID, "valid_access_token", mock.Anything).Return(nil)

	mockCfg := &mockConfig{}
	mockRepo := &mockRepo{}
	mockEvent := &mockEventPublisher{}

	ctxWithCan, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	loginReq := param.LoginRequest{Email: expectedUser.Email, Password: expectedPassword}

	// Act
	s := service.NewAuthService(mockCfg, mockRepo, mockEvent)
	response, err := s.Login(ctxWithCan, loginReq)

	// Assert
	assert.NoError(t, err, "Login should succeed")
	assert.Equal(t, "valid_access_token", response.Tokens.AccessToken)
	assert.Equal(t, "valid_refresh_token", response.Tokens.RefreshToken)
	assert.Equal(t, uint(expectedUser.ID), response.User.ID)
	assert.Equal(t, expectedUser.Email, response.User.Email)

	// Verify mock expectations
	// 	assert.HasCalled(t, mockAuthService.On("ConsumeMessages"), "ConsumeMessages should be called")
	// 	assert.HasCalled(t, mockAuthService.On("comparePassword", expectedUser.Password, expectedPassword), "comparePassword should be called with expected arguments")
	// 	assert.HasCalled(t, mockAuthService.On("CreateAccessToken", expectedUser), "CreateAccessToken should be called with expected user")
	// 	assert.HasCalled(t, mockAuthService.On("RefreshAccessToken", expectedUser), "RefreshAccessToken should be called with expected user")
	// 	assert.HasCalled(t, mockAuthService.On("StoreToken", expectedUser.ID, "valid_access_token", mock.Anything), "StoreToken should be called with expected arguments")

}

type mockConfig struct{ mock.Mock }
type mockRepo struct{ mock.Mock }
type mockEventPublisher struct{ mock.Mock }

func (m *mockConfig) GetBrokerConfig() ports.BrokerConfig {
	args := m.Called()
	return args.Get(0).(ports.BrokerConfig)
}

func (m *mockConfig) GetDatabaseConfig() ports.DatabaseConfig {
	args := m.Called()
	return args.Get(0).(ports.DatabaseConfig)
}

func (m *mockConfig) GetHTTPConfig() ports.HTTPConfig {
	args := m.Called()
	return args.Get(0).(ports.HTTPConfig)
}

func (m *mockRepo) RetrieveToken(userId uint) (*entity.Token, error) {
	args := m.Called(userId)

	if token, ok := args.Get(0).(*entity.Token); ok {
		return token, nil
	}

	return nil, args.Error(1)
}

func (m *mockRepo) StoreToken(userID uint, token string, expiration time.Time) error {
	args := m.Called(userID, token, expiration)
	return args.Error(0)
}

func (m *mockEventPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockEventPublisher) Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	args := m.Called(queueName, consumer, autoAck)
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)

}
