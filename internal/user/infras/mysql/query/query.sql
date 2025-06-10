-- name: CreateUser :execresult
INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?);

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, username, email, password_hash FROM users WHERE email = ?;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, username, email, password_hash FROM users WHERE id = ?;