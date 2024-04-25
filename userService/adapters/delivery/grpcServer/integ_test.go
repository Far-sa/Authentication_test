package grpcserver

import (
	"context"
	"log"
	"testing"
	user "user-svc/ports/protobuf/grpc/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserIntegration(t *testing.T) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := user.NewUserServiceClient(conn)

	// Prepare a CreateUserRequest
	req := &user.CreateUserRequest{
		Email:       "test@example.com",
		Password:    "password123",
		PhoneNumber: "1234567890",
	}

	// Call CreateUser method on the gRPC server
	resp, err := client.CreateUser(context.Background(), req)
	if err != nil {
		// Handle gRPC errors
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("Failed to create user: %v", err)
		}
		t.Fatalf("Failed to create user: %s", st.Message())
	}

	// Verify the response
	assert.True(t, resp.Success)
	assert.Equal(t, req.Email, resp.Email)
}
