-- name: AssignPermissionsToUser :exec
INSERT INTO user_permissions (user_id, perm_id, granter_id, expires_at, created_at)
VALUES (?, ?, ?, ?, NOW())
ON DUPLICATE KEY UPDATE granter_id = ?, expires_at = ?;

-- name: ListPermissionsForUser :many
SELECT p.id, p.name, up.created_at, up.expires_at, up.granter_id
FROM permissions p
         JOIN user_permissions up ON p.id = up.permission_id
WHERE up.user_id = ?;