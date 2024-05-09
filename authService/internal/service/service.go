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

// TODO: add to config
const (
	JwtSignKey                     = "jwt-secret"
	AccessTokenSubject             = "at"
	RefreshTokenSubject            = "rt"
	AccessTokenExpirationDuration  = time.Hour * 24
	RefreshTokenExpirationDuration = time.Hour * 24 * 7
)

type authService struct {
	config   ports.Config
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
func NewAuthService(config ports.Config, authRepo ports.AuthRepository, event ports.EventPublisher) authService {
	return authService{config: config, authRepo: authRepo}
}

func (s authService) Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error) {

	// Hash the password from the incoming request
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// RabbitMQ configuration (replace placeholders with your configuration)
	queueName := s.config.GetBrokerConfig().Queues[0].Name
	msgs, err := s.event.Consume(queueName, "auth_consumer", false)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to consume messages: %w", err)
	}

	// Process messages sequentially and validate credentials
	for msg := range msgs {
		var userData User
		if err := json.Unmarshal([]byte(msg.Body), &userData); err != nil {
			fmt.Println("failed to unmarshal message:", err)
			continue
		}

		// Validate credentials (replace with your validation logic)
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
			User:   param.UserInfo{ID: uint(userData.ID), Email: userData.Email},
			Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
		}, nil
	}
	return param.LoginResponse{}, fmt.Errorf("username/password incorrect")

}

//* Declare exchange (if needed)
// if err := s.event.DeclareExchange("user_data_exchange", "topic"); err != nil {
// 	return fmt.Errorf("failed to declare exchange: %w", err)
// }

//* Create queue
// queue, err := s.event.CreateQueue("auth_queue", true, false)
// if err != nil {
// 	return fmt.Errorf("failed to create queue: %w", err)
// }

// exchangeName := s.config.GetBrokerConfig().Exchanges[0].Name
// routeKey := s.config.GetBrokerConfig().Bindings[0].RoutingKey

//* Create binding
// if err := s.event.CreateBinding(queueName, routeKey, exchangeName); err != nil {
// 	return fmt.Errorf("failed to create binding: %w", err)
// }

// Consume messages -need queue name and routing key

// ? just for signal
// var userMsgReceived = make(chan struct{})

// ! writer

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

//!-------->
// Define a type for the message
// type UserMessage struct {
//     CorrelationId string
//     // Add other fields as needed
// }

// // Modify consumeUserMessages to return a read-only channel of UserMessage
// func (s authService) consumeUserMessages() <-chan UserMessage {
//     // Create a channel for sending user messages
//     userMsgs := make(chan UserMessage)

//     go func() {
//         defer close(userMsgs)

//         queueName := s.config.GetBrokerConfig().Queues[0].Name
//         msgs, err := s.event.Consume(queueName, "auth_consumer", false)
//         if err != nil {
//             log.Printf("failed to consume messages: %v", err)
//             return
//         }

//         for message := range msgs {
//             // Extract relevant information from the message and send it through the channel
//             userMsg := UserMessage{
//                 CorrelationId: message.CorrelationId,
//                 // Extract other fields as needed
//             }
//             userMsgs <- userMsg
//         }
//     }()

//     return userMsgs
// }

// Define a function to process user messages
// func (s authService) processUserMessages(userMsgs <-chan UserMessage) {
//     for userMsg := range userMsgs {
//         // Process the user message
//         log.Printf("Processing user message with CorrelationId: %s\n", userMsg.CorrelationId)

//         // Perform any necessary operations with the user message, such as database updates, authentication, etc.
//         // Example:
//         // if err := s.processUserMessage(userMsg); err != nil {
//         //     log.Printf("Error processing user message: %v\n", err)
//         // }
//     }
// }

//! ------>

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
