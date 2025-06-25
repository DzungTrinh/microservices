-- name: CreatePermission :execresult
INSERT INTO permissions (id, name, created_at)
    VALUES (?, ?, NOW());

-- name: DeletePermission :exec
UPDATE permissions
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: ListPermissions :many
SELECT id, name, created_at, deleted_at
FROM permissions;
