package messaging

// import (
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// // ! Consumer
// type RabbitMQ struct {
// 	conn *amqp.Connection
// 	ch   *amqp.Channel
// }

// func NewRabbit() (*RabbitMQ, error) {
// 	rabbitMQURL := "amqp://guest:guest@rabbitmq:5672/"

// 	conn, err := amqp.Dial(rabbitMQURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		conn.Close()
// 		return nil, err
// 	}

// 	err = ch.ExchangeDeclare(
// 		"topic_exchange", // name
// 		"topic",          // type
// 		true,             // durable
// 		false,            // auto-deleted
// 		false,            // internal
// 		false,            // no-wait
// 		nil,              // arguments
// 	)
// 	if err != nil {
// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}

// 	q, err := ch.QueueDeclare(
// 		"registration_queue", // name
// 		true,                 // durable
// 		false,                // delete when unused
// 		false,                // exclusive
// 		false,                // no-wait
// 		nil,                  // arguments
// 	)
// 	if err != nil {
// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}

// 	err = ch.QueueBind(
// 		q.Name,           // queue name
// 		"registration.*", // routing key
// 		"topic_exchange", // exchange
// 		false,            // no-wait
// 		nil,              // arguments
// 	)
// 	if err != nil {
// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}
// 	msgs, err := ch.Consume(
// 		"registration_queue", // queue
// 		"",                   // consumer
// 		false,                // auto-ack (set to false for manual ack)
// 		false,                // exclusive
// 		false,                // no-local
// 		false,                // no-wait
// 		nil,                  // args
// 	)
// 	if err != nil {
// 		conn.Close()
// 		ch.Close()
// 		return nil, err
// 	}

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		for d := range msgs {
// 			log.Printf("Received a message: %s", d.Body)
// 			// Acknowledge the message
// 			d.Ack(false)
// 		}
// 	}()

// 	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
// 	<-signals

// 	return &RabbitMQ{conn: conn, ch: ch}, nil
// }
