package queue

import (
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

type Rabbit struct {
	User              string
	Password          string
	Host              string
	Port              int
	Vhost             string
	ConsumerQueueName string
	ConsumeName       string
	AutoAck           bool
	Args              amqp.Table
	Channel           *amqp.Channel
}

func NewRabbit() *Rabbit {
	rabbitArgs := amqp.Table{}
	rabbitArgs["x-dead-letter-exchange"] = os.Getenv("RABBIT_DLE")

	rabbitPort := os.Getenv("RABBIT_PORT")
	rabbitPortInt, err := strconv.Atoi(rabbitPort)

	if err != nil {
		rabbitPortInt = 5672
	}

	rabbit := Rabbit{
		User:              os.Getenv("RABBIT_USER"),
		Password:          os.Getenv("RABBIT_PASSWORD"),
		Host:              os.Getenv("RABBIT_HOST"),
		Port:              rabbitPortInt,
		Vhost:             os.Getenv("RABBIT_VHOST"),
		ConsumerQueueName: os.Getenv("RABBIT_CONSUMER_QUEUE_NAME"),
		ConsumeName:       os.Getenv("RABBIT_CONSUMER_NAME"),
		AutoAck:           true,
		Args:              rabbitArgs,
		Channel:           nil,
	}

	return &rabbit
}

func (r *Rabbit) Consume(messageChannel chan amqp.Delivery) error {
	q, err := r.Channel.QueueDeclare(
		r.ConsumerQueueName,
		true,
		false,
		false,
		false,
		r.Args,
	)

	failOnError(err, "Failed to declare a queue")

	incomingMessage, err := r.Channel.Consume(
		q.Name,
		r.ConsumeName,
		r.AutoAck,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer")

	go func() {
		for message := range incomingMessage {
			log.Println("Received a message: ", string(message.Body))
			messageChannel <- message
		}

		log.Println("Consumer closed")
		close(messageChannel)
	}()

	return nil
}

func (r *Rabbit) Notify(message string, contentType string, exchange string, routingKey string) error {
	err := r.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		})

	if err != nil {
		return err
	}
	return nil
}
func failOnError(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}
