package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/usecases/user_role"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
)

type Consumer struct {
	urUC      user_role.UserRoleUseCase
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

func NewConsumer(urUC user_role.UserRoleUseCase) (*Consumer, error) {
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

	// Declare compensating queue
	_, err = channel.QueueDeclare(
		"user-compensation",
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
		var userID string
		var role string

		switch msg.Type {
		case "UserRegistered":
			var payload UserRegisteredPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				c.logger.Errorf("Failed to unmarshal UserRegistered event: %v", err)
				continue
			}
			userID = payload.UserID
			role = constants.RoleUser

		case "AdminUserCreated":
			var payload AdminUserCreatedPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				c.logger.Errorf("Failed to unmarshal AdminUserCreated event: %v", err)
				continue
			}
			userID = payload.UserID
			role = constants.RoleAdmin

		default:
			c.logger.Errorf("Unknown event type: '%s', full message: %s", msg.Type, string(msg.Body))
			continue
		}

		err = c.urUC.AssignRolesToUser(ctx, []domain.UserRole{
			{
				UserID:   userID,
				RoleName: role,
			},
		})
		if err != nil {
			c.logger.Errorf("Failed to assign role %s to user %s: %v", role, userID, err)
			// Publish compensating event
			compPayload := struct {
				UserID string `json:"user_id"`
				Reason string `json:"reason"`
			}{
				UserID: userID,
				Reason: "Failed to assign role",
			}
			compBytes, err := json.Marshal(compPayload)
			if err != nil {
				c.logger.Errorf("Failed to marshal compensating event for user %s: %v", userID, err)
				continue
			}
			err = c.channel.PublishWithContext(
				ctx,
				"",                  // exchange
				"user-compensation", // routing key
				false,               // mandatory
				false,               // immediate
				amqp091.Publishing{
					ContentType: "application/json",
					Body:        compBytes,
					MessageId:   "comp-" + userID,
				},
			)
			if err != nil {
				c.logger.Errorf("Failed to publish compensating event for user %s: %v", userID, err)
			}
			continue
		}

		c.logger.Infof("Assigned role %s to user %s", role, userID)
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
