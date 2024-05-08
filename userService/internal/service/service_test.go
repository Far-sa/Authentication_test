package userService_test

import (
	"context"
	"fmt"
	"testing"
	"user-svc/internal/entity"
	userService "user-svc/internal/service"
	"user-svc/internal/service/param"
	mocks "user-svc/ports/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// !-----

func TestRegisterUser(t *testing.T) {

	//TODO can separate test cases from error cases

	type testCase struct {
		name          string
		req           param.RegisterRequest
		expectedUser  entity.User
		userRepoError error
		eventPubError error
		expectedError string
	}

	mockUserRepo := mocks.NewMockUserRepository()
	mockEventPublisher := mocks.NewMockEventPublisher()
	mockLogger := mocks.NewMockLogger()
	mockConfig := mocks.NewMockConfig()

	cases := []testCase{
		{
			name: "Successful Registration",
			req: param.RegisterRequest{
				PhoneNumber: "1234567890",
				Email:       "test@example.com",
				Password:    "password123",
			},
			expectedUser: entity.User{
				ID:          1,
				PhoneNumber: "1234567890",
				Email:       "test@example.com",
				Password:    "hashedpassword",
			},
			userRepoError: nil,
			eventPubError: nil,
			expectedError: "",
		},
		{
			name: "Failure in User Creation",
			req: param.RegisterRequest{
				PhoneNumber: "9876543210",
				Email:       "failure@example.com",
				Password:    "password456",
			},
			expectedUser:  entity.User{},
			userRepoError: fmt.Errorf("error creating user"),
			eventPubError: nil,
			expectedError: "error creating user",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			//* Arrange
			ctx := context.Background()

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(c.req.Password), 8)
			c.expectedUser.Password = string(hashedPassword)

			mockUserRepo.On("CreateUser", ctx, mock.Anything).Return(c.expectedUser, c.userRepoError)
			if c.userRepoError == nil {
				eventPayload := []byte(userService.MapStringToByte(c.expectedUser.Email))
				mockEventPublisher.On("PublishUserRegisteredEvent", ctx, eventPayload).Return(c.eventPubError)
			}

			mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Once()

			service := userService.NewService(mockConfig, mockUserRepo, mockEventPublisher, mockLogger)

			//* ACT
			resp, err := service.Register(ctx, c.req)

			//* Assert
			assert.NoError(t, err)
			assert.Equal(t, c.expectedUser.ID, resp.User.ID)
			assert.Equal(t, c.expectedUser.Email, resp.User.Email)
			assert.Equal(t, c.expectedUser.PhoneNumber, resp.User.PhoneNumber)

			// Assert that the expected method calls were made
			// mockUserRepo.AssertExpectations(t)
			// mockEventPublisher.AssertExpectations(t)
			// mockLogger.AssertExpectations(t)
		})
	}
}
