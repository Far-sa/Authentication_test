package messaging

import (
	"auth-svc/internal/ports"
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

func (rc *RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}
