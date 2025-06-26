package repo

import (
	"context"
	"microservices/user-management/internal/user/domain"
)

type OutboxRepository interface {
	InsertEvent(ctx context.Context, event *domain.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int32) ([]*domain.OutboxEvent, error)
	MarkEventProcessed(ctx context.Context, id int64) error
}
