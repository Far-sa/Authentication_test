package grpcserver

import (
	"context"
	"testing"
	user "user-svc/ports/protobuf/grpc/user"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ! MockServer is a mock implementation of ServerInterface
type MockServer struct {
	mock.Mock
}

func (m *MockServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*user.CreateUserResponse), args.Error(1)
}

// ! Test with mocking interface
func TestCreateUser(t *testing.T) {
	// Create a new instance of the mock server
	mockServer := new(MockServer)

	// Setup expectations for the mock server's CreateUser method
	mockServer.On("CreateUser", mock.Anything, mock.Anything).Return(&user.CreateUserResponse{Email: "fesgheli@teo.com", Success: true}, nil)

	// Perform the test using the mock server
	res, err := mockServer.CreateUser(context.Background(), &user.CreateUserRequest{Email: "fesgheli@teo.com"})

	// Assert that the response matches the expected values
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "fesgheli@teo.com", res.Email)
	require.True(t, res.Success)

	// Assert that the expectations for the mock server's CreateUser method were met
	mockServer.AssertExpectations(t)
}

// ! real case scenario testing
// func TestServer_CreateUser(t *testing.T) {

// 	lis := bufconn.Listen(1024 * 1024)
// 	t.Cleanup(func() {
// 		lis.Close()
// 	})

// 	srv := grpc.NewServer()
// 	t.Cleanup(func() {
// 		srv.Stop()
// 	})

// 	//? learn : This could happen if the server struct is an unexported type, and you're trying to instantiate it directly in your test function.
// 	svc := server{}
// 	user.RegisterUserServiceServer(srv, &svc)

// 	go func() {
// 		err := srv.Serve(lis)
// 		assert.NoError(t, err, "failed to serve")

// 		// if err := srv.Serve(lis); err != nil {
// 		// 	log.Fatalf("service server error: %v", err)
// 		// }
// 	}()

// 	//* Test
// 	dialer := func(context.Context, string) (net.Conn, error) {
// 		return lis.Dial()
// 	}

// 	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
// 	require.NoError(t, err) // Ensure no dialing error

// 	t.Cleanup(func() {
// 		conn.Close()
// 	})

// 	client := user.NewUserServiceClient(conn)
// 	res, err := client.CreateUser(context.Background(), &user.CreateUserRequest{Email: "fesgheli@teo.com"})
// 	require.NoError(t, err) // Ensure no RPC error

// 	require.Equal(t, "fesgheli@teo.com", res.Email)
// 	require.True(t, res.Success)

// }
