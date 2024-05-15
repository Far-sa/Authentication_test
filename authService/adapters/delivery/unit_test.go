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
