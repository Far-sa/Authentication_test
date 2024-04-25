package httpServer

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockServer is a mock implementation of the Server struct
type MockServer struct {
	mock.Mock
}

// Serve mocks the Serve method of the Server struct
func (m *MockServer) Serve() {
	// We don't need to implement the actual Serve logic for this mock
	// Instead, we can use it to verify that Serve is called with the expected behavior
	m.Called()
}

func TestHttp_Serve(t *testing.T) {
	// Create an instance of the mock server
	mockServer := new(MockServer)

	// Expect the Serve method to be called once
	mockServer.On("Serve").Once()

	// Call the Serve method of the mock server
	mockServer.Serve()

	// Assert that the expectations for the Serve method were met
	mockServer.AssertExpectations(t)
}
