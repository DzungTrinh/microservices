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
