-- name: AssignPermissionsToUser :exec
INSERT INTO user_permissions (user_id, permission_id, granter_id, expires_at, created_at)
SELECT ?, ?, ?, ?, NOW()
FROM permissions
WHERE permissions.id = ? AND permissions.deleted_at IS NULL
ON DUPLICATE KEY UPDATE
                     granter_id = ?,
                     expires_at = ?,
                     deleted_at = NULL;

-- name: ListPermissionsForUser :many
SELECT p.id, p.name, p.created_at, COALESCE(p.deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM permissions p
         JOIN user_permissions up ON p.id = up.permission_id
WHERE up.user_id = ? AND up.deleted_at IS NULL AND p.deleted_at IS NULL;