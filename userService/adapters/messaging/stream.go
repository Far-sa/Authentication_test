package messaging

// import (
// 	"context"
// 	"fmt"
// 	"user-svc/ports"

// 	"github.com/rabbitmq/amqp091-go"
// 	amqp "github.com/rabbitmq/amqp091-go"
// 	//"github.com/rabbitmq/amqp091-go/amqp"
// 	// "github.com/streadway/amqp" // Import the streadway/amqp library
// )

// type RabbitClient struct {
// 	config ports.Config
// 	conn   *amqp.Connection
// }

// func NewRabbitCL(config ports.Config) (*RabbitClient, error) {
// 	cfg := config.GetBrokerConfig()
// 	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)

// 	conn, err := amqp.Dial(dsn)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
// 	}

// 	return &RabbitClient{config: config, conn: conn}, nil
// }

// func (rc *RabbitClient) GetChannel() (*amqp.Channel, error) {
// 	ch, err := rc.conn.Channel()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ch, nil
// }

// func (rc *RabbitClient) Close() error {
// 	return rc.conn.Close() // Close the underlying connection
// }

// func (rc *RabbitClient) CreateQueue(name string) (string, error) {
// 	ch, err := rc.GetChannel()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer ch.Close() // Close the channel after use

// 	queue, err := ch.QueueDeclare("events", true, false, false, true, amqp091.Table{
// 		"x-queue-type":                    "stream",
// 		"x-stream-max-segment-size-bytes": 30000,  // EACH SEGMENT FILE IS ALLOWED 0.03 MB
// 		"x-max-length-bytes":              150000, // TOTAL STREAM SIZE IS 0.15 MB
// 		// "x-max-age" : "10s"
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	return queue.Name, nil
// }

// func (rc *RabbitClient) PublishMessage(ctx context.Context, exchangeName string, routingKey string, options amqp.Publishing) error {
// 	ch, err := rc.GetChannel()
// 	if err != nil {
// 		return err
// 	}
// 	defer ch.Close() // Close the channel after use

// 	// body, err := json.Marshal(message) // Marshal the message to JSON
// 	// if err != nil {
// 	// 	return fmt.Errorf("failed to marshal message: %w", err)
// 	// }

// 	err = ch.PublishWithContext(
// 		ctx,
// 		exchangeName, // Name of the exchange
// 		routingKey,   // Routing key for message
// 		false,        // Mandatory (if true, message is rejected if no queue is bound)
// 		false,        // Immediate (if true, delivery happens now, or fails)
// 		options,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
