package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// import (
// 	"strings"
// 	"time"

// 	"github.com/golang-jwt/jwt/v4"
// )

// import (
// 	"auth-svc/internal/param"
// 	"auth-svc/internal/ports"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"syscall"
// 	"time"

// 	"github.com/golang-jwt/jwt/v4"
// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"golang.org/x/crypto/bcrypt"
// 	//"github.com/dgrijalva/jwt-go"
// )

// type Config struct {
// 	JwtSignKey                     string
// 	AccessTokenSubject             string
// 	RefreshTokenSubject            string
// 	AccessTokenExpirationDuration  time.Duration
// 	RefreshTokenExpirationDuration time.Duration
// }

// TODO: add to config
// const (
// 	JwtSignKey                     = "jwt-secret"
// 	AccessTokenSubject             = "at"
// 	RefreshTokenSubject            = "rt"
// 	AccessTokenExpirationDuration  = time.Hour * 24
// 	RefreshTokenExpirationDuration = time.Hour * 24 * 7
// )

// type authService struct {
// 	config   ports.Config
// 	authRepo ports.AuthRepository
// 	event    ports.EventPublisher
// 	// event    ports.EventPublisher
// }

// User represents the user data received from the message
// type User struct {
// 	ID       uint   `json:"id"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// // NewTokenHandler creates a new TokenHandler with the given authService
// func NewAuthService(config ports.Config, authRepo ports.AuthRepository, event ports.EventPublisher) authService {
// 	return authService{config: config, authRepo: authRepo, event: event}
// }

// func (s authService) Login(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {

// 	//! sequntial

// 	err := s.consumeMessages()
// 	if err != nil {
// 		return param.LoginResponse{}, fmt.Errorf("failed to consume messages: %w", err)
// 	}

// 	var user User

// 	Iterate over userData slice and extract data
// 	for _, u := range userData {

// 		user = User{
// 			ID:       u.ID,
// 			Email:    u.Email,
// 			Password: u.Password,
// 		}
// 		log.Println("User ID:", u.ID)
// 		log.Println("User Email:", u.Email)
// 		log.Println("Password", u.Password)
// 		// Extract other fields as needed
// 	}

// 	hashedPassword, err := HashPassword(req.Password)
// 	if err != nil {
// 		return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
// 		return param.LoginResponse{}, fmt.Errorf("failed to compare passwords: %w", err)

// 	}

// 	accessToken, err := s.createAccessToken(user)
// 	if err != nil {
// 		return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
// 	}

// 	refreshToken, err := s.refreshAccessToken(user)
// 	if err != nil {
// 		return param.LoginResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
// 	}

// 	if err := s.authRepo.StoreToken(int(user.ID), accessToken, time.Now().Add(72*time.Hour)); err != nil {
// 		fmt.Println("Error storing token:", err)
// 	}

// 	return param.LoginResponse{
// 		User:   param.UserInfo{ID: uint(user.ID), Email: user.Email},
// 		Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
// 	}, nil

// }

// TODO Bug- (change exchange)
// func (s authService) consumeMessages() ([]User, error) {

// 	var allUsers []User // Declare a slice to store processed users

// 	if err := s.event.DeclareExchange("user_events", "topic"); err != nil {
// 		return nil, fmt.Errorf("failed to create exchange: %w", err) // Propagate error
// 	}

// 	queue, err := s.event.CreateQueue("registrations_queue", true, false)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create queue: %w", err) // Propagate error
// 	}

// 	if err := s.event.CreateBinding(queue.Name, "registration.*", "user_events"); err != nil {
// 		return nil, fmt.Errorf("failed to bind queue: %w", err) // Propagate error
// 	}

// 	msgs, err := s.event.Consume(queue.Name, "auth_consumer", false)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to consume messages: %w", err)
// 	}

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
// 	go func() {
// 		for message := range msgs {
// 			log.Printf("Received a message: %s", message.Body)
// 			message.Ack(false)

// 		}
// 	}()

// 	log.Println("Consuming, to close the program press CTRL+C")
// 	<-signals

// 	var wg sync.WaitGroup
// 	const numWorkers = 5
// 	resultChan := make(chan User, numWorkers) // Buffered channel for processed users

// 	// Start worker goroutines to process messages concurrently
// 	for i := 0; i < numWorkers; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for msg := range msgs {
// 				log.Printf("Received a message: %s", msg.Body)
// 				userData, err := s.processMessages(msg)
// 				if err != nil {
// 					log.Printf("Error processing message: %v", err)
// 					// Handle error or re-queue message if needed
// 					continue
// 				}
// 				resultChan <- userData // Send processed user to the channel
// 				// allUsers = append(allUsers, userData) // Add processed user to the slice
// 			}
// 		}()
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(resultChan) // Close the channel once all workers finish
// 	}()

// 	// Collect results from the channel
// 	for userData := range resultChan {
// 		allUsers = append(allUsers, userData)
// 	}

// 	return allUsers, nil
// }

// func (s authService) processMessages(msg amqp.Delivery) (User, error) {
// 	var userData User

// 	err := json.Unmarshal(msg.Body, &userData)
// 	if err != nil {
// 		fmt.Printf("Error unmarshalling message: %v\n", err)
// 		// Optional: You can potentially re-queue the message here
// 		return User{}, err
// 	}

// 	fmt.Printf("processing user: %v\n", userData)
// 	// Validate credentials (replace with your validation logic)

// 	//! Publish the new message with processed data and message ID
// 	processedData, err := json.Marshal(userData)
// 	if err != nil {
// 		return User{}, err
// 	}

// 	err = s.event.PublishMessage("user_processed_events", "auth_routing_key", amqp.Publishing{
// 		ContentType:   "application/json",
// 		DeliveryMode:  amqp.Persistent,
// 		Body:          []byte(processedData),
// 		CorrelationId: msg.MessageId,
// 	})
// 	if err != nil {
// 		fmt.Printf("Error publishing processed message: %v\n", err)
// 	}

// 	return userData, nil
// }

// ! helper function
// func HashPassword(password string) (string, error) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(hashedPassword), nil
// }

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
