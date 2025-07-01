-- name: CreateRole :exec
INSERT INTO roles (id, name, built_in, created_at)
VALUES (?, ?, ?, NOW());

-- name: GetRoleByName :one
SELECT id, name, built_in, created_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM roles
WHERE name = ? AND deleted_at IS NULL;

-- name: ListRoles :many
SELECT id, name, built_in, created_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM roles
WHERE deleted_at IS NULL;

-- name: UpdateRole :exec
UPDATE roles
SET name = ?, built_in = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: DeleteRole :exec
UPDATE roles
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;