-- name: CreateRole :execresult
INSERT INTO roles (id, name, built_in, created_at)
    VALUES (?, ?, ?, NOW());

-- name: GetRoleByID :one
SELECT id, name, built_in, created_at
FROM roles
WHERE id = ? AND deleted_at IS NULL;

-- name: ListRoles :many
SELECT id, name, built_in, created_at, deleted_at
FROM roles
WHERE deleted_at IS NULL;

-- name: UpdateRole :exec
UPDATE roles
SET name = ?, built_in = ?, updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: DeleteRole :exec
UPDATE roles
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;