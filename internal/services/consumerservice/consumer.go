package consumerservice

import (
	"encoding/json"
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/packages/utils/utlogger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	ConsumerService struct {
		ch    *amqp.Channel
		cfg   *configs.Config
		smail mailservice.IMailService
	}

	IConsumerService interface {
		ConsumeMail()
	}
)

func NewConsumerService(
	ch *amqp.Channel,
	cfg *configs.Config,
	smail mailservice.IMailService,
) *ConsumerService {
	return &ConsumerService{
		ch:    ch,
		cfg:   cfg,
		smail: smail,
	}
}

func (s *ConsumerService) ConsumeTask() {
	s.ConsumeMail()
}

func (s *ConsumerService) ConsumeMail() {
	var (
		qname = s.cfg.Queue.QueueMail.QueueName
	)
	// Listen to Queue
	messages, err := s.ch.Consume(
		qname, // queue
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)

	utlogger.Fatal(err)

	go func() {
		for d := range messages {
			utlogger.Info(fmt.Sprintf("Received a message: %s", d.Body))

			var (
				data requests.SendEmail
			)

			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				utlogger.Error(fmt.Errorf("Unable to unmarshal message: %w", err))
				return
			}

			utlogger.Info(fmt.Sprintf("Reference data mail: %v", data.Data))

			err = s.smail.SendEmail(data)
			if err != nil {
				utlogger.Error(fmt.Errorf("Unable to send email: %v", err))
			}

			ok := err == nil
			utlogger.Info(fmt.Sprintf("Send mail ok: %v", ok))
		}
	}()

	utlogger.Info(fmt.Sprintf("Service for %s is running, waiting for messages from queue!", qname))
}
