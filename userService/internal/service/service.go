package userService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"user-svc/internal/entity"
	"user-svc/internal/service/param"
	"user-svc/ports"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// TODO: add config - zap logger as singleton

type Config struct{}

// Service represents the usecase for user operations
type Service struct {
	Config         ports.Config
	userRepo       ports.UserRepository
	eventPublisher ports.EventPublisher
	logger         ports.Logger
	// Cash           ports.Caching
}

// NewService creates a new instance of Service
func NewService(cfg ports.Config, repo ports.UserRepository, publisher ports.EventPublisher, logger ports.Logger) Service {
	return Service{Config: cfg, userRepo: repo, eventPublisher: publisher, logger: logger}
}

// RegisterUser handles user registration
func (us Service) Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error) {

	// hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	hashed, err := HashPassword(req.Password)
	if err != nil {
		us.logger.Error("Error generating hashed password", zap.Error(err))
		return param.RegisterResponse{}, fmt.Errorf("error generating hashed password: %v", err)
	}

	user := entity.User{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    string(hashed),
	}

	createdUser, err := us.userRepo.CreateUser(ctx, user)
	if err != nil {
		us.logger.Error("Error creating user", zap.Error(err))
		return param.RegisterResponse{}, fmt.Errorf("error creating user: %v", err)
	}

	//? publish event
	// Call separate function to publish user data
	if err := us.publishUserData(ctx, createdUser); err != nil {
		us.logger.Error("Failed to publish user credential to auth-service", zap.Error(err))
		return param.RegisterResponse{}, fmt.Errorf("error publishing user data: %w", err)
	}

	us.logger.Info("User created successfully", zap.Any("user", createdUser))

	return param.RegisterResponse{
		User: param.UserInfo{
			ID:          createdUser.ID,
			Email:       createdUser.Email,
			PhoneNumber: createdUser.PhoneNumber,
		},
	}, nil
}

func (us Service) publishUserData(ctx context.Context, createdUser interface{}) error {

	if err := us.eventPublisher.DeclareExchange("user_events", "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := us.eventPublisher.CreateQueue("registration_queue", true, false)
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err) // Propagate error
	}

	if err := us.eventPublisher.CreateBinding(queue.Name, "registration.*", "user_events"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err) // Propagate error
	}

	data, jErr := json.Marshal(createdUser)
	if jErr != nil {
		return fmt.Errorf("failed to serialize user data: %w", jErr)
	}

	if err := us.eventPublisher.Publish(ctx, "user_events", "registration.new", amqp.Publishing{
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

func (us Service) CheckUserExistence(ctx context.Context) (param.RegisterResponse, error) {
	userChan, err := us.consumeMessages()
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("failed to consume messages: %w", err)
	}
	_ = userChan
	//TODO add logic to compare data between database and req
	select {
	case <-ctx.Done():
		// Handle timeout or context cancellation
		return param.RegisterResponse{}, errors.New("timed out waiting for user information")
	case data, ok := <-userChan:
		if !ok {
			return param.RegisterResponse{}, errors.New("user not found") // No matching user found
		}
		user, ok := data.(entity.User) // Cast data to the User struct
		if !ok {
			return param.RegisterResponse{}, errors.New("invalid data received from queue")
		}
		// if err != nil {
		// 	return param.RegisterResponse{}, fmt.Errorf("failed to compare password: %w", err)
		// }

		// if !valid {
		// 	return param.RegisterResponse{}, errors.New("invalid password") // Indicate invalid password
		// }

		if err := us.userRepo.CheckUserData(user.ID, accessToken, time.Now().Add(72*time.Hour)); err != nil {
			fmt.Println("Error storing token:", err)
		}

		return param.RegisterResponse{}, nil
	}
}

func (us Service) consumeMessages() (<-chan interface{}, error) {

	msgs, err := us.eventPublisher.Consume("login_request", "user_consumer", false)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	userChannel := make(chan interface{})

	go func() {
		defer close(userChannel)

		for d := range msgs {
			// Print the raw data received from the queue
			fmt.Println("Raw data:", string(d.Body))

			var user entity.User
			user, err := UnmarshalUser(d.Body) // Call unmarshalUser function
			if err != nil {
				log.Println("Error unmarshalling data:", err)
				// Handle the error accordingly
			} else {
				// Process the user data as needed
				userChannel <- user
				// Acknowledge the message after processing
				d.Ack(false)
			}
		}
	}()

	log.Println("Consuming, to close the program press CTRL+C")
	return userChannel, nil

}

func UnmarshalUser(data []byte) (entity.User, error) {
	var user entity.User
	err := json.Unmarshal(data, &user)
	return user, err
}

// GenerateToken generates a JWT token
func (us Service) GenerateToken(userID uint, expiration time.Duration) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(expiration).Unix()

	// Sign the token with a secret
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ! Helper functions
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
