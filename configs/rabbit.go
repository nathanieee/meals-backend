package configs

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rmq *Queue) Init() (*amqp.Channel, func()) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rmq.Username, rmq.Password, rmq.Host, rmq.Port)
	conn, err := amqp.Dial(url)
	FailOnError(err, "Failed to connect to Queue")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	return ch, func() {
		conn.Close()
		ch.Close()
	}
}

func (rmq *Queue) SetupRabbitMQ(ch *amqp.Channel, cfg *Config) {
	rmq.SetupMailQueue(ch, cfg.Queue)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: %s", msg, err))
	}
}

func (rmq *Queue) SetupMailQueue(ch *amqp.Channel, cfg Queue) {
	exchangeName := cfg.Mail.ExchangeName
	exchangeType := cfg.Mail.ExchangeType
	queueName := cfg.Mail.QueueName
	bindingKey := cfg.Mail.BindingKey

	// Declare Exchange
	err := ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare an exchange")

	// Declare Queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	// Binding Exchange and Queue
	err = ch.QueueBind(
		q.Name,       // queue name
		bindingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	FailOnError(err, "Failed to bind a queue")
}
