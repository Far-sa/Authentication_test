package userService

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"user-svc/internal/entity"
	"user-svc/internal/service/param"
	"user-svc/ports"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	Config        ports.Config
	userRepo      ports.UserRepository
	messageBroker ports.EventPublisher
	logger        ports.Logger
}

func NewService(cfg ports.Config, repo ports.UserRepository, publisher ports.EventPublisher, logger ports.Logger) *userService {
	return &userService{Config: cfg, userRepo: repo, messageBroker: publisher, logger: logger}
}

func (s *userService) Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error) {

	// hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	hashed, err := HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Error generating hashed password", zap.Error(err))
		return param.RegisterResponse{}, fmt.Errorf("error generating hashed password: %v", err)
	}

	user := entity.User{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    string(hashed),
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("Error creating user", zap.Error(err))
		return param.RegisterResponse{}, fmt.Errorf("error creating user: %v", err)
	}

	//? publish event
	// Call separate function to publish user data
	// if err := us.publishUserData(ctx, createdUser); err != nil {
	// 	us.logger.Error("Failed to publish user credential to auth-service", zap.Error(err))
	// 	return param.RegisterResponse{}, fmt.Errorf("error publishing user data: %w", err)
	// }

	s.logger.Info("User created successfully", zap.Any("user", createdUser))

	return param.RegisterResponse{
		User: param.UserInfo{
			ID:          createdUser.ID,
			Email:       createdUser.Email,
			PhoneNumber: createdUser.PhoneNumber,
		},
	}, nil
}

func (s *userService) StartMessageListener(ctx context.Context) error {
	//TODO implement in auth svc
	msgs, err := s.messageBroker.Consume("login_requests", "user_service", false)
	if err != nil {
		return fmt.Errorf("failed to start message listener: %w", err)
	}
	go func() {
		for d := range msgs {
			err := s.processMessage(ctx, d)
			if err != nil {
				log.Printf("failed to process message: %v", err)
			}
		}
	}()

	log.Printf("Waiting for login requests. To exit press CTRL+C")
	return nil
}

func (s *userService) processMessage(ctx context.Context, d amqp.Delivery) error {
	// Unmarshal data
	var loginReq param.LoginRequest
	err := json.Unmarshal(d.Body, &loginReq)
	if err != nil {
		return fmt.Errorf("failed to unmarshal login request: %w", err)
	}

	//! Check data with database
	userResponse, err := s.CheckUserInDatabase(ctx, loginReq)
	if err != nil {
		return fmt.Errorf("checking user in database failed : %w", err)
	}

	response := param.LoginResponse{}
	// Assuming UserExists is a field in param.LoginResponse indicating user existence
	if userResponse.UserExist {
		// User exists, process the response
		response.Error = "user found"
		response.UserExist = true
		response.User = param.UserInfo{ID: response.User.ID, Email: response.User.Email}
	} else {
		// User not found
		response.Error = "user not found: please register"
		response.UserExist = false
	}

	// Publish to user-service queue
	err = s.publishUserData(ctx, response)
	if err != nil {
		return fmt.Errorf("failed to publish response: %w", err)
	}

	return nil
}

func (s *userService) CheckUserInDatabase(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error) {
	// Implement logic to check user data based on email using injected UserRepository
	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		// Handle potential errors from userRepo.GetUserByEmail
		if errors.Is(err, sql.ErrNoRows) { // Handle user not found error
			return param.LoginResponse{}, fmt.Errorf("user not found")
		}
		return param.LoginResponse{}, fmt.Errorf("error fetching user data: %w", err)
	}

	// Additional checks or data manipulation on the existing user object (optional)

	return param.LoginResponse{UserExist: true, User: param.UserInfo{ID: existingUser.ID, Email: existingUser.Email}}, nil
}

func (s *userService) publishUserData(ctx context.Context, userData interface{}) error {

	if err := s.messageBroker.DeclareExchange("user_events", "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := s.messageBroker.CreateQueue("user_info", true, false)
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err) // Propagate error
	}

	if err := s.messageBroker.CreateBinding(queue.Name, "users.*", "user_events"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err) // Propagate error
	}

	data, jErr := json.Marshal(userData)
	if jErr != nil {
		return fmt.Errorf("failed to serialize user data: %w", jErr)
	}

	if err := s.messageBroker.Publish(ctx, "user_events", "users.new", amqp.Publishing{
		ContentType:   "text/plain",
		DeliveryMode:  amqp.Persistent,
		Body:          data,
		CorrelationId: uuid.NewString(),
	}); err != nil {
		return fmt.Errorf("failed to publish user to auth-service: %w", err)
	}
	log.Println("User event published successfully")

	return nil
}

// ! Helper functions

func UnmarshalUser(data []byte) (entity.User, error) {
	var user entity.User
	err := json.Unmarshal(data, &user)
	return user, err
}

func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func mapUintToByte(num uint) byte {
	// Since byte is an alias for uint8, we can directly cast uint to byte
	return byte(num)
}

func MapStringToByte(str string) []byte {
	// Convert string to byte slice ([]byte)
	return []byte(str)
}
func mapByteToString(bytes []byte) string {
	// Convert byte slice ([]byte) to string
	return string(bytes)
}

//!!!!!>

//!!-- imp
// func main() {
// 	rabbitMQ, err := NewRabbitMQ("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		log.Fatalf("failed to connect to RabbitMQ: %v", err)
// 	}
// 	defer rabbitMQ.Close()

// 	userSvc := NewUserService(rabbitMQ)

// 	// Start the message listener
// 	err = userSvc.startMessageListener()
// 	if err != nil {
// 		log.Fatalf("failed to start message listener: %v", err)
// 	}

// 	// Set up HTTP server and handlers
// 	http.HandleFunc("/register", userSvc.registerUserHandler)
// 	http.HandleFunc("/update", userSvc.updateUserHandler)
// 	http.HandleFunc("/delete", userSvc.deleteUserHandler)
// 	http.HandleFunc("/info", userSvc.getUserInfoHandler)

// 	log.Println("Starting HTTP server on :8080")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatalf("failed to start HTTP server: %v", err)
// 	}
// }
