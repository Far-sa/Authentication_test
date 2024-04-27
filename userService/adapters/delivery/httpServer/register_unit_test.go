package httpServer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"user-svc/internal/service/param"
	mocks "user-svc/ports/mock"

	"github.com/labstack/echo/v4"
	//"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO modify mockLogger methods
func TestRegister_Success(t *testing.T) {
	// Create mock objects
	mockLogger := mocks.NewMockLogger()
	if mockLogger == nil {
		t.Fatal("mockLogger is nil")
	}

	mockMetrics := mocks.NewMockHTTPMetrics()
	if mockMetrics == nil {
		t.Fatal("mockMetrics is nil")
	}

	mockService := mocks.NewMockService()
	if mockService == nil {
		t.Fatal("mockService is nil")
	}
	//mockConfig := mocks.NewMockConfig()

	// Setup expectations
	mockLogger.On("Info", "Handling register request", mock.Anything).Once()

	// mockMetrics.On("RegisterHTTPDurationHistogram").Return(&prometheus.HistogramVec{})
	// mockMetrics.On("RegisterHTTPErrorCounter").Return(&prometheus.CounterVec{})
	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Once()

	//mockConfig.On("").Once()
	mockService.On("Register", mock.Anything, mock.Anything).Return(
		param.RegisterResponse{User: param.UserInfo{Email: ""}}, nil).Once()

	// Create a new instance of the Server struct with mocked dependencies
	server := server{
		logger: mockLogger,
		// metrics: mockMetrics,
		userSvc: mockService,
		Router:  echo.New(),
	}

	// Create a new echo.Context (you may need to set up additional dependencies for the context)
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email": "test@example.com"}`))
	rec := httptest.NewRecorder()
	c := server.Router.NewContext(req, rec)

	err := server.Register(c)

	mockLogger.On("Info", " register success", mock.Anything).Once()

	// Assert that no error occurred
	assert.NoError(t, err)

	// Verify that all expectations were met
	mockLogger.AssertExpectations(t)
	mockMetrics.AssertExpectations(t)
	mockService.AssertExpectations(t)
}
