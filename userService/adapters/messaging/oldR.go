package messaging

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ! Producer
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbit() (*RabbitMQ, error) {
	rabbitMQURL := "amqp://guest:guest@rabbitmq:5672/"

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	err = ch.ExchangeDeclare(
		"topic_exchange", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	body := "New user registered: John Doe"
	err = ch.Publish(
		"topic_exchange",   // exchange
		"registration.new", // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent %s", body)

	return &RabbitMQ{conn: conn, ch: ch}, nil

}
