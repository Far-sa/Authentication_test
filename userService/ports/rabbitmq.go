package ports

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	GetChannel() (*amqp.Channel, error)
	Close() error

	DeclareExchange(name, kind string) error
	Publish(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error
	CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error)
	CreateBinding(queueName, routingKey, exchangeName string) error
	// Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
}
