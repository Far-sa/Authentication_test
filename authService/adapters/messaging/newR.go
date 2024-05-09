package messaging

import (
	"auth-svc/internal/ports"
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	//"github.com/rabbitmq/amqp091-go/amqp"
	// "github.com/streadway/amqp" // Import the streadway/amqp library
)

type RabbitMQConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	//Vhost    string
	//Url string
}

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
func (rc *RabbitMQClient) CreateExchange(name string, kind string) error {
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

// CreateQueue declares a new queue on the RabbitMQ server
func (rc *RabbitMQClient) CreateQueue(name string) (string, error) {
	ch, err := rc.GetChannel()
	if err != nil {
		return "", err
	}
	defer ch.Close() // Close the channel after use

	queue, err := ch.QueueDeclare(
		name,  // Name of the queue
		true,  // Durable (survives server restarts)
		false, // Delete when unused
		false, // Exclusive (only this connection can access)
		false,
		nil, // Arguments
	)
	if err != nil {
		return "", err
	}

	return queue.Name, nil
}

// BindQueue binds an existing queue to an existing exchange with a routing key
func (rc *RabbitMQClient) BindQueue(queueName string, exchangeName string, routingKey string) error {
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

// PublishMessage sends a message to a specific exchange with a routing key
func (rc *RabbitMQClient) PublishMessage(ctx context.Context, exchangeName string, routingKey string, options amqp.Publishing) error {
	ch, err := rc.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close() // Close the channel after use

	// body, err := json.Marshal(message) // Marshal the message to JSON
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal message: %w", err)
	// }

	confirmation, err := ch.PublishWithDeferredConfirmWithContext(
		ctx,
		exchangeName, // Name of the exchange
		routingKey,   // Routing key for message
		false,        // Mandatory (if true, message is rejected if no queue is bound)
		false,        // Immediate (if true, delivery happens now, or fails)
		options,
	)

	if err != nil {
		return err
	}

	log.Println(confirmation.Wait())
	// confirmation.Wait()
	return nil

}

// ! Consume
func (rc *RabbitMQClient) Consume(ctx context.Context, queueName string, callback func(message amqp.Delivery) error) error {
	ch, err := rc.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close() // Close the channel after use

	msgs, err := ch.Consume(
		queueName,
		"",    // Consumer tag (can be left empty)
		false, // Auto-ack (set to false for manual ack)
		false, // Exclusive (only this consumer can access the queue)
		false, // No local (only deliver to this server)
		false, // No wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	//!--->
	go func() {

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Consume cancelled due to context")
				return
			case msg := <-msgs:
				// Process the message and handle ack
				if err := callback(msg); err != nil {
					fmt.Printf("Error processing message: %v\n", err)
					// Optional: You can potentially re-queue the message here
				}
			}
		}
	}()

	return nil
}

// This function can be implemented in your service layer to handle the received message data with ctx
func callbackFunction(msg amqp.Delivery) error {
	var data map[string]interface{}
	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		// Optional: You can potentially re-queue the message here
		return err
	}

	// Process the data from the message (implement your business logic here)
	fmt.Printf("Received message: %v\n", data)
	// ... your message processing logic ...

	// Acknowledge the message after successful processing
	err = msg.Ack(false)
	if err != nil {
		fmt.Printf("Error acknowledging message: %v\n", err)
	}
	return err
}
