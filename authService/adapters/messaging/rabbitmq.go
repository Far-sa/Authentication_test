package messaging

import (
	"auth-svc/internal/ports"
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// type RabbitMQConfig struct {
// 	Host     string
// 	User     string
// 	Password string
// 	Port     string
// 	//Url string
// }

type RabbitClient struct {
	config ports.Config
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func NewRabbitMQClient(config ports.Config) (RabbitClient, error) {

	cfg := config.GetBrokerConfig()

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %s", err)
		return RabbitClient{}, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening channel: %s", err)
		conn.Close() // Close the connection if channel opening fails
		return RabbitClient{}, fmt.Errorf("error opening channel: %w", err)
	}

	if ch == nil {
		log.Fatal("RabbitMQ channel is nil")
	}

	// Defer closing the connection to ensure it's closed even in case of errors
	if err := ch.Confirm(false); err != nil {
		return RabbitClient{}, nil
	}

	return RabbitClient{
		config: config,
		conn:   conn,
		ch:     ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) DeclareExchange(name, kind string) error {
	return rc.ch.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error) {
	q, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	if err != nil {
		return amqp.Queue{}, nil
	}
	return q, err
}

// * for binding exchange to queue
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

func (rc RabbitClient) Send(
	ctx context.Context,
	exchange, routingKey string,
	options amqp.Publishing,
) error {
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		routingKey,
		true,
		false,
		options,
	)
	if err != nil {
		return err
	}

	log.Println(confirmation.Wait())
	// confirmation.Wait()
	return nil
}
