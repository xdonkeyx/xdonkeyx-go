package rabbit

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"xdonkeyx.com/sample/common"
)

type RabbitConfig struct {
	amqpServerURL string

	// RabbitMQ connection
	connectRabbitMQ *amqp.Connection

	// RabbitMQ channel
	channelRabbitMQ *amqp.Channel
}

func StartRabbit() *RabbitConfig {
	var err error

	config := &RabbitConfig{}

	// Define RabbitMQ server URL.
	// amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	config.amqpServerURL = "amqp://guest:guest@localhost:5672/"

	// Create a new RabbitMQ connection.
	//connectRabbitMQ, err = amqp.Dial(amqpServerURL)
	config.connectRabbitMQ, err = amqp.DialConfig(config.amqpServerURL, amqp.Config{
		Heartbeat: 5 * time.Second,
	})
	common.FailOnError(err, "Failed to connect rabbitMQ")
	//defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already established.
	config.channelRabbitMQ, err = config.connectRabbitMQ.Channel()
	common.FailOnError(err, "Failed to open a channel")
	//defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = config.channelRabbitMQ.QueueDeclare(
		"DriverEvent", // queue name
		true,          // durable
		false,         // auto delete
		false,         // exclusive
		false,         // no wait
		nil,           // arguments
	)
	common.FailOnError(err, "Failed to declare a queue")

	return config
}

// send messages to RabbitMQ
func SendMessage(config *RabbitConfig, msg string) error {
	// Create a message to publish.
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	}

	fmt.Println(config)

	// Attempt to publish a message to the queue.
	if err := config.channelRabbitMQ.Publish(
		"",            // exchange
		"DriverEvent", // queue name
		false,         // mandatory
		false,         // immediate
		message,       // message to publish
	); err != nil {
		return err
	}

	return nil
}

func RegisterConsumer(config *RabbitConfig) {
	msgs, err := config.channelRabbitMQ.Consume(
		"DriverEvent", // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	common.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
