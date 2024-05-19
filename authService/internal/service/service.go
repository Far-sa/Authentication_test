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
	return authService{config: config, authRepo: authRepo, event: event}
}

func (s authService) Login(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {

	//! sequntial

	userChan, err := s.consumeMessages()
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to consume messages: %w", err)
	}

	select {
	case <-ctx.Done():
		// Handle timeout or context cancellation
		return param.LoginResponse{}, errors.New("timed out waiting for user information")
	case data, ok := <-userChan:
		if !ok {
			return param.LoginResponse{}, errors.New("user not found") // No matching user found
		}
		user, ok := data.(User) // Cast data to the User struct
		if !ok {
			return param.LoginResponse{}, errors.New("invalid data received from queue")
		}
		// case data := <-userChan:
		// 	user, ok := data.(User) // Cast data to the User struct
		// 	if !ok {
		// 		return param.LoginResponse{}, errors.New("invalid data received from queue")
		// 	}

		valid, err := ComparePassword(user.Password, req.Password)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to compare password: %w", err)
		}

		if !valid {
			return param.LoginResponse{}, errors.New("invalid password") // Indicate invalid password
		}

		accessToken, err := s.createAccessToken(user)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
		}

		refreshToken, err := s.refreshAccessToken(user)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
		}

		if err := s.authRepo.StoreToken(user.ID, accessToken, time.Now().Add(72*time.Hour)); err != nil {
			fmt.Println("Error storing token:", err)
		}

		return param.LoginResponse{
			User:   param.UserInfo{ID: uint(user.ID), Email: user.Email},
			Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
		}, nil
	}
}

func (s authService) consumeMessages() (<-chan interface{}, error) {

	msgs, err := s.event.Consume("registration_queue", "auth_consumer", false)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	userChannel := make(chan interface{})

	go func() {
		defer close(userChannel)

		for d := range msgs {
			// Print the raw data received from the queue
			fmt.Println("Raw data:", string(d.Body))

			var user User
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

func UnmarshalUser(data []byte) (User, error) {
	var user User
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

//!!
// var wg sync.WaitGroup
// const numWorkers = 5
// resultChan := make(chan User, numWorkers) // Buffered channel for processed users

// // Start worker goroutines to process messages concurrently
// for i := 0; i < numWorkers; i++ {
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for msg := range msgs {
// 			log.Printf("Received a message: %s", msg.Body)
// 			userData, err := s.processMessages(msg)
// 			if err != nil {
// 				log.rintf("Error processing message: %v", err)
// 				// Handle error or re-queue message if needed
// 				continue
// 			}
// 			resultChan <- userData // Send processed user to the channel
// 			// allUsers = append(allUsers, userData) // Add processed user to the slice
// 		}
// 	}()
// }

// go func() {
// 	wg.Wait()
// 	close(resultChan) // Close the channel once all workers finish
// }()

// // Collect results from the channel
// for userData := range resultChan {
// 	allUsers = append(allUsers, userData)
// }
