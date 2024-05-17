package delivery_test

import (
	"auth-svc/adapters/config"
	"auth-svc/adapters/delivery"
	"auth-svc/adapters/messaging"
	"auth-svc/internal/entity"
	"auth-svc/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) StoreToken(userID uint, token string, expiration time.Time) error {
	args := m.Called(userID, token, expiration)
	return args.Error(0)
}

func (m *mockRepository) RetrieveToken(userID uint) (*entity.Token, error) {
	args := m.Called(userID)
	return args.Get(0).(*entity.Token), args.Error(1)
}

func TestLoginHandleIntegration(t *testing.T) {

	type testCase struct {
		name       string
		reqBody    []byte
		statusCode int
		respBody   string
	}

	cases := []testCase{
		{name: "Success",
			reqBody:    []byte(`{"username": "johndoe", "password": "secret"}`),
			statusCode: http.StatusOK,
			respBody:   `{"token": "valid_token"}`,
		},
		{
			name:       "Invalid Request Body (Missing Username)",
			reqBody:    []byte(`{"password": "secret"}`),
			statusCode: http.StatusBadRequest,
			respBody:   "", // No specific expected res
		},
		{
			name:       "AuthService Error",
			reqBody:    []byte(`{"username": "johndoe", "password": "secret"}`),
			statusCode: http.StatusInternalServerError, // Or specific error code returned by AuthService
			respBody:   "",                             // No specific expe
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			mockRepo := new(mockRepository)
			mockRepo.On("StoreToken", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(nil)

			cfg, _ := config.NewViperAdapter()
			event, _ := messaging.NewRabbitMQClient(cfg)
			svc := service.NewAuthService(cfg, mockRepo, event)
			handler := delivery.NewAuthHandler(svc)

			e := echo.New()
			e.POST("login", handler.Login)

			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			assert.NoError(t, err, "Failed to create HTTP request")
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			resp := rec.Result()
			defer resp.Body.Close()

			// Verify response status code
			assert.Equal(t, http.StatusOK, rec.Code, "Unexpected response status code")

			// Assert response body
			actualBody := rec.Body.String()
			assert.Equal(t, c.respBody, actualBody, "Unexpected response body")

			// Verify mock expectations
			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}

}
