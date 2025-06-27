-- name: CreateUser :exec
INSERT INTO users (id, email, username, email_verified, created_at, updated_at)
VALUES (?, ?, ?, ?, NOW(), NOW());

-- name: GetUserByEmail :one
SELECT id, email, username, email_verified, created_at, updated_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM users
WHERE email = ? AND deleted_at IS NULL;

-- name: GetUserById :one
SELECT id, email, username, email_verified, created_at, updated_at, COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00') AS deleted_at
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: GetAllUsers :many
SELECT id, email, username, email_verified, created_at, updated_at
FROM users
WHERE deleted_at IS NULL;