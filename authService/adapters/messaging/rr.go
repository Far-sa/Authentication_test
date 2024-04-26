package messaging

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageBroker interface defines the operations needed for interacting with a message broker
type MessageBroker interface {
	DeclareExchange(name, kind string) error
	DeclareQueue(name string) error
	BindQueue(queueName, exchangeName, routingKey string) error
	Publish(exchange, routingKey string, body []byte) error
}

// RabbitMQAdapter is an implementation of the MessageBroker interface for RabbitMQ
type RabbitMQAdapter struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewRabbitMQAdapter creates a new instance of RabbitMQAdapter
func NewRabbitMQAdapter(amqpURL string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQAdapter{conn: conn, ch: ch}, nil
}

// DeclareExchange declares a new exchange on RabbitMQ
func (r *RabbitMQAdapter) DeclareExchange(name, kind string) error {
	return r.ch.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue declares a new queue on RabbitMQ
func (r *RabbitMQAdapter) DeclareQueue(name string) error {
	_, err := r.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// BindQueue binds a queue to an exchange on RabbitMQ
func (r *RabbitMQAdapter) BindQueue(queueName, exchangeName, routingKey string) error {
	return r.ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
}

// Publish publishes a message to an exchange on RabbitMQ
func (r *RabbitMQAdapter) Publish(exchange, routingKey string, body []byte) error {
	return r.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func main() {
	// Initialize RabbitMQAdapter
	amqpURL := "amqp://guest:guest@localhost:5672/"
	rabbitMQ, err := NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ adapter: %v", err)
	}
	defer rabbitMQ.conn.Close()
	defer rabbitMQ.ch.Close()

	// Declare exchange, queue, and bind them
	exchangeName := "example_topic_exchange"
	queueName := "example_topic_queue"
	routingKey := "example.*"

	err = rabbitMQ.DeclareExchange(exchangeName, amqp.ExchangeTopic)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	err = rabbitMQ.DeclareQueue(queueName)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = rabbitMQ.BindQueue(queueName, exchangeName, routingKey)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	// Publish a message to the exchange
	message := "Hello, RabbitMQ Topic Exchange!"
	err = rabbitMQ.Publish(exchangeName, routingKey, []byte(message))
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Message published: %s", message)
}
