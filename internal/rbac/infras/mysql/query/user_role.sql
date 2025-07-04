-- name: AssignRolesToUser :exec
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT ?, ?, NOW()
FROM roles
WHERE roles.id = ? AND roles.deleted_at IS NULL
ON DUPLICATE KEY UPDATE user_id = user_id;

-- name: ListRolesForUser :many
SELECT r.id, r.name, r.built_in, r.created_at, COALESCE(r.deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM roles r
         JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = ? AND ur.deleted_at IS NULL AND r.deleted_at IS NULL;

-- name: SoftDeleteUserRole :exec
UPDATE user_roles
SET deleted_at = NOW()
WHERE user_id = ? AND role_id = ? AND deleted_at IS NULL;
