package queue

import (
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

func (r *Rabbit) Consume() error {
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

	for msg := range incomingMessage {
		r.Channel.Ack(msg.DeliveryTag, false)
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}
