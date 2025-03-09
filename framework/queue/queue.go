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
