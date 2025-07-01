package events

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/pkg/logger"
	"time"
)

type Consumer struct {
	conn       *amqp091.Channel
	queueName  string
	logger     *logger.LoggerService
	handlers   map[string]EventHandler
	maxRetries int
	retryDelay time.Duration
}

func NewConsumer() (*Consumer, error) {
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

	// Declare main queue
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

	// Declare dead-letter queue
	dlqName := queueName + ".dlq"
	_, err = channel.QueueDeclare(
		dlqName,
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

	return &Consumer{
		conn:       channel,
		queueName:  queueName,
		logger:     logger.GetInstance(),
		handlers:   make(map[string]EventHandler),
		maxRetries: 5,
		retryDelay: 1 * time.Second,
	}, nil
}

// RegisterHandler adds a handler for a specific event type.
func (c *Consumer) RegisterHandler(eventType string, handler EventHandler) {
	c.handlers[eventType] = handler
}

// ConsumeEvents processes messages from the queue.
func (c *Consumer) ConsumeEvents(ctx context.Context) error {
	msgs, err := c.conn.Consume(
		c.queueName, // queue
		"",          // consumer
		false,       // autoAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // args
	)
	if err != nil {
		c.logger.Errorf("Failed to start consuming messages: %v", err)
		return err
	}

	for msg := range msgs {
		handler, exists := c.handlers[msg.Type]
		if !exists {
			c.logger.Errorf("Unknown event type: '%s', full message: %s", msg.Type, string(msg.Body))
			if err := msg.Ack(true); err != nil {
				c.logger.Errorf("Failed to ack unknown event: %v", err)
			}
			continue
		}

		err := handler.Handle(ctx, msg)
		if err != nil {
			c.logger.Errorf("Failed to process event %s: %v", msg.Type, err)
			// Send to DLQ
			if err := c.conn.PublishWithContext(ctx,
				"",                 // exchange
				c.queueName+".dlq", // routing key
				false,              // mandatory
				false,              // immediate
				amqp091.Publishing{
					ContentType: "application/json",
					Type:        msg.Type,
					Body:        msg.Body,
				}); err != nil {
				c.logger.Errorf("Failed to publish to DLQ: %v", err)
			}
			if err := msg.Ack(true); err != nil {
				c.logger.Errorf("Failed to ack failed event: %v", err)
			}
			continue
		}

		if err := msg.Ack(true); err != nil {
			c.logger.Errorf("Failed to ack event %s: %v", msg.Type, err)
		}
	}
	return nil
}

// Close shuts down the consumer.
func (c *Consumer) Close() error {
	if err := c.conn.Close(); err != nil {
		c.logger.Errorf("Failed to close RabbitMQ channel: %v", err)
		return err
	}
	return nil
}
