package domain

import "time"

type OutboxEvent struct {
	ID            int64
	AggregateType string
	AggregateID   string
	Type          string
	Payload       string // assuming JSON string
	Status        string
	CreatedAt     time.Time
	ProcessedAt   *time.Time
}
