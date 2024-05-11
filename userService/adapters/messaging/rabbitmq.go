package messaging

import (
	"context"
	"fmt"
	"log"
	"user-svc/ports"

	amqp "github.com/rabbitmq/amqp091-go"
	//"github.com/rabbitmq/amqp091-go/amqp"
	// "github.com/streadway/amqp" // Import the streadway/amqp library
)

type RabbitMQClient struct {
	config ports.Config
	conn   *amqp.Connection
}

func NewRabbitClient(config ports.Config) (*RabbitMQClient, error) {
	cfg := config.GetBrokerConfig()
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &RabbitMQClient{config: config, conn: conn}, nil
}

func (rc *RabbitMQClient) GetChannel() (*amqp.Channel, error) {
	ch, err := rc.conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (rc *RabbitMQClient) Close() error {
	return rc.conn.Close() // Close the underlying connection
}

// CreateExchange declares a new exchange on the RabbitMQ server
func (rc *RabbitMQClient) DeclareExchange(name, kind string) error {
	ch, err := rc.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close() // Close the channel after use

	return ch.ExchangeDeclare(
		name,  // Name of the exchange
		kind,  // Type of exchange (e.g., "fanout", "direct", "topic")
		true,  // Durable (survives server restarts)
		false, // Delete when unused
		false, // Exclusive (only this connection can access)
		false,
		nil, // Arguments
	)
}

// ! PublishMessage sends a message to a specific exchange with a routing key
func (rc *RabbitMQClient) Publish(ctx context.Context, exchangeName string, routingKey string, options amqp.Publishing) error {
	ch, err := rc.GetChannel()
	if err != nil {
		log.Printf("Error getting channel: %v\n", err)
		return err
	}

	defer func() {
		if closeErr := ch.Close(); closeErr != nil {
			log.Printf("Error closing channel: %v\n", closeErr)
		}
	}()

	confirmation, err := ch.PublishWithDeferredConfirmWithContext(
		ctx,
		exchangeName, // Name of the exchange
		routingKey,   // Routing key for message
		false,        // Mandatory (if true, message is rejected if no queue is bound)
		false,        // Immediate (if true, delivery happens now, or fails)
		options,
	)

	if err != nil {
		log.Printf("Error publishing message: %v\n", err)
		return err
	}

	// confirmation.Wait()
	log.Printf("Message published successfully. Confirmation: %v\n", confirmation.Wait())
	return nil

}

// CreateQueue declares a new queue on the RabbitMQ server
func (rc *RabbitMQClient) CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error) {
	ch, err := rc.GetChannel()
	if err != nil {
		return amqp.Queue{}, err
	}
	defer ch.Close() // Close the channel after use

	queue, err := ch.QueueDeclare(
		queueName,  // Name of the queue
		durable,    // Durable (survives server restarts)
		autodelete, // Exclusive (only this connection can access)
		false,      // Delete when unused
		false,
		nil, // Arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	return queue, nil

}

// BindQueue binds an existing queue to an existing exchange with a routing key
func (rc *RabbitMQClient) CreateBinding(queueName, routingKey, exchangeName string) error {
	ch, err := rc.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close() // Close the channel after use

	return ch.QueueBind(
		queueName,    // Name of the queue to bind
		routingKey,   // Routing key for messages
		exchangeName, // Name of the exchange to bind to
		false,        // No wait
		nil,          // Arguments
	)
}

// ! Consume
func (rc *RabbitMQClient) Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	ch, err := rc.GetChannel()
	if err != nil {
		return nil, err
	}
	defer ch.Close() // Close the channel after use

	return ch.Consume(
		queueName,
		consumer, // Consumer tag (can be left empty)
		autoAck,  // Auto-ack (set to false for manual ack)
		false,    // Exclusive (only this consumer can access the queue)
		false,    // No local (only deliver to this server)
		false,    // No wait
		nil,      // Arguments
	)
}

//! QOS
// func (rc RabbitClient) ApplyQos(count, size int, global bool) error {
// 	return rc.ch.Qos(count, size, global)
// }
