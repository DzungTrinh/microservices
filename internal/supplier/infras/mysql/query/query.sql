-- name: CreateUser :execresult
INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?);

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, username, email, password, role FROM users WHERE email = ?;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, username, email, password, role FROM users WHERE id = ?;

-- name: GetAllUsers :many
SELECT id, created_at, updated_at, username, email, password, role FROM users;