-- name: InsertOutboxEvent :exec
INSERT INTO outbox_events (
    aggregate_type,
    aggregate_id,
    type,
    payload,
    status,
    created_at
) VALUES (?, ?, ?, ?, ?, NOW());

-- name: GetPendingOutboxEvents :many
SELECT
    id,
    aggregate_type,
    aggregate_id,
    type,
    payload,
    status,
    created_at,
    processed_at
FROM outbox_events
WHERE status = 'pending'
ORDER BY created_at
LIMIT ?;

-- name: MarkOutboxEventProcessed :exec
UPDATE outbox_events
SET status = 'processed', processed_at = NOW()
WHERE id = ?;