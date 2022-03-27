package rabbit

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
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
	if err != nil {
		panic(err)
	}
	//defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	config.channelRabbitMQ, err = config.connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	//defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = config.channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

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
