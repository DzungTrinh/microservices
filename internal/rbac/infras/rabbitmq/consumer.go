package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
)

type Consumer struct {
	client    rbacv1.RBACServiceClient
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	queueName string
	logger    *logger.LoggerService
}

type UserRegisteredPayload struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func NewConsumer(client rbacv1.RBACServiceClient) (*Consumer, error) {
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

	return &Consumer{
		client:    client,
		conn:      conn,
		channel:   channel,
		queueName: queueName,
		logger:    logger.GetInstance(),
	}, nil
}

func (c *Consumer) ConsumeEvents(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		true,        // autoAck
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
		var payload UserRegisteredPayload
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			c.logger.Errorf("Failed to unmarshal event: %v", err)
			continue
		}

		_, err = c.client.AssignRolesToUser(ctx, &rbacv1.AssignRolesToUserRequest{
			UserId:  payload.UserID,
			RoleIds: []string{payload.Role},
		})
		if err != nil {
			c.logger.Errorf("Failed to assign role %s to user %s: %v", payload.Role, payload.UserID, err)
			// TODO: Publish compensating event (e.g., UserRegistrationFailed)
			continue
		}

		c.logger.Infof("Assigned role %s to user %s", payload.Role, payload.UserID)
	}
	return nil
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		c.logger.Errorf("Failed to close RabbitMQ channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		c.logger.Errorf("Failed to close RabbitMQ connection: %v", err)
	}
	return nil
}
