-- name: AssignRolesToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE user_id = user_id;

-- name: ListRolesForUser :many
SELECT r.id, r.name, r.built_in, r.created_at
FROM roles r
         JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = ?;