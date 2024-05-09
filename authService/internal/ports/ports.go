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

	StoreToken(userID int, token string, expiration time.Time) error
	RetrieveToken(userID int) (*entity.Token, error)
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
	DeclareExchange(name, kind string) error
	CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error)
	CreateBinding(name, binding, exchange string) error
	Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
	// PublishUser(userInfo *UserInfo) error

}

type RabbitMQ interface {
	// GetChannel opens a new AMQP channel with context support
	GetChannel() (*amqp.Channel, error)

	// Close closes the underlying AMQP connection
	Close() error

	// CreateExchange declares a new exchange on the RabbitMQ server
	CreateExchange(name string, kind string) error

	// CreateQueue declares a new queue on the RabbitMQ server
	CreateQueue(name string) (string, error)

	// BindQueue binds an existing queue to an existing exchange with a routing key
	BindQueue(queueName string, exchangeName string, routingKey string) error

	// PublishMessage sends a message to a specific exchange with a routing key with context support for cancellation
	PublishMessage(ctx context.Context, exchangeName string, routingKey string, options amqp.Publishing) error
	Consume(ctx context.Context, queueName string, callback func(message amqp.Delivery) error) error
}
