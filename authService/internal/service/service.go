package service

import (
	"auth-svc/internal/param"
	"auth-svc/internal/ports"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
	//"github.com/dgrijalva/jwt-go"
)

type Config struct {
	JwtSignKey                     string
	AccessTokenSubject             string
	RefreshTokenSubject            string
	AccessTokenExpirationDuration  time.Duration
	RefreshTokenExpirationDuration time.Duration
}

// TODO: add to config
const (
	JwtSignKey                     = "jwt-secret"
	AccessTokenSubject             = "at"
	RefreshTokenSubject            = "rt"
	AccessTokenExpirationDuration  = time.Hour * 24
	RefreshTokenExpirationDuration = time.Hour * 24 * 7
)

type authService struct {
	config         ports.Config
	authRepo       ports.AuthRepository
	eventPublisher ports.EventPublisher
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewAuthService(config ports.Config, authRepo ports.AuthRepository, event ports.EventPublisher) authService {
	return authService{config: config, authRepo: authRepo, eventPublisher: event}
}

func (s authService) Login(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {

	// Create a LoginRequest struct (assuming it has Email and Password fields)
	loginReq := param.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.publishLoginRequest(ctx, loginReq)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to publish login request: %w", err)
	}

	//! Get response from user service
	fmt.Println("wait for user service response")

	msgs, err := s.waitForConsumeMessages()
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to consume messages: %w", err)
	}

	select {
	case <-ctx.Done():
		// Handle timeout or context cancellation
		return param.LoginResponse{}, errors.New("timed out waiting for user information")
	case data, ok := <-msgs:
		if !ok {
			return param.LoginResponse{}, errors.New("user not found") // No matching user found
		}
		user, ok := data.(param.UserResponse) // Cast data to the User struct
		if !ok {
			return param.LoginResponse{}, errors.New("invalid data received from queue")
		}
		// valid, err := ComparePassword(user.Password, req.Password)
		// if err != nil {
		// 	return param.LoginResponse{}, fmt.Errorf("failed to compare password: %w", err)
		// }

		// if !valid {
		// 	return param.LoginResponse{}, errors.New("invalid password") // Indicate invalid password
		// }

		accessToken, err := s.createAccessToken(user)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
		}

		refreshToken, err := s.refreshAccessToken(user)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
		}

		if err := s.authRepo.StoreToken(user.User.ID, accessToken, time.Now().Add(72*time.Hour)); err != nil {
			fmt.Println("Error storing token:", err)
		}

		return param.LoginResponse{
			User:   param.UserInfo{ID: uint(user.User.ID), Email: user.User.Email},
			Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
		}, nil
	}
}

func (s authService) publishLoginRequest(ctx context.Context, req param.LoginRequest) error {

	// if err := s.eventPublisher.DeclareExchange("auth_exchange", "direct"); err != nil {
	// 	return fmt.Errorf("failed to declare exchange: %w", err)
	// }

	// queue, err := s.eventPublisher.CreateQueue("login_requests", true, false)
	// if err != nil {
	// 	return fmt.Errorf("failed to create queue: %w", err) // Propagate error
	// }

	// if err := s.eventPublisher.CreateBinding(queue.Name, "login", "auth_exchange"); err != nil {
	// 	return fmt.Errorf("failed to bind queue: %w", err) // Propagate error
	// }

	data, jErr := json.Marshal(req)
	if jErr != nil {
		return fmt.Errorf("failed to marshal login request: %w", jErr)
	}

	if err := s.eventPublisher.Publish(ctx, "auth_exchange", "login", amqp.Publishing{
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

func (s authService) waitForConsumeMessages() (<-chan interface{}, error) {

	//TODO for consuming!
	if err := s.eventPublisher.DeclareExchange("auth_exchange", "direct"); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := s.eventPublisher.CreateQueue("user_service_responses", true, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err) // Propagate error
	}

	if err := s.eventPublisher.CreateBinding(queue.Name, "user_response", "auth_exchange"); err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err) // Propagate error

	}

	msgs, err := s.eventPublisher.Consume(queue.Name, "auth_service", false)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	userChannel := make(chan interface{})

	go func() {
		defer close(userChannel)

		for d := range msgs {
			// data UserResponse
			var data interface{}
			data, err := UnmarshalData(d.Body) // Unmarshal to generic interface
			if err != nil {
				log.Println("Error unmarshalling data:", err)
				// Handle the error accordingly
				continue
			} else {
				userChannel <- data
			}

			// The existing auth service consumer (listening to "login_requests") can receive
			// both the initial login request and the user service response on the same queue.
			// switch msg := data.(type) {
			// case param.LoginRequest: // Handle login request
			// 	// You can access login request data directly from msg (assuming correct type)
			// 	userChannel <- msg // Send login request for further processing (optional)
			// 	// ... (optional logic for handling login request within auth service) ...
			// case string: // Handle user service response (assuming string message)
			// 	if msg == "user_validated" {
			// 		// User validated, generate tokens
			// 		// ... (logic for token generation and response) ...
			// 	} else if msg == "user_not_found" {
			// 		// User not found, return error response
			// 		returnError := errors.New("user not found")
			// 		userChannel <- returnError // Send error message (optional)
			// 	} else {
			// 		// Handle unexpected message type
			// 		log.Printf("Unknown message type received: %s", msg)
			// 	}
			// default:
			// 	// Handle unexpected data type
			// 	log.Printf("Unexpected data type received: %T", data)
			// }

			d.Ack(false) // Acknowledge the message after processing
		}
	}()

	return userChannel, nil

}

func UnmarshalData(data []byte) (param.UserResponse, error) {
	var user param.UserResponse
	err := json.Unmarshal(data, &user)
	return user, err
}

func ComparePassword(hashedPassword, reqPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(reqPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil // Indicate invalid password without revealing hashing details
		}
		return false, fmt.Errorf("failed to compare password: %w", err)
	}
	return true, nil // Password matches
}

func (s authService) createAccessToken(user param.UserResponse) (string, error) {
	return s.createToken(user.User.ID, AccessTokenSubject, AccessTokenExpirationDuration)
}

func (s authService) refreshAccessToken(user param.UserResponse) (string, error) {
	return s.createToken(user.User.ID, RefreshTokenSubject, RefreshTokenExpirationDuration)
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

// In a microservices architecture, it's generally a good practice to have each service be responsible for
// creating the queues it will consume from. This ensures that the service can function independently and
// that all necessary resources are in place when the service starts.
