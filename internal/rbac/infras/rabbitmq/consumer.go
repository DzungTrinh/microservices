package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/internal/rbac/usecases/user_role"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
	"time"
)

type Consumer struct {
	urUC      user_role.UserRoleUseCase
	rUC       role.RoleUseCase
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	queueName string
	logger    *logger.LoggerService
}

type UserRegisteredPayload struct {
	UserID string `json:"user_id"`
}

type AdminUserCreatedPayload struct {
	UserID string `json:"user_id"`
}

func NewConsumer(urUC user_role.UserRoleUseCase, rUC role.RoleUseCase) (*Consumer, error) {
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
		urUC:      urUC,
		rUC:       rUC,
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
		var userID string
		var role string

		switch msg.Type {
		case "UserRegistered":
			var payload UserRegisteredPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				c.logger.Errorf("Failed to unmarshal UserRegistered event: %v", err)
				msg.Nack(false, false)
				continue
			}
			userID = payload.UserID
			role = constants.RoleUser

		case "AdminUserCreated":
			var payload AdminUserCreatedPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				c.logger.Errorf("Failed to unmarshal AdminUserCreated event: %v", err)
				msg.Nack(false, false)
				continue
			}
			userID = payload.UserID
			role = constants.RoleAdmin

		default:
			c.logger.Errorf("Unknown event type: '%s', full message: %s", msg.Type, string(msg.Body))
			msg.Ack(true) // Acknowledge to avoid requeueing unknown types
			continue
		}

		// Fetch role ID by name
		res, err := c.rUC.GetRoleByName(ctx, role)
		if err != nil {
			logger.GetInstance().Errorf("Failed to get role %s: %v", role, err)
			return err
		}

		// Retry logic
		const maxRetries = 5
		const retryDelay = 1 * time.Second
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = c.urUC.AssignRolesToUser(ctx, []domain.UserRole{
				{
					UserID: userID,
					RoleID: res.ID,
				},
			})
			if err == nil {
				c.logger.Infof("Assigned role %s to user %s", role, userID)
				msg.Ack(true) // Acknowledge successful processing
				break
			}

			c.logger.Errorf("Attempt %d/%d: Failed to assign role %s to user %s: %v", attempt, maxRetries, role, userID, err)
			if attempt == maxRetries {
				c.logger.Errorf("Exhausted retries for assigning role %s to user %s", role, userID)
				msg.Ack(true) // Acknowledge to prevent infinite retries
				break
			}
			time.Sleep(retryDelay)
		}
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
