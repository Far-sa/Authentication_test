// internal/adapter/messaging/rabbitmq.go

package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQClient(username, password, host, vhost string) (RabbitClient, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
	if err != nil {
		return RabbitClient{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	// if err := ch.Confirm(false); err != nil {
	// 	return RabbitClient{}, nil
	// }
	// TODO: create exchange for binding
	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

// ! Stream Publisher
// func (rc RabbitClient) PublishUserRegisteredEvent(ctx context.Context, data []byte) error {
// 	correlationID := uuid.NewString()

// 	err := rc.ch.PublishWithContext(ctx, "", "events", false, false, amqp091.Publishing{
// 		Body:          data,
// 		CorrelationId: correlationID,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }

func (rc RabbitClient) Publish(
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

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

// ! Create Queue
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

//! creatae connection for further processing
// func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
// 	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
// }
