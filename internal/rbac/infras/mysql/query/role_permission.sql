-- name: AssignPermissionsToRole :exec
INSERT INTO role_permissions (role_id, perm_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE role_id = role_id;

-- name: ListPermissionsForRole :many
SELECT p.id, p.name, p.created_at
FROM permissions p
         JOIN role_permissions rp ON p.id = rp.perm_id
WHERE rp.role_id = ? AND p.deleted_at IS NULL;