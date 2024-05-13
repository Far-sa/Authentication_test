package service

import (
	"auth-svc/internal/param"
	"auth-svc/internal/ports"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
	//"github.com/dgrijalva/jwt-go"
)

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
	case data := <-userChan:
		user, ok := data.(User) // Cast data to the User struct
		if !ok {
			return param.LoginResponse{}, errors.New("invalid data received from queue")
		}

		// hashedPassword, err := HashPassword(req.Password)
		// if err != nil {
		// 	return param.LoginResponse{}, fmt.Errorf("failed to hash password: %w", err)
		// }

		// if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		// 	return param.LoginResponse{}, fmt.Errorf("failed to compare passwords: %w", err)

		// }

		accessToken, err := s.createAccessToken(user)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to create access token: %w", err)
		}

		refreshToken, err := s.refreshAccessToken(user)
		if err != nil {
			return param.LoginResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
		}

		if err := s.authRepo.StoreToken(int(user.ID), accessToken, time.Now().Add(72*time.Hour)); err != nil {
			fmt.Println("Error storing token:", err)
		}

		return param.LoginResponse{
			User:   param.UserInfo{ID: uint(user.ID), Email: user.Email},
			Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
		}, nil
	}

}

// TODO Bug- (change exchange)
func (s authService) consumeMessages() (<-chan interface{}, error) {
	rabbitMQURL := "amqp://guest:guest@rabbitmq:5672/"

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// err = ch.ExchangeDeclare(
	// 	"topic_exchange", // name
	// 	"topic",          // type
	// 	true,             // durable
	// 	false,            // auto-deleted
	// 	false,            // internal
	// 	false,            // no-wait
	// 	nil,              // arguments
	// )
	// if err != nil {
	// 	conn.Close()
	// 	ch.Close()
	// 	return nil, err
	// }

	// q, err := ch.QueueDeclare(
	// 	"registration_queue", // name
	// 	true,                 // durable
	// 	false,                // delete when unused
	// 	false,                // exclusive
	// 	false,                // no-wait
	// 	nil,                  // arguments
	// )
	// if err != nil {
	// 	conn.Close()
	// 	ch.Close()
	// 	return nil, err
	// }

	// err = ch.QueueBind(
	// 	q.Name,           // queue name
	// 	"registration.*", // routing key
	// 	"topic_exchange", // exchange
	// 	false,            // no-wait
	// 	nil,              // arguments
	// )
	// if err != nil {
	// 	conn.Close()
	// 	ch.Close()
	// 	return nil, err
	// }
	msgs, err := ch.Consume(
		"registration_queue", // queue
		"",                   // consumer
		false,                // auto-ack (set to false for manual ack)
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Fatal("Failed to consume messages:", err)

		conn.Close()
		ch.Close()
		return nil, err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	userChannel := make(chan interface{})

	go func() {
		for d := range msgs {
			var user User
			err := json.Unmarshal(d.Body, &user)
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

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-signals
	return userChannel, nil
}

// ! helper function
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
