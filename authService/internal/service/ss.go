package service

// import (
// 	"auth-svc/internal/param"
// 	"auth-svc/internal/ports"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"strings"
// 	"time"

// 	"github.com/golang-jwt/jwt/v4"
// 	"github.com/google/uuid"
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

// // TODO: add to config
// const (
// 	JwtSignKey                     = "jwt-secret"
// 	AccessTokenSubject             = "at"
// 	RefreshTokenSubject            = "rt"
// 	AccessTokenExpirationDuration  = time.Hour * 24
// 	RefreshTokenExpirationDuration = time.Hour * 24 * 7
// )

// type authService struct {
// 	config         ports.Config
// 	authRepo       ports.AuthRepository
// 	eventPublisher ports.EventPublisher
// }

// func NewAuthService(config ports.Config, authRepo ports.AuthRepository, event ports.EventPublisher) authService {
// 	return authService{config: config, authRepo: authRepo, eventPublisher: event}
// }

// func (s authService) Login(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {

//     // Create LoginRequest struct
//     loginReq := req

//     err := s.publishLoginRequest(ctx, loginReq)
//     if err != nil {
//         return param.LoginResponse{}, fmt.Errorf("failed to publish login request: %w", err)
//     }

//     log.Println("Waiting for user service response")

//     // Increase timeout for testing (adjust as needed)
//     ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
//     defer cancel()

//     responseChan := make(chan param.LoginResponse)
//     errChan := make(chan error)

//     go func() {
//         queue, err := s.createTempQueue(ctx)
//         if err != nil {
//             errChan <- fmt.Errorf("failed to create temporary queue: %w", err)
//             return
//         }
//         defer s.eventPublisher.DeleteQueue(queue.Name, false, false) // Clean up on exit

//         if err := s.eventPublisher.CreateBinding(queue.Name, "user_response", "auth_exchange"); err != nil {
//             errChan <- fmt.Errorf("failed to bind queue: %w", err)
//             return
//         }

//         consumerTag := fmt.Sprintf("auth_service_%s", uuid.NewString())
//         msgs, err := s.eventPublisher.Consume(queue.Name, consumerTag, false)
//         if err != nil {
//             errChan <- fmt.Errorf("failed to consume messages: %w", err)
//             return
//         }

//         for d := range msgs {
//             var response param.LoginResponse
//             err := json.Unmarshal(d.Body, &response)
//             if err != nil {
//                 errChan <- fmt.Errorf("failed to unmarshal response: %w", err)
//                 return
//             }

//             if response.User.Email == loginReq.Email { // Filter by request email
//                 responseChan <- response
//                 d.Ack(false)
//                 return
//             }
//         }

//         errChan <- fmt.Errorf("no matching response found") // Handle no match
//     }()

//     select {
//     case <-ctx.Done():
//         log.Println("Context canceled while waiting for user service response")
//         return param.LoginResponse{}, fmt.Errorf("context canceled")
//     case response := <-responseChan:
//         log.Println("Received user service response")

//         // Process user data (e.g., validate, create tokens)
//         accessToken, err := s.createAccessToken(response)
//         if err != nil {
//             return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
//         }

//         refreshToken, err := s.refreshAccessToken(response)
//         if err != nil {
//             return param.LoginResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
//         }

//         // Save tokens to database (replace with your actual persistence logic)
//         if err := s.authRepo.StoreToken(response.User.ID, accessToken, time.Now().Add(72*time.Hour)); err != nil {
//             return param.LoginResponse{}, fmt.Errorf("failed to store access token: %w", err)
//         }
//         if err := s.authRepo.StoreToken(response.User.ID, refreshToken, time.Now().Add(365*24*time.Hour)); err != nil { // Adjust refresh token expiry as needed
//             return param.LoginResponse{}, fmt.Errorf("failed to store refresh token: %w", err)
//         }

//         return param.LoginResponse{
//             User:   param.UserInfo{ID: response.User.ID, Email: response.User.Email},
//             Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
//         }, nil
//     case err := <-errChan:
//         log.Printf("Error occurred while waiting for user service response: %v", err)
//         return param.LoginResponse{}, err
//     }
// }
// func (s authService) publishLoginRequest(ctx context.Context, req param.LoginRequest) error {

// 	data, jErr := json.Marshal(req)
// 	if jErr != nil {
// 		return fmt.Errorf("failed to marshal login request: %w", jErr)
// 	}

// 	if err := s.eventPublisher.Publish(ctx, "auth_exchange", "login", amqp.Publishing{
// 		ContentType:   "application/json",
// 		DeliveryMode:  amqp.Persistent,
// 		Body:          data,
// 		CorrelationId: uuid.NewString(),
// 	}); err != nil {
// 		return fmt.Errorf("failed to publish user to auth-service: %w", err)
// 	}
// 	log.Println("User event published successfully")

// 	return nil
// }
