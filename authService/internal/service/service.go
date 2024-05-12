package service

import (
	"auth-svc/internal/param"
	"auth-svc/internal/ports"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
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
	// event    ports.EventPublisher
}

// User represents the user data received from the message
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewAuthService(config ports.Config, authRepo ports.AuthRepository, event ports.EventPublisher) authService {
	return authService{config: config, authRepo: authRepo, event: event}
}

func (s authService) Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error) {

	//! sequntial
	// hashedPassword, err := HashPassword(user.Password)
	// if err != nil {
	// 	return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
	// }

	userData, err := s.consumeMessages()
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to consume messages: %w", err)
	}

	// Iterate over userData slice and extract data
	for _, user := range userData {
		fmt.Println("User ID:", user.ID)
		fmt.Println("User Email:", user.Email)
		// Extract other fields as needed
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(hashedPassword)); err != nil {
	// 	return param.LoginResponse{}, fmt.Errorf("failed to compare passwords: %w", err)

	// }

	// accessToken, err := s.createAccessToken(userData)
	// if err != nil {
	// 	return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
	// }

	// refreshToken, err := s.refreshAccessToken(userData)
	// if err != nil {
	// 	return param.LoginResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
	// }

	// if err := s.authRepo.StoreToken(int(userData.ID), accessToken, time.Now().Add(72*time.Second)); err != nil {
	// 	fmt.Println("Error storing token:", err)
	// }

	// return param.LoginResponse{
	// 	User:   param.UserInfo{ID: uint(userData.ID), Email: userData.Email},
	// 	Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	// }, nil

	return param.LoginResponse{User: param.UserInfo{Email: user.Email}}, nil
}

func (s authService) consumeMessages() ([]User, error) {
	var allUsers []User // Declare a slice to store processed users

	if err := s.event.DeclareExchange("user_events", "direct"); err != nil {
		return nil, fmt.Errorf("failed to create exchange: %w", err) // Propagate error
	}

	queue, err := s.event.CreateQueue("user_registrations", true, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err) // Propagate error
	}

	if err := s.event.CreateBinding(queue.Name, "auth_routing_key", "user_events"); err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err) // Propagate error
	}

	msgs, err := s.event.Consume(queue.Name, "auth_consumer", false)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	var wg sync.WaitGroup
	const numWorkers = 5
	// Create a buffered channel to store messages
	msgChan := make(chan amqp.Delivery, numWorkers)

	// Start worker goroutines to process messages concurrently
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range msgChan {
				fmt.Println("Received a message: %s", msg.Body)
				userData, err := s.processMessages(msg)
				if err != nil {
					log.Printf("Error processing message: %v", err)
					// Handle error or re-queue message if needed
					continue
				}
				allUsers = append(allUsers, userData) // Add processed user to the slice
			}
		}()
	}

	log.Println("Waiting for messages. To exit press CTRL+C")

	// Start a goroutine to read messages from the event and send them to the channel
	go func() {
		for msg := range msgs {
			msgChan <- msg
		}
	}()

	// Wait for all worker goroutines to finish
	wg.Wait()
	// Print the data after all messages have been processed
	for _, user := range allUsers {
		fmt.Println("User ID:", user.ID)
		fmt.Println("User Email:", user.Email)
		// Print other user data as needed
	}
	return allUsers, nil

}

func (s authService) processMessages(msg amqp.Delivery) (User, error) {
	var userData User
	err := json.Unmarshal(msg.Body, &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		// Optional: You can potentially re-queue the message here
		return User{}, err
	}

	fmt.Printf("processing user: %v\n", userData)
	// Validate credentials (replace with your validation logic)

	// Acknowledge the message after successful processing
	err = msg.Ack(false)
	if err != nil {
		fmt.Printf("Error acknowledging message: %v\n", err)
	}
	return userData, nil
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
