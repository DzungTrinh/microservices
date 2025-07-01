package events

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

// EventHandler defines the interface for handling RabbitMQ events.
type EventHandler interface {
	Handle(ctx context.Context, msg amqp091.Delivery) error
}

// EventHandlerFunc is a function type that implements EventHandler.
type EventHandlerFunc func(ctx context.Context, msg amqp091.Delivery) error

func (f EventHandlerFunc) Handle(ctx context.Context, msg amqp091.Delivery) error {
	return f(ctx, msg)
}
