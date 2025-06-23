-- name: CreateUser :execresult
INSERT INTO users (id, username, email, password, created_at, updated_at)
VALUES (?, ?, ?, ?, NOW(), NOW());

-- name: CreateUserRole :exec
INSERT INTO user_roles (user_id, role_id) VALUES (?, ?);

-- name: GetRoleIDByName :one
SELECT id FROM roles WHERE name = ?;

-- name: GetUserByEmail :one
SELECT u.id, u.created_at, u.updated_at, u.username, u.email, u.password,
       GROUP_CONCAT(r.name) AS roles
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.email = ?
GROUP BY u.id;

-- name: GetUserByID :one
SELECT u.id, u.created_at, u.updated_at, u.username, u.email, u.password,
       GROUP_CONCAT(r.name) AS roles
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.id = ?
GROUP BY u.id;

-- name: GetAllUsers :many
SELECT u.id, u.created_at, u.updated_at, u.username, u.email, u.password,
       GROUP_CONCAT(r.name) AS roles
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
GROUP BY u.id;

-- name: DeleteUserRoles :exec
DELETE FROM user_roles WHERE user_id = ?;

-- name: GetUserRoles :many
SELECT GROUP_CONCAT(r.name) AS roles
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.id = ?
GROUP BY u.id;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token, user_agent, ip_address, expires_at, revoked)
VALUES (?, ?, ?, ?, ?, ?, 0);

-- name: GetRefreshToken :one
SELECT id, user_id, token, user_agent, ip_address, created_at, expires_at, revoked
FROM refresh_tokens
WHERE token = ? AND expires_at > NOW() AND revoked = 0;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = 1
WHERE token = ?;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW() OR revoked = 1;