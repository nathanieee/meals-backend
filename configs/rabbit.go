package configs

import (
	"fmt"
	"project-skbackend/packages/utils/utlogger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rmq *Queue) Init() (*amqp.Channel, func()) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rmq.Username, rmq.Password, rmq.Host, rmq.Port)
	conn, err := amqp.Dial(url)
	utlogger.Fatal(err)

	ch, err := conn.Channel()
	utlogger.Fatal(err)

	return ch, func() {
		conn.Close()
		ch.Close()
	}
}

func (rmq *Queue) SetupRabbitMQ(ch *amqp.Channel, cfg *Config) {
	rmq.SetupMailQueue(ch, cfg.Queue)
}

func (rmq *Queue) SetupMailQueue(ch *amqp.Channel, cfg Queue) {
	xname := cfg.Mail.ExchangeName
	xtype := cfg.Mail.ExchangeType
	qname := cfg.Mail.QueueName
	bkey := cfg.Mail.BindingKey

	// Declare Exchange
	err := ch.ExchangeDeclare(
		xname, // name
		xtype, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	utlogger.Fatal(err)

	// Declare Queue
	q, err := ch.QueueDeclare(
		qname, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utlogger.Fatal(err)

	// Binding Exchange and Queue
	err = ch.QueueBind(
		q.Name, // queue name
		bkey,   // routing key
		xname,  // exchange
		false,
		nil)
	utlogger.Fatal(err)
}
