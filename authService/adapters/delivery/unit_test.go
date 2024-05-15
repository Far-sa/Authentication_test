package delivery_test

import (
	"auth-svc/adapters/delivery"
	"auth-svc/internal/param"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(param.LoginResponse), args.Error(1)
}

func TestLoginHandler(t *testing.T) {

	mockAuthSvc := new(mockAuthService)

	// Define expected behavior for AuthService.Login
	expectedResp := param.LoginResponse{ /* fill with appropriate values */ }
	mockAuthSvc.On("Login", mock.Anything, mock.AnythingOfType("param.LoginRequest")).Return(expectedResp, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := delivery.NewAuthHandler(mockAuthSvc)
	err := handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestLoginHandler_TableDriven(t *testing.T) {
	type testCase struct {
		name         string
		requestBody  interface{}
		expectedCode int
		expectedResp param.LoginResponse
	}

	cases := []testCase{
		{
			name:         "Valid Login Request",
			requestBody:  map[string]string{"username": "test_user", "password": "test_password"},
			expectedCode: http.StatusOK,
			expectedResp: param.LoginResponse{User: param.UserInfo{ID: 1, PhoneNumber: "", Email: ""}, Tokens: param.Tokens{AccessToken: "", RefreshToken: ""}},
		},
		{
			name:         "Invalid JSON Request Body",
			requestBody:  "{", // Invalid JSON
			expectedCode: http.StatusBadRequest,
			expectedResp: param.LoginResponse{},
		},
		{
			name:         "Missing Username",
			requestBody:  map[string]string{"password": "test_password"},
			expectedCode: http.StatusBadRequest,
			expectedResp: param.LoginResponse{},
		},
		{
			name:         "AuthService Login Error",
			requestBody:  map[string]string{"username": "test_user", "password": "test_password"},
			expectedCode: http.StatusBadRequest,
			expectedResp: param.LoginResponse{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// define mock auth service
			mockAuthSvc := new(mockAuthService)

			// Define expected behavior for AuthService.Login
			mockAuthSvc.On("Login", mock.Anything, mock.AnythingOfType("param.LoginRequest")).Return(c.expectedResp, nil)

			//TODO: complete assertion
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/login", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := delivery.NewAuthHandler(mockAuthSvc)
			err := handler.Login(c)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}
