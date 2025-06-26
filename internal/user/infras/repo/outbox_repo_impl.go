package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/domain/repo"
	"microservices/user-management/internal/user/infras/mysql"
)

type outboxRepository struct {
	db *sql.DB
	q  *mysql.Queries
}

func NewOutboxRepository(db *sql.DB) repo.OutboxRepository {
	return &outboxRepository{
		db: db,
		q:  mysql.New(db),
	}
}

func (r *outboxRepository) InsertEvent(ctx context.Context, e *domain.OutboxEvent) error {
	return r.q.InsertOutboxEvent(ctx, mysql.InsertOutboxEventParams{
		AggregateType: e.AggregateType,
		AggregateID:   e.AggregateID,
		Type:          e.Type,
		Payload:       json.RawMessage(e.Payload),
		Status:        e.Status, // usually "pending"
	})
}

func (r *outboxRepository) GetPendingEvents(ctx context.Context, limit int32) ([]*domain.OutboxEvent, error) {
	dbEvents, err := r.q.GetPendingOutboxEvents(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.OutboxEvent, 0, len(dbEvents))
	for _, db := range dbEvents {
		e := &domain.OutboxEvent{
			ID:            db.ID,
			AggregateType: db.AggregateType,
			AggregateID:   db.AggregateID,
			Type:          db.Type,
			Payload:       string(db.Payload),
			Status:        db.Status,
			CreatedAt:     db.CreatedAt,
		}
		e.ProcessedAt = db.ProcessedAt
		result = append(result, e)
	}

	return result, nil
}

func (r *outboxRepository) MarkEventProcessed(ctx context.Context, id int64) error {
	return r.q.MarkOutboxEventProcessed(ctx, id)
}
