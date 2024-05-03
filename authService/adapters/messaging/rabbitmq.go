package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	//Url string
}

type RabbitClient struct {
	config RabbitMQConfig
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func NewRabbitMQClient(config RabbitMQConfig) (RabbitClient, error) {

	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/", config.User, config.Password, config.Host, config.Port))
	if err != nil {
		return RabbitClient{}, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, fmt.Errorf("error opening channel: %w", err)
	}
	// if err := ch.Confirm(false); err != nil {
	// 	return RabbitClient{}, nil
	// }

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

// * for bindig echange to queue
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
