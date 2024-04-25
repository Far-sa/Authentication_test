package grpcserver

import (
	"context"
	"log"
	"net"
	"user-svc/internal/service/param"
	"user-svc/ports"
	user "user-svc/ports/protobuf/grpc/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	user.UnimplementedUserServiceServer
	userSvc ports.Service
	//server  *grpc.server

	//config  ports.Config
}

func New(userSvc ports.Service) server {
	return server{
		UnimplementedUserServiceServer: user.UnimplementedUserServiceServer{},
		userSvc:                        userSvc,
	}
}

// TODO gracefully shutdown

func (s server) Start() {

	// port := s.config.GetHTTPConfig().Port
	// address := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// grpc server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// register user service to  grpc server
	user.RegisterUserServiceServer(grpcServer, &s)

	// serve grpc server
	log.Println("User gRPC server started on port :8000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

//TODO: implement methods

// * CreateUser creates a new user via gRPC
func (s *server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {

	// Deserialize request
	email := req.GetEmail()
	password := req.GetPassword()
	phonNumber := req.GetPhoneNumber()

	createdUser, err := s.userSvc.Register(ctx, param.RegisterRequest{PhoneNumber: phonNumber, Email: email, Password: password})

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return &user.CreateUserResponse{Success: false},
			status.Errorf(codes.FailedPrecondition, "failed to create user %v:", req.Email)
	}

	// Serialize response
	return &user.CreateUserResponse{Success: true, Email: createdUser.User.Email}, nil
}
