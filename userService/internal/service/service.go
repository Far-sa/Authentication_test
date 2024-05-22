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

	//? publish event if necessary	 in case of  scenario changed
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

func (s *userService) GetUserProfile(ctx context.Context, userID uint) (param.UserInfo, error) {
	userInfo, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return param.UserInfo{}, err
	}
	return param.UserInfo{ID: userInfo.ID}, nil
}

func (s *userService) StartMessageListener(ctx context.Context) error {
	//TODO need to create exchange,queue,binding first

	log.Println("Starting message listener...")

	if err := s.messageBroker.DeclareExchange("auth_exchange", "direct"); err != nil {
		log.Printf("Failed to declare exchange: %v", err)
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	log.Println("Exchange declared.")

	queue, err := s.messageBroker.CreateQueue("login_requests", true, false)
	if err != nil {
		log.Printf("Failed to create queue: %v", err)
		return fmt.Errorf("failed to create queue: %w", err) // Propagate error
	}
	log.Println("Queue created.")

	if err := s.messageBroker.CreateBinding(queue.Name, "login", "auth_exchange"); err != nil {
		log.Printf("Failed to bind queue: %v", err)
		return fmt.Errorf("failed to bind queue: %w", err) // Propagate error
	}
	log.Println("Queue bound to exchange.")

	msgs, err := s.messageBroker.Consume(queue.Name, "user_service", false)
	if err != nil {
		log.Printf("Failed to start message listener: %v", err)
		return fmt.Errorf("failed to start message listener: %w", err)
	}

	// Start consuming messages in a goroutine
	go func() {
		log.Println("Consumer started inside go routine, waiting for messages...")

		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			// Process the received message
			err := s.processMessage(ctx, d)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
			}
		}
		log.Println("Exited message consumption loop")
	}()

	log.Printf("Waiting for login requests. To exit press CTRL+C")
	return nil
}

func (s *userService) processMessage(ctx context.Context, d amqp.Delivery) error {
	// Unmarshal data
	var loginReq param.LoginRequest
	err := json.Unmarshal(d.Body, &loginReq)
	if err != nil {
		log.Printf("Failed to unmarshal login request: %v", err)
		return fmt.Errorf("failed to unmarshal login request: %w", err)
	}
	log.Printf("Processing login request for email: %s", loginReq.Email)

	userExists, err := s.CheckUserInDatabase(ctx, loginReq)
	if err != nil {
		log.Printf("Error checking user in database: %v", err)
		return fmt.Errorf("checking user in database failed : %w", err)
	}

	response := param.LoginResponse{}
	// Assuming UserExists is a field in param.LoginResponse indicating user existence
	if userExists.UserExist {
		log.Printf("User found: %s", userExists.User.Email)
		response.UserExist = true
		response.User = userExists.User
		response.Error = "user found"
	} else {
		// User not found
		log.Printf("User not found: %s", loginReq.Email)
		response.UserExist = false
		response.Error = "user not found: please register"
	}

	// Publish to user-service queue
	err = s.publishUserData(ctx, response)
	if err != nil {
		log.Printf("Failed to publish user response: %v", err)
		return fmt.Errorf("failed to publish response: %w", err)
	}

	if err := d.Ack(false); err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}

	log.Printf("Message processed and acknowledged for email: %s", loginReq.Email)
	return nil
}

func (s *userService) CheckUserInDatabase(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error) {
	// Implement logic to check user data based on email using injected UserRepository
	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User not found in database: %s", user.Email)
			return param.LoginResponse{Error: "user not found: please register"}, fmt.Errorf("user not found: %w", err)
		}
		log.Printf("Error fetching user data: %v", err)
		return param.LoginResponse{}, fmt.Errorf("error fetching user data: %w", err)
	}

	log.Printf("User found in database: %s", existingUser.Email)
	return param.LoginResponse{UserExist: true, User: param.UserInfo{ID: existingUser.ID, Email: existingUser.Email}}, nil
}

func (s *userService) publishUserData(ctx context.Context, userData interface{}) error {

	// Log the userData before serialization
	log.Printf("Publishing user data: %+v", userData)

	data, jErr := json.Marshal(userData)
	if jErr != nil {
		return fmt.Errorf("failed to serialize user data: %w", jErr)
	}

	// Log the serialized data and its length
	log.Printf("Serialized user data: %s", data)
	log.Printf("Serialized data length: %d", len(data))

	// Ensure the data is not empty
	if len(data) == 0 {
		return fmt.Errorf("serialized user data is empty")
	}

	if err := s.messageBroker.Publish(ctx, "auth_exchange", "user_response", amqp.Publishing{
		ContentType:   "application/json",
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
