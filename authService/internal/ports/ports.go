package ports

import (
	"auth-svc/internal/entity"
	"auth-svc/internal/param"
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthRepository interface {
	// AddRevokedToken(tokenID string) error
	// IsTokenRevoked(tokenID string) bool

	StoreToken(userID uint, token string, expiration time.Time) error
	RetrieveToken(userID uint) (*entity.Token, error)
}

type AuthService interface {
	// RevokeToken(ctx context.Context, tokenID string) error
	Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error)
}

//	type ConsumeResult struct {
//		Messages <-chan amqp.Delivery
//		Closed   bool
//	}
type EventPublisher interface {
	// GetChannel() (*amqp.Channel, error)
	Close() error

	DeclareExchange(name, kind string) error
	CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error)
	CreateBinding(queueName, routingKey, exchangeName string) error
	Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
	Publish(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error
	// PublishUser(userInfo *UserInfo) error

}
