package producerservice

import (
	"context"
	"encoding/json"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	ProducerService struct {
		ch  *amqp.Channel
		cfg *configs.Config
		ctx context.Context
	}

	IProducerService interface {
		PublishEmail(message requests.SendEmail) error
	}
)

func NewProducerService(
	ch *amqp.Channel,
	cfg *configs.Config,
	ctx context.Context,
) *ProducerService {
	return &ProducerService{
		ch:  ch,
		cfg: cfg,
		ctx: ctx,
	}
}

func (s *ProducerService) PublishEmail(message requests.SendEmail) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = s.ch.PublishWithContext(
		s.ctx,
		s.cfg.Queue.Mail.ExchangeName, // exchange
		s.cfg.Queue.Mail.BindingKey,   // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonData,
		})
	if err != nil {
		return err
	}

	return nil
}
