package httpServer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-svc/adapters/config"
	"user-svc/adapters/logger"
	"user-svc/adapters/messaging"
	"user-svc/internal/entity"
	userService "user-svc/internal/service"
	"user-svc/internal/service/param"
	mocks "user-svc/ports/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// !---
func TestRegisterHandler(t *testing.T) {
	type testCase struct {
		name           string
		reqBody        interface{}
		expectedStatus int
		expectedError  string
	}

	cases := []testCase{
		{
			name:           "Successful Registration",
			reqBody:        param.RegisterRequest{PhoneNumber: "1234567890", Email: "test@example.com", Password: "password123"},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name:           "Failed Request Body Binding",
			reqBody:        "invalid", // Invalid request body to trigger binding error
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Failed to bind request",
		},
		// Add more test cases as needed
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			//* Arrange
			configAdapter, _ := config.NewViperAdapter()
			// err := configAdapter.LoadConfig("../../../config.yaml")
			// if err != nil {
			// 	panic("failed to load configuration")
			// }

			userRepoMock := mocks.NewMockUserRepository()
			userRepoMock.On("CreateUser").Return(entity.User{
				ID:          1,
				PhoneNumber: "1234567890",
				Email:       "test@example.com",
				Password:    "password123",
			}, nil)

			zapLogger, _ := logger.NewZapLogger(configAdapter)

			//connectionString := "amqp://puppet:password@localhost:5672/"
			publisher, err := messaging.NewRabbitMQClient(configAdapter)
			if err != nil {
				log.Fatalf("failed to connect to RabbitMQ: %v", err)
			}

			//prometheusAdapter := metrics.NewPrometheus()
			userSvc := userService.NewService(configAdapter, userRepoMock, publisher, zapLogger)

			userHandler := New(configAdapter, userSvc, zapLogger)

			//* http://localhost:5000/register
			e := echo.New()
			e.POST("/register", userHandler.Register)

			//* Act
			reqBody, _ := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))

			// Create a new HTTP response recorder
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			resp := rec.Result()
			defer resp.Body.Close()

			// Set up expectations for the logger mock
			// mockLogger.On("Info", "Handling register request").Once()
			// mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

			//* Assertions
			assert.NoError(t, err)

			assert.Equal(t, c.expectedStatus, rec.Code)

			// Assert that the expected method calls were made
			// mockLogger.AssertExpectations(t)
			// mockUserSvc.AssertExpectations(t)
		})
	}
}
