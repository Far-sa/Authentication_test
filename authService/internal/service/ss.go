package service

import (
	"auth-svc/internal/param"
	"auth-svc/internal/ports"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// import (
// 	"auth-svc/internal/param"
// 	"auth-svc/internal/ports"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"golang.org/x/crypto/bcrypt"
// 	//"github.com/dgrijalva/jwt-go"
// )

// // TODO: add to config
// const (
// 	JwtSignKey                     = "jwt-secret"
// 	AccessTokenSubject             = "at"
// 	RefreshTokenSubject            = "rt"
// 	AccessTokenExpirationDuration  = time.Hour * 24
// 	RefreshTokenExpirationDuration = time.Hour * 24 * 7
// )

type authSvc struct {
	config         ports.Config
	authRepo       ports.AuthRepository
	eventPublisher ports.EventPublisher
	// event    ports.EventPublisher
}

// // User represents the user data received from the message
type UserL struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewAuthSvc(config ports.Config, authRepo ports.AuthRepository, eventPublisher ports.EventPublisher) authSvc {
	return authSvc{config: config, authRepo: authRepo, eventPublisher: eventPublisher}
}

func (s authSvc) SignUp(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {

	// Create a LoginRequest struct (assuming it has Email and Password fields)
	loginReq := param.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.publishLoginRequest(ctx, loginReq)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to publish login request: %w", err)
	}

	userChan, err := s.ConsumeMessages()
	if err != nil {

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

		if err := s.authRepo.StoreToken(user.ID, accessToken, time.Now().Add(72*time.Hour)); err != nil {
			fmt.Println("Error storing token:", err)
		}

		return param.LoginResponse{
			User:   param.UserInfo{ID: uint(user.ID), Email: user.Email},
			Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
		}, nil
	}
	// Return a message indicating login request is processing
}

func (s authSvc) publishLoginRequest(ctx context.Context, req param.LoginRequest) error {

	if err := s.eventPublisher.DeclareExchange("login_requests_exchange", "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := s.eventPublisher.CreateQueue("login_requests", true, false)
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err) // Propagate error
	}

	if err := s.eventPublisher.CreateBinding(queue.Name, "login_request.*", "login_requests_exchange"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err) // Propagate error
	}

	data, jErr := json.Marshal(req)
	if jErr != nil {
		return fmt.Errorf("failed to marshal login request: %w", jErr)
	}

	if err := s.eventPublisher.Publish(ctx, "user_events", "registration.new", amqp.Publishing{
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

func (s authSvc) ConsumeMessages() (<-chan interface{}, error) {

	msgs, err := s.eventPublisher.Consume("login_requests", "auth_consumer", false)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	userChannel := make(chan interface{})

	go func() {
		defer close(userChannel)

		for d := range msgs {
			var data interface{}
			err := json.Unmarshal(d.Body, &data) // Unmarshal to generic interface
			if err != nil {
				log.Println("Error unmarshalling data:", err)
				// Handle the error accordingly
				continue
			}

			// The existing auth service consumer (listening to "login_requests") can receive
			// both the initial login request and the user service response on the same queue.
			switch msg := data.(type) {
			case param.LoginRequest: // Handle login request
				// You can access login request data directly from msg (assuming correct type)
				userChannel <- msg // Send login request for further processing (optional)
				// ... (optional logic for handling login request within auth service) ...
			case string: // Handle user service response (assuming string message)
				if msg == "user_validated" {
					// User validated, generate tokens
					// ... (logic for token generation and response) ...
				} else if msg == "user_not_found" {
					// User not found, return error response
					returnError := errors.New("user not found")
					userChannel <- returnError // Send error message (optional)
				} else {
					// Handle unexpected message type
					log.Printf("Unknown message type received: %s", msg)
				}
			default:
				// Handle unexpected data type
				log.Printf("Unexpected data type received: %T", data)
			}

			d.Ack(false) // Acknowledge the message after processing
		}
	}()

	return userChannel, nil
}

// func (s authService) Login(ctx context.Context, req param.LoginRequest) (param.LoginResponse, error) {

// 	//! sequntial

// 	userChan, err := s.consumeMessages()
// 	if err != nil {
// 		return param.LoginResponse{}, fmt.Errorf("failed to consume messages: %w", err)
// 	}

// 	select {
// 	case <-ctx.Done():
// 		// Handle timeout or context cancellation
// 		return param.LoginResponse{}, errors.New("timed out waiting for user information")
// 	case data := <-userChan:
// 		user, ok := data.(User) // Cast data to the User struct
// 		if !ok {
// 			return param.LoginResponse{}, errors.New("invalid data received from queue")
// 		}

// 		hashedPassword, err := HashPassword(req.Password)
// 		if err != nil {
// 			return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
// 		}

// 		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
// 			return param.LoginResponse{}, fmt.Errorf("failed to compare passwords: %w", err)

// 		}

// 		accessToken, err := s.createAccessToken(user)
// 		if err != nil {
// 			return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
// 		}

// 		refreshToken, err := s.refreshAccessToken(user)
// 		if err != nil {
// 			return param.LoginResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
// 		}

// 		if err := s.authRepo.StoreToken(int(user.ID), accessToken, time.Now().Add(72*time.Hour)); err != nil {
// 			fmt.Println("Error storing token:", err)
// 		}

// 		return param.LoginResponse{
// 			User:   param.UserInfo{ID: uint(user.ID), Email: user.Email},
// 			Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
// 		}, nil
// 	}

// }

// // TODO Bug- (change exchange)
// func (s authService) consumeMessages() (<-chan interface{}, error) {
// 	rabbitMQURL := "amqp://guest:guest@rabbitmq:5672/"

// 	conn, err := amqp.Dial(rabbitMQURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		conn.Close()
// 		return nil, err
// 	}

// 	err = ch.ExchangeDeclare(
// 		"topic_exchange", // name
// 		"topic",          // type
// 		true,             // durable
// 		false,            // auto-deleted
// 		false,            // internal
// 		false,            // no-wait
// 		nil,              // arguments
// 	)
// 	if err != nil {
// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}

// 	q, err := ch.QueueDeclare(
// 		"registration_queue", // name
// 		true,                 // durable
// 		false,                // delete when unused
// 		false,                // exclusive
// 		false,                // no-wait
// 		nil,                  // arguments
// 	)
// 	if err != nil {
// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}

// 	err = ch.QueueBind(
// 		q.Name,           // queue name
// 		"registration.*", // routing key
// 		"topic_exchange", // exchange
// 		false,            // no-wait
// 		nil,              // arguments
// 	)
// 	if err != nil {
// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}

// 	msgs, err := ch.Consume(
// 		"registration_queue", // queue
// 		"",                   // consumer
// 		false,                // auto-ack (set to false for manual ack)
// 		false,                // exclusive
// 		false,                // no-local
// 		false,                // no-wait
// 		nil,                  // args
// 	)
// 	if err != nil {
// 		log.Fatal("Failed to consume messages:", err)

// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

// 	userChannel := make(chan interface{})

// 	go func() {
// 		for d := range msgs {
// 			var user User
// 			err := json.Unmarshal(d.Body, &user)
// 			if err != nil {
// 				log.Println("Error unmarshalling data:", err)
// 				// Handle the error accordingly
// 			} else {
// 				// Process the user data as needed
// 				userChannel <- user
// 				// Acknowledge the message after processing
// 				d.Ack(false)
// 			}
// 		}
// 	}()

// 	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
// 	<-signals
// 	return userChannel, nil
// }

// // ! helper function
// func HashPassword(password string) (string, error) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(hashedPassword), nil
// }
