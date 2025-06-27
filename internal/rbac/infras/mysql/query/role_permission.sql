-- name: AssignPermissionsToRole :exec
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT ?, ?, NOW();

-- name: ListPermissionsForRole :many
SELECT p.id, p.name, p.created_at, COALESCE(p.deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM permissions p
         JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = ? AND rp.deleted_at IS NULL AND p.deleted_at IS NULL;