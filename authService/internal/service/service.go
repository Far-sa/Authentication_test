package service

import (
	"auth-svc/internal/param"
	"auth-svc/internal/ports"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
	//"github.com/dgrijalva/jwt-go"
)

// type Config struct {
// 	JwtSignKey                     string
// 	AccessTokenSubject             string
// 	RefreshTokenSubject            string
// 	AccessTokenExpirationDuration  time.Duration
// 	RefreshTokenExpirationDuration time.Duration
// }

const (
	JwtSignKey                     = "jwt-secret"
	AccessTokenSubject             = "at"
	RefreshTokenSubject            = "rt"
	AccessTokenExpirationDuration  = time.Hour * 24
	RefreshTokenExpirationDuration = time.Hour * 24 * 7
)

type authService struct {
	authRepo ports.AuthRepository
	event    ports.EventPublisher
}

// User represents the user data received from the message
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewAuthService(authRepo ports.AuthRepository, event ports.EventPublisher) authService {
	return authService{authRepo: authRepo}
}

func (s authService) Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error) {

	// Hash the password from the incoming request
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	//! Consume user info messages from RabbitMQ
	go s.consumeUserMessages()

	// Wait for user data to be available before proceeding
	select {
	case <-time.After(5 * time.Second): // Timeout after 5 seconds (adjust as needed)
		return param.LoginResponse{}, fmt.Errorf("timeout waiting for user data")
	case userData := <-userMsgReceived: // Receive user data from the channel
		// Proceed with user data
		fmt.Printf("Received user data in Login function: %+v\n", userData)

		// Validate user credentials
		if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(hashedPassword)); err != nil {
			return param.LoginResponse{}, fmt.Errorf("username/password incorrect")
		}

		// Create tokens
		accessToken, err := s.createAccessToken(userData)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
		}

		refreshToken, err := s.refreshAccessToken(userData)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to refresh access token: %w", err)
		}

		if err := s.authRepo.StoreToken(int(userData.ID), accessToken,
			time.Now().Add(72*time.Second)); err != nil {
			fmt.Println("Error store token", err)
		}

		// Return login response with tokens
		return param.LoginResponse{
			User:   param.UserInfo{ID: userData.ID, Email: userData.Email},
			Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
		}, nil
	}
}

func (s authService) consumeUserMessages() error {
	// Declare exchange (if needed)
	if err := s.event.DeclareExchange("user_data_exchange", "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Create queue
	queue, err := s.event.CreateQueue("auth_queue", true, false)
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}

	// Create binding
	if err := s.event.CreateBinding(queue.Name, "auth_routing_key", "user_data_exchange"); err != nil {
		return fmt.Errorf("failed to create binding: %w", err)
	}

	// Consume messages
	msgs, err := s.event.Consume(queue.Name, "auth_consumer", false)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	// Process messages
	for message := range msgs {
		go s.processUserMessage(message)
	}
	return nil
}

// ? just for signal
// var userMsgReceived = make(chan struct{})
// TODO: move global channel to writer
var userMsgReceived = make(chan User)

func (s authService) processUserMessage(message amqp.Delivery) {
	// Unmarshal the message body into the user struct
	var user User
	if err := json.Unmarshal(message.Body, &user); err != nil {
		fmt.Println("failed to unmarshal message:", err)
		// Handle the error appropriately (e.g., logging, error reporting)
		// Acknowledge or reject the message, depending on your requirements
		message.Ack(false)
		return
	}

	// Process the user data
	fmt.Printf("Received user data: %+v\n", user)

	// Signal that user data is available
	go func() {
		userMsgReceived <- user
		//userMsgReceived <- struct{}{}
	}()

	// Acknowledge the message to RabbitMQ to indicate successful processing
	message.Ack(false)
}

// ! helper function
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s authService) createAccessToken(user User) (string, error) {
	return s.createToken(user.ID, AccessTokenSubject, AccessTokenExpirationDuration)
}

func (s authService) refreshAccessToken(user User) (string, error) {
	return s.createToken(user.ID, RefreshTokenSubject, RefreshTokenExpirationDuration)
}

func (s authService) VerifyToken(bearerToken string) (*Claims, error) {

	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})

	if err != nil {
		return nil, err
	}

	// convert interface to conceret object
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil

	} else {
		return nil, err
	}
}

// "github.com/golang-jwt/jwt/v4"
type Claims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

func (s authService) createToken(userID uint, subject string, expiresDuration time.Duration) (string, error) {
	// set our claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
			// set the expire time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresDuration)),
		},
		UserID: userID,
	}

	// TODO add sign method to config
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := accessToken.SignedString([]byte(JwtSignKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// func (s authService) AddRevokedToken(tokenID string) error {
// 	panic("")
// }

// func (s authService) IsRevokedToken(tokenID string) error {
// 	panic("")
// }

// ! Middleware for token validation in Traefik
// func TokenValidationMiddleware(next http.Handler, authService authService) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Extract token from request
// 		token := ExtractTokenFromRequest(r)

// 		// Validate token
// 		if isValid := authService.ValidateToken(token); !isValid {
// 			// Token is invalid or revoked
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Token is valid, proceed to next handler
// 		next.ServeHTTP(w, r)
// 	})
// }

// // ExtractTokenFromRequest extracts JWT token from request
// func ExtractTokenFromRequest(r *http.Request) string {
// 	// Extract token from request headers, cookies, or query parameters
// 	// Example: Authorization: Bearer <token>
// 	token := r.Header.Get("Authorization")
// 	if token != "" {
// 		return strings.TrimPrefix(token, "Bearer ")
// 	}

// 	// Extract token from cookies or query parameters if needed

// 	return ""
// }
