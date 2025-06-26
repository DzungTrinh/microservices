package rabbitmq

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"microservices/user-management/cmd/user/config"
	"microservices/user-management/internal/user/domain/repo"
	"microservices/user-management/internal/user/dto"
	"microservices/user-management/pkg/logger"
	"time"
)

type Publisher struct {
	repo      repo.OutboxRepository
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	queueName string
	logger    *logger.LoggerService
}

func NewPublisher(repo repo.OutboxRepository) (*Publisher, error) {
	amqpURL := config.GetInstance().RabbitmqUrl
	queueName := config.GetInstance().RabbitmqQueue

	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &Publisher{
		repo:      repo,
		conn:      conn,
		channel:   channel,
		queueName: queueName,
		logger:    logger.GetInstance(),
	}, nil
}

func (p *Publisher) PublishEvents(ctx context.Context) error {
	events, err := p.repo.GetPendingEvents(ctx, 100)
	if err != nil {
		p.logger.Errorf("Failed to get pending events: %v", err)
		return err
	}

	for _, event := range events {
		if event.Status != dto.OutboxPending {
			continue
		}

		err = p.channel.PublishWithContext(
			ctx,
			"",          // exchange
			p.queueName, // routing key
			false,       // mandatory
			false,       // immediate
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        []byte(event.Payload),
				MessageId:   string(event.ID),
			},
		)
		if err != nil {
			p.logger.Errorf("Failed to publish event %d: %v", event.ID, err)
			continue
		}

		err = p.repo.MarkEventProcessed(ctx, event.ID)
		if err != nil {
			p.logger.Errorf("Failed to mark event %d as processed: %v", event.ID, err)
			continue
		}
		p.logger.Infof("Published event %d for user %s", event.ID, event.AggregateID)
	}
	return nil
}

func (p *Publisher) StartOutboxWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Outbox worker shutting down")
			return
		case <-ticker.C:
			err := p.PublishEvents(ctx)
			if err != nil {
				p.logger.Errorf("Worker error: %v", err)
			}
		}
	}
}

func (p *Publisher) Close() error {
	if err := p.channel.Close(); err != nil {
		p.logger.Errorf("Failed to close RabbitMQ channel: %v", err)
	}
	if err := p.conn.Close(); err != nil {
		p.logger.Errorf("Failed to close RabbitMQ connection: %v", err)
	}
	return nil
}
