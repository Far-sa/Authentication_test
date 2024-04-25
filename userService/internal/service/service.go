package userService

import (
	"context"
	"fmt"
	"time"
	"user-svc/internal/entity"
	"user-svc/internal/service/param"
	"user-svc/ports"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// TODO: add config - zap logger as singleton

type Config struct{}

// Service represents the usecase for user operations
type Service struct {
	userRepo       ports.UserRepository
	eventPublisher ports.EventPublisher
	logger         ports.Logger
	// Cash           ports.Caching
}

// NewService creates a new instance of Service
func NewService(repo ports.UserRepository, publisher ports.EventPublisher, logger ports.Logger) Service {
	return Service{userRepo: repo, eventPublisher: publisher, logger: logger}
}

// RegisterUser handles user registration
func (us Service) Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error) {

	// hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
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

	// TODO generate token

	//* Create user event payload as binary
	eventPayload := []byte(MapStringToByte(createdUser.Email))

	//! Publish user created event
	err = us.eventPublisher.PublishUserRegisteredEvent(ctx, eventPayload)
	if err != nil {
		us.logger.Error("Failed to publish user created event", zap.Error(err))
		return param.RegisterResponse{}, fmt.Errorf("failed to publish user created event: %w", err)
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

func (us Service) Login(user entity.User) error {
	panic("unimplemented")
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
