package messaging

// import (
// 	"auth-svc/internal/ports"
// 	"fmt"

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

// // func (rc *RabbitClient) Consume(queueName, consumer string, autoAck bool, options amqp091.Table, callback func(amqp.Delivery) error) error {
// // 	ch, err := rc.GetChannel()
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer ch.Close() // Close the channel after use

// // 	if err := ch.Qos(50, 0, false); err != nil {
// // 		return err
// // 	}

// // 	msgs, err := ch.Consume(
// // 		queueName,
// // 		"",      // Consumer tag (can be left empty)
// // 		false,   // Auto-ack (set to false for manual ack)
// // 		false,   // Exclusive (only this consumer can access the queue)
// // 		false,   // No local (only deliver to this server)
// // 		false,   // No wait
// // 		options, // Arguments
// // 	)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	fmt.Println("Starting to consume stream")
// // 	for event := range msgs {
// // 		fmt.Printf("Event: %s\n", event.CorrelationId)
// // 		fmt.Printf("Headers : %v\n", event.Headers)
// // 		// Payload
// // 		fmt.Printf("Data : %v\n", string(event.Body))

// // 		if err := callback(event); err != nil {
// // 			fmt.Printf("Error processing message: %v\n", err)
// // 			// Optional: You can potentially re-queue the message here
// // 		}
// // 	}

// // 	return nil
// // }

// //!!

// func (rc *RabbitClient) Consume(queueName, consumer string, autoAck bool, options amqp091.Table) (<-chan amqp.Delivery, error) {
// 	ch, err := rc.GetChannel()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer ch.Close() // Close the channel after use

// 	if err := ch.Qos(50, 0, false); err != nil {
// 		panic(err)
// 	}

// 	msgs, err := ch.Consume(
// 		queueName,
// 		"",      // Consumer tag (can be left empty)
// 		false,   // Auto-ack (set to false for manual ack)
// 		false,   // Exclusive (only this consumer can access the queue)
// 		false,   // No local (only deliver to this server)
// 		false,   // No wait
// 		options, // Arguments
// 	)
// 	if err != nil {
// 		return msgs, err
// 	}

// 	fmt.Println("Starting to consume stream")
// 	for event := range msgs {
// 		fmt.Printf("Event: %s\n", event.CorrelationId)
// 		fmt.Printf("Headers : %v\n", event.Headers)
// 		// Payload
// 		fmt.Printf("Data : %v\n", string(event.Body))
// 	}

// 	//!--->
// 	// go func() {

// 	// 	for {
// 	// 		select {
// 	// 		case <-ctx.Done():
// 	// 			fmt.Println("Consume cancelled due to context")
// 	// 			return
// 	// 		case msg := <-msgs:
// 	// 			// Process the message and handle ack
// 	// 			if err := callback(msg); err != nil {
// 	// 				fmt.Printf("Error processing message: %v\n", err)
// 	// 				// Optional: You can potentially re-queue the message here
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }()

// 	return msgs, nil
// }
