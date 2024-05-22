package messaging

import (
	"auth-svc/internal/ports"
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	config ports.Config
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func NewRabbitMQClient(config ports.Config) (*RabbitClient, error) {
	cfg := config.GetBrokerConfig()
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // Close the connection on channel opening error
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Enable publisher confirms (optional)
	if err := ch.Confirm(false); err != nil {
		ch.Close()   // Close the channel on Confirm error
		conn.Close() // Also close the connection
		return nil, fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	return &RabbitClient{
		config: config,
		conn:   conn,
		ch:     ch,
	}, nil
}

func (rc *RabbitClient) Close() error {
	return rc.ch.Close()
}

// CreateExchange declares a new exchange on the RabbitMQ server
func (rc RabbitClient) DeclareExchange(name, kind string) error {

	return rc.ch.ExchangeDeclare(
		name,  // Name of the exchange
		kind,  // Type of exchange (e.g., "fanout", "direct", "topic")
		true,  // Durable (survives server restarts)
		false, // Delete when unused
		false, // Exclusive (only this connection can access)
		false,
		nil, // Arguments
	)
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error) {
	q, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	if err != nil {
		return amqp.Queue{}, nil
	}
	return q, err
}

func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

// ! PublishMessage sends a message to a specific exchange with a routing key
func (rc RabbitClient) Publish(ctx context.Context, exchangeName string, routingKey string, options amqp.Publishing) error {

	return rc.ch.PublishWithContext(
		ctx,
		exchangeName, // Name of the exchange
		routingKey,   // Routing key for message
		false,        // Mandatory (if true, message is rejected if no queue is bound)
		false,        // Immediate (if true, delivery happens now, or fails)
		options,
	)

}

func (rc *RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

func (rc RabbitClient) DeleteQueue(name string) error {
	_, err := rc.ch.QueueDelete(name, false, false, false)
	return err
}
