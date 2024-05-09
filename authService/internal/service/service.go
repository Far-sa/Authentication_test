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
	"github.com/rabbitmq/amqp091-go"
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
	event    ports.RabbitMQ
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
	return authService{config: config, authRepo: authRepo}
}

func (s authService) Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error) {

	//! sequntial
	// Hash the password from the incoming request
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	queueName := s.config.GetBrokerConfig().Queues[0].Name
	exName := s.config.GetBrokerConfig().Exchanges[0].Name
	routeKey := s.config.GetBrokerConfig().Bindings[0].RoutingKey

	var userData User
	var accessToken, refreshToken string
	userData = User{}

	// q, err := s.event.CreateQueue(queueName)
	// if err != nil {
	// 	fmt.Println("Error binding queue", err)

	// }

	if err := s.event.BindQueue(queueName, exName, routeKey); err != nil {
		fmt.Println("Error binding queue", err)
	}

	err = s.event.Consume(ctx, queueName, func(msg amqp091.Delivery) error {
		err := json.Unmarshal(msg.Body, &userData)
		if err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			return err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(hashedPassword)); err != nil {
			return err
		}

		accessToken, err = s.createAccessToken(userData)
		if err != nil {
			return err
		}

		refreshToken, err = s.refreshAccessToken(userData)
		if err != nil {
			return err
		}

		if err := s.authRepo.StoreToken(int(userData.ID), accessToken, time.Now().Add(72*time.Second)); err != nil {
			fmt.Println("Error storing token:", err)
		}

		err = msg.Ack(false)
		if err != nil {
			fmt.Printf("Error acknowledging message: %v\n", err)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error consuming messages:", err)
	}

	return param.LoginResponse{
		User:   param.UserInfo{ID: uint(userData.ID), Email: userData.Email},
		Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	}, nil

	//!----> handle wirh go routines
	// Hash the password from the incoming request
	// hashedPassword, err := HashPassword(user.Password)
	// if err != nil {
	// 	return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
	// }

	// ctxWithCan, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// wg := sync.WaitGroup{}
	// wg.Add(1)

	// queueName := s.config.GetBrokerConfig().Queues[0].Name

	// // Initiate message consumption in a background goroutine
	// go func() {
	// 	defer wg.Done()

	// 	cErr := s.event.Consume(ctxWithCan, queueName, func(msg amqp.Delivery) error {
	// 		var userData User
	// 		err := json.Unmarshal(msg.Body, &userData)
	// 		if err != nil {
	// 			fmt.Printf("Error unmarshalling message: %v\n", err)
	// 			// Optional: You can potentially re-queue the message here
	// 			return err
	// 		}

	// 		// Validate credentials against data from the message
	// 		if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(hashedPassword)); err != nil {
	// 			return err
	// 		}

	// 		// ... rest of your processing logic using userData ...
	// 		accessToken, err := s.createAccessToken(userData) // Assuming a helper function
	// 		if err != nil {
	// 			return err
	// 		}

	// 		refreshToken, err := s.refreshAccessToken(userData) // Assuming a helper function
	// 		if err != nil {
	// 			return err
	// 		}

	// 		if err := s.authRepo.StoreToken(int(userData.ID), accessToken,
	// 			time.Now().Add(72*time.Second)); err != nil {
	// 			fmt.Println("Error store token", err)
	// 		}

	// 		// Acknowledge the message after successful processing
	// 		err = msg.Ack(false)
	// 		if err != nil {
	// 			fmt.Printf("Error acknowledging message: %v\n", err)
	// 		}
	// 		return err
	// 	})
	// 	if cErr != nil {
	// 		fmt.Println("Error consuming messages:", cErr)
	// 		return cErr // Propagate error from Consume
	// 	}
	// }()

	// // Login function doesn't return anything until message processing finishes
	// wg.Wait() // Wait for the goroutine, blocking the login call

	// // Consider redesigning Login flow to be asynchronous (optional)

	// // Handle potential errors from the goroutine (optional)
	// if err != nil {
	// 	return param.LoginResponse{}, fmt.Errorf("login processing failed: %w", err)
	// }

	// // Placeholder for successful login response (assuming success from the goroutine)
	// return param.LoginResponse{
	// 	User:   param.UserInfo{ID: uint(userData.ID), Email: userData.Email}, // Might need modification
	// 	Tokens: param.Tokens{AccessToken: "", RefreshToken: ""},              // Placeholder values
	// }, nil
	//!----->
}

// ?------->
func processMessages(msg amqp091.Delivery) error {
	var userData User
	err := json.Unmarshal(msg.Body, &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		// Optional: You can potentially re-queue the message here
		return err
	}

	// Process the data from the message (implement your business logic here)
	fmt.Printf("processing user: %v\n", userData)
	// ... your message processing logic ...
	// Validate credentials (replace with your validation logic)
	// Acknowledge the message after successful processing
	err = msg.Ack(false)
	if err != nil {
		fmt.Printf("Error acknowledging message: %v\n", err)
	}
	return err
}

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
