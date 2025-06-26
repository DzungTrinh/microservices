-- name: CreatePermission :exec
INSERT INTO permissions (id, name, created_at)
VALUES (?, ?, NOW());

-- name: DeletePermission :exec
UPDATE permissions
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: ListPermissions :many
SELECT id, name, created_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM permissions
WHERE deleted_at IS NULL;
