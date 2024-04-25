package ports

import (
	"auth-svc/internal/entity"
	"auth-svc/internal/param"
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthRepository interface {
	AddRevokedToken(tokenID string) error
	IsTokenRevoked(tokenID string) bool

	StoreToken(userID int, token string, expiration time.Time) error
	RetrieveToken(userID int) (entity.Token, error)
}

type AuthService interface {
	// RevokeToken(ctx context.Context, tokenID string) error
	Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error)
}

type EventPublisher interface {
	CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error)
	CreateBinding(name, binding, exchange string) error
	Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
}
