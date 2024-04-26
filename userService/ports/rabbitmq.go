package ports

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	DeclareExchange(name, kind string) error
	CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error)
	CreateBinding(name, binding, exchange string) error
	Publish(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error
	//Consume(ctx context.Context, queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
}
